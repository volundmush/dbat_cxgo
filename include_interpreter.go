package main

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