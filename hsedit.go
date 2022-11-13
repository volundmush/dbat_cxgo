package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

var house_flags [8]*byte = [8]*byte{libc.CString("!GUEST"), libc.CString("FREE"), libc.CString("!IMM"), libc.CString("IMP_ONLY"), libc.CString("RENTFREE"), libc.CString("SAVE_!RENT"), libc.CString("!SAVE"), libc.CString("\n")}
var house_types [5]*byte = [5]*byte{libc.CString("PLAYER_OWNED"), libc.CString("IMM_OWNED"), libc.CString("CLAN_OWNED"), libc.CString("UNOWNED"), libc.CString("\n")}

func hsedit_setup_new(d *descriptor_data) {
	var i int
	d.Olc.House = new(house_control_rec)
	d.Olc.House.Vnum = d.Olc.Number
	d.Olc.House.Owner = 0
	d.Olc.House.Atrium = 0
	d.Olc.House.Exit_num = -1
	d.Olc.House.Built_on = libc.GetTime(nil)
	d.Olc.House.Mode = HOUSE_PRIVATE
	d.Olc.House.Bitvector = 0
	d.Olc.House.Builtby = 0
	d.Olc.House.Last_payment = 0
	d.Olc.House.Num_of_guests = 0
	for i = 0; i < MAX_GUESTS; i++ {
		d.Olc.House.Guests[i] = 0
	}
	d.Olc.Value = 0
	hsedit_disp_menu(d)
}
func hsedit_setup_existing(d *descriptor_data, real_num int) {
	var (
		house *house_control_rec
		i     int
	)
	house = new(house_control_rec)
	*house = house_control[real_num]
	house.Vnum = house_control[real_num].Vnum
	house.Atrium = house_control[real_num].Atrium
	house.Owner = house_control[real_num].Owner
	house.Exit_num = house_control[real_num].Exit_num
	house.Built_on = house_control[real_num].Built_on
	house.Mode = house_control[real_num].Mode
	house.Bitvector = house_control[real_num].Bitvector
	house.Last_payment = house_control[real_num].Last_payment
	house.Num_of_guests = house_control[real_num].Num_of_guests
	for i = 0; i < MAX_GUESTS; i++ {
		house.Guests[i] = house_control[real_num].Guests[i]
	}
	d.Olc.House = house
	d.Olc.Value = 0
	hsedit_disp_menu(d)
}
func hsedit_save_internally(d *descriptor_data) {
	var house_rnum int
	house_rnum = find_house(d.Olc.Number)
	if house_rnum != int(-1) {
		free_house(&house_control[house_rnum])
		house_control[house_rnum] = *d.Olc.House
	} else {
		house_rnum = func() int {
			p := &num_of_houses
			x := *p
			*p++
			return x
		}()
		if house_rnum < MAX_HOUSES {
			house_control[house_rnum] = *d.Olc.House
			house_control[house_rnum].Vnum = d.Olc.Number
		} else {
			send_to_char(d.Character, libc.CString("MAX House limit reached - Unable to save this house!"))
			mudlog(NRM, ADMLVL_BUILDER, TRUE, libc.CString("HSEDIT: Max houses limit reached - Unable to save OLC data"))
		}
	}
	if real_room(d.Olc.House.Vnum) != room_rnum(-1) {
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(d.Olc.House.Vnum))))).Room_flags[int(ROOM_HOUSE/32)] |= bitvector_t(int32(1 << (int(ROOM_HOUSE % 32))))
	}
	if real_room(d.Olc.House.Atrium) != room_rnum(-1) {
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(d.Olc.House.Atrium))))).Room_flags[int(ROOM_ATRIUM/32)] |= bitvector_t(int32(1 << (int(ROOM_ATRIUM % 32))))
	}
}
func hsedit_save_to_disk() {
	House_save_control()
}
func free_house(house *house_control_rec) {
}
func hedit_delete_house(d *descriptor_data, house_vnum int) {
	var (
		i           int
		j           int
		real_atrium room_rnum
	)
	_ = real_atrium
	var real_house room_rnum
	if (func() int {
		i = find_house(room_vnum(house_vnum))
		return i
	}()) == int(-1) {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: hsedit: Invalid house vnum in hedit_delete_house\r\n"))
		cleanup_olc(d, CLEANUP_STRUCTS)
		return
	}
	if (func() room_rnum {
		real_house = real_room(house_control[i].Vnum)
		return real_house
	}()) == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: House %d had invalid vnum %d!"), house_vnum, house_control[i].Vnum)
	} else {
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_house)))).Room_flags[(int(ROOM_HOUSE|ROOM_HOUSE_CRASH))/32] &= bitvector_t(int32(^(1 << ((int(ROOM_HOUSE | ROOM_HOUSE_CRASH)) % 32))))
	}
	House_delete_file(house_control[i].Vnum)
	for j = i; j < num_of_houses-1; j++ {
		house_control[j] = house_control[j+1]
	}
	num_of_houses--
	send_to_char(d.Character, libc.CString("House deleted.\r\n"))
	House_save_control()
	cleanup_olc(d, CLEANUP_ALL)
}
func hsedit_disp_flags_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
		buf1    [64936]byte
	)
	clear_screen(d)
	for counter = 0; counter < HOUSE_NUM_FLAGS; counter++ {
		send_to_char(d.Character, libc.CString("@g%2d@D)@y %-20.20s @n"), counter+1, house_flags[counter])
		if (func() int {
			p := &columns
			*p++
			return *p
		}() % 2) == 0 {
			send_to_char(d.Character, libc.CString("\r\n"))
		}
	}
	sprintbit(bitvector_t(int32(d.Olc.House.Bitvector)), house_flags[:], &buf1[0], uint64(64936))
	send_to_char(d.Character, libc.CString("\r\nHouse flags: @g%s@n\r\nEnter house flags, 0 to quit : "), &buf1[0])
	d.Olc.Mode = HSEDIT_FLAGS
}
func hsedit_owner_menu(d *descriptor_data) {
	var (
		buf   [64936]byte
		house *house_control_rec
	)
	house = d.Olc.House
	stdio.Sprintf(&buf[0], "@g1@D)@g Owner Name : @c%s@n\r\n@g2@D)@g Owner ID   : @c%ld\r\n@gQ@D)@g Back to main menu\r\n@gEnter choice : @n", get_name_by_id(house.Owner), house.Owner)
	send_to_char(d.Character, &buf[0])
	d.Olc.Mode = HSEDIT_OWNER_MENU
}
func hsedit_dir_menu(d *descriptor_data) {
	var (
		buf        [64936]byte
		house      *house_control_rec
		house_rnum int
		newroom    [12]int
		i          int
	)
	house = d.Olc.House
	mudlog(CMP, ADMLVL_BUILDER, TRUE, libc.CString("(LOG) hsedit_dir_menu: house vnum = %d"), house.Vnum)
	house_rnum = int(real_room(house.Vnum))
	if house_rnum < 0 || house_rnum == int(-1) {
		stdio.Sprintf(&buf[0], "WARNING: You cannot set an atium direction before selecting a valid room vnum\r\n(Press Enter)\r\n")
		d.Olc.Mode = HSEDIT_NOVNUM
	} else {
		for i = 0; i < 12; i++ {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(house_rnum)))).Dir_option[i] != nil {
				newroom[i] = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(house_rnum)))).Dir_option[i].To_room)
			} else {
				newroom[i] = -1
			}
		}
		stdio.Sprintf(&buf[0], "@g1@D)@g North  : @D(@c%s@D)\r\n@g2@D)@g East   : @D(@c%s@D)\r\n@g3@D)@g South  : @D(@c%s@D)\r\n@g4@D)@g West   : @D(@c%s@D)\r\n@g5@D)@g Up     : @D(@c%s@D)\r\n@g6@D)@g Down   : @D(@c%s@D)\r\n@g7@D)@g NorthW : @D(@c%s@D)\r\n@g8@D)@g NorthE : @D(@c%s@D)\r\n@g9@D)@g SouthE : @D(@c%s@D)\r\n@g10@D)@g SouthW : @D(@c%s@D)\r\n@g11@D)@g Inside : @D(@c%s@D)\r\n@g12@D)@g Outside: @D(@c%s@D)\r\n@gQ@D)@g Back to main menu\r\n@gEnter atrium direction : @n", func() string {
			if newroom[0] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[0])))).Name)
		}(), func() string {
			if newroom[1] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[1])))).Name)
		}(), func() string {
			if newroom[2] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[2])))).Name)
		}(), func() string {
			if newroom[3] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[3])))).Name)
		}(), func() string {
			if newroom[4] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[4])))).Name)
		}(), func() string {
			if newroom[5] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[5])))).Name)
		}(), func() string {
			if newroom[6] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[6])))).Name)
		}(), func() string {
			if newroom[7] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[7])))).Name)
		}(), func() string {
			if newroom[8] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[8])))).Name)
		}(), func() string {
			if newroom[9] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[9])))).Name)
		}(), func() string {
			if newroom[10] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[10])))).Name)
		}(), func() string {
			if newroom[11] == int(-1) {
				return "NO ROOM"
			}
			return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom[11])))).Name)
		}())
		d.Olc.Mode = HSEDIT_DIR_MENU
	}
	send_to_char(d.Character, &buf[0])
}
func hsedit_disp_type_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_HOUSE_TYPES; counter++ {
		send_to_char(d.Character, libc.CString("@g%2d@D)@y %-20.20s @n"), counter, house_types[counter])
		if (func() int {
			p := &columns
			*p++
			return *p
		}() % 2) == 0 {
			send_to_char(d.Character, libc.CString("\r\n"))
		}
	}
	send_to_char(d.Character, libc.CString("\r\nEnter house type : "))
	d.Olc.Mode = HSEDIT_TYPE
}
func hsedit_disp_guest_menu(d *descriptor_data) {
	var (
		buf     [64936]byte
		not_set [128]byte
		house   *house_control_rec
	)
	house = d.Olc.House
	stdio.Sprintf(&not_set[0], "@D(@yNOT SET@D)@n")
	stdio.Sprintf(&buf[0], " @g1@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g2@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g3@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g4@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g5@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g6@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g7@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g8@D)@c %s @D(@gID: @c%ld@D)@n\r\n @g9@D)@c %s @D(@gID: @c%ld@D)@n\r\n@g10@D)@c %s @D(@gID: @c%ld@D)@n\r\n\r\n@gA@D)@g Add a guest\r\n@gD@D)@g Delete a guest\r\n@gC@D)@g Clear guest list\r\n@gQ@D)@g Back to main menu\r\n@gEnter selection @D(@gA@D/@gD@D/@gC@D/@gQ@D)@n: ", func() *byte {
		if get_name_by_id(house.Guests[0]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[0])
	}(), func() int {
		if house.Guests[0] < 1 {
			return 0
		}
		return house.Guests[0]
	}(), func() *byte {
		if get_name_by_id(house.Guests[1]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[1])
	}(), func() int {
		if house.Guests[1] < 1 {
			return 0
		}
		return house.Guests[1]
	}(), func() *byte {
		if get_name_by_id(house.Guests[2]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[2])
	}(), func() int {
		if house.Guests[2] < 1 {
			return 0
		}
		return house.Guests[2]
	}(), func() *byte {
		if get_name_by_id(house.Guests[3]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[3])
	}(), func() int {
		if house.Guests[3] < 1 {
			return 0
		}
		return house.Guests[3]
	}(), func() *byte {
		if get_name_by_id(house.Guests[4]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[4])
	}(), func() int {
		if house.Guests[4] < 1 {
			return 0
		}
		return house.Guests[4]
	}(), func() *byte {
		if get_name_by_id(house.Guests[5]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[5])
	}(), func() int {
		if house.Guests[5] < 1 {
			return 0
		}
		return house.Guests[5]
	}(), func() *byte {
		if get_name_by_id(house.Guests[6]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[6])
	}(), func() int {
		if house.Guests[6] < 1 {
			return 0
		}
		return house.Guests[6]
	}(), func() *byte {
		if get_name_by_id(house.Guests[7]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[7])
	}(), func() int {
		if house.Guests[7] < 1 {
			return 0
		}
		return house.Guests[7]
	}(), func() *byte {
		if get_name_by_id(house.Guests[8]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[8])
	}(), func() int {
		if house.Guests[8] < 1 {
			return 0
		}
		return house.Guests[8]
	}(), func() *byte {
		if get_name_by_id(house.Guests[9]) == nil {
			return &not_set[0]
		}
		return get_name_by_id(house.Guests[9])
	}(), func() int {
		if house.Guests[9] < 1 {
			return 0
		}
		return house.Guests[9]
	}())
	send_to_char(d.Character, &buf[0])
	d.Olc.Mode = HSEDIT_GUEST_MENU
}
func hsedit_disp_val0_menu(d *descriptor_data) {
	d.Olc.Mode = HSEDIT_VALUE_0
	switch d.Olc.House.Mode {
	case HOUSE_CLAN:
		send_to_char(d.Character, libc.CString("Enter id of the clan: "))
	case HOUSE_UNOWNED:
		send_to_char(d.Character, libc.CString("Done."))
	case HOUSE_GOD:
		hsedit_disp_val3_menu(d)
	default:
		hsedit_disp_menu(d)
	}
}
func hsedit_disp_val1_menu(d *descriptor_data) {
	d.Olc.Mode = HSEDIT_VALUE_1
	switch d.Olc.House.Mode {
	default:
		hsedit_disp_menu(d)
	}
}
func hsedit_disp_val2_menu(d *descriptor_data) {
	d.Olc.Mode = HSEDIT_VALUE_2
	switch d.Olc.House.Mode {
	default:
		hsedit_disp_menu(d)
	}
}
func hsedit_disp_val3_menu(d *descriptor_data) {
	d.Olc.Mode = HSEDIT_VALUE_3
	switch d.Olc.House.Mode {
	case HOUSE_GOD:
		send_to_char(d.Character, libc.CString("Enter minimum level of guests: "))
	default:
		hsedit_disp_menu(d)
	}
}
func hsedit_list_guests(thishouse *house_control_rec, guestlist *byte) *byte {
	var (
		j           int
		num_printed int
		temp        *byte
	)
	if thishouse.Num_of_guests == 0 {
		stdio.Sprintf(guestlist, "NONE")
		return guestlist
	}
	for num_printed = func() int {
		j = 0
		return j
	}(); j < thishouse.Num_of_guests; j++ {
		if (func() *byte {
			temp = get_name_by_id(thishouse.Guests[j])
			return temp
		}()) == nil {
			continue
		}
		num_printed++
		stdio.Sprintf(guestlist, "%s%c%s ", guestlist, unicode.ToUpper(rune(*temp)), (*byte)(unsafe.Add(unsafe.Pointer(temp), 1)))
	}
	if num_printed == 0 {
		stdio.Sprintf(guestlist, "all dead")
	}
	return guestlist
}
func hsedit_disp_menu(d *descriptor_data) {
	var (
		buf      [64936]byte
		buf1     [64936]byte
		built_on [128]byte
		last_pay [128]byte
		buf2     [64936]byte
		timestr  *byte
		no_name  [128]byte
		house    *house_control_rec
	)
	clear_screen(d)
	house = d.Olc.House
	if house.Built_on != 0 {
		timestr = libc.AscTime(libc.LocalTime(&house.Built_on))
		*((*byte)(unsafe.Add(unsafe.Pointer(timestr), 10))) = '\x00'
		strlcpy(&built_on[0], timestr, uint64(128))
	} else {
		libc.StrCpy(&built_on[0], libc.CString("Unknown"))
	}
	if house.Last_payment != 0 {
		timestr = libc.AscTime(libc.LocalTime(&house.Last_payment))
		*((*byte)(unsafe.Add(unsafe.Pointer(timestr), 10))) = '\x00'
		strlcpy(&last_pay[0], timestr, uint64(128))
	} else {
		libc.StrCpy(&last_pay[0], libc.CString("None"))
	}
	buf2[0] = '\x00'
	sprintbit(bitvector_t(int32(house.Bitvector)), house_flags[:], &buf1[0], uint64(64936))
	stdio.Sprintf(&no_name[0], "(NOBODY)")
	stdio.Sprintf(&buf[0], "@D-- @RJamdog's House OLC Editor @D--\r\n@D--@g House number : @D[@c%d@D]     @gHouse zone: @D[@c%d@D]\r\n@g1@D)@g Owner       : @c%ld -- %s\r\n@g2@D)@g Atrium      : @c%d\r\n@g3@D)@g Direction   : @c%s\r\n@g4@D)@g House Type  : @c%s\r\n@g5@D)@g Built on    : @c%s\r\n@g6@D)@g Payment     : @c%s\r\n@g7@D)@g Guests      : @c%s\r\n@g8@D)@g Flags       : @c%s\r\n@gX@D)@g Delete this house\r\n@gQ@D)@g Quit\r\n@gEnter choice : @n", d.Olc.Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, house.Owner, func() *byte {
		if get_name_by_id(house.Owner) == nil {
			return &no_name[0]
		}
		return get_name_by_id(house.Owner)
	}(), house.Atrium, func() *byte {
		if int(house.Exit_num) >= 0 && int(house.Exit_num) <= 11 {
			return dirs[house.Exit_num]
		}
		return libc.CString("NONE")
	}(), house_types[house.Mode], &built_on[0], &last_pay[0], hsedit_list_guests(house, &buf2[0]), &buf1[0])
	send_to_char(d.Character, &buf[0])
	d.Olc.Mode = HSEDIT_MAIN_MENU
}
func hsedit_parse(d *descriptor_data, arg *byte) {
	var (
		number    int = 0
		id        int = 0
		i         int
		room_rnum int
		tmp       *byte
		found     bool = FALSE != 0
	)
	mudlog(CMP, ADMLVL_BUILDER, FALSE, libc.CString("(LOG) hsedit_parse: OLC mode %d"), d.Olc.Mode)
	switch d.Olc.Mode {
	case HSEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			hsedit_save_internally(d)
			mudlog(CMP, ADMLVL_BUILDER, TRUE, libc.CString("OLC: %s edits house %d"), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				hsedit_save_to_disk()
				write_to_output(d, libc.CString("House saved to disk.\r\n"))
			} else {
				write_to_output(d, libc.CString("House saved to memory.\r\n"))
			}
			cleanup_olc(d, CLEANUP_STRUCTS)
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
		default:
			send_to_char(d.Character, libc.CString("Invalid choice!\r\n"))
			send_to_char(d.Character, libc.CString("Do you wish to save this house internally? : "))
		}
		return
	case HSEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				send_to_char(d.Character, libc.CString("Do you wish to save this house internally? : "))
				d.Olc.Mode = HSEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '1':
			hsedit_owner_menu(d)
		case '2':
			if d.Olc.House.Vnum == room_vnum(-1) || real_room(d.Olc.House.Vnum) == room_rnum(-1) {
				send_to_char(d.Character, libc.CString("ERROR: Invalid house VNUM\r\n(Press Enter)\r\n"))
				mudlog(NRM, ADMLVL_GRGOD, TRUE, libc.CString("SYSERR: Invalid house VNUM in hsedit"))
			} else {
				send_to_char(d.Character, libc.CString("Enter atrium room vnum:"))
				d.Olc.Mode = HSEDIT_ATRIUM
			}
		case '3':
			if d.Olc.House.Vnum == room_vnum(-1) || real_room(d.Olc.House.Vnum) == room_rnum(-1) {
				send_to_char(d.Character, libc.CString("ERROR: Invalid house VNUM\r\n(Press Enter)\r\n"))
				mudlog(NRM, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Invalid house VNUM in hsedit"))
			} else {
				hsedit_dir_menu(d)
			}
		case '4':
			hsedit_disp_type_menu(d)
		case '5':
			send_to_char(d.Character, libc.CString("Set build date to now? (Y/N):"))
			d.Olc.Mode = HSEDIT_BUILD_DATE
		case '6':
			send_to_char(d.Character, libc.CString("Set last payment as now? (Y/N) : "))
			d.Olc.Mode = HSEDIT_PAYMENT
		case '7':
			hsedit_disp_guest_menu(d)
		case '8':
			hsedit_disp_flags_menu(d)
		case 'x':
			fallthrough
		case 'X':
			send_to_char(d.Character, libc.CString("Are you sure you want to delete this house? (Y/N) : "))
			d.Olc.Mode = HSEDIT_DELETE
		default:
			send_to_char(d.Character, libc.CString("Invalid choice!\r\n"))
			hsedit_disp_menu(d)
		}
		return
	case HSEDIT_OWNER_MENU:
		switch *arg {
		case '1':
			send_to_char(d.Character, libc.CString("Enter the name of the owner : "))
			d.Olc.Mode = HSEDIT_OWNER_NAME
		case '2':
			send_to_char(d.Character, libc.CString("Enter the user id of the owner : "))
			d.Olc.Mode = HSEDIT_OWNER_ID
		case 'Q':
			hsedit_disp_menu(d)
		}
		return
	case HSEDIT_OWNER_NAME:
		if (func() int {
			id = get_id_by_name(arg)
			return id
		}()) < 0 {
			send_to_char(d.Character, libc.CString("There is no such player.\r\n"))
			hsedit_owner_menu(d)
			return
		} else {
			d.Olc.House.Owner = id
		}
	case HSEDIT_OWNER_ID:
		id = libc.Atoi(libc.GoString(arg))
		if (func() *byte {
			tmp = get_name_by_id(id)
			return tmp
		}()) == nil {
			send_to_char(d.Character, libc.CString("There is no such player.\r\n"))
			hsedit_owner_menu(d)
			return
		} else {
			d.Olc.House.Owner = id
		}
	case HSEDIT_ATRIUM:
		number = libc.Atoi(libc.GoString(arg))
		if number == 0 {
			hsedit_disp_menu(d)
			return
		}
		room_rnum = int(real_room(d.Olc.House.Vnum))
		if real_room(room_vnum(number)) == room_rnum(-1) {
			send_to_char(d.Character, libc.CString("Room VNUM does not exist.\r\nEnter a valid room VNUM for this atrium (0 to exit) : "))
			return
		} else {
			for i = 0; i < 12; i++ {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_rnum)))).Dir_option[i] != nil {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_rnum)))).Dir_option[i].To_room == real_room(room_vnum(number)) {
						found = TRUE != 0
						id = i
					}
				}
			}
			if int(libc.BoolToInt(found)) == FALSE {
				send_to_char(d.Character, libc.CString("Atrium MUST be an adjoining room.\r\nEnter a valid room VNUM for this atrium (0 to exit) : "))
				return
			} else {
				d.Olc.House.Atrium = room_vnum(number)
				d.Olc.House.Exit_num = int16(id)
			}
		}
	case HSEDIT_DIR_MENU:
		number = libc.Atoi(libc.GoString(arg)) - 1
		if *arg == 'q' || *arg == 'Q' || number == -1 {
			hsedit_disp_menu(d)
			return
		}
		if number < 0 || number > 12 {
			send_to_char(d.Character, libc.CString("Invalid choice, Please select a direction (1-12, Q to quit) : "))
			return
		}
		id = int(real_room(d.Olc.House.Vnum))
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(id)))).Dir_option[number]) == nil {
			send_to_char(d.Character, libc.CString("You cannot set the atrium to a room that doesn't exist!\r\n"))
			hsedit_dir_menu(d)
			return
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(id)))).Dir_option[number].To_room == room_rnum(-1) {
			send_to_char(d.Character, libc.CString("You cannot set the atrium to nowhere!\r\n"))
			hsedit_dir_menu(d)
			return
		} else {
			d.Olc.House.Exit_num = int16(number)
			room_rnum = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(id)))).Dir_option[number].To_room)
			d.Olc.House.Atrium = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_rnum)))).Number
		}
	case HSEDIT_NOVNUM:
	case HSEDIT_BUILD_DATE:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.House.Built_on = libc.GetTime(nil)
		case 'n':
			fallthrough
		case 'N':
			send_to_char(d.Character, libc.CString("Build Date not changed\r\n"))
			hsedit_disp_menu(d)
			return
		}
	case HSEDIT_DELETE:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			hedit_delete_house(d, int(d.Olc.House.Vnum))
			return
		case 'n':
			fallthrough
		case 'N':
			send_to_char(d.Character, libc.CString("House not deleted!\r\n"))
			hsedit_disp_menu(d)
			return
		}
	case HSEDIT_BUILDER:
		if (func() int {
			id = get_id_by_name(arg)
			return id
		}()) < 0 {
			send_to_char(d.Character, libc.CString("No such player.\r\n"))
			return
		} else {
			d.Olc.House.Builtby = id
			send_to_char(d.Character, libc.CString("Builder changed.\r\n"))
		}
	case HSEDIT_PAYMENT:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.House.Last_payment = libc.GetTime(nil)
		case 'n':
			fallthrough
		case 'N':
			send_to_char(d.Character, libc.CString("Last Payment Date not changed\r\n"))
			hsedit_disp_menu(d)
			return
		}
	case HSEDIT_GUEST_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			if d.Olc.House.Num_of_guests > (int(MAX_GUESTS - 1)) {
				send_to_char(d.Character, libc.CString("Guest List Full! - delete some before adding more\r\nEnter selection (A/D/C/Q) : "))
			} else {
				send_to_char(d.Character, libc.CString("Name of guest to add: "))
				d.Olc.Mode = HSEDIT_GUEST_ADD
			}
		case 'd':
			fallthrough
		case 'D':
			if d.Olc.House.Num_of_guests < 1 {
				send_to_char(d.Character, libc.CString("Guest List Empty! - add a guest before trying to delete one\r\nEnter selection (A/D/C/Q) : "))
			} else {
				send_to_char(d.Character, libc.CString("Name of guest to delete : "))
				d.Olc.Mode = HSEDIT_GUEST_DELETE
			}
		case 'c':
			fallthrough
		case 'C':
			send_to_char(d.Character, libc.CString("Clear guest list? (Y/N) : "))
			d.Olc.Mode = HSEDIT_GUEST_CLEAR
		case 'q':
			fallthrough
		case 'Q':
			hsedit_disp_menu(d)
		default:
			send_to_char(d.Character, libc.CString("Invalid choice!\r\n\r\n"))
			hsedit_disp_guest_menu(d)
		}
		return
	case HSEDIT_GUEST_ADD:
		if (func() int {
			id = get_id_by_name(arg)
			return id
		}()) < 0 {
			send_to_char(d.Character, libc.CString("No such player.\r\n"))
			hsedit_disp_guest_menu(d)
			return
		} else if id == int(d.Character.Idnum) {
			send_to_char(d.Character, libc.CString("House owner should not be in the guest list!\r\n"))
			hsedit_disp_guest_menu(d)
			return
		} else {
			for i = 0; i < d.Olc.House.Num_of_guests; i++ {
				if d.Olc.House.Guests[i] == id {
					send_to_char(d.Character, libc.CString("That player is already in the guest list!.\r\n"))
					hsedit_disp_guest_menu(d)
					return
				}
			}
			i = func() int {
				p := &d.Olc.House.Num_of_guests
				x := *p
				*p++
				return x
			}()
			d.Olc.House.Guests[i] = id
			send_to_char(d.Character, libc.CString("Guest added.\r\n"))
			hsedit_disp_guest_menu(d)
			return
		}
	case HSEDIT_GUEST_DELETE:
		if (func() int {
			id = get_id_by_name(arg)
			return id
		}()) < 0 {
			send_to_char(d.Character, libc.CString("No such player.\r\n"))
			hsedit_disp_guest_menu(d)
			return
		} else if id == int(d.Character.Idnum) {
			send_to_char(d.Character, libc.CString("House owner should not be in the guest list!\r\n"))
			hsedit_disp_guest_menu(d)
			return
		} else {
			for i = 0; i < d.Olc.House.Num_of_guests; i++ {
				if d.Olc.House.Guests[i] == id {
					for ; i < d.Olc.House.Num_of_guests; i++ {
						d.Olc.House.Guests[i] = d.Olc.House.Guests[i+1]
					}
					d.Olc.House.Num_of_guests--
					send_to_char(d.Character, libc.CString("Guest deleted.\r\n"))
					d.Olc.Value = 1
					hsedit_disp_guest_menu(d)
					return
				}
			}
			send_to_char(d.Character, libc.CString("That player isn't in the guest list!\r\n"))
			hsedit_disp_guest_menu(d)
			return
		}
	case HSEDIT_GUEST_CLEAR:
		switch *arg {
		case 'n':
			fallthrough
		case 'N':
			send_to_char(d.Character, libc.CString("Invalid choice!"))
			return
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.House.Num_of_guests = 0
			for i = 0; i < MAX_GUESTS; i++ {
				d.Olc.House.Guests[i] = 0
			}
			send_to_char(d.Character, libc.CString("Guest List Cleared!"))
			d.Olc.Value = 1
			hsedit_disp_guest_menu(d)
		default:
			send_to_char(d.Character, libc.CString("Invalid choice!\r\nClear Guest List? (Y/N) : "))
			return
		}
	case HSEDIT_TYPE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number >= NUM_HOUSE_TYPES {
			send_to_char(d.Character, libc.CString("Invalid choice!"))
			hsedit_disp_type_menu(d)
			return
		} else {
			d.Olc.House.Mode = number
		}
	case HSEDIT_FLAGS:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > HOUSE_NUM_FLAGS {
			send_to_char(d.Character, libc.CString("That's not a valid choice!\r\n"))
			hsedit_disp_flags_menu(d)
		} else {
			if number == 0 {
				break
			} else {
				if (d.Olc.House.Bitvector & (1 << (number - 1))) != 0 {
					d.Olc.House.Bitvector &= ^(1 << (number - 1))
				} else {
					d.Olc.House.Bitvector |= 1 << (number - 1)
				}
				hsedit_disp_flags_menu(d)
			}
		}
		return
	default:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Reached default case in parse_hsedit"))
	}
	d.Olc.Value = 1
	hsedit_disp_menu(d)
}
func hsedit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	}
}
func do_oasis_hsedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int = int(-1)
		save     int = 0
		real_num int
		d        *descriptor_data
		buf3     *byte
	)
	_ = buf3
	var buf1 [64936]byte
	var buf2 [64936]byte
	buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			number = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number)
		} else {
			number = -1
		}
	} else if !unicode.IsDigit(rune(buf1[0])) {
		if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) != 0 {
			send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
			return
		}
		save = TRUE
		if is_number(&buf2[0]) != 0 {
			number = libc.Atoi(libc.GoString(&buf2[0]))
		} else if ch.Player_specials.Olc_zone > 0 {
			var zlok zone_rnum
			if (func() zone_rnum {
				zlok = real_zone(zone_vnum(ch.Player_specials.Olc_zone))
				return zlok
			}()) == zone_rnum(-1) {
				number = -1
			} else {
				number = int(genolc_zone_bottom(zlok))
			}
		}
		if number == int(-1) {
			send_to_char(ch, libc.CString("Save which zone?\r\n"))
			return
		}
	}
	if number == int(-1) {
		number = libc.Atoi(libc.GoString(&buf1[0]))
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_HSEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That house is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: do_oasis_hsedit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if save != 0 {
		d.Olc.Zone_num = real_zone(zone_vnum(number))
	} else {
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(number))
	}
	if d.Olc.Zone_num == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for that number!\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 {
		send_to_char(ch, libc.CString(" You do not have permission to edit zone %d. Try zone %d.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("OLC: %s tried to edit zone %d allowed zone %d"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all houses in zone %d.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves house info for zone %d."), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		hsedit_save_to_disk()
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	real_num = find_house(room_vnum(number))
	if real_num == int(-1) {
		if num_of_houses >= MAX_HOUSES {
			send_to_char(ch, libc.CString("MAX houses limit reached (%d) - Unable to create more.\r\n"), MAX_HOUSES)
			mudlog(NRM, ADMLVL_BUILDER, TRUE, libc.CString("HSEDIT: MAX houses limit reached (%d)\r\n"), MAX_HOUSES)
			return
		} else {
			hsedit_setup_new(d)
		}
	} else {
		hsedit_setup_existing(d, real_num)
	}
	d.Connected = CON_HSEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	mudlog(CMP, ADMLVL_BUILDER, TRUE, libc.CString("OLC: (hsedit) %s starts editing zone %d allowed zone %d"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
}
