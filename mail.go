package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const MIN_MAIL_LEVEL = 3
const STAMP_PRICE = 10
const MAX_MAIL_SIZE = 4096
const BLOCK_SIZE = 256
const HEADER_BLOCK = -1
const LAST_BLOCK = -2
const DELETED_BLOCK = -3

type header_data_type struct {
	Next_block int
	From       int
	To         int
	Mail_time  libc.Time
}
type header_block_type_d struct {
	Block_type  int
	Header_data header_data_type
	Txt         [216]byte
}
type data_block_type_d struct {
	Block_type int
	Txt        [248]byte
}
type header_block_type header_block_type_d
type data_block_type data_block_type_d
type position_list_type_d struct {
	Position int
	Next     *position_list_type_d
}
type position_list_type position_list_type_d
type mail_index_type_d struct {
	Recipient  int
	List_start *position_list_type
	Next       *mail_index_type_d
}
type mail_index_type mail_index_type_d

var mail_index *mail_index_type = nil
var free_list *position_list_type = nil
var file_end_pos int = 0

func free_mail_index() {
	var tmp *mail_index_type
	for mail_index != nil {
		tmp = (*mail_index_type)(unsafe.Pointer(mail_index.Next))
		if mail_index.List_start != nil {
			var (
				i *position_list_type
				j *position_list_type
			)
			i = mail_index.List_start
			for i != nil {
				j = (*position_list_type)(unsafe.Pointer(i.Next))
				libc.Free(unsafe.Pointer(i))
				i = j
			}
		}
		libc.Free(unsafe.Pointer(mail_index))
		mail_index = tmp
	}
}
func mail_recip_ok(name *byte) int {
	var (
		player_i int
		ret      int = FALSE
	)
	if (func() int {
		player_i = get_ptable_by_name(name)
		return player_i
	}()) >= 0 {
		if !IS_SET(bitvector_t(int32(player_table[player_i].Flags)), 1<<0) {
			ret = TRUE
		}
	}
	return ret
}
func push_free_list(pos int) {
	var new_pos *position_list_type
	new_pos = new(position_list_type)
	new_pos.Position = pos
	new_pos.Next = (*position_list_type_d)(unsafe.Pointer(free_list))
	free_list = new_pos
}
func pop_free_list() int {
	var (
		old_pos      *position_list_type
		return_value int
	)
	if (func() *position_list_type {
		old_pos = free_list
		return old_pos
	}()) == nil {
		return file_end_pos
	}
	return_value = free_list.Position
	free_list = (*position_list_type)(unsafe.Pointer(old_pos.Next))
	libc.Free(unsafe.Pointer(old_pos))
	return return_value
}
func clear_free_list() {
	for free_list != nil {
		pop_free_list()
	}
}
func find_char_in_index(searchee int) *mail_index_type {
	var tmp *mail_index_type
	if searchee < 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- non fatal error #1 (searchee == %ld)."), searchee)
		return nil
	}
	for tmp = mail_index; tmp != nil && tmp.Recipient != searchee; tmp = (*mail_index_type)(unsafe.Pointer(tmp.Next)) {
	}
	return tmp
}
func write_to_file(buf unsafe.Pointer, size int, filepos int) {
	var mail_file *stdio.File
	if filepos%BLOCK_SIZE != 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- fatal error #2!!! (invalid file position %ld)"), filepos)
		no_mail = TRUE
		return
	}
	if (func() *stdio.File {
		mail_file = stdio.FOpen(LIB_ETC, "r+b")
		return mail_file
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Unable to open mail file '%s'."), LIB_ETC)
		no_mail = TRUE
		return
	}
	mail_file.Seek(int64(filepos), stdio.SEEK_SET)
	mail_file.WriteN((*byte)(buf), size, 1)
	mail_file.Seek(0, stdio.SEEK_END)
	file_end_pos = int(mail_file.Tell())
	mail_file.Close()
	return
}
func read_from_file(buf unsafe.Pointer, size int, filepos int) {
	var mail_file *stdio.File
	if filepos%BLOCK_SIZE != 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- fatal error #3!!! (invalid filepos read %ld)"), filepos)
		no_mail = TRUE
		return
	}
	if (func() *stdio.File {
		mail_file = stdio.FOpen(LIB_ETC, "r+b")
		return mail_file
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Unable to open mail file '%s'."), LIB_ETC)
		no_mail = TRUE
		return
	}
	mail_file.Seek(int64(filepos), stdio.SEEK_SET)
	mail_file.ReadN((*byte)(buf), size, 1)
	mail_file.Close()
	return
}
func index_mail(id_to_index int, pos int) {
	var (
		new_index    *mail_index_type
		new_position *position_list_type
	)
	if id_to_index < 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- non-fatal error #4. (id_to_index == %ld)"), id_to_index)
		return
	}
	if (func() *mail_index_type {
		new_index = find_char_in_index(id_to_index)
		return new_index
	}()) == nil {
		new_index = new(mail_index_type)
		new_index.Recipient = id_to_index
		new_index.List_start = nil
		new_index.Next = (*mail_index_type_d)(unsafe.Pointer(mail_index))
		mail_index = new_index
	}
	new_position = new(position_list_type)
	new_position.Position = pos
	new_position.Next = (*position_list_type_d)(unsafe.Pointer(new_index.List_start))
	new_index.List_start = new_position
}
func scan_file() int {
	var (
		mail_file      *stdio.File
		next_block     header_block_type
		total_messages int = 0
		block_num      int = 0
	)
	if (func() *stdio.File {
		mail_file = stdio.FOpen(LIB_ETC, "rb")
		return mail_file
	}()) == nil {
		basic_mud_log(libc.CString("   Mail file non-existant... creating new file."))
		touch(libc.CString(LIB_ETC))
		return 1
	}
	for int(mail_file.ReadN((*byte)(unsafe.Pointer(&next_block)), int(unsafe.Sizeof(header_block_type{})), 1)) != 0 {
		if next_block.Block_type == int(-1) {
			index_mail(next_block.Header_data.To, block_num*BLOCK_SIZE)
			total_messages++
		} else if next_block.Block_type == int(-3) {
			push_free_list(block_num * BLOCK_SIZE)
		} else {
			push_free_list(block_num * BLOCK_SIZE)
		}
		block_num++
	}
	file_end_pos = int(mail_file.Tell())
	mail_file.Close()
	basic_mud_log(libc.CString("   %ld bytes read."), file_end_pos)
	if file_end_pos%BLOCK_SIZE != 0 {
		basic_mud_log(libc.CString("SYSERR: Error booting mail system -- Mail file corrupt!"))
		basic_mud_log(libc.CString("SYSERR: Mail disabled!"))
		return 0
	}
	basic_mud_log(libc.CString("   Mail file read -- %d messages."), total_messages)
	return 1
}
func has_mail(recipient int) int {
	return int(libc.BoolToInt(find_char_in_index(recipient) != nil))
}
func store_mail(to int, from int, message_pointer *byte) {
	var (
		header         header_block_type
		data           data_block_type
		last_address   int
		target_address int
		msg_txt        *byte = message_pointer
		bytes_written  int
		total_length   int = libc.StrLen(message_pointer)
	)
	if unsafe.Sizeof(header_block_type{}) != unsafe.Sizeof(data_block_type{}) || BLOCK_SIZE != unsafe.Sizeof(header_block_type{}) {
		core_dump_real(libc.CString("__FILE__"), 0)
		return
	}
	if from < 0 && from != -1337 || to < 0 || *message_pointer == 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- non-fatal error #5. (from == %ld, to == %ld)"), from, to)
		return
	}
	*(*header_block_type)(unsafe.Pointer((*byte)(unsafe.Pointer(&header)))) = header_block_type{}
	header.Block_type = -1
	header.Header_data.Next_block = -2
	header.Header_data.From = from
	header.Header_data.To = to
	header.Header_data.Mail_time = libc.GetTime(nil)
	libc.StrNCpy(&header.Txt[0], msg_txt, int(BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(header_data_type{})-unsafe.Sizeof(int8(0))))
	header.Txt[BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(header_data_type{})-unsafe.Sizeof(int8(0))] = '\x00'
	target_address = pop_free_list()
	index_mail(to, target_address)
	write_to_file(unsafe.Pointer(&header), BLOCK_SIZE, target_address)
	if libc.StrLen(msg_txt) <= int(BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(header_data_type{})-unsafe.Sizeof(int8(0))) {
		return
	}
	bytes_written = int(BLOCK_SIZE - unsafe.Sizeof(int(0)) - unsafe.Sizeof(header_data_type{}) - unsafe.Sizeof(int8(0)))
	msg_txt = (*byte)(unsafe.Add(unsafe.Pointer(msg_txt), BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(header_data_type{})-unsafe.Sizeof(int8(0))))
	last_address = target_address
	target_address = pop_free_list()
	header.Header_data.Next_block = target_address
	write_to_file(unsafe.Pointer(&header), BLOCK_SIZE, last_address)
	*(*data_block_type)(unsafe.Pointer((*byte)(unsafe.Pointer(&data)))) = data_block_type{}
	data.Block_type = -2
	libc.StrNCpy(&data.Txt[0], msg_txt, int(BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(int8(0))))
	data.Txt[BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(int8(0))] = '\x00'
	write_to_file(unsafe.Pointer(&data), BLOCK_SIZE, target_address)
	bytes_written += libc.StrLen(&data.Txt[0])
	msg_txt = (*byte)(unsafe.Add(unsafe.Pointer(msg_txt), libc.StrLen(&data.Txt[0])))
	for bytes_written < total_length {
		last_address = target_address
		target_address = pop_free_list()
		data.Block_type = target_address
		write_to_file(unsafe.Pointer(&data), BLOCK_SIZE, last_address)
		data.Block_type = -2
		libc.StrNCpy(&data.Txt[0], msg_txt, int(BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(int8(0))))
		data.Txt[BLOCK_SIZE-unsafe.Sizeof(int(0))-unsafe.Sizeof(int8(0))] = '\x00'
		write_to_file(unsafe.Pointer(&data), BLOCK_SIZE, target_address)
		bytes_written += libc.StrLen(&data.Txt[0])
		msg_txt = (*byte)(unsafe.Add(unsafe.Pointer(msg_txt), libc.StrLen(&data.Txt[0])))
	}
}
func read_delete(recipient int, from **byte) *byte {
	var (
		header           header_block_type
		data             data_block_type
		mail_pointer     *mail_index_type
		prev_mail        *mail_index_type
		position_pointer *position_list_type
		mail_address     int
		following_block  int
		tmstr            *byte
		buf              [4352]byte
		to               *byte
	)
	if recipient < 0 {
		basic_mud_log(libc.CString("SYSERR: Mail system -- non-fatal error #6. (recipient: %ld)"), recipient)
		return nil
	}
	if (func() *mail_index_type {
		mail_pointer = find_char_in_index(recipient)
		return mail_pointer
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Mail system -- post office spec_proc error?  Error #7. (invalid character in index)"))
		return nil
	}
	if (func() *position_list_type {
		position_pointer = mail_pointer.List_start
		return position_pointer
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Mail system -- non-fatal error #8. (invalid position pointer %p)"), position_pointer)
		return nil
	}
	if position_pointer.Next == nil {
		mail_address = position_pointer.Position
		libc.Free(unsafe.Pointer(position_pointer))
		if mail_index == mail_pointer {
			mail_index = (*mail_index_type)(unsafe.Pointer(mail_pointer.Next))
			libc.Free(unsafe.Pointer(mail_pointer))
		} else {
			for prev_mail = mail_index; unsafe.Pointer(prev_mail.Next) != unsafe.Pointer(mail_pointer); prev_mail = (*mail_index_type)(unsafe.Pointer(prev_mail.Next)) {
			}
			prev_mail.Next = mail_pointer.Next
			libc.Free(unsafe.Pointer(mail_pointer))
		}
	} else {
		for position_pointer.Next.Next != nil {
			position_pointer = (*position_list_type)(unsafe.Pointer(position_pointer.Next))
		}
		mail_address = position_pointer.Next.Position
		libc.Free(unsafe.Pointer(position_pointer.Next))
		position_pointer.Next = nil
	}
	read_from_file(unsafe.Pointer(&header), BLOCK_SIZE, mail_address)
	if header.Block_type != int(-1) {
		basic_mud_log(libc.CString("SYSERR: Oh dear. (Header block %ld != %d)"), header.Block_type, -1)
		no_mail = TRUE
		basic_mud_log(libc.CString("SYSERR: Mail system disabled!  -- Error #9. (Invalid header block.)"))
		return nil
	}
	tmstr = libc.AscTime(libc.LocalTime(&header.Header_data.Mail_time))
	*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(tmstr), libc.StrLen(tmstr)))), -1))) = '\x00'
	if header.Header_data.From != -1337 {
		*from = get_name_by_id(header.Header_data.From)
	} else {
		*from = libc.CString("Auctioneer")
	}
	to = get_name_by_id(recipient)
	if *from == nil {
		*from = libc.CString("Unknown")
	}
	stdio.Snprintf(&buf[0], int(4352), " @D* * * * @CGalactic Mail System @D* * * *\r\n@cDate@D:@w %s\r\n  @cTo@D:@G %s\r\n@cFrom@D:@R %s\r\n\r\n@w%s@n", tmstr, func() *byte {
		if to != nil {
			return CAP(to)
		}
		return libc.CString("Unknown")
	}(), func() *byte {
		if *from != nil {
			return CAP(*from)
		}
		return libc.CString("Unknown")
	}(), &header.Txt[0])
	following_block = header.Header_data.Next_block
	header.Block_type = -3
	write_to_file(unsafe.Pointer(&header), BLOCK_SIZE, mail_address)
	push_free_list(mail_address)
	for following_block != int(-2) {
		read_from_file(unsafe.Pointer(&data), BLOCK_SIZE, following_block)
		libc.StrCat(&buf[0], &data.Txt[0])
		mail_address = following_block
		following_block = data.Block_type
		data.Block_type = -3
		write_to_file(unsafe.Pointer(&data), BLOCK_SIZE, mail_address)
		push_free_list(mail_address)
	}
	return libc.StrDup(&buf[0])
}
func postmaster(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	if ch.Desc == nil || IS_NPC(ch) {
		return 0
	}
	if libc.StrCmp(libc.CString("mail"), complete_cmd_info[cmd].Command) != 0 && libc.StrCmp(libc.CString("check"), complete_cmd_info[cmd].Command) != 0 && libc.StrCmp(libc.CString("receive"), complete_cmd_info[cmd].Command) != 0 {
		return 0
	}
	if no_mail != 0 {
		send_to_char(ch, libc.CString("Sorry, the mail system is having technical difficulties.\r\n"))
		return 0
	}
	if libc.StrCmp(libc.CString("mail"), complete_cmd_info[cmd].Command) == 0 {
		postmaster_send_mail(ch, (*char_data)(me), cmd, argument)
		return 1
	} else if libc.StrCmp(libc.CString("check"), complete_cmd_info[cmd].Command) == 0 {
		postmaster_check_mail(ch, (*char_data)(me), cmd, argument)
		return 1
	} else if libc.StrCmp(libc.CString("receive"), complete_cmd_info[cmd].Command) == 0 {
		postmaster_receive_mail(ch, (*char_data)(me), cmd, argument)
		return 1
	} else {
		return 0
	}
}
func postmaster_send_mail(ch *char_data, mailman *char_data, cmd int, arg *byte) {
	var (
		recipient int
		buf       [2048]byte
		mailwrite **byte
	)
	if GET_LEVEL(ch) < MIN_MAIL_LEVEL && ch.Admlevel < ADMLVL_IMMORT {
		stdio.Snprintf(&buf[0], int(2048), "$n tells you, 'Sorry, you have to be level %d to send mail!'", MIN_MAIL_LEVEL)
		act(&buf[0], FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
		return
	}
	one_argument(arg, &buf[0])
	if buf[0] == 0 {
		act(libc.CString("$n tells you, 'You need to specify an addressee!'"), FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
		return
	}
	if ch.Gold < STAMP_PRICE && !ADM_FLAGGED(ch, ADM_MONEY) {
		stdio.Snprintf(&buf[0], int(2048), "$n tells you, 'A stamp costs %d zenni.'\r\n$n tells you, '...which I see you can't afford.'", STAMP_PRICE)
		act(&buf[0], FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
		return
	}
	if (func() int {
		recipient = get_id_by_name(&buf[0])
		return recipient
	}()) < 0 || mail_recip_ok(&buf[0]) == 0 {
		act(libc.CString("$n tells you, 'No one by that name is registered here!'"), FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
		return
	}
	act(libc.CString("$n starts to write some mail."), TRUE, ch, nil, nil, TO_ROOM)
	stdio.Snprintf(&buf[0], int(2048), "$n tells you, 'I'll take %d zenni for the stamp.'\r\n$n tells you, 'Write your message. (/s saves /h for help).'", STAMP_PRICE)
	act(&buf[0], FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
	act(libc.CString("@C$n@w starts writing a letter.@n"), TRUE, ch, nil, nil, TO_ROOM)
	ch.Gold -= STAMP_PRICE
	SET_BIT_AR(ch.Act[:], PLR_MAILING)
	mailwrite = new(*byte)
	string_write(ch.Desc, mailwrite, MAX_MAIL_SIZE, recipient, nil)
}
func postmaster_check_mail(ch *char_data, mailman *char_data, cmd int, arg *byte) {
	if has_mail(int(ch.Idnum)) != 0 {
		act(libc.CString("$n tells you, 'You have mail waiting.'"), FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
	} else {
		act(libc.CString("$n tells you, 'Sorry, you don't have any mail waiting.'"), FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
	}
}
func postmaster_receive_mail(ch *char_data, mailman *char_data, cmd int, arg *byte) {
	var (
		buf  [256]byte
		obj  *obj_data
		y    int
		from *byte
	)
	if has_mail(int(ch.Idnum)) == 0 && mailman != nil {
		stdio.Snprintf(&buf[0], int(256), "$n tells you, 'Sorry, you don't have any mail waiting.'")
		act(&buf[0], FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
		return
	}
	for has_mail(int(ch.Idnum)) != 0 {
		obj = create_obj()
		obj.Item_number = -1
		obj.Type_flag = ITEM_NOTE
		for y = 0; y < TW_ARRAY_MAX; y++ {
			obj.Wear_flags[y] = 0
		}
		SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_TAKE)
		SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_HOLD)
		obj.Weight = 1
		obj.Cost = 30
		obj.Cost_per_day = 10
		obj.Action_description = read_delete(int(ch.Idnum), &from)
		var bla [256]byte
		var blm [256]byte
		stdio.Sprintf(&bla[0], "@WA piece of mail@n")
		stdio.Sprintf(&blm[0], "@WSomeone has left a piece of mail here@n")
		obj.Short_description = libc.StrDup(&bla[0])
		obj.Description = libc.StrDup(&blm[0])
		stdio.Sprintf(&bla[0], "mail paper letter")
		obj.Name = libc.StrDup(&bla[0])
		bla[0] = '\x00'
		blm[0] = '\x00'
		SET_BIT_AR(obj.Extra_flags[:], ITEM_UNIQUE_SAVE)
		add_unique_id(obj)
		obj.Value[VAL_ALL_MATERIAL] = MATERIAL_PAPER
		obj.Value[VAL_NOTE_HEALTH] = 100
		obj.Value[VAL_NOTE_MAXHEALTH] = 100
		if obj.Action_description == nil {
			obj.Action_description = libc.CString("Mail system error - please report.  Error #11.\r\n")
		}
		SET_BIT_AR(obj.Extra_flags[:], ITEM_UNIQUE_SAVE)
		if IS_PLAYING(ch.Desc) && mailman != nil {
			obj_to_char(obj, ch)
			act(libc.CString("$n gives you a piece of mail."), FALSE, mailman, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("$N gives $n a piece of mail."), FALSE, ch, nil, unsafe.Pointer(mailman), TO_ROOM)
		} else {
			extract_obj(obj)
		}
	}
}
func notify_if_playing(from *char_data, recipient_id int) {
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if IS_PLAYING(d) && int(d.Character.Idnum) == recipient_id && has_mail(int(d.Character.Idnum)) != 0 {
			send_to_char(d.Character, libc.CString("\r\n\a\a\a@G@lYou have new mudmail from %s.@n\r\n"), GET_NAME(from))
		}
	}
}
