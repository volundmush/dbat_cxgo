package main

import "github.com/gotranspile/cxgo/runtime/libc"

func grow_plants() {
	var k *obj_data
	for k = object_list; k != nil; k = k.Next {
		if k.In_room == int(-1) {
			continue
		} else if k.Carried_by != nil || k.In_obj != nil {
			continue
		} else if ROOM_FLAGGED(k.In_room, ROOM_GARDEN1) || ROOM_FLAGGED(k.In_room, ROOM_GARDEN2) {
			if (k.Value[VAL_WATERLEVEL]) < 0 && (k.Value[VAL_WATERLEVEL]) > -10 {
				k.Value[VAL_WATERLEVEL] -= 1
				if (k.Value[VAL_WATERLEVEL]) > -10 {
					send_to_room(k.In_room, libc.CString("%s@y withers a bit.\r\n"), k.Short_description)
				} else {
					send_to_room(k.In_room, libc.CString("%s@y has withered to a dried up dead husk.\r\n"), k.Short_description)
				}
			} else if (k.Value[VAL_WATERLEVEL]) >= 0 {
				k.Value[VAL_WATERLEVEL] -= 1
				if (k.Value[VAL_GROWTH]) < (k.Value[VAL_MATGOAL]) && (k.Value[VAL_MATURITY]) < (k.Value[VAL_MAXMATURE]) {
					k.Value[VAL_GROWTH] += 1
					if (k.Value[VAL_GROWTH]) >= (k.Value[VAL_MATGOAL]) {
						k.Value[VAL_GROWTH] = 0
						k.Value[VAL_MATURITY] += 1
					}
					if (k.Value[VAL_MATURITY]) >= (k.Value[VAL_MAXMATURE]) {
						send_to_room(k.In_room, libc.CString("%s@G is now fully grown!@n\r\n"), k.Short_description)
					}
				}
			}
		}
	}
}
func weather_and_time(mode int) {
	another_hour(mode)
	grow_plants()
	if mode != 0 {
		weather_change()
	}
}
func another_hour(mode int) {
	time_info.Hours++
	if mode != 0 {
		switch time_info.Hours {
		case 4:
			if MOON_DATE() {
				send_to_moon(libc.CString("The full moon disappears.\r\n"))
				MOON_UP = 0
				oozaru_drop(nil)
			} else if time_info.Day == 22 {
				send_to_moon(libc.CString("The full moon disappears.\r\n"))
				MOON_UP = 0
				oozaru_drop(nil)
			}
		case 5:
			weather_info.Sunlight = SUN_RISE
			send_to_outdoor(libc.CString("The sun rises in the east.\r\n"))
			if time_info.Day <= 14 {
				star_phase(nil, 1)
			} else if time_info.Day <= 21 {
				star_phase(nil, 2)
			} else {
				star_phase(nil, 0)
			}
		case 6:
			weather_info.Sunlight = SUN_LIGHT
			send_to_outdoor(libc.CString("The day has begun.\r\n"))
		case 19:
			weather_info.Sunlight = SUN_SET
			send_to_outdoor(libc.CString("The sun slowly disappears in the west.\r\n"))
		case 20:
			weather_info.Sunlight = SUN_DARK
			send_to_outdoor(libc.CString("The night has begun.\r\n"))
		case 21:
			if MOON_DATE() {
				send_to_moon(libc.CString("The full moon has risen.\r\n"))
				MOON_UP = 1
				oozaru_add(nil)
			}
		default:
		}
	}
	if time_info.Hours > 23 {
		time_info.Hours -= 24
		time_info.Day++
		if time_info.Day > 29 {
			time_info.Day = 0
			time_info.Month++
			if time_info.Month > 11 {
				time_info.Month = 0
				time_info.Year++
			}
		}
	}
}
func weather_change() {
	var (
		diff   int
		change int
	)
	if time_info.Month >= 9 && time_info.Month <= 16 {
		if weather_info.Pressure > 985 {
			diff = -2
		} else {
			diff = 2
		}
	} else if weather_info.Pressure > 1015 {
		diff = -2
	} else {
		diff = 2
	}
	weather_info.Change += dice(1, 4)*diff + dice(2, 6) - dice(2, 6)
	weather_info.Change = int(MIN(int64(weather_info.Change), 12))
	weather_info.Change = int(MAX(int64(weather_info.Change), -12))
	weather_info.Pressure += weather_info.Change
	weather_info.Pressure = int(MIN(int64(weather_info.Pressure), 1040))
	weather_info.Pressure = int(MAX(int64(weather_info.Pressure), 960))
	change = 0
	switch weather_info.Sky {
	case SKY_CLOUDLESS:
		if weather_info.Pressure < 990 {
			change = 1
		} else if weather_info.Pressure < 1010 {
			if dice(1, 4) == 1 {
				change = 1
			}
		}
	case SKY_CLOUDY:
		if weather_info.Pressure < 970 {
			change = 2
		} else if weather_info.Pressure < 990 {
			if dice(1, 4) == 1 {
				change = 2
			} else {
				change = 0
			}
		} else if weather_info.Pressure > 1030 {
			if dice(1, 4) == 1 {
				change = 3
			}
		}
	case SKY_RAINING:
		if weather_info.Pressure < 970 {
			if dice(1, 4) == 1 {
				change = 4
			} else {
				change = 0
			}
		} else if weather_info.Pressure > 1030 {
			change = 5
		} else if weather_info.Pressure > 1010 {
			if dice(1, 4) == 1 {
				change = 5
			}
		}
	case SKY_LIGHTNING:
		if weather_info.Pressure > 1010 {
			change = 6
		} else if weather_info.Pressure > 990 {
			if dice(1, 4) == 1 {
				change = 6
			}
		}
	default:
		change = 0
		weather_info.Sky = SKY_CLOUDLESS
	}
	switch change {
	case 0:
	case 1:
		send_to_outdoor(libc.CString("The sky starts to get cloudy.\r\n"))
		weather_info.Sky = SKY_CLOUDY
	case 2:
		send_to_outdoor(libc.CString("It starts to rain.\r\n"))
		weather_info.Sky = SKY_RAINING
	case 3:
		send_to_outdoor(libc.CString("The clouds disappear.\r\n"))
		weather_info.Sky = SKY_CLOUDLESS
	case 4:
		send_to_outdoor(libc.CString("Lightning starts to show in the sky.\r\n"))
		weather_info.Sky = SKY_LIGHTNING
	case 5:
		send_to_outdoor(libc.CString("The rain stops.\r\n"))
		weather_info.Sky = SKY_CLOUDY
	case 6:
		send_to_outdoor(libc.CString("The lightning stops.\r\n"))
		weather_info.Sky = SKY_RAINING
	default:
	}
}
func oozaru_add(tch *char_data) {
	var d *descriptor_data
	if tch == nil {
		for d = descriptor_list; d != nil; d = d.Next {
			if !IS_PLAYING(d) {
				continue
			}
			if MOON_OK(d.Character) && !PLR_FLAGGED(d.Character, PLR_OOZARU) {
				act(libc.CString("@rLooking up at the moon your heart begins to beat loudly. Sudden rage begins to fill your mind while your body begins to grow. Hair sprouts  all over your body and your teeth become sharp as your body takes on the Oozaru form!@n"), 1, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@R$n@r looks up at the moon as $s eyes turn red and $s heart starts to beat loudly. Hair starts to grow all over $s body as $e starts screaming. The scream turns into a roar as $s body begins to grow into a giant ape!@n"), 1, d.Character, nil, nil, TO_ROOM)
				if d.Character.Kaioken > 0 {
					do_kaioken(d.Character, libc.CString("0"), 0, 0)
				}
				SET_BIT_AR(d.Character.Act[:], PLR_OOZARU)
				var add int = 10000
				var mult int = 2
				d.Character.Max_hit = (d.Character.Basepl + int64(add)) * int64(mult)
				if (d.Character.Hit+int64(add))*int64(mult) <= d.Character.Max_hit {
					d.Character.Hit = (d.Character.Hit + int64(add)) * int64(mult)
				} else if (d.Character.Hit+int64(add))*int64(mult) > d.Character.Max_hit {
					d.Character.Hit = d.Character.Max_hit
				}
				d.Character.Max_mana = (d.Character.Baseki + int64(add)) * int64(mult)
				if (d.Character.Mana+int64(add))*int64(mult) <= d.Character.Max_mana {
					d.Character.Mana = (d.Character.Mana + int64(add)) * int64(mult)
				} else if (d.Character.Mana+int64(add))*int64(mult) > d.Character.Max_mana {
					d.Character.Mana = d.Character.Max_mana
				}
				d.Character.Max_move = (d.Character.Basest + int64(add)) * int64(mult)
				if (d.Character.Move+int64(add))*int64(mult) <= d.Character.Max_move {
					d.Character.Move = (d.Character.Move + int64(add)) * int64(mult)
				} else if (d.Character.Move+int64(add))*int64(mult) > d.Character.Max_move {
					d.Character.Move = d.Character.Max_move
				}
				continue
			}
		}
	} else {
		if MOON_OK(tch) && !PLR_FLAGGED(tch, PLR_OOZARU) {
			act(libc.CString("@rLooking up at the moon your heart begins to beat loudly. Sudden rage begins to fill your mind while your body begins to grow. Hair sprouts  all over your body and your teeth become sharp as your body takes on the Oozaru form!@n"), 1, tch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n@r looks up at the moon as $s eyes turn red and $s heart starts to beat loudly. Hair starts to grow all over $s body as $e starts screaming. The scream turns into a roar as $s body begins to grow into a giant ape!@n"), 1, tch, nil, nil, TO_ROOM)
			SET_BIT_AR(tch.Act[:], PLR_OOZARU)
			var add int = 10000
			var mult int = 2
			tch.Max_hit = (tch.Basepl + int64(add)) * int64(mult)
			if (tch.Hit+int64(add))*int64(mult) <= tch.Max_hit {
				tch.Hit = (tch.Hit + int64(add)) * int64(mult)
			} else if (tch.Hit+int64(add))*int64(mult) > tch.Max_hit {
				tch.Hit = tch.Max_hit
			}
			tch.Max_mana = (tch.Baseki + int64(add)) * int64(mult)
			if (tch.Mana+int64(add))*int64(mult) <= tch.Max_mana {
				tch.Mana = (tch.Mana + int64(add)) * int64(mult)
			} else if (tch.Mana+int64(add))*int64(mult) > tch.Max_mana {
				tch.Mana = tch.Max_mana
			}
			tch.Max_move = (tch.Basest + int64(add)) * int64(mult)
			if (tch.Move+int64(add))*int64(mult) <= tch.Max_move {
				tch.Move = (tch.Move + int64(add)) * int64(mult)
			} else if (tch.Move+int64(add))*int64(mult) > tch.Max_move {
				tch.Move = tch.Max_move
			}
		}
	}
}
func oozaru_drop(tch *char_data) {
	var d *descriptor_data
	if tch == nil {
		for d = descriptor_list; d != nil; d = d.Next {
			if !IS_PLAYING(d) {
				continue
			}
			if PLR_FLAGGED(d.Character, PLR_OOZARU) {
				act(libc.CString("@CYour body begins to shrink back to its normal form as the power of the Oozaru leaves you. You fall asleep shortly after returning to normal!@n"), 1, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C's body begins to shrink and return to normal. Their giant ape features fading back into humanoid features until $e is left normal and asleep.@n"), 1, d.Character, nil, nil, TO_ROOM)
				REMOVE_BIT_AR(d.Character.Act[:], PLR_OOZARU)
				d.Character.Position = POS_SLEEPING
				d.Character.Hit = (d.Character.Hit / 2) - 10000
				d.Character.Mana = (d.Character.Mana / 2) - 10000
				d.Character.Move = (d.Character.Move / 2) - 10000
				d.Character.Max_hit = d.Character.Basepl
				d.Character.Max_mana = d.Character.Baseki
				d.Character.Max_move = d.Character.Basest
				if d.Character.Move < 1 {
					d.Character.Move = 1
				}
				if d.Character.Mana < 1 {
					d.Character.Mana = 1
				}
				if d.Character.Hit < 1 {
					d.Character.Hit = 1
				}
			}
		}
	} else {
		if PLR_FLAGGED(tch, PLR_OOZARU) {
			act(libc.CString("@CYour body begins to shrink back to its normal form as the power of the Oozaru leaves you. You fall asleep shortly after returning to normal!@n"), 1, tch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C's body begins to shrink and return to normal. Their giant ape features fading back into humanoid features until $e is left normal and asleep.@n"), 1, tch, nil, nil, TO_ROOM)
			REMOVE_BIT_AR(tch.Act[:], PLR_OOZARU)
			tch.Position = POS_SLEEPING
			tch.Hit = (tch.Hit / 2) - 10000
			tch.Mana = (tch.Mana / 2) - 10000
			tch.Move = (tch.Move / 2) - 10000
			tch.Max_hit = tch.Basepl
			tch.Max_mana = tch.Baseki
			tch.Max_move = tch.Basest
			if tch.Move < 1 {
				tch.Move = 1
			}
			if tch.Mana < 1 {
				tch.Mana = 1
			}
			if tch.Hit < 1 {
				tch.Hit = 1
			}
		}
	}
}
func star_phase(ch *char_data, type_ int) {
	var d *descriptor_data
	if ch == nil {
		for d = descriptor_list; d != nil; d = d.Next {
			if !IS_PLAYING(d) {
				continue
			}
			if IS_NPC(d.Character) {
				continue
			}
			if GET_LEVEL(d.Character) < 2 {
				continue
			}
			if int(d.Character.Race) == RACE_HOSHIJIN {
				ch = d.Character
				switch type_ {
				case 0:
					if ch.Starphase > 0 {
						act(libc.CString("@WYour eyes and the glyphs on your skin slowly start to lose their glow. You feel the power received from the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W drain away from your body. It has apparently entered the @rDeath Phase@W of its cycle...@n"), 1, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@W's eyes and the glyphs on $s skin slowly start to lose their glow. You notice that $e seems weaker now for some reason.@n"), 1, ch, nil, nil, TO_ROOM)
						phase_powerup(ch, 0, ch.Starphase)
					}
				case 1:
					if ch.Starphase != 1 {
						act(libc.CString("@WYou suddenly feel a @RSURGE@W of power through your body. You feel the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W come into its @CBirth Phase@W and its power is flowing into your body! Finally your eyes and the glyphs on your skin begin to glow an electric @bb@Bl@Cue@W!@n"), 1, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@W suddenly seems to grow stronger for some reason. You notice $s eyes begin to glow an electric @bb@Bl@Cue@W. Suddenly glyphs start to appear all over $s skin and glow with the same light!@n"), 1, ch, nil, nil, TO_ROOM)
						phase_powerup(ch, 0, ch.Starphase)
						phase_powerup(ch, 1, 1)
					}
				case 2:
					if ch.Starphase != 2 {
						act(libc.CString("@WYou suddenly feel a @RSURGE@W of power through your body. You feel the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W come into its @GLife Phase@W and its power is flowing into your body! Finally your eyes and the glyphs on your skin begin to glow an fiery @Rr@re@Rd@W!@n"), 1, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@W suddenly seems to grow stronger for some reason. You notice $s eyes begin to glow a fiery @rR@Re@rd@W. Suddenly glyphs start to appear all over $s skin and glow with the same light!@n"), 1, ch, nil, nil, TO_ROOM)
						phase_powerup(ch, 0, ch.Starphase)
						phase_powerup(ch, 1, 2)
					}
				default:
					send_to_imm(libc.CString("Strange Error in star_phase by: %s"), GET_NAME(ch))
				}
			}
		}
		return
	} else if ch != nil && !IS_NPC(ch) && GET_LEVEL(ch) > 1 {
		if int(ch.Race) == RACE_HOSHIJIN {
			switch type_ {
			case 0:
				if ch.Starphase > 0 {
					act(libc.CString("@WYour eyes and the glyphs on your skin slowly start to lose their glow. You feel the power received from the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W drain away from your body. It has apparently entered the @rDeath Phase@W of its cycle...@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@c$n@W's eyes and the glyphs on $s skin slowly start to lose their glow. You notice that $e seems weaker now for some reason.@n"), 1, ch, nil, nil, TO_ROOM)
					phase_powerup(ch, 0, ch.Starphase)
				}
			case 1:
				if ch.Starphase != 1 {
					act(libc.CString("@WYou suddenly feel a @RSURGE@W of power through your body. You feel the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W come into its @CBirth Phase@W and its power is flowing into your body! Finally your eyes and the glyphs on your skin begin to glow an electric @bb@Bl@Cue@W!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@c$n@W suddenly seems to grow stronger for some reason. You notice $s eyes begin to glow an electric @bb@Bl@Cue@W. Suddenly glyphs start to appear all over $s skin and glow with the same light!@n"), 1, ch, nil, nil, TO_ROOM)
					phase_powerup(ch, 0, ch.Starphase)
					phase_powerup(ch, 1, 1)
				}
			case 2:
				if ch.Starphase != 2 {
					act(libc.CString("@WYou suddenly feel a @RSURGE@W of power through your body. You feel the @GE@gl@Dd@wri@Dt@gc@Gh @YS@yta@Yr@W come into its @GLife Phase@W and its power is flowing into your body! Finally your eyes and the glyphs on your skin begin to glow an fiery @Rr@re@Rd@W!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@c$n@W suddenly seems to grow stronger for some reason. You notice $s eyes begin to glow a fiery @rR@Re@rd@W. Suddenly glyphs start to appear all over $s skin and glow with the same light!@n"), 1, ch, nil, nil, TO_ROOM)
					phase_powerup(ch, 0, ch.Starphase)
					phase_powerup(ch, 1, 2)
				}
			default:
				send_to_imm(libc.CString("Strange Error in star_phase by: %s"), GET_NAME(ch))
			}
			return
		} else {
			return
		}
	}
	return
}
func phase_powerup(ch *char_data, type_ int, phase int) {
	if ch == nil {
		return
	}
	if IS_NPC(ch) {
		return
	}
	var change int = 0
	var bonus int = 0
	var mult float64 = 0.0
	switch phase {
	case 0:
		return
	case 1:
		change = 2
		mult = 4.0
		bonus = 5
	case 2:
		change = 3
		mult = 8.0
		bonus = 8
	default:
		send_to_imm(libc.CString("Error: phase_powerup called with GET_PHASE equal to zero by: %s"), GET_NAME(ch))
		return
	}
	if ETHER_STREAM(ch) {
		mult += 0.5
	}
	if type_ == 0 {
		ch.Hit = int64((float64(ch.Hit) - (float64(ch.Basepl)*0.1)*mult) / float64(change))
		ch.Mana = int64((float64(ch.Mana) - (float64(ch.Baseki)*0.1)*mult) / float64(change))
		ch.Move = int64((float64(ch.Move) - (float64(ch.Basest)*0.1)*mult) / float64(change))
		if ch.Hit < 0 {
			ch.Hit = 1
		}
		if ch.Mana < 0 {
			ch.Mana = 1
		}
		if ch.Move < 0 {
			ch.Move = 1
		}
		ch.Max_hit = ch.Basepl
		ch.Max_mana = ch.Baseki
		ch.Max_move = ch.Basest
		if (ch.Bonuses[BONUS_WIMP]) > 0 && int(ch.Aff_abils.Str) < 25 {
			ch.Real_abils.Str -= int8(bonus)
		}
		if (ch.Bonuses[BONUS_SLOW]) > 0 && int(ch.Aff_abils.Cha) < 25 {
			ch.Real_abils.Cha -= int8(bonus)
		}
		ch.Starphase = 0
	} else {
		ch.Hit = int64((float64(ch.Hit) + (float64(ch.Basepl)*0.1)*mult) * float64(change))
		ch.Mana = int64((float64(ch.Mana) + (float64(ch.Baseki)*0.1)*mult) * float64(change))
		ch.Move = int64((float64(ch.Move) + (float64(ch.Basest)*0.1)*mult) * float64(change))
		ch.Max_hit = int64((float64(ch.Basepl) + (float64(ch.Basepl)*0.1)*mult) * float64(change))
		ch.Max_mana = int64((float64(ch.Baseki) + (float64(ch.Baseki)*0.1)*mult) * float64(change))
		ch.Max_move = int64((float64(ch.Basest) + (float64(ch.Basest)*0.1)*mult) * float64(change))
		if ch.Hit > ch.Max_hit {
			ch.Hit = ch.Max_hit
		}
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		if (ch.Bonuses[BONUS_WIMP]) > 0 && int(ch.Aff_abils.Str)+bonus <= 25 {
			ch.Real_abils.Str += int8(bonus)
		}
		if (ch.Bonuses[BONUS_SLOW]) > 0 && int(ch.Aff_abils.Cha)+bonus <= 25 {
			ch.Real_abils.Cha += int8(bonus)
		}
		ch.Starphase = phase
	}
	save_char(ch)
	return
}
