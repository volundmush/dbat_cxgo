package main

import "unsafe"

const DG_SCRIPT_VERSION = "DG Scripts 1.0.14"
const MOB_TRIGGER = 0
const OBJ_TRIGGER = 1
const WLD_TRIGGER = 2
const DG_CASTER_PROXY = 1
const DG_SPELL_LEVEL = 25
const ACTOR_ROOM_IS_UID = 1
const MTRIG_GLOBAL = 1
const MTRIG_RANDOM = 2
const MTRIG_COMMAND = 4
const MTRIG_SPEECH = 8
const MTRIG_ACT = 16
const MTRIG_DEATH = 32
const MTRIG_GREET = 64
const MTRIG_GREET_ALL = 128
const MTRIG_ENTRY = 256
const MTRIG_RECEIVE = 512
const MTRIG_FIGHT = 1024
const MTRIG_HITPRCNT = 2048
const MTRIG_BRIBE = 4096
const MTRIG_LOAD = 8192
const MTRIG_MEMORY = 0x4000
const MTRIG_CAST = 0x8000
const MTRIG_LEAVE = 0x10000
const MTRIG_DOOR = 0x20000
const MTRIG_TIME = 0x80000
const OTRIG_GLOBAL = 1
const OTRIG_RANDOM = 2
const OTRIG_COMMAND = 4
const OTRIG_TIMER = 32
const OTRIG_GET = 64
const OTRIG_DROP = 128
const OTRIG_GIVE = 256
const OTRIG_WEAR = 512
const OTRIG_REMOVE = 2048
const OTRIG_LOAD = 8192
const OTRIG_CAST = 0x8000
const OTRIG_LEAVE = 0x10000
const OTRIG_CONSUME = 0x40000
const OTRIG_TIME = 0x80000
const WTRIG_GLOBAL = 1
const WTRIG_RANDOM = 2
const WTRIG_COMMAND = 4
const WTRIG_SPEECH = 8
const WTRIG_RESET = 32
const WTRIG_ENTER = 64
const WTRIG_DROP = 128
const WTRIG_CAST = 0x8000
const WTRIG_LEAVE = 0x10000
const WTRIG_DOOR = 0x20000
const WTRIG_TIME = 0x80000
const OCMD_EQUIP = 1
const OCMD_INVEN = 2
const OCMD_ROOM = 4
const OCMD_EAT = 1
const OCMD_DRINK = 2
const OCMD_QUAFF = 3
const TRIG_NEW = 0
const TRIG_RESTART = 1
const PULSE_DG_SCRIPT = 130
const MAX_SCRIPT_DEPTH = 10
const SCRIPT_ERROR_CODE = -9999999
const DG_ALLOW_GODS = 1
const UID_CHAR = 125
const MOB_ID_BASE = 50000
const ROOM_ID_BASE = 1050000
const OBJ_ID_BASE = 1300000

type cmdlist_element struct {
	Cmd      *byte
	Original *cmdlist_element
	Next     *cmdlist_element
}
type trig_var_data struct {
	Name    *byte
	Value   *byte
	Context int
	Next    *trig_var_data
}
type trig_data struct {
	Nr            int
	Attach_type   int8
	Data_type     int8
	Name          *byte
	Trigger_type  int
	Cmdlist       *cmdlist_element
	Curr_state    *cmdlist_element
	Narg          int
	Arglist       *byte
	Depth         int
	Loops         int
	Wait_event    *event
	Purged        bool
	Var_list      *trig_var_data
	Next          *trig_data
	Next_in_world *trig_data
}
type script_data struct {
	Types       int
	Trig_list   *trig_data
	Global_vars *trig_var_data
	Purged      bool
	Context     int
	Next        *script_data
}
type wait_event_data struct {
	Trigger *trig_data
	Gohere  unsafe.Pointer
	Type    int
}
type script_memory struct {
	Id   int
	Cmd  *byte
	Next *script_memory
}
