package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const sick_fail = 2

func barrier_shed(ch *char_data) {
	if !AFF_FLAGGED(ch, AFF_SANCTUARY) {
		return
	}
	if GET_SKILL(ch, SKILL_AQUA_BARRIER) > 0 {
		return
	}
	var chance int = axion_dice(0)
	var barrier int = GET_SKILL(ch, SKILL_BARRIER)
	var concentrate int = GET_SKILL(ch, SKILL_CONCENTRATION)
	var rate float64 = 0.3
	if barrier >= 100 {
		rate = 0.01
	} else if barrier >= 95 {
		rate = 0.02
	} else if barrier >= 90 {
		rate = 0.04
	} else if barrier >= 80 {
		rate = 0.08
	} else if barrier >= 70 {
		rate = 0.1
	} else if barrier >= 60 {
		rate = 0.15
	} else if barrier >= 50 {
		rate = 0.2
	} else if barrier >= 40 {
		rate = 0.25
	} else if barrier >= 30 {
		rate = 0.27
	} else if barrier >= 20 {
		rate = 0.29
	}
	var loss int64 = int64(float64(ch.Barrier) * rate)
	var recharge int64 = 0
	if concentrate >= chance {
		recharge = int64(float64(loss) * 0.5)
	}
	ch.Barrier -= loss
	if ch.Barrier <= 0 {
		ch.Barrier = 0
		act(libc.CString("@cYour barrier disappears.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@c's barrier disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
	} else {
		act(libc.CString("@cYour barrier loses some energy.@n"), TRUE, ch, nil, nil, TO_CHAR)
		send_to_char(ch, libc.CString("@D[@C%s@D]@n\r\n"), add_commas(loss))
		act(libc.CString("@c$n@c's barrier sends some sparks into the air as it seems to get a bit weaker.@n"), TRUE, ch, nil, nil, TO_ROOM)
	}
	if recharge > 0 && ch.Mana < ch.Max_mana {
		ch.Mana += recharge
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		send_to_char(ch, libc.CString("@CYou reabsorb some of the energy lost into your body!@n\r\n"))
	}
}
func healthy_check(ch *char_data) {
	if (ch.Bonuses[BONUS_HEALTHY]) == 0 || int(ch.Position) != POS_SLEEPING {
		return
	}
	var chance int = 70
	var roll int = rand_number(1, 100)
	var change int = FALSE
	if AFF_FLAGGED(ch, AFF_SHOCKED) && roll >= chance {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_SHOCKED)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_MBREAK) && roll >= chance {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_MBREAK)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_WITHER) && roll >= chance {
		ch.Real_abils.Str += 3
		ch.Real_abils.Cha += 3
		save_char(ch)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_WITHER)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_CURSE) && roll >= chance {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_CURSE)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_POISON) && roll >= chance {
		null_affect(ch, AFF_POISON)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_PARALYZE) && roll >= chance {
		null_affect(ch, AFF_PARALYZE)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_PARA) && roll >= chance {
		null_affect(ch, AFF_PARA)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_BLIND) && roll >= chance {
		null_affect(ch, AFF_BLIND)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_HYDROZAP) && roll >= chance {
		ch.Real_abils.Dex += 4
		ch.Real_abils.Con += 4
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_HYDROZAP)
		save_char(ch)
		change = TRUE
	}
	if AFF_FLAGGED(ch, AFF_KNOCKED) && roll >= chance {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_KNOCKED)
		ch.Position = POS_SITTING
		change = TRUE
	}
	if change == TRUE {
		send_to_char(ch, libc.CString("@CYou feel your body recover from all its ailments!@n\r\n"))
	}
	return
}
func wearing_stardust(ch *char_data) int {
	var (
		count int = 0
		i     int
	)
	for i = 1; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			var obj *obj_data = (ch.Equipment[i])
			switch GET_OBJ_VNUM(obj) {
			case 1110:
				fallthrough
			case 1111:
				fallthrough
			case 1112:
				fallthrough
			case 1113:
				fallthrough
			case 1114:
				fallthrough
			case 1115:
				fallthrough
			case 1116:
				fallthrough
			case 1117:
				fallthrough
			case 1118:
				fallthrough
			case 1119:
				count += 1
			}
		}
	}
	if count == 26 {
		return 1
	} else {
		return 0
	}
}
func mana_gain(ch *char_data) int64 {
	var gain int64 = 0
	if IS_NPC(ch) {
		gain = ch.Max_mana / 70
	} else {
		if ROOM_FLAGGED(ch.In_room, ROOM_REGEN) || (ch.Bonuses[BONUS_DESTROYER]) > 0 && world[ch.In_room].Dmg >= 75 {
			if int(ch.Race) == RACE_KONATSU {
				gain = ch.Max_mana / 12
			}
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_mana / 11
			}
			if int(ch.Race) == RACE_ARLIAN {
				gain = ch.Max_mana / 30
			}
			if int(ch.Race) != RACE_KONATSU && int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_mana / 10
			}
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_REGEN) {
			if int(ch.Race) == RACE_KONATSU {
				gain = ch.Max_mana / 15
			}
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_mana / 13
			}
			if int(ch.Race) != RACE_KONATSU && int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_mana / 12
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_BEDROOM) {
				gain += int64(float64(gain) * 0.25)
			}
			if int(ch.Race) == RACE_ARLIAN {
				gain = ch.Max_mana / 40
			}
		}
		switch ch.Position {
		case POS_STANDING:
			if int(ch.Race) != RACE_HOSHIJIN || int(ch.Race) == RACE_HOSHIJIN && ch.Starphase <= 0 {
				gain = gain / 4
			} else {
				gain += gain / 2
			}
		case POS_FIGHTING:
			gain = gain / 4
		case POS_SLEEPING:
			if ch.Sits == nil {
				gain *= 2
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 {
				gain *= 3
				gain += int64(float64(gain) * 0.1)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A94 {
				gain *= 3
				gain += int64(float64(gain) * 0.3)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain *= 3
			}
		case POS_RESTING:
			if ch.Sits == nil {
				gain += gain / 2
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain *= 2
				gain += int64(float64(gain) * 0.1)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A94 && int(ch.Race) != RACE_ARLIAN {
				gain *= 2
				gain += int64(float64(gain) * 0.3)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain *= 2
			}
		case POS_SITTING:
			if ch.Sits == nil {
				gain += gain / 4
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 {
				gain += int64(float64(gain) * 0.6)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A94 {
				gain += int64(float64(gain) * 0.8)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain += int64(float64(gain) * 0.5)
			}
		}
	}
	if ch.In_room != room_rnum(-1) {
		if cook_element(ch.In_room) == 1 {
			gain += int64(float64(gain) * 0.2)
		}
	}
	if int(ch.Race) == RACE_ARLIAN && int(ch.Sex) == SEX_FEMALE && OUTSIDE(ch) {
		gain *= 4
	}
	if int(ch.Race) == RACE_KANASSAN && weather_info.Sky == SKY_RAINING && OUTSIDE(ch) {
		gain += int64(float64(gain) * 0.1)
	}
	if int(ch.Race) == RACE_KANASSAN && SUNKEN(ch.In_room) {
		gain *= 16
	}
	if int(ch.Race) == RACE_HOSHIJIN && ch.Starphase > 0 {
		gain *= 2
	}
	if PLR_FLAGGED(ch, PLR_HEALT) && ch.Sits != nil {
		gain *= 20
	}
	if PLR_FLAGGED(ch, PLR_POSE) && axion_dice(0) > GET_SKILL(ch, SKILL_POSE) {
		REMOVE_BIT_AR(ch.Act[:], PLR_POSE)
		send_to_char(ch, libc.CString("You feel slightly less confident now.\r\n"))
		ch.Real_abils.Str -= 8
		ch.Real_abils.Dex -= 8
		save_char(ch)
	}
	if AFF_FLAGGED(ch, AFF_HYDROZAP) && rand_number(1, 4) >= 4 {
		ch.Real_abils.Dex += 4
		ch.Real_abils.Con += 4
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_HYDROZAP)
		save_char(ch)
	}
	if GET_SKILL(ch, SKILL_CONCENTRATION) >= 100 {
		gain += gain / 2
	} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 75 {
		gain += gain / 4
	} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 50 {
		gain += gain / 6
	} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 25 {
		gain += gain / 8
	} else if GET_SKILL(ch, SKILL_CONCENTRATION) < 25 && GET_SKILL(ch, SKILL_CONCENTRATION) > 0 {
		gain += gain / 10
	}
	if AFF_FLAGGED(ch, AFF_BLESS) {
		gain *= 2
	}
	if AFF_FLAGGED(ch, AFF_CURSE) {
		gain /= 5
	}
	if ch.Foodr > 0 && rand_number(1, 2) == 2 {
		ch.Foodr -= 1
	}
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_HINTS) && rand_number(1, 5) == 5 {
		hint_system(ch, 0)
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		gain /= 4
	}
	if cook_element(ch.In_room) == 1 {
		gain *= 2
	}
	return gain
}
func hit_gain(ch *char_data) int64 {
	var gain int64 = 0
	if IS_NPC(ch) {
		gain = ch.Max_hit / 70
	} else {
		if ROOM_FLAGGED(ch.In_room, ROOM_REGEN) || (ch.Bonuses[BONUS_DESTROYER]) > 0 && world[ch.In_room].Dmg >= 75 {
			if int(ch.Race) == RACE_HUMAN {
				gain = ch.Max_hit / 20
			}
			if int(ch.Race) == RACE_ARLIAN {
				gain = ch.Max_hit / 30
			}
			if int(ch.Race) == RACE_NAMEK {
				gain = ch.Max_hit / 2
			}
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_hit / 11
			}
			if int(ch.Race) != RACE_HUMAN && int(ch.Race) != RACE_NAMEK && int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_hit / 10
			}
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_REGEN) {
			if int(ch.Race) == RACE_HUMAN {
				gain = ch.Max_hit / 30
			}
			if int(ch.Race) == RACE_NAMEK {
				gain = ch.Max_hit / 4
			}
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_hit / 16
			}
			if int(ch.Race) == RACE_ARLIAN {
				gain = ch.Max_hit / 40
			}
			if int(ch.Race) != RACE_HUMAN && int(ch.Race) != RACE_NAMEK && int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_hit / 15
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_BEDROOM) {
				gain += int64(float64(gain) * 0.25)
			}
		}
		switch ch.Position {
		case POS_STANDING:
			if int(ch.Race) != RACE_HOSHIJIN || int(ch.Race) == RACE_HOSHIJIN && ch.Starphase <= 0 {
				gain = gain / 4
			} else if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) {
				gain = gain / 3
			} else {
				gain += gain / 2
			}
		case POS_FIGHTING:
			gain = gain / 4
		case POS_SLEEPING:
			if int(ch.Race) == RACE_ARLIAN {
				gain *= 3
			} else if ch.Sits == nil {
				gain *= 2
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 {
				gain *= 3
				gain += int64(float64(gain) * 0.1)
			} else if ch.Sits != nil {
				gain *= 3
			}
		case POS_RESTING:
			if ch.Sits == nil {
				gain += gain / 2
			} else if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) {
				gain = int64(float64(gain) * 1.5)
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 1.1)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain *= 2
			}
		case POS_SITTING:
			if ch.Sits == nil {
				gain += gain / 4
			} else if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) {
				gain = int64(float64(gain) * 0.5)
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 0.6)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain += int64(float64(gain) * 0.5)
			}
		}
	}
	healthy_check(ch)
	if int(ch.Race) == RACE_ARLIAN && int(ch.Sex) == SEX_FEMALE && OUTSIDE(ch) {
		gain *= 4
	}
	if int(ch.Race) == RACE_KANASSAN && weather_info.Sky == SKY_RAINING && OUTSIDE(ch) {
		gain += int64(float64(gain) * 0.1)
	}
	if int(ch.Race) == RACE_KANASSAN && SUNKEN(ch.In_room) {
		gain *= 16
	}
	if int(ch.Race) == RACE_HOSHIJIN && ch.Starphase > 0 {
		gain *= 2
	}
	if PLR_FLAGGED(ch, PLR_HEALT) && ch.Sits != nil {
		gain *= 20
	}
	if AFF_FLAGGED(ch, AFF_BLESS) {
		gain *= 2
	}
	if AFF_FLAGGED(ch, AFF_CURSE) {
		gain /= 5
	}
	if PLR_FLAGGED(ch, PLR_FURY) {
		send_to_char(ch, libc.CString("Your fury subsides for now. Next time try to take advantage of it before you calm down.\r\n"))
		REMOVE_BIT_AR(ch.Act[:], PLR_FURY)
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		gain /= 4
	}
	if cook_element(ch.In_room) == 1 {
		gain *= 2
	}
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_ABSORB) {
			gain = gain / 8
		}
	}
	if ch.Regen > 0 {
		gain += int64((float64(gain) * 0.01) * float64(ch.Regen))
	}
	return gain
}
func move_gain(ch *char_data) int64 {
	var gain int64 = 0
	if IS_NPC(ch) {
		gain = ch.Max_move / 70
	} else {
		if ROOM_FLAGGED(ch.In_room, ROOM_REGEN) || (ch.Bonuses[BONUS_DESTROYER]) > 0 && world[ch.In_room].Dmg >= 75 {
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_move / 7
			}
			if int(ch.Race) == RACE_ARLIAN {
				gain = ch.Max_move / 4
			}
			if int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_move / 6
			}
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_REGEN) {
			if int(ch.Race) == RACE_MUTANT {
				gain = ch.Max_move / 9
			}
			if int(ch.Race) != RACE_MUTANT {
				gain = ch.Max_move / 8
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_BEDROOM) {
				gain += int64(float64(gain) * 0.25)
			}
		}
		switch ch.Position {
		case POS_STANDING:
			if int(ch.Race) != RACE_HOSHIJIN || int(ch.Race) == RACE_HOSHIJIN && ch.Starphase <= 0 {
				gain = gain / 4
			} else {
				gain += gain / 2
			}
		case POS_FIGHTING:
			gain = gain / 4
		case POS_SLEEPING:
			if ch.Sits == nil {
				gain *= 2
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain *= 3
				gain += int64(float64(gain) * 0.1)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A93 && int(ch.Race) != RACE_ARLIAN {
				gain *= 3
				gain += int64(float64(gain) * 0.3)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain *= 3
			}
		case POS_RESTING:
			if ch.Sits == nil {
				gain += gain / 2
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 1.1)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A93 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 1.3)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain += gain
			}
		case POS_SITTING:
			if ch.Sits == nil {
				gain += gain / 4
			} else if GET_OBJ_VNUM(ch.Sits) == 19090 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 0.6)
			} else if GET_OBJ_VNUM(ch.Sits) == 0x4A93 && int(ch.Race) != RACE_ARLIAN {
				gain += int64(float64(gain) * 0.8)
			} else if ch.Sits != nil || int(ch.Race) == RACE_ARLIAN {
				gain += gain / 2
			}
		}
	}
	if int(ch.Race) == RACE_ARLIAN && int(ch.Sex) == SEX_FEMALE && OUTSIDE(ch) {
		gain *= 2
	}
	if int(ch.Race) == RACE_NAMEK {
		gain = int64(float64(gain) * 0.5)
	}
	if int(ch.Race) == RACE_KANASSAN && weather_info.Sky == SKY_RAINING && OUTSIDE(ch) {
		gain += int64(float64(gain) * 0.1)
	}
	if int(ch.Race) == RACE_KANASSAN && SUNKEN(ch.In_room) {
		gain *= 16
	}
	if int(ch.Race) == RACE_HOSHIJIN && ch.Starphase > 0 {
		gain *= 2
	}
	if PLR_FLAGGED(ch, PLR_HEALT) && ch.Sits != nil {
		gain *= 20
	}
	if AFF_FLAGGED(ch, AFF_BLESS) {
		gain *= 2
	}
	if AFF_FLAGGED(ch, AFF_CURSE) {
		gain /= 5
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		gain /= 4
	}
	if grav_cost(ch, 0) == 0 {
		if !IS_NPC(ch) && int(ch.Chclass) != CLASS_BARDOCK && world[ch.In_room].Gravity >= 10 {
			send_to_char(ch, libc.CString("This gravity is wearing you out!\r\n"))
			gain /= 4
		}
		if !IS_NPC(ch) && int(ch.Chclass) == CLASS_BARDOCK && world[ch.In_room].Gravity > 10 {
			send_to_char(ch, libc.CString("This gravity is wearing you out!\r\n"))
			gain /= 4
		}
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
		gain = ch.Max_move - ch.Move
	}
	if cook_element(ch.In_room) == 1 {
		gain *= 2
	}
	if ch.Regen > 0 {
		gain += int64((float64(gain) * 0.01) * float64(ch.Regen))
	}
	return gain
}
func update_flags(ch *char_data) {
	if ch == nil {
		send_to_imm(libc.CString("ERROR: Empty ch variable sent to update_flags."))
		return
	}
	if (ch.Bonuses[BONUS_LATE]) != 0 && int(ch.Position) == POS_SLEEPING && rand_number(1, 3) == 3 {
		if ch.Hit >= gear_pl(ch) && ch.Move >= ch.Max_move && ch.Mana >= ch.Max_mana {
			send_to_char(ch, libc.CString("You FINALLY wake up.\r\n"))
			act(libc.CString("$n wakes up."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SITTING
		}
	}
	if AFF_FLAGGED(ch, AFF_KNOCKED) && ch.Fighting == nil {
		act(libc.CString("@W$n is no longer senseless, and wakes up.@n"), FALSE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You are no longer knocked out, and wake up!@n\r\n"))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_KNOCKED)
		ch.Position = POS_SITTING
	}
	barrier_shed(ch)
	if AFF_FLAGGED(ch, AFF_FIRESHIELD) && ch.Fighting == nil && rand_number(1, 101) > GET_SKILL(ch, SKILL_FIRESHIELD) {
		send_to_char(ch, libc.CString("Your fireshield disappears.\r\n"))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_FIRESHIELD)
	}
	if AFF_FLAGGED(ch, AFF_ZANZOKEN) && ch.Fighting == nil && rand_number(1, 3) == 2 {
		send_to_char(ch, libc.CString("You lose concentration and no longer are ready to zanzoken.\r\n"))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
	}
	if AFF_FLAGGED(ch, AFF_ENSNARED) && rand_number(1, 3) == 2 {
		send_to_char(ch, libc.CString("The silk ensnaring your arms disolves enough for you to break it!\r\n"))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_ENSNARED)
	}
	if !IS_NPC(ch) && !PLR_FLAGGED(ch, PLR_STAIL) && !PLR_FLAGGED(ch, PLR_NOGROW) && (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) {
		if ch.Player_specials.Racial_pref == 1 && rand_number(1, 50) >= 40 {
			ch.Tail_growth += 1
		} else if ch.Player_specials.Racial_pref != 1 || int(ch.Race) == RACE_SAIYAN {
			ch.Tail_growth += 1
		}
		if ch.Tail_growth == 10 {
			send_to_char(ch, libc.CString("@wYour tail grows back.@n\r\n"))
			act(libc.CString("$n@w's tail grows back.@n"), TRUE, ch, nil, nil, TO_ROOM)
			SET_BIT_AR(ch.Act[:], PLR_STAIL)
			ch.Tail_growth = 0
		}
	}
	if !IS_NPC(ch) && !PLR_FLAGGED(ch, PLR_TAIL) && (int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO) {
		ch.Tail_growth += 1
		if ch.Tail_growth == 10 {
			send_to_char(ch, libc.CString("@wYour tail grows back.@n\r\n"))
			act(libc.CString("$n@w's tail grows back.@n"), TRUE, ch, nil, nil, TO_ROOM)
			SET_BIT_AR(ch.Act[:], PLR_TAIL)
			ch.Tail_growth = 0
		}
	}
	if AFF_FLAGGED(ch, AFF_MBREAK) && rand_number(1, int(sick_fail+3)) == 2 {
		send_to_char(ch, libc.CString("@wYour mind is no longer in turmoil, you can charge ki again.@n\r\n"))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_MBREAK)
		if GET_SKILL(ch, SKILL_TELEPATHY) <= 0 && rand_number(1, 2) == 2 {
			ch.Real_abils.Intel -= 1
			ch.Real_abils.Wis -= 1
			send_to_char(ch, libc.CString("@RDue to the stress you've lost 1 Intelligence and Wisdom!@n\r\n"))
			if int(ch.Real_abils.Wis) < 4 {
				ch.Real_abils.Wis = 4
			}
			if int(ch.Real_abils.Intel) < 4 {
				ch.Real_abils.Intel = 4
			}
		} else if GET_SKILL(ch, SKILL_TELEPATHY) <= 0 && rand_number(1, 20) == 1 {
			ch.Real_abils.Intel -= 1
			ch.Real_abils.Wis -= 1
			send_to_char(ch, libc.CString("@RDue to the stress you've lost 1 Intelligence and Wisdom!@n\r\n"))
			if int(ch.Real_abils.Wis) < 4 {
				ch.Real_abils.Wis = 4
			}
			if int(ch.Real_abils.Intel) < 4 {
				ch.Real_abils.Intel = 4
			}
		}
	}
	if AFF_FLAGGED(ch, AFF_SHOCKED) && rand_number(1, 4) == 4 {
		send_to_char(ch, libc.CString("@wYour mind is no longer shocked.@n\r\n"))
		if GET_SKILL(ch, SKILL_TELEPATHY) > 0 {
			var (
				skill int = GET_SKILL(ch, SKILL_TELEPATHY)
				stop  int = FALSE
			)
			improve_skill(ch, SKILL_TELEPATHY, 0)
			for stop == FALSE {
				if rand_number(1, 8) == 5 {
					stop = TRUE
				} else {
					improve_skill(ch, SKILL_TELEPATHY, 0)
				}
			}
			if skill < GET_SKILL(ch, SKILL_TELEPATHY) {
				send_to_char(ch, libc.CString("Your mental damage and recovery has taught you things about your own mind.\r\n"))
			}
		}
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_SHOCKED)
	}
	if AFF_FLAGGED(ch, AFF_FROZEN) && rand_number(1, 2) == 2 {
		send_to_char(ch, libc.CString("@wYou realize you have thawed enough and break out of the ice holding you prisoner!\r\n"))
		act(libc.CString("$n@W breaks out of the ice holding $m prisoner!"), TRUE, ch, nil, nil, TO_ROOM)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_FROZEN)
	}
	if AFF_FLAGGED(ch, AFF_WITHER) && rand_number(1, int(sick_fail+6)) == 2 {
		send_to_char(ch, libc.CString("@wYour body returns to normal and you beat the withering that plagued you.\r\n"))
		act(libc.CString("$n@W's looks more fit now."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Real_abils.Str += 3
		ch.Real_abils.Cha += 3
		save_char(ch)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_WITHER)
	}
	if wearing_stardust(ch) == 1 {
		SET_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
		send_to_char(ch, libc.CString("The stardust armor blesses you with a free zanzoken when you next need it.\r\n"))
	}
}
func ki_gain(ch *char_data) int {
	var gain int = 0
	if IS_NPC(ch) {
		gain = GET_LEVEL(ch)
	} else {
		gain = int(ch.Max_ki / 12)
		switch ch.Position {
		case POS_SLEEPING:
			gain *= 2
		case POS_RESTING:
			gain += gain / 2
		case POS_SITTING:
			gain += gain / 4
		}
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		gain /= 4
	}
	if ch.Regen > 0 {
		gain += int((float64(gain) * 0.01) * float64(ch.Regen))
	}
	return gain
}
func set_title(ch *char_data, title *byte) {
	if ch != nil {
		send_to_char(ch, libc.CString("Title is disabled for the time being while Iovan works on a brand new and fancier title system.\r\n"))
		return
	}
}
func gain_level(ch *char_data, whichclass int) {
	if whichclass < 0 {
		whichclass = int(ch.Chclass)
	}
	if GET_LEVEL(ch) < 100 && ch.Exp >= int64(level_exp(ch, GET_LEVEL(ch)+1)) {
		ch.Level += 1
		ch.Chclass = int8(whichclass)
		advance_level(ch, whichclass)
		mudlog(BRF, int(MAX(ADMLVL_IMMORT, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("%s advanced level to level %d."), GET_NAME(ch), GET_LEVEL(ch))
		send_to_char(ch, libc.CString("You rise a level!\r\n"))
		ch.Exp -= int64(level_exp(ch, GET_LEVEL(ch)))
		write_aliases(ch)
		save_char(ch)
	}
}
func run_autowiz() {
}
func gain_exp(ch *char_data, gain int64) {
	if gain > 20000000 {
		gain = 20000000
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("EXP CANCEL: You can not gain experience from the arena.\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_WUNJO) {
		gain += int64(float64(gain) * 0.15)
	}
	if PLR_FLAGGED(ch, PLR_IMMORTAL) {
		gain = int64(float64(gain) * 0.95)
	}
	var diff int64 = int64(float64(gain) * 0.15)
	if !IS_NPC(ch) && GET_LEVEL(ch) < 1 {
		return
	}
	if IS_NPC(ch) {
		ch.Exp += gain
		return
	}
	if gain > 0 {
		gain = MIN(int64(config_info.Play.Max_exp_gain), gain)
		if (ch.Equipment[WEAR_SH]) != nil {
			var obj *obj_data = (ch.Equipment[WEAR_SH])
			if GET_OBJ_VNUM(obj) == 1127 {
				var spar int64 = gain
				gain += int64(float64(gain) * 0.25)
				spar = gain - spar
				send_to_char(ch, libc.CString("@D[@BBooster EXP@W: @G+%s@D]\r\n"), add_commas(spar))
			}
		}
		if GET_LEVEL(ch) < 100 {
			if ch.Mindlink != nil && gain > 0 && ch.Linker == 0 {
				if GET_LEVEL(ch)+20 < GET_LEVEL(ch.Mindlink) || GET_LEVEL(ch)-20 > GET_LEVEL(ch.Mindlink) {
					send_to_char(ch.Mindlink, libc.CString("The level difference between the two of you is too great to gain from mind read.\r\n"))
				} else {
					act(libc.CString("@GYou've absorbed some new experiences from @W$n@G!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_VICT)
					var read int = int(float64(gain) * 0.12)
					gain -= int64(read)
					if read == 0 {
						read = 1
					}
					gain_exp(ch.Mindlink, int64(read))
					act(libc.CString("@RYou sense that @W$N@R has stolen some of your experiences with $S mind!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_CHAR)
				}
			}
			var difff int64 = int64(level_exp(ch, GET_LEVEL(ch)+1) * 5)
			if GET_LEVEL(ch) <= 90 && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp+gain) <= (level_exp(ch, GET_LEVEL(ch)+1)-int(difff)) {
				send_to_char(ch, libc.CString("@WYou -@RNEED@W- to @ylevel@W you can't hold any more experience.@n\r\n"))
			} else if GET_LEVEL(ch) >= 91 && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) <= -1 {
				send_to_char(ch, libc.CString("@WYou -@RNEED@W- to @ylevel@W you can't hold any more experience.@n\r\n"))
			} else {
				ch.Exp += gain
			}
		}
		if GET_LEVEL(ch) < 100 && ch.Exp >= int64(level_exp(ch, GET_LEVEL(ch)+1)) {
			send_to_char(ch, libc.CString("@rYou have earned enough experience to gain a @ylevel@r.@n\r\n"))
		}
		if GET_LEVEL(ch) == 100 && ch.Admlevel < 1 {
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				diff = int64(float64(diff) * 1.3)
			}
			if int(ch.Race) == RACE_ANDROID {
				diff = int64(float64(diff) * 1.2)
			}
			if ch.Mindlink != nil && gain > 0 && ch.Linker == 0 {
				if GET_LEVEL(ch)+20 < GET_LEVEL(ch.Mindlink) || GET_LEVEL(ch)-20 > GET_LEVEL(ch.Mindlink) {
					send_to_char(ch.Mindlink, libc.CString("The level difference between the two of you is too great to gain from mind read.\r\n"))
				} else {
					act(libc.CString("@GYou've absorbed some new experiences from @W$n@G!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_VICT)
					var read int64 = int64(float64(gain) * 0.12)
					diff -= int64(float64(read) * 0.15)
					gain -= read
					if read == 0 {
						read = 1
					}
					gain_exp(ch.Mindlink, read)
					act(libc.CString("@RYou sense that @W$N@R has stolen some of your experiences with $S mind!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_CHAR)
				}
			}
			if rand_number(1, 5) >= 2 {
				if int(ch.Race) == RACE_HUMAN {
					ch.Basepl += int64(float64(diff) * 0.8)
					ch.Max_hit += int64(float64(diff) * 0.8)
				} else {
					ch.Basepl += diff
					ch.Max_hit += diff
				}
				send_to_char(ch, libc.CString("@D[@G+@Y%s @RPL@D]@n "), add_commas(diff))
			}
			if rand_number(1, 5) >= 2 {
				if int(ch.Race) == RACE_HALFBREED {
					ch.Basest += int64(float64(diff) * 0.85)
					ch.Max_move += int64(float64(diff) * 0.85)
				} else {
					ch.Basest += diff
					ch.Max_move += diff
				}
				send_to_char(ch, libc.CString("@D[@G+@Y%s @gSTA@D]@n "), add_commas(diff))
			}
			if rand_number(1, 5) >= 2 {
				ch.Baseki += diff
				ch.Max_mana += diff
				send_to_char(ch, libc.CString("@D[@G+@Y%s @CKi@D]@n"), add_commas(diff))
			}
		}
	} else if gain < 0 {
		gain = MAX(int64(-config_info.Play.Max_exp_loss), gain)
		ch.Exp += gain
		if ch.Exp < 0 {
			ch.Exp = 0
		}
	}
}
func gain_exp_regardless(ch *char_data, gain int) {
	var (
		is_altered int = FALSE
		num_levels int = 0
	)
	gain = int(float32(gain) * config_info.Play.Exp_multiplier)
	ch.Exp += int64(gain)
	if ch.Exp < 0 {
		ch.Exp = 0
	}
	if !IS_NPC(ch) {
		for GET_LEVEL(ch) < config_info.Play.Level_cap-1 && ch.Exp >= int64(level_exp(ch, GET_LEVEL(ch)+1)) {
			ch.Level += 1
			num_levels++
			advance_level(ch, int(ch.Chclass))
			is_altered = TRUE
		}
		if is_altered != 0 {
			mudlog(BRF, int(MAX(ADMLVL_IMMORT, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("%s advanced %d level%s to level %d."), GET_NAME(ch), num_levels, func() string {
				if num_levels == 1 {
					return ""
				}
				return "s"
			}(), GET_LEVEL(ch))
			if num_levels == 1 {
				send_to_char(ch, libc.CString("You rise a level!\r\n"))
			} else {
				send_to_char(ch, libc.CString("You rise %d levels!\r\n"), num_levels)
			}
		}
	}
}
func gain_condition(ch *char_data, condition int, value int) {
	var intoxicated bool
	if IS_NPC(ch) {
		return
	} else if int(ch.Race) == RACE_ANDROID {
		return
	} else if int(ch.Player_specials.Conditions[condition]) < 0 {
		return
	} else if ROOM_FLAGGED(ch.In_room, ROOM_RHELL) {
		return
	} else if ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
		return
	} else if AFF_FLAGGED(ch, AFF_SPIRIT) {
		return
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 1 {
		return
	}
	if PLR_FLAGGED(ch, PLR_WRITING) {
		return
	} else {
		intoxicated = int(ch.Player_specials.Conditions[DRUNK]) > 0
		if value > 0 {
			if int(ch.Player_specials.Conditions[condition]) >= 0 {
				if int(ch.Player_specials.Conditions[condition])+value > 48 {
					var prior int = int(ch.Player_specials.Conditions[condition])
					ch.Player_specials.Conditions[condition] = 48
					if condition != DRUNK && prior >= 48 && int(ch.Race) != RACE_MAJIN {
						var (
							pukeroll int = axion_dice(0)
							ocond    int = condition
						)
						_ = ocond
						if condition == HUNGER {
							ocond = THIRST
						} else if condition == THIRST {
							ocond = HUNGER
						}
						if pukeroll > int(ch.Aff_abils.Con)+19 {
							act(libc.CString("@r@6You retch violently until your stomach is empty! Your constitution couldn't handle being that stuffed!@n"), TRUE, ch, nil, nil, TO_CHAR)
							act(libc.CString("@m@6$n@r@6 retches violently! It seems $e stuffed $mself too much!@n"), TRUE, ch, nil, nil, TO_ROOM)
							SET_BIT_AR(ch.Affected_by[:], AFF_PUKED)
							if int(ch.Race) != RACE_NAMEK {
								ch.Player_specials.Conditions[HUNGER] -= 40
								if int(ch.Player_specials.Conditions[HUNGER]) < 0 {
									ch.Player_specials.Conditions[HUNGER] = 0
								}
								if int(ch.Race) == RACE_BIO && ((ch.Genome[0]) == 3 || (ch.Genome[1]) == 3) {
									ch.Player_specials.Conditions[HUNGER] = -1
								}
							}
							if int(ch.Race) != RACE_KANASSAN {
								ch.Player_specials.Conditions[THIRST] -= 30
								if int(ch.Player_specials.Conditions[THIRST]) < 0 {
									ch.Player_specials.Conditions[THIRST] = 0
								}
							} else {
								send_to_char(ch, libc.CString("Through your mastery of your bodily fluids you manage to retain your hydration.\r\n"))
								return
							}
						} else if pukeroll > int(ch.Aff_abils.Con)+9 {
							act(libc.CString("@r@6You puke violently! Your constitution couldn't handle being that stuffed!@n"), TRUE, ch, nil, nil, TO_CHAR)
							act(libc.CString("@m@6$n@r@6 pukes violently! It seems $e stuffed $mself too much!@n"), TRUE, ch, nil, nil, TO_ROOM)
							SET_BIT_AR(ch.Affected_by[:], AFF_PUKED)
							if int(ch.Race) != RACE_NAMEK {
								ch.Player_specials.Conditions[HUNGER] -= 20
								if int(ch.Player_specials.Conditions[HUNGER]) < 0 {
									ch.Player_specials.Conditions[HUNGER] = 0
								}
							}
							if int(ch.Race) != RACE_KANASSAN {
								ch.Player_specials.Conditions[THIRST] -= 15
								if int(ch.Player_specials.Conditions[THIRST]) < 0 {
									ch.Player_specials.Conditions[THIRST] = 0
								}
							} else {
								send_to_char(ch, libc.CString("Through your mastery of your bodily fluids you manage to retain your hydration.\r\n"))
								return
							}
						} else if pukeroll > int(ch.Aff_abils.Con) {
							act(libc.CString("@r@6You puke a little! Your constitution couldn't handle being that stuffed!@n"), TRUE, ch, nil, nil, TO_CHAR)
							act(libc.CString("@m@6$n@r@6 pukes a little! It seems $e stuffed $mself too much!@n"), TRUE, ch, nil, nil, TO_ROOM)
							SET_BIT_AR(ch.Affected_by[:], AFF_PUKED)
							if int(ch.Race) != RACE_NAMEK {
								ch.Player_specials.Conditions[HUNGER] -= 8
								if int(ch.Player_specials.Conditions[HUNGER]) < 0 {
									ch.Player_specials.Conditions[HUNGER] = 0
								}
							}
							if int(ch.Race) != RACE_KANASSAN {
								ch.Player_specials.Conditions[THIRST] -= 8
								if int(ch.Player_specials.Conditions[THIRST]) < 0 {
									ch.Player_specials.Conditions[THIRST] = 0
								}
							} else {
								send_to_char(ch, libc.CString("Through your mastery of your bodily fluids you manage to retain your hydration.\r\n"))
								return
							}
						}
					}
				} else {
					ch.Player_specials.Conditions[condition] += int8(value)
				}
			}
		}
		if !AFF_FLAGGED(ch, AFF_SPIRIT) && (GET_SKILL(ch, SKILL_SURVIVAL) == 0 || GET_SKILL(ch, SKILL_SURVIVAL) < rand_number(1, 140)) {
			if value <= 0 {
				if int(ch.Player_specials.Conditions[condition]) >= 0 {
					if AFF_FLAGGED(ch, AFF_PUKED) {
						REMOVE_BIT_AR(ch.Affected_by[:], AFF_PUKED)
					}
					if int(ch.Player_specials.Conditions[condition])+value < 0 {
						ch.Player_specials.Conditions[condition] = 0
					} else {
						ch.Player_specials.Conditions[condition] += int8(value)
					}
				}
			}
			switch condition {
			case HUNGER:
				switch ch.Player_specials.Conditions[condition] {
				case 0:
					if ch.Move >= ch.Max_move/3 {
						send_to_char(ch, libc.CString("@RYou are starving to death!@n\r\n"))
						ch.Move -= ch.Move / 3
					} else if ch.Move < ch.Max_move/3 {
						send_to_char(ch, libc.CString("@RYou are starving to death!@n\r\n"))
						ch.Move = 0
						if ch.Suppression > 0 {
							send_to_char(ch, libc.CString("@RYou stop suppressing!@n\r\n"))
							ch.Suppressed = 0
							ch.Hit += ch.Suppression
							ch.Suppression = 0
						}
						ch.Hit -= ch.Max_hit / 3
					}
				case 1:
					send_to_char(ch, libc.CString("You are extremely hungry!\r\n"))
				case 2:
					send_to_char(ch, libc.CString("You are very hungry!\r\n"))
				case 3:
					send_to_char(ch, libc.CString("You are pretty hungry!\r\n"))
				case 4:
					send_to_char(ch, libc.CString("You are hungry!\r\n"))
				case 5:
					fallthrough
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					send_to_char(ch, libc.CString("Your stomach is growling!\r\n"))
				case 9:
					fallthrough
				case 10:
					fallthrough
				case 11:
					send_to_char(ch, libc.CString("You could use something to eat.\r\n"))
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					send_to_char(ch, libc.CString("You could use a bite to eat.\r\n"))
				case 18:
					fallthrough
				case 19:
					fallthrough
				case 20:
					send_to_char(ch, libc.CString("You could use a snack.\r\n"))
				default:
				}
			case THIRST:
				switch ch.Player_specials.Conditions[condition] {
				case 0:
					if ch.Move >= ch.Max_move/3 {
						send_to_char(ch, libc.CString("@RYou are dehydrated!@n\r\n"))
						ch.Move -= ch.Move / 3
					} else if ch.Move < ch.Max_move/3 {
						send_to_char(ch, libc.CString("@RYou are dehydrated!@n\r\n"))
						ch.Move = 0
						if ch.Suppression > 0 {
							send_to_char(ch, libc.CString("@RYou stop suppressing!@n\r\n"))
							ch.Suppressed = 0
							ch.Hit += ch.Suppression
							ch.Suppression = 0
						}
						ch.Hit -= ch.Max_hit / 3
					}
				case 1:
					send_to_char(ch, libc.CString("You are extremely thirsty!\r\n"))
				case 2:
					send_to_char(ch, libc.CString("You are very thirsty!\r\n"))
				case 3:
					send_to_char(ch, libc.CString("You are pretty thirsty!\r\n"))
				case 4:
					send_to_char(ch, libc.CString("You are thirsty!\r\n"))
				case 5:
					fallthrough
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					send_to_char(ch, libc.CString("Your throat is pretty dry!\r\n"))
				case 9:
					fallthrough
				case 10:
					fallthrough
				case 11:
					send_to_char(ch, libc.CString("You could use something to drink.\r\n"))
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					send_to_char(ch, libc.CString("Your mouth feels pretty dry.\r\n"))
				case 18:
					fallthrough
				case 19:
					fallthrough
				case 20:
					send_to_char(ch, libc.CString("You could use a sip of water.\r\n"))
				default:
				}
			case DRUNK:
				if intoxicated {
					if int(ch.Player_specials.Conditions[DRUNK]) <= 0 {
						send_to_char(ch, libc.CString("You are now sober.\r\n"))
					}
				}
			default:
			}
			if ch.Hit <= 0 && int(ch.Player_specials.Conditions[HUNGER]) == 0 {
				send_to_char(ch, libc.CString("You have starved to death!\r\n"))
				ch.Move = 0
				act(libc.CString("@W$n@W falls down dead before you...@n"), FALSE, ch, nil, nil, TO_ROOM)
				die(ch, nil)
				if int(ch.Player_specials.Conditions[HUNGER]) != -1 {
					ch.Player_specials.Conditions[HUNGER] = 48
				}
				if int(ch.Player_specials.Conditions[THIRST]) != -1 {
					ch.Player_specials.Conditions[THIRST] = 48
				}
			}
			if ch.Hit <= 0 && int(ch.Player_specials.Conditions[THIRST]) == 0 {
				send_to_char(ch, libc.CString("You have died of dehydration!\r\n"))
				ch.Move = 0
				act(libc.CString("@W$n@W falls down dead before you...@n"), FALSE, ch, nil, nil, TO_ROOM)
				die(ch, nil)
				if int(ch.Player_specials.Conditions[HUNGER]) != -1 {
					ch.Player_specials.Conditions[HUNGER] = 48
				}
				ch.Player_specials.Conditions[THIRST] = 48
			}
		}
	}
}
func check_idling(ch *char_data) {
	if dball_count(ch) != 0 {
		return
	}
	if func() int {
		p := &ch.Timer
		*p++
		return *p
	}() > config_info.Play.Idle_void {
		if ch.Was_in_room == room_rnum(-1) && ch.In_room != room_rnum(-1) {
			ch.Was_in_room = ch.In_room
			if ch.Fighting != nil {
				stop_fighting(ch.Fighting)
				stop_fighting(ch)
			}
			if ch.In_room == 0 || ch.In_room == 1 {
				ch.Player_specials.Load_room = ch.Player_specials.Load_room
			}
			if !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && (int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 19800 || int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) > 0x4DBB) {
				ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
				ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(1561))))
			}
			if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 2002 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 2011 {
				ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(1960))))
			}
			if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) == 2069 {
				ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(2017))))
			}
			if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) == 2070 {
				ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(2046))))
			}
			if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 101 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 139 {
				if GET_LEVEL(ch) == 1 {
					ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(100))))
					ch.Exp = 0
				} else {
					if int(ch.Chclass) == CLASS_ROSHI {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(1130))))
					}
					if int(ch.Chclass) == CLASS_KABITO {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(0x2F42))))
					}
					if int(ch.Chclass) == CLASS_NAIL {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(0x2DA3))))
					}
					if int(ch.Chclass) == CLASS_BARDOCK {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(2268))))
					}
					if int(ch.Chclass) == CLASS_KRANE {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(0x32D1))))
					}
					if int(ch.Chclass) == CLASS_TAPION {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(8231))))
					}
					if int(ch.Chclass) == CLASS_PICCOLO {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(1659))))
					}
					if int(ch.Chclass) == CLASS_ANDSIX {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(1713))))
					}
					if int(ch.Chclass) == CLASS_DABURA {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(6486))))
					}
					if int(ch.Chclass) == CLASS_FRIEZA {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(4282))))
					}
					if int(ch.Chclass) == CLASS_GINYU {
						ch.Player_specials.Load_room = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(real_room(4289))))
					}
				}
			}
			act(libc.CString("$n disappears into the void."), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("You have been idle, and are pulled into a void.\r\n"))
			save_char(ch)
			char_from_room(ch)
			char_to_room(ch, 1)
		} else if ch.Timer > config_info.Play.Idle_rent_time {
			if ch.In_room != room_rnum(-1) {
				char_from_room(ch)
				char_to_room(ch, 3)
			}
			if ch.Desc != nil {
				send_to_char(ch, libc.CString("You are idle and are extracted safely from the game.\r\n"))
				ch.Desc.Connected = CON_DISCONNECT
				ch.Desc.Character = nil
				ch.Desc = nil
			}
			Crash_rentsave(ch, 0)
			cp(ch)
			mudlog(CMP, ADMLVL_GOD, TRUE, libc.CString("%s force-rented and extracted (idle)."), GET_NAME(ch))
			extract_char(ch)
		}
	}
}
func heal_limb(ch *char_data) {
	var (
		healrate  int = 0
		recovered int = FALSE
	)
	if PLR_FLAGGED(ch, PLR_BANDAGED) {
		healrate += 10
	}
	if int(ch.Position) == POS_SITTING {
		healrate += 1
	} else if int(ch.Position) == POS_RESTING {
		healrate += 3
	} else if int(ch.Position) == POS_SLEEPING {
		healrate += 5
	}
	if healrate > 0 {
		if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 {
			if (ch.Limb_condition[0])+healrate >= 50 {
				act(libc.CString("You realize your right arm is no longer broken."), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n starts moving $s right arm gingerly for a moment."), TRUE, ch, nil, nil, TO_ROOM)
				ch.Limb_condition[0] += healrate
				recovered = TRUE
			} else {
				ch.Limb_condition[0] += healrate
				send_to_char(ch, libc.CString("Your right arm feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[0], "%", "%")
			}
		} else if (ch.Limb_condition[0])+healrate < 100 {
			ch.Limb_condition[0] += healrate
			send_to_char(ch, libc.CString("Your right arm feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[0], "%", "%")
		} else if (ch.Limb_condition[0]) < 100 && (ch.Limb_condition[0])+healrate >= 100 {
			ch.Limb_condition[0] = 100
			send_to_char(ch, libc.CString("Your right arm has fully recovered.\r\n"))
		}
		if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 {
			if (ch.Limb_condition[1])+healrate >= 50 {
				act(libc.CString("You realize your left arm is no longer broken."), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n starts moving $s left arm gingerly for a moment."), TRUE, ch, nil, nil, TO_ROOM)
				ch.Limb_condition[1] += healrate
				recovered = TRUE
			} else {
				ch.Limb_condition[1] += healrate
				send_to_char(ch, libc.CString("Your left arm feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[0], "%", "%")
			}
		} else if (ch.Limb_condition[1])+healrate < 100 {
			ch.Limb_condition[1] += healrate
			send_to_char(ch, libc.CString("Your left arm feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[1], "%", "%")
		} else if (ch.Limb_condition[1]) < 100 && (ch.Limb_condition[1])+healrate >= 100 {
			ch.Limb_condition[1] = 100
			send_to_char(ch, libc.CString("Your left arm has fully recovered.\r\n"))
		}
		if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 {
			if (ch.Limb_condition[2])+healrate >= 50 {
				act(libc.CString("You realize your right leg is no longer broken."), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n starts moving $s right leg gingerly for a moment."), TRUE, ch, nil, nil, TO_ROOM)
				ch.Limb_condition[2] += healrate
				recovered = TRUE
			} else {
				ch.Limb_condition[2] += healrate
				send_to_char(ch, libc.CString("Your right leg feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[0], "%", "%")
			}
		} else if (ch.Limb_condition[2])+healrate < 100 {
			ch.Limb_condition[2] += healrate
			send_to_char(ch, libc.CString("Your right leg feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[2], "%", "%")
		} else if (ch.Limb_condition[2]) < 100 && (ch.Limb_condition[2])+healrate >= 100 {
			ch.Limb_condition[2] = 100
			send_to_char(ch, libc.CString("Your right leg has fully recovered.\r\n"))
		}
		if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 {
			if (ch.Limb_condition[3])+healrate >= 50 {
				act(libc.CString("You realize your left leg is no longer broken."), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n starts moving $s left leg gingerly for a moment."), TRUE, ch, nil, nil, TO_ROOM)
				ch.Limb_condition[3] += healrate
				recovered = TRUE
			} else {
				ch.Limb_condition[3] += healrate
				send_to_char(ch, libc.CString("Your left leg feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[0], "%", "%")
			}
		} else if (ch.Limb_condition[3])+healrate < 100 {
			ch.Limb_condition[3] += healrate
			send_to_char(ch, libc.CString("Your left leg feels a little better @D[@G%d%s@D/@g100%s@D]@n.\r\n"), ch.Limb_condition[3], "%", "%")
		} else if (ch.Limb_condition[3]) < 100 && (ch.Limb_condition[3])+healrate >= 100 {
			ch.Limb_condition[3] = 100
			send_to_char(ch, libc.CString("Your left leg as fully recovered.\r\n"))
		}
		if !PLR_FLAGGED(ch, PLR_BANDAGED) && recovered == TRUE {
			if axion_dice(-10) > int(ch.Aff_abils.Con) {
				ch.Real_abils.Str -= 1
				ch.Real_abils.Dex -= 1
				ch.Real_abils.Cha -= 1
				send_to_char(ch, libc.CString("@RYou lose 1 Strength, Agility, and Speed!\r\n"))
				if int(ch.Real_abils.Str) < 4 {
					ch.Real_abils.Str = 4
				}
				if int(ch.Real_abils.Con) < 4 {
					ch.Real_abils.Con = 4
				}
				if int(ch.Real_abils.Dex) < 4 {
					ch.Real_abils.Dex = 4
				}
				if int(ch.Real_abils.Cha) < 4 {
					ch.Real_abils.Cha = 4
				}
				save_char(ch)
			}
		}
	}
	if PLR_FLAGGED(ch, PLR_BANDAGED) && recovered == TRUE {
		REMOVE_BIT_AR(ch.Act[:], PLR_BANDAGED)
		send_to_char(ch, libc.CString("You remove your bandages.\r\n"))
		return
	}
}
func point_update() {
	var (
		i           *char_data
		next_char   *char_data
		j           *obj_data
		next_thing  *obj_data
		jj          *obj_data
		next_thing2 *obj_data
		vehicle     *obj_data = nil
	)
	for i = character_list; i != nil; i = next_char {
		next_char = i.Next
		if !IS_NPC(i) && i.In_room != room_rnum(-1) {
			if ROOM_FLAGGED(i.In_room, ROOM_HOUSE) {
				i.Relax_count += 1
			} else if i.Relax_count >= 464 {
				i.Relax_count -= 4
			} else if i.Relax_count >= 232 {
				i.Relax_count -= 3
			} else if i.Relax_count > 0 && rand_number(1, 3) == 3 {
				i.Relax_count -= 2
			} else {
				i.Relax_count -= 1
			}
			if i.Relax_count < 0 {
				i.Relax_count = 0
			}
		}
		if rand_number(1, 2) == 2 {
			gain_condition(i, HUNGER, -1)
		}
		if rand_number(1, 2) == 2 {
			gain_condition(i, DRUNK, -1)
		}
		if rand_number(1, 2) == 2 {
			gain_condition(i, THIRST, -1)
		}
		if IS_NPC(i) {
			i.Aggtimer = 0
		}
		if int(i.Position) >= POS_STUNNED {
			var change int = FALSE
			update_flags(i)
			if !IS_NPC(i) {
				if i.Hit < gear_pl(i) {
					change = TRUE
				}
				if i.Mana < i.Max_mana {
					change = TRUE
				}
				if i.Move < i.Max_move {
					change = TRUE
				}
			}
			if PLR_FLAGGED(i, PLR_AURALIGHT) {
				if float64(i.Mana-mana_gain(i)) > float64(i.Max_mana)*0.05 {
					send_to_char(i, libc.CString("You send more energy into your aura to keep the light active.\r\n"))
					i.Mana -= mana_gain(i)
					i.Mana -= int64(float64(i.Max_mana) * 0.05)
				} else {
					send_to_char(i, libc.CString("You don't have enough energy to keep the aura active.\r\n"))
					act(libc.CString("$n's aura slowly stops shining and fades.\r\n"), TRUE, i, nil, nil, TO_ROOM)
					REMOVE_BIT_AR(i.Act[:], PLR_AURALIGHT)
					world[i.In_room].Light--
				}
			}
			if int(i.Race) == RACE_MUTANT && ((i.Genome[0]) == 6 || (i.Genome[1]) == 6) {
				mutant_limb_regen(i)
			}
			var x int = (i.Kaioken * 5) + 5
			if i.Sleeptime > 0 && int(i.Position) != POS_SLEEPING {
				i.Sleeptime -= 1
			}
			if i.Sleeptime < 8 && int(i.Position) == POS_SLEEPING {
				i.Sleeptime += rand_number(2, 4)
				if i.Sleeptime > 8 {
					i.Sleeptime = 8
				}
			}
			if i.Kaioken > 0 && (GET_SKILL(i, SKILL_KAIOKEN) < rand_number(1, x) || i.Move <= i.Max_move/10) {
				send_to_char(i, libc.CString("You lose focus and your kaioken disappears.\r\n"))
				act(libc.CString("$n loses focus and $s kaioken aura disappears."), TRUE, i, nil, nil, TO_ROOM)
				if i.Hit-(gear_pl(i)/10)*int64(i.Kaioken) > 0 {
					i.Hit -= (gear_pl(i) / 10) * int64(i.Kaioken)
				} else {
					i.Hit = 1
				}
				i.Kaioken = 0
			} else if i.Kaioken <= 0 && !AFF_FLAGGED(i, AFF_BURNED) {
				if AFF_FLAGGED(i, AFF_METAMORPH) && float64(i.Hit) < float64(gear_pl(i))+float64(gear_pl(i))*0.6 {
					i.Hit += hit_gain(i)
					if float64(i.Hit) > float64(gear_pl(i))+float64(gear_pl(i))*0.6 {
						i.Hit = int64(float64(gear_pl(i)) + float64(gear_pl(i))*0.6)
					}
				} else {
					if !AFF_FLAGGED(i, AFF_METAMORPH) && i.Hit < gear_pl(i) {
						i.Hit += hit_gain(i)
						if i.Hit > gear_pl(i) {
							i.Hit = gear_pl(i)
						}
					}
				}
				if i.Suppression > 0 {
					if float64(i.Hit) > (float64(gear_pl(i))*0.01)*float64(i.Suppression) {
						i.Hit = int64((float64(gear_pl(i)) * 0.01) * float64(i.Suppression))
						i.Suppressed = gear_pl(i) - i.Hit
					}
				}
			}
			if AFF_FLAGGED(i, AFF_BURNED) {
				if rand_number(1, 5) >= 4 {
					send_to_char(i, libc.CString("Your burns are healed now.\r\n"))
					act(libc.CString("$n@w's burns are now healed.@n"), TRUE, i, nil, nil, TO_ROOM)
					REMOVE_BIT_AR(i.Affected_by[:], AFF_BURNED)
				}
			}
			i.Move += move_gain(i)
			i.Mana += mana_gain(i)
			if i.Move > i.Max_move {
				i.Move = i.Max_move
			}
			if i.Mana > i.Max_mana {
				i.Mana = i.Max_mana
			}
			if !IS_NPC(i) {
				heal_limb(i)
			}
			if SECT(i.In_room) == SECT_WATER_NOSWIM && i.Player_specials.Carried_by == nil && int(i.Race) != RACE_KANASSAN {
				if i.Move >= int64(gear_weight(i)) {
					act(libc.CString("@bYou swim in place.@n"), TRUE, i, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@b swims in place.@n"), TRUE, i, nil, nil, TO_ROOM)
					i.Move -= int64(gear_weight(i))
				} else {
					i.Move -= int64(gear_weight(i))
					if i.Move < 0 {
						i.Move = 0
					}
					act(libc.CString("@RYou are drowning!@n"), TRUE, i, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@b gulps water as $e struggles to stay above the water line.@n"), TRUE, i, nil, nil, TO_ROOM)
					if i.Hit-gear_pl(i)/3 <= 0 {
						act(libc.CString("@rYou drown!@n"), TRUE, i, nil, nil, TO_CHAR)
						act(libc.CString("@R$n@r drowns!@n"), TRUE, i, nil, nil, TO_ROOM)
						die(i, nil)
						i.Hit = 1
					} else {
						i.Hit -= gear_pl(i) / 3
					}
				}
			}
			if has_o2(i) == 0 && SUNKEN(i.In_room) && !ROOM_FLAGGED(i.In_room, ROOM_SPACE) {
				if (i.Mana - mana_gain(i)) > i.Max_mana/200 {
					send_to_char(i, libc.CString("Your ki holds an atmosphere around you.\r\n"))
					i.Mana -= mana_gain(i)
					i.Mana -= int64(float64(i.Max_mana) * 0.005)
				} else {
					if i.Suppressed > 0 && float64(i.Suppressed) > float64(gear_pl(i))*0.05 {
						send_to_char(i, libc.CString("You struggle trying to hold your breath!\r\n"))
						i.Suppressed -= int64(float64(i.Max_hit) * 0.05)
					} else if float64(i.Hit-hit_gain(i)) > float64(gear_pl(i))*0.05 {
						send_to_char(i, libc.CString("You struggle trying to hold your breath!\r\n"))
						i.Hit -= hit_gain(i)
						i.Hit -= int64(float64(i.Max_hit) * 0.05)
					} else if i.Hit <= i.Max_hit/20 {
						send_to_char(i, libc.CString("You have drowned!\r\n"))
						i.Hit = 1
						act(libc.CString("@W$n@W drowns right in front of you.@n"), FALSE, i, nil, nil, TO_ROOM)
						die(i, nil)
					}
				}
			}
			if has_o2(i) == 0 && ROOM_FLAGGED(i.In_room, ROOM_SPACE) {
				if float64(i.Mana-mana_gain(i)) > float64(i.Max_mana)*0.005 {
					send_to_char(i, libc.CString("Your ki holds an atmosphere around you.\r\n"))
					i.Mana -= mana_gain(i)
					i.Mana -= int64(float64(i.Max_mana) * 0.005)
				} else {
					if i.Suppressed > 0 && float64(i.Suppressed) > float64(gear_pl(i))*0.05 {
						send_to_char(i, libc.CString("You struggle trying to hold your breath!\r\n"))
						i.Suppressed -= int64(float64(i.Max_hit) * 0.05)
					} else if float64(i.Hit-hit_gain(i)) > float64(gear_pl(i))*0.05 {
						send_to_char(i, libc.CString("You struggle trying to hold your breath!\r\n"))
						i.Hit -= hit_gain(i)
						i.Hit -= int64(float64(i.Max_hit) * 0.05)
					} else if i.Hit <= i.Max_hit/20 {
						send_to_char(i, libc.CString("You have drowned!\r\n"))
						i.Hit = 1
						act(libc.CString("@W$n@W drowns right in front of you.@n"), FALSE, i, nil, nil, TO_ROOM)
						die(i, nil)
					}
				}
			}
			if !AFF_FLAGGED(i, AFF_FLYING) && world[i.In_room].Geffect == 6 && !MOB_FLAGGED(i, MOB_NOKILL) && int(i.Race) != RACE_DEMON {
				act(libc.CString("@rYour legs are burned by the lava!@n"), TRUE, i, nil, nil, TO_CHAR)
				act(libc.CString("@R$n@r's legs are burned by the lava!@n"), TRUE, i, nil, nil, TO_ROOM)
				if IS_NPC(i) && IS_HUMANOID(i) && rand_number(1, 2) == 2 {
					do_fly(i, nil, 0, 0)
				}
				if float64(i.Suppressed) > float64(gear_pl(i))*0.05 {
					i.Suppressed -= int64(float64(gear_pl(i)) * 0.05)
				} else {
					i.Suppressed = 0
					i.Suppression = 0
					i.Hit -= int64(float64(gear_pl(i)) * 0.05)
					if i.Hit < 0 {
						act(libc.CString("@rYou have burned to death!@n"), TRUE, i, nil, nil, TO_CHAR)
						act(libc.CString("@R$n@r has burned to death!@n"), TRUE, i, nil, nil, TO_ROOM)
						die(i, nil)
					}
				}
			}
			if change == TRUE && !AFF_FLAGGED(i, AFF_POISON) {
				if PLR_FLAGGED(i, PLR_HEALT) && i.Sits != nil {
					send_to_char(i, libc.CString("@wThe healing tank works wonders on your injuries.@n\r\n"))
					i.Sits.Healcharge -= rand_number(1, 2)
					if i.Sits.Healcharge == 0 {
						send_to_char(i, libc.CString("@wThe healing tank is now too low on energy to heal you.\r\n"))
						act(libc.CString("You step out of the now empty healing tank."), TRUE, i, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@w steps out of the now empty healing tank.@n"), TRUE, i, nil, nil, TO_ROOM)
						REMOVE_BIT_AR(i.Act[:], PLR_HEALT)
						i.Sits.Sitting = nil
						i.Sits = nil
					} else if i.Hit == gear_pl(i) && i.Mana == i.Max_mana && i.Move == i.Max_move {
						send_to_char(i, libc.CString("@wYou are fully recovered now.\r\n"))
						act(libc.CString("You step out of the now empty healing tank."), TRUE, i, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@w steps out of the now empty healing tank.@n"), TRUE, i, nil, nil, TO_ROOM)
						REMOVE_BIT_AR(i.Act[:], PLR_HEALT)
						i.Sits.Sitting = nil
						i.Sits = nil
					}
				} else if PLR_FLAGGED(i, PLR_HEALT) && i.Sits == nil {
					REMOVE_BIT_AR(i.Act[:], PLR_HEALT)
				} else if int(i.Position) == POS_SLEEPING {
					send_to_char(i, libc.CString("@wYour sleep does you some good.@n\r\n"))
					if int(i.Race) != RACE_ANDROID && i.Fighting == nil {
						i.Lifeforce = int64(GET_LIFEMAX(i))
					}
				} else if int(i.Position) == POS_RESTING {
					send_to_char(i, libc.CString("@wYou feel relaxed and better.@n\r\n"))
					if i.Lifeforce != int64(GET_LIFEMAX(i)) {
						if int(i.Race) != RACE_ANDROID && i.Fighting == nil && i.Suppression <= 0 && i.Hit != gear_pl(i) {
							i.Lifeforce += int64(float64(GET_LIFEMAX(i)) * 0.15)
							if i.Lifeforce > int64(GET_LIFEMAX(i)) {
								i.Lifeforce = int64(GET_LIFEMAX(i))
							}
							send_to_char(i, libc.CString("@CYou feel more lively.@n\r\n"))
						}
					}
				} else if int(i.Position) == POS_SITTING {
					send_to_char(i, libc.CString("@wYou feel rested and better.@n\r\n"))
				} else {
					send_to_char(i, libc.CString("You feel slightly better.\r\n"))
				}
			}
			if i.Hit <= 0 {
				i.Hit = 1
			}
			if AFF_FLAGGED(i, AFF_POISON) {
				var cost float64 = 0.0
				if int(i.Aff_abils.Con) >= 100 {
					cost = 0.01
				} else if int(i.Aff_abils.Con) >= 80 {
					cost = 0.02
				} else if int(i.Aff_abils.Con) >= 50 {
					cost = 0.03
				} else if int(i.Aff_abils.Con) >= 30 {
					cost = 0.04
				} else if int(i.Aff_abils.Con) >= 20 {
					cost = 0.05
				} else {
					cost = 0.06
				}
				if float64(i.Hit)-float64(i.Max_hit)*cost > 0 {
					send_to_char(i, libc.CString("You puke as the poison burns through your blood.\r\n"))
					act(libc.CString("$n shivers and then pukes."), TRUE, i, nil, nil, TO_ROOM)
					i.Hit -= int64(float64(i.Max_hit) * cost)
				} else {
					send_to_char(i, libc.CString("The poison claims your life!\r\n"))
					act(libc.CString("$n pukes up blood and falls down dead!"), TRUE, i, nil, nil, TO_ROOM)
					if i.Poisonby != nil {
						if AFF_FLAGGED(i.Poisonby, AFF_GROUP) {
							group_gain(i.Poisonby, i)
						} else {
							solo_gain(i.Poisonby, i)
						}
						die(i, i.Poisonby)
					} else {
						die(i, nil)
					}
				}
			}
			if int(i.Position) <= POS_STUNNED {
				update_pos(i)
			}
		} else if int(i.Position) == POS_INCAP {
			continue
		} else if int(i.Position) == POS_MORTALLYW {
			continue
		}
		if float64(i.Mana) >= float64(i.Max_mana)*0.5 && float64(i.Charge) < float64(i.Max_mana)*0.1 && i.Preference == PREFERENCE_KI && !PLR_FLAGGED(i, PLR_AURALIGHT) {
			i.Charge = int64(float64(i.Max_mana) * 0.1)
		}
		if !IS_NPC(i) {
			update_char_objects(i)
			update_innate(i)
			if i.Admlevel < config_info.Play.Idle_max_level {
				check_idling(i)
			} else {
				i.Timer++
			}
		}
	}
	for j = object_list; j != nil; j = next_thing {
		next_thing = j.Next
		if OBJ_FLAGGED(j, ITEM_NORENT) && j.Worn_by == nil && j.Carried_by == nil && obj_selling != j && GET_OBJ_VNUM(j) != 7200 {
			var diff libc.Time = 0
			diff = libc.GetTime(nil) - j.Lload
			if diff > 240 && j.Lload > 0 {
				basic_mud_log(libc.CString("No rent object (%s) extracted from room (%d)"), j.Short_description, GET_ROOM_VNUM(j.In_room))
				extract_obj(j)
			}
		}
		if int(j.Type_flag) == ITEM_HATCH {
			if (func() *obj_data {
				vehicle = find_vehicle_by_vnum(j.Value[VAL_HATCH_DEST])
				return vehicle
			}()) != nil {
				j.Value[3] = int(libc.BoolToInt(GET_ROOM_VNUM(vehicle.In_room)))
			}
		}
		if IS_CORPSE(j) {
			if j.Timer > 0 {
				j.Timer--
			}
			if libc.StrStr(j.Name, libc.CString("android")) == nil && libc.StrStr(j.Name, libc.CString("Android")) == nil && !OBJ_FLAGGED(j, ITEM_BURIED) {
				if j.Timer == 5 {
					if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
						act(libc.CString("@DFlies start to gather around $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						act(libc.CString("@DFlies start to gather around $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
					}
				}
				if j.Timer == 3 {
					if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
						act(libc.CString("@DA cloud of flies has formed over $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						act(libc.CString("@DA cloud of flies has formed over $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
					}
				}
				if j.Timer == 2 {
					if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
						act(libc.CString("@DMaggots can be seen crawling all over $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						act(libc.CString("@DMaggots can be seen crawling all over $p@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
					}
				}
				if j.Timer == 1 {
					if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
						act(libc.CString("@DMaggots have nearly stripped $p of all its flesh@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						act(libc.CString("@DMaggots have nearly stripped $p of all its flesh@D.@n"), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
					}
				}
			}
			if j.Timer == 0 {
				if j.Carried_by != nil {
					if libc.StrStr(j.Name, libc.CString("android")) == nil {
						act(libc.CString("$p decays in your hands."), FALSE, j.Carried_by, j, nil, TO_CHAR)
						if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
							act(libc.CString("A quivering horde of maggots consumes $p."), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
							act(libc.CString("A quivering horde of maggots consumes $p."), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						}
					} else {
						act(libc.CString("$p decays in your hands."), FALSE, j.Carried_by, j, nil, TO_CHAR)
						if j.In_room != room_rnum(-1) && world[j.In_room].People != nil {
							act(libc.CString("$p breaks down completely into a pile of junk."), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
							act(libc.CString("$p breaks down completely into a pile of junk."), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
						}
					}
				}
				for jj = j.Contains; jj != nil; jj = next_thing2 {
					next_thing2 = jj.Next_content
					obj_from_obj(jj)
					if j.In_obj != nil {
						obj_to_obj(jj, j.In_obj)
					} else if j.Carried_by != nil {
						obj_to_room(jj, j.Carried_by.In_room)
					} else if j.In_room != room_rnum(-1) {
						obj_to_room(jj, j.In_room)
					} else {
						core_dump_real(libc.CString("__FILE__"), 0)
					}
				}
				extract_obj(j)
			}
		}
		if GET_OBJ_VNUM(j) == 65 {
			if j.Healcharge < 20 && j.Sitting == nil {
				j.Healcharge += rand_number(0, 1)
			}
		}
		if int(j.Type_flag) == ITEM_PORTAL {
			if j.Timer > 0 {
				j.Timer--
			}
			if j.Timer == 0 {
				act(libc.CString("A glowing portal fades from existence."), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
				act(libc.CString("A glowing portal fades from existence."), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
				extract_obj(j)
			}
		} else if GET_OBJ_VNUM(j) == 1306 {
			if j.Timer > 0 {
				j.Timer--
			}
			if j.Timer == 0 {
				act(libc.CString("The $p@n settles to the ground and goes out."), TRUE, world[j.In_room].People, j, nil, TO_ROOM)
				act(libc.CString("A $p@n settles to the ground and goes out."), TRUE, world[j.In_room].People, j, nil, TO_CHAR)
				extract_obj(j)
			}
		} else if OBJ_FLAGGED(j, ITEM_ICE) {
			if GET_OBJ_VNUM(j) == 79 && rand_number(1, 2) == 2 {
				if world[j.In_room].Geffect >= 1 && world[j.In_room].Geffect <= 5 {
					send_to_room(j.In_room, libc.CString("The heat from the lava melts a great deal of the glacial wall and the lava cools a bit in turn.\r\n"))
					world[j.In_room].Geffect -= 1
					if float64(j.Weight)-((float64(j.Weight)*0.025)+5) > 0 {
						j.Weight -= int64((float64(j.Weight) * 0.025) + 5)
					} else {
						send_to_room(j.In_room, libc.CString("The glacial wall blocking off the %s direction melts completely away.\r\n"), dirs[j.Cost])
						extract_obj(j)
					}
				} else if float64(j.Weight)-((float64(j.Weight)*0.025)+5) > 0 {
					j.Weight -= int64((float64(j.Weight) * 0.025) + 5)
					send_to_room(j.In_room, libc.CString("The glacial wall blocking off the %s direction melts some what.\r\n"), dirs[j.Cost])
				} else {
					send_to_room(j.In_room, libc.CString("The glacial wall blocking off the %s direction melts completely away.\r\n"), dirs[j.Cost])
					extract_obj(j)
				}
			} else if GET_OBJ_VNUM(j) != 79 {
				if j.Carried_by != nil && j.In_obj == nil {
					var melt int = int((float64(j.Weight) * 0.02) + 5)
					if float64(j.Weight)-((float64(j.Weight)*0.02)+5) > 0 {
						j.Weight -= int64(melt)
						send_to_char(j.Carried_by, libc.CString("%s @wmelts a little.\r\n"), j.Short_description)
						j.Carried_by.Carry_weight -= melt
					} else {
						send_to_char(j.Carried_by, libc.CString("%s @wmelts completely away.\r\n"), j.Short_description)
						var remainder int = melt - int(j.Weight)
						j.Carried_by.Carry_weight -= melt - remainder
						extract_obj(j)
					}
				} else if j.In_room != room_rnum(-1) {
					if float64(j.Weight)-((float64(j.Weight)*0.02)+5) > 0 {
						j.Weight -= int64((float64(j.Weight) * 0.02) + 5)
						send_to_room(j.In_room, libc.CString("%s @wmelts a little.\r\n"), j.Short_description)
					} else {
						send_to_room(j.In_room, libc.CString("%s @wmelts completely away.\r\n"), j.Short_description)
						extract_obj(j)
					}
				}
			}
		} else if j.Timer > 0 {
			j.Timer--
			if j.Timer == 0 {
				timer_otrigger(j)
			}
		}
	}
}
func timed_dt(ch *char_data) {
	var (
		vict  *char_data
		rrnum room_rnum
	)
	if ch == nil {
		for rrnum = 0; rrnum < top_of_world; rrnum++ {
			world[rrnum].Timed -= int(libc.BoolToInt(world[rrnum].Timed != -1))
		}
		for vict = character_list; vict != nil; vict = vict.Next {
			if IS_NPC(vict) {
				continue
			}
			if vict.In_room == room_rnum(-1) {
				continue
			}
			if !ROOM_FLAGGED(vict.In_room, ROOM_TIMED_DT) {
				continue
			}
			timed_dt(vict)
		}
		return
	}
	if world[ch.In_room].Timed < 0 {
		world[ch.In_room].Timed = rand_number(2, 5)
		return
	}
	if world[ch.In_room].Timed == 0 {
		for vict = world[ch.In_room].People; vict != nil; vict = vict.Next_in_room {
			if IS_NPC(vict) {
				continue
			}
			if vict.Admlevel >= ADMLVL_IMMORT {
				continue
			}
			if PLR_FLAGGED(vict, PLR_NOTDEADYET) {
				continue
			}
			log_death_trap(vict)
			death_cry(vict)
			extract_char(vict)
		}
	}
}
