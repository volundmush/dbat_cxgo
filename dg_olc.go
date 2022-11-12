package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const NUM_TRIG_TYPE_FLAGS = 20
const TRIGEDIT_MAIN_MENU = 0
const TRIGEDIT_TRIGTYPE = 1
const TRIGEDIT_CONFIRM_SAVESTRING = 2
const TRIGEDIT_NAME = 3
const TRIGEDIT_INTENDED = 4
const TRIGEDIT_TYPES = 5
const TRIGEDIT_COMMANDS = 6
const TRIGEDIT_NARG = 7
const TRIGEDIT_ARGUMENT = 8
const TRIGEDIT_COPY = 9
const OLC_SCRIPT_EDIT = 82766
const SCRIPT_MAIN_MENU = 0
const SCRIPT_NEW_TRIGGER = 1
const SCRIPT_DEL_TRIGGER = 2

func do_oasis_trigedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int
		real_num int
		d        *descriptor_data
	)
	skip_spaces(&argument)
	if *argument == 0 || (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument)))))&int(uint16(int16(_ISdigit)))) == 0 {
		send_to_char(ch, libc.CString("Specify a trigger VNUM to edit.\r\n"))
		return
	}
	number = libc.Atoi(libc.GoString(argument))
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_TRIGEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That trigger is currently being edited by %s.\r\n"), GET_NAME(d.Character))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_trigedit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if (func() zone_rnum {
		p := &d.Olc.Zone_num
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(number))
		return *p
	}()) == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for that number!\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 {
		send_cannot_edit(ch, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() int {
		real_num = int(real_trigger(trig_vnum(number)))
		return real_num
	}()) == int(-1) {
		trigedit_setup_new(d)
	} else {
		trigedit_setup_existing(d, real_num)
	}
	var disp int = 0
	if disp == 0 {
		trigedit_disp_menu(d)
		d.Connected = CON_TRIGEDIT
		disp = 1
	}
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(1 << (int(PLR_WRITING % 32)))
	mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d [trigger](allowed zone %d)"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
}
func script_save_to_disk(fp *C.FILE, item unsafe.Pointer, type_ int) {
	var t *trig_proto_list
	if type_ == MOB_TRIGGER {
		t = ((*char_data)(item)).Proto_script
	} else if type_ == OBJ_TRIGGER {
		t = ((*obj_data)(item)).Proto_script
	} else if type_ == WLD_TRIGGER {
		t = ((*room_data)(item)).Proto_script
	} else {
		basic_mud_log(libc.CString("SYSERR: Invalid type passed to script_save_to_disk()"))
		return
	}
	for t != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "T %d\n", t.Vnum)
		t = t.Next
	}
}
func trigedit_setup_new(d *descriptor_data) {
	var trig *trig_data
	trig = new(trig_data)
	trig.Nr = -1
	trig.Name = C.strdup(libc.CString("new trigger"))
	trig.Trigger_type = 1 << 6
	d.Olc.Storage = (*byte)(unsafe.Pointer(&make([]int8, MAX_CMD_LENGTH)[0]))
	C.strncpy(d.Olc.Storage, libc.CString("%echo% This trigger commandlist is not complete!\r\n"), uint64(int(MAX_CMD_LENGTH-1)))
	trig.Narg = 100
	d.Olc.Trig = trig
	d.Olc.Value = 0
}
func trigedit_setup_existing(d *descriptor_data, rtrg_num int) {
	var (
		trig *trig_data
		c    *cmdlist_element
	)
	trig = new(trig_data)
	trig_data_copy(trig, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rtrg_num)))).Proto)
	c = trig.Cmdlist
	d.Olc.Storage = (*byte)(unsafe.Pointer(&make([]int8, MAX_CMD_LENGTH)[0]))
	C.strcpy(d.Olc.Storage, libc.CString(""))
	for c != nil {
		C.strcat(d.Olc.Storage, c.Cmd)
		C.strcat(d.Olc.Storage, libc.CString("\r\n"))
		c = c.Next
	}
	d.Olc.Trig = trig
	d.Olc.Value = 0
}
func trigedit_disp_menu(d *descriptor_data) {
	var (
		trig        *trig_data = d.Olc.Trig
		attach_type *byte
		trgtypes    [256]byte
	)
	if trig.Attach_type == OBJ_TRIGGER {
		attach_type = libc.CString("Objects")
		sprintbit(bitvector_t(trig.Trigger_type), otrig_types[:], &trgtypes[0], uint64(256))
	} else if trig.Attach_type == WLD_TRIGGER {
		attach_type = libc.CString("Rooms")
		sprintbit(bitvector_t(trig.Trigger_type), wtrig_types[:], &trgtypes[0], uint64(256))
	} else {
		attach_type = libc.CString("Mobiles")
		sprintbit(bitvector_t(trig.Trigger_type), trig_types[:], &trgtypes[0], uint64(256))
	}
	clear_screen(d)
	write_to_output(d, libc.CString("Trigger Editor [@c%d@n]\r\n\r\n@g1@n) Name         : @y%s\r\n@g2@n) Intended for : @y%s\r\n@g3@n) Trigger types: @y%s\r\n@g4@n) Numeric Arg  : @y%d\r\n@g5@n) Arguments    : @y%s\r\n@g6@n) Commands:\r\n@c%s\r\n@gW@n) Copy Trigger\r\n@gZ@n) Wiznet\r\n@gQ@n) Quit\r\nEnter Choice :"), d.Olc.Number, trig.Name, attach_type, &trgtypes[0], trig.Narg, func() *byte {
		if trig.Arglist != nil {
			return trig.Arglist
		}
		return libc.CString("")
	}(), d.Olc.Storage)
	d.Olc.Mode = TRIGEDIT_MAIN_MENU
}
func trigedit_disp_types(d *descriptor_data) {
	var (
		i       int
		columns int = 0
		types   **byte
		bitbuf  [64936]byte
	)
	switch d.Olc.Trig.Attach_type {
	case WLD_TRIGGER:
		types = &wtrig_types[0]
	case OBJ_TRIGGER:
		types = &otrig_types[0]
	case MOB_TRIGGER:
		fallthrough
	default:
		types = &trig_types[0]
	}
	clear_screen(d)
	for i = 0; i < NUM_TRIG_TYPE_FLAGS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s  %s"), i+1, *(**byte)(unsafe.Add(unsafe.Pointer(types), unsafe.Sizeof((*byte)(nil))*uintptr(i))), func() string {
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
	sprintbit(bitvector_t(d.Olc.Trig.Trigger_type), ([0]*byte)(types), &bitbuf[0], uint64(64936))
	write_to_output(d, libc.CString("\r\nCurrent types : @c%s@n\r\nEnter type (0 to quit) : "), &bitbuf[0])
}
func trigedit_parse(d *descriptor_data, arg *byte) {
	var i int = 0
	switch d.Olc.Mode {
	case TRIGEDIT_MAIN_MENU:
		switch C.tolower(int(*arg)) {
		case 'q':
			if d.Olc.Value != 0 {
				if d.Olc.Trig.Trigger_type == 0 {
					write_to_output(d, libc.CString("Invalid Trigger Type! Answer a to abort quit!\r\n"))
				}
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = TRIGEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '1':
			d.Olc.Mode = TRIGEDIT_NAME
			write_to_output(d, libc.CString("Name: "))
		case '2':
			d.Olc.Mode = TRIGEDIT_INTENDED
			write_to_output(d, libc.CString("0: Mobiles, 1: Objects, 2: Rooms: "))
		case '3':
			d.Olc.Mode = TRIGEDIT_TYPES
			trigedit_disp_types(d)
		case '4':
			d.Olc.Mode = TRIGEDIT_NARG
			write_to_output(d, libc.CString("Numeric argument: "))
		case '5':
			d.Olc.Mode = TRIGEDIT_ARGUMENT
			write_to_output(d, libc.CString("Argument: "))
		case '6':
			d.Olc.Mode = TRIGEDIT_COMMANDS
			write_to_output(d, libc.CString("Enter trigger commands: (/s saves /h for help)\r\n\r\n"))
			d.Backstr = nil
			if d.Olc.Storage != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Storage)
				d.Backstr = C.strdup(d.Olc.Storage)
			}
			d.Str = &d.Olc.Storage
			d.Max_str = MAX_CMD_LENGTH
			d.Mail_to = 0
			d.Olc.Value = 1
		case 'w':
			fallthrough
		case 'W':
			write_to_output(d, libc.CString("Copy what trigger? "))
			d.Olc.Mode = TRIGEDIT_COPY
		case 'Z':
			fallthrough
		case 'z':
			search_replace(arg, libc.CString("z "), libc.CString(""))
			do_wiznet(d.Character, arg, 0, 0)
		default:
			trigedit_disp_menu(d)
			return
		}
		return
	case TRIGEDIT_CONFIRM_SAVESTRING:
		switch C.tolower(int(*arg)) {
		case 'y':
			trigedit_save(d)
			mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s edits trigger %d"), GET_NAME(d.Character), d.Olc.Number)
			fallthrough
		case 'n':
			cleanup_olc(d, CLEANUP_ALL)
			return
		case 'a':
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\n"))
			write_to_output(d, libc.CString("Do you wish to save the trigger? : "))
			return
		}
	case TRIGEDIT_NAME:
		smash_tilde(arg)
		if d.Olc.Trig.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Trig.Name))
		}
		d.Olc.Trig.Name = C.strdup(func() *byte {
			if arg != nil && *arg != 0 {
				return arg
			}
			return libc.CString("undefined")
		}())
		d.Olc.Value++
	case TRIGEDIT_INTENDED:
		if libc.Atoi(libc.GoString(arg)) >= MOB_TRIGGER || libc.Atoi(libc.GoString(arg)) <= WLD_TRIGGER {
			d.Olc.Trig.Attach_type = int8(libc.Atoi(libc.GoString(arg)))
		}
		d.Olc.Value++
	case TRIGEDIT_NARG:
		d.Olc.Trig.Narg = MIN(100, MAX(libc.Atoi(libc.GoString(arg)), 0))
		d.Olc.Value++
	case TRIGEDIT_ARGUMENT:
		smash_tilde(arg)
		if *arg != 0 {
			d.Olc.Trig.Arglist = C.strdup(arg)
		} else {
			d.Olc.Trig.Arglist = nil
		}
		d.Olc.Value++
	case TRIGEDIT_TYPES:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) == 0 {
			break
		} else if i >= 0 && i <= NUM_TRIG_TYPE_FLAGS {
			d.Olc.Trig.Trigger_type ^= 1 << (i - 1)
		}
		d.Olc.Value++
		trigedit_disp_types(d)
		return
	case TRIGEDIT_COMMANDS:
	case TRIGEDIT_COPY:
		if (func() int {
			i = int(real_trigger(trig_vnum(libc.Atoi(libc.GoString(arg)))))
			return i
		}()) != int(-1) {
			trigedit_setup_existing(d, i)
		} else {
			write_to_output(d, libc.CString("That trigger does not exist.\r\n"))
		}
	}
	d.Olc.Mode = TRIGEDIT_MAIN_MENU
	trigedit_disp_menu(d)
}
func trigedit_save(d *descriptor_data) {
	var (
		i         int
		rnum      trig_rnum
		found     int = 0
		s         *byte
		proto     *trig_data
		trig      *trig_data = d.Olc.Trig
		live_trig *trig_data
		cmd       *cmdlist_element
		next_cmd  *cmdlist_element
		new_index **index_data
		dsc       *descriptor_data
		trig_file *C.FILE
		zone      int
		top       int
		buf       [16384]byte
		bitBuf    [2048]byte
		fname     [2048]byte
	)
	if (func() trig_rnum {
		rnum = real_trigger(trig_vnum(d.Olc.Number))
		return rnum
	}()) != trig_rnum(-1) {
		proto = (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto
		for cmd = proto.Cmdlist; cmd != nil; cmd = next_cmd {
			next_cmd = cmd.Next
			if cmd.Cmd != nil {
				libc.Free(unsafe.Pointer(cmd.Cmd))
			}
			libc.Free(unsafe.Pointer(cmd))
		}
		libc.Free(unsafe.Pointer(proto.Arglist))
		libc.Free(unsafe.Pointer(proto.Name))
		s = d.Olc.Storage
		trig.Cmdlist = new(cmdlist_element)
		if s != nil {
			var t *byte = strtok(s, libc.CString("\n\r"))
			if t != nil {
				trig.Cmdlist.Cmd = C.strdup(t)
			} else {
				trig.Cmdlist.Cmd = C.strdup(libc.CString("* No script"))
			}
			cmd = trig.Cmdlist
			for (func() *byte {
				s = strtok(nil, libc.CString("\n\r"))
				return s
			}()) != nil {
				cmd.Next = new(cmdlist_element)
				cmd = cmd.Next
				cmd.Cmd = C.strdup(s)
			}
		} else {
			trig.Cmdlist.Cmd = C.strdup(libc.CString("* No Script"))
		}
		trig_data_copy(proto, trig)
		live_trig = trigger_list
		for live_trig != nil {
			if live_trig.Nr == int64(rnum) {
				if live_trig.Arglist != nil {
					libc.Free(unsafe.Pointer(live_trig.Arglist))
					live_trig.Arglist = nil
				}
				if live_trig.Name != nil {
					libc.Free(unsafe.Pointer(live_trig.Name))
					live_trig.Name = nil
				}
				if proto.Arglist != nil {
					live_trig.Arglist = C.strdup(proto.Arglist)
				}
				if proto.Name != nil {
					live_trig.Name = C.strdup(proto.Name)
				}
				if live_trig.Wait_event != nil {
					event_cancel(live_trig.Wait_event)
					live_trig.Wait_event = nil
				}
				if live_trig.Var_list != nil {
					free_varlist(live_trig.Var_list)
					live_trig.Var_list = nil
				}
				live_trig.Cmdlist = proto.Cmdlist
				live_trig.Curr_state = live_trig.Cmdlist
				live_trig.Trigger_type = proto.Trigger_type
				live_trig.Attach_type = proto.Attach_type
				live_trig.Narg = proto.Narg
				live_trig.Data_type = proto.Data_type
				live_trig.Depth = 0
			}
			live_trig = live_trig.Next_in_world
		}
	} else {
		new_index = &make([]*index_data, top_of_trigt+2)[0]
		s = d.Olc.Storage
		trig.Cmdlist = new(cmdlist_element)
		if s != nil {
			var t *byte = strtok(s, libc.CString("\n\r"))
			trig.Cmdlist.Cmd = C.strdup(func() *byte {
				if t != nil {
					return t
				}
				return libc.CString("* No script")
			}())
			cmd = trig.Cmdlist
			for (func() *byte {
				s = strtok(nil, libc.CString("\n\r"))
				return s
			}()) != nil {
				cmd.Next = new(cmdlist_element)
				cmd = cmd.Next
				cmd.Cmd = C.strdup(s)
			}
		} else {
			trig.Cmdlist.Cmd = C.strdup(libc.CString("* No Script"))
		}
		for i = 0; i < top_of_trigt; i++ {
			if found == 0 {
				if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum > mob_vnum(d.Olc.Number) {
					found = TRUE
					rnum = trig_rnum(i)
					*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum))) = new(index_data)
					d.Olc.Trig.Nr = int64(rnum)
					(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Vnum = mob_vnum(d.Olc.Number)
					(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Number = 0
					(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Func = nil
					proto = new(trig_data)
					(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto = proto
					trig_data_copy(proto, trig)
					*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum+1))) = *(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))
					proto = (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto
					proto.Nr = int64(rnum + 1)
				} else {
					*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i))) = *(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))
				}
			} else {
				*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i+1))) = *(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))
				proto = (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto
				proto.Nr = int64(i + 1)
			}
		}
		if found == 0 {
			rnum = trig_rnum(i)
			*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum))) = new(index_data)
			d.Olc.Trig.Nr = int64(rnum)
			(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Vnum = mob_vnum(d.Olc.Number)
			(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Number = 0
			(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Func = nil
			proto = new(trig_data)
			(*(**index_data)(unsafe.Add(unsafe.Pointer(new_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto = proto
			trig_data_copy(proto, trig)
		}
		libc.Free(unsafe.Pointer(trig_index))
		trig_index = new_index
		top_of_trigt++
		for live_trig = trigger_list; live_trig != nil; live_trig = live_trig.Next_in_world {
			live_trig.Nr += int64(libc.BoolToInt(live_trig.Nr != int64(-1) && live_trig.Nr > int64(rnum)))
		}
		for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
			if dsc.Connected == CON_TRIGEDIT {
				if dsc.Olc.Trig.Nr >= int64(rnum) {
					dsc.Olc.Trig.Nr++
				}
			}
		}
	}
	zone = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
	top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Top)
	stdio.Snprintf(&fname[0], int(2048), "%s/%i.new", LIB_WORLD, zone)
	if (func() *C.FILE {
		trig_file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fname[0]), "w")))
		return trig_file
	}()) == nil {
		mudlog(BRF, MAX(ADMLVL_GOD, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("SYSERR: OLC: Can't open trig file \"%s\""), &fname[0])
		return
	}
	for i = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Bot); i <= top; i++ {
		if (func() trig_rnum {
			rnum = real_trigger(trig_vnum(i))
			return rnum
		}()) != trig_rnum(-1) {
			trig = (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(rnum)))).Proto
			if stdio.Fprintf((*stdio.File)(unsafe.Pointer(trig_file)), "#%d\n", i) < 0 {
				mudlog(BRF, MAX(ADMLVL_GOD, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("SYSERR: OLC: Can't write trig file!"))
				C.fclose(trig_file)
				return
			}
			sprintascii(&bitBuf[0], bitvector_t(trig.Trigger_type))
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(trig_file)), "%s%c\n%d %s %d\n%s%c\n", func() *byte {
				if trig.Name != nil {
					return trig.Name
				}
				return libc.CString("unknown trigger")
			}(), STRING_TERMINATOR, trig.Attach_type, &func() [2048]byte {
				if bitBuf[0] != 0 {
					return bitBuf
				}
				return func() [2048]byte {
					var t [2048]byte
					copy(t[:], []byte("0"))
					return t
				}()
			}()[0], trig.Narg, func() *byte {
				if trig.Arglist != nil {
					return trig.Arglist
				}
				return libc.CString("")
			}(), STRING_TERMINATOR)
			C.strcpy(&buf[0], libc.CString(""))
			for cmd = trig.Cmdlist; cmd != nil; cmd = cmd.Next {
				C.strcat(&buf[0], cmd.Cmd)
				C.strcat(&buf[0], libc.CString("\n"))
			}
			if buf[0] == 0 {
				C.strcpy(&buf[0], libc.CString("* Empty script"))
			}
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(trig_file)), "%s%c\n", &buf[0], STRING_TERMINATOR)
			buf[0] = '\x00'
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(trig_file)), "$%c\n", STRING_TERMINATOR)
	C.fclose(trig_file)
	stdio.Snprintf(&buf[0], int(16384), "%s%d.trg", LIB_WORLD, zone)
	stdio.Remove(libc.GoString(&buf[0]))
	stdio.Rename(libc.GoString(&fname[0]), libc.GoString(&buf[0]))
	write_to_output(d, libc.CString("Trigger saved to disk.\r\n"))
	create_world_index(zone, libc.CString("trg"))
}
func dg_olc_script_copy(d *descriptor_data) {
	var (
		origscript *trig_proto_list
		editscript *trig_proto_list
	)
	if d.Olc.Item_type == MOB_TRIGGER {
		origscript = d.Olc.Mob.Proto_script
	} else if d.Olc.Item_type == OBJ_TRIGGER {
		origscript = d.Olc.Obj.Proto_script
	} else {
		origscript = d.Olc.Room.Proto_script
	}
	if origscript != nil {
		editscript = new(trig_proto_list)
		d.Olc.Script = editscript
		for origscript != nil {
			editscript.Vnum = origscript.Vnum
			origscript = origscript.Next
			if origscript != nil {
				editscript.Next = new(trig_proto_list)
			}
			editscript = editscript.Next
		}
	} else {
		d.Olc.Script = nil
	}
}
func dg_script_menu(d *descriptor_data) {
	var (
		editscript *trig_proto_list
		i          int = 0
	)
	d.Olc.Mode = OLC_SCRIPT_EDIT
	d.Olc.Script_mode = SCRIPT_MAIN_MENU
	clear_screen(d)
	write_to_output(d, libc.CString("     Script Editor\r\n\r\n     Trigger List:\r\n"))
	editscript = d.Olc.Script
	for editscript != nil {
		write_to_output(d, libc.CString("     %2d) [@c%d@n] @c%s@n"), func() int {
			p := &i
			*p++
			return *p
		}(), editscript.Vnum, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(real_trigger(trig_vnum(editscript.Vnum)))))).Proto.Name)
		if int((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(real_trigger(trig_vnum(editscript.Vnum)))))).Proto.Attach_type) != d.Olc.Item_type {
			write_to_output(d, libc.CString("   @g** Mis-matched Trigger Type **@n\r\n"))
		} else {
			write_to_output(d, libc.CString("\r\n"))
		}
		editscript = editscript.Next
	}
	if i == 0 {
		write_to_output(d, libc.CString("     <none>\r\n"))
	}
	write_to_output(d, libc.CString("\r\n @gN@n)  New trigger for this script\r\n @gD@n)  Delete a trigger in this script\r\n @gX@n)  Exit Script Editor\r\n\r\n     Enter choice :"))
}
func dg_script_edit_parse(d *descriptor_data, arg *byte) int {
	var (
		trig     *trig_proto_list
		currtrig *trig_proto_list
		count    int
		pos      int
		vnum     int
	)
	switch d.Olc.Script_mode {
	case SCRIPT_MAIN_MENU:
		switch C.tolower(int(*arg)) {
		case 'x':
			return 0
		case 'n':
			write_to_output(d, libc.CString("\r\nPlease enter position, vnum   (ex: 1, 200):"))
			d.Olc.Script_mode = SCRIPT_NEW_TRIGGER
		case 'd':
			write_to_output(d, libc.CString("     Which entry should be deleted?  0 to abort :"))
			d.Olc.Script_mode = SCRIPT_DEL_TRIGGER
		default:
			dg_script_menu(d)
		}
		return 1
	case SCRIPT_NEW_TRIGGER:
		vnum = -1
		count = __isoc99_sscanf(arg, libc.CString("%d, %d"), &pos, &vnum)
		if count == 1 {
			vnum = pos
			pos = 999
		}
		if pos <= 0 {
			break
		}
		if vnum == 0 {
			break
		}
		if real_trigger(trig_vnum(vnum)) == trig_rnum(-1) {
			write_to_output(d, libc.CString("Invalid Trigger VNUM!\r\nPlease enter position, vnum   (ex: 1, 200):"))
			return 1
		}
		currtrig = d.Olc.Script
		trig = new(trig_proto_list)
		trig.Vnum = vnum
		if pos == 1 || currtrig == nil {
			trig.Next = d.Olc.Script
			d.Olc.Script = trig
		} else {
			for currtrig.Next != nil && func() int {
				p := &pos
				*p--
				return *p
			}() != 0 {
				currtrig = currtrig.Next
			}
			trig.Next = currtrig.Next
			currtrig.Next = trig
		}
		d.Olc.Value++
	case SCRIPT_DEL_TRIGGER:
		pos = libc.Atoi(libc.GoString(arg))
		if pos <= 0 {
			break
		}
		if pos == 1 && d.Olc.Script != nil {
			d.Olc.Value++
			currtrig = d.Olc.Script
			d.Olc.Script = currtrig.Next
			libc.Free(unsafe.Pointer(currtrig))
			break
		}
		pos--
		currtrig = d.Olc.Script
		for func() int {
			p := &pos
			*p--
			return *p
		}() != 0 && currtrig != nil {
			currtrig = currtrig.Next
		}
		if currtrig != nil && currtrig.Next != nil {
			d.Olc.Value++
			trig = currtrig.Next
			currtrig.Next = trig.Next
			libc.Free(unsafe.Pointer(trig))
		}
	}
	dg_script_menu(d)
	return 1
}
func trigedit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case TRIGEDIT_COMMANDS:
		trigedit_disp_menu(d)
	}
}
func format_script(d *descriptor_data) int {
	var (
		nsc         [16384]byte
		t           *byte
		line        [256]byte
		sc          *byte
		len_        uint64 = 0
		nlen        uint64 = 0
		llen        uint64 = 0
		indent      int    = 0
		indent_next int    = FALSE
		found_case  int    = FALSE
		i           int
		line_num    int = 0
	)
	if d.Str == nil || *d.Str == nil {
		return FALSE
	}
	sc = C.strdup(*d.Str)
	t = strtok(sc, libc.CString("\n\r"))
	nsc[0] = '\x00'
	for t != nil {
		line_num++
		skip_spaces(&t)
		if C.strncasecmp(t, libc.CString("if "), 3) == 0 || C.strncasecmp(t, libc.CString("switch "), 7) == 0 {
			indent_next = TRUE
		} else if C.strncasecmp(t, libc.CString("while "), 6) == 0 {
			found_case = TRUE
			indent_next = TRUE
		} else if C.strncasecmp(t, libc.CString("end"), 3) == 0 || C.strncasecmp(t, libc.CString("done"), 4) == 0 {
			if indent == 0 {
				write_to_output(d, libc.CString("Unmatched 'end' or 'done' (line %d)!\r\n"), line_num)
				libc.Free(unsafe.Pointer(sc))
				return FALSE
			}
			indent--
			indent_next = FALSE
		} else if C.strncasecmp(t, libc.CString("else"), 4) == 0 {
			if indent == 0 {
				write_to_output(d, libc.CString("Unmatched 'else' (line %d)!\r\n"), line_num)
				libc.Free(unsafe.Pointer(sc))
				return FALSE
			}
			indent--
			indent_next = TRUE
		} else if C.strncasecmp(t, libc.CString("case"), 4) == 0 || C.strncasecmp(t, libc.CString("default"), 7) == 0 {
			if indent == 0 {
				write_to_output(d, libc.CString("Case/default outside switch (line %d)!\r\n"), line_num)
				libc.Free(unsafe.Pointer(sc))
				return FALSE
			}
			if found_case == 0 {
				indent_next = TRUE
			}
			found_case = TRUE
		} else if C.strncasecmp(t, libc.CString("break"), 5) == 0 {
			if found_case == 0 || indent == 0 {
				write_to_output(d, libc.CString("Break not in case (line %d)!\r\n"), line_num)
				libc.Free(unsafe.Pointer(sc))
				return FALSE
			}
			found_case = FALSE
			indent--
		}
		line[0] = '\x00'
		for func() int {
			nlen = 0
			return func() int {
				i = 0
				return i
			}()
		}(); i < indent; i++ {
			C.strncat(&line[0], libc.CString("  "), uint64(256-1))
			nlen += 2
		}
		llen = uint64(stdio.Snprintf(&line[nlen], int(256-uintptr(nlen)), "%s\r\n", t))
		if llen < 0 || llen+nlen+len_ > d.Max_str-1 {
			write_to_output(d, libc.CString("String too long, formatting aborted\r\n"))
			libc.Free(unsafe.Pointer(sc))
			return FALSE
		}
		len_ = len_ + nlen + llen
		C.strcat(&nsc[0], &line[0])
		if indent_next != 0 {
			indent++
			indent_next = FALSE
		}
		t = strtok(nil, libc.CString("\n\r"))
	}
	if indent != 0 {
		write_to_output(d, libc.CString("Unmatched if, while or switch ignored.\r\n"))
	}
	libc.Free(unsafe.Pointer(*d.Str))
	*d.Str = C.strdup(&nsc[0])
	libc.Free(unsafe.Pointer(sc))
	return TRUE
}