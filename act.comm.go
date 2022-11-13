package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

var languages [9]*byte = [9]*byte{libc.CString("common"), libc.CString("elven"), libc.CString("gnomish"), libc.CString("dwarven"), libc.CString("halfling"), libc.CString("orc"), libc.CString("druid"), libc.CString("draconic"), libc.CString("\n")}

func list_languages(ch *char_data) {
	var (
		a int = 0
		i int
	)
	send_to_char(ch, libc.CString("Languages:\r\n["))
	for i = SKILL_LANG_COMMON; i <= SKILL_LANG_DRACONIC; i++ {
		if GET_SKILL(ch, i) != 0 {
			send_to_char(ch, libc.CString("%s %s%s%s"), func() string {
				if func() int {
					p := &a
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), func() string {
				if ch.Player_specials.Speaking == i {
					return "@r"
				}
				return "@n"
			}(), languages[i-SKILL_LANG_COMMON], "@n")
		}
	}
	send_to_char(ch, libc.CString("%s ]\r\n"), func() string {
		if a == 0 {
			return " None!"
		}
		return ""
	}())
}
func do_voice(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	if IS_NPC(ch) {
		return
	}
	if (ch.Bonuses[BONUS_MUTE]) > 0 {
		send_to_char(ch, libc.CString("You're mute. You don't need to describe your voice.\r\n"))
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("What are you changing your voice description to?\r\n"))
		return
	} else if libc.StrLen(argument) > 75 {
		send_to_char(ch, libc.CString("Your voice description can not be longer than 75 characters.\r\n"))
		return
	} else if libc.StrStr(argument, libc.CString("@")) != nil {
		send_to_char(ch, libc.CString("You can not use colorcode in voice descriptions.\r\n"))
		return
	} else if ch.Voice != nil && ch.Rp < 1 {
		send_to_char(ch, libc.CString("Your voice has already been set. You will need at least 1 RPP to be able to change it.\r\n"))
		return
	} else if ch.Voice != nil {
		send_to_char(ch, libc.CString("Your voice has now been set to: %s\r\n"), argument)
		if ch.Voice != nil {
			libc.Free(unsafe.Pointer(ch.Voice))
		}
		ch.Voice = libc.StrDup(argument)
		ch.Rp -= 1
		ch.Desc.Rpp = ch.Rp
		userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		send_to_char(ch, libc.CString("@D(@cRPP@W: @w-1@D)@n\n\n"))
		return
	} else {
		send_to_char(ch, libc.CString("Your voice has now been set to: %s\r\n"), argument)
		if ch.Voice != nil {
			libc.Free(unsafe.Pointer(ch.Voice))
		}
		ch.Voice = libc.StrDup(argument)
		return
	}
}
func do_languages(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i     int
		found int = FALSE
		arg   [64936]byte
	)
	if config_info.Play.Enable_languages != 0 {
		one_argument(argument, &arg[0])
		if arg[0] == 0 {
			list_languages(ch)
		} else {
			for i = SKILL_LANG_COMMON; i <= SKILL_LANG_DRACONIC; i++ {
				if search_block(&arg[0], &languages[0], FALSE) == i-SKILL_LANG_COMMON && GET_SKILL(ch, i) != 0 {
					ch.Player_specials.Speaking = i
					send_to_char(ch, libc.CString("You now speak %s.\r\n"), languages[i-SKILL_LANG_COMMON])
					found = TRUE
					break
				}
			}
			if found == 0 {
				send_to_char(ch, libc.CString("You do not know of any such language.\r\n"))
				return
			}
		}
	} else {
		send_to_char(ch, libc.CString("But everyone already understands everyone else!\r\n"))
		return
	}
}
func garble_text(string_ *byte, known int, lang int) {
	var (
		letters [50]byte = func() [50]byte {
			var t [50]byte
			copy(t[:], []byte(""))
			return t
		}()
		i int
	)
	switch lang {
	case SKILL_LANG_DWARVEN:
		libc.StrCpy(&letters[0], libc.CString("hprstwxyz"))
	case SKILL_LANG_ELVEN:
		libc.StrCpy(&letters[0], libc.CString("aefhilnopstu"))
	default:
		libc.StrCpy(&letters[0], libc.CString("aehiopstuwxyz"))
	}
	for i = 0; i < libc.StrLen(string_); i++ {
		if libc.IsAlpha(rune(*(*byte)(unsafe.Add(unsafe.Pointer(string_), i)))) && known == 0 {
			*(*byte)(unsafe.Add(unsafe.Pointer(string_), i)) = letters[rand_number(0, libc.StrLen(&letters[0])-1)]
		}
	}
}
func do_osay(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	if IS_NPC(ch) {
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Yes, but WHAT do you want to osay?\r\n"))
		return
	} else {
		var (
			buf  [2048]byte
			buf2 [2048]byte
		)
		stdio.Sprintf(&buf[0], "@WYou @D[@mOSAY@D] @W'@w%s@W'@n", argument)
		if !PRF_FLAGGED(ch, PRF_HIDE) {
			stdio.Sprintf(&buf2[0], "@W%s @D[@mOSAY@D] @W'@w%s@W'@n", func() *byte {
				if ch.Admlevel > 0 {
					return GET_NAME(ch)
				}
				return ch.Desc.User
			}(), argument)
		}
		if PRF_FLAGGED(ch, PRF_HIDE) {
			stdio.Sprintf(&buf2[0], "@WAnonymous @D[@mOSAY@D] @W'@w%s@W'@n", argument)
		}
		act(&buf[0], FALSE, ch, nil, nil, TO_CHAR)
		act(&buf2[0], FALSE, ch, nil, nil, TO_ROOM)
	}
}
func do_say(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d       *descriptor_data
		wch     *char_data = nil
		wch2    *char_data = nil
		wch3    *char_data = nil
		tch     *char_data = nil
		sch     *char_data = nil
		obj     *obj_data  = nil
		granted int        = FALSE
		found   int        = FALSE
	)
	_ = found
	var buf2 [2048]byte
	buf2[0] = '\x00'
	skip_spaces(&argument)
	if (ch.Bonuses[BONUS_MUTE]) > 0 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) > 160 {
		send_to_char(ch, libc.CString("You are mute and unable to talk though.\r\n"))
		return
	} else if (ch.Bonuses[BONUS_MUTE]) > 0 {
		send_to_char(ch, libc.CString("You are mute and unable to talk though. You will be allowed to just for MUD School."))
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Yes, but WHAT do you want to say?\r\n"))
		return
	} else {
		var (
			buf  [2118]byte
			verb [10]byte
		)
		if *(*byte)(unsafe.Add(unsafe.Pointer(argument), libc.StrLen(argument)-1)) == '!' {
			libc.StrCpy(&verb[0], libc.CString("exclaim"))
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(argument), libc.StrLen(argument)-1)) == '?' {
			libc.StrCpy(&verb[0], libc.CString("ask"))
		} else {
			libc.StrCpy(&verb[0], libc.CString("say"))
		}
		for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = tch.Next_in_room {
			if tch != ch && tch.Desc != nil {
				var sayto [100]byte
				stdio.Sprintf(&sayto[0], "to %s ", GET_NAME(tch))
				if libc.StrStr(argument, &sayto[0]) != nil {
					var saytoo [200]byte
					verb[0] = '\x00'
					stdio.Sprintf(&saytoo[0], "says to @g%s@W", GET_NAME(tch))
					search_replace(argument, &sayto[0], libc.CString(""))
					libc.StrCpy(&verb[0], &saytoo[0])
					sch = tch
				} else if !IS_NPC(tch) && !IS_NPC(ch) {
					if readIntro(ch, tch) == 1 {
						stdio.Sprintf(&sayto[0], "to %s ", get_i_name(ch, tch))
					}
					if libc.StrStr(argument, &sayto[0]) != nil {
						var saytoo [200]byte
						verb[0] = '\x00'
						stdio.Sprintf(&saytoo[0], "says to @g%s@W", GET_NAME(tch))
						search_replace(argument, &sayto[0], libc.CString(""))
						libc.StrCpy(&verb[0], &saytoo[0])
						sch = tch
					}
				}
			}
		}
		if sch == nil {
			stdio.Snprintf(&buf[0], int(2118), "@w$n @W%ss, '@C%s@W'@n", &verb[0], argument)
			act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
		} else {
			stdio.Snprintf(&buf[0], int(2118), "@w$n @Wsays to @g$N@W, '@C%s@W'@n", argument)
			stdio.Snprintf(&buf2[0], int(2048), "@w$n @Wsays to @gyou@W, '@C%s@W'@n", argument)
			act(&buf2[0], TRUE, ch, nil, unsafe.Pointer(sch), TO_VICT)
			act(&buf[0], TRUE, ch, nil, unsafe.Pointer(sch), TO_NOTVICT)
		}
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else {
			if libc.StrStr(&verb[0], libc.CString("says to")) != nil {
				var saytoo [200]byte
				verb[0] = '\x00'
				stdio.Sprintf(&saytoo[0], "say to @g%s@W", GET_NAME(sch))
				libc.StrCpy(&verb[0], &saytoo[0])
			}
			stdio.Snprintf(&buf[0], int(2118), "@WYou %s, '@C%s@W'@n\r\n", &verb[0], argument)
			send_to_char(ch, libc.CString("%s"), &buf[0])
			add_history(ch, &buf[0], HIST_SAY)
			if SHENRON == TRUE {
				if (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) == room_vnum(DRAGONR) && (func() room_vnum {
					if EDRAGON.In_room != room_rnum(-1) && EDRAGON.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(EDRAGON.In_room)))).Number
					}
					return -1
				}()) == room_vnum(DRAGONR) {
					if libc.StrStr(argument, libc.CString("wish")) != nil {
						for d = descriptor_list; d != nil; d = d.Next {
							if d.Connected != CON_PLAYING {
								continue
							}
							if libc.StrStr(argument, GET_NAME(d.Character)) != nil && wch == nil {
								wch = d.Character
								found = TRUE
							} else if libc.StrStr(argument, GET_NAME(d.Character)) != nil && wch2 == nil {
								wch2 = d.Character
							} else if libc.StrStr(argument, GET_NAME(d.Character)) != nil && wch3 == nil {
								wch3 = d.Character
							}
						}
						if wch == nil && libc.StrStr(argument, libc.CString("myself")) != nil {
							wch = ch
						}
						if wch == nil {
							return
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("knowledge")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s now has more knowledge!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Player_specials.Class_skill_points[wch.Chclass] += rand_number(2000, 5000)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a knowledge wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("speed")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now faster!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Cha += 10
								if int(wch.Real_abils.Cha) > 100 {
									wch.Real_abils.Cha = 100
								}
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a speed wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("tough")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now tougher!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Armor += 5000
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a tough wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("strength")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s has more strength!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Str += 10
								if int(wch.Real_abils.Str) > 100 {
									wch.Real_abils.Str = 100
								}
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a strength wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("intelligence")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now smarter!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Intel += 10
								if int(wch.Real_abils.Intel) > 100 {
									wch.Real_abils.Intel = 100
								}
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a intelligence wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("wisdom")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now wiser!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Wis += 10
								if int(wch.Real_abils.Wis) > 100 {
									wch.Real_abils.Wis = 100
								}
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a wisdom wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("agility")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now more agile!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Dex += 10
								if int(wch.Real_abils.Dex) > 100 {
									wch.Real_abils.Dex = 100
								}
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a agility wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("constitution")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s has more guts!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Real_abils.Con += 10
								if int(wch.Real_abils.Con) > 100 {
									wch.Real_abils.Con = 100
								}
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a constitutionwish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("skill")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s has more skill!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								var roll int = rand_number(1, 3)
								send_to_char(wch, libc.CString("@GYou suddenly feel like you could learn %d more skills!@n\r\n"), roll)
								wch.Skill_slots += roll
								save_char(wch)
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a skill wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("power")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish cannot be granted, You might want to try something else instead, mortal!@w'@n\r\n"))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("money")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s now has become richer!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Bank_gold += 1000000
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a money wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("immunity")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s now has immunity to Burn, Freezing, Mind Break, Poison, Blindness, Yoikominminken, and Paralysis!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Affected_by[int(AFF_IMMUNITY/32)] |= 1 << (int(AFF_IMMUNITY % 32))
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a immunity wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("vitality")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish cannot be granted, You might want to try something else instead, mortal!%s@w'@n\r\n"))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("revive")) != nil {
							var count int = 0
							if wch != nil {
								count += 1
							}
							if wch2 != nil {
								count += 1
							}
							if wch3 != nil {
								count += 1
							}
							if count == 1 {
								if !AFF_FLAGGED(wch, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch))
								} else {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s has returned to life!%s@w'@n\r\n"), GET_NAME(wch), func() string {
										if WISH[0] != 0 {
											return ""
										}
										return " Now make your second wish."
									}())
									if real_room(wch.Droom) == room_rnum(-1) {
										wch.Droom = 300
									}
									if real_room(wch.Droom) != room_rnum(-1) {
										char_from_room(wch)
										if wch.Droom > 0 {
											char_to_room(wch, real_room(wch.Droom))
										} else {
											char_to_room(wch, real_room(300))
										}
										look_at_room(wch.In_room, wch, 0)
										send_to_char(wch, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									granted = TRUE
									SELFISHMETER -= 2
									mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a revive wish on %s."), GET_NAME(ch), GET_NAME(wch))
								}
							}
							if count == 2 {
								if !AFF_FLAGGED(wch, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch))
								}
								if !AFF_FLAGGED(wch2, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch2))
								} else if AFF_FLAGGED(wch, AFF_SPIRIT) && AFF_FLAGGED(wch2, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s and %s have returned to life!%s@w'@n\r\n"), GET_NAME(wch), GET_NAME(wch2), func() string {
										if WISH[0] != 0 {
											return ""
										}
										return " Now make your second wish."
									}())
									if real_room(wch.Droom) == room_rnum(-1) {
										wch.Droom = 300
									}
									if real_room(wch.Droom) != room_rnum(-1) {
										char_from_room(wch)
										char_to_room(wch, real_room(wch.Droom))
										look_at_room(wch.In_room, wch, 0)
										send_to_char(wch, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									if real_room(wch2.Droom) == room_rnum(-1) {
										wch2.Droom = 300
									}
									if real_room(wch2.Droom) != room_rnum(-1) {
										char_from_room(wch2)
										char_to_room(wch2, real_room(wch2.Droom))
										look_at_room(wch2.In_room, wch2, 0)
										send_to_char(wch2, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch2.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch2.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									granted = TRUE
									SELFISHMETER -= 3
									mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a revive wish on %s."), GET_NAME(ch), GET_NAME(wch2))
									WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
								}
							}
							if count == 3 {
								if !AFF_FLAGGED(wch, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch))
								}
								if !AFF_FLAGGED(wch2, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch2))
								}
								if !AFF_FLAGGED(wch3, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@C%s is not dead, and can not be revived.@w'@n\r\n"), GET_NAME(wch3))
								} else if AFF_FLAGGED(wch, AFF_SPIRIT) && AFF_FLAGGED(wch2, AFF_SPIRIT) && AFF_FLAGGED(wch3, AFF_SPIRIT) {
									send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s, %s, and %s have returned to life!!%s@w'@n\r\n"), GET_NAME(wch), GET_NAME(wch2), GET_NAME(wch3), func() string {
										if WISH[0] != 0 {
											return ""
										}
										return " Now make your second wish."
									}())
									if real_room(wch.Droom) == room_rnum(-1) {
										wch.Droom = 300
									}
									if real_room(wch.Droom) != room_rnum(-1) {
										char_from_room(wch)
										char_to_room(wch, real_room(wch.Droom))
										look_at_room(wch.In_room, wch, 0)
										send_to_char(wch, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									if real_room(wch2.Droom) == room_rnum(-1) {
										wch2.Droom = 300
									}
									if real_room(wch2.Droom) != room_rnum(-1) {
										char_from_room(wch2)
										char_to_room(wch2, real_room(wch2.Droom))
										look_at_room(wch2.In_room, wch2, 0)
										send_to_char(wch2, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch2.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch2.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									if real_room(wch3.Droom) == room_rnum(-1) {
										wch3.Droom = 300
									}
									if real_room(wch3.Droom) != room_rnum(-1) {
										char_from_room(wch3)
										char_to_room(wch3, real_room(wch3.Droom))
										look_at_room(wch3.In_room, wch3, 0)
										send_to_char(wch3, libc.CString("@wYou smile as the golden halo above your head disappears! You have returned to life where you had last died!@n\r\n"))
										wch3.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
										wch3.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
									}
									granted = TRUE
									SELFISHMETER -= 3
									mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a revive wish on %s and %s."), GET_NAME(ch), GET_NAME(wch2), GET_NAME(wch3))
									WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
								}
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("immortal")) != nil && WISH[0] == 0 {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now immortal!@w'@n\r\n"), GET_NAME(wch))
								wch.Act[int(PLR_IMMORTAL/32)] |= bitvector_t(int32(1 << (int(PLR_IMMORTAL % 32))))
								WISH[0] = 1
								WISH[1] = 1
								granted = TRUE
								SELFISHMETER += 4
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a immortal wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("immortal")) != nil && WISH[0] == 1 {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CI can not grant that wish, there is not enough remaining power in this summoning!@w'@n\r\n"))
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString(" mortal")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s is now mortal!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								wch.Act[int(PLR_IMMORTAL/32)] &= bitvector_t(int32(^(1 << (int(PLR_IMMORTAL % 32)))))
								granted = TRUE
								SELFISHMETER += 4
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a mortal wish on %s."), GET_NAME(ch), GET_NAME(wch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("senzu")) != nil {
							if wch != nil {
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								obj = read_object(1, VIRTUAL)
								obj_to_char(obj, ch)
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s now possesses 10 senzus!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								granted = TRUE
								SELFISHMETER += 1
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a senzu wish."), GET_NAME(ch))
							}
						}
						if granted == FALSE && libc.StrStr(argument, libc.CString("roleplay")) != nil {
							if wch != nil {
								send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CYour wish has been granted, %s!%s@w'@n\r\n"), GET_NAME(wch), func() string {
									if WISH[0] != 0 {
										return ""
									}
									return " Now make your second wish."
								}())
								granted = TRUE
								mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Shenron: %s has made a roleplay wish."), GET_NAME(ch))
								WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
							}
						}
						if granted == TRUE {
							if WISH[0] == 1 {
								WISH[1] = 1
							} else {
								WISH[0] = 1
							}
							if wch != nil {
								save_char(wch)
							}
							if wch2 != nil {
								save_char(wch)
							}
							if wch3 != nil {
								save_char(wch)
							}
							save_mud_time(&time_info)
						} else if wch == nil {
							send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CThat person does not exist, make another wish.'@n\r\n"))
						} else {
							send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CDo not waste my time with wishes I can not grant...@w'@n\r\n"))
						}
					}
				}
			}
		}
	}
	speech_mtrigger(ch, argument)
	speech_wtrigger(ch, argument)
	if SHENRON == FALSE || SHENRON == TRUE && ch.In_room != real_room(room_vnum(DRAGONR)) {
		mob_talk(ch, argument)
	}
}
func do_gsay(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		k    *char_data
		f    *follow_type
		blah [2048]byte
	)
	skip_spaces(&argument)
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("But you are not the member of a group!\r\n"))
		return
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("Lol, no.\r\n"))
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Yes, but WHAT do you want to group-say?\r\n"))
	} else {
		var buf [64936]byte
		if ch.Master != nil {
			k = ch.Master
		} else {
			k = ch
		}
		libc.StrCpy(&buf[0], argument)
		stdio.Sprintf(&blah[0], "$n@W tells the group @W'@G%s@W'@n\r\n", &buf[0])
		if AFF_FLAGGED(k, AFF_GROUP) && k != ch && AWAKE(k) {
			if config_info.Play.Enable_languages != 0 {
				send_to_char(k, libc.CString("%s@W tells the group%s @W'@G%s@W'@n\r\n"), func() *byte {
					if CAN_SEE(k, ch) {
						return GET_NAME(ch)
					}
					return libc.CString("Someone")
				}(), func() string {
					if GET_SKILL(k, ch.Player_specials.Speaking) != 0 {
						return ","
					}
					return ", in an unfamiliar tongue,"
				}(), &buf[0])
			} else {
				act(&blah[0], TRUE, ch, nil, unsafe.Pointer(k), TO_VICT)
			}
		}
		for f = k.Followers; f != nil; f = f.Next {
			if AFF_FLAGGED(f.Follower, AFF_GROUP) && f.Follower != ch && AWAKE(f.Follower) {
				if !IS_NPC(ch) && !IS_NPC(f.Follower) && config_info.Play.Enable_languages != 0 {
					garble_text(&buf[0], GET_SKILL(f.Follower, ch.Player_specials.Speaking), ch.Player_specials.Speaking)
				} else {
					garble_text(&buf[0], 1, SKILL_LANG_COMMON)
				}
				if config_info.Play.Enable_languages != 0 {
					send_to_char(f.Follower, libc.CString("%s@W tells the group%s @W'%s@W'@n\r\n"), func() *byte {
						if CAN_SEE(f.Follower, ch) {
							return GET_NAME(ch)
						}
						return libc.CString("Someone")
					}(), func() string {
						if GET_SKILL(f.Follower, ch.Player_specials.Speaking) != 0 {
							return ","
						}
						return ", in an unfamiliar tongue,"
					}(), &buf[0])
				} else {
					act(&blah[0], TRUE, ch, nil, unsafe.Pointer(f.Follower), int(TO_VICT|2<<7))
				}
			}
		}
		if PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else {
			send_to_char(ch, libc.CString("@WYou tell the group, '@G%s@W'@n\r\n"), argument)
		}
	}
}
func perform_tell(ch *char_data, vict *char_data, arg *byte) {
	var (
		buf  [64936]byte
		buf2 [64936]byte
	)
	libc.StrCpy(&buf[0], arg)
	if config_info.Play.Enable_languages != 0 {
		stdio.Snprintf(&buf2[0], int(64936), "@[13]%s tells you%s '%s@[13]'@n\r\n", func() *byte {
			if CAN_SEE(vict, ch) {
				return GET_NAME(ch)
			}
			return libc.CString("Someone")
		}(), func() string {
			if GET_SKILL(vict, ch.Player_specials.Speaking) != 0 {
				return ","
			}
			return ", in an unfamiliar tongue,"
		}(), &buf[0])
		send_to_char(vict, libc.CString("%s"), &buf2[0])
		add_history(vict, &buf2[0], HIST_TELL)
	} else if !IS_NPC(ch) && vict.Admlevel < 1 {
		stdio.Snprintf(&buf2[0], int(64936), "@Y%s@Y tells you '%s'@n\r\n", func() *byte {
			if ch.Admlevel > 0 {
				return GET_NAME(ch)
			}
			return ch.Desc.User
		}(), &buf[0])
		send_to_char(vict, libc.CString("%s"), &buf2[0])
		add_history(vict, &buf2[0], HIST_TELL)
	} else if !IS_NPC(ch) && vict.Admlevel >= 1 {
		stdio.Snprintf(&buf2[0], int(64936), "@Y%s(%s)@Y tells you '%s'@n\r\n", ch.Desc.User, GET_NAME(ch), &buf[0])
		send_to_char(vict, libc.CString("%s"), &buf2[0])
		add_history(vict, &buf2[0], HIST_TELL)
	} else if IS_NPC(ch) {
		stdio.Snprintf(&buf2[0], int(64936), "@Y%s@Y tells you '%s'@n\r\n", GET_NAME(ch), &buf[0])
		send_to_char(vict, libc.CString("%s"), &buf2[0])
	}
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	} else {
		if !IS_NPC(ch) {
			stdio.Snprintf(&buf2[0], int(64936), "@YYou tell %s, '%s'@n\r\n", func() *byte {
				if vict.Admlevel > 0 {
					return GET_NAME(vict)
				}
				if vict.Desc.User != nil {
					return vict.Desc.User
				}
				return libc.CString("ERROR")
			}(), arg)
			if ch.Admlevel < 5 && vict.Admlevel < 5 && !IS_NPC(ch) && !IS_NPC(vict) {
				send_to_imm(libc.CString("@GTELL: @C%s@G tells @c%s, @W'@w%s@W'@n"), func() *byte {
					if ch.Admlevel > 0 {
						return GET_NAME(ch)
					}
					return GET_USER(ch)
				}(), func() *byte {
					if vict.Admlevel > 0 {
						return GET_NAME(vict)
					}
					return GET_USER(vict)
				}(), arg)
			}
			send_to_char(ch, libc.CString("%s"), &buf2[0])
			add_history(ch, &buf2[0], HIST_TELL)
		} else {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		}
	}
	if !IS_NPC(vict) && !IS_NPC(ch) {
		vict.Player_specials.Last_tell = ch.Idnum
	}
}
func is_tell_ok(ch *char_data, vict *char_data) int {
	if ch == vict {
		send_to_char(ch, libc.CString("You try to tell yourself something.\r\n"))
	} else if !IS_NPC(ch) && GET_LEVEL(ch) < 3 && vict.Admlevel < 1 {
		send_to_char(ch, libc.CString("You need to be level 3 or higher to send or receive tells"))
	} else if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_AFK) && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("You can't send tells when AFK.\r\n"))
	} else if !IS_NPC(ch) && PRF_FLAGGED(vict, PRF_AFK) && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("They are AFK right now, try later.\r\n"))
	} else if !IS_NPC(ch) && PRF_FLAGGED(vict, PRF_AFK) && ch.Admlevel >= 1 {
		return TRUE
	} else if !IS_NPC(vict) && GET_LEVEL(vict) < 3 && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("They need to be level 3 or higher to send or receive tells"))
	} else if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOTELL) && vict.Admlevel < 1 {
		send_to_char(ch, libc.CString("You can't tell other people while you have notell on.\r\n"))
	} else if ROOM_FLAGGED(ch.In_room, ROOM_SOUNDPROOF) && vict.Admlevel < 1 {
		send_to_char(ch, libc.CString("The walls seem to absorb your words.\r\n"))
	} else if IS_NPC(vict) {
		send_to_char(ch, libc.CString("You can't send tells to mobs.\r\n"))
	} else if !IS_NPC(vict) && vict.Desc == nil {
		act(libc.CString("$E's linkless at the moment."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
	} else if PLR_FLAGGED(vict, PLR_WRITING) {
		act(libc.CString("$E's writing a message right now; try again later."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
	} else if !IS_NPC(vict) && ch.Admlevel < 1 && PRF_FLAGGED(vict, PRF_NOTELL) || ROOM_FLAGGED(vict.In_room, ROOM_SOUNDPROOF) {
		act(libc.CString("$E can't hear you."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
	} else {
		return TRUE
	}
	return FALSE
}
func do_tell(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		buf  [2048]byte
		buf2 [2048]byte
	)
	half_chop(argument, &buf[0], &buf2[0])
	if buf[0] == 0 || buf2[0] == 0 {
		send_to_char(ch, libc.CString("Who do you wish to tell what??\r\n"))
		return
	}
	var k *descriptor_data
	var found int = FALSE
	if !IS_NPC(ch) {
		stdio.Sprintf(&buf[0], "%s", CAP(&buf[0]))
		for k = descriptor_list; k != nil; k = k.Next {
			if IS_NPC(k.Character) {
				continue
			}
			if k.Connected != CON_PLAYING {
				continue
			}
			if k.User == nil {
				continue
			}
			if found == FALSE && !IS_NPC(ch) && (libc.StrCaseCmp(k.User, &buf[0]) == 0 || libc.StrStr(k.User, &buf[0]) != nil) {
				vict = k.Character
				found = TRUE
			} else if !IS_NPC(ch) && found == FALSE && (libc.StrCaseCmp(GET_NAME(k.Character), &buf[0]) == 0 || libc.StrStr(GET_NAME(k.Character), &buf[0]) != nil) && k.Character.Admlevel > 0 {
				vict = k.Character
				found = TRUE
			}
		}
	}
	if found == FALSE && !IS_NPC(ch) {
		send_to_char(ch, libc.CString("No user around with that name."))
	} else if IS_NPC(ch) && (func() *char_data {
		vict = get_player_vis(ch, &buf[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if is_tell_ok(ch, vict) != 0 {
		perform_tell(ch, vict, &buf2[0])
	}
}
func do_reply(ch *char_data, argument *byte, cmd int, subcmd int) {
	var tch *char_data = character_list
	if IS_NPC(ch) {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
		send_to_char(ch, libc.CString("This is a different dimension!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
		send_to_char(ch, libc.CString("This is the past, you can't send tells!\r\n"))
		return
	}
	skip_spaces(&argument)
	if int(ch.Player_specials.Last_tell) == int(-1) {
		send_to_char(ch, libc.CString("You have nobody to reply to!\r\n"))
	} else if *argument == 0 {
		send_to_char(ch, libc.CString("What is your reply?\r\n"))
	} else {
		for tch != nil && (IS_NPC(tch) || int(tch.Idnum) != int(ch.Player_specials.Last_tell)) {
			tch = tch.Next
		}
		if tch == nil {
			send_to_char(ch, libc.CString("They are no longer playing.\r\n"))
		} else if is_tell_ok(ch, tch) != 0 {
			perform_tell(ch, tch, argument)
		}
	}
}
func do_spec_comm(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf           [2048]byte
		buf2          [2048]byte
		vict          *char_data
		action_sing   *byte
		action_plur   *byte
		action_others *byte
	)
	if (ch.Bonuses[BONUS_MUTE]) > 0 {
		send_to_char(ch, libc.CString("You are mute and unable to talk though.\r\n"))
		return
	}
	switch subcmd {
	case SCMD_WHISPER:
		action_sing = libc.CString("whisper to")
		action_plur = libc.CString("whispers to")
		action_others = libc.CString("$n whispers something to $N.")
	case SCMD_ASK:
		action_sing = libc.CString("ask")
		action_plur = libc.CString("asks")
		action_others = libc.CString("$n asks $N a question.")
	default:
		action_sing = libc.CString("oops")
		action_plur = libc.CString("oopses")
		action_others = libc.CString("$n is tongue-tied trying to speak with $N.")
	}
	half_chop(argument, &buf[0], &buf2[0])
	if buf[0] == 0 || buf2[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you want to %s.. and what??\r\n"), action_sing)
	} else if (func() *char_data {
		vict = get_char_vis(ch, &buf[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if vict == ch {
		send_to_char(ch, libc.CString("You can't get your mouth close enough to your ear...\r\n"))
	} else {
		var (
			buf1 [64936]byte
			obuf [64936]byte
		)
		if config_info.Play.Enable_languages != 0 {
			libc.StrCpy(&obuf[0], &buf2[0])
			garble_text(&obuf[0], GET_SKILL(vict, ch.Player_specials.Speaking), ch.Player_specials.Speaking)
			stdio.Snprintf(&buf1[0], int(64936), "$n %s you%s '%s'", action_plur, func() string {
				if GET_SKILL(vict, ch.Player_specials.Speaking) != 0 {
					return ","
				}
				return ", in an unfamiliar tongue,"
			}(), &obuf[0])
		} else {
			stdio.Snprintf(&buf1[0], int(64936), "@c$n @W%s you '@m%s@W'@n", action_plur, &buf2[0])
		}
		act(&buf1[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else {
			var blum [2048]byte
			stdio.Sprintf(&blum[0], "@WYou %s @C$N@W, '@m%s@W'@n\r\n", action_sing, &buf2[0])
			act(&blum[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		}
		if subcmd == SCMD_WHISPER {
			act(action_others, FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			handle_whisper(&buf2[0], ch, vict)
		} else {
			act(action_others, FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		}
	}
}
func handle_whisper(buf *byte, ch *char_data, vict *char_data) {
	var tch *char_data
	for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = tch.Next_in_room {
		if IS_NPC(tch) {
			continue
		}
		if tch == ch {
			continue
		}
		if tch == vict {
			continue
		}
		if GET_SKILL(tch, SKILL_LISTEN) == 0 {
			continue
		}
		if GET_SKILL(tch, SKILL_LISTEN) != 0 {
			var (
				skill int = GET_SKILL(tch, SKILL_LISTEN)
				roll1 int = rand_number(10, 30)
				roll  int = rand_number(roll1, 110)
			)
			if skill >= roll {
				send_to_char(tch, libc.CString("@WYou overhear everything whispered, @W'@m%s@W'@n\r\n"), overhear(buf, 3))
			} else if skill+10 >= roll {
				send_to_char(tch, libc.CString("@WYou overhear a lot of what is whispered, @W'@m%s@W'@n\r\n"), overhear(buf, 2))
			} else if skill+20 >= roll {
				send_to_char(tch, libc.CString("@WYou overhear some of what is whispered, @W'@m%s@W'@n\r\n"), overhear(buf, 1))
			} else if skill+30 >= roll {
				send_to_char(tch, libc.CString("@WYou overhear little of what is whispered, @W'@m%s@W'@n\r\n"), overhear(buf, 0))
			} else {
				send_to_char(tch, libc.CString("@WYou were unable to overhear anything that was whispered.@n\r\n"))
			}
		}
	}
}
func overhear(buf *byte, type_ int) *byte {
	switch type_ {
	case 0:
		if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("a"), libc.CString(".."))
			search_replace(buf, libc.CString("A"), libc.CString(".."))
			search_replace(buf, libc.CString("h"), libc.CString(".."))
			search_replace(buf, libc.CString("H"), libc.CString(".."))
			search_replace(buf, libc.CString("e"), libc.CString(".."))
			search_replace(buf, libc.CString("E"), libc.CString(".."))
			search_replace(buf, libc.CString("m"), libc.CString(".."))
			search_replace(buf, libc.CString("M"), libc.CString(".."))
			search_replace(buf, libc.CString("o"), libc.CString(".."))
			search_replace(buf, libc.CString("O"), libc.CString(".."))
			search_replace(buf, libc.CString("p"), libc.CString(".."))
			search_replace(buf, libc.CString("P"), libc.CString(".."))
			search_replace(buf, libc.CString("y"), libc.CString(".."))
			search_replace(buf, libc.CString("Y"), libc.CString(".."))
			search_replace(buf, libc.CString("j"), libc.CString(".."))
			search_replace(buf, libc.CString("J"), libc.CString(".."))
			search_replace(buf, libc.CString("k"), libc.CString(".."))
			search_replace(buf, libc.CString("K"), libc.CString(".."))
			search_replace(buf, libc.CString("d"), libc.CString(".."))
			search_replace(buf, libc.CString("D"), libc.CString(".."))
			search_replace(buf, libc.CString("w"), libc.CString(".."))
			search_replace(buf, libc.CString("W"), libc.CString(".."))
		} else if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("e"), libc.CString(".."))
			search_replace(buf, libc.CString("E"), libc.CString(".."))
			search_replace(buf, libc.CString("r"), libc.CString(".."))
			search_replace(buf, libc.CString("R"), libc.CString(".."))
			search_replace(buf, libc.CString("k"), libc.CString(".."))
			search_replace(buf, libc.CString("K"), libc.CString(".."))
			search_replace(buf, libc.CString("m"), libc.CString(".."))
			search_replace(buf, libc.CString("M"), libc.CString(".."))
			search_replace(buf, libc.CString("o"), libc.CString(".."))
			search_replace(buf, libc.CString("O"), libc.CString(".."))
			search_replace(buf, libc.CString("p"), libc.CString(".."))
			search_replace(buf, libc.CString("P"), libc.CString(".."))
			search_replace(buf, libc.CString("y"), libc.CString(".."))
			search_replace(buf, libc.CString("Y"), libc.CString(".."))
			search_replace(buf, libc.CString("j"), libc.CString(".."))
			search_replace(buf, libc.CString("J"), libc.CString(".."))
			search_replace(buf, libc.CString("k"), libc.CString(".."))
			search_replace(buf, libc.CString("K"), libc.CString(".."))
			search_replace(buf, libc.CString("d"), libc.CString(".."))
			search_replace(buf, libc.CString("D"), libc.CString(".."))
			search_replace(buf, libc.CString("w"), libc.CString(".."))
			search_replace(buf, libc.CString("W"), libc.CString(".."))
		} else {
			search_replace(buf, libc.CString("s"), libc.CString(".."))
			search_replace(buf, libc.CString("S"), libc.CString(".."))
			search_replace(buf, libc.CString("r"), libc.CString(".."))
			search_replace(buf, libc.CString("R"), libc.CString(".."))
			search_replace(buf, libc.CString("c"), libc.CString(".."))
			search_replace(buf, libc.CString("C"), libc.CString(".."))
			search_replace(buf, libc.CString("q"), libc.CString(".."))
			search_replace(buf, libc.CString("Q"), libc.CString(".."))
			search_replace(buf, libc.CString("l"), libc.CString(".."))
			search_replace(buf, libc.CString("L"), libc.CString(".."))
			search_replace(buf, libc.CString("u"), libc.CString(".."))
			search_replace(buf, libc.CString("U"), libc.CString(".."))
			search_replace(buf, libc.CString("i"), libc.CString(".."))
			search_replace(buf, libc.CString("I"), libc.CString(".."))
			search_replace(buf, libc.CString("z"), libc.CString(".."))
			search_replace(buf, libc.CString("Z"), libc.CString(".."))
			search_replace(buf, libc.CString("t"), libc.CString(".."))
			search_replace(buf, libc.CString("T"), libc.CString(".."))
		}
		return buf
	case 1:
		if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("b"), libc.CString(".."))
			search_replace(buf, libc.CString("B"), libc.CString(".."))
			search_replace(buf, libc.CString("f"), libc.CString(".."))
			search_replace(buf, libc.CString("F"), libc.CString(".."))
			search_replace(buf, libc.CString("g"), libc.CString(".."))
			search_replace(buf, libc.CString("G"), libc.CString(".."))
			search_replace(buf, libc.CString("v"), libc.CString(".."))
			search_replace(buf, libc.CString("V"), libc.CString(".."))
			search_replace(buf, libc.CString("j"), libc.CString(".."))
			search_replace(buf, libc.CString("J"), libc.CString(".."))
			search_replace(buf, libc.CString("k"), libc.CString(".."))
			search_replace(buf, libc.CString("K"), libc.CString(".."))
		} else if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("d"), libc.CString(".."))
			search_replace(buf, libc.CString("D"), libc.CString(".."))
			search_replace(buf, libc.CString("y"), libc.CString(".."))
			search_replace(buf, libc.CString("Y"), libc.CString(".."))
			search_replace(buf, libc.CString("m"), libc.CString(".."))
			search_replace(buf, libc.CString("M"), libc.CString(".."))
			search_replace(buf, libc.CString("h"), libc.CString(".."))
			search_replace(buf, libc.CString("H"), libc.CString(".."))
			search_replace(buf, libc.CString("s"), libc.CString(".."))
			search_replace(buf, libc.CString("S"), libc.CString(".."))
			search_replace(buf, libc.CString("t"), libc.CString(".."))
			search_replace(buf, libc.CString("T"), libc.CString(".."))
		} else {
			search_replace(buf, libc.CString("a"), libc.CString(".."))
			search_replace(buf, libc.CString("A"), libc.CString(".."))
			search_replace(buf, libc.CString("r"), libc.CString(".."))
			search_replace(buf, libc.CString("R"), libc.CString(".."))
			search_replace(buf, libc.CString("n"), libc.CString(".."))
			search_replace(buf, libc.CString("N"), libc.CString(".."))
			search_replace(buf, libc.CString("o"), libc.CString(".."))
			search_replace(buf, libc.CString("O"), libc.CString(".."))
		}
		return buf
	case 2:
		if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("q"), libc.CString(".."))
			search_replace(buf, libc.CString("Q"), libc.CString(".."))
			search_replace(buf, libc.CString("o"), libc.CString(".."))
			search_replace(buf, libc.CString("O"), libc.CString(".."))
			search_replace(buf, libc.CString("i"), libc.CString(".."))
			search_replace(buf, libc.CString("I"), libc.CString(".."))
			search_replace(buf, libc.CString("g"), libc.CString(".."))
			search_replace(buf, libc.CString("G"), libc.CString(".."))
		} else if rand_number(1, 10) >= 5 {
			search_replace(buf, libc.CString("a"), libc.CString(".."))
			search_replace(buf, libc.CString("A"), libc.CString(".."))
			search_replace(buf, libc.CString("e"), libc.CString(".."))
			search_replace(buf, libc.CString("E"), libc.CString(".."))
			search_replace(buf, libc.CString("i"), libc.CString(".."))
			search_replace(buf, libc.CString("I"), libc.CString(".."))
			search_replace(buf, libc.CString("o"), libc.CString(".."))
			search_replace(buf, libc.CString("O"), libc.CString(".."))
		} else {
			search_replace(buf, libc.CString("k"), libc.CString(".."))
			search_replace(buf, libc.CString("K"), libc.CString(".."))
			search_replace(buf, libc.CString("m"), libc.CString(".."))
			search_replace(buf, libc.CString("M"), libc.CString(".."))
			search_replace(buf, libc.CString("b"), libc.CString(".."))
			search_replace(buf, libc.CString("B"), libc.CString(".."))
		}
		return buf
	case 3:
		return buf
	}
	return libc.CString("Nothing")
}
func do_write(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		paper     *obj_data
		pen       *obj_data = nil
		obj       *obj_data
		papername *byte
		penname   *byte
		buf1      [64936]byte
		buf2      [64936]byte
	)
	for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
		if int(obj.Type_flag) == ITEM_BOARD {
			break
		}
	}
	if obj == nil {
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj.Next_content {
			if int(obj.Type_flag) == ITEM_BOARD {
				break
			}
		}
	}
	if obj != nil {
		write_board_message(GET_OBJ_VNUM(obj), ch, argument)
		act(libc.CString("$n begins to write a note on $p."), TRUE, ch, obj, nil, TO_ROOM)
		return
	}
	papername = &buf1[0]
	penname = &buf2[0]
	two_arguments(argument, papername, penname)
	if ch.Desc == nil {
		return
	}
	if *papername == 0 {
		send_to_char(ch, libc.CString("write on [what] with [what pen?]\r\n"))
		return
	}
	if *penname != 0 {
		if (func() *obj_data {
			paper = get_obj_in_list_vis(ch, papername, nil, ch.Carrying)
			return paper
		}()) == nil {
			send_to_char(ch, libc.CString("You have no %s.\r\n"), papername)
			return
		}
		if (func() *obj_data {
			pen = get_obj_in_list_vis(ch, penname, nil, ch.Carrying)
			return pen
		}()) == nil {
			send_to_char(ch, libc.CString("You have no %s.\r\n"), penname)
			return
		}
	} else {
		if (func() *obj_data {
			paper = get_obj_in_list_vis(ch, papername, nil, ch.Carrying)
			return paper
		}()) == nil {
			send_to_char(ch, libc.CString("There is no %s in your inventory.\r\n"), papername)
			return
		}
		if int(paper.Type_flag) == ITEM_PEN {
			pen = paper
			paper = nil
		} else if int(paper.Type_flag) != ITEM_NOTE {
			send_to_char(ch, libc.CString("That thing has nothing to do with writing.\r\n"))
			return
		}
		if (ch.Equipment[WEAR_WIELD2]) == nil {
			send_to_char(ch, libc.CString("You can't write with %s %s alone.\r\n"), AN(papername), papername)
			return
		}
		if !CAN_SEE_OBJ(ch, ch.Equipment[WEAR_WIELD2]) {
			send_to_char(ch, libc.CString("The stuff in your hand is invisible!  Yeech!!\r\n"))
			return
		}
		if pen != nil {
			paper = ch.Equipment[WEAR_WIELD2]
		} else {
			pen = ch.Equipment[WEAR_WIELD2]
		}
	}
	if int(pen.Type_flag) != ITEM_PEN {
		act(libc.CString("$p is no good for writing with."), FALSE, ch, pen, nil, TO_CHAR)
	} else if int(paper.Type_flag) != ITEM_NOTE {
		act(libc.CString("You can't write on $p."), FALSE, ch, paper, nil, TO_CHAR)
	} else {
		var backstr *byte = nil
		if paper.Action_description != nil {
			backstr = libc.StrDup(paper.Action_description)
			send_to_char(ch, libc.CString("There's something written on it already:\r\n"))
			send_to_char(ch, libc.CString("%s"), paper.Action_description)
		}
		act(libc.CString("$n begins to jot down a note."), TRUE, ch, nil, nil, TO_ROOM)
		paper.Extra_flags[int(ITEM_UNIQUE_SAVE/32)] |= bitvector_t(int32(1 << (int(ITEM_UNIQUE_SAVE % 32))))
		send_editor_help(ch.Desc)
		string_write(ch.Desc, &paper.Action_description, MAX_NOTE_LENGTH, 0, unsafe.Pointer(backstr))
	}
}
func do_page(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d    *descriptor_data
		vict *char_data
		buf2 [2048]byte
		arg  [2048]byte
	)
	half_chop(argument, &arg[0], &buf2[0])
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Monsters can't page.. go away.\r\n"))
	} else if arg[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to page?\r\n"))
	} else {
		var buf [64936]byte
		stdio.Snprintf(&buf[0], int(64936), "\a\a*$n* %s", &buf2[0])
		if libc.StrCaseCmp(&arg[0], libc.CString("all")) == 0 {
			if ADM_FLAGGED(ch, ADM_TELLALL) {
				for d = descriptor_list; d != nil; d = d.Next {
					if d.Connected == CON_PLAYING && d.Character != nil {
						act(&buf[0], FALSE, ch, nil, unsafe.Pointer(d.Character), TO_VICT)
					}
				}
			} else {
				send_to_char(ch, libc.CString("You will never be godly enough to do that!\r\n"))
			}
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}()) != nil {
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			if PRF_FLAGGED(ch, PRF_NOREPEAT) {
				send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
			} else {
				act(&buf[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			}
		} else {
			send_to_char(ch, libc.CString("There is no such person in the game!\r\n"))
		}
	}
}
func do_gen_comm(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i        *descriptor_data
		color_on [24]byte
		buf1     [2048]byte
		buf2     [2048]byte
		msg      *byte
	)
	buf1[0] = '\x00'
	buf2[0] = '\x00'
	var channels [6]int = [6]int{PRF_NOMUSIC, PRF_DEAF, PRF_NOGOSS, PRF_NOAUCT, PRF_NOGRATZ, 0}
	var hist_type [5]int = [5]int{HIST_HOLLER, HIST_SHOUT, HIST_GOSSIP, HIST_AUCTION, HIST_GRATS}
	var com_msgs [5][4]*byte = [5][4]*byte{{libc.CString("You cannot music!!\r\n"), libc.CString("@D[@mMUSIC@D]"), libc.CString("You aren't even on the channel!\r\n"), libc.CString("@[10]")}, {libc.CString("You cannot shout!!\r\n"), libc.CString("shout"), libc.CString("Turn off your noshout flag first!\r\n"), libc.CString("@[9]")}, {libc.CString("You cannot ooc!!\r\n"), libc.CString("@D[@BOOC@D]"), libc.CString("You aren't even on the channel!\r\n"), libc.CString("@[10]")}, {libc.CString("You cannot newbie!!\r\n"), libc.CString("newbie"), libc.CString("You aren't even on the channel!\r\n"), libc.CString("@[11]")}, {libc.CString("You cannot congratulate!\r\n"), libc.CString("congrat"), libc.CString("You aren't even on the channel!\r\n"), libc.CString("@[12]")}}
	if ch.Desc == nil {
		return
	}
	if PLR_FLAGGED(ch, PLR_NOSHOUT) {
		send_to_char(ch, libc.CString("%s"), com_msgs[subcmd][0])
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_SOUNDPROOF) {
		send_to_char(ch, libc.CString("The walls seem to absorb your words.\r\n"))
		return
	}
	if subcmd == SCMD_SHOUT && (ch.Bonuses[BONUS_MUTE]) > 0 {
		send_to_char(ch, libc.CString("You are mute and are incapable of speech.\r\n"))
		return
	}
	skip_spaces(&argument)
	if subcmd == SCMD_GOSSIP && *argument == '*' {
		subcmd = SCMD_GEMOTE
	}
	if subcmd == SCMD_GEMOTE {
		var do_gmote func(ch *char_data, argument *byte, cmd int, subcmd int)
		if *argument == '*' || *argument == ':' {
			do_gmote(ch, (*byte)(unsafe.Add(unsafe.Pointer(argument), 1)), 0, 1)
		} else {
			do_gmote(ch, argument, 0, 1)
		}
		return
	}
	if GET_LEVEL(ch) < config_info.Play.Level_can_shout {
		send_to_char(ch, libc.CString("You must be at least level %d before you can %s.\r\n"), config_info.Play.Level_can_shout, com_msgs[subcmd][1])
		return
	}
	if !IS_NPC(ch) && PRF_FLAGGED(ch, bitvector_t(int32(channels[subcmd]))) {
		send_to_char(ch, libc.CString("%s"), com_msgs[subcmd][2])
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Yes, %s, fine, %s we must, but WHAT???\r\n"), com_msgs[subcmd][1], com_msgs[subcmd][1])
		return
	}
	delete_doubledollar(argument)
	strlcpy(&color_on[0], com_msgs[subcmd][3], uint64(24))
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	} else {
		stdio.Snprintf(&buf2[0], int(2048), "%s@WYou %s@W, '@w%s@W'%s@n\r\n", &color_on[0], com_msgs[subcmd][1], argument, &color_on[0])
		send_to_char(ch, libc.CString("%s"), &buf2[0])
		add_history(ch, &buf2[0], hist_type[subcmd])
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected == CON_PLAYING && i != ch.Desc && i.Character != nil && (IS_NPC(i.Character) || !PRF_FLAGGED(i.Character, bitvector_t(int32(channels[subcmd])))) && (IS_NPC(i.Character) || !PLR_FLAGGED(i.Character, PLR_WRITING)) && !ROOM_FLAGGED(i.Character.In_room, ROOM_SOUNDPROOF) {
			if subcmd == SCMD_SHOUT && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.Character.In_room)))).Zone || !AWAKE(i.Character)) {
				continue
			}
			if config_info.Play.Enable_languages != 0 {
				garble_text(argument, GET_SKILL(i.Character, ch.Player_specials.Speaking), ch.Player_specials.Speaking)
				stdio.Snprintf(&buf1[0], int(2048), "%s%s %ss%s '%s@n'%s", &color_on[0], func() *byte {
					if ch.Admlevel > 0 {
						return GET_NAME(ch)
					}
					return ch.Desc.User
				}(), com_msgs[subcmd][1], func() string {
					if GET_SKILL(i.Character, ch.Player_specials.Speaking) != 0 {
						return ","
					}
					return ", in an unfamiliar tongue,"
				}(), argument, &color_on[0])
			} else if subcmd == SCMD_SHOUT && i.Character.In_room != ch.In_room {
				stdio.Snprintf(&buf1[0], int(2048), "%s@WSomeone nearby %ss@W, '@w%s@W'@n%s", &color_on[0], com_msgs[subcmd][1], argument, &color_on[0])
			} else if subcmd == SCMD_SHOUT && i.Character.In_room == ch.In_room {
				stdio.Snprintf(&buf1[0], int(2048), "%s@W$n@W %ss@W, '@w%s@W'@n%s", &color_on[0], com_msgs[subcmd][1], argument, &color_on[0])
			} else {
				if ch.Admlevel > 0 {
					stdio.Snprintf(&buf1[0], int(2048), "%s@W$n %ss@W, '@w%s@W'@n%s", &color_on[0], com_msgs[subcmd][1], argument, &color_on[0])
				} else if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_HIDE) && i.Character.Admlevel < ADMLVL_IMMORT && ch != i.Character {
					stdio.Snprintf(&buf1[0], int(2048), "%s@WAnonymous Player %ss@W, '@w%s@W'@n%s", &color_on[0], com_msgs[subcmd][1], argument, &color_on[0])
				} else if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_HIDE) && i.Character.Admlevel >= ADMLVL_IMMORT && ch != i.Character {
					stdio.Snprintf(&buf1[0], int(2048), "%s@W%s(H) %ss@W, '@w%s@W'@n%s", &color_on[0], func() *byte {
						if i.Character.Admlevel > 0 {
							return GET_NAME(ch)
						}
						return ch.Desc.User
					}(), com_msgs[subcmd][1], argument, &color_on[0])
				} else if i.Character.Admlevel > 0 {
					stdio.Snprintf(&buf1[0], int(2048), "%s@W%s ($n) %ss@W, '@w%s@W'@n%s", &color_on[0], ch.Desc.User, com_msgs[subcmd][1], argument, &color_on[0])
				} else {
					stdio.Snprintf(&buf1[0], int(2048), "%s@W%s %ss@W, '@w%s@W'@n%s", &color_on[0], func() *byte {
						if i.Character.Admlevel > 0 {
							return GET_NAME(ch)
						}
						return ch.Desc.User
					}(), com_msgs[subcmd][1], argument, &color_on[0])
				}
			}
			msg = act(&buf1[0], FALSE, ch, nil, unsafe.Pointer(i.Character), int(TO_VICT|2<<7))
			add_history(i.Character, msg, hist_type[subcmd])
		}
	}
	if ch.Spam >= 3 && ch.Admlevel < 1 {
		send_to_imm(libc.CString("SPAMMING: %s has been frozen for spamming!\r\n"), GET_NAME(ch))
		send_to_all(libc.CString("@rSPAMMING@D: @C%s@w has been frozen for spamming, let that be a lesson to 'em.@n\r\n"), GET_NAME(ch))
		ch.Act[int(PLR_FROZEN/32)] |= bitvector_t(int32(1 << (int(PLR_FROZEN % 32))))
		ch.Player_specials.Freeze_level = 1
	} else if ch.Spam < 3 {
		ch.Spam += 1
	}
}
func do_qcomm(ch *char_data, argument *byte, cmd int, subcmd int) {
	if !PRF_FLAGGED(ch, PRF_QUEST) {
		send_to_char(ch, libc.CString("You aren't even part of the quest!\r\n"))
		return
	}
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("%c%s?  Yes, fine, %s we must, but WHAT??\r\n"), unicode.ToUpper(rune(*(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command)), (*byte)(unsafe.Add(unsafe.Pointer((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command), 1)), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command)
	} else {
		var (
			buf [64936]byte
			i   *descriptor_data
		)
		if PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else if subcmd == SCMD_QSAY {
			stdio.Snprintf(&buf[0], int(64936), "You quest-say, '%s'", argument)
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(argument), TO_CHAR)
		} else {
			act(argument, FALSE, ch, nil, unsafe.Pointer(argument), TO_CHAR)
		}
		if subcmd == SCMD_QSAY {
			stdio.Snprintf(&buf[0], int(64936), "$n quest-says, '%s'", argument)
		} else {
			strlcpy(&buf[0], argument, uint64(64936))
		}
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Connected == CON_PLAYING && i != ch.Desc && PRF_FLAGGED(i.Character, PRF_QUEST) {
				act(&buf[0], 0, ch, nil, unsafe.Pointer(i.Character), int(TO_VICT|2<<7))
			}
		}
	}
}
func do_respond(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		found  int = 0
		mnum   int = 0
		obj    *obj_data
		number [64936]byte
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("As a mob, you never bothered to learn to read or write.\r\n"))
		return
	}
	for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
		if int(obj.Type_flag) == ITEM_BOARD {
			found = 1
			break
		}
	}
	if obj == nil {
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj.Next_content {
			if int(obj.Type_flag) == ITEM_BOARD {
				found = 1
				break
			}
		}
	}
	if obj != nil {
		argument = one_argument(argument, &number[0])
		if number[0] == 0 {
			send_to_char(ch, libc.CString("Respond to what?\r\n"))
			return
		}
		if !unicode.IsDigit(rune(number[0])) || (func() int {
			mnum = libc.Atoi(libc.GoString(&number[0]))
			return mnum
		}()) == 0 {
			send_to_char(ch, libc.CString("You must type the number of the message you wish to reply to.\r\n"))
			return
		}
		board_respond(int(GET_OBJ_VNUM(obj)), ch, mnum)
	}
	if found == 0 {
		send_to_char(ch, libc.CString("Sorry, you may only reply to messages posted on a board.\r\n"))
	}
}
