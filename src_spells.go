package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const SUMMON_FAIL = "You failed.\r\n"

func spell_create_water(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var water int
	if ch == nil || obj == nil {
		return
	}
	if int(obj.Type_flag) == ITEM_DRINKCON {
		if (obj.Value[VAL_DRINKCON_LIQUID]) != LIQ_WATER && (obj.Value[VAL_DRINKCON_HOWFULL]) != 0 {
			name_from_drinkcon(obj)
			obj.Value[VAL_DRINKCON_LIQUID] = LIQ_SLIME
			name_to_drinkcon(obj, LIQ_SLIME)
		} else {
			water = int(MAX(int64((obj.Value[VAL_DRINKCON_CAPACITY])-(obj.Value[VAL_DRINKCON_HOWFULL])), 0))
			if water > 0 {
				if (obj.Value[VAL_DRINKCON_HOWFULL]) >= 0 {
					name_from_drinkcon(obj)
				}
				obj.Value[VAL_DRINKCON_LIQUID] = LIQ_WATER
				obj.Value[VAL_DRINKCON_HOWFULL] += water
				name_to_drinkcon(obj, LIQ_WATER)
				weight_change_object(obj, water)
				act(libc.CString("$p is filled."), 0, ch, obj, nil, TO_CHAR)
			}
		}
	}
}
func spell_recall(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if victim == nil || IS_NPC(victim) {
		return
	}
	act(libc.CString("$n disappears."), 1, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, real_room(config_info.Room_nums.Mortal_start_room))
	act(libc.CString("$n appears in the middle of the room."), 1, victim, nil, nil, TO_ROOM)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_teleport(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var to_room int
	if victim == nil || IS_NPC(victim) {
		return
	}
	for {
		to_room = rand_number(0, top_of_world)
		if !ROOM_FLAGGED(to_room, uint32(int32(int(ROOM_PRIVATE|ROOM_DEATH)|ROOM_GODROOM))) {
			break
		}
	}
	act(libc.CString("$n slowly fades out of existence and is gone."), 0, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, to_room)
	act(libc.CString("$n slowly fades into existence."), 0, victim, nil, nil, TO_ROOM)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_summon(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if ch == nil || victim == nil {
		return
	}
	if GET_LEVEL(victim) > level+3 {
		send_to_char(ch, libc.CString("%s"), SUMMON_FAIL)
		return
	}
	if config_info.Play.Pk_allowed == 0 {
		if MOB_FLAGGED(victim, MOB_AGGRESSIVE) {
			act(libc.CString("As the words escape your lips and $N travels\r\nthrough time and space towards you, you realize that $E is\r\naggressive and might harm you, so you wisely send $M back."), 0, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			return
		}
		if !IS_NPC(victim) && !PRF_FLAGGED(victim, PRF_SUMMONABLE) && !PLR_FLAGGED(victim, PLR_KILLER) {
			send_to_char(victim, libc.CString("%s just tried to summon you to: %s.\r\n%s failed because you have summon protection on.\r\nType NOSUMMON to allow other players to summon you.\r\n"), GET_NAME(ch), world[ch.In_room].Name, func() string {
				if int(ch.Sex) == SEX_MALE {
					return "He"
				}
				return "She"
			}())
			send_to_char(ch, libc.CString("You failed because %s has summon protection on.\r\n"), GET_NAME(victim))
			mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("%s failed summoning %s to %s."), GET_NAME(ch), GET_NAME(victim), world[ch.In_room].Name)
			return
		}
	}
	if MOB_FLAGGED(victim, MOB_NOSUMMON) || IS_NPC(victim) && mag_newsaves(ch, victim, SPELL_SUMMON, level, int(ch.Aff_abils.Intel)) {
		send_to_char(ch, libc.CString("%s"), SUMMON_FAIL)
		return
	}
	act(libc.CString("$n disappears suddenly."), 1, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, ch.In_room)
	act(libc.CString("$n arrives suddenly."), 1, victim, nil, nil, TO_ROOM)
	act(libc.CString("$n has summoned you!"), 0, ch, nil, unsafe.Pointer(victim), TO_VICT)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_locate_object(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		i    *obj_data
		name [2048]byte
		j    int
	)
	if obj == nil {
		send_to_char(ch, libc.CString("You sense nothing.\r\n"))
		return
	}
	strlcpy(&name[0], fname(obj.Name), uint64(2048))
	j = level / 2
	for i = object_list; i != nil && j > 0; i = i.Next {
		if !isname(&name[0], i.Name) {
			continue
		}
		send_to_char(ch, libc.CString("%c%s"), unicode.ToUpper(rune(*i.Short_description)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_description), 1)))
		if i.Carried_by != nil {
			send_to_char(ch, libc.CString(" is being carried by %s.\r\n"), PERS(i.Carried_by, ch))
		} else if i.In_room != int(-1) {
			send_to_char(ch, libc.CString(" is in %s.\r\n"), world[i.In_room].Name)
		} else if i.In_obj != nil {
			send_to_char(ch, libc.CString(" is in %s.\r\n"), i.In_obj.Short_description)
		} else if i.Worn_by != nil {
			send_to_char(ch, libc.CString(" is being worn by %s.\r\n"), PERS(i.Worn_by, ch))
		} else {
			send_to_char(ch, libc.CString("'s location is uncertain.\r\n"))
		}
		j--
	}
	if j == level/2 {
		send_to_char(ch, libc.CString("You sense nothing.\r\n"))
	}
}
func spell_charm(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var af affected_type
	if victim == nil || ch == nil {
		return
	}
	if victim == ch {
		send_to_char(ch, libc.CString("You like yourself even better!\r\n"))
	} else if !IS_NPC(victim) && !PRF_FLAGGED(victim, PRF_SUMMONABLE) {
		send_to_char(ch, libc.CString("You fail because SUMMON protection is on!\r\n"))
	} else if AFF_FLAGGED(victim, AFF_SANCTUARY) {
		send_to_char(ch, libc.CString("Your victim is protected by sanctuary!\r\n"))
	} else if MOB_FLAGGED(victim, MOB_NOCHARM) {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else if AFF_FLAGGED(ch, AFF_CHARM) {
		send_to_char(ch, libc.CString("You can't have any followers of your own!\r\n"))
	} else if AFF_FLAGGED(victim, AFF_CHARM) || level < GET_LEVEL(victim) {
		send_to_char(ch, libc.CString("You fail.\r\n"))
	} else if config_info.Play.Pk_allowed == 0 && !IS_NPC(victim) {
		send_to_char(ch, libc.CString("You fail - shouldn't be doing it anyway.\r\n"))
	} else if int(victim.Race) == RACE_SAIYAN && rand_number(1, 100) <= 90 {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else if circle_follow(victim, ch) {
		send_to_char(ch, libc.CString("Sorry, following in circles cannot be allowed.\r\n"))
	} else if mag_newsaves(ch, victim, SPELL_CHARM, level, int(ch.Aff_abils.Intel)) {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else {
		if victim.Master != nil {
			stop_follower(victim)
		}
		add_follower(victim, ch)
		victim.Master_id = ch.Idnum
		af.Type = SPELL_CHARM
		af.Duration = 24 * 2
		if int(ch.Aff_abils.Cha) != 0 {
			af.Duration *= int16(ch.Aff_abils.Cha)
		}
		if int(victim.Aff_abils.Intel) != 0 {
			af.Duration /= int16(victim.Aff_abils.Intel)
		}
		af.Modifier = 0
		af.Location = 0
		af.Bitvector = AFF_CHARM
		affect_to_char(victim, &af)
		act(libc.CString("Isn't $n just such a nice fellow?"), 0, ch, nil, unsafe.Pointer(victim), TO_VICT)
		if IS_NPC(victim) {
			REMOVE_BIT_AR(victim.Act[:], MOB_SPEC)
		}
	}
}
func spell_identify(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		i     int
		found int
		len_  uint64
	)
	if obj != nil {
		var (
			bitbuf [64936]byte
			buf2   [64936]byte
		)
		sprinttype(int(obj.Type_flag), item_types[:], &bitbuf[0], uint64(64936))
		send_to_char(ch, libc.CString("You feel informed:\r\nObject '%s', Item type: %s\r\n"), obj.Short_description, &bitbuf[0])
		if obj.Bitvector[0] > 0 || obj.Bitvector[1] > 0 || obj.Bitvector[2] > 0 || obj.Bitvector[3] > 0 {
			sprintbitarray(obj.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &bitbuf[0])
			send_to_char(ch, libc.CString("Item will give you following abilities:  %s\r\n"), &bitbuf[0])
		}
		sprintbitarray(obj.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, &bitbuf[0])
		send_to_char(ch, libc.CString("Item is: %s\r\n"), &bitbuf[0])
		send_to_char(ch, libc.CString("Weight: %lld, Value: %d, Rent: %d, Min Level: %d\r\n"), obj.Weight, obj.Cost, obj.Cost_per_day, obj.Level)
		switch obj.Type_flag {
		case ITEM_SCROLL:
			fallthrough
		case ITEM_POTION:
			len_ = uint64(func() int {
				i = 0
				return i
			}())
			if (obj.Value[VAL_SCROLL_SPELL1]) >= 1 {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL1]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			if (obj.Value[VAL_SCROLL_SPELL2]) >= 1 && len_ < uint64(64936) {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL2]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			if (obj.Value[VAL_SCROLL_SPELL3]) >= 1 && len_ < uint64(64936) {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL3]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			send_to_char(ch, libc.CString("This %s casts: %s\r\n"), item_types[int(obj.Type_flag)], &bitbuf[0])
		case ITEM_WAND:
			fallthrough
		case ITEM_STAFF:
			send_to_char(ch, libc.CString("This %s casts: %s\r\nIt has %d maximum charge%s and %d remaining.\r\n"), item_types[int(obj.Type_flag)], skill_name(obj.Value[VAL_WAND_SPELL]), obj.Value[VAL_WAND_MAXCHARGES], func() string {
				if (obj.Value[VAL_WAND_MAXCHARGES]) == 1 {
					return ""
				}
				return "s"
			}(), obj.Value[VAL_WAND_CHARGES])
		case ITEM_WEAPON:
			send_to_char(ch, libc.CString("Damage Dice is '%dD%d' for an average per-round damage of %.1f.\r\n"), obj.Value[VAL_WEAPON_DAMDICE], obj.Value[VAL_WEAPON_DAMSIZE], (float64((obj.Value[VAL_WEAPON_DAMSIZE])+1)/2.0)*float64(obj.Value[VAL_WEAPON_DAMDICE]))
		case ITEM_ARMOR:
			send_to_char(ch, libc.CString("AC-apply is %.1f\r\n"), (float32(obj.Value[VAL_ARMOR_APPLYAC]))/10)
		}
		found = 0
		for i = 0; i < MAX_OBJ_AFFECT; i++ {
			if obj.Affected[i].Location != APPLY_NONE && obj.Affected[i].Modifier != 0 {
				if found == 0 {
					send_to_char(ch, libc.CString("Can affect you as :\r\n"))
					found = 1
				}
				sprinttype(obj.Affected[i].Location, apply_types[:], &bitbuf[0], uint64(64936))
				switch obj.Affected[i].Location {
				case APPLY_FEAT:
					stdio.Snprintf(&buf2[0], int(64936), " (%s)", feat_list[obj.Affected[i].Specific].Name)
				case APPLY_SKILL:
					stdio.Snprintf(&buf2[0], int(64936), " (%s)", spell_info[obj.Affected[i].Specific].Name)
				default:
					buf2[0] = 0
				}
				send_to_char(ch, libc.CString("   Affects: %s%s By %d\r\n"), &bitbuf[0], &buf2[0], obj.Affected[i].Modifier)
			}
		}
	}
}
func spell_enchant_weapon(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var i int
	if ch == nil || obj == nil {
		return
	}
	if int(obj.Type_flag) != ITEM_WEAPON || OBJ_FLAGGED(obj, ITEM_MAGIC) {
		return
	}
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location != APPLY_NONE {
			return
		}
	}
	SET_BIT_AR(obj.Extra_flags[:], ITEM_MAGIC)
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location == APPLY_NONE {
			obj.Affected[i].Location = APPLY_ACCURACY
			obj.Affected[i].Modifier = int(libc.BoolToInt(level >= 18)) + 1
			break
		}
	}
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location == APPLY_NONE {
			obj.Affected[i].Location = APPLY_DAMAGE
			obj.Affected[i].Modifier = int(libc.BoolToInt(level >= 20)) + 1
			break
		}
	}
	if IS_GOOD(ch) {
		SET_BIT_AR(obj.Extra_flags[:], ITEM_ANTI_EVIL)
		act(libc.CString("$p glows blue."), 0, ch, obj, nil, TO_CHAR)
	} else if IS_EVIL(ch) {
		SET_BIT_AR(obj.Extra_flags[:], ITEM_ANTI_GOOD)
		act(libc.CString("$p glows red."), 0, ch, obj, nil, TO_CHAR)
	} else {
		act(libc.CString("$p glows yellow."), 0, ch, obj, nil, TO_CHAR)
	}
}
func spell_detect_poison(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if victim != nil {
		if victim == ch {
			if AFF_FLAGGED(victim, AFF_POISON) {
				send_to_char(ch, libc.CString("You can sense poison in your blood.\r\n"))
			} else {
				send_to_char(ch, libc.CString("You feel healthy.\r\n"))
			}
		} else {
			if AFF_FLAGGED(victim, AFF_POISON) {
				act(libc.CString("You sense that $E is poisoned."), 0, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			} else {
				act(libc.CString("You sense that $E is healthy."), 0, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			}
		}
	}
	if obj != nil {
		switch obj.Type_flag {
		case ITEM_DRINKCON:
			fallthrough
		case ITEM_FOUNTAIN:
			fallthrough
		case ITEM_FOOD:
			if (obj.Value[VAL_FOOD_POISON]) != 0 {
				act(libc.CString("You sense that $p has been contaminated."), 0, ch, obj, nil, TO_CHAR)
			} else {
				act(libc.CString("You sense that $p is safe for consumption."), 0, ch, obj, nil, TO_CHAR)
			}
		default:
			send_to_char(ch, libc.CString("You sense that it should not be consumed.\r\n"))
		}
	}
}
func spell_portal(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		portal  *obj_data
		tportal *obj_data
		rm      *room_data
	)
	rm = &world[victim.In_room]
	if ch == nil || victim == nil {
		return
	}
	if !can_edit_zone(ch, rm.Zone) && ZONE_FLAGGED(rm.Zone, ZONE_QUEST) {
		send_to_char(ch, libc.CString("That target is in a quest zone.\r\n"))
		return
	}
	if ZONE_FLAGGED(rm.Zone, ZONE_CLOSED) && ch.Admlevel < ADMLVL_IMMORT {
		send_to_char(ch, libc.CString("That target is in a closed zone.\r\n"))
		return
	}
	if ZONE_FLAGGED(rm.Zone, ZONE_NOIMMORT) && ch.Admlevel < ADMLVL_GRGOD {
		send_to_char(ch, libc.CString("That target is in a zone closed to all.\r\n"))
		return
	}
	portal = read_object(portal_object, VIRTUAL)
	portal.Value[VAL_PORTAL_DEST] = int(libc.BoolToInt(GET_ROOM_VNUM(victim.In_room)))
	portal.Value[VAL_PORTAL_HEALTH] = 100
	portal.Value[VAL_PORTAL_MAXHEALTH] = 100
	portal.Timer = level / 10
	add_unique_id(portal)
	obj_to_room(portal, ch.In_room)
	act(libc.CString("$n opens a portal in thin air."), 1, ch, nil, nil, TO_ROOM)
	act(libc.CString("You open a portal out of thin air."), 1, ch, nil, nil, TO_CHAR)
	tportal = read_object(portal_object, VIRTUAL)
	tportal.Value[VAL_PORTAL_DEST] = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
	tportal.Value[VAL_PORTAL_HEALTH] = 100
	tportal.Value[VAL_PORTAL_MAXHEALTH] = 100
	tportal.Timer = level / 10
	add_unique_id(portal)
	obj_to_room(tportal, victim.In_room)
	act(libc.CString("A shimmering portal appears out of thin air."), 1, victim, nil, nil, TO_ROOM)
	act(libc.CString("A shimmering portal opens here for you."), 1, victim, nil, nil, TO_CHAR)
}
func art_abundant_step(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		steps    int
		i        int = 0
		j        int
		rep      int
		max      int
		r        int
		nextroom int
		buf      [2048]byte
		tc       int8
		p        *byte
	)
	steps = 0
	r = ch.In_room
	p = arg
	max = ((ch.Chclasses[CLASS_KABITO])+(ch.Epicclasses[CLASS_KABITO]))/2 + 10
	for p != nil && *p != 0 && !unicode.IsDigit(rune(*p)) && !libc.IsAlpha(rune(*p)) {
		p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
	}
	if p == nil || *p == 0 {
		send_to_char(ch, libc.CString("You must give directions from your current location. Examples:\r\n  w w nw n e\r\n  2w nw n e\r\n"))
		return
	}
	for *p != 0 {
		for *p != 0 && !unicode.IsDigit(rune(*p)) && !libc.IsAlpha(rune(*p)) {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
		if unicode.IsDigit(rune(*p)) {
			rep = libc.Atoi(libc.GoString(p))
			for unicode.IsDigit(rune(*p)) {
				p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
			}
		} else {
			rep = 1
		}
		if libc.IsAlpha(rune(*p)) {
			for i = 0; libc.IsAlpha(rune(*p)); func() *byte {
				i++
				return func() *byte {
					p := &p
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()
			}() {
				buf[i] = byte(int8(unicode.ToLower(rune(*p))))
			}
			j = i
			tc = int8(buf[i])
			buf[i] = 0
			for i = 1; libc.FuncAddr(complete_cmd_info[i].Command_pointer) == libc.FuncAddr(do_move) && libc.StrCmp(complete_cmd_info[i].Sort_as, &buf[0]) != 0; i++ {
			}
			if libc.FuncAddr(complete_cmd_info[i].Command_pointer) == libc.FuncAddr(do_move) {
				i = complete_cmd_info[i].Subcmd - 1
			} else {
				i = -1
			}
			buf[j] = byte(tc)
		}
		if i > -1 {
			for func() int {
				p := &rep
				x := *p
				*p--
				return x
			}() != 0 {
				if func() int {
					p := &steps
					*p++
					return *p
				}() > max {
					break
				}
				if (world[r].Dir_option[i]) == nil {
					send_to_char(ch, libc.CString("Invalid step. Skipping.\r\n"))
					break
				}
				nextroom = (world[r].Dir_option[i]).To_room
				if nextroom == int(-1) {
					break
				}
				r = nextroom
			}
		}
		if steps > max {
			break
		}
	}
	send_to_char(ch, libc.CString("Your will bends reality as you travel through the ethereal plane.\r\n"))
	act(libc.CString("$n is suddenly absent."), 1, ch, nil, nil, TO_ROOM)
	char_from_room(ch)
	char_to_room(ch, r)
	act(libc.CString("$n is suddenly present."), 1, ch, nil, nil, TO_ROOM)
	look_at_room(ch.In_room, ch, 0)
	return
}
func roll_skill(ch *char_data, snum int) int {
	var (
		roll  int
		skval int
		i     int
	)
	if !IS_NPC(ch) {
		skval = GET_SKILL(ch, snum)
		if SKILL_SPOT == snum {
			if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				skval += 5
			}
		} else if SKILL_HIDE == snum {
			if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
				skval += 5
			} else if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 5 || (ch.Genome[1]) == 5) {
				skval += 10
			}
		}
	} else if IS_NPC(ch) {
		var numb int = 0
		if GET_LEVEL(ch) <= 10 {
			numb = rand_number(15, 30)
		}
		if GET_LEVEL(ch) <= 20 {
			numb = rand_number(20, 40)
		}
		if GET_LEVEL(ch) <= 30 {
			numb = rand_number(40, 60)
		}
		if GET_LEVEL(ch) <= 60 {
			numb = rand_number(60, 80)
		}
		if GET_LEVEL(ch) <= 80 {
			numb = rand_number(70, 90)
		}
		if GET_LEVEL(ch) <= 90 {
			numb = rand_number(80, 95)
		}
		if GET_LEVEL(ch) <= 100 {
			numb = rand_number(90, 100)
		}
		skval = numb
	}
	if snum == SKILL_SPOT && GET_SKILL(ch, SKILL_LISTEN) != 0 {
		skval += GET_SKILL(ch, SKILL_LISTEN) / 10
	}
	if snum < 0 || snum >= SKILL_TABLE_SIZE {
		return 0
	}
	if IS_SET(uint32(int32(spell_info[snum].Skilltype)), 1<<0) {
		for func() int {
			i = 0
			return func() int {
				roll = 0
				return roll
			}()
		}(); i < NUM_CLASSES; i++ {
			if ((ch.Chclasses[i])+(ch.Epicclasses[i])) != 0 && spell_info[snum].Min_level[i] < ((ch.Chclasses[i])+(ch.Epicclasses[i])) {
				roll += (ch.Chclasses[i]) + (ch.Epicclasses[i])
			}
		}
		return roll + rand_number(1, 20)
	} else if IS_SET(uint32(int32(spell_info[snum].Skilltype)), 1<<1) {
		if skval == 0 && IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<0) {
			return -1
		} else {
			roll = skval
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<1) {
				roll += int(ability_mod_value(int(ch.Aff_abils.Str)))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<2) {
				roll += int(dex_mod_capped(ch))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<3) {
				roll += int(ability_mod_value(int(ch.Aff_abils.Con)))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<4) {
				roll += int(ability_mod_value(int(ch.Aff_abils.Intel)))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<5) {
				roll += int(ability_mod_value(int(ch.Aff_abils.Wis)))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<6) {
				roll += int(ability_mod_value(int(ch.Aff_abils.Cha)))
			}
			if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<8) {
				roll -= int(ch.Armorcheckall)
			} else if IS_SET(uint32(int32(spell_info[snum].Flags)), 1<<7) {
				roll -= int(ch.Armorcheck)
			}
			return roll + rand_number(1, 20)
		}
	} else {
		basic_mud_log(libc.CString("Trying to roll uncategorized skill/spell #%d for %s"), snum, GET_NAME(ch))
		return 0
	}
}
func roll_resisted(actor *char_data, sact int, resistor *char_data, sres int) int {
	return int(libc.BoolToInt(roll_skill(actor, sact) >= roll_skill(resistor, sres)))
}
