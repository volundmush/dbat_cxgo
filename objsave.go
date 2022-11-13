package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

const RENT_FACTOR = 1
const CRYO_FACTOR = 4
const LOC_INVENTORY = 0
const MAX_BAG_ROWS = 5

func delete_inv_backup(ch *char_data) {
	var (
		source      *stdio.File
		source_file [20480]byte
		alpha       [2048]byte
		name        [2048]byte
	)
	stdio.Sprintf(&name[0], libc.GoString(GET_NAME(ch)))
	if name[0] == 'a' || name[0] == 'A' || name[0] == 'b' || name[0] == 'B' || name[0] == 'c' || name[0] == 'C' || name[0] == 'd' || name[0] == 'D' || name[0] == 'e' || name[0] == 'E' {
		stdio.Sprintf(&alpha[0], "A-E")
	} else if name[0] == 'f' || name[0] == 'F' || name[0] == 'g' || name[0] == 'G' || name[0] == 'h' || name[0] == 'H' || name[0] == 'i' || name[0] == 'I' || name[0] == 'j' || name[0] == 'J' {
		stdio.Sprintf(&alpha[0], "F-J")
	} else if name[0] == 'k' || name[0] == 'K' || name[0] == 'l' || name[0] == 'L' || name[0] == 'm' || name[0] == 'M' || name[0] == 'n' || name[0] == 'N' || name[0] == 'o' || name[0] == 'O' {
		stdio.Sprintf(&alpha[0], "K-O")
	} else if name[0] == 'p' || name[0] == 'P' || name[0] == 'q' || name[0] == 'Q' || name[0] == 'r' || name[0] == 'R' || name[0] == 's' || name[0] == 'S' || name[0] == 't' || name[0] == 'T' {
		stdio.Sprintf(&alpha[0], "P-T")
	} else if name[0] == 'u' || name[0] == 'U' || name[0] == 'v' || name[0] == 'V' || name[0] == 'w' || name[0] == 'W' || name[0] == 'x' || name[0] == 'X' || name[0] == 'y' || name[0] == 'Y' || name[0] == 'z' || name[0] == 'Z' {
		stdio.Sprintf(&alpha[0], "U-Z")
	}
	stdio.Sprintf(&source_file[0], "/home/m053car2/dbat/lib/plrobjs/%s/%s.copy", &alpha[0], ch.Name)
	if (func() *stdio.File {
		source = stdio.FOpen(libc.GoString(&source_file[0]), "r")
		return source
	}()) == nil {
		return
	}
	source.Close()
	if stdio.Remove(libc.GoString(&source_file[0])) < 0 && libc.Errno != ENOENT {
		basic_mud_log(libc.CString("ERROR: Couldn't delete backup inv."))
	}
	return
}
func load_inv_backup(ch *char_data) int {
	if GET_LEVEL(ch) < 2 {
		return -1
	}
	var chx int8
	var source *stdio.File
	var target *stdio.File
	var source_file [20480]byte
	var target_file [20480]byte
	var buf2 [20480]byte
	var alpha [2048]byte
	var name [2048]byte
	stdio.Sprintf(&name[0], libc.GoString(GET_NAME(ch)))
	if name[0] == 'a' || name[0] == 'A' || name[0] == 'b' || name[0] == 'B' || name[0] == 'c' || name[0] == 'C' || name[0] == 'd' || name[0] == 'D' || name[0] == 'e' || name[0] == 'E' {
		stdio.Sprintf(&alpha[0], "A-E")
	} else if name[0] == 'f' || name[0] == 'F' || name[0] == 'g' || name[0] == 'G' || name[0] == 'h' || name[0] == 'H' || name[0] == 'i' || name[0] == 'I' || name[0] == 'j' || name[0] == 'J' {
		stdio.Sprintf(&alpha[0], "F-J")
	} else if name[0] == 'k' || name[0] == 'K' || name[0] == 'l' || name[0] == 'L' || name[0] == 'm' || name[0] == 'M' || name[0] == 'n' || name[0] == 'N' || name[0] == 'o' || name[0] == 'O' {
		stdio.Sprintf(&alpha[0], "K-O")
	} else if name[0] == 'p' || name[0] == 'P' || name[0] == 'q' || name[0] == 'Q' || name[0] == 'r' || name[0] == 'R' || name[0] == 's' || name[0] == 'S' || name[0] == 't' || name[0] == 'T' {
		stdio.Sprintf(&alpha[0], "P-T")
	} else if name[0] == 'u' || name[0] == 'U' || name[0] == 'v' || name[0] == 'V' || name[0] == 'w' || name[0] == 'W' || name[0] == 'x' || name[0] == 'X' || name[0] == 'y' || name[0] == 'Y' || name[0] == 'z' || name[0] == 'Z' {
		stdio.Sprintf(&alpha[0], "U-Z")
	}
	stdio.Sprintf(&source_file[0], "/home/m053car2/dbat/lib/plrobjs/%s/%s.copy", &alpha[0], ch.Name)
	if get_filename(&buf2[0], uint64(20480), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return -1
	}
	stdio.Sprintf(&target_file[0], "/home/m053car2/dbat/lib/%s", &buf2[0])
	if (func() *stdio.File {
		source = stdio.FOpen(libc.GoString(&source_file[0]), "r")
		return source
	}()) == nil {
		basic_mud_log(libc.CString("Source in load_inv_backup failed to load."))
		basic_mud_log(&source_file[0])
		return -1
	}
	if (func() *stdio.File {
		target = stdio.FOpen(libc.GoString(&target_file[0]), "w")
		return target
	}()) == nil {
		basic_mud_log(libc.CString("Target in load_inv_backup failed to load."))
		basic_mud_log(&target_file[0])
		return -1
	}
	for int(func() int8 {
		chx = int8(source.GetC())
		return chx
	}()) != stdio.EOF {
		target.PutC(int(chx))
	}
	basic_mud_log(libc.CString("Inventory backup restore successful."))
	source.Close()
	target.Close()
	return 1
}
func inv_backup(ch *char_data) int {
	var (
		backup *stdio.File
		buf    [20480]byte
		alpha  [2048]byte
		name   [2048]byte
	)
	stdio.Sprintf(&name[0], libc.GoString(GET_NAME(ch)))
	if name[0] == 'a' || name[0] == 'A' || name[0] == 'b' || name[0] == 'B' || name[0] == 'c' || name[0] == 'C' || name[0] == 'd' || name[0] == 'D' || name[0] == 'e' || name[0] == 'E' {
		stdio.Sprintf(&alpha[0], "A-E")
	} else if name[0] == 'f' || name[0] == 'F' || name[0] == 'g' || name[0] == 'G' || name[0] == 'h' || name[0] == 'H' || name[0] == 'i' || name[0] == 'I' || name[0] == 'j' || name[0] == 'J' {
		stdio.Sprintf(&alpha[0], "F-J")
	} else if name[0] == 'k' || name[0] == 'K' || name[0] == 'l' || name[0] == 'L' || name[0] == 'm' || name[0] == 'M' || name[0] == 'n' || name[0] == 'N' || name[0] == 'o' || name[0] == 'O' {
		stdio.Sprintf(&alpha[0], "K-O")
	} else if name[0] == 'p' || name[0] == 'P' || name[0] == 'q' || name[0] == 'Q' || name[0] == 'r' || name[0] == 'R' || name[0] == 's' || name[0] == 'S' || name[0] == 't' || name[0] == 'T' {
		stdio.Sprintf(&alpha[0], "P-T")
	} else if name[0] == 'u' || name[0] == 'U' || name[0] == 'v' || name[0] == 'V' || name[0] == 'w' || name[0] == 'W' || name[0] == 'x' || name[0] == 'X' || name[0] == 'y' || name[0] == 'Y' || name[0] == 'z' || name[0] == 'Z' {
		stdio.Sprintf(&alpha[0], "U-Z")
	}
	stdio.Sprintf(&buf[0], "/home/m053car2/dbat/lib/plrobjs/%s/%s.copy", &alpha[0], ch.Name)
	if (func() *stdio.File {
		backup = stdio.FOpen(libc.GoString(&buf[0]), "r")
		return backup
	}()) == nil {
		return -1
	}
	backup.Close()
	return 1
}
func cp(ch *char_data) int {
	var (
		chx         int8
		source      *stdio.File
		target      *stdio.File
		source_file [20480]byte
		target_file [20480]byte
		buf2        [20480]byte
		alpha       [2048]byte
		name        [2048]byte
	)
	stdio.Sprintf(&name[0], libc.GoString(GET_NAME(ch)))
	if name[0] == 'a' || name[0] == 'A' || name[0] == 'b' || name[0] == 'B' || name[0] == 'c' || name[0] == 'C' || name[0] == 'd' || name[0] == 'D' || name[0] == 'e' || name[0] == 'E' {
		stdio.Sprintf(&alpha[0], "A-E")
	} else if name[0] == 'f' || name[0] == 'F' || name[0] == 'g' || name[0] == 'G' || name[0] == 'h' || name[0] == 'H' || name[0] == 'i' || name[0] == 'I' || name[0] == 'j' || name[0] == 'J' {
		stdio.Sprintf(&alpha[0], "F-J")
	} else if name[0] == 'k' || name[0] == 'K' || name[0] == 'l' || name[0] == 'L' || name[0] == 'm' || name[0] == 'M' || name[0] == 'n' || name[0] == 'N' || name[0] == 'o' || name[0] == 'O' {
		stdio.Sprintf(&alpha[0], "K-O")
	} else if name[0] == 'p' || name[0] == 'P' || name[0] == 'q' || name[0] == 'Q' || name[0] == 'r' || name[0] == 'R' || name[0] == 's' || name[0] == 'S' || name[0] == 't' || name[0] == 'T' {
		stdio.Sprintf(&alpha[0], "P-T")
	} else if name[0] == 'u' || name[0] == 'U' || name[0] == 'v' || name[0] == 'V' || name[0] == 'w' || name[0] == 'W' || name[0] == 'x' || name[0] == 'X' || name[0] == 'y' || name[0] == 'Y' || name[0] == 'z' || name[0] == 'Z' {
		stdio.Sprintf(&alpha[0], "U-Z")
	}
	stdio.Sprintf(&target_file[0], "/home/m053car2/dbat/lib/plrobjs/%s/%s.copy", &alpha[0], ch.Name)
	if get_filename(&buf2[0], uint64(20480), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return -1
	}
	stdio.Sprintf(&source_file[0], "/home/m053car2/dbat/lib/%s", &buf2[0])
	if (func() *stdio.File {
		source = stdio.FOpen(libc.GoString(&source_file[0]), "r")
		return source
	}()) == nil {
		basic_mud_log(libc.CString("Source failed to load."))
		basic_mud_log(&source_file[0])
		return -1
	}
	if (func() *stdio.File {
		target = stdio.FOpen(libc.GoString(&target_file[0]), "w")
		return target
	}()) == nil {
		basic_mud_log(libc.CString("Target failed to load."))
		basic_mud_log(&target_file[0])
		return -1
	}
	for int(func() int8 {
		chx = int8(source.GetC())
		return chx
	}()) != stdio.EOF {
		target.PutC(int(chx))
	}
	source.Close()
	target.Close()
	return 1
}
func Obj_to_store(obj *obj_data, fl *stdio.File, location int) int {
	my_obj_save_to_disk(fl, obj, location)
	return 1
}
func auto_equip(ch *char_data, obj *obj_data, location int) {
	var j int
	if location > 0 {
		switch func() int {
			j = location - 1
			return j
		}() {
		case WEAR_UNUSED0:
			j = WEAR_WIELD2
		case WEAR_FINGER_R:
			fallthrough
		case WEAR_FINGER_L:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_FINGER) {
				location = LOC_INVENTORY
			}
		case WEAR_NECK_1:
			fallthrough
		case WEAR_NECK_2:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_NECK) {
				location = LOC_INVENTORY
			}
		case WEAR_BODY:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_BODY) {
				location = LOC_INVENTORY
			}
		case WEAR_HEAD:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_HEAD) {
				location = LOC_INVENTORY
			}
		case WEAR_LEGS:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_LEGS) {
				location = LOC_INVENTORY
			}
		case WEAR_FEET:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_FEET) {
				location = LOC_INVENTORY
			}
		case WEAR_HANDS:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_HANDS) {
				location = LOC_INVENTORY
			}
		case WEAR_ARMS:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_ARMS) {
				location = LOC_INVENTORY
			}
		case WEAR_UNUSED1:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_SHIELD) {
				location = LOC_INVENTORY
			}
			j = WEAR_WIELD2
		case WEAR_ABOUT:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_ABOUT) {
				location = LOC_INVENTORY
			}
		case WEAR_WAIST:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_WAIST) {
				location = LOC_INVENTORY
			}
		case WEAR_WRIST_R:
			fallthrough
		case WEAR_WRIST_L:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_WRIST) {
				location = LOC_INVENTORY
			}
		case WEAR_WIELD1:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_WIELD) {
				location = LOC_INVENTORY
			}
		case WEAR_WIELD2:
		case WEAR_EYE:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_EYE) {
				location = LOC_INVENTORY
			}
		case WEAR_BACKPACK:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_PACK) {
				location = LOC_INVENTORY
			}
		case WEAR_SH:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_SH) {
				location = LOC_INVENTORY
			}
		case WEAR_EAR_R:
			fallthrough
		case WEAR_EAR_L:
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_EAR) {
				location = LOC_INVENTORY
			}
		default:
			location = LOC_INVENTORY
		}
		if location > 0 {
			if (ch.Equipment[j]) == nil {
				if invalid_align(ch, obj) != 0 || invalid_class(ch, obj) != 0 {
					location = LOC_INVENTORY
				} else {
					equip_char(ch, obj, j)
				}
			} else {
				mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: autoeq: '%s' already equipped in position %d."), GET_NAME(ch), location)
				location = LOC_INVENTORY
			}
		}
	}
	if location <= 0 {
		obj_to_char(obj, ch)
	}
}
func Crash_delete_file(name *byte) int {
	var (
		filename [50]byte
		fl       *stdio.File
	)
	if get_filename(&filename[0], uint64(50), NEW_OBJ_FILES, name) == 0 {
		return 0
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&filename[0]), "rb")
		return fl
	}()) == nil {
		if libc.Errno != ENOENT {
			basic_mud_log(libc.CString("SYSERR: deleting crash file %s (1): %s"), &filename[0], libc.StrError(libc.Errno))
		}
		return 0
	}
	fl.Close()
	if stdio.Remove(libc.GoString(&filename[0])) < 0 && libc.Errno != ENOENT {
		basic_mud_log(libc.CString("SYSERR: deleting crash file %s (2): %s"), &filename[0], libc.StrError(libc.Errno))
	}
	return 1
}
func Crash_delete_crashfile(ch *char_data) int {
	var (
		filename [2048]byte
		fl       *stdio.File
		rentcode int
		timed    int
		netcost  int
		gold     int
		account  int
		nitems   int
		line     [2048]byte
	)
	if get_filename(&filename[0], uint64(2048), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return 0
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&filename[0]), "rb")
		return fl
	}()) == nil {
		if libc.Errno != ENOENT {
			basic_mud_log(libc.CString("SYSERR: checking for crash file %s (3): %s"), &filename[0], libc.StrError(libc.Errno))
		}
		return 0
	}
	if int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
	}
	stdio.Sscanf(&line[0], "%d %d %d %d %d %d", &rentcode, &timed, &netcost, &gold, &account, &nitems)
	fl.Close()
	if rentcode == RENT_CRASH {
		Crash_delete_file(GET_NAME(ch))
	}
	return 1
}
func Crash_clean_file(name *byte) int {
	var (
		filename [64936]byte
		fl       *stdio.File
		rentcode int
		timed    int
		netcost  int
		gold     int
		account  int
		nitems   int
		line     [64936]byte
	)
	if get_filename(&filename[0], uint64(64936), NEW_OBJ_FILES, name) == 0 {
		return 0
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&filename[0]), "r+b")
		return fl
	}()) == nil {
		if libc.Errno != ENOENT {
			basic_mud_log(libc.CString("SYSERR: OPENING OBJECT FILE %s (4): %s"), &filename[0], libc.StrError(libc.Errno))
		}
		return 0
	}
	if int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%d %d %d %d %d %d", &rentcode, &timed, &netcost, &gold, &account, &nitems)
		fl.Close()
		if rentcode == RENT_CRASH || rentcode == RENT_FORCED || rentcode == RENT_TIMEDOUT {
			if timed < int(libc.GetTime(nil))-config_info.Csd.Crash_file_timeout*((int(SECS_PER_REAL_MIN*60))*24) {
				var filetype *byte
				Crash_delete_file(name)
				switch rentcode {
				case RENT_CRASH:
					filetype = libc.CString("crash")
				case RENT_FORCED:
					filetype = libc.CString("forced rent")
				case RENT_TIMEDOUT:
					filetype = libc.CString("idlesave")
				default:
					filetype = libc.CString("UNKNOWN!")
				}
				basic_mud_log(libc.CString("    Deleting %s's %s file."), name, filetype)
				return 1
			}
		} else if rentcode == RENT_RENTED {
			if timed < int(libc.GetTime(nil))-config_info.Csd.Rent_file_timeout*((int(SECS_PER_REAL_MIN*60))*24) {
				Crash_delete_file(name)
				basic_mud_log(libc.CString("    Deleting %s's rent file."), name)
				return 1
			}
		}
	}
	return 0
}
func update_obj_file() {
	var i int
	for i = 0; i <= top_of_p_table; i++ {
		if *(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name != 0 {
			Crash_clean_file((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name)
		}
	}
}
func Crash_listrent(ch *char_data, name *byte) {
	var (
		fl       *stdio.File = nil
		filename [64936]byte
		buf      [64936]byte
		obj      *obj_data
		rentcode int
		timed    int
		netcost  int
		gold     int
		account  int
		nitems   int
		len_     int
		t        [10]int
		nr       int
		line     [64936]byte
		sdesc    *byte
	)
	if get_filename(&filename[0], uint64(64936), NEW_OBJ_FILES, name) != 0 {
		fl = stdio.FOpen(libc.GoString(&filename[0]), "rb")
	}
	if fl == nil {
		send_to_char(ch, libc.CString("%s has no rent file.\r\n"), name)
		return
	}
	send_to_char(ch, libc.CString("%s\r\n"), &filename[0])
	if int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%d %d %d %d %d %d", &rentcode, &timed, &netcost, &gold, &account, &nitems)
	}
	switch rentcode {
	case RENT_RENTED:
		send_to_char(ch, libc.CString("Rent\r\n"))
	case RENT_CRASH:
		send_to_char(ch, libc.CString("Crash\r\n"))
	case RENT_CRYO:
		send_to_char(ch, libc.CString("Cryo\r\n"))
	case RENT_TIMEDOUT:
		fallthrough
	case RENT_FORCED:
		send_to_char(ch, libc.CString("TimedOut\r\n"))
	default:
		send_to_char(ch, libc.CString("Undef\r\n"))
	}
	buf[0] = 0
	len_ = 0
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		if line[0] == '#' {
			stdio.Sscanf(&line[0], "#%d", &nr)
			if nr != int(-1) {
				if real_object(obj_vnum(nr)) != obj_rnum(-1) {
					obj = read_object(obj_vnum(nr), VIRTUAL)
					if len_+math.MaxUint8 < int(64936) {
						len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[%5d] (%5dau) %-20s\r\n", nr, obj.Cost_per_day, obj.Short_description)
					} else {
						stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "** Excessive rent listing. **\r\n")
						break
					}
					extract_obj(obj)
				} else {
					if len_+math.MaxUint8 < int(64936) {
						len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%s[-----] NONEXISTANT OBJECT #%d\r\n", &buf[0], nr)
					} else {
						stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "** Excessive rent listing. **\r\n")
						break
					}
				}
			} else {
				get_line(fl, &line[0])
				get_line(fl, &line[0])
				fread_string(fl, libc.CString(", listrent reading name"))
				sdesc = fread_string(fl, libc.CString(", listrent reading sdesc"))
				fread_string(fl, libc.CString(", listrent reading desc"))
				fread_string(fl, libc.CString(", listrent reading adesc"))
				get_line(fl, &line[0])
				stdio.Sscanf(&line[0], "%d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4])
				if len_+math.MaxUint8 < int(64936) {
					len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%s[%5d] (%5dau) %-20s\r\n", &buf[0], nr, t[4], sdesc)
				} else {
					stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "** Excessive rent listing. **\r\n")
					break
				}
			}
		}
	}
	page_string(ch.Desc, &buf[0], 0)
	fl.Close()
}
func Crash_save(obj *obj_data, fp *stdio.File, location int) int {
	var (
		tmp    *obj_data
		result int
	)
	if obj != nil {
		Crash_save(obj.Next_content, fp, location)
		Crash_save(obj.Contains, fp, MIN(0, location)-1)
		result = Obj_to_store(obj, fp, location)
		for tmp = obj.In_obj; tmp != nil; tmp = tmp.In_obj {
			tmp.Weight -= obj.Weight
		}
		if result == 0 {
			return FALSE
		}
	}
	return TRUE
}
func Crash_restore_weight(obj *obj_data) {
	if obj != nil {
		Crash_restore_weight(obj.Contains)
		Crash_restore_weight(obj.Next_content)
		if obj.In_obj != nil {
			obj.In_obj.Weight += obj.Weight
		}
	}
}
func Crash_extract_norent_eq(ch *char_data) {
	var j int
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) == nil {
			continue
		}
		if Crash_is_unrentable(ch.Equipment[j]) != 0 {
			obj_to_char(unequip_char(ch, j), ch)
		} else {
			Crash_extract_norents(ch.Equipment[j])
		}
	}
}
func Crash_extract_objs(obj *obj_data) {
	if obj != nil {
		Crash_extract_objs(obj.Contains)
		Crash_extract_objs(obj.Next_content)
		extract_obj(obj)
	}
}
func Crash_is_unrentable(obj *obj_data) int {
	if obj == nil {
		return 0
	}
	if OBJ_FLAGGED(obj, ITEM_NORENT) || obj.Cost_per_day < 0 || obj.Item_number <= obj_vnum(-1) && !IS_SET_AR(obj.Extra_flags[:], ITEM_UNIQUE_SAVE) {
		return 1
	}
	return 0
}
func Crash_extract_norents(obj *obj_data) {
	if obj != nil {
		Crash_extract_norents(obj.Contains)
		Crash_extract_norents(obj.Next_content)
		if Crash_is_unrentable(obj) != 0 {
			extract_obj(obj)
		}
	}
}
func Crash_extract_expensive(obj *obj_data) {
	var (
		tobj *obj_data
		max  *obj_data
	)
	max = obj
	for tobj = obj; tobj != nil; tobj = tobj.Next_content {
		if tobj.Cost_per_day > max.Cost_per_day {
			max = tobj
		}
	}
	extract_obj(max)
}
func Crash_calculate_rent(obj *obj_data, cost *int) {
	if obj != nil {
		*cost += MAX(0, obj.Cost_per_day)
		Crash_calculate_rent(obj.Contains, cost)
		Crash_calculate_rent(obj.Next_content, cost)
	}
}
func Crash_crashsave(ch *char_data) {
	var (
		buf [2048]byte
		j   int
		fp  *stdio.File
	)
	if IS_NPC(ch) {
		return
	}
	if get_filename(&buf[0], uint64(2048), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fp
	}()) == nil {
		return
	}
	stdio.Fprintf(fp, "%d %d %d %d %d %d\r\n", RENT_CRASH, int(libc.GetTime(nil)), 0, ch.Gold, ch.Bank_gold, 0)
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) != nil {
			if Crash_save(ch.Equipment[j], fp, j+1) == 0 {
				fp.Close()
				return
			}
			Crash_restore_weight(ch.Equipment[j])
		}
	}
	if Crash_save(ch.Carrying, fp, 0) == 0 {
		fp.Close()
		return
	}
	Crash_restore_weight(ch.Carrying)
	fp.Close()
	ch.Act[int(PLR_CRASH/32)] &= bitvector_t(int32(^(1 << (int(PLR_CRASH % 32)))))
}
func Crash_idlesave(ch *char_data) {
	var (
		buf     [2048]byte
		j       int
		cost    int
		cost_eq int
		fp      *stdio.File
	)
	if IS_NPC(ch) {
		return
	}
	if get_filename(&buf[0], uint64(2048), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fp
	}()) == nil {
		return
	}
	Crash_extract_norent_eq(ch)
	Crash_extract_norents(ch.Carrying)
	cost = 0
	Crash_calculate_rent(ch.Carrying, &cost)
	cost_eq = 0
	for j = 0; j < NUM_WEARS; j++ {
		Crash_calculate_rent(ch.Equipment[j], &cost_eq)
	}
	cost += cost_eq
	cost *= 2
	if cost > ch.Gold+ch.Bank_gold {
		for j = 0; j < NUM_WEARS; j++ {
			if (ch.Equipment[j]) != nil {
				obj_to_char(unequip_char(ch, j), ch)
			}
		}
		for cost > ch.Gold+ch.Bank_gold && ch.Carrying != nil {
			Crash_extract_expensive(ch.Carrying)
			cost = 0
			Crash_calculate_rent(ch.Carrying, &cost)
			cost *= 2
		}
	}
	if ch.Carrying == nil {
		for j = 0; j < NUM_WEARS && (ch.Equipment[j]) == nil; j++ {
		}
		if j == NUM_WEARS {
			fp.Close()
			Crash_delete_file(GET_NAME(ch))
			return
		}
	}
	stdio.Fprintf(fp, "%d %d %d %d %d %d\r\n", RENT_TIMEDOUT, int(libc.GetTime(nil)), cost, ch.Gold, ch.Bank_gold, 0)
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) != nil {
			if Crash_save(ch.Equipment[j], fp, j+1) == 0 {
				fp.Close()
				return
			}
			Crash_restore_weight(ch.Equipment[j])
			Crash_extract_objs(ch.Equipment[j])
		}
	}
	if Crash_save(ch.Carrying, fp, 0) == 0 {
		fp.Close()
		return
	}
	fp.Close()
	Crash_extract_objs(ch.Carrying)
}
func Crash_rentsave(ch *char_data, cost int) {
	var (
		buf [2048]byte
		j   int
		fp  *stdio.File
	)
	if IS_NPC(ch) {
		return
	}
	if get_filename(&buf[0], uint64(2048), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fp
	}()) == nil {
		return
	}
	Crash_extract_norent_eq(ch)
	Crash_extract_norents(ch.Carrying)
	stdio.Fprintf(fp, "%d %d %d %d %d %d\r\n", RENT_RENTED, int(libc.GetTime(nil)), cost, ch.Gold, ch.Bank_gold, 0)
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) != nil {
			if Crash_save(ch.Equipment[j], fp, j+1) == 0 {
				fp.Close()
				return
			}
			Crash_restore_weight(ch.Equipment[j])
			Crash_extract_objs(ch.Equipment[j])
		}
	}
	if Crash_save(ch.Carrying, fp, 0) == 0 {
		fp.Close()
		return
	}
	fp.Close()
	Crash_extract_objs(ch.Carrying)
}
func Crash_cryosave(ch *char_data, cost int) {
	var (
		buf [2048]byte
		j   int
		fp  *stdio.File
	)
	if IS_NPC(ch) {
		return
	}
	if get_filename(&buf[0], uint64(2048), CRASH_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fp
	}()) == nil {
		return
	}
	Crash_extract_norent_eq(ch)
	Crash_extract_norents(ch.Carrying)
	ch.Gold = MAX(0, ch.Gold-cost)
	stdio.Fprintf(fp, "%d %d %d %d %d %d\r\n", RENT_CRYO, int(libc.GetTime(nil)), 0, ch.Gold, ch.Bank_gold, 0)
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) != nil {
			if Crash_save(ch.Equipment[j], fp, j+1) == 0 {
				fp.Close()
				return
			}
			Crash_restore_weight(ch.Equipment[j])
			Crash_extract_objs(ch.Equipment[j])
		}
	}
	if Crash_save(ch.Carrying, fp, 0) == 0 {
		fp.Close()
		return
	}
	fp.Close()
	Crash_extract_objs(ch.Carrying)
	ch.Act[int(PLR_CRYO/32)] |= bitvector_t(int32(1 << (int(PLR_CRYO % 32))))
}
func Crash_rent_deadline(ch *char_data, recep *char_data, cost int) {
	var (
		buf           [256]byte
		rent_deadline int
	)
	if cost == 0 {
		return
	}
	rent_deadline = (ch.Gold + ch.Bank_gold) / cost
	stdio.Snprintf(&buf[0], int(256), "$n tells you, 'You can rent for %ld day%s with the gold you have\r\non hand and in the bank.'\r\n", rent_deadline, func() string {
		if rent_deadline != 1 {
			return "s"
		}
		return ""
	}())
	act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
}
func Crash_report_unrentables(ch *char_data, recep *char_data, obj *obj_data) int {
	var has_norents int = 0
	if obj != nil {
		if Crash_is_unrentable(obj) != 0 {
			var buf [128]byte
			has_norents = 1
			stdio.Snprintf(&buf[0], int(128), "$n tells you, 'You cannot store %s.'", OBJS(obj, ch))
			act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		}
		has_norents += Crash_report_unrentables(ch, recep, obj.Contains)
		has_norents += Crash_report_unrentables(ch, recep, obj.Next_content)
	}
	return has_norents
}
func Crash_report_rent(ch *char_data, recep *char_data, obj *obj_data, cost *int, nitems *int, display int, factor int) {
	if obj != nil {
		if Crash_is_unrentable(obj) == 0 {
			(*nitems)++
			*cost += MAX(0, obj.Cost_per_day*factor)
			if display != 0 {
				var buf [256]byte
				stdio.Snprintf(&buf[0], int(256), "$n tells you, '%5d zenni for %s..'", obj.Cost_per_day*factor, OBJS(obj, ch))
				act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
			}
		}
		Crash_report_rent(ch, recep, obj.Contains, cost, nitems, display, factor)
		Crash_report_rent(ch, recep, obj.Next_content, cost, nitems, display, factor)
	}
}
func Crash_offer_rent(ch *char_data, recep *char_data, display int, factor int) int {
	var (
		i         int
		totalcost int = 0
		numitems  int = 0
		norent    int
	)
	norent = Crash_report_unrentables(ch, recep, ch.Carrying)
	for i = 0; i < NUM_WEARS; i++ {
		norent += Crash_report_unrentables(ch, recep, ch.Equipment[i])
	}
	if norent != 0 {
		return 0
	}
	totalcost = config_info.Csd.Min_rent_cost * factor
	Crash_report_rent(ch, recep, ch.Carrying, &totalcost, &numitems, display, factor)
	for i = 0; i < NUM_WEARS; i++ {
		Crash_report_rent(ch, recep, ch.Equipment[i], &totalcost, &numitems, display, factor)
	}
	if numitems == 0 {
		act(libc.CString("$n tells you, 'But you are not carrying anything!  Just quit!'"), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		return 0
	}
	if numitems > config_info.Csd.Max_obj_save {
		var buf [256]byte
		stdio.Snprintf(&buf[0], int(256), "$n tells you, 'Sorry, but I cannot store more than %d items.'", config_info.Csd.Max_obj_save)
		act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		return 0
	}
	if display != 0 {
		var buf [256]byte
		stdio.Snprintf(&buf[0], int(256), "$n tells you, 'Plus, my %d zenni fee..'", config_info.Csd.Min_rent_cost*factor)
		act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		stdio.Snprintf(&buf[0], int(256), "$n tells you, 'For a total of %ld zenni.'", totalcost)
		act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		if totalcost > ch.Gold+ch.Bank_gold {
			act(libc.CString("$n tells you, '...which I see you can't afford.'"), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
			return 0
		} else if factor == RENT_FACTOR {
			Crash_rent_deadline(ch, recep, totalcost)
		}
	}
	return totalcost
}
func gen_receptionist(ch *char_data, recep *char_data, cmd int, arg *byte, mode int) int {
	var (
		cost         int
		action_table [9]*byte = [9]*byte{libc.CString("smile"), libc.CString("dance"), libc.CString("sigh"), libc.CString("blush"), libc.CString("burp"), libc.CString("cough"), libc.CString("fart"), libc.CString("twiddle"), libc.CString("yawn")}
	)
	if cmd == 0 && rand_number(0, 5) == 0 {
		do_action(recep, nil, find_command(action_table[rand_number(0, 8)]), 0)
		return FALSE
	}
	if ch.Desc == nil || IS_NPC(ch) {
		return FALSE
	}
	if libc.StrCmp(libc.CString("offer"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) != 0 && libc.StrCmp(libc.CString("rent"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) != 0 {
		return FALSE
	}
	if !AWAKE(recep) {
		send_to_char(ch, libc.CString("%s is unable to talk to you...\r\n"), HSSH(recep))
		return TRUE
	}
	if !CAN_SEE(recep, ch) {
		act(libc.CString("$n says, 'I don't deal with people I can't see!'"), FALSE, recep, nil, nil, TO_ROOM)
		return TRUE
	}
	if config_info.Csd.Free_rent != 0 {
		act(libc.CString("$n tells you, 'Rent is free here.  Just quit, and your objects will be saved!'"), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		return 1
	}
	if libc.StrCmp(libc.CString("rent"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		var buf [128]byte
		if (func() int {
			cost = Crash_offer_rent(ch, recep, FALSE, mode)
			return cost
		}()) == 0 {
			return TRUE
		}
		if mode == RENT_FACTOR {
			stdio.Snprintf(&buf[0], int(128), "$n tells you, 'Rent will cost you %d zenni per day.'", cost)
		} else if mode == CRYO_FACTOR {
			stdio.Snprintf(&buf[0], int(128), "$n tells you, 'It will cost you %d zenni to be frozen.'", cost)
		}
		act(&buf[0], FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
		if cost > ch.Gold+ch.Bank_gold {
			act(libc.CString("$n tells you, '...which I see you can't afford.'"), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
			return TRUE
		}
		if cost != 0 && mode == RENT_FACTOR {
			Crash_rent_deadline(ch, recep, cost)
		}
		if mode == RENT_FACTOR {
			act(libc.CString("$n stores your belongings and helps you into your private chamber."), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
			Crash_rentsave(ch, cost)
			mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s has rented (%d/day, %d tot.)"), GET_NAME(ch), cost, ch.Gold+ch.Bank_gold)
		} else {
			act(libc.CString("$n stores your belongings and helps you into your private chamber.\r\nA white mist appears in the room, chilling you to the bone...\r\nYou begin to lose consciousness..."), FALSE, recep, nil, unsafe.Pointer(ch), TO_VICT)
			Crash_cryosave(ch, cost)
			mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s has cryo-rented."), GET_NAME(ch))
			ch.Act[int(PLR_CRYO/32)] |= bitvector_t(int32(1 << (int(PLR_CRYO % 32))))
		}
		act(libc.CString("$n helps $N into $S private chamber."), FALSE, recep, nil, unsafe.Pointer(ch), TO_NOTVICT)
		extract_char(ch)
	} else {
		Crash_offer_rent(ch, recep, TRUE, mode)
		act(libc.CString("$N gives $n an offer."), FALSE, ch, nil, unsafe.Pointer(recep), TO_ROOM)
	}
	return TRUE
}
func receptionist(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	return gen_receptionist(ch, (*char_data)(me), cmd, argument, RENT_FACTOR)
}
func cryogenicist(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	return gen_receptionist(ch, (*char_data)(me), cmd, argument, CRYO_FACTOR)
}
func Crash_save_all() {
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_PLAYING && !IS_NPC(d.Character) {
			if PLR_FLAGGED(d.Character, PLR_CRASH) {
				Crash_crashsave(d.Character)
				save_char(d.Character)
				d.Character.Act[int(PLR_CRASH/32)] &= bitvector_t(int32(^(1 << (int(PLR_CRASH % 32)))))
			}
		}
	}
}
func Crash_load(ch *char_data) int {
	var (
		fl      *stdio.File
		cmfname [64936]byte
		buf1    [64936]byte
		buf2    [64936]byte
		line    [256]byte
		t       [30]int
		danger  int
	)
	_ = danger
	var zwei int = 0
	var num_of_days int
	_ = num_of_days
	var orig_rent_code int
	var temp *obj_data
	var locate int = 0
	var j int
	var nr int
	var k int
	_ = k
	var cost int
	_ = cost
	var num_objs int = 0
	var obj1 *obj_data
	var cont_row [5]*obj_data
	var new_descr *extra_descr_data
	var rentcode int
	var timed int
	var netcost int
	var gold int
	var account int
	var nitems int
	var f1 [256]byte
	var f2 [256]byte
	var f3 [256]byte
	var f4 [256]byte
	if get_filename(&cmfname[0], uint64(64936), NEW_OBJ_FILES, GET_NAME(ch)) == 0 {
		return 1
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&cmfname[0]), "r+b")
		return fl
	}()) == nil {
		if libc.Errno != ENOENT {
			stdio.Sprintf(&buf1[0], "SYSERR: READING OBJECT FILE %s (5)", &cmfname[0])
			perror(&buf1[0])
			send_to_char(ch, libc.CString("\r\n********************* NOTICE *********************\r\nThere was a problem loading your objects from disk.\r\nContact a God for assistance.\r\n"))
		}
		if GET_LEVEL(ch) > 1 {
			mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s entering game with no equipment. Loading backup."), GET_NAME(ch))
		}
		if inv_backup(ch) == 0 {
			return -1
		} else {
			if load_inv_backup(ch) == 0 {
				return -1
			} else if (func() *stdio.File {
				fl = stdio.FOpen(libc.GoString(&cmfname[0]), "r+b")
				return fl
			}()) == nil {
				if libc.Errno != ENOENT {
					stdio.Sprintf(&buf1[0], "SYSERR: READING OBJECT FILE %s (5)", &cmfname[0])
					perror(&buf1[0])
					send_to_char(ch, libc.CString("\r\n********************* NOTICE *********************\r\nThere was a problem loading your objects from disk.\r\nContact a God for assistance.\r\n"))
				}
				return -1
			}
		}
	}
	if int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
	}
	stdio.Sscanf(&line[0], "%d %d %d %d %d %d", &rentcode, &timed, &netcost, &gold, &account, &nitems)
	if rentcode == RENT_RENTED || rentcode == RENT_TIMEDOUT {
		num_of_days = int(float32(int(libc.GetTime(nil))-timed) / float32((int(SECS_PER_REAL_MIN*60))*24))
		cost = 0
		save_char(ch)
	}
	switch func() int {
		orig_rent_code = rentcode
		return orig_rent_code
	}() {
	case RENT_RENTED:
		mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s un-renting and entering game."), GET_NAME(ch))
	case RENT_CRASH:
		mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s retrieving crash-saved items and entering game."), GET_NAME(ch))
	case RENT_CRYO:
		mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s un-cryo'ing and entering game."), GET_NAME(ch))
	case RENT_FORCED:
		fallthrough
	case RENT_TIMEDOUT:
		mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s retrieving force-saved items and entering game."), GET_NAME(ch))
	default:
		mudlog(BRF, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("WARNING: %s entering game with undefined rent code."), GET_NAME(ch))
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
				temp.Size = SIZE_MEDIUM
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
			temp.Extra_flags[3] = asciiflag_conv(&f4[0])
			temp.Value[8] = t[13]
			temp.Value[9] = t[14]
			temp.Value[10] = t[15]
			temp.Value[11] = t[16]
			temp.Value[12] = t[17]
			temp.Value[13] = t[18]
			temp.Value[14] = t[19]
			temp.Value[15] = t[20]
			get_line(fl, &line[0])
			if libc.StrCmp(libc.CString("XAP"), &line[0]) == 0 {
				if (func() *byte {
					p := &temp.Name
					temp.Name = fread_string(fl, libc.CString("rented object name"))
					return *p
				}()) == nil {
					temp.Name = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Short_description
					temp.Short_description = fread_string(fl, libc.CString("rented object short desc"))
					return *p
				}()) == nil {
					temp.Short_description = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Description
					temp.Description = fread_string(fl, libc.CString("rented object desc"))
					return *p
				}()) == nil {
					temp.Description = libc.CString("undefined")
				}
				if (func() *byte {
					p := &temp.Action_description
					temp.Action_description = fread_string(fl, libc.CString("rented object adesc"))
					return *p
				}()) == nil {
					temp.Action_description = nil
				}
				if get_line(fl, &line[0]) == 0 || stdio.Sscanf(&line[0], "%d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7]) != 8 {
					stdio.Fprintf(stdio.Stderr(), "Format error in first numeric line (expecting _x_ args)")
					return 0
				}
				temp.Type_flag = int8(t[0])
				temp.Wear_flags[0] = t[1]
				temp.Wear_flags[1] = t[2]
				temp.Wear_flags[2] = t[3]
				temp.Wear_flags[3] = t[4]
				temp.Weight = int64(t[5])
				temp.Cost = t[6]
				temp.Cost_per_day = t[7]
				for j = 0; j < MAX_OBJ_AFFECT; j++ {
					temp.Affected[j].Location = APPLY_NONE
					temp.Affected[j].Modifier = 0
					temp.Affected[j].Specific = 0
				}
				if int(temp.Type_flag) == ITEM_SPELLBOOK {
					if temp.Sbinfo == nil {
						temp.Sbinfo = &make([]obj_spellbook_spell, SPELLBOOK_SIZE)[0]
						libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(temp.Sbinfo))), 0, int(SPELLBOOK_SIZE*unsafe.Sizeof(obj_spellbook_spell{})))
					}
					for j = 0; j < SPELLBOOK_SIZE; j++ {
						(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(j)))).Spellname = 0
						(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(j)))).Pages = 0
					}
					(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*0))).Spellname = SPELL_DETECT_MAGIC
					(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*0))).Pages = 1
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
						stdio.Sprintf(&buf2[0], "rented object edesc keyword for object #%d", nr)
						new_descr.Keyword = fread_string(fl, &buf2[0])
						stdio.Sprintf(&buf2[0], "rented object edesc text for object #%d keyword %s", nr, new_descr.Keyword)
						new_descr.Description = fread_string(fl, &buf2[0])
						new_descr.Next = temp.Ex_description
						temp.Ex_description = new_descr
						get_line(fl, &line[0])
					case 'A':
						if j >= MAX_OBJ_AFFECT {
							basic_mud_log(libc.CString("SYSERR: Too many object affectations in loading rent file"))
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
							temp.Sbinfo = &make([]obj_spellbook_spell, SPELLBOOK_SIZE)[0]
							libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(temp.Sbinfo))), 0, int(SPELLBOOK_SIZE*unsafe.Sizeof(obj_spellbook_spell{})))
						}
						(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(j)))).Spellname = t[0]
						(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(temp.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(j)))).Pages = t[1]
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
				check_unique_id(temp)
				add_unique_id(temp)
				if int(temp.Type_flag) == ITEM_DRINKCON {
					name_from_drinkcon(temp)
					if (temp.Value[1]) != 0 {
						name_to_drinkcon(temp, temp.Value[2])
					}
				}
				if GET_OBJ_VNUM(temp) == 0x4E83 || GET_OBJ_VNUM(temp) == 0x4E82 {
					if OBJ_FLAGGED(temp, ITEM_UNBREAKABLE) {
						temp.Extra_flags[int(ITEM_UNBREAKABLE/32)] &= bitvector_t(int32(^(1 << (int(ITEM_UNBREAKABLE % 32)))))
					}
				}
				auto_equip(ch, temp, locate)
			} else {
				continue
			}
			if locate > 0 {
				for j = int(MAX_BAG_ROWS - 1); j > 0; j-- {
					if cont_row[j] != nil {
						for ; cont_row[j] != nil; cont_row[j] = obj1 {
							obj1 = cont_row[j].Next_content
							obj_to_char(cont_row[j], ch)
						}
						cont_row[j] = nil
					}
				}
				if cont_row[0] != nil {
					if int(temp.Type_flag) == ITEM_CONTAINER {
						temp = unequip_char(ch, locate-1)
						temp.Contains = nil
						for ; cont_row[0] != nil; cont_row[0] = obj1 {
							obj1 = cont_row[0].Next_content
							obj_to_obj(cont_row[0], temp)
						}
						equip_char(ch, temp, locate-1)
					} else {
						for ; cont_row[0] != nil; cont_row[0] = obj1 {
							obj1 = cont_row[0].Next_content
							obj_to_char(cont_row[0], ch)
						}
						cont_row[0] = nil
					}
				}
			} else {
				for j = int(MAX_BAG_ROWS - 1); j > -locate; j-- {
					if cont_row[j] != nil {
						for ; cont_row[j] != nil; cont_row[j] = obj1 {
							obj1 = cont_row[j].Next_content
							obj_to_char(cont_row[j], ch)
						}
						cont_row[j] = nil
					}
				}
				if j == -locate && cont_row[j] != nil {
					if int(temp.Type_flag) == ITEM_CONTAINER {
						obj_from_char(temp)
						temp.Contains = nil
						for ; cont_row[j] != nil; cont_row[j] = obj1 {
							obj1 = cont_row[j].Next_content
							obj_to_obj(cont_row[j], temp)
						}
						obj_to_char(temp, ch)
					} else {
						for ; cont_row[j] != nil; cont_row[j] = obj1 {
							obj1 = cont_row[j].Next_content
							obj_to_char(cont_row[j], ch)
						}
						cont_row[j] = nil
					}
				}
				if locate < 0 && locate >= -MAX_BAG_ROWS {
					obj_from_char(temp)
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
			}
		} else {
			get_line(fl, &line[0])
		}
	}
	mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s (level %d) has %d objects (max %d)."), GET_NAME(ch), GET_LEVEL(ch), num_objs, config_info.Csd.Max_obj_save)
	fl.Close()
	Crash_crashsave(ch)
	if orig_rent_code == RENT_RENTED || orig_rent_code == RENT_CRYO {
		return 0
	} else {
		return 1
	}
}
