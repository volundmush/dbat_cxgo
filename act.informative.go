package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unicode"
	"unsafe"
)

const MAX_PORTAL_TYPES = 6
const SHOW_OBJ_LONG = 0
const SHOW_OBJ_SHORT = 1
const SHOW_OBJ_ACTION = 2
const HIST_LENGTH = 100
const WHO_FORMAT = "Usage: who [minlev[-maxlev]] [-k] [-n name] [-q] [-r] [-s] [-z]\r\n"
const USERS_FORMAT = "format: users [-l minlevel[-maxlevel]] [-n name] [-h host] [-o] [-p]\r\n"

func do_evolve(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Race) != RACE_ARLIAN || IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are not an arlian!\r\n"))
		return
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	var plcost int64 = int64(GET_LEVEL(ch))
	var stcost int64 = int64(GET_LEVEL(ch))
	var kicost int64 = int64(GET_LEVEL(ch))
	plcost += int64(float64(molt_threshold(ch))*0.65 + float64(ch.Max_hit)*0.15)
	kicost += int64((float64(molt_threshold(ch)) * 0.5) + float64(ch.Max_mana)*0.22)
	stcost += int64((float64(molt_threshold(ch)) * 0.5) + float64(ch.Max_move)*0.15)
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@D-=@YConvert Evolution Points To What?@D=-@n\r\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------@n\r\n"))
		send_to_char(ch, libc.CString("@CPowerlevel  @D: @Y%s @Wpts\r\n"), add_commas(plcost))
		send_to_char(ch, libc.CString("@CKi          @D: @Y%s @Wpts\r\n"), add_commas(kicost))
		send_to_char(ch, libc.CString("@CStamina     @D: @Y%s @Wpts\r\n"), add_commas(stcost))
		send_to_char(ch, libc.CString("@D[@Y%s @Wpts currently@D]@n\r\n"), add_commas(ch.Moltexp))
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("powerlevel")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("pl")) == 0 {
		if plcost > molt_threshold(ch) {
			send_to_char(ch, libc.CString("You need a few more evolution levels before you can start upgrading powerlevel.\r\n"))
			return
		} else if ch.Moltexp < plcost {
			send_to_char(ch, libc.CString("You do not have enough evolution experience.\r\n"))
			return
		} else {
			var plgain int64 = int64(float64(ch.Basepl) * 0.01)
			if plgain <= 0 {
				plgain = int64(rand_number(1, 5))
			} else {
				plgain = int64(rand_number(int(plgain), int(float64(plgain)*0.5)))
			}
			ch.Hit += plgain
			ch.Max_hit += plgain
			ch.Basepl += plgain
			ch.Moltexp -= plcost
			send_to_char(ch, libc.CString("Your body evolves to make better use of the way it is now, and you feel that your body has strengthened. @D[@RPL@D: @Y+%s@D]@n\r\n"), add_commas(plgain))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("ki")) == 0 {
		if kicost > molt_threshold(ch) {
			send_to_char(ch, libc.CString("You need a few more evolution levels before you can start upgrading ki.\r\n"))
			return
		} else if ch.Moltexp < kicost {
			send_to_char(ch, libc.CString("You do not have enough evolution experience.\r\n"))
			return
		} else {
			var kigain int64 = int64(float64(ch.Baseki) * 0.01)
			if kigain <= 0 {
				kigain = int64(rand_number(1, 5))
			} else {
				kigain = int64(rand_number(int(kigain), int(float64(kigain)*0.5)))
			}
			ch.Mana += kigain
			ch.Max_mana += kigain
			ch.Baseki += kigain
			ch.Moltexp -= kicost
			send_to_char(ch, libc.CString("Your body evolves to make better use of the way it is now, and you feel that your spirit has strengthened. @D[@CKi@D: @Y+%s@D]@n\r\n"), add_commas(kigain))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("stamina")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("st")) == 0 {
		if stcost > molt_threshold(ch) {
			send_to_char(ch, libc.CString("You need a few more evolution levels before you can start upgrading stamina.\r\n"))
			return
		} else if ch.Moltexp < stcost {
			send_to_char(ch, libc.CString("You do not have enough evolution experience.\r\n"))
			return
		} else {
			var stgain int64 = int64(float64(ch.Basest) * 0.01)
			if stgain <= 0 {
				stgain = int64(rand_number(1, 5))
			} else {
				stgain = int64(rand_number(int(stgain), int(float64(stgain)*0.5)))
			}
			ch.Move += stgain
			ch.Max_move += stgain
			ch.Basest += stgain
			ch.Moltexp -= stcost
			send_to_char(ch, libc.CString("Your body evolves to make better use of the way it is now, and you feel that your body has more stamina. @D[@GST@D: @Y+%s@D]@n\r\n"), add_commas(stgain))
		}
	}
}
func see_plant(obj *obj_data, ch *char_data) {
	var water int = (obj.Value[VAL_WATERLEVEL])
	if water >= 0 {
		switch obj.Value[VAL_MATURITY] {
		case 0:
			send_to_char(ch, libc.CString("@wA @G%s@y seed@w has been planted here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 1:
			send_to_char(ch, libc.CString("@wA very young @G%s@w has sprouted from a planter here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 2:
			send_to_char(ch, libc.CString("@wA half grown @G%s@w is in a planter here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 3:
			send_to_char(ch, libc.CString("@wA mature @G%s@w is growing in a planter here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 4:
			send_to_char(ch, libc.CString("@wA mature @G%s@w is flowering in a planter here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 5:
			send_to_char(ch, libc.CString("@wA mature @G%s@w that is close to harvestable is here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		case 6:
			send_to_char(ch, libc.CString("@wA @Rharvestable @G%s@w is in the planter here. @D(@C%d Water Hours@D)@n\r\n"), obj.Short_description, water)
		default:
		}
	} else {
		if water > -4 {
			send_to_char(ch, libc.CString("@yA @G%s@y that is looking a bit @rdry@y, is here.@n\r\n"), obj.Short_description)
		} else if water > -10 {
			send_to_char(ch, libc.CString("@yA @G%s@y that is looking extremely @rdry@y, is here.@n\r\n"), obj.Short_description)
		} else if water <= -10 {
			send_to_char(ch, libc.CString("@yA @G%s@y that is completely @rdead@y and @rwithered@y, is here.@n\r\n"), obj.Short_description)
		}
	}
}
func terrain_bonus(ch *char_data) float64 {
	var bonus float64 = 0.0
	switch func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}() {
	case SECT_FOREST:
		bonus += 0.5
	case SECT_SPACE:
		bonus += -0.5
	case SECT_WATER_NOSWIM:
		bonus += 0.25
	case SECT_MOUNTAIN:
		bonus += 0.1
	default:
		bonus += 0.0
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		bonus += -0.5
	}
	return bonus
}
func search_room(ch *char_data) {
	var (
		vict    *char_data
		next_v  *char_data
		perc    int     = int((float64(ch.Aff_abils.Intel) * 0.6) + float64(GET_SKILL(ch, SKILL_SPOT)) + float64(GET_SKILL(ch, SKILL_SEARCH)) + float64(GET_SKILL(ch, SKILL_LISTEN)))
		prob    int     = 0
		found   int     = 0
		bonus   float64 = 1.0
		terrain float64 = 1.0
	)
	if float64(ch.Move) < float64(ch.Max_move)*0.001 {
		send_to_char(ch, libc.CString("You do not have enough stamina.\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_SENSE) != 0 {
		bonus += float64(GET_SKILL(ch, SKILL_SENSE)) * 0.01
	}
	reveal_hiding(ch, 0)
	act(libc.CString("@y$n@Y begins searching the room carefully.@n"), TRUE, ch, nil, nil, TO_ROOM)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if AFF_FLAGGED(vict, AFF_HIDE) && vict != ch {
			if vict.Suppression >= 1 {
				perc *= int(float64(vict.Suppression) * 0.01)
			}
			prob = int(float64(vict.Aff_abils.Dex) + float64(vict.Aff_abils.Intel)*0.6 + float64(GET_SKILL(vict, SKILL_HIDE)) + float64(GET_SKILL(vict, SKILL_MOVE_SILENTLY)))
			if AFF_FLAGGED(vict, AFF_LIQUEFIED) {
				prob += int(float64(prob) * .5)
			}
			if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				perc += 5
			}
			if int(vict.Race) == RACE_MUTANT && ((vict.Genome[0]) == 5 || (vict.Genome[1]) == 5) {
				prob += 10
			}
			terrain += terrain_bonus(vict)
			if float64(perc)*bonus >= float64(prob)*terrain {
				act(libc.CString("@YYou find @y$N@Y hiding nearby!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@y$n@Y has found your hiding spot!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@y$n@Y has found @y$N's@Y hiding spot!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				reveal_hiding(vict, 4)
				found++
			}
		}
	}
	var obj *obj_data = nil
	for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj.Next_content {
		if OBJ_FLAGGED(obj, ITEM_BURIED) && float64(perc)*bonus > float64(rand_number(50, 200)) {
			act(libc.CString("@YYou uncover @y$p@Y, which had been burried here.@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("@y$n@Y uncovers @y$p@Y, which had burried here.@n"), TRUE, ch, obj, nil, TO_ROOM)
			obj.Extra_flags[int(ITEM_BURIED/32)] &= bitvector_t(int32(^(1 << (int(ITEM_BURIED % 32)))))
			found++
		}
	}
	ch.Move -= int64(float64(ch.Max_move) * 0.001)
	if found == 0 {
		send_to_char(ch, libc.CString("You find nothing hidden.\r\n"))
		return
	}
}

var weapon_disp [6]*byte = [6]*byte{libc.CString("Sword"), libc.CString("Dagger"), libc.CString("Spear"), libc.CString("Club"), libc.CString("Gun"), libc.CString("Brawling")}

func yesrace(num int) int {
	var okay int = TRUE
	switch num {
	case 2:
		fallthrough
	case 4:
		fallthrough
	case 8:
		fallthrough
	case 10:
		fallthrough
	case 11:
		fallthrough
	case 14:
		fallthrough
	case 20:
		okay = FALSE
	}
	return okay
}
func do_mimic(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if GET_SKILL(ch, SKILL_MIMIC) == 0 {
		send_to_char(ch, libc.CString("You do not know how to mimic the appearance of other races.\r\n"))
		return
	}
	var count int = 0
	var x int = 0
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@CMimic Menu\n@c--------------------@W\r\n"))
		for x = 0; x < NUM_RACES; x++ {
			if race_ok_gender[int(ch.Sex)][x] {
				if yesrace(x) != 0 {
					if count == 2 {
						send_to_char(ch, libc.CString("%s\n"), pc_race_types[x])
						count = 0
					} else {
						send_to_char(ch, libc.CString("%s\n"), pc_race_types[x])
						count += 1
					}
				}
			}
		}
		send_to_char(ch, libc.CString("Stop@n\r\n"))
		return
	}
	var prob int = GET_SKILL(ch, SKILL_MIMIC)
	var perc int = axion_dice(0)
	var mult float64 = float64(1 / prob)
	var cost int64 = int64(float64(ch.Max_mana) * mult)
	var israce int = FALSE
	var change int = -1
	x = 0
	for x = 0; x < NUM_RACES; x++ {
		if race_ok_gender[int(ch.Sex)][x] {
			if yesrace(x) != 0 {
				if libc.StrCaseCmp(&arg[0], pc_race_types[x]) == 0 {
					if ch.Mimic == x+1 {
						israce = TRUE
						x = int(NUM_RACES + 1)
					} else {
						change = x + 1
						x = int(NUM_RACES + 1)
					}
				}
			}
		}
	}
	if israce == TRUE {
		send_to_char(ch, libc.CString("You are already mimicing that race. To stop enter 'mimic stop'\r\n"))
		return
	} else if change > -1 && ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki to perform the technique.\r\n"))
		return
	} else if change > -1 && prob < perc {
		ch.Mana -= cost
		act(libc.CString("@mYou concentrate and attempt to create an illusion to obscure your racial features. However you frown as you realize you have failed.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n@m concentrates and the light around them seems to shift and blur. It stops a moment later and $e frowns.@n"), TRUE, ch, nil, nil, TO_ROOM)
		return
	} else if change > -1 {
		var buf [64936]byte
		ch.Mimic = change
		ch.Mana -= cost
		stdio.Sprintf(&buf[0], "@M$n@m concentrates for a moment and $s features start to blur as light bends around $m. Now $e appears to be %s @M%s!@n", AN(JUGGLERACE(ch)), JUGGLERACELOWER(ch))
		send_to_char(ch, libc.CString("@mYou concentrate for a moment and your features start to blur as you use your ki to bend the light around your body. You now appear to be %s %s.@n\r\n"), AN(JUGGLERACE(ch)), JUGGLERACELOWER(ch))
		act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("stop")) == 0 {
		act(libc.CString("@mYou concentrate for a moment and release the illusion that was mimicing another race.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n@m concentrates for a moment and SUDDENLY $s appearance changes some what!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mimic = 0
	} else {
		send_to_char(ch, libc.CString("That is not a race you can change into. Enter mimic without arugments for the mimic menu.\r\n"))
		return
	}
}
func do_kyodaika(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if int(ch.Race) != RACE_NAMEK {
		send_to_char(ch, libc.CString("You are not a namek!\r\n"))
		return
	}
	if int(ch.Real_abils.Str)+5 > 25 && (ch.Bonuses[BONUS_WIMP]) > 0 && (ch.Genome[0]) == 0 {
		send_to_char(ch, libc.CString("You can't handle having your strength increased beyond 25.\r\n"))
		return
	}
	if (ch.Genome[0]) == 0 {
		act(libc.CString("@GYou growl as your body grows to ten times its normal size!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n@G growls as $s body grows to ten times its normal size!@n"), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("@cStrength@D: @C+5\r\n@cSpeed@D: @c-2@n\r\n"))
		ch.Real_abils.Str += 5
		ch.Real_abils.Cha -= 2
		ch.Genome[0] = 11
		save_char(ch)
		return
	} else {
		act(libc.CString("@GYou growl as your body shrinks to its normal size!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n@G growls as $s body shrinks to its normal size!@n"), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("@cStrength@D: @C-5\r\n@cSpeed@D: @c+2@n\r\n"))
		ch.Real_abils.Str -= 5
		ch.Real_abils.Cha += 2
		ch.Genome[0] = 0
		save_char(ch)
		return
	}
}
func do_table(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj  *obj_data = nil
		obj2 *obj_data = nil
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: table (red | blue | green | yellow) (card name)"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't see that table here.\r\n"))
		return
	}
	if (func() *obj_data {
		obj2 = get_obj_in_list_vis(ch, &arg2[0], nil, obj.Contains)
		return obj2
	}()) == nil {
		send_to_char(ch, libc.CString("That card doesn't seem to be on that table.\r\n"))
		return
	}
	var buf [200]byte
	stdio.Sprintf(&buf[0], "$n looks at %s on %s.\r\n", obj2.Short_description, obj.Short_description)
	act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
	send_to_char(ch, libc.CString("%s"), obj2.Action_description)
}
func do_draw(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Sits == nil {
		send_to_char(ch, libc.CString("You are not sitting at a duel table.\r\n"))
		return
	}
	if GET_OBJ_VNUM(ch.Sits) < 604 || GET_OBJ_VNUM(ch.Sits) > 607 {
		send_to_char(ch, libc.CString("You need to be sitting at an official table to play.\r\n"))
		return
	}
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	var obj3 *obj_data = nil
	var next_obj *obj_data = nil
	var drawn int = FALSE
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, libc.CString("case"), nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't have a case.\r\n"))
		return
	}
	for obj2 = obj.Contains; obj2 != nil; obj2 = next_obj {
		next_obj = obj2.Next_content
		if drawn == FALSE {
			obj_from_obj(obj2)
			obj_to_char(obj2, ch)
			obj3 = obj2
			drawn = TRUE
		}
	}
	if drawn == FALSE {
		send_to_char(ch, libc.CString("You don't have any cards in the case!\r\n"))
		return
	} else {
		act(libc.CString("$n draws a card from $s $p.\r\n"), TRUE, ch, obj, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You draw a card.\r\n%s\r\n"), obj3.Action_description)
		return
	}
}
func do_shuffle(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Sits == nil {
		send_to_char(ch, libc.CString("You are not sitting at a duel table.\r\n"))
		return
	}
	if GET_OBJ_VNUM(ch.Sits) < 604 || GET_OBJ_VNUM(ch.Sits) > 607 {
		send_to_char(ch, libc.CString("You need to be sitting at an official table to play.\r\n"))
		return
	}
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	var next_obj *obj_data = nil
	var count int = 0
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, libc.CString("case"), nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't have a case.\r\n"))
		return
	}
	for obj2 = obj.Contains; obj2 != nil; obj2 = next_obj {
		next_obj = obj2.Next_content
		if !OBJ_FLAGGED(obj2, ITEM_ANTI_HIEROPHANT) {
			continue
		}
		count += 1
	}
	if count <= 0 {
		send_to_char(ch, libc.CString("You don't have any cards in the case!\r\n"))
		return
	}
	var total int = count
	for obj2 = obj.Contains; obj2 != nil; obj2 = next_obj {
		next_obj = obj2.Next_content
		obj_from_obj(obj2)
		obj_to_room(obj2, real_room(48))
	}
	for count > 0 {
		for obj2 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(48))))).Contents; obj2 != nil; obj2 = next_obj {
			next_obj = obj2.Next_content
			if !OBJ_FLAGGED(obj2, ITEM_ANTI_HIEROPHANT) {
				continue
			}
			if obj2 != nil && count > 1 && rand_number(1, 4) == 3 {
				count -= 1
				obj_from_room(obj2)
				obj_to_obj(obj2, obj)
			} else if obj2 != nil && count == 1 {
				count -= 1
				obj_from_room(obj2)
				obj_to_obj(obj2, obj)
			}
		}
	}
	send_to_char(ch, libc.CString("You shuffle the cards carefully.\r\n"))
	act(libc.CString("$n shuffles their deck."), TRUE, ch, nil, nil, TO_ROOM)
	send_to_room(ch.In_room, libc.CString("There were %d cards in the deck.\r\n"), total)
}
func do_hand(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj      *obj_data
		next_obj *obj_data
		arg      [2048]byte
		count    int = 0
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: hand (look | show)\r\n"))
		return
	}
	if libc.StrCaseCmp(libc.CString("look"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("@CYour hand contains:\r\n@D---------------------------@n\r\n"))
		for obj = ch.Carrying; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil && !OBJ_FLAGGED(obj, ITEM_ANTI_HIEROPHANT) {
				continue
			}
			if obj != nil {
				count += 1
				send_to_char(ch, libc.CString("%s\r\n"), obj.Short_description)
			}
		}
		act(libc.CString("$n looks at $s hand."), TRUE, ch, nil, nil, TO_ROOM)
		if count == 0 {
			send_to_char(ch, libc.CString("No cards."))
			act(libc.CString("There were no cards."), TRUE, ch, nil, nil, TO_ROOM)
		} else if count > 7 {
			act(libc.CString("You have more than seven cards in your hand."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n has more than seven cards in $s hand."), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			var buf [200]byte
			stdio.Sprintf(&buf[0], "There are %d cards in the hand.", count)
			act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if libc.StrCaseCmp(libc.CString("show"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("You show off your hand to the room.\r\n"))
		act(libc.CString("@C$n's hand contains:\r\n@D---------------------------@n"), TRUE, ch, nil, nil, TO_ROOM)
		for obj = ch.Carrying; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil && !OBJ_FLAGGED(obj, ITEM_ANTI_HIEROPHANT) {
				continue
			}
			if obj != nil {
				count += 1
				act(libc.CString("$p"), TRUE, ch, obj, nil, TO_ROOM)
			}
		}
		if count == 0 {
			act(libc.CString("No cards."), TRUE, ch, nil, nil, TO_ROOM)
		}
		if count > 7 {
			act(libc.CString("You have more than seven cards in your hand."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n has more than seven cards in $s hand."), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: hand (look | show)\r\n"))
		return
	}
}
func do_post(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
		obj  *obj_data
		obj2 *obj_data
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: post (obj name)\n        post (obj name) (target obj name)\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't seem to have that.\r\n"))
		return
	}
	if int(obj.Type_flag) != ITEM_NOTE {
		send_to_char(ch, libc.CString("You can only post notepaper.\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
		send_to_char(ch, libc.CString("You can not post on things in a garden.\r\n"))
		return
	}
	if arg2[0] == 0 {
		if (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_INSIDE && (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_CITY {
			send_to_char(ch, libc.CString("You are not near any general structure you can post it on.\r\n"))
			return
		}
		act(libc.CString("@WYou post $p@W on a nearby structure.@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W posts $p@W on a nearby structure.@n"), TRUE, ch, obj, nil, TO_ROOM)
		obj_from_char(obj)
		obj_to_room(obj, ch.In_room)
		obj.Posttype = 1
		return
	} else {
		if (func() *obj_data {
			obj2 = get_obj_in_list_vis(ch, &arg2[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj2
		}()) == nil {
			send_to_char(ch, libc.CString("You can't seem to find the thing you want to post it on.\r\n"))
			return
		} else if obj2.Posted_to != nil {
			send_to_char(ch, libc.CString("It already has something posted on it. Get that first if you want to post.\r\n"))
			return
		} else if int(obj2.Type_flag) == ITEM_BOARD {
			send_to_char(ch, libc.CString("Boards come with their own means of posting messages.\r\n"))
			return
		} else {
			var buf [64936]byte
			stdio.Sprintf(&buf[0], "@C$n@W posts %s@W on %s@W.@n", obj.Short_description, obj2.Short_description)
			send_to_char(ch, libc.CString("@WYou post %s@W on %s@W.@n\r\n"), obj.Short_description, obj2.Short_description)
			act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
			obj_from_char(obj)
			obj_to_room(obj, ch.In_room)
			obj.Posttype = 2
			obj.Posted_to = obj2
			obj2.Posted_to = obj
			return
		}
	}
}
func do_play(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Position) != POS_SITTING {
		send_to_char(ch, libc.CString("You need to be sitting at an official table to play.\r\n"))
		return
	}
	if ch.Sits == nil {
		send_to_char(ch, libc.CString("You need to be sitting at an official table to play.\r\n"))
		return
	}
	if GET_OBJ_VNUM(ch.Sits) < 604 || GET_OBJ_VNUM(ch.Sits) > 607 {
		send_to_char(ch, libc.CString("You need to be sitting at an official table to play.\r\n"))
		return
	}
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	var obj3 *obj_data = nil
	var next_obj *obj_data = nil
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: play (card name)"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't have that card to play.\r\n"))
		return
	}
	for obj3 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj3 != nil; obj3 = next_obj {
		next_obj = obj3.Next_content
		if GET_OBJ_VNUM(obj3) == GET_OBJ_VNUM(ch.Sits)-4 {
			obj2 = obj3
		}
	}
	if obj2 == nil {
		send_to_char(ch, libc.CString("Your table is missing. Inform an immortal of this problem.\r\n"))
		return
	}
	act(libc.CString("You play $p on your table."), TRUE, ch, obj, nil, TO_CHAR)
	act(libc.CString("$n plays $p on $s table."), TRUE, ch, obj, nil, TO_ROOM)
	obj_from_char(obj)
	obj_to_obj(obj, obj2)
}
func do_nickname(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj  *obj_data = nil
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: nickname (object) (nickname)\n"))
		send_to_char(ch, libc.CString("Syntax: nickname ship (nickname)\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("ship")) != 0 {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have that item to nickname.\r\n"))
			return
		}
	}
	if libc.StrLen(&arg2[0]) > 20 {
		send_to_char(ch, libc.CString("You can't nickname items with any name longer than 20 characters.\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("ship")) == 0 {
		var (
			ship     *obj_data = nil
			next_obj *obj_data = nil
			ship2    *obj_data = nil
			found    int       = FALSE
		)
		for ship = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; ship != nil; ship = next_obj {
			next_obj = ship.Next_content
			if GET_OBJ_VNUM(ship) >= 45000 && GET_OBJ_VNUM(ship) <= 0xB3AF && found == FALSE {
				found = TRUE
				ship2 = ship
			}
		}
		if found == TRUE {
			if libc.StrStr(&arg2[0], libc.CString("@")) != nil {
				send_to_char(ch, libc.CString("You can't nickname a ship and use color codes. Sorry.\r\n"))
				return
			} else {
				var nick [2048]byte
				stdio.Sprintf(&nick[0], "%s", CAP(&arg2[0]))
				ship2.Action_description = libc.StrDup(&nick[0])
				var k *obj_data
				for k = object_list; k != nil; k = k.Next {
					if GET_OBJ_VNUM(k) == GET_OBJ_VNUM(ship2)+1000 {
						extract_obj(k)
						var was_in int = int(func() room_vnum {
							if ship2.In_room != room_rnum(-1) && ship2.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ship2.In_room)))).Number
							}
							return -1
						}())
						obj_from_room(ship2)
						obj_to_room(ship2, real_room(room_vnum(was_in)))
					}
				}
			}
		}
		return
	}
	if libc.StrStr(obj.Short_description, libc.CString("nicknamed")) != nil {
		send_to_char(ch, libc.CString("%s@w has already been nicknamed.@n\r\n"), obj.Short_description)
		return
	} else if libc.StrStr(obj.Name, libc.CString("corpse")) != nil {
		send_to_char(ch, libc.CString("%s@w is a corpse!@n\r\n"), obj.Short_description)
		return
	} else {
		send_to_char(ch, libc.CString("@wYou nickname %s@w as '@C%s@w'.@n\r\n"), obj.Short_description, &arg2[0])
		var nick [2048]byte
		var nick2 [2048]byte
		stdio.Sprintf(&nick[0], "%s @wnicknamed @D(@C%s@D)@n", obj.Short_description, CAP(&arg2[0]))
		stdio.Sprintf(&nick2[0], "%s %s", obj.Name, &arg2[0])
		obj.Short_description = libc.StrDup(&nick[0])
		obj.Name = libc.StrDup(&nick2[0])
		return
	}
}

var cmd_sort_info *int
var portal_appearance [7]*byte = [7]*byte{libc.CString("All you can see is the glow of the portal."), libc.CString("You see an image of yourself in the room - my, you are looking attractive today."), libc.CString("All you can see is a swirling grey mist."), libc.CString("The scene is of the surrounding countryside, but somehow blurry and lacking focus."), libc.CString("The blackness appears to stretch on forever."), libc.CString("Suddenly, out of the blackness a flaming red eye appears and fixes its gaze upon you."), libc.CString("\n")}

func do_showoff(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj  *obj_data  = nil
		vict *char_data = nil
		arg  [2048]byte
		arg2 [2048]byte
	)
	arg[0] = '\x00'
	arg2[0] = '\x00'
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("You want to show what item to what character?\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't seem to have that.\r\n"))
		return
	} else if (func() *char_data {
		vict = get_player_vis(ch, &arg2[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("There is no such person around.\r\n"))
		return
	} else {
		act(libc.CString("@WYou hold up $p@W for @C$N@W to see:@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W holds up $p@W for you to see:@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W holds up $p@W for @c$N@W to see.@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		show_obj_to_char(obj, vict, SHOW_OBJ_ACTION)
		return
	}
}
func introCreate(ch *char_data) {
	var (
		fname [40]byte
		fl    *stdio.File
	)
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save user, %s, to filename, %s."), GET_NAME(ch), &fname[0])
		return
	}
	stdio.Fprintf(fl, "Gibbles Gibbles\n")
	fl.Close()
	return
}
func readIntro(ch *char_data, vict *char_data) int {
	var (
		fname  [40]byte
		filler [50]byte
		scrap  [100]byte
		line   [256]byte
		known  int = FALSE
		fl     *stdio.File
	)
	if vict == nil {
		return 0
	}
	if IS_NPC(ch) || IS_NPC(vict) {
		return 1
	}
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		introCreate(ch)
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
		stdio.Sscanf(&line[0], "%s %s\n", &filler[0], &scrap[0])
		if libc.StrCaseCmp(GET_NAME(vict), &filler[0]) == 0 {
			known = TRUE
		}
	}
	fl.Close()
	if known == TRUE {
		return 1
	} else {
		return 0
	}
}
func introWrite(ch *char_data, vict *char_data, name *byte) {
	var (
		file   *stdio.File
		fname  [40]byte
		filler [50]byte
		scrap  [100]byte
		line   [256]byte
		names  [500]*byte = [500]*byte{0: libc.CString("")}
		alias  [500]*byte = [500]*byte{0: libc.CString("")}
		fl     *stdio.File
		count  int = 0
		x      int = 0
	)
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		introCreate(ch)
	}
	if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		return
	}
	for int(file.IsEOF()) == 0 || count < 498 {
		get_line(file, &line[0])
		stdio.Sscanf(&line[0], "%s %s\n", &filler[0], &scrap[0])
		names[count] = libc.StrDup(&filler[0])
		alias[count] = libc.StrDup(&scrap[0])
		count++
		filler[0] = '\x00'
		scrap[0] = '\x00'
	}
	file.Close()
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save intro file, %s, to filename, %s."), GET_NAME(ch), &fname[0])
		return
	}
	for x < count {
		if x == 0 || libc.StrCaseCmp(names[x-1], names[x]) != 0 {
			if libc.StrCaseCmp(names[x], GET_NAME(vict)) != 0 {
				stdio.Fprintf(fl, "%s %s\n", names[x], alias[x])
			}
		}
		x++
	}
	x = 0
	for x < count {
		if names[x] != nil {
			libc.Free(unsafe.Pointer(names[x]))
		}
		if alias[x] != nil {
			libc.Free(unsafe.Pointer(alias[x]))
		}
		x++
	}
	stdio.Fprintf(fl, "%s %s\n", GET_NAME(vict), CAP(name))
	fl.Close()
	return
}
func do_intro(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	var vict *char_data
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: dub (target) (name)\r\nWho do you want to dub and what do you want to name them?\r\n"))
		return
	}
	if arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: dub (target) (name)\r\nWhat name do you wish to know them by?\r\n"))
		return
	}
	if libc.StrLen(&arg2[0]) > 20 {
		send_to_char(ch, libc.CString("Limit the name to 20 characters.\r\n"))
		return
	}
	if libc.StrLen(&arg2[0]) < 3 {
		send_to_char(ch, libc.CString("Limit the name to at least 3 characters.\r\n"))
		return
	}
	if libc.StrStr(&arg2[0], libc.CString("$")) != nil || libc.StrStr(&arg2[0], libc.CString("@")) != nil || libc.StrStr(&arg2[0], libc.CString("%")) != nil {
		send_to_char(ch, libc.CString("Illegal character. No symbols.\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_player_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("There is no such person.\r\n"))
		return
	}
	if vict == ch {
		send_to_char(ch, libc.CString("That seems rather odd.\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("That seems rather unwise.\r\n"))
		return
	}
	if readIntro(vict, ch) == 2 {
		send_to_char(ch, libc.CString("There seems to have been an error, report this to Iovan.\r\n"))
		return
	} else if readIntro(ch, vict) == 1 && libc.StrStr(JUGGLERACE(vict), &arg[0]) != nil {
		send_to_char(ch, libc.CString("You have already dubbed them a name. If you want to redub them target the name you know them by.\r\n"))
		return
	} else {
		introWrite(ch, vict, &arg2[0])
		act(libc.CString("You decide to call $M, $N."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n seems to decide something about you."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n seems to decide something about $N."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		return
	}
}
func bringdesc(ch *char_data, tch *char_data) {
	if ch != nil && tch != nil && IS_HUMANOID(tch) {
		if ch != tch && PLR_FLAGGED(tch, PLR_DISGUISED) {
			send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WHidden.         @D]@n\r\n"))
			send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WHidden.         @D]@n\r\n"))
			send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WHidden.         @D]@n\r\n"))
			send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WHidden.         @D]@n\r\n"))
			if int(tch.Skin) == SKIN_WHITE {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WWhite.        @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_TAN {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WTan.          @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_BLACK {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WBlack.        @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_GREEN {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WGreen.        @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_ORANGE {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WOrange.       @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_YELLOW {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WYellow.       @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_RED {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WRed.          @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_GREY {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WGrey.         @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_BLUE {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WBlue.         @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_AQUA {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WAqua.         @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_PINK {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WPink.         @D]@n\r\n"))
			} else if int(tch.Skin) == SKIN_PURPLE {
				send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WPurple.       @D]@n\r\n"))
			}
			return
		}
		if int(tch.Race) == RACE_HUMAN || int(tch.Race) == RACE_SAIYAN || int(tch.Race) == RACE_KONATSU || int(tch.Race) == RACE_MUTANT || int(tch.Race) == RACE_ANDROID || int(tch.Race) == RACE_KAI || int(tch.Race) == RACE_HALFBREED || int(tch.Race) == RACE_TRUFFLE || int(tch.Race) == RACE_HOSHIJIN {
			if int(tch.Race) != RACE_SAIYAN && int(tch.Race) != RACE_HALFBREED || (int(tch.Race) == RACE_SAIYAN || int(tch.Race) == RACE_HALFBREED) && !IS_TRANSFORMED(tch) {
				if int(tch.Hairl) == HAIRL_LONG {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WLong.         @D]@n\r\n"))
				} else if int(tch.Hairl) == HAIRL_BALD {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WBald.         @D]@n\r\n"))
				} else if int(tch.Hairl) == HAIRL_SHORT {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WShort.        @D]@n\r\n"))
				} else if int(tch.Hairl) == HAIRL_MEDIUM {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WMedium.       @D]@n\r\n"))
				} else if int(tch.Hairl) == HAIRL_RLONG {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WReally Long.  @D]@n\r\n"))
				}
				if int(tch.Hairs) == HAIRS_PLAIN {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WPlain.        @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_MOHAWK {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WMohawk.       @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_SPIKY {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WSpiky.        @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_CURLY {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WCurly.        @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_UNEVEN {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WUneven.       @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_PONYTAIL {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WPony Tail.    @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_AFRO {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WAfro.         @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_FADE {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WFade.         @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_CREW {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WCrew Cut.     @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_FEATHERED {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WFeathered.    @D]@n\r\n"))
				} else if int(tch.Hairs) == HAIRS_DRED {
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WDread Locks.  @D]@n\r\n"))
				}
				if int(tch.Hairc) == HAIRC_BLACK {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WBlack.        @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_BROWN {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WBrown.        @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_BLONDE {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WBlonde.       @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_GREY {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WGrey.         @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_RED {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WRed.          @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_ORANGE {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WOrange.       @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_GREEN {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WGreen.        @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_BLUE {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WBlue.         @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_PINK {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WPink.         @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_PURPLE {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WPurple.       @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_SILVER {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WSilver.       @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_CRIMSON {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WCrimson.      @D]@n\r\n"))
				} else if int(tch.Hairc) == HAIRC_WHITE {
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WWhite.        @D]@n\r\n"))
				}
			} else if int(tch.Race) == RACE_SAIYAN || int(tch.Race) == RACE_HALFBREED {
				if PLR_FLAGGED(tch, PLR_TRANS1) {
					if int(tch.Hairl) == HAIRL_LONG {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WLong.         @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_BALD {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WBald.         @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_SHORT {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WShort.        @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_MEDIUM {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WMedium.       @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_RLONG {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WReally Long.  @D]@n\r\n"))
					}
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WSpiky.        @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WGolden.       @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WEmerald.      @D]@n\r\n"))
				} else if PLR_FLAGGED(tch, PLR_TRANS2) {
					if int(tch.Hairl) == HAIRL_LONG {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WLong.         @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_BALD {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WBald.         @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_SHORT {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WShort.        @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_MEDIUM {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WMedium.       @D]@n\r\n"))
					} else if int(tch.Hairl) == HAIRL_RLONG {
						send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WReally Long.  @D]@n\r\n"))
					}
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WSharp Spikes. @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WGolden.       @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WEmerald.      @D]@n\r\n"))
				} else if PLR_FLAGGED(tch, PLR_TRANS3) {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WReally Long.  @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WSpiky.        @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WGolden.       @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WAqua Green.   @D]@n\r\n"))
				} else if PLR_FLAGGED(tch, PLR_TRANS4) {
					send_to_char(ch, libc.CString("            @D[@cHair Length @D: @WLong.        @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Style  @D: @WSoft Spikes. @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cHair Color  @D: @WBlack.       @D]@n\r\n"))
					send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WAmber.       @D]@n\r\n"))
				}
			}
		}
		if int(tch.Race) == RACE_DEMON || int(tch.Race) == RACE_ICER {
			if int(tch.Hairl) == HAIRL_BALD {
				send_to_char(ch, libc.CString("            @D[@cHorn Length @D: @WNone.         @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_SHORT {
				send_to_char(ch, libc.CString("            @D[@cHorn Length @D: @WShort.        @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_MEDIUM {
				send_to_char(ch, libc.CString("            @D[@cHorn Length @D: @WMedium.       @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_LONG {
				send_to_char(ch, libc.CString("            @D[@cHorn Length @D: @WLong.         @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_RLONG {
				send_to_char(ch, libc.CString("            @D[@cHorn Length @D: @WReally Long.  @D]@n\r\n"))
			}
		}
		if int(tch.Race) == RACE_NAMEK || int(tch.Race) == RACE_ARLIAN {
			if int(tch.Hairl) == HAIRL_BALD {
				send_to_char(ch, libc.CString("            @D[@cAnt. Length @D: @WTiny.        @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_SHORT {
				send_to_char(ch, libc.CString("            @D[@cAnt. Length @D: @WShort.       @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_MEDIUM {
				send_to_char(ch, libc.CString("            @D[@cAnt. Length @D: @WMedium.      @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_LONG {
				send_to_char(ch, libc.CString("            @D[@cAnt. Length @D: @WLong.        @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_RLONG {
				send_to_char(ch, libc.CString("            @D[@cAnt. Length @D: @WR. Long.     @D]@n\r\n"))
			}
		}
		if int(tch.Race) == RACE_ARLIAN && int(tch.Sex) == SEX_FEMALE {
			if int(tch.Hairc) == HAIRC_BLACK {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WBlack.        @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_BROWN {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WBrown.        @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_BLONDE {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WBlonde.       @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_GREY {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WGrey.         @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_RED {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WRed.          @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_ORANGE {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WOrange.       @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_GREEN {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WGreen.        @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_BLUE {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WBlue.         @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_PINK {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WPink.         @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_PURPLE {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WPurple.       @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_SILVER {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WSilver.       @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_CRIMSON {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WCrimson.      @D]@n\r\n"))
			} else if int(tch.Hairc) == HAIRC_WHITE {
				send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WWhite.        @D]@n\r\n"))
			}
		} else if int(tch.Race) == RACE_ARLIAN && int(tch.Sex) != SEX_FEMALE {
			send_to_char(ch, libc.CString("            @D[@cWing Color  @D: @WWhite.        @D]@n\r\n"))
		}
		if int(tch.Race) == RACE_MAJIN {
			if int(tch.Hairl) == HAIRL_BALD {
				send_to_char(ch, libc.CString("            @D[@cFor. Length @D: @WTiny.         @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_SHORT {
				send_to_char(ch, libc.CString("            @D[@cFor. Length @D: @WShort.        @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_MEDIUM {
				send_to_char(ch, libc.CString("            @D[@cFor. Length @D: @WMedium.       @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_LONG {
				send_to_char(ch, libc.CString("            @D[@cFor. Length @D: @WLong.         @D]@n\r\n"))
			}
			if int(tch.Hairl) == HAIRL_RLONG {
				send_to_char(ch, libc.CString("            @D[@cFor. Length @D: @WR. Long.      @D]@n\r\n"))
			}
		}
		if int(tch.Race) != RACE_SAIYAN && int(tch.Race) != RACE_HALFBREED || (int(tch.Race) == RACE_SAIYAN || int(tch.Race) == RACE_HALFBREED) && !IS_TRANSFORMED(tch) {
			if int(tch.Eye) == EYE_BLUE {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WBlue.         @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_BLACK {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WBlack.        @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_GREEN {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WGreen.        @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_BROWN {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WBrown.        @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_RED {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WRed.          @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_AQUA {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WAqua.         @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_PINK {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WPink.         @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_PURPLE {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WPurple.       @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_CRIMSON {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WCrimson.      @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_GOLD {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WGold.         @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_AMBER {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WAmber.        @D]@n\r\n"))
			} else if int(tch.Eye) == EYE_EMERALD {
				send_to_char(ch, libc.CString("            @D[@cEye Color   @D: @WEmerald.      @D]@n\r\n"))
			}
		}
		if int(tch.Skin) == SKIN_WHITE {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WWhite.        @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_TAN {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WTan.          @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_BLACK {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WBlack.        @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_GREEN {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WGreen.        @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_ORANGE {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WOrange.       @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_YELLOW {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WYellow.       @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_RED {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WRed.          @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_GREY {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WGrey.         @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_BLUE {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WBlue.         @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_AQUA {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WAqua.         @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_PINK {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WPink.         @D]@n\r\n"))
		} else if int(tch.Skin) == SKIN_PURPLE {
			send_to_char(ch, libc.CString("            @D[@cSkin Color  @D: @WPurple.       @D]@n\r\n"))
		}
		if tch.Majinize != 0 && tch.Majinize != 3 {
			send_to_char(ch, libc.CString("            @D[@cForehead    @D: @mMajin Symbol  @D]@n\r\n"))
		}
	} else if !IS_HUMANOID(tch) {
		return
	} else {
		send_to_char(ch, libc.CString("Error in bring-desc, please report.\r\n"))
	}
}
func map_draw_room(map_ [9][10]byte, x int, y int, rnum room_rnum, ch *char_data) {
	var door int
	for door = 0; door < NUM_OF_DIRS; door++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door] != nil && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door], 1<<1) && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door], 1<<4)) {
			switch door {
			case NORTH:
				map_[y-1][x] = '8'
			case EAST:
				map_[y][x+1] = '8'
			case SOUTH:
				map_[y+1][x] = '8'
			case WEST:
				map_[y][x-1] = '8'
			case NORTHEAST:
				map_[y-1][x+1] = '8'
			case NORTHWEST:
				map_[y-1][x-1] = '8'
			case SOUTHEAST:
				map_[y+1][x+1] = '8'
			case SOUTHWEST:
				map_[y+1][x-1] = '8'
			}
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door] != nil && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door], 1<<1) {
			switch door {
			case NORTH:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y-1][x] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '2'
					} else {
						map_[y-1][x] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '2'
					} else {
						map_[y-1][x] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '7'
					} else {
						map_[y-1][x] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '1'
					} else {
						map_[y-1][x] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '6'
					} else {
						map_[y-1][x] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '5'
					} else {
						map_[y-1][x] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x] = '3'
					} else {
						map_[y-1][x] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y-1][x] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y-1][x] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y-1][x] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y-1][x] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y-1][x] = '*'
				} else {
					map_[y-1][x] = '-'
				}
			case EAST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y][x+1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '2'
					} else {
						map_[y][x+1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '2'
					} else {
						map_[y][x+1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '7'
					} else {
						map_[y][x+1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '1'
					} else {
						map_[y][x+1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '6'
					} else {
						map_[y][x+1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '5'
					} else {
						map_[y][x+1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x+1] = '3'
					} else {
						map_[y][x+1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y][x+1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y][x+1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y][x+1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y][x+1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y][x+1] = '*'
				} else {
					map_[y][x+1] = '-'
				}
			case SOUTH:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y+1][x] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '2'
					} else {
						map_[y+1][x] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '2'
					} else {
						map_[y+1][x] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '7'
					} else {
						map_[y+1][x] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '1'
					} else {
						map_[y+1][x] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '6'
					} else {
						map_[y+1][x] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '5'
					} else {
						map_[y+1][x] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '3'
					} else {
						map_[y+1][x] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y+1][x] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y+1][x] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y+1][x] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y+1][x] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y+1][x] = '*'
				} else {
					map_[y+1][x] = '-'
				}
			case WEST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y][x-1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x] = '2'
					} else {
						map_[y][x-1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '2'
					} else {
						map_[y][x-1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '7'
					} else {
						map_[y][x-1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '1'
					} else {
						map_[y][x-1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '6'
					} else {
						map_[y][x-1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '5'
					} else {
						map_[y][x-1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y][x-1] = '3'
					} else {
						map_[y][x-1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y][x-1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y][x-1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y][x-1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y][x-1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y][x-1] = '*'
				} else {
					map_[y][x-1] = '-'
				}
			case NORTHEAST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y-1][x+1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '2'
					} else {
						map_[y-1][x+1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '2'
					} else {
						map_[y-1][x+1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '7'
					} else {
						map_[y-1][x+1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '1'
					} else {
						map_[y-1][x+1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '6'
					} else {
						map_[y-1][x+1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '5'
					} else {
						map_[y-1][x+1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x+1] = '3'
					} else {
						map_[y-1][x+1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y-1][x+1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y-1][x+1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y-1][x+1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y-1][x+1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y-1][x+1] = '*'
				} else {
					map_[y-1][x+1] = '-'
				}
			case NORTHWEST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y-1][x-1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '2'
					} else {
						map_[y-1][x-1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '2'
					} else {
						map_[y-1][x-1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '7'
					} else {
						map_[y-1][x-1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '1'
					} else {
						map_[y-1][x-1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '6'
					} else {
						map_[y-1][x-1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '5'
					} else {
						map_[y-1][x-1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y-1][x-1] = '3'
					} else {
						map_[y-1][x-1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y-1][x-1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y-1][x-1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y-1][x-1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y-1][x-1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y-1][x-1] = '*'
				} else {
					map_[y-1][x-1] = '-'
				}
			case SOUTHEAST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y+1][x+1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '2'
					} else {
						map_[y+1][x+1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '2'
					} else {
						map_[y+1][x+1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '7'
					} else {
						map_[y+1][x+1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '1'
					} else {
						map_[y+1][x+1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '6'
					} else {
						map_[y+1][x+1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '5'
					} else {
						map_[y+1][x+1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x+1] = '3'
					} else {
						map_[y+1][x+1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y+1][x+1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y+1][x+1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y+1][x+1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y+1][x+1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y+1][x+1] = '*'
				} else {
					map_[y+1][x+1] = '-'
				}
			case SOUTHWEST:
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect < 0 || (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_UNDERWATER {
					map_[y+1][x-1] = '='
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_INSIDE {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '2'
					} else {
						map_[y+1][x-1] = 'i'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FIELD {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '2'
					} else {
						map_[y+1][x-1] = 'p'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_DESERT {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '7'
					} else {
						map_[y+1][x-1] = '!'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_CITY {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '1'
					} else {
						map_[y+1][x-1] = '('
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FOREST {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '6'
					} else {
						map_[y+1][x-1] = 'f'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_MOUNTAIN {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '5'
					} else {
						map_[y+1][x-1] = '^'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_HILLS {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Geffect >= 1 {
						map_[y+1][x-1] = '3'
					} else {
						map_[y+1][x-1] = 'h'
					}
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_FLYING {
					map_[y+1][x-1] = 's'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_NOSWIM {
					map_[y+1][x-1] = '`'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_WATER_SWIM {
					map_[y+1][x-1] = '+'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_SHOP {
					map_[y+1][x-1] = '&'
				} else if (func() int {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[door].To_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) == SECT_IMPORTANT {
					map_[y+1][x-1] = '*'
				} else {
					map_[y+1][x-1] = '-'
				}
			}
		}
	}
}
func do_map(ch *char_data, argument *byte, cmd int, subcmd int) {
	gen_map(ch, 1)
}
func gen_map(ch *char_data, num int) {
	var (
		door int
		i    int
		map_ [9][10]byte = [9][10]byte{0: {0: '-'}, 1: {0: '-'}}
		buf2 [2048]byte
	)
	if num == 1 {
		send_to_char(ch, libc.CString("@W               @D-[@CArea Map@D]-\r\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@w\r\n"))
		send_to_char(ch, libc.CString("@WC = City, @wI@W = Inside, @GP@W = Plain, @gF@W = Forest\r\n"))
		send_to_char(ch, libc.CString("@DM@W = Mountain, @yH@W = Hills, @CS@W = Sky, @BW@W = Water\r\n"))
		send_to_char(ch, libc.CString("@bU@W = Underwater, @m$@W = Shop, @m#@W = Important,\r\n"))
		send_to_char(ch, libc.CString("@YD@W = Desert, @c~@W = Shallow Water, @4 @n@W = Lava,\r\n"))
		send_to_char(ch, libc.CString("@WLastly @RX@W = You.\r\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------\r\n"))
		send_to_char(ch, libc.CString("@D                  @CNorth@w\r\n"))
		send_to_char(ch, libc.CString("@D                    @c^@w\r\n"))
		send_to_char(ch, libc.CString("@D             @CWest @c< O > @CEast@w\r\n"))
		send_to_char(ch, libc.CString("@D                    @cv@w\r\n"))
		send_to_char(ch, libc.CString("@D                  @CSouth@w\r\n"))
		send_to_char(ch, libc.CString("@D                ---------@w\r\n"))
	}
	for i = 0; i < 9; i++ {
		libc.StrCpy(&map_[i][0], libc.CString("         "))
	}
	map_draw_room(map_, 4, 4, ch.In_room, ch)
	for door = 0; door < NUM_OF_DIRS; door++ {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room != room_rnum(-1) && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1) {
			switch door {
			case NORTH:
				map_draw_room(map_, 4, 3, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case EAST:
				map_draw_room(map_, 5, 4, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case SOUTH:
				map_draw_room(map_, 4, 5, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case WEST:
				map_draw_room(map_, 3, 4, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case NORTHEAST:
				map_draw_room(map_, 5, 3, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case NORTHWEST:
				map_draw_room(map_, 3, 3, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case SOUTHEAST:
				map_draw_room(map_, 5, 5, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			case SOUTHWEST:
				map_draw_room(map_, 3, 5, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ch)
			}
		}
	}
	map_[4][4] = 'x'
	var key int = 0
	buf2[0] = '\x00'
	for i = 2; i < 9; i++ {
		if i > 6 {
			continue
		}
		if num == 1 {
			stdio.Sprintf(&buf2[0], "@w                %s\r\n", &map_[i][0])
		} else {
			if i == 2 {
				stdio.Sprintf(&buf2[0], "@w       @w|%s@w|           %s", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[0]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[0], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[0], 1<<1) {
							return " @rN "
						}
						return " @CN "
					}
					return "   "
				}(), &map_[i][0])
			}
			if i == 3 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|     %s", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[6]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[6], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[6], 1<<1) {
							return " @rNW"
						}
						return " @CNW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[4]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[4], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[4], 1<<1) {
							return " @yU "
						}
						return " @YU "
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[7]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[7], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[7], 1<<4) {
							return "@rNE "
						}
						return "@CNE "
					}
					return "   "
				}(), &map_[i][0])
			}
			if i == 4 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|     %s", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[3]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[3], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[3], 1<<1) {
							return "  @rW"
						}
						return "  @CW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[10]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[10], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[10], 1<<1) {
							return " @rI "
						}
						return " @mI "
					}
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[11]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[11], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[11], 1<<1) {
							return "@rOUT"
						}
						return "@mOUT"
					}
					return "@r{ }"
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[1]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[1], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[1], 1<<1) {
							return "@rE  "
						}
						return "@CE  "
					}
					return "   "
				}(), &map_[i][0])
			}
			if i == 5 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|     %s", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[9]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[9], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[9], 1<<1) {
							return " @rSW"
						}
						return " @CSW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5], 1<<1) {
							return " @yD "
						}
						return " @YD "
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[8]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[8], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[8], 1<<4) {
							return "@rSE "
						}
						return "@CSE "
					}
					return "   "
				}(), &map_[i][0])
			}
			if i == 6 {
				stdio.Sprintf(&buf2[0], "@w       @w|%s@w|           %s", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[2]) != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[2], 1<<4) {
						if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[2], 1<<1) {
							return " @rS "
						}
						return " @CS "
					}
					return "   "
				}(), &map_[i][0])
			}
		}
		search_replace(&buf2[0], libc.CString("x"), libc.CString("@RX"))
		search_replace(&buf2[0], libc.CString("&"), libc.CString("@m$"))
		search_replace(&buf2[0], libc.CString("*"), libc.CString("@m#"))
		search_replace(&buf2[0], libc.CString("+"), libc.CString("@c~"))
		search_replace(&buf2[0], libc.CString("s"), libc.CString("@CS"))
		search_replace(&buf2[0], libc.CString("i"), libc.CString("@wI"))
		search_replace(&buf2[0], libc.CString("("), libc.CString("@WC"))
		search_replace(&buf2[0], libc.CString("^"), libc.CString("@DM"))
		search_replace(&buf2[0], libc.CString("h"), libc.CString("@yH"))
		search_replace(&buf2[0], libc.CString("`"), libc.CString("@BW"))
		search_replace(&buf2[0], libc.CString("="), libc.CString("@bU"))
		search_replace(&buf2[0], libc.CString("p"), libc.CString("@GP"))
		search_replace(&buf2[0], libc.CString("f"), libc.CString("@gF"))
		search_replace(&buf2[0], libc.CString("!"), libc.CString("@YD"))
		search_replace(&buf2[0], libc.CString("-"), libc.CString("@w:"))
		search_replace(&buf2[0], libc.CString("1"), libc.CString("@4@YC@n"))
		search_replace(&buf2[0], libc.CString("2"), libc.CString("@4@YP@n"))
		search_replace(&buf2[0], libc.CString("3"), libc.CString("@4@YH@n"))
		search_replace(&buf2[0], libc.CString("7"), libc.CString("@4@YD@n"))
		search_replace(&buf2[0], libc.CString("5"), libc.CString("@4@YM@n"))
		search_replace(&buf2[0], libc.CString("6"), libc.CString("@4@YF@n"))
		search_replace(&buf2[0], libc.CString("8"), libc.CString("@1 @n"))
		if num != 1 {
			if key == 0 {
				send_to_char(ch, libc.CString("%s    @WC: City, @wI@W: Inside, @GP@W: Plain@n\r\n"), &buf2[0])
			}
			if key == 1 {
				send_to_char(ch, libc.CString("%s    @gF@W: Forest, @DM@W: Mountain, @yH@W: Hills@n\r\n"), &buf2[0])
			}
			if key == 2 {
				send_to_char(ch, libc.CString("%s    @CS@W: Sky, @BW@W: Water, @bU@W: Underwater@n\r\n"), &buf2[0])
			}
			if key == 3 {
				send_to_char(ch, libc.CString("%s    @m$@W: Shop, @m#@W: Important, @YD@W: Desert@n\r\n"), &buf2[0])
			}
			if key == 4 {
				send_to_char(ch, libc.CString("%s    @c~@W: Shallow Water, @4 @n@W: Lava, @RX@W: You@n\r\n"), &buf2[0])
			}
			key += 1
		} else {
			send_to_char(ch, &buf2[0])
		}
	}
	if num == 1 {
		send_to_char(ch, libc.CString("@D                ---------@w\r\n"))
	}
}
func display_spells(ch *char_data, obj *obj_data) {
	var i int
	send_to_char(ch, libc.CString("The spellbook contains the following spells:\r\n"))
	send_to_char(ch, libc.CString("@c---@wSpell Name@c------------------------------------@w# of pages@c-----@n\r\n"))
	if obj.Sbinfo == nil {
		return
	}
	for i = 0; i < SPELLBOOK_SIZE; i++ {
		if (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(i)))).Spellname != 0 {
			if (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(i)))).Spellname > SPELL_SENSU {
				continue
			}
			send_to_char(ch, libc.CString("@y%-20s@n\t\t\t\t\t[@R%2d@n]\r\n"), spell_info[(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(i)))).Spellname].Name, (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(i)))).Pages)
		}
	}
	return
}
func display_scroll(ch *char_data, obj *obj_data) {
	send_to_char(ch, libc.CString("The scroll contains the following spell:\r\n"))
	send_to_char(ch, libc.CString("@c---@wSpell Name@c---------------------------------------------------@n\r\n"))
	send_to_char(ch, libc.CString("@y%-20s@n\r\n"), skill_name(obj.Value[VAL_SCROLL_SPELL1]))
	return
}
func show_obj_to_char(obj *obj_data, ch *char_data, mode int) {
	if obj == nil || ch == nil {
		basic_mud_log(libc.CString("SYSERR: NULL pointer in show_obj_to_char(): obj=%p ch=%p"), obj, ch)
		return
	}
	var spotted int = FALSE
	if GET_SKILL(ch, SKILL_SPOT) > rand_number(20, 110) {
		spotted = TRUE
	}
	switch mode {
	case SHOW_OBJ_LONG:
		if *obj.Description == '.' && (IS_NPC(ch) || !PRF_FLAGGED(ch, PRF_HOLYLIGHT)) {
			return
		}
		if int(obj.Type_flag) == ITEM_VEHICLE && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) == room_vnum(obj.Value[0]) {
			return
		}
		if obj.Sitting != nil && ch.Admlevel < 1 {
			return
		}
		if obj.Sitting != nil && ch.Admlevel >= 1 {
			send_to_char(ch, libc.CString("@D(@YBeing Used@D)@w"))
		}
		if int(obj.Type_flag) == ITEM_PLANT && (ROOM_FLAGGED(obj.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(obj.In_room, ROOM_GARDEN2)) {
			see_plant(obj, ch)
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BURIED) {
			var bury [2048]byte
			if !IS_CORPSE(obj) {
				if obj.Weight < 10 {
					stdio.Sprintf(&bury[0], "small mound of")
				} else if obj.Weight < 50 {
					stdio.Sprintf(&bury[0], "medium sized mound of")
				} else if obj.Weight < 1000 {
					stdio.Sprintf(&bury[0], "large mound of")
				} else {
					stdio.Sprintf(&bury[0], "gigantic mound of")
				}
			} else {
				stdio.Sprintf(&bury[0], "recent grave covered by")
			}
			if spotted == TRUE && (func() int {
				if obj.In_room != room_rnum(-1) && obj.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) != SECT_DESERT {
				send_to_char(ch, libc.CString("@yA %s soft dirt is here.@n\r\n"), &bury[0])
				return
			} else if spotted == TRUE && (func() int {
				if obj.In_room != room_rnum(-1) && obj.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) == SECT_DESERT {
				send_to_char(ch, libc.CString("@YA %s soft sand is here.@n\r\n"), &bury[0])
				return
			} else {
				return
			}
		}
		if GET_OBJ_VNUM(obj) == 11 {
			send_to_char(ch, libc.CString("@wA gravity generator, set to %sx gravity, is built here"), add_commas(obj.Weight))
		} else if GET_OBJ_VNUM(obj) == 79 {
			send_to_char(ch, libc.CString("@wA @cG@Cl@wa@cc@Ci@wa@cl @wW@ca@Cl@wl @D[@C%s@D]@w is blocking access to the @G%s@w direction"), add_commas(obj.Weight), dirs[obj.Cost])
		} else {
			send_to_char(ch, libc.CString("@w"))
			if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_ROOMFLAGS) {
				if obj.Posted_to == nil {
					send_to_char(ch, libc.CString("@D[@G%d@D]@w "), GET_OBJ_VNUM(obj))
					if obj.Script != nil {
						send_to_char(ch, libc.CString("@D[@wT%d@D]@w "), obj.Proto_script.Vnum)
					}
				} else {
					if obj.Posttype <= 0 {
						send_to_char(ch, libc.CString("@D[@G%d@D]@w "), GET_OBJ_VNUM(obj))
						if obj.Script != nil {
							send_to_char(ch, libc.CString("@D[@wT%d@D]@w "), obj.Proto_script.Vnum)
						}
					}
				}
			}
			if obj.Posttype > 0 {
				if obj.Posted_to != nil {
					return
				} else {
					send_to_char(ch, libc.CString("%s@w, has been posted here.@n"), obj.Short_description)
				}
			} else {
				if !OBJ_FLAGGED(obj, ITEM_BURIED) {
					send_to_char(ch, libc.CString("%s@n"), obj.Description)
				}
			}
			if int(obj.Type_flag) == ITEM_VEHICLE {
				if !OBJVAL_FLAGGED(obj, 1<<2) && GET_OBJ_VNUM(obj) > 0x4AFF {
					send_to_char(ch, libc.CString("\r\n@c...its outer hatch is open@n"))
				} else if !OBJVAL_FLAGGED(obj, 1<<2) && GET_OBJ_VNUM(obj) <= 0x4AFF {
					send_to_char(ch, libc.CString("\r\n@c...its door is open@n"))
				}
			}
			if int(obj.Type_flag) == ITEM_CONTAINER && !IS_CORPSE(obj) {
				if !OBJVAL_FLAGGED(obj, 1<<2) && !OBJ_FLAGGED(obj, ITEM_SHEATH) {
					send_to_char(ch, libc.CString(". @D[@G-open-@D]@n"))
				} else if !OBJ_FLAGGED(obj, ITEM_SHEATH) {
					send_to_char(ch, libc.CString(". @D[@rclosed@D]@n"))
				}
			}
			if int(obj.Type_flag) == ITEM_HATCH {
				if !OBJVAL_FLAGGED(obj, 1<<2) {
					send_to_char(ch, libc.CString(", it is open"))
				} else if OBJVAL_FLAGGED(obj, 1<<2) {
					send_to_char(ch, libc.CString(", it is closed"))
				}
				if OBJVAL_FLAGGED(obj, 1<<3) {
					send_to_char(ch, libc.CString(" and locked@n"))
				} else {
					send_to_char(ch, libc.CString("@n"))
				}
			}
			if int(obj.Type_flag) == ITEM_FOOD {
				if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob {
					send_to_char(ch, libc.CString(", and it has been ate on@n"))
				}
			}
		}
	case SHOW_OBJ_SHORT:
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_ROOMFLAGS) {
			send_to_char(ch, libc.CString("[%d] "), GET_OBJ_VNUM(obj))
			if obj.Script != nil {
				send_to_char(ch, libc.CString("[T%d] "), obj.Proto_script.Vnum)
			}
		}
		if PRF_FLAGGED(ch, PRF_IHEALTH) {
			send_to_char(ch, libc.CString("@D<@gH@D: @C%d@D>@w %s"), obj.Value[VAL_ALL_HEALTH], obj.Short_description)
		} else {
			send_to_char(ch, libc.CString("%s"), obj.Short_description)
		}
		if int(obj.Type_flag) == ITEM_FOOD {
			if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob {
				send_to_char(ch, libc.CString(", and it has been ate on.@n"))
			}
		}
		if GET_OBJ_VNUM(obj) == math.MaxUint8 {
			switch obj.Value[0] {
			case 0:
				fallthrough
			case 1:
				send_to_char(ch, libc.CString(" @D[@wQuality @RC@D]@n"))
			case 2:
				send_to_char(ch, libc.CString(" @D[@wQuality @RC+@D]@n"))
			case 3:
				send_to_char(ch, libc.CString(" @D[@wQuality @yC++@D]@n"))
			case 4:
				send_to_char(ch, libc.CString(" @D[@wQuality @yB@D]@n"))
			case 5:
				send_to_char(ch, libc.CString(" @D[@wQuality @CB+@D]@n"))
			case 6:
				send_to_char(ch, libc.CString(" @D[@wQuality @CB++@D]@n"))
			case 7:
				send_to_char(ch, libc.CString(" @D[@wQuality @CA@D]@n"))
			case 8:
				send_to_char(ch, libc.CString(" @D[@wQuality @GA+@D]@n"))
			}
		}
		if GET_OBJ_VNUM(obj) == 3424 {
			send_to_char(ch, libc.CString(" @D[@bInk Remaining@D: @w%d@D]@n"), obj.Value[6])
		}
		if GET_OBJ_VNUM(obj) == 3423 {
			send_to_char(ch, libc.CString(" @D[@B%d@D/@B24 Inks@D]@n"), obj.Value[6])
		}
		if OBJ_FLAGGED(obj, ITEM_THROW) {
			send_to_char(ch, libc.CString(" @D[@RThrow Only@D]@n"))
		}
		if int(obj.Type_flag) == ITEM_PLANT && !OBJ_FLAGGED(obj, ITEM_MATURE) {
			if (obj.Value[VAL_WATERLEVEL]) < -9 {
				send_to_char(ch, libc.CString("@D[@RDead@D]@n"))
			} else {
				switch obj.Value[VAL_MATURITY] {
				case 0:
					send_to_char(ch, libc.CString(" @D[@ySeed@D]@n"))
				case 1:
					send_to_char(ch, libc.CString(" @D[@GSprout@D]@n"))
				case 2:
					send_to_char(ch, libc.CString(" @D[@GYoung@D]@n"))
				case 3:
					send_to_char(ch, libc.CString(" @D[@GMature@D]@n"))
				case 4:
					send_to_char(ch, libc.CString(" @D[@GBudding@D]@n"))
				case 5:
					send_to_char(ch, libc.CString("@D[@GClose Harvest@D]@n"))
				case 6:
					send_to_char(ch, libc.CString("@D[@gHarvest@D]@n"))
				}
			}
		}
		if int(obj.Type_flag) == ITEM_CONTAINER && !IS_CORPSE(obj) {
			if !OBJVAL_FLAGGED(obj, 1<<2) && !OBJ_FLAGGED(obj, ITEM_SHEATH) {
				send_to_char(ch, libc.CString(" @D[@G-open-@D]@n"))
			} else if !OBJ_FLAGGED(obj, ITEM_SHEATH) {
				send_to_char(ch, libc.CString(" @D[@rclosed@D]@n"))
			}
		}
		if OBJ_FLAGGED(obj, ITEM_DUPLICATE) {
			send_to_char(ch, libc.CString(" @D[@YDuplicate@D]@n"))
		}
	case SHOW_OBJ_ACTION:
		switch obj.Type_flag {
		case ITEM_NOTE:
			if obj.Action_description != nil {
				var notebuf [6000]byte
				stdio.Snprintf(&notebuf[0], int(6000), "There is something written on it:\r\n\r\n%s", obj.Action_description)
				page_string(ch.Desc, &notebuf[0], TRUE)
			} else {
				send_to_char(ch, libc.CString("There appears to be nothing written on it.\r\n"))
			}
			return
		case ITEM_BOARD:
			show_board(GET_OBJ_VNUM(obj), ch)
		case ITEM_CONTROL:
			send_to_char(ch, libc.CString("@RFUEL@D: %s%s@n\r\n"), func() string {
				if (obj.Value[2]) >= 200 {
					return "@G"
				}
				if (obj.Value[2]) >= 100 {
					return "@Y"
				}
				return "@r"
			}(), add_commas(int64(obj.Value[2])))
		case ITEM_DRINKCON:
			send_to_char(ch, libc.CString("It looks like a drink container.\r\n"))
		case ITEM_LIGHT:
			if (obj.Value[VAL_LIGHT_HOURS]) == -1 {
				send_to_char(ch, libc.CString("Light Cycles left: Infinite\r\n"))
			} else {
				send_to_char(ch, libc.CString("Light Cycles left: [%d]\r\n"), obj.Value[VAL_LIGHT_HOURS])
			}
		case ITEM_FOOD:
			if obj.Foob >= 4 {
				if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob/4 {
					send_to_char(ch, libc.CString("Condition of the food: Almost gone.\r\n"))
				} else if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob/2 {
					send_to_char(ch, libc.CString("Condition of the food: Half Eaten."))
				} else if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob {
					send_to_char(ch, libc.CString("Condition of the food: Partially Eaten."))
				} else if (obj.Value[VAL_FOOD_FOODVAL]) == obj.Foob {
					send_to_char(ch, libc.CString("Condition of the food: Whole."))
				}
			} else if obj.Foob > 0 {
				if (obj.Value[VAL_FOOD_FOODVAL]) < obj.Foob {
					send_to_char(ch, libc.CString("Condition of the food: Almost gone."))
				} else if (obj.Value[VAL_FOOD_FOODVAL]) == obj.Foob {
					send_to_char(ch, libc.CString("Condition of the food: Whole."))
				}
			} else {
				send_to_char(ch, libc.CString("Condition of the food: Insignificant."))
			}
		case ITEM_SPELLBOOK:
			send_to_char(ch, libc.CString("It looks like an arcane tome.\r\n"))
			display_spells(ch, obj)
		case ITEM_SCROLL:
			send_to_char(ch, libc.CString("It looks like an arcane scroll.\r\n"))
			display_scroll(ch, obj)
		case ITEM_VEHICLE:
			if GET_OBJ_VNUM(obj) > 0x4AFF {
				send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @COpen hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CClose hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CEnter hatch\r\n"))
			} else {
				send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @COpen door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CClose door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CEnter door\r\n"))
			}
		case ITEM_HATCH:
			if GET_OBJ_VNUM(obj) > 0x4AFF {
				send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @COpen hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CClose hatch\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CLeave@n\r\n"))
			} else {
				send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @COpen door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CClose door\r\n"))
				send_to_char(ch, libc.CString("@YSyntax@D: @CEnter door\r\n"))
			}
		case ITEM_WINDOW:
			look_out_window(ch, obj.Name)
			return
		default:
			if !IS_CORPSE(obj) {
				send_to_char(ch, libc.CString("You see nothing special..\r\n"))
			} else {
				var mention int = FALSE
				send_to_char(ch, libc.CString("This corpse has "))
				if (obj.Value[VAL_CORPSE_HEAD]) == 0 {
					send_to_char(ch, libc.CString("no head,"))
					mention = TRUE
				}
				if (obj.Value[VAL_CORPSE_RARM]) == 0 {
					send_to_char(ch, libc.CString("no right arm, "))
					mention = TRUE
				} else if (obj.Value[VAL_CORPSE_RARM]) == 2 {
					send_to_char(ch, libc.CString("a broken right arm, "))
					mention = TRUE
				}
				if (obj.Value[VAL_CORPSE_LARM]) == 0 {
					send_to_char(ch, libc.CString("no left arm, "))
					mention = TRUE
				} else if (obj.Value[VAL_CORPSE_LARM]) == 2 {
					send_to_char(ch, libc.CString("a broken left arm, "))
					mention = TRUE
				}
				if (obj.Value[VAL_CORPSE_RLEG]) == 0 {
					send_to_char(ch, libc.CString("no right leg, "))
					mention = TRUE
				} else if (obj.Value[VAL_CORPSE_RLEG]) == 2 {
					send_to_char(ch, libc.CString("a broken right leg, "))
					mention = TRUE
				}
				if (obj.Value[VAL_CORPSE_LLEG]) == 0 {
					send_to_char(ch, libc.CString("no left leg, "))
					mention = TRUE
				} else if (obj.Value[VAL_CORPSE_LLEG]) == 2 {
					send_to_char(ch, libc.CString("a broken left leg, "))
					mention = TRUE
				}
				if mention == FALSE {
					send_to_char(ch, libc.CString("nothing missing from it but life."))
				} else {
					send_to_char(ch, libc.CString("and is dead."))
				}
				send_to_char(ch, libc.CString("\r\n"))
			}
		}
		if int(obj.Type_flag) == ITEM_WEAPON {
			var num int = 0
			if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
				num = 1
			} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
				num = 0
			} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_CRUSH-TYPE_HIT) {
				num = 3
			} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
				num = 2
			} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
				num = 4
			} else {
				num = 5
			}
			send_to_char(ch, libc.CString("The weapon type of %s@n is '%s'.\r\n"), obj.Short_description, weapon_disp[num])
			send_to_char(ch, libc.CString("You could wield it %s.\r\n"), wield_names[wield_type(get_size(ch), obj)])
		}
		diag_obj_to_char(obj, ch)
		send_to_char(ch, libc.CString("It appears to be made of %s, and weighs %s"), material_names[obj.Value[VAL_ALL_MATERIAL]], add_commas(obj.Weight))
	default:
		basic_mud_log(libc.CString("SYSERR: Bad display mode (%d) in show_obj_to_char()."), mode)
		return
	}
	if show_obj_modifiers(obj, ch) != 0 || mode != SHOW_OBJ_ACTION {
		send_to_char(ch, libc.CString("\r\n"))
	}
}
func show_obj_modifiers(obj *obj_data, ch *char_data) int {
	var found int = FALSE
	if OBJ_FLAGGED(obj, ITEM_INVISIBLE) {
		send_to_char(ch, libc.CString(" (invisible)"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_BLESS) && AFF_FLAGGED(ch, AFF_DETECT_ALIGN) {
		send_to_char(ch, libc.CString(" ..It glows blue!"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_MAGIC) && AFF_FLAGGED(ch, AFF_DETECT_MAGIC) {
		send_to_char(ch, libc.CString(" ..It glows yellow!"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_GLOW) {
		send_to_char(ch, libc.CString(" @D(@GGlowing@D)@n"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_HOT) {
		send_to_char(ch, libc.CString(" @D(@RHOT@D)@n"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_HUM) {
		send_to_char(ch, libc.CString(" @D(@RHumming@D)@n"))
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_SLOT2) {
		if OBJ_FLAGGED(obj, ITEM_SLOT_ONE) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString(" @D[@m1/2 Tokens@D]@n"))
		} else if OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString(" @D[@m2/2 Tokens@D]@n"))
		} else {
			send_to_char(ch, libc.CString(" @D[@m0/2 Tokens@D]@n"))
		}
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_SLOT1) {
		if OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString(" @D[@m1/1 Tokens@D]@n"))
		} else {
			send_to_char(ch, libc.CString(" @D[@m0/1 Tokens@D]@n"))
		}
		found++
	}
	if obj.Kicharge > 0 {
		var num int = (obj.Distance * 20) + rand_number(1, 5)
		send_to_char(ch, libc.CString(" %d meters away"), num)
		found++
	}
	if OBJ_FLAGGED(obj, ITEM_CUSTOM) {
		send_to_char(ch, libc.CString(" @D(@YCUSTOM@D)@n"))
	}
	if OBJ_FLAGGED(obj, ITEM_RESTRING) {
		send_to_char(ch, libc.CString(" @D(@R%s@D)@n"), func() *byte {
			if ch.Admlevel > 0 {
				return obj.Name
			}
			return libc.CString("*")
		}())
	}
	if OBJ_FLAGGED(obj, ITEM_BROKEN) {
		if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_STEEL || (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_MITHRIL || (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_METAL {
			send_to_char(ch, libc.CString(", and appears to be twisted and broken."))
		} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_WOOD {
			send_to_char(ch, libc.CString(", and is broken into hundreds of splinters."))
		} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_GLASS {
			send_to_char(ch, libc.CString(", and is shattered on the ground."))
		} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_STONE {
			send_to_char(ch, libc.CString(", and is a pile of rubble."))
		} else {
			send_to_char(ch, libc.CString(", and is broken."))
		}
		found++
	} else {
		if int(obj.Type_flag) != ITEM_BOARD {
			if int(obj.Type_flag) != ITEM_CONTAINER {
				send_to_char(ch, libc.CString("."))
			}
			if !IS_NPC(ch) && obj.Posted_to != nil && obj.Posttype <= 0 {
				var (
					obj2  *obj_data = obj.Posted_to
					dvnum [200]byte
				)
				dvnum[0] = '\x00'
				stdio.Sprintf(&dvnum[0], "@D[@G%d@D] @w", GET_OBJ_VNUM(obj2))
				send_to_char(ch, libc.CString("\n...%s%s has been posted to it."), &func() [200]byte {
					if PRF_FLAGGED(ch, PRF_ROOMFLAGS) {
						return dvnum
					}
					return func() [200]byte {
						var t [200]byte
						copy(t[:], []byte(""))
						return t
					}()
				}()[0], obj2.Short_description)
			}
		}
		found++
	}
	return found
}
func list_obj_to_char(list *obj_data, ch *char_data, mode int, show int) {
	var (
		i     *obj_data
		j     *obj_data
		d     *obj_data
		found bool = FALSE != 0
		num   int
	)
	for i = list; i != nil; i = i.Next_content {
		if i.Description == nil {
			continue
		}
		if libc.StrCaseCmp(i.Description, libc.CString("undefined")) == 0 {
			continue
		}
		num = 0
		d = i
		if config_info.Play.Stack_objs != 0 {
			for j = list; j != i; j = j.Next_content {
				if libc.StrCaseCmp(j.Short_description, i.Short_description) == 0 && libc.StrCaseCmp(j.Description, i.Description) == 0 && j.Item_number == i.Item_number && (OBJ_FLAGGED(j, ITEM_BROKEN) && OBJ_FLAGGED(i, ITEM_BROKEN) || !OBJ_FLAGGED(j, ITEM_BROKEN) && !OBJ_FLAGGED(i, ITEM_BROKEN)) {
					if j.Sitting == nil && i.Sitting == nil {
						if (j.Value[6]) == (i.Value[6]) {
							if int(j.Type_flag) != ITEM_PLANT && int(i.Type_flag) != ITEM_PLANT || int(j.Type_flag) == ITEM_PLANT && int(i.Type_flag) == ITEM_PLANT && (j.Value[VAL_MATURITY]) == (i.Value[VAL_MATURITY]) && (j.Value[VAL_WATERLEVEL]) == (i.Value[VAL_WATERLEVEL]) {
								if !OBJ_FLAGGED(j, ITEM_DUPLICATE) && !OBJ_FLAGGED(i, ITEM_DUPLICATE) || OBJ_FLAGGED(j, ITEM_DUPLICATE) && OBJ_FLAGGED(i, ITEM_DUPLICATE) {
									if j.Posttype == 0 && i.Posttype == 0 {
										if j.Fellow_wall == nil && i.Fellow_wall == nil {
											if (j.Value[0]) == (i.Value[0]) && GET_OBJ_VNUM(j) == math.MaxUint8 && GET_OBJ_VNUM(i) == math.MaxUint8 || GET_OBJ_VNUM(j) != math.MaxUint8 && GET_OBJ_VNUM(i) != math.MaxUint8 {
												break
											}
										}
									}
								}
							}
						}
					}
				}
			}
			if j != i {
				continue
			}
			for d = func() *obj_data {
				j = i
				return j
			}(); j != nil; j = j.Next_content {
				if libc.StrCaseCmp(j.Short_description, i.Short_description) == 0 && libc.StrCaseCmp(j.Description, i.Description) == 0 && j.Item_number == i.Item_number && (OBJ_FLAGGED(j, ITEM_BROKEN) && OBJ_FLAGGED(i, ITEM_BROKEN) || !OBJ_FLAGGED(j, ITEM_BROKEN) && !OBJ_FLAGGED(i, ITEM_BROKEN)) {
					if j.Sitting == nil && i.Sitting == nil {
						if j.Posttype == 0 && i.Posttype == 0 {
							if (j.Value[6]) == (i.Value[6]) {
								if int(j.Type_flag) != ITEM_PLANT && int(i.Type_flag) != ITEM_PLANT || int(j.Type_flag) == ITEM_PLANT && int(i.Type_flag) == ITEM_PLANT && (j.Value[VAL_MATURITY]) == (i.Value[VAL_MATURITY]) && (j.Value[VAL_WATERLEVEL]) == (i.Value[VAL_WATERLEVEL]) {
									if !OBJ_FLAGGED(j, ITEM_DUPLICATE) && !OBJ_FLAGGED(i, ITEM_DUPLICATE) || OBJ_FLAGGED(j, ITEM_DUPLICATE) && OBJ_FLAGGED(i, ITEM_DUPLICATE) {
										if j.Fellow_wall == nil && i.Fellow_wall == nil {
											if (j.Value[0]) == (i.Value[0]) && GET_OBJ_VNUM(i) == math.MaxUint8 && GET_OBJ_VNUM(j) == math.MaxUint8 || GET_OBJ_VNUM(j) != math.MaxUint8 && GET_OBJ_VNUM(i) != math.MaxUint8 {
												if CAN_SEE_OBJ(ch, j) {
													num++
													if d == i && !CAN_SEE_OBJ(ch, d) {
														d = j
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if CAN_SEE_OBJ(ch, d) && (*d.Description != '.' && *d.Short_description != '.' || PRF_FLAGGED(ch, PRF_HOLYLIGHT)) || int(d.Type_flag) == ITEM_LIGHT {
			if num > 1 {
				send_to_char(ch, libc.CString("@D(@Rx@Y%2i@D)@n "), num)
			}
			show_obj_to_char(d, ch, mode)
			found = TRUE != 0
		}
	}
	if !found && show != 0 {
		send_to_char(ch, libc.CString(" Nothing.\r\n"))
	}
}
func diag_obj_to_char(obj *obj_data, ch *char_data) {
	var (
		diagnosis [8]struct {
			Percent int
			Text    *byte
		} = [8]struct {
			Percent int
			Text    *byte
		}{{Percent: 100, Text: libc.CString("is in excellent condition.")}, {Percent: 90, Text: libc.CString("has a few scuffs.")}, {Percent: 75, Text: libc.CString("has some small scuffs and scratches.")}, {Percent: 50, Text: libc.CString("has quite a few scratches.")}, {Percent: 30, Text: libc.CString("has some big nasty scrapes and scratches.")}, {Percent: 15, Text: libc.CString("looks pretty damaged.")}, {Percent: 0, Text: libc.CString("is in awful condition.")}, {Percent: -1, Text: libc.CString("is in need of repair.")}}
		percent  int
		ar_index int
		objs     *byte = OBJS(obj, ch)
	)
	if (obj.Value[VAL_ALL_MAXHEALTH]) > 0 {
		percent = ((obj.Value[VAL_ALL_HEALTH]) * 100) / (obj.Value[VAL_ALL_MAXHEALTH])
	} else {
		percent = 0
	}
	for ar_index = 0; diagnosis[ar_index].Percent >= 0; ar_index++ {
		if percent >= diagnosis[ar_index].Percent {
			break
		}
	}
	send_to_char(ch, libc.CString("\r\n%c%s %s\r\n"), unicode.ToUpper(rune(*objs)), (*byte)(unsafe.Add(unsafe.Pointer(objs), 1)), diagnosis[ar_index].Text)
}
func diag_char_to_char(i *char_data, ch *char_data) {
	var (
		diagnosis [12]struct {
			Percent int
			Text    *byte
		} = [12]struct {
			Percent int
			Text    *byte
		}{{Percent: 100, Text: libc.CString("@wis in @Gexcellent@w condition.")}, {Percent: 90, Text: libc.CString("@whas a few @Rscratches@w.")}, {Percent: 80, Text: libc.CString("@whas some small @Rwounds@w and @Rbruises@w.")}, {Percent: 70, Text: libc.CString("@whas quite a few @Rwounds@w.")}, {Percent: 60, Text: libc.CString("@whas some big @rnasty wounds@w and @Rscratches@w.")}, {Percent: 50, Text: libc.CString("@wlooks pretty @rhurt@w.")}, {Percent: 40, Text: libc.CString("@wis mainly @rinjured@w.")}, {Percent: 30, Text: libc.CString("@wis a @rmess@w of @rinjuries@w.")}, {Percent: 20, Text: libc.CString("@wis @rstruggling@w to @msurvive@w.")}, {Percent: 10, Text: libc.CString("@wis in @mawful condition@w.")}, {Percent: 0, Text: libc.CString("@Ris barely alive.@w")}, {Percent: -1, Text: libc.CString("@ris nearly dead.@w")}}
		percent  int
		ar_index int
		hit      int64 = i.Hit
		max      int64 = gear_pl(i)
		total    int64 = 0
	)
	if i.Suppression > 0 {
		total = max - i.Suppressed
	} else {
		total = max
	}
	if hit == total {
		percent = 100
	} else if hit < total && float64(hit) >= (float64(total)*0.9) {
		percent = 90
	} else if hit < total && float64(hit) >= (float64(total)*0.8) {
		percent = 80
	} else if hit < total && float64(hit) >= (float64(total)*0.7) {
		percent = 70
	} else if hit < total && float64(hit) >= (float64(total)*0.6) {
		percent = 60
	} else if hit < total && float64(hit) >= (float64(total)*0.5) {
		percent = 50
	} else if hit < total && float64(hit) >= (float64(total)*0.4) {
		percent = 40
	} else if hit < total && float64(hit) >= (float64(total)*0.3) {
		percent = 30
	} else if hit < total && float64(hit) >= (float64(total)*0.2) {
		percent = 20
	} else if hit < total && float64(hit) >= (float64(total)*0.1) {
		percent = 10
	} else if float64(hit) < float64(total)*0.1 {
		percent = 0
	} else {
		percent = -1
	}
	for ar_index = 0; diagnosis[ar_index].Percent >= 0; ar_index++ {
		if percent >= diagnosis[ar_index].Percent {
			break
		}
	}
	send_to_char(ch, libc.CString("%s\r\n"), diagnosis[ar_index].Text)
}
func look_at_char(i *char_data, ch *char_data) {
	var (
		j       int
		found   int
		clan    int = FALSE
		buf     [100]byte
		tmp_obj *obj_data
	)
	if ch.Desc == nil {
		return
	}
	if i.Description != nil {
		send_to_char(ch, libc.CString("%s"), i.Description)
	}
	if !MOB_FLAGGED(i, MOB_JUSTDESC) {
		bringdesc(ch, i)
	}
	send_to_char(ch, libc.CString("\r\n"))
	if !IS_NPC(i) {
		if (i.Limb_condition[0]) >= 50 && !PLR_FLAGGED(i, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), i.Limb_condition[0], "%", "%")
		} else if (i.Limb_condition[0]) > 0 && !PLR_FLAGGED(i, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), i.Limb_condition[0], "%", "%")
		} else if (i.Limb_condition[0]) > 0 && PLR_FLAGGED(i, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), i.Limb_condition[0], "%", "%")
		} else if (i.Limb_condition[0]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @rMissing.            @D]@n\r\n"))
		}
		if (i.Limb_condition[1]) >= 50 && !PLR_FLAGGED(i, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), i.Limb_condition[1], "%", "%")
		} else if (i.Limb_condition[1]) > 0 && !PLR_FLAGGED(i, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), i.Limb_condition[1], "%", "%")
		} else if (i.Limb_condition[1]) > 0 && PLR_FLAGGED(i, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), i.Limb_condition[1], "%", "%")
		} else if (i.Limb_condition[1]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @rMissing.            @D]@n\r\n"))
		}
		if (i.Limb_condition[2]) >= 50 && !PLR_FLAGGED(i, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), i.Limb_condition[2], "%", "%")
		} else if (i.Limb_condition[2]) > 0 && !PLR_FLAGGED(i, PLR_CRLEG) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), i.Limb_condition[2], "%", "%")
		} else if (i.Limb_condition[2]) > 0 && PLR_FLAGGED(i, PLR_CRLEG) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), i.Limb_condition[2], "%", "%")
		} else if (i.Limb_condition[2]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @rMissing.            @D]@n\r\n"))
		}
		if (i.Limb_condition[3]) >= 50 && !PLR_FLAGGED(i, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), i.Limb_condition[3], "%", "%")
		} else if (i.Limb_condition[3]) > 0 && !PLR_FLAGGED(i, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), i.Limb_condition[3], "%", "%")
		} else if (i.Limb_condition[3]) > 0 && PLR_FLAGGED(i, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), i.Limb_condition[3], "%", "%")
		} else if (i.Limb_condition[3]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @rMissing.             @D]@n\r\n"))
		}
		if PLR_FLAGGED(i, PLR_HEAD) {
			send_to_char(ch, libc.CString("            @D[@cHead        @D: @GHas.                 @D]@n\r\n"))
		}
		if !PLR_FLAGGED(i, PLR_HEAD) {
			send_to_char(ch, libc.CString("            @D[@cHead        @D: @rMissing.             @D]@n\r\n"))
		}
		if (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) && PLR_FLAGGED(i, PLR_STAIL) && !PLR_FLAGGED(i, PLR_TAILHIDE) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @GHas.                 @D]@n\r\n"))
		}
		if (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) && !PLR_FLAGGED(i, PLR_STAIL) && !PLR_FLAGGED(i, PLR_TAILHIDE) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @rMissing.             @D]@n\r\n"))
		}
		if (int(i.Race) == RACE_ICER || int(i.Race) == RACE_BIO) && PLR_FLAGGED(i, PLR_TAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @GHas.                 @D]@n\r\n"))
		}
		if (int(i.Race) == RACE_ICER || int(i.Race) == RACE_BIO) && !PLR_FLAGGED(i, PLR_TAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @rMissing.             @D]@n\r\n"))
		}
	}
	send_to_char(ch, libc.CString("\r\n"))
	if i.Clan != nil && unsafe.Pointer(libc.StrStr(i.Clan, libc.CString("None"))) == unsafe.Pointer(uintptr(FALSE)) {
		stdio.Sprintf(&buf[0], "%s", i.Clan)
		clan = TRUE
	}
	if i.Clan == nil {
		clan = FALSE
	}
	if !IS_NPC(i) {
		send_to_char(ch, libc.CString("            @D[@mClan        @D: @W%-20s@D]@n\r\n"), &func() [100]byte {
			if clan != 0 {
				return buf
			}
			return func() [100]byte {
				var t [100]byte
				copy(t[:], []byte("None."))
				return t
			}()
		}()[0])
	}
	if !IS_NPC(i) {
		send_to_char(ch, libc.CString("\r\n         @D----------------------------------------@n\r\n"))
		trans_check(ch, i)
		send_to_char(ch, libc.CString("         @D----------------------------------------@n\r\n"))
	}
	send_to_char(ch, libc.CString("\r\n"))
	if !PLR_FLAGGED(i, PLR_DISGUISED) && (readIntro(ch, i) == 1 && !IS_NPC(i)) {
		if int(i.Sex) == SEX_NEUTRAL {
			send_to_char(ch, libc.CString("%s appears to be %s %s, "), get_i_name(ch, i), AN(JUGGLERACE(i)), JUGGLERACELOWER(i))
		} else {
			send_to_char(ch, libc.CString("%s appears to be %s %s %s, "), get_i_name(ch, i), AN(MAFE(i)), MAFE(i), JUGGLERACELOWER(i))
		}
	} else if ch == i || IS_NPC(i) {
		if int(i.Sex) == SEX_NEUTRAL {
			send_to_char(ch, libc.CString("%c%s appears to be %s %s, "), unicode.ToUpper(rune(*GET_NAME(i))), (*byte)(unsafe.Add(unsafe.Pointer(GET_NAME(i)), 1)), AN(JUGGLERACE(i)), JUGGLERACELOWER(i))
		} else {
			send_to_char(ch, libc.CString("%c%s appears to be %s %s %s, "), unicode.ToUpper(rune(*GET_NAME(i))), (*byte)(unsafe.Add(unsafe.Pointer(GET_NAME(i)), 1)), AN(MAFE(i)), MAFE(i), JUGGLERACELOWER(i))
		}
	} else {
		if int(i.Sex) == SEX_NEUTRAL {
			send_to_char(ch, libc.CString("Appears to be %s %s, "), AN(JUGGLERACE(i)), JUGGLERACELOWER(i))
		} else {
			send_to_char(ch, libc.CString("Appears to be %s %s %s, "), AN(MAFE(i)), MAFE(i), JUGGLERACELOWER(i))
		}
	}
	if IS_NPC(i) {
		send_to_char(ch, libc.CString("is %s sized, and\r\n"), size_names[get_size(i)])
	}
	if !IS_NPC(i) {
		if !PLR_FLAGGED(i, PLR_OOZARU) && (int(i.Race) != RACE_ICER || !IS_TRANSFORMED(i)) && (i.Genome[0]) < 11 {
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], GET_PC_HEIGHT(i), GET_PC_WEIGHT(i))
		} else if int(i.Race) == RACE_ICER && PLR_FLAGGED(i, PLR_TRANS1) {
			var (
				num1 int = GET_PC_HEIGHT(i) * 3
				num2 int = GET_PC_WEIGHT(i) * 4
			)
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], num1, num2)
		} else if int(i.Race) == RACE_ICER && PLR_FLAGGED(i, PLR_TRANS2) {
			var (
				num1 int = GET_PC_HEIGHT(i) * 3
				num2 int = GET_PC_WEIGHT(i) * 4
			)
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], num1, num2)
		} else if int(i.Race) == RACE_ICER && PLR_FLAGGED(i, PLR_TRANS3) {
			var (
				num1 int = int(float64(GET_PC_HEIGHT(i)) * 1.5)
				num2 int = GET_PC_WEIGHT(i) * 2
			)
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], num1, num2)
		} else if int(i.Race) == RACE_ICER && PLR_FLAGGED(i, PLR_TRANS4) {
			var (
				num1 int = GET_PC_HEIGHT(i) * 2
				num2 int = GET_PC_WEIGHT(i) * 3
			)
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], num1, num2)
		} else if PLR_FLAGGED(i, PLR_OOZARU) || (i.Genome[0]) == 11 {
			var (
				num1 int = GET_PC_HEIGHT(i) * 10
				num2 int = GET_PC_WEIGHT(i) * 50
			)
			send_to_char(ch, libc.CString("is %s sized, about %dcm tall,\r\nabout %dkg heavy,"), size_names[get_size(i)], num1, num2)
		}
		if i == ch {
			send_to_char(ch, libc.CString(" and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+30 {
			send_to_char(ch, libc.CString(" appears to be very much younger than you, and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+25 {
			send_to_char(ch, libc.CString(" appears to be much younger than you, and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+15 {
			send_to_char(ch, libc.CString(" appears to be a good amount younger than you, and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+10 {
			send_to_char(ch, libc.CString(" appears to be about a decade younger than you, and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+5 {
			send_to_char(ch, libc.CString(" appears to be several years younger than you, and "))
		} else if int(age(ch).Year) >= int(age(i).Year)+2 {
			send_to_char(ch, libc.CString(" appears to be a bit younger than you, and "))
		} else if int(age(ch).Year) > int(age(i).Year) {
			send_to_char(ch, libc.CString(" appears to be slightly younger than you, and "))
		} else if int(age(ch).Year) == int(age(i).Year) {
			send_to_char(ch, libc.CString(" appears to be the same age as you, and "))
		}
		if int(age(i).Year) >= int(age(ch).Year)+30 {
			send_to_char(ch, libc.CString(" appears to be very much older than you, and "))
		} else if int(age(i).Year) >= int(age(ch).Year)+25 {
			send_to_char(ch, libc.CString(" appears to be much older than you, and "))
		} else if int(age(i).Year) >= int(age(ch).Year)+15 {
			send_to_char(ch, libc.CString(" appears to be a good amount older than you, and "))
		} else if int(age(i).Year) >= int(age(ch).Year)+10 {
			send_to_char(ch, libc.CString(" appears to be about a decade older than you, and "))
		} else if int(age(i).Year) >= int(age(ch).Year)+5 {
			send_to_char(ch, libc.CString(" appears to be several years older than you, and "))
		} else if int(age(i).Year) >= int(age(ch).Year)+2 {
			send_to_char(ch, libc.CString(" appears to be a bit older than you, and "))
		} else if int(age(i).Year) > int(age(ch).Year) {
			send_to_char(ch, libc.CString(" appears to be slightly older than you, and "))
		}
	}
	diag_char_to_char(i, ch)
	found = FALSE
	for j = 0; found == 0 && j < NUM_WEARS; j++ {
		if (i.Equipment[j]) != nil && CAN_SEE_OBJ(ch, i.Equipment[j]) {
			found = TRUE
		}
	}
	if found != 0 && (!IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NOEQSEE)) {
		send_to_char(ch, libc.CString("\r\n"))
		if !PLR_FLAGGED(i, PLR_DISGUISED) {
			act(libc.CString("$n is using:"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else {
			act(libc.CString("The disguised person is using:"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		for j = 0; j < NUM_WEARS; j++ {
			if (i.Equipment[j]) != nil && CAN_SEE_OBJ(ch, i.Equipment[j]) && (j != WEAR_WIELD1 && j != WEAR_WIELD2) {
				send_to_char(ch, libc.CString("%s"), wear_where[j])
				show_obj_to_char(i.Equipment[j], ch, SHOW_OBJ_SHORT)
				if OBJ_FLAGGED(i.Equipment[j], ITEM_SHEATH) {
					var (
						obj2     *obj_data = nil
						next_obj *obj_data = nil
						sheath   *obj_data = (i.Equipment[j])
					)
					for obj2 = sheath.Contains; obj2 != nil; obj2 = next_obj {
						next_obj = obj2.Next_content
						if obj2 != nil {
							send_to_char(ch, libc.CString("@D  ---- @YSheathed@D ----@c> @n"))
							show_obj_to_char(obj2, ch, SHOW_OBJ_SHORT)
						}
					}
					obj2 = nil
				}
			} else if (i.Equipment[j]) != nil && CAN_SEE_OBJ(ch, i.Equipment[j]) && !PLR_FLAGGED(i, PLR_THANDW) {
				send_to_char(ch, libc.CString("%s"), wear_where[j])
				show_obj_to_char(i.Equipment[j], ch, SHOW_OBJ_SHORT)
				if OBJ_FLAGGED(i.Equipment[j], ITEM_SHEATH) {
					var (
						obj2     *obj_data = nil
						next_obj *obj_data = nil
						sheath   *obj_data = (i.Equipment[j])
					)
					for obj2 = sheath.Contains; obj2 != nil; obj2 = next_obj {
						next_obj = obj2.Next_content
						if obj2 != nil {
							send_to_char(ch, libc.CString("@D  ---- @YSheathed@D ----@c> @n"))
							show_obj_to_char(obj2, ch, SHOW_OBJ_SHORT)
						}
					}
					obj2 = nil
				}
			} else if (i.Equipment[j]) != nil && CAN_SEE_OBJ(ch, i.Equipment[j]) && PLR_FLAGGED(i, PLR_THANDW) {
				send_to_char(ch, libc.CString("@c<@CWielded by B. Hands@c>@n "))
				show_obj_to_char(i.Equipment[j], ch, SHOW_OBJ_SHORT)
			}
		}
	}
	if ch != i && (GET_SKILL(ch, SKILL_KEEN) != 0 && AFF_FLAGGED(ch, AFF_SNEAK) || ch.Admlevel != 0) {
		found = FALSE
		act(libc.CString("\r\nYou attempt to peek at $s inventory:"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		if CAN_SEE(i, ch) {
			act(libc.CString("$n tries to evaluate what you have in your inventory."), TRUE, ch, nil, unsafe.Pointer(i), TO_VICT)
		}
		if GET_SKILL(ch, SKILL_KEEN) > axion_dice(0) && (!IS_NPC(i) || ch.Admlevel > 1) {
			for tmp_obj = i.Carrying; tmp_obj != nil; tmp_obj = tmp_obj.Next_content {
				if CAN_SEE_OBJ(ch, tmp_obj) && (ADM_FLAGGED(ch, ADM_SEEINV) || rand_number(0, 20) < GET_LEVEL(ch)) {
					show_obj_to_char(tmp_obj, ch, SHOW_OBJ_SHORT)
					found = TRUE
				}
			}
			improve_skill(ch, SKILL_KEEN, 1)
		} else if IS_NPC(i) && ch.Admlevel < 2 {
			return
		} else {
			act(libc.CString("You are unsure about $s inventory."), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
			if CAN_SEE(i, ch) {
				act(libc.CString("$n didn't seem to get a good enough look."), TRUE, ch, nil, unsafe.Pointer(i), TO_VICT)
			}
			improve_skill(ch, SKILL_KEEN, 1)
			return
		}
		if found == 0 {
			send_to_char(ch, libc.CString("You can't see anything.\r\n"))
			improve_skill(ch, SKILL_KEEN, 1)
		}
	}
}
func list_one_char(i *char_data, ch *char_data) {
	var (
		chair     *obj_data = nil
		count     int       = FALSE
		positions [9]*byte  = [9]*byte{libc.CString(" is dead"), libc.CString(" is mortally wounded"), libc.CString(" is lying here, incapacitated"), libc.CString(" is lying here, stunned"), libc.CString(" is sleeping here"), libc.CString(" is resting here"), libc.CString(" is sitting here"), libc.CString("!FIGHTING!"), libc.CString(" is standing here")}
	)
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_ROOMFLAGS) && IS_NPC(i) {
		send_to_char(ch, libc.CString("@D[@G%d@D]@w %s"), GET_MOB_VNUM(i), func() string {
			if i.Script != nil {
				return "[TRIG] "
			}
			return ""
		}())
	}
	if IS_NPC(i) && i.Long_descr != nil && int(i.Position) == int(i.Mob_specials.Default_pos) && i.Fighting == nil {
		send_to_char(ch, libc.CString("%s"), i.Long_descr)
		if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.9 && i.Hit != gear_pl(i) {
			act(libc.CString("@R...Some slight wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.8 && float64(i.Hit) < float64(gear_pl(i))*0.9 {
			act(libc.CString("@R...A few wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.7 && float64(i.Hit) < float64(gear_pl(i))*0.8 {
			act(libc.CString("@R...Many wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.6 && float64(i.Hit) < float64(gear_pl(i))*0.7 {
			act(libc.CString("@R...Quite a few wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.5 && float64(i.Hit) < float64(gear_pl(i))*0.6 {
			act(libc.CString("@R...Horrible wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.4 && float64(i.Hit) < float64(gear_pl(i))*0.5 {
			act(libc.CString("@R...Blood is seeping from the wounds on $s body.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.3 && float64(i.Hit) < float64(gear_pl(i))*0.4 {
			act(libc.CString("@R...$s body is in terrible shape.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.2 && float64(i.Hit) < float64(gear_pl(i))*0.3 {
			act(libc.CString("@R...Is absolutely covered in wounds.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) >= float64(gear_pl(i))*0.1 && float64(i.Hit) < float64(gear_pl(i))*0.2 {
			act(libc.CString("@R...Is on $s last leg.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		} else if IS_NPC(i) && float64(i.Hit) < float64(gear_pl(i))*0.1 {
			act(libc.CString("@R...Should be DEAD soon.@w"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if i.Listenroom > 0 {
			var eaves [300]byte
			stdio.Sprintf(&eaves[0], "@w...$e is spying on everything to the @c%s@w.", dirs[i.Eavesdir])
			act(&eaves[0], TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_FLYING) && i.Altitude == 1 {
			act(libc.CString("...$e is in the air!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_FLYING) && i.Altitude == 2 {
			act(libc.CString("...$e is high in the air!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_SANCTUARY) && GET_SKILL(i, SKILL_AQUA_BARRIER) == 0 {
			act(libc.CString("...$e has a barrier around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_FIRESHIELD) {
			act(libc.CString("...$e has @rf@Rl@Ya@rm@Re@Ys@w around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_SANCTUARY) && GET_SKILL(i, SKILL_AQUA_BARRIER) != 0 {
			act(libc.CString("...$e has a @Gbarrier@w of @cwater@w and @Cki@w around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if !IS_NPC(i) && PLR_FLAGGED(i, PLR_SPIRAL) {
			act(libc.CString("...$e is spinning in a vortex!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if i.Charge != 0 {
			act(libc.CString("...$e has a bright %s aura around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_METAMORPH) {
			act(libc.CString("@w...$e has a dark, @rred@w aura and menacing presence."), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_HAYASA) {
			act(libc.CString("@w...$e has a soft @cblue@w glow around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if AFF_FLAGGED(i, AFF_BLIND) {
			act(libc.CString("...$e is groping around blindly!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if affected_by_spell(i, SPELL_FAERIE_FIRE) {
			act(libc.CString("@m...$e @mis outlined with purple fire!@m"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		if i.Feature != nil {
			var woo [64936]byte
			stdio.Sprintf(&woo[0], "@C%s@n", i.Feature)
			act(&woo[0], FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
		return
	}
	if IS_NPC(i) && i.Fighting == nil && int(i.Position) != POS_SITTING && int(i.Position) != POS_SLEEPING {
		send_to_char(ch, libc.CString("@w%c%s"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Grappled != nil && i.Grappled == ch {
		send_to_char(ch, libc.CString("@w%c%s is being grappled with by YOU!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Grappled != nil && i.Grappled != ch {
		send_to_char(ch, libc.CString("@w%c%s is being absorbed from by %s!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)), func() *byte {
			if readIntro(ch, i.Grappled) == 1 {
				return get_i_name(ch, i.Grappled)
			}
			return AN(JUGGLERACE(i.Grappled))
		}())
	} else if IS_NPC(i) && i.Absorbby != nil && i.Absorbby == ch {
		send_to_char(ch, libc.CString("@w%c%s is being absorbed from by YOU!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Absorbby != nil && i.Absorbby != ch {
		send_to_char(ch, libc.CString("@w%c%s is being absorbed from by %s!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)), func() *byte {
			if readIntro(ch, i.Absorbby) == 1 {
				return get_i_name(ch, i.Absorbby)
			}
			return AN(JUGGLERACE(i.Absorbby))
		}())
	} else if IS_NPC(i) && i.Fighting != nil && i.Fighting != ch && int(i.Position) != POS_SITTING && int(i.Position) != POS_SLEEPING && is_sparring(i) {
		send_to_char(ch, libc.CString("@w%c%s is sparring with %s!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)), func() *byte {
			if ch.Admlevel != 0 {
				return GET_NAME(i.Fighting)
			}
			if readIntro(ch, i.Fighting) == 1 {
				return get_i_name(ch, i.Fighting)
			}
			return JUGGLERACELOWER(i.Fighting)
		}())
	} else if IS_NPC(i) && i.Fighting != nil && is_sparring(i) && i.Fighting == ch && int(i.Position) != POS_SITTING && int(i.Position) != POS_SLEEPING {
		send_to_char(ch, libc.CString("@w%c%s is sparring with you!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Fighting != nil && i.Fighting != ch && int(i.Position) != POS_SITTING && int(i.Position) != POS_SLEEPING {
		send_to_char(ch, libc.CString("@w%c%s is fighting %s!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)), func() *byte {
			if ch.Admlevel != 0 {
				return GET_NAME(i.Fighting)
			}
			if readIntro(ch, i.Fighting) == 1 {
				return get_i_name(ch, i.Fighting)
			}
			return JUGGLERACELOWER(i.Fighting)
		}())
	} else if IS_NPC(i) && i.Fighting != nil && i.Fighting == ch && int(i.Position) != POS_SITTING && int(i.Position) != POS_SLEEPING {
		send_to_char(ch, libc.CString("@w%c%s is fighting YOU!"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Fighting != nil && int(i.Position) == POS_SITTING {
		send_to_char(ch, libc.CString("@w%c%s is sitting here."), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) && i.Fighting != nil && int(i.Position) == POS_SLEEPING {
		send_to_char(ch, libc.CString("@w%c%s is sleeping here."), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if IS_NPC(i) {
		send_to_char(ch, libc.CString("@w%c%s"), unicode.ToUpper(rune(*i.Short_descr)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_descr), 1)))
	} else if !IS_NPC(i) {
		if int(i.Race) == RACE_MAJIN && AFF_FLAGGED(i, AFF_LIQUEFIED) {
			send_to_char(ch, libc.CString("@wSeveral blobs of %s colored goo spread out here.@n\n"), skin_types[int(i.Skin)])
			return
		}
		if ch.Admlevel > 0 || i.Admlevel > 0 || IS_NPC(ch) {
			send_to_char(ch, libc.CString("@w%s"), i.Name)
		} else if !PLR_FLAGGED(i, PLR_DISGUISED) && readIntro(ch, i) == 1 {
			send_to_char(ch, libc.CString("@w%s"), get_i_name(ch, i))
		} else if !PLR_FLAGGED(i, PLR_DISGUISED) && readIntro(ch, i) != 1 {
			if int(i.Distfea) == DISTFEA_EYE {
				send_to_char(ch, libc.CString("@wA %s eyed %s %s"), eye_types[int(i.Eye)], MAFE(i), JUGGLERACELOWER(i))
			} else if int(i.Distfea) == DISTFEA_HAIR {
				if int(i.Race) == RACE_MAJIN {
					send_to_char(ch, libc.CString("@wA %s majin, with a %s forelock,"), MAFE(i), FHA_types[int(i.Hairl)])
				} else if int(i.Race) == RACE_NAMEK {
					send_to_char(ch, libc.CString("@wA namek, with %s antennae,"), FHA_types[int(i.Hairl)])
				} else if int(i.Race) == RACE_ARLIAN {
					send_to_char(ch, libc.CString("@wA arlian, with %s antennae,"), FHA_types[int(i.Hairl)])
				} else if int(i.Race) == RACE_ICER || int(i.Race) == RACE_DEMON {
					send_to_char(ch, libc.CString("@wA %s %s, with %s horns"), MAFE(i), JUGGLERACELOWER(i), FHA_types[int(i.Hairl)])
				} else {
					var blarg [2048]byte
					stdio.Sprintf(&blarg[0], "%s %s hair %s", hairl_types[int(i.Hairl)], hairc_types[int(i.Hairc)], hairs_types[int(i.Hairs)])
					send_to_char(ch, libc.CString("@wA %s %s, with %s"), MAFE(i), JUGGLERACELOWER(i), func() string {
						if int(i.Hairl) == 0 {
							return "a bald head"
						}
						return libc.GoString(&blarg[0])
					}())
				}
			} else if int(i.Distfea) == DISTFEA_SKIN {
				send_to_char(ch, libc.CString("@wA %s skinned %s %s"), skin_types[int(i.Skin)], MAFE(i), JUGGLERACELOWER(i))
			} else if int(i.Distfea) == DISTFEA_HEIGHT {
				var height *byte
				if int(i.Race) == RACE_TRUFFLE {
					if GET_PC_HEIGHT(i) > 70 {
						height = libc.CString("very tall")
					} else if GET_PC_HEIGHT(i) > 55 {
						height = libc.CString("tall")
					} else if GET_PC_HEIGHT(i) > 35 {
						height = libc.CString("average height")
					} else {
						height = libc.CString("short")
					}
				} else if PLR_FLAGGED(i, PLR_OOZARU) || (i.Genome[0]) == 11 {
					if GET_PC_HEIGHT(i)*10 > 2000 {
						height = libc.CString("very tall")
					} else if GET_PC_HEIGHT(i)*10 > 1800 {
						height = libc.CString("tall")
					} else if GET_PC_HEIGHT(i)*10 > 1500 {
						height = libc.CString("average height")
					} else {
						height = libc.CString("short")
					}
				} else {
					if GET_PC_HEIGHT(i) > 200 {
						height = libc.CString("very tall")
					} else if GET_PC_HEIGHT(i) > 180 {
						height = libc.CString("tall")
					} else if GET_PC_HEIGHT(i) > 150 {
						height = libc.CString("average height")
					} else if GET_PC_HEIGHT(i) > 120 {
						height = libc.CString("short")
					} else {
						height = libc.CString("very short")
					}
				}
				send_to_char(ch, libc.CString("@wA %s %s %s"), height, MAFE(i), JUGGLERACELOWER(i))
				if height != nil {
					libc.Free(unsafe.Pointer(height))
				}
			} else if int(i.Distfea) == DISTFEA_WEIGHT {
				var height *byte
				if int(i.Race) == RACE_TRUFFLE {
					if GET_PC_WEIGHT(i) > 35 {
						height = libc.CString("very heavy")
					} else if GET_PC_WEIGHT(i) > 25 {
						height = libc.CString("heavy")
					} else if GET_PC_WEIGHT(i) > 15 {
						height = libc.CString("average weight")
					} else {
						height = libc.CString("welterweight")
					}
				} else if PLR_FLAGGED(i, PLR_OOZARU) || (i.Genome[0]) == 11 {
					if GET_PC_WEIGHT(i)*50 > 6000 {
						height = libc.CString("very heavy")
					} else if GET_PC_WEIGHT(i)*50 > 5000 {
						height = libc.CString("heavy")
					} else if GET_PC_WEIGHT(i)*50 > 4000 {
						height = libc.CString("average weight")
					} else if GET_PC_WEIGHT(i)*50 > 3000 {
						height = libc.CString("lightweight")
					} else {
						height = libc.CString("welterweight")
					}
				} else {
					if GET_PC_WEIGHT(i) > 120 {
						height = libc.CString("very heavy")
					} else if GET_PC_WEIGHT(i) > 100 {
						height = libc.CString("heavy")
					} else if GET_PC_WEIGHT(i) > 80 {
						height = libc.CString("average weight")
					} else if GET_PC_WEIGHT(i) > 60 {
						height = libc.CString("lightweight")
					} else {
						height = libc.CString("welterweight")
					}
				}
				send_to_char(ch, libc.CString("@wA %s %s %s"), height, MAFE(i), JUGGLERACELOWER(i))
				if height != nil {
					libc.Free(unsafe.Pointer(height))
				}
			}
		} else {
			send_to_char(ch, libc.CString("@wA disguised %s %s"), MAFE(i), JUGGLERACELOWER(i))
		}
	}
	if !IS_NPC(i) || i.Fighting == nil {
		if AFF_FLAGGED(i, AFF_INVISIBLE) {
			send_to_char(ch, libc.CString(", is invisible"))
			count = TRUE
		}
		if AFF_FLAGGED(i, AFF_ETHEREAL) {
			send_to_char(ch, libc.CString(", has a halo"))
			count = TRUE
		}
		if AFF_FLAGGED(i, AFF_HIDE) && i != ch {
			send_to_char(ch, libc.CString(", is hiding"))
			if GET_SKILL(i, SKILL_HIDE) != 0 && !IS_NPC(ch) && i != ch {
				improve_skill(i, SKILL_HIDE, 1)
			}
			count = TRUE
		}
		if !IS_NPC(i) && i.Desc == nil {
			send_to_char(ch, libc.CString(", has a blank stare"))
			count = TRUE
		}
		if !IS_NPC(i) && PLR_FLAGGED(i, PLR_WRITING) {
			send_to_char(ch, libc.CString(", is writing"))
			count = TRUE
		}
		if !IS_NPC(i) && PRF_FLAGGED(i, PRF_BUILDWALK) {
			send_to_char(ch, libc.CString(", is buildwalking"))
			count = TRUE
		}
		if !IS_NPC(i) && i.Absorbing != nil && i.Absorbing != ch {
			send_to_char(ch, libc.CString(", is absorbing from %s"), GET_NAME(i.Absorbing))
			count = TRUE
		}
		if !IS_NPC(i) && i.Grappling != nil && i.Grappling != ch {
			send_to_char(ch, libc.CString(", is grappling with %s"), func() *byte {
				if readIntro(ch, i.Grappling) == 1 {
					return get_i_name(ch, i.Grappling)
				}
				return introd_calc(i.Grappling)
			}())
			count = TRUE
		}
		if !IS_NPC(i) && i.Player_specials.Carrying != nil && i.Player_specials.Carrying != ch {
			send_to_char(ch, libc.CString(", is carrying %s"), func() *byte {
				if readIntro(ch, i.Player_specials.Carrying) == 1 {
					return get_i_name(ch, i.Player_specials.Carrying)
				}
				return introd_calc(i.Player_specials.Carrying)
			}())
			count = TRUE
		}
		if !IS_NPC(i) && i.Player_specials.Carried_by != nil && i.Player_specials.Carried_by != ch {
			send_to_char(ch, libc.CString(", is being carried by %s"), func() *byte {
				if readIntro(ch, i.Player_specials.Carried_by) == 1 {
					return get_i_name(ch, i.Player_specials.Carried_by)
				}
				return introd_calc(i.Player_specials.Carried_by)
			}())
			count = TRUE
		}
		if !IS_NPC(i) && i.Grappling != nil && i.Grappling == ch {
			send_to_char(ch, libc.CString(", is grappling with YOU"))
			count = TRUE
		}
		if !IS_NPC(i) && i.Absorbing != nil && i.Absorbing == ch {
			send_to_char(ch, libc.CString(", is absorbing from YOU"))
			count = TRUE
		}
		if !IS_NPC(i) && ch.Absorbing != nil && ch.Absorbing == i {
			send_to_char(ch, libc.CString(", is being absorbed from by YOU"))
			count = TRUE
		}
		if !IS_NPC(i) && ch.Grappling != nil && ch.Grappling == i {
			send_to_char(ch, libc.CString(", is being grappled with by YOU"))
			count = TRUE
		}
		if !IS_NPC(i) && ch.Player_specials.Carrying != nil && ch.Player_specials.Carrying == i {
			send_to_char(ch, libc.CString(", is being carried by you"))
			count = TRUE
		}
		if !IS_NPC(ch) && !IS_NPC(i) && i.Fighting != nil {
			if !PLR_FLAGGED(i, PLR_SPAR) || PLR_FLAGGED(i, PLR_SPAR) && (!PLR_FLAGGED(i.Fighting, PLR_SPAR) || IS_NPC(i.Fighting)) {
				send_to_char(ch, libc.CString(", is here fighting "))
			}
			if PLR_FLAGGED(i, PLR_SPAR) && PLR_FLAGGED(i.Fighting, PLR_SPAR) {
				send_to_char(ch, libc.CString(", is here sparring "))
			}
			if i.Fighting == ch {
				send_to_char(ch, libc.CString("@rYOU@w"))
				count = TRUE
			} else {
				if i.In_room == i.Fighting.In_room {
					send_to_char(ch, libc.CString("%s"), func() *byte {
						if ch.Admlevel != 0 {
							return GET_NAME(i.Fighting)
						}
						if readIntro(ch, i.Fighting) == 1 {
							return get_i_name(ch, i.Fighting)
						}
						return JUGGLERACELOWER(i.Fighting)
					}())
					count = TRUE
				} else {
					send_to_char(ch, libc.CString("someone who has already left!"))
				}
			}
		}
	}
	if i.Sits != nil {
		chair = i.Sits
		if PLR_FLAGGED(i, PLR_HEALT) {
			send_to_char(ch, libc.CString("@w is floating inside a healing tank."))
		} else if count == TRUE {
			send_to_char(ch, libc.CString(",@w and%s on %s."), positions[int(i.Position)], chair.Short_description)
		} else if count == FALSE {
			send_to_char(ch, libc.CString("@w%s on %s."), positions[int(i.Position)], chair.Short_description)
		}
	} else if !PLR_FLAGGED(i, PLR_PILOTING) && i.Sits == nil && (!IS_NPC(i) || i.Fighting == nil) {
		if count == TRUE {
			send_to_char(ch, libc.CString("@w, and%s."), positions[int(i.Position)])
		}
		if count == FALSE {
			send_to_char(ch, libc.CString("@w%s."), positions[int(i.Position)])
		}
	} else if PLR_FLAGGED(i, PLR_PILOTING) {
		send_to_char(ch, libc.CString("@w, is sitting in the pilot's chair.\r\n"))
	} else {
		if i.Fighting != nil && !IS_NPC(ch) && !IS_NPC(i) {
			if !PLR_FLAGGED(i, PLR_SPAR) {
				send_to_char(ch, libc.CString(", is here fighting "))
			}
			if PLR_FLAGGED(i, PLR_SPAR) {
				send_to_char(ch, libc.CString(", is here sparring "))
			}
			if i.Fighting == ch {
				send_to_char(ch, libc.CString("@rYOU@w!"))
			} else {
				if i.In_room == i.Fighting.In_room {
					send_to_char(ch, libc.CString("%s!"), func() *byte {
						if ch.Admlevel != 0 {
							return GET_NAME(i.Fighting)
						}
						if readIntro(ch, i.Fighting) == 1 {
							return get_i_name(ch, i.Fighting)
						}
						return JUGGLERACELOWER(i.Fighting)
					}())
				} else {
					send_to_char(ch, libc.CString("someone who has already left!"))
				}
			}
		} else if !IS_NPC(i) {
			send_to_char(ch, libc.CString(" is here struggling with thin air."))
		}
	}
	if AFF_FLAGGED(ch, AFF_DETECT_ALIGN) {
		if IS_EVIL(i) {
			send_to_char(ch, libc.CString(" (@rRed@[3] Aura)"))
		} else if IS_GOOD(i) {
			send_to_char(ch, libc.CString(" (@bBlue@[3] Aura)"))
		}
	}
	if !IS_NPC(i) && PRF_FLAGGED(i, PRF_AFK) {
		send_to_char(ch, libc.CString(" @D(@RAFK@D)"))
	} else if !IS_NPC(i) && i.Timer > 3 {
		send_to_char(ch, libc.CString(" @D(@RIDLE@D)"))
	}
	send_to_char(ch, libc.CString("@n\r\n"))
	if i.Listenroom > 0 {
		var eaves [300]byte
		stdio.Sprintf(&eaves[0], "@w...$e is spying on everything to the @c%s@w.", dirs[i.Eavesdir])
		act(&eaves[0], TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if !IS_NPC(i) {
		if PLR_FLAGGED(i, PLR_FISHING) {
			act(libc.CString("@w...$e is @Cfishing@w.@n"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
	}
	if PLR_FLAGGED(i, PLR_AURALIGHT) {
		var bloom [2048]byte
		stdio.Sprintf(&bloom[0], "...is surrounded by a bright %s aura.@n", aura_types[i.Aura])
		act(&bloom[0], TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_SANCTUARY) && GET_SKILL(i, SKILL_AQUA_BARRIER) == 0 {
		act(libc.CString("@w...$e has a @bbarrier@w around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_FIRESHIELD) {
		act(libc.CString("@w...$e has @rf@Rl@Ya@rm@Re@Ys@w around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_HEALGLOW) {
		act(libc.CString("@w...$e has a serene @Cblue@Y glow@w around $s body."), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_EARMOR) {
		act(libc.CString("@w...$e has ghostly @Ggreen@w ethereal armor around $s body."), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_SANCTUARY) && GET_SKILL(i, SKILL_AQUA_BARRIER) != 0 {
		act(libc.CString("@w...$e has a @bbarrier@w of @cwater@w and @CKi@w around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_FLYING) && i.Altitude == 1 {
		act(libc.CString("@w...$e is in the air!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_FLYING) && i.Altitude == 2 {
		act(libc.CString("@w...$e is high in the air!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if i.Kaioken > 0 {
		act(libc.CString("@w...@r$e has a red aura around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if !IS_NPC(i) && PLR_FLAGGED(i, PLR_SPIRAL) {
		act(libc.CString("@w...$e is spinning in a vortex!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if IS_TRANSFORMED(i) && int(i.Race) != RACE_ANDROID && int(i.Race) != RACE_SAIYAN && int(i.Race) != RACE_HALFBREED {
		act(libc.CString("@w...$e has energy crackling around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if i.Charge != 0 && int(i.Race) != RACE_SAIYAN && int(i.Race) != RACE_HALFBREED {
		var aura [2048]byte
		stdio.Sprintf(&aura[0], "@w...$e has a @Ybright@w %s aura around $s body!", aura_types[i.Aura])
		act(&aura[0], TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if !PLR_FLAGGED(i, PLR_OOZARU) && i.Charge != 0 && IS_TRANSFORMED(i) && (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) {
		act(libc.CString("@w...$e has a @Ybright @Yg@yo@Yl@yd@Ye@yn@w aura around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if !PLR_FLAGGED(i, PLR_OOZARU) && i.Charge != 0 && !IS_TRANSFORMED(i) && (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) {
		var aura [2048]byte
		stdio.Sprintf(&aura[0], "@w...$e has a @Ybright@w %s aura around $s body!", aura_types[i.Aura])
		act(&aura[0], TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if !PLR_FLAGGED(i, PLR_OOZARU) && i.Charge == 0 && IS_TRANSFORMED(i) && (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) {
		act(libc.CString("@w...$e has energy crackling around $s body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if PLR_FLAGGED(i, PLR_OOZARU) && i.Charge != 0 && (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) {
		act(libc.CString("@w...$e is in the form of a @rgreat ape@w!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if (i.Genome[0]) == 11 {
		act(libc.CString("@w...$e has expanded $s body size@w!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if AFF_FLAGGED(i, AFF_HAYASA) {
		act(libc.CString("@w...$e has a soft @cblue@w glow around $s body!"), FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if PLR_FLAGGED(i, PLR_OOZARU) && i.Charge == 0 && (int(i.Race) == RACE_SAIYAN || int(i.Race) == RACE_HALFBREED) {
		act(libc.CString("@w...$e has energy crackling around $s @rgreat ape@w body!"), TRUE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if i.Feature != nil {
		var woo [64936]byte
		stdio.Sprintf(&woo[0], "@C%s@n", i.Feature)
		act(&woo[0], FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
	}
	if i.Rdisplay != nil {
		if i.Rdisplay != libc.CString("Empty") {
			var rdis [64936]byte
			stdio.Sprintf(&rdis[0], "...%s", i.Rdisplay)
			act(&rdis[0], FALSE, i, nil, unsafe.Pointer(ch), TO_VICT)
		}
	}
}

type hide_node struct {
	Next   *hide_node
	Hidden *char_data
}

func list_char_to_char(list *char_data, ch *char_data) {
	var (
		i        *char_data
		j        *char_data
		hideinfo *hide_node
		lasthide *hide_node
		tmphide  *hide_node
		num      int
	)
	hideinfo = func() *hide_node {
		lasthide = nil
		return lasthide
	}()
	for i = list; i != nil; i = i.Next_in_room {
		if AFF_FLAGGED(i, AFF_HIDE) && roll_resisted(i, SKILL_HIDE, ch, SKILL_SPOT) != 0 {
			if GET_SKILL(i, SKILL_HIDE) != 0 && !IS_NPC(ch) && i != ch {
				improve_skill(i, SKILL_HIDE, 1)
			}
			tmphide = new(hide_node)
			tmphide.Next = nil
			tmphide.Hidden = i
			if lasthide == nil {
				hideinfo = func() *hide_node {
					lasthide = tmphide
					return lasthide
				}()
			} else {
				lasthide.Next = tmphide
				lasthide = tmphide
			}
			continue
		}
	}
	for i = list; i != nil; i = i.Next_in_room {
		if ch == i || !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_HOLYLIGHT) && IS_NPC(i) && i.Long_descr != nil && *i.Long_descr == '.' {
			continue
		}
		for tmphide = hideinfo; tmphide != nil; tmphide = tmphide.Next {
			if tmphide.Hidden == i {
				break
			}
		}
		if tmphide != nil {
			continue
		}
		if CAN_SEE(ch, i) {
			num = 0
			if config_info.Play.Stack_mobs != 0 {
				for j = list; j != i; j = j.Next_in_room {
					if i.Nr == j.Nr && int(i.Position) == int(j.Position) && i.Affected_by[0] == j.Affected_by[0] && i.Affected_by[1] == j.Affected_by[1] && i.Affected_by[2] == j.Affected_by[2] && i.Affected_by[3] == j.Affected_by[3] && (i.Fighting == nil && j.Fighting == nil) && (i.Hit == gear_pl(i) && j.Hit == gear_pl(j)) && libc.StrCmp(GET_NAME(i), GET_NAME(j)) == 0 {
						for tmphide = hideinfo; tmphide != nil; tmphide = tmphide.Next {
							if tmphide.Hidden == j {
								break
							}
						}
						if tmphide == nil {
							break
						}
					}
				}
				if j != i {
					continue
				}
				for j = i; j != nil; j = j.Next_in_room {
					if i.Nr == j.Nr && int(i.Position) == int(j.Position) && i.Affected_by[0] == j.Affected_by[0] && i.Affected_by[1] == j.Affected_by[1] && i.Affected_by[2] == j.Affected_by[2] && i.Affected_by[3] == j.Affected_by[3] && (i.Fighting == nil && j.Fighting == nil) && (i.Hit == i.Max_hit && j.Hit == j.Max_hit) && libc.StrCmp(GET_NAME(i), GET_NAME(j)) == 0 {
						for tmphide = hideinfo; tmphide != nil; tmphide = tmphide.Next {
							if tmphide.Hidden == j {
								break
							}
						}
						if tmphide == nil {
							num++
						}
					}
				}
			}
			send_to_char(ch, libc.CString("@w"))
			if num > 1 {
				send_to_char(ch, libc.CString("@D(@Rx@Y%2i@D)@n "), num)
			}
			list_one_char(i, ch)
			send_to_char(ch, libc.CString("@n"))
		} else if room_is_dark(ch.In_room) != 0 && !CAN_SEE_IN_DARK(ch) && AFF_FLAGGED(i, AFF_INFRAVISION) {
			send_to_char(ch, libc.CString("@wYou see a pair of glowing red eyes looking your way.@n\r\n"))
		}
	}
}
func do_auto_exits(target_room room_rnum, ch *char_data, exit_mode int) {
	var (
		door       int
		door_found int = 0
		has_light  int = FALSE
		i          int
		dlist1     [500]byte
		dlist2     [500]byte
		dlist3     [500]byte
		dlist4     [500]byte
		dlist5     [500]byte
		dlist6     [500]byte
		dlist7     [500]byte
		dlist8     [500]byte
		dlist9     [500]byte
		dlist10    [500]byte
		dlist11    [500]byte
		dlist12    [500]byte
	)
	dlist1[0] = '\x00'
	dlist2[0] = '\x00'
	dlist3[0] = '\x00'
	dlist4[0] = '\x00'
	dlist5[0] = '\x00'
	dlist6[0] = '\x00'
	dlist7[0] = '\x00'
	dlist8[0] = '\x00'
	dlist9[0] = '\x00'
	dlist10[0] = '\x00'
	dlist11[0] = '\x00'
	dlist12[0] = '\x00'
	if exit_mode == EXIT_OFF {
		send_to_char(ch, libc.CString("@D------------------------------------------------------------------------@n\r\n"))
	}
	var space int = FALSE
	if (func() int {
		if target_room != room_rnum(-1) && target_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_SPACE && (func() room_vnum {
		if target_room != room_rnum(-1) && target_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
		}
		return -1
	}()) >= 20000 {
		space = TRUE
	}
	if exit_mode == EXIT_NORMAL && space == FALSE && ch.In_room == target_room {
		send_to_char(ch, libc.CString("@D------------------------------------------------------------------------@n\r\n"))
		send_to_char(ch, libc.CString("@w      Compass           Auto-Map            Map Key\r\n"))
		send_to_char(ch, libc.CString("@R     ---------         ----------   -----------------------------\r\n"))
		gen_map(ch, 0)
		send_to_char(ch, libc.CString("@D------------------------------------------------------------------------@n\r\n"))
	}
	if exit_mode == EXIT_NORMAL && space == TRUE {
		send_to_char(ch, libc.CString("@D------------------------------[@CRadar@D]---------------------------------@n\r\n"))
		printmap(int(target_room), ch, 1, -1)
		send_to_char(ch, libc.CString("     @D[@wTurn autoexit complete on for directions instead of radar@D]@n\r\n"))
		send_to_char(ch, libc.CString("@D------------------------------------------------------------------------@n\r\n"))
	}
	if exit_mode == EXIT_COMPLETE || exit_mode == EXIT_NORMAL && space == FALSE && ch.In_room != target_room {
		send_to_char(ch, libc.CString("@D----------------------------[@gObvious Exits@D]-----------------------------@n\r\n"))
		if AFF_FLAGGED(ch, AFF_BLIND) {
			send_to_char(ch, libc.CString("You can't see a damned thing, you're blind!\r\n"))
			return
		}
		if PLR_FLAGGED(ch, PLR_EYEC) {
			send_to_char(ch, libc.CString("You can't see a damned thing, your eyes are closed!\r\n"))
			return
		}
		if room_is_dark(ch.In_room) != 0 && !CAN_SEE_IN_DARK(ch) && !PLR_FLAGGED(ch, PLR_AURALIGHT) {
			send_to_char(ch, libc.CString("It is pitch black...\r\n"))
			return
		}
		for i = 0; i < NUM_WEARS; i++ {
			if (ch.Equipment[i]) != nil {
				if int((ch.Equipment[i]).Type_flag) == ITEM_LIGHT {
					if ((ch.Equipment[i]).Value[VAL_LIGHT_HOURS]) != 0 {
						has_light = TRUE
					}
				}
			}
		}
		if PLR_FLAGGED(ch, PLR_AURALIGHT) {
			has_light = TRUE
		}
		for door = 0; door < NUM_OF_DIRS; door++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]) != nil && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room != room_rnum(-1) {
				if ADM_FLAGGED(ch, ADM_SEESECRET) || ch.Admlevel > 4 {
					door_found++
					var blam [9]byte
					stdio.Sprintf(&blam[0], "%s", dirs[door])
					blam[0] = byte(int8(unicode.ToUpper(rune(blam[0]))))
					if door == 6 {
						stdio.Sprintf(&dlist1[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							var argh [100]byte
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION NORTHWEST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction NORTHWEST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist1[libc.StrLen(&dlist1[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 0 {
						stdio.Sprintf(&dlist2[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION NORTH@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction NORTH at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [200]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist2[libc.StrLen(&dlist2[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 7 {
						stdio.Sprintf(&dlist3[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION NORTHEAST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction NORTHEAST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist3[libc.StrLen(&dlist3[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 1 {
						stdio.Sprintf(&dlist4[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION EAST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction EAST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist4[libc.StrLen(&dlist4[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 8 {
						stdio.Sprintf(&dlist5[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION SOUTHEAST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction SOUTHEAST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist5[libc.StrLen(&dlist5[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 2 {
						stdio.Sprintf(&dlist6[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION SOUTH@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction SOUTH at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist6[libc.StrLen(&dlist6[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 9 {
						stdio.Sprintf(&dlist7[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION SOUTHWEST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction SOUTHWEST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist7[libc.StrLen(&dlist7[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 3 {
						stdio.Sprintf(&dlist8[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION WEST@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction WEST at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist8[libc.StrLen(&dlist8[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 4 {
						stdio.Sprintf(&dlist9[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION UP@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction UP at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist9[libc.StrLen(&dlist9[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 5 {
						stdio.Sprintf(&dlist10[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION DOWN@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction DOWN at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist10[libc.StrLen(&dlist10[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 10 {
						stdio.Sprintf(&dlist11[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION INSIDE@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction INSIDE at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist11[libc.StrLen(&dlist11[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
					if door == 11 {
						stdio.Sprintf(&dlist12[0], "@c%-9s @D- [@Y%5d@D]@w %s.\r\n", &blam[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
						if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<0)) != 0 || (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) != 0 {
							if fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword) == nil {
								send_to_char(ch, libc.CString("@RREPORT THIS ERROR IMMEADIATLY FOR DIRECTION OUTSIDE@n\r\n"))
								basic_mud_log(libc.CString("ERROR: %s found error direction OUTSIDE at room %d"), GET_NAME(ch), func() room_vnum {
									if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
										return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
									}
									return -1
								}())
								return
							}
							var argh [100]byte
							stdio.Sprintf(&argh[0], "%s ", func() *byte {
								if libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
							stdio.Sprintf(&dlist12[libc.StrLen(&dlist12[0])], "                    The %s%s %s %s %s%s.\r\n", func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 4)) != 0 {
									return "@rsecret@w "
								}
								return ""
							}(), func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil && libc.StrCaseCmp(fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword), libc.CString("undefined")) != 0 {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}(), func() string {
								if libc.StrStr(&argh[0], libc.CString("s ")) != nil {
									return "are"
								}
								return "is"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) != 0 {
									return "closed"
								}
								return "open"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 2)) != 0 {
									return "and locked"
								}
								return "and unlocked"
							}(), func() string {
								if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 3)) != 0 {
									return " (pickproof)"
								}
								return ""
							}())
						}
					}
				} else {
					if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info & (1 << 1)) == 0 {
						door_found++
						var blam [9]byte
						stdio.Sprintf(&blam[0], "%s", dirs[door])
						blam[0] = byte(int8(unicode.ToUpper(rune(blam[0]))))
						if door == 6 {
							stdio.Sprintf(&dlist1[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 0 {
							stdio.Sprintf(&dlist2[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 7 {
							stdio.Sprintf(&dlist3[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 1 {
							stdio.Sprintf(&dlist4[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 8 {
							stdio.Sprintf(&dlist5[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 2 {
							stdio.Sprintf(&dlist6[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 9 {
							stdio.Sprintf(&dlist7[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 3 {
							stdio.Sprintf(&dlist8[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 4 {
							stdio.Sprintf(&dlist9[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 5 {
							stdio.Sprintf(&dlist10[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 10 {
							stdio.Sprintf(&dlist11[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
						if door == 11 {
							stdio.Sprintf(&dlist12[0], "@c%-9s @D-@w %s\r\n", &blam[0], func() string {
								if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room) != 0 && !CAN_SEE_IN_DARK(ch) && has_light == 0 {
									return "@bToo dark to tell.@w"
								}
								return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room)))).Name)
							}())
						}
					} else if config_info.Play.Disp_closed_doors != 0 && (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Exit_info&(1<<4)) == 0 {
						door_found++
						var blam [9]byte
						stdio.Sprintf(&blam[0], "%s", dirs[door])
						blam[0] = byte(int8(unicode.ToUpper(rune(blam[0]))))
						if door == 6 {
							stdio.Sprintf(&dlist1[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 0 {
							stdio.Sprintf(&dlist2[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 7 {
							stdio.Sprintf(&dlist3[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 1 {
							stdio.Sprintf(&dlist4[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 8 {
							stdio.Sprintf(&dlist5[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 2 {
							stdio.Sprintf(&dlist6[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 9 {
							stdio.Sprintf(&dlist7[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 3 {
							stdio.Sprintf(&dlist8[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 4 {
							stdio.Sprintf(&dlist9[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 5 {
							stdio.Sprintf(&dlist10[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 10 {
							stdio.Sprintf(&dlist11[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
						if door == 11 {
							stdio.Sprintf(&dlist12[0], "@c%-9s @D-@w The %s appears @rclosed.@n\r\n", &blam[0], func() *byte {
								if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword != nil {
									return fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).Keyword)
								}
								return libc.CString("opening")
							}())
						}
					}
				}
			}
		}
		if door_found == 0 {
			send_to_char(ch, libc.CString(" None.\r\n"))
		}
		if libc.StrStr(&dlist1[0], libc.CString("Northwest")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist1[0])
			dlist1[0] = '\x00'
		}
		if libc.StrStr(&dlist2[0], libc.CString("North")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist2[0])
			dlist2[0] = '\x00'
		}
		if libc.StrStr(&dlist3[0], libc.CString("Northeast")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist3[0])
			dlist3[0] = '\x00'
		}
		if libc.StrStr(&dlist4[0], libc.CString("East")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist4[0])
			dlist4[0] = '\x00'
		}
		if libc.StrStr(&dlist5[0], libc.CString("Southeast")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist5[0])
			dlist5[0] = '\x00'
		}
		if libc.StrStr(&dlist6[0], libc.CString("South")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist6[0])
			dlist6[0] = '\x00'
		}
		if libc.StrStr(&dlist7[0], libc.CString("Southwest")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist7[0])
			dlist7[0] = '\x00'
		}
		if libc.StrStr(&dlist8[0], libc.CString("West")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist8[0])
			dlist8[0] = '\x00'
		}
		if libc.StrStr(&dlist9[0], libc.CString("Up")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist9[0])
			dlist9[0] = '\x00'
		}
		if libc.StrStr(&dlist10[0], libc.CString("Down")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist10[0])
			dlist10[0] = '\x00'
		}
		if libc.StrStr(&dlist11[0], libc.CString("Inside")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist11[0])
			dlist11[0] = '\x00'
		}
		if libc.StrStr(&dlist12[0], libc.CString("Outside")) != nil {
			send_to_char(ch, libc.CString("%s"), &dlist12[0])
			dlist12[0] = '\x00'
		}
		send_to_char(ch, libc.CString("@D------------------------------------------------------------------------@n\r\n"))
		if ROOM_FLAGGED(target_room, ROOM_HOUSE) && !ROOM_FLAGGED(target_room, ROOM_GARDEN1) && !ROOM_FLAGGED(target_room, ROOM_GARDEN2) {
			send_to_char(ch, libc.CString("@D[@GItems Stored@D: @g%d@D]@n\r\n"), check_saveroom_count(ch, nil))
		}
		if ROOM_FLAGGED(target_room, ROOM_HOUSE) && ROOM_FLAGGED(target_room, ROOM_GARDEN1) && !ROOM_FLAGGED(target_room, ROOM_GARDEN2) {
			send_to_char(ch, libc.CString("@D[@GPlants Planted@D: @g%d@W, @GMAX@D: @R8@D]@n\r\n"), check_saveroom_count(ch, nil))
		}
		if ROOM_FLAGGED(target_room, ROOM_HOUSE) && !ROOM_FLAGGED(target_room, ROOM_GARDEN1) && ROOM_FLAGGED(target_room, ROOM_GARDEN2) {
			send_to_char(ch, libc.CString("@D[@GPlants Planted@D: @g%d@W, @GMAX@D: @R20@D]@n\r\n"), check_saveroom_count(ch, nil))
		}
		if ch.Radar1 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar2 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar3 != (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CTwo of your buoys are floating here.@n\r\n"))
		} else if ch.Radar1 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar2 != (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar3 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CTwo of your buoys are floating here.@n\r\n"))
		} else if ch.Radar1 != (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar2 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar3 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CTwo of your buoys are floating here.@n\r\n"))
		} else if ch.Radar1 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar2 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && ch.Radar3 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) && target_room != 0 {
			send_to_char(ch, libc.CString("@CAll three of your buoys are floating here. Why?@n\r\n"))
		} else if ch.Radar1 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CYour @cBuoy #1@C is floating here.@n\r\n"))
		} else if ch.Radar2 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CYour @cBuoy #2@C is floating here.@n\r\n"))
		} else if ch.Radar3 == (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) {
			send_to_char(ch, libc.CString("@CYour @cBuoy #3@C is floating here.@n\r\n"))
		}
	}
}
func do_auto_exits2(target_room room_rnum, ch *char_data) {
	var (
		door int
		slen int = 0
	)
	send_to_char(ch, libc.CString("\nExits: "))
	for door = 0; door < NUM_OF_DIRS; door++ {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]) == nil || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door]).To_room == room_rnum(-1) {
			continue
		}
		if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dir_option[door], 1<<1) {
			continue
		}
		send_to_char(ch, libc.CString("%s "), abbr_dirs[door])
		slen++
	}
	send_to_char(ch, libc.CString("%s\r\n"), func() string {
		if slen != 0 {
			return ""
		}
		return "None!"
	}())
}
func do_exits(ch *char_data, argument *byte, cmd int, subcmd int) {
	if !PRF_FLAGGED(ch, PRF_NODEC) {
		do_auto_exits(ch.In_room, ch, EXIT_COMPLETE)
	} else {
		do_auto_exits2(ch.In_room, ch)
	}
}

var exitlevels [5]*byte = [5]*byte{libc.CString("off"), libc.CString("normal"), libc.CString("n/a"), libc.CString("complete"), libc.CString("\n")}

func do_autoexit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		tp  int
	)
	if IS_NPC(ch) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Your current autoexit level is %s.\r\n"), exitlevels[func() int {
			if !IS_NPC(ch) {
				return (func() int {
					if PRF_FLAGGED(ch, PRF_AUTOEXIT) {
						return 1
					}
					return 0
				}()) + (func() int {
					if PRF_FLAGGED(ch, PRF_FULL_EXIT) {
						return 2
					}
					return 0
				}())
			}
			return 0
		}()])
		return
	}
	if (func() int {
		tp = search_block(&arg[0], &exitlevels[0], FALSE)
		return tp
	}()) == -1 {
		send_to_char(ch, libc.CString("Usage: Autoexit { Off | Normal | Complete }\r\n"))
		return
	}
	switch tp {
	case EXIT_OFF:
		ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)] &= bitvector_t(int32(^(1 << (int(PRF_AUTOEXIT % 32)))))
		ch.Player_specials.Pref[int(PRF_FULL_EXIT/32)] &= bitvector_t(int32(^(1 << (int(PRF_FULL_EXIT % 32)))))
	case EXIT_NORMAL:
		ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)] |= bitvector_t(int32(1 << (int(PRF_AUTOEXIT % 32))))
		ch.Player_specials.Pref[int(PRF_FULL_EXIT/32)] &= bitvector_t(int32(^(1 << (int(PRF_FULL_EXIT % 32)))))
	case EXIT_COMPLETE:
		ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)] |= bitvector_t(int32(1 << (int(PRF_AUTOEXIT % 32))))
		ch.Player_specials.Pref[int(PRF_FULL_EXIT/32)] |= bitvector_t(int32(1 << (int(PRF_FULL_EXIT % 32))))
	}
	send_to_char(ch, libc.CString("Your @rautoexit level@n is now %s.\r\n"), exitlevels[func() int {
		if !IS_NPC(ch) {
			return (func() int {
				if PRF_FLAGGED(ch, PRF_AUTOEXIT) {
					return 1
				}
				return 0
			}()) + (func() int {
				if PRF_FLAGGED(ch, PRF_FULL_EXIT) {
					return 2
				}
				return 0
			}())
		}
		return 0
	}()])
}
func look_at_room(target_room room_rnum, ch *char_data, ignore_brief int) {
	var (
		rm *room_data = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))
		t  *trig_data
	)
	if ch.Desc == nil {
		return
	}
	if room_is_dark(target_room) != 0 && !CAN_SEE_IN_DARK(ch) && !PLR_FLAGGED(ch, PLR_AURALIGHT) {
		send_to_char(ch, libc.CString("It is pitch black...\r\n"))
		return
	} else if AFF_FLAGGED(ch, AFF_BLIND) {
		send_to_char(ch, libc.CString("You see nothing but infinite darkness...\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_EYEC) {
		send_to_char(ch, libc.CString("You can't see a damned thing, your eyes are closed!\r\n"))
		return
	}
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_ROOMFLAGS) {
		var (
			buf  [64936]byte
			buf2 [64936]byte
			buf3 [64936]byte
		)
		sprintbitarray((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Room_flags[:], room_bits[:], RF_ARRAY_MAX, &buf[0])
		sprinttype(rm.Sector_type, sector_types[:], &buf2[0], uint64(64936))
		if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NODEC) {
			send_to_char(ch, libc.CString("\r\n@wO----------------------------------------------------------------------O@n\r\n"))
		}
		send_to_char(ch, libc.CString("@wLocation: @G%-70s@w\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Name)
		if rm.Script != nil {
			send_to_char(ch, libc.CString("@D[@GTriggers"))
			for t = rm.Script.Trig_list; t != nil; t = t.Next {
				send_to_char(ch, libc.CString(" %d"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(t.Nr)))).Vnum)
			}
			send_to_char(ch, libc.CString("@D] "))
		}
		stdio.Sprintf(&buf3[0], "@D[ @G%s@D] @wSector: @D[ @G%s @D] @wVnum: @D[@G%5d@D]@n", &buf[0], &buf2[0], func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}())
		send_to_char(ch, libc.CString("@wFlags: %-70s@w\r\n"), &buf3[0])
		if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NODEC) {
			send_to_char(ch, libc.CString("@wO----------------------------------------------------------------------O@n\r\n"))
		}
	} else {
		if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NODEC) {
			send_to_char(ch, libc.CString("@wO----------------------------------------------------------------------O@n\r\n"))
		}
		send_to_char(ch, libc.CString("@wLocation: %-70s@n\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Name)
		if ROOM_FLAGGED(target_room, ROOM_EARTH) {
			send_to_char(ch, libc.CString("@wPlanet: @GEarth@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_CERRIA) {
			send_to_char(ch, libc.CString("@wPlanet: @RCerria@n\r\n"))
		} else if (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) >= 3400 && (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) <= 3599 || (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) >= 62900 && (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) <= 0xF617 || (func() room_vnum {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Number
			}
			return -1
		}()) == 19600 {
			send_to_char(ch, libc.CString("@wPlanet: @BZenith@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_AETHER) {
			send_to_char(ch, libc.CString("@wPlanet: @MAether@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_FRIGID) {
			send_to_char(ch, libc.CString("@wPlanet: @CFrigid@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_SPACE) {
			send_to_char(ch, libc.CString("@wPlanet: @DNone@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_VEGETA) {
			send_to_char(ch, libc.CString("@wPlanet: @YVegeta@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_NAMEK) {
			send_to_char(ch, libc.CString("@wPlanet: @gNamek@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_KONACK) {
			send_to_char(ch, libc.CString("@wPlanet: @MKonack@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_NEO) {
			send_to_char(ch, libc.CString("@wPlanet: @WNeo Nirvana@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_AL) {
			send_to_char(ch, libc.CString("@wDimension: @yA@Yf@yt@Ye@yr@Yl@yi@Yf@ye@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_HELL) {
			send_to_char(ch, libc.CString("@wDimension: @RPunishment Hell@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_RHELL) {
			send_to_char(ch, libc.CString("@wDimension: @RH@re@Dl@Rl@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_YARDRAT) {
			send_to_char(ch, libc.CString("@wPlanet: @mYardrat@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_KANASSA) {
			send_to_char(ch, libc.CString("@wPlanet: @BKanassa@n\r\n"))
		} else if ROOM_FLAGGED(target_room, ROOM_ARLIA) {
			send_to_char(ch, libc.CString("@wPlanet: @GArlia@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("@wPlanet: @WUNKNOWN@n\r\n"))
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity <= 0 {
			send_to_char(ch, libc.CString("@wGravity: @WNormal@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 10 {
			send_to_char(ch, libc.CString("@wGravity: @W10x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 20 {
			send_to_char(ch, libc.CString("@wGravity: @W20x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 30 {
			send_to_char(ch, libc.CString("@wGravity: @W30x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 40 {
			send_to_char(ch, libc.CString("@wGravity: @W40x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 50 {
			send_to_char(ch, libc.CString("@wGravity: @W50x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 100 {
			send_to_char(ch, libc.CString("@wGravity: @W100x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 200 {
			send_to_char(ch, libc.CString("@wGravity: @W200x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 300 {
			send_to_char(ch, libc.CString("@wGravity: @W300x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 400 {
			send_to_char(ch, libc.CString("@wGravity: @W400x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 500 {
			send_to_char(ch, libc.CString("@wGravity: @W500x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 1000 {
			send_to_char(ch, libc.CString("@wGravity: @W1,000x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 5000 {
			send_to_char(ch, libc.CString("@wGravity: @W5,000x@n\r\n"))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Gravity == 10000 {
			send_to_char(ch, libc.CString("@wGravity: @W10,000x@n\r\n"))
		}
		if ROOM_FLAGGED(target_room, ROOM_REGEN) {
			send_to_char(ch, libc.CString("@CA feeling of calm and relaxation fills this room.@n\r\n"))
		}
		if ROOM_FLAGGED(target_room, ROOM_AURA) {
			send_to_char(ch, libc.CString("@GAn aura of @gregeneration@G surrounds this area.@n\r\n"))
		}
		if ROOM_FLAGGED(target_room, ROOM_HBTC) {
			send_to_char(ch, libc.CString("@rThis room feels like it opperates in a different time frame.@n\r\n"))
		}
		if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NODEC) {
			send_to_char(ch, libc.CString("@wO----------------------------------------------------------------------O@n\r\n"))
		}
	}
	if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_BRIEF) || ROOM_FLAGGED(target_room, ROOM_DEATH) {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 99 {
			send_to_char(ch, libc.CString("@w%s@n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Description)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg == 100 && ((func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_WATER_SWIM || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Geffect < 0 || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_UNDERWATER) || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FLYING || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_SHOP || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_IMPORTANT) {
			send_to_char(ch, libc.CString("@w%s@n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Description)
		}
		if (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_INSIDE && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg > 0 {
			send_to_char(ch, libc.CString("\r\n"))
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 2 {
				send_to_char(ch, libc.CString("@wA small hole with chunks of debris that can be seen scarring the floor.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 4 {
				send_to_char(ch, libc.CString("@wA couple small holes with chunks of debris that can be seen scarring the floor.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 6 {
				send_to_char(ch, libc.CString("@wA few small holes with chunks of debris that can be seen scarring the floor.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 10 {
				send_to_char(ch, libc.CString("@wThere are several small holes with chunks of debris that can be seen scarring the floor.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 20 {
				send_to_char(ch, libc.CString("@wMany holes fill the floor of this area, many of which have burn marks.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 30 {
				send_to_char(ch, libc.CString("@wThe floor is severely damaged with many large holes.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 50 {
				send_to_char(ch, libc.CString("@wBattle damage covers the entire area. Displayed as a tribute to the battles that have\r\nbeen waged here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 75 {
				send_to_char(ch, libc.CString("@wThis entire area is falling apart, it has been damaged so badly.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 99 {
				send_to_char(ch, libc.CString("@wThis area can not withstand much more damage. Everything has been damaged so badly it\r\nis hard to recognise any particular details about their former quality.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg >= 100 {
				send_to_char(ch, libc.CString("@wThis area is completely destroyed. Nothing is recognisable. Chunks of debris\r\nlitter the ground, filling up holes, and overflowing onto what is left of the\r\nfloor. A haze of smoke is wafting through the air, creating a chilling atmosphere..@n"))
			}
			send_to_char(ch, libc.CString("\r\n"))
		} else if ((func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_CITY || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FIELD || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_HILLS || (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_IMPORTANT) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg > 0 {
			send_to_char(ch, libc.CString("\r\n"))
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 2 {
				send_to_char(ch, libc.CString("@wA small hole with chunks of debris that can be seen scarring the ground.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 4 {
				send_to_char(ch, libc.CString("@wA couple small craters with chunks of debris that can be seen scarring the ground.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 6 {
				send_to_char(ch, libc.CString("@wA few small craters with chunks of debris that can be seen scarring the ground.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 10 {
				send_to_char(ch, libc.CString("@wThere are several small craters with chunks of debris that can be seen scarring the ground.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 20 {
				send_to_char(ch, libc.CString("@wMany craters fill the ground of this area, many of which have burn marks.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 30 {
				send_to_char(ch, libc.CString("@wThe ground is severely damaged with many large craters.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 50 {
				send_to_char(ch, libc.CString("@wBattle damage covers the entire area. Displayed as a tribute to the battles that have\r\nbeen waged here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 75 {
				send_to_char(ch, libc.CString("@wThis entire area is falling apart, it has been damaged so badly.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 99 {
				send_to_char(ch, libc.CString("@wThis area can not withstand much more damage. Everything has been damaged so badly it\r\nis hard to recognise any particular details about their former quality.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg >= 100 {
				send_to_char(ch, libc.CString("@wThis area is completely destroyed. Nothing is recognisable. Chunks of debris\r\nlitter the ground, filling up craters, and overflowing onto what is left of the\r\nground. A haze of smoke is wafting through the air, creating a chilling atmosphere..@n"))
			}
			send_to_char(ch, libc.CString("\r\n"))
		} else if (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FOREST && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg > 0 {
			send_to_char(ch, libc.CString("\r\n"))
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 2 {
				send_to_char(ch, libc.CString("@wA small tree sits in a little crater here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 4 {
				send_to_char(ch, libc.CString("@wTrees have been uprooted by craters in the ground.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 6 {
				send_to_char(ch, libc.CString("@wSeveral trees have been reduced to chunks of debris and are\r\nlaying in a few craters here. @n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 10 {
				send_to_char(ch, libc.CString("@wA large patch of trees have been destroyed and are laying in craters here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 20 {
				send_to_char(ch, libc.CString("@wSeveral craters have merged into one large crater in one part of this forest.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 30 {
				send_to_char(ch, libc.CString("@wThe open sky can easily be seen through a hole of trees destroyed\r\nand resting at the bottom of several craters here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 50 {
				send_to_char(ch, libc.CString("@wA good deal of burning tree pieces can be found strewn across the cratered ground here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 75 {
				send_to_char(ch, libc.CString("@wVery few trees are left standing in this area, replaced instead by large craters.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 99 {
				send_to_char(ch, libc.CString("@wSingle solitary trees can be found still standing here or there in the area.\r\nThe rest have been almost completely obliterated in recent conflicts.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg >= 100 {
				send_to_char(ch, libc.CString("@w  One massive crater fills this area. This desolate crater leaves no\r\nevidence of what used to be found in the area. Smoke slowly wafts into\r\nthe sky from the central point of the crater, creating an oppressive\r\natmosphere.@n"))
			}
			send_to_char(ch, libc.CString("\r\n"))
		} else if (func() int {
			if target_room != room_rnum(-1) && target_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_MOUNTAIN && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg > 0 {
			send_to_char(ch, libc.CString("\r\n"))
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 2 {
				send_to_char(ch, libc.CString("@wA small crater has been burned into the side of this mountain.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 4 {
				send_to_char(ch, libc.CString("@wA couple craters have been burned into the side of this mountain.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 6 {
				send_to_char(ch, libc.CString("@wBurned bits of boulders can be seen lying at the bottom of a few nearby craters.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 10 {
				send_to_char(ch, libc.CString("@wSeveral bad craters can be seen in the side of the mountain here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 20 {
				send_to_char(ch, libc.CString("@wLarge boulders have rolled down the mountain side and collected in many nearby craters.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 30 {
				send_to_char(ch, libc.CString("@wMany craters are covering the mountainside here.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 50 {
				send_to_char(ch, libc.CString("@wThe mountain side has partially collapsed, shedding rubble down towards its base.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 75 {
				send_to_char(ch, libc.CString("@wA peak of the mountain has been blown off, leaving behind a smoldering tip.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg <= 99 {
				send_to_char(ch, libc.CString("@wThe mountain side here has completely collapsed, shedding dangerous rubble down to its base.@n"))
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Dmg >= 100 {
				send_to_char(ch, libc.CString("@w  Half the mountain has been blown away, leaving a scarred and jagged\r\nrock in its place. Billowing smoke wafts up from several parts of the\r\nmountain, filling the nearby skies and blotting out the sun.@n"))
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Geffect >= 1 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Geffect <= 5 {
			send_to_char(ch, libc.CString("@rLava@w is pooling in someplaces here...@n\r\n"))
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Geffect >= 6 {
			send_to_char(ch, libc.CString("@RLava@r covers pretty much the entire area!@n\r\n"))
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Geffect < 0 {
			send_to_char(ch, libc.CString("@cThe entire area is flooded with a @Cmystical@c cube of @Bwater!@n\r\n"))
		}
	}
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NODEC) {
		do_auto_exits2(target_room, ch)
	}
	if !IS_NPC(ch) && !PRF_FLAGGED(ch, PRF_NODEC) {
		do_auto_exits(target_room, ch, func() int {
			if !IS_NPC(ch) {
				return (func() int {
					if PRF_FLAGGED(ch, PRF_AUTOEXIT) {
						return 1
					}
					return 0
				}()) + (func() int {
					if PRF_FLAGGED(ch, PRF_FULL_EXIT) {
						return 2
					}
					return 0
				}())
			}
			return 0
		}())
	}
	if ROOM_FLAGGED(target_room, ROOM_GARDEN1) {
		send_to_char(ch, libc.CString("@D[@GPlants Planted@D: @g%d@W, @GMAX@D: @R8@D]@n\r\n"), check_saveroom_count(ch, nil))
	} else if ROOM_FLAGGED(target_room, ROOM_GARDEN2) {
		send_to_char(ch, libc.CString("@D[@GPlants Planted@D: @g%d@W, @GMAX@D: @R20@D]@n\r\n"), check_saveroom_count(ch, nil))
	} else if ROOM_FLAGGED(target_room, ROOM_HOUSE) {
		send_to_char(ch, libc.CString("@D[@GItems Stored@D: @g%d@D]@n\r\n"), check_saveroom_count(ch, nil))
	}
	list_obj_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).Contents, ch, SHOW_OBJ_LONG, FALSE)
	list_char_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(target_room)))).People, ch)
}
func look_in_direction(ch *char_data, dir int) {
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) != nil {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).General_description != nil {
			send_to_char(ch, libc.CString("%s"), ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).General_description)
		} else {
			var (
				obj      *obj_data
				next_obj *obj_data
				founded  int = FALSE
			)
			for obj = ch.Carrying; obj != nil; obj = next_obj {
				next_obj = obj.Next_content
				if GET_OBJ_VNUM(obj) == 17 {
					founded = TRUE
				}
			}
			if founded == FALSE {
				send_to_char(ch, libc.CString("You were unable to discern anything about that direction. Try looking again...\r\n"))
				var obj *obj_data
				obj = read_object(17, VIRTUAL)
				obj_to_char(obj, ch)
			}
		}
		if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<0) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword != nil {
			if !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<4) && EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<1) {
				send_to_char(ch, libc.CString("The %s is closed.\r\n"), fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword))
			} else if !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<1) {
				send_to_char(ch, libc.CString("The %s is open.\r\n"), fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword))
			}
		}
	} else {
		send_to_char(ch, libc.CString("Nothing special there...\r\n"))
	}
}
func look_in_obj(ch *char_data, arg *byte) {
	var (
		obj   *obj_data  = nil
		dummy *char_data = nil
		amt   int
		bits  int
	)
	if *arg == 0 {
		send_to_char(ch, libc.CString("Look in what?\r\n"))
	} else if (func() int {
		bits = generic_find(arg, (1<<2)|1<<3|1<<5, ch, &dummy, &obj)
		return bits
	}()) == 0 {
		send_to_char(ch, libc.CString("There doesn't seem to be %s %s here.\r\n"), AN(arg), arg)
	} else if find_exdesc(arg, obj.Ex_description) != nil && bits == 0 {
		send_to_char(ch, libc.CString("There's nothing inside that!\r\n"))
	} else if int(obj.Type_flag) == ITEM_PORTAL && !OBJVAL_FLAGGED(obj, 1<<0) {
		if (obj.Value[VAL_PORTAL_APPEAR]) < 0 {
			var portal_dest room_rnum = real_room(room_vnum(obj.Value[VAL_PORTAL_DEST]))
			if portal_dest == room_rnum(-1) {
				send_to_char(ch, libc.CString("You see nothing but infinite darkness...\r\n"))
			} else if room_is_dark(portal_dest) != 0 && !CAN_SEE_IN_DARK(ch) && !PLR_FLAGGED(ch, PLR_AURALIGHT) {
				send_to_char(ch, libc.CString("You see nothing but infinite darkness...\r\n"))
			} else {
				send_to_char(ch, libc.CString("After seconds of concentration you see the image of %s.\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(portal_dest)))).Name)
			}
		} else if (obj.Value[VAL_PORTAL_APPEAR]) < MAX_PORTAL_TYPES {
			send_to_char(ch, libc.CString("%s\r\n"), portal_appearance[obj.Value[VAL_PORTAL_APPEAR]])
		} else {
			send_to_char(ch, libc.CString("All you can see is the glow of the portal.\r\n"))
		}
	} else if int(obj.Type_flag) == ITEM_VEHICLE {
		if OBJVAL_FLAGGED(obj, 1<<2) {
			send_to_char(ch, libc.CString("It is closed.\r\n"))
		} else if (obj.Value[VAL_VEHICLE_APPEAR]) < 0 {
			var vehicle_inside room_rnum = real_room(room_vnum(obj.Value[VAL_VEHICLE_ROOM]))
			if vehicle_inside == room_rnum(-1) {
				send_to_char(ch, libc.CString("You cannot see inside that.\r\n"))
			} else if room_is_dark(vehicle_inside) != 0 && !CAN_SEE_IN_DARK(ch) && !PLR_FLAGGED(ch, PLR_AURALIGHT) {
				send_to_char(ch, libc.CString("It is pitch black...\r\n"))
			} else {
				send_to_char(ch, libc.CString("You look inside and see:\r\n"))
				look_at_room(vehicle_inside, ch, 0)
			}
		} else {
			send_to_char(ch, libc.CString("You cannot see inside that.\r\n"))
		}
	} else if int(obj.Type_flag) == ITEM_WINDOW {
		look_out_window(ch, arg)
	} else if int(obj.Type_flag) != ITEM_DRINKCON && int(obj.Type_flag) != ITEM_FOUNTAIN && int(obj.Type_flag) != ITEM_CONTAINER && int(obj.Type_flag) != ITEM_PORTAL {
		send_to_char(ch, libc.CString("There's nothing inside that!\r\n"))
	} else if int(obj.Type_flag) == ITEM_CONTAINER || int(obj.Type_flag) == ITEM_PORTAL {
		if OBJVAL_FLAGGED(obj, 1<<2) {
			send_to_char(ch, libc.CString("It is closed.\r\n"))
		} else {
			send_to_char(ch, libc.CString("%s"), obj.Short_description)
			if int(obj.Type_flag) == ITEM_CONTAINER && (GET_OBJ_VNUM(obj) == 697 || GET_OBJ_VNUM(obj) == 698 || GET_OBJ_VNUM(obj) == 682 || GET_OBJ_VNUM(obj) == 683 || GET_OBJ_VNUM(obj) == 684) {
				act(libc.CString("$n looks in $p."), TRUE, ch, obj, nil, TO_ROOM)
			}
			switch bits {
			case (1 << 2):
				send_to_char(ch, libc.CString(" (carried): \r\n"))
			case (1 << 3):
				send_to_char(ch, libc.CString(" (here): \r\n"))
			case (1 << 5):
				send_to_char(ch, libc.CString(" (used): \r\n"))
			}
			list_obj_to_char(obj.Contains, ch, SHOW_OBJ_SHORT, TRUE)
		}
	} else {
		if (obj.Value[VAL_DRINKCON_HOWFULL]) <= 0 && (obj.Value[VAL_DRINKCON_CAPACITY]) == 0 {
			send_to_char(ch, libc.CString("It is empty.\r\n"))
		} else {
			if (obj.Value[VAL_DRINKCON_CAPACITY]) < 0 {
				var buf2 [64936]byte
				sprinttype(obj.Value[VAL_DRINKCON_LIQUID], color_liquid[:], &buf2[0], uint64(64936))
				send_to_char(ch, libc.CString("It's full of a %s liquid.\r\n"), &buf2[0])
			} else if (obj.Value[VAL_DRINKCON_HOWFULL]) > (obj.Value[VAL_DRINKCON_CAPACITY]) {
				send_to_char(ch, libc.CString("Its contents seem somewhat murky.\r\n"))
			} else {
				var buf2 [64936]byte
				amt = obj.Value[VAL_DRINKCON_CAPACITY]
				var leftin int = (obj.Value[VAL_DRINKCON_HOWFULL])
				sprinttype(obj.Value[VAL_DRINKCON_LIQUID], color_liquid[:], &buf2[0], uint64(64936))
				if leftin == amt {
					send_to_char(ch, libc.CString("It's full of a %s liquid.\r\n"), &buf2[0])
				} else if float64(leftin) >= float64(amt)*0.8 {
					send_to_char(ch, libc.CString("It's almost full of a %s liquid.\r\n"), &buf2[0])
				} else if float64(leftin) >= float64(amt)*0.5 {
					send_to_char(ch, libc.CString("It's about half full of a %s liquid.\r\n"), &buf2[0])
				} else if float64(leftin) >= float64(amt)*0.2 {
					send_to_char(ch, libc.CString("It's less than half full of a %s liquid.\r\n"), &buf2[0])
				} else if leftin > 0 {
					send_to_char(ch, libc.CString("It's barely filled with a %s liquid.\r\n"), &buf2[0])
				} else {
					send_to_char(ch, libc.CString("It's empty.\r\n"))
				}
			}
		}
	}
}
func find_exdesc(word *byte, list *extra_descr_data) *byte {
	var i *extra_descr_data
	for i = list; i != nil; i = i.Next {
		if func() int {
			if *i.Keyword == '.' {
				return isname(word, (*byte)(unsafe.Add(unsafe.Pointer(i.Keyword), 1)))
			}
			return isname(word, i.Keyword)
		}() != 0 {
			return i.Description
		}
	}
	return nil
}
func look_at_target(ch *char_data, arg *byte, cmread int) {
	var (
		bits       int
		found      int = FALSE
		j          int
		fnum       int
		i          int        = 0
		msg        int        = 1
		found_char *char_data = nil
		obj        *obj_data
		found_obj  *obj_data = nil
		desc       *byte
		number     [64936]byte
	)
	if ch.Desc == nil {
		return
	}
	if *arg == 0 {
		send_to_char(ch, libc.CString("Look at what?\r\n"))
		return
	}
	if cmread != 0 {
		for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
			if int(obj.Type_flag) == ITEM_BOARD {
				found = TRUE
				break
			}
		}
		if obj == nil {
			for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj.Next_content {
				if int(obj.Type_flag) == ITEM_BOARD {
					found = TRUE
					break
				}
			}
		}
		if obj != nil {
			arg = one_argument(arg, &number[0])
			if number[0] == 0 {
				send_to_char(ch, libc.CString("Read what?\r\n"))
				return
			}
			if isname(&number[0], obj.Name) != 0 {
				show_board(GET_OBJ_VNUM(obj), ch)
			} else if !unicode.IsDigit(rune(number[0])) || (func() int {
				msg = libc.Atoi(libc.GoString(&number[0]))
				return msg
			}()) == 0 || libc.StrChr(&number[0], '.') != nil {
				stdio.Sprintf(arg, "%s %s", &number[0], arg)
				look_at_target(ch, arg, 0)
			} else {
				board_display_msg(GET_OBJ_VNUM(obj), ch, msg)
			}
		}
	} else {
		bits = generic_find(arg, (1<<2)|1<<3|1<<5|1<<0, ch, &found_char, &found_obj)
		if found_char != nil {
			look_at_char(found_char, ch)
			if ch != found_char {
				if !AFF_FLAGGED(ch, AFF_HIDE) {
					act(libc.CString("$n looks at you."), TRUE, ch, nil, unsafe.Pointer(found_char), TO_VICT)
					act(libc.CString("$n looks at $N."), TRUE, ch, nil, unsafe.Pointer(found_char), TO_NOTVICT)
				}
			}
			return
		}
		if (func() int {
			fnum = get_number(&arg)
			return fnum
		}()) == 0 {
			send_to_char(ch, libc.CString("Look at what?\r\n"))
			return
		}
		if (func() *byte {
			desc = find_exdesc(arg, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Ex_description)
			return desc
		}()) != nil && func() int {
			p := &i
			*p++
			return *p
		}() == fnum {
			page_string(ch.Desc, desc, FALSE)
			return
		}
		for j = 0; j < NUM_WEARS && found == 0; j++ {
			if (ch.Equipment[j]) != nil && CAN_SEE_OBJ(ch, ch.Equipment[j]) {
				if (func() *byte {
					desc = find_exdesc(arg, (ch.Equipment[j]).Ex_description)
					return desc
				}()) != nil && func() int {
					p := &i
					*p++
					return *p
				}() == fnum {
					send_to_char(ch, libc.CString("%s"), desc)
					if isname(arg, (ch.Equipment[j]).Name) != 0 {
						if int((ch.Equipment[j]).Type_flag) == ITEM_WEAPON {
							send_to_char(ch, libc.CString("The weapon type of %s is a %s.\r\n"), (ch.Equipment[j]).Short_description, weapon_type[(ch.Equipment[j]).Value[VAL_WEAPON_SKILL]])
						}
						if int((ch.Equipment[j]).Type_flag) == ITEM_SPELLBOOK {
							display_spells(ch, ch.Equipment[j])
						}
						if int((ch.Equipment[j]).Type_flag) == ITEM_SCROLL {
							display_scroll(ch, ch.Equipment[j])
						}
						diag_obj_to_char(ch.Equipment[j], ch)
						send_to_char(ch, libc.CString("It appears to be made of %s"), material_names[(ch.Equipment[j]).Value[VAL_ALL_MATERIAL]])
					}
					found = TRUE
				}
			}
		}
		for obj = ch.Carrying; obj != nil && found == 0; obj = obj.Next_content {
			if CAN_SEE_OBJ(ch, obj) {
				if (func() *byte {
					desc = find_exdesc(arg, obj.Ex_description)
					return desc
				}()) != nil && func() int {
					p := &i
					*p++
					return *p
				}() == fnum {
					if int(obj.Type_flag) == ITEM_BOARD {
						show_board(GET_OBJ_VNUM(obj), ch)
					} else {
						send_to_char(ch, libc.CString("%s"), desc)
						if isname(arg, obj.Name) != 0 {
							if int(obj.Type_flag) == ITEM_WEAPON {
								send_to_char(ch, libc.CString("The weapon type of %s is a %s.\r\n"), obj.Short_description, weapon_type[obj.Value[VAL_WEAPON_SKILL]])
							}
							if int(obj.Type_flag) == ITEM_SPELLBOOK {
								display_spells(ch, obj)
							}
							if int(obj.Type_flag) == ITEM_SCROLL {
								display_scroll(ch, obj)
							}
							diag_obj_to_char(obj, ch)
							send_to_char(ch, libc.CString("It appears to be made of %s, and weights %s"), material_names[obj.Value[VAL_ALL_MATERIAL]], add_commas(obj.Weight))
						}
					}
					found = TRUE
				}
			}
		}
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil && found == 0; obj = obj.Next_content {
			if CAN_SEE_OBJ(ch, obj) {
				if (func() *byte {
					desc = find_exdesc(arg, obj.Ex_description)
					return desc
				}()) != nil && func() int {
					p := &i
					*p++
					return *p
				}() == fnum {
					if int(obj.Type_flag) == ITEM_BOARD {
						show_board(GET_OBJ_VNUM(obj), ch)
					} else {
						send_to_char(ch, libc.CString("%s"), desc)
						if int(obj.Type_flag) == ITEM_VEHICLE {
							send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @COpen hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CClose hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CEnter hatch\r\n"))
						} else if int(obj.Type_flag) == ITEM_HATCH {
							send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @COpen hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CClose hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CUnlock hatch\r\n"))
							send_to_char(ch, libc.CString("@YSyntax@D: @CLeave@n\r\n"))
						} else if int(obj.Type_flag) == ITEM_WINDOW {
							look_out_window(ch, obj.Name)
						}
						if int(obj.Type_flag) == ITEM_CONTROL {
							send_to_char(ch, libc.CString("@RFUEL@D: %s%s@n\r\n"), func() string {
								if (obj.Value[2]) >= 200 {
									return "@G"
								}
								if (obj.Value[2]) >= 100 {
									return "@Y"
								}
								return "@r"
							}(), add_commas(int64(obj.Value[2])))
						}
						if int(obj.Type_flag) == ITEM_WEAPON {
							send_to_char(ch, libc.CString("The weapon type of %s is a %s.\r\n"), obj.Short_description, weapon_type[obj.Value[VAL_WEAPON_SKILL]])
						}
						diag_obj_to_char(obj, ch)
						send_to_char(ch, libc.CString("It appears to be made of %s, and weights %s"), material_names[obj.Value[VAL_ALL_MATERIAL]], add_commas(obj.Weight))
					}
					found = TRUE
				}
			}
		}
		if bits != 0 {
			if found == 0 {
				show_obj_to_char(found_obj, ch, SHOW_OBJ_ACTION)
			} else {
				if show_obj_modifiers(found_obj, ch) != 0 {
					send_to_char(ch, libc.CString("\r\n"))
				}
			}
		} else if found == 0 {
			send_to_char(ch, libc.CString("You do not see that here.\r\n"))
		}
	}
}
func look_out_window(ch *char_data, arg *byte) {
	var (
		i           *obj_data
		viewport    *obj_data  = nil
		vehicle     *obj_data  = nil
		dummy       *char_data = nil
		target_room room_rnum  = room_rnum(-1)
		bits        int
		door        int
	)
	if *arg != 0 {
		if (func() int {
			bits = generic_find(arg, (1<<3)|1<<2|1<<5, ch, &dummy, &viewport)
			return bits
		}()) == 0 {
			send_to_char(ch, libc.CString("You don't see that here.\r\n"))
			return
		} else if int(viewport.Type_flag) != ITEM_WINDOW {
			send_to_char(ch, libc.CString("You can't look out that!\r\n"))
			return
		}
	} else if OUTSIDE(ch) {
		send_to_char(ch, libc.CString("But you are already outside.\r\n"))
		return
	} else {
		for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
			if int(i.Type_flag) == ITEM_WINDOW && isname(libc.CString("window"), i.Name) != 0 {
				viewport = i
				continue
			}
		}
	}
	if viewport == nil {
		send_to_char(ch, libc.CString("You don't seem to be able to see outside.\r\n"))
	} else if OBJVAL_FLAGGED(viewport, 1<<0) && OBJVAL_FLAGGED(viewport, 1<<2) {
		send_to_char(ch, libc.CString("It is closed.\r\n"))
	} else {
		if (viewport.Value[VAL_WINDOW_UNUSED1]) < 0 {
			if (viewport.Value[VAL_WINDOW_UNUSED4]) < 0 {
				for door = 0; door < NUM_OF_DIRS; door++ {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
						if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room != room_rnum(-1) {
							if !ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ROOM_INDOORS) {
								target_room = ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room
								continue
							}
						}
					}
				}
			} else {
				target_room = real_room(room_vnum(viewport.Value[VAL_WINDOW_UNUSED4]))
			}
		} else {
			if (func() *obj_data {
				vehicle = find_vehicle_by_vnum(viewport.Value[VAL_WINDOW_UNUSED1])
				return vehicle
			}()) != nil {
				target_room = vehicle.In_room
			}
		}
		if target_room == room_rnum(-1) {
			send_to_char(ch, libc.CString("You don't seem to be able to see outside.\r\n"))
		} else {
			if viewport.Action_description != nil {
				act(viewport.Action_description, TRUE, ch, viewport, nil, TO_CHAR)
			} else {
				act(libc.CString("$n looks out the window."), TRUE, ch, nil, nil, TO_ROOM)
			}
			send_to_char(ch, libc.CString("You look outside and see:\r\n"))
			look_at_room(target_room, ch, 0)
		}
	}
}
func do_finger(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What user are you wanting to look at?\r\n"))
		return
	}
	if readUserIndex(&arg[0]) == 0 {
		send_to_char(ch, libc.CString("That user does not exist\r\n"))
		return
	} else {
		fingerUser(ch, &arg[0])
		return
	}
}
func do_rptrans(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		k    *descriptor_data
		amt  int = 0
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: exchange (target) (amount)\r\n"))
		return
	}
	amt = libc.Atoi(libc.GoString(&arg2[0]))
	if amt <= 0 {
		send_to_char(ch, libc.CString("Are you being funny?\r\n"))
		return
	}
	if amt > ch.Rp {
		send_to_char(ch, libc.CString("@WYou only have @C%d@W RPP!@n\r\n"), ch.Rp)
		return
	}
	if readUserIndex(&arg[0]) == 0 {
		send_to_char(ch, libc.CString("That is not a recognised user file.\r\n"))
		return
	}
	for k = descriptor_list; k != nil; k = k.Next {
		if IS_NPC(k.Character) {
			continue
		}
		if k.Connected != CON_PLAYING {
			continue
		}
		if libc.StrCaseCmp(k.User, &arg[0]) == 0 {
			vict = k.Character
		}
	}
	if vict == nil {
		userWrite(nil, 0, amt, 0, &arg[0])
	} else {
		vict.Rp += amt
		vict.Desc.Rpp += amt
		userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		vict.Trp += amt
		save_char(vict)
		send_to_char(vict, libc.CString("@W%s gives @C%d@W of their RPP to you. How nice!\r\n"), GET_NAME(ch), amt)
	}
	ch.Rp -= amt
	ch.Desc.Rpp -= amt
	userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
	send_to_char(ch, libc.CString("@WYou exchange @C%d@W RPP to user @c%s@W for a warm fuzzy feeling.\r\n"), amt, CAP(&arg[0]))
	mudlog(NRM, MAX(ADMLVL_IMPL, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("EXCHANGE: %s gave %d RPP to user %s"), GET_NAME(ch), amt, &arg[0])
	save_char(ch)
}
func do_rpbank(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		amt int = 0
		arg [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: rpbank <amount>\r\n"))
		send_to_char(ch, libc.CString("Current RPP Bank: %d\r\n"), ch.Rbank)
		return
	}
	amt = libc.Atoi(libc.GoString(&arg[0]))
	if amt <= 0 {
		send_to_char(ch, libc.CString("You cannot withdraw from the RPP Bank.\r\n"))
		return
	}
	if amt > ch.Rp {
		send_to_char(ch, libc.CString("You do not have that much RPP to send to the Bank!\r\n"))
		return
	}
	ch.Rp -= amt
	ch.Desc.Rpp -= amt
	ch.Rbank += amt
	ch.Desc.Rbank += amt
	send_to_char(ch, libc.CString("You send %d to your RPP Bank. Your total is now %d.\r\n"), amt, ch.Rbank)
	mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("RPP Bank: %s has put %d RPP into their bank"), GET_NAME(ch), amt)
}
func do_rbanktrans(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		k    *descriptor_data
		amt  int = 0
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: bexchange (target) (amount)\r\n"))
		return
	}
	amt = libc.Atoi(libc.GoString(&arg2[0]))
	if amt <= 0 {
		send_to_char(ch, libc.CString("Are you being funny?\r\n"))
		return
	}
	if amt > ch.Rbank {
		send_to_char(ch, libc.CString("@WYou only have @C%d@W Banked RPP!@n\r\n"), ch.Rbank)
		return
	}
	if readUserIndex(&arg[0]) == 0 {
		send_to_char(ch, libc.CString("That is not a recognised user file.\r\n"))
		return
	}
	for k = descriptor_list; k != nil; k = k.Next {
		if IS_NPC(k.Character) {
			continue
		}
		if k.Connected != CON_PLAYING {
			continue
		}
		if libc.StrCaseCmp(k.User, &arg[0]) == 0 {
			vict = k.Character
		}
	}
	if vict == nil {
		userWrite(nil, 0, amt, 0, &arg[0])
	} else {
		vict.Rbank += amt
		vict.Desc.Rbank += amt
		userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		save_char(vict)
		send_to_char(vict, libc.CString("@W%s gives @C%d@W of their Banked RPP to you. How nice!\r\n"), GET_NAME(ch), amt)
	}
	ch.Rbank -= amt
	ch.Desc.Rbank -= amt
	userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
	send_to_char(ch, libc.CString("@wWELL @xGOLLY @xGEE @xWILLICKERS@x! @wYOU @RZ@YIM @RZ@YAMMED @C%d @wRIPROOZLES TO @C%s!@n\r\n"), amt, CAP(&arg[0]))
	mudlog(NRM, MAX(ADMLVL_IMPL, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("EXCHANGE: %s gave %d Banked RPP to user %s"), GET_NAME(ch), amt, &arg[0])
	save_char(ch)
}
func do_rdisplay(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	if IS_NPC(ch) {
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Clearing room display.\r\n"))
		ch.Rdisplay = libc.CString("Empty")
	} else {
		var derp [64936]byte
		libc.StrCpy(&derp[0], argument)
		send_to_char(ch, libc.CString("You set your display to; %s\r\n"), &derp[0])
		ch.Rdisplay = libc.StrDup(&derp[0])
	}
}
func perf_skill(skill int) int {
	if skill == 0 {
		return 0
	}
	switch skill {
	case 464:
		return 1
	case 441:
		return 1
	case 444:
		return 1
	case 475:
		return 1
	case 474:
		return 1
	case 476:
		return 1
	case 488:
		return 1
	case 472:
		return 1
	case 485:
		return 1
	case 442:
		return 1
	case 510:
		return 1
	case 533:
		return 1
	default:
		return 0
	}
}
func do_perf(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		arg2  [2048]byte
		i     int
		skill int = 1
		found int = FALSE
		type_ int = 0
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if IS_NPC(ch) || ch.Admlevel > 0 {
		send_to_char(ch, libc.CString("I don't think so.\r\n"))
		return
	}
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("@WType @G1@D: @wOver Charged@n\r\n"))
		send_to_char(ch, libc.CString("@WType @G2@D: @wAccurate@n\r\n"))
		send_to_char(ch, libc.CString("@WType @G3@D: @wEfficient@n\r\n"))
		send_to_char(ch, libc.CString("Syntax: perfect (skillname) (type 1/2/or 3)\r\n"))
		return
	}
	if libc.StrLen(&arg[0]) < 4 {
		send_to_char(ch, libc.CString("The skill name should be longer than 3 characters...\r\n"))
		return
	}
	for i = 1; i <= SKILL_TABLE_SIZE; i++ {
		if spell_info[i].Skilltype != (1 << 1) {
			continue
		}
		if found == TRUE {
			continue
		}
		if libc.StrStr(spell_info[i].Name, &arg[0]) != nil {
			skill = i
			found = TRUE
		}
	}
	if found == FALSE {
		send_to_char(ch, libc.CString("The skill %s doesn't exist.\r\n"), &arg[0])
		return
	}
	if GET_SKILL(ch, skill) == 0 {
		send_to_char(ch, libc.CString("You don't know %s.\r\n"), &arg[0])
		return
	}
	if GET_SKILL(ch, skill) < 100 {
		send_to_char(ch, libc.CString("You have not mastered the skill %s and thus can't perfect it.\r\n"), &arg[0])
		return
	}
	if int(ch.Skillperfs[skill]) > 0 {
		send_to_char(ch, libc.CString("You have already mastered the skill %s and chosen how to perfect it.\r\n"), &arg[0])
		return
	}
	if perf_skill(skill) == 0 {
		send_to_char(ch, libc.CString("You can't perfect that type of skill.\r\n"))
		return
	}
	if libc.Atoi(libc.GoString(&arg2[0])) < 1 || libc.Atoi(libc.GoString(&arg2[0])) > 3 {
		send_to_char(ch, libc.CString("@WType @G1@D: @wOver Charged@n\r\n"))
		send_to_char(ch, libc.CString("@WType @G2@D: @wAccurate@n\r\n"))
		send_to_char(ch, libc.CString("@WType @G3@D: @wEfficient@n\r\n"))
		send_to_char(ch, libc.CString("@RType must be a number between 1 and 3.@n\r\n"))
		return
	} else {
		type_ = libc.Atoi(libc.GoString(&arg2[0]))
		switch type_ {
		case 1:
			send_to_char(ch, libc.CString("You perfect the skill %s so that you can over charge it!\r\n"), spell_info[skill].Name)
			ch.Skillperfs[skill] = 1
		case 2:
			send_to_char(ch, libc.CString("You perfect the skill %s so that you have supreme accuracy with it!\r\n"), spell_info[skill].Name)
			ch.Skillperfs[skill] = 2
		case 3:
			send_to_char(ch, libc.CString("You perfect the skill %s so that you require a lower minimum charge for it!\r\n"), spell_info[skill].Name)
			ch.Skillperfs[skill] = 3
		}
	}
}
func do_look(ch *char_data, argument *byte, cmd int, subcmd int) {
	var look_type int
	if ch.Desc == nil {
		return
	}
	if int(ch.Position) < POS_SLEEPING {
		send_to_char(ch, libc.CString("You can't see anything but stars!\r\n"))
	} else if AFF_FLAGGED(ch, AFF_BLIND) {
		send_to_char(ch, libc.CString("You can't see a damned thing, you're blind!\r\n"))
	} else if PLR_FLAGGED(ch, PLR_EYEC) {
		send_to_char(ch, libc.CString("You can't see a damned thing, your eyes are closed!\r\n"))
	} else if room_is_dark(ch.In_room) != 0 && !CAN_SEE_IN_DARK(ch) && !PLR_FLAGGED(ch, PLR_AURALIGHT) {
		send_to_char(ch, libc.CString("It is pitch black...\r\n"))
		list_char_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People, ch)
	} else {
		var (
			arg  [2048]byte
			arg2 [200]byte
		)
		if subcmd == SCMD_READ {
			one_argument(argument, &arg[0])
			if arg[0] == 0 {
				send_to_char(ch, libc.CString("Read what?\r\n"))
			} else {
				look_at_target(ch, &arg[0], 1)
			}
			return
		}
		argument = any_one_arg(argument, &arg[0])
		one_argument(argument, &arg2[0])
		if arg[0] == 0 {
			if subcmd == SCMD_SEARCH {
				search_room(ch)
			} else {
				look_at_room(ch.In_room, ch, 1)
				if ch.Admlevel < 1 && !AFF_FLAGGED(ch, AFF_HIDE) {
					act(libc.CString("@w$n@w looks around the room.@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
		} else if is_abbrev(&arg[0], libc.CString("inside")) != 0 && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[INDIR]) != nil && arg2[0] == 0 {
			if subcmd == SCMD_SEARCH {
				search_in_direction(ch, INDIR)
			} else {
				look_in_direction(ch, INDIR)
			}
		} else if is_abbrev(&arg[0], libc.CString("inside")) != 0 && subcmd == SCMD_SEARCH && arg2[0] == 0 {
			search_in_direction(ch, INDIR)
		} else if is_abbrev(&arg[0], libc.CString("inside")) != 0 || is_abbrev(&arg[0], libc.CString("into")) != 0 || is_abbrev(&arg[0], libc.CString("onto")) != 0 {
			look_in_obj(ch, &arg2[0])
		} else if (is_abbrev(&arg[0], libc.CString("outside")) != 0 || is_abbrev(&arg[0], libc.CString("through")) != 0 || is_abbrev(&arg[0], libc.CString("thru")) != 0) && subcmd == SCMD_LOOK && arg2[0] != 0 {
			look_out_window(ch, &arg2[0])
		} else if is_abbrev(&arg[0], libc.CString("outside")) != 0 && subcmd == SCMD_LOOK && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[OUTDIR]) == nil {
			look_out_window(ch, &arg2[0])
		} else if (func() int {
			look_type = search_block(&arg[0], &dirs[0], FALSE)
			return look_type
		}()) >= 0 || (func() int {
			look_type = search_block(&arg[0], &abbr_dirs[0], FALSE)
			return look_type
		}()) >= 0 {
			if subcmd == SCMD_SEARCH {
				search_in_direction(ch, look_type)
			} else {
				look_in_direction(ch, look_type)
			}
		} else if is_abbrev(&arg[0], libc.CString("towards")) != 0 && ((func() int {
			look_type = search_block(&arg2[0], &dirs[0], FALSE)
			return look_type
		}()) >= 0 || (func() int {
			look_type = search_block(&arg2[0], &abbr_dirs[0], FALSE)
			return look_type
		}()) >= 0) {
			if subcmd == SCMD_SEARCH {
				search_in_direction(ch, look_type)
			} else {
				look_in_direction(ch, look_type)
			}
		} else if is_abbrev(&arg[0], libc.CString("at")) != 0 {
			if subcmd == SCMD_SEARCH {
				send_to_char(ch, libc.CString("That is not a direction!\r\n"))
			} else {
				look_at_target(ch, &arg2[0], 0)
			}
		} else if is_abbrev(&arg[0], libc.CString("around")) != 0 {
			var (
				i     *extra_descr_data
				found int = 0
			)
			for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Ex_description; i != nil; i = i.Next {
				if *i.Keyword != '.' {
					send_to_char(ch, libc.CString("%s%s:\r\n%s"), func() string {
						if found != 0 {
							return "\r\n"
						}
						return ""
					}(), i.Keyword, i.Description)
					found = 1
				}
			}
			if found == 0 {
				send_to_char(ch, libc.CString("You couldn't find anything noticeable.\r\n"))
			}
		} else if find_exdesc(&arg[0], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Ex_description) != nil {
			look_at_target(ch, &arg[0], 0)
		} else {
			if subcmd == SCMD_SEARCH {
				send_to_char(ch, libc.CString("That is not a direction!\r\n"))
			} else {
				look_at_target(ch, &arg[0], 0)
			}
		}
	}
}
func do_examine(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		tmp_char   *char_data
		tmp_object *obj_data
		tempsave   [2048]byte
		arg        [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Examine what?\r\n"))
		return
	}
	look_at_target(ch, libc.StrCpy(&tempsave[0], &arg[0]), 0)
	generic_find(&arg[0], (1<<2)|1<<3|1<<0|1<<5, ch, &tmp_char, &tmp_object)
	if tmp_object != nil {
		if int(tmp_object.Type_flag) == ITEM_DRINKCON || int(tmp_object.Type_flag) == ITEM_FOUNTAIN || int(tmp_object.Type_flag) == ITEM_CONTAINER {
			send_to_char(ch, libc.CString("When you look inside, you see:\r\n"))
			look_in_obj(ch, &arg[0])
		}
	}
}
func do_gold(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Gold == 0 {
		send_to_char(ch, libc.CString("You're broke!\r\n"))
	} else if ch.Gold == 1 {
		send_to_char(ch, libc.CString("You have one little zenni.\r\n"))
	} else {
		send_to_char(ch, libc.CString("You have %d zenni.\r\n"), ch.Gold)
	}
}
func do_score(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var view int = 0
	var full int = 5
	var personal int = 1
	var health int = 2
	var stats int = 3
	var other int = 4
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		view = full
	} else if libc.StrStr(libc.CString("personal"), &arg[0]) != nil || libc.StrStr(libc.CString("Personal"), &arg[0]) != nil {
		view = personal
	} else if libc.StrStr(libc.CString("health"), &arg[0]) != nil || libc.StrStr(libc.CString("Health"), &arg[0]) != nil {
		view = health
	} else if libc.StrStr(libc.CString("statistics"), &arg[0]) != nil || libc.StrStr(libc.CString("Statistics"), &arg[0]) != nil {
		view = stats
	} else if libc.StrStr(libc.CString("other"), &arg[0]) != nil || libc.StrStr(libc.CString("Other"), &arg[0]) != nil {
		view = other
	} else {
		send_to_char(ch, libc.CString("Syntax: score, or... score (personal, health, statistics, other)\r\n"))
		return
	}
	if view == full || view == personal {
		send_to_char(ch, libc.CString("  @cO@D-----------------------------[  @cPersonal  @D]-----------------------------@cO@n\n"))
		send_to_char(ch, libc.CString("  @D|  @CName@D: @W%15s@D,   @CTitle@D: @W%-38s@D|@n\n"), GET_NAME(ch), GET_TITLE(ch))
		if int(ch.Race) == RACE_ANDROID {
			var (
				model   [100]byte
				version [100]byte
				absorb  int = 0
			)
			if PLR_FLAGGED(ch, PLR_ABSORB) {
				stdio.Sprintf(&model[0], "@CAbsorption")
			} else if PLR_FLAGGED(ch, PLR_REPAIR) {
				stdio.Sprintf(&model[0], "@GSelf Repairing")
			} else if PLR_FLAGGED(ch, PLR_SENSEM) {
				stdio.Sprintf(&model[0], "@RSensor Equiped")
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				stdio.Sprintf(&version[0], "Beta 1.0")
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				stdio.Sprintf(&version[0], "ANS 2.0")
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				stdio.Sprintf(&version[0], "ANS 3.0")
			} else if PLR_FLAGGED(ch, PLR_TRANS4) {
				stdio.Sprintf(&version[0], "ANS 4.0")
			} else if PLR_FLAGGED(ch, PLR_TRANS5) {
				stdio.Sprintf(&version[0], "ANS 5.0")
			} else if PLR_FLAGGED(ch, PLR_TRANS6) {
				stdio.Sprintf(&version[0], "ANS 6.0")
			} else {
				stdio.Sprintf(&version[0], "Alpha 0.5")
			}
			send_to_char(ch, libc.CString("  @D| @CModel@D: %15s@D,    @CUGP@D: @G%15s@D,  @CVersion@D: @r%-12s@D|@n\n"), &model[0], func() string {
				if absorb > 0 {
					return "@RN/A"
				}
				return libc.GoString(add_commas(int64(ch.Upgrade)))
			}(), &version[0])
		}
		if ch.Clan != nil {
			send_to_char(ch, libc.CString("  @D|  @CClan@D: @W%-64s@D|@n\n"), ch.Clan)
		}
		send_to_char(ch, libc.CString("  @D|  @CRace@D: @W%10s@D,  @CSensei@D: @W%15s@D,     @CArt@D: @W%-17s@D|@n\n"), pc_race_types[int(ch.Race)], pc_class_types[int(ch.Chclass)], sensei_style[ch.Chclass])
		var hei [300]byte
		var wei [300]byte
		stdio.Sprintf(&hei[0], "%dcm", get_measure(ch, GET_PC_HEIGHT(ch), 0))
		stdio.Sprintf(&wei[0], "%dkg", get_measure(ch, 0, GET_PC_WEIGHT(ch)))
		send_to_char(ch, libc.CString("  @D|   @CAge@D: @W%10s@D,  @CHeight@D: @W%15s@D,  @CWeight@D: @W%-17s@D|@n\n"), add_commas(int64(age(ch).Year)), &hei[0], &wei[0])
		send_to_char(ch, libc.CString("  @D|@CGender@D: @W%10s@D,  @C  Size@D: @W%15s@D,  @C Align@D: @W%-17s@D|@n\n"), genders[int(ch.Sex)], size_names[get_size(ch)], disp_align(ch))
	}
	if view == full || view == health {
		send_to_char(ch, libc.CString("  @cO@D-----------------------------@D[   @cHealth   @D]-----------------------------@cO@n\n"))
		send_to_char(ch, libc.CString("                 @D<@rPowerlevel@D>          <@BKi@D>             <@GStamina@D>@n\n"))
		send_to_char(ch, libc.CString("    @wCurrent   @D-[@R%-16s@D]-[@R%-16s@D]-[@R%-16s@D]@n\n"), add_commas(ch.Hit), add_commas(ch.Mana), add_commas(ch.Move))
		send_to_char(ch, libc.CString("    @wMaximum   @D-[@r%-16s@D]-[@r%-16s@D]-[@r%-16s@D]@n\n"), add_commas(ch.Max_hit), add_commas(ch.Max_mana), add_commas(ch.Max_move))
		send_to_char(ch, libc.CString("    @wBase      @D-[@m%-16s@D]-[@m%-16s@D]-[@m%-16s@D]@n\n"), add_commas(ch.Basepl), add_commas(ch.Baseki), add_commas(ch.Basest))
		if int(ch.Race) != RACE_ANDROID && ch.Lifeforce > 0 {
			send_to_char(ch, libc.CString("    @wLife Force@D-[@C%16s@D%s@c%16s@D]- @wLife Percent@D-[@Y%3d%s@D]@n\n"), add_commas(ch.Lifeforce), "/", add_commas(int64(GET_LIFEMAX(ch))), ch.Lifeperc, "%")
		} else if int(ch.Race) != RACE_ANDROID {
			send_to_char(ch, libc.CString("    @wLife Force@D-[@C%16s@D%s@c%16s@D]- @wLife Percent@D-[@Y%3d%s@D]@n\n"), add_commas(0), "/", add_commas(int64(GET_LIFEMAX(ch))), ch.Lifeperc, "%")
		}
	}
	if view == full || view == stats {
		send_to_char(ch, libc.CString("  @cO@D-----------------------------@D[ @cStatistics @D]-----------------------------@cO@n\n"))
		send_to_char(ch, libc.CString("      @D<@wCharacter Level@D: @w%-3d@D> <@wRPP@D: @w%-3d@D> <@wRPP Bank@D: @w%-3d@D>@n\n"), GET_LEVEL(ch), ch.Rp, ch.Rbank)
		send_to_char(ch, libc.CString("      @D<@wSpeed Index@D: @w%-15s@D> <@wArmor Index@D: @w%-15s@D>@n\n"), add_commas(int64(GET_SPEEDI(ch))), add_commas(int64(ch.Armor)))
		send_to_char(ch, libc.CString("        @D[    @RStrength@D|@G%3d@D] [     @YAgility@D|@G%3d@D] [       @BSpeed@D|@G%3d@D]@n\n"), ch.Aff_abils.Str, ch.Aff_abils.Dex, ch.Aff_abils.Cha)
		send_to_char(ch, libc.CString("        @D[@gConstitution@D|@G%3d@D] [@CIntelligence@D|@G%3d@D] [      @MWisdom@D|@G%3d@D]@n\n"), ch.Aff_abils.Con, ch.Aff_abils.Intel, ch.Aff_abils.Wis)
	}
	if view == full || view == other {
		send_to_char(ch, libc.CString("  @cO@D-----------------------------@D[   @cOther    @D]-----------------------------@cO@n\n"))
		send_to_char(ch, libc.CString("                @D<@YZenni@D>                 <@rInventory Weight@D>@n\n"))
		send_to_char(ch, libc.CString("      @D[   @CCarried@D| @W%-15s@D] [   @CCarried@D| @W%-15s@D]@n\n"), add_commas(int64(ch.Gold)), add_commas(int64(gear_weight(ch))))
		send_to_char(ch, libc.CString("      @D[      @CBank@D| @W%-15s@D] [ @CMax Carry@D| @W%-15s@D]@n\n"), add_commas(int64(ch.Bank_gold)), add_commas(max_carry_weight(ch)))
		send_to_char(ch, libc.CString("      @D[ @CMax Carry@D| @W%-15s@D]@n\n"), add_commas(int64(GOLD_CARRY(ch))))
		var numb int = 0
		if ch.Bank_gold > 99 {
			numb = (ch.Bank_gold / 100) * 2
		} else if ch.Bank_gold > 0 {
			numb = 1
		} else {
			numb = 0
		}
		if numb >= 7500 {
			numb = 7500
		}
		send_to_char(ch, libc.CString("      @D[  @CInterest@D| @W%-15s@D]\n"), add_commas(int64(numb)))
		if int(ch.Race) == RACE_ARLIAN {
			send_to_char(ch, libc.CString("                             @D<@GEvolution @D>@n\n"))
			send_to_char(ch, libc.CString("      @D[ @CEvo Level@D| @W%-15d@D] [   @CEvo Exp@D| @W%-15s@D]\n"), ch.Moltlevel, add_commas(ch.Moltexp))
			send_to_char(ch, libc.CString("      @D[ @CThreshold@D| @W%-15s@D]@n\n"), add_commas(molt_threshold(ch)))
		}
		if GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("                             @D<@gAdvancement@D>@n\n"))
		}
		if GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("      @D[@CExperience@D| @W%-15s@D] [@CNext Level@D| @W%-15s@D]@n\n"), add_commas(ch.Exp), add_commas(int64(level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp))))
			send_to_char(ch, libc.CString("      @D[  @CRpp Cost@D| @W%-15d@D]@n\n"), rpp_to_level(ch))
		}
		send_to_char(ch, libc.CString("\n     @D<@wPlayed@D: @yYears @D(@W%2d@D) @yWeeks @D(@W%2d@D) @yDays @D(@W%2d@D) @yHours @D(@W%2d@D) @yMinutes @D(@W%2d@D)>@n\n"), int(ch.Time.Played)/31536000, int((ch.Time.Played%31536000)/604800), int((ch.Time.Played%604800)/86400), int((ch.Time.Played%86400)/3600), int((ch.Time.Played%3600)/60))
	}
	send_to_char(ch, libc.CString("  @cO@D------------------------------------------------------------------------@cO@n\n"))
}
func trans_check(ch *char_data, vict *char_data) {
	if int(vict.Race) == RACE_HUMAN {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Human First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Human Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Human Third@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Human Fourth@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_HOSHIJIN {
		if vict.Mimic == 0 || vict == ch {
			if vict.Starphase == 1 {
				send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CBirth Phase@n\r\n"))
			} else if vict.Starphase == 2 {
				send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CLife Phase@n\r\n"))
			} else {
				send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
			}
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if (int(vict.Race) == RACE_SAIYAN || int(vict.Race) == RACE_HALFBREED) && !PLR_FLAGGED(vict, PLR_LSSJ) {
		if PLR_FLAGGED(vict, PLR_TRANS1) && !PLR_FLAGGED(vict, PLR_FPSSJ) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Saiyan First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS1) && PLR_FLAGGED(vict, PLR_FPSSJ) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @YFull Powered @CSuper Saiyan@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Saiyan Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Saiyan Third@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Saiyan Fourth@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_OOZARU) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @COozaru@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_SAIYAN && PLR_FLAGGED(vict, PLR_LSSJ) {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Saiyan First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @YLegendary @CSuper Saiyan@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_OOZARU) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @COozaru@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_NAMEK {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Namek First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Namek Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Namek Third@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Namek Fourth@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_ICER {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CTransform First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CTransform Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CTransform Third@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CTransform Fourth@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_KONATSU {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CShadow First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CShadow Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CShadow Third@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CShadow Fourth@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_MUTANT {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMutate First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMutate Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMutate Third@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_BIO {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMature@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSemi-perfect@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CPerfect@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper Perfect@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_ANDROID {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 1.0@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 2.0@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 3.0@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS4) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 4.0@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS5) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 5.0@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS6) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSeries 6.0@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_MAJIN {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CAffinity@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CSuper@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CTrue@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_TRUFFLE {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CAscend First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CAscend Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CAscend Third@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else if int(vict.Race) == RACE_KAI {
		if PLR_FLAGGED(vict, PLR_TRANS1) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMystic First@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS2) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMystic Second@n\r\n"))
		} else if PLR_FLAGGED(vict, PLR_TRANS3) {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @CMystic Third@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("         @cCurrent Transformation@D: @wNone@n\r\n"))
	}
}
func do_status(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg          [2048]byte
		aff          *affected_type
		forget_level [7]*byte = [7]*byte{libc.CString("@GRemembered Well@n"), libc.CString("@GRemembered Well Enough@n"), libc.CString("@RGetting Foggy@n"), libc.CString("@RHalf Forgotten@n"), libc.CString("@rAlmost Forgotten@n"), libc.CString("@rForgotten@n"), libc.CString("\n")}
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@D<@b------------------------@D[@YYour Status@D]@b-------------------------@D>@n\r\n\r\n"))
		send_to_char(ch, libc.CString("            @D---------------@CAppearance@D---------------\n"))
		bringdesc(ch, ch)
		send_to_char(ch, libc.CString("            @D---------------@RAppendages@D---------------\n"))
		if PLR_FLAGGED(ch, PLR_HEAD) {
			send_to_char(ch, libc.CString("            @D[@cHead        @D: @GHave.          @D]@n\r\n"))
		}
		if !PLR_FLAGGED(ch, PLR_HEAD) {
			send_to_char(ch, libc.CString("            @D[@cHead        @D: @rMissing.         @D]@n\r\n"))
		}
		if (ch.Limb_condition[0]) >= 50 && !PLR_FLAGGED(ch, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), ch.Limb_condition[0], "%", "%")
		} else if (ch.Limb_condition[0]) > 0 && !PLR_FLAGGED(ch, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), ch.Limb_condition[0], "%", "%")
		} else if (ch.Limb_condition[0]) > 0 && PLR_FLAGGED(ch, PLR_CRARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), ch.Limb_condition[0], "%", "%")
		} else if (ch.Limb_condition[0]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cRight Arm   @D: @rMissing.         @D]@n\r\n"))
		}
		if (ch.Limb_condition[1]) >= 50 && !PLR_FLAGGED(ch, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), ch.Limb_condition[1], "%", "%")
		} else if (ch.Limb_condition[1]) > 0 && !PLR_FLAGGED(ch, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), ch.Limb_condition[1], "%", "%")
		} else if (ch.Limb_condition[1]) > 0 && PLR_FLAGGED(ch, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), ch.Limb_condition[1], "%", "%")
		} else if (ch.Limb_condition[1]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cLeft Arm    @D: @rMissing.         @D]@n\r\n"))
		}
		if (ch.Limb_condition[2]) >= 50 && !PLR_FLAGGED(ch, PLR_CLARM) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), ch.Limb_condition[2], "%", "%")
		} else if (ch.Limb_condition[2]) > 0 && !PLR_FLAGGED(ch, PLR_CRLEG) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), ch.Limb_condition[2], "%", "%")
		} else if (ch.Limb_condition[2]) > 0 && PLR_FLAGGED(ch, PLR_CRLEG) {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), ch.Limb_condition[2], "%", "%")
		} else if (ch.Limb_condition[2]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cRight Leg   @D: @rMissing.         @D]@n\r\n"))
		}
		if (ch.Limb_condition[3]) >= 50 && !PLR_FLAGGED(ch, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @G%2d%s@D/@g100%s        @D]@n\r\n"), ch.Limb_condition[3], "%", "%")
		} else if (ch.Limb_condition[3]) > 0 && !PLR_FLAGGED(ch, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @rBroken @y%2d%s@D/@g100%s @D]@n\r\n"), ch.Limb_condition[3], "%", "%")
		} else if (ch.Limb_condition[3]) > 0 && PLR_FLAGGED(ch, PLR_CLLEG) {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @cCybernetic @G%2d%s@D/@G100%s@D]@n\r\n"), ch.Limb_condition[3], "%", "%")
		} else if (ch.Limb_condition[3]) <= 0 {
			send_to_char(ch, libc.CString("            @D[@cLeft Leg    @D: @rMissing.         @D]@n\r\n"))
		}
		if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && PLR_FLAGGED(ch, PLR_STAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @GHave.            @D]@n\r\n"))
		}
		if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && !PLR_FLAGGED(ch, PLR_STAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @rMissing.         @D]@n\r\n"))
		}
		if (int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO) && PLR_FLAGGED(ch, PLR_TAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @GHave.            @D]@n\r\n"))
		}
		if (int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO) && !PLR_FLAGGED(ch, PLR_TAIL) {
			send_to_char(ch, libc.CString("            @D[@cTail        @D: @rMissing.         @D]@n\r\n"))
		}
		send_to_char(ch, libc.CString("\r\n"))
		send_to_char(ch, libc.CString("         @D-----------------@YHunger@D/@yThirst@D-----------------@n\r\n"))
		if int(ch.Player_specials.Conditions[HUNGER]) >= 48 {
			send_to_char(ch, libc.CString("         You are full.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 40 {
			send_to_char(ch, libc.CString("         You are nearly full.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 30 {
			send_to_char(ch, libc.CString("         You are not hungry.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 21 {
			send_to_char(ch, libc.CString("         You wouldn't mind a snack.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 15 {
			send_to_char(ch, libc.CString("         You are slightly hungry.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 10 {
			send_to_char(ch, libc.CString("         You are partially hungry.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 5 {
			send_to_char(ch, libc.CString("         You are really hungry.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 2 {
			send_to_char(ch, libc.CString("         You are extremely hungry.\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) >= 0 {
			send_to_char(ch, libc.CString("         You are starving!\r\n"))
		} else if int(ch.Player_specials.Conditions[HUNGER]) < 0 {
			send_to_char(ch, libc.CString("         You need not eat.\r\n"))
		}
		if int(ch.Player_specials.Conditions[THIRST]) >= 48 {
			send_to_char(ch, libc.CString("         You are not thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 40 {
			send_to_char(ch, libc.CString("         You are nearly quenched.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 30 {
			send_to_char(ch, libc.CString("         You are not thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 21 {
			send_to_char(ch, libc.CString("         You wouldn't mind a drink.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 15 {
			send_to_char(ch, libc.CString("         You are slightly thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 10 {
			send_to_char(ch, libc.CString("         You are partially thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 5 {
			send_to_char(ch, libc.CString("         You are really thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 2 {
			send_to_char(ch, libc.CString("         You are extremely thirsty.\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) >= 0 {
			send_to_char(ch, libc.CString("         You are dehydrated!\r\n"))
		} else if int(ch.Player_specials.Conditions[THIRST]) < 0 {
			send_to_char(ch, libc.CString("         You need not drink.\r\n"))
		}
		send_to_char(ch, libc.CString("         @D--------------------@D[@GInfo@D]---------------------@n\r\n"))
		trans_check(ch, ch)
		send_to_char(ch, libc.CString("         You have died %d times.\r\n"), ch.Dcount)
		if PLR_FLAGGED(ch, PLR_NOSHOUT) {
			send_to_char(ch, libc.CString("         You have been @rmuted@n on public channels.\r\n"))
		}
		if ch.In_room == real_room(9) {
			send_to_char(ch, libc.CString("         You are in punishment hell, so sad....\r\n"))
		}
		if !PRF_FLAGGED(ch, PRF_HINTS) {
			send_to_char(ch, libc.CString("         You have hints turned off.\r\n"))
		}
		if NEWSUPDATE > ch.Lastpl {
			send_to_char(ch, libc.CString("         Check the 'news', it has been updated recently.\r\n"))
		}
		if has_mail(int(ch.Idnum)) != 0 {
			send_to_char(ch, libc.CString("         Check your mail at the nearest postmaster.\r\n"))
		}
		if PRF_FLAGGED(ch, PRF_HIDE) {
			send_to_char(ch, libc.CString("         You are hidden from who and ooc.\r\n"))
		}
		if ch.Voice != nil {
			send_to_char(ch, libc.CString("         Your voice desc: '%s'\r\n"), ch.Voice)
		}
		if int(ch.Distfea) == DISTFEA_EYE {
			send_to_char(ch, libc.CString("         Your eyes are your most distinctive feature.\r\n"))
		}
		if ch.Preference == 0 {
			send_to_char(ch, libc.CString("         You preferred a balanced form of fighting.\r\n"))
		} else if ch.Preference == PREFERENCE_KI {
			send_to_char(ch, libc.CString("         You preferred a ki dominate form of fighting.\r\n"))
		} else if ch.Preference == PREFERENCE_WEAPON {
			send_to_char(ch, libc.CString("         You preferred a weapon dominate form of fighting.\r\n"))
		} else if ch.Preference == PREFERENCE_H2H {
			send_to_char(ch, libc.CString("         You preferred a body dominate form of fighting.\r\n"))
		} else if ch.Preference == PREFERENCE_THROWING {
			send_to_char(ch, libc.CString("         You preferred a throwing dominate form of fighting.\r\n"))
		}
		if int(ch.Distfea) == DISTFEA_HAIR && int(ch.Race) != RACE_DEMON && int(ch.Race) != RACE_MAJIN && int(ch.Race) != RACE_ICER && int(ch.Race) != RACE_NAMEK {
			send_to_char(ch, libc.CString("         Your hair is your most distinctive feature.\r\n"))
		} else if int(ch.Distfea) == DISTFEA_HAIR && int(ch.Race) == RACE_DEMON {
			send_to_char(ch, libc.CString("         Your horns are your most distinctive feature.\r\n"))
		} else if int(ch.Distfea) == DISTFEA_HAIR && int(ch.Race) == RACE_MAJIN {
			send_to_char(ch, libc.CString("         Your forelock is your most distinctive feature.\r\n"))
		} else if int(ch.Distfea) == DISTFEA_HAIR && int(ch.Race) == RACE_ICER {
			send_to_char(ch, libc.CString("         Your horns are your most distinctive feature.\r\n"))
		} else if int(ch.Distfea) == DISTFEA_HAIR && int(ch.Race) == RACE_NAMEK {
			send_to_char(ch, libc.CString("         Your antennae are your most distinctive feature.\r\n"))
		}
		if int(ch.Distfea) == DISTFEA_SKIN {
			send_to_char(ch, libc.CString("         Your skin is your most distinctive feature.\r\n"))
		}
		if int(ch.Distfea) == DISTFEA_HEIGHT {
			send_to_char(ch, libc.CString("         Your height is your most distinctive feature.\r\n"))
		}
		if int(ch.Distfea) == DISTFEA_WEIGHT {
			send_to_char(ch, libc.CString("         Your weight is your most distinctive feature.\r\n"))
		}
		if (ch.Equipment[WEAR_EYE]) != nil {
			var obj *obj_data = (ch.Equipment[WEAR_EYE])
			if obj.Scoutfreq == 0 {
				obj.Scoutfreq = 1
			}
			send_to_char(ch, libc.CString("         Your scouter is on frequency @G%d@n\r\n"), obj.Scoutfreq)
			obj = nil
		}
		if ch.Charge > 0 {
			send_to_char(ch, libc.CString("         You have @C%s@n ki charged.\r\n"), add_commas(ch.Charge))
		}
		if ch.Kaioken > 0 {
			send_to_char(ch, libc.CString("         You are focusing Kaioken x %d.\r\n"), ch.Kaioken)
		}
		if AFF_FLAGGED(ch, AFF_SANCTUARY) {
			send_to_char(ch, libc.CString("         You are surrounded by a barrier @D(@Y%s@D)@n\r\n"), add_commas(ch.Barrier))
		}
		if AFF_FLAGGED(ch, AFF_FIRESHIELD) {
			send_to_char(ch, libc.CString("         You are surrounded by flames!@n\r\n"))
		}
		if ch.Suppression > 0 {
			send_to_char(ch, libc.CString("         You are suppressing current PL to %lld.\r\n"), ch.Suppression)
		}
		if int(ch.Race) == RACE_MAJIN {
			send_to_char(ch, libc.CString("         You have ingested %d people.\r\n"), ch.Absorbs)
		}
		if int(ch.Race) == RACE_BIO {
			send_to_char(ch, libc.CString("         You have %d absorbs left.\r\n"), ch.Absorbs)
		}
		send_to_char(ch, libc.CString("         You have %s colored aura.\r\n"), aura_types[ch.Aura])
		if GET_LEVEL(ch) < 100 {
			if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) || int(ch.Race) != RACE_ANDROID && int(ch.Race) != RACE_BIO && int(ch.Race) != RACE_MAJIN {
				send_to_char(ch, libc.CString("         @R%s@n to SC a stat this level.\r\n"), add_commas(show_softcap(ch)))
			} else {
				send_to_char(ch, libc.CString("         @R%s@n in PL/KI/ST combined to SC this level.\r\n"), add_commas(show_softcap(ch)))
			}
		} else {
			send_to_char(ch, libc.CString("         Your strengths are potentially limitless.\r\n"))
		}
		if ch.Forgeting != 0 {
			send_to_char(ch, libc.CString("         @MForgetting @D[@m%s - %s@D]@n\r\n"), spell_info[ch.Forgeting].Name, forget_level[ch.Forgetcount])
		} else {
			send_to_char(ch, libc.CString("         @MForgetting @D[@mNothing.@D]@n\r\n"))
		}
		if GET_SKILL(ch, SKILL_DAGGER) > 0 {
			if ch.Backstabcool > 0 {
				send_to_char(ch, libc.CString("         @yYou can't preform a backstab yet.@n\r\n"))
			} else {
				send_to_char(ch, libc.CString("         @YYou can backstab.@n\r\n"))
			}
		}
		if ch.Feature != nil {
			send_to_char(ch, libc.CString("         Extra Feature: @C%s@n\r\n"), ch.Feature)
		}
		if ch.Rdisplay != nil {
			if ch.Rdisplay != libc.CString("Empty") {
				send_to_char(ch, libc.CString("         Room Display: @C...%s@n\r\n"), ch.Rdisplay)
			}
		}
		send_to_char(ch, libc.CString("\r\n@D<@b-------------------------@D[@BCondition@D]@b--------------------------@D>@n\r\n"))
		if (ch.Bonuses[BONUS_INSOMNIAC]) != 0 {
			send_to_char(ch, libc.CString("You can not sleep.\r\n"))
		} else {
			if ch.Sleeptime > 6 && int(ch.Position) != POS_SLEEPING {
				send_to_char(ch, libc.CString("You are well rested.\r\n"))
			} else if ch.Sleeptime > 6 && int(ch.Position) == POS_SLEEPING {
				send_to_char(ch, libc.CString("You are getting the rest you need.\r\n"))
			} else if ch.Sleeptime > 4 {
				send_to_char(ch, libc.CString("You are rested.\r\n"))
			} else if ch.Sleeptime > 2 {
				send_to_char(ch, libc.CString("You are not sleepy.\r\n"))
			} else if ch.Sleeptime >= 1 {
				send_to_char(ch, libc.CString("You are getting a little sleepy.\r\n"))
			} else if ch.Sleeptime == 0 {
				send_to_char(ch, libc.CString("You could sleep at any time.\r\n"))
			}
		}
		if ch.Relax_count > 464 {
			send_to_char(ch, libc.CString("You are far too at ease to train hard like you should. Get out of the house more often.\r\n"))
		} else if ch.Relax_count > 232 {
			send_to_char(ch, libc.CString("You are too at ease to train hard like you should. Get out of the house more often.\r\n"))
		} else if ch.Relax_count > 116 {
			send_to_char(ch, libc.CString("You are a bit at ease and your training suffers. Get out of the house more often.\r\n"))
		}
		if ch.Mimic > 0 {
			send_to_char(ch, libc.CString("You are mimicing the general appearance of %s %s\r\n"), AN(JUGGLERACELOWER(ch)), JUGGLERACELOWER(ch))
		}
		if int(ch.Race) == RACE_MUTANT {
			send_to_char(ch, libc.CString("Your Mutations:\r\n"))
			if (ch.Genome[0]) == 1 {
				send_to_char(ch, libc.CString("  Extreme Speed.\r\n"))
			}
			if (ch.Genome[0]) == 2 {
				send_to_char(ch, libc.CString("  Increased Cell Regeneration.\r\n"))
			}
			if (ch.Genome[0]) == 3 {
				send_to_char(ch, libc.CString("  Extreme Reflexes.\r\n"))
			}
			if (ch.Genome[0]) == 4 {
				send_to_char(ch, libc.CString("  Infravision.\r\n"))
			}
			if (ch.Genome[0]) == 5 {
				send_to_char(ch, libc.CString("  Natural Camo.\r\n"))
			}
			if (ch.Genome[0]) == 6 {
				send_to_char(ch, libc.CString("  Limb Regen.\r\n"))
			}
			if (ch.Genome[0]) == 7 {
				send_to_char(ch, libc.CString("  Poisonous (you can use the bite attack).\r\n"))
			}
			if (ch.Genome[0]) == 8 {
				send_to_char(ch, libc.CString("  Rubbery Body.\r\n"))
			}
			if (ch.Genome[0]) == 9 {
				send_to_char(ch, libc.CString("  Innate Telepathy.\r\n"))
			}
			if (ch.Genome[0]) == 10 {
				send_to_char(ch, libc.CString("  Natural Energy.\r\n"))
			}
			if (ch.Genome[1]) == 1 {
				send_to_char(ch, libc.CString("  Extreme Speed.\r\n"))
			}
			if (ch.Genome[1]) == 2 {
				send_to_char(ch, libc.CString("  Increased Cell Regeneration.\r\n"))
			}
			if (ch.Genome[1]) == 3 {
				send_to_char(ch, libc.CString("  Extreme Reflexes.\r\n"))
			}
			if (ch.Genome[1]) == 4 {
				send_to_char(ch, libc.CString("  Infravision.\r\n"))
			}
			if (ch.Genome[1]) == 5 {
				send_to_char(ch, libc.CString("  Natural Camo.\r\n"))
			}
			if (ch.Genome[1]) == 6 {
				send_to_char(ch, libc.CString("  Limb Regen.\r\n"))
			}
			if (ch.Genome[1]) == 7 {
				send_to_char(ch, libc.CString("  Poisonous (you can use the bite attack).\r\n"))
			}
			if (ch.Genome[1]) == 8 {
				send_to_char(ch, libc.CString("  Rubbery Body.\r\n"))
			}
			if (ch.Genome[1]) == 9 {
				send_to_char(ch, libc.CString("  Innate Telepathy.\r\n"))
			}
			if (ch.Genome[1]) == 10 {
				send_to_char(ch, libc.CString("  Natural Energy.\r\n"))
			}
		}
		if int(ch.Race) == RACE_BIO {
			send_to_char(ch, libc.CString("Your genes carry:\r\n"))
			if (ch.Genome[0]) == 1 {
				send_to_char(ch, libc.CString("  Human DNA.\r\n"))
			}
			if (ch.Genome[0]) == 2 {
				send_to_char(ch, libc.CString("  Saiyan DNA.\r\n"))
			}
			if (ch.Genome[0]) == 3 {
				send_to_char(ch, libc.CString("  Namek DNA.\r\n"))
			}
			if (ch.Genome[0]) == 4 {
				send_to_char(ch, libc.CString("  Icer DNA.\r\n"))
			}
			if (ch.Genome[0]) == 5 {
				send_to_char(ch, libc.CString("  Truffle DNA.\r\n"))
			}
			if (ch.Genome[0]) == 6 {
				send_to_char(ch, libc.CString("  Arlian DNA.\r\n"))
			}
			if (ch.Genome[0]) == 7 {
				send_to_char(ch, libc.CString("  Kai DNA.\r\n"))
			}
			if (ch.Genome[0]) == 8 {
				send_to_char(ch, libc.CString("  Konatsu DNA.\r\n"))
			}
			if (ch.Genome[1]) == 1 {
				send_to_char(ch, libc.CString("  Human DNA.\r\n"))
			}
			if (ch.Genome[1]) == 2 {
				send_to_char(ch, libc.CString("  Saiyan DNA.\r\n"))
			}
			if (ch.Genome[1]) == 3 {
				send_to_char(ch, libc.CString("  Namek DNA.\r\n"))
			}
			if (ch.Genome[1]) == 4 {
				send_to_char(ch, libc.CString("  Icer DNA.\r\n"))
			}
			if (ch.Genome[1]) == 5 {
				send_to_char(ch, libc.CString("  Truffle DNA.\r\n"))
			}
			if (ch.Genome[1]) == 6 {
				send_to_char(ch, libc.CString("  Arlian DNA.\r\n"))
			}
			if (ch.Genome[1]) == 7 {
				send_to_char(ch, libc.CString("  Kai DNA.\r\n"))
			}
			if (ch.Genome[1]) == 8 {
				send_to_char(ch, libc.CString("  Konatsu DNA.\r\n"))
			}
		}
		if (ch.Genome[0]) == 11 {
			send_to_char(ch, libc.CString("You have used kyodaika.\r\n"))
		}
		if PRF_FLAGGED(ch, PRF_NOPARRY) {
			send_to_char(ch, libc.CString("You have decided not to parry attacks.\r\n"))
		}
		switch ch.Position {
		case POS_DEAD:
			send_to_char(ch, libc.CString("You are DEAD!\r\n"))
		case POS_MORTALLYW:
			send_to_char(ch, libc.CString("You are mortally wounded! You should seek help!\r\n"))
		case POS_INCAP:
			send_to_char(ch, libc.CString("You are incapacitated, slowly fading away...\r\n"))
		case POS_STUNNED:
			send_to_char(ch, libc.CString("You are stunned! You can't move!\r\n"))
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You are sleeping.\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("You are resting.\r\n"))
		case POS_SITTING:
			send_to_char(ch, libc.CString("You are sitting.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("You are fighting %s.\r\n"), func() *byte {
				if ch.Fighting != nil {
					return PERS(ch.Fighting, ch)
				}
				return libc.CString("thin air")
			}())
		case POS_STANDING:
			send_to_char(ch, libc.CString("You are standing.\r\n"))
		default:
			send_to_char(ch, libc.CString("You are floating.\r\n"))
		}
		if has_group(ch) != 0 {
			send_to_char(ch, libc.CString("@GGroup Victories@D: @w%s@n\r\n"), add_commas(int64(ch.Combatexpertise)))
		}
		if PLR_FLAGGED(ch, PLR_EYEC) {
			send_to_char(ch, libc.CString("Your eyes are closed.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_SNEAK) {
			send_to_char(ch, libc.CString("You are prepared to sneak where ever you go.\r\n"))
		}
		if PLR_FLAGGED(ch, PLR_DISGUISED) {
			send_to_char(ch, libc.CString("You have disguised your facial features.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_FLYING) {
			send_to_char(ch, libc.CString("You are flying.\r\n"))
		}
		if PLR_FLAGGED(ch, PLR_PILOTING) {
			send_to_char(ch, libc.CString("You are busy piloting a ship.\r\n"))
		}
		if ch.Powerattack > 0 {
			send_to_char(ch, libc.CString("You are playing @y'@Y%s@y'@n.\r\n"), song_types[ch.Powerattack])
		}
		if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
			send_to_char(ch, libc.CString("You are prepared to zanzoken.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_HASS) {
			send_to_char(ch, libc.CString("Your arms are moving fast.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_INFUSE) {
			send_to_char(ch, libc.CString("Your ki will be infused in your next physical attack.\r\n"))
		}
		if PLR_FLAGGED(ch, PLR_TAILHIDE) {
			send_to_char(ch, libc.CString("Your tail is hidden!\r\n"))
		}
		if PLR_FLAGGED(ch, PLR_NOGROW) {
			send_to_char(ch, libc.CString("Your tail is no longer regrowing!\r\n"))
		}
		if PLR_FLAGGED(ch, PLR_POSE) {
			send_to_char(ch, libc.CString("You are feeling confident from your pose earlier.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_HYDROZAP) {
			send_to_char(ch, libc.CString("You are effected by Kanso Suru.\r\n"))
		}
		if int(ch.Player_specials.Conditions[DRUNK]) > 15 {
			send_to_char(ch, libc.CString("You are extremely drunk.\r\n"))
		} else if int(ch.Player_specials.Conditions[DRUNK]) > 10 {
			send_to_char(ch, libc.CString("You are pretty drunk.\r\n"))
		} else if int(ch.Player_specials.Conditions[DRUNK]) > 4 {
			send_to_char(ch, libc.CString("You are drunk.\r\n"))
		} else if int(ch.Player_specials.Conditions[DRUNK]) > 0 {
			send_to_char(ch, libc.CString("You have an alcoholic buzz.\r\n"))
		}
		if ch.Affected != nil {
			var lasttype int = 0
			for aff = ch.Affected; aff != nil; aff = aff.Next {
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("runic")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Kenaz rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("punch")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Algiz rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("knee")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Oagaz rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("slam")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Wunjo rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("heeldrop")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Purisaz rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("special beam cannon")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Laguz rune is still in effect! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("might")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your muscles are pumped! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("flex")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You are more agile right now! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("bless")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have been blessed! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("curse")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have been cursed! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("healing glow")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have a healing glow enveloping your body! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("genius")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You are smarter right now! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("enlighten")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You are wiser right now! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("yoikominminken")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have been lulled to sleep! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("solar flare")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have been blinded! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("spirit control")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have full control of your spirit! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("!UNUSED!")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You feel poison burning through your blood! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("tough skin")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have toughened skin right now! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("poison")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("You have been poisoned! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("warp pool")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Weakened State! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("dark metamorphosis")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your Dark Metamorphosis is still in effect. (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
				if libc.StrCaseCmp(skill_name(int(aff.Type)), libc.CString("hayasa")) == 0 && int(aff.Type) != lasttype {
					lasttype = int(aff.Type)
					send_to_char(ch, libc.CString("Your body has been infused to move faster! (%2d Mud Hours)\r\n"), int(aff.Duration)+1)
				}
			}
		}
		if AFF_FLAGGED(ch, AFF_KNOCKED) {
			send_to_char(ch, libc.CString("You have been knocked unconcious!\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_INVISIBLE) {
			send_to_char(ch, libc.CString("You are invisible.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_DETECT_INVIS) {
			send_to_char(ch, libc.CString("You are sensitive to the presence of invisible things.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_MBREAK) {
			send_to_char(ch, libc.CString("Your mind has been broken!\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_WITHER) {
			send_to_char(ch, libc.CString("You've been withered! You feel so weak...\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_SHOCKED) {
			send_to_char(ch, libc.CString("Your mind has been shocked!\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_CHARM) {
			send_to_char(ch, libc.CString("You have been charmed!\r\n"))
		}
		if affected_by_spell(ch, SPELL_MAGE_ARMOR) {
			send_to_char(ch, libc.CString("You feel protected.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_INFRAVISION) {
			send_to_char(ch, libc.CString("You can see in darkness with infravision.\r\n"))
		}
		if PRF_FLAGGED(ch, PRF_SUMMONABLE) {
			send_to_char(ch, libc.CString("You are summonable by other players.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_DETECT_ALIGN) {
			send_to_char(ch, libc.CString("You see into the hearts of others.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_DETECT_MAGIC) {
			send_to_char(ch, libc.CString("You are sensitive to the magical nature of things.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_SPIRIT) {
			send_to_char(ch, libc.CString("You have died and are part of the SPIRIT world!\r\n"))
		}
		if PRF_FLAGGED(ch, PRF_NOGIVE) {
			send_to_char(ch, libc.CString("You are not accepting items being handed to you right now.\r\n"))
		}
		if AFF_FLAGGED(ch, AFF_ETHEREAL) {
			send_to_char(ch, libc.CString("You are ethereal and cannot interact with normal space!\r\n"))
		}
		if ch.Regen > 0 {
			send_to_char(ch, libc.CString("Something is augmenting your regen rate by %s%d%s!\r\n"), func() string {
				if ch.Regen > 0 {
					return "+"
				}
				return "-"
			}(), ch.Regen, "%")
		}
		if ch.Asb > 0 {
			send_to_char(ch, libc.CString("Something is augmenting your auto-skill training rate by %s%d%s!\r\n"), func() string {
				if ch.Asb > 0 {
					return "+"
				}
				return "-"
			}(), ch.Asb, "%")
		}
		if ch.Lifebonus > 0 {
			send_to_char(ch, libc.CString("Something is augmenting your Life Force Max by %s%d%s!\r\n"), func() string {
				if ch.Lifebonus > 0 {
					return "+"
				}
				return "-"
			}(), ch.Lifebonus, "%")
		}
		if PLR_FLAGGED(ch, PLR_FISHING) {
			send_to_char(ch, libc.CString("Current Fishing Pole Bonus @D[@C%d@D]@n\r\n"), ch.Accuracy)
		}
		if PLR_FLAGGED(ch, PLR_AURALIGHT) {
			send_to_char(ch, libc.CString("Aura Light is active.\r\n"))
		}
		send_to_char(ch, libc.CString("@D<@b--------------------------------------------------------------@D>@n\r\n"))
		send_to_char(ch, libc.CString("To view your bonus/negative traits enter: status traits\r\n"))
	} else if libc.StrCaseCmp(&arg[0], libc.CString("traits")) == 0 {
		bonus_status(ch)
	} else {
		send_to_char(ch, libc.CString("The only argument status takes is 'traits'. If you just want your status do not use an argument.\r\n"))
	}
}

var list_bonuses [52]*byte = [52]*byte{libc.CString("Thrifty     - -10% Shop Buy Cost and +10% Shop Sell Cost          "), libc.CString("Prodigy     - +25% Experience Gained Until Level 80               "), libc.CString("Quick Study - Character auto-trains skills faster                 "), libc.CString("Die Hard    - Life Force's PL regen doubled, but cost is the same "), libc.CString("Brawler     - Physical attacks do 20% more damage                 "), libc.CString("Destroyer   - Damaged Rooms act as regen rooms for you            "), libc.CString("Hard Worker - Physical activity bonuses + drains less stamina     "), libc.CString("Healer      - Heal/First-aid/Vigor/Repair restore +10%            "), libc.CString("Loyal       - +20% Experience When Grouped As Follower            "), libc.CString("Brawny      - Strength gains +2 every 10 levels, Train STR + 75%  "), libc.CString("Scholarly   - Intelligence gains +2 every 10 levels, Train INT + 75%"), libc.CString("Sage        - Wisdom gains +2 every 10 levels, Train WIS + 75%    "), libc.CString("Agile       - Agility gains +2 every 10 levels, Train AGL + 75%   "), libc.CString("Quick       - Speed gains +2 every 10 levels, Train SPD + 75%     "), libc.CString("Sturdy      - Constitution +2 every 10 levels, Train CON + 75%    "), libc.CString("Thick Skin  - -20% Physical and -10% ki dmg received              "), libc.CString("Recipe Int. - Food cooked by you lasts longer/heals better        "), libc.CString("Fireproof   - -50% Fire Dmg taken, -10% ki, immunity to burn      "), libc.CString("Powerhitter - 15% critical hits will be x4 instead of x2          "), libc.CString("Healthy     - 40% chance to recover from ill effects when sleeping"), libc.CString("Insomniac   - Can't Sleep. Immune to yoikominminken and paralysis "), libc.CString("Evasive     - +15% to dodge rolls                                 "), libc.CString("The Wall    - +20% chance to block                                "), libc.CString("Accurate    - +20% chance to hit physical, +10% to hit with ki     "), libc.CString("Energy Leech- -2% ki damage received for every 5 character levels,\n                  @cas long as you can take that ki to your charge pool.@D        "), libc.CString("Good Memory - +2 Skill Slots initially, +1 every 20 levels after  "), libc.CString("Soft Touch  - Half damage for all hit locations                   "), libc.CString("Late Sleeper- Can only wake automatically. 33% every hour if maxed"), libc.CString("Impulse Shop- +25% shop costs                                     "), libc.CString("Sickly      - Suffer from harmful effects longer                  "), libc.CString("Punching Bag- -15% to dodge rolls                                 "), libc.CString("Pushover    - -20% block chance                                   "), libc.CString("Poor D. Perc- -20% chance to hit with physical, -10% with ki       "), libc.CString("Thin Skin   - +20% physical and +10% ki damage received           "), libc.CString("Fireprone   - +50% Fire Dmg taken, +10% ki, always burned         "), libc.CString("Energy Int. - +2% ki damage received for every 5 character levels,\n                  @rif you have ki charged you have 10% chance to lose   \n                  it and to take 1/4th damage equal to it.@D                    "), libc.CString("Coward      - Can't Attack Enemy With 150% Your Powerlevel        "), libc.CString("Arrogant    - Cannot Suppress                                     "), libc.CString("Unfocused   - Charge concentration randomly breaks                "), libc.CString("Slacker     - Physical activity drains more stamina               "), libc.CString("Slow Learner- Character auto-trains skills slower                 "), libc.CString("Masochistic - Defense Skills Cap At 75                            "), libc.CString("Mute        - Can't use IC speech related commands                "), libc.CString("Wimp        - Strength is capped at 45                            "), libc.CString("Dull        - Intelligence is capped at 45                        "), libc.CString("Foolish     - Wisdom is capped at 45                              "), libc.CString("Clumsy      - Agility is capped at 45                             "), libc.CString("Slow        - Speed is capped at 45                               "), libc.CString("Frail       - Constitution capped at 45                           "), libc.CString("Sadistic    - Half Experience Gained For Quick Kills              "), libc.CString("Loner       - Can't Group with anyone, +5% train and +10% Phys    "), libc.CString("Bad Memory  - -5 Skill Slots                                      ")}

func bonus_status(ch *char_data) {
	var (
		i     int
		max   int = 52
		count int = 0
	)
	if IS_NPC(ch) {
		return
	}
	send_to_char(ch, libc.CString("@CYour Traits@n\n@D-----------------------------@w\n"))
	for i = 0; i < max; i++ {
		if i < 26 {
			if (ch.Bonuses[i]) != 0 {
				send_to_char(ch, libc.CString("@c%s@n\n"), list_bonuses[i])
				count++
			}
		} else {
			if i == 26 {
				send_to_char(ch, libc.CString("\r\n"))
			}
			if (ch.Bonuses[i]) != 0 {
				send_to_char(ch, libc.CString("@r%s@n\n"), list_bonuses[i])
				count++
			}
		}
	}
	if count <= 0 {
		send_to_char(ch, libc.CString("@wNone.\r\n"))
	}
	send_to_char(ch, libc.CString("@D-----------------------------@n\r\n"))
	return
}
func do_inventory(ch *char_data, argument *byte, cmd int, subcmd int) {
	send_to_char(ch, libc.CString("@w              @YInventory\r\n@D-------------------------------------@w\r\n"))
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_STOLEN) {
			ch.Act[int(PLR_STOLEN/32)] &= bitvector_t(int32(^(1 << (int(PLR_STOLEN % 32)))))
			send_to_char(ch, libc.CString("@r   --------------------------------------------------@n\n"))
			send_to_char(ch, libc.CString("@R    You notice that you have been robbed sometime recently!\n"))
			send_to_char(ch, libc.CString("@r   --------------------------------------------------@n\n"))
			return
		}
	}
	list_obj_to_char(ch.Carrying, ch, SHOW_OBJ_SHORT, TRUE)
	send_to_char(ch, libc.CString("\n"))
}
func do_equipment(ch *char_data, argument *byte, cmd int, subcmd int) {
	var i int
	send_to_char(ch, libc.CString("        @YEquipment Being Worn\r\n@D-------------------------------------@w\r\n"))
	for i = 1; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			if CAN_SEE_OBJ(ch, ch.Equipment[i]) && (i != WEAR_WIELD1 && i != WEAR_WIELD2) {
				send_to_char(ch, libc.CString("%s"), wear_where[i])
				show_obj_to_char(ch.Equipment[i], ch, SHOW_OBJ_SHORT)
				if OBJ_FLAGGED(ch.Equipment[i], ITEM_SHEATH) {
					var (
						obj2     *obj_data = nil
						next_obj *obj_data = nil
						sheath   *obj_data = (ch.Equipment[i])
					)
					for obj2 = sheath.Contains; obj2 != nil; obj2 = next_obj {
						next_obj = obj2.Next_content
						if obj2 != nil {
							send_to_char(ch, libc.CString("@D  ---- @YSheathed@D ----@c> @n"))
							show_obj_to_char(obj2, ch, SHOW_OBJ_SHORT)
						}
					}
					obj2 = nil
				}
			} else if CAN_SEE_OBJ(ch, ch.Equipment[i]) && !PLR_FLAGGED(ch, PLR_THANDW) {
				send_to_char(ch, libc.CString("%s"), wear_where[i])
				show_obj_to_char(ch.Equipment[i], ch, SHOW_OBJ_SHORT)
				if OBJ_FLAGGED(ch.Equipment[i], ITEM_SHEATH) {
					var (
						obj2     *obj_data = nil
						next_obj *obj_data = nil
						sheath   *obj_data = (ch.Equipment[i])
					)
					for obj2 = sheath.Contains; obj2 != nil; obj2 = next_obj {
						next_obj = obj2.Next_content
						if obj2 != nil {
							send_to_char(ch, libc.CString("@D  ---- @YSheathed@D ----> @n"))
							show_obj_to_char(obj2, ch, SHOW_OBJ_SHORT)
						}
					}
					obj2 = nil
				}
			} else if CAN_SEE_OBJ(ch, ch.Equipment[i]) && PLR_FLAGGED(ch, PLR_THANDW) {
				send_to_char(ch, libc.CString("@c<@CWielded by B. Hands@c>@n "))
				show_obj_to_char(ch.Equipment[i], ch, SHOW_OBJ_SHORT)
			} else {
				send_to_char(ch, libc.CString("%s"), wear_where[i])
				send_to_char(ch, libc.CString("Something.\r\n"))
			}
		} else {
			if BODY_FLAGGED(ch, bitvector_t(int32(i))) && i != WEAR_WIELD2 {
				send_to_char(ch, libc.CString("%s@wNothing.@n\r\n"), wear_where[i])
			} else if BODY_FLAGGED(ch, bitvector_t(int32(i))) && (i == WEAR_WIELD2 && !PLR_FLAGGED(ch, PLR_THANDW)) {
				send_to_char(ch, libc.CString("%s@wNothing.@n\r\n"), wear_where[i])
			}
		}
	}
}
func do_time(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		suf     *byte
		weekday int
		day     int
	)
	day = time_info.Day + 1
	weekday = day % 6
	send_to_char(ch, libc.CString("It is %d o'clock %s, on %s.\r\n"), func() int {
		if time_info.Hours%12 == 0 {
			return 12
		}
		return time_info.Hours % 12
	}(), func() string {
		if time_info.Hours >= 12 {
			return "PM"
		}
		return "AM"
	}(), weekdays[weekday])
	suf = libc.CString("th")
	if ((day % 100) / 10) != 1 {
		switch day % 10 {
		case 1:
			suf = libc.CString("st")
		case 2:
			suf = libc.CString("nd")
		case 3:
			suf = libc.CString("rd")
		}
	}
	send_to_char(ch, libc.CString("The %d%s Day of the %s, Year %d.\r\n"), day, suf, month_name[time_info.Month], time_info.Year)
}
func do_weather(ch *char_data, argument *byte, cmd int, subcmd int) {
	var sky_look [4]*byte = [4]*byte{libc.CString("cloudless"), libc.CString("cloudy"), libc.CString("rainy"), libc.CString("lit by flashes of lightning")}
	if OUTSIDE(ch) {
		send_to_char(ch, libc.CString("The sky is %s and %s.\r\n"), sky_look[weather_info.Sky], func() string {
			if weather_info.Change >= 0 {
				return "you feel a warm wind from south"
			}
			return "your foot tells you bad weather is due"
		}())
		if ADM_FLAGGED(ch, ADM_KNOWWEATHER) {
			send_to_char(ch, libc.CString("Pressure: %d (change: %d), Sky: %d (%s)\r\n"), weather_info.Pressure, weather_info.Change, weather_info.Sky, sky_look[weather_info.Sky])
		}
	} else {
		send_to_char(ch, libc.CString("You have no feeling about the weather at all.\r\n"))
	}
}
func space_to_minus(str *byte) {
	for (func() *byte {
		str = libc.StrChr(str, ' ')
		return str
	}()) != nil {
		*str = '-'
	}
}
func search_help(argument *byte, level int) int {
	var (
		chk    int
		bot    int
		top    int
		mid    int
		minlen int
	)
	bot = 0
	top = top_of_helpt
	minlen = libc.StrLen(argument)
	for bot <= top {
		mid = (bot + top) / 2
		if (func() int {
			chk = libc.StrNCaseCmp(argument, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Keywords, minlen)
			return chk
		}()) == 0 {
			for mid > 0 && libc.StrNCaseCmp(argument, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid-1)))).Keywords, minlen) == 0 {
				mid--
			}
			for level < (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Min_level && mid < (bot+top)/2 {
				mid++
			}
			if libc.StrNCaseCmp(argument, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Keywords, minlen) != 0 {
				break
			}
			return mid
		} else if chk > 0 {
			bot = mid + 1
		} else {
			top = mid - 1
		}
	}
	return -1
}
func do_help(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf [259744]byte
		mid int = 0
	)
	if ch.Desc == nil {
		return
	}
	skip_spaces(&argument)
	if help_table == nil {
		send_to_char(ch, libc.CString("No help available.\r\n"))
		return
	}
	if *argument == 0 {
		if ch.Admlevel < ADMLVL_IMMORT {
			page_string(ch.Desc, help, 0)
		} else {
			page_string(ch.Desc, ihelp, 0)
		}
		return
	}
	space_to_minus(argument)
	if (func() int {
		mid = search_help(argument, ch.Admlevel)
		return mid
	}()) == int(-1) {
		var (
			i     int
			found int = 0
		)
		send_to_char(ch, libc.CString("There is no help on that word.\r\n"))
		if ch.Admlevel < 3 {
			mudlog(NRM, MAX(ADMLVL_IMPL, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s tried to get help on %s"), GET_NAME(ch), argument)
		}
		for i = 0; i <= top_of_helpt; i++ {
			if (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Min_level > ch.Admlevel {
				continue
			}
			if *argument != *(*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Keywords {
				continue
			}
			if levenshtein_distance(argument, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Keywords) <= 2 {
				if found == 0 {
					send_to_char(ch, libc.CString("\r\nDid you mean:\r\n"))
					found = 1
				}
				send_to_char(ch, libc.CString("  %s\r\n"), (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Keywords)
			}
		}
		return
	}
	if (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Min_level > ch.Admlevel {
		send_to_char(ch, libc.CString("There is no help on that word.\r\n"))
		return
	}
	stdio.Sprintf(&buf[0], "@b~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s", (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Entry)
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@b~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\n")
	if ch.Admlevel > 0 {
		stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@WHelp File Level@w: @D(@R%d@D)@n\n", (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Min_level)
	}
	page_string(ch.Desc, &buf[0], 0)
}
func do_who(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d           *descriptor_data
		tch         *char_data
		i           int
		num_can_see int = 0
		name_search [2048]byte
		buf         [2048]byte
		low         int   = 0
		high        int   = config_info.Play.Level_cap
		localwho    int   = 0
		questwho    int   = 0
		hide        int   = 0
		showclass   int   = 0
		short_list  int   = 0
		outlaws     int   = 0
		who_room    int   = 0
		showgroup   int   = 0
		showleader  int   = 0
		line_color  *byte = libc.CString("@n")
	)
	skip_spaces(&argument)
	libc.StrCpy(&buf[0], argument)
	name_search[0] = '\x00'
	var rank [3]struct {
		Disp      *byte
		Min_level int
		Max_level int
		Count     int
	} = [3]struct {
		Disp      *byte
		Min_level int
		Max_level int
		Count     int
	}{{Disp: libc.CString("\r\n               @c------------  @D[    @gI@Gm@Wm@Do@Gr@Dt@Wa@Gl@gs   @D]  @c------------@n\r\n"), Min_level: ADMLVL_IMMORT, Max_level: ADMLVL_IMPL, Count: 0}, {Disp: libc.CString("\r\n@D[@wx@D]@yxxxxxxxxxx@W  [    @GImmortals   @W]  @yxxxxxxxxxx@D[@wx@D]@n\r\n"), Min_level: int(ADMLVL_IMMORT + 8), Max_level: int(ADMLVL_GRGOD + 8), Count: 0}, {Disp: libc.CString("\r\n               @c------------  @D[     @DM@ro@Rr@wt@Ra@rl@Ds    ]  @c------------@n\r\n"), Min_level: 0, Max_level: int(ADMLVL_IMMORT - 1), Count: 0}}
	var tmstr *byte
	tmstr = libc.AscTime(libc.LocalTime(&PCOUNTDATE))
	*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
	var num_ranks int = int(unsafe.Sizeof([3]struct {
		Disp      *byte
		Min_level int
		Max_level int
		Count     int
	}{}) / unsafe.Sizeof(struct {
		Disp      *byte
		Min_level int
		Max_level int
		Count     int
	}{}))
	send_to_char(ch, libc.CString("\r\n      @r{@b===============  @D[  @DD@wr@ca@Cg@Y(@R*@Y)@Wn@cB@Da@cl@Cl @DA@wd@cv@Ce@Wnt @DT@wr@cu@Ct@Wh@n  @D]  @b===============@r}      @n\r\n"))
	for d = descriptor_list; d != nil && short_list == 0; d = d.Next {
		if !IS_PLAYING(d) {
			continue
		}
		if d.Original != nil {
			tch = d.Original
		} else if (func() *char_data {
			tch = d.Character
			return tch
		}()) == nil {
			continue
		}
		if tch.Admlevel >= ADMLVL_IMMORT {
			line_color = libc.CString("@w")
		} else {
			line_color = libc.CString("@w")
		}
		if CAN_SEE(ch, tch) && IS_PLAYING(d) {
			if name_search[0] != 0 && libc.StrCaseCmp(GET_NAME(tch), &name_search[0]) != 0 && libc.StrStr(GET_TITLE(tch), &name_search[0]) == nil {
				continue
			}
			if !CAN_SEE(ch, tch) || GET_LEVEL(tch) < low || GET_LEVEL(tch) > high {
				continue
			}
			if outlaws != 0 && !PLR_FLAGGED(tch, PLR_KILLER) && !PLR_FLAGGED(tch, PLR_THIEF) {
				continue
			}
			if questwho != 0 && !PRF_FLAGGED(tch, PRF_QUEST) {
				continue
			}
			if localwho != 0 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(tch.In_room)))).Zone {
				continue
			}
			if PRF_FLAGGED(tch, PRF_HIDE) && tch != ch && ch.Admlevel < ADMLVL_IMMORT {
				hide += 1
				continue
			}
			if who_room != 0 && tch.In_room != ch.In_room {
				continue
			}
			if showclass != 0 && (showclass&(1<<int(tch.Chclass))) == 0 {
				continue
			}
			if showgroup != 0 && (tch.Master == nil || !AFF_FLAGGED(tch, AFF_GROUP)) {
				continue
			}
			for i = 0; i < num_ranks; i++ {
				if tch.Admlevel >= rank[i].Min_level && tch.Admlevel <= rank[i].Max_level {
					rank[i].Count++
				}
			}
		}
	}
	for i = 0; i < num_ranks; i++ {
		if rank[i].Count == 0 && short_list == 0 {
			continue
		}
		if short_list != 0 {
			send_to_char(ch, libc.CString("Players\r\n-------\r\n"))
		} else {
			send_to_char(ch, rank[i].Disp)
		}
		for d = descriptor_list; d != nil; d = d.Next {
			if !IS_PLAYING(d) {
				continue
			}
			if d.Original != nil {
				tch = d.Original
			} else if (func() *char_data {
				tch = d.Character
				return tch
			}()) == nil {
				continue
			}
			if (tch.Admlevel < rank[i].Min_level || tch.Admlevel > rank[i].Max_level) && short_list == 0 {
				continue
			}
			if !IS_PLAYING(d) {
				continue
			}
			if name_search[0] != 0 && libc.StrCaseCmp(GET_NAME(tch), &name_search[0]) != 0 && libc.StrStr(GET_TITLE(tch), &name_search[0]) == nil {
				continue
			}
			if !CAN_SEE(ch, tch) || GET_LEVEL(tch) < low || GET_LEVEL(tch) > high {
				continue
			}
			if outlaws != 0 && !PLR_FLAGGED(tch, PLR_KILLER) && !PLR_FLAGGED(tch, PLR_THIEF) {
				continue
			}
			if questwho != 0 && !PRF_FLAGGED(tch, PRF_QUEST) {
				continue
			}
			if localwho != 0 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(tch.In_room)))).Zone {
				continue
			}
			if who_room != 0 && tch.In_room != ch.In_room {
				continue
			}
			if PRF_FLAGGED(tch, PRF_HIDE) && tch != ch && ch.Admlevel < ADMLVL_IMMORT {
				continue
			}
			if showclass != 0 && (showclass&(1<<int(tch.Chclass))) == 0 {
				continue
			}
			if showgroup != 0 && (tch.Master == nil || !AFF_FLAGGED(tch, AFF_GROUP)) {
				continue
			}
			if showleader != 0 && (tch.Followers == nil || !AFF_FLAGGED(tch, AFF_GROUP)) {
				continue
			}
			if short_list != 0 {
				send_to_char(ch, libc.CString("               @B[@W%3d @Y%s @C%s@B]@W %-12.12s@n%s@n"), GET_LEVEL(tch), race_abbrevs[int(tch.Race)], class_abbrevs[int(tch.Chclass)], GET_NAME(tch), func() string {
					if (func() int {
						p := &num_can_see
						*p++
						return *p
					}() % 4) == 0 {
						return "\r\n"
					}
					return ""
				}())
			} else {
				num_can_see++
				var usr [100]byte
				stdio.Sprintf(&usr[0], "@W(@R%s@W)%s", tch.Desc.User, func() string {
					if PLR_FLAGGED(tch, PLR_BIOGR) {
						return ""
					}
					if SPOILED(tch) {
						return " @R*@n"
					}
					return ""
				}())
				send_to_char(ch, libc.CString("%s               @D<@C%-12s@D> %s@w%s"), line_color, func() *byte {
					if ch.Admlevel > 0 {
						return GET_NAME(tch)
					}
					if tch.Admlevel > 0 {
						return GET_NAME(tch)
					}
					if GET_USER(tch) != nil {
						return GET_USER(tch)
					}
					return libc.CString("NULL")
				}(), &func() [100]byte {
					if ch.Admlevel > 0 {
						return usr
					}
					return func() [100]byte {
						var t [100]byte
						copy(t[:], []byte(""))
						return t
					}()
				}()[0], line_color)
				if tch.Admlevel != 0 {
					send_to_char(ch, libc.CString(" (%s)"), admin_level_names[tch.Admlevel])
				}
				if d.Snooping != nil && d.Snooping.Character != ch && ch.Admlevel >= 3 {
					send_to_char(ch, libc.CString(" (Snoop: %s)"), GET_NAME(d.Snooping.Character))
				}
				if int(tch.Player_specials.Invis_level) != 0 {
					send_to_char(ch, libc.CString(" (i%d)"), tch.Player_specials.Invis_level)
				} else if AFF_FLAGGED(tch, AFF_INVISIBLE) {
					send_to_char(ch, libc.CString(" (invis)"))
				}
				if PLR_FLAGGED(tch, PLR_MAILING) {
					send_to_char(ch, libc.CString(" (mailing)"))
				} else if d.Olc != nil {
					send_to_char(ch, libc.CString(" (OLC)"))
				} else if PLR_FLAGGED(tch, PLR_WRITING) {
					send_to_char(ch, libc.CString(" (writing)"))
				}
				if d.Original != nil {
					send_to_char(ch, libc.CString(" (out of body)"))
				}
				if d.Connected == CON_OEDIT {
					send_to_char(ch, libc.CString(" (O Edit)"))
				}
				if d.Connected == CON_MEDIT {
					send_to_char(ch, libc.CString(" (M Edit)"))
				}
				if d.Connected == CON_ZEDIT {
					send_to_char(ch, libc.CString(" (Z Edit)"))
				}
				if d.Connected == CON_SEDIT {
					send_to_char(ch, libc.CString(" (S Edit)"))
				}
				if d.Connected == CON_REDIT {
					send_to_char(ch, libc.CString(" (R Edit)"))
				}
				if d.Connected == CON_TEDIT {
					send_to_char(ch, libc.CString(" (T Edit)"))
				}
				if d.Connected == CON_TRIGEDIT {
					send_to_char(ch, libc.CString(" (T Edit)"))
				}
				if d.Connected == CON_AEDIT {
					send_to_char(ch, libc.CString(" (S Edit)"))
				}
				if d.Connected == CON_CEDIT {
					send_to_char(ch, libc.CString(" (C Edit)"))
				}
				if d.Connected == CON_HEDIT {
					send_to_char(ch, libc.CString(" (H Edit)"))
				}
				if PRF_FLAGGED(tch, PRF_DEAF) {
					send_to_char(ch, libc.CString(" (DEAF)"))
				}
				if PRF_FLAGGED(tch, PRF_NOTELL) {
					send_to_char(ch, libc.CString(" (NO TELL)"))
				}
				if PRF_FLAGGED(tch, PRF_NOGOSS) {
					send_to_char(ch, libc.CString(" (NO OOC)"))
				}
				if PLR_FLAGGED(tch, PLR_NOSHOUT) {
					send_to_char(ch, libc.CString(" (MUTED)"))
				}
				if PRF_FLAGGED(tch, PRF_HIDE) {
					send_to_char(ch, libc.CString(" (WH)"))
				}
				if PRF_FLAGGED(tch, PRF_BUILDWALK) {
					send_to_char(ch, libc.CString(" (Buildwalking)"))
				}
				if PRF_FLAGGED(tch, PRF_AFK) {
					send_to_char(ch, libc.CString(" (AFK)"))
				}
				if PLR_FLAGGED(tch, PLR_FISHING) && ch.Admlevel >= ADMLVL_IMMORT {
					send_to_char(ch, libc.CString(" (@BFISHING@n)"))
				}
				if PRF_FLAGGED(tch, PRF_NOWIZ) {
					send_to_char(ch, libc.CString(" (NO WIZ)"))
				}
				send_to_char(ch, libc.CString("@n\r\n"))
			}
		}
		send_to_char(ch, libc.CString("\r\n"))
		if short_list != 0 {
			break
		}
	}
	if num_can_see == 0 {
		send_to_char(ch, libc.CString("                            Nobody at all!\r\n"))
	} else if num_can_see == 1 {
		send_to_char(ch, libc.CString("                         One lonely character displayed.\r\n"))
	} else {
		send_to_char(ch, libc.CString("                           @Y%d@w characters displayed.\r\n"), num_can_see)
		if hide > 0 {
			var bam int = FALSE
			if hide > 1 {
				bam = TRUE
			}
			send_to_char(ch, libc.CString("                           and @Y%d@w character%s hidden.\r\n"), hide, func() string {
				if bam != 0 {
					return "s"
				}
				return ""
			}())
		}
	}
	if circle_restrict > 0 && circle_restrict <= 100 {
		send_to_char(ch, libc.CString("                      @rThe mud has been wizlocked to lvl %d@n\r\n"), circle_restrict)
	}
	if circle_restrict == 101 {
		send_to_char(ch, libc.CString("                      @rThe mud has been wizlocked to IMMs only.@n\r\n"))
	}
	send_to_char(ch, libc.CString("      @r{@b=================================================================@r}@n\r\n"))
	send_to_char(ch, libc.CString("           @cHighest Logon Count Ever@D: @Y%d@w, on %s\r\n"), HIGHPCOUNT, tmstr)
	send_to_char(ch, libc.CString("                        @cHighest Logon Count Today@D: @Y%d@n\r\n"), PCOUNT)
}
func do_users(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		line        [200]byte
		line2       [220]byte
		idletime    [10]byte
		state       [30]byte
		timeptr     *byte
		mode        int8
		name_search [2048]byte
		host_search [2048]byte
		tch         *char_data
		d           *descriptor_data
		low         int = 0
		high        int = config_info.Play.Level_cap
		num_can_see int = 0
		showclass   int = 0
		outlaws     int = 0
		playing     int = 0
		deadweight  int = 0
		showrace    int = 0
		buf         [2048]byte
		arg         [2048]byte
	)
	host_search[0] = func() byte {
		p := &name_search[0]
		name_search[0] = '\x00'
		return *p
	}()
	libc.StrCpy(&buf[0], argument)
	for buf[0] != 0 {
		var buf1 [2048]byte
		half_chop(&buf[0], &arg[0], &buf1[0])
		if arg[0] == '-' {
			mode = int8(arg[1])
			switch mode {
			case 'o':
				fallthrough
			case 'k':
				outlaws = 1
				playing = 1
				libc.StrCpy(&buf[0], &buf1[0])
			case 'p':
				playing = 1
				libc.StrCpy(&buf[0], &buf1[0])
			case 'd':
				deadweight = 1
				libc.StrCpy(&buf[0], &buf1[0])
			case 'l':
				playing = 1
				half_chop(&buf1[0], &arg[0], &buf[0])
				stdio.Sscanf(&arg[0], "%d-%d", &low, &high)
			case 'n':
				playing = 1
				half_chop(&buf1[0], &name_search[0], &buf[0])
			case 'h':
				playing = 1
				half_chop(&buf1[0], &host_search[0], &buf[0])
			default:
				send_to_char(ch, libc.CString("%s"), USERS_FORMAT)
				return
			}
		} else {
			send_to_char(ch, libc.CString("%s"), USERS_FORMAT)
			return
		}
	}
	send_to_char(ch, libc.CString("Num Name                 User-name            State          Idl Login    C\r\n--- -------------------- -------------------- -------------- --- -------- -\r\n"))
	one_argument(argument, &arg[0])
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected != CON_PLAYING && playing != 0 {
			continue
		}
		if d.Connected == CON_PLAYING && deadweight != 0 {
			continue
		}
		if IS_PLAYING(d) {
			if d.Original != nil {
				tch = d.Original
			} else if (func() *char_data {
				tch = d.Character
				return tch
			}()) == nil {
				continue
			}
			if host_search[0] != 0 && libc.StrStr(&d.Host[0], &host_search[0]) == nil {
				continue
			}
			if name_search[0] != 0 && libc.StrCaseCmp(GET_NAME(tch), &name_search[0]) != 0 {
				continue
			}
			if !CAN_SEE(ch, tch) || GET_LEVEL(tch) < low || GET_LEVEL(tch) > high {
				continue
			}
			if PRF_FLAGGED(tch, PRF_HIDE) && tch != ch && ch.Admlevel < ADMLVL_IMMORT {
				continue
			}
			if outlaws != 0 && !PLR_FLAGGED(tch, PLR_KILLER) && !PLR_FLAGGED(tch, PLR_THIEF) {
				continue
			}
			if showclass != 0 && (showclass&(1<<int(tch.Chclass))) == 0 {
				continue
			}
			if showrace != 0 && (showrace&(1<<int(tch.Race))) == 0 {
				continue
			}
			if int(tch.Player_specials.Invis_level) > ch.Admlevel {
				continue
			}
		}
		timeptr = libc.AscTime(libc.LocalTime(&d.Login_time))
		timeptr = (*byte)(unsafe.Add(unsafe.Pointer(timeptr), 11))
		*((*byte)(unsafe.Add(unsafe.Pointer(timeptr), 8))) = '\x00'
		if d.Connected == CON_PLAYING && d.Original != nil {
			libc.StrCpy(&state[0], libc.CString("Switched"))
		} else {
			libc.StrCpy(&state[0], connected_types[d.Connected])
		}
		if d.Character != nil && d.Connected == CON_PLAYING && d.Character.Admlevel <= ch.Admlevel {
			stdio.Sprintf(&idletime[0], "%3d", d.Character.Timer*SECS_PER_MUD_HOUR/SECS_PER_REAL_MIN)
		} else {
			libc.StrCpy(&idletime[0], libc.CString(""))
		}
		stdio.Sprintf(&line[0], "%3d %-20s %-20s %-14s %-3s %-8s %1s ", d.Desc_num, func() *byte {
			if d.Original != nil && d.Original.Name != nil {
				return d.Original.Name
			}
			if d.Character != nil && d.Character.Name != nil {
				return d.Character.Name
			}
			return libc.CString("UNDEFINED")
		}(), func() *byte {
			if d.User != nil {
				return d.User
			}
			return libc.CString("UNKNOWN")
		}(), &state[0], &idletime[0], timeptr, "N")
		if d.Host != nil && d.Host[0] != 0 {
			stdio.Sprintf(&line[libc.StrLen(&line[0])], "\n%3d [%s Site: %s]\r\n", d.Desc_num, func() *byte {
				if d.User != nil {
					return d.User
				}
				return libc.CString("UNKNOWN")
			}(), &d.Host[0])
		} else {
			stdio.Sprintf(&line[libc.StrLen(&line[0])], "\n%3d [%s Site: Hostname unknown]\r\n", d.Desc_num, func() *byte {
				if d.User != nil {
					return d.User
				}
				return libc.CString("UNKNOWN")
			}())
		}
		if d.Connected != CON_PLAYING {
			stdio.Sprintf(&line2[0], "@g%s@n", &line[0])
			libc.StrCpy(&line[0], &line2[0])
		}
		if d.Connected != CON_PLAYING || d.Connected == CON_PLAYING && CAN_SEE(ch, d.Character) {
			send_to_char(ch, libc.CString("%s"), &line[0])
			num_can_see++
		}
	}
	send_to_char(ch, libc.CString("\r\n%d visible sockets connected.\r\n"), num_can_see)
}
func do_gen_ps(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		bum [10000]byte
	)
	one_argument(argument, &arg[0])
	switch subcmd {
	case SCMD_CREDITS:
		page_string(ch.Desc, credits, 0)
	case SCMD_NEWS:
		page_string(ch.Desc, news, 0)
		ch.Lastpl = libc.GetTime(nil)
	case SCMD_INFO:
		page_string(ch.Desc, info, 0)
	case SCMD_WIZLIST:
		page_string(ch.Desc, wizlist, 0)
	case SCMD_IMMLIST:
		page_string(ch.Desc, immlist, 0)
	case SCMD_HANDBOOK:
		page_string(ch.Desc, handbook, 0)
	case SCMD_POLICIES:
		stdio.Sprintf(&bum[0], "--------------------\r\n%s\r\n--------------------\r\n", policies)
		page_string(ch.Desc, &bum[0], 0)
	case SCMD_MOTD:
		page_string(ch.Desc, motd, 0)
	case SCMD_IMOTD:
		page_string(ch.Desc, imotd, 0)
	case SCMD_CLEAR:
		send_to_char(ch, libc.CString("\x1b[H\x1b[J"))
	case SCMD_VERSION:
		send_to_char(ch, libc.CString("%s\r\n"), circlemud_version)
		send_to_char(ch, libc.CString("%s\r\n"), oasisolc_version)
		send_to_char(ch, libc.CString("%s\r\n"), DG_SCRIPT_VERSION)
		send_to_char(ch, libc.CString("%s\r\n"), CWG_VERSION)
		send_to_char(ch, libc.CString("%s\r\n"), DBAT_VERSION)
	case SCMD_WHOAMI:
		send_to_char(ch, libc.CString("%s\r\n"), GET_NAME(ch))
	default:
		basic_mud_log(libc.CString("SYSERR: Unhandled case in do_gen_ps. (%d)"), subcmd)
		return
	}
}
func perform_mortal_where(ch *char_data, arg *byte) {
	var (
		i *char_data
		d *descriptor_data
	)
	if *arg == 0 {
		send_to_char(ch, libc.CString("Players in your Zone\r\n--------------------\r\n"))
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Connected != CON_PLAYING || d.Character == ch {
				continue
			}
			if (func() *char_data {
				i = func() *char_data {
					if d.Original != nil {
						return d.Original
					}
					return d.Character
				}()
				return i
			}()) == nil {
				continue
			}
			if i.In_room == room_rnum(-1) || !CAN_SEE(ch, i) {
				continue
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Zone {
				continue
			}
			send_to_char(ch, libc.CString("%-20s - %s\r\n"), GET_NAME(i), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Name)
		}
	} else {
		for i = character_list; i != nil; i = i.Next {
			if i.In_room == room_rnum(-1) || i == ch {
				continue
			}
			if !CAN_SEE(ch, i) || (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone {
				continue
			}
			if isname(arg, i.Name) == 0 {
				continue
			}
			send_to_char(ch, libc.CString("%-25s - %s\r\n"), GET_NAME(i), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Name)
			return
		}
		send_to_char(ch, libc.CString("Nobody around by that name.\r\n"))
	}
}
func print_object_location(num int, obj *obj_data, ch *char_data, recur int) {
	if num > 0 {
		send_to_char(ch, libc.CString("O%3d. %-25s - "), num, obj.Short_description)
	} else {
		send_to_char(ch, libc.CString("%33s"), " - ")
	}
	if obj.Script != nil {
		send_to_char(ch, libc.CString("[T%d]"), obj.Proto_script.Vnum)
	}
	if obj.In_room != room_rnum(-1) {
		send_to_char(ch, libc.CString("[%5d] %s\r\n"), func() room_vnum {
			if obj.In_room != room_rnum(-1) && obj.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).Number
			}
			return -1
		}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).Name)
	} else if obj.Carried_by != nil {
		send_to_char(ch, libc.CString("carried by %s in room [%d]\r\n"), PERS(obj.Carried_by, ch), func() room_vnum {
			if obj.Carried_by.In_room != room_rnum(-1) && obj.Carried_by.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.Carried_by.In_room)))).Number
			}
			return -1
		}())
	} else if obj.Worn_by != nil {
		send_to_char(ch, libc.CString("worn by %s in room [%d]\r\n"), PERS(obj.Worn_by, ch), func() room_vnum {
			if obj.Worn_by.In_room != room_rnum(-1) && obj.Worn_by.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.Worn_by.In_room)))).Number
			}
			return -1
		}())
	} else if obj.In_obj != nil {
		send_to_char(ch, libc.CString("inside %s%s\r\n"), obj.In_obj.Short_description, func() string {
			if recur != 0 {
				return ", which is"
			}
			return " "
		}())
		if recur != 0 {
			print_object_location(0, obj.In_obj, ch, recur)
		}
	} else {
		send_to_char(ch, libc.CString("in an unknown location\r\n"))
	}
}
func perform_immort_where(ch *char_data, arg *byte) {
	var (
		i      *char_data
		k      *obj_data
		d      *descriptor_data
		num    int       = 0
		num2   int       = 0
		found  int       = 0
		planet [11]*byte = [11]*byte{libc.CString("@GEarth@n"), libc.CString("@CFrigid@n"), libc.CString("@YVegeta@n"), libc.CString("@MKonack@n"), libc.CString("@gNamek@n"), libc.CString("@mAether@n"), libc.CString("@mArlia@n"), libc.CString("@CZenith@n"), libc.CString("@YYardrat@n"), libc.CString("@cKanassa@n"), libc.CString("@RUNKOWN@n")}
	)
	if *arg == 0 {
		mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("GODCMD: %s has checked where to check player locations"), GET_NAME(ch))
		send_to_char(ch, libc.CString("Players                  Vnum    Planet        Location\r\n-------                 ------   ----------    ----------------\r\n"))
		for d = descriptor_list; d != nil; d = d.Next {
			if IS_PLAYING(d) {
				if d.Character.In_room != room_rnum(-1) {
					if ROOM_FLAGGED(d.Character.In_room, ROOM_EARTH) {
						num2 = 0
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_FRIGID) {
						num2 = 1
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_VEGETA) {
						num2 = 2
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_KONACK) {
						num2 = 3
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_NAMEK) {
						num2 = 4
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_AETHER) {
						num2 = 5
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_ARLIA) {
						num2 = 6
					} else if (func() room_vnum {
						if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
						}
						return -1
					}()) >= 3400 && (func() room_vnum {
						if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
						}
						return -1
					}()) <= 3599 || (func() room_vnum {
						if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
						}
						return -1
					}()) >= 62900 && (func() room_vnum {
						if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
						}
						return -1
					}()) <= 0xF617 || (func() room_vnum {
						if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
						}
						return -1
					}()) == 19600 {
						num2 = 7
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_YARDRAT) {
						num2 = 8
					} else if ROOM_FLAGGED(d.Character.In_room, ROOM_KANASSA) {
						num2 = 9
					} else {
						num2 = 10
					}
				}
				if d.Original != nil {
					i = d.Original
				} else {
					i = d.Character
				}
				if i != nil && CAN_SEE(ch, i) && i.In_room != room_rnum(-1) {
					if d.Original != nil {
						send_to_char(ch, libc.CString("%-20s - [%5d]   %s (in %s)\r\n"), GET_NAME(i), func() room_vnum {
							if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number
							}
							return -1
						}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Name, GET_NAME(d.Character))
					} else {
						send_to_char(ch, libc.CString("%-20s - [%5d]   %-14s %s\r\n"), GET_NAME(i), func() room_vnum {
							if i.In_room != room_rnum(-1) && i.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Number
							}
							return -1
						}(), planet[num2], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Name)
					}
				}
			}
		}
	} else {
		mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("GODCMD: %s has checked where for the location of %s"), GET_NAME(ch), arg)
		for i = character_list; i != nil; i = i.Next {
			if CAN_SEE(ch, i) && i.In_room != room_rnum(-1) && isname(arg, i.Name) != 0 {
				found = 1
				send_to_char(ch, libc.CString("M%3d. %-25s - [%5d] %-25s"), func() int {
					p := &num
					*p++
					return *p
				}(), GET_NAME(i), func() room_vnum {
					if i.In_room != room_rnum(-1) && i.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Number
					}
					return -1
				}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Name)
				if IS_NPC(i) && i.Script != nil {
					if i.Script.Trig_list.Next == nil {
						send_to_char(ch, libc.CString("[T%5d] "), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i.Script.Trig_list.Nr)))).Vnum)
					} else {
						send_to_char(ch, libc.CString("[TRIGS] "))
					}
				}
				send_to_char(ch, libc.CString("\r\n"))
			}
		}
		for k = object_list; k != nil; k = k.Next {
			if CAN_SEE_OBJ(ch, k) && isname(arg, k.Name) != 0 {
				found = 1
				print_object_location(func() int {
					p := &num
					*p++
					return *p
				}(), k, ch, TRUE)
			}
		}
		if found == 0 {
			send_to_char(ch, libc.CString("Couldn't find any such thing.\r\n"))
		} else {
			send_to_char(ch, libc.CString("\r\nFound %d matches.\r\n"), num)
		}
	}
}
func do_where(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if ADM_FLAGGED(ch, ADM_FULLWHERE) || ch.Admlevel > 4 {
		perform_immort_where(ch, &arg[0])
	} else {
		perform_mortal_where(ch, &arg[0])
	}
}
func do_levels(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [64936]byte
		i    uint64
		len_ uint64 = 0
		nlen uint64
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You ain't nothin' but a hound-dog.\r\n"))
		return
	}
	for i = 1; i < 101; i++ {
		if i == 100 {
			nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[100] %8s          : \r\n", add_commas(int64(level_exp(ch, 100)))))
		} else {
			nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[%2lld] %8s-%-8s : \r\n", i, add_commas(int64(level_exp(ch, int(i)))), add_commas(int64(level_exp(ch, int(i+1))-1))))
		}
		if len_+nlen >= uint64(64936) || nlen < 0 {
			break
		}
		len_ += nlen
	}
	page_string(ch.Desc, &buf[0], TRUE)
}
func do_consider(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		victim *char_data
		diff   int
	)
	one_argument(argument, &buf[0])
	if (func() *char_data {
		victim = get_char_vis(ch, &buf[0], nil, 1<<0)
		return victim
	}()) == nil {
		send_to_char(ch, libc.CString("Consider killing who?\r\n"))
		return
	}
	if victim == ch {
		send_to_char(ch, libc.CString("Easy!  Very easy indeed!\r\n"))
		return
	}
	diff = GET_LEVEL(victim) - GET_LEVEL(ch)
	if diff <= -10 {
		send_to_char(ch, libc.CString("Now where did that chicken go?\r\n"))
	} else if diff <= -5 {
		send_to_char(ch, libc.CString("You could do it with a needle!\r\n"))
	} else if diff <= -2 {
		send_to_char(ch, libc.CString("Easy.\r\n"))
	} else if diff <= -1 {
		send_to_char(ch, libc.CString("Fairly easy.\r\n"))
	} else if diff == 0 {
		send_to_char(ch, libc.CString("The perfect match!\r\n"))
	} else if diff <= 1 {
		send_to_char(ch, libc.CString("You could probably manage it.\r\n"))
	} else if diff <= 2 {
		send_to_char(ch, libc.CString("You might take a beating.\r\n"))
	} else if diff <= 3 {
		send_to_char(ch, libc.CString("You MIGHT win, maybe.\r\n"))
	} else if diff <= 5 {
		send_to_char(ch, libc.CString("Do you feel lucky? You better.\r\n"))
	} else if diff <= 10 {
		send_to_char(ch, libc.CString("Better bring some tough backup!\r\n"))
	} else if diff <= 25 {
		send_to_char(ch, libc.CString("Maybe if they are allergic to you, otherwise your last words will be 'Oh shit.'\r\n"))
	} else {
		send_to_char(ch, libc.CString("No chance.\r\n"))
	}
}
func do_diagnose(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		vict *char_data
	)
	one_argument(argument, &buf[0])
	if buf[0] != 0 {
		if (func() *char_data {
			vict = get_char_vis(ch, &buf[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		} else {
			send_to_char(ch, libc.CString("%s"), func() string {
				if int(vict.Sex) == SEX_MALE {
					return "He "
				}
				if int(vict.Sex) == SEX_FEMALE {
					return "She "
				}
				return "It "
			}())
			diag_char_to_char(vict, ch)
		}
	} else {
		if ch.Fighting != nil {
			send_to_char(ch, libc.CString("%s"), func() string {
				if int(ch.Fighting.Sex) == SEX_MALE {
					return "He "
				}
				if int(ch.Fighting.Sex) == SEX_FEMALE {
					return "She "
				}
				return "It "
			}())
			diag_char_to_char(ch.Fighting, ch)
		} else {
			send_to_char(ch, libc.CString("Diagnose who?\r\n"))
		}
	}
}

var ctypes [3]*byte = [3]*byte{libc.CString("off"), libc.CString("on"), libc.CString("\n")}

func cchoice_to_str(col *byte) *byte {
	var (
		buf    [256]byte
		s      *byte = nil
		i      int   = 0
		fg     int   = 0
		needfg int   = 0
		bold   int   = 0
	)
	if col == nil {
		buf[0] = 0
		return &buf[0]
	}
	for *col != 0 {
		if libc.StrChr(libc.CString(ANSISTART), *col) != nil {
			col = (*byte)(unsafe.Add(unsafe.Pointer(col), 1))
		} else {
			switch *col {
			case ANSISEP:
				fallthrough
			case ANSIEND:
				s = nil
			case '0':
				s = nil
			case '1':
				bold = 1
				s = nil
			case '5':
				s = libc.CString("blinking")
			case '7':
				s = libc.CString("reverse")
			case '8':
				s = libc.CString("invisible")
			case '3':
				col = (*byte)(unsafe.Add(unsafe.Pointer(col), 1))
				fg = 1
				switch *col {
				case '0':
					s = libc.CString(func() string {
						if bold != 0 {
							return "grey"
						}
						return "black"
					}())
					bold = 0
					fg = 1
				case '1':
					s = libc.CString("red")
					fg = 1
				case '2':
					s = libc.CString("green")
					fg = 1
				case '3':
					s = libc.CString("yellow")
					fg = 1
				case '4':
					s = libc.CString("blue")
					fg = 1
				case '5':
					s = libc.CString("magenta")
					fg = 1
				case '6':
					s = libc.CString("cyan")
					fg = 1
				case '7':
					s = libc.CString("white")
					fg = 1
				case 0:
					s = nil
				}
			case '4':
				col = (*byte)(unsafe.Add(unsafe.Pointer(col), 1))
				switch *col {
				case '0':
					s = libc.CString("on black")
					needfg = 1
					bold = 0
					fallthrough
				case '1':
					s = libc.CString("on red")
					needfg = 1
					bold = 0
					fallthrough
				case '2':
					s = libc.CString("on green")
					needfg = 1
					bold = 0
					fallthrough
				case '3':
					s = libc.CString("on yellow")
					needfg = 1
					bold = 0
					fallthrough
				case '4':
					s = libc.CString("on blue")
					needfg = 1
					bold = 0
					fallthrough
				case '5':
					s = libc.CString("on magenta")
					needfg = 1
					bold = 0
					fallthrough
				case '6':
					s = libc.CString("on cyan")
					needfg = 1
					bold = 0
					fallthrough
				case '7':
					s = libc.CString("on white")
					needfg = 1
					bold = 0
					fallthrough
				default:
					s = libc.CString("underlined")
				}
			default:
				s = nil
			}
			if s != nil {
				if needfg != 0 && fg == 0 {
					i += stdio.Snprintf(&buf[i], int(256-uintptr(i)), "%snormal", func() string {
						if i != 0 {
							return " "
						}
						return ""
					}())
					fg = 1
				}
				if i != 0 {
					i += stdio.Snprintf(&buf[i], int(256-uintptr(i)), " ")
				}
				if bold != 0 {
					i += stdio.Snprintf(&buf[i], int(256-uintptr(i)), "bright ")
					bold = 0
				}
				i += stdio.Snprintf(&buf[i], int(256-uintptr(i)), "%s", func() *byte {
					if s != nil {
						return s
					}
					return libc.CString("null 1")
				}())
				s = nil
			}
			col = (*byte)(unsafe.Add(unsafe.Pointer(col), 1))
		}
	}
	if fg == 0 {
		i += stdio.Snprintf(&buf[i], int(256-uintptr(i)), "%snormal", func() string {
			if i != 0 {
				return " "
			}
			return ""
		}())
	}
	return &buf[0]
}
func str_to_cchoice(str *byte, choice *byte) int {
	var (
		buf     [64936]byte
		bold    int = 0
		blink   int = 0
		uline   int = 0
		rev     int = 0
		invis   int = 0
		fg      int = 0
		bg      int = 0
		error   int = 0
		i       int
		len_    int = MAX_INPUT_LENGTH
		attribs [7]struct {
			Name *byte
			Ptr  *int
		} = [7]struct {
			Name *byte
			Ptr  *int
		}{{Name: libc.CString("bright"), Ptr: &bold}, {Name: libc.CString("bold"), Ptr: &bold}, {Name: libc.CString("underlined"), Ptr: &uline}, {Name: libc.CString("reverse"), Ptr: &rev}, {Name: libc.CString("blinking"), Ptr: &blink}, {Name: libc.CString("invisible"), Ptr: &invis}, {}}
		colors [13]struct {
			Name *byte
			Val  int
			Bold int
		} = [13]struct {
			Name *byte
			Val  int
			Bold int
		}{{Name: libc.CString("default"), Val: -1, Bold: 0}, {Name: libc.CString("normal"), Val: -1, Bold: 0}, {Name: libc.CString("black"), Val: 0, Bold: 0}, {Name: libc.CString("red"), Val: 1, Bold: 0}, {Name: libc.CString("green"), Val: 2, Bold: 0}, {Name: libc.CString("yellow"), Val: 3, Bold: 0}, {Name: libc.CString("blue"), Val: 4, Bold: 0}, {Name: libc.CString("magenta"), Val: 5, Bold: 0}, {Name: libc.CString("cyan"), Val: 6, Bold: 0}, {Name: libc.CString("white"), Val: 7, Bold: 0}, {Name: libc.CString("grey"), Val: 0, Bold: 1}, {Name: libc.CString("gray"), Val: 0, Bold: 1}, {}}
	)
	skip_spaces(&str)
	if unicode.IsDigit(rune(*str)) {
		libc.StrCpy(choice, str)
		for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(choice), i)) != 0 && (unicode.IsDigit(rune(*(*byte)(unsafe.Add(unsafe.Pointer(choice), i)))) || *(*byte)(unsafe.Add(unsafe.Pointer(choice), i)) == ';'); i++ {
		}
		error = int(libc.BoolToInt(*(*byte)(unsafe.Add(unsafe.Pointer(choice), i)) != 0))
		*(*byte)(unsafe.Add(unsafe.Pointer(choice), i)) = 0
		return error
	}
	for *str != 0 {
		str = any_one_arg(str, &buf[0])
		if libc.StrCmp(&buf[0], libc.CString("on")) == 0 {
			bg = 1
			continue
		}
		if fg == 0 {
			for i = 0; attribs[i].Name != nil; i++ {
				if libc.StrNCmp(attribs[i].Name, &buf[0], libc.StrLen(&buf[0])) == 0 {
					break
				}
			}
			if attribs[i].Name != nil {
				*attribs[i].Ptr = 1
				continue
			}
		}
		for i = 0; colors[i].Name != nil; i++ {
			if libc.StrNCmp(colors[i].Name, &buf[0], libc.StrLen(&buf[0])) == 0 {
				break
			}
		}
		if colors[i].Name == nil {
			error = 1
			continue
		}
		if colors[i].Val != -1 {
			if bg == 1 {
				bg = colors[i].Val + 40
			} else {
				fg = colors[i].Val + 30
				if colors[i].Bold != 0 {
					bold = 1
				}
			}
		}
	}
	*choice = byte(int8(func() int {
		i = 0
		return i
	}()))
	if bold != 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_BOLD)
	}
	if uline != 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_UNDERLINE)
	}
	if blink != 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_BLINK)
	}
	if rev != 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_REVERSE)
	}
	if invis != 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_INVIS)
	}
	if i == 0 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%s", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), AA_NORMAL)
	}
	if fg != 0 && fg != -1 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%d", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), fg)
	}
	if bg != 0 && bg != -1 {
		i += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(choice), i)), len_-i, "%s%d", func() string {
			if i != 0 {
				return ANSISEPSTR
			}
			return ""
		}(), bg)
	}
	return error
}

var default_color_choices [17]*byte = [17]*byte{libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), nil}

func do_color(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		p   *byte
	)
	_ = p
	var tp int
	p = any_one_arg(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Usage: color [ off | on ]\r\n"))
		return
	}
	if (func() int {
		tp = search_block(&arg[0], &ctypes[0], FALSE)
		return tp
	}()) == -1 {
		send_to_char(ch, libc.CString("Usage: color [ off | on ]\r\n"))
		return
	}
	switch tp {
	case C_OFF:
		ch.Player_specials.Pref[int(PRF_COLOR/32)] &= bitvector_t(int32(^(1 << (int(PRF_COLOR % 32)))))
	case C_ON:
		ch.Player_specials.Pref[int(PRF_COLOR/32)] |= bitvector_t(int32(1 << (int(PRF_COLOR % 32))))
	}
	send_to_char(ch, libc.CString("Your color is now @o%s@n.\r\n"), ctypes[tp])
}
func do_toggle(ch *char_data, argument *byte, cmd int, subcmd int) {
	var buf2 [4]byte
	if IS_NPC(ch) {
		return
	}
	if ch.Player_specials.Wimp_level == 0 {
		libc.StrCpy(&buf2[0], libc.CString("OFF"))
	} else {
		stdio.Sprintf(&buf2[0], "%-3.3d", ch.Player_specials.Wimp_level)
	}
	if ch.Admlevel != 0 {
		send_to_char(ch, libc.CString("      Buildwalk: %-3s    Clear Screen in OLC: %-3s\r\n"), func() string {
			if PRF_FLAGGED(ch, PRF_BUILDWALK) {
				return "ON"
			}
			return "OFF"
		}(), func() string {
			if PRF_FLAGGED(ch, PRF_CLS) {
				return "ON"
			}
			return "OFF"
		}())
		send_to_char(ch, libc.CString("      No Hassle: %-3s          Holylight: %-3s         Room Flags: %-3s\r\n"), func() string {
			if PRF_FLAGGED(ch, PRF_NOHASSLE) {
				return "ON"
			}
			return "OFF"
		}(), func() string {
			if PRF_FLAGGED(ch, PRF_HOLYLIGHT) {
				return "ON"
			}
			return "OFF"
		}(), func() string {
			if PRF_FLAGGED(ch, PRF_ROOMFLAGS) {
				return "ON"
			}
			return "OFF"
		}())
	}
	send_to_char(ch, libc.CString("Hit Pnt Display: %-3s         Brief Mode: %-3s     Summon Protect: %-3s\r\n   Move Display: %-3s       Compact Mode: %-3s           On Quest: %-3s\r\n    Exp Display: %-3s             NoTell: %-3s       Repeat Comm.: %-3s\r\n     Ki Display: %-3s               Deaf: %-3s         Wimp Level: %-3s\r\n Gossip Channel: %-3s    Auction Channel: %-3s      Grats Channel: %-3s\r\n      Auto Loot: %-3s          Auto Gold: %-3s        Color Level: %s\r\n     Auto Split: %-3s           Auto Sac: %-3s           Auto Mem: %-3s\r\n     View Order: %-3s        Auto Assist: %-3s     Auto Show Exit: %-3s\r\n    TNL Display: %-3s    "), func() string {
		if PRF_FLAGGED(ch, PRF_DISPHP) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_BRIEF) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if !PRF_FLAGGED(ch, PRF_SUMMONABLE) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_DISPMOVE) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_COMPACT) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_QUEST) {
			return "YES"
		}
		return "NO"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_DISPEXP) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_NOTELL) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if !PRF_FLAGGED(ch, PRF_NOREPEAT) {
			return "YES"
		}
		return "NO"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_DISPKI) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_DEAF) {
			return "YES"
		}
		return "NO"
	}(), &buf2[0], func() string {
		if !PRF_FLAGGED(ch, PRF_NOGOSS) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if !PRF_FLAGGED(ch, PRF_NOAUCT) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if !PRF_FLAGGED(ch, PRF_NOGRATZ) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOLOOT) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOGOLD) {
			return "ON"
		}
		return "OFF"
	}(), ctypes[func() int {
		if !IS_NPC(ch) {
			if PRF_FLAGGED(ch, PRF_COLOR) {
				return 1
			}
			return 0
		}
		return 0
	}()], func() string {
		if PRF_FLAGGED(ch, PRF_AUTOSPLIT) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOSAC) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOMEM) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_VIEWORDER) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOASSIST) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_AUTOEXIT) {
			return "ON"
		}
		return "OFF"
	}(), func() string {
		if PRF_FLAGGED(ch, PRF_DISPTNL) {
			return "ON"
		}
		return "OFF"
	}())
	if config_info.Play.Enable_compression != 0 {
		send_to_char(ch, libc.CString("    Compression: %-3s\r\n"), func() string {
			if !PRF_FLAGGED(ch, PRF_NOCOMPRESS) {
				return "ON"
			}
			return "OFF"
		}())
	}
}
func sort_commands_helper(a unsafe.Pointer, b unsafe.Pointer) int {
	return libc.StrCmp((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(*(*int)(a))))).Sort_as, (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(*(*int)(b))))).Sort_as)
}
func sort_commands() {
	var (
		a           int
		num_of_cmds int = 0
	)
	for *(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(num_of_cmds)))).Command != '\n' {
		num_of_cmds++
	}
	num_of_cmds++
	cmd_sort_info = &make([]int, num_of_cmds)[0]
	for a = 0; a < num_of_cmds; a++ {
		*(*int)(unsafe.Add(unsafe.Pointer(cmd_sort_info), unsafe.Sizeof(int(0))*uintptr(a))) = a
	}
	libc.Sort(unsafe.Pointer((*int)(unsafe.Add(unsafe.Pointer(cmd_sort_info), unsafe.Sizeof(int(0))*1))), uint32(int32(num_of_cmds-2)), uint32(unsafe.Sizeof(int(0))), func(arg1 unsafe.Pointer, arg2 unsafe.Pointer) int32 {
		return int32(sort_commands_helper(arg1, arg2))
	})
}
func do_commands(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		no      int
		i       int
		cmd_num int
		wizhelp int = 0
		socials int = 0
		vict    *char_data
		arg     [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] != 0 {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}()) == nil || IS_NPC(vict) {
			send_to_char(ch, libc.CString("Who is that?\r\n"))
			return
		}
		if GET_LEVEL(ch) < GET_LEVEL(vict) {
			send_to_char(ch, libc.CString("You can't see the commands of people above your level.\r\n"))
			return
		}
	} else {
		vict = ch
	}
	if subcmd == SCMD_SOCIALS {
		socials = 1
	} else if subcmd == SCMD_WIZHELP {
		wizhelp = 1
	}
	send_to_char(ch, libc.CString("The following %s%s are available to %s:\r\n"), func() string {
		if wizhelp != 0 {
			return "privileged "
		}
		return ""
	}(), func() string {
		if socials != 0 {
			return "socials"
		}
		return "commands"
	}(), func() string {
		if vict == ch {
			return "you"
		}
		return libc.GoString(GET_NAME(vict))
	}())
	for func() int {
		no = 1
		return func() int {
			cmd_num = 1
			return cmd_num
		}()
	}(); *(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(*(*int)(unsafe.Add(unsafe.Pointer(cmd_sort_info), unsafe.Sizeof(int(0))*uintptr(cmd_num))))))).Command != '\n'; cmd_num++ {
		i = *(*int)(unsafe.Add(unsafe.Pointer(cmd_sort_info), unsafe.Sizeof(int(0))*uintptr(cmd_num)))
		if int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_level) < 0 || GET_LEVEL(vict) < int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_level) {
			continue
		}
		if int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_admlevel) < 0 || vict.Admlevel < int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_admlevel) {
			continue
		}
		if int(libc.BoolToInt(int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_admlevel) >= ADMLVL_IMMORT)) != wizhelp {
			continue
		}
		if wizhelp == 0 && socials != int(libc.BoolToInt(libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command_pointer) == libc.FuncAddr(do_action) || libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command_pointer) == libc.FuncAddr(do_insult))) {
			continue
		}
		if check_disabled((*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))) != 0 {
			stdio.Sprintf(&arg[0], "(%s)", (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command)
		} else {
			stdio.Sprintf(&arg[0], "%s", (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command)
		}
		send_to_char(ch, libc.CString("%-11s%s"), &arg[0], func() string {
			if func() int {
				p := &no
				x := *p
				*p++
				return x
			}()%7 == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	if no%7 != 1 {
		send_to_char(ch, libc.CString("\r\n"))
	}
}
func free_history(ch *char_data, type_ int) {
	var (
		tmp  *txt_block = (ch.Player_specials.Comm_hist[type_])
		ftmp *txt_block
	)
	for (func() *txt_block {
		ftmp = tmp
		return ftmp
	}()) != nil {
		tmp = tmp.Next
		if ftmp.Text != nil {
			libc.Free(unsafe.Pointer(ftmp.Text))
		}
		libc.Free(unsafe.Pointer(ftmp))
	}
	ch.Player_specials.Comm_hist[type_] = nil
}
func do_history(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		type_ int
	)
	one_argument(argument, &arg[0])
	type_ = search_block(&arg[0], &history_types[0], FALSE)
	if arg[0] == 0 || type_ < 0 {
		var i int
		send_to_char(ch, libc.CString("Usage: history <"))
		for i = 0; *history_types[i] != '\n'; i++ {
			if i != 3 && ch.Admlevel <= 0 || ch.Admlevel >= 1 {
				send_to_char(ch, libc.CString(" %s "), history_types[i])
			}
			if *history_types[i+1] == '\n' {
				send_to_char(ch, libc.CString(">\r\n"))
			} else {
				if i != 3 && ch.Admlevel <= 0 || ch.Admlevel >= 1 {
					send_to_char(ch, libc.CString("|"))
				}
			}
		}
		return
	}
	if (ch.Player_specials.Comm_hist[type_]) != nil && (ch.Player_specials.Comm_hist[type_]).Text != nil && *(ch.Player_specials.Comm_hist[type_]).Text != 0 {
		var tmp *txt_block
		for tmp = ch.Player_specials.Comm_hist[type_]; tmp != nil; tmp = tmp.Next {
			send_to_char(ch, libc.CString("%s"), tmp.Text)
		}
	} else {
		send_to_char(ch, libc.CString("You have no history in that channel.\r\n"))
	}
}
func add_history(ch *char_data, str *byte, type_ int) {
	var (
		i        int = 0
		time_str [64936]byte
		buf      [64936]byte
		tmp      *txt_block
		ct       libc.Time
	)
	if IS_NPC(ch) {
		return
	}
	tmp = ch.Player_specials.Comm_hist[type_]
	ct = libc.GetTime(nil)
	strftime(&time_str[0], uint64(64936), libc.CString("%H:%M "), libc.LocalTime(&ct))
	stdio.Sprintf(&buf[0], "%s%s", &time_str[0], str)
	if tmp == nil {
		ch.Player_specials.Comm_hist[type_] = new(txt_block)
		(ch.Player_specials.Comm_hist[type_]).Text = libc.StrDup(&buf[0])
	} else {
		for tmp.Next != nil {
			tmp = tmp.Next
		}
		tmp.Next = new(txt_block)
		tmp.Next.Text = libc.StrDup(&buf[0])
		for tmp = ch.Player_specials.Comm_hist[type_]; tmp != nil; func() int {
			tmp = tmp.Next
			return func() int {
				p := &i
				x := *p
				*p++
				return x
			}()
		}() {
		}
		for ; i > HIST_LENGTH && (ch.Player_specials.Comm_hist[type_]) != nil; i-- {
			tmp = ch.Player_specials.Comm_hist[type_]
			ch.Player_specials.Comm_hist[type_] = tmp.Next
			if tmp.Text != nil {
				libc.Free(unsafe.Pointer(tmp.Text))
			}
			libc.Free(unsafe.Pointer(tmp))
		}
	}
	if type_ != HIST_ALL {
		add_history(ch, str, HIST_ALL)
	}
}
func do_scan(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i        int
		newroom  int
		dirnames [12]*byte = [12]*byte{libc.CString("North"), libc.CString("East"), libc.CString("South"), libc.CString("West"), libc.CString("Up"), libc.CString("Down"), libc.CString("Northwest"), libc.CString("Northeast"), libc.CString("Southeast"), libc.CString("Southwest"), libc.CString("Inside"), libc.CString("Outside")}
	)
	if int(ch.Position) < POS_SLEEPING {
		send_to_char(ch, libc.CString("You can't see anything but stars!\n\r"))
		return
	}
	if !AWAKE(ch) {
		send_to_char(ch, libc.CString("You must be dreaming.\n\r"))
		return
	}
	if AFF_FLAGGED(ch, AFF_BLIND) {
		send_to_char(ch, libc.CString("You can't see a damn thing, you're blind!\n\r"))
		return
	}
	if PLR_FLAGGED(ch, PLR_EYEC) {
		send_to_char(ch, libc.CString("You can't see a damned thing, your eyes are closed!\r\n"))
		return
	}
	for i = 0; i < 10; i++ {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]) != nil {
			if room_is_dark(ch.In_room) != 0 && ch.Admlevel < ADMLVL_IMMORT && !AFF_FLAGGED(ch, AFF_INFRAVISION) {
				send_to_char(ch, libc.CString("%s: DARK\n\r"), dirnames[i])
				continue
			}
			if CAN_GO(ch, i) {
				newroom = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i].To_room)
				send_to_char(ch, libc.CString("@w-----------------------------------------@n\r\n"))
				send_to_char(ch, libc.CString("          %s%s: %s %s\n\r"), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KCYN
					}
					return KNUL
				}(), dirnames[i], func() *byte {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Name != nil {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Name
					}
					return libc.CString("You don't think you saw what you just saw.")
				}(), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KNRM
					}
					return KNUL
				}())
				send_to_char(ch, libc.CString("@W          -----------------          @n\r\n"))
				list_obj_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Contents, ch, SHOW_OBJ_LONG, FALSE)
				list_char_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).People, ch)
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect >= 1 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect <= 5 {
					send_to_char(ch, libc.CString("@rLava@w is pooling in someplaces here...@n\r\n"))
				}
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect >= 6 {
					send_to_char(ch, libc.CString("@RLava@r covers pretty much the entire area!@n\r\n"))
				}
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]).To_room)))).Dir_option[i]) != nil && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]).To_room)))).Dir_option[i]).To_room != 0 {
					newroom = int(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]).To_room)))).Dir_option[i]).To_room)
					if newroom != int(-1) && (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]).To_room)))).Dir_option[i]).Exit_info&(1<<1)) == 0 {
						if room_is_dark(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[i]).To_room)))).Dir_option[i]).To_room) == 0 {
							send_to_char(ch, libc.CString("@w-----------------------------------------@n\r\n"))
							send_to_char(ch, libc.CString("          %sFar %s: %s %s\n\r"), func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KCYN
								}
								return KNUL
							}(), dirnames[i], func() *byte {
								if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Name != nil {
									return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Name
								}
								return libc.CString("You don't think you saw what you just saw.")
							}(), func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KNRM
								}
								return KNUL
							}())
							send_to_char(ch, libc.CString("@W          -----------------          @n\r\n"))
							list_obj_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Contents, ch, SHOW_OBJ_LONG, FALSE)
							list_char_to_char((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).People, ch)
							if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect >= 1 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect <= 5 {
								send_to_char(ch, libc.CString("@rLava@w is pooling in someplaces here...@n\r\n"))
							}
							if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(newroom)))).Geffect >= 6 {
								send_to_char(ch, libc.CString("@RLava@r covers pretty much the entire area!@n\r\n"))
							}
						} else {
							send_to_char(ch, libc.CString("%s<-> %sFar %s: Too dark to tell! %s<->%s\r\n"), func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KMAG
								}
								return KNUL
							}(), func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KCYN
								}
								return KNUL
							}(), dirnames[i], func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KMAG
								}
								return KNUL
							}(), func() string {
								if (func() int {
									if !IS_NPC(ch) {
										if PRF_FLAGGED(ch, PRF_COLOR) {
											return 1
										}
										return 0
									}
									return 0
								}()) >= C_ON {
									return KNRM
								}
								return KNUL
							}())
						}
					}
				}
			}
		}
	}
	send_to_char(ch, libc.CString("@w-----------------------------------------@n\r\n"))
}
func do_toplist(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var file *stdio.File
	var fname [40]byte
	var filler [50]byte
	var line [256]byte
	var points [25]int64 = [25]int64{}
	_ = points
	var stats int64
	var title [25]*byte = [25]*byte{0: libc.CString("")}
	var count int = 0
	var x int = 0
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("toplist")) == 0 {
		send_to_char(ch, libc.CString("The toplist file does not exist."))
		return
	} else if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		send_to_char(ch, libc.CString("The toplist file does not exist."))
		return
	}
	for int(file.IsEOF()) == 0 || count < 25 {
		get_line(file, &line[0])
		switch count {
		default:
			stdio.Sscanf(&line[0], "%s %lld\n", &filler[0], &stats)
		}
		title[count] = libc.StrDup(&filler[0])
		points[count] = stats
		count++
		filler[0] = '\x00'
	}
	send_to_char(ch, libc.CString("@D-=[@BDBAT Top Lists for @REra@C %d@D]=-@n\r\n"), CURRENT_ERA)
	for x <= count {
		switch x {
		case 0:
			send_to_char(ch, libc.CString("       @D-@RPowerlevel@D-@n\r\n"))
			send_to_char(ch, libc.CString("    @D|@c1@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 1:
			send_to_char(ch, libc.CString("    @D|@c2@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 2:
			send_to_char(ch, libc.CString("    @D|@c3@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 3:
			send_to_char(ch, libc.CString("    @D|@c4@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 4:
			send_to_char(ch, libc.CString("    @D|@c5@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 5:
			send_to_char(ch, libc.CString("       @D-@BKi        @D-@n\r\n"))
			send_to_char(ch, libc.CString("    @D|@c1@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 6:
			send_to_char(ch, libc.CString("    @D|@c2@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 7:
			send_to_char(ch, libc.CString("    @D|@c3@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 8:
			send_to_char(ch, libc.CString("    @D|@c4@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 9:
			send_to_char(ch, libc.CString("    @D|@c5@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 10:
			send_to_char(ch, libc.CString("       @D-@GStamina   @D-@n\r\n"))
			send_to_char(ch, libc.CString("    @D|@c1@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 11:
			send_to_char(ch, libc.CString("    @D|@c2@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 12:
			send_to_char(ch, libc.CString("    @D|@c3@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 13:
			send_to_char(ch, libc.CString("    @D|@c4@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 14:
			send_to_char(ch, libc.CString("    @D|@c5@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 15:
			send_to_char(ch, libc.CString("       @D-@gZenni     @D-@n\r\n"))
			send_to_char(ch, libc.CString("    @D|@c1@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 16:
			send_to_char(ch, libc.CString("    @D|@c2@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 17:
			send_to_char(ch, libc.CString("    @D|@c3@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 18:
			send_to_char(ch, libc.CString("    @D|@c4@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 19:
			send_to_char(ch, libc.CString("    @D|@c5@W: @C%13s@D|@n\r\n"), title[x])
			libc.Free(unsafe.Pointer(title[x]))
		case 20:
			libc.Free(unsafe.Pointer(title[x]))
		case 21:
			libc.Free(unsafe.Pointer(title[x]))
		case 22:
			libc.Free(unsafe.Pointer(title[x]))
		case 23:
			libc.Free(unsafe.Pointer(title[x]))
		case 24:
			libc.Free(unsafe.Pointer(title[x]))
		}
		x++
	}
	file.Close()
}
func do_whois(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf       [2048]byte
		clan      int        = FALSE
		immlevels [8]*byte   = [8]*byte{0: libc.CString("[Mortal]"), 1: libc.CString("[Enforcer]"), 2: libc.CString("[First Class Enforcer]"), 3: libc.CString("[High Enforcer]"), 4: libc.CString("[Vice Admin]"), 5: libc.CString("[Administrator]"), 6: libc.CString("[Implementor]")}
		victim    *char_data = nil
	)
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Who?\r\n"))
	} else {
		victim = new(char_data)
		clear_char(victim)
		victim.Player_specials = new(player_special_data)
		if load_char(argument, victim) >= 0 {
			if victim.Clan != nil {
				if libc.StrStr(victim.Clan, libc.CString("None")) == nil {
					stdio.Sprintf(&buf[0], "%s", victim.Clan)
					clan = TRUE
				}
				if libc.StrStr(victim.Clan, libc.CString("Applying")) != nil {
					stdio.Sprintf(&buf[0], "%s", victim.Clan)
					clan = TRUE
				}
			}
			if victim.Clan == nil || libc.StrStr(victim.Clan, libc.CString("None")) != nil {
				clan = FALSE
			}
			send_to_char(ch, libc.CString("@D~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\r\n"))
			if victim.Admlevel >= ADMLVL_IMMORT {
				send_to_char(ch, libc.CString("@cName     @D: @G%s\r\n"), GET_NAME(victim))
				send_to_char(ch, libc.CString("@cImm Level@D: @G%s\r\n"), immlevels[victim.Admlevel])
				send_to_char(ch, libc.CString("@cTitle    @D: @G%s\r\n"), GET_TITLE(victim))
			} else {
				send_to_char(ch, libc.CString("@cName  @D: @w%s\r\n@cSensei@D: @w%s\r\n@cRace  @D: @w%s\r\n@cTitle @D: @w%s@n\r\n@cClan  @D: @w%s@n\r\n"), GET_NAME(victim), pc_class_types[int(victim.Chclass)], pc_race_types[int(victim.Race)], GET_TITLE(victim), &func() [2048]byte {
					if clan != 0 {
						return buf
					}
					return func() [2048]byte {
						var t [2048]byte
						copy(t[:], []byte("None."))
						return t
					}()
				}()[0])
				if clan == TRUE && libc.StrStr(victim.Clan, libc.CString("Applying")) == nil {
					if checkCLAN(victim) == TRUE {
						clanRANKD(victim.Clan, ch, victim)
					}
				}
			}
			send_to_char(ch, libc.CString("@D~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("There is no such player.\r\n"))
		}
		libc.Free(unsafe.Pointer(victim))
	}
}
func search_in_direction(ch *char_data, dir int) {
	var (
		check     int = FALSE
		skill_lvl int
		dchide    int = 20
	)
	send_to_char(ch, libc.CString("You search for secret doors.\r\n"))
	act(libc.CString("$n searches the area intently."), TRUE, ch, nil, nil, TO_ROOM)
	skill_lvl = GET_SKILL(ch, SKILL_SEARCH)
	if int(ch.Race) == RACE_TRUFFLE || int(ch.Race) == RACE_HUMAN {
		skill_lvl = skill_lvl + 2
	}
	if int(ch.Race) == RACE_HALFBREED {
		skill_lvl = skill_lvl + 1
	}
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) != nil {
		dchide = ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dchide
	}
	if skill_lvl > dchide {
		check = TRUE
	}
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) != nil {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).General_description != nil && !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<4) {
			send_to_char(ch, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).General_description)
		} else if !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<4) {
			send_to_char(ch, libc.CString("There is a normal exit there.\r\n"))
		} else if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<0) && EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<4) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword != nil && check == TRUE {
			send_to_char(ch, libc.CString("There is a hidden door keyword: '%s' %sthere.\r\n"), fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword), func() string {
				if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<1) {
					return ""
				}
				return "open "
			}())
		} else {
			send_to_char(ch, libc.CString("There is no exit there.\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("There is no exit there.\r\n"))
	}
}
