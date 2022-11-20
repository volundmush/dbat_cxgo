package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func strlcpy(dest *byte, source *byte, totalsize uint64) uint64 {
	libc.StrNCpy(dest, source, int(totalsize-1))
	*(*byte)(unsafe.Add(unsafe.Pointer(dest), totalsize-1)) = '\x00'
	return uint64(libc.StrLen(source))
}
