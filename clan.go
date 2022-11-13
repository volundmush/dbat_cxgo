package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const LVL_CLAN_MOD = 32
const DEFAULT_OPEN_JOIN = 0
const DEFAULT_OPEN_LEAVE = 0
const DEFAULT_CLAN_INFO = "little is known about this clan, currently."

var num_clans int = 0
var clan **clan_data = nil

type clan_member struct {
	Next *clan_member
	Id   int
}
type clan_data struct {
	Name       *byte
	Info       *byte
	Highrank   *byte
	Midrank    *byte
	Modlist    [1000]byte
	Memlist    [1000]byte
	Applist    [1000]byte
	Moderators *clan_member
	Members    *clan_member
	Applicants *clan_member
	Open_join  int
	Open_leave int
	Bank       int
	Bany       int
}

func clanMemberFromList(id int, list *clan_member) *clan_member {
	for ; list != nil; list = list.Next {
		if id == list.Id {
			return list
		}
	}
	return list
}
func writeClanMasterlist() {
	var (
		i   int
		fl  *stdio.File
		buf [64936]byte
	)
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not open clan masterlist for writing."))
		return
	}
	stdio.Sprintf(&buf[0], "%d\n", num_clans)
	stdio.Fprintf(fl, libc.GoString(&buf[0]))
	for i = 0; i < num_clans; i++ {
		stdio.Fprintf(fl, "%s%d.cla\n", LIB_ETC, i)
	}
	fl.Close()
}
func fgetlinetomax(file *stdio.File, p *byte, maxlen int) int {
	var count int = 0
	for int(file.IsEOF()) == 0 && count < maxlen-1 {
		*(*byte)(unsafe.Add(unsafe.Pointer(p), count)) = byte(int8(file.GetC()))
		if *(*byte)(unsafe.Add(unsafe.Pointer(p), count)) == '\n' {
			break
		}
		count++
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(p), count)) = '\x00'
	return count
}
func clanFilename(S *clan_data) *byte {
	var i int
	for i = 0; i < num_clans; i++ {
		if libc.StrCmp(S.Name, (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name) == 0 {
			break
		}
	}
	if i == num_clans {
		return nil
	} else {
		var buf [64936]byte
		stdio.Sprintf(&buf[0], "%s%d.cla", LIB_ETC, i)
		return libc.StrDup(&buf[0])
	}
}
func clanLoad(filename *byte) *clan_data {
	var (
		fl   *stdio.File
		line [64936]byte
		info *byte
	)
	_ = info
	var id int
	var infolen int
	_ = infolen
	var S *clan_data
	if filename == nil {
		basic_mud_log(libc.CString("ERROR: passed null pointer to clanLoad"))
		return nil
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(filename), "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not open file, %s, in clanLoad."), filename)
		return nil
	}
	S = new(clan_data)
	stdio.Sprintf(&S.Modlist[0], "@D---@CLeaders@D---\n")
	stdio.Sprintf(&S.Memlist[0], "@D---@cMembers@D---\n")
	stdio.Sprintf(&S.Applist[0], "@D---@YApplicants@D---\n")
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	stdio.Sscanf(&line[0], "%d %d", &S.Open_join, &S.Open_leave)
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	stdio.Sscanf(&line[0], "%ld", &S.Bank)
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	stdio.Sscanf(&line[0], "%d", &S.Bany)
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	S.Name = libc.StrDup(&line[0])
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	S.Highrank = libc.StrDup(&line[0])
	fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	S.Midrank = libc.StrDup(&line[0])
	var memcount int = 0
	for TRUE != 0 {
		var moderator *clan_member
		fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
		if libc.StrCmp(&line[0], libc.CString("~")) == 0 {
			break
		}
		stdio.Sscanf(&line[0], "%d", &id)
		moderator = new(clan_member)
		moderator.Id = id
		moderator.Next = S.Moderators
		S.Moderators = moderator
		if get_name_by_id(id) != nil {
			memcount += 1
			stdio.Sprintf(&S.Modlist[libc.StrLen(&S.Modlist[0])], "@D[@G%2d@D]@W %s\n", memcount, get_name_by_id(id))
		}
	}
	for TRUE != 0 {
		var member *clan_member
		fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
		if libc.StrCmp(&line[0], libc.CString("~")) == 0 {
			break
		}
		stdio.Sscanf(&line[0], "%d", &id)
		member = new(clan_member)
		member.Id = id
		member.Next = S.Members
		S.Members = member
		if get_name_by_id(id) != nil {
			memcount += 1
			stdio.Sprintf(&S.Memlist[libc.StrLen(&S.Memlist[0])], "@D[@G%2d@D]@W %s\n", memcount, get_name_by_id(id))
		}
	}
	for TRUE != 0 {
		var applicant *clan_member
		fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
		if libc.StrCmp(&line[0], libc.CString("~")) == 0 {
			break
		}
		stdio.Sscanf(&line[0], "%d", &id)
		applicant = new(clan_member)
		applicant.Id = id
		applicant.Next = S.Applicants
		S.Applicants = applicant
		if get_name_by_id(id) != nil {
			stdio.Sprintf(&S.Applist[libc.StrLen(&S.Applist[0])], "@W%s\n", get_name_by_id(id))
		}
	}
	infolen = 0
	info = libc.CString("")
	libc.StrCpy(&line[0], libc.CString(""))
	S.Info = fread_string(fl, &line[0])
	if libc.StrLen(&line[0]) > 0 {
	}
	fl.Close()
	return S
}
func clanSave(S *clan_data, filename *byte) bool {
	var (
		fl   *stdio.File
		list *clan_member
	)
	if filename == nil {
		basic_mud_log(libc.CString("ERROR: passed null pointer to clanSave when saving %s"), S.Name)
		return FALSE != 0
	}
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(filename), "w")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("ERROR: could not save clan, %s, to filename, %s."), S.Name, filename)
		return FALSE != 0
	}
	stdio.Fprintf(fl, "%d %d\n", S.Open_join, S.Open_leave)
	stdio.Fprintf(fl, "%ld\n", S.Bank)
	stdio.Fprintf(fl, "%d\n", S.Bany)
	stdio.Fprintf(fl, "%s\n", S.Name)
	stdio.Fprintf(fl, "%s\n", S.Highrank)
	stdio.Fprintf(fl, "%s\n", S.Midrank)
	for list = S.Moderators; list != nil; list = list.Next {
		stdio.Fprintf(fl, "%d\n", list.Id)
	}
	stdio.Fprintf(fl, "~\n")
	for list = S.Members; list != nil; list = list.Next {
		stdio.Fprintf(fl, "%d\n", list.Id)
	}
	stdio.Fprintf(fl, "~\n")
	for list = S.Applicants; list != nil; list = list.Next {
		stdio.Fprintf(fl, "%d\n", list.Id)
	}
	stdio.Fprintf(fl, "~\n")
	stdio.Fprintf(fl, "%s~\n", S.Info)
	fl.Close()
	return TRUE != 0
}
func clanDelete(S *clan_data) {
	var (
		next   *clan_member
		member *clan_member
	)
	if S.Moderators != nil {
		for member = S.Moderators; member != nil; member = next {
			next = member.Next
			libc.Free(unsafe.Pointer(member))
		}
	}
	if S.Members != nil {
		for member = S.Members; member != nil; member = next {
			next = member.Next
			libc.Free(unsafe.Pointer(member))
		}
	}
	if S.Applicants != nil {
		for member = S.Applicants; member != nil; member = next {
			next = member.Next
			libc.Free(unsafe.Pointer(member))
		}
	}
	libc.Free(unsafe.Pointer(S.Name))
	libc.Free(unsafe.Pointer(S.Info))
	if S.Highrank != nil {
		libc.Free(unsafe.Pointer(S.Highrank))
	}
	if S.Midrank != nil {
		libc.Free(unsafe.Pointer(S.Midrank))
	}
	libc.Free(unsafe.Pointer(S))
}
func clanRemove(S *clan_data) {
	var (
		i int
		j int
	)
	for i = 0; i < num_clans; i++ {
		if *(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))) == S {
			break
		}
	}
	if i == num_clans {
		basic_mud_log(libc.CString("ERROR: tried to remove clan, %s, which did not formally exist."), S.Name)
		clanDelete(S)
		return
	}
	num_clans--
	for j = i; j < num_clans; j++ {
		*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(j))) = *(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(j+1)))
	}
	for ; i < num_clans; i++ {
		clanSave(*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))), clanFilename(*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))))
	}
	clanDelete(S)
	writeClanMasterlist()
}
func clanAdd(S *clan_data) {
	var (
		i       int
		oldList **clan_data = clan
	)
	clan = (**clan_data)(libc.Malloc((num_clans + 1) * int(unsafe.Sizeof((*clan_data)(nil)))))
	for i = 0; i < num_clans; i++ {
		*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))) = *(**clan_data)(unsafe.Add(unsafe.Pointer(oldList), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))
	}
	*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(num_clans))) = S
	num_clans++
	clanSave(S, clanFilename(S))
	libc.Free(unsafe.Pointer(oldList))
	writeClanMasterlist()
}
func clanGet(name *byte) *clan_data {
	var (
		i       int
		newname *byte = strlwr(libc.StrDup(name))
	)
	for i = 0; i < num_clans; i++ {
		if libc.StrCmp(newname, strlwr(libc.StrDup((*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name))) == 0 {
			libc.Free(unsafe.Pointer(newname))
			return *(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))
		}
	}
	libc.Free(unsafe.Pointer(newname))
	return nil
}
func clanReload(name *byte) bool {
	var (
		i int
		S *clan_data
	)
	if (func() *clan_data {
		S = clanGet(name)
		return S
	}()) == nil {
		return FALSE != 0
	}
	for i = 0; i < num_clans; i++ {
		if S == *(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))) {
			var buf [64936]byte
			clanDelete(*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))))
			stdio.Sprintf(&buf[0], "%s%d.cla", LIB_ETC, i)
			*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))) = clanLoad(&buf[0])
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func clanBoot() {
	var (
		fl   *stdio.File
		i    int
		len_ int
		line [64936]byte
	)
	if (func() *stdio.File {
		fl = stdio.FOpen(LIB_ETC, "r")
		return fl
	}()) == nil {
		basic_mud_log(libc.CString("  Could not open clan masterlist. Aborting."))
		return
	}
	if int(fl.IsEOF()) != 0 {
		basic_mud_log(libc.CString("  Clan masterlist contained no data! Aborting."))
		return
	}
	len_ = fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
	stdio.Sscanf(&line[0], "%d", &num_clans)
	if num_clans <= 0 {
		basic_mud_log(libc.CString("  No clans have formed yet."))
		clan = nil
		return
	}
	clan = (**clan_data)(libc.Malloc(num_clans * int(unsafe.Sizeof((*clan_data)(nil)))))
	for i = 0; i < num_clans; i++ {
		if (func() int {
			len_ = fgetlinetomax(fl, &line[0], MAX_STRING_LENGTH)
			return len_
		}()) > 0 {
			basic_mud_log(libc.CString("  Loading clan: %s"), &line[0])
			*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))) = clanLoad(&line[0])
		} else {
			basic_mud_log(libc.CString("  Found blank line while looking for clan names. Aborting."))
			for i--; i >= 0; i-- {
				clanDelete(*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i))))
			}
			libc.Free(unsafe.Pointer(clan))
			clan = nil
			num_clans = 0
			fl.Close()
			return
		}
	}
	fl.Close()
}
func isClan(name *byte) bool {
	return clanGet(name) != nil
}
func clanCreate(name *byte) bool {
	var S *clan_data
	if isClan(name) {
		return FALSE != 0
	}
	S = new(clan_data)
	S.Name = libc.StrDup(name)
	S.Info = libc.CString(DEFAULT_CLAN_INFO)
	S.Moderators = nil
	S.Members = nil
	S.Applicants = nil
	S.Highrank = libc.CString("Captain")
	S.Midrank = libc.CString("Lieutenant")
	clanAdd(S)
	return TRUE != 0
}
func clanINFOW(name *byte, ch *char_data) {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return
	} else {
		var backstr *byte = nil
		act(libc.CString("$n begins to edit a clan's info."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
		send_editor_help(ch.Desc)
		write_to_output(ch.Desc, libc.CString("@rYou are limited to 1000 characters for the clan info.@n\r\n"))
		backstr = libc.StrDup(S.Info)
		write_to_output(ch.Desc, libc.CString("%s\r\n"), S.Info)
		string_write(ch.Desc, &S.Info, 1000, 0, unsafe.Pointer(backstr))
		clanSave(S, clanFilename(S))
	}
}
func clan_update() {
	var i int
	if num_clans < 1 {
		return
	}
	for i = 0; i < num_clans; i++ {
		clanSAFE((*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name)
	}
	return
}
func clanSAFE(name *byte) {
	var S *clan_data = clanGet(name)
	if S == nil {
		return
	} else {
		clanSave(S, clanFilename(S))
	}
}
func clanDestroy(name *byte) {
	var S *clan_data = clanGet(name)
	if S != nil {
		clanRemove(S)
	}
}
func clanApply(name *byte, ch *char_data) bool {
	var (
		buf [2048]byte
		S   *clan_data = clanGet(name)
	)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if clanMemberFromList(int(ch.Idnum), S.Moderators) != nil || clanMemberFromList(int(ch.Idnum), S.Members) != nil {
		return FALSE != 0
	}
	if clanMemberFromList(int(ch.Idnum), S.Applicants) != nil {
		return TRUE != 0
	}
	var new *clan_member
	new = new(clan_member)
	new.Id = int(ch.Idnum)
	stdio.Sprintf(&buf[0], "Applying for %s", S.Name)
	set_clan(ch, &buf[0])
	new.Next = S.Applicants
	S.Applicants = new
	clanSave(S, clanFilename(S))
	return TRUE != 0
}
func clanHIGHRANK(name *byte, ch *char_data, rank *byte) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	} else {
		S.Highrank = libc.StrDup(rank)
		clanSave(S, clanFilename(S))
		return TRUE != 0
	}
}
func clanMIDRANK(name *byte, ch *char_data, rank *byte) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	} else {
		S.Midrank = libc.StrDup(rank)
		clanSave(S, clanFilename(S))
		return TRUE != 0
	}
}
func clanRANK(name *byte, ch *char_data, vict *char_data, num int) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	} else {
		vict.Crank = num
		return TRUE != 0
	}
}
func clanRANKD(name *byte, ch *char_data, vict *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	} else {
		send_to_char(ch, libc.CString("@cClan Rank@D: @w"))
		if vict.Crank == 0 && clanMemberFromList(int(vict.Idnum), S.Moderators) == nil {
			send_to_char(ch, libc.CString("Member@n\r\n"))
		} else if vict.Crank == 1 && clanMemberFromList(int(vict.Idnum), S.Moderators) == nil {
			send_to_char(ch, libc.CString("%s@n\r\n"), S.Midrank)
		} else if vict.Crank == 2 && clanMemberFromList(int(vict.Idnum), S.Moderators) == nil {
			send_to_char(ch, libc.CString("%s@n\r\n"), S.Highrank)
		} else {
			send_to_char(ch, libc.CString("Leader@n\r\n"))
		}
		return TRUE != 0
	}
}
func clanBANY(name *byte, ch *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if S.Bany <= 0 {
		return FALSE != 0
	} else {
		return TRUE != 0
	}
}
func clanBSET(name *byte, ch *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if S.Bany > 0 {
		S.Bany = 0
		send_to_char(ch, libc.CString("The clan bank will now only be accessible from its room.\r\n"))
		clanSave(S, clanFilename(S))
		return TRUE != 0
	} else {
		S.Bany = 1
		send_to_char(ch, libc.CString("The clan bank will now be accessible from anywhere.\r\n"))
		clanSave(S, clanFilename(S))
		return TRUE != 0
	}
}
func clanBANKADD(name *byte, ch *char_data, amt int) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	} else {
		S.Bank += amt
		clanSave(S, clanFilename(S))
		return TRUE != 0
	}
}
func clanBANK(name *byte, ch *char_data) int {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE
	} else {
		var amt int = 0
		amt = S.Bank
		return amt
	}
}
func clanBANKSUB(name *byte, ch *char_data, amt int) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if S.Bank-amt < 0 {
		return FALSE != 0
	} else {
		S.Bank -= amt
		clanSave(S, clanFilename(S))
		return TRUE != 0
	}
}
func clanInduct(name *byte, ch *char_data) bool {
	var (
		m    *clan_member
		temp *clan_member
		S    *clan_data = clanGet(name)
		buf  [2048]byte
	)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if clanMemberFromList(int(ch.Idnum), S.Moderators) != nil || clanMemberFromList(int(ch.Idnum), S.Members) != nil {
		return TRUE != 0
	}
	if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Applicants)
		return m
	}()) != nil {
		if m == S.Applicants {
			S.Applicants = m.Next
		} else {
			temp = S.Applicants
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	}
	stdio.Sprintf(&buf[0], "%s", S.Name)
	set_clan(ch, &buf[0])
	m = new(clan_member)
	m.Id = int(ch.Idnum)
	m.Next = S.Members
	S.Members = m
	clanSave(S, clanFilename(S))
	clanReload(name)
	return TRUE != 0
}
func set_clan(ch *char_data, clan *byte) {
	if ch.Clan != nil {
		libc.Free(unsafe.Pointer(ch.Clan))
	}
	ch.Clan = libc.StrDup(clan)
	ch.Crank = 0
}
func remove_clan(ch *char_data) {
	if ch.Clan != nil {
		libc.Free(unsafe.Pointer(ch.Clan))
	}
	ch.Clan = libc.CString("None.")
}
func clanMakeModerator(name *byte, ch *char_data) bool {
	var (
		m    *clan_member
		temp *clan_member
		S    *clan_data = clanGet(name)
		buf  [2048]byte
	)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if clanMemberFromList(int(ch.Idnum), S.Moderators) != nil {
		return TRUE != 0
	}
	if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Members)
		return m
	}()) != nil {
		if m == S.Members {
			S.Members = m.Next
		} else {
			temp = S.Members
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	} else if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Applicants)
		return m
	}()) != nil {
		if m == S.Applicants {
			S.Applicants = m.Next
		} else {
			temp = S.Applicants
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	}
	stdio.Sprintf(&buf[0], "%s", S.Name)
	m = new(clan_member)
	set_clan(ch, &buf[0])
	m.Id = int(ch.Idnum)
	m.Next = S.Moderators
	S.Moderators = m
	clanSave(S, clanFilename(S))
	return TRUE != 0
}
func clanExpel(name *byte, ch *char_data) {
	var (
		m    *clan_member
		temp *clan_member
		S    *clan_data = clanGet(name)
	)
	if S == nil || IS_NPC(ch) {
		return
	}
	remove_clan(ch)
	if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Moderators)
		return m
	}()) != nil {
		if m == S.Moderators {
			S.Moderators = m.Next
		} else {
			temp = S.Moderators
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	} else if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Members)
		return m
	}()) != nil {
		if m == S.Members {
			S.Members = m.Next
		} else {
			temp = S.Members
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	}
	clanSave(S, clanFilename(S))
	clanReload(name)
}
func clanDecline(name *byte, ch *char_data) {
	var (
		m    *clan_member
		temp *clan_member
		S    *clan_data = clanGet(name)
	)
	if S == nil || IS_NPC(ch) {
		return
	}
	if (func() *clan_member {
		m = clanMemberFromList(int(ch.Idnum), S.Applicants)
		return m
	}()) != nil {
		if m == S.Applicants {
			S.Applicants = m.Next
		} else {
			temp = S.Applicants
			for temp != nil && temp.Next != m {
				temp = temp.Next
			}
			if temp != nil {
				temp.Next = m.Next
			}
		}
		libc.Free(unsafe.Pointer(m))
	}
	clanSave(S, clanFilename(S))
}
func handle_clan_member_list(ch *char_data) {
	if IS_NPC(ch) {
		return
	}
	if ch.Clan == nil || ch.Clan == nil {
		return
	}
	if libc.StrStr(ch.Clan, libc.CString("None")) != nil {
		return
	}
	var S *clan_data = clanGet(ch.Clan)
	if S == nil {
		return
	}
	send_to_char(ch, &S.Modlist[0])
	send_to_char(ch, &S.Memlist[0])
	send_to_char(ch, &S.Applist[0])
	send_to_char(ch, libc.CString("@n"))
}
func clanIsMember(name *byte, ch *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	if clanMemberFromList(int(ch.Idnum), S.Moderators) != nil || clanMemberFromList(int(ch.Idnum), S.Members) != nil {
		return TRUE != 0
	} else {
		return FALSE != 0
	}
}
func clanIsModerator(name *byte, ch *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	return clanMemberFromList(int(ch.Idnum), S.Moderators) != nil
}
func clanIsApplicant(name *byte, ch *char_data) bool {
	var S *clan_data = clanGet(name)
	if S == nil || IS_NPC(ch) {
		return FALSE != 0
	}
	return clanMemberFromList(int(ch.Idnum), S.Applicants) != nil
}
func clanOpenJoin(name *byte) bool {
	var S *clan_data = clanGet(name)
	return S != nil && S.Open_join == TRUE
}
func clanOpenLeave(name *byte) bool {
	var S *clan_data = clanGet(name)
	return S != nil && S.Open_leave == TRUE
}
func clanSetOpenJoin(name *byte, val int) bool {
	var S *clan_data = clanGet(name)
	if S == nil {
		return FALSE != 0
	}
	if val == FALSE {
		S.Open_join = FALSE
	} else {
		S.Open_join = TRUE
	}
	clanSave(S, clanFilename(S))
	return TRUE != 0
}
func clanSetOpenLeave(name *byte, val int) bool {
	var S *clan_data = clanGet(name)
	if S == nil {
		return FALSE != 0
	}
	if val == FALSE {
		S.Open_leave = FALSE
	} else {
		S.Open_leave = TRUE
	}
	clanSave(S, clanFilename(S))
	return TRUE != 0
}
func listClanInfo(name *byte, ch *char_data) {
	var S *clan_data = clanGet(name)
	if S == nil {
		send_to_char(ch, libc.CString("%s is not a formal clan.\r\n"), name)
		return
	}
	send_to_char(ch, libc.CString("@cClan Name        @D: @C%s\n@cJoin Restriction @D: @C%s\n@cLeave Restriction@D: @C%s\n@D---@YClan Ranks@D---@n\n@cLeader@n\n@c%s@n\n@c%s@n\n@cMember@n\n\n%s@n\n"), S.Name, func() string {
		if S.Open_join == FALSE {
			return "Players must be enrolled to join this clan"
		}
		return "Players may join this clan as they please"
	}(), func() string {
		if S.Open_leave == FALSE {
			return "Players must be expelled to leave this clan"
		}
		return "Players may leave this clan as they please"
	}(), S.Highrank, S.Midrank, S.Info)
}
func listClansOfVictToChar(vict *char_data, ch *char_data) {
	var (
		i          int
		clan_found bool = FALSE != 0
	)
	if !IS_NPC(vict) {
		for i = 0; i < num_clans; i++ {
			if clanMemberFromList(int(vict.Idnum), (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Moderators) != nil || clanMemberFromList(int(vict.Idnum), (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Members) != nil {
				if int(libc.BoolToInt(clan_found)) == FALSE {
					clan_found = TRUE != 0
					send_to_char(ch, libc.CString("Clans %s belongs to:\r\n"), GET_NAME(vict))
				}
				send_to_char(ch, libc.CString("  %s\r\n"), (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name)
			}
		}
	}
	if !clan_found {
		send_to_char(ch, libc.CString("%s does not belong to any clans.\r\n"), GET_NAME(vict))
	}
}
func listClans(ch *char_data) {
	var i int
	if num_clans < 1 {
		send_to_char(ch, libc.CString("Presently, no clans have formally created.\r\n"))
		return
	}
	send_to_char(ch, libc.CString("The list of clans on Dragonball Advent Truth:\r\n"))
	for i = 0; i < num_clans; i++ {
		send_to_char(ch, libc.CString("  %s\r\n"), (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name)
	}
}
func checkCLAN(ch *char_data) int {
	var i int
	if num_clans < 1 {
		return FALSE
	}
	if ch.Clan == nil {
		return FALSE
	}
	for i = 0; i < num_clans; i++ {
		if libc.StrStr(ch.Clan, (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name) != nil {
			return TRUE
		}
	}
	return FALSE
}
func checkAPP(ch *char_data) {
	var i int
	if num_clans < 1 {
		return
	}
	for i = 0; i < num_clans; i++ {
		if libc.StrStr(ch.Clan, (*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name) != nil {
			ch.Clan = libc.StrDup((*(**clan_data)(unsafe.Add(unsafe.Pointer(clan), unsafe.Sizeof((*clan_data)(nil))*uintptr(i)))).Name)
			return
		}
	}
	return
}
