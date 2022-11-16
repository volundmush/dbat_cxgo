package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func is_instrument(obj *obj_data) bool {
	return GET_OBJ_VNUM(obj) == 8802 || GET_OBJ_VNUM(obj) == 8807
}
func do_spiritcontrol(ch *char_data, argument *byte, cmd int, subcmd int) {
	if GET_SKILL(ch, SKILL_SPIRITCONTROL) == 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique.\r\n"))
		return
	} else {
		if AFF_FLAGGED(ch, AFF_SPIRITCONTROL) {
			send_to_char(ch, libc.CString("You have already concentrated and have full control of your spirit.\r\n"))
			return
		} else {
			var cost int64 = int64(float64(ch.Max_mana) * 0.2)
			if ch.Move < cost {
				send_to_char(ch, libc.CString("You need at least 20%s of your max ki in stamina to prepare this skill.\r\n"), "%")
				return
			} else {
				ch.Move -= cost
				act(libc.CString("@YYou concentrate and quantify every last bit of your spiritual and mental energies. You have full control of them and can bring them forth in an instant.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@y$n@Y seems to concentrate hard for a moment.@n"), TRUE, ch, nil, nil, TO_ROOM)
				var duration int = rand_number(2, 4)
				assign_affect(ch, AFF_SPIRITCONTROL, SKILL_SPIRITCONTROL, duration, 0, 0, 0, 0, 0, 0)
			}
		}
	}
}
func do_tailhide(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if int(ch.Race) != RACE_SAIYAN && int(ch.Race) != RACE_HALFBREED {
		send_to_char(ch, libc.CString("You have no need to hide your tail!\r\n"))
	}
	if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && !PLR_FLAGGED(ch, PLR_TAILHIDE) {
		SET_BIT_AR(ch.Act[:], PLR_TAILHIDE)
		send_to_char(ch, libc.CString("You have decided to hide your tail!\r\n"))
	} else if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && PLR_FLAGGED(ch, PLR_TAILHIDE) {
		REMOVE_BIT_AR(ch.Act[:], PLR_TAILHIDE)
		send_to_char(ch, libc.CString("You have decided to display your tail for all to see!\r\n"))
	}
}
func do_nogrow(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if int(ch.Race) != RACE_SAIYAN && int(ch.Race) != RACE_HALFBREED {
		send_to_char(ch, libc.CString("What do you mean?\r\n"))
	}
	if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && !PLR_FLAGGED(ch, PLR_NOGROW) {
		SET_BIT_AR(ch.Act[:], PLR_NOGROW)
		send_to_char(ch, libc.CString("You have decided to halt your tail growth!\r\n"))
	} else if (int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED) && PLR_FLAGGED(ch, PLR_NOGROW) {
		REMOVE_BIT_AR(ch.Act[:], PLR_NOGROW)
		send_to_char(ch, libc.CString("You have decided to regrow your tail!\r\n"))
	}
}
func do_restring(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		obj *obj_data
		pay int = 0
	)
	one_argument(argument, &arg[0])
	if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 178 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 184 {
		pay = 5000
		if ch.Gold < pay {
			send_to_char(ch, libc.CString("You need at least 5,000 zenni to initiate an equipment restring.\r\n"))
			return
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have a that equipment to restring in your inventory.\r\n"))
			send_to_char(ch, libc.CString("Syntax: restring (obj name)\r\n"))
			return
		} else if OBJ_FLAGGED(obj, ITEM_CUSTOM) {
			send_to_char(ch, libc.CString("You can not restring a custom piece. Why? Because you already restrung it you dummy.\r\n"))
			return
		} else {
			ch.Desc.Connected = CON_POBJ
			var thename [2048]byte
			var theshort [2048]byte
			var thelong [2048]byte
			thename[0] = '\x00'
			theshort[0] = '\x00'
			thelong[0] = '\x00'
			stdio.Sprintf(&thename[0], "%s", obj.Name)
			stdio.Sprintf(&theshort[0], "%s", obj.Short_description)
			stdio.Sprintf(&thelong[0], "%s", obj.Description)
			ch.Desc.Obj_name = libc.StrDup(&thename[0])
			ch.Desc.Obj_was = libc.StrDup(&theshort[0])
			ch.Desc.Obj_short = libc.StrDup(&theshort[0])
			ch.Desc.Obj_long = libc.StrDup(&thelong[0])
			ch.Desc.Obj_point = obj
			ch.Desc.Obj_type = 1
			ch.Desc.Obj_weapon = 0
			disp_restring_menu(ch.Desc)
			ch.Desc.Obj_editflag = EDIT_RESTRING
			ch.Desc.Obj_editval = EDIT_RESTRING_MAIN
			return
		}
	}
}
func do_multiform(ch *char_data, argument *byte, cmd int, subcmd int) {
	if !IS_NPC(ch) && GET_SKILL(ch, SKILL_MULTIFORM) == 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique.\r\n"))
		return
	}
	var tch *char_data = nil
	var next_v *char_data = nil
	var multi1 *char_data = nil
	var multi2 *char_data = nil
	var multi3 *char_data = nil
	var num int = 0
	for tch = world[ch.In_room].People; tch != nil; tch = next_v {
		next_v = tch.Next_in_room
		if tch == ch {
			continue
		}
		if !IS_NPC(tch) {
			continue
		}
		if GET_MOB_VNUM(tch) == 25 {
			if tch.Original == ch {
				num += 1
				if multi1 == nil {
					multi1 = tch
				} else if multi2 == nil {
					multi2 = tch
				} else if multi3 == nil {
					multi3 = tch
				} else {
					basic_mud_log(libc.CString("Error: More clones processed than allowed.\r\n"))
					return
				}
			}
		}
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if num >= 3 {
		if libc.StrCaseCmp(&arg[0], libc.CString("merge")) == 0 {
			extract_char(multi1)
			extract_char(multi2)
			extract_char(multi3)
			return
		} else {
			send_to_char(ch, libc.CString("You have the maximum number of bodies multiform can produce.\nTo merge use: multiform merge.\r\n"))
			return
		}
	} else {
		if libc.StrCaseCmp(&arg[0], libc.CString("merge")) == 0 {
			if multi1 != nil {
				extract_char(multi1)
			}
			if multi2 != nil {
				extract_char(multi2)
			}
			if multi3 != nil {
				extract_char(multi3)
			}
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("split")) == 0 {
			var (
				cost    int64 = int64((float64(ch.Max_mana) * 0.005) + float64(ch.Max_move)*0.005 + 2)
				penalty int   = 0
			)
			if ch.Fighting != nil {
				penalty = rand_number(8, 15)
			}
			var roll int = axion_dice(penalty)
			cost *= int64(float64(GET_SKILL(ch, SKILL_MULTIFORM)) * 0.2)
			if ch.Mana < cost {
				send_to_char(ch, libc.CString("You do not have enough ki to split!\r\n"))
				return
			} else if ch.Move < cost {
				send_to_char(ch, libc.CString("You do not have enough stamina to split!\r\n"))
				return
			}
			if float64(ch.Hit) < float64(gear_pl(ch))*0.4 {
				send_to_char(ch, libc.CString("Your powerlevel is too weakened to split! It can't be below 40%s of your weighted max.\r\n"), "%")
				return
			} else if GET_SKILL(ch, SKILL_MULTIFORM) < roll {
				act(libc.CString("@YYou focus your ki into your body while concentrating on the image of your body splitting into two. @yYou lose your concentration and fail to split though...@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@y$n@Y seems to concentrate really hard for a moment, before relaxing.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Mana -= cost
				ch.Move -= cost
				improve_skill(ch, SKILL_MULTIFORM, 1)
			} else {
				act(libc.CString("@YYou focus your ki into your body while concentrating on the image of your body splitting into two. Another you splits out of your body!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@YSuddenly @y$n@Y seems to concentrates really and after a brief moment splits into two copies of $mself!@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Mana -= cost
				ch.Move -= cost
				if multi1 == nil {
					generate_multiform(ch, multi1, multi2, multi3)
				} else if multi2 == nil {
					generate_multiform(ch, multi1, multi2, multi3)
				} else if multi3 == nil {
					generate_multiform(ch, multi1, multi2, multi3)
				}
				improve_skill(ch, SKILL_MULTIFORM, 1)
			}
		}
	}
}
func generate_multiform(ch *char_data, multi1 *char_data, multi2 *char_data, multi3 *char_data) {
	var (
		clone *char_data = nil
		r_num mob_rnum
	)
	if (func() mob_rnum {
		r_num = real_mobile(25)
		return r_num
	}()) == mob_rnum(-1) {
		send_to_imm(libc.CString("Multiform Clone doesn't exist!"))
		return
	}
	clone = read_mobile(mob_vnum(r_num), REAL)
	char_to_room(clone, ch.In_room)
	var buf [2048]byte
	clone.Name = (*byte)(unsafe.Pointer(uintptr('\x00')))
	buf[0] = '\x00'
	stdio.Snprintf(&buf[0], int(2048), "%s's Clone", GET_NAME(ch))
	clone.Name = libc.StrDup(&buf[0])
	clone.Short_descr = (*byte)(unsafe.Pointer(uintptr('\x00')))
	buf[0] = '\x00'
	stdio.Snprintf(&buf[0], int(2048), "%s's @CClone@n", GET_NAME(ch))
	clone.Short_descr = libc.StrDup(&buf[0])
	clone.Long_descr = (*byte)(unsafe.Pointer(uintptr('\x00')))
	buf[0] = '\x00'
	stdio.Snprintf(&buf[0], int(2048), "%s's @CClone@w is standing here.@n\n", GET_NAME(ch))
	clone.Long_descr = libc.StrDup(&buf[0])
	clone.Description = (*byte)(unsafe.Pointer(uintptr('\x00')))
	buf[0] = '\x00'
	stdio.Snprintf(&buf[0], int(2048), "%s", ch.Description)
	clone.Description = libc.StrDup(&buf[0])
	clone.Race = ch.Race
	clone.Chclass = ch.Chclass
	var multi_forms int = 0
	_ = multi_forms
	if multi1 != nil {
		multi_forms += 1
	}
	if multi2 != nil {
		multi_forms += 1
	}
	if multi3 != nil {
		multi_forms += 1
	}
	ch.Clones += 1
	var mult float64 = 1
	if int(ch.Race) == RACE_BIO {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 3.5
		} else if PLR_FLAGGED(ch, PLR_TRANS4) {
			mult = 4
		}
	} else if int(ch.Race) == RACE_MAJIN {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 4.5
		}
	} else if int(ch.Race) == RACE_TRUFFLE {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 4
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 5
		}
	} else if PLR_FLAGGED(ch, PLR_TRANS1) || PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) || PLR_FLAGGED(ch, PLR_TRANS5) || PLR_FLAGGED(ch, PLR_TRANS6) {
		do_transform(ch, libc.CString("revert"), 0, 0)
	}
	ch.Max_hit -= int64((float64(ch.Basepl) * 0.25) * mult)
	ch.Max_move -= int64((float64(ch.Basest) * 0.25) * mult)
	ch.Max_mana -= int64((float64(ch.Baseki) * 0.25) * mult)
	var blamo [2048]byte
	stdio.Sprintf(&blamo[0], "p.%s", GET_NAME(ch))
	do_follow(clone, &blamo[0], 0, 0)
	if ch.Hit > gear_pl(ch) {
		ch.Hit = gear_pl(ch)
	}
	if ch.Mana > ch.Max_mana {
		ch.Mana = ch.Max_mana
	}
	if ch.Move > ch.Max_move {
		ch.Move = ch.Max_move
	}
	if multi1 != nil {
		multi1.Hit = ch.Hit
		multi1.Max_hit = ch.Max_hit
		multi1.Move = ch.Move
		multi1.Max_move = ch.Max_move
		multi1.Mana = ch.Mana
		multi1.Max_mana = ch.Max_mana
	}
	if multi2 != nil {
		multi2.Hit = ch.Hit
		multi2.Max_hit = ch.Max_hit
		multi2.Move = ch.Move
		multi2.Max_move = ch.Max_move
		multi2.Mana = ch.Mana
		multi2.Max_mana = ch.Max_mana
	}
	if multi3 != nil {
		multi3.Hit = ch.Hit
		multi3.Max_hit = ch.Max_hit
		multi3.Move = ch.Move
		multi3.Max_move = ch.Max_move
		multi3.Mana = ch.Mana
		multi3.Max_mana = ch.Max_mana
	}
	clone.Hit = ch.Hit
	clone.Max_hit = ch.Max_hit
	clone.Move = ch.Move
	clone.Max_move = ch.Max_move
	clone.Mana = ch.Mana
	clone.Max_mana = ch.Max_mana
	clone.Original = ch
}
func handle_multi_merge(form *char_data) {
	var ch *char_data = form.Original
	if ch == nil {
		return
	}
	send_to_char(ch, libc.CString("@YYou merge with one of your forms!@n\r\n"))
	act(libc.CString("@y$n@Y merges with one of his multiforms!@n\r\n"), TRUE, ch, nil, nil, TO_ROOM)
	var mult float64 = 1
	if int(ch.Race) == RACE_BIO {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 3.5
		} else if PLR_FLAGGED(ch, PLR_TRANS4) {
			mult = 4
		}
	} else if int(ch.Race) == RACE_MAJIN {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 4.5
		}
	} else if int(ch.Race) == RACE_TRUFFLE {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 4
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 5
		}
	}
	ch.Clones -= 1
	ch.Max_hit += int64((float64(ch.Basepl) * 0.25) * mult)
	ch.Max_mana += int64((float64(ch.Baseki) * 0.25) * mult)
	ch.Max_move += int64((float64(ch.Basest) * 0.25) * mult)
	ch.Hit += int64((float64(ch.Basepl) * 0.25) * mult)
	ch.Mana += int64((float64(ch.Baseki) * 0.25) * mult)
	ch.Move += int64((float64(ch.Basest) * 0.25) * mult)
	if ch.Hit > gear_pl(ch) {
		ch.Hit = gear_pl(ch)
	}
	if ch.Mana > ch.Max_mana {
		ch.Mana = ch.Max_mana
	}
	if ch.Move > ch.Max_move {
		ch.Max_move = ch.Max_move
	}
	extract_char(form)
}
func handle_songs() {
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if !IS_PLAYING(d) {
			continue
		}
		if d.Character == nil {
			continue
		}
		if d.Character.Powerattack > 0 {
			resolve_song(d.Character)
		}
	}
}
func resolve_song(ch *char_data) {
	var (
		vict     *char_data = nil
		next_v   *char_data = nil
		obj2     *obj_data  = nil
		next_obj *obj_data
	)
	_ = next_obj
	var diceroll int = axion_dice(0)
	var skill int = GET_SKILL(ch, SKILL_MYSTICMUSIC)
	var stopplaying int = FALSE
	_ = stopplaying
	var buf [2048]byte
	if ch.Powerattack <= 0 {
		return
	}
	var instrument vnum = 0
	_ = instrument
	obj2 = find_obj_in_list_lambda(ch.Carrying, is_instrument)
	if obj2 == nil {
		send_to_char(ch, libc.CString("You do not have an instrument.\r\n"))
		act(libc.CString("@c$n@C stops playing $s song.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Powerattack = 0
		return
	}
	instrument = vnum(GET_OBJ_VNUM(obj2))
	if skill > diceroll {
		stdio.Sprintf(&buf[0], "@c$n@C continues playing @y'@Y%s@y'@C.@n", func() string {
			if ch.Powerattack == SONG_SAFETY {
				return "Song of Safety"
			}
			if ch.Powerattack == SONG_SHIELDING {
				return "Song of Shielding"
			}
			if ch.Powerattack == SONG_SHADOW_STITCH {
				return "Shadow Stitch Minuet"
			}
			return "Teleportation Melody"
		}())
		act(libc.CString("@CYou continue playing your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
	} else {
		act(libc.CString("@CYou mess up a portion of the song, but continue playing.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C messes up a portion of $s song, but continues to play.@n"), TRUE, ch, nil, nil, TO_ROOM)
		return
	}
	for vict = world[ch.In_room].People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		switch ch.Powerattack {
		case SONG_SAFETY:
			if ch.Master == vict.Master || ch == vict.Master || vict == ch.Master || vict == ch {
				if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) || vict == ch {
					if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master || vict == ch {
						if skill > diceroll {
							var restore int64 = int64(float64(skill*10) + (float64(ch.Max_mana)*0.0004)*float64(skill))
							if vict != ch {
								act(libc.CString("@CYour skillfully playing of the Song of Safety has an effect on @c$N@C.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							} else {
								act(libc.CString("@CYour skillfully playing of the Song of Safety has an effect on your own body@C.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							}
							vict.Hit += restore
							if vict.Hit > gear_pl(vict) {
								vict.Hit = gear_pl(vict)
							}
							vict.Move += int64(float64(restore) * 0.5)
							if vict.Move > vict.Max_move {
								vict.Move = vict.Max_move
							}
							if (vict.Limb_condition[0]) < 100 {
								vict.Limb_condition[0] += int((float64(skill) * 0.1) + 1)
								if (vict.Limb_condition[0]) > 100 {
									send_to_char(vict, libc.CString("Your right arm is no longer broken!@n\r\n"))
									vict.Limb_condition[0] = 100
								}
							}
							if (vict.Limb_condition[1]) < 100 {
								vict.Limb_condition[1] += int((float64(skill) * 0.1) + 1)
								if (vict.Limb_condition[1]) > 100 {
									send_to_char(vict, libc.CString("Your left arm is no longer broken!@n\r\n"))
									vict.Limb_condition[1] = 100
								}
							}
							if (vict.Limb_condition[2]) < 100 {
								vict.Limb_condition[2] += int((float64(skill) * 0.1) + 1)
								if (vict.Limb_condition[2]) > 100 {
									send_to_char(vict, libc.CString("Your right leg is no longer broken!@n\r\n"))
									vict.Limb_condition[2] = 100
								}
							}
							if (vict.Limb_condition[0]) < 100 {
								vict.Limb_condition[3] += int((float64(skill) * 0.1) + 1)
								if (vict.Limb_condition[3]) > 100 {
									send_to_char(vict, libc.CString("Your left leg is no longer broken!@n\r\n"))
									vict.Limb_condition[3] = 100
								}
							}
							if vict != ch {
								act(libc.CString("@c$n's@C soothing Song of Safety has recovered some of your powerlevel, stamina, and limb condition."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
							}
							act(libc.CString("@c$n@C continues playing $s ocarina!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
							improve_skill(ch, SKILL_MYSTICMUSIC, 2)
							ch.Mana -= int64((float64(ch.Max_mana) * 0.0003) + float64(skill))
						}
					}
				}
			}
			if ch.Mana <= 0 {
				send_to_char(ch, libc.CString("You no longer have the ki necessary to play your song.\r\n"))
				act(libc.CString("@c$n@C stops playing $s song.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = 0
				return
			}
		case SONG_SHADOW_STITCH:
			if ch.Master != nil && vict.Master != nil {
				if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
					if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
						continue
					} else if skill > diceroll+10 {
						act(libc.CString("@CYour forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C forboding music has caused YOUR shadows to stitch into YOUR body, slow YOUR actions down!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						if !IS_NPC(vict) {
							WAIT_STATE(vict, (int(1000000/OPT_USEC))*2)
							ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
						} else {
							vict.Real_abils.Cha -= 2
							ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
							if int(vict.Real_abils.Cha) < 3 {
								vict.Real_abils.Cha = 3
							}
						}
					}
				} else if skill > diceroll+10 {
					act(libc.CString("@CYour forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n's@C forboding music has caused YOUR shadows to stitch into YOUR body, slow YOUR actions down!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n's@C forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if !IS_NPC(vict) {
						WAIT_STATE(vict, (int(1000000/OPT_USEC))*2)
						ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
					} else {
						ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
						vict.Real_abils.Cha -= 2
						if int(vict.Real_abils.Cha) < 3 {
							vict.Real_abils.Cha = 3
						}
					}
				}
			} else if skill > diceroll+10 {
				act(libc.CString("@CYour forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n's@C forboding music has caused YOUR shadows to stitch into YOUR body, slow YOUR actions down!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@C forboding music has caused @c$N's@C shadows to stitch into $S body, slowing $S actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !IS_NPC(vict) {
					WAIT_STATE(vict, (int(1000000/OPT_USEC))*2)
					ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
				} else {
					vict.Real_abils.Cha -= 2
					ch.Mana -= int64((float64(ch.Max_mana) * 0.001) + float64(skill))
					if int(vict.Real_abils.Cha) < 3 {
						vict.Real_abils.Cha = 3
					}
				}
			}
			if ch.Mana <= 0 {
				send_to_char(ch, libc.CString("You no longer have the ki necessary to play your song.\r\n"))
				act(libc.CString("@c$n@C stops playing $s song.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = 0
				return
			}
		case SONG_TELEPORT_EARTH:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Earth in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Earth in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(300))
					}
				}
			}
		case SONG_TELEPORT_VEGETA:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Vegeta in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Vegeta in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(2234))
					}
				}
			}
		case SONG_TELEPORT_FRIGID:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Frigid in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Frigid in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(4047))
					}
				}
			}
		case SONG_TELEPORT_KONACK:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Konack in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Konack in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(8003))
					}
				}
			}
		case SONG_TELEPORT_NAMEK:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Namek in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Namek in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(0x27C6))
					}
				}
			}
		case SONG_TELEPORT_ARLIA:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Arlia in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Arlia in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(0x3ED7))
					}
				}
			}
		case SONG_TELEPORT_AETHER:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Aether in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Aether in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(0x2EF9))
					}
				}
			}
		case SONG_TELEPORT_KANASSA:
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master {
					if skill > diceroll {
						act(libc.CString("@CYour Teleportation Melody has transported @c$N@C to Kanassa in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@c$n's@C Teleportation Melody has transported you to Kanassa in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n's@C Teleportation Melody has transported @c$N@C away in a flash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						char_from_room(vict)
						char_to_room(vict, real_room(14910))
					}
				}
			}
		case SONG_SHIELDING:
			if vict == ch || AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
				if ch == vict.Master || ch.Master == vict || ch.Master == vict.Master || vict == ch {
					if skill > diceroll {
						if vict != ch {
							act(libc.CString("@CYour triumphant and soaring music has powered a barrier around @c$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("@c$n's@C triumphant and soaring music has powered a barrier around you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						} else {
							act(libc.CString("@CYour triumphant and soaring music has powered a barrier around yourself@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						}
						act(libc.CString("@c$n's@C triumphant and soaring music has powered a barrier around @c$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Barrier += int64(((float64(ch.Max_mana) * 0.005) * (float64(skill) * 0.25)) + float64(skill))
						if float64(vict.Barrier) >= float64(vict.Max_mana)*0.75 {
							vict.Barrier = int64(float64(vict.Max_mana) * 0.75)
						}
						if !AFF_FLAGGED(vict, AFF_SANCTUARY) {
							SET_BIT_AR(vict.Affected_by[:], AFF_SANCTUARY)
						}
						ch.Mana -= int64((float64(ch.Max_mana) * 0.02) + float64(skill))
					}
				}
			}
			if ch.Mana <= 0 {
				send_to_char(ch, libc.CString("You no longer have the ki necessary to play your song.\r\n"))
				act(libc.CString("@c$n@C stops playing $s song.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = 0
				return
			}
		}
	}
	if ch.Powerattack >= 4 && skill > diceroll {
		switch ch.Powerattack {
		case SONG_TELEPORT_EARTH:
			char_from_room(ch)
			char_to_room(ch, real_room(300))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Earth and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_VEGETA:
			char_from_room(ch)
			char_to_room(ch, real_room(2234))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Vegeta and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_FRIGID:
			char_from_room(ch)
			char_to_room(ch, real_room(4047))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Frigid and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_NAMEK:
			char_from_room(ch)
			char_to_room(ch, real_room(0x27C6))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Namek and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_KANASSA:
			char_from_room(ch)
			char_to_room(ch, real_room(14910))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Kanassa and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_AETHER:
			char_from_room(ch)
			char_to_room(ch, real_room(0x2EF9))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Aether and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_ARLIA:
			char_from_room(ch)
			char_to_room(ch, real_room(0x3ED7))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Arlia and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		case SONG_TELEPORT_KONACK:
			char_from_room(ch)
			char_to_room(ch, real_room(8003))
			ch.Powerattack = 0
			act(libc.CString("@CFinally as the last of your comrades has been teleported you teleport yourself to Konack and stop your song.@n"), TRUE, ch, nil, nil, TO_CHAR)
		}
	}
}
func do_song(ch *char_data, argument *byte, cmd int, subcmd int) {
	if GET_SKILL(ch, SKILL_MYSTICMUSIC) == 0 {
		send_to_char(ch, libc.CString("You do not know how to play mystical music.\r\n"))
		return
	}
	var obj2 *obj_data = find_obj_in_list_lambda(ch.Carrying, is_instrument)
	if obj2 == nil {
		send_to_char(ch, libc.CString("You do not have an instrument.\r\n"))
		return
	}
	var instrument vnum = vnum(GET_OBJ_VNUM(obj2))
	if ch.Powerattack != 0 {
		act(libc.CString("@cYou stop playing your ocarina.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n stops playing their ocarina.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Powerattack = 0
		return
	} else {
		var (
			arg   [2048]byte
			arg2  [2048]byte
			skill int = GET_SKILL(ch, SKILL_MYSTICMUSIC)
		)
		two_arguments(argument, &arg[0], &arg2[0])
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("@YSongs Known\n@c-------------------@n\r\n"))
			send_to_char(ch, libc.CString("@W%s%s%sSong of Safety\r\n"), func() string {
				if skill > 99 {
					return "Song of Shielding\n"
				}
				return ""
			}(), func() string {
				if skill > 80 {
					return "Melody of Teleportation\n"
				}
				return ""
			}(), func() string {
				if skill > 50 {
					return "Shadow Stitch Minuet\n"
				}
				return ""
			}())
			send_to_char(ch, libc.CString("@wSyntax: song (shielding | safety | teleport | shadow )\r\n"))
			return
		} else {
			var (
				cost     int64 = int64(float64(ch.Max_mana) * 0.01)
				modifier int   = 1
			)
			if libc.StrCaseCmp(&arg[0], libc.CString("shielding")) == 0 {
				modifier = 20
				cost *= int64(modifier)
			} else if libc.StrCaseCmp(&arg[0], libc.CString("teleport")) == 0 {
				modifier = 50
				cost *= int64(modifier)
			} else if libc.StrCaseCmp(&arg[0], libc.CString("shadow")) == 0 {
				modifier = 8
				cost *= int64(modifier)
			} else if libc.StrCaseCmp(&arg[0], libc.CString("safety")) == 0 {
				modifier = 3
				cost *= int64(modifier)
			}
			if instrument == 8802 {
				cost -= int64(float64(cost) * 0.5)
			}
			if modifier == 0 {
				send_to_char(ch, libc.CString("@wSyntax: song (shielding | safety | teleport | shadow )\r\n"))
				return
			} else if modifier == 3 && ch.Mana < cost {
				send_to_char(ch, libc.CString("@wYou do not have enough ki to power the instrument for that song!@n\r\n"))
				return
			} else if modifier == 3 {
				act(libc.CString("@CYou begin to play the Song of Safety! Your fingers lightly glide over the ocarina and as you blow into it sweet music similar to a lullaby issues forth from the intrument.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C begins to play a song on $s ocarina. The music seems to be some sort of lullaby.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = SONG_SAFETY
				ch.Mana -= cost
				return
			} else if modifier == 8 && ch.Mana < cost {
				send_to_char(ch, libc.CString("@wYou do not have enough ki to power the instrument for that song!@n\r\n"))
				return
			} else if modifier == 8 && skill <= 49 {
				send_to_char(ch, libc.CString("You do not posess the skill to play such a song!\r\n"))
				return
			} else if modifier == 8 {
				act(libc.CString("@CYou begin to play the Shadow Stitch Minuet! Your fingers lightly glide over the ocarina and as you blow into it forboding low toned music issues forth.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C begins to play a song on $s ocarina. Depressing low toned music issues forth from the ocarina.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = SONG_SHADOW_STITCH
				ch.Mana -= cost
				return
			} else if modifier == 50 && ch.Mana < cost {
				send_to_char(ch, libc.CString("@wYou do not have enough ki to power the instrument for that song!@n\r\n"))
				return
			} else if modifier == 50 && skill <= 79 {
				send_to_char(ch, libc.CString("You do not posess the skill to play such a song!\r\n"))
				return
			} else if modifier == 50 {
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("Where would you like to teleport to?\nSyntax: song teleport (earth | vegeta | kanassa | arlia | aether | namek | konack| frigid)\r\n"))
					return
				} else if AFF_FLAGGED(ch, AFF_SPIRIT) {
					send_to_char(ch, libc.CString("Not while you're dead!\r\n"))
					return
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("earth")) == 0 {
					ch.Powerattack = SONG_TELEPORT_EARTH
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("frigid")) == 0 {
					ch.Powerattack = SONG_TELEPORT_FRIGID
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("vegeta")) == 0 {
					ch.Powerattack = SONG_TELEPORT_VEGETA
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("namek")) == 0 {
					ch.Powerattack = SONG_TELEPORT_NAMEK
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("arlia")) == 0 {
					ch.Powerattack = SONG_TELEPORT_ARLIA
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("kanassa")) == 0 {
					ch.Powerattack = SONG_TELEPORT_KANASSA
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("konack")) == 0 {
					ch.Powerattack = SONG_TELEPORT_KONACK
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("aether")) == 0 {
					ch.Powerattack = SONG_TELEPORT_AETHER
				} else {
					send_to_char(ch, libc.CString("Syntax: song teleport (earth | vegeta | namek | aether | konack | kanassa | arlia | frigid)\r\n"))
					return
				}
				act(libc.CString("@CYou begin to play the Melody of Teleportation! Your fingers lightly glide over the ocarina and as you blow into it a repeating light hearted melody issues forth.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C begins to play a song on $s ocarina. A light hearted melody can be heard sounding from the ocarina as it is played.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Mana -= cost
				return
			} else if modifier == 20 && ch.Mana < cost {
				send_to_char(ch, libc.CString("@wYou do not have enough ki to power the instrument for that song!@n\r\n"))
				return
			} else if modifier == 20 && skill <= 98 {
				send_to_char(ch, libc.CString("You do not posess the skill to play such a song!\r\n"))
				return
			} else if modifier == 20 {
				act(libc.CString("@CYou begin to play the Song of Shielding! Your fingers lightly glide over the ocarina and as you blow into it a triumphant series of notes issues forth.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C begins to play a song on $s ocarina. A triumphant song full of soaring sounds from the ocarina as it is played.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Powerattack = SONG_SHIELDING
				ch.Mana -= cost
				return
			}
		}
	}
}
func do_preference(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: preference (throw | weapon | hand | ki)\r\n"))
		return
	}
	if ch.Preference > 0 {
		send_to_char(ch, libc.CString("You've already chosen a specialization. No going back.\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("throw")) == 0 {
		send_to_char(ch, libc.CString("You will now favor throwing weapons as fighting specialization. You're sure to nail it.\r\n"))
		ch.Preference = PREFERENCE_THROWING
		if int(ch.Skills[SKILL_THROW]) <= 90 {
			ch.Skills[SKILL_THROW] += 10
		} else if int(ch.Skills[SKILL_THROW]) < 100 {
			ch.Skills[SKILL_THROW] = 100
		}
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("hand")) == 0 {
		send_to_char(ch, libc.CString("You will now favor your body as your fighting specialization. Your body is ready.\r\n"))
		ch.Preference = PREFERENCE_H2H
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("ki")) == 0 {
		send_to_char(ch, libc.CString("You will now favor your ki energy as your fighting specialization. I expect more than a few smoldering craters.\r\n"))
		ch.Preference = PREFERENCE_KI
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("weapon")) == 0 {
		send_to_char(ch, libc.CString("You will now favor your weapons as your fighting specialization. Let the blood fly!\r\n"))
		ch.Preference = PREFERENCE_WEAPON
		return
	} else {
		send_to_char(ch, libc.CString("Syntax: preference (throw | weapon | hand | ki)\r\n"))
		return
	}
}
func do_moondust(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		cost int64 = int64(float64(ch.Max_move) * 0.02)
		heal int64 = 0
	)
	if int(ch.Race) != RACE_ARLIAN || int(ch.Sex) != SEX_FEMALE {
		send_to_char(ch, libc.CString("You are not an arlian female.\r\n"))
		return
	}
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("You need to be in a group to use this skill!\r\n"))
		return
	}
	cost += int64(float64(GET_LIFEMAX(ch)) * 0.02)
	heal = cost * 3
	if float64(ch.Hit) >= float64(gear_pl(ch))*0.8 {
		cost = int64(float64(cost) * 0.5)
	}
	if ch.Move < cost {
		send_to_char(ch, libc.CString("You do not have enough stamina to perform this technique.\r\n"))
		return
	}
	var chance int = axion_dice(0)
	if chance > int(ch.Aff_abils.Wis)+rand_number(1, 10) {
		act(libc.CString("@GYou spread your wings and begin to concentrate. Your wings begin to glow a soft sea green color. As you prepare to release a cloud of your charged wing dust you lose focus and the power you had begun to charge into your wings dissipates.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n@G spreads $s wings and seems to concentrate for a moment. Suddenly $s wings begin to glow a soft sea green color. This soft glow grows brighter for a second before fading completely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Move -= cost
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		return
	}
	ch.Hit += heal
	if ch.Hit > gear_pl(ch) {
		ch.Hit = gear_pl(ch)
	}
	ch.Move -= cost
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
	act(libc.CString("@GYou spread your wings and begin to concentrate. Your wings begin to glow a soft sea green color. As your wings grow brighter you focus your charged bio energy in a shockwave the unleashes a cloud of glowing green dust. You breath in the dust and feel it rejuvinate your body's cells!@n"), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("@g$n@G spreads $s wings and seems to concentrate for a moment. Suddenly $s wings begin to glow a soft sea green color. This soft glow grows brighter and as $e flexes $s wings to their full extent a shockwave of energy explodes outward. Carried on this shockwave is a cloud of glowing dust! You notice some of the dust being breathed in by $s!@n"), TRUE, ch, nil, nil, TO_ROOM)
	send_to_char(ch, libc.CString("@RHeal@Y: @C%s@n\r\n"), add_commas(heal))
	var vict *char_data = nil
	var next_v *char_data = nil
	for vict = world[ch.In_room].People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		}
		if AFF_FLAGGED(vict, AFF_GROUP) {
			if ch.Master == vict.Master || vict.Master == ch || ch.Master == vict {
				vict.Hit += heal
				if vict.Hit > gear_pl(vict) {
					vict.Hit = gear_pl(vict)
				}
				act(libc.CString("@CYou breathe in the dust and are healed by it somewhat!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C breathes in the dust and is healed somewhat!@n"), TRUE, vict, nil, nil, TO_ROOM)
				send_to_char(vict, libc.CString("@RHeal@Y: @C%s@n\r\n"), add_commas(heal))
			}
		}
	}
}
func do_shell(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Race) != RACE_ARLIAN {
		send_to_char(ch, libc.CString("You are not capable of doing that!\r\n"))
		return
	}
	if int(ch.Sex) == SEX_FEMALE {
		send_to_char(ch, libc.CString("Sorry, you can't do that.\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SHELL) {
		act(libc.CString("@mYou quickly absorb the armor carapace covering your body back inside.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n's@m armored carapce retreats back to its original size.@n"), TRUE, ch, nil, nil, TO_ROOM)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_SHELL)
		return
	}
	if float64(ch.Move) < float64(ch.Max_move)*0.2 {
		send_to_char(ch, libc.CString("You do not have enough stamina to grow your armored carapace.@n\r\n"))
		return
	} else if axion_dice(0) > int(ch.Aff_abils.Con)+rand_number(1, 10) {
		act(libc.CString("@mYou crouch down and begin to focus on your body's carapace cells encouraging them to multiply! However your control is lacking and you ultimately fail to grow your armor very much.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n@m crouches down and seems to strain for a moment before giving up and resuming $s normal stance.@n"), TRUE, ch, nil, nil, TO_ROOM)
		return
	} else {
		act(libc.CString("@mYou crouch down and begin to focus on your body's carapace cells, encouraging them to multiply! Very quickly millions of new carapace cells have been born and your armored carapace extends over all parts of your body!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@M$n@m crouches down and after a few moments of straining $s body's carapace armor starts to grow thicker and extends to cover all parts of $s body!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Move -= int64(float64(ch.Max_move) * 0.2)
		SET_BIT_AR(ch.Affected_by[:], AFF_SHELL)
		return
	}
}
func do_liquefy(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Race) != RACE_MAJIN {
		send_to_char(ch, libc.CString("You are not capable of liquefying yourself right now. Try finding a giant blender maybe?\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
		act(libc.CString("@MSuddenly large chunks of goo start to hover up slowly. These very same chunks quickly begin to fly into each other, piling on as the ball of goo grows. Suddenly @m$n@M emerges as the ball of goo takes $s shape!@n"), TRUE, ch, nil, nil, TO_ROOM)
		act(libc.CString("@MYou begin to pull the liquid chunks of your body together. Those chunks hover upward and merge into each other until a large ball of goo is formed. Slowly your body emerges as the pieces of your body take on their old form!@n"), TRUE, ch, nil, nil, TO_CHAR)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_LIQUEFIED)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: liquefy hide\nSyntax: liquefy explode (target)\r\n"))
		return
	}
	if float64(ch.Mana) < (float64(ch.Max_mana)*0.002)+150 {
		send_to_char(ch, libc.CString("You do not have enough ki to manage this level of body control!\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("hide")) == 0 {
		if ch.Grappled != nil {
			ch.Grappled.Grappling = nil
			ch.Grappled = nil
		}
		if ch.Grappling != nil {
			ch.Grappling.Grappled = nil
			ch.Grappling = nil
		}
		if ch.Drag != nil {
			ch.Drag.Dragged = nil
			ch.Drag = nil
		}
		if axion_dice(0) > GET_LEVEL(ch) {
			act(libc.CString("@MYour body starts to become loose and sag, but you lose focus and it reverts to its original shape!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@m$n@M's body starts to become loose and sag, but $e seems to return normal a moment later.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			return
		} else {
			act(libc.CString("@MYour body starts to become loose and sag. It continues to droop down until it begins to run down like a river of goo flowing from where your body was.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@m$n@M's body starts to become loose and sag. Much of $s body begins to pour down and scatter around as pools of goo.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			SET_BIT_AR(ch.Affected_by[:], AFF_LIQUEFIED)
			return
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("explode")) == 0 {
		var vict *char_data
		if ch.Grappled != nil {
			ch.Grappled.Grappling = nil
			ch.Grappled = nil
		}
		if ch.Grappling != nil {
			ch.Grappling.Grappled = nil
			ch.Grappling = nil
		}
		if ch.Drag != nil {
			ch.Drag.Dragged = nil
			ch.Drag = nil
		}
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: liquefy hide\nSyntax: liquefy explode (target)\r\n"))
			return
		} else if float64(ch.Mana) < (float64(ch.Max_mana)*0.1)+150 {
			send_to_char(ch, libc.CString("You do not have enough ki for that action!@n\r\n"))
			return
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("That target isn't here.\r\n"))
			return
		} else if can_kill(ch, vict, nil, 1) == 0 {
			send_to_char(ch, libc.CString("You can't kill them!\r\n"))
			return
		} else if axion_dice(0) > GET_LEVEL(ch) {
			act(libc.CString("@MYour body starts to become loose and sag, but you lose focus and it reverts to its original shape!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@m$n@M's body starts to become loose and sag, but $e seems to return normal a moment later.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if GET_SPEEDI(ch) < GET_SPEEDI(vict) {
			act(libc.CString("@MYour body rapidly turns to liquid and flies for @R$N's@M open mouth! However $E easily dodges and avoids your attempt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@m$n@M's body rapidly turns to liquid and flies for @RYOUR@M open mouth! However you are faster and managed to dodge the attempt.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@m$n@M's body rapidly turns into liquid and flies for @R$N's@M open mouth! However $E easily dodges and avoids @m$n's@M attempt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ch.Hit < vict.Hit*2 {
			act(libc.CString("@MYour body rapidly turns to liquid and flies for @R$N's@M open mouth! However as you force yourself in through $S mouth $E successfully resists and forces your back out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@m$n@M's body rapidly turns to liquid and flies for @RYOUR@M open mouth! However you think quickly and force $m out before $e has a chance to get fully into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@m$n@M's body rapidly turns into liquid and flies for @R$N's@M open mouth! However as $e forces $mself in through @R$N's@M mouth $E manages to resist and force @m$n@M back out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			var dmg int64 = int64(float64(ch.Max_hit) * 0.08)
			hurt(0, 0, ch, vict, nil, dmg, 0)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else {
			act(libc.CString("@MYour body rapidly turns to liquid and flies for @R$N's@M open mouth! As you fill $S body you expand outward until $s body explodes into a gory mess!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@m$n@M's body rapidly turns to liquid and flies for @RYOUR@M open mouth! As $e fills your body it begins to expand until it is unable to take the strain any longer and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@m$n@M's body rapidly turns into liquid and flies for @R$N's@M open mouth! As $e forces $mself in through @R$N's@M mouth $S body begins to expand until it can't take the strain any longer and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Mana -= int64((float64(ch.Max_mana) * 0.002) + 150)
			if AFF_FLAGGED(ch, AFF_GROUP) {
				group_gain(ch, vict)
			} else {
				solo_gain(ch, vict)
			}
			die(vict, ch)
			SET_BIT_AR(ch.Affected_by[:], AFF_LIQUEFIED)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			handle_cooldown(ch, 9)
			return
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: liquefy hide\nSyntax: liquefy explode (target)\r\n"))
		return
	}
}
func do_lifeforce(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg     [2048]byte
		setting int = 0
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: life (0 - 99)\n0 is off.\r\n"))
		return
	}
	setting = libc.Atoi(libc.GoString(&arg[0]))
	if setting > 99 {
		send_to_char(ch, libc.CString("Syntax: life (1 - 99)\n%s isn't an acceptable percent.\r\n"), add_commas(int64(setting)))
		return
	} else if setting <= 0 {
		send_to_char(ch, libc.CString("Your will just isn't in the fight, huh?\nYou will not use up life force to maintain your PL period.\r\n"))
		ch.Lifeperc = 0
		return
	} else {
		send_to_char(ch, libc.CString("Your life force will automatically kick in at %d%s of your optimal PL.\r\n"), setting, "%")
		ch.Lifeperc = setting
		return
	}
}
func do_defend(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 && ch.Defending == nil {
		send_to_char(ch, libc.CString("Defend who?\r\n"))
		return
	} else if arg[0] == 0 && ch.Defending != nil {
		act(libc.CString("@YYou stop defending @y$N@Y.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Defending), TO_CHAR)
		act(libc.CString("@y$n@Y stops defending you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Defending), TO_VICT)
		act(libc.CString("@y$n@Y stops defending @y$N@Y.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Defending), TO_NOTVICT)
		ch.Defending.Defender = nil
		ch.Defending = nil
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("You can't seem to find that person.\r\n"))
		return
	} else if vict == ch {
		send_to_char(ch, libc.CString("Well hopefully you are smart enough to defend yourself.\r\n"))
		return
	} else {
		act(libc.CString("@YYou start defending @y$N@Y.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@y$n@Y starts defending you.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@y$n@Y starts defending @y$N@Y.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Defender = ch
		ch.Defending = vict
		return
	}
}
func do_fish(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: fish ( cast | hook | reel | apply | stop)\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("cast")) == 0 {
		if PLR_FLAGGED(ch, PLR_FISHING) {
			send_to_char(ch, libc.CString("You are already fishing! Syntax: fish stop\r\n"))
			return
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_FISHING) {
			send_to_char(ch, libc.CString("This is not an area you can fish at.\r\n"))
			return
		} else if AFF_FLAGGED(ch, AFF_FLYING) {
			send_to_char(ch, libc.CString("You can't fish while flying.\r\n"))
			return
		} else if ch.Fighting != nil {
			send_to_char(ch, libc.CString("You can't fish while fighting!\r\n"))
			return
		} else if (ch.Equipment[WEAR_WIELD2]) == nil {
			send_to_char(ch, libc.CString("You are not holding a fishing pole.\r\n"))
			return
		} else {
			var pole *obj_data = (ch.Equipment[WEAR_WIELD2])
			if int(pole.Type_flag) != ITEM_FISHPOLE {
				send_to_char(ch, libc.CString("You do not have a fishing pole in your hand!\r\n"))
				return
			} else if (pole.Value[0]) == 0 {
				send_to_char(ch, libc.CString("There is no bait on your line!\r\n"))
				return
			}
			reveal_hiding(ch, 0)
			act(libc.CString("@CYou pull your arm back and then spring it forward, casting the baited line. A moment later there is a splash as the hook enters the water.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C pulls $s arm back and then springs it foward, casting the line of $s fishing pole into the water.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Accuracy_mod = rand_number(30, 80)
			SET_BIT_AR(ch.Act[:], PLR_FISHING)
			send_to_char(ch, libc.CString("@D[@wDistance@D: @Y%d@D]@n\r\n"), ch.Accuracy_mod)
			return
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("hook")) == 0 {
		if !PLR_FLAGGED(ch, PLR_FISHING) {
			send_to_char(ch, libc.CString("You are not even fishing!\r\n"))
			return
		} else if ch.Fishstate == FISH_NOFISH {
			send_to_char(ch, libc.CString("You do not have a fish biting on your line.\r\n"))
			return
		} else if ch.Fishstate == FISH_REELING {
			send_to_char(ch, libc.CString("You are already trying to reel the fish in!\r\n"))
			return
		} else if ch.Fishstate == FISH_HOOKED {
			send_to_char(ch, libc.CString("You already have the fish hooked! Reel it in!\r\n"))
			return
		} else if axion_dice(-18) > ch.Accuracy {
			reveal_hiding(ch, 0)
			act(libc.CString("@CYou pull hard but the fish spits the hook out a second before you pull! You return to waiting for a bite...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C pulls hard on $s fishing line, but a moment later $e frowns and returns to fishing.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Fishstate = FISH_NOFISH
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@CYou pull hard on your line and feel that you have managed to hook the fish! Better @Greel@C it in!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C pulls hard on $s fishing line and starts to struggle with the fish on the other end!@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Fishstate = FISH_HOOKED
			return
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("reel")) == 0 {
		if !PLR_FLAGGED(ch, PLR_FISHING) {
			send_to_char(ch, libc.CString("You are not even fishing!\r\n"))
			return
		} else if ch.Fishstate == FISH_NOFISH {
			send_to_char(ch, libc.CString("You do not have a fish biting on your line.\r\n"))
			return
		} else if ch.Fishstate == FISH_REELING {
			send_to_char(ch, libc.CString("You are already trying to reel the fish in!\r\n"))
			return
		} else if ch.Fishstate != FISH_HOOKED {
			send_to_char(ch, libc.CString("You don't have a fish hooked!\r\n"))
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@CYou begin reeling the fish in slowly.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C begins to reel the line on $s pole in.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Fishstate = FISH_REELING
			return
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("apply")) == 0 {
		if (ch.Equipment[WEAR_WIELD2]) == nil {
			send_to_char(ch, libc.CString("You are not holding a fishing pole.\r\n"))
			return
		} else {
			var pole *obj_data = (ch.Equipment[WEAR_WIELD2])
			if int(pole.Type_flag) != ITEM_FISHPOLE {
				send_to_char(ch, libc.CString("You do not have a fishing pole in your hand!\r\n"))
				return
			} else if (pole.Value[0]) != 0 {
				send_to_char(ch, libc.CString("Your fishing pole already has bait on its hook.\r\n"))
				return
			} else {
				var bait *obj_data
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("Syntax: fish apply (bait)\r\n"))
					return
				}
				if (func() *obj_data {
					bait = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
					return bait
				}()) == nil {
					send_to_char(ch, libc.CString("You don't have that bait.\r\n"))
					return
				} else if int(bait.Type_flag) != ITEM_FISHBAIT {
					send_to_char(ch, libc.CString("That isn't fishing bait!\r\n"))
					return
				} else {
					reveal_hiding(ch, 0)
					act(libc.CString("@CYou carefully apply the $p@C to your hook.@n"), TRUE, ch, bait, nil, TO_CHAR)
					act(libc.CString("@c$n@C carefully applies $p@C to $s fishing pole's hook.@n"), TRUE, ch, bait, nil, TO_ROOM)
					pole.Value[0] = bait.Cost
					extract_obj(bait)
					return
				}
			}
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("stop")) == 0 {
		if !PLR_FLAGGED(ch, PLR_FISHING) {
			send_to_char(ch, libc.CString("You are not even fishing!\r\n"))
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@CYou reel in your line and stop fishing.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reels in $s fishing line and stops fishing.@n"), TRUE, ch, nil, nil, TO_ROOM)
			REMOVE_BIT_AR(ch.Act[:], PLR_FISHING)
			ch.Fishstate = FISH_NOFISH
			ch.Accuracy_mod = 0
			return
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: fish ( cast | hook | reel | apply | stop )\r\n"))
		return
	}
}
func has_pole(ch *char_data) int {
	if (ch.Equipment[WEAR_WIELD2]) != nil {
		var pole *obj_data = (ch.Equipment[WEAR_WIELD2])
		if int(pole.Type_flag) == ITEM_FISHPOLE {
			return TRUE
		}
	}
	return FALSE
}
func fish_update() {
	var (
		i         *char_data
		next_char *char_data
		ch        *char_data = nil
		quality   int        = 0
	)
	for i = character_list; i != nil; i = next_char {
		next_char = i.Next
		if ROOM_FLAGGED(i.In_room, ROOM_FISHING) {
			if PLR_FLAGGED(i, PLR_FISHING) && has_pole(i) == TRUE {
				ch = i
				if ch.Accuracy_mod <= 0 && ch.Fishstate == FISH_REELING {
					if ch.Accuracy >= rand_number(60, 100) {
						quality = rand_number(0, 3) + rand_number(0, 3) + rand_number(0, 3)
					} else if ch.Accuracy >= rand_number(45, 60) {
						quality = rand_number(0, 3) + rand_number(0, 3)
					} else {
						quality = rand_number(0, 3)
					}
					catch_fish(ch, quality)
				} else if rand_number(1, 5) >= 3 {
					if ch.Fishstate == FISH_REELING && rand_number(1, 100) <= 80 {
						if ch.Accuracy >= 80 {
							ch.Accuracy_mod -= rand_number(6, 10)
						} else if ch.Accuracy >= 40 {
							ch.Accuracy_mod -= rand_number(5, 8)
						} else {
							ch.Accuracy_mod -= rand_number(1, 4)
						}
						act(libc.CString("@CYou reel the line on your pole some.@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C reels the line on $s pole slowly.@n"), TRUE, ch, nil, nil, TO_ROOM)
						send_to_char(ch, libc.CString("@D[@wDistance@D: @Y%d@D]@n\r\n"), func() int {
							if ch.Accuracy_mod > 0 {
								return ch.Accuracy_mod
							}
							return 0
						}())
					} else if ch.Fishstate == FISH_REELING && rand_number(1, 58) <= 55 {
						act(libc.CString("@CYou struggle as the fish fights against your attempts to reel it in!@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C struggles with the fish on the end of $s pole!@n"), TRUE, ch, nil, nil, TO_ROOM)
					} else if ch.Fishstate == FISH_REELING {
						act(libc.CString("@CYou feel the line go slack and realize you've lost the fish! You reel your line back in...@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C frowns and then begins to reel in $s line.@n"), TRUE, ch, nil, nil, TO_ROOM)
						ch.Accuracy_mod = 0
						ch.Fishstate = FISH_NOFISH
						REMOVE_BIT_AR(ch.Act[:], PLR_FISHING)
						if has_pole(ch) == TRUE {
							var pole *obj_data = (ch.Equipment[WEAR_WIELD2])
							pole.Value[0] = 0
						}
					} else if ch.Fishstate == FISH_HOOKED && rand_number(1, 20) >= 12 {
						act(libc.CString("@CYou feel the line go slack and realize you've lost the fish! You reel your line back in...@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C frowns and then begins to reel in $s line.@n"), TRUE, ch, nil, nil, TO_ROOM)
						ch.Accuracy_mod = 0
						ch.Fishstate = FISH_NOFISH
						REMOVE_BIT_AR(ch.Act[:], PLR_FISHING)
					} else if ch.Fishstate == FISH_BITE && rand_number(1, 20) >= 12 {
						act(libc.CString("@CYou feel as if the fish has stopped biting...@n"), TRUE, ch, nil, nil, TO_CHAR)
						ch.Fishstate = FISH_NOFISH
					} else if ch.Fishstate != FISH_HOOKED && ch.Fishstate != FISH_BITE && (ROOM_FLAGGED(ch.In_room, ROOM_FISHFRESH) && rand_number(1, 10) >= 8 || !ROOM_FLAGGED(ch.In_room, ROOM_FISHFRESH) && rand_number(1, 20) >= 18) {
						act(libc.CString("@CYou feel a fish biting on your line! Better @Ghook@C it!@n"), TRUE, ch, nil, nil, TO_CHAR)
						ch.Fishstate = FISH_BITE
					}
				}
			} else if PLR_FLAGGED(i, PLR_FISHING) && has_pole(i) == FALSE {
				REMOVE_BIT_AR(i.Act[:], PLR_FISHING)
				i.Accuracy_mod = 0
				i.Fishstate = FISH_NOFISH
			}
		} else {
			if PLR_FLAGGED(i, PLR_FISHING) {
				REMOVE_BIT_AR(i.Act[:], PLR_FISHING)
				i.Accuracy_mod = 0
				i.Fishstate = FISH_NOFISH
			}
		}
	}
}
func catch_fish(ch *char_data, quality int) {
	var (
		fish *obj_data = nil
		num  int       = 1000
	)
	if ROOM_FLAGGED(ch.In_room, ROOM_FISHFRESH) {
		if ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
			switch rand_number(1, 10) {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				num = 1000
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				num = 1001
			case 8:
				fallthrough
			case 9:
				num = 1002
			case 10:
				num = 1003
			}
		} else if ROOM_FLAGGED(ch.In_room, ROOM_AETHER) {
			switch rand_number(1, 10) {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				num = 1012
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				num = 1013
			case 8:
				fallthrough
			case 9:
				num = 1014
			case 10:
				num = 1015
			}
		}
	} else {
		if ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
			switch rand_number(1, 10) {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				num = 1004
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				num = 1005
			case 8:
				fallthrough
			case 9:
				num = 1006
			case 10:
				num = 1007
			}
		} else if ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) {
			switch rand_number(1, 10) {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				num = 1008
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				num = 1009
			case 8:
				fallthrough
			case 9:
				num = 1010
			case 10:
				num = 1011
			}
		}
	}
	fish = read_object(obj_vnum(num), VIRTUAL)
	if fish == nil {
		send_to_imm(libc.CString("Error: Fish success with no fish! Report to Iovan!\r\n"))
		return
	}
	act(libc.CString("@CYou manage to pull a $p@C from the water and onto the ground in front of you!@n"), TRUE, ch, fish, nil, TO_CHAR)
	act(libc.CString("@c$n@C manages to pull a $p@C from the water and onto the ground in front of $m!@n"), TRUE, ch, fish, nil, TO_ROOM)
	var weight int = 1
	if quality <= 0 && rand_number(1, 20) >= 17 {
		quality += rand_number(2, 7)
	}
	var pole *obj_data = (ch.Equipment[WEAR_WIELD2])
	if (pole.Value[0])*2 >= axion_dice(0) {
		quality += 2
	} else if (pole.Value[0]) >= axion_dice(0) {
		quality += 1
	}
	switch quality {
	case 0:
		fallthrough
	case 1:
		weight = rand_number(0, 2)
	case 2:
		fallthrough
	case 3:
		weight = rand_number(3, 4)
		fish.Cost += int(float64(fish.Cost) * 0.2)
		fish.Value[0] += 1
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		weight = rand_number(5, 9)
		fish.Cost += int(float64(fish.Cost) * 0.5)
		fish.Value[0] += 3
	default:
		weight = rand_number(10, 15)
		fish.Cost += fish.Cost * 2
		fish.Value[0] += 5
	}
	fish.Weight += int64(weight)
	pole.Value[0] = 0
	obj_to_room(fish, ch.In_room)
	do_get(ch, libc.CString("fish"), 0, 0)
	send_to_char(ch, libc.CString("@D[@cFish Weight@D: @G%lld@D]@n\r\n"), fish.Weight)
	REMOVE_BIT_AR(ch.Act[:], PLR_FISHING)
	ch.Accuracy_mod = 0
	ch.Fishstate = FISH_NOFISH
}
func do_extract(ch *char_data, argument *byte, cmd int, subcmd int) {
	if GET_SKILL(ch, SKILL_EXTRACT) == 0 {
		send_to_char(ch, libc.CString("You do not know how to extract!\r\n"))
		return
	}
	var arg [2048]byte
	var argu [2048]byte
	var arg2 [2048]byte
	var arg3 [2048]byte
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	var skill int = GET_SKILL(ch, SKILL_EXTRACT)
	var chance int = axion_dice(0)
	half_chop(argument, &arg[0], &argu[0])
	two_arguments(&argu[0], &arg2[0], &arg3[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: extract (object)\r\n"))
		send_to_char(ch, libc.CString("Syntax: extract combine (bottle1) (bottle2)\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("combine")) == 0 {
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You do not have the first bottle that you were wanting to combine.\r\n"))
			return
		}
		if (func() *obj_data {
			obj2 = get_obj_in_list_vis(ch, &arg3[0], nil, ch.Carrying)
			return obj2
		}()) == nil {
			send_to_char(ch, libc.CString("You do not have the second bottle that you were wanting to combine.\r\n"))
			return
		}
		if obj != nil && obj2 != nil {
			if GET_OBJ_VNUM(obj) != 3424 {
				send_to_char(ch, libc.CString("That is not an ink bottle!\r\n"))
				return
			} else if GET_OBJ_VNUM(obj2) != 3424 {
				send_to_char(ch, libc.CString("That is not an ink bottle!\r\n"))
				return
			} else if (obj.Value[6]) <= 0 {
				send_to_char(ch, libc.CString("There isn't any ink in the first bottle!\r\n"))
				return
			} else if (obj2.Value[6]) <= 0 {
				send_to_char(ch, libc.CString("There isn't any ink in the second bottle!\r\n"))
				return
			} else {
				if (obj.Value[6]) >= (obj2.Value[6]) {
					obj.Value[6] += obj2.Value[6]
					if (obj.Value[6]) > 24 {
						obj.Value[6] = 24
					}
					send_to_char(ch, libc.CString("You combine the ink of the two bottles into one bottle, and discard the leftovers.\r\n"))
					extract_obj(obj2)
				} else if (obj2.Value[6]) > (obj.Value[6]) {
					obj2.Value[6] += obj.Value[6]
					if (obj2.Value[6]) > 24 {
						obj2.Value[6] = 24
					}
					send_to_char(ch, libc.CString("You combine the ink of the two bottles into one bottle, and discard the leftovers.\r\n"))
					extract_obj(obj)
				}
				return
			}
		}
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You do not have that item.\r\n"))
		return
	} else {
		if GET_OBJ_VNUM(obj) == 3425 {
			if (obj.Value[VAL_MAXMATURE]) != 0 && (obj.Value[VAL_MATURITY]) < (obj.Value[VAL_MAXMATURE]) {
				send_to_char(ch, libc.CString("It's not mature enough to extract from!\r\n"))
				return
			}
			var bottle *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 3423)
			var cost int64 = int64((float64(ch.Max_mana) * 0.35) + 500)
			if bottle == nil {
				send_to_char(ch, libc.CString("You do not have an empty bottle to put the extracted ink in.\r\n"))
				return
			}
			var extra int64 = 0
			if (bottle.Value[6])+4 >= 24 {
				extra = int64(float64(ch.Max_mana) * 0.5)
			}
			cost += extra
			if ch.Mana < cost {
				send_to_char(ch, libc.CString("You do not have enough ki! @D[@rNeeded@D: @R%s@D]@n\r\n"), add_commas(cost))
				return
			} else if skill < chance {
				ch.Mana -= cost
				act(libc.CString("@WWith your ki flowing carefully into your hands you take a hold of the @G$p@W and begin to strip it of its leaves. Once it has been stripped you go to squeeze the ink carefully from the leaves into the bottle, but unfortunately the ink explodes into a mess instead!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W takes a hold of the @G$p@W and begins to strip it of its leaves. Once it has been stripped $e bundles up the leaves in $s hands and begins to squeeze. A nasty explosion of a mess is all that follows!@n"), TRUE, ch, obj, nil, TO_ROOM)
				improve_skill(ch, SKILL_EXTRACT, 0)
				extract_obj(obj)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				return
			} else {
				ch.Mana -= cost
				act(libc.CString("@WWith your ki flowing carefully into your hands you take a hold of the @G$p@W and begin to strip it of its leaves. Once it has been stripped you go to squeeze the ink carefully from the leaves into the bottle, and manage to get every last drop of ink into it.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W takes a hold of the @G$p@W and begins to strip it of its leaves. Once it has been stripped $e bundles up the leaves in $s hands and begins to squeeze ink carefully from the leaves into a bottle.@n"), TRUE, ch, obj, nil, TO_ROOM)
				extract_obj(obj)
				bottle.Value[6] += rand_number(4, 6)
				if (bottle.Value[6]) >= 24 {
					var filled *obj_data = read_object(3424, VIRTUAL)
					extract_obj(bottle)
					filled.Value[6] = 24
					obj_to_char(filled, ch)
					ch.Mana -= 0
					act(libc.CString("@GAs the last of the ink fills the bottle you infuse a final burst of ki into the bottle.@n"), TRUE, ch, filled, nil, TO_CHAR)
					act(libc.CString("@GAs the last of the ink fills the bottle @g$n@G infuses a final burst of ki into the bottle.@n "), TRUE, ch, filled, nil, TO_ROOM)
				} else {
					send_to_char(ch, libc.CString("You will need to fill the bottle before giving it a final infusion of ki to complete the process.\r\n"))
				}
				improve_skill(ch, SKILL_EXTRACT, 0)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			}
		} else {
			send_to_char(ch, libc.CString("That is not something you can extract from.\r\n"))
			return
		}
	}
}
func do_runic(ch *char_data, argument *byte, cmd int, subcmd int) {
	if GET_SKILL(ch, SKILL_RUNIC) == 0 {
		send_to_char(ch, libc.CString("You do not know how to write down runes.\r\n"))
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	var skill int = GET_SKILL(ch, SKILL_RUNIC)
	var bonus int = 0
	two_arguments(argument, &arg[0], &arg2[0])
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are too busy fighting to write runes!\r\n"))
		return
	}
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: runic (target) (skill)\r\n"))
		send_to_char(ch, libc.CString("@D----@GRunic Skills@D----@n\r\n"))
		send_to_char(ch, libc.CString("@Rkenaz\n%s\n%s\n%s\n%s\n%s\n%s@n\n"), func() string {
			if skill >= 40 {
				return "@Galgiz"
			}
			return ""
		}(), func() string {
			if skill >= 40 {
				return "@moagaz"
			}
			return ""
		}(), func() string {
			if skill >= 50 {
				return "@CLaguz"
			}
			return ""
		}(), func() string {
			if skill >= 60 {
				return "@Ywunjo"
			}
			return ""
		}(), func() string {
			if skill >= 80 {
				return "@rpurisaz"
			}
			return ""
		}(), func() string {
			if skill >= 100 {
				return "@mgebo"
			}
			return ""
		}())
		return
	}
	var bottle *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 3424)
	if bottle == nil {
		send_to_char(ch, libc.CString("You do not have a bottle with enough ink in it.\r\n"))
		return
	}
	var brush *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 3427)
	if brush == nil {
		send_to_char(ch, libc.CString("You do not have a brush!\r\n"))
		return
	}
	var cost int64 = int64(float64(ch.Max_mana) * 0.05)
	var inkcost int = 0
	if int(ch.Race) == RACE_HOSHIJIN {
		bonus = 10
	} else {
		inkcost += 2
	}
	var vict *char_data
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("You can't seem to find that person.\r\n"))
		return
	} else if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki to write runes.\r\n"))
		return
	} else if skill+bonus < axion_dice(0) && rand_number(1, 5) == 5 {
		act(libc.CString("@BYou dip your brush into the ink, but as you infuse your ki you balance the flow wrong and end up destroying the ink bottle!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@b$n@B dips $s runic brush into a bottle filled with shimmering ink. @b$n@B appears to concentrate for a moment before a look of panic dons $s face. Just at that moment the bottle of ink explodes! Strange...@n"), TRUE, ch, nil, nil, TO_ROOM)
		extract_obj(bottle)
		improve_skill(ch, SKILL_RUNIC, 1)
		ch.Mana -= cost
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if skill+bonus < axion_dice(0) {
		ch.Mana -= cost
		act(libc.CString("@BYou dip your brush into the ink, but as you infuse your ki you balance the flow wrong and end up evaporating some ink!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@b$n@B dips $s runic brush into a bottle filled with shimmering ink. @b$n@B appears to concentrate for a moment before some ink evaporates. Strange...@n"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_RUNIC, 1)
		bottle.Value[6] -= rand_number(1, 3)
		if (bottle.Value[6]) < 0 {
			bottle.Value[6] = 0
		}
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("kenaz")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Kenaz")) == 0 {
		inkcost += 1
		if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CKenaz@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CKenaz@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.16)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou can now see in the dark! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_INFRAVISION, SKILL_RUNIC, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CKenaz@D'@B rune on @b$N's@B skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CKenaz@D'@B rune on @RYOUR@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CKenaz@D'@B rune on @b$N's@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.16)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou can now see in the dark! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_INFRAVISION, SKILL_RUNIC, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("algiz")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Algiz")) == 0 {
		inkcost += 2
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CAlgiz@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CAlgiz@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.05)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou now have Ethereal Armor! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_EARMOR, SKILL_PUNCH, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CAlgiz@D'@B rune on @b$N's@B skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CAlgiz@D'@B rune on @RYOUR@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CAlgiz@D'@B rune on @b$N's@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.05)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou now have Ethereal Armor! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_EARMOR, SKILL_PUNCH, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("oagaz")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Oagaz")) == 0 {
		inkcost += 3
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@COagaz@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@COagaz@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.04)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou now are protected by Ethereal Chains! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_ECHAINS, SKILL_KNEE, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("laguz")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Laguz")) == 0 {
		inkcost += 4
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CLaguz@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CLaguz@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.04)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou now have water breathing! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_WATERBREATH, SKILL_SBC, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@COagaz@D'@B rune on @b$N's@B skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@COagaz@D'@B rune on @RYOUR@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@COagaz@D'@B rune on @b$N's@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.04)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou now are protected by Ethereal Chains! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_ECHAINS, SKILL_KNEE, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("wunjo")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Wunjo")) == 0 {
		inkcost += 4
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CWunjo@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CWunjo@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.08)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou are now blessed with a deeper understanding of things you experience! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_WUNJO, SKILL_SLAM, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CWunjo@D'@B rune on @b$N's@B skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CWunjo@D'@B rune on @RYOUR@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CWunjo@D'@B rune on @b$N's@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.08)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou are now blessed with a deeper understanding of things you experience! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_WUNJO, SKILL_SLAM, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("purisaz")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Purisaz")) == 0 {
		inkcost += 4
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CPurisaz@D'@B rune on your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CPurisaz@D'@B rune on $s skin.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.06)
			if duration < 1 {
				duration = 1
			}
			send_to_char(vict, libc.CString("@GYou feel as if your inner energy is more potent! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			assign_affect(vict, AFF_POTENT, SKILL_HEELDROP, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CPurisaz@D'@B rune on @b$N's@B skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CPurisaz@D'@B rune on @RYOUR@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CPUrisaz@D'@B rune on @b$N's@B skin.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			var duration int = int(float64(skill) * 0.06)
			send_to_char(vict, libc.CString("@GYou feel as if your inner energy is more potent! @D(@WLasts@D: @w%d@D)@n\r\n"), duration)
			if duration < 1 {
				duration = 1
			}
			assign_affect(vict, AFF_POTENT, SKILL_HEELDROP, duration, 0, 0, 0, 0, 0, 0)
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else if libc.StrCaseCmp(&arg2[0], libc.CString("gebo")) == 0 || libc.StrCaseCmp(&arg2[0], libc.CString("Gebo")) == 0 {
		inkcost += 10
		if (bottle.Value[6]) < inkcost {
			send_to_char(ch, libc.CString("You do not have a bottle with enough ink. @D[@bInkcost@D: @R%d@D]@n\r\n"), inkcost)
			return
		} else if vict == ch {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CGebo@D'@B rune on your skin! The rune flashes out of existence immediately!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CGebo@D'@B rune on $s skin. The rune flashes out of existence immediately!@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			vict.Player_specials.Class_skill_points[vict.Chclass] += 125
			send_to_char(vict, libc.CString("@GYou feel like you've just gained a lot of knowledge. Now if only you could apply it. @D[@m+125 PS@D]@n\r\n"))
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		} else {
			ch.Mana -= cost
			act(libc.CString("@BYou dip your brush into the ink and infuse your ki skillfully into it. You pull the brush out and paint the @D'@CGebo@D'@B rune on @b$N's@B skin! The rune flashes out of existence immediately!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CGebo@D'@B rune on @RYOUR@B skin. The rune flashes out of existence immediately!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@b$n@B dips $s brush into a bottle of ink and at the same time the ink starts to glow. Skillfully $e then writes the @D'@CGebo@D'@B rune on @b$N's@B skin. The rune flashes out of existence immediately!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D[@B%d@b ink used.@D]@n\r\n"), inkcost)
			vict.Player_specials.Class_skill_points[vict.Chclass] += 125
			send_to_char(vict, libc.CString("@GYou feel like you've just gained a lot of knowledge. Now if only you could apply it. @D[@m+125 PS@D]@n\r\n"))
			bottle.Value[6] -= inkcost
			if (bottle.Value[6]) <= 0 {
				extract_obj(bottle)
				var empty *obj_data = read_object(3423, VIRTUAL)
				obj_to_char(empty, ch)
			}
		}
		improve_skill(ch, SKILL_RUNIC, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else {
		send_to_char(ch, libc.CString("Syntax: runic (target) (skill)\r\n"))
		send_to_char(ch, libc.CString("@D----@GRunic Skills@D----@n\r\n"))
		send_to_char(ch, libc.CString("@Rkenaz\n%s\n%s\n%s\n%s\n%s@n\n"), func() string {
			if skill >= 40 {
				return "@Galgiz"
			}
			return ""
		}(), func() string {
			if skill >= 40 {
				return "@moagaz"
			}
			return ""
		}(), func() string {
			if skill >= 60 {
				return "@Ywunjo"
			}
			return ""
		}(), func() string {
			if skill >= 80 {
				return "@rpurisaz"
			}
			return ""
		}(), func() string {
			if skill >= 100 {
				return "@mgebo"
			}
			return ""
		}())
		return
	}
}
func do_scry(ch *char_data, argument *byte, cmd int, subcmd int) {
	if libc.StrCaseCmp(CAP(GET_NAME(ch)), libc.CString("Galeos")) != 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique.\r\n"))
		return
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: scry (target)\r\n"))
		return
	}
	var vict *char_data
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Who are you using Oracle Scry on?\r\n"))
		return
	}
	if vict == ch {
		send_to_char(ch, libc.CString("You can't do that to yourself!\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("No using this on mobs!\r\n"))
		return
	}
	var cost int = 2000
	if (ch.Player_specials.Class_skill_points[ch.Chclass]) < cost {
		send_to_char(ch, libc.CString("You do not have enough PS to Oracle Scry!\r\n"))
		return
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("@GYou focus your mind and begin to allow the flood of images and energy to roar through your mind. You then allow those thoughts to make their way into the mind of @c$N@G. You can hardly comprehend the vastness of the information flooding in, yet still glimpse bits and pieces of your own destiny.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@GYou see @C$n@G begin to focus, and then without warning, your mind is flooded painfully with images, energy and information. The data streams in a mad torrent through your psyche, and just when you think snapping is possible, the voice of @C$n@G comes to you and eases and guides you. You see images of potential futures, information not yet known, knowledge yet undiscovered. Though you could not fully  grasp what is to come, you feel more prepared at facing the unknown.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W appears to be performing some sort of ritual or something with @c$N@W.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		var boost int64 = int64(float64(ch.Aff_abils.Intel) * 0.5)
		var mult float64 = 1
		if int(vict.Race) == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 3
		} else if int(vict.Race) == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 4
		} else if int(vict.Race) == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 5
		} else if int(vict.Race) == RACE_HOSHIJIN && vict.Starphase == 1 {
			mult = 2
		} else if int(vict.Race) == RACE_HOSHIJIN && vict.Starphase == 2 {
			mult = 3
		} else if int(vict.Race) == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if int(vict.Race) == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if int(vict.Race) == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 3.5
		} else if int(vict.Race) == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS4) {
			mult = 4
		}
		vict.Max_hit += int64(((float64(vict.Basepl) * 0.01) * float64(boost)) * mult)
		vict.Basepl += int64((float64(vict.Basepl) * 0.01) * float64(boost))
		vict.Max_mana += int64(((float64(vict.Baseki) * 0.01) * float64(boost)) * mult)
		vict.Baseki += int64((float64(vict.Baseki) * 0.01) * float64(boost))
		vict.Max_move += int64(((float64(vict.Basest) * 0.01) * float64(boost)) * mult)
		vict.Basest += int64((float64(vict.Basest) * 0.01) * float64(boost))
		send_to_char(vict, libc.CString("Your Powerlevel, Ki, and Stamina have improved drastically! On top of that your Intelligence and Wisdom have improved permanantly!\r\n"))
		vict.Real_abils.Intel += 2
		vict.Real_abils.Wis += 2
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 2000
		if GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("@D[@mPractice Sessions@D:@R -2000@D]@n\r\n"))
			if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(level_exp(ch, GET_LEVEL(ch)+1) - int(ch.Exp))
				send_to_char(ch, libc.CString("The remaining experience needed for your next level up has been gained!@n\r\n"))
			} else {
				send_to_char(ch, libc.CString("Due to already having enough experience to level up you gain no expereince.\r\n"))
			}
		} else {
			ch.Max_hit += int64(float64(ch.Basepl) * 0.025)
			ch.Basepl += int64(float64(ch.Basepl) * 0.025)
			ch.Max_move += int64(float64(ch.Basest) * 0.025)
			ch.Basest += int64(float64(ch.Basest) * 0.025)
			ch.Max_mana += int64(float64(ch.Baseki) * 0.025)
			ch.Baseki += int64(float64(ch.Baseki) * 0.025)
			send_to_char(ch, libc.CString("Your Powerlevel, Ki, and Stamina have improved!\r\n"))
		}
	}
}
func ash_burn(ch *char_data) {
	var (
		obj      *obj_data
		next_obj *obj_data
	)
	if ch != nil && ch.In_room != room_rnum(-1) {
		for obj = world[ch.In_room].Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if GET_OBJ_VNUM(obj) == 1306 {
				if axion_dice(0) > int(ch.Aff_abils.Con) {
					if int(ch.Race) != RACE_ANDROID && int(ch.Race) != RACE_DEMON && int(ch.Race) != RACE_ICER {
						reveal_hiding(ch, 0)
						ch.Move -= int64(((float64(ch.Max_move) * 0.005) + 20) * float64(obj.Cost))
						if ch.Move < 0 {
							ch.Move = 0
						}
						act(libc.CString("@RYou choke on the the burning hot @Da@Ws@wh@Dc@Wl@wo@Du@Wd@R!@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@r$n@R chokes on the burning hot @Da@Ws@wh@Dc@Wl@wo@Du@Wd@R!@n"), TRUE, ch, nil, nil, TO_ROOM)
					}
					if int(ch.Race) != RACE_ANDROID && int(ch.Race) != RACE_DEMON && !IS_NPC(ch) {
						if !PLR_FLAGGED(ch, PLR_EYEC) && !AFF_FLAGGED(ch, AFF_BLIND) {
							reveal_hiding(ch, 0)
							act(libc.CString("@DYour eyes sting from the hot ash! You can't see!@n"), TRUE, ch, nil, nil, TO_CHAR)
							act(libc.CString("@r$n@D eyes appear to have been hurt by the ash!@n"), TRUE, ch, nil, nil, TO_ROOM)
							var duration int = 1
							assign_affect(ch, AFF_BLIND, SKILL_SOLARF, duration, 0, 0, 0, 0, 0, 0)
						}
					}
				}
			}
		}
	}
}
func do_ashcloud(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Race) != RACE_DEMON {
		send_to_char(ch, libc.CString("You are not trained in the use of ash and fire!\r\n"))
		return
	}
	var level int = 1
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: ashcloud (1 | 2 | 3)\r\n"))
		return
	}
	var ash *obj_data = find_obj_in_list_vnum(ch.Carrying, 1305)
	var there *obj_data = find_obj_in_list_vnum(world[ch.In_room].Contents, 1306)
	if there != nil {
		send_to_char(ch, libc.CString("You can not pile more ash into the air without causing it to clump together and settle.\r\n"))
		return
	}
	if ash == nil {
		send_to_char(ch, libc.CString("You do not have any ash!\r\n"))
		return
	}
	level = libc.Atoi(libc.GoString(&arg[0]))
	var mult int64 = 5
	var initial float64 = 0.0
	switch level {
	case 1:
		mult = 20
		initial = 0.25
	case 2:
		mult = 10
		initial = 0.1
	case 3:
		mult = 5
		initial = 0.05
	default:
		send_to_char(ch, libc.CString("Syntax: ashcloud (1 | 2 | 3)\r\n"))
		return
	}
	var cost int64 = int64((float64(ch.Max_mana) * initial) + float64(int64(ch.Aff_abils.Intel)*mult))
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki!\r\n"))
		return
	} else if SUNKEN(ch.In_room) {
		send_to_char(ch, libc.CString("You can not create an ashcloud here, because it is too wet.\r\n"))
		return
	} else if SECT(ch.In_room) == SECT_SPACE {
		send_to_char(ch, libc.CString("You can not create an ashcloud in space.\r\n"))
		return
	} else if int(ch.Aff_abils.Intel) < axion_dice(-10) {
		reveal_hiding(ch, 0)
		act(libc.CString("@RYou take a handful of ashes, and when you go to blow flames at it you lose focus. The ashes are blown from your hands by your huge gust of breath.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@r$n@R takes a handful of ashes from $s belongings and blows it out of $s hands with a strong gust of air. @YStrange.@n"), TRUE, ch, nil, nil, TO_ROOM)
		extract_obj(ash)
		ch.Mana -= cost
		return
	} else {
		var ashcloud *obj_data
		reveal_hiding(ch, 0)
		if level == 3 {
			ch.Mana -= cost
			act(libc.CString("@RYou take a handful of ashes and you create a fierce heat within your lungs. With the heat ready you breathe ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R takes a handful of ashes and $e breathes ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_room(ch.In_room, libc.CString("@WThe ashes ripple with an intense aftershock of power.@n\r\n"))
			ashcloud = read_object(1306, VIRTUAL)
			obj_to_room(ashcloud, ch.In_room)
			extract_obj(ash)
			ashcloud.Timer = 4
			ashcloud.Cost = 3
			ash_burn(ch)
		} else if level == 2 {
			ch.Mana -= cost
			act(libc.CString("@RYou take a handful of ashes and you create a fierce heat within your lungs. With the heat ready you breathe ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R takes a handful of ashes and $e breathes ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_room(ch.In_room, libc.CString("@WThe ashes ripple with a strong aftershock of power.@n\r\n"))
			ashcloud = read_object(1306, VIRTUAL)
			obj_to_room(ashcloud, ch.In_room)
			ashcloud.Timer = 2
			ashcloud.Cost = 2
			extract_obj(ash)
			ash_burn(ch)
		} else {
			ch.Mana -= cost
			act(libc.CString("@RYou take a handful of ashes and you create a fierce heat within your lungs. With the heat ready you breathe ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R takes a handful of ashes and $e breathes ki infused flames at the pile of ashes! The flames and ashes mix and fill the surrounding area with a hot burning ash!@n"), TRUE, ch, nil, nil, TO_ROOM)
			ashcloud = read_object(1306, VIRTUAL)
			obj_to_room(ashcloud, ch.In_room)
			extract_obj(ash)
			ashcloud.Timer = 1
			ashcloud.Cost = 1
			ash_burn(ch)
		}
	}
}
func do_resize(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
		obj  *obj_data
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if GET_SKILL(ch, SKILL_BUILD) == 0 {
		send_to_char(ch, libc.CString("You do not have the skill to resize equipment!\r\n"))
		return
	} else if GET_SKILL(ch, SKILL_BUILD) < 80 {
		send_to_char(ch, libc.CString("Your build skill must be at least level 80 before you can resize equipment.\r\n"))
		return
	} else {
		if arg[0] == 0 || arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: resize (obj) (small | medium)\r\n"))
			return
		}
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have that object!\r\n"))
			return
		} else {
			if wearable_obj(obj) == 0 {
				send_to_char(ch, libc.CString("That is not equipment! You can only resize equipment.\r\n"))
				return
			} else {
				if ch.Move < obj.Weight+ch.Max_move/40 {
					send_to_char(ch, libc.CString("You do not have enough stamina to resize this object at this time.\r\n"))
					return
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("small")) == 0 {
					if obj.Size == SIZE_SMALL {
						send_to_char(ch, libc.CString("The equipment is already small sized.\r\n"))
						return
					} else {
						act(libc.CString("@WYou carefully adjust the size of @c$p@W.@n"), TRUE, ch, obj, nil, TO_CHAR)
						act(libc.CString("@C$n@W carefully adjusts the size of @c$p@W.@n"), TRUE, ch, obj, nil, TO_ROOM)
						obj.Size = SIZE_SMALL
						ch.Move -= obj.Weight + ch.Max_move/40
					}
				} else if libc.StrCaseCmp(&arg2[0], libc.CString("medium")) == 0 {
					if obj.Size == SIZE_MEDIUM {
						send_to_char(ch, libc.CString("The equipment is already medium sized.\r\n"))
						return
					} else {
						act(libc.CString("@WYou carefully adjust the size of @c$p@W.@n"), TRUE, ch, obj, nil, TO_CHAR)
						act(libc.CString("@C$n@W carefully adjusts the size of @c$p@W.@n"), TRUE, ch, obj, nil, TO_ROOM)
						obj.Size = SIZE_MEDIUM
						ch.Move -= obj.Weight + ch.Max_move/40
					}
				} else {
					send_to_char(ch, libc.CString("Syntax: resize (obj) (small | medium)\r\n"))
				}
			}
		}
	}
}
func do_healglow(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if GET_SKILL(ch, SKILL_HEALGLOW) == 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique.\r\n"))
		return
	}
	if arg[0] == 0 {
		vict = ch
	} else if GET_SKILL(ch, SKILL_HEALGLOW) < 100 {
		send_to_char(ch, libc.CString("You can not target anyone except yourself unless you are a master of this technique.\nSyntax: healingglow\r\n"))
		return
	} else if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Nobody around by that name.\r\n"))
		return
	}
	if AFF_FLAGGED(vict, AFF_HEALGLOW) && vict == ch {
		send_to_char(ch, libc.CString("You already have a healing glow surrounding your body.\r\n"))
		return
	} else if AFF_FLAGGED(vict, AFF_HEALGLOW) {
		send_to_char(ch, libc.CString("They already have a healing glow surrounding their body.\r\n"))
		return
	}
	if vict.Fighting != nil && vict == ch {
		send_to_char(ch, libc.CString("You are too busy fighting!@n\r\n"))
		return
	} else if vict.Fighting != nil {
		send_to_char(ch, libc.CString("They are too busy fighting!@n\r\n"))
		return
	}
	var cost int64 = int64(float64(ch.Max_mana) * 0.5)
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki. It requires at least 50%s of your ki in cost.\r\n"), "%")
		return
	} else {
		if vict == ch {
			act(libc.CString("@CPlacing your hands on your body you begin to focus your energies. Slowly a strong blue glow glistens and shines across your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C places $s hands on $s body. Slowly a strong blue glow glistens and shines across $s skin!@n"), TRUE, ch, nil, nil, TO_ROOM)
			SET_BIT_AR(vict.Affected_by[:], AFF_HEALGLOW)
			var duration int = int(float64(GET_SKILL(ch, SKILL_HEALGLOW)) * 0.1)
			if duration <= 0 {
				duration = 1
			}
			assign_affect(ch, AFF_HEALGLOW, SKILL_HEALGLOW, duration, 0, 0, 0, 0, 0, 0)
			ch.Mana -= cost
		} else {
			act(libc.CString("@CPlacing your hands on @c$N's@C body you begin to focus your energies. Slowly a strong blue glow glistens and shines across $S skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@c$n@C places $s hands on YOUR body. Slowly a strong blue glow glistens and shines across your skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$n@C places $s hands on @c$N's@C body. Slowly a strong blue glow glistens and shines across $S skin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			var duration int = int(float64(GET_SKILL(ch, SKILL_HEALGLOW)) * 0.1)
			duration += rand_number(-2, 1)
			if duration <= 0 {
				duration = 1
			}
			assign_affect(ch, AFF_HEALGLOW, SKILL_HEALGLOW, duration, 0, 0, 0, 0, 0, 0)
			ch.Mana -= cost
		}
	}
}
func do_amnisiac(ch *char_data, argument *byte, cmd int, subcmd int) {
	if libc.StrCaseCmp(GET_NAME(ch), libc.CString("Kanashimi")) != 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique.\r\n"))
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	var vict *char_data
	var skill int
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: amnesiac (target) (skill)\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Perform amnesiac kiss on whom?\r\n"))
		return
	}
	skill = find_skill_num(&arg2[0], 1<<1)
	if skill <= 0 {
		send_to_char(ch, libc.CString("That is not a skill\r\n"))
		return
	}
	var chance int = axion_dice(0)
	var perc int = int(ch.Aff_abils.Intel) + 10
	var cost int = int(float64(ch.Max_mana) * 0.18)
	if cost > int(ch.Mana) {
		send_to_char(ch, libc.CString("You do not have enough ki!\r\n"))
		return
	} else if perc < chance {
		act(libc.CString("@WYou attempt to grab @C$N@W to kiss $M, but $E evades!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@M$n@W attempts to grab you and leans in with puckered lips, but you managed to evade!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@M$n@W attempts to grab @C$N@W and leans in with puckered lips, but $E manages to evade!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= int64(cost)
		return
	} else {
		act(libc.CString("@WYou reach out quickly, grabbing @C$N@W and pulling them in close to you. Just as quick, you pull their head forcefully towards yours, planting a deep and heavy kiss. The fool wobbles a bit, shocked.  It is unlikely that @c$E@W will be able to focus very well on that skill, not with the thought of your lips on their mind."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@M$n@W grabs you, giving you a deep, passionate kiss. Your mind is suddenly overwhelmed, and in your shock you seem to forget a few of your tricks."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@WYou see @C$n@W quickly grab @C$N@W, pulling them into a deep, almost passionate kiss. @C$N@W seems shocked, and wobbles a bit, grabbing at @c$s@W head once @C$n@W lets go."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= int64(cost)
		vict.Stupidkiss = int16(skill)
		return
	}
}
func do_metamorph(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_METAMORPH) == 0 {
		return
	}
	if ch.Alignment >= 51 {
		send_to_char(ch, libc.CString("Your heart is too pure to use that technique!\r\n"))
		return
	}
	var cost int64 = int64(float64(ch.Max_mana) * 0.16)
	if AFF_FLAGGED(ch, AFF_METAMORPH) {
		send_to_char(ch, libc.CString("You are already surrounded by a dark aura!\r\n"))
		return
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki. You need %s.\r\n"), add_commas(cost))
		return
	}
	var chance int = axion_dice(0)
	var perc int = (int(ch.Aff_abils.Wis) * 2)
	if perc < 100 && perc > 60 {
		perc += 100 - perc
	} else if perc < 100 {
		perc += 10
	}
	if perc < chance {
		act(libc.CString("@WYou focus your energies and prepare your @RDark Metamorphisis@W but screw up your focus!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@WA dark @Rred@W glow starts to surround @C$n@W, but it fades quickly.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= cost / 2
		return
	} else {
		act(libc.CString("'@RDark@W...' An explosion of sanguine aura erupts over the surface of your body, your eyes darkening to a bleeding crimson. The flaring glow emanating from your body pronounces the shadows cast, a darkening umbrage that threatens a malicious promise. Fists clench tightly, muscles bulking as you hiss; You complete the transition, relaxing visibly, '...@RMetamorphosis@W'@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("'@RDark@W...' An explosion of sanguine aura erupts over the surface of @C$n@W's body, $s eyes darkening to a bleeding crimson. The flaring glow emanating from $s body pronounces the shadows cast, a darkening umbrage that threatens a malicious promise. Fists clench tightly, muscles bulking as $e hisses; $e completes the transition, relaxing visibly, '...@RMetamorphosis@W'@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= cost
		var duration int = int(ch.Aff_abils.Intel) / 12
		assign_affect(ch, AFF_METAMORPH, SKILL_METAMORPH, duration, 0, 0, 0, 0, 0, 0)
		ch.Hit += int64(float64(gear_pl(ch)) * 0.6)
		return
	}
}
func do_shimmer(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		skill    int        = 0
		perc     int        = 0
		location int        = 0
		cost     int64      = 0
		tar      *char_data = nil
		arg      [2048]byte = func() [2048]byte {
			var t [2048]byte
			copy(t[:], []byte(""))
			return t
		}()
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
		if PRF_FLAGGED(ch, PRF_ARENAWATCH) {
			REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_ARENAWATCH)
			ch.Arenawatch = -1
			send_to_char(ch, libc.CString("You stop watching the arena action.\r\n"))
		}
	}
	if libc.StrCaseCmp(GET_NAME(ch), libc.CString("Anubis")) != 0 {
		send_to_char(ch, libc.CString("You do not even know how to perform that skill!\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 19800 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 0x4DBB {
		send_to_char(ch, libc.CString("@rYou are in a pocket dimension!@n\r\n"))
		return
	} else if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who or where do you want to shimmer to? [target | planet-(planet name) | afterlife]\r\n"))
		send_to_char(ch, libc.CString("Example: shimmer goku\nExample 2: shimmer planet-earth\r\n"))
		return
	}
	cost = ch.Max_mana / 40
	if ch.Mana-cost < 0 {
		send_to_char(ch, libc.CString("You do not have enough ki to instantaneously move.\r\n"))
		return
	}
	perc = axion_dice(0)
	skill = 100
	if libc.StrCaseCmp(&arg[0], libc.CString("planet-earth")) == 0 {
		location = 300
	} else if libc.StrCaseCmp(&arg[0], libc.CString("planet-namek")) == 0 {
		location = 0x27EE
	} else if libc.StrCaseCmp(&arg[0], libc.CString("planet-frigid")) == 0 {
		location = 4017
	} else if libc.StrCaseCmp(&arg[0], libc.CString("planet-vegeta")) == 0 {
		location = 2200
	} else if libc.StrCaseCmp(&arg[0], libc.CString("planet-konack")) == 0 {
		location = 8006
	} else if libc.StrCaseCmp(&arg[0], libc.CString("planet-aether")) == 0 {
		location = 0x2EF8
	} else if libc.StrCaseCmp(&arg[0], libc.CString("afterlife")) == 0 {
		location = 6000
	} else if (func() *char_data {
		tar = get_char_vis(ch, &arg[0], nil, 1<<1)
		return tar
	}()) == nil {
		send_to_char(ch, libc.CString("@RThat target doesn't exist.@n\r\n"))
		send_to_char(ch, libc.CString("Who or where do you want to shimmer to? [target | planet-(planet name) | afterlife]\r\n"))
		send_to_char(ch, libc.CString("Example: shimmer goku\nExample 2: shimmer planet-earth\r\n"))
		return
	}
	if skill < perc || ch.Fighting != nil && rand_number(1, 2) <= 1 {
		if tar != nil {
			if tar != ch {
				send_to_char(ch, libc.CString("You prepare to move instantly but mess up the process and waste some of your ki!\r\n"))
				ch.Mana -= cost
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			} else {
				send_to_char(ch, libc.CString("Moving to yourself would be kinda impossible wouldn't it? If not that then it would at least be pointless.\r\n"))
				return
			}
		} else {
			send_to_char(ch, libc.CString("You prepare to move instantly but mess up the process and waste some of your ki!\r\n"))
			ch.Mana -= cost
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		}
	}
	reveal_hiding(ch, 0)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	if tar != nil {
		if tar == ch {
			send_to_char(ch, libc.CString("Moving to yourself would be kinda impossible wouldn't it? If not that then it would at least be pointless.\r\n"))
			return
		} else if ch.Grappling != nil && ch.Grappling == tar {
			send_to_char(ch, libc.CString("You are already in the same room with them and are grappling with them!\r\n"))
			return
		} else if tar.Admlevel > 0 && ch.Admlevel < 1 {
			send_to_char(ch, libc.CString("That immortal prevents you from reaching them.\r\n"))
			return
		} else if ROOM_FLAGGED(tar.In_room, ROOM_NOINSTANT) {
			send_to_char(ch, libc.CString("You can not go there as it is a protected area!\r\n"))
			return
		} else if ch.Grappling != nil && AFF_FLAGGED(ch.Grappling, AFF_SPIRIT) {
			send_to_char(ch, libc.CString("You can not take the dead with you!\r\n"))
			return
		} else if ch.Drag != nil && AFF_FLAGGED(ch.Drag, AFF_SPIRIT) {
			send_to_char(ch, libc.CString("You can not take the dead with you!\r\n"))
			return
		} else if ch.Grappled != nil && AFF_FLAGGED(ch.Grappled, AFF_SPIRIT) {
			send_to_char(ch, libc.CString("You can not take the dead with you!\r\n"))
			return
		}
		ch.Mana -= cost
		act(libc.CString("@wYour body begins to fade away almost appearing ghost like, before a ripple passes through your image and your are gone in an instant!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_CHAR)
		act(libc.CString("@w$n@w appears in an instant out of nowhere right next to you!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_VICT)
		act(libc.CString("@w$n@w body begins to fade away almost appearing ghost like, before a ripple passes through $s image and $e is gone in an instant!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_NOTVICT)
		SET_BIT_AR(ch.Act[:], PLR_TRANSMISSION)
		handle_teleport(ch, tar, 0)
	} else {
		ch.Mana -= cost
		act(libc.CString("@wYour body begins to fade away almost appearing ghost like, before a ripple passes through your image and your are gone in an instant!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_CHAR)
		act(libc.CString("@w$n@w body begins to fade away almost appearing ghost like, before a ripple passes through $s image and $e is gone in an instant!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_NOTVICT)
		handle_teleport(ch, nil, location)
	}
}
func is_cold_ruby(obj *obj_data) bool {
	return GET_OBJ_VNUM(obj) == 6600 && !OBJ_FLAGGED(obj, ITEM_HOT)
}
func do_channel(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Chclass) != CLASS_DABURA || int(ch.Skills[SKILL_STYLE]) <= 0 {
		send_to_char(ch, libc.CString("You do not know how to do that!\r\n"))
		return
	}
	var cost int64 = int64(float64(ch.Max_mana) * 0.15)
	var chance int = axion_dice(0)
	var skill int = GET_SKILL(ch, SKILL_STYLE)
	if cost > ch.Mana {
		send_to_char(ch, libc.CString("You do not have enough ki to channel with!\r\n"))
		return
	}
	var ruby *obj_data = find_obj_in_list_lambda(ch.Carrying, is_cold_ruby)
	if ruby == nil {
		send_to_char(ch, libc.CString("You do not have any uncharged blood rubies.\r\n"))
		return
	}
	if world[ch.In_room].Geffect <= 0 {
		send_to_char(ch, libc.CString("There is no lava here!\r\n"))
		return
	}
	if ruby != nil {
		if skill < chance {
			act(libc.CString("@RAs you move your ki through the lava you begin to draw heat away from it into the ruby. You screw up the rate of heating though and cause the ruby to crumble to dust!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@RAs $n@R moves $s ki through the lava $e begins to draw heat away from it into a blood ruby. However $e screws up the rate of heating and causes the ruby to crumble to dust!@n"), TRUE, ch, nil, nil, TO_ROOM)
			extract_obj(ruby)
		} else {
			act(libc.CString("@RAs you move your ki through the lava you begin to draw heat away from it into the ruby. You do so at an even rate and end up with a glowing red hot blood ruby!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@RAs $n@R moves $s ki through the lava $e begins to draw heat away from it into a blood ruby. The ruby glows red hot as $e finishes the process of channeling the heat!@n"), TRUE, ch, nil, nil, TO_ROOM)
			world[ch.In_room].Geffect = 0
			SET_BIT_AR(ruby.Extra_flags[:], ITEM_HOT)
		}
		ch.Mana -= cost
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
	}
}
func do_hydromancy(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Chclass) != CLASS_TSUNA || int(ch.Skills[SKILL_STYLE]) <= 0 {
		send_to_char(ch, libc.CString("You know nothing about hydromancy!\r\n"))
		return
	}
	var skill int = int(ch.Skills[SKILL_STYLE])
	var chance int = axion_dice(0)
	var cost int64 = 0
	cost = (ch.Max_mana / 12) - int64(int(ch.Aff_abils.Intel)*GET_LEVEL(ch))
	if world[ch.In_room].Geffect >= 0 && SECT(ch.In_room) != SECT_WATER_SWIM && SECT(ch.In_room) != SECT_WATER_NOSWIM {
		if SECT(ch.In_room) != SECT_UNDERWATER {
			send_to_char(ch, libc.CString("There is not sufficient water here.\r\n"))
			return
		} else {
			send_to_char(ch, libc.CString("There is too much water here to control!\r\n"))
			return
		}
	}
	if cost <= 0 {
		cost = 100
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki to manipulate any water around you.\r\n"))
		return
	}
	if ch.Con_cooldown > 0 {
		send_to_char(ch, libc.CString("You must wait a short period before concentrating again.\r\n"))
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax 1: hydromancy flood (direction)\r\n"))
		send_to_char(ch, libc.CString("Example: hydromancy flood nw\r\n"))
		send_to_char(ch, libc.CString("\nSyntax 2: hydromancy spike\r\n"))
		return
	}
	var attempt int = 0
	if libc.StrCaseCmp(&arg[0], libc.CString("spike")) == 0 {
		var obj *obj_data
		cost = int64((float64(GET_SKILL(ch, SKILL_STYLE)) / ((float64(ch.Max_mana) * 0.5) + 1)) + 100)
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to form an ice spike.\r\n"))
			return
		}
		if skill < chance {
			ch.Mana -= cost
			act(libc.CString("@CYou press your palms together in front of your body but you fail to produce the proper control to form the spike!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C presses $s palms together and then slowly pulls them apart. Nothing important appears to have happened."), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_STYLE, 2)
			return
		}
		if skill >= 100 {
			obj = read_object(0x4A72, VIRTUAL)
		} else if skill >= 50 {
			obj = read_object(0x4A71, VIRTUAL)
		} else if skill >= 1 {
			obj = read_object(0x4A70, VIRTUAL)
		}
		ch.Mana -= cost
		act(libc.CString("@CYou press your palms together in front of your body and focusing ki you force water up along your body. That water pools between your palms and as pull your palms apart a @c$p@C forms!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@c$n@C presses $s palms together in front of $s body and water begins to flow up $s body and pools between $s palms. Slowly pulling them apart reveals a @c$p@C as it forms between them!@n"), TRUE, ch, obj, nil, TO_VICT)
		if obj.Weight+int64(gear_weight(ch)) <= max_carry_weight(ch) {
			obj_to_char(obj, ch)
		} else {
			send_to_char(ch, libc.CString("You are unable to hold it and so let it go at your feet.\r\n"))
			act(libc.CString("@C$n@w drops an ice spike.@n"), TRUE, ch, nil, nil, TO_ROOM)
			obj_to_room(obj, ch.In_room)
		}
		improve_skill(ch, SKILL_STYLE, 1)
		ch.Con_cooldown = 10
	} else if libc.StrCaseCmp(&arg[0], libc.CString("flood")) == 0 {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax 1: hydromancy flood (direction)\r\n"))
			send_to_char(ch, libc.CString("Example: hydromancy flood nw\r\n"))
			send_to_char(ch, libc.CString("\nSyntax 2: hydromancy spike\r\n"))
			return
		}
		attempt = search_block(&arg2[0], &dirs[0], FALSE)
		if CAN_GO(ch, attempt) {
			var (
				vict   *char_data
				next_v *char_data
				last   int = ch.Lastattack
			)
			ch.Lastattack = 500
			var bun [64936]byte
			var bunn [64936]byte
			if skill < chance {
				act(libc.CString("@BUsing your ki you attempt to create a rush of water! @RYou fail!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@b$n@B seems to attempt to create water with $s ki! @RHowever, $e fails!@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Mana -= cost
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			} else {
				ch.Mana -= cost
				stdio.Sprintf(&bun[0], "@BUsing your ki you create a rush of water flooding away toward the @C%s@B!@n", dirs[attempt])
				stdio.Sprintf(&bunn[0], "@B$n@B uses $s ki to create a rush of water flooding away toward the @C%s@B!@n", dirs[attempt])
				act(&bun[0], TRUE, ch, nil, nil, TO_CHAR)
				act(&bunn[0], TRUE, ch, nil, nil, TO_ROOM)
				for vict = world[ch.In_room].People; vict != nil; vict = next_v {
					next_v = vict.Next_in_room
					if vict == ch {
						continue
					}
					if can_kill(ch, vict, nil, 1) == 0 {
						act(libc.CString("@CYou are protected from the water!@n"), TRUE, vict, nil, nil, TO_VICT)
						act(libc.CString("@C$n@C is protected from the water!@n"), TRUE, vict, nil, nil, TO_ROOM)
					} else if int(vict.Race) == RACE_KANASSAN {
						act(libc.CString("@CYou effortlessly swim against the current.@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@C effortlessly swims against the current.@n"), TRUE, vict, nil, nil, TO_ROOM)
					} else if int(vict.Skills[SKILL_BALANCE]) >= axion_dice(-10) {
						act(libc.CString("@CYou manage to keep your balance and are not swept away!@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@C manages to keep $s balance and is not swept away!@n"), TRUE, vict, nil, nil, TO_ROOM)
					} else if AFF_FLAGGED(ch, AFF_FLYING) {
						act(libc.CString("@CYou fly above the rushing waters and are untouched.@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@C flies above the rushing waters and is untouched.@n"), TRUE, vict, nil, nil, TO_ROOM)
					} else {
						act(libc.CString("@cYou are caught by the rushing waters and sent tumbling away!@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@c is caught by the rushing waters and sent tumbling away!@n"), TRUE, vict, nil, nil, TO_ROOM)
						do_simple_move(vict, attempt, TRUE)
						hurt(0, 0, ch, vict, nil, cost*4, 1)
					}
				}
				world[(world[ch.In_room].Dir_option[attempt]).To_room].Geffect = -3
				ch.Lastattack = last
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				ch.Con_cooldown = 15
			}
		} else {
			send_to_char(ch, libc.CString("You can not flood the water that direction!\r\n"))
			return
		}
	} else {
		send_to_char(ch, libc.CString("Syntax 1: hydromancy (flood) (direction)\r\n"))
		send_to_char(ch, libc.CString("Example: hydromancy flood nw\r\n"))
		send_to_char(ch, libc.CString("\nSyntax 2: hydromancy spike\r\n"))
		return
	}
}
func do_kanso(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if libc.StrCaseCmp(GET_NAME(ch), libc.CString("Levanthoth")) != 0 {
		send_to_char(ch, libc.CString("You do not know how to perform that technique. \r\n"))
		return
	}
	var vict *char_data
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Perform Kanso Suru on who?\r\n"))
		return
	}
	if int(vict.Race) == RACE_ANDROID {
		send_to_char(ch, libc.CString("Mechanical beings are not effected by this technique.\r\n"))
		return
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	var cost int64 = ch.Max_mana / int64(ch.Aff_abils.Intel)
	var dice int = axion_dice(-5)
	var skill int = int(ch.Aff_abils.Intel)
	var pdice int = axion_dice(0)
	var dam int = rand_number(1, 4)
	var af affected_type
	if int(ch.Aff_abils.Wis) > axion_dice(-5) {
		dam += 1
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	if skill < dice {
		act(libc.CString("You close your eyes and focus, before bounding effortlessly toward $N. Closing the distance, you place your hands on $N's chest but nothing happens!\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n closes $s eyes and bounds effortlessly towards you. Closing the distance, $e places $s hands on your chest but nothing happens!\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n closes $s eyes and bounds toward $N. Smirking, $e puts $m hands on $N's chest but nothing seems to happen.\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= cost
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		act(libc.CString("You close your eyes and focus, before effortlessly bounding towards $N. Closing the distance, you smirk at $N and place your hands on their chest. Electricity flows into their body as you draw water out of their very cells!\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n closes $s eyes and focuses, before effortlessly bounding toward you. Closing the distance, $e smirks and places both of $s hands on your chest. Electricity begins to pulse through your body and a great thirst takes hold, as if $n is drawing the water from your body!\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n closes $s eyes, before effortlessly bounding toward $N. Closing the distance, $n smirks and places both $s hands on $N's chest. Electricity seems to pass from $n's body to $N's!\r\n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= cost
		if int(vict.Player_specials.Conditions[THIRST])-dam >= 0 {
			vict.Player_specials.Conditions[THIRST] -= int8(dam)
		} else {
			vict.Player_specials.Conditions[THIRST] = 0
		}
		if int(ch.Player_specials.Conditions[THIRST])+dam <= 48 {
			ch.Player_specials.Conditions[THIRST] += int8(dam)
		} else {
			ch.Player_specials.Conditions[THIRST] = 48
		}
		if float64(ch.Hit)+(float64(gear_pl(ch))*0.01)*float64(dam) <= float64(gear_pl(ch)) {
			ch.Hit += int64((float64(gear_pl(ch)) * 0.01) * float64(dam))
		} else {
			ch.Hit = gear_pl(ch)
		}
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		if !AFF_FLAGGED(vict, AFF_HYDROZAP) {
			send_to_char(vict, libc.CString("@RYou feel less agile and your muscles ache!@n\r\n"))
			SET_BIT_AR(vict.Affected_by[:], AFF_HYDROZAP)
			vict.Real_abils.Dex -= 4
			vict.Real_abils.Con -= 4
			save_char(vict)
		}
		if skill > pdice && !AFF_FLAGGED(vict, AFF_PARA) {
			act(libc.CString("@R$N@W is paralyzed by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RYou are paralyzed by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$N@Wis paralyzed by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			af.Type = SKILL_PARALYZE
			af.Duration = int16(rand_number(1, 3))
			af.Modifier = 0
			af.Location = APPLY_NONE
			af.Bitvector = AFF_PARA
			affect_join(vict, &af, FALSE != 0, FALSE != 0, FALSE != 0, FALSE != 0)
		}
	}
}
func rpp_feature(ch *char_data, arg *byte) {
	var (
		cost   int = 0
		change int = FALSE
	)
	if *arg == 0 {
		send_to_char(ch, libc.CString("Syntax: rpp 13 (description)\nExample: rpp 13 a large red scar on his face\nDisplayed to others: He has a large red scar on his face.\r\n"))
		return
	}
	if libc.StrLen(arg) > 60 {
		send_to_char(ch, libc.CString("Please limit it to 60 characters.\r\n"))
		return
	}
	if ch.Feature == nil {
		cost = 2
	} else {
		cost = 2
		change = TRUE
	}
	if cost > ch.Rbank {
		send_to_char(ch, libc.CString("You do not have enough RPP in your Bank for that!\r\n"))
		return
	} else {
		var (
			sex  [128]byte
			buf8 [2048]byte
		)
		stdio.Sprintf(&sex[0], "%s", func() string {
			if int(ch.Sex) == SEX_FEMALE {
				return "She"
			}
			if int(ch.Sex) == SEX_MALE {
				return "He"
			}
			return "It"
		}())
		ch.Rbank -= cost
		stdio.Sprintf(&buf8[0], "...%s has %s.", &sex[0], arg)
		send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
		send_to_char(ch, libc.CString("You now have the following line underneath you when someone sees you in a room as:\n@C%s@n\r\n"), &buf8[0])
		ch.Feature = libc.StrDup(&buf8[0])
		if change == TRUE {
			send_to_imm(libc.CString("%s has altered their extra description. Make sure the reason is legit! If it is then reimb them 2 RPP.\r\n"), GET_USER(ch))
			send_to_char(ch, libc.CString("The immortals have been notified about this change. It had better have been for a good reason.\r\n"))
		}
		basic_mud_log(libc.CString("%s RPP Feature: '%s' Check for rule compliance."), GET_USER(ch), &buf8[0])
		return
	}
}
func do_instill(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var obj *obj_data
	var token *obj_data
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: instill (token) (target)\r\n"))
		return
	}
	if (func() *obj_data {
		token = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return token
	}()) == nil {
		send_to_char(ch, libc.CString("Syntax: instill (token) (target)\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("Syntax: instill (token) (target)\r\n"))
		return
	}
	if !OBJ_FLAGGED(token, ITEM_TOKEN) {
		send_to_char(ch, libc.CString("That is not a token.\r\n"))
		return
	}
	if OBJ_FLAGGED(token, ITEM_FORGED) {
		send_to_char(ch, libc.CString("That token is a forgery!\r\n"))
		return
	}
	if wearable_obj(obj) == 0 {
		send_to_char(ch, libc.CString("You can only instill tokens into equipment.\r\n"))
		return
	}
	if !OBJ_FLAGGED(obj, ITEM_SLOT1) && !OBJ_FLAGGED(obj, ITEM_SLOT2) {
		send_to_char(ch, libc.CString("That piece of equipment does not have any slots.\r\n"))
		return
	}
	if OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("That piece of equipment has already had its token slots filled. This can not be reversed."))
		return
	} else {
		var (
			stat  int = 0
			raise int = 0
		)
		stat = token.Affected[0].Location
		if obj.Affected[0].Location != 0 && obj.Affected[1].Location != 0 && obj.Affected[2].Location != 0 && obj.Affected[3].Location != 0 && obj.Affected[4].Location != 0 && obj.Affected[5].Location != 0 {
			if obj.Affected[0].Location != stat && obj.Affected[1].Location != stat && obj.Affected[2].Location != stat && obj.Affected[3].Location != stat && obj.Affected[4].Location != stat && obj.Affected[5].Location != stat {
				send_to_char(ch, libc.CString("This already has as many different stats as it can hold.\r\n"))
				return
			}
		}
		act(libc.CString("@GYou instill the token into @g$p@G. It glows @ggreen@G for a moment before returning to normal. The token disappears with the glow.@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@g$n@G instills a token into @g$p@G. It glows @ggreen@G for a moment before returning to normal. The token disappears with the glow.@n"), TRUE, ch, obj, nil, TO_ROOM)
		raise = token.Affected[0].Modifier
		extract_obj(token)
		if OBJ_FLAGGED(obj, ITEM_SLOT1) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOTS_FILLED)
		} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && !OBJ_FLAGGED(obj, ITEM_SLOT_ONE) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOT_ONE)
		} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && OBJ_FLAGGED(obj, ITEM_SLOT_ONE) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOTS_FILLED)
		}
		if obj.Affected[0].Location == stat {
			obj.Affected[0].Modifier += raise
		} else if obj.Affected[1].Location == stat {
			obj.Affected[1].Modifier += raise
		} else if obj.Affected[2].Location == stat {
			obj.Affected[2].Modifier += raise
		} else if obj.Affected[3].Location == stat {
			obj.Affected[3].Modifier += raise
		} else if obj.Affected[4].Location == stat {
			obj.Affected[4].Modifier += raise
		} else if obj.Affected[5].Location == stat {
			obj.Affected[5].Modifier += raise
		} else if obj.Affected[0].Location == 0 {
			obj.Affected[0].Location = stat
			obj.Affected[0].Modifier = raise
		} else if obj.Affected[1].Location == 0 {
			obj.Affected[1].Location = stat
			obj.Affected[1].Modifier = raise
		} else if obj.Affected[2].Location == 0 {
			obj.Affected[2].Location = stat
			obj.Affected[2].Modifier = raise
		} else if obj.Affected[3].Location == 0 {
			obj.Affected[3].Location = stat
			obj.Affected[3].Modifier = raise
		} else if obj.Affected[4].Location == 0 {
			obj.Affected[4].Location = stat
			obj.Affected[4].Modifier = raise
		} else if obj.Affected[5].Location == 0 {
			obj.Affected[5].Location = stat
			obj.Affected[5].Modifier = raise
		}
	}
}
func do_hayasa(ch *char_data, argument *byte, cmd int, subcmd int) {
	if !IS_NPC(ch) && GET_SKILL(ch, SKILL_HAYASA) == 0 {
		send_to_char(ch, libc.CString("You do not know how to perform this technique!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_HAYASA) {
		send_to_char(ch, libc.CString("You are already focusing ki to continually speed up your movements.\r\n"))
		return
	}
	var skill int = GET_SKILL(ch, SKILL_HAYASA)
	var prob int = axion_dice(0)
	var cost int64 = ch.Max_mana / int64(skill/2)
	var duration int = 1
	if skill >= 100 {
		duration = 6
	} else if skill >= 80 {
		duration = 5
	} else if skill >= 50 {
		duration = 4
	} else if skill >= 25 {
		duration = 3
	} else {
		duration = 2
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	} else if skill < prob {
		ch.Mana -= cost
		act(libc.CString("@CYou close your eyes for a brief moment and focus your ki around your body as a soft blue glow. The glow disappears though as you fail to maintain the effect...@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C closes $s eyes for a brief moment and a soft blue glow begins to form around $s body. The glow disappears a second later though and $e frowns.@n"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_HAYASA, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	} else {
		var af affected_type
		ch.Mana -= cost
		af.Type = SPELL_HAYASA
		af.Duration = int16(duration)
		af.Modifier = 0
		af.Location = APPLY_NONE
		af.Bitvector = AFF_HAYASA
		affect_join(ch, &af, FALSE != 0, FALSE != 0, FALSE != 0, FALSE != 0)
		ch.Speedboost = int(float64(GET_SPEEDCALC(ch)) * 0.5)
		reveal_hiding(ch, 0)
		act(libc.CString("@CYou close your eyes for a brief moment and focus your ki around your body as a soft blue glow. All your movements are faster now!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C closes $s eyes for a brief moment and a soft blue glow begins to form around $s body. The glow pulsates gently as $e opens his eyes and smiles.@n"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_HAYASA, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	}
}
func do_bury(ch *char_data, argument *byte, cmd int, subcmd int) {
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are busy grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil || ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are busy struggling with someone!\r\n"))
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: dig [bury (item) | uncover]\r\n"))
		return
	}
	if SECT(ch.In_room) != SECT_FIELD && SECT(ch.In_room) != SECT_HILLS && SECT(ch.In_room) != SECT_FOREST && SECT(ch.In_room) != SECT_DESERT && SECT(ch.In_room) != SECT_MOUNTAIN {
		send_to_char(ch, libc.CString("You are not in a room with enough available dirt or sand to dig.\r\n"))
		return
	}
	var obj *obj_data = nil
	var buried *obj_data = nil
	_ = buried
	var fobj *obj_data = find_obj_in_list_flag(world[ch.In_room].Contents, ITEM_BURIED)
	if libc.StrCaseCmp(&arg[0], libc.CString("bury")) == 0 {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Bury what?\r\n"))
			return
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have that object to bury.\r\n"))
			return
		} else if fobj != nil {
			send_to_char(ch, libc.CString("There is already something buried near here.\r\n"))
			return
		} else {
			if SECT(ch.In_room) != SECT_DESERT {
				act(libc.CString("@yYou start digging in a spot of soft dirt. Once you have an appropriately sized hole you drop @G$p@y in and then cover it.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@y starts digging in a spot of soft dirt. Once $e has an appropriately sized hole $e drops @G$p@y in and then covers it.@n"), TRUE, ch, obj, nil, TO_ROOM)
			} else {
				act(libc.CString("@YYou start digging in a spot of soft sand. Once you have an appropriately sized hole you drop @G$p@Y in and then cover it.@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@Y starts digging in a spot of soft sand. Once $e has an appropriately sized hole $e drops @G$p@Y in and then covers it.@n"), TRUE, ch, obj, nil, TO_ROOM)
			}
			obj_from_char(obj)
			obj_to_room(obj, ch.In_room)
			SET_BIT_AR(obj.Extra_flags[:], ITEM_BURIED)
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("uncover")) == 0 {
		if fobj == nil {
			send_to_char(ch, libc.CString("There is nothing buried here.\r\n"))
			return
		} else {
			if SECT(ch.In_room) != SECT_DESERT {
				act(libc.CString("@yYou slowly dig and reveal @G$p@y buried in the dirt! You pull it out and set it on the ground before covering the hole back up.@n"), TRUE, ch, fobj, nil, TO_CHAR)
				act(libc.CString("@C$n@y starts digging and shortly reveals @G$p@y buried in the dirt! Quickly $e pulls it out and sets it on the ground before covering the hole back up.@n"), TRUE, ch, fobj, nil, TO_ROOM)
			} else {
				act(libc.CString("@YYou slowly dig and reveal @G$p@Y buried in the sand! You pull it out and set it on the ground before covering the hole back up.@n"), TRUE, ch, fobj, nil, TO_CHAR)
				act(libc.CString("@C$n@Y starts digging and shortly reveals @G$p@Y buried in the sand! Quickly $e pulls it out and sets it on the ground before covering the hole back up.@n"), TRUE, ch, fobj, nil, TO_ROOM)
			}
			REMOVE_BIT_AR(fobj.Extra_flags[:], ITEM_BURIED)
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: dig [bury (item) | uncover]\r\n"))
		return
	}
}
func do_arena(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("You are too busy competing to be a spectator.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: arena (fighter number of participant)\r\n        arena look\r\n        arena scan\r\n        arena stop\r\n"))
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("stop")) == 0 {
		send_to_char(ch, libc.CString("You stop viewing what's going on in the arena.\r\n"))
		REMOVE_BIT_AR(ch.Player_specials.Pref[:], PRF_ARENAWATCH)
		ch.Arenawatch = -1
		return
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) != 0x45D3 {
		send_to_char(ch, libc.CString("You are not close enough to the arena floor to see it.\r\n"))
		return
	} else if libc.StrCaseCmp(&arg[0], libc.CString("look")) == 0 {
		if !PRF_FLAGGED(ch, PRF_ARENAWATCH) {
			send_to_char(ch, libc.CString("You are not even watching anyone in the arena.\r\n"))
			return
		} else if arena_watch(ch) != int(-1) {
			look_at_room(real_room(room_vnum(arena_watch(ch))), ch, 0)
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("scan")) == 0 {
		if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) == 0x45D3 {
			var (
				found int = FALSE
				d     *descriptor_data
			)
			send_to_char(ch, libc.CString("@D---@CFighters in the arena@D---@n\r\n"))
			for d = descriptor_list; d != nil; d = d.Next {
				if d.Connected != CON_PLAYING {
					continue
				}
				if IN_ARENA(d.Character) {
					var buf [2048]byte
					stdio.Sprintf(&buf[0], "@YFighter Number@D: @w%d, $N.@n", d.Character.Idnum)
					act(&buf[0], TRUE, ch, nil, unsafe.Pointer(d.Character), TO_CHAR)
					found = TRUE
				}
			}
			if found == FALSE {
				send_to_char(ch, libc.CString("@wNone.@n\r\n"))
			}
		} else {
			send_to_char(ch, libc.CString("You are not close enough to see what fighters are in the arena.\r\n"))
			return
		}
	} else {
		var num int = libc.Atoi(libc.GoString(&arg[0]))
		if num < 0 {
			send_to_char(ch, libc.CString("That is not a valid fighter number\r\n"))
			return
		} else {
			var (
				d     *descriptor_data
				found int = FALSE
			)
			for d = descriptor_list; d != nil; d = d.Next {
				if d.Connected != CON_PLAYING {
					continue
				}
				if int(d.Character.Idnum) == num {
					if IN_ARENA(d.Character) {
						found = TRUE
					}
				}
			}
			if found == TRUE {
				act(libc.CString("@wYou start watching the action surrounding that particular fighter in the arena.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@w starts watching the action in the arena.@n"), TRUE, ch, nil, nil, TO_ROOM)
				SET_BIT_AR(ch.Player_specials.Pref[:], PRF_ARENAWATCH)
				ch.Arenawatch = num
			} else {
				send_to_char(ch, libc.CString("A fighter with such a number was not found in the arena.\r\n"))
				return
			}
		}
	}
}
func is_good_silk(obj *obj_data) bool {
	return valid_silk(obj) != 0 && !OBJ_FLAGGED(obj, ITEM_FORGED)
}
func do_ensnare(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_ENSNARE) == 0 {
		return
	}
	var weave *obj_data = find_obj_in_list_lambda(ch.Carrying, is_good_silk)
	if weave == nil {
		send_to_char(ch, libc.CString("You do not have a bundle of silk to ensnare an opponent with!\r\n"))
		return
	} else {
		var (
			prob int = GET_SKILL(ch, SKILL_ENSNARE)
			perc int = axion_dice(0)
			arg  [2048]byte
			vict *char_data
		)
		one_argument(argument, &arg[0])
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: ensnare (target)\r\n"))
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Who are you trying to target with ensnare?\r\n"))
			return
		} else if AFF_FLAGGED(vict, AFF_ENSNARED) {
			send_to_char(ch, libc.CString("They are already ensnared!\r\n"))
			return
		} else if !HAS_ARMS(vict) {
			send_to_char(ch, libc.CString("They don't have arms to ensnare!\r\n"))
			return
		} else if prob <= perc {
			act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Unfortunately you miss and lose the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Fortunately $e misses and loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Fortunately $e misses and loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		} else if AFF_FLAGGED(vict, AFF_ZANZOKEN) && !AFF_FLAGGED(ch, AFF_ZANZOKEN) {
			act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Unfortunately @c$N@W zanzokens away avoiding it and you lose the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Fortunately you zanzoken away avoiding it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Fortunately @c$N@W zanzokens away avoiding it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			REMOVE_BIT_AR(vict.Affected_by[:], AFF_ZANZOKEN)
		} else if AFF_FLAGGED(vict, AFF_ZANZOKEN) && AFF_FLAGGED(ch, AFF_ZANZOKEN) {
			if GET_SPEEDI(ch)+rand_number(1, 100) < GET_SPEEDI(vict)+rand_number(1, 100) {
				act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! You both zanzoken! Unfortunately @c$N@W manages to avoid it and you lose the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! You both zanzoken! Fortunately you manage to avoid it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! They both zanzoken! Fortunately @c$N@W manages to avoid it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				REMOVE_BIT_AR(vict.Affected_by[:], AFF_ZANZOKEN)
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
			} else {
				act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Fortunately you manage to hit $M! You both zanzoken! Quickly you spin around $M and ensnare $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Unfortunately $e manages to hit YOU! You both zanzoken! Quickly $e spins around you and ensnares your arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Unfortunately $e manages to hit $M! They both zanzoken! Quickly $e spins around @c$N@W and ensnares $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				SET_BIT_AR(vict.Affected_by[:], AFF_ENSNARED)
				REMOVE_BIT_AR(vict.Affected_by[:], AFF_ZANZOKEN)
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
			}
		} else if AFF_FLAGGED(ch, AFF_ZANZOKEN) && !AFF_FLAGGED(vict, AFF_ZANZOKEN) {
			act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Fortunately you manage to hit $M! Quickly you zanzoken and spin around $M and ensnare $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Unfortunately $e manages to hit YOU! Quickly $e zanzokens and spins around you and ensnares your arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Unfortunately $e manages to hit $M! Quickly $e zanzokens and spins around @c$N@W and ensnares $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_ENSNARE, 0)
			REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
		} else if GET_SPEEDI(ch)+rand_number(1, 100) < GET_SPEEDI(vict)+rand_number(1, 100) {
			act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Unfortunately @c$N@W manages to avoid it and you lose the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Fortunately you manage to avoid it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Fortunately @c$N@W manages to avoid it and @C$n@W loses the bundle...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		} else {
			act(libc.CString("@WYou unwind your bundle of silk and grab a loose end of it. Splitting that end to reveal the sticky innards of the strand you swing the strand at @c$N@W! Fortunately you manage to hit $M! Quickly you spin around $M and ensnare $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at YOU! Unfortunately $e manages to hit YOU! Quickly $e spins around you and ensnares your arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W unwinds a bundle of silk and grabs a loose end of it. Splitting that end to reveal the sticky innards of the strand $e swings the strand at @c$N@W! Unfortunately $e manages to hit $M! Quickly $e spins around @c$N@W and ensnares $S arms with the silk!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			SET_BIT_AR(vict.Affected_by[:], AFF_ENSNARED)
		}
		extract_obj(weave)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		improve_skill(ch, SKILL_ENSNARE, 0)
	}
}
func valid_silk(obj *obj_data) int {
	var value int = 0
	switch GET_OBJ_VNUM(obj) {
	case 16700:
		fallthrough
	case 0x413D:
		fallthrough
	case 0x413E:
		fallthrough
	case 0x413F:
		fallthrough
	case 0x4140:
		fallthrough
	case 0x4144:
		value = 1
	}
	return value
}
func do_silk(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_SILK) == 0 {
		return
	}
	var obj *obj_data = nil
	var weave *obj_data = nil
	_ = weave
	var next_obj *obj_data = nil
	_ = next_obj
	var weaved *obj_data = nil
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	var prob int = GET_SKILL(ch, SKILL_SILK)
	var perc int = rand_number(1, 120)
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: silk (weave | bundle)\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("weave")) == 0 {
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: silk weave (head | wrist | belt)\r\n"))
			return
		}
		obj = find_obj_in_list_lambda(ch.Carrying, is_good_silk)
		var found int = FALSE
		_ = found
		var armor int = 500
		var str int = 0
		var intel int = 0
		var olevel int = 0
		var price float64 = 1
		if obj == nil {
			send_to_char(ch, libc.CString("You do not have an acceptable bundle of silk in your inventory!\r\n"))
			return
		} else {
			if libc.StrCaseCmp(&arg2[0], libc.CString("head")) == 0 {
				if prob <= perc {
					act(libc.CString("@WYou attempt to weave $p@W into the desired piece but end up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave $p@W into some type of clothing but ends up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_ROOM)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
					return
				} else {
					weaved = read_object(0x4141, VIRTUAL)
					obj_to_room(weaved, ch.In_room)
					if GET_OBJ_VNUM(obj) == 0x4144 {
						armor *= 20
						str = 4
						intel = 4
						price = 4
						olevel = 80
					} else if GET_OBJ_VNUM(obj) == 16700 {
						armor *= 15
						str = 3
						intel = 3
						price = 4
						olevel = 75
					} else if GET_OBJ_VNUM(obj) == 0x413D {
						armor *= 10
						str = 3
						intel = 3
						price = 3
						olevel = 50
					} else if GET_OBJ_VNUM(obj) == 0x413E {
						armor *= 6
						str = 2
						intel = 2
						price = 2
						olevel = 25
					} else if GET_OBJ_VNUM(obj) == 0x413F {
						armor *= 4
						str = 1
						intel = 1
						price = 1.5
						olevel = 5
					}
					weaved.Affected[0].Location = 17
					weaved.Affected[0].Modifier = armor
					weaved.Cost *= int(price)
					weaved.Value[0] = olevel
					weaved.Level = olevel
					if str > 0 {
						weaved.Affected[1].Location = 1
						weaved.Affected[1].Modifier = str
					}
					if intel > 0 {
						weaved.Affected[2].Location = 3
						weaved.Affected[2].Modifier = intel
					}
					act(libc.CString("@WYou attempt to weave the bundle and manage to create $p@W!@n"), TRUE, ch, weaved, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave a bundle into something and manages to create $p@W!@n"), TRUE, ch, weaved, nil, TO_ROOM)
					do_get(ch, libc.CString("headsash"), 0, 0)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				}
			} else if libc.StrCaseCmp(&arg2[0], libc.CString("wrist")) == 0 {
				if prob <= perc {
					act(libc.CString("@WYou attempt to weave $p@W into the desired piece but end up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave $p@W into some type of clothing but ends up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_ROOM)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
					return
				} else {
					weaved = read_object(0x4142, VIRTUAL)
					obj_to_room(weaved, ch.In_room)
					if GET_OBJ_VNUM(obj) == 0x4144 {
						armor *= 20
						str = 4
						intel = 4
						price = 4
						olevel = 80
					} else if GET_OBJ_VNUM(obj) == 16700 {
						armor *= 15
						str = 3
						intel = 3
						price = 4
						olevel = 75
					} else if GET_OBJ_VNUM(obj) == 0x413D {
						armor *= 10
						str = 3
						intel = 3
						price = 3
						olevel = 50
					} else if GET_OBJ_VNUM(obj) == 0x413E {
						armor *= 6
						str = 2
						intel = 2
						price = 2
						olevel = 25
					} else if GET_OBJ_VNUM(obj) == 0x413F {
						armor *= 4
						str = 1
						intel = 1
						price = 1.5
						olevel = 5
					}
					weaved.Affected[0].Location = 17
					weaved.Affected[0].Modifier = armor
					weaved.Cost *= int(price)
					weaved.Value[0] = olevel
					weaved.Level = olevel
					if str > 0 {
						weaved.Affected[1].Location = 1
						weaved.Affected[1].Modifier = str
					}
					if intel > 0 {
						weaved.Affected[2].Location = 3
						weaved.Affected[2].Modifier = intel
					}
					act(libc.CString("@WYou attempt to weave the bundle and manage to create $p@W!@n"), TRUE, ch, weaved, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave a bundle into something and manages to create $p@W!@n"), TRUE, ch, weaved, nil, TO_ROOM)
					do_get(ch, libc.CString("wristband"), 0, 0)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				}
			} else if libc.StrCaseCmp(&arg2[0], libc.CString("belt")) == 0 {
				if prob <= perc {
					act(libc.CString("@WYou attempt to weave $p@W into the desired piece but end up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave $p@W into some type of clothing but ends up ruining the entire bundle instead.@n"), TRUE, ch, obj, nil, TO_ROOM)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
					return
				} else {
					weaved = read_object(0x4143, VIRTUAL)
					obj_to_room(weaved, ch.In_room)
					if GET_OBJ_VNUM(obj) == 0x4144 {
						armor *= 20
						str = 4
						intel = 4
						price = 4
						olevel = 80
					} else if GET_OBJ_VNUM(obj) == 16700 {
						armor *= 15
						str = 3
						intel = 3
						price = 4
						olevel = 75
					} else if GET_OBJ_VNUM(obj) == 0x413D {
						armor *= 10
						str = 3
						intel = 3
						price = 3
						olevel = 50
					} else if GET_OBJ_VNUM(obj) == 0x413E {
						armor *= 6
						str = 2
						intel = 2
						price = 2
						olevel = 25
					} else if GET_OBJ_VNUM(obj) == 0x413F {
						armor *= 4
						str = 1
						intel = 1
						price = 1.5
						olevel = 5
					}
					weaved.Affected[0].Location = 17
					weaved.Affected[0].Modifier = armor
					weaved.Cost *= int(price)
					weaved.Value[0] = olevel
					weaved.Level = olevel
					if str > 0 {
						weaved.Affected[1].Location = 1
						weaved.Affected[1].Modifier = str
					}
					if intel > 0 {
						weaved.Affected[2].Location = 3
						weaved.Affected[2].Modifier = intel
					}
					act(libc.CString("@WYou attempt to weave the bundle and manage to create $p@W!@n"), TRUE, ch, weaved, nil, TO_CHAR)
					act(libc.CString("@C$n@W attempts to weave a bundle into something and manages to create $p@W!@n"), TRUE, ch, weaved, nil, TO_ROOM)
					do_get(ch, libc.CString("belt"), 0, 0)
					extract_obj(obj)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				}
			} else {
				send_to_char(ch, libc.CString("Syntax: silk weave (head | wrist | belt)"))
				return
			}
			return
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("bundle")) == 0 {
		var cost int64 = int64(((float64(ch.Max_mana) * 0.01) * (float64(prob) * 0.2)) + float64(int(ch.Aff_abils.Intel)*GET_LEVEL(ch)))
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to weave any bundles of silk.\r\n"))
			return
		} else {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			var super int = FALSE
			var superoll int = rand_number(1, 100)
			if int(ch.Chclass) == CLASS_KURZAK {
				if GET_SKILL(ch, SKILL_SILK) >= 100 {
					if 8 > superoll {
						super = TRUE
					}
				} else if GET_SKILL(ch, SKILL_SILK) >= 60 {
					if 6 > superoll {
						super = TRUE
					}
				} else if GET_SKILL(ch, SKILL_SILK) >= 40 {
					if 3 > superoll {
						super = TRUE
					}
				}
			}
			if super == TRUE {
				obj = read_object(0x4144, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@YYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@Y piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				send_to_char(ch, libc.CString("@YIt's SUPER grand!@n\r\n"))
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else if prob > perc && prob >= 100 {
				obj = read_object(16700, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@W piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else if prob > perc && prob >= 90 {
				obj = read_object(0x413D, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@W piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else if prob > perc && prob >= 80 {
				obj = read_object(0x413E, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@W piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else if prob > perc && prob >= 50 {
				obj = read_object(0x413F, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@W piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else if prob > perc {
				obj = read_object(0x4140, VIRTUAL)
				obj_to_room(obj, ch.In_room)
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You gently weave the silk and in no time at all you have a $p@W piled at your feet!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a golden colored silk from $s mouth. Gently $e weaves the silk and in no time at all $e has a $p@W piled at $s feet!@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
			} else {
				act(libc.CString("@WYou concentrate your ki into your silk sacs and begin to spit silk out of your mouth. You end up making a poorly formed puddle of goo...@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W seems to concentrate for a moment before spitting out a poorly formed puddle of goo...@n"), TRUE, ch, obj, nil, TO_ROOM)
				ch.Mana -= cost
				improve_skill(ch, SKILL_SILK, 1)
			}
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: silk (weave | bundle)\r\n"))
		return
	}
}
func do_adrenaline(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Race) != RACE_ARLIAN && (int(ch.Race) != RACE_BIO || int(ch.Race) == RACE_BIO && ((ch.Genome[0]) != 6 && (ch.Genome[1]) != 6)) {
		send_to_char(ch, libc.CString("You are not an arlian and do not possess this ability\r\n"))
		return
	} else {
		var (
			arg  [2048]byte
			arg2 [2048]byte
		)
		two_arguments(argument, &arg[0], &arg2[0])
		if arg[0] == 0 || arg2[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: adrenaline (pl or ki) (percent)\r\nExample: adrenaline pl 10\r\n"))
			return
		} else {
			if libc.Atoi(libc.GoString(&arg2[0])) < 0 || libc.Atoi(libc.GoString(&arg2[0])) > 100 {
				send_to_char(ch, libc.CString("The percent must be between 1 and 100%s.\r\n"), "%")
				return
			}
			var percent float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
			if float64(ch.Move)-float64(gear_pl(ch))*percent < 0 {
				send_to_char(ch, libc.CString("You do not have enough stamina to trade for adrenaline!\r\n"))
				return
			}
			var trade int64 = int64(float64(gear_pl(ch)) * percent)
			if libc.StrCaseCmp(&arg[0], libc.CString("pl")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("PL")) == 0 {
				act(libc.CString("@GYou focus your mind and begin to overwork your powerful adrenal glands and your wounds begin to heal!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@g$n@G seems to concentrate and $s wounds begin to heal!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Hit+trade > gear_pl(ch) {
					send_to_char(ch, libc.CString("Some of your stamina was wasted because your powerlevel maxed out.\r\n"))
					ch.Hit = gear_pl(ch)
					ch.Move -= trade
				} else {
					ch.Hit += trade
					ch.Move -= trade
				}
			} else if libc.StrCaseCmp(&arg[0], libc.CString("ki")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("KI")) == 0 {
				act(libc.CString("@GYou focus your mind and begin to overwork your powerful adrenal glands and you feel your ki replenish!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@g$n@G seems to concentrate and $e appears energized!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Mana+trade > ch.Max_mana {
					send_to_char(ch, libc.CString("Some of your stamina was wasted because your ki maxed out.\r\n"))
					ch.Mana = ch.Max_mana
					ch.Move -= trade
				} else {
					ch.Mana += trade
					ch.Move -= trade
				}
			}
		}
	}
}
func disp_rpp_store(ch *char_data) {
	send_to_char(ch, libc.CString("@m                        RPP Item Store@n\n"))
	send_to_char(ch, libc.CString("@D~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\n"))
	send_to_char(ch, libc.CString("@GItem Name                      @gRPP Cost        @cChoice Number   @yMin Lvl@n\n"))
	send_to_char(ch, libc.CString("@WStardust Equipment Set         @D[@Y20@D]            @D[ @C1@D]            @w50@n\n"))
	send_to_char(ch, libc.CString("@WPlatinum Masamune (Sword Skill)@D[@Y 5@D]            @D[ @C2@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WObsidian Dirk (Dagger Skill)   @D[@Y 5@D]            @D[ @C3@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WEmerald Javelin (Spear Skill)  @D[@Y 5@D]            @D[ @C4@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WIvory Cane (Club Skill)        @D[@Y 5@D]            @D[ @C5@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WHyper X65 Cannon (Gun Skill)   @D[@Y 5@D]            @D[ @C6@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WJagged Rock (Brawl skill)      @D[@Y 5@D]            @D[ @C7@D]            @w40@n\n"))
	send_to_char(ch, libc.CString("@WKachin Mountain                @D[@Y 8@D]            @D[ @C8@D]@n\n"))
	send_to_char(ch, libc.CString("@WSpar Booster                   @D[@Y15@D]            @D[ @C9@D]@n\n"))
	send_to_char(ch, libc.CString("@D~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~@n\n"))
	send_to_char(ch, libc.CString("@wSyntax: rpp 12 (choice number)@n\r\n"))
}
func handle_rpp_store(ch *char_data, choice int) {
	var (
		obj    *obj_data
		objnum int = 0
		cost   int = 0
	)
	switch choice {
	case 1:
		cost = 20
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 7:
		cost = 5
	case 8:
		cost = 8
	case 9:
		cost = 15
	default:
		send_to_char(ch, libc.CString("That is not a selection option!\r\n"))
		return
	}
	if ch.Rp < cost {
		send_to_char(ch, libc.CString("You do not have enough RPP to afford that option.\r\n"))
		return
	} else {
		switch choice {
		case 1:
			if ch.Carry_weight+26 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+13 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 50 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				for objnum = 1110; objnum < 1120; objnum++ {
					if objnum <= 1116 {
						obj = read_object(obj_vnum(objnum), VIRTUAL)
						obj_to_char(obj, ch)
						obj.Size = get_size(ch)
						obj = nil
					} else {
						obj = read_object(obj_vnum(objnum), VIRTUAL)
						obj_to_char(obj, ch)
						obj.Size = get_size(ch)
						obj = nil
						obj = read_object(obj_vnum(objnum), VIRTUAL)
						obj_to_char(obj, ch)
						obj.Size = get_size(ch)
					}
				}
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 2:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1120, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP from your Bank paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 3:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1121, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 4:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1122, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 5:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("@R%d@W RPP from your Bank paid for your selection. Enjoy!@n\r\n"), cost)
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1123, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 6:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1124, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 7:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else if GET_LEVEL(ch) < 40 {
				send_to_char(ch, libc.CString("You are below the minimum level to equip it.\r\n"))
			} else {
				obj = read_object(1125, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 8:
			if ch.Carry_weight+10000000 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else {
				obj = read_object(1126, VIRTUAL)
				obj.Weight = 10000000
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		case 9:
			if ch.Carry_weight+2 > int(max_carry_weight(ch)) {
				send_to_char(ch, libc.CString("You can not carry that much weight at this moment.\r\n"))
			} else if int(ch.Carry_items)+1 > 50 {
				send_to_char(ch, libc.CString("You have too many items on you to carry anymore at this moment.\r\n"))
			} else {
				obj = read_object(1127, VIRTUAL)
				obj_to_char(obj, ch)
				obj.Size = get_size(ch)
				ch.Rp -= cost
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), cost)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), cost)
			}
		}
	}
}
func valid_recipe(ch *char_data, recipe int, type_ int) int {
	var (
		tomato    int = -1
		cucumber  int = -1
		onion     int = -1
		greenbean int = -1
		garlic    int = -1
		redpep    int = -1
		potato    int = -1
		carrot    int = -1
		brownmush int = -1
		lettuce   int = -1
		normmeat  int = -1
		goodmeat  int = -1
		normfish  int = -1
		goodfish  int = -1
		greatfish int = -1
		bestfish  int = -1
		rice      int = -1
		flour     int = -1
		appleplum int = -1
		fberry    int = -1
		carambola int = -1
		obj2      *obj_data
		next_obj  *obj_data
	)
	_ = next_obj
	var pass int = FALSE
	switch recipe {
	case RECIPE_TOMATO_SOUP:
		tomato = 2
	case RECIPE_STEAK:
		normmeat = 1
	case RECIPE_POTATO_SOUP:
		potato = 2
	case RECIPE_VEGETABLE_SOUP:
		potato = 1
		tomato = 1
		carrot = 1
		greenbean = 1
		onion = 1
	case RECIPE_MEAT_STEW:
		normmeat = 1
		potato = 1
		tomato = 1
		garlic = 1
	case RECIPE_ROAST:
		normmeat = 1
		potato = 2
		garlic = 1
		onion = 1
		greenbean = 3
	case RECIPE_CHILI_SOUP:
		normmeat = 1
		redpep = 4
		tomato = 2
	case RECIPE_GRILLED_NORMFISH:
		normfish = 1
	case RECIPE_GRILLED_GOODFISH:
		goodfish = 1
	case RECIPE_GRILLED_GREATFISH:
		greatfish = 1
	case RECIPE_GRILLED_BESTFISH:
		bestfish = 1
	case RECIPE_COOKED_RICE:
		rice = 1
	case RECIPE_SUSHI:
		rice = 1
		normfish = 1
	case RECIPE_BREAD:
		flour = 1
	case RECIPE_SALAD:
		tomato = 1
		cucumber = 1
		carrot = 1
		lettuce = 1
	case RECIPE_APPLEPLUM:
		flour = 1
		appleplum = 1
	case RECIPE_FBERRY_MUFFIN:
		flour = 1
		fberry = 1
	case RECIPE_CARAMBOLA_BREAD:
		flour = 1
		carambola = 1
	}
	if type_ == 0 {
		for obj2 = ch.Carrying; obj2 != nil; obj2 = obj2.Next_content {
			switch GET_OBJ_VNUM(obj2) {
			case RCP_TOMATO:
				if tomato > 0 {
					tomato -= 1
				}
			case RCP_NORMAL_MEAT:
				if normmeat > 0 {
					normmeat -= 1
				}
			case RCP_POTATO:
				if potato > 0 {
					potato -= 1
				}
			case RCP_ONION:
				if onion > 0 {
					onion -= 1
				}
			case RCP_CUCUMBER:
				if cucumber > 0 {
					cucumber -= 1
				}
			case RCP_CHILIPEPPER:
				if redpep > 0 {
					redpep -= 1
				}
			case RCP_CARROT:
				if carrot > 0 {
					carrot -= 1
				}
			case RCP_GREENBEAN:
				if greenbean > 0 {
					greenbean -= 1
				}
			case RCP_BLACKBASS:
				fallthrough
			case RCP_FLOUNDER:
				fallthrough
			case RCP_NARRI:
				fallthrough
			case RCP_GRAVELREBOI:
				if normfish > 0 {
					normfish -= 1
				}
			case RCP_SILVERTROUT:
				fallthrough
			case RCP_SILVEREEL:
				fallthrough
			case RCP_VALBISH:
				fallthrough
			case RCP_VOOSPIKE:
				if goodfish > 0 {
					goodfish -= 1
				}
			case RCP_STRIPEDBASS:
				fallthrough
			case RCP_COBIA:
				fallthrough
			case RCP_GUSBLAT:
				fallthrough
			case RCP_SHADOWFISH:
				if greatfish > 0 {
					greatfish -= 1
				}
			case RCP_BLUECATFISH:
				fallthrough
			case RCP_TAMBOR:
				fallthrough
			case RCP_REPEEIL:
				fallthrough
			case RCP_SHADEEEL:
				if bestfish > 0 {
					bestfish -= 1
				}
			case RCP_BROWNMUSH:
				if brownmush > 0 {
					brownmush -= 1
				}
			case RCP_GARLIC:
				if garlic > 0 {
					garlic -= 1
				}
			case RCP_RICE:
				if rice > 0 {
					rice -= 1
				}
			case RCP_FLOUR:
				if flour > 0 {
					flour -= 1
				}
			case RCP_LETTUCE:
				if lettuce > 0 {
					lettuce -= 1
				}
			case RCP_APPLEPLUM:
				if appleplum > 0 {
					appleplum -= 1
				}
			case RCP_FROZENBERRY:
				if fberry > 0 {
					fberry -= 1
				}
			case RCP_CARAMBOLA:
				if carambola > 0 {
					carambola -= 1
				}
			}
		}
	} else {
		for obj2 = ch.Carrying; obj2 != nil; obj2 = obj2.Next_content {
			switch GET_OBJ_VNUM(obj2) {
			case RCP_TOMATO:
				if tomato > 0 {
					tomato -= 1
					extract_obj(obj2)
				}
			case RCP_NORMAL_MEAT:
				if normmeat > 0 {
					normmeat -= 1
					extract_obj(obj2)
				}
			case RCP_POTATO:
				if potato > 0 {
					potato -= 1
					extract_obj(obj2)
				}
			case RCP_ONION:
				if onion > 0 {
					onion -= 1
					extract_obj(obj2)
				}
			case RCP_CUCUMBER:
				if cucumber > 0 {
					cucumber -= 1
					extract_obj(obj2)
				}
			case RCP_CHILIPEPPER:
				if redpep > 0 {
					redpep -= 1
					extract_obj(obj2)
				}
			case RCP_CARROT:
				if carrot > 0 {
					carrot -= 1
					extract_obj(obj2)
				}
			case RCP_GREENBEAN:
				if greenbean > 0 {
					greenbean -= 1
					extract_obj(obj2)
				}
			case RCP_BLACKBASS:
				fallthrough
			case RCP_FLOUNDER:
				fallthrough
			case RCP_NARRI:
				fallthrough
			case RCP_GRAVELREBOI:
				if normfish > 0 {
					normfish -= 1
					extract_obj(obj2)
				}
			case RCP_SILVERTROUT:
				fallthrough
			case RCP_SILVEREEL:
				fallthrough
			case RCP_VALBISH:
				fallthrough
			case RCP_VOOSPIKE:
				if goodfish > 0 {
					goodfish -= 1
					extract_obj(obj2)
				}
			case RCP_STRIPEDBASS:
				fallthrough
			case RCP_COBIA:
				fallthrough
			case RCP_GUSBLAT:
				fallthrough
			case RCP_SHADOWFISH:
				if greatfish > 0 {
					greatfish -= 1
					extract_obj(obj2)
				}
			case RCP_BLUECATFISH:
				fallthrough
			case RCP_TAMBOR:
				fallthrough
			case RCP_REPEEIL:
				fallthrough
			case RCP_SHADEEEL:
				if bestfish > 0 {
					bestfish -= 1
					extract_obj(obj2)
				}
			case RCP_BROWNMUSH:
				if brownmush > 0 {
					brownmush -= 1
					extract_obj(obj2)
				}
			case RCP_GARLIC:
				if garlic > 0 {
					garlic -= 1
					extract_obj(obj2)
				}
			case RCP_RICE:
				if rice > 0 {
					rice -= 1
					extract_obj(obj2)
				}
			case RCP_FLOUR:
				if flour > 0 {
					flour -= 1
					extract_obj(obj2)
				}
			case RCP_LETTUCE:
				if lettuce > 0 {
					lettuce -= 1
					extract_obj(obj2)
				}
			case RCP_APPLEPLUM:
				if appleplum > 0 {
					appleplum -= 1
					extract_obj(obj2)
				}
			case RCP_FROZENBERRY:
				if fberry > 0 {
					fberry -= 1
					extract_obj(obj2)
				}
			case RCP_CARAMBOLA:
				if carambola > 0 {
					carambola -= 1
					extract_obj(obj2)
				}
			}
		}
		return TRUE
	}
	switch recipe {
	case RECIPE_TOMATO_SOUP:
		if tomato == 0 {
			pass = TRUE
		}
	case RECIPE_STEAK:
		if normmeat == 0 {
			pass = TRUE
		}
	case RECIPE_POTATO_SOUP:
		if potato == 0 {
			pass = TRUE
		}
	case RECIPE_VEGETABLE_SOUP:
		if potato == 0 && tomato == 0 && carrot == 0 && greenbean == 0 && onion == 0 {
			pass = TRUE
		}
	case RECIPE_MEAT_STEW:
		if normmeat == 0 && potato == 0 && tomato == 0 && garlic == 0 {
			pass = TRUE
		}
	case RECIPE_ROAST:
		if normmeat == 0 && potato == 0 && garlic == 0 && onion == 0 && greenbean == 0 {
			pass = TRUE
		}
	case RECIPE_CHILI_SOUP:
		if normmeat == 0 && redpep == 0 && tomato == 0 {
			pass = TRUE
		}
	case RECIPE_GRILLED_NORMFISH:
		if normfish == 0 {
			pass = TRUE
		}
	case RECIPE_GRILLED_GOODFISH:
		if goodfish == 0 {
			pass = TRUE
		}
	case RECIPE_GRILLED_GREATFISH:
		if greatfish == 0 {
			pass = TRUE
		}
	case RECIPE_GRILLED_BESTFISH:
		if bestfish == 0 {
			pass = TRUE
		}
	case RECIPE_COOKED_RICE:
		if rice == 0 {
			pass = TRUE
		}
	case RECIPE_SUSHI:
		if rice == 0 && normfish == 0 {
			pass = TRUE
		}
	case RECIPE_BREAD:
		if flour == 0 {
			pass = TRUE
		}
	case RECIPE_SALAD:
		if tomato == 0 && cucumber == 0 && carrot == 0 && lettuce == 0 {
			pass = TRUE
		}
	case RECIPE_APPLEPLUM:
		if flour == 0 && appleplum == 0 {
			pass = TRUE
		}
	case RECIPE_FBERRY_MUFFIN:
		if flour == 0 && fberry == 0 {
			pass = TRUE
		}
	case RECIPE_CARAMBOLA_BREAD:
		if flour == 0 && carambola == 0 {
			pass = TRUE
		}
	}
	if pass == FALSE {
		if tomato > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W tomato%s for this recipe.@n\r\n"), tomato, func() string {
				if tomato > 1 {
					return "es"
				}
				return ""
			}())
		}
		if potato > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W potato%s for this recipe.@n\r\n"), potato, func() string {
				if potato > 1 {
					return "es"
				}
				return ""
			}())
		}
		if onion > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W onion%s for this recipe.@n\r\n"), onion, func() string {
				if onion > 1 {
					return "s"
				}
				return ""
			}())
		}
		if appleplum > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W appleplum%s for this recipe.@n\r\n"), appleplum, func() string {
				if appleplum > 1 {
					return "s"
				}
				return ""
			}())
		}
		if fberry > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W frozen berry%s for this recipe.@n\r\n"), fberry, func() string {
				if fberry > 1 {
					return "s"
				}
				return ""
			}())
		}
		if carambola > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W carambola%s for this recipe.@n\r\n"), carambola, func() string {
				if carambola > 1 {
					return "s"
				}
				return ""
			}())
		}
		if lettuce > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W head%s of lettuce for this recipe.@n\r\n"), lettuce, func() string {
				if lettuce > 1 {
					return "s"
				}
				return ""
			}())
		}
		if flour > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W cup%s of white flour for this recipe.@n\r\n"), flour, func() string {
				if flour > 1 {
					return "s"
				}
				return ""
			}())
		}
		if rice > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W cup%s of white rice for this recipe.@n\r\n"), rice, func() string {
				if rice > 1 {
					return "s"
				}
				return ""
			}())
		}
		if garlic > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W garlic clove%s for this recipe.@n\r\n"), garlic, func() string {
				if garlic > 1 {
					return "s"
				}
				return ""
			}())
		}
		if carrot > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W carrot%s for this recipe.@n\r\n"), carrot, func() string {
				if carrot > 1 {
					return "s"
				}
				return ""
			}())
		}
		if cucumber > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W cucumber%s for this recipe.@n\r\n"), cucumber, func() string {
				if cucumber > 1 {
					return "s"
				}
				return ""
			}())
		}
		if greenbean > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W green bean%s for this recipe.@n\r\n"), greenbean, func() string {
				if greenbean > 1 {
					return "s"
				}
				return ""
			}())
		}
		if normmeat > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W normal raw steak%s for this recipe.@n\r\n"), normmeat, func() string {
				if normmeat > 1 {
					return "s"
				}
				return ""
			}())
		}
		if goodmeat > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W good raw steak%s for this recipe.@n\r\n"), goodmeat, func() string {
				if goodmeat > 1 {
					return "s"
				}
				return ""
			}())
		}
		if redpep > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W chili pepper%s for this recipe.@n\r\n"), redpep, func() string {
				if redpep > 1 {
					return "s"
				}
				return ""
			}())
		}
		if normfish > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W black bass, flounder, narri, or gravel reboi for this recipe.@n\r\n"), normfish)
		}
		if goodfish > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W silver trout, silver eel, valbish, or voos pike for this recipe.@n\r\n"), goodfish)
		}
		if greatfish > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W striped bass, cobia, gusblat, or shadowfish for this recipe.@n\r\n"), greatfish)
		}
		if bestfish > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W blue catfish, tambor, repeeil, or shadeeel for this recipe.@n\r\n"), bestfish)
		}
		if brownmush > 0 {
			send_to_char(ch, libc.CString("@WYou need @m%d@W brown mushroom%s for this recipe.@n\r\n"), brownmush, func() string {
				if brownmush > 1 {
					return "s"
				}
				return ""
			}())
		}
		return FALSE
	} else {
		return TRUE
	}
}
func campfire_cook(recipe int) int {
	switch recipe {
	case RECIPE_STEAK:
		fallthrough
	case RECIPE_GRILLED_NORMFISH:
		fallthrough
	case RECIPE_GRILLED_GOODFISH:
		fallthrough
	case RECIPE_GRILLED_GREATFISH:
		fallthrough
	case RECIPE_GRILLED_BESTFISH:
		fallthrough
	case RECIPE_ROAST:
		return TRUE
	}
	return FALSE
}
func do_cook(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if cook_element(ch.In_room) == 0 {
		send_to_char(ch, libc.CString("You need a campfire or Flambus Stove nearby to cook.\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_COOKING) == 0 {
		send_to_char(ch, libc.CString("You don't even know the basics!\r\n"))
		return
	}
	var skill int = GET_SKILL(ch, SKILL_COOKING)
	var prob int = axion_dice(0)
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@D---------------------@RCooking@D---------------------@n\r\n"))
		send_to_char(ch, libc.CString("@Y 1@B) @CCooked Steak\t\t@Y17@B) @CCarambola Bread@n\r\n"))
		send_to_char(ch, libc.CString("@Y 2@B) @CTomato Soup\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 3@B) @CPotato Soup\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 4@B) @CVegetable Soup\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 5@B) @CMeat Stew\t\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 6@B) @CChili Soup\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 7@B) @CGrilled Fish\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 8@B) @CGood Grilled Fish\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y 9@B) @CGreat Grilled Fish\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y10@B) @CMagnificent Grilled Fish\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y11@B) @CCooked White Rice\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y12@B) @CSushi\t\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y13@B) @CWhite Bread\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y14@B) @CBasic Salad\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y15@B) @CAppleplum Chasan\t\t@n\r\n"))
		send_to_char(ch, libc.CString("@Y16@B) @CFrozen Berry Muffin\t@n\r\n"))
		send_to_char(ch, libc.CString("@wSyntax: cook (recipe number)@n\r\n"))
		return
	} else {
		var (
			num    int       = libc.Atoi(libc.GoString(&arg[0]))
			pass   int       = FALSE
			meal   *obj_data = nil
			recipe int       = -1
		)
		switch num {
		case 1:
			recipe = RECIPE_STEAK
			prob += 8
		case 2:
			recipe = RECIPE_TOMATO_SOUP
		case 3:
			recipe = RECIPE_POTATO_SOUP
		case 4:
			recipe = RECIPE_VEGETABLE_SOUP
		case 5:
			recipe = RECIPE_MEAT_STEW
		case 6:
			recipe = RECIPE_CHILI_SOUP
		case 7:
			recipe = RECIPE_GRILLED_NORMFISH
			prob += 6
		case 8:
			recipe = RECIPE_GRILLED_GOODFISH
			prob += 10
		case 9:
			recipe = RECIPE_GRILLED_GREATFISH
			prob += 12
		case 10:
			recipe = RECIPE_GRILLED_BESTFISH
			prob += 16
		case 11:
			recipe = RECIPE_COOKED_RICE
		case 12:
			recipe = RECIPE_SUSHI
		case 13:
			recipe = RECIPE_BREAD
		case 14:
			recipe = RECIPE_SALAD
		case 15:
			recipe = RECIPE_APPLEPLUM
		case 16:
			recipe = RECIPE_FBERRY_MUFFIN
		case 17:
			recipe = RECIPE_CARAMBOLA_BREAD
		}
		if recipe == -1 {
			send_to_char(ch, libc.CString("That is not a valid dish!\r\n"))
			return
		}
		if valid_recipe(ch, recipe, 0) == 0 {
			return
		} else if cook_element(ch.In_room) == 1 && campfire_cook(recipe) == 0 {
			send_to_char(ch, libc.CString("You can not cook that dish over a campfire.\r\n"))
			return
		} else {
			valid_recipe(ch, recipe, 1)
			pass = TRUE
		}
		if pass == TRUE {
			if skill < prob {
				act(libc.CString("@wYou screw up the preparation of the recipe and end up wasting the ingredients!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@w starts to prepare some food, but ends up ruining the ingredients instead!@n"), TRUE, ch, nil, nil, TO_ROOM)
				improve_skill(ch, SKILL_COOKING, 0)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			}
			var psbonus int = 0
			var expbonus int = 0
			switch num {
			case 1:
				meal = read_object(MEAL_STEAK, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 5
			case 2:
				meal = read_object(MEAL_TOMATO_SOUP, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 2
				expbonus = 15
			case 3:
				meal = read_object(MEAL_POTATO_SOUP, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 20
			case 4:
				meal = read_object(MEAL_VEGETABLE_SOUP, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 3
				expbonus = 45
			case 5:
				meal = read_object(MEAL_MEAT_STEW, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 2
				expbonus = 50
			case 6:
				meal = read_object(MEAL_CHILI_SOUP, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 0
				expbonus = 100
			case 7:
				meal = read_object(MEAL_NORM_FISH, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 2
				expbonus = 12
			case 8:
				meal = read_object(MEAL_GOOD_FISH, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 3
				expbonus = 40
			case 9:
				meal = read_object(MEAL_GREAT_FISH, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 5
				expbonus = 80
			case 10:
				meal = read_object(MEAL_BEST_FISH, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 7
				expbonus = 125
			case 11:
				meal = read_object(MEAL_COOKED_RICE, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 8
			case 12:
				meal = read_object(MEAL_SUSHI, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 2
				expbonus = 20
			case 13:
				meal = read_object(MEAL_BREAD, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 8
			case 14:
				meal = read_object(MEAL_SALAD, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 5
				expbonus = 8
			case 15:
				meal = read_object(MEAL_APPLEPLUM, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 9
			case 16:
				meal = read_object(MEAL_FBERRY_MUFFIN, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 3
				expbonus = 12
			case 17:
				meal = read_object(MEAL_CARAMBOLA_BREAD, VIRTUAL)
				obj_to_char(meal, ch)
				psbonus = 1
				expbonus = 9
			default:
				send_to_char(ch, libc.CString("That is not a valid dish!\r\n"))
				return
			}
			if (ch.Bonuses[BONUS_RECIPE]) != 0 {
				psbonus += 1
				expbonus += 3
			}
			act(libc.CString("@wYou carefully prepare the ingredients and then start cooking them. After a while of patience  and skillful care you successfully make @D'@C$p@D'@w!@n"), TRUE, ch, meal, nil, TO_CHAR)
			act(libc.CString("@C$n@w carefully prepares some ingredients and starts cooking them. After a while of patience and skillful care $e succeeds in making @D'@C$p@D'@w!@n"), TRUE, ch, meal, nil, TO_ROOM)
			improve_skill(ch, SKILL_COOKING, 0)
			if psbonus > 0 {
				if float64(skill)*0.1 > 0 {
					psbonus = int((float64(skill) * 0.1) * float64(psbonus))
				}
			}
			if expbonus > 0 {
				expbonus = skill * expbonus
			}
			meal.Value[1] = psbonus
			meal.Value[2] = expbonus
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
	}
}
func do_fireshield(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_FIRESHIELD) == 0 {
		return
	}
	if AFF_FLAGGED(ch, AFF_FIRESHIELD) {
		send_to_char(ch, libc.CString("You are already covered in a fireshield!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SANCTUARY) {
		send_to_char(ch, libc.CString("You are covered in a barrier!\r\n"))
		return
	}
	if SUNKEN(ch.In_room) {
		send_to_char(ch, libc.CString("There is way too much water here!\r\n"))
		return
	}
	var cost int64 = int64(float64(ch.Max_mana) * 0.03)
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki!\r\n"))
		return
	}
	var skill int = init_skill(ch, SKILL_FIRESHIELD)
	var prob int = axion_dice(0)
	if skill <= prob {
		act(libc.CString("@WYou hold your hands up in front of you on either side and try to summon defensive @rf@Rl@Ya@rm@Re@Ys@W to cover your body. Yet you screw up and the technique fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@W holds $s hands up in front of $m on either side and tries to summon defensive @rf@Rl@Ya@rm@Re@Ys@W to cover $s body. Yet $e seems to screw up and the technique fails!@n"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_FIRESHIELD, 0)
		ch.Mana -= cost
		return
	} else {
		act(libc.CString("@WYou hold your hands up in front of you on either side and try to summon defensive @rf@Rl@Ya@rm@Re@ys@W to cover your body. The ki you have gathered pours out of your body and creates intense black @rf@Rl@Ya@rm@Re@Ys@W that cover your entire body in a protective layer!"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@W holds $s hands up in front of $m on either side and tries to summon defensive @rf@Rl@Ya@rm@Re@ys@W to cover $s body. The ki $e has gathered pours out of $s body and creates intense black @rf@Rl@Ya@rm@Re@Ys@W that cover $s entire body in a protective layer!"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_FIRESHIELD, 0)
		ch.Mana -= cost
		SET_BIT_AR(ch.Affected_by[:], AFF_FIRESHIELD)
		return
	}
}
func do_warppool(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_WARP) == 0 {
		return
	}
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil || ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are struggling with someone!\r\n"))
		return
	}
	if ch.Sits != nil {
		send_to_char(ch, libc.CString("You should get up first.\r\n"))
		return
	}
	var perc int = GET_SKILL(ch, SKILL_WARP)
	var prob int = axion_dice(0)
	var cost int = int(ch.Max_mana / 20)
	var pass int = FALSE
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What planet are you wanting to warp to?\n[ earth | frigid | kanassa | namek | aether ]\r\n"))
		return
	}
	if ch.Mana < int64(cost) {
		send_to_char(ch, libc.CString("You do not have enough ki to perform the technique.\r\n"))
		return
	}
	if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 4600 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 4700 {
		pass = TRUE
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 795 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 1099 {
		pass = TRUE
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 15100 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 0x3BC3 {
		pass = TRUE
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 0x3363 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 0x338F {
		pass = TRUE
	} else if ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) && SECT(ch.In_room) == SECT_WATER_NOSWIM {
		pass = TRUE
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 0x2F47 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 0x3001 {
		pass = TRUE
	}
	if pass == FALSE {
		send_to_char(ch, libc.CString("You must be on or in a sea or ocean for warp pool to work.\r\n"))
		return
	}
	if libc.StrCaseCmp(libc.CString("earth"), &arg[0]) == 0 && ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
		send_to_char(ch, libc.CString("You are already on Earth!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("frigid"), &arg[0]) == 0 && ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) {
		send_to_char(ch, libc.CString("You are already on Frigid!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("kanassa"), &arg[0]) == 0 && ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) {
		send_to_char(ch, libc.CString("You are already on Kanasssa!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("namek"), &arg[0]) == 0 && ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) {
		send_to_char(ch, libc.CString("You are already on Namek!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("aether"), &arg[0]) == 0 && ROOM_FLAGGED(ch.In_room, ROOM_AETHER) {
		send_to_char(ch, libc.CString("You are already on Aether!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("earth"), &arg[0]) == 0 {
		if prob > perc {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. You lose your concentration and the ritual fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly a puzzled look comes across @c$n's @Cface and the water returns to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
			improve_skill(ch, SKILL_WARP, 1)
		} else {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. As you complete the ritual you connect the water you disturbed with the water you envisioned and warp between the two points!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly @c$n@C vanishes into this water! A moment later the waters return to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_WARP, 1)
			char_from_room(ch)
			char_to_room(ch, real_room(850))
			act(libc.CString("@CSuddenly a large whirlpool of flashing water begins to form nearby. After a few seconds @c$n@C pops out of the center of the pool! The water then return to normal a moment laterr...@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
		}
	} else if libc.StrCaseCmp(libc.CString("frigid"), &arg[0]) == 0 {
		if prob > perc {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. You lose your concentration and the ritual fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly a puzzled look comes across @c$n's @Cface and the water returns to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
			improve_skill(ch, SKILL_WARP, 1)
		} else {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. As you complete the ritual you connect the water you disturbed with the water you envisioned and warp between the two points!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly @c$n@C vanishes into this water! A moment later the waters return to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_WARP, 1)
			char_from_room(ch)
			char_to_room(ch, real_room(4609))
			act(libc.CString("@CSuddenly a large whirlpool of flashing water begins to form nearby. After a few seconds @c$n@C pops out of the center of the pool! The water then return to normal a moment laterr...@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
		}
	} else if libc.StrCaseCmp(libc.CString("namek"), &arg[0]) == 0 {
		if prob > perc {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. You lose your concentration and the ritual fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly a puzzled look comes across @c$n's @Cface and the water returns to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
			improve_skill(ch, SKILL_WARP, 1)
		} else {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. As you complete the ritual you connect the water you disturbed with the water you envisioned and warp between the two points!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly @c$n@C vanishes into this water! A moment later the waters return to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_WARP, 1)
			char_from_room(ch)
			char_to_room(ch, real_room(0x2A98))
			act(libc.CString("@CSuddenly a large whirlpool of flashing water begins to form nearby. After a few seconds @c$n@C pops out of the center of the pool! The water then return to normal a moment laterr...@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
		}
	} else if libc.StrCaseCmp(libc.CString("kanassa"), &arg[0]) == 0 {
		if prob > perc {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. You lose your concentration and the ritual fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly a puzzled look comes across @c$n's @Cface and the water returns to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
			improve_skill(ch, SKILL_WARP, 1)
		} else {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. As you complete the ritual you connect the water you disturbed with the water you envisioned and warp between the two points!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly @c$n@C vanishes into this water! A moment later the waters return to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_WARP, 1)
			char_from_room(ch)
			char_to_room(ch, real_room(15100))
			act(libc.CString("@CSuddenly a large whirlpool of flashing water begins to form nearby. After a few seconds @c$n@C pops out of the center of the pool! The water then return to normal a moment laterr...@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
		}
	} else if libc.StrCaseCmp(libc.CString("aether"), &arg[0]) == 0 {
		if prob > perc {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. You lose your concentration and the ritual fails!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly a puzzled look comes across @c$n's @Cface and the water returns to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
			improve_skill(ch, SKILL_WARP, 1)
		} else {
			act(libc.CString("@CYou reach your hand out and begin to swirl nearby water with it. At the same time you release ki into the water and focus your mind on sensing out the distant body of water you wish to travel to. As you complete the ritual you connect the water you disturbed with the water you envisioned and warp between the two points!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C reaches $s hand out and begins to swirl nearby water with it. The water that is being swirled begins to glow @wbright@B blue@C and has a distinct separation from the rest of the waters. Suddenly @c$n@C vanishes into this water! A moment later the waters return to normal.@n"), TRUE, ch, nil, nil, TO_ROOM)
			improve_skill(ch, SKILL_WARP, 1)
			char_from_room(ch)
			char_to_room(ch, real_room(0x2FDC))
			act(libc.CString("@CSuddenly a large whirlpool of flashing water begins to form nearby. After a few seconds @c$n@C pops out of the center of the pool! The water then return to normal a moment laterr...@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= int64(cost)
		}
	} else {
		send_to_char(ch, libc.CString("That is not an acceptable choice. It must be a planet with a large body of water.\n[ earth | frigid | kanassa | namek | aether ]\r\n"))
		return
	}
}
func do_obstruct(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_HYOGA_KABE) == 0 {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
		send_to_char(ch, libc.CString("You can not use this in such a peaceful area.\r\n"))
		return
	}
	if SECT(ch.In_room) == SECT_SPACE || ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		send_to_char(ch, libc.CString("You can not wall off the vastness of space.\r\n"))
		return
	}
	if SECT(ch.In_room) == SECT_FLYING {
		send_to_char(ch, libc.CString("You can not create gravity defying glacial walls.\r\n"))
		return
	}
	var arg [2048]byte
	var skill int = GET_SKILL(ch, SKILL_HYOGA_KABE)
	var prob int = axion_dice(0)
	var cost int = int(float64(ch.Max_mana/int64(skill)) * 2.5)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What direction are you wanting to block off?\n[ N | E | S | W | NE | NW | SE | SW | U | D | I | O ]\r\n"))
		return
	}
	if ch.Mana < int64(cost) {
		send_to_char(ch, libc.CString("You do not have enough ki to perform the technique.\r\n"))
		return
	}
	var dir int = -1
	var dir2 int = -1
	if libc.StrCaseCmp(libc.CString("n"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("N"), &arg[0]) == 0 {
		dir = 0
		dir2 = 2
	} else if libc.StrCaseCmp(libc.CString("e"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("E"), &arg[0]) == 0 {
		dir = 1
		dir2 = 3
	} else if libc.StrCaseCmp(libc.CString("s"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("S"), &arg[0]) == 0 {
		dir = 2
		dir2 = 0
	} else if libc.StrCaseCmp(libc.CString("w"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("W"), &arg[0]) == 0 {
		dir = 3
		dir2 = 1
	} else if libc.StrCaseCmp(libc.CString("u"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("U"), &arg[0]) == 0 {
		dir = 4
		dir2 = 5
	} else if libc.StrCaseCmp(libc.CString("d"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("D"), &arg[0]) == 0 {
		dir = 5
		dir2 = 4
	} else if libc.StrCaseCmp(libc.CString("i"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("I"), &arg[0]) == 0 {
		dir = 10
		dir2 = 11
	} else if libc.StrCaseCmp(libc.CString("o"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("O"), &arg[0]) == 0 {
		dir = 11
		dir2 = 10
	} else if libc.StrCaseCmp(libc.CString("nw"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("NW"), &arg[0]) == 0 {
		dir = 6
		dir2 = 8
	} else if libc.StrCaseCmp(libc.CString("ne"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("NE"), &arg[0]) == 0 {
		dir = 7
		dir2 = 9
	} else if libc.StrCaseCmp(libc.CString("se"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("SE"), &arg[0]) == 0 {
		dir = 8
		dir2 = 6
	} else if libc.StrCaseCmp(libc.CString("sw"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("SW"), &arg[0]) == 0 {
		dir = 9
		dir2 = 7
	} else {
		send_to_char(ch, libc.CString("That is not an acceptable direction.\n[ N | E | S | W | NE | NW | SE | SW | U | D | I | O ]\r\n"))
		return
	}
	if (world[ch.In_room].Dir_option[dir]) == nil {
		send_to_char(ch, libc.CString("That direction does not exist here.\r\n"))
		return
	} else if skill < prob {
		act(libc.CString("@CYou channel your ki and start to create a wall of water, but lose your concentration and the water promptly disappears.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C channels $s ki and starts to create a wall of water, but loses $s concentration and the water promptly disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= int64(cost)
		improve_skill(ch, SKILL_HYOGA_KABE, 0)
		return
	} else {
		var (
			obj     *obj_data
			newroom int = int(world[ch.In_room].Dir_option[dir].To_room)
		)
		if ROOM_FLAGGED(room_rnum(newroom), ROOM_PEACEFUL) {
			send_to_char(ch, libc.CString("You can not block off a peaceful area.\r\n"))
			return
		}
		for obj = world[newroom].Contents; obj != nil; obj = obj.Next_content {
			if GET_OBJ_VNUM(obj) == 79 {
				if obj.Cost == dir2 {
					if skill < prob {
						act(libc.CString("@CYou place your hands on the glacial wall and concentrate. You fail to undo the composition of the wall!@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C places $s hands on the glacial wall and concentrates. Nothing happens...@n"), TRUE, ch, nil, nil, TO_ROOM)
						ch.Mana -= int64(cost / 2)
					} else {
						act(libc.CString("@CYou place your hands on the glacial wall and concentrate. You unfreeze the wall and evaporate the water effortlessly.@n"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C places $s hands on the glacial wall and concentrates. Suddenly the wall melts and then evaporates!@n"), TRUE, ch, nil, nil, TO_ROOM)
						ch.Mana -= int64(cost / 2)
						extract_obj(obj)
					}
					return
				}
			}
		}
		var obj2 *obj_data
		var obj3 *obj_data
		obj2 = read_object(79, VIRTUAL)
		obj_to_room(obj2, room_rnum(newroom))
		obj3 = read_object(79, VIRTUAL)
		obj_to_room(obj3, ch.In_room)
		var strength int64 = int64(float64(((int(ch.Aff_abils.Intel)*skill)*int(ch.Aff_abils.Wis))*20) + float64(ch.Max_mana)*0.001)
		if strength > ch.Max_hit*20 {
			strength = ch.Max_hit + strength/20
		} else if strength > ch.Max_hit*15 {
			strength = ch.Max_hit + strength/15
		} else if strength > ch.Max_hit*10 {
			strength = ch.Max_hit + strength/10
		} else if strength > ch.Max_hit*5 {
			strength = ch.Max_hit + strength/5
		} else if strength > ch.Max_hit*2 {
			strength = ch.Max_hit + strength/2
		}
		obj2.Cost = dir2
		obj2.Weight = strength
		obj3.Cost = dir
		obj3.Weight = strength
		obj2.Fellow_wall = obj3
		obj3.Fellow_wall = obj2
		act(libc.CString("@CYou concentrate and channel your ki. A wall of water starts to form in such a way to block off the direction of your choice. As the wall becomes complete it freezes solid by your will!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C concentrates and channels $s ki. A wall of water starts to form in such a way to block off one of the directions of this area. As the wall becomes complete it freezes solid by @c$n's@C will!@n"), TRUE, ch, nil, nil, TO_ROOM)
		send_to_room(room_rnum(newroom), libc.CString("@cA wall of water forms slowly upward blocking off the %s direction. This wall of water then freezes instantly once it stops growing.@n\r\n"), dirs[dir2])
		improve_skill(ch, SKILL_HYOGA_KABE, 0)
		ch.Mana -= int64(cost)
		return
	}
}
func do_dimizu(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_DIMIZU) == 0 {
		return
	}
	var skill int = GET_SKILL(ch, SKILL_DIMIZU)
	var prob int = axion_dice(0)
	if world[ch.In_room].Geffect < 0 {
		act(libc.CString("@CYou concentrate and distabilie the water, separating the hydrogen and oxygen. The gases dissipate quickly."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C concentrates and the water filling the area seems to shudder. Suddenly the water begins to evaporate as the hydrogen and oxygen are separated."), TRUE, ch, nil, nil, TO_ROOM)
		world[ch.In_room].Geffect = 0
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		return
	} else if SECT(ch.In_room) == SECT_UNDERWATER {
		send_to_char(ch, libc.CString("The area is already underwater!\r\n"))
		return
	} else if SECT(ch.In_room) == SECT_SPACE || ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		send_to_char(ch, libc.CString("You can't flood space!\r\n"))
		return
	} else if ch.Mana < ch.Max_mana/12 {
		send_to_char(ch, libc.CString("You do not have enough ki to perform the technique.\r\n"))
		return
	} else if skill < prob {
		act(libc.CString("@CYou gather your ki and concentrate on creating water from it. Water begins to flow upward around the entire area, but you lose your concentration and it all goes flooding away!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C gathers $s ki and concentrates on creating water from it. Water begins to flow upward around the entire area, but $e loses $s concentration and all the water goes flooding away!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= ch.Max_mana / 12
		improve_skill(ch, SKILL_DIMIZU, 0)
		return
	} else {
		act(libc.CString("@CYou gather your ki and concentrate on creating water from it. Water begins to flow upward around the entire area. You form the water into a perfect cube with barely any ripples in its walls. It will maintain this form for a while.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C gathers $s ki and concentrates on creating water from it. Water begins to flow upward around the entire area. @c$n@C forms the water into a perfect cube with barely any ripples in its walls. It appears the water will maintain this form for a while.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= ch.Max_mana / 12
		world[ch.In_room].Geffect = -3
		improve_skill(ch, SKILL_DIMIZU, 0)
		return
	}
}
func do_beacon(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if AFF_FLAGGED(ch, AFF_SPIRIT) {
		send_to_char(ch, libc.CString("You are dead. You can not stake out a room to return to upon revival.\r\n"))
		return
	} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 0 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 14 {
		send_to_char(ch, libc.CString("You can not stake out an immortal room to be revived in.\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("You stake out the room you are in and will return to it if you die and are revived.\r\n"))
		ch.Droom = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
		return
	}
}
func do_feed(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	var vict *char_data
	var obj *obj_data
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Feed a senzu to whom?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if int(vict.Race) == RACE_ANDROID {
		send_to_char(ch, libc.CString("They are unaffected by senzu beans.\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You need to give them a senzu.\r\n"))
		return
	}
	if int(obj.Type_flag) != ITEM_POTION {
		send_to_char(ch, libc.CString("You can only feed senzu beans.\r\n"))
		return
	}
	if OBJ_FLAGGED(obj, ITEM_FORGED) {
		send_to_char(ch, libc.CString("They can't swallow that, it is fake!\r\n"))
		return
	}
	if OBJ_FLAGGED(obj, ITEM_BROKEN) {
		send_to_char(ch, libc.CString("They can't swallow that, it is broken!\r\n"))
		return
	}
	if vict.Fighting != nil {
		send_to_char(ch, libc.CString("They are a bit busy at the moment!\r\n"))
		return
	}
	if vict.Master != ch && ch.Master != vict && ch.Master != vict.Master {
		send_to_char(ch, libc.CString("You need to be grouped with them first.\r\n"))
		return
	}
	if !AFF_FLAGGED(vict, AFF_GROUP) || !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("You need to be grouped with them first.\r\n"))
		return
	}
	act(libc.CString("@WYou take $p@W and pop it into @C$N@W's mouth!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
	act(libc.CString("@C$n@W takes $p@W and pops it into YOUR mouth!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
	act(libc.CString("@C$n@W takes $p@W and pops it into @c$N@W's mouth!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
	mag_objectmagic(vict, obj, libc.CString(""))
}
func do_spoil(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	var obj *obj_data
	var type_ int = 0
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What corpse do you want to decapitate?\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("No corpse around here by that name.\r\n"))
		return
	}
	if (obj.Value[VAL_CORPSE_HEAD]) == 0 {
		send_to_char(ch, libc.CString("That corpse is already missing its head.\r\n"))
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil {
		if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
			type_ = 1
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
			type_ = 1
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
			type_ = 1
		}
	} else if (ch.Equipment[WEAR_WIELD2]) != nil {
		if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
			type_ = 2
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
			type_ = 2
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
			type_ = 2
		}
	}
	if type_ == 0 {
		act(libc.CString("@C$n@W reaches down and @rtears@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		act(libc.CString("@WYou reach down and @rtear@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
	} else if type_ == 1 {
		act(libc.CString("@C$n@W reaches down and @rcuts@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		act(libc.CString("@WYou reach down and @rcut@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
	} else if type_ == 2 {
		act(libc.CString("@C$n@W reaches down and @rcuts@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		act(libc.CString("@WYou reach down and @rcut@W the head off of @R$p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
	}
	obj.Value[VAL_CORPSE_HEAD] = 0
	var body_part *obj_data
	var part [1000]byte
	var buf [1000]byte
	var buf2 [1000]byte
	var buf3 [1000]byte
	part[0] = '\x00'
	buf[0] = '\x00'
	buf2[0] = '\x00'
	buf3[0] = '\x00'
	body_part = create_obj()
	body_part.Item_number = -1
	body_part.In_room = -1
	stdio.Snprintf(&part[0], int(1000), "%s", obj.Name)
	search_replace(&part[0], libc.CString("headless"), libc.CString(""))
	search_replace(&part[0], libc.CString("corpse"), libc.CString(""))
	search_replace(&part[0], libc.CString("half"), libc.CString(""))
	search_replace(&part[0], libc.CString("burnt"), libc.CString(""))
	search_replace(&part[0], libc.CString("chunks"), libc.CString(""))
	search_replace(&part[0], libc.CString("beaten"), libc.CString(""))
	search_replace(&part[0], libc.CString("bloody"), libc.CString(""))
	trim(&part[0])
	stdio.Snprintf(&buf[0], int(1000), "bloody head %s", &part[0])
	stdio.Snprintf(&buf2[0], int(1000), "@wThe bloody head of %s@w is lying here@n", &part[0])
	stdio.Snprintf(&buf3[0], int(1000), "@wThe bloody head of %s@w@n", &part[0])
	body_part.Name = libc.StrDup(&buf[0])
	body_part.Description = libc.StrDup(&buf2[0])
	body_part.Short_description = libc.StrDup(&buf3[0])
	body_part.Type_flag = ITEM_OTHER
	SET_BIT_AR(body_part.Wear_flags[:], ITEM_WEAR_TAKE)
	SET_BIT_AR(body_part.Extra_flags[:], ITEM_UNIQUE_SAVE)
	body_part.Value[0] = 0
	body_part.Value[1] = 0
	body_part.Value[2] = 0
	body_part.Value[3] = 0
	body_part.Value[4] = 1
	body_part.Value[5] = 1
	body_part.Weight = int64(rand_number(4, 10))
	body_part.Cost_per_day = 0
	add_unique_id(body_part)
	obj_to_room(body_part, ch.In_room)
	obj_from_room(body_part)
	obj_to_char(body_part, ch)
}
