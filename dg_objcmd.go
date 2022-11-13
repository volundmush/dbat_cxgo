package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const SCMD_OSEND = 0
const SCMD_OECHOAROUND = 1

var do_oecho func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_oforce func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_ozoneecho func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_osend func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_orecho func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_otimer func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_otransform func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_opurge func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_dupe func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_oteleport func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_dgoload func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_odamage func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_oasound func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_ogoto func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_odoor func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_osetval func(obj *obj_data, argument *byte, cmd int, subcmd int)
var do_oat func(obj *obj_data, argument *byte, cmd int, subcmd int)

type obj_command_info struct {
	Command         *byte
	Command_pointer func(obj *obj_data, argument *byte, cmd int, subcmd int)
	Subcmd          int
}

func obj_log(obj *obj_data, format *byte, _rest ...interface{}) {
	var (
		args   libc.ArgList
		output [64936]byte
	)
	stdio.Snprintf(&output[0], int(64936), "Obj (%s, VNum %d):: %s", obj.Short_description, GET_OBJ_VNUM(obj), format)
	args.Start(format, _rest)
	script_vlog(&output[0], args)
	args.End()
}
func obj_room(obj *obj_data) room_rnum {
	if obj.In_room != room_rnum(-1) {
		return obj.In_room
	} else if obj.Carried_by != nil {
		return obj.Carried_by.In_room
	} else if obj.Worn_by != nil {
		return obj.Worn_by.In_room
	} else if obj.In_obj != nil {
		return obj_room(obj.In_obj)
	} else {
		return -1
	}
}
func find_obj_target_room(obj *obj_data, rawroomstr *byte) room_rnum {
	var (
		tmp        int
		location   room_rnum
		target_mob *char_data
		target_obj *obj_data
		roomstr    [2048]byte
	)
	one_argument(rawroomstr, &roomstr[0])
	if roomstr[0] == 0 {
		return -1
	}
	if unicode.IsDigit(rune(roomstr[0])) && libc.StrChr(&roomstr[0], '.') == nil {
		tmp = libc.Atoi(libc.GoString(&roomstr[0]))
		if (func() room_rnum {
			location = real_room(room_vnum(tmp))
			return location
		}()) == room_rnum(-1) {
			return -1
		}
	} else if (func() *char_data {
		target_mob = get_char_by_obj(obj, &roomstr[0])
		return target_mob
	}()) != nil {
		location = target_mob.In_room
	} else if (func() *obj_data {
		target_obj = get_obj_by_obj(obj, &roomstr[0])
		return target_obj
	}()) != nil {
		if target_obj.In_room != room_rnum(-1) {
			location = target_obj.In_room
		} else {
			return -1
		}
	} else {
		return -1
	}
	if ROOM_FLAGGED(location, ROOM_GODROOM) || ROOM_FLAGGED(location, ROOM_PRIVATE) {
		return -1
	}
	return location
}
func do_oecho(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var room int
	skip_spaces(&argument)
	if *argument == 0 {
		obj_log(obj, libc.CString("oecho called with no args"))
	} else if (func() int {
		room = int(obj_room(obj))
		return room
	}()) != int(-1) {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).People != nil {
			sub_write(argument, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).People, TRUE, TO_ROOM)
			sub_write(argument, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).People, TRUE, TO_CHAR)
		}
	} else {
		obj_log(obj, libc.CString("oecho called by object in NOWHERE"))
	}
}
func do_oforce(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		ch      *char_data
		next_ch *char_data
		room    int
		arg1    [2048]byte
		line    *byte
	)
	line = one_argument(argument, &arg1[0])
	if arg1[0] == 0 || *line == 0 {
		obj_log(obj, libc.CString("oforce called with too few args"))
		return
	}
	if libc.StrCaseCmp(&arg1[0], libc.CString("all")) == 0 {
		if (func() int {
			room = int(obj_room(obj))
			return room
		}()) == int(-1) {
			obj_log(obj, libc.CString("oforce called by object in NOWHERE"))
		} else {
			for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).People; ch != nil; ch = next_ch {
				next_ch = ch.Next_in_room
				if valid_dg_target(ch, 0) != 0 {
					command_interpreter(ch, line)
				}
			}
		}
	} else {
		if (func() *char_data {
			ch = get_char_by_obj(obj, &arg1[0])
			return ch
		}()) != nil {
			if valid_dg_target(ch, 0) != 0 {
				command_interpreter(ch, line)
			}
		} else {
			obj_log(obj, libc.CString("oforce: no target found"))
		}
	}
}
func do_ozoneecho(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		zone        int
		room_number [2048]byte
		buf         [2048]byte
		msg         *byte
	)
	msg = any_one_arg(argument, &room_number[0])
	skip_spaces(&msg)
	if room_number[0] == 0 || *msg == 0 {
		obj_log(obj, libc.CString("ozoneecho called with too few args"))
	} else if (func() int {
		zone = int(real_zone_by_thing(room_vnum(libc.Atoi(libc.GoString(&room_number[0])))))
		return zone
	}()) == int(-1) {
		obj_log(obj, libc.CString("ozoneecho called for nonexistant zone"))
	} else {
		stdio.Sprintf(&buf[0], "%s\r\n", msg)
		send_to_zone(&buf[0], zone_rnum(zone))
	}
}
func do_osend(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		buf [2048]byte
		msg *byte
		ch  *char_data
	)
	msg = any_one_arg(argument, &buf[0])
	if buf[0] == 0 {
		obj_log(obj, libc.CString("osend called with no args"))
		return
	}
	skip_spaces(&msg)
	if *msg == 0 {
		obj_log(obj, libc.CString("osend called without a message"))
		return
	}
	if (func() *char_data {
		ch = get_char_by_obj(obj, &buf[0])
		return ch
	}()) != nil {
		if subcmd == SCMD_OSEND {
			sub_write(msg, ch, TRUE, TO_CHAR)
		} else if subcmd == SCMD_OECHOAROUND {
			var buf [64936]byte
			stdio.Sprintf(&buf[0], libc.GoString(msg))
			search_replace(&buf[0], GET_NAME(ch), libc.CString("$n"))
			act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
		}
	} else {
		obj_log(obj, libc.CString("no target found for osend"))
	}
}
func do_orecho(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		start  [2048]byte
		finish [2048]byte
		msg    *byte
	)
	msg = two_arguments(argument, &start[0], &finish[0])
	skip_spaces(&msg)
	if *msg == 0 || start[0] == 0 || finish[0] == 0 || is_number(&start[0]) == 0 || is_number(&finish[0]) == 0 {
		obj_log(obj, libc.CString("orecho: too few args"))
	} else {
		send_to_range(room_vnum(libc.Atoi(libc.GoString(&start[0]))), room_vnum(libc.Atoi(libc.GoString(&finish[0]))), libc.CString("%s\r\n"), msg)
	}
}
func do_otimer(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		obj_log(obj, libc.CString("otimer: missing argument"))
	} else if !unicode.IsDigit(rune(arg[0])) {
		obj_log(obj, libc.CString("otimer: bad argument"))
	} else {
		obj.Timer = libc.Atoi(libc.GoString(&arg[0]))
	}
}
func do_otransform(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		o      *obj_data
		tmpobj obj_data
		wearer *char_data = nil
		pos    int        = 0
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		obj_log(obj, libc.CString("otransform: missing argument"))
	} else if !unicode.IsDigit(rune(arg[0])) {
		obj_log(obj, libc.CString("otransform: bad argument"))
	} else {
		o = read_object(obj_vnum(libc.Atoi(libc.GoString(&arg[0]))), VIRTUAL)
		if o == nil {
			obj_log(obj, libc.CString("otransform: bad object vnum"))
			return
		}
		if obj.Worn_by != nil {
			pos = int(obj.Worn_on)
			wearer = obj.Worn_by
			unequip_char(obj.Worn_by, pos)
		}
		libc.MemCpy(unsafe.Pointer(&tmpobj), unsafe.Pointer(o), int(unsafe.Sizeof(obj_data{})))
		tmpobj.In_room = obj.In_room
		tmpobj.Carried_by = obj.Carried_by
		tmpobj.Worn_by = obj.Worn_by
		tmpobj.Worn_on = obj.Worn_on
		tmpobj.In_obj = obj.In_obj
		tmpobj.Contains = obj.Contains
		tmpobj.Id = obj.Id
		tmpobj.Proto_script = obj.Proto_script
		tmpobj.Script = obj.Script
		tmpobj.Next_content = obj.Next_content
		tmpobj.Next = obj.Next
		libc.MemCpy(unsafe.Pointer(obj), unsafe.Pointer(&tmpobj), int(unsafe.Sizeof(obj_data{})))
		if wearer != nil {
			equip_char(wearer, obj, pos)
		}
		extract_obj(o)
	}
}
func do_dupe(obj *obj_data, argument *byte, cmd int, subcmd int) {
	obj.Extra_flags[int(ITEM_DUPLICATE/32)] |= bitvector_t(int32(1 << (int(ITEM_DUPLICATE % 32))))
}
func do_opurge(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		ch       *char_data
		next_ch  *char_data
		o        *obj_data
		next_obj *obj_data
		rm       int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		if (func() int {
			rm = int(obj_room(obj))
			return rm
		}()) != int(-1) {
			for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).People; ch != nil; ch = next_ch {
				next_ch = ch.Next_in_room
				if IS_NPC(ch) {
					extract_char(ch)
				}
			}
			for o = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).Contents; o != nil; o = next_obj {
				next_obj = o.Next_content
				if o != obj {
					extract_obj(o)
				}
			}
		}
		return
	}
	ch = get_char_by_obj(obj, &arg[0])
	if ch == nil {
		o = get_obj_by_obj(obj, &arg[0])
		if o != nil {
			if o == obj {
				dg_owner_purged = 1
			}
			extract_obj(o)
		} else {
			obj_log(obj, libc.CString("opurge: bad argument"))
		}
		return
	}
	if !IS_NPC(ch) {
		obj_log(obj, libc.CString("opurge: purging a PC"))
		return
	}
	extract_char(ch)
}
func do_ogoto(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		target room_rnum
		arg1   [2048]byte
	)
	one_argument(argument, &arg1[0])
	if arg1[0] == 0 {
		obj_log(obj, libc.CString("ogoto called with too few args"))
		return
	}
	target = find_obj_target_room(obj, &arg1[0])
	if target == room_rnum(-1) {
		obj_log(obj, libc.CString("ogoto target is an invalid room"))
	} else if obj.In_room == room_rnum(-1) {
		obj_log(obj, libc.CString("ogoto tried to leave nowhere"))
	} else {
		obj_from_room(obj)
		obj_to_room(obj, target)
	}
}
func do_oteleport(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		ch      *char_data
		next_ch *char_data
		target  room_rnum
		rm      room_rnum
		arg1    [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 {
		obj_log(obj, libc.CString("oteleport called with too few args"))
		return
	}
	target = find_obj_target_room(obj, &arg2[0])
	if target == room_rnum(-1) {
		obj_log(obj, libc.CString("oteleport target is an invalid room"))
	} else if libc.StrCaseCmp(&arg1[0], libc.CString("all")) == 0 {
		rm = obj_room(obj)
		if target == rm {
			obj_log(obj, libc.CString("oteleport target is itself"))
		}
		for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).People; ch != nil; ch = next_ch {
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
			ch = get_char_by_obj(obj, &arg1[0])
			return ch
		}()) != nil {
			if valid_dg_target(ch, 1<<0) != 0 {
				char_from_room(ch)
				char_to_room(ch, target)
				enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
			}
		} else {
			obj_log(obj, libc.CString("oteleport: no target found"))
		}
	}
}
func do_dgoload(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1   [2048]byte
		arg2   [2048]byte
		number int = 0
		room   int
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
		obj_log(obj, libc.CString("oload: bad syntax"))
		return
	}
	if (func() int {
		room = int(obj_room(obj))
		return room
	}()) == int(-1) {
		obj_log(obj, libc.CString("oload: object in NOWHERE trying to load"))
		return
	}
	if is_abbrev(&arg1[0], libc.CString("mob")) != 0 {
		var rnum room_rnum
		if target == nil || *target == 0 {
			rnum = room_rnum(room)
		} else {
			if !unicode.IsDigit(rune(*target)) || (func() room_rnum {
				rnum = real_room(room_vnum(libc.Atoi(libc.GoString(target))))
				return rnum
			}()) == room_rnum(-1) {
				obj_log(obj, libc.CString("oload: room target vnum doesn't exist (loading mob vnum %d to room %s)"), number, target)
				return
			}
		}
		if (func() *char_data {
			mob = read_mobile(mob_vnum(number), VIRTUAL)
			return mob
		}()) == nil {
			obj_log(obj, libc.CString("oload: bad mob vnum"))
			return
		}
		char_to_room(mob, rnum)
		if obj.Script != nil {
			var buf [2048]byte
			stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, mob.Id)
			add_var(&obj.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
		}
		load_mtrigger(mob)
	} else if is_abbrev(&arg1[0], libc.CString("obj")) != 0 {
		if (func() *obj_data {
			object = read_object(obj_vnum(number), VIRTUAL)
			return object
		}()) == nil {
			obj_log(obj, libc.CString("oload: bad object vnum"))
			return
		}
		if obj.Script != nil {
			var buf [2048]byte
			stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, object.Id)
			add_var(&obj.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
		}
		if target == nil || *target == 0 {
			add_unique_id(object)
			obj_to_room(object, room_rnum(room))
			load_otrigger(object)
			return
		}
		two_arguments(target, &arg1[0], &arg2[0])
		tch = get_char_near_obj(obj, &arg1[0])
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
		cnt = get_obj_near_obj(obj, &arg1[0])
		if cnt != nil && int(cnt.Type_flag) == ITEM_CONTAINER {
			obj_to_obj(object, cnt)
			load_otrigger(object)
			return
		}
		add_unique_id(object)
		obj_to_room(object, room_rnum(room))
		load_otrigger(object)
		return
	} else {
		obj_log(obj, libc.CString("oload: bad type"))
	}
}
func do_odamage(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		name   [2048]byte
		amount [2048]byte
		dam    int = 0
		ch     *char_data
	)
	two_arguments(argument, &name[0], &amount[0])
	if name[0] == 0 || amount[0] == 0 {
		obj_log(obj, libc.CString("odamage: bad syntax"))
		return
	}
	dam = libc.Atoi(libc.GoString(&amount[0]))
	ch = get_char_by_obj(obj, &name[0])
	if ch == nil {
		obj_log(obj, libc.CString("odamage: target not found"))
		return
	}
	script_damage(ch, dam)
}
func do_oasound(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		room room_rnum
		door int
	)
	skip_spaces(&argument)
	if *argument == 0 {
		obj_log(obj, libc.CString("oasound called with no args"))
		return
	}
	if (func() room_rnum {
		room = obj_room(obj)
		return room
	}()) == room_rnum(-1) {
		obj_log(obj, libc.CString("oecho called by object in NOWHERE"))
		return
	}
	for door = 0; door < NUM_OF_DIRS; door++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door] != nil && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door]).To_room != room && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door]).To_room)))).People != nil {
			sub_write(argument, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door]).To_room)))).People, TRUE, TO_ROOM)
			sub_write(argument, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).Dir_option[door]).To_room)))).People, TRUE, TO_CHAR)
		}
	}
}
func do_odoor(obj *obj_data, argument *byte, cmd int, subcmd int) {
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
		obj_log(obj, libc.CString("odoor called with too few args"))
		return
	}
	if (func() *room_data {
		rm = get_room(&target[0])
		return rm
	}()) == nil {
		obj_log(obj, libc.CString("odoor: invalid target"))
		return
	}
	if (func() int {
		dir = search_block(&direction[0], &dirs[0], FALSE)
		return dir
	}()) == -1 {
		obj_log(obj, libc.CString("odoor: invalid direction"))
		return
	}
	if (func() int {
		fd = search_block(&field[0], &door_field[0], FALSE)
		return fd
	}()) == -1 {
		obj_log(obj, libc.CString("odoor: invalid field"))
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
				obj_log(obj, libc.CString("odoor: invalid door target"))
			}
		}
	}
}
func do_osetval(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1      [2048]byte
		arg2      [2048]byte
		position  int
		new_value int
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if arg1 == nil || arg1[0] == 0 || arg2 == nil || arg2[0] == 0 || is_number(&arg1[0]) == 0 || is_number(&arg2[0]) == 0 {
		obj_log(obj, libc.CString("osetval: bad syntax"))
		return
	}
	position = libc.Atoi(libc.GoString(&arg1[0]))
	new_value = libc.Atoi(libc.GoString(&arg2[0]))
	if position >= 0 && position < NUM_OBJ_VAL_POSITIONS {
		obj.Value[position] = new_value
	} else {
		obj_log(obj, libc.CString("osetval: position out of bounds!"))
	}
}
func do_oat(obj *obj_data, argument *byte, cmd int, subcmd int) {
	var (
		loc     room_rnum = room_rnum(-1)
		ch      *char_data
		object  *obj_data
		arg     [2048]byte
		command *byte
	)
	command = any_one_arg(argument, &arg[0])
	if arg[0] == 0 {
		obj_log(obj, libc.CString("oat called with no args"))
		return
	}
	skip_spaces(&command)
	if *command == 0 {
		obj_log(obj, libc.CString("oat called without a command"))
		return
	}
	if unicode.IsDigit(rune(arg[0])) {
		loc = real_room(room_vnum(libc.Atoi(libc.GoString(&arg[0]))))
	} else if (func() *char_data {
		ch = get_char_by_obj(obj, &arg[0])
		return ch
	}()) != nil {
		loc = ch.In_room
	}
	if loc == room_rnum(-1) {
		obj_log(obj, libc.CString("oat: location not found (%s)"), &arg[0])
		return
	}
	if (func() *obj_data {
		object = read_object(GET_OBJ_VNUM(obj), VIRTUAL)
		return object
	}()) == nil {
		return
	}
	add_unique_id(object)
	obj_to_room(object, loc)
	obj_command_interpreter(object, command)
	if object.In_room == loc {
		extract_obj(object)
	}
}

var obj_cmd_info [20]obj_command_info = [20]obj_command_info{{Command: libc.CString("RESERVED"), Command_pointer: nil, Subcmd: 0}, {Command: libc.CString("oasound "), Command_pointer: do_oasound, Subcmd: 0}, {Command: libc.CString("oat "), Command_pointer: do_oat, Subcmd: 0}, {Command: libc.CString("odoor "), Command_pointer: do_odoor, Subcmd: 0}, {Command: libc.CString("odupe "), Command_pointer: do_dupe, Subcmd: 0}, {Command: libc.CString("odamage "), Command_pointer: do_odamage, Subcmd: 0}, {Command: libc.CString("oecho "), Command_pointer: do_oecho, Subcmd: 0}, {Command: libc.CString("oechoaround "), Command_pointer: do_osend, Subcmd: SCMD_OECHOAROUND}, {Command: libc.CString("oforce "), Command_pointer: do_oforce, Subcmd: 0}, {Command: libc.CString("ogoto "), Command_pointer: do_ogoto, Subcmd: 0}, {Command: libc.CString("oload "), Command_pointer: do_dgoload, Subcmd: 0}, {Command: libc.CString("opurge "), Command_pointer: do_opurge, Subcmd: 0}, {Command: libc.CString("orecho "), Command_pointer: do_orecho, Subcmd: 0}, {Command: libc.CString("osend "), Command_pointer: do_osend, Subcmd: SCMD_OSEND}, {Command: libc.CString("osetval "), Command_pointer: do_osetval, Subcmd: 0}, {Command: libc.CString("oteleport "), Command_pointer: do_oteleport, Subcmd: 0}, {Command: libc.CString("otimer "), Command_pointer: do_otimer, Subcmd: 0}, {Command: libc.CString("otransform "), Command_pointer: do_otransform, Subcmd: 0}, {Command: libc.CString("ozoneecho "), Command_pointer: do_ozoneecho, Subcmd: 0}, {Command: libc.CString("\n"), Command_pointer: nil, Subcmd: 0}}

func obj_command_interpreter(obj *obj_data, argument *byte) {
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
	}(); *obj_cmd_info[cmd].Command != '\n'; cmd++ {
		if libc.StrNCmp(obj_cmd_info[cmd].Command, &arg[0], length) == 0 {
			break
		}
	}
	if *obj_cmd_info[cmd].Command == '\n' {
		obj_log(obj, libc.CString("Unknown object cmd: '%s'"), argument)
	} else {
		(obj_cmd_info[cmd].Command_pointer)(obj, line, cmd, obj_cmd_info[cmd].Subcmd)
	}
}
