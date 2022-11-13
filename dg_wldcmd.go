package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const SCMD_WSEND = 0
const SCMD_WECHOAROUND = 1

var do_wasound func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wecho func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wsend func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wzoneecho func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wrecho func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wdoor func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wteleport func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wforce func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wpurge func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wload func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wdamage func(room *room_data, argument *byte, cmd int, subcmd int)
var do_wat func(room *room_data, argument *byte, cmd int, subcmd int)
var do_weffect func(room *room_data, argument *byte, cmd int, subcmd int)

type wld_command_info struct {
	Command         *byte
	Command_pointer func(room *room_data, argument *byte, cmd int, subcmd int)
	Subcmd          int
}

func wld_log(room *room_data, format *byte, _rest ...interface{}) {
	var (
		args   libc.ArgList
		output [64936]byte
	)
	stdio.Snprintf(&output[0], int(64936), "Room %d :: %s", room.Number, format)
	args.Start(format, _rest)
	script_vlog(&output[0], args)
	args.End()
}
func act_to_room(str *byte, room *room_data) {
	if room.People == nil {
		return
	}
	act(str, FALSE, room.People, nil, nil, TO_ROOM)
	act(str, FALSE, room.People, nil, nil, TO_CHAR)
}
func do_weffect(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		arg2   [2048]byte
		num    int = 0
		target room_rnum
		nr     room_rnum
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		wld_log(room, libc.CString("weffect called without type argument"))
		return
	}
	if arg2[0] == 0 {
		wld_log(room, libc.CString("weffect called without setting argument"))
		return
	}
	num = libc.Atoi(libc.GoString(&arg2[0]))
	nr = room_rnum(num)
	target = real_room(room_vnum(nr))
	if libc.StrCaseCmp(&arg[0], libc.CString("gravity")) == 0 {
		if num < 0 || num > 10000 {
			wld_log(room, libc.CString("weffect setting out of bounds, 0 - 10000 only."))
			return
		} else {
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(room.Number))))).Gravity = num
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("light")) == 0 {
		if target == room_rnum(-1) {
			wld_log(room, libc.CString("weffect target is NOWHERE."))
			return
		} else {
			if !ROOM_FLAGGED(target, ROOM_INDOORS) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target)))).Room_flags[int(ROOM_INDOORS/32)] |= bitvector_t(int32(1 << (int(ROOM_INDOORS % 32))))
			} else {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target)))).Room_flags[int(ROOM_INDOORS/32)] &= bitvector_t(int32(^(1 << (int(ROOM_INDOORS % 32)))))
			}
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("lava")) == 0 {
		if target == room_rnum(-1) {
			wld_log(room, libc.CString("weffect target is NOWHERE."))
			return
		} else {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target)))).Geffect != 0 {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target)))).Geffect = 5
			} else {
				wld_log(room, libc.CString("weffect target already has lava."))
				return
			}
		}
	}
}
func do_wasound(room *room_data, argument *byte, cmd int, subcmd int) {
	var door int
	skip_spaces(&argument)
	if *argument == 0 {
		wld_log(room, libc.CString("wasound called with no argument"))
		return
	}
	for door = 0; door < NUM_OF_DIRS; door++ {
		var newexit *room_direction_data
		if (func() *room_direction_data {
			newexit = room.Dir_option[door]
			return newexit
		}()) != nil && newexit.To_room != room_rnum(-1) && room != (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newexit.To_room))) {
			act_to_room(argument, (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newexit.To_room))))
		}
	}
}
func do_wecho(room *room_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	if *argument == 0 {
		wld_log(room, libc.CString("wecho called with no args"))
	} else {
		act_to_room(argument, room)
	}
}
func do_wsend(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		buf [2048]byte
		msg *byte
		ch  *char_data
	)
	msg = any_one_arg(argument, &buf[0])
	if buf[0] == 0 {
		wld_log(room, libc.CString("wsend called with no args"))
		return
	}
	skip_spaces(&msg)
	if *msg == 0 {
		wld_log(room, libc.CString("wsend called without a message"))
		return
	}
	if (func() *char_data {
		ch = get_char_by_room(room, &buf[0])
		return ch
	}()) != nil {
		if subcmd == SCMD_WSEND {
			sub_write(msg, ch, TRUE, TO_CHAR)
		} else if subcmd == SCMD_WECHOAROUND {
			sub_write(msg, ch, TRUE, TO_ROOM)
		}
	} else {
		wld_log(room, libc.CString("no target found for wsend"))
	}
}
func do_wzoneecho(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		zone     zone_rnum
		room_num [2048]byte
		buf      [2048]byte
		msg      *byte
	)
	msg = any_one_arg(argument, &room_num[0])
	skip_spaces(&msg)
	if room_num[0] == 0 || *msg == 0 {
		wld_log(room, libc.CString("wzoneecho called with too few args"))
	} else if (func() zone_rnum {
		zone = real_zone_by_thing(room_vnum(libc.Atoi(libc.GoString(&room_num[0]))))
		return zone
	}()) == zone_rnum(-1) {
		wld_log(room, libc.CString("wzoneecho called for nonexistant zone"))
	} else {
		stdio.Sprintf(&buf[0], "%s\r\n", msg)
		send_to_zone(&buf[0], zone)
	}
}
func do_wrecho(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		start  [2048]byte
		finish [2048]byte
		msg    *byte
	)
	msg = two_arguments(argument, &start[0], &finish[0])
	skip_spaces(&msg)
	if *msg == 0 || start[0] == 0 || finish[0] == 0 || is_number(&start[0]) == 0 || is_number(&finish[0]) == 0 {
		wld_log(room, libc.CString("wrecho: too few args"))
	} else {
		send_to_range(room_vnum(libc.Atoi(libc.GoString(&start[0]))), room_vnum(libc.Atoi(libc.GoString(&finish[0]))), libc.CString("%s\r\n"), msg)
	}
}
func do_wdoor(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		target     [2048]byte
		direction  [2048]byte
		field      [2048]byte
		value      *byte
		rm         *room_data
		newexit    *room_direction_data
		dir        int
		fd         int
		to_room    int
		door_field [7]*byte = [7]*byte{libc.CString("purge"), libc.CString("description"), libc.CString("flags"), libc.CString("key"), libc.CString("name"), libc.CString("room"), libc.CString("\n")}
	)
	argument = two_arguments(argument, &target[0], &direction[0])
	value = one_argument(argument, &field[0])
	skip_spaces(&value)
	if target[0] == 0 || direction[0] == 0 || field[0] == 0 {
		wld_log(room, libc.CString("wdoor called with too few args"))
		return
	}
	if (func() *room_data {
		rm = get_room(&target[0])
		return rm
	}()) == nil {
		wld_log(room, libc.CString("wdoor: invalid target"))
		return
	}
	if libc.Atoi(libc.GoString(&direction[0])) >= 0 && libc.Atoi(libc.GoString(&direction[0])) <= 11 {
		dir = libc.Atoi(libc.GoString(&direction[0]))
	} else if libc.Atoi(libc.GoString(&direction[0])) < 0 && libc.Atoi(libc.GoString(&direction[0])) > 11 {
		wld_log(room, libc.CString("wdoor: invalid direction"))
		return
	}
	if (func() int {
		fd = search_block(&field[0], &door_field[0], FALSE)
		return fd
	}()) == -1 {
		wld_log(room, libc.CString("wdoor: invalid field"))
		return
	}
	newexit = rm.Dir_option[dir]
	if fd == 0 {
		if newexit != nil {
			if newexit.General_description != nil {
				libc.Free(unsafe.Pointer(newexit.General_description))
			}
			if newexit.Keyword != nil {
				libc.Free(unsafe.Pointer(newexit.Keyword))
			}
			libc.Free(unsafe.Pointer(newexit))
			rm.Dir_option[dir] = nil
		}
	} else {
		if newexit == nil {
			newexit = new(room_direction_data)
			rm.Dir_option[dir] = newexit
		}
		switch fd {
		case 1:
			if newexit.General_description != nil {
				libc.Free(unsafe.Pointer(newexit.General_description))
			}
			newexit.General_description = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(value)+3)[0]))
			libc.StrCpy(newexit.General_description, value)
			libc.StrCat(newexit.General_description, libc.CString("\r\n"))
		case 2:
			newexit.Exit_info = bitvector_t(int16(uint16(asciiflag_conv(value))))
		case 3:
			newexit.Key = obj_vnum(libc.Atoi(libc.GoString(value)))
		case 4:
			if newexit.Keyword != nil {
				libc.Free(unsafe.Pointer(newexit.Keyword))
			}
			newexit.Keyword = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(value)+1)[0]))
			libc.StrCpy(newexit.Keyword, value)
		case 5:
			if (func() int {
				to_room = int(real_room(room_vnum(libc.Atoi(libc.GoString(value)))))
				return to_room
			}()) != int(-1) {
				newexit.To_room = room_rnum(to_room)
			} else {
				wld_log(room, libc.CString("wdoor: invalid door target"))
			}
		}
	}
}
func do_wteleport(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		ch      *char_data
		next_ch *char_data
		target  room_rnum
		nr      room_rnum
		arg1    [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 {
		wld_log(room, libc.CString("wteleport called with too few args"))
		return
	}
	nr = room_rnum(libc.Atoi(libc.GoString(&arg2[0])))
	target = real_room(room_vnum(nr))
	if target == room_rnum(-1) {
		wld_log(room, libc.CString("wteleport target is an invalid room"))
	} else if libc.StrCaseCmp(&arg1[0], libc.CString("all")) == 0 {
		if nr == room_rnum(room.Number) {
			wld_log(room, libc.CString("wteleport all target is itself"))
			return
		}
		for ch = room.People; ch != nil; ch = next_ch {
			next_ch = ch.Next_in_room
			if valid_dg_target(ch, 1<<0) == 0 {
				continue
			}
			char_from_room(ch)
			char_to_room(ch, target)
			enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
		}
	} else {
		if (func() *char_data {
			ch = get_char_by_room(room, &arg1[0])
			return ch
		}()) != nil {
			if valid_dg_target(ch, 1<<0) != 0 {
				char_from_room(ch)
				char_to_room(ch, target)
				enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
			}
		} else {
			wld_log(room, libc.CString("wteleport: no target found"))
		}
	}
}
func do_wforce(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		ch      *char_data
		next_ch *char_data
		arg1    [2048]byte
		line    *byte
	)
	line = one_argument(argument, &arg1[0])
	if arg1[0] == 0 || *line == 0 {
		wld_log(room, libc.CString("wforce called with too few args"))
		return
	}
	if libc.StrCaseCmp(&arg1[0], libc.CString("all")) == 0 {
		for ch = room.People; ch != nil; ch = next_ch {
			next_ch = ch.Next_in_room
			if valid_dg_target(ch, 0) != 0 {
				command_interpreter(ch, line)
			}
		}
	} else {
		if (func() *char_data {
			ch = get_char_by_room(room, &arg1[0])
			return ch
		}()) != nil {
			if valid_dg_target(ch, 0) != 0 {
				command_interpreter(ch, line)
			}
		} else {
			wld_log(room, libc.CString("wforce: no target found"))
		}
	}
}
func do_wpurge(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		ch       *char_data
		next_ch  *char_data
		obj      *obj_data
		next_obj *obj_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		for ch = room.People; ch != nil; ch = next_ch {
			next_ch = ch.Next_in_room
			if IS_NPC(ch) {
				extract_char(ch)
			}
		}
		for obj = room.Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			extract_obj(obj)
		}
		return
	}
	if arg[0] == UID_CHAR {
		ch = get_char(&arg[0])
	} else {
		ch = get_char_in_room(room, &arg[0])
	}
	if ch == nil {
		if arg[0] == UID_CHAR {
			obj = get_obj(&arg[0])
		} else {
			obj = get_obj_in_room(room, &arg[0])
		}
		if obj != nil {
			extract_obj(obj)
		} else {
			wld_log(room, libc.CString("wpurge: bad argument"))
		}
		return
	}
	if !IS_NPC(ch) {
		wld_log(room, libc.CString("wpurge: purging a PC"))
		return
	}
	extract_char(ch)
}
func do_wload(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1   [2048]byte
		arg2   [2048]byte
		number int = 0
		mob    *char_data
		object *obj_data
		target *byte
		tch    *char_data
		cnt    *obj_data
		pos    int
	)
	target = two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 || is_number(&arg2[0]) == 0 || (func() int {
		number = libc.Atoi(libc.GoString(&arg2[0]))
		return number
	}()) < 0 {
		wld_log(room, libc.CString("wload: bad syntax"))
		return
	}
	if is_abbrev(&arg1[0], libc.CString("mob")) != 0 {
		var rnum room_rnum
		if target == nil || *target == 0 {
			rnum = real_room(room.Number)
		} else {
			if !unicode.IsDigit(rune(*target)) || (func() room_rnum {
				rnum = real_room(room_vnum(libc.Atoi(libc.GoString(target))))
				return rnum
			}()) == room_rnum(-1) {
				wld_log(room, libc.CString("wload: room target vnum doesn't exist (loading mob vnum %d to room %s)"), number, target)
				return
			}
		}
		if (func() *char_data {
			mob = read_mobile(mob_vnum(number), VIRTUAL)
			return mob
		}()) == nil {
			wld_log(room, libc.CString("mload: bad mob vnum"))
			return
		}
		char_to_room(mob, rnum)
		if room.Script != nil {
			var buf [2048]byte
			stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, mob.Id)
			add_var(&room.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
		}
		load_mtrigger(mob)
	} else if is_abbrev(&arg1[0], libc.CString("obj")) != 0 {
		if (func() *obj_data {
			object = read_object(obj_vnum(number), VIRTUAL)
			return object
		}()) == nil {
			wld_log(room, libc.CString("wload: bad object vnum"))
			return
		}
		if target == nil || *target == 0 {
			add_unique_id(object)
			obj_to_room(object, real_room(room.Number))
			if room.Script != nil {
				var buf [2048]byte
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, object.Id)
				add_var(&room.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
			}
			load_otrigger(object)
			return
		}
		two_arguments(target, &arg1[0], &arg2[0])
		tch = get_char_in_room(room, &arg1[0])
		if tch != nil {
			if arg2 != nil && arg2[0] != 0 && (func() int {
				pos = find_eq_pos_script(&arg2[0])
				return pos
			}()) >= 0 && (tch.Equipment[pos]) == nil && can_wear_on_pos(object, pos) != 0 {
				equip_char(tch, object, pos)
				load_otrigger(object)
				return
			}
			obj_to_char(object, tch)
			load_otrigger(object)
			return
		}
		cnt = get_obj_in_room(room, &arg1[0])
		if cnt != nil && int(cnt.Type_flag) == ITEM_CONTAINER {
			obj_to_obj(object, cnt)
			load_otrigger(object)
			return
		}
		add_unique_id(object)
		obj_to_room(object, real_room(room.Number))
		load_otrigger(object)
		return
	} else {
		wld_log(room, libc.CString("wload: bad type"))
	}
}
func do_wdamage(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		name   [2048]byte
		amount [2048]byte
		dam    int = 0
		ch     *char_data
	)
	two_arguments(argument, &name[0], &amount[0])
	if name[0] == 0 || amount[0] == 0 {
		wld_log(room, libc.CString("wdamage: bad syntax"))
		return
	}
	dam = libc.Atoi(libc.GoString(&amount[0]))
	ch = get_char_by_room(room, &name[0])
	if ch == nil {
		wld_log(room, libc.CString("wdamage: target not found"))
		return
	}
	script_damage(ch, dam)
}
func do_wat(room *room_data, argument *byte, cmd int, subcmd int) {
	var (
		loc     room_rnum = room_rnum(-1)
		ch      *char_data
		arg     [2048]byte
		command *byte
	)
	command = any_one_arg(argument, &arg[0])
	if arg[0] == 0 {
		wld_log(room, libc.CString("wat called with no args"))
		return
	}
	skip_spaces(&command)
	if *command == 0 {
		wld_log(room, libc.CString("wat called without a command"))
		return
	}
	if unicode.IsDigit(rune(arg[0])) {
		loc = real_room(room_vnum(libc.Atoi(libc.GoString(&arg[0]))))
	} else if (func() *char_data {
		ch = get_char_by_room(room, &arg[0])
		return ch
	}()) != nil {
		loc = ch.In_room
	}
	if loc == room_rnum(-1) {
		wld_log(room, libc.CString("wat: location not found (%s)"), &arg[0])
		return
	}
	wld_command_interpreter((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(loc))), command)
}

var wld_cmd_info [16]wld_command_info = [16]wld_command_info{{Command: libc.CString("RESERVED"), Command_pointer: nil, Subcmd: 0}, {Command: libc.CString("wasound "), Command_pointer: do_wasound, Subcmd: 0}, {Command: libc.CString("wdoor "), Command_pointer: do_wdoor, Subcmd: 0}, {Command: libc.CString("wecho "), Command_pointer: do_wecho, Subcmd: 0}, {Command: libc.CString("wechoaround "), Command_pointer: do_wsend, Subcmd: SCMD_WECHOAROUND}, {Command: libc.CString("wforce "), Command_pointer: do_wforce, Subcmd: 0}, {Command: libc.CString("wload "), Command_pointer: do_wload, Subcmd: 0}, {Command: libc.CString("wpurge "), Command_pointer: do_wpurge, Subcmd: 0}, {Command: libc.CString("wrecho "), Command_pointer: do_wrecho, Subcmd: 0}, {Command: libc.CString("wsend "), Command_pointer: do_wsend, Subcmd: SCMD_WSEND}, {Command: libc.CString("wteleport "), Command_pointer: do_wteleport, Subcmd: 0}, {Command: libc.CString("wzoneecho "), Command_pointer: do_wzoneecho, Subcmd: 0}, {Command: libc.CString("wdamage "), Command_pointer: do_wdamage, Subcmd: 0}, {Command: libc.CString("wat "), Command_pointer: do_wat, Subcmd: 0}, {Command: libc.CString("weffect "), Command_pointer: do_weffect, Subcmd: 0}, {Command: libc.CString("\n"), Command_pointer: nil, Subcmd: 0}}

func wld_command_interpreter(room *room_data, argument *byte) {
	var (
		cmd    int
		length int
		line   *byte
		arg    [2048]byte
	)
	skip_spaces(&argument)
	if *argument == 0 {
		return
	}
	line = any_one_arg(argument, &arg[0])
	for func() int {
		length = libc.StrLen(&arg[0])
		return func() int {
			cmd = 0
			return cmd
		}()
	}(); *wld_cmd_info[cmd].Command != '\n'; cmd++ {
		if libc.StrNCmp(wld_cmd_info[cmd].Command, &arg[0], length) == 0 {
			break
		}
	}
	if *wld_cmd_info[cmd].Command == '\n' {
		wld_log(room, libc.CString("Unknown world cmd: '%s'"), argument)
	} else {
		(wld_cmd_info[cmd].Command_pointer)(room, line, cmd, wld_cmd_info[cmd].Subcmd)
	}
}
