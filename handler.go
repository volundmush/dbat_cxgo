package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unicode"
	"unsafe"
)

const FIND_INDIV = 0
const FIND_ALL = 1
const FIND_ALLDOT = 2
const FIND_CHAR_ROOM = 1
const FIND_CHAR_WORLD = 2
const FIND_OBJ_INV = 4
const FIND_OBJ_ROOM = 8
const FIND_OBJ_WORLD = 16
const FIND_OBJ_EQUIP = 32
const WHITESPACE = " \t"

var extractions_pending int = 0

func get_i_name(ch *char_data, vict *char_data) *byte {
	var (
		fname  [40]byte
		filler [50]byte
		scrap  [100]byte
		line   [256]byte
		name   [50]byte
		known  int = FALSE
		fl     *stdio.File
	)
	if vict == nil {
		return libc.CString("")
	}
	if IS_NPC(ch) || IS_NPC(vict) {
		return JUGGLERACE(vict)
	}
	if vict == ch {
		return libc.CString("")
	}
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		return JUGGLERACE(vict)
	} else if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		return JUGGLERACE(vict)
	}
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%s %s\n", &filler[0], &scrap[0])
		if libc.StrCaseCmp(GET_NAME(vict), &filler[0]) == 0 {
			stdio.Sprintf(&name[0], "%s", &scrap[0])
			known = TRUE
		}
	}
	fl.Close()
	if known == TRUE {
		return &name[0]
	} else {
		return JUGGLERACE(vict)
	}
}
func fname(namelist *byte) *byte {
	var (
		holder [256]byte
		point  *byte
	)
	for point = &holder[0]; libc.IsAlpha(rune(*namelist)); func() *byte {
		namelist = (*byte)(unsafe.Add(unsafe.Pointer(namelist), 1))
		return func() *byte {
			p := &point
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
	}() {
		*point = *namelist
	}
	*point = '\x00'
	return &holder[0]
}
func is_name(str *byte, namelist *byte) int {
	var (
		curname *byte
		curstr  *byte
	)
	if *str == 0 || *namelist == 0 || str == nil || namelist == nil {
		return 0
	}
	curname = namelist
	for {
		for curstr = str; ; func() *byte {
			curstr = (*byte)(unsafe.Add(unsafe.Pointer(curstr), 1))
			return func() *byte {
				p := &curname
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		}() {
			if *curstr == 0 && !libc.IsAlpha(rune(*curname)) {
				return 1
			}
			if *curname == 0 {
				return 0
			}
			if *curstr == 0 || *curname == ' ' {
				break
			}
			if unicode.ToLower(rune(*curstr)) != unicode.ToLower(rune(*curname)) {
				break
			}
		}
		for ; libc.IsAlpha(rune(*curname)); curname = (*byte)(unsafe.Add(unsafe.Pointer(curname), 1)) {
		}
		if *curname == 0 {
			return 0
		}
		curname = (*byte)(unsafe.Add(unsafe.Pointer(curname), 1))
	}
}
func isname(str *byte, namelist *byte) int {
	basic_mud_log(libc.CString("PANIC! REIMPLEMENT THIS!"))
	os.Exit(-1)
	return 0
}
func aff_apply_modify(ch *char_data, loc int, mod int, spec int, msg *byte) {
	switch loc {
	case APPLY_NONE:
	case APPLY_STR:
		ch.Aff_abils.Str += int8(mod)
	case APPLY_DEX:
		ch.Aff_abils.Dex += int8(mod)
	case APPLY_INT:
		ch.Aff_abils.Intel += int8(mod)
	case APPLY_WIS:
		ch.Aff_abils.Wis += int8(mod)
	case APPLY_CON:
		ch.Aff_abils.Con += int8(mod)
	case APPLY_CHA:
		ch.Aff_abils.Cha += int8(mod)
	case APPLY_CLASS:
	case APPLY_LEVEL:
	case APPLY_AGE:
		ch.Time.Birth -= libc.Time(mod * (((int(SECS_PER_MUD_HOUR * 24)) * 30) * 12))
	case APPLY_CHAR_WEIGHT:
		ch.Weight += uint8(int8(mod))
	case APPLY_CHAR_HEIGHT:
		ch.Height += uint8(int8(mod))
	case APPLY_MANA:
		ch.Max_mana += int64(mod)
	case APPLY_HIT:
		ch.Max_hit += int64(mod)
	case APPLY_MOVE:
		ch.Max_move += int64(mod)
	case APPLY_KI:
		ch.Max_ki += int64(mod)
	case APPLY_GOLD:
	case APPLY_EXP:
	case APPLY_AC:
		ch.Armor += mod
	case APPLY_ACCURACY:
		ch.Accuracy += mod
	case APPLY_DAMAGE:
		ch.Damage_mod += mod
	case APPLY_REGEN:
		ch.Regen += mod
	case APPLY_TRAIN:
		ch.Asb += mod
	case APPLY_LIFEMAX:
		ch.Lifebonus += mod
	case APPLY_UNUSED3:
		fallthrough
	case APPLY_UNUSED4:
	case APPLY_RACE:
	case APPLY_TURN_LEVEL:
		ch.Player_specials.Tlevel += mod
	case APPLY_SPELL_LVL_0:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_0] += mod
		}
	case APPLY_SPELL_LVL_1:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_1] += mod
		}
	case APPLY_SPELL_LVL_2:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_2] += mod
		}
	case APPLY_SPELL_LVL_3:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_3] += mod
		}
	case APPLY_SPELL_LVL_4:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_4] += mod
		}
	case APPLY_SPELL_LVL_5:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_5] += mod
		}
	case APPLY_SPELL_LVL_6:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_6] += mod
		}
	case APPLY_SPELL_LVL_7:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_7] += mod
		}
	case APPLY_SPELL_LVL_8:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_8] += mod
		}
	case APPLY_SPELL_LVL_9:
		if !IS_NPC(ch) {
			ch.Player_specials.Spell_level[SPELL_LEVEL_9] += mod
		}
	case APPLY_FORTITUDE:
		ch.Apply_saving_throw[SAVING_FORTITUDE] += int16(mod)
	case APPLY_REFLEX:
		ch.Apply_saving_throw[SAVING_REFLEX] += int16(mod)
	case APPLY_WILL:
		ch.Apply_saving_throw[SAVING_WILL] += int16(mod)
	case APPLY_SKILL:
		for {
			ch.Skillmods[spec] = int8(int(ch.Skillmods[spec]) + mod)
			if true {
				break
			}
		}
	case APPLY_FEAT:
		ch.Feats[spec] += int8(mod)
	case APPLY_ALLSAVES:
		ch.Apply_saving_throw[SAVING_FORTITUDE] += int16(mod)
		ch.Apply_saving_throw[SAVING_REFLEX] += int16(mod)
		ch.Apply_saving_throw[SAVING_WILL] += int16(mod)
	case APPLY_ALL_STATS:
		ch.Aff_abils.Str += int8(mod)
		ch.Aff_abils.Intel += int8(mod)
		ch.Aff_abils.Wis += int8(mod)
		ch.Aff_abils.Dex += int8(mod)
		ch.Aff_abils.Con += int8(mod)
		ch.Aff_abils.Cha += int8(mod)
	case APPLY_RESISTANCE:
	default:
		basic_mud_log(libc.CString("SYSERR: Unknown apply adjust %d attempt (%s, affect_modify)."), loc, "__FILE__")
	}
}
func affect_modify(ch *char_data, loc int, mod int, spec int, bitv bitvector_t, add bool) {
	if add {
		if bitv != AFF_INFRAVISION || int(ch.Race) != RACE_ANDROID {
			SET_BIT_AR(ch.Affected_by[:], bitvector_t(int32(bitv)))
		}
	} else {
		if bitv != AFF_INFRAVISION || int(ch.Race) != RACE_ANDROID {
			REMOVE_BIT_AR(ch.Affected_by[:], bitvector_t(int32(bitv)))
			mod = -mod
		}
	}
	aff_apply_modify(ch, loc, mod, spec, libc.CString("affect_modify"))
}
func affect_modify_ar(ch *char_data, loc int, mod int, spec int, bitv []bitvector_t, add bool) {
	var (
		i int
		j int
	)
	if add {
		for i = 0; i < AF_ARRAY_MAX; i++ {
			for j = 0; j < 32; j++ {
				if IS_SET_AR(bitv, bitvector_t(int32((i*32)+j))) {
					if (i*32)+j != AFF_INFRAVISION || int(ch.Race) != RACE_ANDROID {
						SET_BIT_AR(ch.Affected_by[:], bitvector_t(int32((i*32)+j)))
					}
				}
			}
		}
	} else {
		for i = 0; i < AF_ARRAY_MAX; i++ {
			for j = 0; j < 32; j++ {
				if IS_SET_AR(bitv, bitvector_t(int32((i*32)+j))) {
					if (i*32)+j != AFF_INFRAVISION || int(ch.Race) != RACE_ANDROID {
						REMOVE_BIT_AR(ch.Affected_by[:], bitvector_t(int32((i*32)+j)))
					}
				}
			}
		}
		mod = -mod
	}
	aff_apply_modify(ch, loc, mod, spec, libc.CString("affect_modify_ar"))
}
func affect_total(ch *char_data) {
	var (
		af *affected_type
		i  int
		j  int
	)
	ch.Spellfail = func() int16 {
		p := &ch.Armorcheck
		ch.Armorcheck = func() int16 {
			p := &ch.Armorcheckall
			ch.Armorcheckall = 0
			return *p
		}()
		return *p
	}()
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			for j = 0; j < MAX_OBJ_AFFECT; j++ {
				affect_modify_ar(ch, (ch.Equipment[i]).Affected[j].Location, (ch.Equipment[i]).Affected[j].Modifier, (ch.Equipment[i]).Affected[j].Specific, ch.Equipment[i].Bitvector[:], false)
			}
		}
	}
	for af = ch.Affected; af != nil; af = af.Next {
		affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, FALSE != 0)
	}
	ch.Aff_abils = ch.Real_abils
	ch.Apply_saving_throw[SAVING_FORTITUDE] = int16(int(ch.Feats[FEAT_GREAT_FORTITUDE]) * 3)
	ch.Apply_saving_throw[SAVING_REFLEX] = int16(int(ch.Feats[FEAT_LIGHTNING_REFLEXES]) * 3)
	ch.Apply_saving_throw[SAVING_WILL] = int16(int(ch.Feats[FEAT_IRON_WILL]) * 3)
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			if int((ch.Equipment[i]).Type_flag) == ITEM_ARMOR {
				ch.Spellfail += int16((ch.Equipment[i]).Value[VAL_ARMOR_SPELLFAIL])
				ch.Armorcheckall += int16((ch.Equipment[i]).Value[VAL_ARMOR_CHECK])
				if is_proficient_with_armor(ch, (ch.Equipment[i]).Value[VAL_ARMOR_SKILL]) == 0 {
					ch.Armorcheck += int16((ch.Equipment[i]).Value[VAL_ARMOR_CHECK])
				}
			}
			for j = 0; j < MAX_OBJ_AFFECT; j++ {
				affect_modify_ar(ch, (ch.Equipment[i]).Affected[j].Location, (ch.Equipment[i]).Affected[j].Modifier, ch.Equipment[i].Affected[j].Specific, ch.Equipment[i].Bitvector[:], TRUE != 0)
			}
		}
	}
	for af = ch.Affected; af != nil; af = af.Next {
		affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, TRUE != 0)
	}
	if (ch.Bonuses[BONUS_WIMP]) > 0 {
		ch.Aff_abils.Str = int8(MAX(0, MIN(int64(ch.Aff_abils.Str), 45)))
	} else {
		ch.Aff_abils.Str = int8(MAX(0, MIN(int64(ch.Aff_abils.Str), 100)))
	}
	if (ch.Bonuses[BONUS_DULL]) > 0 {
		ch.Aff_abils.Intel = int8(MAX(0, MIN(int64(ch.Aff_abils.Intel), 45)))
	} else {
		ch.Aff_abils.Intel = int8(MAX(0, MIN(int64(ch.Aff_abils.Intel), 100)))
	}
	if (ch.Bonuses[BONUS_FOOLISH]) > 0 {
		ch.Aff_abils.Wis = int8(MAX(0, MIN(int64(ch.Aff_abils.Wis), 45)))
	} else {
		ch.Aff_abils.Wis = int8(MAX(0, MIN(int64(ch.Aff_abils.Wis), 100)))
	}
	if (ch.Bonuses[BONUS_SLOW]) > 0 {
		ch.Aff_abils.Cha = int8(MAX(0, MIN(int64(ch.Aff_abils.Cha), 45)))
	} else {
		ch.Aff_abils.Cha = int8(MAX(0, MIN(int64(ch.Aff_abils.Cha), 100)))
	}
	if (ch.Bonuses[BONUS_CLUMSY]) > 0 {
		ch.Aff_abils.Dex = int8(MAX(0, MIN(int64(ch.Aff_abils.Dex), 45)))
	} else {
		ch.Aff_abils.Dex = int8(MAX(0, MIN(int64(ch.Aff_abils.Dex), 100)))
	}
	if (ch.Bonuses[BONUS_FRAIL]) > 0 {
		ch.Aff_abils.Con = int8(MAX(0, MIN(int64(ch.Aff_abils.Con), 45)))
	} else {
		ch.Aff_abils.Con = int8(MAX(0, MIN(int64(ch.Aff_abils.Con), 100)))
	}
}
func affect_to_char(ch *char_data, af *affected_type) {
	var affected_alloc *affected_type
	affected_alloc = new(affected_type)
	if ch.Affected == nil {
		ch.Next_affect = affect_list
		affect_list = ch
	}
	*affected_alloc = *af
	affected_alloc.Next = ch.Affected
	ch.Affected = affected_alloc
	affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, TRUE != 0)
	affect_total(ch)
}
func affect_remove(ch *char_data, af *affected_type) {
	var cmtemp *affected_type
	if ch.Affected == nil {
		core_dump_real(libc.CString("__FILE__"), 0)
		return
	}
	affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, FALSE != 0)
	if af == ch.Affected {
		ch.Affected = af.Next
	} else {
		cmtemp = ch.Affected
		for cmtemp != nil && cmtemp.Next != af {
			cmtemp = cmtemp.Next
		}
		if cmtemp != nil {
			cmtemp.Next = af.Next
		}
	}
	libc.Free(unsafe.Pointer(af))
	affect_total(ch)
	if ch.Affected == nil {
		var temp *char_data
		if ch == affect_list {
			affect_list = ch.Next_affect
		} else {
			temp = affect_list
			for temp != nil && temp.Next_affect != ch {
				temp = temp.Next_affect
			}
			if temp != nil {
				temp.Next_affect = ch.Next_affect
			}
		}
		ch.Next_affect = nil
	}
}
func affect_from_char(ch *char_data, type_ int) {
	var (
		hjp  *affected_type
		next *affected_type
	)
	for hjp = ch.Affected; hjp != nil; hjp = next {
		next = hjp.Next
		if int(hjp.Type) == type_ {
			affect_remove(ch, hjp)
		}
	}
}
func affectv_from_char(ch *char_data, type_ int) {
	var (
		hjp  *affected_type
		next *affected_type
	)
	for hjp = ch.Affectedv; hjp != nil; hjp = next {
		next = hjp.Next
		if int(hjp.Type) == type_ {
			affectv_remove(ch, hjp)
		}
	}
}
func affected_by_spell(ch *char_data, type_ int) bool {
	var hjp *affected_type
	for hjp = ch.Affected; hjp != nil; hjp = hjp.Next {
		if int(hjp.Type) == type_ {
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func affectedv_by_spell(ch *char_data, type_ int) bool {
	var hjp *affected_type
	for hjp = ch.Affectedv; hjp != nil; hjp = hjp.Next {
		if int(hjp.Type) == type_ {
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func affect_join(ch *char_data, af *affected_type, add_dur bool, avg_dur bool, add_mod bool, avg_mod bool) {
	var (
		hjp   *affected_type
		next  *affected_type
		found bool = FALSE != 0
	)
	for hjp = ch.Affected; !found && hjp != nil; hjp = next {
		next = hjp.Next
		if int(hjp.Type) == int(af.Type) && hjp.Location == af.Location {
			if add_dur {
				af.Duration += hjp.Duration
			}
			if avg_dur {
				af.Duration /= 2
			}
			if add_mod {
				af.Modifier += hjp.Modifier
			}
			if avg_mod {
				af.Modifier /= 2
			}
			affect_remove(ch, hjp)
			affect_to_char(ch, af)
			found = TRUE != 0
		}
	}
	if !found {
		affect_to_char(ch, af)
	}
}
func char_from_room(ch *char_data) {
	var (
		temp *char_data
		i    int
	)
	if ch == nil || ch.In_room == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: NULL character or NOWHERE in %s, char_from_room"), "__FILE__")
		return
	}
	if ch.Fighting != nil && !AFF_FLAGGED(ch, AFF_PURSUIT) {
		stop_fighting(ch)
	}
	if AFF_FLAGGED(ch, AFF_PURSUIT) && ch.Fighting == nil {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_PURSUIT)
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			if int((ch.Equipment[i]).Type_flag) == ITEM_LIGHT {
				if ((ch.Equipment[i]).Value[VAL_LIGHT_HOURS]) != 0 {
					world[ch.In_room].Light--
				}
			}
		}
	}
	if PLR_FLAGGED(ch, PLR_AURALIGHT) {
		world[ch.In_room].Light--
	}
	if ch == world[ch.In_room].People {
		world[ch.In_room].People = ch.Next_in_room
	} else {
		temp = world[ch.In_room].People
		for temp != nil && temp.Next_in_room != ch {
			temp = temp.Next_in_room
		}
		if temp != nil {
			temp.Next_in_room = ch.Next_in_room
		}
	}
	ch.In_room = -1
	ch.Next_in_room = nil
}
func char_to_room(ch *char_data, room room_rnum) {
	var i int
	if ch == nil || room == room_rnum(-1) || room > top_of_world {
		basic_mud_log(libc.CString("SYSERR: Illegal value(s) passed to char_to_room. (Room: %d/%d Ch: %p"), room, top_of_world, ch)
	} else {
		ch.Next_in_room = world[room].People
		world[room].People = ch
		ch.In_room = room
		for i = 0; i < NUM_WEARS; i++ {
			if (ch.Equipment[i]) != nil {
				if int((ch.Equipment[i]).Type_flag) == ITEM_LIGHT {
					if ((ch.Equipment[i]).Value[VAL_LIGHT_HOURS]) != 0 {
						world[room].Light++
					}
				}
			}
		}
		if PLR_FLAGGED(ch, PLR_AURALIGHT) {
			world[room].Light++
		}
		if ch.Fighting != nil && ch.In_room != ch.Fighting.In_room && !AFF_FLAGGED(ch, AFF_PURSUIT) {
			stop_fighting(ch.Fighting)
			stop_fighting(ch)
		}
		if !IS_NPC(ch) {
			if PRF_FLAGGED(ch, PRF_ARENAWATCH) {
				REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_ARENAWATCH)
				ch.Arenawatch = -1
			}
		}
	}
}
func obj_to_char(object *obj_data, ch *char_data) {
	if object != nil && ch != nil {
		object.Next_content = ch.Carrying
		ch.Carrying = object
		object.Carried_by = ch
		object.In_room = -1
		ch.Carry_weight += int(object.Weight)
		ch.Carry_items++
		if ch.Kaioken <= 0 && !AFF_FLAGGED(ch, AFF_METAMORPH) && !OBJ_FLAGGED(object, ITEM_THROW) {
			if ch.Hit > gear_pl(ch) {
				ch.Hit = gear_pl(ch)
			}
			if ch.Hit <= 0 {
				ch.Hit = 1
			}
		} else if ch.Hit > gear_pl(ch) {
			if ch.Kaioken > 0 {
				send_to_char(ch, libc.CString("@RThe strain of the weight has reduced your kaioken somewhat!@n\n"))
				ch.Hit -= object.Weight * 5
				if ch.Hit <= 0 {
					ch.Hit = 1
				}
			} else if AFF_FLAGGED(ch, AFF_METAMORPH) {
				send_to_char(ch, libc.CString("@RYour metamorphosis strains under the additional weight!@n\n"))
				ch.Hit -= object.Weight * 5
				if ch.Hit <= 0 {
					ch.Hit = 1
				}
			}
		} else if ch.Hit <= gear_pl(ch) && ch.Kaioken > 0 {
			send_to_char(ch, libc.CString("You've dropped out of kaioken due to the weight!\r\n"))
			ch.Kaioken = 0
			if ch.Hit <= 0 {
				ch.Hit = 1
			}
		}
		if (object.Value[0]) != 0 {
			if GET_OBJ_VNUM(object) == 0x4141 || GET_OBJ_VNUM(object) == 0x4142 || GET_OBJ_VNUM(object) == 0x4143 {
				object.Level = object.Value[0]
			}
		}
		if !IS_NPC(ch) {
			SET_BIT_AR(ch.Act[:], PLR_CRASH)
		}
	} else {
		basic_mud_log(libc.CString("SYSERR: NULL obj (%p) or char (%p) passed to obj_to_char."), object, ch)
	}
}
func obj_from_char(object *obj_data) {
	var temp *obj_data
	if object == nil {
		basic_mud_log(libc.CString("SYSERR: NULL object passed to obj_from_char."))
		return
	}
	if object == object.Carried_by.Carrying {
		object.Carried_by.Carrying = object.Next_content
	} else {
		temp = object.Carried_by.Carrying
		for temp != nil && temp.Next_content != object {
			temp = temp.Next_content
		}
		if temp != nil {
			temp.Next_content = object.Next_content
		}
	}
	if !IS_NPC(object.Carried_by) {
		SET_BIT_AR(object.Carried_by.Act[:], PLR_CRASH)
	}
	var previous int64 = gear_pl(object.Carried_by)
	object.Carried_by.Carry_weight -= int(object.Weight)
	object.Carried_by.Carry_items--
	if object.Carried_by.Kaioken <= 0 && object.Carried_by.Hit >= gear_pl(object.Carried_by) {
		if gear_pl_restore(object.Carried_by, previous) > 0 {
			object.Carried_by.Hit += gear_pl_restore(object.Carried_by, previous)
			if object.Carried_by.Hit > object.Carried_by.Max_hit {
				object.Carried_by.Hit = object.Carried_by.Max_hit
			}
		}
	}
	if (object.Value[0]) != 0 {
		if GET_OBJ_VNUM(object) == 0x4141 || GET_OBJ_VNUM(object) == 0x4142 || GET_OBJ_VNUM(object) == 0x4143 {
			object.Level = object.Value[0]
		}
	}
	object.Carried_by = nil
	object.Next_content = nil
}
func apply_ac(ch *char_data, eq_pos int) int {
	if (ch.Equipment[eq_pos]) == nil {
		core_dump_real(libc.CString("__FILE__"), 0)
		return 0
	}
	if int((ch.Equipment[eq_pos]).Type_flag) != ITEM_ARMOR {
		return 0
	}
	return (ch.Equipment[eq_pos]).Value[VAL_ARMOR_APPLYAC]
}
func invalid_align(ch *char_data, obj *obj_data) int {
	if OBJ_FLAGGED(obj, ITEM_ANTI_EVIL) && IS_EVIL(ch) {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_GOOD) && IS_GOOD(ch) {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_NEUTRAL) && IS_NEUTRAL(ch) {
		return TRUE
	}
	return FALSE
}
func equip_char(ch *char_data, obj *obj_data, pos int) {
	var j int
	if pos < 0 || pos >= NUM_WEARS {
		core_dump_real(libc.CString("__FILE__"), 0)
		return
	}
	if (ch.Equipment[pos]) != nil {
		basic_mud_log(libc.CString("SYSERR: Char is already equipped: %s, %s"), GET_NAME(ch), obj.Short_description)
		return
	}
	if obj.Carried_by != nil {
		basic_mud_log(libc.CString("SYSERR: EQUIP: Obj is carried_by when equip."))
		return
	}
	if obj.In_room != room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: EQUIP: Obj is in_room when equip."))
		return
	}
	if invalid_align(ch, obj) != 0 || invalid_class(ch, obj) != 0 || invalid_race(ch, obj) != 0 {
		act(libc.CString("You stop wearing $p as something prevents you."), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("$n stops wearing $p as something prevents $m."), FALSE, ch, obj, nil, TO_ROOM)
		obj_to_char(obj, ch)
		return
	}
	ch.Equipment[pos] = obj
	obj.Worn_by = ch
	obj.Worn_on = int16(pos)
	if int(obj.Type_flag) == ITEM_ARMOR {
		ch.Armor += apply_ac(ch, pos)
	}
	if ch.In_room != room_rnum(-1) {
		if int(obj.Type_flag) == ITEM_LIGHT {
			if (obj.Value[VAL_LIGHT_HOURS]) != 0 {
				world[ch.In_room].Light++
			}
		}
	} else {
		basic_mud_log(libc.CString("SYSERR: IN_ROOM(ch) = NOWHERE when equipping char %s."), GET_NAME(ch))
	}
	for j = 0; j < MAX_OBJ_AFFECT; j++ {
		affect_modify_ar(ch, obj.Affected[j].Location, obj.Affected[j].Modifier, obj.Affected[j].Specific, obj.Bitvector[:], TRUE != 0)
	}
	affect_total(ch)
}
func unequip_char(ch *char_data, pos int) *obj_data {
	var (
		j   int
		obj *obj_data
	)
	if pos < 0 || pos >= NUM_WEARS || (ch.Equipment[pos]) == nil {
		core_dump_real(libc.CString("__FILE__"), 0)
		return nil
	}
	obj = ch.Equipment[pos]
	obj.Worn_by = nil
	obj.Worn_on = -1
	if int(obj.Type_flag) == ITEM_ARMOR {
		ch.Armor -= apply_ac(ch, pos)
	}
	if ch.In_room != room_rnum(-1) {
		if int(obj.Type_flag) == ITEM_LIGHT {
			if (obj.Value[VAL_LIGHT_HOURS]) != 0 {
				world[ch.In_room].Light--
			}
		}
	} else {
		basic_mud_log(libc.CString("SYSERR: IN_ROOM(ch) = NOWHERE when unequipping char %s."), GET_NAME(ch))
	}
	ch.Equipment[pos] = nil
	for j = 0; j < MAX_OBJ_AFFECT; j++ {
		affect_modify_ar(ch, obj.Affected[j].Location, obj.Affected[j].Modifier, obj.Affected[j].Specific, obj.Bitvector[:], FALSE != 0)
	}
	affect_total(ch)
	return obj
}
func get_number(name **byte) int {
	var (
		i      int
		ppos   *byte
		number [2048]byte
	)
	number[0] = '\x00'
	if (func() *byte {
		ppos = libc.StrChr(*name, '.')
		return ppos
	}()) != nil {
		*func() *byte {
			p := &ppos
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}() = '\x00'
		strlcpy(&number[0], *name, uint64(2048))
		libc.StrCpy(*name, ppos)
		for i = 0; number[i] != 0; i++ {
			if !unicode.IsDigit(rune(number[i])) {
				return 0
			}
		}
		return libc.Atoi(libc.GoString(&number[0]))
	}
	return 1
}
func get_obj_in_list_num(num int, list *obj_data) *obj_data {
	var i *obj_data
	for i = list; i != nil; i = i.Next_content {
		if i.Item_number == obj_vnum(num) {
			return i
		}
	}
	return nil
}
func get_obj_num(nr obj_rnum) *obj_data {
	var i *obj_data
	for i = object_list; i != nil; i = i.Next {
		if i.Item_number == obj_vnum(nr) {
			return i
		}
	}
	return nil
}
func get_char_room(name *byte, number *int, room room_rnum) *char_data {
	var (
		i   *char_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	if *number == 0 {
		return nil
	}
	for i = world[room].People; i != nil && *number != 0; i = i.Next_in_room {
		if isname(name, i.Name) != 0 {
			if func() int {
				p := number
				*p--
				return *p
			}() == 0 {
				return i
			}
		}
	}
	return nil
}
func get_char_num(nr mob_rnum) *char_data {
	var i *char_data
	for i = character_list; i != nil; i = i.Next {
		if i.Nr == nr {
			return i
		}
	}
	return nil
}
func obj_to_room(object *obj_data, room room_rnum) {
	var vehicle *obj_data = nil
	if object == nil || room == room_rnum(-1) || room > top_of_world {
		basic_mud_log(libc.CString("SYSERR: Illegal value(s) passed to obj_to_room. (Room #%d/%d, obj %p)"), room, top_of_world, object)
	} else {
		if ROOM_FLAGGED(room, ROOM_GARDEN1) || ROOM_FLAGGED(room, ROOM_GARDEN2) {
			if int(object.Type_flag) != ITEM_PLANT {
				send_to_room(room, libc.CString("%s @wDisappears in a puff of smoke! It seems the room was designed to vaporize anything not plant related. Strange...@n\r\n"), object.Short_description)
				extract_obj(object)
				return
			}
		}
		if room == real_room(80) {
			auc_load(object)
		}
		object.Next_content = world[room].Contents
		world[room].Contents = object
		object.In_room = room
		object.Carried_by = nil
		object.Lload = libc.GetTime(nil)
		if int(object.Type_flag) == ITEM_VEHICLE && !OBJ_FLAGGED(object, ITEM_UNBREAKABLE) && GET_OBJ_VNUM(object) > 0x4AFF {
			SET_BIT_AR(object.Extra_flags[:], ITEM_UNBREAKABLE)
		}
		if int(object.Type_flag) == ITEM_HATCH && GET_OBJ_VNUM(object) <= 0x4AFF {
			if GET_OBJ_VNUM(object) <= 0x4A37 && GET_OBJ_VNUM(object) >= 18800 || GET_OBJ_VNUM(object) <= 0x4AFF && GET_OBJ_VNUM(object) >= 19100 {
				var (
					hnum  int       = (object.Value[0])
					house *obj_data = read_object(obj_vnum(hnum), VIRTUAL)
				)
				obj_to_room(house, real_room(room_vnum(object.Value[6])))
				object.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
				object.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
			}
		}
		if int(object.Type_flag) == ITEM_HATCH && (object.Value[0]) > 1 && GET_OBJ_VNUM(object) > 0x4AFF {
			if (func() *obj_data {
				vehicle = find_vehicle_by_vnum(object.Value[VAL_HATCH_DEST])
				return vehicle
			}()) == nil {
				if real_room(room_vnum(object.Value[3])) != room_rnum(-1) {
					vehicle = read_object(obj_vnum(object.Value[0]), VIRTUAL)
					obj_to_room(vehicle, real_room(room_vnum(object.Value[3])))
					if object.Action_description != nil {
						if libc.StrLen(object.Action_description) != 0 {
							var (
								nick  [2048]byte
								nick2 [2048]byte
								nick3 [2048]byte
							)
							if GET_OBJ_VNUM(vehicle) <= 0xB413 && GET_OBJ_VNUM(vehicle) >= 46000 {
								stdio.Sprintf(&nick[0], "Saiyan Pod %s", object.Action_description)
								stdio.Sprintf(&nick2[0], "@wA @Ys@ya@Yi@yy@Ya@yn @Dp@Wo@Dd@w named @D(@C%s@D)@w", object.Action_description)
							} else if GET_OBJ_VNUM(vehicle) >= 46100 && GET_OBJ_VNUM(vehicle) <= 0xB477 {
								stdio.Sprintf(&nick[0], "EDI Xenofighter MK. II %s", object.Action_description)
								stdio.Sprintf(&nick2[0], "@wAn @YE@yD@YI @CX@ce@Wn@Do@Cf@ci@Wg@Dh@Wt@ce@Cr @RMK. II @wnamed @D(@C%s@D)@w", object.Action_description)
							}
							stdio.Sprintf(&nick3[0], "%s is resting here@w", &nick2[0])
							vehicle.Name = libc.StrDup(&nick[0])
							vehicle.Short_description = libc.StrDup(&nick2[0])
							vehicle.Description = libc.StrDup(&nick3[0])
						}
					}
					object.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
					object.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
				} else {
					basic_mud_log(libc.CString("Hatch load: Hatch with no vehicle load room: #%d!"), GET_OBJ_VNUM(object))
				}
			}
		}
		if (world[object.In_room].Dir_option[5]) != nil && (SECT(object.In_room) == SECT_UNDERWATER || SECT(object.In_room) == SECT_WATER_NOSWIM) {
			act(libc.CString("$p @Bsinks to deeper waters.@n"), TRUE, nil, object, nil, TO_ROOM)
			var numb int = int(libc.BoolToInt(GET_ROOM_VNUM((world[object.In_room].Dir_option[5]).To_room)))
			obj_from_room(object)
			obj_to_room(object, real_room(room_vnum(numb)))
		}
		if (world[object.In_room].Dir_option[5]) != nil && SECT(object.In_room) == SECT_FLYING && (GET_OBJ_VNUM(object) < 80 || GET_OBJ_VNUM(object) > 83) {
			act(libc.CString("$p @Cfalls down.@n"), TRUE, nil, object, nil, TO_ROOM)
			var numb int = int(libc.BoolToInt(GET_ROOM_VNUM((world[object.In_room].Dir_option[5]).To_room)))
			obj_from_room(object)
			obj_to_room(object, real_room(room_vnum(numb)))
			if SECT(object.In_room) != SECT_FLYING {
				act(libc.CString("$p @Cfalls down and smacks the ground.@n"), TRUE, nil, object, nil, TO_ROOM)
			}
		}
		if (object.Value[0]) != 0 {
			if GET_OBJ_VNUM(object) == 0x4141 || GET_OBJ_VNUM(object) == 0x4142 || GET_OBJ_VNUM(object) == 0x4143 {
				object.Level = object.Value[0]
			}
		}
		if ROOM_FLAGGED(room, ROOM_HOUSE) {
			SET_BIT_AR(world[room].Room_flags[:], ROOM_HOUSE_CRASH)
		}
	}
}
func obj_from_room(object *obj_data) {
	var temp *obj_data
	if object == nil || object.In_room == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: NULL object (%p) or obj not in a room (%d) passed to obj_from_room"), object, object.In_room)
		return
	}
	if object.Posted_to != nil && object.In_obj == nil {
		var obj *obj_data = object.Posted_to
		if object.Posttype <= 0 {
			send_to_room(obj.In_room, libc.CString("%s@W shakes loose from %s@W.@n\r\n"), obj.Short_description, object.Short_description)
		} else {
			send_to_room(obj.In_room, libc.CString("%s@W comes loose from %s@W.@n\r\n"), object.Short_description, obj.Short_description)
		}
		obj.Posted_to = nil
		obj.Posttype = 0
		object.Posted_to = nil
		object.Posttype = 0
	}
	if object == world[object.In_room].Contents {
		world[object.In_room].Contents = object.Next_content
	} else {
		temp = world[object.In_room].Contents
		for temp != nil && temp.Next_content != object {
			temp = temp.Next_content
		}
		if temp != nil {
			temp.Next_content = object.Next_content
		}
	}
	if ROOM_FLAGGED(object.In_room, ROOM_HOUSE) {
		SET_BIT_AR(world[object.In_room].Room_flags[:], ROOM_HOUSE_CRASH)
	}
	object.In_room = -1
	object.Next_content = nil
}
func obj_to_obj(obj *obj_data, obj_to *obj_data) {
	var tmp_obj *obj_data
	if obj == nil || obj_to == nil || obj == obj_to {
		basic_mud_log(libc.CString("SYSERR: NULL object (%p) or same source (%p) and target (%p VNUM: %d) obj passed to obj_to_obj."), obj, obj, obj_to, func() obj_vnum {
			if obj_to != nil {
				return GET_OBJ_VNUM(obj_to)
			}
			return -1
		}())
		return
	}
	obj.Next_content = obj_to.Contains
	obj_to.Contains = obj
	obj.In_obj = obj_to
	tmp_obj = obj.In_obj
	if (obj.In_obj.Value[VAL_CONTAINER_CAPACITY]) > 0 {
		for tmp_obj = obj.In_obj; tmp_obj.In_obj != nil; tmp_obj = tmp_obj.In_obj {
			tmp_obj.Weight += obj.Weight
		}
		tmp_obj.Weight += obj.Weight
		if tmp_obj.Carried_by != nil {
			tmp_obj.Carried_by.Carry_weight += int(obj.Weight)
		}
	}
	if obj_to.In_room != room_rnum(-1) && ROOM_FLAGGED(obj_to.In_room, ROOM_HOUSE) {
		SET_BIT_AR(world[obj_to.In_room].Room_flags[:], ROOM_HOUSE_CRASH)
	}
}
func obj_from_obj(obj *obj_data) {
	var (
		temp     *obj_data
		obj_from *obj_data
	)
	if obj.In_obj == nil {
		basic_mud_log(libc.CString("SYSERR: (%s): trying to illegally extract obj from obj."), "__FILE__")
		return
	}
	obj_from = obj.In_obj
	temp = obj.In_obj
	if obj == obj_from.Contains {
		obj_from.Contains = obj.Next_content
	} else {
		temp = obj_from.Contains
		for temp != nil && temp.Next_content != obj {
			temp = temp.Next_content
		}
		if temp != nil {
			temp.Next_content = obj.Next_content
		}
	}
	if (obj.In_obj.Value[VAL_CONTAINER_CAPACITY]) > 0 {
		for temp = obj.In_obj; temp.In_obj != nil; temp = temp.In_obj {
			temp.Weight -= obj.Weight
		}
		temp.Weight -= obj.Weight
		if temp.Carried_by != nil {
			temp.Carried_by.Carry_weight -= int(obj.Weight)
		}
	}
	if obj_from.In_room != room_rnum(-1) && ROOM_FLAGGED(obj_from.In_room, ROOM_HOUSE) {
		SET_BIT_AR(world[obj_from.In_room].Room_flags[:], ROOM_HOUSE_CRASH)
	}
	obj.In_obj = nil
	obj.Next_content = nil
}
func object_list_new_owner(list *obj_data, ch *char_data) {
	if list != nil {
		object_list_new_owner(list.Contains, ch)
		object_list_new_owner(list.Next_content, ch)
		list.Carried_by = ch
	}
}
func extract_obj(obj *obj_data) {
	var (
		temp *obj_data
		ch   *char_data
	)
	if obj.Worn_by != nil {
		if unequip_char(obj.Worn_by, int(obj.Worn_on)) != obj {
			basic_mud_log(libc.CString("SYSERR: Inconsistent worn_by and worn_on pointers!!"))
		}
	}
	if obj.In_room != room_rnum(-1) {
		obj_from_room(obj)
	} else if obj.Carried_by != nil {
		obj_from_char(obj)
	} else if obj.In_obj != nil {
		obj_from_obj(obj)
	}
	if obj.Fellow_wall != nil && GET_OBJ_VNUM(obj) == 79 {
		var trash *obj_data
		trash = obj.Fellow_wall
		obj.Fellow_wall = nil
		trash.Fellow_wall = nil
		extract_obj(trash)
	}
	if obj.Sitting != nil {
		ch = obj.Sitting
		obj.Sitting = nil
		ch.Sits = nil
	}
	if obj.Posted_to != nil && obj.In_obj == nil {
		var obj2 *obj_data = obj.Posted_to
		obj2.Posted_to = nil
		obj2.Posttype = 0
		obj.Posted_to = nil
	}
	if obj.Target != nil {
		obj.Target = nil
	}
	if obj.User != nil {
		obj.User = nil
	}
	for obj.Contains != nil {
		extract_obj(obj.Contains)
	}
	if obj == object_list {
		object_list = obj.Next
	} else {
		temp = object_list
		for temp != nil && temp.Next != obj {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = obj.Next
		}
	}
	if obj.Item_number != obj_vnum(-1) {
		obj_index[obj.Item_number].Number--
	}
	if obj.Script != nil {
		extract_script(unsafe.Pointer(obj), OBJ_TRIGGER)
	}
	if GET_OBJ_VNUM(obj) != 80 && GET_OBJ_VNUM(obj) != 81 {
		if obj.Item_number == obj_vnum(-1) || obj.Proto_script != obj_proto[obj.Item_number].Proto_script {
			free_proto_script(unsafe.Pointer(obj), OBJ_TRIGGER)
		}
	}
	free_obj(obj)
}
func update_object(obj *obj_data, use int) {
	if obj == nil {
		return
	}
	if (obj.Script == nil || !IS_SET(bitvector_t(int32(obj.Script.Types)), 1<<5)) && obj.Timer > 0 {
		obj.Timer -= use
	}
	if obj.Contains != nil {
		update_object(obj.Contains, use)
	}
	if obj.Next_content != nil {
		update_object(obj.Next_content, use)
	}
}
func update_char_objects(ch *char_data) {
	var (
		i int
		j int
	)
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			if int((ch.Equipment[i]).Type_flag) == ITEM_LIGHT && ((ch.Equipment[i]).Value[VAL_LIGHT_HOURS]) > 0 && ((ch.Equipment[i]).Value[VAL_LIGHT_TIME]) <= 0 {
				j = func() int {
					p := &((ch.Equipment[i]).Value[VAL_LIGHT_HOURS])
					*p--
					return *p
				}()
				(ch.Equipment[i]).Value[VAL_LIGHT_TIME] = 3
				if j == 1 {
					send_to_char(ch, libc.CString("Your light begins to flicker and fade.\r\n"))
					act(libc.CString("$n's light begins to flicker and fade."), FALSE, ch, nil, nil, TO_ROOM)
				} else if j == 0 {
					send_to_char(ch, libc.CString("Your light sputters out and dies.\r\n"))
					act(libc.CString("$n's light sputters out and dies."), FALSE, ch, nil, nil, TO_ROOM)
					world[ch.In_room].Light--
				}
			} else if int((ch.Equipment[i]).Type_flag) == ITEM_LIGHT && ((ch.Equipment[i]).Value[VAL_LIGHT_HOURS]) > 0 {
				(ch.Equipment[i]).Value[VAL_LIGHT_TIME] -= 1
			}
			update_object(ch.Equipment[i], 2)
		}
	}
	if ch.Carrying != nil {
		update_object(ch.Carrying, 1)
	}
}
func extract_char_final(ch *char_data) {
	var (
		k     *char_data
		temp  *char_data
		chair *obj_data
		d     *descriptor_data
		obj   *obj_data
		i     int
	)
	if ch.In_room == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: NOWHERE extracting char %s. (%s, extract_char_final)"), GET_NAME(ch), "__FILE__")
		os.Exit(1)
	}
	if !IS_NPC(ch) && ch.Desc == nil {
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Original == ch {
				do_return(d.Character, nil, 0, 0)
				break
			}
		}
	}
	if ch.Desc != nil {
		if ch.Desc.Original != nil {
			do_return(ch, nil, 0, 0)
		} else {
			for d = descriptor_list; d != nil; d = d.Next {
				if d == ch.Desc {
					continue
				}
				if d.Character != nil && int(ch.Idnum) == int(d.Character.Idnum) {
					d.Connected = CON_CLOSE
				}
			}
			ch.Desc.Connected = CON_MENU
			write_to_output(ch.Desc, libc.CString("%s"), config_info.Operation.MENU)
		}
	}
	if ch.Followers != nil || ch.Master != nil {
		die_follower(ch)
	}
	if ch.Sits != nil {
		chair = ch.Sits
		chair.Sitting = nil
		ch.Sits = nil
	}
	if IS_NPC(ch) && GET_MOB_VNUM(ch) == 25 {
		if ch.Original != nil {
			handle_multi_merge(ch)
		}
	}
	if !IS_NPC(ch) && int(ch.Clones) > 0 {
		var clone *char_data = nil
		for clone = character_list; clone != nil; clone = clone.Next {
			if IS_NPC(clone) {
				if GET_MOB_VNUM(clone) == 25 {
					if clone.Original == ch {
						handle_multi_merge(clone)
					}
				}
			}
		}
	}
	purge_homing(ch)
	if ch.Mindlink != nil {
		ch.Mindlink.Mindlink = nil
		ch.Mindlink = nil
	}
	if ch.Grappling != nil {
		act(libc.CString("@WYou stop grappling with @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappling), TO_CHAR)
		act(libc.CString("@C$n@W stops grappling with @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappling), TO_ROOM)
		ch.Grappling.Grap = -1
		ch.Grappling.Grappled = nil
		ch.Grappling = nil
		ch.Grap = -1
	}
	if ch.Grappled != nil {
		act(libc.CString("@WYou stop being grappled with by @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappled), TO_CHAR)
		act(libc.CString("@C$n@W stops being grappled with by @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappled), TO_ROOM)
		ch.Grappled.Grap = -1
		ch.Grappled.Grappling = nil
		ch.Grappled = nil
		ch.Grap = -1
	}
	if ch.Player_specials.Carrying != nil {
		carry_drop(ch, 3)
	}
	if ch.Player_specials.Carried_by != nil {
		carry_drop(ch.Player_specials.Carried_by, 3)
	}
	if ch.Poisonby != nil {
		ch.Poisonby = nil
	}
	if ch.Drag != nil {
		act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		ch.Drag.Dragged = nil
		ch.Drag = nil
	}
	if ch.Dragged != nil {
		act(libc.CString("@WYou stop being dragged by @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Dragged), TO_CHAR)
		act(libc.CString("@C$n@W stops being dragged by @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Dragged), TO_ROOM)
		ch.Dragged.Drag = nil
		ch.Dragged = nil
	}
	if ch.Defender != nil {
		ch.Defender.Defending = nil
		ch.Defender = nil
	}
	if ch.Defending != nil {
		ch.Defending.Defender = nil
		ch.Defending = nil
	}
	if ch.Blocked != nil {
		ch.Blocked.Blocks = nil
		ch.Blocked = nil
	}
	if ch.Blocks != nil {
		ch.Blocks.Blocked = nil
		ch.Blocks = nil
	}
	if ch.Absorbing != nil {
		ch.Absorbing.Absorbby = nil
		ch.Absorbing = nil
	}
	if ch.Absorbby != nil {
		ch.Absorbby.Absorbing = nil
		ch.Absorbby = nil
	}
	for ch.Carrying != nil {
		obj = ch.Carrying
		obj_from_char(obj)
		obj_to_room(obj, ch.In_room)
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			obj_to_room(unequip_char(ch, i), ch.In_room)
		}
	}
	if ch.Fighting != nil {
		stop_fighting(ch)
	}
	for k = combat_list; k != nil; k = temp {
		temp = k.Next_fighting
		if k.Fighting == ch {
			stop_fighting(k)
		}
	}
	char_from_room(ch)
	if IS_NPC(ch) {
		if ch.Nr != mob_rnum(-1) {
			mob_index[ch.Nr].Number--
		}
		clearMemory(ch)
		if ch.Script != nil {
			extract_script(unsafe.Pointer(ch), MOB_TRIGGER)
		}
		if ch.Memory != nil {
			extract_script_mem(ch.Memory)
		}
	} else {
		save_char(ch)
		Crash_delete_crashfile(ch)
	}
	if IS_NPC(ch) || ch.Desc == nil {
		free_char(ch)
	}
}
func extract_char(ch *char_data) {
	var (
		foll *follow_type
		i    int
		obj  *obj_data
	)
	if IS_NPC(ch) {
		if !IS_SET_AR(ch.Act[:], MOB_NOTDEADYET) {
			SET_BIT_AR(ch.Act[:], MOB_NOTDEADYET)
		} else {
			return
		}
	} else {
		if !IS_SET_AR(ch.Act[:], PLR_NOTDEADYET) {
			SET_BIT_AR(ch.Act[:], PLR_NOTDEADYET)
		} else {
			return
		}
	}
	for foll = ch.Followers; foll != nil; foll = foll.Next {
		if IS_NPC(foll.Follower) && AFF_FLAGGED(foll.Follower, AFF_CHARM) && (foll.Follower.In_room == ch.In_room || ch.In_room == 1) {
			for foll.Follower.Carrying != nil {
				obj = foll.Follower.Carrying
				obj_from_char(obj)
				obj_to_char(obj, ch)
			}
			for i = 0; i < NUM_WEARS; i++ {
				if (foll.Follower.Equipment[i]) != nil {
					obj = unequip_char(foll.Follower, i)
					obj_to_char(obj, ch)
				}
			}
			extract_char(foll.Follower)
		}
	}
	extractions_pending++
}
func extract_pending_chars() {
	var (
		vict      *char_data
		next_vict *char_data
		prev_vict *char_data
		temp      *char_data
	)
	if extractions_pending < 0 {
		basic_mud_log(libc.CString("SYSERR: Negative (%d) extractions pending."), extractions_pending)
	}
	for func() *char_data {
		vict = character_list
		return func() *char_data {
			prev_vict = nil
			return prev_vict
		}()
	}(); vict != nil && extractions_pending != 0; vict = next_vict {
		next_vict = vict.Next
		if MOB_FLAGGED(vict, MOB_NOTDEADYET) {
			REMOVE_BIT_AR(vict.Act[:], MOB_NOTDEADYET)
		} else if PLR_FLAGGED(vict, PLR_NOTDEADYET) {
			REMOVE_BIT_AR(vict.Act[:], PLR_NOTDEADYET)
		} else {
			prev_vict = vict
			continue
		}
		if vict == affect_list {
			affect_list = vict.Next_affect
		} else {
			temp = affect_list
			for temp != nil && temp.Next_affect != vict {
				temp = temp.Next_affect
			}
			if temp != nil {
				temp.Next_affect = vict.Next_affect
			}
		}
		if vict == affectv_list {
			affectv_list = vict.Next_affectv
		} else {
			temp = affectv_list
			for temp != nil && temp.Next_affectv != vict {
				temp = temp.Next_affectv
			}
			if temp != nil {
				temp.Next_affectv = vict.Next_affectv
			}
		}
		extract_char_final(vict)
		extractions_pending--
		if prev_vict != nil {
			prev_vict.Next = next_vict
		} else {
			character_list = next_vict
		}
	}
	if extractions_pending > 0 {
		basic_mud_log(libc.CString("SYSERR: Couldn't find %d extractions as counted."), extractions_pending)
	}
	extractions_pending = 0
}
func get_player_vis(ch *char_data, name *byte, number *int, inroom room_rnum) *char_data {
	var (
		i   *char_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	for i = character_list; i != nil; i = i.Next {
		if IS_NPC(i) {
			continue
		}
		if inroom == (1<<0) && i.In_room != ch.In_room {
			continue
		}
		if ch.Admlevel < 1 && i.Admlevel < 1 && !IS_NPC(ch) && !IS_NPC(i) {
			if libc.StrCaseCmp(JUGGLERACE(i), name) != 0 && libc.StrStr(JUGGLERACE(i), name) == nil {
				if readIntro(ch, i) == 1 {
					if libc.StrCaseCmp(get_i_name(ch, i), name) != 0 && libc.StrStr(get_i_name(ch, i), name) == nil {
						continue
					}
				} else {
					continue
				}
			}
		}
		if ch.Admlevel >= 1 || i.Admlevel >= 1 || IS_NPC(ch) || IS_NPC(i) {
			if libc.StrCaseCmp(i.Name, name) != 0 && libc.StrStr(i.Name, name) == nil {
				if libc.StrCaseCmp(JUGGLERACE(i), name) != 0 && libc.StrStr(JUGGLERACE(i), name) == nil {
					if !IS_NPC(ch) && !IS_NPC(i) && readIntro(ch, i) == 1 {
						if libc.StrCaseCmp(get_i_name(ch, i), name) != 0 && libc.StrStr(get_i_name(ch, i), name) == nil {
							continue
						}
					} else {
						continue
					}
				}
			}
		}
		if !CAN_SEE(ch, i) {
			continue
		}
		if func() int {
			p := number
			*p--
			return *p
		}() != 0 {
			continue
		}
		return i
	}
	return nil
}
func get_char_room_vis(ch *char_data, name *byte, number *int) *char_data {
	var (
		i   *char_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	if libc.StrCaseCmp(name, libc.CString("self")) == 0 || libc.StrCaseCmp(name, libc.CString("me")) == 0 {
		return ch
	}
	if *number == 0 {
		return get_player_vis(ch, name, nil, 1<<0)
	}
	for i = world[ch.In_room].People; i != nil && *number != 0; i = i.Next_in_room {
		if libc.StrCaseCmp(name, libc.CString("last")) == 0 && i.Lasthit != 0 && i.Lasthit == int(ch.Idnum) {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if isname(name, i.Name) != 0 && (IS_NPC(i) || IS_NPC(ch) || i.Admlevel > 0 || ch.Admlevel > 0) && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if isname(name, i.Name) != 0 && i == ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && !IS_NPC(ch) && libc.StrCaseCmp(get_i_name(ch, i), CAP(name)) == 0 && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && !IS_NPC(ch) && libc.StrStr(get_i_name(ch, i), CAP(name)) != nil && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && libc.StrCmp(JUGGLERACE(i), CAP(name)) == 0 && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && libc.StrStr(JUGGLERACE(i), CAP(name)) != nil && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && libc.StrCmp(JUGGLERACE(i), name) == 0 && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		} else if !IS_NPC(i) && libc.StrStr(JUGGLERACE(i), name) != nil && i != ch {
			if CAN_SEE(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		}
	}
	return nil
}
func get_char_world_vis(ch *char_data, name *byte, number *int) *char_data {
	var (
		i   *char_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	if (func() *char_data {
		i = get_char_room_vis(ch, name, number)
		return i
	}()) != nil {
		return i
	}
	if *number == 0 {
		return get_player_vis(ch, name, nil, 0)
	}
	for i = character_list; i != nil && *number != 0; i = i.Next {
		if ch.In_room == i.In_room {
			continue
		}
		if ch.Admlevel < 1 && i.Admlevel < 1 && !IS_NPC(ch) && !IS_NPC(i) {
			if libc.StrCaseCmp(JUGGLERACE(i), name) != 0 && libc.StrStr(JUGGLERACE(i), name) == nil {
				if readIntro(ch, i) == 1 {
					if libc.StrCaseCmp(get_i_name(ch, i), name) != 0 && libc.StrStr(get_i_name(ch, i), name) == nil {
						continue
					}
				} else {
					continue
				}
			}
		}
		if ch.Admlevel >= 1 || i.Admlevel >= 1 || IS_NPC(ch) || IS_NPC(i) {
			if libc.StrCaseCmp(i.Name, name) != 0 && libc.StrStr(i.Name, name) == nil {
				if libc.StrCaseCmp(JUGGLERACE(i), name) != 0 && libc.StrStr(JUGGLERACE(i), name) == nil {
					if !IS_NPC(ch) && !IS_NPC(i) && readIntro(ch, i) == 1 {
						if libc.StrCaseCmp(get_i_name(ch, i), name) != 0 && libc.StrStr(get_i_name(ch, i), name) == nil {
							continue
						}
					} else {
						continue
					}
				}
			}
		}
		if func() int {
			p := number
			*p--
			return *p
		}() != 0 {
			continue
		}
		return i
	}
	return nil
}
func get_char_vis(ch *char_data, name *byte, number *int, where int) *char_data {
	if where == (1 << 0) {
		return get_char_room_vis(ch, name, number)
	} else if where == (1 << 1) {
		return get_char_world_vis(ch, name, number)
	} else {
		return nil
	}
}
func get_obj_in_list_vis(ch *char_data, name *byte, number *int, list *obj_data) *obj_data {
	var (
		i   *obj_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	if *number == 0 {
		return nil
	}
	for i = list; i != nil && *number != 0; i = i.Next_content {
		if isname(name, i.Name) != 0 {
			if CAN_SEE_OBJ(ch, i) || int(i.Type_flag) == ITEM_LIGHT {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		}
	}
	return nil
}
func get_obj_vis(ch *char_data, name *byte, number *int) *obj_data {
	var (
		i   *obj_data
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&name)
	}
	if *number == 0 {
		return nil
	}
	if (func() *obj_data {
		i = get_obj_in_list_vis(ch, name, number, ch.Carrying)
		return i
	}()) != nil {
		return i
	}
	if (func() *obj_data {
		i = get_obj_in_list_vis(ch, name, number, world[ch.In_room].Contents)
		return i
	}()) != nil {
		return i
	}
	for i = object_list; i != nil && *number != 0; i = i.Next {
		if isname(name, i.Name) != 0 {
			if CAN_SEE_OBJ(ch, i) {
				if func() int {
					p := number
					*p--
					return *p
				}() == 0 {
					return i
				}
			}
		}
	}
	return nil
}
func get_obj_in_equip_vis(ch *char_data, arg *byte, number *int, equipment []*obj_data) *obj_data {
	var (
		j   int
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&arg)
	}
	if *number == 0 {
		return nil
	}
	for j = 0; j < NUM_WEARS; j++ {
		if equipment[j] != nil && CAN_SEE_OBJ(ch, equipment[j]) && isname(arg, equipment[j].Name) != 0 {
			if func() int {
				p := number
				*p--
				return *p
			}() == 0 {
				return equipment[j]
			}
		}
	}
	return nil
}
func get_obj_pos_in_equip_vis(ch *char_data, arg *byte, number *int, equipment []*obj_data) int {
	var (
		j   int
		num int
	)
	if number == nil {
		number = &num
		num = get_number(&arg)
	}
	if *number == 0 {
		return -1
	}
	for j = 0; j < NUM_WEARS; j++ {
		if equipment[j] != nil && CAN_SEE_OBJ(ch, equipment[j]) && isname(arg, equipment[j].Name) != 0 {
			if func() int {
				p := number
				*p--
				return *p
			}() == 0 {
				return j
			}
		}
	}
	return -1
}
func money_desc(amount int) *byte {
	var (
		cnt         int
		money_table [15]struct {
			Limit       int
			Description *byte
		} = [15]struct {
			Limit       int
			Description *byte
		}{{Limit: 1, Description: libc.CString("a single zenni")}, {Limit: 10, Description: libc.CString("a tiny pile of zenni")}, {Limit: 20, Description: libc.CString("a handful of zenni")}, {Limit: 75, Description: libc.CString("a little pile of zenni")}, {Limit: 150, Description: libc.CString("a small pile of zenni")}, {Limit: 250, Description: libc.CString("a pile of zenni")}, {Limit: 500, Description: libc.CString("a big pile of zenni")}, {Limit: 1000, Description: libc.CString("a large heap of zenni")}, {Limit: 5000, Description: libc.CString("a huge mound of zenni")}, {Limit: 10000, Description: libc.CString("an enormous mound of zenni")}, {Limit: 15000, Description: libc.CString("a small mountain of zenni")}, {Limit: 20000, Description: libc.CString("a mountain of zenni")}, {Limit: 25000, Description: libc.CString("a huge mountain of zenni")}, {Limit: 50000, Description: libc.CString("an enormous mountain of zenni")}, {}}
	)
	if amount <= 0 {
		basic_mud_log(libc.CString("SYSERR: Try to create negative or 0 money (%d)."), amount)
		return nil
	}
	for cnt = 0; money_table[cnt].Limit != 0; cnt++ {
		if amount <= money_table[cnt].Limit {
			return money_table[cnt].Description
		}
	}
	return libc.CString("an absolutely colossal mountain of zenni")
}
func create_money(amount int) *obj_data {
	var (
		obj       *obj_data
		new_descr *extra_descr_data
		buf       [200]byte
		y         int
	)
	if amount <= 0 {
		basic_mud_log(libc.CString("SYSERR: Try to create negative or 0 money. (%d)"), amount)
		return nil
	}
	obj = create_obj()
	new_descr = new(extra_descr_data)
	if amount == 1 {
		obj.Name = libc.CString("zenni money")
		obj.Short_description = libc.CString("a single zenni")
		obj.Description = libc.CString("One miserable zenni is lying here")
		new_descr.Keyword = libc.CString("zenni money")
		new_descr.Description = libc.CString("It's just one miserable little zenni.")
	} else {
		obj.Name = libc.CString("zenni money")
		obj.Short_description = libc.StrDup(money_desc(amount))
		stdio.Snprintf(&buf[0], int(200), "%s is lying here", money_desc(amount))
		obj.Description = libc.StrDup(CAP(&buf[0]))
		new_descr.Keyword = libc.CString("zenni money")
		if amount < 10 {
			stdio.Snprintf(&buf[0], int(200), "There is %d zenni.", amount)
		} else if amount < 100 {
			stdio.Snprintf(&buf[0], int(200), "There is about %d zenni.", (amount/10)*10)
		} else if amount < 1000 {
			stdio.Snprintf(&buf[0], int(200), "It looks to be about %d zenni.", (amount/100)*100)
		} else if amount < 100000 {
			stdio.Snprintf(&buf[0], int(200), "You guess there is, maybe, %d zenni.", ((amount/1000)+rand_number(0, amount/1000))*1000)
		} else {
			libc.StrCpy(&buf[0], libc.CString("There are is LOT of zenni."))
		}
		new_descr.Description = libc.StrDup(&buf[0])
	}
	new_descr.Next = nil
	obj.Ex_description = new_descr
	obj.Type_flag = ITEM_MONEY
	obj.Value[VAL_ALL_MATERIAL] = MATERIAL_GOLD
	obj.Value[VAL_ALL_MAXHEALTH] = 100
	obj.Value[VAL_ALL_HEALTH] = 100
	for y = 0; y < TW_ARRAY_MAX; y++ {
		obj.Wear_flags[y] = 0
	}
	SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_TAKE)
	obj.Value[VAL_MONEY_SIZE] = amount
	obj.Cost = amount
	obj.Item_number = -1
	return obj
}
func generic_find(arg *byte, bitvector bitvector_t, ch *char_data, tar_ch **char_data, tar_obj **obj_data) int {
	var (
		i        int
		found    int
		number   int
		name_val [2048]byte
		name     *byte = &name_val[0]
	)
	*tar_ch = nil
	*tar_obj = nil
	one_argument(arg, name)
	if *name == 0 {
		return 0
	}
	if (func() int {
		number = get_number(&name)
		return number
	}()) == 0 {
		return 0
	}
	if IS_SET(bitvector, 1<<0) {
		if (func() *char_data {
			p := tar_ch
			*tar_ch = get_char_room_vis(ch, name, &number)
			return *p
		}()) != nil {
			return 1 << 0
		}
	}
	if IS_SET(bitvector, 1<<1) {
		if (func() *char_data {
			p := tar_ch
			*tar_ch = get_char_world_vis(ch, name, &number)
			return *p
		}()) != nil {
			return 1 << 1
		}
	}
	if IS_SET(bitvector, 1<<5) {
		for func() int {
			found = FALSE
			return func() int {
				i = 0
				return i
			}()
		}(); i < NUM_WEARS && found == 0; i++ {
			if (ch.Equipment[i]) != nil && isname(name, (ch.Equipment[i]).Name) != 0 && func() int {
				p := &number
				*p--
				return *p
			}() == 0 {
				*tar_obj = ch.Equipment[i]
				found = TRUE
			}
		}
		if found != 0 {
			return 1 << 5
		}
	}
	if IS_SET(bitvector, 1<<2) {
		if (func() *obj_data {
			p := tar_obj
			*tar_obj = get_obj_in_list_vis(ch, name, &number, ch.Carrying)
			return *p
		}()) != nil {
			return 1 << 2
		}
	}
	if IS_SET(bitvector, 1<<3) {
		if (func() *obj_data {
			p := tar_obj
			*tar_obj = get_obj_in_list_vis(ch, name, &number, world[ch.In_room].Contents)
			return *p
		}()) != nil {
			return 1 << 3
		}
	}
	if IS_SET(bitvector, 1<<4) {
		if (func() *obj_data {
			p := tar_obj
			*tar_obj = get_obj_vis(ch, name, &number)
			return *p
		}()) != nil {
			return 1 << 4
		}
	}
	return 0
}
func find_all_dots(arg *byte) int {
	if libc.StrCmp(arg, libc.CString("all")) == 0 {
		return FIND_ALL
	} else if libc.StrNCmp(arg, libc.CString("all."), 4) == 0 {
		libc.StrCpy(arg, (*byte)(unsafe.Add(unsafe.Pointer(arg), 4)))
		return FIND_ALLDOT
	} else {
		return FIND_INDIV
	}
}
func affectv_to_char(ch *char_data, af *affected_type) {
	var affected_alloc *affected_type
	affected_alloc = new(affected_type)
	if ch.Affectedv == nil {
		ch.Next_affectv = affectv_list
		affectv_list = ch
	}
	*affected_alloc = *af
	affected_alloc.Next = ch.Affectedv
	ch.Affectedv = affected_alloc
	affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, TRUE != 0)
	affect_total(ch)
}
func affectv_remove(ch *char_data, af *affected_type) {
	var cmtemp *affected_type
	if ch.Affectedv == nil {
		core_dump_real(libc.CString("__FILE__"), 0)
		return
	}
	affect_modify(ch, af.Location, af.Modifier, af.Specific, af.Bitvector, FALSE != 0)
	if af == ch.Affectedv {
		ch.Affectedv = af.Next
	} else {
		cmtemp = ch.Affectedv
		for cmtemp != nil && cmtemp.Next != af {
			cmtemp = cmtemp.Next
		}
		if cmtemp != nil {
			cmtemp.Next = af.Next
		}
	}
	libc.Free(unsafe.Pointer(af))
	affect_total(ch)
	if ch.Affectedv == nil {
		var temp *char_data
		if ch == affectv_list {
			affectv_list = ch.Next_affectv
		} else {
			temp = affectv_list
			for temp != nil && temp.Next_affectv != ch {
				temp = temp.Next_affectv
			}
			if temp != nil {
				temp.Next_affectv = ch.Next_affectv
			}
		}
		ch.Next_affectv = nil
	}
}
func affectv_join(ch *char_data, af *affected_type, add_dur bool, avg_dur bool, add_mod bool, avg_mod bool) {
	var (
		hjp   *affected_type
		next  *affected_type
		found bool = FALSE != 0
	)
	for hjp = ch.Affectedv; !found && hjp != nil; hjp = next {
		next = hjp.Next
		if int(hjp.Type) == int(af.Type) && hjp.Location == af.Location {
			if add_dur {
				af.Duration += hjp.Duration
			}
			if avg_dur {
				af.Duration /= 2
			}
			if add_mod {
				af.Modifier += hjp.Modifier
			}
			if avg_mod {
				af.Modifier /= 2
			}
			affectv_remove(ch, hjp)
			affectv_to_char(ch, af)
			found = TRUE != 0
		}
	}
	if !found {
		affectv_to_char(ch, af)
	}
}
func is_better(object *obj_data, object2 *obj_data) int {
	var (
		value1 int = 0
		value2 int = 0
	)
	switch object.Type_flag {
	case ITEM_ARMOR:
		value1 = object.Value[VAL_ARMOR_APPLYAC]
		value2 = object2.Value[VAL_ARMOR_APPLYAC]
	case ITEM_WEAPON:
		value1 = ((object.Value[VAL_WEAPON_DAMSIZE]) + 1) * (object.Value[VAL_WEAPON_DAMDICE])
		value2 = ((object2.Value[VAL_WEAPON_DAMSIZE]) + 1) * (object2.Value[VAL_WEAPON_DAMDICE])
	default:
	}
	if value1 > value2 {
		return 1
	} else {
		return 0
	}
}
func item_check(object *obj_data, ch *char_data) {
	var where int = 0
	if IS_HUMANOID(ch) && libc.FuncAddr(mob_index[ch.Nr].Func) != libc.FuncAddr(shop_keeper) {
		if invalid_align(ch, object) != 0 || invalid_class(ch, object) != 0 {
			return
		}
		switch object.Type_flag {
		case ITEM_WEAPON:
			if (ch.Equipment[WEAR_WIELD1]) == nil {
				perform_wear(ch, object, WEAR_WIELD1)
			} else {
				if is_better(object, ch.Equipment[WEAR_WIELD1]) != 0 {
					perform_remove(ch, WEAR_WIELD1)
					perform_wear(ch, object, WEAR_WIELD1)
				}
			}
		case ITEM_ARMOR:
			fallthrough
		case ITEM_WORN:
			where = find_eq_pos(ch, object, nil)
			if (ch.Equipment[where]) == nil {
				perform_wear(ch, object, where)
			} else {
				if is_better(object, ch.Equipment[where]) != 0 {
					perform_remove(ch, where)
					perform_wear(ch, object, where)
				}
			}
		default:
		}
	}
}
