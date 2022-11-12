package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func do_oasis_aedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		d   *descriptor_data
		i   int
	)
	if config_info.Operation.Use_new_socials == 0 {
		send_to_char(ch, libc.CString("Socials cannot be edited at the moment.\r\n"))
		return
	}
	if ch.Player_specials.Olc_zone != AEDIT_PERMISSION && ch.Admlevel < ADMLVL_BUILDER {
		send_to_char(ch, libc.CString("You don't have access to editing socials.\r\n"))
		return
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_AEDIT {
			send_to_char(ch, libc.CString("Sorry, only one can edit socials at a time.\r\n"))
			return
		}
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Please specify a social to edit.\r\n"))
		return
	}
	d = ch.Desc
	if C.strcasecmp(libc.CString("save"), &arg[0]) == 0 {
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves socials."), GET_NAME(ch))
		send_to_char(ch, libc.CString("Writing social file..\r\n"))
		aedit_save_to_disk(d)
		send_to_char(ch, libc.CString("Done.\r\n"))
		return
	}
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	d.Olc.Number = 0
	d.Olc.Storage = C.strdup(&arg[0])
	for d.Olc.Zone_num = 0; d.Olc.Zone_num <= zone_rnum(top_of_socialt); d.Olc.Zone_num++ {
		if is_abbrev(d.Olc.Storage, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Command) != 0 {
			break
		}
	}
	if d.Olc.Zone_num > zone_rnum(top_of_socialt) {
		if (func() int {
			i = aedit_find_command(d.Olc.Storage)
			return i
		}()) != -1 {
			send_to_char(ch, libc.CString("The '%s' command already exists (%s).\r\n"), d.Olc.Storage, (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command)
			cleanup_olc(d, CLEANUP_ALL)
			return
		}
		send_to_char(ch, libc.CString("Do you wish to add the '%s' action? "), d.Olc.Storage)
		d.Olc.Mode = AEDIT_CONFIRM_ADD
	} else {
		send_to_char(ch, libc.CString("Do you wish to edit the '%s' action? "), (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Command)
		d.Olc.Mode = AEDIT_CONFIRM_EDIT
	}
	d.Connected = CON_AEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(1 << (int(PLR_WRITING % 32)))
	mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing actions."), GET_NAME(ch))
}
func aedit_setup_new(d *descriptor_data) {
	d.Olc.Action = new(social_messg)
	d.Olc.Action.Command = C.strdup(d.Olc.Storage)
	d.Olc.Action.Sort_as = C.strdup(d.Olc.Storage)
	d.Olc.Action.Hide = 0
	d.Olc.Action.Min_victim_position = POS_STANDING
	d.Olc.Action.Min_char_position = POS_STANDING
	d.Olc.Action.Min_level_char = 0
	d.Olc.Action.Char_no_arg = C.strdup(libc.CString("This action is unfinished."))
	d.Olc.Action.Others_no_arg = C.strdup(libc.CString("This action is unfinished."))
	d.Olc.Action.Char_found = nil
	d.Olc.Action.Others_found = nil
	d.Olc.Action.Vict_found = nil
	d.Olc.Action.Not_found = nil
	d.Olc.Action.Char_auto = nil
	d.Olc.Action.Others_auto = nil
	d.Olc.Action.Char_body_found = nil
	d.Olc.Action.Others_body_found = nil
	d.Olc.Action.Vict_body_found = nil
	d.Olc.Action.Char_obj_found = nil
	d.Olc.Action.Others_obj_found = nil
	aedit_disp_menu(d)
	d.Olc.Value = 0
}
func aedit_setup_existing(d *descriptor_data, real_num int) {
	d.Olc.Action = new(social_messg)
	d.Olc.Action.Command = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Command)
	d.Olc.Action.Sort_as = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Sort_as)
	d.Olc.Action.Hide = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Hide
	d.Olc.Action.Min_victim_position = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Min_victim_position
	d.Olc.Action.Min_char_position = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Min_char_position
	d.Olc.Action.Min_level_char = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Min_level_char
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_no_arg != nil {
		d.Olc.Action.Char_no_arg = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_no_arg)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_no_arg != nil {
		d.Olc.Action.Others_no_arg = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_no_arg)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_found != nil {
		d.Olc.Action.Char_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_found != nil {
		d.Olc.Action.Others_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Vict_found != nil {
		d.Olc.Action.Vict_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Vict_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Not_found != nil {
		d.Olc.Action.Not_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Not_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_auto != nil {
		d.Olc.Action.Char_auto = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_auto)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_auto != nil {
		d.Olc.Action.Others_auto = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_auto)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_body_found != nil {
		d.Olc.Action.Char_body_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_body_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_body_found != nil {
		d.Olc.Action.Others_body_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_body_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Vict_body_found != nil {
		d.Olc.Action.Vict_body_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Vict_body_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_obj_found != nil {
		d.Olc.Action.Char_obj_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Char_obj_found)
	}
	if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_obj_found != nil {
		d.Olc.Action.Others_obj_found = C.strdup((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(real_num)))).Others_obj_found)
	}
	d.Olc.Value = 0
	aedit_disp_menu(d)
}
func aedit_save_internally(d *descriptor_data) {
	var (
		new_soc_mess_list *social_messg = nil
		i                 int
	)
	if d.Olc.Zone_num > zone_rnum(top_of_socialt) {
		new_soc_mess_list = &make([]social_messg, top_of_socialt+2)[0]
		for i = 0; i <= top_of_socialt; i++ {
			*(*social_messg)(unsafe.Add(unsafe.Pointer(new_soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i))) = *(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))
		}
		*(*social_messg)(unsafe.Add(unsafe.Pointer(new_soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(func() int {
			p := &top_of_socialt
			*p++
			return *p
		}()))) = *d.Olc.Action
		libc.Free(unsafe.Pointer(soc_mess_list))
		soc_mess_list = new_soc_mess_list
	} else {
		i = aedit_find_command(d.Olc.Action.Command)
		d.Olc.Action.Act_nr = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Act_nr
		free_action((*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num))))
		*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num))) = *d.Olc.Action
	}
	create_command_list()
	sort_commands()
	add_to_save_list(AEDIT_PERMISSION, int(SL_GLD+1))
	aedit_save_to_disk(d)
}
func aedit_save_to_disk(d *descriptor_data) {
	var (
		fp *C.FILE
		i  int
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_MISC, "w+")))
		return fp
	}()) == nil {
		var error [64936]byte
		stdio.Snprintf(&error[0], int(64936), "Can't open socials file '%s'", LIB_MISC)
		C.perror(&error[0])
		C.exit(1)
	}
	for i = 0; i <= top_of_socialt; i++ {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "~%s %s %d %d %d %d\n", (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Command, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Sort_as, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Hide, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_char_position, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_victim_position, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_level_char)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n%s\n%s\n%s\n", func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_no_arg != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_no_arg
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_no_arg != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_no_arg
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_found
			}
			return libc.CString("#")
		}())
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n%s\n%s\n%s\n", func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Not_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Not_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_auto != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_auto
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_auto != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_auto
			}
			return libc.CString("#")
		}())
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n%s\n%s\n", func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_body_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_body_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_body_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_body_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_body_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_body_found
			}
			return libc.CString("#")
		}())
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n%s\n\n", func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_obj_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_obj_found
			}
			return libc.CString("#")
		}(), func() *byte {
			if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_obj_found != nil {
				return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_obj_found
			}
			return libc.CString("#")
		}())
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "$\n")
	C.fclose(fp)
	remove_from_save_list(AEDIT_PERMISSION, int(SL_GLD+1))
}
func aedit_disp_menu(d *descriptor_data) {
	var action *social_messg = d.Olc.Action
	write_to_output(d, libc.CString("@n-- Action editor\r\n@gn@n) Command         : @y%-15.15s@n @g1@n) Sort as Command  : @y%-15.15s@n\r\n@g2@n) Min Position[CH]: @c%-8.8s        @g3@n) Min Position [VT]: @c%-8.8s\r\n@g4@n) Min Level   [CH]: @c%-3d             @g5@n) Show if Invisible: @c%s\r\n@ga@n) Char    [NO ARG]: @c%s\r\n@gb@n) Others  [NO ARG]: @c%s\r\n@gc@n) Char [NOT FOUND]: @c%s\r\n@gd@n) Char  [ARG SELF]: @c%s\r\n@ge@n) Others[ARG SELF]: @c%s\r\n@gf@n) Char      [VICT]: @c%s\r\n@gg@n) Others    [VICT]: @c%s\r\n@gh@n) Victim    [VICT]: @c%s\r\n@gi@n) Char  [BODY PRT]: @c%s\r\n@gj@n) Others[BODY PRT]: @c%s\r\n@gk@n) Victim[BODY PRT]: @c%s\r\n@gl@n) Char       [OBJ]: @c%s\r\n@gm@n) Others     [OBJ]: @c%s\r\n@gq@n) Quit\r\nEnter Choice:"), action.Command, action.Sort_as, position_types[action.Min_char_position], position_types[action.Min_victim_position], action.Min_level_char, func() string {
		if action.Hide != 0 {
			return "HIDDEN"
		}
		return "NOT HIDDEN"
	}(), func() *byte {
		if action.Char_no_arg != nil {
			return action.Char_no_arg
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Others_no_arg != nil {
			return action.Others_no_arg
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Not_found != nil {
			return action.Not_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Char_auto != nil {
			return action.Char_auto
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Others_auto != nil {
			return action.Others_auto
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Char_found != nil {
			return action.Char_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Others_found != nil {
			return action.Others_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Vict_found != nil {
			return action.Vict_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Char_body_found != nil {
			return action.Char_body_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Others_body_found != nil {
			return action.Others_body_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Vict_body_found != nil {
			return action.Vict_body_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Char_obj_found != nil {
			return action.Char_obj_found
		}
		return libc.CString("<Null>")
	}(), func() *byte {
		if action.Others_obj_found != nil {
			return action.Others_obj_found
		}
		return libc.CString("<Null>")
	}())
	d.Olc.Mode = AEDIT_MAIN_MENU
}
func aedit_parse(d *descriptor_data, arg *byte) {
	var i int
	switch d.Olc.Mode {
	case AEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			aedit_save_internally(d)
			mudlog(CMP, ADMLVL_IMPL, TRUE, libc.CString("OLC: %s edits action %s"), GET_NAME(d.Character), d.Olc.Action.Command)
			cleanup_olc(d, CLEANUP_STRUCTS)
			write_to_output(d, libc.CString("Action saved to disk.\r\n"))
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to save your changes? : "))
		}
		return
	case AEDIT_CONFIRM_EDIT:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			aedit_setup_existing(d, int(d.Olc.Zone_num))
		case 'q':
			fallthrough
		case 'Q':
			cleanup_olc(d, CLEANUP_ALL)
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Zone_num++
			for ; d.Olc.Zone_num <= zone_rnum(top_of_socialt); d.Olc.Zone_num++ {
				if is_abbrev(d.Olc.Storage, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Command) != 0 {
					break
				}
			}
			if d.Olc.Zone_num > zone_rnum(top_of_socialt) {
				if aedit_find_command(d.Olc.Storage) != -1 {
					cleanup_olc(d, CLEANUP_ALL)
					break
				}
				write_to_output(d, libc.CString("Do you wish to add the '%s' action? "), d.Olc.Storage)
				d.Olc.Mode = AEDIT_CONFIRM_ADD
			} else {
				write_to_output(d, libc.CString("Do you wish to edit the '%s' action? "), (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Command)
				d.Olc.Mode = AEDIT_CONFIRM_EDIT
			}
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to edit the '%s' action? "), (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(d.Olc.Zone_num)))).Command)
		}
		return
	case AEDIT_CONFIRM_ADD:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			aedit_setup_new(d)
		case 'n':
			fallthrough
		case 'N':
			fallthrough
		case 'q':
			fallthrough
		case 'Q':
			cleanup_olc(d, CLEANUP_ALL)
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to add the '%s' action? "), d.Olc.Storage)
		}
		return
	case AEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = AEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
		case 'n':
			write_to_output(d, libc.CString("Enter action name: "))
			d.Olc.Mode = AEDIT_ACTION_NAME
			return
		case '1':
			write_to_output(d, libc.CString("Enter sort info for this action (for the command listing): "))
			d.Olc.Mode = AEDIT_SORT_AS
			return
		case '2':
			write_to_output(d, libc.CString("Enter the minimum position the Character has to be in to activate social:\r\n"))
			for i = POS_DEAD; i <= POS_STANDING; i++ {
				write_to_output(d, libc.CString("   %d) %s\r\n"), i, position_types[i])
			}
			write_to_output(d, libc.CString("Enter choice: "))
			d.Olc.Mode = AEDIT_MIN_CHAR_POS
			return
		case '3':
			write_to_output(d, libc.CString("Enter the minimum position the Victim has to be in to activate social:\r\n"))
			for i = POS_DEAD; i <= POS_STANDING; i++ {
				write_to_output(d, libc.CString("   %d) %s\r\n"), i, position_types[i])
			}
			write_to_output(d, libc.CString("Enter choice: "))
			d.Olc.Mode = AEDIT_MIN_VICT_POS
			return
		case '4':
			write_to_output(d, libc.CString("Enter new minimum level for social: "))
			d.Olc.Mode = AEDIT_MIN_CHAR_LEVEL
			return
		case '5':
			d.Olc.Action.Hide = int(libc.BoolToInt(d.Olc.Action.Hide == 0))
			aedit_disp_menu(d)
			d.Olc.Value = 1
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter social shown to the Character when there is no argument supplied.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Char_no_arg != nil {
					return d.Olc.Action.Char_no_arg
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_NOVICT_CHAR
			return
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter social shown to Others when there is no argument supplied.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Others_no_arg != nil {
					return d.Olc.Action.Others_no_arg
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_NOVICT_OTHERS
			return
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter text shown to the Character when his victim isnt found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Not_found != nil {
					return d.Olc.Action.Not_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_NOT_FOUND
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("Enter social shown to the Character when it is its own victim.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Char_auto != nil {
					return d.Olc.Action.Char_auto
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_SELF_CHAR
			return
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Enter social shown to Others when the Char is its own victim.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Others_auto != nil {
					return d.Olc.Action.Others_auto
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_SELF_OTHERS
			return
		case 'f':
			fallthrough
		case 'F':
			write_to_output(d, libc.CString("Enter normal social shown to the Character when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Char_found != nil {
					return d.Olc.Action.Char_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_CHAR_FOUND
			return
		case 'g':
			fallthrough
		case 'G':
			write_to_output(d, libc.CString("Enter normal social shown to Others when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Others_found != nil {
					return d.Olc.Action.Others_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_OTHERS_FOUND
			return
		case 'h':
			fallthrough
		case 'H':
			write_to_output(d, libc.CString("Enter normal social shown to the Victim when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Vict_found != nil {
					return d.Olc.Action.Vict_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_VICT_FOUND
			return
		case 'i':
			fallthrough
		case 'I':
			write_to_output(d, libc.CString("Enter 'body part' social shown to the Character when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Char_body_found != nil {
					return d.Olc.Action.Char_body_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_CHAR_BODY_FOUND
			return
		case 'j':
			fallthrough
		case 'J':
			write_to_output(d, libc.CString("Enter 'body part' social shown to Others when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Others_body_found != nil {
					return d.Olc.Action.Others_body_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_OTHERS_BODY_FOUND
			return
		case 'k':
			fallthrough
		case 'K':
			write_to_output(d, libc.CString("Enter 'body part' social shown to the Victim when the victim is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Vict_body_found != nil {
					return d.Olc.Action.Vict_body_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_VICT_VICT_BODY_FOUND
			return
		case 'l':
			fallthrough
		case 'L':
			write_to_output(d, libc.CString("Enter 'object' social shown to the Character when the object is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Char_obj_found != nil {
					return d.Olc.Action.Char_obj_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_OBJ_CHAR_FOUND
			return
		case 'm':
			fallthrough
		case 'M':
			write_to_output(d, libc.CString("Enter 'object' social shown to the Room when the object is found.\r\n[OLD]: %s\r\n[NEW]: "), func() *byte {
				if d.Olc.Action.Others_obj_found != nil {
					return d.Olc.Action.Others_obj_found
				}
				return libc.CString("NULL")
			}())
			d.Olc.Mode = AEDIT_OBJ_OTHERS_FOUND
			return
		default:
			aedit_disp_menu(d)
		}
		return
	case AEDIT_ACTION_NAME:
		if *arg == 0 || C.strchr(arg, ' ') != nil {
			aedit_disp_menu(d)
			return
		}
		if d.Olc.Action.Command != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Command))
		}
		d.Olc.Action.Command = C.strdup(arg)
	case AEDIT_SORT_AS:
		if *arg == 0 || C.strchr(arg, ' ') != nil {
			aedit_disp_menu(d)
			return
		}
		if d.Olc.Action.Sort_as != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Sort_as))
			d.Olc.Action.Sort_as = C.strdup(arg)
		}
	case AEDIT_MIN_CHAR_POS:
		fallthrough
	case AEDIT_MIN_VICT_POS:
		if *arg == 0 {
			aedit_disp_menu(d)
			return
		}
		i = libc.Atoi(libc.GoString(arg))
		if i < POS_DEAD && i > POS_STANDING {
			aedit_disp_menu(d)
			return
		}
		if d.Olc.Mode == AEDIT_MIN_CHAR_POS {
			d.Olc.Action.Min_char_position = i
		} else {
			d.Olc.Action.Min_victim_position = i
		}
	case AEDIT_MIN_CHAR_LEVEL:
		if *arg == 0 {
			aedit_disp_menu(d)
			return
		}
		i = libc.Atoi(libc.GoString(arg))
		if i < 0 {
			aedit_disp_menu(d)
			return
		}
		d.Olc.Action.Min_level_char = i
	case AEDIT_NOVICT_CHAR:
		if d.Olc.Action.Char_no_arg != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Char_no_arg))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Char_no_arg = C.strdup(arg)
		} else {
			d.Olc.Action.Char_no_arg = nil
		}
	case AEDIT_NOVICT_OTHERS:
		if d.Olc.Action.Others_no_arg != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Others_no_arg))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Others_no_arg = C.strdup(arg)
		} else {
			d.Olc.Action.Others_no_arg = nil
		}
	case AEDIT_VICT_CHAR_FOUND:
		if d.Olc.Action.Char_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Char_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Char_found = C.strdup(arg)
		} else {
			d.Olc.Action.Char_found = nil
		}
	case AEDIT_VICT_OTHERS_FOUND:
		if d.Olc.Action.Others_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Others_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Others_found = C.strdup(arg)
		} else {
			d.Olc.Action.Others_found = nil
		}
	case AEDIT_VICT_VICT_FOUND:
		if d.Olc.Action.Vict_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Vict_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Vict_found = C.strdup(arg)
		} else {
			d.Olc.Action.Vict_found = nil
		}
	case AEDIT_VICT_NOT_FOUND:
		if d.Olc.Action.Not_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Not_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Not_found = C.strdup(arg)
		} else {
			d.Olc.Action.Not_found = nil
		}
	case AEDIT_SELF_CHAR:
		if d.Olc.Action.Char_auto != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Char_auto))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Char_auto = C.strdup(arg)
		} else {
			d.Olc.Action.Char_auto = nil
		}
	case AEDIT_SELF_OTHERS:
		if d.Olc.Action.Others_auto != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Others_auto))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Others_auto = C.strdup(arg)
		} else {
			d.Olc.Action.Others_auto = nil
		}
	case AEDIT_VICT_CHAR_BODY_FOUND:
		if d.Olc.Action.Char_body_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Char_body_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Char_body_found = C.strdup(arg)
		} else {
			d.Olc.Action.Char_body_found = nil
		}
	case AEDIT_VICT_OTHERS_BODY_FOUND:
		if d.Olc.Action.Others_body_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Others_body_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Others_body_found = C.strdup(arg)
		} else {
			d.Olc.Action.Others_body_found = nil
		}
	case AEDIT_VICT_VICT_BODY_FOUND:
		if d.Olc.Action.Vict_body_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Vict_body_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Vict_body_found = C.strdup(arg)
		} else {
			d.Olc.Action.Vict_body_found = nil
		}
	case AEDIT_OBJ_CHAR_FOUND:
		if d.Olc.Action.Char_obj_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Char_obj_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Char_obj_found = C.strdup(arg)
		} else {
			d.Olc.Action.Char_obj_found = nil
		}
	case AEDIT_OBJ_OTHERS_FOUND:
		if d.Olc.Action.Others_obj_found != nil {
			libc.Free(unsafe.Pointer(d.Olc.Action.Others_obj_found))
		}
		if *arg != 0 {
			delete_doubledollar(arg)
			d.Olc.Action.Others_obj_found = C.strdup(arg)
		} else {
			d.Olc.Action.Others_obj_found = nil
		}
	default:
	}
	d.Olc.Value = 1
	aedit_disp_menu(d)
}
func do_astat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i    int
		real int = FALSE
		arg  [2048]byte
	)
	if IS_NPC(ch) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Astat which social?\r\n"))
		return
	}
	for i = 0; i <= top_of_socialt; i++ {
		if is_abbrev(&arg[0], (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Command) != 0 {
			real = TRUE
			break
		}
	}
	if real == 0 {
		send_to_char(ch, libc.CString("No such social.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("n) Command         : @y%-15.15s@n 1) Sort as Command : @y%-15.15s@n\r\n2) Min Position[CH]: @c%-8.8s@n        3) Min Position[VT]: @c%-8.8s@n\r\n4) Min Level   [CH]: @c%-3d@n             5) Show if Invis   : @c%s@n\r\na) Char    [NO ARG]: @c%s@n\r\nb) Others  [NO ARG]: @c%s@n\r\nc) Char [NOT FOUND]: @c%s@n\r\nd) Char  [ARG SELF]: @c%s@n\r\ne) Others[ARG SELF]: @c%s@n\r\nf) Char      [VICT]: @c%s@n\r\ng) Others    [VICT]: @c%s@n\r\nh) Victim    [VICT]: @c%s@n\r\ni) Char  [BODY PRT]: @c%s@n\r\nj) Others[BODY PRT]: @c%s@n\r\nk) Victim[BODY PRT]: @c%s@n\r\nl) Char       [OBJ]: @c%s@n\r\nm) Others     [OBJ]: @c%s@n\r\n"), (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Command, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Sort_as, position_types[(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_char_position], position_types[(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_victim_position], (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Min_level_char, func() string {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Hide != 0 {
			return "HIDDEN"
		}
		return "NOT HIDDEN"
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_no_arg != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_no_arg
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_no_arg != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_no_arg
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Not_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Not_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_auto != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_auto
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_auto != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_auto
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_body_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_body_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_body_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_body_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_body_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Vict_body_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_obj_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Char_obj_found
		}
		return libc.CString("")
	}(), func() *byte {
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_obj_found != nil {
			return (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Others_obj_found
		}
		return libc.CString("")
	}())
}
func aedit_find_command(txt *byte) int {
	var cmd int
	for cmd = 1; *(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command != '\n'; cmd++ {
		if C.strncmp((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Sort_as, txt, uint64(C.strlen(txt))) == 0 || C.strcmp((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command, txt) == 0 {
			return cmd
		}
	}
	return -1
}
