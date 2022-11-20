package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func send_editor_help(d *descriptor_data) {
	if CONFIG_IMPROVED_EDITOR != 0 {
		write_to_output(d, libc.CString("Instructions: /s to save, /h for more options.\r\n"))
	} else {
		write_to_output(d, libc.CString("Instructions: Type @ on a line by itself to end.\r\n"))
	}
}
func improved_editor_execute(d *descriptor_data, str *byte) int {
	var actions [2048]byte
	if *str != '/' {
		return STRINGADD_OK
	}
	libc.StrNCpy(&actions[0], (*byte)(unsafe.Add(unsafe.Pointer(str), 2)), int(2048-1))
	actions[2048-1] = '\x00'
	*str = '\x00'
	switch *(*byte)(unsafe.Add(unsafe.Pointer(str), 1)) {
	case 'a':
		return STRINGADD_ABORT
	case 'c':
		if *d.Str != nil {
			libc.Free(unsafe.Pointer(*d.Str))
			*d.Str = nil
			write_to_output(d, libc.CString("Current buffer cleared.\r\n"))
		} else {
			write_to_output(d, libc.CString("Current buffer empty.\r\n"))
		}
	case 'd':
		parse_edit_action(PARSE_DELETE, &actions[0], d)
	case 'e':
		parse_edit_action(PARSE_EDIT, &actions[0], d)
	case 'f':
		if *d.Str != nil {
			parse_edit_action(PARSE_FORMAT, &actions[0], d)
		} else {
			write_to_output(d, libc.CString("Current buffer empty.\r\n"))
		}
	case 'i':
		if *d.Str != nil {
			parse_edit_action(PARSE_INSERT, &actions[0], d)
		} else {
			write_to_output(d, libc.CString("Current buffer empty.\r\n"))
		}
	case 'h':
		parse_edit_action(PARSE_HELP, &actions[0], d)
	case 'l':
		if *d.Str != nil {
			parse_edit_action(PARSE_LIST_NORM, &actions[0], d)
		} else {
			write_to_output(d, libc.CString("Current buffer empty.\r\n"))
		}
	case 'n':
		if *d.Str != nil {
			parse_edit_action(PARSE_LIST_NUM, &actions[0], d)
		} else {
			write_to_output(d, libc.CString("Current buffer empty.\r\n"))
		}
	case 'r':
		parse_edit_action(PARSE_REPLACE, &actions[0], d)
	case 's':
		return STRINGADD_SAVE
	default:
		write_to_output(d, libc.CString("Invalid option.\r\n"))
	}
	return STRINGADD_ACTION
}
func parse_edit_action(command int, string_ *byte, d *descriptor_data) {
	var (
		indent    int = 0
		rep_all   int = 0
		flags     int = 0
		replaced  int
		i         int
		line_low  int
		line_high int
		j         int = 0
		total_len uint
		s         *byte
		t         *byte
		temp      int8
		buf       [64936]byte
		buf2      [64936]byte
	)
	switch command {
	case PARSE_HELP:
		write_to_output(d, libc.CString("Editor command formats: /<letter>\r\n\r\n/a         -  aborts editor\r\n/c         -  clears buffer\r\n/d#        -  deletes a line #\r\n/e# <text> -  changes the line at # with <text>\r\n/f         -  formats text\r\n/fi        -  indented formatting of text\r\n/fi#       -  indented formatting on a specific line\r\n/fi #-#    -  indented formating on specific lines\r\n/h         -  list text editor commands\r\n/i# <text> -  inserts <text> before line #\r\n/l         -  lists buffer\r\n/n         -  lists buffer with line numbers\r\n/r 'a' 'b' -  replace 1st occurance of text <a> in buffer with text <b>\r\n/ra 'a' 'b'-  replace all occurances of text <a> within buffer with text <b>\r\n              usage: /r[a] 'pattern' 'replacement'\r\n/s         -  saves text\r\n"))
	case PARSE_FORMAT:
		if d.Connected == CON_TRIGEDIT {
			write_to_output(d, libc.CString("Script %sformatted.\r\n"), func() string {
				if format_script(d) != 0 {
					return ""
				}
				return "not "
			}())
			return
		}
		for libc.IsAlpha(rune(*(*byte)(unsafe.Add(unsafe.Pointer(string_), j)))) && j < 2 {
			if *(*byte)(unsafe.Add(unsafe.Pointer(string_), func() int {
				p := &j
				x := *p
				*p++
				return x
			}())) == 'i' && indent == 0 {
				indent = 1
				flags += 1 << 0
			}
		}
		switch stdio.Sscanf(func() *byte {
			if indent != 0 {
				return (*byte)(unsafe.Add(unsafe.Pointer(string_), 1))
			}
			return string_
		}(), " %d - %d ", &line_low, &line_high) {
		case -1:
			fallthrough
		case 0:
			line_low = 1
			line_high = 0xF423F
		case 1:
			line_high = line_low
		case 2:
			if line_high < line_low {
				write_to_output(d, libc.CString("That range is invalid.\\r\\n"))
				return
			}
		}
		line_low = int(MAX(1, int64(line_low)))
		switch stdio.Sscanf(func() *byte {
			if indent != 0 {
				return (*byte)(unsafe.Add(unsafe.Pointer(string_), 1))
			}
			return string_
		}(), " %d - %d ", &line_low, &line_high) {
		case -1:
			fallthrough
		case 0:
			line_low = 1
			line_high = 0xF423F
		case 1:
			line_high = line_low
		case 2:
			if line_high < line_low {
				write_to_output(d, libc.CString("That range is invalid.\r\n"))
				return
			}
		}
		line_low = int(MAX(1, int64(line_low)))
		if format_text(d.Str, flags, d, uint(d.Max_str), line_low, line_high) {
			write_to_output(d, libc.CString("Text formatted with%s indent.\r\n"), func() string {
				if indent != 0 {
					return ""
				}
				return "out"
			}())
		}
	case PARSE_REPLACE:
		for libc.IsAlpha(rune(*(*byte)(unsafe.Add(unsafe.Pointer(string_), j)))) && j < 2 {
			if *(*byte)(unsafe.Add(unsafe.Pointer(string_), func() int {
				p := &j
				x := *p
				*p++
				return x
			}())) == 'a' && indent == 0 {
				rep_all = 1
			}
		}
		if (func() *byte {
			s = libc.StrTok(string_, libc.CString("'"))
			return s
		}()) == nil {
			write_to_output(d, libc.CString("Invalid format.\r\n"))
			return
		} else if (func() *byte {
			s = libc.StrTok(nil, libc.CString("'"))
			return s
		}()) == nil {
			write_to_output(d, libc.CString("Target string must be enclosed in single quotes.\r\n"))
			return
		} else if (func() *byte {
			t = libc.StrTok(nil, libc.CString("'"))
			return t
		}()) == nil {
			write_to_output(d, libc.CString("No replacement string.\r\n"))
			return
		} else if (func() *byte {
			t = libc.StrTok(nil, libc.CString("'"))
			return t
		}()) == nil {
			write_to_output(d, libc.CString("Replacement string must be enclosed in single quotes.\r\n"))
			return
		} else if *d.Str == nil {
			return
		} else if (func() uint {
			total_len = uint((libc.StrLen(t) - libc.StrLen(s)) + libc.StrLen(*d.Str))
			return total_len
		}()) <= uint(d.Max_str) {
			if (func() int {
				replaced = replace_str(d.Str, s, t, rep_all, uint(d.Max_str))
				return replaced
			}()) > 0 {
				write_to_output(d, libc.CString("Replaced %d occurance%sof '%s' with '%s'.\r\n"), replaced, func() string {
					if replaced != 1 {
						return "s "
					}
					return " "
				}(), s, t)
			} else if replaced == 0 {
				write_to_output(d, libc.CString("String '%s' not found.\r\n"), s)
			} else {
				write_to_output(d, libc.CString("ERROR: Replacement string causes buffer overflow, aborted replace.\r\n"))
			}
		} else {
			write_to_output(d, libc.CString("Not enough space left in buffer.\r\n"))
		}
	case PARSE_DELETE:
		switch stdio.Sscanf(string_, " %d - %d ", &line_low, &line_high) {
		case 0:
			write_to_output(d, libc.CString("You must specify a line number or range to delete.\r\n"))
			return
		case 1:
			line_high = line_low
		case 2:
			if line_high < line_low {
				write_to_output(d, libc.CString("That range is invalid.\r\n"))
				return
			}
		}
		i = 1
		total_len = 1
		if (func() *byte {
			s = *d.Str
			return s
		}()) == nil {
			write_to_output(d, libc.CString("Buffer is empty.\r\n"))
			return
		} else if line_low > 0 {
			for s != nil && i < line_low {
				if (func() *byte {
					s = libc.StrChr(s, '\n')
					return s
				}()) != nil {
					i++
					s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				}
			}
			if s == nil || i < line_low {
				write_to_output(d, libc.CString("Line(s) out of range; not deleting.\r\n"))
				return
			}
			t = s
			for s != nil && i < line_high {
				if (func() *byte {
					s = libc.StrChr(s, '\n')
					return s
				}()) != nil {
					i++
					total_len++
					s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				}
			}
			if s != nil && (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				for *(func() *byte {
					p := &s
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return *p
				}()) != 0 {
					*(func() *byte {
						p := &t
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}()) = *s
				}
			} else {
				total_len--
			}
			*t = '\x00'
			*d.Str = (*byte)(libc.Realloc(unsafe.Pointer(*d.Str), libc.StrLen(*d.Str)*int(unsafe.Sizeof(int8(0)))+3))
			write_to_output(d, libc.CString("%d line%sdeleted.\r\n"), total_len, func() string {
				if total_len != 1 {
					return "s "
				}
				return " "
			}())
		} else {
			write_to_output(d, libc.CString("Invalid, line numbers to delete must be higher than 0.\r\n"))
			return
		}
	case PARSE_LIST_NORM:
		buf[0] = '\x00'
		if *string_ != 0 {
			switch stdio.Sscanf(string_, " %d - %d ", &line_low, &line_high) {
			case 0:
				line_low = 1
				line_high = 0xF423F
			case 1:
				line_high = line_low
			}
		} else {
			line_low = 1
			line_high = 0xF423F
		}
		if line_low < 1 {
			write_to_output(d, libc.CString("Line numbers must be greater than 0.\r\n"))
			return
		} else if line_high < line_low {
			write_to_output(d, libc.CString("That range is invalid.\r\n"))
			return
		}
		buf[0] = '\x00'
		if line_high < 0xF423F || line_low > 1 {
			stdio.Sprintf(&buf[0], "Current buffer range [%d - %d]:\r\n", line_low, line_high)
		}
		i = 1
		total_len = 0
		s = *d.Str
		for s != nil && i < line_low {
			if (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				i++
				s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
			}
		}
		if i < line_low || s == nil {
			write_to_output(d, libc.CString("Line(s) out of range; no buffer listing.\r\n"))
			return
		}
		t = s
		for s != nil && i <= line_high {
			if (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				i++
				total_len++
				s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
			}
		}
		if s != nil {
			temp = int8(*s)
			*s = '\x00'
			libc.StrCat(&buf[0], t)
			*s = byte(temp)
		} else {
			libc.StrCat(&buf[0], t)
		}
		stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\r\n%d line%sshown.\r\n", total_len, func() string {
			if total_len != 1 {
				return "s "
			}
			return " "
		}())
		page_string(d, &buf[0], 1)
	case PARSE_LIST_NUM:
		buf[0] = '\x00'
		if *string_ != 0 {
			switch stdio.Sscanf(string_, " %d - %d ", &line_low, &line_high) {
			case 0:
				line_low = 1
				line_high = 0xF423F
			case 1:
				line_high = line_low
			}
		} else {
			line_low = 1
			line_high = 0xF423F
		}
		if line_low < 1 {
			write_to_output(d, libc.CString("Line numbers must be greater than 0.\r\n"))
			return
		}
		if line_high < line_low {
			write_to_output(d, libc.CString("That range is invalid.\r\n"))
			return
		}
		buf[0] = '\x00'
		i = 1
		total_len = 0
		s = *d.Str
		for s != nil && i < line_low {
			if (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				i++
				s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
			}
		}
		if i < line_low || s == nil {
			write_to_output(d, libc.CString("Line(s) out of range; no buffer listing.\r\n"))
			return
		}
		t = s
		for s != nil && i <= line_high {
			if (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				i++
				total_len++
				s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				temp = int8(*s)
				*s = '\x00'
				stdio.Sprintf(&buf[0], "%s%4d: ", &buf[0], i-1)
				libc.StrCat(&buf[0], t)
				*s = byte(temp)
				t = s
			}
		}
		if s != nil && t != nil {
			temp = int8(*s)
			*s = '\x00'
			libc.StrCat(&buf[0], t)
			*s = byte(temp)
		} else if t != nil {
			libc.StrCat(&buf[0], t)
		}
		page_string(d, &buf[0], 1)
	case PARSE_INSERT:
		half_chop(string_, &buf[0], &buf2[0])
		if buf[0] == '\x00' {
			write_to_output(d, libc.CString("You must specify a line number before which to insert text.\r\n"))
			return
		}
		line_low = libc.Atoi(libc.GoString(&buf[0]))
		libc.StrCat(&buf2[0], libc.CString("\r\n"))
		i = 1
		buf[0] = '\x00'
		if (func() *byte {
			s = *d.Str
			return s
		}()) == nil {
			write_to_output(d, libc.CString("Buffer is empty, nowhere to insert.\r\n"))
			return
		}
		if line_low > 0 {
			for s != nil && i < line_low {
				if (func() *byte {
					s = libc.StrChr(s, '\n')
					return s
				}()) != nil {
					i++
					s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				}
			}
			if i < line_low || s == nil {
				write_to_output(d, libc.CString("Line number out of range; insert aborted.\r\n"))
				return
			}
			temp = int8(*s)
			*s = '\x00'
			if (libc.StrLen(*d.Str) + libc.StrLen(&buf2[0]) + libc.StrLen((*byte)(unsafe.Add(unsafe.Pointer(s), 1))) + 3) > int(d.Max_str) {
				*s = byte(temp)
				write_to_output(d, libc.CString("Insert text pushes buffer over maximum size, insert aborted.\r\n"))
				return
			}
			if *d.Str != nil && **d.Str != 0 {
				libc.StrCat(&buf[0], *d.Str)
			}
			*s = byte(temp)
			libc.StrCat(&buf[0], &buf2[0])
			if s != nil && *s != 0 {
				libc.StrCat(&buf[0], s)
			}
			*d.Str = (*byte)(libc.Realloc(unsafe.Pointer(*d.Str), libc.StrLen(&buf[0])*int(unsafe.Sizeof(int8(0)))+3))
			libc.StrCpy(*d.Str, &buf[0])
			write_to_output(d, libc.CString("Line inserted.\r\n"))
		} else {
			write_to_output(d, libc.CString("Line number must be higher than 0.\r\n"))
			return
		}
	case PARSE_EDIT:
		half_chop(string_, &buf[0], &buf2[0])
		if buf[0] == '\x00' {
			write_to_output(d, libc.CString("You must specify a line number at which to change text.\r\n"))
			return
		}
		line_low = libc.Atoi(libc.GoString(&buf[0]))
		libc.StrCat(&buf2[0], libc.CString("\r\n"))
		i = 1
		buf[0] = '\x00'
		if (func() *byte {
			s = *d.Str
			return s
		}()) == nil {
			write_to_output(d, libc.CString("Buffer is empty, nothing to change.\r\n"))
			return
		}
		if line_low > 0 {
			for s != nil && i < line_low {
				if (func() *byte {
					s = libc.StrChr(s, '\n')
					return s
				}()) != nil {
					i++
					s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				}
			}
			if s == nil || i < line_low {
				write_to_output(d, libc.CString("Line number out of range; change aborted.\r\n"))
				return
			}
			if s != *d.Str {
				temp = int8(*s)
				*s = '\x00'
				libc.StrCat(&buf[0], *d.Str)
				*s = byte(temp)
			}
			libc.StrCat(&buf[0], &buf2[0])
			if (func() *byte {
				s = libc.StrChr(s, '\n')
				return s
			}()) != nil {
				s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
				libc.StrCat(&buf[0], s)
			}
			if libc.StrLen(&buf[0]) > int(d.Max_str) {
				write_to_output(d, libc.CString("Change causes new length to exceed buffer maximum size, aborted.\r\n"))
				return
			}
			*d.Str = (*byte)(libc.Realloc(unsafe.Pointer(*d.Str), libc.StrLen(&buf[0])*int(unsafe.Sizeof(int8(0)))+3))
			libc.StrCpy(*d.Str, &buf[0])
			write_to_output(d, libc.CString("Line changed.\r\n"))
		} else {
			write_to_output(d, libc.CString("Line number must be higher than 0.\r\n"))
			return
		}
	default:
		write_to_output(d, libc.CString("Invalid option.\r\n"))
		mudlog(BRF, ADMLVL_IMPL, 1, libc.CString("SYSERR: invalid command passed to parse_edit_action"))
		return
	}
}
func format_text(ptr_string **byte, mode int, d *descriptor_data, maxlen uint, low int, high int) bool {
	var (
		line_chars    int
		cap_next      int = 1
		cap_next_next int = 0
		color_chars   int = 0
		i             int
		pass_line     int = 0
		flow          *byte
		start         *byte = nil
		temp          int8
		formatted     [64936]byte = func() [64936]byte {
			var t [64936]byte
			copy(t[:], []byte(""))
			return t
		}()
	)
	if d.Max_str > MAX_STRING_LENGTH {
		basic_mud_log(libc.CString("SYSERR: format_text: max_str is greater than buffer size."))
		return false
	}
	if (func() *byte {
		flow = *ptr_string
		return flow
	}()) == nil {
		return false
	}
	var str [64936]byte
	libc.StrCpy(&str[0], flow)
	for i = 0; i < low-1; i++ {
		start = libc.StrTok(&str[0], libc.CString("\n"))
		if start == nil {
			write_to_output(d, libc.CString("There aren't that many lines!\r\n"))
			return false
		}
		libc.StrCat(&formatted[0], libc.StrCat(start, libc.CString("\n")))
		flow = libc.StrStr(flow, libc.CString("\n"))
		libc.StrCpy(&str[0], func() *byte {
			p := &flow
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return *p
		}())
	}
	if IS_SET(uint32(int32(mode)), 1<<0) {
		libc.StrCat(&formatted[0], libc.CString("   "))
		line_chars = 3
	} else {
		line_chars = 0
	}
	for *flow != 0 && i < high {
		for *flow != 0 && libc.StrChr(libc.CString("\n\r\f\t\v "), *flow) != nil {
			if *flow == '\n' && pass_line == 0 {
				if func() int {
					p := &i
					x := *p
					*p++
					return x
				}() >= high {
					pass_line = 1
					break
				}
			}
			flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
		}
		if *flow != 0 {
			start = flow
			for *flow != 0 && libc.StrChr(libc.CString("\n\r\f\t\v .?!"), *flow) == nil {
				if *flow == '@' {
					if *((*byte)(unsafe.Add(unsafe.Pointer(flow), 1))) == '@' {
						color_chars++
					} else {
						color_chars += 2
					}
					flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
				}
				flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
			}
			if cap_next_next != 0 {
				cap_next_next = 0
				cap_next = 1
			}
			for libc.StrChr(libc.CString(".!?"), *flow) != nil {
				cap_next_next = 1
				flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
			}
			if libc.StrChr(libc.CString("\n\r"), *flow) != nil {
				*flow = '\x00'
				flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
				if *flow == '\n' && func() int {
					p := &i
					x := *p
					*p++
					return x
				}() >= high {
					pass_line = 1
				}
				for *flow != 0 && libc.StrChr(libc.CString("\n\r"), *flow) != nil && pass_line == 0 {
					flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
					if *flow == '\n' && func() int {
						p := &i
						x := *p
						*p++
						return x
					}() >= high {
						pass_line = 1
					}
				}
				temp = int8(*flow)
			} else {
				temp = int8(*flow)
				*flow = '\x00'
			}
			if line_chars+libc.StrLen(start)+1-color_chars > PAGE_WIDTH {
				libc.StrCat(&formatted[0], libc.CString("\r\n"))
				line_chars = 0
				color_chars = count_color_chars(start)
			}
			if cap_next == 0 {
				if line_chars > 0 {
					libc.StrCat(&formatted[0], libc.CString(" "))
					line_chars++
				}
			} else {
				cap_next = 0
				CAP(start)
			}
			line_chars += libc.StrLen(start)
			libc.StrCat(&formatted[0], start)
			*flow = byte(temp)
		}
		if cap_next_next != 0 && *flow != 0 {
			if line_chars+3-color_chars > PAGE_WIDTH {
				libc.StrCat(&formatted[0], libc.CString("\r\n"))
				line_chars = 0
				color_chars = count_color_chars(start)
			} else if *flow == '"' || *flow == '\'' {
				var buf [64936]byte
				stdio.Sprintf(&buf[0], "%c ", *flow)
				libc.StrCat(&formatted[0], &buf[0])
				flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), 1))
				line_chars++
			} else {
				libc.StrCat(&formatted[0], libc.CString(" "))
				line_chars += 2
			}
		}
	}
	if *flow != 0 {
		libc.StrCat(&formatted[0], libc.CString("\r\n"))
	}
	libc.StrCat(&formatted[0], flow)
	if *flow == 0 {
		libc.StrCat(&formatted[0], libc.CString("\r\n"))
	}
	if libc.StrLen(&formatted[0])+1 > int(maxlen) {
		formatted[maxlen-1] = '\x00'
	}
	*ptr_string = (*byte)(libc.Realloc(unsafe.Pointer(*ptr_string), int(MIN(int64(maxlen), int64(libc.StrLen(&formatted[0])+1))*int64(unsafe.Sizeof(int8(0))))))
	libc.StrCpy(*ptr_string, &formatted[0])
	return true
}
func replace_str(string_ **byte, pattern *byte, replacement *byte, rep_all int, max_size uint) int {
	var (
		replace_buffer *byte = nil
		flow           *byte
		jetsam         *byte
		temp           int8
		len_           int
		i              int
	)
	if (libc.StrLen(*string_)-libc.StrLen(pattern))+libc.StrLen(replacement) > int(max_size) {
		return -1
	}
	replace_buffer = (*byte)(unsafe.Pointer(&make([]int8, int(max_size))[0]))
	i = 0
	jetsam = *string_
	flow = *string_
	*replace_buffer = '\x00'
	if rep_all != 0 {
		for (func() *byte {
			flow = libc.StrStr(flow, pattern)
			return flow
		}()) != nil {
			i++
			temp = int8(*flow)
			*flow = '\x00'
			if (libc.StrLen(replace_buffer) + libc.StrLen(jetsam) + libc.StrLen(replacement)) > int(max_size) {
				i = -1
				break
			}
			libc.StrCat(replace_buffer, jetsam)
			libc.StrCat(replace_buffer, replacement)
			*flow = byte(temp)
			flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), libc.StrLen(pattern)))
			jetsam = flow
		}
		libc.StrCat(replace_buffer, jetsam)
	} else {
		if (func() *byte {
			flow = libc.StrStr(*string_, pattern)
			return flow
		}()) != nil {
			i++
			flow = (*byte)(unsafe.Add(unsafe.Pointer(flow), libc.StrLen(pattern)))
			len_ = int((int64(uintptr(unsafe.Pointer(flow)) - uintptr(unsafe.Pointer(*string_)))) - int64(libc.StrLen(pattern)))
			libc.StrNCpy(replace_buffer, *string_, len_)
			libc.StrCat(replace_buffer, replacement)
			libc.StrCat(replace_buffer, flow)
		}
	}
	if i <= 0 {
		return 0
	} else {
		*string_ = (*byte)(libc.Realloc(unsafe.Pointer(*string_), libc.StrLen(replace_buffer)*int(unsafe.Sizeof(int8(0)))+3))
		libc.StrCpy(*string_, replace_buffer)
	}
	libc.Free(unsafe.Pointer(replace_buffer))
	return i
}
