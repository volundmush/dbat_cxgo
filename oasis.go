package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const _OASISOLC = 518
const AEDIT_PERMISSION = 999
const HEDIT_PERMISSION = 888
const NUM_ZONE_FLAGS = 36
const NUM_GENDERS = 3
const NUM_SHOP_FLAGS = 3
const MAX_ROOM_NAME = 100
const MAX_MOB_NAME = 50
const MAX_OBJ_NAME = 50
const MAX_ROOM_DESC = 4096
const MAX_EXIT_DESC = 256
const MAX_EXTRA_DESC = 512
const MAX_MOB_DESC = 1024
const MAX_OBJ_DESC = 512
const MAX_DUPLICATES = 10000
const MAX_FROM_ROOM = 50
const MAX_WEAPON_SDICE = 50
const MAX_WEAPON_NDICE = 50
const MAX_OBJ_WEIGHT = 100000000
const MAX_OBJ_COST = 2000000
const MAX_OBJ_RENT = 2000000
const MAX_CONTAINER_SIZE = 100000000
const MAX_MOB_GOLD = 1000000
const MAX_MOB_EXP = 1500000
const MAX_OBJ_TIMER = 1071000
const BIT_STRING_LENGTH = 33
const OASIS_WLD = 0
const OASIS_MOB = 1
const OASIS_OBJ = 2
const OASIS_ZON = 3
const OASIS_EXI = 4
const OASIS_CFG = 5

const CLEANUP_ALL = 1
const CLEANUP_STRUCTS = 2
const CLEANUP_CONFIG = 3
const AEDIT_CONFIRM_SAVESTRING = 0
const AEDIT_CONFIRM_EDIT = 1
const AEDIT_CONFIRM_ADD = 2
const AEDIT_MAIN_MENU = 3
const AEDIT_ACTION_NAME = 4
const AEDIT_SORT_AS = 5
const AEDIT_MIN_CHAR_POS = 6
const AEDIT_MIN_VICT_POS = 7
const AEDIT_HIDDEN_FLAG = 8
const AEDIT_MIN_CHAR_LEVEL = 9
const AEDIT_NOVICT_CHAR = 10
const AEDIT_NOVICT_OTHERS = 11
const AEDIT_VICT_CHAR_FOUND = 12
const AEDIT_VICT_OTHERS_FOUND = 13
const AEDIT_VICT_VICT_FOUND = 14
const AEDIT_VICT_NOT_FOUND = 15
const AEDIT_SELF_CHAR = 16
const AEDIT_SELF_OTHERS = 17
const AEDIT_VICT_CHAR_BODY_FOUND = 18
const AEDIT_VICT_OTHERS_BODY_FOUND = 19
const AEDIT_VICT_VICT_BODY_FOUND = 20
const AEDIT_OBJ_CHAR_FOUND = 21
const AEDIT_OBJ_OTHERS_FOUND = 22
const OEDIT_MAIN_MENU = 1
const OEDIT_EDIT_NAMELIST = 2
const OEDIT_SHORTDESC = 3
const OEDIT_LONGDESC = 4
const OEDIT_ACTDESC = 5
const OEDIT_TYPE = 6
const OEDIT_EXTRAS = 7
const OEDIT_WEAR = 8
const OEDIT_WEIGHT = 9
const OEDIT_COST = 10
const OEDIT_COSTPERDAY = 11
const OEDIT_TIMER = 12
const OEDIT_VALUE_1 = 13
const OEDIT_VALUE_2 = 14
const OEDIT_VALUE_3 = 15
const OEDIT_VALUE_4 = 16
const OEDIT_APPLY = 17
const OEDIT_APPLYMOD = 18
const OEDIT_EXTRADESC_KEY = 19
const OEDIT_CONFIRM_SAVEDB = 20
const OEDIT_CONFIRM_SAVESTRING = 21
const OEDIT_PROMPT_APPLY = 22
const OEDIT_EXTRADESC_DESCRIPTION = 23
const OEDIT_EXTRADESC_MENU = 24
const OEDIT_LEVEL = 25
const OEDIT_PERM = 26
const OEDIT_VALUE_5 = 27
const OEDIT_VALUE_6 = 28
const OEDIT_VALUE_7 = 29
const OEDIT_VALUE_8 = 30
const OEDIT_MATERIAL = 31
const OEDIT_VALUE_9 = 32
const OEDIT_VALUE_10 = 33
const OEDIT_VALUE_11 = 34
const OEDIT_VALUE_12 = 35
const OEDIT_VALUE_13 = 36
const OEDIT_VALUE_14 = 37
const OEDIT_VALUE_15 = 38
const OEDIT_VALUE_16 = 39
const OEDIT_SIZE = 40
const OEDIT_APPLYSPEC = 41
const OEDIT_PROMPT_SPELLBOOK = 42
const OEDIT_SPELLBOOK = 43
const OEDIT_COPY = 44
const OEDIT_DELETE = 45
const REDIT_MAIN_MENU = 1
const REDIT_NAME = 2
const REDIT_DESC = 3
const REDIT_FLAGS = 4
const REDIT_SECTOR = 5
const REDIT_EXIT_MENU = 6
const REDIT_CONFIRM_SAVEDB = 7
const REDIT_CONFIRM_SAVESTRING = 8
const REDIT_EXIT_NUMBER = 9
const REDIT_EXIT_DESCRIPTION = 10
const REDIT_EXIT_KEYWORD = 11
const REDIT_EXIT_KEY = 12
const REDIT_EXIT_DOORFLAGS = 13
const REDIT_EXTRADESC_MENU = 14
const REDIT_EXTRADESC_KEY = 15
const REDIT_EXTRADESC_DESCRIPTION = 16
const REDIT_DELETE = 17
const REDIT_EXIT_DCLOCK = 18
const REDIT_EXIT_DCHIDE = 19
const REDIT_EXIT_DCSKILL = 20
const REDIT_EXIT_DCMOVE = 21
const REDIT_EXIT_SAVETYPE = 22
const REDIT_EXIT_DCSAVE = 23
const REDIT_EXIT_FAILROOM = 24
const REDIT_EXIT_TOTALFAILROOM = 25
const REDIT_COPY = 26
const ZEDIT_MAIN_MENU = 0
const ZEDIT_DELETE_ENTRY = 1
const ZEDIT_NEW_ENTRY = 2
const ZEDIT_CHANGE_ENTRY = 3
const ZEDIT_COMMAND_TYPE = 4
const ZEDIT_IF_FLAG = 5
const ZEDIT_ARG1 = 6
const ZEDIT_ARG2 = 7
const ZEDIT_ARG3 = 8
const ZEDIT_ARG4 = 9
const ZEDIT_ARG5 = 10
const ZEDIT_ZONE_NAME = 11
const ZEDIT_ZONE_LIFE = 12
const ZEDIT_ZONE_BOT = 13
const ZEDIT_ZONE_TOP = 14
const ZEDIT_ZONE_RESET = 15
const ZEDIT_CONFIRM_SAVESTRING = 16
const ZEDIT_ZONE_BUILDERS = 17
const ZEDIT_SARG1 = 18
const ZEDIT_SARG2 = 19
const ZEDIT_ZONE_FLAGS = 20
const ZEDIT_MIN_LEVEL = 21
const ZEDIT_MAX_LEVEL = 22
const MEDIT_MAIN_MENU = 0
const MEDIT_ALIAS = 1
const MEDIT_S_DESC = 2
const MEDIT_L_DESC = 3
const MEDIT_D_DESC = 4
const MEDIT_NPC_FLAGS = 5
const MEDIT_AFF_FLAGS = 6
const MEDIT_CONFIRM_SAVESTRING = 7
const MEDIT_NUMERICAL_RESPONSE = 10
const MEDIT_SEX = 11
const MEDIT_ACCURACY = 12
const MEDIT_DAMAGE = 13
const MEDIT_NDD = 14
const MEDIT_SDD = 15
const MEDIT_NUM_HP_DICE = 16
const MEDIT_SIZE_HP_DICE = 17
const MEDIT_ADD_HP = 18
const MEDIT_AC = 19
const MEDIT_EXP = 20
const MEDIT_GOLD = 21
const MEDIT_POS = 22
const MEDIT_DEFAULT_POS = 23
const MEDIT_ATTACK = 24
const MEDIT_LEVEL = 25
const MEDIT_ALIGNMENT = 26
const MEDIT_CLASS = 33
const MEDIT_RACE = 34
const MEDIT_SIZE = 35
const MEDIT_COPY = 36
const MEDIT_DELETE = 37
const MEDIT_PERSONALITY = 38
const SEDIT_MAIN_MENU = 0
const SEDIT_CONFIRM_SAVESTRING = 1
const SEDIT_NOITEM1 = 2
const SEDIT_NOITEM2 = 3
const SEDIT_NOCASH1 = 4
const SEDIT_NOCASH2 = 5
const SEDIT_NOBUY = 6
const SEDIT_BUY = 7
const SEDIT_SELL = 8
const SEDIT_PRODUCTS_MENU = 11
const SEDIT_ROOMS_MENU = 12
const SEDIT_NAMELIST_MENU = 13
const SEDIT_NAMELIST = 14
const SEDIT_NUMERICAL_RESPONSE = 20
const SEDIT_OPEN1 = 21
const SEDIT_OPEN2 = 22
const SEDIT_CLOSE1 = 23
const SEDIT_CLOSE2 = 24
const SEDIT_KEEPER = 25
const SEDIT_BUY_PROFIT = 26
const SEDIT_SELL_PROFIT = 27
const SEDIT_TYPE_MENU = 29
const SEDIT_DELETE_TYPE = 30
const SEDIT_DELETE_PRODUCT = 31
const SEDIT_NEW_PRODUCT = 32
const SEDIT_DELETE_ROOM = 33
const SEDIT_NEW_ROOM = 34
const SEDIT_SHOP_FLAGS = 35
const SEDIT_NOTRADE = 36
const SEDIT_COPY = 37
const CEDIT_MAIN_MENU = 0
const CEDIT_CONFIRM_SAVESTRING = 1
const CEDIT_GAME_OPTIONS_MENU = 2
const CEDIT_CRASHSAVE_OPTIONS_MENU = 3
const CEDIT_OPERATION_OPTIONS_MENU = 4
const CEDIT_DISP_EXPERIENCE_MENU = 5
const CEDIT_ROOM_NUMBERS_MENU = 6
const CEDIT_AUTOWIZ_OPTIONS_MENU = 7
const CEDIT_OK = 8
const CEDIT_NOPERSON = 9
const CEDIT_NOEFFECT = 10
const CEDIT_DFLT_IP = 11
const CEDIT_DFLT_DIR = 12
const CEDIT_LOGNAME = 13
const CEDIT_MENU = 14
const CEDIT_WELC_MESSG = 15
const CEDIT_START_MESSG = 16
const CEDIT_ADVANCE_OPTIONS_MENU = 17
const CEDIT_NUMERICAL_RESPONSE = 20
const CEDIT_LEVEL_CAN_SHOUT = 21
const CEDIT_HOLLER_MOVE_COST = 22
const CEDIT_TUNNEL_SIZE = 23
const CEDIT_MAX_EXP_GAIN = 24
const CEDIT_MAX_EXP_LOSS = 25
const CEDIT_MAX_NPC_CORPSE_TIME = 26
const CEDIT_MAX_PC_CORPSE_TIME = 27
const CEDIT_IDLE_VOID = 28
const CEDIT_IDLE_RENT_TIME = 29
const CEDIT_IDLE_MAX_LEVEL = 30
const CEDIT_DTS_ARE_DUMPS = 31
const CEDIT_LOAD_INTO_INVENTORY = 32
const CEDIT_TRACK_THROUGH_DOORS = 33
const CEDIT_LEVEL_CAP = 34
const CEDIT_MAX_OBJ_SAVE = 35
const CEDIT_MIN_RENT_COST = 36
const CEDIT_AUTOSAVE_TIME = 37
const CEDIT_CRASH_FILE_TIMEOUT = 38
const CEDIT_RENT_FILE_TIMEOUT = 39
const CEDIT_MORTAL_START_ROOM = 40
const CEDIT_IMMORT_START_ROOM = 41
const CEDIT_FROZEN_START_ROOM = 42
const CEDIT_DONATION_ROOM_1 = 43
const CEDIT_DONATION_ROOM_2 = 44
const CEDIT_DONATION_ROOM_3 = 45
const CEDIT_DFLT_PORT = 46
const CEDIT_MAX_PLAYING = 47
const CEDIT_MAX_FILESIZE = 48
const CEDIT_MAX_BAD_PWS = 49
const CEDIT_SITEOK_EVERYONE = 50
const CEDIT_NAMESERVER_IS_SLOW = 51
const CEDIT_USE_AUTOWIZ = 52
const CEDIT_MIN_WIZLIST_LEV = 53
const CEDIT_ALLOW_MULTICLASS = 54
const CEDIT_EXP_MULTIPLIER = 55
const CEDIT_PULSE_VIOLENCE = 56
const CEDIT_PULSE_MOBILE = 57
const CEDIT_PULSE_ZONE = 58
const CEDIT_PULSE_CURRENT = 59
const CEDIT_PULSE_IDLEPWD = 60
const CEDIT_PULSE_USAGE = 61
const CEDIT_PULSE_SANITY = 62
const CEDIT_PULSE_AUTOSAVE = 63
const CEDIT_PULSE_TIMESAVE = 64
const CEDIT_TICKS_OPTIONS_MENU = 65
const CEDIT_CREATION_OPTIONS_MENU = 66
const CEDIT_CREATION_MENU = 67
const CEDIT_POINTS_MENU = 68
const ASSEDIT_DO_NOT_USE = 0
const ASSEDIT_MAIN_MENU = 1
const ASSEDIT_ADD_COMPONENT = 2
const ASSEDIT_EDIT_COMPONENT = 3
const ASSEDIT_DELETE_COMPONENT = 4
const ASSEDIT_EDIT_EXTRACT = 5
const ASSEDIT_EDIT_INROOM = 6
const ASSEDIT_EDIT_TYPES = 7
const CEDIT_CREATION_METHOD_1 = 0
const CEDIT_CREATION_METHOD_2 = 1
const CEDIT_CREATION_METHOD_3 = 2
const CEDIT_CREATION_METHOD_4 = 3
const CEDIT_CREATION_METHOD_5 = 4
const HEDIT_CONFIRM_SAVESTRING = 0
const HEDIT_CONFIRM_EDIT = 1
const HEDIT_CONFIRM_ADD = 2
const HEDIT_MAIN_MENU = 3
const HEDIT_ENTRY = 4
const HEDIT_KEYWORDS = 5
const HEDIT_MIN_LEVEL = 6
const HSEDIT_MAIN_MENU = 0
const HSEDIT_CONFIRM_SAVESTRING = 1
const HSEDIT_OWNER_MENU = 2
const HSEDIT_OWNER_NAME = 3
const HSEDIT_OWNER_ID = 4
const HSEDIT_ROOM = 5
const HSEDIT_ATRIUM = 6
const HSEDIT_DIR_MENU = 7
const HSEDIT_GUEST_MENU = 8
const HSEDIT_GUEST_ADD = 9
const HSEDIT_GUEST_DELETE = 10
const HSEDIT_GUEST_CLEAR = 11
const HSEDIT_FLAGS = 12
const HSEDIT_BUILD_DATE = 13
const HSEDIT_PAYMENT = 14
const HSEDIT_TYPE = 15
const HSEDIT_DELETE = 16
const HSEDIT_VALUE_0 = 17
const HSEDIT_VALUE_1 = 18
const HSEDIT_VALUE_2 = 19
const HSEDIT_VALUE_3 = 20
const HSEDIT_NOVNUM = 21
const HSEDIT_BUILDER = 22
const CONTEXT_HELP_STRING = "help"
const CONTEXT_OEDIT_MAIN_MENU = 1
const CONTEXT_OEDIT_EDIT_NAMELIST = 2
const CONTEXT_OEDIT_SHORTDESC = 3
const CONTEXT_OEDIT_LONGDESC = 4
const CONTEXT_OEDIT_ACTDESC = 5
const CONTEXT_OEDIT_TYPE = 6
const CONTEXT_OEDIT_EXTRAS = 7
const CONTEXT_OEDIT_WEAR = 8
const CONTEXT_OEDIT_WEIGHT = 9
const CONTEXT_OEDIT_COST = 10
const CONTEXT_OEDIT_COSTPERDAY = 11
const CONTEXT_OEDIT_TIMER = 12
const CONTEXT_OEDIT_VALUE_1 = 13
const CONTEXT_OEDIT_VALUE_2 = 14
const CONTEXT_OEDIT_VALUE_3 = 15
const CONTEXT_OEDIT_VALUE_4 = 16
const CONTEXT_OEDIT_APPLY = 17
const CONTEXT_OEDIT_APPLYMOD = 18
const CONTEXT_OEDIT_EXTRADESC_KEY = 19
const CONTEXT_OEDIT_CONFIRM_SAVEDB = 20
const CONTEXT_OEDIT_CONFIRM_SAVESTRING = 21
const CONTEXT_OEDIT_PROMPT_APPLY = 22
const CONTEXT_OEDIT_EXTRADESC_DESCRIPTION = 23
const CONTEXT_OEDIT_EXTRADESC_MENU = 24
const CONTEXT_OEDIT_LEVEL = 25
const CONTEXT_OEDIT_PERM = 26
const CONTEXT_REDIT_MAIN_MENU = 27
const CONTEXT_REDIT_NAME = 28
const CONTEXT_REDIT_DESC = 29
const CONTEXT_REDIT_FLAGS = 30
const CONTEXT_REDIT_SECTOR = 31
const CONTEXT_REDIT_EXIT_MENU = 32
const CONTEXT_REDIT_CONFIRM_SAVEDB = 33
const CONTEXT_REDIT_CONFIRM_SAVESTRING = 34
const CONTEXT_REDIT_EXIT_NUMBER = 35
const CONTEXT_REDIT_EXIT_DESCRIPTION = 36
const CONTEXT_REDIT_EXIT_KEYWORD = 37
const CONTEXT_REDIT_EXIT_KEY = 38
const CONTEXT_REDIT_EXIT_DOORFLAGS = 39
const CONTEXT_REDIT_EXTRADESC_MENU = 40
const CONTEXT_REDIT_EXTRADESC_KEY = 41
const CONTEXT_REDIT_EXTRADESC_DESCRIPTION = 42
const CONTEXT_ZEDIT_MAIN_MENU = 43
const CONTEXT_ZEDIT_DELETE_ENTRY = 44
const CONTEXT_ZEDIT_NEW_ENTRY = 45
const CONTEXT_ZEDIT_CHANGE_ENTRY = 46
const CONTEXT_ZEDIT_COMMAND_TYPE = 47
const CONTEXT_ZEDIT_IF_FLAG = 48
const CONTEXT_ZEDIT_ARG1 = 49
const CONTEXT_ZEDIT_ARG2 = 50
const CONTEXT_ZEDIT_ARG3 = 51
const CONTEXT_ZEDIT_ZONE_NAME = 52
const CONTEXT_ZEDIT_ZONE_LIFE = 53
const CONTEXT_ZEDIT_ZONE_BOT = 54
const CONTEXT_ZEDIT_ZONE_TOP = 55
const CONTEXT_ZEDIT_ZONE_RESET = 56
const CONTEXT_ZEDIT_CONFIRM_SAVESTRING = 57
const CONTEXT_ZEDIT_SARG1 = 58
const CONTEXT_ZEDIT_SARG2 = 59
const CONTEXT_MEDIT_MAIN_MENU = 60
const CONTEXT_MEDIT_ALIAS = 61
const CONTEXT_MEDIT_S_DESC = 62
const CONTEXT_MEDIT_L_DESC = 63
const CONTEXT_MEDIT_D_DESC = 64
const CONTEXT_MEDIT_NPC_FLAGS = 65
const CONTEXT_MEDIT_AFF_FLAGS = 66
const CONTEXT_MEDIT_CONFIRM_SAVESTRING = 67
const CONTEXT_MEDIT_SEX = 68
const CONTEXT_MEDIT_ACCURACY = 69
const CONTEXT_MEDIT_DAMAGE = 70
const CONTEXT_MEDIT_NDD = 71
const CONTEXT_MEDIT_SDD = 72
const CONTEXT_MEDIT_NUM_HP_DICE = 73
const CONTEXT_MEDIT_SIZE_HP_DICE = 74
const CONTEXT_MEDIT_ADD_HP = 75
const CONTEXT_MEDIT_AC = 76
const CONTEXT_MEDIT_EXP = 77
const CONTEXT_MEDIT_GOLD = 78
const CONTEXT_MEDIT_POS = 79
const CONTEXT_MEDIT_DEFAULT_POS = 80
const CONTEXT_MEDIT_ATTACK = 81
const CONTEXT_MEDIT_LEVEL = 82
const CONTEXT_MEDIT_ALIGNMENT = 83
const CONTEXT_SEDIT_MAIN_MENU = 84
const CONTEXT_SEDIT_CONFIRM_SAVESTRING = 85
const CONTEXT_SEDIT_NOITEM1 = 86
const CONTEXT_SEDIT_NOITEM2 = 87
const CONTEXT_SEDIT_NOCASH1 = 88
const CONTEXT_SEDIT_NOCASH2 = 89
const CONTEXT_SEDIT_NOBUY = 90
const CONTEXT_SEDIT_BUY = 91
const CONTEXT_SEDIT_SELL = 92
const CONTEXT_SEDIT_PRODUCTS_MENU = 93
const CONTEXT_SEDIT_ROOMS_MENU = 94
const CONTEXT_SEDIT_NAMELIST_MENU = 95
const CONTEXT_SEDIT_NAMELIST = 96
const CONTEXT_SEDIT_OPEN1 = 97
const CONTEXT_SEDIT_OPEN2 = 98
const CONTEXT_SEDIT_CLOSE1 = 99
const CONTEXT_SEDIT_CLOSE2 = 100
const CONTEXT_SEDIT_KEEPER = 101
const CONTEXT_SEDIT_BUY_PROFIT = 102
const CONTEXT_SEDIT_SELL_PROFIT = 103
const CONTEXT_SEDIT_TYPE_MENU = 104
const CONTEXT_SEDIT_DELETE_TYPE = 105
const CONTEXT_SEDIT_DELETE_PRODUCT = 106
const CONTEXT_SEDIT_NEW_PRODUCT = 107
const CONTEXT_SEDIT_DELETE_ROOM = 108
const CONTEXT_SEDIT_NEW_ROOM = 109
const CONTEXT_SEDIT_SHOP_FLAGS = 110
const CONTEXT_SEDIT_NOTRADE = 111
const CONTEXT_TRIGEDIT_MAIN_MENU = 112
const CONTEXT_TRIGEDIT_TRIGTYPE = 113
const CONTEXT_TRIGEDIT_CONFIRM_SAVESTRING = 114
const CONTEXT_TRIGEDIT_NAME = 115
const CONTEXT_TRIGEDIT_INTENDED = 116
const CONTEXT_TRIGEDIT_TYPES = 117
const CONTEXT_TRIGEDIT_COMMANDS = 118
const CONTEXT_TRIGEDIT_NARG = 119
const CONTEXT_TRIGEDIT_ARGUMENT = 120
const CONTEXT_SCRIPT_MAIN_MENU = 121
const CONTEXT_SCRIPT_NEW_TRIGGER = 122
const CONTEXT_SCRIPT_DEL_TRIGGER = 123
const CONTEXT_ZEDIT_ARG4 = 124
const CONTEXT_GEDIT_MAIN_MENU = 125
const CONTEXT_GEDIT_CONFIRM_SAVESTRING = 126
const CONTEXT_GEDIT_NO_CASH uint8 = math.MaxInt8
const CONTEXT_GEDIT_NO_SKILL = 128
const CONTEXT_GEDIT_NUMERICAL_RESPONSE = 129
const CONTEXT_GEDIT_CHARGE = 130
const CONTEXT_GEDIT_OPEN = 131
const CONTEXT_GEDIT_CLOSE = 132
const CONTEXT_GEDIT_TRAINER = 133
const CONTEXT_GEDIT_NO_TRAIN = 134
const CONTEXT_GEDIT_MINLVL = 135
const CONTEXT_GEDIT_SELECT_SPELLS = 136
const CONTEXT_GEDIT_SELECT_SKILLS = 137
const CONTEXT_GEDIT_SELECT_WPS = 138
const CONTEXT_GEDIT_SELECT_LANGS = 139
const NUM_CONTEXTS = 140
const GEDIT_MAIN_MENU = 0
const GEDIT_CONFIRM_SAVESTRING = 1
const GEDIT_NO_CASH = 2
const GEDIT_NO_SKILL = 3
const GEDIT_NUMERICAL_RESPONSE = 5
const GEDIT_CHARGE = 6
const GEDIT_OPEN = 7
const GEDIT_CLOSE = 8
const GEDIT_TRAINER = 9
const GEDIT_NO_TRAIN = 10
const GEDIT_MINLVL = 11
const GEDIT_SELECT_SPELLS = 12
const GEDIT_SELECT_SKILLS = 13
const GEDIT_SELECT_WPS = 14
const GEDIT_SELECT_LANGS = 15
const GEDIT_SELECT_FEATS = 16
const STAT_GET_STR = 0
const STAT_GET_INT = 1
const STAT_GET_WIS = 2
const STAT_GET_DEX = 3
const STAT_GET_CON = 4
const STAT_GET_CHA = 5
const STAT_QUIT = 6
const STAT_PARSE_MENU = 7

type oasis_olc_data struct {
	Mode             int
	Zone_num         zone_rnum
	Number           room_vnum
	Value            int
	Storage          *byte
	Mob              *char_data
	Room             *room_data
	Obj              *obj_data
	Iobj             *obj_data
	Zone             *zone_data
	Shop             *shop_data
	House            *house_control_rec
	Config           *config_data
	Desc             *extra_descr_data
	Action           *social_messg
	Trig             *trig_data
	Script_mode      int
	Trigger_position int
	Item_type        int
	Script           *trig_proto_list
	OlcAssembly      *assembly_data
	Guild            *guild_data
	Help             *help_index_element
}
type olc_scmd_info_t struct {
	Text     *byte
	Con_type int
}

var olc_scmd_info [12]olc_scmd_info_t = [12]olc_scmd_info_t{{Text: libc.CString("room"), Con_type: CON_REDIT}, {Text: libc.CString("object"), Con_type: CON_OEDIT}, {Text: libc.CString("zone"), Con_type: CON_ZEDIT}, {Text: libc.CString("mobile"), Con_type: CON_MEDIT}, {Text: libc.CString("shop"), Con_type: CON_SEDIT}, {Text: libc.CString("config"), Con_type: CON_CEDIT}, {Text: libc.CString("trigger"), Con_type: CON_TRIGEDIT}, {Text: libc.CString("action"), Con_type: CON_AEDIT}, {Text: libc.CString("guild"), Con_type: CON_GEDIT}, {Text: libc.CString("help"), Con_type: CON_HEDIT}, {Text: libc.CString("house"), Con_type: CON_HSEDIT}, {Text: libc.CString("\n"), Con_type: -1}}

func clear_screen(d *descriptor_data) {
	if PRF_FLAGGED(d.Character, PRF_CLS) {
		write_to_output(d, libc.CString("\x1b[H\x1b[J"))
	}
}
func do_oasis(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) || ch.Desc == nil {
		return
	}
	if ch.Desc.Connected != CON_PLAYING {
		return
	}
	switch subcmd {
	case SCMD_OASIS_CEDIT:
		do_oasis_cedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_ZEDIT:
		do_oasis_zedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_REDIT:
		do_oasis_redit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_OEDIT:
		do_oasis_oedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_MEDIT:
		do_oasis_medit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_SEDIT:
		do_oasis_sedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_AEDIT:
		do_oasis_aedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_HEDIT:
		do_oasis_hedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_HSEDIT:
		do_oasis_hsedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_RLIST:
		fallthrough
	case SCMD_OASIS_MLIST:
		fallthrough
	case SCMD_OASIS_OLIST:
		fallthrough
	case SCMD_OASIS_SLIST:
		fallthrough
	case SCMD_OASIS_ZLIST:
		fallthrough
	case SCMD_OASIS_TLIST:
		fallthrough
	case SCMD_OASIS_GLIST:
		do_oasis_list(ch, argument, cmd, subcmd)
	case SCMD_OASIS_TRIGEDIT:
		do_oasis_trigedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_LINKS:
		do_oasis_links(ch, argument, cmd, subcmd)
	case SCMD_OASIS_GEDIT:
		do_oasis_gedit(ch, argument, cmd, subcmd)
	default:
		basic_mud_log(libc.CString("SYSERR: (OLC) Invalid subcmd passed to do_oasis, subcmd - (%d)"), subcmd)
		return
	}
	return
}
func cleanup_olc(d *descriptor_data, cleanup_type int8) {
	if d.Olc == nil {
		return
	}
	if d.Olc.Room != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_proto_script(unsafe.Pointer(d.Olc.Room), WLD_TRIGGER)
			free_room(d.Olc.Room)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Room))
		case CLEANUP_CONFIG:
			free_config(d.Olc.Config)
		default:
			basic_mud_log(libc.CString("SYSERR: cleanup_olc: Unknown type!"))
		}
	}
	if d.Olc.Obj != nil {
		free_object_strings(d.Olc.Obj)
		libc.Free(unsafe.Pointer(d.Olc.Obj))
	}
	if d.Olc.Mob != nil {
		free_mobile(d.Olc.Mob)
	}
	if d.Olc.Zone != nil {
		if d.Olc.Zone.Builders != nil {
			libc.Free(unsafe.Pointer(d.Olc.Zone.Builders))
		}
		if d.Olc.Zone.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Zone.Name))
		}
		if d.Olc.Zone.Cmd != nil {
			libc.Free(unsafe.Pointer(d.Olc.Zone.Cmd))
		}
		libc.Free(unsafe.Pointer(d.Olc.Zone))
	}
	if d.Olc.Shop != nil {
		free_shop(d.Olc.Shop)
	}
	if d.Olc.Guild != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_guild(d.Olc.Guild)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Guild))
		default:
		}
	}
	if d.Olc.House != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_house(d.Olc.House)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.House))
		default:
		}
	}
	if d.Olc.Action != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_action(d.Olc.Action)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Action))
		default:
		}
	}
	if d.Olc.Storage != nil {
		libc.Free(unsafe.Pointer(d.Olc.Storage))
		d.Olc.Storage = nil
	}
	if d.Olc.Trig != nil {
		free_trigger(d.Olc.Trig)
		d.Olc.Trig = nil
	}
	if d.Character != nil {
		d.Character.Act[int(PLR_WRITING/32)] &= bitvector_t(int32(^(1 << (int(PLR_WRITING % 32)))))
		act(libc.CString("$n stops using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
		if int(cleanup_type) == CLEANUP_CONFIG {
			mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s stops editing the game configuration"), GET_NAME(d.Character))
		} else if d.Connected == CON_TEDIT {
			mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s stops editing text files."), GET_NAME(d.Character))
		} else if d.Connected == CON_HEDIT {
			mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s stops editing help files."), GET_NAME(d.Character))
		} else {
			mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s stops editing zone %d allowed zone %d"), GET_NAME(d.Character), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, d.Character.Player_specials.Olc_zone)
		}
		d.Connected = CON_PLAYING
	}
	libc.Free(unsafe.Pointer(d.Olc))
	d.Olc = nil
}
func split_argument(argument *byte, tag *byte) {
	var (
		tmp  *byte = argument
		ttag *byte = tag
		wrt  *byte = argument
		i    int
	)
	for i = 0; *tmp != 0; func() int {
		tmp = (*byte)(unsafe.Add(unsafe.Pointer(tmp), 1))
		return func() int {
			p := &i
			x := *p
			*p++
			return x
		}()
	}() {
		if *tmp != ' ' && *tmp != '=' {
			*(func() *byte {
				p := &ttag
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *tmp
		} else if *tmp == '=' {
			break
		}
	}
	*ttag = '\x00'
	for *tmp == '=' || *tmp == ' ' {
		tmp = (*byte)(unsafe.Add(unsafe.Pointer(tmp), 1))
	}
	for *tmp != 0 {
		*(func() *byte {
			p := &wrt
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(func() *byte {
			p := &tmp
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	*wrt = '\x00'
}
func free_config(data *config_data) {
	free_strings(unsafe.Pointer(data), OASIS_CFG)
	libc.Free(unsafe.Pointer(data))
}
func can_edit_zone(ch *char_data, rnum zone_rnum) int {
	if ch.Desc == nil || IS_NPC(ch) || rnum == zone_rnum(-1) {
		return FALSE
	}
	if rnum == HEDIT_PERMISSION {
		return TRUE
	}
	if ch.Admlevel >= ADMLVL_GRGOD {
		return TRUE
	}
	if is_name(GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Builders) != 0 {
		return TRUE
	}
	if ch.Player_specials.Olc_zone == int(-1) {
		return FALSE
	}
	if ch.Admlevel < ADMLVL_BUILDER {
		return FALSE
	}
	if real_zone(zone_vnum(ch.Player_specials.Olc_zone)) == rnum {
		return TRUE
	}
	return FALSE
}
func send_cannot_edit(ch *char_data, zone zone_vnum) {
	send_to_char(ch, libc.CString("You do not have permission to edit zone %d."), zone)
	if ch.Player_specials.Olc_zone != int(-1) {
		send_to_char(ch, libc.CString("  Try zone %d."), ch.Player_specials.Olc_zone)
	}
	send_to_char(ch, libc.CString("\r\n"))
	mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("OLC: %s tried to edit zone %d allowed zone %d"), GET_NAME(ch), zone, ch.Player_specials.Olc_zone)
}
