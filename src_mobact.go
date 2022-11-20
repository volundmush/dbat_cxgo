package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const MOB_AGGR_TO_ALIGN = 11

func mob_absorb(ch *char_data, vict *char_data) {
	if ch.Absorbing != nil {
		act(libc.CString("@R$n@w releases YOU from $s grip!@n"), 1, ch, nil, unsafe.Pointer(ch.Absorbing), TO_VICT)
		act(libc.CString("@R$n@w releases @R$N@w from $s grip!@n"), 1, ch, nil, unsafe.Pointer(ch.Absorbing), TO_NOTVICT)
		ch.Absorbing.Absorbby = nil
		ch.Absorbing = nil
		return
	}
	var zanzo int = 0
	var roll int = 0
	var chance int = int(float64(GET_LEVEL(ch)) * 0.5)
	var chance2 int = GET_LEVEL(ch) + 10
	if chance2 > 118 {
		chance2 = 118
	}
	if GET_LEVEL(ch) < 2 {
		return
	} else {
		roll = rand_number(chance, chance2)
	}
	if vict == nil {
		return
	}
	if int(vict.Race) == RACE_ANDROID {
		return
	}
	if AFF_FLAGGED(vict, AFF_ZANZOKEN) {
		if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
			if GET_SPEEDI(ch) < GET_SPEEDI(vict) {
				zanzo = 1
			} else {
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
			}
		} else {
			zanzo = 1
		}
		if zanzo == 1 {
			act(libc.CString("@R$n@c tries to grab @RYOU@c but you @Czanzoken@c out of the way!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$n@ctries to grab @R$N@c but $E @Czanzokens@c out of the way!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			REMOVE_BIT_AR(ch.Affected_by[:], AFF_ZANZOKEN)
			REMOVE_BIT_AR(vict.Affected_by[:], AFF_ZANZOKEN)
			return
		} else {
			act(libc.CString("@cYou try to @Czanzoken@c out of @R$n's@c reach, but $e is too fast!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N tries to @Czanzoken@c out of @R$n's@c reach, but $e is too fast!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			REMOVE_BIT_AR(vict.Affected_by[:], AFF_ZANZOKEN)
		}
	}
	if roll < check_def(vict) {
		act(libc.CString("@R$n@r tries to grab YOU, but you manage to evade $s grasp!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@R$n@r tries to grab @R$N@r, but @R$N@r manages to evade!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		return
	} else {
		act(libc.CString("@R$n@r grabs onto YOU and starts to absorb your energy!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@R$n@r grabs onto @R$N@r and starts to absorb your energy!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Absorbing = vict
		vict.Absorbby = ch
		return
	}
}
func player_present(ch *char_data) int {
	var (
		vict   *char_data
		next_v *char_data
		found  int = 0
	)
	if ch.In_room == int(-1) {
		return 0
	}
	for vict = world[ch.In_room].People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if !IS_NPC(vict) {
			found = 1
		}
	}
	return found
}
func mobile_activity() {
	var (
		ch       *char_data
		next_ch  *char_data
		vict     *char_data
		obj      *obj_data
		best_obj *obj_data
		door     int
		found    int
		max      int
		names    *memory_rec_struct
	)
	for ch = character_list; ch != nil; ch = next_ch {
		next_ch = ch.Next
		if !IS_MOB(ch) {
			continue
		}
		if MOB_FLAGGED(ch, MOB_SPEC) && no_specials == 0 {
			if mob_index[ch.Nr].Func == nil {
				basic_mud_log(libc.CString("SYSERR: %s (#%d): Attempting to call non-existing mob function."), GET_NAME(ch), GET_MOB_VNUM(ch))
				REMOVE_BIT_AR(ch.Act[:], MOB_SPEC)
			} else {
				var actbuf [2048]byte = func() [2048]byte {
					var t [2048]byte
					copy(t[:], []byte(""))
					return t
				}()
				if mob_index[ch.Nr].Func(ch, unsafe.Pointer(ch), 0, &actbuf[0]) {
					continue
				}
			}
		}
		if !AWAKE(ch) {
			continue
		}
		if IS_HUMANOID(ch) && ch.Fighting == nil && AWAKE(ch) && !MOB_FLAGGED(ch, MOB_NOSCAVENGER) && !MOB_FLAGGED(ch, MOB_NOKILL) && (player_present(ch) == 0 || axion_dice(0) > 118) {
			if world[ch.In_room].Contents != nil && rand_number(1, 100) >= 95 {
				max = 1
				best_obj = nil
				for obj = world[ch.In_room].Contents; obj != nil; obj = obj.Next_content {
					if CAN_GET_OBJ(ch, obj) && obj.Cost > max {
						best_obj = obj
						max = obj.Cost
					}
				}
				if best_obj != nil && CAN_GET_OBJ(ch, best_obj) && int(best_obj.Type_flag) != ITEM_BED && best_obj.Posted_to == nil && !OBJ_FLAGGED(best_obj, ITEM_NOPICKUP) {
					switch rand_number(1, 5) {
					case 1:
						act(libc.CString("$n@W says, '@CFinders keepers, losers weepers.@W'@n"), 1, ch, nil, nil, TO_ROOM)
					case 2:
						act(libc.CString("$n@W says, '@CPeople always leaving their garbage JUST LYING AROUND. The nerve....@W'@n"), 1, ch, nil, nil, TO_ROOM)
					case 3:
						act(libc.CString("$n@W says, '@CWho would leave this here? Oh well..@W'@n"), 1, ch, nil, nil, TO_ROOM)
					case 4:
						act(libc.CString("$n@W says, '@CI always wanted one of these.@W'@n"), 1, ch, nil, nil, TO_ROOM)
					case 5:
						act(libc.CString("$n@W looks around quickly to see if anyone is paying attention.@n"), 1, ch, nil, nil, TO_ROOM)
					}
					perform_get_from_room(ch, best_obj)
				}
			}
		}
		if !MOB_FLAGGED(ch, MOB_SENTINEL) && int(ch.Position) == POS_STANDING && ch.Fighting == nil && !AFF_FLAGGED(ch, AFF_TAMED) && ch.Absorbby == nil && (func() int {
			door = rand_number(0, 18)
			return door
		}()) < NUM_OF_DIRS && CAN_GO(ch, door) && !ROOM_FLAGGED((world[ch.In_room].Dir_option[door]).To_room, ROOM_NOMOB) && !ROOM_FLAGGED((world[ch.In_room].Dir_option[door]).To_room, ROOM_DEATH) && (!MOB_FLAGGED(ch, MOB_STAY_ZONE) || world[(world[ch.In_room].Dir_option[door]).To_room].Zone == world[ch.In_room].Zone) {
			if rand_number(1, 2) == 2 && !AFF_FLAGGED(ch, 18) && block_calc(ch) {
				perform_move(ch, door, 1)
			}
		}
		var hugeatk *obj_data = nil
		var next_huge *obj_data = nil
		for hugeatk = world[ch.In_room].Contents; hugeatk != nil; hugeatk = next_huge {
			next_huge = hugeatk.Next_content
			if ch.Fighting != nil {
				continue
			}
			if MOB_FLAGGED(ch, MOB_NOKILL) {
				continue
			}
			if GET_OBJ_VNUM(hugeatk) == 82 || GET_OBJ_VNUM(hugeatk) == 83 {
				if hugeatk.User != nil {
					act(libc.CString("@W$n@R leaps at @C$N@R desperately!@n"), 1, ch, nil, unsafe.Pointer(hugeatk.User), TO_ROOM)
					act(libc.CString("@W$n@R leaps at YOU desperately!@n"), 1, ch, nil, unsafe.Pointer(hugeatk.User), TO_VICT)
					if IS_HUMANOID(ch) {
						var tar [2048]byte
						stdio.Sprintf(&tar[0], "%s", GET_NAME(hugeatk.User))
						do_punch(ch, &tar[0], 0, 0)
					} else {
						var tar [2048]byte
						stdio.Sprintf(&tar[0], "%s", GET_NAME(hugeatk.User))
						do_bite(ch, &tar[0], 0, 0)
					}
				}
			}
		}
		if MOB_FLAGGED(ch, MOB_AGGRESSIVE) && !AFF_FLAGGED(ch, 18) {
			var spot_roll int = rand_number(1, GET_LEVEL(ch)+10)
			found = 0
			for vict = world[ch.In_room].People; vict != nil && found == 0; vict = vict.Next_in_room {
				if vict == ch {
					continue
				} else if ch.Fighting != nil {
					continue
				} else if !CAN_SEE(ch, vict) {
					continue
				} else if IS_NPC(vict) {
					continue
				} else if PRF_FLAGGED(vict, PRF_NOHASSLE) {
					continue
				} else if MOB_FLAGGED(ch, MOB_AGGR_EVIL) && vict.Alignment < 50 {
					continue
				} else if MOB_FLAGGED(ch, MOB_AGGR_GOOD) && vict.Alignment > -50 {
					continue
				} else if GET_LEVEL(vict) < 5 {
					continue
				} else if AFF_FLAGGED(vict, AFF_HIDE) && GET_SKILL(vict, SKILL_HIDE) > spot_roll {
					continue
				} else if AFF_FLAGGED(vict, AFF_SNEAK) && GET_SKILL(vict, SKILL_MOVE_SILENTLY) > spot_roll {
					continue
				} else if ch.Aggtimer < 8 {
					ch.Aggtimer += 1
				} else {
					ch.Aggtimer = 0
					var tar [2048]byte
					stdio.Sprintf(&tar[0], "%s", GET_NAME(vict))
					if IS_HUMANOID(ch) {
						if !AFF_FLAGGED(vict, AFF_HIDE) && !AFF_FLAGGED(vict, AFF_SNEAK) {
							act(libc.CString("@w'I am going to get you!' @C$n@w shouts at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@w'I am going to get you!' @C$n@w shouts at @c$N@w!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						} else {
							act(libc.CString("@C$n@w notices YOU.\n@w'I am going to get you!' @C$n@w shouts at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@C$n@w notices @c$N@w.\n@w'I am going to get you!' @C$n@w shouts at @c$N@w!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						}
						if AFF_FLAGGED(vict, AFF_FLYING) && !AFF_FLAGGED(ch, AFF_FLYING) && IS_HUMANOID(ch) && GET_LEVEL(ch) > 10 {
							do_fly(ch, nil, 0, 0)
							continue
						}
						if !AFF_FLAGGED(vict, AFF_FLYING) && AFF_FLAGGED(ch, AFF_FLYING) {
							do_fly(ch, nil, 0, 0)
							continue
						}
						do_punch(ch, &tar[0], 0, 0)
					}
					if !IS_HUMANOID(ch) {
						if AFF_FLAGGED(vict, AFF_FLYING) && !AFF_FLAGGED(ch, AFF_FLYING) && IS_HUMANOID(ch) && GET_LEVEL(ch) > 10 {
							do_fly(ch, nil, 0, 0)
							continue
						}
						if !AFF_FLAGGED(vict, AFF_FLYING) && AFF_FLAGGED(ch, AFF_FLYING) {
							do_fly(ch, nil, 0, 0)
							continue
						}
						if !AFF_FLAGGED(vict, AFF_HIDE) && !AFF_FLAGGED(vict, AFF_SNEAK) {
							act(libc.CString("@C$n @wgrowls viciously at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@C$n @wgrowls viciously at @c$N@w!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						} else {
							act(libc.CString("@C$n@w notices YOU.\n@C$n @wgrowls viciously at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@C$n@w notices @c$N@w.\n@C$n @wgrowls viciously at @c$N@w!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						}
						do_bite(ch, &tar[0], 0, 0)
					}
					found = 1
				}
			}
		}
		if ch.Original != nil && rand_number(1, 5) >= 4 {
			var original *char_data = ch.Original
			if original.Fighting != nil && ch.Fighting == nil {
				var (
					target [2048]byte
					targ   *char_data = original.Fighting
				)
				stdio.Sprintf(&target[0], "%s", targ.Name)
				if rand_number(1, 5) >= 4 {
					do_kick(ch, &target[0], 0, 0)
				} else if rand_number(1, 5) >= 4 {
					do_elbow(ch, &target[0], 0, 0)
				} else {
					do_punch(ch, &target[0], 0, 0)
				}
			}
		}
		if IS_HUMANOID(ch) && !MOB_FLAGGED(ch, MOB_NOKILL) {
			var (
				vict   *char_data
				next_v *char_data
				done   int = 0
			)
			for vict = world[ch.In_room].People; vict != nil; vict = next_v {
				next_v = vict.Next_in_room
				if vict == ch {
					continue
				}
				if IS_NPC(vict) && vict.Fighting != nil && done == 0 {
					if float64(vict.Hit) < float64(ch.Hit)*0.6 && axion_dice(0) >= 90 {
						act(libc.CString("@c$n@C rushes to @c$N's@C aid!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						var buf [2048]byte
						stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
						if int(ch.Chclass) == CLASS_KABITO || int(ch.Chclass) == CLASS_NAIL {
							do_heal(ch, &buf[0], 0, 0)
						} else {
							do_rescue(ch, &buf[0], 0, 0)
							if rand_number(1, 6) == 2 {
								var tar [2048]byte
								stdio.Sprintf(&tar[0], "%s", GET_NAME(vict.Fighting))
								do_kiblast(ch, &tar[0], 0, 0)
							} else if rand_number(1, 6) >= 4 {
								var tar [2048]byte
								stdio.Sprintf(&tar[0], "%s", GET_NAME(vict.Fighting))
								do_slam(ch, &tar[0], 0, 0)
							} else {
								var tar [2048]byte
								stdio.Sprintf(&tar[0], "%s", GET_NAME(vict.Fighting))
								do_punch(ch, &tar[0], 0, 0)
							}
						}
					}
				}
			}
		}
		if ch.Fighting == nil && rand_number(1, 20) >= 14 && IS_HUMANOID(ch) && !MOB_FLAGGED(ch, MOB_NOKILL) {
			var (
				vict   *char_data
				next_v *char_data
				done   int = 0
			)
			for vict = world[ch.In_room].People; vict != nil; vict = next_v {
				next_v = vict.Next_in_room
				if vict == ch {
					continue
				}
				if IS_NPC(vict) && vict.Fighting != nil && done == 0 {
					if float64(vict.Hit) < float64(ch.Hit)*0.6 && axion_dice(0) >= 70 {
						act(libc.CString("@c$n@C rushes to @c$N's@C aid!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						var buf [2048]byte
						stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
						if int(ch.Chclass) == CLASS_KABITO || int(ch.Chclass) == CLASS_NAIL {
							do_heal(ch, &buf[0], 0, 0)
							done = 1
						} else {
							do_rescue(ch, &buf[0], 0, 0)
							done = 1
						}
					}
				}
			}
		}
		if ch.Absorbby != nil && rand_number(1, 3) == 3 {
			do_escape(ch, nil, 0, 0)
		}
		if int(ch.Position) == POS_SLEEPING && rand_number(1, 3) == 3 {
			do_wake(ch, nil, 0, 0)
		}
		if libc.FuncAddr(GET_MOB_SPEC(ch)) == libc.FuncAddr(shop_keeper) {
			var diff libc.Time = 0
			diff = libc.GetTime(nil) - ch.Lastpl
			if diff > 86400 {
				var (
					sobj     *obj_data
					next_obj *obj_data
					shop_nr  int
					shopnr   int = -1
				)
				ch.Lastpl = libc.GetTime(nil)
				for shop_nr = 0; shop_nr <= top_shop; shop_nr++ {
					if shop_index[shop_nr].Keeper == ch.Nr {
						shopnr = shop_nr
					}
				}
				for sobj = ch.Carrying; sobj != nil; sobj = next_obj {
					next_obj = sobj.Next_content
					if sobj != nil && !shop_producing(sobj, shopnr) {
						ch.Gold += sobj.Cost
						extract_obj(sobj)
					}
				}
			}
		}
		if IS_HUMANOID(ch) && ch.Mob_specials.Memory != nil && !MOB_FLAGGED(ch, MOB_DUMMY) && !AFF_FLAGGED(ch, 18) {
			found = 0
			for vict = world[ch.In_room].People; vict != nil && found == 0; vict = vict.Next_in_room {
				if IS_NPC(vict) || !CAN_SEE(ch, vict) || PRF_FLAGGED(vict, PRF_NOHASSLE) {
					continue
				}
				if ch.Fighting != nil {
					continue
				}
				if ch.Hit <= ch.Max_hit/100 {
					continue
				}
				for names = ch.Mob_specials.Memory; names != nil && found == 0; names = names.Next {
					if int(names.Id) != int(vict.Idnum) {
						continue
					}
					found = 1
					act(libc.CString("'Hey!  You're the fiend that attacked me!!!', exclaims $n."), 0, ch, nil, nil, TO_ROOM)
					var tar [2048]byte
					stdio.Sprintf(&tar[0], "%s", GET_NAME(vict))
					do_punch(ch, &tar[0], 0, 0)
				}
			}
		}
		if ch.Fighting != nil && rand_number(1, 30) >= 25 {
			mob_taunt(ch)
		}
		if MOB_FLAGGED(ch, MOB_HELPER) && !AFF_FLAGGED(ch, AFF_BLIND) && !AFF_FLAGGED(ch, AFF_CHARM) {
			found = 0
			for vict = world[ch.In_room].People; vict != nil && found == 0; vict = vict.Next_in_room {
				if ch == vict || !IS_NPC(vict) || vict.Fighting == nil {
					continue
				}
				if IS_NPC(vict.Fighting) || ch == vict.Fighting {
					continue
				}
				if IS_HUMANOID(vict) {
					act(libc.CString("$n jumps to the aid of $N!"), 0, ch, nil, unsafe.Pointer(vict), TO_ROOM)
					var tar [2048]byte
					stdio.Sprintf(&tar[0], "%s", GET_NAME(vict.Fighting))
					do_punch(ch, &tar[0], 0, 0)
					found = 1
				}
			}
		}
		if int(ch.Chclass) == CLASS_KABITO {
			var shop_nr int
			found = 0
			for vict = world[ch.In_room].People; vict != nil && found == 0; vict = vict.Next_in_room {
				if libc.FuncAddr(GET_MOB_SPEC(vict)) == libc.FuncAddr(shop_keeper) {
					for shop_nr = 0; shop_nr <= top_shop; shop_nr++ {
						if shop_index[shop_nr].Keeper == vict.Nr {
							break
						}
					}
					if shop_nr <= top_shop {
						if ok_shop_room(shop_nr, vict.In_room) {
							if !IS_SET(shop_index[shop_nr].Bitvector, 1<<2) {
								found = 1
							}
						}
					}
				}
			}
			for vict = world[ch.In_room].People; vict != nil && found == 0; vict = vict.Next_in_room {
				if vict == ch {
					continue
				}
				if MOB_FLAGGED(ch, MOB_WIMPY) && AWAKE(vict) {
					continue
				}
				if !IS_HUMANOID(vict) {
					continue
				}
				if IS_NPC(vict) && MOB_FLAGGED(vict, MOB_NOKILL) {
					continue
				}
				if GET_MOB_VNUM(ch) == GET_MOB_VNUM(vict) {
					continue
				}
				if GET_LEVEL(ch) >= GET_LEVEL(vict) {
					if roll_skill(ch, SKILL_SLEIGHT_OF_HAND) != 0 {
						npc_steal(ch, vict)
						found = 1
					}
				}
			}
		}
	}
}
func mob_taunt(ch *char_data) {
	var message int = 1
	if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		return
	}
	if ch.Fighting == nil {
		return
	}
	var vict *char_data = ch.Fighting
	if vict == nil {
		return
	}
	if !IS_HUMANOID(ch) && !SUNKEN(ch.In_room) {
		message = rand_number(1, 12)
		switch message {
		case 1:
			act(libc.CString("@C$n@W growls viciously at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W growls viciously at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 2:
			act(libc.CString("@C$n@W snaps $s jaws at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W snaps $s jaws at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 3:
			act(libc.CString("@C$n@W is panting heavily from $s struggle with @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W is panting heavily from $s struggle with you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 4:
			act(libc.CString("@C$n@W circles around @c$N@W trying to get a better position!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W circles around you trying to find a weak spot!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 5:
			act(libc.CString("@C$n@W jumps up slightly in an attempt to threaten @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W jumps up slightly in an attempt to threaten you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 6:
			act(libc.CString("@C$n@W turns sideways while facing @c$N@W in an attempt to appear larger and more threatening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W turns sideways while facing you in an attempt to appear larger and more threatening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 7:
			act(libc.CString("@C$n@W roars with the full power of its lungs at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W roars with the full power of its lungs at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
			fallthrough
		case 8:
			act(libc.CString("@C$n@W staggers from the strain of fighting.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W staggers from the strain of fighting.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 9:
			act(libc.CString("@C$n@W slumps down for a moment before regaining $s guard against @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W slumps down for a moment before regaining $s guard against you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 10:
			act(libc.CString("@C$n's@W eyes dart around as $e seems to look for safe places to run.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n's@W eyes dart around as $e seems to look for safe places to run.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 11:
			act(libc.CString("@C$n@W jumps past @c$N@W before turning and facing $M again!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W jumps past you before turning and facing you again!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		default:
			act(libc.CString("@C$n@W watches @c$N@W with a threatening gaze while $e looks for a weakness!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W watches you with a threatening gaze while $e looks for a weakness!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		}
	} else if !IS_HUMANOID(ch) {
		message = rand_number(1, 7)
		switch message {
		case 1:
			act(libc.CString("@C$n@W snaps $s jaws at @c$N@W which causes a torrent of bubbles to float upward!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W snaps $s jaws at you which causes a torrent of bubbles to float upward!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 2:
			act(libc.CString("@C$n@W thrashes around in the water!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W thrashes around in the water!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 3:
			act(libc.CString("@C$n@W swims past @c$N@W before turning and facing $M again!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W swims past you before turning and facing you again!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 4:
			act(libc.CString("@C$n@W begins to slowly circle @c$N@W while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W begins to slowly circle you while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 5:
			act(libc.CString("@C$n@W swims backward in an attempt to gain a safe distance from @C$N's@W aggression.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W swims backward in an attempt to gain a safe distance from you.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		case 6:
			act(libc.CString("@C$n@W swims toward the side of @C$N@W in an attempt to flank $M!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W swims toward the side of you in an attempt to flank you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		default:
			act(libc.CString("@C$n@W swims upward before darting down past @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@C$n@W swims upward before darting down past you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		}
	} else if !MOB_FLAGGED(ch, MOB_DUMMY) {
		message = rand_number(1, 10)
		if !SUNKEN(ch.In_room) {
			if AFF_FLAGGED(ch, AFF_FLYING) {
				switch message {
				case 1:
					act(libc.CString("@C$n@W flies around @c$N@W slowly while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W flies around you slowly while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 2:
					act(libc.CString("@C$n@W floats slowly while scowling at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W floats slowly while scowling at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 3:
					act(libc.CString("@C$n@W spits at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W spits at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 4:
					act(libc.CString("@C$n@W looks at @c$N@W as if $e is weighing $s options.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W looks at you as if $e is weighing $s options.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 5:
					act(libc.CString("@C$n@W scowls at @c$N@W while changing $s position carefully!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W scowls at you while changing $s position carefully!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 6:
					act(libc.CString("@C$n@W flips backward a short way away from @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W flips backward a short way away from you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 7:
					act(libc.CString("@C$n@W moves slowly to the side of @c$N@W while watching $M carefully.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W moves slowly to the side of you while watching you carefully.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 8:
					act(libc.CString("@C$n@W flexes $s arms in an attempt to threaten @C$N@W.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W flexes $s arms threaten in an attempt to threaten you@W.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 9:
					act(libc.CString("@C$n@W raises an arm in front of $s body as a defense.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W raises an arm in front of $s body as a defense.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				default:
					act(libc.CString("@C$n@W feints a punch toward @c$N@W that misses by a mile.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W feints a punch toward you that misses by a mile.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				}
			} else {
				message = rand_number(1, 13)
				switch message {
				case 1:
					act(libc.CString("@C$n@W shuffles around @c$N@W slowly while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W shuffles around you slowly while looking for an opening!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 2:
					act(libc.CString("@C$n@W scowls @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W scowls at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 3:
					if int(ch.Race) == RACE_ANDROID {
						act(libc.CString("@C$n@W has sparks come off them that land on @c$N@W!@n@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						act(libc.CString("@C$n@W has sparks come off them that land on you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
					} else {
						act(libc.CString("@C$n@W spits at @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						act(libc.CString("@C$n@W spits at you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
					}
				case 4:
					act(libc.CString("@C$n@W looks at @c$N@W as if $e is weighing $s options.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W looks at you as if $e is weighing $s options.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 5:
					act(libc.CString("@C$n@W scowls at @c$N@W while changing $s position carefully!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W scowls at you while changing $s position carefully!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 6:
					act(libc.CString("@C$n@W flips backward a short way away from @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W flips backward a short way away from you!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 7:
					act(libc.CString("@C$n@W moves slowly to the side of @c$N@W while watching $M carefully.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W moves slowly to the side of you while watching you carefully.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 8:
					act(libc.CString("@C$n@W crouches down cautiously.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W crouches down cautiously.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 9:
					act(libc.CString("@C$n@W moves $s feet slowly to achieve a better balance.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W moves $s feet slowly to achieve a better balance.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 10:
					act(libc.CString("@C$n@W leaps to a more defensible spot.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W leaps to a more defensible spot.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 11:
					act(libc.CString("@C$n@W runs a short distance away before skidding to a halt and resuming $s fighting stance.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W runs a short distance away before skidding to a halt and resuming $s fighting stance.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				case 12:
					act(libc.CString("@C$n@W stands up to $s full height and glares at @C$N@W with burning eyes.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W stands up to $s full height and glares at you with intense burning eyes.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				default:
					act(libc.CString("@C$n@W feints a punch toward @c$N@W that misses by a mile.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("@C$n@W feints a punch toward you that misses by a mile.@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
				}
			}
		}
	}
}
func remember(ch *char_data, victim *char_data) {
	var (
		tmp     *memory_rec_struct
		present bool = false
	)
	if !IS_NPC(ch) || IS_NPC(victim) || PRF_FLAGGED(victim, PRF_NOHASSLE) {
		return
	}
	for tmp = ch.Mob_specials.Memory; tmp != nil && !present; tmp = tmp.Next {
		if int(tmp.Id) == int(victim.Idnum) {
			present = true
		}
	}
	if !present && !MOB_FLAGGED(ch, MOB_SPAR) && !PLR_FLAGGED(victim, PLR_SPAR) {
		tmp = new(memory_rec_struct)
		tmp.Next = ch.Mob_specials.Memory
		tmp.Id = victim.Idnum
		ch.Mob_specials.Memory = tmp
	}
}
func forget(ch *char_data, victim *char_data) {
	var (
		curr *memory_rec_struct
		prev *memory_rec_struct = nil
	)
	if (func() *memory_rec_struct {
		curr = ch.Mob_specials.Memory
		return curr
	}()) == nil {
		return
	}
	for curr != nil && int(curr.Id) != int(victim.Idnum) {
		prev = curr
		curr = curr.Next
	}
	if curr == nil {
		return
	}
	if curr == ch.Mob_specials.Memory {
		ch.Mob_specials.Memory = curr.Next
	} else {
		prev.Next = curr.Next
	}
	libc.Free(unsafe.Pointer(curr))
}
func clearMemory(ch *char_data) {
	var (
		curr *memory_rec_struct
		next *memory_rec_struct
	)
	curr = ch.Mob_specials.Memory
	for curr != nil {
		next = curr.Next
		libc.Free(unsafe.Pointer(curr))
		curr = next
	}
	ch.Mob_specials.Memory = nil
}
func aggressive_mob_on_a_leash(slave *char_data, master *char_data, attack *char_data) bool {
	var (
		snarl_cmd int
		dieroll   int
	)
	if master == nil || !AFF_FLAGGED(slave, AFF_CHARM) {
		return false
	}
	if snarl_cmd == 0 {
		snarl_cmd = find_command(libc.CString("snarl"))
	}
	dieroll = rand_number(1, 20)
	if dieroll != 1 && (dieroll == 20 || dieroll > 10-int(master.Aff_abils.Cha)+int(slave.Aff_abils.Intel)) {
		if snarl_cmd > 0 && attack != nil && rand_number(0, 3) == 0 {
			var victbuf [21]byte
			libc.StrNCpy(&victbuf[0], GET_NAME(attack), int(21))
			victbuf[21-1] = '\x00'
			do_action(slave, &victbuf[0], snarl_cmd, 0)
		}
		return true
	}
	return false
}
