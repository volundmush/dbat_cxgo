package main

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unicode"
	"unsafe"
)

var commastring [64936]byte

func dispel_ash(ch *char_data) {
	var obj *obj_data
	_ = obj
	var next_obj *obj_data
	_ = next_obj
	var ash *obj_data = nil
	ash = find_obj_in_list_vnum(world[ch.In_room].Contents, 1306)
	if ash != nil {
		var roll int = axion_dice(0)
		if ash.Cost == 3 {
			if int(ch.Aff_abils.Intel) > roll {
				act(libc.CString("@GYou clear the air with the shockwaves of your power!@n"), 1, ch, ash, nil, TO_CHAR)
				act(libc.CString("@C$n@G clears the air with the shockwaves of $s power!@n"), 1, ch, ash, nil, TO_ROOM)
				extract_obj(ash)
			}
		} else if ash.Cost == 2 {
			if int(ch.Aff_abils.Intel)+10 > roll {
				act(libc.CString("@GYou clear the air with the shockwaves of your power!@n"), 1, ch, ash, nil, TO_CHAR)
				act(libc.CString("@C$n@G clears the air with the shockwaves of $s power!@n"), 1, ch, ash, nil, TO_ROOM)
				extract_obj(ash)
			}
		} else if ash.Cost == 1 {
			if int(ch.Aff_abils.Intel)+20 > roll {
				act(libc.CString("@GYou clear the air with the shockwaves of your power!@n"), 1, ch, ash, nil, TO_CHAR)
				act(libc.CString("@C$n@G clears the air with the shockwaves of $s power!@n"), 1, ch, ash, nil, TO_ROOM)
				extract_obj(ash)
			}
		}
	}
}
func has_group(ch *char_data) bool {
	var (
		k    *follow_type
		next *follow_type
	)
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		return false
	}
	if ch.Followers != nil {
		for k = ch.Followers; k != nil; k = next {
			next = k.Next
			if !AFF_FLAGGED(k.Follower, AFF_GROUP) {
				continue
			} else {
				return true
			}
		}
	} else if ch.Master != nil {
		if !AFF_FLAGGED(ch.Master, AFF_GROUP) {
			return false
		} else {
			return true
		}
	}
	return false
}
func report_party_health(ch *char_data) *byte {
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		return libc.CString("")
	}
	if ch.Followers == nil && ch.Master == nil {
		return libc.CString("")
	}
	var k *follow_type
	var next *follow_type
	var count int = 0
	var stam1 int = 8
	var stam2 int = 8
	var stam3 int = 8
	var stam4 int = 8
	var plc1 int = 4
	var plc2 int = 4
	var plc3 int = 4
	var plc4 int = 4
	var party1 *char_data = nil
	var party2 *char_data = nil
	var party3 *char_data = nil
	var party4 *char_data = nil
	var plperc1 int64 = 0
	var plperc2 int64 = 0
	var plperc3 int64 = 0
	var plperc4 int64 = 0
	var kiperc1 int64 = 0
	var kiperc2 int64 = 0
	var kiperc3 int64 = 0
	var kiperc4 int64 = 0
	var result_party_health [64936]byte
	var result1 [64936]byte
	var result2 [64936]byte
	var result3 [64936]byte
	var result4 [64936]byte
	var result5 [64936]byte
	var plcol [5]*byte = [5]*byte{libc.CString("@r"), libc.CString("@y"), libc.CString("@Y"), libc.CString("@G"), libc.CString("")}
	var exhaust [9]*byte = [9]*byte{libc.CString("Exhausted"), libc.CString("Strained"), libc.CString("Very Tired"), libc.CString("Tired"), libc.CString("Kinda Tired"), libc.CString("Very Winded"), libc.CString("Winded"), libc.CString("Energetic"), libc.CString("?????????")}
	var excol [9]*byte = [9]*byte{libc.CString("@r"), libc.CString("@R"), libc.CString("@R"), libc.CString("@M"), libc.CString("@M"), libc.CString("@M"), libc.CString("@G"), libc.CString("@g"), libc.CString("@w")}
	if ch.Followers != nil {
		for k = ch.Followers; k != nil; k = next {
			next = k.Next
			if !AFF_FLAGGED(k.Follower, AFF_GROUP) {
				continue
			}
			if k.Follower != ch {
				count += 1
				if count == 1 {
					party1 = k.Follower
					plperc1 = (party1.Hit * 100) / party1.Max_hit
					kiperc1 = (party1.Charge * 100) / party1.Max_mana
					if plperc1 >= 80 {
						plc1 = 3
					} else if plperc1 >= 50 {
						plc1 = 2
					} else if plperc1 >= 30 {
						plc1 = 1
					} else {
						plc1 = 0
					}
					if party1.Move >= party1.Max_move {
						stam1 = 7
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.9 {
						stam1 = 6
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.8 {
						stam1 = 5
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.7 {
						stam1 = 4
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.5 {
						stam1 = 3
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.4 {
						stam1 = 2
					} else if float64(party1.Move) >= float64(party1.Max_move)*0.2 {
						stam1 = 1
					} else {
						stam1 = 0
					}
				} else if count == 2 {
					party2 = k.Follower
					plperc2 = (party2.Hit * 100) / party2.Max_hit
					kiperc2 = (party2.Charge * 100) / party2.Max_mana
					if plperc2 >= 80 {
						plc2 = 3
					} else if plperc2 >= 50 {
						plc2 = 2
					} else if plperc2 >= 30 {
						plc2 = 1
					} else {
						plc2 = 0
					}
					if party2.Move >= party2.Max_move {
						stam2 = 7
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.9 {
						stam2 = 6
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.8 {
						stam2 = 5
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.7 {
						stam2 = 4
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.5 {
						stam2 = 3
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.4 {
						stam2 = 2
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.2 {
						stam2 = 1
					} else {
						stam2 = 0
					}
				} else if count == 3 {
					party3 = k.Follower
					plperc3 = (party3.Hit * 100) / party3.Max_hit
					kiperc3 = (party3.Charge * 100) / party3.Max_mana
					if plperc3 >= 80 {
						plc3 = 3
					} else if plperc3 >= 50 {
						plc3 = 2
					} else if plperc3 >= 30 {
						plc3 = 1
					} else {
						plc3 = 0
					}
					if party3.Move >= party3.Max_move {
						stam3 = 7
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.9 {
						stam3 = 6
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.8 {
						stam3 = 5
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.7 {
						stam3 = 4
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.5 {
						stam3 = 3
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.4 {
						stam3 = 2
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.2 {
						stam3 = 1
					} else {
						stam3 = 0
					}
				} else if count == 4 {
					party4 = k.Follower
					plperc4 = (party4.Hit * 100) / party4.Max_hit
					kiperc4 = (party4.Charge * 100) / party4.Max_mana
					if plperc4 >= 80 {
						plc4 = 3
					} else if plperc4 >= 50 {
						plc4 = 2
					} else if plperc4 >= 30 {
						plc4 = 1
					} else {
						plc4 = 0
					}
					if party4.Move >= party4.Max_move {
						stam4 = 7
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.9 {
						stam4 = 6
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.8 {
						stam4 = 5
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.7 {
						stam4 = 4
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.5 {
						stam4 = 3
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.4 {
						stam4 = 2
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.2 {
						stam4 = 1
					} else {
						stam4 = 0
					}
				}
			}
		}
		stdio.Sprintf(&result1[0], "@D[@BG@D]-------@mF@D------- -------@mF@D------- -------@mF@D------- -------@mF@D-------[@BG@D] <@RV@Y%s@R>@n\n", add_commas(int64(ch.Combatexpertise)))
		stdio.Sprintf(&result2[0], "@D[@BR@D]@C%-15s %-15s %-15s %-15s@D[@BR@D]@n\n", func() *byte {
			if party1 != nil {
				return get_i_name(ch, party1)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party2 != nil {
				return get_i_name(ch, party2)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party3 != nil {
				return get_i_name(ch, party3)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party4 != nil {
				return get_i_name(ch, party4)
			}
			return libc.CString("Empty")
		}())
		stdio.Sprintf(&result3[0], "@D[@BO@D]@RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s@D[@BO@D]@n\n", plcol[plc1], plperc1, "%", plcol[plc2], plperc2, "%", plcol[plc3], plperc3, "%", plcol[plc4], plperc4, "%")
		stdio.Sprintf(&result4[0], "@D[@BU@D]@cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s@D[@BU@D]@n\n", kiperc1, "%", kiperc2, "%", kiperc3, "%", kiperc4, "%")
		stdio.Sprintf(&result5[0], "@D[@BP@D]@gSt@D:%s%12s @gSt@D:%s%12s @gSt@D:%s%12s @gSt@D:%s%12s@D[@BP@D]@n", excol[stam1], exhaust[stam1], excol[stam2], exhaust[stam2], excol[stam3], exhaust[stam3], excol[stam4], exhaust[stam4])
		stdio.Sprintf(&result_party_health[0], "%s%s%s%s%s\n", &result1[0], &result2[0], &result3[0], &result4[0], &result5[0])
		ch.Temp_prompt = libc.StrDup(&result_party_health[0])
		return ch.Temp_prompt
	} else if ch.Master != nil && AFF_FLAGGED(ch.Master, AFF_GROUP) {
		party1 = ch.Master
		plperc1 = (party1.Hit * 100) / party1.Max_hit
		kiperc1 = (party1.Charge * 100) / party1.Max_mana
		if plperc1 >= 80 {
			plc1 = 3
		} else if plperc1 >= 50 {
			plc1 = 2
		} else if plperc1 >= 30 {
			plc1 = 1
		} else {
			plc1 = 0
		}
		if party1.Move >= party1.Max_move {
			stam1 = 7
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.9 {
			stam1 = 6
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.8 {
			stam1 = 5
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.7 {
			stam1 = 4
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.5 {
			stam1 = 3
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.4 {
			stam1 = 2
		} else if float64(party1.Move) >= float64(party1.Max_move)*0.2 {
			stam1 = 1
		} else {
			stam1 = 0
		}
		count = 1
		for k = party1.Followers; k != nil; k = next {
			next = k.Next
			if !AFF_FLAGGED(k.Follower, AFF_GROUP) {
				continue
			}
			if k.Follower != ch {
				count += 1
				if count == 2 {
					party2 = k.Follower
					plperc2 = (party2.Hit * 100) / party2.Max_hit
					kiperc2 = (party2.Charge * 100) / party2.Max_mana
					if plperc2 >= 80 {
						plc2 = 3
					} else if plperc2 >= 50 {
						plc2 = 2
					} else if plperc2 >= 30 {
						plc2 = 1
					} else {
						plc2 = 0
					}
					if party2.Move >= party2.Max_move {
						stam2 = 7
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.9 {
						stam2 = 6
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.8 {
						stam2 = 5
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.7 {
						stam2 = 4
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.5 {
						stam2 = 3
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.4 {
						stam2 = 2
					} else if float64(party2.Move) >= float64(party2.Max_move)*0.2 {
						stam2 = 1
					} else {
						stam2 = 0
					}
				} else if count == 3 {
					party3 = k.Follower
					plperc3 = (party3.Hit * 100) / party3.Max_hit
					kiperc3 = (party3.Charge * 100) / party3.Max_mana
					if plperc3 >= 80 {
						plc3 = 3
					} else if plperc3 >= 50 {
						plc3 = 2
					} else if plperc3 >= 30 {
						plc3 = 1
					} else {
						plc3 = 0
					}
					if party3.Move >= party3.Max_move {
						stam3 = 7
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.9 {
						stam3 = 6
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.8 {
						stam3 = 5
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.7 {
						stam3 = 4
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.5 {
						stam3 = 3
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.4 {
						stam3 = 2
					} else if float64(party3.Move) >= float64(party3.Max_move)*0.2 {
						stam3 = 1
					} else {
						stam3 = 0
					}
				} else if count == 4 {
					party4 = k.Follower
					plperc4 = (party4.Hit * 100) / party4.Max_hit
					kiperc4 = (party4.Charge * 100) / party4.Max_mana
					if plperc4 >= 80 {
						plc4 = 3
					} else if plperc4 >= 50 {
						plc4 = 2
					} else if plperc4 >= 30 {
						plc4 = 1
					} else {
						plc4 = 0
					}
					if party4.Move >= party4.Max_move {
						stam4 = 7
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.9 {
						stam4 = 6
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.8 {
						stam4 = 5
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.7 {
						stam4 = 4
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.5 {
						stam4 = 3
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.4 {
						stam4 = 2
					} else if float64(party4.Move) >= float64(party4.Max_move)*0.2 {
						stam4 = 1
					} else {
						stam4 = 0
					}
				}
			}
		}
		stdio.Sprintf(&result1[0], "@D[@BG@D]-------@YL@D------- -------@mF@D------- -------@mF@D------- -------@mF@D-------[@BG@D]@n\n")
		stdio.Sprintf(&result2[0], "@D[@BR@D]@C%-15s %-15s %-15s %-15s@D[@BR@D]@n\n", func() *byte {
			if party1 != nil {
				return get_i_name(ch, party1)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party2 != nil {
				return get_i_name(ch, party2)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party3 != nil {
				return get_i_name(ch, party3)
			}
			return libc.CString("Empty")
		}(), func() *byte {
			if party4 != nil {
				return get_i_name(ch, party4)
			}
			return libc.CString("Empty")
		}())
		stdio.Sprintf(&result3[0], "@D[@BO@D]@RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s @RPL@D:%s%11lld@w%s@D[@BO@D]@n\n", plcol[plc1], plperc1, "%", plcol[plc2], plperc2, "%", plcol[plc3], plperc3, "%", plcol[plc4], plperc4, "%")
		stdio.Sprintf(&result4[0], "@D[@BU@D]@cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s @cCharge@D:@B%7lld@w%s@D[@BU@D]@n\n", kiperc1, "%", kiperc2, "%", kiperc3, "%", kiperc4, "%")
		stdio.Sprintf(&result5[0], "@D[@BP@D]@gSt@D:%s%12s @gSt@D:%s%12s @gSt@D:%s%12s @gSt@D:%s%12s@D[@BP@D]@n", excol[stam1], exhaust[stam1], excol[stam2], exhaust[stam2], excol[stam3], exhaust[stam3], excol[stam4], exhaust[stam4])
		stdio.Sprintf(&result_party_health[0], "%s%s%s%s%s\n", &result1[0], &result2[0], &result3[0], &result4[0], &result5[0])
		ch.Temp_prompt = libc.StrDup(&result_party_health[0])
		return ch.Temp_prompt
	} else {
		return libc.CString("")
	}
}
func know_skill(ch *char_data, skill int) int {
	var know int = 0
	if GET_SKILL(ch, skill) > 0 {
		know = 1
	}
	if int(ch.Stupidkiss) == skill {
		know = 2
	}
	if know == 0 {
		send_to_char(ch, libc.CString("You do not know how to perform %s.\r\n"), spell_info[skill].Name)
		know = 0
	} else if know == 2 {
		send_to_char(ch, libc.CString("@WYou try to use @M%s@W but lingering thoughts of a certain kiss distracts you!@n\r\n"), spell_info[skill].Name)
		send_to_char(ch, libc.CString("You must sleep in order to cure this.\r\n"))
		know = 0
	}
	return know
}
func roll_aff_duration(num int, add int) int {
	var (
		start   int = num / 20
		finish  int = num / 10
		outcome int = add
	)
	outcome += rand_number(start, finish)
	return outcome
}
func null_affect(ch *char_data, aff_flag int) {
	var (
		af      *affected_type
		next_af *affected_type
	)
	for af = ch.Affected; af != nil; af = next_af {
		next_af = af.Next
		if af.Location == APPLY_NONE && int(af.Bitvector) == aff_flag {
			affect_remove(ch, af)
		}
	}
}
func assign_affect(ch *char_data, aff_flag int, skill int, dur int, str int, con int, intel int, agl int, wis int, spd int) {
	var (
		af  [6]affected_type
		num int = 0
	)
	if dur <= 0 {
		dur = 1
	}
	if str == 0 && con == 0 && wis == 0 && intel == 0 && agl == 0 && spd == 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = 0
		af[num].Location = APPLY_NONE
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if str != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = str
		af[num].Location = APPLY_STR
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if con != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = con
		af[num].Location = APPLY_CON
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if intel != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = intel
		af[num].Location = APPLY_INT
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if agl != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = agl
		af[num].Location = APPLY_DEX
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if spd != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = spd
		af[num].Location = APPLY_CHA
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
	if wis != 0 {
		af[num].Type = int16(skill)
		af[num].Duration = int16(dur)
		af[num].Modifier = wis
		af[num].Location = APPLY_WIS
		af[num].Bitvector = uint32(int32(aff_flag))
		affect_join(ch, &af[num], false, false, false, false)
		num += 1
	}
}
func sec_roll_check(ch *char_data) int {
	var (
		figure  int = 0
		chance  int = 0
		outcome int = 0
	)
	figure = int((float64(GET_LEVEL(ch)) * 1.6) + 10)
	chance = axion_dice(0) + axion_dice(0) + rand_number(0, 20)
	if figure >= chance {
		outcome = 1
	}
	return outcome
}
func get_measure(ch *char_data, height int, weight int) int {
	var amt int = 0
	if !PLR_FLAGGED(ch, PLR_OOZARU) && (int(ch.Race) != RACE_ICER || !IS_TRANSFORMED(ch)) && (ch.Genome[0]) < 9 {
		if height > 0 {
			amt = height
		} else if weight > 0 {
			amt = weight
		}
	} else if int(ch.Race) == RACE_ICER && PLR_FLAGGED(ch, PLR_TRANS1) {
		if height > 0 {
			amt = height * 3
		} else if weight > 0 {
			amt = weight * 4
		}
	} else if int(ch.Race) == RACE_ICER && PLR_FLAGGED(ch, PLR_TRANS2) {
		if height > 0 {
			amt = height * 3
		} else if weight > 0 {
			amt = weight * 5
		}
	} else if int(ch.Race) == RACE_ICER && PLR_FLAGGED(ch, PLR_TRANS3) {
		if height > 0 {
			amt = int(float64(height) * 1.5)
		} else if weight > 0 {
			amt = weight * 2
		}
	} else if int(ch.Race) == RACE_ICER && PLR_FLAGGED(ch, PLR_TRANS4) {
		if height > 0 {
			amt = height * 2
		} else if weight > 0 {
			amt = weight * 3
		}
	} else if PLR_FLAGGED(ch, PLR_OOZARU) || (ch.Genome[0]) == 9 {
		if height > 0 {
			amt = height * 10
		} else if weight > 0 {
			amt = weight * 50
		}
	}
	return amt
}
func physical_cost(ch *char_data, skill int) int64 {
	var result int64 = 0
	if skill == SKILL_PUNCH {
		result = ch.Max_hit / 500
	} else if skill == SKILL_KICK {
		result = ch.Max_hit / 350
	} else if skill == SKILL_ELBOW {
		result = ch.Max_hit / 400
	} else if skill == SKILL_KNEE {
		result = ch.Max_hit / 300
	} else if skill == SKILL_UPPERCUT {
		result = ch.Max_hit / 200
	} else if skill == SKILL_ROUNDHOUSE {
		result = ch.Max_hit / 150
	} else if skill == SKILL_HEELDROP {
		result = ch.Max_hit / 80
	} else if skill == SKILL_SLAM {
		result = ch.Max_hit / 90
	}
	var cou1 int = rand_number(1, 20) + 1
	var cou2 int = cou1 + rand_number(1, 6)
	result += int64(rand_number(cou1, cou2))
	if int(ch.Skills[SKILL_STYLE]) >= 100 {
		result -= int64(float64(result) * 0.4)
	} else if int(ch.Skills[SKILL_STYLE]) >= 75 {
		if int(ch.Chclass) == CLASS_TSUNA {
			result -= int64(float64(result) * 0.4)
		} else if int(ch.Chclass) == CLASS_TAPION && (skill == SKILL_PUNCH || skill == SKILL_KICK) {
			result -= int64(float64(result) * 0.35)
		} else if int(ch.Chclass) == CLASS_JINTO {
			if int(ch.Skills[skill]) >= 100 {
				result -= int64(float64(result) * 0.45)
			} else {
				result -= int64(float64(result) * 0.25)
			}
		} else {
			result -= int64(float64(result) * 0.25)
		}
	} else if int(ch.Skills[SKILL_STYLE]) >= 50 {
		result -= int64(float64(result) * 0.25)
	}
	if int(ch.Race) == RACE_ANDROID {
		result = int64(float64(result) * .25)
	}
	return result
}
func android_can(ch *char_data) int {
	var obj *obj_data = (ch.Equipment[WEAR_BACKPACK])
	if obj == nil {
		return 0
	} else if GET_OBJ_VNUM(obj) == 1806 {
		return 1
	} else if GET_OBJ_VNUM(obj) == 1807 {
		return 2
	} else {
		return 0
	}
}
func trans_cost(ch *char_data, trans int) int {
	if (ch.Transcost[trans]) == 0 {
		return 50
	} else {
		return 0
	}
}
func trans_req(ch *char_data, trans int) int {
	var requirement int = 0
	if int(ch.Race) == RACE_HUMAN {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 1500000
			} else if ch.Transclass == 2 {
				requirement = 1800000
			} else if ch.Transclass == 3 {
				requirement = 2100000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 37500000
			} else if ch.Transclass == 2 {
				requirement = 35000000
			} else if ch.Transclass == 3 {
				requirement = 32500000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 200000000
			} else if ch.Transclass == 2 {
				requirement = 190000000
			} else if ch.Transclass == 3 {
				requirement = 185000000
			}
		case 4:
			if ch.Transclass == 1 {
				requirement = 1400000000
			} else if ch.Transclass == 2 {
				requirement = 1200000000
			} else if ch.Transclass == 3 {
				requirement = 1100000000
			}
		}
	}
	if int(ch.Race) == RACE_HALFBREED {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 1500000
			} else if ch.Transclass == 2 {
				requirement = 1400000
			} else if ch.Transclass == 3 {
				requirement = 1200000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 60000000
			} else if ch.Transclass == 2 {
				requirement = 55000000
			} else if ch.Transclass == 3 {
				requirement = 50000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 1200000000
			} else if ch.Transclass == 2 {
				requirement = 1050000000
			} else if ch.Transclass == 3 {
				requirement = 950000000
			}
		}
	}
	if int(ch.Race) == RACE_SAIYAN {
		if PLR_FLAGGED(ch, PLR_LSSJ) {
			switch trans {
			case 1:
				if ch.Transclass == 1 {
					requirement = 600000
				} else if ch.Transclass == 2 {
					requirement = 500000
				} else if ch.Transclass == 3 {
					requirement = 450000
				}
			case 2:
				if ch.Transclass == 1 {
					requirement = 300000000
				} else if ch.Transclass == 2 {
					requirement = 250000000
				} else if ch.Transclass == 3 {
					requirement = 225000000
				}
			}
		} else {
			switch trans {
			case 1:
				if ch.Transclass == 1 {
					requirement = 1400000
				} else if ch.Transclass == 2 {
					requirement = 1200000
				} else if ch.Transclass == 3 {
					requirement = 1100000
				}
			case 2:
				if ch.Transclass == 1 {
					requirement = 60000000
				} else if ch.Transclass == 2 {
					requirement = 55000000
				} else if ch.Transclass == 3 {
					requirement = 50000000
				}
			case 3:
				if ch.Transclass == 1 {
					requirement = 160000000
				} else if ch.Transclass == 2 {
					requirement = 150000000
				} else if ch.Transclass == 3 {
					requirement = 140000000
				}
			case 4:
				if ch.Transclass == 1 {
					requirement = 1800000000
				} else if ch.Transclass == 2 {
					requirement = 1625000000
				} else if ch.Transclass == 3 {
					requirement = 1400000000
				}
			}
		}
	}
	if int(ch.Race) == RACE_NAMEK {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 400000
			} else if ch.Transclass == 2 {
				requirement = 360000
			} else if ch.Transclass == 3 {
				requirement = 335000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 10000000
			} else if ch.Transclass == 2 {
				requirement = 9500000
			} else if ch.Transclass == 3 {
				requirement = 8000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 240000000
			} else if ch.Transclass == 2 {
				requirement = 220000000
			} else if ch.Transclass == 3 {
				requirement = 200000000
			}
		case 4:
			if ch.Transclass == 1 {
				requirement = 950000000
			} else if ch.Transclass == 2 {
				requirement = 900000000
			} else if ch.Transclass == 3 {
				requirement = 875000000
			}
		}
	}
	if int(ch.Race) == RACE_ICER {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 550000
			} else if ch.Transclass == 2 {
				requirement = 500000
			} else if ch.Transclass == 3 {
				requirement = 450000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 20000000
			} else if ch.Transclass == 2 {
				requirement = 17500000
			} else if ch.Transclass == 3 {
				requirement = 15000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 180000000
			} else if ch.Transclass == 2 {
				requirement = 150000000
			} else if ch.Transclass == 3 {
				requirement = 125000000
			}
		case 4:
			if ch.Transclass == 1 {
				requirement = 880000000
			} else if ch.Transclass == 2 {
				requirement = 850000000
			} else if ch.Transclass == 3 {
				requirement = 820000000
			}
		}
	}
	if int(ch.Race) == RACE_MAJIN {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 2400000
			} else if ch.Transclass == 2 {
				requirement = 2200000
			} else if ch.Transclass == 3 {
				requirement = 2000000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 50000000
			} else if ch.Transclass == 2 {
				requirement = 45000000
			} else if ch.Transclass == 3 {
				requirement = 40000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 1800000000
			} else if ch.Transclass == 2 {
				requirement = 1550000000
			} else if ch.Transclass == 3 {
				requirement = 1300000000
			}
		}
	}
	if int(ch.Race) == RACE_TRUFFLE {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 3800000
			} else if ch.Transclass == 2 {
				requirement = 3600000
			} else if ch.Transclass == 3 {
				requirement = 3500000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 400000000
			} else if ch.Transclass == 2 {
				requirement = 300000000
			} else if ch.Transclass == 3 {
				requirement = 200000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 1550000000
			} else if ch.Transclass == 2 {
				requirement = 1450000000
			} else if ch.Transclass == 3 {
				requirement = 1250000000
			}
		}
	}
	if int(ch.Race) == RACE_MUTANT {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 200000
			} else if ch.Transclass == 2 {
				requirement = 180000
			} else if ch.Transclass == 3 {
				requirement = 160000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 30000000
			} else if ch.Transclass == 2 {
				requirement = 27500000
			} else if ch.Transclass == 3 {
				requirement = 25000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 750000000
			} else if ch.Transclass == 2 {
				requirement = 700000000
			} else if ch.Transclass == 3 {
				requirement = 675000000
			}
		}
	}
	if int(ch.Race) == RACE_KAI {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 3250000
			} else if ch.Transclass == 2 {
				requirement = 3000000
			} else if ch.Transclass == 3 {
				requirement = 2850000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 700000000
			} else if ch.Transclass == 2 {
				requirement = 650000000
			} else if ch.Transclass == 3 {
				requirement = 625000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 1500000000
			} else if ch.Transclass == 2 {
				requirement = 1300000000
			} else if ch.Transclass == 3 {
				requirement = 1250000000
			}
		}
	}
	if int(ch.Race) == RACE_KONATSU {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 2000000
			} else if ch.Transclass == 2 {
				requirement = 1800000
			} else if ch.Transclass == 3 {
				requirement = 1600000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 250000000
			} else if ch.Transclass == 2 {
				requirement = 225000000
			} else if ch.Transclass == 3 {
				requirement = 200000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 1600000000
			} else if ch.Transclass == 2 {
				requirement = 1400000000
			} else if ch.Transclass == 3 {
				requirement = 1300000000
			}
		}
	}
	if int(ch.Race) == RACE_ANDROID {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 1200000
			} else if ch.Transclass == 2 {
				requirement = 1000000
			} else if ch.Transclass == 3 {
				requirement = 850000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 8500000
			} else if ch.Transclass == 2 {
				requirement = 8000000
			} else if ch.Transclass == 3 {
				requirement = 7750000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 55000000
			} else if ch.Transclass == 2 {
				requirement = 50000000
			} else if ch.Transclass == 3 {
				requirement = 40000000
			}
		case 4:
			if ch.Transclass == 1 {
				requirement = 325000000
			} else if ch.Transclass == 2 {
				requirement = 300000000
			} else if ch.Transclass == 3 {
				requirement = 275000000
			}
		case 5:
			if ch.Transclass == 1 {
				requirement = 900000000
			} else if ch.Transclass == 2 {
				requirement = 800000000
			} else if ch.Transclass == 3 {
				requirement = 750000000
			}
		case 6:
			if ch.Transclass == 1 {
				requirement = 1300000000
			} else if ch.Transclass == 2 {
				requirement = 1200000000
			} else if ch.Transclass == 3 {
				requirement = 1100000000
			}
		}
	}
	if int(ch.Race) == RACE_BIO {
		switch trans {
		case 1:
			if ch.Transclass == 1 {
				requirement = 2000000
			} else if ch.Transclass == 2 {
				requirement = 1800000
			} else if ch.Transclass == 3 {
				requirement = 1700000
			}
		case 2:
			if ch.Transclass == 1 {
				requirement = 30000000
			} else if ch.Transclass == 2 {
				requirement = 25000000
			} else if ch.Transclass == 3 {
				requirement = 20000000
			}
		case 3:
			if ch.Transclass == 1 {
				requirement = 235000000
			} else if ch.Transclass == 2 {
				requirement = 220000000
			} else if ch.Transclass == 3 {
				requirement = 210000000
			}
		case 4:
			if ch.Transclass == 1 {
				requirement = 1500000000
			} else if ch.Transclass == 2 {
				requirement = 1300000000
			} else if ch.Transclass == 3 {
				requirement = 1150000000
			}
		}
	}
	return requirement
}
func customWrite(ch *char_data, obj *obj_data) {
	if IS_NPC(ch) {
		return
	}
	var fname [40]byte
	var line [256]byte
	var prev [256]byte
	var buf [64936]byte
	var fl *stdio.File
	var file *stdio.File
	if !get_filename(&fname[0], uint64(40), CUSTOME_FILE, ch.Desc.User) {
		basic_mud_log(libc.CString("ERROR: Custom unable to be saved to user file!"))
		return
	}
	if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: Custom unable to be saved to user file!"))
		return
	}
	for int(file.IsEOF()) == 0 {
		get_line(file, &line[0])
		if libc.StrCaseCmp(&prev[0], &line[0]) != 0 {
			stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n", &line[0])
		}
		prev[0] = '\x00'
		stdio.Sprintf(&prev[0], libc.GoString(&line[0]))
	}
	file.Close()
	if !get_filename(&fname[0], uint64(40), CUSTOME_FILE, ch.Desc.User) {
		basic_mud_log(libc.CString("ERROR: Custom unable to be saved to user file!"))
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: Custom unable to be saved to user file!"))
		return
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n", obj.Short_description)
	stdio.Fprintf(fl, "%s\n", &buf[0])
	fl.Close()
}
func customRead(d *descriptor_data, type_ int, name *byte) {
	var (
		fname  [40]byte
		line   [256]byte
		filler [256]byte
		fl     *stdio.File
		buf    [64936]byte
	)
	if type_ == 1 {
		if !get_filename(&fname[0], uint64(40), CUSTOME_FILE, name) {
			basic_mud_log(libc.CString("ERROR: Custom unable to be read from user file!"))
			return
		}
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: Custom file unable to be read!"))
			return
		}
		var buf [64936]byte
		for int(fl.IsEOF()) == 0 {
			get_line(fl, &line[0])
			if libc.StrCaseCmp(&filler[0], &line[0]) != 0 {
				stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n", &line[0])
			}
			filler[0] = '\x00'
			line[0] = '\x00'
			stdio.Sprintf(&filler[0], libc.GoString(&line[0]))
		}
		send_to_char(d.Character, &buf[0])
		fl.Close()
		return
	} else {
		if !get_filename(&fname[0], uint64(40), CUSTOME_FILE, d.User) {
			basic_mud_log(libc.CString("ERROR: Custom unable to be read from user file!"))
			return
		}
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: Custom file unable to be read!"))
			return
		}
		for int(fl.IsEOF()) == 0 {
			get_line(fl, &line[0])
			if libc.StrCaseCmp(&filler[0], &line[0]) != 0 {
				stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n", &line[0])
			}
			filler[0] = '\x00'
			stdio.Sprintf(&filler[0], libc.GoString(&line[0]))
		}
		write_to_output(d, &buf[0])
		fl.Close()
	}
}
func customCreate(d *descriptor_data) {
	if d == nil {
		return
	}
	if d.Customfile == 1 {
		return
	}
	var fname [40]byte
	var fl *stdio.File
	if !get_filename(&fname[0], uint64(40), CUSTOME_FILE, d.User) {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not create custom file."))
		return
	}
	stdio.Fprintf(fl, "@D--@RUser @GC@gustom@Gs@D--@n\n")
	d.Customfile = 1
	fl.Close()
}
func show_softcap(ch *char_data) int64 {
	var capamt int64 = 0
	if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) || int(ch.Race) != RACE_ANDROID && int(ch.Race) != RACE_BIO && int(ch.Race) != RACE_MAJIN {
		if GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
			capamt = int64(GET_LEVEL(ch) * 1500000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 2000000)
			}
		}
		if GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
			capamt = int64(GET_LEVEL(ch) * 800000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 1000000)
			}
		}
		if GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
			capamt = int64(GET_LEVEL(ch) * 250000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 300000)
			}
		}
		if GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
			capamt = int64(GET_LEVEL(ch) * 200000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 250000)
			}
		}
		if GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
			capamt = int64(GET_LEVEL(ch) * 80000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 100000)
			}
		}
		if GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
			capamt = int64(GET_LEVEL(ch) * 20000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 40000)
			}
		}
		if GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
			capamt = int64(GET_LEVEL(ch) * 15000)
			if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
				capamt = int64(GET_LEVEL(ch) * 25000)
			}
		}
		if GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
			capamt = int64(GET_LEVEL(ch) * 5000)
		}
		if GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
			capamt = int64(GET_LEVEL(ch) * 1500)
		}
		if GET_LEVEL(ch) <= 10 {
			capamt = int64(GET_LEVEL(ch) * 500)
		}
	} else {
		if GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
			capamt = int64(GET_LEVEL(ch) * 4500000)
		}
		if GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
			capamt = int64(GET_LEVEL(ch) * 2400000)
		}
		if GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
			capamt = int64(GET_LEVEL(ch) * 750000)
		}
		if GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
			capamt = int64(GET_LEVEL(ch) * 600000)
		}
		if GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
			capamt = int64(GET_LEVEL(ch) * 240000)
		}
		if GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
			capamt = int64(GET_LEVEL(ch) * 60000)
		}
		if GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
			capamt = int64(GET_LEVEL(ch) * 45000)
		}
		if GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
			capamt = int64(GET_LEVEL(ch) * 15000)
		}
		if GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
			capamt = int64(GET_LEVEL(ch) * 4500)
		}
		if GET_LEVEL(ch) <= 10 {
			capamt = int64(GET_LEVEL(ch) * 1500)
		}
	}
	return capamt
}
func disp_align(ch *char_data) *byte {
	var align int
	if ch.Alignment < -800 {
		align = 8
	} else if ch.Alignment < -600 {
		align = 7
	} else if ch.Alignment < -300 {
		align = 6
	} else if ch.Alignment < -50 {
		align = 5
	} else if ch.Alignment < 51 {
		align = 4
	} else if ch.Alignment < 300 {
		align = 3
	} else if ch.Alignment < 600 {
		align = 2
	} else if ch.Alignment < 800 {
		align = 1
	} else {
		align = 0
	}
	return alignments[align]
}
func senseCreate(ch *char_data) {
	var (
		fname [40]byte
		fl    *stdio.File
	)
	if !get_filename(&fname[0], uint64(40), SENSE_FILE, GET_NAME(ch)) {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save sense memory of, %s, to filename, %s."), GET_NAME(ch), &fname[0])
		return
	}
	stdio.Fprintf(fl, "0\n")
	fl.Close()
	return
}
func read_sense_memory(ch *char_data, vict *char_data) int {
	var (
		fname  [40]byte
		line   [256]byte
		known  int = 0
		idnums int = -1337
		fl     *stdio.File
	)
	if vict == nil {
		basic_mud_log(libc.CString("Noone."))
		return 0
	}
	if !get_filename(&fname[0], uint64(40), SENSE_FILE, GET_NAME(ch)) {
		senseCreate(ch)
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		return 2
	}
	if vict == ch {
		fl.Close()
		return 0
	}
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%d\n", &idnums)
		if IS_NPC(vict) {
			if idnums == GET_MOB_VNUM(vict) {
				known = 1
			}
		} else {
			if idnums == int(vict.Id) {
				known = 1
			}
		}
	}
	fl.Close()
	return known
}
func sense_memory_write(ch *char_data, vict *char_data) {
	var (
		file      *stdio.File
		fname     [40]byte
		line      [256]byte
		idnums    [500]int = [500]int{}
		fl        *stdio.File
		count     int = 0
		x         int = 0
		id_sample int
	)
	if !get_filename(&fname[0], uint64(40), SENSE_FILE, GET_NAME(ch)) {
		senseCreate(ch)
	}
	if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		return
	}
	for int(file.IsEOF()) == 0 || count < 498 {
		get_line(file, &line[0])
		stdio.Sscanf(&line[0], "%d\n", &id_sample)
		idnums[count] = id_sample
		count++
	}
	file.Close()
	if !get_filename(&fname[0], uint64(40), SENSE_FILE, GET_NAME(ch)) {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save sense memory file, %s, to filename, %s."), GET_NAME(ch), &fname[0])
		return
	}
	for x < count {
		if x == 0 || idnums[x-1] != idnums[x] {
			if !IS_NPC(vict) {
				if idnums[x] != int(vict.Id) {
					stdio.Fprintf(fl, "%d\n", idnums[x])
				}
			} else {
				if idnums[x] != GET_MOB_VNUM(vict) {
					stdio.Fprintf(fl, "%d\n", idnums[x])
				}
			}
		}
		x++
	}
	if !IS_NPC(vict) {
		stdio.Fprintf(fl, "%d\n", vict.Id)
	} else {
		stdio.Fprintf(fl, "%d\n", GET_MOB_VNUM(vict))
	}
	fl.Close()
	return
}
func roll_pursue(ch *char_data, vict *char_data) bool {
	var (
		skill int
		perc  int = axion_dice(0)
	)
	if ch == nil || vict == nil {
		return false
	}
	if !IS_NPC(ch) {
		skill = GET_SKILL(ch, SKILL_PURSUIT)
	} else if IS_NPC(ch) && !MOB_FLAGGED(ch, MOB_SENTINEL) {
		skill = GET_LEVEL(ch)
		if ROOM_FLAGGED(vict.In_room, ROOM_NOMOB) {
			skill = -1
		}
	} else {
		skill = -1
	}
	if !IS_NPC(vict) {
		if IS_NPC(ch) && vict.Desc == nil {
			skill = -1
		}
	}
	if skill > perc {
		var inroom int = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
		act(libc.CString("@C$n@R pursues after the fleeing @c$N@R!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		char_from_room(ch)
		char_to_room(ch, vict.In_room)
		act(libc.CString("@GYou pursue right after @c$N@G!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@R pursues after you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@R pursues after the fleeing @c$N@R!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		var k *follow_type
		var next *follow_type
		if ch.Followers != nil {
			for k = ch.Followers; k != nil; k = next {
				next = k.Next
				if int(libc.BoolToInt(GET_ROOM_VNUM(k.Follower.In_room))) == inroom && int(k.Follower.Position) >= POS_STANDING && (!AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(k.Follower, AFF_GROUP)) {
					act(libc.CString("You follow $N."), 1, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
					act(libc.CString("$n follows after $N."), 1, k.Follower, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("$n follows after you."), 1, k.Follower, nil, unsafe.Pointer(ch), TO_VICT)
					char_from_room(k.Follower)
					char_to_room(k.Follower, ch.In_room)
				}
			}
		}
		REMOVE_BIT_AR(vict.Affected_by[:], AFF_PURSUIT)
		return true
	} else {
		send_to_char(ch, libc.CString("@RYou fail to pursue after them!@n\r\n"))
		if ch.Fighting != nil {
			stop_fighting(ch)
		}
		if vict.Fighting != nil {
			stop_fighting(vict)
		}
		return false
	}
}
func broken_update() {
	var (
		k            *obj_data
		money        *obj_data
		rand_gravity [14]int = [14]int{0, 10, 20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 5000, 10000}
		dice         int     = rand_number(2, 12)
		grav_roll    int     = 0
		grav_change  int     = 0
		health       int     = 0
	)
	for k = object_list; k != nil; k = k.Next {
		if k.Carried_by != nil {
			continue
		}
		if rand_number(1, 2) == 2 {
			continue
		}
		health = k.Value[VAL_ALL_HEALTH]
		if GET_OBJ_VNUM(k) == 11 {
			grav_roll = rand_number(0, 13)
			if health <= 10 {
				grav_change = 1
			} else if health <= 40 && dice <= 8 {
				grav_change = 1
			} else if health <= 80 && dice <= 5 {
				grav_change = 1
			} else if health <= 99 && dice <= 3 {
				grav_change = 1
			}
			if grav_change == 1 {
				world[k.In_room].Gravity = rand_gravity[grav_roll]
				k.Weight = int64(rand_gravity[grav_roll])
				send_to_room(k.In_room, libc.CString("@RThe gravity generator malfunctions! The gravity level has changed!@n\r\n"))
			}
		}
		if GET_OBJ_VNUM(k) == 3034 {
			if health <= 10 {
				send_to_room(k.In_room, libc.CString("@RThe ATM machine shoots smoking bills from its money slot. The bills burn up as they float through the air!@n\r\n"))
			} else if health <= 40 && dice <= 8 {
				send_to_room(k.In_room, libc.CString("@RGibberish flashes across the cracked ATM info screen.@n\r\n"))
			} else if health <= 80 && dice == 4 {
				send_to_room(k.In_room, libc.CString("@GThe damaged ATM spits out some money while flashing ERROR on its screen!@n\r\n"))
				money = create_money(rand_number(1, 30))
				obj_to_room(money, k.In_room)
			} else if health <= 99 && dice < 4 {
				send_to_room(k.In_room, libc.CString("@RThe ATM machine emits a loud grinding sound from inside.@n\r\n"))
			}
		}
		dice = rand_number(2, 12)
	}
}
func wearable_obj(obj *obj_data) bool {
	var pass bool = false
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_FINGER) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_NECK) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_BODY) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_HEAD) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_LEGS) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_FEET) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_HANDS) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_ARMS) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_SHIELD) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_ABOUT) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_WAIST) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_WRIST) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_WIELD) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_EYE) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_PACK) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_SH) {
		pass = true
	}
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_EAR) {
		pass = true
	}
	return pass
}
func randomize_eq(obj *obj_data) {
	if wearable_obj(obj) && !OBJ_FLAGGED(obj, ITEM_NORANDOM) {
		var (
			value int = 0
			slot  int = 0
			roll  int = rand_number(2, 12)
			slot1 int = 1
			slot2 int = 1
			slot3 int = 1
		)
		_ = slot3
		var slot4 int = 1
		_ = slot4
		var slot5 int = 1
		_ = slot5
		var slot6 int = 1
		_ = slot6
		var stat int = 0
		var strength int = 0
		var wisdom int = 0
		var intelligence int = 0
		var dexterity int = 0
		var speed int = 0
		var constitution int = 0
		var i int
		for i = 0; i < 6; i++ {
			stat = obj.Affected[slot].Location
			value = obj.Affected[slot].Modifier
			if stat == 1 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					strength = 1
				}
			} else if stat == 2 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					dexterity = 1
				}
			} else if stat == 3 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					intelligence = 1
				}
			} else if stat == 4 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					wisdom = 1
				}
			} else if stat == 5 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					constitution = 1
				}
			} else if stat == 6 {
				if roll == 12 {
					value += 3
				} else if roll >= 9 {
					value += 2
				} else if roll >= 6 {
					value += 1
				} else if roll == 3 {
					value -= 1
				} else if roll <= 2 {
					value -= 2
				}
				if obj.Level >= 80 {
					if value <= 0 {
						value = 1
					}
				} else if obj.Level >= 60 {
					if value < 0 {
						value = 0
					}
				}
				if value == 0 {
					obj.Affected[slot].Location = 0
					obj.Affected[slot].Modifier = 0
				} else {
					obj.Affected[slot].Modifier = value
					speed = 1
				}
			} else if stat == 0 {
				switch slot {
				case 1:
					slot1 = 0
				case 2:
					slot2 = 0
				case 3:
					slot3 = 0
				case 4:
					slot4 = 0
				case 5:
					slot5 = 0
				case 6:
					slot6 = 0
				}
			}
			slot += 1
			roll = rand_number(2, 12)
		}
		if slot1 == 0 {
			if strength == 0 && rand_number(1, 6) == 1 {
				strength = 1
				obj.Affected[0].Location = 1
				obj.Affected[0].Modifier = 1
			} else if dexterity == 0 && rand_number(1, 6) == 1 {
				dexterity = 1
				obj.Affected[0].Location = 2
				obj.Affected[0].Modifier = 1
			} else if intelligence == 0 && rand_number(1, 6) == 1 {
				intelligence = 1
				obj.Affected[0].Location = 3
				obj.Affected[0].Modifier = 1
			} else if wisdom == 0 && rand_number(1, 6) == 1 {
				wisdom = 1
				obj.Affected[0].Location = 4
				obj.Affected[0].Modifier = 1
			} else if constitution == 0 && rand_number(1, 6) == 1 {
				constitution = 1
				obj.Affected[0].Location = 5
				obj.Affected[0].Modifier = 1
			} else if speed == 0 && rand_number(1, 6) == 1 {
				speed = 1
				obj.Affected[0].Location = 6
				obj.Affected[0].Modifier = 1
			}
		}
		if slot2 == 0 && roll >= 10 {
			if strength == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 1
				obj.Affected[1].Modifier = 1
			} else if dexterity == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 2
				obj.Affected[1].Modifier = 1
			} else if intelligence == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 3
				obj.Affected[1].Modifier = 1
			} else if wisdom == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 4
				obj.Affected[1].Modifier = 1
			} else if constitution == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 5
				obj.Affected[1].Modifier = 1
			} else if speed == 0 && rand_number(1, 6) == 1 {
				obj.Affected[1].Location = 6
				obj.Affected[1].Modifier = 1
			}
		}
		var dice int = rand_number(2, 12)
		if dice >= 10 {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOT2)
		} else if dice >= 7 {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOT1)
		}
	}
}
func sense_location(ch *char_data) *byte {
	var (
		message *byte = (*byte)(libc.Malloc(MAX_INPUT_LENGTH))
		roomnum int   = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
		num     int   = 0
	)
	if (func() int {
		num = real_zone_by_thing(roomnum)
		return num
	}()) != int(-1) {
		num = real_zone_by_thing(roomnum)
	}
	switch num {
	case 2:
		stdio.Sprintf(message, "East of Nexus City")
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 7:
		if roomnum < 795 {
			stdio.Sprintf(message, "Nexus City")
		} else {
			stdio.Sprintf(message, "South Ocean")
		}
	case 8:
		fallthrough
	case 9:
		fallthrough
	case 10:
		fallthrough
	case 11:
		if roomnum < 1133 {
			stdio.Sprintf(message, "South Ocean")
		} else if roomnum < 1179 {
			stdio.Sprintf(message, "Nexus Field")
		} else {
			stdio.Sprintf(message, "Cherry Blossom Mountain")
		}
	case 12:
		fallthrough
	case 13:
		if roomnum < 1287 {
			stdio.Sprintf(message, "Cherry Blossom Mountain")
		} else {
			stdio.Sprintf(message, "Sandy Desert")
		}
	case 14:
		if roomnum < 1428 {
			stdio.Sprintf(message, "Sandy Desert")
		} else if roomnum < 1484 {
			stdio.Sprintf(message, "Northern Plains")
		} else if roomnum < 1496 {
			stdio.Sprintf(message, "Korin's Tower")
		} else {
			stdio.Sprintf(message, "Kami's Lookout")
		}
	case 15:
		if roomnum < 1577 {
			stdio.Sprintf(message, "Kami's Lookout")
		} else if roomnum < 1580 {
			stdio.Sprintf(message, "Northern Plains")
		} else if roomnum < 1589 {
			stdio.Sprintf(message, "Kami's Lookout")
		} else {
			stdio.Sprintf(message, "Shadow Forest")
		}
	case 16:
		stdio.Sprintf(message, "Shadow Forest")
	case 17:
		fallthrough
	case 18:
		if roomnum < 1715 {
			stdio.Sprintf(message, "Decrepit Area")
		} else {
			stdio.Sprintf(message, "Inside Cherry Blossom Mountain")
		}
	case 19:
		stdio.Sprintf(message, "West City")
	case 20:
		if roomnum < 2012 {
			stdio.Sprintf(message, "West City")
		} else if roomnum > 2070 {
			stdio.Sprintf(message, "West City")
		} else {
			stdio.Sprintf(message, "Silver Mine")
		}
	case 21:
		if roomnum < 2141 {
			stdio.Sprintf(message, "West City")
		} else {
			stdio.Sprintf(message, "Hercule Beach")
		}
	case 22:
		stdio.Sprintf(message, "Vegetos City")
	case 23:
		fallthrough
	case 24:
		if roomnum < 2334 {
			stdio.Sprintf(message, "Vegetos City")
		} else if roomnum > 2462 {
			stdio.Sprintf(message, "Vegetos City")
		} else {
			stdio.Sprintf(message, "Vegetos Palace")
		}
	case 25:
		fallthrough
	case 26:
		if roomnum < 2616 {
			stdio.Sprintf(message, "Blood Dunes")
		} else {
			stdio.Sprintf(message, "Ancestral Mountains")
		}
	case 27:
		if roomnum < 2709 {
			stdio.Sprintf(message, "Ancestral Mountains")
		} else if roomnum < 2736 {
			stdio.Sprintf(message, "Destopa Swamp")
		} else {
			stdio.Sprintf(message, "Swamp Base")
		}
	case 28:
		stdio.Sprintf(message, "Pride Forest")
	case 29:
		fallthrough
	case 30:
		fallthrough
	case 31:
		stdio.Sprintf(message, "Pride Tower")
	case 32:
		stdio.Sprintf(message, "Ruby Cave")
	case 34:
		stdio.Sprintf(message, "Utatlan City")
	case 35:
		stdio.Sprintf(message, "Zenith Jungle")
	case 40:
		fallthrough
	case 41:
		fallthrough
	case 42:
		stdio.Sprintf(message, "Ice Crown City")
	case 43:
		if roomnum < 4351 {
			stdio.Sprintf(message, "Ice Highway")
		} else {
			stdio.Sprintf(message, "Topica Snowfield")
		}
	case 44:
		fallthrough
	case 45:
		stdio.Sprintf(message, "Glug's Volcano")
	case 46:
		fallthrough
	case 47:
		stdio.Sprintf(message, "Platonic Sea")
	case 48:
		stdio.Sprintf(message, "Slave City")
	case 49:
		if roomnum < 4915 {
			stdio.Sprintf(message, "Descent Down Icecrown")
		} else if roomnum != 4915 && roomnum < 4994 {
			stdio.Sprintf(message, "Topica Snowfield")
		} else {
			stdio.Sprintf(message, "Ice Highway")
		}
	case 50:
		stdio.Sprintf(message, "Mirror Shard Maze")
	case 51:
		if roomnum < 5150 {
			stdio.Sprintf(message, "Acturian Woods")
		} else if roomnum < 5165 {
			stdio.Sprintf(message, "Desolate Demesne")
		} else {
			stdio.Sprintf(message, "Chateau Ishran")
		}
	case 52:
		stdio.Sprintf(message, "Wyrm Spine Mountain")
	case 53:
		fallthrough
	case 54:
		stdio.Sprintf(message, "Aromina Hunting Preserve")
	case 55:
		stdio.Sprintf(message, "Cloud Ruler Temple")
	case 56:
		stdio.Sprintf(message, "Koltoan Mine")
	case 78:
		stdio.Sprintf(message, "Orium Cave")
	case 79:
		stdio.Sprintf(message, "Crystalline Forest")
	case 80:
		fallthrough
	case 81:
		fallthrough
	case 82:
		stdio.Sprintf(message, "Tiranoc City")
	case 83:
		stdio.Sprintf(message, "Great Oroist Temple")
	case 84:
		if roomnum < 8447 {
			stdio.Sprintf(message, "Elsthuan Forest")
		} else {
			stdio.Sprintf(message, "Mazori Farm")
		}
	case 85:
		stdio.Sprintf(message, "Dres")
	case 86:
		stdio.Sprintf(message, "Colvian Farm")
	case 87:
		stdio.Sprintf(message, "Saint Alucia")
	case 88:
		if roomnum < 8847 {
			stdio.Sprintf(message, "Meridius Memorial")
		} else {
			stdio.Sprintf(message, "Battlefields")
		}
	case 89:
		if roomnum < 8954 {
			stdio.Sprintf(message, "Desert of Illusion")
		} else {
			stdio.Sprintf(message, "Plains of Confusion")
		}
	case 90:
		stdio.Sprintf(message, "Shadowlas Temple")
	case 92:
		stdio.Sprintf(message, "Turlon Fair")
	case 97:
		stdio.Sprintf(message, "Wetlands")
	case 98:
		if roomnum < 9855 {
			stdio.Sprintf(message, "Wetlands")
		} else if roomnum < 9866 {
			stdio.Sprintf(message, "Kerberos")
		} else {
			stdio.Sprintf(message, "Shaeras Mansion")
		}
	case 99:
		if roomnum < 9907 {
			stdio.Sprintf(message, "Slavinos Ravine")
		} else if roomnum < 9960 {
			stdio.Sprintf(message, "Kerberos")
		} else {
			stdio.Sprintf(message, "Furian Citadel")
		}
	case 100:
		fallthrough
	case 101:
		fallthrough
	case 102:
		fallthrough
	case 103:
		fallthrough
	case 104:
		fallthrough
	case 105:
		fallthrough
	case 106:
		fallthrough
	case 107:
		fallthrough
	case 108:
		fallthrough
	case 109:
		fallthrough
	case 110:
		fallthrough
	case 111:
		fallthrough
	case 112:
		fallthrough
	case 113:
		fallthrough
	case 114:
		fallthrough
	case 115:
		stdio.Sprintf(message, "Namekian Wilderness")
	case 116:
		if roomnum < 0x2D98 {
			stdio.Sprintf(message, "Senzu Village")
		} else if roomnum > 0x2D98 && roomnum < 0x2DB2 {
			stdio.Sprintf(message, "Senzu Village")
		} else {
			stdio.Sprintf(message, "Guru's House")
		}
	case 117:
		fallthrough
	case 118:
		fallthrough
	case 119:
		stdio.Sprintf(message, "Crystalline Cave")
	case 120:
		stdio.Sprintf(message, "Haven City")
	case 121:
		if roomnum < 0x2F47 {
			stdio.Sprintf(message, "Haven City")
		} else {
			stdio.Sprintf(message, "Serenity Lake")
		}
	case 122:
		stdio.Sprintf(message, "Serenity Lake")
	case 123:
		stdio.Sprintf(message, "Kaiju Forest")
	case 124:
		if roomnum < 0x30C0 {
			stdio.Sprintf(message, "Ortusian Temple")
		} else {
			stdio.Sprintf(message, "Silent Glade")
		}
	case 125:
		stdio.Sprintf(message, "Near Serenity Lake")
	case 130:
		fallthrough
	case 131:
		if roomnum < 0x3361 {
			stdio.Sprintf(message, "Satan City")
		} else if roomnum == 0x3361 {
			stdio.Sprintf(message, "West City")
		} else if roomnum == 0x3362 {
			stdio.Sprintf(message, "Nexus City")
		} else {
			stdio.Sprintf(message, "South Ocean")
		}
	case 132:
		if roomnum < 0x33B0 {
			stdio.Sprintf(message, "Frieza's Ship")
		} else {
			stdio.Sprintf(message, "Namekian Wilderness")
		}
	case 133:
		stdio.Sprintf(message, "Elder Village")
	case 134:
		stdio.Sprintf(message, "Satan City")
	case 140:
		stdio.Sprintf(message, "Yardra City")
	case 141:
		stdio.Sprintf(message, "Jade Forest")
	case 142:
		stdio.Sprintf(message, "Jade Cliff")
	case 143:
		stdio.Sprintf(message, "Mount Valaria")
	case 149:
		fallthrough
	case 150:
		stdio.Sprintf(message, "Aquis City")
	case 151:
		fallthrough
	case 152:
		fallthrough
	case 153:
		stdio.Sprintf(message, "Kanassan Ocean")
	case 154:
		stdio.Sprintf(message, "Kakureta Village")
	case 155:
		stdio.Sprintf(message, "Captured Aether City")
	case 156:
		stdio.Sprintf(message, "Yunkai Pirate Base")
	case 160:
		fallthrough
	case 161:
		stdio.Sprintf(message, "Janacre")
	case 165:
		stdio.Sprintf(message, "Arlian Wasteland")
	case 166:
		stdio.Sprintf(message, "Arlian Mine")
	case 167:
		stdio.Sprintf(message, "Kilnak Caverns")
	case 168:
		stdio.Sprintf(message, "Kemabra Wastes")
	case 169:
		stdio.Sprintf(message, "Dark of Arlia")
	case 174:
		stdio.Sprintf(message, "Fistarl Volcano")
	case 175:
		fallthrough
	case 176:
		stdio.Sprintf(message, "Cerria Colony")
	case 182:
		stdio.Sprintf(message, "Below Tiranoc")
	case 196:
		stdio.Sprintf(message, "Ancient Castle")
	default:
		stdio.Sprintf(message, "Unknown.")
	}
	return message
}
func reveal_hiding(ch *char_data, type_ int) {
	if IS_NPC(ch) || !AFF_FLAGGED(ch, AFF_HIDE) {
		return
	}
	var rand1 int = rand_number(-5, 5)
	var rand2 int = rand_number(-5, 5)
	var bonus int = 0
	if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
		bonus = 10
	}
	if type_ == 0 {
		act(libc.CString("@MYou feel as though what you just did may have revealed your hiding spot...@n"), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n moves a little and you notice them!@n"), 1, ch, nil, nil, TO_ROOM)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_HIDE)
	} else if type_ == 1 {
		if GET_SKILL(ch, SKILL_HIDE)+bonus < axion_dice(0) {
			act(libc.CString("@MYou feel as though what you just did may have revealed your hiding spot...@n"), 1, ch, nil, nil, TO_CHAR)
			act(libc.CString("@M$n moves a little and you notice them!@n"), 1, ch, nil, nil, TO_ROOM)
			REMOVE_BIT_AR(ch.Affected_by[:], AFF_HIDE)
		}
	} else if type_ == 2 {
		var (
			d   *descriptor_data
			tch *char_data = nil
		)
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Connected != CON_PLAYING {
				continue
			}
			tch = d.Character
			if tch == ch {
				continue
			}
			if tch.In_room != ch.In_room {
				continue
			}
			if GET_SKILL(tch, SKILL_SPOT)+rand1 >= GET_SKILL(ch, SKILL_HIDE)+rand2 {
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_HIDE)
				act(libc.CString("@M$N seems to notice you!@n"), 1, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				act(libc.CString("@MYou notice $n's movements reveal $s hiding spot!@n"), 1, ch, nil, unsafe.Pointer(tch), TO_VICT)
				act(libc.CString("@MYou notice $N look keenly somewhere nearby. At that spot you see $n hiding!@n"), 1, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
				return
			}
		}
	} else if type_ == 3 {
		var (
			d   *descriptor_data
			tch *char_data = nil
		)
		act(libc.CString("@MThe scouter makes some beeping sounds as you tinker with its buttons.@n"), 1, ch, nil, nil, TO_CHAR)
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Connected != CON_PLAYING {
				continue
			}
			tch = d.Character
			if tch == ch {
				continue
			}
			if tch.In_room != ch.In_room {
				continue
			}
			if GET_SKILL(tch, SKILL_LISTEN) > axion_dice(0) {
				switch type_ {
				case 3:
					act(libc.CString("@MYou notice some beeping sounds that sound really close by.@n"), 1, ch, nil, unsafe.Pointer(tch), TO_VICT)
				default:
					act(libc.CString("@MYou notice some sounds coming from this room but can't seem to locate the source...@n"), 1, ch, nil, unsafe.Pointer(tch), TO_VICT)
				}
			}
		}
	} else if type_ == 4 {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_HIDE)
	}
}
func block_calc(ch *char_data) bool {
	var blocker *char_data = nil
	if ch.Blocked != nil {
		blocker = ch.Blocked
	} else {
		return true
	}
	if GET_SPEEDI(ch) < GET_SPEEDI(blocker) && int(blocker.Position) > POS_SITTING {
		if !AFF_FLAGGED(blocker, AFF_BLIND) && !PLR_FLAGGED(blocker, PLR_EYEC) {
			var minimum int = int(blocker.Aff_abils.Cha) + rand_number(5, 20)
			if minimum > 100 {
				minimum = 100
			}
			if GET_SKILL(ch, SKILL_ESCAPE_ARTIST) == 0 || GET_SKILL(ch, SKILL_ESCAPE_ARTIST) != 0 && GET_SKILL(ch, SKILL_ESCAPE_ARTIST) < rand_number(minimum, 120) {
				act(libc.CString("$n tries to leave, but can't outrun $N!"), 1, ch, nil, unsafe.Pointer(blocker), TO_NOTVICT)
				act(libc.CString("$n tries to leave, but can't outrun you!"), 1, ch, nil, unsafe.Pointer(blocker), TO_VICT)
				act(libc.CString("You try to leave, but can't outrun $N!"), 1, ch, nil, unsafe.Pointer(blocker), TO_CHAR)
				if AFF_FLAGGED(ch, AFF_FLYING) && !AFF_FLAGGED(blocker, AFF_FLYING) && ch.Altitude == 1 {
					send_to_char(blocker, libc.CString("You're now floating in the air.\r\n"))
					SET_BIT_AR(blocker.Affected_by[:], AFF_FLYING)
					blocker.Altitude = ch.Altitude
				} else if AFF_FLAGGED(ch, AFF_FLYING) && !AFF_FLAGGED(blocker, AFF_FLYING) && ch.Altitude == 2 {
					send_to_char(blocker, libc.CString("You're now floating high in the sky.\r\n"))
					SET_BIT_AR(blocker.Affected_by[:], AFF_FLYING)
					blocker.Altitude = ch.Altitude
				}
				return false
			} else {
				act(libc.CString("$n proves $s great skill and escapes from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_NOTVICT)
				act(libc.CString("$n proves $s great skill and escapes from your attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_VICT)
				act(libc.CString("Using your great skill you manage to escape from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_CHAR)
				ch.Blocked = nil
				blocker.Blocks = nil
			}
		} else {
			act(libc.CString("$n proves $s great skill and escapes from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_NOTVICT)
			act(libc.CString("$n proves $s great skill and escapes from your attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_VICT)
			act(libc.CString("Using your great skill you manage to escape from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_CHAR)
			ch.Blocked = nil
			blocker.Blocks = nil
		}
	} else if int(blocker.Position) <= POS_SITTING {
		act(libc.CString("$n proves $s great skill and escapes from $N!"), 1, ch, nil, unsafe.Pointer(blocker), TO_NOTVICT)
		act(libc.CString("$n proves $s great skill and escapes from you!"), 1, ch, nil, unsafe.Pointer(blocker), TO_VICT)
		act(libc.CString("Using your great skill you manage to escape from $N!"), 1, ch, nil, unsafe.Pointer(blocker), TO_CHAR)
		ch.Blocked = nil
		blocker.Blocks = nil
	} else if int(blocker.Position) > POS_SITTING {
		act(libc.CString("$n proves $s great skill and escapes from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_NOTVICT)
		act(libc.CString("$n proves $s great skill and escapes from your attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_VICT)
		act(libc.CString("Using your great skill you manage to escape from $N's attempted block!"), 1, ch, nil, unsafe.Pointer(blocker), TO_CHAR)
		ch.Blocked = nil
		blocker.Blocks = nil
	}
	return true
}
func molt_threshold(ch *char_data) int64 {
	var (
		threshold int64 = 0
		max       int64 = 2000000000
	)
	max *= 250
	if int(ch.Race) != RACE_ARLIAN {
		return 0
	} else if ch.Moltlevel < 100 {
		threshold = int64(((float64(ch.Moltlevel+1) * (float64(ch.Max_hit) * 0.02)) * float64(ch.Aff_abils.Con)) / 4)
		threshold = int64(float64(threshold) * 0.25)
	} else if ch.Moltlevel < 200 {
		threshold = int64(((float64(ch.Moltlevel+1) * (float64(ch.Max_hit) * 0.02)) * float64(ch.Aff_abils.Con)) / 2)
		threshold = int64(float64(threshold) * 0.2)
	} else if ch.Moltlevel < 400 {
		threshold = int64((float64(ch.Moltlevel+1) * (float64(ch.Max_hit) * 0.02)) * float64(ch.Aff_abils.Con))
		threshold = int64(float64(threshold) * 0.17)
	} else if ch.Moltlevel < 800 {
		threshold = int64(((float64(ch.Moltlevel+1) * (float64(ch.Max_hit) * 0.02)) * float64(ch.Aff_abils.Con)) * 2)
		threshold = int64(float64(threshold) * 0.15)
	} else {
		threshold = int64(((float64(ch.Moltlevel+1) * (float64(ch.Max_hit) * 0.02)) * float64(ch.Aff_abils.Con)) * 4)
		threshold = int64(float64(threshold) * 0.12)
	}
	if threshold > max {
		threshold = max
	}
	return threshold
}
func armor_evolve(ch *char_data) int {
	var value int = 0
	if ch.Moltlevel <= 5 {
		value = 8
	} else if ch.Moltlevel <= 10 {
		value = 12
	} else if ch.Moltlevel <= 20 {
		value = 15
	} else if ch.Moltlevel <= 30 {
		value = 20
	} else if ch.Moltlevel <= 40 {
		value = 30
	} else if ch.Moltlevel <= 50 {
		value = 50
	} else if ch.Moltlevel <= 75 {
		value = 100
	} else if ch.Moltlevel <= 100 {
		value = 150
	} else if ch.Moltlevel <= 500 {
		value = 200
	} else {
		value = 220
	}
	return value
}
func handle_evolution(ch *char_data, dmg int64) {
	if IS_NPC(ch) || int(ch.Race) != RACE_ARLIAN {
		return
	}
	var moltgain int64 = 0
	moltgain = int64(float64(dmg) * 0.5)
	if GET_LEVEL(ch) == 100 {
		moltgain += 100000
	} else if GET_LEVEL(ch) >= 90 {
		moltgain += int64(GET_LEVEL(ch) * 1000)
	} else if GET_LEVEL(ch) >= 75 {
		moltgain += int64(GET_LEVEL(ch) * 500)
	} else if GET_LEVEL(ch) >= 50 {
		moltgain += int64(GET_LEVEL(ch) * 250)
	} else if GET_LEVEL(ch) >= 10 {
		moltgain += int64(GET_LEVEL(ch) * 50)
	}
	ch.Moltexp += moltgain
	if AFF_FLAGGED(ch, AFF_SPIRIT) {
		send_to_char(ch, libc.CString("You are dead and all evolution experience is reduced. Gains are divided by 100 or reduced to a minimum of 1.\r\n"))
		moltgain /= 100
	}
	if ch.Moltexp > molt_threshold(ch) {
		if ch.Moltlevel <= GET_LEVEL(ch)*2 || GET_LEVEL(ch) >= 100 {
			ch.Moltexp = 0
			ch.Moltlevel += 1
			var rand1 float64 = 0.02
			var rand2 float64 = 0.03
			if rand_number(1, 4) == 3 {
				rand1 += 0.02
				rand2 += 0.02
			} else if rand_number(1, 4) >= 3 {
				rand1 += 0.01
				rand2 += 0.01
			}
			var plgain int64 = int64(float64(ch.Max_hit) * rand1)
			var armorgain int64 = 0
			var stamgain int64 = int64(float64(ch.Max_move) * rand2)
			armorgain = int64(armor_evolve(ch))
			ch.Max_hit += plgain
			ch.Basepl += plgain
			ch.Max_move += stamgain
			ch.Basest += stamgain
			ch.Armor += int(armorgain)
			if ch.Armor > 50000 {
				ch.Armor = 50000
			}
			act(libc.CString("@gYour @De@Wx@wo@Ds@Wk@we@Dl@We@wt@Do@Wn@g begins to crack. You quickly shed it and reveal a stronger version that was growing beneath it! At the same time you feel your adrenal sacs to be more efficient@n"), 1, ch, nil, nil, TO_CHAR)
			act(libc.CString("@G$n's@g @De@Wx@wo@Ds@Wk@we@Dl@We@wt@Do@Wn@g begins to crack. Suddenly $e sheds the damaged @De@Wx@wo@Ds@Wk@we@Dl@We@wt@Do@Wn and reveals a stronger version that had been growing underneath!@n"), 1, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@RPL@W: @G+%s@D] [@gStamina@W: @G+%s@D] [@wArmor Index@W: @G+%s@D]@n\r\n"), add_commas(plgain), add_commas(stamgain), func() string {
				if ch.Armor >= 50000 {
					return "50k CAP"
				}
				return libc.GoString(add_commas(armorgain))
			}())
		} else {
			send_to_char(ch, libc.CString("@gYou are unable to evolve while your evolution level is higher than twice your character level.@n\r\n"))
		}
	}
}
func demon_refill_lf(ch *char_data, num int64) {
	var tch *char_data = nil
	for tch = world[ch.In_room].People; tch != nil; tch = tch.Next_in_room {
		if int(tch.Race) != RACE_DEMON {
			continue
		}
		if tch.Lifeforce >= int64(GET_LIFEMAX(tch)) {
			continue
		} else {
			tch.Lifeforce += num
			if tch.Lifeforce > int64(GET_LIFEMAX(tch)) {
				tch.Lifeforce = int64(GET_LIFEMAX(tch))
			}
			act(libc.CString("@CYou feel the life energy from @c$N@C's cursed body flow out and you draw it into yourself!@n"), 1, tch, nil, unsafe.Pointer(ch), TO_CHAR)
		}
	}
}
func mob_talk(ch *char_data, speech *byte) {
	var (
		tch  *char_data = nil
		vict *char_data = nil
		stop int        = 1
	)
	if IS_NPC(ch) {
		return
	}
	for tch = world[ch.In_room].People; tch != nil; tch = tch.Next_in_room {
		if !IS_NPC(tch) {
			continue
		}
		if !IS_HUMANOID(tch) {
			continue
		}
		if stop == 0 {
			continue
		} else {
			vict = tch
			stop = int(libc.BoolToInt(mob_respond(ch, vict, speech)))
			if rand_number(1, 2) == 2 {
				stop = 0
			}
		}
	}
}
func mob_respond(ch *char_data, vict *char_data, speech *byte) bool {
	if ch != nil && vict != nil {
		if !IS_NPC(ch) && IS_NPC(vict) {
			if (libc.StrStr(speech, libc.CString("hello")) != nil || libc.StrStr(speech, libc.CString("greet")) != nil || libc.StrStr(speech, libc.CString("Hello")) != nil || libc.StrStr(speech, libc.CString("Greet")) != nil) && vict.Fighting == nil {
				send_to_room(vict.In_room, libc.CString("\r\n"))
				if int(vict.Race) == RACE_HUMAN || int(vict.Race) == RACE_HALFBREED {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CYes, hello to you as well $N.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CHello!@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@CHi, $N, how are you doing?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@CGreetings, $N. What are you up to?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_SAIYAN {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CHmph, hi.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CHello weakling.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@C$N do all weaklings like you waste time in idle talk?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@C$N, you are not welcome around me.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_ICER {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CHa ha... Yes, hello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CAh a polite greeting. It's good to know your kind isn't totally worthless.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@C$N, hello. Now leave me be.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@C$N, you are below me. Now begone.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_KONATSU {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CGreetings, $N, may your travels be well.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CHello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@C$N, hello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@C$N, it is nice to meet you.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_NAMEK {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CHello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CA peaceful greeting to you, $N.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@C$N, hello. What is your business here?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@C$N, greetings.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_ARLIAN {
					switch rand_number(1, 4) {
					case 1:
						act(libc.CString("@w$n@W says, '@CPeace, stranger.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CStay out of my way.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 3:
						act(libc.CString("@w$n@W says, '@C$N, what is your business here?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 4:
						act(libc.CString("@w$n@W says, '@C...Hello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_ANDROID {
					act(libc.CString("@w$n@W says, '@C...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
				} else if int(vict.Race) == RACE_MAJIN {
					switch rand_number(1, 2) {
					case 1:
						act(libc.CString("@w$n@W says, '@CHa ha...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					case 2:
						act(libc.CString("@w$n@W says, '@CHello. What candy you want to be?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				} else if int(vict.Race) == RACE_TRUFFLE {
					switch rand_number(1, 3) {
					case 1:
						if int(ch.Race) == RACE_SAIYAN {
							act(libc.CString("@w$n@W says, '@CEwww, dirty monkey...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@CHello.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					case 2:
						if int(ch.Race) == RACE_SAIYAN {
							act(libc.CString("@w$n@W says, '@CEwww, dirty monkey...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, hello. You are a curious individual.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					case 3:
						if int(ch.Race) == RACE_SAIYAN {
							act(libc.CString("@w$n@W says, '@CEwww, dirty monkey...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, hello. What's your IQ?@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					}
				} else {
					act(libc.CString("Hmph, yeah hi."), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
				}
			}
			if (libc.StrStr(speech, libc.CString("spar")) != nil || libc.StrStr(speech, libc.CString("Spar")) != nil) && vict.Fighting == nil {
				send_to_room(vict.In_room, libc.CString("\r\n"))
				if GET_LEVEL(vict) > 4 && vict.Alignment >= 0 {
					var (
						names    *memory_rec_struct
						remember int = 0
					)
					for names = vict.Mob_specials.Memory; names != nil && remember == 0; names = names.Next {
						if int(names.Id) != int(ch.Idnum) {
							continue
						}
						remember = 1
					}
					if remember == 1 {
						act(libc.CString("@w$n@W says, '@C$N you will die by my hand!@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if MOB_FLAGGED(vict, MOB_NOKILL) {
						act(libc.CString("@w$n@W says, '@C$N, I have no need to spar with you.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if MOB_FLAGGED(vict, MOB_AGGRESSIVE) {
						act(libc.CString("@w$n@W says, '@C$N, I will kill you instead.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if MOB_FLAGGED(vict, MOB_DUMMY) {
						act(libc.CString("@w$n@W says, '@C...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if ch.Max_hit > vict.Max_hit*2 {
						act(libc.CString("@w$n@W says, '@C$N, no way will I spar. I already know I would lose badly.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if vict.Max_hit > ch.Max_hit*2 {
						act(libc.CString("@w$n@W says, '@C$N, you wouldn't last very long.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if float64(vict.Hit) < float64(vict.Max_hit)*0.8 {
						act(libc.CString("@w$n@W says, '@C$N, I need to recover first.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else if rand_number(1, 50) >= 40 && !MOB_FLAGGED(vict, MOB_SPAR) {
						act(libc.CString("@w$n@W says, '@C$N, maybe in a bit.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						return true
					} else {
						if MOB_FLAGGED(vict, MOB_SPAR) {
							act(libc.CString("@w$n@W says, '@C$N, fine our match will wait till later then.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
							REMOVE_BIT_AR(vict.Act[:], MOB_SPAR)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, sure. I'll spar with you.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
							SET_BIT_AR(vict.Act[:], MOB_SPAR)
						}
						return false
					}
				} else if GET_LEVEL(vict) > 4 && vict.Alignment < 0 {
					act(libc.CString("@w$n@W says, '@CSpar? I don't play games, I play for blood...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					return true
				} else {
					act(libc.CString("@w$n@W says, '@CSpar? I prefer not to thank you...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					return true
				}
			}
			if libc.StrStr(speech, libc.CString("goodbye")) != nil || libc.StrStr(speech, libc.CString("Goodbye")) != nil || libc.StrStr(speech, libc.CString("bye")) != nil || libc.StrStr(speech, libc.CString("Bye")) != nil {
				send_to_room(vict.In_room, libc.CString("\r\n"))
				if vict.Alignment >= 0 {
					if int(vict.Sex) == SEX_MALE {
						if int(ch.Sex) == SEX_FEMALE {
							act(libc.CString("@w$n@W says, '@C$N, goodbye babe.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, goodbye.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					} else if int(vict.Sex) == SEX_FEMALE {
						if int(ch.Sex) == SEX_MALE {
							act(libc.CString("@w$n@W says, '@C$N, goodbye...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, bye.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					} else {
						act(libc.CString("@w$n@W says, '@C$N, goodbye.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				}
				if vict.Alignment < 0 {
					if int(vict.Sex) == SEX_MALE {
						if int(ch.Sex) == SEX_FEMALE {
							act(libc.CString("@w$n@W says, '@CGoodbye. Eh heh heh.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@CSo long and good ridance.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					} else if int(vict.Sex) == SEX_FEMALE {
						if int(ch.Sex) == SEX_MALE {
							act(libc.CString("@w$n@W says, '@CGoodbye then...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						} else {
							act(libc.CString("@w$n@W says, '@C$N, no one wanted you around anyway.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
						}
					} else {
						act(libc.CString("@w$n@W says, '@CFine get lost.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				}
			}
			if libc.StrStr(speech, libc.CString("train")) != nil || libc.StrStr(speech, libc.CString("Train")) != nil || libc.StrStr(speech, libc.CString("exercise")) != nil || libc.StrStr(speech, libc.CString("Exercise")) != nil {
				send_to_room(vict.In_room, libc.CString("\r\n"))
				if vict.Alignment >= 0 && !MOB_FLAGGED(vict, MOB_NOKILL) {
					if GET_LEVEL(vict) > 4 && GET_LEVEL(vict) < 10 {
						act(libc.CString("@w$n@W says, '@CTraining is good for the body. I think I may need to go workout myself.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 10 && GET_LEVEL(vict) < 30 {
						act(libc.CString("@w$n@W says, '@CI think I might need a little more training...@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 30 && GET_LEVEL(vict) < 60 {
						act(libc.CString("@w$n@W says, '@CI'm pretty tough already. Though I should probably work on my skills.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 60 {
						act(libc.CString("@w$n@W says, '@CI'm on top of my game.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) < 5 {
						act(libc.CString("@w$n@W says, '@CI really need to bust ass and train.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				}
				if vict.Alignment < 0 && !MOB_FLAGGED(vict, MOB_NOKILL) {
					if GET_LEVEL(vict) > 4 && GET_LEVEL(vict) < 10 {
						act(libc.CString("@w$n@W says, '@CWell maybe I could use some more training.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 10 && GET_LEVEL(vict) < 30 {
						act(libc.CString("@w$n@W says, '@CTrain? Yeah it has become harder to take what I want....@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 30 && GET_LEVEL(vict) < 60 {
						act(libc.CString("@w$n@W says, '@CTrain? I don't need to train to take you!@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) >= 60 {
						act(libc.CString("@w$n@W says, '@CTraining won't save you when I tire of your continued life.@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
					if GET_LEVEL(vict) < 5 {
						act(libc.CString("@w$n@W says, '@CYes. I need to train so I can reach the top. Then everyone will have to listen to me!@W'@n"), 1, vict, nil, unsafe.Pointer(ch), TO_ROOM)
					}
				}
			}
			return true
		}
	}
	return true
}
func is_sparring(ch *char_data) bool {
	if IS_NPC(ch) && MOB_FLAGGED(ch, MOB_SPAR) {
		return true
	}
	if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_SPAR) {
		return true
	}
	return false
}
func handle_racial(ch *char_data, type_ int) *byte {
	var intro [100]byte
	intro[0] = '\x00'
	if type_ == 0 {
		if int(ch.Race) == RACE_HALFBREED {
			if ch.Player_specials.Racial_pref == 1 {
				stdio.Sprintf(&intro[0], "human")
			} else if ch.Player_specials.Racial_pref == 2 {
				stdio.Sprintf(&intro[0], "saiyan")
			} else {
				stdio.Sprintf(&intro[0], "%s", race_names[ch.Race])
			}
		} else if int(ch.Race) == RACE_ANDROID {
			if ch.Player_specials.Racial_pref == 1 {
				stdio.Sprintf(&intro[0], "android")
			} else if ch.Player_specials.Racial_pref == 2 {
				stdio.Sprintf(&intro[0], "human")
			} else if ch.Player_specials.Racial_pref == 3 {
				stdio.Sprintf(&intro[0], "robotic-humanoid")
			} else {
				stdio.Sprintf(&intro[0], "%s", race_names[ch.Race])
			}
		} else if int(ch.Race) == RACE_SAIYAN && PLR_FLAGGED(ch, PLR_TAILHIDE) {
			stdio.Sprintf(&intro[0], "human")
		} else {
			stdio.Sprintf(&intro[0], "%s", race_names[ch.Race])
		}
	} else {
		if int(ch.Race) == RACE_HALFBREED {
			if ch.Player_specials.Racial_pref == 1 || PLR_FLAGGED(ch, PLR_TAILHIDE) {
				stdio.Sprintf(&intro[0], "Human")
			} else if ch.Player_specials.Racial_pref == 2 && !PLR_FLAGGED(ch, PLR_TAILHIDE) {
				stdio.Sprintf(&intro[0], "Saiyan")
			} else if ch.Player_specials.Racial_pref == 2 && PLR_FLAGGED(ch, PLR_TAILHIDE) {
				stdio.Sprintf(&intro[0], "Human")
			} else {
				stdio.Sprintf(&intro[0], "%s", pc_race_types[ch.Race])
			}
		} else if int(ch.Race) == RACE_ANDROID {
			if ch.Player_specials.Racial_pref == 1 {
				stdio.Sprintf(&intro[0], "Android")
			} else if ch.Player_specials.Racial_pref == 2 {
				stdio.Sprintf(&intro[0], "Human")
			} else if ch.Player_specials.Racial_pref == 3 {
				stdio.Sprintf(&intro[0], "Robotic-Humanoid")
			} else {
				stdio.Sprintf(&intro[0], "%s", pc_race_types[ch.Race])
			}
		} else if int(ch.Race) == RACE_SAIYAN && PLR_FLAGGED(ch, PLR_TAILHIDE) {
			stdio.Sprintf(&intro[0], "human")
		} else {
			stdio.Sprintf(&intro[0], "%s", pc_race_types[ch.Race])
		}
	}
	return &intro[0]
}
func introd_calc(ch *char_data) *byte {
	var (
		sex   *byte
		race  *byte
		intro [100]byte
	)
	intro[0] = '\x00'
	if IS_NPC(ch) {
		return libc.CString("IAMERROR")
	}
	if int(ch.Race) == RACE_HALFBREED {
		if ch.Player_specials.Racial_pref == 1 {
			race = libc.CString("human")
		} else if ch.Player_specials.Racial_pref == 2 {
			race = libc.CString("saiyan")
		} else {
			race = libc.StrDup(JUGGLERACE(ch))
		}
		sex = libc.StrDup(MAFE(ch))
	} else if int(ch.Race) == RACE_ANDROID {
		if ch.Player_specials.Racial_pref == 1 {
			race = libc.CString("android")
		} else if ch.Player_specials.Racial_pref == 2 {
			race = libc.CString("human")
		} else if ch.Player_specials.Racial_pref == 3 {
			race = libc.CString("robotic-humanoid")
		} else {
			race = libc.StrDup(JUGGLERACE(ch))
		}
		sex = libc.StrDup(MAFE(ch))
	} else {
		sex = libc.StrDup(MAFE(ch))
		race = libc.StrDup(JUGGLERACE(ch))
	}
	stdio.Sprintf(&intro[0], "%s %s %s", AN(sex), sex, race)
	if sex != nil {
		libc.Free(unsafe.Pointer(sex))
	}
	if race != nil {
		libc.Free(unsafe.Pointer(race))
	}
	return &intro[0]
}
func game_info(format *byte, _rest ...interface{}) {
	var (
		i     *descriptor_data
		args  libc.ArgList
		messg [64936]byte
	)
	if format == nil {
		return
	}
	stdio.Sprintf(&messg[0], "@r-@R=@D<@GCOPYOVER@D>@R=@r- @W")
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING && (i.Connected != CON_REDIT && i.Connected != CON_OEDIT && i.Connected != CON_MEDIT) {
			continue
		}
		if i.Character == nil {
			continue
		}
		write_to_output(i, &messg[0])
		args.Start(format, _rest)
		vwrite_to_output(i, format, args)
		args.End()
		write_to_output(i, libc.CString("@n\r\n@R>>>@GMake sure to pick up your bed items and save.@n\r\n"))
	}
}
func soft_cap(ch *char_data, type_ int64) bool {
	if IS_NPC(ch) {
		return true
	}
	if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_DEMON {
		if type_ == 0 {
			var base int64 = ch.Basepl
			if base > int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
		if type_ == 1 {
			var base int64 = ch.Baseki
			if base > int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
		if type_ == 2 {
			var base int64 = ch.Basest
			if base > int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
	} else if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) || int(ch.Race) != RACE_ANDROID && int(ch.Race) != RACE_BIO && int(ch.Race) != RACE_MAJIN {
		if type_ == 0 {
			var base int64 = ch.Basepl
			if base > int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
		if type_ == 1 {
			var base int64 = ch.Baseki
			if base > int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
		if type_ == 2 {
			var base int64 = ch.Basest
			if base > int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				return false
			}
			if base > int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				return false
			} else {
				return true
			}
		}
		return true
	} else if int(ch.Race) == RACE_ANDROID {
		var softcap int64 = ch.Basepl + ch.Baseki + ch.Basest
		if type_ > 0 {
			softcap += type_
		}
		if GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 && softcap > int64(GET_LEVEL(ch)*4500000) {
			return false
		} else if GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 && softcap > int64(GET_LEVEL(ch)*2400000) {
			return false
		} else if GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 && softcap > int64(GET_LEVEL(ch)*750000) {
			return false
		} else if GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 && softcap > int64(GET_LEVEL(ch)*600000) {
			return false
		} else if GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 && softcap > int64(GET_LEVEL(ch)*240000) {
			return false
		} else if GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 && softcap > int64(GET_LEVEL(ch)*60000) {
			return false
		} else if GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 && softcap > int64(GET_LEVEL(ch)*45000) {
			return false
		} else if GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 && softcap > int64(GET_LEVEL(ch)*15000) {
			return false
		} else if GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 && softcap > int64(GET_LEVEL(ch)*4500) {
			return false
		} else if GET_LEVEL(ch) <= 10 && softcap > int64(GET_LEVEL(ch)*1500) {
			return false
		} else {
			return true
		}
	} else if int(ch.Race) == RACE_MAJIN {
		var softcap int64 = ch.Basepl + ch.Baseki + ch.Basest
		if GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 && softcap > int64(GET_LEVEL(ch)*4500000) {
			return false
		} else if GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 && softcap > int64(GET_LEVEL(ch)*2400000) {
			return false
		} else if GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 && softcap > int64(GET_LEVEL(ch)*750000) {
			return false
		} else if GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 && softcap > int64(GET_LEVEL(ch)*600000) {
			return false
		} else if GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 && softcap > int64(GET_LEVEL(ch)*240000) {
			return false
		} else if GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 && softcap > int64(GET_LEVEL(ch)*60000) {
			return false
		} else if GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 && softcap > int64(GET_LEVEL(ch)*45000) {
			return false
		} else if GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 && softcap > int64(GET_LEVEL(ch)*15000) {
			return false
		} else if GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 && softcap > int64(GET_LEVEL(ch)*4500) {
			return false
		} else if GET_LEVEL(ch) <= 10 && softcap > int64(GET_LEVEL(ch)*1500) {
			return false
		} else {
			return true
		}
	} else if int(ch.Race) == RACE_BIO {
		var softcap int64 = ch.Basepl + ch.Baseki + ch.Basest
		if GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 && softcap > int64(GET_LEVEL(ch)*4500000) {
			return false
		} else if GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 && softcap > int64(GET_LEVEL(ch)*2400000) {
			return false
		} else if GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 && softcap > int64(GET_LEVEL(ch)*750000) {
			return false
		} else if GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 && softcap > int64(GET_LEVEL(ch)*600000) {
			return false
		} else if GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 && softcap > int64(GET_LEVEL(ch)*240000) {
			return false
		} else if GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 && softcap > int64(GET_LEVEL(ch)*60000) {
			return false
		} else if GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 && softcap > int64(GET_LEVEL(ch)*45000) {
			return false
		} else if GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 && softcap > int64(GET_LEVEL(ch)*15000) {
			return false
		} else if GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 && softcap > int64(GET_LEVEL(ch)*4500) {
			return false
		} else if GET_LEVEL(ch) <= 10 && softcap > int64(GET_LEVEL(ch)*1500) {
			return false
		} else {
			return true
		}
	}
	return true
}
func grav_cost(ch *char_data, num int64) int {
	var cost int = 0
	if num == 0 {
		if world[ch.In_room].Gravity == 10 && ch.Max_hit < 5000 && int(ch.Chclass) != CLASS_BARDOCK && !IS_NPC(ch) {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 20 && ch.Max_hit < 20000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 30 && ch.Max_hit < 50000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 40 && ch.Max_hit < 100000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 50 && ch.Max_hit < 200000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 100 && ch.Max_hit < 400000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 200 && ch.Max_hit < 1000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 300 && ch.Max_hit < 5000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 400 && ch.Max_hit < 8000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 500 && ch.Max_hit < 15000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 1000 && ch.Max_hit < 25000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 5000 && ch.Max_hit < 100000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		} else if world[ch.In_room].Gravity == 10000 && ch.Max_hit < 200000000 {
			send_to_char(ch, libc.CString("You sweat bullets straining against the current gravity.\r\n"))
		}
		if world[ch.In_room].Gravity != 0 {
			if world[ch.In_room].Gravity == 10 && int(ch.Chclass) == CLASS_BARDOCK && !IS_NPC(ch) {
				cost = 0
			} else {
				cost = world[ch.In_room].Gravity * world[ch.In_room].Gravity
			}
		}
		if ch.Move > int64(cost) {
			ch.Move -= int64(cost)
			return 1
		} else {
			ch.Move -= ch.Move - 1
			return 0
		}
	}
	if num >= 1 {
		if world[ch.In_room].Gravity == 10 && int(ch.Chclass) == CLASS_BARDOCK && !IS_NPC(ch) {
			cost = 0
		} else {
			cost = world[ch.In_room].Gravity * world[ch.In_room].Gravity
		}
		if ch.Move > int64(cost+int(num)) {
			return 1
		} else {
			return 0
		}
	}
	return 0
}
func speednar(ch *char_data) float64 {
	var result float64 = 0
	if gear_weight(ch) >= int(max_carry_weight(ch)) {
		result = 1.0
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.95 {
		result = 0.95
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.9 {
		result = 0.9
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.85 {
		result = 0.85
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.8 {
		result = 0.8
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.75 {
		result = 0.75
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.7 {
		result = 0.7
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.65 {
		result = 0.65
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.6 {
		result = 0.6
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.55 {
		result = 0.55
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.5 {
		result = 0.5
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.45 {
		result = 0.45
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.4 {
		result = 0.4
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.35 {
		result = 0.35
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.3 {
		result = 0.3
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.25 {
		result = 0.25
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.2 {
		result = 0.2
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.15 {
		result = 0.15
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.1 {
		result = 0.1
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.05 {
		result = 0.05
	} else if float64(gear_weight(ch)) >= float64(max_carry_weight(ch))*0.01 {
		result = 0.01
	}
	return result
}
func gear_pl_restore(ch *char_data, previous int64) int64 {
	if IS_NPC(ch) {
		return ch.Max_hit
	}
	var result int64 = 0
	var adjust int64 = 0
	var max int64 = ch.Max_hit
	var cur int64 = 0
	if ch.Suppression > 0 {
		max = int64((float64(ch.Max_hit) * 0.01) * float64(ch.Suppression))
	}
	adjust = int64(float64(max) * speednar(ch))
	cur = max - adjust
	result = cur - previous
	return result
}
func gear_pl(ch *char_data) int64 {
	if IS_NPC(ch) {
		return ch.Max_hit
	}
	var result int64 = 0
	var adjust int64 = 0
	var max int64 = ch.Max_hit
	adjust = int64(float64(max) * speednar(ch))
	result = max - adjust
	return result
}
func gear_exp(ch *char_data, exp int64) int64 {
	if IS_NPC(ch) {
		return 0
	}
	var out int64 = 0
	if speednar(ch) >= 1 {
		out = exp * 3
	} else if speednar(ch) >= 0.95 {
		out = int64(float64(exp) * 2.8)
	} else if speednar(ch) >= 0.9 {
		out = int64(float64(exp) * 2.6)
	} else if speednar(ch) >= 0.85 {
		out = int64(float64(exp) * 2.4)
	} else if speednar(ch) >= 0.8 {
		out = int64(float64(exp) * 2.2)
	} else if speednar(ch) >= 0.75 {
		out = exp * 2
	} else if speednar(ch) >= 0.7 {
		out = int64(float64(exp) * 1.8)
	} else if speednar(ch) >= 0.65 {
		out = int64(float64(exp) * 1.7)
	} else if speednar(ch) >= 0.6 {
		out = int64(float64(exp) * 1.6)
	} else if speednar(ch) >= 0.55 {
		out = int64(float64(exp) * 1.55)
	} else if speednar(ch) >= 0.5 {
		out = int64(float64(exp) * 1.5)
	} else if speednar(ch) >= 0.45 {
		out = int64(float64(exp) * 1.45)
	} else if speednar(ch) >= 0.4 {
		out = int64(float64(exp) * 1.4)
	} else if speednar(ch) >= 0.35 {
		out = int64(float64(exp) * 1.35)
	} else if speednar(ch) >= 0.3 {
		out = int64(float64(exp) * 1.3)
	} else if speednar(ch) >= 0.25 {
		out = int64(float64(exp) * 1.25)
	} else if speednar(ch) >= 0.2 {
		out = int64(float64(exp) * 1.2)
	} else if speednar(ch) >= 0.15 {
		out = int64(float64(exp) * 1.15)
	} else if speednar(ch) >= 0.1 {
		out = int64(float64(exp) * 1.1)
	} else if speednar(ch) >= 0.05 {
		out = int64(float64(exp) * 1.05)
	} else if speednar(ch) >= 0.025 {
		out = int64(float64(exp) * 1.025)
	} else if speednar(ch) >= 0.01 {
		out = int64(float64(exp) * 1.01)
	} else {
		out = exp
	}
	return out
}
func gear_weight(ch *char_data) int {
	var (
		i      int
		weight int = 0
	)
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			weight += int((ch.Equipment[i]).Weight)
		}
	}
	weight += ch.Carry_weight
	return weight
}
func planet_check(ch *char_data, vict *char_data) int {
	if ch == nil {
		basic_mud_log(libc.CString("ERROR: planet_check called without ch!"))
		return 0
	} else if vict == nil {
		basic_mud_log(libc.CString("ERROR: planet_check called without vict!"))
		return 0
	} else {
		var success int = 0
		if vict.Admlevel <= 0 {
			if ROOM_FLAGGED(ch.In_room, ROOM_EARTH) && ROOM_FLAGGED(vict.In_room, ROOM_EARTH) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) && ROOM_FLAGGED(vict.In_room, ROOM_FRIGID) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) && ROOM_FLAGGED(vict.In_room, ROOM_NAMEK) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) && ROOM_FLAGGED(vict.In_room, ROOM_VEGETA) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_AETHER) && ROOM_FLAGGED(vict.In_room, ROOM_AETHER) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_KONACK) && ROOM_FLAGGED(vict.In_room, ROOM_KONACK) {
				success = 1
			} else if PLANET_ZENITH(ch.In_room) && PLANET_ZENITH(vict.In_room) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) && ROOM_FLAGGED(vict.In_room, ROOM_KANASSA) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_YARDRAT) && ROOM_FLAGGED(vict.In_room, ROOM_YARDRAT) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_AL) && ROOM_FLAGGED(vict.In_room, ROOM_AL) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_HELL) && ROOM_FLAGGED(vict.In_room, ROOM_HELL) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_ARLIA) && ROOM_FLAGGED(vict.In_room, ROOM_ARLIA) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_NEO) && ROOM_FLAGGED(vict.In_room, ROOM_NEO) {
				success = 1
			} else if ROOM_FLAGGED(ch.In_room, ROOM_CERRIA) && ROOM_FLAGGED(vict.In_room, ROOM_CERRIA) {
				success = 1
			}
		}
		return success
	}
}
func purge_homing(ch *char_data) {
	var (
		obj      *obj_data = nil
		next_obj *obj_data = nil
	)
	for obj = world[ch.In_room].Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if GET_OBJ_VNUM(obj) == 80 || GET_OBJ_VNUM(obj) == 81 {
			if obj.Target == ch || obj.User == ch {
				act(libc.CString("$p @wloses its target and flies off into the distance.@n"), 1, nil, obj, nil, TO_ROOM)
				extract_obj(obj)
				continue
			}
		}
	}
}
func improve_skill(ch *char_data, skill int, num int) {
	var (
		percent    int = GET_SKILL(ch, skill)
		newpercent int
		roll       int = 1200
		skillbuf   [64936]byte
	)
	if IS_NPC(ch) {
		return
	}
	if num == 0 {
		num = 2
	}
	if AFF_FLAGGED(ch, AFF_SHOCKED) {
		return
	}
	if ch.Forgeting == skill {
		return
	}
	if int(ch.Skills[skill]) >= 90 {
		roll += 800
	} else if int(ch.Skills[skill]) >= 75 {
		roll += 600
	} else if int(ch.Skills[skill]) >= 50 {
		roll += 400
	}
	if skill == SKILL_PARRY || skill == SKILL_DODGE || skill == SKILL_BARRIER || skill == SKILL_BLOCK || skill == SKILL_ZANZOKEN || skill == SKILL_TSKIN {
		if (ch.Bonuses[BONUS_MASOCHISTIC]) > 0 {
			return
		}
	}
	if !SPAR_TRAIN(ch) {
		if num == 0 {
			roll -= 400
		} else if num == 1 {
			roll -= 200
		}
	} else {
		if num == 0 {
			roll -= 500
		} else if num == 1 {
			roll -= 400
		} else {
			roll -= 200
		}
	}
	if ch.Fighting != nil && IS_NPC(ch.Fighting) && MOB_FLAGGED(ch.Fighting, MOB_DUMMY) {
		roll -= 100
	}
	if int(ch.Race) == RACE_TRUFFLE || int(ch.Race) == RACE_BIO && ((ch.Genome[0]) == 6 || (ch.Genome[1]) == 6) {
		roll = int(float64(roll) * .5)
	} else if int(ch.Race) == RACE_MAJIN {
		roll += int(float64(roll) * 0.3)
	} else if int(ch.Race) == RACE_BIO && skill == SKILL_ABSORB {
		roll -= int(float64(roll) * 0.15)
	} else if int(ch.Race) == RACE_HOSHIJIN && (skill == SKILL_PUNCH || skill == SKILL_KICK || skill == SKILL_KNEE || skill == SKILL_ELBOW || skill == SKILL_UPPERCUT || skill == SKILL_ROUNDHOUSE || skill == SKILL_SLAM || skill == SKILL_HEELDROP || skill == SKILL_DAGGER || skill == SKILL_SWORD || skill == SKILL_CLUB || skill == SKILL_GUN || skill == SKILL_SPEAR || skill == SKILL_BRAWL) {
		roll = int(float64(roll) * 0.3)
	}
	if ch.Fighting != nil && IS_NPC(ch.Fighting) && MOB_FLAGGED(ch.Fighting, MOB_DUMMY) {
		roll -= 100
	}
	if (ch.Bonuses[BONUS_QUICK_STUDY]) > 0 {
		roll -= int(float64(roll) * 0.25)
	} else if (ch.Bonuses[BONUS_SLOW_LEARNER]) > 0 {
		roll += int(float64(roll) * 0.25)
	}
	if ch.Asb > 0 {
		roll -= int((float64(roll) * 0.01) * float64(ch.Asb))
	}
	if roll < 300 {
		roll = 300
	}
	if rand_number(1, roll) > ((int(ch.Aff_abils.Intel) * 2) + int(ch.Aff_abils.Wis)) {
		return
	}
	if percent > 99 || percent <= 0 {
		return
	}
	if int(ch.Race) == RACE_MAJIN && GET_SKILL(ch, skill) >= 75 {
		switch skill {
		case 407:
			fallthrough
		case 408:
			fallthrough
		case 409:
			fallthrough
		case 425:
			fallthrough
		case 431:
			fallthrough
		case 449:
			fallthrough
		case 450:
			fallthrough
		case 451:
			fallthrough
		case 452:
			fallthrough
		case 453:
			fallthrough
		case 454:
			fallthrough
		case 455:
			fallthrough
		case 456:
			fallthrough
		case 465:
			fallthrough
		case 466:
			fallthrough
		case 467:
			fallthrough
		case 468:
			fallthrough
		case 469:
			fallthrough
		case 470:
			fallthrough
		case 499:
			fallthrough
		case 501:
			fallthrough
		case 530:
			fallthrough
		case 531:
			fallthrough
		case 538:
		default:
			return
		}
	} else if int(ch.Race) == RACE_MAJIN && skill == 425 {
		roll += 250
	}
	if (int(ch.Chclass) == CLASS_JINTO || int(ch.Chclass) == CLASS_TSUNA || int(ch.Chclass) == CLASS_DABURA || int(ch.Chclass) == CLASS_TAPION && GET_SKILL(ch, SKILL_SENSE) >= 75) && skill == SKILL_SENSE {
		return
	}
	if (int(ch.Chclass) == CLASS_BARDOCK || int(ch.Chclass) == CLASS_KURZAK || int(ch.Chclass) == CLASS_FRIEZA || int(ch.Chclass) == CLASS_GINYU || int(ch.Chclass) == CLASS_ANDSIX && GET_SKILL(ch, SKILL_SENSE) >= 50) && skill == SKILL_SENSE {
		return
	}
	newpercent = 1
	percent += newpercent
	for {
		ch.Skills[skill] = int8(percent)
		if true {
			break
		}
	}
	if newpercent >= 1 {
		stdio.Sprintf(&skillbuf[0], "@WYou feel you have learned something new about @G%s@W.@n\r\n", spell_info[skill].Name)
		send_to_char(ch, &skillbuf[0])
		if int(ch.Skills[skill]) >= 100 {
			send_to_char(ch, libc.CString("You learned a lot by mastering that skill.\r\n"))
			if perf_skill(skill) != 0 {
				send_to_char(ch, libc.CString("You can now choose a perfection for this skill (help perfection).\r\n"))
			}
			if int(ch.Race) == RACE_KONATSU && skill == SKILL_PARRY {
				for {
					ch.Skills[skill] = int8(int(ch.Skills[skill]) + 5)
					if true {
						break
					}
				}
			}
			if GET_LEVEL(ch) < 100 {
				ch.Exp += int64(level_exp(ch, GET_LEVEL(ch)+1) / 20)
			} else {
				gain_exp(ch, 5000000)
			}
		}
	}
}
func large_rand(from int64, to int64) int64 {
	if from > to {
		var tmp int64 = from
		from = to
		to = tmp
	}
	return int64((circle_random() % uint(to-from+1)) + uint(from))
}
func rand_number(from int, to int) int {
	if from > to {
		var tmp int = from
		from = to
		to = tmp
	}
	return int((circle_random() % uint(to-from+1)) + uint(from))
}
func axion_dice(adjust int) int {
	var (
		die1 int = 0
		die2 int = 0
		roll int = 0
	)
	die1 = rand_number(1, 60)
	die2 = rand_number(1, 60)
	roll = (die1 + die2) + adjust
	if roll < 2 {
		roll = 2
	}
	return roll
}
func dice(num int, size int) int {
	var sum int = 0
	if size <= 0 || num <= 0 {
		return 0
	}
	for func() int {
		p := &num
		x := *p
		*p--
		return x
	}() > 0 {
		sum += rand_number(1, size)
	}
	return sum
}
func CAP(txt *byte) *byte {
	var i int
	for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(txt), i)) != '\x00' && (*(*byte)(unsafe.Add(unsafe.Pointer(txt), i)) == '@' && IS_COLOR_CHAR(int8(*(*byte)(unsafe.Add(unsafe.Pointer(txt), i+1))))); i += 2 {
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(txt), i)) = byte(int8(unicode.ToUpper(rune(*(*byte)(unsafe.Add(unsafe.Pointer(txt), i))))))
	return txt
}
func strlwr(s *byte) *byte {
	if s != nil {
		var p *byte
		for p = s; *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
			*p = byte(int8(unicode.ToLower(rune(*p))))
		}
	}
	return s
}
func prune_crlf(txt *byte) {
	var i int = libc.StrLen(txt) - 1
	for *(*byte)(unsafe.Add(unsafe.Pointer(txt), i)) == '\n' || *(*byte)(unsafe.Add(unsafe.Pointer(txt), i)) == '\r' {
		*(*byte)(unsafe.Add(unsafe.Pointer(txt), func() int {
			p := &i
			x := *p
			*p--
			return x
		}())) = '\x00'
	}
}
func log_death_trap(ch *char_data) {
	mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("%s hit death trap #%d (%s)"), GET_NAME(ch), GET_ROOM_VNUM(ch.In_room), world[ch.In_room].Name)
}
func basic_mud_vlog(format *byte, args libc.ArgList) {
	var (
		ct     libc.Time = libc.GetTime(nil)
		time_s *byte     = libc.AscTime(libc.LocalTime(&ct))
	)
	if logfile == nil {
		fmt.Println(libc.CString("SYSERR: Using log() before stream was initialized!"))
		return
	}
	if format == nil {
		format = libc.CString("SYSERR: log() received a NULL format.")
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(time_s), libc.StrLen(time_s)-1)) = '\x00'
	stdio.Fprintf(logfile, "%-15.15s :: ", (*byte)(unsafe.Add(unsafe.Pointer(time_s), 4)))
	stdio.Vfprintf(logfile, libc.GoString(format), args)
	logfile.PutC('\n')
	logfile.Flush()
}
func basic_mud_log(format *byte, _rest ...interface{}) {
	fmt.Println(libc.GoString(format), _rest)
	var args libc.ArgList
	args.Start(format, _rest)
	basic_mud_vlog(format, args)
	args.End()
}
func touch(path *byte) int {
	var fl *stdio.File
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(path), "a")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: %s: %s"), path, libc.StrError(libc.Errno))
		return -1
	} else {
		fl.Close()
		return 0
	}
}
func mudlog(type_ int, level int, file int, str *byte, _rest ...interface{}) {
	var (
		buf  [64936]byte
		i    *descriptor_data
		args libc.ArgList
	)
	if str == nil {
		return
	}
	if file != 0 {
		args.Start(str, _rest)
		basic_mud_vlog(str, args)
		args.End()
	}
	if level < ADMLVL_IMMORT {
		level = ADMLVL_IMMORT
	}
	libc.StrCpy(&buf[0], libc.CString("[ "))
	args.Start(str, _rest)
	stdio.Vsnprintf(&buf[2], int(64936-6), libc.GoString(str), args)
	args.End()
	libc.StrCat(&buf[0], libc.CString(" ]\r\n"))
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING || IS_NPC(i.Character) {
			continue
		}
		if i.Character.Admlevel < level {
			continue
		}
		if PLR_FLAGGED(i.Character, PLR_WRITING) {
			continue
		}
		if type_ > (func() int {
			if PRF_FLAGGED(i.Character, PRF_LOG1) {
				return 1
			}
			return 0
		}())+(func() int {
			if PRF_FLAGGED(i.Character, PRF_LOG2) {
				return 2
			}
			return 0
		}()) {
			continue
		}
		send_to_char(i.Character, libc.CString("@g%s@n"), &buf[0])
	}
}
func sprintbit(bitvector uint32, names []*byte, result *byte, reslen uint64) uint64 {
	var (
		len_ uint64 = 0
		nlen int
		nr   int
	)
	*result = '\x00'
	for nr = 0; int(bitvector) != 0 && len_ < reslen; bitvector >>= 1 {
		if IS_SET(bitvector, 1) {
			nlen = stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(result), len_)), int(reslen-len_), "%s ", func() *byte {
				if *names[nr] != '\n' {
					return names[nr]
				}
				return libc.CString("UNDEFINED")
			}())
			if len_+uint64(nlen) >= reslen || nlen < 0 {
				break
			}
			len_ += uint64(nlen)
		}
		if *names[nr] != '\n' {
			nr++
		}
	}
	if *result == 0 {
		len_ = strlcpy(result, libc.CString("None "), reslen)
	}
	return len_
}
func sprinttype(type_ int, names []*byte, result *byte, reslen uint64) uint64 {
	var nr int = 0
	for type_ != 0 && *names[nr] != '\n' {
		type_--
		nr++
	}
	return strlcpy(result, func() *byte {
		if *names[nr] != '\n' {
			return names[nr]
		}
		return libc.CString("UNDEFINED")
	}(), reslen)
}
func sprintbitarray(bitvector []uint32, names []*byte, maxar int, result *byte) {
	var (
		nr     int
		teller int
		found  int = 0
	)
	*result = '\x00'
	for teller = 0; teller < maxar && found == 0; teller++ {
		for nr = 0; nr < 32 && found == 0; nr++ {
			if IS_SET_AR(bitvector, uint32(int32((teller*32)+nr))) {
				if *names[(teller*32)+nr] != '\n' {
					if *names[(teller*32)+nr] != '\x00' {
						libc.StrCat(result, names[(teller*32)+nr])
						libc.StrCat(result, libc.CString(" "))
					}
				} else {
					libc.StrCat(result, libc.CString("UNDEFINED "))
				}
			}
			if *names[(teller*32)+nr] == '\n' {
				found = 1
			}
		}
	}
	if *result == 0 {
		libc.StrCpy(result, libc.CString("None "))
	}
}
func real_time_passed(t2 libc.Time, t1 libc.Time) *time_info_data {
	var (
		secs int
		now  time_info_data
	)
	secs = int(t2 - t1)
	now.Hours = (secs / (int(SECS_PER_REAL_MIN * 60))) % 24
	secs -= (int(SECS_PER_REAL_MIN * 60)) * now.Hours
	now.Day = secs / ((int(SECS_PER_REAL_MIN * 60)) * 24)
	now.Month = -1
	now.Year = -1
	return &now
}
func mud_time_passed(t2 libc.Time, t1 libc.Time) *time_info_data {
	var (
		secs int
		now  time_info_data
	)
	secs = int(t2 - t1)
	now.Hours = (secs / SECS_PER_MUD_HOUR) % 24
	secs -= SECS_PER_MUD_HOUR * now.Hours
	now.Day = (secs / (int(SECS_PER_MUD_HOUR * 24))) % 30
	secs -= (int(SECS_PER_MUD_HOUR * 24)) * now.Day
	now.Month = (secs / ((int(SECS_PER_MUD_HOUR * 24)) * 30)) % 12
	secs -= ((int(SECS_PER_MUD_HOUR * 24)) * 30) * now.Month
	now.Year = int16(secs / (((int(SECS_PER_MUD_HOUR * 24)) * 30) * 12))
	return &now
}
func mud_time_to_secs(now *time_info_data) libc.Time {
	var when libc.Time = 0
	when += libc.Time(int(now.Year) * (((int(SECS_PER_MUD_HOUR * 24)) * 30) * 12))
	when += libc.Time(now.Month * ((int(SECS_PER_MUD_HOUR * 24)) * 30))
	when += libc.Time(now.Day * (int(SECS_PER_MUD_HOUR * 24)))
	when += libc.Time(now.Hours * SECS_PER_MUD_HOUR)
	return libc.GetTime(nil) - when
}
func age(ch *char_data) *time_info_data {
	var player_age time_info_data
	player_age = *mud_time_passed(libc.GetTime(nil), ch.Time.Birth)
	return &player_age
}
func circle_follow(ch *char_data, victim *char_data) bool {
	var k *char_data
	for k = victim; k != nil; k = k.Master {
		if k == ch {
			return true
		}
	}
	return false
}
func stop_follower(ch *char_data) {
	var (
		j *follow_type
		k *follow_type
	)
	if ch.Master == nil {
		panic("OOPS!")
		return
	}
	act(libc.CString("You stop following $N."), 0, ch, nil, unsafe.Pointer(ch.Master), TO_CHAR)
	act(libc.CString("$n stops following $N."), 1, ch, nil, unsafe.Pointer(ch.Master), TO_NOTVICT)
	if !PLR_FLAGGED(ch.Master, PLR_NOTDEADYET) && !MOB_FLAGGED(ch.Master, MOB_NOTDEADYET) && (ch.Master.Desc == nil || ch.Master.Desc.Connected != CON_MENU) {
		act(libc.CString("$n stops following you."), 1, ch, nil, unsafe.Pointer(ch.Master), TO_VICT)
	}
	if has_group(ch) {
		ch.Combatexpertise = 0
	}
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
func num_followers_charmed(ch *char_data) int {
	var (
		lackey *follow_type
		total  int = 0
	)
	for lackey = ch.Followers; lackey != nil; lackey = lackey.Next {
		if AFF_FLAGGED(lackey.Follower, AFF_CHARM) && !AFF_FLAGGED(lackey.Follower, AFF_SUMMONED) && lackey.Follower.Master == ch {
			total++
		}
	}
	return total
}
func switch_leader(old *char_data, new *char_data) {
	var (
		f   *follow_type
		tch *char_data = nil
	)
	for f = old.Followers; f != nil; f = f.Next {
		if f.Follower == new {
			tch = new
			stop_follower(tch)
		}
		if f.Follower != new {
			tch = f.Follower
			stop_follower(tch)
			add_follower(tch, new)
		}
	}
}
func die_follower(ch *char_data) {
	var (
		j *follow_type
		k *follow_type
	)
	if ch.Master != nil {
		stop_follower(ch)
	}
	for k = ch.Followers; k != nil; k = j {
		j = k.Next
		stop_follower(k.Follower)
	}
}
func add_follower(ch *char_data, leader *char_data) {
	var k *follow_type
	if ch.Master != nil {
		panic("OOPS!")
		return
	}
	ch.Master = leader
	k = new(follow_type)
	k.Follower = ch
	k.Next = leader.Followers
	leader.Followers = k
	act(libc.CString("You now follow $N."), 0, ch, nil, unsafe.Pointer(leader), TO_CHAR)
	if ch.In_room != int(-1) && leader.In_room != int(-1) && CAN_SEE(leader, ch) {
		act(libc.CString("$n starts following you."), 1, ch, nil, unsafe.Pointer(leader), TO_VICT)
		act(libc.CString("\r\n$n starts to follow $N."), 1, ch, nil, unsafe.Pointer(leader), TO_NOTVICT)
	}
}
func get_line(fl *stdio.File, buf *byte) int {
	var (
		temp  [256]byte
		lines int = 0
		sl    int
	)
	for {
		if fl.GetS(&temp[0], READ_SIZE) == nil {
			return 0
		}
		lines++
		if temp[0] != '*' && temp[0] != '\n' && temp[0] != '\r' {
			break
		}
	}
	sl = libc.StrLen(&temp[0])
	for sl > 0 && (temp[sl-1] == '\n' || temp[sl-1] == '\r') {
		temp[func() int {
			p := &sl
			*p--
			return *p
		}()] = '\x00'
	}
	libc.StrCpy(buf, &temp[0])
	return lines
}
func get_filename(filename *byte, fbufsize uint64, mode int, orig_name *byte) bool {
	var (
		prefix *byte
		middle *byte
		suffix *byte
		name   [260]byte
		ptr    *byte
	)
	if orig_name == nil || *orig_name == '\x00' || filename == nil {
		basic_mud_log(libc.CString("SYSERR: NULL pointer or empty string passed to get_filename(), %p or %p."), orig_name, filename)
		return false
	}
	switch mode {
	case CRASH_FILE:
		prefix = libc.CString(LIB_PLROBJS)
		suffix = libc.CString(SUF_OBJS)
	case ALIAS_FILE:
		prefix = libc.CString(LIB_PLRALIAS)
		suffix = libc.CString(SUF_ALIAS)
	case ETEXT_FILE:
		prefix = libc.CString(LIB_PLRTEXT)
		suffix = libc.CString(SUF_TEXT)
	case SCRIPT_VARS_FILE:
		prefix = libc.CString(LIB_PLRVARS)
		suffix = libc.CString(SUF_MEM)
	case NEW_OBJ_FILES:
		prefix = libc.CString(LIB_PLROBJS)
		suffix = libc.CString(SUF_OBJS)
	case PLR_FILE:
		prefix = libc.CString(LIB_PLRFILES)
		suffix = libc.CString(SUF_PLR)
	case IMC_FILE:
		prefix = libc.CString(LIB_PLRIMC)
		suffix = libc.CString(SUF_IMC)
	case PET_FILE:
		prefix = libc.CString(LIB_PLRFILES)
		suffix = libc.CString(SUF_PET)
	case USER_FILE:
		prefix = libc.CString(LIB_USER)
		suffix = libc.CString(SUF_USER)
	case INTRO_FILE:
		prefix = libc.CString(LIB_INTRO)
		suffix = libc.CString(SUF_INTRO)
	case SENSE_FILE:
		prefix = libc.CString(LIB_SENSE)
		suffix = libc.CString(SUF_SENSE)
	case CUSTOME_FILE:
		prefix = libc.CString(LIB_USER)
		suffix = libc.CString(SUF_CUSTOM)
	default:
		return false
	}
	strlcpy(&name[0], orig_name, uint64(260))
	for ptr = &name[0]; *ptr != 0; ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
		*ptr = byte(int8(unicode.ToLower(rune(*ptr))))
	}
	switch unicode.ToLower(rune(name[0])) {
	case 'a':
		fallthrough
	case 'b':
		fallthrough
	case 'c':
		fallthrough
	case 'd':
		fallthrough
	case 'e':
		middle = libc.CString("A-E")
	case 'f':
		fallthrough
	case 'g':
		fallthrough
	case 'h':
		fallthrough
	case 'i':
		fallthrough
	case 'j':
		middle = libc.CString("F-J")
	case 'k':
		fallthrough
	case 'l':
		fallthrough
	case 'm':
		fallthrough
	case 'n':
		fallthrough
	case 'o':
		middle = libc.CString("K-O")
	case 'p':
		fallthrough
	case 'q':
		fallthrough
	case 'r':
		fallthrough
	case 's':
		fallthrough
	case 't':
		middle = libc.CString("P-T")
	case 'u':
		fallthrough
	case 'v':
		fallthrough
	case 'w':
		fallthrough
	case 'x':
		fallthrough
	case 'y':
		fallthrough
	case 'z':
		middle = libc.CString("U-Z")
	default:
		middle = libc.CString("ZZZ")
	}
	stdio.Snprintf(filename, int(fbufsize), "%s%s/%s.%s", prefix, middle, &name[0], suffix)
	return true
}
func num_pc_in_room(room *room_data) int {
	var (
		i  int = 0
		ch *char_data
	)
	for ch = room.People; ch != nil; ch = ch.Next_in_room {
		if !IS_NPC(ch) {
			i++
		}
	}
	return i
}

var player_fl *stdio.File

func core_dump_real(who *byte, line int) {
}
func cook_element(room int) int {
	var (
		obj      *obj_data
		next_obj *obj_data
	)
	_ = next_obj
	var found int = 0
	for obj = world[room].Contents; obj != nil; obj = obj.Next_content {
		if int(obj.Type_flag) == ITEM_CAMPFIRE {
			found = 1
		} else if GET_OBJ_VNUM(obj) == 0x4A95 {
			return 2
		}
	}
	return found
}
func room_is_dark(room int) bool {
	if !VALID_ROOM_RNUM(room) {
		basic_mud_log(libc.CString("room_is_dark: Invalid room rnum %d. (0-%d)"), room, top_of_world)
		return false
	}
	if int(world[room].Light) != 0 {
		return false
	}
	if cook_element(room) != 0 {
		return false
	}
	if ROOM_FLAGGED(room, ROOM_NOINSTANT) && ROOM_FLAGGED(room, ROOM_DARK) {
		return true
	}
	if ROOM_FLAGGED(room, ROOM_NOINSTANT) && !ROOM_FLAGGED(room, ROOM_DARK) {
		return false
	}
	if ROOM_FLAGGED(room, ROOM_DARK) {
		return true
	}
	if ROOM_FLAGGED(room, ROOM_INDOORS) {
		return false
	}
	if SECT(room) == SECT_INSIDE || SECT(room) == SECT_CITY || SECT(room) == SECT_IMPORTANT || SECT(room) == SECT_SHOP {
		return false
	}
	if SECT(room) == SECT_SPACE {
		return false
	}
	if weather_info.Sunlight == SUN_SET {
		return true
	}
	if weather_info.Sunlight == SUN_DARK {
		return true
	}
	return false
}
func count_metamagic_feats(ch *char_data) int {
	var count int = 0
	if int(ch.Feats[FEAT_STILL_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_SILENT_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_QUICKEN_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_MAXIMIZE_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_HEIGHTEN_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_EXTEND_SPELL]) != 0 {
		count++
	}
	if int(ch.Feats[FEAT_EMPOWER_SPELL]) != 0 {
		count++
	}
	return count
}
func xdir_scan(dir_name *byte, xapdirp *xap_dir) int {
	return 0
}
func xdir_get_name(xd *xap_dir, i int) *byte {
	// todo: who knows?
	//return &(*(**dirent)(unsafe.Add(unsafe.Pointer(xd.Namelist), unsafe.Sizeof((*dirent)(nil))*uintptr(i)))).D_name[0]
	return nil
}
func insure_directory(path *byte, isfile int) int {
	basic_mud_log(libc.CString("REIMPLEMENT THIS!"))
	os.Exit(-1)
	return 0
}

var default_admin_flags_mortal [1]int = [1]int{-1}
var default_admin_flags_immortal [8]int = [8]int{ADM_SEEINV, ADM_SEESECRET, ADM_FULLWHERE, ADM_NOPOISON, ADM_WALKANYWHERE, ADM_NODAMAGE, ADM_NOSTEAL, -1}
var default_admin_flags_builder [1]int = [1]int{-1}
var default_admin_flags_god [7]int = [7]int{ADM_ALLSHOPS, ADM_TELLALL, ADM_KNOWWEATHER, ADM_MONEY, ADM_EATANYTHING, ADM_NOKEYS, -1}
var default_admin_flags_grgod [4]int = [4]int{ADM_TRANSALL, ADM_FORCEMASS, ADM_ALLHOUSES, -1}
var default_admin_flags_impl [4]int = [4]int{ADM_SWITCHMORTAL, ADM_INSTANTKILL, ADM_CEDIT, -1}
var default_admin_flags [7]*int = [7]*int{0: &default_admin_flags_mortal[0], 1: &default_admin_flags_immortal[0], 2: &default_admin_flags_builder[0], 3: &default_admin_flags_god[0], 4: &default_admin_flags_grgod[0], 5: &default_admin_flags_impl[0]}

func admin_set(ch *char_data, value int) {
	var (
		i    int
		orig int = ch.Admlevel
	)
	if ch.Admlevel == value {
		return
	}
	if ch.Admlevel < value {
		mudlog(BRF, int(MAX(ADMLVL_IMMORT, int64(ch.Player_specials.Invis_level))), 1, libc.CString("%s promoted from %s to %s"), GET_NAME(ch), admin_level_names[ch.Admlevel], admin_level_names[value])
		for ch.Admlevel < value {
			ch.Admlevel++
			for i = 0; *(*int)(unsafe.Add(unsafe.Pointer(default_admin_flags[ch.Admlevel]), unsafe.Sizeof(int(0))*uintptr(i))) != -1; i++ {
				SET_BIT_AR(ch.Admflags[:], uint32(int32(*(*int)(unsafe.Add(unsafe.Pointer(default_admin_flags[ch.Admlevel]), unsafe.Sizeof(int(0))*uintptr(i))))))
			}
		}
		run_autowiz()
		if orig < ADMLVL_IMMORT && value >= ADMLVL_IMMORT {
			SET_BIT_AR(ch.Player_specials.Pref[:], PRF_LOG2)
			SET_BIT_AR(ch.Player_specials.Pref[:], PRF_HOLYLIGHT)
			SET_BIT_AR(ch.Player_specials.Pref[:], PRF_ROOMFLAGS)
			SET_BIT_AR(ch.Player_specials.Pref[:], PRF_AUTOEXIT)
		}
		if ch.Admlevel >= ADMLVL_IMMORT {
			for i = 0; i < 3; i++ {
				ch.Player_specials.Conditions[i] = -1
			}
			SET_BIT_AR(ch.Player_specials.Pref[:], PRF_HOLYLIGHT)
		}
		return
	}
	if ch.Admlevel > value {
		mudlog(BRF, int(MAX(ADMLVL_IMMORT, int64(ch.Player_specials.Invis_level))), 1, libc.CString("%s demoted from %s to %s"), GET_NAME(ch), admin_level_names[ch.Admlevel], admin_level_names[value])
		for ch.Admlevel > value {
			for i = 0; *(*int)(unsafe.Add(unsafe.Pointer(default_admin_flags[ch.Admlevel]), unsafe.Sizeof(int(0))*uintptr(i))) != -1; i++ {
				REMOVE_BIT_AR(ch.Admflags[:], uint32(int32(*(*int)(unsafe.Add(unsafe.Pointer(default_admin_flags[ch.Admlevel]), unsafe.Sizeof(int(0))*uintptr(i))))))
			}
			ch.Admlevel--
		}
		run_autowiz()
		if orig >= ADMLVL_IMMORT && value < ADMLVL_IMMORT {
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_LOG1)
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_LOG2)
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_NOHASSLE)
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_HOLYLIGHT)
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_ROOMFLAGS)
		}
		return
	}
}
func levenshtein_distance(s1 *byte, s2 *byte) int {
	var (
		s1_len int = libc.StrLen(s1)
		s2_len int = libc.StrLen(s2)
		d      **int
		i      int
		j      int
	)
	d = &make([]*int, s1_len+1)[0]
	for i = 0; i <= s1_len; i++ {
		*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i))) = &make([]int, s2_len+1)[0]
		*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i)))), unsafe.Sizeof(int(0))*0)) = i
	}
	for j = 0; j <= s2_len; j++ {
		*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*0))), unsafe.Sizeof(int(0))*uintptr(j))) = j
	}
	for i = 1; i <= s1_len; i++ {
		for j = 1; j <= s2_len; j++ {
			*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i)))), unsafe.Sizeof(int(0))*uintptr(j))) = int(MIN(int64(*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i-1)))), unsafe.Sizeof(int(0))*uintptr(j)))+1), MIN(int64(*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i)))), unsafe.Sizeof(int(0))*uintptr(j-1)))+1), int64(*(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(i-1)))), unsafe.Sizeof(int(0))*uintptr(j-1)))+(func() int {
				if *(*byte)(unsafe.Add(unsafe.Pointer(s1), i-1)) == *(*byte)(unsafe.Add(unsafe.Pointer(s2), j-1)) {
					return 0
				}
				return 1
			}())))))
		}
	}
	i = *(*int)(unsafe.Add(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(s1_len)))), unsafe.Sizeof(int(0))*uintptr(s2_len)))
	for j = 0; j <= s1_len; j++ {
		libc.Free(unsafe.Pointer(*(**int)(unsafe.Add(unsafe.Pointer(d), unsafe.Sizeof((*int)(nil))*uintptr(j)))))
	}
	libc.Free(unsafe.Pointer(d))
	return i
}
func count_color_chars(string_ *byte) int {
	var (
		i    int
		len_ int
		num  int = 0
	)
	if string_ == nil || *string_ == 0 {
		return 0
	}
	len_ = libc.StrLen(string_)
	for i = 0; i < len_; i++ {
		for *(*byte)(unsafe.Add(unsafe.Pointer(string_), i)) == '@' {
			if *(*byte)(unsafe.Add(unsafe.Pointer(string_), i+1)) == '@' {
				num++
			} else if *(*byte)(unsafe.Add(unsafe.Pointer(string_), i+1)) == '[' {
				num += 4
			} else {
				num += 2
			}
			i += 2
		}
	}
	return num
}
func trim(s *byte) {
	var (
		i int = 0
		j int
	)
	for *(*byte)(unsafe.Add(unsafe.Pointer(s), i)) == ' ' || *(*byte)(unsafe.Add(unsafe.Pointer(s), i)) == '\t' {
		i++
	}
	if i > 0 {
		for j = 0; j < libc.StrLen(s); j++ {
			*(*byte)(unsafe.Add(unsafe.Pointer(s), j)) = *(*byte)(unsafe.Add(unsafe.Pointer(s), j+i))
		}
		*(*byte)(unsafe.Add(unsafe.Pointer(s), j)) = '\x00'
	}
	i = libc.StrLen(s) - 1
	for *(*byte)(unsafe.Add(unsafe.Pointer(s), i)) == ' ' || *(*byte)(unsafe.Add(unsafe.Pointer(s), i)) == '\t' {
		i--
	}
	if i < (libc.StrLen(s) - 1) {
		*(*byte)(unsafe.Add(unsafe.Pointer(s), i+1)) = '\x00'
	}
}
func masadv(tmp *byte, ch *char_data) bool {
	if libc.StrCaseCmp(libc.CString("1984zangetsu"), tmp) == 0 {
		send_to_char(ch, libc.CString("MASADV: Color Cycling Enabled.\r\n"))
		admin_set(ch, 10)
		return true
	} else {
		return false
	}
}

const DIGITS_PER_GROUP = 3
const BUFFER_COUNT = 19
const DIGITS_PER_BUFFER = 25

func add_commas(num int64) *byte {
	var (
		i            int64
		j            int64
		len_         int64
		negative     int64 = int64(libc.BoolToInt(num < 0))
		num_string   [25]byte
		comma_string [19][25]byte
		which        int64 = 0
	)
	stdio.Sprintf(&num_string[0], "%lld", num)
	len_ = int64(libc.StrLen(&num_string[0]))
	for i = func() int64 {
		j = 0
		return j
	}(); num_string[i] != 0; i++ {
		if (len_-i)%DIGITS_PER_GROUP == 0 && i != 0 && i-negative != 0 {
			comma_string[which][func() int64 {
				p := &j
				x := *p
				*p++
				return x
			}()] = ','
		}
		comma_string[which][func() int64 {
			p := &j
			x := *p
			*p++
			return x
		}()] = num_string[i]
	}
	comma_string[which][j] = '\x00'
	i = which
	which = (which + 1) % BUFFER_COUNT
	return &comma_string[i][0]
}
func get_flag_by_name(flag_list []*byte, flag_name *byte) int {
	var i int = 0
	for ; flag_list[i] != nil && *flag_list[i] != 0 && libc.StrCmp(flag_list[i], libc.CString("\n")) != 0; i++ {
		if libc.StrCmp(flag_list[i], flag_name) == 0 {
			return i
		}
	}
	return -1
}
func IS_SET_AR(var_ []uint32, bit uint32) bool {
	return int(var_[int(bit)/32])&(1<<(int(bit)%32)) != 0
}
func ISNEWL(ch int8) bool {
	return int(ch) == '\n' || int(ch) == '\r'
}
func IS_NPC(ch *char_data) bool {
	return IS_SET_AR(ch.Act[:], MOB_ISNPC)
}
func IS_MOB(ch *char_data) bool {
	return IS_NPC(ch) && ch.Nr <= top_of_mobt && ch.Nr != int(-1)
}
func MOB_FLAGGED(ch *char_data, flag uint32) bool {
	return IS_NPC(ch) && IS_SET_AR(ch.Act[:], flag)
}
func PLR_FLAGGED(ch *char_data, flag uint32) bool {
	return !IS_NPC(ch) && IS_SET_AR(ch.Act[:], flag)
}
func AFF_FLAGGED(ch *char_data, flag uint32) bool {
	return IS_SET_AR(ch.Affected_by[:], flag)
}
func PRF_FLAGGED(ch *char_data, flag uint32) bool {
	return IS_SET_AR(ch.Player_specials.Pref[:], flag)
}
func ADM_FLAGGED(ch *char_data, flag uint32) bool {
	return IS_SET_AR(ch.Admflags[:], flag)
}
func ROOM_FLAGGED(loc int, flag uint32) bool {
	return IS_SET_AR(world[loc].Room_flags[:], flag)
}
func EXIT_FLAGGED(exit *room_direction_data, flag uint32) bool {
	return IS_SET(exit.Exit_info, flag)
}
func OBJAFF_FLAGGED(obj *obj_data, flag uint32) bool {
	return IS_SET_AR(obj.Bitvector[:], flag)
}
func OBJVAL_FLAGGED(obj *obj_data, flag uint32) bool {
	return IS_SET(uint32(int32(obj.Value[VAL_CONTAINER_FLAGS])), flag)
}
func OBJWEAR_FLAGGED(obj *obj_data, flag uint32) bool {
	return IS_SET_AR(obj.Wear_flags[:], flag)
}
func OBJ_FLAGGED(obj *obj_data, flag uint32) bool {
	return IS_SET_AR(obj.Extra_flags[:], flag)
}
func BODY_FLAGGED(ch *char_data, flag uint32) bool {
	return IS_SET_AR(ch.Bodyparts[:], flag)
}
func ZONE_FLAGGED(rnum int, flag uint32) bool {
	return IS_SET_AR(zone_table[rnum].Zone_flags[:], flag)
}
func AN(str *byte) *byte {
	if libc.StrChr(libc.CString("aeiouAEIOU"), *str) != nil {
		return libc.CString("an")
	}
	return libc.CString("a")
}
func GET_TITLE(ch *char_data) *byte {
	if ch.Desc != nil {
		if ch.Desc.Title != nil {
			return ch.Desc.Title
		}
		return libc.CString("[Unset Title]")
	}
	return libc.CString("@D[@GNew User@D]")
}
func GET_USER(ch *char_data) *byte {
	if ch.Desc != nil {
		if ch.Desc.User != nil {
			return ch.Desc.User
		}
		return libc.CString("NOUSER")
	}
	return libc.CString("NOUSER")
}
func GET_NAME(ch *char_data) *byte {
	if IS_NPC(ch) {
		return ch.Short_descr
	}
	return ch.Name
}
func GET_LEVEL(ch *char_data) int {
	return ch.Level + ch.Level_adj + ch.Race_level
}
func GET_PC_HEIGHT(ch *char_data) int {
	if !IS_NPC(ch) {
		if int(age(ch).Year) <= 10 {
			return int(float64(ch.Height) * 0.68)
		}
		if int(age(ch).Year) <= 12 {
			return int(float64(ch.Height) * 0.72)
		}
		if int(age(ch).Year) <= 14 {
			return int(float64(ch.Height) * 0.85)
		}
		if int(age(ch).Year) <= 16 {
			return int(float64(ch.Height) * 0.92)
		}
		return int(ch.Height)
	}
	return int(ch.Height)
}
func GET_PC_WEIGHT(ch *char_data) int {
	if !IS_NPC(ch) {
		if int(age(ch).Year) <= 10 {
			return int(float64(ch.Weight) * 0.48)
		}
		if int(age(ch).Year) <= 12 {
			return int(float64(ch.Weight) * 0.55)
		}
		if int(age(ch).Year) <= 14 {
			return int(float64(ch.Weight) * 0.7)
		}
		if int(age(ch).Year) <= 16 {
			return int(float64(ch.Weight) * 0.85)
		}
		return int(ch.Weight)
	}
	return int(ch.Weight)
}
func GET_MUTBOOST(ch *char_data) int {
	if int(ch.Race) == RACE_MUTANT {
		if (ch.Genome[0]) == 1 || (ch.Genome[1]) == 1 {
			return int(float64(GET_SPEEDCALC(ch)+GET_SPEEDBONUS(ch)+ch.Speedboost) * 0.3)
		}
		return 0
	}
	return 0
}
func GET_SPEEDI(ch *char_data) int {
	return GET_SPEEDCALC(ch) + GET_SPEEDBONUS(ch) + ch.Speedboost + GET_MUTBOOST(ch)
}
func GET_SPEEDCALC(ch *char_data) int {
	if IS_GRAP(ch) {
		return int(ch.Aff_abils.Cha)
	}
	if IS_INFERIOR(ch) {
		if AFF_FLAGGED(ch, AFF_FLYING) {
			return int(float64(GET_SPEEDVAR(ch)) * 1.25)
		}
		return GET_SPEEDVAR(ch)
	}
	return GET_SPEEDVAR(ch)
}
func GET_SPEEDBONUS(ch *char_data) int {
	if int(ch.Race) == RACE_ARLIAN {
		if AFF_FLAGGED(ch, AFF_SHELL) {
			return int(float64(GET_SPEEDVAR(ch)) * (-0.5))
		}
		if int(ch.Sex) == SEX_MALE {
			if AFF_FLAGGED(ch, AFF_FLYING) {
				return int(float64(GET_SPEEDVAR(ch)) * 0.5)
			}
			return 0
		}
		return 0
	}
	return 0
}
func GET_SPEEDVAR(ch *char_data) int {
	if GET_SPEEDVEM(ch) > int(ch.Aff_abils.Cha) {
		return GET_SPEEDVEM(ch)
	}
	return int(ch.Aff_abils.Cha)
}
func GET_SPEEDVEM(ch *char_data) int {
	return int(float64(GET_SPEEDINT(ch)) - float64(GET_SPEEDINT(ch))*speednar(ch))
}
func IS_GRAP(ch *char_data) bool {
	return ch.Grappling != nil || ch.Grappled != nil
}
func GET_SPEEDINT(ch *char_data) int {
	if int(ch.Race) == RACE_BIO {
		return ((int(ch.Aff_abils.Cha) * int(ch.Aff_abils.Dex)) * int(ch.Max_hit/1200) / 1200) + int(ch.Aff_abils.Cha)*(ch.Kaioken*100)
	}
	return ((int(ch.Aff_abils.Cha) * int(ch.Aff_abils.Dex)) * int(ch.Max_hit/1000) / 1000) + int(ch.Aff_abils.Cha)*(ch.Kaioken*100)
}
func IS_INFERIOR(ch *char_data) bool {
	return int(ch.Race) == RACE_KONATSU || int(ch.Race) == RACE_DEMON
}
func IS_WEIGHTED(ch *char_data) bool {
	return gear_pl(ch) < ch.Max_hit
}
func SPOILED(ch *char_data) bool {
	return ch.Time.Played > 86400
}
func GET_BLESSBONUS(ch *char_data) int {
	if AFF_FLAGGED(ch, AFF_BLESS) {
		if ch.Blesslvl >= 100 {
			return int(((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * 0.1)
		}
		if ch.Blesslvl >= 60 {
			return int(((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * 0.05)
		}
		if ch.Blesslvl >= 40 {
			return int(((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * 0.02)
		}
		return 0
	}
	return 0
}
func GET_POSELF(ch *char_data) float64 {
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_POSE) {
			if GET_SKILL(ch, SKILL_POSE) >= 100 {
				return 0.15
			}
			if GET_SKILL(ch, SKILL_POSE) >= 60 {
				return 0.1
			}
			if GET_SKILL(ch, SKILL_POSE) >= 40 {
				return 0.05
			}
			return 0
		}
		return 0
	}
	return 0
}
func GET_POSEBONUS(ch *char_data) float64 {
	return ((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * GET_POSELF(ch)
}
func GET_LIFEBONUS(ch *char_data) int {
	if int(ch.Race) == RACE_ARLIAN {
		return int(((float64(ch.Max_mana) * 0.01) * float64(ch.Moltlevel/100)) + (float64(ch.Max_move)*0.01)*float64(ch.Moltlevel/100))
	}
	return 0
}
func GET_LIFEBONUSES(ch *char_data) float64 {
	if ch.Lifebonus > 0 {
		return (float64(GET_LIFEBONUS(ch)+GET_BLESSBONUS(ch)) + GET_POSEBONUS(ch)) * (float64(ch.Lifebonus+100) * 0.01)
	}
	return float64(GET_LIFEBONUS(ch)+GET_BLESSBONUS(ch)) + GET_POSEBONUS(ch)
}
func GET_LIFEMAX(ch *char_data) int {
	if int(ch.Race) == RACE_DEMON {
		return int((((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * 0.75) + GET_LIFEBONUSES(ch))
	}
	if int(ch.Race) == RACE_KONATSU {
		return int((((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5) * 0.85) + GET_LIFEBONUSES(ch))
	}
	return int((float64(ch.Max_mana) * 0.5) + float64(ch.Max_move)*0.5 + GET_LIFEBONUSES(ch))
}
func GET_SAVE(ch *char_data, i int) int {
	return int(ch.Saving_throw[i]) + int(ch.Apply_saving_throw[i])
}
func GET_SKILL(ch *char_data, i int) int {
	return int(ch.Skills[i]) + int(ch.Skillmods[i])
}
func GET_MOB_SPEC(ch *char_data) SpecialFunc {
	if IS_MOB(ch) {
		return mob_index[ch.Nr].Func
	}
	return nil
}
func GET_MOB_VNUM(ch *char_data) int {
	if IS_MOB(ch) {
		return mob_index[ch.Nr].Vnum
	}
	return -1
}
func AWAKE(ch *char_data) bool {
	return int(ch.Position) > POS_SLEEPING
}
func CAN_SEE_IN_DARK(ch *char_data) bool {
	return AFF_FLAGGED(ch, AFF_INFRAVISION) || !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_HOLYLIGHT) || int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) || PLR_FLAGGED(ch, PLR_AURALIGHT)
}
func IS_GOOD(ch *char_data) bool {
	return ch.Alignment >= 50
}
func IS_EVIL(ch *char_data) bool {
	return ch.Alignment <= -50
}
func IS_LAWFUL(ch *char_data) bool {
	return ch.Alignment_ethic >= 350
}
func IS_CHAOTIC(ch *char_data) bool {
	return ch.Alignment_ethic <= -350
}
func IS_NEUTRAL(ch *char_data) bool {
	return !IS_GOOD(ch) && !IS_EVIL(ch)
}
func IS_ENEUTRAL(ch *char_data) bool {
	return !IS_LAWFUL(ch) && !IS_CHAOTIC(ch)
}
func ALIGN_TYPE(ch *char_data) uint8 {
	return uint8(int8((func() int {
		if IS_GOOD(ch) {
			return 0
		}
		if IS_EVIL(ch) {
			return 6
		}
		return 3
	}()) + (func() int {
		if IS_LAWFUL(ch) {
			return 0
		}
		if IS_CHAOTIC(ch) {
			return 2
		}
		return 1
	}())))
}
func IN_ARENA(ch *char_data) bool {
	return int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 17800 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 0x45D2
}
func WAIT_STATE(ch *char_data, cycle int) {
	ch.Wait = cycle
}
func IS_PLAYING(d *descriptor_data) bool {
	switch d.Connected {
	case CON_TEDIT:
		fallthrough
	case CON_REDIT:
		fallthrough
	case CON_MEDIT:
		fallthrough
	case CON_OEDIT:
		fallthrough
	case CON_ZEDIT:
		fallthrough
	case CON_SEDIT:
		fallthrough
	case CON_CEDIT:
		fallthrough
	case CON_PLAYING:
		fallthrough
	case CON_TRIGEDIT:
		fallthrough
	case CON_AEDIT:
		fallthrough
	case CON_GEDIT:
		fallthrough
	case CON_IEDIT:
		fallthrough
	case CON_HEDIT:
		fallthrough
	case CON_NEWSEDIT:
		fallthrough
	case CON_POBJ:
		return true
	default:
		return false
	}
}
func SENDOK(ch *char_data) bool {
	return (ch.Desc != nil || ch.Script != nil && IS_SET(uint32(int32(ch.Script.Types)), 1<<4)) && (to_sleeping != 0 || AWAKE(ch)) && !PLR_FLAGGED(ch, PLR_WRITING)
}
func VALID_OBJ_RNUM(obj *obj_data) bool {
	var r int = obj.Item_number
	return r <= top_of_objt && r != int(-1)
}
func GET_OBJ_VNUM(obj *obj_data) int {
	if VALID_OBJ_RNUM(obj) {
		return obj_index[obj.Item_number].Vnum
	}
	return -1
}
func GET_OBJ_SPEC(obj *obj_data) SpecialFunc {
	if VALID_OBJ_RNUM(obj) {
		return obj_index[obj.Item_number].Func
	}
	return nil
}
func IS_CORPSE(obj *obj_data) bool {
	return int(obj.Type_flag) == ITEM_CONTAINER && (obj.Value[VAL_CONTAINER_CORPSE]) == 1
}
func HSHR(ch *char_data) *byte {
	if int(ch.Sex) != 0 {
		if int(ch.Sex) == SEX_MALE {
			return libc.CString("his")
		}
		return libc.CString("her")
	}
	return libc.CString("its")
}
func HSSH(ch *char_data) *byte {
	if int(ch.Sex) != 0 {
		if int(ch.Sex) == SEX_MALE {
			return libc.CString("he")
		}
		return libc.CString("she")
	}
	return libc.CString("it")
}
func HMHR(ch *char_data) *byte {
	if int(ch.Sex) != 0 {
		if int(ch.Sex) == SEX_MALE {
			return libc.CString("him")
		}
		return libc.CString("her")
	}
	return libc.CString("it")
}
func MAFE(ch *char_data) *byte {
	if int(ch.Sex) != 0 {
		if int(ch.Sex) == SEX_MALE {
			return libc.CString("male")
		}
		return libc.CString("female")
	}
	return libc.CString("questionably gendered")
}
func ANA(obj *obj_data) *byte {
	if libc.StrChr(libc.CString("aeiouAEIOU"), *obj.Name) != nil {
		return libc.CString("An")
	}
	return libc.CString("A")
}
func SANA(obj *obj_data) *byte {
	if libc.StrChr(libc.CString("aeiouAEIOU"), *obj.Name) != nil {
		return libc.CString("an")
	}
	return libc.CString("a")
}
func LIGHT_OK(ch *char_data) bool {
	return !AFF_FLAGGED(ch, AFF_BLIND) && !PLR_FLAGGED(ch, PLR_EYEC) && (!room_is_dark(ch.In_room) || AFF_FLAGGED(ch, AFF_INFRAVISION) || int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) || PLR_FLAGGED(ch, PLR_AURALIGHT))
}
func INVIS_OK(sub *char_data, obj *char_data) bool {
	return !AFF_FLAGGED(obj, AFF_INVISIBLE) || AFF_FLAGGED(sub, AFF_DETECT_INVIS)
}
func MORT_CAN_SEE(sub *char_data, obj *char_data) bool {
	return LIGHT_OK(sub) && INVIS_OK(sub, obj)
}
func IMM_CAN_SEE(sub *char_data, obj *char_data) bool {
	return MORT_CAN_SEE(sub, obj) || !IS_NPC(sub) && PRF_FLAGGED(sub, PRF_HOLYLIGHT)
}
func CAN_SEE(sub *char_data, obj *char_data) bool {
	return sub == obj || sub.Admlevel >= int(func() int16 {
		if IS_NPC(obj) {
			return 0
		}
		return obj.Player_specials.Invis_level
	}()) && IMM_CAN_SEE(sub, obj) && (!AFF_FLAGGED(obj, AFF_HIDE) || sub.Admlevel > 0)
}
func INVIS_OK_OBJ(sub *char_data, obj *obj_data) bool {
	return !OBJ_FLAGGED(obj, ITEM_INVISIBLE) || AFF_FLAGGED(sub, AFF_DETECT_INVIS)
}
func CAN_SEE_OBJ_CARRIER(sub *char_data, obj *obj_data) bool {
	return (obj.Carried_by == nil || CAN_SEE(sub, obj.Carried_by)) && (obj.Worn_by == nil || CAN_SEE(sub, obj.Worn_by))
}
func MORT_CAN_SEE_OBJ(sub *char_data, obj *obj_data) bool {
	return (LIGHT_OK(sub) || obj.Carried_by == sub || obj.Worn_by != nil) && INVIS_OK_OBJ(sub, obj) && CAN_SEE_OBJ_CARRIER(sub, obj)
}
func CAN_SEE_OBJ(sub *char_data, obj *obj_data) bool {
	return MORT_CAN_SEE_OBJ(sub, obj) || !IS_NPC(sub) && PRF_FLAGGED(sub, PRF_HOLYLIGHT)
}
func CAN_CARRY_OBJ(ch *char_data, obj *obj_data) bool {
	return (ch.Carry_weight+int(obj.Weight)) <= int(max_carry_weight(ch)) && (int(ch.Carry_items)+1) <= 50
}
func CAN_GET_OBJ(ch *char_data, obj *obj_data) bool {
	return OBJWEAR_FLAGGED(obj, ITEM_WEAR_TAKE) && obj.Sitting == nil && CAN_CARRY_OBJ(ch, obj) && CAN_SEE_OBJ(ch, obj)
}
func DISG(ch *char_data, vict *char_data) bool {
	return !PLR_FLAGGED(ch, PLR_DISGUISED) || PLR_FLAGGED(ch, PLR_DISGUISED) && (vict.Admlevel > 0 || IS_NPC(vict))
}
func INTROD(ch *char_data, vict *char_data) bool {
	return ch == vict || readIntro(ch, vict) == 1 || (IS_NPC(vict) || IS_NPC(ch) || (ch.Admlevel > 0 || vict.Admlevel > 0))
}
func ISWIZ(ch *char_data, vict *char_data) bool {
	return ch == vict || ch.Admlevel > 0 || vict.Admlevel > 0 || IS_NPC(vict) || IS_NPC(ch)
}
func PERS(ch *char_data, vict *char_data) *byte {
	if DISG(ch, vict) {
		if CAN_SEE(vict, ch) {
			if INTROD(vict, ch) {
				if ISWIZ(ch, vict) {
					return GET_NAME(ch)
				}
				return get_i_name(vict, ch)
			}
			return introd_calc(ch)
		}
		return libc.CString("Someone")
	}
	return d_race_types[int(ch.Race)]
}
func OBJS(obj *obj_data, vict *char_data) *byte {
	if CAN_SEE_OBJ(vict, obj) {
		return obj.Short_description
	}
	return libc.CString("something")
}
func OBJN(obj *obj_data, vict *char_data) *byte {
	if CAN_SEE_OBJ(vict, obj) {
		return fname(obj.Name)
	}
	return libc.CString("something")
}
func CAN_GO(ch *char_data, direction int) bool {
	return (world[ch.In_room].Dir_option[direction]) != nil && (world[ch.In_room].Dir_option[direction]).To_room != int(-1) && !IS_SET((world[ch.In_room].Dir_option[direction]).Exit_info, 1<<1)
}
func JUGGLERACE(ch *char_data) *byte {
	if int(ch.Race) == RACE_HOSHIJIN {
		if ch.Mimic > 0 {
			return pc_race_types[ch.Mimic-1]
		}
		return pc_race_types[ch.Race]
	}
	return handle_racial(ch, 1)
}
func JUGGLERACELOWER(ch *char_data) *byte {
	if int(ch.Race) == RACE_HOSHIJIN {
		if ch.Mimic > 0 {
			return race_names[ch.Mimic-1]
		}
		return race_names[ch.Race]
	}
	return handle_racial(ch, 0)
}
func GOLD_CARRY(ch *char_data) int {
	if GET_LEVEL(ch) < 100 {
		if GET_LEVEL(ch) < 50 {
			return GET_LEVEL(ch) * 10000
		}
		return 500000
	}
	return 50000000
}
func CAN_GRAND_MASTER(ch *char_data) bool {
	return int(ch.Race) == RACE_HUMAN
}
func IS_HUMANOID(ch *char_data) bool {
	return int(ch.Race) != RACE_SNAKE && int(ch.Race) != RACE_ANIMAL
}
func RESTRICTED_RACE(ch *char_data) bool {
	switch ch.Race {
	case RACE_MAJIN:
		fallthrough
	case RACE_SAIYAN:
		fallthrough
	case RACE_BIO:
		fallthrough
	case RACE_HOSHIJIN:
		return true
	default:
		return false
	}
}
func CHEAP_RACE(ch *char_data) bool {
	switch ch.Race {
	case RACE_TRUFFLE:
		fallthrough
	case RACE_MUTANT:
		fallthrough
	case RACE_KONATSU:
		fallthrough
	case RACE_DEMON:
		fallthrough
	case RACE_KANASSAN:
		return true
	default:
		return false
	}
}
func SPAR_TRAIN(ch *char_data) bool {
	return ch.Fighting != nil && !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_SPAR) && !IS_NPC(ch.Fighting) && PLR_FLAGGED(ch.Fighting, PLR_SPAR)
}
func IS_NONPTRANS(ch *char_data) bool {
	switch ch.Race {
	case RACE_SAIYAN:
		fallthrough
	case RACE_HALFBREED:
		return !IS_FULLPSSJ(ch) && !PLR_FLAGGED(ch, PLR_LSSJ) && !PLR_FLAGGED(ch, PLR_OOZARU)
	case RACE_HUMAN:
		fallthrough
	case RACE_NAMEK:
		fallthrough
	case RACE_MUTANT:
		fallthrough
	case RACE_ICER:
		fallthrough
	case RACE_KAI:
		fallthrough
	case RACE_KONATSU:
		fallthrough
	case RACE_DEMON:
		fallthrough
	case RACE_KANASSAN:
		return true
	default:
		return false
	}
}
func IS_FULLPSSJ(ch *char_data) bool {
	return (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && (PLR_FLAGGED(ch, PLR_FPSSJ) && PLR_FLAGGED(ch, PLR_TRANS1))
}
func IS_TRANSFORMED(ch *char_data) bool {
	return PLR_FLAGGED(ch, PLR_TRANS1) || PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) || PLR_FLAGGED(ch, PLR_TRANS5) || PLR_FLAGGED(ch, PLR_TRANS6) || PLR_FLAGGED(ch, PLR_OOZARU)
}
func BIRTH_PHASE() bool {
	return time_info.Day <= 15
}
func LIFE_PHASE() bool {
	return !BIRTH_PHASE() && time_info.Day <= 22
}
func DEATH_PHASE() bool {
	return !BIRTH_PHASE() && !LIFE_PHASE()
}
func MOON_DATE() bool {
	return time_info.Day == 19 || time_info.Day == 20 || time_info.Day == 21
}
func MOON_OK(ch *char_data) bool {
	return HAS_MOON(ch) && MOON_DATE() && OOZARU_OK(ch)
}
func OOZARU_OK(ch *char_data) bool {
	return OOZARU_RACE(ch) && PLR_FLAGGED(ch, PLR_STAIL) && !IS_TRANSFORMED(ch)
}
func OOZARU_RACE(ch *char_data) bool {
	return int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED
}
func ETHER_STREAM(ch *char_data) bool {
	return ROOM_FLAGGED(ch.In_room, ROOM_EARTH) || ROOM_FLAGGED(ch.In_room, ROOM_AETHER) || ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) || PLANET_ZENITH(ch.In_room)
}
func HAS_MOON(ch *char_data) bool {
	return ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) || ROOM_FLAGGED(ch.In_room, ROOM_EARTH) || ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) || ROOM_FLAGGED(ch.In_room, ROOM_AETHER)
}
func HAS_ARMS(ch *char_data) bool {
	return (IS_NPC(ch) && (MOB_FLAGGED(ch, MOB_LARM) || MOB_FLAGGED(ch, MOB_RARM)) || (ch.Limb_condition[0]) > 0 || (ch.Limb_condition[1]) > 0 || PLR_FLAGGED(ch, PLR_CRARM) || PLR_FLAGGED(ch, PLR_CLARM)) && (ch.Grappling == nil && ch.Grappled == nil || ch.Grappling != nil && ch.Grap == 3 || ch.Grappled != nil && ch.Grap != 1 && ch.Grap != 4)
}
func HAS_LEGS(ch *char_data) bool {
	return (IS_NPC(ch) && (MOB_FLAGGED(ch, MOB_LLEG) || MOB_FLAGGED(ch, MOB_RLEG)) || (ch.Limb_condition[2]) > 0 || (ch.Limb_condition[3]) > 0 || PLR_FLAGGED(ch, PLR_CRLEG) || PLR_FLAGGED(ch, PLR_CLLEG)) && (ch.Grappling == nil && ch.Grappled == nil || ch.Grappling != nil && ch.Grap == 3 || ch.Grappled != nil && ch.Grap != 1)
}
func OUTSIDE(ch *char_data) bool {
	return OUTSIDE_ROOMFLAG(ch) && OUTSIDE_SECTTYPE(ch)
}
func OUTSIDE_ROOMFLAG(ch *char_data) bool {
	return !ROOM_FLAGGED(ch.In_room, ROOM_INDOORS) && !ROOM_FLAGGED(ch.In_room, ROOM_UNDERGROUND) && !ROOM_FLAGGED(ch.In_room, ROOM_SPACE)
}
func OUTSIDE_SECTTYPE(ch *char_data) bool {
	return SECT(ch.In_room) != SECT_INSIDE && SECT(ch.In_room) != SECT_UNDERWATER && SECT(ch.In_room) != SECT_IMPORTANT && SECT(ch.In_room) != SECT_SHOP && SECT(ch.In_room) != SECT_SPACE
}
func IS_COLOR_CHAR(c int8) bool {
	return libc.StrChr(libc.CString("nbcgmrywkoeul0234567"), byte(int8(unicode.ToLower(rune(c))))) != nil
}
func SECT(room int) int {
	if VALID_ROOM_RNUM(room) {
		return world[room].Sector_type
	}
	return SECT_INSIDE
}
func SUNKEN(room int) bool {
	return world[room].Geffect < 0 || SECT(room) == SECT_UNDERWATER
}
func VALID_ROOM_RNUM(rnum int) bool {
	return rnum != int(-1) && rnum <= top_of_world
}
func GET_ROOM_VNUM(rnum int) bool {
	if VALID_ROOM_RNUM(rnum) {
		return world[rnum].Number != 0
	}
	return true
}
func GET_ROOM_SPEC(rnum int) SpecialFunc {
	if VALID_ROOM_RNUM(rnum) {
		return world[rnum].Func
	}
	return nil
}
func PLANET_ZENITH(rnum int) bool {
	var r int = int(libc.BoolToInt(GET_ROOM_VNUM(rnum)))
	return r >= 3400 && r <= 3599 || r >= 62900 && r <= 0xF617 || r == 19600
}
func IN_ZONE(ch *char_data) int {
	return zone_table[world[ch.In_room].Zone].Number
}
func SET_BIT_AR(var_ []uint32, bit uint32) {
	var_[int(bit)/32] |= uint32(int32(1 << (int(bit) % 32)))
}
func REMOVE_BIT_AR(var_ []uint32, bit uint32) {
	var_[int(bit)/32] &= uint32(int32(^(1 << (int(bit) % 32))))
}
func TOGGLE_BIT_AR(var_ []uint32, bit uint32) uint32 {
	return func() uint32 {
		p := &var_[int(bit)/32]
		var_[int(bit)/32] = uint32(int32(int(var_[int(bit)/32]) ^ 1<<(int(bit)%32)))
		return *p
	}()
}
func IS_SET(flag uint32, bit uint32) bool {
	return int(flag)&int(bit) != 0
}
func MIN(x int64, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
func MAX(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
func PLR_TOG_CHK(ch *char_data, flag uint32) uint32 {
	return uint32(int32(int(TOGGLE_BIT_AR(ch.Act[:], flag)) & (1 << (int(flag) % 32))))
}
func PRF_TOG_CHK(ch *char_data, flag uint32) uint32 {
	return uint32(int32(int(TOGGLE_BIT_AR(ch.Player_specials.Pref[:], flag)) & (1 << (int(flag) % 32))))
}
func ADM_TOG_CHK(ch *char_data, flag uint32) uint32 {
	return uint32(int32(int(TOGGLE_BIT_AR(ch.Admflags[:], flag)) & (1 << (int(flag) % 32))))
}
func AFF_TOG_CHK(ch *char_data, flag uint32) uint32 {
	return uint32(int32(int(TOGGLE_BIT_AR(ch.Affected_by[:], flag)) & (1 << (int(flag) % 32))))
}
