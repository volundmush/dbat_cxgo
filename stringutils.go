package main

import "unsafe"

func strlcpy(dest *byte, source *byte, totalsize uint64) uint64 {
	C.strncpy(dest, source, totalsize-1)
	*(*byte)(unsafe.Add(unsafe.Pointer(dest), totalsize-1)) = '\x00'
	return uint64(C.strlen(source))
}
