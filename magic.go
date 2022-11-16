package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const MAX_SPELL_AFFECTS = 5
const MOB_ELEMENTAL_BASE = 20
const MOB_ZOMBIE = 11
const MOB_AERIALSERVANT = 19

func affect_update() {
	var (
		af   *affected_type
		next *affected_type
		i    *char_data
	)
	for i = affect_list; i != nil; i = i.Next_affect {
		for af = i.Affected; af != nil; af = next {
			next = af.Next
			if int(af.Duration) >= 1 {
				af.Duration--
			} else if int(af.Duration) == 0 {
				if int(af.Type) > 0 {
					if af.Next == nil || int(af.Next.Type) != int(af.Type) || int(af.Next.Duration) > 0 {
						if spell_info[af.Type].Wear_off_msg != nil {
							send_to_char(i, libc.CString("%s\r\n"), spell_info[af.Type].Wear_off_msg)
						}
						if i.Speedboost > 0 && int(af.Type) == SPELL_HAYASA {
							i.Speedboost = 0
						}
						if int(af.Type) == SKILL_METAMORPH {
							if i.Hit > gear_pl(i) {
								i.Hit = gear_pl(i)
							}
						}
					}
				}
				affect_remove(i, af)
			}
		}
	}
}
func mag_materials(ch *char_data, item0 int, item1 int, item2 int, extract int, verbose int) int {
	var (
		tobj *obj_data
		obj0 *obj_data = nil
		obj1 *obj_data = nil
		obj2 *obj_data = nil
	)
	for tobj = ch.Carrying; tobj != nil; tobj = tobj.Next_content {
		if item0 > 0 && GET_OBJ_VNUM(tobj) == obj_vnum(item0) {
			obj0 = tobj
			item0 = -1
		} else if item1 > 0 && GET_OBJ_VNUM(tobj) == obj_vnum(item1) {
			obj1 = tobj
			item1 = -1
		} else if item2 > 0 && GET_OBJ_VNUM(tobj) == obj_vnum(item2) {
			obj2 = tobj
			item2 = -1
		}
	}
	if item0 > 0 || item1 > 0 || item2 > 0 {
		if verbose != 0 {
			switch rand_number(0, 2) {
			case 0:
				send_to_char(ch, libc.CString("A wart sprouts on your nose.\r\n"))
			case 1:
				send_to_char(ch, libc.CString("Your hair falls out in clumps.\r\n"))
			case 2:
				send_to_char(ch, libc.CString("A huge corn develops on your big toe.\r\n"))
			}
		}
		return FALSE
	}
	if extract != 0 {
		if item0 < 0 {
			extract_obj(obj0)
		}
		if item1 < 0 {
			extract_obj(obj1)
		}
		if item2 < 0 {
			extract_obj(obj2)
		}
	}
	if verbose != 0 {
		send_to_char(ch, libc.CString("A puff of smoke rises from your pack.\r\n"))
		act(libc.CString("A puff of smoke rises from $n's pack."), TRUE, ch, nil, nil, TO_ROOM)
	}
	return TRUE
}
func mag_newsaves(ch *char_data, victim *char_data, spellnum int, level int, cast_stat int) int {
	var (
		stype int
		dc    int
		total int
	)
	if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<0) {
		stype = SAVING_FORTITUDE
	} else if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<1) {
		stype = SAVING_REFLEX
	} else if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<2) {
		stype = SAVING_WILL
	} else {
		return FALSE
	}
	total = GET_SAVE(victim, stype) + rand_number(1, 20)
	dc = spell_info[spellnum].Spell_level + level + int(ability_mod_value(cast_stat))
	if ch != nil {
		if IS_SET(bitvector_t(int32(ch.School_feats[CFEAT_SPELL_FOCUS])), bitvector_t(int32(spell_info[spellnum].School))) {
			dc++
		}
		if IS_SET(bitvector_t(int32(ch.School_feats[CFEAT_GREATER_SPELL_FOCUS])), bitvector_t(int32(spell_info[spellnum].School))) {
			dc++
		}
	}
	if total >= dc {
		return TRUE
	}
	return FALSE
}
func mag_damage(level int, ch *char_data, victim *char_data, spellnum int) int {
	var dam int = 0
	_ = dam
	if victim == nil || ch == nil {
		return 0
	}
	switch spellnum {
	case SPELL_MAGIC_MISSILE:
		dam = dice(int(MIN(int64(level), 5)), 4) + int(MIN(int64(level), 5))
	case SPELL_CHILL_TOUCH:
		dam = dice(1, 6)
	case SPELL_BURNING_HANDS:
		dam = dice(int(MIN(int64(level), 5)), 4)
	case SPELL_SHOCKING_GRASP:
		dam = dice(1, 8) + int(MIN(int64(level), 20))
	case SPELL_LIGHTNING_BOLT:
		fallthrough
	case SPELL_FIREBALL:
		dam = dice(int(MIN(int64(level), 10)), 6)
	case SPELL_COLOR_SPRAY:
		dam = dice(1, 1) + 1
	case SPELL_DISPEL_EVIL:
		dam = dice(6, 8) + 6
		if IS_EVIL(ch) {
			victim = ch
			dam = int(ch.Hit - 1)
		} else if IS_GOOD(victim) {
			act(libc.CString("The gods protect $N."), FALSE, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			return 0
		}
	case SPELL_DISPEL_GOOD:
		dam = dice(6, 8) + 6
		if IS_GOOD(ch) {
			victim = ch
			dam = int(ch.Hit - 1)
		} else if IS_EVIL(victim) {
			act(libc.CString("The gods protect $N."), FALSE, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			return 0
		}
	case SPELL_CALL_LIGHTNING:
		dam = dice(int(MIN(int64(level), 10)), 10)
	case SPELL_HARM:
		dam = dice(8, 8) + 8
	case SPELL_ENERGY_DRAIN:
		if GET_LEVEL(victim) <= 2 {
			dam = 100
		} else {
			dam = dice(1, 10)
		}
	case SPELL_EARTHQUAKE:
		dam = dice(2, 8) + level
	case SPELL_INFLICT_LIGHT:
		dam = dice(1, 8) + int(MIN(int64(level), 5))
	case SPELL_INFLICT_CRITIC:
		dam = dice(4, 8) + int(MIN(int64(level), 20))
	case SPELL_ACID_SPLASH:
		fallthrough
	case SPELL_RAY_OF_FROST:
		dam = dice(1, 3)
	case SPELL_DISRUPT_UNDEAD:
		if AFF_FLAGGED(victim, AFF_UNDEAD) {
			dam = dice(1, 6)
		} else {
			send_to_char(ch, libc.CString("This magic only affects the undead!\r\n"))
			dam = 0
		}
	case SPELL_ICE_STORM:
		fallthrough
	case SPELL_SHOUT:
		dam = dice(5, 6)
	case SPELL_CONE_OF_COLD:
		dam = dice(int(MIN(int64(level), 15)), 6)
	}
	if mag_newsaves(ch, victim, spellnum, level, int(ch.Aff_abils.Intel)) != 0 {
		if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<4) {
			send_to_char(victim, libc.CString("@g*save*@y You avoid any injury.@n\r\n"))
			dam = 0
		} else if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<3) {
			if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<1) && int(victim.Feats[FEAT_EVASION]) != 0 && ((victim.Equipment[WEAR_BODY]) == nil || int((victim.Equipment[WEAR_BODY]).Type_flag) != ITEM_ARMOR || ((victim.Equipment[WEAR_BODY]).Value[VAL_ARMOR_SKILL]) < ARMOR_TYPE_MEDIUM) {
				send_to_char(victim, libc.CString("@g*save*@y Your evasion ability allows you to avoid ANY injury.@n\r\n"))
				dam = 0
			} else {
				send_to_char(victim, libc.CString("@g*save*@y You take half damage.@n\r\n"))
				dam /= 2
			}
		}
	} else if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<3) && IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), 1<<1) && int(victim.Feats[FEAT_IMPROVED_EVASION]) != 0 && ((victim.Equipment[WEAR_BODY]) == nil || int((victim.Equipment[WEAR_BODY]).Type_flag) != ITEM_ARMOR || ((victim.Equipment[WEAR_BODY]).Value[VAL_ARMOR_SKILL]) < ARMOR_TYPE_MEDIUM) {
		send_to_char(victim, libc.CString("@r*save*@y Your improved evasion prevents full damage even on failure.@n\r\n"))
		dam /= 2
	}
	return 1
}
func mag_affects(level int, ch *char_data, victim *char_data, spellnum int) {
	var (
		af             [5]affected_type
		accum_affect   bool  = FALSE != 0
		accum_duration bool  = FALSE != 0
		to_vict        *byte = nil
		to_room        *byte = nil
		i              int
	)
	if victim == nil || ch == nil {
		return
	}
	for i = 0; i < MAX_SPELL_AFFECTS; i++ {
		af[i].Type = int16(spellnum)
		af[i].Bitvector = 0
		af[i].Modifier = 0
		af[i].Location = APPLY_NONE
	}
	if mag_newsaves(ch, victim, spellnum, level, int(ch.Aff_abils.Intel)) != 0 {
		if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), (1<<5)|1<<4) {
			send_to_char(victim, libc.CString("@g*save*@y You avoid any lasting affects.@n\r\n"))
			return
		}
	}
	switch spellnum {
	case SPELL_CHILL_TOUCH:
		af[0].Location = APPLY_STR
		af[0].Duration = 24
		af[0].Modifier = -1
		accum_duration = TRUE != 0
		to_vict = libc.CString("You feel your strength wither!")
	case SPELL_MAGE_ARMOR:
		af[0].Location = APPLY_AC
		af[0].Modifier = 40
		af[0].Duration = int16(GET_LEVEL(ch) * 1)
		accum_duration = FALSE != 0
		to_vict = libc.CString("You feel someone protecting you.")
	case SPELL_BLESS:
		af[0].Location = APPLY_ACCURACY
		af[0].Modifier = 2
		af[0].Duration = 6
		af[1].Location = APPLY_WILL
		af[1].Modifier = 1
		af[1].Duration = 6
		accum_duration = TRUE != 0
		to_vict = libc.CString("You feel righteous.")
	case SPELL_BLINDNESS:
		if MOB_FLAGGED(victim, MOB_NOBLIND) {
			send_to_char(ch, libc.CString("You fail.\r\n"))
			return
		}
		af[0].Location = APPLY_ACCURACY
		af[0].Modifier = -4
		af[0].Duration = 2
		af[0].Bitvector = AFF_BLIND
		af[1].Location = APPLY_AC
		af[1].Modifier = -4
		af[1].Duration = 2
		af[1].Bitvector = AFF_BLIND
		to_room = libc.CString("$n seems to be blinded!")
		to_vict = libc.CString("You have been blinded!")
	case SPELL_BANE:
		af[0].Location = APPLY_ACCURACY
		af[0].Duration = int16((level / 2) + 1)
		af[0].Modifier = -1
		af[0].Bitvector = AFF_CURSE
		af[1].Location = APPLY_WILL
		af[1].Duration = int16((level / 2) + 1)
		af[1].Modifier = -1
		af[1].Bitvector = AFF_CURSE
		accum_duration = TRUE != 0
		accum_affect = TRUE != 0
		to_room = libc.CString("$n briefly glows red!")
		to_vict = libc.CString("You feel very uncomfortable.")
	case SPELL_BESTOW_CURSE:
		af[0].Location = APPLY_STR
		af[0].Duration = -1
		af[0].Modifier = -6
		af[0].Bitvector = AFF_CURSE
		af[1].Location = APPLY_ACCURACY
		af[1].Duration = -1
		af[1].Modifier = -4
		af[1].Bitvector = AFF_CURSE
		accum_duration = FALSE != 0
		accum_affect = FALSE != 0
		to_room = libc.CString("$n briefly glows red!")
		to_vict = libc.CString("You feel very uncomfortable.")
	case SPELL_DETECT_ALIGN:
		af[0].Duration = int16(level + 12)
		af[0].Bitvector = AFF_DETECT_ALIGN
		accum_duration = TRUE != 0
		to_vict = libc.CString("Your eyes tingle.")
	case SPELL_SEE_INVIS:
		af[0].Duration = int16(level + 12)
		af[0].Bitvector = AFF_DETECT_INVIS
		accum_duration = TRUE != 0
		to_vict = libc.CString("Your eyes tingle.")
	case SPELL_DETECT_MAGIC:
		af[0].Duration = int16(level + 12)
		af[0].Bitvector = AFF_DETECT_MAGIC
		accum_duration = TRUE != 0
		to_vict = libc.CString("Your eyes tingle.")
	case SPELL_FAERIE_FIRE:
		af[0].Location = APPLY_AC
		af[0].Modifier = -1
		af[0].Duration = 3
		accum_duration = FALSE != 0
		to_vict = libc.CString("Your body flickers with a purplish light.")
		to_room = libc.CString("$n's body flickers with with a purplish light.")
	case SPELL_DARKVISION:
		af[0].Duration = int16(level + 12)
		af[0].Bitvector = AFF_INFRAVISION
		accum_duration = TRUE != 0
		to_vict = libc.CString("Your eyes glow red.")
		to_room = libc.CString("$n's eyes glow red.")
	case SPELL_INVISIBLE:
		if victim == nil {
			victim = ch
		}
		af[0].Duration = int16((level / 4) + 12)
		af[0].Modifier = 4
		af[0].Location = APPLY_AC
		af[0].Bitvector = AFF_INVISIBLE
		accum_duration = TRUE != 0
		to_vict = libc.CString("You vanish.")
		to_room = libc.CString("$n slowly fades out of existence.")
	case SPELL_POISON:
		af[0].Location = APPLY_STR
		af[0].Duration = int16(level)
		af[0].Modifier = -2
		af[0].Bitvector = AFF_POISON
		to_vict = libc.CString("You feel very sick.")
		to_room = libc.CString("$n gets violently ill!")
	case SPELL_PROT_FROM_EVIL:
		af[0].Duration = 24
		af[0].Bitvector = AFF_PROTECT_GOOD
		accum_duration = TRUE != 0
		to_vict = libc.CString("You feel invulnerable!")
	case SPELL_SANCTUARY:
		af[0].Duration = 4
		af[0].Bitvector = AFF_SANCTUARY
		accum_duration = TRUE != 0
		to_vict = libc.CString("A white aura momentarily surrounds you.")
		to_room = libc.CString("$n is surrounded by a white aura.")
	case SPELL_SLEEP:
		if config_info.Play.Pk_allowed == 0 && !IS_NPC(ch) && !IS_NPC(victim) {
			return
		}
		if MOB_FLAGGED(victim, MOB_NOSLEEP) {
			return
		}
		af[0].Duration = int16((level / 4) + 4)
		af[0].Bitvector = AFF_SLEEP
		if int(victim.Position) > POS_SLEEPING {
			send_to_char(victim, libc.CString("You feel very sleepy...  Zzzz......\r\n"))
			act(libc.CString("$n goes to sleep."), TRUE, victim, nil, nil, TO_ROOM)
			victim.Position = POS_SLEEPING
		}
	case SPELL_HAYASA:
		if config_info.Play.Pk_allowed == 0 && !IS_NPC(ch) && !IS_NPC(victim) {
			return
		}
		if MOB_FLAGGED(victim, MOB_NOSLEEP) {
			return
		}
		af[0].Duration = int16((level / 4) + 4)
		af[0].Bitvector = AFF_SLEEP
		if int(victim.Position) > POS_SLEEPING {
			send_to_char(victim, libc.CString("You feel very sleepy...  Zzzz......\r\n"))
			act(libc.CString("$n goes to sleep."), TRUE, victim, nil, nil, TO_ROOM)
			victim.Position = POS_SLEEPING
		}
	case SPELL_BULL_STRENGTH:
		af[0].Location = APPLY_STR
		af[0].Duration = int16(level)
		af[0].Modifier = rand_number(1, 4) + 1
		accum_duration = FALSE != 0
		accum_affect = FALSE != 0
		to_vict = libc.CString("You feel stronger!")
	case SPELL_SENSE_LIFE:
		to_vict = libc.CString("Your feel your awareness improve.")
		af[0].Duration = int16(level)
		af[0].Bitvector = AFF_SENSE_LIFE
		accum_duration = TRUE != 0
	case SPELL_WATERWALK:
		af[0].Duration = 24
		af[0].Bitvector = AFF_WATERWALK
		accum_duration = TRUE != 0
		to_vict = libc.CString("You feel webbing between your toes.")
	case SPELL_STONESKIN:
		af[0].Duration = int16(GET_LEVEL(ch) * 1)
		accum_duration = FALSE != 0
		to_vict = libc.CString("Your skin hardens into stone!")
	}
	if IS_NPC(victim) && !affected_by_spell(victim, spellnum) {
		for i = 0; i < MAX_SPELL_AFFECTS; i++ {
			if AFF_FLAGGED(victim, af[i].Bitvector) {
				send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
				return
			}
		}
	}
	if affected_by_spell(victim, spellnum) && (!accum_duration && !accum_affect) {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
		return
	}
	for i = 0; i < MAX_SPELL_AFFECTS; i++ {
		if af[i].Bitvector != 0 || af[i].Location != APPLY_NONE {
			affect_join(victim, &af[i], accum_duration, FALSE != 0, accum_affect, FALSE != 0)
		}
	}
	if to_vict != nil {
		act(to_vict, FALSE, victim, nil, unsafe.Pointer(ch), TO_CHAR)
	}
	if to_room != nil {
		act(to_room, TRUE, victim, nil, unsafe.Pointer(ch), TO_ROOM)
	}
}
func perform_mag_groups(level int, ch *char_data, tch *char_data, spellnum int) {
	switch spellnum {
	case SPELL_MASS_HEAL:
		mag_points(level, ch, tch, SPELL_HEAL)
	case SPELL_GROUP_ARMOR:
		mag_affects(level, ch, tch, SPELL_MAGE_ARMOR)
	case SPELL_GROUP_RECALL:
		spell_recall(level, ch, tch, nil, nil)
	}
}
func mag_groups(level int, ch *char_data, spellnum int) {
	var (
		tch    *char_data
		k      *char_data
		f      *follow_type
		f_next *follow_type
	)
	if ch == nil {
		return
	}
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		return
	}
	if ch.Master != nil {
		k = ch.Master
	} else {
		k = ch
	}
	for f = k.Followers; f != nil; f = f_next {
		f_next = f.Next
		tch = f.Follower
		if tch.In_room != ch.In_room {
			continue
		}
		if !AFF_FLAGGED(tch, AFF_GROUP) {
			continue
		}
		if ch == tch {
			continue
		}
		perform_mag_groups(level, ch, tch, spellnum)
	}
	if k != ch && AFF_FLAGGED(k, AFF_GROUP) {
		perform_mag_groups(level, ch, k, spellnum)
	}
	perform_mag_groups(level, ch, ch, spellnum)
}
func mag_masses(level int, ch *char_data, spellnum int) {
	var (
		tch      *char_data
		tch_next *char_data
	)
	for tch = world[ch.In_room].People; tch != nil; tch = tch_next {
		tch_next = tch.Next_in_room
		if tch == ch {
			continue
		}
		switch spellnum {
		}
	}
}
func mag_areas(level int, ch *char_data, spellnum int) {
	var (
		tch      *char_data
		next_tch *char_data
		to_char  *byte = nil
		to_room  *byte = nil
	)
	if ch == nil {
		return
	}
	switch spellnum {
	case SPELL_EARTHQUAKE:
		to_char = libc.CString("You gesture and the earth begins to shake all around you!")
		to_room = libc.CString("$n gracefully gestures and the earth begins to shake violently!")
	}
	if to_char != nil {
		act(to_char, FALSE, ch, nil, nil, TO_CHAR)
	}
	if to_room != nil {
		act(to_room, FALSE, ch, nil, nil, TO_ROOM)
	}
	for tch = world[ch.In_room].People; tch != nil; tch = next_tch {
		next_tch = tch.Next_in_room
		if tch == ch {
			continue
		}
		if ADM_FLAGGED(tch, ADM_NODAMAGE) {
			continue
		}
		if config_info.Play.Pk_allowed == 0 && !IS_NPC(ch) && !IS_NPC(tch) {
			continue
		}
		if !IS_NPC(ch) && IS_NPC(tch) && AFF_FLAGGED(tch, AFF_CHARM) {
			continue
		}
		mag_damage(level, ch, tch, spellnum)
	}
}

var monsum_list_lg_1 [4]mob_vnum = [4]mob_vnum{300, 301, 302, mob_vnum(-1)}
var monsum_list_ng_1 [6]mob_vnum = [6]mob_vnum{300, 301, 302, 303, 304, mob_vnum(-1)}
var monsum_list_cg_1 [4]mob_vnum = [4]mob_vnum{302, 303, 304, mob_vnum(-1)}
var monsum_list_ln_1 [5]mob_vnum = [5]mob_vnum{300, 301, 305, 306, mob_vnum(-1)}
var monsum_list_nn_1 [4]mob_vnum = [4]mob_vnum{302, 307, 308, mob_vnum(-1)}
var monsum_list_cn_1 [6]mob_vnum = [6]mob_vnum{303, 304, 309, 310, 311, mob_vnum(-1)}
var monsum_list_le_1 [5]mob_vnum = [5]mob_vnum{305, 306, 307, 308, mob_vnum(-1)}
var monsum_list_ne_1 [8]mob_vnum = [8]mob_vnum{305, 306, 307, 308, 309, 310, 311, mob_vnum(-1)}
var monsum_list_ce_1 [6]mob_vnum = [6]mob_vnum{307, 308, 309, 310, 311, mob_vnum(-1)}
var monsum_list_lg_2 [4]mob_vnum = [4]mob_vnum{312, 313, 314, mob_vnum(-1)}
var monsum_list_ng_2 [4]mob_vnum = [4]mob_vnum{312, 313, 314, mob_vnum(-1)}
var monsum_list_cg_2 [5]mob_vnum = [5]mob_vnum{312, 313, 314, 315, mob_vnum(-1)}
var monsum_list_ln_2 [3]mob_vnum = [3]mob_vnum{312, 317, mob_vnum(-1)}
var monsum_list_nn_2 [4]mob_vnum = [4]mob_vnum{315, 316, 317, mob_vnum(-1)}
var monsum_list_cn_2 [4]mob_vnum = [4]mob_vnum{313, 316, 318, mob_vnum(-1)}
var monsum_list_le_2 [3]mob_vnum = [3]mob_vnum{316, 317, mob_vnum(-1)}
var monsum_list_ne_2 [3]mob_vnum = [3]mob_vnum{318, 319, mob_vnum(-1)}
var monsum_list_ce_2 [3]mob_vnum = [3]mob_vnum{320, 321, mob_vnum(-1)}
var monsum_list_lg_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_3 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_4 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_5 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_6 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_7 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_8 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_lg_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ng_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cg_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ln_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_nn_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_cn_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_le_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ne_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list_ce_9 [1]mob_vnum = [1]mob_vnum{mob_vnum(-1)}
var monsum_list [9][9]*mob_vnum = [9][9]*mob_vnum{{&monsum_list_lg_1[0], &monsum_list_ng_1[0], &monsum_list_cg_1[0], &monsum_list_ln_1[0], &monsum_list_nn_1[0], &monsum_list_cn_1[0], &monsum_list_le_1[0], &monsum_list_ne_1[0], &monsum_list_ce_1[0]}, {&monsum_list_lg_2[0], &monsum_list_ng_2[0], &monsum_list_cg_2[0], &monsum_list_ln_2[0], &monsum_list_nn_2[0], &monsum_list_cn_2[0], &monsum_list_le_2[0], &monsum_list_ne_2[0], &monsum_list_ce_2[0]}, {&monsum_list_lg_3[0], &monsum_list_ng_3[0], &monsum_list_cg_3[0], &monsum_list_ln_3[0], &monsum_list_nn_3[0], &monsum_list_cn_3[0], &monsum_list_le_3[0], &monsum_list_ne_3[0], &monsum_list_ce_3[0]}, {&monsum_list_lg_4[0], &monsum_list_ng_4[0], &monsum_list_cg_4[0], &monsum_list_ln_4[0], &monsum_list_nn_4[0], &monsum_list_cn_4[0], &monsum_list_le_4[0], &monsum_list_ne_4[0], &monsum_list_ce_4[0]}, {&monsum_list_lg_5[0], &monsum_list_ng_5[0], &monsum_list_cg_5[0], &monsum_list_ln_5[0], &monsum_list_nn_5[0], &monsum_list_cn_5[0], &monsum_list_le_5[0], &monsum_list_ne_5[0], &monsum_list_ce_5[0]}, {&monsum_list_lg_6[0], &monsum_list_ng_6[0], &monsum_list_cg_6[0], &monsum_list_ln_6[0], &monsum_list_nn_6[0], &monsum_list_cn_6[0], &monsum_list_le_6[0], &monsum_list_ne_6[0], &monsum_list_ce_6[0]}, {&monsum_list_lg_7[0], &monsum_list_ng_7[0], &monsum_list_cg_7[0], &monsum_list_ln_7[0], &monsum_list_nn_7[0], &monsum_list_cn_7[0], &monsum_list_le_7[0], &monsum_list_ne_7[0], &monsum_list_ce_7[0]}, {&monsum_list_lg_8[0], &monsum_list_ng_8[0], &monsum_list_cg_8[0], &monsum_list_ln_8[0], &monsum_list_nn_8[0], &monsum_list_cn_8[0], &monsum_list_le_8[0], &monsum_list_ne_8[0], &monsum_list_ce_8[0]}, {&monsum_list_lg_9[0], &monsum_list_ng_9[0], &monsum_list_cg_9[0], &monsum_list_ln_9[0], &monsum_list_nn_9[0], &monsum_list_cn_9[0], &monsum_list_le_9[0], &monsum_list_ne_9[0], &monsum_list_ce_9[0]}}
var mag_summon_msgs [3]*byte = [3]*byte{libc.CString("\r\n"), libc.CString("$n animates a corpse!"), libc.CString("$n summons extraplanar assistance!")}

func mag_summons(level int, ch *char_data, obj *obj_data, spellnum int, arg *byte) {
	var (
		mob           *char_data = nil
		tobj          *obj_data
		next_obj      *obj_data
		msg           int = 0
		num           int = 1
		handle_corpse int = FALSE
		affs          int = 0
		affvs         int = 0
		assist        int = 0
		i             int
		j             int
		count         int
		buf           *byte = nil
		buf2          [2048]byte
		lev           int
		mob_num       mob_vnum
	)
	if ch == nil {
		return
	}
	lev = spell_info[spellnum].Spell_level
	switch spellnum {
	case SPELL_ANIMATE_DEAD:
		if obj == nil {
			send_to_char(ch, libc.CString("With what corpse?\r\n"))
			return
		}
		if !IS_CORPSE(obj) {
			send_to_char(ch, libc.CString("That's not a corpse!\r\n"))
			return
		}
		handle_corpse = TRUE
		msg = 11
		mob_num = MOB_ZOMBIE
	case SPELL_SUMMON_MONSTER_I:
		fallthrough
	case SPELL_SUMMON_MONSTER_II:
		fallthrough
	case SPELL_SUMMON_MONSTER_III:
		fallthrough
	case SPELL_SUMMON_MONSTER_IV:
		fallthrough
	case SPELL_SUMMON_MONSTER_V:
		fallthrough
	case SPELL_SUMMON_MONSTER_VI:
		fallthrough
	case SPELL_SUMMON_MONSTER_VII:
		fallthrough
	case SPELL_SUMMON_MONSTER_VIII:
		fallthrough
	case SPELL_SUMMON_MONSTER_IX:
		mob_num = -1
		affvs = 1
		assist = 1
		if arg != nil {
			buf = arg
			skip_spaces(&buf)
			if *buf == 0 {
				buf = nil
			}
		}
		j = int(ALIGN_TYPE(ch))
		if buf != nil {
			buf = any_one_arg(buf, &buf2[0])
			for i = lev - 1; i >= 0; i-- {
				for count = 0; *(*mob_vnum)(unsafe.Add(unsafe.Pointer(monsum_list[i][j]), unsafe.Sizeof(mob_vnum(0))*uintptr(count))) != mob_vnum(-1); count++ {
					mob_num = *(*mob_vnum)(unsafe.Add(unsafe.Pointer(monsum_list[i][j]), unsafe.Sizeof(mob_vnum(0))*uintptr(count)))
					if real_mobile(mob_num) == mob_rnum(-1) {
						mob_num = -1
					} else if is_name(&buf2[0], mob_proto[real_mobile(mob_num)].Name) == 0 {
						mob_num = -1
					} else {
						break
					}
				}
				if mob_num != mob_vnum(-1) {
					break
				}
			}
			if mob_num == mob_vnum(-1) {
				send_to_char(ch, libc.CString("That's not a name for a monster you can summon. Summoning something else.\r\n"))
			} else {
				basic_mud_log(libc.CString("lev=%d, i=%d, ngen=%d"), lev, i, lev-i)
				switch lev - i {
				case 1:
					num = 1
				case 2:
					num = rand_number(1, 3)
				default:
					num = rand_number(1, 4) + 1
				}
			}
		}
		if mob_num == mob_vnum(-1) {
			num = 1
			for count = 0; *(*mob_vnum)(unsafe.Add(unsafe.Pointer(monsum_list[lev-1][j]), unsafe.Sizeof(mob_vnum(0))*uintptr(count))) != mob_vnum(-1); count++ {
			}
			if count == 0 {
				basic_mud_log(libc.CString("No monsums for spell level %d align %s"), lev, alignments[j])
				return
			}
			count--
			mob_num = *(*mob_vnum)(unsafe.Add(unsafe.Pointer(monsum_list[lev-1][j]), unsafe.Sizeof(mob_vnum(0))*uintptr(rand_number(0, count))))
		}
	default:
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) {
		send_to_char(ch, libc.CString("You are too giddy to have any followers!\r\n"))
		return
	}
	for i = 0; i < num; i++ {
		if (func() *char_data {
			mob = read_mobile(mob_num, VIRTUAL)
			return mob
		}()) == nil {
			send_to_char(ch, libc.CString("You don't quite remember how to summon that creature.\r\n"))
			return
		}
		char_to_room(mob, ch.In_room)
		if affs != 0 {
			mag_affects(level, ch, mob, spellnum)
		}
		if affvs != 0 {
			mag_affectsv(level, ch, mob, spellnum)
		}
		mob.Carry_weight = 0
		mob.Carry_items = 0
		SET_BIT_AR(mob.Affected_by[:], AFF_CHARM)
		act(mag_summon_msgs[msg], FALSE, ch, nil, unsafe.Pointer(mob), TO_ROOM)
		load_mtrigger(mob)
		add_follower(mob, ch)
		if assist != 0 && ch.Fighting != nil {
			set_fighting(mob, ch.Fighting)
		}
		mob.Master_id = ch.Idnum
	}
	if handle_corpse != 0 {
		for tobj = obj.Contains; tobj != nil; tobj = next_obj {
			next_obj = tobj.Next_content
			obj_from_obj(tobj)
			obj_to_char(tobj, mob)
		}
		extract_obj(obj)
	}
}
func mag_points(level int, ch *char_data, victim *char_data, spellnum int) {
	var (
		healing int = 0
		tmp     int
	)
	if victim == nil {
		return
	}
	switch spellnum {
	case SPELL_CURE_LIGHT:
		healing = dice(1, 8) + int(MIN(int64(level), 5))
		send_to_char(victim, libc.CString("You feel better.\r\n"))
	case SPELL_CURE_CRITIC:
		healing = dice(4, 8) + int(MIN(int64(level), 20))
		send_to_char(victim, libc.CString("You feel a lot better!\r\n"))
	case SPELL_HEAL:
		healing = dice(3, 8) + 100
		send_to_char(victim, libc.CString("A warm feeling floods your body.\r\n"))
		if AFF_FLAGGED(ch, AFF_CDEATH) {
			affectv_from_char(ch, ART_QUIVERING_PALM)
			send_to_char(ch, libc.CString("Your nerves settle slightly\r\n"))
		}
	case SPELL_SENSU:
		if int(victim.Player_specials.Conditions[HUNGER]) > -1 {
			victim.Player_specials.Conditions[HUNGER] = 48
		}
		if victim.Kaioken <= 0 {
			victim.Hit = victim.Max_hit
		}
		victim.Mana = victim.Max_mana
		victim.Move = victim.Max_move
		if AFF_FLAGGED(victim, AFF_KNOCKED) {
			act(libc.CString("@W$n@W is no longer senseless, and wakes up.@n"), FALSE, ch, nil, nil, TO_ROOM)
			send_to_char(victim, libc.CString("You are no longer knocked out, and wake up!@n\r\n"))
			REMOVE_BIT_AR(victim.Affected_by[:], AFF_KNOCKED)
			victim.Position = POS_SITTING
		}
		if victim.Suppression > 0 && victim.Hit > ((victim.Max_hit/100)*victim.Suppression) {
			victim.Hit = (victim.Max_hit / 100) * victim.Suppression
			send_to_char(victim, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
		}
		send_to_char(victim, libc.CString("@GYour wounds heal and your strength returns.@n\r\n"))
		act(libc.CString("@C$n@W suddenly looks a lot better!@b"), FALSE, victim, nil, nil, TO_NOTVICT)
		affect_from_char(victim, SPELL_POISON)
		if AFF_FLAGGED(victim, AFF_BURNED) {
			send_to_char(victim, libc.CString("Your burns are healed now.\r\n"))
			act(libc.CString("$n@w's burns are now healed.@n"), TRUE, victim, nil, nil, TO_ROOM)
			REMOVE_BIT_AR(victim.Affected_by[:], AFF_BURNED)
		}
		if (victim.Limb_condition[0]) <= 0 {
			send_to_char(victim, libc.CString("Your right arm grows back!\r\n"))
			victim.Limb_condition[0] = 100
		} else if (victim.Limb_condition[0]) < 50 {
			send_to_char(victim, libc.CString("Your right arm is no longer broken!\r\n"))
			victim.Limb_condition[0] = 100
		}
		if (victim.Limb_condition[1]) <= 0 {
			send_to_char(victim, libc.CString("Your left arm grows back!\r\n"))
			victim.Limb_condition[1] = 100
		} else if (victim.Limb_condition[1]) < 50 {
			send_to_char(victim, libc.CString("Your left arm is no longer broken!\r\n"))
			victim.Limb_condition[1] = 100
		}
		if (victim.Limb_condition[2]) <= 0 {
			send_to_char(victim, libc.CString("Your right leg grows back!\r\n"))
			victim.Limb_condition[2] = 100
		} else if (victim.Limb_condition[2]) < 50 {
			send_to_char(victim, libc.CString("Your right leg is no longer broken!\r\n"))
			victim.Limb_condition[2] = 100
		}
		if (victim.Limb_condition[3]) <= 0 {
			send_to_char(victim, libc.CString("Your left leg grows back!\r\n"))
			victim.Limb_condition[3] = 100
		} else if (victim.Limb_condition[3]) < 50 {
			send_to_char(victim, libc.CString("Your left leg is no longer broken!\r\n"))
			victim.Limb_condition[3] = 100
		}
	case ART_WHOLENESS_OF_BODY:
		healing = int(victim.Max_hit - victim.Hit)
		healing = int(MAX(0, int64(healing)))
		tmp = int(ch.Ki / 2)
		if tmp > healing {
			tmp = healing
		} else {
			healing = tmp
		}
		ch.Ki -= int64(tmp * 2)
	}
	update_pos(victim)
}
func mag_unaffects(level int, ch *char_data, victim *char_data, spellnum int) {
	var (
		spell            int   = 0
		msg_not_affected int   = TRUE
		to_vict          *byte = nil
		to_room          *byte = nil
	)
	if victim == nil {
		return
	}
	switch spellnum {
	case SPELL_HEAL:
		msg_not_affected = FALSE
		fallthrough
	case SPELL_REMOVE_BLINDNESS:
		spell = SPELL_BLINDNESS
		to_vict = libc.CString("Your vision returns!")
		to_room = libc.CString("There's a momentary gleam in $n's eyes.")
	case SPELL_NEUTRALIZE_POISON:
		spell = SPELL_POISON
		to_vict = libc.CString("A warm feeling runs through your body!")
		to_room = libc.CString("$n looks better.")
	case SPELL_REMOVE_CURSE:
		spell = SPELL_BESTOW_CURSE
		to_vict = libc.CString("You don't feel so unlucky.")
	default:
		return
	}
	if !affected_by_spell(victim, spell) {
		if msg_not_affected != 0 {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
		}
		return
	}
	affect_from_char(victim, spell)
	if to_vict != nil {
		act(to_vict, FALSE, victim, nil, unsafe.Pointer(ch), TO_CHAR)
	}
	if to_room != nil {
		act(to_room, TRUE, victim, nil, unsafe.Pointer(ch), TO_ROOM)
	}
}
func mag_alter_objs(level int, ch *char_data, obj *obj_data, spellnum int) {
	var (
		to_char *byte = nil
		to_room *byte = nil
	)
	if obj == nil {
		return
	}
	switch spellnum {
	case SPELL_BLESS:
		if !OBJ_FLAGGED(obj, ITEM_BLESS) && obj.Weight <= int64(level*5) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_BLESS)
			to_char = libc.CString("$p glows briefly.")
		}
	case SPELL_INVISIBLE:
		if !OBJ_FLAGGED(obj, bitvector_t(int32(int(ITEM_NOINVIS|ITEM_INVISIBLE)))) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_INVISIBLE)
			to_char = libc.CString("$p vanishes.")
		}
	case SPELL_POISON:
		if (int(obj.Type_flag) == ITEM_DRINKCON || int(obj.Type_flag) == ITEM_FOUNTAIN || int(obj.Type_flag) == ITEM_FOOD) && (obj.Value[VAL_FOOD_POISON]) == 0 {
			obj.Value[VAL_FOOD_POISON] = 1
			to_char = libc.CString("$p steams briefly.")
		}
	case SPELL_REMOVE_CURSE:
		if OBJ_FLAGGED(obj, ITEM_NODROP) {
			REMOVE_BIT_AR(obj.Extra_flags[:], ITEM_NODROP)
			if int(obj.Type_flag) == ITEM_WEAPON {
				(obj.Value[VAL_WEAPON_DAMSIZE])++
			}
			to_char = libc.CString("$p briefly glows blue.")
		}
	case SPELL_NEUTRALIZE_POISON:
		if (int(obj.Type_flag) == ITEM_DRINKCON || int(obj.Type_flag) == ITEM_FOUNTAIN || int(obj.Type_flag) == ITEM_FOOD) && (obj.Value[VAL_FOOD_POISON]) != 0 {
			obj.Value[VAL_FOOD_POISON] = 0
			to_char = libc.CString("$p steams briefly.")
		}
	}
	if to_char == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
	} else {
		act(to_char, TRUE, ch, obj, nil, TO_CHAR)
	}
	if to_room != nil {
		act(to_room, TRUE, ch, obj, nil, TO_ROOM)
	} else if to_char != nil {
		act(to_char, TRUE, ch, obj, nil, TO_ROOM)
	}
}
func mag_creations(level int, ch *char_data, spellnum int) {
	var (
		tobj *obj_data
		z    obj_vnum
	)
	if ch == nil {
		return
	}
	switch spellnum {
	case SPELL_CREATE_FOOD:
		z = 10
	default:
		send_to_char(ch, libc.CString("Spell unimplemented, it would seem.\r\n"))
		return
	}
	if (func() *obj_data {
		tobj = read_object(z, VIRTUAL)
		return tobj
	}()) == nil {
		send_to_char(ch, libc.CString("I seem to have goofed.\r\n"))
		basic_mud_log(libc.CString("SYSERR: spell_creations, spell %d, obj %d: obj not found"), spellnum, z)
		return
	}
	add_unique_id(tobj)
	obj_to_char(tobj, ch)
	act(libc.CString("$n creates $p."), FALSE, ch, tobj, nil, TO_ROOM)
	act(libc.CString("You create $p."), FALSE, ch, tobj, nil, TO_CHAR)
	load_otrigger(tobj)
}
func affect_update_violence() {
	var (
		af     *affected_type
		next   *affected_type
		i      *char_data
		dam    int
		maxdam int
	)
	for i = affectv_list; i != nil; i = i.Next_affectv {
		for af = i.Affectedv; af != nil; af = next {
			next = af.Next
			if int(af.Duration) >= 1 {
				af.Duration--
				switch af.Type {
				case ART_EMPTY_BODY:
					if i.Ki >= 10 {
						i.Ki -= 10
					} else {
						af.Duration = 0
					}
				}
			} else if int(af.Duration) == -1 {
				continue
			}
			if int(af.Duration) == 0 {
				if int(af.Type) > 0 && int(af.Type) < SKILL_TABLE_SIZE {
					if af.Next == nil || int(af.Next.Type) != int(af.Type) || int(af.Next.Duration) > 0 {
						if spell_info[af.Type].Wear_off_msg != nil {
							send_to_char(i, libc.CString("%s\r\n"), spell_info[af.Type].Wear_off_msg)
						}
					}
				}
				if af.Bitvector == AFF_SUMMONED {
					stop_follower(i)
					if !PLR_FLAGGED(i, PLR_NOTDEADYET) && !MOB_FLAGGED(i, MOB_NOTDEADYET) {
						extract_char(i)
					}
				}
				if int(af.Type) == ART_QUIVERING_PALM {
					maxdam = int(i.Hit + 8)
					dam = int(i.Max_hit * 3 / 4)
					dam = int(MIN(int64(dam), int64(maxdam)))
					dam = int(MAX(0, int64(dam)))
					basic_mud_log(libc.CString("Creeping death strike doing %d dam"), dam)
				}
				affectv_remove(i, af)
			}
		}
	}
}
func mag_affectsv(level int, ch *char_data, victim *char_data, spellnum int) {
	var (
		af             [5]affected_type
		accum_affect   bool  = FALSE != 0
		accum_duration bool  = FALSE != 0
		to_vict        *byte = nil
		to_room        *byte = nil
		i              int
	)
	if victim == nil || ch == nil {
		return
	}
	for i = 0; i < MAX_SPELL_AFFECTS; i++ {
		af[i].Type = int16(spellnum)
		af[i].Bitvector = 0
		af[i].Modifier = 0
		af[i].Location = APPLY_NONE
	}
	if mag_newsaves(ch, victim, spellnum, level, int(ch.Aff_abils.Intel)) != 0 {
		if IS_SET(bitvector_t(int32(spell_info[spellnum].Save_flags)), (1<<5)|1<<4) {
			send_to_char(victim, libc.CString("@g*save*@y You avoid any lasting affects.@n\r\n"))
			return
		}
	}
	switch spellnum {
	case SPELL_PARALYZE:
		af[0].Duration = int16(level / 2)
		af[0].Bitvector = AFF_PARALYZE
		accum_duration = FALSE != 0
		to_vict = libc.CString("You feel your limbs freeze!")
		to_room = libc.CString("$n suddenly freezes in place!")
	case ART_STUNNING_FIST:
		af[0].Duration = 1
		af[0].Bitvector = AFF_STUNNED
		accum_duration = FALSE != 0
		to_vict = libc.CString("You are in a stunned daze!")
		to_room = libc.CString("$n is stunned.")
	case ART_EMPTY_BODY:
		af[0].Duration = int16(ch.Ki / 10)
		af[0].Bitvector = AFF_ETHEREAL
		accum_duration = FALSE != 0
		to_vict = libc.CString("You switch to the ethereal plane.")
		to_room = libc.CString("$n disappears.")
	case ART_QUIVERING_PALM:
		if GET_LEVEL(ch) <= GET_LEVEL(victim) {
			send_to_char(ch, libc.CString("They are too high level for that.\r\n"))
			return
		}
		af[0].Duration = int16(MAX(6, int64(20-level)))
		af[0].Bitvector = AFF_CDEATH
		accum_duration = FALSE != 0
		to_vict = libc.CString("You feel death closing in.")
	case SPELL_RESISTANCE:
		af[0].Duration = 12
		af[0].Location = APPLY_ALLSAVES
		af[0].Modifier = 1
		accum_duration = FALSE != 0
		to_vict = libc.CString("You glow briefly with a silvery light.")
		to_room = libc.CString("$n glows briefly with a silvery light.")
	case SPELL_DAZE:
		if victim.Race_level < 5 {
			af[0].Bitvector = AFF_NEXTNOACTION
		}
		accum_duration = FALSE != 0
		to_vict = libc.CString("You are struck dumb by a flash of bright light!")
		to_room = libc.CString("$n is struck dumb by a flash of bright light!")
	case SPELL_SUMMON_MONSTER_I:
		fallthrough
	case SPELL_SUMMON_MONSTER_II:
		fallthrough
	case SPELL_SUMMON_MONSTER_III:
		fallthrough
	case SPELL_SUMMON_MONSTER_IV:
		fallthrough
	case SPELL_SUMMON_MONSTER_V:
		fallthrough
	case SPELL_SUMMON_MONSTER_VI:
		fallthrough
	case SPELL_SUMMON_MONSTER_VII:
		fallthrough
	case SPELL_SUMMON_MONSTER_VIII:
		fallthrough
	case SPELL_SUMMON_MONSTER_IX:
		af[0].Duration = int16(level)
		af[0].Bitvector = AFF_SUMMONED
		accum_duration = FALSE != 0
		to_vict = libc.CString("You are summoned to assist $N!")
		to_room = libc.CString("$n appears, ready for action.")
	case SPELL_FLARE:
		if MOB_FLAGGED(victim, MOB_NOBLIND) {
			send_to_char(ch, libc.CString("You fail.\r\n"))
			return
		}
		af[0].Location = APPLY_ACCURACY
		af[0].Modifier = -1
		af[0].Duration = 2
		af[0].Bitvector = AFF_BLIND
		to_room = libc.CString("$n seems to be dazzled!")
		to_vict = libc.CString("You have been dazzled!")
	case SPELL_FIRE_SHIELD:
		af[0].Duration = int16(level * 5)
		af[0].Bitvector = AFF_FIRE_SHIELD
		to_room = libc.CString("$n is engulfed in a firey shield!")
		to_vict = libc.CString("You are engulfed in a firey shield!")
	}
	if IS_NPC(victim) && !affected_by_spell(victim, spellnum) {
		for i = 0; i < MAX_SPELL_AFFECTS; i++ {
			if AFF_FLAGGED(victim, af[i].Bitvector) {
				send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
				return
			}
		}
	}
	if affected_by_spell(victim, spellnum) && (!accum_duration && !accum_affect) {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOEFFECT)
		return
	}
	for i = 0; i < MAX_SPELL_AFFECTS; i++ {
		if af[i].Bitvector != 0 || af[i].Location != APPLY_NONE {
			affectv_join(victim, &af[i], accum_duration, FALSE != 0, accum_affect, FALSE != 0)
		}
	}
	if to_vict != nil {
		act(to_vict, FALSE, victim, nil, unsafe.Pointer(ch), TO_CHAR)
	}
	if to_room != nil {
		act(to_room, TRUE, victim, nil, unsafe.Pointer(ch), TO_ROOM)
	}
}
