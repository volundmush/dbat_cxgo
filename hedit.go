package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func do_oasis_hedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		d   *descriptor_data
	)
	if IS_NPC(ch) || ch.Desc == nil || ch.Desc.Connected != CON_PLAYING {
		return
	}
	if HEDITS == TRUE {
		send_to_char(ch, libc.CString("Sorry, only one person can edit help files at a time.\r\n"))
		return
	}
	if ch.Admlevel < 4 && (libc.StrCaseCmp(libc.CString("Tepsih"), GET_NAME(ch)) == 0 && libc.StrCaseCmp(libc.CString("Rogoshen"), GET_NAME(ch)) == 0) {
		send_to_char(ch, libc.CString("Sorry you are incapable of editing help files at this time.\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Please specify a help entry to edit.\r\n"))
		return
	}
	d = ch.Desc
	if libc.StrCaseCmp(libc.CString("save"), &arg[0]) == 0 {
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves help files."), GET_NAME(ch))
		hedit_save_to_disk(d)
		send_to_char(ch, libc.CString("Saving help files.\r\n"))
		return
	}
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	d.Olc.Number = 0
	d.Olc.Storage = libc.StrDup(&arg[0])
	d.Olc.Zone_num = zone_rnum(search_help(d.Olc.Storage, ADMLVL_IMPL))
	if d.Olc.Zone_num == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Do you wish to add the '%s' help file? "), d.Olc.Storage)
		d.Olc.Mode = HEDIT_CONFIRM_ADD
	} else {
		send_to_char(ch, libc.CString("Do you wish to edit the '%s' help file? "), (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(d.Olc.Zone_num)))).Keywords)
		d.Olc.Mode = HEDIT_CONFIRM_EDIT
	}
	d.Connected = CON_HEDIT
	act(libc.CString("$n starts using OLC."), TRUE, ch, nil, nil, TO_ROOM)
	HEDITS = TRUE
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing help files."), GET_NAME(ch))
}
func hedit_setup_new(d *descriptor_data) {
	d.Olc.Help = new(help_index_element)
	var buf3 [2048]byte
	stdio.Sprintf(&buf3[0], "<<X<< Put helpfile keywords here in caps")
	d.Olc.Help.Keywords = libc.StrDup(&buf3[0])
	d.Olc.Help.Entry = libc.CString("\r\nThis help file is unfinished.\r\n")
	d.Olc.Help.Min_level = 0
	d.Olc.Help.Duplicate = 0
	d.Olc.Value = 0
	hedit_disp_menu(d)
}
func hedit_setup_existing(d *descriptor_data, rnum int) {
	d.Olc.Help = new(help_index_element)
	d.Olc.Help.Keywords = str_udup((*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(rnum)))).Keywords)
	d.Olc.Help.Entry = str_udup((*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(rnum)))).Entry)
	d.Olc.Help.Duplicate = (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(rnum)))).Duplicate
	d.Olc.Help.Min_level = (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(rnum)))).Min_level
	d.Olc.Value = 0
	hedit_disp_menu(d)
}
func hedit_save_internally(d *descriptor_data) {
	var new_help_table *help_index_element = nil
	if d.Olc.Zone_num == zone_rnum(-1) {
		var i int
		new_help_table = &make([]help_index_element, top_of_helpt+2)[0]
		for i = 0; i < top_of_helpt; i++ {
			*(*help_index_element)(unsafe.Add(unsafe.Pointer(new_help_table), unsafe.Sizeof(help_index_element{})*uintptr(i))) = *(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))
		}
		*(*help_index_element)(unsafe.Add(unsafe.Pointer(new_help_table), unsafe.Sizeof(help_index_element{})*uintptr(func() int {
			p := &top_of_helpt
			x := *p
			*p++
			return x
		}()))) = *d.Olc.Help
		libc.Free(unsafe.Pointer(help_table))
		help_table = new_help_table
	} else {
		*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(d.Olc.Zone_num))) = *d.Olc.Help
	}
	add_to_save_list(HEDIT_PERMISSION, int(SL_GLD+2))
	hedit_save_to_disk(d)
}
func hedit_save_to_disk(d *descriptor_data) {
	var (
		fp         *stdio.File
		buf1       [64936]byte
		index_name [256]byte
		i          int
	)
	stdio.Snprintf(&index_name[0], int(256), "%s%s", LIB_TEXT, HELP_FILE)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&index_name[0]), "w")
		return fp
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Could not write help index file"))
		return
	}
	for i = 0; i < top_of_helpt; i++ {
		if (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Duplicate != 0 {
			continue
		}
		libc.StrNCpy(&buf1[0], func() *byte {
			if (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Entry != nil {
				return (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Entry
			}
			return libc.CString("Empty\r\n")
		}(), int(64936-1))
		strip_cr(&buf1[0])
		stdio.Fprintf(fp, "%s#%d\n", &buf1[0], (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Min_level)
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	remove_from_save_list(HEDIT_PERMISSION, int(SL_GLD+2))
	free_help_table()
	index_boot(DB_BOOT_HLP)
}
func hedit_disp_menu(d *descriptor_data) {
	write_to_output(d, libc.CString("@n-- Help file editor\r\n@g1@n) Keywords    : @y%s\n@g2@n) Entry       :\r\n@y%s@g3@n) Min Level   : @y%d\r\n@gQ@n) Quit\r\nEnter choice : "), d.Olc.Help.Keywords, d.Olc.Help.Entry, d.Olc.Help.Min_level)
	d.Olc.Mode = HEDIT_MAIN_MENU
}
func hedit_parse(d *descriptor_data, arg *byte) {
	var (
		oldtext *byte = (*byte)(unsafe.Pointer(uintptr('\x00')))
		number  int
		change  int = TRUE
	)
	switch d.Olc.Mode {
	case HEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			if d.Olc.Help.Keywords == nil {
				hedit_disp_menu(d)
				write_to_output(d, libc.CString("\n@RYou must fill in the keywords before you save.@n\r\n"))
			} else if libc.StrStr(d.Olc.Help.Keywords, libc.CString("undefined")) != nil {
				hedit_disp_menu(d)
				write_to_output(d, libc.CString("\n@RYou must fill in the keywords before you save.@n\r\n"))
			} else if libc.StrStr(d.Olc.Help.Keywords, libc.CString("<<X<<")) != nil {
				hedit_disp_menu(d)
				write_to_output(d, libc.CString("\n@RYou must fill in the keywords before you save.@n\r\n"))
			} else {
				write_to_output(d, libc.CString("Help saved to disk.\r\n"))
				var buf [2048]byte
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "%s", d.Olc.Help.Keywords)
				send_to_imm(libc.CString("@gHedit@D: @w%s@G has just edited and saved, @Y%s@G.@n"), d.Character.Name, &buf[0])
				hedit_save_internally(d)
				cleanup_olc(d, CLEANUP_STRUCTS)
				HEDITS = FALSE
			}
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
			HEDITS = FALSE
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to save your changes? : \r\n"))
		}
		return
	case HEDIT_CONFIRM_EDIT:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			hedit_setup_existing(d, int(d.Olc.Zone_num))
		case 'q':
			fallthrough
		case 'Q':
			cleanup_olc(d, CLEANUP_ALL)
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Zone_num++
			for ; d.Olc.Zone_num < zone_rnum(top_of_helpt); d.Olc.Zone_num++ {
				if is_abbrev(d.Olc.Storage, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(d.Olc.Zone_num)))).Keywords) != 0 {
					break
				} else {
					d.Olc.Zone_num = zone_rnum(top_of_helpt + 1)
				}
			}
			if d.Olc.Zone_num > zone_rnum(top_of_helpt) {
				write_to_output(d, libc.CString("Do you wish to add the '%s' help file? "), d.Olc.Storage)
				d.Olc.Mode = HEDIT_CONFIRM_ADD
			} else {
				write_to_output(d, libc.CString("Do you wish to edit the '%s' help file? "), (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(d.Olc.Zone_num)))).Keywords)
				d.Olc.Mode = HEDIT_CONFIRM_EDIT
			}
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to edit the '%s' help file? "), (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(d.Olc.Zone_num)))).Keywords)
		}
		return
	case HEDIT_CONFIRM_ADD:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			hedit_setup_new(d)
		case 'n':
			fallthrough
		case 'N':
			fallthrough
		case 'q':
			fallthrough
		case 'Q':
			cleanup_olc(d, CLEANUP_ALL)
			HEDITS = FALSE
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to add the '%s' help file? "), d.Olc.Storage)
		}
		return
	case HEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = HEDIT_CONFIRM_SAVESTRING
			} else {
				write_to_output(d, libc.CString("No changes made.\r\n"))
				cleanup_olc(d, CLEANUP_ALL)
				HEDITS = FALSE
			}
		case '1':
			d.Olc.Mode = HEDIT_KEYWORDS
			clear_screen(d)
			write_to_output(d, libc.CString("Enter help file keywords: "))
		case '2':
			d.Olc.Mode = HEDIT_ENTRY
			clear_screen(d)
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter help entry: (/s saves /h for help)\r\n"))
			if d.Olc.Help.Entry != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Help.Entry)
				oldtext = libc.StrDup(d.Olc.Help.Entry)
			}
			string_write(d, &d.Olc.Help.Entry, MAX_MESSAGE_LENGTH, 0, unsafe.Pointer(oldtext))
			d.Olc.Value = 1
		case '3':
			write_to_output(d, libc.CString("Enter min level : "))
			d.Olc.Mode = HEDIT_MIN_LEVEL
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\n"))
			hedit_disp_menu(d)
		}
		return
	case HEDIT_KEYWORDS:
		if d.Olc.Help.Keywords != nil {
			libc.Free(unsafe.Pointer(d.Olc.Help.Keywords))
		}
		if libc.StrLen(arg) > MAX_HELP_KEYWORDS {
			*(*byte)(unsafe.Add(unsafe.Pointer(arg), int(MAX_HELP_KEYWORDS-1))) = '\x00'
		}
		strip_cr(arg)
		d.Olc.Help.Keywords = str_udup(arg)
		var buf4 [8192]byte
		if libc.StrStr(d.Olc.Help.Keywords, libc.CString("undefined")) != nil {
			d.Olc.Mode = HEDIT_KEYWORDS
			clear_screen(d)
			write_to_output(d, libc.CString("@RYou must at least enter SOME keywords.@n\n"))
			write_to_output(d, libc.CString("Keywords: "))
			change = FALSE
		} else if libc.StrStr(d.Olc.Help.Keywords, libc.CString("<<X<<")) != nil {
			d.Olc.Mode = HEDIT_KEYWORDS
			clear_screen(d)
			write_to_output(d, libc.CString("@RLet's not joke around with help files now.@n\n"))
			write_to_output(d, libc.CString("Keywords: "))
			change = FALSE
		} else if libc.StrStr(d.Olc.Help.Keywords, libc.CString("<<x<<")) != nil {
			d.Olc.Mode = HEDIT_KEYWORDS
			clear_screen(d)
			write_to_output(d, libc.CString("@RLet's not joke around with help files now.@n\n"))
			write_to_output(d, libc.CString("Keywords: "))
			change = FALSE
		} else {
			stdio.Sprintf(&buf4[0], "%s\r\n----------\r\n\r\n%s", d.Olc.Help.Keywords, d.Olc.Help.Entry)
			d.Olc.Help.Entry = libc.StrDup(&buf4[0])
		}
	case HEDIT_ENTRY:
		mudlog(TRUE, ADMLVL_BUILDER, BRF, libc.CString("SYSERR: Reached HEDIT_ENTRY case in parse_hedit"))
	case HEDIT_MIN_LEVEL:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > ADMLVL_IMPL {
			write_to_output(d, libc.CString("That is not a valid choice!\r\nEnter min level:-\r\n] "))
		} else {
			d.Olc.Help.Min_level = number
			break
		}
		return
	default:
		mudlog(TRUE, ADMLVL_BUILDER, BRF, libc.CString("SYSERR: Reached default case in parse_hedit"))
	}
	if change == TRUE {
		d.Olc.Value = 1
		hedit_disp_menu(d)
	}
}
func hedit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case HEDIT_ENTRY:
		hedit_disp_menu(d)
	}
}
func do_helpcheck(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		do_action func(ch *char_data, argument *byte, cmd int, subcmd int)
		buf       [64936]byte
		i         int
		count     int    = 0
		len_      uint64 = 0
		nlen      uint64
	)
	send_to_char(ch, libc.CString("Commands without help entries:\r\n"))
	for i = 1; *(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command != '\n'; i++ {
		if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command_pointer) != libc.FuncAddr(do_action) && int((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Minimum_level) >= 0 {
			if search_help((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command, ADMLVL_IMPL) == int(-1) {
				nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%-20.20s%s", (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command, func() string {
					if func() int {
						p := &count
						*p++
						return *p
					}()%3 != 0 {
						return ""
					}
					return "\r\n"
				}()))
				if len_+nlen >= uint64(64936) {
					break
				}
				len_ += nlen
			}
		}
	}
	if count%3 != 0 && len_ < uint64(64936) {
		nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "\r\n"))
	}
	if ch.Desc != nil {
		page_string(ch.Desc, &buf[0], TRUE)
	}
	buf[0] = '\x00'
}
func do_hindex(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		len_  int
		count int = 0
		i     int
		num   int = 0
		buf   [64936]byte
	)
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Usage: hindex <string>\r\n"))
		for i = 0; i < top_of_helpt; i++ {
			num++
		}
		if num > 0 && ch.Admlevel > 0 {
			send_to_char(ch, libc.CString("\r\n@D[@Y%d@y Help files in index.@D]@n\r\n"), num)
		}
		return
	}
	len_ = stdio.Sprintf(&buf[0], "Help index entries based on '%s':\r\n", argument)
	for i = 0; i < top_of_helpt; i++ {
		num++
		if is_abbrev(argument, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Keywords) != 0 && ch.Admlevel >= (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Min_level {
			len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%-20.20s%s", (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(i)))).Keywords, func() string {
				if func() int {
					p := &count
					*p++
					return *p
				}()%3 != 0 {
					return ""
				}
				return "\r\n"
			}())
		}
	}
	if count%3 != 0 {
		len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "\r\n")
	}
	if count == 0 {
		len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "  None.\r\n")
	}
	if count > 0 && ch.Admlevel > 0 {
		len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "  %d Help files in index.\r\n", count)
	}
	page_string(ch.Desc, &buf[0], TRUE)
}
