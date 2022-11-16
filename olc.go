package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
	"unsafe"
)

const OLC_SET = 0
const OLC_SHOW = 1
const OLC_REPEAT = 2
const OLC_ROOM_TYPE = 3
const OLC_MOB_TYPE = 4
const OLC_OBJ_TYPE = 5
const OLC_COPY = 0
const OLC_NAME = 1
const OLC_DESC_TYPE = 2
const OLC_ALIASES = 3
const MAX_ROOM_NAME = 512
const MAX_MOB_NAME = 512
const MAX_OBJ_NAME = 512
const MAX_ROOM_DESC = 4096
const MAX_MOB_DESC = 512
const MAX_OBJ_DESC = 512
const OLC_USAGE = "Usage: olc { . | set | show | obj | mob | room} [args]\r\n"

var olc_ch *char_data

var olc_modes [8]*byte = [8]*byte{libc.CString("set"), libc.CString("show"), libc.CString("."), libc.CString("room"), libc.CString("mobile"), libc.CString("object"), libc.CString("assedit"), libc.CString("\n")}
var olc_commands [5]*byte = [5]*byte{libc.CString("copy"), libc.CString("name"), libc.CString("description"), libc.CString("aliases"), libc.CString("\n")}

func do_olc(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		olc_targ unsafe.Pointer = nil
		mode_arg [2048]byte
		arg      [2048]byte
		rnum     room_rnum
		vnum     room_vnum = room_vnum(-1)
		olc_mode int
	)
	if libc.StrCmp(GET_NAME(ch), libc.CString("Ras")) != 0 {
		send_to_char(ch, libc.CString("OLC is not yet complete.  Sorry.\r\n"))
		return
	}
	half_chop(argument, &mode_arg[0], argument)
	if (func() int {
		olc_mode = search_block(&mode_arg[0], &olc_modes[0], FALSE)
		return olc_mode
	}()) < 0 {
		send_to_char(ch, libc.CString("Invalid mode '%s'.\r\n%s"), &mode_arg[0], OLC_USAGE)
		return
	}
	switch olc_mode {
	case OLC_SET:
		fallthrough
	case OLC_SHOW:
		olc_set_show(ch, olc_mode, argument)
		return
	case OLC_REPEAT:
		if (func() int {
			olc_mode = ch.Player_specials.Last_olc_mode
			return olc_mode
		}()) == 0 || (func() unsafe.Pointer {
			olc_targ = ch.Player_specials.Last_olc_targ
			return olc_targ
		}()) == nil {
			send_to_char(ch, libc.CString("No last OLC operation!\r\n"))
			return
		}
	case OLC_ROOM_TYPE:
		if unicode.IsDigit(rune(*argument)) {
			argument = one_argument(argument, &arg[0])
			if is_number(&arg[0]) == 0 {
				send_to_char(ch, libc.CString("Invalid room vnum '%s'.\r\n"), &arg[0])
				return
			}
			vnum = room_vnum(libc.Atoi(libc.GoString(&arg[0])))
			if (func() room_rnum {
				rnum = real_room(vnum)
				return rnum
			}()) == room_rnum(-1) {
				send_to_char(ch, libc.CString("No such room!\r\n"))
				return
			}
		} else {
			rnum = ch.In_room
			vnum = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
			send_to_char(ch, libc.CString("(Using current room %d)\r\n"), vnum)
		}
		olc_targ = unsafe.Pointer(&(world[rnum]))
	case OLC_MOB_TYPE:
		argument = one_argument(argument, &arg[0])
		if is_number(&arg[0]) == 0 {
			send_to_char(ch, libc.CString("Invalid mob vnum '%s'.\r\n"), &arg[0])
			return
		}
		vnum = room_vnum(libc.Atoi(libc.GoString(&arg[0])))
		if (func() room_rnum {
			rnum = room_rnum(real_mobile(mob_vnum(vnum)))
			return rnum
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("No such mobile vnum.\r\n"))
		} else {
			olc_targ = unsafe.Pointer(&(mob_proto[rnum]))
		}
	case OLC_OBJ_TYPE:
		argument = one_argument(argument, &arg[0])
		if is_number(&arg[0]) == 0 {
			send_to_char(ch, libc.CString("Invalid obj vnum '%s'\r\n"), &arg[0])
			return
		}
		vnum = room_vnum(libc.Atoi(libc.GoString(&arg[0])))
		if (func() room_rnum {
			rnum = room_rnum(real_object(obj_vnum(vnum)))
			return rnum
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("No object with vnum %d.\r\n"), vnum)
		} else {
			olc_targ = unsafe.Pointer(&(obj_proto[rnum]))
		}
	default:
		send_to_char(ch, libc.CString("Usage: olc {.|set|show|obj|mob|room} [args]\r\n"))
		return
	}
	if olc_targ == nil {
		return
	}
	if can_modify(ch, int(vnum)) == 0 {
		send_to_char(ch, libc.CString("You can't modify that.\r\n"))
		return
	}
	ch.Player_specials.Last_olc_mode = olc_mode
	ch.Player_specials.Last_olc_targ = olc_targ
	olc_ch = ch
	olc_interpreter(olc_targ, olc_mode, argument)
}
func olc_interpreter(targ unsafe.Pointer, mode int, arg *byte) {
	var error int = 0
	_ = error
	var command int
	var command_string [2048]byte
	var olc_mob *char_data = nil
	var olc_room *room_data = nil
	var olc_obj *obj_data = nil
	half_chop(arg, &command_string[0], arg)
	if (func() int {
		command = search_block(&command_string[0], &olc_commands[0], FALSE)
		return command
	}()) < 0 {
		send_to_char(olc_ch, libc.CString("Invalid OLC command '%s'.\r\n"), &command_string[0])
		return
	}
	switch mode {
	case OLC_ROOM_TYPE:
		olc_room = (*room_data)(targ)
	case OLC_MOB_TYPE:
		olc_mob = (*char_data)(targ)
	case OLC_OBJ_TYPE:
		olc_obj = (*obj_data)(targ)
	default:
		basic_mud_log(libc.CString("SYSERR: Invalid OLC mode %d passed to interp."), mode)
		return
	}
	switch command {
	case OLC_COPY:
		switch mode {
		case OLC_ROOM_TYPE:
		case OLC_MOB_TYPE:
		case OLC_OBJ_TYPE:
		default:
			error = 1
		}
	case OLC_NAME:
		switch mode {
		case OLC_ROOM_TYPE:
			olc_string(&olc_room.Name, MAX_ROOM_NAME, arg)
		case OLC_MOB_TYPE:
			olc_string(&olc_mob.Short_descr, MAX_MOB_NAME, arg)
		case OLC_OBJ_TYPE:
			olc_string(&olc_obj.Short_description, MAX_OBJ_NAME, arg)
		default:
			error = 1
		}
	case OLC_DESC_TYPE:
		switch mode {
		case OLC_ROOM_TYPE:
			olc_string(&olc_room.Description, MAX_ROOM_DESC, arg)
		case OLC_MOB_TYPE:
			olc_string(&olc_mob.Long_descr, MAX_MOB_DESC, arg)
		case OLC_OBJ_TYPE:
			olc_string(&olc_obj.Description, MAX_OBJ_DESC, arg)
		default:
			error = 1
		}
	case OLC_ALIASES:
		switch mode {
		case OLC_ROOM_TYPE:
		case OLC_MOB_TYPE:
		case OLC_OBJ_TYPE:
		default:
			error = 1
		}
	}
}
func can_modify(ch *char_data, vnum int) int {
	return 1
}
func olc_string(string_ **byte, maxlen uint64, arg *byte) {
	skip_spaces(&arg)
	if *arg == 0 {
		send_to_char(olc_ch, libc.CString("Enter new string (max of %d characters); use '@' on a new line when done.\r\n"), int(maxlen))
		**string_ = '\x00'
		string_write(olc_ch.Desc, string_, maxlen, 0, nil)
	} else {
		if libc.StrLen(arg) > int(maxlen) {
			send_to_char(olc_ch, libc.CString("String too long (cannot be more than %d chars).\r\n"), int(maxlen))
		} else {
			if *string_ != nil {
				libc.Free(unsafe.Pointer(*string_))
			}
			*string_ = libc.StrDup(arg)
			send_to_char(olc_ch, libc.CString("%s"), config_info.Play.OK)
		}
	}
}
func olc_bitvector(bv *bitvector_t, names **byte, arg *byte) {
	var (
		newbv     bitvector_t
		flagnum   int
		doremove  int = 0
		this_name *byte
		buf       [64936]byte
	)
	skip_spaces(&arg)
	if *arg == 0 {
		send_to_char(olc_ch, libc.CString("Flag list or flag modifiers required.\r\n"))
		return
	}
	if *arg == '+' || *arg == '-' {
		newbv = *bv
	} else {
		newbv = 0
	}
	for *arg != 0 {
		arg = one_argument(arg, &buf[0])
		for this_name = &buf[0]; *this_name != 0; this_name = (*byte)(unsafe.Add(unsafe.Pointer(this_name), 1)) {
			CAP(this_name)
		}
		if buf[0] == '+' || buf[0] == '-' {
			this_name = &buf[1]
			if buf[0] == '-' {
				doremove = TRUE
			} else {
				doremove = FALSE
			}
		} else {
			this_name = &buf[0]
			doremove = FALSE
		}
		if (func() int {
			flagnum = search_block(this_name, names, TRUE)
			return flagnum
		}()) < 0 {
			send_to_char(olc_ch, libc.CString("Unknown flag: %s\r\n"), this_name)
		} else {
			if doremove != 0 {
				newbv &= ^(1 << flagnum)
			} else {
				newbv |= 1 << flagnum
			}
		}
	}
	*bv = newbv
	//sprintbit(bitvector_t(int32(newbv)), ([]*byte)(names), &buf[0], uint64(64936))
	send_to_char(olc_ch, libc.CString("Flags now set to: %s\r\n"), &buf[0])
}
func olc_set_show(ch *char_data, olc_mode int, arg *byte) {
}
