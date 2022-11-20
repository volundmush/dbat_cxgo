package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func free_var_el(var_ *trig_var_data) {
	if var_.Name != nil {
		libc.Free(unsafe.Pointer(var_.Name))
	}
	if var_.Value != nil {
		libc.Free(unsafe.Pointer(var_.Value))
	}
	libc.Free(unsafe.Pointer(var_))
}
func free_varlist(vd *trig_var_data) {
	var (
		i *trig_var_data
		j *trig_var_data
	)
	for i = vd; i != nil; {
		j = i
		i = i.Next
		free_var_el(j)
	}
}
func remove_var(var_list **trig_var_data, name *byte) bool {
	var (
		i *trig_var_data
		j *trig_var_data
	)
	for func() *trig_var_data {
		j = nil
		return func() *trig_var_data {
			i = *var_list
			return i
		}()
	}(); i != nil && libc.StrCaseCmp(name, i.Name) != 0; func() *trig_var_data {
		j = i
		return func() *trig_var_data {
			i = i.Next
			return i
		}()
	}() {
	}
	if i != nil {
		if j != nil {
			j.Next = i.Next
			free_var_el(i)
		} else {
			*var_list = i.Next
			free_var_el(i)
		}
		return true
	}
	return false
}
func free_trigger(trig *trig_data) {
	libc.Free(unsafe.Pointer(trig.Name))
	trig.Name = nil
	if trig.Arglist != nil {
		libc.Free(unsafe.Pointer(trig.Arglist))
		trig.Arglist = nil
	}
	if trig.Var_list != nil {
		free_varlist(trig.Var_list)
		trig.Var_list = nil
	}
	if trig.Wait_event != nil {
		event_cancel(trig.Wait_event)
	}
	libc.Free(unsafe.Pointer(trig))
}
func extract_trigger(trig *trig_data) {
	var temp *trig_data
	if trig.Wait_event != nil {
		event_cancel(trig.Wait_event)
		trig.Wait_event = nil
	}
	trig_index[trig.Nr].Number--
	if trig == trigger_list {
		trigger_list = trig.Next_in_world
	} else {
		temp = trigger_list
		for temp != nil && temp.Next_in_world != trig {
			temp = temp.Next_in_world
		}
		if temp != nil {
			temp.Next_in_world = trig.Next_in_world
		}
	}
	free_trigger(trig)
}
func extract_script(thing unsafe.Pointer, type_ int) {
	var (
		sc        *script_data = nil
		trig      *trig_data
		next_trig *trig_data
		mob       *char_data
		obj       *obj_data
		room      *room_data
	)
	switch type_ {
	case MOB_TRIGGER:
		mob = (*char_data)(thing)
		sc = mob.Script
		mob.Script = nil
	case OBJ_TRIGGER:
		obj = (*obj_data)(thing)
		sc = obj.Script
		obj.Script = nil
	case WLD_TRIGGER:
		room = (*room_data)(thing)
		sc = room.Script
		room.Script = nil
	}
	{
		var (
			i *char_data = character_list
			j *obj_data  = object_list
			k int
		)
		if sc != nil {
			for ; i != nil; i = i.Next {
				if sc == i.Script {
					panic("assert failed")
				}
			}
			for ; j != nil; j = j.Next {
				if sc == j.Script {
					panic("assert failed")
				}
			}
			for k = 0; k < top_of_world; k++ {
				if sc == world[k].Script {
					panic("assert failed")
				}
			}
		}
	}
	for trig = sc.Trig_list; trig != nil; trig = next_trig {
		next_trig = trig.Next
		extract_trigger(trig)
	}
	sc.Trig_list = nil
	free_varlist(sc.Global_vars)
	libc.Free(unsafe.Pointer(sc))
}
func extract_script_mem(sc *script_memory) {
	var next *script_memory
	for sc != nil {
		next = sc.Next
		if sc.Cmd != nil {
			libc.Free(unsafe.Pointer(sc.Cmd))
		}
		libc.Free(unsafe.Pointer(sc))
		sc = next
	}
}
func free_proto_script(thing unsafe.Pointer, type_ int) {
	var (
		proto  *trig_proto_list = nil
		fproto *trig_proto_list
		mob    *char_data
		obj    *obj_data
		room   *room_data
	)
	switch type_ {
	case MOB_TRIGGER:
		mob = (*char_data)(thing)
		proto = mob.Proto_script
		mob.Proto_script = nil
	case OBJ_TRIGGER:
		obj = (*obj_data)(thing)
		proto = obj.Proto_script
		obj.Proto_script = nil
	case WLD_TRIGGER:
		room = (*room_data)(thing)
		proto = room.Proto_script
		room.Proto_script = nil
	}
	{
		var (
			i *char_data = character_list
			j *obj_data  = object_list
			k int
		)
		if proto != nil {
			for ; i != nil; i = i.Next {
				if proto == i.Proto_script {
					panic("assert failed")
				}
			}
			for ; j != nil; j = j.Next {
				if proto == j.Proto_script {
					panic("assert failed")
				}
			}
			for k = 0; k < top_of_world; k++ {
				if proto == world[k].Proto_script {
					panic("assert failed")
				}
			}
		}
	}
	for proto != nil {
		fproto = proto
		proto = proto.Next
		libc.Free(unsafe.Pointer(fproto))
	}
}
func copy_proto_script(source unsafe.Pointer, dest unsafe.Pointer, type_ int) {
	var (
		tp_src *trig_proto_list = nil
		tp_dst *trig_proto_list = nil
	)
	switch type_ {
	case MOB_TRIGGER:
		tp_src = ((*char_data)(source)).Proto_script
	case OBJ_TRIGGER:
		tp_src = ((*obj_data)(source)).Proto_script
	case WLD_TRIGGER:
		tp_src = ((*room_data)(source)).Proto_script
	}
	if tp_src != nil {
		tp_dst = new(trig_proto_list)
		switch type_ {
		case MOB_TRIGGER:
			((*char_data)(dest)).Proto_script = tp_dst
		case OBJ_TRIGGER:
			((*obj_data)(dest)).Proto_script = tp_dst
		case WLD_TRIGGER:
			((*room_data)(dest)).Proto_script = tp_dst
		}
		for tp_src != nil {
			tp_dst.Vnum = tp_src.Vnum
			tp_src = tp_src.Next
			if tp_src != nil {
				tp_dst.Next = new(trig_proto_list)
			}
			tp_dst = tp_dst.Next
		}
	}
}
func delete_variables(charname *byte) {
	var filename [260]byte
	if !get_filename(&filename[0], uint64(260), SCRIPT_VARS_FILE, charname) {
		return
	}
	if stdio.Remove(libc.GoString(&filename[0])) < 0 && libc.Errno != 2 {
		basic_mud_log(libc.CString("SYSERR: deleting variable file %s: %s"), &filename[0], libc.StrError(libc.Errno))
	}
}
func update_wait_events(to *room_data, from *room_data) {
	var trig *trig_data
	if from.Script == nil {
		return
	}
	for trig = from.Script.Trig_list; trig != nil; trig = trig.Next {
		if trig.Wait_event == nil {
			continue
		}
		((*wait_event_data)(trig.Wait_event.Event_obj)).Gohere = unsafe.Pointer(to)
	}
}
