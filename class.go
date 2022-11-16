package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const CONFIG_LEVEL_VERSION = 0
const CONFIG_LEVEL_EXPERIENCE = 1
const CONFIG_LEVEL_VERNUM = 2
const CONFIG_LEVEL_FORTITUDE = 3
const CONFIG_LEVEL_REFLEX = 4
const CONFIG_LEVEL_WILL = 5
const CONFIG_LEVEL_BASEHIT = 6
const SAVE_MANUAL = 0
const SAVE_LOW = 1
const SAVE_HIGH = 2
const BASEHIT_MANUAL = 0
const BASEHIT_LOW = 1
const BASEHIT_MEDIUM = 2
const BASEHIT_HIGH = 3

var class_abbrevs [32]*byte = [32]*byte{libc.CString("Ro"), libc.CString("Pi"), libc.CString("Kr"), libc.CString("Na"), libc.CString("Ba"), libc.CString("Gi"), libc.CString("Fr"), libc.CString("Ta"), libc.CString("An"), libc.CString("Da"), libc.CString("Ki"), libc.CString("Ji"), libc.CString("Ts"), libc.CString("Ku"), libc.CString("As"), libc.CString("Bl"), libc.CString("Dd"), libc.CString("Du"), libc.CString("Dw"), libc.CString("Ek"), libc.CString("Ht"), libc.CString("Hw"), libc.CString("Lo"), libc.CString("Mt"), libc.CString("Sh"), libc.CString("Th"), libc.CString("Ex"), libc.CString("Ad"), libc.CString("Co"), libc.CString("Ar"), libc.CString("Wa"), libc.CString("\n")}
var pc_class_types [32]*byte = [32]*byte{libc.CString("Roshi"), libc.CString("Piccolo"), libc.CString("Krane"), libc.CString("Nail"), libc.CString("Bardock"), libc.CString("Ginyu"), libc.CString("Frieza"), libc.CString("Tapion"), libc.CString("Android 16"), libc.CString("Dabura"), libc.CString("Kibito"), libc.CString("Jinto"), libc.CString("Tsuna"), libc.CString("Kurzak"), libc.CString("Assassin"), libc.CString("Blackguard"), libc.CString("Dragon Disciple"), libc.CString("Duelist"), libc.CString("Dwarven Defender"), libc.CString("Eldritch Knight"), libc.CString("Hierophant"), libc.CString("Horizon Walker"), libc.CString("Loremaster"), libc.CString("Mystic Theurge"), libc.CString("Shinobi"), libc.CString("Thaumaturgist"), libc.CString("Expert"), libc.CString("Adept"), libc.CString("Commoner"), libc.CString("Aristrocrat"), libc.CString("Warrior"), libc.CString("\n")}
var class_names [32]*byte = [32]*byte{libc.CString("roshi"), libc.CString("piccolo"), libc.CString("krane"), libc.CString("nail"), libc.CString("bardock"), libc.CString("ginyu"), libc.CString("frieza"), libc.CString("tapion"), libc.CString("android 16"), libc.CString("dabura"), libc.CString("kibito"), libc.CString("jinto"), libc.CString("tsuna"), libc.CString("kurzak"), libc.CString("assassin"), libc.CString("blackguard"), libc.CString("dragon disciple"), libc.CString("duelist"), libc.CString("dwarven defender"), libc.CString("eldritch knight"), libc.CString("hierophant"), libc.CString("horizon walker"), libc.CString("loremaster"), libc.CString("mystic theurge"), libc.CString("shadowdancer"), libc.CString("thaumaturgist"), libc.CString("artisan"), libc.CString("magi"), libc.CString("normal"), libc.CString("noble"), libc.CString("soldier"), libc.CString("\n")}
var class_display [31]*byte = [31]*byte{libc.CString("@B1@W) @MRoshi\r\n"), libc.CString("@B2@W) @WPiccolo\r\n"), libc.CString("@B3@W) @YKrane\r\n"), libc.CString("@B4@W) @BNail\r\n"), libc.CString("@B5@W) @BBardock\r\n"), libc.CString("@B6@W) @BGinyu\r\n"), libc.CString("@B7@W) @WFrieza\r\n"), libc.CString("@B8@W) @YTapion\r\n"), libc.CString("@B9@W) @BAndroid 16\r\n"), libc.CString("@B10@W) @BDabura\r\n"), libc.CString("@B11@W) @BKibito\r\n"), libc.CString("@B12@W) @BJinto\r\n"), libc.CString("@B13@W) @BTsuna\r\n"), libc.CString("@B14@W) @BKurzak\r\n"), libc.CString("assassin (P)\r\n"), libc.CString("blackguard (P)\r\n"), libc.CString("dragon disciple (P)\r\n"), libc.CString("duelist (P)\r\n"), libc.CString("dwarven defender (P)\r\n"), libc.CString("eldritch knight (P)\r\n"), libc.CString("hierophant (P)\r\n"), libc.CString("horizon walker (P)\r\n"), libc.CString("loremaster (P)\r\n"), libc.CString("mystic theurge (P)\r\n"), libc.CString("shadowdancer (P)\r\n"), libc.CString("thaumaturgist (P)\r\n"), libc.CString("Artisan NPC\r\n"), libc.CString("Magi NPC\r\n"), libc.CString("Normal NPC\r\n"), libc.CString("Noble NPC\r\n"), libc.CString("Soldier NPC\r\n")}
var class_ok_race [24][31]bool = [24][31]bool{{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: false, 7: false, 8: true, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: true, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: true, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: false, 10: true, 11: false, 12: false, 13: true, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}, {0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true, 8: true, 9: true, 10: true, 11: false, 12: true, 13: false, 14: false, 15: false, 16: true, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false}}
var class_ok_align [9][31]bool = [9][31]bool{{true, true, true, true, true, true, true, false, false, true, false, true, false, true, false, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, true, true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, false, true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, true, false, true, true, false, true, false, true, false, true, false, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, true, true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, true, true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, true, false, true, false, false, true, false, true, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, true, true, true, true, true, true, true, true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}, {true, true, true, true, false, false, true, false, true, true, true, true, true, true, true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, true, true}}
var favored_class [24]int = [24]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
var prestige_classes [31]bool = [31]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, true, true, true, true, true, true, true, true, true, true, true, false, false, false, false, false}
var class_max_ranks [31]int = [31]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 10, 10, 10, 10, 10, 10, 5, 10, 10, 10, 10, 5, -1, -1, -1, -1, -1}

func parse_class(ch *char_data, arg int) int {
	var chclass int = -1
	switch arg {
	case 1:
		chclass = CLASS_ROSHI
	case 2:
		chclass = CLASS_PICCOLO
	case 3:
		chclass = CLASS_KRANE
	case 4:
		chclass = CLASS_NAIL
	case 5:
		chclass = CLASS_BARDOCK
	case 6:
		chclass = CLASS_GINYU
	case 7:
		chclass = CLASS_FRIEZA
	case 8:
		chclass = CLASS_TAPION
	case 9:
		chclass = CLASS_ANDSIX
	case 10:
		chclass = CLASS_DABURA
	case 11:
		chclass = CLASS_KABITO
	case 12:
		chclass = CLASS_JINTO
	case 13:
		chclass = CLASS_TSUNA
	case 14:
		chclass = CLASS_KURZAK
	default:
		chclass = -1
	}
	if chclass >= 0 && chclass < 14 {
		if !class_ok_race[int(ch.Race)][chclass] {
			chclass = -1
		}
	}
	return chclass
}

var guild_info [6]guild_info_type = [6]guild_info_type{{Pc_class: CLASS_ROSHI, Guild_room: 3017, Direction: SCMD_EAST}, {Pc_class: CLASS_PICCOLO, Guild_room: 3004, Direction: SCMD_NORTH}, {Pc_class: CLASS_KRANE, Guild_room: 3027, Direction: SCMD_EAST}, {Pc_class: CLASS_NAIL, Guild_room: 3021, Direction: SCMD_EAST}, {Pc_class: -999, Guild_room: 5065, Direction: SCMD_WEST}, {Pc_class: -1, Guild_room: room_vnum(-1), Direction: -1}}
var config_sect [8]*byte = [8]*byte{libc.CString("version"), libc.CString("experience"), libc.CString("vernum"), libc.CString("fortitude"), libc.CString("reflex"), libc.CString("will"), libc.CString("basehit"), libc.CString("\n")}
var level_version [256]byte
var level_vernum int = 0
var save_classes [3][31]int
var basehit_classes [31]int
var save_type_names [3]*byte = [3]*byte{libc.CString("manual"), libc.CString("low"), libc.CString("high")}
var basehit_type_names [4]*byte = [4]*byte{libc.CString("manual"), libc.CString("low"), libc.CString("medium"), libc.CString("high")}
var class_template [14][6]int = [14][6]int{0: {10, 13, 13, 18, 16, 10}, 1: {13, 10, 13, 10, 18, 16}, 2: {13, 18, 13, 16, 10, 10}, 3: {18, 13, 16, 10, 13, 10}, 4: {13, 16, 13, 10, 18, 10}, 5: {18, 10, 13, 10, 16, 13}, 6: {10, 13, 13, 18, 16, 10}, 7: {13, 10, 13, 10, 18, 16}, 8: {13, 18, 13, 16, 10, 10}, 9: {13, 13, 13, 10, 16, 10}, 10: {13, 12, 18, 10, 12, 10}}
var race_template [24][6]int = [24][6]int{{13, 13, 13, 13, 13, 13}, {16, 12, 14, 10, 12, 12}, {14, 14, 12, 12, 12, 12}, {10, 16, 10, 13, 14, 14}, {14, 12, 13, 12, 14, 12}, {12, 12, 15, 13, 13, 13}, {10, 14, 10, 15, 13, 10}, {14, 13, 14, 12, 13, 12}, {15, 10, 15, 12, 12, 10}, {14, 14, 14, 12, 10, 12}, {14, 13, 14, 10, 12, 10}, {15, 10, 15, 10, 12, 14}, {11, 14, 10, 14, 14, 11}, {10, 14, 10, 16, 16, 12}, {13, 13, 13, 13, 13, 13}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}, {13, 13, 13, 12, 12, 14}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}, {10, 10, 10, 10, 10, 10}}

func cedit_creation(ch *char_data) {
	switch config_info.Creation.Method {
	case CEDIT_CREATION_METHOD_3:
	case CEDIT_CREATION_METHOD_4:
	case CEDIT_CREATION_METHOD_5:
	case CEDIT_CREATION_METHOD_2:
		fallthrough
	case CEDIT_CREATION_METHOD_1:
		fallthrough
	default:
	}
	racial_ability_modifiers(ch)
	racial_body_parts(ch)
	ch.Aff_abils = ch.Real_abils
}

var class_hit_die_size [31]int = [31]int{4, 8, 6, 10, 8, 10, 4, 8, 6, 8, 12, 8, 4, 4, 6, 10, 12, 10, 12, 6, 8, 8, 4, 4, 8, 4, 6, 6, 4, 8, 8}

func do_start(ch *char_data) {
	var (
		punch int
		obj   *obj_data
	)
	ch.Level = 1
	ch.Race_level = 0
	ch.Level_adj = 0
	ch.Chclasses[ch.Chclass] = 1
	ch.Exp = 1
	if int(ch.Race) == RACE_ANDROID {
		ch.Player_specials.Conditions[HUNGER] = -1
		ch.Player_specials.Conditions[THIRST] = -1
		ch.Player_specials.Conditions[DRUNK] = -1
	} else if int(ch.Race) == RACE_BIO && ((ch.Genome[0]) == 3 || (ch.Genome[1]) == 3) {
		ch.Player_specials.Conditions[HUNGER] = -1
		ch.Player_specials.Conditions[DRUNK] = 0
		ch.Player_specials.Conditions[THIRST] = 48
	} else if int(ch.Race) == RACE_NAMEK {
		ch.Player_specials.Conditions[HUNGER] = -1
		ch.Player_specials.Conditions[DRUNK] = 0
		ch.Player_specials.Conditions[THIRST] = 48
	} else {
		ch.Player_specials.Conditions[THIRST] = 48
		ch.Player_specials.Conditions[HUNGER] = 48
		ch.Player_specials.Conditions[DRUNK] = 0
	}
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_AUTOEXIT)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_HINTS)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_NOMUSIC)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_DISPHP)
	ch.Limb_condition[0] = 100
	ch.Limb_condition[1] = 100
	ch.Limb_condition[2] = 100
	ch.Limb_condition[3] = 100
	SET_BIT_AR(ch.Act[:], PLR_HEAD)
	ch.Skill_slots = 30
	if int(ch.Race) == RACE_HUMAN {
		ch.Skill_slots += 1
	} else if int(ch.Race) == RACE_SAIYAN {
		ch.Skill_slots -= 1
	} else if int(ch.Race) == RACE_TRUFFLE {
		ch.Skill_slots += 2
	} else if int(ch.Race) == RACE_HALFBREED {
		ch.Skill_slots += 1
	} else if int(ch.Race) == RACE_MAJIN {
		ch.Skill_slots -= 1
	} else if int(ch.Race) == RACE_KAI {
		ch.Skill_slots += 4
	}
	if int(ch.Chclass) == CLASS_TSUNA || int(ch.Chclass) == CLASS_KABITO || int(ch.Chclass) == CLASS_NAIL {
		ch.Skill_slots += 5
	}
	if (ch.Bonuses[BONUS_GMEMORY]) != 0 {
		ch.Skill_slots += 2
	}
	if (ch.Bonuses[BONUS_BMEMORY]) != 0 {
		ch.Skill_slots -= 5
	}
	if int(ch.Race) == RACE_SAIYAN || int(ch.Race) == RACE_HALFBREED {
		if int(ch.Race) != RACE_HALFBREED || int(ch.Race) == RACE_HALFBREED && ch.Player_specials.Racial_pref != 1 {
			SET_BIT_AR(ch.Act[:], PLR_STAIL)
		}
	}
	if int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO {
		SET_BIT_AR(ch.Act[:], PLR_TAIL)
	}
	if int(ch.Race) == RACE_MAJIN {
		ch.Absorbs = 0
		ch.IngestLearned = 0
	}
	if int(ch.Race) == RACE_BIO {
		ch.Absorbs = 3
	}
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_VIEWORDER)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_DISPMOVE)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_DISPKI)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_DISPEXP)
	SET_BIT_AR(ch.Player_specials.Pref[:], PRF_DISPTNL)
	if !PLR_FLAGGED(ch, PLR_FORGET) {
		if ch.Choice == 1 {
			punch = rand_number(30, 40)
			for {
				ch.Skills[SKILL_BLOCK] = int8(punch)
				if true {
					break
				}
			}
		}
		if ch.Choice == 2 {
			punch = rand_number(10, 20)
			for {
				ch.Skills[SKILL_PUNCH] = int8(punch)
				if true {
					break
				}
			}
		}
		if ch.Choice == 3 {
			punch = rand_number(30, 40)
			for {
				ch.Skills[SKILL_KICK] = int8(punch)
				if true {
					break
				}
			}
		}
		if ch.Choice == 4 {
			punch = rand_number(20, 30)
			for {
				ch.Skills[SKILL_SLAM] = int8(punch)
				if true {
					break
				}
			}
		}
		if ch.Choice == 5 {
			punch = rand_number(20, 30)
			for {
				ch.Skills[SKILL_FOCUS] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_HUMAN {
			punch = rand_number(5, 15)
			for {
				ch.Skills[SKILL_BUILD] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_TRUFFLE {
			punch = rand_number(15, 25)
			for {
				ch.Skills[SKILL_BUILD] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_KONATSU {
			punch = rand_number(50, 60)
			for {
				ch.Skills[SKILL_SWORD] = int8(punch)
				if true {
					break
				}
			}
			punch = rand_number(10, 30)
			for {
				ch.Skills[SKILL_MOVE_SILENTLY] = int8(punch)
				if true {
					break
				}
			}
			punch = rand_number(10, 30)
			for {
				ch.Skills[SKILL_HIDE] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Chclass) == CLASS_TAPION || int(ch.Chclass) == CLASS_GINYU || int(ch.Chclass) == CLASS_DABURA || int(ch.Chclass) == CLASS_KURZAK {
			punch = rand_number(30, 40)
			if int(ch.Chclass) == CLASS_KURZAK || int(ch.Chclass) == CLASS_TAPION {
				punch += rand_number(5, 10)
			}
			for {
				ch.Skills[SKILL_THROW] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_KAI || int(ch.Race) == RACE_KANASSAN {
			punch = rand_number(40, 60)
			for {
				ch.Skills[SKILL_FOCUS] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_KANASSAN {
			punch = rand_number(40, 60)
			for {
				ch.Skills[SKILL_CONCENTRATION] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_KAI {
			punch = rand_number(30, 50)
			for {
				ch.Skills[SKILL_HEAL] = int8(punch)
				if true {
					break
				}
			}
		}
		if int(ch.Race) == RACE_DEMON {
			punch = rand_number(50, 60)
			for {
				ch.Skills[SKILL_SPEAR] = int8(punch)
				if true {
					break
				}
			}
		}
		punch = 0
		punch = rand_number(50, 70)
		for {
			ch.Skills[SKILL_PUNCH] = int8(int(ch.Skills[SKILL_PUNCH]) + punch)
			if true {
				break
			}
		}
	} else {
		REMOVE_BIT_AR(ch.Act[:], PLR_FORGET)
	}
	if int(ch.Race) == RACE_KAI || int(ch.Race) == RACE_KANASSAN {
		punch = rand_number(15, 30)
		for {
			ch.Skills[SKILL_TELEPATHY] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Race) == RACE_MAJIN || int(ch.Race) == RACE_NAMEK || int(ch.Race) == RACE_BIO {
		punch = rand_number(10, 16)
		for {
			ch.Skills[SKILL_REGENERATE] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_ABSORB) {
		punch = rand_number(25, 35)
		for {
			ch.Skills[SKILL_ABSORB] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Race) == RACE_BIO {
		punch = rand_number(15, 25)
		for {
			ch.Skills[SKILL_ABSORB] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Race) == RACE_ARLIAN {
		punch = rand_number(30, 50)
		for {
			ch.Skills[SKILL_SEISHOU] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Race) == RACE_ICER {
		punch = rand_number(20, 30)
		for {
			ch.Skills[SKILL_TAILWHIP] = int8(punch)
			if true {
				break
			}
		}
	}
	if int(ch.Chclass) < 0 || int(ch.Chclass) > NUM_CLASSES {
		basic_mud_log(libc.CString("Unknown character class %d in do_start, resetting."), ch.Chclass)
		ch.Chclass = 0
	}
	if ch.Alignment < 51 && ch.Alignment > -51 {
		set_title(ch, libc.CString("the Warrior"))
	}
	if ch.Alignment >= 51 {
		set_title(ch, libc.CString("the Hero"))
	}
	if ch.Alignment <= -51 {
		set_title(ch, libc.CString("The Villain"))
	}
	if ch.Gold <= 0 {
		ch.Gold = dice(3, 6) * 10
	}
	switch ch.Race {
	case RACE_HUMAN:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_SAIYAN:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_HALFBREED:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_ICER:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_KONATSU:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_NAMEK:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_MUTANT:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	case RACE_ANDROID:
		SET_BIT_AR(ch.Affected_by[:], AFF_INFRAVISION)
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	default:
		for {
			ch.Skills[SKILL_LANG_COMMON] = 1
			if true {
				break
			}
		}
	}
	ch.Player_specials.Speaking = SKILL_LANG_COMMON
	ch.Lifeperc = 75
	obj = read_object(17, VIRTUAL)
	obj_to_char(obj, ch)
	if int(ch.Race) == RACE_HOSHIJIN {
		obj = read_object(3428, VIRTUAL)
		obj_to_char(obj, ch)
	}
	var obj2 *obj_data
	obj2 = read_object(0x464E, VIRTUAL)
	obj_to_char(obj2, ch)
	if int(ch.Chclass) == CLASS_TAPION || int(ch.Chclass) == CLASS_GINYU {
		var throw *obj_data
		throw = read_object(19050, VIRTUAL)
		obj_to_char(throw, ch)
		if rand_number(1, 2) == 2 {
			throw = nil
			throw = read_object(19050, VIRTUAL)
			obj_to_char(throw, ch)
		}
		if rand_number(1, 2) == 2 {
			throw = nil
			throw = read_object(19050, VIRTUAL)
			obj_to_char(throw, ch)
		}
	} else if int(ch.Chclass) == CLASS_DABURA {
		var throw *obj_data
		throw = read_object(0x4A6F, VIRTUAL)
		obj_to_char(throw, ch)
		throw = nil
		throw = read_object(0x4A6F, VIRTUAL)
		obj_to_char(throw, ch)
	}
	send_to_imm(libc.CString("New character created, %s, by user, %s."), GET_NAME(ch), GET_USER(ch))
	advance_level(ch, int(ch.Chclass))
	if ch.Max_hit < 100 {
		ch.Max_hit = 100
	}
	if ch.Max_mana < 100 {
		ch.Max_mana = 100
	}
	if ch.Max_move < 100 {
		ch.Max_move = 100
	}
	if int(ch.Race) == RACE_ANDROID && PLR_FLAGGED(ch, PLR_SENSEM) {
		for {
			ch.Skills[SKILL_SENSE] = 100
			if true {
				break
			}
		}
		ch.Max_hit += int64(rand_number(400, 500))
		ch.Max_mana += int64(rand_number(400, 500))
		ch.Max_move += int64(rand_number(400, 500))
	}
	ch.Hit = ch.Max_hit
	ch.Mana = ch.Max_mana
	ch.Move = ch.Max_move
	ch.Basepl = ch.Max_hit
	ch.Baseki = ch.Max_mana
	ch.Basest = ch.Max_move
	if int(ch.Real_abils.Str) > 20 {
		ch.Real_abils.Str = 20
	}
	if int(ch.Real_abils.Str) < 8 {
		ch.Real_abils.Str = 8
	}
	if int(ch.Real_abils.Con) > 20 {
		ch.Real_abils.Con = 20
	}
	if int(ch.Real_abils.Con) < 8 {
		ch.Real_abils.Con = 8
	}
	if int(ch.Real_abils.Intel) > 20 {
		ch.Real_abils.Intel = 20
	}
	if int(ch.Real_abils.Intel) < 8 {
		ch.Real_abils.Intel = 8
	}
	if int(ch.Real_abils.Cha) > 20 {
		ch.Real_abils.Cha = 20
	}
	if int(ch.Real_abils.Cha) < 8 {
		ch.Real_abils.Cha = 8
	}
	if int(ch.Real_abils.Dex) > 20 {
		ch.Real_abils.Dex = 20
	}
	if int(ch.Real_abils.Dex) < 8 {
		ch.Real_abils.Dex = 8
	}
	if int(ch.Real_abils.Wis) > 20 {
		ch.Real_abils.Wis = 20
	}
	if int(ch.Real_abils.Wis) < 8 {
		ch.Real_abils.Wis = 8
	}
	ch.Transclass = rand_number(1, 3)
	if config_info.Operation.Siteok_everyone != 0 {
		SET_BIT_AR(ch.Act[:], PLR_SITEOK)
	}
	if int(ch.Race) == RACE_SAIYAN && rand_number(1, 100) >= 95 {
		SET_BIT_AR(ch.Act[:], PLR_LSSJ)
		write_to_output(ch.Desc, libc.CString("@GYou were one of the few born a Legendary Super Saiyan!@n\r\n"))
	}
	ch.Player_specials.Olc_zone = -1
	save_char(ch)
}

var free_start_feats_wizard [3]int = [3]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_SCRIBE_SCROLL, 0}
var free_start_feats_sorcerer [2]int = [2]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, 0}
var free_start_feats_cleric [6]int = [6]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_HEAVY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_MEDIUM, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_rogue [3]int = [3]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_LIGHT, 0}
var free_start_feats_fighter [7]int = [7]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_MARTIAL_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_HEAVY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_MEDIUM, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_paladin [7]int = [7]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_MARTIAL_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_HEAVY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_MEDIUM, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_barbarian [6]int = [6]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_MARTIAL_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_MEDIUM, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_bard [4]int = [4]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_ranger [5]int = [5]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_MARTIAL_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var free_start_feats_monk [4]int = [4]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_MARTIAL_WEAPON_PROFICIENCY, FEAT_IMPROVED_GRAPPLE, 0}
var free_start_feats_druid [5]int = [5]int{FEAT_SIMPLE_WEAPON_PROFICIENCY, FEAT_ARMOR_PROFICIENCY_LIGHT, FEAT_ARMOR_PROFICIENCY_MEDIUM, FEAT_ARMOR_PROFICIENCY_SHIELD, 0}
var no_free_start_feats [1]int = [1]int{}
var free_start_feats [26]*int = [26]*int{&free_start_feats_wizard[0], &free_start_feats_cleric[0], &free_start_feats_rogue[0], &free_start_feats_fighter[0], &free_start_feats_monk[0], &free_start_feats_paladin[0], &free_start_feats_sorcerer[0], &free_start_feats_druid[0], &free_start_feats_bard[0], &free_start_feats_ranger[0], &free_start_feats_barbarian[0], &no_free_start_feats[0], &no_free_start_feats[0], &free_start_feats_wizard[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0], &no_free_start_feats[0]}

func advance_level(ch *char_data, whichclass int) {
	var (
		llog            *levelup_data
		add_hp          int64 = 0
		add_move        int64 = 0
		add_mana        int64 = 0
		add_ki          int64 = 0
		add_prac        int   = 1
		add_train       int
		i               int
		j               int = 0
		ranks           int
		add_gen_feats   int = 0
		add_class_feats int = 0
		buf             [64936]byte
	)
	if whichclass < 0 || whichclass >= NUM_CLASSES {
		basic_mud_log(libc.CString("Invalid class %d passed to advance_level, resetting."), whichclass)
		whichclass = 0
	}
	if config_info.Advance.Allow_multiclass == 0 && whichclass != int(ch.Chclass) {
		basic_mud_log(libc.CString("Attempt to gain a second class without multiclass enabled for %s"), GET_NAME(ch))
		whichclass = int(ch.Chclass)
	}
	ranks = (ch.Chclasses[whichclass]) + (ch.Epicclasses[whichclass])
	llog = new(levelup_data)
	llog.Next = ch.Level_info
	llog.Prev = nil
	if llog.Next != nil {
		llog.Next.Prev = llog
	}
	ch.Level_info = llog
	llog.Skills = func() *level_learn_entry {
		p := &llog.Feats
		llog.Feats = nil
		return *p
	}()
	llog.Type = LEVELTYPE_CLASS
	llog.Spec = int8(whichclass)
	llog.Level = int8(GET_LEVEL(ch))
	switch ranks {
	case 1:
		for i = 0; (func() int {
			j = *(*int)(unsafe.Add(unsafe.Pointer(free_start_feats[whichclass]), unsafe.Sizeof(int(0))*uintptr(i)))
			return j
		}()) != 0; i++ {
			ch.Feats[j] = 1
		}
	case 2:
		switch whichclass {
		case CLASS_ROSHI:
		case CLASS_PICCOLO:
		}
	case 6:
		switch whichclass {
		case CLASS_ROSHI:
		}
	case 9:
		switch whichclass {
		case CLASS_ROSHI:
		}
	}
	if ch.Level == 1 && ch.Race_level < 2 {
		ch.Race_level = 0
		ch.Saving_throw[SAVING_FORTITUDE] = 0
		ch.Saving_throw[SAVING_REFLEX] = 0
		ch.Saving_throw[SAVING_WILL] = 0
	}
	if ranks >= int(LVL_EPICSTART*20) {
		j = ranks - 20
		switch whichclass {
		case CLASS_ROSHI:
			fallthrough
		case CLASS_JINTO:
			fallthrough
		case CLASS_DRAGON_DISCIPLE:
			fallthrough
		case CLASS_ELDRITCH_KNIGHT:
			fallthrough
		case CLASS_HORIZON_WALKER:
			if (j % 4) == 0 {
				add_class_feats++
			}
		case CLASS_PICCOLO:
			if (j % 2) == 0 {
				add_class_feats++
			}
		case CLASS_KRANE:
			if (j % 5) == 0 {
				add_class_feats++
			}
		case CLASS_NPC_EXPERT:
			fallthrough
		case CLASS_NPC_ADEPT:
			fallthrough
		case CLASS_NPC_COMMONER:
			fallthrough
		case CLASS_NPC_ARISTOCRAT:
			fallthrough
		case CLASS_NPC_WARRIOR:
		case CLASS_DWARVEN_DEFENDER:
			fallthrough
		case CLASS_MYSTIC_THEURGE:
			if (j % 6) == 0 {
				add_class_feats++
			}
		default:
			if (j % 3) == 0 {
				add_class_feats++
			}
		}
	} else {
		switch whichclass {
		case CLASS_ROSHI:
			if ranks == 1 || (ranks%2) == 0 {
				add_class_feats++
			}
		case CLASS_PICCOLO:
			if ranks > 9 && (ranks%3) == 0 {
				add_class_feats++
			}
		case CLASS_KRANE:
			if (ranks % 5) == 0 {
				add_class_feats++
			}
		default:
		}
	}
	if GET_LEVEL(ch) >= 1 {
		var (
			pl_percent  float64 = 10
			ki_percent  float64 = 10
			st_percent  float64 = 10
			prac_reward float64 = float64(ch.Aff_abils.Wis)
		)
		if prac_reward < 10 {
			prac_reward = 10
		}
		if GET_LEVEL(ch) >= 91 {
			pl_percent -= 7.2
			ki_percent -= 7.2
			st_percent -= 7.2
		} else if GET_LEVEL(ch) >= 81 {
			pl_percent -= 5.8
			ki_percent -= 5.8
			st_percent -= 5.8
		} else if GET_LEVEL(ch) >= 71 {
			pl_percent -= 5
			ki_percent -= 5
			st_percent -= 5
		} else if GET_LEVEL(ch) >= 61 {
			pl_percent -= 4
			ki_percent -= 4
			st_percent -= 4
		} else if GET_LEVEL(ch) >= 51 {
			pl_percent -= 3
			ki_percent -= 3
			st_percent -= 3
		} else if GET_LEVEL(ch) >= 41 {
			pl_percent -= 1.7
			ki_percent -= 1.7
			st_percent -= 1.7
		} else if GET_LEVEL(ch) >= 31 {
			pl_percent -= 0.8
			ki_percent -= 0.8
			st_percent -= 0.8
		} else if GET_LEVEL(ch) >= 21 {
			pl_percent -= 0.35
			ki_percent -= 0.35
			st_percent -= 0.35
		} else if GET_LEVEL(ch) >= 11 {
			pl_percent -= 0.15
			ki_percent -= 0.15
			st_percent -= 0.15
		}
		switch ch.Race {
		case RACE_HUMAN:
			ki_percent += 3
		case RACE_SAIYAN:
			fallthrough
		case RACE_MAJIN:
			pl_percent += 3
		case RACE_MUTANT:
			st_percent += 3
		case RACE_HALFBREED:
			pl_percent += 1.5
			ki_percent += 1.5
			prac_reward -= prac_reward * 0.4
		case RACE_TRUFFLE:
			prac_reward += prac_reward * 0.5
		}
		if int(ch.Race) != RACE_HUMAN {
			add_hp = int64((float64(ch.Basepl) * 0.01) * pl_percent)
		} else if int(ch.Race) == RACE_HUMAN {
			add_hp = int64(((float64(ch.Basepl) * 0.01) * pl_percent) * 0.8)
		}
		add_mana = int64((float64(ch.Baseki) * 0.01) * ki_percent)
		add_move = int64((float64(ch.Basest) * 0.01) * st_percent)
		add_prac = int(prac_reward + float64(ch.Aff_abils.Intel))
	}
	if add_hp >= 300000 && add_hp < 600000 {
		add_hp += int64(float64(add_hp) * 0.75)
		if add_hp < 300000 {
			add_hp = int64(rand_number(300000, 330000))
		}
	} else if add_hp >= 600000 && add_hp < 1000000 {
		add_hp += int64(float64(add_hp) * 0.7)
		if add_hp < 600000 {
			add_hp = int64(rand_number(600000, 650000))
		}
	} else if add_hp >= 1000000 && add_hp < 2000000 {
		add_hp += int64(float64(add_hp) * 0.65)
		if add_hp < 1000000 {
			add_hp = int64(rand_number(1000000, 1250000))
		}
	} else if add_hp >= 2000000 {
		add_hp += int64(float64(add_hp) * 0.45)
		if add_hp < 2000000 {
			add_hp = int64(rand_number(2000000, 2250000))
		}
	}
	if add_hp >= 15000000 {
		add_hp = 15000000
	}
	if add_move >= 300000 && add_move < 600000 {
		add_move += int64(float64(add_move) * 0.75)
		if add_move < 300000 {
			add_move = int64(rand_number(300000, 330000))
		}
	} else if add_move >= 600000 && add_move < 1000000 {
		add_move += int64(float64(add_move) * 0.7)
		if add_move < 600000 {
			add_move = int64(rand_number(600000, 650000))
		}
	} else if add_move >= 1000000 && add_move < 2000000 {
		add_move += int64(float64(add_move) * 0.65)
		if add_move < 1000000 {
			add_move = int64(rand_number(1000000, 1250000))
		}
	} else if add_move >= 2000000 {
		add_move += int64(float64(add_move) * 0.45)
		if add_move < 2000000 {
			add_move = int64(rand_number(2000000, 2250000))
		}
	}
	if add_move >= 15000000 {
		add_move = 15000000
	}
	if add_mana >= 300000 && add_mana < 600000 {
		add_mana += int64(float64(add_mana) * 0.75)
		if add_mana < 300000 {
			add_mana = int64(rand_number(300000, 330000))
		}
	} else if add_mana >= 600000 && add_mana < 1000000 {
		add_mana += int64(float64(add_mana) * 0.7)
		if add_mana < 600000 {
			add_mana = int64(rand_number(600000, 650000))
		}
	} else if add_mana >= 1000000 && add_mana < 2000000 {
		add_mana += int64(float64(add_mana) * 0.65)
		if add_mana < 1000000 {
			add_mana = int64(rand_number(1000000, 1250000))
		}
	} else if add_mana >= 2000000 {
		add_mana += int64(float64(add_mana) * 0.45)
		if add_mana < 2000000 {
			add_mana = int64(rand_number(2000000, 2250000))
		}
	}
	if add_mana >= 15000000 {
		add_mana = 15000000
	}
	switch GET_LEVEL(ch) {
	case 5:
		add_hp += int64(rand_number(600, 1000))
		add_move += int64(rand_number(600, 1000))
		add_mana += int64(rand_number(600, 1000))
	case 10:
		add_hp += int64(rand_number(5000, 8000))
		add_move += int64(rand_number(5000, 8000))
		add_mana += int64(rand_number(5000, 8000))
	case 20:
		add_hp += int64(rand_number(15000, 18000))
		add_move += int64(rand_number(15000, 18000))
		add_mana += int64(rand_number(15000, 18000))
	case 30:
		add_hp += int64(rand_number(20000, 30000))
		add_move += int64(rand_number(20000, 30000))
		add_mana += int64(rand_number(20000, 30000))
	case 40:
		add_hp += int64(rand_number(50000, 60000))
		add_move += int64(rand_number(50000, 60000))
		add_mana += int64(rand_number(50000, 60000))
	case 50:
		add_hp += int64(rand_number(60000, 70000))
		add_move += int64(rand_number(60000, 70000))
		add_mana += int64(rand_number(60000, 70000))
	case 60:
		add_hp += int64(rand_number(80000, 100000))
		add_move += int64(rand_number(80000, 100000))
		add_mana += int64(rand_number(80000, 100000))
	case 70:
		add_hp += int64(rand_number(100000, 150000))
		add_move += int64(rand_number(100000, 150000))
		add_mana += int64(rand_number(100000, 150000))
	case 80:
		add_hp += int64(rand_number(150000, 200000))
		add_move += int64(rand_number(150000, 200000))
		add_mana += int64(rand_number(150000, 200000))
	case 90:
		add_hp += int64(rand_number(150000, 200000))
		add_move += int64(rand_number(150000, 200000))
		add_mana += int64(rand_number(150000, 200000))
	case 99:
		add_hp += int64(rand_number(1500000, 2000000))
		add_move += int64(rand_number(1500000, 2000000))
		add_mana += int64(rand_number(1500000, 2000000))
	case 100:
		add_hp += int64(rand_number(5000000, 6000000))
		add_move += int64(rand_number(5000000, 6000000))
		add_mana += int64(rand_number(5000000, 6000000))
	}
	if GET_LEVEL(ch) == 1 || (GET_LEVEL(ch)%3) == 0 {
		add_gen_feats += 1
	}
	if int(ch.Race) == RACE_HUMAN {
		add_prac += 2
	}
	i = int(ability_mod_value(int(ch.Aff_abils.Con)))
	if GET_LEVEL(ch) > 1 {
	} else {
		ch.Max_hit += int64(rand_number(1, 20))
		if ch.Max_hit > 250 {
			ch.Max_hit = 250
		}
		if ch.Max_mana > 250 {
			ch.Max_mana = 250
		}
		if ch.Max_move > 250 {
			ch.Max_move = 250
		}
		ch.Basepl = ch.Max_hit
		ch.Baseki = ch.Max_mana
		ch.Basest = ch.Max_hit
		add_prac = 5
		if PLR_FLAGGED(ch, PLR_SKILLP) {
			REMOVE_BIT_AR(ch.Act[:], PLR_SKILLP)
			add_prac *= 5
		} else {
			add_prac *= 2
		}
	}
	llog.Hp_roll = int8(j)
	if rand_number(1, 8) == 2 {
		add_train = 1
		if add_train != 0 {
			ch.Player_specials.Ability_trains += add_train
		}
	}
	if rand_number(1, 4) == 4 {
		send_to_char(ch, libc.CString("@D[@mPractice Session Bonus!@D]@n\r\n"))
		add_prac += rand_number(4, 12)
	}
	if (int(ch.Race) == RACE_DEMON || int(ch.Race) == RACE_KANASSAN) && GET_LEVEL(ch) > 80 {
		add_hp *= 2
		add_mana *= 2
		add_move *= 2
	} else if (int(ch.Race) == RACE_DEMON || int(ch.Race) == RACE_KANASSAN) && GET_LEVEL(ch) > 60 {
		add_hp += int64(float64(add_hp) * 1.75)
		add_mana += int64(float64(add_mana) * 1.75)
		add_move += int64(float64(add_move) * 1.75)
	} else if (int(ch.Race) == RACE_DEMON || int(ch.Race) == RACE_KANASSAN) && GET_LEVEL(ch) > 50 {
		add_hp += int64(float64(add_hp) * 1.5)
		add_mana += int64(float64(add_mana) * 1.5)
		add_move += int64(float64(add_move) * 1.5)
	} else if (int(ch.Race) == RACE_DEMON || int(ch.Race) == RACE_KANASSAN) && GET_LEVEL(ch) > 40 {
		add_hp += int64(float64(add_hp) * 1.25)
		add_mana += int64(float64(add_mana) * 1.25)
		add_move += int64(float64(add_move) * 1.25)
	}
	llog.Mana_roll = int8(add_mana)
	llog.Move_roll = int8(add_move)
	llog.Ki_roll = int8(add_ki)
	llog.Add_skill = int8(add_prac)
	ch.Player_specials.Class_skill_points[whichclass] += add_prac
	ch.Basepl += add_hp
	ch.Baseki += add_mana
	ch.Basest += add_move
	var nhp int = int(add_hp)
	var nma int = int(add_mana)
	var nmo int = int(add_move)
	if int(ch.Race) == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS1) {
		add_hp *= 3
		add_move *= 3
		add_mana *= 3
	} else if int(ch.Race) == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS2) {
		add_hp *= 4
		add_move *= 4
		add_mana *= 4
	} else if int(ch.Race) == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS3) {
		add_hp *= 5
		add_move *= 5
		add_mana *= 5
	} else if int(ch.Race) == RACE_HOSHIJIN && ch.Starphase == 1 {
		add_hp *= 2
		add_move *= 2
		add_mana *= 2
	} else if int(ch.Race) == RACE_HOSHIJIN && ch.Starphase == 2 {
		add_hp *= 3
		add_move *= 3
		add_mana *= 3
	} else if int(ch.Race) == RACE_BIO && PLR_FLAGGED(ch, PLR_TRANS1) {
		add_hp *= 2
		add_move *= 2
		add_mana *= 2
	} else if int(ch.Race) == RACE_BIO && PLR_FLAGGED(ch, PLR_TRANS2) {
		add_hp *= 3
		add_move *= 3
		add_mana *= 3
	} else if int(ch.Race) == RACE_BIO && PLR_FLAGGED(ch, PLR_TRANS3) {
		add_hp += int64(float64(add_hp) * 3.5)
		add_move += int64(float64(add_move) * 3.5)
		add_mana += int64(float64(add_mana) * 3.5)
	} else if int(ch.Race) == RACE_BIO && PLR_FLAGGED(ch, PLR_TRANS4) {
		add_hp *= 4
		add_move *= 4
		add_mana *= 4
	} else if int(ch.Race) == RACE_MAJIN && PLR_FLAGGED(ch, PLR_TRANS1) {
		add_hp *= 2
		add_move *= 2
		add_mana *= 2
	} else if int(ch.Race) == RACE_MAJIN && PLR_FLAGGED(ch, PLR_TRANS2) {
		add_hp *= 3
		add_move *= 3
		add_mana *= 3
	} else if int(ch.Race) == RACE_MAJIN && PLR_FLAGGED(ch, PLR_TRANS3) {
		add_hp += int64(float64(add_hp) * 4.5)
		add_move += int64(float64(add_move) * 4.5)
		add_mana += int64(float64(add_mana) * 4.5)
	}
	ch.Max_hit += add_hp
	ch.Max_move += add_move
	ch.Max_mana += add_mana
	add_hp = int64(nhp)
	add_mana = int64(nma)
	add_move = int64(nmo)
	if ch.Level >= LVL_EPICSTART {
		ch.Player_specials.Epic_feat_points += add_gen_feats
		llog.Add_epic_feats = int8(add_gen_feats)
		ch.Player_specials.Epic_class_feat_points[whichclass] += add_class_feats
		llog.Add_class_epic_feats = int8(add_class_feats)
	} else {
		ch.Player_specials.Feat_points += add_gen_feats
		llog.Add_gen_feats = int8(add_gen_feats)
		ch.Player_specials.Class_feat_points[whichclass] += add_class_feats
		llog.Add_class_feats = int8(add_class_feats)
	}
	if ch.Admlevel >= ADMLVL_IMMORT {
		for i = 0; i < 3; i++ {
			ch.Player_specials.Conditions[i] = -1
		}
		SET_BIT_AR(ch.Player_specials.Pref[:], PRF_HOLYLIGHT)
	}
	stdio.Sprintf(&buf[0], "@D[@YGain@D: @RPl@D(@G%s@D) @gSt@D(@G%s@D) @CKi@D(@G%s@D) @bPS@D(@G%s@D)]", add_commas(add_hp), add_commas(add_move), add_commas(add_mana), add_commas(int64(add_prac)))
	if (ch.Bonuses[BONUS_GMEMORY]) != 0 && (GET_LEVEL(ch) == 20 || GET_LEVEL(ch) == 40 || GET_LEVEL(ch) == 60 || GET_LEVEL(ch) == 80 || GET_LEVEL(ch) == 100) {
		ch.Skill_slots += 1
		send_to_char(ch, libc.CString("@CYou feel like you could remember a new skill!@n\r\n"))
	}
	if int(ch.Race) == RACE_NAMEK && rand_number(1, 100) <= 5 {
		ch.Skill_slots += 1
		send_to_char(ch, libc.CString("@CYou feel as though you could learn another skill.@n\r\n"))
	}
	if int(ch.Race) == RACE_ICER && rand_number(1, 100) <= 25 {
		bring_to_cap(ch)
		send_to_char(ch, libc.CString("@GYou feel your body obtain its current optimal strength!@n\r\n"))
	}
	var gain_stat int = FALSE
	switch GET_LEVEL(ch) {
	case 10:
		fallthrough
	case 20:
		fallthrough
	case 30:
		fallthrough
	case 40:
		fallthrough
	case 50:
		fallthrough
	case 60:
		fallthrough
	case 70:
		fallthrough
	case 80:
		fallthrough
	case 90:
		fallthrough
	case 100:
		gain_stat = TRUE
	}
	if gain_stat == TRUE {
		var (
			raise     int = FALSE
			stat_fail int = 0
		)
		if int(ch.Race) == RACE_KONATSU {
			for raise == FALSE {
				if int(ch.Real_abils.Dex) < 100 && rand_number(1, 2) == 2 && stat_fail != 1 {
					if int(ch.Real_abils.Dex) < 45 || (ch.Bonuses[BONUS_CLUMSY]) <= 0 {
						ch.Real_abils.Dex += 1
						send_to_char(ch, libc.CString("@GYou feel your agility increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 1
					}
				} else if int(ch.Real_abils.Cha) < 100 && raise == FALSE && stat_fail < 2 {
					if int(ch.Real_abils.Cha) < 45 || (ch.Bonuses[BONUS_SLOW]) > 0 {
						ch.Real_abils.Cha += 1
						send_to_char(ch, libc.CString("@GYou feel your speed increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 2
					}
				} else if stat_fail == 3 {
					send_to_char(ch, libc.CString("@RBoth agility and speed are capped!@n"))
					raise = TRUE
				}
			}
		} else if int(ch.Race) == RACE_MUTANT {
			for raise == FALSE {
				if int(ch.Real_abils.Con) < 100 && rand_number(1, 2) == 2 && stat_fail != 1 {
					if int(ch.Real_abils.Con) < 45 || (ch.Bonuses[BONUS_FRAIL]) <= 0 {
						ch.Real_abils.Con += 1
						send_to_char(ch, libc.CString("@GYou feel your constitution increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 1
					}
				} else if int(ch.Real_abils.Cha) < 100 && raise == FALSE && stat_fail < 2 {
					if int(ch.Real_abils.Cha) < 45 || (ch.Bonuses[BONUS_SLOW]) > 0 {
						ch.Real_abils.Cha += 1
						send_to_char(ch, libc.CString("@GYou feel your speed increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 2
					}
				} else if stat_fail == 3 {
					send_to_char(ch, libc.CString("@RBoth constitution and speed are capped!@n"))
					raise = TRUE
				}
			}
		} else if int(ch.Race) == RACE_HOSHIJIN {
			for raise == FALSE {
				if int(ch.Real_abils.Str) < 100 && rand_number(1, 2) == 2 && stat_fail != 1 {
					if int(ch.Real_abils.Str) < 45 || (ch.Bonuses[BONUS_WIMP]) <= 0 {
						ch.Real_abils.Str += 1
						send_to_char(ch, libc.CString("@GYou feel your strength increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 1
					}
				} else if int(ch.Real_abils.Dex) < 100 && raise == FALSE && stat_fail < 2 {
					if int(ch.Real_abils.Dex) < 45 || (ch.Bonuses[BONUS_SLOW]) > 0 {
						ch.Real_abils.Dex += 1
						send_to_char(ch, libc.CString("@GYou feel your agility increase!@n\r\n"))
						raise = TRUE
					} else {
						stat_fail += 2
					}
				} else if stat_fail == 3 {
					send_to_char(ch, libc.CString("@RBoth strength and agility are capped!@n"))
					raise = TRUE
				}
			}
		}
	}
	libc.StrCat(&buf[0], libc.CString(".\r\n"))
	send_to_char(ch, libc.CString("%s"), &buf[0])
	if GET_SKILL(ch, SKILL_POTENTIAL) != 0 && rand_number(1, 4) == 4 {
		send_to_char(ch, libc.CString("You can now perform another Potential Release.\r\n"))
		ch.Boosts += 1
	}
	if int(ch.Race) == RACE_MAJIN && GET_LEVEL(ch) == 25 {
		send_to_char(ch, libc.CString("You can now perform another Majinization.\r\n"))
		ch.Boosts += 1
	}
	if int(ch.Race) == RACE_MAJIN && GET_LEVEL(ch) == 50 {
		send_to_char(ch, libc.CString("You can now perform another Majinization.\r\n"))
		ch.Boosts += 1
	}
	if int(ch.Race) == RACE_MAJIN && GET_LEVEL(ch) == 75 {
		send_to_char(ch, libc.CString("You can now perform another Majinization.\r\n"))
		ch.Boosts += 1
	}
	if int(ch.Race) == RACE_MAJIN && GET_LEVEL(ch) == 100 {
		send_to_char(ch, libc.CString("You can now perform another Majinization.\r\n"))
		ch.Boosts += 1
	}
	switch GET_LEVEL(ch) {
	case 10:
		fallthrough
	case 20:
		fallthrough
	case 30:
		fallthrough
	case 40:
		fallthrough
	case 50:
		fallthrough
	case 60:
		fallthrough
	case 70:
		fallthrough
	case 80:
		fallthrough
	case 90:
		fallthrough
	case 100:
		if (ch.Bonuses[BONUS_BRAWNY]) > 0 {
			ch.Real_abils.Str += 2
			send_to_char(ch, libc.CString("@GYour muscles have grown stronger!@n\r\n"))
		}
		if (ch.Bonuses[BONUS_SCHOLARLY]) > 0 {
			ch.Real_abils.Intel += 2
			send_to_char(ch, libc.CString("@GYour mind has grown sharper!@n\r\n"))
		}
		if (ch.Bonuses[BONUS_SAGE]) > 0 {
			ch.Real_abils.Wis += 2
			send_to_char(ch, libc.CString("@GYour understanding about life has improved!@n\r\n"))
		}
		if (ch.Bonuses[BONUS_AGILE]) > 0 {
			ch.Real_abils.Dex += 2
			send_to_char(ch, libc.CString("@GYour body has grown more agile!@n\r\n"))
		}
		if (ch.Bonuses[BONUS_QUICK]) > 0 {
			ch.Real_abils.Cha += 2
			send_to_char(ch, libc.CString("@GYou feel like your speed has improved!@n\r\n"))
		}
		if (ch.Bonuses[BONUS_STURDY]) > 0 {
			ch.Real_abils.Con += 2
			send_to_char(ch, libc.CString("@GYour body feels tougher now!@n\r\n"))
		}
	}
	if GET_LEVEL(ch) == 1 {
		ch.Armor = 0
	}
	if GET_LEVEL(ch) == 2 {
		ERAPLAYERS += 1
	}
	snoop_check(ch)
	save_char(ch)
}
func invalid_class(ch *char_data, obj *obj_data) int {
	if OBJ_FLAGGED(obj, ITEM_ANTI_WIZARD) && int(ch.Chclass) == CLASS_ROSHI {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_CLERIC) && int(ch.Chclass) == CLASS_PICCOLO {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_FIGHTER) && int(ch.Chclass) == CLASS_NAIL {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ROGUE) && int(ch.Chclass) == CLASS_KRANE {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_MONK) && int(ch.Chclass) == CLASS_BARDOCK {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_MONK) && int(ch.Chclass) != CLASS_BARDOCK {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_WIZARD) && int(ch.Chclass) != CLASS_ROSHI {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_JINTO) && int(ch.Chclass) != CLASS_JINTO {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_CLERIC) && int(ch.Chclass) != CLASS_PICCOLO {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_ROGUE) && int(ch.Chclass) != CLASS_KRANE {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ASSASSIN) && int(ch.Chclass) != CLASS_KURZAK {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_FIGHTER) && int(ch.Chclass) != CLASS_NAIL {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_PALADIN) && int(ch.Chclass) != CLASS_GINYU {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_WIZARD) && int(ch.Chclass) == CLASS_ROSHI {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_PALADIN) && int(ch.Chclass) == CLASS_GINYU {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_BARBARIAN) && int(ch.Chclass) == CLASS_KABITO {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_BARD) && int(ch.Chclass) == CLASS_ANDSIX {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_BARD) && int(ch.Chclass) != CLASS_ANDSIX {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_RANGER) && int(ch.Chclass) == CLASS_DABURA {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_DRUID) && int(ch.Chclass) == CLASS_TAPION {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ARCANE_ARCHER) && int(ch.Chclass) == CLASS_JINTO {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ARCANE_TRICKSTER) && int(ch.Chclass) == CLASS_TSUNA {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ARCHMAGE) && int(ch.Chclass) == CLASS_KURZAK {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_BLACKGUARD) && ((ch.Chclasses[CLASS_BLACKGUARD])+(ch.Epicclasses[CLASS_BLACKGUARD])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_DRAGON_DISCIPLE) && ((ch.Chclasses[CLASS_DRAGON_DISCIPLE])+(ch.Epicclasses[CLASS_DRAGON_DISCIPLE])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_DUELIST) && ((ch.Chclasses[CLASS_DUELIST])+(ch.Epicclasses[CLASS_DUELIST])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_DWARVEN_DEFENDER) && ((ch.Chclasses[CLASS_DWARVEN_DEFENDER])+(ch.Epicclasses[CLASS_DWARVEN_DEFENDER])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ELDRITCH_KNIGHT) && ((ch.Chclasses[CLASS_ELDRITCH_KNIGHT])+(ch.Epicclasses[CLASS_ELDRITCH_KNIGHT])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_HIEROPHANT) && ((ch.Chclasses[CLASS_HIEROPHANT])+(ch.Epicclasses[CLASS_HIEROPHANT])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_HORIZON_WALKER) && ((ch.Chclasses[CLASS_HORIZON_WALKER])+(ch.Epicclasses[CLASS_HORIZON_WALKER])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_LOREMASTER) && ((ch.Chclasses[CLASS_LOREMASTER])+(ch.Epicclasses[CLASS_LOREMASTER])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_MYSTIC_THEURGE) && ((ch.Chclasses[CLASS_MYSTIC_THEURGE])+(ch.Epicclasses[CLASS_MYSTIC_THEURGE])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_SHADOWDANCER) && ((ch.Chclasses[CLASS_SHADOWDANCER])+(ch.Epicclasses[CLASS_SHADOWDANCER])) > 0 {
		return TRUE
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_THAUMATURGIST) && ((ch.Chclasses[CLASS_THAUMATURGIST])+(ch.Epicclasses[CLASS_THAUMATURGIST])) > 0 {
		return TRUE
	}
	return FALSE
}
func level_exp(ch *char_data, level int) int {
	var req int = 1
	switch level {
	case 0:
		req = 0
	case 1:
		req = 1
	case 2:
		req = 999
	case 3:
		req = 2000
	case 4:
		req = 4500
	case 5:
		req = 6000
	case 6:
		req = 10000
	case 7:
		req = 15000
	case 8:
		req = 25000
	case 9:
		req = 35000
	case 10:
		req = 50000
	case 11:
		req = 75000
	case 12:
		req = 100000
	case 13:
		req = 125000
	case 14:
		req = 150000
	case 15:
		req = 175000
	case 16:
		req = 225000
	case 17:
		req = 300000
	case 18:
		req = 350000
	case 19:
		req = 400000
	case 20:
		req = 500000
	case 21:
		req = 650000
	case 22:
		req = 800000
	case 23:
		req = 1000000
	case 24:
		req = 1250000
	case 25:
		req = 1500000
	case 26:
		req = 1750000
	case 27:
		req = 2000000
	case 28:
		req = 2250000
	case 29:
		req = 2500000
	case 30:
		req = 2750000
	case 31:
		req = 3000000
	case 32:
		req = 3250000
	case 33:
		req = 3500000
	case 34:
		req = 3750000
	case 35:
		req = 4000000
	case 36:
		req = 4250000
	case 37:
		req = 4500000
	case 38:
		req = 4750000
	case 39:
		req = 5000000
	case 40:
		req = 5500000
	case 41:
		req = 6000000
	case 42:
		req = 6500000
	case 43:
		req = 7000000
	case 44:
		req = 7500000
	case 45:
		req = 8000000
	case 46:
		req = 8500000
	case 47:
		req = 9000000
	case 48:
		req = 9500000
	case 49:
		req = 10000000
	case 50:
		req = 10500000
	case 51:
		req = 11000000
	case 52:
		req = 11500000
	case 53:
		req = 12000000
	case 54:
		req = 12500000
	case 55:
		req = 13000000
	case 56:
		req = 13500000
	case 57:
		req = 14000000
	case 58:
		req = 14500000
	case 59:
		req = 15000000
	case 60:
		req = 18000000
	case 61:
		req = 25000000
	case 62:
		req = 28000000
	case 63:
		req = 31000000
	case 64:
		req = 34000000
	case 65:
		req = 37000000
	case 66:
		req = 40000000
	case 67:
		req = 43000000
	case 68:
		req = 46000000
	case 69:
		req = 49000000
	case 70:
		req = 52000000
	case 71:
		req = 55000000
	case 72:
		req = 58000000
	case 73:
		req = 61000000
	case 74:
		req = 64000000
	case 75:
		req = 67000000
	case 76:
		req = 70000000
	case 77:
		req = 73000000
	case 78:
		req = 76000000
	case 79:
		req = 79000000
	case 80:
		req = 82000000
	case 81:
		req = 88000000
	case 82:
		req = 94000000
	case 83:
		req = 100000000
	case 84:
		req = 106000000
	case 85:
		req = 112000000
	case 86:
		req = 118000000
	case 87:
		req = 124000000
	case 88:
		req = 130000000
	case 89:
		req = 136000000
	case 90:
		req = 142000000
	case 91:
		req = 150000000
	case 92:
		req = 175000000
	case 93:
		req = 200000000
	case 94:
		req = 225000000
	case 95:
		req = 250000000
	case 96:
		req = 300000000
	case 97:
		req = 400000000
	case 98:
		req = 500000000
	case 99:
		req = 600000000
	case 100:
		req = 800000000
	}
	if int(ch.Race) == RACE_KAI {
		req += int(float64(req) * 0.15)
	}
	return req
}
func ability_mod_value(abil int) int8 {
	return int8((abil / 2) - 5)
}
func dex_mod_capped(ch *char_data) int8 {
	var (
		mod   int8
		armor *obj_data
	)
	mod = ability_mod_value(int(ch.Aff_abils.Dex))
	armor = ch.Equipment[WEAR_BODY]
	if armor != nil && int(armor.Type_flag) == ITEM_ARMOR {
		mod = int8(MIN(int64(mod), int64(armor.Value[VAL_ARMOR_MAXDEXMOD])))
	}
	return mod
}

var cabbr_ranktable [31]int

func comp_rank(a unsafe.Pointer, b unsafe.Pointer) int {
	var (
		first  int
		second int
	)
	first = *(*int)(a)
	second = *(*int)(b)
	return cabbr_ranktable[second] - cabbr_ranktable[first]
}

func total_skill_levels(ch *char_data, skill int) int {
	var (
		i     int = 0
		j     int
		total int = 0
	)
	for i = 0; i < NUM_CLASSES; i++ {
		j = ((ch.Chclasses[i]) + (ch.Epicclasses[i])) + 1 - spell_info[skill].Min_level[i]
		if j > 0 {
			total += j
		}
	}
	return total
}
func load_levels() int {
	var (
		fp        *stdio.File
		line      [256]byte
		sect_name [256]byte = [256]byte{0: '\x00'}
		ptr       *byte
		linenum   int = 0
		tp        int
		cls       int
		sect_type int = -1
	)
	if (func() *stdio.File {
		fp = stdio.FOpen(LIB_ETC, "r")
		return fp
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Could not open level configuration file, error: %s!"), libc.StrError(libc.Errno))
		return -1
	}
	for cls = 0; cls < NUM_CLASSES; cls++ {
		for tp = 0; tp <= SAVING_WILL; tp++ {
			save_classes[tp][cls] = 0
		}
		basehit_classes[cls] = 0
	}
	for {
		linenum++
		if fp.GetS(&line[0], READ_SIZE) == nil {
			basic_mud_log(libc.CString("SYSERR: Unexpected EOF in file %s."), LIB_ETC)
			return -1
		} else if line[0] == '$' {
			break
		} else if line[0] == '*' {
			continue
		} else if line[0] == '#' {
			if (func() int {
				tp = stdio.Sscanf(&line[0], "#%s", &sect_name[0])
				return tp
			}()) != 1 {
				basic_mud_log(libc.CString("SYSERR: Format error in file %s, line number %d - text: %s."), LIB_ETC, linenum, &line[0])
				return -1
			} else if (func() int {
				sect_type = search_block(&sect_name[0], &config_sect[0], FALSE)
				return sect_type
			}()) == -1 {
				basic_mud_log(libc.CString("SYSERR: Invalid section in file %s, line number %d: %s."), LIB_ETC, linenum, &sect_name[0])
				return -1
			}
		} else {
			if sect_type == CONFIG_LEVEL_VERSION {
				if libc.StrNCmp(&line[0], libc.CString("Suntzu"), 6) == 0 {
					basic_mud_log(libc.CString("SYSERR: Suntzu %s config files are not compatible with rasputin"), LIB_ETC)
					return -1
				} else {
					libc.StrCpy(&level_version[0], &line[0])
				}
			} else if sect_type == CONFIG_LEVEL_VERNUM {
				level_vernum = libc.Atoi(libc.GoString(&line[0]))
			} else if sect_type == CONFIG_LEVEL_EXPERIENCE {
				tp = libc.Atoi(libc.GoString(&line[0]))
				exp_multiplier = float32(tp)
			} else if sect_type >= CONFIG_LEVEL_FORTITUDE && sect_type <= CONFIG_LEVEL_WILL || sect_type == CONFIG_LEVEL_BASEHIT {
				for ptr = &line[0]; ptr != nil && *ptr != 0 && !unicode.IsDigit(rune(*ptr)); ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
				}
				if ptr == nil || *ptr == 0 || !unicode.IsDigit(rune(*ptr)) {
					basic_mud_log(libc.CString("SYSERR: Cannot find class number in file %s, line number %d, section %s."), LIB_ETC, linenum, &sect_name[0])
					return -1
				}
				cls = libc.Atoi(libc.GoString(ptr))
				for ; ptr != nil && *ptr != 0 && unicode.IsDigit(rune(*ptr)); ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
				}
				for ; ptr != nil && *ptr != 0 && !unicode.IsDigit(rune(*ptr)); ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
				}
				if ptr != nil && *ptr != 0 && !unicode.IsDigit(rune(*ptr)) {
					basic_mud_log(libc.CString("SYSERR: Non-numeric entry in file %s, line number %d, section %s."), LIB_ETC, linenum, &sect_name[0])
					return -1
				}
				if ptr != nil && *ptr != 0 {
					tp = libc.Atoi(libc.GoString(ptr))
				} else {
					basic_mud_log(libc.CString("SYSERR: Need 1 value in %s, line number %d, section %s."), LIB_ETC, linenum, &sect_name[0])
					return -1
				}
				if cls < 0 || cls >= NUM_CLASSES {
					basic_mud_log(libc.CString("SYSERR: Invalid class number %d in file %s, line number %d."), cls, LIB_ETC, linenum)
					return -1
				} else {
					if sect_type == CONFIG_LEVEL_BASEHIT {
						basehit_classes[cls] = tp
					} else {
						save_classes[SAVING_FORTITUDE+sect_type-CONFIG_LEVEL_FORTITUDE][cls] = tp
					}
				}
			} else {
				basic_mud_log(libc.CString("Unsupported level config option"))
			}
		}
	}
	fp.Close()
	for cls = 0; cls < NUM_CLASSES; cls++ {
		basic_mud_log(libc.CString("Base hit for class %s: %s"), class_names[cls], basehit_type_names[basehit_classes[cls]])
	}
	for cls = 0; cls < NUM_CLASSES; cls++ {
		basic_mud_log(libc.CString("Saves for class %s: fort=%s, reflex=%s, will=%s"), class_names[cls], save_type_names[save_classes[SAVING_FORTITUDE][cls]], save_type_names[save_classes[SAVING_REFLEX][cls]], save_type_names[save_classes[SAVING_WILL][cls]])
	}
	return 0
}
func highest_skill_value(level int, type_ int) int {
	if level >= 60 {
		return 100
	} else if level >= 20 {
		return level + 40
	} else if level >= 10 {
		return level + 30
	} else if level >= 1 {
		return level + 25
	} else {
		return 0
	}
}
func calc_penalty_exp(ch *char_data, gain int) int {
	return gain
}

var size_scaling_table [9][4]int = [9][4]int{{-10, -2, -2, 0}, {-10, -2, -2, 0}, {-8, -2, -2, 0}, {-4, -2, -2, 0}, {}, {8, -2, 4, 2}, {16, -4, 8, 5}, {24, -4, 12, 9}, {32, -4, 16, 14}}

func birth_age(ch *char_data) libc.Time {
	var tmp int
	tmp = rand_number(16, 18)
	return libc.Time(tmp)
}
func max_age(ch *char_data) libc.Time {
	var aging *aging_data
	_ = aging
	var tmp uint64
	if ch.Time.Maxage != 0 {
		return ch.Time.Maxage - ch.Time.Birth
	}
	aging = &racial_aging_data[ch.Race]
	tmp = 120
	return libc.Time(uint32(tmp))
}

var class_feats_wizard [1]int = [1]int{FEAT_UNDEFINED}
var class_feats_rogue [1]int = [1]int{FEAT_UNDEFINED}
var class_feats_fighter [1]int = [1]int{FEAT_UNDEFINED}
var no_class_feats [1]int = [1]int{FEAT_UNDEFINED}
var class_bonus_feats [31]*int = [31]*int{0: &class_feats_wizard[0], 1: &no_class_feats[0], 2: &class_feats_rogue[0], 3: &class_feats_fighter[0], 4: &no_class_feats[0], 5: &no_class_feats[0], 6: &no_class_feats[0], 7: &no_class_feats[0], 8: &no_class_feats[0], 9: &no_class_feats[0], 10: &no_class_feats[0]}
