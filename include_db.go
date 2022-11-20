package main

import "github.com/gotranspile/cxgo/runtime/libc"

const DB_BOOT_WLD = 0
const DB_BOOT_MOB = 1
const DB_BOOT_OBJ = 2
const DB_BOOT_ZON = 3
const DB_BOOT_SHP = 4
const DB_BOOT_HLP = 5
const DB_BOOT_TRG = 6
const DB_BOOT_GLD = 7

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
	Arg1    int
	Arg2    int
	Arg3    int
	Arg4    int
	Arg5    int
	Line    int
	Sarg1   *byte
	Sarg2   *byte
}
type zone_data struct {
	Name       *byte
	Builders   *byte
	Lifespan   int
	Age        int
	Bot        int
	Top        int
	Reset_mode int
	Number     int
	Cmd        []reset_com
	Min_level  int
	Max_level  int
	Zone_flags [4]uint32
}
type reset_q_element struct {
	Zone_to_reset int
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
