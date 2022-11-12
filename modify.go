package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

var string_fields [7]*byte = [7]*byte{libc.CString("name"), libc.CString("short"), libc.CString("long"), libc.CString("description"), libc.CString("title"), libc.CString("delete-description"), libc.CString("\n")}
var length [5]int = [5]int{15, 60, 256, 240, 60}

func smash_tilde(str *byte) {
	var p *byte = str
	for ; *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if *p == '~' && (*((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\r' || *((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\n' || *((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\x00') {
			*p = ' '
		}
	}
	for (func() *byte {
		str = C.strchr(str, '~')
		return str
	}()) != nil {
		*str = ' '
	}
}
func smash_numb(str *byte) {
	var p *byte = str
	for ; *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		if *p == '#' && (*((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\r' || *((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\n' || *((*byte)(unsafe.Add(unsafe.Pointer(p), 1))) == '\x00') {
			*p = ' '
		}
	}
	for (func() *byte {
		str = C.strchr(str, '#')
		return str
	}()) != nil {
		*str = ' '
	}
}
func string_write(d *descriptor_data, writeto **byte, len_ uint64, mailto int, data unsafe.Pointer) {
	if d.Character != nil && !IS_NPC(d.Character) {
		d.Character.Act[int(PLR_WRITING/32)] |= bitvector_t(1 << (int(PLR_WRITING % 32)))
	}
	if CONFIG_IMPROVED_EDITOR != 0 {
		d.Backstr = (*byte)(data)
	} else if data != nil {
		data = nil
	}
	d.Str = writeto
	d.Max_str = len_
	d.Mail_to = int32(mailto)
}
func string_add(d *descriptor_data, str *byte) {
	var action int
	delete_doubledollar(str)
	smash_tilde(str)
	smash_numb(str)
	if (func() int {
		action = int(libc.BoolToInt(*str == '@' && *(*byte)(unsafe.Add(unsafe.Pointer(str), 1)) == 0))
		return action
	}()) != 0 {
		*str = '\x00'
	} else if (func() int {
		action = improved_editor_execute(d, str)
		return action
	}()) == STRINGADD_ACTION {
		return
	}
	if action != STRINGADD_OK {
	} else if (*d.Str) == nil {
		if uint64(C.strlen(str)+3) > d.Max_str {
			send_to_char(d.Character, libc.CString("String too long - Truncated.\r\n"))
			C.strcpy((*byte)(unsafe.Add(unsafe.Pointer(str), d.Max_str-3)), libc.CString("\r\n"))
			*d.Str = (*byte)(unsafe.Pointer(&make([]int8, int(d.Max_str))[0]))
			C.strcpy(*d.Str, str)
			if CONFIG_IMPROVED_EDITOR == 0 {
				action = STRINGADD_SAVE
			}
		} else {
			*d.Str = (*byte)(unsafe.Pointer(&make([]int8, int(C.strlen(str)+3))[0]))
			C.strcpy(*d.Str, str)
		}
	} else {
		if uint64(C.strlen(str)+C.strlen(*d.Str)+3) > d.Max_str {
			send_to_char(d.Character, libc.CString("String too long.  Last line skipped.\r\n"))
			if CONFIG_IMPROVED_EDITOR == 0 {
				action = STRINGADD_SAVE
			} else if action == STRINGADD_OK {
				action = STRINGADD_ACTION
			}
		} else {
			*d.Str = (*byte)(libc.Realloc(unsafe.Pointer(*d.Str), int(C.strlen(*d.Str)*int64(unsafe.Sizeof(int8(0)))+C.strlen(str)+3)))
			C.strcat(*d.Str, str)
		}
	}
	switch action {
	case STRINGADD_ABORT:
		switch d.Connected {
		case CON_CEDIT:
			fallthrough
		case CON_TEDIT:
			fallthrough
		case CON_NEWSEDIT:
			fallthrough
		case CON_REDIT:
			fallthrough
		case CON_MEDIT:
			fallthrough
		case CON_OEDIT:
			fallthrough
		case CON_IEDIT:
			fallthrough
		case CON_EXDESC:
			fallthrough
		case CON_TRIGEDIT:
			fallthrough
		case CON_HEDIT:
			libc.Free(unsafe.Pointer(*d.Str))
			*d.Str = d.Backstr
			d.Backstr = nil
			d.Str = nil
		case CON_PLAYING:
		default:
			basic_mud_log(libc.CString("SYSERR: string_add: Aborting write from unknown origin."))
		}
	case STRINGADD_SAVE:
		if d.Str != nil && *d.Str != nil && **d.Str == '\x00' {
			libc.Free(unsafe.Pointer(*d.Str))
			*d.Str = C.strdup(libc.CString("Nothing.\r\n"))
		}
		if d.Backstr != nil {
			libc.Free(unsafe.Pointer(d.Backstr))
		}
		d.Backstr = nil
	case STRINGADD_ACTION:
	}
	if action == STRINGADD_SAVE || action == STRINGADD_ABORT {
		var (
			i             int
			cleanup_modes [12]struct {
				Mode int
				Func func(dsc *descriptor_data, todo int)
			} = [12]struct {
				Mode int
				Func func(dsc *descriptor_data, todo int)
			}{{Mode: CON_CEDIT, Func: cedit_string_cleanup}, {Mode: CON_MEDIT, Func: medit_string_cleanup}, {Mode: CON_OEDIT, Func: oedit_string_cleanup}, {Mode: CON_REDIT, Func: redit_string_cleanup}, {Mode: CON_TEDIT, Func: tedit_string_cleanup}, {Mode: CON_TRIGEDIT, Func: trigedit_string_cleanup}, {Mode: CON_EXDESC, Func: exdesc_string_cleanup}, {Mode: CON_PLAYING, Func: playing_string_cleanup}, {Mode: CON_IEDIT, Func: oedit_string_cleanup}, {Mode: CON_HEDIT, Func: hedit_string_cleanup}, {Mode: CON_NEWSEDIT, Func: news_string_cleanup}, {Mode: -1, Func: nil}}
		)
		for i = 0; cleanup_modes[i].Func != nil; i++ {
			if d.Connected == cleanup_modes[i].Mode {
				(cleanup_modes[i].Func)(d, action)
			}
		}
		d.Str = nil
		d.Mail_to = 0
		d.Max_str = 0
		if d.Character != nil && !IS_NPC(d.Character) {
			d.Character.Act[int(PLR_MAILING/32)] &= bitvector_t(^(1 << (int(PLR_MAILING % 32))))
			d.Character.Act[int(PLR_WRITING/32)] &= bitvector_t(^(1 << (int(PLR_WRITING % 32))))
		}
	} else if action != STRINGADD_ACTION && uint64(C.strlen(*d.Str)+3) <= d.Max_str {
		C.strcat(*d.Str, libc.CString("\r\n"))
	}
}
func playing_string_cleanup(d *descriptor_data, action int) {
	var (
		board *board_info
		fore  *board_msg
		cur   *board_msg
		aft   *board_msg
	)
	if PLR_FLAGGED(d.Character, PLR_MAILING) {
		if action == STRINGADD_SAVE && *d.Str != nil {
			store_mail(int(d.Mail_to), int(d.Character.Idnum), *d.Str)
			write_to_output(d, libc.CString("Message sent!\r\n"))
			notify_if_playing(d.Character, int(d.Mail_to))
		} else {
			write_to_output(d, libc.CString("Mail aborted.\r\n"))
			libc.Free(unsafe.Pointer(*d.Str))
			libc.Free(unsafe.Pointer(d.Str))
		}
	}
	if PLR_FLAGGED(d.Character, PLR_WRITING) {
		if d.Mail_to >= BOARD_MAGIC {
			if action == STRINGADD_ABORT {
				board = locate_board(obj_vnum(d.Mail_to - BOARD_MAGIC))
				fore = func() *board_msg {
					cur = func() *board_msg {
						aft = nil
						return aft
					}()
					return cur
				}()
				for cur = board.Messages; cur != nil; cur = aft {
					aft = cur.Next
					if cur.Data == *d.Str {
						if board.Messages == cur {
							if cur.Next != nil {
								board.Messages = cur.Next
							} else {
								board.Messages = nil
							}
						}
						if fore != nil {
							fore.Next = aft
						}
						if aft != nil {
							aft.Prev = fore
						}
						libc.Free(unsafe.Pointer(cur.Subject))
						libc.Free(unsafe.Pointer(cur.Data))
						libc.Free(unsafe.Pointer(cur))
						board.Num_messages--
						write_to_output(d, libc.CString("Post aborted.\r\n"))
						return
					}
					fore = cur
				}
				write_to_output(d, libc.CString("Unable to find your message to delete it!\r\n"))
			} else {
				write_to_output(d, libc.CString("\r\nPost saved.\r\n"))
				save_board(locate_board(obj_vnum(d.Mail_to - BOARD_MAGIC)))
			}
		}
	}
}
func exdesc_string_cleanup(d *descriptor_data, action int) {
	if action == STRINGADD_ABORT {
		write_to_output(d, libc.CString("Description aborted.\r\n"))
	}
	write_to_output(d, config_info.Operation.MENU)
	d.Connected = CON_MENU
}
func do_skillset(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict  *char_data
		name  [2048]byte
		buf   [2048]byte
		help  [64936]byte
		skill int
		value int
		i     int = 0
		qend  int
	)
	argument = one_argument(argument, &name[0])
	if name[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: skillset <name> '<skill>' <value>\r\nSkill being one of the following:\r\n"))
		for func() int {
			qend = 0
			return func() int {
				i = 0
				return i
			}()
		}(); i < SKILL_TABLE_SIZE; i++ {
			if spell_info[i].Name == unused_spellname {
				continue
			}
			send_to_char(ch, libc.CString("%18s"), spell_info[i].Name)
			if func() int {
				p := &qend
				x := *p
				*p++
				return x
			}()%4 == 3 {
				send_to_char(ch, libc.CString("\r\n"))
			}
		}
		if qend%4 != 0 {
			send_to_char(ch, libc.CString("\r\n"))
		}
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &name[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		return
	}
	skip_spaces(&argument)
	if *argument == 0 {
		i = stdio.Snprintf(&help[0], int(64936-uintptr(i)), "\r\nSkills:\r\n")
		i += print_skills_by_type(vict, &help[i], int(64936-uintptr(i)), 1<<1, nil)
		i += stdio.Snprintf(&help[i], int(64936-uintptr(i)), "\r\nSpells:\r\n")
		i += print_skills_by_type(vict, &help[i], int(64936-uintptr(i)), 1<<0, nil)
		if config_info.Play.Enable_languages != 0 {
			i += stdio.Snprintf(&help[i], int(64936-uintptr(i)), "\r\nLanguages:\r\n")
			i += print_skills_by_type(vict, &help[i], int(64936-uintptr(i)), (1<<1)|1<<2, nil)
		}
		if i >= int(64936) {
			C.strcpy((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&help[64936]), -C.strlen(libc.CString("** OVERFLOW **"))))), -1)), libc.CString("** OVERFLOW **"))
		}
		page_string(ch.Desc, &help[0], TRUE)
		return
	}
	if *argument != '\'' {
		send_to_char(ch, libc.CString("Skill must be enclosed in: ''\r\n"))
		return
	}
	for qend = 1; *(*byte)(unsafe.Add(unsafe.Pointer(argument), qend)) != 0 && *(*byte)(unsafe.Add(unsafe.Pointer(argument), qend)) != '\''; qend++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(argument), qend)) = byte(int8(C.tolower(int(*(*byte)(unsafe.Add(unsafe.Pointer(argument), qend))))))
	}
	if *(*byte)(unsafe.Add(unsafe.Pointer(argument), qend)) != '\'' {
		send_to_char(ch, libc.CString("Skill must be enclosed in: ''\r\n"))
		return
	}
	C.strcpy(&help[0], (*byte)(unsafe.Add(unsafe.Pointer(argument), 1)))
	help[qend-1] = '\x00'
	if (func() int {
		skill = find_skill_num(&help[0], 1<<1)
		return skill
	}()) <= 0 {
		send_to_char(ch, libc.CString("Unrecognized skill.\r\n"))
		return
	}
	argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), qend+1))
	argument = one_argument(argument, &buf[0])
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("Learned value expected.\r\n"))
		return
	}
	value = libc.Atoi(libc.GoString(&buf[0]))
	if value < 0 {
		send_to_char(ch, libc.CString("Minimum value for learned is 0.\r\n"))
		return
	}
	for {
		vict.Skills[skill] = int8(value)
		if true {
			break
		}
	}
	mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("skillset: %s changed %s's '%s' to %d."), GET_NAME(ch), GET_NAME(vict), spell_info[skill].Name, value)
	send_to_char(ch, libc.CString("You change %s's %s to %d.\r\n"), GET_NAME(vict), spell_info[skill].Name, value)
}
func next_page(str *byte, ch *char_data) *byte {
	var (
		col  int = 1
		line int = 1
	)
	for ; ; str = (*byte)(unsafe.Add(unsafe.Pointer(str), 1)) {
		if *str == '\x00' {
			return nil
		} else if line > (int(ch.Player_specials.Page_length) - (func() int {
			if PRF_FLAGGED(ch, PRF_COMPACT) {
				return 1
			}
			return 2
		}())) {
			return str
		} else if *str == '\x1b' {
			str = (*byte)(unsafe.Add(unsafe.Pointer(str), 1))
		} else if *str == '@' {
			if *((*byte)(unsafe.Add(unsafe.Pointer(str), 1))) != '@' {
				str = (*byte)(unsafe.Add(unsafe.Pointer(str), 1))
			}
		} else {
			if *str == '\r' {
				col = 1
			} else if *str == '\n' {
				line++
			} else if func() int {
				p := &col
				x := *p
				*p++
				return x
			}() > PAGE_WIDTH {
				col = 1
				line++
			}
		}
	}
}
func count_pages(str *byte, ch *char_data) int {
	var pages int
	for pages = 1; (func() *byte {
		str = next_page(str, ch)
		return str
	}()) != nil; pages++ {
	}
	return pages
}
func paginate_string(str *byte, d *descriptor_data) {
	var i int
	if d.Showstr_count != 0 {
		*d.Showstr_vector = str
	}
	for i = 1; i < d.Showstr_count && str != nil; i++ {
		str = func() *byte {
			p := (**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(i)))
			*(**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(i))) = next_page(str, d.Character)
			return *p
		}()
	}
	d.Showstr_page = 0
}
func page_string(d *descriptor_data, str *byte, keep_internal int) {
	var actbuf [2048]byte = func() [2048]byte {
		var t [2048]byte
		copy(t[:], []byte(""))
		return t
	}()
	if d == nil {
		return
	}
	if str == nil || *str == 0 {
		return
	}
	if d.Character.Player_specials.Page_length < 5 || d.Character.Player_specials.Page_length > 50 {
		d.Character.Player_specials.Page_length = PAGE_LENGTH
	}
	d.Showstr_count = count_pages(str, d.Character)
	d.Showstr_vector = &make([]*byte, d.Showstr_count)[0]
	if keep_internal != 0 {
		d.Showstr_head = C.strdup(str)
		paginate_string(d.Showstr_head, d)
	} else {
		paginate_string(str, d)
	}
	show_string(d, &actbuf[0])
}
func show_string(d *descriptor_data, input *byte) {
	var (
		buffer [64936]byte
		buf    [2048]byte
		diff   int
	)
	any_one_arg(input, &buf[0])
	if C.tolower(int(buf[0])) == 'q' {
		libc.Free(unsafe.Pointer(d.Showstr_vector))
		d.Showstr_vector = nil
		d.Showstr_count = 0
		if d.Showstr_head != nil {
			libc.Free(unsafe.Pointer(d.Showstr_head))
			d.Showstr_head = nil
		}
		return
	} else if C.tolower(int(buf[0])) == 'r' {
		d.Showstr_page = MAX(0, d.Showstr_page-1)
	} else if C.tolower(int(buf[0])) == 'b' {
		d.Showstr_page = MAX(0, d.Showstr_page-2)
	} else if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(buf[0]))))) & int(uint16(int16(_ISdigit)))) != 0 {
		d.Showstr_page = MAX(0, MIN(libc.Atoi(libc.GoString(&buf[0]))-1, d.Showstr_count-1))
	} else if buf[0] != 0 {
		send_to_char(d.Character, libc.CString("Valid commands while paging are RETURN, Q, R, B, or a numeric value.\r\n"))
		return
	}
	if d.Showstr_page+1 >= d.Showstr_count {
		send_to_char(d.Character, libc.CString("%s"), *(**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(d.Showstr_page))))
		libc.Free(unsafe.Pointer(d.Showstr_vector))
		d.Showstr_vector = nil
		d.Showstr_count = 0
		if d.Showstr_head != nil {
			libc.Free(unsafe.Pointer(d.Showstr_head))
			d.Showstr_head = nil
		}
	} else {
		diff = int(int64(uintptr(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(d.Showstr_page+1))))) - uintptr(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(d.Showstr_page)))))))
		if diff > int(MAX_STRING_LENGTH-3) {
			diff = int(MAX_STRING_LENGTH - 3)
		}
		C.strncpy(&buffer[0], *(**byte)(unsafe.Add(unsafe.Pointer(d.Showstr_vector), unsafe.Sizeof((*byte)(nil))*uintptr(d.Showstr_page))), uint64(diff))
		if buffer[diff-2] == '\r' && buffer[diff-1] == '\n' {
			buffer[diff] = '\x00'
		} else if buffer[diff-2] == '\n' && buffer[diff-1] == '\r' {
			C.strcpy((*byte)(unsafe.Add(unsafe.Pointer(&buffer[diff]), -2)), libc.CString("\r\n"))
		} else if buffer[diff-1] == '\r' || buffer[diff-1] == '\n' {
			C.strcpy((*byte)(unsafe.Add(unsafe.Pointer(&buffer[diff]), -1)), libc.CString("\r\n"))
		} else {
			C.strcpy(&buffer[diff], libc.CString("\r\n"))
		}
		send_to_char(d.Character, libc.CString("%s"), &buffer[0])
		d.Showstr_page++
	}
}
