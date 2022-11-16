package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func write_aliases(ch *char_data) {
	var (
		file *stdio.File
		fn   [64936]byte
		temp *alias_data
	)
	get_filename(&fn[0], uint64(64936), ALIAS_FILE, GET_NAME(ch))
	stdio.Remove(libc.GoString(&fn[0]))
	if ch.Player_specials.Aliases == nil {
		return
	}
	if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&fn[0]), "w")
		return file
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Couldn't save aliases for %s in '%s': %s"), GET_NAME(ch), &fn[0], libc.StrError(libc.Errno))
		return
	}
	for temp = ch.Player_specials.Aliases; temp != nil; temp = temp.Next {
		var (
			aliaslen int = libc.StrLen(temp.Alias)
			repllen  int = libc.StrLen(temp.Replacement) - 1
		)
		stdio.Fprintf(file, "%d\n%s\n%d\n%s\n%d\n", aliaslen, temp.Alias, repllen, (*byte)(unsafe.Add(unsafe.Pointer(temp.Replacement), 1)), temp.Type)
	}
	file.Close()
}
func read_aliases(ch *char_data) {
	var (
		file   *stdio.File
		xbuf   [64936]byte
		t2     *alias_data
		prev   *alias_data = nil
		length int
	)
	get_filename(&xbuf[0], uint64(64936), ALIAS_FILE, GET_NAME(ch))
	if (func() *stdio.File {
		file = stdio.FOpen(libc.GoString(&xbuf[0]), "r")
		return file
	}()) == nil {
		if libc.Errno != 2 {
			basic_mud_log(libc.CString("SYSERR: Couldn't open alias file '%s' for %s: %s"), &xbuf[0], GET_NAME(ch), libc.StrError(libc.Errno))
		}
		return
	}
	ch.Player_specials.Aliases = new(alias_data)
	t2 = ch.Player_specials.Aliases
	for {
		if stdio.Fscanf(file, "%d\n", &length) != 1 {
			goto read_alias_error
		}
		file.GetS(&xbuf[0], int32(length+1))
		t2.Alias = libc.StrDup(&xbuf[0])
		if stdio.Fscanf(file, "%d\n", &length) != 1 {
			goto read_alias_error
		}
		xbuf[0] = ' '
		file.GetS(&xbuf[1], int32(length+1))
		t2.Replacement = libc.StrDup(&xbuf[0])
		if stdio.Fscanf(file, "%d\n", &length) != 1 {
			goto read_alias_error
		}
		t2.Type = length
		if int(file.IsEOF()) != 0 {
			break
		}
		t2.Next = new(alias_data)
		prev = t2
		t2 = t2.Next
	}
	file.Close()
	return
read_alias_error:
	if t2.Alias != nil {
		libc.Free(unsafe.Pointer(t2.Alias))
	}
	libc.Free(unsafe.Pointer(t2))
	if prev != nil {
		prev.Next = nil
	}
	file.Close()
}
func delete_aliases(charname *byte) {
	var filename [260]byte
	if get_filename(&filename[0], uint64(260), ALIAS_FILE, charname) == 0 {
		return
	}
	if stdio.Remove(libc.GoString(&filename[0])) < 0 && libc.Errno != 2 {
		basic_mud_log(libc.CString("SYSERR: deleting alias file %s: %s"), &filename[0], libc.StrError(libc.Errno))
	}
}
