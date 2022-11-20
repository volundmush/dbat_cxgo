package main

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unsafe"
)

func add_room(room *room_data) int {
	basic_mud_log(libc.CString("REIMPLEMENT THIS!"))
	os.Exit(-1)
	return 0
}
func delete_room(rnum int) bool {
	var (
		i        int
		j        int
		ppl      *char_data
		next_ppl *char_data
		obj      *obj_data
		next_obj *obj_data
		room     *room_data
	)
	if rnum <= 0 || rnum > top_of_world {
		return false
	}
	room = &world[rnum]
	add_to_save_list(zone_table[room.Zone].Number, SL_WLD)
	htree_del(room_htree, room.Number)
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
	for obj = world[rnum].Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		obj_from_room(obj)
		obj_to_room(obj, 0)
	}
	for ppl = world[rnum].People; ppl != nil; ppl = next_ppl {
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
			if (world[i].Dir_option[j]) == nil {
				continue
			} else if (world[i].Dir_option[j]).To_room > rnum {
				(world[i].Dir_option[j]).To_room -= int(libc.BoolToInt((world[i].Dir_option[j]).To_room != int(-1)))
			} else if (world[i].Dir_option[j]).To_room == rnum {
				if ((world[i].Dir_option[j]).Keyword == nil || *(world[i].Dir_option[j]).Keyword == 0) && ((world[i].Dir_option[j]).General_description == nil || *(world[i].Dir_option[j]).General_description == 0) {
					if (world[i].Dir_option[j]).Keyword != nil {
						libc.Free(unsafe.Pointer((world[i].Dir_option[j]).Keyword))
					}
					if (world[i].Dir_option[j]).General_description != nil {
						libc.Free(unsafe.Pointer((world[i].Dir_option[j]).General_description))
					}
					libc.Free(unsafe.Pointer(world[i].Dir_option[j]))
					world[i].Dir_option[j] = nil
				} else {
					(world[i].Dir_option[j]).To_room = -1
				}
			}
		}
		if i <= 0 {
			break
		}
	}
	for i = 0; i <= top_of_zone_table; i++ {
		for j = 0; int(zone_table[i].Cmd[j].Command) != 'S'; j++ {
			switch zone_table[i].Cmd[j].Command {
			case 'M':
				fallthrough
			case 'O':
				fallthrough
			case 'T':
				fallthrough
			case 'V':
				if zone_table[i].Cmd[j].Arg3 == rnum {
					zone_table[i].Cmd[j].Command = '*'
				} else if zone_table[i].Cmd[j].Arg3 > rnum {
					zone_table[i].Cmd[j].Arg3 -= int(libc.BoolToInt(zone_table[i].Cmd[j].Arg3 != int(-1)))
				}
			case 'D':
				fallthrough
			case 'R':
				if zone_table[i].Cmd[j].Arg1 == rnum {
					zone_table[i].Cmd[j].Command = '*'
				} else if zone_table[i].Cmd[j].Arg1 > rnum {
					zone_table[i].Cmd[j].Arg1 -= int(libc.BoolToInt(zone_table[i].Cmd[j].Arg1 != int(-1)))
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
				mudlog(BRF, ADMLVL_GOD, 1, libc.CString("SYSERR: GenOLC: delete_room: Unknown zone entry found!"))
			}
		}
	}
	{
		for i = 0; i < top_shop; i++ {
			for j = 0; (shop_index[i].In_room[j]) != int(-1); j++ {
				if (shop_index[i].In_room[j]) == world[rnum].Number {
					shop_index[i].In_room[j] = 0
				}
			}
		}
	}
	for i = rnum; i < top_of_world; i++ {
		world[i] = world[i+1]
		update_wait_events(&world[i], &world[i+1])
		for ppl = world[i].People; ppl != nil; ppl = ppl.Next_in_room {
			ppl.In_room -= int(libc.BoolToInt(ppl.In_room != int(-1)))
		}
		for obj = world[i].Contents; obj != nil; obj = obj.Next_content {
			obj.In_room -= int(libc.BoolToInt(obj.In_room != int(-1)))
		}
	}
	top_of_world--
	// todo: fix this
	// world = []room_data((*room_data)(libc.Realloc(unsafe.Pointer(&world[0]), top_of_world*int(unsafe.Sizeof(room_data{}))+1)))
	return true
}
func save_rooms(zone_num int) bool {
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
		return false
	}
	basic_mud_log(libc.CString("GenOLC: save_rooms: Saving rooms in zone #%d (%d-%d)."), zone_table[zone_num].Number, genolc_zone_bottom(zone_num), zone_table[zone_num].Top)
	stdio.Snprintf(&filename[0], int(128), "%s%d.new", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		sf = stdio.FOpen(libc.GoString(&filename[0]), "w")
		return sf
	}()) == nil {
		fmt.Println(libc.CString("SYSERR: save_rooms"))
		return false
	}
	for i = genolc_zone_bottom(zone_num); i <= zone_table[zone_num].Top; i++ {
		var rnum int
		if (func() int {
			rnum = real_room(i)
			return rnum
		}()) != int(-1) {
			var j int
			room = &world[rnum]
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
			}(), STRING_TERMINATOR, &buf[0], STRING_TERMINATOR, zone_table[room.Zone].Number, &rbuf1[0], &rbuf2[0], &rbuf3[0], &rbuf4[0], room.Sector_type)
			for j = 0; j < NUM_OF_DIRS; j++ {
				if (room.Dir_option[j]) != nil {
					var dflag int
					if (room.Dir_option[j]).General_description != nil {
						libc.StrNCpy(&buf[0], (room.Dir_option[j]).General_description, int(64936-1))
						strip_cr(&buf[0])
					} else {
						buf[0] = '\x00'
					}
					if IS_SET((room.Dir_option[j]).Exit_info, 1<<0) {
						if IS_SET((room.Dir_option[j]).Exit_info, 1<<4) && IS_SET((room.Dir_option[j]).Exit_info, 1<<3) {
							dflag = 4
						} else if IS_SET((room.Dir_option[j]).Exit_info, 1<<4) {
							dflag = 3
						} else if IS_SET((room.Dir_option[j]).Exit_info, 1<<3) {
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
					stdio.Fprintf(sf, "D%d\n%s~\n%s~\n%d %d %d %d %d %d %d %d %d %d %d\n", j, &buf[0], &buf1[0], dflag, func() int {
						if (room.Dir_option[j]).Key != int(-1) {
							return (room.Dir_option[j]).Key
						}
						return -1
					}(), func() int {
						if (room.Dir_option[j]).To_room != int(-1) {
							return world[(room.Dir_option[j]).To_room].Number
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
	stdio.Snprintf(&buf[0], int(64936), "%s%d.wld", LIB_WORLD, zone_table[zone_num].Number)
	stdio.Remove(libc.GoString(&buf[0]))
	stdio.Rename(libc.GoString(&filename[0]), libc.GoString(&buf[0]))
	if in_save_list(zone_table[zone_num].Number, SL_WLD) {
		remove_from_save_list(zone_table[zone_num].Number, SL_WLD)
		create_world_index(zone_table[zone_num].Number, libc.CString("wld"))
		basic_mud_log(libc.CString("GenOLC: save_rooms: Saving rooms '%s'"), &buf[0])
	}
	return true
}
func copy_room(to *room_data, from *room_data) bool {
	free_room_strings(to)
	*to = *from
	copy_room_strings(to, from)
	from.People = nil
	from.Contents = nil
	return true
}
func copy_room_strings(dest *room_data, source *room_data) bool {
	var i int
	if dest == nil || source == nil {
		basic_mud_log(libc.CString("SYSERR: GenOLC: copy_room_strings: NULL values passed."))
		return false
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
	return true
}
func free_room_strings(room *room_data) bool {
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
	return true
}
