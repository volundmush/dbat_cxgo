package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const MAX_INVALID_NAMES = 200

var ban_list *ban_list_element = nil
var ban_types [5]*byte = [5]*byte{libc.CString("no"), libc.CString("new"), libc.CString("select"), libc.CString("all"), libc.CString("ERROR")}

func load_banned() {
	var (
		fl        *C.FILE
		i         int
		date      int
		site_name [51]byte
		ban_type  [100]byte
		name      [21]byte
		next_node *ban_list_element
	)
	ban_list = nil
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_ETC, "r")))
		return fl
	}()) == nil {
		if (*__errno_location()) != ENOENT {
			basic_mud_log(libc.CString("SYSERR: Unable to open banfile '%s': %s"), LIB_ETC, C.strerror(*__errno_location()))
		} else {
			basic_mud_log(libc.CString("   Ban file '%s' doesn't exist."), LIB_ETC)
		}
		return
	}
	for __isoc99_fscanf(fl, libc.CString(" %s %s %d %s "), &ban_type[0], &site_name[0], &date, &name[0]) == 4 {
		next_node = new(ban_list_element)
		C.strncpy(&next_node.Site[0], &site_name[0], BANNED_SITE_LENGTH)
		next_node.Site[BANNED_SITE_LENGTH] = '\x00'
		C.strncpy(&next_node.Name[0], &name[0], MAX_NAME_LENGTH)
		next_node.Name[MAX_NAME_LENGTH] = '\x00'
		next_node.Date = int64(date)
		for i = BAN_NOT; i <= BAN_ALL; i++ {
			if C.strcmp(&ban_type[0], ban_types[i]) == 0 {
				next_node.Type = i
			}
		}
		next_node.Next = ban_list
		ban_list = next_node
	}
	C.fclose(fl)
}
func isbanned(hostname *byte) int {
	var (
		i           int
		banned_node *ban_list_element
		nextchar    *byte
	)
	if hostname == nil || *hostname == 0 {
		return 0
	}
	i = 0
	for nextchar = hostname; *nextchar != 0; nextchar = (*byte)(unsafe.Add(unsafe.Pointer(nextchar), 1)) {
		*nextchar = byte(int8(C.tolower(int(*nextchar))))
	}
	for banned_node = ban_list; banned_node != nil; banned_node = banned_node.Next {
		if C.strstr(hostname, &banned_node.Site[0]) != nil {
			i = MAX(i, banned_node.Type)
		}
	}
	return i
}
func _write_one_node(fp *C.FILE, node *ban_list_element) {
	if node != nil {
		_write_one_node(fp, node.Next)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s %s %ld %s\n", ban_types[node.Type], &node.Site[0], int(node.Date), &node.Name[0])
	}
}
func write_ban_list() {
	var fl *C.FILE
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_ETC, "w")))
		return fl
	}()) == nil {
		C.perror(libc.CString("SYSERR: Unable to open 'etc/badsites' for writing"))
		return
	}
	_write_one_node(fl, ban_list)
	C.fclose(fl)
	return
}
func do_ban(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		flag     [2048]byte
		site     [2048]byte
		nextchar *byte
		timestr  [16]byte
		i        int
		ban_node *ban_list_element
	)
	if *argument == 0 {
		if ban_list == nil {
			send_to_char(ch, libc.CString("No sites are banned.\r\n"))
			return
		}
		send_to_char(ch, libc.CString(BAN_LIST_FORMAT), "Banned Site Name", "Ban Type", "Banned On", "Banned By")
		send_to_char(ch, libc.CString(BAN_LIST_FORMAT), "---------------------------------", "---------------------------------", "---------------------------------", "---------------------------------")
		for ban_node = ban_list; ban_node != nil; ban_node = ban_node.Next {
			if ban_node.Date != 0 {
				strlcpy(&timestr[0], C.asctime(C.localtime(&ban_node.Date)), 10)
				timestr[10] = '\x00'
			} else {
				C.strcpy(&timestr[0], libc.CString("Unknown"))
			}
			send_to_char(ch, libc.CString(BAN_LIST_FORMAT), &ban_node.Site[0], ban_types[ban_node.Type], &timestr[0], &ban_node.Name[0])
		}
		return
	}
	two_arguments(argument, &flag[0], &site[0])
	if site[0] == 0 || flag[0] == 0 {
		send_to_char(ch, libc.CString("Usage: ban {all | select | new} site_name\r\n"))
		return
	}
	if C.strcasecmp(&flag[0], libc.CString("select")) != 0 && C.strcasecmp(&flag[0], libc.CString("all")) != 0 && C.strcasecmp(&flag[0], libc.CString("new")) != 0 {
		send_to_char(ch, libc.CString("Flag must be ALL, SELECT, or NEW.\r\n"))
		return
	}
	for ban_node = ban_list; ban_node != nil; ban_node = ban_node.Next {
		if C.strcasecmp(&ban_node.Site[0], &site[0]) == 0 {
			send_to_char(ch, libc.CString("That site has already been banned -- unban it to change the ban type.\r\n"))
			return
		}
	}
	ban_node = new(ban_list_element)
	C.strncpy(&ban_node.Site[0], &site[0], BANNED_SITE_LENGTH)
	for nextchar = &ban_node.Site[0]; *nextchar != 0; nextchar = (*byte)(unsafe.Add(unsafe.Pointer(nextchar), 1)) {
		*nextchar = byte(int8(C.tolower(int(*nextchar))))
	}
	ban_node.Site[BANNED_SITE_LENGTH] = '\x00'
	C.strncpy(&ban_node.Name[0], GET_NAME(ch), MAX_NAME_LENGTH)
	ban_node.Name[MAX_NAME_LENGTH] = '\x00'
	ban_node.Date = C.time(nil)
	for i = BAN_NEW; i <= BAN_ALL; i++ {
		if C.strcasecmp(&flag[0], ban_types[i]) == 0 {
			ban_node.Type = i
		}
	}
	ban_node.Next = ban_list
	ban_list = ban_node
	mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s has banned %s for %s players."), GET_NAME(ch), &site[0], ban_types[ban_node.Type])
	send_to_char(ch, libc.CString("Site banned.\r\n"))
	write_ban_list()
}
func do_unban(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		site     [2048]byte
		ban_node *ban_list_element
		temp     *ban_list_element
		found    int = 0
	)
	one_argument(argument, &site[0])
	if site[0] == 0 {
		send_to_char(ch, libc.CString("A site to unban might help.\r\n"))
		return
	}
	ban_node = ban_list
	for ban_node != nil && found == 0 {
		if C.strcasecmp(&ban_node.Site[0], &site[0]) == 0 {
			found = 1
		} else {
			ban_node = ban_node.Next
		}
	}
	if found == 0 {
		send_to_char(ch, libc.CString("That site is not currently banned.\r\n"))
		return
	}
	if ban_node == ban_list {
		ban_list = ban_node.Next
	} else {
		temp = ban_list
		for temp != nil && temp.Next != ban_node {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = ban_node.Next
		}
	}
	send_to_char(ch, libc.CString("Site unbanned.\r\n"))
	mudlog(NRM, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s removed the %s-player ban on %s."), GET_NAME(ch), ban_types[ban_node.Type], &ban_node.Site[0])
	libc.Free(unsafe.Pointer(ban_node))
	write_ban_list()
}

var invalid_list [200]*byte
var num_invalid int = 0

func Valid_Name(newname *byte) int {
	var (
		i        int
		wovels   int = 0
		dt       *descriptor_data
		tempname [2048]byte
	)
	for dt = descriptor_list; dt != nil; dt = dt.Next {
		if dt.Character != nil && GET_NAME(dt.Character) != nil && C.strcasecmp(GET_NAME(dt.Character), newname) == 0 {
			if dt.Character.Idnum == -1 {
				return int(libc.BoolToInt(IS_PLAYING(dt)))
			}
		}
	}
	for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(newname), i)) != 0; i++ {
		if C.strchr(libc.CString("aeiouyAEIOUY"), int(*(*byte)(unsafe.Add(unsafe.Pointer(newname), i)))) != nil {
			wovels++
		}
	}
	if wovels == 0 {
		return 0
	}
	if invalid_list == nil || num_invalid < 1 {
		return 1
	}
	strlcpy(&tempname[0], newname, uint64(2048))
	for i = 0; tempname[i] != 0; i++ {
		tempname[i] = byte(int8(C.tolower(int(tempname[i]))))
	}
	for i = 0; i < num_invalid; i++ {
		if C.strstr(&tempname[0], invalid_list[i]) != nil {
			return 0
		}
	}
	return 1
}
func Free_Invalid_List() {
	var invl int
	for invl = 0; invl < num_invalid; invl++ {
		libc.Free(unsafe.Pointer(invalid_list[invl]))
	}
	num_invalid = 0
}
func Read_Invalid_List() {
	var (
		fp   *C.FILE
		temp [256]byte
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_MISC, "r")))
		return fp
	}()) == nil {
		C.perror(libc.CString("SYSERR: Unable to open 'misc/xnames' for reading"))
		return
	}
	num_invalid = 0
	for get_line(fp, &temp[0]) != 0 && num_invalid < MAX_INVALID_NAMES {
		invalid_list[func() int {
			p := &num_invalid
			x := *p
			*p++
			return x
		}()] = C.strdup(&temp[0])
	}
	if num_invalid >= MAX_INVALID_NAMES {
		basic_mud_log(libc.CString("SYSERR: Too many invalid names; change MAX_INVALID_NAMES in ban.c"))
		C.exit(1)
	}
	C.fclose(fp)
}
