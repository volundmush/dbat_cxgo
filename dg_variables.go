package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func add_var(var_list **trig_var_data, name *byte, value *byte, id int) {
	var vd *trig_var_data
	if libc.StrChr(name, '.') != nil {
		basic_mud_log(libc.CString("add_var() : Attempt to add illegal var: %s"), name)
		return
	}
	for vd = *var_list; vd != nil && libc.StrCaseCmp(vd.Name, name) != 0; vd = vd.Next {
	}
	if vd != nil && (vd.Context == 0 || vd.Context == id) {
		libc.Free(unsafe.Pointer(vd.Value))
		vd.Value = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(value)+1)[0]))
	} else {
		vd = new(trig_var_data)
		vd.Name = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(name)+1)[0]))
		libc.StrCpy(vd.Name, name)
		vd.Value = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(value)+1)[0]))
		vd.Next = *var_list
		vd.Context = id
		*var_list = vd
	}
	libc.StrCpy(vd.Value, value)
}
func skill_percent(ch *char_data, skill *byte) *byte {
	var (
		retval   [16]byte
		skillnum int
	)
	skillnum = find_skill_num(skill, 1<<1)
	if skillnum <= 0 {
		return libc.CString("unknown skill")
	}
	stdio.Snprintf(&retval[0], int(16), "%d", GET_SKILL(ch, skillnum))
	return &retval[0]
}
func item_in_list(item *byte, list *obj_data) int {
	var (
		i     *obj_data
		count int = 0
	)
	if item == nil || *item == 0 {
		return 0
	}
	if *item == UID_CHAR {
		var id int = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(item), 1))))
		for i = list; i != nil; i = i.Next_content {
			if id == int(i.Id) {
				count++
			}
			if int(i.Type_flag) == ITEM_CONTAINER {
				count += item_in_list(item, i.Contains)
			}
		}
	} else if is_number(item) > -1 {
		var ovnum obj_vnum = obj_vnum(libc.Atoi(libc.GoString(item)))
		for i = list; i != nil; i = i.Next_content {
			if GET_OBJ_VNUM(i) == ovnum {
				count++
			}
			if int(i.Type_flag) == ITEM_CONTAINER {
				count += item_in_list(item, i.Contains)
			}
		}
	} else {
		for i = list; i != nil; i = i.Next_content {
			if isname(item, i.Name) != 0 {
				count++
			}
			if int(i.Type_flag) == ITEM_CONTAINER {
				count += item_in_list(item, i.Contains)
			}
		}
	}
	return count
}
func char_has_item(item *byte, ch *char_data) int {
	if get_object_in_equip(ch, item) != nil {
		return 1
	}
	if item_in_list(item, ch.Carrying) == 0 {
		return 0
	} else {
		return 1
	}
}
func text_processed(field *byte, subfield *byte, vd *trig_var_data, str *byte, slen uint64) int {
	var (
		p      *byte
		p2     *byte
		tmpvar [64936]byte
	)
	if libc.StrCaseCmp(field, libc.CString("strlen")) == 0 {
		var limit [200]byte
		stdio.Sprintf(&limit[0], "%lld", libc.StrLen(vd.Value))
		stdio.Snprintf(str, int(slen), "%d", libc.Atoi(libc.GoString(&limit[0])))
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("trim")) == 0 {
		stdio.Snprintf(&tmpvar[0], int(64936-1), "%s", vd.Value)
		p = &tmpvar[0]
		p2 = (*byte)(unsafe.Add(unsafe.Pointer(&tmpvar[libc.StrLen(&tmpvar[0])]), -1))
		for *p != 0 && unicode.IsSpace(rune(*p)) {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
		for uintptr(unsafe.Pointer(p)) <= uintptr(unsafe.Pointer(p2)) && unicode.IsSpace(rune(*p2)) {
			p2 = (*byte)(unsafe.Add(unsafe.Pointer(p2), -1))
		}
		if uintptr(unsafe.Pointer(p)) > uintptr(unsafe.Pointer(p2)) {
			*str = '\x00'
			return TRUE
		}
		*(func() *byte {
			p := &p2
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return *p
		}()) = '\x00'
		stdio.Snprintf(str, int(slen), "%s", p)
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("contains")) == 0 {
		if str_str(vd.Value, subfield) != nil {
			libc.StrCpy(str, libc.CString("1"))
		} else {
			libc.StrCpy(str, libc.CString("0"))
		}
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("car")) == 0 {
		var car *byte = vd.Value
		for *car != 0 && !unicode.IsSpace(rune(*car)) {
			*func() *byte {
				p := &str
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &car
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		}
		*str = '\x00'
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("cdr")) == 0 {
		var cdr *byte = vd.Value
		for *cdr != 0 && !unicode.IsSpace(rune(*cdr)) {
			cdr = (*byte)(unsafe.Add(unsafe.Pointer(cdr), 1))
		}
		for *cdr != 0 && unicode.IsSpace(rune(*cdr)) {
			cdr = (*byte)(unsafe.Add(unsafe.Pointer(cdr), 1))
		}
		stdio.Snprintf(str, int(slen), "%s", cdr)
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("charat")) == 0 {
		var (
			len_    uint64 = uint64(libc.StrLen(vd.Value))
			dgindex uint64 = uint64(libc.Atoi(libc.GoString(subfield)))
		)
		if dgindex > len_ || dgindex < 1 {
			libc.StrCpy(str, libc.CString(""))
		} else {
			stdio.Snprintf(str, int(slen), "%c", *(*byte)(unsafe.Add(unsafe.Pointer(vd.Value), dgindex-1)))
		}
		return TRUE
	} else if libc.StrCaseCmp(field, libc.CString("mudcommand")) == 0 {
		var (
			length int
			cmd    int
		)
		for func() int {
			length = libc.StrLen(vd.Value)
			return func() int {
				cmd = 0
				return cmd
			}()
		}(); *cmd_info[cmd].Command != '\n'; cmd++ {
			if libc.StrNCmp(cmd_info[cmd].Command, vd.Value, length) == 0 {
				break
			}
		}
		if *cmd_info[cmd].Command == '\n' {
			*str = '\x00'
		} else {
			stdio.Snprintf(str, int(slen), "%s", cmd_info[cmd].Command)
		}
		return TRUE
	}
	return FALSE
}
func find_replacement(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, var_ *byte, field *byte, subfield *byte, str *byte, slen uint64) {
	var (
		vd             *trig_var_data = nil
		ch             *char_data
		c              *char_data = nil
		rndm           *char_data
		obj            *obj_data
		o              *obj_data = nil
		room           *room_data
		r              *room_data = nil
		name           *byte
		num            int
		count          int
		i              int
		j              int
		doors          int
		send_cmd       [3]*byte = [3]*byte{libc.CString("msend "), libc.CString("osend "), libc.CString("wsend ")}
		echo_cmd       [3]*byte = [3]*byte{libc.CString("mecho "), libc.CString("oecho "), libc.CString("wecho ")}
		echoaround_cmd [3]*byte = [3]*byte{libc.CString("mechoaround "), libc.CString("oechoaround "), libc.CString("wechoaround ")}
		door           [3]*byte = [3]*byte{libc.CString("mdoor "), libc.CString("odoor "), libc.CString("wdoor ")}
		force          [3]*byte = [3]*byte{libc.CString("mforce "), libc.CString("oforce "), libc.CString("wforce ")}
		load           [3]*byte = [3]*byte{libc.CString("mload "), libc.CString("oload "), libc.CString("wload ")}
		purge          [3]*byte = [3]*byte{libc.CString("mpurge "), libc.CString("opurge "), libc.CString("wpurge ")}
		teleport       [3]*byte = [3]*byte{libc.CString("mteleport "), libc.CString("oteleport "), libc.CString("wteleport ")}
		xdamage        [3]*byte = [3]*byte{libc.CString("mdamage "), libc.CString("odamage "), libc.CString("wdamage ")}
		zoneecho       [3]*byte = [3]*byte{libc.CString("mzoneecho "), libc.CString("ozoneecho "), libc.CString("wzoneecho ")}
		asound         [3]*byte = [3]*byte{libc.CString("masound "), libc.CString("oasound "), libc.CString("wasound ")}
		at             [3]*byte = [3]*byte{libc.CString("mat "), libc.CString("oat "), libc.CString("wat ")}
		transform      [3]*byte = [3]*byte{libc.CString("mtransform "), libc.CString("otransform "), libc.CString("wecho ")}
		recho          [3]*byte = [3]*byte{libc.CString("mrecho "), libc.CString("orecho "), libc.CString("wrecho ")}
	)
	*str = '\x00'
	if trig != nil {
		for vd = trig.Var_list; vd != nil; vd = vd.Next {
			if libc.StrCaseCmp(vd.Name, var_) == 0 {
				break
			}
		}
	}
	if vd == nil && sc != nil {
		for vd = sc.Global_vars; vd != nil; vd = vd.Next {
			if libc.StrCaseCmp(vd.Name, var_) == 0 && (vd.Context == 0 || vd.Context == sc.Context) {
				break
			}
		}
	}
	if *field == 0 {
		if vd != nil {
			stdio.Snprintf(str, int(slen), "%s", vd.Value)
		} else {
			if libc.StrCaseCmp(var_, libc.CString("self")) == 0 {
				switch type_ {
				case MOB_TRIGGER:
					stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, ((*char_data)(gohere)).Id)
				case OBJ_TRIGGER:
					stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, ((*obj_data)(gohere)).Id)
				case WLD_TRIGGER:
					stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, ((*room_data)(gohere)).Number+ROOM_ID_BASE)
				}
			} else if libc.StrCaseCmp(var_, libc.CString("global")) == 0 {
				stdio.Snprintf(str, int(slen), "%d", ROOM_ID_BASE)
				return
			} else if libc.StrCaseCmp(var_, libc.CString("ctime")) == 0 {
				stdio.Snprintf(str, int(slen), "%ld", libc.GetTime(nil))
			} else if libc.StrCaseCmp(var_, libc.CString("door")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", door[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("force")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", force[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("load")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", load[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("purge")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", purge[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("teleport")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", teleport[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("damage")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", xdamage[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("send")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", send_cmd[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("echo")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", echo_cmd[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("echoaround")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", echoaround_cmd[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("zoneecho")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", zoneecho[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("asound")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", asound[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("at")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", at[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("transform")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", transform[type_])
			} else if libc.StrCaseCmp(var_, libc.CString("recho")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", recho[type_])
			} else {
				*str = '\x00'
			}
		}
		return
	} else {
		if vd != nil {
			name = vd.Value
			switch type_ {
			case MOB_TRIGGER:
				ch = (*char_data)(gohere)
				if (func() *obj_data {
					o = get_object_in_equip(ch, name)
					return o
				}()) != nil {
				} else if (func() *obj_data {
					o = get_obj_in_list(name, ch.Carrying)
					return o
				}()) != nil {
				} else if ch.In_room != room_rnum(-1) && (func() *char_data {
					c = get_char_in_room(&world[ch.In_room], name)
					return c
				}()) != nil {
				} else if (func() *obj_data {
					o = get_obj_in_list(name, world[ch.In_room].Contents)
					return o
				}()) != nil {
				} else if (func() *char_data {
					c = get_char(name)
					return c
				}()) != nil {
				} else if (func() *obj_data {
					o = get_obj(name)
					return o
				}()) != nil {
				} else if (func() *room_data {
					r = get_room(name)
					return r
				}()) != nil {
				}
			case OBJ_TRIGGER:
				obj = (*obj_data)(gohere)
				if (func() *char_data {
					c = get_char_by_obj(obj, name)
					return c
				}()) != nil {
				} else if (func() *obj_data {
					o = get_obj_by_obj(obj, name)
					return o
				}()) != nil {
				} else if (func() *room_data {
					r = get_room(name)
					return r
				}()) != nil {
				}
			case WLD_TRIGGER:
				room = (*room_data)(gohere)
				if (func() *char_data {
					c = get_char_by_room(room, name)
					return c
				}()) != nil {
				} else if (func() *obj_data {
					o = get_obj_by_room(room, name)
					return o
				}()) != nil {
				} else if (func() *room_data {
					r = get_room(name)
					return r
				}()) != nil {
				}
			}
		} else {
			if libc.StrCaseCmp(var_, libc.CString("self")) == 0 {
				switch type_ {
				case MOB_TRIGGER:
					c = (*char_data)(gohere)
					r = nil
					o = nil
				case OBJ_TRIGGER:
					o = (*obj_data)(gohere)
					c = nil
					r = nil
				case WLD_TRIGGER:
					r = (*room_data)(gohere)
					c = nil
					o = nil
				}
			} else if libc.StrCaseCmp(var_, libc.CString("global")) == 0 {
				var thescript *script_data = world[0].Script
				*str = '\x00'
				if thescript == nil {
					script_log(libc.CString("Attempt to find global var. Apparently the void has no script."))
					return
				}
				for vd = thescript.Global_vars; vd != nil; vd = vd.Next {
					if libc.StrCaseCmp(vd.Name, field) == 0 {
						break
					}
				}
				if vd != nil {
					stdio.Snprintf(str, int(slen), "%s", vd.Value)
				}
				return
			} else if libc.StrCaseCmp(var_, libc.CString("people")) == 0 {
				stdio.Snprintf(str, int(slen), "%d", func() int {
					if (func() int {
						num = libc.Atoi(libc.GoString(field))
						return num
					}()) > 0 {
						return trgvar_in_room(room_vnum(num))
					}
					return 0
				}())
				return
			} else if libc.StrCaseCmp(var_, libc.CString("time")) == 0 {
				if libc.StrCaseCmp(field, libc.CString("hour")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", time_info.Hours)
				} else if libc.StrCaseCmp(field, libc.CString("day")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", time_info.Day+1)
				} else if libc.StrCaseCmp(field, libc.CString("month")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", time_info.Month+1)
				} else if libc.StrCaseCmp(field, libc.CString("year")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", time_info.Year)
				} else {
					*str = '\x00'
				}
				return
			} else if libc.StrCaseCmp(var_, libc.CString("findmob")) == 0 {
				if field == nil || *field == 0 || subfield == nil || *subfield == 0 {
					script_log(libc.CString("findmob.vnum(mvnum) - illegal syntax"))
					libc.StrCpy(str, libc.CString("0"))
				} else {
					var (
						rrnum room_rnum = real_room(room_vnum(libc.Atoi(libc.GoString(field))))
						mvnum mob_vnum  = mob_vnum(libc.Atoi(libc.GoString(subfield)))
					)
					if rrnum == room_rnum(-1) {
						script_log(libc.CString("findmob.vnum(ovnum): No room with vnum %d"), libc.Atoi(libc.GoString(field)))
						libc.StrCpy(str, libc.CString("0"))
					} else {
						for func() *char_data {
							i = 0
							return func() *char_data {
								ch = world[rrnum].People
								return ch
							}()
						}(); ch != nil; ch = ch.Next_in_room {
							if GET_MOB_VNUM(ch) == mvnum {
								i++
							}
						}
						stdio.Snprintf(str, int(slen), "%d", i)
					}
				}
			} else if libc.StrCaseCmp(var_, libc.CString("findobj")) == 0 {
				if field == nil || *field == 0 || subfield == nil || *subfield == 0 {
					script_log(libc.CString("findobj.vnum(ovnum) - illegal syntax"))
					libc.StrCpy(str, libc.CString("0"))
				} else {
					var rrnum room_rnum = real_room(room_vnum(libc.Atoi(libc.GoString(field))))
					if rrnum == room_rnum(-1) {
						script_log(libc.CString("findobj.vnum(ovnum): No room with vnum %d"), libc.Atoi(libc.GoString(field)))
						libc.StrCpy(str, libc.CString("0"))
					} else {
						stdio.Snprintf(str, int(slen), "%d", item_in_list(subfield, world[rrnum].Contents))
					}
				}
			} else if libc.StrCaseCmp(var_, libc.CString("random")) == 0 {
				if libc.StrCaseCmp(field, libc.CString("char")) == 0 {
					rndm = nil
					count = 0
					if type_ == MOB_TRIGGER {
						ch = (*char_data)(gohere)
						for c = world[ch.In_room].People; c != nil; c = c.Next_in_room {
							if c != ch && valid_dg_target(c, 1<<0) != 0 && CAN_SEE(ch, c) {
								if rand_number(0, count) == 0 {
									rndm = c
								}
								count++
							}
						}
					} else if type_ == OBJ_TRIGGER {
						for c = world[obj_room((*obj_data)(gohere))].People; c != nil; c = c.Next_in_room {
							if valid_dg_target(c, 1<<0) != 0 {
								if rand_number(0, count) == 0 {
									rndm = c
								}
								count++
							}
						}
					} else if type_ == WLD_TRIGGER {
						for c = ((*room_data)(gohere)).People; c != nil; c = c.Next_in_room {
							if valid_dg_target(c, 1<<0) != 0 {
								if rand_number(0, count) == 0 {
									rndm = c
								}
								count++
							}
						}
					}
					if rndm != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, rndm.Id)
					} else {
						*str = '\x00'
					}
				} else if libc.StrCaseCmp(field, libc.CString("dir")) == 0 {
					var in_room room_rnum = room_rnum(-1)
					switch type_ {
					case WLD_TRIGGER:
						in_room = real_room(((*room_data)(gohere)).Number)
					case OBJ_TRIGGER:
						in_room = obj_room((*obj_data)(gohere))
					case MOB_TRIGGER:
						in_room = ((*char_data)(gohere)).In_room
					}
					if in_room == room_rnum(-1) {
						*str = '\x00'
					} else {
						doors = 0
						room = &world[in_room]
						for i = 0; i < NUM_OF_DIRS; i++ {
							if (room.Dir_option[i]) != nil {
								doors++
							}
						}
						if doors == 0 {
							*str = '\x00'
						} else {
							for {
								doors = rand_number(0, int(NUM_OF_DIRS-1))
								if (room.Dir_option[doors]) != nil {
									break
								}
							}
							stdio.Snprintf(str, int(slen), "%s", dirs[doors])
						}
					}
				} else {
					stdio.Snprintf(str, int(slen), "%d", func() int {
						if (func() int {
							num = libc.Atoi(libc.GoString(field))
							return num
						}()) > 0 {
							return rand_number(1, num)
						}
						return 0
					}())
				}
				return
			}
		}
		if c != nil {
			if text_processed(field, subfield, vd, str, slen) != 0 {
				return
			} else if libc.StrCaseCmp(field, libc.CString("global")) == 0 {
				if IS_NPC(c) && c.Script != nil {
					find_replacement(gohere, c.Script, nil, MOB_TRIGGER, subfield, nil, nil, str, slen)
				}
			}
			*str = '\x01'
			switch unicode.ToLower(rune(*field)) {
			case 'a':
				if libc.StrCaseCmp(field, libc.CString("aaaaa")) == 0 {
					libc.StrCpy(str, libc.CString("0"))
				} else if libc.StrCaseCmp(field, libc.CString("affect")) == 0 {
					if subfield != nil && *subfield != 0 {
						var affect int = get_flag_by_name(affected_bits[:], subfield)
						if affect != int(-1) && AFF_FLAGGED(c, bitvector_t(int32(affect))) {
							libc.StrCpy(str, libc.CString("1"))
						} else {
							libc.StrCpy(str, libc.CString("0"))
						}
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("alias")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", c.Name)
				} else if libc.StrCaseCmp(field, libc.CString("align")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Alignment = int(MAX(-1000, MIN(int64(addition), 1000)))
					}
					stdio.Snprintf(str, int(slen), "%d", c.Alignment)
				}
			case 'b':
				if libc.StrCaseCmp(field, libc.CString("bank")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Bank_gold += addition
					}
					stdio.Snprintf(str, int(slen), "%d", c.Gold)
				}
			case 'c':
				if libc.StrCaseCmp(field, libc.CString("canbeseen")) == 0 {
					if type_ == MOB_TRIGGER && !CAN_SEE((*char_data)(gohere), c) {
						libc.StrCpy(str, libc.CString("0"))
					} else {
						libc.StrCpy(str, libc.CString("1"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("carry")) == 0 {
					if !IS_NPC(c) && c.Player_specials.Carrying != nil {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("clan")) == 0 {
					if c.Clan != nil && libc.StrStr(c.Clan, subfield) != nil {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("class")) == 0 {
					if !IS_NPC(c) {
						stdio.Snprintf(str, int(slen), "%s", pc_class_types[int(c.Chclass)])
					} else {
						stdio.Snprintf(str, int(slen), "blank")
					}
				} else if libc.StrCaseCmp(field, libc.CString("con")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Con += int8(addition)
						if int(c.Aff_abils.Con) > max {
							c.Aff_abils.Con = int8(max)
						}
						if int(c.Aff_abils.Con) < 3 {
							c.Aff_abils.Con = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Con)
				} else if libc.StrCaseCmp(field, libc.CString("cha")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Cha += int8(addition)
						if int(c.Aff_abils.Cha) > max {
							c.Aff_abils.Cha = int8(max)
						}
						if int(c.Aff_abils.Cha) < 3 {
							c.Aff_abils.Cha = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Cha)
				}
			case 'd':
				if libc.StrCaseCmp(field, libc.CString("dead")) == 0 {
					if AFF_FLAGGED(c, AFF_SPIRIT) {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("death")) == 0 {
					stdio.Snprintf(str, int(slen), "%ld", c.Deathtime)
				} else if libc.StrCaseCmp(field, libc.CString("dex")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Dex += int8(addition)
						if int(c.Aff_abils.Dex) > max {
							c.Aff_abils.Dex = int8(max)
						}
						if int(c.Aff_abils.Dex) < 3 {
							c.Aff_abils.Dex = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Dex)
				} else if libc.StrCaseCmp(field, libc.CString("drag")) == 0 {
					if !IS_NPC(c) && c.Drag != nil {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("drunk")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Player_specials.Conditions[DRUNK] = int8(MAX(-1, MIN(int64(addition), 24)))
					}
					stdio.Snprintf(str, int(slen), "%d", c.Player_specials.Conditions[DRUNK])
				}
			case 'e':
				if libc.StrCaseCmp(field, libc.CString("eq")) == 0 {
					var pos int
					if subfield == nil || *subfield == 0 {
						*str = '\x00'
					} else if *subfield == '*' {
						for func() int {
							i = 0
							return func() int {
								j = 0
								return j
							}()
						}(); i < NUM_WEARS; i++ {
							if (c.Equipment[i]) != nil {
								j++
								break
							}
						}
						if j > 0 {
							libc.StrCpy(str, libc.CString("1"))
						} else {
							*str = '\x00'
						}
					} else if (func() int {
						pos = find_eq_pos_script(subfield)
						return pos
					}()) < 0 || (c.Equipment[pos]) == nil {
						*str = '\x00'
					} else {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, (c.Equipment[pos]).Id)
					}
				}
				if libc.StrCaseCmp(field, libc.CString("exp")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = MIN(int64(libc.Atoi(libc.GoString(subfield))), 2100000000)
						gain_exp(c, addition)
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Exp)
				}
			case 'f':
				if libc.StrCaseCmp(field, libc.CString("fighting")) == 0 {
					if c.Fighting != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, c.Fighting.Id)
					} else {
						*str = '\x00'
					}
				} else if libc.StrCaseCmp(field, libc.CString("flying")) == 0 {
					if AFF_FLAGGED(c, AFF_FLYING) {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("follower")) == 0 {
					if c.Followers == nil || c.Followers.Follower == nil {
						*str = '\x00'
					} else {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, c.Followers.Follower.Id)
					}
				}
			case 'g':
				if libc.StrCaseCmp(field, libc.CString("gold")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Gold += addition
					}
					stdio.Snprintf(str, int(slen), "%d", c.Gold)
				}
			case 'h':
				if libc.StrCaseCmp(field, libc.CString("has_item")) == 0 {
					if subfield == nil || *subfield == 0 {
						*str = '\x00'
					} else {
						stdio.Snprintf(str, int(slen), "%d", char_has_item(subfield, c))
					}
				} else if libc.StrCaseCmp(field, libc.CString("hisher")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", HSHR(c))
				} else if libc.StrCaseCmp(field, libc.CString("heshe")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", HSSH(c))
				} else if libc.StrCaseCmp(field, libc.CString("himher")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", HMHR(c))
				} else if libc.StrCaseCmp(field, libc.CString("hitp")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Hit += addition
						update_pos(c)
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Hit)
				} else if libc.StrCaseCmp(field, libc.CString("hunger")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Player_specials.Conditions[HUNGER] = int8(MAX(-1, MIN(int64(addition), 24)))
					}
					stdio.Snprintf(str, int(slen), "%d", c.Player_specials.Conditions[HUNGER])
				}
			case 'i':
				if libc.StrCaseCmp(field, libc.CString("id")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", c.Id)
				} else if libc.StrCaseCmp(field, libc.CString("is_pc")) == 0 {
					if IS_NPC(c) {
						libc.StrCpy(str, libc.CString("0"))
					} else {
						libc.StrCpy(str, libc.CString("1"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("inventory")) == 0 {
					if subfield != nil && *subfield != 0 {
						for obj = c.Carrying; obj != nil; obj = obj.Next_content {
							if GET_OBJ_VNUM(obj) == obj_vnum(libc.Atoi(libc.GoString(subfield))) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, obj.Id)
								return
							}
						}
						if obj == nil {
							*str = '\x00'
						}
					} else {
						if c.Carrying != nil {
							stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, c.Carrying.Id)
						} else {
							*str = '\x00'
						}
					}
				} else if libc.StrCaseCmp(field, libc.CString("is_killer")) == 0 {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(libc.CString("on"), subfield) == 0 {
							SET_BIT_AR(c.Act[:], PLR_KILLER)
						} else if libc.StrCaseCmp(libc.CString("off"), subfield) == 0 {
							REMOVE_BIT_AR(c.Act[:], PLR_KILLER)
						}
					}
					if PLR_FLAGGED(c, PLR_KILLER) {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("is_thief")) == 0 {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(libc.CString("on"), subfield) == 0 {
							SET_BIT_AR(c.Act[:], PLR_THIEF)
						} else if libc.StrCaseCmp(libc.CString("off"), subfield) == 0 {
							REMOVE_BIT_AR(c.Act[:], PLR_THIEF)
						}
					}
					if PLR_FLAGGED(c, PLR_THIEF) {
						libc.StrCpy(str, libc.CString("1"))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("int")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Intel += int8(addition)
						if int(c.Aff_abils.Intel) > max {
							c.Aff_abils.Intel = int8(max)
						}
						if int(c.Aff_abils.Intel) < 3 {
							c.Aff_abils.Intel = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Intel)
				}
			case 'l':
				if libc.StrCaseCmp(field, libc.CString("level")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", GET_LEVEL(c))
				}
			case 'm':
				if libc.StrCaseCmp(field, libc.CString("maxhitp")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Max_hit = MAX(c.Max_hit+addition, 1)
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Max_hit)
				} else if libc.StrCaseCmp(field, libc.CString("mana")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Mana += addition
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Mana)
				} else if libc.StrCaseCmp(field, libc.CString("maxmana")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Max_mana = MAX(c.Max_mana+addition, 1)
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Max_mana)
				} else if libc.StrCaseCmp(field, libc.CString("move")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Move += addition
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Move)
				} else if libc.StrCaseCmp(field, libc.CString("maxmove")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int64 = int64(libc.Atoi(libc.GoString(subfield)))
						c.Max_move = MAX(c.Max_move+addition, 1)
					}
					stdio.Snprintf(str, int(slen), "%lld", c.Max_move)
				} else if libc.StrCaseCmp(field, libc.CString("master")) == 0 {
					if c.Master == nil {
						*str = '\x00'
					} else {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, c.Master.Id)
					}
				}
			case 'n':
				if libc.StrCaseCmp(field, libc.CString("name")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", GET_NAME(c))
				} else if libc.StrCaseCmp(field, libc.CString("next_in_room")) == 0 {
					if c.Next_in_room != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, c.Next_in_room.Id)
					} else {
						*str = '\x00'
					}
				}
			case 'p':
				if libc.StrCaseCmp(field, libc.CString("pos")) == 0 {
					if subfield != nil && *subfield != 0 {
						for i = POS_SLEEPING; i <= POS_STANDING; i++ {
							if libc.StrNCaseCmp(subfield, position_types[i], libc.StrLen(subfield)) == 0 {
								c.Position = int8(i)
								break
							}
						}
					}
					stdio.Snprintf(str, int(slen), "%s", position_types[c.Position])
				} else if libc.StrCaseCmp(field, libc.CString("prac")) == 0 {
					if IS_NPC(c) {
						if c.In_room != room_rnum(-1) {
							send_to_room(c.In_room, libc.CString("Error!: Report this trigger error to the coding authorities!\r\n"))
						}
					}
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Player_specials.Class_skill_points[c.Chclass] = int(MAX(0, int64((c.Player_specials.Class_skill_points[c.Chclass])+addition)))
					}
					stdio.Snprintf(str, int(slen), "%d", c.Player_specials.Class_skill_points[c.Chclass])
				} else if libc.StrCaseCmp(field, libc.CString("plr")) == 0 {
					if subfield != nil && *subfield != 0 {
						var plr int = get_flag_by_name(player_bits[:], subfield)
						if plr != int(-1) && PLR_FLAGGED(c, bitvector_t(int32(plr))) {
							libc.StrCpy(str, libc.CString("1"))
						} else {
							libc.StrCpy(str, libc.CString("0"))
						}
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				} else if libc.StrCaseCmp(field, libc.CString("pref")) == 0 {
					if subfield != nil && *subfield != 0 {
						var pref int = get_flag_by_name(preference_bits[:], subfield)
						if pref != int(-1) && PRF_FLAGGED(c, bitvector_t(int32(pref))) {
							libc.StrCpy(str, libc.CString("1"))
						} else {
							libc.StrCpy(str, libc.CString("0"))
						}
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				}
			case 'r':
				if libc.StrCaseCmp(field, libc.CString("room")) == 0 {
					stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, func() room_vnum {
						if c.In_room != room_rnum(-1) {
							return world[c.In_room].Number + ROOM_ID_BASE
						}
						return ROOM_ID_BASE
					}())
				} else if libc.StrCaseCmp(field, libc.CString("race")) == 0 {
					if IS_NPC(c) {
						sprinttype(int(c.Race), race_names[:], str, slen)
					} else {
						sprinttype(int(c.Race), race_names[:], str, slen)
					}
				} else if libc.StrCaseCmp(field, libc.CString("rpp")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Rp += addition
					}
					stdio.Snprintf(str, int(slen), "%d", c.Rp)
				}
			case 's':
				if libc.StrCaseCmp(field, libc.CString("sex")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", genders[int(c.Sex)])
				} else if libc.StrCaseCmp(field, libc.CString("str")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Str += int8(addition)
						if int(c.Aff_abils.Str) > max {
							c.Aff_abils.Str = int8(max)
						}
						if int(c.Aff_abils.Str) < 3 {
							c.Aff_abils.Str = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Str)
				} else if libc.StrCaseCmp(field, libc.CString("size")) == 0 {
					if subfield != nil && *subfield != 0 {
						var ns int
						if (func() int {
							ns = search_block(subfield, &size_names[0], FALSE)
							return ns
						}()) > -1 {
							c.Size = ns
						}
					}
					sprinttype(get_size(c), size_names[:], str, slen)
				} else if libc.StrCaseCmp(field, libc.CString("skill")) == 0 {
					stdio.Snprintf(str, int(slen), "%s", skill_percent(c, subfield))
				} else if libc.StrCaseCmp(field, libc.CString("skillset")) == 0 {
					if !IS_NPC(c) && subfield != nil && *subfield != 0 {
						var (
							skillname [2048]byte
							amount    *byte
						)
						amount = one_word(subfield, &skillname[0])
						skip_spaces(&amount)
						if amount != nil && *amount != 0 && is_number(amount) != 0 {
							var skillnum int = find_skill_num(&skillname[0], 1<<1)
							if skillnum > 0 {
								var new_value int = int(MAX(0, MIN(100, int64(libc.Atoi(libc.GoString(amount))))))
								for {
									c.Skills[skillnum] = int8(new_value)
									if true {
										break
									}
								}
							}
						}
					}
					*str = '\x00'
				} else if libc.StrCaseCmp(field, libc.CString("saving_fortitude")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Apply_saving_throw[SAVING_FORTITUDE] += int16(addition)
					}
					stdio.Snprintf(str, int(slen), "%d", c.Apply_saving_throw[SAVING_FORTITUDE])
				} else if libc.StrCaseCmp(field, libc.CString("saving_reflex")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Apply_saving_throw[SAVING_REFLEX] += int16(addition)
					}
					stdio.Snprintf(str, int(slen), "%d", c.Apply_saving_throw[SAVING_REFLEX])
				} else if libc.StrCaseCmp(field, libc.CString("saving_will")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Apply_saving_throw[SAVING_WILL] += int16(addition)
					}
					stdio.Snprintf(str, int(slen), "%d", c.Apply_saving_throw[SAVING_WILL])
				}
			case 't':
				if libc.StrCaseCmp(field, libc.CString("thirst")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Player_specials.Conditions[THIRST] = int8(MAX(-1, MIN(int64(addition), 24)))
					}
					stdio.Snprintf(str, int(slen), "%d", c.Player_specials.Conditions[THIRST])
				} else if libc.StrCaseCmp(field, libc.CString("tnl")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", level_exp(c, GET_LEVEL(c)+1))
				}
			case 'v':
				if libc.StrCaseCmp(field, libc.CString("vnum")) == 0 {
					if subfield != nil && *subfield != 0 {
						stdio.Snprintf(str, int(slen), "%d", func() int {
							if IS_NPC(c) {
								return int(libc.BoolToInt(GET_MOB_VNUM(c) == mob_vnum(libc.Atoi(libc.GoString(subfield)))))
							}
							return -1
						}())
					} else {
						if IS_NPC(c) {
							stdio.Snprintf(str, int(slen), "%d", GET_MOB_VNUM(c))
						} else {
							libc.StrCpy(str, libc.CString("-1"))
						}
					}
				} else if libc.StrCaseCmp(field, libc.CString("varexists")) == 0 {
					var remote_vd *trig_var_data
					libc.StrCpy(str, libc.CString("0"))
					if c.Script != nil {
						for remote_vd = c.Script.Global_vars; remote_vd != nil; remote_vd = remote_vd.Next {
							if libc.StrCaseCmp(remote_vd.Name, subfield) == 0 {
								break
							}
						}
						if remote_vd != nil {
							libc.StrCpy(str, libc.CString("1"))
						}
					}
				}
			case 'w':
				if libc.StrCaseCmp(field, libc.CString("weight")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", c.Weight)
				} else if libc.StrCaseCmp(field, libc.CString("wis")) == 0 {
					if subfield != nil && *subfield != 0 {
						var (
							addition int = libc.Atoi(libc.GoString(subfield))
							max      int = 100
						)
						c.Aff_abils.Wis += int8(addition)
						if int(c.Aff_abils.Wis) > max {
							c.Aff_abils.Wis = int8(max)
						}
						if int(c.Aff_abils.Wis) < 3 {
							c.Aff_abils.Wis = 3
						}
					}
					stdio.Snprintf(str, int(slen), "%d", c.Aff_abils.Wis)
				}
			case 'z':
				if libc.StrCaseCmp(field, libc.CString("zenni")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						c.Gold += addition
					}
					stdio.Snprintf(str, int(slen), "%d", c.Gold)
				}
			}
			if *str == '\x01' {
				if c.Script != nil {
					for vd = c.Script.Global_vars; vd != nil; vd = vd.Next {
						if libc.StrCaseCmp(vd.Name, field) == 0 {
							break
						}
					}
					if vd != nil {
						stdio.Snprintf(str, int(slen), "%s", vd.Value)
					} else {
						*str = '\x00'
						script_log(libc.CString("Trigger: %s, VNum %d. unknown char field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, field)
					}
				} else {
					*str = '\x00'
					script_log(libc.CString("Trigger: %s, VNum %d. unknown char field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, field)
				}
			}
		} else if o != nil {
			if text_processed(field, subfield, vd, str, slen) != 0 {
				return
			}
			*str = '\x01'
			switch unicode.ToLower(rune(*field)) {
			case 'a':
				if libc.StrCaseCmp(field, libc.CString("affects")) == 0 {
					if subfield != nil && *subfield != 0 {
						if check_flags_by_name_ar(o.Bitvector[:], NUM_AFF_FLAGS, subfield, affected_bits[:]) > 0 {
							stdio.Snprintf(str, int(slen), "1")
						} else {
							stdio.Snprintf(str, int(slen), "0")
						}
					} else {
						stdio.Snprintf(str, int(slen), "0")
					}
				}
			case 'c':
				if libc.StrCaseCmp(field, libc.CString("cost")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						o.Cost = int(MAX(0, int64(addition+o.Cost)))
					}
					stdio.Snprintf(str, int(slen), "%d", o.Cost)
				} else if libc.StrCaseCmp(field, libc.CString("cost_per_day")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						o.Cost_per_day = int(MAX(0, int64(addition+o.Cost_per_day)))
					}
					stdio.Snprintf(str, int(slen), "%d", o.Cost_per_day)
				} else if libc.StrCaseCmp(field, libc.CString("carried_by")) == 0 {
					if o.Carried_by != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, o.Carried_by.Id)
					} else {
						*str = '\x00'
					}
				} else if libc.StrCaseCmp(field, libc.CString("contents")) == 0 {
					if o.Contains != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, o.Contains.Id)
					} else {
						*str = '\x00'
					}
				} else if libc.StrCaseCmp(field, libc.CString("count")) == 0 {
					if int(o.Type_flag) == ITEM_CONTAINER {
						stdio.Snprintf(str, int(slen), "%d", item_in_list(subfield, o.Contains))
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				}
			case 'e':
				if libc.StrCaseCmp(field, libc.CString("extra")) == 0 {
					if subfield != nil && *subfield != 0 {
						if check_flags_by_name_ar(o.Extra_flags[:], NUM_ITEM_FLAGS, subfield, extra_bits[:]) > 0 {
							stdio.Snprintf(str, int(slen), "1")
						} else {
							stdio.Snprintf(str, int(slen), "0")
						}
					} else {
						stdio.Snprintf(str, int(slen), "0")
					}
				} else {
					sprintbitarray(o.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, str)
				}
			case 'h':
				if libc.StrCaseCmp(field, libc.CString("has_in")) == 0 {
					if int(o.Type_flag) == ITEM_CONTAINER {
						stdio.Snprintf(str, int(slen), "%s", func() string {
							if item_in_list(subfield, o.Contains) != 0 {
								return "1"
							}
							return "0"
						}())
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				}
				if libc.StrCaseCmp(field, libc.CString("health")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						o.Value[VAL_ALL_HEALTH] = int(MAX(1, int64(addition+(o.Value[VAL_ALL_HEALTH]))))
						if OBJ_FLAGGED(o, ITEM_BROKEN) && (o.Value[VAL_ALL_HEALTH]) >= 100 {
							REMOVE_BIT_AR(o.Extra_flags[:], ITEM_BROKEN)
						}
					}
					stdio.Snprintf(str, int(slen), "%d", o.Value[VAL_ALL_HEALTH])
				}
			case 'i':
				if libc.StrCaseCmp(field, libc.CString("id")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Id)
				} else if libc.StrCaseCmp(field, libc.CString("is_inroom")) == 0 {
					if o.In_room != room_rnum(-1) {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[o.In_room].Number+ROOM_ID_BASE)
					} else {
						*str = '\x00'
					}
				} else if libc.StrCaseCmp(field, libc.CString("is_pc")) == 0 {
					libc.StrCpy(str, libc.CString("-1"))
				} else if libc.StrCaseCmp(field, libc.CString("itemflag")) == 0 {
					if subfield != nil && *subfield != 0 {
						var item int = get_flag_by_name(extra_bits[:], subfield)
						if item != int(-1) && OBJ_FLAGGED(o, bitvector_t(int32(item))) {
							libc.StrCpy(str, libc.CString("1"))
						} else {
							libc.StrCpy(str, libc.CString("0"))
						}
					} else {
						libc.StrCpy(str, libc.CString("0"))
					}
				}
			case 'l':
				if libc.StrCaseCmp(field, libc.CString("level")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Level)
				}
			case 'n':
				if libc.StrCaseCmp(field, libc.CString("name")) == 0 {
					if subfield == nil || *subfield == 0 {
						stdio.Snprintf(str, int(slen), "%s", o.Name)
					} else {
						var blah [500]byte
						stdio.Sprintf(&blah[0], "%s %s", o.Name, subfield)
						o.Name = libc.StrDup(&blah[0])
					}
				} else if libc.StrCaseCmp(field, libc.CString("next_in_list")) == 0 {
					if o.Next_content != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, o.Next_content.Id)
					} else {
						*str = '\x00'
					}
				}
			case 'r':
				if libc.StrCaseCmp(field, libc.CString("room")) == 0 {
					if obj_room(o) != room_rnum(-1) {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[obj_room(o)].Number+ROOM_ID_BASE)
					} else {
						*str = '\x00'
					}
				}
			case 's':
				if libc.StrCaseCmp(field, libc.CString("shortdesc")) == 0 {
					if subfield == nil || *subfield == 0 {
						stdio.Snprintf(str, int(slen), "%s", o.Short_description)
					} else {
						var blah [500]byte
						stdio.Sprintf(&blah[0], "%s @wnicknamed @D(@C%s@D)@n", o.Short_description, subfield)
						o.Short_description = libc.StrDup(&blah[0])
					}
				} else if libc.StrCaseCmp(field, libc.CString("setaffects")) == 0 {
					if subfield != nil && *subfield != 0 {
						var ns int
						if (func() int {
							ns = check_flags_by_name_ar(o.Bitvector[:], NUM_AFF_FLAGS, subfield, affected_bits[:])
							return ns
						}()) > 0 {
							TOGGLE_BIT_AR(o.Bitvector[:], bitvector_t(int32(ns)))
							stdio.Snprintf(str, int(slen), "1")
						}
					}
				} else if libc.StrCaseCmp(field, libc.CString("setextra")) == 0 {
					if subfield != nil && *subfield != 0 {
						var ns int
						if (func() int {
							ns = check_flags_by_name_ar(o.Extra_flags[:], NUM_ITEM_FLAGS, subfield, extra_bits[:])
							return ns
						}()) > 0 {
							TOGGLE_BIT_AR(o.Extra_flags[:], bitvector_t(int32(ns)))
							stdio.Snprintf(str, int(slen), "1")
						}
					}
				} else if libc.StrCaseCmp(field, libc.CString("size")) == 0 {
					if subfield != nil && *subfield != 0 {
						var ns int
						if (func() int {
							ns = search_block(subfield, &size_names[0], FALSE)
							return ns
						}()) > -1 {
							o.Size = ns
						}
					}
					sprinttype(o.Size, size_names[:], str, slen)
				}
			case 't':
				if libc.StrCaseCmp(field, libc.CString("type")) == 0 {
					sprinttype(int(o.Type_flag), item_types[:], str, slen)
				} else if libc.StrCaseCmp(field, libc.CString("timer")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Timer)
				}
			case 'v':
				if libc.StrCaseCmp(field, libc.CString("vnum")) == 0 {
					if subfield != nil && *subfield != 0 {
						stdio.Snprintf(str, int(slen), "%d", int(libc.BoolToInt(GET_OBJ_VNUM(o) == obj_vnum(libc.Atoi(libc.GoString(subfield))))))
					} else {
						stdio.Snprintf(str, int(slen), "%d", GET_OBJ_VNUM(o))
					}
				} else if libc.StrCaseCmp(field, libc.CString("val0")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[0])
				} else if libc.StrCaseCmp(field, libc.CString("val1")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[1])
				} else if libc.StrCaseCmp(field, libc.CString("val2")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[2])
				} else if libc.StrCaseCmp(field, libc.CString("val3")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[3])
				} else if libc.StrCaseCmp(field, libc.CString("val4")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[4])
				} else if libc.StrCaseCmp(field, libc.CString("val5")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[5])
				} else if libc.StrCaseCmp(field, libc.CString("val6")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[6])
				} else if libc.StrCaseCmp(field, libc.CString("val7")) == 0 {
					stdio.Snprintf(str, int(slen), "%d", o.Value[7])
				}
			case 'w':
				if libc.StrCaseCmp(field, libc.CString("weight")) == 0 {
					if subfield != nil && *subfield != 0 {
						var addition int = libc.Atoi(libc.GoString(subfield))
						if addition < 0 || addition > 0 {
							o.Weight = MAX(0, int64(addition+int(o.Weight)))
						} else {
							o.Weight = 0
						}
					}
					stdio.Snprintf(str, int(slen), "%lld", o.Weight)
				} else if libc.StrCaseCmp(field, libc.CString("worn_by")) == 0 {
					if o.Worn_by != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, o.Worn_by.Id)
					} else {
						*str = '\x00'
					}
				}
			}
			if *str == '\x01' {
				if o.Script != nil {
					for vd = o.Script.Global_vars; vd != nil; vd = vd.Next {
						if libc.StrCaseCmp(vd.Name, field) == 0 {
							break
						}
					}
					if vd != nil {
						stdio.Snprintf(str, int(slen), "%s", vd.Value)
					} else {
						*str = '\x00'
						if libc.StrCaseCmp(trig.Name, libc.CString("Rename Object")) != 0 {
							script_log(libc.CString("Trigger: %s, VNum %d, type: %d. unknown object field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, type_, field)
						}
					}
				} else {
					*str = '\x00'
					if libc.StrCaseCmp(trig.Name, libc.CString("Rename Object")) != 0 {
						script_log(libc.CString("Trigger: %s, VNum %d, type: %d. unknown object field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, type_, field)
					}
				}
			}
		} else if r != nil {
			if text_processed(field, subfield, vd, str, slen) != 0 {
				return
			}
			if r.Number == 0 {
				if r.Script == nil {
					*str = '\x00'
					script_log(libc.CString("Trigger: %s, Vnum %d, type %d. Trying to access Global var list of void. Apparently this has not been set up!"), trig.Name, trig_index[trig.Nr].Vnum, type_)
				} else {
					for vd = r.Script.Global_vars; vd != nil; vd = vd.Next {
						if libc.StrCaseCmp(vd.Name, field) == 0 {
							break
						}
					}
					if vd != nil {
						stdio.Snprintf(str, int(slen), "%s", vd.Value)
					} else {
						*str = '\x00'
					}
				}
			} else if libc.StrCaseCmp(field, libc.CString("name")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", r.Name)
			} else if libc.StrCaseCmp(field, libc.CString("sector")) == 0 {
				sprinttype(r.Sector_type, sector_types[:], str, slen)
			} else if libc.StrCaseCmp(field, libc.CString("gravity")) == 0 {
				stdio.Snprintf(str, int(slen), "%d", r.Gravity)
			} else if libc.StrCaseCmp(field, libc.CString("vnum")) == 0 {
				if subfield != nil && *subfield != 0 {
					stdio.Snprintf(str, int(slen), "%d", int(libc.BoolToInt(r.Number == room_vnum(libc.Atoi(libc.GoString(subfield))))))
				} else {
					stdio.Snprintf(str, int(slen), "%d", r.Number)
				}
			} else if libc.StrCaseCmp(field, libc.CString("contents")) == 0 {
				if subfield != nil && *subfield != 0 {
					for obj = r.Contents; obj != nil; obj = obj.Next_content {
						if GET_OBJ_VNUM(obj) == obj_vnum(libc.Atoi(libc.GoString(subfield))) {
							stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, obj.Id)
							return
						}
					}
					if obj == nil {
						*str = '\x00'
					}
				} else {
					if r.Contents != nil {
						stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, r.Contents.Id)
					} else {
						*str = '\x00'
					}
				}
			} else if libc.StrCaseCmp(field, libc.CString("people")) == 0 {
				if r.People != nil {
					stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, r.People.Id)
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("id")) == 0 {
				var rnum room_rnum = real_room(r.Number)
				if rnum != room_rnum(-1) {
					stdio.Snprintf(str, int(slen), "%d", world[rnum].Number+ROOM_ID_BASE)
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("weather")) == 0 {
				var sky_look [4]*byte = [4]*byte{libc.CString("sunny"), libc.CString("cloudy"), libc.CString("rainy"), libc.CString("lightning")}
				if !IS_SET_AR(r.Room_flags[:], ROOM_INDOORS) {
					stdio.Snprintf(str, int(slen), "%s", sky_look[weather_info.Sky])
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("fishing")) == 0 {
				var thisroom room_rnum = real_room(r.Number)
				if ROOM_FLAGGED(thisroom, ROOM_FISHING) {
					stdio.Snprintf(str, int(slen), "1")
				} else {
					stdio.Snprintf(str, int(slen), "0")
				}
			} else if libc.StrCaseCmp(field, libc.CString("zonenumber")) == 0 {
				stdio.Snprintf(str, int(slen), "%d", zone_table[r.Zone].Number)
			} else if libc.StrCaseCmp(field, libc.CString("zonename")) == 0 {
				stdio.Snprintf(str, int(slen), "%s", zone_table[r.Zone].Name)
			} else if libc.StrCaseCmp(field, libc.CString("roomflag")) == 0 {
				if subfield != nil && *subfield != 0 {
					var thisroom room_rnum = real_room(r.Number)
					if check_flags_by_name_ar(world[thisroom].Room_flags[:], NUM_ROOM_FLAGS, subfield, room_bits[:]) > 0 {
						stdio.Snprintf(str, int(slen), "1")
					} else {
						stdio.Snprintf(str, int(slen), "0")
					}
				} else {
					stdio.Snprintf(str, int(slen), "0")
				}
			} else if libc.StrCaseCmp(field, libc.CString("north")) == 0 {
				if (r.Dir_option[NORTH]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[NORTH]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[NORTH]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[NORTH]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[NORTH]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[NORTH]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[NORTH]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("east")) == 0 {
				if (r.Dir_option[EAST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[EAST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[EAST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[EAST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[EAST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[EAST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[EAST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("south")) == 0 {
				if (r.Dir_option[SOUTH]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[SOUTH]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[SOUTH]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[SOUTH]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[SOUTH]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[SOUTH]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[SOUTH]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("west")) == 0 {
				if (r.Dir_option[WEST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[WEST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[WEST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[WEST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[WEST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[WEST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[WEST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("up")) == 0 {
				if (r.Dir_option[UP]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[UP]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[UP]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[UP]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[UP]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[UP]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[UP]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("down")) == 0 {
				if (r.Dir_option[DOWN]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[DOWN]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[DOWN]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[DOWN]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[DOWN]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[DOWN]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[DOWN]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("northwest")) == 0 {
				if (r.Dir_option[NORTHWEST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[NORTHWEST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[NORTHWEST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[NORTHWEST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[NORTHWEST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[NORTHWEST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[NORTHWEST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("northeast")) == 0 {
				if (r.Dir_option[NORTHEAST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[NORTHEAST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[NORTHEAST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[NORTHEAST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[NORTHEAST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[NORTHEAST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[NORTHEAST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("southwest")) == 0 {
				if (r.Dir_option[SOUTHWEST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[SOUTHWEST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[SOUTHWEST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[SOUTHWEST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[SOUTHWEST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[SOUTHWEST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[SOUTHWEST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("southeast")) == 0 {
				if (r.Dir_option[SOUTHEAST]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[SOUTHEAST]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[SOUTHEAST]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[SOUTHEAST]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[SOUTHEAST]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[SOUTHEAST]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[SOUTHEAST]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("inside")) == 0 {
				if (r.Dir_option[INDIR]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[INDIR]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[INDIR]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[INDIR]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[INDIR]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[INDIR]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[INDIR]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else if libc.StrCaseCmp(field, libc.CString("outside")) == 0 {
				if (r.Dir_option[OUTDIR]) != nil {
					if subfield != nil && *subfield != 0 {
						if libc.StrCaseCmp(subfield, libc.CString("vnum")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", GET_ROOM_VNUM((r.Dir_option[OUTDIR]).To_room))
						} else if libc.StrCaseCmp(subfield, libc.CString("key")) == 0 {
							stdio.Snprintf(str, int(slen), "%d", (r.Dir_option[OUTDIR]).Key)
						} else if libc.StrCaseCmp(subfield, libc.CString("bits")) == 0 {
							sprintbit((r.Dir_option[OUTDIR]).Exit_info, exit_bits[:], str, slen)
						} else if libc.StrCaseCmp(subfield, libc.CString("room")) == 0 {
							if (r.Dir_option[OUTDIR]).To_room != room_rnum(-1) {
								stdio.Snprintf(str, int(slen), "%c%d", UID_CHAR, world[(r.Dir_option[OUTDIR]).To_room].Number+ROOM_ID_BASE)
							} else {
								*str = '\x00'
							}
						}
					} else {
						sprintbit((r.Dir_option[OUTDIR]).Exit_info, exit_bits[:], str, slen)
					}
				} else {
					*str = '\x00'
				}
			} else {
				if r.Script != nil {
					for vd = r.Script.Global_vars; vd != nil; vd = vd.Next {
						if libc.StrCaseCmp(vd.Name, field) == 0 {
							break
						}
					}
					if vd != nil {
						stdio.Snprintf(str, int(slen), "%s", vd.Value)
					} else {
						*str = '\x00'
						script_log(libc.CString("Trigger: %s, VNum %d, type: %d. unknown room field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, type_, field)
					}
				} else {
					*str = '\x00'
					script_log(libc.CString("Trigger: %s, VNum %d, type: %d. unknown room field: '%s'"), trig.Name, trig_index[trig.Nr].Vnum, type_, field)
				}
			}
		}
	}
}
func var_subst(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, line *byte, buf *byte) {
	var (
		tmp         [2048]byte
		repl_str    [2048]byte
		var_        *byte = nil
		field       *byte = nil
		p           *byte = nil
		tmp2        [2048]byte
		subfield_p  *byte
		subfield    [2048]byte
		left        int
		len_        int
		paren_count int = 0
		dots        int = 0
	)
	if libc.StrChr(line, '%') == nil {
		libc.StrCpy(buf, line)
		return
	}
	repl_str[0] = func() byte {
		p := &tmp[0]
		tmp[0] = func() byte {
			p := &tmp2[0]
			tmp2[0] = '\x00'
			return *p
		}()
		return *p
	}()
	p = libc.StrCpy(&tmp[0], line)
	subfield_p = &subfield[0]
	left = int(MAX_INPUT_LENGTH - 1)
	for *p != 0 && left > 0 {
		for *p != 0 && *p != '%' && left > 0 {
			*(func() *byte {
				p := &buf
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *(func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			left--
		}
		*buf = '\x00'
		if *p != 0 && *(func() *byte {
			p := &p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return *p
		}()) == '%' && left > 0 {
			*(func() *byte {
				p := &buf
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *(func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			*buf = '\x00'
			left--
			continue
		} else if *p != 0 && left > 0 {
			for var_ = p; *p != 0 && *p != '%' && *p != '.'; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
			}
			field = p
			if *p == '.' {
				*(func() *byte {
					p := &p
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()) = '\x00'
				dots = 0
				for field = p; *p != 0 && (*p != '%' || paren_count > 0 || dots != 0); p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
					if dots > 0 {
						*subfield_p = '\x00'
						find_replacement(gohere, sc, trig, type_, var_, field, &subfield[0], &repl_str[0], uint64(2048))
						if repl_str[0] != 0 {
							stdio.Snprintf(&tmp2[0], int(2048), "eval tmpvr %s", &repl_str[0])
							process_eval(gohere, sc, trig, type_, &tmp2[0])
							libc.StrCpy(var_, libc.CString("tmpvr"))
							field = p
							dots = 0
							continue
						}
						dots = 0
					} else if *p == '(' {
						*p = '\x00'
						paren_count++
					} else if *p == ')' {
						*p = '\x00'
						paren_count--
					} else if paren_count > 0 {
						*func() *byte {
							p := &subfield_p
							x := *p
							*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
							return x
						}() = *p
					} else if *p == '.' {
						*p = '\x00'
						dots++
					}
				}
			}
			*(func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = '\x00'
			*subfield_p = '\x00'
			if subfield[0] != 0 {
				var_subst(gohere, sc, trig, type_, &subfield[0], &tmp2[0])
				libc.StrCpy(&subfield[0], &tmp2[0])
			}
			find_replacement(gohere, sc, trig, type_, var_, field, &subfield[0], &repl_str[0], uint64(2048))
			libc.StrNCat(buf, &repl_str[0], left)
			len_ = libc.StrLen(&repl_str[0])
			buf = (*byte)(unsafe.Add(unsafe.Pointer(buf), len_))
			left -= len_
		}
	}
}
