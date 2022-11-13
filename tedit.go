package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func news_string_cleanup(d *descriptor_data, terminator int) {
	var (
		fl      *stdio.File
		storage *byte = libc.CString(LIB_TEXT)
	)
	if storage == nil {
		terminator = STRINGADD_ABORT
	}
	switch terminator {
	case STRINGADD_SAVE:
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(storage), "a")
			return fl
		}()) == nil {
			mudlog(CMP, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: Can't write file '%s'."), storage)
		}
		if *d.Str == nil {
			mudlog(CMP, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: Can't write file '%s'."), storage)
		} else {
			var (
				tmstr  *byte
				mytime libc.Time = libc.GetTime(nil)
			)
			tmstr = libc.AscTime(libc.LocalTime(&mytime))
			*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
			stdio.Fprintf(fl, "#%d %s\n@cUpdated By@D: @C%-13s @cDate@D: @Y%s@n\n", TOP_OF_NEWS, d.Newsbuf, GET_NAME(d.Character), tmstr)
			libc.Free(unsafe.Pointer(d.Newsbuf))
			d.Newsbuf = nil
			strip_cr(*d.Str)
			stdio.Fprintf(fl, "%s\n", *d.Str)
			*d.Str = nil
			fl.Close()
			NEWSUPDATE = libc.GetTime(nil)
			save_mud_time(&time_info)
			var i *descriptor_data
			for i = descriptor_list; i != nil; i = i.Next {
				if !IS_PLAYING(i) {
					continue
				}
				if PLR_FLAGGED(i.Character, PLR_WRITING) {
					continue
				}
				if NEWSUPDATE > i.Character.Lastpl {
					send_to_char(i.Character, libc.CString("\r\n@GA news entry has been made by %s, type 'news %d' to see it.@n\r\n"), GET_NAME(d.Character), TOP_OF_NEWS)
				}
			}
			do_reboot(d.Character, libc.CString("news"), 0, 0)
		}
	case STRINGADD_ABORT:
		write_to_output(d, libc.CString("Edit aborted.\r\n"))
		act(libc.CString("$n stops editing the news."), TRUE, d.Character, nil, nil, TO_ROOM)
	default:
		basic_mud_log(libc.CString("SYSERR: news_string_cleanup: Unknown terminator status."))
	}
	d.Connected = CON_PLAYING
}
func tedit_string_cleanup(d *descriptor_data, terminator int) {
	var (
		fl      *stdio.File
		storage *byte = d.Olc.Storage
	)
	if storage == nil {
		terminator = STRINGADD_ABORT
	}
	switch terminator {
	case STRINGADD_SAVE:
		if (func() *stdio.File {
			fl = stdio.FOpen(libc.GoString(storage), "w")
			return fl
		}()) == nil {
			mudlog(CMP, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: Can't write file '%s'."), storage)
		} else {
			if *d.Str != nil && libc.StrCmp(storage, libc.CString("text/news")) == 0 {
				var (
					tmstr  *byte
					mytime libc.Time = libc.GetTime(nil)
				)
				tmstr = libc.AscTime(libc.LocalTime(&mytime))
				*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
				strip_cr(*d.Str)
				stdio.Fprintf(fl, "\n-----------------------------------------------\n@Y%s @cUpdated By@D: @C%s@n\r\n-----------------------------------------------\n%s\n", tmstr, GET_NAME(d.Character), *d.Str)
			} else if *d.Str != nil {
				strip_cr(*d.Str)
				fl.PutS(*d.Str)
			}
			fl.Close()
			mudlog(CMP, ADMLVL_GOD, TRUE, libc.CString("OLC: %s saves '%s'."), GET_NAME(d.Character), storage)
			write_to_output(d, libc.CString("Saved.\r\n"))
			if libc.StrCmp(storage, libc.CString("text/news")) == 0 {
				NEWSUPDATE = libc.GetTime(nil)
				save_mud_time(&time_info)
				var i *descriptor_data
				for i = descriptor_list; i != nil; i = i.Next {
					if !IS_PLAYING(i) {
						continue
					}
					if PLR_FLAGGED(i.Character, PLR_WRITING) {
						continue
					}
					if NEWSUPDATE > i.Character.Lastpl {
						send_to_char(i.Character, libc.CString("\r\n@GThe NEWS file has been updated, type 'news' to see it.@n\r\n"))
					}
				}
				do_reboot(d.Character, libc.CString("all"), 0, 0)
			}
		}
	case STRINGADD_ABORT:
		write_to_output(d, libc.CString("Edit aborted.\r\n"))
		act(libc.CString("$n stops editing some scrolls."), TRUE, d.Character, nil, nil, TO_ROOM)
	default:
		basic_mud_log(libc.CString("SYSERR: tedit_string_cleanup: Unknown terminator status."))
	}
	cleanup_olc(d, CLEANUP_ALL)
	d.Connected = CON_PLAYING
}
func do_tedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		l       int
		i       int = 0
		field   [2048]byte
		backstr *byte = nil
		fields  [11]struct {
			Cmd      *byte
			Level    int8
			Buffer   **byte
			Size     int
			Filename *byte
		} = [11]struct {
			Cmd      *byte
			Level    int8
			Buffer   **byte
			Size     int
			Filename *byte
		}{{Cmd: libc.CString("credits"), Level: ADMLVL_IMPL, Buffer: &credits, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("donottouch"), Level: 6, Buffer: &news, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("motd"), Level: ADMLVL_IMPL, Buffer: &motd, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("imotd"), Level: ADMLVL_IMPL, Buffer: &imotd, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("help"), Level: ADMLVL_GRGOD, Buffer: &help, Size: 24000, Filename: libc.CString(LIB_TEXT_HELP)}, {Cmd: libc.CString("info"), Level: ADMLVL_GRGOD, Buffer: &info, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("background"), Level: ADMLVL_IMPL, Buffer: &background, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("handbook"), Level: ADMLVL_IMPL, Buffer: &handbook, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("update"), Level: ADMLVL_IMPL, Buffer: &policies, Size: 24000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("ihelp"), Level: ADMLVL_GRGOD, Buffer: &ihelp, Size: 24000, Filename: libc.CString(LIB_TEXT_HELP)}, {Cmd: libc.CString("\n"), Level: 0, Buffer: nil, Size: 0, Filename: nil}}
	)
	if ch.Desc == nil {
		return
	}
	one_argument(argument, &field[0])
	if field[0] == 0 {
		send_to_char(ch, libc.CString("Files available to be edited:\r\n"))
		for l = 0; *fields[l].Cmd != '\n'; l++ {
			if ch.Admlevel >= int(fields[l].Level) {
				send_to_char(ch, libc.CString("%-11.11s "), fields[l].Cmd)
				if (func() int {
					p := &i
					*p++
					return *p
				}() % 7) == 0 {
					send_to_char(ch, libc.CString("\r\n"))
				}
			}
		}
		if i%7 != 0 {
			send_to_char(ch, libc.CString("\r\n"))
		}
		if i == 0 {
			send_to_char(ch, libc.CString("None.\r\n"))
		}
		return
	}
	for l = 0; *fields[l].Cmd != '\n'; l++ {
		if libc.StrNCmp(&field[0], fields[l].Cmd, libc.StrLen(&field[0])) == 0 {
			break
		}
	}
	if *fields[l].Cmd == '\n' {
		send_to_char(ch, libc.CString("Invalid text editor option.\r\n"))
		return
	}
	if ch.Admlevel < int(fields[l].Level) {
		send_to_char(ch, libc.CString("You are not godly enough for that!\r\n"))
		return
	}
	clear_screen(ch.Desc)
	send_editor_help(ch.Desc)
	send_to_char(ch, libc.CString("Edit file below:\r\n\r\n"))
	if ch.Desc.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_tedit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(ch.Desc.Olc))
	}
	ch.Desc.Olc = new(oasis_olc_data)
	if *fields[l].Buffer != nil {
		send_to_char(ch, libc.CString("%s"), *fields[l].Buffer)
		backstr = libc.StrDup(*fields[l].Buffer)
	}
	ch.Desc.Olc.Storage = libc.StrDup(fields[l].Filename)
	string_write(ch.Desc, fields[l].Buffer, uint64(fields[l].Size), 0, unsafe.Pointer(backstr))
	act(libc.CString("$n begins editing a text file."), TRUE, ch, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	ch.Desc.Connected = CON_TEDIT
}
