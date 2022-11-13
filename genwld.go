package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func add_room(room *room_data) room_rnum {
	var (
		tch   *char_data
		tobj  *obj_data
		j     int
		found int = FALSE
		i     room_rnum
	)
	if room == nil {
		return -1
	}
	if (func() room_rnum {
		i = real_room(room.Number)
		return i
	}()) != room_rnum(-1) {
		if ((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Script != nil {
			extract_script(unsafe.Pointer((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))), WLD_TRIGGER)
		}
		tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).People
		tobj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Contents
		copy_room((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))), room)
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).People = tch
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Contents = tobj
		add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(room.Zone)))).Number, SL_WLD)
		basic_mud_log(libc.CString("GenOLC: add_room: Updated existing room #%d."), room.Number)
		return i
	}
	world = (*room_data)(libc.Realloc(unsafe.Pointer(world), int(top_of_world*room_rnum(unsafe.Sizeof(room_data{}))+2)))
	top_of_world++
	for i = top_of_world; i > 0; i-- {
		if room.Number > (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i-1)))).Number {
			*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))) = *room
			copy_room_strings((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))), room)
			found = int(i)
			break
		} else {
			*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))) = *(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i-1)))
			update_wait_events((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))), (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i-1))))
			for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).People; tch != nil; tch = tch.Next_in_room {
				tch.In_room += room_rnum(libc.BoolToInt(tch.In_room != room_rnum(-1)))
			}
			for tobj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Contents; tobj != nil; tobj = tobj.Next_content {
				tobj.In_room += room_rnum(libc.BoolToInt(tobj.In_room != room_rnum(-1)))
			}
		}
		htree_add(room_htree, int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number), int64(i))
	}
	if found == 0 {
		*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*0)) = *room
		copy_room_strings((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*0)), room)
	}
	basic_mud_log(libc.CString("GenOLC: add_room: Added room %d at index #%d."), room.Number, found)
	for i = room_rnum(room.Zone); i <= room_rnum(top_of_zone_table); i++ {
		for j = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command) != 'S'; j++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command {
			case 'M':
				fallthrough
			case 'O':
				fallthrough
			case 'T':
				fallthrough
			case 'V':
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 != vnum(-1) && (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 >= vnum(found)))
			case 'D':
				fallthrough
			case 'R':
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 != vnum(-1) && (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 >= vnum(found)))
				fallthrough
			case 'G':
				fallthrough
			case 'P':
				fallthrough
			case 'E':
				fallthrough
			case '*':
			default:
				mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("SYSERR: GenOLC: add_room: Unknown zone entry found!"))
			}
		}
	}
	r_mortal_start_room += room_rnum(libc.BoolToInt(r_mortal_start_room >= room_rnum(found)))
	r_immort_start_room += room_rnum(libc.BoolToInt(r_immort_start_room >= room_rnum(found)))
	r_frozen_start_room += room_rnum(libc.BoolToInt(r_frozen_start_room >= room_rnum(found)))
	i = top_of_world + 1
	for {
		i--
		for j = 0; j < NUM_OF_DIRS; j++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]) != nil && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room != room_rnum(-1) {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room += room_rnum(libc.BoolToInt(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room >= room_rnum(found)))
			}
		}
		if i <= 0 {
			break
		}
	}
	add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(room.Zone)))).Number, SL_WLD)
	return room_rnum(found)
}
func delete_room(rnum room_rnum) int {
	var (
		i        room_rnum
		j        int
		ppl      *char_data
		next_ppl *char_data
		obj      *obj_data
		next_obj *obj_data
		room     *room_data
	)
	if rnum <= 0 || rnum > top_of_world {
		return FALSE
	}
	room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))
	add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(room.Zone)))).Number, SL_WLD)
	htree_del(room_htree, int64(room.Number))
	basic_mud_log(libc.CString("GenOLC: delete_room: Deleting room #%d (%s)."), room.Number, room.Name)
	if r_mortal_start_room == rnum {
		basic_mud_log(libc.CString("WARNING: GenOLC: delete_room: Deleting mortal start room!"))
		r_mortal_start_room = 0
	}
	if r_immort_start_room == rnum {
		basic_mud_log(libc.CString("WARNING: GenOLC: delete_room: Deleting immortal start room!"))
		r_immort_start_room = 0
	}
	if r_frozen_start_room == rnum {
		basic_mud_log(libc.CString("WARNING: GenOLC: delete_room: Deleting frozen start room!"))
		r_frozen_start_room = 0
	}
	for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		obj_from_room(obj)
		obj_to_room(obj, 0)
	}
	for ppl = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).People; ppl != nil; ppl = next_ppl {
		next_ppl = ppl.Next_in_room
		char_from_room(ppl)
		char_to_room(ppl, 0)
	}
	free_room_strings(room)
	if room.Script != nil {
		extract_script(unsafe.Pointer(room), WLD_TRIGGER)
	}
	free_proto_script(unsafe.Pointer(room), WLD_TRIGGER)
	i = top_of_world + 1
	for {
		i--
		for j = 0; j < NUM_OF_DIRS; j++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]) == nil {
				continue
			} else if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room > rnum {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room -= room_rnum(libc.BoolToInt(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room != room_rnum(-1)))
			} else if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room == rnum {
				if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).Keyword == nil || *((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).Keyword == 0) && (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).General_description == nil || *((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).General_description == 0) {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).Keyword != nil {
						libc.Free(unsafe.Pointer(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).Keyword))
					}
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).General_description != nil {
						libc.Free(unsafe.Pointer(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).General_description))
					}
					libc.Free(unsafe.Pointer((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]))
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j] = nil
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room = -1
				}
			}
		}
		if i <= 0 {
			break
		}
	}
	for i = 0; i <= room_rnum(top_of_zone_table); i++ {
		for j = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command) != 'S'; j++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command {
			case 'M':
				fallthrough
			case 'O':
				fallthrough
			case 'T':
				fallthrough
			case 'V':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 == vnum(rnum) {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command = '*'
				} else if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 > vnum(rnum) {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 -= vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 != vnum(-1)))
				}
			case 'D':
				fallthrough
			case 'R':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 == vnum(rnum) {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command = '*'
				} else if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 > vnum(rnum) {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 -= vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 != vnum(-1)))
				}
				fallthrough
			case 'G':
				fallthrough
			case 'P':
				fallthrough
			case 'E':
				fallthrough
			case '*':
			default:
				mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("SYSERR: GenOLC: delete_room: Unknown zone entry found!"))
			}
		}
	}
	{
		for i = 0; i < room_rnum(top_shop); i++ {
			for j = 0; (*(*room_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).In_room), unsafe.Sizeof(room_vnum(0))*uintptr(j)))) != room_vnum(-1); j++ {
				if (*(*room_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).In_room), unsafe.Sizeof(room_vnum(0))*uintptr(j)))) == (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Number {
					*(*room_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).In_room), unsafe.Sizeof(room_vnum(0))*uintptr(j))) = 0
				}
			}
		}
	}
	for i = rnum; i < top_of_world; i++ {
		*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))) = *(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i+1)))
		update_wait_events((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i))), (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i+1))))
		for ppl = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).People; ppl != nil; ppl = ppl.Next_in_room {
			ppl.In_room -= room_rnum(libc.BoolToInt(ppl.In_room != room_rnum(-1)))
		}
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Contents; obj != nil; obj = obj.Next_content {
			obj.In_room -= room_rnum(libc.BoolToInt(obj.In_room != room_rnum(-1)))
		}
	}
	top_of_world--
	world = (*room_data)(libc.Realloc(unsafe.Pointer(world), int(top_of_world*room_rnum(unsafe.Sizeof(room_data{}))+1)))
	return TRUE
}
func save_rooms(zone_num zone_rnum) int {
	var (
		i        int
		room     *room_data
		sf       *stdio.File
		filename [128]byte
		buf      [64936]byte
		buf1     [64936]byte
		rbuf1    [64936]byte
		rbuf2    [64936]byte
		rbuf3    [64936]byte
		rbuf4    [64936]byte
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_rooms: Invalid zone number %d passed! (0-%d)"), zone_num, top_of_zone_table)
		return FALSE
	}
	basic_mud_log(libc.CString("GenOLC: save_rooms: Saving rooms in zone #%d (%d-%d)."), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, genolc_zone_bottom(zone_num), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Top)
	stdio.Snprintf(&filename[0], int(128), "%s%d.new", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	if (func() *stdio.File {
		sf = stdio.FOpen(libc.GoString(&filename[0]), "w")
		return sf
	}()) == nil {
		perror(libc.CString("SYSERR: save_rooms"))
		return FALSE
	}
	for i = int(genolc_zone_bottom(zone_num)); i <= int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Top); i++ {
		var rnum room_rnum
		if (func() room_rnum {
			rnum = real_room(room_vnum(i))
			return rnum
		}()) != room_rnum(-1) {
			var j int
			room = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))
			libc.StrNCpy(&buf[0], func() *byte {
				if room.Description != nil {
					return room.Description
				}
				return libc.CString("Empty room.")
			}(), int(64936-1))
			strip_cr(&buf[0])
			sprintascii(&rbuf1[0], room.Room_flags[0])
			sprintascii(&rbuf2[0], room.Room_flags[1])
			sprintascii(&rbuf3[0], room.Room_flags[2])
			sprintascii(&rbuf4[0], room.Room_flags[3])
			stdio.Fprintf(sf, "#%d\n%s%c\n%s%c\n%d %s %s %s %s %d\n", room.Number, func() *byte {
				if room.Name != nil {
					return room.Name
				}
				return libc.CString("Untitled")
			}(), STRING_TERMINATOR, &buf[0], STRING_TERMINATOR, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(room.Zone)))).Number, &rbuf1[0], &rbuf2[0], &rbuf3[0], &rbuf4[0], room.Sector_type)
			for j = 0; j < NUM_OF_DIRS; j++ {
				if (room.Dir_option[j]) != nil {
					var dflag int
					if (room.Dir_option[j]).General_description != nil {
						libc.StrNCpy(&buf[0], (room.Dir_option[j]).General_description, int(64936-1))
						strip_cr(&buf[0])
					} else {
						buf[0] = '\x00'
					}
					if ((room.Dir_option[j]).Exit_info & (1 << 0)) != 0 {
						if ((room.Dir_option[j]).Exit_info&(1<<4)) != 0 && ((room.Dir_option[j]).Exit_info&(1<<3)) != 0 {
							dflag = 4
						} else if ((room.Dir_option[j]).Exit_info & (1 << 4)) != 0 {
							dflag = 3
						} else if ((room.Dir_option[j]).Exit_info & (1 << 3)) != 0 {
							dflag = 2
						} else {
							dflag = 1
						}
					} else {
						dflag = 0
					}
					if (room.Dir_option[j]).Keyword != nil {
						libc.StrNCpy(&buf1[0], (room.Dir_option[j]).Keyword, int(64936-1))
					} else {
						buf1[0] = '\x00'
					}
					stdio.Fprintf(sf, "D%d\n%s~\n%s~\n%d %d %d %d %d %d %d %d %d %d %d\n", j, &buf[0], &buf1[0], dflag, func() obj_vnum {
						if (room.Dir_option[j]).Key != obj_vnum(-1) {
							return (room.Dir_option[j]).Key
						}
						return -1
					}(), func() room_vnum {
						if (room.Dir_option[j]).To_room != room_rnum(-1) {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((room.Dir_option[j]).To_room)))).Number
						}
						return -1
					}(), (room.Dir_option[j]).Dclock, (room.Dir_option[j]).Dchide, (room.Dir_option[j]).Dcskill, (room.Dir_option[j]).Dcmove, (room.Dir_option[j]).Failsavetype, (room.Dir_option[j]).Dcfailsave, (room.Dir_option[j]).Failroom, (room.Dir_option[j]).Totalfailroom)
				}
			}
			if room.Ex_description != nil {
				var xdesc *extra_descr_data
				for xdesc = room.Ex_description; xdesc != nil; xdesc = xdesc.Next {
					libc.StrNCpy(&buf[0], xdesc.Description, int(64936))
					strip_cr(&buf[0])
					stdio.Fprintf(sf, "E\n%s~\n%s~\n", xdesc.Keyword, &buf[0])
				}
			}
			stdio.Fprintf(sf, "S\n")
			script_save_to_disk(sf, unsafe.Pointer(room), WLD_TRIGGER)
		}
	}
	stdio.Fprintf(sf, "$~\n")
	sf.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.wld", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	stdio.Remove(libc.GoString(&buf[0]))
	stdio.Rename(libc.GoString(&filename[0]), libc.GoString(&buf[0]))
	if in_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_WLD) != 0 {
		remove_from_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_WLD)
		create_world_index(int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number), libc.CString("wld"))
		basic_mud_log(libc.CString("GenOLC: save_rooms: Saving rooms '%s'"), &buf[0])
	}
	return TRUE
}
func copy_room(to *room_data, from *room_data) int {
	free_room_strings(to)
	*to = *from
	copy_room_strings(to, from)
	from.People = nil
	from.Contents = nil
	return TRUE
}
func copy_room_strings(dest *room_data, source *room_data) int {
	var i int
	if dest == nil || source == nil {
		basic_mud_log(libc.CString("SYSERR: GenOLC: copy_room_strings: NULL values passed."))
		return FALSE
	}
	dest.Description = str_udup(source.Description)
	dest.Name = str_udup(source.Name)
	for i = 0; i < NUM_OF_DIRS; i++ {
		if (source.Dir_option[i]) == nil {
			continue
		}
		dest.Dir_option[i] = new(room_direction_data)
		*(dest.Dir_option[i]) = *(source.Dir_option[i])
		if (source.Dir_option[i]).General_description != nil {
			(dest.Dir_option[i]).General_description = libc.StrDup((source.Dir_option[i]).General_description)
		}
		if (source.Dir_option[i]).Keyword != nil {
			(dest.Dir_option[i]).Keyword = libc.StrDup((source.Dir_option[i]).Keyword)
		}
	}
	if source.Ex_description != nil {
		copy_ex_descriptions(&dest.Ex_description, source.Ex_description)
	}
	return TRUE
}
func free_room_strings(room *room_data) int {
	var i int
	if room.Name != nil {
		libc.Free(unsafe.Pointer(room.Name))
	}
	if room.Description != nil {
		libc.Free(unsafe.Pointer(room.Description))
	}
	if room.Ex_description != nil {
		free_ex_descriptions(room.Ex_description)
	}
	for i = 0; i < NUM_OF_DIRS; i++ {
		if room.Dir_option[i] != nil {
			if room.Dir_option[i].General_description != nil {
				libc.Free(unsafe.Pointer(room.Dir_option[i].General_description))
			}
			if room.Dir_option[i].Keyword != nil {
				libc.Free(unsafe.Pointer(room.Dir_option[i].Keyword))
			}
			libc.Free(unsafe.Pointer(room.Dir_option[i]))
			room.Dir_option[i] = nil
		}
	}
	return TRUE
}
