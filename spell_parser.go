package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

var spell_info [1000]spell_info_type

type syllable struct {
	Org  *byte
	News *byte
}

var syls [55]syllable = [55]syllable{{Org: libc.CString(" "), News: libc.CString(" ")}, {Org: libc.CString("ar"), News: libc.CString("abra")}, {Org: libc.CString("ate"), News: libc.CString("i")}, {Org: libc.CString("cau"), News: libc.CString("kada")}, {Org: libc.CString("blind"), News: libc.CString("nose")}, {Org: libc.CString("bur"), News: libc.CString("mosa")}, {Org: libc.CString("cu"), News: libc.CString("judi")}, {Org: libc.CString("de"), News: libc.CString("oculo")}, {Org: libc.CString("dis"), News: libc.CString("mar")}, {Org: libc.CString("ect"), News: libc.CString("kamina")}, {Org: libc.CString("en"), News: libc.CString("uns")}, {Org: libc.CString("gro"), News: libc.CString("cra")}, {Org: libc.CString("light"), News: libc.CString("dies")}, {Org: libc.CString("lo"), News: libc.CString("hi")}, {Org: libc.CString("magi"), News: libc.CString("kari")}, {Org: libc.CString("mon"), News: libc.CString("bar")}, {Org: libc.CString("mor"), News: libc.CString("zak")}, {Org: libc.CString("move"), News: libc.CString("sido")}, {Org: libc.CString("ness"), News: libc.CString("lacri")}, {Org: libc.CString("ning"), News: libc.CString("illa")}, {Org: libc.CString("per"), News: libc.CString("duda")}, {Org: libc.CString("ra"), News: libc.CString("gru")}, {Org: libc.CString("re"), News: libc.CString("candus")}, {Org: libc.CString("son"), News: libc.CString("sabru")}, {Org: libc.CString("tect"), News: libc.CString("infra")}, {Org: libc.CString("tri"), News: libc.CString("cula")}, {Org: libc.CString("ven"), News: libc.CString("nofo")}, {Org: libc.CString("word of"), News: libc.CString("inset")}, {Org: libc.CString("a"), News: libc.CString("i")}, {Org: libc.CString("b"), News: libc.CString("v")}, {Org: libc.CString("c"), News: libc.CString("q")}, {Org: libc.CString("d"), News: libc.CString("m")}, {Org: libc.CString("e"), News: libc.CString("o")}, {Org: libc.CString("f"), News: libc.CString("y")}, {Org: libc.CString("g"), News: libc.CString("t")}, {Org: libc.CString("h"), News: libc.CString("p")}, {Org: libc.CString("i"), News: libc.CString("u")}, {Org: libc.CString("j"), News: libc.CString("y")}, {Org: libc.CString("k"), News: libc.CString("t")}, {Org: libc.CString("l"), News: libc.CString("r")}, {Org: libc.CString("m"), News: libc.CString("w")}, {Org: libc.CString("n"), News: libc.CString("b")}, {Org: libc.CString("o"), News: libc.CString("a")}, {Org: libc.CString("p"), News: libc.CString("s")}, {Org: libc.CString("q"), News: libc.CString("d")}, {Org: libc.CString("r"), News: libc.CString("f")}, {Org: libc.CString("s"), News: libc.CString("g")}, {Org: libc.CString("t"), News: libc.CString("h")}, {Org: libc.CString("u"), News: libc.CString("e")}, {Org: libc.CString("v"), News: libc.CString("z")}, {Org: libc.CString("w"), News: libc.CString("x")}, {Org: libc.CString("x"), News: libc.CString("n")}, {Org: libc.CString("y"), News: libc.CString("l")}, {Org: libc.CString("z"), News: libc.CString("k")}, {Org: libc.CString(""), News: libc.CString("")}}
var unused_spellname *byte = libc.CString("!UNUSED!")

func mag_manacost(ch *char_data, spellnum int) int {
	var whichclass int
	_ = whichclass
	var i int
	var min int
	var tval int
	if config_info.Advance.Allow_multiclass != 0 {
		min = MAX(spell_info[spellnum].Mana_max-spell_info[spellnum].Mana_change*(GET_LEVEL(ch)-spell_info[spellnum].Min_level[(ch.Chclasses[ch.Chclass])+(ch.Epicclasses[ch.Chclass])]), spell_info[spellnum].Mana_min)
		whichclass = int(ch.Chclass)
		for i = 0; i < NUM_CLASSES; i++ {
			if ((ch.Chclasses[i]) + (ch.Epicclasses[i])) == 0 {
				continue
			}
			tval = MAX(spell_info[spellnum].Mana_max-spell_info[spellnum].Mana_change*(((ch.Chclasses[i])+(ch.Epicclasses[i]))-spell_info[spellnum].Min_level[i]), spell_info[spellnum].Mana_min)
			if tval < min {
				min = tval
				whichclass = i
			}
		}
		return min
	} else {
		return MAX(spell_info[spellnum].Mana_max-spell_info[spellnum].Mana_change*(GET_LEVEL(ch)-spell_info[spellnum].Min_level[int(ch.Chclass)]), spell_info[spellnum].Mana_min)
	}
}
func mag_kicost(ch *char_data, spellnum int) int {
	var whichclass int
	_ = whichclass
	var i int
	var min int
	var tval int
	if config_info.Advance.Allow_multiclass != 0 {
		min = MAX(spell_info[spellnum].Ki_max-spell_info[spellnum].Ki_change*(GET_LEVEL(ch)-spell_info[spellnum].Min_level[(ch.Chclasses[ch.Chclass])+(ch.Epicclasses[ch.Chclass])]), spell_info[spellnum].Ki_min)
		whichclass = int(ch.Chclass)
		for i = 0; i < NUM_CLASSES; i++ {
			if ((ch.Chclasses[i]) + (ch.Epicclasses[i])) == 0 {
				continue
			}
			tval = MAX(spell_info[spellnum].Ki_max-spell_info[spellnum].Ki_change*(((ch.Chclasses[i])+(ch.Epicclasses[i]))-spell_info[spellnum].Min_level[i]), spell_info[spellnum].Ki_min)
			if tval < min {
				min = tval
				whichclass = i
			}
		}
		return min
	} else {
		return MAX(spell_info[spellnum].Ki_max-spell_info[spellnum].Ki_change*(GET_LEVEL(ch)-spell_info[spellnum].Min_level[int(ch.Chclass)]), spell_info[spellnum].Ki_min)
	}
}
func mag_nextstrike(level int, caster *char_data, spellnum int) {
	if caster == nil {
		return
	}
	if caster.Actq != nil {
		send_to_char(caster, libc.CString("You can't perform more than one special attack at a time!"))
		return
	}
	caster.Actq = new(queued_act)
	caster.Actq.Level = level
	caster.Actq.Spellnum = spellnum
}
func say_spell(ch *char_data, spellnum int, tch *char_data, tobj *obj_data) {
	var (
		lbuf   [256]byte
		buf    [256]byte
		buf1   [256]byte
		buf2   [256]byte
		format *byte
		i      *char_data
		j      int
		ofs    int = 0
	)
	buf[0] = '\x00'
	strlcpy(&lbuf[0], skill_name(spellnum), uint64(256))
	for lbuf[ofs] != 0 {
		for j = 0; *syls[j].Org != 0; j++ {
			if libc.StrNCmp(syls[j].Org, &lbuf[ofs], libc.StrLen(syls[j].Org)) == 0 {
				libc.StrCat(&buf[0], syls[j].News)
				ofs += libc.StrLen(syls[j].Org)
				break
			}
		}
		if *syls[j].Org == 0 {
			basic_mud_log(libc.CString("No entry in syllable table for substring of '%s'"), &lbuf[0])
			ofs++
		}
	}
	if tch != nil && tch.In_room == ch.In_room {
		if tch == ch {
			format = libc.CString("$n closes $s eyes and utters the words, '%s'.")
		} else {
			format = libc.CString("$n stares at $N and utters the words, '%s'.")
		}
	} else if tobj != nil && (tobj.In_room == ch.In_room || tobj.Carried_by == ch) {
		format = libc.CString("$n stares at $p and utters the words, '%s'.")
	} else {
		format = libc.CString("$n utters the words, '%s'.")
	}
	stdio.Snprintf(&buf1[0], int(256), libc.GoString(format), skill_name(spellnum))
	stdio.Snprintf(&buf2[0], int(256), libc.GoString(format), &buf[0])
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; i != nil; i = i.Next_in_room {
		if i == ch || i == tch || i.Desc == nil || !AWAKE(i) {
			continue
		}
		if ((i.Chclasses[ch.Chclass])+(i.Epicclasses[ch.Chclass])) != 0 || i.Admlevel >= ADMLVL_IMMORT {
			perform_act(&buf1[0], ch, tobj, unsafe.Pointer(tch), i)
		} else {
			perform_act(&buf2[0], ch, tobj, unsafe.Pointer(tch), i)
		}
	}
	if tch != nil && tch != ch && tch.In_room == ch.In_room {
		stdio.Snprintf(&buf1[0], int(256), "$n stares at you and utters the words, '%s'.", func() *byte {
			if ((tch.Chclasses[ch.Chclass]) + (tch.Epicclasses[ch.Chclass])) != 0 {
				return skill_name(spellnum)
			}
			return &buf[0]
		}())
		act(&buf1[0], FALSE, ch, nil, unsafe.Pointer(tch), TO_VICT)
	}
}
func skill_name(num int) *byte {
	if num > 0 && num < SKILL_TABLE_SIZE {
		return spell_info[num].Name
	} else if num == -1 {
		return libc.CString("UNUSED")
	} else {
		return libc.CString("UNDEFINED")
	}
}
func find_skill_num(name *byte, sktype int) int {
	var (
		skindex int
		ok      int
		temp    *byte
		temp2   *byte
		first   [256]byte
		first2  [256]byte
		tempbuf [256]byte
	)
	for skindex = 1; skindex < SKILL_TABLE_SIZE; skindex++ {
		if is_abbrev(name, spell_info[skindex].Name) != 0 && (spell_info[skindex].Skilltype&sktype) != 0 {
			return skindex
		}
		ok = TRUE
		strlcpy(&tempbuf[0], spell_info[skindex].Name, uint64(256))
		temp = any_one_arg(&tempbuf[0], &first[0])
		temp2 = any_one_arg(name, &first2[0])
		for first[0] != 0 && first2[0] != 0 && ok != 0 {
			if is_abbrev(&first2[0], &first[0]) == 0 {
				ok = FALSE
			}
			temp = any_one_arg(temp, &first[0])
			temp2 = any_one_arg(temp2, &first2[0])
		}
		if ok != 0 && first2[0] == 0 && (spell_info[skindex].Skilltype&sktype) != 0 {
			return skindex
		}
	}
	return -1
}
func call_magic(caster *char_data, cvict *char_data, ovict *obj_data, spellnum int, level int, casttype int, arg *byte) int {
	if spellnum < 1 || spellnum > SKILL_TABLE_SIZE {
		return 0
	}
	if cast_wtrigger(caster, cvict, ovict, spellnum) == 0 {
		return 0
	}
	if cast_otrigger(caster, ovict, spellnum) == 0 {
		return 0
	}
	if cast_mtrigger(caster, cvict, spellnum) == 0 {
		return 0
	}
	if ROOM_FLAGGED(caster.In_room, ROOM_PEACEFUL) && caster.Admlevel < ADMLVL_IMPL && (int(spell_info[spellnum].Violent) != 0 || (spell_info[spellnum].Routines&(1<<0)) != 0) {
		send_to_char(caster, libc.CString("A flash of white light fills the room, dispelling your violent magic!\r\n"))
		act(libc.CString("White light from no particular source suddenly fills the room, then vanishes."), FALSE, caster, nil, nil, TO_ROOM)
		return 0
	}
	if (spell_info[spellnum].Routines&(1<<15)) != 0 && casttype != CAST_STRIKE {
		mag_nextstrike(level, caster, spellnum)
		return 1
	}
	if (spell_info[spellnum].Routines & (1 << 0)) != 0 {
		if mag_damage(level, caster, cvict, spellnum) == -1 {
			return -1
		}
	}
	if (spell_info[spellnum].Routines & (1 << 1)) != 0 {
		mag_affects(level, caster, cvict, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 2)) != 0 {
		mag_unaffects(level, caster, cvict, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 3)) != 0 {
		mag_points(level, caster, cvict, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 4)) != 0 {
		mag_alter_objs(level, caster, ovict, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 5)) != 0 {
		mag_groups(level, caster, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 6)) != 0 {
		mag_masses(level, caster, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 7)) != 0 {
		mag_areas(level, caster, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 8)) != 0 {
		mag_summons(level, caster, ovict, spellnum, arg)
	}
	if (spell_info[spellnum].Routines & (1 << 9)) != 0 {
		mag_creations(level, caster, spellnum)
	}
	if (spell_info[spellnum].Routines & (1 << 10)) != 0 {
		switch spellnum {
		case SPELL_CHARM:
			spell_charm(level, caster, cvict, ovict, arg)
		case SPELL_CREATE_WATER:
			spell_create_water(level, caster, cvict, ovict, arg)
		case SPELL_DETECT_POISON:
			spell_detect_poison(level, caster, cvict, ovict, arg)
		case SPELL_ENCHANT_WEAPON:
			spell_enchant_weapon(level, caster, cvict, ovict, arg)
		case SPELL_IDENTIFY:
			spell_identify(level, caster, cvict, ovict, arg)
		case SPELL_LOCATE_OBJECT:
			spell_locate_object(level, caster, cvict, ovict, arg)
		case SPELL_SUMMON:
			spell_summon(level, caster, cvict, ovict, arg)
		case SPELL_WORD_OF_RECALL:
			spell_recall(level, caster, cvict, ovict, arg)
		case SPELL_TELEPORT:
			spell_teleport(level, caster, cvict, ovict, arg)
		case SPELL_PORTAL:
			spell_portal(level, caster, cvict, ovict, arg)
		case ART_ABUNDANT_STEP:
			art_abundant_step(level, caster, cvict, ovict, arg)
		}
	}
	if (spell_info[spellnum].Routines & (1 << 11)) != 0 {
		mag_affectsv(level, caster, cvict, spellnum)
	}
	return 1
}
func mag_objectmagic(ch *char_data, obj *obj_data, argument *byte) {
	var (
		arg      [2048]byte
		i        int
		k        int
		tch      *char_data = nil
		next_tch *char_data
		tobj     *obj_data = nil
	)
	one_argument(argument, &arg[0])
	k = generic_find(&arg[0], (1<<0)|1<<2|1<<3|1<<5, ch, &tch, &tobj)
	switch obj.Type_flag {
	case ITEM_STAFF:
		act(libc.CString("You tap $p three times on the ground."), FALSE, ch, obj, nil, TO_CHAR)
		if obj.Action_description != nil {
			act(obj.Action_description, FALSE, ch, obj, nil, TO_ROOM)
		} else {
			act(libc.CString("$n taps $p three times on the ground."), FALSE, ch, obj, nil, TO_ROOM)
		}
		if (obj.Value[VAL_STAFF_CHARGES]) <= 0 {
			send_to_char(ch, libc.CString("It seems powerless.\r\n"))
			act(libc.CString("Nothing seems to happen."), FALSE, ch, obj, nil, TO_ROOM)
		} else {
			(obj.Value[VAL_STAFF_CHARGES])--
			ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
			if (obj.Value[VAL_STAFF_LEVEL]) != 0 {
				k = obj.Value[VAL_STAFF_LEVEL]
			} else {
				k = DEFAULT_STAFF_LVL
			}
			if (spell_info[obj.Value[VAL_STAFF_SPELL]].Routines & ((1 << 6) | 1<<7)) != 0 {
				for func() *char_data {
					i = 0
					return func() *char_data {
						tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People
						return tch
					}()
				}(); tch != nil; tch = tch.Next_in_room {
					i++
				}
				for func() int {
					p := &i
					x := *p
					*p--
					return x
				}() > 0 {
					call_magic(ch, nil, nil, obj.Value[VAL_STAFF_SPELL], k, CAST_STAFF, nil)
				}
			} else {
				for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = next_tch {
					next_tch = tch.Next_in_room
					if ch != tch {
						call_magic(ch, tch, nil, obj.Value[VAL_STAFF_SPELL], k, CAST_STAFF, nil)
					}
				}
			}
		}
	case ITEM_WAND:
		if k == (1 << 0) {
			if tch == ch {
				act(libc.CString("You point $p at yourself."), FALSE, ch, obj, nil, TO_CHAR)
				act(libc.CString("$n points $p at $mself."), FALSE, ch, obj, nil, TO_ROOM)
			} else {
				act(libc.CString("You point $p at $N."), FALSE, ch, obj, unsafe.Pointer(tch), TO_CHAR)
				if obj.Action_description != nil {
					act(obj.Action_description, FALSE, ch, obj, unsafe.Pointer(tch), TO_ROOM)
				} else {
					act(libc.CString("$n points $p at $N."), TRUE, ch, obj, unsafe.Pointer(tch), TO_ROOM)
				}
			}
		} else if tobj != nil {
			act(libc.CString("You point $p at $P."), FALSE, ch, obj, unsafe.Pointer(tobj), TO_CHAR)
			if obj.Action_description != nil {
				act(obj.Action_description, FALSE, ch, obj, unsafe.Pointer(tobj), TO_ROOM)
			} else {
				act(libc.CString("$n points $p at $P."), TRUE, ch, obj, unsafe.Pointer(tobj), TO_ROOM)
			}
		} else if (spell_info[obj.Value[VAL_WAND_SPELL]].Routines & ((1 << 7) | 1<<6)) != 0 {
			act(libc.CString("You point $p outward."), FALSE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n points $p outward."), TRUE, ch, obj, nil, TO_ROOM)
		} else {
			act(libc.CString("At what should $p be pointed?"), FALSE, ch, obj, nil, TO_CHAR)
			return
		}
		if (obj.Value[VAL_WAND_CHARGES]) <= 0 {
			send_to_char(ch, libc.CString("It seems powerless.\r\n"))
			act(libc.CString("Nothing seems to happen."), FALSE, ch, obj, nil, TO_ROOM)
			return
		}
		(obj.Value[VAL_WAND_CHARGES])--
		ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
		if (obj.Value[VAL_WAND_LEVEL]) != 0 {
			call_magic(ch, tch, tobj, obj.Value[VAL_WAND_SPELL], obj.Value[VAL_WAND_LEVEL], CAST_WAND, nil)
		} else {
			call_magic(ch, tch, tobj, obj.Value[VAL_WAND_SPELL], DEFAULT_WAND_LVL, CAST_WAND, nil)
		}
	case ITEM_SCROLL:
		if arg[0] != 0 {
			if k == 0 {
				act(libc.CString("There is nothing to here to affect with $p."), FALSE, ch, obj, nil, TO_CHAR)
				return
			}
		} else {
			tch = ch
		}
		act(libc.CString("You recite $p which dissolves."), TRUE, ch, obj, nil, TO_CHAR)
		if obj.Action_description != nil {
			act(obj.Action_description, FALSE, ch, obj, unsafe.Pointer(tch), TO_ROOM)
		} else {
			act(libc.CString("$n recites $p."), FALSE, ch, obj, nil, TO_ROOM)
		}
		ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
		for i = 1; i <= 3; i++ {
			if call_magic(ch, tch, tobj, obj.Value[i], obj.Value[VAL_SCROLL_LEVEL], CAST_SCROLL, nil) <= 0 {
				break
			}
		}
		if obj != nil {
			extract_obj(obj)
		}
	case ITEM_POTION:
		tch = ch
		if consume_otrigger(obj, ch, OCMD_QUAFF) == 0 {
			return
		}
		act(libc.CString("You swallow $p."), FALSE, ch, obj, nil, TO_CHAR)
		if obj.Action_description != nil {
			act(obj.Action_description, FALSE, ch, obj, nil, TO_ROOM)
		} else {
			act(libc.CString("$n swallows $p."), TRUE, ch, obj, nil, TO_ROOM)
		}
		ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
		for i = 1; i <= 3; i++ {
			if call_magic(ch, ch, nil, obj.Value[i], obj.Value[VAL_POTION_LEVEL], CAST_POTION, nil) <= 0 {
				break
			}
		}
		if obj != nil {
			extract_obj(obj)
		}
	default:
		basic_mud_log(libc.CString("SYSERR: Unknown object_type %d in mag_objectmagic."), obj.Type_flag)
	}
}
func cast_spell(ch *char_data, tch *char_data, tobj *obj_data, spellnum int, arg *byte) int {
	var whichclass int = -1
	_ = whichclass
	var i int
	var j int
	var diff int = -1
	var lvl int = 1
	if spellnum < 0 || spellnum >= SKILL_TABLE_SIZE {
		basic_mud_log(libc.CString("SYSERR: cast_spell trying to call out of range spellnum %d/%d."), spellnum, SKILL_TABLE_SIZE)
		return 0
	}
	if (spell_info[spellnum].Skilltype & ((1 << 0) | 1<<4)) == 0 {
		basic_mud_log(libc.CString("SYSERR: cast_spell trying to call nonspell spellnum %d/%d."), spellnum, SKILL_TABLE_SIZE)
		return 0
	}
	if int(ch.Position) < int(spell_info[spellnum].Min_position) {
		switch ch.Position {
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You dream about great magical powers.\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("You cannot concentrate while resting.\r\n"))
		case POS_SITTING:
			send_to_char(ch, libc.CString("You can't do this sitting!\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Impossible!  You can't concentrate enough!\r\n"))
		default:
			send_to_char(ch, libc.CString("You can't do much of anything like this!\r\n"))
		}
		return 0
	}
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master == tch {
		send_to_char(ch, libc.CString("You are afraid you might hurt your master!\r\n"))
		return 0
	}
	if tch != ch && (spell_info[spellnum].Targets&(1<<5)) != 0 {
		send_to_char(ch, libc.CString("You can only cast this spell upon yourself!\r\n"))
		return 0
	}
	if tch == ch && (spell_info[spellnum].Targets&(1<<6)) != 0 {
		send_to_char(ch, libc.CString("You cannot cast this spell upon yourself!\r\n"))
		return 0
	}
	if (spell_info[spellnum].Routines&(1<<5)) != 0 && !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("You can't cast this spell if you're not in a group!\r\n"))
		return 0
	}
	if (spell_info[spellnum].Skilltype & (1 << 0)) != 0 {
		for i = 0; i < NUM_CLASSES; i++ {
			j = ((ch.Chclasses[i]) + (ch.Epicclasses[i])) - spell_info[spellnum].Min_level[i]
			if j > diff {
				whichclass = i
				diff = j
				lvl = (ch.Chclasses[i]) + (ch.Epicclasses[i])
			}
		}
	} else if (spell_info[spellnum].Skilltype & (1 << 4)) != 0 {
		lvl = ((ch.Chclasses[CLASS_KABITO]) + (ch.Epicclasses[CLASS_KABITO])) / 2
	} else {
		lvl = GET_LEVEL(ch)
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	if (spell_info[spellnum].Skilltype & (1 << 0)) != 0 {
		say_spell(ch, spellnum, tch, tobj)
	}
	return call_magic(ch, tch, tobj, spellnum, lvl, CAST_SPELL, arg)
}
func do_cast(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		tch      *char_data = nil
		tobj     *obj_data  = nil
		s        *byte
		t        *byte
		buffer   [25]byte
		ki       int = 0
		spellnum int
		i        int
		target   int = 0
		innate   int = FALSE
	)
	_ = innate
	s = libc.StrTok(argument, libc.CString("'"))
	if s == nil {
		if subcmd == SCMD_ART {
			send_to_char(ch, libc.CString("Use what ability?\r\n"))
		} else {
			send_to_char(ch, libc.CString("Cast what where?\r\n"))
		}
		return
	}
	s = libc.StrTok(nil, libc.CString("'"))
	if s == nil {
		if subcmd == SCMD_ART {
			send_to_char(ch, libc.CString("You must enclose the ability name in quotes: '\r\n"))
		} else {
			send_to_char(ch, libc.CString("Spell names must be enclosed in the Holy Magic Symbols: '\r\n"))
		}
		return
	}
	t = libc.StrTok(nil, libc.CString("\x00"))
	spellnum = find_skill_num(s, (1<<0)|1<<4)
	stdio.Sprintf(&buffer[0], "%d", spellnum)
	if subcmd == SCMD_ART {
		if int(ch.Chclass) != CLASS_KABITO && ch.Admlevel < ADMLVL_IMMORT {
			send_to_char(ch, libc.CString("You are not trained in that.\r\n"))
			return
		}
		if spellnum < 1 || spellnum >= SKILL_TABLE_SIZE || (spell_info[spellnum].Skilltype&(1<<4)) == 0 || *s == 0 {
			send_to_char(ch, libc.CString("I don't recognize that martial art or ability.\r\n"))
			return
		}
		if GET_SKILL(ch, spellnum) == 0 {
			send_to_char(ch, libc.CString("You do not have that ability.\r\n"))
			return
		}
	} else {
		if ch.Admlevel < ADMLVL_IMMORT {
			send_to_char(ch, libc.CString("You are not able to cast spells.\r\n"))
			return
		}
		if spellnum < 1 || spellnum >= SKILL_TABLE_SIZE || (spell_info[spellnum].Skilltype&(1<<0)) == 0 {
			send_to_char(ch, libc.CString("Cast what?!?\r\n"))
			return
		}
	}
	if subcmd == SCMD_ART {
		ki = mag_kicost(ch, spellnum)
		if ki > 0 && ch.Ki < int64(ki) && ch.Admlevel < ADMLVL_IMMORT {
			send_to_char(ch, libc.CString("You haven't the energy to cast that spell!\r\n"))
			return
		}
	}
	if t != nil {
		var arg [2048]byte
		strlcpy(&arg[0], t, uint64(2048))
		one_argument(&arg[0], t)
		skip_spaces(&t)
	}
	if (spell_info[spellnum].Targets & (1 << 0)) != 0 {
		target = TRUE
	} else if t != nil && *t != 0 {
		if target == 0 && (spell_info[spellnum].Targets&(1<<1)) != 0 {
			if (func() *char_data {
				tch = get_char_vis(ch, t, nil, 1<<0)
				return tch
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<2)) != 0 {
			if (func() *char_data {
				tch = get_char_vis(ch, t, nil, 1<<1)
				return tch
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<7)) != 0 {
			if (func() *obj_data {
				tobj = get_obj_in_list_vis(ch, t, nil, ch.Carrying)
				return tobj
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<10)) != 0 {
			for i = 0; target == 0 && i < NUM_WEARS; i++ {
				if (ch.Equipment[i]) != nil && isname(t, (ch.Equipment[i]).Name) != 0 {
					tobj = ch.Equipment[i]
					target = TRUE
				}
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<8)) != 0 {
			if (func() *obj_data {
				tobj = get_obj_in_list_vis(ch, t, nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
				return tobj
			}()) != nil {
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<9)) != 0 {
			if (func() *obj_data {
				tobj = get_obj_vis(ch, t, nil)
				return tobj
			}()) != nil {
				target = TRUE
			}
		}
	} else {
		if target == 0 && (spell_info[spellnum].Targets&(1<<3)) != 0 {
			if ch.Fighting != nil {
				tch = ch
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<4)) != 0 {
			if ch.Fighting != nil {
				tch = ch.Fighting
				target = TRUE
			}
		}
		if target == 0 && (spell_info[spellnum].Targets&(1<<1)) != 0 && int(spell_info[spellnum].Violent) == 0 {
			tch = ch
			target = TRUE
		}
		if target == 0 {
			send_to_char(ch, libc.CString("Upon %s should the spell be cast?\r\n"), func() string {
				if (spell_info[spellnum].Targets & ((1 << 8) | 1<<7 | 1<<9 | 1<<10)) != 0 {
					return "what"
				}
				return "who"
			}())
			return
		}
	}
	if target != 0 && tch == ch && int(spell_info[spellnum].Violent) != 0 {
		send_to_char(ch, libc.CString("You shouldn't cast that on yourself -- could be bad for your health!\r\n"))
		return
	}
	if target == 0 {
		send_to_char(ch, libc.CString("Cannot find the target of your spell!\r\n"))
		return
	}
	if is_innate(ch, spellnum) != 0 {
		innate = TRUE
	}
	if int(spell_info[spellnum].Violent) != 0 && tch != nil && IS_NPC(tch) {
		if tch.Fighting == nil {
			set_fighting(tch, ch)
		}
		if ch.Fighting == nil {
			set_fighting(ch, tch)
		}
	}
	if (spell_info[spellnum].Comp_flags&(1<<4)) != 0 && rand_number(1, 100) <= int(ch.Spellfail) {
		if (spell_info[spellnum].Routines & ((1 << 14) | 1<<13)) != 0 {
			ch.Affected_by[int(AFF_NEXTPARTIAL/32)] |= 1 << (int(AFF_NEXTPARTIAL % 32))
		} else if (spell_info[spellnum].Routines & ((1 << 14) | 1<<14)) != 0 {
			ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
		}
		send_to_char(ch, libc.CString("Your armor interferes with your casting, and you fail!\r\n"))
	} else {
		if ki > 0 {
			ch.Ki = int64(MAX(0, MIN(int(ch.Max_ki), int(ch.Ki-int64(ki)))))
		}
		if cast_spell(ch, tch, tobj, spellnum, t) != 0 && ch.Admlevel < ADMLVL_IMMORT {
			if (spell_info[spellnum].Routines & ((1 << 14) | 1<<13)) != 0 {
				ch.Affected_by[int(AFF_NEXTPARTIAL/32)] |= 1 << (int(AFF_NEXTPARTIAL % 32))
			} else if (spell_info[spellnum].Routines & ((1 << 14) | 1<<14)) != 0 {
				ch.Affected_by[int(AFF_NEXTNOACTION/32)] |= 1 << (int(AFF_NEXTNOACTION % 32))
			}
			if subcmd == SCMD_CAST {
				send_to_char(ch, libc.CString("The magical energy from the spell leaves your mind.\r\n"))
				if PRF_FLAGGED(ch, PRF_AUTOMEM) {
					send_to_char(ch, libc.CString("You begin to commit the spell again to your mind.\r\n"))
				}
			}
		}
	}
}
func skill_race_class(spell int, race int, learntype int) {
	var bad int = 0
	if spell < 0 || spell >= SKILL_TABLE_SIZE {
		basic_mud_log(libc.CString("SYSERR: attempting assign to illegal spellnum %d/%d"), spell, SKILL_TABLE_SIZE)
		return
	}
	if race < 0 || race >= NUM_RACES {
		basic_mud_log(libc.CString("SYSERR: assigning '%s' to illegal race %d/%d."), skill_name(spell), race, int(NUM_RACES-1))
		bad = 1
	}
	if bad == 0 {
		spell_info[spell].Race_can_learn[race] = learntype
	}
}
func spell_level(spell int, chclass int, level int) {
	var bad int = 0
	if spell < 0 || spell > SKILL_TABLE_SIZE {
		basic_mud_log(libc.CString("SYSERR: attempting assign to illegal spellnum %d/%d"), spell, SKILL_TABLE_SIZE)
		return
	}
	if chclass < 0 || chclass >= NUM_CLASSES {
		basic_mud_log(libc.CString("SYSERR: assigning '%s' to illegal class %d/%d."), skill_name(spell), chclass, int(NUM_CLASSES-1))
		bad = 1
	}
	if level < 1 {
		basic_mud_log(libc.CString("SYSERR: assigning '%s' to illegal level %d."), skill_name(spell), level)
		bad = 1
	}
	if bad == 0 {
		spell_info[spell].Min_level[chclass] = level
	}
}
func skill_class(skill int, chclass int, learntype int) {
	var bad int = 0
	if skill < 0 || skill > SKILL_TABLE_SIZE {
		basic_mud_log(libc.CString("SYSERR: attempting assign to illegal skillnum %d/%d"), skill, SKILL_TABLE_SIZE)
		return
	}
	if chclass < 0 || chclass >= NUM_CLASSES {
		basic_mud_log(libc.CString("SYSERR: assigning '%s' to illegal class %d/%d."), skill_name(skill), chclass, int(NUM_CLASSES-1))
		bad = 1
	}
	if learntype < 0 || learntype > SKLEARN_CLASS {
		basic_mud_log(libc.CString("SYSERR: assigning skill '%s' illegal learn type %d for class %d."), skill_name(skill), learntype, chclass)
		bad = 1
	}
	if bad == 0 {
		spell_info[skill].Can_learn_skill[chclass] = int8(learntype)
	}
}
func skill_type(snum int) int {
	return spell_info[snum].Skilltype
}
func set_skill_type(snum int, sktype int) {
	spell_info[snum].Skilltype = sktype
}
func spello(spl int, name *byte, max_mana int, min_mana int, mana_change int, minpos int, targets int, violent int, routines int, save_flags int, comp_flags int, wearoff *byte, cmspell_level int, school int, domain int) {
	var i int
	for i = 0; i < NUM_CLASSES; i++ {
		spell_info[spl].Min_level[i] = config_info.Play.Level_cap
	}
	for i = 0; i < NUM_RACES; i++ {
		spell_info[spl].Race_can_learn[i] = config_info.Play.Level_cap
	}
	spell_info[spl].Mana_max = max_mana
	spell_info[spl].Mana_min = min_mana
	spell_info[spl].Mana_change = mana_change
	spell_info[spl].Ki_max = 0
	spell_info[spl].Ki_min = 0
	spell_info[spl].Ki_change = 0
	spell_info[spl].Min_position = int8(minpos)
	spell_info[spl].Targets = targets
	spell_info[spl].Violent = int8(violent)
	spell_info[spl].Routines = routines
	spell_info[spl].Name = name
	spell_info[spl].Wear_off_msg = wearoff
	spell_info[spl].Skilltype = 1 << 0
	spell_info[spl].Flags = 0
	spell_info[spl].Save_flags = save_flags
	spell_info[spl].Comp_flags = comp_flags
	spell_info[spl].Spell_level = cmspell_level
	spell_info[spl].School = school
	spell_info[spl].Domain = domain
}
func arto(spl int, name *byte, max_ki int, min_ki int, ki_change int, minpos int, targets int, violent int, routines int, save_flags int, comp_flags int, wearoff *byte) {
	spello(spl, name, 0, 0, 0, minpos, targets, violent, routines, save_flags, comp_flags, wearoff, 0, 0, 0)
	set_skill_type(spl, 1<<4)
	spell_info[spl].Ki_max = max_ki
	spell_info[spl].Ki_min = min_ki
	spell_info[spl].Ki_change = ki_change
}
func unused_spell(spl int) {
	var i int
	for i = 0; i < NUM_CLASSES; i++ {
		spell_info[spl].Min_level[i] = config_info.Play.Level_cap
		spell_info[spl].Can_learn_skill[i] = SKLEARN_CROSSCLASS
	}
	for i = 0; i < NUM_RACES; i++ {
		spell_info[spl].Race_can_learn[i] = SKLEARN_CROSSCLASS
	}
	spell_info[spl].Mana_max = 0
	spell_info[spl].Mana_min = 0
	spell_info[spl].Mana_change = 0
	spell_info[spl].Ki_max = 0
	spell_info[spl].Ki_min = 0
	spell_info[spl].Ki_change = 0
	spell_info[spl].Min_position = 0
	spell_info[spl].Targets = 0
	spell_info[spl].Violent = 0
	spell_info[spl].Routines = 0
	spell_info[spl].Name = unused_spellname
	spell_info[spl].Skilltype = SKTYPE_NONE
	spell_info[spl].Flags = 0
	spell_info[spl].Save_flags = 0
	spell_info[spl].Comp_flags = 0
	spell_info[spl].Spell_level = 0
	spell_info[spl].School = 0
	spell_info[spl].Domain = 0
}
func skillo(skill int, name *byte, flags int) {
	spello(skill, name, 0, 0, 0, 0, 0, 0, 0, 0, 0, nil, 0, 0, 0)
	spell_info[skill].Skilltype = 1 << 1
	spell_info[skill].Flags = flags
}
func mag_assign_spells() {
	var i int
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		unused_spell(i)
	}
	spello(SPELL_ANIMATE_DEAD, libc.CString("animate dead"), 35, 10, 3, POS_STANDING, 1<<8, FALSE, (1<<14)|1<<8, 0, 0, nil, 3, SCHOOL_NECROMANCY, DOMAIN_DEATH)
	spello(SPELL_MAGE_ARMOR, libc.CString("mage armor"), 30, 15, 3, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("You feel less protected."), 1, SCHOOL_CONJURATION, -1)
	spello(SPELL_BLESS, libc.CString("blessed"), 35, 5, 3, POS_STANDING, (1<<1)|1<<7, FALSE, (1<<14)|1<<1|1<<4, 0, 0, libc.CString("You feel less righteous."), 1, -1, DOMAIN_UNIVERSAL)
	spello(SKILL_SPIRITCONTROL, libc.CString("spirit control"), 35, 5, 3, POS_STANDING, (1<<1)|1<<7, FALSE, (1<<14)|1<<1|1<<4, 0, 0, libc.CString("You no longer have full control of your spirit."), 1, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_BLINDNESS, libc.CString("blindness"), 35, 25, 1, POS_STANDING, (1<<1)|1<<6, FALSE, (1<<14)|1<<1, (1<<0)|1<<4, 0, libc.CString("You feel a cloak of blindness dissolve."), 2, SCHOOL_TRANSMUTATION, DOMAIN_UNIVERSAL)
	spello(SKILL_SOLARF, libc.CString("blind"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("You are no longer blind!"), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SPELL_BURNING_HANDS, libc.CString("burning hands"), 30, 10, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<1)|1<<3, 0, nil, 1, SCHOOL_TRANSMUTATION, DOMAIN_FIRE)
	spello(SPELL_CALL_LIGHTNING, libc.CString("call lightning"), 40, 25, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<1)|1<<3, 0, nil, 3, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_INFLICT_CRITIC, libc.CString("inflict critic"), 30, 10, 2, POS_FIGHTING, 1<<1, TRUE, (1<<14)|1<<0, (1<<2)|1<<3, 0, nil, 4, -1, DOMAIN_HEALING)
	spello(SPELL_INFLICT_LIGHT, libc.CString("inflict light"), 30, 10, 2, POS_FIGHTING, 1<<1, TRUE, (1<<14)|1<<0, (1<<2)|1<<3, 0, nil, 1, -1, DOMAIN_HEALING)
	spello(SPELL_CHARM, libc.CString("charm person"), 75, 50, 2, POS_FIGHTING, (1<<1)|1<<6, TRUE, (1<<14)|1<<10, (1<<2)|1<<4, 0, libc.CString("You feel more self-confident."), 1, SCHOOL_ENCHANTMENT, -1)
	spello(SPELL_CHILL_TOUCH, libc.CString("chill touch"), 30, 10, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0|1<<1, (1<<0)|1<<5, 0, libc.CString("You feel your strength return."), 1, SCHOOL_NECROMANCY, -1)
	spello(SPELL_COLOR_SPRAY, libc.CString("color spray"), 30, 15, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<2)|1<<4, 0, nil, 1, SCHOOL_ILLUSION, -1)
	spello(SPELL_CONTROL_WEATHER, libc.CString("control weather"), 75, 25, 5, POS_STANDING, 1<<0, FALSE, (1<<14)|1<<10, 0, 0, nil, 7, SCHOOL_TRANSMUTATION, DOMAIN_AIR)
	spello(SPELL_CREATE_FOOD, libc.CString("create food"), 30, 5, 4, POS_STANDING, 1<<0, FALSE, (1<<14)|1<<9, 0, 0, nil, 3, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_CREATE_WATER, libc.CString("create water"), 30, 5, 4, POS_STANDING, (1<<7)|1<<10, FALSE, (1<<14)|1<<10, 0, 0, nil, 0, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_REMOVE_BLINDNESS, libc.CString("remove blindness"), 30, 5, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<2, 0, 0, nil, 3, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_CURE_CRITIC, libc.CString("cure critic"), 30, 10, 2, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<3, 0, 0, nil, 4, -1, DOMAIN_HEALING)
	spello(SPELL_CURE_LIGHT, libc.CString("cure light"), 30, 10, 2, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<3, 0, 0, nil, 1, -1, DOMAIN_HEALING)
	spello(SPELL_BESTOW_CURSE, libc.CString("bestow curse"), 80, 50, 2, POS_FIGHTING, 1<<1, TRUE, 1<<1, (1<<2)|1<<4, 0, libc.CString("You feel more optimistic."), 8, SCHOOL_NECROMANCY, DOMAIN_DESTRUCTION)
	spello(SPELL_BANE, libc.CString("bane"), 80, 50, 2, POS_FIGHTING, 1<<1, TRUE, 1<<1, (1<<2)|1<<4, 0, libc.CString("You feel more optimistic."), 8, SCHOOL_ENCHANTMENT, DOMAIN_CHARM)
	spello(SPELL_DETECT_ALIGN, libc.CString("detect alignment"), 20, 10, 2, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("You feel less aware."), 1, SCHOOL_DIVINATION, DOMAIN_UNIVERSAL)
	spello(SPELL_SEE_INVIS, libc.CString("see invisibility"), 20, 10, 2, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("Your eyes stop tingling."), 2, SCHOOL_DIVINATION, -1)
	spello(SPELL_DETECT_MAGIC, libc.CString("detect magic"), 20, 10, 2, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("The detect magic wears off."), 0, SCHOOL_UNIVERSAL, DOMAIN_UNIVERSAL)
	spello(SPELL_DETECT_POISON, libc.CString("detect poison"), 15, 5, 1, POS_STANDING, (1<<1)|1<<7|1<<8, FALSE, (1<<14)|1<<10, 0, 0, libc.CString("The detect poison wears off."), 0, SCHOOL_DIVINATION, DOMAIN_UNIVERSAL)
	spello(SPELL_DISPEL_EVIL, libc.CString("dispel evil"), 40, 25, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<2)|1<<4, 0, nil, 5, -1, DOMAIN_GOOD)
	spello(SPELL_DISPEL_GOOD, libc.CString("dispel good"), 40, 25, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<2)|1<<4, 0, nil, 5, -1, DOMAIN_EVIL)
	spello(SPELL_EARTHQUAKE, libc.CString("earthquake"), 40, 25, 3, POS_FIGHTING, 1<<0, TRUE, (1<<14)|1<<7, (1<<1)|1<<3, 0, nil, 8, -1, int(DOMAIN_DESTRUCTION|DOMAIN_EARTH))
	spello(SPELL_ENCHANT_WEAPON, libc.CString("enchant weapon"), 150, 100, 10, POS_STANDING, 1<<7, FALSE, (1<<14)|1<<10, 0, 0, nil, 9, SCHOOL_TRANSMUTATION, -1)
	spello(SPELL_ENERGY_DRAIN, libc.CString("energy drain"), 40, 25, 1, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<0)|1<<10, (1<<0)|1<<4, 0, nil, 9, SCHOOL_NECROMANCY, DOMAIN_UNIVERSAL)
	spello(SPELL_GROUP_ARMOR, libc.CString("group armor"), 50, 30, 2, POS_STANDING, 1<<0, FALSE, (1<<14)|1<<5, 0, 0, nil, 5, SCHOOL_CONJURATION, -1)
	spello(SPELL_FAERIE_FIRE, libc.CString("faerie fire"), 20, 10, 2, POS_STANDING, (1<<2)|1<<6, FALSE, (1<<14)|1<<10, 0, 0, nil, 1, SCHOOL_EVOCATION, DOMAIN_UNIVERSAL)
	spello(SPELL_FIREBALL, libc.CString("fireball"), 40, 30, 2, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<1)|1<<3, 0, nil, 3, SCHOOL_EVOCATION, -1)
	spello(SPELL_MASS_HEAL, libc.CString("mass heal"), 80, 60, 5, POS_STANDING, 1<<0, FALSE, (1<<14)|1<<5, 0, 0, nil, 6, SCHOOL_CONJURATION, DOMAIN_HEALING)
	spello(SPELL_HARM, libc.CString("harm"), 75, 45, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<0)|1<<4, 0, nil, 6, -1, DOMAIN_DESTRUCTION)
	spello(SPELL_HEAL, libc.CString("heal"), 60, 40, 3, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<3|1<<2, 0, 0, nil, 6, -1, DOMAIN_HEALING)
	spello(SPELL_SENSU, libc.CString("sensu"), 1, 0, 3, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<3, 0, 0, nil, 6, SCHOOL_DIVINATION, DOMAIN_HEALING)
	spello(SPELL_IDENTIFY, libc.CString("identify"), 50, 25, 5, POS_STANDING, (1<<7)|1<<8, FALSE, (1<<14)|1<<10, 0, 0, nil, 2, SCHOOL_DIVINATION, DOMAIN_MAGIC)
	spello(SPELL_DARKVISION, libc.CString("darkvision"), 25, 10, 1, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("Your night vision seems to fade."), 2, -1, -1)
	spello(SPELL_INVISIBLE, libc.CString("invisibility"), 35, 25, 1, POS_STANDING, (1<<1)|1<<7|1<<8, FALSE, (1<<14)|1<<1|1<<4, 0, 0, libc.CString("You feel yourself exposed."), 2, SCHOOL_ILLUSION, DOMAIN_TRICKERY)
	spello(SPELL_LIGHTNING_BOLT, libc.CString("lightning bolt"), 30, 15, 1, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, (1<<1)|1<<3, 0, nil, 3, SCHOOL_EVOCATION, -1)
	spello(SPELL_LOCATE_OBJECT, libc.CString("locate object"), 25, 20, 1, POS_STANDING, 1<<9, FALSE, (1<<14)|1<<10, 0, 0, nil, 3, SCHOOL_DIVINATION, DOMAIN_TRAVEL)
	spello(SPELL_MAGIC_MISSILE, libc.CString("magic missile"), 25, 10, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, 0, 0, nil, 1, SCHOOL_EVOCATION, -1)
	spello(SPELL_PARALYZE, libc.CString("stone"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("Your body is no longer petrified."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_HASSHUKEN, libc.CString("hasshuken"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("Your arms slow down."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_CURSE, libc.CString("curse"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("You are no longer cursed!"), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_MIGHT, libc.CString("might"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("Your strength fades."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_PARALYZE, libc.CString("paralyze"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("Your feel like you are able to move again."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_POISON, libc.CString("poison"), 50, 20, 3, POS_STANDING, (1<<1)|1<<6|1<<7, TRUE, (1<<14)|1<<1|1<<4, (1<<0)|1<<4, 0, libc.CString("You feel like you got over something."), 4, SCHOOL_NECROMANCY, DOMAIN_UNIVERSAL)
	spello(SKILL_POISON, libc.CString("dark metamorphosis"), 50, 20, 3, POS_STANDING, (1<<1)|1<<6|1<<7, TRUE, (1<<14)|1<<1|1<<4, (1<<0)|1<<4, 0, libc.CString("Your dark metamorphosis fades."), 4, SCHOOL_NECROMANCY, DOMAIN_UNIVERSAL)
	spello(SKILL_POISON, libc.CString("healing glow"), 50, 20, 3, POS_STANDING, (1<<1)|1<<6|1<<7, TRUE, (1<<14)|1<<1|1<<4, (1<<0)|1<<4, 0, libc.CString("Your healing glow fades."), 4, SCHOOL_NECROMANCY, DOMAIN_UNIVERSAL)
	spello(SKILL_ENLIGHTEN, libc.CString("enlighten"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("You feel less wise."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_GENIUS, libc.CString("genius"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("You am dumb dumbner now."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SKILL_FLEX, libc.CString("flex"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("You feel less agile."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SPELL_PORTAL, libc.CString("portal"), 75, 75, 0, POS_STANDING, (1<<2)|1<<6, FALSE, (1<<14)|1<<10, 0, 0, nil, 7, SCHOOL_CONJURATION, -1)
	spello(SPELL_PROT_FROM_EVIL, libc.CString("protection from evil"), 40, 10, 3, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("You feel less protected."), 1, SCHOOL_ABJURATION, DOMAIN_GOOD)
	spello(SPELL_REMOVE_CURSE, libc.CString("remove curse"), 45, 25, 5, POS_STANDING, (1<<1)|1<<7|1<<10, FALSE, (1<<14)|1<<2|1<<4, 0, 0, nil, 3, SCHOOL_ABJURATION, DOMAIN_UNIVERSAL)
	spello(SPELL_NEUTRALIZE_POISON, libc.CString("neutralize poison"), 40, 8, 4, POS_STANDING, (1<<1)|1<<7|1<<8, FALSE, (1<<14)|1<<2|1<<4, 0, 0, nil, 4, SCHOOL_CONJURATION, DOMAIN_UNIVERSAL)
	spello(SPELL_SANCTUARY, libc.CString("sanctuary"), 110, 85, 5, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("The white aura around your body fades."), 9, -1, DOMAIN_PROTECTION)
	spello(SPELL_SENSE_LIFE, libc.CString("sense life"), 20, 10, 2, POS_STANDING, (1<<1)|1<<5, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("You feel less aware of your surroundings."), 2, SCHOOL_DIVINATION, DOMAIN_UNIVERSAL)
	spello(SPELL_SHOCKING_GRASP, libc.CString("shocking grasp"), 30, 15, 3, POS_FIGHTING, (1<<1)|1<<4, TRUE, (1<<14)|1<<0, 0, 0, nil, 1, SCHOOL_TRANSMUTATION, -1)
	spello(SPELL_SLEEP, libc.CString("sleep"), 40, 25, 5, POS_FIGHTING, 1<<1, TRUE, (1<<14)|1<<1, (1<<2)|1<<4, 0, libc.CString("You feel like you can wake up again."), 1, SCHOOL_ENCHANTMENT, -1)
	spello(SPELL_HAYASA, libc.CString("hayasa"), 40, 25, 5, POS_FIGHTING, 1<<1, TRUE, (1<<14)|1<<1, (1<<2)|1<<4, 0, libc.CString("You feel your speed decrease as Hayasa fades."), 1, SCHOOL_ENCHANTMENT, -1)
	spello(SKILL_TSKIN, libc.CString("tough skin"), 25, 10, 1, POS_STANDING, (1<<1)|1<<6, TRUE, (1<<14)|1<<11, (1<<0)|1<<4, 0, libc.CString("Your skin isn't quite so thick anymore."), 2, SCHOOL_TRANSMUTATION, -1)
	spello(SPELL_BULL_STRENGTH, libc.CString("bull strength"), 35, 30, 1, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("You feel weaker."), 2, SCHOOL_TRANSMUTATION, DOMAIN_STRENGTH)
	spello(SPELL_SUMMON, libc.CString("summon"), 75, 50, 3, POS_STANDING, (1<<2)|1<<6, FALSE, (1<<14)|1<<10, 0, 0, nil, 7, SCHOOL_CONJURATION, DOMAIN_UNIVERSAL)
	spello(SPELL_TELEPORT, libc.CString("teleport"), 75, 50, 3, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<10, 0, 0, nil, 5, SCHOOL_TRANSMUTATION, DOMAIN_TRAVEL)
	spello(SPELL_WATERWALK, libc.CString("waterwalk"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, libc.CString("Your feet seem less buoyant."), 3, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_WORD_OF_RECALL, libc.CString("word of recall"), 20, 10, 2, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<10, 0, 0, nil, 6, -1, DOMAIN_UNIVERSAL)
	spello(SPELL_RESISTANCE, libc.CString("resistance"), 40, 20, 0, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<11, 0, (1<<3)|1<<4|1<<5|1<<0, nil, 0, SCHOOL_ABJURATION, -1)
	spello(SPELL_ACID_SPLASH, libc.CString("acid splash"), 40, 20, 0, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, (1<<4)|1<<5, nil, 0, SCHOOL_CONJURATION, -1)
	spello(SPELL_DAZE, libc.CString("daze"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<11, (1<<2)|1<<4, (1<<3)|1<<4|1<<5, nil, 0, SCHOOL_ENCHANTMENT, -1)
	spello(SPELL_FLARE, libc.CString("flare"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<11, (1<<0)|1<<4, 1<<5, nil, 0, SCHOOL_EVOCATION, -1)
	spello(SPELL_RAY_OF_FROST, libc.CString("ray of frost"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, (1<<4)|1<<5, nil, 0, SCHOOL_EVOCATION, -1)
	spello(SPELL_DISRUPT_UNDEAD, libc.CString("disrupt undead"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, (1<<4)|1<<5, nil, 0, SCHOOL_NECROMANCY, -1)
	spello(SPELL_SUMMON_MONSTER_I, libc.CString("summon monster i"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 1, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_II, libc.CString("summon monster ii"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 2, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_III, libc.CString("summon monster iii"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 3, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_IV, libc.CString("summon monster iv"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 4, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_V, libc.CString("summon monster v"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 5, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_VI, libc.CString("summon monster vi"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 6, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_VII, libc.CString("summon monster vii"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 7, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_VIII, libc.CString("summon monster viii"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 8, SCHOOL_CONJURATION, -1)
	spello(SPELL_SUMMON_MONSTER_IX, libc.CString("summon monster ix"), 40, 20, 2, POS_FIGHTING, 1<<0, FALSE, (1<<14)|1<<8, 0, 0, nil, 9, SCHOOL_CONJURATION, -1)
	spello(SPELL_FIRE_SHIELD, libc.CString("fire shield"), 40, 20, 2, POS_FIGHTING, 1<<1, FALSE, (1<<14)|1<<11, 0, 0, nil, 4, SCHOOL_EVOCATION, DOMAIN_FIRE)
	spello(SPELL_ICE_STORM, libc.CString("ice storm"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_SHOUT, libc.CString("shout"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_FEAR, libc.CString("fear"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_CLOUDKILL, libc.CString("cloudkill"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_MAJOR_CREATION, libc.CString("major creation"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_HOLD_MONSTER, libc.CString("hold monster"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_CONE_OF_COLD, libc.CString("cone of cold"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_ANIMAL_GROWTH, libc.CString("animal growth"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_BALEFUL_POLYMORPH, libc.CString("baleful polymorph"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_PASSWALL, libc.CString("passwall"), 40, 20, 2, POS_STANDING, 1<<1, FALSE, (1<<14)|1<<1, 0, 0, nil, 0, -1, -1)
	spello(SPELL_FIRE_BREATH, libc.CString("fire breath"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	spello(SPELL_GAS_BREATH, libc.CString("gas breath"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	spello(SPELL_FROST_BREATH, libc.CString("frost breath"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	spello(SPELL_ACID_BREATH, libc.CString("acid breath"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	spello(SPELL_LIGHTNING_BREATH, libc.CString("lightning breath"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	spello(SPELL_DG_AFFECT, libc.CString("Script-inflicted"), 0, 0, 0, POS_SITTING, 1<<0, TRUE, 0, 0, 0, nil, 0, 0, 0)
	skillo(SKILL_FLEX, libc.CString("flex"), (1<<6)|1<<0)
	skillo(SKILL_GENIUS, libc.CString("genius"), (1<<4)|1<<0)
	skillo(SKILL_ENLIGHTEN, libc.CString("enlighten"), (1<<4)|1<<0)
	skillo(SKILL_TSKIN, libc.CString("tough skin"), (1<<1)|1<<8)
	skillo(SKILL_KAIOKEN, libc.CString("kaioken"), (1<<4)|1<<0)
	skillo(SKILL_BLESS, libc.CString("bless"), (1<<5)|1<<0)
	skillo(SKILL_CURSE, libc.CString("curse"), (1<<5)|1<<0)
	skillo(SKILL_POISON, libc.CString("poison"), (1<<5)|1<<0)
	skillo(SKILL_VIGOR, libc.CString("vigor"), (1<<5)|1<<0)
	skillo(SKILL_POSE, libc.CString("special pose"), (1<<4)|1<<0)
	skillo(SKILL_HASSHUKEN, libc.CString("hasshuken"), (1<<4)|1<<0)
	skillo(SKILL_GARDENING, libc.CString("gardening"), (1<<4)|1<<0)
	skillo(SKILL_EXTRACT, libc.CString("extract"), (1<<4)|1<<0)
	skillo(SKILL_RUNIC, libc.CString("runic"), (1<<4)|1<<0)
	skillo(SKILL_COMMUNE, libc.CString("commune"), (1<<4)|1<<0)
	skillo(SKILL_SOLARF, libc.CString("solar flare"), (1<<4)|1<<0)
	skillo(SKILL_MIGHT, libc.CString("might"), (1<<1)|1<<0)
	skillo(SKILL_BALANCE, libc.CString("balance"), (1<<2)|1<<8)
	skillo(SKILL_BUILD, libc.CString("build"), (1<<4)|1<<0)
	skillo(SKILL_CONCENTRATION, libc.CString("concentration"), 1<<3)
	skillo(SKILL_SPOT, libc.CString("spot"), 1<<5)
	skillo(SKILL_FIRST_AID, libc.CString("first aid"), (1<<5)|1<<0)
	skillo(SKILL_DISGUISE, libc.CString("disguise"), 1<<6)
	skillo(SKILL_ESCAPE_ARTIST, libc.CString("escape"), (1<<2)|1<<8)
	skillo(SKILL_APPRAISE, libc.CString("appraise"), 1<<4)
	skillo(SKILL_HEAL, libc.CString("heal"), (1<<5)|1<<7)
	skillo(SKILL_FORGERY, libc.CString("forgery"), 1<<4)
	skillo(SKILL_HIDE, libc.CString("hide"), (1<<2)|1<<8)
	skillo(SKILL_LISTEN, libc.CString("listen"), 1<<5)
	skillo(SKILL_EAVESDROP, libc.CString("eavesdrop"), 1<<4)
	skillo(SKILL_CURE, libc.CString("cure poison"), (1<<5)|1<<0)
	skillo(SKILL_OPEN_LOCK, libc.CString("open lock"), (1<<2)|1<<0|1<<7)
	skillo(SKILL_REGENERATE, libc.CString("regenerate"), (1<<3)|1<<0)
	skillo(SKILL_KEEN, libc.CString("keen sight"), (1<<4)|1<<0)
	skillo(SKILL_SEARCH, libc.CString("search"), 1<<4)
	skillo(SKILL_MOVE_SILENTLY, libc.CString("move silently"), (1<<2)|1<<8)
	skillo(SKILL_ABSORB, libc.CString("absorb"), (1<<4)|1<<0)
	skillo(SKILL_SLEIGHT_OF_HAND, libc.CString("sleight of hand"), (1<<2)|1<<8)
	skillo(SKILL_INGEST, libc.CString("ingest"), (1<<1)|1<<0)
	skillo(SKILL_REPAIR, libc.CString("fix"), (1<<4)|1<<0)
	skillo(SKILL_SENSE, libc.CString("sense"), (1<<4)|1<<0)
	skillo(SKILL_SURVIVAL, libc.CString("survival"), (1<<5)|1<<0)
	skillo(SKILL_YOIK, libc.CString("yoikominminken"), (1<<4)|1<<0)
	skillo(SKILL_CREATE, libc.CString("create"), (1<<4)|1<<0)
	skillo(SKILL_SPIT, libc.CString("stone spit"), (1<<4)|1<<0)
	skillo(SKILL_POTENTIAL, libc.CString("potential release"), (1<<4)|1<<0)
	skillo(SKILL_TELEPATHY, libc.CString("telepathy"), (1<<4)|1<<0)
	skillo(SKILL_FOCUS, libc.CString("focus"), (1<<4)|1<<0)
	skillo(SKILL_INSTANTT, libc.CString("instant transmission"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_SWORD, libc.CString("sword"), (1<<4)|1<<0)
	skillo(SKILL_DAGGER, libc.CString("dagger"), (1<<4)|1<<0)
	skillo(SKILL_CLUB, libc.CString("club"), (1<<4)|1<<0)
	skillo(SKILL_SPEAR, libc.CString("spear"), (1<<4)|1<<0)
	skillo(SKILL_GUN, libc.CString("gun"), (1<<4)|1<<0)
	skillo(SKILL_BRAWL, libc.CString("brawl"), (1<<4)|1<<0)
	skillo(SKILL_DODGE, libc.CString("dodge"), (1<<6)|1<<0)
	skillo(SKILL_PARRY, libc.CString("parry"), (1<<2)|1<<0)
	skillo(SKILL_BLOCK, libc.CString("block"), (1<<2)|1<<0)
	skillo(SKILL_ZANZOKEN, libc.CString("zanzoken"), (1<<4)|1<<0)
	skillo(SKILL_BARRIER, libc.CString("barrier"), (1<<4)|1<<0)
	skillo(SKILL_THROW, libc.CString("throw"), (1<<2)|1<<0)
	skillo(SKILL_PUNCH, libc.CString("punch"), (1<<1)|1<<0)
	skillo(SKILL_KICK, libc.CString("kick"), (1<<1)|1<<0)
	skillo(SKILL_ELBOW, libc.CString("elbow"), (1<<1)|1<<0)
	skillo(SKILL_KNEE, libc.CString("knee"), (1<<1)|1<<0)
	skillo(SKILL_ROUNDHOUSE, libc.CString("roundhouse"), (1<<1)|1<<0)
	skillo(SKILL_UPPERCUT, libc.CString("uppercut"), (1<<1)|1<<0)
	skillo(SKILL_SLAM, libc.CString("slam"), (1<<1)|1<<0|1<<10)
	skillo(SKILL_HEELDROP, libc.CString("heeldrop"), (1<<1)|1<<0|1<<10)
	skillo(SKILL_KIBALL, libc.CString("kiball"), (1<<4)|1<<0)
	skillo(SKILL_KIBLAST, libc.CString("kiblast"), (1<<4)|1<<0)
	skillo(SKILL_BEAM, libc.CString("beam"), (1<<4)|1<<0)
	skillo(SKILL_TSUIHIDAN, libc.CString("tsuihidan"), (1<<4)|1<<0)
	skillo(SKILL_SHOGEKIHA, libc.CString("shogekiha"), (1<<4)|1<<0)
	skillo(SKILL_RENZO, libc.CString("renzokou energy dan"), (1<<4)|1<<0)
	skillo(SKILL_MASENKO, libc.CString("masenko"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_DODONPA, libc.CString("dodonpa"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_GALIKGUN, libc.CString("galik gun"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_KAMEHAMEHA, libc.CString("kamehameha"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_DEATHBEAM, libc.CString("deathbeam"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_ERASER, libc.CString("eraser cannon"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_TSLASH, libc.CString("twin slash"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_PSYBLAST, libc.CString("psychic blast"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_HONOO, libc.CString("honoo"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_DUALBEAM, libc.CString("dual beam"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_ROGAFUFUKEN, libc.CString("rogafufuken"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_BAKUHATSUHA, libc.CString("bakuhatsuha"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_KIENZAN, libc.CString("kienzan"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_TRIBEAM, libc.CString("tribeam"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_SBC, libc.CString("special beam cannon"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_FINALFLASH, libc.CString("final flash"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_CRUSHER, libc.CString("crusher ball"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_DDSLASH, libc.CString("darkness dragon slash"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_PBARRAGE, libc.CString("psychic barrage"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_HELLFLASH, libc.CString("hell flash"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_HELLSPEAR, libc.CString("hell spear blast"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_KAKUSANHA, libc.CString("kakusanha"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_SCATTER, libc.CString("scatter shot"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_BIGBANG, libc.CString("big bang"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_PSLASH, libc.CString("phoenix slash"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_DEATHBALL, libc.CString("deathball"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_SPIRITBALL, libc.CString("spirit ball"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_GENKIDAMA, libc.CString("genki dama"), (1<<4)|1<<0|1<<13)
	skillo(SKILL_GENOCIDE, libc.CString("genocide"), (1<<4)|1<<0|1<<13)
	skillo(SKILL_DUALWIELD, libc.CString("dual wield"), (1<<4)|1<<0)
	skillo(SKILL_TWOHAND, libc.CString("twohand"), (1<<4)|1<<0)
	skillo(SKILL_STYLE, libc.CString("fighting arts"), 1<<4)
	skillo(SKILL_KURA, libc.CString("kuraiiro seiki"), (1<<4)|1<<0)
	skillo(SKILL_TAILWHIP, libc.CString("tailwhip"), (1<<4)|1<<0|1<<9)
	skillo(SKILL_KOUSENGAN, libc.CString("kousengan"), (1<<4)|1<<0|1<<9)
	skillo(SKILL_TAISHA, libc.CString("taisha reiki"), (1<<4)|1<<0)
	skillo(SKILL_PARALYZE, libc.CString("paralyze"), (1<<4)|1<<0)
	skillo(SKILL_INFUSE, libc.CString("infuse"), (1<<4)|1<<0)
	skillo(SKILL_ROLL, libc.CString("roll"), (1<<4)|1<<0)
	skillo(SKILL_TRIP, libc.CString("trip"), (1<<4)|1<<0)
	skillo(SKILL_GRAPPLE, libc.CString("grapple"), (1<<4)|1<<0)
	skillo(SKILL_WSPIKE, libc.CString("water spikes"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_SELFD, libc.CString("self destruct"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_SPIRAL, libc.CString("spiral comet"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_BREAKER, libc.CString("star breaker"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_MIMIC, libc.CString("mimic"), (1<<4)|1<<0)
	skillo(SKILL_WRAZOR, libc.CString("water razor"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_KOTEIRU, libc.CString("koteiru bakuha"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_DIMIZU, libc.CString("dimizu toride"), (1<<4)|1<<0)
	skillo(SKILL_HYOGA_KABE, libc.CString("hyoga kabe"), (1<<4)|1<<0)
	skillo(SKILL_WELLSPRING, libc.CString("wellspring"), (1<<4)|1<<0)
	skillo(SKILL_AQUA_BARRIER, libc.CString("aqua barrier"), (1<<4)|1<<0)
	skillo(SKILL_WARP, libc.CString("warp pool"), (1<<4)|1<<0)
	skillo(SKILL_HSPIRAL, libc.CString("hell spiral"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_ARMOR, libc.CString("nanite armor"), (1<<4)|1<<0)
	skillo(SKILL_FIRESHIELD, libc.CString("fireshield"), (1<<4)|1<<0)
	skillo(SKILL_COOKING, libc.CString("cooking"), (1<<4)|1<<0)
	skillo(SKILL_SEISHOU, libc.CString("seishou enko"), (1<<4)|1<<0|1<<10)
	skillo(SKILL_SILK, libc.CString("silk"), (1<<4)|1<<0)
	skillo(SKILL_BASH, libc.CString("bash"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_HEADBUTT, libc.CString("headbutt"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_ENSNARE, libc.CString("ensnare"), (1<<4)|1<<0)
	skillo(SKILL_STARNOVA, libc.CString("starnova"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_MALICE, libc.CString("malice breaker"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_ZEN, libc.CString("zen blade strike"), (1<<4)|1<<0|1<<11)
	skillo(SKILL_SUNDER, libc.CString("sundering force"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_WITHER, libc.CString("wither"), (1<<4)|1<<0)
	skillo(SKILL_METAMORPH, libc.CString("dark metamorphosis"), (1<<4)|1<<0)
	skillo(SKILL_HAYASA, libc.CString("hayasa"), (1<<4)|1<<0)
	skillo(SKILL_ENERGIZE, libc.CString("energize throwing"), (1<<4)|1<<0)
	skillo(SKILL_PURSUIT, libc.CString("pursuit"), (1<<4)|1<<0)
	skillo(SKILL_HEALGLOW, libc.CString("healing glow"), (1<<4)|1<<0)
	skillo(SKILL_HANDLING, libc.CString("handling"), (1<<4)|1<<0)
	skillo(SKILL_MYSTICMUSIC, libc.CString("mystic music"), (1<<4)|1<<0)
	skillo(SKILL_LIGHTGRENADE, libc.CString("light grenade"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_MULTIFORM, libc.CString("multiform"), (1<<4)|1<<0|1<<9)
	skillo(SKILL_SPIRITCONTROL, libc.CString("spirit control"), (1<<4)|1<<0|1<<9)
	skillo(SKILL_BALEFIRE, libc.CString("balefire"), (1<<4)|1<<0|1<<12)
	skillo(SKILL_BLESSEDHAMMER, libc.CString("blessed hammer"), (1<<4)|1<<0|1<<9)
}
