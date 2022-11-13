package main

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unicode"
	"unsafe"
)

const EXE_FILE = "bin/circle"
const PC = 1
const NPC = 2
const BOTH = 3
const MISC = 0
const BINARY = 1
const NUMBER = 2
const PLIST_FORMAT = "players [minlev[-maxlev]] [-n name] [-d days] [-h hours] [-m]"
const MAX_LEVEL_ALLOWED = 100
const MAX_MOB_DAM_ALLOWED = 500
const MAX_DAM_ALLOWED = 50
const MAX_AFFECTS_ALLOWED = 3
const TOTAL_WEAR_CHECKS = 17
const CAN_WEAR_WEAPONS = 0
const MAX_APPLIES_LIMIT = 1
const CHECK_ITEM_RENT = 0
const CHECK_ITEM_COST = 0
const MAX_APPLY_ACCURCY_MOD_TOTAL = 5
const MAX_APPLY_DAMAGE_MOD_TOTAL = 5
const MIN_ROOM_DESC_LENGTH = 80
const MAX_COLOUMN_WIDTH = 80

var copyover_timer int = 0

func do_lag(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d    *descriptor_data
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: lag (target) (number of seconds)\r\n"))
		return
	}
	var found int = FALSE
	var num int = libc.Atoi(libc.GoString(&arg2[0]))
	if num <= 0 || num > 5 {
		send_to_char(ch, libc.CString("Keep it between 1 to 5 seconds please.\r\n"))
		return
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if libc.StrCaseCmp(CAP(GET_NAME(d.Character)), CAP(&arg[0])) == 0 {
			if d.Character.Admlevel > ch.Admlevel {
				send_to_char(ch, libc.CString("Sorry, you've been outranked.\r\n"))
				return
			}
			switch num {
			case 1:
				WAIT_STATE(d.Character, (int(1000000/OPT_USEC))*1)
			case 2:
				WAIT_STATE(d.Character, (int(1000000/OPT_USEC))*2)
			case 3:
				WAIT_STATE(d.Character, (int(1000000/OPT_USEC))*3)
			case 4:
				WAIT_STATE(d.Character, (int(1000000/OPT_USEC))*4)
			case 5:
				WAIT_STATE(d.Character, (int(1000000/OPT_USEC))*5)
			}
			found = TRUE
		}
	}
	if found == FALSE {
		send_to_char(ch, libc.CString("That player isn't around.\r\n"))
		return
	}
}
func update_space() {
	var (
		mapfile    *stdio.File
		rowcounter int
		colcounter int
		vnum_read  int
	)
	basic_mud_log(libc.CString("Updated Space Map. "))
	mapfile = stdio.FOpen("../lib/surface.map", "r")
	for rowcounter = 0; rowcounter <= MAP_ROWS; rowcounter++ {
		for colcounter = 0; colcounter <= MAP_COLS; colcounter++ {
			stdio.Fscanf(mapfile, "%d", &vnum_read)
			mapnums[rowcounter][colcounter] = int(real_room(room_vnum(vnum_read)))
		}
	}
	mapfile.Close()
}
func do_news(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var fl *stdio.File
	var filename *byte
	var line [256]byte
	var arg [2048]byte
	var buf [64936]byte
	var title [256]byte
	var lastline [256]byte
	var entries int = 0
	var lookup int = 0
	var nr int
	var found int = FALSE
	var first int = TRUE
	var exit int = FALSE
	one_argument(argument, &arg[0])
	filename = libc.CString(LIB_TEXT)
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: news (number | list)\r\n"))
		return
	}
	lookup = libc.Atoi(libc.GoString(&arg[0]))
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(filename), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: opening news file for reading"))
		return
	}
	if lookup > 0 {
		for int(fl.IsEOF()) == 0 && exit == FALSE {
			get_line(fl, &line[0])
			if line[0] == '#' {
				if stdio.Sscanf(&line[0], "#%d", &nr) != 1 {
					continue
				} else if nr != lookup {
					entries++
					continue
				} else {
					stdio.Sscanf(&line[0], "#%d %50[0-9a-zA-Z,.!' ]s\n", &nr, &title[0])
					stdio.Sprintf(&buf[0], "@w--------------------------------------------------------------\n@cNum@W: @D(@G%3d@D)                @cTitle@W: @g%-50s@n\n", nr, &title[0])
					found = TRUE
					for int(fl.IsEOF()) == 0 && exit == FALSE {
						get_line(fl, &line[0])
						if line[0] != '#' {
							if first == TRUE {
								first = FALSE
								stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n@w--------------------------------------------------------------\n", &line[0])
								stdio.Sprintf(&lastline[0], "%s", &line[0])
							} else if libc.StrCaseCmp(&line[0], &lastline[0]) == 0 {
								continue
							} else {
								stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%s\n", &line[0])
								stdio.Sprintf(&lastline[0], "%s", &line[0])
							}
						} else {
							exit = TRUE
						}
					}
				}
			}
		}
		fl.Close()
	} else if libc.StrCaseCmp(&arg[0], libc.CString("list")) == 0 {
		for int(fl.IsEOF()) == 0 {
			get_line(fl, &line[0])
			if line[0] == '#' {
				entries++
				if stdio.Sscanf(&line[0], "#%d", &nr) != 1 {
					continue
				} else {
					if first == TRUE {
						stdio.Sscanf(&line[0], "#%d %50[0-9a-zA-Z,.!' ]s\n", &nr, &title[0])
						stdio.Sprintf(&buf[0], "@wNews Entries (Newest at the bottom, to read an entry use 'news (number)')\n@D[@cNum@W: @D(@G%3d@D) @cTitle@W: @g%-50s@D]@n\n", nr, &title[0])
						first = FALSE
					} else {
						stdio.Sscanf(&line[0], "#%d %50[0-9a-zA-Z,.!' ]s\n", &nr, &title[0])
						stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@D[@cNum@W: @D(@G%3d@D) @cTitle@W: @g%-50s@D]@n\n", nr, &title[0])
					}
				}
			}
		}
		fl.Close()
		if entries > 0 {
			ch.Lastpl = libc.GetTime(nil)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			page_string(ch.Desc, &buf[0], TRUE)
		} else {
			send_to_char(ch, libc.CString("The news file is empty right now.\r\n"))
		}
		buf[0] = '\x00'
		title[0] = '\x00'
		lastline[0] = '\x00'
		return
	} else {
		fl.Close()
		send_to_char(ch, libc.CString("Syntax: news (number | list)\r\n"))
		return
	}
	if found == TRUE {
		send_to_char(ch, libc.CString("%s\r\n"), &buf[0])
		ch.Lastpl = libc.GetTime(nil)
		buf[0] = '\x00'
		title[0] = '\x00'
		lastline[0] = '\x00'
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		return
	} else {
		send_to_char(ch, libc.CString("That news entry does not exist.\r\n"))
		return
	}
}
func do_newsedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Admlevel < 1 || IS_NPC(ch) {
		return
	}
	var fl *stdio.File
	var filename *byte
	var line [256]byte
	var entries int = 0
	var lookup int = 0
	var lastentry int = 0
	var nr int
	filename = libc.CString(LIB_TEXT)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Syntax: newsedit (title)\r\n"))
		return
	} else if libc.StrLen(argument) > 50 {
		send_to_char(ch, libc.CString("Limit of 50 characters for title.\r\n"))
		return
	} else if libc.StrStr(argument, libc.CString("#")) != nil {
		send_to_char(ch, libc.CString("# is a forbidden character for news entries as it is used by the file system.\r\n"))
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(filename), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Couldn't open news file for reading"))
		return
	}
	for int(fl.IsEOF()) == 0 {
		get_line(fl, &line[0])
		if line[0] == '#' {
			entries++
			if stdio.Sscanf(&line[0], "#%d", &nr) != 1 {
				continue
			} else {
				lastentry = nr
			}
		}
	}
	fl.Close()
	var fields [2]struct {
		Cmd      *byte
		Level    int8
		Buffer   **byte
		Size     int
		Filename *byte
	} = [2]struct {
		Cmd      *byte
		Level    int8
		Buffer   **byte
		Size     int
		Filename *byte
	}{{Cmd: libc.CString("news"), Level: ADMLVL_IMMORT, Buffer: &immlist, Size: 2000, Filename: libc.CString(LIB_TEXT)}, {Cmd: libc.CString("\n"), Level: 0, Buffer: nil, Size: 0, Filename: nil}}
	var tmstr *byte
	var mytime libc.Time = libc.GetTime(nil)
	tmstr = libc.AscTime(libc.LocalTime(&mytime))
	*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
	if lastentry == 0 {
		if entries == 0 {
			lookup = 1
		} else {
			send_to_imm(libc.CString("ERROR: News file entries are disorganized. Report to Iovan for analysis.\r\n"))
			return
		}
	} else {
		lookup = lastentry + 1
	}
	var backstr *byte = nil
	act(libc.CString("$n begins to edit the news."), TRUE, ch, nil, nil, TO_ROOM)
	send_to_char(ch, libc.CString("@D----------------------=[@GNews Edit@D]=----------------------@n\n"))
	send_to_char(ch, libc.CString(" @RRemember that using # in newsedit is not possible. That\n"))
	send_to_char(ch, libc.CString("character will be eaten because it is required for the news\n"))
	send_to_char(ch, libc.CString("file as a delimiter. Also if you want to create an empty line\n"))
	send_to_char(ch, libc.CString("between paragraphs you will need to enter a single space and\n"))
	send_to_char(ch, libc.CString("not just push enter. Happy editing!@n\n"))
	send_to_char(ch, libc.CString("@D---------------------------------------------------------@n\n"))
	send_editor_help(ch.Desc)
	skip_spaces(&argument)
	ch.Desc.Newsbuf = libc.StrDup(argument)
	TOP_OF_NEWS = lookup
	LASTNEWS = lookup
	string_write(ch.Desc, fields[0].Buffer, 2000, 0, unsafe.Pointer(backstr))
	ch.Desc.Connected = CON_NEWSEDIT
}
func print_lockout(ch *char_data) {
	if IS_NPC(ch) {
		return
	}
	var file *stdio.File
	var fname [40]byte
	var filler [50]byte
	var line [256]byte
	var buf [259744]byte
	var count int = 0
	var first int = TRUE
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("lockout")) == 0 {
		send_to_char(ch, libc.CString("The lockout file does not exist."))
		return
	} else if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		send_to_char(ch, libc.CString("The lockout file does not exist."))
		return
	}
	stdio.Sprintf(&buf[0], "@b------------------[ @RLOCKOUT @b]------------------@n\n")
	for int(file.IsEOF()) == 0 {
		get_line(file, &line[0])
		stdio.Sscanf(&line[0], "%s\n", &filler[0])
		if first != TRUE && libc.StrStr(&buf[0], &filler[0]) == nil {
			if count == 0 {
				stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%-23s@D|@n", &filler[0])
				count = 1
			} else {
				stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%-23s\n", &filler[0])
				count = 0
			}
		} else if first == TRUE {
			stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "%-23s@D|@n", &filler[0])
			first = FALSE
			count = 1
		}
		filler[0] = '\x00'
	}
	if count == 1 {
		stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "\n")
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@b------------------[ @RLOCKOUT @b]------------------@n\n")
	page_string(ch.Desc, &buf[0], 0)
	file.Close()
}
func do_approve(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		vict *char_data = nil
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What player do you want to approve as having an acceptable bio?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That player is not in the game.\r\n"))
		return
	}
	if PLR_FLAGGED(vict, PLR_BIOGR) {
		send_to_char(ch, libc.CString("They have already been approved. If this was made in error inform Iovan.\r\n"))
		return
	} else {
		vict.Act[int(PLR_BIOGR/32)] |= bitvector_t(int32(1 << (int(PLR_BIOGR % 32))))
		send_to_char(ch, libc.CString("They have now been approved.\r\n"))
		return
	}
}
func lockWrite(ch *char_data, name *byte) {
	var (
		file   *stdio.File
		fname  [40]byte
		filler [50]byte
		line   [256]byte
		names  [500]*byte = [500]*byte{0: libc.CString("")}
		fl     *stdio.File
		count  int = 0
		x      int = 0
		found  int = FALSE
	)
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("lockout")) == 0 {
		send_to_char(ch, libc.CString("The lockout file does not exist."))
		return
	} else if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fname[0]), "r")
		return file
	}()) == nil {
		send_to_char(ch, libc.CString("The lockout file does not exist."))
		return
	}
	for int(file.IsEOF()) == 0 || count < 498 {
		get_line(file, &line[0])
		stdio.Sscanf(&line[0], "%s\n", &filler[0])
		names[count] = libc.StrDup(&filler[0])
		count++
		filler[0] = '\x00'
	}
	file.Close()
	if get_filename(&fname[0], uint64(40), INTRO_FILE, libc.CString("lockout")) == 0 {
		return
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save Lockout File, %s."), &fname[0])
		return
	}
	for x < count {
		if x == 0 || libc.StrCaseCmp(names[x-1], names[x]) != 0 {
			if libc.StrCaseCmp(names[x], CAP(name)) != 0 {
				stdio.Fprintf(fl, "%s\n", CAP(names[x]))
			} else {
				found = TRUE
			}
		}
		x++
	}
	if found == FALSE {
		stdio.Fprintf(fl, "%s\n", CAP(name))
		send_to_all(libc.CString("@rLOCKOUT@D: @WThe character, @C%s@W, was locked out of the MUD by @c%s@W.@n\r\n"), CAP(name), GET_NAME(ch))
		basic_mud_log(libc.CString("LOCKOUT: %s sentenced by %s."), CAP(name), GET_NAME(ch))
		log_imm_action(libc.CString("LOCKOUT: %s sentenced by %s."), CAP(name), GET_NAME(ch))
	} else {
		send_to_all(libc.CString("@rLOCKOUT@D: @WThe character, @C%s@W, has had lockout removed by @c%s@W.@n\r\n"), CAP(name), GET_NAME(ch))
		basic_mud_log(libc.CString("LOCKOUT: %s sentenced by %s."), CAP(name), GET_NAME(ch))
		log_imm_action(libc.CString("LOCKOUT: %s sentenced by %s."), CAP(name), GET_NAME(ch))
	}
	fl.Close()
	return
}
func do_reward(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		k    *descriptor_data
		amt  int = 0
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: reward (target) (amount)\r\nThis is either a positive number or a negative.\r\n"))
		return
	}
	if readUserIndex(&arg[0]) == 0 {
		send_to_char(ch, libc.CString("That is not a recognised user file.\r\n"))
		return
	}
	amt = libc.Atoi(libc.GoString(&arg2[0]))
	if amt == 0 {
		send_to_char(ch, libc.CString("That is pointless don't you think? Try an amount higher than 0.\r\n"))
		return
	}
	for k = descriptor_list; k != nil; k = k.Next {
		if IS_NPC(k.Character) {
			continue
		}
		if k.Connected != CON_PLAYING {
			continue
		}
		if libc.StrCaseCmp(k.User, &arg[0]) == 0 {
			vict = k.Character
		}
	}
	if amt > 0 {
		if vict != nil {
			send_to_char(ch, libc.CString("@WYou award @C%s @D(@G%d@D)@W RP points.@n\r\n"), GET_NAME(vict), amt)
			send_to_char(vict, libc.CString("@D[@YROLEPLAY@D] @WYou have been awarded @D(@G%d@D)@W RP points by @C%s@W.@n\r\n"), amt, GET_NAME(ch))
			send_to_imm(libc.CString("ROLEPLAY: %s has been awarded %d RP points by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has been awarded %d RP points by %s."), &arg[0], amt, GET_NAME(ch))
			vict.Rp += amt
			vict.Rp = vict.Rp
			if amt <= 29 {
				vict.Trp += amt
			}
			vict.Desc.Rpp += amt
			userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			send_to_char(ch, libc.CString("@WYou award user @C%s @D(@G%d@D)@W RP points.@n\r\n"), &arg[0], amt)
			send_to_imm(libc.CString("ROLEPLAY: User %s has been awarded %d RP points by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has been awarded %d RP points by %s."), &arg[0], amt, GET_NAME(ch))
			userWrite(nil, 0, amt, 0, &arg[0])
		}
	} else {
		if vict != nil {
			send_to_char(ch, libc.CString("@WYou deduct @D(@G%d@D)@W RP points from @C%s@W.@n\r\n"), amt, GET_NAME(vict))
			send_to_char(vict, libc.CString("@D[@YROLEPLAY@D] @C%s@W deducts @D(@G%d@D)@W RP points from you.@n\r\n"), GET_NAME(ch), amt)
			send_to_imm(libc.CString("ROLEPLAY: %s has had %d RP points deducted by %s."), GET_NAME(vict), amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has had %d RP points deducted by %s."), GET_NAME(vict), amt, GET_NAME(ch))
			vict.Rp += amt
			vict.Rp = vict.Rp
			vict.Desc.Rpp += amt
			userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			send_to_char(ch, libc.CString("@WYou deduct @D(@G%d@D)@W RP points from user @C%s@W.@n\r\n"), amt, &arg[0])
			send_to_imm(libc.CString("ROLEPLAY: User %s has had %d RP points deducted by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has been awarded %d RP points by %s."), &arg[0], amt, GET_NAME(ch))
			userWrite(nil, 0, amt, 0, &arg[0])
		}
	}
}
func do_rbank(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		k    *descriptor_data
		amt  int = 0
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: rbank (target) (amount)\r\nThis is either a positive number or a negative.\r\n"))
		return
	}
	if readUserIndex(&arg[0]) == 0 {
		send_to_char(ch, libc.CString("That is not a recognised user file.\r\n"))
		return
	}
	amt = libc.Atoi(libc.GoString(&arg2[0]))
	if amt == 0 {
		send_to_char(ch, libc.CString("That is pointless don't you think? Try an amount higher than 0.\r\n"))
		return
	}
	for k = descriptor_list; k != nil; k = k.Next {
		if IS_NPC(k.Character) {
			continue
		}
		if k.Connected != CON_PLAYING {
			continue
		}
		if libc.StrCaseCmp(k.User, &arg[0]) == 0 {
			vict = k.Character
		}
	}
	if amt > 0 {
		if vict != nil {
			send_to_char(ch, libc.CString("@WYou put @D(@G%d@D)@W RP points into @C%s@W's Bank.@n\r\n"), amt, GET_NAME(vict))
			send_to_char(vict, libc.CString("@D[@YROLEPLAY@D] @WYou have had @D(@G%d@D)@W RP points added to your RPP Bank by @C%s@W.@n\r\n"), amt, GET_NAME(ch))
			send_to_imm(libc.CString("ROLEPLAY: %s has had %d RP points put into their RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has had %d RP points put into their RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			vict.Rbank += amt
			vict.Rbank = vict.Rbank
			if amt <= 29 {
				vict.Trp += amt
			}
			vict.Desc.Rbank += amt
			userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			send_to_char(ch, libc.CString("@WYou put @D(@G%d@D)@W RP points into @C%s@W's Bank.@n\r\n"), amt, &arg[0])
			send_to_imm(libc.CString("ROLEPLAY: %s has had %d RP points put into their RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has had %d RP points put into their RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			userWrite(nil, 0, 0, amt, &arg[0])
		}
	} else {
		if vict != nil {
			send_to_char(ch, libc.CString("@WYou deduct @D(@G%d@D)@W RP points from @C%s@W's Bank.@n\r\n"), amt, GET_NAME(vict))
			send_to_char(vict, libc.CString("@D[@YROLEPLAY@D] @C%s@W deducts @D(@G%d@D)@W RP points from your RPP Bank.@n\r\n"), GET_NAME(ch), amt)
			send_to_imm(libc.CString("ROLEPLAY: %s has had %d RP points deducted from their RPP Bank by %s."), GET_NAME(vict), amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has had %d RP points deducted from their RPP Bank by %s."), GET_NAME(vict), amt, GET_NAME(ch))
			vict.Rbank += amt
			vict.Rbank = vict.Rbank
			vict.Desc.Rbank += amt
			userWrite(vict.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			send_to_char(ch, libc.CString("@WYou deduct @D(@G%d@D)@W RP points from @C%s@W's Bank.@n\r\n"), amt, &arg[0])
			send_to_imm(libc.CString("ROLEPLAY: %s has had %d RP points deducted from their RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			log_imm_action(libc.CString("ROLEPLAY: %s has had %d RP points put deducted from RPP Bank by %s."), &arg[0], amt, GET_NAME(ch))
			userWrite(nil, 0, 0, amt, &arg[0])
		}
	}
}
func do_permission(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("You want to @Grestrict@n or @Gunrestrict@n?\r\n"))
		return
	}
	if arg2[0] == 0 && libc.StrCaseCmp(libc.CString("unrestrict"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("You want to unrestrict which race? @Gsaiyan @nor @Gmajin@n?\r\n"))
		return
	}
	if libc.StrCaseCmp(libc.CString("unrestrict"), &arg[0]) == 0 {
		if libc.StrCaseCmp(libc.CString("saiyan"), &arg2[0]) == 0 {
			send_to_char(ch, libc.CString("You have unrestricted saiyans for the very next character creation.\r\n"))
			send_to_imm(libc.CString("PERMISSION: %s unrestricted saiyans."), GET_NAME(ch))
			SAIYAN_ALLOWED = TRUE
		} else if libc.StrCaseCmp(libc.CString("majin"), &arg2[0]) == 0 {
			send_to_char(ch, libc.CString("You have unrestricted majins for the very next character creation.\r\n"))
			send_to_imm(libc.CString("PERMISSION: %s unrestricted majins."), GET_NAME(ch))
			MAJIN_ALLOWED = TRUE
		} else {
			send_to_char(ch, libc.CString("You want to unrestrict which race? @Gsaiyan @nor @Gmajin@n?\r\n"))
			return
		}
	} else if libc.StrCaseCmp(libc.CString("restrict"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("You have restricted character creation to standard race slection.\r\n"))
		send_to_imm(libc.CString("PERMISSION: %s restricted races again."), GET_NAME(ch))
		MAJIN_ALLOWED = FALSE
	} else {
		send_to_char(ch, libc.CString("You want to @Grestrict@n or @Gunrestrict@n?\r\n"))
		return
	}
}
func do_transobj(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj  *obj_data
		vict *char_data
		d    *descriptor_data
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if !IS_NPC(ch) && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("Huh!?"))
		return
	}
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: transo (object) (target)\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You want to send what?\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("all"), &arg2[0]) == 0 {
		var (
			num  int       = int(GET_OBJ_VNUM(obj))
			obj2 *obj_data = nil
		)
		act(libc.CString("You send $p to everyone in the game."), TRUE, ch, obj, nil, TO_CHAR)
		for d = descriptor_list; d != nil; d = d.Next {
			if IS_NPC(d.Character) {
				continue
			} else if !IS_PLAYING(d) {
				continue
			} else if d.Character == ch {
				continue
			} else {
				act(libc.CString("$N sends $p across the universe to you."), TRUE, d.Character, obj, unsafe.Pointer(ch), TO_CHAR)
				obj2 = read_object(obj_vnum(num), VIRTUAL)
				obj_to_char(obj2, d.Character)
			}
		}
	} else if (func() *char_data {
		vict = get_char_vis(ch, &arg2[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That player is not in the game.\r\n"))
		return
	} else {
		act(libc.CString("You send $p to $N."), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n sends $p across the universe to you."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		obj_from_char(obj)
		obj_to_char(obj, vict)
		return
	}
}
func search_replace(string_ *byte, find *byte, replace *byte) {
	var (
		final [64936]byte
		temp  [2]byte
		start uint64
		end   uint64
		i     uint64
	)
	for libc.StrStr(string_, find) != nil {
		final[0] = '\x00'
		start = uint64(int64(uintptr(unsafe.Pointer(libc.StrStr(string_, find))) - uintptr(unsafe.Pointer(string_))))
		end = start + uint64(libc.StrLen(find))
		temp[1] = '\x00'
		libc.StrNCat(&final[0], string_, int(start))
		libc.StrCat(&final[0], replace)
		for i = end; *(*byte)(unsafe.Add(unsafe.Pointer(string_), i)) != '\x00'; i++ {
			temp[0] = *(*byte)(unsafe.Add(unsafe.Pointer(string_), i))
			libc.StrCat(&final[0], &temp[0])
		}
		stdio.Sprintf(string_, libc.GoString(&final[0]))
	}
	return
}
func do_interest(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Admlevel < 5 {
		send_to_char(ch, libc.CString("Huh!?\r\n"))
		return
	} else {
		if INTERESTTIME > 0 {
			var tmstr *byte
			tmstr = libc.AscTime(libc.LocalTime(&INTERESTTIME))
			*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
			send_to_char(ch, libc.CString("INTEREST TIME: [%s]\r\n"), tmstr)
			return
		}
		send_to_char(ch, libc.CString("Interest time has been initiated!\r\n"))
		INTERESTTIME = libc.GetTime(nil) + 86400
		LASTINTEREST = libc.GetTime(nil) + 86400
		return
	}
}
func do_finddoor(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d        int
		vnum     int = int(-1)
		num      int = 0
		len_     uint64
		nlen     uint64
		i        room_rnum
		arg      [2048]byte
		buf      [64936]byte = [64936]byte{}
		tmp_char *char_data
		obj      *obj_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Format: finddoor <obj/vnum>\r\n"))
	} else if is_number(&arg[0]) != 0 {
		vnum = libc.Atoi(libc.GoString(&arg[0]))
		obj = (*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(real_object(obj_vnum(vnum)))))
	} else {
		generic_find(&arg[0], (1<<2)|1<<3|1<<4|1<<5, ch, &tmp_char, &obj)
		if obj == nil {
			send_to_char(ch, libc.CString("What key do you want to find a door for?\r\n"))
		} else {
			vnum = int(GET_OBJ_VNUM(obj))
		}
	}
	if vnum != int(-1) {
		len_ = uint64(stdio.Snprintf(&buf[0], int(64936), "Doors unlocked by key [%d] %s are:\r\n", vnum, obj.Short_description))
		for i = 0; i <= top_of_world; i++ {
			for d = 0; d < NUM_OF_DIRS; d++ {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[d] != nil && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[d].Key != 0 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[d].Key == obj_vnum(vnum) {
					nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[%3d] Room %d, %s (%s)\r\n", func() int {
						p := &num
						*p++
						return *p
					}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number, dirs[d], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[d].Keyword))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				}
			}
		}
		if num > 0 {
			page_string(ch.Desc, &buf[0], 1)
		} else {
			send_to_char(ch, libc.CString("No doors were found for key [%d] %s.\r\n"), vnum, obj.Short_description)
		}
	}
}
func do_recall(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("You are not an immortal!\r\n"))
	} else {
		send_to_char(ch, libc.CString("You disappear in a burst of light!\r\n"))
		act(libc.CString("$n disappears in a burst of light!"), FALSE, ch, nil, nil, TO_ROOM)
		if real_room(2) != room_rnum(-1) {
			char_from_room(ch)
			char_to_room(ch, real_room(2))
			look_at_room(ch.In_room, ch, 0)
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			} else {
				ch.Player_specials.Load_room = -1
			}
		}
	}
}
func do_hell(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: lockout (character)\n        lockout list\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("Iovan")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("iovan")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Fahl")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("fahl")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Xyron")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("xyron")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Samael")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("samael")) == 0 {
		send_to_char(ch, libc.CString("What are you smoking? You can't lockout senior imms.\r\n"))
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("list")) == 0 {
		print_lockout(ch)
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		lockWrite(ch, &arg[0])
		return
	} else {
		var d *descriptor_data = vict.Desc
		Crash_rentsave(vict, 0)
		extract_char(vict)
		lockWrite(ch, GET_NAME(vict))
		if d != nil && d.Connected != CON_PLAYING {
			d.Connected = CON_CLOSE
			vict.Desc.Character = nil
			vict.Desc = nil
		}
		return
	}
	return
}
func do_echo(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	var NoName bool = FALSE != 0
	if *argument == 0 {
		send_to_char(ch, libc.CString("Yes.. but what?\r\n"))
	} else {
		var (
			buf    [2052]byte
			name   [128]byte
			found  int        = FALSE
			trunc  int        = 0
			vict   *char_data = nil
			next_v *char_data = nil
			tch    *char_data = nil
		)
		if libc.StrLen(argument) > 1000 {
			trunc = libc.StrLen(argument) - 1000
			*(*byte)(unsafe.Add(unsafe.Pointer(argument), libc.StrLen(argument)-trunc)) = '\x00'
			stdio.Sprintf(argument, "%s\n@D(@gMessage truncated to 1000 characters@D)@n\n", argument)
		}
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
			next_v = vict.Next_in_room
			if vict == ch {
				continue
			}
			if found == FALSE {
				stdio.Sprintf(&name[0], "*%s", GET_NAME(vict))
				if libc.StrStr(argument, CAP(&name[0])) != nil {
					found = TRUE
					tch = vict
				}
				if found == FALSE && !IS_NPC(vict) {
					if readIntro(ch, vict) == 1 {
						stdio.Sprintf(&name[0], "*%s", get_i_name(ch, vict))
						if libc.StrStr(argument, CAP(&name[0])) != nil {
							found = TRUE
							tch = vict
						}
					}
				}
			}
		}
		if subcmd == SCMD_SMOTE {
			if libc.StrStr(argument, libc.CString("#")) == nil {
				NoName = TRUE != 0
			}
			strlcpy(&buf[0], argument, uint64(2052))
			search_replace(&buf[0], libc.CString("#"), libc.CString("$n"))
			search_replace(&buf[0], libc.CString("&1"), libc.CString("'@C"))
			search_replace(&buf[0], libc.CString("&2"), libc.CString("@w'"))
			if found == TRUE {
				search_replace(&buf[0], &name[0], libc.CString("$N"))
			} else if libc.StrStr(&buf[0], libc.CString("*")) != nil {
				search_replace(&buf[0], libc.CString("*"), libc.CString(""))
			}
		} else if subcmd == SCMD_EMOTE {
			stdio.Snprintf(&buf[0], int(2052), "$n %s", argument)
			search_replace(&buf[0], libc.CString("#"), libc.CString("$n"))
			search_replace(&buf[0], libc.CString("&1"), libc.CString("'@C"))
			search_replace(&buf[0], libc.CString("&2"), libc.CString("@w'"))
			if found == TRUE {
				search_replace(&buf[0], &name[0], libc.CString("$N"))
			}
		} else {
			stdio.Snprintf(&buf[0], int(2052), "%s", argument)
		}
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		}
		if int(libc.BoolToInt(NoName)) == TRUE {
			var blom [2048]byte
			stdio.Sprintf(&blom[0], "@D(@GOOC@W: @gSmote by user %s@D)@n", func() *byte {
				if IS_NPC(ch) {
					return GET_NAME(ch)
				}
				if ch.Desc.User == nil {
					return libc.CString("ERROR REPORT")
				}
				return ch.Desc.User
			}())
			act(&blom[0], FALSE, ch, nil, nil, TO_ROOM)
		}
		if found == FALSE {
			act(&buf[0], FALSE, ch, nil, nil, TO_CHAR)
			act(&buf[0], FALSE, ch, nil, nil, TO_ROOM)
		} else {
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
			search_replace(&buf[0], libc.CString("$N"), libc.CString("you"))
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(tch), TO_VICT)
		}
	}
}
func do_send(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		buf  [2048]byte
		vict *char_data
	)
	half_chop(argument, &arg[0], &buf[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Send what to who?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		return
	}
	send_to_char(vict, libc.CString("%s\r\n"), &buf[0])
	if PRF_FLAGGED(ch, PRF_NOREPEAT) {
		send_to_char(ch, libc.CString("Sent.\r\n"))
	} else {
		send_to_char(ch, libc.CString("You send '%s' to %s.\r\n"), &buf[0], GET_NAME(vict))
	}
}
func find_target_room(ch *char_data, rawroomstr *byte) room_rnum {
	var (
		location room_rnum = room_rnum(-1)
		roomstr  [2048]byte
		rm       *room_data
	)
	one_argument(rawroomstr, &roomstr[0])
	if roomstr[0] == 0 {
		send_to_char(ch, libc.CString("You must supply a room number or name.\r\n"))
		return -1
	}
	if unicode.IsDigit(rune(roomstr[0])) && libc.StrChr(&roomstr[0], '.') == nil {
		if (func() room_rnum {
			location = real_room(room_vnum(libc.Atoi(libc.GoString(&roomstr[0]))))
			return location
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("No room exists with that number.\r\n"))
			return -1
		}
	} else {
		var (
			target_mob *char_data
			target_obj *obj_data
			mobobjstr  *byte = &roomstr[0]
			num        int
		)
		num = get_number(&mobobjstr)
		if (func() *char_data {
			target_mob = get_char_vis(ch, mobobjstr, &num, 1<<1)
			return target_mob
		}()) != nil {
			if (func() room_rnum {
				location = target_mob.In_room
				return location
			}()) == room_rnum(-1) {
				send_to_char(ch, libc.CString("That character is currently lost.\r\n"))
				return -1
			}
		} else if (func() *obj_data {
			target_obj = get_obj_vis(ch, mobobjstr, &num)
			return target_obj
		}()) != nil {
			if target_obj.In_room != room_rnum(-1) {
				location = target_obj.In_room
			} else if target_obj.Carried_by != nil && target_obj.Carried_by.In_room != room_rnum(-1) {
				location = target_obj.Carried_by.In_room
			} else if target_obj.Worn_by != nil && target_obj.Worn_by.In_room != room_rnum(-1) {
				location = target_obj.Worn_by.In_room
			}
			if location == room_rnum(-1) {
				send_to_char(ch, libc.CString("That object is currently not in a room.\r\n"))
				return -1
			}
		}
		if location == room_rnum(-1) {
			send_to_char(ch, libc.CString("Nothing exists by that name.\r\n"))
			return -1
		}
	}
	if ch.Admlevel >= ADMLVL_VICE {
		return location
	}
	rm = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(location)))
	if can_edit_zone(ch, rm.Zone) == 0 && ch.Admlevel < ADMLVL_GOD && ZONE_FLAGGED(rm.Zone, ZONE_QUEST) {
		send_to_char(ch, libc.CString("This target is in a quest zone.\r\n"))
		return -1
	}
	if ch.Admlevel < ADMLVL_VICE && ZONE_FLAGGED(rm.Zone, ZONE_NOIMMORT) {
		send_to_char(ch, libc.CString("This target is in a zone closed to all.\r\n"))
		return -1
	}
	if ROOM_FLAGGED(location, ROOM_GODROOM) {
		send_to_char(ch, libc.CString("You are not godly enough to use that room!\r\n"))
	} else {
		return location
	}
	return -1
}
func do_at(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		command      [2048]byte
		buf          [2048]byte
		location     room_rnum
		original_loc room_rnum
	)
	half_chop(argument, &buf[0], &command[0])
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("You must supply a room number or a name.\r\n"))
		return
	}
	if command[0] == 0 {
		send_to_char(ch, libc.CString("What do you want to do there?\r\n"))
		return
	}
	if (func() room_rnum {
		location = find_target_room(ch, &buf[0])
		return location
	}()) == room_rnum(-1) {
		return
	}
	original_loc = ch.In_room
	char_from_room(ch)
	char_to_room(ch, location)
	command_interpreter(ch, &command[0])
	if ch.In_room == location {
		char_from_room(ch)
		char_to_room(ch, original_loc)
	}
}
func do_goto(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf      [64936]byte
		location room_rnum
	)
	if (func() room_rnum {
		location = find_target_room(ch, argument)
		return location
	}()) == room_rnum(-1) {
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("They are inside a healing tank!\r\n"))
		return
	}
	stdio.Snprintf(&buf[0], int(64936), "$n %s", func() *byte {
		if ch.Player_specials.Poofout != nil {
			return ch.Player_specials.Poofout
		}
		return libc.CString("disappears in a puff of smoke.")
	}())
	act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
	char_from_room(ch)
	char_to_room(ch, location)
	stdio.Snprintf(&buf[0], int(64936), "$n %s", func() *byte {
		if ch.Player_specials.Poofin != nil {
			return ch.Player_specials.Poofin
		}
		return libc.CString("appears with an ear-splitting bang.")
	}())
	act(&buf[0], TRUE, ch, nil, nil, TO_ROOM)
	look_at_room(ch.In_room, ch, 0)
	enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, -1)
}
func do_trans(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		i      *descriptor_data
		victim *char_data
	)
	one_argument(argument, &buf[0])
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to transfer?\r\n"))
	} else if libc.StrCaseCmp(libc.CString("all"), &buf[0]) != 0 {
		if (func() *char_data {
			victim = get_char_vis(ch, &buf[0], nil, 1<<1)
			return victim
		}()) == nil {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		} else if victim == ch {
			send_to_char(ch, libc.CString("That doesn't make much sense, does it?\r\n"))
		} else {
			if ch.Admlevel < victim.Admlevel && !IS_NPC(victim) {
				send_to_char(ch, libc.CString("Go transfer someone your own size.\r\n"))
				return
			}
			if PLR_FLAGGED(victim, PLR_HEALT) {
				send_to_char(ch, libc.CString("They are inside a healing tank!\r\n"))
				return
			}
			act(libc.CString("$n disappears in a mushroom cloud."), FALSE, victim, nil, nil, TO_ROOM)
			char_from_room(victim)
			char_to_room(victim, ch.In_room)
			act(libc.CString("$n arrives from a puff of smoke."), FALSE, victim, nil, nil, TO_ROOM)
			act(libc.CString("$n has transferred you!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
			look_at_room(victim.In_room, victim, 0)
			enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(victim.In_room))), victim, -1)
		}
	} else {
		if !ADM_FLAGGED(ch, ADM_TRANSALL) {
			send_to_char(ch, libc.CString("I think not.\r\n"))
			return
		}
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Connected == CON_PLAYING && i.Character != nil && i.Character != ch {
				victim = i.Character
				if victim.Admlevel >= ch.Admlevel {
					continue
				}
				act(libc.CString("$n disappears in a mushroom cloud."), FALSE, victim, nil, nil, TO_ROOM)
				char_from_room(victim)
				char_to_room(victim, ch.In_room)
				act(libc.CString("$n arrives from a puff of smoke."), FALSE, victim, nil, nil, TO_ROOM)
				act(libc.CString("$n has transferred you!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
				look_at_room(victim.In_room, victim, 0)
				enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(victim.In_room))), victim, -1)
			}
		}
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	}
}
func do_teleport(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		buf2   [2048]byte
		victim *char_data
		target room_rnum
	)
	two_arguments(argument, &buf[0], &buf2[0])
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to teleport?\r\n"))
	} else if (func() *char_data {
		victim = get_char_vis(ch, &buf[0], nil, 1<<1)
		return victim
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if victim == ch {
		send_to_char(ch, libc.CString("Use 'goto' to teleport yourself.\r\n"))
	} else if victim.Admlevel >= ch.Admlevel {
		send_to_char(ch, libc.CString("Maybe you shouldn't do that.\r\n"))
	} else if buf2[0] == 0 {
		send_to_char(ch, libc.CString("Where do you wish to send this person?\r\n"))
	} else if (func() room_rnum {
		target = find_target_room(ch, &buf2[0])
		return target
	}()) != room_rnum(-1) {
		if PLR_FLAGGED(victim, PLR_HEALT) {
			send_to_char(ch, libc.CString("They are inside a healing tank!\r\n"))
			return
		}
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		act(libc.CString("$n disappears in a puff of smoke."), FALSE, victim, nil, nil, TO_ROOM)
		char_from_room(victim)
		char_to_room(victim, target)
		act(libc.CString("$n arrives from a puff of smoke."), FALSE, victim, nil, nil, TO_ROOM)
		act(libc.CString("$n has teleported you!"), FALSE, ch, nil, unsafe.Pointer((*byte)(unsafe.Pointer(victim))), TO_VICT)
		look_at_room(victim.In_room, victim, 0)
		enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(victim.In_room))), victim, -1)
	}
}
func do_vnum(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		buf2 [2048]byte
	)
	half_chop(argument, &buf[0], &buf2[0])
	if buf[0] == 0 || buf2[0] == 0 || is_abbrev(&buf[0], libc.CString("mob")) == 0 && is_abbrev(&buf[0], libc.CString("obj")) == 0 && is_abbrev(&buf[0], libc.CString("mat")) == 0 && is_abbrev(&buf[0], libc.CString("wtype")) == 0 && is_abbrev(&buf[0], libc.CString("atype")) == 0 {
		send_to_char(ch, libc.CString("Usage: vnum { atype | material | mob | obj | wtype } <name>\r\n"))
		return
	}
	if is_abbrev(&buf[0], libc.CString("mob")) != 0 {
		if vnum_mobile(&buf2[0], ch) == 0 {
			send_to_char(ch, libc.CString("No mobiles by that name.\r\n"))
		}
	}
	if is_abbrev(&buf[0], libc.CString("obj")) != 0 {
		if vnum_object(&buf2[0], ch) == 0 {
			send_to_char(ch, libc.CString("No objects by that name.\r\n"))
		}
	}
	if is_abbrev(&buf[0], libc.CString("mat")) != 0 {
		if vnum_material(&buf2[0], ch) == 0 {
			send_to_char(ch, libc.CString("No materials by that name.\r\n"))
		}
	}
	if is_abbrev(&buf[0], libc.CString("wtype")) != 0 {
		if vnum_weapontype(&buf2[0], ch) == 0 {
			send_to_char(ch, libc.CString("No weapon types by that name.\r\n"))
		}
	}
	if is_abbrev(&buf[0], libc.CString("atype")) != 0 {
		if vnum_armortype(&buf2[0], ch) == 0 {
			send_to_char(ch, libc.CString("No armor types by that name.\r\n"))
		}
	}
}
func list_zone_commands_room(ch *char_data, rvnum room_vnum) {
	var (
		zrnum    zone_rnum = real_zone_by_thing(rvnum)
		rrnum    room_rnum = real_room(rvnum)
		cmd_room room_rnum = room_rnum(-1)
		subcmd   int       = 0
		count    int       = 0
	)
	if zrnum == zone_rnum(-1) || rrnum == room_rnum(-1) {
		send_to_char(ch, libc.CString("No zone information available.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Zone commands in this room:@y\r\n"))
	for int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Command) != 'S' {
		switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Command {
		case 'M':
			fallthrough
		case 'O':
			fallthrough
		case 'T':
			fallthrough
		case 'V':
			cmd_room = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3)
		case 'D':
			fallthrough
		case 'R':
			cmd_room = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)
		default:
		}
		if cmd_room == rrnum {
			count++
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Command {
			case 'M':
				send_to_char(ch, libc.CString("%sLoad %s@y [@c%d@y], MaxMud : %d, MaxR : %d, Chance : %d\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Short_descr, (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg4, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg5)
			case 'G':
				send_to_char(ch, libc.CString("%sGive it %s@y [@c%d@y], Max : %d, Chance : %d\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg5)
			case 'O':
				send_to_char(ch, libc.CString("%sLoad %s@y [@c%d@y], Max : %d, MaxR : %d, Chance : %d\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg4, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg5)
			case 'E':
				send_to_char(ch, libc.CString("%sEquip with %s@y [@c%d@y], %s, Max : %d, Chance : %d\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Vnum, equipment_types[(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3], (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg5)
			case 'P':
				send_to_char(ch, libc.CString("%sPut %s@y [@c%d@y] in %s@y [@c%d@y], Max : %d, Chance : %d\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1)))).Vnum, (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg5)
			case 'R':
				send_to_char(ch, libc.CString("%sRemove %s@y [@c%d@y] from room.\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2)))).Short_description, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2)))).Vnum)
			case 'D':
				send_to_char(ch, libc.CString("%sSet door %s as %s.\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), dirs[(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2], func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3 != 0 {
						if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg3 == 1 {
							return "closed"
						}
						return "locked"
					}
					return "open"
				}())
			case 'T':
				send_to_char(ch, libc.CString("%sAttach trigger @c%s@y [@c%d@y] to %s\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2)))).Proto.Name, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2)))).Vnum, func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == MOB_TRIGGER {
						return "mobile"
					}
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == OBJ_TRIGGER {
						return "object"
					}
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == WLD_TRIGGER {
						return "room"
					}
					return "????"
				}())
			case 'V':
				send_to_char(ch, libc.CString("%sAssign global %s:%d to %s = %s\r\n"), func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).If_flag {
						return " then "
					}
					return ""
				}(), (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Sarg1, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg2, func() string {
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == MOB_TRIGGER {
						return "mobile"
					}
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == OBJ_TRIGGER {
						return "object"
					}
					if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Arg1 == WLD_TRIGGER {
						return "room"
					}
					return "????"
				}(), (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(subcmd)))).Sarg2)
			default:
				send_to_char(ch, libc.CString("<Unknown Command>\r\n"))
			}
		}
		subcmd++
	}
	send_to_char(ch, libc.CString("@n"))
	if count == 0 {
		send_to_char(ch, libc.CString("None!\r\n"))
	}
}
func do_stat_room(ch *char_data) {
	var (
		buf2   [64936]byte
		desc   *extra_descr_data
		rm     *room_data = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))
		i      int
		found  int
		column int
		j      *obj_data
		k      *char_data
	)
	send_to_char(ch, libc.CString("Room name: @c%s@n\r\n"), rm.Name)
	sprinttype(rm.Sector_type, sector_types[:], &buf2[0], uint64(64936))
	send_to_char(ch, libc.CString("Zone: [%3d], VNum: [@g%5d@n], RNum: [%5d], IDNum: [%5ld], Type: %s\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rm.Zone)))).Number, rm.Number, ch.In_room, int(rm.Number)+ROOM_ID_BASE, &buf2[0])
	sprintbitarray(rm.Room_flags[:], room_bits[:], RF_ARRAY_MAX, &buf2[0])
	send_to_char(ch, libc.CString("Room Damage: %d, Room Effect: %d\r\n"), rm.Dmg, rm.Geffect)
	send_to_char(ch, libc.CString("SpecProc: %s, Flags: %s\r\n"), func() string {
		if rm.Func == nil {
			return "None"
		}
		return "Exists"
	}(), &buf2[0])
	send_to_char(ch, libc.CString("Description:\r\n%s"), func() *byte {
		if rm.Description != nil {
			return rm.Description
		}
		return libc.CString("  None.\r\n")
	}())
	if rm.Ex_description != nil {
		send_to_char(ch, libc.CString("Extra descs:"))
		for desc = rm.Ex_description; desc != nil; desc = desc.Next {
			send_to_char(ch, libc.CString(" [@c%s@n]"), desc.Keyword)
		}
		send_to_char(ch, libc.CString("\r\n"))
	}
	send_to_char(ch, libc.CString("Chars present:"))
	column = 14
	for func() *char_data {
		found = FALSE
		return func() *char_data {
			k = rm.People
			return k
		}()
	}(); k != nil; k = k.Next_in_room {
		if !CAN_SEE(ch, k) {
			continue
		}
		column += int(send_to_char(ch, libc.CString("%s @y%s@n(%s)"), func() string {
			if func() int {
				p := &found
				x := *p
				*p++
				return x
			}() != 0 {
				return ","
			}
			return ""
		}(), GET_NAME(k), func() string {
			if !IS_NPC(k) {
				return "PC"
			}
			if !IS_MOB(k) {
				return "NPC"
			}
			return "MOB"
		}()))
		if column >= 62 {
			send_to_char(ch, libc.CString("%s\r\n"), func() string {
				if k.Next_in_room != nil {
					return ","
				}
				return ""
			}())
			found = FALSE
			column = 0
		}
	}
	if rm.Contents != nil {
		send_to_char(ch, libc.CString("Contents:@g"))
		column = 9
		for func() *obj_data {
			found = 0
			return func() *obj_data {
				j = rm.Contents
				return j
			}()
		}(); j != nil; j = j.Next_content {
			if !CAN_SEE_OBJ(ch, j) {
				continue
			}
			column += int(send_to_char(ch, libc.CString("%s %s"), func() string {
				if func() int {
					p := &found
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), j.Short_description))
			if column >= 62 {
				send_to_char(ch, libc.CString("%s\r\n"), func() string {
					if j.Next_content != nil {
						return ","
					}
					return ""
				}())
				found = FALSE
				column = 0
			}
		}
		send_to_char(ch, libc.CString("@n"))
	}
	for i = 0; i < NUM_OF_DIRS; i++ {
		var buf1 [128]byte
		if rm.Dir_option[i] == nil {
			continue
		}
		if rm.Dir_option[i].To_room == room_rnum(-1) {
			stdio.Snprintf(&buf1[0], int(128), " @cNONE@n")
		} else {
			stdio.Snprintf(&buf1[0], int(128), "@c%5d@n", func() room_vnum {
				if rm.Dir_option[i].To_room != room_rnum(-1) && rm.Dir_option[i].To_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rm.Dir_option[i].To_room)))).Number
				}
				return -1
			}())
		}
		sprintbit(rm.Dir_option[i].Exit_info, exit_bits[:], &buf2[0], uint64(64936))
		send_to_char(ch, libc.CString("Exit @c%-5s@n:  To: [%s], Key: [%5d], Keywrd: %s, Type: %s\r\n  DC Lock: [%2d], DC Hide: [%2d], DC Skill: [%4s], DC Move: [%2d]\r\n%s"), dirs[i], &buf1[0], func() obj_vnum {
			if rm.Dir_option[i].Key == obj_vnum(-1) {
				return -1
			}
			return rm.Dir_option[i].Key
		}(), func() *byte {
			if rm.Dir_option[i].Keyword != nil {
				return rm.Dir_option[i].Keyword
			}
			return libc.CString("None")
		}(), &buf2[0], rm.Dir_option[i].Dclock, rm.Dir_option[i].Dchide, func() string {
			if rm.Dir_option[i].Dcskill == 0 {
				return "None"
			}
			return libc.GoString(spell_info[rm.Dir_option[i].Dcskill].Name)
		}(), rm.Dir_option[i].Dcmove, func() *byte {
			if rm.Dir_option[i].General_description != nil {
				return rm.Dir_option[i].General_description
			}
			return libc.CString("  No exit description.\r\n")
		}())
	}
	do_sstat_room(ch, rm)
	list_zone_commands_room(ch, rm.Number)
}
func do_stat_object(ch *char_data, j *obj_data) {
	var (
		i      int
		found  int
		vnum   obj_vnum
		j2     *obj_data
		sitter *char_data
		desc   *extra_descr_data
		buf    [64936]byte
	)
	vnum = GET_OBJ_VNUM(j)
	if j.Lload > 0 {
		var tmstr *byte
		tmstr = libc.AscTime(libc.LocalTime(&j.Lload))
		*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
		send_to_char(ch, libc.CString("LOADED DROPPED: [%s]\r\n"), tmstr)
	}
	if GET_OBJ_VNUM(j) == 65 {
		send_to_char(ch, libc.CString("Healing Tank Charge Level: [%d]\r\n"), j.Healcharge)
	}
	send_to_char(ch, libc.CString("Name: '%s', Keywords: %s, Size: %s\r\n"), func() *byte {
		if j.Short_description != nil {
			return j.Short_description
		}
		return libc.CString("<None>")
	}(), j.Name, size_names[j.Size])
	sprinttype(int(j.Type_flag), item_types[:], &buf[0], uint64(64936))
	send_to_char(ch, libc.CString("VNum: [@g%5d@n], RNum: [%5d], Idnum: [%5d], Type: %s, SpecProc: %s\r\n"), vnum, j.Item_number, j.Id, &buf[0], func() string {
		if GET_OBJ_SPEC(j) != nil {
			return "Exists"
		}
		return "None"
	}())
	send_to_char(ch, libc.CString("Generation time: @g%s@nUnique ID: @g%lld@n\r\n"), ctime(&j.Generation), j.Unique_id)
	send_to_char(ch, libc.CString("Object Hit Points: [ @g%3d@n/@g%3d@n]\r\n"), j.Value[VAL_ALL_HEALTH], j.Value[VAL_ALL_MAXHEALTH])
	send_to_char(ch, libc.CString("Object loaded in room: @y%d@n\r\n"), j.Room_loaded)
	send_to_char(ch, libc.CString("Object Material: @y%s@n\r\n"), material_names[j.Value[VAL_ALL_MATERIAL]])
	if j.Sitting != nil {
		sitter = j.Sitting
		send_to_char(ch, libc.CString("HOLDING: %s\r\n"), GET_NAME(sitter))
	}
	if j.Ex_description != nil {
		send_to_char(ch, libc.CString("Extra descs:"))
		for desc = j.Ex_description; desc != nil; desc = desc.Next {
			send_to_char(ch, libc.CString(" [@g%s@n]"), desc.Keyword)
		}
		send_to_char(ch, libc.CString("\r\n"))
	}
	sprintbitarray(j.Wear_flags[:], wear_bits[:], TW_ARRAY_MAX, &buf[0])
	send_to_char(ch, libc.CString("Can be worn on: %s\r\n"), &buf[0])
	sprintbitarray(j.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &buf[0])
	send_to_char(ch, libc.CString("Set char bits : %s\r\n"), &buf[0])
	sprintbitarray(j.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, &buf[0])
	send_to_char(ch, libc.CString("Extra flags   : %s\r\n"), &buf[0])
	send_to_char(ch, libc.CString("Weight: %lld, Value: %d, Cost/day: %d, Timer: %d, Min Level: %d\r\n"), j.Weight, j.Cost, j.Cost_per_day, j.Timer, j.Level)
	send_to_char(ch, libc.CString("In room: %d (%s), "), func() room_vnum {
		if j.In_room != room_rnum(-1) && j.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(j.In_room)))).Number
		}
		return -1
	}(), func() string {
		if j.In_room == room_rnum(-1) {
			return "Nowhere"
		}
		return libc.GoString((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(j.In_room)))).Name)
	}())
	send_to_char(ch, libc.CString("In object: %s, "), func() *byte {
		if j.In_obj != nil {
			return j.In_obj.Short_description
		}
		return libc.CString("None")
	}())
	send_to_char(ch, libc.CString("Carried by: %s, "), func() *byte {
		if j.Carried_by != nil {
			return GET_NAME(j.Carried_by)
		}
		return libc.CString("Nobody")
	}())
	send_to_char(ch, libc.CString("Worn by: %s\r\n"), func() *byte {
		if j.Worn_by != nil {
			return GET_NAME(j.Worn_by)
		}
		return libc.CString("Nobody")
	}())
	switch j.Type_flag {
	case ITEM_LIGHT:
		if (j.Value[VAL_LIGHT_HOURS]) == -1 {
			send_to_char(ch, libc.CString("Hours left: Infinite\r\n"))
		} else {
			send_to_char(ch, libc.CString("Hours left: [%d]\r\n"), j.Value[VAL_LIGHT_HOURS])
		}
	case ITEM_SCROLL:
		send_to_char(ch, libc.CString("Spell: (Level %d) %s\r\n"), j.Value[VAL_SCROLL_LEVEL], skill_name(j.Value[VAL_SCROLL_SPELL1]))
	case ITEM_POTION:
		send_to_char(ch, libc.CString("Spells: (Level %d) %s, %s, %s\r\n"), j.Value[VAL_POTION_LEVEL], skill_name(j.Value[VAL_POTION_SPELL1]), skill_name(j.Value[VAL_POTION_SPELL2]), skill_name(j.Value[VAL_POTION_SPELL3]))
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		send_to_char(ch, libc.CString("Spell: %s at level %d, %d (of %d) charges remaining\r\n"), skill_name(j.Value[VAL_STAFF_SPELL]), j.Value[VAL_STAFF_LEVEL], j.Value[VAL_STAFF_CHARGES], j.Value[VAL_STAFF_MAXCHARGES])
	case ITEM_WEAPON:
		send_to_char(ch, libc.CString("Weapon Type: %s, Todam: %dd%d, Message type: %d\r\n"), weapon_type[j.Value[VAL_WEAPON_SKILL]], j.Value[VAL_WEAPON_DAMDICE], j.Value[VAL_WEAPON_DAMSIZE], j.Value[VAL_WEAPON_DAMTYPE])
		send_to_char(ch, libc.CString("Average damage per round %.1f\r\n"), (float64((j.Value[VAL_WEAPON_DAMSIZE])+1)/2.0)*float64(j.Value[VAL_WEAPON_DAMDICE]))
		send_to_char(ch, libc.CString("Crit type: %s, Crit range: %d-20\r\n"), crit_type[j.Value[6]], 20-(j.Value[8]))
	case ITEM_ARMOR:
		send_to_char(ch, libc.CString("Armor Type: %s, AC-apply: [%d]\r\n"), armor_type[j.Value[VAL_ARMOR_SKILL]], j.Value[VAL_ARMOR_APPLYAC])
		send_to_char(ch, libc.CString("Max dex bonus: %d, Armor penalty: %d, Spell failure: %d\r\n"), j.Value[VAL_ARMOR_MAXDEXMOD], j.Value[VAL_ARMOR_CHECK], j.Value[VAL_ARMOR_SPELLFAIL])
	case ITEM_TRAP:
		send_to_char(ch, libc.CString("Spell: %d, - Hitpoints: %d\r\n"), j.Value[VAL_TRAP_SPELL], j.Value[VAL_TRAP_HITPOINTS])
	case ITEM_CONTAINER:
		sprintbit(bitvector_t(int32(j.Value[VAL_CONTAINER_FLAGS])), container_bits[:], &buf[0], uint64(64936))
		send_to_char(ch, libc.CString("Weight capacity: %d, Lock Type: %s, Key Num: %d, Corpse: %s\r\n"), j.Value[VAL_CONTAINER_CAPACITY], &buf[0], j.Value[VAL_CONTAINER_KEY], func() string {
			if (j.Value[VAL_CONTAINER_CORPSE]) != 0 {
				return "YES"
			}
			return "NO"
		}())
	case ITEM_DRINKCON:
		fallthrough
	case ITEM_FOUNTAIN:
		sprinttype(j.Value[VAL_DRINKCON_LIQUID], drinks[:], &buf[0], uint64(64936))
		send_to_char(ch, libc.CString("Capacity: %d, Contains: %d, Poisoned: %s, Liquid: %s\r\n"), j.Value[VAL_DRINKCON_CAPACITY], j.Value[VAL_DRINKCON_HOWFULL], func() string {
			if (j.Value[VAL_DRINKCON_POISON]) != 0 {
				return "YES"
			}
			return "NO"
		}(), &buf[0])
	case ITEM_NOTE:
		send_to_char(ch, libc.CString("Tongue: %d\r\n"), j.Value[VAL_NOTE_LANGUAGE])
	case ITEM_KEY:
	case ITEM_FOOD:
		send_to_char(ch, libc.CString("Makes full: %d, Poisoned: %s\r\n"), j.Value[VAL_FOOD_FOODVAL], func() string {
			if (j.Value[VAL_FOOD_POISON]) != 0 {
				return "YES"
			}
			return "NO"
		}())
	case ITEM_MONEY:
		send_to_char(ch, libc.CString("Coins: %d\r\n"), j.Value[VAL_MONEY_SIZE])
	default:
		send_to_char(ch, libc.CString("Values 0-12: [%d] [%d] [%d] [%d] [%d] [%d] [%d] [%d] [%d] [%d] [%d] [%d]\r\n"), j.Value[0], j.Value[1], j.Value[2], j.Value[3], j.Value[4], j.Value[5], j.Value[6], j.Value[7], j.Value[8], j.Value[9], j.Value[10], j.Value[11])
	}
	if j.Contains != nil {
		var column int
		send_to_char(ch, libc.CString("\r\nContents:@g"))
		column = 9
		for func() *obj_data {
			found = 0
			return func() *obj_data {
				j2 = j.Contains
				return j2
			}()
		}(); j2 != nil; j2 = j2.Next_content {
			column += int(send_to_char(ch, libc.CString("%s %s"), func() string {
				if func() int {
					p := &found
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), j2.Short_description))
			if column >= 62 {
				send_to_char(ch, libc.CString("%s\r\n"), func() string {
					if j2.Next_content != nil {
						return ","
					}
					return ""
				}())
				found = FALSE
				column = 0
			}
		}
		send_to_char(ch, libc.CString("@n"))
	}
	found = FALSE
	send_to_char(ch, libc.CString("Affections:"))
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if j.Affected[i].Modifier != 0 {
			sprinttype(j.Affected[i].Location, apply_types[:], &buf[0], uint64(64936))
			send_to_char(ch, libc.CString("%s %+d to %s"), func() string {
				if func() int {
					p := &found
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), j.Affected[i].Modifier, &buf[0])
			switch j.Affected[i].Location {
			case APPLY_FEAT:
				send_to_char(ch, libc.CString(" (%s)"), feat_list[j.Affected[i].Specific].Name)
			case APPLY_SKILL:
				send_to_char(ch, libc.CString(" (%s)"), spell_info[j.Affected[i].Specific].Name)
			}
		}
	}
	if found == 0 {
		send_to_char(ch, libc.CString(" None"))
	}
	send_to_char(ch, libc.CString("\r\n"))
	do_sstat_object(ch, j)
}
func do_stat_character(ch *char_data, k *char_data) {
	var (
		buf    [64936]byte
		buf2   [64936]byte
		i      int
		i2     int
		column int
		found  int = FALSE
		j      *obj_data
		chair  *obj_data
		fol    *follow_type
		aff    *affected_type
	)
	if IS_NPC(k) {
		var tmstr *byte
		tmstr = libc.AscTime(libc.LocalTime(&k.Lastpl))
		*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
		send_to_char(ch, libc.CString("LOADED AT: [%s]\r\n"), tmstr)
	}
	sprinttype(int(k.Sex), genders[:], &buf[0], uint64(64936))
	send_to_char(ch, libc.CString("%s %s '%s'  IDNum: [%5d], In room [%5d], Loadroom : [%5d]\r\n"), &buf[0], func() string {
		if !IS_NPC(k) {
			return "PC"
		}
		if !IS_MOB(k) {
			return "NPC"
		}
		return "MOB"
	}(), GET_NAME(k), func() int {
		if IS_NPC(k) {
			return int(k.Id)
		}
		return int(k.Idnum)
	}(), func() room_vnum {
		if k.In_room != room_rnum(-1) && k.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Number
		}
		return -1
	}(), func() room_vnum {
		if IS_NPC(k) {
			return k.Hometown
		}
		return k.Player_specials.Load_room
	}())
	send_to_char(ch, libc.CString("DROOM: [%5d]\r\n"), k.Droom)
	if IS_MOB(k) {
		if int(k.Master_id) > -1 {
			stdio.Sprintf(&buf[0], ", Master: %s", get_name_by_id(int(k.Master_id)))
		} else {
			buf[0] = 0
		}
		send_to_char(ch, libc.CString("Keyword: %s, VNum: [%5d], RNum: [%5d]%s\r\n"), k.Name, GET_MOB_VNUM(k), k.Nr, &buf[0])
	} else {
		send_to_char(ch, libc.CString("Title: %s\r\n"), func() *byte {
			if k.Title != nil {
				return k.Title
			}
			return libc.CString("<None>")
		}())
	}
	send_to_char(ch, libc.CString("L-Des: %s@n"), func() *byte {
		if k.Long_descr != nil {
			return k.Long_descr
		}
		return libc.CString("<None>\r\n")
	}())
	if config_info.Advance.Allow_multiclass != 0 {
		libc.StrNCpy(&buf[0], class_desc_str(k, 1, 0), int(64936))
	} else {
		sprinttype(int(k.Chclass), pc_class_types[:], &buf[0], uint64(64936))
	}
	sprinttype(int(k.Race), pc_race_types[:], &buf2[0], uint64(64936))
	send_to_char(ch, libc.CString("Class: %s, Race: %s, Lev: [@y%2d(%dHD+%dcl+%d)@n], XP: [@y%lld@n]\r\n"), &buf[0], &buf2[0], GET_LEVEL(k), k.Race_level, k.Level, k.Level_adj, k.Exp)
	if !IS_NPC(k) {
		var (
			buf1   [64]byte
			cmbuf2 [64]byte
		)
		strlcpy(&buf1[0], libc.AscTime(libc.LocalTime(&k.Time.Created)), uint64(64))
		strlcpy(&cmbuf2[0], libc.AscTime(libc.LocalTime(&k.Time.Logon)), uint64(64))
		buf1[10] = func() byte {
			p := &cmbuf2[10]
			cmbuf2[10] = '\x00'
			return *p
		}()
		send_to_char(ch, libc.CString("Created: [%s], Last Logon: [%s], Played [%dh %dm], Age [%d]\r\n"), &buf1[0], &cmbuf2[0], int(k.Time.Played)/3600, int((k.Time.Played%3600)/60), age(k).Year)
		if k.Desc != nil {
			send_to_char(ch, libc.CString("@YOwned by User@D: [@C%s@D]@n\r\n"), GET_USER(k))
		} else {
			send_to_char(ch, libc.CString("@YOwned by User@D: [@C%s@D]@n\r\n"), k.Loguser)
		}
		if !IS_NPC(k) {
			send_to_char(ch, libc.CString("@RCharacter Deaths@D: @r%d@n\r\n"), k.Dcount)
		}
		send_to_char(ch, libc.CString("Hometown: [%d], Align: [%4d], Ethic: [%4d]"), k.Hometown, k.Alignment, k.Alignment_ethic)
		if k.Admlevel >= ADMLVL_BUILDER {
			if k.Player_specials.Olc_zone == AEDIT_PERMISSION {
				send_to_char(ch, libc.CString(", OLC[@cActions@n]"))
			} else if k.Player_specials.Olc_zone == HEDIT_PERMISSION {
				send_to_char(ch, libc.CString(", OLC[@cHedit@n]"))
			} else if k.Player_specials.Olc_zone == int(-1) {
				send_to_char(ch, libc.CString(", OLC[@cOFF@n]"))
			} else {
				send_to_char(ch, libc.CString(", OLC: [@c%d@n]"), k.Player_specials.Olc_zone)
			}
		}
		send_to_char(ch, libc.CString("\r\n"))
	}
	send_to_char(ch, libc.CString("Str: [@c%d@n]  Int: [@c%d@n]  Wis: [@c%d@n]  Dex: [@c%d@n]  Con: [@c%d@n]  Cha: [@c%d@n]\r\n"), k.Aff_abils.Str, k.Aff_abils.Intel, k.Aff_abils.Wis, k.Aff_abils.Dex, k.Aff_abils.Con, k.Aff_abils.Cha)
	send_to_char(ch, libc.CString("PL :[@g%12s@n]  KI :[@g%12s@n]  ST :[@g%12s@n]\r\n"), add_commas(k.Hit), add_commas(k.Mana), add_commas(k.Move))
	send_to_char(ch, libc.CString("MPL:[@g%12s@n]  MKI:[@g%12s@n]  MST:[@g%12s@n]\r\n"), add_commas(k.Max_hit), add_commas(k.Max_mana), add_commas(k.Max_move))
	send_to_char(ch, libc.CString("BPL:[@g%12s@n]  BKI:[@g%12s@n]  BST:[@g%12s@n]\r\n"), add_commas(k.Basepl), add_commas(k.Baseki), add_commas(k.Basest))
	send_to_char(ch, libc.CString("LF :[@g%12s@n]  MLF:[@g%12s@n]  LFP:[@g%3d@n]\r\n"), add_commas(k.Lifeforce), add_commas(int64(GET_LIFEMAX(k))), k.Lifeperc)
	if k.Admlevel != 0 {
		send_to_char(ch, libc.CString("Admin Level: [@y%d - %s@n]\r\n"), k.Admlevel, admin_level_names[k.Admlevel])
	}
	send_to_char(ch, libc.CString("Coins: [%9d], Bank: [%9d] (Total: %d)\r\n"), k.Gold, k.Bank_gold, k.Gold+k.Bank_gold)
	send_to_char(ch, libc.CString("Armor: [%d ], Damage: [%2d], Saving throws: [%d/%d/%d]\r\n"), k.Armor, k.Damage_mod, k.Apply_saving_throw[0], k.Apply_saving_throw[1], k.Apply_saving_throw[2])
	sprinttype(int(k.Position), position_types[:], &buf[0], uint64(64936))
	send_to_char(ch, libc.CString("Pos: %s, Fighting: %s"), &buf[0], func() *byte {
		if k.Fighting != nil {
			return GET_NAME(k.Fighting)
		}
		return libc.CString("Nobody")
	}())
	if k.Desc != nil {
		sprinttype(k.Desc.Connected, connected_types[:], &buf[0], uint64(64936))
		send_to_char(ch, libc.CString(", Connected: %s"), &buf[0])
	}
	if IS_NPC(k) {
		sprinttype(int(k.Mob_specials.Default_pos), position_types[:], &buf[0], uint64(64936))
		send_to_char(ch, libc.CString(", Default position: %s\r\n"), &buf[0])
		sprintbitarray(k.Act[:], action_bits[:], PM_ARRAY_MAX, &buf[0])
		send_to_char(ch, libc.CString("NPC flags: @c%s@n\r\n"), &buf[0])
	} else {
		send_to_char(ch, libc.CString(", Idle Timer (in tics) [%d]\r\n"), k.Timer)
		sprintbitarray(k.Act[:], player_bits[:], PM_ARRAY_MAX, &buf[0])
		send_to_char(ch, libc.CString("PLR: @c%s@n\r\n"), &buf[0])
		sprintbitarray(k.Player_specials.Pref[:], preference_bits[:], PR_ARRAY_MAX, &buf[0])
		send_to_char(ch, libc.CString("PRF: @g%s@n\r\n"), &buf[0])
	}
	if IS_MOB(k) {
		send_to_char(ch, libc.CString("Mob Spec-Proc: %s, NPC Bare Hand Dam: %dd%d\r\n"), func() string {
			if (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(k.Nr)))).Func != nil {
				return "Exists"
			}
			return "None"
		}(), k.Mob_specials.Damnodice, k.Mob_specials.Damsizedice)
		send_to_char(ch, libc.CString("Average damage per round %.1f (%.1f [BHD] + %d [STR MOD] + %d [DMG MOD])\r\n"), (float64(int(k.Mob_specials.Damsizedice)+1)/2.0)*float64(k.Mob_specials.Damnodice)+float64(ability_mod_value(int(k.Aff_abils.Str)))+float64(k.Damage_mod), (float64(int(k.Mob_specials.Damsizedice)+1)/2.0)*float64(k.Mob_specials.Damnodice), ability_mod_value(int(k.Aff_abils.Str)), k.Damage_mod)
	}
	var counts int = 0
	var total int = 0
	for func() *obj_data {
		i = 0
		return func() *obj_data {
			j = k.Carrying
			return j
		}()
	}(); j != nil; func() int {
		j = j.Next_content
		return func() int {
			p := &i
			x := *p
			*p++
			return x
		}()
	}() {
		counts += check_insidebag(j, 0.5)
		counts++
	}
	total = counts
	total += i
	for func() int {
		i = 0
		return func() int {
			i2 = 0
			return i2
		}()
	}(); i < NUM_WEARS; i++ {
		if (k.Equipment[i]) != nil {
			i2++
			total += check_insidebag(k.Equipment[i], 0.5) + 1
		}
	}
	send_to_char(ch, libc.CString("Carried: weight: %d, Total Items (includes bagged items): %d, EQ: %d\r\n"), k.Carry_weight, total, i2)
	if !IS_NPC(k) {
		send_to_char(ch, libc.CString("Hunger: %d, Thirst: %d, Drunk: %d\r\n"), k.Player_specials.Conditions[HUNGER], k.Player_specials.Conditions[THIRST], k.Player_specials.Conditions[DRUNK])
	}
	column = int(send_to_char(ch, libc.CString("Master is: %s, Followers are:"), func() *byte {
		if k.Master != nil {
			return GET_NAME(k.Master)
		}
		return libc.CString("<none>")
	}()))
	if k.Followers == nil {
		send_to_char(ch, libc.CString(" <none>\r\n"))
	} else {
		for fol = k.Followers; fol != nil; fol = fol.Next {
			column += int(send_to_char(ch, libc.CString("%s %s"), func() string {
				if func() int {
					p := &found
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), PERS(fol.Follower, ch)))
			if column >= 62 {
				send_to_char(ch, libc.CString("%s\r\n"), func() string {
					if fol.Next != nil {
						return ","
					}
					return ""
				}())
				found = FALSE
				column = 0
			}
		}
		if column != 0 {
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if k.Sits != nil {
		chair = k.Sits
		send_to_char(ch, libc.CString("Is on: %s@n\r\n"), chair.Short_description)
	}
	sprintbitarray(k.Affected_by[:], affected_bits[:], AF_ARRAY_MAX, &buf[0])
	send_to_char(ch, libc.CString("AFF: @y%s@n\r\n"), &buf[0])
	if k.Affected != nil {
		for aff = k.Affected; aff != nil; aff = aff.Next {
			send_to_char(ch, libc.CString("SPL: (%3dhr) @c%-21s@n "), int(aff.Duration)+1, skill_name(int(aff.Type)))
			if aff.Modifier != 0 {
				send_to_char(ch, libc.CString("%+d to %s"), aff.Modifier, apply_types[aff.Location])
			}
			if aff.Bitvector != 0 {
				if aff.Modifier != 0 {
					send_to_char(ch, libc.CString(", "))
				}
				libc.StrCpy(&buf[0], affected_bits[aff.Bitvector])
				send_to_char(ch, libc.CString("sets %s"), &buf[0])
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if k.Affectedv != nil {
		for aff = k.Affectedv; aff != nil; aff = aff.Next {
			send_to_char(ch, libc.CString("SPL: (%3d rounds) @c%-21s@n "), int(aff.Duration)+1, skill_name(int(aff.Type)))
			if aff.Modifier != 0 {
				send_to_char(ch, libc.CString("%+d to %s"), aff.Modifier, apply_types[aff.Location])
			}
			if aff.Bitvector != 0 {
				if aff.Modifier != 0 {
					send_to_char(ch, libc.CString(", "))
				}
				libc.StrCpy(&buf[0], affected_bits[aff.Bitvector])
				send_to_char(ch, libc.CString("sets %s"), &buf[0])
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if IS_NPC(k) {
		do_sstat_character(ch, k)
		if k.Memory != nil {
			var mem *script_memory = k.Memory
			send_to_char(ch, libc.CString("Script memory:\r\n  Remember             Command\r\n"))
			for mem != nil {
				var mc *char_data = find_char(mem.Id)
				if mc == nil {
					send_to_char(ch, libc.CString("  ** Corrupted!\r\n"))
				} else {
					if mem.Cmd != nil {
						send_to_char(ch, libc.CString("  %-20.20s%s\r\n"), GET_NAME(mc), mem.Cmd)
					} else {
						send_to_char(ch, libc.CString("  %-20.20s <default>\r\n"), GET_NAME(mc))
					}
				}
				mem = mem.Next
			}
		}
	} else {
		var (
			x     int
			track int = 0
		)
		send_to_char(ch, libc.CString("Bonuses/Negatives:\r\n"))
		for x = 0; x < 30; x++ {
			if x < 15 {
				if (k.Bonuses[x]) > 0 {
					send_to_char(ch, libc.CString("@c%s@n\n"), list_bonus[x])
					track += 1
				}
			} else {
				if (k.Bonuses[x]) > 0 {
					send_to_char(ch, libc.CString("@r%s@n\n"), list_bonus[x])
					track += 1
				}
			}
		}
		if track <= 0 {
			send_to_char(ch, libc.CString("@wNone.@n\r\n"))
		}
		send_to_char(ch, libc.CString("To see player variables use varstat now.\r\n"))
	}
}
func do_varstat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That player is not in the game.\r\n"))
		return
	} else if IS_NPC(vict) {
		send_to_char(ch, libc.CString("Just use stat for an NPC\r\n"))
		return
	} else {
		if vict.Script != nil && vict.Script.Global_vars != nil {
			var (
				tv    *trig_var_data
				uname [2048]byte
			)
			send_to_char(ch, libc.CString("%s's Global Variables:\r\n"), GET_NAME(vict))
			for tv = vict.Script.Global_vars; tv != nil; tv = tv.Next {
				if *tv.Value == UID_CHAR {
					find_uid_name(tv.Value, &uname[0], uint64(2048))
					send_to_char(ch, libc.CString("    %10s:  [UID]: %s\r\n"), tv.Name, &uname[0])
				} else {
					send_to_char(ch, libc.CString("    %10s:  %s\r\n"), tv.Name, tv.Value)
				}
			}
		}
	}
}
func do_stat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf1   [2048]byte
		buf2   [2048]byte
		victim *char_data
		object *obj_data
	)
	half_chop(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		send_to_char(ch, libc.CString("Stats on who or what or where?\r\n"))
		return
	} else if is_abbrev(&buf1[0], libc.CString("room")) != 0 {
		do_stat_room(ch)
	} else if is_abbrev(&buf1[0], libc.CString("mob")) != 0 {
		if buf2[0] == 0 {
			send_to_char(ch, libc.CString("Stats on which mobile?\r\n"))
		} else {
			if (func() *char_data {
				victim = get_char_vis(ch, &buf2[0], nil, 1<<1)
				return victim
			}()) != nil {
				do_stat_character(ch, victim)
			} else {
				send_to_char(ch, libc.CString("No such mobile around.\r\n"))
			}
		}
	} else if is_abbrev(&buf1[0], libc.CString("player")) != 0 {
		if buf2[0] == 0 {
			send_to_char(ch, libc.CString("Stats on which player?\r\n"))
		} else {
			if (func() *char_data {
				victim = get_player_vis(ch, &buf2[0], nil, 1<<1)
				return victim
			}()) != nil {
				do_stat_character(ch, victim)
			} else {
				send_to_char(ch, libc.CString("No such player around.\r\n"))
			}
		}
	} else if is_abbrev(&buf1[0], libc.CString("file")) != 0 {
		if buf2[0] == 0 {
			send_to_char(ch, libc.CString("Stats on which player?\r\n"))
		} else if (func() *char_data {
			victim = get_player_vis(ch, &buf2[0], nil, 1<<1)
			return victim
		}()) != nil {
			do_stat_character(ch, victim)
		} else {
			victim = new(char_data)
			clear_char(victim)
			victim.Player_specials = new(player_special_data)
			if load_char(&buf2[0], victim) >= 0 {
				char_to_room(victim, 0)
				if victim.Admlevel > ch.Admlevel {
					send_to_char(ch, libc.CString("Sorry, you can't do that.\r\n"))
				} else {
					do_stat_character(ch, victim)
				}
				extract_char_final(victim)
			} else {
				send_to_char(ch, libc.CString("There is no such player.\r\n"))
				free_char(victim)
			}
		}
	} else if is_abbrev(&buf1[0], libc.CString("object")) != 0 {
		if buf2[0] == 0 {
			send_to_char(ch, libc.CString("Stats on which object?\r\n"))
		} else {
			if (func() *obj_data {
				object = get_obj_vis(ch, &buf2[0], nil)
				return object
			}()) != nil {
				do_stat_object(ch, object)
			} else {
				send_to_char(ch, libc.CString("No such object around.\r\n"))
			}
		}
	} else if is_abbrev(&buf1[0], libc.CString("zone")) != 0 {
		if buf2[0] == 0 {
			send_to_char(ch, libc.CString("Stats on which zone?\r\n"))
			return
		} else {
			print_zone(ch, zone_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
			return
		}
	} else {
		var (
			name   *byte = &buf1[0]
			number int   = get_number(&name)
		)
		if (func() *obj_data {
			object = get_obj_in_equip_vis(ch, name, &number, ch.Equipment[:])
			return object
		}()) != nil {
			do_stat_object(ch, object)
		} else if (func() *obj_data {
			object = get_obj_in_list_vis(ch, name, &number, ch.Carrying)
			return object
		}()) != nil {
			do_stat_object(ch, object)
		} else if (func() *char_data {
			victim = get_char_vis(ch, name, &number, 1<<0)
			return victim
		}()) != nil {
			do_stat_character(ch, victim)
		} else if (func() *obj_data {
			object = get_obj_in_list_vis(ch, name, &number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return object
		}()) != nil {
			do_stat_object(ch, object)
		} else if (func() *char_data {
			victim = get_char_vis(ch, name, &number, 1<<1)
			return victim
		}()) != nil {
			do_stat_character(ch, victim)
		} else if (func() *obj_data {
			object = get_obj_vis(ch, name, &number)
			return object
		}()) != nil {
			do_stat_object(ch, object)
		} else {
			send_to_char(ch, libc.CString("Nothing around by that name.\r\n"))
		}
	}
}
func do_shutdown(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	if subcmd != SCMD_SHUTDOWN {
		send_to_char(ch, libc.CString("If you want to shut something down, say so!\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		basic_mud_log(libc.CString("(GC) Shutdown by %s."), GET_NAME(ch))
		send_to_all(libc.CString("Shutting down.\r\n"))
		circle_shutdown = 1
	} else if libc.StrCaseCmp(&arg[0], libc.CString("reboot")) == 0 {
		basic_mud_log(libc.CString("(GC) Reboot by %s."), GET_NAME(ch))
		send_to_all(libc.CString("Rebooting.. come back in a minute or two.\r\n"))
		touch(libc.CString(FASTBOOT_FILE))
		circle_shutdown = func() int {
			circle_reboot = 1
			return circle_reboot
		}()
	} else if libc.StrCaseCmp(&arg[0], libc.CString("die")) == 0 {
		basic_mud_log(libc.CString("(GC) Shutdown by %s."), GET_NAME(ch))
		send_to_all(libc.CString("Shutting down for maintenance.\r\n"))
		touch(libc.CString(KILLSCRIPT_FILE))
		circle_shutdown = 1
	} else if libc.StrCaseCmp(&arg[0], libc.CString("now")) == 0 {
		basic_mud_log(libc.CString("(GC) Shutdown NOW by %s."), GET_NAME(ch))
		send_to_all(libc.CString("Rebooting.. come back in a minute or two.\r\n"))
		circle_shutdown = 1
		circle_reboot = 2
	} else if libc.StrCaseCmp(&arg[0], libc.CString("pause")) == 0 {
		basic_mud_log(libc.CString("(GC) Shutdown by %s."), GET_NAME(ch))
		send_to_all(libc.CString("Shutting down for maintenance.\r\n"))
		touch(libc.CString(PAUSE_FILE))
		circle_shutdown = 1
	} else {
		send_to_char(ch, libc.CString("Unknown shutdown option.\r\n"))
	}
}
func snoop_check(ch *char_data) {
	if ch == nil || ch.Desc == nil {
		return
	}
	if ch.Desc.Snooping != nil && ch.Desc.Snooping.Character.Admlevel >= ch.Admlevel {
		ch.Desc.Snooping.Snoop_by = nil
		ch.Desc.Snooping = nil
	}
	if ch.Desc.Snoop_by != nil && ch.Admlevel >= ch.Desc.Snoop_by.Character.Admlevel {
		ch.Desc.Snoop_by.Snooping = nil
		ch.Desc.Snoop_by = nil
	}
}
func stop_snooping(ch *char_data) {
	if ch.Desc.Snooping == nil {
		send_to_char(ch, libc.CString("You aren't snooping anyone.\r\n"))
	} else {
		send_to_char(ch, libc.CString("You stop snooping.\r\n"))
		ch.Desc.Snooping.Snoop_by = nil
		ch.Desc.Snooping = nil
	}
}
func do_snoop(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
		tch    *char_data
	)
	if ch.Desc == nil {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		stop_snooping(ch)
	} else if (func() *char_data {
		victim = get_char_vis(ch, &arg[0], nil, 1<<1)
		return victim
	}()) == nil {
		send_to_char(ch, libc.CString("No such person around.\r\n"))
	} else if victim.Desc == nil {
		send_to_char(ch, libc.CString("There's no link.. nothing to snoop.\r\n"))
	} else if victim == ch {
		stop_snooping(ch)
	} else if victim.Desc.Snoop_by != nil {
		send_to_char(ch, libc.CString("Busy already. \r\n"))
	} else if victim.Desc.Snooping == ch.Desc {
		send_to_char(ch, libc.CString("Don't be stupid.\r\n"))
	} else {
		if victim.Desc.Original != nil {
			tch = victim.Desc.Original
		} else {
			tch = victim
		}
		if tch.Admlevel >= ch.Admlevel {
			send_to_char(ch, libc.CString("You can't.\r\n"))
			return
		}
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		if ch.Desc.Snooping != nil {
			ch.Desc.Snooping.Snoop_by = nil
		}
		ch.Desc.Snooping = victim.Desc
		victim.Desc.Snoop_by = ch.Desc
	}
}
func do_switch(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
	)
	one_argument(argument, &arg[0])
	if ch.Desc.Original != nil {
		send_to_char(ch, libc.CString("You're already switched.\r\n"))
	} else if arg[0] == 0 {
		send_to_char(ch, libc.CString("Switch with who?\r\n"))
	} else if (func() *char_data {
		victim = get_char_vis(ch, &arg[0], nil, 1<<1)
		return victim
	}()) == nil {
		send_to_char(ch, libc.CString("No such character.\r\n"))
	} else if ch == victim {
		send_to_char(ch, libc.CString("Hee hee... we are jolly funny today, eh?\r\n"))
	} else if victim.Desc != nil {
		send_to_char(ch, libc.CString("You can't do that, the body is already in use!\r\n"))
	} else if !IS_NPC(victim) && !ADM_FLAGGED(ch, ADM_SWITCHMORTAL) {
		send_to_char(ch, libc.CString("You aren't holy enough to use a mortal's body.\r\n"))
	} else if ch.Admlevel < ADMLVL_VICE && ROOM_FLAGGED(victim.In_room, ROOM_GODROOM) {
		send_to_char(ch, libc.CString("You are not godly enough to use that room!\r\n"))
	} else {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		ch.Desc.Character = victim
		ch.Desc.Original = ch
		victim.Desc = ch.Desc
		ch.Desc = nil
	}
}
func do_return(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Desc != nil && ch.Desc.Original != nil {
		send_to_char(ch, libc.CString("You return to your original body.\r\n"))
		if ch.Desc.Original.Desc != nil {
			ch.Desc.Original.Desc.Character = nil
			ch.Desc.Original.Desc.Connected = CON_DISCONNECT
		}
		ch.Desc.Character = ch.Desc.Original
		ch.Desc.Original = nil
		ch.Desc.Character.Desc = ch.Desc
		ch.Desc = nil
	}
}
func do_load(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		buf2 [2048]byte
		buf3 [2048]byte
		i    int = 0
		n    int = 1
	)
	one_argument(two_arguments(argument, &buf[0], &buf2[0]), &buf3[0])
	if buf[0] == 0 || buf2[0] == 0 || !unicode.IsDigit(rune(buf2[0])) {
		send_to_char(ch, libc.CString("Usage: load { obj | mob } <vnum> (amt)\r\n"))
		return
	}
	if is_number(&buf2[0]) == 0 || is_number(&buf3[0]) == 0 {
		send_to_char(ch, libc.CString("That is not a number.\r\n"))
		return
	}
	if libc.Atoi(libc.GoString(&buf3[0])) > 0 {
		if libc.Atoi(libc.GoString(&buf3[0])) >= 100 {
			n = 100
		} else if libc.Atoi(libc.GoString(&buf3[0])) < 100 {
			n = libc.Atoi(libc.GoString(&buf3[0]))
		}
	} else {
		n = 1
	}
	if is_abbrev(&buf[0], libc.CString("mob")) != 0 {
		var (
			mob   *char_data = nil
			r_num mob_rnum
		)
		if (func() mob_rnum {
			r_num = real_mobile(mob_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
			return r_num
		}()) == mob_rnum(-1) {
			send_to_char(ch, libc.CString("There is no monster with that number.\r\n"))
			return
		}
		for i = 0; i < n; i++ {
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, ch.In_room)
			act(libc.CString("$n makes a quaint, magical gesture with one hand."), TRUE, ch, nil, nil, TO_ROOM)
			act(libc.CString("$n has created $N!"), FALSE, ch, nil, unsafe.Pointer(mob), TO_ROOM)
			act(libc.CString("You create $N."), FALSE, ch, nil, unsafe.Pointer(mob), TO_CHAR)
			load_mtrigger(mob)
		}
	} else if is_abbrev(&buf[0], libc.CString("obj")) != 0 {
		var (
			obj   *obj_data
			r_num obj_rnum
		)
		if (func() obj_rnum {
			r_num = real_object(obj_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
			return r_num
		}()) == obj_rnum(-1) {
			send_to_char(ch, libc.CString("There is no object with that number.\r\n"))
			return
		}
		for i = 0; i < n; i++ {
			obj = read_object(obj_vnum(r_num), REAL)
			add_unique_id(obj)
			if ch.Admlevel > 0 {
				send_to_imm(libc.CString("LOAD: %s has loaded a %s"), GET_NAME(ch), obj.Short_description)
				log_imm_action(libc.CString("LOAD: %s has loaded a %s"), GET_NAME(ch), obj.Short_description)
			}
			if config_info.Play.Load_into_inventory != 0 {
				obj_to_char(obj, ch)
			} else {
				obj_to_room(obj, ch.In_room)
			}
			act(libc.CString("$n makes a strange magical gesture."), TRUE, ch, nil, nil, TO_ROOM)
			act(libc.CString("$n has created $p!"), FALSE, ch, obj, nil, TO_ROOM)
			act(libc.CString("You create $p."), FALSE, ch, obj, nil, TO_CHAR)
			load_otrigger(obj)
		}
	} else {
		send_to_char(ch, libc.CString("That'll have to be either 'obj' or 'mob'.\r\n"))
	}
}
func do_vstat(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		buf2 [2048]byte
	)
	two_arguments(argument, &buf[0], &buf2[0])
	if buf[0] == 0 || buf2[0] == 0 || !unicode.IsDigit(rune(buf2[0])) {
		send_to_char(ch, libc.CString("Usage: vstat { obj | mob } <number>\r\n"))
		return
	}
	if is_number(&buf2[0]) == 0 {
		send_to_char(ch, libc.CString("That's not a valid number.\r\n"))
		return
	}
	if is_abbrev(&buf[0], libc.CString("mob")) != 0 {
		var (
			mob   *char_data
			r_num mob_rnum
		)
		if (func() mob_rnum {
			r_num = real_mobile(mob_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
			return r_num
		}()) == mob_rnum(-1) {
			send_to_char(ch, libc.CString("There is no monster with that number.\r\n"))
			return
		}
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, 0)
		do_stat_character(ch, mob)
		extract_char(mob)
	} else if is_abbrev(&buf[0], libc.CString("obj")) != 0 {
		var (
			obj   *obj_data
			r_num obj_rnum
		)
		if (func() obj_rnum {
			r_num = real_object(obj_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
			return r_num
		}()) == obj_rnum(-1) {
			send_to_char(ch, libc.CString("There is no object with that number.\r\n"))
			return
		}
		obj = read_object(obj_vnum(r_num), REAL)
		do_stat_object(ch, obj)
		extract_obj(obj)
	} else {
		send_to_char(ch, libc.CString("That'll have to be either 'obj' or 'mob'.\r\n"))
	}
}
func do_purge(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		vict *char_data
		obj  *obj_data
	)
	one_argument(argument, &buf[0])
	if buf[0] != 0 {
		if (func() *char_data {
			vict = get_char_vis(ch, &buf[0], nil, 1<<1)
			return vict
		}()) != nil {
			if !IS_NPC(vict) && ch.Admlevel <= vict.Admlevel {
				send_to_char(ch, libc.CString("Fuuuuuuuuu!\r\n"))
				return
			}
			act(libc.CString("$n disintegrates $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if !IS_NPC(vict) {
				send_to_all(libc.CString("@R%s@r purges @R%s's@r sorry ass right off the MUD!@n\r\n"), GET_NAME(ch), GET_NAME(vict))
			}
			if !IS_NPC(vict) {
				mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s has purged %s."), GET_NAME(ch), GET_NAME(vict))
				log_imm_action(libc.CString("PURGED: %s burned %s's sorry ass off the MUD!"), GET_NAME(ch), GET_NAME(vict))
				if vict.Desc != nil {
					vict.Desc.Connected = CON_CLOSE
					vict.Desc.Character = nil
					vict.Desc = nil
				}
			}
			extract_char(vict)
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &buf[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) != nil {
			act(libc.CString("$n destroys $p."), FALSE, ch, obj, nil, TO_ROOM)
			extract_obj(obj)
		} else {
			send_to_char(ch, libc.CString("Nothing here by that name.\r\n"))
			return
		}
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	} else {
		var i int
		act(libc.CString("$n gestures... You are surrounded by scorching flames!"), FALSE, ch, nil, nil, TO_ROOM)
		send_to_room(ch.In_room, libc.CString("The world seems a little cleaner.\r\n"))
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
			if !IS_NPC(vict) {
				continue
			}
			delete_inv_backup(vict)
			for vict.Carrying != nil {
				extract_obj(vict.Carrying)
			}
			for i = 0; i < NUM_WEARS; i++ {
				if (vict.Equipment[i]) != nil {
					extract_obj(vict.Equipment[i])
				}
			}
			extract_char(vict)
		}
		for (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents != nil {
			extract_obj((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
		}
	}
}

var logtypes [5]*byte = [5]*byte{libc.CString("off"), libc.CString("brief"), libc.CString("normal"), libc.CString("complete"), libc.CString("\n")}

func do_syslog(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		tp  int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Your syslog is currently %s.\r\n"), logtypes[(func() int {
			if PRF_FLAGGED(ch, PRF_LOG1) {
				return 1
			}
			return 0
		}())+(func() int {
			if PRF_FLAGGED(ch, PRF_LOG2) {
				return 2
			}
			return 0
		}())])
		return
	}
	if (func() int {
		tp = search_block(&arg[0], &logtypes[0], FALSE)
		return tp
	}()) == -1 {
		send_to_char(ch, libc.CString("Usage: syslog { Off | Brief | Normal | Complete }\r\n"))
		return
	}
	ch.Player_specials.Pref[int(PRF_LOG1/32)] &= bitvector_t(int32(^(1 << (int(PRF_LOG1 % 32)))))
	ch.Player_specials.Pref[int(PRF_LOG2/32)] &= bitvector_t(int32(^(1 << (int(PRF_LOG2 % 32)))))
	if tp&1 != 0 {
		ch.Player_specials.Pref[int(PRF_LOG1/32)] |= bitvector_t(int32(1 << (int(PRF_LOG1 % 32))))
	}
	if tp&2 != 0 {
		ch.Player_specials.Pref[int(PRF_LOG2/32)] |= bitvector_t(int32(1 << (int(PRF_LOG2 % 32))))
	}
	send_to_char(ch, libc.CString("Your syslog is now %s.\r\n"), logtypes[tp])
}
func do_copyover(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		secs int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		execute_copyover()
	} else if is_abbrev(&arg[0], libc.CString("cancel")) != 0 || is_abbrev(&arg[0], libc.CString("stop")) != 0 {
		if copyover_timer == 0 {
			send_to_char(ch, libc.CString("A timed copyover has not been started!\r\n"))
		} else {
			copyover_timer = 0
			game_info(libc.CString("Copyover cancelled"))
		}
	} else if is_abbrev(&arg[0], libc.CString("help")) != 0 {
		send_to_char(ch, libc.CString("COPYOVER\r\n\r\n"))
		send_to_char(ch, libc.CString("Usage: @ycopyover@n           - Perform an immediate copyover\r\n       @ycopyover <seconds>@n - Start a timed copyover\r\n       @ycopyover cancel@n    - Stop a timed copyover\r\n\r\n"))
		send_to_char(ch, libc.CString("A timed copyover will produce an automatic warning when it starts, and then\r\n"))
		send_to_char(ch, libc.CString("every minute.  During the last minute, there will be a warning every 15 seconds.\r\n"))
	} else {
		secs = libc.Atoi(libc.GoString(&arg[0]))
		if secs == 0 || secs < 0 {
			send_to_char(ch, libc.CString("Type @ycopyover help@n for usage info."))
		} else {
			copyover_timer = secs
			basic_mud_log(libc.CString("-- Timed Copyover started by %s - %d seconds until copyover --"), GET_NAME(ch), secs)
			if secs >= 60 {
				if secs%60 != 0 {
					game_info(libc.CString("A copyover will be performed in %d minutes and %d seconds."), copyover_timer/60, copyover_timer%60)
				} else {
					game_info(libc.CString("A copyover will be performed in %d minute%s."), copyover_timer/60, func() string {
						if (copyover_timer / 60) > 1 {
							return "s"
						}
						return ""
					}())
				}
			} else {
				game_info(libc.CString("A copyover will be performed in %d seconds."), copyover_timer)
			}
		}
	}
}
func execute_copyover() {
	var (
		fp     *stdio.File
		d      *descriptor_data
		d_next *descriptor_data
		buf    [100]byte
		buf2   [100]byte
	)
	fp = stdio.FOpen(COPYOVER_FILE, "w")
	if fp == nil {
		send_to_imm(libc.CString("Copyover file not writeable, aborted.\n\r"))
		return
	}
	save_all()
	save_mud_time(&time_info)
	stdio.Sprintf(&buf[0], "\t\x1b[1;31m \a\a\aThe universe stops for a moment as space and time fold.\x1b[0;0m\r\n")
	for d = descriptor_list; d != nil; d = d_next {
		var och *char_data = d.Character
		d_next = d.Next
		if d.Character == nil || d.Connected > CON_PLAYING {
			write_to_descriptor(d.Descriptor, libc.CString("\n\rSorry, we are rebooting. Come back in a few seconds.\n\r"))
			stdio.ByFD(uintptr(uintptr(unsafe.Pointer(d)))).Close()
		} else {
			if (func() room_vnum {
				if och.In_room != room_rnum(-1) && och.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(och.In_room)))).Number
				}
				return -1
			}()) > 1 {
				stdio.Fprintf(fp, "%d %s %s %d %s\n", d.Descriptor, GET_NAME(och), &d.Host[0], func() room_vnum {
					if och.In_room != room_rnum(-1) && och.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(och.In_room)))).Number
					}
					return -1
				}(), d.User)
			} else if (func() room_vnum {
				if och.In_room != room_rnum(-1) && och.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(och.In_room)))).Number
				}
				return -1
			}()) <= 1 && (func() room_vnum {
				if och.Was_in_room != room_rnum(-1) && och.Was_in_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(och.Was_in_room)))).Number
				}
				return -1
			}()) > 1 {
				stdio.Fprintf(fp, "%d %s %s %d %s\n", d.Descriptor, GET_NAME(och), &d.Host[0], func() room_vnum {
					if och.Was_in_room != room_rnum(-1) && och.Was_in_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(och.Was_in_room)))).Number
					}
					return -1
				}(), d.User)
			} else {
				stdio.Fprintf(fp, "%d %s %s 300 %s\n", d.Descriptor, GET_NAME(och), &d.Host[0], d.User)
			}
			basic_mud_log(libc.CString("printing descriptor name and host of connected players"))
			Crash_rentsave(och, 0)
			save_char(och)
			write_to_descriptor(d.Descriptor, &buf[0])
		}
	}
	stdio.Fprintf(fp, "-1\n")
	fp.Close()
	stdio.Sprintf(&buf[0], "%d", port)
	stdio.Chdir(libc.CString(".."))
	execl(libc.CString(EXE_FILE), libc.CString("circle"), &buf2[0], &buf[0], nil)
	perror(libc.CString("do_copyover: execl"))
	send_to_imm(libc.CString("Copyover FAILED!\n\r"))
	os.Exit(1)
}
func copyover_check() {
	if copyover_timer == 0 {
		return
	}
	copyover_timer--
	if copyover_timer == 0 {
		execute_copyover()
	}
	if copyover_timer > 59 {
		if copyover_timer%60 == 0 {
			game_info(libc.CString("A copyover will be performed in %d minute%s."), copyover_timer/60, func() string {
				if (copyover_timer / 60) > 1 {
					return "s"
				}
				return ""
			}())
		}
	} else {
		if copyover_timer%10 == 0 && copyover_timer > 29 {
			game_info(libc.CString("A copyover will be performed in %d seconds."), copyover_timer)
		}
		if copyover_timer%5 == 0 && copyover_timer <= 29 {
			game_info(libc.CString("A copyover will be performed in %d seconds."), copyover_timer)
		}
	}
}
func do_advance(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim   *char_data
		name     [2048]byte
		level    [2048]byte
		newlevel int
		oldlevel int
	)
	two_arguments(argument, &name[0], &level[0])
	if name[0] != 0 {
		if (func() *char_data {
			victim = get_char_vis(ch, &name[0], nil, 1<<1)
			return victim
		}()) == nil {
			send_to_char(ch, libc.CString("That player is not here.\r\n"))
			return
		}
	} else {
		send_to_char(ch, libc.CString("Advance who?\r\n"))
		return
	}
	if IS_NPC(victim) {
		send_to_char(ch, libc.CString("NO!  Not on NPC's.\r\n"))
		return
	}
	if level[0] == 0 {
		send_to_char(ch, libc.CString("[ 1 - 100 | demote ]\r\n"))
		return
	} else if (func() int {
		newlevel = libc.Atoi(libc.GoString(&level[0]))
		return newlevel
	}()) <= 0 {
		if libc.StrCaseCmp(libc.CString("demote"), &level[0]) == 0 {
			victim.Level = 1
			victim.Max_hit = 150
			victim.Max_mana = 150
			victim.Max_move = 150
			victim.Basepl = 150
			victim.Baseki = 150
			victim.Basest = 150
			send_to_char(ch, libc.CString("They have now been demoted!\r\n"))
			send_to_char(victim, libc.CString("You were demoted to level 1!\r\n"))
			return
		} else {
			send_to_char(ch, libc.CString("That's not a level!\r\n"))
			return
		}
	}
	if newlevel > 100 {
		send_to_char(ch, libc.CString("100 is the highest possible level.\r\n"))
		return
	}
	if newlevel == GET_LEVEL(victim) {
		send_to_char(ch, libc.CString("They are already at that level.\r\n"))
		return
	}
	oldlevel = GET_LEVEL(victim)
	if newlevel < GET_LEVEL(victim) {
		send_to_char(ch, libc.CString("You cannot demote a player.\r\n"))
	} else {
		act(libc.CString("$n makes some strange gestures.\r\nA strange feeling comes upon you, like a giant hand, light comes down\r\nfrom above, grabbing your body, which begins to pulse with colored\r\nlights from inside.\r\n\r\nYour head seems to be filled with demons from another plane as your\r\nbody dissolves to the elements of time and space itself.\r\n\r\nSuddenly a silent explosion of light snaps you back to reality.\r\n\r\nYou feel slightly different."), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	if newlevel < oldlevel {
		basic_mud_log(libc.CString("(GC) %s demoted %s from level %d to %d."), GET_NAME(ch), GET_NAME(victim), oldlevel, newlevel)
	} else {
		basic_mud_log(libc.CString("(GC) %s has advanced %s to level %d (from %d)"), GET_NAME(ch), GET_NAME(victim), newlevel, oldlevel)
	}
	gain_exp_regardless(victim, level_exp(victim, newlevel)-int(victim.Exp))
	save_char(victim)
}
func do_handout(ch *char_data, argument *byte, cmd int, subcmd int) {
	var j *descriptor_data
	if ch.Admlevel < 3 {
		send_to_char(ch, libc.CString("You can't do that.\r\n"))
		return
	}
	for j = descriptor_list; j != nil; j = j.Next {
		if !IS_PLAYING(j) || ch == j.Character || j.Character.Admlevel > 0 {
			continue
		}
		if IS_NPC(j.Character) {
			continue
		} else {
			j.Character.Player_specials.Class_skill_points[j.Character.Chclass] += 10
		}
	}
	send_to_all(libc.CString("@g%s@G hands out 10 practice sessions to everyone!@n\r\n"), GET_NAME(ch))
	basic_mud_log(libc.CString("%s gave a handout of 10 PS to everyone."), GET_NAME(ch))
	log_imm_action(libc.CString("HANDOUT: %s has handed out 10 PS to everyone."), GET_NAME(ch))
}
func do_restore(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf  [2048]byte
		vict *char_data
		j    *descriptor_data
		i    int
	)
	one_argument(argument, &buf[0])
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to restore?\r\n"))
	} else if is_abbrev(&buf[0], libc.CString("all")) != 0 {
		send_to_imm(libc.CString("[Log: %s restored all.]"), GET_NAME(ch))
		log_imm_action(libc.CString("RESTORE: %s has restored all players."), GET_NAME(ch))
		for j = descriptor_list; j != nil; j = j.Next {
			if !IS_PLAYING(j) || (func() *char_data {
				vict = j.Character
				return vict
			}()) == nil {
				continue
			}
			if vict.Hit < gear_pl(vict) {
				vict.Hit = gear_pl(vict)
			}
			vict.Mana = vict.Max_mana
			vict.Move = vict.Max_move
			update_pos(vict)
			act(libc.CString("You have been fully healed by $N!"), FALSE, vict, nil, unsafe.Pointer(ch), int(TO_CHAR|2<<7))
			if vict.Suppression > 0 && vict.Hit > ((gear_pl(vict)/100)*vict.Suppression) {
				vict.Hit = (gear_pl(vict) / 100) * vict.Suppression
				send_to_char(vict, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
			}
		}
		send_to_char(ch, libc.CString("Okay.\r\n"))
	} else if (func() *char_data {
		vict = get_char_vis(ch, &buf[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if !IS_NPC(vict) && ch != vict && vict.Admlevel >= ch.Admlevel {
		send_to_char(ch, libc.CString("They don't need your help.\r\n"))
	} else {
		vict.Hit = gear_pl(vict)
		vict.Mana = vict.Max_mana
		vict.Move = vict.Max_move
		vict.Ki = vict.Max_ki
		vict.Affected_by[int(AFF_BLIND/32)] &= ^(1 << (int(AFF_BLIND % 32)))
		vict.Limb_condition[1] = 100
		vict.Limb_condition[2] = 100
		vict.Limb_condition[3] = 100
		vict.Limb_condition[4] = 100
		vict.Act[int(PLR_HEAD/32)] |= bitvector_t(int32(1 << (int(PLR_HEAD % 32))))
		if !IS_NPC(vict) && ch.Admlevel >= ADMLVL_VICE {
			if vict.Admlevel >= ADMLVL_IMMORT {
				for i = 1; i <= MAX_SKILLS; i++ {
					for {
						vict.Skills[i] = 100
						if true {
							break
						}
					}
				}
			}
			if vict.Admlevel >= ADMLVL_GRGOD {
				vict.Real_abils.Intel = 25
				vict.Real_abils.Wis = 25
				vict.Real_abils.Dex = 25
				vict.Real_abils.Str = 25
				vict.Real_abils.Con = 25
				vict.Real_abils.Cha = 25
			}
		}
		update_pos(vict)
		affect_total(vict)
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		send_to_imm(libc.CString("[Log: %s restored %s.]"), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("RESTORE: %s has restored %s."), GET_NAME(ch), GET_NAME(vict))
		act(libc.CString("You have been fully healed by $N!"), FALSE, vict, nil, unsafe.Pointer(ch), TO_CHAR)
		if vict.Suppression > 0 && vict.Hit > ((gear_pl(vict)/100)*vict.Suppression) {
			vict.Hit = (gear_pl(vict) / 100) * vict.Suppression
			send_to_char(vict, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
		}
	}
}
func perform_immort_vis(ch *char_data) {
	ch.Player_specials.Invis_level = 0
}
func perform_immort_invis(ch *char_data, level int) {
	var tch *char_data
	for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = tch.Next_in_room {
		if tch == ch {
			continue
		}
		if tch.Admlevel >= int(ch.Player_specials.Invis_level) && tch.Admlevel < level {
			act(libc.CString("You blink and suddenly realize that $n is gone."), FALSE, ch, nil, unsafe.Pointer(tch), TO_VICT)
		}
		if tch.Admlevel < int(ch.Player_specials.Invis_level) && tch.Admlevel >= level {
			act(libc.CString("You suddenly realize that $n is standing beside you."), FALSE, ch, nil, unsafe.Pointer(tch), TO_VICT)
		}
	}
	ch.Player_specials.Invis_level = int16(level)
	send_to_char(ch, libc.CString("Your invisibility level is %d.\r\n"), level)
}
func do_invis(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		level int
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You can't do that!\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		if int(ch.Player_specials.Invis_level) > 0 {
			perform_immort_vis(ch)
		} else {
			perform_immort_invis(ch, ch.Admlevel)
		}
	} else {
		level = libc.Atoi(libc.GoString(&arg[0]))
		if level > ch.Admlevel {
			send_to_char(ch, libc.CString("You can't go invisible above your own level.\r\n"))
		} else if level < 1 {
			perform_immort_vis(ch)
		} else {
			perform_immort_invis(ch, level)
		}
	}
}
func do_gecho(ch *char_data, argument *byte, cmd int, subcmd int) {
	var pt *descriptor_data
	skip_spaces(&argument)
	delete_doubledollar(argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("That must be a mistake...\r\n"))
	} else {
		for pt = descriptor_list; pt != nil; pt = pt.Next {
			if IS_PLAYING(pt) && pt.Character != nil && pt.Character != ch {
				send_to_char(pt.Character, libc.CString("%s\r\n"), argument)
			}
		}
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else {
			send_to_char(ch, libc.CString("%s\r\n"), argument)
		}
	}
}
func do_ginfo(ch *char_data, argument *byte, cmd int, subcmd int) {
	var pt *descriptor_data
	skip_spaces(&argument)
	delete_doubledollar(argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("That must be a mistake...\r\n"))
	} else {
		for pt = descriptor_list; pt != nil; pt = pt.Next {
			if IS_PLAYING(pt) && pt.Character != nil && pt.Character != ch {
				send_to_char(pt.Character, libc.CString("@D[@GINFO@D] @g%s@n\r\n"), argument)
			}
		}
		if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_NOREPEAT) {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		} else {
			send_to_char(ch, libc.CString("@D[@GINFO@D] @g%s@n\r\n"), argument)
		}
	}
}
func do_poofset(ch *char_data, argument *byte, cmd int, subcmd int) {
	var msg **byte
	switch subcmd {
	case SCMD_POOFIN:
		msg = &ch.Player_specials.Poofin
	case SCMD_POOFOUT:
		msg = &ch.Player_specials.Poofout
	default:
		return
	}
	skip_spaces(&argument)
	if *msg != nil {
		libc.Free(unsafe.Pointer(*msg))
	}
	if *argument == 0 {
		*msg = nil
	} else {
		*msg = libc.StrDup(argument)
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
}
func do_dc(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg       [2048]byte
		d         *descriptor_data
		num_to_dc int
	)
	one_argument(argument, &arg[0])
	if (func() int {
		num_to_dc = libc.Atoi(libc.GoString(&arg[0]))
		return num_to_dc
	}()) == 0 {
		send_to_char(ch, libc.CString("Usage: DC <user number> (type USERS for a list)\r\n"))
		return
	}
	for d = descriptor_list; d != nil && d.Desc_num != num_to_dc; d = d.Next {
	}
	if d == nil {
		send_to_char(ch, libc.CString("No such connection.\r\n"))
		return
	}
	if d.Character != nil && d.Character.Admlevel >= ch.Admlevel {
		if !CAN_SEE(ch, d.Character) {
			send_to_char(ch, libc.CString("No such connection.\r\n"))
		} else {
			send_to_char(ch, libc.CString("Umm.. maybe that's not such a good idea...\r\n"))
		}
		return
	}
	if d.Connected == CON_DISCONNECT || d.Connected == CON_CLOSE {
		send_to_char(ch, libc.CString("They're already being disconnected.\r\n"))
	} else {
		if d.Connected == CON_PLAYING {
			d.Connected = CON_DISCONNECT
		} else {
			d.Connected = CON_CLOSE
		}
		send_to_char(ch, libc.CString("Connection #%d closed.\r\n"), num_to_dc)
		basic_mud_log(libc.CString("(GC) Connection closed by %s."), GET_NAME(ch))
	}
}
func do_wizlock(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		value int
		when  *byte
	)
	one_argument(argument, &arg[0])
	if arg[0] != 0 {
		value = libc.Atoi(libc.GoString(&arg[0]))
		if value < 0 || value > 101 {
			send_to_char(ch, libc.CString("Invalid wizlock value.\r\n"))
			return
		}
		circle_restrict = value
		when = libc.CString("now")
	} else {
		when = libc.CString("currently")
	}
	if arg[0] != 0 {
		switch circle_restrict {
		case 0:
			send_to_char(ch, libc.CString("The game is %s completely open.\r\n"), when)
			send_to_all(libc.CString("@RWIZLOCK@D: @WThe game has been completely opened by @C%s@W.@n"), GET_NAME(ch))
			basic_mud_log(libc.CString("WIZLOCK: The game has been completely opened by %s."), GET_NAME(ch))
		case 1:
			send_to_char(ch, libc.CString("The game is %s closed to new players.\r\n"), when)
			send_to_all(libc.CString("@RWIZLOCK@D: @WThe game is %s closed to new players by @C%s@W.@n"), when, GET_NAME(ch))
			basic_mud_log(libc.CString("WIZLOCK: The game is %s closed to new players by %s."), when, GET_NAME(ch))
		case 101:
			send_to_char(ch, libc.CString("The game is %s closed to non-imms.\r\n"), when)
			send_to_all(libc.CString("@RWIZLOCK@D: @WThe game is %s closed to non-imms by @C%s@W.@n"), when, GET_NAME(ch))
			basic_mud_log(libc.CString("WIZLOCK: The game is %s closed to non-imms by %s."), when, GET_NAME(ch))
		default:
			send_to_char(ch, libc.CString("Only level %d+ may enter the game %s.\r\n"), circle_restrict, when)
			send_to_all(libc.CString("@RWIZLOCK@D: @WLevel %d+ only can enter the game %s, thanks to @C%s@W.@n"), circle_restrict, when, GET_NAME(ch))
			basic_mud_log(libc.CString("WIZLOCK: Level %d+ only can enter the game %s, thanks to %s."), circle_restrict, when, GET_NAME(ch))
		}
	}
	if arg[0] == 0 {
		switch circle_restrict {
		case 0:
			send_to_char(ch, libc.CString("The game is %s completely open.\r\n"), when)
		case 1:
			send_to_char(ch, libc.CString("The game is %s closed to new players.\r\n"), when)
		default:
			send_to_char(ch, libc.CString("Only level %d and above may enter the game %s.\r\n"), circle_restrict, when)
		}
	}
}
func do_date(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		tmstr  *byte
		mytime libc.Time
		d      int
		h      int
		m      int
	)
	if subcmd == SCMD_DATE {
		mytime = libc.GetTime(nil)
	} else {
		mytime = boot_time
	}
	tmstr = libc.AscTime(libc.LocalTime(&mytime))
	*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
	if subcmd == SCMD_DATE {
		send_to_char(ch, libc.CString("Current machine time: %s\r\n"), tmstr)
	} else {
		mytime = libc.GetTime(nil) - boot_time
		d = int(mytime / 86400)
		h = int((mytime / 3600) % 24)
		m = int((mytime / 60) % 60)
		send_to_char(ch, libc.CString("Up since %s: %d day%s, %d:%02d\r\n"), tmstr, d, func() string {
			if d == 1 {
				return ""
			}
			return "s"
		}(), h, m)
	}
}
func do_last(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		vict *char_data = nil
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("For whom do you wish to search?\r\n"))
		return
	}
	vict = new(char_data)
	clear_char(vict)
	vict.Player_specials = new(player_special_data)
	if load_char(&arg[0], vict) < 0 {
		send_to_char(ch, libc.CString("There is no such player.\r\n"))
		free_char(vict)
		return
	}
	if vict.Admlevel > ch.Admlevel && ch.Admlevel < ADMLVL_IMPL {
		send_to_char(ch, libc.CString("You are not sufficiently godly for that!\r\n"))
		return
	}
	send_to_char(ch, libc.CString("[%5d] [%2d %s %s] %-12s : %-18s : %-20s\r\n"), vict.Idnum, GET_LEVEL(vict), race_abbrevs[int(vict.Race)], class_abbrevs[int(vict.Chclass)], GET_NAME(vict), func() *byte {
		if vict.Player_specials.Host != nil && *vict.Player_specials.Host != 0 {
			return vict.Player_specials.Host
		}
		return libc.CString("(NOHOST)")
	}(), ctime(&vict.Time.Logon))
	free_char(vict)
}
func do_force(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i          *descriptor_data
		next_desc  *descriptor_data
		vict       *char_data
		next_force *char_data
		arg        [2048]byte
		to_force   [2048]byte
		buf1       [2080]byte
	)
	half_chop(argument, &arg[0], &to_force[0])
	stdio.Snprintf(&buf1[0], int(2080), "$n has forced you to '%s'.", &to_force[0])
	if arg[0] == 0 || to_force[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to force do what?\r\n"))
	} else if !ADM_FLAGGED(ch, ADM_FORCEMASS) || libc.StrCaseCmp(libc.CString("all"), &arg[0]) != 0 && libc.StrCaseCmp(libc.CString("room"), &arg[0]) != 0 {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		} else if !IS_NPC(vict) && ch.Admlevel <= vict.Admlevel {
			send_to_char(ch, libc.CString("No, no, no!\r\n"))
		} else {
			send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
			act(&buf1[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s forced %s to %s"), GET_NAME(ch), GET_NAME(vict), &to_force[0])
			command_interpreter(vict, &to_force[0])
		}
	} else if libc.StrCaseCmp(libc.CString("room"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s forced room %d to %s"), GET_NAME(ch), func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}(), &to_force[0])
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_force {
			next_force = vict.Next_in_room
			if !IS_NPC(vict) && vict.Admlevel >= ch.Admlevel {
				continue
			}
			act(&buf1[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			command_interpreter(vict, &to_force[0])
		}
	} else {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s forced all to %s"), GET_NAME(ch), &to_force[0])
		for i = descriptor_list; i != nil; i = next_desc {
			next_desc = i.Next
			if i.Connected != CON_PLAYING || (func() *char_data {
				vict = i.Character
				return vict
			}()) == nil || !IS_NPC(vict) && vict.Admlevel >= ch.Admlevel {
				continue
			}
			act(&buf1[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			command_interpreter(vict, &to_force[0])
		}
	}
}
func do_wiznet(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf1  [2100]byte
		buf2  [2100]byte
		msg   *byte
		d     *descriptor_data
		emote int8 = FALSE
		any   int8 = FALSE
		level int  = ADMLVL_IMMORT
	)
	skip_spaces(&argument)
	delete_doubledollar(argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Usage: wiznet <text> | #<level> <text> | *<emotetext> |\r\n        wiznet @<level> *<emotetext> | wiz @\r\n"))
		return
	}
	switch *argument {
	case '*':
		emote = TRUE
		fallthrough
	case '#':
		one_argument((*byte)(unsafe.Add(unsafe.Pointer(argument), 1)), &buf1[0])
		if is_number(&buf1[0]) != 0 {
			half_chop((*byte)(unsafe.Add(unsafe.Pointer(argument), 1)), &buf1[0], argument)
			level = MAX(libc.Atoi(libc.GoString(&buf1[0])), ADMLVL_IMMORT)
			if level > ch.Admlevel {
				send_to_char(ch, libc.CString("You can't wizline above your own level.\r\n"))
				return
			}
		} else if int(emote) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
	case '@':
		send_to_char(ch, libc.CString("God channel status:\r\n"))
		for func() *descriptor_data {
			any = 0
			return func() *descriptor_data {
				d = descriptor_list
				return d
			}()
		}(); d != nil; d = d.Next {
			if d.Connected != CON_PLAYING || d.Character.Admlevel < ADMLVL_IMMORT {
				continue
			}
			if !CAN_SEE(ch, d.Character) {
				continue
			}
			send_to_char(ch, libc.CString("  %-*s%s%s%s\r\n"), MAX_NAME_LENGTH, GET_NAME(d.Character), func() string {
				if PLR_FLAGGED(d.Character, PLR_WRITING) {
					return " (Writing)"
				}
				return ""
			}(), func() string {
				if PLR_FLAGGED(d.Character, PLR_MAILING) {
					return " (Writing mail)"
				}
				return ""
			}(), func() string {
				if PRF_FLAGGED(d.Character, PRF_NOWIZ) {
					return " (Offline)"
				}
				return ""
			}())
		}
		return
	case '\\':
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	default:
	}
	if PRF_FLAGGED(ch, PRF_NOWIZ) {
		send_to_char(ch, libc.CString("You are offline!\r\n"))
		return
	}
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Don't bother the gods like that!\r\n"))
		return
	}
	if level > ADMLVL_IMMORT {
		stdio.Snprintf(&buf1[0], int(2100), "@c%s@D: <@C%d@D> @G%s%s@n\r\n", GET_NAME(ch), level, func() string {
			if int(emote) != 0 {
				return "<--- "
			}
			return ""
		}(), argument)
		stdio.Snprintf(&buf2[0], int(2100), "@cSomeone@D: <@C%d@D> @G%s%s@n\r\n", level, func() string {
			if int(emote) != 0 {
				return "<--- "
			}
			return ""
		}(), argument)
	} else {
		stdio.Snprintf(&buf1[0], int(2100), "@c%s@D: @G%s%s@n\r\n", GET_NAME(ch), func() string {
			if int(emote) != 0 {
				return "<--- "
			}
			return ""
		}(), argument)
		stdio.Snprintf(&buf2[0], int(2100), "@cSomeone@D: @G%s%s@n\r\n", func() string {
			if int(emote) != 0 {
				return "<--- "
			}
			return ""
		}(), argument)
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if IS_PLAYING(d) && d.Character.Admlevel >= level && !PRF_FLAGGED(d.Character, PRF_NOWIZ) && !PLR_FLAGGED(d.Character, bitvector_t(int32(int(PLR_WRITING|PLR_MAILING)))) && (d != ch.Desc || !PRF_FLAGGED(d.Character, PRF_NOREPEAT)) {
			if CAN_SEE(d.Character, ch) {
				msg = libc.StrDup(&buf1[0])
				send_to_char(d.Character, libc.CString("%s"), &buf1[0])
			} else {
				msg = libc.StrDup(&buf2[0])
				send_to_char(d.Character, libc.CString("%s"), &buf2[0])
			}
			add_history(d.Character, msg, HIST_WIZNET)
		}
	}
	if PRF_FLAGGED(ch, PRF_NOREPEAT) {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	}
}
func do_zreset(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		i   zone_rnum
		j   zone_vnum
	)
	one_argument(argument, &arg[0])
	if arg[0] == '*' {
		if ch.Admlevel < ADMLVL_VICE {
			send_to_char(ch, libc.CString("You do not have permission to reset the entire world.\r\n"))
			return
		} else {
			for i = 0; i <= top_of_zone_table; i++ {
				if i < 200 {
					reset_zone(i)
				}
			}
			send_to_char(ch, libc.CString("Reset world.\r\n"))
			mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s reset all MUD zones."), GET_NAME(ch))
			log_imm_action(libc.CString("RESET: %s has reset all MUD zones."), GET_NAME(ch))
			return
		}
	} else if arg[0] == '.' || arg[0] == 0 {
		i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone
	} else {
		j = zone_vnum(libc.Atoi(libc.GoString(&arg[0])))
		for i = 0; i <= top_of_zone_table; i++ {
			if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number == j {
				break
			}
		}
	}
	if i <= top_of_zone_table && (can_edit_zone(ch, i) != 0 || ch.Admlevel > ADMLVL_IMMORT) {
		reset_zone(i)
		send_to_char(ch, libc.CString("Reset zone #%d: %s.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Name)
		mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s reset zone %d (%s)"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Name)
		log_imm_action(libc.CString("RESET: %s has reset zone #%d: %s."), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Name)
	} else {
		send_to_char(ch, libc.CString("You do not have permission to reset this zone. Try %d.\r\n"), ch.Player_specials.Olc_zone)
	}
}
func do_wizutil(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg     [2048]byte
		vict    *char_data
		taeller int
		result  int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Yes, but for whom?!?\r\n"))
	} else if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("There is no such player.\r\n"))
	} else if IS_NPC(vict) {
		send_to_char(ch, libc.CString("You can't do that to a mob!\r\n"))
	} else if vict.Admlevel > ch.Admlevel {
		send_to_char(ch, libc.CString("Hmmm...you'd better not.\r\n"))
	} else {
		switch subcmd {
		case SCMD_REROLL:
			send_to_char(ch, libc.CString("Rerolling is not possible at this time, bug Iovan about it...\r\n"))
			basic_mud_log(libc.CString("(GC) %s has rerolled %s."), GET_NAME(ch), GET_NAME(vict))
			send_to_char(ch, libc.CString("New stats: Str %d, Int %d, Wis %d, Dex %d, Con %d, Cha %d\r\n"), vict.Aff_abils.Str, vict.Aff_abils.Intel, vict.Aff_abils.Wis, vict.Aff_abils.Dex, vict.Aff_abils.Con, vict.Aff_abils.Cha)
		case SCMD_PARDON:
			if !PLR_FLAGGED(vict, PLR_THIEF) && !PLR_FLAGGED(vict, PLR_KILLER) {
				send_to_char(ch, libc.CString("Your victim is not flagged.\r\n"))
				return
			}
			vict.Act[int(PLR_THIEF/32)] &= bitvector_t(int32(^(1 << (int(PLR_THIEF % 32)))))
			vict.Act[int(PLR_KILLER/32)] &= bitvector_t(int32(^(1 << (int(PLR_KILLER % 32)))))
			send_to_char(ch, libc.CString("Pardoned.\r\n"))
			send_to_char(vict, libc.CString("You have been pardoned by the Gods!\r\n"))
			mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s pardoned by %s"), GET_NAME(vict), GET_NAME(ch))
		case SCMD_NOTITLE:
			result = int(func() bitvector_t {
				p := &vict.Act[int(PLR_NOTITLE/32)]
				vict.Act[int(PLR_NOTITLE/32)] = bitvector_t(int32(int(vict.Act[int(PLR_NOTITLE/32)]) ^ 1<<(int(PLR_NOTITLE%32))))
				return *p
			}()) & (1 << (int(PLR_NOTITLE % 32)))
			mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) Notitle %s for %s by %s."), func() string {
				if result != 0 {
					return "ON"
				}
				return "OFF"
			}(), GET_NAME(vict), GET_NAME(ch))
			send_to_char(ch, libc.CString("(GC) Notitle %s for %s by %s.\r\n"), func() string {
				if result != 0 {
					return "ON"
				}
				return "OFF"
			}(), GET_NAME(vict), GET_NAME(ch))
		case SCMD_SQUELCH:
			result = int(func() bitvector_t {
				p := &vict.Act[int(PLR_NOSHOUT/32)]
				vict.Act[int(PLR_NOSHOUT/32)] = bitvector_t(int32(int(vict.Act[int(PLR_NOSHOUT/32)]) ^ 1<<(int(PLR_NOSHOUT%32))))
				return *p
			}()) & (1 << (int(PLR_NOSHOUT % 32)))
			mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) Squelch %s for %s by %s."), func() string {
				if result != 0 {
					return "ON"
				}
				return "OFF"
			}(), GET_NAME(vict), GET_NAME(ch))
			send_to_char(ch, libc.CString("(GC) Mute turned %s for %s by %s.\r\n"), func() string {
				if result != 0 {
					return "ON"
				}
				return "OFF"
			}(), GET_NAME(vict), GET_NAME(ch))
			send_to_all(libc.CString("@D[@RMUTE@D] @C%s@W has had mute turned @r%s@W by @C%s@W.\r\n"), GET_NAME(vict), func() string {
				if result != 0 {
					return "ON"
				}
				return "OFF"
			}(), GET_NAME(ch))
		case SCMD_FREEZE:
			if ch == vict {
				send_to_char(ch, libc.CString("Oh, yeah, THAT'S real smart...\r\n"))
				return
			}
			if ch.Admlevel <= vict.Admlevel {
				send_to_char(ch, libc.CString("Pfft...\r\n"))
				return
			}
			if PLR_FLAGGED(vict, PLR_FROZEN) {
				send_to_char(ch, libc.CString("Your victim is already pretty cold.\r\n"))
				return
			}
			vict.Act[int(PLR_FROZEN/32)] |= bitvector_t(int32(1 << (int(PLR_FROZEN % 32))))
			vict.Player_specials.Freeze_level = int8(ch.Admlevel)
			send_to_char(vict, libc.CString("A bitter wind suddenly rises and drains every erg of heat from your body!\r\nYou feel frozen!\r\n"))
			send_to_char(ch, libc.CString("Frozen.\r\n"))
			act(libc.CString("A sudden cold wind conjured from nowhere freezes $n!"), FALSE, vict, nil, nil, TO_ROOM)
			mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s frozen by %s."), GET_NAME(vict), GET_NAME(ch))
		case SCMD_THAW:
			if !PLR_FLAGGED(vict, PLR_FROZEN) {
				send_to_char(ch, libc.CString("Sorry, your victim is not morbidly encased in ice at the moment.\r\n"))
				return
			}
			if int(vict.Player_specials.Freeze_level) > ch.Admlevel {
				send_to_char(ch, libc.CString("Sorry, a level %d God froze %s... you can't unfreeze %s.\r\n"), vict.Player_specials.Freeze_level, GET_NAME(vict), HMHR(vict))
				return
			}
			mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s un-frozen by %s."), GET_NAME(vict), GET_NAME(ch))
			vict.Act[int(PLR_FROZEN/32)] &= bitvector_t(int32(^(1 << (int(PLR_FROZEN % 32)))))
			send_to_char(vict, libc.CString("A fireball suddenly explodes in front of you, melting the ice!\r\nYou feel thawed.\r\n"))
			send_to_char(ch, libc.CString("Thawed.\r\n"))
			act(libc.CString("A sudden fireball conjured from nowhere thaws $n!"), FALSE, vict, nil, nil, TO_ROOM)
		case SCMD_UNAFFECT:
			if vict.Affected != nil || vict.Affected_by != nil || vict.Affectedv != nil {
				for vict.Affected != nil {
					affect_remove(vict, vict.Affected)
				}
				for taeller = 0; taeller < AF_ARRAY_MAX; taeller++ {
					ch.Affected_by[taeller] = 0
				}
				for vict.Affectedv != nil {
					affectv_remove(vict, vict.Affectedv)
				}
				for taeller = 0; taeller < AF_ARRAY_MAX; taeller++ {
					ch.Affected_by[taeller] = 0
				}
				send_to_char(vict, libc.CString("There is a brief flash of light!\r\nYou feel slightly different.\r\n"))
				send_to_char(ch, libc.CString("All spells removed.\r\n"))
			} else {
				send_to_char(ch, libc.CString("Your victim does not have any affections!\r\n"))
				return
			}
		default:
			basic_mud_log(libc.CString("SYSERR: Unknown subcmd %d passed to do_wizutil (%s)"), subcmd, __FILE__)
		}
		save_char(vict)
	}
}
func print_zone_to_buf(bufptr *byte, left uint64, zone zone_rnum, listall int) uint64 {
	var tmp uint64
	if listall != 0 {
		var (
			i int
			j int
			k int
			l int
			m int
			n int
			o int
		)
		tmp = uint64(stdio.Snprintf(bufptr, int(left), "%3d %-30.30s By: %-10.10s Age: %3d; Reset: %3d (%1d); Range: %5d-%5d\r\n", (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Name, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Builders, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Age, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Lifespan, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Reset_mode, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top))
		i = func() int {
			j = func() int {
				k = func() int {
					l = func() int {
						m = func() int {
							n = func() int {
								o = 0
								return o
							}()
							return n
						}()
						return m
					}()
					return l
				}()
				return k
			}()
			return j
		}()
		for i = 0; i < int(top_of_world); i++ {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number >= (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number <= (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top {
				j++
			}
		}
		for i = 0; i < int(top_of_objt); i++ {
			if (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot) && (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top) {
				k++
			}
		}
		for i = 0; i < int(top_of_mobt); i++ {
			if (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot) && (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top) {
				l++
			}
		}
		m = count_shops(shop_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot), shop_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top))
		for i = 0; i < top_of_trigt; i++ {
			if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum >= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot) && (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum <= mob_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top) {
				n++
			}
		}
		o = count_guilds(guild_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot), guild_vnum((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top))
		tmp += uint64(stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(bufptr), tmp)), int(left-tmp), "       Zone stats:\r\n       ---------------\r\n         Rooms:    %2d\r\n         Objects:  %2d\r\n         Mobiles:  %2d\r\n         Shops:    %2d\r\n         Triggers: %2d\r\n         Guilds:   %2d\r\n", j, k, l, m, n, o))
		return tmp
	}
	return uint64(stdio.Snprintf(bufptr, int(left), "%3d %-*s By: %-10.10s Range: %5d-%5d\r\n", (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Number, count_color_chars((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Name)+30, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Name, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Builders, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Bot, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Top))
}
func do_show(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i     int
		j     int
		k     int
		l     int
		con   int
		len_  uint64
		nlen  uint64
		zrn   zone_rnum
		zvn   zone_vnum
		low   int
		high  int
		self  int8       = FALSE
		vict  *char_data = nil
		obj   *obj_data
		d     *descriptor_data
		aff   *affected_type
		field [2048]byte
		value [2048]byte
		strp  *byte
		arg   [2048]byte
		buf   [64936]byte
	)
	type show_struct struct {
		Cmd   *byte
		Level int8
	}
	var fields [18]show_struct = [18]show_struct{{Cmd: libc.CString("nothing"), Level: 0}, {Cmd: libc.CString("zones"), Level: ADMLVL_IMMORT}, {Cmd: libc.CString("player"), Level: ADMLVL_GOD}, {Cmd: libc.CString("rent"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("stats"), Level: ADMLVL_IMMORT}, {Cmd: libc.CString("errors"), Level: ADMLVL_IMPL}, {Cmd: libc.CString("death"), Level: ADMLVL_GOD}, {Cmd: libc.CString("godrooms"), Level: ADMLVL_IMMORT}, {Cmd: libc.CString("shops"), Level: ADMLVL_IMMORT}, {Cmd: libc.CString("houses"), Level: ADMLVL_GOD}, {Cmd: libc.CString("snoop"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("assemblies"), Level: ADMLVL_IMMORT}, {Cmd: libc.CString("guilds"), Level: ADMLVL_GOD}, {Cmd: libc.CString("levels"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("uniques"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("affect"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("affectv"), Level: ADMLVL_GRGOD}, {Cmd: libc.CString("\n"), Level: 0}}
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Game Info options:\r\n"))
		for func() int {
			j = 0
			return func() int {
				i = 1
				return i
			}()
		}(); int(fields[i].Level) != 0; i++ {
			if int(fields[i].Level) <= ch.Admlevel {
				send_to_char(ch, libc.CString("%-15s%s"), fields[i].Cmd, func() string {
					if (func() int {
						p := &j
						*p++
						return *p
					}() % 5) == 0 {
						return "\r\n"
					}
					return ""
				}())
			}
		}
		send_to_char(ch, libc.CString("\r\n"))
		return
	}
	libc.StrCpy(&arg[0], two_arguments(argument, &field[0], &value[0]))
	for l = 0; *fields[l].Cmd != '\n'; l++ {
		if libc.StrNCmp(&field[0], fields[l].Cmd, libc.StrLen(&field[0])) == 0 {
			break
		}
	}
	if ch.Admlevel < int(fields[l].Level) {
		send_to_char(ch, libc.CString("You are not godly enough for that!\r\n"))
		return
	}
	if libc.StrCmp(&value[0], libc.CString(".")) == 0 {
		self = TRUE
	}
	buf[0] = '\x00'
	switch l {
	case 1:
		if int(self) != 0 {
			print_zone_to_buf(&buf[0], uint64(64936), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone, 1)
		} else if value[0] != 0 && is_number(&value[0]) != 0 {
			for func() zone_rnum {
				zvn = zone_vnum(libc.Atoi(libc.GoString(&value[0])))
				return func() zone_rnum {
					zrn = 0
					return zrn
				}()
			}(); (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrn)))).Number != zvn && zrn <= top_of_zone_table; zrn++ {
			}
			if zrn <= top_of_zone_table {
				print_zone_to_buf(&buf[0], uint64(64936), zrn, 1)
			} else {
				send_to_char(ch, libc.CString("That is not a valid zone.\r\n"))
				return
			}
		} else {
			for len_ = uint64(func() zone_rnum {
				zrn = 0
				return zrn
			}()); zrn <= top_of_zone_table; zrn++ {
				nlen = print_zone_to_buf(&buf[len_], uint64(64936-uintptr(len_)), zrn, 0)
				if len_+nlen >= uint64(64936) || nlen < 0 {
					break
				}
				len_ += nlen
			}
		}
		page_string(ch.Desc, &buf[0], TRUE)
	case 2:
		if value[0] == 0 {
			send_to_char(ch, libc.CString("A name would help.\r\n"))
			return
		}
		vict = new(char_data)
		clear_char(vict)
		vict.Player_specials = new(player_special_data)
		if load_char(&value[0], vict) < 0 {
			send_to_char(ch, libc.CString("There is no such player.\r\n"))
			free_char(vict)
			return
		}
		send_to_char(ch, libc.CString("Player: %-12s (%s) [%2d %s %s]\r\n"), GET_NAME(vict), genders[int(vict.Sex)], GET_LEVEL(vict), class_abbrevs[int(vict.Chclass)], race_abbrevs[int(vict.Race)])
		send_to_char(ch, libc.CString("Au: %-8d  Bal: %-8d  Exp: %lld  Align: %-5d  Ethic: %-5d\r\n"), vict.Gold, vict.Bank_gold, vict.Exp, vict.Alignment, vict.Alignment_ethic)
		if config_info.Advance.Allow_multiclass != 0 {
			send_to_char(ch, libc.CString("Class ranks: %s\r\n"), class_desc_str(vict, 1, 0))
		}
		send_to_char(ch, libc.CString("Started: %-20.16s  "), ctime(&vict.Time.Created))
		send_to_char(ch, libc.CString("Last: %-20.16s  Played: %3dh %2dm\r\n"), ctime(&vict.Time.Logon), int(vict.Time.Played/3600), int(vict.Time.Played/60%60))
		free_char(vict)
	case 3:
		if value[0] == 0 {
			send_to_char(ch, libc.CString("A name would help.\r\n"))
			return
		}
		Crash_listrent(ch, &value[0])
	case 4:
		i = 0
		j = 0
		k = 0
		con = 0
		for vict = character_list; vict != nil; vict = vict.Next {
			if IS_NPC(vict) {
				j++
			} else if CAN_SEE(ch, vict) {
				i++
				if vict.Desc != nil {
					con++
				}
			}
		}
		for obj = object_list; obj != nil; obj = obj.Next {
			k++
		}
		send_to_char(ch, libc.CString("             @D---   @CCore Stats   @D---\r\n  @Y%5d@W players in game  @y%5d@W connected\r\n  @Y%5d@W registered\r\n  @Y%5d@W mobiles          @y%5d@W prototypes\r\n  @Y%5d@W objects          @y%5d@W prototypes\r\n  @Y%5d@W rooms            @y%5d@W zones\r\n  @Y%5d@W triggers\r\n  @Y%5d@W large bufs\r\n  @Y%5d@W buf switches     @y%5d@W overflows\r\n             @D--- @CMiscellaneous  @D---\r\n  @Y%5s@W Mob ki attacks this boot\r\n  @Y%5s@W Asssassins Generated@n\r\n  @Y%5d@W Wish Selfishness Meter@n\r\n"), i, con, top_of_p_table+1, j, top_of_mobt+1, k, top_of_objt+1, top_of_world+1, top_of_zone_table+1, top_of_trigt+1, buf_largecount, buf_switches, buf_overflows, add_commas(int64(mob_specials_used)), add_commas(int64(number_of_assassins)), SELFISHMETER)
	case 5:
		len_ = strlcpy(&buf[0], libc.CString("Errant Rooms\r\n------------\r\n"), uint64(64936))
		for func() int {
			i = 0
			return func() int {
				k = 0
				return k
			}()
		}(); i <= int(top_of_world); i++ {
			for j = 0; j < NUM_OF_DIRS; j++ {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]) == nil {
					continue
				}
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room == 0 {
					nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%2d: (void   ) [%5d] %-*s%s (%s)\r\n", func() int {
						p := &k
						*p++
						return *p
					}(), func() room_vnum {
						if i != int(-1) && i <= int(top_of_world) {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number
						}
						return -1
					}(), count_color_chars((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name)+40, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name, func() string {
						if (func() int {
							if !IS_NPC(ch) {
								if PRF_FLAGGED(ch, PRF_COLOR) {
									return 1
								}
								return 0
							}
							return 0
						}()) >= C_ON {
							return KNRM
						}
						return KNUL
					}(), dirs[j]))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				}
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room == room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).General_description == nil {
					nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%2d: (Nowhere) [%5d] %-*s%s (%s)\r\n", func() int {
						p := &k
						*p++
						return *p
					}(), func() room_vnum {
						if i != int(-1) && i <= int(top_of_world) {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number
						}
						return -1
					}(), count_color_chars((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name)+40, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name, func() string {
						if (func() int {
							if !IS_NPC(ch) {
								if PRF_FLAGGED(ch, PRF_COLOR) {
									return 1
								}
								return 0
							}
							return 0
						}()) >= C_ON {
							return KNRM
						}
						return KNUL
					}(), dirs[j]))
					if len_+nlen >= uint64(64936) || nlen < 0 {
						break
					}
					len_ += nlen
				}
			}
		}
		page_string(ch.Desc, &buf[0], TRUE)
	case 6:
		len_ = strlcpy(&buf[0], libc.CString("Death Traps\r\n-----------\r\n"), uint64(64936))
		for func() int {
			i = 0
			return func() int {
				j = 0
				return j
			}()
		}(); i <= int(top_of_world); i++ {
			if ROOM_FLAGGED(room_rnum(i), ROOM_DEATH) {
				nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%2d: [%5d] %s\r\n", func() int {
					p := &j
					*p++
					return *p
				}(), func() room_vnum {
					if i != int(-1) && i <= int(top_of_world) {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number
					}
					return -1
				}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name))
				if len_+nlen >= uint64(64936) || nlen < 0 {
					break
				}
				len_ += nlen
			}
		}
		page_string(ch.Desc, &buf[0], TRUE)
	case 7:
		len_ = strlcpy(&buf[0], libc.CString("Godrooms\r\n--------------------------\r\n"), uint64(64936))
		for func() int {
			i = 0
			return func() int {
				j = 0
				return j
			}()
		}(); i <= int(top_of_world); i++ {
			if ROOM_FLAGGED(room_rnum(i), ROOM_GODROOM) {
				nlen = uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%2d: [%5d] %s\r\n", func() int {
					p := &j
					*p++
					return *p
				}(), func() room_vnum {
					if i != int(-1) && i <= int(top_of_world) {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number
					}
					return -1
				}(), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name))
				if len_+nlen >= uint64(64936) || nlen < 0 {
					break
				}
				len_ += nlen
			}
		}
		page_string(ch.Desc, &buf[0], TRUE)
	case 8:
		show_shops(ch, &value[0])
	case 9:
		hcontrol_list_houses(ch)
	case 10:
		i = 0
		send_to_char(ch, libc.CString("People currently snooping:\r\n--------------------------\r\n"))
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Snooping == nil || d.Character == nil {
				continue
			}
			if d.Connected != CON_PLAYING || ch.Admlevel < d.Character.Admlevel {
				continue
			}
			if !CAN_SEE(ch, d.Character) || d.Character.In_room == room_rnum(-1) {
				continue
			}
			i++
			send_to_char(ch, libc.CString("%-10s - snooped by %s.\r\n"), GET_NAME(d.Snooping.Character), GET_NAME(d.Character))
		}
		if i == 0 {
			send_to_char(ch, libc.CString("No one is currently snooping.\r\n"))
		}
	case 11:
		assemblyListToChar(ch)
	case 12:
		show_guild(ch, &value[0])
	case 13:
		send_to_char(ch, libc.CString("This is not used currently.\r\n"))
	case 14:
		if value != nil && value[0] != 0 {
			if stdio.Sscanf(&value[0], "%d-%d", &low, &high) != 2 {
				if stdio.Sscanf(&value[0], "%d", &low) != 1 {
					send_to_char(ch, libc.CString("Usage: show uniques, show uniques [vnum], or show uniques [low-high]\r\n"))
					return
				} else {
					high = low
				}
			}
		} else {
			low = -1
			high = 0x98967F
		}
		strp = sprintuniques(low, high)
		page_string(ch.Desc, strp, TRUE)
		libc.Free(unsafe.Pointer(strp))
	case 15:
		fallthrough
	case 16:
		if value[0] == 0 {
			low = 1
			if l == 15 {
				vict = affect_list
			} else {
				vict = affectv_list
			}
		} else {
			low = 0
			if (func() *char_data {
				vict = get_char_world_vis(ch, &value[0], nil)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Cannot find that character.\r\n"))
				return
			}
		}
		k = MAX_STRING_LENGTH
		strp = (*byte)(unsafe.Pointer(&make([]int8, k)[0]))
		*strp = byte(int8(func() int {
			j = 0
			return j
		}()))
		if vict == nil {
			send_to_char(ch, libc.CString("None.\r\n"))
			return
		}
		for {
			if (k - j) < (int(MAX_INPUT_LENGTH * 8)) {
				k *= 2
				strp = (*byte)(libc.Realloc(unsafe.Pointer(strp), k*int(unsafe.Sizeof(int8(0)))))
			}
			j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, "Name: %s\r\n", GET_NAME(vict))
			if l == 15 {
				aff = vict.Affected
			} else {
				aff = vict.Affectedv
			}
			for ; aff != nil; aff = aff.Next {
				j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, "SPL: (%3d%s) @c%-21s@n ", int(aff.Duration)+1, func() string {
					if l == 15 {
						return "hr"
					}
					return "rd"
				}(), skill_name(int(aff.Type)))
				if aff.Modifier != 0 {
					j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, "%+d to %s", aff.Modifier, apply_types[aff.Location])
				}
				if aff.Bitvector != 0 {
					if aff.Modifier != 0 {
						j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, ", ")
					}
					libc.StrCpy(&field[0], affected_bits[aff.Bitvector])
					j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, "sets %s", &field[0])
				}
				j += stdio.Snprintf((*byte)(unsafe.Add(unsafe.Pointer(strp), j)), k-j, "\r\n")
			}
			if l == 15 {
				vict = vict.Next_affect
			} else {
				vict = vict.Next_affectv
			}
			if low == 0 || vict == nil {
				break
			}
		}
		page_string(ch.Desc, strp, TRUE)
		libc.Free(unsafe.Pointer(strp))
	default:
		send_to_char(ch, libc.CString("Sorry, I don't understand that.\r\n"))
	}
}

type set_struct struct {
	Cmd   *byte
	Level int8
	Pcnpc int8
	Type  int8
}

var set_fields [83]set_struct = [83]set_struct{{Cmd: libc.CString("brief"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("invstart"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("title"), Level: ADMLVL_GOD, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("nosummon"), Level: ADMLVL_GRGOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("maxpl"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("maxki"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("maxst"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("pl"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("ki"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("sta"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("align"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("str"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("stradd"), Level: ADMLVL_IMPL, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("int"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("wis"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("dex"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("con"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("cha"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("ac"), Level: ADMLVL_IMPL, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("zenni"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("bank"), Level: ADMLVL_BUILDER, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("exp"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("accuracy"), Level: ADMLVL_VICE, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("damage"), Level: ADMLVL_VICE, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("invis"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("nohassle"), Level: ADMLVL_VICE, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("frozen"), Level: ADMLVL_GRGOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("practices"), Level: ADMLVL_BUILDER, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("lessons"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("drunk"), Level: ADMLVL_GOD, Pcnpc: BOTH, Type: MISC}, {Cmd: libc.CString("hunger"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: MISC}, {Cmd: libc.CString("thirst"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: MISC}, {Cmd: libc.CString("killer"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("thief"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("level"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("room"), Level: ADMLVL_IMPL, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("roomflag"), Level: ADMLVL_VICE, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("siteok"), Level: ADMLVL_VICE, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("deleted"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("class"), Level: ADMLVL_VICE, Pcnpc: BOTH, Type: MISC}, {Cmd: libc.CString("nowizlist"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("quest"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("loadroom"), Level: ADMLVL_VICE, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("color"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("idnum"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("passwd"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("nodelete"), Level: ADMLVL_GOD, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("sex"), Level: ADMLVL_VICE, Pcnpc: BOTH, Type: MISC}, {Cmd: libc.CString("age"), Level: ADMLVL_VICE, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("height"), Level: ADMLVL_GOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("weight"), Level: ADMLVL_GOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("olc"), Level: ADMLVL_GRGOD, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("race"), Level: ADMLVL_VICE, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("trains"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("feats"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("ethic"), Level: ADMLVL_GOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("unused1"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("unused2"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("adminlevel"), Level: ADMLVL_GRGOD, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("hairl"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("hairs"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("hairc"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("skin"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("eye"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("basepl"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("baseki"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("basest"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("droom"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("absorbs"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("ugp"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("aura"), Level: ADMLVL_IMMORT, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("trp"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("boost"), Level: ADMLVL_GRGOD, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("multi"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("deaths"), Level: ADMLVL_BUILDER, Pcnpc: BOTH, Type: NUMBER}, {Cmd: libc.CString("user"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: MISC}, {Cmd: libc.CString("phase"), Level: ADMLVL_IMMORT, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("racial"), Level: ADMLVL_IMMORT, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("slots"), Level: ADMLVL_IMMORT, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("feature"), Level: ADMLVL_IMMORT, Pcnpc: PC, Type: BINARY}, {Cmd: libc.CString("tclass"), Level: ADMLVL_VICE, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("clones"), Level: ADMLVL_IMPL, Pcnpc: PC, Type: NUMBER}, {Cmd: libc.CString("\n"), Level: 0, Pcnpc: BOTH, Type: MISC}}

func perform_set(ch *char_data, vict *char_data, mode int, val_arg *byte) int {
	var (
		i     int
		on    int   = 0
		off   int   = 0
		value int64 = 0
		rnum  room_rnum
		rvnum room_vnum
	)
	if ch.Admlevel != ADMLVL_IMPL {
		if !IS_NPC(vict) && ch.Admlevel <= vict.Admlevel && vict != ch {
			send_to_char(ch, libc.CString("Maybe that's not such a great idea...\r\n"))
			return 0
		}
	}
	if ch.Admlevel < int(set_fields[mode].Level) {
		send_to_char(ch, libc.CString("You are not godly enough for that!\r\n"))
		return 0
	}
	if IS_NPC(vict) && (int(set_fields[mode].Pcnpc)&NPC) == 0 {
		send_to_char(ch, libc.CString("You can't do that to a beast!\r\n"))
		return 0
	} else if !IS_NPC(vict) && (int(set_fields[mode].Pcnpc)&PC) == 0 {
		send_to_char(ch, libc.CString("That can only be done to a beast!\r\n"))
		return 0
	}
	if int(set_fields[mode].Type) == BINARY {
		if libc.StrCmp(val_arg, libc.CString("on")) == 0 || libc.StrCmp(val_arg, libc.CString("yes")) == 0 {
			on = 1
		} else if libc.StrCmp(val_arg, libc.CString("off")) == 0 || libc.StrCmp(val_arg, libc.CString("no")) == 0 {
			off = 1
		}
		if on == 0 && off == 0 {
			send_to_char(ch, libc.CString("Value must be 'on' or 'off'.\r\n"))
			return 0
		}
		send_to_char(ch, libc.CString("%s %s for %s.\r\n"), set_fields[mode].Cmd, func() string {
			if on != 0 {
				return "ON"
			}
			return "OFF"
		}(), GET_NAME(vict))
	} else if int(set_fields[mode].Type) == NUMBER {
		var ptr *byte
		value = strtoll(val_arg, &ptr, 10)
		send_to_char(ch, libc.CString("%s's %s set to %lld.\r\n"), GET_NAME(vict), set_fields[mode].Cmd, value)
	} else {
		send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
	}
	switch mode {
	case 0:
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_BRIEF/32)] |= bitvector_t(int32(1 << (int(PRF_BRIEF % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_BRIEF/32)] &= bitvector_t(int32(^(1 << (int(PRF_BRIEF % 32)))))
		}
	case 1:
		if on != 0 {
			vict.Act[int(PLR_INVSTART/32)] |= bitvector_t(int32(1 << (int(PLR_INVSTART % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_INVSTART/32)] &= bitvector_t(int32(^(1 << (int(PLR_INVSTART % 32)))))
		}
	case 2:
		set_title(vict, val_arg)
		send_to_char(ch, libc.CString("%s's title is now: %s\r\n"), GET_NAME(vict), GET_TITLE(vict))
	case 3:
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_SUMMONABLE/32)] |= bitvector_t(int32(1 << (int(PRF_SUMMONABLE % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_SUMMONABLE/32)] &= bitvector_t(int32(^(1 << (int(PRF_SUMMONABLE % 32)))))
		}
		send_to_char(ch, libc.CString("Nosummon %s for %s.\r\n"), func() string {
			if on == 0 {
				return "ON"
			}
			return "OFF"
		}(), GET_NAME(vict))
	case 4:
		vict.Max_hit = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set maxpl for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set maxpl for %s."), GET_NAME(ch), GET_NAME(vict))
	case 5:
		vict.Max_mana = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set maxki for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set maxki for %s."), GET_NAME(ch), GET_NAME(vict))
	case 6:
		vict.Max_move = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set maxsta for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set maxsta for %s."), GET_NAME(ch), GET_NAME(vict))
	case 7:
		vict.Hit = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set pl for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set pl for %s."), GET_NAME(ch), GET_NAME(vict))
	case 8:
		vict.Mana = value
		affect_total(vict)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set ki for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set ki for %s."), GET_NAME(ch), GET_NAME(vict))
	case 9:
		vict.Move = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set st for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set st for %s."), GET_NAME(ch), GET_NAME(vict))
	case 10:
		vict.Alignment = int(func() int64 {
			value = int64(MAX(int(-1000), MIN(1000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set align for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set align for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 11:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Str = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set str for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set str for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 12:
		send_to_char(ch, libc.CString("Setting str_add does nothing now.\r\n"))
		fallthrough
	case 13:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Intel = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set intel for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set intel for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 14:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Wis = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set wis for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set wis for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 15:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Dex = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set dex for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set dex for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 16:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Con = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set con for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set con for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 17:
		value = int64(MAX(0, MIN(100, int(value))))
		vict.Real_abils.Cha = int8(value)
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set speed for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set speed for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 18:
		vict.Armor = int(func() int64 {
			value = int64(MAX(int(-100), MIN(500, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set armor index for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set armor index for %s."), GET_NAME(ch), GET_NAME(vict))
		affect_total(vict)
	case 19:
		vict.Gold = int(func() int64 {
			value = int64(MAX(0, MIN(100000000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set zenni for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set zenni for %s."), GET_NAME(ch), GET_NAME(vict))
	case 20:
		vict.Bank_gold = int(func() int64 {
			value = int64(MAX(0, MIN(100000000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set bank for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set bank for %s."), GET_NAME(ch), GET_NAME(vict))
	case 21:
		vict.Exp = func() int64 {
			value = int64(MAX(0, MIN(50000000, int(value))))
			return value
		}()
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set exp for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set exp for %s."), GET_NAME(ch), GET_NAME(vict))
	case 22:
		send_to_char(ch, libc.CString("This does nothing at the moment.\r\n"))
	case 23:
		vict.Damage_mod = int(func() int64 {
			value = int64(MAX(int(-20), MIN(20, int(value))))
			return value
		}())
		affect_total(vict)
	case 24:
		if ch.Admlevel < ADMLVL_IMPL && ch != vict {
			send_to_char(ch, libc.CString("You aren't godly enough for that!\r\n"))
			return 0
		}
		vict.Player_specials.Invis_level = int16(func() int64 {
			value = int64(MAX(0, MIN(vict.Admlevel, int(value))))
			return value
		}())
	case 25:
		if ch.Admlevel < ADMLVL_IMPL && ch != vict {
			send_to_char(ch, libc.CString("You aren't godly enough for that!\r\n"))
			return 0
		}
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_NOHASSLE/32)] |= bitvector_t(int32(1 << (int(PRF_NOHASSLE % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_NOHASSLE/32)] &= bitvector_t(int32(^(1 << (int(PRF_NOHASSLE % 32)))))
		}
	case 26:
		if ch == vict && on != 0 {
			send_to_char(ch, libc.CString("Better not -- could be a long winter!\r\n"))
			return 0
		}
		if on != 0 {
			vict.Act[int(PLR_FROZEN/32)] |= bitvector_t(int32(1 << (int(PLR_FROZEN % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_FROZEN/32)] &= bitvector_t(int32(^(1 << (int(PLR_FROZEN % 32)))))
		}
	case 27:
		fallthrough
	case 28:
		if vict.Level != 0 {
			vict.Player_specials.Class_skill_points[vict.Chclass] = int(func() int64 {
				value = int64(MAX(0, MIN(10000, int(value))))
				return value
			}())
			mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set PS for %s."), GET_NAME(ch), GET_NAME(vict))
			log_imm_action(libc.CString("SET: %s has set PS for %s."), GET_NAME(ch), GET_NAME(vict))
		} else {
			vict.Player_specials.Skill_points = int(func() int64 {
				value = int64(MAX(0, MIN(10000, int(value))))
				return value
			}())
			mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set PS for %s."), GET_NAME(ch), GET_NAME(vict))
			log_imm_action(libc.CString("SET: %s has set PS for %s."), GET_NAME(ch), GET_NAME(vict))
		}
	case 29:
		fallthrough
	case 30:
		fallthrough
	case 31:
		if libc.StrCaseCmp(val_arg, libc.CString("off")) == 0 {
			vict.Player_specials.Conditions[mode-29] = -1
			send_to_char(ch, libc.CString("%s's %s now off.\r\n"), GET_NAME(vict), set_fields[mode].Cmd)
		} else if is_number(val_arg) != 0 {
			value = int64(libc.Atoi(libc.GoString(val_arg)))
			value = int64(MAX(0, MIN(24, int(value))))
			vict.Player_specials.Conditions[mode-29] = int8(value)
			send_to_char(ch, libc.CString("%s's %s set to %lld.\r\n"), GET_NAME(vict), set_fields[mode].Cmd, value)
		} else {
			send_to_char(ch, libc.CString("Must be 'off' or a value from 0 to 24.\r\n"))
			return 0
		}
	case 32:
		if on != 0 {
			vict.Act[int(PLR_KILLER/32)] |= bitvector_t(int32(1 << (int(PLR_KILLER % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_KILLER/32)] &= bitvector_t(int32(^(1 << (int(PLR_KILLER % 32)))))
		}
	case 33:
		if on != 0 {
			vict.Act[int(PLR_THIEF/32)] |= bitvector_t(int32(1 << (int(PLR_THIEF % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_THIEF/32)] &= bitvector_t(int32(^(1 << (int(PLR_THIEF % 32)))))
		}
	case 34:
		if !IS_NPC(vict) && value > 100 {
			send_to_char(ch, libc.CString("You can't do that.\r\n"))
			return 0
		}
		value = int64(MAX(0, int(value)))
		vict.Level = int(value)
	case 35:
		if (func() room_rnum {
			rnum = real_room(room_vnum(value))
			return rnum
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("No room exists with that number.\r\n"))
			return 0
		}
		if vict.In_room != room_rnum(-1) {
			char_from_room(vict)
		}
		char_to_room(vict, rnum)
	case 36:
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_ROOMFLAGS/32)] |= bitvector_t(int32(1 << (int(PRF_ROOMFLAGS % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_ROOMFLAGS/32)] &= bitvector_t(int32(^(1 << (int(PRF_ROOMFLAGS % 32)))))
		}
	case 37:
		if on != 0 {
			vict.Act[int(PLR_SITEOK/32)] |= bitvector_t(int32(1 << (int(PLR_SITEOK % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_SITEOK/32)] &= bitvector_t(int32(^(1 << (int(PLR_SITEOK % 32)))))
		}
	case 38:
		if on != 0 {
			vict.Act[int(PLR_DELETED/32)] |= bitvector_t(int32(1 << (int(PLR_DELETED % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_DELETED/32)] &= bitvector_t(int32(^(1 << (int(PLR_DELETED % 32)))))
		}
	case 39:
		if (func() int {
			i = search_block(val_arg, &class_names[0], FALSE)
			return i
		}()) < 0 {
			send_to_char(ch, libc.CString("That is not a class.\r\n"))
			return 0
		}
		value = int64((vict.Chclasses[vict.Chclass]) + (vict.Epicclasses[vict.Chclass]))
		vict.Chclass = int8(i)
	case 40:
		if on != 0 {
			vict.Act[int(PLR_NOWIZLIST/32)] |= bitvector_t(int32(1 << (int(PLR_NOWIZLIST % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_NOWIZLIST/32)] &= bitvector_t(int32(^(1 << (int(PLR_NOWIZLIST % 32)))))
		}
	case 41:
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_QUEST/32)] |= bitvector_t(int32(1 << (int(PRF_QUEST % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_QUEST/32)] &= bitvector_t(int32(^(1 << (int(PRF_QUEST % 32)))))
		}
	case 42:
		if libc.StrCaseCmp(val_arg, libc.CString("off")) == 0 {
			vict.Act[int(PLR_LOADROOM/32)] &= bitvector_t(int32(^(1 << (int(PLR_LOADROOM % 32)))))
			vict.Player_specials.Load_room = -1
		} else if is_number(val_arg) != 0 {
			rvnum = room_vnum(libc.Atoi(libc.GoString(val_arg)))
			if real_room(rvnum) != room_rnum(-1) {
				vict.Act[int(PLR_LOADROOM/32)] |= bitvector_t(int32(1 << (int(PLR_LOADROOM % 32))))
				vict.Player_specials.Load_room = rvnum
				send_to_char(ch, libc.CString("%s will enter at room #%d.\r\n"), GET_NAME(vict), vict.Player_specials.Load_room)
			} else {
				send_to_char(ch, libc.CString("That room does not exist!\r\n"))
				return 0
			}
		} else {
			send_to_char(ch, libc.CString("Must be 'off' or a room's virtual number.\r\n"))
			return 0
		}
	case 43:
		if on != 0 {
			vict.Player_specials.Pref[int(PRF_COLOR/32)] |= bitvector_t(int32(1 << (int(PRF_COLOR % 32))))
		} else if off != 0 {
			vict.Player_specials.Pref[int(PRF_COLOR/32)] &= bitvector_t(int32(^(1 << (int(PRF_COLOR % 32)))))
		}
	case 44:
		if int(ch.Idnum) == 0 || IS_NPC(vict) {
			return 0
		}
		vict.Idnum = int32(value)
	case 45:
		if int(ch.Idnum) > 1 {
			send_to_char(ch, libc.CString("Please don't use this command, yet.\r\n"))
			return 0
		}
		if ch.Admlevel < 10 {
			send_to_char(ch, libc.CString("NO.\r\n"))
			return 0
		}
	case 46:
		if on != 0 {
			vict.Act[int(PLR_NODELETE/32)] |= bitvector_t(int32(1 << (int(PLR_NODELETE % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_NODELETE/32)] &= bitvector_t(int32(^(1 << (int(PLR_NODELETE % 32)))))
		}
	case 47:
		if (func() int {
			i = search_block(val_arg, &genders[0], FALSE)
			return i
		}()) < 0 {
			send_to_char(ch, libc.CString("Must be 'male', 'female', or 'neutral'.\r\n"))
			return 0
		}
		vict.Sex = int8(i)
	case 48:
		if value < 2 || value > 20000 {
			send_to_char(ch, libc.CString("Ages 2 to 20000 accepted.\r\n"))
			return 0
		}
		vict.Time.Birth = libc.Time(int64(libc.GetTime(nil)) - value*int64(((int(SECS_PER_MUD_HOUR*24))*30)*12))
	case 49:
		vict.Height = uint8(int8(value))
		affect_total(vict)
	case 50:
		vict.Weight = uint8(int8(value))
		affect_total(vict)
	case 51:
		if is_abbrev(val_arg, libc.CString("socials")) != 0 || is_abbrev(val_arg, libc.CString("actions")) != 0 {
			vict.Player_specials.Olc_zone = AEDIT_PERMISSION
		} else if is_abbrev(val_arg, libc.CString("hedit")) != 0 {
			vict.Player_specials.Olc_zone = HEDIT_PERMISSION
		} else if is_abbrev(val_arg, libc.CString("off")) != 0 {
			vict.Player_specials.Olc_zone = -1
		} else if is_number(val_arg) == 0 {
			send_to_char(ch, libc.CString("Value must be either 'socials', 'actions', 'hedit', 'off' or a zone number.\r\n"))
			return 0
		} else {
			vict.Player_specials.Olc_zone = libc.Atoi(libc.GoString(val_arg))
		}
	case 52:
		if (func() int {
			i = search_block(val_arg, &race_names[0], FALSE)
			return i
		}()) < 0 {
			send_to_char(ch, libc.CString("That is not a race.\r\n"))
			return 0
		}
		vict.Race = int8(i)
		racial_body_parts(vict)
	case 53:
		vict.Player_specials.Ability_trains = int(func() int64 {
			value = int64(MAX(0, MIN(500, int(value))))
			return value
		}())
	case 54:
		vict.Player_specials.Feat_points = int(func() int64 {
			value = int64(MAX(0, MIN(500, int(value))))
			return value
		}())
	case 55:
		vict.Alignment_ethic = int(func() int64 {
			value = int64(MAX(int(-1000), MIN(1000, int(value))))
			return value
		}())
		affect_total(vict)
	case 56:
		vict.Max_ki = func() int64 {
			value = int64(MAX(1, MIN(5000, int(value))))
			return value
		}()
		affect_total(vict)
	case 57:
		vict.Ki = func() int64 {
			value = int64(MAX(0, MIN(int(vict.Max_ki), int(value))))
			return value
		}()
		affect_total(vict)
	case 58:
		if vict.Admlevel >= ch.Admlevel && vict != ch {
			send_to_char(ch, libc.CString("Permission denied.\r\n"))
			return 0
		}
		if value < ADMLVL_NONE || value > int64(ch.Admlevel) {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		if vict.Admlevel == int(value) {
			return 1
		}
		admin_set(vict, int(value))
	case 59:
		if value < 0 || value >= 6 {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		vict.Hairl = int8(value)
		return 1
	case 60:
		if value < 0 || value >= 13 {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		vict.Hairs = int8(value)
		return 1
	case 61:
		if value < 0 || value >= 15 {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		vict.Hairc = int8(value)
		return 1
	case 62:
		if value < 0 || value >= 12 {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		vict.Skin = int8(value)
		return 1
	case 63:
		if value < 0 || value >= 13 {
			send_to_char(ch, libc.CString("You can't set it to that.\r\n"))
			return 0
		}
		vict.Eye = int8(value)
		return 1
	case 64:
		vict.Basepl = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set basepl for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set basepl for %s."), GET_NAME(ch), GET_NAME(vict))
	case 65:
		vict.Baseki = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set baseki for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set baseki for %s."), GET_NAME(ch), GET_NAME(vict))
	case 66:
		vict.Basest = value
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set basest for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set basest for %s."), GET_NAME(ch), GET_NAME(vict))
	case 67:
		vict.Droom = room_vnum(func() int64 {
			value = int64(MAX(0, MIN(20000, int(value))))
			return value
		}())
	case 68:
		vict.Absorbs = int(func() int64 {
			value = int64(MAX(0, MIN(3, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set absorbs for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set absorbs for %s."), GET_NAME(ch), GET_NAME(vict))
	case 69:
		vict.Upgrade += int(func() int64 {
			value = int64(MAX(1, MIN(1000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set upgrade points for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set upgrade points for %s."), GET_NAME(ch), GET_NAME(vict))
	case 70:
		vict.Aura = int(func() int64 {
			value = int64(MAX(0, MIN(8, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set aura for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set aura for %s."), GET_NAME(ch), GET_NAME(vict))
	case 71:
		send_to_char(ch, libc.CString("Use the reward command.\r\n"))
	case 72:
		vict.Boosts = int(func() int64 {
			value = int64(MAX(int(-1000), MIN(1000, int(value))))
			return value
		}())
	case 73:
		if on != 0 {
			vict.Act[int(PLR_MULTP/32)] |= bitvector_t(int32(1 << (int(PLR_MULTP % 32))))
		} else if off != 0 {
			vict.Act[int(PLR_MULTP/32)] &= bitvector_t(int32(^(1 << (int(PLR_MULTP % 32)))))
		}
	case 74:
		vict.Dcount = int(func() int64 {
			value = int64(MAX(int(-1000), MIN(1000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set death count for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set death count for %s."), GET_NAME(ch), GET_NAME(vict))
	case 75:
		send_to_char(ch, libc.CString("No."))
	case 76:
		if vict.Desc != nil {
			star_phase(vict, int(func() int64 {
				value = int64(MAX(0, MIN(2, int(value))))
				return value
			}()))
		} else {
			send_to_char(ch, libc.CString("They aren't even in the game!\r\n"))
		}
	case 77:
		vict.Player_specials.Racial_pref = int(func() int64 {
			value = int64(MAX(1, MIN(3, int(value))))
			return value
		}())
	case 78:
		vict.Skill_slots = int(func() int64 {
			value = int64(MAX(1, MIN(1000, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set skill slots for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set skill slots for %s."), GET_NAME(ch), GET_NAME(vict))
	case 79:
		vict.Feature = (*byte)(unsafe.Pointer(uintptr('\x00')))
	case 80:
		vict.Transclass = int(func() int64 {
			value = int64(MAX(1, MIN(3, int(value))))
			return value
		}())
		mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("SET: %s has set transformation class for %s."), GET_NAME(ch), GET_NAME(vict))
		log_imm_action(libc.CString("SET: %s has set transformation class for %s."), GET_NAME(ch), GET_NAME(vict))
	case 81:
		vict.Clones = int16(func() int64 {
			value = int64(MAX(1, MIN(3, int(value))))
			return value
		}())
		send_to_char(ch, libc.CString("Done.\r\n"))
	default:
		send_to_char(ch, libc.CString("Can't set that!\r\n"))
		return 0
	}
	return 1
}
func do_set(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict      *char_data = nil
		cbuf      *char_data = nil
		field     [2048]byte
		name      [2048]byte
		buf       [2048]byte
		mode      int
		len_      int
		player_i  int = 0
		retval    int
		is_file   int8 = 0
		is_player int8 = 0
	)
	half_chop(argument, &name[0], &buf[0])
	if libc.StrCmp(&name[0], libc.CString("file")) == 0 {
		is_file = 1
		half_chop(&buf[0], &name[0], &buf[0])
	} else if libc.StrCaseCmp(&name[0], libc.CString("player")) == 0 {
		is_player = 1
		half_chop(&buf[0], &name[0], &buf[0])
	} else if libc.StrCaseCmp(&name[0], libc.CString("mob")) == 0 {
		half_chop(&buf[0], &name[0], &buf[0])
	}
	half_chop(&buf[0], &field[0], &buf[0])
	if name[0] == 0 || field[0] == 0 {
		send_to_char(ch, libc.CString("Usage: set <victim> <field> <value>\r\n"))
		return
	}
	if int(is_file) == 0 {
		if int(is_player) != 0 {
			if (func() *char_data {
				vict = get_player_vis(ch, &name[0], nil, 1<<1)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("There is no such player.\r\n"))
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<1)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("There is no such creature.\r\n"))
				return
			}
		}
	} else if int(is_file) != 0 {
		cbuf = new(char_data)
		clear_char(cbuf)
		cbuf.Player_specials = new(player_special_data)
		if (func() int {
			player_i = load_char(&name[0], cbuf)
			return player_i
		}()) > -1 {
			if cbuf.Admlevel >= ch.Admlevel {
				free_char(cbuf)
				send_to_char(ch, libc.CString("Sorry, you can't do that.\r\n"))
				return
			}
			vict = cbuf
		} else {
			free_char(cbuf)
			send_to_char(ch, libc.CString("There is no such player.\r\n"))
			return
		}
	}
	len_ = libc.StrLen(&field[0])
	for mode = 0; *set_fields[mode].Cmd != '\n'; mode++ {
		if libc.StrCmp(&field[0], set_fields[mode].Cmd) == 0 {
			break
		}
	}
	if *set_fields[mode].Cmd == '\n' {
		for mode = 0; *set_fields[mode].Cmd != '\n'; mode++ {
			if libc.StrNCmp(&field[0], set_fields[mode].Cmd, len_) == 0 {
				break
			}
		}
	}
	retval = perform_set(ch, vict, mode, &buf[0])
	if retval != 0 {
		if int(is_file) == 0 && !IS_NPC(vict) {
			save_char(vict)
		}
		if int(is_file) != 0 {
			cbuf.Pfilepos = player_i
			save_char(cbuf)
			send_to_char(ch, libc.CString("Saved in file.\r\n"))
		}
	}
	if int(is_file) != 0 {
		free_char(cbuf)
	}
}
func do_saveall(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Admlevel < ADMLVL_BUILDER {
		send_to_char(ch, libc.CString("You are not holy enough to use this privelege.\n\r"))
	} else {
		save_all()
		House_save_all()
		send_to_char(ch, libc.CString("World and house files saved.\n\r"))
	}
}
func do_plist(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i           int
		len_        int = 0
		count       int = 0
		mode        int8
		buf         [64936]byte
		name_search [20]byte
		time_str    [64936]byte
		time_away   time_info_data
		low         int = 0
		high        int = 100
		low_day     int = 0
		high_day    int = 10000
		low_hr      int = 0
		high_hr     int = 24
	)
	skip_spaces(&argument)
	libc.StrCpy(&buf[0], argument)
	name_search[0] = '\x00'
	for buf[0] != 0 {
		var (
			arg  [2048]byte
			buf1 [2048]byte
		)
		half_chop(&buf[0], &arg[0], &buf1[0])
		if unicode.IsDigit(rune(arg[0])) {
			if stdio.Sscanf(&arg[0], "%d-%d", &low, &high) == 1 {
				high = low
			}
			libc.StrCpy(&buf[0], &buf1[0])
		} else if arg[0] == '-' {
			mode = int8(arg[1])
			switch mode {
			case 'l':
				half_chop(&buf1[0], &arg[0], &buf[0])
				stdio.Sscanf(&arg[0], "%d-%d", &low, &high)
			case 'n':
				half_chop(&buf1[0], &name_search[0], &buf[0])
			case 'i':
				libc.StrCpy(&buf[0], &buf1[0])
				low = 1
			case 'm':
				libc.StrCpy(&buf[0], &buf1[0])
				high = 100
			case 'd':
				half_chop(&buf1[0], &arg[0], &buf[0])
				if stdio.Sscanf(&arg[0], "%d-%d", &low_day, &high_day) == 1 {
					high_day = low_day
				}
			case 'h':
				half_chop(&buf1[0], &arg[0], &buf[0])
				if stdio.Sscanf(&arg[0], "%d-%d", &low_hr, &high_hr) == 1 {
					high_hr = low_hr
				}
			default:
				send_to_char(ch, libc.CString("%s\r\n"), PLIST_FORMAT)
				return
			}
		} else {
			send_to_char(ch, libc.CString("%s\r\n"), PLIST_FORMAT)
			return
		}
	}
	len_ = 0
	len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "@W[ Id] (Lv) Name         Last@n\r\n%s-----------------------------------------------%s\r\n", func() string {
		if (func() int {
			if !IS_NPC(ch) {
				if PRF_FLAGGED(ch, PRF_COLOR) {
					return 1
				}
				return 0
			}
			return 0
		}()) >= C_ON {
			return KCYN
		}
		return KNUL
	}(), func() string {
		if (func() int {
			if !IS_NPC(ch) {
				if PRF_FLAGGED(ch, PRF_COLOR) {
					return 1
				}
				return 0
			}
			return 0
		}()) >= C_ON {
			return KNRM
		}
		return KNUL
	}())
	for i = 0; i <= top_of_p_table; i++ {
		if ((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Flags & (1 << 0)) != 0 {
			len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[%3ld] <DELETED> --Will be removed next boot.\r\n", (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id)
			continue
		}
		if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level < low || (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level > high {
			continue
		}
		time_away = *real_time_passed(libc.GetTime(nil), (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Last)
		if name_search[0] != 0 && libc.StrCaseCmp(&name_search[0], (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name) != 0 {
			continue
		}
		if time_away.Day > high_day || time_away.Day < low_day {
			continue
		}
		if time_away.Hours > high_hr || time_away.Hours < low_hr {
			continue
		}
		libc.StrCpy(&time_str[0], libc.AscTime(libc.LocalTime(&(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Last)))
		time_str[libc.StrLen(&time_str[0])-1] = '\x00'
		len_ += stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "[%3ld] (%2d) %-12s %s\r\n", (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level, CAP(libc.StrDup((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name)), &time_str[0])
		count++
	}
	stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%s-----------------------------------------------%s\r\n%d players listed.\r\n", func() string {
		if (func() int {
			if !IS_NPC(ch) {
				if PRF_FLAGGED(ch, PRF_COLOR) {
					return 1
				}
				return 0
			}
			return 0
		}()) >= C_ON {
			return KCYN
		}
		return KNUL
	}(), func() string {
		if (func() int {
			if !IS_NPC(ch) {
				if PRF_FLAGGED(ch, PRF_COLOR) {
					return 1
				}
				return 0
			}
			return 0
		}()) >= C_ON {
			return KNRM
		}
		return KNUL
	}(), count)
	page_string(ch.Desc, &buf[0], TRUE)
}
func do_peace(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict   *char_data
		next_v *char_data
	)
	send_to_room(ch.In_room, libc.CString("Everything is quite peaceful now.\r\n"))
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict.Admlevel > ch.Admlevel {
			continue
		}
		stop_fighting(vict)
		vict.Position = POS_SITTING
	}
	stop_fighting(ch)
	ch.Position = POS_STANDING
}
func do_wizupdate(ch *char_data, argument *byte, cmd int, subcmd int) {
	run_autowiz()
	send_to_char(ch, libc.CString("Wizlists updated.\n\r"))
}
func do_raise(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		name [2048]byte
	)
	one_argument(argument, &name[0])
	if ch.Admlevel < ADMLVL_BUILDER && !IS_NPC(ch) {
		return
	}
	if (func() *char_data {
		vict = get_player_vis(ch, &name[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("There is no such player.\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("Sorry, only players get spirits.\r\n"))
		return
	}
	if !AFF_FLAGGED(vict, AFF_SPIRIT) {
		send_to_char(ch, libc.CString("But they aren't even dead!\r\n"))
		return
	}
	send_to_char(ch, libc.CString("@wYou return %s from the @Bspirit@w world, to the world of the living!@n\r\n"), GET_NAME(vict))
	send_to_char(vict, libc.CString("@wYour @Bspirit@w has been returned to the world of the living by %s!@n\r\n"), GET_NAME(ch))
	vict.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
	vict.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
	send_to_imm(libc.CString("Log: %s has raised %s from the dead."), GET_NAME(ch), GET_NAME(vict))
	log_imm_action(libc.CString("RAISE: %s has raised %s from the dead."), GET_NAME(ch), GET_NAME(vict))
	var statpunish int = FALSE
	if ch.Admlevel <= 0 {
		statpunish = TRUE
	}
	if vict.Hit < 1 {
		vict.Hit = gear_pl(vict)
	}
	vict.Mana = vict.Max_mana
	vict.Move = vict.Max_move
	vict.Limb_condition[1] = 100
	vict.Limb_condition[2] = 100
	vict.Limb_condition[3] = 100
	vict.Limb_condition[4] = 100
	vict.Act[int(PLR_HEAD/32)] |= bitvector_t(int32(1 << (int(PLR_HEAD % 32))))
	vict.Act[int(PLR_PDEATH/32)] &= bitvector_t(int32(^(1 << (int(PLR_PDEATH % 32)))))
	char_from_room(vict)
	if vict.Droom != room_vnum(-1) && vict.Droom != 0 && vict.Droom != 1 {
		char_to_room(vict, real_room(vict.Droom))
	} else if int(vict.Chclass) == CLASS_ROSHI {
		char_to_room(vict, real_room(1130))
	} else if int(vict.Chclass) == CLASS_KABITO {
		char_to_room(vict, real_room(0x2F42))
	} else if int(vict.Chclass) == CLASS_NAIL {
		char_to_room(vict, real_room(0x2DA3))
	} else if int(vict.Chclass) == CLASS_BARDOCK {
		char_to_room(vict, real_room(2268))
	} else if int(vict.Chclass) == CLASS_KRANE {
		char_to_room(vict, real_room(0x32D1))
	} else if int(vict.Chclass) == CLASS_TAPION {
		char_to_room(vict, real_room(8231))
	} else if int(vict.Chclass) == CLASS_PICCOLO {
		char_to_room(vict, real_room(1659))
	} else if int(vict.Chclass) == CLASS_ANDSIX {
		char_to_room(vict, real_room(1713))
	} else if int(vict.Chclass) == CLASS_DABURA {
		char_to_room(vict, real_room(6486))
	} else if int(vict.Chclass) == CLASS_FRIEZA {
		char_to_room(vict, real_room(4282))
	} else if int(vict.Chclass) == CLASS_GINYU {
		char_to_room(vict, real_room(4289))
	} else if int(vict.Chclass) == CLASS_JINTO {
		char_to_room(vict, real_room(3499))
	} else if int(vict.Chclass) == CLASS_TSUNA {
		char_to_room(vict, real_room(15000))
	} else if int(vict.Chclass) == CLASS_KURZAK {
		char_to_room(vict, real_room(16100))
	} else {
		char_to_room(vict, real_room(300))
		send_to_imm(libc.CString("ERROR: Player %s without acceptable sensei.\r\n"), GET_NAME(vict))
	}
	look_at_room(vict.In_room, vict, 0)
	var losschance int = axion_dice(0)
	if GET_LEVEL(vict) > 9 && statpunish == TRUE {
		send_to_char(vict, libc.CString("@RThe the strain of this type of revival has caused you to be in a weakened state for 100 hours (Game time)! Strength, constitution, wisdom, intelligence, speed, and agility have been reduced by 8 points for the duration.@n\r\n"))
		var str int = -8
		var intel int = -8
		var wis int = -8
		var spd int = -8
		var con int = -8
		var agl int = -8
		var dur int = 100
		if vict.Dcount >= 8 && vict.Dcount < 10 {
			dur = 90
		} else if vict.Dcount >= 5 && vict.Dcount < 8 {
			dur = 75
		} else if vict.Dcount >= 3 && vict.Dcount < 5 {
			dur = 60
		} else if vict.Dcount >= 1 && vict.Dcount < 3 {
			dur = 40
		}
		if int(vict.Real_abils.Intel) <= 16 {
			intel = -4
		}
		if int(vict.Real_abils.Cha) <= 16 {
			spd = -4
		}
		if int(vict.Real_abils.Dex) <= 16 {
			agl = -4
		}
		if int(vict.Real_abils.Wis) <= 16 {
			wis = -4
		}
		if int(vict.Real_abils.Con) <= 16 {
			con = -4
		}
		assign_affect(vict, AFF_WEAKENED_STATE, SKILL_WARP, dur, str, con, intel, agl, wis, spd)
		if losschance >= 100 {
			var psloss int = rand_number(100, 300)
			vict.Player_specials.Class_skill_points[vict.Chclass] -= psloss
			send_to_char(vict, libc.CString("@R...and a loss of @r%d@R PS!@n"), psloss)
			if (vict.Player_specials.Class_skill_points[vict.Chclass]) < 0 {
				vict.Player_specials.Class_skill_points[vict.Chclass] = 0
			}
		}
	}
	ch.Deathtime = 0
	act(libc.CString("$n's body forms in a pool of @Bblue light@n."), TRUE, vict, nil, nil, TO_ROOM)
}
func do_chown(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		victim *char_data
		obj    *obj_data
		buf2   [80]byte
		buf3   [80]byte
		i      int
		k      int = 0
	)
	two_arguments(argument, &buf2[0], &buf3[0])
	if buf2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: chown <object> <character>.\r\n"))
	} else if (func() *char_data {
		victim = get_char_vis(ch, &buf3[0], nil, 1<<1)
		return victim
	}()) == nil {
		send_to_char(ch, libc.CString("No one by that name here.\r\n"))
	} else if victim == ch {
		send_to_char(ch, libc.CString("Are you sure you're feeling ok?\r\n"))
	} else if GET_LEVEL(victim) >= GET_LEVEL(ch) {
		send_to_char(ch, libc.CString("That's really not such a good idea.\r\n"))
	} else if buf3[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: chown <object> <character>.\r\n"))
	} else {
		for i = 0; i < NUM_WEARS; i++ {
			if (victim.Equipment[i]) != nil && CAN_SEE_OBJ(ch, victim.Equipment[i]) && isname(&buf2[0], (victim.Equipment[i]).Name) != 0 {
				obj_to_char(unequip_char(victim, i), victim)
				k = 1
			}
		}
		if (func() *obj_data {
			obj = get_obj_in_list_vis(victim, &buf2[0], nil, victim.Carrying)
			return obj
		}()) == nil {
			if k == 0 && (func() *obj_data {
				obj = get_obj_in_list_vis(victim, &buf2[0], nil, victim.Carrying)
				return obj
			}()) == nil {
				send_to_char(ch, libc.CString("%s does not appear to have the %s.\r\n"), GET_NAME(victim), &buf2[0])
				return
			}
		}
		act(libc.CString("@n$n makes a magical gesture and $p@n flies from $N to $m."), FALSE, ch, obj, unsafe.Pointer(victim), TO_NOTVICT)
		act(libc.CString("@n$n makes a magical gesture and $p@n flies away from you to $m."), FALSE, ch, obj, unsafe.Pointer(victim), TO_VICT)
		act(libc.CString("@nYou make a magical gesture and $p@n flies away from $N to you."), FALSE, ch, obj, unsafe.Pointer(victim), TO_CHAR)
		obj_from_char(obj)
		obj_to_char(obj, ch)
		save_char(ch)
		save_char(victim)
	}
}
func do_zpurge(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj      *obj_data
		next_obj *obj_data
		mob      *char_data
		next_mob *char_data
		i        int
		stored   int = -1
		zone     int
		found    int = FALSE
		room     int
		arg      [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		zone = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone)))).Number)
	} else {
		zone = libc.Atoi(libc.GoString(&arg[0]))
	}
	for i = 0; i <= int(top_of_zone_table) && found == 0; i++ {
		if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number == zone_vnum(zone) {
			stored = i
			found = TRUE
		}
	}
	if found == 0 || can_edit_zone(ch, zone_rnum(zone)) == 0 {
		send_to_char(ch, libc.CString("You cannot purge that zone. Try %d.\r\n"), ch.Player_specials.Olc_zone)
		return
	}
	for room = int(genolc_zone_bottom(zone_rnum(stored))); room <= int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(stored)))).Top); room++ {
		if (func() int {
			i = int(real_room(room_vnum(room)))
			return i
		}()) != int(-1) {
			for mob = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).People; mob != nil; mob = next_mob {
				next_mob = mob.Next_in_room
				if IS_NPC(mob) {
					extract_char(mob)
				}
			}
			for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Contents; obj != nil; obj = next_obj {
				next_obj = obj.Next_content
				extract_obj(obj)
			}
		}
	}
	send_to_char(ch, libc.CString("All mobiles and objects in zone %d purged.\r\n"), zone)
	mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s has purged zone %d."), GET_NAME(ch), zone)
}

type zcheck_armor struct {
	Bitvector  bitvector_t
	Ac_allowed int
	Message    *byte
}

var zarmor [17]zcheck_armor = [17]zcheck_armor{{Bitvector: ITEM_WEAR_FINGER, Ac_allowed: 10, Message: libc.CString("Ring")}, {Bitvector: ITEM_WEAR_NECK, Ac_allowed: 10, Message: libc.CString("Necklace")}, {Bitvector: ITEM_WEAR_BODY, Ac_allowed: 10, Message: libc.CString("Body armor")}, {Bitvector: ITEM_WEAR_HEAD, Ac_allowed: 10, Message: libc.CString("Head gear")}, {Bitvector: ITEM_WEAR_LEGS, Ac_allowed: 10, Message: libc.CString("Legwear")}, {Bitvector: ITEM_WEAR_FEET, Ac_allowed: 10, Message: libc.CString("Footwear")}, {Bitvector: ITEM_WEAR_HANDS, Ac_allowed: 10, Message: libc.CString("Glove")}, {Bitvector: ITEM_WEAR_ARMS, Ac_allowed: 10, Message: libc.CString("Armwear")}, {Bitvector: ITEM_WEAR_SHIELD, Ac_allowed: 10, Message: libc.CString("Shield")}, {Bitvector: ITEM_WEAR_ABOUT, Ac_allowed: 10, Message: libc.CString("Cloak")}, {Bitvector: ITEM_WEAR_WAIST, Ac_allowed: 10, Message: libc.CString("Belt")}, {Bitvector: ITEM_WEAR_WRIST, Ac_allowed: 10, Message: libc.CString("Wristwear")}, {Bitvector: ITEM_WEAR_HOLD, Ac_allowed: 10, Message: libc.CString("Held item")}, {Bitvector: ITEM_WEAR_PACK, Ac_allowed: 10, Message: libc.CString("Backpack item")}, {Bitvector: ITEM_WEAR_EAR, Ac_allowed: 10, Message: libc.CString("Earring item")}, {Bitvector: ITEM_WEAR_SH, Ac_allowed: 10, Message: libc.CString("Shoulder item")}, {Bitvector: ITEM_WEAR_EYE, Ac_allowed: 10, Message: libc.CString("Eye item")}}

type zcheck_affs struct {
	Aff_type int
	Min_aff  int
	Max_aff  int
	Message  *byte
}

var zaffs [45]zcheck_affs = [45]zcheck_affs{{Aff_type: APPLY_NONE, Min_aff: 0, Max_aff: -99, Message: libc.CString("unused0")}, {Aff_type: APPLY_STR, Min_aff: -6, Max_aff: 6, Message: libc.CString("strength")}, {Aff_type: APPLY_DEX, Min_aff: -6, Max_aff: 6, Message: libc.CString("dexterity")}, {Aff_type: APPLY_INT, Min_aff: -6, Max_aff: 6, Message: libc.CString("intelligence")}, {Aff_type: APPLY_WIS, Min_aff: -6, Max_aff: 6, Message: libc.CString("wisdom")}, {Aff_type: APPLY_CON, Min_aff: -6, Max_aff: 6, Message: libc.CString("constitution")}, {Aff_type: APPLY_CHA, Min_aff: -6, Max_aff: 6, Message: libc.CString("charisma")}, {Aff_type: APPLY_CLASS, Min_aff: 0, Max_aff: 0, Message: libc.CString("class")}, {Aff_type: APPLY_LEVEL, Min_aff: 0, Max_aff: 0, Message: libc.CString("level")}, {Aff_type: APPLY_AGE, Min_aff: -10, Max_aff: 10, Message: libc.CString("age")}, {Aff_type: APPLY_CHAR_WEIGHT, Min_aff: -50, Max_aff: 50, Message: libc.CString("character weight")}, {Aff_type: APPLY_CHAR_HEIGHT, Min_aff: -50, Max_aff: 50, Message: libc.CString("character height")}, {Aff_type: APPLY_MANA, Min_aff: -50, Max_aff: 50, Message: libc.CString("mana")}, {Aff_type: APPLY_HIT, Min_aff: -50, Max_aff: 50, Message: libc.CString("hit points")}, {Aff_type: APPLY_MOVE, Min_aff: -50, Max_aff: 50, Message: libc.CString("movement")}, {Aff_type: APPLY_GOLD, Min_aff: 0, Max_aff: 0, Message: libc.CString("gold")}, {Aff_type: APPLY_EXP, Min_aff: 0, Max_aff: 0, Message: libc.CString("experience")}, {Aff_type: APPLY_AC, Min_aff: -10, Max_aff: 10, Message: libc.CString("magical AC")}, {Aff_type: APPLY_ACCURACY, Min_aff: 0, Max_aff: -99, Message: libc.CString("accuracy")}, {Aff_type: APPLY_DAMAGE, Min_aff: 0, Max_aff: -99, Message: libc.CString("damage")}, {Aff_type: APPLY_REGEN, Min_aff: 0, Max_aff: 0, Message: libc.CString("regen")}, {Aff_type: APPLY_TRAIN, Min_aff: 0, Max_aff: 0, Message: libc.CString("train")}, {Aff_type: APPLY_LIFEMAX, Min_aff: 0, Max_aff: 0, Message: libc.CString("lifemax")}, {Aff_type: APPLY_UNUSED3, Min_aff: 0, Max_aff: 0, Message: libc.CString("unused")}, {Aff_type: APPLY_UNUSED4, Min_aff: 0, Max_aff: 0, Message: libc.CString("unused")}, {Aff_type: APPLY_RACE, Min_aff: 0, Max_aff: 0, Message: libc.CString("race")}, {Aff_type: APPLY_TURN_LEVEL, Min_aff: -6, Max_aff: 6, Message: libc.CString("turn level")}, {Aff_type: APPLY_SPELL_LVL_0, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 0")}, {Aff_type: APPLY_SPELL_LVL_1, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 1")}, {Aff_type: APPLY_SPELL_LVL_2, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 2")}, {Aff_type: APPLY_SPELL_LVL_3, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 3")}, {Aff_type: APPLY_SPELL_LVL_4, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 4")}, {Aff_type: APPLY_SPELL_LVL_5, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 5")}, {Aff_type: APPLY_SPELL_LVL_6, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 6")}, {Aff_type: APPLY_SPELL_LVL_7, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 7")}, {Aff_type: APPLY_SPELL_LVL_8, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 8")}, {Aff_type: APPLY_SPELL_LVL_9, Min_aff: 0, Max_aff: 0, Message: libc.CString("spell level 9")}, {Aff_type: APPLY_KI, Min_aff: 0, Max_aff: 0, Message: libc.CString("ki")}, {Aff_type: APPLY_FORTITUDE, Min_aff: -4, Max_aff: 4, Message: libc.CString("fortitude")}, {Aff_type: APPLY_REFLEX, Min_aff: -4, Max_aff: 4, Message: libc.CString("reflex")}, {Aff_type: APPLY_WILL, Min_aff: -4, Max_aff: 4, Message: libc.CString("will")}, {Aff_type: APPLY_SKILL, Min_aff: -10, Max_aff: 10, Message: libc.CString("skill")}, {Aff_type: APPLY_FEAT, Min_aff: -10, Max_aff: 10, Message: libc.CString("feat")}, {Aff_type: APPLY_ALLSAVES, Min_aff: -4, Max_aff: 4, Message: libc.CString("all 3 save types")}, {Aff_type: APPLY_RESISTANCE, Min_aff: -4, Max_aff: 4, Message: libc.CString("resistance")}}
var offlimit_zones [5]int = [5]int{0, 12, 13, 14, -1}

func do_zcheck(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		zrnum   zone_rnum
		obj     *obj_data
		mob     *char_data = nil
		exroom  room_vnum  = 0
		ac      int        = 0
		affs    int        = 0
		tohit   int
		todam   int
		value   int
		i       int = 0
		j       int = 0
		k       int = 0
		l       int = 0
		m       int = 0
		found   int = 0
		buf     [64936]byte
		avg_dam float32
		len_    uint64 = 0
		ext     *extra_descr_data
		ext2    *extra_descr_data
	)
	one_argument(argument, &buf[0])
	if buf == nil || buf[0] == 0 || libc.StrCmp(&buf[0], libc.CString(".")) == 0 {
		zrnum = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone
	} else {
		zrnum = real_zone(zone_vnum(libc.Atoi(libc.GoString(&buf[0]))))
	}
	if zrnum == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Check what zone ?\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("Checking zone %d!\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Number)
	}
	send_to_char(ch, libc.CString("Checking Mobs for limits...\r\n"))
	for i = 0; i < int(top_of_mobt); i++ {
		if real_zone_by_thing(room_vnum((*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum)) == zrnum {
			mob = (*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))
			if libc.StrCmp(mob.Name, libc.CString("mob unfinished")) == 0 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Alias hasn't been set.\r\n"))
			}
			if libc.StrCmp(mob.Short_descr, libc.CString("the unfinished mob")) == 0 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Short description hasn't been set.\r\n"))
			}
			if libc.StrNCmp(mob.Long_descr, libc.CString("An unfinished mob stands here."), 30) == 0 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Long description hasn't been set.\r\n"))
			}
			if mob.Description != nil && *mob.Description != 0 {
				if libc.StrNCmp(mob.Description, libc.CString("It looks unfinished."), 20) == 0 && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Description hasn't been set.\r\n"))
				} else if libc.StrNCmp(mob.Description, libc.CString("   "), 3) != 0 && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Description hasn't been formatted. (/fi)\r\n"))
				}
			}
			if GET_LEVEL(mob) > 100 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Is level %d (limit: 1-%d)\r\n", GET_LEVEL(mob), 100))
			}
			if mob.Damage_mod > MAX(GET_LEVEL(mob)/5, 2) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Damage mod of %d is too high (limit: %d)\r\n", mob.Damage_mod, MAX(GET_LEVEL(mob)/5, 2)))
			}
			avg_dam = float32(((float64(mob.Mob_specials.Damsizedice) / 2.0) * float64(mob.Mob_specials.Damnodice)) + float64(mob.Damage_mod))
			if avg_dam > MAX_MOB_DAM_ALLOWED && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- average damage of %4.1f is too high (limit: %d)\r\n", avg_dam, MAX_MOB_DAM_ALLOWED))
			}
			if int(mob.Mob_specials.Damsizedice) == 1 && int(mob.Mob_specials.Damnodice) == 1 && GET_LEVEL(mob) == 0 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Needs to be fixed - %sAutogenerate!%s\r\n", func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KYEL
					}
					return KNUL
				}(), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KNRM
					}
					return KNUL
				}()))
			}
			if MOB_FLAGGED(mob, MOB_AGGRESSIVE) && MOB_FLAGGED(mob, bitvector_t(int32(int(MOB_AGGR_GOOD|MOB_AGGR_EVIL)|MOB_AGGR_NEUTRAL))) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Both aggresive and agressive to align.\r\n"))
			}
			if mob.Gold > GET_LEVEL(mob)*3000 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Set to %d Gold (limit : %d).\r\n", mob.Gold, GET_LEVEL(mob)*3000))
			}
			if mob.Exp > int64(GET_LEVEL(mob)*GET_LEVEL(mob)*120) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Has %lld experience (limit: %d)\r\n", mob.Exp, GET_LEVEL(mob)*GET_LEVEL(mob)*120))
			}
			if AFF_FLAGGED(mob, bitvector_t(int32(int(AFF_GROUP|AFF_CHARM)|AFF_POISON))) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Has illegal affection bits set (%s %s %s)\r\n", func() string {
					if AFF_FLAGGED(mob, AFF_GROUP) {
						return "GROUP"
					}
					return ""
				}(), func() string {
					if AFF_FLAGGED(mob, AFF_CHARM) {
						return "CHARM"
					}
					return ""
				}(), func() string {
					if AFF_FLAGGED(mob, AFF_POISON) {
						return "POISON"
					}
					return ""
				}()))
			}
			if MOB_FLAGGED(mob, MOB_SPEC) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- SPEC flag needs to be removed.\r\n"))
			}
			if found != 0 {
				send_to_char(ch, libc.CString("%s[%5d]%s %-30s: %s\r\n%s"), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KCYN
					}
					return KNUL
				}(), GET_MOB_VNUM(mob), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KYEL
					}
					return KNUL
				}(), GET_NAME(mob), func() string {
					if (func() int {
						if !IS_NPC(ch) {
							if PRF_FLAGGED(ch, PRF_COLOR) {
								return 1
							}
							return 0
						}
						return 0
					}()) >= C_ON {
						return KNRM
					}
					return KNUL
				}(), &buf[0])
			}
			libc.StrCpy(&buf[0], libc.CString(""))
			found = 0
			len_ = 0
		}
	}
	send_to_char(ch, libc.CString("\r\nChecking Objects for limits...\r\n"))
	for i = 0; i < int(top_of_objt); i++ {
		if real_zone_by_thing(room_vnum((*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum)) == zrnum {
			obj = (*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))
			switch obj.Type_flag {
			case ITEM_MONEY:
				if (func() int {
					value = obj.Value[1]
					return value
				}()) > GET_LEVEL(mob)*3000 && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Is worth %d (money limit %d coins).\r\n", value, GET_LEVEL(mob)*3000))
				}
			case ITEM_WEAPON:
				if (obj.Value[3]) >= NUM_ATTACK_TYPES && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has out of range attack type %d.\r\n", obj.Value[3]))
				}
				if ((float64((obj.Value[2])+1)/2.0)*float64(obj.Value[1])) > MAX_DAM_ALLOWED && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Damroll is %2.1f (limit %d)\r\n", (float64((obj.Value[2])+1)/2.0)*float64(obj.Value[1]), MAX_DAM_ALLOWED))
				}
			case ITEM_ARMOR:
				ac = obj.Value[0]
				for j = 0; j < (int(NUM_ITEM_WEARS - 2)); j++ {
					if OBJWEAR_FLAGGED(obj, zarmor[j].Bitvector) && ac > zarmor[j].Ac_allowed && (func() int {
						found = 1
						return found
					}()) != 0 {
						len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Has AC %d (%s limit is %d)\r\n", ac, zarmor[j].Message, zarmor[j].Ac_allowed))
					}
				}
			}
			if !OBJWEAR_FLAGGED(obj, ITEM_WEAR_TAKE) {
				if (obj.Cost != 0 || obj.Weight != 0 && int(obj.Type_flag) != ITEM_FOUNTAIN || obj.Cost_per_day != 0) && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- is NO_TAKE, but has cost (%d) weight (%lld) or rent (%d) set.\r\n", obj.Cost, obj.Weight, obj.Cost_per_day))
				}
			} else {
				if obj.Cost == 0 && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has 0 cost (min. 1).\r\n"))
				}
				if obj.Weight == 0 && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has 0 weight (min. 1).\r\n"))
				}
				if obj.Weight > MAX_OBJ_WEIGHT && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "  Weight is too high: %lld (limit  %d).\r\n", obj.Weight, MAX_OBJ_WEIGHT))
				}
				if obj.Cost > MAX_OBJ_COST && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has %d cost (max %d).\r\n", obj.Cost, MAX_OBJ_COST))
				}
			}
			if obj.Level > int(ADMLVL_IMMORT-1) && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has min level set to %d (max %d).\r\n", obj.Level, int(ADMLVL_IMMORT-1)))
			}
			if obj.Action_description != nil && *obj.Action_description != 0 && int(obj.Type_flag) != ITEM_STAFF && int(obj.Type_flag) != ITEM_WAND && int(obj.Type_flag) != ITEM_SCROLL && int(obj.Type_flag) != ITEM_NOTE && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has action_description set, but is inappropriate type.\r\n"))
			}
			for func() int {
				affs = 0
				return func() int {
					j = 0
					return j
				}()
			}(); j < MAX_OBJ_AFFECT; j++ {
				if obj.Affected[j].Modifier != 0 {
					affs++
				}
			}
			if affs > MAX_AFFECTS_ALLOWED && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has %d affects (limit %d).\r\n", affs, MAX_AFFECTS_ALLOWED))
			}
			for j = 0; j < MAX_OBJ_AFFECT; j++ {
				if zaffs[obj.Affected[j].Location].Max_aff != -99 && (obj.Affected[j].Modifier > zaffs[obj.Affected[j].Location].Max_aff || obj.Affected[j].Modifier < zaffs[obj.Affected[j].Location].Min_aff || zaffs[obj.Affected[j].Location].Min_aff == zaffs[obj.Affected[j].Location].Max_aff) && (func() int {
					found = 1
					return found
				}()) != 0 {
					len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- apply to %s is %d (limit %d - %d).\r\n", zaffs[obj.Affected[j].Location].Message, obj.Affected[j].Modifier, zaffs[obj.Affected[j].Location].Min_aff, zaffs[obj.Affected[j].Location].Max_aff))
				}
			}
			for func() int {
				todam = 0
				tohit = 0
				return func() int {
					j = 0
					return j
				}()
			}(); j < MAX_OBJ_AFFECT; j++ {
				if obj.Affected[j].Location == APPLY_DAMAGE {
					todam += obj.Affected[j].Modifier
				}
			}
			if cmath.Abs(int64(todam)) > MAX_APPLY_DAMAGE_MOD_TOTAL && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- total damage mod %d out of range (limit +/-%d.\r\n", todam, MAX_APPLY_DAMAGE_MOD_TOTAL))
			}
			if cmath.Abs(int64(tohit)) > MAX_APPLY_ACCURCY_MOD_TOTAL && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- total accurcy mod %d out of range (limit +/-%d).\r\n", tohit, MAX_APPLY_ACCURCY_MOD_TOTAL))
			}
			for func() *extra_descr_data {
				ext2 = nil
				return func() *extra_descr_data {
					ext = obj.Ex_description
					return ext
				}()
			}(); ext != nil; ext = ext.Next {
				if libc.StrNCmp(ext.Description, libc.CString("   "), 3) != 0 {
					ext2 = ext
				}
			}
			if ext2 != nil && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has unformatted extra description\r\n"))
			}
			if found != 0 {
				send_to_char(ch, libc.CString("[%5d] %-30s: \r\n%s"), GET_OBJ_VNUM(obj), obj.Short_description, &buf[0])
			}
			libc.StrCpy(&buf[0], libc.CString(""))
			len_ = 0
			found = 0
		}
	}
	send_to_char(ch, libc.CString("\r\nChecking Rooms for limits...\r\n"))
	for i = 0; i < int(top_of_world); i++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Zone == zrnum {
			for j = 0; j < NUM_OF_DIRS; j++ {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j] == nil {
					continue
				}
				exroom = room_vnum((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j].To_room)
				if exroom == room_vnum(-1) {
					continue
				}
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(exroom)))).Zone == zrnum {
					continue
				}
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(exroom)))).Zone == (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Zone {
					continue
				}
				for k = 0; offlimit_zones[k] != -1; k++ {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(exroom)))).Zone == real_zone(zone_vnum(offlimit_zones[k])) && (func() int {
						found = 1
						return found
					}()) != 0 {
						len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Exit %s cannot connect to %d (zone off limits).\r\n", dirs[j], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(exroom)))).Number))
					}
				}
			}
			if ROOM_FLAGGED(room_rnum(i), bitvector_t(int32(int(ROOM_ATRIUM|ROOM_HOUSE)|ROOM_HOUSE_CRASH|ROOM_OLC|ROOM_BFS_MARK))) {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Has illegal affection bits set (%s %s %s %s %s)\r\n", func() string {
					if ROOM_FLAGGED(room_rnum(i), ROOM_ATRIUM) {
						return "ATRIUM"
					}
					return ""
				}(), func() string {
					if ROOM_FLAGGED(room_rnum(i), ROOM_HOUSE) {
						return "HOUSE"
					}
					return ""
				}(), func() string {
					if ROOM_FLAGGED(room_rnum(i), ROOM_HOUSE_CRASH) {
						return "HCRSH"
					}
					return ""
				}(), func() string {
					if ROOM_FLAGGED(room_rnum(i), ROOM_OLC) {
						return "OLC"
					}
					return ""
				}(), func() string {
					if ROOM_FLAGGED(room_rnum(i), ROOM_BFS_MARK) {
						return "*"
					}
					return ""
				}()))
			}
			if MIN_ROOM_DESC_LENGTH != 0 && libc.StrLen((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Description) < MIN_ROOM_DESC_LENGTH && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Room description is too short. (%4.4lld of min. %d characters).\r\n", libc.StrLen((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Description), MIN_ROOM_DESC_LENGTH))
			}
			if libc.StrNCmp((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Description, libc.CString("   "), 3) != 0 && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Room description not formatted with indent (/fi in the editor).\r\n"))
			}
			if libc.StrCSpn((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Description, libc.CString("\r\n")) > MAX_COLOUMN_WIDTH && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- Room description not wrapped at %d chars (/fi in the editor).\r\n", MAX_COLOUMN_WIDTH))
			}
			for func() *extra_descr_data {
				ext2 = nil
				return func() *extra_descr_data {
					ext = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Ex_description
					return ext
				}()
			}(); ext != nil; ext = ext.Next {
				if libc.StrNCmp(ext.Description, libc.CString("   "), 3) != 0 {
					ext2 = ext
				}
			}
			if ext2 != nil && (func() int {
				found = 1
				return found
			}()) != 0 {
				len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "- has unformatted extra description\r\n"))
			}
			if found != 0 {
				send_to_char(ch, libc.CString("[%5d] %-30s: \r\n%s"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number, func() *byte {
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name != nil {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name
					}
					return libc.CString("An unnamed room")
				}(), &buf[0])
				libc.StrCpy(&buf[0], libc.CString(""))
				len_ = 0
				found = 0
			}
		}
	}
	for i = 0; i < int(top_of_world); i++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Zone == zrnum {
			m++
			for func() int {
				j = 0
				return func() int {
					k = 0
					return k
				}()
			}(); j < NUM_OF_DIRS; j++ {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j] == nil {
					k++
				}
			}
			if k == NUM_OF_DIRS {
				l++
			}
		}
	}
	if l*3 > m {
		send_to_char(ch, libc.CString("More than 1/3 of the rooms are not linked.\r\n"))
	}
}
func mob_checkload(ch *char_data, mvnum mob_vnum) {
	var (
		cmd_no int
		count  int = 0
		zone   zone_rnum
		mrnum  mob_rnum = real_mobile(mvnum)
	)
	if mrnum == mob_rnum(-1) {
		send_to_char(ch, libc.CString("That mob does not exist.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Checking load info for the mob [%d] %s...\r\n"), mvnum, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(mrnum)))).Short_descr)
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			if int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'M' {
				continue
			}
			if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(mrnum) {
				send_to_char(ch, libc.CString("  [%5d] %s (%d MAX)\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Name, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2)
				count += 1
			}
		}
	}
	if count > 0 {
		send_to_char(ch, libc.CString("@D[@nTotal counted: %s.@D]@n\r\n"), add_commas(int64(count)))
	}
}
func obj_checkload(ch *char_data, ovnum obj_vnum) {
	var (
		cmd_no     int
		count      int = 0
		zone       zone_rnum
		ornum      obj_rnum  = real_object(ovnum)
		lastroom_v room_vnum = 0
		lastroom_r room_rnum = 0
		lastmob_r  mob_rnum  = 0
	)
	if ornum == obj_rnum(-1) {
		send_to_char(ch, libc.CString("That object does not exist.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Checking load info for the obj [%d] %s...\r\n"), ovnum, (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(ornum)))).Short_description)
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command {
			case 'M':
				lastroom_v = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Number
				lastroom_r = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)
				lastmob_r = mob_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'O':
				lastroom_v = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Number
				lastroom_r = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(ornum) {
					send_to_char(ch, libc.CString("  [%5d] %s (%d Max)\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2)
					count += 1
				}
			case 'P':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(ornum) {
					send_to_char(ch, libc.CString("  [%5d] %s (Put in another object [%d Max])\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2)
					count += 1
				}
			case 'G':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(ornum) {
					send_to_char(ch, libc.CString("  [%5d] %s (Given to %s [%d][%d Max])\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(lastmob_r)))).Short_descr, (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(lastmob_r)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2)
					count += 1
				}
			case 'E':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(ornum) {
					send_to_char(ch, libc.CString("  [%5d] %s (Equipped to %s [%d][%d Max])\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(lastmob_r)))).Short_descr, (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(lastmob_r)))).Vnum, (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2)
					count += 1
				}
			case 'R':
				lastroom_v = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)))).Number
				lastroom_r = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 == vnum(ornum) {
					send_to_char(ch, libc.CString("  [%5d] %s (Removed from room)\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name)
					count += 1
				}
			}
		}
	}
	if count > 0 {
		send_to_char(ch, libc.CString("@D[@nTotal counted: %s.@D]@n\r\n"), add_commas(int64(count)))
	}
}
func trg_checkload(ch *char_data, tvnum trig_vnum) {
	var (
		cmd_no     int
		found      int = 0
		zone       zone_rnum
		trnum      trig_rnum = real_trigger(tvnum)
		lastroom_v room_vnum = 0
		lastroom_r room_rnum = 0
		k          room_rnum
		lastmob_r  mob_rnum = 0
		i          mob_rnum
		lastobj_r  obj_rnum = 0
		j          obj_rnum
		tpl        *trig_proto_list
	)
	if trnum == trig_rnum(-1) {
		send_to_char(ch, libc.CString("That trigger does not exist.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Checking load info for the %s trigger [%d] '%s':\r\n"), func() string {
		if int((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trnum)))).Proto.Attach_type) == MOB_TRIGGER {
			return "mobile"
		}
		if int((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trnum)))).Proto.Attach_type) == OBJ_TRIGGER {
			return "object"
		}
		return "room"
	}(), tvnum, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(trnum)))).Proto.Name)
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command {
			case 'M':
				lastroom_v = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Number
				lastroom_r = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)
				lastmob_r = mob_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'O':
				lastroom_v = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)))).Number
				lastroom_r = room_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3)
				lastobj_r = obj_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'P':
				lastobj_r = obj_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'G':
				lastobj_r = obj_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'E':
				lastobj_r = obj_rnum((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1)
			case 'R':
				lastroom_v = 0
				lastroom_r = 0
				lastobj_r = 0
				lastmob_r = 0
				fallthrough
			case 'T':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 != vnum(trnum) {
					break
				}
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == MOB_TRIGGER {
					send_to_char(ch, libc.CString("mob [%5d] %-60s (zedit room %5d)\r\n"), (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(lastmob_r)))).Vnum, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(lastmob_r)))).Short_descr, lastroom_v)
					found = 1
				} else if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == OBJ_TRIGGER {
					send_to_char(ch, libc.CString("obj [%5d] %-60s  (zedit room %d)\r\n"), (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(lastobj_r)))).Vnum, (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(lastobj_r)))).Short_description, lastroom_v)
					found = 1
				} else if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == WLD_TRIGGER {
					send_to_char(ch, libc.CString("room [%5d] %-60s (zedit)\r\n"), lastroom_v, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(lastroom_r)))).Name)
					found = 1
				}
			}
		}
	}
	for i = 0; i < top_of_mobt; i++ {
		if (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Proto_script == nil {
			continue
		}
		for tpl = (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Proto_script; tpl != nil; tpl = tpl.Next {
			if tpl.Vnum == int(tvnum) {
				send_to_char(ch, libc.CString("mob [%5d] %s\r\n"), (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Short_descr)
				found = 1
			}
		}
	}
	for j = 0; j < top_of_objt; j++ {
		if (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(j)))).Proto_script == nil {
			continue
		}
		for tpl = (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(j)))).Proto_script; tpl != nil; tpl = tpl.Next {
			if tpl.Vnum == int(tvnum) {
				send_to_char(ch, libc.CString("obj [%5d] %s\r\n"), (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(j)))).Vnum, (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(j)))).Short_description)
				found = 1
			}
		}
	}
	for k = 0; k < top_of_world; k++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k)))).Proto_script == nil {
			continue
		}
		for tpl = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k)))).Proto_script; tpl != nil; tpl = tpl.Next {
			if tpl.Vnum == int(tvnum) {
				send_to_char(ch, libc.CString("room[%5d] %s\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k)))).Number, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k)))).Name)
				found = 1
			}
		}
	}
	if found == 0 {
		send_to_char(ch, libc.CString("This trigger is not attached to anything.\r\n"))
	}
}
func do_checkloadstatus(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf1 [2048]byte
		buf2 [2048]byte
	)
	two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 || buf2[0] == 0 || !unicode.IsDigit(rune(buf2[0])) {
		send_to_char(ch, libc.CString("Checkload <M | O | T> <vnum>\r\n"))
		return
	}
	if unicode.ToLower(rune(buf1[0])) == 'm' {
		mob_checkload(ch, mob_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
		return
	}
	if unicode.ToLower(rune(buf1[0])) == 'o' {
		obj_checkload(ch, obj_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
		return
	}
	if unicode.ToLower(rune(buf1[0])) == 't' {
		trg_checkload(ch, trig_vnum(libc.Atoi(libc.GoString(&buf2[0]))))
		return
	}
}
func do_findkey(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		dir int
		key int
		arg [2048]byte
		buf [64936]byte
	)
	any_one_arg(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Format: findkey <dir>\r\n"))
	} else if (func() int {
		dir = search_block(&arg[0], &dirs[0], FALSE)
		return dir
	}()) >= 0 || (func() int {
		dir = search_block(&arg[0], &abbr_dirs[0], FALSE)
		return dir
	}()) >= 0 {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) == nil {
			send_to_char(ch, libc.CString("There's no exit in that direction!\r\n"))
		} else if (func() int {
			key = int(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Key)
			return key
		}()) == int(-1) || key == 0 {
			send_to_char(ch, libc.CString("There's no key for that exit.\r\n"))
		} else {
			stdio.Sprintf(&buf[0], "obj %d", key)
			do_checkloadstatus(ch, &buf[0], 0, 0)
		}
	} else {
		send_to_char(ch, libc.CString("What direction is that?!?\r\n"))
	}
}
func do_spells(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i    int
		qend int
	)
	send_to_char(ch, libc.CString("The following spells are in the game:\r\n"))
	for func() int {
		qend = 0
		return func() int {
			i = 0
			return i
		}()
	}(); i < SPELL_SENSU; i++ {
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
func do_boom(ch *char_data, argument *byte, cmd int, subcmd int) {
	if int(ch.Idnum) != 1 {
		send_to_char(ch, libc.CString("Sorry, only the Founder may use the boom command.\r\n"))
		return
	}
	send_to_outdoor(libc.CString("%s shakes the world with a mighty boom!\r\n"), GET_NAME(ch))
}
