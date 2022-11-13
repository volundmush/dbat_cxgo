package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
	"unsafe"
)

func do_oasis_redit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var buf3 *byte
	_ = buf3
	var buf1 [64936]byte
	var buf2 [64936]byte
	var number int = int(-1)
	var save int = 0
	var real_num int
	var d *descriptor_data
	if ch.Admlevel > 0 {
		buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	}
	if ch.Admlevel < 1 && !ROOM_FLAGGED(ch.In_room, ROOM_CANREMODEL) {
		send_to_char(ch, libc.CString("You can not remodel this room.\r\n"))
		return
	}
	var capsule *obj_data = nil
	var next_obj *obj_data = nil
	var remove *obj_data = nil
	var remodeling int = FALSE
	for capsule = ch.Carrying; capsule != nil; capsule = next_obj {
		next_obj = capsule.Next_content
		if remove != nil {
			continue
		} else if GET_OBJ_VNUM(capsule) == 0x4A96 {
			remove = capsule
		}
	}
	if remove == nil && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("You do not have a R.A.D. Remodeling Assistance Droid.\r\n"))
		return
	} else if ch.Admlevel < 1 {
		act(libc.CString("@GYou open up the computer panel on the droid and begin to program its remodeling routine.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n@G opens up a computer panel on some kind of small spherical droid and begins to program it.@n"), TRUE, ch, nil, nil, TO_ROOM)
		extract_obj(remove)
		remodeling = TRUE
	}
	if buf1[0] == 0 || ch.Admlevel < 1 {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			number = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number)
		} else {
			number = -1
		}
	} else if !unicode.IsDigit(rune(buf1[0])) {
		if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) != 0 {
			send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
			return
		}
		save = TRUE
		if is_number(&buf2[0]) != 0 {
			number = libc.Atoi(libc.GoString(&buf2[0]))
		} else if ch.Player_specials.Olc_zone >= 0 {
			var zlok zone_rnum
			if (func() zone_rnum {
				zlok = real_zone(zone_vnum(ch.Player_specials.Olc_zone))
				return zlok
			}()) == zone_rnum(-1) {
				number = -1
			} else {
				number = int(genolc_zone_bottom(zlok))
			}
		}
		if number == int(-1) {
			send_to_char(ch, libc.CString("Save which zone?\r\n"))
			return
		}
	}
	if number == int(-1) {
		number = libc.Atoi(libc.GoString(&buf1[0]))
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_REDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That room is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_redit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if save != 0 {
		d.Olc.Zone_num = real_zone(zone_vnum(number))
	} else {
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(number))
	}
	if d.Olc.Zone_num == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for that number!\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 && remodeling == FALSE {
		send_cannot_edit(ch, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all rooms in zone %d.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves room info for zone %d."), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		save_rooms(d.Olc.Zone_num)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() int {
		real_num = int(real_room(room_vnum(number)))
		return real_num
	}()) != int(-1) {
		redit_setup_existing(d, real_num)
	} else {
		redit_setup_new(d)
	}
	redit_disp_menu(d)
	d.Connected = CON_REDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
}
func redit_setup_new(d *descriptor_data) {
	d.Olc.Room = new(room_data)
	d.Olc.Room.Name = libc.CString("An unfinished room")
	d.Olc.Room.Description = libc.CString("You are in an unfinished room.\r\n")
	d.Olc.Room.Number = -1
	d.Olc.Item_type = WLD_TRIGGER
	d.Olc.Room.Proto_script = func() *trig_proto_list {
		p := &d.Olc.Script
		d.Olc.Script = nil
		return *p
	}()
	d.Olc.Value = 0
}
func redit_setup_existing(d *descriptor_data, real_num int) {
	var (
		room    *room_data
		counter int
	)
	room = new(room_data)
	*room = *(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))
	room.Name = str_udup((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Name)
	room.Description = str_udup((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Description)
	for counter = 0; counter < NUM_OF_DIRS; counter++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter] != nil {
			room.Dir_option[counter] = new(room_direction_data)
			*room.Dir_option[counter] = *(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter]
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter].General_description != nil {
				room.Dir_option[counter].General_description = libc.StrDup((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter].General_description)
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter].Keyword != nil {
				room.Dir_option[counter].Keyword = libc.StrDup((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Dir_option[counter].Keyword)
			}
		}
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Ex_description != nil {
		var (
			tdesc *extra_descr_data
			temp  *extra_descr_data
			temp2 *extra_descr_data
		)
		temp = new(extra_descr_data)
		room.Ex_description = temp
		for tdesc = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_num)))).Ex_description; tdesc != nil; tdesc = tdesc.Next {
			temp.Keyword = libc.StrDup(tdesc.Keyword)
			temp.Description = libc.StrDup(tdesc.Description)
			if tdesc.Next != nil {
				temp2 = new(extra_descr_data)
				temp.Next = temp2
				temp = temp2
			} else {
				temp.Next = nil
			}
		}
	}
	d.Olc.Room = room
	d.Olc.Value = 0
	d.Olc.Item_type = WLD_TRIGGER
	dg_olc_script_copy(d)
	room.Proto_script = nil
	room.Script = nil
}
func redit_save_internally(d *descriptor_data) {
	var (
		j        int
		room_num int
		new_room int = FALSE
		dsc      *descriptor_data
	)
	if d.Olc.Room.Number == room_vnum(-1) {
		new_room = TRUE
	}
	d.Olc.Room.Number = d.Olc.Number
	d.Olc.Room.Zone = d.Olc.Zone_num
	if (func() int {
		room_num = int(add_room(d.Olc.Room))
		return room_num
	}()) == int(-1) {
		write_to_output(d, libc.CString("Something went wrong...\r\n"))
		basic_mud_log(libc.CString("SYSERR: redit_save_internally: Something failed! (%d)"), room_num)
		return
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_num)))).Proto_script != nil && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_num)))).Proto_script != d.Olc.Script {
		free_proto_script(unsafe.Pointer((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_num)))), WLD_TRIGGER)
	}
	(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_num)))).Proto_script = d.Olc.Script
	assign_triggers(unsafe.Pointer((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room_num)))), WLD_TRIGGER)
	if new_room == 0 {
		return
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		if dsc == d {
			continue
		}
		if dsc.Connected == CON_ZEDIT {
			for j = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command) != 'S'; j++ {
				switch (*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Command {
				case 'O':
					fallthrough
				case 'M':
					fallthrough
				case 'T':
					fallthrough
				case 'V':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg3 >= vnum(room_num)))
				case 'D':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg2 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg2 >= vnum(room_num)))
					fallthrough
				case 'R':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(j)))).Arg1 >= vnum(room_num)))
				}
			}
		} else if dsc.Connected == CON_REDIT {
			for j = 0; j < NUM_OF_DIRS; j++ {
				if dsc.Olc.Room.Dir_option[j] != nil {
					if dsc.Olc.Room.Dir_option[j].To_room >= room_rnum(room_num) {
						dsc.Olc.Room.Dir_option[j].To_room++
					}
				}
			}
		}
	}
}
func redit_save_to_disk(zone_num zone_vnum) {
	save_rooms(zone_rnum(zone_num))
	update_space()
}
func free_room(room *room_data) {
	free_room_strings(room)
	if room.Script != nil {
		extract_script(unsafe.Pointer(room), WLD_TRIGGER)
	}
	free_proto_script(unsafe.Pointer(room), WLD_TRIGGER)
	libc.Free(unsafe.Pointer(room))
}
func redit_disp_extradesc_menu(d *descriptor_data) {
	var extra_desc *extra_descr_data = d.Olc.Desc
	clear_screen(d)
	write_to_output(d, libc.CString("@g1@n) Keyword: @y%s\r\n@g2@n) Description:\r\n@y%s\r\n@g3@n) Goto next description: "), func() *byte {
		if extra_desc.Keyword != nil {
			return extra_desc.Keyword
		}
		return libc.CString("<NONE>")
	}(), func() *byte {
		if extra_desc.Description != nil {
			return extra_desc.Description
		}
		return libc.CString("<NONE>")
	}())
	write_to_output(d, libc.CString(func() string {
		if extra_desc.Next == nil {
			return "Not Set.\r\n"
		}
		return "Set.\r\n"
	}()))
	write_to_output(d, libc.CString("Enter choice (0 to quit) : "))
	d.Olc.Mode = REDIT_EXTRADESC_MENU
}
func redit_disp_exit_menu(d *descriptor_data) {
	var door_buf [40]byte
	if (d.Olc.Room.Dir_option[d.Olc.Value]) == nil {
		d.Olc.Room.Dir_option[d.Olc.Value] = new(room_direction_data)
	}
	sprintbit((d.Olc.Room.Dir_option[d.Olc.Value]).Exit_info, exit_bits[:], &door_buf[0], uint64(40))
	clear_screen(d)
	write_to_output(d, libc.CString("@g1@n) Exit to     \t\t: @c%d\r\n@g2@n) Description \t\t:-\r\n@y%s\r\n@g3@n) Door name   \t\t: @y%s\r\n@g4@n) Key         \t\t: @c%d\r\n@g5@n) Door flags  \t\t: @c%s\r\n@g6@n) Purge exit.\r\n@g7@n) DC Lock\t\t: @c%d\r\n@g8@n) DC Hide     \t\t: @c%d\r\n@g9@n) DC Skill\t\t: @c%s\r\n@gA@n) DC Move\t\t: @c%d\r\n@gB@n) Skill Fail Save Type\t\t: @c%d\r\n@gC@n) DC Skill Save\t\t: @c%d\r\n@gD@n) Minor Fail Dest. Room\t: @c%d\r\n@gE@n) Major Fail Dest. Room\t: @c%d@n\r\nEnter choice, 0 to quit : "), func() room_vnum {
		if (d.Olc.Room.Dir_option[d.Olc.Value]).To_room != room_rnum(-1) {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((d.Olc.Room.Dir_option[d.Olc.Value]).To_room)))).Number
		}
		return -1
	}(), func() *byte {
		if (d.Olc.Room.Dir_option[d.Olc.Value]).General_description != nil {
			return (d.Olc.Room.Dir_option[d.Olc.Value]).General_description
		}
		return libc.CString("<NONE>")
	}(), func() *byte {
		if (d.Olc.Room.Dir_option[d.Olc.Value]).Keyword != nil {
			return (d.Olc.Room.Dir_option[d.Olc.Value]).Keyword
		}
		return libc.CString("<NONE>")
	}(), func() obj_vnum {
		if (d.Olc.Room.Dir_option[d.Olc.Value]).Key != obj_vnum(-1) {
			return (d.Olc.Room.Dir_option[d.Olc.Value]).Key
		}
		return -1
	}(), &door_buf[0], (d.Olc.Room.Dir_option[d.Olc.Value]).Dclock, (d.Olc.Room.Dir_option[d.Olc.Value]).Dchide, func() *byte {
		if (d.Olc.Room.Dir_option[d.Olc.Value]).Dcskill != 0 {
			return spell_info[(d.Olc.Room.Dir_option[d.Olc.Value]).Dcskill].Name
		}
		return libc.CString("<NONE>")
	}(), (d.Olc.Room.Dir_option[d.Olc.Value]).Dcmove, (d.Olc.Room.Dir_option[d.Olc.Value]).Failsavetype, (d.Olc.Room.Dir_option[d.Olc.Value]).Dcfailsave, (d.Olc.Room.Dir_option[d.Olc.Value]).Failroom, (d.Olc.Room.Dir_option[d.Olc.Value]).Totalfailroom)
	d.Olc.Mode = REDIT_EXIT_MENU
}
func redit_disp_exit_flag_menu(d *descriptor_data) {
	var (
		bits    [2048]byte
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_EXIT_FLAGS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter+1, exit_bits[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 2) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbit((d.Olc.Room.Dir_option[d.Olc.Value]).Exit_info, exit_bits[:], &bits[0], uint64(2048))
	write_to_output(d, libc.CString("\r\nExit flags: @c%s@n\r\nEnter exit flags, 0 to quit : "), &bits[0])
	d.Olc.Mode = REDIT_EXIT_DOORFLAGS
}
func redit_disp_flag_menu(d *descriptor_data) {
	var (
		bits    [64936]byte
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ROOM_FLAGS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter+1, room_bits[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 2) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbitarray(d.Olc.Room.Room_flags[:], room_bits[:], RF_ARRAY_MAX, &bits[0])
	write_to_output(d, libc.CString("\r\nRoom flags: @c%s@n\r\nEnter room flags, 0 to quit : "), &bits[0])
	d.Olc.Mode = REDIT_FLAGS
}
func redit_disp_sector_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ROOM_SECTORS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter, sector_types[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter sector type : "))
	d.Olc.Mode = REDIT_SECTOR
}
func redit_disp_menu(d *descriptor_data) {
	var (
		buf1 [64936]byte
		buf2 [64936]byte
		room *room_data
	)
	clear_screen(d)
	room = d.Olc.Room
	sprintbitarray(room.Room_flags[:], room_bits[:], RF_ARRAY_MAX, &buf1[0])
	sprinttype(room.Sector_type, sector_types[:], &buf2[0], uint64(64936))
	if d.Character.Admlevel > 0 {
		write_to_output(d, libc.CString("-- Room number : [@c%d@n]  \tRoom zone: [@c%d@n]\r\n@g1@n) Name        : @y%s\r\n@g2@n) Description :\r\n@y%s@g3@n) Room flags  : @c%s\r\n@g4@n) Sector type : @c%s\r\n@g5@n) Exit north  : [@c%6d@n],  @gB@n) Exit northwest : [@c%6d@n]\r\n@g6@n) Exit east   : [@c%6d@n],  @gC@n) Exit northeast : [@c%6d@n]\r\n@g7@n) Exit south  : [@c%6d@n],  @gD@n) Exit southeast : [@c%6d@n]\r\n@g8@n) Exit west   : [@c%6d@n],  @gE@n) Exit southwest : [@c%6d@n]\r\n@g9@n) Exit up     : [@c%6d@n],  @gF@n) Exit in        : [@c%6d@n]\r\n@gA@n) Exit down   : [@c%6d@n],  @gG@n) Exit out       : [@c%6d@n]\r\n@gH@n) Extra descriptions menu\r\n@gW@n) Copy Room\r\n@gX@n) Delete Room\r\n@gS@n) Script      : @c%s\r\n@gZ@n) Wiznet      :\r\n@gQ@n) Quit\r\nEnter choice : "), d.Olc.Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, room.Name, room.Description, &buf1[0], &buf2[0], func() room_vnum {
			if room.Dir_option[NORTH] != nil && room.Dir_option[NORTH].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[NORTH].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[NORTHWEST] != nil && room.Dir_option[NORTHWEST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[NORTHWEST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[EAST] != nil && room.Dir_option[EAST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[EAST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[NORTHEAST] != nil && room.Dir_option[NORTHEAST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[NORTHEAST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[SOUTH] != nil && room.Dir_option[SOUTH].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[SOUTH].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[SOUTHEAST] != nil && room.Dir_option[SOUTHEAST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[SOUTHEAST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[WEST] != nil && room.Dir_option[WEST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[WEST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[SOUTHWEST] != nil && room.Dir_option[SOUTHWEST].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[SOUTHWEST].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[UP] != nil && room.Dir_option[UP].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[UP].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[INDIR] != nil && room.Dir_option[INDIR].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[INDIR].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[DOWN] != nil && room.Dir_option[DOWN].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[DOWN].To_room)))).Number
			}
			return -1
		}(), func() room_vnum {
			if room.Dir_option[OUTDIR] != nil && room.Dir_option[OUTDIR].To_room != room_rnum(-1) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room.Dir_option[OUTDIR].To_room)))).Number
			}
			return -1
		}(), func() string {
			if d.Olc.Script != nil {
				return "Set."
			}
			return "Not Set."
		}())
	} else {
		write_to_output(d, libc.CString("-- Room number : [@c%d@n]    Room zone: [@c%d@n]\r\n@g1@n) Location Designation (Room Name)        : @y%s\r\n@g2@n) Remodeling Routine (Room Description)   :\r\n@y%s@gH@n) Extra descriptions menu\r\n@gQ@n) Quit\r\nEnter choice : "), d.Olc.Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, room.Name, room.Description)
	}
	d.Olc.Mode = REDIT_MAIN_MENU
}
func redit_parse(d *descriptor_data, arg *byte) {
	var (
		number  int
		oldtext *byte = nil
	)
	switch d.Olc.Mode {
	case REDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			redit_save_internally(d)
			mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s edits room %d."), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				redit_save_to_disk(zone_vnum(real_zone_by_thing(d.Olc.Number)))
				if d.Character.Admlevel < 1 {
					write_to_output(d, libc.CString("@GThe remodeling droid quickly launches about with remodeling the room to your specifications. It seems to finish in no time at all...@n\r\n"))
					act(libc.CString("@GThe remodeling droid quickly launches about with remodeling the room as to @g$n's@G specifications. It seems to finish in no time at all...@n"), TRUE, d.Character, nil, nil, TO_ROOM)
				} else {
					write_to_output(d, libc.CString("Room saved to disk.\r\n"))
				}
			} else {
				write_to_output(d, libc.CString("Room saved to memory.\r\n"))
			}
			cleanup_olc(d, CLEANUP_ALL)
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Room.Proto_script = d.Olc.Script
			cleanup_olc(d, CLEANUP_ALL)
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("@GYou close the droids interface and put it back among the rest of your things.@n\r\n"))
				act(libc.CString("@g$n@G stops typing information into the droid in $s hands and closes it back up.@n\r\n"), TRUE, d.Character, nil, nil, TO_ROOM)
				var obj *obj_data
				obj = read_object(0x4A96, VIRTUAL)
				obj_to_char(obj, d.Character)
			}
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to save your changes ? : "))
		}
		return
	case REDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = REDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '1':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("Enter Location Name:\r\n"))
			} else {
				write_to_output(d, libc.CString("Enter Room Name:-\r\n] "))
			}
			d.Olc.Mode = REDIT_NAME
		case '2':
			d.Olc.Mode = REDIT_DESC
			clear_screen(d)
			send_editor_help(d)
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("Set Remodel Parameters (Enter Room Description):-\r\n] "))
			} else {
				write_to_output(d, libc.CString("Enter room description:\r\n\r\n"))
			}
			if d.Olc.Room.Description != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Room.Description)
				oldtext = libc.StrDup(d.Olc.Room.Description)
			}
			string_write(d, &d.Olc.Room.Description, MAX_ROOM_DESC, 0, unsafe.Pointer(oldtext))
			d.Olc.Value = 1
		case '3':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			redit_disp_flag_menu(d)
		case '4':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			redit_disp_sector_menu(d)
		case '5':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = NORTH
			redit_disp_exit_menu(d)
		case '6':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = EAST
			redit_disp_exit_menu(d)
		case '7':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = SOUTH
			redit_disp_exit_menu(d)
		case '8':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = WEST
			redit_disp_exit_menu(d)
		case '9':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = UP
			redit_disp_exit_menu(d)
		case 'a':
			fallthrough
		case 'A':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = DOWN
			redit_disp_exit_menu(d)
		case 'b':
			fallthrough
		case 'B':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = NORTHWEST
			redit_disp_exit_menu(d)
		case 'c':
			fallthrough
		case 'C':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = NORTHEAST
			redit_disp_exit_menu(d)
		case 'd':
			fallthrough
		case 'D':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = SOUTHEAST
			redit_disp_exit_menu(d)
		case 'e':
			fallthrough
		case 'E':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = SOUTHWEST
			redit_disp_exit_menu(d)
		case 'f':
			fallthrough
		case 'F':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = INDIR
			redit_disp_exit_menu(d)
		case 'g':
			fallthrough
		case 'G':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Value = OUTDIR
			redit_disp_exit_menu(d)
		case 'h':
			fallthrough
		case 'H':
			if d.Olc.Room.Ex_description == nil {
				d.Olc.Room.Ex_description = new(extra_descr_data)
			}
			d.Olc.Desc = d.Olc.Room.Ex_description
			redit_disp_extradesc_menu(d)
		case 'w':
			fallthrough
		case 'W':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			write_to_output(d, libc.CString("Disabled for the time being!"))
		case 'x':
			fallthrough
		case 'X':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			write_to_output(d, libc.CString("Are you sure you want to delete this room? "))
			d.Olc.Mode = REDIT_DELETE
		case 's':
			fallthrough
		case 'S':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			d.Olc.Script_mode = SCRIPT_MAIN_MENU
			dg_script_menu(d)
			return
		case 'Z':
			fallthrough
		case 'z':
			if d.Character.Admlevel < 1 {
				write_to_output(d, libc.CString("That options isn't available to non-builders.\r\n"))
				break
			}
			search_replace(arg, libc.CString("z "), libc.CString(""))
			do_wiznet(d.Character, arg, 0, 0)
		default:
			write_to_output(d, libc.CString("Invalid choice!"))
			redit_disp_menu(d)
		}
		return
	case OLC_SCRIPT_EDIT:
		if dg_script_edit_parse(d, arg) != 0 {
			return
		}
	case REDIT_NAME:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Room.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Room.Name))
		}
		*(*byte)(unsafe.Add(unsafe.Pointer(arg), MAX_ROOM_NAME)) = '\x00'
		d.Olc.Room.Name = str_udup(arg)
	case REDIT_DESC:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Reached REDIT_DESC case in parse_redit()."))
		write_to_output(d, libc.CString("Oops, in REDIT_DESC.\r\n"))
	case REDIT_FLAGS:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > NUM_ROOM_FLAGS {
			write_to_output(d, libc.CString("That is not a valid choice!\r\n"))
			redit_disp_flag_menu(d)
		} else if number == 0 {
			break
		} else {
			d.Olc.Room.Room_flags[(number-1)/32] = bitvector_t(int32(int(d.Olc.Room.Room_flags[(number-1)/32]) ^ 1<<((number-1)%32)))
			redit_disp_flag_menu(d)
		}
		return
	case REDIT_SECTOR:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number >= NUM_ROOM_SECTORS {
			write_to_output(d, libc.CString("Invalid choice!"))
			redit_disp_sector_menu(d)
			return
		}
		d.Olc.Room.Sector_type = number
	case REDIT_EXIT_MENU:
		switch *arg {
		case '0':
		case '1':
			d.Olc.Mode = REDIT_EXIT_NUMBER
			write_to_output(d, libc.CString("Exit to room number : "))
			return
		case '2':
			d.Olc.Mode = REDIT_EXIT_DESCRIPTION
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter exit description:\r\n\r\n"))
			if (d.Olc.Room.Dir_option[d.Olc.Value]).General_description != nil {
				write_to_output(d, libc.CString("%s"), (d.Olc.Room.Dir_option[d.Olc.Value]).General_description)
				oldtext = libc.StrDup((d.Olc.Room.Dir_option[d.Olc.Value]).General_description)
			}
			string_write(d, &(d.Olc.Room.Dir_option[d.Olc.Value]).General_description, MAX_EXIT_DESC, 0, unsafe.Pointer(oldtext))
			return
		case '3':
			d.Olc.Mode = REDIT_EXIT_KEYWORD
			write_to_output(d, libc.CString("Enter keywords : "))
			return
		case '4':
			d.Olc.Mode = REDIT_EXIT_KEY
			write_to_output(d, libc.CString("Enter key number : "))
			return
		case '5':
			d.Olc.Mode = REDIT_EXIT_DOORFLAGS
			redit_disp_exit_flag_menu(d)
			return
		case '6':
			if (d.Olc.Room.Dir_option[d.Olc.Value]).Keyword != nil {
				libc.Free(unsafe.Pointer((d.Olc.Room.Dir_option[d.Olc.Value]).Keyword))
			}
			if (d.Olc.Room.Dir_option[d.Olc.Value]).General_description != nil {
				libc.Free(unsafe.Pointer((d.Olc.Room.Dir_option[d.Olc.Value]).General_description))
			}
			if (d.Olc.Room.Dir_option[d.Olc.Value]) != nil {
				libc.Free(unsafe.Pointer(d.Olc.Room.Dir_option[d.Olc.Value]))
			}
			d.Olc.Room.Dir_option[d.Olc.Value] = nil
		case '7':
			d.Olc.Mode = REDIT_EXIT_DCLOCK
			write_to_output(d, libc.CString("Enter lock DC number : "))
			return
		case '8':
			d.Olc.Mode = REDIT_EXIT_DCHIDE
			write_to_output(d, libc.CString("Enter door search DC number : "))
			return
		case '9':
			d.Olc.Mode = REDIT_EXIT_DCSKILL
			write_to_output(d, libc.CString("Enter skill to be checked to pass through exit : "))
			return
		case 'a':
			fallthrough
		case 'A':
			d.Olc.Mode = REDIT_EXIT_DCMOVE
			write_to_output(d, libc.CString("Enter DC for skill required to pass through exit : "))
			return
		case 'b':
			fallthrough
		case 'B':
			d.Olc.Mode = REDIT_EXIT_SAVETYPE
			write_to_output(d, libc.CString("Enter the SAVE TYPE for failed skill checks : "))
			return
		case 'c':
			fallthrough
		case 'C':
			d.Olc.Mode = REDIT_EXIT_DCSAVE
			write_to_output(d, libc.CString("Enter the DC to beat for SAVE from failed skill checks : "))
			return
		case 'd':
			fallthrough
		case 'D':
			d.Olc.Mode = REDIT_EXIT_FAILROOM
			write_to_output(d, libc.CString("Enter the room to send play to for minor save failure : "))
			return
		case 'e':
			fallthrough
		case 'E':
			d.Olc.Mode = REDIT_EXIT_TOTALFAILROOM
			write_to_output(d, libc.CString("Enter the room to send play to for major save failure : "))
			return
		default:
			write_to_output(d, libc.CString("Try again : "))
			return
		}
	case REDIT_EXIT_NUMBER:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) != -1 {
			if (func() int {
				number = int(real_room(room_vnum(number)))
				return number
			}()) == int(-1) {
				write_to_output(d, libc.CString("That room does not exist, try again : "))
				return
			}
		}
		(d.Olc.Room.Dir_option[d.Olc.Value]).To_room = room_rnum(number)
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DESCRIPTION:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Reached REDIT_EXIT_DESC case in parse_redit"))
		write_to_output(d, libc.CString("Oops, in REDIT_EXIT_DESCRIPTION.\r\n"))
	case REDIT_EXIT_KEYWORD:
		if (d.Olc.Room.Dir_option[d.Olc.Value]).Keyword != nil {
			libc.Free(unsafe.Pointer((d.Olc.Room.Dir_option[d.Olc.Value]).Keyword))
		}
		(d.Olc.Room.Dir_option[d.Olc.Value]).Keyword = str_udup(arg)
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_KEY:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Key = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Key = obj_vnum(number)
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DOORFLAGS:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > NUM_EXIT_FLAGS {
			write_to_output(d, libc.CString("That's not a valid choice!\r\n"))
			redit_disp_exit_flag_menu(d)
		} else if number == 0 {
			redit_disp_exit_menu(d)
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Exit_info ^= bitvector_t(int32(1 << (number - 1)))
			redit_disp_exit_flag_menu(d)
		}
		return
	case REDIT_EXTRADESC_KEY:
		if genolc_checkstring(d, arg) != 0 {
			if d.Olc.Desc.Keyword != nil {
				libc.Free(unsafe.Pointer(d.Olc.Desc.Keyword))
			}
			d.Olc.Desc.Keyword = str_udup(arg)
		}
		redit_disp_extradesc_menu(d)
		return
	case REDIT_EXTRADESC_MENU:
		switch func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}() {
		case 0:
			if d.Olc.Desc.Keyword == nil || d.Olc.Desc.Description == nil {
				var temp *extra_descr_data
				if d.Olc.Desc.Keyword != nil {
					libc.Free(unsafe.Pointer(d.Olc.Desc.Keyword))
				}
				if d.Olc.Desc.Description != nil {
					libc.Free(unsafe.Pointer(d.Olc.Desc.Description))
				}
				if d.Olc.Desc == d.Olc.Room.Ex_description {
					d.Olc.Room.Ex_description = d.Olc.Desc.Next
				} else {
					temp = d.Olc.Room.Ex_description
					for temp != nil && temp.Next != d.Olc.Desc {
						temp = temp.Next
					}
					if temp != nil {
						temp.Next = d.Olc.Desc.Next
					}
				}
				libc.Free(unsafe.Pointer(d.Olc.Desc))
			}
		case 1:
			d.Olc.Mode = REDIT_EXTRADESC_KEY
			write_to_output(d, libc.CString("Enter keywords, separated by spaces : "))
			return
		case 2:
			d.Olc.Mode = REDIT_EXTRADESC_DESCRIPTION
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter extra description:\r\n\r\n"))
			if d.Olc.Desc.Description != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Desc.Description)
				oldtext = libc.StrDup(d.Olc.Desc.Description)
			}
			string_write(d, &d.Olc.Desc.Description, MAX_MESSAGE_LENGTH, 0, unsafe.Pointer(oldtext))
			return
		case 3:
			if d.Olc.Desc.Keyword == nil || d.Olc.Desc.Description == nil {
				write_to_output(d, libc.CString("You can't edit the next extra description without completing this one.\r\n"))
				redit_disp_extradesc_menu(d)
			} else {
				var new_extra *extra_descr_data
				if d.Olc.Desc.Next != nil {
					d.Olc.Desc = d.Olc.Desc.Next
				} else {
					new_extra = new(extra_descr_data)
					d.Olc.Desc.Next = new_extra
					d.Olc.Desc = new_extra
				}
				redit_disp_extradesc_menu(d)
			}
			return
		}
	case REDIT_COPY:
		if (func() int {
			number = int(real_room(room_vnum(libc.Atoi(libc.GoString(arg)))))
			return number
		}()) != int(-1) {
			redit_setup_existing(d, number)
		} else {
			write_to_output(d, libc.CString("That room does not exist.\r\n"))
		}
	case REDIT_DELETE:
		if *arg == 'y' || *arg == 'Y' {
			if delete_room(real_room(d.Olc.Room.Number)) != 0 {
				write_to_output(d, libc.CString("Room deleted.\r\n"))
			} else {
				write_to_output(d, libc.CString("Couldn't delete the room!\r\n"))
			}
			cleanup_olc(d, CLEANUP_ALL)
			return
		} else if *arg == 'n' || *arg == 'N' {
			redit_disp_menu(d)
			d.Olc.Mode = REDIT_MAIN_MENU
			return
		} else {
			write_to_output(d, libc.CString("Please answer 'Y' or 'N': "))
		}
	case REDIT_EXIT_DCLOCK:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dclock = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dclock = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DCHIDE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dchide = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dchide = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DCSKILL:
		number = find_skill_num(arg, 1<<1)
		if number < 1 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcskill = 0
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcskill = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DCMOVE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcmove = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcmove = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_SAVETYPE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Failsavetype = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Failsavetype = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_DCSAVE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcfailsave = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Dcfailsave = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_FAILROOM:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Failroom = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Failroom = number
		}
		redit_disp_exit_menu(d)
		return
	case REDIT_EXIT_TOTALFAILROOM:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Totalfailroom = -1
		} else {
			(d.Olc.Room.Dir_option[d.Olc.Value]).Totalfailroom = number
		}
		redit_disp_exit_menu(d)
		return
	default:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: Reached default case in parse_redit"))
	}
	d.Olc.Value = 1
	redit_disp_menu(d)
}
func redit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case REDIT_DESC:
		redit_disp_menu(d)
	case REDIT_EXIT_DESCRIPTION:
		redit_disp_exit_menu(d)
	case REDIT_EXTRADESC_DESCRIPTION:
		redit_disp_extradesc_menu(d)
	}
}
