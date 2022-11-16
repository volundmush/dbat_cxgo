package main

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const MAX_HOUSES = 1000
const MAX_GUESTS = 10
const HOUSE_PRIVATE = 0
const HOUSE_GOD = 1
const HOUSE_CLAN = 2
const HOUSE_UNOWNED = 3
const NUM_HOUSE_TYPES = 4
const HOUSE_NOGUESTS = 1
const HOUSE_FREE = 2
const HOUSE_NOIMMS = 4
const HOUSE_IMPONLY = 8
const HOUSE_RENTFREE = 16
const HOUSE_SAVENORENT = 32
const HOUSE_NOSAVE = 64
const HOUSE_NUM_FLAGS = 7

type house_control_rec struct {
	Vnum          room_vnum
	Atrium        room_vnum
	Exit_num      int16
	Built_on      libc.Time
	Mode          int
	Owner         int
	Num_of_guests int
	Guests        [10]int
	Last_payment  libc.Time
	Bitvector     int
	Builtby       int
	Spare2        int
	Spare3        int
	Spare4        int
	Spare5        int
	Spare6        int
	Spare7        int
}

var house_control [1000]house_control_rec
var num_of_houses int = 0

func House_get_filename(vnum room_vnum, filename *byte, maxlen uint64) int {
	if vnum == room_vnum(-1) {
		return 0
	}
	stdio.Snprintf(filename, int(maxlen), LIB_HOUSE, vnum)
	return 1
}
func House_save(obj *obj_data, fp *stdio.File, location int) int {
	var (
		tmp    *obj_data
		result int
	)
	if obj != nil {
		if OBJ_FLAGGED(obj, ITEM_NORENT) {
			obj = obj.Next_content
		}
	}
	if obj != nil {
		House_save(obj.Next_content, fp, location)
		House_save(obj.Contains, fp, int(MIN(0, int64(location))-1))
		result = Obj_to_store(obj, fp, location)
		if result == 0 {
			return 0
		}
		for tmp = obj.In_obj; tmp != nil; tmp = tmp.In_obj {
			tmp.Weight -= obj.Weight
		}
	}
	return 1
}
func House_restore_weight(obj *obj_data) {
	if obj != nil {
		House_restore_weight(obj.Contains)
		House_restore_weight(obj.Next_content)
		if obj.In_obj != nil {
			obj.In_obj.Weight += obj.Weight
		}
	}
}
func House_crashsave(vnum room_vnum) {
	var (
		rnum int
		buf  [64936]byte
		fp   *stdio.File
	)
	if (func() int {
		rnum = int(real_room(vnum))
		return rnum
	}()) == int(-1) {
		return
	}
	if House_get_filename(vnum, &buf[0], uint64(64936)) == 0 {
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fp
	}()) == nil {
		fmt.Println(libc.CString("SYSERR: Error saving house file"))
		return
	}
	if House_save(world[rnum].Contents, fp, 0) == 0 {
		fp.Close()
		return
	}
	fp.Close()
	House_restore_weight(world[rnum].Contents)
	REMOVE_BIT_AR(world[rnum].Room_flags[:], ROOM_HOUSE_CRASH)
}
func House_delete_file(vnum room_vnum) {
	var (
		filename [2048]byte
		fl       *stdio.File
	)
	if House_get_filename(vnum, &filename[0], uint64(2048)) == 0 {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&filename[0]), "rb")
		return fl
	}()) == nil {
		if libc.Errno != 2 {
			basic_mud_log(libc.CString("SYSERR: Error deleting house file #%d. (1): %s"), vnum, libc.StrError(libc.Errno))
		}
		return
	}
	fl.Close()
	if stdio.Remove(libc.GoString(&filename[0])) < 0 {
		basic_mud_log(libc.CString("SYSERR: Error deleting house file #%d. (2): %s"), vnum, libc.StrError(libc.Errno))
	}
}
func find_house(vnum room_vnum) int {
	var i int
	for i = 0; i < num_of_houses; i++ {
		if house_control[i].Vnum == vnum {
			return i
		}
	}
	return -1
}
func House_save_control() {
	var fl *stdio.File
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "wb")
		return fl
	}()) == nil {
		fmt.Println(libc.CString("SYSERR: Unable to open house control file."))
		return
	}
	fl.WriteN((*byte)(unsafe.Pointer(&house_control[0])), int(unsafe.Sizeof(house_control_rec{})), num_of_houses)
	fl.Close()
}
func House_boot() {
	var (
		temp_house house_control_rec
		real_house room_rnum
		fl         *stdio.File
	)
	libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(&house_control[0]))), 0, int(MAX_HOUSES*unsafe.Sizeof(house_control_rec{})))
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "rb")
		return fl
	}()) == nil {
		if libc.Errno == 2 {
			basic_mud_log(libc.CString("   No houses to load. File '%s' does not exist."), LIB_ETC)
		} else {
			fmt.Println(libc.CString("SYSERR: etc/hcontrol"))
		}
		return
	}
	for int(fl.IsEOF()) == 0 && num_of_houses < MAX_HOUSES {
		fl.ReadN((*byte)(unsafe.Pointer(&temp_house)), int(unsafe.Sizeof(house_control_rec{})), 1)
		if int(fl.IsEOF()) != 0 {
			break
		}
		if get_name_by_id(temp_house.Owner) == nil {
			continue
		}
		if (func() room_rnum {
			real_house = real_room(temp_house.Vnum)
			return real_house
		}()) == room_rnum(-1) {
			continue
		}
		if find_house(temp_house.Vnum) != int(-1) {
			continue
		}
		house_control[func() int {
			p := &num_of_houses
			x := *p
			*p++
			return x
		}()] = temp_house
		SET_BIT_AR(world[real_house].Room_flags[:], ROOM_HOUSE)
		House_load(temp_house.Vnum)
	}
	fl.Close()
	House_save_control()
}

var HCONTROL_FORMAT *byte = libc.CString("Usage: hcontrol build <house vnum> <exit direction> <player name>\r\n       hcontrol destroy <house vnum>\r\n       hcontrol pay <house vnum>\r\n       hcontrol show\r\n")

func hcontrol_list_houses(ch *char_data) {
	var (
		i        int
		timestr  *byte
		temp     *byte
		built_on [128]byte
		last_pay [128]byte
		own_name [21]byte
	)
	if num_of_houses == 0 {
		send_to_char(ch, libc.CString("No houses have been defined.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Address  Atrium  Build Date  Guests  Owner        Last Paymt\r\n-------  ------  ----------  ------  ------------ ----------\r\n"))
	for i = 0; i < num_of_houses; i++ {
		if (func() *byte {
			temp = get_name_by_id(house_control[i].Owner)
			return temp
		}()) == nil {
			continue
		}
		if house_control[i].Built_on != 0 {
			timestr = libc.AscTime(libc.LocalTime(&house_control[i].Built_on))
			*((*byte)(unsafe.Add(unsafe.Pointer(timestr), 10))) = '\x00'
			strlcpy(&built_on[0], timestr, uint64(128))
		} else {
			libc.StrCpy(&built_on[0], libc.CString("Unknown"))
		}
		if house_control[i].Last_payment != 0 {
			timestr = libc.AscTime(libc.LocalTime(&house_control[i].Last_payment))
			*((*byte)(unsafe.Add(unsafe.Pointer(timestr), 10))) = '\x00'
			strlcpy(&last_pay[0], timestr, uint64(128))
		} else {
			libc.StrCpy(&last_pay[0], libc.CString("None"))
		}
		libc.StrCpy(&own_name[0], temp)
		send_to_char(ch, libc.CString("%7d %-10s    %2d    %-12s %s\r\n"), house_control[i].Vnum, &built_on[0], house_control[i].Num_of_guests, CAP(&own_name[0]), &last_pay[0])
		House_list_guests(ch, i, TRUE)
	}
}
func hcontrol_build_house(ch *char_data, arg *byte) {
	var (
		arg1       [2048]byte
		temp_house house_control_rec
		virt_house room_vnum
		real_house room_rnum
		exit_num   int16
		owner      int
	)
	if num_of_houses >= MAX_HOUSES {
		send_to_char(ch, libc.CString("Max houses already defined.\r\n"))
		return
	}
	arg = one_argument(arg, &arg1[0])
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
		return
	}
	virt_house = room_vnum(libc.Atoi(libc.GoString(&arg1[0])))
	if (func() room_rnum {
		real_house = real_room(virt_house)
		return real_house
	}()) == room_rnum(-1) {
		send_to_char(ch, libc.CString("No such room exists.\r\n"))
		return
	}
	if find_house(virt_house) != int(-1) {
		send_to_char(ch, libc.CString("House already exists.\r\n"))
		return
	}
	arg = one_argument(arg, &arg1[0])
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
		return
	}
	if int(func() int16 {
		exit_num = int16(search_block(&arg1[0], &dirs[0], FALSE))
		return exit_num
	}()) < 0 && int(func() int16 {
		exit_num = int16(search_block(&arg1[0], &abbr_dirs[0], FALSE))
		return exit_num
	}()) < 0 {
		send_to_char(ch, libc.CString("'%s' is not a valid direction.\r\n"), &arg1[0])
		return
	}
	if (func() room_rnum {
		if world[real_house].Dir_option[exit_num] != nil {
			return world[real_house].Dir_option[exit_num].To_room
		}
		return -1
	}()) == room_rnum(-1) {
		send_to_char(ch, libc.CString("There is no exit %s from room %d.\r\n"), dirs[exit_num], virt_house)
		return
	}
	one_argument(arg, &arg1[0])
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
		return
	}
	if (func() int {
		owner = get_id_by_name(&arg1[0])
		return owner
	}()) < 0 {
		send_to_char(ch, libc.CString("Unknown player '%s'.\r\n"), &arg1[0])
		return
	}
	temp_house.Mode = HOUSE_PRIVATE
	temp_house.Vnum = virt_house
	temp_house.Exit_num = exit_num
	temp_house.Built_on = libc.GetTime(nil)
	temp_house.Last_payment = 0
	temp_house.Owner = owner
	temp_house.Num_of_guests = 0
	house_control[func() int {
		p := &num_of_houses
		x := *p
		*p++
		return x
	}()] = temp_house
	SET_BIT_AR(world[real_house].Room_flags[:], ROOM_HOUSE)
	House_crashsave(virt_house)
	send_to_char(ch, libc.CString("House built.  Mazel tov!\r\n"))
	House_save_control()
}
func hcontrol_destroy_house(ch *char_data, arg *byte) {
	var (
		i           int
		j           int
		real_atrium room_rnum
		real_house  room_rnum
	)
	if *arg == 0 {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
		return
	}
	if (func() int {
		i = find_house(room_vnum(libc.Atoi(libc.GoString(arg))))
		return i
	}()) == int(-1) {
		send_to_char(ch, libc.CString("Unknown house.\r\n"))
		return
	}
	if (func() room_rnum {
		real_atrium = real_room(house_control[i].Atrium)
		return real_atrium
	}()) == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: House %d had invalid atrium %d!"), libc.Atoi(libc.GoString(arg)), house_control[i].Atrium)
	} else {
		REMOVE_BIT_AR(world[real_atrium].Room_flags[:], ROOM_ATRIUM)
	}
	if (func() room_rnum {
		real_house = real_room(house_control[i].Vnum)
		return real_house
	}()) == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: House %d had invalid vnum %d!"), libc.Atoi(libc.GoString(arg)), house_control[i].Vnum)
	} else {
		REMOVE_BIT_AR(world[real_house].Room_flags[:], ROOM_HOUSE)
		REMOVE_BIT_AR(world[real_house].Room_flags[:], ROOM_HOUSE_CRASH)
	}
	House_delete_file(house_control[i].Vnum)
	for j = i; j < num_of_houses-1; j++ {
		house_control[j] = house_control[j+1]
	}
	num_of_houses--
	send_to_char(ch, libc.CString("House deleted.\r\n"))
	House_save_control()
	for i = 0; i < num_of_houses; i++ {
		if (func() room_rnum {
			real_atrium = real_room(house_control[i].Atrium)
			return real_atrium
		}()) != room_rnum(-1) {
			SET_BIT_AR(world[real_atrium].Room_flags[:], ROOM_ATRIUM)
		}
	}
}
func hcontrol_pay_house(ch *char_data, arg *byte) {
	var i int
	if *arg == 0 {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
	} else if (func() int {
		i = find_house(room_vnum(libc.Atoi(libc.GoString(arg))))
		return i
	}()) == int(-1) {
		send_to_char(ch, libc.CString("Unknown house.\r\n"))
	} else {
		mudlog(NRM, int(MAX(ADMLVL_IMMORT, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("Payment for house %s collected by %s."), arg, GET_NAME(ch))
		house_control[i].Last_payment = libc.GetTime(nil)
		House_save_control()
		send_to_char(ch, libc.CString("Payment recorded.\r\n"))
	}
}
func do_hcontrol(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1 [2048]byte
		arg2 [2048]byte
	)
	half_chop(argument, &arg1[0], &arg2[0])
	if is_abbrev(&arg1[0], libc.CString("build")) != 0 {
		hcontrol_build_house(ch, &arg2[0])
	} else if is_abbrev(&arg1[0], libc.CString("destroy")) != 0 {
		hcontrol_destroy_house(ch, &arg2[0])
	} else if is_abbrev(&arg1[0], libc.CString("pay")) != 0 {
		hcontrol_pay_house(ch, &arg2[0])
	} else if is_abbrev(&arg1[0], libc.CString("show")) != 0 {
		hcontrol_list_houses(ch)
	} else {
		send_to_char(ch, libc.CString("%s"), HCONTROL_FORMAT)
	}
}
func do_house(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		i   int
		j   int
		id  int
	)
	one_argument(argument, &arg[0])
	if !ROOM_FLAGGED(ch.In_room, ROOM_HOUSE) {
		send_to_char(ch, libc.CString("You must be in your house to set guests.\r\n"))
	} else if (func() int {
		i = find_house(room_vnum(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))))
		return i
	}()) == int(-1) {
		send_to_char(ch, libc.CString("Um.. this house seems to be screwed up.\r\n"))
	} else if int(ch.Idnum) != house_control[i].Owner {
		send_to_char(ch, libc.CString("Only the primary owner can set guests.\r\n"))
	} else if arg[0] == 0 {
		House_list_guests(ch, i, FALSE)
	} else if (func() int {
		id = get_id_by_name(&arg[0])
		return id
	}()) < 0 {
		send_to_char(ch, libc.CString("No such player.\r\n"))
	} else if id == int(ch.Idnum) {
		send_to_char(ch, libc.CString("It's your house!\r\n"))
	} else {
		for j = 0; j < house_control[i].Num_of_guests; j++ {
			if house_control[i].Guests[j] == id {
				for ; j < house_control[i].Num_of_guests; j++ {
					house_control[i].Guests[j] = house_control[i].Guests[j+1]
				}
				house_control[i].Num_of_guests--
				House_save_control()
				send_to_char(ch, libc.CString("Guest deleted.\r\n"))
				return
			}
		}
		if house_control[i].Num_of_guests == MAX_GUESTS {
			send_to_char(ch, libc.CString("You have too many guests.\r\n"))
			return
		}
		j = func() int {
			p := &house_control[i].Num_of_guests
			x := *p
			*p++
			return x
		}()
		house_control[i].Guests[j] = id
		House_save_control()
		send_to_char(ch, libc.CString("Guest added.\r\n"))
	}
}
func House_save_all() {
	var (
		i          int
		real_house room_rnum
	)
	for i = 0; i < num_of_houses; i++ {
		if (func() room_rnum {
			real_house = real_room(house_control[i].Vnum)
			return real_house
		}()) != room_rnum(-1) {
			House_crashsave(house_control[i].Vnum)
		}
	}
}
func House_can_enter(ch *char_data, house room_vnum) int {
	var (
		i int
		j int
	)
	if ADM_FLAGGED(ch, ADM_ALLHOUSES) || (func() int {
		i = find_house(house)
		return i
	}()) == int(-1) {
		return 1
	}
	switch house_control[i].Mode {
	case HOUSE_CLAN:
		fallthrough
	case HOUSE_UNOWNED:
		return 1
	case HOUSE_PRIVATE:
		if int(ch.Idnum) == house_control[i].Owner {
			return 1
		}
		for j = 0; j < house_control[i].Num_of_guests; j++ {
			if int(ch.Idnum) == house_control[i].Guests[j] {
				return 1
			}
		}
	}
	return 0
}
func House_list_guests(ch *char_data, i int, quiet int) {
	var (
		j           int
		num_printed int
		temp        *byte
	)
	if house_control[i].Num_of_guests == 0 {
		if quiet == 0 {
			send_to_char(ch, libc.CString("  Guests: None\r\n"))
		}
		return
	}
	send_to_char(ch, libc.CString("  Guests: "))
	for num_printed = func() int {
		j = 0
		return j
	}(); j < house_control[i].Num_of_guests; j++ {
		if (func() *byte {
			temp = get_name_by_id(house_control[i].Guests[j])
			return temp
		}()) == nil {
			continue
		}
		num_printed++
		send_to_char(ch, libc.CString("%c%s "), unicode.ToUpper(rune(*temp)), (*byte)(unsafe.Add(unsafe.Pointer(temp), 1)))
	}
	if num_printed == 0 {
		send_to_char(ch, libc.CString("all dead"))
	}
	send_to_char(ch, libc.CString("\r\n"))
}
func House_load(rvnum room_vnum) int {
	var (
		fl      *stdio.File
		f1      [256]byte
		f2      [256]byte
		f3      [256]byte
		f4      [256]byte
		cmfname [64936]byte
		buf1    [64936]byte
		buf2    [64936]byte
		line    [256]byte
		t       [21]int
		danger  int
	)
	_ = danger
	var zwei int = 0
	var temp *obj_data
	var locate int = 0
	var j int
	var nr int
	var k int
	_ = k
	var num_objs int = 0
	var obj1 *obj_data
	var cont_row [5]*obj_data
	var new_descr *extra_descr_data
	var rrnum room_rnum
	if (func() room_rnum {
		rrnum = real_room(rvnum)
		return rrnum
	}()) == room_rnum(-1) {
		return 0
	}
	if House_get_filename(rvnum, &cmfname[0], uint64(64936)) == 0 {
		return 0
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&cmfname[0]), "r+b")
		return fl
	}()) == nil {
		if libc.Errno != 2 {
			stdio.Sprintf(&buf1[0], "SYSERR: READING HOUSE FILE %s (5)", &cmfname[0])
			fmt.Println(&buf1[0])
		}
		return 0
	}
	for j = 0; j < MAX_BAG_ROWS; j++ {
		cont_row[j] = nil
	}
	if int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
	}
	for int(fl.IsEOF()) == 0 {
		temp = nil
		if line[0] == '#' {
			if stdio.Sscanf(&line[0], "#%d", &nr) != 1 {
				continue
			}
			if nr == int(-1) {
				temp = create_obj()
				temp.Item_number = -1
			} else if nr < 0 {
				continue
			} else {
				if nr >= 0xF423F {
					continue
				}
				temp = read_object(obj_vnum(nr), VIRTUAL)
				if temp == nil {
					get_line(fl, &line[0])
					continue
				}
			}
			get_line(fl, &line[0])
			stdio.Sscanf(&line[0], "%d %d %d %d %d %d %d %d %d %s %s %s %s %d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7], &t[8], &f1[0], &f2[0], &f3[0], &f4[0], &t[13], &t[14], &t[15], &t[16], &t[17], &t[18], &t[19], &t[20])
			locate = t[0]
			temp.Value[0] = t[1]
			temp.Value[1] = t[2]
			temp.Value[2] = t[3]
			temp.Value[3] = t[4]
			temp.Value[4] = t[5]
			temp.Value[5] = t[6]
			temp.Value[6] = t[7]
			temp.Value[7] = t[8]
			temp.Extra_flags[0] = asciiflag_conv(&f1[0])
			temp.Extra_flags[1] = asciiflag_conv(&f2[0])
			temp.Extra_flags[2] = asciiflag_conv(&f3[0])
			temp.Extra_flags[3] = asciiflag_conv(&f3[0])
			temp.Value[8] = t[13]
			temp.Value[9] = t[14]
			temp.Value[10] = t[15]
			temp.Value[11] = t[16]
			temp.Value[12] = t[17]
			temp.Value[13] = t[18]
			temp.Value[14] = t[19]
			temp.Value[15] = t[20]
			temp.Posted_to = nil
			temp.Posttype = 0
			get_line(fl, &line[0])
			if libc.StrCmp(libc.CString("XAP"), &line[0]) == 0 {
				if (func() *byte {
					p := &temp.Name
					temp.Name = fread_string(fl, &buf2[0])
					return *p
				}()) == nil {
					temp.Name = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Short_description
					temp.Short_description = fread_string(fl, &buf2[0])
					return *p
				}()) == nil {
					temp.Short_description = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Description
					temp.Description = fread_string(fl, &buf2[0])
					return *p
				}()) == nil {
					temp.Description = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Action_description
					temp.Action_description = fread_string(fl, &buf2[0])
					return *p
				}()) == nil {
					temp.Action_description = nil
				}
				if get_line(fl, &line[0]) == 0 || stdio.Sscanf(&line[0], "%d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7]) != 8 {
					stdio.Fprintf(stdio.Stderr(), "Format error in first numeric line (expecting _x_ args)")
					return 0
				}
				temp.Type_flag = int8(t[0])
				temp.Wear_flags[0] = bitvector_t(int32(t[1]))
				temp.Wear_flags[1] = bitvector_t(int32(t[2]))
				temp.Wear_flags[2] = bitvector_t(int32(t[3]))
				temp.Wear_flags[3] = bitvector_t(int32(t[4]))
				temp.Weight = int64(t[5])
				temp.Cost = t[6]
				temp.Cost_per_day = t[7]
				libc.StrCat(&buf2[0], libc.CString(", after numeric constants (expecting E/#xxx)"))
				for j = 0; j < MAX_OBJ_AFFECT; j++ {
					temp.Affected[j].Location = APPLY_NONE
					temp.Affected[j].Modifier = 0
					temp.Affected[j].Specific = 0
				}
				temp.Ex_description = nil
				get_line(fl, &line[0])
				for k = func() int {
					j = func() int {
						zwei = 0
						return zwei
					}()
					return j
				}(); zwei == 0 && int(fl.IsEOF()) == 0; {
					switch line[0] {
					case 'E':
						new_descr = new(extra_descr_data)
						new_descr.Keyword = fread_string(fl, &buf2[0])
						new_descr.Description = fread_string(fl, &buf2[0])
						new_descr.Next = temp.Ex_description
						temp.Ex_description = new_descr
						get_line(fl, &line[0])
					case 'A':
						if j >= MAX_OBJ_AFFECT {
							basic_mud_log(libc.CString("SYSERR: Too many object affectations in loading house file"))
							danger = 1
						}
						get_line(fl, &line[0])
						stdio.Sscanf(&line[0], "%d %d %d", &t[0], &t[1], &t[2])
						temp.Affected[j].Location = t[0]
						temp.Affected[j].Modifier = t[1]
						temp.Affected[j].Specific = t[2]
						j++
						get_line(fl, &line[0])
					case 'G':
						get_line(fl, &line[0])
						stdio.Sscanf(&line[0], "%ld", &temp.Generation)
						get_line(fl, &line[0])
					case 'U':
						get_line(fl, &line[0])
						stdio.Sscanf(&line[0], "%lld", &temp.Unique_id)
						get_line(fl, &line[0])
					case 'S':
						if j >= SPELLBOOK_SIZE {
							basic_mud_log(libc.CString("SYSERR: Too many spells in spellbook loading rent file"))
							danger = 1
						}
						get_line(fl, &line[0])
						stdio.Sscanf(&line[0], "%d %d", &t[0], &t[1])
						if temp.Sbinfo == nil {
							temp.Sbinfo = make([]obj_spellbook_spell, SPELLBOOK_SIZE)
							libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(&temp.Sbinfo[0]))), 0, int(SPELLBOOK_SIZE*unsafe.Sizeof(obj_spellbook_spell{})))
						}
						temp.Sbinfo[j].Spellname = t[0]
						temp.Sbinfo[j].Pages = t[1]
						j++
						get_line(fl, &line[0])
					case 'Z':
						get_line(fl, &line[0])
						stdio.Sscanf(&line[0], "%d", &temp.Size)
						get_line(fl, &line[0])
					case '$':
						fallthrough
					case '#':
						zwei = 1
					default:
						zwei = 1
					}
				}
			}
			if temp != nil {
				num_objs++
				obj_to_room(temp, rrnum)
			} else {
				continue
			}
			for j = int(MAX_BAG_ROWS - 1); j > -locate; j-- {
				if cont_row[j] != nil {
					for ; cont_row[j] != nil; cont_row[j] = obj1 {
						obj1 = cont_row[j].Next_content
						obj_to_room(cont_row[j], rrnum)
					}
					cont_row[j] = nil
				}
			}
			if j == -locate && cont_row[j] != nil {
				if int(temp.Type_flag) == ITEM_CONTAINER {
					obj_from_room(temp)
					temp.Contains = nil
					for ; cont_row[j] != nil; cont_row[j] = obj1 {
						obj1 = cont_row[j].Next_content
						obj_to_obj(cont_row[j], temp)
					}
					obj_to_room(temp, rrnum)
				} else {
					for ; cont_row[j] != nil; cont_row[j] = obj1 {
						obj1 = cont_row[j].Next_content
						obj_to_room(cont_row[j], rrnum)
					}
					cont_row[j] = nil
				}
			}
			if locate < 0 && locate >= -MAX_BAG_ROWS {
				obj_from_room(temp)
				if (func() *obj_data {
					obj1 = cont_row[-locate-1]
					return obj1
				}()) != nil {
					for obj1.Next_content != nil {
						obj1 = obj1.Next_content
					}
					obj1.Next_content = temp
				} else {
					cont_row[-locate-1] = temp
				}
			}
		} else {
			get_line(fl, &line[0])
		}
	}
	fl.Close()
	return 1
}
