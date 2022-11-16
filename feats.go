package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const FEAT_UNDEFINED = 0
const FEAT_ALERTNESS = 1
const FEAT_ARMOR_PROFICIENCY_HEAVY = 3
const FEAT_ARMOR_PROFICIENCY_LIGHT = 4
const FEAT_ARMOR_PROFICIENCY_MEDIUM = 5
const FEAT_BLIND_FIGHT = 6
const FEAT_BREW_POTION = 7
const FEAT_CLEAVE = 8
const FEAT_COMBAT_CASTING = 9
const FEAT_COMBAT_REFLEXES = 10
const FEAT_CRAFT_MAGICAL_ARMS_AND_ARMOR = 11
const FEAT_CRAFT_ROD = 12
const FEAT_CRAFT_STAFF = 13
const FEAT_CRAFT_WAND = 14
const FEAT_CRAFT_WONDEROUS_ITEM = 15
const FEAT_DEFLECT_ARROWS = 16
const FEAT_DODGE = 17
const FEAT_EMPOWER_SPELL = 18
const FEAT_ENDURANCE = 19
const FEAT_ENLARGE_SPELL = 20
const FEAT_WEAPON_PROFICIENCY_BASTARD_SWORD = 21
const FEAT_EXPERTISE = 22
const FEAT_EXTEND_SPELL = 23
const FEAT_EXTRA_TURNING = 24
const FEAT_FAR_SHOT = 25
const FEAT_FORGE_RING = 26
const FEAT_GREAT_CLEAVE = 27
const FEAT_GREAT_FORTITUDE = 28
const FEAT_HEIGHTEN_SPELL = 29
const FEAT_IMPROVED_BULL_RUSH = 30
const FEAT_IMPROVED_CRITICAL = 31
const FEAT_IMPROVED_DISARM = 61
const FEAT_IMPROVED_INITIATIVE = 62
const FEAT_IMPROVED_TRIP = 63
const FEAT_IMPROVED_TWO_WEAPON_FIGHTING = 64
const FEAT_IMPROVED_UNARMED_STRIKE = 65
const FEAT_IRON_WILL = 66
const FEAT_LEADERSHIP = 67
const FEAT_LIGHTNING_REFLEXES = 68
const FEAT_MARTIAL_WEAPON_PROFICIENCY = 69
const FEAT_MAXIMIZE_SPELL = 70
const FEAT_MOBILITY = 71
const FEAT_MOUNTED_ARCHERY = 72
const FEAT_MOUNTED_COMBAT = 73
const FEAT_POINT_BLANK_SHOT = 74
const FEAT_POWER_ATTACK = 75
const FEAT_PRECISE_SHOT = 76
const FEAT_QUICK_DRAW = 77
const FEAT_QUICKEN_SPELL = 78
const FEAT_RAPID_SHOT = 79
const FEAT_RIDE_BY_ATTACK = 80
const FEAT_RUN = 81
const FEAT_SCRIBE_SCROLL = 82
const FEAT_SHIELD_PROFICIENCY = 83
const FEAT_SHOT_ON_THE_RUN = 84
const FEAT_SILENT_SPELL = 85
const FEAT_SIMPLE_WEAPON_PROFICIENCY = 86
const FEAT_SKILL_FOCUS = 87
const FEAT_SPELL_FOCUS = 88
const FEAT_SPELL_MASTERY = 96
const FEAT_SPELL_PENETRATION = 97
const FEAT_SPIRITED_CHARGE = 98
const FEAT_SPRING_ATTACK = 99
const FEAT_STILL_SPELL = 100
const FEAT_STUNNING_FIST = 101
const FEAT_SUNDER = 102
const FEAT_TOUGHNESS = 103
const FEAT_TRACK = 104
const FEAT_TRAMPLE = 105
const FEAT_TWO_WEAPON_FIGHTING = 106
const FEAT_WEAPON_FINESSE = 107
const FEAT_WEAPON_FOCUS = 137
const FEAT_WEAPON_SPECIALIZATION = 167
const FEAT_WHIRLWIND_ATTACK = 197
const FEAT_WEAPON_PROFICIENCY_DRUID = 198
const FEAT_WEAPON_PROFICIENCY_ROGUE = 199
const FEAT_WEAPON_PROFICIENCY_MONK = 200
const FEAT_WEAPON_PROFICIENCY_WIZARD = 201
const FEAT_WEAPON_PROFICIENCY_ELF = 202
const FEAT_ARMOR_PROFICIENCY_SHIELD = 203
const FEAT_SNEAK_ATTACK = 204
const FEAT_EVASION = 205
const FEAT_IMPROVED_EVASION = 206
const FEAT_ACROBATIC = 207
const FEAT_AGILE = 208
const FEAT_ANIMAL_AFFINITY = 209
const FEAT_ATHLETIC = 210
const FEAT_AUGMENT_SUMMONING = 211
const FEAT_COMBAT_EXPERTISE = 212
const FEAT_DECEITFUL = 213
const FEAT_DEFT_HANDS = 214
const FEAT_DIEHARD = 215
const FEAT_DILIGENT = 216
const FEAT_ESCHEW_MATERIALS = 217
const FEAT_EXOTIC_WEAPON_PROFICIENCY = 218
const FEAT_GREATER_SPELL_FOCUS = 219
const FEAT_GREATER_SPELL_PENETRATION = 220
const FEAT_GREATER_TWO_WEAPON_FIGHTING = 221
const FEAT_GREATER_WEAPON_FOCUS = 222
const FEAT_GREATER_WEAPON_SPECIALIZATION = 223
const FEAT_IMPROVED_COUNTERSPELL = 224
const FEAT_IMPROVED_FAMILIAR = 225
const FEAT_IMPROVED_FEINT = 226
const FEAT_IMPROVED_GRAPPLE = 227
const FEAT_IMPROVED_OVERRUN = 228
const FEAT_IMPROVED_PRECISE_SHOT = 229
const FEAT_IMPROVED_SHIELD_BASH = 230
const FEAT_IMPROVED_SUNDER = 231
const FEAT_IMPROVED_TURNING = 232
const FEAT_INVESTIGATOR = 233
const FEAT_MAGICAL_APTITUDE = 234
const FEAT_MANYSHOT = 235
const FEAT_NATURAL_SPELL = 236
const FEAT_NEGOTIATOR = 237
const FEAT_NIMBLE_FINGERS = 238
const FEAT_PERSUASIVE = 239
const FEAT_RAPID_RELOAD = 240
const FEAT_SELF_SUFFICIENT = 241
const FEAT_STEALTHY = 242
const FEAT_ARMOR_PROFICIENCY_TOWER_SHIELD = 243
const FEAT_TWO_WEAPON_DEFENSE = 244
const FEAT_WIDEN_SPELL = 245
const FEAT_CRIPPLING_STRIKE = 246
const FEAT_DEFENSIVE_ROLL = 247
const FEAT_OPPORTUNIST = 248
const FEAT_SKILL_MASTERY = 249
const FEAT_SLIPPERY_MIND = 250
const FEAT_KI_STRIKE = 251
const FEAT_SNATCH_ARROWS = 252
const FEAT_SENSE = 253
const WEAPON_TYPE_UNDEFINED = 0
const WEAPON_TYPE_UNARMED = 1
const WEAPON_TYPE_DAGGER = 2
const WEAPON_TYPE_MACE = 3
const WEAPON_TYPE_SICKLE = 4
const WEAPON_TYPE_SPEAR = 5
const WEAPON_TYPE_STAFF = 6
const WEAPON_TYPE_CROSSBOW = 7
const WEAPON_TYPE_LONGBOW = 8
const WEAPON_TYPE_SHORTBOW = 9
const WEAPON_TYPE_SLING = 10
const WEAPON_TYPE_THROWN = 11
const WEAPON_TYPE_HAMMER = 12
const WEAPON_TYPE_LANCE = 13
const WEAPON_TYPE_FLAIL = 14
const WEAPON_TYPE_LONGSWORD = 15
const WEAPON_TYPE_SHORTSWORD = 16
const WEAPON_TYPE_GREATSWORD = 17
const WEAPON_TYPE_RAPIER = 18
const WEAPON_TYPE_SCIMITAR = 19
const WEAPON_TYPE_POLEARM = 20
const WEAPON_TYPE_CLUB = 21
const WEAPON_TYPE_BASTARD_SWORD = 22
const WEAPON_TYPE_MONK_WEAPON = 23
const WEAPON_TYPE_DOUBLE_WEAPON = 24
const WEAPON_TYPE_AXE = 25
const WEAPON_TYPE_WHIP = 26
const ARMOR_TYPE_UNDEFINED = 0
const ARMOR_TYPE_LIGHT = 1
const ARMOR_TYPE_MEDIUM = 2
const ARMOR_TYPE_HEAVY = 3
const ARMOR_TYPE_SHIELD = 4
const HRANK_ANY = 0
const HRANK_ARCANE = 1
const HRANK_DIVINE = 2
const HRANK_CASTER = 3
const BASE_DC = 10

type feat_info struct {
	Name      *byte
	In_game   int
	Can_learn int
	Can_stack int
}

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
		feat_list[i].In_game = FALSE
		feat_list[i].Can_learn = FALSE
		feat_list[i].Can_stack = FALSE
	}
	feato(FEAT_ALERTNESS, libc.CString("alertness"), TRUE, FALSE, FALSE)
	feato(FEAT_ARMOR_PROFICIENCY_HEAVY, libc.CString("heavy armor proficiency"), FALSE, TRUE, FALSE)
	feato(FEAT_ARMOR_PROFICIENCY_LIGHT, libc.CString("light armor proficiency"), FALSE, TRUE, FALSE)
	feato(FEAT_ARMOR_PROFICIENCY_MEDIUM, libc.CString("medium armor proficiency"), FALSE, TRUE, FALSE)
	feato(FEAT_BLIND_FIGHT, libc.CString("blind fighting"), TRUE, TRUE, FALSE)
	feato(FEAT_BREW_POTION, libc.CString("brew potion"), FALSE, TRUE, FALSE)
	feato(FEAT_CLEAVE, libc.CString("cleave"), FALSE, TRUE, FALSE)
	feato(FEAT_COMBAT_CASTING, libc.CString("combat casting"), FALSE, TRUE, FALSE)
	feato(FEAT_COMBAT_REFLEXES, libc.CString("combat reflexes"), FALSE, TRUE, FALSE)
	feato(FEAT_CRAFT_MAGICAL_ARMS_AND_ARMOR, libc.CString("craft magical arms and armor"), FALSE, TRUE, FALSE)
	feato(FEAT_CRAFT_ROD, libc.CString("craft rod"), FALSE, TRUE, FALSE)
	feato(FEAT_CRAFT_STAFF, libc.CString("craft staff"), FALSE, TRUE, FALSE)
	feato(FEAT_CRAFT_WAND, libc.CString("craft wand"), FALSE, TRUE, FALSE)
	feato(FEAT_CRAFT_WONDEROUS_ITEM, libc.CString("craft wonderous item"), FALSE, TRUE, FALSE)
	feato(FEAT_DEFLECT_ARROWS, libc.CString("deflect arrows"), FALSE, FALSE, FALSE)
	feato(FEAT_DODGE, libc.CString("dodge"), TRUE, TRUE, FALSE)
	feato(FEAT_EMPOWER_SPELL, libc.CString("empower spell"), FALSE, TRUE, FALSE)
	feato(FEAT_ENDURANCE, libc.CString("endurance"), FALSE, TRUE, FALSE)
	feato(FEAT_ENLARGE_SPELL, libc.CString("enlarge spell"), FALSE, FALSE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_BASTARD_SWORD, libc.CString("weapon proficiency - bastard sword"), FALSE, TRUE, FALSE)
	feato(FEAT_EXTEND_SPELL, libc.CString("extend spell"), FALSE, TRUE, FALSE)
	feato(FEAT_EXTRA_TURNING, libc.CString("extra turning"), FALSE, TRUE, FALSE)
	feato(FEAT_FAR_SHOT, libc.CString("far shot"), FALSE, FALSE, FALSE)
	feato(FEAT_FORGE_RING, libc.CString("forge ring"), FALSE, TRUE, FALSE)
	feato(FEAT_GREAT_CLEAVE, libc.CString("great cleave"), FALSE, FALSE, FALSE)
	feato(FEAT_GREAT_FORTITUDE, libc.CString("great fortitude"), TRUE, TRUE, FALSE)
	feato(FEAT_HEIGHTEN_SPELL, libc.CString("heighten spell"), FALSE, TRUE, FALSE)
	feato(FEAT_IMPROVED_BULL_RUSH, libc.CString("improved bull rush"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_CRITICAL, libc.CString("improved critical"), TRUE, TRUE, TRUE)
	feato(FEAT_IMPROVED_DISARM, libc.CString("improved disarm"), FALSE, TRUE, FALSE)
	feato(FEAT_IMPROVED_INITIATIVE, libc.CString("improved initiative"), TRUE, TRUE, FALSE)
	feato(FEAT_IMPROVED_TRIP, libc.CString("improved trip"), TRUE, TRUE, FALSE)
	feato(FEAT_IMPROVED_TWO_WEAPON_FIGHTING, libc.CString("improved two weapon fighting"), TRUE, TRUE, FALSE)
	feato(FEAT_IMPROVED_UNARMED_STRIKE, libc.CString("improved unarmed strike"), FALSE, FALSE, FALSE)
	feato(FEAT_IRON_WILL, libc.CString("iron will"), TRUE, TRUE, FALSE)
	feato(FEAT_LEADERSHIP, libc.CString("leadership"), FALSE, FALSE, FALSE)
	feato(FEAT_LIGHTNING_REFLEXES, libc.CString("lightning reflexes"), TRUE, TRUE, FALSE)
	feato(FEAT_MARTIAL_WEAPON_PROFICIENCY, libc.CString("martial weapon proficiency"), FALSE, TRUE, FALSE)
	feato(FEAT_MAXIMIZE_SPELL, libc.CString("maximize spell"), FALSE, TRUE, FALSE)
	feato(FEAT_MOBILITY, libc.CString("mobility"), TRUE, TRUE, FALSE)
	feato(FEAT_MOUNTED_ARCHERY, libc.CString("mounted archery"), FALSE, FALSE, FALSE)
	feato(FEAT_MOUNTED_COMBAT, libc.CString("mounted combat"), FALSE, FALSE, FALSE)
	feato(FEAT_POINT_BLANK_SHOT, libc.CString("point blank shot"), FALSE, FALSE, FALSE)
	feato(FEAT_POWER_ATTACK, libc.CString("power attack"), TRUE, TRUE, FALSE)
	feato(FEAT_PRECISE_SHOT, libc.CString("precise shot"), FALSE, FALSE, FALSE)
	feato(FEAT_QUICK_DRAW, libc.CString("quick draw"), FALSE, FALSE, FALSE)
	feato(FEAT_QUICKEN_SPELL, libc.CString("quicken spell"), FALSE, TRUE, FALSE)
	feato(FEAT_RAPID_SHOT, libc.CString("rapid shot"), FALSE, FALSE, FALSE)
	feato(FEAT_RIDE_BY_ATTACK, libc.CString("ride by attack"), FALSE, FALSE, FALSE)
	feato(FEAT_RUN, libc.CString("run"), FALSE, FALSE, FALSE)
	feato(FEAT_SCRIBE_SCROLL, libc.CString("scribe scroll"), FALSE, TRUE, FALSE)
	feato(FEAT_SHOT_ON_THE_RUN, libc.CString("shot on the run"), FALSE, FALSE, FALSE)
	feato(FEAT_SILENT_SPELL, libc.CString("silent spell"), FALSE, TRUE, FALSE)
	feato(FEAT_SIMPLE_WEAPON_PROFICIENCY, libc.CString("simple weapon proficiency"), TRUE, TRUE, FALSE)
	feato(FEAT_SKILL_FOCUS, libc.CString("skill focus"), TRUE, TRUE, TRUE)
	feato(FEAT_SPELL_FOCUS, libc.CString("spell focus"), FALSE, TRUE, TRUE)
	feato(FEAT_SPELL_MASTERY, libc.CString("spell mastery"), FALSE, TRUE, TRUE)
	feato(FEAT_SPELL_PENETRATION, libc.CString("spell penetration"), FALSE, TRUE, FALSE)
	feato(FEAT_SPIRITED_CHARGE, libc.CString("spirited charge"), FALSE, FALSE, FALSE)
	feato(FEAT_SPRING_ATTACK, libc.CString("spring attack"), TRUE, FALSE, FALSE)
	feato(FEAT_STILL_SPELL, libc.CString("still spell"), FALSE, TRUE, FALSE)
	feato(FEAT_STUNNING_FIST, libc.CString("stunning fist"), FALSE, TRUE, FALSE)
	feato(FEAT_SUNDER, libc.CString("sunder"), FALSE, TRUE, FALSE)
	feato(FEAT_TOUGHNESS, libc.CString("toughness"), TRUE, TRUE, TRUE)
	feato(FEAT_TRACK, libc.CString("track"), FALSE, TRUE, FALSE)
	feato(FEAT_TRAMPLE, libc.CString("trample"), FALSE, FALSE, FALSE)
	feato(FEAT_TWO_WEAPON_FIGHTING, libc.CString("two weapon fighting"), TRUE, TRUE, FALSE)
	feato(FEAT_WEAPON_FINESSE, libc.CString("weapon finesse"), TRUE, TRUE, TRUE)
	feato(FEAT_WEAPON_FOCUS, libc.CString("weapon focus"), FALSE, TRUE, TRUE)
	feato(FEAT_WEAPON_SPECIALIZATION, libc.CString("weapon specialization"), FALSE, TRUE, TRUE)
	feato(FEAT_WHIRLWIND_ATTACK, libc.CString("whirlwind attack"), FALSE, TRUE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_DRUID, libc.CString("weapon proficiency - druids"), FALSE, FALSE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_ROGUE, libc.CString("weapon proficiency - rogues"), FALSE, FALSE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_MONK, libc.CString("weapon proficiency - monks"), FALSE, FALSE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_WIZARD, libc.CString("weapon proficiency - wizards"), FALSE, FALSE, FALSE)
	feato(FEAT_WEAPON_PROFICIENCY_ELF, libc.CString("weapon proficiency - elves"), FALSE, FALSE, FALSE)
	feato(FEAT_ARMOR_PROFICIENCY_SHIELD, libc.CString("shield armor proficiency"), FALSE, FALSE, FALSE)
	feato(FEAT_SNEAK_ATTACK, libc.CString("sneak attack"), TRUE, FALSE, TRUE)
	feato(FEAT_EVASION, libc.CString("evasion"), TRUE, FALSE, FALSE)
	feato(FEAT_IMPROVED_EVASION, libc.CString("improved evasion"), TRUE, FALSE, FALSE)
	feato(FEAT_ACROBATIC, libc.CString("acrobatic"), TRUE, TRUE, FALSE)
	feato(FEAT_AGILE, libc.CString("agile"), TRUE, TRUE, FALSE)
	feato(FEAT_ALERTNESS, libc.CString("alertness"), TRUE, FALSE, FALSE)
	feato(FEAT_ANIMAL_AFFINITY, libc.CString("animal affinity"), FALSE, TRUE, FALSE)
	feato(FEAT_ATHLETIC, libc.CString("athletic"), TRUE, TRUE, FALSE)
	feato(FEAT_AUGMENT_SUMMONING, libc.CString("augment summoning"), FALSE, FALSE, FALSE)
	feato(FEAT_COMBAT_EXPERTISE, libc.CString("combat expertise"), FALSE, FALSE, FALSE)
	feato(FEAT_DECEITFUL, libc.CString("deceitful"), TRUE, TRUE, FALSE)
	feato(FEAT_DEFT_HANDS, libc.CString("deft hands"), FALSE, TRUE, FALSE)
	feato(FEAT_DIEHARD, libc.CString("diehard"), TRUE, FALSE, FALSE)
	feato(FEAT_DILIGENT, libc.CString("diligent"), TRUE, TRUE, FALSE)
	feato(FEAT_ESCHEW_MATERIALS, libc.CString("eschew materials"), FALSE, FALSE, FALSE)
	feato(FEAT_EXOTIC_WEAPON_PROFICIENCY, libc.CString("exotic weapon proficiency"), FALSE, FALSE, FALSE)
	feato(FEAT_GREATER_SPELL_FOCUS, libc.CString("greater spell focus"), FALSE, FALSE, TRUE)
	feato(FEAT_GREATER_SPELL_PENETRATION, libc.CString("greater spell penetration"), FALSE, FALSE, FALSE)
	feato(FEAT_GREATER_TWO_WEAPON_FIGHTING, libc.CString("greater two weapon fighting"), TRUE, FALSE, FALSE)
	feato(FEAT_GREATER_WEAPON_FOCUS, libc.CString("greater weapon focus"), FALSE, TRUE, TRUE)
	feato(FEAT_GREATER_WEAPON_SPECIALIZATION, libc.CString("greater weapon specialization"), FALSE, TRUE, TRUE)
	feato(FEAT_IMPROVED_COUNTERSPELL, libc.CString("improved counterspell"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_FAMILIAR, libc.CString("improved familiar"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_FEINT, libc.CString("improved feint"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_GRAPPLE, libc.CString("improved grapple"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_OVERRUN, libc.CString("improved overrun"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_PRECISE_SHOT, libc.CString("improved precise shot"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_SHIELD_BASH, libc.CString("improved shield bash"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_SUNDER, libc.CString("improved sunder"), FALSE, FALSE, FALSE)
	feato(FEAT_IMPROVED_TURNING, libc.CString("improved turning"), FALSE, FALSE, FALSE)
	feato(FEAT_INVESTIGATOR, libc.CString("investigator"), FALSE, TRUE, FALSE)
	feato(FEAT_MAGICAL_APTITUDE, libc.CString("magical aptitude"), FALSE, TRUE, FALSE)
	feato(FEAT_MANYSHOT, libc.CString("manyshot"), FALSE, FALSE, FALSE)
	feato(FEAT_NATURAL_SPELL, libc.CString("natural spell"), FALSE, FALSE, FALSE)
	feato(FEAT_NEGOTIATOR, libc.CString("negotiator"), FALSE, TRUE, FALSE)
	feato(FEAT_NIMBLE_FINGERS, libc.CString("nimble fingers"), FALSE, TRUE, FALSE)
	feato(FEAT_PERSUASIVE, libc.CString("persuasive"), FALSE, TRUE, FALSE)
	feato(FEAT_RAPID_RELOAD, libc.CString("rapid reload"), FALSE, FALSE, FALSE)
	feato(FEAT_SELF_SUFFICIENT, libc.CString("self sufficient"), FALSE, TRUE, FALSE)
	feato(FEAT_STEALTHY, libc.CString("stealthy"), TRUE, TRUE, FALSE)
	feato(FEAT_ARMOR_PROFICIENCY_TOWER_SHIELD, libc.CString("tower shield armor proficiency"), FALSE, FALSE, FALSE)
	feato(FEAT_TWO_WEAPON_DEFENSE, libc.CString("two weapon defense"), FALSE, FALSE, FALSE)
	feato(FEAT_WIDEN_SPELL, libc.CString("widen spell"), FALSE, FALSE, FALSE)
}
func feat_is_available(ch *char_data, featnum int, iarg int, sarg *byte) int {
	if featnum > NUM_FEATS_DEFINED {
		return FALSE
	}
	if int(ch.Feats[featnum]) != 0 && feat_list[featnum].Can_stack == 0 {
		return FALSE
	}
	switch featnum {
	case FEAT_ARMOR_PROFICIENCY_HEAVY:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_MEDIUM]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_ARMOR_PROFICIENCY_MEDIUM:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_LIGHT]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_AUGMENT_SUMMONING:
		if int(ch.Feats[FEAT_SPELL_FOCUS]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_SHIELD_BASH:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_SHIELD]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_DODGE:
		if int(ch.Aff_abils.Dex) >= 13 {
			return TRUE
		}
		return FALSE
	case FEAT_COMBAT_EXPERTISE:
		if int(ch.Aff_abils.Intel) >= 13 {
			return TRUE
		}
		return FALSE
	case FEAT_DIEHARD:
		if int(ch.Feats[FEAT_ENDURANCE]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_SUNDER:
		fallthrough
	case FEAT_IMPROVED_OVERRUN:
		fallthrough
	case FEAT_IMPROVED_BULL_RUSH:
		if int(ch.Feats[FEAT_POWER_ATTACK]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_MOBILITY:
		if int(ch.Feats[FEAT_DODGE]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_EXOTIC_WEAPON_PROFICIENCY:
		if ch.Accuracy >= 1 {
			return TRUE
		}
		return FALSE
	case FEAT_WEAPON_PROFICIENCY_BASTARD_SWORD:
		if ch.Accuracy >= 1 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_FEINT:
		fallthrough
	case FEAT_IMPROVED_DISARM:
		fallthrough
	case FEAT_IMPROVED_TRIP:
		if int(ch.Feats[FEAT_COMBAT_EXPERTISE]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_GRAPPLE:
		fallthrough
	case FEAT_DEFLECT_ARROWS:
		if int(ch.Feats[FEAT_IMPROVED_UNARMED_STRIKE]) != 0 {
			if int(ch.Aff_abils.Dex) >= 13 {
				return TRUE
			}
		}
		return FALSE
	case FEAT_STUNNING_FIST:
		if int(ch.Aff_abils.Dex) >= 13 {
			if int(ch.Aff_abils.Wis) >= 13 {
				if int(ch.Feats[FEAT_IMPROVED_UNARMED_STRIKE]) != 0 {
					if ch.Accuracy >= 8 {
						return TRUE
					}
				}
			}
		}
		return FALSE
	case FEAT_POWER_ATTACK:
		if int(ch.Aff_abils.Str) >= 13 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_PRECISE_SHOT:
		if int(ch.Aff_abils.Dex) >= 19 {
			if int(ch.Feats[FEAT_POINT_BLANK_SHOT]) != 0 {
				if ch.Accuracy >= 11 {
					return TRUE
				}
			}
		}
		return FALSE
	case FEAT_CLEAVE:
		if int(ch.Feats[FEAT_POWER_ATTACK]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_TWO_WEAPON_FIGHTING:
		if int(ch.Aff_abils.Dex) >= 15 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_TWO_WEAPON_FIGHTING:
		if int(ch.Aff_abils.Dex) >= 17 && int(ch.Feats[FEAT_TWO_WEAPON_FIGHTING]) != 0 && ch.Accuracy >= 6 {
			return TRUE
		}
		return FALSE
	case FEAT_IMPROVED_CRITICAL:
		if ch.Accuracy < 8 {
			return FALSE
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_FAR_SHOT:
		if int(ch.Feats[FEAT_POINT_BLANK_SHOT]) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_WEAPON_FINESSE:
		fallthrough
	case FEAT_WEAPON_FOCUS:
		if ch.Accuracy < 1 {
			return FALSE
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_WEAPON_SPECIALIZATION:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 4 {
			return FALSE
		}
		if iarg == 0 || is_proficient_with_weapon(ch, iarg) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_GREATER_WEAPON_FOCUS:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 8 {
			return FALSE
		}
		if iarg == 0 {
			return TRUE
		}
		if is_proficient_with_weapon(ch, iarg) != 0 && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_FOCUS][:], bitvector_t(int32(iarg))) {
			return TRUE
		}
		return FALSE
	case FEAT_WHIRLWIND_ATTACK:
		if int(ch.Feats[FEAT_COMBAT_EXPERTISE]) != 0 {
			if int(ch.Feats[FEAT_DODGE]) != 0 {
				if int(ch.Feats[FEAT_MOBILITY]) != 0 {
					if int(ch.Feats[FEAT_SPRING_ATTACK]) != 0 {
						if ch.Accuracy >= 4 {
							return TRUE
						}
					}
				}
			}
		}
		return FALSE
	case FEAT_GREATER_WEAPON_SPECIALIZATION:
		if ((ch.Chclasses[CLASS_NAIL]) + (ch.Epicclasses[CLASS_NAIL])) < 12 {
			return FALSE
		}
		if iarg == 0 {
			return TRUE
		}
		if is_proficient_with_weapon(ch, iarg) != 0 && IS_SET_AR(ch.Combat_feats[CFEAT_GREATER_WEAPON_FOCUS][:], bitvector_t(int32(iarg))) && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_SPECIALIZATION][:], bitvector_t(int32(iarg))) && IS_SET_AR(ch.Combat_feats[CFEAT_WEAPON_FOCUS][:], bitvector_t(int32(iarg))) {
			return TRUE
		}
		return FALSE
	case FEAT_SPELL_FOCUS:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_SPELL_PENETRATION:
		if GET_LEVEL(ch) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_BREW_POTION:
		if GET_LEVEL(ch) >= 3 {
			return TRUE
		}
		return FALSE
	case FEAT_CRAFT_MAGICAL_ARMS_AND_ARMOR:
		if GET_LEVEL(ch) >= 5 {
			return TRUE
		}
		return FALSE
	case FEAT_CRAFT_ROD:
		if GET_LEVEL(ch) >= 9 {
			return TRUE
		}
		return FALSE
	case FEAT_CRAFT_STAFF:
		if GET_LEVEL(ch) >= 12 {
			return TRUE
		}
		return FALSE
	case FEAT_CRAFT_WAND:
		if GET_LEVEL(ch) >= 5 {
			return TRUE
		}
		return FALSE
	case FEAT_FORGE_RING:
		if GET_LEVEL(ch) >= 5 {
			return TRUE
		}
		return FALSE
	case FEAT_SCRIBE_SCROLL:
		if GET_LEVEL(ch) >= 1 {
			return TRUE
		}
		return FALSE
	case FEAT_EMPOWER_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_EXTEND_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_HEIGHTEN_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_MAXIMIZE_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_QUICKEN_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_SILENT_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_STILL_SPELL:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_EXTRA_TURNING:
		if ((ch.Chclasses[CLASS_PICCOLO]) + (ch.Epicclasses[CLASS_PICCOLO])) != 0 {
			return TRUE
		}
		return FALSE
	case FEAT_SPELL_MASTERY:
		if ((ch.Chclasses[CLASS_ROSHI]) + (ch.Epicclasses[CLASS_ROSHI])) != 0 {
			return TRUE
		}
		return FALSE
	default:
		return TRUE
	}
}
func is_proficient_with_armor(ch *char_data, cmarmor_type int) int {
	switch cmarmor_type {
	case ARMOR_TYPE_LIGHT:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_LIGHT]) != 0 {
			return TRUE
		}
	case ARMOR_TYPE_MEDIUM:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_MEDIUM]) != 0 {
			return TRUE
		}
	case ARMOR_TYPE_HEAVY:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_HEAVY]) != 0 {
			return TRUE
		}
	case ARMOR_TYPE_SHIELD:
		if int(ch.Feats[FEAT_ARMOR_PROFICIENCY_SHIELD]) != 0 {
			return TRUE
		}
	}
	return FALSE
}
func is_proficient_with_weapon(ch *char_data, cmweapon_type int) int {
	switch cmweapon_type {
	case WEAPON_TYPE_UNARMED:
		return 1
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
			return TRUE
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
			return TRUE
		}
	default:
		return FALSE
	}
	return FALSE
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
		none_shown int = TRUE
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
				none_shown = FALSE
			case FEAT_TOUGHNESS:
				temp_value = int(ch.Feats[FEAT_TOUGHNESS])
				added_hp = temp_value * 3
				stdio.Sprintf(&buf[0], "%-20s (+%d hp)\r\n", feat_list[i].Name, added_hp)
				libc.StrCat(&buf2[0], &buf[0])
				none_shown = FALSE
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
					if IS_SET_AR(ch.Combat_feats[feat_to_subfeat(i)][:], bitvector_t(int32(j))) {
						stdio.Sprintf(&buf[0], "%-20s (%s)\r\n", feat_list[i].Name, weapon_type[j])
						libc.StrCat(&buf2[0], &buf[0])
						none_shown = FALSE
					}
				}
			default:
				stdio.Sprintf(&buf[0], "%-20s\r\n", feat_list[i].Name)
				libc.StrCat(&buf2[0], &buf[0])
				none_shown = FALSE
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
		none_shown int = TRUE
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
		if feat_is_available(ch, i, 0, nil) != 0 && feat_list[i].In_game != 0 && feat_list[i].Can_learn != 0 {
			stdio.Sprintf(&buf[0], "%-20s\r\n", feat_list[i].Name)
			libc.StrCat(&buf2[0], &buf[0])
			none_shown = FALSE
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
		none_shown int = TRUE
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
			none_shown = FALSE
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
		if is_abbrev(name, feat_list[ftindex].Name) != 0 {
			return ftindex
		}
		ok = TRUE
		temp = any_one_arg(feat_list[ftindex].Name, &first[0])
		temp2 = any_one_arg(name, &first2[0])
		for first[0] != 0 && first2[0] != 0 && ok != 0 {
			if is_abbrev(&first2[0], &first[0]) == 0 {
				ok = FALSE
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
	if is_abbrev(&arg[0], libc.CString("known")) != 0 || arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax is \"feats <available | complete | known>\".\r\n"))
		list_feats_known(ch)
	} else if is_abbrev(&arg[0], libc.CString("available")) != 0 {
		list_feats_available(ch)
	} else if is_abbrev(&arg[0], libc.CString("complete")) != 0 {
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
