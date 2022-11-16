package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unicode"
	"unsafe"
)

const GW_ARRAY_MAX = 4

type guild_data struct {
	Vnum            room_vnum
	Skills          [1000]int
	Charge          float32
	No_such_skill   *byte
	Not_enough_gold *byte
	Minlvl          int
	Gm              mob_rnum
	With_who        [4]bitvector_t
	Open            int
	Close           int
	Func            SpecialFunc
	Feats           [252]int
}

var spell_sort_info [1001]int
var top_guild int = -1
var guild_index []guild_data

func calculate_skill_cost(ch *char_data, skill int) int {
	var cost int = 0
	if IS_SET(bitvector_t(int32(spell_info[skill].Flags)), 1<<10) {
		cost = 8
	} else if IS_SET(bitvector_t(int32(spell_info[skill].Flags)), 1<<11) {
		cost = 15
	} else if IS_SET(bitvector_t(int32(spell_info[skill].Flags)), 1<<12) {
		if int(ch.Skills[skill]) == 0 {
			cost = 200
		} else {
			cost = 25
		}
	} else if IS_SET(bitvector_t(int32(spell_info[skill].Flags)), 1<<13) {
		if int(ch.Skills[skill]) == 0 {
			cost = 300
		} else {
			cost = 40
		}
	} else {
		cost = 4
	}
	if int(ch.Skills[skill]) > 90 {
		cost += 12
	} else if int(ch.Skills[skill]) > 80 {
		cost += 10
	} else if int(ch.Skills[skill]) > 70 {
		cost += 8
	} else if int(ch.Skills[skill]) > 50 {
		cost += 6
	} else if int(ch.Skills[skill]) > 40 {
		cost += 2
	} else if int(ch.Skills[skill]) > 30 {
		cost += 1
	}
	if ch.Forgeting != 0 {
		cost += 6
	}
	if skill == SKILL_RUNIC {
		cost += 6
	}
	if skill == SKILL_EXTRACT {
		cost += 3
	}
	if int(ch.Race) == RACE_HOSHIJIN && (skill == SKILL_PUNCH || skill == SKILL_KICK || skill == SKILL_KNEE || skill == SKILL_ELBOW || skill == SKILL_UPPERCUT || skill == SKILL_ROUNDHOUSE || skill == SKILL_SLAM || skill == SKILL_HEELDROP || skill == SKILL_DAGGER || skill == SKILL_SWORD || skill == SKILL_CLUB || skill == SKILL_GUN || skill == SKILL_SPEAR || skill == SKILL_BRAWL) {
		cost += 5
	}
	if skill == SKILL_INSTANTT {
		if int(ch.Skills[skill]) == 0 {
			cost = 2000
		} else {
			cost = 50
		}
	}
	if skill == SKILL_MYSTICMUSIC {
		cost = int(float64(cost) * 1.5)
	}
	return cost
}
func handle_ingest_learn(ch *char_data, vict *char_data) {
	var i int = 1
	send_to_char(ch, libc.CString("@YAll your current skills improve somewhat!@n\r\n"))
	for i = 1; i <= SKILL_TABLE_SIZE; i++ {
		if int(ch.Skills[i]) > 0 && int(vict.Skills[i]) > 0 && i != 141 {
			send_to_char(ch, libc.CString("@YYou gained a lot of new knowledge about @y%s@Y!@n\r\n"), spell_info[i].Name)
			if int(ch.Skills[i])+10 < 100 {
				ch.Skills[i] += 10
			} else if int(ch.Skills[i]) > 0 && int(ch.Skills[i]) < 100 {
				ch.Skills[i] += 1
			} else {
				ch.Skills[i] = 100
			}
		}
		if (i >= 481 && i <= 489 || i == 517 || i == 535) && (int(ch.Skills[i]) <= 0 && int(vict.Skills[i]) > 0) {
			for {
				ch.Skills[i] = int8(int(ch.Skills[i]) + rand_number(10, 25))
				if true {
					break
				}
			}
			send_to_char(ch, libc.CString("@YYou learned @y%s@Y from ingesting your target!@n\r\n"), spell_info[i].Name)
			ch.Skill_slots += 1
			ch.IngestLearned = 1
		}
	}
}
func do_teach(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var arg [2048]byte
	var arg2 [2048]byte
	var skill int = 100
	var vict *char_data
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What skill are you wanting to teach?\r\n"))
		return
	}
	skill = find_skill_num(&arg[0], 1<<1)
	if int(ch.Skills[skill]) < 101 {
		send_to_char(ch, libc.CString("You are not a Grand Master in that skill!\r\n"))
		send_to_char(ch, libc.CString("@wSyntax: teach (skill) (target)@n\r\n"))
		return
	}
	if arg2[0] == 0 {
		send_to_char(ch, libc.CString("@wWho are you wanting to teach @C%s@w to?@n\r\n"), spell_info[skill].Name)
		send_to_char(ch, libc.CString("@wSyntax: teach (skill) (target)@n\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("@wTeach who?@n\r\n"))
		send_to_char(ch, libc.CString("@wSyntax: teach (skill) (target)@n\r\n"))
		return
	}
	var cost int = calculate_skill_cost(vict, skill)
	var free int = FALSE
	if int(ch.Skills[skill]) >= 103 {
		cost = int(float64(cost) * 0.5)
		if rand_number(1, 4) == 4 {
			free = TRUE
		}
	} else if int(ch.Skills[skill]) == 102 {
		cost = int(float64(cost) * 0.5)
	} else {
		cost = int(float64(cost) * 0.75)
	}
	if cost == 0 {
		cost = 1
	}
	if vict.Master == nil {
		send_to_char(ch, libc.CString("They must be following you in order for you to teach them.\r\n"))
		return
	} else if vict.Master != ch {
		send_to_char(ch, libc.CString("They must be following you in order for you to teach them.\r\n"))
		return
	} else if vict.Forgeting == skill {
		send_to_char(ch, libc.CString("They are trying to forget that skill!\r\n"))
		return
	} else if (vict.Player_specials.Class_skill_points[vict.Chclass]) < cost {
		send_to_char(ch, libc.CString("They do not have enough practice sessions for you to teach them.\r\n"))
		return
	} else if int(vict.Skills[skill]) >= 80 {
		send_to_char(ch, libc.CString("You can not teach them anymore.\r\n"))
		return
	} else if int(vict.Skills[skill]) > 0 {
		var (
			tochar  [64936]byte
			tovict  [64936]byte
			toother [64936]byte
		)
		stdio.Sprintf(&tochar[0], "@YYou instruct @y$N@Y in the finer points of @C%s@Y.@n\r\n", spell_info[skill].Name)
		stdio.Sprintf(&tovict[0], "@y$n@Y instructs you in the finer points of @C%s@Y.@n\r\n", spell_info[skill].Name)
		stdio.Sprintf(&toother[0], "@y$n@Y instructs @y$N@Y in the finer points of @C%s@Y.@n\r\n", spell_info[skill].Name)
		act(&tochar[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(&tovict[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(&toother[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		for {
			vict.Skills[skill] = int8(int(vict.Skills[skill]) + 1)
			if true {
				break
			}
		}
		if free == FALSE {
			vict.Player_specials.Class_skill_points[vict.Chclass] -= cost
		} else {
			send_to_char(ch, libc.CString("@GYou teach your lesson so well that it cost them nothing to learn from you!@n\r\n"))
			send_to_char(vict, libc.CString("@GYour teacher taught you the lesson so well that it cost you nothing!@n\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("They do not even know the basics. It's a waste of your teaching skills.\r\n"))
		return
	}
}
func how_good(percent int) *byte {
	if percent < 0 {
		return libc.CString(" error)")
	}
	if percent == 0 {
		return libc.CString("(@Mnot@n)")
	}
	if percent <= 10 {
		return libc.CString("(@rawful@n)")
	}
	if percent <= 20 {
		return libc.CString("(@Rbad@n)")
	}
	if percent <= 40 {
		return libc.CString("(@ypoor@n)")
	}
	if percent <= 55 {
		return libc.CString("(@Yaverage@n)")
	}
	if percent <= 70 {
		return libc.CString("(@gfair@n)")
	}
	if percent <= 80 {
		return libc.CString("(@Ggood@n)")
	}
	if percent <= 85 {
		return libc.CString("(@bgreat@n)")
	}
	if percent <= 100 {
		return libc.CString("(@Bsuperb@n)")
	}
	return libc.CString("(@rinate@n)")
}

var prac_types [2]*byte = [2]*byte{libc.CString("spell"), libc.CString("skill")}

func compare_spells(x unsafe.Pointer, y unsafe.Pointer) int {
	var (
		a int = *(*int)(x)
		b int = *(*int)(y)
	)
	return libc.StrCmp(spell_info[a].Name, spell_info[b].Name)
}
func print_skills_by_type(ch *char_data, buf *byte, maxsz int, sktype int, argument *byte) int {
	var (
		arg     [1000]byte
		len_    uint64 = 0
		i       int
		t       int
		known   int
		nlen    int = 0
		count   int = 0
		canknow int = 0
		buf2    [256]byte
	)
	one_argument(argument, &arg[0])
	for i = 1; i <= SKILL_TABLE_SIZE; i++ {
		t = spell_info[i].Skilltype
		if t != sktype {
			continue
		}
		if (t&(1<<1)) != 0 || (t&(1<<0)) != 0 {
			for func() int {
				nlen = 0
				return func() int {
					known = 0
					return known
				}()
			}(); nlen <= NUM_CLASSES; nlen++ {
				if ((ch.Chclasses[nlen])+(ch.Epicclasses[nlen])) > 0 && int(spell_info[i].Can_learn_skill[nlen]) > SKLEARN_CANT {
					known = int(spell_info[i].Can_learn_skill[nlen])
				}
			}
		} else {
			known = 0
		}
		if GET_SKILL(ch, i) <= 0 {
			known = 0
		}
		if arg[0] != 0 {
			if libc.Atoi(libc.GoString(&arg[0])) <= 0 && unsafe.Pointer(libc.StrStr(spell_info[i].Name, &arg[0])) == unsafe.Pointer(uintptr(FALSE)) {
				known = 0
			} else if libc.Atoi(libc.GoString(&arg[0])) > GET_SKILL(ch, i) {
				known = 0
			}
		}
		if known != 0 {
			if t&(1<<2) != 0 {
				nlen = stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(buf), len_)), maxsz-int(len_), "%-20s  (%s)\r\n", spell_info[i].Name, func() string {
					if int(ch.Skills[i]) != 0 {
						return "known"
					}
					return "unknown"
				}())
			} else if t&(1<<1) != 0 {
				if int(ch.Skillmods[i]) != 0 {
					stdio.Snprintf(&buf2[0], int(256), " (base %d + bonus %d)", ch.Skills[i], ch.Skillmods[i])
				} else {
					buf2[0] = 0
				}
				if known == SKLEARN_CROSSCLASS {
					count++
					canknow = highest_skill_value(GET_LEVEL(ch), GET_SKILL(ch, i))
					nlen = stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(buf), len_)), maxsz-int(len_), "@y(@Y%2d@y) @W%-30s  @y(@Y%2d@y) @C%3d@D/@c%3d   %s@n%s%s\r\n", count, spell_info[i].Name, count, GET_SKILL(ch, i), canknow, func() string {
						if int(ch.Skillperfs[i]) > 0 {
							if int(ch.Skillperfs[i]) == 1 {
								return "@ROver Charge"
							}
							if int(ch.Skillperfs[i]) == 2 {
								return "@BAccurate"
							}
							return "@GEfficient"
						}
						return ""
					}(), func() string {
						if int(ch.Skills[i]) > 100 {
							return " @D(@YGrand Master@D)@n"
						}
						return ""
					}(), &buf2[0])
				} else {
					count++
					canknow = highest_skill_value(GET_LEVEL(ch), GET_SKILL(ch, i))
					nlen = stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(buf), len_)), maxsz-int(len_), "@y(@Y%2d@y) @W%-30s  @y(@Y%2d@y) @C%3d@D/@c%d3   %s@n%s%s\r\n", count, spell_info[i].Name, count, GET_SKILL(ch, i), canknow, func() string {
						if int(ch.Skillperfs[i]) > 0 {
							if int(ch.Skillperfs[i]) == 1 {
								return "@ROver Charge"
							}
							if int(ch.Skillperfs[i]) == 2 {
								return "@BAccurate"
							}
							return "@GEfficient"
						}
						return ""
					}(), func() string {
						if int(ch.Skills[i]) > 100 {
							return " @D(@YGrand Master@D)@n"
						}
						return ""
					}(), &buf2[0])
				}
			}
			if len_+uint64(nlen) >= uint64(maxsz) || nlen < 0 {
				break
			}
			len_ += uint64(nlen)
		}
	}
	return int(len_)
}
func slot_count(ch *char_data) int {
	var (
		i         int
		skills    int = -1
		fail      int = FALSE
		punch     int = FALSE
		kick      int = FALSE
		knee      int = FALSE
		elbow     int = FALSE
		kiball    int = FALSE
		kiblast   int = FALSE
		beam      int = FALSE
		renzo     int = FALSE
		shogekiha int = FALSE
	)
	for i = 1; i <= SKILL_TABLE_SIZE; i++ {
		if GET_SKILL(ch, i) > 0 {
			switch i {
			case SKILL_PUNCH:
				fail = TRUE
				punch = TRUE
			case SKILL_KICK:
				fail = TRUE
				kick = TRUE
			case SKILL_ELBOW:
				fail = TRUE
				elbow = TRUE
			case SKILL_KNEE:
				fail = TRUE
				knee = TRUE
			case SKILL_KIBALL:
				fail = TRUE
				kiball = TRUE
			case SKILL_KIBLAST:
				fail = TRUE
				kiblast = TRUE
			case SKILL_BEAM:
				fail = TRUE
				beam = TRUE
			case SKILL_SHOGEKIHA:
				fail = TRUE
				shogekiha = TRUE
			case SKILL_RENZO:
				fail = TRUE
				renzo = TRUE
			case SKILL_TELEPATHY:
				if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_KAI {
					fail = TRUE
				}
			case SKILL_ABSORB:
				if int(ch.Race) == RACE_BIO || int(ch.Race) == RACE_ANDROID {
					fail = TRUE
				}
			case SKILL_TAILWHIP:
				if int(ch.Race) == RACE_ICER {
					fail = TRUE
				}
			case SKILL_SEISHOU:
				if int(ch.Race) == RACE_ARLIAN {
					fail = TRUE
				}
			case SKILL_REGENERATE:
				if int(ch.Race) == RACE_MAJIN || int(ch.Race) == RACE_NAMEK || int(ch.Race) == RACE_BIO {
					fail = TRUE
				}
			}
			if fail == FALSE {
				skills += 1
			}
			fail = FALSE
		}
	}
	if punch == TRUE && kick == TRUE && elbow == TRUE && knee == TRUE {
		skills += 1
	}
	if kiball == TRUE && kiblast == TRUE && beam == TRUE && shogekiha == TRUE && renzo == TRUE {
		skills += 1
	}
	return skills
}
func list_skills(ch *char_data, arg *byte) {
	var (
		overflow *byte  = libc.CString("\r\n**OVERFLOW**\r\n")
		len_     uint64 = 0
		slots    int    = FALSE
		buf2     [64936]byte
	)
	len_ = uint64(stdio.Snprintf(&buf2[0], int(64936), "You have %d practice session%s remaining.\r\n", ch.Player_specials.Class_skill_points[ch.Chclass], func() string {
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) == 1 {
			return ""
		}
		return "s"
	}()))
	len_ += uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "\r\nYou know the following skills:     @CKnown@D/@cPrac. Max@n\r\n@w-------------------------------------------------------@n\r\n"))
	len_ += uint64(print_skills_by_type(ch, &buf2[len_], int(64936-uintptr(len_)), 1<<1, arg))
	if slots == FALSE {
		len_ += uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "\r\n@DSkill Slots@W: @M%d@W/@m%d", slot_count(ch), ch.Skill_slots))
	}
	if len_ >= uint64(64936) {
		libc.StrCpy((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&buf2[64935]), -libc.StrLen(overflow)))), -1)), overflow)
	}
	page_string(ch.Desc, &buf2[0], TRUE)
}
func is_guild_open(keeper *char_data, guild_nr int, msg int) int {
	var buf [200]byte
	buf[0] = 0
	if guild_index[guild_nr].Open > time_info.Hours && guild_index[guild_nr].Close < time_info.Hours {
		strlcpy(&buf[0], libc.CString(MSG_TRAINER_NOT_OPEN), uint64(200))
	}
	if buf[0] == 0 {
		return TRUE
	}
	if msg != 0 {
		do_say(keeper, &buf[0], cmd_tell, 0)
	}
	return FALSE
}
func is_guild_ok_char(keeper *char_data, ch *char_data, guild_nr int) int {
	var buf [200]byte
	if !CAN_SEE(keeper, ch) {
		do_say(keeper, libc.CString(MSG_TRAINER_NO_SEE_CH), cmd_say, 0)
		return FALSE
	}
	if GET_LEVEL(ch) < guild_index[guild_nr].Minlvl {
		stdio.Snprintf(&buf[0], int(200), "%s %s", GET_NAME(ch), MSG_TRAINER_MINLVL)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return FALSE
	}
	if IS_GOOD(ch) && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOGOOD) || IS_EVIL(ch) && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOEVIL) || IS_NEUTRAL(ch) && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NONEUTRAL) {
		stdio.Snprintf(&buf[0], int(200), "%s %s", GET_NAME(ch), MSG_TRAINER_DISLIKE_ALIGN)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return FALSE
	}
	if IS_NPC(ch) {
		return FALSE
	}
	if int(ch.Chclass) == CLASS_ROSHI && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOWIZARD) || int(ch.Chclass) == CLASS_PICCOLO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOCLERIC) || int(ch.Chclass) == CLASS_KRANE && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOROGUE) || int(ch.Chclass) == CLASS_NAIL && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOFIGHTER) || int(ch.Chclass) == CLASS_GINYU && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOPALADIN) || int(ch.Chclass) == CLASS_FRIEZA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOSORCERER) || int(ch.Chclass) == CLASS_TAPION && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NODRUID) || int(ch.Chclass) == CLASS_ANDSIX && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOBARD) || int(ch.Chclass) == CLASS_DABURA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NORANGER) || int(ch.Chclass) == CLASS_BARDOCK && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOMONK) || int(ch.Chclass) == CLASS_KABITO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOBARBARIAN) || int(ch.Chclass) == CLASS_JINTO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOARCANE_ARCHER) || int(ch.Chclass) == CLASS_TSUNA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOARCANE_TRICKSTER) || int(ch.Chclass) == CLASS_KURZAK && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOARCHMAGE) || ((ch.Chclasses[CLASS_ASSASSIN])+(ch.Epicclasses[CLASS_ASSASSIN])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOASSASSIN) || ((ch.Chclasses[CLASS_BLACKGUARD])+(ch.Epicclasses[CLASS_BLACKGUARD])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOBLACKGUARD) || ((ch.Chclasses[CLASS_DRAGON_DISCIPLE])+(ch.Epicclasses[CLASS_DRAGON_DISCIPLE])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NODRAGON_DISCIPLE) || ((ch.Chclasses[CLASS_DUELIST])+(ch.Epicclasses[CLASS_DUELIST])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NODUELIST) || ((ch.Chclasses[CLASS_DWARVEN_DEFENDER])+(ch.Epicclasses[CLASS_DWARVEN_DEFENDER])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NODWARVEN_DEFENDER) || ((ch.Chclasses[CLASS_ELDRITCH_KNIGHT])+(ch.Epicclasses[CLASS_ELDRITCH_KNIGHT])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOELDRITCH_KNIGHT) || ((ch.Chclasses[CLASS_HIEROPHANT])+(ch.Epicclasses[CLASS_HIEROPHANT])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOHIEROPHANT) || ((ch.Chclasses[CLASS_HORIZON_WALKER])+(ch.Epicclasses[CLASS_HORIZON_WALKER])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOHORIZON_WALKER) || ((ch.Chclasses[CLASS_LOREMASTER])+(ch.Epicclasses[CLASS_LOREMASTER])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOLOREMASTER) || ((ch.Chclasses[CLASS_MYSTIC_THEURGE])+(ch.Epicclasses[CLASS_MYSTIC_THEURGE])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOMYSTIC_THEURGE) || ((ch.Chclasses[CLASS_SHADOWDANCER])+(ch.Epicclasses[CLASS_SHADOWDANCER])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOSHADOWDANCER) || ((ch.Chclasses[CLASS_THAUMATURGIST])+(ch.Epicclasses[CLASS_THAUMATURGIST])) > 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOTHAUMATURGIST) {
		stdio.Snprintf(&buf[0], int(200), "%s %s", GET_NAME(ch), MSG_TRAINER_DISLIKE_CLASS)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return FALSE
	}
	if int(ch.Chclass) != CLASS_ROSHI && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYWIZARD) || int(ch.Chclass) != CLASS_PICCOLO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYCLERIC) || int(ch.Chclass) != CLASS_KRANE && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYROGUE) || int(ch.Chclass) != CLASS_BARDOCK && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYMONK) || int(ch.Chclass) != CLASS_GINYU && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYPALADIN) || int(ch.Chclass) != CLASS_NAIL && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYFIGHTER) || int(ch.Chclass) != CLASS_FRIEZA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYSORCERER) || int(ch.Chclass) != CLASS_TAPION && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYDRUID) || int(ch.Chclass) != CLASS_ANDSIX && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYBARD) || int(ch.Chclass) != CLASS_DABURA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYRANGER) || int(ch.Chclass) != CLASS_KABITO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYBARBARIAN) || int(ch.Chclass) != CLASS_JINTO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYARCANE_ARCHER) || int(ch.Chclass) != CLASS_TSUNA && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYARCANE_TRICKSTER) || int(ch.Chclass) != CLASS_KURZAK && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYARCHMAGE) || ((ch.Chclasses[CLASS_ASSASSIN])+(ch.Epicclasses[CLASS_ASSASSIN])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYASSASSIN) || ((ch.Chclasses[CLASS_BLACKGUARD])+(ch.Epicclasses[CLASS_BLACKGUARD])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYBLACKGUARD) || ((ch.Chclasses[CLASS_DRAGON_DISCIPLE])+(ch.Epicclasses[CLASS_DRAGON_DISCIPLE])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYDRAGON_DISCIPLE) || ((ch.Chclasses[CLASS_DUELIST])+(ch.Epicclasses[CLASS_DUELIST])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYDUELIST) || ((ch.Chclasses[CLASS_DWARVEN_DEFENDER])+(ch.Epicclasses[CLASS_DWARVEN_DEFENDER])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYDWARVEN_DEFENDER) || ((ch.Chclasses[CLASS_ELDRITCH_KNIGHT])+(ch.Epicclasses[CLASS_ELDRITCH_KNIGHT])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYELDRITCH_KNIGHT) || ((ch.Chclasses[CLASS_HIEROPHANT])+(ch.Epicclasses[CLASS_HIEROPHANT])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYHIEROPHANT) || ((ch.Chclasses[CLASS_HORIZON_WALKER])+(ch.Epicclasses[CLASS_HORIZON_WALKER])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYHORIZON_WALKER) || ((ch.Chclasses[CLASS_LOREMASTER])+(ch.Epicclasses[CLASS_LOREMASTER])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYLOREMASTER) || ((ch.Chclasses[CLASS_MYSTIC_THEURGE])+(ch.Epicclasses[CLASS_MYSTIC_THEURGE])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYMYSTIC_THEURGE) || ((ch.Chclasses[CLASS_SHADOWDANCER])+(ch.Epicclasses[CLASS_SHADOWDANCER])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYSHADOWDANCER) || ((ch.Chclasses[CLASS_THAUMATURGIST])+(ch.Epicclasses[CLASS_THAUMATURGIST])) <= 0 && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_ONLYTHAUMATURGIST) {
		stdio.Snprintf(&buf[0], int(200), "%s %s", GET_NAME(ch), MSG_TRAINER_DISLIKE_CLASS)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return FALSE
	}
	if int(ch.Race) == RACE_HUMAN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOHUMAN) || int(ch.Race) == RACE_SAIYAN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOSAIYAN) || int(ch.Race) == RACE_ICER && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOICER) || int(ch.Race) == RACE_KONATSU && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOKONATSU) || int(ch.Race) == RACE_NAMEK && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NONAMEK) || int(ch.Race) == RACE_MUTANT && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOMUTANT) || int(ch.Race) == RACE_KANASSAN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOKANASSAN) || int(ch.Race) == RACE_ANDROID && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOANDROID) || int(ch.Race) == RACE_BIO && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOBIO) || int(ch.Race) == RACE_DEMON && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NODEMON) || int(ch.Race) == RACE_MAJIN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOMAJIN) || int(ch.Race) == RACE_KAI && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOKAI) || int(ch.Race) == RACE_TRUFFLE && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOTRUFFLE) || int(ch.Race) == RACE_HOSHIJIN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOGOBLIN) || int(ch.Race) == RACE_ANIMAL && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOANIMAL) || int(ch.Race) == RACE_ORC && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOORC) || int(ch.Race) == RACE_SNAKE && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOSNAKE) || int(ch.Race) == RACE_TROLL && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOTROLL) || int(ch.Race) == RACE_HALFBREED && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOHALFBREED) || int(ch.Race) == RACE_MINOTAUR && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOMINOTAUR) || int(ch.Race) == RACE_ARLIAN && IS_SET_AR(guild_index[guild_nr].With_who[:], TRADE_NOKOBOLD) {
		stdio.Snprintf(&buf[0], int(200), "%s %s", GET_NAME(ch), MSG_TRAINER_DISLIKE_RACE)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return FALSE
	}
	return TRUE
}
func is_guild_ok(keeper *char_data, ch *char_data, guild_nr int) int {
	if is_guild_open(keeper, guild_nr, TRUE) != 0 {
		return is_guild_ok_char(keeper, ch, guild_nr)
	}
	return FALSE
}
func does_guild_know(guild_nr int, i int) int {
	return guild_index[guild_nr].Skills[i]
}
func does_guild_know_feat(guild_nr int, i int) int {
	return guild_index[guild_nr].Feats[i]
}
func sort_spells() {
	var a int
	for a = 1; a < SKILL_TABLE_SIZE; a++ {
		spell_sort_info[a] = a
	}
	libc.Sort(unsafe.Pointer(&spell_sort_info[1]), SKILL_TABLE_SIZE, uint32(unsafe.Sizeof(int(0))), func(arg1 unsafe.Pointer, arg2 unsafe.Pointer) int32 {
		return int32(compare_spells(arg1, arg2))
	})
}
func what_does_guild_know(guild_nr int, ch *char_data) {
	var (
		overflow *byte = libc.CString("\r\n**OVERFLOW**\r\n")
		buf2     [64936]byte
		i        int
		sortpos  int
		canknow  int
		j        int
		k        int
		count    int    = 0
		cost     int    = 0
		nlen     uint64 = 0
		len_     uint64 = 0
	)
	len_ = uint64(stdio.Snprintf(&buf2[0], int(64936), "You have %d practice session%s remaining.\r\n", ch.Player_specials.Class_skill_points[ch.Chclass], func() string {
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) == 1 {
			return ""
		}
		return "s"
	}()))
	nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "You can practice these skills:     @CKnown@D/@cPrac. Max    @GPS Cost@n\r\n@w-------------------------------------------------------------@n\r\n"))
	len_ += nlen
	for sortpos = 0; sortpos < SKILL_TABLE_SIZE; sortpos++ {
		i = sortpos
		if does_guild_know(guild_nr, i) != 0 && skill_type(i) == (1<<1) {
			for func() int {
				canknow = 0
				k = 0
				return func() int {
					j = 0
					return j
				}()
			}(); j <= NUM_CLASSES; j++ {
				if ((ch.Chclasses[j])+(ch.Epicclasses[j])) > 0 && int(spell_info[i].Can_learn_skill[j]) > SKLEARN_CANT {
					k = int(spell_info[i].Can_learn_skill[j])
				}
			}
			canknow = highest_skill_value(GET_LEVEL(ch), k)
			count++
			cost = calculate_skill_cost(ch, i)
			if k == SKLEARN_CLASS {
				if int(ch.Skills[i]) < canknow {
					nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "@y(@Y%2d@y) @W%-30s @y(@Y%2d@y) @C%d@D/@c%3d        @g%d%s@n\r\n", count, spell_info[i].Name, count, ch.Skills[i], canknow, cost, func() string {
						if int(ch.Skills[i]) > 100 {
							return "  @D(@YGrand Master@D)@n"
						}
						return ""
					}()))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				} else {
					nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "@y(@Y%2d@y) @W%-30s @y(@Y%2d@y) @C%d@D/@c%3d        @g%d%s@n\r\n", count, spell_info[i].Name, count, ch.Skills[i], canknow, cost, func() string {
						if int(ch.Skills[i]) > 100 {
							return "  @D(@YGrand Master@D)@n"
						}
						return ""
					}()))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				}
			} else {
				nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "@y(@Y%2d@y) @W%-30s @y(@Y%2d@y) @C%d@D/@c%3d        @g%d%s@n\r\n", count, spell_info[i].Name, count, ch.Skills[i], canknow, cost, func() string {
					if int(ch.Skills[i]) > 100 {
						return "  @D(@YGrand Master@D)@n"
					}
					return ""
				}()))
				if len_+nlen >= uint64(64936) || nlen < 0 {
					break
				}
				len_ += nlen
			}
		}
	}
	for sortpos = 1; sortpos <= NUM_FEATS_DEFINED; sortpos++ {
		i = feat_sort_info[sortpos]
		if does_guild_know_feat(guild_nr, i) != 0 && feat_is_available(ch, i, 0, nil) != 0 && feat_list[i].In_game != 0 && feat_list[i].Can_learn != 0 {
			nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "@b%-20s@n\r\n", feat_list[i].Name))
			if len_+nlen >= uint64(64936) || nlen < 0 {
				break
			}
			len_ += nlen
		}
	}
	if config_info.Play.Enable_languages != 0 {
		len_ += uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "\r\nand the following languages:\r\n"))
		for sortpos = 0; sortpos < SKILL_TABLE_SIZE; sortpos++ {
			i = sortpos
			if does_guild_know(guild_nr, i) != 0 && IS_SET(bitvector_t(int32(skill_type(i))), 1<<2) {
				for func() int {
					canknow = 0
					return func() int {
						j = 0
						return j
					}()
				}(); j < NUM_CLASSES; j++ {
					if int(spell_info[i].Can_learn_skill[j]) > canknow {
						canknow = int(spell_info[i].Can_learn_skill[j])
					}
				}
				canknow = highest_skill_value(GET_LEVEL(ch), canknow)
				if int(ch.Skills[i]) < canknow {
					nlen = uint64(stdio.Snprintf(&buf2[len_], int(64936-uintptr(len_)), "%-20s %s\r\n", spell_info[i].Name, func() string {
						if int(ch.Skills[i]) != 0 {
							return "known"
						}
						return "unknown"
					}()))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				}
			}
		}
	}
	if len_ >= uint64(64936) {
		libc.StrCpy((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&buf2[64935]), -libc.StrLen(overflow)))), -1)), overflow)
	}
	page_string(ch.Desc, &buf2[0], TRUE)
}
func prereq_pass(ch *char_data, snum int) int {
	if snum == SKILL_KOUSENGAN || snum == SKILL_TSUIHIDAN || snum == SKILL_RENZO || snum == SKILL_SHOGEKIHA {
		if int(ch.Skills[SKILL_KIBALL]) < 40 || int(ch.Skills[SKILL_KIBLAST]) < 40 || int(ch.Skills[SKILL_BEAM]) < 40 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Kiball, Kiblast, and Beam to Skill LVL 40."))
			return 0
		}
	} else if snum == SKILL_INSTANTT {
		if int(ch.Skills[SKILL_FOCUS]) < 90 || int(ch.Skills[SKILL_CONCENTRATION]) < 90 || int(ch.Skills[SKILL_ZANZOKEN]) < 90 {
			send_to_char(ch, libc.CString("You can not train instant transmission until you have Focus, Concentration, and Zanzoken up to Skill LVL 90."))
			return 0
		}
	} else if snum == SKILL_SLAM {
		if int(ch.Skills[SKILL_UPPERCUT]) < 50 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained uppercut to Skill LVL 50."))
			return 0
		}
	} else if snum == SKILL_UPPERCUT {
		if int(ch.Skills[SKILL_ELBOW]) < 40 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained elbow to Skill LVL 40."))
			return 0
		}
	} else if snum == SKILL_HEELDROP {
		if int(ch.Skills[SKILL_ROUNDHOUSE]) < 50 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained roundhouse to Skill LVL 50."))
			return 0
		}
	} else if snum == SKILL_ROUNDHOUSE {
		if int(ch.Skills[SKILL_KNEE]) < 40 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained knee to Skill LVL 40."))
			return 0
		}
	} else if snum == SKILL_KIBALL || snum == SKILL_KIBLAST || snum == SKILL_BEAM {
		if int(ch.Skills[SKILL_FOCUS]) < 30 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained focus to Skill LVL 30."))
			return 0
		}
	} else if IS_SET(bitvector_t(int32(spell_info[snum].Flags)), 1<<10) || IS_SET(bitvector_t(int32(spell_info[snum].Flags)), 1<<11) {
		if snum != 530 && snum != 531 {
			if int(ch.Skills[SKILL_TSUIHIDAN]) < 40 || int(ch.Skills[SKILL_RENZO]) < 40 || int(ch.Skills[SKILL_SHOGEKIHA]) < 40 {
				send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Tsuihidan, Renzokou Energy Dan, and Shogekiha to Skill LVL 40."))
				return 0
			}
		}
	} else if IS_SET(bitvector_t(int32(spell_info[snum].Flags)), 1<<12) {
		if int(ch.Chclass) == CLASS_ROSHI && (int(ch.Skills[SKILL_KAMEHAMEHA]) < 40 || int(ch.Skills[SKILL_KIENZAN]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Kamehameha and Kienzan to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_TSUNA && (int(ch.Skills[SKILL_WRAZOR]) < 40 || int(ch.Skills[SKILL_WSPIKE]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Water Razor and Water Spikes to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_PICCOLO && (int(ch.Skills[SKILL_MASENKO]) < 40 || int(ch.Skills[SKILL_SBC]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Masenko and Special Beam Cannon to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_FRIEZA && (int(ch.Skills[SKILL_DEATHBEAM]) < 40 || int(ch.Skills[SKILL_KIENZAN]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Deathbeam and Kienzan to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_GINYU && (int(ch.Skills[SKILL_CRUSHER]) < 40 || int(ch.Skills[SKILL_ERASER]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Crusher Ball and Eraser Cannon to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_BARDOCK && (int(ch.Skills[SKILL_GALIKGUN]) < 40 || int(ch.Skills[SKILL_FINALFLASH]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Galik Gun and Final Flash to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_TAPION && (int(ch.Skills[SKILL_TSLASH]) < 40 || int(ch.Skills[SKILL_DDSLASH]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Twin Slash and Darkness Dragon Slash to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_NAIL && (int(ch.Skills[SKILL_MASENKO]) < 40 || int(ch.Skills[SKILL_KOUSENGAN]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Masenko and Kousengan to Skill LVL 40."))
			return 0
		}
		if int(ch.Race) == RACE_ANDROID && (int(ch.Skills[SKILL_DUALBEAM]) < 40 || int(ch.Skills[SKILL_HELLFLASH]) < 40) {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Dual Beam and Hell Flash to Skill LVL 40."))
			return 0
		}
		if int(ch.Chclass) == CLASS_JINTO && int(ch.Skills[SKILL_BREAKER]) < 40 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained Star Breaker to Skill LVL 40."))
			return 0
		}
	} else if IS_SET(bitvector_t(int32(spell_info[snum].Flags)), 1<<13) {
		if int(ch.Skills[SKILL_FOCUS]) < 60 || int(ch.Skills[SKILL_CONCENTRATION]) < 80 {
			send_to_char(ch, libc.CString("You can not train that skill until you at least have trained focus to Skill LVL 60 and concentration to Skill LVL 80."))
			return 0
		}
	}
	return 1
}
func handle_forget(keeper *char_data, guild_nr int, ch *char_data, argument *byte, unused int) {
	var skill_num int
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("What skill do you want to start to forget?\r\n"))
		return
	}
	skill_num = find_skill_num(argument, 1<<1)
	if int(ch.Skills[skill_num]) > 30 {
		send_to_char(ch, libc.CString("@MYou can not forget that skill, you know too much about it.@n\r\n"))
		return
	} else if skill_num == SKILL_MIMIC && ch.Mimic > 0 {
		send_to_char(ch, libc.CString("@MYou can not forget mimic while you are using it!\r\n"))
	} else if skill_num == SKILL_FOCUS {
		send_to_char(ch, libc.CString("@MYou can not forget such a fundamental skill!@n\r\n"))
	} else if int(ch.Skills[skill_num]) <= 0 {
		send_to_char(ch, libc.CString("@MYou can not forget a skill you don't know!@n\r\n"))
	} else if ch.Forgeting == skill_num {
		send_to_char(ch, libc.CString("@MYou stop forgetting %s@n\r\n"), spell_info[skill_num].Name)
		ch.Forgetcount = 0
		ch.Forgeting = 0
	} else if ch.Forgeting != 0 {
		send_to_char(ch, libc.CString("@MYou stop forgetting %s, and start trying to forget %s.@n\r\n"), spell_info[ch.Forgeting].Name, spell_info[skill_num].Name)
		ch.Forgetcount = 0
		ch.Forgeting = skill_num
	} else {
		send_to_char(ch, libc.CString("@MYou start trying to forget %s.@n\r\n"), spell_info[skill_num].Name)
		ch.Forgetcount = 0
		ch.Forgeting = skill_num
	}
}
func handle_grand(keeper *char_data, guild_nr int, ch *char_data, argument *byte, unused int) {
	var skill_num int
	skip_spaces(&argument)
	if !CAN_GRAND_MASTER(ch) {
		send_to_char(ch, libc.CString("Your race can not become a Grand Master in a skill through this process.\r\n"))
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("What skill do you want to become a Grand Master in?"))
		return
	}
	skill_num = find_skill_num(argument, 1<<1)
	var buf [64936]byte
	if does_guild_know(guild_nr, skill_num) == 0 {
		stdio.Snprintf(&buf[0], int(64936), libc.GoString(guild_index[guild_nr].No_such_skill), GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if int(ch.Skills[skill_num]) <= 0 {
		send_to_char(ch, libc.CString("You do not know that skill!\r\n"))
		return
	} else if int(ch.Skills[skill_num]) < 100 {
		send_to_char(ch, libc.CString("You haven't even mastered that skill. How can you become a Grand Master in it?\r\n"))
		return
	} else if int(ch.Skills[skill_num]) >= 103 {
		send_to_char(ch, libc.CString("You have already become a Grand Master in that skill and have progessed as far as possible in it.\r\n"))
		return
	} else if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1000 {
		send_to_char(ch, libc.CString("You need at least 1,000 practice sessions to rank up beyond 100 in a skill.\r\n"))
		return
	} else {
		if int(ch.Skills[skill_num]) == 100 {
			send_to_char(ch, libc.CString("@YYou have ascended to Grand Master in the skill, @C%s@Y.\r\n"), spell_info[skill_num].Name)
		} else {
			send_to_char(ch, libc.CString("@YYou have ranked up in your Grand Mastery of the skill, @C%s@Y.\r\n"), spell_info[skill_num].Name)
		}
		for {
			ch.Skills[skill_num] = int8(int(ch.Skills[skill_num]) + 1)
			if true {
				break
			}
		}
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1000
	}
}
func handle_practice(keeper *char_data, guild_nr int, ch *char_data, argument *byte, skill int) {
	var percent int = GET_SKILL(ch, skill)
	_ = percent
	var skill_num int
	var learntype int
	var pointcost int
	var highest int
	var i int
	var buf [64936]byte
	skip_spaces(&argument)
	if *argument == 0 {
		what_does_guild_know(guild_nr, ch)
		return
	}
	if (ch.Player_specials.Class_skill_points[ch.Chclass]) <= 0 {
		send_to_char(ch, libc.CString("You do not seem to be able to practice now.\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SHOCKED) {
		send_to_char(ch, libc.CString("You can not practice while your mind is shocked!\r\n"))
		return
	}
	skill_num = find_skill_num(argument, 1<<1)
	if libc.StrStr(sensei_style[ch.Chclass], argument) != nil {
		skill_num = 539
	}
	if skill_num == ch.Forgeting {
		send_to_char(ch, libc.CString("You can't practice that! You are trying to forget it!@n\r\n"))
		return
	}
	if does_guild_know(guild_nr, skill_num) == 0 {
		stdio.Snprintf(&buf[0], int(64936), libc.GoString(guild_index[guild_nr].No_such_skill), GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if IS_SET(bitvector_t(int32(spell_info[skill_num].Skilltype)), 1<<1) {
		for func() int {
			learntype = 0
			return func() int {
				i = 0
				return i
			}()
		}(); i < NUM_CLASSES; i++ {
			if int(spell_info[skill_num].Can_learn_skill[i]) > learntype {
				learntype = int(spell_info[skill_num].Can_learn_skill[i])
			}
		}
		switch learntype {
		case SKLEARN_CANT:
			stdio.Snprintf(&buf[0], int(64936), libc.GoString(guild_index[guild_nr].No_such_skill), GET_NAME(ch))
			do_tell(keeper, &buf[0], cmd_tell, 0)
			return
		case SKLEARN_CROSSCLASS:
			highest = highest_skill_value(GET_LEVEL(ch), learntype)
		case SKLEARN_CLASS:
			highest = highest_skill_value(GET_LEVEL(ch), learntype)
		default:
			basic_mud_log(libc.CString("Unknown SKLEARN type for skill %d in practice"), skill_num)
			send_to_char(ch, libc.CString("You can't learn that.\r\n"))
			return
		}
		pointcost = calculate_skill_cost(ch, skill_num)
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) >= pointcost {
			if prereq_pass(ch, skill_num) == 0 {
				return
			}
			if int(ch.Skills[skill_num]) >= highest {
				send_to_char(ch, libc.CString("You cannot increase that skill again until you progress further.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 75 && (ch.Bonuses[BONUS_MASOCHISTIC]) > 0 {
				if skill_num == SKILL_PARRY || skill_num == SKILL_ZANZOKEN || skill_num == SKILL_DODGE || skill_num == SKILL_BARRIER || skill_num == SKILL_BLOCK || skill_num == SKILL_TSKIN {
					send_to_char(ch, libc.CString("You cannot increase that skill again because it would deny you the pain you enjoy.\r\n"))
					return
				}
			}
			if int(ch.Skills[skill_num]) >= 75 && int(ch.Chclass) == CLASS_TAPION && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 75 && int(ch.Chclass) == CLASS_DABURA && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 75 && int(ch.Chclass) == CLASS_JINTO && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 75 && int(ch.Chclass) == CLASS_TSUNA && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 50 && int(ch.Chclass) == CLASS_FRIEZA && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 50 && int(ch.Chclass) == CLASS_ANDSIX && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 50 && int(ch.Chclass) == CLASS_KURZAK && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 50 && int(ch.Chclass) == CLASS_GINYU && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 50 && int(ch.Chclass) == CLASS_BARDOCK && skill_num == SKILL_SENSE {
				send_to_char(ch, libc.CString("You cannot practice that anymore.\r\n"))
				return
			}
			if int(ch.Skills[skill_num]) >= 100 {
				send_to_char(ch, libc.CString("You know everything about that skill.\r\n"))
				return
			} else {
				if int(ch.Skills[skill_num]) == 0 {
					if slot_count(ch) < ch.Skill_slots {
						if skill_num != 539 {
							send_to_char(ch, libc.CString("You practice and master the basics!\r\n"))
						} else {
							send_to_char(ch, libc.CString("You practice the basics of %s\r\n"), sensei_style[ch.Chclass])
						}
						for {
							ch.Skills[skill_num] = int8(int(ch.Skills[skill_num]) + rand_number(10, 25))
							if true {
								break
							}
						}
						ch.Player_specials.Class_skill_points[ch.Chclass] -= pointcost
						if ch.Forgeting != 0 && int(ch.Skills[ch.Forgeting]) < 30 {
							ch.Forgetcount += 1
							if ch.Forgetcount >= 5 {
								for {
									ch.Skills[ch.Forgeting] = 0
									if true {
										break
									}
								}
								send_to_char(ch, libc.CString("@MYou have finally forgotten what little you knew of %s@n\r\n"), spell_info[ch.Forgeting].Name)
								ch.Forgeting = 0
								ch.Forgetcount = 0
								save_char(ch)
							}
						} else if int(ch.Skills[ch.Forgeting]) < 30 {
							ch.Forgeting = 0
						}
					} else {
						send_to_char(ch, libc.CString("You already know the maximum number of skills you can for the time being!\r\n"))
						return
					}
				} else {
					if skill_num != 539 {
						send_to_char(ch, libc.CString("You practice for a while and manage to advance your technique.\r\n"))
					} else {
						send_to_char(ch, libc.CString("You practice the basics of %s\r\n"), sensei_style[ch.Chclass])
					}
					for {
						ch.Skills[skill_num] = int8(int(ch.Skills[skill_num]) + 1)
						if true {
							break
						}
					}
					ch.Player_specials.Class_skill_points[ch.Chclass] -= pointcost
					if int(ch.Skills[skill_num]) >= 100 {
						send_to_char(ch, libc.CString("You learned a lot by mastering that skill.\r\n"))
						if int(ch.Race) == RACE_KONATSU && skill_num == SKILL_PARRY {
							for {
								ch.Skills[skill_num] = int8(int(ch.Skills[skill_num]) + 5)
								if true {
									break
								}
							}
						}
						gain_exp(ch, int64(level_exp(ch, GET_LEVEL(ch)+1)/20))
					}
					if ch.Forgeting != 0 {
						ch.Forgetcount += 1
						if ch.Forgetcount >= 5 {
							for {
								ch.Skills[ch.Forgeting] = 0
								if true {
									break
								}
							}
							send_to_char(ch, libc.CString("@MYou have finally forgotten what little you knew of %s@n\r\n"), spell_info[ch.Forgeting].Name)
							ch.Forgeting = 0
							ch.Forgetcount = 0
							save_char(ch)
						}
					}
				}
			}
		} else {
			send_to_char(ch, libc.CString("You need %d practice session%s to increase that skill.\r\n"), pointcost, func() string {
				if pointcost == 1 {
					return ""
				}
				return "s"
			}())
		}
	} else {
		stdio.Snprintf(&buf[0], int(64936), libc.GoString(guild_index[guild_nr].No_such_skill), GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
	}
}
func handle_train(keeper *char_data, guild_nr int, ch *char_data, argument *byte) {
	skip_spaces(&argument)
	if argument == nil || *argument == 0 {
		send_to_char(ch, libc.CString("Training sessions remaining: %d\r\nStats: strength constitution agility intelligence wisdom speed\r\n"), ch.Player_specials.Ability_trains)
	} else if ch.Player_specials.Ability_trains == 0 {
		send_to_char(ch, libc.CString("You have no ability training sessions.\r\n"))
	} else if libc.StrNCaseCmp(libc.CString("strength"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Str += 1
		ch.Player_specials.Ability_trains -= 1
	} else if libc.StrNCaseCmp(libc.CString("constitution"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Con += 1
		if (int(ch.Real_abils.Con) % 2) == 0 {
			ch.Max_hit += int64(GET_LEVEL(ch))
		}
		ch.Player_specials.Ability_trains -= 1
	} else if libc.StrNCaseCmp(libc.CString("agility"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Dex += 1
		ch.Player_specials.Ability_trains -= 1
	} else if libc.StrNCaseCmp(libc.CString("intelligence"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Intel += 1
		if (int(ch.Real_abils.Intel) % 2) == 0 {
			ch.Player_specials.Class_skill_points[ch.Chclass] += 1
		}
		ch.Player_specials.Ability_trains -= 1
	} else if libc.StrNCaseCmp(libc.CString("wisdom"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Wis += 1
		ch.Player_specials.Ability_trains -= 1
	} else if libc.StrNCaseCmp(libc.CString("speed"), argument, libc.StrLen(argument)) == 0 {
		send_to_char(ch, config_info.Play.OK)
		ch.Real_abils.Cha += 1
		ch.Player_specials.Ability_trains -= 1
	} else {
		send_to_char(ch, libc.CString("Stats: strength constitution agility intelligence wisdom speed\r\n"))
	}
	affect_total(ch)
	return
}
func handle_gain(keeper *char_data, guild_nr int, ch *char_data, argument *byte, unused int) {
	var whichclass int = int(ch.Chclass)
	skip_spaces(&argument)
	if GET_LEVEL(ch) < 100 && ch.Exp >= int64(level_exp(ch, GET_LEVEL(ch)+1)) {
		if ch.Rp < rpp_to_level(ch) {
			send_to_char(ch, libc.CString("You need at least %d RPP to gain the next level.\r\n"), rpp_to_level(ch))
		} else if rpp_to_level(ch) <= ch.Rp {
			ch.Rp -= rpp_to_level(ch)
			ch.Desc.Rpp = ch.Rp
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
			send_to_char(ch, libc.CString("@D(@cRPP@W: @w-%d@D)@n\n\n"), rpp_to_level(ch))
			gain_level(ch, whichclass)
		} else {
			gain_level(ch, whichclass)
		}
	} else {
		send_to_char(ch, libc.CString("You are not yet ready for further advancement.\r\n"))
	}
	return
}
func rpp_to_level(ch *char_data) int {
	if GET_LEVEL(ch) == 2 {
		if int(ch.Race) == RACE_SAIYAN {
			return 60
		} else if int(ch.Race) == RACE_BIO {
			return 35
		} else if int(ch.Race) == RACE_MAJIN {
			return 55
		} else if int(ch.Race) == RACE_HOSHIJIN {
			return 30
		}
	} else if GET_LEVEL(ch) >= 90 {
		if GET_LEVEL(ch) == 91 {
			return 3
		} else if GET_LEVEL(ch) == 92 {
			return 3
		} else if GET_LEVEL(ch) == 93 {
			return 3
		} else if GET_LEVEL(ch) == 94 {
			return 3
		} else if GET_LEVEL(ch) == 95 {
			return 3
		} else if GET_LEVEL(ch) == 96 {
			return 4
		} else if GET_LEVEL(ch) == 97 {
			return 4
		} else if GET_LEVEL(ch) == 98 {
			return 4
		} else if GET_LEVEL(ch) == 99 {
			return 5
		}
	}
	return 0
}
func handle_exp(keeper *char_data, guild_nr int, ch *char_data, argument *byte) {
	if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 25 {
		send_to_char(ch, libc.CString("You need at least 25 practice sessions to learn.\r\n"))
		return
	}
	if ch.Exp > int64(level_exp(ch, GET_LEVEL(ch)+1)) && GET_LEVEL(ch) != 100 {
		send_to_char(ch, libc.CString("You can't learn with negative TNL.\r\n"))
		return
	} else {
		var amt int64 = int64(level_exp(ch, GET_LEVEL(ch)+1) / 100)
		if GET_LEVEL(ch) == 100 {
			amt = 400000
		}
		act(libc.CString("@c$n@W spends time training you in $s fighting style.@n"), TRUE, keeper, nil, unsafe.Pointer(ch), TO_VICT)
		act(libc.CString("@c$n@W spends time training @C$N@W in $s fighting style.@n"), TRUE, keeper, nil, unsafe.Pointer(ch), TO_NOTVICT)
		send_to_char(ch, libc.CString("@wExperience Gained: @C%s@n\r\n"), add_commas(amt))
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 25
		if int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED {
			amt = int64(float64(amt) + float64(amt)*0.3)
		}
		if int(ch.Race) == RACE_ICER {
			amt = int64(float64(amt) - float64(amt)*0.1)
		}
		gain_exp(ch, amt)
		return
	}
}
func handle_study(keeper *char_data, guild_nr int, ch *char_data, argument *byte, unused int) {
	var (
		expcost    int = 25000
		goldcost   int = 750
		fail       int = FALSE
		reward     int = 25
		goldadjust int = 0
		expadjust  int = 0
	)
	if GET_LEVEL(ch) >= 100 {
		goldadjust = 500
		expadjust = 15000
	} else if GET_LEVEL(ch) >= 91 {
		goldadjust = 450
		expadjust = 12500
	} else if GET_LEVEL(ch) >= 81 {
		goldadjust = 400
		expadjust = 10000
	} else if GET_LEVEL(ch) >= 71 {
		goldadjust = 350
		expadjust = 7500
	} else if GET_LEVEL(ch) >= 61 {
		goldadjust = 300
		expadjust = 5000
	} else if GET_LEVEL(ch) >= 51 {
		goldadjust = 250
		expadjust = 2500
	} else if GET_LEVEL(ch) >= 41 {
		goldadjust = 200
	} else if GET_LEVEL(ch) >= 31 {
		goldadjust = 150
	} else if GET_LEVEL(ch) >= 21 {
		goldadjust = 100
	} else if GET_LEVEL(ch) >= 11 {
		goldadjust = 50
	}
	goldcost += goldadjust
	expcost += expadjust
	if ch.Exp < int64(expcost) {
		send_to_char(ch, libc.CString("You do not have enough experience to study. @D[@wCost@W: @G%s@D]@n\r\n"), add_commas(int64(expcost)))
		fail = TRUE
	}
	if ch.Gold < goldcost {
		send_to_char(ch, libc.CString("You do not have enough zenni to study. @D[@wCost@W: @Y%s@D]@n\r\n"), add_commas(int64(goldcost)))
		fail = TRUE
	}
	if fail == TRUE {
		return
	}
	ch.Exp -= int64(expcost)
	ch.Gold -= goldcost
	ch.Player_specials.Class_skill_points[ch.Chclass] += 25
	act(libc.CString("@c$N@W spends time lecturing you on various subjects.@n"), TRUE, ch, nil, unsafe.Pointer(keeper), TO_CHAR)
	act(libc.CString("@c$N@W spends time lecturing @C$n@W on various subjects.@n"), TRUE, ch, nil, unsafe.Pointer(keeper), TO_ROOM)
	send_to_char(ch, libc.CString("@wYou have gained %d practice sessions in exchange for %s EXP and %s zenni.\r\n"), reward, add_commas(int64(expcost)), add_commas(int64(goldcost)))
}
func handle_learn(keeper *char_data, guild_nr int, ch *char_data, argument *byte) {
	var (
		feat_num int
		subval   int
		sftype   int
		subfeat  int
		ptr      *byte
		cptr     *byte
		buf      [64936]byte
	)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Which feat would you like to learn?\r\n"))
		return
	}
	if ch.Player_specials.Feat_points < 1 {
		send_to_char(ch, libc.CString("You can't learn any new feats right now.\r\n"))
		return
	}
	ptr = libc.StrChr(argument, ':')
	if ptr != nil {
		*ptr = 0
	}
	feat_num = find_feat_num(argument)
	if ptr != nil {
		*ptr = ':'
	}
	if does_guild_know_feat(guild_nr, feat_num) == 0 {
		stdio.Snprintf(&buf[0], int(64936), libc.GoString(guild_index[guild_nr].No_such_skill), GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if int(ch.Feats[feat_num]) != 0 && feat_list[feat_num].Can_stack == 0 {
		send_to_char(ch, libc.CString("You already know the %s feat.\r\n"), feat_list[feat_num].Name)
		return
	}
	if feat_is_available(ch, feat_num, 0, nil) == 0 || feat_list[feat_num].In_game == 0 || feat_list[feat_num].Can_learn == 0 {
		send_to_char(ch, libc.CString("The %s feat is not available to you at this time.\r\n"), argument)
		return
	}
	sftype = 2
	switch feat_num {
	case FEAT_GREATER_WEAPON_SPECIALIZATION:
		fallthrough
	case FEAT_GREATER_WEAPON_FOCUS:
		fallthrough
	case FEAT_WEAPON_SPECIALIZATION:
		fallthrough
	case FEAT_WEAPON_FOCUS:
		fallthrough
	case FEAT_WEAPON_FINESSE:
		fallthrough
	case FEAT_IMPROVED_CRITICAL:
		sftype = 1
		fallthrough
	case FEAT_SPELL_FOCUS:
		fallthrough
	case FEAT_GREATER_SPELL_FOCUS:
		subfeat = feat_to_subfeat(feat_num)
		if subfeat == -1 {
			basic_mud_log(libc.CString("guild: Unconfigured subfeat '%s', check feat_to_subfeat()"), feat_list[feat_num].Name)
			send_to_char(ch, libc.CString("That feat is not yet ready for use.\r\n"))
			return
		}
		if ptr == nil || *ptr == 0 {
			if sftype == 2 {
				cptr = libc.CString("spell school")
			} else {
				cptr = libc.CString("weapon type")
			}
			subfeat = stdio.Snprintf(&buf[0], int(64936), "No ':' found. You must specify a %s to improve. Example:\r\n learn %s: %s\r\nAvailable %s:\r\n", cptr, feat_list[feat_num].Name, func() *byte {
				if sftype == 2 {
					return spell_schools[0]
				}
				return weapon_type[0]
			}(), cptr)
			for subval = 1; subval <= (func() int {
				if sftype == 2 {
					return NUM_SCHOOLS
				}
				return MAX_WEAPON_TYPES
			}()); subval++ {
				if sftype == 2 {
					cptr = spell_schools[subval]
				} else {
					cptr = weapon_type[subval]
				}
				subfeat += stdio.Snprintf(&buf[subfeat], int(64936-uintptr(subfeat)), "  %s\r\n", cptr)
			}
			page_string(ch.Desc, &buf[0], TRUE)
			return
		}
		if *ptr == ':' {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
		skip_spaces(&ptr)
		if ptr == nil || *ptr == 0 {
			if sftype == 2 {
				cptr = libc.CString("spell school")
			} else {
				cptr = libc.CString("weapon type")
			}
			subfeat = stdio.Snprintf(&buf[0], int(64936), "No %s found. You must specify a %s to improve.\r\n\r\nExample:\r\n learn %s: %s\r\n\r\nAvailable %s:\r\n", cptr, cptr, feat_list[feat_num].Name, func() *byte {
				if sftype == 2 {
					return spell_schools[0]
				}
				return weapon_type[0]
			}(), cptr)
			for subval = 1; subval <= (func() int {
				if sftype == 2 {
					return NUM_SCHOOLS
				}
				return MAX_WEAPON_TYPES
			}()); subval++ {
				if sftype == 2 {
					cptr = spell_schools[subval]
				} else {
					cptr = weapon_type[subval]
				}
				subfeat += stdio.Snprintf(&buf[subfeat], int(64936-uintptr(subfeat)), "  %s\r\n", cptr)
			}
			page_string(ch.Desc, &buf[0], TRUE)
			return
		}
		// todo: figure thi sout
		//subval = search_block(ptr, func() [11]*byte {
		//	if sftype == 2 {
		//		return spell_schools
		//	}
		//	return ([11]*byte)(weapon_type)
		//}()[0], FALSE)
		if subval == -1 {
			basic_mud_log(libc.CString("bad subval: %s"), ptr)
			if sftype == 2 {
				ptr = libc.CString("spell school")
			} else {
				ptr = libc.CString("weapon type")
			}
			subfeat = stdio.Snprintf(&buf[0], int(64936), "That is not a known %s. Available %s:\r\n", ptr, ptr)
			for subval = 1; subval <= (func() int {
				if sftype == 2 {
					return NUM_SCHOOLS
				}
				return MAX_WEAPON_TYPES
			}()); subval++ {
				if sftype == 2 {
					cptr = spell_schools[subval]
				} else {
					cptr = weapon_type[subval]
				}
				subfeat += stdio.Snprintf(&buf[subfeat], int(64936-uintptr(subfeat)), "  %s\r\n", cptr)
			}
			page_string(ch.Desc, &buf[0], TRUE)
			return
		}
		if feat_is_available(ch, feat_num, subval, nil) == 0 {
			send_to_char(ch, libc.CString("You do not satisfy the prerequisites for that feat.\r\n"))
			return
		}
		if sftype == 1 {
			if IS_SET_AR(ch.Combat_feats[subfeat][:], bitvector_t(int32(subval))) {
				send_to_char(ch, libc.CString("You already have that weapon feat.\r\n"))
				return
			}
			SET_BIT_AR(ch.Combat_feats[subfeat][:], bitvector_t(int32(subval)))
		} else if sftype == 2 {
			if IS_SET(bitvector_t(int32(ch.School_feats[subfeat])), bitvector_t(int32(subval))) {
				send_to_char(ch, libc.CString("You already have that spell school feat.\r\n"))
				return
			}
			ch.School_feats[subfeat] |= bitvector_t(subval)
		} else {
			basic_mud_log(libc.CString("unknown feat subtype %d in subfeat code"), sftype)
			send_to_char(ch, libc.CString("That feat is not yet ready for use.\r\n"))
			return
		}
		for {
			ch.Feats[feat_num] = int8(int(ch.Feats[feat_num]) + 1)
			if true {
				break
			}
		}
	case FEAT_GREAT_FORTITUDE:
		for {
			ch.Feats[feat_num] = 1
			if true {
				break
			}
		}
		ch.Apply_saving_throw[SAVING_FORTITUDE] += 2
	case FEAT_IRON_WILL:
		for {
			ch.Feats[feat_num] = 1
			if true {
				break
			}
		}
		ch.Apply_saving_throw[SAVING_WILL] += 2
	case FEAT_LIGHTNING_REFLEXES:
		for {
			ch.Feats[feat_num] = 1
			if true {
				break
			}
		}
		ch.Apply_saving_throw[SAVING_REFLEX] += 2
	case FEAT_TOUGHNESS:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		ch.Max_hit += 3
	case FEAT_SKILL_FOCUS:
		if ptr == nil || *ptr == 0 {
			send_to_char(ch, libc.CString("You must specify a skill to improve. Syntax:\r\n  learn skill focus: skill\r\n"))
			return
		}
		if *ptr == ':' {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
		skip_spaces(&ptr)
		if ptr == nil || *ptr == 0 {
			send_to_char(ch, libc.CString("You must specify a skill to improve. Syntax:\r\n  learn skill focus: skill\r\n"))
			return
		}
		if GET_LEVEL(ch) <= 49 {
			send_to_char(ch, libc.CString("You must be at least level 50 to gain this feat on a skill.\r\n"))
			return
		}
		subval = find_skill_num(ptr, 1<<1)
		if subval < 0 {
			send_to_char(ch, libc.CString("I don't recognize that skill.\r\n"))
			return
		}
		var snum int = GET_SKILL(ch, subval)
		if snum > 100 {
			send_to_char(ch, libc.CString("You have already focused that skill as high as possible.\r\n"))
			return
		}
		for {
			ch.Skillmods[subval] = int8(int(ch.Skillmods[subval]) + 5)
			if true {
				break
			}
		}
		for {
			ch.Feats[feat_num] = int8(int(ch.Feats[feat_num]) + 1)
			if true {
				break
			}
		}
	case FEAT_SPELL_MASTERY:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		ch.Player_specials.Spell_mastery_points += int(MAX(1, int64(ability_mod_value(int(ch.Aff_abils.Intel)))))
	case FEAT_ACROBATIC:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_AGILE:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_BALANCE] = int8(int(ch.Skillmods[SKILL_BALANCE]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_ESCAPE_ARTIST] = int8(int(ch.Skillmods[SKILL_ESCAPE_ARTIST]) + 2)
			if true {
				break
			}
		}
	case FEAT_ALERTNESS:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_LISTEN] = int8(int(ch.Skillmods[SKILL_LISTEN]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_SPOT] = int8(int(ch.Skillmods[SKILL_SPOT]) + 2)
			if true {
				break
			}
		}
	case FEAT_ANIMAL_AFFINITY:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_ATHLETIC:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_DECEITFUL:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_DISGUISE] = int8(int(ch.Skillmods[SKILL_DISGUISE]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_FORGERY] = int8(int(ch.Skillmods[SKILL_FORGERY]) + 2)
			if true {
				break
			}
		}
	case FEAT_DEFT_HANDS:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_SLEIGHT_OF_HAND] = int8(int(ch.Skillmods[SKILL_SLEIGHT_OF_HAND]) + 2)
			if true {
				break
			}
		}
	case FEAT_DILIGENT:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_APPRAISE] = int8(int(ch.Skillmods[SKILL_APPRAISE]) + 2)
			if true {
				break
			}
		}
	case FEAT_INVESTIGATOR:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_EAVESDROP] = int8(int(ch.Skillmods[SKILL_EAVESDROP]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_SEARCH] = int8(int(ch.Skillmods[SKILL_SEARCH]) + 2)
			if true {
				break
			}
		}
	case FEAT_MAGICAL_APTITUDE:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_NEGOTIATOR:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_NIMBLE_FINGERS:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_OPEN_LOCK] = int8(int(ch.Skillmods[SKILL_OPEN_LOCK]) + 2)
			if true {
				break
			}
		}
	case FEAT_PERSUASIVE:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
	case FEAT_SELF_SUFFICIENT:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_HEAL] = int8(int(ch.Skillmods[SKILL_HEAL]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_SURVIVAL] = int8(int(ch.Skillmods[SKILL_SURVIVAL]) + 2)
			if true {
				break
			}
		}
	case FEAT_STEALTHY:
		subval = int(ch.Feats[feat_num]) + 1
		for {
			ch.Feats[feat_num] = int8(subval)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_HIDE] = int8(int(ch.Skillmods[SKILL_HIDE]) + 2)
			if true {
				break
			}
		}
		for {
			ch.Skillmods[SKILL_MOVE_SILENTLY] = int8(int(ch.Skillmods[SKILL_MOVE_SILENTLY]) + 2)
			if true {
				break
			}
		}
	default:
		for {
			ch.Feats[feat_num] = TRUE
			if true {
				break
			}
		}
	}
	save_char(ch)
	ch.Player_specials.Feat_points--
	send_to_char(ch, libc.CString("Your training has given you the %s feat!\r\n"), feat_list[feat_num].Name)
	return
}

type GuildCmd struct {
	Cmd  *byte
	Func func(*char_data, int, *char_data, *byte, int)
}

var guild_cmd_tab = []GuildCmd{{Cmd: libc.CString("practice"), Func: handle_practice}, {Cmd: libc.CString("gain"), Func: handle_gain}, {Cmd: libc.CString("forget"), Func: handle_forget}, {Cmd: libc.CString("study"), Func: handle_study}, {Cmd: libc.CString("grand"), Func: handle_grand}}

func guild(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		arg      [2048]byte
		guild_nr int
		i        int
		keeper   *char_data = (*char_data)(me)
	)
	for guild_nr = 0; guild_nr <= top_guild; guild_nr++ {
		if guild_index[guild_nr].Gm == keeper.Nr {
			break
		}
	}
	if guild_nr > top_guild {
		return FALSE
	}
	if guild_index[guild_nr].Func != nil {
		if guild_index[guild_nr].Func(ch, me, cmd, &arg[0]) != 0 {
			return TRUE
		}
	}
	if !AWAKE(keeper) {
		return FALSE
	}
	for i = 0; guild_cmd_tab[i].Cmd != nil; i++ {
		if libc.StrCmp(guild_cmd_tab[i].Cmd, complete_cmd_info[cmd].Command) == 0 {
			break
		}
	}
	if guild_cmd_tab[i].Cmd == nil {
		return FALSE
	}
	if is_guild_ok(keeper, ch, guild_nr) == 0 {
		return TRUE
	}
	guild_cmd_tab[i].Func(keeper, guild_nr, ch, argument, 0)
	return TRUE
}
func clear_skills(gdindex int) {
	var i int
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		guild_index[gdindex].Skills[i] = 0
	}
}
func read_guild_line(gm_f *stdio.File, string_ *byte, data unsafe.Pointer, type_ *byte) {
	var buf [64936]byte
	if get_line(gm_f, &buf[0]) == 0 || stdio.Sscanf(&buf[0], libc.GoString(string_), data) == 0 {
		stdio.Fprintf(stdio.Stderr(), "Error in guild #%d, Could not get %s\n", guild_index[top_guild].Vnum, type_)
		os.Exit(1)
	}
}
func boot_the_guilds(gm_f *stdio.File, filename *byte, rec_count int) {
	var (
		buf  *byte
		buf2 [256]byte
		p    *byte
		buf3 [256]byte
		temp int
		val  int
		t1   int
		t2   int
		rv   int
		done int = FALSE
	)
	stdio.Snprintf(&buf2[0], int(256), "beginning of GM file %s", filename)
	buf = fread_string(gm_f, &buf2[0])
	for done == 0 {
		if *buf == '#' {
			stdio.Sscanf(buf, "#%d\n", &temp)
			stdio.Snprintf(&buf2[0], int(256), "GM #%d in GM file %s", temp, filename)
			libc.Free(unsafe.Pointer(buf))
			top_guild++
			if top_guild == 0 {
				guild_index = make([]guild_data, rec_count)
			}
			guild_index[top_guild].Vnum = room_vnum(temp)
			clear_skills(top_guild)
			get_line(gm_f, &buf3[0])
			rv = stdio.Sscanf(&buf3[0], "%d %d", &t1, &t2)
			for t1 > -1 {
				if rv == 1 {
					guild_index[top_guild].Skills[t1] = 1
				} else if rv == 2 {
					if t2 == 1 {
						guild_index[top_guild].Skills[t1] = 1
					} else if t2 == 2 {
						guild_index[top_guild].Feats[t1] = 1
					} else {
						basic_mud_log(libc.CString("SYSERR: Invalid 2nd arg in guild file!"))
						os.Exit(1)
					}
				} else {
					basic_mud_log(libc.CString("SYSERR: Invalid format in guild file. Expecting 2 args but got %d!"), rv)
					os.Exit(1)
				}
				get_line(gm_f, &buf3[0])
				rv = stdio.Sscanf(&buf3[0], "%d %d", &t1, &t2)
			}
			read_guild_line(gm_f, libc.CString("%f"), unsafe.Pointer(&guild_index[top_guild].Charge), libc.CString("GM_CHARGE"))
			guild_index[top_guild].No_such_skill = fread_string(gm_f, &buf2[0])
			guild_index[top_guild].Not_enough_gold = fread_string(gm_f, &buf2[0])
			read_guild_line(gm_f, libc.CString("%d"), unsafe.Pointer(&guild_index[top_guild].Minlvl), libc.CString("GM_MINLVL"))
			read_guild_line(gm_f, libc.CString("%d"), unsafe.Pointer(&guild_index[top_guild].Gm), libc.CString("GM_TRAINER"))
			guild_index[top_guild].Gm = real_mobile(mob_vnum(guild_index[top_guild].Gm))
			read_guild_line(gm_f, libc.CString("%d"), unsafe.Pointer(&guild_index[top_guild].With_who[0]), libc.CString("GM_WITH_WHO"))
			read_guild_line(gm_f, libc.CString("%d"), unsafe.Pointer(&guild_index[top_guild].Open), libc.CString("GM_OPEN"))
			read_guild_line(gm_f, libc.CString("%d"), unsafe.Pointer(&guild_index[top_guild].Close), libc.CString("GM_CLOSE"))
			guild_index[top_guild].Func = nil
			buf = (*byte)(unsafe.Pointer(&make([]int8, READ_SIZE)[0]))
			get_line(gm_f, buf)
			if buf != nil && *buf != '#' && *buf != '$' {
				p = buf
				for temp = 1; temp < GW_ARRAY_MAX; temp++ {
					if p == nil || *p == 0 {
						break
					}
					if stdio.Sscanf(p, "%d", &val) != 1 {
						basic_mud_log(libc.CString("SYSERR: Can't parse GM_WITH_WHO line in %s: '%s'"), &buf2[0], buf)
						break
					}
					guild_index[top_guild].With_who[temp] = bitvector_t(int32(val))
					for unicode.IsDigit(rune(*p)) || *p == '-' {
						p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
					}
					for *p != 0 && (!unicode.IsDigit(rune(*p)) && *p != '-') {
						p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
					}
				}
				for temp < GW_ARRAY_MAX {
					guild_index[top_guild].With_who[func() int {
						p := &temp
						x := *p
						*p++
						return x
					}()] = 0
				}
				libc.Free(unsafe.Pointer(buf))
				buf = fread_string(gm_f, &buf2[0])
			}
		} else {
			if *buf == '$' {
				done = TRUE
			}
			libc.Free(unsafe.Pointer(buf))
		}
	}
}
func assign_the_guilds() {
	var gdindex int
	cmd_say = find_command(libc.CString("say"))
	cmd_tell = find_command(libc.CString("tell"))
	for gdindex = 0; gdindex <= top_guild; gdindex++ {
		if guild_index[gdindex].Gm == mob_rnum(-1) {
			continue
		}
		if mob_index[guild_index[gdindex].Gm].Func != nil && libc.FuncAddr(mob_index[guild_index[gdindex].Gm].Func) != libc.FuncAddr(guild) {
			guild_index[gdindex].Func = mob_index[guild_index[gdindex].Gm].Func
		}
		mob_index[guild_index[gdindex].Gm].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
				return guild(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	}
}
func guild_customer_string(guild_nr int, detailed int) *byte {
	var (
		gindex int = 0
		flag   int = 0
		nlen   int
		len_   uint64 = 0
		buf    [64936]byte
	)
	for *trade_letters[gindex] != '\n' && len_+1 < uint64(64936) {
		if detailed != 0 {
			if !IS_SET_AR(guild_index[guild_nr].With_who[:], bitvector_t(int32(flag))) {
				nlen = stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), ", %s", trade_letters[gindex])
				if len_+uint64(nlen) >= uint64(64936) || nlen < 0 {
					break
				}
				len_ += uint64(nlen)
			}
		} else {
			buf[func() uint64 {
				p := &len_
				x := *p
				*p++
				return x
			}()] = func() byte {
				if IS_SET_AR(guild_index[guild_nr].With_who[:], bitvector_t(int32(flag))) {
					return '_'
				}
				return *trade_letters[gindex]
			}()
			buf[len_] = '\x00'
			if len_ >= uint64(64936) {
				break
			}
		}
		gindex++
		flag += 1
	}
	buf[64936-1] = '\x00'
	return &buf[0]
}
func list_all_guilds(ch *char_data) {
	var (
		list_all_guilds_header *byte = libc.CString("Virtual   G.Master\tCharge   Members\r\n----------------------------------------------------------------------\r\n")
		gm_nr                  int
		headerlen              int    = libc.StrLen(list_all_guilds_header)
		len_                   uint64 = 0
		buf                    [64936]byte
		buf1                   [16]byte
	)
	buf[0] = '\x00'
	for gm_nr = 0; gm_nr <= top_guild && len_ < uint64(64936); gm_nr++ {
		if (gm_nr % (int(PAGE_LENGTH - 2))) == 0 {
			if len_+uint64(headerlen)+1 >= uint64(64936) {
				break
			}
			libc.StrCpy(&buf[len_], list_all_guilds_header)
			len_ += uint64(headerlen)
		}
		if guild_index[gm_nr].Gm == mob_rnum(-1) {
			libc.StrCpy(&buf1[0], libc.CString("<NONE>"))
		} else {
			stdio.Sprintf(&buf1[0], "%6d", mob_index[guild_index[gm_nr].Gm].Vnum)
		}
		len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%6d\t%s\t\t%5.2f\t%s\r\n", guild_index[gm_nr].Vnum, &buf1[0], guild_index[gm_nr].Charge, guild_customer_string(gm_nr, FALSE)))
	}
	page_string(ch.Desc, &buf[0], TRUE)
}
func list_detailed_guild(ch *char_data, gm_nr int) {
	var (
		i    int
		buf  [64936]byte
		buf1 [64936]byte
		buf2 [64936]byte
	)
	if guild_index[gm_nr].Gm < mob_rnum(-1) {
		libc.StrCpy(&buf1[0], libc.CString("<NONE>"))
	} else {
		stdio.Sprintf(&buf1[0], "%6d   ", mob_index[guild_index[gm_nr].Gm].Vnum)
	}
	stdio.Sprintf(&buf[0], " Guild Master: %s\r\n", &buf1[0])
	stdio.Sprintf(&buf[0], "%s Hours: %4d to %4d,  Surcharge: %5.2f\r\n", &buf[0], guild_index[gm_nr].Open, guild_index[gm_nr].Close, guild_index[gm_nr].Charge)
	stdio.Sprintf(&buf[0], "%s Min Level will train: %d\r\n", &buf[0], guild_index[gm_nr].Minlvl)
	stdio.Sprintf(&buf[0], "%s Whom will train: %s\r\n", &buf[0], guild_customer_string(gm_nr, TRUE))
	stdio.Sprintf(&buf[0], "%s The GM can teach the following:\r\n", &buf[0])
	buf2[0] = '\x00'
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if does_guild_know(gm_nr, i) != 0 {
			stdio.Sprintf(&buf2[0], "%s %s \r\n", &buf2[0], spell_info[i].Name)
		}
	}
	libc.StrCat(&buf[0], &buf2[0])
	page_string(ch.Desc, &buf[0], 1)
}
func show_guild(ch *char_data, arg *byte) {
	var (
		gm_nr  int
		gm_num int
	)
	if *arg == 0 {
		list_all_guilds(ch)
	} else {
		if is_number(arg) != 0 {
			gm_num = libc.Atoi(libc.GoString(arg))
		} else {
			gm_num = -1
		}
		if gm_num > 0 {
			for gm_nr = 0; gm_nr <= top_guild; gm_nr++ {
				if gm_num == int(guild_index[gm_nr].Vnum) {
					break
				}
			}
			if gm_num < 0 || gm_nr > top_guild {
				send_to_char(ch, libc.CString("Illegal guild master number.\n\r"))
				return
			}
			list_detailed_guild(ch, gm_nr)
		}
	}
}
func list_guilds(ch *char_data, rnum zone_rnum, vmin guild_vnum, vmax guild_vnum) {
	var (
		i       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != zone_rnum(-1) {
		bottom = int(zone_table[rnum].Bot)
		top = int(zone_table[rnum].Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("Index VNum    Guild Master\r\n----- ------- ---------------------------------------------\r\n"))
	if top_guild == 0 {
		return
	}
	for i = 0; i <= top_guild; i++ {
		if guild_index[i].Vnum >= room_vnum(bottom) && guild_index[i].Vnum <= room_vnum(top) {
			counter++
			send_to_char(ch, libc.CString("@g%4d@n) [@c%-5d@n]"), counter, guild_index[i].Vnum)
			send_to_char(ch, libc.CString(" @c[@y%d@c]@y %s@n"), func() mob_vnum {
				if guild_index[i].Gm == mob_rnum(-1) {
					return -1
				}
				return mob_index[guild_index[i].Gm].Vnum
			}(), func() string {
				if guild_index[i].Gm == mob_rnum(-1) {
					return ""
				}
				return libc.GoString(mob_proto[guild_index[i].Gm].Short_descr)
			}())
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("None found.\r\n"))
	}
}
func destroy_guilds() {
	var cnt int64
	if guild_index == nil {
		return
	}
	for cnt = 0; cnt <= int64(top_guild); cnt++ {
		if guild_index[cnt].No_such_skill != nil {
			libc.Free(unsafe.Pointer(guild_index[cnt].No_such_skill))
		}
		if guild_index[cnt].Not_enough_gold != nil {
			libc.Free(unsafe.Pointer(guild_index[cnt].Not_enough_gold))
		}
	}
	libc.Free(unsafe.Pointer(&guild_index[0]))
	guild_index = nil
	top_guild = -1
}
func count_guilds(low guild_vnum, high guild_vnum) int {
	var (
		i int
		j int
	)
	for i = func() int {
		j = 0
		return j
	}(); guild_index[i].Vnum <= room_vnum(high) && i <= top_guild; i++ {
		if guild_index[i].Vnum >= room_vnum(low) {
			j++
		}
	}
	return j
}
func levelup_parse(d *descriptor_data, arg *byte) {
}
