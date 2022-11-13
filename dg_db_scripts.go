package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

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
	Nr            int64
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

func parse_trigger(trig_f *stdio.File, nr int) {
	var (
		t           [2]int
		k           int
		attach_type int
		line        [256]byte
		cmds        *byte
		s           *byte
		flags       [256]byte
		errors      [2048]byte
		cle         *cmdlist_element
		t_index     *index_data
		trig        *trig_data
	)
	trig = new(trig_data)
	t_index = new(index_data)
	t_index.Vnum = mob_vnum(nr)
	t_index.Number = 0
	t_index.Func = nil
	t_index.Proto = trig
	stdio.Snprintf(&errors[0], int(2048), "trig vnum %d", nr)
	trig.Nr = int64(top_of_trigt)
	trig.Name = fread_string(trig_f, &errors[0])
	get_line(trig_f, &line[0])
	k = stdio.Sscanf(&line[0], "%d %s %d", &attach_type, &flags[0], &t[0])
	trig.Attach_type = int8(attach_type)
	trig.Trigger_type = int(asciiflag_conv(&flags[0]))
	if k == 3 {
		trig.Narg = t[0]
	} else {
		trig.Narg = 0
	}
	trig.Arglist = fread_string(trig_f, &errors[0])
	cmds = func() *byte {
		s = fread_string(trig_f, &errors[0])
		return s
	}()
	trig.Cmdlist = new(cmdlist_element)
	trig.Cmdlist.Cmd = libc.StrDup(libc.StrTok(s, libc.CString("\n\r")))
	cle = trig.Cmdlist
	for (func() *byte {
		s = libc.StrTok(nil, libc.CString("\n\r"))
		return s
	}()) != nil {
		cle.Next = new(cmdlist_element)
		cle = cle.Next
		cle.Cmd = libc.StrDup(s)
	}
	libc.Free(unsafe.Pointer(cmds))
	*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(func() int {
		p := &top_of_trigt
		x := *p
		*p++
		return x
	}()))) = t_index
}
func read_trigger(nr int) *trig_data {
	var (
		t_index *index_data
		trig    *trig_data
	)
	if nr >= top_of_trigt {
		return nil
	}
	if (func() *index_data {
		t_index = *(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(nr)))
		return t_index
	}()) == nil {
		return nil
	}
	trig = new(trig_data)
	trig_data_copy(trig, t_index.Proto)
	t_index.Number++
	return trig
}
func trig_data_init(this_data *trig_data) {
	this_data.Nr = -1
	this_data.Data_type = 0
	this_data.Name = nil
	this_data.Trigger_type = 0
	this_data.Cmdlist = nil
	this_data.Curr_state = nil
	this_data.Narg = 0
	this_data.Arglist = nil
	this_data.Depth = 0
	this_data.Wait_event = nil
	this_data.Purged = FALSE != 0
	this_data.Var_list = nil
	this_data.Next = nil
}
func trig_data_copy(this_data *trig_data, trg *trig_data) {
	trig_data_init(this_data)
	this_data.Nr = trg.Nr
	this_data.Attach_type = trg.Attach_type
	this_data.Data_type = trg.Data_type
	if trg.Name != nil {
		this_data.Name = libc.StrDup(trg.Name)
	} else {
		this_data.Name = libc.CString("unnamed trigger")
		basic_mud_log(libc.CString("Trigger with no name! (%d)"), trg.Nr)
	}
	this_data.Trigger_type = trg.Trigger_type
	this_data.Cmdlist = trg.Cmdlist
	this_data.Narg = trg.Narg
	if trg.Arglist != nil {
		this_data.Arglist = libc.StrDup(trg.Arglist)
	}
}
func dg_read_trigger(fp *stdio.File, proto unsafe.Pointer, type_ int) {
	var (
		line      [256]byte
		junk      [8]byte
		vnum      int
		rnum      int
		count     int
		mob       *char_data
		room      *room_data
		trg_proto *trig_proto_list
		new_trg   *trig_proto_list
	)
	get_line(fp, &line[0])
	count = stdio.Sscanf(&line[0], "%7s %d", &junk[0], &vnum)
	if count != 2 {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Error assigning trigger! - Line was\n  %s"), &line[0])
		return
	}
	rnum = int(real_trigger(trig_vnum(vnum)))
	if rnum == int(-1) {
		switch type_ {
		case MOB_TRIGGER:
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: dg_read_trigger: Trigger vnum #%d asked for but non-existant! (mob: %s - %d)"), vnum, GET_NAME((*char_data)(proto)), GET_MOB_VNUM((*char_data)(proto)))
		case WLD_TRIGGER:
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: dg_read_trigger: Trigger vnum #%d asked for but non-existant! (room:%d)"), vnum, func() room_vnum {
				if ((*room_data)(proto)).Number != room_vnum(-1) && ((*room_data)(proto)).Number <= room_vnum(top_of_world) {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*room_data)(proto)).Number)))).Number
				}
				return -1
			}())
		default:
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: dg_read_trigger: Trigger vnum #%d asked for but non-existant! (?)"), vnum)
		}
		return
	}
	switch type_ {
	case MOB_TRIGGER:
		new_trg = new(trig_proto_list)
		new_trg.Vnum = vnum
		new_trg.Next = nil
		mob = (*char_data)(proto)
		trg_proto = mob.Proto_script
		if trg_proto == nil {
			mob.Proto_script = func() *trig_proto_list {
				trg_proto = new_trg
				return trg_proto
			}()
		} else {
			for trg_proto.Next != nil {
				trg_proto = trg_proto.Next
			}
			trg_proto.Next = new_trg
		}
	case WLD_TRIGGER:
		new_trg = new(trig_proto_list)
		new_trg.Vnum = vnum
		new_trg.Next = nil
		room = (*room_data)(proto)
		trg_proto = room.Proto_script
		if trg_proto == nil {
			room.Proto_script = func() *trig_proto_list {
				trg_proto = new_trg
				return trg_proto
			}()
		} else {
			for trg_proto.Next != nil {
				trg_proto = trg_proto.Next
			}
			trg_proto.Next = new_trg
		}
		if rnum != int(-1) {
			if room.Script == nil {
				room.Script = new(script_data)
			}
			add_trigger(room.Script, read_trigger(rnum), -1)
		} else {
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: non-existant trigger #%d assigned to room #%d"), vnum, room.Number)
		}
	default:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Trigger vnum #%d assigned to non-mob/obj/room"), vnum)
	}
}
func dg_obj_trigger(line *byte, obj *obj_data) {
	var (
		junk      [8]byte
		vnum      int
		rnum      int
		count     int
		trg_proto *trig_proto_list
		new_trg   *trig_proto_list
	)
	count = stdio.Sscanf(line, "%s %d", &junk[0], &vnum)
	if count != 2 {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: dg_obj_trigger() : Error assigning trigger! - Line was:\n  %s"), line)
		return
	}
	rnum = int(real_trigger(trig_vnum(vnum)))
	if rnum == int(-1) {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Trigger vnum #%d asked for but non-existant! (Object: %s - %d)"), vnum, obj.Short_description, GET_OBJ_VNUM(obj))
		return
	}
	new_trg = new(trig_proto_list)
	new_trg.Vnum = vnum
	new_trg.Next = nil
	trg_proto = obj.Proto_script
	if trg_proto == nil {
		obj.Proto_script = func() *trig_proto_list {
			trg_proto = new_trg
			return trg_proto
		}()
	} else {
		for trg_proto.Next != nil {
			trg_proto = trg_proto.Next
		}
		trg_proto.Next = new_trg
	}
}
func assign_triggers(i unsafe.Pointer, type_ int) {
	var (
		mob       *char_data = nil
		obj       *obj_data  = nil
		room      *room_data = nil
		rnum      int
		trg_proto *trig_proto_list
	)
	switch type_ {
	case MOB_TRIGGER:
		mob = (*char_data)(i)
		trg_proto = mob.Proto_script
		for trg_proto != nil {
			rnum = int(real_trigger(trig_vnum(trg_proto.Vnum)))
			if rnum == int(-1) {
				mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: trigger #%d non-existant, for mob #%d"), trg_proto.Vnum, (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(mob.Nr)))).Vnum)
			} else {
				if mob.Script == nil {
					mob.Script = new(script_data)
				}
				add_trigger(mob.Script, read_trigger(rnum), -1)
			}
			trg_proto = trg_proto.Next
		}
	case OBJ_TRIGGER:
		obj = (*obj_data)(i)
		trg_proto = obj.Proto_script
		for trg_proto != nil {
			rnum = int(real_trigger(trig_vnum(trg_proto.Vnum)))
			if rnum == int(-1) {
				basic_mud_log(libc.CString("SYSERR: trigger #%d non-existant, for obj #%d"), trg_proto.Vnum, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(obj.Item_number)))).Vnum)
			} else {
				if obj.Script == nil {
					obj.Script = new(script_data)
				}
				add_trigger(obj.Script, read_trigger(rnum), -1)
			}
			trg_proto = trg_proto.Next
		}
	case WLD_TRIGGER:
		room = (*room_data)(i)
		trg_proto = room.Proto_script
		for trg_proto != nil {
			rnum = int(real_trigger(trig_vnum(trg_proto.Vnum)))
			if rnum == int(-1) {
				mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: trigger #%d non-existant, for room #%d"), trg_proto.Vnum, room.Number)
			} else {
				if room.Script == nil {
					room.Script = new(script_data)
				}
				add_trigger(room.Script, read_trigger(rnum), -1)
			}
			trg_proto = trg_proto.Next
		}
	default:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: unknown type for assign_triggers()"))
	}
}
