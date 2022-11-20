package main

import (
	"github.com/gotranspile/cxgo/runtime/csys"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unsafe"
)

var bboards *board_info = nil

func init_boards() {
	var (
		i          int
		j          int
		board_vnum int
		xd         xap_dir
		tmp_board  *board_info
		dir_name   [128]byte
	)
	if insure_directory(libc.CString("etc/boards/"), 0) == 0 {
		basic_mud_log(libc.CString("Unable to open/create directory '%s' - Exiting"), "etc/boards/")
		os.Exit(1)
	}
	libc.StrCpy(&dir_name[0], libc.CString("etc/boards"))
	if (func() int {
		i = xdir_scan(&dir_name[0], &xd)
		return i
	}()) <= 0 {
		basic_mud_log(libc.CString("Funny, no board files found.\n"))
		return
	}
	for j = 0; j < i; j++ {
		if libc.StrCmp(libc.CString(".."), xdir_get_name(&xd, j)) != 0 && libc.StrCmp(libc.CString("."), xdir_get_name(&xd, j)) != 0 && libc.StrCmp(libc.CString(".cvsignore"), xdir_get_name(&xd, j)) != 0 {
			stdio.Sscanf(xdir_get_name(&xd, j), "%ld", &board_vnum)
			if (func() *board_info {
				tmp_board = load_board(board_vnum)
				return tmp_board
			}()) != nil {
				tmp_board.Next = bboards
				bboards = tmp_board
			}
		}
	}
	look_at_boards()
}
func create_new_board(board_vnum int) *board_info {
	var (
		buf    [512]byte
		fl     *stdio.File
		temp   *board_info = nil
		backup *board_info
		obj    *obj_data = nil
	)
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&buf[0]), "r")
		return fl
	}()) != nil {
		fl.Close()
		basic_mud_log(libc.CString("Preexisting board file when attempting to create new board [vnum: %d]. Attempting to correct."), board_vnum)
		stdio.Unlink(&buf[0])
		for func() *board_info {
			temp = bboards
			return func() *board_info {
				backup = nil
				return backup
			}()
		}(); temp != nil && backup == nil; temp = temp.Next {
			if temp.Vnum == board_vnum {
				backup = temp
			}
		}
		if backup != nil {
			if backup == bboards {
				bboards = backup.Next
			} else {
				temp = bboards
				for temp != nil && temp.Next != backup {
					temp = temp.Next
				}
				if temp != nil {
					temp.Next = backup.Next
				}
			}
			clear_one_board(backup)
		}
	}
	temp = new(board_info)
	if real_object(board_vnum) == int(-1) {
		basic_mud_log(libc.CString("Creating board [vnum: %d] though no associated object with that vnum can be found. Using defaults."), board_vnum)
		temp.Read_lvl = config_info.Play.Level_cap
		temp.Write_lvl = config_info.Play.Level_cap
		temp.Remove_lvl = config_info.Play.Level_cap
	} else {
		obj = &(obj_proto[real_object(board_vnum)])
		temp.Read_lvl = obj.Value[VAL_BOARD_READ]
		temp.Write_lvl = obj.Value[VAL_BOARD_WRITE]
		temp.Remove_lvl = obj.Value[VAL_BOARD_ERASE]
	}
	temp.Vnum = board_vnum
	temp.Num_messages = 0
	temp.Version = CURRENT_BOARD_VER
	temp.Next = nil
	temp.Messages = nil
	if !save_board(temp) {
		basic_mud_log(libc.CString("Hm. Error while creating new board file [vnum: %d]. Unable to create new file."), board_vnum)
		libc.Free(unsafe.Pointer(temp))
		return nil
	}
	return temp
}
func save_board(ts *board_info) bool {
	var (
		message  *board_msg
		memboard *board_memory
		fl       *stdio.File
		buf      [512]byte
		i        int = 1
	)
	stdio.Sprintf(&buf[0], "%s%d", "etc/boards/", ts.Vnum)
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&buf[0]), "wb")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("Hm. Error while creating attempting to save board [vnum: %d].  Unable to create file '%s'"), ts.Vnum, &buf[0])
		return false
	}
	stdio.Fprintf(fl, "Board File\n%d %d %d %d %d\n", ts.Read_lvl, ts.Write_lvl, ts.Remove_lvl, ts.Num_messages, CURRENT_BOARD_VER)
	for message = ts.Messages; message != nil; message = message.Next {
		if ts.Version != CURRENT_BOARD_VER {
			message.Name = get_name_by_id(message.Poster)
		}
		if message != nil {
			stdio.Fprintf(fl, "#%d\n%s\n%ld\n%s\n%s~\n", func() int {
				p := &i
				x := *p
				*p++
				return x
			}(), message.Name, message.Timestamp, message.Subject, message.Data)
		}
	}
	for i = 0; i != 301; i++ {
		memboard = ts.Memory[i]
		for memboard != nil {
			stdio.Fprintf(fl, "S%d %s %d\n", i, memboard.Name, +memboard.Timestamp)
			memboard = memboard.Next
		}
	}
	fl.Close()
	return true
}
func load_board(board_vnum int) *board_info {
	var (
		temp_board  *board_info
		bmsg        *board_msg
		obj         *obj_data = nil
		st          csys.StatRes
		memboard    *board_memory
		list        *board_memory
		t           [10]int
		mnum        int
		poster      int
		timestamp   int
		msg_num     int
		retval      int = 0
		filebuf     [512]byte
		buf         [512]byte
		poster_name [128]byte
		fl          *stdio.File
		sflag       int
	)
	stdio.Sprintf(&filebuf[0], "%s%d", "etc/boards/", board_vnum)
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(&filebuf[0]), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("Request to open board [vnum %d] failed - unable to open file '%s'."), board_vnum, &filebuf[0])
		return nil
	}
	get_line(fl, &buf[0])
	if libc.StrCmp(libc.CString("Board File"), &buf[0]) != 0 {
		basic_mud_log(libc.CString("Invalid board file '%s' [vnum: %d] - failed to load."), &filebuf[0], board_vnum)
		return nil
	}
	temp_board = new(board_info)
	temp_board.Vnum = board_vnum
	get_line(fl, &buf[0])
	if (func() int {
		retval = stdio.Sscanf(&buf[0], "%d %d %d %d %d", &t[0], &t[1], &t[2], &t[3], &t[4])
		return retval
	}()) != 5 {
		if retval == 4 {
			basic_mud_log(libc.CString("Parse error on board [vnum: %d], file '%s' - attempting to correct [4] args expecting 5."), board_vnum, &filebuf[0])
			t[4] = 1
		} else if retval != 4 {
			basic_mud_log(libc.CString("Parse error on board [vnum: %d], file '%s' - attempting to correct [< 4] args expecting 5."), board_vnum, &filebuf[0])
			t[0] = func() int {
				p := &t[1]
				t[1] = func() int {
					p := &t[2]
					t[2] = config_info.Play.Level_cap
					return *p
				}()
				return *p
			}()
			t[3] = -1
			t[4] = 1
		}
	}
	if real_object(board_vnum) == int(-1) {
		basic_mud_log(libc.CString("No associated object exists when attempting to create a board [vnum %d]."), board_vnum)
		csys.Stat(&filebuf[0], &st)
		// todo : fix this
		//if libc.TimeVal(libc.GetTime(nil))-st.MTime > libc.TimeVal(60*60*24*7) {
		if false {
			basic_mud_log(libc.CString("Deleting old board file '%s' [vnum %d].  7 days without modification & no associated object."), &filebuf[0], board_vnum)
			stdio.Unlink(&filebuf[0])
			libc.Free(unsafe.Pointer(temp_board))
			return nil
		}
		temp_board.Read_lvl = t[0]
		temp_board.Write_lvl = t[1]
		temp_board.Remove_lvl = t[2]
		temp_board.Num_messages = t[3]
		temp_board.Version = t[4]
		basic_mud_log(libc.CString("Board vnum %d, Version %d"), temp_board.Vnum, temp_board.Version)
	} else {
		obj = &(obj_proto[real_object(board_vnum)])
		if t[0] != (obj.Value[VAL_BOARD_READ]) || t[1] != (obj.Value[VAL_BOARD_WRITE]) || t[2] != (obj.Value[VAL_BOARD_ERASE]) {
			basic_mud_log(libc.CString("Mismatch in board <-> object read/write/remove settings for board [vnum: %d]. Correcting."), board_vnum)
		}
		temp_board.Read_lvl = obj.Value[VAL_BOARD_READ]
		temp_board.Write_lvl = obj.Value[VAL_BOARD_WRITE]
		temp_board.Remove_lvl = obj.Value[VAL_BOARD_ERASE]
		temp_board.Num_messages = t[3]
		temp_board.Version = t[4]
	}
	temp_board.Next = nil
	temp_board.Messages = nil
	msg_num = 0
	for get_line(fl, &buf[0]) != 0 {
		if buf[0] == 'S' && temp_board.Version != CURRENT_BOARD_VER {
			if stdio.Sscanf(&buf[0], "S %d %d %d ", &mnum, &poster, &timestamp) == 3 {
				memboard = new(board_memory)
				memboard.Reader = poster
				memboard.Timestamp = timestamp
			}
		} else if buf[0] == 'S' && temp_board.Version == CURRENT_BOARD_VER {
			if stdio.Sscanf(&buf[0], "S %d %s %d ", &mnum, &poster_name[0], &timestamp) == 3 {
				memboard = new(board_memory)
				memboard.Name = libc.StrDup(&poster_name[0])
				memboard.Timestamp = timestamp
				if get_name_by_id(poster) == nil && temp_board.Version != CURRENT_BOARD_VER {
					libc.Free(unsafe.Pointer(memboard))
				} else if temp_board.Version == CURRENT_BOARD_VER {
					libc.Free(unsafe.Pointer(memboard))
				} else {
					if temp_board.Version == CURRENT_BOARD_VER {
						for func() int {
							bmsg = temp_board.Messages
							return func() int {
								sflag = 0
								return sflag
							}()
						}(); bmsg != nil && sflag == 0; bmsg = bmsg.Next {
							if int(bmsg.Timestamp) == memboard.Timestamp && mnum == ((int(bmsg.Timestamp%301)+get_id_by_name(bmsg.Name)%301)%301) {
								sflag = 1
							}
						}
					} else {
						for func() int {
							bmsg = temp_board.Messages
							return func() int {
								sflag = 0
								return sflag
							}()
						}(); bmsg != nil && sflag == 0; bmsg = bmsg.Next {
							if int(bmsg.Timestamp) == memboard.Timestamp && mnum == ((int(bmsg.Timestamp%301)+bmsg.Poster%301)%301) {
								sflag = 1
							}
						}
					}
					if sflag != 0 {
						if (temp_board.Memory[mnum]) != nil {
							list = temp_board.Memory[mnum]
							temp_board.Memory[mnum] = memboard
							memboard.Next = list
						} else {
							temp_board.Memory[mnum] = memboard
							memboard.Next = nil
						}
					} else {
						libc.Free(unsafe.Pointer(memboard))
					}
				}
			}
		} else if buf[0] == '#' {
			if parse_message(fl, temp_board) != 0 {
				msg_num++
			}
		}
	}
	fl.Close()
	if msg_num != temp_board.Num_messages {
		basic_mud_log(libc.CString("Board [vnum: %d] message count (%d) not equal to actual message count (%d). Correcting."), temp_board.Vnum, temp_board.Num_messages, msg_num)
		temp_board.Num_messages = msg_num
	}
	save_board(temp_board)
	return temp_board
}
func parse_message(fl *stdio.File, temp_board *board_info) int {
	var (
		tmsg    *board_msg
		t2msg   *board_msg
		subject [81]byte
		buf     [4097]byte
		poster  [128]byte
	)
	tmsg = new(board_msg)
	if temp_board.Version != CURRENT_BOARD_VER {
		if stdio.Fscanf(fl, "%ld\n", &tmsg.Poster) != 1 || stdio.Fscanf(fl, "%ld\n", &tmsg.Timestamp) != 1 {
			basic_mud_log(libc.CString("Parse error in message for board [vnum: %d].  Skipping."), temp_board.Vnum)
			libc.Free(unsafe.Pointer(tmsg))
			return 0
		}
	} else {
		if stdio.Fscanf(fl, "%s\n", &poster[0]) != 1 || stdio.Fscanf(fl, "%ld\n", &tmsg.Timestamp) != 1 {
			basic_mud_log(libc.CString("Parse error in message for board [vnum: %d].  Skipping."), temp_board.Vnum)
			libc.Free(unsafe.Pointer(tmsg))
			return 0
		}
		tmsg.Name = libc.StrDup(&poster[0])
	}
	get_line(fl, &subject[0])
	tmsg.Subject = libc.StrDup(&subject[0])
	tmsg.Data = fread_string(fl, &buf[0])
	tmsg.Next = nil
	tmsg.Next = func() *board_msg {
		p := &tmsg.Prev
		tmsg.Prev = nil
		return *p
	}()
	if temp_board.Messages != nil {
		t2msg = temp_board.Messages
		for t2msg.Next != nil {
			t2msg = t2msg.Next
		}
		t2msg.Next = tmsg
		tmsg.Prev = t2msg
	} else {
		tmsg.Prev = nil
		temp_board.Messages = tmsg
	}
	return 1
}
func look_at_boards() {
	var (
		counter  int
		messages int         = 0
		tboard   *board_info = bboards
		msg      *board_msg
	)
	for counter = 0; tboard != nil; counter++ {
		msg = tboard.Messages
		for msg != nil {
			messages++
			msg = msg.Next
		}
		tboard = tboard.Next
	}
	basic_mud_log(libc.CString("There are %d boards located; %d messages"), counter, messages)
}
func clear_boards() {
	var (
		tmp  *board_info
		tmp2 *board_info
	)
	for tmp = bboards; tmp != nil; tmp = tmp2 {
		tmp2 = tmp.Next
		clear_one_board(tmp)
	}
}
func clear_one_board(tmp *board_info) {
	var (
		m1   *board_msg
		m2   *board_msg
		mem1 *board_memory
		mem2 *board_memory
		i    int
	)
	for m1 = tmp.Messages; m1 != nil; m1 = m2 {
		m2 = m1.Next
		libc.Free(unsafe.Pointer(m1.Subject))
		libc.Free(unsafe.Pointer(m1.Data))
		libc.Free(unsafe.Pointer(m1))
	}
	for i = 0; i < 301; i++ {
		for mem1 = tmp.Memory[i]; mem1 != nil; mem1 = mem2 {
			mem2 = mem1.Next
			libc.Free(unsafe.Pointer(mem1))
		}
	}
	libc.Free(unsafe.Pointer(tmp))
	tmp = nil
}
func show_board(board_vnum int, ch *char_data) {
	var (
		thisboard *board_info
		message   *board_msg
		tmstr     *byte
		msgcount  int = 0
		yesno     int = 0
		bnum      int = 0
		buf       [64936]byte
		buf2      [64936]byte
		name      [127]byte
	)
	buf[0] = '\x00'
	buf2[0] = '\x00'
	name[0] = '\x00'
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Gosh.. now .. if only mobs could read.. you'd be doing good.\r\n"))
		return
	}
	thisboard = locate_board(board_vnum)
	if thisboard == nil {
		basic_mud_log(libc.CString("Creating new board - board #%d"), board_vnum)
		thisboard = create_new_board(board_vnum)
		thisboard.Next = bboards
		bboards = thisboard
	}
	if ch.Admlevel < thisboard.Read_lvl {
		send_to_char(ch, libc.CString("You try but fail to understand the holy words.\r\n"))
		return
	}
	var obj *obj_data
	var num int = board_vnum
	if (func() int {
		board_vnum = real_object(num)
		return board_vnum
	}()) == int(-1) {
		basic_mud_log(libc.CString("SYSERR: DEFUNCT BOARD VNUM.\r\n"))
		send_to_char(ch, libc.CString("@W                  This is a bulletin board.\r\n"))
		send_to_char(ch, libc.CString("@rO@b============================================================================@rO@n\n"))
		send_to_char(ch, libc.CString("     @D[@GX@D] means you have read the message, @D[@RX@D] means you have not.\r\n     @WUsage@D:@CREAD@D/@cREMOVE @D<@Wmessg #@D>@W, @CRESPOND @D<@Wmessg #@D>@W, @CWRITE @D<@Wheader@D>@W.@n\r\n     @CVieworder@W, this changes the order in which posts are listed to you.@n\r\n"))
	} else {
		obj = read_object(board_vnum, REAL)
		bnum = GET_OBJ_VNUM(obj)
		var clan [120]byte
		if OBJ_FLAGGED(obj, ITEM_CBOARD) {
			if ch.Clan != nil {
				stdio.Sprintf(&clan[0], "%s", ch.Clan)
			}
			if libc.StrStr(obj.Action_description, &clan[0]) == nil {
				send_to_char(ch, libc.CString("You are incapable of reading this board!\r\n"))
				return
			}
		}
		send_to_char(ch, libc.CString("@W                  This is the %20s\r\n"), obj.Short_description)
		send_to_char(ch, libc.CString("@rO@b============================================================================@rO@n\n     @D[@GX@D] means you have read the message, @D[@RX@D] means you have not.\r\n     @WUsage@D:@CREAD@D/@cREMOVE @D<@Wmessg #@D>@W, @CRESPOND @D<@Wmessg #@D>@W, @CWRITE @D<@Wheader@D>@W.@n\r\n     @CVieworder@W, this changes the order in which posts are listed to you.@n\r\n     @D----------------------------------------------------------------\n"))
		extract_obj(obj)
	}
	if thisboard.Num_messages == 0 || thisboard.Messages == nil {
		stdio.Sprintf(&buf[0], "                  @WThe board is empty.@n\r\n")
		send_to_char(ch, &buf[0])
		return
	} else {
		send_to_char(ch, libc.CString("                  @WThere %s %d %s on the board.@n\r\n"), func() string {
			if thisboard.Num_messages == 1 {
				return "is"
			}
			return "are"
		}(), thisboard.Num_messages, func() string {
			if thisboard.Num_messages == 1 {
				return "message"
			}
			return "messages"
		}())
	}
	message = thisboard.Messages
	if PRF_FLAGGED(ch, PRF_VIEWORDER) {
		for message.Next != nil {
			message = message.Next
		}
	}
	for message != nil {
		tmstr = libc.AscTime(libc.LocalTime(&message.Timestamp))
		*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
		yesno = int(libc.BoolToInt(mesglookup(message, ch, thisboard)))
		if thisboard.Version != CURRENT_BOARD_VER {
			stdio.Snprintf(&name[0], int(127), "%s", get_name_by_id(message.Poster))
		} else {
			stdio.Snprintf(&name[0], int(127), "%s", message.Name)
		}
		if msgcount < 1 {
			stdio.Sprintf(&buf[0], "@D[%s] (@C%2d@D) : @W%6.10s @D(@G%-10s@D) ::@w %-45s\r\n", func() string {
				if yesno != 0 {
					return "@GX@D"
				}
				return "@RX@D"
			}(), func() int {
				p := &msgcount
				*p++
				return *p
			}(), tmstr, CAP(&name[0]), func() *byte {
				if message.Subject != nil {
					return message.Subject
				}
				return libc.CString("No Subject")
			}())
		} else {
			stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@D[%s] (@C%2d@D) : @W%6.10s @D(@G%-10s@D) ::@w %-45s\r\n", func() string {
				if yesno != 0 {
					return "@GX@D"
				}
				return "@RX@D"
			}(), func() int {
				p := &msgcount
				*p++
				return *p
			}(), tmstr, CAP(&name[0]), func() *byte {
				if message.Subject != nil {
					return message.Subject
				}
				return libc.CString("No Subject")
			}())
		}
		if PRF_FLAGGED(ch, PRF_VIEWORDER) {
			message = message.Prev
		} else {
			message = message.Next
		}
	}
	stdio.Sprintf(&buf[libc.StrLen(&buf[0])], "@rO@b============================================================================@rO@n\n")
	libc.StrCpy(&buf2[0], &buf[0])
	page_string(ch.Desc, &buf2[0], 1)
	if bnum == 3092 {
		ch.Lboard[0] = libc.GetTime(nil)
	}
	if bnum == 3099 {
		ch.Lboard[3] = libc.GetTime(nil)
	}
	if bnum == 3098 {
		ch.Lboard[1] = libc.GetTime(nil)
	}
	if bnum == 3090 {
		ch.Lboard[4] = libc.GetTime(nil)
	}
	save_char(ch)
}
func board_display_msg(board_vnum int, ch *char_data, arg int) {
	var (
		thisboard   *board_info = bboards
		message     *board_msg
		tmstr       *byte
		msgcount    int
		mem         int
		sflag       int
		name        [127]byte
		mboard_type *board_memory
		list        *board_memory
		buf         [64937]byte
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Silly mob - reading is for pcs!\r\n"))
		return
	}
	thisboard = locate_board(board_vnum)
	if thisboard == nil {
		basic_mud_log(libc.CString("Creating new board - board #%d"), board_vnum)
		thisboard = create_new_board(board_vnum)
	}
	if ch.Admlevel < thisboard.Read_lvl {
		send_to_char(ch, libc.CString("You try but fail to understand the holy words.\r\n"))
		return
	}
	if thisboard.Messages == nil {
		send_to_char(ch, libc.CString("The board is empty!\r\n"))
		return
	}
	var obj *obj_data
	var num int = board_vnum
	var bnum int = 0
	if (func() int {
		board_vnum = real_object(num)
		return board_vnum
	}()) == int(-1) {
		send_to_imm(libc.CString("Error with %d board, object doesn't exist."), board_vnum)
	} else {
		obj = read_object(board_vnum, REAL)
		bnum = GET_OBJ_VNUM(obj)
		var clan [200]byte
		if OBJ_FLAGGED(obj, ITEM_CBOARD) {
			if ch.Clan != nil {
				stdio.Sprintf(&clan[0], "%s", ch.Clan)
			}
			if libc.StrStr(obj.Action_description, &clan[0]) == nil {
				send_to_char(ch, libc.CString("You are incapable of reading this board!\r\n"))
				return
			}
		}
		extract_obj(obj)
	}
	message = thisboard.Messages
	if arg < 1 {
		send_to_char(ch, libc.CString("You must specify the (positive) number of the message to be read!\r\n"))
		return
	}
	if PRF_FLAGGED(ch, PRF_VIEWORDER) {
		for message.Next != nil {
			message = message.Next
		}
	}
	for msgcount = arg; message != nil && msgcount != 1; msgcount-- {
		if PRF_FLAGGED(ch, PRF_VIEWORDER) {
			message = message.Prev
		} else {
			message = message.Next
		}
	}
	if message == nil {
		send_to_char(ch, libc.CString("That message exists only in your imagination.\r\n"))
		return
	}
	if thisboard.Version != CURRENT_BOARD_VER {
		mem = (int(message.Timestamp%301) + message.Poster%301) % 301
	} else {
		mem = (int(message.Timestamp%301) + get_id_by_name(message.Name)%301) % 301
	}
	mboard_type = new(board_memory)
	if thisboard.Version != CURRENT_BOARD_VER {
		mboard_type.Reader = int(ch.Idnum)
	} else {
		mboard_type.Name = libc.StrDup(GET_NAME(ch))
	}
	mboard_type.Timestamp = int(message.Timestamp)
	mboard_type.Next = nil
	list = thisboard.Memory[mem]
	sflag = 1
	for list != nil && sflag != 0 {
		if thisboard.Version != CURRENT_BOARD_VER {
			if list.Reader == mboard_type.Reader && list.Timestamp == mboard_type.Timestamp {
				sflag = 0
			}
		} else {
			if libc.StrCmp(list.Name, mboard_type.Name) == 0 && list.Timestamp == mboard_type.Timestamp {
				sflag = 0
			}
		}
		list = list.Next
	}
	if sflag != 0 {
		list = thisboard.Memory[mem]
		thisboard.Memory[mem] = mboard_type
		mboard_type.Next = list
	} else {
		if mboard_type != nil {
		}
	}
	tmstr = libc.AscTime(libc.LocalTime(&message.Timestamp))
	*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
	if thisboard.Version != CURRENT_BOARD_VER {
		stdio.Snprintf(&name[0], int(127), "%s", get_name_by_id(message.Poster))
	} else {
		stdio.Snprintf(&name[0], int(127), "%s", message.Name)
	}
	stdio.Sprintf(&buf[0], "@r_____________________________________________________________________________@n\r\n\r\n@gMessage @W[@Y%2d@W] @D: @c%6.10s @D(@C%s@D)\r\n@GTopic        @D: @w%s\r\n@r-------------------------------------------------------------------------@n\r\n\r\n%s\n@r_____________________________________________________________________________@n\r\n", arg, tmstr, CAP(&name[0]), func() *byte {
		if message.Subject != nil {
			return message.Subject
		}
		return libc.CString("No Subject")
	}(), func() *byte {
		if message.Data != nil {
			return message.Data
		}
		return libc.CString("Looks like this message is empty.")
	}())
	page_string(ch.Desc, &buf[0], 1)
	if bnum == 3092 {
		ch.Lboard[0] = libc.GetTime(nil)
	}
	if bnum == 3099 {
		ch.Lboard[3] = libc.GetTime(nil)
	}
	if bnum == 3098 {
		ch.Lboard[1] = libc.GetTime(nil)
	}
	if bnum == 3090 {
		ch.Lboard[4] = libc.GetTime(nil)
	}
	if sflag != 0 {
		save_board(thisboard)
	}
}
func mesglookup(message *board_msg, ch *char_data, board *board_info) bool {
	var (
		mem         int = 0
		mboard_type *board_memory
		tempname    *byte = nil
	)
	if board.Version != CURRENT_BOARD_VER {
		mem = (int(message.Timestamp%301) + message.Poster%301) % 301
	} else {
		mem = (int(message.Timestamp%301) + get_id_by_name(message.Name)%301) % 301
	}
	mboard_type = board.Memory[mem]
	for mboard_type != nil && board.Version != CURRENT_BOARD_VER {
		if mboard_type.Reader == int(ch.Idnum) && mboard_type.Timestamp == int(message.Timestamp) {
			return true
		} else {
			mboard_type = mboard_type.Next
		}
	}
	tempname = libc.StrDup(GET_NAME(ch))
	for mboard_type != nil && board.Version == CURRENT_BOARD_VER {
		if libc.StrCmp(mboard_type.Name, tempname) == 0 && mboard_type.Timestamp == int(message.Timestamp) {
			return true
		} else {
			mboard_type = mboard_type.Next
		}
	}
	libc.Free(unsafe.Pointer(tempname))
	return false
}
func write_board_message(board_vnum int, ch *char_data, arg *byte) {
	var (
		thisboard *board_info = bboards
		message   *board_msg
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Orwellian police thwart your attempt at free speech.\r\n"))
		return
	}
	thisboard = locate_board(board_vnum)
	if thisboard == nil {
		send_to_char(ch, libc.CString("Error: Your board could not be found. Please report.\n"))
		basic_mud_log(libc.CString("Error in write_board_msg - board #%d"), board_vnum)
		return
	}
	if ch.Admlevel < thisboard.Write_lvl {
		send_to_char(ch, libc.CString("You are not holy enough to write on this board.\r\n"))
		return
	}
	if *arg == 0 || arg == nil {
		stdio.Sprintf(arg, "No Subject")
	}
	if libc.StrLen(arg) > 46 {
		send_to_char(ch, libc.CString("Your subject can only be 45 characters long(including colorcode).\r\n"))
		return
	}
	act(libc.CString("@C$n@w starts writing on the board.@n"), 1, ch, nil, nil, TO_ROOM)
	skip_spaces(&arg)
	delete_doubledollar(arg)
	*(*byte)(unsafe.Add(unsafe.Pointer(arg), 81)) = '\x00'
	message = new(board_msg)
	message.Name = libc.StrDup(GET_NAME(ch))
	message.Timestamp = libc.GetTime(nil)
	message.Subject = libc.StrDup(arg)
	message.Next = nil
	message.Prev = nil
	message.Data = nil
	thisboard.Num_messages = int(MAX(int64(thisboard.Num_messages+1), 1))
	message.Next = thisboard.Messages
	if thisboard.Messages != nil {
		thisboard.Messages.Prev = message
	}
	thisboard.Messages = message
	send_to_char(ch, libc.CString("Write your message.  (/s saves /h for help)\r\n"))
	SET_BIT_AR(ch.Act[:], PLR_WRITING)
	string_write(ch.Desc, &message.Data, MAX_MESSAGE_LENGTH, board_vnum+BOARD_MAGIC, nil)
	if board_vnum == 3092 {
		BOARDNEWMORT = libc.GetTime(nil)
	}
	if board_vnum == 3098 {
		BOARDNEWIMM = libc.GetTime(nil)
	}
	if board_vnum == 3099 {
		BOARDNEWDUO = libc.GetTime(nil)
	}
	if board_vnum == 3090 {
		BOARDNEWBUI = libc.GetTime(nil)
	}
	save_mud_time(&time_info)
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if !IS_PLAYING(d) {
			continue
		}
		if PLR_FLAGGED(d.Character, PLR_WRITING) {
			continue
		}
		if d.Character.Admlevel >= 1 && BOARDNEWIMM > (d.Character.Lboard[1]) && board_vnum == 3098 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Immortal Board Post.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWBUI > (d.Character.Lboard[4]) && board_vnum == 3090 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Builder Board Post.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWDUO > (d.Character.Lboard[3]) && board_vnum == 3099 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Punishment Board Post.@n\r\n"))
		}
		if BOARDNEWMORT > (d.Character.Lboard[0]) && board_vnum == 3092 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Mortal Board Post.@n\r\n"))
		}
	}
	return
}
func board_respond(board_vnum int, ch *char_data, mnum int) {
	var (
		thisboard *board_info = bboards
		message   *board_msg
		other     *board_msg
		number    [64936]byte
		buf       [64936]byte
		gcount    int = 0
	)
	thisboard = locate_board(board_vnum)
	if thisboard == nil {
		send_to_char(ch, libc.CString("Error: Your board could not be found. Please report.\n"))
		basic_mud_log(libc.CString("Error in board_respond - board #%ld"), board_vnum)
		return
	}
	if ch.Admlevel < thisboard.Write_lvl {
		send_to_char(ch, libc.CString("You are not holy enough to write on this board.\r\n"))
		return
	}
	if ch.Admlevel < thisboard.Read_lvl {
		send_to_char(ch, libc.CString("You are not holy enough to respond to this board.\r\n"))
		return
	}
	if PRF_FLAGGED(ch, PRF_VIEWORDER) {
		mnum = (thisboard.Num_messages - mnum) + 1
	}
	if mnum < 0 || mnum > thisboard.Num_messages {
		send_to_char(ch, libc.CString("You can only respond to an actual message.\r\n"))
		return
	}
	other = thisboard.Messages
	for gcount = 0; other != nil && gcount != (mnum-1); gcount++ {
		other = other.Next
	}
	message = new(board_msg)
	message.Name = libc.StrDup(GET_NAME(ch))
	message.Timestamp = libc.GetTime(nil)
	stdio.Sprintf(&buf[0], "Re: %s", other.Subject)
	message.Subject = libc.StrDup(&buf[0])
	message.Next = func() *board_msg {
		p := &message.Prev
		message.Prev = nil
		return *p
	}()
	message.Data = nil
	thisboard.Num_messages = thisboard.Num_messages + 1
	message.Next = thisboard.Messages
	if thisboard.Messages != nil {
		thisboard.Messages.Prev = message
	}
	thisboard.Messages = message
	send_to_char(ch, libc.CString("Write your message.  (/s saves /h for help)\r\n\r\n"))
	act(libc.CString("@C$n@w starts writing on the board.@n"), 1, ch, nil, nil, TO_ROOM)
	if !IS_NPC(ch) {
		SET_BIT_AR(ch.Act[:], PLR_WRITING)
	}
	stdio.Sprintf(&number[0], "\t@D------- @cQuoted message @D-------@w\r\n%s\t@D   ------- @cEnd Quote @D-------@w\r\n", other.Data)
	message.Data = libc.StrDup(&number[0])
	ch.Desc.Backstr = libc.StrDup(&number[0])
	write_to_output(ch.Desc, &number[0])
	string_write(ch.Desc, &message.Data, MAX_MESSAGE_LENGTH, board_vnum+BOARD_MAGIC, nil)
	if board_vnum == 3092 {
		BOARDNEWMORT = libc.GetTime(nil)
	}
	if board_vnum == 3098 {
		BOARDNEWIMM = libc.GetTime(nil)
	}
	if board_vnum == 3099 {
		BOARDNEWDUO = libc.GetTime(nil)
	}
	if board_vnum == 3090 {
		BOARDNEWBUI = libc.GetTime(nil)
	}
	save_mud_time(&time_info)
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if !IS_PLAYING(d) {
			continue
		}
		if PLR_FLAGGED(d.Character, PLR_WRITING) {
			continue
		}
		if d.Character.Admlevel >= 1 && BOARDNEWIMM > (d.Character.Lboard[1]) && board_vnum == 3098 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Immortal Board Post.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWBUI > (d.Character.Lboard[4]) && board_vnum == 3090 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Builder Board Post.@n\r\n"))
		}
		if d.Character.Admlevel >= 1 && BOARDNEWDUO > (d.Character.Lboard[3]) && board_vnum == 3099 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Punishment Board Post.@n\r\n"))
		}
		if BOARDNEWMORT > (d.Character.Lboard[0]) && board_vnum == 3092 {
			send_to_char(d.Character, libc.CString("\r\n@GThere is a new Mortal Board Post.@n\r\n"))
		}
	}
	return
}
func locate_board(board_vnum int) *board_info {
	var thisboard *board_info = bboards
	for thisboard != nil {
		if thisboard.Vnum == board_vnum {
			return thisboard
		}
		thisboard = thisboard.Next
	}
	return nil
}
func remove_board_msg(board_vnum int, ch *char_data, arg int) {
	var (
		thisboard *board_info
		cur       *board_msg
		temp      *board_msg
		d         *descriptor_data
		obj       *obj_data
		msgcount  int
		buf       [64937]byte
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Nuts.. looks like you forgot your eraser back in mobland...\r\n"))
		return
	}
	thisboard = locate_board(board_vnum)
	if thisboard == nil {
		send_to_char(ch, libc.CString("Error: Your board could not be found. Please report.\n"))
		basic_mud_log(libc.CString("Error in Board_remove_msg - board #%d"), board_vnum)
		return
	}
	cur = thisboard.Messages
	if arg < 1 {
		send_to_char(ch, libc.CString("You must specify the (positive) number of the message to be read!\r\n"))
		return
	}
	if PRF_FLAGGED(ch, PRF_VIEWORDER) {
		arg = thisboard.Num_messages - arg + 1
	}
	for msgcount = arg; cur != nil && msgcount != 1; msgcount-- {
		cur = cur.Next
	}
	if cur == nil {
		send_to_char(ch, libc.CString("That message exists only in your imagination.\r\n"))
		return
	}
	var num int = board_vnum
	if (func() int {
		board_vnum = real_object(num)
		return board_vnum
	}()) == int(-1) {
		basic_mud_log(libc.CString("Board doesn't exists! Weird."))
		return
	} else {
		var clan [120]byte
		obj = read_object(board_vnum, REAL)
		if OBJ_FLAGGED(obj, ITEM_CBOARD) {
			if ch.Clan != nil {
				stdio.Sprintf(&clan[0], "%s", ch.Clan)
			}
			if clanIsModerator(&clan[0], ch) && libc.StrStr(obj.Action_description, &clan[0]) != nil {
				send_to_char(ch, libc.CString("Exercising your clan leader powers....\r\n"))
			} else if ch.Admlevel < thisboard.Remove_lvl && libc.StrCmp(GET_NAME(ch), cur.Name) != 0 {
				send_to_char(ch, libc.CString("You can't remove other people's messages.\r\n"))
				extract_obj(obj)
				return
			}
		} else if !OBJ_FLAGGED(obj, ITEM_CBOARD) {
			if ch.Admlevel < thisboard.Remove_lvl && libc.StrCmp(GET_NAME(ch), cur.Name) != 0 {
				send_to_char(ch, libc.CString("You can't remove other people's messages.\r\n"))
				extract_obj(obj)
				return
			}
		}
		extract_obj(obj)
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == 0 && d.Str == &cur.Data {
			send_to_char(ch, libc.CString("At least wait until the author is finished before removing it!\r\n"))
			return
		}
	}
	if cur == thisboard.Messages {
		thisboard.Messages = cur.Next
		if thisboard.Messages != nil {
			thisboard.Messages.Prev = nil
		}
	} else {
		temp = thisboard.Messages
		for temp != nil && temp.Next != cur {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = cur.Next
			if cur.Next != nil {
				cur.Next.Prev = temp
			}
		}
	}
	libc.Free(unsafe.Pointer(cur))
	cur = nil
	thisboard.Num_messages = thisboard.Num_messages - 1
	send_to_char(ch, libc.CString("Message removed.\r\n"))
	stdio.Sprintf(&buf[0], "$n just removed message %d.", arg)
	act(&buf[0], 0, ch, nil, nil, TO_ROOM)
	save_board(thisboard)
	return
}
