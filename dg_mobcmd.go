package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func mob_log(mob *char_data, format *byte, _rest ...interface{}) {
	var (
		args   libc.ArgList
		output [64936]byte
	)
	stdio.Snprintf(&output[0], int(64936), "Mob (%s, VNum %d):: %s", mob.Short_descr, GET_MOB_VNUM(mob), format)
	args.Start(format, _rest)
	script_vlog(&output[0], args)
	args.End()
}
func do_masound(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		was_in_room room_rnum
		door        int
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if *argument == 0 {
		mob_log(ch, libc.CString("masound called with no argument"))
		return
	}
	skip_spaces(&argument)
	was_in_room = ch.In_room
	for door = 0; door < NUM_OF_DIRS; door++ {
		var newexit *room_direction_data
		if (func() *room_direction_data {
			newexit = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in_room)))).Dir_option[door]
			return newexit
		}()) != nil && newexit.To_room != room_rnum(-1) && newexit.To_room != was_in_room {
			ch.In_room = newexit.To_room
			sub_write(argument, ch, TRUE, TO_ROOM)
		}
	}
	ch.In_room = was_in_room
}
func do_mheal(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		mob_log(ch, libc.CString("mheal called without an argument"))
		return
	}
	var amount int64 = 0
	_ = amount
	var num float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
	var perc float64 = num * 0.01
	amount = int64(float64(ch.Max_hit) * perc)
	if libc.StrCaseCmp(&arg[0], libc.CString("pl")) == 0 {
		ch.Hit += int64(float64(ch.Max_hit) * num)
		if ch.Hit > ch.Max_hit {
			ch.Hit = ch.Max_hit
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("ki")) == 0 {
		ch.Mana += int64(float64(ch.Max_mana) * num)
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("st")) == 0 {
		ch.Move += int64(float64(ch.Max_move) * num)
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
	} else {
		mob_log(ch, libc.CString("mheal called with wrong argument [pl | ki | st]"))
		return
	}
}
func do_mkill(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
		buf    [2048]byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mkill called with no argument"))
		return
	}
	if arg[0] == UID_CHAR {
		if (func() *char_data {
			victim = get_char(&arg[0])
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("mkill: victim (%s) not found"), &arg[0])
			return
		}
	} else if (func() *char_data {
		victim = get_char_room_vis(ch, &arg[0], nil)
		return victim
	}()) == nil {
		mob_log(ch, libc.CString("mkill: victim (%s) not found"), &arg[0])
		return
	}
	if victim == ch {
		mob_log(ch, libc.CString("mkill: victim is self"))
		return
	}
	if !IS_NPC(victim) && PRF_FLAGGED(victim, PRF_NOHASSLE) {
		mob_log(ch, libc.CString("mkill: target has nohassle on"))
		return
	}
	stdio.Sprintf(&buf[0], "%s", GET_NAME(victim))
	if IS_HUMANOID(ch) {
		switch rand_number(1, 7) {
		case 1:
			do_punch(ch, &buf[0], 0, 0)
		case 2:
			do_kick(ch, &buf[0], 0, 0)
		case 3:
			do_elbow(ch, &buf[0], 0, 0)
		case 4:
			do_knee(ch, &buf[0], 0, 0)
		case 5:
			do_kick(ch, &buf[0], 0, 0)
		default:
			do_punch(ch, &buf[0], 0, 0)
		}
	} else {
		do_bite(ch, &buf[0], 0, 0)
	}
	return
}
func do_mjunk(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		pos      int
		junk_all int = 0
		obj      *obj_data
		obj_next *obj_data
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mjunk called with no argument"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("all")) == 0 {
		junk_all = 1
	}
	if find_all_dots(&arg[0]) != FIND_INDIV && junk_all == 0 {
		if (func() int {
			pos = get_obj_pos_in_equip_vis(ch, &arg[0], nil, ch.Equipment[:])
			return pos
		}()) >= 0 {
			extract_obj(unequip_char(ch, pos))
			return
		}
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) != nil {
			extract_obj(obj)
		}
		return
	} else {
		for obj = ch.Carrying; obj != nil; obj = obj_next {
			obj_next = obj.Next_content
			if arg[3] == '\x00' || isname(&arg[4], obj.Name) != 0 {
				extract_obj(obj)
			}
		}
		for (func() int {
			pos = get_obj_pos_in_equip_vis(ch, &arg[0], nil, ch.Equipment[:])
			return pos
		}()) >= 0 {
			extract_obj(unequip_char(ch, pos))
		}
	}
	return
}
func do_mechoaround(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
		p      *byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	p = one_argument(argument, &arg[0])
	skip_spaces(&p)
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mechoaround called with no argument"))
		return
	}
	if arg[0] == UID_CHAR {
		if (func() *char_data {
			victim = get_char(&arg[0])
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("mechoaround: victim (%s) does not exist"), &arg[0])
			return
		}
	} else if (func() *char_data {
		victim = get_char_room_vis(ch, &arg[0], nil)
		return victim
	}()) == nil {
		mob_log(ch, libc.CString("mechoaround: victim (%s) does not exist"), &arg[0])
		return
	}
	var buf [64936]byte
	stdio.Sprintf(&buf[0], libc.GoString(p))
	search_replace(&buf[0], GET_NAME(victim), libc.CString("$n"))
	act(&buf[0], TRUE, victim, nil, nil, TO_ROOM)
}
func do_msend(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
		p      *byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	p = one_argument(argument, &arg[0])
	skip_spaces(&p)
	if arg[0] == 0 {
		mob_log(ch, libc.CString("msend called with no argument"))
		return
	}
	if arg[0] == UID_CHAR {
		if (func() *char_data {
			victim = get_char(&arg[0])
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("msend: victim (%s) does not exist"), &arg[0])
			return
		}
	} else if (func() *char_data {
		victim = get_char_room_vis(ch, &arg[0], nil)
		return victim
	}()) == nil {
		mob_log(ch, libc.CString("msend: victim (%s) does not exist"), &arg[0])
		return
	}
	sub_write(p, victim, TRUE, TO_CHAR)
}
func do_mecho(ch *char_data, argument *byte, cmd int, subcmd int) {
	var p *byte
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if *argument == 0 {
		mob_log(ch, libc.CString("mecho called with no arguments"))
		return
	}
	p = argument
	skip_spaces(&p)
	sub_write(p, ch, TRUE, TO_ROOM)
}
func do_mzoneecho(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		zone        int
		room_number [2048]byte
		buf         [2048]byte
		msg         *byte
	)
	msg = any_one_arg(argument, &room_number[0])
	skip_spaces(&msg)
	if room_number[0] == 0 || *msg == 0 {
		mob_log(ch, libc.CString("mzoneecho called with too few args"))
	} else if (func() int {
		zone = int(real_zone_by_thing(room_vnum(libc.Atoi(libc.GoString(&room_number[0])))))
		return zone
	}()) == int(-1) {
		mob_log(ch, libc.CString("mzoneecho called for nonexistant zone"))
	} else {
		stdio.Sprintf(&buf[0], "%s\r\n", msg)
		send_to_zone(&buf[0], zone_rnum(zone))
	}
}
func do_mload(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		return
	}
	target = two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 || is_number(&arg2[0]) == 0 || (func() int {
		number = libc.Atoi(libc.GoString(&arg2[0]))
		return number
	}()) < 0 {
		mob_log(ch, libc.CString("mload: bad syntax"))
		return
	}
	if is_abbrev(&arg1[0], libc.CString("mob")) != 0 {
		var rnum room_rnum
		if target == nil || *target == 0 {
			rnum = ch.In_room
		} else {
			if !unicode.IsDigit(rune(*target)) || (func() room_rnum {
				rnum = real_room(room_vnum(libc.Atoi(libc.GoString(target))))
				return rnum
			}()) == room_rnum(-1) {
				mob_log(ch, libc.CString("mload: room target vnum doesn't exist (loading mob vnum %d to room %s)"), number, target)
				return
			}
		}
		if (func() *char_data {
			mob = read_mobile(mob_vnum(number), VIRTUAL)
			return mob
		}()) == nil {
			mob_log(ch, libc.CString("mload: bad mob vnum"))
			return
		}
		char_to_room(mob, rnum)
		if ch.Script != nil {
			var buf [2048]byte
			stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, mob.Id)
			add_var(&ch.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
		}
		load_mtrigger(mob)
	} else if is_abbrev(&arg1[0], libc.CString("obj")) != 0 {
		if (func() *obj_data {
			object = read_object(obj_vnum(number), VIRTUAL)
			return object
		}()) == nil {
			mob_log(ch, libc.CString("mload: bad object vnum"))
			return
		}
		if ch.Script != nil {
			var buf [2048]byte
			stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, object.Id)
			add_var(&ch.Script.Global_vars, libc.CString("lastloaded"), &buf[0], 0)
		}
		randomize_eq(object)
		if target == nil || *target == 0 {
			add_unique_id(object)
			if OBJWEAR_FLAGGED(object, ITEM_WEAR_TAKE) {
				obj_to_char(object, ch)
			} else {
				obj_to_room(object, ch.In_room)
			}
			load_otrigger(object)
			return
		}
		two_arguments(target, &arg1[0], &arg2[0])
		if arg1 != nil && arg1[0] == UID_CHAR {
			tch = get_char(&arg1[0])
		} else {
			tch = get_char_room_vis(ch, &arg1[0], nil)
		}
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
		if arg1 != nil && arg1[0] == UID_CHAR {
			cnt = get_obj(&arg1[0])
		} else {
			cnt = get_obj_vis(ch, &arg1[0], nil)
		}
		if cnt != nil && int(cnt.Type_flag) == ITEM_CONTAINER {
			obj_to_obj(object, cnt)
			load_otrigger(object)
			return
		}
		add_unique_id(object)
		obj_to_room(object, ch.In_room)
		load_otrigger(object)
		return
	} else {
		mob_log(ch, libc.CString("mload: bad type"))
	}
}
func do_mpurge(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
		obj    *obj_data
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		var (
			vnext    *char_data
			obj_next *obj_data
		)
		for victim = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; victim != nil; victim = vnext {
			vnext = victim.Next_in_room
			if IS_NPC(victim) && victim != ch {
				extract_char(victim)
			}
		}
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj_next {
			obj_next = obj.Next_content
			extract_obj(obj)
		}
		return
	}
	if arg[0] == UID_CHAR {
		victim = get_char(&arg[0])
	} else {
		victim = get_char_room_vis(ch, &arg[0], nil)
	}
	if victim == nil {
		if arg[0] == UID_CHAR {
			obj = get_obj(&arg[0])
		} else {
			obj = get_obj_vis(ch, &arg[0], nil)
		}
		if obj != nil {
			extract_obj(obj)
			obj = nil
		} else {
			mob_log(ch, libc.CString("mpurge: bad argument"))
		}
		return
	}
	if !IS_NPC(victim) {
		mob_log(ch, libc.CString("mpurge: purging a PC"))
		return
	}
	if victim == ch {
		dg_owner_purged = 1
	}
	extract_char(victim)
}
func do_mgoto(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		location room_rnum
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mgoto called with no argument"))
		return
	}
	if (func() room_rnum {
		location = find_target_room(ch, &arg[0])
		return location
	}()) == room_rnum(-1) && GET_MOB_VNUM(ch) != 3 {
		mob_log(ch, libc.CString("mgoto: invalid location"))
		return
	} else if (func() room_rnum {
		location = find_target_room(ch, &arg[0])
		return location
	}()) == room_rnum(-1) {
		return
	}
	if ch.Fighting != nil {
		stop_fighting(ch)
	}
	char_from_room(ch)
	char_to_room(ch, location)
	enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
}
func do_mat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		location room_rnum
		original room_rnum
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	argument = one_argument(argument, &arg[0])
	if arg[0] == 0 || *argument == 0 {
		mob_log(ch, libc.CString("mat: bad argument"))
		return
	}
	if (func() room_rnum {
		location = find_target_room(ch, &arg[0])
		return location
	}()) == room_rnum(-1) {
		mob_log(ch, libc.CString("mat: invalid location"))
		return
	}
	original = ch.In_room
	char_from_room(ch)
	char_to_room(ch, location)
	command_interpreter(ch, argument)
	if ch.In_room == location {
		char_from_room(ch)
		char_to_room(ch, original)
	}
}
func do_mteleport(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1    [2048]byte
		arg2    [2048]byte
		target  room_rnum
		vict    *char_data
		next_ch *char_data
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	argument = two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 {
		mob_log(ch, libc.CString("mteleport: bad syntax"))
		return
	}
	target = find_target_room(ch, &arg2[0])
	if target == room_rnum(-1) {
		mob_log(ch, libc.CString("mteleport target is an invalid room"))
		return
	}
	if libc.StrCaseCmp(&arg1[0], libc.CString("all")) == 0 {
		if target == ch.In_room {
			mob_log(ch, libc.CString("mteleport all target is itself"))
			return
		}
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_ch {
			next_ch = vict.Next_in_room
			if valid_dg_target(vict, 1<<0) != 0 {
				char_from_room(vict)
				char_to_room(vict, target)
				enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
			}
		}
	} else {
		if arg1[0] == UID_CHAR {
			if (func() *char_data {
				vict = get_char(&arg1[0])
				return vict
			}()) == nil {
				mob_log(ch, libc.CString("mteleport: victim (%s) does not exist"), &arg1[0])
				return
			}
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg1[0], nil, 1<<1)
			return vict
		}()) == nil {
			mob_log(ch, libc.CString("mteleport: victim (%s) does not exist"), &arg1[0])
			return
		}
		if valid_dg_target(ch, 1<<0) != 0 {
			char_from_room(vict)
			char_to_room(vict, target)
			enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
		}
	}
}
func do_mdamage(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		name   [2048]byte
		amount [2048]byte
		dam    int = 0
		vict   *char_data
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	two_arguments(argument, &name[0], &amount[0])
	if name[0] == 0 || amount[0] == 0 {
		mob_log(ch, libc.CString("mdamage: bad syntax"))
		return
	}
	dam = libc.Atoi(libc.GoString(&amount[0]))
	if name[0] == UID_CHAR {
		if (func() *char_data {
			vict = get_char(&name[0])
			return vict
		}()) == nil {
			mob_log(ch, libc.CString("mdamage: victim (%s) does not exist"), &name[0])
			return
		}
	} else if (func() *char_data {
		vict = get_char_room_vis(ch, &name[0], nil)
		return vict
	}()) == nil {
		mob_log(ch, libc.CString("mdamage: victim (%s) does not exist"), &name[0])
		return
	}
	script_damage(vict, dam)
}
func do_mforce(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		return
	}
	argument = one_argument(argument, &arg[0])
	if arg[0] == 0 || *argument == 0 {
		mob_log(ch, libc.CString("mforce: bad syntax"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("all")) == 0 {
		var (
			i   *descriptor_data
			vch *char_data
		)
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Character != ch && i.Connected == 0 && i.Character.In_room == ch.In_room {
				vch = i.Character
				if GET_LEVEL(vch) < GET_LEVEL(ch) && CAN_SEE(ch, vch) && valid_dg_target(vch, 0) != 0 {
					command_interpreter(vch, argument)
				}
			}
		}
	} else {
		var victim *char_data
		if arg[0] == UID_CHAR {
			if (func() *char_data {
				victim = get_char(&arg[0])
				return victim
			}()) == nil {
				mob_log(ch, libc.CString("mforce: victim (%s) does not exist"), &arg[0])
				return
			}
		} else if (func() *char_data {
			victim = get_char_room_vis(ch, &arg[0], nil)
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("mforce: no such victim"))
			return
		}
		if victim == ch {
			mob_log(ch, libc.CString("mforce: forcing self"))
			return
		}
		if valid_dg_target(victim, 0) != 0 {
			command_interpreter(victim, argument)
		}
	}
}
func do_mremember(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim *char_data
		mem    *script_memory
		arg    [2048]byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		return
	}
	argument = one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mremember: bad syntax"))
		return
	}
	if arg[0] == UID_CHAR {
		if (func() *char_data {
			victim = get_char(&arg[0])
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("mremember: victim (%s) does not exist"), &arg[0])
			return
		}
	} else if (func() *char_data {
		victim = get_char_vis(ch, &arg[0], nil, 1<<1)
		return victim
	}()) == nil {
		mob_log(ch, libc.CString("mremember: victim (%s) does not exist"), &arg[0])
		return
	}
	mem = new(script_memory)
	if ch.Memory == nil {
		ch.Memory = mem
	} else {
		var tmpmem *script_memory = ch.Memory
		for tmpmem.Next != nil {
			tmpmem = tmpmem.Next
		}
		tmpmem.Next = mem
	}
	mem.Id = int(victim.Id)
	if argument != nil && *argument != 0 {
		mem.Cmd = libc.StrDup(argument)
	}
}
func do_mforget(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim *char_data
		mem    *script_memory
		prev   *script_memory
		arg    [2048]byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mforget: bad syntax"))
		return
	}
	if arg[0] == UID_CHAR {
		if (func() *char_data {
			victim = get_char(&arg[0])
			return victim
		}()) == nil {
			mob_log(ch, libc.CString("mforget: victim (%s) does not exist"), &arg[0])
			return
		}
	} else if (func() *char_data {
		victim = get_char_vis(ch, &arg[0], nil, 1<<1)
		return victim
	}()) == nil {
		mob_log(ch, libc.CString("mforget: victim (%s) does not exist"), &arg[0])
		return
	}
	mem = ch.Memory
	prev = nil
	for mem != nil {
		if mem.Id == int(victim.Id) {
			if mem.Cmd != nil {
				libc.Free(unsafe.Pointer(mem.Cmd))
			}
			if prev == nil {
				ch.Memory = mem.Next
				libc.Free(unsafe.Pointer(mem))
				mem = ch.Memory
			} else {
				prev.Next = mem.Next
				libc.Free(unsafe.Pointer(mem))
				mem = prev.Next
			}
		} else {
			prev = mem
			mem = mem.Next
		}
	}
}
func do_mtransform(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg       [2048]byte
		m         *char_data
		tmpmob    char_data
		obj       [23]*obj_data
		this_rnum mob_rnum = ch.Nr
		pos       int
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	if ch.Desc != nil {
		send_to_char(ch, libc.CString("You've got no VNUM to return to, dummy! try 'switch'\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		mob_log(ch, libc.CString("mtransform: missing argument"))
	} else if !unicode.IsDigit(rune(arg[0])) && arg[0] != '-' {
		mob_log(ch, libc.CString("mtransform: bad argument"))
	} else {
		if unicode.IsDigit(rune(arg[0])) {
			m = read_mobile(mob_vnum(libc.Atoi(libc.GoString(&arg[0]))), VIRTUAL)
		} else {
			m = read_mobile(mob_vnum(libc.Atoi(libc.GoString(&arg[1]))), VIRTUAL)
		}
		if m == nil {
			mob_log(ch, libc.CString("mtransform: bad mobile vnum"))
			return
		}
		for pos = 0; pos < NUM_WEARS; pos++ {
			if (ch.Equipment[pos]) != nil {
				obj[pos] = unequip_char(ch, pos)
			} else {
				obj[pos] = nil
			}
		}
		char_to_room(m, ch.In_room)
		libc.MemCpy(unsafe.Pointer(&tmpmob), unsafe.Pointer(m), int(unsafe.Sizeof(char_data{})))
		if m.Name != nil {
			tmpmob.Name = libc.StrDup(m.Name)
		}
		if m.Title != nil {
			tmpmob.Title = libc.StrDup(m.Title)
		}
		if m.Short_descr != nil {
			tmpmob.Short_descr = libc.StrDup(m.Short_descr)
		}
		if m.Long_descr != nil {
			tmpmob.Long_descr = libc.StrDup(m.Long_descr)
		}
		if m.Description != nil {
			tmpmob.Description = libc.StrDup(m.Description)
		}
		tmpmob.Id = ch.Id
		tmpmob.Affected = ch.Affected
		tmpmob.Carrying = ch.Carrying
		tmpmob.Proto_script = ch.Proto_script
		tmpmob.Script = ch.Script
		tmpmob.Memory = ch.Memory
		tmpmob.Next_in_room = ch.Next_in_room
		tmpmob.Next = ch.Next
		tmpmob.Next_fighting = ch.Next_fighting
		tmpmob.Followers = ch.Followers
		tmpmob.Master = ch.Master
		tmpmob.Was_in_room = ch.Was_in_room
		tmpmob.Gold = ch.Gold
		tmpmob.Position = ch.Position
		tmpmob.Carry_weight = ch.Carry_weight
		tmpmob.Carry_items = ch.Carry_items
		tmpmob.Fighting = ch.Fighting
		libc.MemCpy(unsafe.Pointer(ch), unsafe.Pointer(&tmpmob), int(unsafe.Sizeof(char_data{})))
		for pos = 0; pos < NUM_WEARS; pos++ {
			if obj[pos] != nil {
				equip_char(ch, obj[pos], pos)
			}
		}
		ch.Nr = this_rnum
		extract_char(m)
	}
}
func do_mdoor(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	argument = two_arguments(argument, &target[0], &direction[0])
	value = one_argument(argument, &field[0])
	skip_spaces(&value)
	if target[0] == 0 || direction[0] == 0 || field[0] == 0 {
		mob_log(ch, libc.CString("mdoor called with too few args"))
		return
	}
	if (func() *room_data {
		rm = get_room(&target[0])
		return rm
	}()) == nil {
		mob_log(ch, libc.CString("mdoor: invalid target"))
		return
	}
	if (func() int {
		dir = search_block(&direction[0], &dirs[0], FALSE)
		return dir
	}()) == -1 {
		mob_log(ch, libc.CString("mdoor: invalid direction"))
		return
	}
	if (func() int {
		fd = search_block(&field[0], &door_field[0], FALSE)
		return fd
	}()) == -1 {
		mob_log(ch, libc.CString("odoor: invalid field"))
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
				mob_log(ch, libc.CString("mdoor: invalid door target"))
			}
		}
	}
}
func do_mfollow(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		leader *char_data
		j      *follow_type
		k      *follow_type
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	one_argument(argument, &buf[0])
	if buf[0] == 0 {
		mob_log(ch, libc.CString("mfollow: bad syntax"))
		return
	}
	if buf[0] == UID_CHAR {
		if (func() *char_data {
			leader = get_char(&buf[0])
			return leader
		}()) == nil {
			mob_log(ch, libc.CString("mfollow: victim (%s) does not exist"), &buf[0])
			return
		}
	} else if (func() *char_data {
		leader = get_char_vis(ch, &buf[0], nil, 1<<0)
		return leader
	}()) == nil {
		mob_log(ch, libc.CString("mfollow: victim (%s) not found"), &buf[0])
		return
	}
	if ch.Master == leader {
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master != nil {
		return
	}
	if ch.Master != nil {
		if ch.Master.Followers.Follower == ch {
			k = ch.Master.Followers
			ch.Master.Followers = k.Next
			libc.Free(unsafe.Pointer(k))
		} else {
			for k = ch.Master.Followers; k.Next.Follower != ch; k = k.Next {
			}
			j = k.Next
			k.Next = j.Next
			libc.Free(unsafe.Pointer(j))
		}
		ch.Master = nil
	}
	if ch == leader {
		return
	}
	if circle_follow(ch, leader) {
		mob_log(ch, libc.CString("mfollow: Following in circles."))
		return
	}
	ch.Master = leader
	k = new(follow_type)
	k.Follower = ch
	k.Next = leader.Followers
	leader.Followers = k
}
func do_mrecho(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		start  [2048]byte
		finish [2048]byte
		msg    *byte
	)
	if !IS_NPC(ch) || ch.Desc != nil && ch.Desc.Original.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("Huh?!?\r\n"))
		return
	}
	msg = two_arguments(argument, &start[0], &finish[0])
	skip_spaces(&msg)
	if *msg == 0 || start[0] == 0 || finish[0] == 0 || is_number(&start[0]) == 0 || is_number(&finish[0]) == 0 {
		mob_log(ch, libc.CString("mrecho called with too few args"))
	} else {
		send_to_range(room_vnum(libc.Atoi(libc.GoString(&start[0]))), room_vnum(libc.Atoi(libc.GoString(&finish[0]))), libc.CString("%s\r\n"), msg)
	}
}
