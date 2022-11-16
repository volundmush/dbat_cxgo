package main

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"os"
	"unicode"
	"unsafe"
)

const DB_BOOT_WLD = 0
const DB_BOOT_MOB = 1
const DB_BOOT_OBJ = 2
const DB_BOOT_ZON = 3
const DB_BOOT_SHP = 4
const DB_BOOT_HLP = 5
const DB_BOOT_TRG = 6
const DB_BOOT_GLD = 7

const LIB_USER = "user/"
const LIB_INTRO = "intro/"
const LIB_SENSE = "sense/"
const LIB_WORLD = "world/"
const LIB_TEXT = "text/"
const LIB_TEXT_HELP = "text/help/"
const LIB_MISC = "misc/"
const LIB_ETC = "etc/"
const LIB_PLRTEXT = "plrtext/"
const LIB_PLROBJS = "plrobjs/"
const LIB_PLRVARS = "plrvars/"
const LIB_PLRALIAS = "plralias/"
const LIB_PLRFILES = "plrfiles/"
const LIB_HOUSE = "house/"
const LIB_PLRIMC = "plrimc/"
const SLASH = "/"
const SUF_OBJS = "new"
const SUF_TEXT = "text"
const SUF_ALIAS = "alias"
const SUF_MEM = "mem"
const SUF_PLR = "plr"
const SUF_PET = "pet"
const SUF_IMC = "imc"
const SUF_USER = "usr"
const SUF_INTRO = "itr"
const SUF_SENSE = "sen"
const SUF_CUSTOM = "cus"

const FASTBOOT_FILE = "../.fastboot"
const KILLSCRIPT_FILE = "../.killscript"
const PAUSE_FILE = "../pause"

const INDEX_FILE = "index"
const MINDEX_FILE = "index.mini"
const HELP_FILE = "help.hlp"
const PINDEX_DELETED = 1
const PINDEX_NODELETE = 2
const PINDEX_SELFDELETE = 4
const PINDEX_NOWIZLIST = 8
const REAL = 0
const VIRTUAL = 1
const CUR_WORLD_VERSION = 1
const CUR_ZONE_VERSION = 2
const BAN_NOT = 0
const BAN_NEW = 1
const BAN_SELECT = 2
const BAN_ALL = 3
const BANNED_SITE_LENGTH = 50
const DISABLED_FILE = "disabled.cmds"
const END_MARKER = "END"

const WLD_PREFIX = LIB_WORLD + "wld" + SLASH /* room definitions	*/
const MOB_PREFIX = LIB_WORLD + "mob" + SLASH /* monster prototypes	*/
const OBJ_PREFIX = LIB_WORLD + "obj" + SLASH /* object prototypes	*/
const ZON_PREFIX = LIB_WORLD + "zon" + SLASH /* zon defs & command tables */
const SHP_PREFIX = LIB_WORLD + "shp" + SLASH /* shop definitions	*/
const HLP_PREFIX = LIB_TEXT + "help" + SLASH /* for HELP <keyword>	*/
const TRG_PREFIX = LIB_WORLD + "trg" + SLASH /* trigger files	*/
const GLD_PREFIX = LIB_WORLD + "gld" + SLASH /* guild files		*/

const CREDITS_FILE = LIB_TEXT + "credits"          /* for the 'credits' command	*/
const NEWS_FILE = LIB_TEXT + "news"                /* for the 'news' command	*/
const MOTD_FILE = LIB_TEXT + "motd"                /* messages of the day / mortal	*/
const IMOTD_FILE = LIB_TEXT + "imotd"              /* messages of the day / immort	*/
const GREETINGS_FILE = LIB_TEXT + "greetings"      /* The opening screen.	*/
const GREETANSI_FILE = LIB_TEXT + "greetansi"      /* The opening screen.	*/
const HELP_PAGE_FILE = LIB_TEXT_HELP + "screen"    /* for HELP <CR>	*/
const CONTEXT_HELP_FILE = LIB_TEXT + "contexthelp" /* The opening screen.	*/
const INFO_FILE = LIB_TEXT + "info"                /* for INFO		*/
const WIZLIST_FILE = LIB_TEXT + "wizlist"          /* for WIZLIST		*/
const IMMLIST_FILE = LIB_TEXT + "immlist"          /* for IMMLIST		*/
const BACKGROUND_FILE = LIB_TEXT + "background"    /* for the background story	*/
const POLICIES_FILE = LIB_TEXT + "policies"        /* player policies/rules	*/
const HANDBOOK_FILE = LIB_TEXT + "handbook"        /* handbook for new immorts	*/
const IHELP_PAGE_FILE = LIB_TEXT_HELP + "iscreen"  /* for HELP <CR>	*/

const IDEA_FILE = LIB_MISC + "ideas"              /* for the 'idea'-command	*/
const TYPO_FILE = LIB_MISC + "typos"              /*         'typo'		*/
const BUG_FILE = LIB_MISC + "bugs"                /*         'bug'		*/
const REQUEST_FILE = LIB_MISC + "request"         /*      RPP Requests         */
const CUSTOM_FILE = LIB_MISC + "customs"          /*      Custom EQ            */
const MESS_FILE = LIB_MISC + "messages"           /* damage messages		*/
const SOCMESS_FILE = LIB_MISC + "socials"         /* messages for social acts	*/
const SOCMESS_FILE_NEW = LIB_MISC + "socials.new" /* messages for social acts with aedit patch*/
const XNAME_FILE = LIB_MISC + "xnames"            /* invalid name substrings	*/

const CONFIG_FILE = LIB_ETC + "config"         /* OasisOLC * GAME CONFIG FL */
const PLAYER_FILE = LIB_ETC + "players"        /* the player database	*/
const MAIL_FILE = LIB_ETC + "plrmail"          /* for the mudmail system	*/
const BAN_FILE = LIB_ETC + "badsites"          /* for the siteban system	*/
const HCONTROL_FILE = LIB_ETC + "hcontrol"     /* for the house system	*/
const TIME_FILE = LIB_ETC + "time"             /* for calendar system	*/
const AUCTION_FILE = LIB_ETC + "auction"       /* for the auction house system */
const ASSEMBLIES_FILE = LIB_ETC + "assemblies" /* for assemblies system 	*/
const LEVEL_CONFIG = LIB_ETC + "levels"        /* set various level values  */

type help_index_element struct {
	Index     *byte
	Keywords  *byte
	Entry     *byte
	Duplicate int
	Min_level int
}
type reset_com struct {
	Command int8
	If_flag bool
	Arg1    vnum
	Arg2    vnum
	Arg3    vnum
	Arg4    vnum
	Arg5    vnum
	Line    int
	Sarg1   *byte
	Sarg2   *byte
}
type zone_data struct {
	Name       *byte
	Builders   *byte
	Lifespan   int
	Age        int
	Bot        room_vnum
	Top        room_vnum
	Reset_mode int
	Number     zone_vnum
	Cmd        []reset_com
	Min_level  int
	Max_level  int
	Zone_flags [4]bitvector_t
}
type reset_q_element struct {
	Zone_to_reset zone_rnum
	Next          *reset_q_element
}
type reset_q_type struct {
	Head *reset_q_element
	Tail *reset_q_element
}
type player_index_element struct {
	Name     *byte
	Id       int
	Level    int
	Admlevel int
	Flags    int
	Last     libc.Time
	Ship     int
	Shiproom int
	Played   libc.Time
	Clan     *byte
}
type ban_list_element struct {
	Site [51]byte
	Type int
	Date libc.Time
	Name [21]byte
	Next *ban_list_element
}

type disabled_data struct {
	Next        *disabled_data
	Command     *command_info
	Disabled_by *byte
	Level       int16
	Subcmd      int
}

const NUM_OBJ_UNIQUE_POOLS = 5000
const ZO_DEAD = 999

var config_info config_data
var world []room_data = nil
var top_of_world room_rnum = 0
var room_htree *htree_node = nil
var character_list *char_data = nil
var affect_list *char_data = nil
var affectv_list *char_data = nil
var mob_index []index_data
var mob_proto []char_data
var top_of_mobt mob_rnum = 0
var mob_htree *htree_node = nil
var object_list *obj_data = nil
var obj_index []index_data
var obj_proto []obj_data
var top_of_objt obj_rnum = 0
var obj_htree *htree_node = nil
var zone_table []zone_data
var top_of_zone_table zone_rnum = 0
var fight_messages [100]message_list
var trig_index []*index_data
var trigger_list *trig_data = nil
var top_of_trigt int = 0
var max_mob_id int = MOB_ID_BASE
var max_obj_id int = OBJ_ID_BASE
var dg_owner_purged int
var no_mail int = 0
var mini_mud int = 0
var no_rent_check int = 0
var boot_time libc.Time = 0
var circle_restrict int = 0
var dballtime int = 0
var SHENRON int = FALSE
var DRAGONR int = 0
var DRAGONZ int = 0
var WISH [2]int = [2]int{}
var DRAGONC int = 0
var EDRAGON *char_data = nil
var r_mortal_start_room room_rnum
var r_immort_start_room room_rnum
var r_frozen_start_room room_rnum
var xap_objs int = 0
var converting int = FALSE
var credits *byte = nil
var news *byte = nil
var motd *byte = nil
var imotd *byte = nil
var GREETINGS *byte = nil
var GREETANSI *byte = nil
var help *byte = nil
var info *byte = nil
var wizlist *byte = nil
var immlist *byte = nil
var background *byte = nil
var handbook *byte = nil
var policies *byte = nil
var ihelp *byte = nil
var help_table []help_index_element = nil
var top_of_helpt int = 0
var soc_mess_list []social_messg = nil
var top_of_socialt int = -1
var time_info time_info_data
var weather_info weather_data
var dummy_mob player_special_data
var reset_q reset_q_type

func dragon_level(ch *char_data) {
	var (
		d     *descriptor_data
		level int = 0
		count int = 0
	)
	for d = descriptor_list; d != nil; d = d.Next {
		if IS_PLAYING(d) && d.Character.Admlevel < 1 {
			level += GET_LEVEL(d.Character)
			count += 1
		}
	}
	if level > 0 && count > 0 {
		level = level / count
	} else {
		level = rand_number(60, 110)
	}
	if level < 50 {
		level = rand_number(40, 60)
	}
	ch.Race_level = 0
	ch.Race_level = level + rand_number(5, 20)
}
func mob_stats(mob *char_data) {
	var (
		start  int = int(float64(GET_LEVEL(mob)) * 0.5)
		finish int = GET_LEVEL(mob)
	)
	if finish < 20 {
		finish = 20
	}
	if !IS_HUMANOID(mob) {
		mob.Real_abils.Str = int8(rand_number(start, finish))
		mob.Real_abils.Intel = int8(rand_number(start, finish) - 30)
		mob.Real_abils.Wis = int8(rand_number(start, finish) - 30)
		mob.Real_abils.Dex = int8(rand_number(start+5, finish))
		mob.Real_abils.Con = int8(rand_number(start+5, finish))
		mob.Real_abils.Cha = int8(rand_number(start, finish))
	} else {
		if int(mob.Race) == RACE_SAIYAN {
			mob.Real_abils.Str = int8(rand_number(start+10, finish))
			mob.Real_abils.Intel = int8(rand_number(start, finish-10))
			mob.Real_abils.Wis = int8(rand_number(start, finish-5))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start+5, finish))
			mob.Real_abils.Cha = int8(rand_number(start+5, finish))
		} else if int(mob.Race) == RACE_KONATSU {
			mob.Real_abils.Str = int8(rand_number(start, finish-10))
			mob.Real_abils.Intel = int8(rand_number(start, finish))
			mob.Real_abils.Wis = int8(rand_number(start, finish))
			mob.Real_abils.Dex = int8(rand_number(start+10, finish))
			mob.Real_abils.Con = int8(rand_number(start, finish))
			mob.Real_abils.Cha = int8(rand_number(start, finish))
		} else if int(mob.Race) == RACE_ANDROID {
			mob.Real_abils.Str = int8(rand_number(start, finish))
			mob.Real_abils.Intel = int8(rand_number(start, finish))
			mob.Real_abils.Wis = int8(rand_number(start, finish-10))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start, finish))
			mob.Real_abils.Cha = int8(rand_number(start, finish))
		} else if int(mob.Race) == RACE_MAJIN {
			mob.Real_abils.Str = int8(rand_number(start, finish))
			mob.Real_abils.Intel = int8(rand_number(start, finish-10))
			mob.Real_abils.Wis = int8(rand_number(start, finish-5))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start+15, finish))
			mob.Real_abils.Cha = int8(rand_number(start, finish))
		} else if int(mob.Race) == RACE_TRUFFLE {
			mob.Real_abils.Str = int8(rand_number(start, finish-10))
			mob.Real_abils.Intel = int8(rand_number(start+15, finish))
			mob.Real_abils.Wis = int8(rand_number(start, finish))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start, finish))
			mob.Real_abils.Cha = int8(rand_number(start, finish))
		} else if int(mob.Race) == RACE_ICER {
			mob.Real_abils.Str = int8(rand_number(start+5, finish))
			mob.Real_abils.Intel = int8(rand_number(start, finish))
			mob.Real_abils.Wis = int8(rand_number(start, finish))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start, finish))
			mob.Real_abils.Cha = int8(rand_number(start+10, finish))
		} else {
			mob.Real_abils.Str = int8(rand_number(start, finish))
			mob.Real_abils.Intel = int8(rand_number(start, finish))
			mob.Real_abils.Wis = int8(rand_number(start, finish))
			mob.Real_abils.Dex = int8(rand_number(start, finish))
			mob.Real_abils.Con = int8(rand_number(start, finish))
			mob.Real_abils.Cha = int8(rand_number(start, finish))
		}
	}
	if int(mob.Real_abils.Str) > 100 {
		mob.Real_abils.Str = 100
	} else if int(mob.Real_abils.Str) < 5 {
		mob.Real_abils.Str = int8(rand_number(5, 8))
	}
	if int(mob.Real_abils.Intel) > 100 {
		mob.Real_abils.Intel = 100
	} else if int(mob.Real_abils.Intel) < 5 {
		mob.Real_abils.Intel = int8(rand_number(5, 8))
	}
	if int(mob.Real_abils.Wis) > 100 {
		mob.Real_abils.Wis = 100
	} else if int(mob.Real_abils.Wis) < 5 {
		mob.Real_abils.Wis = int8(rand_number(5, 8))
	}
	if int(mob.Real_abils.Con) > 100 {
		mob.Real_abils.Con = 100
	} else if int(mob.Real_abils.Con) < 5 {
		mob.Real_abils.Con = int8(rand_number(5, 8))
	}
	if int(mob.Real_abils.Cha) > 100 {
		mob.Real_abils.Cha = 100
	} else if int(mob.Real_abils.Cha) < 5 {
		mob.Real_abils.Cha = int8(rand_number(5, 8))
	}
	if int(mob.Real_abils.Dex) > 100 {
		mob.Real_abils.Dex = 100
	} else if int(mob.Real_abils.Dex) < 5 {
		mob.Real_abils.Dex = int8(rand_number(5, 8))
	}
}
func suntzu_armor_convert(obj *obj_data) int {
	var (
		i            int
		conv         int       = 0
		conv_table   [9][3]int = [9][3]int{{100, 0, 0}, {8, 0, 5}, {6, 0, 10}, {5, 1, 15}, {4, 2, 20}, {2, 5, 30}, {0, 7, 40}, {0, 7, 40}, {1, 6, 35}}
		shield_table [9][2]int = [9][2]int{{}, {1, 5}, {2, 15}, {3, 30}, {4, 40}, {5, 50}, {6, 60}, {7, 70}, {8, 80}}
	)
	i = obj.Value[0]
	if i != 0 && i < 10 {
		obj.Value[0] = i * 10
		conv = 1
	} else {
		i /= 10
	}
	i = int(MAX(0, MIN(8, int64(i))))
	if OBJWEAR_FLAGGED(obj, ITEM_WEAR_SHIELD) {
		if (obj.Value[6]) != 0 {
			return conv
		}
		obj.Value[1] = ARMOR_TYPE_SHIELD
		obj.Value[2] = 100
		obj.Value[3] = shield_table[i][0]
		obj.Value[6] = shield_table[i][1]
		conv = 1
	} else if OBJWEAR_FLAGGED(obj, ITEM_WEAR_BODY) {
		if (obj.Value[6]) != 0 {
			return conv
		}
		obj.Value[2] = conv_table[i][0]
		obj.Value[3] = conv_table[i][1]
		obj.Value[6] = conv_table[i][2]
		conv = 1
	} else if (obj.Value[2]) != 0 || (obj.Value[3]) != 0 {
		return conv
	} else {
		obj.Value[2] = 100
		obj.Value[3] = 0
		obj.Value[6] = 0
		conv = 1
	}
	//basic_mud_log(libc.CString("Converted armor #%d [%s] armor=%d i=%d maxdex=%d acheck=%d sfail=%d"), obj_index[obj-&obj_proto[0]*(667*667)].Vnum, obj.Short_description, obj.Value[0], i, obj.Value[2], obj.Value[3], obj.Value[6])
	return conv
}
func suntzu_weapon_convert(wp_type int) int {
	var new_type int
	switch wp_type {
	case 170:
		new_type = WEAPON_TYPE_DAGGER
	case 171:
		new_type = WEAPON_TYPE_SHORTSWORD
	case 172:
		new_type = WEAPON_TYPE_LONGSWORD
	case 173:
		new_type = WEAPON_TYPE_GREATSWORD
	case 174:
		new_type = WEAPON_TYPE_MACE
	case 175:
		new_type = WEAPON_TYPE_AXE
	case 176:
		new_type = WEAPON_TYPE_WHIP
	case 177:
		new_type = WEAPON_TYPE_SPEAR
	case 178:
		new_type = WEAPON_TYPE_POLEARM
	case 179:
		new_type = WEAPON_TYPE_UNARMED
	case 180:
		new_type = WEAPON_TYPE_FLAIL
	case 181:
		new_type = WEAPON_TYPE_STAFF
	case 182:
		new_type = WEAPON_TYPE_HAMMER
	default:
		new_type = WEAPON_TYPE_UNDEFINED
	}
	basic_mud_log(libc.CString("Converted weapon from [%d] to [%d]."), wp_type, new_type)
	return new_type
}
func reboot_wizlists() {
	file_to_string_alloc(libc.CString(LIB_TEXT), &wizlist)
	file_to_string_alloc(libc.CString(LIB_TEXT), &immlist)
}
func free_text_files() {
	var (
		textfiles [15]**byte = [15]**byte{&wizlist, &immlist, &news, &credits, &motd, &imotd, &help, &info, &policies, &handbook, &background, &GREETINGS, &GREETANSI, &ihelp, nil}
		rf        int
	)
	for rf = 0; textfiles[rf] != nil; rf++ {
		if *textfiles[rf] != nil {
			libc.Free(unsafe.Pointer(*textfiles[rf]))
			*textfiles[rf] = nil
		}
	}
}
func do_reboot(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if libc.StrCaseCmp(&arg[0], libc.CString("all")) == 0 || arg[0] == '*' {
		if load_levels() < 0 {
			send_to_char(ch, libc.CString("Cannot read level configurations\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETINGS) == 0 {
			prune_crlf(GREETINGS)
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETANSI) == 0 {
			prune_crlf(GREETANSI)
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &wizlist) < 0 {
			send_to_char(ch, libc.CString("Cannot read wizlist\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &immlist) < 0 {
			send_to_char(ch, libc.CString("Cannot read immlist\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &news) < 0 {
			send_to_char(ch, libc.CString("Cannot read news\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &credits) < 0 {
			send_to_char(ch, libc.CString("Cannot read credits\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &motd) < 0 {
			send_to_char(ch, libc.CString("Cannot read motd\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &imotd) < 0 {
			send_to_char(ch, libc.CString("Cannot read imotd\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT_HELP), &help) < 0 {
			send_to_char(ch, libc.CString("Cannot read help front page\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &info) < 0 {
			send_to_char(ch, libc.CString("Cannot read info file\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &policies) < 0 {
			send_to_char(ch, libc.CString("Cannot read policies\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &handbook) < 0 {
			send_to_char(ch, libc.CString("Cannot read handbook\r\n"))
		}
		if file_to_string_alloc(libc.CString(LIB_TEXT), &background) < 0 {
			send_to_char(ch, libc.CString("Cannot read background\r\n"))
		}
		if help_table != nil {
			free_help_table()
		}
		index_boot(DB_BOOT_HLP)
	} else if libc.StrCaseCmp(&arg[0], libc.CString("levels")) == 0 {
		if load_levels() < 0 {
			send_to_char(ch, libc.CString("Cannot read level configurations\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("wizlist")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &wizlist) < 0 {
			send_to_char(ch, libc.CString("Cannot read wizlist\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("immlist")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &immlist) < 0 {
			send_to_char(ch, libc.CString("Cannot read immlist\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("news")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &news) < 0 {
			send_to_char(ch, libc.CString("Cannot read news\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("credits")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &credits) < 0 {
			send_to_char(ch, libc.CString("Cannot read credits\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("motd")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &motd) < 0 {
			send_to_char(ch, libc.CString("Cannot read motd\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("imotd")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &imotd) < 0 {
			send_to_char(ch, libc.CString("Cannot read imotd\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("help")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT_HELP), &help) < 0 {
			send_to_char(ch, libc.CString("Cannot read help front page\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("info")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &info) < 0 {
			send_to_char(ch, libc.CString("Cannot read info\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("policy")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &policies) < 0 {
			send_to_char(ch, libc.CString("Cannot read policy\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("handbook")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &handbook) < 0 {
			send_to_char(ch, libc.CString("Cannot read handbook\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("background")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &background) < 0 {
			send_to_char(ch, libc.CString("Cannot read background\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("greetings")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETINGS) == 0 {
			prune_crlf(GREETINGS)
		} else {
			send_to_char(ch, libc.CString("Cannot read greetings.\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("greetansi")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETANSI) == 0 {
			prune_crlf(GREETANSI)
		} else {
			send_to_char(ch, libc.CString("Cannot read greetings.\r\n"))
		}
	} else if libc.StrCaseCmp(&arg[0], libc.CString("xhelp")) == 0 {
		if help_table != nil {
			free_help_table()
		}
		index_boot(DB_BOOT_HLP)
	} else if libc.StrCaseCmp(&arg[0], libc.CString("ihelp")) == 0 {
		if file_to_string_alloc(libc.CString(LIB_TEXT_HELP), &ihelp) < 0 {
			send_to_char(ch, libc.CString("Cannot read help front page\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("Unknown reload option.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
}
func boot_world() {
	basic_mud_log(libc.CString("Loading level tables."))
	load_levels()
	basic_mud_log(libc.CString("Loading zone table."))
	index_boot(DB_BOOT_ZON)
	basic_mud_log(libc.CString("Loading triggers and generating index."))
	index_boot(DB_BOOT_TRG)
	basic_mud_log(libc.CString("Loading rooms."))
	index_boot(DB_BOOT_WLD)
	basic_mud_log(libc.CString("Renumbering rooms."))
	renum_world()
	basic_mud_log(libc.CString("Checking start rooms."))
	check_start_rooms()
	basic_mud_log(libc.CString("Loading mobs and generating index."))
	index_boot(DB_BOOT_MOB)
	basic_mud_log(libc.CString("Loading objs and generating index."))
	index_boot(DB_BOOT_OBJ)
	basic_mud_log(libc.CString("Renumbering zone table."))
	renum_zone_table()
	basic_mud_log(libc.CString("Loading disabled commands list..."))
	load_disabled()
	if converting != 0 {
		basic_mud_log(libc.CString("Saving converted worldfiles to disk."))
		save_all()
	}
	if no_specials == 0 {
		basic_mud_log(libc.CString("Loading shops."))
		index_boot(DB_BOOT_SHP)
		basic_mud_log(libc.CString("Loading guild masters."))
		index_boot(DB_BOOT_GLD)
	}
	if SELFISHMETER >= 10 {
		basic_mud_log(libc.CString("Loading Shadow Dragons."))
		load_shadow_dragons()
	}
}
func free_extra_descriptions(edesc *extra_descr_data) {
	var enext *extra_descr_data
	for ; edesc != nil; edesc = enext {
		enext = edesc.Next
		libc.Free(unsafe.Pointer(edesc.Keyword))
		libc.Free(unsafe.Pointer(edesc.Description))
		libc.Free(unsafe.Pointer(edesc))
	}
}
func destroy_db() {
	var (
		cnt    int64
		itr    int64
		chtmp  *char_data
		objtmp *obj_data
	)
	for character_list != nil {
		chtmp = character_list
		character_list = character_list.Next
		if chtmp.Master != nil {
			stop_follower(chtmp)
		}
		free_char(chtmp)
	}
	for object_list != nil {
		objtmp = object_list
		object_list = object_list.Next
		free_obj(objtmp)
	}
	for cnt = 0; cnt <= int64(top_of_world); cnt++ {
		if world[cnt].Name != nil {
			libc.Free(unsafe.Pointer(world[cnt].Name))
		}
		if world[cnt].Description != nil {
			libc.Free(unsafe.Pointer(world[cnt].Description))
		}
		free_extra_descriptions(world[cnt].Ex_description)
		if world[cnt].Script != nil {
			extract_script(unsafe.Pointer(&world[cnt]), WLD_TRIGGER)
		}
		free_proto_script(unsafe.Pointer(&world[cnt]), WLD_TRIGGER)
		for itr = 0; itr < NUM_OF_DIRS; itr++ {
			if world[cnt].Dir_option[itr] == nil {
				continue
			}
			if world[cnt].Dir_option[itr].General_description != nil {
				libc.Free(unsafe.Pointer(world[cnt].Dir_option[itr].General_description))
			}
			if world[cnt].Dir_option[itr].Keyword != nil {
				libc.Free(unsafe.Pointer(world[cnt].Dir_option[itr].Keyword))
			}
			libc.Free(unsafe.Pointer(world[cnt].Dir_option[itr]))
		}
	}
	libc.Free(unsafe.Pointer(&world[0]))
	top_of_world = 0
	htree_free(room_htree)
	for cnt = 0; cnt <= int64(top_of_objt); cnt++ {
		if obj_proto[cnt].Name != nil {
			libc.Free(unsafe.Pointer(obj_proto[cnt].Name))
		}
		if obj_proto[cnt].Description != nil {
			libc.Free(unsafe.Pointer(obj_proto[cnt].Description))
		}
		if obj_proto[cnt].Short_description != nil {
			libc.Free(unsafe.Pointer(obj_proto[cnt].Short_description))
		}
		if obj_proto[cnt].Action_description != nil {
			libc.Free(unsafe.Pointer(obj_proto[cnt].Action_description))
		}
		if obj_proto[cnt].Ex_description != nil {
			free_extra_descriptions(obj_proto[cnt].Ex_description)
		}
		free_proto_script(unsafe.Pointer(&obj_proto[cnt]), OBJ_TRIGGER)
	}
	libc.Free(unsafe.Pointer(&obj_proto[0]))
	libc.Free(unsafe.Pointer(&obj_index[0]))
	htree_free(obj_htree)
	for cnt = 0; cnt <= int64(top_of_mobt); cnt++ {
		if mob_proto[cnt].Name != nil {
			libc.Free(unsafe.Pointer(mob_proto[cnt].Name))
		}
		if mob_proto[cnt].Title != nil {
			libc.Free(unsafe.Pointer(mob_proto[cnt].Title))
		}
		if mob_proto[cnt].Short_descr != nil {
			libc.Free(unsafe.Pointer(mob_proto[cnt].Short_descr))
		}
		if mob_proto[cnt].Long_descr != nil {
			libc.Free(unsafe.Pointer(mob_proto[cnt].Long_descr))
		}
		if mob_proto[cnt].Description != nil {
			libc.Free(unsafe.Pointer(mob_proto[cnt].Description))
		}
		free_proto_script(unsafe.Pointer(&mob_proto[cnt]), MOB_TRIGGER)
		for mob_proto[cnt].Affected != nil {
			affect_remove(&mob_proto[cnt], mob_proto[cnt].Affected)
		}
	}
	libc.Free(unsafe.Pointer(&mob_proto[0]))
	libc.Free(unsafe.Pointer(&mob_index[0]))
	htree_free(mob_htree)
	destroy_shops()
	destroy_guilds()
	if reset_q.Head != nil {
		var (
			ftemp *reset_q_element = reset_q.Head
			temp  *reset_q_element
		)
		for ftemp != nil {
			temp = ftemp.Next
			libc.Free(unsafe.Pointer(ftemp))
			ftemp = temp
		}
	}
	for cnt = 0; cnt <= int64(top_of_zone_table); cnt++ {
		if zone_table[cnt].Name != nil {
			libc.Free(unsafe.Pointer(zone_table[cnt].Name))
		}
		if zone_table[cnt].Builders != nil {
			libc.Free(unsafe.Pointer(zone_table[cnt].Builders))
		}
		if zone_table[cnt].Cmd != nil {
			for itr = 0; int(zone_table[cnt].Cmd[itr].Command) != 'S'; itr++ {
				if int(zone_table[cnt].Cmd[itr].Command) == 'V' {
					if zone_table[cnt].Cmd[itr].Sarg1 != nil {
						libc.Free(unsafe.Pointer(zone_table[cnt].Cmd[itr].Sarg1))
					}
					if zone_table[cnt].Cmd[itr].Sarg2 != nil {
						libc.Free(unsafe.Pointer(zone_table[cnt].Cmd[itr].Sarg2))
					}
				}
			}
			libc.Free(unsafe.Pointer(&zone_table[cnt].Cmd[0]))
		}
	}
	libc.Free(unsafe.Pointer(&zone_table[0]))
	if reset_q.Head != nil {
		var (
			ftemp *reset_q_element = reset_q.Head
			temp  *reset_q_element
		)
		for ftemp != nil {
			temp = ftemp.Next
			libc.Free(unsafe.Pointer(ftemp))
			ftemp = temp
		}
	}
	for cnt = 0; cnt < int64(top_of_trigt); cnt++ {
		if trig_index[cnt].Proto != nil {
			if trig_index[cnt].Proto.Cmdlist != nil {
				var (
					i *cmdlist_element
					j *cmdlist_element
				)
				i = trig_index[cnt].Proto.Cmdlist
				for i != nil {
					j = i.Next
					if i.Cmd != nil {
						libc.Free(unsafe.Pointer(i.Cmd))
					}
					libc.Free(unsafe.Pointer(i))
					i = j
				}
			}
			free_trigger(trig_index[cnt].Proto)
		}
		libc.Free(unsafe.Pointer(trig_index[cnt]))
	}
	libc.Free(unsafe.Pointer(&trig_index[0]))
	event_free_all()
	free_context_help()
	free_feats()
	free_obj_unique_hash()
	htree_shutdown()
	basic_mud_log(libc.CString("Freeing Assemblies."))
	free_assemblies()
}

var obj_unique_hash_pools **obj_unique_hash_elem = nil

func init_obj_unique_hash() {
	var i int
	obj_unique_hash_pools = &make([]*obj_unique_hash_elem, NUM_OBJ_UNIQUE_POOLS)[0]
	for i = 0; i < NUM_OBJ_UNIQUE_POOLS; i++ {
		*(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i))) = nil
	}
}
func boot_db() {
	var i zone_rnum
	basic_mud_log(libc.CString("Boot db -- BEGIN."))
	basic_mud_log(libc.CString("Resetting the game time:"))
	reset_time()
	basic_mud_log(libc.CString("Reading news, credits, help, ihelp, bground, info & motds."))
	file_to_string_alloc(libc.CString(LIB_TEXT), &news)
	file_to_string_alloc(libc.CString(LIB_TEXT), &credits)
	file_to_string_alloc(libc.CString(LIB_TEXT), &motd)
	file_to_string_alloc(libc.CString(LIB_TEXT), &imotd)
	file_to_string_alloc(libc.CString(LIB_TEXT), &info)
	file_to_string_alloc(libc.CString(LIB_TEXT), &wizlist)
	file_to_string_alloc(libc.CString(LIB_TEXT), &immlist)
	file_to_string_alloc(libc.CString(LIB_TEXT), &policies)
	file_to_string_alloc(libc.CString(LIB_TEXT), &handbook)
	file_to_string_alloc(libc.CString(LIB_TEXT), &background)
	file_to_string_alloc(libc.CString(LIB_TEXT_HELP), &help)
	file_to_string_alloc(libc.CString(LIB_TEXT_HELP), &ihelp)
	if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETINGS) == 0 {
		prune_crlf(GREETINGS)
	}
	if file_to_string_alloc(libc.CString(LIB_TEXT), &GREETANSI) == 0 {
		prune_crlf(GREETANSI)
	}
	basic_mud_log(libc.CString("Loading spell definitions."))
	mag_assign_spells()
	basic_mud_log(libc.CString("Loading feats."))
	assign_feats()
	boot_world()
	htree_test()
	basic_mud_log(libc.CString("Loading help entries."))
	index_boot(DB_BOOT_HLP)
	basic_mud_log(libc.CString("Setting up context sensitive help system for OLC"))
	boot_context_help()
	basic_mud_log(libc.CString("Generating player index."))
	build_player_index()
	if ERAPLAYERS <= 0 {
		ERAPLAYERS = top_of_p_table + 1
	}
	insure_directory(libc.CString(LIB_PLROBJS), 0)
	basic_mud_log(libc.CString("Booting mail system."))
	if scan_file() == 0 {
		basic_mud_log(libc.CString("    Mail boot failed -- Mail system disabled"))
		no_mail = 1
	}
	if auto_pwipe != 0 {
		basic_mud_log(libc.CString("Cleaning out inactive players."))
		clean_pfiles()
	}
	basic_mud_log(libc.CString("Loading social messages."))
	boot_social_messages()
	basic_mud_log(libc.CString("Loading Clans."))
	clanBoot()
	basic_mud_log(libc.CString("Building command list."))
	create_command_list()
	basic_mud_log(libc.CString("Assigning function pointers:"))
	if no_specials == 0 {
		basic_mud_log(libc.CString("   Mobiles."))
		assign_mobiles()
		basic_mud_log(libc.CString("   Shopkeepers."))
		assign_the_shopkeepers()
		basic_mud_log(libc.CString("   Objects."))
		assign_objects()
		basic_mud_log(libc.CString("   Rooms."))
		assign_rooms()
		basic_mud_log(libc.CString("   Guildmasters."))
		assign_the_guilds()
	}
	basic_mud_log(libc.CString("Init Object Unique Hash"))
	init_obj_unique_hash()
	basic_mud_log(libc.CString("Booting assembled objects."))
	assemblyBootAssemblies()
	basic_mud_log(libc.CString("Sorting command list and spells."))
	sort_commands()
	sort_spells()
	sort_feats()
	basic_mud_log(libc.CString("Booting boards system."))
	init_boards()
	basic_mud_log(libc.CString("Reading banned site and invalid-name list."))
	load_banned()
	Read_Invalid_List()
	if no_rent_check == 0 {
		basic_mud_log(libc.CString("Deleting timed-out crash and rent files:"))
		update_obj_file()
		basic_mud_log(libc.CString("   Done."))
	}
	if mini_mud == 0 {
		basic_mud_log(libc.CString("Booting houses."))
		House_boot()
	}
	for i = 0; i <= top_of_zone_table; i++ {
		basic_mud_log(libc.CString("Resetting #%d: %s (rooms %d-%d)."), zone_table[i].Number, zone_table[i].Name, zone_table[i].Bot, zone_table[i].Top)
		reset_zone(i)
	}
	reset_q.Head = func() *reset_q_element {
		p := &reset_q.Tail
		reset_q.Tail = nil
		return *p
	}()
	boot_time = libc.GetTime(nil)
	basic_mud_log(libc.CString("Boot db -- DONE."))
}
func auc_save() {
	var fl *stdio.File
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Can't write to '%s' auction file."), LIB_ETC)
	} else {
		var obj *obj_data
		for obj = world[real_room(80)].Contents; obj != nil; obj = obj.Next_content {
			if obj != nil {
				stdio.Fprintf(fl, "%lld %s %d %d %d %d %ld\n", obj.Unique_id, obj.Auctname, obj.Aucter, obj.CurBidder, obj.Startbid, obj.Bid, obj.AucTime)
			}
		}
		stdio.Fprintf(fl, "~END~\n")
		fl.Close()
	}
}
func auc_load(obj *obj_data) {
	var (
		line   [500]byte
		filler [50]byte
		oID    int64
		timer  libc.Time
		aID    int
		bID    int
		cost   int
		startc int
		fl     *stdio.File
	)
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Can't read from '%s' auction file."), LIB_ETC)
	} else {
		for int(fl.IsEOF()) == 0 {
			get_line(fl, &line[0])
			stdio.Sscanf(&line[0], "%lld %s %d %d %d %d %ld\n", &oID, &filler[0], &aID, &bID, &startc, &cost, &timer)
			if obj.Unique_id == oID {
				obj.Auctname = libc.StrDup(&filler[0])
				obj.Aucter = int32(aID)
				obj.CurBidder = int32(bID)
				obj.Startbid = startc
				obj.Bid = cost
				obj.AucTime = timer
			}
		}
		fl.Close()
	}
}
func reset_time() {
	var (
		beginning_of_time libc.Time = 0
		bgtime            *stdio.File
	)
	if (func() *stdio.File {
		bgtime = stdio.FOpen(TIME_FILE, "r")
		return bgtime
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Can't read from '%s' time file."), TIME_FILE)
	} else {
		stdio.Fscanf(bgtime, "%d\n", &beginning_of_time)
		stdio.Fscanf(bgtime, "%d\n", &NEWSUPDATE)
		stdio.Fscanf(bgtime, "%d\n", &BOARDNEWMORT)
		stdio.Fscanf(bgtime, "%d\n", &BOARDNEWDUO)
		stdio.Fscanf(bgtime, "%d\n", &BOARDNEWCOD)
		stdio.Fscanf(bgtime, "%d\n", &BOARDNEWBUI)
		stdio.Fscanf(bgtime, "%d\n", &BOARDNEWIMM)
		stdio.Fscanf(bgtime, "%d\n", &INTERESTTIME)
		stdio.Fscanf(bgtime, "%d\n", &LASTINTEREST)
		stdio.Fscanf(bgtime, "%d\n", &HIGHPCOUNT)
		stdio.Fscanf(bgtime, "%d\n", &PCOUNTDATE)
		stdio.Fscanf(bgtime, "%d\n", &WISHTIME)
		stdio.Fscanf(bgtime, "%d\n", &PCOUNT)
		stdio.Fscanf(bgtime, "%d\n", &LASTPAYOUT)
		stdio.Fscanf(bgtime, "%d\n", &LASTPAYTYPE)
		stdio.Fscanf(bgtime, "%d\n", &LASTNEWS)
		stdio.Fscanf(bgtime, "%d\n", &dballtime)
		stdio.Fscanf(bgtime, "%d\n", &SELFISHMETER)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON1)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON2)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON3)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON4)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON5)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON6)
		stdio.Fscanf(bgtime, "%d\n", &SHADOW_DRAGON7)
		stdio.Fscanf(bgtime, "%d\n", &ERAPLAYERS)
		bgtime.Close()
	}
	if dballtime == 0 {
		dballtime = 604800
	}
	if beginning_of_time == 0 {
		beginning_of_time = 0x26C359CB
	}
	time_info = *mud_time_passed(libc.GetTime(nil), beginning_of_time)
	if time_info.Hours <= 4 {
		weather_info.Sunlight = SUN_DARK
	} else if time_info.Hours == 5 {
		weather_info.Sunlight = SUN_RISE
	} else if time_info.Hours <= 20 {
		weather_info.Sunlight = SUN_LIGHT
	} else if time_info.Hours == 21 {
		weather_info.Sunlight = SUN_SET
	} else {
		weather_info.Sunlight = SUN_DARK
	}
	basic_mud_log(libc.CString("   Current Gametime: %dH %dD %dM %dY."), time_info.Hours, time_info.Day, time_info.Month, time_info.Year)
	weather_info.Pressure = 960
	if time_info.Month >= 7 && time_info.Month <= 12 {
		weather_info.Pressure += dice(1, 50)
	} else {
		weather_info.Pressure += dice(1, 80)
	}
	weather_info.Change = 0
	if weather_info.Pressure <= 980 {
		weather_info.Sky = SKY_LIGHTNING
	} else if weather_info.Pressure <= 1000 {
		weather_info.Sky = SKY_RAINING
	} else if weather_info.Pressure <= 1020 {
		weather_info.Sky = SKY_CLOUDY
	} else {
		weather_info.Sky = SKY_CLOUDLESS
	}
}
func save_mud_time(when *time_info_data) {
	var bgtime *stdio.File
	if (func() *stdio.File {
		bgtime = stdio.FOpen(LIB_ETC, "w")
		return bgtime
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Can't write to '%s' time file."), LIB_ETC)
	} else {
		stdio.Fprintf(bgtime, "%ld\n", mud_time_to_secs(when))
		stdio.Fprintf(bgtime, "%ld\n", NEWSUPDATE)
		stdio.Fprintf(bgtime, "%ld\n", BOARDNEWMORT)
		stdio.Fprintf(bgtime, "%ld\n", BOARDNEWDUO)
		stdio.Fprintf(bgtime, "%ld\n", BOARDNEWCOD)
		stdio.Fprintf(bgtime, "%ld\n", BOARDNEWBUI)
		stdio.Fprintf(bgtime, "%ld\n", BOARDNEWIMM)
		stdio.Fprintf(bgtime, "%ld\n", INTERESTTIME)
		stdio.Fprintf(bgtime, "%ld\n", LASTINTEREST)
		stdio.Fprintf(bgtime, "%d\n", HIGHPCOUNT)
		stdio.Fprintf(bgtime, "%ld\n", PCOUNTDATE)
		stdio.Fprintf(bgtime, "%d\n", WISHTIME)
		stdio.Fprintf(bgtime, "%d\n", PCOUNT)
		stdio.Fprintf(bgtime, "%ld\n", LASTPAYOUT)
		stdio.Fprintf(bgtime, "%d\n", LASTPAYTYPE)
		stdio.Fprintf(bgtime, "%d\n", LASTNEWS)
		stdio.Fprintf(bgtime, "%d\n", dballtime)
		stdio.Fprintf(bgtime, "%d\n", SELFISHMETER)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON1)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON2)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON3)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON4)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON5)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON6)
		stdio.Fprintf(bgtime, "%d\n", SHADOW_DRAGON7)
		stdio.Fprintf(bgtime, "%d\n", ERAPLAYERS)
		bgtime.Close()
	}
}
func count_alias_records(fl *stdio.File) int {
	var (
		key            [256]byte
		next_key       [256]byte
		line           [256]byte
		scan           *byte
		total_keywords int = 0
	)
	get_one_line(fl, &key[0])
	for key[0] != '$' {
		for {
			get_one_line(fl, &line[0])
			if int(fl.IsEOF()) != 0 {
				goto ackeof
			}
			if line[0] == '#' {
				break
			}
		}
		scan = &key[0]
		for {
			scan = one_word(scan, &next_key[0])
			if next_key[0] != 0 {
				total_keywords++
			}
			if next_key[0] == 0 {
				break
			}
		}
		get_one_line(fl, &key[0])
		if int(fl.IsEOF()) != 0 {
			goto ackeof
		}
	}
	return total_keywords
ackeof:
	basic_mud_log(libc.CString("SYSERR: Unexpected end of help file."))
	os.Exit(1)
	return -1
}
func count_hash_records(fl *stdio.File) int {
	var (
		buf   [128]byte
		count int = 0
	)
	for fl.GetS(&buf[0], 128) != nil {
		if buf[0] == '#' {
			count++
		}
	}
	return count
}
func index_boot(mode int) {
	var (
		index_filename *byte
		prefix         *byte = nil
		db_index       *stdio.File
		db_file        *stdio.File
		rec_count      int = 0
		size           [2]int
		buf2           [260]byte
		buf1           [64936]byte
	)
	switch mode {
	case DB_BOOT_WLD:
		prefix = libc.CString(WLD_PREFIX)
	case DB_BOOT_MOB:
		prefix = libc.CString(MOB_PREFIX)
	case DB_BOOT_OBJ:
		prefix = libc.CString(OBJ_PREFIX)
	case DB_BOOT_ZON:
		prefix = libc.CString(ZON_PREFIX)
	case DB_BOOT_SHP:
		prefix = libc.CString(SHP_PREFIX)
	case DB_BOOT_HLP:
		prefix = libc.CString(HLP_PREFIX)
	case DB_BOOT_TRG:
		prefix = libc.CString(TRG_PREFIX)
	case DB_BOOT_GLD:
		prefix = libc.CString(GLD_PREFIX)
	default:
		basic_mud_log(libc.CString("SYSERR: Unknown subcommand %d to index_boot!"), mode)
		os.Exit(1)
	}
	if mini_mud != 0 {
		index_filename = libc.CString(MINDEX_FILE)
	} else {
		index_filename = libc.CString(INDEX_FILE)
	}
	stdio.Snprintf(&buf2[0], int(260), "%s%s", prefix, index_filename)
	if (func() *stdio.File {
		db_index = stdio.FOpen(libc.GoString(&buf2[0]), "r")
		return db_index
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: opening index file '%s': %s"), &buf2[0], libc.StrError(libc.Errno))
		os.Exit(1)
	}
	stdio.Fscanf(db_index, "%s\r\n", &buf1[0])
	for buf1[0] != '$' {
		stdio.Snprintf(&buf2[0], int(260), "%s%s", prefix, &buf1[0])
		if (func() *stdio.File {
			db_file = stdio.FOpen(libc.GoString(&buf2[0]), "r")
			return db_file
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: File '%s' listed in '%s%s': %s"), &buf2[0], prefix, index_filename, libc.StrError(libc.Errno))
			stdio.Fscanf(db_index, "%s\n", &buf1[0])
			continue
		} else {
			if mode == DB_BOOT_ZON {
				rec_count++
			} else if mode == DB_BOOT_HLP {
				rec_count += count_alias_records(db_file)
			} else {
				rec_count += count_hash_records(db_file)
			}
		}
		db_file.Close()
		stdio.Fscanf(db_index, "%s\n", &buf1[0])
	}
	if rec_count == 0 {
		if mode == DB_BOOT_SHP || mode == DB_BOOT_GLD {
			return
		}
		basic_mud_log(libc.CString("SYSERR: boot error - 0 records counted in %s/%s."), prefix, index_filename)
		os.Exit(1)
	}
	switch mode {
	case DB_BOOT_TRG:
		trig_index = make([]*index_data, rec_count)
	case DB_BOOT_WLD:
		world = make([]room_data, rec_count)
		size[0] = rec_count * int(unsafe.Sizeof(room_data{}))
		basic_mud_log(libc.CString("   %d rooms, %d bytes."), rec_count, size[0])
	case DB_BOOT_MOB:
		mob_proto = make([]char_data, rec_count)
		mob_index = make([]index_data, rec_count)
		size[0] = rec_count * int(unsafe.Sizeof(index_data{}))
		size[1] = rec_count * int(unsafe.Sizeof(char_data{}))
		basic_mud_log(libc.CString("   %d mobs, %d bytes in index, %d bytes in prototypes."), rec_count, size[0], size[1])
	case DB_BOOT_OBJ:
		obj_proto = make([]obj_data, rec_count)
		obj_index = make([]index_data, rec_count)
		size[0] = rec_count * int(unsafe.Sizeof(index_data{}))
		size[1] = rec_count * int(unsafe.Sizeof(obj_data{}))
		basic_mud_log(libc.CString("   %d objs, %d bytes in index, %d bytes in prototypes."), rec_count, size[0], size[1])
	case DB_BOOT_ZON:
		zone_table = make([]zone_data, rec_count)
		size[0] = rec_count * int(unsafe.Sizeof(zone_data{}))
		basic_mud_log(libc.CString("   %d zones, %d bytes."), rec_count, size[0])
	case DB_BOOT_HLP:
		help_table = make([]help_index_element, rec_count)
		size[0] = rec_count * int(unsafe.Sizeof(help_index_element{}))
		basic_mud_log(libc.CString("   %d entries, %d bytes."), rec_count, size[0])
	}
	db_index.Seek(0, 0)
	stdio.Fscanf(db_index, "%s\n", &buf1[0])
	for buf1[0] != '$' {
		stdio.Snprintf(&buf2[0], int(260), "%s%s", prefix, &buf1[0])
		if (func() *stdio.File {
			db_file = stdio.FOpen(libc.GoString(&buf2[0]), "r")
			return db_file
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: %s: %s"), &buf2[0], libc.StrError(libc.Errno))
			os.Exit(1)
		}
		switch mode {
		case DB_BOOT_WLD:
			fallthrough
		case DB_BOOT_OBJ:
			fallthrough
		case DB_BOOT_MOB:
			fallthrough
		case DB_BOOT_TRG:
			discrete_load(db_file, mode, &buf2[0])
		case DB_BOOT_ZON:
			load_zones(db_file, &buf2[0])
		case DB_BOOT_HLP:
			load_help(db_file, &buf2[0])
		case DB_BOOT_SHP:
			boot_the_shops(db_file, &buf2[0], rec_count)
		case DB_BOOT_GLD:
			boot_the_guilds(db_file, &buf2[0], rec_count)
		}
		db_file.Close()
		stdio.Fscanf(db_index, "%s\n", &buf1[0])
	}
	db_index.Close()
	if mode == DB_BOOT_HLP {
		libc.Sort(unsafe.Pointer(&help_table[0]), uint32(int32(top_of_helpt)), uint32(unsafe.Sizeof(help_index_element{})), func(arg1 unsafe.Pointer, arg2 unsafe.Pointer) int32 {
			return int32(hsort(arg1, arg2))
		})
		top_of_helpt--
	}
}
func discrete_load(fl *stdio.File, mode int, filename *byte) {
	var (
		nr    int = -1
		last  int
		line  [256]byte
		modes [7]*byte = [7]*byte{libc.CString("world"), libc.CString("mob"), libc.CString("obj"), libc.CString("ZON"), libc.CString("SHP"), libc.CString("HLP"), libc.CString("trg")}
	)
	for {
		if mode != DB_BOOT_OBJ || nr < 0 {
			if get_line(fl, &line[0]) == 0 {
				if nr == -1 {
					basic_mud_log(libc.CString("SYSERR: %s file %s is empty!"), modes[mode], filename)
				} else {
					basic_mud_log(libc.CString("SYSERR: Format error in %s after %s #%d\n...expecting a new %s, but file ended!\n(maybe the file is not terminated with '$'?)"), filename, modes[mode], nr, modes[mode])
				}
				os.Exit(1)
			}
		}
		if line[0] == '$' {
			return
		}
		if line[0] == '#' {
			last = nr
			if stdio.Sscanf(&line[0], "#%d", &nr) != 1 {
				basic_mud_log(libc.CString("SYSERR: Format error after %s #%d"), modes[mode], last)
				os.Exit(1)
			}
			if nr >= 99999 {
				return
			} else {
				switch mode {
				case DB_BOOT_WLD:
					parse_room(fl, nr)
				case DB_BOOT_MOB:
					parse_mobile(fl, nr)
				case DB_BOOT_TRG:
					parse_trigger(fl, nr)
				case DB_BOOT_OBJ:
					strlcpy(&line[0], parse_object(fl, nr), uint64(256))
				}
			}
		} else {
			basic_mud_log(libc.CString("SYSERR: Format error in %s file %s near %s #%d"), modes[mode], filename, modes[mode], nr)
			basic_mud_log(libc.CString("SYSERR: ... offending line: '%s'"), &line[0])
			os.Exit(1)
		}
	}
}
func fread_letter(fp *stdio.File) int8 {
	var c int8
	for {
		c = int8(fp.GetC())
		if !unicode.IsSpace(rune(c)) {
			break
		}
	}
	return c
}
func asciiflag_conv(flag *byte) bitvector_t {
	var (
		flags  bitvector_t = 0
		is_num int         = TRUE
		p      *byte
	)
	for p = flag; *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if unicode.IsLower(rune(*p)) {
			flags |= bitvector_t(1 << (*p - 'a'))
		} else if unicode.IsUpper(rune(*p)) {
			flags |= bitvector_t(1 << ((*p - 'A') + 26))
		}
		if !unicode.IsDigit(rune(*p)) && *p != '-' {
			is_num = FALSE
		}
	}
	if is_num != 0 {
		flags = bitvector_t(int32(libc.Atoi(libc.GoString(flag))))
	}
	return flags
}
func asciiflag_conv_aff(flag *byte) bitvector_t {
	var (
		flags  bitvector_t = 0
		is_num int         = TRUE
		p      *byte
	)
	for p = flag; *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if unicode.IsLower(rune(*p)) {
			flags |= bitvector_t(1 << ((*p - 'a') + 1))
		} else if unicode.IsUpper(rune(*p)) {
			flags |= bitvector_t(1 << ((*p - 'A') + 26))
		}
		if !unicode.IsDigit(rune(*p)) && *p != '-' {
			is_num = FALSE
		}
	}
	if is_num != 0 {
		flags = bitvector_t(int32(libc.Atoi(libc.GoString(flag))))
	}
	return flags
}
func parse_room(fl *stdio.File, virtual_nr int) {
	var (
		room_nr   int = 0
		zone      int = 0
		t         [10]int
		i         int
		retval    int
		line      [256]byte
		flags     [128]byte
		flags2    [128]byte
		flags3    [128]byte
		flags4    [128]byte
		buf2      [64936]byte
		buf       [128]byte
		new_descr *extra_descr_data
		letter    int8
	)
	stdio.Snprintf(&buf2[0], int(64936), "room #%d", virtual_nr)
	if virtual_nr < int(zone_table[zone].Bot) {
		basic_mud_log(libc.CString("SYSERR: Room #%d is below zone %d."), virtual_nr, zone)
		os.Exit(1)
	}
	for virtual_nr > int(zone_table[zone].Top) {
		if func() int {
			p := &zone
			*p++
			return *p
		}() > int(top_of_zone_table) {
			basic_mud_log(libc.CString("SYSERR: Room %d is outside of any zone."), virtual_nr)
			os.Exit(1)
		}
	}
	world[room_nr].Zone = zone_rnum(zone)
	world[room_nr].Number = room_vnum(virtual_nr)
	world[room_nr].Name = fread_string(fl, &buf2[0])
	world[room_nr].Description = fread_string(fl, &buf2[0])
	if room_htree == nil {
		room_htree = htree_init()
	}
	htree_add(room_htree, int64(virtual_nr), int64(room_nr))
	if get_line(fl, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Expecting roomflags/sector type of room #%d but file ended!"), virtual_nr)
		os.Exit(1)
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], " %d %s %s %s %s %d ", &t[0], &flags[0], &flags2[0], &flags3[0], &flags4[0], &t[2])
		return retval
	}()) == 3 && bitwarning == TRUE {
		basic_mud_log(libc.CString("WARNING: Conventional worldfiles detected. Please read 128bit.readme."))
		os.Exit(1)
	} else if retval == 3 && bitwarning == FALSE {
		basic_mud_log(libc.CString("Converting room #%d to 128bits.."), virtual_nr)
		world[room_nr].Room_flags[0] = asciiflag_conv(&flags[0])
		world[room_nr].Room_flags[1] = 0
		world[room_nr].Room_flags[2] = 0
		world[room_nr].Room_flags[3] = 0
		stdio.Sprintf(&flags[0], "room #%d", virtual_nr)
		check_bitvector_names(world[room_nr].Room_flags[0], room_bits_count, &flags[0], libc.CString("room"))
		if bitsavetodisk != 0 {
			add_to_save_list(zone_table[real_zone_by_thing(room_vnum(virtual_nr))].Number, 3)
			converting = TRUE
		}
		basic_mud_log(libc.CString("   done."))
	} else if retval == 6 {
		var taeller int
		world[room_nr].Room_flags[0] = asciiflag_conv(&flags[0])
		world[room_nr].Room_flags[1] = asciiflag_conv(&flags2[0])
		world[room_nr].Room_flags[2] = asciiflag_conv(&flags3[0])
		world[room_nr].Room_flags[3] = asciiflag_conv(&flags4[0])
		world[room_nr].Sector_type = t[2]
		stdio.Sprintf(&flags[0], "object #%d", virtual_nr)
		for taeller = 0; taeller < AF_ARRAY_MAX; taeller++ {
			check_bitvector_names(world[room_nr].Room_flags[taeller], room_bits_count, &flags[0], libc.CString("room"))
		}
	} else {
		basic_mud_log(libc.CString("SYSERR: Format error in roomflags/sector type of room #%d"), virtual_nr)
		os.Exit(1)
	}
	world[room_nr].Func = nil
	world[room_nr].Contents = nil
	world[room_nr].People = nil
	world[room_nr].Light = 0
	world[room_nr].Timed = -1
	world[room_nr].Dmg = 0
	world[room_nr].Gravity = 0
	if ROOM_FLAGGED(room_rnum(room_nr), ROOM_VEGETA) || ROOM_FLAGGED(room_rnum(room_nr), ROOM_GRAVITYX10) {
		world[room_nr].Gravity = 10
	}
	if world[room_nr].Number >= 19800 && world[room_nr].Number <= 0x4DBB {
		world[room_nr].Gravity = 1000
	}
	if world[room_nr].Number >= 64000 && world[room_nr].Number <= 0xFA06 {
		world[room_nr].Gravity = 100
	}
	if world[room_nr].Number >= 0xFA07 && world[room_nr].Number <= 0xFA10 {
		world[room_nr].Gravity = 300
	}
	if world[room_nr].Number >= 0xFA11 && world[room_nr].Number <= 64030 {
		world[room_nr].Gravity = 500
	}
	if world[room_nr].Number >= 0xFA1F && world[room_nr].Number <= 0xFA30 {
		world[room_nr].Gravity = 1000
	}
	if world[room_nr].Number >= 0xFA31 && world[room_nr].Number <= 64070 {
		world[room_nr].Gravity = 5000
	}
	if world[room_nr].Number >= 0xFA47 && world[room_nr].Number <= 0xFA60 {
		world[room_nr].Gravity = 10000
	}
	if world[room_nr].Number == 0xFA61 {
		world[room_nr].Gravity = 1000
	}
	for i = 0; i < NUM_OF_DIRS; i++ {
		world[room_nr].Dir_option[i] = nil
	}
	world[room_nr].Ex_description = nil
	stdio.Snprintf(&buf[0], int(128), "SYSERR: Format error in room #%d (expecting D/E/S)", virtual_nr)
	for {
		if get_line(fl, &line[0]) == 0 {
			basic_mud_log(libc.CString("%s"), &buf[0])
			os.Exit(1)
		}
		switch line[0] {
		case 'D':
			setup_dir(fl, room_nr, libc.Atoi(libc.GoString(&line[1])))
		case 'E':
			new_descr = new(extra_descr_data)
			new_descr.Keyword = fread_string(fl, &buf2[0])
			new_descr.Description = fread_string(fl, &buf2[0])
			{
				var tmp *byte = libc.StrChr(new_descr.Description, '\x00')
				if uintptr(unsafe.Pointer(tmp)) > uintptr(unsafe.Pointer(new_descr.Description)) && *((*byte)(unsafe.Add(unsafe.Pointer(tmp), -1))) != '\n' {
					tmp = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(new_descr.Description)+3)[0]))
					stdio.Sprintf(tmp, "%s\r\n", new_descr.Description)
					libc.Free(unsafe.Pointer(new_descr.Description))
					new_descr.Description = tmp
				}
			}
			new_descr.Next = world[room_nr].Ex_description
			world[room_nr].Ex_description = new_descr
		case 'S':
			letter = fread_letter(fl)
			fl.UnGetC(int(letter))
			for int(letter) == 'T' {
				dg_read_trigger(fl, unsafe.Pointer(&world[room_nr]), WLD_TRIGGER)
				letter = fread_letter(fl)
				fl.UnGetC(int(letter))
			}
			top_of_world = room_rnum(func() int {
				p := &room_nr
				x := *p
				*p++
				return x
			}())
			return
		default:
			basic_mud_log(libc.CString("%s"), &buf[0])
			os.Exit(1)
		}
	}
}
func setup_dir(fl *stdio.File, room int, dir int) {
	var (
		t      [11]int
		retval int
		line   [256]byte
		buf2   [128]byte
	)
	stdio.Snprintf(&buf2[0], int(128), "room #%d, direction D%d", int(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(room))))+1, dir)
	world[room].Dir_option[dir] = new(room_direction_data)
	world[room].Dir_option[dir].General_description = fread_string(fl, &buf2[0])
	world[room].Dir_option[dir].Keyword = fread_string(fl, &buf2[0])
	if get_line(fl, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Format error, %s"), &buf2[0])
		os.Exit(1)
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], " %d %d %d %d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7], &t[8], &t[9], &t[10])
		return retval
	}()) == 3 && bitwarning == TRUE {
		basic_mud_log(libc.CString("SYSERR: Format error, %s"), &buf2[0])
		os.Exit(1)
	} else if bitwarning == FALSE {
		if t[0] == 1 {
			world[room].Dir_option[dir].Exit_info = 1 << 0
		} else if t[0] == 2 {
			world[room].Dir_option[dir].Exit_info = (1 << 0) | 1<<3
		} else if t[0] == 3 {
			world[room].Dir_option[dir].Exit_info = (1 << 0) | 1<<4
		} else if t[0] == 4 {
			world[room].Dir_option[dir].Exit_info = (1 << 0) | 1<<3 | 1<<4
		} else {
			world[room].Dir_option[dir].Exit_info = 0
		}
		if t[1] == -1 || t[1] == math.MaxUint16 {
			world[room].Dir_option[dir].Key = -1
		} else {
			world[room].Dir_option[dir].Key = obj_vnum(t[1])
		}
		if t[2] == -1 || t[2] == math.MaxUint16 {
			world[room].Dir_option[dir].To_room = -1
		} else {
			world[room].Dir_option[dir].To_room = room_rnum(t[2])
		}
		if retval == 3 {
			basic_mud_log(libc.CString("Converting world files to include DC add ons."))
			world[room].Dir_option[dir].Dclock = 20
			world[room].Dir_option[dir].Dchide = 20
			world[room].Dir_option[dir].Dcskill = 0
			world[room].Dir_option[dir].Dcmove = 0
			world[room].Dir_option[dir].Failsavetype = 0
			world[room].Dir_option[dir].Dcfailsave = 0
			world[room].Dir_option[dir].Failroom = -1
			world[room].Dir_option[dir].Totalfailroom = -1
			if bitsavetodisk != 0 {
				add_to_save_list(zone_table[world[room].Zone].Number, 3)
				converting = TRUE
			}
		} else if retval == 5 {
			world[room].Dir_option[dir].Dclock = t[3]
			world[room].Dir_option[dir].Dchide = t[4]
			world[room].Dir_option[dir].Dcskill = 0
			world[room].Dir_option[dir].Dcmove = 0
			world[room].Dir_option[dir].Failsavetype = 0
			world[room].Dir_option[dir].Dcfailsave = 0
			world[room].Dir_option[dir].Failroom = -1
			world[room].Dir_option[dir].Totalfailroom = -1
			if bitsavetodisk != 0 {
				add_to_save_list(zone_table[world[room].Zone].Number, 3)
				converting = TRUE
			}
		} else if retval == 7 {
			world[room].Dir_option[dir].Dclock = t[3]
			world[room].Dir_option[dir].Dchide = t[4]
			world[room].Dir_option[dir].Dcskill = t[5]
			world[room].Dir_option[dir].Dcmove = t[6]
			world[room].Dir_option[dir].Failsavetype = 0
			world[room].Dir_option[dir].Dcfailsave = 0
			world[room].Dir_option[dir].Failroom = -1
			world[room].Dir_option[dir].Totalfailroom = -1
			if bitsavetodisk != 0 {
				add_to_save_list(zone_table[world[room].Zone].Number, 3)
				converting = TRUE
			}
		} else if retval == 11 {
			world[room].Dir_option[dir].Dclock = t[3]
			world[room].Dir_option[dir].Dchide = t[4]
			world[room].Dir_option[dir].Dcskill = t[5]
			world[room].Dir_option[dir].Dcmove = t[6]
			world[room].Dir_option[dir].Failsavetype = t[7]
			world[room].Dir_option[dir].Dcfailsave = t[8]
			world[room].Dir_option[dir].Failroom = t[9]
			world[room].Dir_option[dir].Totalfailroom = t[10]
		}
	}
}
func check_start_rooms() {
	if (func() room_rnum {
		r_mortal_start_room = real_room(config_info.Room_nums.Mortal_start_room)
		return r_mortal_start_room
	}()) == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR:  Mortal start room does not exist.  Change mortal_start_room in lib/etc/config."))
		os.Exit(1)
	}
	if (func() room_rnum {
		r_immort_start_room = real_room(config_info.Room_nums.Immort_start_room)
		return r_immort_start_room
	}()) == room_rnum(-1) {
		if mini_mud == 0 {
			basic_mud_log(libc.CString("SYSERR:  Warning: Immort start room does not exist.  Change immort_start_room in /lib/etc/config."))
		}
		r_immort_start_room = r_mortal_start_room
	}
	if (func() room_rnum {
		r_frozen_start_room = real_room(config_info.Room_nums.Frozen_start_room)
		return r_frozen_start_room
	}()) == room_rnum(-1) {
		if mini_mud == 0 {
			basic_mud_log(libc.CString("SYSERR:  Warning: Frozen start room does not exist.  Change frozen_start_room in /lib/etc/config."))
		}
		r_frozen_start_room = r_mortal_start_room
	}
}
func renum_world() {
	var (
		room int
		door int
	)
	for room = 0; room <= int(top_of_world); room++ {
		for door = 0; door < NUM_OF_DIRS; door++ {
			if world[room].Dir_option[door] != nil {
				if world[room].Dir_option[door].To_room != room_rnum(-1) {
					world[room].Dir_option[door].To_room = real_room(room_vnum(world[room].Dir_option[door].To_room))
				}
			}
		}
	}
}
func renum_zone_table() {
	var (
		cmd_no int
		a      room_rnum
		b      room_rnum
		c      room_rnum
		olda   room_rnum
	)
	_ = olda
	var oldb room_rnum
	_ = oldb
	var oldc room_rnum
	_ = oldc
	var zone zone_rnum
	var buf [128]byte
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int(zone_table[zone].Cmd[cmd_no].Command) != 'S'; cmd_no++ {
			a = func() room_rnum {
				b = func() room_rnum {
					c = 0
					return c
				}()
				return b
			}()
			olda = room_rnum(zone_table[zone].Cmd[cmd_no].Arg1)
			oldb = room_rnum(zone_table[zone].Cmd[cmd_no].Arg2)
			oldc = room_rnum(zone_table[zone].Cmd[cmd_no].Arg3)
			switch zone_table[zone].Cmd[cmd_no].Command {
			case 'M':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_mobile(mob_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
				c = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg3
					zone_table[zone].Cmd[cmd_no].Arg3 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg3)))
					return *p
				}())
			case 'O':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
				if zone_table[zone].Cmd[cmd_no].Arg3 != vnum(-1) {
					c = room_rnum(func() vnum {
						p := &zone_table[zone].Cmd[cmd_no].Arg3
						zone_table[zone].Cmd[cmd_no].Arg3 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg3)))
						return *p
					}())
				}
			case 'G':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
			case 'E':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
			case 'P':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
				c = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg3
					zone_table[zone].Cmd[cmd_no].Arg3 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg3)))
					return *p
				}())
			case 'D':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
			case 'R':
				a = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg1
					zone_table[zone].Cmd[cmd_no].Arg1 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg1)))
					return *p
				}())
				b = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg2
					zone_table[zone].Cmd[cmd_no].Arg2 = vnum(real_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg2)))
					return *p
				}())
			case 'T':
				b = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg2
					zone_table[zone].Cmd[cmd_no].Arg2 = vnum(real_trigger(trig_vnum(zone_table[zone].Cmd[cmd_no].Arg2)))
					return *p
				}())
				c = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg3
					zone_table[zone].Cmd[cmd_no].Arg3 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg3)))
					return *p
				}())
			case 'V':
				b = room_rnum(func() vnum {
					p := &zone_table[zone].Cmd[cmd_no].Arg3
					zone_table[zone].Cmd[cmd_no].Arg3 = vnum(real_room(room_vnum(zone_table[zone].Cmd[cmd_no].Arg3)))
					return *p
				}())
			}
			if a == room_rnum(-1) || b == room_rnum(-1) || c == room_rnum(-1) {
				if mini_mud == 0 {
					stdio.Snprintf(&buf[0], int(128), "Invalid vnum %lld, cmd disabled", func() room_rnum {
						if a == room_rnum(-1) {
							return olda
						}
						if b == room_rnum(-1) {
							return oldb
						}
						return oldc
					}())
					log_zone_error(zone, cmd_no, &buf[0])
				}
				zone_table[zone].Cmd[cmd_no].Command = '*'
			}
		}
	}
}
func mob_autobalance(ch *char_data) {
	ch.Hit = 0
	ch.Mana = 0
	ch.Move = 0
	ch.Exp = 0
	ch.Armor = 0
	ch.Mob_specials.Damnodice = 0
	ch.Mob_specials.Damsizedice = 0
	ch.Damage_mod = 0
}
func parse_simple_mob(mob_f *stdio.File, ch *char_data, nr int) int {
	var (
		j    int
		t    [10]int
		line [256]byte
	)
	ch.Real_abils.Str = 0
	ch.Real_abils.Intel = 0
	ch.Real_abils.Wis = 0
	ch.Real_abils.Dex = 0
	ch.Real_abils.Con = 0
	ch.Real_abils.Cha = 0
	if get_line(mob_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Format error in mob #%d, file ended after S flag!"), nr)
		return 0
	}
	if stdio.Sscanf(&line[0], " %d %d %d %dd%d+%d %dd%d+%d ", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7], &t[8]) != 9 {
		basic_mud_log(libc.CString("SYSERR: Format error in mob #%d, first line after S flag\n...expecting line of form '# # # #d#+# #d#+#'"), nr)
		return 0
	}
	ch.Race_level = t[0]
	ch.Level_adj = 0
	ch.Level = 0
	ch.Armor = (10 - t[2]) * 10
	ch.Max_hit = 0
	ch.Hit = int64(t[3])
	ch.Mana = int64(t[4])
	ch.Move = int64(t[5])
	ch.Mob_specials.Damnodice = int8(t[6])
	ch.Mob_specials.Damsizedice = int8(t[7])
	ch.Damage_mod = t[8]
	if get_line(mob_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Format error in mob #%d, second line after S flag\n...expecting line of form '# #', but file ended!"), nr)
		return 0
	}
	if stdio.Sscanf(&line[0], " %d %d %d %d", &t[0], &t[1], &t[2], &t[3]) != 4 {
		basic_mud_log(libc.CString("SYSERR: Format error in mob #%d, second line after S flag\n...expecting line of form '# # # #'"), nr)
		return 0
	}
	ch.Gold = t[0]
	ch.Exp = 0
	ch.Race = int8(t[2])
	ch.Chclass = int8(t[3])
	ch.Saving_throw[SAVING_FORTITUDE] = 0
	ch.Saving_throw[SAVING_REFLEX] = 0
	ch.Saving_throw[SAVING_WILL] = 0
	if int(ch.Race) != RACE_HUMAN {
		if !AFF_FLAGGED(ch, AFF_INFRAVISION) {
			SET_BIT_AR(ch.Affected_by[:], AFF_INFRAVISION)
		}
	}
	ch.Player_specials.Speaking = SKILL_LANG_COMMON
	if get_line(mob_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Format error in last line of mob #%d\n...expecting line of form '# # #', but file ended!"), nr)
		return 0
	}
	if stdio.Sscanf(&line[0], " %d %d %d ", &t[0], &t[1], &t[2]) != 3 {
		basic_mud_log(libc.CString("SYSERR: Format error in last line of mob #%d\n...expecting line of form '# # #'"), nr)
		return 0
	}
	ch.Position = int8(t[0])
	ch.Mob_specials.Default_pos = int8(t[1])
	ch.Sex = int8(t[2])
	ch.Player_specials.Speaking = SKILL_LANG_COMMON
	set_height_and_weight_by_race(ch)
	for j = 0; j < 3; j++ {
		ch.Apply_saving_throw[j] = 0
	}
	if MOB_FLAGGED(ch, MOB_AUTOBALANCE) {
		mob_autobalance(ch)
	}
	return 1
}
func interpret_espec(keyword *byte, value *byte, ch *char_data, nr int) {
	var (
		num_arg int = 0
		matched int = FALSE
		num     int
		num2    int
		num3    int
		num4    int
		num5    int
		num6    int
		af      affected_type
	)
	if value != nil {
		num_arg = libc.Atoi(libc.GoString(value))
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("BareHandAttack")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99, int64(num_arg))))
		ch.Mob_specials.Attack_type = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Size")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(int64(-1), MIN(int64(int(NUM_SIZES-1)), int64(num_arg))))
		ch.Size = num_arg
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Str")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Str = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("StrAdd")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		basic_mud_log(libc.CString("mob #%d trying to set StrAdd, rebalance its strength."), GET_MOB_VNUM(ch))
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Int")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Intel = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Wis")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Wis = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Dex")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Dex = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Con")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Con = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Cha")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(200, int64(num_arg))))
		ch.Real_abils.Cha = int8(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Hit")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Hit = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("MaxHit")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Max_hit = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Mana")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Mana = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("MaxMana")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Max_mana = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Moves")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Move = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("MaxMoves")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num_arg = int(MAX(0, MIN(99999, int64(num_arg))))
		ch.Max_move = int64(num_arg)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Affect")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num = func() int {
			num2 = func() int {
				num3 = func() int {
					num4 = func() int {
						num5 = func() int {
							num6 = 0
							return num6
						}()
						return num5
					}()
					return num4
				}()
				return num3
			}()
			return num2
		}()
		stdio.Sscanf(value, "%d %d %d %d %d %d", &num, &num2, &num3, &num4, &num5, &num6)
		if num > 0 {
			af.Type = int16(num)
			af.Duration = int16(num2)
			af.Modifier = num3
			af.Location = num4
			af.Bitvector = bitvector_t(int32(num5))
			af.Specific = num6
			affect_to_char(ch, &af)
		}
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("AffectV")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		num = func() int {
			num2 = func() int {
				num3 = func() int {
					num4 = func() int {
						num5 = func() int {
							num6 = 0
							return num6
						}()
						return num5
					}()
					return num4
				}()
				return num3
			}()
			return num2
		}()
		stdio.Sscanf(value, "%d %d %d %d %d %d", &num, &num2, &num3, &num4, &num5, &num6)
		if num > 0 {
			af.Type = int16(num)
			af.Duration = int16(num2)
			af.Modifier = num3
			af.Location = num4
			af.Bitvector = bitvector_t(int32(num5))
			af.Specific = num6
			affectv_to_char(ch, &af)
		}
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Feat")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		stdio.Sscanf(value, "%d %d", &num, &num2)
		ch.Feats[num] = int8(num2)
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Skill")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		stdio.Sscanf(value, "%d %d", &num, &num2)
		for {
			ch.Skills[num] = int8(num2)
			if true {
				break
			}
		}
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("SkillMod")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		stdio.Sscanf(value, "%d %d", &num, &num2)
		for {
			ch.Skillmods[num] = int8(num2)
			if true {
				break
			}
		}
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("Class")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		stdio.Sscanf(value, "%d %d", &num, &num2)
		ch.Chclasses[num] = num2
		ch.Level += num2
	}
	if value != nil && matched == 0 && libc.StrCaseCmp(keyword, libc.CString("EpicClass")) == 0 && (func() int {
		matched = TRUE
		return matched
	}()) != 0 {
		stdio.Sscanf(value, "%d %d", &num, &num2)
		ch.Epicclasses[num] = num2
		ch.Level += num2
	}
	if matched == 0 {
		basic_mud_log(libc.CString("SYSERR: Warning: unrecognized espec keyword %s in mob #%d"), keyword, nr)
	}
}
func parse_espec(buf *byte, ch *char_data, nr int) {
	var ptr *byte
	if (func() *byte {
		ptr = libc.StrChr(buf, ':')
		return ptr
	}()) != nil {
		*(func() *byte {
			p := &ptr
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = '\x00'
		for unicode.IsSpace(rune(*ptr)) {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
	}
	interpret_espec(buf, ptr, ch, nr)
}
func parse_enhanced_mob(mob_f *stdio.File, ch *char_data, nr int) int {
	var line [256]byte
	parse_simple_mob(mob_f, ch, nr)
	for get_line(mob_f, &line[0]) != 0 {
		if libc.StrCmp(&line[0], libc.CString("E")) == 0 {
			return 1
		} else if line[0] == '#' {
			basic_mud_log(libc.CString("SYSERR: Unterminated E section in mob #%d"), nr)
			return 0
		} else {
			parse_espec(&line[0], ch, nr)
		}
	}
	basic_mud_log(libc.CString("SYSERR: Unexpected end of file reached after mob #%d"), nr)
	return 0
}
func parse_mobile_from_file(mob_f *stdio.File, ch *char_data) int {
	var (
		j      int
		t      [10]int
		retval int
		line   [256]byte
		tmpptr *byte
		letter int8
		f1     [128]byte
		f2     [128]byte
		f3     [128]byte
		f4     [128]byte
		f5     [128]byte
		f6     [128]byte
		f7     [128]byte
		f8     [128]byte
		buf2   [128]byte
		nr     mob_vnum = mob_index[ch.Nr].Vnum
	)
	ch.Player_specials = &dummy_mob
	stdio.Sprintf(&buf2[0], "mob vnum %d", nr)
	ch.Name = fread_string(mob_f, &buf2[0])
	tmpptr = func() *byte {
		p := &ch.Short_descr
		ch.Short_descr = fread_string(mob_f, &buf2[0])
		return *p
	}()
	if tmpptr != nil && *tmpptr != 0 {
		if libc.StrCaseCmp(fname(tmpptr), libc.CString("a")) == 0 || libc.StrCaseCmp(fname(tmpptr), libc.CString("an")) == 0 || libc.StrCaseCmp(fname(tmpptr), libc.CString("the")) == 0 {
			*tmpptr = byte(int8(unicode.ToLower(rune(*tmpptr))))
		}
	}
	ch.Long_descr = fread_string(mob_f, &buf2[0])
	ch.Description = fread_string(mob_f, &buf2[0])
	if get_line(mob_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Format error after string section of mob #%d\n...expecting line of form '# # # {S | E}', but file ended!"), nr)
		return 0
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], "%s %s %s %s %s %s %s %s %d %c", &f1[0], &f2[0], &f3[0], &f4[0], &f5[0], &f6[0], &f7[0], &f8[0], &t[2], &letter)
		return retval
	}()) == 10 && bitwarning == TRUE {
		basic_mud_log(libc.CString("WARNING: Conventional mobilefiles detected. Please read 128bit.readme."))
		return 0
	} else if retval == 4 && bitwarning == FALSE {
		basic_mud_log(libc.CString("Converting mobile #%d to 128bits.."), nr)
		ch.Act[0] = asciiflag_conv(&f1[0])
		ch.Act[1] = 0
		ch.Act[2] = 0
		ch.Act[3] = 0
		check_bitvector_names(ch.Act[0], action_bits_count, &buf2[0], libc.CString("mobile"))
		ch.Affected_by[0] = asciiflag_conv_aff(&f2[0])
		ch.Affected_by[1] = 0
		ch.Affected_by[2] = 0
		ch.Affected_by[3] = 0
		ch.Alignment = libc.Atoi(libc.GoString(&f3[0]))
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_CHARM)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_POISON)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_GROUP)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_SLEEP)
		if MOB_FLAGGED(ch, MOB_AGGRESSIVE) && MOB_FLAGGED(ch, MOB_AGGR_GOOD) {
			REMOVE_BIT_AR(ch.Act[:], MOB_AGGR_GOOD)
		}
		if MOB_FLAGGED(ch, MOB_AGGRESSIVE) && MOB_FLAGGED(ch, MOB_AGGR_NEUTRAL) {
			REMOVE_BIT_AR(ch.Act[:], MOB_AGGR_NEUTRAL)
		}
		if MOB_FLAGGED(ch, MOB_AGGRESSIVE) && MOB_FLAGGED(ch, MOB_AGGR_EVIL) {
			REMOVE_BIT_AR(ch.Act[:], MOB_AGGR_EVIL)
		}
		check_bitvector_names(ch.Affected_by[0], affected_bits_count, &buf2[0], libc.CString("mobile affect"))
		letter = int8(f4[0])
		if bitsavetodisk != 0 {
			add_to_save_list(zone_table[real_zone_by_thing(room_vnum(nr))].Number, 0)
			converting = TRUE
		}
		basic_mud_log(libc.CString("   done."))
	} else if retval == 10 {
		var taeller int
		ch.Act[0] = asciiflag_conv(&f1[0])
		ch.Act[1] = asciiflag_conv(&f2[0])
		ch.Act[2] = asciiflag_conv(&f3[0])
		ch.Act[3] = asciiflag_conv(&f4[0])
		for taeller = 0; taeller < AF_ARRAY_MAX; taeller++ {
			check_bitvector_names(ch.Act[taeller], action_bits_count, &buf2[0], libc.CString("mobile"))
		}
		ch.Affected_by[0] = asciiflag_conv(&f5[0])
		ch.Affected_by[1] = asciiflag_conv(&f6[0])
		ch.Affected_by[2] = asciiflag_conv(&f7[0])
		ch.Affected_by[3] = asciiflag_conv(&f8[0])
		ch.Alignment = t[2]
		for taeller = 0; taeller < AF_ARRAY_MAX; taeller++ {
			check_bitvector_names(ch.Affected_by[taeller], affected_bits_count, &buf2[0], libc.CString("mobile affect"))
		}
	} else {
		basic_mud_log(libc.CString("SYSERR: Format error after string section of mob #%d\n...expecting line of form '# # # {S | E}'"), nr)
		os.Exit(1)
	}
	SET_BIT_AR(ch.Act[:], MOB_ISNPC)
	if MOB_FLAGGED(ch, MOB_NOTDEADYET) {
		basic_mud_log(libc.CString("SYSERR: Mob #%d has reserved bit MOB_NOTDEADYET set."), nr)
		REMOVE_BIT_AR(ch.Act[:], MOB_NOTDEADYET)
	}
	switch unicode.ToUpper(rune(letter)) {
	case 'S':
		parse_simple_mob(mob_f, ch, int(nr))
	case 'E':
		parse_enhanced_mob(mob_f, ch, int(nr))
		mob_stats(ch)
	default:
		basic_mud_log(libc.CString("SYSERR: Unsupported mob type '%c' in mob #%d"), letter, nr)
		os.Exit(1)
	}
	letter = fread_letter(mob_f)
	mob_f.UnGetC(int(letter))
	for int(letter) == 'T' {
		dg_read_trigger(mob_f, unsafe.Pointer(ch), MOB_TRIGGER)
		letter = fread_letter(mob_f)
		mob_f.UnGetC(int(letter))
	}
	ch.Aff_abils = ch.Real_abils
	for j = 0; j < NUM_WEARS; j++ {
		ch.Equipment[j] = nil
	}
	return 1
}
func parse_mobile(mob_f *stdio.File, nr int) {
	var i int = 0
	mob_index[i].Vnum = mob_vnum(nr)
	mob_index[i].Number = 0
	mob_index[i].Func = nil
	clear_char(&mob_proto[i])
	mob_proto[i].Nr = mob_rnum(i)
	mob_proto[i].Desc = nil
	if parse_mobile_from_file(mob_f, &mob_proto[i]) != 0 {
		if mob_htree == nil {
			mob_htree = htree_init()
		}
		htree_add(mob_htree, int64(nr), int64(i))
		top_of_mobt = mob_rnum(func() int {
			p := &i
			x := *p
			*p++
			return x
		}())
	} else {
		os.Exit(1)
	}
}
func parse_object(obj_f *stdio.File, nr int) *byte {
	var (
		i         int = 0
		line      [256]byte
		t         [18]int
		j         int
		retval    int
		tmpptr    *byte
		buf2      [128]byte
		f1        [256]byte
		f2        [256]byte
		f3        [256]byte
		f4        [256]byte
		f5        [256]byte
		f6        [256]byte
		f7        [256]byte
		f8        [256]byte
		f9        [256]byte
		f10       [256]byte
		f11       [256]byte
		f12       [256]byte
		new_descr *extra_descr_data
	)
	obj_index[i].Vnum = mob_vnum(nr)
	obj_index[i].Number = 0
	obj_index[i].Func = nil
	if obj_htree == nil {
		obj_htree = htree_init()
	}
	htree_add(obj_htree, int64(nr), int64(i))
	clear_object(&obj_proto[i])
	obj_proto[i].Item_number = obj_vnum(i)
	stdio.Sprintf(&buf2[0], "object #%d", nr)
	if (func() *byte {
		p := &obj_proto[i].Name
		obj_proto[i].Name = fread_string(obj_f, &buf2[0])
		return *p
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Null obj name or format error at or near %s"), &buf2[0])
		os.Exit(1)
	}
	tmpptr = func() *byte {
		p := &obj_proto[i].Short_description
		obj_proto[i].Short_description = fread_string(obj_f, &buf2[0])
		return *p
	}()
	if tmpptr != nil && *tmpptr != 0 {
		if libc.StrCaseCmp(fname(tmpptr), libc.CString("a")) == 0 || libc.StrCaseCmp(fname(tmpptr), libc.CString("an")) == 0 || libc.StrCaseCmp(fname(tmpptr), libc.CString("the")) == 0 {
			*tmpptr = byte(int8(unicode.ToLower(rune(*tmpptr))))
		}
	}
	tmpptr = func() *byte {
		p := &obj_proto[i].Description
		obj_proto[i].Description = fread_string(obj_f, &buf2[0])
		return *p
	}()
	if tmpptr != nil && *tmpptr != 0 {
		CAP(tmpptr)
	}
	obj_proto[i].Action_description = fread_string(obj_f, &buf2[0])
	if get_line(obj_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Expecting first numeric line of %s, but file ended!"), &buf2[0])
		os.Exit(1)
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], " %d %s %s %s %s %s %s %s %s %s %s %s %s", &t[0], &f1[0], &f2[0], &f3[0], &f4[0], &f5[0], &f6[0], &f7[0], &f8[0], &f9[0], &f10[0], &f11[0], &f12[0])
		return retval
	}()) == 4 && bitwarning == TRUE {
		basic_mud_log(libc.CString("WARNING: Conventional objectfiles detected. Please read 128bit.readme."))
		os.Exit(1)
	} else if (retval == 4 || retval == 3) && bitwarning == FALSE {
		if retval == 3 {
			t[3] = 0
		} else if retval == 4 {
			t[3] = int(asciiflag_conv_aff(&f3[0]))
		}
		basic_mud_log(libc.CString("Converting object #%d to 128bits.."), nr)
		obj_proto[i].Extra_flags[0] = asciiflag_conv(&f1[0])
		obj_proto[i].Extra_flags[1] = 0
		obj_proto[i].Extra_flags[2] = 0
		obj_proto[i].Extra_flags[3] = 0
		obj_proto[i].Wear_flags[0] = asciiflag_conv(&f2[0])
		obj_proto[i].Wear_flags[1] = 0
		obj_proto[i].Wear_flags[2] = 0
		obj_proto[i].Wear_flags[3] = 0
		obj_proto[i].Bitvector[0] = asciiflag_conv_aff(&f3[0])
		obj_proto[i].Bitvector[1] = 0
		obj_proto[i].Bitvector[2] = 0
		obj_proto[i].Bitvector[3] = 0
		if bitsavetodisk != 0 {
			add_to_save_list(zone_table[real_zone_by_thing(room_vnum(nr))].Number, 1)
			converting = TRUE
		}
		basic_mud_log(libc.CString("   done."))
	} else if retval == 13 {
		obj_proto[i].Extra_flags[0] = asciiflag_conv(&f1[0])
		obj_proto[i].Extra_flags[1] = asciiflag_conv(&f2[0])
		obj_proto[i].Extra_flags[2] = asciiflag_conv(&f3[0])
		obj_proto[i].Extra_flags[3] = asciiflag_conv(&f4[0])
		obj_proto[i].Wear_flags[0] = asciiflag_conv(&f5[0])
		obj_proto[i].Wear_flags[1] = asciiflag_conv(&f6[0])
		obj_proto[i].Wear_flags[2] = asciiflag_conv(&f7[0])
		obj_proto[i].Wear_flags[3] = asciiflag_conv(&f8[0])
		obj_proto[i].Bitvector[0] = asciiflag_conv(&f9[0])
		obj_proto[i].Bitvector[1] = asciiflag_conv(&f10[0])
		obj_proto[i].Bitvector[2] = asciiflag_conv(&f11[0])
		obj_proto[i].Bitvector[3] = asciiflag_conv(&f12[0])
	} else {
		basic_mud_log(libc.CString("SYSERR: Format error in first numeric line (expecting 13 args, got %d), %s"), retval, &buf2[0])
		os.Exit(1)
	}
	obj_proto[i].Type_flag = int8(t[0])
	if get_line(obj_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Expecting second numeric line of %s, but file ended!"), &buf2[0])
		os.Exit(1)
	}
	for j = 0; j < NUM_OBJ_VAL_POSITIONS; j++ {
		t[j] = 0
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], "%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7], &t[8], &t[9], &t[10], &t[11], &t[12], &t[13], &t[14], &t[15])
		return retval
	}()) > NUM_OBJ_VAL_POSITIONS {
		basic_mud_log(libc.CString("SYSERR: Format error in second numeric line (expecting <=%d args, got %d), %s"), NUM_OBJ_VAL_POSITIONS, retval, &buf2[0])
		os.Exit(1)
	}
	for j = 0; j < NUM_OBJ_VAL_POSITIONS; j++ {
		obj_proto[i].Value[j] = t[j]
	}
	if (int(obj_proto[i].Type_flag) == ITEM_PORTAL || int(obj_proto[i].Type_flag) == ITEM_HATCH) && ((obj_proto[i].Value[VAL_DOOR_DCLOCK]) == 0 || (obj_proto[i].Value[VAL_DOOR_DCHIDE]) == 0) {
		obj_proto[i].Value[VAL_DOOR_DCLOCK] = 20
		obj_proto[i].Value[VAL_DOOR_DCHIDE] = 20
		if bitsavetodisk != 0 {
			add_to_save_list(zone_table[real_zone_by_thing(room_vnum(nr))].Number, 1)
			converting = TRUE
		}
	}
	if int(obj_proto[i].Type_flag) == ITEM_WEAPON && (obj_proto[i].Value[0]) > 169 {
		obj_proto[i].Value[0] = suntzu_weapon_convert(t[0])
		if bitsavetodisk != 0 {
			add_to_save_list(zone_table[real_zone_by_thing(room_vnum(nr))].Number, 1)
			converting = TRUE
		}
	}
	if get_line(obj_f, &line[0]) == 0 {
		basic_mud_log(libc.CString("SYSERR: Expecting third numeric line of %s, but file ended!"), &buf2[0])
		os.Exit(1)
	}
	if (func() int {
		retval = stdio.Sscanf(&line[0], "%d %d %d %d", &t[0], &t[1], &t[2], &t[3])
		return retval
	}()) != 4 {
		if retval == 3 {
			t[3] = 0
		} else {
			basic_mud_log(libc.CString("SYSERR: Format error in third numeric line (expecting 4 args, got %d), %s"), retval, &buf2[0])
			os.Exit(1)
		}
	}
	obj_proto[i].Weight = int64(t[0])
	obj_proto[i].Cost = t[1]
	obj_proto[i].Cost_per_day = t[2]
	obj_proto[i].Level = t[3]
	obj_proto[i].Size = SIZE_MEDIUM
	if int(obj_proto[i].Type_flag) == ITEM_DRINKCON || int(obj_proto[i].Type_flag) == ITEM_FOUNTAIN {
		if obj_proto[i].Weight < int64(obj_proto[i].Value[1]) {
			obj_proto[i].Weight = int64((obj_proto[i].Value[1]) + 5)
		}
	}
	if int(obj_proto[i].Type_flag) == ITEM_PORTAL {
		obj_proto[i].Timer = -1
	}
	for j = 0; j < MAX_OBJ_AFFECT; j++ {
		obj_proto[i].Affected[j].Location = APPLY_NONE
		obj_proto[i].Affected[j].Modifier = 0
		obj_proto[i].Affected[j].Specific = 0
	}
	libc.StrCat(&buf2[0], libc.CString(", after numeric constants\n...expecting 'E', 'A', '$', or next object number"))
	j = 0
	for {
		if get_line(obj_f, &line[0]) == 0 {
			basic_mud_log(libc.CString("SYSERR: Format error in %s"), &buf2[0])
			os.Exit(1)
		}
		switch line[0] {
		case 'E':
			new_descr = new(extra_descr_data)
			new_descr.Keyword = fread_string(obj_f, &buf2[0])
			new_descr.Description = fread_string(obj_f, &buf2[0])
			new_descr.Next = obj_proto[i].Ex_description
			obj_proto[i].Ex_description = new_descr
		case 'A':
			if j >= MAX_OBJ_AFFECT {
				basic_mud_log(libc.CString("SYSERR: Too many A fields (%d max), %s"), MAX_OBJ_AFFECT, &buf2[0])
				os.Exit(1)
			}
			if get_line(obj_f, &line[0]) == 0 {
				basic_mud_log(libc.CString("SYSERR: Format error in 'A' field, %s\n...expecting 2 numeric constants but file ended!"), &buf2[0])
				os.Exit(1)
			}
			t[1] = 0
			if (func() int {
				retval = stdio.Sscanf(&line[0], " %d %d %d ", &t[0], &t[1], &t[2])
				return retval
			}()) != 3 {
				if retval != 2 {
					basic_mud_log(libc.CString("SYSERR: Format error in 'A' field, %s\n...expecting 2 numeric arguments, got %d\n...offending line: '%s'"), &buf2[0], retval, &line[0])
					os.Exit(1)
				}
			}
			if t[0] >= APPLY_UNUSED3 && t[0] <= APPLY_UNUSED4 {
				basic_mud_log(libc.CString("Warning: object #%d (%s) uses deprecated saving throw applies"), nr, obj_proto[i].Short_description)
			}
			obj_proto[i].Affected[j].Location = t[0]
			obj_proto[i].Affected[j].Modifier = t[1]
			obj_proto[i].Affected[j].Specific = t[2]
			j++
		case 'S':
			if j >= SPELLBOOK_SIZE {
				basic_mud_log(libc.CString("SYSERR: Unknown spellbook slot in S field, %s"), &buf2[0])
				os.Exit(1)
			}
			if get_line(obj_f, &line[0]) == 0 {
				basic_mud_log(libc.CString("SYSERR: Format error in 'S' field, %s\n...expecting 2 numeric constants but file ended!"), &buf2[0])
				os.Exit(1)
			}
			if (func() int {
				retval = stdio.Sscanf(&line[0], " %d %d ", &t[0], &t[1])
				return retval
			}()) != 2 {
				basic_mud_log(libc.CString("SYSERR: Format error in 'S' field, %s\n...expecting 2 numeric arguments, got %d\n...offending line: '%s'"), &buf2[0], retval, &line[0])
				os.Exit(1)
			}
			if obj_proto[i].Sbinfo == nil {
				obj_proto[i].Sbinfo = make([]obj_spellbook_spell, SPELLBOOK_SIZE)
				libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(&obj_proto[i].Sbinfo[0]))), 0, int(SPELLBOOK_SIZE*unsafe.Sizeof(obj_spellbook_spell{})))
			}
			obj_proto[i].Sbinfo[j].Spellname = t[0]
			obj_proto[i].Sbinfo[j].Pages = t[1]
			j++
		case 'T':
			dg_obj_trigger(&line[0], &obj_proto[i])
		case 'Z':
			if get_line(obj_f, &line[0]) == 0 {
				basic_mud_log(libc.CString("SYSERR: Format error in 'Z' field, %s\n...expecting numeric constant but file ended!"), &buf2[0])
				os.Exit(1)
			}
			if stdio.Sscanf(&line[0], "%d", &t[0]) != 1 {
				basic_mud_log(libc.CString("SYSERR: Format error in 'Z' field, %s\n...expecting numeric argument\n...offending line: '%s'"), &buf2[0], &line[0])
				os.Exit(1)
			}
			obj_proto[i].Size = t[0]
		case '$':
			fallthrough
		case '#':
			if OBJAFF_FLAGGED(&obj_proto[i], AFF_CHARM) {
				basic_mud_log(libc.CString("SYSERR: Object #%d has reserved bit AFF_CHARM set."), nr)
				REMOVE_BIT_AR(obj_proto[i].Bitvector[:], AFF_CHARM)
			}
			top_of_objt = obj_rnum(i)
			check_object(&obj_proto[i])
			i++
			return &line[0]
		default:
			basic_mud_log(libc.CString("SYSERR: Format error in (%c): %s"), line[0], &buf2[0])
			os.Exit(1)
		}
	}
	return &line[0]
}
func load_zones(fl *stdio.File, zonename *byte) {
	var (
		zone        zone_rnum = 0
		cmd_no      int
		num_of_cmds int = 0
		line_num    int = 0
		tmp         int
		error       int
		arg_num     int
		version     int = 1
		ptr         *byte
		buf         [256]byte
		zname       [256]byte
		buf2        [64936]byte
		zone_fix    int = FALSE
		t1          [80]byte
		t2          [80]byte
		line        [64936]byte
	)
	strlcpy(&zname[0], zonename, uint64(256))
	for tmp = 0; tmp < 3; tmp++ {
		get_line(fl, &buf[0])
	}
	for get_line(fl, &buf[0]) != 0 {
		if libc.StrChr(libc.CString("MOPGERDTV"), buf[0]) != nil && buf[1] == ' ' || buf[0] == 'S' && buf[1] == '\x00' {
			num_of_cmds++
		}
	}
	fl.Seek(0, 0)
	if num_of_cmds == 0 {
		basic_mud_log(libc.CString("SYSERR: %s is empty!"), &zname[0])
		os.Exit(1)
	} else {
		zone_table[zone].Cmd = make([]reset_com, num_of_cmds)
	}
	line_num += get_line(fl, &buf[0])
	if buf[0] == '@' {
		if stdio.Sscanf(&buf[0], "@Version: %d", &version) != 1 {
			basic_mud_log(libc.CString("SYSERR: Format error in %s (version)"), &zname[0])
			basic_mud_log(libc.CString("SYSERR: ...Line: %s"), &line[0])
			os.Exit(1)
		}
		line_num += get_line(fl, &buf[0])
	}
	if stdio.Sscanf(&buf[0], "#%hd", &zone_table[zone].Number) != 1 {
		basic_mud_log(libc.CString("SYSERR: Format error in %s, line %d"), &zname[0], line_num)
		os.Exit(1)
	}
	stdio.Snprintf(&buf2[0], int(64936), "beginning of zone #%d", zone_table[zone].Number)
	line_num += get_line(fl, &buf[0])
	if (func() *byte {
		ptr = libc.StrChr(&buf[0], '~')
		return ptr
	}()) != nil {
		*ptr = '\x00'
	}
	zone_table[zone].Builders = libc.StrDup(&buf[0])
	line_num += get_line(fl, &buf[0])
	if (func() *byte {
		ptr = libc.StrChr(&buf[0], '~')
		return ptr
	}()) != nil {
		*ptr = '\x00'
	}
	zone_table[zone].Name = libc.StrDup(&buf[0])
	line_num += get_line(fl, &buf[0])
	if version >= 2 {
		var (
			zbuf1 [64936]byte
			zbuf2 [64936]byte
			zbuf3 [64936]byte
			zbuf4 [64936]byte
		)
		if stdio.Sscanf(&buf[0], " %hd %hd %d %d %s %s %s %s %d %d", &zone_table[zone].Bot, &zone_table[zone].Top, &zone_table[zone].Lifespan, &zone_table[zone].Reset_mode, &zbuf1[0], &zbuf2[0], &zbuf3[0], &zbuf4[0], &zone_table[zone].Min_level, &zone_table[zone].Max_level) != 10 {
			basic_mud_log(libc.CString("SYSERR: Format error in 10-constant line of %s"), &zname[0])
			os.Exit(1)
		}
		zone_table[zone].Zone_flags[0] = asciiflag_conv(&zbuf1[0])
		zone_table[zone].Zone_flags[1] = asciiflag_conv(&zbuf2[0])
		zone_table[zone].Zone_flags[2] = asciiflag_conv(&zbuf3[0])
		zone_table[zone].Zone_flags[3] = asciiflag_conv(&zbuf4[0])
	} else if stdio.Sscanf(&buf[0], " %hd %hd %d %d ", &zone_table[zone].Bot, &zone_table[zone].Top, &zone_table[zone].Lifespan, &zone_table[zone].Reset_mode) != 4 {
		basic_mud_log(libc.CString("SYSERR: Format error in numeric constant line of %s, attempting to fix."), &zname[0])
		if stdio.Sscanf(zone_table[zone].Name, " %hd %hd %d %d ", &zone_table[zone].Bot, &zone_table[zone].Top, &zone_table[zone].Lifespan, &zone_table[zone].Reset_mode) != 4 {
			basic_mud_log(libc.CString("SYSERR: Could not fix previous error, aborting game."))
			os.Exit(1)
		} else {
			libc.Free(unsafe.Pointer(zone_table[zone].Name))
			zone_table[zone].Name = libc.StrDup(zone_table[zone].Builders)
			libc.Free(unsafe.Pointer(zone_table[zone].Builders))
			zone_table[zone].Builders = libc.CString("None.")
			zone_fix = TRUE
		}
	}
	if zone_table[zone].Bot > zone_table[zone].Top {
		basic_mud_log(libc.CString("SYSERR: Zone %d bottom (%d) > top (%d)."), zone_table[zone].Number, zone_table[zone].Bot, zone_table[zone].Top)
		os.Exit(1)
	}
	cmd_no = 0
	for {
		if zone_fix != TRUE {
			if (func() int {
				tmp = get_line(fl, &buf[0])
				return tmp
			}()) == 0 {
				basic_mud_log(libc.CString("SYSERR: Format error in %s - premature end of file"), &zname[0])
				os.Exit(1)
			}
		} else {
			zone_fix = FALSE
		}
		line_num += tmp
		ptr = &buf[0]
		skip_spaces(&ptr)
		if int(func() int8 {
			p := &zone_table[zone].Cmd[cmd_no].Command
			zone_table[zone].Cmd[cmd_no].Command = int8(*ptr)
			return *p
		}()) == '*' {
			continue
		}
		ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		if int(zone_table[zone].Cmd[cmd_no].Command) == 'S' || int(zone_table[zone].Cmd[cmd_no].Command) == '$' {
			zone_table[zone].Cmd[cmd_no].Command = 'S'
			break
		}
		error = 0
		if libc.StrChr(libc.CString("MOEPDTVG"), byte(zone_table[zone].Cmd[cmd_no].Command)) == nil {
			if stdio.Sscanf(ptr, " %d %d %d %d ", &tmp, &zone_table[zone].Cmd[cmd_no].Arg1, &zone_table[zone].Cmd[cmd_no].Arg2, &zone_table[zone].Cmd[cmd_no].Arg3) != 4 {
				error = 1
			}
		} else if int(zone_table[zone].Cmd[cmd_no].Command) == 'V' {
			if stdio.Sscanf(ptr, " %d %d %d %d %d %d %79s %79[^\f\n\r\t\v]", &tmp, &zone_table[zone].Cmd[cmd_no].Arg1, &zone_table[zone].Cmd[cmd_no].Arg2, &zone_table[zone].Cmd[cmd_no].Arg3, &zone_table[zone].Cmd[cmd_no].Arg4, &zone_table[zone].Cmd[cmd_no].Arg5, &t1[0], &t2[0]) != 8 {
				error = 1
			} else {
				zone_table[zone].Cmd[cmd_no].Sarg1 = libc.StrDup(&t1[0])
				zone_table[zone].Cmd[cmd_no].Sarg2 = libc.StrDup(&t2[0])
			}
		} else {
			if (func() int {
				arg_num = stdio.Sscanf(ptr, " %d %d %d %d %d %d ", &tmp, &zone_table[zone].Cmd[cmd_no].Arg1, &zone_table[zone].Cmd[cmd_no].Arg2, &zone_table[zone].Cmd[cmd_no].Arg3, &zone_table[zone].Cmd[cmd_no].Arg4, &zone_table[zone].Cmd[cmd_no].Arg5)
				return arg_num
			}()) != 6 {
				if arg_num != 5 {
					error = 1
				} else {
					zone_table[zone].Cmd[cmd_no].Arg5 = 0
				}
			}
		}
		zone_table[zone].Cmd[cmd_no].If_flag = tmp != 0
		if error != 0 {
			basic_mud_log(libc.CString("SYSERR: Format error in %s, line %d: '%s'"), &zname[0], line_num, &buf[0])
			os.Exit(1)
		}
		zone_table[zone].Cmd[cmd_no].Line = line_num
		cmd_no++
	}
	if num_of_cmds != cmd_no+1 {
		basic_mud_log(libc.CString("SYSERR: Zone command count mismatch for %s. Estimated: %d, Actual: %d"), &zname[0], num_of_cmds, cmd_no+1)
		os.Exit(1)
	}
	top_of_zone_table = func() zone_rnum {
		p := &zone
		x := *p
		*p++
		return x
	}()
}
func get_one_line(fl *stdio.File, buf *byte) {
	if fl.GetS(buf, READ_SIZE) == nil {
		basic_mud_log(libc.CString("SYSERR: error reading help file: not terminated with $?"))
		os.Exit(1)
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(buf), libc.StrLen(buf)-1)) = '\x00'
}
func free_help(cmhelp *help_index_element) {
	if cmhelp.Keywords != nil {
		libc.Free(unsafe.Pointer(cmhelp.Keywords))
	}
	if cmhelp.Entry != nil && cmhelp.Duplicate == 0 {
		libc.Free(unsafe.Pointer(cmhelp.Entry))
	}
	libc.Free(unsafe.Pointer(cmhelp))
}
func free_help_table() {
	if help_table != nil {
		var hp int
		for hp = 0; hp < top_of_helpt; hp++ {
			if help_table[hp].Keywords != nil {
				libc.Free(unsafe.Pointer(help_table[hp].Keywords))
			}
			if help_table[hp].Entry != nil && help_table[hp].Duplicate == 0 {
				libc.Free(unsafe.Pointer(help_table[hp].Entry))
			}
		}
		libc.Free(unsafe.Pointer(&help_table[0]))
		help_table = nil
	}
	top_of_helpt = 0
}
func load_help(fl *stdio.File, name *byte) {
	var (
		key      [257]byte
		next_key [257]byte
		entry    [32384]byte
		entrylen uint64
		line     [257]byte
		hname    [257]byte
		scan     *byte
		el       help_index_element
	)
	strlcpy(&hname[0], name, uint64(257))
	get_one_line(fl, &key[0])
	for key[0] != '$' {
		libc.StrCat(&key[0], libc.CString("\r\n"))
		entrylen = strlcpy(&entry[0], &key[0], uint64(32384))
		get_one_line(fl, &line[0])
		for line[0] != '#' && entrylen < uint64(32384-1) {
			entrylen += strlcpy(&entry[entrylen], &line[0], uint64(32384-uintptr(entrylen)))
			if entrylen+2 < uint64(32384-1) {
				libc.StrCpy(&entry[entrylen], libc.CString("\r\n"))
				entrylen += 2
			}
			get_one_line(fl, &line[0])
		}
		if entrylen >= uint64(32384-1) {
			var (
				keysize  int
				truncmsg *byte = libc.CString("\r\n*TRUNCATED*\r\n")
			)
			libc.StrCpy((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&entry[32383]), -libc.StrLen(truncmsg)))), -1)), truncmsg)
			keysize = libc.StrLen(&key[0]) - 2
			basic_mud_log(libc.CString("SYSERR: Help entry exceeded buffer space: %.*s"), keysize, &key[0])
			for line[0] != '#' {
				get_one_line(fl, &line[0])
			}
		}
		if line[0] == '#' {
			if stdio.Sscanf(&line[0], "#%d", &el.Min_level) != 1 {
				basic_mud_log(libc.CString("SYSERR: Help entry does not have a min level. %s"), &key[0])
				el.Min_level = 0
			}
		}
		el.Duplicate = 0
		el.Entry = libc.StrDup(&entry[0])
		scan = one_word(&key[0], &next_key[0])
		for next_key[0] != 0 {
			el.Keywords = libc.StrDup(&next_key[0])
			help_table[func() int {
				p := &top_of_helpt
				x := *p
				*p++
				return x
			}()] = el
			el.Duplicate++
			scan = one_word(scan, &next_key[0])
		}
		get_one_line(fl, &key[0])
	}
}
func hsort(a unsafe.Pointer, b unsafe.Pointer) int {
	var (
		a1 *help_index_element
		b1 *help_index_element
	)
	a1 = (*help_index_element)(a)
	b1 = (*help_index_element)(b)
	return libc.StrCaseCmp(a1.Keywords, b1.Keywords)
}
func vnum_mobile(searchname *byte, ch *char_data) int {
	var (
		nr    int
		found int = 0
	)
	for nr = 0; nr <= int(top_of_mobt); nr++ {
		if isname(searchname, mob_proto[nr].Name) != 0 {
			send_to_char(ch, libc.CString("%3d. [%5d] %-40s %s\r\n"), func() int {
				p := &found
				*p++
				return *p
			}(), mob_index[nr].Vnum, mob_proto[nr].Short_descr, func() string {
				if mob_proto[nr].Proto_script != nil {
					return "[TRIG]"
				}
				return ""
			}())
		}
	}
	return found
}
func vnum_object(searchname *byte, ch *char_data) int {
	var (
		nr    int
		found int = 0
	)
	for nr = 0; nr <= int(top_of_objt); nr++ {
		if isname(searchname, obj_proto[nr].Name) != 0 {
			send_to_char(ch, libc.CString("%3d. [%5d] %-40s %s\r\n"), func() int {
				p := &found
				*p++
				return *p
			}(), obj_index[nr].Vnum, obj_proto[nr].Short_description, func() string {
				if obj_proto[nr].Proto_script != nil {
					return "[TRIG]"
				}
				return ""
			}())
		}
	}
	return found
}
func vnum_material(searchname *byte, ch *char_data) int {
	var (
		nr    int
		found int = 0
	)
	for nr = 0; nr <= int(top_of_objt); nr++ {
		if isname(searchname, material_names[obj_proto[nr].Value[VAL_ALL_MATERIAL]]) != 0 {
			send_to_char(ch, libc.CString("%3d. [%5d] %-40s %s\r\n"), func() int {
				p := &found
				*p++
				return *p
			}(), obj_index[nr].Vnum, obj_proto[nr].Short_description, func() string {
				if obj_proto[nr].Proto_script != nil {
					return "[TRIG]"
				}
				return ""
			}())
		}
	}
	return found
}
func vnum_weapontype(searchname *byte, ch *char_data) int {
	var (
		nr    int
		found int = 0
	)
	for nr = 0; nr <= int(top_of_objt); nr++ {
		if int(obj_proto[nr].Type_flag) == ITEM_WEAPON {
			if isname(searchname, weapon_type[obj_proto[nr].Value[VAL_WEAPON_SKILL]]) != 0 {
				send_to_char(ch, libc.CString("%3d. [%5d] %-40s %s\r\n"), func() int {
					p := &found
					*p++
					return *p
				}(), obj_index[nr].Vnum, obj_proto[nr].Short_description, func() string {
					if obj_proto[nr].Proto_script != nil {
						return "[TRIG]"
					}
					return ""
				}())
			}
		}
	}
	return found
}
func vnum_armortype(searchname *byte, ch *char_data) int {
	var (
		nr    int
		found int = 0
	)
	for nr = 0; nr <= int(top_of_objt); nr++ {
		if int(obj_proto[nr].Type_flag) == ITEM_ARMOR {
			if isname(searchname, armor_type[obj_proto[nr].Value[VAL_ARMOR_SKILL]]) != 0 {
				send_to_char(ch, libc.CString("%3d. [%5d] %-40s %s\r\n"), func() int {
					p := &found
					*p++
					return *p
				}(), obj_index[nr].Vnum, obj_proto[nr].Short_description, func() string {
					if obj_proto[nr].Proto_script != nil {
						return "[TRIG]"
					}
					return ""
				}())
			}
		}
	}
	return found
}
func create_char() *char_data {
	var ch *char_data
	ch = new(char_data)
	clear_char(ch)
	ch.Next = character_list
	character_list = ch
	ch.Next_affect = nil
	ch.Next_affectv = nil
	ch.Id = int32(func() int {
		p := &max_mob_id
		x := *p
		*p++
		return x
	}())
	add_to_lookup_table(int(ch.Id), unsafe.Pointer(ch))
	return ch
}
func read_mobile(nr mob_vnum, type_ int) *char_data {
	var (
		i   mob_rnum
		mob *char_data
	)
	if type_ == VIRTUAL {
		if (func() mob_rnum {
			i = real_mobile(nr)
			return i
		}()) == mob_rnum(-1) {
			basic_mud_log(libc.CString("WARNING: Mobile vnum %d does not exist in database."), nr)
			return nil
		}
	} else {
		i = mob_rnum(nr)
	}
	mob = new(char_data)
	clear_char(mob)
	*mob = mob_proto[i]
	mob.Next = character_list
	character_list = mob
	mob.Next_affect = nil
	mob.Next_affectv = nil
	if int(mob.Race) == RACE_HOSHIJIN && int(mob.Sex) == SEX_MALE {
		mob.Hairl = 0
		mob.Hairc = 0
		mob.Hairs = 0
	} else {
		mob.Hairl = int8(rand_number(0, 4))
		mob.Hairc = int8(rand_number(1, 13))
		mob.Hairs = int8(rand_number(1, 11))
	}
	mob.Eye = int8(rand_number(0, 11))
	mob.Absorbs = 0
	mob.Absorbing = nil
	mob.Absorbby = nil
	mob.Sits = nil
	mob.Blocked = nil
	mob.Blocks = nil
	if int(mob.Race) != RACE_HUMAN && int(mob.Race) != RACE_SAIYAN && int(mob.Race) != RACE_HALFBREED && int(mob.Race) != RACE_NAMEK {
		mob.Skin = int8(rand_number(0, 11))
	}
	if int(mob.Race) == RACE_NAMEK {
		mob.Skin = 2
	}
	if int(mob.Race) == RACE_HUMAN || int(mob.Race) == RACE_SAIYAN || int(mob.Race) == RACE_HALFBREED {
		if rand_number(1, 5) <= 2 {
			mob.Skin = int8(rand_number(0, 1))
		} else if rand_number(1, 5) <= 4 {
			mob.Skin = int8(rand_number(4, 5))
		} else if rand_number(1, 5) <= 5 {
			mob.Skin = int8(rand_number(9, 10))
		}
	}
	if int(mob.Race) == RACE_SAIYAN {
		mob.Hairc = int8(rand_number(1, 2))
		mob.Eye = 1
	}
	if GET_MOB_VNUM(mob) >= 81 && GET_MOB_VNUM(mob) <= 87 {
		dragon_level(mob)
	}
	var mult int64 = 0
	switch GET_LEVEL(mob) {
	case 1:
		mult = int64(rand_number(50, 80))
	case 2:
		mult = int64(rand_number(90, 120))
	case 3:
		mult = int64(rand_number(100, 140))
	case 4:
		mult = int64(rand_number(120, 180))
	case 5:
		mult = int64(rand_number(200, 250))
	case 6:
		mult = int64(rand_number(240, 300))
	case 7:
		mult = int64(rand_number(280, 350))
	case 8:
		mult = int64(rand_number(320, 400))
	case 9:
		mult = int64(rand_number(380, 480))
	case 10:
		mult = int64(rand_number(500, 600))
	case 11:
		fallthrough
	case 12:
		fallthrough
	case 13:
		fallthrough
	case 14:
		fallthrough
	case 15:
		mult = int64(rand_number(1200, 1600))
	case 16:
		fallthrough
	case 17:
		fallthrough
	case 18:
		fallthrough
	case 19:
		fallthrough
	case 20:
		mult = int64(rand_number(2400, 3000))
	case 21:
		fallthrough
	case 22:
		fallthrough
	case 23:
		fallthrough
	case 24:
		fallthrough
	case 25:
		mult = int64(rand_number(5500, 8000))
	case 26:
		fallthrough
	case 27:
		fallthrough
	case 28:
		fallthrough
	case 29:
		fallthrough
	case 30:
		mult = int64(rand_number(10000, 14000))
	case 31:
		fallthrough
	case 32:
		fallthrough
	case 33:
		fallthrough
	case 34:
		fallthrough
	case 35:
		mult = int64(rand_number(16000, 20000))
	case 36:
		fallthrough
	case 37:
		fallthrough
	case 38:
		fallthrough
	case 39:
		fallthrough
	case 40:
		mult = int64(rand_number(22000, 30000))
	case 41:
		fallthrough
	case 42:
		fallthrough
	case 43:
		fallthrough
	case 44:
		fallthrough
	case 45:
		mult = int64(rand_number(50000, 70000))
	case 46:
		fallthrough
	case 47:
		fallthrough
	case 48:
		fallthrough
	case 49:
		fallthrough
	case 50:
		mult = int64(rand_number(95000, 140000))
	case 51:
		fallthrough
	case 52:
		fallthrough
	case 53:
		fallthrough
	case 54:
		fallthrough
	case 55:
		mult = int64(rand_number(180000, 250000))
	case 56:
		fallthrough
	case 57:
		fallthrough
	case 58:
		fallthrough
	case 59:
		fallthrough
	case 60:
		mult = int64(rand_number(400000, 480000))
	case 61:
		fallthrough
	case 62:
		fallthrough
	case 63:
		fallthrough
	case 64:
		fallthrough
	case 65:
		mult = int64(rand_number(700000, 900000))
	case 66:
		fallthrough
	case 67:
		fallthrough
	case 68:
		fallthrough
	case 69:
		fallthrough
	case 70:
		mult = int64(rand_number(1400000, 1600000))
	case 71:
		fallthrough
	case 72:
		fallthrough
	case 73:
		fallthrough
	case 74:
		fallthrough
	case 75:
		mult = int64(rand_number(2200000, 2500000))
	case 76:
		fallthrough
	case 77:
		fallthrough
	case 78:
		fallthrough
	case 79:
		fallthrough
	case 80:
		mult = int64(rand_number(3000000, 3500000))
	case 81:
		fallthrough
	case 82:
		fallthrough
	case 83:
		fallthrough
	case 84:
		fallthrough
	case 85:
		mult = int64(rand_number(4250000, 4750000))
	case 86:
		fallthrough
	case 87:
		fallthrough
	case 88:
		fallthrough
	case 89:
		fallthrough
	case 90:
		mult = int64(rand_number(6500000, 8500000))
	case 91:
		fallthrough
	case 92:
		fallthrough
	case 93:
		fallthrough
	case 94:
		fallthrough
	case 95:
		mult = int64(rand_number(15000000, 18000000))
	case 96:
		fallthrough
	case 97:
		fallthrough
	case 98:
		fallthrough
	case 99:
		fallthrough
	case 100:
		mult = int64(rand_number(22000000, 30000000))
	case 101:
		mult = int64(rand_number(32000000, 40000000))
	case 102:
		mult = int64(rand_number(42000000, 55000000))
	case 103:
		mult = int64(rand_number(80000000, 95000000))
	case 104:
		mult = int64(rand_number(150000000, 200000000))
	case 105:
		mult = int64(rand_number(220000000, 250000000))
	case 106:
		fallthrough
	case 107:
		fallthrough
	case 108:
		fallthrough
	case 109:
		fallthrough
	case 110:
		mult = int64(rand_number(500000000, 750000000))
	case 111:
		fallthrough
	case 112:
		fallthrough
	case 113:
		fallthrough
	case 114:
		fallthrough
	case 115:
		fallthrough
	case 116:
		fallthrough
	case 117:
		fallthrough
	case 118:
		fallthrough
	case 119:
		fallthrough
	case 120:
		mult = int64(rand_number(800000000, 900000000))
	default:
		if GET_LEVEL(mob) >= 150 {
			mult = int64(rand_number(1500000000, 2000000000))
		} else {
			mult = int64(rand_number(1250000000, 1500000000))
		}
	}
	mob.Lastpl = libc.GetTime(nil)
	if mob.Max_hit <= 1 {
		mob.Max_hit = int64(GET_LEVEL(mob) * int(mult))
		if GET_LEVEL(mob) > 140 {
			mob.Max_hit *= 8
		} else if GET_LEVEL(mob) > 130 {
			mob.Max_hit *= 6
		} else if GET_LEVEL(mob) > 120 {
			mob.Max_hit *= 3
		} else if GET_LEVEL(mob) > 110 {
			mob.Max_hit *= 2
		}
		mob.Hit = mob.Max_hit
		mob.Basepl = mob.Max_hit
	}
	if mob.Max_mana <= 1 {
		mob.Max_mana = int64(GET_LEVEL(mob) * int(mult))
		if GET_LEVEL(mob) > 140 {
			mob.Max_mana *= 8
		} else if GET_LEVEL(mob) > 130 {
			mob.Max_mana *= 6
		} else if GET_LEVEL(mob) > 120 {
			mob.Max_mana *= 3
		} else if GET_LEVEL(mob) > 110 {
			mob.Max_mana *= 2
		}
		mob.Mana = mob.Max_mana
		mob.Baseki = mob.Max_mana
	}
	if mob.Max_move <= 1 {
		mob.Max_move = int64(GET_LEVEL(mob) * int(mult))
		if GET_LEVEL(mob) > 140 {
			mob.Max_move *= 8
		} else if GET_LEVEL(mob) > 130 {
			mob.Max_move *= 6
		} else if GET_LEVEL(mob) > 120 {
			mob.Max_move *= 3
		} else if GET_LEVEL(mob) > 110 {
			mob.Max_move *= 2
		}
		mob.Move = mob.Max_move
		mob.Basest = mob.Max_move
	}
	if GET_MOB_VNUM(mob) == 2245 {
		mob.Max_hit = int64(rand_number(1, 4))
		mob.Hit = mob.Max_hit
		mob.Basepl = mob.Max_hit
		mob.Max_mana = int64(rand_number(1, 4))
		mob.Mana = mob.Max_mana
		mob.Baseki = mob.Max_mana
		mob.Max_move = int64(rand_number(1, 4))
		mob.Move = mob.Max_move
		mob.Basest = mob.Max_move
	}
	var base int = 0
	switch GET_LEVEL(mob) {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		base = rand_number(80, 120)
	case 6:
		base = rand_number(200, 280)
	case 7:
		base = rand_number(250, 350)
	case 8:
		base = rand_number(275, 375)
	case 9:
		base = rand_number(300, 400)
	case 10:
		base = rand_number(325, 450)
	case 11:
		fallthrough
	case 12:
		fallthrough
	case 13:
		fallthrough
	case 14:
		fallthrough
	case 15:
		base = rand_number(500, 700)
	case 16:
		fallthrough
	case 17:
		fallthrough
	case 18:
		fallthrough
	case 19:
		fallthrough
	case 20:
		base = rand_number(700, 1000)
	case 21:
		fallthrough
	case 22:
		fallthrough
	case 23:
		fallthrough
	case 24:
		fallthrough
	case 25:
		base = rand_number(1000, 1200)
	case 26:
		fallthrough
	case 27:
		fallthrough
	case 28:
		fallthrough
	case 29:
		fallthrough
	case 30:
		base = rand_number(1200, 1400)
	case 31:
		fallthrough
	case 32:
		fallthrough
	case 33:
		fallthrough
	case 34:
		fallthrough
	case 35:
		base = rand_number(1400, 1600)
	case 36:
		fallthrough
	case 37:
		fallthrough
	case 38:
		fallthrough
	case 39:
		fallthrough
	case 40:
		base = rand_number(1600, 1800)
	case 41:
		fallthrough
	case 42:
		fallthrough
	case 43:
		fallthrough
	case 44:
		fallthrough
	case 45:
		base = rand_number(1800, 2000)
	case 46:
		fallthrough
	case 47:
		fallthrough
	case 48:
		fallthrough
	case 49:
		fallthrough
	case 50:
		base = rand_number(2000, 2200)
	case 51:
		fallthrough
	case 52:
		fallthrough
	case 53:
		fallthrough
	case 54:
		fallthrough
	case 55:
		base = rand_number(2200, 2500)
	case 56:
		fallthrough
	case 57:
		fallthrough
	case 58:
		fallthrough
	case 59:
		fallthrough
	case 60:
		base = rand_number(2500, 2800)
	case 61:
		fallthrough
	case 62:
		fallthrough
	case 63:
		fallthrough
	case 64:
		fallthrough
	case 65:
		base = rand_number(2800, 3000)
	case 66:
		fallthrough
	case 67:
		fallthrough
	case 68:
		fallthrough
	case 69:
		fallthrough
	case 70:
		base = rand_number(3000, 3200)
	case 71:
		fallthrough
	case 72:
		fallthrough
	case 73:
		fallthrough
	case 74:
		fallthrough
	case 75:
		base = rand_number(3200, 3500)
	case 76:
		fallthrough
	case 77:
		fallthrough
	case 78:
		fallthrough
	case 79:
		base = rand_number(3500, 3800)
	case 80:
		fallthrough
	case 81:
		fallthrough
	case 82:
		fallthrough
	case 83:
		fallthrough
	case 84:
		fallthrough
	case 85:
		base = rand_number(4000, 4500)
	case 86:
		fallthrough
	case 87:
		fallthrough
	case 88:
		fallthrough
	case 89:
		fallthrough
	case 90:
		base = rand_number(4500, 5500)
	case 91:
		fallthrough
	case 92:
		fallthrough
	case 93:
		fallthrough
	case 94:
		fallthrough
	case 95:
		base = rand_number(5500, 7000)
	case 96:
		fallthrough
	case 97:
		fallthrough
	case 98:
		fallthrough
	case 99:
		base = rand_number(8000, 10000)
	case 100:
		base = rand_number(10000, 15000)
	case 101:
		base = rand_number(15000, 25000)
	case 102:
		base = rand_number(35000, 40000)
	case 103:
		base = rand_number(40000, 50000)
	case 104:
		base = rand_number(60000, 80000)
	case 105:
		base = rand_number(80000, 100000)
	default:
		base = rand_number(130000, 180000)
	}
	mob.Cooldown = 0
	if mob.Gold <= 0 && !MOB_FLAGGED(mob, MOB_DUMMY) {
		if GET_LEVEL(mob) < 4 {
			mob.Gold = GET_LEVEL(mob) * rand_number(1, 2)
		} else if GET_LEVEL(mob) < 10 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(1, 2)) - 1
		} else if GET_LEVEL(mob) < 20 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(1, 3)) - 2
		} else if GET_LEVEL(mob) < 30 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(1, 3)) - 4
		} else if GET_LEVEL(mob) < 40 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(1, 3)) - 6
		} else if GET_LEVEL(mob) < 50 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(2, 3)) - 25
		} else if GET_LEVEL(mob) < 60 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(2, 3)) - 40
		} else if GET_LEVEL(mob) < 70 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(2, 3)) - 50
		} else if GET_LEVEL(mob) < 80 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(2, 4)) - 60
		} else if GET_LEVEL(mob) < 90 {
			mob.Gold = (GET_LEVEL(mob) * rand_number(2, 4)) - 70
		} else {
			mob.Gold = (GET_LEVEL(mob) * rand_number(3, 4)) - 85
		}
		if !IS_HUMANOID(mob) {
			mob.Gold = int(float64(mob.Gold) * 0.5)
			if mob.Gold <= 0 {
				mob.Gold = 1
			}
		}
	}
	if mob.Exp <= 0 && !MOB_FLAGGED(mob, MOB_DUMMY) {
		mob.Exp = int64(GET_LEVEL(mob) * base)
		mob.Exp = int64(float64(mob.Exp) * 0.9)
		mob.Exp += int64(GET_LEVEL(mob) / 2)
		mob.Exp += int64(GET_LEVEL(mob) / 3)
		if int(mob.Race) == RACE_LIZARDFOLK {
			mob.Exp += int64(float64(mob.Exp) * .4)
		} else if int(mob.Race) == RACE_ANDROID {
			mob.Exp += int64(float64(mob.Exp) * .25)
		} else if int(mob.Race) == RACE_SAIYAN {
			mob.Exp += int64(float64(mob.Exp) * .1)
		} else if int(mob.Race) == RACE_BIO {
			mob.Exp += int64(float64(mob.Exp) * .2)
		} else if int(mob.Race) == RACE_MAJIN {
			mob.Exp += int64(float64(mob.Exp) * .25)
		} else if int(mob.Race) == RACE_DEMON {
			mob.Exp += int64(float64(mob.Exp) * .1)
		} else if int(mob.Chclass) == CLASS_SHADOWDANCER {
			mob.Exp *= 2
		}
		if int(mob.Chclass) == CLASS_NPC_COMMONER && IS_HUMANOID(mob) && int(mob.Race) != RACE_LIZARDFOLK {
			if int(mob.Race) != RACE_ANDROID && int(mob.Race) != RACE_SAIYAN && int(mob.Race) != RACE_BIO && int(mob.Race) != RACE_MAJIN {
				mob.Exp = int64(float64(mob.Exp) * 0.75)
			}
		}
		if GET_LEVEL(mob) > 90 {
			mob.Exp = int64(float64(mob.Exp) * 0.7)
		} else if GET_LEVEL(mob) > 80 {
			mob.Exp = int64(float64(mob.Exp) * 0.75)
		} else if GET_LEVEL(mob) > 70 {
			mob.Exp = int64(float64(mob.Exp) * 0.8)
		} else if GET_LEVEL(mob) > 60 {
			mob.Exp = int64(float64(mob.Exp) * 0.85)
		} else if GET_LEVEL(mob) > 40 {
			mob.Exp = int64(float64(mob.Exp) * 0.9)
		} else if GET_LEVEL(mob) > 30 {
			mob.Exp = int64(float64(mob.Exp) * 0.95)
		}
		if mob.Exp > 20000000 {
			mob.Exp = 20000000
		}
	}
	mob.Hit = mob.Max_hit
	mob.Mana = mob.Max_hit
	mob.Move = mob.Max_hit
	mob.Time.Birth = libc.GetTime(nil) - birth_age(mob)
	mob.Time.Created = func() libc.Time {
		p := &mob.Time.Logon
		mob.Time.Logon = libc.GetTime(nil)
		return *p
	}()
	mob.Time.Maxage = mob.Time.Birth + max_age(mob)
	mob.Time.Played = 0
	mob.Time.Logon = libc.GetTime(nil)
	mob.Hometown = -1
	if IS_HUMANOID(mob) {
		SET_BIT_AR(mob.Act[:], MOB_RARM)
		SET_BIT_AR(mob.Act[:], MOB_LARM)
		SET_BIT_AR(mob.Act[:], MOB_RLEG)
		SET_BIT_AR(mob.Act[:], MOB_LLEG)
	}
	mob_index[i].Number++
	mob.Id = int32(func() int {
		p := &max_mob_id
		x := *p
		*p++
		return x
	}())
	add_to_lookup_table(int(mob.Id), unsafe.Pointer(mob))
	copy_proto_script(unsafe.Pointer(&mob_proto[i]), unsafe.Pointer(mob), MOB_TRIGGER)
	assign_triggers(unsafe.Pointer(mob), MOB_TRIGGER)
	racial_body_parts(mob)
	if GET_MOB_VNUM(mob) >= 800 && GET_MOB_VNUM(mob) <= 805 {
		number_of_assassins += 1
	}
	return mob
}

type obj_unique_hash_elem struct {
	Generation libc.Time
	Unique_id  int64
	Obj        *obj_data
	Next_e     *obj_unique_hash_elem
}

func free_obj_unique_hash() {
	var (
		i         int
		elem      *obj_unique_hash_elem
		next_elem *obj_unique_hash_elem
	)
	if obj_unique_hash_pools != nil {
		for i = 0; i < NUM_OBJ_UNIQUE_POOLS; i++ {
			elem = *(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i)))
			for elem != nil {
				next_elem = elem.Next_e
				libc.Free(unsafe.Pointer(elem))
				elem = next_elem
			}
		}
		libc.Free(unsafe.Pointer(obj_unique_hash_pools))
	}
}
func add_unique_id(obj *obj_data) {
	var (
		elem *obj_unique_hash_elem
		i    int
	)
	if obj_unique_hash_pools == nil {
		init_obj_unique_hash()
	}
	if obj.Unique_id == -1 {
		if unsafe.Sizeof(int64(0)) > unsafe.Sizeof(int(0)) {
			obj.Unique_id = ((int64(circle_random())) << int64(unsafe.Sizeof(int64(0))*4)) + int64(circle_random())
		} else {
			obj.Unique_id = int64(circle_random())
		}
	}
	if config_info.Play.All_items_unique != 0 {
		if !OBJ_FLAGGED(obj, ITEM_UNIQUE_SAVE) {
			SET_BIT_AR(obj.Extra_flags[:], ITEM_UNIQUE_SAVE)
		}
	}
	elem = new(obj_unique_hash_elem)
	elem.Generation = obj.Generation
	elem.Unique_id = obj.Unique_id
	elem.Obj = obj
	i = int(obj.Unique_id % NUM_OBJ_UNIQUE_POOLS)
	elem.Next_e = *(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i)))
	*(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i))) = elem
}
func remove_unique_id(obj *obj_data) {
	var (
		elem *obj_unique_hash_elem
		ptr  **obj_unique_hash_elem
		tmp  *obj_unique_hash_elem
	)
	if obj == nil || obj.Unique_id < 0 {
		return
	}
	ptr = (**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(obj.Unique_id%NUM_OBJ_UNIQUE_POOLS)))
	if ptr == nil || *ptr == nil {
		return
	}
	elem = *ptr
	for elem != nil {
		tmp = elem.Next_e
		if elem.Obj == obj {
			libc.Free(unsafe.Pointer(elem))
			*ptr = tmp
		} else {
			ptr = &elem.Next_e
		}
		elem = tmp
	}
}
func log_dupe_objects(obj1 *obj_data, obj2 *obj_data) {
	mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("DUPE: Dupe object found: %s [%d] [%ld:%lld]"), func() *byte {
		if obj1.Short_description != nil {
			return obj1.Short_description
		}
		return libc.CString("<No name>")
	}(), GET_OBJ_VNUM(obj1), obj1.Generation, obj1.Unique_id)
	mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("DUPE: First: In room: %d (%s), In object: %s, Carried by: %s, Worn by: %s"), GET_ROOM_VNUM(obj1.In_room), func() string {
		if obj1.In_room == room_rnum(-1) {
			return "Nowhere"
		}
		return libc.GoString(world[obj1.In_room].Name)
	}(), func() *byte {
		if obj1.In_obj != nil {
			return obj1.In_obj.Short_description
		}
		return libc.CString("None")
	}(), func() *byte {
		if obj1.Carried_by != nil {
			return GET_NAME(obj1.Carried_by)
		}
		return libc.CString("Nobody")
	}(), func() *byte {
		if obj1.Worn_by != nil {
			return GET_NAME(obj1.Worn_by)
		}
		return libc.CString("Nobody")
	}())
	mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("DUPE: Newer: In room: %d (%s), In object: %s, Carried by: %s, Worn by: %s"), GET_ROOM_VNUM(obj2.In_room), func() string {
		if obj2.In_room == room_rnum(-1) {
			return "Nowhere"
		}
		return libc.GoString(world[obj2.In_room].Name)
	}(), func() *byte {
		if obj2.In_obj != nil {
			return obj2.In_obj.Short_description
		}
		return libc.CString("None")
	}(), func() *byte {
		if obj2.Carried_by != nil {
			return GET_NAME(obj2.Carried_by)
		}
		return libc.CString("Nobody")
	}(), func() *byte {
		if obj2.Worn_by != nil {
			return GET_NAME(obj2.Worn_by)
		}
		return libc.CString("Nobody")
	}())
}
func check_unique_id(obj *obj_data) {
	var elem *obj_unique_hash_elem
	if obj == nil || obj.Unique_id == -1 {
		return
	}
	elem = *(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(obj.Unique_id%NUM_OBJ_UNIQUE_POOLS)))
	for elem != nil {
		if elem.Obj == obj {
			basic_mud_log(libc.CString("SYSERR: check_unique_id checking for existing object?!"))
		}
		if elem.Generation == obj.Generation && elem.Unique_id == obj.Unique_id {
			log_dupe_objects(elem.Obj, obj)
			SET_BIT_AR(obj.Extra_flags[:], ITEM_PURGE)
		}
		elem = elem.Next_e
	}
}
func sprintuniques(low int, high int) *byte {
	var (
		i      int
		count  int = 0
		remain int
		header int
		q      *obj_unique_hash_elem
		str    *byte
		ptr    *byte
	)
	remain = 40
	for i = 0; i < NUM_OBJ_UNIQUE_POOLS; i++ {
		q = *(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i)))
		remain += 40
		for q != nil {
			count++
			remain += (func() int {
				if q.Obj.Short_description != nil {
					return libc.StrLen(q.Obj.Short_description)
				}
				return 20
			}()) + 80
			q = q.Next_e
		}
	}
	if count < 1 {
		return libc.CString("No objects in unique hash.\r\n")
	}
	str = (*byte)(unsafe.Pointer(&make([]int8, remain+1)[0]))
	ptr = str
	count = stdio.Snprintf(ptr, remain, "Unique object hashes (vnums %d - %d)\r\n", low, high)
	ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), count))
	remain -= count
	for i = 0; i < NUM_OBJ_UNIQUE_POOLS; i++ {
		header = 0
		q = *(**obj_unique_hash_elem)(unsafe.Add(unsafe.Pointer(obj_unique_hash_pools), unsafe.Sizeof((*obj_unique_hash_elem)(nil))*uintptr(i)))
		for q != nil {
			if GET_OBJ_VNUM(q.Obj) >= obj_vnum(low) && GET_OBJ_VNUM(q.Obj) <= obj_vnum(high) {
				if header == 0 {
					header = 1
					count = stdio.Snprintf(ptr, remain, "|-Hash %d\r\n", i)
					ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), count))
					remain -= count
				}
				count = stdio.Snprintf(ptr, remain, "| |- [@g%6d@n] - [@y%10ld:%-19lld@n] - %s\r\n", GET_OBJ_VNUM(q.Obj), q.Generation, q.Unique_id, func() *byte {
					if q.Obj.Short_description != nil {
						return q.Obj.Short_description
					}
					return libc.CString("<Unknown>")
				}())
				ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), count))
				remain -= count
			}
			q = q.Next_e
		}
	}
	return str
}
func create_obj() *obj_data {
	var obj *obj_data
	obj = new(obj_data)
	clear_object(obj)
	obj.Next = object_list
	object_list = obj
	obj.Id = int32(func() int {
		p := &max_obj_id
		x := *p
		*p++
		return x
	}())
	add_to_lookup_table(int(obj.Id), unsafe.Pointer(obj))
	obj.Generation = libc.GetTime(nil)
	obj.Unique_id = -1
	assign_triggers(unsafe.Pointer(obj), OBJ_TRIGGER)
	add_to_lookup_table(int(obj.Id), unsafe.Pointer(obj))
	return obj
}
func read_object(nr obj_vnum, type_ int) *obj_data {
	var (
		obj *obj_data
		i   obj_rnum
	)
	if type_ == VIRTUAL {
		i = real_object(nr)
	} else {
		i = obj_rnum(nr)
	}
	var j int
	if i == obj_rnum(-1) || i > top_of_objt {
		basic_mud_log(libc.CString("Object (%c) %d does not exist in database."), func() int {
			if type_ == VIRTUAL {
				return 'V'
			}
			return 'R'
		}(), nr)
		return nil
	}
	obj = new(obj_data)
	clear_object(obj)
	*obj = obj_proto[i]
	obj.Next = object_list
	object_list = obj
	obj.Room_loaded = -1
	obj_index[i].Number++
	obj.Id = int32(func() int {
		p := &max_obj_id
		x := *p
		*p++
		return x
	}())
	add_to_lookup_table(int(obj.Id), unsafe.Pointer(obj))
	obj.Generation = libc.GetTime(nil)
	obj.Unique_id = -1
	if obj_proto[i].Sbinfo != nil {
		obj.Sbinfo = make([]obj_spellbook_spell, SPELLBOOK_SIZE)
		for j = 0; j < SPELLBOOK_SIZE; j++ {
			obj.Sbinfo[j].Spellname = obj_proto[i].Sbinfo[j].Spellname
			obj.Sbinfo[j].Pages = obj_proto[i].Sbinfo[j].Pages
		}
	}
	copy_proto_script(unsafe.Pointer(&obj_proto[i]), unsafe.Pointer(obj), OBJ_TRIGGER)
	assign_triggers(unsafe.Pointer(obj), OBJ_TRIGGER)
	if GET_OBJ_VNUM(obj) == 65 {
		obj.Healcharge = 20
	}
	if int(obj.Type_flag) == ITEM_FOOD {
		if (obj.Value[1]) == 0 {
			obj.Value[1] = obj.Value[VAL_FOOD_FOODVAL]
		}
		obj.Foob = obj.Value[1]
	}
	return obj
}
func zone_update() {
	var (
		i        int
		update_u *reset_q_element
		temp     *reset_q_element
		timer    int = 0
	)
	if ((func() int {
		p := &timer
		*p++
		return *p
	}() * (config_info.Ticks.Pulse_zone * (int(1000000 / OPT_USEC)))) / (int(1000000 / OPT_USEC))) >= 60 {
		timer = 0
		for i = 0; i <= int(top_of_zone_table); i++ {
			if zone_table[i].Age < zone_table[i].Lifespan && zone_table[i].Reset_mode != 0 {
				zone_table[i].Age++
			}
			if zone_table[i].Age >= zone_table[i].Lifespan && zone_table[i].Age < ZO_DEAD && zone_table[i].Reset_mode != 0 {
				update_u = new(reset_q_element)
				update_u.Zone_to_reset = zone_rnum(i)
				update_u.Next = nil
				if reset_q.Head == nil {
					reset_q.Head = func() *reset_q_element {
						p := &reset_q.Tail
						reset_q.Tail = update_u
						return *p
					}()
				} else {
					reset_q.Tail.Next = update_u
					reset_q.Tail = update_u
				}
				zone_table[i].Age = ZO_DEAD
			}
		}
	}
	for update_u = reset_q.Head; update_u != nil; update_u = update_u.Next {
		if zone_table[update_u.Zone_to_reset].Reset_mode == 2 || is_empty(update_u.Zone_to_reset) != 0 {
			reset_zone(update_u.Zone_to_reset)
			mudlog(CMP, ADMLVL_GOD, FALSE, libc.CString("Auto zone reset: %s (Zone %d)"), zone_table[update_u.Zone_to_reset].Name, zone_table[update_u.Zone_to_reset].Number)
			if update_u == reset_q.Head {
				reset_q.Head = reset_q.Head.Next
			} else {
				for temp = reset_q.Head; temp.Next != update_u; temp = temp.Next {
				}
				if update_u.Next == nil {
					reset_q.Tail = temp
				}
				temp.Next = update_u.Next
			}
			libc.Free(unsafe.Pointer(update_u))
			break
		}
	}
}
func log_zone_error(zone zone_rnum, cmd_no int, message *byte) {
	mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: zone file: %s"), message)
	mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: ...offending cmd: '%c' cmd in zone #%d, line %d"), zone_table[zone].Cmd[cmd_no].Command, zone_table[zone].Number, zone_table[zone].Cmd[cmd_no].Line)
}
func reset_zone(zone zone_rnum) {
	var (
		cmd_no   int
		last_cmd int        = 0
		mob      *char_data = nil
		obj      *obj_data
		obj_to   *obj_data
		rvnum    room_vnum
		rrnum    room_rnum
		tmob     *char_data = nil
		tobj     *obj_data  = nil
		mob_load int        = FALSE
		obj_load int        = FALSE
	)
	if int(libc.BoolToInt(pre_reset(zone_table[zone].Number))) == FALSE {
		for cmd_no = 0; int(zone_table[zone].Cmd[cmd_no].Command) != 'S'; cmd_no++ {
			if zone_table[zone].Cmd[cmd_no].If_flag && last_cmd == 0 && mob_load == 0 && obj_load == 0 {
				continue
			}
			if !zone_table[zone].Cmd[cmd_no].If_flag {
				mob_load = FALSE
				obj_load = FALSE
			}
			switch zone_table[zone].Cmd[cmd_no].Command {
			case '*':
				last_cmd = 0
			case 'M':
				if mob_index[zone_table[zone].Cmd[cmd_no].Arg1].Number < int(zone_table[zone].Cmd[cmd_no].Arg2) && rand_number(1, 100) >= int(zone_table[zone].Cmd[cmd_no].Arg5) {
					var (
						room_max int = 0
						i        *char_data
					)
					mob = read_mobile(mob_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
					if zone_table[zone].Cmd[cmd_no].Arg4 > 0 {
						for i = character_list; i != nil; i = i.Next {
							if i.Hometown == room_vnum(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3)))) && GET_MOB_VNUM(i) == GET_MOB_VNUM(mob) {
								room_max++
							}
						}
					}
					char_to_room(mob, room_rnum(zone_table[zone].Cmd[cmd_no].Arg3))
					if room_max != 0 && room_max >= int(zone_table[zone].Cmd[cmd_no].Arg4) {
						extract_char(mob)
						extract_pending_chars()
						break
					}
					mob.Hometown = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3))))
					load_mtrigger(mob)
					tmob = mob
					last_cmd = 1
					mob_load = TRUE
				} else {
					last_cmd = 0
				}
				tobj = nil
			case 'O':
				if obj_index[zone_table[zone].Cmd[cmd_no].Arg1].Number < int(zone_table[zone].Cmd[cmd_no].Arg2) && rand_number(1, 100) >= int(zone_table[zone].Cmd[cmd_no].Arg5) {
					if zone_table[zone].Cmd[cmd_no].Arg3 != vnum(-1) {
						var (
							room_max int = 0
							k        *obj_data
						)
						obj = read_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
						if zone_table[zone].Cmd[cmd_no].Arg4 > 0 {
							for k = object_list; k != nil; k = k.Next {
								if k.Room_loaded == room_vnum(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3)))) && GET_OBJ_VNUM(k) == GET_OBJ_VNUM(obj) || GET_OBJ_VNUM(k) == GET_OBJ_VNUM(obj) && GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3)) == GET_ROOM_VNUM(k.In_room) {
									if k.In_room == room_rnum(-1) || GET_ROOM_VNUM(k.In_room) != GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3)) {
										continue
									}
									room_max++
								}
							}
						}
						add_unique_id(obj)
						obj_to_room(obj, room_rnum(zone_table[zone].Cmd[cmd_no].Arg3))
						if room_max != 0 && room_max >= int(zone_table[zone].Cmd[cmd_no].Arg4) {
							extract_obj(obj)
							break
						}
						obj.Room_loaded = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(zone_table[zone].Cmd[cmd_no].Arg3))))
						last_cmd = 1
						load_otrigger(obj)
						tobj = obj
						obj_load = TRUE
					} else {
						obj = read_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
						add_unique_id(obj)
						obj.In_room = -1
						last_cmd = 1
						tobj = obj
						obj_load = TRUE
					}
				} else {
					last_cmd = 0
				}
				tmob = nil
			case 'P':
				if obj_index[zone_table[zone].Cmd[cmd_no].Arg1].Number < int(zone_table[zone].Cmd[cmd_no].Arg2) && obj_load != 0 && rand_number(1, 100) >= int(zone_table[zone].Cmd[cmd_no].Arg5) {
					obj = read_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
					if (func() *obj_data {
						obj_to = get_obj_num(obj_rnum(zone_table[zone].Cmd[cmd_no].Arg3))
						return obj_to
					}()) == nil {
						log_zone_error(zone, cmd_no, libc.CString("target obj not found, command disabled"))
						last_cmd = 0
						zone_table[zone].Cmd[cmd_no].Command = '*'
						break
					}
					add_unique_id(obj)
					obj_to_obj(obj, obj_to)
					last_cmd = 1
					load_otrigger(obj)
					tobj = obj
				} else {
					last_cmd = 0
				}
				tmob = nil
			case 'G':
				if mob == nil {
					log_zone_error(zone, cmd_no, libc.CString("attempt to give obj to non-existant mob, command disabled"))
					last_cmd = 0
					zone_table[zone].Cmd[cmd_no].Command = '*'
					break
				}
				if obj_index[zone_table[zone].Cmd[cmd_no].Arg1].Number < int(zone_table[zone].Cmd[cmd_no].Arg2) && mob_load != 0 && rand_number(1, 100) >= int(zone_table[zone].Cmd[cmd_no].Arg5) {
					obj = read_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
					add_unique_id(obj)
					obj_to_char(obj, mob)
					if libc.FuncAddr(GET_MOB_SPEC(mob)) != libc.FuncAddr(shop_keeper) {
						randomize_eq(obj)
					}
					last_cmd = 1
					load_otrigger(obj)
					tobj = obj
				} else {
					last_cmd = 0
				}
				tmob = nil
			case 'E':
				if mob == nil {
					log_zone_error(zone, cmd_no, libc.CString("trying to equip non-existant mob, command disabled"))
					last_cmd = 0
					zone_table[zone].Cmd[cmd_no].Command = '*'
					break
				}
				if obj_index[zone_table[zone].Cmd[cmd_no].Arg1].Number < int(zone_table[zone].Cmd[cmd_no].Arg2) && mob_load != 0 && rand_number(1, 100) >= int(zone_table[zone].Cmd[cmd_no].Arg5) {
					if zone_table[zone].Cmd[cmd_no].Arg3 < 0 || zone_table[zone].Cmd[cmd_no].Arg3 >= NUM_WEARS {
						log_zone_error(zone, cmd_no, libc.CString("invalid equipment pos number"))
						last_cmd = 0
					} else {
						obj = read_object(obj_vnum(zone_table[zone].Cmd[cmd_no].Arg1), REAL)
						add_unique_id(obj)
						obj.In_room = mob.In_room
						load_otrigger(obj)
						if wear_otrigger(obj, mob, int(zone_table[zone].Cmd[cmd_no].Arg3)) != 0 {
							obj.In_room = -1
							equip_char(mob, obj, int(zone_table[zone].Cmd[cmd_no].Arg3))
						} else {
							obj_to_char(obj, mob)
						}
						tobj = obj
						last_cmd = 1
					}
				} else {
					last_cmd = 0
				}
				tmob = nil
			case 'R':
				if (func() *obj_data {
					obj = get_obj_in_list_num(int(zone_table[zone].Cmd[cmd_no].Arg2), world[zone_table[zone].Cmd[cmd_no].Arg1].Contents)
					return obj
				}()) != nil {
					extract_obj(obj)
				}
				last_cmd = 1
				tmob = nil
				tobj = nil
			case 'D':
				if zone_table[zone].Cmd[cmd_no].Arg2 < 0 || zone_table[zone].Cmd[cmd_no].Arg2 >= NUM_OF_DIRS || world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2] == nil {
					log_zone_error(zone, cmd_no, libc.CString("door does not exist, command disabled"))
					last_cmd = 0
					zone_table[zone].Cmd[cmd_no].Command = '*'
				} else {
					switch zone_table[zone].Cmd[cmd_no].Arg3 {
					case 0:
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info &= ^(bitvector_t(1) << 2)
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info &= ^(bitvector_t(1) << 1)
					case 1:
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info |= 1 << 1
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info &= ^(bitvector_t(1) << 2)
					case 2:
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info |= 1 << 2
						world[zone_table[zone].Cmd[cmd_no].Arg1].Dir_option[zone_table[zone].Cmd[cmd_no].Arg2].Exit_info |= 1 << 1
					}
				}
				last_cmd = 1
				tmob = nil
				tobj = nil
			case 'T':
				if zone_table[zone].Cmd[cmd_no].Arg1 == MOB_TRIGGER && tmob != nil {
					if tmob.Script == nil {
						tmob.Script = new(script_data)
					}
					add_trigger(tmob.Script, read_trigger(int(zone_table[zone].Cmd[cmd_no].Arg2)), -1)
					last_cmd = 1
				} else if zone_table[zone].Cmd[cmd_no].Arg1 == OBJ_TRIGGER && tobj != nil {
					if tobj.Script == nil {
						tobj.Script = new(script_data)
					}
					add_trigger(tobj.Script, read_trigger(int(zone_table[zone].Cmd[cmd_no].Arg2)), -1)
					last_cmd = 1
				} else if zone_table[zone].Cmd[cmd_no].Arg1 == WLD_TRIGGER {
					if zone_table[zone].Cmd[cmd_no].Arg3 == vnum(-1) || zone_table[zone].Cmd[cmd_no].Arg3 > vnum(top_of_world) {
						log_zone_error(zone, cmd_no, libc.CString("Invalid room number in trigger assignment"))
						last_cmd = 0
					}
					if world[zone_table[zone].Cmd[cmd_no].Arg3].Script == nil {
						world[zone_table[zone].Cmd[cmd_no].Arg3].Script = new(script_data)
					}
					add_trigger(world[zone_table[zone].Cmd[cmd_no].Arg3].Script, read_trigger(int(zone_table[zone].Cmd[cmd_no].Arg2)), -1)
					last_cmd = 1
				}
			case 'V':
				if zone_table[zone].Cmd[cmd_no].Arg1 == MOB_TRIGGER && tmob != nil {
					if tmob.Script == nil {
						log_zone_error(zone, cmd_no, libc.CString("Attempt to give variable to scriptless mobile"))
						last_cmd = 0
					} else {
						add_var(&tmob.Script.Global_vars, zone_table[zone].Cmd[cmd_no].Sarg1, zone_table[zone].Cmd[cmd_no].Sarg2, int(zone_table[zone].Cmd[cmd_no].Arg3))
					}
					last_cmd = 1
				} else if zone_table[zone].Cmd[cmd_no].Arg1 == OBJ_TRIGGER && tobj != nil {
					if tobj.Script == nil {
						log_zone_error(zone, cmd_no, libc.CString("Attempt to give variable to scriptless object"))
						last_cmd = 0
					} else {
						add_var(&tobj.Script.Global_vars, zone_table[zone].Cmd[cmd_no].Sarg1, zone_table[zone].Cmd[cmd_no].Sarg2, int(zone_table[zone].Cmd[cmd_no].Arg3))
					}
					last_cmd = 1
				} else if zone_table[zone].Cmd[cmd_no].Arg1 == WLD_TRIGGER {
					if zone_table[zone].Cmd[cmd_no].Arg3 == vnum(-1) || zone_table[zone].Cmd[cmd_no].Arg3 > vnum(top_of_world) {
						log_zone_error(zone, cmd_no, libc.CString("Invalid room number in variable assignment"))
						last_cmd = 0
					} else {
						if world[zone_table[zone].Cmd[cmd_no].Arg3].Script == nil {
							log_zone_error(zone, cmd_no, libc.CString("Attempt to give variable to scriptless object"))
							last_cmd = 0
						} else {
							add_var(&world[zone_table[zone].Cmd[cmd_no].Arg3].Script.Global_vars, zone_table[zone].Cmd[cmd_no].Sarg1, zone_table[zone].Cmd[cmd_no].Sarg2, int(zone_table[zone].Cmd[cmd_no].Arg2))
						}
						last_cmd = 1
					}
				}
			default:
				log_zone_error(zone, cmd_no, libc.CString("unknown cmd in reset table; cmd disabled"))
				last_cmd = 0
				zone_table[zone].Cmd[cmd_no].Command = '*'
			}
		}
		zone_table[zone].Age = 0
		rvnum = zone_table[zone].Bot
		for rvnum <= zone_table[zone].Top {
			rrnum = real_room(rvnum)
			if rrnum != room_rnum(-1) {
				reset_wtrigger(&world[rrnum])
				if ROOM_FLAGGED(rrnum, ROOM_AURA) && rand_number(1, 5) >= 4 {
					send_to_room(rrnum, libc.CString("The aura of regeneration covering the surrounding area disappears.\r\n"))
					REMOVE_BIT_AR(world[rrnum].Room_flags[:], ROOM_AURA)
				}
				if SECT(rrnum) == SECT_LAVA {
					world[rrnum].Geffect = 5
				}
				if world[rrnum].Geffect < -1 {
					send_to_room(rrnum, libc.CString("The area loses some of the water flooding it.\r\n"))
					world[rrnum].Geffect += 1
				} else if world[rrnum].Geffect == -1 {
					send_to_room(rrnum, libc.CString("The area loses the last of the water flooding it in one large rush.\r\n"))
					world[rrnum].Geffect = 0
				}
				if world[rrnum].Dmg >= 100 {
					send_to_room(rrnum, libc.CString("The area gets rebuilt a little.\r\n"))
					world[rrnum].Dmg -= rand_number(5, 10)
				} else if world[rrnum].Dmg >= 50 {
					send_to_room(rrnum, libc.CString("The area gets rebuilt a little.\r\n"))
					world[rrnum].Dmg -= rand_number(1, 10)
				} else if world[rrnum].Dmg >= 10 {
					send_to_room(rrnum, libc.CString("The area gets rebuilt a little.\r\n"))
					world[rrnum].Dmg -= rand_number(1, 10)
				} else if world[rrnum].Dmg > 1 {
					send_to_room(rrnum, libc.CString("The area gets rebuilt a little.\r\n"))
					world[rrnum].Dmg -= rand_number(1, world[rrnum].Dmg)
				} else if world[rrnum].Dmg > 0 {
					send_to_room(rrnum, libc.CString("The area gets rebuilt a little.\r\n"))
					world[rrnum].Dmg--
				}
				if world[rrnum].Geffect >= 1 && rand_number(1, 4) == 4 && !SUNKEN(rrnum) && SECT(rrnum) != SECT_LAVA {
					send_to_room(rrnum, libc.CString("The lava has cooled and become solid rock.\r\n"))
					world[rrnum].Geffect = 0
				} else if world[rrnum].Geffect >= 1 && rand_number(1, 2) == 2 && SUNKEN(rrnum) && SECT(rrnum) != SECT_LAVA {
					send_to_room(rrnum, libc.CString("The water has cooled the lava and it has become solid rock.\r\n"))
					world[rrnum].Geffect = 0
				}
			}
			rvnum++
		}
	} else {
		zone_table[zone].Age = 0
	}
	post_reset(zone_table[zone].Number)
}
func is_empty(zone_nr zone_rnum) int {
	var i *descriptor_data
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		if i.Character.In_room == room_rnum(-1) {
			continue
		}
		if world[i.Character.In_room].Zone != zone_nr {
			continue
		}
		if IS_NPC(i.Character) {
			continue
		}
		if i.Character.Admlevel >= ADMLVL_IMMORT && PRF_FLAGGED(i.Character, PRF_NOHASSLE) {
			continue
		}
		return 0
	}
	return 1
}
func fread_string(fl *stdio.File, error *byte) *byte {
	var (
		buf        [64936]byte
		tmp        [520]byte
		point      *byte
		done       int = 0
		length     int = 0
		templength int
	)
	buf[0] = func() byte {
		p := &tmp[0]
		tmp[0] = '\x00'
		return *p
	}()
	for {
		if fl.GetS(&tmp[0], 512) == nil {
			basic_mud_log(libc.CString("SYSERR: fread_string: format error at string (pos %ld): %s at or near %s"), fl.Tell(), func() string {
				if int(fl.IsEOF()) != 0 {
					return "EOF"
				}

				return "unknown error"
			}(), error)
			os.Exit(1)
		}
		for point = &tmp[0]; *point != 0 && *point != '\r' && *point != '\n'; point = (*byte)(unsafe.Add(unsafe.Pointer(point), 1)) {
		}
		if uintptr(unsafe.Pointer(point)) > uintptr(unsafe.Pointer(&tmp[0])) && *(*byte)(unsafe.Add(unsafe.Pointer(point), -1)) == '~' {
			*(func() *byte {
				p := &point
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), -1))
				return *p
			}()) = '\x00'
			done = 1
		} else {
			*point = '\r'
			*(func() *byte {
				p := &point
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()) = '\n'
			*(func() *byte {
				p := &point
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()) = '\x00'
		}
		// TODO: figure this out
		//templength = int(uintptr(unsafe.Pointer(point - tmp)))
		if length+templength >= MAX_STRING_LENGTH {
			basic_mud_log(libc.CString("SYSERR: fread_string: string too large (db.c)"))
			basic_mud_log(libc.CString("%s"), error)
			os.Exit(1)
		} else {
			libc.StrCat(&buf[length], &tmp[0])
			length += templength
		}
		if done != 0 {
			break
		}
	}
	if libc.StrLen(&buf[0]) != 0 {
		return libc.StrDup(&buf[0])
	}
	return nil
}
func free_followers(k *follow_type) {
	if k == nil {
		return
	}
	if k.Next != nil {
		free_followers(k.Next)
	}
	k.Follower = nil
	libc.Free(unsafe.Pointer(k))
}
func free_char(ch *char_data) {
	var (
		i          int
		a          *alias_data
		data       *levelup_data
		next_data  *levelup_data
		learn      *level_learn_entry
		next_learn *level_learn_entry
	)
	if ch.Player_specials != nil && ch.Player_specials != &dummy_mob {
		for (func() *alias_data {
			a = ch.Player_specials.Aliases
			return a
		}()) != nil {
			ch.Player_specials.Aliases = ch.Player_specials.Aliases.Next
			free_alias(a)
		}
		if ch.Player_specials.Poofin != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Poofin))
		}
		if ch.Player_specials.Poofout != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Poofout))
		}
		if ch.Player_specials.Host != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Host))
		}
		for i = 0; i < NUM_COLOR; i++ {
			if ch.Player_specials.Color_choices[i] != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Color_choices[i]))
			}
		}
		if IS_NPC(ch) {
			basic_mud_log(libc.CString("SYSERR: Mob %s (#%d) had player_specials allocated!"), GET_NAME(ch), GET_MOB_VNUM(ch))
		}
	}
	if !IS_NPC(ch) || IS_NPC(ch) && ch.Nr == mob_rnum(-1) {
		if GET_NAME(ch) != nil {
			libc.Free(unsafe.Pointer(GET_NAME(ch)))
		}
		if ch.Voice != nil {
			libc.Free(unsafe.Pointer(ch.Voice))
		}
		if ch.Clan != nil {
			libc.Free(unsafe.Pointer(ch.Clan))
		}
		if ch.Title != nil {
			libc.Free(unsafe.Pointer(ch.Title))
		}
		if ch.Short_descr != nil {
			libc.Free(unsafe.Pointer(ch.Short_descr))
		}
		if ch.Long_descr != nil {
			libc.Free(unsafe.Pointer(ch.Long_descr))
		}
		if ch.Description != nil {
			libc.Free(unsafe.Pointer(ch.Description))
		}
		for i = 0; i < NUM_HIST; i++ {
			if (ch.Player_specials.Comm_hist[i]) != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Comm_hist[i]))
			}
		}
		if ch.Player_specials != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials))
		}
		free_proto_script(unsafe.Pointer(ch), MOB_TRIGGER)
	} else if (func() int {
		i = int(ch.Nr)
		return i
	}()) != int(-1) {
		if ch.Name != nil && ch.Name != mob_proto[i].Name {
			libc.Free(unsafe.Pointer(ch.Name))
		}
		if ch.Title != nil && ch.Title != mob_proto[i].Title {
			libc.Free(unsafe.Pointer(ch.Title))
		}
		if ch.Short_descr != nil && ch.Short_descr != mob_proto[i].Short_descr {
			libc.Free(unsafe.Pointer(ch.Short_descr))
		}
		if ch.Long_descr != nil && ch.Long_descr != mob_proto[i].Long_descr {
			libc.Free(unsafe.Pointer(ch.Long_descr))
		}
		if ch.Description != nil && ch.Description != mob_proto[i].Description {
			libc.Free(unsafe.Pointer(ch.Description))
		}
		if ch.Proto_script != nil && ch.Proto_script != mob_proto[i].Proto_script {
			free_proto_script(unsafe.Pointer(ch), MOB_TRIGGER)
		}
	}
	for ch.Affected != nil {
		affect_remove(ch, ch.Affected)
	}
	if ch.Script != nil {
		extract_script(unsafe.Pointer(ch), MOB_TRIGGER)
	}
	free_followers(ch.Followers)
	if ch.Desc != nil {
		ch.Desc.Character = nil
	}
	if ch.Level_info != nil {
		for data = ch.Level_info; data != nil; data = next_data {
			next_data = data.Next
			for learn = data.Skills; learn != nil; learn = next_learn {
				next_learn = learn.Next
				libc.Free(unsafe.Pointer(learn))
			}
			for learn = data.Feats; learn != nil; learn = next_learn {
				next_learn = learn.Next
				libc.Free(unsafe.Pointer(learn))
			}
			libc.Free(unsafe.Pointer(data))
		}
	}
	ch.Level_info = nil
	if int(ch.Id) != 0 {
		remove_from_lookup_table(int(ch.Id))
	}
	libc.Free(unsafe.Pointer(ch))
}
func free_obj(obj *obj_data) {
	remove_unique_id(obj)
	if obj.Item_number == obj_vnum(-1) {
		free_object_strings(obj)
		free_proto_script(unsafe.Pointer(obj), OBJ_TRIGGER)
	} else {
		free_object_strings_proto(obj)
		if obj.Proto_script != obj_proto[obj.Item_number].Proto_script {
			free_proto_script(unsafe.Pointer(obj), OBJ_TRIGGER)
		}
	}
	if obj.Auctname != nil {
		libc.Free(unsafe.Pointer(obj.Auctname))
	}
	if obj.Script != nil {
		extract_script(unsafe.Pointer(obj), OBJ_TRIGGER)
	}
	remove_from_lookup_table(int(obj.Id))
	if obj.Sbinfo != nil {
		libc.Free(unsafe.Pointer(&obj.Sbinfo[0]))
	}
	libc.Free(unsafe.Pointer(obj))
}
func file_to_string_alloc(name *byte, buf **byte) int {
	var (
		temppage int
		temp     [64936]byte
		in_use   *descriptor_data
	)
	for in_use = descriptor_list; in_use != nil; in_use = in_use.Next {
		if in_use.Showstr_vector != nil && *in_use.Showstr_vector == *buf {
			return -1
		}
	}
	if file_to_string(name, &temp[0]) < 0 {
		return -1
	}
	for in_use = descriptor_list; in_use != nil; in_use = in_use.Next {
		if in_use.Showstr_count == 0 || *in_use.Showstr_vector != *buf {
			continue
		}
		temppage = in_use.Showstr_page
		paginate_string(func() *byte {
			p := &in_use.Showstr_head
			in_use.Showstr_head = libc.StrDup(*in_use.Showstr_vector)
			return *p
		}(), in_use)
		in_use.Showstr_page = temppage
	}
	if *buf != nil {
		libc.Free(unsafe.Pointer(*buf))
	}
	*buf = libc.StrDup(&temp[0])
	return 0
}
func file_to_string(name *byte, buf *byte) int {
	var (
		fl   *stdio.File
		tmp  [259]byte
		len_ int
	)
	*buf = '\x00'
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(name), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: reading %s: %s"), name, libc.StrError(libc.Errno))
		return -1
	}
	for {
		if fl.GetS(&tmp[0], READ_SIZE) == nil {
			break
		}
		if (func() int {
			len_ = libc.StrLen(&tmp[0])
			return len_
		}()) > 0 {
			tmp[len_-1] = '\x00'
		}
		libc.StrCat(&tmp[0], libc.CString("\r\n"))
		if libc.StrLen(buf)+libc.StrLen(&tmp[0])+1 > MAX_STRING_LENGTH {
			basic_mud_log(libc.CString("SYSERR: %s: string too big (%d max)"), name, MAX_STRING_LENGTH)
			*buf = '\x00'
			fl.Close()
			return -1
		}
		libc.StrCat(buf, &tmp[0])
	}
	fl.Close()
	return 0
}
func reset_char(ch *char_data) {
	var i int
	for i = 0; i < NUM_WEARS; i++ {
		ch.Equipment[i] = nil
	}
	ch.Followers = nil
	ch.Master = nil
	ch.In_room = -1
	ch.Carrying = nil
	ch.Next = nil
	ch.Next_fighting = nil
	ch.Next_in_room = nil
	ch.Fighting = nil
	ch.Position = POS_STANDING
	ch.Mob_specials.Default_pos = POS_STANDING
	ch.Carry_weight = 0
	ch.Carry_items = 0
	ch.Time.Logon = libc.GetTime(nil)
	if ch.Hit <= 0 {
		ch.Hit = 1
	}
	if ch.Move <= 0 {
		ch.Move = 1
	}
	if ch.Mana <= 0 {
		ch.Mana = 1
	}
	if ch.Ki < 0 {
		ch.Ki = 0
	}
	ch.Player_specials.Last_tell = -1
}
func clear_char(ch *char_data) {
	*(*char_data)(unsafe.Pointer((*byte)(unsafe.Pointer(ch)))) = char_data{}
	ch.In_room = -1
	ch.Pfilepos = -1
	ch.Nr = -1
	ch.Was_in_room = -1
	ch.Position = POS_STANDING
	ch.Mob_specials.Default_pos = POS_STANDING
	ch.Size = -1
	ch.Armor = 0
}
func clear_object(obj *obj_data) {
	*(*obj_data)(unsafe.Pointer((*byte)(unsafe.Pointer(obj)))) = obj_data{}
	obj.Item_number = -1
	obj.In_room = -1
	obj.Worn_on = -1
}
func init_char(ch *char_data) {
	var i int
	if ch.Player_specials == nil {
		ch.Player_specials = new(player_special_data)
		*(*player_special_data)(unsafe.Pointer(ch.Player_specials)) = player_special_data{}
	}
	ch.Admlevel = ADMLVL_NONE
	ch.Crank = 0
	ch.Clan = libc.CString("None.")
	ch.Absorbs = 0
	ch.Absorbing = nil
	ch.Absorbby = nil
	ch.Sits = nil
	ch.Blocked = nil
	ch.Blocks = nil
	if top_of_p_table == 0 {
		admin_set(ch, ADMLVL_IMPL)
		ch.Chclasses[ch.Chclass] = GET_LEVEL(ch)
		ch.Max_hit = 1000
		ch.Max_mana = 1000
		ch.Max_move = 1000
		ch.Max_ki = 1000
		ch.Hit = ch.Max_hit
		ch.Mana = ch.Max_mana
		ch.Move = ch.Max_move
		ch.Ki = ch.Max_ki
	}
	set_title(ch, nil)
	ch.Short_descr = nil
	ch.Long_descr = nil
	ch.Description = nil
	ch.Time.Logon = func() libc.Time {
		p := &ch.Time.Created
		ch.Time.Created = libc.GetTime(nil)
		return *p
	}()
	ch.Time.Maxage = ch.Time.Birth + max_age(ch)
	ch.Time.Played = 0
	ch.Hometown = 1
	ch.Armor = 0
	set_height_and_weight_by_race(ch)
	if (func() int {
		i = get_ptable_by_name(GET_NAME(ch))
		return i
	}()) != -1 {
		player_table[i].Id = int(func() int32 {
			p := &ch.Idnum
			ch.Idnum = int32(func() int {
				p := &top_idnum
				*p++
				return *p
			}())
			return *p
		}())
	} else {
		basic_mud_log(libc.CString("SYSERR: init_char: Character '%s' not found in player table."), GET_NAME(ch))
	}
	for i = 1; i < SKILL_TABLE_SIZE; i++ {
		if ch.Admlevel < ADMLVL_IMPL {
			for {
				ch.Skills[i] = 0
				if true {
					break
				}
			}
		} else {
			for {
				ch.Skills[i] = 100
				if true {
					break
				}
			}
		}
		for {
			ch.Skillmods[i] = 0
			if true {
				break
			}
		}
	}
	for i = 0; i < AF_ARRAY_MAX; i++ {
		ch.Affected_by[i] = 0
	}
	for i = 0; i < 3; i++ {
		ch.Apply_saving_throw[i] = 0
	}
	for i = 0; i < 3; i++ {
		if ch.Admlevel == ADMLVL_IMPL {
			ch.Player_specials.Conditions[i] = -1
		} else {
			ch.Player_specials.Conditions[i] = 24
		}
	}
	ch.Player_specials.Load_room = -1
	ch.Player_specials.Speaking = SKILL_LANG_COMMON
	ch.Player_specials.Feat_points = 1
}
func real_room(vnum room_vnum) room_rnum {
	var (
		bot      room_rnum
		top      room_rnum
		mid      room_rnum
		i        room_rnum
		last_top room_rnum
	)
	i = room_rnum(htree_find(room_htree, int64(vnum)))
	if i != room_rnum(-1) && world[i].Number == vnum {
		return i
	} else {
		bot = 0
		top = top_of_world
		for {
			last_top = top
			mid = (bot + top) / 2
			if world[mid].Number == vnum {
				basic_mud_log(libc.CString("room_htree sync fix: %d: %d -> %d"), vnum, i, mid)
				htree_add(room_htree, int64(vnum), int64(mid))
				return mid
			}
			if bot >= top {
				return -1
			}
			if world[mid].Number > vnum {
				top = mid - 1
			} else {
				bot = mid + 1
			}
			if top > last_top {
				return -1
			}
		}
	}
}
func real_mobile(vnum mob_vnum) mob_rnum {
	var (
		bot      mob_rnum
		top      mob_rnum
		mid      mob_rnum
		i        mob_rnum
		last_top mob_rnum
	)
	i = mob_rnum(htree_find(mob_htree, int64(vnum)))
	if i != mob_rnum(-1) && mob_index[i].Vnum == vnum {
		return i
	} else {
		bot = 0
		top = top_of_mobt
		for {
			last_top = top
			mid = (bot + top) / 2
			if mob_index[mid].Vnum == vnum {
				basic_mud_log(libc.CString("mob_htree sync fix: %d: %d -> %d"), vnum, i, mid)
				htree_add(mob_htree, int64(vnum), int64(mid))
				return mid
			}
			if bot >= top {
				return -1
			}
			if mob_index[mid].Vnum > vnum {
				top = mid - 1
			} else {
				bot = mid + 1
			}
			if top > last_top {
				return -1
			}
		}
	}
}
func real_object(vnum obj_vnum) obj_rnum {
	var (
		bot      obj_rnum
		top      obj_rnum
		mid      obj_rnum
		i        obj_rnum
		last_top obj_rnum
	)
	i = obj_rnum(htree_find(obj_htree, int64(vnum)))
	if i != obj_rnum(-1) && obj_index[i].Vnum == mob_vnum(vnum) {
		return i
	} else {
		bot = 0
		top = top_of_objt
		for {
			last_top = top
			mid = (bot + top) / 2
			if obj_index[mid].Vnum == mob_vnum(vnum) {
				basic_mud_log(libc.CString("obj_htree sync fix: %d: %d -> %d"), vnum, i, mid)
				htree_add(obj_htree, int64(vnum), int64(mid))
				return mid
			}
			if bot >= top {
				return -1
			}
			if obj_index[mid].Vnum > mob_vnum(vnum) {
				top = mid - 1
			} else {
				bot = mid + 1
			}
			if top > last_top {
				return -1
			}
		}
	}
}
func real_zone(vnum zone_vnum) zone_rnum {
	var (
		bot      zone_rnum
		top      zone_rnum
		mid      zone_rnum
		last_top zone_rnum
	)
	bot = 0
	top = top_of_zone_table
	for {
		last_top = top
		mid = (bot + top) / 2
		if zone_table[mid].Number == vnum {
			return mid
		}
		if bot >= top {
			return -1
		}
		if zone_table[mid].Number > vnum {
			top = mid - 1
		} else {
			bot = mid + 1
		}
		if top > last_top {
			return -1
		}
	}
}
func check_object(obj *obj_data) int {
	var (
		objname [2080]byte
		error   int = FALSE
		y       int
	)
	if obj.Weight < 0 && (func() int {
		error = TRUE
		return error
	}()) != 0 {
		basic_mud_log(libc.CString("SYSERR: Object #%d (%s) has negative weight (%lld)."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Weight)
	}
	if obj.Cost_per_day < 0 && (func() int {
		error = TRUE
		return error
	}()) != 0 {
		basic_mud_log(libc.CString("SYSERR: Object #%d (%s) has negative cost/day (%d)."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Cost_per_day)
	}
	stdio.Snprintf(&objname[0], int(2080), "Object #%d (%s)", GET_OBJ_VNUM(obj), obj.Short_description)
	for y = 0; y < TW_ARRAY_MAX; y++ {
		error |= check_bitvector_names(obj.Wear_flags[y], wear_bits_count, &objname[0], libc.CString("object wear"))
		error |= check_bitvector_names(obj.Extra_flags[y], extra_bits_count, &objname[0], libc.CString("object extra"))
		error |= check_bitvector_names(obj.Bitvector[y], affected_bits_count, &objname[0], libc.CString("object affect"))
	}
	switch obj.Type_flag {
	case ITEM_DRINKCON:
		var (
			onealias [2048]byte
			space    *byte = libc.StrRChr(obj.Name, ' ')
		)
		strlcpy(&onealias[0], func() *byte {
			if space != nil {
				return (*byte)(unsafe.Add(unsafe.Pointer(space), 1))
			}
			return obj.Name
		}(), uint64(2048))
		if search_block(&onealias[0], &drinknames[0], TRUE) < 0 && (func() int {
			error = TRUE
			return error
		}()) != 0 {
		}
		fallthrough
	case ITEM_FOUNTAIN:
		if (obj.Value[0]) > 0 && ((obj.Value[1]) > (obj.Value[0]) && (func() int {
			error = TRUE
			return error
		}()) != 0) {
			basic_mud_log(libc.CString("SYSERR: Object #%d (%s) contains (%d) more than maximum (%d)."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Value[1], obj.Value[0])
		}
	case ITEM_SCROLL:
		fallthrough
	case ITEM_POTION:
		error |= check_object_level(obj, 0)
		error |= check_object_spell_number(obj, 1)
		error |= check_object_spell_number(obj, 2)
		error |= check_object_spell_number(obj, 3)
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		error |= check_object_level(obj, 0)
		error |= check_object_spell_number(obj, 3)
		if (obj.Value[2]) > (obj.Value[1]) && (func() int {
			error = TRUE
			return error
		}()) != 0 {
			basic_mud_log(libc.CString("SYSERR: Object #%d (%s) has more charges (%d) than maximum (%d)."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Value[2], obj.Value[1])
		}
	}
	return error
}
func check_object_spell_number(obj *obj_data, val int) int {
	var (
		error     int = FALSE
		spellname *byte
	)
	if (obj.Value[val]) == -1 || (obj.Value[val]) == 0 {
		return error
	}
	if (obj.Value[val]) < 0 {
		error = TRUE
	}
	if (obj.Value[val]) >= SKILL_TABLE_SIZE {
		error = TRUE
	}
	if skill_type(obj.Value[val]) != (1 << 0) {
		error = TRUE
	}
	if error != 0 {
		basic_mud_log(libc.CString("SYSERR: Object #%d (%s) has out of range spell #%d."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Value[val])
	}
	if scheck != 0 {
		return error
	}
	spellname = skill_name(obj.Value[val])
	if (spellname == unused_spellname || libc.StrCaseCmp(libc.CString("UNDEFINED"), spellname) == 0) && (func() int {
		error = TRUE
		return error
	}()) != 0 {
		basic_mud_log(libc.CString("SYSERR: Object #%d (%s) uses '%s' spell #%d."), GET_OBJ_VNUM(obj), obj.Short_description, spellname, obj.Value[val])
	}
	return error
}
func check_object_level(obj *obj_data, val int) int {
	var error int = FALSE
	if (obj.Value[val]) < 0 && (func() int {
		error = TRUE
		return error
	}()) != 0 {
		basic_mud_log(libc.CString("SYSERR: Object #%d (%s) has out of range level #%d."), GET_OBJ_VNUM(obj), obj.Short_description, obj.Value[val])
	}
	return error
}
func check_bitvector_names(bits bitvector_t, namecount uint64, whatami *byte, whatbits *byte) int {
	var (
		flagnum uint
		error   bool = FALSE != 0
	)
	if uintptr(bits) <= (uintptr(^bitvector_t(0)) >> (unsafe.Sizeof(bitvector_t(0))*8 - uintptr(namecount))) {
		return FALSE
	}
	for flagnum = uint(namecount); flagnum < uint(unsafe.Sizeof(bitvector_t(0))*8); flagnum++ {
		if (1<<flagnum)&uint(bits) != 0 {
			basic_mud_log(libc.CString("SYSERR: %s has unknown %s flag, bit %d (0 through %lld known)."), whatami, whatbits, flagnum, namecount-1)
			error = TRUE != 0
		}
	}
	return int(libc.BoolToInt(error))
}
func my_obj_save_to_disk(fp *stdio.File, obj *obj_data, locate int) int {
	var (
		counter2 int
		i        int
		ex_desc  *extra_descr_data
		buf1     [64937]byte
		ebuf0    [64936]byte
		ebuf1    [64936]byte
		ebuf2    [64936]byte
		ebuf3    [64936]byte
	)
	if obj.Action_description != nil {
		libc.StrCpy(&buf1[0], obj.Action_description)
		strip_string(&buf1[0])
	} else {
		buf1[0] = 0
	}
	sprintascii(&ebuf0[0], obj.Extra_flags[0])
	sprintascii(&ebuf1[0], obj.Extra_flags[1])
	sprintascii(&ebuf2[0], obj.Extra_flags[2])
	sprintascii(&ebuf3[0], obj.Extra_flags[3])
	stdio.Fprintf(fp, "#%d\n%d %d %d %d %d %d %d %d %d %s %s %s %s %d %d %d %d %d %d %d %d\n", GET_OBJ_VNUM(obj), locate, obj.Value[0], obj.Value[1], obj.Value[2], obj.Value[3], obj.Value[4], obj.Value[5], obj.Value[6], obj.Value[7], &ebuf0[0], &ebuf1[0], &ebuf2[0], &ebuf3[0], obj.Value[8], obj.Value[9], obj.Value[10], obj.Value[11], obj.Value[12], obj.Value[13], obj.Value[14], obj.Value[15])
	if !OBJ_FLAGGED(obj, ITEM_UNIQUE_SAVE) && int(libc.BoolToInt(int(obj.Type_flag) == 0)) == ITEM_SPELLBOOK {
		return 1
	}
	stdio.Fprintf(fp, "XAP\n%s~\n%s~\n%s~\n%s~\n%d %d %d %d %d %lld %d %d\n", func() *byte {
		if obj.Name != nil {
			return obj.Name
		}
		return libc.CString("undefined")
	}(), func() *byte {
		if obj.Short_description != nil {
			return obj.Short_description
		}
		return libc.CString("undefined")
	}(), func() *byte {
		if obj.Description != nil {
			return obj.Description
		}
		return libc.CString("undefined")
	}(), &buf1[0], obj.Type_flag, obj.Wear_flags[0], obj.Wear_flags[1], obj.Wear_flags[2], obj.Wear_flags[3], obj.Weight, obj.Cost, obj.Cost_per_day)
	if obj.Generation != 0 {
		stdio.Fprintf(fp, "G\n%ld\n", obj.Generation)
	}
	if obj.Unique_id != 0 {
		stdio.Fprintf(fp, "U\n%lld\n", obj.Unique_id)
	}
	stdio.Fprintf(fp, "Z\n%d\n", obj.Size)
	for counter2 = 0; counter2 < MAX_OBJ_AFFECT; counter2++ {
		if obj.Affected[counter2].Modifier != 0 {
			stdio.Fprintf(fp, "A\n%d %d %d\n", obj.Affected[counter2].Location, obj.Affected[counter2].Modifier, obj.Affected[counter2].Specific)
		}
	}
	if obj.Ex_description != nil {
		for ex_desc = obj.Ex_description; ex_desc != nil; ex_desc = ex_desc.Next {
			if *ex_desc.Keyword == 0 || *ex_desc.Description == 0 {
				continue
			}
			libc.StrCpy(&buf1[0], ex_desc.Description)
			strip_string(&buf1[0])
			stdio.Fprintf(fp, "E\n%s~\n%s~\n", ex_desc.Keyword, &buf1[0])
		}
	}
	if obj.Sbinfo != nil {
		for i = 0; i < SPELLBOOK_SIZE; i++ {
			if obj.Sbinfo[i].Spellname == 0 {
				break
			}
			stdio.Fprintf(fp, "S\n%d %d\n", obj.Sbinfo[i].Spellname, obj.Sbinfo[i].Pages)
			continue
		}
	}
	return 1
}
func strip_string(buffer *byte) {
	var (
		ptr *byte
		str *byte
	)
	ptr = buffer
	str = ptr
	for (func() byte {
		p := str
		*str = *ptr
		return *p
	}()) != 0 {
		str = (*byte)(unsafe.Add(unsafe.Pointer(str), 1))
		ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		if *ptr == '\r' {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
	}
}
func load_default_config() {
	config_info.Play.Pk_allowed = pk_allowed
	config_info.Play.Pt_allowed = pt_allowed
	config_info.Play.Level_can_shout = level_can_shout
	config_info.Play.Holler_move_cost = holler_move_cost
	config_info.Play.Tunnel_size = tunnel_size
	config_info.Play.Max_exp_gain = max_exp_gain
	config_info.Play.Max_exp_loss = max_exp_loss
	config_info.Play.Max_npc_corpse_time = max_npc_corpse_time
	config_info.Play.Max_pc_corpse_time = max_pc_corpse_time
	config_info.Play.Idle_void = idle_void
	config_info.Play.Idle_rent_time = idle_rent_time
	config_info.Play.Idle_max_level = idle_max_level
	config_info.Play.Dts_are_dumps = dts_are_dumps
	config_info.Play.Load_into_inventory = load_into_inventory
	config_info.Play.OK = libc.StrDup(OK)
	config_info.Play.NOPERSON = libc.StrDup(NOPERSON)
	config_info.Play.NOEFFECT = libc.StrDup(NOEFFECT)
	config_info.Play.Track_through_doors = track_through_doors
	config_info.Play.Level_cap = level_cap
	config_info.Play.Stack_mobs = show_mob_stacking
	config_info.Play.Stack_objs = show_obj_stacking
	config_info.Play.Mob_fighting = mob_fighting
	config_info.Play.Disp_closed_doors = disp_closed_doors
	config_info.Play.Reroll_player = reroll_status
	config_info.Play.Initial_points = initial_points
	config_info.Play.Enable_compression = enable_compression
	config_info.Play.Enable_languages = enable_languages
	config_info.Play.All_items_unique = all_items_unique
	config_info.Play.Exp_multiplier = exp_multiplier
	config_info.Csd.Free_rent = free_rent
	config_info.Csd.Max_obj_save = max_obj_save
	config_info.Csd.Min_rent_cost = min_rent_cost
	config_info.Csd.Auto_save = auto_save
	config_info.Csd.Autosave_time = autosave_time
	config_info.Csd.Crash_file_timeout = crash_file_timeout
	config_info.Csd.Rent_file_timeout = rent_file_timeout
	config_info.Room_nums.Mortal_start_room = mortal_start_room
	config_info.Room_nums.Immort_start_room = immort_start_room
	config_info.Room_nums.Frozen_start_room = frozen_start_room
	config_info.Room_nums.Donation_room_1 = donation_room_1
	config_info.Room_nums.Donation_room_2 = donation_room_2
	config_info.Room_nums.Donation_room_3 = donation_room_3
	config_info.Operation.DFLT_PORT = DFLT_PORT
	if DFLT_IP != nil {
		config_info.Operation.DFLT_IP = libc.StrDup(DFLT_IP)
	} else {
		config_info.Operation.DFLT_IP = nil
	}
	config_info.Operation.DFLT_DIR = libc.StrDup(DFLT_DIR)
	if LOGNAME != nil {
		config_info.Operation.LOGNAME = libc.StrDup(LOGNAME)
	} else {
		config_info.Operation.LOGNAME = nil
	}
	config_info.Operation.Max_playing = max_playing
	config_info.Operation.Max_filesize = max_filesize
	config_info.Operation.Max_bad_pws = max_bad_pws
	config_info.Operation.Siteok_everyone = siteok_everyone
	config_info.Operation.Nameserver_is_slow = nameserver_is_slow
	config_info.Operation.Use_new_socials = use_new_socials
	config_info.Operation.Auto_save_olc = auto_save_olc
	config_info.Operation.MENU = libc.StrDup(MENU)
	config_info.Operation.WELC_MESSG = libc.StrDup(WELC_MESSG)
	config_info.Operation.START_MESSG = libc.StrDup(START_MESSG)
	config_info.Operation.Imc_enabled = imc_is_enabled
	config_info.Play.Exp_multiplier = 1.0
	config_info.Autowiz.Use_autowiz = use_autowiz
	config_info.Autowiz.Min_wizlist_lev = min_wizlist_lev
	config_info.Advance.Allow_multiclass = allow_multiclass
	config_info.Advance.Allow_prestige = allow_prestige
	config_info.Ticks.Pulse_violence = pulse_violence
	config_info.Ticks.Pulse_mobile = pulse_mobile
	config_info.Ticks.Pulse_zone = pulse_zone
	config_info.Ticks.Pulse_current = pulse_current
	config_info.Ticks.Pulse_sanity = pulse_sanity
	config_info.Ticks.Pulse_idlepwd = pulse_idlepwd
	config_info.Ticks.Pulse_autosave = pulse_autosave
	config_info.Ticks.Pulse_usage = pulse_usage
	config_info.Ticks.Pulse_timesave = pulse_timesave
	config_info.Creation.Method = method
}
func load_config() {
	var (
		fl   *stdio.File
		line [64936]byte
		tag  [2048]byte
		num  int
		fum  float32
		buf  [2048]byte
	)
	load_default_config()

	stdio.Snprintf(&buf[0], int(2048), "%s/%s", DFLT_DIR, config_info.CONFFILE)
	if (func() *stdio.File {
		fl = stdio.FOpen(config_info.CONFFILE, "r")
		return fl
	}()) == nil && (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&buf[0]), "r")
		return fl
	}()) == nil {
		stdio.Snprintf(&buf[0], int(2048), "Game Config File: %s", config_info.CONFFILE)
		fmt.Println(&buf[0])
		return
	}

	for get_line(fl, &line[0]) != 0 {
		split_argument(&line[0], &tag[0])
		num = libc.Atoi(libc.GoString(&line[0]))
		fum = float32(libc.Atof(libc.GoString(&line[0])))
		switch unicode.ToLower(rune(tag[0])) {
		case 'a':
			if libc.StrCaseCmp(&tag[0], libc.CString("auto_save")) == 0 {
				config_info.Csd.Auto_save = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("autosave_time")) == 0 {
				config_info.Csd.Autosave_time = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("auto_save_olc")) == 0 {
				config_info.Operation.Auto_save_olc = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("allow_multiclass")) == 0 {
				config_info.Advance.Allow_multiclass = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("allow_prestige")) == 0 {
				config_info.Advance.Allow_prestige = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("auto_level")) == 0 {
				basic_mud_log(libc.CString("ignoring obsolete config option auto_level"))
			} else if libc.StrCaseCmp(&tag[0], libc.CString("all_items_unique")) == 0 {
				config_info.Play.All_items_unique = num
			}
		case 'c':
			if libc.StrCaseCmp(&tag[0], libc.CString("crash_file_timeout")) == 0 {
				config_info.Csd.Crash_file_timeout = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("compression")) == 0 {
				config_info.Play.Enable_compression = num
			}
		case 'd':
			if libc.StrCaseCmp(&tag[0], libc.CString("disp_closed_doors")) == 0 {
				config_info.Play.Disp_closed_doors = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("dts_are_dumps")) == 0 {
				config_info.Play.Dts_are_dumps = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("donation_room_1")) == 0 {
				if num == -1 {
					config_info.Room_nums.Donation_room_1 = -1
				} else {
					config_info.Room_nums.Donation_room_1 = room_vnum(num)
				}
			} else if libc.StrCaseCmp(&tag[0], libc.CString("donation_room_2")) == 0 {
				if num == -1 {
					config_info.Room_nums.Donation_room_2 = -1
				} else {
					config_info.Room_nums.Donation_room_2 = room_vnum(num)
				}
			} else if libc.StrCaseCmp(&tag[0], libc.CString("donation_room_3")) == 0 {
				if num == -1 {
					config_info.Room_nums.Donation_room_3 = -1
				} else {
					config_info.Room_nums.Donation_room_3 = room_vnum(num)
				}
			} else if libc.StrCaseCmp(&tag[0], libc.CString("dflt_dir")) == 0 {
				if config_info.Operation.DFLT_DIR != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.DFLT_DIR))
				}
				if line[0] != 0 {
					config_info.Operation.DFLT_DIR = libc.StrDup(&line[0])
				} else {
					config_info.Operation.DFLT_DIR = libc.StrDup(DFLT_DIR)
				}
			} else if libc.StrCaseCmp(&tag[0], libc.CString("dflt_ip")) == 0 {
				if config_info.Operation.DFLT_IP != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.DFLT_IP))
				}
				if line[0] != 0 {
					config_info.Operation.DFLT_IP = libc.StrDup(&line[0])
				} else {
					config_info.Operation.DFLT_IP = nil
				}
			} else if libc.StrCaseCmp(&tag[0], libc.CString("dflt_port")) == 0 {
				config_info.Operation.DFLT_PORT = uint16(int16(num))
			}
		case 'e':
			if libc.StrCaseCmp(&tag[0], libc.CString("enable_languages")) == 0 {
				config_info.Play.Enable_languages = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("exp_multiplier")) == 0 {
				config_info.Play.Exp_multiplier = fum
			}
		case 'f':
			if libc.StrCaseCmp(&tag[0], libc.CString("free_rent")) == 0 {
				config_info.Csd.Free_rent = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("frozen_start_room")) == 0 {
				config_info.Room_nums.Frozen_start_room = room_vnum(num)
			}
		case 'h':
			if libc.StrCaseCmp(&tag[0], libc.CString("holler_move_cost")) == 0 {
				config_info.Play.Holler_move_cost = num
			}
		case 'i':
			if libc.StrCaseCmp(&tag[0], libc.CString("idle_void")) == 0 {
				config_info.Play.Idle_void = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("idle_rent_time")) == 0 {
				config_info.Play.Idle_rent_time = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("idle_max_level")) == 0 {
				if num >= config_info.Play.Level_cap {
					num += 1 - config_info.Play.Level_cap
				}
				config_info.Play.Idle_max_level = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("immort_level_ok")) == 0 {
				basic_mud_log(libc.CString("Ignoring immort_level_ok obsolete config"))
			} else if libc.StrCaseCmp(&tag[0], libc.CString("immort_start_room")) == 0 {
				config_info.Room_nums.Immort_start_room = room_vnum(num)
			} else if libc.StrCaseCmp(&tag[0], libc.CString("imc_enabled")) == 0 {
				config_info.Operation.Imc_enabled = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("initial_points")) == 0 {
				config_info.Play.Initial_points = num
			}
		case 'l':
			if libc.StrCaseCmp(&tag[0], libc.CString("level_can_shout")) == 0 {
				config_info.Play.Level_can_shout = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("level_cap")) == 0 {
				config_info.Play.Level_cap = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("load_into_inventory")) == 0 {
				config_info.Play.Load_into_inventory = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("logname")) == 0 {
				if config_info.Operation.LOGNAME != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.LOGNAME))
				}
				if line[0] != 0 {
					config_info.Operation.LOGNAME = libc.StrDup(&line[0])
				} else {
					config_info.Operation.LOGNAME = nil
				}
			}
		case 'm':
			if libc.StrCaseCmp(&tag[0], libc.CString("max_bad_pws")) == 0 {
				config_info.Operation.Max_bad_pws = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_exp_gain")) == 0 {
				config_info.Play.Max_exp_gain = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_exp_loss")) == 0 {
				config_info.Play.Max_exp_loss = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_filesize")) == 0 {
				config_info.Operation.Max_filesize = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_npc_corpse_time")) == 0 {
				config_info.Play.Max_npc_corpse_time = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_obj_save")) == 0 {
				config_info.Csd.Max_obj_save = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_pc_corpse_time")) == 0 {
				config_info.Play.Max_pc_corpse_time = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("max_playing")) == 0 {
				config_info.Operation.Max_playing = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("menu")) == 0 {
				if config_info.Operation.MENU != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.MENU))
				}
				libc.StrNCpy(&buf[0], libc.CString("Reading menu in load_config()"), int(2048))
				config_info.Operation.MENU = fread_string(fl, &buf[0])
			} else if libc.StrCaseCmp(&tag[0], libc.CString("min_rent_cost")) == 0 {
				config_info.Csd.Min_rent_cost = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("min_wizlist_lev")) == 0 {
				if num >= config_info.Play.Level_cap {
					num += 1 - config_info.Play.Level_cap
				}
				config_info.Autowiz.Min_wizlist_lev = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("mob_fighting")) == 0 {
				config_info.Play.Mob_fighting = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("mortal_start_room")) == 0 {
				config_info.Room_nums.Mortal_start_room = room_vnum(num)
			} else if libc.StrCaseCmp(&tag[0], libc.CString("method")) == 0 {
				config_info.Creation.Method = num
			}
		case 'n':
			if libc.StrCaseCmp(&tag[0], libc.CString("nameserver_is_slow")) == 0 {
				config_info.Operation.Nameserver_is_slow = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("noperson")) == 0 {
				var tmp [256]byte
				if config_info.Play.NOPERSON != nil {
					libc.Free(unsafe.Pointer(config_info.Play.NOPERSON))
				}
				stdio.Snprintf(&tmp[0], int(256), "%s\r\n", &line[0])
				config_info.Play.NOPERSON = libc.StrDup(&tmp[0])
			} else if libc.StrCaseCmp(&tag[0], libc.CString("noeffect")) == 0 {
				var tmp [256]byte
				if config_info.Play.NOEFFECT != nil {
					libc.Free(unsafe.Pointer(config_info.Play.NOEFFECT))
				}
				stdio.Snprintf(&tmp[0], int(256), "%s\r\n", &line[0])
				config_info.Play.NOEFFECT = libc.StrDup(&tmp[0])
			}
		case 'o':
			if libc.StrCaseCmp(&tag[0], libc.CString("ok")) == 0 {
				var tmp [256]byte
				if config_info.Play.OK != nil {
					libc.Free(unsafe.Pointer(config_info.Play.OK))
				}
				stdio.Snprintf(&tmp[0], int(256), "%s\r\n", &line[0])
				config_info.Play.OK = libc.StrDup(&tmp[0])
			}
		case 'p':
			if libc.StrCaseCmp(&tag[0], libc.CString("pk_allowed")) == 0 {
				config_info.Play.Pk_allowed = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pt_allowed")) == 0 {
				config_info.Play.Pt_allowed = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_viol")) == 0 {
				config_info.Ticks.Pulse_violence = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_mobile")) == 0 {
				config_info.Ticks.Pulse_mobile = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_current")) == 0 {
				config_info.Ticks.Pulse_current = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_zone")) == 0 {
				config_info.Ticks.Pulse_zone = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_autosave")) == 0 {
				config_info.Ticks.Pulse_autosave = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_usage")) == 0 {
				config_info.Ticks.Pulse_usage = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_sanity")) == 0 {
				config_info.Ticks.Pulse_sanity = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_timesave")) == 0 {
				config_info.Ticks.Pulse_timesave = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("pulse_idlepwd")) == 0 {
				config_info.Ticks.Pulse_idlepwd = num
			}
		case 'r':
			if libc.StrCaseCmp(&tag[0], libc.CString("rent_file_timeout")) == 0 {
				config_info.Csd.Rent_file_timeout = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("reroll_stats")) == 0 {
				config_info.Play.Reroll_player = num
			}
		case 's':
			if libc.StrCaseCmp(&tag[0], libc.CString("siteok_everyone")) == 0 {
				config_info.Operation.Siteok_everyone = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("start_messg")) == 0 {
				libc.StrNCpy(&buf[0], libc.CString("Reading start message in load_config()"), int(2048))
				if config_info.Operation.START_MESSG != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.START_MESSG))
				}
				config_info.Operation.START_MESSG = fread_string(fl, &buf[0])
			} else if libc.StrCaseCmp(&tag[0], libc.CString("stack_mobs")) == 0 {
				config_info.Play.Stack_mobs = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("stack_objs")) == 0 {
				config_info.Play.Stack_objs = num
			}
		case 't':
			if libc.StrCaseCmp(&tag[0], libc.CString("tunnel_size")) == 0 {
				config_info.Play.Tunnel_size = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("track_through_doors")) == 0 {
				config_info.Play.Track_through_doors = num
			}
		case 'u':
			if libc.StrCaseCmp(&tag[0], libc.CString("use_autowiz")) == 0 {
				config_info.Autowiz.Use_autowiz = num
			} else if libc.StrCaseCmp(&tag[0], libc.CString("use_new_socials")) == 0 {
				config_info.Operation.Use_new_socials = num
			}
		case 'w':
			if libc.StrCaseCmp(&tag[0], libc.CString("welc_messg")) == 0 {
				libc.StrNCpy(&buf[0], libc.CString("Reading welcome message in load_config()"), int(2048))
				if config_info.Operation.WELC_MESSG != nil {
					libc.Free(unsafe.Pointer(config_info.Operation.WELC_MESSG))
				}
				config_info.Operation.WELC_MESSG = fread_string(fl, &buf[0])
			}
		default:
		}
	}
	fl.Close()
}
func read_level_data(ch *char_data, fl *stdio.File) {
	var (
		buf   [256]byte
		p     *byte
		i     int = 1
		t     [16]int
		curr  *levelup_data = nil
		learn *level_learn_entry
	)
	ch.Level_info = nil
	for int(fl.IsEOF()) == 0 {
		i++
		if get_line(fl, &buf[0]) == 0 {
			basic_mud_log(libc.CString("read_level_data: get_line() failed reading level data line %d for %s"), i, GET_NAME(ch))
			return
		}
		for p = &buf[0]; *p != 0 && *p != ' '; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if libc.StrCmp(&buf[0], libc.CString("end")) == 0 {
			return
		}
		if *p == 0 {
			basic_mud_log(libc.CString("read_level_data: malformed line reading level data line %d for %s: %s"), i, GET_NAME(ch), &buf[0])
			return
		}
		*(func() *byte {
			p := &p
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = 0
		if libc.StrCmp(&buf[0], libc.CString("level")) == 0 {
			if stdio.Sscanf(p, "%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5], &t[6], &t[7], &t[8], &t[9], &t[10], &t[11], &t[12], &t[13], &t[14], &t[15]) != 16 {
				basic_mud_log(libc.CString("read_level_data: missing fields on level_data line %d for %s"), i, GET_NAME(ch))
				curr = nil
				continue
			}
			curr = new(levelup_data)
			curr.Prev = nil
			curr.Next = ch.Level_info
			if (func() *levelup_data {
				p := &curr.Next
				curr.Next = ch.Level_info
				return *p
			}()) != nil {
				curr.Next.Prev = curr
			}
			ch.Level_info = curr
			curr.Type = int8(t[0])
			curr.Spec = int8(t[1])
			curr.Level = int8(t[2])
			curr.Hp_roll = int8(t[3])
			curr.Mana_roll = int8(t[4])
			curr.Ki_roll = int8(t[5])
			curr.Move_roll = int8(t[6])
			curr.Fort = int8(t[8])
			curr.Reflex = int8(t[9])
			curr.Will = int8(t[10])
			curr.Add_skill = int8(t[11])
			curr.Add_gen_feats = int8(t[12])
			curr.Add_epic_feats = int8(t[13])
			curr.Add_class_feats = int8(t[14])
			curr.Add_class_epic_feats = int8(t[15])
			curr.Skills = func() *level_learn_entry {
				p := &curr.Feats
				curr.Feats = nil
				return *p
			}()
			continue
		}
		if curr == nil {
			basic_mud_log(libc.CString("read_level_data: found continuation entry without current level for %s"), GET_NAME(ch))
			continue
		}
		if stdio.Sscanf(p, "%d %d %d", &t[0], &t[1], &t[2]) != 3 {
			basic_mud_log(libc.CString("read_level_data: missing fields on level_data %s line %d for %s"), &buf[0], i, GET_NAME(ch))
			continue
		}
		learn = new(level_learn_entry)
		learn.Location = t[0]
		learn.Specific = t[1]
		learn.Value = int8(t[2])
		if libc.StrCmp(&buf[0], libc.CString("skill")) == 0 {
			learn.Next = curr.Skills
			curr.Skills = learn
		} else if libc.StrCmp(&buf[0], libc.CString("feat")) == 0 {
			learn.Next = curr.Feats
			curr.Feats = learn
		}
	}
	basic_mud_log(libc.CString("read_level_data: EOF reached reading level_data for %s"), GET_NAME(ch))
	return
}
func write_level_data(ch *char_data, fl *stdio.File) {
	var (
		lev   *levelup_data
		learn *level_learn_entry
	)
	for lev = ch.Level_info; lev != nil && lev.Next != nil; lev = lev.Next {
	}
	for lev != nil {
		stdio.Fprintf(fl, "level %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n", lev.Type, lev.Spec, lev.Level, lev.Hp_roll, lev.Mana_roll, lev.Ki_roll, lev.Move_roll, lev.Accuracy, lev.Fort, lev.Reflex, lev.Will, lev.Add_skill, lev.Add_gen_feats, lev.Add_epic_feats, lev.Add_class_feats, lev.Add_class_epic_feats)
		for learn = lev.Skills; learn != nil; learn = learn.Next {
			stdio.Fprintf(fl, "skill %d %d %d", learn.Location, learn.Specific, learn.Value)
		}
		for learn = lev.Feats; learn != nil; learn = learn.Next {
			stdio.Fprintf(fl, "feat %d %d %d", learn.Location, learn.Specific, learn.Value)
		}
		lev = lev.Prev
	}
	stdio.Fprintf(fl, "end\n")
}
