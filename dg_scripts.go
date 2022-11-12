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

const PULSES_PER_MUD_HOUR = 9000
const BUCKET_COUNT = 64
const UID_OUT_OF_RANGE = 1000000000

var trig_wait_event func(event_obj unsafe.Pointer) int
var do_attach func(ch *char_data, argument *byte, cmd int, subcmd int)
var do_detach func(ch *char_data, argument *byte, cmd int, subcmd int)
var do_vdelete func(ch *char_data, argument *byte, cmd int, subcmd int)
var do_tstat func(ch *char_data, argument *byte, cmd int, subcmd int)

func str_str(cs *byte, ct *byte) *byte {
	var (
		s *byte
		t *byte
	)
	if cs == nil || ct == nil || *ct == 0 {
		return nil
	}
	for *cs != 0 {
		t = ct
		for *cs != 0 && C.tolower(int(*cs)) != C.tolower(int(*t)) {
			cs = (*byte)(unsafe.Add(unsafe.Pointer(cs), 1))
		}
		s = cs
		for *t != 0 && *cs != 0 && C.tolower(int(*cs)) == C.tolower(int(*t)) {
			t = (*byte)(unsafe.Add(unsafe.Pointer(t), 1))
			cs = (*byte)(unsafe.Add(unsafe.Pointer(cs), 1))
		}
		if *t == 0 {
			return s
		}
	}
	return nil
}
func trgvar_in_room(vnum room_vnum) int {
	var (
		rnum room_rnum = real_room(vnum)
		i    int       = 0
		ch   *char_data
	)
	if rnum == room_rnum(-1) {
		script_log(libc.CString("people.vnum: world[rnum] does not exist"))
		return -1
	}
	for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).People; ch != nil; ch = ch.Next_in_room {
		i++
	}
	return i
}
func get_obj_in_list(name *byte, list *obj_data) *obj_data {
	var (
		i  *obj_data
		id int
	)
	if *name == UID_CHAR {
		id = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1))))
		for i = list; i != nil; i = i.Next_content {
			if id == int(i.Id) {
				return i
			}
		}
	} else {
		for i = list; i != nil; i = i.Next_content {
			if isname(name, i.Name) != 0 {
				return i
			}
		}
	}
	return nil
}
func get_object_in_equip(ch *char_data, name *byte) *obj_data {
	var (
		j       int
		n       int = 0
		number  int
		obj     *obj_data
		tmpname [2048]byte
		tmp     *byte = &tmpname[0]
		id      int
	)
	if *name == UID_CHAR {
		id = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1))))
		for j = 0; j < NUM_WEARS; j++ {
			if (func() *obj_data {
				obj = ch.Equipment[j]
				return obj
			}()) != nil {
				if id == int(obj.Id) {
					return obj
				}
			}
		}
	} else if is_number(name) != 0 {
		var ovnum obj_vnum = obj_vnum(libc.Atoi(libc.GoString(name)))
		for j = 0; j < NUM_WEARS; j++ {
			if (func() *obj_data {
				obj = ch.Equipment[j]
				return obj
			}()) != nil {
				if GET_OBJ_VNUM(obj) == ovnum {
					return obj
				}
			}
		}
	} else {
		stdio.Snprintf(&tmpname[0], int(2048), "%s", name)
		if (func() int {
			number = get_number(&tmp)
			return number
		}()) == 0 {
			return nil
		}
		for j = 0; j < NUM_WEARS && n <= number; j++ {
			if (func() *obj_data {
				obj = ch.Equipment[j]
				return obj
			}()) != nil {
				if isname(tmp, obj.Name) != 0 {
					if func() int {
						p := &n
						*p++
						return *p
					}() == number {
						return obj
					}
				}
			}
		}
	}
	return nil
}
func find_eq_pos_script(arg *byte) int {
	var i int
	type eq_pos_list struct {
		Pos   *byte
		Where int
	}
	var eq_pos [25]eq_pos_list = [25]eq_pos_list{{Pos: libc.CString("hold"), Where: WEAR_WIELD2}, {Pos: libc.CString("held"), Where: WEAR_WIELD2}, {Pos: libc.CString("light"), Where: WEAR_WIELD2}, {Pos: libc.CString("wield"), Where: WEAR_WIELD1}, {Pos: libc.CString("rfinger"), Where: WEAR_FINGER_R}, {Pos: libc.CString("lfinger"), Where: WEAR_FINGER_L}, {Pos: libc.CString("neck1"), Where: WEAR_NECK_1}, {Pos: libc.CString("neck2"), Where: WEAR_NECK_2}, {Pos: libc.CString("body"), Where: WEAR_BODY}, {Pos: libc.CString("head"), Where: WEAR_HEAD}, {Pos: libc.CString("legs"), Where: WEAR_LEGS}, {Pos: libc.CString("feet"), Where: WEAR_FEET}, {Pos: libc.CString("hands"), Where: WEAR_HANDS}, {Pos: libc.CString("arms"), Where: WEAR_ARMS}, {Pos: libc.CString("shield"), Where: WEAR_WIELD2}, {Pos: libc.CString("about"), Where: WEAR_ABOUT}, {Pos: libc.CString("waist"), Where: WEAR_WAIST}, {Pos: libc.CString("rwrist"), Where: WEAR_WRIST_R}, {Pos: libc.CString("lwrist"), Where: WEAR_WRIST_L}, {Pos: libc.CString("backpack"), Where: WEAR_BACKPACK}, {Pos: libc.CString("rear"), Where: WEAR_EAR_R}, {Pos: libc.CString("lear"), Where: WEAR_EAR_L}, {Pos: libc.CString("shoulders"), Where: WEAR_SH}, {Pos: libc.CString("scouter"), Where: WEAR_EYE}, {Pos: libc.CString("none"), Where: -1}}
	if is_number(arg) != 0 && (func() int {
		i = libc.Atoi(libc.GoString(arg))
		return i
	}()) >= 0 && i < NUM_WEARS {
		return i
	}
	for i = 0; eq_pos[i].Where != -1; i++ {
		if C.strcasecmp(eq_pos[i].Pos, arg) == 0 {
			return eq_pos[i].Where
		}
	}
	return -1
}
func can_wear_on_pos(obj *obj_data, pos int) int {
	switch pos {
	case WEAR_WIELD1:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_WIELD)))
	case WEAR_WIELD2:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_TAKE)))
	case WEAR_FINGER_R:
		fallthrough
	case WEAR_FINGER_L:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_FINGER)))
	case WEAR_NECK_1:
		fallthrough
	case WEAR_NECK_2:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_NECK)))
	case WEAR_BODY:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_BODY)))
	case WEAR_HEAD:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_HEAD)))
	case WEAR_LEGS:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_LEGS)))
	case WEAR_FEET:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_FEET)))
	case WEAR_HANDS:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_HANDS)))
	case WEAR_ARMS:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_ARMS)))
	case WEAR_ABOUT:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_ABOUT)))
	case WEAR_WAIST:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_WAIST)))
	case WEAR_WRIST_R:
		fallthrough
	case WEAR_WRIST_L:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_WRIST)))
	case WEAR_BACKPACK:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_PACK)))
	case WEAR_EAR_R:
		fallthrough
	case WEAR_EAR_L:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_EAR)))
	case WEAR_SH:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_SH)))
	case WEAR_EYE:
		return int(libc.BoolToInt(OBJWEAR_FLAGGED(obj, ITEM_WEAR_EYE)))
	default:
		return FALSE
	}
}
func find_char(n int) *char_data {
	if n >= ROOM_ID_BASE {
		return nil
	}
	return find_char_by_uid_in_lookup_table(n)
}
func find_obj(n int) *obj_data {
	if n < OBJ_ID_BASE {
		return nil
	}
	return find_obj_by_uid_in_lookup_table(n)
}
func find_room(n int) *room_data {
	var rnum room_rnum
	n -= ROOM_ID_BASE
	if n < 0 {
		return nil
	}
	rnum = real_room(room_vnum(n))
	if rnum != room_rnum(-1) {
		return (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))
	}
	return nil
}
func get_char(name *byte) *char_data {
	var i *char_data
	if *name == UID_CHAR {
		i = find_char(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
		if i != nil && valid_dg_target(i, 1<<0) != 0 {
			return i
		}
	} else {
		for i = character_list; i != nil; i = i.Next {
			if isname(name, i.Name) != 0 && valid_dg_target(i, 1<<0) != 0 {
				return i
			}
		}
	}
	return nil
}
func get_char_near_obj(obj *obj_data, name *byte) *char_data {
	var ch *char_data
	if *name == UID_CHAR {
		ch = find_char(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
		if ch != nil && valid_dg_target(ch, 1<<0) != 0 {
			return ch
		}
	} else {
		var num room_rnum
		if (func() room_rnum {
			num = obj_room(obj)
			return num
		}()) != room_rnum(-1) {
			for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(num)))).People; ch != nil; ch = ch.Next_in_room {
				if isname(name, ch.Name) != 0 && valid_dg_target(ch, 1<<0) != 0 {
					return ch
				}
			}
		}
	}
	return nil
}
func get_char_in_room(room *room_data, name *byte) *char_data {
	var ch *char_data
	if *name == UID_CHAR {
		ch = find_char(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
		if ch != nil && valid_dg_target(ch, 1<<0) != 0 {
			return ch
		}
	} else {
		for ch = room.People; ch != nil; ch = ch.Next_in_room {
			if isname(name, ch.Name) != 0 && valid_dg_target(ch, 1<<0) != 0 {
				return ch
			}
		}
	}
	return nil
}
func get_obj_near_obj(obj *obj_data, name *byte) *obj_data {
	var (
		i  *obj_data = nil
		ch *char_data
		rm int
		id int
	)
	if C.strcasecmp(name, libc.CString("self")) == 0 || C.strcasecmp(name, libc.CString("me")) == 0 {
		return obj
	}
	if obj.Contains != nil && (func() *obj_data {
		i = get_obj_in_list(name, obj.Contains)
		return i
	}()) != nil {
		return i
	}
	if obj.In_obj != nil {
		if *name == UID_CHAR {
			id = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1))))
			if id == int(obj.In_obj.Id) {
				return obj.In_obj
			}
		} else if isname(name, obj.In_obj.Name) != 0 {
			return obj.In_obj
		}
	} else if obj.Worn_by != nil && (func() *obj_data {
		i = get_object_in_equip(obj.Worn_by, name)
		return i
	}()) != nil {
		return i
	} else if obj.Carried_by != nil && (func() *obj_data {
		i = get_obj_in_list(name, obj.Carried_by.Carrying)
		return i
	}()) != nil {
		return i
	} else if (func() int {
		rm = int(obj_room(obj))
		return rm
	}()) != int(-1) {
		if (func() *obj_data {
			i = get_obj_in_list(name, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).Contents)
			return i
		}()) != nil {
			return i
		}
		for ch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).People; ch != nil; ch = ch.Next_in_room {
			if (func() *obj_data {
				i = get_object_in_equip(ch, name)
				return i
			}()) != nil {
				return i
			}
		}
	}
	return nil
}
func get_obj(name *byte) *obj_data {
	var obj *obj_data
	if *name == UID_CHAR {
		return find_obj(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
	} else {
		for obj = object_list; obj != nil; obj = obj.Next {
			if isname(name, obj.Name) != 0 {
				return obj
			}
		}
	}
	return nil
}
func get_room(name *byte) *room_data {
	var nr room_rnum
	if *name == UID_CHAR {
		return find_room(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
	} else if (func() room_rnum {
		nr = real_room(room_vnum(libc.Atoi(libc.GoString(name))))
		return nr
	}()) == room_rnum(-1) {
		return nil
	} else {
		return (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))
	}
}
func get_char_by_obj(obj *obj_data, name *byte) *char_data {
	var ch *char_data
	if *name == UID_CHAR {
		ch = find_char(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
		if ch != nil && valid_dg_target(ch, 1<<0) != 0 {
			return ch
		}
	} else {
		if obj.Carried_by != nil && isname(name, obj.Carried_by.Name) != 0 && valid_dg_target(obj.Carried_by, 1<<0) != 0 {
			return obj.Carried_by
		}
		if obj.Worn_by != nil && isname(name, obj.Worn_by.Name) != 0 && valid_dg_target(obj.Worn_by, 1<<0) != 0 {
			return obj.Worn_by
		}
		for ch = character_list; ch != nil; ch = ch.Next {
			if isname(name, ch.Name) != 0 && valid_dg_target(ch, 1<<0) != 0 {
				return ch
			}
		}
	}
	return nil
}
func get_char_by_room(room *room_data, name *byte) *char_data {
	var ch *char_data
	if *name == UID_CHAR {
		ch = find_char(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
		if ch != nil && valid_dg_target(ch, 1<<0) != 0 {
			return ch
		}
	} else {
		for ch = room.People; ch != nil; ch = ch.Next_in_room {
			if isname(name, ch.Name) != 0 && valid_dg_target(ch, 1<<0) != 0 {
				return ch
			}
		}
		for ch = character_list; ch != nil; ch = ch.Next {
			if isname(name, ch.Name) != 0 && valid_dg_target(ch, 1<<0) != 0 {
				return ch
			}
		}
	}
	return nil
}
func get_obj_by_obj(obj *obj_data, name *byte) *obj_data {
	var (
		i  *obj_data = nil
		rm int
	)
	if *name == UID_CHAR {
		return find_obj(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
	}
	if C.strcasecmp(name, libc.CString("self")) == 0 || C.strcasecmp(name, libc.CString("me")) == 0 {
		return obj
	}
	if obj.Contains != nil && (func() *obj_data {
		i = get_obj_in_list(name, obj.Contains)
		return i
	}()) != nil {
		return i
	}
	if obj.In_obj != nil && isname(name, obj.In_obj.Name) != 0 {
		return obj.In_obj
	}
	if obj.Worn_by != nil && (func() *obj_data {
		i = get_object_in_equip(obj.Worn_by, name)
		return i
	}()) != nil {
		return i
	}
	if obj.Carried_by != nil && (func() *obj_data {
		i = get_obj_in_list(name, obj.Carried_by.Carrying)
		return i
	}()) != nil {
		return i
	}
	if (func() int {
		rm = int(obj_room(obj))
		return rm
	}()) != int(-1) && (func() *obj_data {
		i = get_obj_in_list(name, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm)))).Contents)
		return i
	}()) != nil {
		return i
	}
	return get_obj(name)
}
func get_obj_in_room(room *room_data, name *byte) *obj_data {
	var (
		obj *obj_data
		id  int
	)
	if *name == UID_CHAR {
		id = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1))))
		for obj = room.Contents; obj != nil; obj = obj.Next_content {
			if id == int(obj.Id) {
				return obj
			}
		}
	} else {
		for obj = room.Contents; obj != nil; obj = obj.Next_content {
			if isname(name, obj.Name) != 0 {
				return obj
			}
		}
	}
	return nil
}
func get_obj_by_room(room *room_data, name *byte) *obj_data {
	var obj *obj_data
	if *name == UID_CHAR {
		return find_obj(libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1)))))
	}
	for obj = room.Contents; obj != nil; obj = obj.Next_content {
		if isname(name, obj.Name) != 0 {
			return obj
		}
	}
	for obj = object_list; obj != nil; obj = obj.Next {
		if isname(name, obj.Name) != 0 {
			return obj
		}
	}
	return nil
}
func script_trigger_check() {
	var (
		ch   *char_data
		obj  *obj_data
		room *room_data = nil
		nr   int
		sc   *script_data
	)
	for ch = character_list; ch != nil; ch = ch.Next {
		if ch.Script != nil {
			sc = ch.Script
			if (sc.Types&(1<<1)) != 0 && (is_empty((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone) == 0 || (sc.Types&(1<<0)) != 0) {
				random_mtrigger(ch)
			}
		}
	}
	for obj = object_list; obj != nil; obj = obj.Next {
		if obj.Script != nil {
			sc = obj.Script
			if (sc.Types & (1 << 1)) != 0 {
				random_otrigger(obj)
			}
		}
	}
	for nr = 0; nr <= int(top_of_world); nr++ {
		if ((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Script != nil {
			room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))
			sc = room.Script
			if (sc.Types&(1<<1)) != 0 && (is_empty(room.Zone) == 0 || (sc.Types&(1<<0)) != 0) {
				random_wtrigger(room)
			}
		}
	}
}
func check_timed_triggers() {
	var (
		ch   *char_data
		obj  *obj_data
		room *room_data = nil
		nr   int
		sc   *script_data
	)
	for ch = character_list; ch != nil; ch = ch.Next {
		if ch.Script != nil {
			sc = ch.Script
			if (sc.Types&(1<<19)) != 0 && (is_empty((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone) == 0 || (sc.Types&(1<<0)) != 0) {
				time_mtrigger(ch)
			}
		}
	}
	for obj = object_list; obj != nil; obj = obj.Next {
		if obj.Script != nil {
			sc = obj.Script
			if (sc.Types & (1 << 19)) != 0 {
				time_otrigger(obj)
			}
		}
	}
	for nr = 0; nr <= int(top_of_world); nr++ {
		if ((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Script != nil {
			room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))
			sc = room.Script
			if (sc.Types&(1<<19)) != 0 && (is_empty(room.Zone) == 0 || (sc.Types&(1<<0)) != 0) {
				time_wtrigger(room)
			}
		}
	}
}
func trig_wait_event(event_obj unsafe.Pointer) int {
	var (
		wait_event_obj *wait_event_data = (*wait_event_data)(event_obj)
		trig           *trig_data
		gohere         unsafe.Pointer
		type_          int
	)
	trig = wait_event_obj.Trigger
	gohere = wait_event_obj.Gohere
	type_ = wait_event_obj.Type
	libc.Free(unsafe.Pointer(wait_event_obj))
	trig.Wait_event = nil
	{
		var found int = FALSE
		if type_ == MOB_TRIGGER {
			var tch *char_data
			for tch = character_list; tch != nil && found == 0; tch = tch.Next {
				if tch == (*char_data)(gohere) {
					found = TRUE
				}
			}
		} else if type_ == OBJ_TRIGGER {
			var obj *obj_data
			for obj = object_list; obj != nil && found == 0; obj = obj.Next {
				if obj == (*obj_data)(gohere) {
					found = TRUE
				}
			}
		} else {
			var i room_rnum
			for i = 0; i < top_of_world && found == 0; i++ {
				if (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))) == (*room_data)(gohere) {
					found = TRUE
				}
			}
		}
		if found == 0 {
			basic_mud_log(libc.CString("Trigger restarted on unknown entity. Vnum: %d"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
			basic_mud_log(libc.CString("Type: %s trigger"), func() string {
				if type_ == MOB_TRIGGER {
					return "Mob"
				}
				if type_ == OBJ_TRIGGER {
					return "Obj"
				}
				return "Room"
			}())
			basic_mud_log(libc.CString("attached %d places"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Number)
			script_log(libc.CString("Trigger restart attempt on unknown entity."))
			return 0
		}
	}
	script_driver(unsafe.Pointer(&gohere), trig, type_, TRIG_RESTART)
	return 0
}
func do_stat_trigger(ch *char_data, trig *trig_data) {
	var (
		cmd_list *cmdlist_element
		sb       [64936]byte
		buf      [64936]byte
		len_     int = 0
	)
	if trig == nil {
		basic_mud_log(libc.CString("SYSERR: NULL trigger passed to do_stat_trigger."))
		return
	}
	len_ += stdio.Snprintf(&sb[0], int(64936), "Name: '@y%s@n',  VNum: [@g%5d@n], RNum: [%5d]\r\n", trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, trig.Nr)
	if trig.Attach_type == OBJ_TRIGGER {
		len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "Trigger Intended Assignment: Objects\r\n")
		sprintbit(bitvector_t(trig.Trigger_type), otrig_types[:], &buf[0], uint64(64936))
	} else if trig.Attach_type == WLD_TRIGGER {
		len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "Trigger Intended Assignment: Rooms\r\n")
		sprintbit(bitvector_t(trig.Trigger_type), wtrig_types[:], &buf[0], uint64(64936))
	} else {
		len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "Trigger Intended Assignment: Mobiles\r\n")
		sprintbit(bitvector_t(trig.Trigger_type), trig_types[:], &buf[0], uint64(64936))
	}
	len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "Trigger Type: %s, Numeric Arg: %d, Arg list: %s\r\n", &buf[0], trig.Narg, func() *byte {
		if trig.Arglist != nil && *trig.Arglist != 0 {
			return trig.Arglist
		}
		return libc.CString("None")
	}())
	len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "Commands:\r\n")
	cmd_list = trig.Cmdlist
	for cmd_list != nil {
		if cmd_list.Cmd != nil {
			len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "%s\r\n", cmd_list.Cmd)
		}
		if len_ > int(MAX_STRING_LENGTH-80) {
			len_ += stdio.Snprintf(&sb[len_], int(64936-uintptr(len_)), "*** Overflow - script too long! ***\r\n")
			break
		}
		cmd_list = cmd_list.Next
	}
	page_string(ch.Desc, &sb[0], 1)
}
func find_uid_name(uid *byte, name *byte, nlen uint64) {
	var (
		ch  *char_data
		obj *obj_data
	)
	if (func() *char_data {
		ch = get_char(uid)
		return ch
	}()) != nil {
		stdio.Snprintf(name, int(nlen), "%s", ch.Name)
	} else if (func() *obj_data {
		obj = get_obj(uid)
		return obj
	}()) != nil {
		stdio.Snprintf(name, int(nlen), "%s", obj.Name)
	} else {
		stdio.Snprintf(name, int(nlen), "uid = %s, (not found)", (*byte)(unsafe.Add(unsafe.Pointer(uid), 1)))
	}
}
func script_stat(ch *char_data, sc *script_data) {
	var (
		tv      *trig_var_data
		t       *trig_data
		name    [2048]byte
		namebuf [512]byte
		buf1    [64936]byte
	)
	send_to_char(ch, libc.CString("Global Variables: %s\r\n"), func() string {
		if sc.Global_vars != nil {
			return ""
		}
		return "None"
	}())
	send_to_char(ch, libc.CString("Global context: %ld\r\n"), sc.Context)
	for tv = sc.Global_vars; tv != nil; tv = tv.Next {
		stdio.Snprintf(&namebuf[0], int(512), "%s:%ld", tv.Name, tv.Context)
		if *tv.Value == UID_CHAR {
			find_uid_name(tv.Value, &name[0], uint64(2048))
			send_to_char(ch, libc.CString("    %15s:  %s\r\n"), func() *byte {
				if tv.Context != 0 {
					return &namebuf[0]
				}
				return tv.Name
			}(), &name[0])
		} else {
			send_to_char(ch, libc.CString("    %15s:  %s\r\n"), func() *byte {
				if tv.Context != 0 {
					return &namebuf[0]
				}
				return tv.Name
			}(), tv.Value)
		}
	}
	for t = sc.Trig_list; t != nil; t = t.Next {
		send_to_char(ch, libc.CString("\r\n  Trigger: @y%s@n, VNum: [@y%5d@n], RNum: [%5d]\r\n"), t.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(t.Nr)))).Vnum, t.Nr)
		if t.Attach_type == OBJ_TRIGGER {
			send_to_char(ch, libc.CString("  Trigger Intended Assignment: Objects\r\n"))
			sprintbit(bitvector_t(t.Trigger_type), otrig_types[:], &buf1[0], uint64(64936))
		} else if t.Attach_type == WLD_TRIGGER {
			send_to_char(ch, libc.CString("  Trigger Intended Assignment: Rooms\r\n"))
			sprintbit(bitvector_t(t.Trigger_type), wtrig_types[:], &buf1[0], uint64(64936))
		} else {
			send_to_char(ch, libc.CString("  Trigger Intended Assignment: Mobiles\r\n"))
			sprintbit(bitvector_t(t.Trigger_type), trig_types[:], &buf1[0], uint64(64936))
		}
		send_to_char(ch, libc.CString("  Trigger Type: %s, Numeric Arg: %d, Arg list: %s\r\n"), &buf1[0], t.Narg, func() *byte {
			if t.Arglist != nil && *t.Arglist != 0 {
				return t.Arglist
			}
			return libc.CString("None")
		}())
		if t.Wait_event != nil {
			send_to_char(ch, libc.CString("    Wait: %ld, Current line: %s\r\n"), event_time(t.Wait_event), func() *byte {
				if t.Curr_state != nil {
					return t.Curr_state.Cmd
				}
				return libc.CString("End of Script")
			}())
			send_to_char(ch, libc.CString("  Variables: %s\r\n"), func() string {
				if t.Var_list != nil {
					return ""
				}
				return "None"
			}())
			for tv = t.Var_list; tv != nil; tv = tv.Next {
				if *tv.Value == UID_CHAR {
					find_uid_name(tv.Value, &name[0], uint64(2048))
					send_to_char(ch, libc.CString("    %15s:  %s\r\n"), tv.Name, &name[0])
				} else {
					send_to_char(ch, libc.CString("    %15s:  %s\r\n"), tv.Name, tv.Value)
				}
			}
		}
	}
}
func do_sstat_room(ch *char_data, rm *room_data) {
	send_to_char(ch, libc.CString("Triggers:\r\n"))
	if rm.Script == nil {
		send_to_char(ch, libc.CString("  None.\r\n"))
		return
	}
	script_stat(ch, rm.Script)
}
func do_sstat_object(ch *char_data, j *obj_data) {
	send_to_char(ch, libc.CString("Triggers:\r\n"))
	if j.Script == nil {
		send_to_char(ch, libc.CString("  None.\r\n"))
		return
	}
	script_stat(ch, j.Script)
}
func do_sstat_character(ch *char_data, k *char_data) {
	send_to_char(ch, libc.CString("Script information:\r\n"))
	if k.Script == nil {
		send_to_char(ch, libc.CString("  None.\r\n"))
		return
	}
	script_stat(ch, k.Script)
}
func add_trigger(sc *script_data, t *trig_data, loc int) {
	var (
		i *trig_data
		n int
	)
	for func() *trig_data {
		n = loc
		return func() *trig_data {
			i = sc.Trig_list
			return i
		}()
	}(); i != nil && i.Next != nil && n != 0; func() *trig_data {
		n--
		return func() *trig_data {
			i = i.Next
			return i
		}()
	}() {
	}
	if loc == 0 {
		t.Next = sc.Trig_list
		sc.Trig_list = t
	} else if i == nil {
		sc.Trig_list = t
	} else {
		t.Next = i.Next
		i.Next = t
	}
	sc.Types |= t.Trigger_type
	t.Next_in_world = trigger_list
	trigger_list = t
}
func do_attach(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim    *char_data
		object    *obj_data
		room      *room_data
		trig      *trig_data
		targ_name [2048]byte
		trig_name [2048]byte
		loc_name  [2048]byte
		arg       [2048]byte
		loc       int
		tn        int
		rn        int
		num_arg   int
		rnum      room_rnum
	)
	argument = two_arguments(argument, &arg[0], &trig_name[0])
	two_arguments(argument, &targ_name[0], &loc_name[0])
	if arg[0] == 0 || targ_name[0] == 0 || trig_name[0] == 0 {
		send_to_char(ch, libc.CString("Usage: attach { mob | obj | room } { trigger } { name } [ location ]\r\n"))
		return
	}
	num_arg = libc.Atoi(libc.GoString(&targ_name[0]))
	tn = libc.Atoi(libc.GoString(&trig_name[0]))
	if (loc_name[0]) != 0 {
		loc = libc.Atoi(libc.GoString(&loc_name[0]))
	} else {
		loc = -1
	}
	if is_abbrev(&arg[0], libc.CString("mobile")) != 0 || is_abbrev(&arg[0], libc.CString("mtr")) != 0 {
		victim = get_char_vis(ch, &targ_name[0], nil, 1<<1)
		if victim == nil {
			for victim = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; victim != nil; victim = victim.Next_in_room {
				if GET_MOB_VNUM(victim) == mob_vnum(num_arg) {
					break
				}
			}
			if victim == nil {
				send_to_char(ch, libc.CString("That mob does not exist.\r\n"))
				return
			}
		}
		if !IS_NPC(victim) {
			send_to_char(ch, libc.CString("Players can't have scripts.\r\n"))
			return
		}
		if can_edit_zone(ch, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone) == 0 {
			send_to_char(ch, libc.CString("You can only attach triggers in your own zone\r\n"))
			return
		}
		rn = int(real_trigger(trig_vnum(tn)))
		if rn == int(-1) || (func() *trig_data {
			trig = read_trigger(rn)
			return trig
		}()) == nil {
			send_to_char(ch, libc.CString("That trigger does not exist.\r\n"))
			return
		}
		if victim.Script == nil {
			victim.Script = new(script_data)
		}
		add_trigger(victim.Script, trig, loc)
		send_to_char(ch, libc.CString("Trigger %d (%s) attached to %s [%d].\r\n"), tn, trig.Name, victim.Short_descr, GET_MOB_VNUM(victim))
	} else if is_abbrev(&arg[0], libc.CString("object")) != 0 || is_abbrev(&arg[0], libc.CString("otr")) != 0 {
		object = get_obj_vis(ch, &targ_name[0], nil)
		if object == nil {
			for object = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; object != nil; object = object.Next_content {
				if GET_OBJ_VNUM(object) == obj_vnum(num_arg) {
					break
				}
			}
			if object == nil {
				for object = ch.Carrying; object != nil; object = object.Next_content {
					if GET_OBJ_VNUM(object) == obj_vnum(num_arg) {
						break
					}
				}
				if object == nil {
					send_to_char(ch, libc.CString("That object does not exist.\r\n"))
					return
				}
			}
		}
		if can_edit_zone(ch, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone) == 0 {
			send_to_char(ch, libc.CString("You can only attach triggers in your own zone\r\n"))
			return
		}
		rn = int(real_trigger(trig_vnum(tn)))
		if rn == int(-1) || (func() *trig_data {
			trig = read_trigger(rn)
			return trig
		}()) == nil {
			send_to_char(ch, libc.CString("That trigger does not exist.\r\n"))
			return
		}
		if object.Script == nil {
			object.Script = new(script_data)
		}
		add_trigger(object.Script, trig, loc)
		send_to_char(ch, libc.CString("Trigger %d (%s) attached to %s [%d].\r\n"), tn, trig.Name, func() *byte {
			if object.Short_description != nil {
				return object.Short_description
			}
			return object.Name
		}(), GET_OBJ_VNUM(object))
	} else if is_abbrev(&arg[0], libc.CString("room")) != 0 || is_abbrev(&arg[0], libc.CString("wtr")) != 0 {
		if C.strchr(&targ_name[0], '.') != nil {
			rnum = ch.In_room
		} else if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(targ_name[0]))))) & int(uint16(int16(_ISdigit)))) != 0 {
			rnum = find_target_room(ch, &targ_name[0])
		} else {
			rnum = -1
		}
		if rnum == room_rnum(-1) {
			send_to_char(ch, libc.CString("You need to supply a room number or . for current room.\r\n"))
			return
		}
		if can_edit_zone(ch, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Zone) == 0 {
			send_to_char(ch, libc.CString("You can only attach triggers in your own zone\r\n"))
			return
		}
		rn = int(real_trigger(trig_vnum(tn)))
		if rn == int(-1) || (func() *trig_data {
			trig = read_trigger(rn)
			return trig
		}()) == nil {
			send_to_char(ch, libc.CString("That trigger does not exist.\r\n"))
			return
		}
		room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))
		if room.Script == nil {
			room.Script = new(script_data)
		}
		add_trigger(room.Script, trig, loc)
		send_to_char(ch, libc.CString("Trigger %d (%s) attached to room %d.\r\n"), tn, trig.Name, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Number)
	} else {
		send_to_char(ch, libc.CString("Please specify 'mob', 'obj', or 'room'.\r\n"))
	}
}
func remove_trigger(sc *script_data, name *byte) int {
	var (
		i       *trig_data
		j       *trig_data
		num     int = 0
		string_ int = FALSE
		n       int
		cname   *byte
	)
	if sc == nil {
		return 0
	}
	if (func() *byte {
		cname = C.strstr(name, libc.CString("."))
		return cname
	}()) != nil || (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*name)))))&int(uint16(int16(_ISdigit)))) == 0 {
		string_ = TRUE
		if cname != nil {
			*cname = '\x00'
			num = libc.Atoi(libc.GoString(name))
			name = func() *byte {
				p := &cname
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()
		}
	} else {
		num = libc.Atoi(libc.GoString(name))
	}
	for func() *trig_data {
		n = 0
		j = nil
		return func() *trig_data {
			i = sc.Trig_list
			return i
		}()
	}(); i != nil; func() *trig_data {
		j = i
		return func() *trig_data {
			i = i.Next
			return i
		}()
	}() {
		if string_ != 0 {
			if isname(name, i.Name) != 0 {
				if func() int {
					p := &n
					*p++
					return *p
				}() >= num {
					break
				}
			}
		} else if func() int {
			p := &n
			*p++
			return *p
		}() >= num {
			break
		} else if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i.Nr)))).Vnum == mob_vnum(num) {
			break
		}
	}
	if i != nil {
		if j != nil {
			j.Next = i.Next
			extract_trigger(i)
		} else {
			sc.Trig_list = i.Next
			extract_trigger(i)
		}
		sc.Types = 0
		for i = sc.Trig_list; i != nil; i = i.Next {
			sc.Types |= i.Trigger_type
		}
		return 1
	} else {
		return 0
	}
}
func do_detach(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim  *char_data = nil
		object  *obj_data  = nil
		room    *room_data
		arg1    [2048]byte
		arg2    [2048]byte
		arg3    [2048]byte
		trigger *byte = nil
		num_arg int
	)
	argument = two_arguments(argument, &arg1[0], &arg2[0])
	one_argument(argument, &arg3[0])
	if arg1[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Usage: detach [ mob | object | room ] { target } { trigger | 'all' }\r\n"))
		return
	}
	num_arg = libc.Atoi(libc.GoString(&arg2[0]))
	if C.strcasecmp(&arg1[0], libc.CString("room")) == 0 || C.strcasecmp(&arg1[0], libc.CString("wtr")) == 0 {
		room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))
		if can_edit_zone(ch, room.Zone) == 0 {
			send_to_char(ch, libc.CString("You can only detach triggers in your own zone\r\n"))
			return
		}
		if room.Script == nil {
			send_to_char(ch, libc.CString("This room does not have any triggers.\r\n"))
		} else if C.strcasecmp(&arg2[0], libc.CString("all")) == 0 {
			extract_script(unsafe.Pointer(room), WLD_TRIGGER)
			send_to_char(ch, libc.CString("All triggers removed from room.\r\n"))
		} else if remove_trigger(room.Script, &arg2[0]) != 0 {
			send_to_char(ch, libc.CString("Trigger removed.\r\n"))
			if room.Script.Trig_list == nil {
				extract_script(unsafe.Pointer(room), WLD_TRIGGER)
			}
		} else {
			send_to_char(ch, libc.CString("That trigger was not found.\r\n"))
		}
	} else {
		if is_abbrev(&arg1[0], libc.CString("mobile")) != 0 || C.strcasecmp(&arg1[0], libc.CString("mtr")) == 0 {
			victim = get_char_vis(ch, &arg2[0], nil, 1<<1)
			if victim == nil {
				for victim = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; victim != nil; victim = victim.Next_in_room {
					if GET_MOB_VNUM(victim) == mob_vnum(num_arg) {
						break
					}
				}
				if victim == nil {
					send_to_char(ch, libc.CString("No such mobile around.\r\n"))
					return
				}
			}
			if arg3 == nil || arg3[0] == 0 {
				send_to_char(ch, libc.CString("You must specify a trigger to remove.\r\n"))
			} else {
				trigger = &arg3[0]
			}
		} else if is_abbrev(&arg1[0], libc.CString("object")) != 0 || C.strcasecmp(&arg1[0], libc.CString("otr")) == 0 {
			object = get_obj_vis(ch, &arg2[0], nil)
			if object == nil {
				for object = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; object != nil; object = object.Next_content {
					if GET_OBJ_VNUM(object) == obj_vnum(num_arg) {
						break
					}
				}
				if object == nil {
					for object = ch.Carrying; object != nil; object = object.Next_content {
						if GET_OBJ_VNUM(object) == obj_vnum(num_arg) {
							break
						}
					}
					if object == nil {
						send_to_char(ch, libc.CString("No such object around.\r\n"))
						return
					}
				}
			}
			if arg3 == nil || arg3[0] == 0 {
				send_to_char(ch, libc.CString("You must specify a trigger to remove.\r\n"))
			} else {
				trigger = &arg3[0]
			}
		} else {
			if (func() *obj_data {
				object = get_obj_in_equip_vis(ch, &arg1[0], nil, ch.Equipment[:])
				return object
			}()) != nil {
			} else if (func() *obj_data {
				object = get_obj_in_list_vis(ch, &arg1[0], nil, ch.Carrying)
				return object
			}()) != nil {
			} else if (func() *char_data {
				victim = get_char_room_vis(ch, &arg1[0], nil)
				return victim
			}()) != nil {
			} else if (func() *obj_data {
				object = get_obj_in_list_vis(ch, &arg1[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
				return object
			}()) != nil {
			} else if (func() *char_data {
				victim = get_char_vis(ch, &arg1[0], nil, 1<<1)
				return victim
			}()) != nil {
			} else if (func() *obj_data {
				object = get_obj_vis(ch, &arg1[0], nil)
				return object
			}()) != nil {
			} else {
				send_to_char(ch, libc.CString("Nothing around by that name.\r\n"))
			}
			trigger = &arg2[0]
		}
		if victim != nil {
			if !IS_NPC(victim) {
				send_to_char(ch, libc.CString("Players don't have triggers.\r\n"))
			} else if victim.Script == nil {
				send_to_char(ch, libc.CString("That mob doesn't have any triggers.\r\n"))
			} else if can_edit_zone(ch, real_zone_by_thing(room_vnum(GET_MOB_VNUM(victim)))) == 0 {
				send_to_char(ch, libc.CString("You can only detach triggers in your own zone\r\n"))
				return
			} else if trigger != nil && C.strcasecmp(trigger, libc.CString("all")) == 0 {
				extract_script(unsafe.Pointer(victim), MOB_TRIGGER)
				send_to_char(ch, libc.CString("All triggers removed from %s.\r\n"), victim.Short_descr)
			} else if trigger != nil && remove_trigger(victim.Script, trigger) != 0 {
				send_to_char(ch, libc.CString("Trigger removed.\r\n"))
				if victim.Script.Trig_list == nil {
					extract_script(unsafe.Pointer(victim), MOB_TRIGGER)
				}
			} else {
				send_to_char(ch, libc.CString("That trigger was not found.\r\n"))
			}
		} else if object != nil {
			if object.Script == nil {
				send_to_char(ch, libc.CString("That object doesn't have any triggers.\r\n"))
			} else if can_edit_zone(ch, real_zone_by_thing(room_vnum(GET_OBJ_VNUM(object)))) == 0 {
				send_to_char(ch, libc.CString("You can only detach triggers in your own zone\r\n"))
				return
			} else if trigger != nil && C.strcasecmp(trigger, libc.CString("all")) == 0 {
				extract_script(unsafe.Pointer(object), OBJ_TRIGGER)
				send_to_char(ch, libc.CString("All triggers removed from %s.\r\n"), func() *byte {
					if object.Short_description != nil {
						return object.Short_description
					}
					return object.Name
				}())
			} else if remove_trigger(object.Script, trigger) != 0 {
				send_to_char(ch, libc.CString("Trigger removed.\r\n"))
				if object.Script.Trig_list == nil {
					extract_script(unsafe.Pointer(object), OBJ_TRIGGER)
				}
			} else {
				send_to_char(ch, libc.CString("That trigger was not found.\r\n"))
			}
		}
	}
}
func script_vlog(format *byte, args libc.ArgList) {
	var (
		output [64936]byte
		i      *descriptor_data
	)
	stdio.Snprintf(&output[0], int(64936), "SCRIPT ERR: %s", format)
	basic_mud_vlog(&output[0], args)
	C.strcpy(&output[0], libc.CString("[ "))
	stdio.Vsnprintf(&output[2], int(64936-6), libc.GoString(format), args)
	C.strcat(&output[0], libc.CString(" ]\r\n"))
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING || IS_NPC(i.Character) {
			continue
		}
		if i.Character.Admlevel < ADMLVL_BUILDER {
			continue
		}
		if PLR_FLAGGED(i.Character, PLR_WRITING) {
			continue
		}
		if NRM > (func() int {
			if PRF_FLAGGED(i.Character, PRF_LOG1) {
				return 1
			}
			return 0
		}())+(func() int {
			if PRF_FLAGGED(i.Character, PRF_LOG2) {
				return 2
			}
			return 0
		}()) {
			continue
		}
		send_to_char(i.Character, libc.CString("@g%s@n"), &output[0])
	}
}
func script_log(format *byte, _rest ...interface{}) {
	var args libc.ArgList
	args.Start(format, _rest)
	script_vlog(format, args)
	args.End()
}
func is_num(arg *byte) int {
	if *arg == '\x00' {
		return FALSE
	}
	if *arg == '+' || *arg == '-' {
		arg = (*byte)(unsafe.Add(unsafe.Pointer(arg), 1))
	}
	for ; *arg != '\x00'; arg = (*byte)(unsafe.Add(unsafe.Pointer(arg), 1)) {
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*arg))))) & int(uint16(int16(_ISdigit)))) == 0 {
			return FALSE
		}
	}
	return TRUE
}
func eval_op(op *byte, lhs *byte, rhs *byte, result *byte, gohere unsafe.Pointer, sc *script_data, trig *trig_data) {
	var (
		p *uint8
		n int
	)
	for *lhs != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*lhs)))))&int(uint16(int16(_ISspace)))) != 0 {
		lhs = (*byte)(unsafe.Add(unsafe.Pointer(lhs), 1))
	}
	for *rhs != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*rhs)))))&int(uint16(int16(_ISspace)))) != 0 {
		rhs = (*byte)(unsafe.Add(unsafe.Pointer(rhs), 1))
	}
	for p = (*uint8)(unsafe.Pointer(lhs)); int(*p) != 0; p = (*uint8)(unsafe.Add(unsafe.Pointer(p), 1)) {
	}
	for p = (*uint8)(unsafe.Add(unsafe.Pointer(p), -1)); (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0 && uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(p)))) > uintptr(unsafe.Pointer(lhs)); *func() *uint8 {
		p := &p
		x := *p
		*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), -1))
		return x
	}() = '\x00' {
	}
	for p = (*uint8)(unsafe.Pointer(rhs)); int(*p) != 0; p = (*uint8)(unsafe.Add(unsafe.Pointer(p), 1)) {
	}
	for p = (*uint8)(unsafe.Add(unsafe.Pointer(p), -1)); (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0 && uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(p)))) > uintptr(unsafe.Pointer(rhs)); *func() *uint8 {
		p := &p
		x := *p
		*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), -1))
		return x
	}() = '\x00' {
	}
	if C.strcmp(libc.CString("||"), op) == 0 {
		if (*lhs == 0 || *lhs == '0') && (*rhs == 0 || *rhs == '0') {
			C.strcpy(result, libc.CString("0"))
		} else {
			C.strcpy(result, libc.CString("1"))
		}
	} else if C.strcmp(libc.CString("&&"), op) == 0 {
		if *lhs == 0 || *lhs == '0' || *rhs == 0 || *rhs == '0' {
			C.strcpy(result, libc.CString("0"))
		} else {
			C.strcpy(result, libc.CString("1"))
		}
	} else if C.strcmp(libc.CString("=="), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) == libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs) == 0)
		}
	} else if C.strcmp(libc.CString("!="), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) != libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs))
		}
	} else if C.strcmp(libc.CString("<="), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) <= libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs) <= 0)
		}
	} else if C.strcmp(libc.CString(">="), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) >= libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs) <= 0)
		}
	} else if C.strcmp(libc.CString("<"), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) < libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs) < 0)
		}
	} else if C.strcmp(libc.CString(">"), op) == 0 {
		if is_num(lhs) != 0 && is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs)) > libc.Atoi(libc.GoString(rhs)))
		} else {
			stdio.Sprintf(result, "%d", C.strcasecmp(lhs, rhs) > 0)
		}
	} else if C.strcmp(libc.CString("/="), op) == 0 {
		stdio.Sprintf(result, "%c", func() int {
			if str_str(lhs, rhs) != nil {
				return '1'
			}
			return '0'
		}())
	} else if C.strcmp(libc.CString("*"), op) == 0 {
		stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs))*libc.Atoi(libc.GoString(rhs)))
	} else if C.strcmp(libc.CString("/"), op) == 0 {
		stdio.Sprintf(result, "%d", func() int {
			if (func() int {
				n = libc.Atoi(libc.GoString(rhs))
				return n
			}()) != 0 {
				return libc.Atoi(libc.GoString(lhs)) / n
			}
			return 0
		}())
	} else if C.strcmp(libc.CString("+"), op) == 0 {
		stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs))+libc.Atoi(libc.GoString(rhs)))
	} else if C.strcmp(libc.CString("-"), op) == 0 {
		stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(lhs))-libc.Atoi(libc.GoString(rhs)))
	} else if C.strcmp(libc.CString("!"), op) == 0 {
		if is_num(rhs) != 0 {
			stdio.Sprintf(result, "%d", libc.Atoi(libc.GoString(rhs)) == 0)
		} else {
			stdio.Sprintf(result, "%d", *rhs == 0)
		}
	}
}
func matching_quote(p *byte) *byte {
	for p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)); *p != 0 && *p != '"'; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if *p == '\\' {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
	}
	if *p == 0 {
		p = (*byte)(unsafe.Add(unsafe.Pointer(p), -1))
	}
	return p
}
func matching_paren(p *byte) *byte {
	var i int
	for func() int {
		p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		return func() int {
			i = 1
			return i
		}()
	}(); *p != 0 && i != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if *p == '(' {
			i++
		} else if *p == ')' {
			i--
		} else if *p == '"' {
			p = matching_quote(p)
		}
	}
	return func() *byte {
		p := &p
		*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), -1))
		return *p
	}()
}
func eval_expr(line *byte, result *byte, gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int) {
	var (
		expr [2048]byte
		p    *byte
	)
	_ = p
	for *line != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*line)))))&int(uint16(int16(_ISspace)))) != 0 {
		line = (*byte)(unsafe.Add(unsafe.Pointer(line), 1))
	}
	if eval_lhs_op_rhs(line, result, gohere, sc, trig, type_) != 0 {
	} else if *line == '(' {
		p = C.strcpy(&expr[0], line)
		p = matching_paren(&expr[0])
		*p = '\x00'
		eval_expr(&expr[1], result, gohere, sc, trig, type_)
	} else {
		var_subst(gohere, sc, trig, type_, line, result)
	}
}
func eval_lhs_op_rhs(expr *byte, result *byte, gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int) int {
	var (
		p      *byte
		tokens [2048]*byte
		line   [2048]byte
		lhr    [2048]byte
		rhr    [2048]byte
		i      int
		j      int
		ops    [15]*byte = [15]*byte{libc.CString("||"), libc.CString("&&"), libc.CString("=="), libc.CString("!="), libc.CString("<="), libc.CString(">="), libc.CString("<"), libc.CString(">"), libc.CString("/="), libc.CString("-"), libc.CString("+"), libc.CString("/"), libc.CString("*"), libc.CString("!"), libc.CString("\n")}
	)
	p = C.strcpy(&line[0], expr)
	for j = 0; *p != 0; j++ {
		tokens[j] = p
		if *p == '(' {
			p = (*byte)(unsafe.Add(unsafe.Pointer(matching_paren(p)), 1))
		} else if *p == '"' {
			p = (*byte)(unsafe.Add(unsafe.Pointer(matching_quote(p)), 1))
		} else if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p))))) & int(uint16(int16(_ISalnum)))) != 0 {
			for p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)); *p != 0 && ((int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISalnum)))) != 0 || (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0); p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
			}
		} else {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
	}
	tokens[j] = nil
	for i = 0; *ops[i] != '\n'; i++ {
		for j = 0; tokens[j] != nil; j++ {
			if C.strncasecmp(ops[i], tokens[j], uint64(C.strlen(ops[i]))) == 0 {
				*tokens[j] = '\x00'
				p = (*byte)(unsafe.Add(unsafe.Pointer(tokens[j]), C.strlen(ops[i])))
				eval_expr(&line[0], &lhr[0], gohere, sc, trig, type_)
				eval_expr(p, &rhr[0], gohere, sc, trig, type_)
				eval_op(ops[i], &lhr[0], &rhr[0], result, gohere, sc, trig)
				return 1
			}
		}
	}
	return 0
}
func process_if(cond *byte, gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int) int {
	var (
		result [2048]byte
		p      *byte
	)
	eval_expr(cond, &result[0], gohere, sc, trig, type_)
	p = &result[0]
	skip_spaces(&p)
	if *p == 0 || *p == '0' {
		return 0
	} else {
		return 1
	}
}
func find_end(trig *trig_data, cl *cmdlist_element) *cmdlist_element {
	var (
		c *cmdlist_element
		p *byte
	)
	if cl.Next == nil {
		script_log(libc.CString("Trigger VNum %d has 'if' without 'end'. (error 1)"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return cl
	}
	for c = cl.Next; c != nil; c = c.Next {
		for p = c.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if C.strncasecmp(libc.CString("if "), p, 3) == 0 {
			c = find_end(trig, c)
		} else if C.strncasecmp(libc.CString("end"), p, 3) == 0 {
			return c
		}
		if c.Next == nil {
			script_log(libc.CString("Trigger VNum %d has 'if' without 'end'. (error 2)"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
			return c
		}
	}
	script_log(libc.CString("Trigger VNum %d has 'if' without 'end'. (error 3)"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
	return c
}
func find_else_end(trig *trig_data, cl *cmdlist_element, gohere unsafe.Pointer, sc *script_data, type_ int) *cmdlist_element {
	var (
		c *cmdlist_element
		p *byte
	)
	if cl.Next == nil {
		return cl
	}
	for c = cl.Next; c.Next != nil; c = c.Next {
		for p = c.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if C.strncasecmp(libc.CString("if "), p, 3) == 0 {
			c = find_end(trig, c)
		} else if C.strncasecmp(libc.CString("elseif "), p, 7) == 0 {
			if process_if((*byte)(unsafe.Add(unsafe.Pointer(p), 7)), gohere, sc, trig, type_) != 0 {
				trig.Depth++
				return c
			}
		} else if C.strncasecmp(libc.CString("else"), p, 4) == 0 {
			trig.Depth++
			return c
		} else if C.strncasecmp(libc.CString("end"), p, 3) == 0 {
			return c
		}
		if c.Next == nil {
			script_log(libc.CString("Trigger VNum %d has 'if' without 'end'. (error 4)"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
			return c
		}
	}
	for p = c.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
	}
	if C.strncasecmp(libc.CString("end"), p, 3) != 0 {
		script_log(libc.CString("Trigger VNum %d has 'if' without 'end'. (error 5)"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
	}
	return c
}
func process_wait(gohere unsafe.Pointer, trig *trig_data, type_ int, cmd *byte, cl *cmdlist_element) {
	var (
		buf            [2048]byte
		arg            *byte
		wait_event_obj *wait_event_data
		when           int
		hr             int
		min            int
		ntime          int
		c              int8
	)
	arg = any_one_arg(cmd, &buf[0])
	skip_spaces(&arg)
	if *arg == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. wait w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cl.Cmd)
		return
	}
	if C.strncasecmp(arg, libc.CString("until "), 6) == 0 {
		if __isoc99_sscanf(arg, libc.CString("until %ld:%ld"), &hr, &min) == 2 {
			min += hr * 60
		} else {
			min = (hr % 100) + (hr/100)*60
		}
		ntime = (min * SECS_PER_MUD_HOUR * (int(1000000 / OPT_USEC))) / 60
		when = int((pulse % uint(SECS_PER_MUD_HOUR*(int(1000000/OPT_USEC)))) + uint(time_info.Hours*SECS_PER_MUD_HOUR*(int(1000000/OPT_USEC))))
		if when >= ntime {
			when = ((int(SECS_PER_MUD_HOUR * 24)) * (int(1000000 / OPT_USEC))) - when + ntime
		} else {
			when = ntime - when
		}
	} else {
		if __isoc99_sscanf(arg, libc.CString("%ld %c"), &when, &c) == 2 {
			if int(c) == 't' {
				when *= SECS_PER_MUD_HOUR * (int(1000000 / OPT_USEC))
			} else if int(c) == 's' {
				when *= int(1000000 / OPT_USEC)
			}
		}
	}
	wait_event_obj = new(wait_event_data)
	wait_event_obj.Trigger = trig
	wait_event_obj.Gohere = gohere
	wait_event_obj.Type = type_
	trig.Wait_event = event_create(trig_wait_event, unsafe.Pointer(wait_event_obj), when)
	trig.Curr_state = cl.Next
}
func process_set(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		arg   [2048]byte
		name  [2048]byte
		value *byte
	)
	value = two_arguments(cmd, &arg[0], &name[0])
	skip_spaces(&value)
	if name[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. set w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	add_var(&trig.Var_list, &name[0], value, func() int {
		if sc != nil {
			return sc.Context
		}
		return 0
	}())
}
func process_eval(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, cmd *byte) {
	var (
		arg    [2048]byte
		name   [2048]byte
		result [2048]byte
		expr   *byte
	)
	expr = one_argument(cmd, &arg[0])
	expr = one_argument(expr, &name[0])
	skip_spaces(&expr)
	if name[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. eval w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	eval_expr(expr, &result[0], gohere, sc, trig, type_)
	add_var(&trig.Var_list, &name[0], &result[0], func() int {
		if sc != nil {
			return sc.Context
		}
		return 0
	}())
}
func process_attach(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, cmd *byte) {
	var (
		arg       [2048]byte
		trignum_s [2048]byte
		result    [2048]byte
		id_p      *byte
		newtrig   *trig_data
		c         *char_data = nil
		o         *obj_data  = nil
		r         *room_data = nil
		trignum   int
		id        int
	)
	id_p = two_arguments(cmd, &arg[0], &trignum_s[0])
	skip_spaces(&id_p)
	if trignum_s[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. attach w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	if id_p == nil || *id_p == 0 || libc.Atoi(libc.GoString(id_p)) == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. attach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	eval_expr(id_p, &result[0], gohere, sc, trig, type_)
	if (func() int {
		id = libc.Atoi(libc.GoString(&result[0]))
		return id
	}()) == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. attach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	c = find_char(id)
	if c == nil {
		o = find_obj(id)
		if o == nil {
			r = find_room(id)
			if r == nil {
				script_log(libc.CString("Trigger: %s, VNum %d. attach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
				return
			}
		}
	}
	trignum = int(real_trigger(trig_vnum(libc.Atoi(libc.GoString(&trignum_s[0])))))
	if trignum == int(-1) || (func() *trig_data {
		newtrig = read_trigger(trignum)
		return newtrig
	}()) == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. attach invalid trigger: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &trignum_s[0])
		return
	}
	if c != nil {
		if !IS_NPC(c) {
			script_log(libc.CString("Trigger: %s, VNum %d. attach invalid target: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, GET_NAME(c))
			return
		}
		if c.Script == nil {
			c.Script = new(script_data)
		}
		add_trigger(c.Script, newtrig, -1)
		return
	}
	if o != nil {
		if o.Script == nil {
			o.Script = new(script_data)
		}
		add_trigger(o.Script, newtrig, -1)
		return
	}
	if r != nil {
		if r.Script == nil {
			r.Script = new(script_data)
		}
		add_trigger(r.Script, newtrig, -1)
		return
	}
}
func process_detach(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, cmd *byte) {
	var (
		arg       [2048]byte
		trignum_s [2048]byte
		result    [2048]byte
		id_p      *byte
		c         *char_data = nil
		o         *obj_data  = nil
		r         *room_data = nil
		id        int
	)
	id_p = two_arguments(cmd, &arg[0], &trignum_s[0])
	skip_spaces(&id_p)
	if trignum_s[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. detach w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	if id_p == nil || *id_p == 0 || libc.Atoi(libc.GoString(id_p)) == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. detach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	eval_expr(id_p, &result[0], gohere, sc, trig, type_)
	if (func() int {
		id = libc.Atoi(libc.GoString(&result[0]))
		return id
	}()) == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. detach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	c = find_char(id)
	if c == nil {
		o = find_obj(id)
		if o == nil {
			r = find_room(id)
			if r == nil {
				script_log(libc.CString("Trigger: %s, VNum %d. detach invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
				return
			}
		}
	}
	if c != nil && c.Script != nil {
		if C.strcmp(&trignum_s[0], libc.CString("all")) == 0 {
			extract_script(unsafe.Pointer(c), MOB_TRIGGER)
			return
		}
		if remove_trigger(c.Script, &trignum_s[0]) != 0 {
			if c.Script.Trig_list == nil {
				extract_script(unsafe.Pointer(c), MOB_TRIGGER)
			}
		}
		return
	}
	if o != nil && o.Script != nil {
		if C.strcmp(&trignum_s[0], libc.CString("all")) == 0 {
			extract_script(unsafe.Pointer(o), OBJ_TRIGGER)
			return
		}
		if remove_trigger(o.Script, &trignum_s[0]) != 0 {
			if o.Script.Trig_list == nil {
				extract_script(unsafe.Pointer(o), OBJ_TRIGGER)
			}
		}
		return
	}
	if r != nil && r.Script != nil {
		if C.strcmp(&trignum_s[0], libc.CString("all")) == 0 {
			extract_script(unsafe.Pointer(r), WLD_TRIGGER)
			return
		}
		if remove_trigger(r.Script, &trignum_s[0]) != 0 {
			if r.Script.Trig_list == nil {
				extract_script(unsafe.Pointer(r), WLD_TRIGGER)
			}
		}
		return
	}
}
func dg_room_of_obj(obj *obj_data) *room_data {
	if obj.In_room != room_rnum(-1) {
		return (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))
	}
	if obj.Carried_by != nil {
		return (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.Carried_by.In_room)))
	}
	if obj.Worn_by != nil {
		return (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.Worn_by.In_room)))
	}
	if obj.In_obj != nil {
		return dg_room_of_obj(obj.In_obj)
	}
	return nil
}
func makeuid_var(gohere unsafe.Pointer, sc *script_data, trig *trig_data, type_ int, cmd *byte) {
	var (
		junk    [2048]byte
		varname [2048]byte
		arg     [2048]byte
		name    [2048]byte
		uid     [2048]byte
	)
	uid[0] = '\x00'
	half_chop(cmd, &junk[0], cmd)
	half_chop(cmd, &varname[0], cmd)
	half_chop(cmd, &arg[0], cmd)
	half_chop(cmd, &name[0], cmd)
	if varname[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. makeuid w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	if arg == nil || arg[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. makeuid invalid id arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	if libc.Atoi(libc.GoString(&arg[0])) != 0 {
		var result [2048]byte
		eval_expr(&arg[0], &result[0], gohere, sc, trig, type_)
		stdio.Snprintf(&uid[0], int(2048), "%c%s", UID_CHAR, &result[0])
	} else {
		if name == nil || name[0] == 0 {
			script_log(libc.CString("Trigger: %s, VNum %d. makeuid needs name: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
			return
		}
		if is_abbrev(&arg[0], libc.CString("mob")) != 0 {
			var c *char_data = nil
			switch type_ {
			case WLD_TRIGGER:
				c = get_char_in_room((*room_data)(gohere), &name[0])
			case OBJ_TRIGGER:
				c = get_char_near_obj((*obj_data)(gohere), &name[0])
			case MOB_TRIGGER:
				c = get_char_room_vis((*char_data)(gohere), &name[0], nil)
			}
			if c != nil {
				stdio.Snprintf(&uid[0], int(2048), "%c%d", UID_CHAR, c.Id)
			}
		} else if is_abbrev(&arg[0], libc.CString("obj")) != 0 {
			var o *obj_data = nil
			switch type_ {
			case WLD_TRIGGER:
				o = get_obj_in_room((*room_data)(gohere), &name[0])
			case OBJ_TRIGGER:
				o = get_obj_near_obj((*obj_data)(gohere), &name[0])
			case MOB_TRIGGER:
				if (func() *obj_data {
					o = get_obj_in_list_vis((*char_data)(gohere), &name[0], nil, ((*char_data)(gohere)).Carrying)
					return o
				}()) == nil {
					o = get_obj_in_list_vis((*char_data)(gohere), &name[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*char_data)(gohere)).In_room)))).Contents)
				}
			}
			if o != nil {
				stdio.Snprintf(&uid[0], int(2048), "%c%d", UID_CHAR, o.Id)
			}
		} else if is_abbrev(&arg[0], libc.CString("room")) != 0 {
			var r room_rnum = room_rnum(-1)
			switch type_ {
			case WLD_TRIGGER:
				r = real_room(((*room_data)(gohere)).Number)
			case OBJ_TRIGGER:
				r = obj_room((*obj_data)(gohere))
			case MOB_TRIGGER:
				r = ((*char_data)(gohere)).In_room
			}
			if r != room_rnum(-1) {
				stdio.Snprintf(&uid[0], int(2048), "%c%d", UID_CHAR, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(r)))).Number+ROOM_ID_BASE)
			}
		} else {
			script_log(libc.CString("Trigger: %s, VNum %d. makeuid syntax error: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
			return
		}
	}
	if uid[0] != 0 {
		add_var(&trig.Var_list, &varname[0], &uid[0], func() int {
			if sc != nil {
				return sc.Context
			}
			return 0
		}())
	}
}
func process_return(trig *trig_data, cmd *byte) int {
	var (
		arg1 [2048]byte
		arg2 [2048]byte
	)
	two_arguments(cmd, &arg1[0], &arg2[0])
	if arg2[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. return w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return 1
	}
	return libc.Atoi(libc.GoString(&arg2[0]))
}
func process_unset(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		arg  [2048]byte
		var_ *byte
	)
	var_ = any_one_arg(cmd, &arg[0])
	skip_spaces(&var_)
	if *var_ == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. unset w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	if remove_var(&sc.Global_vars, var_) == 0 {
		remove_var(&trig.Var_list, var_)
	}
}
func process_remote(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		vd        *trig_var_data
		sc_remote *script_data = nil
		line      *byte
		var_      *byte
		uid_p     *byte
		arg       [2048]byte
		buf       [2048]byte
		buf2      [2048]byte
		uid       int
		context   int
		room      *room_data
		mob       *char_data
		obj       *obj_data
	)
	line = any_one_arg(cmd, &arg[0])
	two_arguments(line, &buf[0], &buf2[0])
	var_ = &buf[0]
	uid_p = &buf2[0]
	skip_spaces(&var_)
	skip_spaces(&uid_p)
	if buf[0] == 0 || buf2[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. remote: invalid arguments '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	for vd = trig.Var_list; vd != nil; vd = vd.Next {
		if C.strcasecmp(vd.Name, &buf[0]) == 0 {
			break
		}
	}
	if vd == nil {
		for vd = sc.Global_vars; vd != nil; vd = vd.Next {
			if C.strcasecmp(vd.Name, var_) == 0 && (vd.Context == 0 || vd.Context == sc.Context) {
				break
			}
		}
	}
	if vd == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. local var '%s' not found in remote call"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &buf[0])
		return
	}
	uid = libc.Atoi(libc.GoString(&buf2[0]))
	if uid <= 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. remote: illegal uid '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &buf2[0])
		return
	}
	context = vd.Context
	if (func() *room_data {
		room = find_room(uid)
		return room
	}()) != nil {
		sc_remote = room.Script
	} else if (func() *char_data {
		mob = find_char(uid)
		return mob
	}()) != nil {
		sc_remote = mob.Script
		if !IS_NPC(mob) {
			context = 0
		}
	} else if (func() *obj_data {
		obj = find_obj(uid)
		return obj
	}()) != nil {
		sc_remote = obj.Script
	} else {
		script_log(libc.CString("Trigger: %s, VNum %d. remote: uid '%ld' invalid"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, uid)
		return
	}
	if sc_remote == nil {
		return
	}
	add_var(&sc_remote.Global_vars, vd.Name, vd.Value, context)
}
func do_vdelete(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vd        *trig_var_data
		vd_prev   *trig_var_data = nil
		sc_remote *script_data   = nil
		var_      *byte
		uid_p     *byte
		buf       [2048]byte
		buf2      [2048]byte
		uid       int
		context   int
	)
	_ = context
	var room *room_data
	var mob *char_data
	var obj *obj_data
	argument = two_arguments(argument, &buf[0], &buf2[0])
	var_ = &buf[0]
	uid_p = &buf2[0]
	skip_spaces(&var_)
	skip_spaces(&uid_p)
	if buf[0] == 0 || buf2[0] == 0 {
		send_to_char(ch, libc.CString("Usage: vdelete { <variablename> | * | all } <id>\r\n"))
		return
	}
	uid = libc.Atoi(libc.GoString(&buf2[0]))
	if uid <= 0 {
		send_to_char(ch, libc.CString("vdelete: illegal id specified.\r\n"))
		return
	}
	if (func() *room_data {
		room = find_room(uid)
		return room
	}()) != nil {
		sc_remote = room.Script
	} else if (func() *char_data {
		mob = find_char(uid)
		return mob
	}()) != nil {
		sc_remote = mob.Script
		if !IS_NPC(mob) {
			context = 0
		}
	} else if (func() *obj_data {
		obj = find_obj(uid)
		return obj
	}()) != nil {
		sc_remote = obj.Script
	} else {
		send_to_char(ch, libc.CString("vdelete: cannot resolve specified id.\r\n"))
		return
	}
	if sc_remote == nil {
		send_to_char(ch, libc.CString("That id represents no global variables.(1)\r\n"))
		return
	}
	if sc_remote.Global_vars == nil {
		send_to_char(ch, libc.CString("That id represents no global variables.(2)\r\n"))
		return
	}
	if *var_ == '*' || is_abbrev(var_, libc.CString("all")) != 0 {
		var vd_next *trig_var_data
		for vd = sc_remote.Global_vars; vd != nil; vd = vd_next {
			vd_next = vd.Next
			libc.Free(unsafe.Pointer(vd.Value))
			libc.Free(unsafe.Pointer(vd.Name))
			libc.Free(unsafe.Pointer(vd))
		}
		sc_remote.Global_vars = nil
		send_to_char(ch, libc.CString("All variables deleted from that id.\r\n"))
		return
	}
	for vd = sc_remote.Global_vars; vd != nil; func() *trig_var_data {
		vd_prev = vd
		return func() *trig_var_data {
			vd = vd.Next
			return vd
		}()
	}() {
		if C.strcasecmp(vd.Name, var_) == 0 {
			break
		}
	}
	if vd == nil {
		send_to_char(ch, libc.CString("That variable cannot be located.\r\n"))
		return
	}
	if vd_prev != nil {
		vd_prev.Next = vd.Next
	} else {
		sc_remote.Global_vars = vd.Next
	}
	libc.Free(unsafe.Pointer(vd.Value))
	libc.Free(unsafe.Pointer(vd.Name))
	libc.Free(unsafe.Pointer(vd))
	send_to_char(ch, libc.CString("Deleted.\r\n"))
}
func perform_set_dg_var(ch *char_data, vict *char_data, val_arg *byte) int {
	var (
		var_name  [2048]byte
		var_value *byte
	)
	var_value = any_one_arg(val_arg, &var_name[0])
	if var_name == nil || var_name[0] == 0 || var_value == nil || *var_value == 0 {
		send_to_char(ch, libc.CString("Usage: set <char> <varname> <value>\r\n"))
		return 0
	}
	if vict.Script == nil {
		vict.Script = new(script_data)
	}
	add_var(&vict.Script.Global_vars, &var_name[0], var_value, 0)
	return 1
}
func process_rdelete(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		vd        *trig_var_data
		vd_prev   *trig_var_data = nil
		sc_remote *script_data   = nil
		line      *byte
		var_      *byte
		uid_p     *byte
		arg       [2048]byte
		buf       [64936]byte
		buf2      [64936]byte
		uid       int
		context   int
	)
	_ = context
	var room *room_data
	var mob *char_data
	var obj *obj_data
	line = any_one_arg(cmd, &arg[0])
	two_arguments(line, &buf[0], &buf2[0])
	var_ = &buf[0]
	uid_p = &buf2[0]
	skip_spaces(&var_)
	skip_spaces(&uid_p)
	if buf[0] == 0 || buf2[0] == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. rdelete: invalid arguments '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	uid = libc.Atoi(libc.GoString(&buf2[0]))
	if uid <= 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. rdelete: illegal uid '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, &buf2[0])
		return
	}
	if (func() *room_data {
		room = find_room(uid)
		return room
	}()) != nil {
		sc_remote = room.Script
	} else if (func() *char_data {
		mob = find_char(uid)
		return mob
	}()) != nil {
		sc_remote = mob.Script
		if !IS_NPC(mob) {
			context = 0
		}
	} else if (func() *obj_data {
		obj = find_obj(uid)
		return obj
	}()) != nil {
		sc_remote = obj.Script
	} else {
		script_log(libc.CString("Trigger: %s, VNum %d. remote: uid '%ld' invalid"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, uid)
		return
	}
	if sc_remote == nil {
		return
	}
	if sc_remote.Global_vars == nil {
		return
	}
	for vd = sc_remote.Global_vars; vd != nil; func() *trig_var_data {
		vd_prev = vd
		return func() *trig_var_data {
			vd = vd.Next
			return vd
		}()
	}() {
		if C.strcasecmp(vd.Name, var_) == 0 && (vd.Context == 0 || vd.Context == sc.Context) {
			break
		}
	}
	if vd == nil {
		return
	}
	if vd_prev != nil {
		vd_prev.Next = vd.Next
	} else {
		sc_remote.Global_vars = vd.Next
	}
	libc.Free(unsafe.Pointer(vd.Value))
	libc.Free(unsafe.Pointer(vd.Name))
	libc.Free(unsafe.Pointer(vd))
}
func process_global(sc *script_data, trig *trig_data, cmd *byte, id int) {
	var (
		vd   *trig_var_data
		arg  [2048]byte
		var_ *byte
	)
	var_ = any_one_arg(cmd, &arg[0])
	skip_spaces(&var_)
	if *var_ == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. global w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	for vd = trig.Var_list; vd != nil; vd = vd.Next {
		if C.strcasecmp(vd.Name, var_) == 0 {
			break
		}
	}
	if vd == nil {
		script_log(libc.CString("Trigger: %s, VNum %d. local var '%s' not found in global call"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, var_)
		return
	}
	add_var(&sc.Global_vars, vd.Name, vd.Value, id)
	remove_var(&trig.Var_list, vd.Name)
}
func process_context(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		arg  [2048]byte
		var_ *byte
	)
	var_ = any_one_arg(cmd, &arg[0])
	skip_spaces(&var_)
	if *var_ == 0 {
		script_log(libc.CString("Trigger: %s, VNum %d. context w/o an arg: '%s'"), trig.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum, cmd)
		return
	}
	sc.Context = libc.Atoi(libc.GoString(var_))
}
func extract_value(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		buf  [2048]byte
		buf2 [2048]byte
		buf3 *byte
		to   [128]byte
		num  int
	)
	buf3 = any_one_arg(cmd, &buf[0])
	half_chop(buf3, &buf2[0], &buf[0])
	C.strcpy(&to[0], &buf2[0])
	num = libc.Atoi(libc.GoString(&buf[0]))
	if num < 1 {
		script_log(libc.CString("extract number < 1!"))
		return
	}
	half_chop(&buf[0], buf3, &buf2[0])
	for num > 0 {
		half_chop(&buf2[0], &buf[0], &buf2[0])
		num--
	}
	add_var(&trig.Var_list, &to[0], &buf[0], func() int {
		if sc != nil {
			return sc.Context
		}
		return 0
	}())
}
func dg_letter_value(sc *script_data, trig *trig_data, cmd *byte) {
	var (
		junk    [2048]byte
		varname [2048]byte
		num_s   [2048]byte
		string_ [2048]byte
		num     int
	)
	half_chop(cmd, &junk[0], cmd)
	half_chop(cmd, &varname[0], cmd)
	half_chop(cmd, &num_s[0], &string_[0])
	num = libc.Atoi(libc.GoString(&num_s[0]))
	script_log(libc.CString("The use of dg_letter is deprecated"))
	script_log(libc.CString("- Use 'set <new variable> %%<text/var>.charat(index)%%' instead."))
	if num < 1 {
		script_log(libc.CString("Trigger #%d : dg_letter number < 1!"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	if num > int(C.strlen(&string_[0])) {
		script_log(libc.CString("Trigger #%d : dg_letter number > C.strlen!"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		return
	}
	junk[0] = string_[num-1]
	junk[1] = '\x00'
	add_var(&trig.Var_list, &varname[0], &junk[0], sc.Context)
}
func script_driver(go_adress unsafe.Pointer, trig *trig_data, type_ int, mode int) int {
	var (
		depth   int = 0
		ret_val int = 1
		cl      *cmdlist_element
		cmd     [2048]byte
		p       *byte
		sc      *script_data = nil
		temp    *cmdlist_element
		loops   uint           = 0
		gohere  unsafe.Pointer = nil
	)
	switch type_ {
	case MOB_TRIGGER:
		gohere = unsafe.Pointer(*(**char_data)(go_adress))
		sc = ((*char_data)(gohere)).Script
	case OBJ_TRIGGER:
		gohere = unsafe.Pointer(*(**obj_data)(go_adress))
		sc = ((*obj_data)(gohere)).Script
	case WLD_TRIGGER:
		gohere = unsafe.Pointer(*(**room_data)(go_adress))
		sc = ((*room_data)(gohere)).Script
	}
	if depth > MAX_SCRIPT_DEPTH {
		script_log(libc.CString("Trigger %d recursed beyond maximum allowed depth."), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
		switch type_ {
		case MOB_TRIGGER:
			script_log(libc.CString("It was attached to %s [%d]"), GET_NAME((*char_data)(gohere)), GET_MOB_VNUM((*char_data)(gohere)))
		case OBJ_TRIGGER:
			script_log(libc.CString("It was attached to %s [%d]"), ((*obj_data)(gohere)).Short_description, GET_OBJ_VNUM((*obj_data)(gohere)))
		case WLD_TRIGGER:
			script_log(libc.CString("It was attached to %s [%d]"), ((*room_data)(gohere)).Name, ((*room_data)(gohere)).Number)
		}
		extract_script(gohere, type_)
		return -9999999
	}
	depth++
	if mode == TRIG_NEW {
		trig.Depth = 1
		trig.Loops = 0
		sc.Context = 0
	}
	dg_owner_purged = 0
	if mode == TRIG_NEW {
		cl = trig.Cmdlist
	} else {
		cl = trig.Curr_state
	}
	for cl = cl; cl != nil && trig.Depth != 0; cl = cl.Next {
		for p = cl.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if *p == '*' {
			continue
		} else if C.strncasecmp(p, libc.CString("if "), 3) == 0 {
			if process_if((*byte)(unsafe.Add(unsafe.Pointer(p), 3)), gohere, sc, trig, type_) != 0 {
				trig.Depth++
			} else {
				cl = find_else_end(trig, cl, gohere, sc, type_)
			}
		} else if C.strncasecmp(libc.CString("elseif "), p, 7) == 0 || C.strncasecmp(libc.CString("else"), p, 4) == 0 {
			if trig.Depth == 1 {
				script_log(libc.CString("Trigger VNum %d has 'else' without 'if'."), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
				continue
			}
			cl = find_end(trig, cl)
			trig.Depth--
		} else if C.strncasecmp(libc.CString("while "), p, 6) == 0 {
			temp = find_done(cl)
			if temp == nil {
				script_log(libc.CString("Trigger VNum %d has 'while' without 'done'."), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
				return ret_val
			}
			if process_if((*byte)(unsafe.Add(unsafe.Pointer(p), 6)), gohere, sc, trig, type_) != 0 {
				temp.Original = cl
			} else {
				cl = temp
				loops = 0
			}
		} else if C.strncasecmp(libc.CString("switch "), p, 7) == 0 {
			cl = find_case(trig, cl, gohere, sc, type_, (*byte)(unsafe.Add(unsafe.Pointer(p), 7)))
		} else if C.strncasecmp(libc.CString("end"), p, 3) == 0 {
			if trig.Depth == 1 {
				script_log(libc.CString("Trigger VNum %d has 'end' without 'if'."), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
				continue
			}
			trig.Depth--
		} else if C.strncasecmp(libc.CString("done"), p, 4) == 0 {
			if cl.Original != nil {
				var orig_cmd *byte = cl.Original.Cmd
				for *orig_cmd != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*orig_cmd)))))&int(uint16(int16(_ISspace)))) != 0 {
					orig_cmd = (*byte)(unsafe.Add(unsafe.Pointer(orig_cmd), 1))
				}
				if cl.Original != nil && process_if((*byte)(unsafe.Add(unsafe.Pointer(orig_cmd), 6)), gohere, sc, trig, type_) != 0 {
					cl = cl.Original
					loops++
					trig.Loops++
					if loops == 40 {
						process_wait(gohere, trig, type_, libc.CString("wait 1"), cl)
						depth--
						return ret_val
					}
					if trig.Loops >= 5000 {
						script_log(libc.CString("Trigger VNum %d has looped 5,000 times!!!"), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trig.Nr)))).Vnum)
						break
					}
				} else {
				}
			}
		} else if C.strncasecmp(libc.CString("break"), p, 5) == 0 {
			cl = find_done(cl)
		} else if C.strncasecmp(libc.CString("case"), p, 4) == 0 {
		} else {
			var_subst(gohere, sc, trig, type_, p, &cmd[0])
			if C.strncasecmp(&cmd[0], libc.CString("eval "), 5) == 0 {
				process_eval(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("nop "), 4) == 0 {
			} else if C.strncasecmp(&cmd[0], libc.CString("extract "), 8) == 0 {
				extract_value(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("dg_letter "), 10) == 0 {
				dg_letter_value(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("makeuid "), 8) == 0 {
				makeuid_var(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("halt"), 4) == 0 {
				break
			} else if C.strncasecmp(&cmd[0], libc.CString("dg_cast "), 8) == 0 {
				do_dg_cast(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("dg_affect "), 10) == 0 {
				do_dg_affect(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("global "), 7) == 0 {
				process_global(sc, trig, &cmd[0], sc.Context)
			} else if C.strncasecmp(&cmd[0], libc.CString("context "), 8) == 0 {
				process_context(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("remote "), 7) == 0 {
				process_remote(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("rdelete "), 8) == 0 {
				process_rdelete(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("return "), 7) == 0 {
				ret_val = process_return(trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("set "), 4) == 0 {
				process_set(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("unset "), 6) == 0 {
				process_unset(sc, trig, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("wait "), 5) == 0 {
				process_wait(gohere, trig, type_, &cmd[0], cl)
				depth--
				return ret_val
			} else if C.strncasecmp(&cmd[0], libc.CString("attach "), 7) == 0 {
				process_attach(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("detach "), 7) == 0 {
				process_detach(gohere, sc, trig, type_, &cmd[0])
			} else if C.strncasecmp(&cmd[0], libc.CString("version"), 7) == 0 {
				mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("%s"), DG_SCRIPT_VERSION)
			} else {
				switch type_ {
				case MOB_TRIGGER:
					command_interpreter((*char_data)(gohere), &cmd[0])
				case OBJ_TRIGGER:
					obj_command_interpreter((*obj_data)(gohere), &cmd[0])
				case WLD_TRIGGER:
					wld_command_interpreter((*room_data)(gohere), &cmd[0])
				}
				if dg_owner_purged != 0 {
					depth--
					if type_ == OBJ_TRIGGER {
						*(**obj_data)(go_adress) = nil
					}
					return ret_val
				}
			}
		}
	}
	switch type_ {
	case MOB_TRIGGER:
		sc = ((*char_data)(gohere)).Script
	case OBJ_TRIGGER:
		sc = ((*obj_data)(gohere)).Script
	case WLD_TRIGGER:
		sc = ((*room_data)(gohere)).Script
	}
	if sc != nil {
		free_varlist(trig.Var_list)
	}
	trig.Var_list = nil
	trig.Depth = 0
	depth--
	return ret_val
}
func real_trigger(vnum trig_vnum) trig_rnum {
	var (
		bot trig_rnum
		top trig_rnum
		mid trig_rnum
	)
	bot = 0
	top = trig_rnum(top_of_trigt - 1)
	if top_of_trigt == 0 || (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(bot)))).Vnum > mob_vnum(vnum) || (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(top)))).Vnum < mob_vnum(vnum) {
		return -1
	}
	for bot <= top {
		mid = (bot + top) / 2
		if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(mid)))).Vnum == mob_vnum(vnum) {
			return mid
		}
		if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(mid)))).Vnum > mob_vnum(vnum) {
			top = mid - 1
		} else {
			bot = mid + 1
		}
	}
	return -1
}
func do_tstat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		rnum int
		str  [2048]byte
	)
	half_chop(argument, &str[0], argument)
	if str[0] != 0 {
		rnum = int(real_trigger(trig_vnum(libc.Atoi(libc.GoString(&str[0])))))
		if rnum == int(-1) {
			send_to_char(ch, libc.CString("That vnum does not exist.\r\n"))
			return
		}
		do_stat_trigger(ch, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto)
	} else {
		send_to_char(ch, libc.CString("Usage: tstat <vnum>\r\n"))
	}
}
func find_case(trig *trig_data, cl *cmdlist_element, gohere unsafe.Pointer, sc *script_data, type_ int, cond *byte) *cmdlist_element {
	var (
		result [2048]byte
		c      *cmdlist_element
		p      *byte
		buf    *byte
	)
	eval_expr(cond, &result[0], gohere, sc, trig, type_)
	if cl.Next == nil {
		return cl
	}
	for c = cl.Next; c.Next != nil; c = c.Next {
		for p = c.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if C.strncasecmp(libc.CString("while "), p, 6) == 0 || C.strncasecmp(libc.CString("switch"), p, 6) == 0 {
			c = find_done(c)
		} else if C.strncasecmp(libc.CString("case "), p, 5) == 0 {
			buf = (*byte)(libc.Malloc(MAX_STRING_LENGTH))
			eval_op(libc.CString("=="), &result[0], (*byte)(unsafe.Add(unsafe.Pointer(p), 5)), buf, gohere, sc, trig)
			if *buf != 0 && *buf != '0' {
				libc.Free(unsafe.Pointer(buf))
				return c
			}
			libc.Free(unsafe.Pointer(buf))
		} else if C.strncasecmp(libc.CString("default"), p, 7) == 0 {
			return c
		} else if C.strncasecmp(libc.CString("done"), p, 3) == 0 {
			return c
		}
	}
	return c
}
func find_done(cl *cmdlist_element) *cmdlist_element {
	var (
		c *cmdlist_element
		p *byte
	)
	if cl == nil || cl.Next == nil {
		return cl
	}
	for c = cl.Next; c != nil && c.Next != nil; c = c.Next {
		for p = c.Cmd; *p != 0 && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*p)))))&int(uint16(int16(_ISspace)))) != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		}
		if C.strncasecmp(libc.CString("while "), p, 6) == 0 || C.strncasecmp(libc.CString("switch "), p, 7) == 0 {
			c = find_done(c)
		} else if C.strncasecmp(libc.CString("done"), p, 3) == 0 {
			return c
		}
	}
	return c
}
func fgetline(file *C.FILE, p *byte) int {
	var count int = 0
	for {
		*p = byte(int8(fgetc(file)))
		if *p != '\n' && C.feof(file) == 0 {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
			count++
		}
		if *p == '\n' || C.feof(file) != 0 {
			break
		}
	}
	if *p == '\n' {
		*p = '\x00'
	}
	return count
}
func read_saved_vars(ch *char_data) {
	var (
		file        *C.FILE
		context     int
		fn          [127]byte
		input_line  [1024]byte
		temp        *byte
		p           *byte
		varname     [32]byte
		context_str [16]byte
	)
	if ch.Script != nil {
		return
	}
	ch.Script = new(script_data)
	get_filename(&fn[0], uint64(127), SCRIPT_VARS_FILE, GET_NAME(ch))
	file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fn[0]), "r")))
	if file == nil {
		basic_mud_log(libc.CString("%s had no variable file"), GET_NAME(ch))
		return
	}
	for {
		if get_line(file, &input_line[0]) > 0 {
			p = func() *byte {
				temp = C.strdup(&input_line[0])
				return temp
			}()
			temp = any_one_arg(temp, &varname[0])
			temp = any_one_arg(temp, &context_str[0])
			skip_spaces(&temp)
			context = libc.Atoi(libc.GoString(&context_str[0]))
			add_var(&ch.Script.Global_vars, &varname[0], temp, context)
			libc.Free(unsafe.Pointer(p))
		}
		if C.feof(file) != 0 {
			break
		}
	}
	C.fclose(file)
}
func save_char_vars(ch *char_data) {
	var (
		file *C.FILE
		fn   [127]byte
		vars *trig_var_data
	)
	if ch.Script == nil {
		return
	}
	if IS_NPC(ch) {
		return
	}
	get_filename(&fn[0], uint64(127), SCRIPT_VARS_FILE, GET_NAME(ch))
	unlink(&fn[0])
	if ch.Script.Global_vars == nil {
		return
	}
	vars = ch.Script.Global_vars
	file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fn[0]), "wt")))
	if file == nil {
		mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: Could not open player variable file %s for writing.:%s"), &fn[0], C.strerror(*__errno_location()))
		return
	}
	for vars != nil {
		if *vars.Name != '-' {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(file)), "%s %ld %s\n", vars.Name, vars.Context, vars.Value)
		}
		vars = vars.Next
	}
	C.fclose(file)
}
func read_saved_vars_ascii(file *C.FILE, ch *char_data, count int) {
	var (
		context     int
		input_line  [1024]byte
		temp        *byte
		p           *byte
		varname     [256]byte
		context_str [256]byte
		i           int
	)
	if ch.Script != nil {
		return
	}
	ch.Script = new(script_data)
	for i = 0; i < count; i++ {
		if get_line(file, &input_line[0]) > 0 {
			p = func() *byte {
				temp = C.strdup(&input_line[0])
				return temp
			}()
			temp = any_one_arg(temp, &varname[0])
			temp = any_one_arg(temp, &context_str[0])
			skip_spaces(&temp)
			context = libc.Atoi(libc.GoString(&context_str[0]))
			add_var(&ch.Script.Global_vars, &varname[0], temp, context)
			libc.Free(unsafe.Pointer(p))
		}
	}
}
func save_char_vars_ascii(file *C.FILE, ch *char_data) {
	var (
		vars  *trig_var_data
		count int = 0
	)
	if ch.Script == nil {
		return
	}
	if IS_NPC(ch) {
		return
	}
	if ch.Script.Global_vars == nil {
		return
	}
	for vars = ch.Script.Global_vars; vars != nil; vars = vars.Next {
		if *vars.Name != '-' {
			count++
		}
	}
	if count != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(file)), "Vars: %d\n", count)
		for vars = ch.Script.Global_vars; vars != nil; vars = vars.Next {
			if *vars.Name != '-' {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(file)), "%s %ld %s\n", vars.Name, vars.Context, vars.Value)
			}
		}
	}
}

type lookup_table_t struct {
	Uid  int
	C    unsafe.Pointer
	Next *lookup_table_t
}

var lookup_table [64]lookup_table_t

func init_lookup_table() {
	var i int
	for i = 0; i < BUCKET_COUNT; i++ {
		lookup_table[i].Uid = UID_OUT_OF_RANGE
		lookup_table[i].C = nil
		lookup_table[i].Next = nil
	}
}
func find_char_by_uid_in_lookup_table(uid int) *char_data {
	var (
		bucket int             = (uid & (int(BUCKET_COUNT - 1)))
		lt     *lookup_table_t = &lookup_table[bucket]
	)
	for ; lt != nil && lt.Uid != uid; lt = lt.Next {
	}
	if lt != nil {
		return (*char_data)(lt.C)
	}
	basic_mud_log(libc.CString("find_char_by_uid_in_lookup_table : No entity with number %ld in lookup table"), uid)
	return nil
}
func find_obj_by_uid_in_lookup_table(uid int) *obj_data {
	var (
		bucket int             = (uid & (int(BUCKET_COUNT - 1)))
		lt     *lookup_table_t = &lookup_table[bucket]
	)
	for ; lt != nil && lt.Uid != uid; lt = lt.Next {
	}
	if lt != nil {
		return (*obj_data)(lt.C)
	}
	basic_mud_log(libc.CString("find_obj_by_uid_in_lookup_table : No entity with number %ld in lookup table"), uid)
	return nil
}
func add_to_lookup_table(uid int, c unsafe.Pointer) {
	var (
		bucket int             = (uid & (int(BUCKET_COUNT - 1)))
		lt     *lookup_table_t = &lookup_table[bucket]
	)
	for ; lt.Next != nil; lt = lt.Next {
		if lt.C == c && lt.Uid == uid {
			basic_mud_log(libc.CString("Add_to_lookup failed. Already there. (uid = %ld)"), uid)
			return
		}
	}
	lt.Next = new(lookup_table_t)
	lt.Next.Uid = uid
	lt.Next.C = c
}
func remove_from_lookup_table(uid int) {
	var (
		bucket int             = (uid & (int(BUCKET_COUNT - 1)))
		lt     *lookup_table_t = &lookup_table[bucket]
		flt    *lookup_table_t = nil
	)
	if uid == 0 {
		return
	}
	for ; lt != nil; lt = lt.Next {
		if lt.Uid == uid {
			flt = lt
		}
	}
	if flt != nil {
		for lt = &lookup_table[bucket]; lt.Next != flt; lt = lt.Next {
		}
		lt.Next = flt.Next
		libc.Free(unsafe.Pointer(flt))
		return
	}
	basic_mud_log(libc.CString("remove_from_lookup. UID %ld not found."), uid)
}
func check_flags_by_name_ar(array *int, numflags int, search *byte, namelist [0]*byte) int {
	var (
		i    int
		item int = -1
	)
	for i = 0; i < numflags && item < 0; i++ {
		if C.strcmp(search, namelist[i]) == 0 {
			item = i
		}
	}
	if item < 0 {
		return FALSE
	}
	if IS_SET_AR([0]bitvector_t(array), bitvector_t(item)) {
		return item
	}
	return FALSE
}