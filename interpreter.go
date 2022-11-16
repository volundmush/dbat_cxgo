package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const RESERVE_CMDS = 15
const ALIAS_SIMPLE = 0
const ALIAS_COMPLEX = 1
const ALIAS_SEP_CHAR = 59
const ALIAS_VAR_CHAR = 36
const ALIAS_GLOB_CHAR = 42
const SCMD_NORTH = 1
const SCMD_EAST = 2
const SCMD_SOUTH = 3
const SCMD_WEST = 4
const SCMD_UP = 5
const SCMD_DOWN = 6
const SCMD_NW = 7
const SCMD_NE = 8
const SCMD_SE = 9
const SCMD_SW = 10
const SCMD_IN = 11
const SCMD_OUT = 12
const SCMD_INFO = 0
const SCMD_HANDBOOK = 1
const SCMD_CREDITS = 2
const SCMD_NEWS = 3
const SCMD_WIZLIST = 4
const SCMD_POLICIES = 5
const SCMD_VERSION = 6
const SCMD_IMMLIST = 7
const SCMD_MOTD = 8
const SCMD_IMOTD = 9
const SCMD_CLEAR = 10
const SCMD_WHOAMI = 11
const SCMD_NOSUMMON = 0
const SCMD_NOHASSLE = 1
const SCMD_BRIEF = 2
const SCMD_COMPACT = 3
const SCMD_NOTELL = 4
const SCMD_NOAUCTION = 5
const SCMD_DEAF = 6
const SCMD_NOGOSSIP = 7
const SCMD_NOGRATZ = 8
const SCMD_NOWIZ = 9
const SCMD_QUEST = 10
const SCMD_ROOMFLAGS = 11
const SCMD_NOREPEAT = 12
const SCMD_HOLYLIGHT = 13
const SCMD_SLOWNS = 14
const SCMD_AUTOEXIT = 15
const SCMD_TRACK = 16
const SCMD_BUILDWALK = 17
const SCMD_AFK = 18
const SCMD_AUTOASSIST = 19
const SCMD_AUTOLOOT = 20
const SCMD_AUTOGOLD = 21
const SCMD_CLS = 22
const SCMD_AUTOSPLIT = 23
const SCMD_AUTOSAC = 24
const SCMD_SNEAK = 25
const SCMD_HIDE = 26
const SCMD_AUTOMEM = 27
const SCMD_VIEWORDER = 28
const SCMD_NOCOMPRESS = 29
const SCMD_TEST = 30
const SCMD_WHOHIDE = 31
const SCMD_NMWARN = 32
const SCMD_HINTS = 33
const SCMD_NODEC = 34
const SCMD_NOEQSEE = 35
const SCMD_NOMUSIC = 36
const SCMD_NOPARRY = 37
const SCMD_LKEEP = 38
const SCMD_CARVE = 39
const SCMD_NOGIVE = 40
const SCMD_INSTRUCT = 41
const SCMD_GHEALTH = 42
const SCMD_IHEALTH = 43
const SCMD_REROLL = 0
const SCMD_PARDON = 1
const SCMD_NOTITLE = 2
const SCMD_SQUELCH = 3
const SCMD_FREEZE = 4
const SCMD_THAW = 5
const SCMD_UNAFFECT = 6
const SCMD_WHISPER = 0
const SCMD_ASK = 1
const SCMD_HOLLER = 0
const SCMD_SHOUT = 1
const SCMD_GOSSIP = 2
const SCMD_AUCTION = 3
const SCMD_GRATZ = 4
const SCMD_GEMOTE = 5
const SCMD_SHUTDOW = 0
const SCMD_SHUTDOWN = 1
const SCMD_QUI = 0
const SCMD_QUIT = 1
const SCMD_DATE = 0
const SCMD_UPTIME = 1
const SCMD_COMMANDS = 0
const SCMD_SOCIALS = 1
const SCMD_WIZHELP = 2
const SCMD_DROP = 0
const SCMD_JUNK = 1
const SCMD_DONATE = 2
const SCMD_BUG = 0
const SCMD_TYPO = 1
const SCMD_IDEA = 2
const SCMD_LOOK = 0
const SCMD_READ = 1
const SCMD_SEARCH = 2
const SCMD_QSAY = 0
const SCMD_QECHO = 1
const SCMD_POUR = 0
const SCMD_FILL = 1
const SCMD_POOFIN = 0
const SCMD_POOFOUT = 1
const SCMD_HIT = 0
const SCMD_MURDER = 1
const SCMD_EAT = 0
const SCMD_TASTE = 1
const SCMD_DRINK = 2
const SCMD_SIP = 3
const SCMD_USE = 0
const SCMD_QUAFF = 1
const SCMD_RECITE = 2
const SCMD_ECHO = 0
const SCMD_EMOTE = 1
const SCMD_SMOTE = 2
const SCMD_OPEN = 0
const SCMD_CLOSE = 1
const SCMD_UNLOCK = 2
const SCMD_LOCK = 3
const SCMD_PICK = 4
const SCMD_OASIS_REDIT = 0
const SCMD_OASIS_OEDIT = 1
const SCMD_OASIS_ZEDIT = 2
const SCMD_OASIS_MEDIT = 3
const SCMD_OASIS_SEDIT = 4
const SCMD_OASIS_CEDIT = 5
const SCMD_OLC_SAVEINFO = 7
const SCMD_OASIS_RLIST = 8
const SCMD_OASIS_MLIST = 9
const SCMD_OASIS_OLIST = 10
const SCMD_OASIS_SLIST = 11
const SCMD_OASIS_ZLIST = 12
const SCMD_OASIS_TRIGEDIT = 13
const SCMD_OASIS_AEDIT = 14
const SCMD_OASIS_TLIST = 15
const SCMD_OASIS_LINKS = 16
const SCMD_OASIS_GEDIT = 17
const SCMD_OASIS_GLIST = 18
const SCMD_OASIS_HEDIT = 19
const SCMD_OASIS_HSEDIT = 20
const SCMD_RLIST = 0
const SCMD_OLIST = 1
const SCMD_MLIST = 2
const SCMD_TLIST = 3
const SCMD_SLIST = 4
const SCMD_GLIST = 5
const SCMD_MAKE = 0
const SCMD_BAKE = 1
const SCMD_BREW = 2
const SCMD_ASSEMBLE = 3
const SCMD_CRAFT = 4
const SCMD_FLETCH = 5
const SCMD_KNIT = 6
const SCMD_MIX = 7
const SCMD_THATCH = 8
const SCMD_WEAVE = 9
const SCMD_FORGE = 10
const SCMD_MEMORIZE = 1
const SCMD_FORGET = 2
const SCMD_STOP = 3
const SCMD_WHEN_SLOT = 4
const SCMD_WIMPY = 0
const SCMD_POWERATT = 1
const SCMD_COMBATEXP = 2
const SCMD_CAST = 0
const SCMD_ART = 1
const SCMD_TEDIT = 0
const SCMD_REDIT = 1
const SCMD_OEDIT = 2
const SCMD_MEDIT = 3

type command_info struct {
	Command          *byte
	Sort_as          *byte
	Minimum_position int8
	Command_pointer  CommandFunc
	Minimum_level    int16
	Minimum_admlevel int16
	Subcmd           int
}
type alias_data struct {
	Alias       *byte
	Replacement *byte
	Type        int
	Next        *alias_data
}

const NUM_TOKENS = 9
const RECON = 1
const USURP = 2
const UNSWITCH = 3

var disabled_first *disabled_data = nil
var complete_cmd_info []command_info
var cmd_info []command_info

func initCmdInfo() {
	cmd_info = []command_info{{Command: libc.CString("RESERVED"), Sort_as: libc.CString(""), Minimum_position: 0, Command_pointer: nil, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("north"), Sort_as: libc.CString("n"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NORTH}, {Command: libc.CString("east"), Sort_as: libc.CString("e"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_EAST}, {Command: libc.CString("south"), Sort_as: libc.CString("s"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SOUTH}, {Command: libc.CString("west"), Sort_as: libc.CString("w"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WEST}, {Command: libc.CString("up"), Sort_as: libc.CString("u"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_UP}, {Command: libc.CString("down"), Sort_as: libc.CString("d"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_DOWN}, {Command: libc.CString("northwest"), Sort_as: libc.CString("northw"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NW}, {Command: libc.CString("nw"), Sort_as: libc.CString("nw"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NW}, {Command: libc.CString("northeast"), Sort_as: libc.CString("northe"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NE}, {Command: libc.CString("ne"), Sort_as: libc.CString("ne"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NE}, {Command: libc.CString("southeast"), Sort_as: libc.CString("southe"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SE}, {Command: libc.CString("se"), Sort_as: libc.CString("se"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SE}, {Command: libc.CString("southwest"), Sort_as: libc.CString("southw"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SW}, {Command: libc.CString("sw"), Sort_as: libc.CString("sw"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SW}, {Command: libc.CString("i"), Sort_as: libc.CString("i"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_inventory(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("inside"), Sort_as: libc.CString("in"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_IN}, {Command: libc.CString("outside"), Sort_as: libc.CString("out"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_move(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_OUT}, {Command: libc.CString("absorb"), Sort_as: libc.CString("absor"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_absorb(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("at"), Sort_as: libc.CString("at"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_at(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("adrenaline"), Sort_as: libc.CString("adrenalin"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_adrenaline(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("advance"), Sort_as: libc.CString("adv"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_advance(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("aedit"), Sort_as: libc.CString("aed"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_OASIS_AEDIT}, {Command: libc.CString("alias"), Sort_as: libc.CString("ali"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_alias(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("afk"), Sort_as: libc.CString("afk"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_AFK}, {Command: libc.CString("aid"), Sort_as: libc.CString("aid"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_aid(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("amnesiac"), Sort_as: libc.CString("amnesia"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_amnisiac(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("appraise"), Sort_as: libc.CString("apprais"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_appraise(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("approve"), Sort_as: libc.CString("approve"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_approve(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("arena"), Sort_as: libc.CString("aren"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_arena(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("ashcloud"), Sort_as: libc.CString("ashclou"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ashcloud(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("assedit"), Sort_as: libc.CString("assed"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_assedit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("assist"), Sort_as: libc.CString("assis"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_assist(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("astat"), Sort_as: libc.CString("ast"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_astat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("ask"), Sort_as: libc.CString("ask"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spec_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_ASK}, {Command: libc.CString("attack"), Sort_as: libc.CString("attack"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_attack(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("auction"), Sort_as: libc.CString("auctio"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("augment"), Sort_as: libc.CString("augmen"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("aura"), Sort_as: libc.CString("aura"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_aura(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("autoexit"), Sort_as: libc.CString("autoex"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_autoexit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("autogold"), Sort_as: libc.CString("autogo"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_AUTOGOLD}, {Command: libc.CString("autoloot"), Sort_as: libc.CString("autolo"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_AUTOLOOT}, {Command: libc.CString("autosplit"), Sort_as: libc.CString("autosp"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_AUTOSPLIT}, {Command: libc.CString("bakuhatsuha"), Sort_as: libc.CString("baku"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_baku(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("ban"), Sort_as: libc.CString("ban"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ban(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_VICE, Subcmd: 0}, {Command: libc.CString("balance"), Sort_as: libc.CString("bal"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("balefire"), Sort_as: libc.CString("balef"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_balefire(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("barrage"), Sort_as: libc.CString("barrage"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pbarrage(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("barrier"), Sort_as: libc.CString("barri"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_barrier(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("bash"), Sort_as: libc.CString("bas"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_bash(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("beam"), Sort_as: libc.CString("bea"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_beam(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("bexchange"), Sort_as: libc.CString("bexchan"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rbanktrans(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("bid"), Sort_as: libc.CString("bi"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_bid(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("bigbang"), Sort_as: libc.CString("bigban"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_bigbang(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("bite"), Sort_as: libc.CString("bit"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_bite(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("blessedhammer"), Sort_as: libc.CString("bham"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_blessedhammer(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("block"), Sort_as: libc.CString("block"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_block(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("book"), Sort_as: libc.CString("boo"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_INFO}, {Command: libc.CString("break"), Sort_as: libc.CString("break"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_break(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("brief"), Sort_as: libc.CString("br"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_BRIEF}, {Command: libc.CString("build"), Sort_as: libc.CString("bui"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_assemble(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_BREW}, {Command: libc.CString("buildwalk"), Sort_as: libc.CString("buildwalk"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_BUILDWALK}, {Command: libc.CString("buy"), Sort_as: libc.CString("bu"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("bug"), Sort_as: libc.CString("bug"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_write(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_BUG}, {Command: libc.CString("cancel"), Sort_as: libc.CString("cance"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("candy"), Sort_as: libc.CString("cand"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_candy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("carry"), Sort_as: libc.CString("carr"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_carry(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("carve"), Sort_as: libc.CString("carv"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: SCMD_CARVE}, {Command: libc.CString("cedit"), Sort_as: libc.CString("cedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: SCMD_OASIS_CEDIT}, {Command: libc.CString("channel"), Sort_as: libc.CString("channe"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_channel(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("charge"), Sort_as: libc.CString("char"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_charge(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("check"), Sort_as: libc.CString("ch"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("checkload"), Sort_as: libc.CString("checkl"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_checkloadstatus(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("chown"), Sort_as: libc.CString("cho"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_chown(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("clan"), Sort_as: libc.CString("cla"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_clan(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("clear"), Sort_as: libc.CString("cle"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_CLEAR}, {Command: libc.CString("close"), Sort_as: libc.CString("cl"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_door(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_CLOSE}, {Command: libc.CString("closeeyes"), Sort_as: libc.CString("closeey"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_eyec(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("cls"), Sort_as: libc.CString("cls"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_CLEAR}, {Command: libc.CString("clsolc"), Sort_as: libc.CString("clsolc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_CLS}, {Command: libc.CString("consider"), Sort_as: libc.CString("con"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_consider(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("color"), Sort_as: libc.CString("col"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_color(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("combine"), Sort_as: libc.CString("comb"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_combine(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("compare"), Sort_as: libc.CString("comp"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_compare(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("commands"), Sort_as: libc.CString("com"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_commands(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_COMMANDS}, {Command: libc.CString("commune"), Sort_as: libc.CString("comm"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_commune(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("compact"), Sort_as: libc.CString("compact"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_COMPACT}, {Command: libc.CString("cook"), Sort_as: libc.CString("coo"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_cook(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("copyover"), Sort_as: libc.CString("copyover"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_copyover(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("create"), Sort_as: libc.CString("crea"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_form(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("credits"), Sort_as: libc.CString("cred"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_CREDITS}, {Command: libc.CString("crusher"), Sort_as: libc.CString("crushe"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_crusher(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("date"), Sort_as: libc.CString("da"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_date(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_DATE}, {Command: libc.CString("darkness"), Sort_as: libc.CString("darknes"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ddslash(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("dc"), Sort_as: libc.CString("dc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_dc(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("deathball"), Sort_as: libc.CString("deathbal"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_deathball(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("deathbeam"), Sort_as: libc.CString("deathbea"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_deathbeam(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("decapitate"), Sort_as: libc.CString("decapit"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spoil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("defend"), Sort_as: libc.CString("defen"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_defend(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("deploy"), Sort_as: libc.CString("deplo"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_deploy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("dualbeam"), Sort_as: libc.CString("dualbea"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_dualbeam(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("deposit"), Sort_as: libc.CString("depo"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("diagnose"), Sort_as: libc.CString("diagnos"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_diagnose(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("dimizu"), Sort_as: libc.CString("dimizu"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_dimizu(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("disable"), Sort_as: libc.CString("disa"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_disable(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_VICE, Subcmd: 0}, {Command: libc.CString("disguise"), Sort_as: libc.CString("disguis"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_disguise(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("dig"), Sort_as: libc.CString("dig"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_bury(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("display"), Sort_as: libc.CString("disp"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_display(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("dodonpa"), Sort_as: libc.CString("dodon"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_dodonpa(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("donate"), Sort_as: libc.CString("don"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_DONATE}, {Command: libc.CString("drag"), Sort_as: libc.CString("dra"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drag(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("draw"), Sort_as: libc.CString("dra"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_draw(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("drink"), Sort_as: libc.CString("dri"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drink(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_DRINK}, {Command: libc.CString("drop"), Sort_as: libc.CString("dro"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_DROP}, {Command: libc.CString("dub"), Sort_as: libc.CString("du"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_intro(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("eat"), Sort_as: libc.CString("ea"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_eat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_EAT}, {Command: libc.CString("eavesdrop"), Sort_as: libc.CString("eaves"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_eavesdrop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("echo"), Sort_as: libc.CString("ec"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_echo(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_ECHO}, {Command: libc.CString("elbow"), Sort_as: libc.CString("elb"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_elbow(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("emote"), Sort_as: libc.CString("em"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_echo(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_EMOTE}, {Command: libc.CString("energize"), Sort_as: libc.CString("energiz"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_energize(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString(":"), Sort_as: libc.CString(":"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_echo(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_EMOTE}, {Command: libc.CString("ensnare"), Sort_as: libc.CString("ensnar"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ensnare(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("enter"), Sort_as: libc.CString("ent"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_enter(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("equipment"), Sort_as: libc.CString("eq"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_equipment(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("eraser"), Sort_as: libc.CString("eras"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_eraser(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("escape"), Sort_as: libc.CString("esca"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_escape(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("evolve"), Sort_as: libc.CString("evolv"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_evolve(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("exchange"), Sort_as: libc.CString("exchan"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rptrans(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("exits"), Sort_as: libc.CString("ex"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_exits(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("examine"), Sort_as: libc.CString("exa"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_examine(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("extract"), Sort_as: libc.CString("extrac"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_extract(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("feed"), Sort_as: libc.CString("fee"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_feed(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("fill"), Sort_as: libc.CString("fil"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pour(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_FILL}, {Command: libc.CString("file"), Sort_as: libc.CString("fi"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_file(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("finalflash"), Sort_as: libc.CString("finalflash"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_final(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("finddoor"), Sort_as: libc.CString("findd"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_finddoor(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("findkey"), Sort_as: libc.CString("findk"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_findkey(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("finger"), Sort_as: libc.CString("finge"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_finger(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("fireshield"), Sort_as: libc.CString("firesh"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_fireshield(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("fish"), Sort_as: libc.CString("fis"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_fish(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("fix"), Sort_as: libc.CString("fix"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_fix(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("flee"), Sort_as: libc.CString("fl"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_flee(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("fly"), Sort_as: libc.CString("fly"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_fly(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("focus"), Sort_as: libc.CString("foc"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_focus(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("follow"), Sort_as: libc.CString("fol"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_follow(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("force"), Sort_as: libc.CString("force"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_force(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("forgery"), Sort_as: libc.CString("forg"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_forgery(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("forget"), Sort_as: libc.CString("forg"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("freeze"), Sort_as: libc.CString("freeze"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_FREEZE}, {Command: libc.CString("fury"), Sort_as: libc.CString("fury"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_fury(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("future"), Sort_as: libc.CString("futu"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_future(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("gain"), Sort_as: libc.CString("ga"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("galikgun"), Sort_as: libc.CString("galik"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_galikgun(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("game"), Sort_as: libc.CString("gam"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_show(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("garden"), Sort_as: libc.CString("garde"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_garden(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("genkidama"), Sort_as: libc.CString("genkidam"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_genki(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("genocide"), Sort_as: libc.CString("genocid"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_geno(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("get"), Sort_as: libc.CString("get"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_get(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("gecho"), Sort_as: libc.CString("gecho"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gecho(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("gedit"), Sort_as: libc.CString("gedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_OASIS_GEDIT}, {Command: libc.CString("gemote"), Sort_as: libc.CString("gem"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_GEMOTE}, {Command: libc.CString("generator"), Sort_as: libc.CString("genr"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("glist"), Sort_as: libc.CString("glist"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_OASIS_GLIST}, {Command: libc.CString("give"), Sort_as: libc.CString("giv"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_give(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("goto"), Sort_as: libc.CString("go"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_goto(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("gold"), Sort_as: libc.CString("gol"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gold(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("group"), Sort_as: libc.CString("gro"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_group(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("grab"), Sort_as: libc.CString("grab"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_grab(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("grand"), Sort_as: libc.CString("gran"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("grapple"), Sort_as: libc.CString("grapp"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_grapple(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("grats"), Sort_as: libc.CString("grat"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_GRATZ}, {Command: libc.CString("gravity"), Sort_as: libc.CString("grav"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("gsay"), Sort_as: libc.CString("gsay"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gsay(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("gtell"), Sort_as: libc.CString("gt"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gsay(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hand"), Sort_as: libc.CString("han"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hand(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("handout"), Sort_as: libc.CString("hand"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_handout(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("hasshuken"), Sort_as: libc.CString("hasshuke"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hass(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hayasa"), Sort_as: libc.CString("hayas"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hayasa(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("headbutt"), Sort_as: libc.CString("headbut"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_head(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("heal"), Sort_as: libc.CString("hea"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_heal(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("health"), Sort_as: libc.CString("hea"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_GHEALTH}, {Command: libc.CString("healingglow"), Sort_as: libc.CString("healing"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_healglow(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("heeldrop"), Sort_as: libc.CString("heeldr"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_heeldrop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hellflash"), Sort_as: libc.CString("hellflas"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hellflash(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hellspear"), Sort_as: libc.CString("hellspea"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hellspear(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("help"), Sort_as: libc.CString("h"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_help(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hedit"), Sort_as: libc.CString("hedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_HEDIT}, {Command: libc.CString("hindex"), Sort_as: libc.CString("hind"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hindex(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("helpcheck"), Sort_as: libc.CString("helpch"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_helpcheck(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("handbook"), Sort_as: libc.CString("handb"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_HANDBOOK}, {Command: libc.CString("hide"), Sort_as: libc.CString("hide"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_HIDE}, {Command: libc.CString("hints"), Sort_as: libc.CString("hints"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_HINTS}, {Command: libc.CString("history"), Sort_as: libc.CString("hist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_history(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hold"), Sort_as: libc.CString("hold"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_grab(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("holylight"), Sort_as: libc.CString("holy"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_HOLYLIGHT}, {Command: libc.CString("honoo"), Sort_as: libc.CString("hono"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_honoo(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("house"), Sort_as: libc.CString("house"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_house(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hsedit"), Sort_as: libc.CString("hsedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_OASIS_HSEDIT}, {Command: libc.CString("hspiral"), Sort_as: libc.CString("hspira"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hspiral(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("htank"), Sort_as: libc.CString("htan"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hydromancy"), Sort_as: libc.CString("hydrom"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hydromancy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("hyoga"), Sort_as: libc.CString("hyoga"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_obstruct(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("ihealth"), Sort_as: libc.CString("ihea"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_IHEALTH}, {Command: libc.CString("info"), Sort_as: libc.CString("info"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ginfo(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("infuse"), Sort_as: libc.CString("infus"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_infuse(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("ingest"), Sort_as: libc.CString("inges"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ingest(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("imotd"), Sort_as: libc.CString("imotd"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_IMOTD}, {Command: libc.CString("immlist"), Sort_as: libc.CString("imm"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WIZLIST}, {Command: libc.CString("implant"), Sort_as: libc.CString("implan"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_implant(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("instant"), Sort_as: libc.CString("insta"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_instant(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("instill"), Sort_as: libc.CString("instil"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_instill(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("instruct"), Sort_as: libc.CString("instruc"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: SCMD_INSTRUCT}, {Command: libc.CString("inventory"), Sort_as: libc.CString("inv"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_inventory(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("interest"), Sort_as: libc.CString("inter"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_interest(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("iedit"), Sort_as: libc.CString("ie"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_iedit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("invis"), Sort_as: libc.CString("invi"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_invis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("iwarp"), Sort_as: libc.CString("iwarp"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_warp(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("junk"), Sort_as: libc.CString("junk"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_JUNK}, {Command: libc.CString("kaioken"), Sort_as: libc.CString("kaioken"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kaioken(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kakusanha"), Sort_as: libc.CString("kakusan"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kakusanha(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kamehameha"), Sort_as: libc.CString("kame"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kamehameha(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kanso"), Sort_as: libc.CString("kans"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kanso(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kiball"), Sort_as: libc.CString("kibal"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kiball(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kiblast"), Sort_as: libc.CString("kiblas"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kiblast(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kienzan"), Sort_as: libc.CString("kienza"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kienzan(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kill"), Sort_as: libc.CString("kil"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kill(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("kick"), Sort_as: libc.CString("kic"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kick(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("knee"), Sort_as: libc.CString("kne"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_knee(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("koteiru"), Sort_as: libc.CString("koteiru"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_koteiru(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kousengan"), Sort_as: libc.CString("kousengan"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kousengan(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kuraiiro"), Sort_as: libc.CString("kuraiir"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kura(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("kyodaika"), Sort_as: libc.CString("kyodaik"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_kyodaika(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("look"), Sort_as: libc.CString("lo"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_look(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_LOOK}, {Command: libc.CString("lag"), Sort_as: libc.CString("la"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_lag(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 5, Subcmd: 0}, {Command: libc.CString("land"), Sort_as: libc.CString("lan"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_land(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("languages"), Sort_as: libc.CString("lang"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_languages(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("last"), Sort_as: libc.CString("last"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_last(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("learn"), Sort_as: libc.CString("lear"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("leave"), Sort_as: libc.CString("lea"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_leave(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("levels"), Sort_as: libc.CString("lev"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_levels(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("light"), Sort_as: libc.CString("ligh"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_lightgrenade(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("list"), Sort_as: libc.CString("lis"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("life"), Sort_as: libc.CString("lif"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_lifeforce(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("links"), Sort_as: libc.CString("lin"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_OASIS_LINKS}, {Command: libc.CString("liquefy"), Sort_as: libc.CString("liquef"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_liquefy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("lkeep"), Sort_as: libc.CString("lkee"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_LKEEP}, {Command: libc.CString("lock"), Sort_as: libc.CString("loc"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_door(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_LOCK}, {Command: libc.CString("lockout"), Sort_as: libc.CString("lock"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_hell(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("load"), Sort_as: libc.CString("load"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_load(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("majinize"), Sort_as: libc.CString("majini"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_majinize(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("malice"), Sort_as: libc.CString("malic"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_malice(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("masenko"), Sort_as: libc.CString("masenk"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_masenko(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("motd"), Sort_as: libc.CString("motd"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_MOTD}, {Command: libc.CString("mail"), Sort_as: libc.CString("mail"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 2, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("map"), Sort_as: libc.CString("map"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_map(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("medit"), Sort_as: libc.CString("medit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_MEDIT}, {Command: libc.CString("meditate"), Sort_as: libc.CString("medita"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_meditate(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("metamorph"), Sort_as: libc.CString("metamorp"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_metamorph(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mimic"), Sort_as: libc.CString("mimi"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mimic(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mlist"), Sort_as: libc.CString("mlist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_MLIST}, {Command: libc.CString("moondust"), Sort_as: libc.CString("moondus"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_moondust(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("multiform"), Sort_as: libc.CString("multifor"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_multiform(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mute"), Sort_as: libc.CString("mute"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_SQUELCH}, {Command: libc.CString("music"), Sort_as: libc.CString("musi"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_HOLLER}, {Command: libc.CString("newbie"), Sort_as: libc.CString("newbie"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_AUCTION}, {Command: libc.CString("news"), Sort_as: libc.CString("news"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_news(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("newsedit"), Sort_as: libc.CString("newsedi"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_newsedit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("nickname"), Sort_as: libc.CString("nicknam"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_nickname(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("nocompress"), Sort_as: libc.CString("nocompress"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOCOMPRESS}, {Command: libc.CString("noeq"), Sort_as: libc.CString("noeq"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOEQSEE}, {Command: libc.CString("nolin"), Sort_as: libc.CString("nolin"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NODEC}, {Command: libc.CString("nomusic"), Sort_as: libc.CString("nomusi"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOMUSIC}, {Command: libc.CString("noooc"), Sort_as: libc.CString("noooc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOGOSSIP}, {Command: libc.CString("nogive"), Sort_as: libc.CString("nogiv"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: SCMD_NOGIVE}, {Command: libc.CString("nograts"), Sort_as: libc.CString("nograts"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOGRATZ}, {Command: libc.CString("nogrow"), Sort_as: libc.CString("nogro"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_nogrow(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("nohassle"), Sort_as: libc.CString("nohassle"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_NOHASSLE}, {Command: libc.CString("nomail"), Sort_as: libc.CString("nomail"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NMWARN}, {Command: libc.CString("nonewbie"), Sort_as: libc.CString("nonewbie"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOAUCTION}, {Command: libc.CString("noparry"), Sort_as: libc.CString("noparr"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOPARRY}, {Command: libc.CString("norepeat"), Sort_as: libc.CString("norepeat"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOREPEAT}, {Command: libc.CString("noshout"), Sort_as: libc.CString("noshout"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_DEAF}, {Command: libc.CString("nosummon"), Sort_as: libc.CString("nosummon"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOSUMMON}, {Command: libc.CString("notell"), Sort_as: libc.CString("notell"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_NOTELL}, {Command: libc.CString("notitle"), Sort_as: libc.CString("notitle"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: SCMD_NOTITLE}, {Command: libc.CString("nova"), Sort_as: libc.CString("nov"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_nova(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("nowiz"), Sort_as: libc.CString("nowiz"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_NOWIZ}, {Command: libc.CString("ooc"), Sort_as: libc.CString("ooc"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_GOSSIP}, {Command: libc.CString("offer"), Sort_as: libc.CString("off"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("open"), Sort_as: libc.CString("ope"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_door(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_OPEN}, {Command: libc.CString("olc"), Sort_as: libc.CString("olc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_show_save_list(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("olist"), Sort_as: libc.CString("olist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_OLIST}, {Command: libc.CString("oedit"), Sort_as: libc.CString("oedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_OEDIT}, {Command: libc.CString("osay"), Sort_as: libc.CString("osay"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_osay(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("pack"), Sort_as: libc.CString("pac"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pack(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("page"), Sort_as: libc.CString("pag"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_page(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("paralyze"), Sort_as: libc.CString("paralyz"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_paralyze(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("pagelength"), Sort_as: libc.CString("pagel"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pagelength(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("peace"), Sort_as: libc.CString("pea"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_peace(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("perfect"), Sort_as: libc.CString("perfec"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_perf(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("permission"), Sort_as: libc.CString("permiss"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_permission(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("phoenix"), Sort_as: libc.CString("phoeni"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pslash(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("pick"), Sort_as: libc.CString("pi"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_door(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_PICK}, {Command: libc.CString("pickup"), Sort_as: libc.CString("picku"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("pilot"), Sort_as: libc.CString("pilot"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drive(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("plant"), Sort_as: libc.CString("plan"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_plant(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("play"), Sort_as: libc.CString("pla"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_play(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("players"), Sort_as: libc.CString("play"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_plist(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("poofin"), Sort_as: libc.CString("poofi"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_poofset(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_POOFIN}, {Command: libc.CString("poofout"), Sort_as: libc.CString("poofo"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_poofset(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_POOFOUT}, {Command: libc.CString("pose"), Sort_as: libc.CString("pos"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pose(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("post"), Sort_as: libc.CString("pos"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_post(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("potential"), Sort_as: libc.CString("poten"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_potential(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("pour"), Sort_as: libc.CString("pour"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pour(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_POUR}, {Command: libc.CString("powerup"), Sort_as: libc.CString("poweru"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_powerup(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("preference"), Sort_as: libc.CString("preferenc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_preference(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("program"), Sort_as: libc.CString("progra"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_OASIS_REDIT}, {Command: libc.CString("prompt"), Sort_as: libc.CString("pro"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_display(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("practice"), Sort_as: libc.CString("pra"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_practice(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("psychic"), Sort_as: libc.CString("psychi"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_psyblast(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("punch"), Sort_as: libc.CString("punc"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_punch(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("pushup"), Sort_as: libc.CString("pushu"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_pushup(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("put"), Sort_as: libc.CString("put"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_put(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("purge"), Sort_as: libc.CString("purge"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_purge(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("qui"), Sort_as: libc.CString("qui"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_quit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("quit"), Sort_as: libc.CString("quit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_quit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_QUIT}, {Command: libc.CString("radar"), Sort_as: libc.CString("rada"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_sradar(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("raise"), Sort_as: libc.CString("rai"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_raise(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rbank"), Sort_as: libc.CString("rban"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rbank(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("refuel"), Sort_as: libc.CString("refue"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_refuel(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("resize"), Sort_as: libc.CString("resiz"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_resize(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("restring"), Sort_as: libc.CString("restring"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_restring(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rclone"), Sort_as: libc.CString("rclon"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rcopy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("rcopy"), Sort_as: libc.CString("rcopy"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rcopy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("roomdisplay"), Sort_as: libc.CString("roomdisplay"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rdisplay(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("read"), Sort_as: libc.CString("rea"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_look(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_READ}, {Command: libc.CString("recall"), Sort_as: libc.CString("reca"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_recall(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("recharge"), Sort_as: libc.CString("rechar"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_recharge(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("regenerate"), Sort_as: libc.CString("regen"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_regenerate(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("renzokou"), Sort_as: libc.CString("renzo"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_renzo(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("repair"), Sort_as: libc.CString("repai"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_srepair(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("reply"), Sort_as: libc.CString("rep"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_reply(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rescue"), Sort_as: libc.CString("rescu"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rescue(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rest"), Sort_as: libc.CString("re"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rest(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("reward"), Sort_as: libc.CString("rewar"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_reward(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("reload"), Sort_as: libc.CString("reload"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_reboot(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 5, Subcmd: 0}, {Command: libc.CString("receive"), Sort_as: libc.CString("rece"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("remove"), Sort_as: libc.CString("rem"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_remove(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rent"), Sort_as: libc.CString("rent"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("report"), Sort_as: libc.CString("repor"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_write(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_IDEA}, {Command: libc.CString("reroll"), Sort_as: libc.CString("rero"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: SCMD_REROLL}, {Command: libc.CString("respond"), Sort_as: libc.CString("resp"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_respond(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("restore"), Sort_as: libc.CString("resto"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_restore(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("return"), Sort_as: libc.CString("retu"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_return(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("redit"), Sort_as: libc.CString("redit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_REDIT}, {Command: libc.CString("rip"), Sort_as: libc.CString("ri"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rip(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rlist"), Sort_as: libc.CString("rlist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_RLIST}, {Command: libc.CString("rogafufuken"), Sort_as: libc.CString("rogafu"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rogafufuken(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("roomflags"), Sort_as: libc.CString("roomf"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_ROOMFLAGS}, {Command: libc.CString("roundhouse"), Sort_as: libc.CString("roundhou"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_roundhouse(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rpbank"), Sort_as: libc.CString("rpban"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rpbank(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("rpp"), Sort_as: libc.CString("rpp"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_rpp(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("runic"), Sort_as: libc.CString("runi"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_runic(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("say"), Sort_as: libc.CString("say"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_say(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("'"), Sort_as: libc.CString("'"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_say(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("save"), Sort_as: libc.CString("sav"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_save(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("saveall"), Sort_as: libc.CString("saveall"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_saveall(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("sbc"), Sort_as: libc.CString("sbc"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_sbc(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("scan"), Sort_as: libc.CString("sca"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_scan(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("scatter"), Sort_as: libc.CString("scatte"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_scatter(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("score"), Sort_as: libc.CString("sc"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_score(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("scouter"), Sort_as: libc.CString("scou"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_scouter(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("scry"), Sort_as: libc.CString("scr"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_scry(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("seishou"), Sort_as: libc.CString("seisho"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_seishou(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("shell"), Sort_as: libc.CString("she"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shell(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("shimmer"), Sort_as: libc.CString("shimme"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shimmer(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("shogekiha"), Sort_as: libc.CString("shog"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shogekiha(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("shuffle"), Sort_as: libc.CString("shuff"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shuffle(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("snet"), Sort_as: libc.CString("snet"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_snet(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("search"), Sort_as: libc.CString("sea"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_look(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SEARCH}, {Command: libc.CString("sell"), Sort_as: libc.CString("sell"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("selfdestruct"), Sort_as: libc.CString("selfdest"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_selfd(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("sedit"), Sort_as: libc.CString("sedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_SEDIT}, {Command: libc.CString("send"), Sort_as: libc.CString("send"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_send(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("sense"), Sort_as: libc.CString("sense"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_track(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("set"), Sort_as: libc.CString("set"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_set(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("shout"), Sort_as: libc.CString("sho"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SHOUT}, {Command: libc.CString("show"), Sort_as: libc.CString("show"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_showoff(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("shutdow"), Sort_as: libc.CString("shutdow"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shutdown(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("shutdown"), Sort_as: libc.CString("shutdown"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_shutdown(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: SCMD_SHUTDOWN}, {Command: libc.CString("silk"), Sort_as: libc.CString("sil"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_silk(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("sip"), Sort_as: libc.CString("sip"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_drink(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SIP}, {Command: libc.CString("sit"), Sort_as: libc.CString("sit"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_sit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("situp"), Sort_as: libc.CString("situp"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_situp(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("skills"), Sort_as: libc.CString("skills"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_skills(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("skillset"), Sort_as: libc.CString("skillset"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_skillset(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 5, Subcmd: 0}, {Command: libc.CString("slam"), Sort_as: libc.CString("sla"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_slam(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("sleep"), Sort_as: libc.CString("sl"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_sleep(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("slist"), Sort_as: libc.CString("slist"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_SLIST}, {Command: libc.CString("slowns"), Sort_as: libc.CString("slowns"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: SCMD_SLOWNS}, {Command: libc.CString("smote"), Sort_as: libc.CString("sm"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_echo(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SMOTE}, {Command: libc.CString("sneak"), Sort_as: libc.CString("sneak"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SNEAK}, {Command: libc.CString("snoop"), Sort_as: libc.CString("snoop"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_snoop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("song"), Sort_as: libc.CString("son"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_song(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("spiral"), Sort_as: libc.CString("spiral"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spiral(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("socials"), Sort_as: libc.CString("socials"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_commands(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_SOCIALS}, {Command: libc.CString("solarflare"), Sort_as: libc.CString("solarflare"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_solar(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("spar"), Sort_as: libc.CString("spa"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spar(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("spit"), Sort_as: libc.CString("spi"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("spiritball"), Sort_as: libc.CString("spiritball"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spiritball(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("spiritcontrol"), Sort_as: libc.CString("spiritcontro"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spiritcontrol(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("split"), Sort_as: libc.CString("split"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_split(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("speak"), Sort_as: libc.CString("spe"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_languages(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("spells"), Sort_as: libc.CString("spel"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spells(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("stand"), Sort_as: libc.CString("st"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_stand(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("starbreaker"), Sort_as: libc.CString("starbr"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_breaker(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("stake"), Sort_as: libc.CString("stak"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_beacon(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("stat"), Sort_as: libc.CString("stat"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_stat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("status"), Sort_as: libc.CString("statu"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_status(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 0, Subcmd: 0}, {Command: libc.CString("steal"), Sort_as: libc.CString("ste"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_steal(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("stone"), Sort_as: libc.CString("ston"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("stop"), Sort_as: libc.CString("sto"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_stop(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("study"), Sort_as: libc.CString("stu"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("summon"), Sort_as: libc.CString("summo"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_summon(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("sunder"), Sort_as: libc.CString("sunde"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_sunder(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("suppress"), Sort_as: libc.CString("suppres"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_suppress(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("swallow"), Sort_as: libc.CString("swall"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_use(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_QUAFF}, {Command: libc.CString("switch"), Sort_as: libc.CString("switch"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_switch(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_VICE, Subcmd: 0}, {Command: libc.CString("syslog"), Sort_as: libc.CString("syslog"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_syslog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("tailhide"), Sort_as: libc.CString("tailh"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tailhide(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("table"), Sort_as: libc.CString("tabl"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_table(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("teach"), Sort_as: libc.CString("teac"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_teach(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tell"), Sort_as: libc.CString("tel"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tell(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("take"), Sort_as: libc.CString("tak"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_get(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tailwhip"), Sort_as: libc.CString("tailw"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tailwhip(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("taisha"), Sort_as: libc.CString("taish"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_taisha(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("taste"), Sort_as: libc.CString("tas"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_eat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_TASTE}, {Command: libc.CString("teleport"), Sort_as: libc.CString("tele"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_teleport(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("telepathy"), Sort_as: libc.CString("telepa"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_telepathy(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tedit"), Sort_as: libc.CString("tedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tedit(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GRGOD, Subcmd: 0}, {Command: libc.CString("test"), Sort_as: libc.CString("test"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: SCMD_TEST}, {Command: libc.CString("thaw"), Sort_as: libc.CString("thaw"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_THAW}, {Command: libc.CString("think"), Sort_as: libc.CString("thin"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_think(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("throw"), Sort_as: libc.CString("thro"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_throw(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("title"), Sort_as: libc.CString("title"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_title(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("time"), Sort_as: libc.CString("time"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_time(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("toggle"), Sort_as: libc.CString("toggle"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_toggle(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("toplist"), Sort_as: libc.CString("toplis"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_toplist(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("trackthru"), Sort_as: libc.CString("trackthru"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: SCMD_TRACK}, {Command: libc.CString("train"), Sort_as: libc.CString("train"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_train(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("transfer"), Sort_as: libc.CString("transfer"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_trans(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("transform"), Sort_as: libc.CString("transform"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_transform(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("transo"), Sort_as: libc.CString("trans"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_transobj(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: 5, Subcmd: 0}, {Command: libc.CString("tribeam"), Sort_as: libc.CString("tribe"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tribeam(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("trigedit"), Sort_as: libc.CString("trigedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_TRIGEDIT}, {Command: libc.CString("trip"), Sort_as: libc.CString("trip"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_trip(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tsuihidan"), Sort_as: libc.CString("tsuihida"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tsuihidan(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tunnel"), Sort_as: libc.CString("tunne"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_dig(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("twinslash"), Sort_as: libc.CString("twins"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tslash(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("twohand"), Sort_as: libc.CString("twohand"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_twohand(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("typo"), Sort_as: libc.CString("typo"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_write(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_TYPO}, {Command: libc.CString("unlock"), Sort_as: libc.CString("unlock"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_door(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_UNLOCK}, {Command: libc.CString("ungroup"), Sort_as: libc.CString("ungroup"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_ungroup(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("unban"), Sort_as: libc.CString("unban"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_unban(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GRGOD, Subcmd: 0}, {Command: libc.CString("unaffect"), Sort_as: libc.CString("unaffect"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizutil(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: SCMD_UNAFFECT}, {Command: libc.CString("uppercut"), Sort_as: libc.CString("upperc"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_uppercut(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("upgrade"), Sort_as: libc.CString("upgrad"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_upgrade(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("uptime"), Sort_as: libc.CString("uptime"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_date(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_UPTIME}, {Command: libc.CString("use"), Sort_as: libc.CString("use"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_use(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_USE}, {Command: libc.CString("users"), Sort_as: libc.CString("users"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_users(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("value"), Sort_as: libc.CString("val"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("varstat"), Sort_as: libc.CString("varst"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_varstat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("version"), Sort_as: libc.CString("ver"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_VERSION}, {Command: libc.CString("vieworder"), Sort_as: libc.CString("view"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_VIEWORDER}, {Command: libc.CString("visible"), Sort_as: libc.CString("vis"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_visible(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("vnum"), Sort_as: libc.CString("vnum"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_vnum(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("voice"), Sort_as: libc.CString("voic"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_voice(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("vstat"), Sort_as: libc.CString("vstat"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_vstat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("wake"), Sort_as: libc.CString("wa"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wake(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("warppool"), Sort_as: libc.CString("warppoo"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_warppool(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("waterrazor"), Sort_as: libc.CString("waterraz"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_razor(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("waterspikes"), Sort_as: libc.CString("waterspik"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spike(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("wear"), Sort_as: libc.CString("wea"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wear(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("weather"), Sort_as: libc.CString("weather"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_weather(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("who"), Sort_as: libc.CString("who"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_who(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("whoami"), Sort_as: libc.CString("whoami"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WHOAMI}, {Command: libc.CString("whohide"), Sort_as: libc.CString("whohide"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_tog(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WHOHIDE}, {Command: libc.CString("whois"), Sort_as: libc.CString("whois"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_whois(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("where"), Sort_as: libc.CString("where"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_where(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("whisper"), Sort_as: libc.CString("whisper"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_spec_comm(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WHISPER}, {Command: libc.CString("wield"), Sort_as: libc.CString("wie"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wield(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("will"), Sort_as: libc.CString("wil"), Minimum_position: POS_RESTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_willpower(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("wimpy"), Sort_as: libc.CString("wimpy"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_value(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WIMPY}, {Command: libc.CString("withdraw"), Sort_as: libc.CString("withdraw"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("wire"), Sort_as: libc.CString("wir"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_not_here(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("wiznet"), Sort_as: libc.CString("wiz"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wiznet(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString(";"), Sort_as: libc.CString(";"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wiznet(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("wizhelp"), Sort_as: libc.CString("wizhelp"), Minimum_position: POS_SLEEPING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_commands(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_WIZHELP}, {Command: libc.CString("wizlist"), Sort_as: libc.CString("wizlist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_gen_ps(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: SCMD_WIZLIST}, {Command: libc.CString("wizlock"), Sort_as: libc.CString("wizlock"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizlock(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("wizupdate"), Sort_as: libc.CString("wizupdate"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_wizupdate(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMPL, Subcmd: 0}, {Command: libc.CString("write"), Sort_as: libc.CString("write"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_write(ch, argument, cmd, subcmd)
	}, Minimum_level: 1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("zanzoken"), Sort_as: libc.CString("zanzo"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_zanzoken(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("zen"), Sort_as: libc.CString("ze"), Minimum_position: POS_FIGHTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_zen(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("zcheck"), Sort_as: libc.CString("zcheck"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_zcheck(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GOD, Subcmd: 0}, {Command: libc.CString("zreset"), Sort_as: libc.CString("zreset"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_zreset(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("zedit"), Sort_as: libc.CString("zedit"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_ZEDIT}, {Command: libc.CString("zlist"), Sort_as: libc.CString("zlist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_ZLIST}, {Command: libc.CString("zpurge"), Sort_as: libc.CString("zpurge"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_zpurge(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_GRGOD, Subcmd: 0}, {Command: libc.CString("attach"), Sort_as: libc.CString("attach"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_attach(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("detach"), Sort_as: libc.CString("detach"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_detach(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("detect"), Sort_as: libc.CString("detec"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_radar(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("tlist"), Sort_as: libc.CString("tlist"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_oasis(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: SCMD_OASIS_TLIST}, {Command: libc.CString("tstat"), Sort_as: libc.CString("tstat"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_tstat(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_IMMORT, Subcmd: 0}, {Command: libc.CString("masound"), Sort_as: libc.CString("masound"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_masound(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mheal"), Sort_as: libc.CString("mhea"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mheal(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mkill"), Sort_as: libc.CString("mkill"), Minimum_position: POS_STANDING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mkill(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mjunk"), Sort_as: libc.CString("mjunk"), Minimum_position: POS_SITTING, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mjunk(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mdamage"), Sort_as: libc.CString("mdamage"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mdamage(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mdoor"), Sort_as: libc.CString("mdoor"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mdoor(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mecho"), Sort_as: libc.CString("mecho"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mecho(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mechoaround"), Sort_as: libc.CString("mechoaround"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mechoaround(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("msend"), Sort_as: libc.CString("msend"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_msend(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mload"), Sort_as: libc.CString("mload"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mload(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mpurge"), Sort_as: libc.CString("mpurge"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mpurge(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mgoto"), Sort_as: libc.CString("mgoto"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mgoto(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mat"), Sort_as: libc.CString("mat"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mat(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mteleport"), Sort_as: libc.CString("mteleport"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mteleport(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mforce"), Sort_as: libc.CString("mforce"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mforce(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mremember"), Sort_as: libc.CString("mremember"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mremember(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mforget"), Sort_as: libc.CString("mforget"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mforget(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mtransform"), Sort_as: libc.CString("mtransform"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mtransform(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("mzoneecho"), Sort_as: libc.CString("mzoneecho"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mzoneecho(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("vdelete"), Sort_as: libc.CString("vdelete"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_vdelete(ch, argument, cmd, subcmd)
	}, Minimum_level: 0, Minimum_admlevel: ADMLVL_BUILDER, Subcmd: 0}, {Command: libc.CString("mfollow"), Sort_as: libc.CString("mfollow"), Minimum_position: POS_DEAD, Command_pointer: func(ch *char_data, argument *byte, cmd int, subcmd int) {
		do_mfollow(ch, argument, cmd, subcmd)
	}, Minimum_level: -1, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}, {Command: libc.CString("\n"), Sort_as: libc.CString("zzzzzzz"), Minimum_position: 0, Command_pointer: nil, Minimum_level: 0, Minimum_admlevel: ADMLVL_NONE, Subcmd: 0}}
}

var fill [9]*byte = [9]*byte{libc.CString("in"), libc.CString("into"), libc.CString("from"), libc.CString("with"), libc.CString("the"), libc.CString("on"), libc.CString("at"), libc.CString("to"), libc.CString("\n")}
var reserved [9]*byte = [9]*byte{libc.CString("a"), libc.CString("an"), libc.CString("self"), libc.CString("me"), libc.CString("all"), libc.CString("room"), libc.CString("someone"), libc.CString("something"), libc.CString("\n")}

func roll_stats(ch *char_data, type_ int, bonus int) int {
	var (
		pool       int = 0
		base_num   int = bonus
		max_num    int = bonus
		powerlevel int = 0
		ki         int = 1
		stamina    int = 2
	)
	if type_ == powerlevel {
		base_num = int(ch.Real_abils.Str) * 3
		max_num = int(ch.Real_abils.Str) * 5
	} else if type_ == ki {
		base_num = int(ch.Real_abils.Intel) * 3
		max_num = int(ch.Real_abils.Intel) * 5
	} else if type_ == stamina {
		base_num = int(ch.Real_abils.Con) * 3
		max_num = int(ch.Real_abils.Con) * 5
	}
	pool = rand_number(base_num, max_num) + bonus
	return pool
}
func command_interpreter(ch *char_data, argument *byte) {
	var (
		cmd     int
		length  int
		skip_ld int = 0
		line    *byte
		arg     [2048]byte
	)
	switch ch.Position {
	case POS_DEAD:
		fallthrough
	case POS_INCAP:
		fallthrough
	case POS_MORTALLYW:
		fallthrough
	case POS_STUNNED:
		ch.Position = POS_SITTING
	}
	skip_spaces(&argument)
	if *argument == 0 {
		return
	}
	if !libc.IsAlpha(rune(*argument)) {
		arg[0] = *argument
		arg[1] = '\x00'
		line = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	} else {
		line = any_one_arg(argument, &arg[0])
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("-")) == 0 {
		return
	}
	{
		var cont int
		cont = command_wtrigger(ch, &arg[0], line)
		if cont == 0 {
			cont = command_mtrigger(ch, &arg[0], line)
		}
		if cont == 0 {
			cont = command_otrigger(ch, &arg[0], line)
		}
		if cont != 0 {
			return
		}
	}
	for func() int {
		length = libc.StrLen(&arg[0])
		return func() int {
			cmd = 0
			return cmd
		}()
	}(); *complete_cmd_info[cmd].Command != '\n'; cmd++ {
		if libc.StrNCmp(complete_cmd_info[cmd].Command, &arg[0], length) == 0 {
			if GET_LEVEL(ch) >= int(complete_cmd_info[cmd].Minimum_level) && ch.Admlevel >= int(complete_cmd_info[cmd].Minimum_admlevel) {
				break
			}
		}
	}
	var blah [2048]byte
	stdio.Sprintf(&blah[0], "%s", complete_cmd_info[cmd].Command)
	if libc.StrCaseCmp(&blah[0], libc.CString("throw")) == 0 {
		ch.Throws = rand_number(1, 3)
	}
	if *complete_cmd_info[cmd].Command == '\n' {
		send_to_char(ch, libc.CString("Huh!?!\r\n"))
		return
	} else if command_pass(&blah[0], ch) == 0 && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("It's unfortunate...\r\n"))
	} else if check_disabled(&complete_cmd_info[cmd]) != 0 {
		send_to_char(ch, libc.CString("This command has been temporarily disabled.\r\n"))
	} else if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_GOOP) && ch.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("You only have your internal thoughts until your body has finished regenerating!\r\n"))
	} else if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_FROZEN) && ch.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("You try, but the mind-numbing cold prevents you...\r\n"))
	} else if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_SPIRAL) {
		send_to_char(ch, libc.CString("You are occupied with your Spiral Comet attack!\r\n"))
	} else if complete_cmd_info[cmd].Command_pointer == nil {
		send_to_char(ch, libc.CString("Sorry, that command hasn't been implemented yet.\r\n"))
	} else if IS_NPC(ch) && int(complete_cmd_info[cmd].Minimum_admlevel) >= ADMLVL_IMMORT {
		send_to_char(ch, libc.CString("You can't use immortal commands while switched.\r\n"))
	} else if int(ch.Position) < int(complete_cmd_info[cmd].Minimum_position) && int(ch.Position) != POS_FIGHTING {
		switch ch.Position {
		case POS_DEAD:
			send_to_char(ch, libc.CString("Lie still; you are DEAD!!! :-(\r\n"))
		case POS_INCAP:
			fallthrough
		case POS_MORTALLYW:
			send_to_char(ch, libc.CString("You are in a pretty bad shape, unable to do anything!\r\n"))
		case POS_STUNNED:
			send_to_char(ch, libc.CString("All you can do right now is think about the stars!\r\n"))
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("In your dreams, or what?\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("Nah... You feel too relaxed to do that..\r\n"))
		case POS_SITTING:
			send_to_char(ch, libc.CString("Maybe you should get on your feet first?\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("No way!  You're fighting for your life!\r\n"))
		}
	} else if no_specials != 0 || special(ch, cmd, line) == 0 {
		if skip_ld == 0 {
			(complete_cmd_info[cmd].Command_pointer)(ch, line, cmd, complete_cmd_info[cmd].Subcmd)
		}
	}
}
func find_alias(alias_list *alias_data, str *byte) *alias_data {
	for alias_list != nil {
		if *str == *alias_list.Alias {
			if libc.StrCmp(str, alias_list.Alias) == 0 {
				return alias_list
			}
		}
		alias_list = alias_list.Next
	}
	return nil
}
func free_alias(a *alias_data) {
	if a.Alias != nil {
		libc.Free(unsafe.Pointer(a.Alias))
	}
	if a.Replacement != nil {
		libc.Free(unsafe.Pointer(a.Replacement))
	}
	libc.Free(unsafe.Pointer(a))
}
func do_alias(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		repl *byte
		a    *alias_data
		temp *alias_data
	)
	if IS_NPC(ch) {
		return
	}
	repl = any_one_arg(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Currently defined aliases:\r\n"))
		if (func() *alias_data {
			a = ch.Player_specials.Aliases
			return a
		}()) == nil {
			send_to_char(ch, libc.CString(" None.\r\n"))
		} else {
			for a != nil {
				send_to_char(ch, libc.CString("%-15s %s\r\n"), a.Alias, a.Replacement)
				a = a.Next
			}
		}
	} else {
		if (func() *alias_data {
			a = find_alias(ch.Player_specials.Aliases, &arg[0])
			return a
		}()) != nil {
			if a == ch.Player_specials.Aliases {
				ch.Player_specials.Aliases = a.Next
			} else {
				temp = ch.Player_specials.Aliases
				for temp != nil && temp.Next != a {
					temp = temp.Next
				}
				if temp != nil {
					temp.Next = a.Next
				}
			}
			free_alias(a)
		}
		if *repl == 0 {
			if a == nil {
				send_to_char(ch, libc.CString("No such alias.\r\n"))
			} else {
				send_to_char(ch, libc.CString("Alias deleted.\r\n"))
			}
		} else {
			if libc.StrCaseCmp(&arg[0], libc.CString("alias")) == 0 {
				send_to_char(ch, libc.CString("You can't alias 'alias'.\r\n"))
				return
			}
			a = new(alias_data)
			a.Alias = libc.StrDup(&arg[0])
			delete_doubledollar(repl)
			a.Replacement = libc.StrDup(repl)
			if libc.StrChr(repl, ALIAS_SEP_CHAR) != nil || libc.StrChr(repl, ALIAS_VAR_CHAR) != nil {
				a.Type = ALIAS_COMPLEX
			} else {
				a.Type = ALIAS_SIMPLE
			}
			a.Next = ch.Player_specials.Aliases
			ch.Player_specials.Aliases = a
			send_to_char(ch, libc.CString("Alias added.\r\n"))
		}
	}
}
func perform_complex_alias(input_q *txt_q, orig *byte, a *alias_data) {
	var (
		temp_queue    txt_q
		tokens        [9]*byte
		temp          *byte
		write_point   *byte
		buf2          [4096]byte
		buf           [4096]byte
		num_of_tokens int = 0
		num           int
	)
	libc.StrCpy(&buf2[0], orig)
	temp = libc.StrTok(&buf2[0], libc.CString(" "))
	for temp != nil && num_of_tokens < NUM_TOKENS {
		tokens[func() int {
			p := &num_of_tokens
			x := *p
			*p++
			return x
		}()] = temp
		temp = libc.StrTok(nil, libc.CString(" "))
	}
	write_point = &buf[0]
	temp_queue.Head = func() *txt_block {
		p := &temp_queue.Tail
		temp_queue.Tail = nil
		return *p
	}()
	for temp = a.Replacement; *temp != 0; temp = (*byte)(unsafe.Add(unsafe.Pointer(temp), 1)) {
		if *temp == ALIAS_SEP_CHAR {
			*write_point = '\x00'
			buf[int(MAX_INPUT_LENGTH-1)] = '\x00'
			write_to_q(&buf[0], &temp_queue, 1)
			write_point = &buf[0]
		} else if *temp == ALIAS_VAR_CHAR {
			temp = (*byte)(unsafe.Add(unsafe.Pointer(temp), 1))
			if (func() int {
				num = int(*temp - '1')
				return num
			}()) < num_of_tokens && num >= 0 {
				libc.StrCpy(write_point, tokens[num])
				write_point = (*byte)(unsafe.Add(unsafe.Pointer(write_point), libc.StrLen(tokens[num])))
			} else if *temp == ALIAS_GLOB_CHAR {
				libc.StrCpy(write_point, orig)
				write_point = (*byte)(unsafe.Add(unsafe.Pointer(write_point), libc.StrLen(orig)))
			} else if (func() byte {
				p := (func() *byte {
					p := &write_point
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}())
				*(func() *byte {
					p := &write_point
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()) = *temp
				return *p
			}()) == '$' {
				*(func() *byte {
					p := &write_point
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()) = '$'
			}
		} else {
			*(func() *byte {
				p := &write_point
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *temp
		}
	}
	*write_point = '\x00'
	buf[int(MAX_INPUT_LENGTH-1)] = '\x00'
	write_to_q(&buf[0], &temp_queue, 1)
	if input_q.Head == nil {
		*input_q = temp_queue
	} else {
		temp_queue.Tail.Next = input_q.Head
		input_q.Head = temp_queue.Head
	}
}
func perform_alias(d *descriptor_data, orig *byte, maxlen uint64) int {
	var (
		first_arg [2048]byte
		ptr       *byte
		a         *alias_data
		tmp       *alias_data
	)
	if IS_NPC(d.Character) {
		return 0
	}
	if (func() *alias_data {
		tmp = d.Character.Player_specials.Aliases
		return tmp
	}()) == nil {
		return 0
	}
	ptr = any_one_arg(orig, &first_arg[0])
	if first_arg[0] == 0 {
		return 0
	}
	if (func() *alias_data {
		a = find_alias(tmp, &first_arg[0])
		return a
	}()) == nil {
		return 0
	}
	if a.Type == ALIAS_SIMPLE {
		strlcpy(orig, a.Replacement, maxlen)
		return 0
	} else {
		perform_complex_alias(&d.Input, ptr, a)
		return 1
	}
}
func search_block(arg *byte, list **byte, exact int) int {
	var (
		i int
		l int
	)
	if *arg == '!' {
		return -1
	}
	for l = 0; *((*byte)(unsafe.Add(unsafe.Pointer(arg), l))) != 0; l++ {
		*((*byte)(unsafe.Add(unsafe.Pointer(arg), l))) = byte(int8(unicode.ToLower(rune(*((*byte)(unsafe.Add(unsafe.Pointer(arg), l)))))))
	}
	if exact != 0 {
		for i = 0; **((**byte)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof((*byte)(nil))*uintptr(i)))) != '\n'; i++ {
			if libc.StrCmp(arg, *((**byte)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof((*byte)(nil))*uintptr(i))))) == 0 {
				return i
			}
		}
	} else {
		if l == 0 {
			l = 1
		}
		for i = 0; **((**byte)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof((*byte)(nil))*uintptr(i)))) != '\n'; i++ {
			if libc.StrNCmp(arg, *((**byte)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof((*byte)(nil))*uintptr(i)))), l) == 0 {
				return i
			}
		}
	}
	return -1
}
func is_number(str *byte) int {
	for *str != 0 {
		if !unicode.IsDigit(rune(*(func() *byte {
			p := &str
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()))) {
			return 0
		}
	}
	return 1
}
func skip_spaces(string_ **byte) {
	for ; **string_ != 0 && unicode.IsSpace(rune(**string_)); *string_ = (*byte)(unsafe.Add(unsafe.Pointer(*string_), 1)) {
	}
}
func delete_doubledollar(string_ *byte) *byte {
	var (
		ddread  *byte
		ddwrite *byte
	)
	if (func() *byte {
		ddwrite = libc.StrChr(string_, '$')
		return ddwrite
	}()) == nil {
		return string_
	}
	ddread = ddwrite
	for *ddread != 0 {
		if (func() byte {
			p := (func() *byte {
				p := &ddwrite
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			*(func() *byte {
				p := &ddwrite
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *(func() *byte {
				p := &ddread
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			return *p
		}()) == '$' {
			if *ddread == '$' {
				ddread = (*byte)(unsafe.Add(unsafe.Pointer(ddread), 1))
			}
		}
	}
	*ddwrite = '\x00'
	return string_
}
func fill_word(argument *byte) int {
	return int(libc.BoolToInt(search_block(argument, &fill[0], TRUE) >= 0))
}
func topLoad() {
	var (
		file   *stdio.File
		fname  [40]byte
		line   [256]byte
		filler [50]byte
		x      int = 0
	)
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("toplist")) == 0 {
		basic_mud_log(libc.CString("ERROR: Toplist file does not exist."))
		return
	} else if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: Toplist file does not exist."))
		return
	}
	TOPLOADED = TRUE
	for int(file.IsEOF()) == 0 {
		get_line(file, &line[0])
		stdio.Sscanf(&line[0], "%s %lld\n", &filler[0], &toppoint[x])
		topname[x] = libc.StrDup(&filler[0])
		filler[0] = '\x00'
		x++
	}
	file.Close()
}
func topWrite(ch *char_data) {
	if ch.Admlevel > 0 || IS_NPC(ch) {
		return
	}
	if TOPLOADED == FALSE {
		return
	}
	var fname [40]byte
	var fl *stdio.File
	var positions [25]*byte
	var points [25]int64 = [25]int64{}
	var x int = 0
	var writeEm int = FALSE
	var placed int = FALSE
	var start int = 0
	var finish int = 25
	var location int = -1
	var progress int = FALSE
	if ch == nil {
		return
	}
	if ch.Desc == nil || GET_USER(ch) == nil {
		return
	}
	for x = start; x < finish; x++ {
		positions[x] = libc.StrDup(topname[x])
		points[x] = toppoint[x]
	}
	start = 0
	finish = 5
	for x = start; x < finish; x++ {
		if placed == FALSE {
			if libc.StrCaseCmp(topname[x], GET_NAME(ch)) != 0 {
				if ch.Max_hit > toppoint[x] {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = ch.Max_hit
					topname[x] = libc.StrDup(GET_NAME(ch))
					placed = TRUE
					writeEm = TRUE
					location = x
				}
			} else {
				placed = TRUE
				location = finish
			}
		} else {
			if x < finish && location < finish {
				if libc.StrCaseCmp(positions[location], GET_NAME(ch)) != 0 {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				} else {
					progress = TRUE
					location += 1
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				}
			}
		}
	}
	if progress == TRUE {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas moved up in rank in the powerlevel section.@D]\r\n"), GET_NAME(ch))
	} else if placed == TRUE && location != finish {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas placed in the powerlevel section.@D]\r\n"), GET_NAME(ch))
	}
	location = -1
	placed = FALSE
	progress = FALSE
	start = 5
	finish = 10
	for x = start; x < finish; x++ {
		if placed == FALSE {
			if libc.StrCaseCmp(topname[x], GET_NAME(ch)) != 0 {
				if ch.Max_mana > toppoint[x] {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = ch.Max_mana
					topname[x] = libc.StrDup(GET_NAME(ch))
					placed = TRUE
					writeEm = TRUE
					location = x
				}
			} else {
				placed = TRUE
				location = finish
			}
		} else {
			if x < finish && location < finish {
				if libc.StrCaseCmp(positions[location], GET_NAME(ch)) != 0 {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				} else {
					progress = TRUE
					location += 1
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				}
			}
		}
	}
	if progress == TRUE {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas moved up in rank in the ki section.@D]\r\n"), GET_NAME(ch))
	} else if placed == TRUE && location != finish {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas placed in the ki section.@D]\r\n"), GET_NAME(ch))
	}
	location = -1
	placed = FALSE
	progress = FALSE
	start = 10
	finish = 15
	for x = start; x < finish; x++ {
		if placed == FALSE {
			if libc.StrCaseCmp(topname[x], GET_NAME(ch)) != 0 {
				if ch.Max_move > toppoint[x] {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = ch.Max_move
					topname[x] = libc.StrDup(GET_NAME(ch))
					placed = TRUE
					writeEm = TRUE
					location = x
				}
			} else {
				placed = TRUE
				location = finish
			}
		} else {
			if x < finish && location < finish {
				if libc.StrCaseCmp(positions[location], GET_NAME(ch)) != 0 {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				} else {
					progress = TRUE
					location += 1
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				}
			}
		}
	}
	if progress == TRUE {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas moved up in rank in the stamina section.@D]\r\n"), GET_NAME(ch))
	} else if placed == TRUE && location != finish {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas placed in the stamina section.@D]\r\n"), GET_NAME(ch))
	}
	location = -1
	placed = FALSE
	progress = FALSE
	start = 15
	finish = 20
	for x = start; x < finish; x++ {
		if placed == FALSE {
			if libc.StrCaseCmp(topname[x], GET_NAME(ch)) != 0 {
				if ch.Bank_gold+ch.Gold > int(toppoint[x]) {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = int64(ch.Bank_gold + ch.Gold)
					topname[x] = libc.StrDup(GET_NAME(ch))
					placed = TRUE
					writeEm = TRUE
					location = x
				}
			} else {
				placed = TRUE
				location = finish
			}
		} else {
			if x < finish && location < finish {
				if libc.StrCaseCmp(positions[location], GET_NAME(ch)) != 0 {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				} else {
					progress = TRUE
					location += 1
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				}
			}
		}
	}
	if progress == TRUE {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas moved up in rank in the zenni section.@D]\r\n"), GET_NAME(ch))
	} else if placed == TRUE && location != finish {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas placed in the zenni section.@D]\r\n"), GET_NAME(ch))
	}
	location = -1
	placed = FALSE
	progress = FALSE
	start = 20
	finish = 25
	for x = start; x < finish; x++ {
		if placed == FALSE {
			if libc.StrCaseCmp(topname[x], GET_USER(ch)) != 0 {
				if ch.Trp > int(toppoint[x]) {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = int64(ch.Trp)
					topname[x] = libc.StrDup(GET_USER(ch))
					placed = TRUE
					writeEm = TRUE
					location = x
				}
			} else {
				placed = TRUE
				location = finish
			}
		} else {
			if x < finish && location < finish {
				if libc.StrCaseCmp(positions[location], GET_USER(ch)) != 0 {
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				} else {
					progress = TRUE
					location += 1
					libc.Free(unsafe.Pointer(topname[x]))
					toppoint[x] = points[location]
					topname[x] = libc.StrDup(positions[location])
					location += 1
				}
			}
		}
	}
	if progress == TRUE {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas moved up in rank in the RPP section.@D]\r\n"), GET_USER(ch))
	} else if placed == TRUE && location != finish {
		send_to_all(libc.CString("@D[@GToplist@W: @C%s @Whas placed in the RPP section.@D]\r\n"), GET_USER(ch))
	}
	location = -1
	placed = FALSE
	progress = FALSE
	for x = 0; x < 25; x++ {
		libc.Free(unsafe.Pointer(positions[x]))
	}
	if writeEm == TRUE {
		if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("toplist")) == 0 {
			return
		}
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: could not save Toplist File, %s."), &fname[0])
			return
		}
		x = 0
		for x < 25 {
			stdio.Fprintf(fl, "%s %lld\n", topname[x], toppoint[x])
			x++
		}
		fl.Close()
	}
	return
}
func reserved_word(argument *byte) int {
	return int(libc.BoolToInt(search_block(argument, &reserved[0], TRUE) >= 0))
}
func one_argument(argument *byte, first_arg *byte) *byte {
	var begin *byte = first_arg
	if argument == nil {
		*first_arg = '\x00'
		return nil
	}
	for {
		skip_spaces(&argument)
		first_arg = begin
		for *argument != 0 && !unicode.IsSpace(rune(*argument)) {
			*(func() *byte {
				p := &first_arg
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = byte(int8(unicode.ToLower(rune(*argument))))
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		*first_arg = '\x00'
		if fill_word(begin) == 0 {
			break
		}
	}
	return argument
}
func one_word(argument *byte, first_arg *byte) *byte {
	skip_spaces(&argument)
	if *argument == '"' {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for *argument != 0 && *argument != '"' {
			*(func() *byte {
				p := &first_arg
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = byte(int8(unicode.ToLower(rune(*argument))))
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	} else {
		for *argument != 0 && !unicode.IsSpace(rune(*argument)) {
			*(func() *byte {
				p := &first_arg
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = byte(int8(unicode.ToLower(rune(*argument))))
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
	}
	*first_arg = '\x00'
	return argument
}
func any_one_arg(argument *byte, first_arg *byte) *byte {
	skip_spaces(&argument)
	for *argument != 0 && !unicode.IsSpace(rune(*argument)) {
		*(func() *byte {
			p := &first_arg
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = byte(int8(unicode.ToLower(rune(*argument))))
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	*first_arg = '\x00'
	return argument
}
func two_arguments(argument *byte, first_arg *byte, second_arg *byte) *byte {
	return one_argument(one_argument(argument, first_arg), second_arg)
}
func three_arguments(argument *byte, first_arg *byte, second_arg *byte, third_arg *byte) *byte {
	return one_argument(one_argument(one_argument(argument, first_arg), second_arg), third_arg)
}
func display_races(d *descriptor_data) {
	var (
		x    int
		i    int       = 0
		cost [23]*byte = [23]*byte{libc.CString("Free"), libc.CString("60 RPP"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("35 RPP"), libc.CString("Free"), libc.CString("Free"), libc.CString("55 RPP"), libc.CString("Free"), libc.CString("Free"), libc.CString("30 RPP"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free"), libc.CString("Free")}
	)
	send_to_char(d.Character, libc.CString("\r\n@YRace SELECTION menu:\r\n@D---------------------------------------\r\n@n"))
	for x = 0; x < NUM_RACES; x++ {
		if race_ok_gender[int(d.Character.Sex)][x] {
			send_to_char(d.Character, libc.CString("@B%2d@W) @C%-15s@D[@R%-6s@D]@n%s"), func() int {
				if (x + 1) != 21 {
					return x + 1
				}
				return 16
			}(), pc_race_types[x], cost[x], func() string {
				if (func() int {
					p := &i
					*p++
					return *p
				}() % 2) == 0 {
					return "\r\n"
				}
				return "   "
			}())
		}
	}
	send_to_char(d.Character, libc.CString("\n @BR@W) @CRandom Race Selection!\r\n@n"))
	send_to_char(d.Character, libc.CString("\n @BT@W) @CToggle between SELECTION/HELP Menu\r\n@n"))
	send_to_char(d.Character, libc.CString("\n@WRace: @n"))
}
func display_classes(d *descriptor_data) {
	var (
		x int
		i int = 0
	)
	send_to_char(d.Character, libc.CString("\r\n@YSensei SELECTION menu:\r\n@D--------------------------------------\r\n@n"))
	for x = 0; x < 14; x++ {
		if class_ok_race[int(d.Character.Race)][x] {
			send_to_char(d.Character, libc.CString("@B%2d@W) @C%s@n%s"), x+1, pc_class_types[x], func() string {
				if (func() int {
					p := &i
					*p++
					return *p
				}() % 2) == 0 {
					return "\r\n"
				}
				return "\t"
			}())
		}
	}
	send_to_char(d.Character, libc.CString("\n @BR@W) @CRandom Sensei Selection!\r\n@n"))
	send_to_char(d.Character, libc.CString("\n @BT@W) @CToggle between SELECTION/HELP Menu\r\n@n"))
	send_to_char(d.Character, libc.CString("\n@WSensei: @n"))
}
func display_races_help(d *descriptor_data) {
	var (
		x int
		i int = 0
	)
	send_to_char(d.Character, libc.CString("\r\n@YRace HELP menu:\r\n@G--------------------------------------------\r\n@n"))
	for x = 0; x < NUM_RACES; x++ {
		if race_ok_gender[int(d.Character.Sex)][x] {
			send_to_char(d.Character, libc.CString("@B%2d@W) @C%-15s@n%s"), x+1, pc_race_types[x], func() string {
				if (func() int {
					p := &i
					*p++
					return *p
				}() % 2) == 0 {
					return "\r\n"
				}
				return "\t"
			}())
		}
	}
	send_to_char(d.Character, libc.CString("\n @BT@W) @CToggle between SELECTION/HELP Menu\r\n@n"))
	send_to_char(d.Character, libc.CString("\n@WHelp on Race #: @n"))
}
func display_classes_help(d *descriptor_data) {
	var (
		x int
		i int = 0
	)
	send_to_char(d.Character, libc.CString("\r\n@YClass HELP menu:\r\n@G-------------------------------------------\r\n@n"))
	for x = 0; x < 14; x++ {
		if class_ok_race[int(d.Character.Race)][x] {
			send_to_char(d.Character, libc.CString("@B%2d@W) @C%s@n%s"), x+1, pc_class_types[x], func() string {
				if (func() int {
					p := &i
					*p++
					return *p
				}() % 2) == 0 {
					return "\r\n"
				}
				return "\t"
			}())
		}
	}
	send_to_char(d.Character, libc.CString("\n @BT@W) @CToggle between SELECTION/HELP Menu\r\n@n"))
	send_to_char(d.Character, libc.CString("\n@WHelp on Class #: @n"))
}
func is_abbrev(arg1 *byte, arg2 *byte) int {
	if *arg1 == 0 {
		return 0
	}
	for ; *arg1 != 0 && *arg2 != 0; func() *byte {
		arg1 = (*byte)(unsafe.Add(unsafe.Pointer(arg1), 1))
		return func() *byte {
			p := &arg2
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
	}() {
		if unicode.ToLower(rune(*arg1)) != unicode.ToLower(rune(*arg2)) {
			return 0
		}
	}
	if *arg1 == 0 {
		return 1
	} else {
		return 0
	}
}
func half_chop(string_ *byte, arg1 *byte, arg2 *byte) {
	var temp *byte
	temp = any_one_arg(string_, arg1)
	skip_spaces(&temp)
	if arg2 != temp {
		libc.StrCpy(arg2, temp)
	}
}
func find_command(command *byte) int {
	var cmd int
	for cmd = 0; *complete_cmd_info[cmd].Command != '\n'; cmd++ {
		if libc.StrCmp(complete_cmd_info[cmd].Command, command) == 0 {
			return cmd
		}
	}
	return -1
}
func special(ch *char_data, cmd int, arg *byte) int {
	var (
		i *obj_data
		k *char_data
		j int
	)
	if GET_ROOM_SPEC(ch.In_room) != nil {
		if GET_ROOM_SPEC(ch.In_room)(ch, unsafe.Pointer(&world[ch.In_room]), cmd, arg) != 0 {
			return 1
		}
	}
	for j = 0; j < NUM_WEARS; j++ {
		if (ch.Equipment[j]) != nil && GET_OBJ_SPEC(ch.Equipment[j]) != nil {
			if GET_OBJ_SPEC(ch.Equipment[j])(ch, unsafe.Pointer(ch.Equipment[j]), cmd, arg) != 0 {
				return 1
			}
		}
	}
	for i = ch.Carrying; i != nil; i = i.Next_content {
		if GET_OBJ_SPEC(i) != nil {
			if GET_OBJ_SPEC(i)(ch, unsafe.Pointer(i), cmd, arg) != 0 {
				return 1
			}
		}
	}
	for k = world[ch.In_room].People; k != nil; k = k.Next_in_room {
		if !MOB_FLAGGED(k, MOB_NOTDEADYET) {
			if GET_MOB_SPEC(k) != nil && GET_MOB_SPEC(k)(ch, unsafe.Pointer(k), cmd, arg) != 0 {
				return 1
			}
		}
	}
	for i = world[ch.In_room].Contents; i != nil; i = i.Next_content {
		if GET_OBJ_SPEC(i) != nil {
			if GET_OBJ_SPEC(i)(ch, unsafe.Pointer(i), cmd, arg) != 0 {
				return 1
			}
		}
	}
	return 0
}
func _parse_name(arg *byte, name *byte) int {
	var i int
	skip_spaces(&arg)
	for i = 0; (func() byte {
		p := name
		*name = *arg
		return *p
	}()) != 0; func() *byte {
		arg = (*byte)(unsafe.Add(unsafe.Pointer(arg), 1))
		i++
		return func() *byte {
			p := &name
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
	}() {
		if !libc.IsAlpha(rune(*arg)) {
			return 1
		}
	}
	if i == 0 {
		return 1
	}
	return 0
}
func perform_dupe_check(d *descriptor_data) int {
	var (
		k       *descriptor_data
		next_k  *descriptor_data
		target  *char_data = nil
		ch      *char_data
		next_ch *char_data
		mode    int = 0
		id      int = int(d.Character.Idnum)
	)
	for k = descriptor_list; k != nil; k = next_k {
		next_k = k.Next
		if k == d {
			continue
		}
		if k.Original != nil && int(k.Original.Idnum) == id {
			write_to_output(d, libc.CString("\r\nMultiple login detected -- disconnecting.\r\n"))
			k.Connected = CON_CLOSE
			if target == nil {
				target = k.Original
				mode = UNSWITCH
			}
			if k.Character != nil {
				k.Character.Desc = nil
			}
			k.Character = nil
			k.Original = nil
		} else if k.Character != nil && int(k.Character.Idnum) == id && k.Original != nil {
			do_return(k.Character, nil, 0, 0)
		} else if k.Character != nil && int(k.Character.Idnum) == id {
			if target == nil && k.Connected == CON_PLAYING {
				write_to_output(k, libc.CString("\r\nThis body has been usurped!\r\n"))
				if k.Snoop_by != nil {
					k.Snoop_by.Snooping = d
					d.Snoop_by = k.Snoop_by
					k.Snoop_by = nil
				}
				target = k.Character
				mode = USURP
			}
			k.Character.Desc = nil
			k.Character = nil
			k.Original = nil
			write_to_output(k, libc.CString("\r\nMultiple login detected -- disconnecting.\r\n"))
			k.Connected = CON_CLOSE
		}
	}
	for ch = character_list; ch != nil; ch = next_ch {
		next_ch = ch.Next
		if IS_NPC(ch) {
			continue
		}
		if int(ch.Idnum) != id {
			continue
		}
		if ch.Desc != nil {
			continue
		}
		if ch == target {
			continue
		}
		if target == nil {
			target = ch
			mode = RECON
			continue
		}
		if ch.In_room != room_rnum(-1) {
			char_from_room(ch)
		}
		char_to_room(ch, 1)
		extract_char(ch)
	}
	if target == nil {
		return 0
	}
	free_char(d.Character)
	d.Character = target
	d.Character.Desc = d
	d.Original = nil
	d.Character.Timer = 0
	REMOVE_BIT_AR(d.Character.Act[:], PLR_MAILING)
	REMOVE_BIT_AR(d.Character.Act[:], PLR_WRITING)
	REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_GROUP)
	d.Connected = CON_PLAYING
	switch mode {
	case RECON:
		write_to_output(d, libc.CString("Reconnecting.\r\n"))
		var count int = 0
		var oldcount int = HIGHPCOUNT
		var k *descriptor_data
		for k = descriptor_list; k != nil; k = k.Next {
			if !IS_NPC(k.Character) && GET_LEVEL(k.Character) > 3 {
				count += 1
			}
			if count > PCOUNT {
				PCOUNT = count
			}
			if PCOUNT >= HIGHPCOUNT {
				oldcount = HIGHPCOUNT
				HIGHPCOUNT = PCOUNT
				PCOUNTDATE = libc.GetTime(nil)
			}
		}
		if PCOUNT < HIGHPCOUNT && PCOUNT >= HIGHPCOUNT-4 {
			payout(0)
		}
		if PCOUNT == HIGHPCOUNT {
			payout(1)
		}
		if PCOUNT > oldcount {
			payout(2)
		}
		d.Character.Time.Logon = libc.GetTime(nil)
		act(libc.CString("$n has reconnected."), TRUE, d.Character, nil, nil, TO_ROOM)
		mudlog(NRM, int(MAX(ADMLVL_NONE, int64(d.Character.Player_specials.Invis_level))), TRUE, libc.CString("%s [%s] has reconnected."), GET_NAME(d.Character), &d.Host[0])
		d.Character.Rp = d.Rpp
		if has_mail(int(d.Character.Idnum)) != 0 {
			write_to_output(d, libc.CString("You have mail waiting.\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWIMM > (d.Character.Lboard[1]) {
			send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the immortal board, there is a new post there.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWCOD > (d.Character.Lboard[2]) {
			send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the request file, it has been updated.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWBUI > (d.Character.Lboard[4]) {
			send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the builder board, there is a new post there.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWDUO > (d.Character.Lboard[3]) {
			send_to_char(d.Character, libc.CString("\r\n@GMake sure to check punishment board, there is a new post there.@n\r\n"))
		}
		if BOARDNEWMORT > (d.Character.Lboard[0]) {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new bulletin board post.@n\r\n"))
		}
		if NEWSUPDATE > d.Character.Lastpl {
			send_to_char(d.Character, libc.CString("\r\n@GThe NEWS file has been updated, type 'news %d' to see the latest entry or 'news list' to see available entries.@n\r\n"), LASTNEWS)
		}
		if LASTINTEREST != 0 && LASTINTEREST > d.Character.Lastint {
			var (
				diff int = int(LASTINTEREST - d.Character.Lastint)
				mult int = 0
			)
			for diff > 0 {
				if (diff-86400) < 0 && mult == 0 {
					mult = 1
				} else if (diff - 86400) >= 0 {
					diff -= 86400
					mult++
				} else {
					diff = 0
				}
			}
			if mult > 3 {
				mult = 3
			}
			d.Character.Lastint = LASTINTEREST
			if d.Character.Bank_gold > 0 {
				var inc int = ((d.Character.Bank_gold / 100) * 2)
				if inc >= 7500 {
					inc = 7500
				}
				inc *= mult
				d.Character.Bank_gold += inc
				send_to_char(d.Character, libc.CString("Interest happened while you were away, %d times.\r\n@cBank Interest@D: @Y%s@n\r\n"), mult, add_commas(int64(inc)))
			}
		}
	case USURP:
		write_to_output(d, libc.CString("You take over your own body, already in use!\r\n"))
		act(libc.CString("$n suddenly keels over in pain, surrounded by a white aura...\r\n$n's body has been taken over by a new spirit!"), TRUE, d.Character, nil, nil, TO_ROOM)
		d.Character.Rp = d.Rpp
		mudlog(NRM, int(MAX(ADMLVL_IMMORT, int64(d.Character.Player_specials.Invis_level))), TRUE, libc.CString("%s has re-logged in ... disconnecting old socket."), GET_NAME(d.Character))
	case UNSWITCH:
		write_to_output(d, libc.CString("Reconnecting to unswitched char."))
		mudlog(NRM, int(MAX(ADMLVL_IMMORT, int64(d.Character.Player_specials.Invis_level))), TRUE, libc.CString("%s [%s] has reconnected."), GET_NAME(d.Character), &d.Host[0])
	}
	return 1
}
func enter_player_game(d *descriptor_data) int {
	var (
		load_result int
		load_room   int64
		check       *char_data
	)
	reset_char(d.Character)
	read_aliases(d.Character)
	racial_body_parts(d.Character)
	if PLR_FLAGGED(d.Character, PLR_INVSTART) {
		d.Character.Player_specials.Invis_level = int16(GET_LEVEL(d.Character))
	}
	if (func() int64 {
		load_room = int64(d.Character.Player_specials.Load_room)
		return load_room
	}()) != int64(-1) {
		load_room = int64(real_room(room_vnum(load_room)))
	}
	if load_room == int64(-1) {
		if d.Character.Admlevel != 0 {
			load_room = int64(real_room(config_info.Room_nums.Immort_start_room))
		} else {
			load_room = int64(real_room(config_info.Room_nums.Mortal_start_room))
		}
	}
	if PLR_FLAGGED(d.Character, PLR_FROZEN) {
		load_room = int64(real_room(config_info.Room_nums.Frozen_start_room))
	}
	d.Character.Next = character_list
	character_list = d.Character
	char_to_room(d.Character, room_rnum(load_room))
	load_result = Crash_load(d.Character)
	if d.Character.Player_specials.Host != nil {
		libc.Free(unsafe.Pointer(d.Character.Player_specials.Host))
		d.Character.Player_specials.Host = nil
	}
	d.Character.Player_specials.Host = libc.StrDup(&d.Host[0])
	d.Character.Id = d.Character.Idnum
	add_to_lookup_table(int(d.Character.Id), unsafe.Pointer(d.Character))
	read_saved_vars(d.Character)
	for check = character_list; check != nil; check = check.Next {
		if check.Master == nil && IS_NPC(check) && int(check.Master_id) == int(d.Character.Idnum) && AFF_FLAGGED(check, AFF_CHARM) && !circle_follow(check, d.Character) {
			add_follower(check, d.Character)
		}
	}
	save_char(d.Character)
	if d.Customfile != 1 {
		customCreate(d)
		userWrite(d, 0, 0, 0, libc.CString("index"))
	}
	if d.Character.Lifeforce == 0 || d.Character.Lifeforce > int64(GET_LIFEMAX(d.Character)) {
		d.Character.Lifeforce = int64(GET_LIFEMAX(d.Character))
	}
	if PLR_FLAGGED(d.Character, PLR_RARM) {
		d.Character.Limb_condition[0] = 100
		REMOVE_BIT_AR(d.Character.Act[:], PLR_RARM)
	}
	if PLR_FLAGGED(d.Character, PLR_LARM) {
		d.Character.Limb_condition[1] = 100
		REMOVE_BIT_AR(d.Character.Act[:], PLR_LARM)
	}
	if PLR_FLAGGED(d.Character, PLR_LLEG) {
		d.Character.Limb_condition[3] = 100
		REMOVE_BIT_AR(d.Character.Act[:], PLR_LLEG)
	}
	if PLR_FLAGGED(d.Character, PLR_RLEG) {
		d.Character.Limb_condition[2] = 100
		REMOVE_BIT_AR(d.Character.Act[:], PLR_RLEG)
	}
	d.Character.Combine = -1
	d.Character.Sleeptime = 8
	d.Character.Foodr = 2
	if AFF_FLAGGED(d.Character, AFF_FLYING) {
		d.Character.Altitude = 1
	} else {
		d.Character.Altitude = 0
	}
	if AFF_FLAGGED(d.Character, AFF_POSITION) {
		REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_POSITION)
	}
	if AFF_FLAGGED(d.Character, AFF_SANCTUARY) {
		REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_SANCTUARY)
	}
	if AFF_FLAGGED(d.Character, AFF_ZANZOKEN) {
		REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_ZANZOKEN)
	}
	if PLR_FLAGGED(d.Character, PLR_KNOCKED) {
		REMOVE_BIT_AR(d.Character.Act[:], PLR_KNOCKED)
	}
	if int(d.Character.Race) == RACE_ANDROID && !AFF_FLAGGED(d.Character, AFF_INFRAVISION) {
		SET_BIT_AR(d.Character.Affected_by[:], AFF_INFRAVISION)
	}
	d.Character.Absorbing = nil
	d.Character.Absorbby = nil
	d.Character.Sits = nil
	d.Character.Blocked = nil
	d.Character.Blocks = nil
	d.Character.Overf = FALSE
	d.Character.Spam = 0
	d.Character.Rage_meter = 0
	if d.Character.Affected == nil {
		if AFF_FLAGGED(d.Character, AFF_HEALGLOW) {
			REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_HEALGLOW)
		}
	}
	if AFF_FLAGGED(d.Character, AFF_HAYASA) {
		d.Character.Speedboost = int(float64(GET_SPEEDCALC(d.Character)) * 0.5)
	} else {
		d.Character.Speedboost = 0
	}
	if d.Character.Trp < d.Character.Rp {
		d.Character.Trp = d.Character.Rp
	}
	if int(d.Character.Race) == RACE_NAMEK && int(d.Character.Player_specials.Conditions[HUNGER]) >= 0 {
		d.Character.Player_specials.Conditions[HUNGER] = -1
	}
	if PLR_FLAGGED(d.Character, PLR_HEALT) {
		REMOVE_BIT_AR(d.Character.Act[:], PLR_HEALT)
	}
	if readIntro(d.Character, d.Character) == 2 {
		introCreate(d.Character)
	}
	if read_sense_memory(d.Character, d.Character) == 2 {
		senseCreate(d.Character)
	}
	if d.Character.Admlevel > 0 {
		d.Level = 1
	}
	if d.Character.Clan != nil && libc.StrStr(d.Character.Clan, libc.CString("None")) == nil {
		if !clanIsMember(d.Character.Clan, d.Character) {
			if !clanIsModerator(d.Character.Clan, d.Character) {
				if checkCLAN(d.Character) == 0 {
					write_to_output(d, libc.CString("Your clan no longer exists.\r\n"))
					d.Character.Clan = libc.CString("None.")
				}
			}
		}
	}
	d.Character.Rp = d.Rpp
	d.Character.Rbank = d.Rbank
	if MOON_UP == FALSE && (int(d.Character.Race) == RACE_SAIYAN || int(d.Character.Race) == RACE_HALFBREED) {
		oozaru_drop(d.Character)
	}
	if MOON_UP == TRUE && (int(d.Character.Race) == RACE_SAIYAN || int(d.Character.Race) == RACE_HALFBREED) {
		oozaru_add(d.Character)
	}
	if int(d.Character.Race) == RACE_HOSHIJIN {
		if time_info.Day <= 14 {
			star_phase(d.Character, 1)
		} else if time_info.Day <= 21 {
			star_phase(d.Character, 2)
		} else {
			star_phase(d.Character, 0)
		}
	}
	if int(d.Character.Race) == RACE_ICER && GET_SKILL(d.Character, SKILL_TAILWHIP) == 0 {
		var numb int = rand_number(20, 30)
		for {
			d.Character.Skills[SKILL_TAILWHIP] = int8(numb)
			if true {
				break
			}
		}
	} else if int(d.Character.Race) != RACE_ICER && GET_SKILL(d.Character, SKILL_TAILWHIP) != 0 {
		for {
			d.Character.Skills[SKILL_TAILWHIP] = 0
			if true {
				break
			}
		}
	}
	if int(d.Character.Race) == RACE_MUTANT && ((d.Character.Genome[0]) == 9 || (d.Character.Genome[1]) == 9) && GET_SKILL(d.Character, SKILL_TELEPATHY) == 0 {
		for {
			d.Character.Skills[SKILL_TELEPATHY] = 50
			if true {
				break
			}
		}
	}
	if int(d.Character.Race) == RACE_BIO && ((d.Character.Genome[0]) == 7 || (d.Character.Genome[1]) == 7) && GET_SKILL(d.Character, SKILL_TELEPATHY) == 0 && GET_SKILL(d.Character, SKILL_FOCUS) == 0 {
		for {
			d.Character.Skills[SKILL_TELEPATHY] = 30
			if true {
				break
			}
		}
		for {
			d.Character.Skills[SKILL_FOCUS] = 30
			if true {
				break
			}
		}
	}
	d.Character.Combo = -1
	return load_result
}
func readUserIndex(name *byte) int {
	var (
		fname [40]byte
		fl    *stdio.File
	)
	if get_filename(&fname[0], uint64(40), USER_FILE, name) == 0 {
		return 0
	} else if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		return 0
	}
	fl.Close()
	return 1
}
func payout(num int) {
	var k *descriptor_data
	if LASTPAYOUT == 0 {
		LASTPAYOUT = libc.GetTime(nil) + 86400
		LASTPAYTYPE = num
	} else if num > LASTPAYTYPE {
		LASTPAYOUT = libc.GetTime(nil) + 86400
		LASTPAYTYPE = num
	} else if LASTPAYOUT <= libc.GetTime(nil) {
		LASTPAYOUT = libc.GetTime(nil) + 86400
		LASTPAYTYPE = num
	}
	for k = descriptor_list; k != nil; k = k.Next {
		if k.Character.Admlevel <= 0 && IS_PLAYING(k) && k.Character.Rewtime < LASTPAYOUT {
			if num == 0 {
				k.Rpp += 1
				k.Character.Rp = k.Rpp
				k.Character.Trp += 1
				userWrite(k, 0, 0, 0, libc.CString("index"))
				send_to_char(k.Character, libc.CString("@D[@G+ 1 RPP@D] @cA total logon count within 4 of the highest has been achieved.@n\r\n"))
			} else if num == 1 {
				k.Rpp += 2
				k.Character.Rp = k.Rpp
				k.Character.Trp += 2
				userWrite(k, 0, 0, 0, libc.CString("index"))
				send_to_char(k.Character, libc.CString("@D[@G+ 2 RPP@D] @cThe total logon count has tied with the highest ever.@n\r\n"))
			} else {
				k.Rpp += 3
				k.Character.Rp = k.Rpp
				k.Character.Trp += 3
				userWrite(k, 0, 0, 0, libc.CString("index"))
				send_to_char(k.Character, libc.CString("@D[@G+ 3 RPP@D] @cA new logon count record has been achieved!@n\r\n"))
			}
			k.Character.Rewtime = LASTPAYOUT
		}
	}
}
func command_pass(cmd *byte, ch *char_data) int {
	if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
		if libc.StrCaseCmp(cmd, libc.CString("liquefy")) != 0 && libc.StrCaseCmp(cmd, libc.CString("ingest")) != 0 && libc.StrCaseCmp(cmd, libc.CString("look")) != 0 && libc.StrCaseCmp(cmd, libc.CString("score")) != 0 && libc.StrCaseCmp(cmd, libc.CString("ooc")) != 0 && libc.StrCaseCmp(cmd, libc.CString("osay")) != 0 && libc.StrCaseCmp(cmd, libc.CString("emote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("smote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("status")) != 0 {
			send_to_char(ch, libc.CString("You are not capable of performing that action while liquefied!\r\n"))
			return FALSE
		}
	} else if AFF_FLAGGED(ch, AFF_PARALYZE) {
		if libc.StrCaseCmp(cmd, libc.CString("look")) != 0 && libc.StrCaseCmp(cmd, libc.CString("score")) != 0 && libc.StrCaseCmp(cmd, libc.CString("ooc")) != 0 && libc.StrCaseCmp(cmd, libc.CString("osay")) != 0 && libc.StrCaseCmp(cmd, libc.CString("emote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("smote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("status")) != 0 {
			send_to_char(ch, libc.CString("You are not capable of performing that action while petrified!\r\n"))
			return FALSE
		}
	} else if AFF_FLAGGED(ch, AFF_FROZEN) {
		if libc.StrCaseCmp(cmd, libc.CString("look")) != 0 && libc.StrCaseCmp(cmd, libc.CString("score")) != 0 && libc.StrCaseCmp(cmd, libc.CString("ooc")) != 0 && libc.StrCaseCmp(cmd, libc.CString("osay")) != 0 && libc.StrCaseCmp(cmd, libc.CString("emote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("smote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("status")) != 0 {
			send_to_char(ch, libc.CString("You are not capable of performing that action while a frozen block of ice!\r\n"))
			return FALSE
		}
	} else if AFF_FLAGGED(ch, AFF_PARA) && int(ch.Aff_abils.Intel) < rand_number(1, 60) {
		if libc.StrCaseCmp(cmd, libc.CString("look")) != 0 && libc.StrCaseCmp(cmd, libc.CString("score")) != 0 && libc.StrCaseCmp(cmd, libc.CString("ooc")) != 0 && libc.StrCaseCmp(cmd, libc.CString("osay")) != 0 && libc.StrCaseCmp(cmd, libc.CString("emote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("smote")) != 0 && libc.StrCaseCmp(cmd, libc.CString("status")) != 0 {
			act(libc.CString("@yYou fail to overcome your paralysis!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@Y$n @ystruggles with $s paralysis!@n"), TRUE, ch, nil, nil, TO_ROOM)
			return FALSE
		}
	}
	return TRUE
}
func lockRead(name *byte) int {
	var (
		fname  [40]byte
		filler [50]byte
		line   [256]byte
		known  int = FALSE
		fl     *stdio.File
	)
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("lockout")) == 0 {
		return 0
	} else if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		return 0
	}
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%s\n", &filler[0])
		if libc.StrCaseCmp(CAP(name), CAP(&filler[0])) == 0 {
			known = TRUE
		}
	}
	fl.Close()
	if known == TRUE {
		return 1
	} else {
		return 0
	}
}
func userLoad(d *descriptor_data, name *byte) {
	var (
		fname  [40]byte
		filler [100]byte
		line   [256]byte
		fl     *stdio.File
	)
	if get_filename(&fname[0], uint64(40), USER_FILE, name) == 0 {
		return
	} else if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not load user, %s, from filename, %s."), name, &fname[0])
		return
	}
	var count int = 0
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		count += 1
		switch count {
		case 1:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.User != nil {
				libc.Free(unsafe.Pointer(d.User))
				d.User = nil
			}
			d.User = libc.StrDup(&filler[0])
		case 2:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Email != nil {
				libc.Free(unsafe.Pointer(d.Email))
				d.Email = nil
			}
			d.Email = libc.StrDup(&filler[0])
		case 3:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Pass != nil {
				libc.Free(unsafe.Pointer(d.Pass))
				d.Pass = nil
			}
			d.Pass = libc.StrDup(&filler[0])
		case 4:
			stdio.Sscanf(&line[0], "%d\n", &d.Total)
		case 5:
			stdio.Sscanf(&line[0], "%d\n", &d.Rpp)
		case 6:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Tmp1 != nil {
				libc.Free(unsafe.Pointer(d.Tmp1))
				d.Tmp1 = nil
			}
			d.Tmp1 = libc.StrDup(&filler[0])
		case 7:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Tmp2 != nil {
				libc.Free(unsafe.Pointer(d.Tmp2))
				d.Tmp2 = nil
			}
			d.Tmp2 = libc.StrDup(&filler[0])
		case 8:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Tmp3 != nil {
				libc.Free(unsafe.Pointer(d.Tmp3))
				d.Tmp3 = nil
			}
			d.Tmp3 = libc.StrDup(&filler[0])
		case 9:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Tmp4 != nil {
				libc.Free(unsafe.Pointer(d.Tmp4))
				d.Tmp4 = nil
			}
			d.Tmp4 = libc.StrDup(&filler[0])
		case 10:
			stdio.Sscanf(&line[0], "%s\n", &filler[0])
			if d.Tmp5 != nil {
				libc.Free(unsafe.Pointer(d.Tmp5))
				d.Tmp5 = nil
			}
			d.Tmp5 = libc.StrDup(&filler[0])
		case 11:
			stdio.Sscanf(&line[0], "%d\n", &d.Level)
		case 12:
			stdio.Sscanf(&line[0], "%d\n", &d.Customfile)
		case 13:
			stdio.Sscanf(&line[0], "%d\n", &d.Rbank)
		}
		filler[0] = '\x00'
	}
	fl.Close()
	return
}
func userCreate(d *descriptor_data) {
	var (
		fname [40]byte
		fl    *stdio.File
	)
	if get_filename(&fname[0], uint64(40), USER_FILE, d.User) == 0 {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save user, %s, to filename, %s."), d.User, &fname[0])
		return
	}
	stdio.Fprintf(fl, "%s\n", CAP(d.User))
	stdio.Fprintf(fl, "%s\n", d.Email)
	stdio.Fprintf(fl, "%s\n", d.Pass)
	stdio.Fprintf(fl, "3\n")
	d.Total = 3
	stdio.Fprintf(fl, "0\n")
	d.Rpp = 0
	stdio.Fprintf(fl, "Empty\n")
	stdio.Fprintf(fl, "Empty\n")
	stdio.Fprintf(fl, "Empty\n")
	stdio.Fprintf(fl, "Empty\n")
	stdio.Fprintf(fl, "Empty\n")
	stdio.Fprintf(fl, "0\n")
	stdio.Fprintf(fl, "1\n")
	stdio.Fprintf(fl, "0\n")
	d.Rbank = 0
	fl.Close()
	if d.Tmp1 != nil {
		libc.Free(unsafe.Pointer(d.Tmp1))
		d.Tmp1 = nil
	}
	d.Tmp1 = libc.CString("Empty")
	if d.Tmp2 != nil {
		libc.Free(unsafe.Pointer(d.Tmp2))
		d.Tmp2 = nil
	}
	d.Tmp2 = libc.CString("Empty")
	if d.Tmp3 != nil {
		libc.Free(unsafe.Pointer(d.Tmp3))
		d.Tmp3 = nil
	}
	d.Tmp3 = libc.CString("Empty")
	if d.Tmp4 != nil {
		libc.Free(unsafe.Pointer(d.Tmp4))
		d.Tmp4 = nil
	}
	d.Tmp4 = libc.CString("Empty")
	if d.Tmp5 != nil {
		libc.Free(unsafe.Pointer(d.Tmp5))
		d.Tmp5 = nil
	}
	d.Tmp5 = libc.CString("Empty")
	customCreate(d)
	return
}
func userDelete(d *descriptor_data) {
	var (
		player_i int
		fname    [40]byte
	)
	if get_filename(&fname[0], uint64(40), USER_FILE, d.User) != 0 {
		if libc.StrCaseCmp(d.Tmp1, libc.CString("Empty")) != 0 {
			if (func() int {
				player_i = get_ptable_by_name(d.Tmp1)
				return player_i
			}()) >= 0 {
				remove_player(player_i)
			}
		}
		if libc.StrCaseCmp(d.Tmp2, libc.CString("Empty")) != 0 {
			if (func() int {
				player_i = get_ptable_by_name(d.Tmp2)
				return player_i
			}()) >= 0 {
				remove_player(player_i)
			}
		}
		if libc.StrCaseCmp(d.Tmp3, libc.CString("Empty")) != 0 {
			if (func() int {
				player_i = get_ptable_by_name(d.Tmp3)
				return player_i
			}()) >= 0 {
				remove_player(player_i)
			}
		}
		if libc.StrCaseCmp(d.Tmp4, libc.CString("Empty")) != 0 {
			if (func() int {
				player_i = get_ptable_by_name(d.Tmp4)
				return player_i
			}()) >= 0 {
				remove_player(player_i)
			}
		}
		if libc.StrCaseCmp(d.Tmp5, libc.CString("Empty")) != 0 {
			if (func() int {
				player_i = get_ptable_by_name(d.Tmp5)
				return player_i
			}()) >= 0 {
				remove_player(player_i)
			}
		}
		stdio.Unlink(&fname[0])
		return
	} else {
		write_to_output(d, libc.CString("Error. Your user file doesn't even exist!\n"))
		return
	}
}
func rIntro(ch *char_data, arg *byte) *byte {
	var (
		fname  [40]byte
		filler [50]byte
		scrap  [100]byte
		line   [256]byte
		name   [80]byte
		known  int = FALSE
		fl     *stdio.File
	)
	if IS_NPC(ch) {
		return libc.CString("NOTHING")
	}
	if get_filename(&fname[0], uint64(40), INTRO_FILE, GET_NAME(ch)) == 0 {
		return libc.CString("NOTHING")
	} else if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return fl
	}()) == nil {
		return libc.CString("NOTHING")
	}
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		stdio.Sscanf(&line[0], "%s %s\n", &filler[0], &scrap[0])
		if libc.StrCaseCmp(arg, &scrap[0]) == 0 {
			known = TRUE
			stdio.Sprintf(&name[0], "%s", &filler[0])
		}
	}
	fl.Close()
	if known == TRUE {
		return &name[0]
	} else {
		return libc.CString("NOTHING")
	}
}
func userWrite(d *descriptor_data, setTot int, setRpp int, setRBank int, name *byte) {
	var (
		fname [40]byte
		fl    *stdio.File
	)
	if libc.StrCaseCmp(name, libc.CString("index")) == 0 {
		if d == nil {
			return
		}
		if d.User == nil {
			return
		}
		if get_filename(&fname[0], uint64(40), USER_FILE, d.User) == 0 {
			return
		}
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: could not save user, %s, to filename, %s."), d.User, &fname[0])
			return
		}
		stdio.Fprintf(fl, "%s\n", CAP(d.User))
		stdio.Fprintf(fl, "%s\n", d.Email)
		stdio.Fprintf(fl, "%s\n", d.Pass)
		if setTot <= 3 || setTot > 5 {
			stdio.Fprintf(fl, "%d\n", d.Total)
		} else if setTot > 3 {
			stdio.Fprintf(fl, "%d\n", setTot)
			d.Total = setTot
		}
		if setRpp == 0 {
			stdio.Fprintf(fl, "%d\n", d.Rpp)
		} else {
			d.Rpp += setRpp
			stdio.Fprintf(fl, "%d\n", d.Rpp)
		}
		if d.Tmp1 != nil {
			stdio.Fprintf(fl, "%s\n", d.Tmp1)
		} else {
			stdio.Fprintf(fl, "Empty\n")
		}
		if d.Tmp2 != nil {
			stdio.Fprintf(fl, "%s\n", d.Tmp2)
		} else {
			stdio.Fprintf(fl, "Empty\n")
		}
		if d.Tmp3 != nil {
			stdio.Fprintf(fl, "%s\n", d.Tmp3)
		} else {
			stdio.Fprintf(fl, "Empty\n")
		}
		if d.Tmp4 != nil {
			stdio.Fprintf(fl, "%s\n", d.Tmp4)
		} else {
			stdio.Fprintf(fl, "Empty\n")
		}
		if d.Tmp5 != nil {
			stdio.Fprintf(fl, "%s\n", d.Tmp5)
		} else {
			stdio.Fprintf(fl, "Empty\n")
		}
		stdio.Fprintf(fl, "%d\n", d.Level)
		stdio.Fprintf(fl, "%d\n", d.Customfile)
		if setRBank == 0 {
			stdio.Fprintf(fl, "%d\n", d.Rbank)
		} else {
			d.Rbank += setRBank
			stdio.Fprintf(fl, "%d\n", d.Rbank)
		}
		fl.Close()
		return
	} else if libc.StrCaseCmp(name, libc.CString("index")) != 0 {
		var (
			filename [40]byte
			uname    [100]byte
			email    [100]byte
			pass     [100]byte
			tmp1     [100]byte
			tmp2     [100]byte
			tmp3     [100]byte
			tmp4     [100]byte
			tmp5     [100]byte
			line     [256]byte
			total    int = 0
			rpp      int = 0
			level    int = 0
			custom   int = 0
			rbank    int = 0
			file     *stdio.File
		)
		if get_filename(&filename[0], uint64(40), USER_FILE, name) == 0 {
			return
		} else if (func() *stdio.File {
			file = stdio.FOpen(libc.GoString(&filename[0]), "r")
			return file
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: could not load user, %s, from filename, %s."), name, &filename[0])
			return
		}
		var count int = 0
		for int(file.IsEOF()) == 0 {
			get_line(file, &line[0])
			count += 1
			switch count {
			case 1:
				stdio.Sscanf(&line[0], "%s\n", &uname[0])
			case 2:
				stdio.Sscanf(&line[0], "%s\n", &email[0])
			case 3:
				stdio.Sscanf(&line[0], "%s\n", &pass[0])
			case 4:
				stdio.Sscanf(&line[0], "%d\n", &total)
			case 5:
				stdio.Sscanf(&line[0], "%d\n", &rpp)
			case 6:
				stdio.Sscanf(&line[0], "%s\n", &tmp1[0])
			case 7:
				stdio.Sscanf(&line[0], "%s\n", &tmp2[0])
			case 8:
				stdio.Sscanf(&line[0], "%s\n", &tmp3[0])
			case 9:
				stdio.Sscanf(&line[0], "%s\n", &tmp4[0])
			case 10:
				stdio.Sscanf(&line[0], "%s\n", &tmp5[0])
			case 11:
				stdio.Sscanf(&line[0], "%d\n", &level)
			case 12:
				stdio.Sscanf(&line[0], "%d\n", &custom)
			case 13:
				stdio.Sscanf(&line[0], "%d\n", &rbank)
			}
		}
		file.Close()
		if get_filename(&fname[0], uint64(40), USER_FILE, name) == 0 {
			return
		}
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("ERROR: could not save user, %s, to filename, %s."), name, &fname[0])
			return
		}
		stdio.Fprintf(fl, "%s\n", &uname[0])
		stdio.Fprintf(fl, "%s\n", &email[0])
		stdio.Fprintf(fl, "%s\n", &pass[0])
		if setTot <= 3 || setTot > 5 {
			stdio.Fprintf(fl, "%d\n", total)
		} else if setTot > 3 {
			stdio.Fprintf(fl, "%d\n", setTot)
		}
		if setRpp == 0 {
			stdio.Fprintf(fl, "%d\n", rpp)
		} else if rpp+setRpp < 0 {
			send_to_imm(libc.CString("RPP would be below 0, reward canceled."))
			stdio.Fprintf(fl, "%d\n", rpp)
		} else {
			rpp += setRpp
			stdio.Fprintf(fl, "%d\n", rpp)
		}
		stdio.Fprintf(fl, "%s\n", &tmp1[0])
		stdio.Fprintf(fl, "%s\n", &tmp2[0])
		stdio.Fprintf(fl, "%s\n", &tmp3[0])
		stdio.Fprintf(fl, "%s\n", &tmp4[0])
		stdio.Fprintf(fl, "%s\n", &tmp5[0])
		stdio.Fprintf(fl, "%d\n", level)
		stdio.Fprintf(fl, "%d\n", custom)
		if setRBank == 0 {
			stdio.Fprintf(fl, "%d\n", rbank)
		} else if rbank+setRBank < 0 {
			send_to_imm(libc.CString("RPP Bank would be below 0, reward canceled."))
			stdio.Fprintf(fl, "%d\n", rbank)
		} else {
			rbank += setRBank
			stdio.Fprintf(fl, "%d\n", rbank)
		}
		fl.Close()
		return
	} else {
		send_to_imm(libc.CString("Error with userWrite!"))
		return
	}
}
func fingerUser(ch *char_data, name *byte) {
	var (
		filename [40]byte
		uname    [100]byte
		email    [100]byte
		pass     [100]byte
		tmp1     [100]byte
		tmp2     [100]byte
		tmp3     [100]byte
		tmp4     [100]byte
		tmp5     [100]byte
		line     [256]byte
		total    int = 0
		rpp      int = 0
		rbank    int = 0
		file     *stdio.File
	)
	if get_filename(&filename[0], uint64(40), USER_FILE, name) == 0 {
		send_to_char(ch, libc.CString("That user doesn't exist.\r\n"))
		return
	} else if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&filename[0]), "r")
		return file
	}()) == nil {
		send_to_char(ch, libc.CString("That user is bugged! Report to Iovan.\r\n"))
		return
	}
	var count int = 0
	for int(file.IsEOF()) == 0 {
		get_line(file, &line[0])
		count += 1
		switch count {
		case 1:
			stdio.Sscanf(&line[0], "%s\n", &uname[0])
		case 2:
			stdio.Sscanf(&line[0], "%s\n", &email[0])
		case 3:
			stdio.Sscanf(&line[0], "%s\n", &pass[0])
		case 4:
			stdio.Sscanf(&line[0], "%d\n", &total)
		case 5:
			stdio.Sscanf(&line[0], "%d\n", &rpp)
		case 6:
			stdio.Sscanf(&line[0], "%s\n", &tmp1[0])
		case 7:
			stdio.Sscanf(&line[0], "%s\n", &tmp2[0])
		case 8:
			stdio.Sscanf(&line[0], "%s\n", &tmp3[0])
		case 9:
			stdio.Sscanf(&line[0], "%s\n", &tmp4[0])
		case 10:
			stdio.Sscanf(&line[0], "%s\n", &tmp5[0])
		case 13:
			stdio.Sscanf(&line[0], "%d\n", &rbank)
		}
	}
	file.Close()
	send_to_char(ch, libc.CString("@D[@gUsername   @D: @w%-30s@D]@n\r\n"), &uname[0])
	send_to_char(ch, libc.CString("@D[@gEmail      @D: @w%-30s@D]@n\r\n"), &email[0])
	if ch.Admlevel > 4 {
		send_to_char(ch, libc.CString("@D[@gPass       @D: @w%-30s@D]@n\r\n"), &pass[0])
	} else if ch.Admlevel > 0 {
		send_to_char(ch, libc.CString("@D[@gPass       @D: @w%-30s@D]@n\r\n"), "??????????")
	}
	send_to_char(ch, libc.CString("@D[@gTotal Slots@D: @w%-30d@D]@n\r\n"), total)
	send_to_char(ch, libc.CString("@D[@gRP Points  @D: @w%-30d@D]@n\r\n"), rpp)
	send_to_char(ch, libc.CString("@D[@gRP Bank    @D: @w%-30d@D]@n\r\n"), rbank)
	if ch.Admlevel > 0 {
		send_to_char(ch, libc.CString("@D[@gCh. Slot 1 @D: @w%-30s@D]@n\r\n"), &tmp1[0])
		send_to_char(ch, libc.CString("@D[@gCh. Slot 2 @D: @w%-30s@D]@n\r\n"), &tmp2[0])
		send_to_char(ch, libc.CString("@D[@gCh. Slot 3 @D: @w%-30s@D]@n\r\n"), &tmp3[0])
		send_to_char(ch, libc.CString("@D[@gCh. Slot 4 @D: @w%-30s@D]@n\r\n"), &tmp4[0])
		send_to_char(ch, libc.CString("@D[@gCh. Slot 5 @D: @w%-30s@D]@n\r\n"), &tmp5[0])
		send_to_char(ch, libc.CString("\n"))
		customRead(ch.Desc, 1, name)
	}
	return
}
func userRead(d *descriptor_data) {
	write_to_output(d, libc.CString("                 @RUser Menu@n\n"))
	write_to_output(d, libc.CString("@D=============================================@n\r\n"))
	write_to_output(d, libc.CString("@D|@gUser Account  @D: @w%-27s@D|@n\n"), d.User)
	write_to_output(d, libc.CString("@D|@gEmail Address @D: @w%-27s@D|@n\n"), d.Email)
	write_to_output(d, libc.CString("@D|@gMax Characters@D: @w%-27d@D|@n\n"), d.Total)
	write_to_output(d, libc.CString("@D|@gRP Points     @D: @w%-27d@D|@n\n"), d.Rpp)
	write_to_output(d, libc.CString("@D|@gRP Bank       @D: @w%-27d@D|@n\n"), d.Rbank)
	write_to_output(d, libc.CString("@D=============================================@n\r\n\r\n"))
	write_to_output(d, libc.CString("      @D[@y----@YSelect A Character Slot@y----@D]@n\n"))
	write_to_output(d, libc.CString("                @B(@W1@B) @C%s@n\n"), d.Tmp1)
	write_to_output(d, libc.CString("                @B(@W2@B) @C%s@n\n"), d.Tmp2)
	write_to_output(d, libc.CString("                @B(@W3@B) @C%s@n\n"), d.Tmp3)
	if d.Total > 3 {
		write_to_output(d, libc.CString("                @B(@W4@B) @C%s@n\n"), d.Tmp4)
	}
	if d.Total > 4 {
		write_to_output(d, libc.CString("                @B(@W5@B) @C%s@n\n"), d.Tmp5)
	}
	write_to_output(d, libc.CString("\n"))
	write_to_output(d, libc.CString("      @D[@y---- @YSelect Another Option @y----@D]@n\n"))
	write_to_output(d, libc.CString("                @B(@WB@B) @CBuy New C. Slot @D(@R15 RPP@D)@n\n"))
	write_to_output(d, libc.CString("                @B(@WC@B) @CUser's Customs\n"))
	write_to_output(d, libc.CString("                @B(@WD@B) @RDelete User@n\n"))
	write_to_output(d, libc.CString("                @B(@WE@B) @CEmail@n\n"))
	write_to_output(d, libc.CString("                @B(@WP@B) @CNew Password@n\n"))
	write_to_output(d, libc.CString("                @B(@WQ@B) @CQuit@n\n"))
	write_to_output(d, libc.CString("\r\nMake your choice: \n"))
}
func display_bonus_menu(ch *char_data, type_ int) {
	var (
		BonusCount int       = 26
		NegCount   int       = 52
		x          int       = 0
		y          int       = 0
		bonus      [52]*byte = [52]*byte{libc.CString("Thrifty     - -10% Shop Buy Cost and +10% Shop Sell Cost             @D[@G-2pts @D]"), libc.CString("Prodigy     - +25% Experience Gained Until Level 80                  @D[@G-5pts @D]"), libc.CString("Quick Study - Character auto-trains skills faster                    @D[@G-3pts @D]"), libc.CString("Die Hard    - Life Force's PL regen doubled, but cost is the same    @D[@G-6pts @D]"), libc.CString("Brawler     - Physical attacks do 20% more damage                    @D[@G-4pts @D]"), libc.CString("Destroyer   - Damaged Rooms act as regen rooms for you               @D[@G-3pts @D]"), libc.CString("Hard Worker - Physical rewards better + activity drains less stamina @D[@G-3pts @D]"), libc.CString("Healer      - Heal/First-aid/Vigor/Repair restore +10%               @D[@G-3pts @D]"), libc.CString("Loyal       - +20% Experience When Grouped As Follower               @D[@G-2pts @D]"), libc.CString("Brawny      - Strength gains +2 every 10 levels, Train STR + 75%     @D[@G-5pts @D]"), libc.CString("Scholarly   - Intelligence gains +2 every 10 levels, Train INT + 75% @D[@G-5pts @D]"), libc.CString("Sage        - Wisdom gains +2 every 10 levels, Train WIS + 75%       @D[@G-5pts @D]"), libc.CString("Agile       - Agility gains +2 every 10 levels, Train AGL + 75%      @D[@G-4pts @D]"), libc.CString("Quick       - Speed gains +2 every 10 levels, Train SPD + 75%        @D[@G-6pts @D]"), libc.CString("Sturdy      - Constitution +2 every 10 levels, Train CON + 75%       @D[@G-5pts @D]"), libc.CString("Thick Skin  - -20% Physical and -10% ki dmg received                 @D[@G-5pts @D]"), libc.CString("Recipe Int. - Food cooked by you lasts longer/heals better           @D[@G-2pts @D]"), libc.CString("Fireproof   - -50% Fire Dmg taken, -10% ki, immunity to burn         @D[@G-4pts @D]"), libc.CString("Powerhitter - 15% critical hits will be x4 instead of x2             @D[@G-4pts @D]"), libc.CString("Healthy     - 40% chance to recover from ill effects when sleeping   @D[@G-3pts @D]"), libc.CString("Insomniac   - Can't Sleep. Immune to yoikominminken and paralysis    @D[@G-2pts @D]"), libc.CString("Evasive     - +15% to dodge rolls                                    @D[@G-3pts @D]"), libc.CString("The Wall    - +20% chance to block                                   @D[@G-3pts @D]"), libc.CString("Accurate    - +20% chance to hit physical, +10% to hit with ki       @D[@G-3pts @D]"), libc.CString("Energy Leech- -2% ki damage received for every 5 character levels,   @D[@G-5pts @D]\n                  @cas long as you can take that ki to your charge pool.@D        "), libc.CString("Good Memory - +2 Skill Slots initially, +1 every 20 levels after     @D[@G-6pts @D]"), libc.CString("Soft Touch  - Half damage for all hit locations                      @D[@G+5pts @D]"), libc.CString("Late Sleeper- Can only wake automatically. 33% every hour if maxed   @D[@G+5pts @D]"), libc.CString("Impulse Shop- +25% shop costs                                        @D[@G+3pts @D]"), libc.CString("Sickly      - Suffer from harmful effects longer                     @D[@G+5pts @D]"), libc.CString("Punching Bag- -15% to dodge rolls                                    @D[@G+3pts @D]"), libc.CString("Pushover    - -20% block chance                                      @D[@G+3pts @D]"), libc.CString("Poor D. Perc- -20% chance to hit with physical, -10% with ki         @D[@G+3pts @D]"), libc.CString("Thin Skin   - +20% physical and +10% ki damage received              @D[@G+4pts @D]"), libc.CString("Fireprone   - +50% Fire Dmg taken, +10% ki, always burned            @D[@G+5pts @D]"), libc.CString("Energy Int. - +2% ki damage received for every 5 character levels,   @D[@G+6pts @D]\n                  @rif you have ki charged you have 10% chance to lose   \n                  it and to take 1/4th damage equal to it.@D                    "), libc.CString("Coward      - Can't Attack Enemy With 150% Your Powerlevel           @D[@G+6pts @D]"), libc.CString("Arrogant    - Cannot Suppress                                        @D[@G+1pt  @D]"), libc.CString("Unfocused   - Charge concentration randomly breaks                   @D[@G+3pts @D]"), libc.CString("Slacker     - Physical activity drains more stamina                  @D[@G+3pts @D]"), libc.CString("Slow Learner- Character auto-trains skills slower                    @D[@G+3pts @D]"), libc.CString("Masochistic - Defense Skills Cap At 75                               @D[@G+5pts @D]"), libc.CString("Mute        - Can't use IC speech related commands                   @D[@G+4pts @D]"), libc.CString("Wimp        - Strength is capped at 45                               @D[@G+6pts @D]"), libc.CString("Dull        - Intelligence is capped at 45                           @D[@G+6pts @D]"), libc.CString("Foolish     - Wisdom is capped at 45                                 @D[@G+6pts @D]"), libc.CString("Clumsy      - Agility is capped at 45                                @D[@G+3pts @D]"), libc.CString("Slow        - Speed is capped at 45                                  @D[@G+6pts @D]"), libc.CString("Frail       - Constitution capped at 45                              @D[@G+4pts @D]"), libc.CString("Sadistic    - Half Experience Gained For Quick Kills                 @D[@G+3pts @D]"), libc.CString("Loner       - Can't Group, +5% Train gains, +10% to physical gains   @D[@G+2pts @D]"), libc.CString("Bad Memory  - -5 Skill Slots                                         @D[@G+6pts @D]")}
	)
	if type_ == 0 {
		send_to_char(ch, libc.CString("\r\n@YBonus Trait SELECTION menu:\r\n@D---------------------------------------\r\n@n"))
		for x = 0; x < BonusCount; x++ {
			send_to_char(ch, libc.CString("@C%-2d@D)@c %s <@g%s@D>\n"), x+1, bonus[x], func() string {
				if (ch.Bonuses[x]) > 0 {
					return "X"
				}
				return " "
			}())
		}
		send_to_char(ch, libc.CString("\n"))
	}
	if type_ == 1 {
		y = BonusCount
		send_to_char(ch, libc.CString("@YNegative Trait SELECTION menu:\r\n@D---------------------------------------\r\n@n"))
		for y < NegCount {
			send_to_char(ch, libc.CString("@R%-2d@D)@r %s <@g%s@D>\n"), y-14, bonus[y], func() string {
				if (ch.Bonuses[y]) > 0 {
					return "X"
				}
				return " "
			}())
			y += 1
		}
	}
	if type_ == 0 {
		send_to_char(ch, libc.CString("\n@CN@D)@c Show Negatives@n\n"))
	} else {
		send_to_char(ch, libc.CString("\n@CB@D)@c Show Bonuses@n\n"))
	}
	send_to_char(ch, libc.CString("@CX@D)@c Exit Traits Section and complete your character@n\n"))
	send_to_char(ch, libc.CString("@D---------------------------------------\n[@WCurrent Points Pool@W: @y%d@D] [@WPTS From Neg@W: @y%d@D]@w\n"), ch.Ccpoints, ch.Negcount)
}
func parse_bonuses(arg *byte) int {
	var (
		value int = -1
		ident int = -1
	)
	switch *arg {
	case 'b':
		fallthrough
	case 'B':
		value = 53
	case 'n':
		fallthrough
	case 'N':
		value = 54
	case 'x':
		fallthrough
	case 'X':
		value = 55
	}
	if value < 52 {
		ident = libc.Atoi(libc.GoString(arg))
	}
	switch ident {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 7:
		fallthrough
	case 8:
		fallthrough
	case 9:
		fallthrough
	case 10:
		fallthrough
	case 11:
		fallthrough
	case 12:
		fallthrough
	case 13:
		fallthrough
	case 14:
		fallthrough
	case 15:
		fallthrough
	case 16:
		fallthrough
	case 17:
		fallthrough
	case 18:
		fallthrough
	case 19:
		fallthrough
	case 20:
		fallthrough
	case 21:
		fallthrough
	case 22:
		fallthrough
	case 23:
		fallthrough
	case 24:
		fallthrough
	case 25:
		fallthrough
	case 26:
		fallthrough
	case 27:
		fallthrough
	case 28:
		fallthrough
	case 29:
		fallthrough
	case 30:
		fallthrough
	case 31:
		fallthrough
	case 32:
		fallthrough
	case 33:
		fallthrough
	case 34:
		fallthrough
	case 35:
		fallthrough
	case 36:
		fallthrough
	case 37:
		fallthrough
	case 38:
		fallthrough
	case 39:
		fallthrough
	case 40:
		fallthrough
	case 41:
		fallthrough
	case 42:
		fallthrough
	case 43:
		fallthrough
	case 44:
		fallthrough
	case 45:
		fallthrough
	case 46:
		fallthrough
	case 47:
		fallthrough
	case 48:
		fallthrough
	case 49:
		fallthrough
	case 50:
		fallthrough
	case 51:
		value = ident - 1
	default:
		value = -1
	}
	return value
}

var list_bonus [52]*byte = [52]*byte{libc.CString("Thrifty"), libc.CString("Prodigy"), libc.CString("Quick Study"), libc.CString("Die Hard"), libc.CString("Brawler"), libc.CString("Destroyer"), libc.CString("Hard Worker"), libc.CString("Healer"), libc.CString("Loyal"), libc.CString("Brawny"), libc.CString("Scholarly"), libc.CString("Sage"), libc.CString("Agile"), libc.CString("Quick"), libc.CString("Sturdy"), libc.CString("Thick Skin"), libc.CString("Recipe Int"), libc.CString("Fireproof"), libc.CString("Powerhitter"), libc.CString("Healthy"), libc.CString("Insomniac"), libc.CString("Evasive"), libc.CString("The Wall"), libc.CString("Accurate"), libc.CString("Energy Leech"), libc.CString("Good Memory"), libc.CString("Soft Touch"), libc.CString("Late Sleeper"), libc.CString("Impulse Shop"), libc.CString("Sickly"), libc.CString("Punching Bag"), libc.CString("Pushover"), libc.CString("Poor Depth Perception"), libc.CString("Thin Skin"), libc.CString("Fireprone"), libc.CString("Energy Intollerant"), libc.CString("Coward"), libc.CString("Arrogant"), libc.CString("Unfocused"), libc.CString("Slacker"), libc.CString("Slow Learner"), libc.CString("Masochistic"), libc.CString("Mute"), libc.CString("Wimp"), libc.CString("Dull"), libc.CString("Foolish"), libc.CString("Clumsy"), libc.CString("Slow"), libc.CString("Frail"), libc.CString("Sadistic"), libc.CString("Loner"), libc.CString("Bad Memory")}
var list_bonus_cost [52]int = [52]int{-2, -5, -3, -6, -4, -3, -3, -3, -2, -5, -5, -5, -4, -6, -5, -5, -2, -4, -4, -3, -2, -3, -3, -4, -5, -6, 5, 5, 3, 5, 3, 3, 4, 4, 5, 6, 6, 1, 3, 3, 3, 5, 4, 6, 6, 6, 3, 6, 4, 3, 2, 6}

func opp_bonus(ch *char_data, value int, type_ int) int {
	var give int = TRUE
	switch value {
	case 0:
		if (ch.Bonuses[BONUS_IMPULSE]) != 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[BONUS_IMPULSE], list_bonus[value])
			give = FALSE
		}
	case 2:
		if (ch.Bonuses[40]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[40], list_bonus[value])
			give = FALSE
		}
	case 3:
		if int(ch.Race) == RACE_ANDROID {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("You can not take %s as an android!@n\r\n"), list_bonus[3])
			give = FALSE
		}
	case 6:
		if (ch.Bonuses[39]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[39], list_bonus[value])
			give = FALSE
		}
	case 8:
		if (ch.Bonuses[50]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[50], list_bonus[value])
			give = FALSE
		}
	case 9:
		if (ch.Bonuses[43]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[43], list_bonus[value])
			give = FALSE
		}
	case 10:
		if (ch.Bonuses[44]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[44], list_bonus[value])
			give = FALSE
		}
	case 11:
		if (ch.Bonuses[45]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[45], list_bonus[value])
			give = FALSE
		}
	case 12:
		if (ch.Bonuses[46]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[46], list_bonus[value])
			give = FALSE
		}
	case 13:
		if (ch.Bonuses[47]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[47], list_bonus[value])
			give = FALSE
		}
	case 14:
		if (ch.Bonuses[48]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[48], list_bonus[value])
			give = FALSE
		}
	case 15:
		if (ch.Bonuses[33]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[33], list_bonus[value])
			give = FALSE
		}
	case 16:
		if int(ch.Race) == RACE_ANDROID {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("You are an android and can not suppress anyway.\n\n"))
			give = FALSE
		}
	case 17:
		if int(ch.Race) == RACE_DEMON {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("As a demon you are already fireproof.\r\n"))
			give = FALSE
		} else if (ch.Bonuses[34]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[34], list_bonus[value])
			give = FALSE
		}
	case 18:
		if (ch.Bonuses[26]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[26], list_bonus[value])
			give = FALSE
		}
	case 19:
		if (ch.Bonuses[29]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[29], list_bonus[value])
			give = FALSE
		}
	case 20:
		if (ch.Bonuses[27]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[27], list_bonus[value])
			give = FALSE
		} else if int(ch.Race) == RACE_ANDROID {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("You can not take %s as an android!@n\r\n"), list_bonus[value])
			give = FALSE
		}
	case 21:
		if (ch.Bonuses[30]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[30], list_bonus[value])
			give = FALSE
		}
	case 22:
		if (ch.Bonuses[31]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[31], list_bonus[value])
			give = FALSE
		}
	case 23:
		if (ch.Bonuses[32]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[32], list_bonus[value])
			give = FALSE
		}
	case 24:
		if (ch.Bonuses[35]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[35], list_bonus[value])
			give = FALSE
		}
	case 25:
		if (ch.Bonuses[51]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[51], list_bonus[value])
			give = FALSE
		}
	case 26:
		if (ch.Bonuses[18]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[18], list_bonus[value])
			give = FALSE
		}
	case 27:
		if (ch.Bonuses[20]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[27], list_bonus[value])
			give = FALSE
		} else if int(ch.Race) == RACE_ANDROID {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("You can not take %s as an android!@n\r\n"), list_bonus[value])
			give = FALSE
		}
	case 29:
		if (ch.Bonuses[19]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[19], list_bonus[value])
			give = FALSE
		}
	case 30:
		if (ch.Bonuses[21]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[21], list_bonus[value])
			give = FALSE
		}
	case 31:
		if (ch.Bonuses[22]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[22], list_bonus[value])
			give = FALSE
		}
	case 32:
		if (ch.Bonuses[23]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[23], list_bonus[value])
			give = FALSE
		}
	case 33:
		if (ch.Bonuses[15]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[15], list_bonus[value])
			give = FALSE
		}
	case 34:
		if (ch.Bonuses[17]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[17], list_bonus[value])
			give = FALSE
		}
	case 35:
		if (ch.Bonuses[24]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[24], list_bonus[value])
			give = FALSE
		}
	case 39:
		if (ch.Bonuses[6]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[6], list_bonus[value])
			give = FALSE
		}
	case 40:
		if (ch.Bonuses[2]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[2], list_bonus[value])
			give = FALSE
		}
	case 43:
		if (ch.Bonuses[9]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[9], list_bonus[value])
			give = FALSE
		}
	case 44:
		if (ch.Bonuses[10]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[10], list_bonus[value])
			give = FALSE
		}
	case 45:
		if (ch.Bonuses[11]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[11], list_bonus[value])
			give = FALSE
		}
	case 46:
		if (ch.Bonuses[12]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[12], list_bonus[value])
			give = FALSE
		}
	case 47:
		if (ch.Bonuses[13]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[13], list_bonus[value])
			give = FALSE
		}
	case 48:
		if (ch.Bonuses[14]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[14], list_bonus[value])
			give = FALSE
		}
	case 50:
		if (ch.Bonuses[8]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[8], list_bonus[value])
			give = FALSE
		}
	case 51:
		if (ch.Bonuses[25]) > 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@R%s and %s are mutually exclusive.\n\n"), list_bonus[25], list_bonus[value])
			give = FALSE
		}
	}
	return give
}
func exchange_ccpoints(ch *char_data, value int) {
	var type_ int = 0
	if ch.Desc.Connected == CON_BONUS {
		type_ = 0
	} else {
		type_ = 1
	}
	if (ch.Bonuses[value]) > 0 && ch.Ccpoints-list_bonus_cost[value] < 0 {
		display_bonus_menu(ch, type_)
		send_to_char(ch, libc.CString("@RYou must unselect some bonus traits first.\r\n"))
		return
	} else if (ch.Bonuses[value]) > 0 && ch.Ccpoints-list_bonus_cost[value] >= 0 {
		ch.Ccpoints -= list_bonus_cost[value]
		if list_bonus_cost[value] > 0 {
			ch.Negcount -= list_bonus_cost[value]
		}
		ch.Bonuses[value] = 0
		display_bonus_menu(ch, type_)
		send_to_char(ch, libc.CString("@GYou cancel your selection of %s.\r\n"), list_bonus[value])
		return
	}
	if type_ == 0 {
		if value > 25 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@RYou are not in the negatives menu, enter B to switch.\r\n"))
			return
		} else if ch.Ccpoints+list_bonus_cost[value] < 0 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@RYou do not have enough points for %s.\r\n"), list_bonus[value])
			return
		} else if opp_bonus(ch, value, type_) == 0 {
			return
		} else if list_bonus_cost[value] < 0 {
			ch.Ccpoints += list_bonus_cost[value]
			ch.Bonuses[value] = 1
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@GYou select the bonus %s\r\n"), list_bonus[value])
			return
		}
	} else {
		if value < 26 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@RYou are not in the bonuses menu, enter B to switch.\r\n"))
			return
		}
		var x int
		var count int = 0
		for x = 14; x < 52; x++ {
			if (ch.Bonuses[x]) > 1 {
				count += list_bonus_cost[x]
			}
		}
		if list_bonus_cost[value]+count > 10 {
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@RYou can not have more than +10 points from negatives.\r\n"))
			return
		} else if opp_bonus(ch, value, type_) == 0 {
			return
		} else {
			ch.Ccpoints += list_bonus_cost[value]
			ch.Negcount += list_bonus_cost[value]
			ch.Bonuses[value] = 2
			display_bonus_menu(ch, type_)
			send_to_char(ch, libc.CString("@GYou select the negative %s\r\n"), list_bonus[value])
			return
		}
	}
}
func nanny(d *descriptor_data, arg *byte) {
	var (
		load_result   int = -1
		total         int
		rr            int
		moveon        int = FALSE
		penalty       int = FALSE
		player_i      int
		value         int
		roll          int = rand_number(1, 6)
		k             *descriptor_data
		olc_functions [16]struct {
			State int
			Func  func(*descriptor_data, *byte)
		} = [16]struct {
			State int
			Func  func(*descriptor_data, *byte)
		}{{State: CON_OEDIT, Func: oedit_parse}, {State: CON_IEDIT, Func: oedit_parse}, {State: CON_ZEDIT, Func: zedit_parse}, {State: CON_SEDIT, Func: sedit_parse}, {State: CON_MEDIT, Func: medit_parse}, {State: CON_REDIT, Func: redit_parse}, {State: CON_CEDIT, Func: cedit_parse}, {State: CON_AEDIT, Func: aedit_parse}, {State: CON_TRIGEDIT, Func: trigedit_parse}, {State: CON_ASSEDIT, Func: assedit_parse}, {State: CON_GEDIT, Func: gedit_parse}, {State: CON_LEVELUP, Func: levelup_parse}, {State: CON_HEDIT, Func: hedit_parse}, {State: CON_HSEDIT, Func: hsedit_parse}, {State: CON_POBJ, Func: pobj_edit_parse}, {State: -1, Func: nil}}
	)
	skip_spaces(&arg)
	if d.Character == nil {
		d.Character = new(char_data)
		clear_char(d.Character)
		d.Character.Player_specials = new(player_special_data)
		d.Character.Desc = d
	}
	for player_i = 0; olc_functions[player_i].State >= 0; player_i++ {
		if d.Connected == olc_functions[player_i].State {
			if context_help(d, arg) != 0 {
				return
			}
			(olc_functions[player_i].Func)(d, arg)
			return
		}
	}
	switch d.Connected {
	case CON_GET_NAME:
		if d.Character != nil {
			free_char(d.Character)
		}
		if d.Character == nil {
			d.Character = new(char_data)
			clear_char(d.Character)
			d.Character.Player_specials = new(player_special_data)
			d.Character.Desc = d
			SET_BIT_AR(d.Character.Player_specials.Pref[:], PRF_COLOR)
		}
		var buf [2048]byte
		var tmp_name [2048]byte
		if d.Writenew > 0 {
			if *arg == 0 {
				write_to_output(d, libc.CString("Enter name: "))
				return
			}
			d.Loadplay = libc.StrDup(arg)
		}
		if _parse_name(d.Loadplay, &tmp_name[0]) != 0 || libc.StrLen(&tmp_name[0]) < 2 || libc.StrLen(&tmp_name[0]) > MAX_NAME_LENGTH || Valid_Name(&tmp_name[0]) == 0 || fill_word(libc.StrCpy(&buf[0], &tmp_name[0])) != 0 || reserved_word(&buf[0]) != 0 {
			write_to_output(d, libc.CString("Invalid name, please try another.\r\nName: "))
			return
		}
		if d.Writenew > 0 && (func() int {
			player_i = load_char(&tmp_name[0], d.Character)
			return player_i
		}()) > -1 {
			userRead(d)
			write_to_output(d, libc.CString("That character is already taken.\r\n"))
			d.Writenew = 0
			d.Connected = CON_UMENU
			return
		} else {
			if (func() int {
				player_i = load_char(&tmp_name[0], d.Character)
				return player_i
			}()) > -1 {
				if d.Writenew > 0 {
					write_to_output(d, libc.CString("That character is already taken.\r\n"))
					userRead(d)
					d.Connected = CON_UMENU
					return
				}
				d.Character.Pfilepos = player_i
				if PLR_FLAGGED(d.Character, PLR_DELETED) {
					if (func() int {
						player_i = get_ptable_by_name(&tmp_name[0])
						return player_i
					}()) >= 0 {
						remove_player(player_i)
					}
					free_char(d.Character)
					if Valid_Name(&tmp_name[0]) == 0 {
						write_to_output(d, libc.CString("@YInvalid name@n, please try @Canother.@n\r\nName: "))
						return
					}
					d.Character = new(char_data)
					clear_char(d.Character)
					d.Character.Player_specials = new(player_special_data)
					d.Character.Desc = d
					d.Character.Name = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(&tmp_name[0])+1)[0]))
					libc.StrCpy(d.Character.Name, CAP(&tmp_name[0]))
					d.Character.Pfilepos = player_i
					SET_BIT_AR(d.Character.Player_specials.Pref[:], PRF_COLOR)
					display_races(d)
					d.Character.Rp = d.Rpp
					d.Connected = CON_QRACE
				} else {
					REMOVE_BIT_AR(d.Character.Act[:], PLR_WRITING)
					REMOVE_BIT_AR(d.Character.Act[:], PLR_MAILING)
					REMOVE_BIT_AR(d.Character.Act[:], PLR_CRYO)
					REMOVE_BIT_AR(d.Character.Affected_by[:], AFF_GROUP)
					if isbanned(&d.Host[0]) == BAN_SELECT && !PLR_FLAGGED(d.Character, PLR_SITEOK) {
						write_to_output(d, libc.CString("Sorry, this char has not been cleared for login from your site!\r\n"))
						mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Connection attempt for %s denied from %s"), GET_NAME(d.Character), &d.Host[0])
						d.Connected = CON_CLOSE
						return
					} else if GET_LEVEL(d.Character) < circle_restrict && circle_restrict < 101 {
						userRead(d)
						write_to_output(d, libc.CString("The game is temporarily restricted to at least %d level.\r\n"), circle_restrict)
						d.Connected = CON_UMENU
						mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Request for character load denied for %s [%s] (wizlock)"), GET_NAME(d.Character), &d.Host[0])
						return
					} else if perform_dupe_check(d) != 0 {
						return
					} else {
						d.Idle_tics = 0
						write_to_output(d, libc.CString("\r\n%s\r\n%s"), motd, config_info.Operation.MENU)
						d.Character.Rp = d.Rpp
						d.Connected = CON_MENU
					}
				}
			} else {
				if d.Writenew <= 0 {
					if libc.StrCaseCmp(d.Loadplay, d.Tmp1) == 0 {
						if d.Tmp1 != nil {
							libc.Free(unsafe.Pointer(d.Tmp1))
							d.Tmp1 = nil
						}
						d.Tmp1 = libc.CString("Empty")
					}
					if libc.StrCaseCmp(d.Loadplay, d.Tmp2) == 0 {
						if d.Tmp2 != nil {
							libc.Free(unsafe.Pointer(d.Tmp2))
							d.Tmp2 = nil
						}
						d.Tmp2 = libc.CString("Empty")
					}
					if libc.StrCaseCmp(d.Loadplay, d.Tmp3) == 0 {
						if d.Tmp3 != nil {
							libc.Free(unsafe.Pointer(d.Tmp3))
							d.Tmp3 = nil
						}
						d.Tmp3 = libc.CString("Empty")
					}
					if libc.StrCaseCmp(d.Loadplay, d.Tmp4) == 0 {
						if d.Tmp4 != nil {
							libc.Free(unsafe.Pointer(d.Tmp4))
							d.Tmp4 = nil
						}
						d.Tmp4 = libc.CString("Empty")
					}
					if libc.StrCaseCmp(d.Loadplay, d.Tmp5) == 0 {
						if d.Tmp5 != nil {
							libc.Free(unsafe.Pointer(d.Tmp5))
							d.Tmp5 = nil
						}
						d.Tmp5 = libc.CString("Empty")
					}
					userWrite(d, 0, 0, 0, libc.CString("index"))
					userRead(d)
					write_to_output(d, libc.CString("Character missing. Emptying slot.\r\n"))
					d.Connected = CON_UMENU
				} else {
					if Valid_Name(&tmp_name[0]) == 0 {
						write_to_output(d, libc.CString("Invalid name, please try another.\r\nName: "))
						return
					}
					d.Character.Name = (*byte)(unsafe.Pointer(&make([]int8, libc.StrLen(&tmp_name[0])+1)[0]))
					libc.StrCpy(d.Character.Name, CAP(&tmp_name[0]))
					display_races(d)
					d.Character.Rp = d.Rpp
					switch d.Writenew {
					case 1:
						if d.Tmp1 != nil {
							libc.Free(unsafe.Pointer(d.Tmp1))
							d.Tmp1 = nil
						}
						d.Tmp1 = libc.StrDup(d.Character.Name)
						userWrite(d, 0, 0, 0, libc.CString("index"))
					case 2:
						if d.Tmp2 != nil {
							libc.Free(unsafe.Pointer(d.Tmp2))
							d.Tmp2 = nil
						}
						d.Tmp2 = libc.StrDup(d.Character.Name)
						userWrite(d, 0, 0, 0, libc.CString("index"))
					case 3:
						if d.Tmp3 != nil {
							libc.Free(unsafe.Pointer(d.Tmp3))
							d.Tmp3 = nil
						}
						d.Tmp3 = libc.StrDup(d.Character.Name)
						userWrite(d, 0, 0, 0, libc.CString("index"))
					case 4:
						if d.Tmp4 != nil {
							libc.Free(unsafe.Pointer(d.Tmp4))
							d.Tmp4 = nil
						}
						d.Tmp4 = libc.StrDup(d.Character.Name)
						userWrite(d, 0, 0, 0, libc.CString("index"))
					case 5:
						if d.Tmp5 != nil {
							libc.Free(unsafe.Pointer(d.Tmp5))
							d.Tmp5 = nil
						}
						d.Tmp5 = libc.StrDup(d.Character.Name)
						userWrite(d, 0, 0, 0, libc.CString("index"))
					}
					d.Connected = CON_QRACE
				}
			}
		}
	case CON_NAME_CNFRM:
		if unicode.ToUpper(rune(*arg)) == 'Y' {
			if isbanned(&d.Host[0]) >= BAN_NEW {
				mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Request for new char %s denied from [%s] (siteban)"), d.Character.Name, &d.Host[0])
				write_to_output(d, libc.CString("Sorry, new characters are not allowed from your site!\r\n"))
				d.Connected = CON_CLOSE
				return
			}
			if circle_restrict != 0 {
				write_to_output(d, libc.CString("Sorry, new players can't be created at the moment.\r\n"))
				mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("Request for new char %s denied from [%s] (wizlock)"), d.Character.Name, &d.Host[0])
				d.Connected = CON_CLOSE
				return
			}
			d.Idle_tics = 0
			write_to_output(d, libc.CString("@MNew character.@n\r\nGive me a @gpassword@n for @C%s@n: "), d.Character.Name)
			d.Connected = CON_NEWPASSWD
		} else if *arg == 'n' || *arg == 'N' {
			write_to_output(d, libc.CString("Okay, what IS it, then? "))
			libc.Free(unsafe.Pointer(d.Character.Name))
			d.Character.Name = nil
			d.Connected = CON_GET_NAME
		} else {
			write_to_output(d, libc.CString("Please type Yes or No: "))
		}
	case CON_GET_USER:
		if isbanned(&d.Host[0]) != 0 {
			write_to_output(d, libc.CString("You have been banned. Have a nice day.\r\n"))
			d.Connected = CON_CLOSE
		}
		if *arg == 0 {
			write_to_output(d, libc.CString("Enter your desired username or the username you have already made.\nUsername?\r\n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("index")) == 0 {
			write_to_output(d, libc.CString("Try again, username?\r\n"))
			return
		} else if Valid_Name(CAP(arg)) == 0 {
			write_to_output(d, libc.CString("Invalid name. Username?\r\n"))
			return
		} else if libc.StrStr(arg, libc.CString(" ")) != nil {
			write_to_output(d, libc.CString("No spaces. Username?\r\n"))
			return
		} else if libc.StrStr(arg, libc.CString("1")) != nil || libc.StrStr(arg, libc.CString("2")) != nil || libc.StrStr(arg, libc.CString("3")) != nil || libc.StrStr(arg, libc.CString("4")) != nil || libc.StrStr(arg, libc.CString("5")) != nil || libc.StrStr(arg, libc.CString("6")) != nil || libc.StrStr(arg, libc.CString("7")) != nil || libc.StrStr(arg, libc.CString("8")) != nil || libc.StrStr(arg, libc.CString("9")) != nil || libc.StrStr(arg, libc.CString("0")) != nil {
			write_to_output(d, libc.CString("No numbers. Username?\r\n"))
			return
		} else if libc.StrStr(arg, libc.CString(".")) != nil || libc.StrStr(arg, libc.CString(",")) != nil || libc.StrStr(arg, libc.CString("!")) != nil || libc.StrStr(arg, libc.CString("?")) != nil || libc.StrStr(arg, libc.CString(";")) != nil || libc.StrStr(arg, libc.CString(":")) != nil || libc.StrStr(arg, libc.CString("'")) != nil {
			write_to_output(d, libc.CString("No punctuation. Username?\r\n"))
			return
		} else if libc.StrLen(arg) < 3 {
			write_to_output(d, libc.CString("Name must at least be 3 characters long, username?\r\n"))
			return
		} else if libc.StrLen(arg) > 10 {
			write_to_output(d, libc.CString("Name must be at most 10 characters long, username?\r\n"))
			return
		} else {
			if readUserIndex(arg) == 0 {
				if circle_restrict == 101 {
					write_to_output(d, libc.CString("The mud is locked for mortals at this time.\r\n"))
					write_to_output(d, libc.CString("Please check: www.advent-truth.com/forum for further details.\r\n"))
					send_to_imm(libc.CString("%s rejected by wizlock"), CAP(arg))
					d.Connected = CON_CLOSE
				} else {
					if d.User != nil {
						libc.Free(unsafe.Pointer(d.User))
						d.User = nil
					}
					d.User = libc.StrDup(arg)
					write_to_output(d, libc.CString("You want you user name to be, %s?\r\n"), d.User)
					write_to_output(d, libc.CString("Yes or no: \n"))
					d.Connected = CON_USER_CONF
				}
			} else {
				if circle_restrict == 101 {
					if d.User != nil {
						libc.Free(unsafe.Pointer(d.User))
						d.User = nil
					}
					d.User = libc.StrDup(arg)
					userLoad(d, d.User)
					if d.Level == 0 {
						write_to_output(d, libc.CString("The mud is locked for mortals at this time.\r\n"))
						send_to_imm(libc.CString("%s rejected by wizlock"), d.User)
						d.Connected = CON_CLOSE
					} else {
						write_to_output(d, libc.CString("Password: \r\n"))
						send_to_imm(libc.CString("Username, %s, logging in."), CAP(arg))
						d.Connected = CON_PASSWORD
					}
				} else {
					if d.User != nil {
						libc.Free(unsafe.Pointer(d.User))
						d.User = nil
					}
					userLoad(d, arg)
					write_to_output(d, libc.CString("Password: \r\n"))
					send_to_imm(libc.CString("Username, %s, logging in."), CAP(arg))
					d.Connected = CON_PASSWORD
				}
			}
		}
	case CON_USER_CONF:
		if *arg == 0 {
			write_to_output(d, libc.CString("You want your user name to be, %s?\r\n"), d.User)
			write_to_output(d, libc.CString("Yes or no: \n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("yes")) == 0 || libc.StrCaseCmp(arg, libc.CString("y")) == 0 {
			write_to_output(d, libc.CString("User Account, %s, created.\r\nEnter Email:\n"), d.User)
			write_to_output(d, libc.CString("Remember your email must be valid and matching the example given.\n"))
			write_to_output(d, libc.CString("[Example: iovan@@advent-truth.com]\n"))
			send_to_imm(libc.CString("Username, %s, creating."), CAP(d.User))
			d.Connected = CON_GET_EMAIL
		} else if libc.StrCaseCmp(arg, libc.CString("no")) == 0 || libc.StrCaseCmp(arg, libc.CString("n")) == 0 {
			if d.User != nil {
				libc.Free(unsafe.Pointer(d.User))
				d.User = nil
			}
			d.User = libc.CString("Empty")
			write_to_output(d, libc.CString("Enter Username: \n"))
			d.Connected = CON_GET_USER
		} else {
			write_to_output(d, libc.CString("You want you user name to be, %s?\r\n"), d.User)
			write_to_output(d, libc.CString("Yes or no: \n"))
			return
		}
	case CON_GET_EMAIL:
		if readUserIndex(d.User) != 0 && *arg == 0 {
			write_to_output(d, libc.CString("Email, or M for menu?\r\n"))
			return
		} else if *arg == 0 {
			write_to_output(d, libc.CString("Email?\r\n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("M")) == 0 && readUserIndex(d.User) != 0 {
			userRead(d)
			d.Connected = CON_UMENU
		} else if libc.StrStr(arg, libc.CString(".com")) == nil && libc.StrStr(arg, libc.CString(".net")) == nil && libc.StrStr(arg, libc.CString(".org")) == nil {
			write_to_output(d, libc.CString("Improper email format missing '.com' or '.net' or '.org'. Email?\r\n"))
			return
		} else if readUserIndex(d.User) != 0 {
			if libc.StrStr(arg, libc.CString("@")) != nil {
				search_replace(arg, libc.CString("@"), libc.CString("<AT>"))
			}
			if d.Email != nil {
				libc.Free(unsafe.Pointer(d.Email))
				d.Email = nil
			}
			d.Email = libc.StrDup(arg)
			userWrite(d, 0, 0, 0, libc.CString("index"))
			userRead(d)
			d.Connected = CON_UMENU
		} else {
			if libc.StrStr(arg, libc.CString("@")) != nil {
				search_replace(arg, libc.CString("@"), libc.CString("<AT>"))
			}
			write_to_output(d, libc.CString("Your email is: %s\n"), arg)
			if d.Email != nil {
				libc.Free(unsafe.Pointer(d.Email))
				d.Email = nil
			}
			d.Email = libc.StrDup(arg)
			write_to_output(d, libc.CString("Password: \r\n"))
			d.Connected = CON_NEWPASSWD
		}
	case CON_NEWPASSWD:
		if readUserIndex(d.User) != 0 && *arg == 0 {
			write_to_output(d, libc.CString("Password, or M for menu?\r\n"))
			return
		} else if *arg == 0 {
			write_to_output(d, libc.CString("Password?\r\n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("M")) == 0 && readUserIndex(d.User) != 0 {
			userRead(d)
			d.Connected = CON_UMENU
		} else if libc.StrLen(arg) > MAX_PWD_LENGTH && readUserIndex(d.User) != 0 {
			write_to_output(d, libc.CString("Password is too long. Password, or M for menu?\r\n"))
			return
		} else if libc.StrLen(arg) > MAX_PWD_LENGTH {
			write_to_output(d, libc.CString("Password is too long. Password?\r\n"))
			return
		} else if readUserIndex(d.User) != 0 {
			write_to_output(d, libc.CString("Your password and user account have been fully saved.\r\n"))
			if d.Pass != nil {
				libc.Free(unsafe.Pointer(d.Pass))
				d.Pass = nil
			}
			d.Pass = libc.StrDup(arg)
			userWrite(d, 0, 0, 0, libc.CString("index"))
			userRead(d)
			d.Connected = CON_UMENU
		} else {
			write_to_output(d, libc.CString("Your password and user account have been fully saved.\r\n"))
			if d.Pass != nil {
				libc.Free(unsafe.Pointer(d.Pass))
				d.Pass = nil
			}
			d.Pass = libc.StrDup(arg)
			userCreate(d)
			userRead(d)
			d.Connected = CON_UMENU
		}
	case CON_PASSWORD:
		if *arg == 0 {
			write_to_output(d, libc.CString("Enter Password or Return:\r\n"))
			write_to_output(d, libc.CString("(Return will ask for a different username)\r\n"))
			return
		}
		if libc.StrCaseCmp(libc.CString("return"), arg) == 0 || libc.StrCaseCmp(libc.CString("Return"), arg) == 0 {
			if d.User != nil {
				libc.Free(unsafe.Pointer(d.User))
				d.User = nil
			}
			d.User = libc.CString("Empty")
			write_to_output(d, libc.CString("Username?\r\n"))
			d.Connected = CON_GET_USER
		}
		if libc.StrCaseCmp(d.Pass, arg) == 0 {
			for k = descriptor_list; k != nil; k = k.Next {
				if k == d {
					continue
				}
				if k.User == nil || k.User == nil {
					continue
				}
				if d.User == nil || d.User == nil {
					continue
				}
				if libc.StrCaseCmp(k.User, d.User) == 0 {
					if k.Connected == CON_PLAYING {
						k.Connected = CON_DISCONNECT
						write_to_output(k, libc.CString("Your account has been usurped by someone who knows its password!@n"))
					} else {
						k.Connected = CON_CLOSE
						write_to_output(k, libc.CString("Your account has been usurped by someone who knows its password!@n"))
					}
				}
			}
			userRead(d)
			d.Connected = CON_UMENU
		} else {
			write_to_output(d, libc.CString("Password is wrong. Password or Return?\r\n"))
			write_to_output(d, libc.CString("(Return will ask for a different username)\r\n"))
			send_to_imm(libc.CString("Username, %s, password failure!"), CAP(d.User))
			basic_mud_log(libc.CString("%s BAD PASSWORD: %s"), CAP(d.User), arg)
			return
		}
	case CON_UMENU:
		if *arg == 0 {
			userRead(d)
			return
		}
		if libc.StrCaseCmp(arg, libc.CString("Q")) == 0 {
			write_to_output(d, libc.CString("Thanks for visiting!\n"))
			d.Connected = CON_CLOSE
		} else if libc.StrCaseCmp(arg, libc.CString("P")) == 0 {
			write_to_output(d, libc.CString("Enter New Password, or M for menu::\n"))
			d.Connected = CON_NEWPASSWD
		} else if libc.StrCaseCmp(arg, libc.CString("C")) == 0 {
			write_to_output(d, libc.CString("\n"))
			customRead(d, 0, nil)
			write_to_output(d, libc.CString("\r\n@n--Press Enter--\n@n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("E")) == 0 {
			write_to_output(d, libc.CString("Enter New Email, or M for menu:\n"))
			d.Connected = CON_GET_EMAIL
		} else if libc.StrCaseCmp(arg, libc.CString("D")) == 0 {
			write_to_output(d, libc.CString("Are you sure you want to delete your user file and all its characters? Yes or no:\n"))
			d.Connected = CON_DELCNF1
		} else if libc.StrCaseCmp(arg, libc.CString("B")) == 0 {
			if d.Total == 3 && d.Rpp >= 15 {
				d.Rpp -= 15
				d.Total = 4
				userWrite(d, 0, 0, 0, libc.CString("index"))
				userRead(d)
				write_to_output(d, libc.CString("New slot purchased, -15 RPP.\n"))
				return
			} else if d.Total == 3 && d.Rpp < 15 {
				userRead(d)
				write_to_output(d, libc.CString("You need at least 15 RPP to purchase a new character slot.\n"))
				return
			} else if d.Total == 4 && d.Rpp >= 15 {
				d.Rpp -= 15
				d.Total = 5
				userWrite(d, 0, 0, 0, libc.CString("index"))
				userRead(d)
				write_to_output(d, libc.CString("New slot purchased, -15 RPP.\n"))
				return
			} else if d.Total == 4 && d.Rpp < 15 {
				userRead(d)
				write_to_output(d, libc.CString("You need at least 15 RPP to purchase a new character slot.\n"))
				return
			} else {
				userRead(d)
				write_to_output(d, libc.CString("You can't have more than 5 character slots.\n"))
				return
			}
		} else {
			switch libc.Atoi(libc.GoString(arg)) {
			case 1:
				if libc.StrCaseCmp(d.Tmp1, libc.CString("Empty")) == 0 {
					write_to_output(d, libc.CString("Enter New Character Name: \n"))
					d.Writenew = 1
					d.Connected = CON_GET_NAME
				} else {
					if lockRead(d.Tmp1) != 0 && d.Level <= 0 {
						write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
						return
					} else {
						write_to_output(d, libc.CString("Loading character, %s.\r\n***Press Return***\n"), d.Tmp1)
						d.Loadplay = d.Tmp1
						d.Writenew = 0
						d.Connected = CON_GET_NAME
					}
				}
			case 2:
				if libc.StrCaseCmp(d.Tmp2, libc.CString("Empty")) == 0 {
					write_to_output(d, libc.CString("Enter New Character Name: \n"))
					d.Writenew = 2
					d.Connected = CON_GET_NAME
				} else {
					if lockRead(d.Tmp2) != 0 && d.Level <= 0 {
						write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
						return
					} else {
						write_to_output(d, libc.CString("Loading character, %s.\r\n***Press Return***\n"), d.Tmp2)
						d.Loadplay = d.Tmp2
						d.Writenew = 0
						d.Connected = CON_GET_NAME
					}
				}
			case 3:
				if libc.StrCaseCmp(d.Tmp3, libc.CString("Empty")) == 0 {
					write_to_output(d, libc.CString("Enter New Character Name: \n"))
					d.Writenew = 3
					d.Connected = CON_GET_NAME
				} else {
					if lockRead(d.Tmp3) != 0 && d.Level <= 0 {
						write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
						return
					} else {
						write_to_output(d, libc.CString("Loading character, %s.\r\n***Press Return***\n"), d.Tmp3)
						d.Loadplay = d.Tmp3
						d.Writenew = 0
						d.Connected = CON_GET_NAME
					}
				}
			case 4:
				if d.Total <= 3 {
					userRead(d)
					write_to_output(d, libc.CString("You only have %d character slots avaialable!\r\n"), d.Total)
					return
				}
				if libc.StrCaseCmp(d.Tmp4, libc.CString("Empty")) == 0 {
					write_to_output(d, libc.CString("Enter New Character Name: \n"))
					d.Writenew = 4
					d.Connected = CON_GET_NAME
				} else {
					if lockRead(d.Tmp4) != 0 && d.Level <= 0 {
						write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
						return
					} else {
						write_to_output(d, libc.CString("Loading character, %s.\r\n***Press Return***\n"), d.Tmp4)
						d.Loadplay = d.Tmp4
						d.Writenew = 0
						d.Connected = CON_GET_NAME
					}
				}
			case 5:
				if d.Total <= 4 {
					userRead(d)
					write_to_output(d, libc.CString("You only have %d character slots avaialable!\r\n"), d.Total)
					return
				}
				if libc.StrCaseCmp(d.Tmp5, libc.CString("Empty")) == 0 {
					write_to_output(d, libc.CString("Enter New Character Name: \n"))
					d.Writenew = 5
					d.Connected = CON_GET_NAME
				} else {
					if lockRead(d.Tmp5) != 0 && d.Level <= 0 {
						write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
						return
					} else {
						write_to_output(d, libc.CString("Loading character, %s.\r\n***Press Return***\n"), d.Tmp5)
						d.Loadplay = d.Tmp5
						d.Writenew = 0
						d.Connected = CON_GET_NAME
					}
				}
			default:
				userRead(d)
				write_to_output(d, libc.CString("That is not an option.\n"))
				return
			}
		}
	case CON_QSEX:
		if int(d.Character.Race) != RACE_NAMEK {
			switch *arg {
			case 'm':
				fallthrough
			case 'M':
				d.Character.Sex = SEX_MALE
			case 'f':
				fallthrough
			case 'F':
				d.Character.Sex = SEX_FEMALE
			case 'n':
				fallthrough
			case 'N':
				d.Character.Sex = SEX_NEUTRAL
			default:
				write_to_output(d, libc.CString("That is not a sex..\r\nWhat IS your sex? "))
				return
			}
		}
		if int(d.Character.Race) == RACE_HUMAN || int(d.Character.Race) == RACE_SAIYAN || int(d.Character.Race) == RACE_KONATSU || int(d.Character.Race) == RACE_MUTANT || int(d.Character.Race) == RACE_ANDROID || int(d.Character.Race) == RACE_KAI || int(d.Character.Race) == RACE_HALFBREED || int(d.Character.Race) == RACE_TRUFFLE || int(d.Character.Race) == RACE_HOSHIJIN && int(d.Character.Sex) == SEX_FEMALE {
			write_to_output(d, libc.CString("@YHair Length SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Bald  @B2@W)@C Short  @B3@W)@C Medium\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Long  @B5@W)@C Really Long@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_HAIRL
		} else if int(d.Character.Race) == RACE_DEMON || int(d.Character.Race) == RACE_ICER {
			write_to_output(d, libc.CString("@YHorn Length SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C None  @B2@W)@C Short  @B3@W)@C Medium\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Long  @B5@W)@C Really Long@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_HAIRL
		} else if int(d.Character.Race) == RACE_MAJIN {
			write_to_output(d, libc.CString("@YForelock Length SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Tiny  @B2@W)@C Short  @B3@W)@C Medium\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Long  @B5@W)@C Really Long@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_HAIRL
		} else if int(d.Character.Race) == RACE_NAMEK || int(d.Character.Race) == RACE_ARLIAN {
			write_to_output(d, libc.CString("@YAntenae Length SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Tiny  @B2@W)@C Short  @B3@W)@C Medium\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Long  @B5@W)@C Really Long@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_HAIRL
		} else {
			write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
			write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
			write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_SKIN
		}
	case CON_QRACE:
		switch *arg {
		case 'r':
			fallthrough
		case 'R':
			for load_result == -1 {
				rr = rand_number(1, NUM_RACES)
				load_result = parse_race(d.Character, rr)
			}
		case 't':
			fallthrough
		case 'T':
			display_races_help(d)
			d.Connected = CON_RACE_HELP
			return
		}
		if load_result == -1 {
			load_result = parse_race(d.Character, libc.Atoi(libc.GoString(arg)))
		}
		if load_result == -1 && (libc.Atoi(libc.GoString(arg)) != 2 && libc.Atoi(libc.GoString(arg)) != 12 && libc.Atoi(libc.GoString(arg)) != 9) {
			write_to_output(d, libc.CString("\r\nThat's not a race.\r\nRace: "))
			return
		} else if load_result == -1 && libc.Atoi(libc.GoString(arg)) == 2 {
			write_to_output(d, libc.CString("\r\nThat race is restricted, You need 60 RPP to unlock it.\r\nRace: "))
			return
		} else if load_result == -1 && libc.Atoi(libc.GoString(arg)) == 12 {
			write_to_output(d, libc.CString("\r\nThat race is restricted, You need 55 RPP to unlock it.\r\nRace: "))
			return
		} else if load_result == -1 && libc.Atoi(libc.GoString(arg)) == 9 {
			write_to_output(d, libc.CString("\r\nThat race is restricted, You need 30 RPP to unlock it.\r\nRace: "))
			return
		} else if load_result == -1 && libc.Atoi(libc.GoString(arg)) == 15 {
			write_to_output(d, libc.CString("\r\nThat race is restricted, You need 20 RPP to unlock it.\r\nRace: "))
			return
		} else {
			d.Character.Race = int8(load_result)
		}
		if int(d.Character.Race) == RACE_HALFBREED {
			write_to_output(d, libc.CString("@YWhat race do you prefer to by identified with?\r\n"))
			write_to_output(d, libc.CString("@cThis controls how others first view you and whether you start with\na tail and how fast it regrows when missing.\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Human\n@B2@W)@C Saiyan\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_RACIAL
		} else if int(d.Character.Race) == RACE_ANDROID {
			write_to_output(d, libc.CString("@YWhat do you want to be identified as at first glance?\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Android\n@B2@W)@C Human\n@B3@W)@C Robotic Humanoid\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_RACIAL
		} else if int(d.Character.Race) == RACE_NAMEK {
			d.Character.Sex = SEX_NEUTRAL
			d.Connected = CON_QSEX
		} else {
			write_to_output(d, libc.CString("\r\n@wWhat is your sex @W(@BM@W/@MF@W/@GN@W)@w?@n"))
			d.Connected = CON_QSEX
		}
	case CON_RACIAL:
		switch *arg {
		case '1':
			d.Character.Player_specials.Racial_pref = 1
		case '2':
			d.Character.Player_specials.Racial_pref = 2
		case '3':
			if int(d.Character.Race) == RACE_HALFBREED {
				write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
				return
			} else {
				d.Character.Player_specials.Racial_pref = 3
			}
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@wWhat is your sex @W(@BM@W/@MF@W/@GN@W)@w?@n"))
		d.Connected = CON_QSEX
	case CON_RACE_HELP:
		if *arg == 't' || *arg == 'T' {
			display_races(d)
			d.Connected = CON_QRACE
			return
		}
		if unicode.IsDigit(rune(*arg)) {
			player_i = libc.Atoi(libc.GoString(arg))
			if player_i > NUM_RACES || player_i < 1 {
				write_to_output(d, libc.CString("\r\nThat's not a race.\r\nHelp on Race #: "))
				break
			}
			player_i -= 1
			if race_ok_gender[int(d.Character.Sex)][player_i] {
				show_help(d, race_names[player_i])
			} else {
				write_to_output(d, libc.CString("\r\nThat's not a race.\r\nHelp on Race #: "))
			}
		} else {
			display_races_help(d)
		}
		d.Connected = CON_RACE_HELP
	case CON_CLASS_HELP:
		if *arg == 't' || *arg == 'T' {
			display_classes(d)
			d.Connected = CON_QCLASS
			return
		}
		if unicode.IsDigit(rune(*arg)) {
			player_i = libc.Atoi(libc.GoString(arg))
			if player_i > 14 || player_i < 1 {
				write_to_output(d, libc.CString("\r\nThat's not a sensei.\r\nHelp on Sensei #: "))
				break
			}
			player_i -= 1
			if class_ok_race[int(d.Character.Sex)][player_i] {
				show_help(d, class_names[player_i])
			} else {
				write_to_output(d, libc.CString("\r\nThat's not a sensei.\r\nHelp on Sensei #: "))
			}
		} else {
			display_classes_help(d)
		}
		d.Connected = CON_CLASS_HELP
	case CON_HAIRL:
		if int(d.Character.Race) == RACE_HUMAN || int(d.Character.Race) == RACE_SAIYAN || int(d.Character.Race) == RACE_KONATSU || int(d.Character.Race) == RACE_MUTANT || int(d.Character.Race) == RACE_ANDROID || int(d.Character.Race) == RACE_KAI || int(d.Character.Race) == RACE_HALFBREED || int(d.Character.Race) == RACE_TRUFFLE || int(d.Character.Race) == RACE_HOSHIJIN && int(d.Character.Sex) == SEX_FEMALE {
			switch *arg {
			case '1':
				d.Character.Hairl = HAIRL_BALD
				d.Character.Hairc = HAIRC_NONE
				d.Character.Hairs = HAIRS_NONE
			case '2':
				d.Character.Hairl = HAIRL_SHORT
			case '3':
				d.Character.Hairl = HAIRL_MEDIUM
			case '4':
				d.Character.Hairl = HAIRL_LONG
			case '5':
				d.Character.Hairl = HAIRL_RLONG
			default:
				write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
				return
			}
			if int(d.Character.Hairl) == HAIRL_BALD {
				write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
				write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
				write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
				write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
				write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
				write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
				write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
				d.Connected = CON_SKIN
			} else {
				write_to_output(d, libc.CString("@YHair color SELECTION menu:\r\n"))
				write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
				write_to_output(d, libc.CString("@B1@W)@C Black  @B2@W)@C Brown  @B3@W)@C Blonde\r\n"))
				write_to_output(d, libc.CString("@B4@W)@C Grey   @B5@W)@C Red    @B6@W)@C Orange@n\r\n"))
				write_to_output(d, libc.CString("@B7@W)@C Green  @B8@W)@C Blue   @B9@W)@C Pink\r\n"))
				write_to_output(d, libc.CString("@BA@W)@C Purple @BB@W)@C Silver @BC@W)@C Crimson@n\r\n"))
				write_to_output(d, libc.CString("@BD@W)@C White@n\r\n"))
				write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
				d.Connected = CON_HAIRC
			}
		} else {
			if int(d.Character.Race) == RACE_DEMON || int(d.Character.Race) == RACE_ICER {
				switch *arg {
				case '1':
					d.Character.Hairl = HAIRL_BALD
				case '2':
					d.Character.Hairl = HAIRL_SHORT
				case '3':
					d.Character.Hairl = HAIRL_MEDIUM
				case '4':
					d.Character.Hairl = HAIRL_LONG
				case '5':
					d.Character.Hairl = HAIRL_RLONG
				default:
					write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
					return
				}
				d.Character.Hairc = HAIRC_NONE
				d.Character.Hairs = HAIRS_NONE
				write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
				write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
				write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
				write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
				write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
				write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
				write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
				d.Connected = CON_SKIN
			}
			if int(d.Character.Race) == RACE_MAJIN || int(d.Character.Race) == RACE_NAMEK || int(d.Character.Race) == RACE_ARLIAN {
				switch *arg {
				case '1':
					d.Character.Hairl = HAIRL_BALD
				case '2':
					d.Character.Hairl = HAIRL_SHORT
				case '3':
					d.Character.Hairl = HAIRL_MEDIUM
				case '4':
					d.Character.Hairl = HAIRL_LONG
				case '5':
					d.Character.Hairl = HAIRL_RLONG
				default:
					write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
					return
				}
				if int(d.Character.Race) == RACE_ARLIAN && int(d.Character.Sex) == SEX_FEMALE {
					write_to_output(d, libc.CString("@YWing color SELECTION menu:\r\n"))
					write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
					write_to_output(d, libc.CString("@B1@W)@C Black  @B2@W)@C Brown  @B3@W)@C Blonde\r\n"))
					write_to_output(d, libc.CString("@B4@W)@C Grey   @B5@W)@C Red    @B6@W)@C Orange@n\r\n"))
					write_to_output(d, libc.CString("@B7@W)@C Green  @B8@W)@C Blue   @B9@W)@C Pink\r\n"))
					write_to_output(d, libc.CString("@BA@W)@C Purple @BB@W)@C Silver @BC@W)@C Crimson@n\r\n"))
					write_to_output(d, libc.CString("@BD@W)@C White@n\r\n"))
					write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
					d.Connected = CON_HAIRC
				} else {
					d.Character.Hairc = HAIRC_NONE
					d.Character.Hairs = HAIRS_NONE
					write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
					write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
					write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
					write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
					write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
					write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
					write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
					d.Connected = CON_SKIN
				}
			} else {
				d.Character.Hairc = HAIRC_NONE
				d.Character.Hairs = HAIRS_NONE
				d.Connected = CON_SKIN
			}
		}
	case CON_HAIRC:
		switch *arg {
		case '1':
			d.Character.Hairc = HAIRC_BLACK
		case '2':
			d.Character.Hairc = HAIRC_BROWN
		case '3':
			d.Character.Hairc = HAIRC_BLONDE
		case '4':
			d.Character.Hairc = HAIRC_GREY
		case '5':
			d.Character.Hairc = HAIRC_RED
		case '6':
			d.Character.Hairc = HAIRC_ORANGE
		case '7':
			d.Character.Hairc = HAIRC_GREEN
		case '8':
			d.Character.Hairc = HAIRC_BLUE
		case '9':
			d.Character.Hairc = HAIRC_PINK
		case 'A':
			d.Character.Hairc = HAIRC_PURPLE
		case 'a':
			d.Character.Hairc = HAIRC_PURPLE
		case 'B':
			d.Character.Hairc = HAIRC_SILVER
		case 'b':
			d.Character.Hairc = HAIRC_SILVER
		case 'C':
			d.Character.Hairc = HAIRC_CRIMSON
		case 'c':
			d.Character.Hairc = HAIRC_CRIMSON
		case 'D':
			d.Character.Hairc = HAIRC_WHITE
		case 'd':
			d.Character.Hairc = HAIRC_WHITE
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		if int(d.Character.Race) == RACE_ARLIAN {
			d.Character.Hairs = HAIRS_NONE
			write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
			write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
			write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_SKIN
		} else {
			write_to_output(d, libc.CString("@YHair style SELECTION menu:\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Plain     @B2@W)@C Mohawk    @B3@W)@C Spiky\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Curly     @B5@W)@C Uneven    @B6@W)@C Ponytail@n\r\n"))
			write_to_output(d, libc.CString("@B7@W)@C Afro      @B8@W)@C Fade      @B9@W)@C Crew Cut\r\n"))
			write_to_output(d, libc.CString("@BA@W)@C Feathered @BB@W)@C Dred Locks@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_HAIRS
		}
	case CON_HAIRS:
		switch *arg {
		case '1':
			d.Character.Hairs = HAIRS_PLAIN
		case '2':
			d.Character.Hairs = HAIRS_MOHAWK
		case '3':
			d.Character.Hairs = HAIRS_SPIKY
		case '4':
			d.Character.Hairs = HAIRS_CURLY
		case '5':
			d.Character.Hairs = HAIRS_UNEVEN
		case '6':
			d.Character.Hairs = HAIRS_PONYTAIL
		case '7':
			d.Character.Hairs = HAIRS_AFRO
		case '8':
			d.Character.Hairs = HAIRS_FADE
		case '9':
			d.Character.Hairs = HAIRS_CREW
		case 'A':
			d.Character.Hairs = HAIRS_FEATHERED
		case 'a':
			d.Character.Hairs = HAIRS_FEATHERED
		case 'B':
			d.Character.Hairs = HAIRS_DRED
		case 'b':
			d.Character.Hairs = HAIRS_DRED
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("@YSkin color SELECTION menu:\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Black  @B3@W)@C Green\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C Orange @B5@W)@C Yellow @B6@W)@C Red@n\r\n"))
		write_to_output(d, libc.CString("@B7@W)@C Grey   @B8@W)@C Blue   @B9@W)@C Aqua\r\n"))
		write_to_output(d, libc.CString("@BA@W)@C Pink   @BB@W)@C Purple @BC@W)@C Tan@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_SKIN
	case CON_SKIN:
		switch *arg {
		case '1':
			d.Character.Skin = SKIN_WHITE
		case '2':
			d.Character.Skin = SKIN_BLACK
		case '3':
			d.Character.Skin = SKIN_GREEN
		case '4':
			d.Character.Skin = SKIN_ORANGE
		case '5':
			d.Character.Skin = SKIN_YELLOW
		case '6':
			d.Character.Skin = SKIN_RED
		case '7':
			d.Character.Skin = SKIN_GREY
		case '8':
			d.Character.Skin = SKIN_BLUE
		case '9':
			d.Character.Skin = SKIN_AQUA
		case 'A':
			d.Character.Skin = SKIN_PINK
		case 'a':
			d.Character.Skin = SKIN_PINK
		case 'B':
			d.Character.Skin = SKIN_PURPLE
		case 'b':
			d.Character.Skin = SKIN_PURPLE
		case 'C':
			d.Character.Skin = SKIN_TAN
		case 'c':
			d.Character.Skin = SKIN_TAN
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("@YEye color SELECTION menu:\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C Blue  @B2@W)@C Black  @B3@W)@C Green\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C Brown @B5@W)@C Red    @B6@W)@C Aqua@n\r\n"))
		write_to_output(d, libc.CString("@B7@W)@C Pink  @B8@W)@C Purple @B9@W)@C Crimson\r\n"))
		write_to_output(d, libc.CString("@BA@W)@C Gold  @BB@W)@C Amber  @BC@W)@C Emerald@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_EYE
	case CON_EYE:
		switch *arg {
		case '1':
			d.Character.Eye = EYE_BLUE
		case '2':
			d.Character.Eye = EYE_BLACK
		case '3':
			d.Character.Eye = EYE_GREEN
		case '4':
			d.Character.Eye = EYE_BROWN
		case '5':
			d.Character.Eye = EYE_RED
		case '6':
			d.Character.Eye = EYE_AQUA
		case '7':
			d.Character.Eye = EYE_PINK
		case '8':
			d.Character.Eye = EYE_PURPLE
		case '9':
			d.Character.Eye = EYE_CRIMSON
		case 'A':
			d.Character.Eye = EYE_GOLD
		case 'a':
			d.Character.Eye = EYE_GOLD
		case 'B':
			d.Character.Eye = EYE_AMBER
		case 'b':
			d.Character.Eye = EYE_AMBER
		case 'C':
			d.Character.Eye = EYE_EMERALD
		case 'c':
			d.Character.Eye = EYE_EMERALD
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("@YWhat do you want to be your most distinguishing feature:\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C My Eyes@n\r\n"))
		if int(d.Character.Race) == RACE_MAJIN {
			write_to_output(d, libc.CString("@B2@W)@C My Forelock@n\r\n"))
		} else if int(d.Character.Race) == RACE_NAMEK || int(d.Character.Race) == RACE_ARLIAN {
			write_to_output(d, libc.CString("@B2@W)@C My Antennae@n\r\n"))
		} else if int(d.Character.Race) == RACE_ICER || int(d.Character.Race) == RACE_DEMON {
			write_to_output(d, libc.CString("@B2@W)@C My Horns@n\r\n"))
		} else if int(d.Character.Hairl) == HAIRL_BALD {
			write_to_output(d, libc.CString("@B2@W)@C My Baldness@n\r\n"))
		} else {
			write_to_output(d, libc.CString("@B2@W)@C My Hair@n\r\n"))
		}
		write_to_output(d, libc.CString("@B3@W)@C My Skin@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C My Height@n\r\n"))
		write_to_output(d, libc.CString("@B5@W)@C My Weight@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_DISTFEA
	case CON_DISTFEA:
		switch *arg {
		case '1':
			d.Character.Distfea = 0
		case '2':
			d.Character.Distfea = 1
		case '3':
			d.Character.Distfea = 2
		case '4':
			d.Character.Distfea = 3
		case '5':
			d.Character.Distfea = 4
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("@YWhat Height/Weight Range do you prefer:\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
			write_to_output(d, libc.CString("@B1@W)@C 100-120cm, 25-30kg@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C 120-140cm, 30-35kg@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C 140-160cm, 35-45kg@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C 160-180cm, 45-60kg@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C 180-200cm, 60-80kg@n\r\n"))
			write_to_output(d, libc.CString("@B6@W)@C 200-220cm, 80-100kg@n\r\n"))
		} else if int(d.Character.Race) == RACE_ICER {
			write_to_output(d, libc.CString("@B1@W)@C 100-120cm, 25-30kg@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C 120-140cm, 30-35kg@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C 140-160cm, 35-45kg@n\r\n"))
		} else {
			write_to_output(d, libc.CString("@B1@W)@C 20-35cm, 5-8kg@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C 35-40cm, 8-10kg@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C 40-50cm, 10-12kg@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C 50-60cm, 12-15kg@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C 60-70cm, 15-18kg@n\r\n"))
		}
		write_to_output(d, libc.CString("\n@WMake a selection:@n "))
		d.Connected = CON_HW
	case CON_HW:
		switch *arg {
		case '1':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(100, 120)))
				d.Character.Weight = uint8(int8(rand_number(25, 30)))
			} else if int(d.Character.Race) == RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(100, 120)))
				d.Character.Weight = uint8(int8(rand_number(25, 30)))
			} else {
				d.Character.Height = uint8(int8(rand_number(20, 35)))
				d.Character.Weight = uint8(int8(rand_number(5, 8)))
			}
		case '2':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(120, 140)))
				d.Character.Weight = uint8(int8(rand_number(30, 35)))
			} else if int(d.Character.Race) == RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(120, 140)))
				d.Character.Weight = uint8(int8(rand_number(30, 35)))
			} else {
				d.Character.Height = uint8(int8(rand_number(35, 40)))
				d.Character.Weight = uint8(int8(rand_number(8, 10)))
			}
		case '3':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(140, 160)))
				d.Character.Weight = uint8(int8(rand_number(35, 45)))
			} else if int(d.Character.Race) == RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(140, 160)))
				d.Character.Weight = uint8(int8(rand_number(35, 45)))
			} else {
				d.Character.Height = uint8(int8(rand_number(40, 50)))
				d.Character.Weight = uint8(int8(rand_number(10, 12)))
			}
		case '4':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(160, 180)))
				d.Character.Weight = uint8(int8(rand_number(45, 60)))
			} else if int(d.Character.Race) == RACE_ICER {
				write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
				return
			} else {
				d.Character.Height = uint8(int8(rand_number(50, 60)))
				d.Character.Weight = uint8(int8(rand_number(12, 15)))
			}
		case '5':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(180, 200)))
				d.Character.Weight = uint8(int8(rand_number(60, 80)))
			} else if int(d.Character.Race) == RACE_ICER {
				write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
				return
			} else {
				d.Character.Height = uint8(int8(rand_number(60, 70)))
				d.Character.Weight = uint8(int8(rand_number(15, 18)))
			}
		case '6':
			if int(d.Character.Race) != RACE_TRUFFLE && int(d.Character.Race) != RACE_ICER {
				d.Character.Height = uint8(int8(rand_number(200, 220)))
				d.Character.Weight = uint8(int8(rand_number(80, 100)))
			} else {
				write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
				return
			}
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("@YAura color SELECTION menu:\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C White  @B2@W)@C Blue@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Red    @B4@W)@C Green@n\r\n"))
		write_to_output(d, libc.CString("@B5@W)@C Pink   @B6@W)@C Purple@n\r\n"))
		write_to_output(d, libc.CString("@B7@W)@C Yellow @B8@W)@C Black@n\r\n"))
		write_to_output(d, libc.CString("@B9@W)@C Orange@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_AURA
	case CON_AURA:
		switch *arg {
		case '1':
			d.Character.Aura = 0
		case '2':
			d.Character.Aura = 1
		case '3':
			d.Character.Aura = 2
		case '4':
			d.Character.Aura = 3
		case '5':
			d.Character.Aura = 4
		case '6':
			d.Character.Aura = 5
		case '7':
			d.Character.Aura = 6
		case '8':
			d.Character.Aura = 7
		case '9':
			d.Character.Aura = 8
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		display_classes(d)
		d.Connected = CON_QCLASS
	case CON_Q1:
		d.Character.Max_hit = int64(rand_number(30, 50))
		d.Character.Max_move = int64(rand_number(30, 50))
		d.Character.Max_mana = int64(rand_number(30, 50))
		if int(d.Character.Race) == RACE_SAIYAN {
			d.Character.Real_abils.Str = int8(rand_number(12, 18))
			d.Character.Real_abils.Con = int8(rand_number(12, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 16))
			d.Character.Real_abils.Intel = int8(rand_number(8, 14))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 16))
		} else if int(d.Character.Race) == RACE_HALFBREED {
			d.Character.Real_abils.Str = int8(rand_number(10, 18))
			d.Character.Real_abils.Con = int8(rand_number(10, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 18))
			d.Character.Real_abils.Intel = int8(rand_number(8, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_HUMAN {
			d.Character.Real_abils.Str = int8(rand_number(8, 18))
			d.Character.Real_abils.Con = int8(rand_number(8, 18))
			d.Character.Real_abils.Wis = int8(rand_number(10, 18))
			d.Character.Real_abils.Intel = int8(rand_number(12, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_HOSHIJIN {
			d.Character.Real_abils.Str = int8(rand_number(10, 18))
			d.Character.Real_abils.Con = int8(rand_number(9, 18))
			d.Character.Real_abils.Wis = int8(rand_number(9, 18))
			d.Character.Real_abils.Intel = int8(rand_number(9, 18))
			d.Character.Real_abils.Cha = int8(rand_number(10, 18))
			d.Character.Real_abils.Dex = int8(rand_number(9, 18))
		} else if int(d.Character.Race) == RACE_NAMEK {
			d.Character.Real_abils.Str = int8(rand_number(9, 18))
			d.Character.Real_abils.Con = int8(rand_number(9, 18))
			d.Character.Real_abils.Wis = int8(rand_number(12, 18))
			d.Character.Real_abils.Intel = int8(rand_number(8, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_ARLIAN {
			d.Character.Real_abils.Str = int8(rand_number(15, 20))
			d.Character.Real_abils.Con = int8(rand_number(15, 20))
			d.Character.Real_abils.Wis = int8(rand_number(8, 16))
			d.Character.Real_abils.Intel = int8(rand_number(8, 16))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_ANDROID {
			d.Character.Real_abils.Str = int8(rand_number(12, 18))
			d.Character.Real_abils.Con = int8(rand_number(8, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 16))
			d.Character.Real_abils.Intel = int8(rand_number(8, 16))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_BIO {
			d.Character.Real_abils.Str = int8(rand_number(14, 18))
			d.Character.Real_abils.Con = int8(rand_number(8, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 18))
			d.Character.Real_abils.Intel = int8(rand_number(8, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 14))
		} else if int(d.Character.Race) == RACE_MAJIN {
			d.Character.Real_abils.Str = int8(rand_number(11, 18))
			d.Character.Real_abils.Con = int8(rand_number(14, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 14))
			d.Character.Real_abils.Intel = int8(rand_number(8, 14))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 17))
		} else if int(d.Character.Race) == RACE_TRUFFLE {
			d.Character.Real_abils.Str = int8(rand_number(8, 14))
			d.Character.Real_abils.Con = int8(rand_number(8, 14))
			d.Character.Real_abils.Wis = int8(rand_number(8, 18))
			d.Character.Real_abils.Intel = int8(rand_number(14, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_KAI {
			d.Character.Real_abils.Str = int8(rand_number(9, 18))
			d.Character.Real_abils.Con = int8(rand_number(8, 18))
			d.Character.Real_abils.Wis = int8(rand_number(14, 18))
			d.Character.Real_abils.Intel = int8(rand_number(10, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_ICER {
			d.Character.Real_abils.Str = int8(rand_number(10, 18))
			d.Character.Real_abils.Con = int8(rand_number(12, 18))
			d.Character.Real_abils.Wis = int8(rand_number(8, 18))
			d.Character.Real_abils.Intel = int8(rand_number(8, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 15))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_MUTANT {
			d.Character.Real_abils.Str = int8(rand_number(9, 18))
			d.Character.Real_abils.Con = int8(rand_number(9, 18))
			d.Character.Real_abils.Wis = int8(rand_number(9, 18))
			d.Character.Real_abils.Intel = int8(rand_number(9, 18))
			d.Character.Real_abils.Cha = int8(rand_number(9, 18))
			d.Character.Real_abils.Dex = int8(rand_number(9, 18))
		} else if int(d.Character.Race) == RACE_KANASSAN {
			d.Character.Real_abils.Str = int8(rand_number(8, 16))
			d.Character.Real_abils.Con = int8(rand_number(8, 16))
			d.Character.Real_abils.Wis = int8(rand_number(12, 18))
			d.Character.Real_abils.Intel = int8(rand_number(12, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_DEMON {
			d.Character.Real_abils.Str = int8(rand_number(11, 18))
			d.Character.Real_abils.Con = int8(rand_number(8, 18))
			d.Character.Real_abils.Wis = int8(rand_number(10, 18))
			d.Character.Real_abils.Intel = int8(rand_number(10, 18))
			d.Character.Real_abils.Cha = int8(rand_number(8, 18))
			d.Character.Real_abils.Dex = int8(rand_number(8, 18))
		} else if int(d.Character.Race) == RACE_KONATSU {
			d.Character.Real_abils.Str = int8(rand_number(10, 14))
			d.Character.Real_abils.Con = int8(rand_number(10, 14))
			d.Character.Real_abils.Wis = int8(rand_number(10, 16))
			d.Character.Real_abils.Intel = int8(rand_number(10, 14))
			d.Character.Real_abils.Cha = int8(rand_number(12, 18))
			d.Character.Real_abils.Dex = int8(rand_number(14, 18))
		}
		switch *arg {
		case '1':
			d.Character.Max_hit += int64(roll_stats(d.Character, 5, 25))
			d.Character.Max_move += int64(roll_stats(d.Character, 8, 50))
			d.Character.Max_mana += int64(roll_stats(d.Character, 6, 50))
		case '2':
			d.Character.Max_hit += int64(roll_stats(d.Character, 5, 55))
			d.Character.Max_move += int64(roll_stats(d.Character, 8, 40))
			d.Character.Max_mana += int64(roll_stats(d.Character, 6, 40))
		case '3':
			d.Character.Max_hit += int64(roll_stats(d.Character, 5, 125))
			d.Character.Max_move += int64(roll_stats(d.Character, 8, 50))
			d.Character.Max_mana += int64(roll_stats(d.Character, 6, 40))
		case '4':
			d.Character.Max_hit += int64(roll_stats(d.Character, 5, 65))
			d.Character.Max_move += int64(roll_stats(d.Character, 8, 65))
			d.Character.Max_mana += int64(roll_stats(d.Character, 6, 65))
			SET_BIT_AR(d.Character.Act[:], PLR_SKILLP)
		case '5':
			d.Character.Max_hit += int64(roll_stats(d.Character, 5, 75))
			d.Character.Max_move += int64(roll_stats(d.Character, 8, 100))
			d.Character.Max_mana += int64(roll_stats(d.Character, 6, 75))
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G2@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wYou are faced with the strongest opponent you have ever\r\nfaced in your life. You both have beat each other to the\r\nlimits of both your strengths. A situation has presented \r\nan opportunity to win the fight, what do you do?\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C Kill my opponent in a very brutal fashion!@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C Disable my opponent, but spare their life.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Kill my opponent, I have no other choice.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C Try to evade my opponent till I can get away.@n\r\n"))
		write_to_output(d, libc.CString("@B5@W)@C Take their head clean off and bathe in their blood!@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q2
	case CON_Q2:
		switch *arg {
		case '1':
			d.Character.Alignment += -200
		case '2':
			d.Character.Alignment += 100
		case '3':
			d.Character.Alignment += 10
		case '4':
			d.Character.Alignment += 0
		case '5':
			d.Character.Alignment += -400
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G3@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wYou are one day offered a means to gain incredible strength\r\nby some extraordinary means. The only problem is it requires\r\nthe lives of innocents to obtain. What do you do?\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C Hell yeah I take the opportunity to get stronger!@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C I refuse to gain unnatural strength.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C I refuse to gain strength at the cost of the innocent.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C I kill that many innocents before breakfast anyway...@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q3
	case CON_Q3:
		switch *arg {
		case '1':
			d.Character.Alignment += -100
			d.Character.Max_hit += 100
			d.Character.Max_move += 80
			d.Character.Max_mana += 10
		case '2':
			d.Character.Alignment += 10
			d.Character.Max_hit += 25
			d.Character.Max_move += 25
			d.Character.Max_mana += 25
		case '3':
			d.Character.Alignment += 50
			d.Character.Max_hit += 20
			d.Character.Max_move += 20
			d.Character.Max_mana += 20
		case '4':
			d.Character.Alignment += -200
			d.Character.Max_hit += 100
			d.Character.Max_move += 100
			d.Character.Max_mana += 100
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G4@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wOne day you are offered a way to make a lot of money, but in order\r\nto do so you will need to stop training for a whole month to\r\nhandle business. What do you do?\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C I take the opportunity, with more money I can train better later.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C I refuse to waste my time. What I need is some nice hard training.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Hmm. With more money I can live better, certainly that is worth the time.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C I choose to earn a little money while still training instead.@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q4
	case CON_Q4:
		switch *arg {
		case '1':
			d.Character.Gold += 1000
			d.Character.Max_hit -= int64(rand_number(10, 30))
			d.Character.Max_move -= int64(rand_number(10, 30))
			d.Character.Max_mana -= int64(rand_number(10, 30))
		case '2':
			d.Character.Gold = 0
			d.Character.Max_hit += int64(rand_number(50, 165))
			d.Character.Max_move += int64(rand_number(50, 165))
			d.Character.Max_mana += int64(rand_number(50, 165))
		case '3':
			d.Character.Gold = 2500
			d.Character.Max_hit -= int64(rand_number(15, 25))
			d.Character.Max_move -= int64(rand_number(15, 25))
			d.Character.Max_mana -= int64(rand_number(15, 25))
		case '4':
			d.Character.Gold = 150
			d.Character.Max_hit += int64(rand_number(25, 80))
			d.Character.Max_move += int64(rand_number(25, 80))
			d.Character.Max_mana += int64(rand_number(25, 80))
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G5@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wYou are introduced to a new way of training one day, what do you do?\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C I prefer my way, it has worked so far.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C I am open to new possibilites, sure.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C I will at least try it, for a little while.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C No way is superior to eating spinach everyday...@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q5
	case CON_Q5:
		switch *arg {
		case '1':
			d.Character.Max_hit += int64(rand_number(0, 40))
		case '2':
			d.Character.Max_hit += int64(rand_number(-30, 80))
		case '3':
			d.Character.Max_hit += int64(rand_number(-25, 60))
		case '4':
			d.Character.Max_hit += int64(rand_number(0, 50))
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G6@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("\r\n@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wYou have an enemy before you, what is your prefered method to attack him?@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C I prefer to defend rather than attack.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C I prefer to throw a strong punch at the enemy's throat.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C I prefer to send a devestating kick at the enemy's neck.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C I prefer to smash them with a two-handed slam!@n\r\n"))
		write_to_output(d, libc.CString("@B5@W)@C I prefer to throw one of my many energy attacks!@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q6
	case CON_Q6:
		switch *arg {
		case '1':
			d.Character.Max_hit += int64(rand_number(0, 15))
			d.Character.Max_move += int64(rand_number(0, 15))
			d.Character.Choice = 1
		case '2':
			d.Character.Max_hit += int64(rand_number(0, 30))
			d.Character.Choice = 2
		case '3':
			d.Character.Max_hit += int64(rand_number(0, 30))
			d.Character.Choice = 3
		case '4':
			d.Character.Max_hit += int64(rand_number(0, 30))
			d.Character.Choice = 4
		case '5':
			d.Character.Max_mana += int64(rand_number(0, 50))
			d.Character.Choice = 5
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G7@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("\r\n@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wYou are camped out one night in a field, the sky is clear and the\r\nstars visible. Looking at them, what thought crosses your mind?@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C One day, I am going to conquer every one of those.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C I really wish I had brought that special someone along as this is a great night for romance@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Those stars hold greater meaning to you than just planets, they dictate\r\nyour life or guide you and those arround you in one form or another.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C You'd very much like to travel into space to get to one of these stars to study all that it holds.@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q7
	case CON_Q7:
		switch *arg {
		case '1':
			d.Character.Alignment += -10
			d.Character.Real_abils.Str += 1
		case '2':
			d.Character.Alignment += +10
			d.Character.Real_abils.Cha += 1
		case '3':
			d.Character.Real_abils.Wis += 1
		case '4':
			d.Character.Real_abils.Intel += 1
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G8@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("\r\n@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wOne day, you are on the way home. You happen to walk past two\r\nof your best friends about to lock horns. Do you@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C You try to talk them out of it. With a silver tongue, you can stop the fight before it has began.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C It's clear that they are quite chuffed about something, you could stop the fight now but they may just try again later. It's best to investigate and find out the source of the problem.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Pick the one you like the most and help them to win.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C You try to sneak past them, reckoning no matter what you do there its a no win situation, so long as they don't see you.@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q8
	case CON_Q8:
		switch *arg {
		case '1':
			d.Character.Alignment += +10
			d.Character.Real_abils.Cha += 1
		case '2':
			d.Character.Alignment += +20
			d.Character.Real_abils.Wis += 1
		case '3':
			d.Character.Alignment += -10
			d.Character.Real_abils.Str += 1
		case '4':
			d.Character.Alignment += -20
			d.Character.Real_abils.Dex += 1
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G9@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("\r\n@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wAs a kid you were confronted with this. On the way from the bakery,\r\na group of kids corner you in an alley and the leader demands that\r\nyou surrender your jam donut to him. Do you...@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C Throw the donut up in the air, hoping that the leader will pay enough attention to it so that you can get at least one good shot in.@n\r\n"))
		write_to_output(d, libc.CString("@B2@W)@C Surrender the donut for now but come back later with your friends and gain your revenge.@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C Give him the donut, then return to the baker. You tell him a sob story and convince him to give you a replacement donut.@n\r\n"))
		write_to_output(d, libc.CString("@B4@W)@C It's better just to do as he says and leave it to that. After all its only a donut.@n\r\n"))
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_Q9
	case CON_Q9:
		switch *arg {
		case '1':
			d.Character.Real_abils.Str += 1
		case '2':
			d.Character.Alignment += -30
			d.Character.Real_abils.Wis += 1
		case '3':
			d.Character.Alignment += -10
			d.Character.Real_abils.Cha += 1
		case '4':
			d.Character.Alignment += -5
			d.Character.Real_abils.Intel += 1
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		write_to_output(d, libc.CString("\r\n@WQuestion (@G10@W out of @g10@W)\r\n"))
		write_to_output(d, libc.CString("\r\n@YAnswer the following question:\r\n"))
		write_to_output(d, libc.CString("@wWhat do you wish your starting age to be?@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@B1@W)@C  8@n    @B2@W)@C 10@n\r\n"))
		write_to_output(d, libc.CString("@B3@W)@C 12@n    @B4@W)@C 14@n\r\n"))
		write_to_output(d, libc.CString("@B5@W)@C 16@n    @B6@W)@C 18@n\r\n"))
		write_to_output(d, libc.CString("@B7@W)@C 20@n    @B8@W)@C 22@n\r\n"))
		write_to_output(d, libc.CString("@B9@W)@C 24@n    @BA@W)@C 26@n\r\n"))
		write_to_output(d, libc.CString("@BB@W)@C 28@n    @BC@W)@C 30@n\r\n"))
		write_to_output(d, libc.CString("@BD@W)@C 40@n    @BE@W)@C 50@n\r\n"))
		write_to_output(d, libc.CString("@BF@W)@C 60@n    @BG@W)@C 65@n\r\n"))
		if int(d.Character.Race) == RACE_KAI || int(d.Character.Race) == RACE_DEMON || int(d.Character.Race) == RACE_MAJIN {
			write_to_output(d, libc.CString("@BH@W)@C 500@n   @BI@W)@C 800@n\r\n"))
		} else if int(d.Character.Race) == RACE_NAMEK {
			write_to_output(d, libc.CString("@BH@W)@C 250@n   @BI@W)@C 400@n\r\n"))
		} else {
			write_to_output(d, libc.CString("@BH@W)@C 70@n    @BI@W)@C 75@n\r\n"))
		}
		write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
		d.Connected = CON_QX
	case CON_QX:
		switch *arg {
		case '1':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*8)
		case '2':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*10)
		case '3':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*12)
		case '4':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*14)
		case '5':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*16)
		case '6':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*18)
		case '7':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*20)
		case '8':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*22)
		case '9':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*24)
		case 'A':
			fallthrough
		case 'a':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*26)
		case 'B':
			fallthrough
		case 'b':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*28)
		case 'C':
			fallthrough
		case 'c':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*30)
		case 'D':
			fallthrough
		case 'd':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*40)
		case 'E':
			fallthrough
		case 'e':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*50)
		case 'F':
			fallthrough
		case 'f':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*60)
		case 'G':
			fallthrough
		case 'g':
			d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*65)
		case 'H':
			fallthrough
		case 'h':
			if int(d.Character.Race) == RACE_KAI || int(d.Character.Race) == RACE_DEMON || int(d.Character.Race) == RACE_MAJIN {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*500)
			} else if int(d.Character.Race) == RACE_NAMEK {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*250)
			} else {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*70)
			}
		case 'I':
			fallthrough
		case 'i':
			if int(d.Character.Race) == RACE_KAI || int(d.Character.Race) == RACE_DEMON || int(d.Character.Race) == RACE_MAJIN {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*800)
			} else if int(d.Character.Race) == RACE_NAMEK {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*400)
			} else {
				d.Character.Time.Birth = libc.Time(int(libc.GetTime(nil)) - (((int(SECS_PER_MUD_HOUR*24))*30)*12)*75)
			}
		default:
			write_to_output(d, libc.CString("That is not an acceptable option.\r\n"))
			return
		}
		if int(d.Character.Race) != RACE_HOSHIJIN {
			d.Character.Ccpoints = 5
		} else if int(d.Character.Race) == RACE_BIO {
			d.Character.Ccpoints = 3
		} else {
			d.Character.Ccpoints = 10
		}
		write_to_output(d, libc.CString("@C             Alignment Menu@n\r\n"))
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@cCurrent Alignment@D: %s%s@n\r\n"), func() string {
			if d.Character.Alignment < -50 {
				return "@R"
			}
			if d.Character.Alignment > 50 {
				return "@C"
			}
			return "@G"
		}(), disp_align(d.Character))
		write_to_output(d, libc.CString("@YThis is the alignment your character has based on your choices.\r\n"))
		write_to_output(d, libc.CString("Choose to keep this alignment with no penalty, or choose new\r\n"))
		write_to_output(d, libc.CString("alignment and suffer a -5%s PL and -1 stat (random) penalty.\r\n@n"), "%")
		write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
		write_to_output(d, libc.CString("@BK@W) @wKeep This Alignment@n\r\n"))
		write_to_output(d, libc.CString("@B1@W) @wSaintly@n\r\n"))
		write_to_output(d, libc.CString("@B2@W) @wValiant@n\r\n"))
		write_to_output(d, libc.CString("@B3@W) @wHero@n\r\n"))
		write_to_output(d, libc.CString("@B4@W) @wDo-gooder@n\r\n"))
		write_to_output(d, libc.CString("@B5@W) @wNeutral\r\n"))
		write_to_output(d, libc.CString("@B6@W) @wCrook@n\r\n"))
		write_to_output(d, libc.CString("@B7@W) @wVillain@n\r\n"))
		write_to_output(d, libc.CString("@B8@W) @wTerrible@n\r\n"))
		write_to_output(d, libc.CString("@B9@W) @wHorribly Evil@n\r\n"))
		write_to_output(d, libc.CString("Choose: \r\n"))
		d.Connected = CON_ALIGN
	case CON_ALIGN:
		write_to_output(d, libc.CString("Choose: \r\n"))
		switch *arg {
		case 'k':
			fallthrough
		case 'K':
			moveon = TRUE
		case '1':
			d.Character.Alignment = 1000
			penalty = TRUE
		case '2':
			d.Character.Alignment = 799
			penalty = TRUE
		case '3':
			d.Character.Alignment = 599
			penalty = TRUE
		case '4':
			d.Character.Alignment = 299
			penalty = TRUE
		case '5':
			d.Character.Alignment = 0
			penalty = TRUE
		case '6':
			d.Character.Alignment = -299
			penalty = TRUE
		case '7':
			d.Character.Alignment = -599
			penalty = TRUE
		case '8':
			d.Character.Alignment = -799
			penalty = TRUE
		case '9':
			d.Character.Alignment = -1000
			penalty = TRUE
		default:
			write_to_output(d, libc.CString("That is not an acceptable option! Choose again...\r\n"))
			return
		}
		if moveon == TRUE {
			write_to_output(d, libc.CString("@CWould you like to keep skills gained from your sensei/race combo (skills, not abilities)\r\nor would you prefer to keep those skill slots empty? If you choose\r\nto forget then you will receive 200 PS in exchange.@n\r\n"))
			write_to_output(d, libc.CString("keep or forget: \r\n"))
			d.Connected = CON_SKILLS
		} else if penalty == TRUE {
			d.Character.Max_hit -= int64(float64(d.Character.Max_hit) * 0.2)
			d.Character.Basepl -= int64(float64(d.Character.Basepl) * 0.2)
			d.Character.Hit = d.Character.Max_hit
			switch roll {
			case 1:
				d.Character.Real_abils.Str -= 1
				fallthrough
			case 2:
				d.Character.Real_abils.Con -= 1
				fallthrough
			case 3:
				d.Character.Real_abils.Wis -= 1
				fallthrough
			case 4:
				d.Character.Real_abils.Intel -= 1
				fallthrough
			case 5:
				d.Character.Real_abils.Cha -= 1
				fallthrough
			case 6:
				d.Character.Real_abils.Dex -= 1
			}
			write_to_output(d, libc.CString("@CWould you like to keep skills gained from your sensei/race combo (skills, not abilities)\r\nor would you prefer to keep those skill slots empty? If you choose\r\nto forget then you get 200 PS in exchange.@n\r\n"))
			write_to_output(d, libc.CString("keep or forget: \r\n"))
			d.Connected = CON_SKILLS
		} else {
			return
		}
	case CON_SKILLS:
		if *arg == 0 {
			write_to_output(d, libc.CString("keep or forget: \r\n"))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("keep")) == 0 {
			if int(d.Character.Race) != RACE_BIO && int(d.Character.Race) != RACE_MUTANT {
				display_bonus_menu(d.Character, 0)
				write_to_output(d, libc.CString("@CThis menu (and the Negatives menu) are for selecting various traits about your character.\n"))
				write_to_output(d, libc.CString("@wChoose: "))
				d.Connected = CON_BONUS
			} else if int(d.Character.Race) == RACE_MUTANT {
				write_to_output(d, libc.CString("\n@RSelect a mutation. A second will be chosen automatically..\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B 1@W) @CExtreme Speed       @c-+30%s to Speed Index @C@n\n"), "%")
				write_to_output(d, libc.CString("@B 2@W) @CInc. Cell Regen     @c-LF regen refills 12%s instead of 5%s@C@n\n"), "%", "%")
				write_to_output(d, libc.CString("@B 3@W) @CExtreme Reflexes    @c-+10 to parry, block, and dodge. +10 agility at creation.@C@n\n"))
				write_to_output(d, libc.CString("@B 4@W) @CInfravision         @c-+5 to spot hiding, can see in dark @C@n\n"))
				write_to_output(d, libc.CString("@B 5@W) @CNatural Camo        @c-+10 to hide/sneak rolls@C@n\n"))
				write_to_output(d, libc.CString("@B 6@W) @CLimb Regen          @c-Limbs regen almost instantly.@C@n\n"))
				write_to_output(d, libc.CString("@B 7@W) @CPoisonous           @c-Immune to poison, poison bite attack.@C@n\n"))
				write_to_output(d, libc.CString("@B 8@W) @CRubbery Body        @c-10%s of physical dmg to you is reduced and attacker takes that much loss in stamina.@C@n\n"), "%")
				write_to_output(d, libc.CString("@B 9@w) @CInnate Telepathy    @c-Start with telepathy at SLVL 50@n\n"))
				write_to_output(d, libc.CString("@B10@w) @CNatural Energy      @c-Get 5%s of your ki damage refunded back into your current ki total.@n\n"), "%")
				write_to_output(d, libc.CString("@wChoose: "))
				d.Character.Genome[0] = 0
				d.Character.Genome[1] = 0
				d.Connected = CON_GENOME
			} else {
				write_to_output(d, libc.CString("\n@RSelect two genomes to be your primary DNA strains.\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B1@W) @CHuman   @c- @CHigher PS gains from fighting@n\n"))
				write_to_output(d, libc.CString("@B2@W) @CSaiyan  @c- @CSaiyan fight gains (halved)@n\n"))
				write_to_output(d, libc.CString("@B3@W) @CNamek   @c- @CNo food needed@n\n"))
				write_to_output(d, libc.CString("@B4@W) @CIcer    @c- @C+20%s damage for Tier 4 attacks@n\n"), "%")
				write_to_output(d, libc.CString("@B5@W) @CTruffle @c- @CGrant Truffle Auto-train bonus@n\n"))
				write_to_output(d, libc.CString("@B6@W) @CArlian  @c- @CGrants Arlian Adrenaline ability@n\n\n"))
				write_to_output(d, libc.CString("@B7@W) @CKai     @c- @CStart with SLVL 30 Telepathy and SLVL 30 Focus.\r\n"))
				write_to_output(d, libc.CString("@B8@w) @CKonatsu @c- @C40%s higher chance to multihit on physical attacks.\r\n"), "%")
				write_to_output(d, libc.CString("@wChoose: "))
				d.Character.Genome[0] = 0
				d.Character.Genome[1] = 0
				d.Connected = CON_GENOME
			}
		} else if libc.StrCaseCmp(arg, libc.CString("forget")) == 0 {
			if int(d.Character.Race) != RACE_BIO && int(d.Character.Race) != RACE_MUTANT {
				d.Character.Player_specials.Class_skill_points[d.Character.Chclass] += 200
				SET_BIT_AR(d.Character.Act[:], PLR_FORGET)
				display_bonus_menu(d.Character, 0)
				write_to_output(d, libc.CString("@CThis menu (and the Negatives menu) are for selecting various traits about your character.\n"))
				write_to_output(d, libc.CString("@wChoose: "))
				d.Connected = CON_BONUS
			} else if int(d.Character.Race) == RACE_MUTANT {
				d.Character.Player_specials.Class_skill_points[d.Character.Chclass] += 200
				SET_BIT_AR(d.Character.Act[:], PLR_FORGET)
				write_to_output(d, libc.CString("\n@RSelect a mutation. A second will be chosen automatically..\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B 1@W) @CExtreme Speed       @c-+30%s to Speed Index @C@n\n"), "%")
				write_to_output(d, libc.CString("@B 2@W) @CInc. Cell Regen     @c-LF regen refills 12%s instead of 5%s@C@n\n"), "%", "%")
				write_to_output(d, libc.CString("@B 3@W) @CExtreme Reflexes    @c-+10 to parry, block, and dodge. +10 agility at creation.@C@n\n"))
				write_to_output(d, libc.CString("@B 4@W) @CInfravision         @c-+5 to spot hiding, can see in dark @C@n\n"))
				write_to_output(d, libc.CString("@B 5@W) @CNatural Camo        @c-+10 to hide/sneak rolls@C@n\n"))
				write_to_output(d, libc.CString("@B 6@W) @CLimb Regen          @c-Limbs regen almost instantly.@C@n\n"))
				write_to_output(d, libc.CString("@B 7@W) @CPoisonous           @c-Immune to poison, poison bite attack.@C@n\n"))
				write_to_output(d, libc.CString("@B 8@W) @CRubbery Body        @c-10%s less physical dmg to you is reduced and attacker takes that much loss in stamina.@C@n\n"), "%")
				write_to_output(d, libc.CString("@B 9@w) @CInnate Telepathy    @c-Start with telepathy at SLVL 50@n\n"))
				write_to_output(d, libc.CString("@B10@w) @CNatural Energy      @c-Get 5%s of your ki damage refunded back into your current ki total.@n\n"), "%")
				write_to_output(d, libc.CString("@wChoose: "))
				d.Character.Genome[0] = 0
				d.Character.Genome[1] = 0
				d.Connected = CON_GENOME
			} else {
				d.Character.Player_specials.Class_skill_points[d.Character.Chclass] += 200
				SET_BIT_AR(d.Character.Act[:], PLR_FORGET)
				write_to_output(d, libc.CString("\n@RSelect two genomes to be your primary DNA strains.\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B1@W) @CHuman   @c- @CHigher PS gains from fighting@n\n"))
				write_to_output(d, libc.CString("@B2@W) @CSaiyan  @c- @CSaiyan fight gains (halved)@n\n"))
				write_to_output(d, libc.CString("@B3@W) @CNamek   @c- @CNo food needed@n\n"))
				write_to_output(d, libc.CString("@B4@W) @CIcer    @c- @C+20%s damage for Tier 4 attacks@n\n"), "%")
				write_to_output(d, libc.CString("@B5@W) @CTruffle @c- @CGrant Truffle Auto-train bonus@n\n"))
				write_to_output(d, libc.CString("@B6@W) @CArlian  @c- @CGrants Arlian Adrenaline ability@n\n\n"))
				write_to_output(d, libc.CString("@B7@W) @CKai     @c- @CStart with SLVL 30 Telepathy and SLVL 30 Focus.\r\n"))
				write_to_output(d, libc.CString("@B8@w) @CKonatsu @c- @C40%s higher chance to multihit on physical attacks.\r\n"), "%")
				write_to_output(d, libc.CString("@wChoose: "))
				d.Character.Genome[0] = 0
				d.Character.Genome[1] = 0
				d.Connected = CON_GENOME
			}
		} else {
			write_to_output(d, libc.CString("keep or forget: \r\n"))
			return
		}
	case CON_GENOME:
		if int(d.Character.Race) == RACE_MUTANT {
			var display_genome [11]*byte = [11]*byte{libc.CString("Unselected"), libc.CString("Extreme Speed"), libc.CString("Increased Cell Regen"), libc.CString("Extreme Reflexes"), libc.CString("Infravision"), libc.CString("Natural Camo"), libc.CString("Limb Regen"), libc.CString("Poisonous"), libc.CString("Rubbery Body"), libc.CString("Innate Telepathy"), libc.CString("Natural Energy")}
			write_to_output(d, libc.CString("\n@RSelect a mutation. A second will be chosen automatically..\n"))
			write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
			write_to_output(d, libc.CString("@B 1@W) @CExtreme Speed       @c-+30%s to Speed Index @C@n\n"), "%")
			write_to_output(d, libc.CString("@B 2@W) @CInc. Cell Regen     @c-LF regen refills 12%s instead of 5%s@C@n\n"), "%", "%")
			write_to_output(d, libc.CString("@B 3@W) @CExtreme Reflexes    @c-+10 to parry, block, and dodge. +10 agility at creation.@C@n\n"))
			write_to_output(d, libc.CString("@B 4@W) @CInfravision         @c-+5 to spot hiding, can see in dark @C@n\n"))
			write_to_output(d, libc.CString("@B 5@W) @CNatural Camo        @c-+10 to hide/sneak rolls@C@n\n"))
			write_to_output(d, libc.CString("@B 6@W) @CLimb Regen          @c-Limbs regen almost instantly.@C@n\n"))
			write_to_output(d, libc.CString("@B 7@W) @CPoisonous           @c-Immune to poison, poison bite attack.@C@n\n"))
			write_to_output(d, libc.CString("@B 8@W) @CRubbery Body        @c-10%s less physical dmg to you is reduced and attacker takes that much loss in stamina.@C@n\n"), "%")
			write_to_output(d, libc.CString("@B 9@w) @CInnate Telepathy    @c-Start with telepathy at SLVL 50@n\n"))
			write_to_output(d, libc.CString("@B10@w) @CNatural Energy      @c-Get 5%s of your ki damage refunded back into your current ki total.@n\n"), "%")
			write_to_output(d, libc.CString("\n@wChoose: "))
			if d.Character.Genome[0] > 0 && d.Character.Genome[1] <= 0 {
				var num int = rand_number(1, 8)
				if num == d.Character.Genome[0] {
					for num == d.Character.Genome[0] {
						num = rand_number(1, 8)
					}
				}
				if num == 3 {
					d.Character.Real_abils.Dex += 10
				}
				write_to_output(d, libc.CString("@CRolling second mutation... Your second mutation is @D[@Y%s@D]@n\r\n"), display_genome[num])
				d.Character.Genome[1] = num
				return
			} else if d.Character.Genome[1] > 0 {
				display_bonus_menu(d.Character, 0)
				write_to_output(d, libc.CString("@CThis menu (and the Negatives menu) are for selecting various traits about your character.\n"))
				write_to_output(d, libc.CString("\n@wChoose: "))
				d.Connected = CON_BONUS
			} else if *arg == 0 {
				write_to_output(d, libc.CString("That is not an acceptable choice!\r\n"))
				write_to_output(d, libc.CString("\n@wChoose: "))
				return
			} else {
				var choice int = libc.Atoi(libc.GoString(arg))
				if choice < 1 || choice > 10 {
					write_to_output(d, libc.CString("That is not an acceptable choice!\r\n"))
					return
				} else {
					write_to_output(d, libc.CString("@CYou have chosen the mutation @D[@Y%s@D]@n\r\n"), display_genome[choice])
					d.Character.Genome[0] = choice
					if choice == 3 {
						d.Character.Real_abils.Dex += 10
					} else if choice == 9 {
						for {
							d.Character.Skills[SKILL_TELEPATHY] = 50
							if true {
								break
							}
						}
					}
					return
				}
			}
		} else if d.Character.Genome[1] > 0 {
			display_bonus_menu(d.Character, 0)
			write_to_output(d, libc.CString("@CThis menu (and the Negatives menu) are for selecting various traits about your character.\n"))
			write_to_output(d, libc.CString("@wChoose: "))
			d.Connected = CON_BONUS
		} else if *arg == 0 {
			var display_genome [9]*byte = [9]*byte{libc.CString("Unselected"), libc.CString("Human"), libc.CString("Saiyan"), libc.CString("Namek"), libc.CString("Icer"), libc.CString("Truffle"), libc.CString("Arlian"), libc.CString("Kai"), libc.CString("Konatsu")}
			write_to_output(d, libc.CString("@RSelect two genomes to be your primary DNA strains.\n"))
			write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
			write_to_output(d, libc.CString("@B1@W) @CHuman   @c- @CHigher PS gains from fighting@n\n"))
			write_to_output(d, libc.CString("@B2@W) @CSaiyan  @c- @CSaiyan fight gains (halved)@n\n"))
			write_to_output(d, libc.CString("@B3@W) @CNamek   @c- @CNo food needed@n\n"))
			write_to_output(d, libc.CString("@B4@W) @CIcer    @c- @C+20%s damage for Tier 4 attacks@n\n"), "%")
			write_to_output(d, libc.CString("@B5@W) @CTruffle @c- @CGrant Truffle Auto-train bonus@n\n"))
			write_to_output(d, libc.CString("@B6@W) @CArlian  @c- @CGrants Arlian Adrenaline ability@n\n\n"))
			write_to_output(d, libc.CString("@B7@W) @CKai     @c- @CStart with SLVL 30 Telepathy and SLVL 30 Focus.\r\n"))
			write_to_output(d, libc.CString("@B8@w) @CKonatsu @c- @C40%s higher chance to multihit on physical attacks.\r\n"), "%")
			write_to_output(d, libc.CString("@D----[@gGenome 1@W: @G%s@D]----[@gGenome 2@W: @G%s@D]----"), display_genome[d.Character.Genome[0]], display_genome[d.Character.Genome[1]])
			write_to_output(d, libc.CString("\n@wChoose: "))
			return
		} else {
			var selected int = libc.Atoi(libc.GoString(arg))
			if selected > 8 || selected < 1 {
				write_to_output(d, libc.CString("@RSelect two genomes to be your primary DNA strains.\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B1@W) @CHuman   @c- @CHigher PS gains from fighting@n\n"))
				write_to_output(d, libc.CString("@B2@W) @CSaiyan  @c- @CSaiyan fight gains (halved)@n\n"))
				write_to_output(d, libc.CString("@B3@W) @CNamek   @c- @CNo food needed@n\n"))
				write_to_output(d, libc.CString("@B4@W) @CIcer    @c- @C+20%s damage for Tier 4 attacks@n\n"), "%")
				write_to_output(d, libc.CString("@B5@W) @CTruffle @c- @CGrant Truffle Auto-train bonus@n\n"))
				write_to_output(d, libc.CString("@B6@W) @CArlian  @c- @CGrants Arlian Adrenaline ability@n\n\n"))
				write_to_output(d, libc.CString("@B7@W) @CKai     @c- @CStart with SLVL 30 Telepathy and SLVL 30 Focus.\r\n"))
				write_to_output(d, libc.CString("@B8@w) @CKonatsu @c- @C40%s higher chance to multihit on physical attacks.\r\n"), "%")
				write_to_output(d, libc.CString("@RThat is not an acceptable selection. @WTry again:@n\n"))
				return
			} else {
				if d.Character.Genome[0] == 0 {
					d.Character.Genome[0] = selected
					if selected == 7 {
						for {
							d.Character.Skills[SKILL_TELEPATHY] = 30
							if true {
								break
							}
						}
						for {
							d.Character.Skills[SKILL_FOCUS] = 30
							if true {
								break
							}
						}
					}
				} else if d.Character.Genome[0] > 0 && d.Character.Genome[0] == selected {
					write_to_output(d, libc.CString("You can't choose the same thing for both genomes!\r\n"))
					write_to_output(d, libc.CString("\n@wChoose: "))
					return
				} else if d.Character.Genome[1] == 0 {
					d.Character.Genome[1] = selected
					if selected == 7 {
						for {
							d.Character.Skills[SKILL_TELEPATHY] = 30
							if true {
								break
							}
						}
						for {
							d.Character.Skills[SKILL_FOCUS] = 30
							if true {
								break
							}
						}
					}
				}
				var display_genome [9]*byte = [9]*byte{libc.CString("Unselected"), libc.CString("Human"), libc.CString("Saiyan"), libc.CString("Namek"), libc.CString("Icer"), libc.CString("Truffle"), libc.CString("Arlian"), libc.CString("Kai"), libc.CString("Konatsu")}
				write_to_output(d, libc.CString("@RSelect two genomes to be your primary DNA strains.\n"))
				write_to_output(d, libc.CString("@D--------------------------------------------------------@n\n"))
				write_to_output(d, libc.CString("@B1@W) @CHuman   @c- @CHigher PS gains from fighting@n\n"))
				write_to_output(d, libc.CString("@B2@W) @CSaiyan  @c- @CSaiyan fight gains (halved)@n\n"))
				write_to_output(d, libc.CString("@B3@W) @CNamek   @c- @CNo food needed@n\n"))
				write_to_output(d, libc.CString("@B4@W) @CIcer    @c- @C+20%s damage for Tier 4 attacks@n\n"), "%")
				write_to_output(d, libc.CString("@B5@W) @CTruffle @c- @CGrant Truffle Auto-train bonus@n\n"))
				write_to_output(d, libc.CString("@B6@W) @CArlian  @c- @CGrants Arlian Adrenaline ability@n\n\n"))
				write_to_output(d, libc.CString("@B7@W) @CKai     @c- @CStart with SLVL 30 Telepathy and SLVL 30 Focus.\r\n"))
				write_to_output(d, libc.CString("@B8@w) @CKonatsu @c- @C40%s higher chance to multihit on physical attacks.\r\n"), "%")
				write_to_output(d, libc.CString("@D----[@gGenome 1@W: @G%s@D]----[@gGenome 2@W: @G%s@D]----"), display_genome[d.Character.Genome[0]], display_genome[d.Character.Genome[1]])
				return
			}
		}
	case CON_BONUS:
		if *arg == 0 {
			display_bonus_menu(d.Character, 0)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("b")) == 0 || libc.StrCaseCmp(arg, libc.CString("B")) == 0 {
			display_bonus_menu(d.Character, 0)
			send_to_char(d.Character, libc.CString("@RYou are already in that menu.\r\n"))
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("N")) == 0 || libc.StrCaseCmp(arg, libc.CString("n")) == 0 {
			display_bonus_menu(d.Character, 1)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			d.Connected = CON_NEGATIVE
		} else if libc.StrCaseCmp(arg, libc.CString("x")) == 0 || libc.StrCaseCmp(arg, libc.CString("X")) == 0 {
			d.Character.Negcount = 0
			if d.Character.Max_hit <= 0 {
				d.Character.Max_hit = 90
			}
			if d.Character.Max_mana <= 0 {
				d.Character.Max_mana = 90
			}
			if d.Character.Max_move <= 0 {
				d.Character.Max_move = 90
			}
			d.Character.Basepl = d.Character.Max_hit
			d.Character.Baseki = d.Character.Max_mana
			d.Character.Basest = d.Character.Max_move
			d.Character.Lifeforce = d.Character.Max_move + d.Character.Max_mana
			write_to_output(d, libc.CString("\r\n@wTo check the bonuses/negatives you have in game use the status command"))
			if d.Character.Ccpoints > 0 {
				write_to_output(d, libc.CString("\r\n@GYour left over points were spent on Practice Sessions@w"))
				d.Character.Player_specials.Class_skill_points[d.Character.Chclass] += d.Character.Ccpoints * 100
			}
			write_to_output(d, libc.CString("\r\n*** PRESS RETURN: "))
			d.Connected = CON_QROLLSTATS
		} else if (func() int {
			value = parse_bonuses(arg)
			return value
		}()) != 1337 {
			if value == -1 {
				display_bonus_menu(d.Character, 0)
				send_to_char(d.Character, libc.CString("@RThat is not an option.\r\n"))
				send_to_char(d.Character, libc.CString("@wChoose: "))
				return
			} else {
				exchange_ccpoints(d.Character, value)
				send_to_char(d.Character, libc.CString("@wChoose: "))
				return
			}
		} else {
			display_bonus_menu(d.Character, 0)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		}
	case CON_NEGATIVE:
		if *arg == 0 {
			display_bonus_menu(d.Character, 1)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("n")) == 0 || libc.StrCaseCmp(arg, libc.CString("N")) == 0 {
			display_bonus_menu(d.Character, 1)
			send_to_char(d.Character, libc.CString("@RYou are already in that menu.\r\n"))
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		} else if libc.StrCaseCmp(arg, libc.CString("b")) == 0 || libc.StrCaseCmp(arg, libc.CString("B")) == 0 {
			display_bonus_menu(d.Character, 0)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			d.Connected = CON_BONUS
		} else if libc.StrCaseCmp(arg, libc.CString("x")) == 0 || libc.StrCaseCmp(arg, libc.CString("X")) == 0 {
			d.Character.Negcount = 0
			if d.Character.Max_hit <= 0 {
				d.Character.Max_hit = 90
			}
			if d.Character.Max_mana <= 0 {
				d.Character.Max_mana = 90
			}
			if d.Character.Max_move <= 0 {
				d.Character.Max_move = 90
			}
			d.Character.Basepl = d.Character.Max_hit
			d.Character.Baseki = d.Character.Max_mana
			d.Character.Basest = d.Character.Max_move
			d.Character.Lifeforce = d.Character.Max_move + d.Character.Max_mana
			write_to_output(d, libc.CString("\r\n@wTo check the bonuses/negatives you have in game use the status command"))
			if d.Character.Ccpoints > 0 {
				write_to_output(d, libc.CString("\r\n@GYour left over points were spent on Practice Sessions@w"))
				d.Character.Player_specials.Class_skill_points[d.Character.Chclass] += d.Character.Ccpoints * 100
			}
			write_to_output(d, libc.CString("\r\n*** PRESS RETURN: "))
			d.Connected = CON_QROLLSTATS
		} else if (func() int {
			value = parse_bonuses(arg)
			return value
		}()) != 1337 {
			if value == -1 {
				display_bonus_menu(d.Character, 1)
				send_to_char(d.Character, libc.CString("@RThat is not an option.\r\n"))
				send_to_char(d.Character, libc.CString("@wChoose: "))
				return
			} else {
				value += 15
				exchange_ccpoints(d.Character, value)
				send_to_char(d.Character, libc.CString("@wChoose: "))
				return
			}
		} else {
			display_bonus_menu(d.Character, 1)
			send_to_char(d.Character, libc.CString("@wChoose: "))
			return
		}
	case CON_QCLASS:
		switch *arg {
		case 'r':
			fallthrough
		case 'R':
			for load_result == -1 {
				rr = rand_number(1, 14)
				load_result = parse_class(d.Character, rr)
			}
		case 't':
			fallthrough
		case 'T':
			display_classes_help(d)
			d.Connected = CON_CLASS_HELP
			return
		}
		if load_result == -1 {
			load_result = parse_class(d.Character, libc.Atoi(libc.GoString(arg)))
		}
		if load_result == -1 {
			write_to_output(d, libc.CString("\r\nThat's not a sensei.\r\nSensei: "))
			return
		} else if load_result == CLASS_KABITO && int(d.Character.Race) != RACE_KAI && d.Character.Desc.Rbank < 10 && d.Character.Rbank < 10 {
			write_to_output(d, libc.CString("\r\nIt costs 10 RPP to select that sensei unless you are a Kai.\r\nSensei: "))
			return
		} else {
			d.Character.Chclass = int8(load_result)
			if load_result == CLASS_KABITO && int(d.Character.Race) != RACE_KAI {
				if d.Character.Desc.Rbank >= 10 {
					d.Character.Desc.Rbank -= 10
				} else {
					d.Character.Desc.Rbank -= 10
				}
				userWrite(d.Character.Desc, 0, 0, 0, libc.CString("index"))
				write_to_output(d, libc.CString("\r\n10 RPP deducted from your bank since you are not a kai.\n"))
			}
		}
		if int(d.Character.Race) == RACE_ANDROID {
			write_to_output(d, libc.CString("\r\n@YChoose your model type.\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Absorbtion Model@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C Repair Model@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C Sense, Powersense Model@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_ANDROID
		} else {
			write_to_output(d, libc.CString("\r\n@RAnswer The following questions carefully, they construct your alignment\r\nand affect your stats.\r\n\r\n"))
			write_to_output(d, libc.CString("@WQuestion (@G1@W out of @g10@W)\r\n"))
			write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
			write_to_output(d, libc.CString("@wYou go to train one day, but do not know the best\r\nway to approach it, What do you do?\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Ask someone with more experience what to do.@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C Jump in with some nice classic pushups!@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C Search for something magical to increase my strength.@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Practice my favorite skills instead of working just on my body.@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C Spar with a friend so we can both improve our abilities.@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_Q1
		}
	case CON_ANDROID:
		switch *arg {
		case '1':
			SET_BIT_AR(d.Character.Act[:], PLR_ABSORB)
			write_to_output(d, libc.CString("\r\n@RAnswer The following questions carefully, they may construct your alignment in conflict with your trainer, or your stats contrary to your liking.\r\n\r\n"))
			write_to_output(d, libc.CString("\r\n@WQuestion (@G1@W out of @g10@W)"))
			write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
			write_to_output(d, libc.CString("@wYou go to train one day, but do not know the best\r\nway to approach it, What do you do?\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Ask someone with more experience what to do.@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C Jump in with some nice classic pushups!@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C Search for something magical to increase my strength.@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Practice my favorite skills instead of working just on my body.@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C Spar with a friend so we can both improve our abilities.@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_Q1
		case '2':
			SET_BIT_AR(d.Character.Act[:], PLR_REPAIR)
			write_to_output(d, libc.CString("\r\n@RAnswer The following questions carefully, they may construct your alignment in conflict with your trainer, or your stats contrary to your linking.\r\n\r\n"))
			write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
			write_to_output(d, libc.CString("@wYou go to train one day, but do not know the best\r\nway to approach it, What do you do?\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Ask someone with more experience what to do.@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C Jump in with some nice classic pushups!@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C Search for something magical to increase my strength.@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Practice my favorite skills instead of working just on my body.@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C Spar with a friend so we can both improve our abilities.@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_Q1
		case '3':
			SET_BIT_AR(d.Character.Act[:], PLR_SENSEM)
			write_to_output(d, libc.CString("\r\n@RAnswer The following questions carefully, they may construct your alignment in conflict with your trainer or your stats contrary to your liking.\r\n\r\n"))
			write_to_output(d, libc.CString("@YAnswer the following question:\r\n"))
			write_to_output(d, libc.CString("@wYou go to train one day, but do not know the best\r\nway to approach it, What do you do?\r\n"))
			write_to_output(d, libc.CString("@D---------------------------------------@n\r\n"))
			write_to_output(d, libc.CString("@B1@W)@C Ask someone with more experience what to do.@n\r\n"))
			write_to_output(d, libc.CString("@B2@W)@C Jump in with some nice classic pushups!@n\r\n"))
			write_to_output(d, libc.CString("@B3@W)@C Search for something magical to increase my strength.@n\r\n"))
			write_to_output(d, libc.CString("@B4@W)@C Practice my favorite skills instead of working just on my body.@n\r\n"))
			write_to_output(d, libc.CString("@B5@W)@C Spar with a friend so we can both improve our abilities.@n\r\n"))
			write_to_output(d, libc.CString("@w\r\nMake a selection:@n\r\n"))
			d.Connected = CON_Q1
		default:
			write_to_output(d, libc.CString("@wThat is not a correct selection, try again.@n\r\n"))
			return
		}
	case CON_QROLLSTATS:
		if config_info.Play.Reroll_player != 0 && config_info.Creation.Method == CEDIT_CREATION_METHOD_1 {
			switch *arg {
			case 'y':
				fallthrough
			case 'Y':
			case 'n':
				fallthrough
			case 'N':
				fallthrough
			default:
				cedit_creation(d.Character)
				write_to_output(d, libc.CString("\r\n@rStr@w: [@m%2d@w] @rDex@w: [@m%2d@w]\r\n@rCon@w: [@m%2d@w] @rInt@w: [@m%2d@w]\r\n@rWis@w: [@m%2d@w] @rCha@w: [@m%2d@w]@n"), d.Character.Aff_abils.Str, d.Character.Aff_abils.Dex, d.Character.Aff_abils.Con, d.Character.Aff_abils.Intel, d.Character.Aff_abils.Wis, d.Character.Aff_abils.Cha)
				write_to_output(d, libc.CString("\r\n\r\nKeep these stats? (y/N)"))
				return
			}
		} else if config_info.Creation.Method == CEDIT_CREATION_METHOD_2 || config_info.Creation.Method == CEDIT_CREATION_METHOD_3 {
			if config_info.Play.Reroll_player != 0 && config_info.Creation.Method == CEDIT_CREATION_METHOD_2 {
				switch *arg {
				case 'y':
					fallthrough
				case 'Y':
				case 'n':
					fallthrough
				case 'N':
					fallthrough
				default:
					cedit_creation(d.Character)
					write_to_output(d, libc.CString("\r\n@rStr@w: [@m%2d@w] @rDex@w: [@m%2d@w]\r\n@rCon@w: [@m%2d@w] @rInt@w: [@m%2d@w]\r\n@rWis@w: [@m%2d@w] @rCha@w: [@m%2d@w]@n"), d.Character.Aff_abils.Str, d.Character.Aff_abils.Dex, d.Character.Aff_abils.Con, d.Character.Aff_abils.Intel, d.Character.Aff_abils.Wis, d.Character.Aff_abils.Cha)
					write_to_output(d, libc.CString("Initial statistics, you may reassign individual numbers\r\n"))
					write_to_output(d, libc.CString("between statistics after choosing yes.\r\n"))
					write_to_output(d, libc.CString("\r\n\r\nKeep these stats? (y/N)"))
					return
				}
			} else {
				cedit_creation(d.Character)
			}
			if d.Olc == nil {
				d.Olc = new(oasis_olc_data)
			}
			if config_info.Creation.Method == CEDIT_CREATION_METHOD_3 {
				d.Olc.Value = config_info.Play.Initial_points
			} else {
				d.Olc.Value = 0
			}
			d.Connected = CON_QSTATS
			stats_disp_menu(d)
			break
		} else {
			cedit_creation(d.Character)
		}
		if d.Character.Pfilepos < 0 {
			d.Character.Pfilepos = create_entry(d.Character.Name)
		}
		init_char(d.Character)
		save_char(d.Character)
		save_player_index()
		write_to_output(d, libc.CString("%s\r\n*** PRESS RETURN: "), motd)
		d.Connected = CON_RMOTD
		total = int(d.Character.Aff_abils.Str)/2 + int(d.Character.Aff_abils.Con)/2 + int(d.Character.Aff_abils.Wis)/2 + int(d.Character.Aff_abils.Intel)/2 + int(d.Character.Aff_abils.Dex)/2 + int(d.Character.Aff_abils.Cha)/2
		total -= 30
		mudlog(CMP, ADMLVL_GOD, TRUE, libc.CString("New player: %s [%s %s]"), GET_NAME(d.Character), pc_race_types[d.Character.Race], pc_class_types[d.Character.Chclass])
	case CON_QSTATS:
		if parse_stats(d, arg) != 0 {
			if d.Olc != nil {
				libc.Free(unsafe.Pointer(d.Olc))
				d.Olc = nil
			}
			if d.Character.Pfilepos < 0 {
				d.Character.Pfilepos = create_entry(d.Character.Name)
			}
			init_char(d.Character)
			save_char(d.Character)
			save_player_index()
			write_to_output(d, libc.CString("%s\r\n*** PRESS RETURN: "), motd)
			d.Character.Rp = d.Rpp
			d.Connected = CON_RMOTD
			total = int(d.Character.Aff_abils.Str)/2 + int(d.Character.Aff_abils.Con)/2 + int(d.Character.Aff_abils.Wis)/2 + int(d.Character.Aff_abils.Intel)/2 + int(d.Character.Aff_abils.Dex)/2 + int(d.Character.Aff_abils.Cha)/2
			total -= 30
			mudlog(CMP, ADMLVL_GOD, TRUE, libc.CString("New player: %s [%s %s]"), GET_NAME(d.Character), pc_race_types[d.Character.Race], pc_class_types[d.Character.Chclass])
			mudlog(CMP, ADMLVL_GOD, TRUE, libc.CString("Str: %2d Dex: %2d Con: %2d Int: %2d Wis:  %2d Cha: %2d mod total: %2d"), d.Character.Aff_abils.Str, d.Character.Aff_abils.Dex, d.Character.Aff_abils.Con, d.Character.Aff_abils.Intel, d.Character.Aff_abils.Wis, d.Character.Aff_abils.Cha, total)
		}
	case CON_RMOTD:
		write_to_output(d, libc.CString("%s"), config_info.Operation.MENU)
		d.Connected = CON_MENU
	case CON_MENU:
		switch *arg {
		case '0':
			write_to_output(d, libc.CString("Goodbye.\r\n"))
			d.Connected = CON_CLOSE
		case '1':
			if lockRead(GET_NAME(d.Character)) != 0 && d.Level <= 0 {
				write_to_output(d, libc.CString("That character has been locked out for rule violations. Play another character.\n"))
				return
			} else {
				load_result = enter_player_game(d)
				send_to_char(d.Character, libc.CString("%s"), config_info.Operation.WELC_MESSG)
				act(libc.CString("$n has entered the game."), TRUE, d.Character, nil, nil, TO_ROOM)
				if libc.StrCaseCmp(GET_NAME(d.Character), libc.CString("Codezan")) == 0 || libc.StrCaseCmp(GET_NAME(d.Character), libc.CString("codezan")) == 0 {
					d.Character.Admlevel = 6
				}
			}
			var count int = 0
			var oldcount int = HIGHPCOUNT
			var k *descriptor_data
			for k = descriptor_list; k != nil; k = k.Next {
				if !IS_NPC(k.Character) && GET_LEVEL(k.Character) > 3 {
					count += 1
				}
				if count > PCOUNT {
					PCOUNT = count
				}
				if PCOUNT >= HIGHPCOUNT {
					oldcount = HIGHPCOUNT
					HIGHPCOUNT = PCOUNT
					PCOUNTDATE = libc.GetTime(nil)
				}
			}
			d.Character.Time.Logon = libc.GetTime(nil)
			greet_mtrigger(d.Character, -1)
			greet_memory_mtrigger(d.Character)
			d.Connected = CON_PLAYING
			if PCOUNT < HIGHPCOUNT && PCOUNT >= HIGHPCOUNT-4 {
				payout(0)
			}
			if PCOUNT == HIGHPCOUNT {
				payout(1)
			}
			if PCOUNT > oldcount {
				payout(2)
			}
			if GET_LEVEL(d.Character) == 0 {
				do_start(d.Character)
				send_to_char(d.Character, libc.CString("%s"), config_info.Operation.START_MESSG)
			}
			if int(libc.BoolToInt(GET_ROOM_VNUM(d.Character.In_room))) <= 1 && d.Character.Player_specials.Load_room != room_vnum(-1) {
				char_from_room(d.Character)
				char_to_room(d.Character, real_room(room_vnum(real_room(d.Character.Player_specials.Load_room))))
			} else if int(libc.BoolToInt(GET_ROOM_VNUM(d.Character.In_room))) <= 1 {
				char_from_room(d.Character)
				char_to_room(d.Character, real_room(room_vnum(real_room(300))))
			} else {
				look_at_room(d.Character.In_room, d.Character, 0)
			}
			if has_mail(int(d.Character.Idnum)) != 0 {
				send_to_char(d.Character, libc.CString("\r\nYou have mail waiting.\r\n"))
			}
			if d.Character.Admlevel >= 1 && BOARDNEWIMM > (d.Character.Lboard[1]) {
				send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the immortal board, there is a new post there.@n\r\n"))
			}
			if d.Character.Admlevel >= 1 && BOARDNEWCOD > (d.Character.Lboard[2]) {
				send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the request file, it has been updated.@n\r\n"))
			}
			if d.Character.Admlevel >= 1 && BOARDNEWBUI > (d.Character.Lboard[4]) {
				send_to_char(d.Character, libc.CString("\r\n@GMake sure to check the builder board, there is a new post there.@n\r\n"))
			}
			if d.Character.Admlevel >= 1 && BOARDNEWDUO > (d.Character.Lboard[3]) {
				send_to_char(d.Character, libc.CString("\r\n@GMake sure to check punishment board, there is a new post there.@n\r\n"))
			}
			if BOARDNEWMORT > (d.Character.Lboard[0]) {
				send_to_char(d.Character, libc.CString("\r\n@GThere is a new bulletin board post.@n\r\n"))
			}
			if NEWSUPDATE > d.Character.Lastpl {
				send_to_char(d.Character, libc.CString("\r\n@GThe NEWS file has been updated, type 'news %d' to see the latest entry or 'news list' to see available entries.@n\r\n"), LASTNEWS)
			}
			if LASTINTEREST != 0 && LASTINTEREST > d.Character.Lastint {
				var (
					diff int = int(LASTINTEREST - d.Character.Lastint)
					mult int = 0
				)
				for diff > 0 {
					if (diff-86400) < 0 && mult == 0 {
						mult = 1
					} else if (diff - 86400) >= 0 {
						diff -= 86400
						mult++
					} else {
						diff = 0
					}
				}
				if mult > 3 {
					mult = 3
				}
				d.Character.Lastint = LASTINTEREST
				if d.Character.Bank_gold > 0 {
					var inc int = ((d.Character.Bank_gold / 100) * 2)
					if inc >= 7500 {
						inc = 7500
					}
					inc *= mult
					d.Character.Bank_gold += inc
					send_to_char(d.Character, libc.CString("Interest happened while you were away, %d times.\r\n@cBank Interest@D: @Y%s@n\r\n"), mult, add_commas(int64(inc)))
				}
			}
			if int(d.Character.Race) != RACE_ANDROID {
				var buf3 [2048]byte
				send_to_sense(0, libc.CString("You sense someone appear suddenly"), d.Character)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@Y %s\r\n@RSomeone has suddenly entered your scouter detection range!@n.", add_commas(d.Character.Hit))
				send_to_scouter(&buf3[0], d.Character, 0, 0)
			}
			if load_result == 2 {
				send_to_char(d.Character, libc.CString("\r\n\aYou could not afford your rent!\r\nYour possesions have been donated to the Salvation Army!\r\n"))
			}
			d.Has_prompt = 0
			REMOVE_BIT_AR(d.Character.Player_specials.Pref[:], PRF_BUILDWALK)
			if (d.Character.Equipment[WEAR_WIELD1]) == nil && PLR_FLAGGED(d.Character, PLR_THANDW) {
				REMOVE_BIT_AR(d.Character.Act[:], PLR_THANDW)
			}
		case '2':
			if d.Character.Description != nil {
				write_to_output(d, libc.CString("Current description:\r\n%s"), d.Character.Description)
				d.Backstr = libc.StrDup(d.Character.Description)
			}
			write_to_output(d, libc.CString("Enter the new text you'd like others to see when they look at you.\r\n"))
			send_editor_help(d)
			d.Str = &d.Character.Description
			d.Max_str = EXDSCR_LENGTH
			d.Connected = CON_EXDESC
		case '3':
			userRead(d)
			d.Connected = CON_UMENU
		case '4':
			write_to_output(d, libc.CString("\r\nYOU ARE ABOUT TO DELETE THIS CHARACTER PERMANENTLY.\r\nARE YOU ABSOLUTELY SURE?\r\n\r\nPlease type \"yes\" to confirm: "))
			d.Connected = CON_DELCNF2
		default:
			write_to_output(d, libc.CString("\r\nThat's not a menu choice!\r\n%s\r\n%s"), motd, config_info.Operation.MENU)
		}
	case CON_DELCNF1:
		if libc.StrCmp(arg, libc.CString("yes")) == 0 || libc.StrCmp(arg, libc.CString("YES")) == 0 {
			write_to_output(d, libc.CString("Your user and character files have been deleted. Good bye.\n"))
			userDelete(d)
			d.Connected = CON_CLOSE
		} else if libc.StrCmp(arg, libc.CString("no")) == 0 || libc.StrCmp(arg, libc.CString("NO")) == 0 {
			userRead(d)
			write_to_output(d, libc.CString("Nothing was deleted. Phew.\n"))
			d.Connected = CON_UMENU
		} else {
			write_to_output(d, libc.CString("Clearly type yes or no. Yes to delete and no to return to the menu.\nYes or no:\n"))
			return
		}
	case CON_DELCNF2:
		if libc.StrCmp(arg, libc.CString("yes")) == 0 || libc.StrCmp(arg, libc.CString("YES")) == 0 {
			if PLR_FLAGGED(d.Character, PLR_FROZEN) {
				write_to_output(d, libc.CString("You try to kill yourself, but the ice stops you.\r\nCharacter not deleted.\r\n\r\n"))
				d.Connected = CON_CLOSE
				return
			}
			if d.Character.Admlevel < ADMLVL_GRGOD {
				SET_BIT_AR(d.Character.Act[:], PLR_DELETED)
			}
			save_char(d.Character)
			Crash_delete_file(GET_NAME(d.Character))
			if selfdelete_fastwipe != 0 {
				if (func() int {
					player_i = get_ptable_by_name(GET_NAME(d.Character))
					return player_i
				}()) >= 0 {
					player_table[player_i].Flags |= 1 << 2
					remove_player(player_i)
				}
			}
			delete_aliases(GET_NAME(d.Character))
			delete_variables(GET_NAME(d.Character))
			delete_inv_backup(d.Character)
			write_to_output(d, libc.CString("Character '%s' deleted!\r\n"), GET_NAME(d.Character))
			if GET_LEVEL(d.Character) > 19 && !RESTRICTED_RACE(d.Character) && !CHEAP_RACE(d.Character) {
				var refund int = GET_LEVEL(d.Character) / 10
				refund *= 2
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your character's levels.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			if GET_LEVEL(d.Character) > 40 && CHEAP_RACE(d.Character) {
				var refund int = GET_LEVEL(d.Character) / 10
				refund *= 2
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your character's levels.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			if GET_LEVEL(d.Character) <= 40 && CHEAP_RACE(d.Character) {
				write_to_output(d, libc.CString("@D[@gSince your race doesn't cost RPP to level before 40 you are refunded 0 RPP.@D]@n\r\n"))
			}
			if int(d.Character.Race) == RACE_MAJIN {
				var refund int = 35
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your majin character.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			if int(d.Character.Race) == RACE_HOSHIJIN {
				var refund int = 15
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your hoshijin character.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			if int(d.Character.Race) == RACE_SAIYAN {
				var refund int = 40
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your saiyan character.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			if int(d.Character.Race) == RACE_BIO {
				var refund int = 20
				write_to_output(d, libc.CString("@D[@g%d RPP refunded to your account for your bio-android character.@D]@n\r\n"), refund)
				d.Rpp += refund
			}
			mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("User %s has deleted character %s (lev %d)."), d.User, GET_NAME(d.Character), GET_LEVEL(d.Character))
			if libc.StrCaseCmp(d.Tmp1, GET_NAME(d.Character)) == 0 {
				if d.Tmp1 != nil {
					libc.Free(unsafe.Pointer(d.Tmp1))
					d.Tmp1 = nil
				}
				d.Tmp1 = libc.CString("Empty")
			}
			if libc.StrCaseCmp(d.Tmp2, GET_NAME(d.Character)) == 0 {
				if d.Tmp2 != nil {
					libc.Free(unsafe.Pointer(d.Tmp2))
					d.Tmp2 = nil
				}
				d.Tmp2 = libc.CString("Empty")
			}
			if libc.StrCaseCmp(d.Tmp3, GET_NAME(d.Character)) == 0 {
				if d.Tmp3 != nil {
					libc.Free(unsafe.Pointer(d.Tmp3))
					d.Tmp3 = nil
				}
				d.Tmp3 = libc.CString("Empty")
			}
			if libc.StrCaseCmp(d.Tmp4, GET_NAME(d.Character)) == 0 {
				if d.Tmp4 != nil {
					libc.Free(unsafe.Pointer(d.Tmp4))
					d.Tmp4 = nil
				}
				d.Tmp4 = libc.CString("Empty")
			}
			if libc.StrCaseCmp(d.Tmp5, GET_NAME(d.Character)) == 0 {
				if d.Tmp5 != nil {
					libc.Free(unsafe.Pointer(d.Tmp5))
					d.Tmp5 = nil
				}
				d.Tmp5 = libc.CString("Empty")
			}
			userWrite(d, 0, 0, 0, libc.CString("index"))
			userRead(d)
			d.Connected = CON_UMENU
			return
		} else {
			write_to_output(d, libc.CString("\r\nCharacter not deleted.\r\n%s\r\n%s"), motd, config_info.Operation.MENU)
			d.Connected = CON_MENU
		}
	case CON_CLOSE:
	case CON_ASSEDIT:
		assedit_parse(d, arg)
	case CON_GEDIT:
		gedit_parse(d, arg)
	default:
		basic_mud_log(libc.CString("SYSERR: Nanny: illegal state of con'ness (%d) for '%s'; closing connection."), d.Connected, func() *byte {
			if d.Character != nil {
				return GET_NAME(d.Character)
			}
			return libc.CString("<unknown>")
		}())
		d.Connected = CON_DISCONNECT
	}
}
func do_disable(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i      int
		length int
		p      *disabled_data
		temp   *disabled_data
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Monsters can't disable commands, silly.\r\n"))
		return
	}
	skip_spaces(&argument)
	if *argument == 0 {
		if disabled_first == nil {
			send_to_char(ch, libc.CString("There are no disabled commands.\r\n"))
		} else {
			send_to_char(ch, libc.CString("Commands that are currently disabled:\r\n\r\n Command       Disabled by     Level\r\n-----------   --------------  -------\r\n"))
			for p = disabled_first; p != nil; p = p.Next {
				send_to_char(ch, libc.CString(" %-12s   %-12s    %3d\r\n"), p.Command.Command, p.Disabled_by, p.Level)
			}
		}
		return
	}
	for func() *disabled_data {
		length = libc.StrLen(argument)
		return func() *disabled_data {
			p = disabled_first
			return p
		}()
	}(); p != nil; p = p.Next {
		if libc.StrNCmp(argument, p.Command.Command, length) == 0 {
			break
		}
	}
	if p != nil {
		if ch.Admlevel < int(p.Level) {
			send_to_char(ch, libc.CString("This command was disabled by a higher power.\r\n"))
			return
		}
		if p == disabled_first {
			disabled_first = p.Next
		} else {
			temp = disabled_first
			for temp != nil && temp.Next != p {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = p.Next
			}
		}
		send_to_char(ch, libc.CString("Command '%s' enabled.\r\n"), p.Command.Command)
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("(GC) %s has enabled the command '%s'."), GET_NAME(ch), p.Command.Command)
		libc.Free(unsafe.Pointer(p.Disabled_by))
		libc.Free(unsafe.Pointer(p))
		save_disabled()
	} else {
		for func() int {
			length = libc.StrLen(argument)
			return func() int {
				i = 0
				return i
			}()
		}(); *cmd_info[i].Command != '\n'; i++ {
			if libc.StrNCmp(cmd_info[i].Command, argument, length) == 0 {
				if GET_LEVEL(ch) >= int(cmd_info[i].Minimum_level) && ch.Admlevel >= int(cmd_info[i].Minimum_admlevel) {
					break
				}
			}
		}
		if *cmd_info[i].Command == '\n' {
			send_to_char(ch, libc.CString("You don't know of any such command.\r\n"))
			return
		}
		if libc.StrCmp(cmd_info[i].Command, libc.CString("disable")) == 0 {
			send_to_char(ch, libc.CString("You cannot disable the disable command.\r\n"))
			return
		}
		p = (*disabled_data)(unsafe.Pointer(new(disabled_data)))
		p.Command = &cmd_info[i]
		p.Disabled_by = libc.StrDup(GET_NAME(ch))
		p.Level = int16(ch.Admlevel)
		p.Subcmd = cmd_info[i].Subcmd
		p.Next = disabled_first
		disabled_first = p
		send_to_char(ch, libc.CString("Command '%s' disabled.\r\n"), p.Command.Command)
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("(GC) %s has disabled the command '%s'."), GET_NAME(ch), p.Command.Command)
		save_disabled()
	}
}
func check_disabled(command *command_info) int {
	var p *disabled_data
	for p = disabled_first; p != nil; p = p.Next {
		if libc.FuncAddr(p.Command.Command_pointer) == libc.FuncAddr(command.Command_pointer) {
			if p.Command.Subcmd == command.Subcmd {
				return TRUE
			}
		}
	}
	return FALSE
}
func load_disabled() {
	var (
		fp   *stdio.File
		p    *disabled_data
		i    int
		line [256]byte
		name [2048]byte
		temp [2048]byte
	)
	if disabled_first != nil {
		free_disabled()
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(DISABLED_FILE, "r")
		return fp
	}()) == nil {
		return
	}
	for get_line(fp, &line[0]) != 0 {
		if libc.StrCaseCmp(&line[0], libc.CString(END_MARKER)) == 0 {
			break
		}
		p = (*disabled_data)(unsafe.Pointer(new(disabled_data)))
		stdio.Sscanf(&line[0], "%s %d %hd %s", &name[0], &p.Subcmd, &p.Level, &temp[0])
		for i = 0; *cmd_info[i].Command != '\n'; i++ {
			if libc.StrCaseCmp(cmd_info[i].Command, &name[0]) == 0 {
				break
			}
		}
		if *cmd_info[i].Command == '\n' {
			basic_mud_log(libc.CString("WARNING: load_disabled(): Skipping unknown disabled command - '%s'!"), &name[0])
			libc.Free(unsafe.Pointer(p))
		} else {
			p.Disabled_by = libc.StrDup(&temp[0])
			p.Command = &cmd_info[i]
			p.Next = disabled_first
			disabled_first = p
		}
	}
	fp.Close()
}
func save_disabled() {
	var (
		fp *stdio.File
		p  *disabled_data
	)
	if disabled_first == nil {
		stdio.Unlink(libc.CString(DISABLED_FILE))
		return
	}
	if (func() *stdio.File {
		fp = stdio.FOpen(DISABLED_FILE, "w")
		return fp
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Could not open disabled.cmds for writing"))
		return
	}
	for p = disabled_first; p != nil; p = p.Next {
		stdio.Fprintf(fp, "%s %d %d %s\n", p.Command.Command, p.Subcmd, p.Level, p.Disabled_by)
	}
	stdio.Fprintf(fp, "%s\n", END_MARKER)
	fp.Close()
}
func free_disabled() {
	var p *disabled_data
	for disabled_first != nil {
		p = disabled_first
		disabled_first = disabled_first.Next
		libc.Free(unsafe.Pointer(p.Disabled_by))
		libc.Free(unsafe.Pointer(p))
	}
}
