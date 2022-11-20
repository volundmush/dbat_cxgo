package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

var feat_list [253]feat_info
var feat_sort_info [751]int
var buf3 [64936]byte
var buf4 [64936]byte

func feato(featnum int, name *byte, in_game int, can_learn int, can_stack int) {
	feat_list[featnum].Name = name
	feat_list[featnum].In_game = in_game
	feat_list[featnum].Can_learn = can_learn
	feat_list[featnum].Can_stack = can_stack
}
func free_feats() {
}
func assign_feats() {
	var i int
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		feat_list[i].Name = libc.CString("Unused Feat")
		feat_list[i].In_game = 0
		feat_list[i].Can_learn = 0
		feat_list[i].Can_stack = 0
	}
	feato(FEAT_ALERTNESS, libc.CString("alertness"), 1, 0, 0)
	feato(FEAT_ARMOR_PROFICIENCY_HEAVY, libc.CString("heavy armor proficiency"), 0, 1, 0)
	feato(FEAT_ARMOR_PROFICIENCY_LIGHT, libc.CString("light armor proficiency"), 0, 1, 0)
	feato(FEAT_ARMOR_PROFICIENCY_MEDIUM, libc.CString("medium armor proficiency"), 0, 1, 0)
	feato(FEAT_BLIND_FIGHT, libc.CString("blind fighting"), 1, 1, 0)
	feato(FEAT_BREW_POTION, libc.CString("brew potion"), 0, 1, 0)
	feato(FEAT_CLEAVE, libc.CString("cleave"), 0, 1, 0)
	feato(FEAT_COMBAT_CASTING, libc.CString("combat casting"), 0, 1, 0)
	feato(FEAT_COMBAT_REFLEXES, libc.CString("combat reflexes"), 0, 1, 0)
	feato(FEAT_CRAFT_MAGICAL_ARMS_AND_ARMOR, libc.CString("craft magical arms and armor"), 0, 1, 0)
	feato(FEAT_CRAFT_ROD, libc.CString("craft rod"), 0, 1, 0)
	feato(FEAT_CRAFT_STAFF, libc.CString("craft staff"), 0, 1, 0)
	feato(FEAT_CRAFT_WAND, libc.CString("craft wand"), 0, 1, 0)
	feato(FEAT_CRAFT_WONDEROUS_ITEM, libc.CString("craft wonderous item"), 0, 1, 0)
	feato(FEAT_DEFLECT_ARROWS, libc.CString("deflect arrows"), 0, 0, 0)
	feato(FEAT_DODGE, libc.CString("dodge"), 1, 1, 0)
	feato(FEAT_EMPOWER_SPELL, libc.CString("empower spell"), 0, 1, 0)
	feato(FEAT_ENDURANCE, libc.CString("endurance"), 0, 1, 0)
	feato(FEAT_ENLARGE_SPELL, libc.CString("enlarge spell"), 0, 0, 0)
	feato(FEAT_WEAPON_PROFICIENCY_BASTARD_SWORD, libc.CString("weapon proficiency - bastard sword"), 0, 1, 0)
	feato(FEAT_EXTEND_SPELL, libc.CString("extend spell"), 0, 1, 0)
	feato(FEAT_EXTRA_TURNING, libc.CString("extra turning"), 0, 1, 0)
	feato(FEAT_FAR_SHOT, libc.CString("far shot"), 0, 0, 0)
	feato(FEAT_FORGE_RING, libc.CString("forge ring"), 0, 1, 0)
	feato(FEAT_GREAT_CLEAVE, libc.CString("great cleave"), 0, 0, 0)
	feato(FEAT_GREAT_FORTITUDE, libc.CString("great fortitude"), 1, 1, 0)
	feato(FEAT_HEIGHTEN_SPELL, libc.CString("heighten spell"), 0, 1, 0)
	feato(FEAT_IMPROVED_BULL_RUSH, libc.CString("improved bull rush"), 0, 0, 0)
	feato(FEAT_IMPROVED_CRITICAL, libc.CString("improved critical"), 1, 1, 1)
	feato(FEAT_IMPROVED_DISARM, libc.CString("improved disarm"), 0, 1, 0)
	feato(FEAT_IMPROVED_INITIATIVE, libc.CString("improved initiative"), 1, 1, 0)
	feato(FEAT_IMPROVED_TRIP, libc.CString("improved trip"), 1, 1, 0)
	feato(FEAT_IMPROVED_TWO_WEAPON_FIGHTING, libc.CString("improved two weapon fighting"), 1, 1, 0)
	feato(FEAT_IMPROVED_UNARMED_STRIKE, libc.CString("improved unarmed strike"), 0, 0, 0)
	feato(FEAT_IRON_WILL, libc.CString("iron will"), 1, 1, 0)
	feato(FEAT_LEADERSHIP, libc.CString("leadership"), 0, 0, 0)
	feato(FEAT_LIGHTNING_REFLEXES, libc.CString("lightning reflexes"), 1, 1, 0)
	feato(FEAT_MARTIAL_WEAPON_PROFICIENCY, libc.CString("martial weapon proficiency"), 0, 1, 0)
	feato(FEAT_MAXIMIZE_SPELL, libc.CString("maximize spell"), 0, 1, 0)
	feato(FEAT_MOBILITY, libc.CString("mobility"), 1, 1, 0)
	feato(FEAT_MOUNTED_ARCHERY, libc.CString("mounted archery"), 0, 0, 0)
	feato(FEAT_MOUNTED_COMBAT, libc.CString("mounted combat"), 0, 0, 0)
	feato(FEAT_POINT_BLANK_SHOT, libc.CString("point blank shot"), 0, 0, 0)
	feato(FEAT_POWER_ATTACK, libc.CString("power attack"), 1, 1, 0)
	feato(FEAT_PRECISE_SHOT, libc.CString("precise shot"), 0, 0, 0)
	feato(FEAT_QUICK_DRAW, libc.CString("quick draw"), 0, 0, 0)
	feato(FEAT_QUICKEN_SPELL, libc.CString("quicken spell"), 0, 1, 0)
	feato(FEAT_RAPID_SHOT, libc.CString("rapid shot"), 0, 0, 0)
	feato(FEAT_RIDE_BY_ATTACK, libc.CString("ride by attack"), 0, 0, 0)
	feato(FEAT_RUN, libc.CString("run"), 0, 0, 0)
	feato(FEAT_SCRIBE_SCROLL, libc.CString("scribe scroll"), 0, 1, 0)
	feato(FEAT_SHOT_ON_THE_RUN, libc.CString("shot on the run"), 0, 0, 0)
	feato(FEAT_SILENT_SPELL, libc.CString("silent spell"), 0, 1, 0)
	feato(FEAT_SIMPLE_WEAPON_PROFICIENCY, libc.CString("simple weapon proficiency"), 1, 1, 0)
	feato(FEAT_SKILL_FOCUS, libc.CString("skill focus"), 1, 1, 1)
	feato(FEAT_SPELL_FOCUS, libc.CString("spell focus"), 0, 1, 1)
	feato(FEAT_SPELL_MASTERY, libc.CString("spell mastery"), 0, 1, 1)
	feato(FEAT_SPELL_PENETRATION, libc.CString("spell penetration"), 0, 1, 0)
	feato(FEAT_SPIRITED_CHARGE, libc.CString("spirited charge"), 0, 0, 0)
	feato(FEAT_SPRING_ATTACK, libc.CString("spring attack"), 1, 0, 0)
	feato(FEAT_STILL_SPELL, libc.CString("still spell"), 0, 1, 0)
	feato(FEAT_STUNNING_FIST, libc.CString("stunning fist"), 0, 1, 0)
	feato(FEAT_SUNDER, libc.CString("sunder"), 0, 1, 0)
	feato(FEAT_TOUGHNESS, libc.CString("toughness"), 1, 1, 1)
	feato(FEAT_TRACK, libc.CString("track"), 0, 1, 0)
	feato(FEAT_TRAMPLE, libc.CString("trample"), 0, 0, 0)
	feato(FEAT_TWO_WEAPON_FIGHTING, libc.CString("two weapon fighting"), 1, 1, 0)
	feato(FEAT_WEAPON_FINESSE, libc.CString("weapon finesse"), 1, 1, 1)
	feato(FEAT_WEAPON_FOCUS, libc.CString("weapon focus"), 0, 1, 1)
	feato(FEAT_WEAPON_SPECIALIZATION, libc.CString("weapon specialization"), 0, 1, 1)
	feato(FEAT_WHIRLWIND_ATTACK, libc.CString("whirlwind attack"), 0, 1, 0)
	feato(FEAT_WEAPON_PROFICIENCY_DRUID, libc.CString("weapon proficiency - druids"), 0, 0, 0)
	feato(FEAT_WEAPON_PROFICIENCY_ROGUE, libc.CString("weapon proficiency - rogues"), 0, 0, 0)
	feato(FEAT_WEAPON_PROFICIENCY_MONK, libc.CString("weapon proficiency - monks"), 0, 0, 0)
	feato(FEAT_WEAPON_PROFICIENCY_WIZARD, libc.CString("weapon proficiency - wizards"), 0, 0, 0)
	feato(FEAT_WEAPON_PROFICIENCY_ELF, libc.CString("weapon proficiency - elves"), 0, 0, 0)
	feato(FEAT_ARMOR_PROFICIENCY_SHIELD, libc.CString("shield armor proficiency"), 0, 0, 0)
	feato(FEAT_SNEAK_ATTACK, libc.CString("sneak attack"), 1, 0, 1)
	feato(FEAT_EVASION, libc.CString("evasion"), 1, 0, 0)
	feato(FEAT_IMPROVED_EVASION, libc.CString("improved evasion"), 1, 0, 0)
	feato(FEAT_ACROBATIC, libc.CString("acrobatic"), 1, 1, 0)
	feato(FEAT_AGILE, libc.CString("agile"), 1, 1, 0)
	feato(FEAT_ALERTNESS, libc.CString("alertness"), 1, 0, 0)
	feato(FEAT_ANIMAL_AFFINITY, libc.CString("animal affinity"), 0, 1, 0)
	feato(FEAT_ATHLETIC, libc.CString("athletic"), 1, 1, 0)
	feato(FEAT_AUGMENT_SUMMONING, libc.CString("augment summoning"), 0, 0, 0)
	feato(FEAT_COMBAT_EXPERTISE, libc.CString("combat expertise"), 0, 0, 0)
	feato(FEAT_DECEITFUL, libc.CString("deceitful"), 1, 1, 0)
	feato(FEAT_DEFT_HANDS, libc.CString("deft hands"), 0, 1, 0)
	feato(FEAT_DIEHARD, libc.CString("diehard"), 1, 0, 0)
	feato(FEAT_DILIGENT, libc.CString("diligent"), 1, 1, 0)
	feato(FEAT_ESCHEW_MATERIALS, libc.CString("eschew materials"), 0, 0, 0)
	feato(FEAT_EXOTIC_WEAPON_PROFICIENCY, libc.CString("exotic weapon proficiency"), 0, 0, 0)
	feato(FEAT_GREATER_SPELL_FOCUS, libc.CString("greater spell focus"), 0, 0, 1)
	feato(FEAT_GREATER_SPELL_PENETRATION, libc.CString("greater spell penetration"), 0, 0, 0)
	feato(FEAT_GREATER_TWO_WEAPON_FIGHTING, libc.CString("greater two weapon fighting"), 1, 0, 0)
	feato(FEAT_GREATER_WEAPON_FOCUS, libc.CString("greater weapon focus"), 0, 1, 1)
	feato(FEAT_GREATER_WEAPON_SPECIALIZATION, libc.CString("greater weapon specialization"), 0, 1, 1)
	feato(FEAT_IMPROVED_COUNTERSPELL, libc.CString("improved counterspell"), 0, 0, 0)
	feato(FEAT_IMPROVED_FAMILIAR, libc.CString("improved familiar"), 0, 0, 0)
	feato(FEAT_IMPROVED_FEINT, libc.CString("improved feint"), 0, 0, 0)
	feato(FEAT_IMPROVED_GRAPPLE, libc.CString("improved grapple"), 0, 0, 0)
	feato(FEAT_IMPROVED_OVERRUN, libc.CString("improved overrun"), 0, 0, 0)
	feato(FEAT_IMPROVED_PRECISE_SHOT, libc.CString("improved precise shot"), 0, 0, 0)
	feato(FEAT_IMPROVED_SHIELD_BASH, libc.CString("improved shield bash"), 0, 0, 0)
	feato(FEAT_IMPROVED_SUNDER, libc.CString("improved sunder"), 0, 0, 0)
	feato(FEAT_IMPROVED_TURNING, libc.CString("improved turning"), 0, 0, 0)
	feato(FEAT_INVESTIGATOR, libc.CString("investigator"), 0, 1, 0)
	feato(FEAT_MAGICAL_APTITUDE, libc.CString("magical aptitude"), 0, 1, 0)
	feato(FEAT_MANYSHOT, libc.CString("manyshot"), 0, 0, 0)
	feato(FEAT_NATURAL_SPELL, libc.CString("natural spell"), 0, 0, 0)
	feato(FEAT_NEGOTIATOR, libc.CString("negotiator"), 0, 1, 0)
	feato(FEAT_NIMBLE_FINGERS, libc.CString("nimble fingers"), 0, 1, 0)
	feato(FEAT_PERSUASIVE, libc.CString("persuasive"), 0, 1, 0)
	feato(FEAT_RAPID_RELOAD, libc.CString("rapid reload"), 0, 0, 0)
	feato(FEAT_SELF_SUFFICIENT, libc.CString("self sufficient"), 0, 1, 0)
	feato(FEAT_STEALTHY, libc.CString("stealthy"), 1, 1, 0)
	feato(FEAT_ARMOR_PROFICIENCY_TOWER_SHIELD, libc.CString("tower shield armor proficiency"), 0, 0, 0)
	feato(FEAT_TWO_WEAPON_DEFENSE, libc.CString("two weapon defense"), 0, 0, 0)
	feato(FEAT_WIDEN_SPELL, libc.CString("widen spell"), 0, 0, 0)
}
func feat_is_available(ch *char_data, featnum int, iarg int, sarg *byte) bool {
	if featnum > NUM_FEATS_DEFINED {
		return false
	}
	if int(ch.Feats[featnum]) != 0 && feat_list[featnum].Can_stack == 0 {
		return false
	}
	switch featnum {
	case FEAT_ARMOR_PROFICIENCY_HEAVY:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_MEDIUM]) != 0 {
			return true
		}
		return false
	case FEAT_ARMOR_PROFICIENCY_MEDIUM:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_LIGHT]) != 0 {
			return true
		}
		return false
	case FEAT_AUGMENT_SUMMONING:
		if int(ch.Feats[FEAT_SPELL_FOCUS]) != 0 {
			return true
		}
		return false
	case FEAT_IMPROVED_SHIELD_BASH:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_SHIELD]) != 0 {
			return true
		}
		return false
	case FEAT_DODGE:
		if int(ch.Aff_abils.Dex) >= 13 {
			return true
		}
		return false
	case FEAT_COMBAT_EXPERTISE:
		if int(ch.Aff_abils.Intel) >= 13 {
			return true
		}
		return false
	case FEAT_DIEHARD:
		if int(ch.Feats[FEAT_ENDURANCE]) != 0 {
			return true
		}
		return false
	case FEAT_SUNDER:
		fallthrough
	case FEAT_IMPROVED_OVERRUN:
		fallthrough
	case FEAT_IMPROVED_BULL_RUSH:
		if int(ch.Feats[FEAT_POWER_ATTACK]) != 0 {
			return true
		}
		return false
	case FEAT_MOBILITY:
		if int(ch.Feats[FEAT_DODGE]) != 0 {
			return true
		}
		return false
	case FEAT_EXOTIC_WEAPON_PROFICIENCY:
		if ch.Accuracy >= 1 {
			return true
		}
		return false
	case FEAT_WEAPON_PROFICIENCY_BASTARD_SWORD:
		if ch.Accuracy >= 1 {
			return true
		}
		return false
	case FEAT_IMPROVED_FEINT:
		fallthrough
	case FEAT_IMPROVED_DISARM:
		fallthrough
	case FEAT_IMPROVED_TRIP:
		if int(ch.Feats[FEAT_COMBAT_EXPERTISE]) != 0 {
			return true
		}
		return false
	case FEAT_IMPROVED_GRAPPLE:
		fallthrough
	case FEAT_DEFLECT_ARROWS:
		if int(ch.Feats[FEAT_IMPROVED_UNARMED_STRIKE]) != 0 {
			if int(ch.Aff_abils.Dex) >= 13 {
				return true
			}
		}
		return false
	case FEAT_STUNNING_FIST:
		if int(ch.Aff_abils.Dex) >= 13 {
			if int(ch.Aff_abils.Wis) >= 13 {
				if int(ch.Feats[FEAT_IMPROVED_UNARMED_STRIKE]) != 0 {
					if ch.Accuracy >= 8 {
						return true
					}
				}
			}
		}
		return false
	case FEAT_POWER_ATTACK:
		if int(ch.Aff_abils.Str) >= 13 {
			return true
		}
		return false
	case FEAT_IMPROVED_PRECISE_SHOT:
		if int(ch.Aff_abils.Dex) >= 19 {
			if int(ch.Feats[FEAT_POINT_BLANK_SHOT]) != 0 {
				if ch.Accuracy >= 11 {
					return true
				}
			}
		}
		return false
	case FEAT_CLEAVE:
		if int(ch.Feats[FEAT_POWER_ATTACK]) != 0 {
			return true
		}
		return false
	case FEAT_TWO_WEAPON_FIGHTING:
		if int(ch.Aff_abils.Dex) >= 15 {
			return true
		}
		return false
	case FEAT_IMPROVED_TWO_WEAPON_FIGHTING:
		if int(ch.Aff_abils.Dex) >= 17 && int(ch.Feats[FEAT_TWO_WEAPON_FIGHTING]) != 0 && ch.Accuracy >= 6 {
			return true
		}
		return false
	case FEAT_IMPROVED_CRITICAL:
		if ch.Accuracy < 8 {
			return false
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) {
			return true
		}
		return false
	case FEAT_FAR_SHOT:
		if int(ch.Feats[FEAT_POINT_BLANK_SHOT]) != 0 {
			return true
		}
		return false
	case FEAT_WEAPON_FINESSE:
		fallthrough
	case FEAT_WEAPON_FOCUS:
		if ch.Accuracy < 1 {
			return false
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) {
			return true
		}
		return false
	case FEAT_WEAPON_SPECIALIZATION:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 4 {
			return false
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) {
			return true
		}
		return false
	case FEAT_GREATER_WEAPON_FOCUS:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 8 {
			return false
		}
		if iarg == 0 {
			return true
		}
		if is_proficient_with_weapon(ch, iarg) && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_FOCUS][:], uint32(int32(iarg))) {
			return true
		}
		return false
	case FEAT_WHIRLWIND_ATTACK:
		if int(ch.Feats[FEAT_COMBAT_EXPERTISE]) != 0 {
			if int(ch.Feats[FEAT_DODGE]) != 0 {
				if int(ch.Feats[FEAT_MOBILITY]) != 0 {
					if int(ch.Feats[FEAT_SPRING_ATTACK]) != 0 {
						if ch.Accuracy >= 4 {
							return true
						}
					}
				}
			}
		}
		return false
	case FEAT_GREATER_WEAPON_SPECIALIZATION:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 12 {
			return false
		}
		if iarg == 0 {
			return true
		}
		if is_proficient_with_weapon(ch, iarg) && IS_SET_AR(ch.Combat_feats[CFEAT_GREATER_WEAPON_FOCUS][:], uint32(int32(iarg))) && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_SPECIALIZATION][:], uint32(int32(iarg))) && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_FOCUS][:], uint32(int32(iarg))) {
			return true
		}
		return false
	case FEAT_SPELL_FOCUS:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_SPELL_PENETRATION:
		if GET_LEVEL(ch) != 0 {
			return true
		}
		return false
	case FEAT_BREW_POTION:
		if GET_LEVEL(ch) >= 3 {
			return true
		}
		return false
	case FEAT_CRAFT_MAGICAL_ARMS_AND_ARMOR:
		if GET_LEVEL(ch) >= 5 {
			return true
		}
		return false
	case FEAT_CRAFT_ROD:
		if GET_LEVEL(ch) >= 9 {
			return true
		}
		return false
	case FEAT_CRAFT_STAFF:
		if GET_LEVEL(ch) >= 12 {
			return true
		}
		return false
	case FEAT_CRAFT_WAND:
		if GET_LEVEL(ch) >= 5 {
			return true
		}
		return false
	case FEAT_FORGE_RING:
		if GET_LEVEL(ch) >= 5 {
			return true
		}
		return false
	case FEAT_SCRIBE_SCROLL:
		if GET_LEVEL(ch) >= 1 {
			return true
		}
		return false
	case FEAT_EMPOWER_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_EXTEND_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_HEIGHTEN_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_MAXIMIZE_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_QUICKEN_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_SILENT_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_STILL_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	case FEAT_EXTRA_TURNING:
		if ((ch.Chclasses[CLASS_PICCOLO]) + (ch.Epicclasses[CLASS_PICCOLO])) != 0 {
			return true
		}
		return false
	case FEAT_SPELL_MASTERY:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return true
		}
		return false
	default:
		return true
	}
}
func is_proficient_with_armor(ch *char_data, cmarmor_type int) bool {
	switch cmarmor_type {
	case ARMOR_TYPE_LIGHT:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_LIGHT]) != 0 {
			return true
		}
	case ARMOR_TYPE_MEDIUM:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_MEDIUM]) != 0 {
			return true
		}
	case ARMOR_TYPE_HEAVY:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_HEAVY]) != 0 {
			return true
		}
	case ARMOR_TYPE_SHIELD:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_SHIELD]) != 0 {
			return true
		}
	}
	return false
}
func is_proficient_with_weapon(ch *char_data, cmweapon_type int) bool {
	switch cmweapon_type {
	case WEAPON_TYPE_UNARMED:
		return true
	case WEAPON_TYPE_DAGGER:
		fallthrough
	case WEAPON_TYPE_MACE:
		fallthrough
	case WEAPON_TYPE_SICKLE:
		fallthrough
	case WEAPON_TYPE_SPEAR:
		fallthrough
	case WEAPON_TYPE_STAFF:
		fallthrough
	case WEAPON_TYPE_CROSSBOW:
		fallthrough
	case WEAPON_TYPE_SLING:
		fallthrough
	case WEAPON_TYPE_THROWN:
		fallthrough
	case WEAPON_TYPE_CLUB:
		if int(ch.Feats[FEAT_SIMPLE_WEAPON_PROFICIENCY]) != 0 {
			return true
		}
	case WEAPON_TYPE_SHORTBOW:
		fallthrough
	case WEAPON_TYPE_LONGBOW:
		fallthrough
	case WEAPON_TYPE_HAMMER:
		fallthrough
	case WEAPON_TYPE_LANCE:
		fallthrough
	case WEAPON_TYPE_FLAIL:
		fallthrough
	case WEAPON_TYPE_LONGSWORD:
		fallthrough
	case WEAPON_TYPE_SHORTSWORD:
		fallthrough
	case WEAPON_TYPE_GREATSWORD:
		fallthrough
	case WEAPON_TYPE_RAPIER:
		fallthrough
	case WEAPON_TYPE_SCIMITAR:
		fallthrough
	case WEAPON_TYPE_POLEARM:
		fallthrough
	case WEAPON_TYPE_BASTARD_SWORD:
		fallthrough
	case WEAPON_TYPE_AXE:
		if int(ch.Feats[FEAT_MARTIAL_WEAPON_PROFICIENCY]) != 0 {
			return true
		}
	default:
		return false
	}
	return false
}
func compare_feats(x unsafe.Pointer, y unsafe.Pointer) int {
	var (
		a int = *(*int)(x)
		b int = *(*int)(y)
	)
	return libc.StrCmp(feat_list[a].Name, feat_list[b].Name)
}
func sort_feats() {
	var a int
	for a = 1; a <= NUM_FEATS_DEFINED; a++ {
		feat_sort_info[a] = a
	}
	libc.Sort(unsafe.Pointer(&feat_sort_info[1]), NUM_FEATS_DEFINED, uint32(unsafe.Sizeof(int(0))), func(arg1 unsafe.Pointer, arg2 unsafe.Pointer) int32 {
		return int32(compare_feats(arg1, arg2))
	})
}
func list_feats_known(ch *char_data) {
	var (
		i          int
		j          int
		sortpos    int
		none_shown int = 1
		temp_value int
		added_hp   int = 0
		buf        [64936]byte
		buf2       [64936]byte
	)
	if ch.Player_specials.Feat_points == 0 {
		libc.StrCpy(&buf[0], libc.CString("\r\nYou cannot learn any feats right now.\r\n"))
	} else {
		stdio.Sprintf(&buf[0], "\r\nYou can learn %d feat%s right now.\r\n", ch.Player_specials.Feat_points, func() string {
			if ch.Player_specials.Feat_points == 1 {
				return ""
			}
			return "s"
		}())
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@WFeats Known@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	libc.StrCpy(&buf2[0], &buf[0])
	for sortpos = 1; sortpos <= NUM_FEATS_DEFINED; sortpos++ {
		if libc.StrLen(&buf2[0]) > int(MAX_STRING_LENGTH-32) {
			break
		}
		i = feat_sort_info[sortpos]
		if int(ch.Feats[i]) != 0 && feat_list[i].In_game != 0 {
			switch i {
			case FEAT_SKILL_FOCUS:
				stdio.Sprintf(&buf[0], "%-20s (+%d points overall)\r\n", feat_list[i].Name, int(ch.Feats[i])*2)
				libc.StrCat(&buf2[0], &buf[0])
				none_shown = 0
			case FEAT_TOUGHNESS:
				temp_value = int(ch.Feats[FEAT_TOUGHNESS])
				added_hp = temp_value * 3
				stdio.Sprintf(&buf[0], "%-20s (+%d hp)\r\n", feat_list[i].Name, added_hp)
				libc.StrCat(&buf2[0], &buf[0])
				none_shown = 0
			case FEAT_IMPROVED_CRITICAL:
				fallthrough
			case FEAT_WEAPON_FINESSE:
				fallthrough
			case FEAT_WEAPON_FOCUS:
				fallthrough
			case FEAT_WEAPON_SPECIALIZATION:
				fallthrough
			case FEAT_GREATER_WEAPON_FOCUS:
				fallthrough
			case FEAT_GREATER_WEAPON_SPECIALIZATION:
				for j = 0; j <= MAX_WEAPON_TYPES; j++ {
					if IS_SET_AR(ch.Combat_feats[feat_to_subfeat(i)][:], uint32(int32(j))) {
						stdio.Sprintf(&buf[0], "%-20s (%s)\r\n", feat_list[i].Name, weapon_type[j])
						libc.StrCat(&buf2[0], &buf[0])
						none_shown = 0
					}
				}
			default:
				stdio.Sprintf(&buf[0], "%-20s\r\n", feat_list[i].Name)
				libc.StrCat(&buf2[0], &buf[0])
				none_shown = 0
			}
		}
	}
	if none_shown != 0 {
		stdio.Sprintf(&buf[0], "You do not know any feats at this time.\r\n")
		libc.StrCat(&buf2[0], &buf[0])
	}
	page_string(ch.Desc, &buf2[0], 1)
}
func list_feats_available(ch *char_data) {
	var (
		buf        [64936]byte
		buf2       [64936]byte
		i          int
		sortpos    int
		none_shown int = 1
	)
	if ch.Player_specials.Feat_points == 0 {
		libc.StrCpy(&buf[0], libc.CString("\r\nYou cannot learn any feats right now.\r\n"))
	} else {
		stdio.Sprintf(&buf[0], "\r\nYou can learn %d feat%s right now.\r\n", ch.Player_specials.Feat_points, func() string {
			if ch.Player_specials.Feat_points == 1 {
				return ""
			}
			return "s"
		}())
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@WFeats Available to Learn@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	libc.StrCpy(&buf2[0], &buf[0])
	for sortpos = 1; sortpos <= NUM_FEATS_DEFINED; sortpos++ {
		i = feat_sort_info[sortpos]
		if libc.StrLen(&buf2[0]) >= int(MAX_STRING_LENGTH-32) {
			libc.StrCat(&buf2[0], libc.CString("**OVERFLOW**\r\n"))
			break
		}
		if feat_is_available(ch, i, 0, nil) && feat_list[i].In_game != 0 && feat_list[i].Can_learn != 0 {
			stdio.Sprintf(&buf[0], "%-20s\r\n", feat_list[i].Name)
			libc.StrCat(&buf2[0], &buf[0])
			none_shown = 0
		}
	}
	if none_shown != 0 {
		stdio.Sprintf(&buf[0], "There are no feats available for you to learn at this point.\r\n")
		libc.StrCat(&buf2[0], &buf[0])
	}
	page_string(ch.Desc, &buf2[0], 1)
}
func list_feats_complete(ch *char_data) {
	var (
		buf        [64936]byte
		buf2       [64936]byte
		i          int
		sortpos    int
		none_shown int = 1
	)
	if ch.Player_specials.Feat_points == 0 {
		libc.StrCpy(&buf[0], libc.CString("\r\nYou cannot learn any feats right now.\r\n"))
	} else {
		stdio.Sprintf(&buf[0], "\r\nYou can learn %d feat%s right now.\r\n", ch.Player_specials.Feat_points, func() string {
			if ch.Player_specials.Feat_points == 1 {
				return ""
			}
			return "s"
		}())
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@WComplete Feat List@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@B~@R~@n\r\n")
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n")
	libc.StrCpy(&buf2[0], &buf[0])
	for sortpos = 1; sortpos <= NUM_FEATS_DEFINED; sortpos++ {
		i = feat_sort_info[sortpos]
		if libc.StrLen(&buf2[0]) >= int(MAX_STRING_LENGTH-32) {
			libc.StrCat(&buf2[0], libc.CString("**OVERFLOW**\r\n"))
			break
		}
		if feat_list[i].In_game != 0 {
			stdio.Sprintf(&buf[0], "%-20s\r\n", feat_list[i].Name)
			libc.StrCat(&buf2[0], &buf[0])
			none_shown = 0
		}
	}
	if none_shown != 0 {
		stdio.Sprintf(&buf[0], "There are currently no feats in the game.\r\n")
		libc.StrCat(&buf2[0], &buf[0])
	}
	page_string(ch.Desc, &buf2[0], 1)
}
func find_feat_num(name *byte) int {
	var (
		ftindex int
		ok      int
		temp    *byte
		temp2   *byte
		first   [256]byte
		first2  [256]byte
	)
	for ftindex = 1; ftindex <= NUM_FEATS_DEFINED; ftindex++ {
		if is_abbrev(name, feat_list[ftindex].Name) {
			return ftindex
		}
		ok = 1
		temp = any_one_arg(feat_list[ftindex].Name, &first[0])
		temp2 = any_one_arg(name, &first2[0])
		for first[0] != 0 && first2[0] != 0 && ok != 0 {
			if !is_abbrev(&first2[0], &first[0]) {
				ok = 0
			}
			temp = any_one_arg(temp, &first[0])
			temp2 = any_one_arg(temp2, &first2[0])
		}
		if ok != 0 && first2[0] == 0 {
			return ftindex
		}
	}
	return -1
}
func do_feats(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [80]byte
	one_argument(argument, &arg[0])
	if is_abbrev(&arg[0], libc.CString("known")) || arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax is \"feats <available | complete | known>\".\r\n"))
		list_feats_known(ch)
	} else if is_abbrev(&arg[0], libc.CString("available")) {
		list_feats_available(ch)
	} else if is_abbrev(&arg[0], libc.CString("complete")) {
		list_feats_complete(ch)
	}
}
func feat_to_subfeat(feat int) int {
	switch feat {
	case FEAT_IMPROVED_CRITICAL:
		return CFEAT_IMPROVED_CRITICAL
	case FEAT_WEAPON_FINESSE:
		return CFEAT_WEAPON_FINESSE
	case FEAT_WEAPON_FOCUS:
		return CFEAT_WEAPON_FOCUS
	case FEAT_WEAPON_SPECIALIZATION:
		return CFEAT_WEAPON_SPECIALIZATION
	case FEAT_GREATER_WEAPON_FOCUS:
		return CFEAT_GREATER_WEAPON_FOCUS
	case FEAT_GREATER_WEAPON_SPECIALIZATION:
		return CFEAT_GREATER_WEAPON_SPECIALIZATION
	case FEAT_SPELL_FOCUS:
		return CFEAT_SPELL_FOCUS
	case FEAT_GREATER_SPELL_FOCUS:
		return CFEAT_GREATER_SPELL_FOCUS
	default:
		return -1
	}
}
