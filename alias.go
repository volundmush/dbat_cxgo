package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func write_aliases(ch *char_data) {
	var (
		file *C.FILE
		fn   [64936]byte
		temp *alias_data
	)
	get_filename(&fn[0], uint64(64936), ALIAS_FILE, GET_NAME(ch))
	stdio.Remove(libc.GoString(&fn[0]))
	if ch.Player_specials.Aliases == nil {
		return
	}
	if (func() *C.FILE {
		file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fn[0]), "w")))
		return file
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Couldn't save aliases for %s in '%s': %s"), GET_NAME(ch), &fn[0], C.strerror(*__errno_location()))
		return
	}
	for temp = ch.Player_specials.Aliases; temp != nil; temp = temp.Next {
		var (
			aliaslen int = int(C.strlen(temp.Alias))
			repllen  int = int(C.strlen(temp.Replacement) - 1)
		)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(file)), "%d\n%s\n%d\n%s\n%d\n", aliaslen, temp.Alias, repllen, (*byte)(unsafe.Add(unsafe.Pointer(temp.Replacement), 1)), temp.Type)
	}
	C.fclose(file)
}
func read_aliases(ch *char_data) {
	var (
		file   *C.FILE
		xbuf   [64936]byte
		t2     *alias_data
		prev   *alias_data = nil
		length int
	)
	get_filename(&xbuf[0], uint64(64936), ALIAS_FILE, GET_NAME(ch))
	if (func() *C.FILE {
		file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&xbuf[0]), "r")))
		return file
	}()) == nil {
		if (*__errno_location()) != ENOENT {
			basic_mud_log(libc.CString("SYSERR: Couldn't open alias file '%s' for %s: %s"), &xbuf[0], GET_NAME(ch), C.strerror(*__errno_location()))
		}
		return
	}
	ch.Player_specials.Aliases = new(alias_data)
	t2 = ch.Player_specials.Aliases
	for {
		if __isoc99_fscanf(file, libc.CString("%d\n"), &length) != 1 {
			goto read_alias_error
		}
		C.fgets(&xbuf[0], length+1, file)
		t2.Alias = C.strdup(&xbuf[0])
		if __isoc99_fscanf(file, libc.CString("%d\n"), &length) != 1 {
			goto read_alias_error
		}
		xbuf[0] = ' '
		C.fgets(&xbuf[1], length+1, file)
		t2.Replacement = C.strdup(&xbuf[0])
		if __isoc99_fscanf(file, libc.CString("%d\n"), &length) != 1 {
			goto read_alias_error
		}
		t2.Type = length
		if C.feof(file) != 0 {
			break
		}
		t2.Next = new(alias_data)
		prev = t2
		t2 = t2.Next
	}
	C.fclose(file)
	return
read_alias_error:
	if t2.Alias != nil {
		libc.Free(unsafe.Pointer(t2.Alias))
	}
	libc.Free(unsafe.Pointer(t2))
	if prev != nil {
		prev.Next = nil
	}
	C.fclose(file)
}
func delete_aliases(charname *byte) {
	var filename [4096]byte
	if get_filename(&filename[0], uint64(4096), ALIAS_FILE, charname) == 0 {
		return
	}
	if stdio.Remove(libc.GoString(&filename[0])) < 0 && (*__errno_location()) != ENOENT {
		basic_mud_log(libc.CString("SYSERR: deleting alias file %s: %s"), &filename[0], C.strerror(*__errno_location()))
	}
}
