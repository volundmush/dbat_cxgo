package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const APPLY_TYPE = 1
const AFFECT_TYPE = 2

func do_dg_cast(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, cmd *byte) {
	var (
		caster      *char_data = nil
		tch         *char_data = nil
		tobj        *obj_data  = nil
		caster_room *room_data = nil
		s           *byte
		t           *byte
		spellnum    int
		target      int = 0
		buf2        [64936]byte
		orig_cmd    [2048]byte
	)
	switch type_ {
	case MOB_TRIGGER:
		caster = (*char_data)(gohere)
	case WLD_TRIGGER:
		caster_room = (*room_data)(gohere)
	case OBJ_TRIGGER:
		caster_room = dg_room_of_obj((*obj_data)(gohere))
		if caster_room == nil {
			script_log(libc.CString("dg_do_cast: unknown room for object-caster!"))
			return
		}
	default:
		script_log(libc.CString("dg_do_cast: unknown trigger type!"))
		return
	}
	libc.StrCpy(&orig_cmd[0], cmd)
	s = libc.StrTok(cmd, libc.CString("'"))
	if s == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_cast needs spell name."), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	s = libc.StrTok(nil, libc.CString("'"))
	if s == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_cast needs spell name in `'s."), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	t = libc.StrTok(nil, libc.CString("\x00"))
	spellnum = find_skill_num(s, 1<<0)
	if spellnum < 1 || (spellnum >= SKILL_TABLE_SIZE || (skill_type(spellnum)&(1<<0)) == 0) {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_cast: invalid spell name (%s)"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &orig_cmd[0])
		return
	}
	if t != nil {
		one_argument(libc.StrCpy(&buf2[0], t), t)
		skip_spaces(&t)
	}
	if (spell_info[spellnum].Targets & (1 << 0)) != 0 {
		target = TRUE
	} else if t != nil && *t != 0 {
		if target == 0 && ((spell_info[spellnum].Targets&(1<<1)) != 0 || (spell_info[spellnum].Targets&(1<<2)) != 0) {
			if (func() *char_data {
				tch = get_char(t)
				return tch
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 && ((spell_info[spellnum].Targets&(1<<7)) != 0 || (spell_info[spellnum].Targets&(1<<10)) != 0 || (spell_info[spellnum].Targets&(1<<8)) != 0 || (spell_info[spellnum].Targets&(1<<9)) != 0) {
			if (func() *obj_data {
				tobj = get_obj(t)
				return tobj
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 {
			script_log(libc.CString("Trigger: %s, VNum %d. dg_cast: target not found (%s)"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &orig_cmd[0])
			return
		}
	}
	if (spell_info[spellnum].Routines & (1 << 5)) != 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_cast: group spells not permitted (%s)"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &orig_cmd[0])
		return
	}
	if caster == nil {
		caster = read_mobile(DG_CASTER_PROXY, VIRTUAL)
		if caster == nil {
			script_log(libc.CString("dg_cast: Cannot load the caster mob!"))
			return
		}
		if type_ == OBJ_TRIGGER {
			caster.Short_descr = libc.StrDup(((*obj_data)(gohere)).Short_description)
		} else if type_ == WLD_TRIGGER {
			caster.Short_descr = libc.CString("The gods")
		}
		caster.Next_in_room = caster_room.People
		caster_room.People = caster
		caster.In_room = real_room(caster_room.Number)
		call_magic(caster, tch, tobj, spellnum, DG_SPELL_LEVEL, CAST_SPELL, t)
		extract_char(caster)
	} else {
		call_magic(caster, tch, tobj, spellnum, GET_LEVEL(caster), CAST_SPELL, t)
	}
}
func do_dg_affect(gohere unsafe.Pointer, sc *script_data, trig *trig_data, script_type int, cmd *byte) {
	var (
		ch         *char_data = nil
		value      int        = 0
		duration   int        = 0
		junk       [2048]byte
		charname   [2048]byte
		property   [2048]byte
		value_p    [2048]byte
		duration_p [2048]byte
		i          int = 0
		type_      int = 0
		af         affected_type
	)
	half_chop(cmd, &junk[0], cmd)
	half_chop(cmd, &charname[0], cmd)
	half_chop(cmd, &property[0], cmd)
	half_chop(cmd, &value_p[0], &duration_p[0])
	if charname == nil || charname[0] == 0 || property == nil || property[0] == 0 || value_p == nil || value_p[0] == 0 || duration_p == nil || duration_p[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_affect usage: <target> <property> <value> <duration>"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	value = libc.Atoi(libc.GoString(&value_p[0]))
	duration = libc.Atoi(libc.GoString(&duration_p[0]))
	if duration <= 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_affect: need positive duration!"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		script_log(libc.CString("Line was: dg_affect %s %s %s %s (%d)"), &charname[0], &property[0], &value_p[0], &duration_p[0], duration)
		return
	}
	i = 0
	for libc.StrCaseCmp(apply_types[i], libc.CString("\n")) != 0 {
		if libc.StrCaseCmp(apply_types[i], &property[0]) == 0 {
			type_ = APPLY_TYPE
			break
		}
		i++
	}
	if type_ == 0 {
		i = 0
		for libc.StrCaseCmp(affected_bits[i], libc.CString("\n")) != 0 {
			if libc.StrCaseCmp(affected_bits[i], &property[0]) == 0 {
				type_ = AFFECT_TYPE
				break
			}
			i++
		}
	}
	if type_ == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_affect: unknown property '%s'!"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &property[0])
		return
	}
	ch = get_char(&charname[0])
	if ch == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. dg_affect: cannot locate target!"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	if libc.StrCaseCmp(&value_p[0], libc.CString("off")) == 0 {
		affect_from_char(ch, SPELL_DG_AFFECT)
		return
	}
	af.Type = SPELL_DG_AFFECT
	af.Duration = int16(duration - 1)
	af.Modifier = value
	if type_ == APPLY_TYPE {
		af.Location = i
		af.Bitvector = 0
	} else {
		af.Location = 0
		af.Bitvector = bitvector_t(int32(i))
	}
	affect_to_char(ch, &af)
}
func send_char_pos(ch *char_data, dam int) {
	switch ch.Position {
	case POS_MORTALLYW:
		act(libc.CString("$n is mortally wounded, and will die soon, if not aided."), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You are mortally wounded, and will die soon, if not aided.\r\n"))
	case POS_INCAP:
		act(libc.CString("$n is incapacitated and will slowly die, if not aided."), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You are incapacitated and will slowly die, if not aided.\r\n"))
	case POS_STUNNED:
		act(libc.CString("$n is stunned, but will probably regain consciousness again."), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You're stunned, but will probably regain consciousness again.\r\n"))
	case POS_DEAD:
		act(libc.CString("$n is dead!  R.I.P."), FALSE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You are dead!  Sorry...\r\n"))
	default:
		if dam > int(ch.Max_hit>>2) {
			act(libc.CString("That really did HURT!"), FALSE, ch, nil, nil, TO_CHAR)
		}
		if ch.Hit < (ch.Max_hit >> 2) {
			send_to_char(ch, libc.CString("@rYou wish that your wounds would stop BLEEDING so much!@n\r\n"))
		}
	}
}
func valid_dg_target(ch *char_data, bitvector int) int {
	if IS_NPC(ch) {
		return TRUE
	} else if ch.Admlevel < ADMLVL_IMMORT {
		return TRUE
	} else if (bitvector&(1<<0)) == 0 && (ch.Admlevel >= 2 && !PRF_FLAGGED(ch, PRF_TEST)) {
		return FALSE
	} else if !PRF_FLAGGED(ch, PRF_NOHASSLE) || PRF_FLAGGED(ch, PRF_TEST) {
		return TRUE
	} else {
		return FALSE
	}
}
func script_damage(vict *char_data, dam int) {
	if ADM_FLAGGED(vict, ADM_NODAMAGE) && dam > 0 {
		send_to_char(vict, libc.CString("Being the cool immortal you are, you sidestep a trap, obviously placed to kill you.\r\n"))
		return
	}
	if vict.Suppressed <= 0 {
		vict.Hit -= int64(dam)
		vict.Hit = int64(MIN(int(vict.Hit), int(vict.Max_hit)))
	} else if vict.Suppressed-int64(dam) >= 0 {
		vict.Suppressed -= int64(dam)
	} else {
		dam -= int(vict.Suppressed)
		vict.Hit -= int64(dam)
		vict.Hit = int64(MIN(int(vict.Hit), int(vict.Max_hit)))
	}
	update_pos(vict)
	send_char_pos(vict, dam)
	if int(vict.Position) == POS_DEAD {
		if !IS_NPC(vict) {
			mudlog(BRF, 0, TRUE, libc.CString("%s killed by script at %s"), GET_NAME(vict), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Name)
		}
		die(vict, nil)
	}
}
