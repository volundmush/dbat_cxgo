package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func any_one_name(argument *byte, first_arg *byte) *byte {
	var arg *byte
	for unicode.IsSpace(rune(*argument)) {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	for arg = first_arg; *argument != 0 && !unicode.IsSpace(rune(*argument)) && (!unicode.IsPunct(rune(*argument)) || *argument == '#' || *argument == '-'); func() *byte {
		arg = (*byte)(unsafe.Add(unsafe.Pointer(arg), 1))
		return func() *byte {
			p := &argument
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
	}() {
		*arg = byte(int8(unicode.ToLower(rune(*argument))))
	}
	*arg = '\x00'
	return argument
}
func sub_write_to_char(ch *char_data, tokens []*byte, otokens []unsafe.Pointer, type_ []byte) {
	var (
		sb [64936]byte
		i  int
	)
	libc.StrCpy(&sb[0], libc.CString(""))
	for i = 0; tokens[i+1] != nil; i++ {
		libc.StrCat(&sb[0], tokens[i])
		switch type_[i] {
		case '~':
			if otokens[i] == nil {
				libc.StrCat(&sb[0], libc.CString("someone"))
			} else if (*char_data)(otokens[i]) == ch {
				libc.StrCat(&sb[0], libc.CString("you"))
			} else {
				libc.StrCat(&sb[0], PERS((*char_data)(otokens[i]), ch))
			}
		case '|':
			if otokens[i] == nil {
				libc.StrCat(&sb[0], libc.CString("someone's"))
			} else if (*char_data)(otokens[i]) == ch {
				libc.StrCat(&sb[0], libc.CString("your"))
			} else {
				libc.StrCat(&sb[0], PERS((*char_data)(otokens[i]), ch))
				libc.StrCat(&sb[0], libc.CString("'s"))
			}
		case '^':
			if otokens[i] == nil || !CAN_SEE(ch, (*char_data)(otokens[i])) {
				libc.StrCat(&sb[0], libc.CString("its"))
			} else if otokens[i] == unsafe.Pointer(ch) {
				libc.StrCat(&sb[0], libc.CString("your"))
			} else {
				libc.StrCat(&sb[0], HSHR((*char_data)(otokens[i])))
			}
		case '&':
			if otokens[i] == nil || !CAN_SEE(ch, (*char_data)(otokens[i])) {
				libc.StrCat(&sb[0], libc.CString("it"))
			} else if otokens[i] == unsafe.Pointer(ch) {
				libc.StrCat(&sb[0], libc.CString("you"))
			} else {
				libc.StrCat(&sb[0], HSSH((*char_data)(otokens[i])))
			}
		case '*':
			if otokens[i] == nil || !CAN_SEE(ch, (*char_data)(otokens[i])) {
				libc.StrCat(&sb[0], libc.CString("it"))
			} else if otokens[i] == unsafe.Pointer(ch) {
				libc.StrCat(&sb[0], libc.CString("you"))
			} else {
				libc.StrCat(&sb[0], HMHR((*char_data)(otokens[i])))
			}
		case '\xa8':
			if otokens[i] == nil {
				libc.StrCat(&sb[0], libc.CString("something"))
			} else {
				libc.StrCat(&sb[0], OBJS((*obj_data)(otokens[i]), ch))
			}
		}
	}
	libc.StrCat(&sb[0], tokens[i])
	libc.StrCat(&sb[0], libc.CString("\n\r"))
	sb[0] = byte(int8(unicode.ToUpper(rune(sb[0]))))
	send_to_char(ch, libc.CString("%s"), &sb[0])
}
func sub_write(arg *byte, ch *char_data, find_invis int8, targets int) {
	var (
		str         [4096]byte
		type_       [2048]byte
		name        [2048]byte
		tokens      [2048]*byte
		s           *byte
		p           *byte
		otokens     [2048]unsafe.Pointer
		to          *char_data
		obj         *obj_data
		i           int
		tmp         int
		to_sleeping int = 1
	)
	_ = to_sleeping
	if arg == nil {
		return
	}
	tokens[0] = &str[0]
	for func() *byte {
		i = 0
		p = arg
		return func() *byte {
			s = &str[0]
			return s
		}()
	}(); *p != 0; {
		switch *p {
		case '~':
			fallthrough
		case '|':
			fallthrough
		case '^':
			fallthrough
		case '&':
			fallthrough
		case '*':
			type_[i] = *p
			*s = '\x00'
			p = any_one_name(func() *byte {
				p := &p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}(), &name[0])
			if int(find_invis) != 0 {
				otokens[i] = unsafe.Pointer(get_char_in_room(&world[ch.In_room], &name[0]))
			} else {
				otokens[i] = unsafe.Pointer(get_char_room_vis(ch, &name[0], nil))
			}
			tokens[func() int {
				p := &i
				*p++
				return *p
			}()] = func() *byte {
				p := &s
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()
		case '\xa8':
			type_[i] = *p
			*s = '\x00'
			p = any_one_name(func() *byte {
				p := &p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}(), &name[0])
			if int(find_invis) != 0 {
				obj = get_obj_in_room(&world[ch.In_room], &name[0])
			} else if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &name[0], nil, world[ch.In_room].Contents)
				return obj
			}()) == nil {
			} else if (func() *obj_data {
				obj = get_obj_in_equip_vis(ch, &name[0], &tmp, ch.Equipment[:])
				return obj
			}()) == nil {
			} else {
				obj = get_obj_in_list_vis(ch, &name[0], nil, ch.Carrying)
			}
			otokens[i] = unsafe.Pointer(obj)
			tokens[func() int {
				p := &i
				*p++
				return *p
			}()] = func() *byte {
				p := &s
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()
		case '\\':
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
			*func() *byte {
				p := &s
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		default:
			*func() *byte {
				p := &s
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		}
	}
	*s = '\x00'
	tokens[func() int {
		p := &i
		*p++
		return *p
	}()] = nil
	if IS_SET(bitvector_t(int32(targets)), TO_CHAR) && SENDOK(ch) {
		sub_write_to_char(ch, tokens[:], otokens[:], type_[:])
	}
	if IS_SET(bitvector_t(int32(targets)), TO_ROOM) {
		for to = world[ch.In_room].People; to != nil; to = to.Next_in_room {
			if to != ch && SENDOK(to) {
				sub_write_to_char(to, tokens[:], otokens[:], type_[:])
			}
		}
	}
}
func send_to_zone(messg *byte, zone zone_rnum) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected == 0 && i.Character != nil && AWAKE(i.Character) && i.Character.In_room != room_rnum(-1) && world[i.Character.In_room].Zone == zone {
			write_to_output(i, libc.CString("%s"), messg)
		}
	}
}
func fly_zone(zone zone_rnum, messg *byte, ch *char_data) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected == 0 && i.Character != nil && AWAKE(i.Character) && OUTSIDE(i.Character) && i.Character.In_room != room_rnum(-1) && world[i.Character.In_room].Zone == zone && i.Character != ch {
			if PLR_FLAGGED(i.Character, PLR_DISGUISED) {
				write_to_output(i, libc.CString("A disguised figure %s"), messg)
			} else {
				write_to_output(i, libc.CString("%s%s %s"), func() string {
					if readIntro(i.Character, ch) == 1 {
						return ""
					}
					return "A "
				}(), get_i_name(i.Character, ch), messg)
			}
		}
	}
}
func send_to_sense(type_ int, messg *byte, ch *char_data) {
	var (
		i   *descriptor_data
		tch *char_data
		obj *obj_data
	)
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		tch = i.Character
		obj = tch.Equipment[WEAR_EYE]
		if tch == ch {
			continue
		}
		if GET_SKILL(tch, SKILL_SENSE) == 0 {
			continue
		}
		if world[ch.In_room].Zone != world[tch.In_room].Zone && type_ == 0 || !AWAKE(tch) {
			continue
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_SHIP) {
			continue
		}
		if obj != nil && type_ == 0 {
			continue
		}
		if int(ch.Race) == RACE_ANDROID {
			continue
		}
		if ch.In_room == tch.In_room {
			continue
		} else if float64(ch.Hit) < (float64(tch.Hit)*0.001)+1 {
			continue
		} else if type_ == 0 {
			if ch.Max_hit > tch.Max_hit {
				write_to_output(i, libc.CString("%s who is stronger than you. They are nearby.\r\n"), messg)
			} else if float64(ch.Max_hit) >= float64(tch.Max_hit)*0.9 {
				write_to_output(i, libc.CString("%s who is near your strength. They are nearby.\r\n"), messg)
			} else if float64(ch.Max_hit) >= float64(tch.Max_hit)*0.6 {
				write_to_output(i, libc.CString("%s who is a good bit weaker than you. They are nearby.\r\n"), messg)
			} else if float64(ch.Max_hit) >= float64(tch.Max_hit)*0.4 {
				write_to_output(i, libc.CString("%s who is a lot weaker than you. They are nearby.\r\n"), messg)
			} else {
				continue
			}
			if readIntro(tch, ch) == 1 {
				write_to_output(i, libc.CString("@YYou recognise this signal as @y%s@Y!@n\r\n"), get_i_name(tch, ch))
			} else if read_sense_memory(ch, tch) != 0 {
				write_to_output(i, libc.CString("@YYou recognise this signal, but don't seem to know their name.@n\r\n"))
			}
		} else if planet_check(ch, tch) != 0 {
			var (
				blah  *byte = sense_location(ch)
				power [2048]byte
				align [2048]byte
			)
			if ch.Hit > tch.Hit*10 {
				stdio.Sprintf(&power[0], ", who is @Runbelievably stronger@Y than you")
			} else if ch.Hit > tch.Hit*5 {
				stdio.Sprintf(&power[0], ", who is much @Rstronger@Y than you")
			} else if ch.Hit > tch.Hit*2 {
				stdio.Sprintf(&power[0], ", who is more than twice as @Rstrong@Y as you")
			} else if ch.Hit > tch.Hit {
				stdio.Sprintf(&power[0], ", who is somewhat @mstronger@Y than you")
			} else if ch.Hit*10 < tch.Hit {
				stdio.Sprintf(&power[0], ", who is @Munbelievably weaker@Y than you")
			} else if ch.Hit*5 < tch.Hit {
				stdio.Sprintf(&power[0], ", who is much @Mweaker@Y than you")
			} else if ch.Hit*2 < tch.Hit {
				stdio.Sprintf(&power[0], ", who is more than twice as @Mweak@Y as you")
			} else if ch.Hit < tch.Hit {
				stdio.Sprintf(&power[0], ", who is somewhat @Wweaker@Y than you")
			} else {
				stdio.Sprintf(&power[0], ", who is close to @Cequal@Y with you")
			}
			if ch.Alignment >= 1000 {
				stdio.Sprintf(&align[0], ", with a @wsaintly@Y aura,")
			} else if ch.Alignment >= 500 {
				stdio.Sprintf(&align[0], ", with a very @Cgood@Y aura,")
			} else if ch.Alignment >= 200 {
				stdio.Sprintf(&align[0], ", with a @cgood@Y aura,")
			} else if ch.Alignment > -100 {
				stdio.Sprintf(&align[0], ", with a near @Wneutral@Y aura,")
			} else if ch.Alignment > -200 {
				stdio.Sprintf(&align[0], ", with a sorta @revil@Y aura,")
			} else if ch.Alignment > -500 {
				stdio.Sprintf(&align[0], ", with an @revil@Y aura,")
			} else if ch.Alignment > -900 {
				stdio.Sprintf(&align[0], ", with a @rvery evil@Y aura,")
			} else {
				stdio.Sprintf(&align[0], ", with a @rd@De@Wv@wil@Wi@Ds@rh@Y aura,")
			}
			if libc.StrStr(messg, libc.CString("land")) != nil {
				write_to_output(i, libc.CString("@YYou sense %s%s%s %s! They appear to have landed at...@G%s@n\r\n"), func() *byte {
					if readIntro(tch, ch) == 1 {
						return get_i_name(tch, ch)
					}
					return libc.CString("someone")
				}(), &power[0], &align[0], messg, blah)
			} else {
				write_to_output(i, libc.CString("@YYou sense %s%s%s %s!@n\r\n"), func() *byte {
					if readIntro(tch, ch) == 1 {
						return get_i_name(tch, ch)
					}
					return libc.CString("someone")
				}(), &power[0], &align[0], messg)
			}
		}
	}
}
func send_to_scouter(messg *byte, ch *char_data, num int, type_ int) {
	var (
		i   *descriptor_data
		tch *char_data
		obj *obj_data
	)
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		tch = i.Character
		obj = tch.Equipment[WEAR_EYE]
		if tch == ch {
			continue
		} else {
			if world[ch.In_room].Zone != world[tch.In_room].Zone && type_ == 0 || !AWAKE(tch) {
				continue
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_SHIP) {
				continue
			}
			if int(ch.Player_specials.Invis_level) > tch.Admlevel {
				continue
			}
			if int(ch.Race) == RACE_ANDROID {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_EARTH) && !ROOM_FLAGGED(tch.In_room, ROOM_EARTH) {
				continue
			} else if PLANET_ZENITH(ch.In_room) && !PLANET_ZENITH(tch.In_room) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) && !ROOM_FLAGGED(tch.In_room, ROOM_FRIGID) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) && !ROOM_FLAGGED(tch.In_room, ROOM_NAMEK) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_AL) && !ROOM_FLAGGED(tch.In_room, ROOM_AL) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) && !ROOM_FLAGGED(tch.In_room, ROOM_VEGETA) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_KONACK) && !ROOM_FLAGGED(tch.In_room, ROOM_KONACK) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_NEO) && !ROOM_FLAGGED(tch.In_room, ROOM_NEO) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_YARDRAT) && !ROOM_FLAGGED(tch.In_room, ROOM_YARDRAT) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) && !ROOM_FLAGGED(tch.In_room, ROOM_KANASSA) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_ARLIA) && !ROOM_FLAGGED(tch.In_room, ROOM_ARLIA) {
				continue
			} else if ROOM_FLAGGED(ch.In_room, ROOM_AETHER) && !ROOM_FLAGGED(tch.In_room, ROOM_AETHER) {
				continue
			}
			if obj == nil {
				continue
			} else if ch.In_room == tch.In_room {
				continue
			} else if type_ == 0 {
				if num == 1 {
					var obj *obj_data = (tch.Equipment[WEAR_EYE])
					if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && ch.Hit >= 150000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Rising Powerlevel Detected@D:@Y ??????????\r\n"))
					} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && ch.Hit >= 5000000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Rising Powerlevel Detected@D:@Y ??????????\r\n"))
					} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && ch.Hit >= 15000000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Rising Powerlevel Detected@D:@Y ??????????\r\n"))
					} else {
						write_to_output(i, libc.CString("%s@n"), messg)
					}
				} else {
					if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && ch.Hit >= 150000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Nearby Powerlevel Detected@D:@Y ??????????\r\n"))
					} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && ch.Hit >= 5000000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Nearby Powerlevel Detected@D:@Y ??????????\r\n"))
					} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && ch.Hit >= 15000000 {
						write_to_output(i, libc.CString("@D[@GBlip@D]@r Nearby Powerlevel Detected@D:@Y ??????????\r\n"))
					} else {
						write_to_output(i, libc.CString("%s\r\n"), messg)
					}
				}
			} else if type_ == 1 && GET_SKILL(tch, SKILL_SENSE) < 20 {
				if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && ch.Hit >= 150000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s. @RPL@D:@Y ??????????\r\n"), messg)
				} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && ch.Hit >= 5000000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s. @RPL@D:@Y ??????????\r\n"), messg)
				} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && ch.Hit >= 15000000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s. @RPL@D:@Y ??????????\r\n"), messg)
				} else {
					write_to_output(i, libc.CString("@D[Blip@D]@w %s. @RPL@D:@Y %s@n\r\n\r\n"), messg, add_commas(ch.Hit))
				}
			} else if type_ == 2 && GET_SKILL(tch, SKILL_SENSE) < 20 {
				var blah *byte = sense_location(ch)
				if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && ch.Hit >= 150000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s at... @G%s. @RPL@D:@Y ??????????\r\n"), messg, blah)
				} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && ch.Hit >= 5000000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s at... @G%s. @RPL@D:@Y ??????????\r\n"), messg, blah)
				} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && ch.Hit >= 15000000 {
					write_to_output(i, libc.CString("@D[@GBlip@D]@w %s at... @G%s. @RPL@D:@Y ??????????\r\n"), messg, blah)
				} else {
					write_to_output(i, libc.CString("@D[Blip@D]@w %s at... @G%s. @RPL@D:@Y %s@n\r\n\r\n"), messg, blah, add_commas(ch.Hit))
				}
			}
		}
	}
}
func send_to_worlds(ch *char_data) {
	var (
		i       *descriptor_data
		message [2048]byte
	)
	if ch.Max_hit > 2000000000 {
		stdio.Sprintf(&message[0], "@RThe whole planet begins to quake violently as if the world is ending!@n\r\n")
	} else if ch.Max_hit > 1000000000 {
		stdio.Sprintf(&message[0], "@RThe whole planet begins to quake violently with a thunderous roar!@n\r\n")
	} else if ch.Max_hit > 500000000 {
		stdio.Sprintf(&message[0], "@RThe whole planet begins to quake violently!@n\r\n")
	} else if ch.Max_hit > 100000000 {
		stdio.Sprintf(&message[0], "@RThe whole planet rumbles and shakes!@n\r\n")
	} else if ch.Max_hit > 50000000 {
		stdio.Sprintf(&message[0], "@RThe whole planet rumbles faintly!@n\r\n")
	} else {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		if ROOM_FLAGGED(i.Character.In_room, ROOM_EARTH) && ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_VEGETA) && ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if PLANET_ZENITH(i.Character.In_room) && PLANET_ZENITH(ch.In_room) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_NAMEK) && ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_KONACK) && ROOM_FLAGGED(ch.In_room, ROOM_KONACK) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_YARDRAT) && ROOM_FLAGGED(ch.In_room, ROOM_YARDRAT) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_FRIGID) && ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_KANASSA) && ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_ARLIA) && ROOM_FLAGGED(ch.In_room, ROOM_ARLIA) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		} else if ROOM_FLAGGED(i.Character.In_room, ROOM_AETHER) && ROOM_FLAGGED(ch.In_room, ROOM_AETHER) {
			send_to_char(i.Character, libc.CString("%s"), &message[0])
		}
	}
}
func send_to_imm(messg *byte, _rest ...interface{}) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		} else if i.Character.Admlevel == 0 {
			continue
		} else if !PRF_FLAGGED(i.Character, PRF_LOG2) {
			continue
		} else if PLR_FLAGGED(i.Character, PLR_WRITING) {
			continue
		} else {
			write_to_output(i, libc.CString("@g[ Log: "))
			var args libc.ArgList
			args.Start(messg, _rest)
			vwrite_to_output(i, messg, args)
			write_to_output(i, libc.CString(" ]@n\n"))
			args.End()
		}
	}
	var args libc.ArgList
	args.Start(messg, _rest)
	basic_mud_vlog(messg, args)
	args.End()
}
