package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func do_oasis_copy(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i               int
		src_vnum        int
		src_rnum        int
		dst_vnum        int
		dst_rnum        int
		buf1            [2048]byte
		buf2            [2048]byte
		d               *descriptor_data
		oasis_copy_info [6]struct {
			Con_type       int
			Binary_search  func(vnum int64) int64
			Save_func      func(d *descriptor_data)
			Setup_existing func(d *descriptor_data, rnum int)
			Command        *byte
			Text           *byte
		} = [6]struct {
			Con_type       int
			Binary_search  func(vnum int64) int64
			Save_func      func(d *descriptor_data)
			Setup_existing func(d *descriptor_data, rnum int)
			Command        *byte
			Text           *byte
		}{{Con_type: CON_REDIT, Binary_search: func(vnum int64) int64 {
			return int64(real_room(room_vnum(vnum)))
		}, Save_func: redit_save_internally, Setup_existing: redit_setup_existing, Command: libc.CString("rcopy"), Text: libc.CString("room")}, {Con_type: CON_OEDIT, Binary_search: func(vnum int64) int64 {
			return int64(real_object(obj_vnum(vnum)))
		}, Save_func: oedit_save_internally, Setup_existing: oedit_setup_existing, Command: libc.CString("ocopy"), Text: libc.CString("object")}, {Con_type: CON_MEDIT, Binary_search: func(vnum int64) int64 {
			return int64(real_mobile(mob_vnum(vnum)))
		}, Save_func: medit_save_internally, Setup_existing: medit_setup_existing, Command: libc.CString("mcopy"), Text: libc.CString("mobile")}, {Con_type: CON_SEDIT, Binary_search: func(vnum int64) int64 {
			return int64(real_shop(shop_vnum(vnum)))
		}, Save_func: sedit_save_internally, Setup_existing: sedit_setup_existing, Command: libc.CString("scopy"), Text: libc.CString("shop")}, {Con_type: CON_TRIGEDIT, Binary_search: func(vnum int64) int64 {
			return int64(real_trigger(trig_vnum(vnum)))
		}, Save_func: trigedit_save, Setup_existing: trigedit_setup_existing, Command: libc.CString("tcopy"), Text: libc.CString("trigger")}, {Con_type: -1, Binary_search: nil, Save_func: nil, Setup_existing: nil, Command: libc.CString("\n"), Text: libc.CString("\n")}}
	)
	for i = 0; *oasis_copy_info[i].Text != '\n'; i++ {
		if subcmd == oasis_copy_info[i].Con_type {
			break
		}
	}
	if *oasis_copy_info[i].Text == '\n' {
		return
	}
	if IS_NPC(ch) || ch.Desc == nil || ch.Desc.Connected != CON_PLAYING {
		return
	}
	two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 || buf2[0] == 0 || (is_number(&buf1[0]) == 0 || is_number(&buf2[0]) == 0) {
		send_to_char(ch, libc.CString("Syntax: %s <source vnum> <target vnum>\r\n"), oasis_copy_info[i].Command)
		return
	}
	src_vnum = libc.Atoi(libc.GoString(&buf1[0]))
	src_rnum = int((oasis_copy_info[i].Binary_search)(int64(src_vnum)))
	if src_rnum == int(-1) {
		send_to_char(ch, libc.CString("The source %s (#%d) does not exist.\r\n"), oasis_copy_info[i].Text, src_vnum)
		return
	}
	dst_vnum = libc.Atoi(libc.GoString(&buf2[0]))
	dst_rnum = int((oasis_copy_info[i].Binary_search)(int64(dst_vnum)))
	if dst_rnum != int(-1) {
		send_to_char(ch, libc.CString("The target %s (#%d) already exists.\r\n"), oasis_copy_info[i].Text, dst_vnum)
		return
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == subcmd {
			if d.Olc != nil && d.Olc.Number == room_vnum(dst_vnum) {
				send_to_char(ch, libc.CString("The target %s (#%d) is currently being edited by %s.\r\n"), oasis_copy_info[i].Text, dst_vnum, GET_NAME(d.Character))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_copy: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if (func() zone_rnum {
		p := &d.Olc.Zone_num
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(dst_vnum))
		return *p
	}()) == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for the given vnum (#%d)!\r\n"), dst_vnum)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 {
		send_cannot_edit(ch, zone_table[d.Olc.Zone_num].Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(dst_vnum)
	send_to_char(ch, libc.CString("Copying %s: source: #%d, dest: #%d.\r\n"), oasis_copy_info[i].Text, src_vnum, dst_vnum)
	(oasis_copy_info[i].Setup_existing)(d, src_rnum)
	(oasis_copy_info[i].Save_func)(d)
	cleanup_olc(d, CLEANUP_ALL)
	send_to_char(ch, libc.CString("Done.\r\n"))
}
func do_dig(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		sdir          [2048]byte
		sroom         [2048]byte
		new_room_name *byte
		rvnum         room_vnum = room_vnum(-1)
		rrnum         room_rnum = room_rnum(-1)
		zone          zone_rnum
		dir           int = 0
		rawvnum       int
		d             *descriptor_data = ch.Desc
	)
	new_room_name = two_arguments(argument, &sdir[0], &sroom[0])
	skip_spaces(&new_room_name)
	if sdir[0] == 0 || sroom[0] == 0 {
		send_to_char(ch, libc.CString("Format: tunnel <direction> <room> - to create an exit\r\n        tunnel <direction> -1     - to delete an exit\r\n"))
		return
	}
	rawvnum = libc.Atoi(libc.GoString(&sroom[0]))
	if rawvnum == -1 {
		rvnum = -1
	} else {
		rvnum = room_vnum(rawvnum)
	}
	rrnum = real_room(rvnum)
	if (func() int {
		dir = search_block(&sdir[0], &abbr_dirs[0], FALSE)
		return dir
	}()) < 0 {
		dir = search_block(&sdir[0], &dirs[0], FALSE)
	}
	zone = world[ch.In_room].Zone
	if dir < 0 {
		send_to_char(ch, libc.CString("Cannot create an exit to the '%s'.\r\n"), &sdir[0])
		return
	}
	if zone == zone_rnum(-1) || can_edit_zone(ch, zone) == 0 {
		send_cannot_edit(ch, zone_vnum(zone))
		return
	}
	if rvnum == 0 {
		send_to_char(ch, libc.CString("The target exists, but you can't dig to limbo!\r\n"))
		return
	}
	if rvnum == room_vnum(-1) {
		if (world[ch.In_room].Dir_option[dir]) != nil {
			if (world[ch.In_room].Dir_option[dir]).General_description != nil {
				libc.Free(unsafe.Pointer((world[ch.In_room].Dir_option[dir]).General_description))
			}
			if (world[ch.In_room].Dir_option[dir]).Keyword != nil {
				libc.Free(unsafe.Pointer((world[ch.In_room].Dir_option[dir]).Keyword))
			}
			libc.Free(unsafe.Pointer(world[ch.In_room].Dir_option[dir]))
			world[ch.In_room].Dir_option[dir] = nil
			add_to_save_list(zone_table[world[ch.In_room].Zone].Number, SL_WLD)
			send_to_char(ch, libc.CString("You remove the exit to the %s.\r\n"), dirs[dir])
			return
		}
		send_to_char(ch, libc.CString("There is no exit to the %s.\r\nNo exit removed.\r\n"), dirs[dir])
		return
	}
	if (world[ch.In_room].Dir_option[dir]) != nil {
		send_to_char(ch, libc.CString("There already is an exit to the %s.\r\n"), dirs[dir])
		return
	}
	zone = real_zone_by_thing(rvnum)
	if zone == zone_rnum(-1) {
		send_to_char(ch, libc.CString("You cannot link to a non-existing zone!\r\n"))
		return
	}
	if can_edit_zone(ch, zone) == 0 {
		send_cannot_edit(ch, zone_vnum(zone))
		return
	}
	if rrnum == room_rnum(-1) {
		if d.Olc != nil {
			mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_dig: Player already had olc structure."))
			libc.Free(unsafe.Pointer(d.Olc))
		}
		d.Olc = new(oasis_olc_data)
		d.Olc.Zone_num = zone
		d.Olc.Number = rvnum
		d.Olc.Room = new(room_data)
		if *new_room_name != 0 {
			d.Olc.Room.Name = libc.StrDup(new_room_name)
		} else {
			d.Olc.Room.Name = libc.CString("An unfinished room")
		}
		d.Olc.Room.Description = libc.CString("You are in an unfinished room.\r\n")
		d.Olc.Room.Zone = d.Olc.Zone_num
		d.Olc.Room.Number = -1
		redit_save_internally(d)
		d.Olc.Value = 0
		send_to_char(ch, libc.CString("New room (%d) created.\r\n"), rvnum)
		cleanup_olc(d, CLEANUP_ALL)
		update_space()
		rrnum = real_room(rvnum)
	}
	world[ch.In_room].Dir_option[dir] = new(room_direction_data)
	(world[ch.In_room].Dir_option[dir]).General_description = nil
	(world[ch.In_room].Dir_option[dir]).Keyword = nil
	(world[ch.In_room].Dir_option[dir]).To_room = rrnum
	add_to_save_list(zone_table[world[ch.In_room].Zone].Number, SL_WLD)
	save_rooms(zone_rnum(zone_table[world[rrnum].Zone].Number))
	send_to_char(ch, libc.CString("You make an exit %s to room %d (%s).\r\n"), dirs[dir], rvnum, world[rrnum].Name)
	if (world[rrnum].Dir_option[rev_dir[dir]]) != nil {
		send_to_char(ch, libc.CString("You cannot dig from %d to here. The target room already has an exit to the %s.\r\n"), rvnum, dirs[rev_dir[dir]])
	} else {
		world[rrnum].Dir_option[rev_dir[dir]] = new(room_direction_data)
		(world[rrnum].Dir_option[rev_dir[dir]]).General_description = nil
		(world[rrnum].Dir_option[rev_dir[dir]]).Keyword = nil
		(world[rrnum].Dir_option[rev_dir[dir]]).To_room = ch.In_room
		add_to_save_list(zone_table[world[rrnum].Zone].Number, SL_WLD)
		save_rooms(zone_rnum(zone_table[world[rrnum].Zone].Number))
	}
}
func do_rcopy(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		rvnum room_vnum
		rrnum room_rnum
		trnum room_rnum
		tvnum room_rnum
		zone  zone_rnum
		arg   [2048]byte
		arg2  [2048]byte
		i     int
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: rclone <base rnum> <target rnum>\r\nOr be in the room you wish to be the target and provide base vnum.\r\n"))
		return
	}
	if arg2[0] == 0 {
		tvnum = room_rnum(libc.Atoi(libc.GoString(&arg[0])))
		rvnum = room_vnum(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
	} else if arg2[0] != 0 {
		rvnum = room_vnum(libc.Atoi(libc.GoString(&arg2[0])))
		tvnum = room_rnum(libc.Atoi(libc.GoString(&arg[0])))
	}
	if rvnum == 0 || tvnum == 0 {
		send_to_char(ch, libc.CString("The void is fine as it is, try again.\r\n"))
		return
	}
	if rvnum == room_vnum(tvnum) {
		send_to_char(ch, libc.CString("Don't copy a room to its self!\r\n"))
		return
	}
	trnum = real_room(room_vnum(tvnum))
	rrnum = real_room(rvnum)
	if trnum == room_rnum(-1) {
		send_to_char(ch, libc.CString("Could not find base room: %d\r\n"), tvnum)
		return
	}
	if rrnum == room_rnum(-1) {
		send_to_char(ch, libc.CString("Could not find target room: %d\r\n"), rvnum)
		return
	}
	zone = world[rrnum].Zone
	if zone == zone_rnum(-1) || can_edit_zone(ch, zone) == 0 {
		send_to_char(ch, libc.CString("\r\n"))
		send_cannot_edit(ch, zone_vnum(zone))
		return
	}
	if world[rrnum].Name != nil {
		libc.Free(unsafe.Pointer(world[rrnum].Name))
	}
	if world[rrnum].Description != nil {
		libc.Free(unsafe.Pointer(world[rrnum].Description))
	}
	if world[rrnum].Ex_description != nil {
		free_ex_descriptions(world[rrnum].Ex_description)
	}
	world[rrnum].Sector_type = world[trnum].Sector_type
	world[rrnum].Description = str_udup(world[trnum].Description)
	world[rrnum].Name = str_udup(world[trnum].Name)
	if world[trnum].Ex_description != nil {
		copy_ex_descriptions(&world[rrnum].Ex_description, world[trnum].Ex_description)
	} else {
		world[rrnum].Ex_description = nil
	}
	for i = 0; i < AF_ARRAY_MAX; i++ {
		world[rrnum].Room_flags[i] = world[trnum].Room_flags[i]
	}
	send_to_imm(libc.CString("Log: %s has copied room [%d] to room [%d]."), GET_NAME(ch), tvnum, rvnum)
	add_to_save_list(zone_table[world[rrnum].Zone].Number, SL_WLD)
	save_rooms(zone_rnum(zone_table[world[rrnum].Zone].Number))
	send_to_char(ch, libc.CString("Room [%d] copied to room [%d].\r\n"), tvnum, rvnum)
}
func redit_find_new_vnum(zone zone_rnum) room_vnum {
	var (
		vnum room_vnum = genolc_zone_bottom(zone)
		rnum room_rnum = real_room(vnum)
	)
	if rnum == room_rnum(-1) {
		return -1
	}
	for {
		if vnum > zone_table[zone].Top {
			return -1
		}
		if rnum > top_of_world || world[rnum].Number > vnum {
			break
		}
		rnum++
		vnum++
	}
	return vnum
}
func buildwalk(ch *char_data, dir int) int {
	var (
		buf  [2048]byte
		vnum room_vnum
		rnum room_rnum
	)
	if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_BUILDWALK) && ch.Admlevel >= ADMLVL_IMMORT {
		if can_edit_zone(ch, world[ch.In_room].Zone) == 0 {
			send_cannot_edit(ch, zone_vnum(world[ch.In_room].Zone))
		} else if (func() room_vnum {
			vnum = redit_find_new_vnum(world[ch.In_room].Zone)
			return vnum
		}()) == room_vnum(-1) {
			send_to_char(ch, libc.CString("No free vnums are available in this zone!\r\n"))
		} else {
			var d *descriptor_data = ch.Desc
			if d.Olc != nil {
				mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: buildwalk(): Player already had olc structure."))
				libc.Free(unsafe.Pointer(d.Olc))
			}
			d.Olc = new(oasis_olc_data)
			d.Olc.Zone_num = world[ch.In_room].Zone
			d.Olc.Number = vnum
			d.Olc.Room = new(room_data)
			d.Olc.Room.Name = libc.CString("New BuildWalk Room")
			stdio.Sprintf(&buf[0], "This unfinished room was created by %s.\r\n", GET_NAME(ch))
			d.Olc.Room.Description = libc.StrDup(&buf[0])
			d.Olc.Room.Zone = d.Olc.Zone_num
			d.Olc.Room.Number = -1
			redit_save_internally(d)
			d.Olc.Value = 0
			rnum = real_room(vnum)
			world[ch.In_room].Dir_option[dir] = new(room_direction_data)
			(world[ch.In_room].Dir_option[dir]).To_room = rnum
			world[rnum].Dir_option[rev_dir[dir]] = new(room_direction_data)
			world[rnum].Dir_option[rev_dir[dir]].To_room = ch.In_room
			send_to_char(ch, libc.CString("@yRoom #%d created by BuildWalk.@n\r\n"), vnum)
			cleanup_olc(d, CLEANUP_STRUCTS)
			update_space()
			return 1
		}
	}
	return 0
}
