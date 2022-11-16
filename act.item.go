package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unicode"
	"unsafe"
)

var obj_selling *obj_data = nil
var ch_selling *char_data = nil
var ch_buying *char_data = nil
var curbid int = 0
var aucstat int = AUC_NULL_STATE
var auctioneer [11]*byte = [11]*byte{libc.CString("@D[@CAUCTION@c: @C$n@W puts $p@W up for sale at @Y%d@W zenni.@D]@n"), libc.CString("@D[@CAUCTION@c: @W$p@W at @Y%d@W zenni going once!@D]@n"), libc.CString("@D[@CAUCTION@c: @W$p@W at @Y%d@W zenni going twice!@D]@n"), libc.CString("@D[@CAUCTION@c: @WLast call: $p@W going for @Y%d@W zenni.@D]@n"), libc.CString("@D[@CAUCTION@c: @WUnfortunately $p@W is unsold, returning it to $n. @D]@n"), libc.CString("@D[@CAUCTION@c: @WSOLD! $p@W to @C$n@W for @Y%d@W zenni!@D]@n"), libc.CString("@D[@CAUCTION@c: @WSorry, @C$n@W has cancelled the auction.@D]@n"), libc.CString("@D[@CAUCTION@c: @WSorry, @C$n@W has left us, the auction can't go on.@D]@n"), libc.CString("@D[@CAUCTION@c: @WSorry, $p@W has been confiscated, shame on you $n.@D]@n"), libc.CString("@D[@CAUCTION@c: @C$n@W is selling $p@W for @Y%d@W zenni.@D]@n"), libc.CString("@D[@CAUCTION@c: @C$n@W bids @Y%d@W zenni on $p@W.@D]@n")}
var buf [64936]byte

func do_refuel(ch *char_data, argument *byte, cmd int, subcmd int) {
	var controls *obj_data
	if (func() *obj_data {
		controls = find_control(ch)
		return controls
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou need to be in the cockpit to place a new fuel canister into the ship.\r\n"))
		return
	}
	var fuel *obj_data = find_obj_in_list_vnum(ch.Carrying, 17290)
	if fuel == nil {
		send_to_char(ch, libc.CString("You do not have any fuel canisters on you.\r\n"))
		return
	}
	var max int = 0
	if GET_OBJ_VNUM(controls) >= 44000 && GET_OBJ_VNUM(controls) <= 0xACA7 {
		max = 300
	} else if GET_OBJ_VNUM(controls) >= 44200 && GET_OBJ_VNUM(controls) <= 0xADD3 {
		max = 500
	} else if GET_OBJ_VNUM(controls) >= 44200 && GET_OBJ_VNUM(controls) <= 0xAFC7 {
		max = 1000
	}
	if (controls.Value[2]) == max {
		send_to_char(ch, libc.CString("The ship is full on fuel!\r\n"))
		return
	} else {
		if (controls.Value[2])+int(fuel.Weight*4) > max {
			controls.Value[2] = max
		} else {
			controls.Value[2] += int(fuel.Weight * 4)
		}
		extract_obj(fuel)
		send_to_char(ch, libc.CString("You place the fuel canister into the ship. Within seconds the fuel has been extracted from the canister into the ships' internal tanks.\r\n"))
	}
}
func can_harvest(plant *obj_data) int {
	switch GET_OBJ_VNUM(plant) {
	case 250:
		fallthrough
	case 1129:
		fallthrough
	case 17210:
		fallthrough
	case 0x433B:
		fallthrough
	case 0x433E:
		fallthrough
	case 0x4340:
		fallthrough
	case 0x4342:
		fallthrough
	case 17220:
		fallthrough
	case 0x4346:
		fallthrough
	case 0x4348:
		fallthrough
	case 0x434A:
		fallthrough
	case 3702:
		return TRUE
	}
	return FALSE
}
func harvest_plant(ch *char_data, plant *obj_data) {
	var (
		extract int       = FALSE
		reward  int       = rand_number(5, 15)
		count   int       = reward
		fruit   *obj_data = nil
	)
	if (plant.Value[VAL_SOILQ]) > 7 {
		reward += 10
		send_to_char(ch, libc.CString("@GThe soil seems to have made the plant exteremely bountiful"))
	} else if (plant.Value[VAL_SOILQ]) >= 5 {
		reward += 6
		send_to_char(ch, libc.CString("@GThe soil seems to have made the plant very bountiful"))
	} else if (plant.Value[VAL_SOILQ]) >= 3 {
		reward += 4
		send_to_char(ch, libc.CString("@GThe soil seems to have made the plant bountiful"))
	} else if (plant.Value[VAL_SOILQ]) > 0 {
		reward += 2
		send_to_char(ch, libc.CString("@GThe soil seems to have made the plant a bit more bountiful"))
	}
	var skill int = GET_SKILL(ch, SKILL_GARDENING)
	if skill >= 100 {
		reward += 10
		send_to_char(ch, libc.CString(" and your outstanding skill has helped the plant be more bountiful yet!@n\r\n"))
	} else if skill >= 90 {
		reward += 8
		send_to_char(ch, libc.CString(" and your great skill has also helped the plant be more bountiful yet!@n\r\n"))
	} else if skill >= 80 {
		reward += 5
		send_to_char(ch, libc.CString(" and your good skill has also helped the plant be more bountiful yet!@n\r\n"))
	} else if skill >= 50 {
		reward += 3
		send_to_char(ch, libc.CString(" and your decent skill has also helped the plant be more bountiful yet!@n\r\n"))
	} else if skill >= 40 {
		reward += 2
		send_to_char(ch, libc.CString(" and your mastery of the basics of gardening has also helped the plant be more bountiful yet!@n\r\n"))
	} else if skill >= 30 {
		reward += 1
		send_to_char(ch, libc.CString(" and you somehow managed to make the plant slightly more bountiful with what little you know!@n\r\n"))
	} else {
		send_to_char(ch, libc.CString(".@n\r\n"))
	}
	count = reward
	switch GET_OBJ_VNUM(plant) {
	case 250:
		if reward > 2 {
			reward = 2
			count = 2
		}
		for count > 0 {
			fruit = read_object(1, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = FALSE
	case 1129:
		for count > 0 {
			fruit = read_object(1131, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 17210:
		reward += 2
		count += 2
		for count > 0 {
			fruit = read_object(0x433C, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x433B:
		reward += 2
		count += 2
		for count > 0 {
			fruit = read_object(0x433D, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x433E:
		reward += 1
		count += 1
		for count > 0 {
			fruit = read_object(0x433F, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x4340:
		for count > 0 {
			fruit = read_object(0x4341, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x4342:
		reward += 14
		count += 14
		for count > 0 {
			fruit = read_object(0x4343, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 17220:
		reward -= int(float64(reward) * 0.75)
		count = reward
		for count > 0 {
			fruit = read_object(0x4345, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x4346:
		reward += rand_number(1, 3)
		count = reward
		for count > 0 {
			fruit = read_object(0x4347, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x4348:
		reward += 10
		count = reward
		for count > 0 {
			fruit = read_object(0x4349, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 0x434A:
		reward -= 8
		if reward < 0 {
			reward = 1
		}
		count = reward
		for count > 0 {
			fruit = read_object(0x434B, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	case 3702:
		reward -= 2
		count = reward
		for count > 0 {
			fruit = read_object(3703, VIRTUAL)
			obj_to_char(fruit, ch)
			count -= 1
		}
		send_to_char(ch, libc.CString("@YYou harvest @D[@G%d@D]@Y @g%s@Y!@n\r\n"), reward, fruit.Short_description)
		extract = TRUE
	default:
		send_to_imm(libc.CString("ERROR: Harvest plant called for illegitimate plant, VNUM %d."), GET_OBJ_VNUM(plant))
	}
	if extract == TRUE {
		send_to_char(ch, libc.CString("@wThe harvesting process has killed the plant. Do not worry, this is normal for that type.@n\r\n"))
		extract_obj(plant)
	} else {
		plant.Value[VAL_MATURITY] = 3
		plant.Value[VAL_GROWTH] = 0
	}
}
func do_garden(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
		obj  *obj_data
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if GET_SKILL(ch, SKILL_GARDENING) == 0 && slot_count(ch)+1 <= ch.Skill_slots {
		var numb int = rand_number(8, 16)
		for {
			ch.Skills[SKILL_GARDENING] = int8(numb)
			if true {
				break
			}
		}
		send_to_char(ch, libc.CString("@GYou learn the very basics of gardening.\r\n"))
	} else if GET_SKILL(ch, SKILL_GARDENING) == 0 && slot_count(ch)+1 > ch.Skill_slots {
		send_to_char(ch, libc.CString("You need additional skill slots to pick up the skill linked with this.\r\n"))
		return
	}
	if arg[0] != 0 {
		if libc.StrCaseCmp(&arg[0], libc.CString("collect")) == 0 {
			var shovel *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 254)
			if shovel == nil {
				send_to_char(ch, libc.CString("You need a shovel in order to collect soil.\r\n"))
				return
			}
			if SECT(ch.In_room) != SECT_FOREST && SECT(ch.In_room) != SECT_FIELD && SECT(ch.In_room) != SECT_MOUNTAIN && SECT(ch.In_room) != SECT_HILLS {
				send_to_char(ch, libc.CString("You can not collect soil from this area.\r\n"))
				return
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_FERTILE1) {
				var soil *obj_data = read_object(math.MaxUint8, VIRTUAL)
				obj_to_char(soil, ch)
				act(libc.CString("@yYou sink your shovel into the soft ground and manage to dig up a pile of fertile soil!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@y sinks $s shovel into the soft ground and manages to dig up a pile of fertile soil!@n"), TRUE, ch, nil, nil, TO_ROOM)
				soil.Value[0] = 8
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				return
			} else if ROOM_FLAGGED(ch.In_room, ROOM_FERTILE2) {
				var soil *obj_data = read_object(math.MaxUint8, VIRTUAL)
				obj_to_char(soil, ch)
				act(libc.CString("@yYou sink your shovel into the soft ground and manage to dig up a pile of good soil!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@y sinks $s shovel into the soft ground and manages to dig up a pile of good soil!@n"), TRUE, ch, nil, nil, TO_ROOM)
				soil.Value[0] = rand_number(5, 7)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				return
			} else {
				var soil *obj_data = read_object(math.MaxUint8, VIRTUAL)
				obj_to_char(soil, ch)
				act(libc.CString("@yYou sink your shovel into the soft ground and manage to dig up a pile of soil!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@y sinks $s shovel into the soft ground and manages to dig up a pile of soil!@n"), TRUE, ch, nil, nil, TO_ROOM)
				soil.Value[0] = rand_number(0, 4)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				return
			}
		}
	}
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: garden (plant) ( water | harvest | dig | plant | pick )\r\n"))
		send_to_char(ch, libc.CString("Syntax: garden collect [Will collect soil from a room with soil.\r\n"))
		return
	}
	if !ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) && !ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
		send_to_char(ch, libc.CString("You are not even in a garden!\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg2[0], libc.CString("plant")) == 0 {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("What are you trying to plant?\r\n"))
			send_to_char(ch, libc.CString("Syntax: garden (plant in inventory) plant\r\n"))
			return
		}
	} else {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("That plant doesn't seem to be here.\r\n"))
			return
		}
	}
	if obj == nil {
		send_to_char(ch, libc.CString("What plant are you gardening?\r\n"))
		return
	}
	var cost int64 = int64((float64(ch.Max_move) * 0.005) + float64(rand_number(50, 150)))
	var skill int = GET_SKILL(ch, SKILL_GARDENING)
	if ch.Move < cost {
		send_to_char(ch, libc.CString("@WYou need at least @G%s@W stamina to garden.\r\n"), add_commas(cost))
		return
	} else {
		if libc.StrCaseCmp(&arg2[0], libc.CString("water")) == 0 {
			var water *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 251)
			if water == nil {
				send_to_char(ch, libc.CString("You do not have any grow water!\r\n"))
				return
			} else if (obj.Value[VAL_WATERLEVEL]) >= 500 {
				send_to_char(ch, libc.CString("You stop as you realize that the plant already has enough water.\r\n"))
				return
			} else if (obj.Value[VAL_WATERLEVEL]) <= -10 {
				send_to_char(ch, libc.CString("The plant is dead!\r\n"))
				return
			} else if skill < axion_dice(0) {
				act(libc.CString("@GAs you go to water @g$p@G you end up sloppily wasting about half of it on the ground.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G takes a bottle of grow water and sloshes some of it on @g$p@G.@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				obj.Value[VAL_WATERLEVEL] += 40
				if (obj.Value[VAL_WATERLEVEL]) > 500 {
					obj.Value[VAL_WATERLEVEL] = 500
					send_to_char(ch, libc.CString("@YThe plant is now at full water level.@n\r\n"))
				}
				extract_obj(water)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			} else {
				act(libc.CString("@GYou calmly and expertly pour the grow water on @g$p@G.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G calmly and expertly pours some grow water on @g$p@G.@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				obj.Value[VAL_WATERLEVEL] += 225
				if (obj.Value[VAL_WATERLEVEL]) >= 500 {
					obj.Value[VAL_WATERLEVEL] = 500
					send_to_char(ch, libc.CString("@YThe plant is now at full water level.@n\r\n"))
				}
				extract_obj(water)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			}
		} else if libc.StrCaseCmp(&arg2[0], libc.CString("harvest")) == 0 {
			var clippers *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 253)
			if clippers == nil {
				send_to_char(ch, libc.CString("You do not have any working gardening clippers!\r\n"))
				return
			} else if can_harvest(obj) == FALSE {
				send_to_char(ch, libc.CString("You can not harvest that plant. Instead, Syntax: garden (plant) (pick)\r\n"))
				return
			} else if (obj.Value[VAL_WATERLEVEL]) <= -10 {
				send_to_char(ch, libc.CString("That plant is dead!\r\n"))
				return
			} else if (obj.Value[VAL_MATURITY]) < (obj.Value[VAL_MAXMATURE]) {
				send_to_char(ch, libc.CString("You stop as you realize that the plant isn't mature enough to harvest.\r\n"))
				return
			} else if skill < axion_dice(-5) {
				act(libc.CString("@GAs you go to harvest @g$p@G you end up cutting it in half instead!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G attempts to harvest @g$p@G with $s clippers, but accidently cuts the plant in half!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				extract_obj(obj)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			} else {
				act(libc.CString("@GYou calmly and expertly harvest @g$p@G.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G calmly and expertly harvests @g$p@G.@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				clippers.Value[VAL_ALL_HEALTH] -= 1
				if (clippers.Value[VAL_ALL_HEALTH]) <= 0 {
					send_to_char(ch, libc.CString("The clippers are now too dull to use.\r\n"))
					return
				}
				harvest_plant(ch, obj)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			}
		} else if libc.StrCaseCmp(&arg2[0], libc.CString("dig")) == 0 {
			var shovel *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 254)
			if shovel == nil {
				send_to_char(ch, libc.CString("You do not have any working gardening shovels!\r\n"))
				return
			} else {
				act(libc.CString("@GYou calmly dig up @g$p@G.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G calmly digs up @g$p@G.@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				obj_from_room(obj)
				obj_to_char(obj, ch)
				obj.Value[VAL_SOILQ] = 0
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			}
		} else if libc.StrCaseCmp(&arg2[0], libc.CString("plant")) == 0 {
			var shovel *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 254)
			if shovel == nil {
				send_to_char(ch, libc.CString("You do not have any working gardening shovels!\r\n"))
				return
			}
			var soil *obj_data = find_obj_in_list_vnum_good(ch.Carrying, math.MaxUint8)
			if soil == nil {
				send_to_char(ch, libc.CString("You don't have any real soil.\r\n"))
			} else if check_saveroom_count(ch, nil) > 7 && ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) {
				send_to_char(ch, libc.CString("This room already has all its planters full. Try digging up some plants.\r\n"))
				return
			} else if check_saveroom_count(ch, nil) > 19 && ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
				send_to_char(ch, libc.CString("This room already has all its planters full. Try digging up some plants.\r\n"))
				return
			} else if skill < axion_dice(-5) {
				act(libc.CString("@GYou end up digging a hole too shallow to hold @g$p@G. Better try again.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G digs a very shallow hole in one of the planters and then realizes @g$p@G won't fit in it.@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Move -= cost
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			} else {
				act(libc.CString("@GYou dig a proper sized hole and plant @g$p@G in it.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G digs a proper sized hole in a planter and plants @g$p@G in it.@n"), TRUE, ch, obj, nil, TO_ROOM)
				obj_from_char(obj)
				obj_to_room(obj, ch.In_room)
				ch.Move -= cost
				obj.Value[VAL_MAXMATURE] = 6
				obj.Value[VAL_MATGOAL] = 200
				obj.Value[VAL_SOILQ] = soil.Value[0]
				switch obj.Value[VAL_SOILQ] {
				case 1:
					obj.Value[VAL_MATGOAL] -= 10
				case 2:
					obj.Value[VAL_MATGOAL] -= 15
				case 3:
					obj.Value[VAL_MATGOAL] -= 20
				case 4:
					obj.Value[VAL_MATGOAL] -= 25
				case 5:
					obj.Value[VAL_MATGOAL] -= 50
				case 6:
					obj.Value[VAL_MATGOAL] -= 60
				case 7:
					obj.Value[VAL_MATGOAL] -= 70
				default:
					obj.Value[VAL_MATGOAL] -= 80
				}
				extract_obj(soil)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
			}
		} else if libc.StrCaseCmp(&arg2[0], libc.CString("pick")) == 0 {
			if !OBJ_FLAGGED(obj, ITEM_MATURE) {
				send_to_char(ch, libc.CString("You can't pick that type of plant. Syntax: garden (plant) harvest\r\n"))
				return
			} else if (obj.Value[VAL_MATURITY]) < (obj.Value[VAL_MAXMATURE]) {
				send_to_char(ch, libc.CString("That plant is not mature enough yet.\r\n"))
				return
			} else if skill < axion_dice(-5) {
				act(libc.CString("@GYou end up shredding @g$p@G with your clumsy unskilled hands.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G grabs a hold of @g$p@G and shreds it in an attempt to pick it.@n"), TRUE, ch, obj, nil, TO_ROOM)
				return
			} else {
				act(libc.CString("@GYou grab a hold of @g$p@G and carefully pick it out of the soil.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@g$n@G grabs a hold of @g$p@G and carefully picks it out of the soil.@n"), TRUE, ch, obj, nil, TO_ROOM)
				obj_from_room(obj)
				obj_to_char(obj, ch)
				ch.Move -= cost
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				improve_skill(ch, SKILL_GARDENING, 0)
				return
			}
		} else {
			send_to_char(ch, libc.CString("Syntax: garden (plant) ( water | harvest | dig | plant | pick )\r\n"))
			send_to_char(ch, libc.CString("Syntax: garden collect [Will collect soil from a room with soil.\r\n"))
			return
		}
	}
}
func has_housekey(ch *char_data, obj *obj_data) int {
	var (
		obj2     *obj_data = nil
		next_obj *obj_data
	)
	for obj2 = ch.Carrying; obj2 != nil; obj2 = next_obj {
		next_obj = obj2.Next_content
		if OBJ_FLAGGED(obj, ITEM_DUPLICATE) {
			continue
		}
		if GET_OBJ_VNUM(obj) == 0x4972 {
			if GET_OBJ_VNUM(obj2) == 18800 {
				return 1
			}
		} else {
			if GET_OBJ_VNUM(obj2) == GET_OBJ_VNUM(obj)-1 {
				return 1
			}
		}
	}
	return 0
}
func do_pack(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj  *obj_data
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Pack up which type of house capsule?\nSyntax: pack (target)\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("That house item doesn't seem to be around.\r\n"))
		return
	} else {
		var packed *obj_data = nil
		if GET_OBJ_VNUM(obj) >= 19090 && GET_OBJ_VNUM(obj) <= 0x4A9B || GET_OBJ_VNUM(obj) == 11 {
			act(libc.CString("@CYou push a hidden button on $p@C and a cloud of smoke erupts and covers it. As the smoke clears a small capsule can be seen on the ground.@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("@c$n@C pushes a hidden button on $p@C and a cloud of smokes erupts and covers it. As the smoke clears a small capsule can be seen on the ground.@n"), TRUE, ch, obj, nil, TO_ROOM)
			if GET_OBJ_VNUM(obj) == 11 {
				extract_obj(obj)
				packed = read_object(0x4A8D, VIRTUAL)
				obj_to_room(packed, ch.In_room)
			} else {
				var fnum int = int(GET_OBJ_VNUM(obj) - 10)
				packed = read_object(obj_vnum(fnum), VIRTUAL)
				extract_obj(obj)
				obj_to_room(packed, ch.In_room)
			}
			return
		} else if GET_OBJ_VNUM(obj) >= 18800 && GET_OBJ_VNUM(obj) <= 0x4AFF && int(obj.Type_flag) == ITEM_VEHICLE {
			if arg2[0] == 0 {
				send_to_char(ch, libc.CString("This will sell off your house and delete everything inside. Are you sure? If you are then enter the command again with a yes at the end.\nSyntax: pack (house) yes\r\n"))
				return
			} else if libc.StrCaseCmp(&arg2[0], libc.CString("yes")) != 0 {
				send_to_char(ch, libc.CString("This will sell off your house and delete everything inside. Are you sure? If you are then enter the command again with a yes at the end.\nSyntax: pack (house) yes\r\n"))
				return
			} else if has_housekey(ch, obj) == 0 {
				send_to_char(ch, libc.CString("You do not own this house.\r\n"))
				return
			} else {
				var cont *obj_data = nil
				_ = cont
				act(libc.CString("@CYou push a hidden button on $p@C and a cloud of smoke erupts and covers it. As the smoke clears a pile of money can be seen on the ground!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@c$n@C pushes a hidden button on $p@C and a cloud of smokes erupts and covers it. As the smoke clears a pile of money can be seen on the ground!@n"), TRUE, ch, obj, nil, TO_ROOM)
				var money int = 0
				var count int = 0
				var rnum int = int(GET_OBJ_VNUM(obj))
				if GET_OBJ_VNUM(obj) >= 18800 && GET_OBJ_VNUM(obj) <= 0x49D3 {
					if rnum == 0x4972 {
						rnum = 18800
					} else {
						rnum = rnum - 1
					}
					money = 65000
					for count < 4 {
						for world[real_room(room_vnum(rnum))].Contents != nil {
							extract_obj(world[real_room(room_vnum(rnum))].Contents)
						}
						count++
						rnum++
					}
				} else if GET_OBJ_VNUM(obj) >= 18900 && GET_OBJ_VNUM(obj) <= 0x4A37 {
					rnum = rnum - 1
					money = 150000
					for count < 4 {
						for world[real_room(room_vnum(rnum))].Contents != nil {
							extract_obj(world[real_room(room_vnum(rnum))].Contents)
						}
						count++
						rnum++
					}
				} else if GET_OBJ_VNUM(obj) >= 19100 && GET_OBJ_VNUM(obj) <= 0x4AFF {
					rnum = rnum - 1
					money = 1000000
					for count < 4 {
						for world[real_room(room_vnum(rnum))].Contents != nil {
							extract_obj(world[real_room(room_vnum(rnum))].Contents)
						}
						count++
						rnum++
					}
				}
				var obj2 *obj_data = nil
				var next_obj *obj_data
				_ = next_obj
				for obj2 = ch.Carrying; obj2 != nil; obj2 = obj2.Next_content {
					if GET_OBJ_VNUM(obj) == 0x4972 {
						if GET_OBJ_VNUM(obj2) == 18800 {
							extract_obj(obj2)
						}
					} else {
						if GET_OBJ_VNUM(obj2) == GET_OBJ_VNUM(obj)-1 {
							extract_obj(obj2)
						}
					}
				}
				var money_obj *obj_data = create_money(money)
				obj_to_room(money_obj, ch.In_room)
				extract_obj(obj)
				return
			}
		} else {
			send_to_char(ch, libc.CString("That isn't something you can pack up!\r\n"))
			return
		}
	}
}
func check_insidebag(cont *obj_data, mult float64) int {
	var (
		inside     *obj_data = nil
		next_obj2  *obj_data = nil
		count      int       = 0
		containers int       = 0
	)
	for inside = cont.Contains; inside != nil; inside = next_obj2 {
		next_obj2 = inside.Next_content
		if int(inside.Type_flag) == ITEM_CONTAINER {
			count++
			count += check_insidebag(inside, mult)
			containers++
		} else {
			count++
		}
	}
	count = int(float64(count) * mult)
	count += containers
	return count
}
func check_saveroom_count(ch *char_data, cont *obj_data) int {
	var (
		obj      *obj_data
		next_obj *obj_data = nil
		count    int       = 0
		was      int       = 0
	)
	_ = was
	if ch.In_room == room_rnum(-1) {
		return 0
	} else if !ROOM_FLAGGED(ch.In_room, ROOM_HOUSE) {
		return 0
	}
	for obj = world[ch.In_room].Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		count++
		if !OBJ_FLAGGED(obj, ITEM_CARDCASE) {
			count += check_insidebag(obj, 0.5)
		}
	}
	was = count
	if cont != nil {
		if !OBJ_FLAGGED(cont, ITEM_CARDCASE) {
			count += check_insidebag(cont, 0.5)
		}
		count++
	}
	return count
}
func do_deploy(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj3      *obj_data
		next_obj  *obj_data
		obj4      *obj_data
		obj       *obj_data = nil
		capsule   int       = FALSE
		furniture int       = FALSE
		arg       [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		for obj4 = ch.Carrying; obj4 != nil; obj4 = next_obj {
			next_obj = obj4.Next_content
			if GET_OBJ_VNUM(obj4) == 4 || GET_OBJ_VNUM(obj4) == 5 || GET_OBJ_VNUM(obj4) == 6 {
				obj = obj4
				capsule = TRUE
			}
		}
	} else if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("Syntax: deploy (no argument for houses)\nSyntax: deploy (target) <-- For furniture\r\n"))
		return
	}
	if capsule == FALSE && obj != nil {
		if GET_OBJ_VNUM(obj) >= 19080 && GET_OBJ_VNUM(obj) <= 0x4A9B {
			capsule = TRUE
			furniture = TRUE
		} else {
			send_to_char(ch, libc.CString("That is not a furniture capsule!\r\n"))
			return
		}
	}
	if capsule == FALSE {
		send_to_char(ch, libc.CString("You do not have any house type capsules to deploy.@n\r\n"))
		return
	} else if ch.Rp < 10 && furniture == FALSE {
		send_to_char(ch, libc.CString("You are required to have (not spend) 10 RPP in order to place a house.\r\n"))
		return
	} else if furniture == TRUE && (!ROOM_FLAGGED(ch.In_room, ROOM_HOUSE) || ROOM_FLAGGED(ch.In_room, ROOM_SHIP)) {
		send_to_char(ch, libc.CString("You can't deploy house furniture capsules here.\r\n"))
		return
	} else if furniture == TRUE && (ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2)) {
		send_to_char(ch, libc.CString("You can't deploy house furniture capsules here.\r\n"))
		return
	} else if furniture == FALSE && (SECT(ch.In_room) == SECT_INSIDE || SECT(ch.In_room) == SECT_WATER_NOSWIM || SECT(ch.In_room) == SECT_WATER_SWIM || SECT(ch.In_room) == SECT_SPACE) {
		send_to_char(ch, libc.CString("You can not deploy that in this kind of area. Try an area more suitable for a house.\r\n"))
		return
	}
	if furniture == TRUE {
		var fnum int = 0
		if GET_OBJ_VNUM(obj) == 19080 {
			fnum = 19090
		} else if GET_OBJ_VNUM(obj) == 0x4A89 {
			fnum = 0x4A93
		} else if GET_OBJ_VNUM(obj) == 0x4A8A {
			fnum = 0x4A94
		} else if GET_OBJ_VNUM(obj) == 0x4A8B {
			fnum = 0x4A95
		} else if GET_OBJ_VNUM(obj) == 0x4A8D {
			fnum = 11
		}
		if fnum != 0 {
			var furn *obj_data = read_object(obj_vnum(fnum), VIRTUAL)
			act(libc.CString("@CYou click the capsule's button and toss it to the floor. A puff of smoke erupts immediately and quickly dissipates to reveal, $p@C.@n"), TRUE, ch, furn, nil, TO_CHAR)
			act(libc.CString("@c$n@C clicks a capsule's button and tosses it to the floor. A puff of smoke erupts immediately and quickly dissipates to reveal, $p@C.@n"), TRUE, ch, furn, nil, TO_ROOM)
			obj_to_room(furn, ch.In_room)
			extract_obj(obj)
			return
		} else {
			send_to_imm(libc.CString("ERROR: Furniture failed to deploy at %d."), GET_ROOM_VNUM(ch.In_room))
			return
		}
	}
	var rnum int = 18800
	var giveup int = FALSE
	var cont int = FALSE
	var found int = FALSE
	_ = found
	var type_ int = 0
	if GET_OBJ_VNUM(obj) == 4 {
		type_ = 0
	} else if GET_OBJ_VNUM(obj) == 5 {
		rnum = 18900
		type_ = 1
	} else if GET_OBJ_VNUM(obj) == 6 {
		rnum = 19100
		type_ = 2
	}
	var final int = rnum + 99
	for giveup == FALSE && cont == FALSE {
		obj3 = find_obj_in_list_vnum(world[real_room(room_vnum(rnum))].Contents, 0x4971)
		if obj3 != nil && rnum < final {
			if type_ == 0 {
				rnum += 4
			} else {
				rnum += 5
			}
			found = FALSE
		} else if rnum >= final {
			giveup = TRUE
		} else {
			cont = TRUE
		}
	}
	if cont == TRUE {
		var hnum int = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
		_ = hnum
		var door *obj_data = read_object(0x4971, VIRTUAL)
		door.Value[6] = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
		if rnum != 18800 {
			door.Value[0] = rnum + 1
		} else {
			door.Value[0] = 0x4972
		}
		door.Value[2] = rnum
		obj_to_room(door, real_room(room_vnum(rnum)))
		var key *obj_data = read_object(obj_vnum(rnum), VIRTUAL)
		obj_to_char(key, ch)
		act(libc.CString("@WYou click the capsule and toss it to the ground. A large cloud of smoke erupts from the capsule and after it clears a house is visible in its place!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W clicks a capsule and then tosses it to the ground. A large cloud of smoke erupts from the capsule and after it clears a house is visible in its place!@n"), TRUE, ch, nil, nil, TO_ROOM)
		var foun *obj_data = read_object(0x4973, VIRTUAL)
		obj_to_room(foun, real_room(room_vnum(rnum+1)))
		extract_obj(obj)
	} else {
		send_to_char(ch, libc.CString("@ROOC@D: @wSorry for the inconvenience, but it appears there are no houses available. Please contact Iovan.@n\r\n"))
		return
	}
}
func do_twohand(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are busy grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil || ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are busy struggling with someone!\r\n"))
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) == nil && !PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("You need to wield a sword to use this.\r\n"))
		return
	} else if (ch.Equipment[WEAR_WIELD2]) != nil && !PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("You have something in your offhand already and can't two hand wield your main weapon.\r\n"))
		return
	} else if ((ch.Limb_condition[0]) <= 0 || (ch.Limb_condition[1]) <= 0) && !PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("Kind of hard with only one arm...\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("You stop wielding your weapon with both hands.\r\n"))
		act(libc.CString("$n stops wielding $s weapon with both hands."), TRUE, ch, nil, nil, TO_ROOM)
		REMOVE_BIT_AR(ch.Act[:], PLR_THANDW)
		return
	} else {
		send_to_char(ch, libc.CString("You grab your weapon with both hands.\r\n"))
		act(libc.CString("$n starts wielding $s weapon with both hands."), TRUE, ch, nil, nil, TO_ROOM)
		SET_BIT_AR(ch.Act[:], PLR_THANDW)
		return
	}
}
func start_auction(ch *char_data, obj *obj_data, bid int) {
	obj_from_char(obj)
	obj_selling = obj
	ch_selling = ch
	ch_buying = nil
	curbid = bid
	stdio.Sprintf(&buf[0], "%s magicly flies away from your hands to be auctioned!\r\n", obj_selling.Short_description)
	CAP(&buf[0])
	send_to_char(ch_selling, &buf[0])
	stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_NULL_STATE]), curbid)
	auc_send_to_all(&buf[0], FALSE != 0)
	aucstat = AUC_OFFERING
}
func check_auction() {
	switch aucstat {
	case AUC_NULL_STATE:
		return
	case AUC_OFFERING:
		if obj_selling == nil {
			auc_send_to_all(libc.CString("@RThe auction has stopped because someone has made off with the auctioned object!@n\r\n"), FALSE != 0)
			curbid = 0
			ch_selling = nil
			ch_buying = nil
			aucstat = AUC_NULL_STATE
			return
		}
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_OFFERING]), curbid)
		CAP(&buf[0])
		auc_send_to_all(&buf[0], FALSE != 0)
		aucstat = AUC_GOING_ONCE
		return
	case AUC_GOING_ONCE:
		if obj_selling == nil {
			auc_send_to_all(libc.CString("@RThe auction has stopped because someone has made off with the auctioned object!@n\r\n"), FALSE != 0)
			curbid = 0
			ch_selling = nil
			ch_buying = nil
			aucstat = AUC_NULL_STATE
			return
		}
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_GOING_ONCE]), curbid)
		CAP(&buf[0])
		auc_send_to_all(&buf[0], FALSE != 0)
		aucstat = AUC_GOING_TWICE
		return
	case AUC_GOING_TWICE:
		if obj_selling == nil {
			auc_send_to_all(libc.CString("@RThe auction has stopped because someone has made off with the auctioned object!@n\r\n"), FALSE != 0)
			curbid = 0
			ch_selling = nil
			ch_buying = nil
			aucstat = AUC_NULL_STATE
			return
		}
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_GOING_TWICE]), curbid)
		CAP(&buf[0])
		auc_send_to_all(&buf[0], FALSE != 0)
		aucstat = AUC_LAST_CALL
		return
	case AUC_LAST_CALL:
		if obj_selling == nil {
			auc_send_to_all(libc.CString("@RThe auction has stopped because someone has made off with the auctioned object!@n\r\n"), FALSE != 0)
			curbid = 0
			ch_selling = nil
			ch_buying = nil
			aucstat = AUC_NULL_STATE
			return
		}
		if ch_buying == nil {
			stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_LAST_CALL]), curbid)
			CAP(&buf[0])
			auc_send_to_all(&buf[0], FALSE != 0)
			stdio.Sprintf(&buf[0], "%s flies out the sky and into your hands.\r\n", obj_selling.Short_description)
			CAP(&buf[0])
			send_to_char(ch_selling, &buf[0])
			obj_to_char(obj_selling, ch_selling)
			obj_selling = nil
			ch_selling = nil
			ch_buying = nil
			curbid = 0
			aucstat = AUC_NULL_STATE
			return
		} else {
			stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_SOLD]), curbid)
			auc_send_to_all(&buf[0], TRUE != 0)
			obj_to_char(obj_selling, ch_buying)
			stdio.Sprintf(&buf[0], "%s flies out the sky and into your hands, what a steal!\r\n", obj_selling.Short_description)
			CAP(&buf[0])
			send_to_char(ch_buying, &buf[0])
			stdio.Sprintf(&buf[0], "Congrats! You have sold %s for @Y%d@W zenni!\r\n", obj_selling.Short_description, curbid)
			send_to_char(ch_selling, &buf[0])
			if ch_selling.Gold+curbid > GOLD_CARRY(ch_selling) {
				send_to_char(ch_buying, libc.CString("You couldn't hold all the zenni, so some of it was deposited for you.\r\n"))
				var diff int = 0
				diff = (ch_selling.Gold + curbid) - GOLD_CARRY(ch_selling)
				ch_selling.Gold = GOLD_CARRY(ch_selling)
				ch_selling.Bank_gold += diff
			} else if ch_selling.Gold+curbid <= GOLD_CARRY(ch_selling) {
				ch_selling.Gold += curbid
			}
			obj_selling = nil
			ch_selling = nil
			ch_buying = nil
			curbid = 0
			aucstat = AUC_NULL_STATE
			return
		}
	}
}
func dball_load() {
	var (
		found1 int = FALSE
		found2 int = FALSE
		found3 int = FALSE
		found4 int = FALSE
		found5 int = FALSE
		load   int = FALSE
		num    int = -1
		found6 int = FALSE
		found7 int = FALSE
		room   int = 0
		loaded int = FALSE
	)
	_ = loaded
	var hunter1 int = FALSE
	var hunter2 int = FALSE
	var k *obj_data
	if SELFISHMETER >= 10 {
		return
	}
	if dballtime == 0 {
		var (
			hunter *char_data = nil
			r_num  mob_rnum
		)
		WISHTIME = 0
		for k = object_list; k != nil; k = k.Next {
			if OBJ_FLAGGED(k, ITEM_FORGED) {
				continue
			}
			if GET_OBJ_VNUM(k) == 20 {
				found1 = TRUE
			} else if GET_OBJ_VNUM(k) == 21 {
				found2 = TRUE
			} else if GET_OBJ_VNUM(k) == 22 {
				found3 = TRUE
			} else if GET_OBJ_VNUM(k) == 23 {
				found4 = TRUE
			} else if GET_OBJ_VNUM(k) == 24 {
				found5 = TRUE
			} else if GET_OBJ_VNUM(k) == 25 {
				found6 = TRUE
			} else if GET_OBJ_VNUM(k) == 26 {
				found7 = TRUE
			} else if k.In_room != room_rnum(-1) && world[k.In_room].Geffect == 6 && !OBJ_FLAGGED(k, ITEM_UNBREAKABLE) {
				send_to_room(k.In_room, libc.CString("@R%s@r melts in the lava!@n\r\n"), k.Short_description)
				extract_obj(k)
			} else {
				continue
			}
		}
		if found1 == FALSE {
			load = FALSE
			var zone int = 0
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if (func() int {
						zone = int(real_zone_by_thing(room_vnum(real_room(room_vnum(num)))))
						return zone
					}()) != int(-1) {
						if ZONE_FLAGGED(zone_rnum(zone), ZONE_DBALLS) {
							room = num
							load = TRUE
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(200, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(20, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(20, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(20, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(20, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found2 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(21, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(21, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(21, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(21, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found3 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(22, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(22, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(22, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(22, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found4 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(23, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(23, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(23, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(23, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found5 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(24, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(24, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(24, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(24, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found6 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(25, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(25, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(25, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(25, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		if found7 == FALSE {
			load = FALSE
			for load == FALSE {
				if real_room(room_vnum(num)) != room_rnum(-1) {
					if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
						room = num
						load = TRUE
						num = rand_number(200, 20000)
					} else {
						num = rand_number(200, 20000)
					}
				} else {
					num = rand_number(20, 20000)
				}
			}
			if rand_number(1, 10) > 8 {
				if hunter1 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER1_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter1 = TRUE
					DBALL_HUNTER1 = room
					k = read_object(26, VIRTUAL)
					obj_to_char(k, hunter)
				} else if hunter2 == FALSE {
					if (func() mob_rnum {
						r_num = real_mobile(mob_vnum(DBALL_HUNTER2_VNUM))
						return r_num
					}()) == mob_rnum(-1) {
						return
					}
					hunter = read_mobile(mob_vnum(r_num), REAL)
					char_to_room(hunter, real_room(room_vnum(room)))
					hunter2 = TRUE
					DBALL_HUNTER2 = room
					k = read_object(26, VIRTUAL)
					obj_to_char(k, hunter)
				} else {
					k = read_object(26, VIRTUAL)
					obj_to_room(k, real_room(room_vnum(room)))
				}
			} else {
				k = read_object(26, VIRTUAL)
				obj_to_room(k, real_room(room_vnum(room)))
			}
			loaded = TRUE
		}
		dballtime = 604800
	} else if dballtime == 0x7E900 || dballtime == 432000 || dballtime == 0x54600 || dballtime == 259200 || dballtime == 0x2A300 || dballtime == 86400 {
		dballtime -= 1
	} else {
		if WISHTIME == 0 {
			WISHTIME = dballtime - 1
		} else if WISHTIME > 0 && dballtime != WISHTIME {
			dballtime = WISHTIME
		}
		WISHTIME -= 1
		dballtime -= 1
	}
}
func do_auction(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1 [2048]byte
		arg2 [2048]byte
		obj  *obj_data
		bid  int = 0
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
		send_to_char(ch, libc.CString("This is a different dimension!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
		send_to_char(ch, libc.CString("You are in the past!\r\n"))
		return
	}
	if PRF_FLAGGED(ch, PRF_HIDE) {
		send_to_char(ch, libc.CString("The auctioneer will not accept items from hidden people.\r\n"))
		return
	}
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("Auction what?\r\n"))
		send_to_char(ch, libc.CString("[ Auction: <item> | <cancel> ]\r\n"))
		return
	} else if is_abbrev(&arg1[0], libc.CString("cancel")) != 0 || is_abbrev(&arg1[0], libc.CString("stop")) != 0 {
		if ch != ch_selling && ch.Admlevel <= ADMLVL_GRGOD || aucstat == AUC_NULL_STATE {
			send_to_char(ch, libc.CString("You're not even selling anything!\r\n"))
			return
		} else if ch == ch_selling {
			stop_auction(AUC_NORMAL_CANCEL, nil)
			return
		} else {
			stop_auction(AUC_WIZ_CANCEL, ch)
		}
	} else if is_abbrev(&arg1[0], libc.CString("stats")) != 0 || is_abbrev(&arg1[0], libc.CString("identify")) != 0 {
		auc_stat(ch, obj_selling)
		return
	} else if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		stdio.Sprintf(&buf[0], "You don't seem to have %s %s.\r\n", AN(&arg1[0]), &arg1[0])
		send_to_char(ch, &buf[0])
		return
	} else if arg2[0] == 0 {
		stdio.Sprintf(&buf[0], "What should be the minimum bid?\r\n")
		send_to_char(ch, &buf[0])
		return
	} else if arg2[0] != 0 && (func() int {
		bid = libc.Atoi(libc.GoString(&arg2[0]))
		return bid
	}()) <= 0 {
		send_to_char(ch, libc.CString("Come on? One zenni at least?\r\n"))
		return
	} else if aucstat != AUC_NULL_STATE {
		stdio.Sprintf(&buf[0], "Sorry, but %s is already auctioning %s at @Y%d@W zenni!\r\n", GET_NAME(ch_selling), obj_selling.Short_description, bid)
		send_to_char(ch, &buf[0])
		return
	} else if OBJ_FLAGGED(obj, ITEM_NOSELL) {
		send_to_char(ch, libc.CString("Sorry but you can't sell that!\r\n"))
		return
	} else if (obj.Value[VAL_CONTAINER_CORPSE]) == 1 {
		send_to_char(ch, libc.CString("Sorry but you can't sell that!\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("Ok.\r\n"))
		start_auction(ch, obj, bid)
		return
	}
}
func do_bid(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj        *obj_data
		next_obj   *obj_data
		obj2       *obj_data = nil
		arg        [2048]byte
		arg2       [2048]byte
		found      int = FALSE
		list       int = 0
		masterList int = 0
		auct_room  room_vnum
	)
	auct_room = room_vnum(real_room(80))
	if IS_NPC(ch) {
		return
	}
	if (ch.Equipment[WEAR_EYE]) == nil {
		send_to_char(ch, libc.CString("You need a scouter to make an auction bid.\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
		send_to_char(ch, libc.CString("This is a different dimension!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
		send_to_char(ch, libc.CString("This is the past, nothing is being auctioned!\r\n"))
		return
	}
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: bid [ list | # ] (amt)\r\nOr...\r\nSyntax: bid appraise (list number)\r\n"))
		return
	}
	for obj = world[auct_room].Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if obj != nil {
			list++
		}
	}
	masterList = list
	list = 0
	if libc.StrCaseCmp(&arg[0], libc.CString("list")) == 0 {
		send_to_char(ch, libc.CString("@Y                                   Auction@n\r\n"))
		send_to_char(ch, libc.CString("@c------------------------------------------------------------------------------@n\r\n"))
		for obj = world[auct_room].Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil {
				if int(obj.Aucter) <= 0 {
					continue
				}
				list++
				if obj.AucTime+86400 > libc.GetTime(nil) && int(obj.CurBidder) <= -1 {
					send_to_char(ch, libc.CString("@D[@R#@W%3d@D][@mOwner@W: @w%10s@D][@GItem Name@W: @w%-*s@D][@GCost@W: @Y%s@D]@n\r\n"), list, func() *byte {
						if get_name_by_id(int(obj.Aucter)) != nil {
							return CAP(get_name_by_id(int(obj.Aucter)))
						}
						return libc.CString("Nobody")
					}(), count_color_chars(obj.Short_description)+30, obj.Short_description, add_commas(int64(obj.Bid)))
				} else if obj.AucTime+86400 > libc.GetTime(nil) && int(obj.CurBidder) > -1 {
					send_to_char(ch, libc.CString("@D[@R#@W%3d@D][@mOwner@W: @w%10s@D][@GItem Name@W: @w%-*s@D][@RTop Bid@W: %s @Y%s@D]@n\r\n"), list, func() *byte {
						if get_name_by_id(int(obj.Aucter)) != nil {
							return CAP(get_name_by_id(int(obj.Aucter)))
						}
						return libc.CString("Nobody")
					}(), count_color_chars(obj.Short_description)+30, obj.Short_description, func() *byte {
						if get_name_by_id(int(obj.CurBidder)) != nil {
							return CAP(get_name_by_id(int(obj.CurBidder)))
						}
						return libc.CString("Nobody")
					}(), add_commas(int64(obj.Bid)))
				} else if obj.AucTime+86400 < libc.GetTime(nil) && int(obj.CurBidder) > -1 {
					send_to_char(ch, libc.CString("@D[@R#@W%3d@D][@mOwner@W: @w%10s@D][@GItem Name@W: @w%-*s@D][@RBid Winner@W: %s @Y%s@D]@n\r\n"), list, func() *byte {
						if get_name_by_id(int(obj.Aucter)) != nil {
							return CAP(get_name_by_id(int(obj.Aucter)))
						}
						return libc.CString("Nobody")
					}(), count_color_chars(obj.Short_description)+30, obj.Short_description, func() *byte {
						if get_name_by_id(int(obj.CurBidder)) != nil {
							return CAP(get_name_by_id(int(obj.CurBidder)))
						}
						return libc.CString("Nobody")
					}(), add_commas(int64(obj.Bid)))
				} else {
					send_to_char(ch, libc.CString("@D[@R#@W%3d@D][@mOwner@W: @w%10s@D][@GItem Name@W: @w%-*s@D][@RClosed@D]@n\r\n"), list, func() *byte {
						if get_name_by_id(int(obj.Aucter)) != nil {
							return CAP(get_name_by_id(int(obj.Aucter)))
						}
						return libc.CString("Nobody")
					}(), count_color_chars(obj.Short_description)+30, obj.Short_description)
				}
				found = TRUE
			}
		}
		if found == FALSE {
			send_to_char(ch, libc.CString("No items are currently being auctioned.\r\n"))
		}
		send_to_char(ch, libc.CString("@c------------------------------------------------------------------------------@n\r\n"))
	} else if libc.StrCaseCmp(&arg[0], libc.CString("appraise")) == 0 {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: bid [ list | # ] (amt)\r\nOr...\r\nSyntax: bid appraise (list number)\r\n"))
			send_to_char(ch, libc.CString("What item number did you want to appraise?\r\n"))
			return
		} else if libc.Atoi(libc.GoString(&arg2[0])) < 0 || libc.Atoi(libc.GoString(&arg2[0])) > masterList {
			send_to_char(ch, libc.CString("Syntax: bid [ list | # ] (amt)\r\nOr...\r\nSyntax: bid appraise (list number)\r\n"))
			send_to_char(ch, libc.CString("That item number doesn't exist.\r\n"))
			return
		}
		for obj = world[auct_room].Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil {
				if int(obj.Aucter) <= 0 {
					continue
				}
				list++
				if libc.Atoi(libc.GoString(&arg2[0])) == list {
					obj2 = obj
				}
			}
		}
		if obj2 == nil {
			send_to_char(ch, libc.CString("That item number is not found.\r\n"))
			return
		} else {
			if GET_SKILL(ch, SKILL_APPRAISE) == 0 {
				send_to_char(ch, libc.CString("You are unskilled at appraising.\r\n"))
				return
			}
			improve_skill(ch, SKILL_APPRAISE, 1)
			if GET_SKILL(ch, SKILL_APPRAISE) < rand_number(1, 101) {
				send_to_char(ch, libc.CString("You look at the images for %s and fail to perceive its worth..\r\n"), obj2.Short_description)
				act(libc.CString("@c$n@w looks stumped about something they viewed on their scouter screen.@n"), TRUE, ch, nil, nil, TO_ROOM)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			} else {
				send_to_char(ch, libc.CString("You look at images of the object on your scouter.\r\n"))
				act(libc.CString("@c$n@w looks at something on their scouter screen.@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_char(ch, libc.CString("@c------------------------------------------------------------------------\n"))
				send_to_char(ch, libc.CString("@GOwner       @W: @w%s@n\n"), func() *byte {
					if get_name_by_id(int(obj2.Aucter)) != nil {
						return CAP(get_name_by_id(int(obj2.Aucter)))
					}
					return libc.CString("Nobody")
				}())
				send_to_char(ch, libc.CString("@GItem Name   @W: @w%s@n\n"), obj2.Short_description)
				send_to_char(ch, libc.CString("@GCurrent Bid @W: @Y%s@n\n"), add_commas(int64(obj2.Bid)))
				send_to_char(ch, libc.CString("@GStore Value @W: @Y%s@n\n"), add_commas(int64(obj2.Cost)))
				send_to_char(ch, libc.CString("@GItem Min LVL@W: @w%d@n\n"), obj2.Level)
				if (obj2.Value[VAL_ALL_HEALTH]) >= 100 {
					send_to_char(ch, libc.CString("@GCondition   @W: @C%d%s@n\n"), obj2.Value[VAL_ALL_HEALTH], "%")
				} else if (obj2.Value[VAL_ALL_HEALTH]) >= 50 {
					send_to_char(ch, libc.CString("@GCondition   @W: @y%d%s@n\n"), obj2.Value[VAL_ALL_HEALTH], "%")
				} else if (obj2.Value[VAL_ALL_HEALTH]) >= 1 {
					send_to_char(ch, libc.CString("@GCondition   @W: @r%d%s@n\n"), obj2.Value[VAL_ALL_HEALTH], "%")
				} else {
					send_to_char(ch, libc.CString("@GCondition   @W: @D%d%s@n\n"), obj2.Value[VAL_ALL_HEALTH], "%")
				}
				send_to_char(ch, libc.CString("@GItem Weight @W: @w%s@n\n"), add_commas(obj2.Weight))
				var bits [64936]byte
				sprintbitarray(obj2.Wear_flags[:], wear_bits[:], TW_ARRAY_MAX, &bits[0])
				search_replace(&bits[0], libc.CString("TAKE"), libc.CString(""))
				send_to_char(ch, libc.CString("@GWear Loc.   @W:@w%s\n"), &bits[0])
				if int(obj2.Type_flag) == ITEM_WEAPON {
					if OBJ_FLAGGED(obj2, ITEM_WEAPLVL1) {
						send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w5%s@D]@n\r\n"), "%")
					} else if OBJ_FLAGGED(obj2, ITEM_WEAPLVL2) {
						send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w10%s@D]@n\r\n"), "%")
					} else if OBJ_FLAGGED(obj2, ITEM_WEAPLVL3) {
						send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w20%s@D]@n\r\n"), "%")
					} else if OBJ_FLAGGED(obj2, ITEM_WEAPLVL4) {
						send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w30%s@D]@n\r\n"), "%")
					} else if OBJ_FLAGGED(obj2, ITEM_WEAPLVL5) {
						send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w50%s@D]@n\r\n"), "%")
					}
				}
				var i int
				var found int = FALSE
				send_to_char(ch, libc.CString("@GItem Size   @W:@w %s@n\r\n"), size_names[obj2.Size])
				send_to_char(ch, libc.CString("@GItem Bonuses@W:@w"))
				for i = 0; i < MAX_OBJ_AFFECT; i++ {
					if obj2.Affected[i].Modifier != 0 {
						sprinttype(obj2.Affected[i].Location, apply_types[:], &buf[0], uint64(64936))
						send_to_char(ch, libc.CString("%s %+d to %s"), func() string {
							if func() int {
								p := &found
								x := *p
								*p++
								return x
							}() != 0 {
								return ","
							}
							return ""
						}(), obj2.Affected[i].Modifier, &buf[0])
						switch obj2.Affected[i].Location {
						case APPLY_FEAT:
							send_to_char(ch, libc.CString(" (%s)"), feat_list[obj2.Affected[i].Specific].Name)
						case APPLY_SKILL:
							send_to_char(ch, libc.CString(" (%s)"), spell_info[obj2.Affected[i].Specific].Name)
						}
					}
				}
				if found == 0 {
					send_to_char(ch, libc.CString(" None@n"))
				} else {
					send_to_char(ch, libc.CString("@n"))
				}
				var buf2 [64936]byte
				sprintbitarray(obj2.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &buf2[0])
				send_to_char(ch, libc.CString("\n@GSpecial     @W:@w %s\n"), &buf2[0])
				send_to_char(ch, libc.CString("@c------------------------------------------------------------------------\n"))
				return
			}
		}
	} else {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: bid [ list | # ] (amt)\r\nOr...\r\nSyntax: bid appraise (list number)\r\n"))
			send_to_char(ch, libc.CString("What amount did you want to bid?\r\n"))
			return
		} else if libc.Atoi(libc.GoString(&arg[0])) < 0 || libc.Atoi(libc.GoString(&arg[0])) > masterList {
			send_to_char(ch, libc.CString("Syntax: bid [ list | # ] (amt)\r\nOr...\r\nSyntax: bid appraise (list number)\r\n"))
			send_to_char(ch, libc.CString("That item number is not found.\r\n"))
			return
		}
		for obj = world[auct_room].Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil {
				if int(obj.Aucter) <= 0 {
					continue
				}
				list++
				if libc.Atoi(libc.GoString(&arg[0])) == list {
					obj2 = obj
				}
			}
		}
		if obj2 == nil {
			send_to_char(ch, libc.CString("That item number is not found.\r\n"))
			return
		} else if int(obj2.CurBidder) == int(ch.Id) {
			send_to_char(ch, libc.CString("You already have the highest bid.\r\n"))
			return
		} else if int(obj2.Aucter) == int(ch.Id) {
			send_to_char(ch, libc.CString("You auctioned the item, go to the auction house and cancel if you can.\r\n"))
			return
		} else if int(obj2.CurBidder) > 0 && float64(libc.Atoi(libc.GoString(&arg2[0]))) <= (float64(obj2.Bid)+float64(obj2.Bid)*0.1) && int(obj2.CurBidder) > -1 {
			send_to_char(ch, libc.CString("You have to bid at least 10 percent over the current bid.\r\n"))
			return
		} else if libc.Atoi(libc.GoString(&arg2[0])) < obj2.Bid && int(obj2.CurBidder) <= -1 {
			send_to_char(ch, libc.CString("You have to bid at least the starting bid.\r\n"))
			return
		} else if libc.Atoi(libc.GoString(&arg2[0])) > (((ch.Gold+ch.Bank_gold)/100)*50)+(ch.Gold+ch.Bank_gold) {
			send_to_char(ch, libc.CString("You can not bid more than 150%s of your total money (on hand and in the bank).\r\n"), "%")
			return
		} else if obj2.AucTime+86400 <= libc.GetTime(nil) {
			send_to_char(ch, libc.CString("Bidding on that object has been closed.\r\n"))
			return
		} else {
			obj2.Bid = libc.Atoi(libc.GoString(&arg2[0]))
			obj2.CurBidder = ch.Id
			auc_save()
			var d *descriptor_data
			var bid int = libc.Atoi(libc.GoString(&arg2[0]))
			basic_mud_log(libc.CString("AUCTION: %s has bid %s on %s"), GET_NAME(ch), obj2.Short_description, add_commas(int64(bid)))
			for d = descriptor_list; d != nil; d = d.Next {
				if d.Connected != CON_PLAYING || IS_NPC(d.Character) {
					continue
				}
				if d.Character == ch {
					if (d.Character.Equipment[WEAR_EYE]) != nil {
						send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GYou have bid @Y%s@G on @w%s@G@n\r\n"), add_commas(int64(obj2.Bid)), obj2.Short_description)
					}
					continue
				}
				if (d.Character.Equipment[WEAR_EYE]) != nil {
					send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GThe bid on, @w%s@G, has been raised to @Y%s@n\r\n"), obj2.Short_description, add_commas(int64(obj2.Bid)))
				}
			}
		}
	}
}
func stop_auction(type_ int, ch *char_data) {
	if obj_selling == nil {
		auc_send_to_all(libc.CString("@RThe auction has stopped because someone has made off with the auctioned object!@n\r\n"), FALSE != 0)
		curbid = 0
		ch_selling = nil
		ch_buying = nil
		aucstat = AUC_NULL_STATE
		return
	}
	switch type_ {
	case AUC_NORMAL_CANCEL:
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_NORMAL_CANCEL]))
		auc_send_to_all(&buf[0], FALSE != 0)
	case AUC_QUIT_CANCEL:
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_QUIT_CANCEL]))
		auc_send_to_all(&buf[0], FALSE != 0)
	case AUC_WIZ_CANCEL:
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_WIZ_CANCEL]))
		auc_send_to_all(&buf[0], FALSE != 0)
	default:
		send_to_char(ch, libc.CString("Sorry, that is an unrecognised cancel command, please report."))
		return
	}
	if type_ != AUC_WIZ_CANCEL {
		stdio.Sprintf(&buf[0], "%s flies out the sky and into your hands.\r\n", obj_selling.Short_description)
		CAP(&buf[0])
		send_to_char(ch_selling, &buf[0])
		obj_to_char(obj_selling, ch_selling)
	} else {
		stdio.Sprintf(&buf[0], "%s flies out the sky and into your hands.\r\n", obj_selling.Short_description)
		CAP(&buf[0])
		send_to_char(ch, &buf[0])
		obj_to_char(obj_selling, ch)
	}
	if ch_buying != nil {
		ch_buying.Gold += curbid
	}
	obj_selling = nil
	ch_selling = nil
	ch_buying = nil
	curbid = 0
	aucstat = AUC_NULL_STATE
}
func auc_stat(ch *char_data, obj *obj_data) {
	if aucstat == AUC_NULL_STATE {
		send_to_char(ch, libc.CString("Nothing is being auctioned!\r\n"))
		return
	} else if ch == ch_selling {
		send_to_char(ch, libc.CString("You should have found that out BEFORE auctioning it!\r\n"))
		return
	} else if ch.Gold < 500 {
		send_to_char(ch, libc.CString("You can't afford to find the stats on that, it costs 500 zenni!\r\n"))
		return
	} else {
		stdio.Sprintf(&buf[0], libc.GoString(auctioneer[AUC_STAT]), curbid)
		act(&buf[0], TRUE, ch_selling, obj, unsafe.Pointer(ch), int(TO_VICT|2<<7))
		ch.Gold -= 500
	}
}
func auc_send_to_all(messg *byte, buyer bool) {
	var i *descriptor_data
	if messg == nil {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		if ROOM_FLAGGED(i.Character.In_room, ROOM_HBTC) {
			continue
		}
		if ROOM_FLAGGED(i.Character.In_room, ROOM_PAST) {
			continue
		}
		if buyer {
			act(messg, TRUE, ch_buying, obj_selling, unsafe.Pointer(i.Character), int(TO_VICT|2<<7))
		} else {
			act(messg, TRUE, ch_selling, obj_selling, unsafe.Pointer(i.Character), int(TO_VICT|2<<7))
		}
	}
}
func lambda_toolsearch(obj *obj_data) bool {
	return GET_OBJ_VNUM(obj) == 386 && (obj.Value[VAL_ALL_HEALTH]) > 0
}
func do_assemble(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		lVnum   int       = int(-1)
		pObject *obj_data = nil
		buf     [64936]byte
		roll    int = 0
	)
	skip_spaces(&argument)
	if *argument == '\x00' {
		send_to_char(ch, libc.CString("What would you like to %s?\r\n"), complete_cmd_info[cmd].Command)
		return
	} else if (func() int {
		lVnum = assemblyFindAssembly(argument)
		return lVnum
	}()) < 0 {
		send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
		return
	} else if assemblyGetType(lVnum) != subcmd {
		send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
		return
	} else if !assemblyCheckComponents(lVnum, ch, FALSE) {
		send_to_char(ch, libc.CString("You haven't got all the things you need.\r\n"))
		return
	} else if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		send_to_char(ch, libc.CString("You can't do that in space."))
		return
	} else if GET_SKILL(ch, SKILL_SURVIVAL) == 0 && libc.StrCaseCmp(argument, libc.CString("campfire")) == 0 {
		send_to_char(ch, libc.CString("You know nothing about building campfires.\r\n"))
		return
	}
	if libc.StrStr(argument, libc.CString("Signal")) != nil || libc.StrStr(argument, libc.CString("signal")) != nil {
		if GET_SKILL(ch, SKILL_BUILD) < 70 {
			send_to_char(ch, libc.CString("You need at least a build skill level of 70.\r\n"))
			return
		}
	}
	var tool *obj_data = find_obj_in_list_lambda(ch.Carrying, lambda_toolsearch)
	if tool != nil {
		act(libc.CString("@WYou open up your toolkit and take out the necessary tools.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W opens up $s toolkit and takes out the necessary tools.@n"), TRUE, ch, nil, nil, TO_ROOM)
	} else {
		send_to_char(ch, libc.CString("You wish you had tools, but make the best out of what you do have anyway...\r\n"))
		roll = 20
	}
	var survival int = GET_SKILL(ch, SKILL_SURVIVAL)
	if libc.StrCaseCmp(argument, libc.CString("campfire")) != 0 {
		if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) || SECT(ch.In_room) == SECT_WATER_NOSWIM || SUNKEN(ch.In_room) {
			send_to_char(ch, libc.CString("This area will not allow a fire to burn properly.\r\n"))
			return
		}
		if survival >= 90 {
			roll += axion_dice(0)
		} else if survival < 90 {
			roll += axion_dice(0)
		}
		improve_skill(ch, SKILL_BUILD, 1)
		if GET_SKILL(ch, SKILL_BUILD) <= roll {
			if (func() *obj_data {
				pObject = read_object(obj_vnum(lVnum), VIRTUAL)
				return pObject
			}()) == nil {
				send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
				return
			}
			extract_obj(pObject)
			send_to_char(ch, libc.CString("You start to %s %s %s, but mess up royally!\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
			if assemblyCheckComponents(lVnum, ch, TRUE) {
				roll = 9001
			}
			return
		}
	} else {
		if survival >= 90 {
			roll += axion_dice(0)
		} else if survival < 90 {
			roll += axion_dice(-10)
		}
		improve_skill(ch, SKILL_BUILD, 1)
		if survival <= roll {
			if (func() *obj_data {
				pObject = read_object(obj_vnum(lVnum), VIRTUAL)
				return pObject
			}()) == nil {
				send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
				return
			}
			extract_obj(pObject)
			send_to_char(ch, libc.CString("You start to %s %s %s, but mess up royally!\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
			if tool != nil && rand_number(1, 3) == 3 && (tool.Value[VAL_ALL_HEALTH]) > 0 {
				tool.Value[VAL_ALL_HEALTH] -= rand_number(1, 5)
				act(libc.CString("@RYour toolset is looking a bit more worn.@n"), TRUE, ch, nil, nil, TO_CHAR)
				if (tool.Value[VAL_ALL_HEALTH]) <= 0 {
					tool.Value[VAL_ALL_HEALTH] = 0
				}
			}
			if assemblyCheckComponents(lVnum, ch, TRUE) {
				roll = 9001
			}
			return
		}
	}
	if axion_dice(0)-int(ch.Aff_abils.Intel)/5 > 95 {
		if (func() *obj_data {
			pObject = read_object(obj_vnum(lVnum), VIRTUAL)
			return pObject
		}()) == nil {
			send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
			return
		}
		extract_obj(pObject)
		send_to_char(ch, libc.CString("You start to %s %s %s, but forget a couple of steps. You take it apart and give up.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		return
	} else if rand_number(1, 100) >= 92 {
		if (func() *obj_data {
			pObject = read_object(obj_vnum(lVnum), VIRTUAL)
			return pObject
		}()) == nil {
			send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
			return
		}
		extract_obj(pObject)
		send_to_char(ch, libc.CString("You start to %s %s %s, but put it together wrong and have to stop. You take it apart and give up.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		return
	}
	if (func() *obj_data {
		pObject = read_object(obj_vnum(lVnum), VIRTUAL)
		return pObject
	}()) == nil {
		send_to_char(ch, libc.CString("You can't %s %s %s.\r\n"), complete_cmd_info[cmd].Command, AN(argument), argument)
		return
	}
	add_unique_id(pObject)
	if GET_OBJ_VNUM(pObject) != 1611 {
		obj_to_char(pObject, ch)
		if int(ch.Race) == RACE_TRUFFLE {
			var (
				count    int = 0
				plused   int = FALSE
				hasstat  int = 0
				failsafe int = 0
			)
			for count < 6 {
				if pObject.Affected[count].Location > 0 && rand_number(1, 6) <= 2 && plused == FALSE {
					pObject.Affected[count].Modifier += rand_number(1, 3)
					plused = TRUE
				} else if pObject.Affected[count].Location > 0 {
					hasstat += 1
				}
				failsafe++
				count++
				if plused == TRUE {
					send_to_char(ch, libc.CString("@YYour intuitive skill with building has made this item even better!@n\r\n"))
					count = 6
				} else if failsafe >= 12 {
					send_to_char(ch, libc.CString("@yIt seems this item could not be upgraded with your truffle knowledge...@n\r\n"))
					count = 6
				} else if failsafe == 11 && plused == FALSE {
					send_to_char(ch, libc.CString("@YYour intuitive skill with building has made this item even better!@n\r\n"))
					pObject.Affected[count].Location = rand_number(1, 6)
					pObject.Affected[count].Modifier += rand_number(1, 3)
					plused = TRUE
					count = 6
				} else if count == 6 && hasstat > 0 {
					count = 0
				}
			}
		}
	} else {
		obj_to_room(pObject, ch.In_room)
		pObject.Timer = int(float64(GET_SKILL(ch, SKILL_SURVIVAL)) * 0.12)
	}
	stdio.Sprintf(&buf[0], "You %s $p.", complete_cmd_info[cmd].Command)
	act(&buf[0], FALSE, ch, pObject, nil, TO_CHAR)
	stdio.Sprintf(&buf[0], "$n %ss $p.", complete_cmd_info[cmd].Command)
	act(&buf[0], FALSE, ch, pObject, nil, TO_ROOM)
	if assemblyCheckComponents(lVnum, ch, TRUE) {
		roll = 9001
	}
	if int(ch.Race) != RACE_TRUFFLE && axion_dice(8) > GET_SKILL(ch, SKILL_BUILD) {
		send_to_char(ch, libc.CString("@yYou've made an inferior product. Its value will be somewhat less.@n\r\n"))
		pObject.Cost -= int(float64(pObject.Cost) * 0.25)
	} else if int(ch.Race) == RACE_TRUFFLE && axion_dice(18) > GET_SKILL(ch, SKILL_BUILD) {
		send_to_char(ch, libc.CString("@yYou've made an inferior product. Its value will be somewhat less.@n\r\n"))
		pObject.Cost -= int(float64(pObject.Cost) * 0.12)
	} else if int(libc.BoolToInt(int(ch.Race) == RACE_TRUFFLE)) < GET_SKILL(ch, SKILL_BUILD) {
		send_to_char(ch, libc.CString("@YYou've made an excellent product. Its value will be somewhat more.@n\r\n"))
		pObject.Cost += int(float64(pObject.Cost) * 0.12)
	}
	if int(ch.Race) == RACE_TRUFFLE && rand_number(1, 5) >= 4 && pObject.Cost >= 500 {
		if GET_LEVEL(ch) < 100 && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
			var gain int64 = int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.011)
			send_to_char(ch, libc.CString("@RExp Bonus@D: @G%s@n\r\n"), add_commas(gain))
			gain_exp(ch, gain)
		} else {
			gain_exp(ch, 1375000)
			send_to_char(ch, libc.CString("@RExp Bonus@D: @G%s@n\r\n"), add_commas(1375000))
		}
	}
}
func perform_put(ch *char_data, obj *obj_data, cont *obj_data) {
	var dball [7]int = [7]int{20, 21, 22, 23, 24, 25, 26}
	if drop_otrigger(obj, ch) == 0 {
		return
	}
	if obj == nil {
		return
	}
	if OBJ_FLAGGED(cont, ITEM_FORGED) {
		act(libc.CString("$P is forged and won't hold anything."), FALSE, ch, nil, unsafe.Pointer(cont), TO_CHAR)
		return
	}
	if OBJ_FLAGGED(cont, ITEM_SHEATH) && int(obj.Type_flag) != ITEM_WEAPON {
		send_to_char(ch, libc.CString("That is made to only hold weapons.\r\n"))
		return
	}
	if OBJ_FLAGGED(cont, ITEM_SHEATH) {
		var (
			obj2     *obj_data = nil
			next_obj *obj_data = nil
		)
		_ = next_obj
		var count int = 0
		var minus int = 0
		for obj2 = cont.Contains; obj2 != nil; obj2 = obj2.Next_content {
			minus += int(obj2.Weight)
			count++
		}
		obj2 = nil
		var holds int = int(cont.Weight - int64(minus))
		if count >= holds {
			send_to_char(ch, libc.CString("It can only hold %d weapon%s at a time.\r\n"), holds, func() string {
				if holds > 1 {
					return "s"
				}
				return ""
			}())
			return
		}
	}
	if int(cont.Type_flag) == ITEM_CONTAINER && (cont.Value[VAL_CONTAINER_CAPACITY]) == 0 {
		act(libc.CString("$p won't fit in $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
	} else if GET_OBJ_VNUM(cont) >= 600 && GET_OBJ_VNUM(cont) <= 603 {
		send_to_char(ch, libc.CString("You can't put cards on a duel table. You have to @Gplay@n them.\r\n"))
	} else if (GET_OBJ_VNUM(cont) == 697 || GET_OBJ_VNUM(cont) == 698 || GET_OBJ_VNUM(cont) == 682 || GET_OBJ_VNUM(cont) == 683 || GET_OBJ_VNUM(cont) == 684 || OBJ_FLAGGED(cont, ITEM_CARDCASE)) && !OBJ_FLAGGED(obj, ITEM_ANTI_HIEROPHANT) {
		send_to_char(ch, libc.CString("You can only put cards in a case.\r\n"))
	} else if int(cont.Type_flag) == ITEM_CONTAINER && (cont.Value[VAL_CONTAINER_CAPACITY]) > 0 && cont.Weight+obj.Weight > int64(cont.Value[VAL_CONTAINER_CAPACITY]) {
		act(libc.CString("$p won't fit in $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
	} else if OBJ_FLAGGED(obj, ITEM_NODROP) && cont.In_room != room_rnum(-1) {
		act(libc.CString("You can't get $p out of your hand."), FALSE, ch, obj, nil, TO_CHAR)
	} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[0]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[1]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[2]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[3]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[4]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[5]) || GET_OBJ_VNUM(obj) == obj_vnum(dball[6]) {
		send_to_char(ch, libc.CString("You can not bag dragon balls.\r\n"))
	} else if OBJ_FLAGGED(obj, ITEM_NORENT) {
		send_to_char(ch, libc.CString("That isn't worth bagging. Better keep that close if you wanna keep it at all.\r\n"))
	} else if cont.Carried_by == nil && check_saveroom_count(ch, obj) > 150 {
		send_to_char(ch, libc.CString("The save room can not hold anymore items. (150 max, count of items in containers is halved)\r\n"))
	} else {
		obj_from_char(obj)
		obj_to_obj(obj, cont)
		if !OBJ_FLAGGED(obj, ITEM_ANTI_HIEROPHANT) {
			act(libc.CString("$n puts $p in $P."), TRUE, ch, obj, unsafe.Pointer(cont), TO_ROOM)
		} else {
			act(libc.CString("$n puts an @DA@wd@cv@Ce@Wnt @DD@wu@ce@Cl @mC@Ma@Wr@wd@n in $P."), TRUE, ch, obj, unsafe.Pointer(cont), TO_ROOM)
		}
		if OBJ_FLAGGED(obj, ITEM_NODROP) && !OBJ_FLAGGED(cont, ITEM_NODROP) {
			SET_BIT_AR(cont.Extra_flags[:], ITEM_NODROP)
			act(libc.CString("You get a strange feeling as you put $p in $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
		} else {
			act(libc.CString("You put $p in $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
		}
		if int(cont.Type_flag) == ITEM_PORTAL || int(cont.Type_flag) == ITEM_VEHICLE {
			obj_from_obj(obj)
			obj_to_room(obj, real_room(room_vnum(cont.Value[VAL_CONTAINER_CAPACITY])))
			if int(cont.Type_flag) == ITEM_PORTAL {
				act(libc.CString("What? $U$p disappears from $P in a puff of smoke!"), TRUE, ch, obj, unsafe.Pointer(cont), TO_ROOM)
				act(libc.CString("What? $U$p disappears from $P in a puff of smoke!"), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
			}
		}
	}
}
func do_put(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1         [2048]byte
		arg2         [2048]byte
		arg3         [2048]byte
		obj          *obj_data
		next_obj     *obj_data
		cont         *obj_data
		tmp_char     *char_data
		obj_dotmode  int
		cont_dotmode int
		found        int = 0
		howmany      int = 1
		theobj       *byte
		thecont      *byte
	)
	one_argument(two_arguments(argument, &arg1[0], &arg2[0]), &arg3[0])
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if arg3[0] != 0 && is_number(&arg1[0]) != 0 {
		howmany = libc.Atoi(libc.GoString(&arg1[0]))
		theobj = &arg2[0]
		thecont = &arg3[0]
	} else {
		theobj = &arg1[0]
		thecont = &arg2[0]
	}
	obj_dotmode = find_all_dots(theobj)
	cont_dotmode = find_all_dots(thecont)
	if *theobj == 0 {
		send_to_char(ch, libc.CString("Put what in what?\r\n"))
	} else if cont_dotmode != FIND_INDIV {
		send_to_char(ch, libc.CString("You can only put things into one container at a time.\r\n"))
	} else if *thecont == 0 {
		send_to_char(ch, libc.CString("What do you want to put %s in?\r\n"), func() string {
			if obj_dotmode == FIND_INDIV {
				return "it"
			}
			return "them"
		}())
	} else {
		generic_find(thecont, (1<<2)|1<<5|1<<3, ch, &tmp_char, &cont)
		if cont == nil {
			send_to_char(ch, libc.CString("You don't see %s %s here.\r\n"), AN(thecont), thecont)
		} else if int(cont.Type_flag) != ITEM_CONTAINER && int(cont.Type_flag) != ITEM_PORTAL && int(cont.Type_flag) != ITEM_VEHICLE {
			act(libc.CString("$p is not a container."), FALSE, ch, cont, nil, TO_CHAR)
		} else if OBJVAL_FLAGGED(cont, 1<<2) {
			send_to_char(ch, libc.CString("You'd better open it first!\r\n"))
		} else {
			if obj_dotmode == FIND_INDIV {
				if (func() *obj_data {
					obj = get_obj_in_list_vis(ch, theobj, nil, ch.Carrying)
					return obj
				}()) == nil {
					send_to_char(ch, libc.CString("You aren't carrying %s %s.\r\n"), AN(theobj), theobj)
				} else if obj == cont && howmany == 1 {
					send_to_char(ch, libc.CString("You attempt to fold it into itself, but fail.\r\n"))
				} else {
					for obj != nil && howmany != 0 {
						next_obj = obj.Next_content
						if obj != cont {
							howmany--
							perform_put(ch, obj, cont)
						}
						obj = get_obj_in_list_vis(ch, theobj, nil, next_obj)
					}
				}
			} else {
				for obj = ch.Carrying; obj != nil; obj = next_obj {
					next_obj = obj.Next_content
					if obj != cont && CAN_SEE_OBJ(ch, obj) && (obj_dotmode == FIND_ALL || isname(theobj, obj.Name) != 0) {
						found = 1
						perform_put(ch, obj, cont)
					}
				}
				if found == 0 {
					if obj_dotmode == FIND_ALL {
						send_to_char(ch, libc.CString("You don't seem to have anything to put in it.\r\n"))
					} else {
						send_to_char(ch, libc.CString("You don't seem to have any %ss.\r\n"), theobj)
					}
				}
			}
		}
	}
}
func can_take_obj(ch *char_data, obj *obj_data) int {
	if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_TAKE) {
		act(libc.CString("$p: you can't take that!"), FALSE, ch, obj, nil, TO_CHAR)
		return 0
	} else if int(ch.Carry_items) >= 50 {
		act(libc.CString("$p: your arms are full!"), FALSE, ch, obj, nil, TO_CHAR)
		return 0
	} else if (gear_weight(ch) + int(obj.Weight)) > int(max_carry_weight(ch)) {
		act(libc.CString("$p: you can't carry that much weight."), FALSE, ch, obj, nil, TO_CHAR)
		return 0
	} else if (obj.Weight+int64(world[ch.In_room].Gravity))+int64(gear_weight(ch)) > max_carry_weight(ch) {
		act(libc.CString("$p: you can't carry that much weight because of the gravity."), FALSE, ch, obj, nil, TO_CHAR)
		return 0
	}
	return 1
}
func get_check_money(ch *char_data, obj *obj_data) {
	var value int = (obj.Value[VAL_MONEY_SIZE])
	if int(obj.Type_flag) != ITEM_MONEY || value <= 0 {
		return
	}
	if ch.Gold+value > GOLD_CARRY(ch) {
		send_to_char(ch, libc.CString("You can only carry %s zenni at your current level, and leave the rest.\r\n"), add_commas(int64(GOLD_CARRY(ch))))
		act(libc.CString("@w$n @wdrops some onto the ground.@n"), FALSE, ch, nil, nil, TO_ROOM)
		extract_obj(obj)
		var diff int = 0
		diff = (ch.Gold + value) - GOLD_CARRY(ch)
		obj = create_money(diff)
		obj_to_room(obj, ch.In_room)
		ch.Gold = GOLD_CARRY(ch)
		return
	}
	ch.Gold += value
	extract_obj(obj)
	if value == 1 {
		send_to_char(ch, libc.CString("There was 1 zenni.\r\n"))
	} else {
		send_to_char(ch, libc.CString("There were %d zenni.\r\n"), value)
		if AFF_FLAGGED(ch, AFF_GROUP) && PRF_FLAGGED(ch, PRF_AUTOSPLIT) {
			var split [2048]byte
			stdio.Sprintf(&split[0], "%d", value)
			do_split(ch, &split[0], 0, 0)
		}
	}
}
func perform_get_from_container(ch *char_data, obj *obj_data, cont *obj_data, mode int) {
	if mode == (1<<2) || mode == (1<<5) || can_take_obj(ch, obj) != 0 {
		if int(ch.Carry_items) >= 50 {
			act(libc.CString("$p: you can't hold any more items."), FALSE, ch, obj, nil, TO_CHAR)
			return
		}
		if ch.Sits != nil && GET_OBJ_VNUM(ch.Sits) > 603 && GET_OBJ_VNUM(ch.Sits) < 608 && GET_OBJ_VNUM(ch.Sits)-4 != GET_OBJ_VNUM(cont) && GET_OBJ_VNUM(cont) > 599 && GET_OBJ_VNUM(cont) < 604 {
			send_to_char(ch, libc.CString("You aren't playing at that table!\r\n"))
			return
		} else if get_otrigger(obj, ch) != 0 {
			obj_from_obj(obj)
			obj_to_char(obj, ch)
			if OBJ_FLAGGED(cont, ITEM_SHEATH) {
				act(libc.CString("You draw $p from $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
				act(libc.CString("$n draws $p from $P."), TRUE, ch, obj, unsafe.Pointer(cont), TO_ROOM)
			} else {
				act(libc.CString("You get $p from $P."), FALSE, ch, obj, unsafe.Pointer(cont), TO_CHAR)
				act(libc.CString("$n gets $p from $P."), TRUE, ch, obj, unsafe.Pointer(cont), TO_ROOM)
			}
			if OBJ_FLAGGED(obj, ITEM_HOT) {
				if (ch.Bonuses[BONUS_FIREPROOF]) <= 0 && int(ch.Race) != RACE_DEMON {
					ch.Hit -= int64(float64(ch.Hit) * 0.25)
					if (ch.Bonuses[BONUS_FIREPRONE]) > 0 {
						ch.Hit = 1
					}
					SET_BIT_AR(ch.Affected_by[:], AFF_BURNED)
					act(libc.CString("@RYou are burned by it!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n@R is burned by it!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			if IS_NPC(ch) {
				item_check(obj, ch)
			}
			get_check_money(ch, obj)
		}
	}
}
func get_from_container(ch *char_data, cont *obj_data, arg *byte, mode int, howmany int) {
	var (
		obj         *obj_data
		next_obj    *obj_data
		obj_dotmode int
		found       int = 0
	)
	obj_dotmode = find_all_dots(arg)
	if OBJVAL_FLAGGED(cont, 1<<2) {
		act(libc.CString("$p is closed."), FALSE, ch, cont, nil, TO_CHAR)
	} else if obj_dotmode == FIND_INDIV {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, arg, nil, cont.Contains)
			return obj
		}()) == nil {
			var buf [64936]byte
			stdio.Snprintf(&buf[0], int(64936), "There doesn't seem to be %s %s in $p.", AN(arg), arg)
			act(&buf[0], FALSE, ch, cont, nil, TO_CHAR)
		} else {
			var obj_next *obj_data
			for obj != nil && func() int {
				p := &howmany
				x := *p
				*p--
				return x
			}() != 0 {
				obj_next = obj.Next_content
				perform_get_from_container(ch, obj, cont, mode)
				obj = get_obj_in_list_vis(ch, arg, nil, obj_next)
			}
		}
	} else {
		if obj_dotmode == FIND_ALLDOT && *arg == 0 {
			send_to_char(ch, libc.CString("Get all of what?\r\n"))
			return
		}
		for obj = cont.Contains; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if CAN_SEE_OBJ(ch, obj) && (obj_dotmode == FIND_ALL || isname(arg, obj.Name) != 0) {
				found = 1
				perform_get_from_container(ch, obj, cont, mode)
			}
		}
		if found == 0 {
			if obj_dotmode == FIND_ALL {
				act(libc.CString("$p seems to be empty."), FALSE, ch, cont, nil, TO_CHAR)
			} else {
				var buf [64936]byte
				stdio.Snprintf(&buf[0], int(64936), "You can't seem to find any %ss in $p.", arg)
				act(&buf[0], FALSE, ch, cont, nil, TO_CHAR)
			}
		}
	}
}
func perform_get_from_room(ch *char_data, obj *obj_data) int {
	if obj.Sitting != nil {
		send_to_char(ch, libc.CString("Someone is on that!\r\n"))
		return 0
	}
	if OBJ_FLAGGED(obj, ITEM_BURIED) {
		send_to_char(ch, libc.CString("Get what?\r\n"))
		return 0
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
		send_to_char(ch, libc.CString("You can't get things from a garden. Help garden.\r\n"))
		return 0
	}
	if can_take_obj(ch, obj) != 0 && get_otrigger(obj, ch) != 0 {
		obj_from_room(obj)
		obj_to_char(obj, ch)
		act(libc.CString("You get $p."), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("$n gets $p."), TRUE, ch, obj, nil, TO_ROOM)
		if OBJ_FLAGGED(obj, ITEM_HOT) {
			if (ch.Bonuses[BONUS_FIREPROOF]) <= 0 && int(ch.Race) != RACE_DEMON {
				ch.Hit -= int64(float64(ch.Hit) * 0.25)
				if (ch.Bonuses[BONUS_FIREPRONE]) > 0 {
					ch.Hit = 1
				}
				SET_BIT_AR(ch.Affected_by[:], AFF_BURNED)
				act(libc.CString("@RYou are burned by it!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n@R is burned by it!@n"), TRUE, ch, nil, nil, TO_ROOM)
			}
		}
		if IS_NPC(ch) {
			item_check(obj, ch)
		}
		get_check_money(ch, obj)
	}
	return 0
}
func find_exdesc_keywords(word *byte, list *extra_descr_data) *byte {
	var i *extra_descr_data
	for i = list; i != nil; i = i.Next {
		if isname(word, i.Keyword) != 0 {
			return i.Keyword
		}
	}
	return nil
}
func get_from_room(ch *char_data, arg *byte, howmany int) {
	var (
		obj      *obj_data
		next_obj *obj_data
		dotmode  int
		found    int = 0
		descword *byte
	)
	if find_exdesc(arg, world[ch.In_room].Ex_description) != nil {
		send_to_char(ch, libc.CString("You can't take %s %s.\r\n"), AN(arg), arg)
		return
	}
	dotmode = find_all_dots(arg)
	if dotmode == FIND_INDIV {
		if (func() *byte {
			descword = find_exdesc_keywords(arg, world[ch.In_room].Ex_description)
			return descword
		}()) != nil {
			send_to_char(ch, libc.CString("%s: you can't take that!\r\n"), fname(descword))
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, arg, nil, world[ch.In_room].Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't see %s %s here.\r\n"), AN(arg), arg)
		} else {
			var obj_next *obj_data
			for obj != nil && func() int {
				p := &howmany
				x := *p
				*p--
				return x
			}() != 0 {
				obj_next = obj.Next_content
				perform_get_from_room(ch, obj)
				obj = get_obj_in_list_vis(ch, arg, nil, obj_next)
			}
		}
	} else {
		if dotmode == FIND_ALLDOT && *arg == 0 {
			send_to_char(ch, libc.CString("Get all of what?\r\n"))
			return
		}
		for obj = world[ch.In_room].Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if CAN_SEE_OBJ(ch, obj) && (dotmode == FIND_ALL || isname(arg, obj.Name) != 0) {
				found = 1
				perform_get_from_room(ch, obj)
			}
		}
		if found == 0 {
			if dotmode == FIND_ALL {
				send_to_char(ch, libc.CString("There doesn't seem to be anything here.\r\n"))
			} else {
				send_to_char(ch, libc.CString("You don't see any %ss here.\r\n"), arg)
			}
		}
	}
}
func do_get(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1         [2048]byte
		arg2         [2048]byte
		arg3         [2048]byte
		cont_dotmode int
		found        int = 0
		mode         int
		cont         *obj_data
		tmp_char     *char_data
	)
	one_argument(two_arguments(argument, &arg1[0], &arg2[0]), &arg3[0])
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("Get what?\r\n"))
	} else if arg2[0] == 0 {
		get_from_room(ch, &arg1[0], 1)
	} else if is_number(&arg1[0]) != 0 && arg3[0] == 0 {
		get_from_room(ch, &arg2[0], libc.Atoi(libc.GoString(&arg1[0])))
	} else {
		var amount int = 1
		if is_number(&arg1[0]) != 0 {
			amount = libc.Atoi(libc.GoString(&arg1[0]))
			libc.StrCpy(&arg1[0], &arg2[0])
			libc.StrCpy(&arg2[0], &arg3[0])
		}
		cont_dotmode = find_all_dots(&arg2[0])
		if cont_dotmode == FIND_INDIV {
			mode = generic_find(&arg2[0], (1<<2)|1<<5|1<<3, ch, &tmp_char, &cont)
			if cont == nil {
				send_to_char(ch, libc.CString("You don't have %s %s.\r\n"), AN(&arg2[0]), &arg2[0])
			} else if int(cont.Type_flag) == ITEM_VEHICLE {
				send_to_char(ch, libc.CString("You will need to enter it first.\r\n"))
			} else if int(cont.Type_flag) != ITEM_CONTAINER && (int(cont.Type_flag) != ITEM_PORTAL || !OBJVAL_FLAGGED(cont, 1<<0)) {
				act(libc.CString("$p is not a container."), FALSE, ch, cont, nil, TO_CHAR)
			} else {
				get_from_container(ch, cont, &arg1[0], mode, amount)
			}
		} else {
			if cont_dotmode == FIND_ALLDOT && arg2[0] == 0 {
				send_to_char(ch, libc.CString("Get from all of what?\r\n"))
				return
			}
			for cont = ch.Carrying; cont != nil; cont = cont.Next_content {
				if CAN_SEE_OBJ(ch, cont) && (cont_dotmode == FIND_ALL || isname(&arg2[0], cont.Name) != 0) {
					if int(cont.Type_flag) == ITEM_CONTAINER {
						found = 1
						get_from_container(ch, cont, &arg1[0], 1<<2, amount)
					} else if cont_dotmode == FIND_ALLDOT {
						found = 1
						act(libc.CString("$p is not a container."), FALSE, ch, cont, nil, TO_CHAR)
					}
				}
			}
			for cont = world[ch.In_room].Contents; cont != nil; cont = cont.Next_content {
				if CAN_SEE_OBJ(ch, cont) && (cont_dotmode == FIND_ALL || isname(&arg2[0], cont.Name) != 0) {
					if int(cont.Type_flag) == ITEM_CONTAINER {
						get_from_container(ch, cont, &arg1[0], 1<<3, amount)
						found = 1
					} else if cont_dotmode == FIND_ALLDOT {
						act(libc.CString("$p is not a container."), FALSE, ch, cont, nil, TO_CHAR)
						found = 1
					}
				}
			}
			if found == 0 {
				if cont_dotmode == FIND_ALL {
					send_to_char(ch, libc.CString("You can't seem to find any containers.\r\n"))
				} else {
					send_to_char(ch, libc.CString("You can't seem to find any %ss here.\r\n"), &arg2[0])
				}
			}
		}
	}
}
func perform_drop_gold(ch *char_data, amount int, mode int8, RDR room_rnum) {
	var obj *obj_data
	if amount <= 0 {
		send_to_char(ch, libc.CString("Heh heh heh.. we are jolly funny today, eh?\r\n"))
	} else if ch.Gold < amount {
		send_to_char(ch, libc.CString("You don't have that many zenni!\r\n"))
	} else {
		if int(mode) != SCMD_JUNK {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			obj = create_money(amount)
			if int(mode) == SCMD_DONATE {
				send_to_char(ch, libc.CString("You throw some zenni into the air where it disappears in a puff of smoke!\r\n"))
				act(libc.CString("$n throws some gold into the air where it disappears in a puff of smoke!"), FALSE, ch, nil, nil, TO_ROOM)
				obj_to_room(obj, RDR)
				act(libc.CString("$p suddenly appears in a puff of orange smoke!"), 0, nil, obj, nil, TO_ROOM)
			} else {
				var buf [64936]byte
				if drop_wtrigger(obj, ch) == 0 {
					extract_obj(obj)
					return
				}
				if drop_wtrigger(obj, ch) == 0 && obj != nil {
					extract_obj(obj)
					return
				}
				stdio.Snprintf(&buf[0], int(64936), "$n drops %s.", money_desc(amount))
				act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
				send_to_char(ch, libc.CString("You drop some zenni.\r\n"))
				obj_to_room(obj, ch.In_room)
				if ch.Admlevel > 0 {
					send_to_imm(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
					log_imm_action(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
					if check_insidebag(obj, 0.0) > 1 {
						send_to_imm(libc.CString("IMM DROP: Object contains %d other items."), check_insidebag(obj, 0.0))
						log_imm_action(libc.CString("IMM DROP: Object contains %d other items."), check_insidebag(obj, 0.0))
					}
				}
			}
		} else {
			var buf [64936]byte
			stdio.Snprintf(&buf[0], int(64936), "$n drops %s which disappears in a puff of smoke!", money_desc(amount))
			act(&buf[0], FALSE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("You drop some zenni which disappears in a puff of smoke!\r\n"))
		}
		ch.Gold -= amount
	}
}
func perform_drop(ch *char_data, obj *obj_data, mode int8, sname *byte, RDR room_rnum) int {
	var (
		buf   [64936]byte
		value int
	)
	if drop_otrigger(obj, ch) == 0 {
		return 0
	}
	if int(mode) == SCMD_DROP && drop_wtrigger(obj, ch) == 0 {
		return 0
	}
	if GET_OBJ_VNUM(obj) == 17 || GET_OBJ_VNUM(obj) == 0x464E {
		stdio.Snprintf(&buf[0], int(64936), "You can't %s $p, it is grafted into your soul :P", sname)
		act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
		return 0
	}
	if GET_OBJ_VNUM(obj) == 20 || GET_OBJ_VNUM(obj) == 21 || GET_OBJ_VNUM(obj) == 22 || GET_OBJ_VNUM(obj) == 23 || GET_OBJ_VNUM(obj) == 24 || GET_OBJ_VNUM(obj) == 25 || GET_OBJ_VNUM(obj) == 26 {
		if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
			stdio.Snprintf(&buf[0], int(64936), "You can't %s $p in space!", sname)
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			return 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
			stdio.Snprintf(&buf[0], int(64936), "You can't %s $p in here. Read help garden.", sname)
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			return 0
		}
		if int(mode) == SCMD_DROP && OBJ_FLAGGED(obj, ITEM_NORENT) {
			stdio.Snprintf(&buf[0], int(64936), "You drop $p but it gets lost on the ground!")
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			obj_from_char(obj)
			extract_obj(obj)
			return 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_NOINSTANT) {
			stdio.Snprintf(&buf[0], int(64936), "You can't %s $p in this protected area!", sname)
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			return 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_SHIP) {
			stdio.Snprintf(&buf[0], int(64936), "You can't %s $p on a private ship!", sname)
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			return 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HOUSE) {
			stdio.Snprintf(&buf[0], int(64936), "You can't %s $p in a private house!", sname)
			act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
			return 0
		}
	}
	if OBJ_FLAGGED(obj, ITEM_NODROP) && ch.Admlevel < 1 {
		stdio.Snprintf(&buf[0], int(64936), "You can't %s $p, it must be CURSED!", sname)
		act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
		return 0
	}
	if (int(mode) == SCMD_DONATE || int(mode) == SCMD_JUNK) && OBJ_FLAGGED(obj, ITEM_NODONATE) {
		stdio.Snprintf(&buf[0], int(64936), "You can't %s $p!", sname)
		act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
		return 0
	}
	stdio.Snprintf(&buf[0], int(64936), "You %s $p.", sname)
	act(&buf[0], FALSE, ch, obj, nil, TO_CHAR)
	stdio.Snprintf(&buf[0], int(64936), "$n %ss $p.", sname)
	act(&buf[0], TRUE, ch, obj, nil, TO_ROOM)
	obj_from_char(obj)
	switch mode {
	case SCMD_DROP:
		if !OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) && world[ch.In_room].Geffect == 6 {
			act(libc.CString("$p melts in the lava!"), FALSE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$p melts in the lava!"), FALSE, ch, obj, nil, TO_ROOM)
			extract_obj(obj)
		} else if world[ch.In_room].Geffect == 6 {
			act(libc.CString("$p plops down on some cooled lava!"), FALSE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$p plops down on some cooled lava!"), FALSE, ch, obj, nil, TO_ROOM)
			obj_to_room(obj, ch.In_room)
			if ch.Admlevel > 0 {
				send_to_imm(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
				log_imm_action(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
			}
		} else {
			obj_to_room(obj, ch.In_room)
			if ch.Admlevel > 0 {
				send_to_imm(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
				log_imm_action(libc.CString("IMM DROP: %s dropped %s in room [%d]"), GET_NAME(ch), obj.Short_description, GET_ROOM_VNUM(obj.In_room))
			}
		}
		return 0
	case SCMD_DONATE:
		obj_to_room(obj, RDR)
		act(libc.CString("$p suddenly appears in a puff a smoke!"), FALSE, nil, obj, nil, TO_ROOM)
		return 0
	case SCMD_JUNK:
		value = int(MAX(1, MIN(200, int64(obj.Cost/16))))
		extract_obj(obj)
		return value
	default:
		basic_mud_log(libc.CString("SYSERR: Incorrect argument %d passed to perform_drop."), mode)
	}
	return 0
}
func do_drop(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		obj      *obj_data
		next_obj *obj_data
		RDR      room_rnum = 0
		mode     int8      = SCMD_DROP
		dotmode  int
		amount   int = 0
	)
	_ = amount
	var multi int
	var num_don_rooms int
	var sname *byte
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(ch.In_room, ROOM_GARDEN2) {
		send_to_char(ch, libc.CString("You can not do that in a garden.\r\n"))
		return
	}
	switch subcmd {
	case SCMD_JUNK:
		sname = libc.CString("junk")
		mode = SCMD_JUNK
	case SCMD_DONATE:
		sname = libc.CString("donate")
		mode = SCMD_DONATE
		num_don_rooms = int(libc.BoolToInt(config_info.Room_nums.Donation_room_1 != room_vnum(-1)))*2 + int(libc.BoolToInt(config_info.Room_nums.Donation_room_2 != room_vnum(-1))) + int(libc.BoolToInt(config_info.Room_nums.Donation_room_3 != room_vnum(-1))) + 1
		switch rand_number(0, num_don_rooms) {
		case 0:
			mode = SCMD_JUNK
		case 1:
			fallthrough
		case 2:
			RDR = real_room(config_info.Room_nums.Donation_room_1)
		case 3:
			RDR = real_room(config_info.Room_nums.Donation_room_2)
		case 4:
			RDR = real_room(config_info.Room_nums.Donation_room_3)
		}
		if RDR == room_rnum(-1) {
			send_to_char(ch, libc.CString("Sorry, you can't donate anything right now.\r\n"))
			return
		}
	default:
		sname = libc.CString("drop")
	}
	argument = one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What do you want to %s?\r\n"), sname)
		return
	} else if is_number(&arg[0]) != 0 {
		multi = libc.Atoi(libc.GoString(&arg[0]))
		one_argument(argument, &arg[0])
		if libc.StrCaseCmp(libc.CString("zenni"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("gold"), &arg[0]) == 0 {
			perform_drop_gold(ch, multi, mode, RDR)
		} else if multi <= 0 {
			send_to_char(ch, libc.CString("Yeah, that makes sense.\r\n"))
		} else if arg[0] == 0 {
			send_to_char(ch, libc.CString("What do you want to %s %d of?\r\n"), sname, multi)
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't seem to have any %ss.\r\n"), &arg[0])
		} else if check_saveroom_count(ch, obj) > 150 {
			send_to_char(ch, libc.CString("The save room you are in can not hold anymore items! (150 max, count of items in containers is halved)\r\n"))
		} else {
			for {
				next_obj = get_obj_in_list_vis(ch, &arg[0], nil, obj.Next_content)
				amount += perform_drop(ch, obj, mode, sname, RDR)
				obj = next_obj
				if obj == nil || func() int {
					p := &multi
					*p--
					return *p
				}() == 0 {
					break
				}
			}
		}
	} else {
		dotmode = find_all_dots(&arg[0])
		if dotmode == FIND_ALL && (subcmd == SCMD_JUNK || subcmd == SCMD_DONATE) {
			if subcmd == SCMD_JUNK {
				send_to_char(ch, libc.CString("Go to the dump if you want to junk EVERYTHING!\r\n"))
			} else {
				send_to_char(ch, libc.CString("Go do the donation room if you want to donate EVERYTHING!\r\n"))
			}
			return
		}
		if dotmode == FIND_ALL {
			var fail int = FALSE
			if ch.Carrying == nil {
				send_to_char(ch, libc.CString("You don't seem to be carrying anything.\r\n"))
			} else {
				for obj = ch.Carrying; obj != nil; obj = next_obj {
					next_obj = obj.Next_content
					if check_saveroom_count(ch, obj) > 150 {
						fail = TRUE
					} else {
						amount += perform_drop(ch, obj, mode, sname, RDR)
					}
				}
				if fail == TRUE {
					send_to_char(ch, libc.CString("Some of the items couldn't be dropped into this save room. It is too full. (150 max, containers half the count inside)\r\n"))
				}
			}
		} else if dotmode == FIND_ALLDOT {
			var fail int = FALSE
			if arg[0] == 0 {
				send_to_char(ch, libc.CString("What do you want to %s all of?\r\n"), sname)
				return
			}
			if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
				return obj
			}()) == nil {
				send_to_char(ch, libc.CString("You don't seem to have any %ss.\r\n"), &arg[0])
			}
			for obj != nil {
				next_obj = get_obj_in_list_vis(ch, &arg[0], nil, obj.Next_content)
				if check_saveroom_count(ch, obj) > 150 {
					fail = TRUE
				} else {
					amount += perform_drop(ch, obj, mode, sname, RDR)
				}
				obj = next_obj
			}
			if fail == TRUE {
				send_to_char(ch, libc.CString("Some of the items couldn't be dropped into this save room. It is too full. (150 max, containers half the count inside)\r\n"))
			}
		} else {
			if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
				return obj
			}()) == nil {
				send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
			} else if check_saveroom_count(ch, obj) > 150 {
				send_to_char(ch, libc.CString("The item couldn't be dropped into this save room. It is too full. (150 max, containers half the count inside)\r\n"))
			} else {
				amount += perform_drop(ch, obj, mode, sname, RDR)
			}
		}
	}
}
func perform_give(ch *char_data, vict *char_data, obj *obj_data) {
	if give_otrigger(obj, ch, vict) == 0 {
		return
	}
	if receive_mtrigger(vict, ch, obj) == 0 {
		return
	}
	if OBJ_FLAGGED(obj, ITEM_NODROP) {
		act(libc.CString("You can't let go of $p!!  Yeech!"), FALSE, ch, obj, nil, TO_CHAR)
		return
	}
	if int(vict.Carry_items) >= 50 {
		act(libc.CString("$N seems to have $S hands full."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		if IS_NPC(ch) {
			act(libc.CString("$n@n drops $p because you can't carry anymore."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n@n drops $p on the ground since $N's unable to carry it."), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			obj_from_char(obj)
			obj_to_room(obj, ch.In_room)
		}
		return
	}
	if IS_NPC(vict) && (OBJ_FLAGGED(obj, ITEM_FORGED) || OBJ_FLAGGED(obj, ITEM_FORGED)) {
		act(libc.CString("$n tries to hand $p to $N."), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		do_say(vict, libc.CString("I don't want that piece of junk."), 0, 0)
		return
	}
	if obj.Weight+int64(gear_weight(vict)) > max_carry_weight(vict) {
		act(libc.CString("$E can't carry that much weight."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		if IS_NPC(ch) {
			act(libc.CString("$n@n drops $p because you can't carry anymore."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n@n drops $p on the ground since $N's unable to carry it."), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			obj_from_char(obj)
			obj_to_room(obj, ch.In_room)
		}
		return
	}
	if (obj.Weight+int64(world[vict.In_room].Gravity))+int64(gear_weight(vict)) > max_carry_weight(vict) {
		act(libc.CString("$E can't carry that much weight because of the gravity."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		if IS_NPC(ch) {
			act(libc.CString("$n@n drops $p because you can't carry anymore."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n@n drops $p on the ground since $N's unable to carry it."), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			obj_from_char(obj)
			obj_to_room(obj, ch.In_room)
		}
		return
	}
	if !IS_NPC(vict) && !IS_NPC(ch) {
		if PRF_FLAGGED(vict, PRF_NOGIVE) {
			act(libc.CString("$N refuses to accept $p at this time."), FALSE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("$n tries to give you $p, but you are refusing to be handed things."), FALSE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n tries to give $N, $p, but $E is refusing to be handed things right now."), FALSE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			return
		}
	}
	obj_from_char(obj)
	obj_to_char(obj, vict)
	act(libc.CString("You give $p to $N."), FALSE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
	act(libc.CString("$n gives you $p."), FALSE, ch, obj, unsafe.Pointer(vict), TO_VICT)
	act(libc.CString("$n gives $p to $N."), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
	if OBJ_FLAGGED(obj, ITEM_HOT) {
		if (vict.Bonuses[BONUS_FIREPROOF]) <= 0 && int(vict.Race) != RACE_DEMON {
			vict.Hit -= int64(float64(vict.Hit) * 0.25)
			if (vict.Bonuses[BONUS_FIREPRONE]) > 0 {
				vict.Hit = 1
			}
			SET_BIT_AR(vict.Affected_by[:], AFF_BURNED)
			act(libc.CString("@RYou are burned by it!@n"), TRUE, vict, nil, nil, TO_CHAR)
			act(libc.CString("@R$n@R is burned by it!@n"), TRUE, vict, nil, nil, TO_ROOM)
		}
	}
}
func give_find_vict(ch *char_data, arg *byte) *char_data {
	var vict *char_data
	skip_spaces(&arg)
	if *arg == 0 {
		send_to_char(ch, libc.CString("To who?\r\n"))
	} else if (func() *char_data {
		vict = get_char_vis(ch, arg, nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		if IS_NPC(ch) {
			send_to_imm(libc.CString("Mob Give: Victim, %s, doesn't exist."), arg)
		}
	} else if vict == ch {
		send_to_char(ch, libc.CString("What's the point of that?\r\n"))
	} else {
		return vict
	}
	return nil
}
func perform_give_gold(ch *char_data, vict *char_data, amount int) {
	var buf [64936]byte
	if amount <= 0 {
		send_to_char(ch, libc.CString("Heh heh heh ... we are jolly funny today, eh?\r\n"))
		return
	}
	if ch.Gold < amount && (IS_NPC(ch) || !ADM_FLAGGED(ch, ADM_MONEY)) {
		send_to_char(ch, libc.CString("You don't have that much zenni!\r\n"))
		return
	}
	if vict.Gold+amount > GOLD_CARRY(vict) {
		send_to_char(ch, libc.CString("They can't carry that much zenni.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	stdio.Snprintf(&buf[0], int(64936), "$n gives you %d zenni.", amount)
	act(&buf[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
	stdio.Snprintf(&buf[0], int(64936), "$n gives %s to $N.", money_desc(amount))
	act(&buf[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	if IS_NPC(ch) || !ADM_FLAGGED(ch, ADM_MONEY) {
		ch.Gold -= amount
	}
	vict.Gold += amount
	bribe_mtrigger(vict, ch, amount)
}
func do_give(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [64936]byte
		amount   int
		dotmode  int
		vict     *char_data
		obj      *obj_data
		next_obj *obj_data
	)
	argument = one_argument(argument, &arg[0])
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Give what to who?\r\n"))
	} else if is_number(&arg[0]) != 0 {
		amount = libc.Atoi(libc.GoString(&arg[0]))
		argument = one_argument(argument, &arg[0])
		if libc.StrCaseCmp(libc.CString("zenni"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("gold"), &arg[0]) == 0 {
			one_argument(argument, &arg[0])
			if (func() *char_data {
				vict = give_find_vict(ch, &arg[0])
				return vict
			}()) != nil {
				perform_give_gold(ch, vict, amount)
				if ch.Admlevel > 0 && !IS_NPC(vict) {
					send_to_imm(libc.CString("IMM GIVE: %s has given %s zenni to %s."), GET_NAME(ch), add_commas(int64(amount)), GET_NAME(vict))
					log_imm_action(libc.CString("IMM GIVE: %s has given %s zenni to %s."), GET_NAME(ch), add_commas(int64(amount)), GET_NAME(vict))
				}
			}
			return
		} else if arg[0] == 0 {
			send_to_char(ch, libc.CString("What do you want to give %d of?\r\n"), amount)
		} else if (func() *char_data {
			vict = give_find_vict(ch, argument)
			return vict
		}()) == nil {
			return
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't seem to have any %ss.\r\n"), &arg[0])
		} else {
			for obj != nil && func() int {
				p := &amount
				x := *p
				*p--
				return x
			}() != 0 {
				next_obj = get_obj_in_list_vis(ch, &arg[0], nil, obj.Next_content)
				perform_give(ch, vict, obj)
				if ch.Admlevel > 0 && !IS_NPC(vict) {
					send_to_imm(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
					log_imm_action(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
				}
				obj = next_obj
			}
		}
	} else {
		var buf1 [2048]byte
		one_argument(argument, &buf1[0])
		if (func() *char_data {
			vict = give_find_vict(ch, &buf1[0])
			return vict
		}()) == nil {
			return
		}
		dotmode = find_all_dots(&arg[0])
		if dotmode == FIND_INDIV {
			if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
				return obj
			}()) == nil {
				send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
			} else {
				perform_give(ch, vict, obj)
				if ch.Admlevel > 0 && !IS_NPC(vict) {
					send_to_imm(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
					log_imm_action(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
				}
			}
		} else {
			if dotmode == FIND_ALLDOT && arg[0] == 0 {
				send_to_char(ch, libc.CString("All of what?\r\n"))
				return
			}
			if ch.Carrying == nil {
				send_to_char(ch, libc.CString("You don't seem to be holding anything.\r\n"))
			} else {
				for obj = ch.Carrying; obj != nil; obj = next_obj {
					next_obj = obj.Next_content
					if CAN_SEE_OBJ(ch, obj) && (dotmode == FIND_ALL || isname(&arg[0], obj.Name) != 0) {
						perform_give(ch, vict, obj)
						if ch.Admlevel > 0 && !IS_NPC(vict) {
							send_to_imm(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
							log_imm_action(libc.CString("IMM GIVE: %s has given %s to %s."), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
						}
					}
				}
			}
		}
	}
}
func weight_change_object(obj *obj_data, weight int) {
	var (
		tmp_obj *obj_data
		tmp_ch  *char_data
	)
	if obj.In_room != room_rnum(-1) {
		obj.Weight += int64(weight)
	} else if (func() *char_data {
		tmp_ch = obj.Carried_by
		return tmp_ch
	}()) != nil {
		obj_from_char(obj)
		obj.Weight += int64(weight)
		obj_to_char(obj, tmp_ch)
	} else if (func() *obj_data {
		tmp_obj = obj.In_obj
		return tmp_obj
	}()) != nil {
		obj_from_obj(obj)
		obj.Weight += int64(weight)
		obj_to_obj(obj, tmp_obj)
	} else {
		basic_mud_log(libc.CString("SYSERR: Unknown attempt to subtract weight from an object."))
	}
}
func name_from_drinkcon(obj *obj_data) {
	var (
		new_name *byte
		cur_name *byte
		next     *byte
		liqname  *byte
		liqlen   int
		cpylen   int
	)
	if obj == nil || int(obj.Type_flag) != ITEM_DRINKCON && int(obj.Type_flag) != ITEM_FOUNTAIN {
		return
	}
	liqname = drinknames[obj.Value[VAL_DRINKCON_LIQUID]]
	if isname(liqname, obj.Name) == 0 {
		return
	}
	liqlen = libc.StrLen(liqname)
	new_name = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(obj.Name)-libc.StrLen(liqname))[0]))
	for cur_name = obj.Name; cur_name != nil; cur_name = next {
		if *cur_name == ' ' {
			cur_name = (*byte)(unsafe.Add(unsafe.Pointer(cur_name), 1))
		}
		if (func() *byte {
			next = libc.StrChr(cur_name, ' ')
			return next
		}()) != nil {
			cpylen = int(int64(uintptr(unsafe.Pointer(next)) - uintptr(unsafe.Pointer(cur_name))))
		} else {
			cpylen = libc.StrLen(cur_name)
		}
		if libc.StrNCaseCmp(cur_name, liqname, liqlen) == 0 {
			continue
		}
		if *new_name != 0 {
			libc.StrCat(new_name, libc.CString(" "))
		}
		libc.StrNCat(new_name, cur_name, cpylen)
	}
	if obj.Item_number == obj_vnum(-1) || obj.Name != obj_proto[obj.Item_number].Name {
		libc.Free(unsafe.Pointer(obj.Name))
	}
	obj.Name = new_name
}
func name_to_drinkcon(obj *obj_data, type_ int) {
	var new_name *byte
	if obj == nil || int(obj.Type_flag) != ITEM_DRINKCON && int(obj.Type_flag) != ITEM_FOUNTAIN {
		return
	}
	new_name = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(obj.Name)+libc.StrLen(drinknames[type_])+2)[0]))
	stdio.Sprintf(new_name, "%s %s", obj.Name, drinknames[type_])
	if obj.Item_number == obj_vnum(-1) || obj.Name != obj_proto[obj.Item_number].Name {
		libc.Free(unsafe.Pointer(obj.Name))
	}
	obj.Name = new_name
}
func do_drink(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg        [2048]byte
		temp       *obj_data
		af         affected_type
		amount     int
		weight     int
		wasthirsty int = 0
		on_ground  int = 0
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if int(ch.Race) == RACE_ANDROID || int(ch.Player_specials.Conditions[THIRST]) < 0 {
		send_to_char(ch, libc.CString("You need not drink!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if int(ch.Player_specials.Conditions[HUNGER]) <= 1 && int(ch.Player_specials.Conditions[THIRST]) >= 2 && int(ch.Race) != RACE_NAMEK && (ch.Genome[0]) != 3 && (ch.Genome[1]) != 3 {
		send_to_char(ch, libc.CString("You need to eat first!\r\n"))
		return
	}
	wasthirsty = int(ch.Player_specials.Conditions[THIRST])
	if arg[0] == 0 && !IS_NPC(ch) {
		var buf [64936]byte
		switch SECT(ch.In_room) {
		case SECT_WATER_SWIM:
			fallthrough
		case SECT_WATER_NOSWIM:
			fallthrough
		case SECT_UNDERWATER:
			stdio.Snprintf(&buf[0], int(64936), "$n takes a refreshing drink from the surrounding water.")
			act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("You take a refreshing drink from the surrounding water.\r\n"))
			gain_condition(ch, THIRST, 1)
			if GET_SKILL(ch, SKILL_WELLSPRING) != 0 && ch.Mana < ch.Max_mana && wasthirsty <= 30 && subcmd != SCMD_SIP {
				ch.Mana += int64(((float64(ch.Max_mana) * 0.005) + float64(int(ch.Aff_abils.Wis)*rand_number(80, 100))) * float64(GET_SKILL(ch, SKILL_WELLSPRING)))
				if ch.Mana > ch.Max_mana {
					ch.Mana = ch.Max_mana
				}
				send_to_char(ch, libc.CString("You feel your ki return to full strength.\r\n"))
			}
			if int(ch.Player_specials.Conditions[THIRST]) >= 48 {
				send_to_char(ch, libc.CString("You don't feel thirsty anymore.\r\n"))
			}
			return
		default:
			if !SUNKEN(ch.In_room) {
				send_to_char(ch, libc.CString("Drink from what?\r\n"))
				return
			} else {
				stdio.Snprintf(&buf[0], int(64936), "$n takes a refreshing drink from the surrounding water.")
				act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
				send_to_char(ch, libc.CString("You take a refreshing drink from the surrounding water.\r\n"))
				gain_condition(ch, THIRST, 1)
				if GET_SKILL(ch, SKILL_WELLSPRING) != 0 && ch.Mana < ch.Max_mana && wasthirsty <= 30 && subcmd != SCMD_SIP {
					ch.Mana += int64(((float64(ch.Max_mana) * 0.005) + float64(int(ch.Aff_abils.Wis)*rand_number(80, 100))) * float64(GET_SKILL(ch, SKILL_WELLSPRING)))
					if ch.Mana > ch.Max_mana {
						ch.Mana = ch.Max_mana
						send_to_char(ch, libc.CString("You feel your ki return to full strength.\r\n"))
					} else {
						send_to_char(ch, libc.CString("You feel your ki has rejuvenated.\r\n"))
					}
				}
				if int(ch.Player_specials.Conditions[THIRST]) >= 48 {
					send_to_char(ch, libc.CString("You don't feel thirsty anymore.\r\n"))
				}
			}
			return
		}
	}
	if (func() *obj_data {
		temp = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return temp
	}()) == nil {
		if (func() *obj_data {
			temp = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
			return temp
		}()) == nil {
			send_to_char(ch, libc.CString("You can't find it!\r\n"))
			return
		} else {
			on_ground = 1
		}
	}
	if int(temp.Type_flag) != ITEM_DRINKCON && int(temp.Type_flag) != ITEM_FOUNTAIN {
		if GET_OBJ_VNUM(temp) == 86 && on_ground != 1 {
			act(libc.CString("@wYou uncork the $p and tip it to your lips. Drinking it down you feel a warmth flow through your body and your ki returns to full strength .@n"), TRUE, ch, temp, nil, TO_CHAR)
			act(libc.CString("@C$n@w uncorks the $p and tips it to $s lips. Drinking it down and then smiling.@n"), TRUE, ch, temp, nil, TO_ROOM)
			obj_from_char(temp)
			extract_obj(temp)
			ch.Mana = ch.Max_mana
			ch.Player_specials.Conditions[THIRST] += 8
			return
		} else if GET_OBJ_VNUM(temp) == 86 && on_ground == 1 {
			send_to_char(ch, libc.CString("You need to pick that up first.\r\n"))
			return
		} else {
			send_to_char(ch, libc.CString("You can't drink from that!\r\n"))
			return
		}
	}
	if on_ground != 0 && int(temp.Type_flag) == ITEM_DRINKCON {
		send_to_char(ch, libc.CString("You have to be holding that to drink from it.\r\n"))
		return
	}
	if OBJ_FLAGGED(temp, ITEM_BROKEN) {
		send_to_char(ch, libc.CString("Seems like it's broken!\r\n"))
		return
	}
	if OBJ_FLAGGED(temp, ITEM_FORGED) {
		send_to_char(ch, libc.CString("Seems like it doesn't work, maybe it is fake...\r\n"))
		return
	}
	if IS_NPC(ch) {
		act(libc.CString("$n@w drinks from $p."), TRUE, ch, temp, nil, TO_ROOM)
		obj_from_char(temp)
		extract_obj(temp)
		return
	}
	if int(ch.Player_specials.Conditions[DRUNK]) > 10 && int(ch.Player_specials.Conditions[THIRST]) > 0 {
		send_to_char(ch, libc.CString("You can't seem to get close enough to your mouth.\r\n"))
		act(libc.CString("$n tries to drink but misses $s mouth!"), TRUE, ch, nil, nil, TO_ROOM)
		return
	}
	if (temp.Value[VAL_DRINKCON_HOWFULL]) <= 0 && (temp.Value[VAL_DRINKCON_CAPACITY]) >= 1 {
		send_to_char(ch, libc.CString("It's empty.\r\n"))
		return
	}
	if consume_otrigger(temp, ch, OCMD_DRINK) == 0 {
		return
	}
	if subcmd == SCMD_DRINK {
		var buf [64936]byte
		stdio.Snprintf(&buf[0], int(64936), "$n drinks %s from $p.", drinks[temp.Value[VAL_DRINKCON_LIQUID]])
		act(&buf[0], TRUE, ch, temp, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You drink the %s.\r\n"), drinks[temp.Value[VAL_DRINKCON_LIQUID]])
		if temp.Action_description != nil {
			act(temp.Action_description, TRUE, ch, temp, nil, TO_CHAR)
		}
		amount = 4
	} else {
		act(libc.CString("$n sips from $p."), TRUE, ch, temp, nil, TO_ROOM)
		send_to_char(ch, libc.CString("It tastes like %s.\r\n"), drinks[temp.Value[VAL_DRINKCON_LIQUID]])
		if (temp.Value[VAL_DRINKCON_POISON]) != 0 {
			send_to_char(ch, libc.CString("It has a sickening taste! Better not eat it..."))
			return
		}
		amount = 1
	}
	amount = int(MIN(int64(amount), int64(temp.Value[VAL_DRINKCON_HOWFULL])))
	if (temp.Value[VAL_DRINKCON_CAPACITY]) > 0 {
		weight = int(MIN(int64(amount), temp.Weight))
		weight_change_object(temp, -weight)
	}
	gain_condition(ch, DRUNK, drink_aff[temp.Value[VAL_DRINKCON_LIQUID]][DRUNK]*amount)
	gain_condition(ch, HUNGER, drink_aff[temp.Value[VAL_DRINKCON_LIQUID]][HUNGER]*amount)
	gain_condition(ch, THIRST, drink_aff[temp.Value[VAL_DRINKCON_LIQUID]][THIRST]*amount)
	if ch.Foodr == 0 && subcmd != SCMD_SIP {
		ch.Move += (ch.Max_move / 100) * int64(amount)
		ch.Foodr = 2
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		send_to_char(ch, libc.CString("You feel rejuvinated by it.\r\n"))
	}
	if GET_SKILL(ch, SKILL_WELLSPRING) != 0 && ch.Mana < ch.Max_mana && wasthirsty <= 30 && subcmd != SCMD_SIP {
		if (temp.Value[VAL_DRINKCON_LIQUID]) == 0 || (temp.Value[VAL_DRINKCON_LIQUID]) == 14 || (temp.Value[VAL_DRINKCON_LIQUID]) == 15 {
			ch.Mana += int64(((float64(ch.Max_mana) * 0.005) + float64(int(ch.Aff_abils.Wis)*rand_number(80, 100))) * float64(GET_SKILL(ch, SKILL_WELLSPRING)))
			if ch.Mana > ch.Max_mana {
				ch.Mana = ch.Max_mana
				send_to_char(ch, libc.CString("You feel your ki return to full strength.\r\n"))
			} else {
				send_to_char(ch, libc.CString("You feel your ki has rejuvenated.\r\n"))
			}
		}
	}
	if int(ch.Player_specials.Conditions[DRUNK]) > 10 {
		send_to_char(ch, libc.CString("You feel drunk.\r\n"))
	}
	if int(ch.Player_specials.Conditions[THIRST]) >= 48 {
		send_to_char(ch, libc.CString("You don't feel thirsty anymore.\r\n"))
	}
	if (temp.Value[VAL_DRINKCON_POISON]) != 0 && (int(ch.Race) != RACE_MUTANT || (ch.Genome[0]) != 7 && (ch.Genome[1]) != 7) {
		send_to_char(ch, libc.CString("Oops, it tasted rather strange!\r\n"))
		act(libc.CString("$n chokes and utters some strange sounds."), TRUE, ch, nil, nil, TO_ROOM)
		af.Type = SPELL_POISON
		af.Duration = int16(amount * 3)
		af.Modifier = 0
		af.Location = APPLY_NONE
		af.Bitvector = AFF_POISON
		affect_join(ch, &af, FALSE != 0, FALSE != 0, FALSE != 0, FALSE != 0)
	}
	if (temp.Value[VAL_DRINKCON_CAPACITY]) > 0 {
		temp.Value[VAL_DRINKCON_HOWFULL] -= amount
		if (temp.Value[VAL_DRINKCON_HOWFULL]) == 0 {
			name_from_drinkcon(temp)
			temp.Value[VAL_DRINKCON_POISON] = 0
		}
	}
	return
}
func do_eat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		food   *obj_data
		af     affected_type
		amount int
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if int(ch.Race) == RACE_ANDROID || int(ch.Player_specials.Conditions[HUNGER]) < 0 {
		send_to_char(ch, libc.CString("You need not eat!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		send_to_char(ch, libc.CString("You feel too sick from the poison to eat!\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Eat what?\r\n"))
		return
	}
	if (func() *obj_data {
		food = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return food
	}()) == nil {
		send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
		return
	}
	if int(ch.Player_specials.Conditions[THIRST]) <= 1 && int(ch.Player_specials.Conditions[HUNGER]) >= 2 {
		send_to_char(ch, libc.CString("You need to drink first!\r\n"))
		return
	}
	if subcmd == SCMD_TASTE && (int(food.Type_flag) == ITEM_DRINKCON || int(food.Type_flag) == ITEM_FOUNTAIN) {
		do_drink(ch, argument, 0, SCMD_SIP)
		return
	}
	if int(food.Type_flag) != ITEM_FOOD {
		send_to_char(ch, libc.CString("You can't eat THAT!\r\n"))
		return
	}
	if consume_otrigger(food, ch, OCMD_EAT) == 0 {
		return
	}
	if subcmd == SCMD_EAT {
		act(libc.CString("You eat some of $p."), FALSE, ch, food, nil, TO_CHAR)
		if food.Action_description != nil {
			act(food.Action_description, FALSE, ch, food, nil, TO_CHAR)
		}
		act(libc.CString("$n eats some of $p."), TRUE, ch, food, nil, TO_ROOM)
	} else {
		act(libc.CString("You nibble a little bit of $p."), FALSE, ch, food, nil, TO_CHAR)
		act(libc.CString("$n tastes a little bit of $p."), TRUE, ch, food, nil, TO_ROOM)
		if (food.Value[VAL_FOOD_POISON]) != 0 {
			send_to_char(ch, libc.CString("It has a sickening taste! Better not eat it..."))
			return
		}
	}
	var foob int = 48 - int(ch.Player_specials.Conditions[HUNGER])
	if subcmd == SCMD_EAT {
		amount = food.Value[VAL_FOOD_FOODVAL]
	} else {
		amount = 1
	}
	gain_condition(ch, HUNGER, amount)
	if ch.Foodr == 0 && subcmd != SCMD_TASTE {
		ch.Move += (ch.Max_move / 100) * int64(amount)
		ch.Foodr = 2
		if OBJ_FLAGGED(food, ITEM_YUM) {
			ch.Move += int64(float64(ch.Max_move) * 0.25)
		}
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		send_to_char(ch, libc.CString("You feel rejuvinated by it.\r\n"))
	}
	if GET_OBJ_VNUM(food) >= MEAL_START && GET_OBJ_VNUM(food) <= MEAL_LAST && int(ch.Player_specials.Conditions[HUNGER]) < 48 && (!ROOM_FLAGGED(ch.In_room, ROOM_AL) && !ROOM_FLAGGED(ch.In_room, ROOM_RHELL)) {
		if subcmd != SCMD_TASTE {
			var (
				psbonus  int = (food.Value[1])
				expbonus int = int(float64(food.Value[2]) * ((float64(GET_LEVEL(ch)) * 0.4) + 1))
				capped   int = FALSE
				pscapped int = FALSE
			)
			if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) <= 0 && GET_LEVEL(ch) < 100 {
				expbonus = 1
				capped = TRUE
			} else if expbonus > GET_LEVEL(ch)*1500 && GET_LEVEL(ch) < 100 {
				expbonus = GET_LEVEL(ch) * 1000
			}
			if (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 500 {
				psbonus = 0
				pscapped = TRUE
			}
			if !AFF_FLAGGED(ch, AFF_PUKED) {
				gain_exp(ch, int64(expbonus))
				ch.Player_specials.Class_skill_points[ch.Chclass] += psbonus
				send_to_char(ch, libc.CString("That was exceptionally delicious! @D[@mPS@D: @C+%d@D] [@gEXP@D: @G+%s@D]@n\r\n"), psbonus, add_commas(int64(expbonus)))
				if capped == TRUE {
					send_to_char(ch, libc.CString("Experience capped due to negative TNL.\r\n"))
				}
				if pscapped == TRUE {
					send_to_char(ch, libc.CString("Practice Sessions capped for food at 500 PS.\r\n"))
				}
			} else {
				send_to_char(ch, libc.CString("You have recently puked. You must wait a while for your body to adjust before excellent food gives you any bonuses.\r\n"))
			}
		}
		if (food.Value[VAL_FOOD_POISON]) == 0 && ch.Hit < gear_pl(ch) && subcmd != SCMD_TASTE {
			var suppress int64 = int64((float64(gear_pl(ch)) * 0.01) * float64(ch.Suppression))
			if food.Weight < 6 {
				ch.Hit += int64(float64(gear_pl(ch)) * 0.05)
			} else {
				ch.Hit += int64(float64(gear_pl(ch)) * 0.1)
			}
			if OBJ_FLAGGED(food, ITEM_YUM) {
				ch.Hit += int64(float64(gear_pl(ch)) * 0.2)
			}
			if ch.Hit > gear_pl(ch) {
				ch.Hit = gear_pl(ch)
			}
			if suppress > 0 {
				if ch.Hit > suppress {
					ch.Hit = suppress
				}
			}
			send_to_char(ch, libc.CString("@MYou feel some of your strength return!@n\r\n"))
		}
	}
	if int(ch.Player_specials.Conditions[HUNGER]) >= 48 && int(ch.Race) != RACE_MAJIN {
		send_to_char(ch, libc.CString("You are full, but may continue to stuff yourself.\r\n"))
	}
	if (food.Value[VAL_FOOD_POISON]) != 0 && !ADM_FLAGGED(ch, ADM_NOPOISON) {
		send_to_char(ch, libc.CString("Oops, that tasted rather strange!\r\n"))
		act(libc.CString("$n coughs and utters some strange sounds."), FALSE, ch, nil, nil, TO_ROOM)
		af.Type = SPELL_POISON
		af.Duration = int16(amount * 2)
		af.Modifier = 0
		af.Location = APPLY_NONE
		af.Bitvector = AFF_POISON
		affect_join(ch, &af, FALSE != 0, FALSE != 0, FALSE != 0, FALSE != 0)
	}
	if subcmd == SCMD_EAT {
		if foob >= (food.Value[VAL_FOOD_FOODVAL]) {
			send_to_char(ch, libc.CString("You finish the last bite.\r\n"))
			if GET_OBJ_VNUM(food) == 53 {
				ch.Move += ch.Max_move / 30
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.01)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, -1)
				}
			}
			if GET_OBJ_VNUM(food) == 93 {
				ch.Move += ch.Max_move / 20
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.01)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 0)
				}
			}
			if GET_OBJ_VNUM(food) == 94 {
				ch.Move += ch.Max_move / 10
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.01)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 1)
				}
			}
			if GET_OBJ_VNUM(food) == 95 {
				ch.Move += ch.Max_move / 10
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.02)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 2)
				}
			}
			extract_obj(food)
		} else {
			food.Value[VAL_FOOD_FOODVAL] -= foob
			if GET_OBJ_VNUM(food) == 53 {
				ch.Move += ch.Max_move / 30
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.01)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, -1)
					extract_obj(food)
				}
			}
			if GET_OBJ_VNUM(food) == 93 {
				ch.Move += ch.Max_move / 20
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.01)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 0)
					extract_obj(food)
				}
			}
			if GET_OBJ_VNUM(food) == 94 {
				ch.Move += ch.Max_move / 10
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.02)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 1)
					extract_obj(food)
				}
			}
			if GET_OBJ_VNUM(food) == 95 {
				ch.Move += ch.Max_move / 10
				ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.03)
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
					ch.Lifeforce = int64(GET_LIFEMAX(ch))
				}
				if OBJ_FLAGGED(food, ITEM_FORGED) {
					send_to_char(ch, libc.CString("This is a forgery. You gain nothing!\r\n"))
				} else {
					send_to_char(ch, libc.CString("You feel more energetic!\r\n"))
					majin_gain(ch, 2)
					extract_obj(food)
				}
			}
		}
	} else {
		if (func() int {
			p := &(food.Value[VAL_FOOD_FOODVAL])
			*p--
			return *p
		}()) == 0 {
			send_to_char(ch, libc.CString("There's nothing left now.\r\n"))
			extract_obj(food)
		}
	}
}
func majin_gain(ch *char_data, type_ int) {
	if int(ch.Race) != RACE_MAJIN || IS_NPC(ch) {
		return
	}
	if !soft_cap(ch, 0) {
		send_to_char(ch, libc.CString("You can not gain anymore from candy consumption at your current level.\r\n"))
		return
	} else if type_ == -1 {
		var (
			st int = int((ch.Basest / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			pl int = int((ch.Basepl / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			ki int = int((ch.Baseki / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
		)
		if st > 300000 {
			st = 300000
		}
		if ki > 300000 {
			ki = 300000
		}
		if pl > 300000 {
			pl = 300000
		}
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_hit += int64(pl * 2)
			ch.Max_mana += int64(ki * 2)
			ch.Max_move += int64(st * 2)
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_hit += int64(pl * 3)
			ch.Max_mana += int64(ki * 3)
			ch.Max_move += int64(st * 3)
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_hit += int64(float64(pl) * 4.5)
			ch.Max_mana += int64(float64(ki) * 4.5)
			ch.Max_move += int64(float64(st) * 4.5)
		} else {
			ch.Max_hit += int64(pl)
			ch.Max_mana += int64(ki)
			ch.Max_move += int64(st)
		}
		ch.Basepl += int64(pl)
		ch.Baseki += int64(ki)
		ch.Basest += int64(st)
		send_to_char(ch, libc.CString("@mYou feel stronger after consuming the candy @D[@RPL@W: @r%s @CKi@D: @c%s @GSt@D: @g%s@D]@m!@n\r\n"), add_commas(int64(pl)), add_commas(int64(ki)), add_commas(int64(st)))
		return
	} else if type_ == 0 {
		var (
			st int = int((ch.Basest / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			pl int = int((ch.Basepl / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			ki int = int((ch.Baseki / 1200) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
		)
		if st > 500000 {
			st = 500000
		}
		if ki > 500000 {
			ki = 500000
		}
		if pl > 500000 {
			pl = 500000
		}
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_hit += int64(pl * 2)
			ch.Max_mana += int64(ki * 2)
			ch.Max_move += int64(st * 2)
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_hit += int64(pl * 3)
			ch.Max_mana += int64(ki * 3)
			ch.Max_move += int64(st * 3)
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_hit += int64(float64(pl) * 4.5)
			ch.Max_mana += int64(float64(ki) * 4.5)
			ch.Max_move += int64(float64(st) * 4.5)
		} else {
			ch.Max_hit += int64(pl)
			ch.Max_mana += int64(ki)
			ch.Max_move += int64(st)
		}
		ch.Basepl += int64(pl)
		ch.Baseki += int64(ki)
		ch.Basest += int64(st)
		send_to_char(ch, libc.CString("@mYou feel stronger after consuming the candy @D[@RPL@W: @r%s @CKi@D: @c%s @GSt@D: @g%s@D]@m!@n\r\n"), add_commas(int64(pl)), add_commas(int64(ki)), add_commas(int64(st)))
		return
	} else if type_ == 1 {
		var (
			st int = int((ch.Basest / 1000) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			pl int = int((ch.Basepl / 1000) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			ki int = int((ch.Baseki / 1000) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
		)
		if st > 1200000 {
			st = 1200000
		}
		if ki > 1200000 {
			ki = 1200000
		}
		if pl > 1200000 {
			pl = 1200000
		}
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_hit += int64(pl * 2)
			ch.Max_mana += int64(ki * 2)
			ch.Max_move += int64(st * 2)
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_hit += int64(pl * 3)
			ch.Max_mana += int64(ki * 3)
			ch.Max_move += int64(st * 3)
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_hit += int64(float64(pl) * 4.5)
			ch.Max_mana += int64(float64(ki) * 4.5)
			ch.Max_move += int64(float64(st) * 4.5)
		} else {
			ch.Max_hit += int64(pl)
			ch.Max_mana += int64(ki)
			ch.Max_move += int64(st)
		}
		ch.Basepl += int64(pl)
		ch.Baseki += int64(ki)
		ch.Basest += int64(st)
		send_to_char(ch, libc.CString("@mYou feel stronger after consuming the candy @D[@RPL@W: @r%s @CKi@D: @c%s @GSt@D: @g%s@D]@m!@n\r\n"), add_commas(int64(pl)), add_commas(int64(ki)), add_commas(int64(st)))
		return
	} else if type_ == 2 {
		var (
			st int = int((ch.Basest / 900) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			pl int = int((ch.Basepl / 900) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
			ki int = int((ch.Baseki / 900) + int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2)))
		)
		if st > 1500000 {
			st = 1500000
		}
		if ki > 1500000 {
			ki = 1500000
		}
		if pl > 1500000 {
			pl = 1500000
		}
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_hit += int64(pl * 2)
			ch.Max_mana += int64(ki * 2)
			ch.Max_move += int64(st * 2)
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_hit += int64(pl * 3)
			ch.Max_mana += int64(ki * 3)
			ch.Max_move += int64(st * 3)
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_hit += int64(float64(pl) * 4.5)
			ch.Max_mana += int64(float64(ki) * 4.5)
			ch.Max_move += int64(float64(st) * 4.5)
		} else {
			ch.Max_hit += int64(pl)
			ch.Max_mana += int64(ki)
			ch.Max_move += int64(st)
		}
		ch.Basepl += int64(pl)
		ch.Baseki += int64(ki)
		ch.Basest += int64(st)
		send_to_char(ch, libc.CString("@mYou feel stronger after consuming the candy @D[@RPL@W: @r%s @CKi@D: @c%s @GSt@D: @g%s@D]@m!@n\r\n"), add_commas(int64(pl)), add_commas(int64(ki)), add_commas(int64(st)))
		return
	}
}
func do_pour(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1     [2048]byte
		arg2     [2048]byte
		from_obj *obj_data = nil
		to_obj   *obj_data = nil
		amount   int       = 0
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if subcmd == SCMD_POUR {
		if arg1[0] == 0 {
			send_to_char(ch, libc.CString("From what do you want to pour?\r\n"))
			return
		}
		if (func() *obj_data {
			from_obj = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
			return from_obj
		}()) == nil {
			send_to_char(ch, libc.CString("You can't find it!\r\n"))
			return
		}
		if int(from_obj.Type_flag) != ITEM_DRINKCON {
			send_to_char(ch, libc.CString("You can't pour from that!\r\n"))
			return
		}
	}
	if subcmd == SCMD_FILL {
		if arg1[0] == 0 {
			send_to_char(ch, libc.CString("What do you want to fill?  And what are you filling it from?\r\n"))
			return
		}
		if (func() *obj_data {
			to_obj = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
			return to_obj
		}()) == nil {
			send_to_char(ch, libc.CString("You can't find it!\r\n"))
			return
		}
		if int(to_obj.Type_flag) != ITEM_DRINKCON {
			act(libc.CString("You can't fill $p!"), FALSE, ch, to_obj, nil, TO_CHAR)
			return
		}
		if arg2[0] == 0 {
			act(libc.CString("What do you want to fill $p from?"), FALSE, ch, to_obj, nil, TO_CHAR)
			return
		}
		if (func() *obj_data {
			from_obj = get_obj_in_list_vis(ch, &arg2[0], nil, world[ch.In_room].Contents)
			return from_obj
		}()) == nil {
			send_to_char(ch, libc.CString("There doesn't seem to be %s %s here.\r\n"), AN(&arg2[0]), &arg2[0])
			return
		}
		if int(from_obj.Type_flag) != ITEM_FOUNTAIN && !OBJ_FLAGGED(from_obj, ITEM_BROKEN) {
			act(libc.CString("You can't fill something from $p."), FALSE, ch, from_obj, nil, TO_CHAR)
			return
		} else if int(from_obj.Type_flag) == ITEM_FOUNTAIN && OBJ_FLAGGED(from_obj, ITEM_BROKEN) {
			act(libc.CString("You can't fill something from a broken fountain."), FALSE, ch, from_obj, nil, TO_CHAR)
			return
		}
	}
	if (from_obj.Value[VAL_DRINKCON_HOWFULL]) == 0 {
		act(libc.CString("The $p is empty."), FALSE, ch, from_obj, nil, TO_CHAR)
		return
	}
	if subcmd == SCMD_POUR {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Where do you want it?  Out or in what?\r\n"))
			return
		}
		if libc.StrCaseCmp(&arg2[0], libc.CString("out")) == 0 {
			if (from_obj.Value[VAL_DRINKCON_CAPACITY]) > 0 {
				act(libc.CString("$n empties $p."), TRUE, ch, from_obj, nil, TO_ROOM)
				act(libc.CString("You empty $p."), FALSE, ch, from_obj, nil, TO_CHAR)
				weight_change_object(from_obj, -(from_obj.Value[VAL_DRINKCON_HOWFULL]))
				name_from_drinkcon(from_obj)
				from_obj.Value[VAL_DRINKCON_HOWFULL] = 0
				from_obj.Value[VAL_DRINKCON_LIQUID] = 0
				from_obj.Value[VAL_DRINKCON_POISON] = 0
			} else {
				send_to_char(ch, libc.CString("You can't possibly pour that container out!\r\n"))
			}
			return
		}
		if (func() *obj_data {
			to_obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
			return to_obj
		}()) == nil {
			send_to_char(ch, libc.CString("You can't find it!\r\n"))
			return
		}
		if int(to_obj.Type_flag) != ITEM_DRINKCON && int(to_obj.Type_flag) != ITEM_FOUNTAIN {
			send_to_char(ch, libc.CString("You can't pour anything into that.\r\n"))
			return
		}
	}
	if to_obj == from_obj {
		send_to_char(ch, libc.CString("A most unproductive effort.\r\n"))
		return
	}
	if (to_obj.Value[VAL_DRINKCON_HOWFULL]) != 0 && (to_obj.Value[VAL_DRINKCON_LIQUID]) != (from_obj.Value[VAL_DRINKCON_LIQUID]) {
		send_to_char(ch, libc.CString("There is already another liquid in it!\r\n"))
		return
	}
	if (to_obj.Value[VAL_DRINKCON_CAPACITY]) < 0 || (to_obj.Value[VAL_DRINKCON_HOWFULL]) >= (to_obj.Value[VAL_DRINKCON_CAPACITY]) {
		send_to_char(ch, libc.CString("There is no room for more.\r\n"))
		return
	}
	if subcmd == SCMD_POUR {
		send_to_char(ch, libc.CString("You pour the %s into the %s."), drinks[from_obj.Value[VAL_DRINKCON_LIQUID]], &arg2[0])
	}
	if subcmd == SCMD_FILL {
		act(libc.CString("You gently fill $p from $P."), FALSE, ch, to_obj, unsafe.Pointer(from_obj), TO_CHAR)
		act(libc.CString("$n gently fills $p from $P."), TRUE, ch, to_obj, unsafe.Pointer(from_obj), TO_ROOM)
	}
	if (to_obj.Value[VAL_DRINKCON_HOWFULL]) == 0 {
		name_to_drinkcon(to_obj, from_obj.Value[VAL_DRINKCON_LIQUID])
	}
	to_obj.Value[VAL_DRINKCON_LIQUID] = from_obj.Value[VAL_DRINKCON_LIQUID]
	if (from_obj.Value[VAL_DRINKCON_CAPACITY]) > 0 {
		from_obj.Value[VAL_DRINKCON_HOWFULL] -= func() int {
			amount = (to_obj.Value[VAL_DRINKCON_CAPACITY]) - (to_obj.Value[VAL_DRINKCON_HOWFULL])
			return amount
		}()
		to_obj.Value[VAL_DRINKCON_HOWFULL] = to_obj.Value[VAL_DRINKCON_CAPACITY]
		if (from_obj.Value[VAL_DRINKCON_HOWFULL]) < 0 {
			to_obj.Value[VAL_DRINKCON_HOWFULL] += from_obj.Value[VAL_DRINKCON_HOWFULL]
			amount += from_obj.Value[VAL_DRINKCON_HOWFULL]
			name_from_drinkcon(from_obj)
			from_obj.Value[VAL_DRINKCON_HOWFULL] = 0
			from_obj.Value[VAL_DRINKCON_LIQUID] = 0
			from_obj.Value[VAL_DRINKCON_POISON] = 0
		}
	} else {
		to_obj.Value[VAL_DRINKCON_HOWFULL] = to_obj.Value[VAL_DRINKCON_CAPACITY]
	}
	to_obj.Value[VAL_DRINKCON_POISON] = int(libc.BoolToInt((to_obj.Value[VAL_DRINKCON_POISON]) != 0 || (from_obj.Value[VAL_DRINKCON_POISON]) != 0))
	if (from_obj.Value[VAL_DRINKCON_CAPACITY]) > 0 {
		weight_change_object(from_obj, -amount)
	}
	weight_change_object(to_obj, amount)
}
func wear_message(ch *char_data, obj *obj_data, where int) {
	var wear_messages [23][2]*byte = [23][2]*byte{{libc.CString("$n lights $p and holds it."), libc.CString("You light $p and hold it.")}, {libc.CString("$n slides $p on to $s right ring finger."), libc.CString("You slide $p on to your right ring finger.")}, {libc.CString("$n slides $p on to $s left ring finger."), libc.CString("You slide $p on to your left ring finger.")}, {libc.CString("$n wears $p around $s neck."), libc.CString("You wear $p around your neck.")}, {libc.CString("$n wears $p around $s neck."), libc.CString("You wear $p around your neck.")}, {libc.CString("$n wears $p on $s body."), libc.CString("You wear $p on your body.")}, {libc.CString("$n wears $p on $s head."), libc.CString("You wear $p on your head.")}, {libc.CString("$n puts $p on $s legs."), libc.CString("You put $p on your legs.")}, {libc.CString("$n wears $p on $s feet."), libc.CString("You wear $p on your feet.")}, {libc.CString("$n puts $p on $s hands."), libc.CString("You put $p on your hands.")}, {libc.CString("$n wears $p on $s arms."), libc.CString("You wear $p on your arms.")}, {libc.CString("$n straps $p around $s arm as a shield."), libc.CString("You start to use $p as a shield.")}, {libc.CString("$n wears $p about $s body."), libc.CString("You wear $p around your body.")}, {libc.CString("$n wears $p around $s waist."), libc.CString("You wear $p around your waist.")}, {libc.CString("$n puts $p on around $s right wrist."), libc.CString("You put $p on around your right wrist.")}, {libc.CString("$n puts $p on around $s left wrist."), libc.CString("You put $p on around your left wrist.")}, {libc.CString("$n wields $p."), libc.CString("You wield $p.")}, {libc.CString("$n grabs $p."), libc.CString("You grab $p.")}, {libc.CString("$n wears $p on $s back."), libc.CString("You wear $p on your back.")}, {libc.CString("$n puts $p in $s right ear."), libc.CString("You put $p in your right ear.")}, {libc.CString("$n puts $p in $s left ear."), libc.CString("You put $p in your left ear.")}, {libc.CString("$n wears $p as a cape."), libc.CString("You wear $p as a cape.")}, {libc.CString("$n covers $s left eye with $p."), libc.CString("You wear $p over your left eye.")}}
	act(wear_messages[where][0], TRUE, ch, obj, nil, TO_ROOM)
	act(wear_messages[where][1], FALSE, ch, obj, nil, TO_CHAR)
}
func hands(ch *char_data) int {
	var x int
	if (ch.Equipment[WEAR_WIELD1]) != nil {
		if OBJ_FLAGGED(ch.Equipment[WEAR_WIELD1], ITEM_2H) || wield_type(get_size(ch), ch.Equipment[WEAR_WIELD1]) == WIELD_TWOHAND {
			x = 2
		} else {
			x = 1
		}
	} else {
		x = 0
	}
	if (ch.Equipment[WEAR_WIELD2]) != nil {
		if OBJ_FLAGGED(ch.Equipment[WEAR_WIELD2], ITEM_2H) || wield_type(get_size(ch), ch.Equipment[WEAR_WIELD2]) == WIELD_TWOHAND {
			x += 2
		} else {
			x += 1
		}
	}
	return x
}
func perform_wear(ch *char_data, obj *obj_data, where int) {
	var (
		wear_bitvectors [23]int   = [23]int{ITEM_WEAR_TAKE, ITEM_WEAR_FINGER, ITEM_WEAR_FINGER, ITEM_WEAR_NECK, ITEM_WEAR_NECK, ITEM_WEAR_BODY, ITEM_WEAR_HEAD, ITEM_WEAR_LEGS, ITEM_WEAR_FEET, ITEM_WEAR_HANDS, ITEM_WEAR_ARMS, ITEM_WEAR_SHIELD, ITEM_WEAR_ABOUT, ITEM_WEAR_WAIST, ITEM_WEAR_WRIST, ITEM_WEAR_WRIST, ITEM_WEAR_TAKE, ITEM_WEAR_TAKE, ITEM_WEAR_PACK, ITEM_WEAR_EAR, ITEM_WEAR_EAR, ITEM_WEAR_SH, ITEM_WEAR_EYE}
		already_wearing [23]*byte = [23]*byte{libc.CString("You're already using a light.\r\n"), libc.CString("YOU SHOULD NEVER SEE THIS MESSAGE.  PLEASE REPORT.\r\n"), libc.CString("You're already wearing something on both of your ring fingers.\r\n"), libc.CString("YOU SHOULD NEVER SEE THIS MESSAGE.  PLEASE REPORT.\r\n"), libc.CString("You can't wear anything else around your neck.\r\n"), libc.CString("You're already wearing something on your body.\r\n"), libc.CString("You're already wearing something on your head.\r\n"), libc.CString("You're already wearing something on your legs.\r\n"), libc.CString("You're already wearing something on your feet.\r\n"), libc.CString("You're already wearing something on your hands.\r\n"), libc.CString("You're already wearing something on your arms.\r\n"), libc.CString("You're already using a shield.\r\n"), libc.CString("You're already wearing something about your body.\r\n"), libc.CString("You already have something around your waist.\r\n"), libc.CString("YOU SHOULD NEVER SEE THIS MESSAGE.  PLEASE REPORT.\r\n"), libc.CString("You're already wearing something around both of your wrists.\r\n"), libc.CString("You're already wielding a weapon.\r\n"), libc.CString("You're already holding something.\r\n"), libc.CString("You're already wearing something on your back.\r\n"), libc.CString("YOU SHOULD NEVER SEE THIS MESSAGE.  PLEASE REPORT.\r\n"), libc.CString("You're already wearing something in both ears.\r\n"), libc.CString("You're already wearing something on your shoulders.\r\n"), libc.CString("You're already wearing something as a scouter.\r\n")}
	)
	if !OBJWEAR_FLAGGED(obj, bitvector_t(int32(wear_bitvectors[where]))) {
		act(libc.CString("You can't wear $p there."), FALSE, ch, obj, nil, TO_CHAR)
		return
	}
	if !BODY_FLAGGED(ch, bitvector_t(int32(where))) {
		send_to_char(ch, libc.CString("Seems like your body type doesn't really allow that.\r\n"))
		return
	}
	if where == WEAR_FINGER_R || where == WEAR_NECK_1 || where == WEAR_WRIST_R || where == WEAR_EAR_R || where == WEAR_WIELD1 {
		if (ch.Equipment[where]) != nil {
			where++
		}
	}
	if (OBJ_FLAGGED(obj, ITEM_2H) || wield_type(get_size(ch), obj) == WIELD_TWOHAND) && hands(ch) > 0 {
		send_to_char(ch, libc.CString("Seems like you might not have enough free hands.\r\n"))
		return
	}
	if where == WEAR_WIELD2 && PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("Seems like you might not have enough free hands.\r\n"))
		return
	}
	if (where == WEAR_WIELD1 || where == WEAR_WIELD2) && hands(ch) > 1 {
		send_to_char(ch, libc.CString("Seems like you might not have enough free hands.\r\n"))
		return
	}
	if (ch.Equipment[where]) != nil {
		send_to_char(ch, libc.CString("%s"), already_wearing[where])
		return
	}
	if wear_otrigger(obj, ch, where) == 0 || obj.Carried_by != ch {
		return
	}
	if int(obj.Type_flag) == ITEM_WEAPON && OBJ_FLAGGED(obj, ITEM_CUSTOM) {
		if GET_LEVEL(ch) < 20 {
			send_to_char(ch, libc.CString("You are not experienced enough to hold that.\r\n"))
			return
		}
	}
	if int(obj.Type_flag) != ITEM_LIGHT && obj.Size > get_size(ch) {
		send_to_char(ch, libc.CString("Seems like it is too big for you.\r\n"))
		return
	}
	if obj.Size < get_size(ch) && int(obj.Type_flag) != ITEM_LIGHT {
		send_to_char(ch, libc.CString("Seems like it is too small for you.\r\n"))
		return
	}
	wear_message(ch, obj, where)
	obj_from_char(obj)
	equip_char(ch, obj, where)
}
func find_eq_pos(ch *char_data, obj *obj_data, arg *byte) int {
	var (
		where    int       = -1
		keywords [24]*byte = [24]*byte{libc.CString("!RESERVED!"), libc.CString("finger"), libc.CString("!RESERVED!"), libc.CString("neck"), libc.CString("!RESERVED!"), libc.CString("body"), libc.CString("head"), libc.CString("legs"), libc.CString("feet"), libc.CString("hands"), libc.CString("arms"), libc.CString("shield"), libc.CString("about"), libc.CString("waist"), libc.CString("wrist"), libc.CString("!RESERVED!"), libc.CString("!RESERVED!"), libc.CString("!RESERVED!"), libc.CString("back"), libc.CString("ear"), libc.CString("\r!RESERVED!"), libc.CString("shoulders"), libc.CString("scouter"), libc.CString("\n")}
	)
	if arg == nil || *arg == 0 {
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_FINGER) && BODY_FLAGGED(ch, WEAR_FINGER_R) {
			return WEAR_FINGER_R
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_NECK) && BODY_FLAGGED(ch, WEAR_NECK_1) {
			return WEAR_NECK_1
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_BODY) && BODY_FLAGGED(ch, WEAR_BODY) {
			return WEAR_BODY
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_HEAD) && BODY_FLAGGED(ch, WEAR_HEAD) {
			return WEAR_HEAD
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_LEGS) && BODY_FLAGGED(ch, WEAR_LEGS) {
			return WEAR_LEGS
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_FEET) && BODY_FLAGGED(ch, WEAR_FEET) {
			return WEAR_FEET
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_HANDS) && BODY_FLAGGED(ch, WEAR_HANDS) {
			return WEAR_HANDS
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_ARMS) && BODY_FLAGGED(ch, WEAR_ARMS) {
			return WEAR_ARMS
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_SHIELD) && BODY_FLAGGED(ch, WEAR_WIELD2) {
			return WEAR_WIELD2
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_ABOUT) && BODY_FLAGGED(ch, WEAR_ABOUT) {
			return WEAR_ABOUT
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_WAIST) && BODY_FLAGGED(ch, WEAR_WAIST) {
			return WEAR_WAIST
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_WRIST) && BODY_FLAGGED(ch, WEAR_WRIST_R) {
			return WEAR_WRIST_R
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_HOLD) && BODY_FLAGGED(ch, WEAR_WIELD2) {
			return WEAR_WIELD2
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_PACK) && BODY_FLAGGED(ch, WEAR_BACKPACK) {
			return WEAR_BACKPACK
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_EAR) && BODY_FLAGGED(ch, WEAR_EAR_R) {
			return WEAR_EAR_R
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_SH) && BODY_FLAGGED(ch, WEAR_SH) {
			return WEAR_SH
		}
		if OBJWEAR_FLAGGED(obj, ITEM_WEAR_EYE) && BODY_FLAGGED(ch, WEAR_EYE) {
			return WEAR_EYE
		}
	} else if (func() int {
		where = search_block(arg, &keywords[0], FALSE)
		return where
	}()) < 0 || *arg == '!' {
		send_to_char(ch, libc.CString("'%s'?  What part of your body is THAT?\r\n"), arg)
		return -1
	}
	return where
}
func do_wear(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1       [2048]byte
		arg2       [2048]byte
		obj        *obj_data
		next_obj   *obj_data
		where      int
		dotmode    int
		items_worn int = 0
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 {
		send_to_char(ch, libc.CString("Wear what?\r\n"))
		return
	}
	dotmode = find_all_dots(&arg1[0])
	if arg2[0] != 0 && dotmode != FIND_INDIV {
		send_to_char(ch, libc.CString("You can't specify the same body location for more than one item!\r\n"))
		return
	}
	if dotmode == FIND_ALL {
		for obj = ch.Carrying; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if CAN_SEE_OBJ(ch, obj) && (func() int {
				where = find_eq_pos(ch, obj, nil)
				return where
			}()) >= 0 {
				if GET_LEVEL(ch) < obj.Level {
					act(libc.CString("$p: you are not experienced enough to use that."), FALSE, ch, obj, nil, TO_CHAR)
					send_to_char(ch, libc.CString("You need to be at least %d level to use it.\r\n"), obj.Level)
				} else if OBJ_FLAGGED(obj, ITEM_BROKEN) {
					act(libc.CString("$p: it seems to be broken."), FALSE, ch, obj, nil, TO_CHAR)
				} else if OBJ_FLAGGED(obj, ITEM_FORGED) {
					act(libc.CString("$p: it seems to be fake..."), FALSE, ch, obj, nil, TO_CHAR)
				} else {
					items_worn++
					if is_proficient_with_armor(ch, obj.Value[VAL_ARMOR_SKILL]) == 0 && int(obj.Type_flag) == ITEM_ARMOR {
						send_to_char(ch, libc.CString("You have no proficiency with this type of armor.\r\nYour fighting and physical skills will be greatly impeded.\r\n"))
					}
					perform_wear(ch, obj, where)
				}
			}
		}
		if items_worn == 0 {
			send_to_char(ch, libc.CString("You don't seem to have anything wearable.\r\n"))
		}
	} else if dotmode == FIND_ALLDOT {
		if arg1[0] == 0 {
			send_to_char(ch, libc.CString("Wear all of what?\r\n"))
			return
		}
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't seem to have any %ss.\r\n"), &arg1[0])
		} else if GET_LEVEL(ch) < obj.Level {
			send_to_char(ch, libc.CString("You are not experienced enough to use that.\r\n"))
		} else {
			for obj != nil {
				next_obj = get_obj_in_list_vis(ch, &arg1[0], nil, obj.Next_content)
				if (func() int {
					where = find_eq_pos(ch, obj, nil)
					return where
				}()) >= 0 {
					if is_proficient_with_armor(ch, obj.Value[VAL_ARMOR_SKILL]) == 0 && int(obj.Type_flag) == ITEM_ARMOR {
						send_to_char(ch, libc.CString("You have no proficiency with this type of armor.\r\nYour fighting and physical skills will be greatly impeded.\r\n"))
					}
					perform_wear(ch, obj, where)
				} else {
					act(libc.CString("You can't wear $p."), FALSE, ch, obj, nil, TO_CHAR)
				}
				obj = next_obj
			}
		}
	} else {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg1[0]), &arg1[0])
		} else if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("But it seems to be broken!\r\n"))
		} else if OBJ_FLAGGED(obj, ITEM_FORGED) {
			send_to_char(ch, libc.CString("But it seems to be fake!\r\n"))
		} else if GET_LEVEL(ch) < obj.Level {
			send_to_char(ch, libc.CString("You are not experienced enough to use that.\r\n"))
		} else {
			if (func() int {
				where = find_eq_pos(ch, obj, &arg2[0])
				return where
			}()) >= 0 {
				if is_proficient_with_armor(ch, obj.Value[VAL_ARMOR_SKILL]) == 0 && int(obj.Type_flag) == ITEM_ARMOR {
					send_to_char(ch, libc.CString("You have no proficiency with this type of armor.\r\nYour fighting and physical skills will be greatly impeded.\r\n"))
				}
				perform_wear(ch, obj, where)
			} else if arg2[0] == 0 {
				act(libc.CString("You can't wear $p."), FALSE, ch, obj, nil, TO_CHAR)
			}
		}
	}
}
func do_wield(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		obj *obj_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Wield what?\r\n"))
	} else if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
	} else {
		if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_WIELD) {
			send_to_char(ch, libc.CString("You can't wield that.\r\n"))
		} else if obj.Weight > max_carry_weight(ch) {
			send_to_char(ch, libc.CString("It's too heavy for you to use.\r\n"))
		} else if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("But it seems to be broken!\r\n"))
		} else if OBJ_FLAGGED(obj, ITEM_FORGED) {
			send_to_char(ch, libc.CString("But it seems to be fake!\r\n"))
		} else if GET_LEVEL(ch) < obj.Level {
			send_to_char(ch, libc.CString("You are not experienced enough to use that.\r\n"))
		} else if PLR_FLAGGED(ch, PLR_THANDW) {
			send_to_char(ch, libc.CString("You are holding a weapon with two hands right now!\r\n"))
		} else {
			if !IS_NPC(ch) && is_proficient_with_weapon(ch, obj.Value[VAL_WEAPON_SKILL]) == 0 && int(obj.Type_flag) == ITEM_ARMOR {
				send_to_char(ch, libc.CString("You have no proficiency with this type of weapon.\r\nYour attack accuracy will be greatly reduced.\r\n"))
			}
			perform_wear(ch, obj, WEAR_WIELD1)
		}
	}
}
func do_grab(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		obj *obj_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Hold what?\r\n"))
	} else if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
	} else if GET_LEVEL(ch) < obj.Level {
		send_to_char(ch, libc.CString("You are not experienced enough to use that.\r\n"))
	} else if PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("You are wielding a weapon with both hands currently.\r\n"))
	} else {
		if int(obj.Type_flag) == ITEM_LIGHT {
			perform_wear(ch, obj, WEAR_WIELD2)
			if (obj.Value[VAL_LIGHT_HOURS]) > 0 || (obj.Value[VAL_LIGHT_HOURS]) < 0 {
				act(libc.CString("@wYou light $p@w."), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@w lights $p@w."), TRUE, ch, obj, nil, TO_ROOM)
			}
			if (obj.Value[VAL_LIGHT_HOURS]) == 0 {
				act(libc.CString("@wYou try to light $p@w but it is burnt out."), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@w tries to light $p@w but nothing happens."), TRUE, ch, obj, nil, TO_ROOM)
			}
		} else {
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_HOLD) && int(obj.Type_flag) != ITEM_WAND && int(obj.Type_flag) != ITEM_STAFF && int(obj.Type_flag) != ITEM_SCROLL && int(obj.Type_flag) != ITEM_POTION {
				send_to_char(ch, libc.CString("You can't hold that.\r\n"))
			} else {
				perform_wear(ch, obj, WEAR_WIELD2)
			}
		}
	}
}
func perform_remove(ch *char_data, pos int) {
	var (
		obj      *obj_data
		previous int64 = ch.Hit
	)
	if (func() *obj_data {
		obj = ch.Equipment[pos]
		return obj
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: perform_remove: bad pos %d passed."), pos)
	} else if OBJ_FLAGGED(obj, ITEM_NODROP) && ch.Admlevel < 1 {
		act(libc.CString("You can't remove $p, it must be CURSED!"), FALSE, ch, obj, nil, TO_CHAR)
	} else if int(ch.Carry_items) >= 50 {
		act(libc.CString("$p: your arms are full!"), FALSE, ch, obj, nil, TO_CHAR)
	} else {
		if remove_otrigger(obj, ch) == 0 {
			return
		}
		if pos == WEAR_WIELD1 && PLR_FLAGGED(ch, PLR_THANDW) {
			REMOVE_BIT_AR(ch.Act[:], PLR_THANDW)
		}
		obj_to_char(unequip_char(ch, pos), ch)
		act(libc.CString("You stop using $p."), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("$n stops using $p."), TRUE, ch, obj, nil, TO_ROOM)
		if previous > ch.Hit {
			var drop [2048]byte
			stdio.Sprintf(&drop[0], "@RYour powerlevel has dropped from removing $p@R! @D[@r-%s@D]\r\n", add_commas(previous-ch.Hit))
			act(&drop[0], FALSE, ch, obj, nil, TO_CHAR)
		}
	}
}
func do_remove(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj     *obj_data
		arg     [2048]byte
		i       int
		dotmode int
		found   int = 0
		msg     int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Remove what?\r\n"))
		return
	}
	obj = find_obj_in_list_type(ch.Carrying, ITEM_BOARD)
	if obj == nil {
		obj = find_obj_in_list_type(world[ch.In_room].Contents, ITEM_BOARD)
	}
	if found != 0 {
		if !unicode.IsDigit(rune(arg[0])) || (func() int {
			msg = libc.Atoi(libc.GoString(&arg[0]))
			return msg
		}()) == 0 {
			found = 0
		} else {
			remove_board_msg(GET_OBJ_VNUM(obj), ch, msg)
		}
	}
	if found == 0 {
		dotmode = find_all_dots(&arg[0])
		if dotmode == FIND_ALL {
			found = 0
			for i = 0; i < NUM_WEARS; i++ {
				if (ch.Equipment[i]) != nil {
					perform_remove(ch, i)
					found = 1
				}
			}
			if found == 0 {
				send_to_char(ch, libc.CString("You're not using anything.\r\n"))
			}
		} else if dotmode == FIND_ALLDOT {
			if arg[0] == 0 {
				send_to_char(ch, libc.CString("Remove all of what?\r\n"))
			} else {
				found = 0
				for i = 0; i < NUM_WEARS; i++ {
					if (ch.Equipment[i]) != nil && CAN_SEE_OBJ(ch, ch.Equipment[i]) && isname(&arg[0], (ch.Equipment[i]).Name) != 0 {
						perform_remove(ch, i)
						found = 1
					}
				}
				if found == 0 {
					send_to_char(ch, libc.CString("You don't seem to be using any %ss.\r\n"), &arg[0])
				}
			}
		} else {
			if (func() int {
				i = get_obj_pos_in_equip_vis(ch, &arg[0], nil, ch.Equipment[:])
				return i
			}()) < 0 {
				send_to_char(ch, libc.CString("You don't seem to be using %s %s.\r\n"), AN(&arg[0]), &arg[0])
			} else {
				perform_remove(ch, i)
			}
		}
	}
}
func do_sac(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		j   *obj_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Sacrifice what?\n\r"))
		return
	}
	if (func() *obj_data {
		j = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
		return j
	}()) == nil && (func() *obj_data {
		j = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return j
	}()) == nil {
		send_to_char(ch, libc.CString("It doesn't seem to be here.\n\r"))
		return
	}
	if !OBJWEAR_FLAGGED(j, ITEM_WEAR_TAKE) {
		send_to_char(ch, libc.CString("You can't sacrifice that!\n\r"))
		return
	}
	act(libc.CString("$n sacrifices $p."), FALSE, ch, j, nil, TO_ROOM)
	if j.Cost == 0 && !IS_CORPSE(j) {
		send_to_char(ch, libc.CString("Zizazat mocks your sacrifice. Try again, try harder.\r\n"))
		return
	}
	if !IS_CORPSE(j) {
		switch rand_number(0, 5) {
		case 0:
			send_to_char(ch, libc.CString("You sacrifice %s to the Gods.\r\nYou receive one zenni for your humility.\r\n"), j.Short_description)
			ch.Gold += 1
		case 1:
			send_to_char(ch, libc.CString("You sacrifice %s to the Gods.\r\nThe Gods ignore your sacrifice.\r\n"), j.Short_description)
		case 2:
			send_to_char(ch, libc.CString("You sacrifice %s to the Gods.\r\nZizazat gives you %d experience points.\r\n"), j.Short_description, j.Cost*2)
			ch.Exp += int64(j.Cost * 2)
		case 3:
			send_to_char(ch, libc.CString("You sacrifice %s to the Gods.\r\nYou receive %d experience points.\r\n"), j.Short_description, j.Cost)
			ch.Exp += int64(j.Cost)
		case 4:
			send_to_char(ch, libc.CString("Your sacrifice to the Gods is rewarded with %d zenni.\r\n"), j.Cost)
			ch.Gold += j.Cost
		case 5:
			send_to_char(ch, libc.CString("Your sacrifice to the Gods is rewarded with %d zenni\r\n"), j.Cost*2)
			ch.Gold += j.Cost * 2
		default:
			send_to_char(ch, libc.CString("You sacrifice %s to the Gods.\r\nYou receive one zenni for your humility.\r\n"), j.Short_description)
			ch.Gold += 1
		}
	} else {
		send_to_char(ch, libc.CString("You send the corpse on to the next life!\r\n"))
	}
	extract_obj(j)
}

var max_carry_load [31]int = [31]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 115, 130, 150, 175, 200, 230, 260, 300, 350, 400, 460, 520, 600, 700, 800, 920, 1040, 1200, 1400, 1640}

func max_carry_weight(ch *char_data) int64 {
	var (
		abil  int64
		total int
	)
	abil = (ch.Max_hit / 200) + int64(int(ch.Aff_abils.Str)*50)
	total = 1
	return int64(total * int(abil))
}
