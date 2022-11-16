package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
	"unsafe"
)

func do_oasis_zedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int = int(-1)
		save     int = 0
		real_num int
		d        *descriptor_data
		buf3     *byte
		buf1     [64936]byte
		buf2     [64936]byte
	)
	buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		number = int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))
	} else if !unicode.IsDigit(rune(buf1[0])) {
		if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) == 0 {
			save = TRUE
			if is_number(&buf2[0]) != 0 {
				number = libc.Atoi(libc.GoString(&buf2[0]))
			} else if ch.Player_specials.Olc_zone > 0 {
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
		} else if ch.Admlevel >= ADMLVL_IMPL {
			if libc.StrCaseCmp(libc.CString("new"), &buf1[0]) != 0 || buf3 == nil || *buf3 == 0 {
				send_to_char(ch, libc.CString("Format: zedit new <zone number> <bottom-room> <upper-room>\r\n"))
			} else {
				var (
					sbot   [2048]byte
					stop   [2048]byte
					bottom room_vnum
					top    room_vnum
				)
				skip_spaces(&buf3)
				two_arguments(buf3, &sbot[0], &stop[0])
				number = libc.Atoi(libc.GoString(&buf2[0]))
				if number < 0 {
					number = -1
				}
				bottom = room_vnum(libc.Atoi(libc.GoString(&sbot[0])))
				top = room_vnum(libc.Atoi(libc.GoString(&stop[0])))
				zedit_new_zone(ch, zone_vnum(number), bottom, top)
			}
			return
		} else {
			send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
			return
		}
	}
	if number == int(-1) {
		number = libc.Atoi(libc.GoString(&buf1[0]))
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_ZEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That zone is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_zedit: Player already had olc structure."))
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
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 {
		send_cannot_edit(ch, zone_table[d.Olc.Zone_num].Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all zone information for zone %d.\r\n"), zone_table[d.Olc.Zone_num].Number)
		mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("OLC: %s saves zone information for zone %d."), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number)
		save_zone(d.Olc.Zone_num)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() int {
		real_num = int(real_room(room_vnum(number)))
		return real_num
	}()) == int(-1) {
		write_to_output(d, libc.CString("That room does not exist.\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	zedit_setup(d, real_num)
	d.Connected = CON_ZEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	SET_BIT_AR(ch.Act[:], PLR_WRITING)
	mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number, ch.Player_specials.Olc_zone)
}
func zedit_setup(d *descriptor_data, room_num int) {
	var (
		zone     *zone_data
		subcmd   int = 0
		count    int = 0
		cmd_room int = int(-1)
	)
	zone = new(zone_data)
	zone.Name = libc.StrDup(zone_table[d.Olc.Zone_num].Name)
	if zone_table[d.Olc.Zone_num].Builders != nil {
		zone.Builders = libc.StrDup(zone_table[d.Olc.Zone_num].Builders)
	}
	zone.Lifespan = zone_table[d.Olc.Zone_num].Lifespan
	zone.Bot = zone_table[d.Olc.Zone_num].Bot
	zone.Top = zone_table[d.Olc.Zone_num].Top
	zone.Reset_mode = zone_table[d.Olc.Zone_num].Reset_mode
	zone.Zone_flags[0] = zone_table[d.Olc.Zone_num].Zone_flags[0]
	zone.Zone_flags[1] = zone_table[d.Olc.Zone_num].Zone_flags[1]
	zone.Zone_flags[2] = zone_table[d.Olc.Zone_num].Zone_flags[2]
	zone.Zone_flags[3] = zone_table[d.Olc.Zone_num].Zone_flags[3]
	zone.Min_level = zone_table[d.Olc.Zone_num].Min_level
	zone.Max_level = zone_table[d.Olc.Zone_num].Max_level
	zone.Number = 0
	zone.Age = 0
	zone.Cmd = make([]reset_com, 0)
	zone.Cmd[0].Command = 'S'
	for int(zone_table[d.Olc.Zone_num].Cmd[subcmd].Command) != 'S' {
		switch zone_table[d.Olc.Zone_num].Cmd[subcmd].Command {
		case 'M':
			fallthrough
		case 'O':
			fallthrough
		case 'T':
			fallthrough
		case 'V':
			cmd_room = int(zone_table[d.Olc.Zone_num].Cmd[subcmd].Arg3)
		case 'D':
			fallthrough
		case 'R':
			cmd_room = int(zone_table[d.Olc.Zone_num].Cmd[subcmd].Arg1)
		default:
		}
		if cmd_room == room_num {
			add_cmd_to_list((**reset_com)(unsafe.Pointer(&zone.Cmd[0])), &zone_table[d.Olc.Zone_num].Cmd[subcmd], count)
			count++
		}
		subcmd++
	}
	d.Olc.Zone = zone
	zedit_disp_menu(d)
}
func zedit_disp_flag_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
		bits    [64936]byte
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ZONE_FLAGS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter+1, zone_bits[counter], func() string {
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
	sprintbitarray(d.Olc.Zone.Zone_flags[:], zone_bits[:], ZF_ARRAY_MAX, &bits[0])
	write_to_output(d, libc.CString("\r\nZone flags: @c%s@n\r\nEnter Zone flags, 0 to quit : "), &bits[0])
	d.Olc.Mode = ZEDIT_ZONE_FLAGS
}
func zedit_new_zone(ch *char_data, vzone_num zone_vnum, bottom room_vnum, top room_vnum) {
	var (
		result int
		error  *byte
		dsc    *descriptor_data
	)
	if (func() int {
		result = int(create_new_zone(vzone_num, bottom, top, &error))
		return result
	}()) == int(-1) {
		write_to_output(ch.Desc, error)
		return
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		switch dsc.Connected {
		case CON_REDIT:
			dsc.Olc.Room.Zone += zone_rnum(libc.BoolToInt(dsc.Olc.Zone_num >= zone_rnum(result)))
			fallthrough
		case CON_ZEDIT:
			fallthrough
		case CON_MEDIT:
			fallthrough
		case CON_SEDIT:
			fallthrough
		case CON_OEDIT:
			fallthrough
		case CON_TRIGEDIT:
			fallthrough
		case CON_GEDIT:
			dsc.Olc.Zone_num += zone_rnum(libc.BoolToInt(dsc.Olc.Zone_num >= zone_rnum(result)))
		default:
		}
	}
	zedit_save_to_disk(result)
	mudlog(BRF, int(MAX(ADMLVL_BUILDER, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("OLC: %s creates new zone #%d"), GET_NAME(ch), vzone_num)
	write_to_output(ch.Desc, libc.CString("Zone created successfully.\r\n"))
}
func zedit_save_internally(d *descriptor_data) {
	var (
		mobloaded int = FALSE
		objloaded int = FALSE
		subcmd    int
		room_num  room_rnum = real_room(d.Olc.Number)
	)
	if room_num == room_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: zedit_save_internally: OLC_NUM(d) room %d not found."), d.Olc.Number)
		return
	}
	remove_room_zone_commands(d.Olc.Zone_num, room_num)
	for subcmd = 0; int((d.Olc.Zone.Cmd[subcmd]).Command) != 'S'; subcmd++ {
		switch (d.Olc.Zone.Cmd[subcmd]).Command {
		case 'G':
			fallthrough
		case 'E':
			if mobloaded != 0 {
				break
			}
			write_to_output(d, libc.CString("Equip/Give command not saved since no mob was loaded first.\r\n"))
			continue
		case 'P':
			if objloaded != 0 {
				break
			}
			write_to_output(d, libc.CString("Put command not saved since another object was not loaded first.\r\n"))
			continue
		case 'M':
			mobloaded = TRUE
		case 'O':
			objloaded = TRUE
		default:
			mobloaded = func() int {
				objloaded = FALSE
				return objloaded
			}()
		}
		add_cmd_to_list((**reset_com)(unsafe.Pointer(&zone_table[d.Olc.Zone_num].Cmd[0])), &(d.Olc.Zone.Cmd[subcmd]), subcmd)
	}
	if d.Olc.Zone.Number != 0 {
		libc.Free(unsafe.Pointer(zone_table[d.Olc.Zone_num].Name))
		libc.Free(unsafe.Pointer(zone_table[d.Olc.Zone_num].Builders))
		zone_table[d.Olc.Zone_num].Name = libc.StrDup(d.Olc.Zone.Name)
		zone_table[d.Olc.Zone_num].Builders = libc.StrDup(d.Olc.Zone.Builders)
		zone_table[d.Olc.Zone_num].Bot = d.Olc.Zone.Bot
		zone_table[d.Olc.Zone_num].Top = d.Olc.Zone.Top
		zone_table[d.Olc.Zone_num].Reset_mode = d.Olc.Zone.Reset_mode
		zone_table[d.Olc.Zone_num].Lifespan = d.Olc.Zone.Lifespan
		zone_table[d.Olc.Zone_num].Zone_flags[0] = d.Olc.Zone.Zone_flags[0]
		zone_table[d.Olc.Zone_num].Zone_flags[1] = d.Olc.Zone.Zone_flags[1]
		zone_table[d.Olc.Zone_num].Zone_flags[2] = d.Olc.Zone.Zone_flags[2]
		zone_table[d.Olc.Zone_num].Zone_flags[3] = d.Olc.Zone.Zone_flags[3]
		zone_table[d.Olc.Zone_num].Min_level = d.Olc.Zone.Min_level
		zone_table[d.Olc.Zone_num].Max_level = d.Olc.Zone.Max_level
	}
	add_to_save_list(zone_table[d.Olc.Zone_num].Number, SL_ZON)
}
func zedit_save_to_disk(zone int) {
	save_zone(zone_rnum(zone))
}
func start_change_command(d *descriptor_data, pos int) int {
	if pos < 0 || pos >= count_commands(d.Olc.Zone.Cmd) {
		return 0
	}
	d.Olc.Value = pos
	return 1
}
func zedit_disp_menu(d *descriptor_data) {
	var (
		subcmd int = 0
		room   int
	)
	_ = room
	var counter int = 0
	var buf1 [64936]byte
	clear_screen(d)
	room = int(real_room(d.Olc.Number))
	sprintbitarray(d.Olc.Zone.Zone_flags[:], zone_bits[:], ZF_ARRAY_MAX, &buf1[0])
	send_to_char(d.Character, libc.CString("Room number: [@c%d@n]\t\tRoom zone: @c%d\r\n@g1@n) Builders       : @y%s\r\n@gA@n) Zone name      : @y%s\r\n@gL@n) Lifespan       : @y%d minutes\r\n@gB@n) Bottom of zone : @y%d\r\n@gT@n) Top of zone    : @y%d\r\n@gR@n) Reset Mode     : @y%s@n\r\n@gF@n) Zone Flags     : @y%s@n\r\n@gM@n) Min Level      : @y%d@n\r\n@gX@n) Max Level      : @y%d@n\r\n@gZ@n) Wiznet         :\r\n[Command list]\r\n"), d.Olc.Number, zone_table[d.Olc.Zone_num].Number, func() *byte {
		if d.Olc.Zone.Builders != nil {
			return d.Olc.Zone.Builders
		}
		return libc.CString("None.")
	}(), func() *byte {
		if d.Olc.Zone.Name != nil {
			return d.Olc.Zone.Name
		}
		return libc.CString("<NONE!>")
	}(), d.Olc.Zone.Lifespan, d.Olc.Zone.Bot, d.Olc.Zone.Top, func() string {
		if d.Olc.Zone.Reset_mode != 0 {
			if d.Olc.Zone.Reset_mode == 1 {
				return "Reset when no players are in zone."
			}
			return "Normal reset."
		}
		return "Never reset"
	}(), &buf1[0], d.Olc.Zone.Min_level, d.Olc.Zone.Max_level)
	for int((d.Olc.Zone.Cmd[subcmd]).Command) != 'S' {
		write_to_output(d, libc.CString("@n%d - @y"), func() int {
			p := &counter
			x := *p
			*p++
			return x
		}())
		switch (d.Olc.Zone.Cmd[subcmd]).Command {
		case 'M':
			write_to_output(d, libc.CString("%sLoad %s@y [@c%d@y], Max : %d, MaxR %d, Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), mob_proto[(d.Olc.Zone.Cmd[subcmd]).Arg1].Short_descr, mob_index[(d.Olc.Zone.Cmd[subcmd]).Arg1].Vnum, (d.Olc.Zone.Cmd[subcmd]).Arg2, (d.Olc.Zone.Cmd[subcmd]).Arg4, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'G':
			write_to_output(d, libc.CString("%sGive it %s@y [@c%d@y], Max : %d, Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg1].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg1].Vnum, (d.Olc.Zone.Cmd[subcmd]).Arg2, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'O':
			write_to_output(d, libc.CString("%sLoad %s@y [@c%d@y], Max : %d, MaxR %d, Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg1].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg1].Vnum, (d.Olc.Zone.Cmd[subcmd]).Arg2, (d.Olc.Zone.Cmd[subcmd]).Arg4, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'E':
			write_to_output(d, libc.CString("%sEquip with %s@y [@c%d@n], %s, Max : %d, Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg1].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg1].Vnum, equipment_types[(d.Olc.Zone.Cmd[subcmd]).Arg3], (d.Olc.Zone.Cmd[subcmd]).Arg2, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'P':
			write_to_output(d, libc.CString("%sPut %s@y [@c%d@n] in %s [@c%d@n], Max : %d, %% Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg1].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg1].Vnum, obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg3].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg3].Vnum, (d.Olc.Zone.Cmd[subcmd]).Arg2, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'R':
			write_to_output(d, libc.CString("%sRemove %s@y [@c%d@n] from room."), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), obj_proto[(d.Olc.Zone.Cmd[subcmd]).Arg2].Short_description, obj_index[(d.Olc.Zone.Cmd[subcmd]).Arg2].Vnum)
		case 'D':
			write_to_output(d, libc.CString("%sSet door %s@y as %s."), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), dirs[(d.Olc.Zone.Cmd[subcmd]).Arg2], func() string {
				if (d.Olc.Zone.Cmd[subcmd]).Arg3 != 0 {
					if (d.Olc.Zone.Cmd[subcmd]).Arg3 == 1 {
						return "closed"
					}
					return "locked"
				}
				return "open"
			}())
		case 'T':
			write_to_output(d, libc.CString("%sAttach trigger @c%s@y [@c%d@y] to %s, %% Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), trig_index[(d.Olc.Zone.Cmd[subcmd]).Arg2].Proto.Name, trig_index[(d.Olc.Zone.Cmd[subcmd]).Arg2].Vnum, func() string {
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == MOB_TRIGGER {
					return "mobile"
				}
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == OBJ_TRIGGER {
					return "object"
				}
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == WLD_TRIGGER {
					return "room"
				}
				return "????"
			}(), (d.Olc.Zone.Cmd[subcmd]).Arg5)
		case 'V':
			write_to_output(d, libc.CString("%sAssign global %s:%d to %s = %s, %% Chance %d"), func() string {
				if (d.Olc.Zone.Cmd[subcmd]).If_flag {
					return " then "
				}
				return ""
			}(), (d.Olc.Zone.Cmd[subcmd]).Sarg1, (d.Olc.Zone.Cmd[subcmd]).Arg2, func() string {
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == MOB_TRIGGER {
					return "mobile"
				}
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == OBJ_TRIGGER {
					return "object"
				}
				if (d.Olc.Zone.Cmd[subcmd]).Arg1 == WLD_TRIGGER {
					return "room"
				}
				return "????"
			}(), (d.Olc.Zone.Cmd[subcmd]).Sarg2, (d.Olc.Zone.Cmd[subcmd]).Arg5)
		default:
			write_to_output(d, libc.CString("<Unknown Command>"))
		}
		write_to_output(d, libc.CString("\r\n"))
		subcmd++
	}
	write_to_output(d, libc.CString("@n%d - <END OF LIST>\r\n@gN@n) Insert new command.\r\n@gE@n) Edit a command.\r\n@gD@n) Delete a command.\r\n@gQ@n) Quit\r\nEnter your choice : "), counter)
	d.Olc.Mode = ZEDIT_MAIN_MENU
}
func zedit_disp_comtype(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n@gM@n) Load Mobile to room             @gO@n) Load Object to room\r\n@gE@n) Equip mobile with object        @gG@n) Give an object to a mobile\r\n@gP@n) Put object in another object    @gD@n) Open/Close/Lock a Door\r\n@gR@n) Remove an object from the room\r\n@gT@n) Assign a trigger                @gV@n) Set a global variable\r\n\r\nWhat sort of command will this be? : "))
	d.Olc.Mode = ZEDIT_COMMAND_TYPE
}
func zedit_disp_arg1(d *descriptor_data) {
	write_to_output(d, libc.CString("\r\n"))
	switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
	case 'M':
		write_to_output(d, libc.CString("Input mob's vnum : "))
		d.Olc.Mode = ZEDIT_ARG1
	case 'O':
		fallthrough
	case 'E':
		fallthrough
	case 'P':
		fallthrough
	case 'G':
		write_to_output(d, libc.CString("Input object vnum : "))
		d.Olc.Mode = ZEDIT_ARG1
	case 'D':
		fallthrough
	case 'R':
		(d.Olc.Zone.Cmd[d.Olc.Value]).Arg1 = vnum(real_room(d.Olc.Number))
		zedit_disp_arg2(d)
	case 'T':
		fallthrough
	case 'V':
		write_to_output(d, libc.CString("Input trigger type (0:mob, 1:obj, 2:room) :"))
		d.Olc.Mode = ZEDIT_ARG1
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_disp_arg1(): Help!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
		return
	}
}
func zedit_disp_arg2(d *descriptor_data) {
	var i int
	switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
	case 'M':
		fallthrough
	case 'O':
		fallthrough
	case 'E':
		fallthrough
	case 'P':
		fallthrough
	case 'G':
		write_to_output(d, libc.CString("Input the maximum number that can exist on the mud : "))
	case 'D':
		for i = 0; *dirs[i] != '\n'; i++ {
			write_to_output(d, libc.CString("%d) Exit %s.\r\n"), i, dirs[i])
		}
		write_to_output(d, libc.CString("Enter exit number for door : "))
	case 'R':
		write_to_output(d, libc.CString("Input object's vnum : "))
	case 'T':
		write_to_output(d, libc.CString("Enter the trigger VNum : "))
	case 'V':
		write_to_output(d, libc.CString("Global's context (0 for none) : "))
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_disp_arg2(): Help!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
		return
	}
	d.Olc.Mode = ZEDIT_ARG2
}
func zedit_disp_arg3(d *descriptor_data) {
	var i int = 0
	write_to_output(d, libc.CString("\r\n"))
	switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
	case 'E':
		for *equipment_types[i] != '\n' {
			write_to_output(d, libc.CString("%2d) %26.26s %2d) %26.26s\r\n"), i, equipment_types[i], i+1, func() *byte {
				if *equipment_types[i+1] != '\n' {
					return equipment_types[i+1]
				}
				return libc.CString("")
			}())
			if *equipment_types[i+1] != '\n' {
				i += 2
			} else {
				break
			}
		}
		write_to_output(d, libc.CString("Location to equip : "))
	case 'P':
		write_to_output(d, libc.CString("Virtual number of the container : "))
	case 'D':
		write_to_output(d, libc.CString("0)  Door open\r\n1)  Door closed\r\n2)  Door locked\r\nEnter state of the door : "))
	case 'V':
		fallthrough
	case 'T':
		fallthrough
	case 'M':
		fallthrough
	case 'O':
		fallthrough
	case 'R':
		fallthrough
	case 'G':
		fallthrough
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_disp_arg3(): Help!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
		return
	}
	d.Olc.Mode = ZEDIT_ARG3
}
func zedit_disp_arg4(d *descriptor_data) {
	write_to_output(d, libc.CString("\r\n"))
	switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
	case 'M':
		fallthrough
	case 'O':
		write_to_output(d, libc.CString("Input the max allowed to load from this room (Pressing enter == 0(ignore)) : "))
	case 'E':
		fallthrough
	case 'P':
		fallthrough
	case 'G':
		fallthrough
	case 'T':
		fallthrough
	case 'V':
		zedit_disp_arg5(d)
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_disp_arg4(): Help!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
		return
	}
	d.Olc.Mode = ZEDIT_ARG4
}
func zedit_disp_arg5(d *descriptor_data) {
	write_to_output(d, libc.CString("\r\n"))
	switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
	case 'M':
		fallthrough
	case 'O':
		fallthrough
	case 'E':
		fallthrough
	case 'P':
		fallthrough
	case 'G':
		write_to_output(d, libc.CString("Input the percentage chance of the load NOT occurring : "))
	case 'T':
		fallthrough
	case 'V':
		fallthrough
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_disp_arg5(): Help!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
		return
	}
	d.Olc.Mode = ZEDIT_ARG5
}
func zedit_parse(d *descriptor_data, arg *byte) {
	var (
		pos    int
		i      int = 0
		number int
	)
	switch d.Olc.Mode {
	case ZEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			zedit_save_internally(d)
			if config_info.Operation.Auto_save_olc != 0 {
				write_to_output(d, libc.CString("Saving zone info to disk.\r\n"))
				zedit_save_to_disk(int(d.Olc.Zone_num))
			} else {
				write_to_output(d, libc.CString("Saving zone info in memory.\r\n"))
			}
			mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(d.Character.Player_specials.Invis_level))), TRUE, libc.CString("OLC: %s edits zone info for room %d."), GET_NAME(d.Character), d.Olc.Number)
			fallthrough
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\n"))
			write_to_output(d, libc.CString("Do you wish to save your changes? : "))
		}
	case ZEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Zone.Age != 0 || d.Olc.Zone.Number != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = ZEDIT_CONFIRM_SAVESTRING
			} else {
				write_to_output(d, libc.CString("No changes made.\r\n"))
				cleanup_olc(d, CLEANUP_ALL)
			}
		case 'n':
			fallthrough
		case 'N':
			if int(d.Olc.Zone.Cmd[0].Command) == 'S' {
				if new_command(d.Olc.Zone, 0) != 0 && start_change_command(d, 0) != 0 {
					zedit_disp_comtype(d)
					d.Olc.Zone.Age = 1
					break
				}
			}
			write_to_output(d, libc.CString("What number in the list should the new command be? : "))
			d.Olc.Mode = ZEDIT_NEW_ENTRY
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Which command do you wish to change? : "))
			d.Olc.Mode = ZEDIT_CHANGE_ENTRY
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("Which command do you wish to delete? : "))
			d.Olc.Mode = ZEDIT_DELETE_ENTRY
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter new zone name : "))
			d.Olc.Mode = ZEDIT_ZONE_NAME
		case '1':
			if d.Character.Admlevel <= ADMLVL_BUILDER {
				d.Olc.Mode = ZEDIT_MAIN_MENU
				write_to_output(d, libc.CString("Access Denied.\r\n"))
			} else {
				write_to_output(d, libc.CString("Enter new builders list : "))
				d.Olc.Mode = ZEDIT_ZONE_BUILDERS
			}
		case 'b':
			fallthrough
		case 'B':
			if d.Character.Admlevel < ADMLVL_IMPL {
				zedit_disp_menu(d)
			} else {
				write_to_output(d, libc.CString("Enter new bottom of zone : "))
				d.Olc.Mode = ZEDIT_ZONE_BOT
			}
		case 't':
			fallthrough
		case 'T':
			if d.Character.Admlevel < ADMLVL_IMPL {
				zedit_disp_menu(d)
			} else {
				write_to_output(d, libc.CString("Enter new top of zone : "))
				d.Olc.Mode = ZEDIT_ZONE_TOP
			}
		case 'l':
			fallthrough
		case 'L':
			write_to_output(d, libc.CString("Enter new zone lifespan : "))
			d.Olc.Mode = ZEDIT_ZONE_LIFE
		case 'r':
			fallthrough
		case 'R':
			write_to_output(d, libc.CString("\r\n0) Never reset\r\n1) Reset only when no players in zone\r\n2) Normal reset\r\nEnter new zone reset type : "))
			d.Olc.Mode = ZEDIT_ZONE_RESET
		case 'm':
			fallthrough
		case 'M':
			write_to_output(d, libc.CString("Enter Minimum level to enter zone : \r\n"))
			d.Olc.Mode = ZEDIT_MIN_LEVEL
		case 'f':
			fallthrough
		case 'F':
			zedit_disp_flag_menu(d)
		case 'Z':
			fallthrough
		case 'z':
			search_replace(arg, libc.CString("z "), libc.CString(""))
			do_wiznet(d.Character, arg, 0, 0)
		case 'x':
			fallthrough
		case 'X':
			write_to_output(d, libc.CString("Enter Maximum level to enter zone : \r\n"))
			d.Olc.Mode = ZEDIT_MAX_LEVEL
		default:
			zedit_disp_menu(d)
		}
	case ZEDIT_NEW_ENTRY:
		pos = libc.Atoi(libc.GoString(arg))
		if unicode.IsDigit(rune(*arg)) && new_command(d.Olc.Zone, pos) != 0 {
			if start_change_command(d, pos) != 0 {
				zedit_disp_comtype(d)
				d.Olc.Zone.Age = 1
			}
		} else {
			zedit_disp_menu(d)
		}
	case ZEDIT_DELETE_ENTRY:
		pos = libc.Atoi(libc.GoString(arg))
		if unicode.IsDigit(rune(*arg)) {
			delete_zone_command(d.Olc.Zone, pos)
			d.Olc.Zone.Age = 1
		}
		zedit_disp_menu(d)
	case ZEDIT_CHANGE_ENTRY:
		if unicode.ToUpper(rune(*arg)) == 'A' {
			if int((d.Olc.Zone.Cmd[d.Olc.Value]).Command) == 'N' {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Command = '*'
			}
			zedit_disp_menu(d)
			break
		}
		pos = libc.Atoi(libc.GoString(arg))
		if unicode.IsDigit(rune(*arg)) && start_change_command(d, pos) != 0 {
			zedit_disp_comtype(d)
			d.Olc.Zone.Age = 1
		} else {
			zedit_disp_menu(d)
		}
	case ZEDIT_COMMAND_TYPE:
		(d.Olc.Zone.Cmd[d.Olc.Value]).Command = int8(unicode.ToUpper(rune(*arg)))
		if int((d.Olc.Zone.Cmd[d.Olc.Value]).Command) == 0 || libc.StrChr(libc.CString("MOPEDGRTV"), byte((d.Olc.Zone.Cmd[d.Olc.Value]).Command)) == nil {
			write_to_output(d, libc.CString("Invalid choice, try again : "))
		} else {
			if d.Olc.Value != 0 {
				if int((d.Olc.Zone.Cmd[d.Olc.Value]).Command) == 'T' || int((d.Olc.Zone.Cmd[d.Olc.Value]).Command) == 'V' {
					(d.Olc.Zone.Cmd[d.Olc.Value]).If_flag = true
					zedit_disp_arg1(d)
				} else {
					write_to_output(d, libc.CString("Is this command dependent on the success of the previous one? (y/n)\r\n"))
					d.Olc.Mode = ZEDIT_IF_FLAG
				}
			} else {
				(d.Olc.Zone.Cmd[d.Olc.Value]).If_flag = false
				zedit_disp_arg1(d)
			}
		}
	case ZEDIT_IF_FLAG:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			(d.Olc.Zone.Cmd[d.Olc.Value]).If_flag = true
		case 'n':
			fallthrough
		case 'N':
			(d.Olc.Zone.Cmd[d.Olc.Value]).If_flag = false
		default:
			write_to_output(d, libc.CString("Try again : "))
			return
		}
		zedit_disp_arg1(d)
	case ZEDIT_ARG1:
		if !unicode.IsDigit(rune(*arg)) {
			write_to_output(d, libc.CString("Must be a numeric value, try again : "))
			return
		}
		switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
		case 'M':
			if (func() int {
				pos = int(real_mobile(mob_vnum(libc.Atoi(libc.GoString(arg)))))
				return pos
			}()) != int(-1) {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg1 = vnum(pos)
				zedit_disp_arg2(d)
			} else {
				write_to_output(d, libc.CString("That mobile does not exist, try again : "))
			}
		case 'O':
			fallthrough
		case 'P':
			fallthrough
		case 'E':
			fallthrough
		case 'G':
			if (func() int {
				pos = int(real_object(obj_vnum(libc.Atoi(libc.GoString(arg)))))
				return pos
			}()) != int(-1) {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg1 = vnum(pos)
				zedit_disp_arg2(d)
			} else {
				write_to_output(d, libc.CString("That object does not exist, try again : "))
			}
		case 'T':
			fallthrough
		case 'V':
			if libc.Atoi(libc.GoString(arg)) < MOB_TRIGGER || libc.Atoi(libc.GoString(arg)) > WLD_TRIGGER {
				write_to_output(d, libc.CString("Invalid input."))
			} else {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg1 = vnum(libc.Atoi(libc.GoString(arg)))
				zedit_disp_arg2(d)
			}
		case 'D':
			fallthrough
		case 'R':
			fallthrough
		default:
			cleanup_olc(d, CLEANUP_ALL)
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): case ARG1: Ack!"))
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
	case ZEDIT_ARG2:
		if !unicode.IsDigit(rune(*arg)) {
			write_to_output(d, libc.CString("Must be a numeric value, try again : "))
			return
		}
		switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
		case 'M':
			fallthrough
		case 'O':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(MIN(MAX_DUPLICATES, int64(libc.Atoi(libc.GoString(arg)))))
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(real_room(d.Olc.Number))
			zedit_disp_arg4(d)
		case 'G':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(MIN(MAX_DUPLICATES, int64(libc.Atoi(libc.GoString(arg)))))
			zedit_disp_arg5(d)
		case 'P':
			fallthrough
		case 'E':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(MIN(MAX_DUPLICATES, int64(libc.Atoi(libc.GoString(arg)))))
			zedit_disp_arg3(d)
		case 'V':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(libc.Atoi(libc.GoString(arg)))
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(real_room(d.Olc.Number))
			write_to_output(d, libc.CString("Enter the global name : "))
			d.Olc.Mode = ZEDIT_SARG1
		case 'T':
			if real_trigger(trig_vnum(libc.Atoi(libc.GoString(arg)))) != trig_rnum(-1) {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(real_trigger(trig_vnum(libc.Atoi(libc.GoString(arg)))))
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(real_room(d.Olc.Number))
				zedit_disp_menu(d)
			} else {
				write_to_output(d, libc.CString("That trigger does not exist, try again : "))
			}
		case 'D':
			pos = libc.Atoi(libc.GoString(arg))
			if pos < 0 || pos > NUM_OF_DIRS {
				write_to_output(d, libc.CString("Try again : "))
			} else {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(pos)
				zedit_disp_arg3(d)
			}
		case 'R':
			if (func() int {
				pos = int(real_object(obj_vnum(libc.Atoi(libc.GoString(arg)))))
				return pos
			}()) != int(-1) {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg2 = vnum(pos)
				zedit_disp_menu(d)
			} else {
				write_to_output(d, libc.CString("That object does not exist, try again : "))
			}
		default:
			cleanup_olc(d, CLEANUP_ALL)
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): case ARG2: Ack!"))
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
	case ZEDIT_ARG3:
		if !unicode.IsDigit(rune(*arg)) {
			write_to_output(d, libc.CString("Must be a numeric value, try again : "))
			return
		}
		switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
		case 'E':
			pos = libc.Atoi(libc.GoString(arg))
			for *equipment_types[i] != '\n' {
				i++
			}
			if pos < 0 || pos > i {
				write_to_output(d, libc.CString("Try again : "))
			} else {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(pos)
				zedit_disp_arg5(d)
			}
		case 'P':
			if (func() int {
				pos = int(real_object(obj_vnum(libc.Atoi(libc.GoString(arg)))))
				return pos
			}()) != int(-1) {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(pos)
				zedit_disp_arg5(d)
			} else {
				write_to_output(d, libc.CString("That object does not exist, try again : "))
			}
		case 'D':
			pos = libc.Atoi(libc.GoString(arg))
			if pos < 0 || pos > 2 {
				write_to_output(d, libc.CString("Try again : "))
			} else {
				(d.Olc.Zone.Cmd[d.Olc.Value]).Arg3 = vnum(pos)
				zedit_disp_menu(d)
			}
		case 'M':
			fallthrough
		case 'O':
			fallthrough
		case 'G':
			fallthrough
		case 'R':
			fallthrough
		case 'T':
			fallthrough
		case 'V':
			fallthrough
		default:
			cleanup_olc(d, CLEANUP_ALL)
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): case ARG3: Ack!"))
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
	case ZEDIT_ARG4:
		if *arg == 0 {
			arg = libc.CString("0")
		}
		if !unicode.IsDigit(rune(*arg)) {
			write_to_output(d, libc.CString("Must be a numeric value, try again : "))
			return
		}
		switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
		case 'M':
			fallthrough
		case 'O':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg4 = vnum(MIN(MAX_FROM_ROOM, int64(libc.Atoi(libc.GoString(arg)))))
			zedit_disp_arg5(d)
		case 'G':
			fallthrough
		case 'P':
			fallthrough
		case 'E':
			fallthrough
		case 'V':
			fallthrough
		case 'T':
			fallthrough
		default:
			cleanup_olc(d, CLEANUP_ALL)
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): case ARG4: Ack!"))
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
	case ZEDIT_ARG5:
		if !unicode.IsDigit(rune(*arg)) {
			write_to_output(d, libc.CString("Must be a numeric value, try again : "))
			return
		}
		if libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 100 {
			write_to_output(d, libc.CString("Value must be from 0 to 100"))
			return
		}
		switch (d.Olc.Zone.Cmd[d.Olc.Value]).Command {
		case 'M':
			fallthrough
		case 'O':
			fallthrough
		case 'G':
			fallthrough
		case 'P':
			fallthrough
		case 'E':
			(d.Olc.Zone.Cmd[d.Olc.Value]).Arg5 = vnum(libc.Atoi(libc.GoString(arg)))
			zedit_disp_menu(d)
		case 'V':
			fallthrough
		case 'T':
			fallthrough
		default:
			cleanup_olc(d, CLEANUP_ALL)
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): case ARG5: Ack!"))
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
	case ZEDIT_SARG1:
		if libc.StrLen(arg) != 0 {
			if (d.Olc.Zone.Cmd[d.Olc.Value]).Sarg1 != nil {
				libc.Free(unsafe.Pointer((d.Olc.Zone.Cmd[d.Olc.Value]).Sarg1))
			}
			(d.Olc.Zone.Cmd[d.Olc.Value]).Sarg1 = libc.StrDup(arg)
			d.Olc.Mode = ZEDIT_SARG2
			write_to_output(d, libc.CString("Enter the global value : "))
		} else {
			write_to_output(d, libc.CString("Must have some name to assign : "))
		}
	case ZEDIT_SARG2:
		if libc.StrLen(arg) != 0 {
			(d.Olc.Zone.Cmd[d.Olc.Value]).Sarg2 = libc.StrDup(arg)
			zedit_disp_arg4(d)
		} else {
			write_to_output(d, libc.CString("Must have some value to set it to :"))
		}
	case ZEDIT_ZONE_NAME:
		if genolc_checkstring(d, arg) != 0 {
			if d.Olc.Zone.Name != nil {
				libc.Free(unsafe.Pointer(d.Olc.Zone.Name))
			} else {
				basic_mud_log(libc.CString("SYSERR: OLC: ZEDIT_ZONE_NAME: no name to free!"))
			}
			d.Olc.Zone.Name = libc.StrDup(arg)
			d.Olc.Zone.Number = 1
		}
		zedit_disp_menu(d)
	case ZEDIT_ZONE_BUILDERS:
		if genolc_checkstring(d, arg) != 0 {
			if d.Olc.Zone.Builders != nil {
				libc.Free(unsafe.Pointer(d.Olc.Zone.Builders))
			} else {
				basic_mud_log(libc.CString("SYSERR: OLC: ZEDIT_ZONE_BUILDERS: no builders list to free!"))
			}
			d.Olc.Zone.Builders = libc.StrDup(arg)
			d.Olc.Zone.Number = 1
		}
		zedit_disp_menu(d)
	case ZEDIT_ZONE_RESET:
		pos = libc.Atoi(libc.GoString(arg))
		if !unicode.IsDigit(rune(*arg)) || pos < 0 || pos > 2 {
			write_to_output(d, libc.CString("Try again (0-2) : "))
		} else {
			d.Olc.Zone.Reset_mode = pos
			d.Olc.Zone.Number = 1
			zedit_disp_menu(d)
		}
	case ZEDIT_ZONE_LIFE:
		pos = libc.Atoi(libc.GoString(arg))
		if !unicode.IsDigit(rune(*arg)) || pos < 0 || pos > 240 {
			write_to_output(d, libc.CString("Try again (0-240) : "))
		} else {
			d.Olc.Zone.Lifespan = pos
			d.Olc.Zone.Number = 1
			zedit_disp_menu(d)
		}
	case ZEDIT_ZONE_BOT:
		if d.Olc.Zone_num == 0 {
			d.Olc.Zone.Bot = room_vnum(MIN(int64(d.Olc.Zone.Top), MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
		} else {
			d.Olc.Zone.Bot = room_vnum(MIN(int64(d.Olc.Zone.Top), MAX(int64(libc.Atoi(libc.GoString(arg))), int64(zone_table[d.Olc.Zone_num-1].Top+1))))
		}
		d.Olc.Zone.Number = 1
		zedit_disp_menu(d)
	case ZEDIT_ZONE_TOP:
		if d.Olc.Zone_num == top_of_zone_table {
			d.Olc.Zone.Top = room_vnum(MIN(65000, MAX(int64(libc.Atoi(libc.GoString(arg))), int64(genolc_zonep_bottom(d.Olc.Zone)))))
		} else {
			d.Olc.Zone.Top = room_vnum(MIN(int64(genolc_zone_bottom(d.Olc.Zone_num+1)-1), MAX(int64(libc.Atoi(libc.GoString(arg))), int64(genolc_zonep_bottom(d.Olc.Zone)))))
		}
		d.Olc.Zone.Number = 1
		zedit_disp_menu(d)
	case ZEDIT_ZONE_FLAGS:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > NUM_ZONE_FLAGS {
			write_to_output(d, libc.CString("That is not a valid choice!\r\n"))
			zedit_disp_flag_menu(d)
		} else if number == 0 {
			zedit_disp_menu(d)
			break
		} else {
			TOGGLE_BIT_AR(d.Olc.Zone.Zone_flags[:], bitvector_t(int32(number-1)))
			d.Olc.Zone.Number = 1
			zedit_disp_flag_menu(d)
		}
		return
	case ZEDIT_MIN_LEVEL:
		pos = libc.Atoi(libc.GoString(arg))
		if !unicode.IsDigit(rune(*arg)) || pos < 0 || pos > config_info.Play.Level_cap {
			write_to_output(d, libc.CString("Try again (0 - CONFIG_LEVEL_CAP) : \r\n"))
		} else {
			d.Olc.Zone.Min_level = pos
			d.Olc.Zone.Number = 1
			zedit_disp_menu(d)
		}
	case ZEDIT_MAX_LEVEL:
		pos = libc.Atoi(libc.GoString(arg))
		if !unicode.IsDigit(rune(*arg)) || pos < 0 || pos > config_info.Play.Level_cap {
			write_to_output(d, libc.CString("Try again (0 - CONFIG_LEVEL_CAP) : \r\n"))
		} else if d.Olc.Zone.Min_level > pos {
			write_to_output(d, libc.CString("Max level can not be lower the Min level! :\r\n"))
		} else {
			d.Olc.Zone.Max_level = pos
			d.Olc.Zone.Number = 1
			zedit_disp_menu(d)
		}
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: zedit_parse(): Reached default case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
}
