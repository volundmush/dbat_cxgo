package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

const NUM_ZBUCKETS = 256
const MAX_ZDUMP_SIZE = 32

var beginPad [4]uint8 = [4]uint8{222, 173, 222, 173}
var endPad [4]uint8 = [4]uint8{222, 173, 222, 173}
var zfd *C.FILE = nil

type meminfo struct {
	Next  *meminfo
	Size  int
	Addr  *uint8
	Frees int
	File  *byte
	Line  int
}

var memlist [256]*meminfo
var zmalloclogging int = 2

func zmalloc_init() {
	var i int
	for i = 0; i < NUM_ZBUCKETS; i++ {
		memlist[i] = nil
	}
	zfd = (*C.FILE)(unsafe.Pointer(stdio.FOpen("zmalloc.log", "w+")))
}
func zdump(m *meminfo) {
	var (
		hextab  *uint8 = (*uint8)(unsafe.Pointer(libc.CString("0123456789ABCDEF")))
		hexline [37]uint8
		ascline [17]uint8
		hexp    *uint8
		ascp    *uint8
		inp     *uint8
		len_    int
		c       int = 1
	)
	if m.Addr == nil || m.Size <= 0 {
		return
	}
	hexp = &hexline[0]
	ascp = &ascline[0]
	inp = m.Addr
	if m.Size > MAX_ZDUMP_SIZE {
		len_ = MAX_ZDUMP_SIZE
	} else {
		len_ = m.Size
	}
	for ; len_ > 0; func() int {
		len_--
		inp = (*uint8)(unsafe.Add(unsafe.Pointer(inp), 1))
		return func() int {
			p := &c
			x := *p
			*p++
			return x
		}()
	}() {
		*(func() *uint8 {
			p := &hexp
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(*uint8)(unsafe.Add(unsafe.Pointer(hextab), (int(*inp)&240)>>4))
		*(func() *uint8 {
			p := &hexp
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(*uint8)(unsafe.Add(unsafe.Pointer(hextab), int(*inp)&15))
		if c%4 == 0 {
			*(func() *uint8 {
				p := &hexp
				x := *p
				*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = ' '
		}
		*(func() *uint8 {
			p := &ascp
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = uint8(int8(func() int {
			if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*inp))))) & int(uint16(int16(_ISprint)))) != 0 {
				return int(*inp)
			}
			return '.'
		}()))
		if c%16 == 0 || len_ <= 1 {
			*hexp = '\x00'
			*ascp = '\x00'
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "     %-40.40s%s\n", &hexline[0], &ascline[0])
			hexp = &hexline[0]
			ascp = &ascline[0]
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "\n")
}
func zmalloc(len_ int, file *byte, line int) *uint8 {
	var (
		ret *uint8
		m   *meminfo
	)
	ret = (*uint8)(libc.Calloc(1, len_+int(4)+int(4)))
	if ret == nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: malloc FAILED")
		return nil
	}
	libc.MemCpy(unsafe.Pointer(ret), unsafe.Pointer(&beginPad[0]), int(4))
	ret = (*uint8)(unsafe.Add(unsafe.Pointer(ret), 4))
	libc.MemCpy(unsafe.Add(unsafe.Pointer(ret), len_), unsafe.Pointer(&endPad[0]), int(4))
	if zmalloclogging > 2 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: 0x%4p  %d bytes %s:%d\n", ret, len_, file, line)
	}
	m = new(meminfo)
	if m == nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: FAILED mem alloc for zmalloc struct... bailing!\n")
		return nil
	}
	m.Addr = ret
	m.Size = len_
	m.Frees = 0
	m.File = C.strdup(file)
	if m.File == nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: FAILED mem alloc for zmalloc struct... bailing!\n")
		libc.Free(unsafe.Pointer(m))
		return nil
	}
	m.Line = line
	m.Next = memlist[(uint64(uintptr(unsafe.Pointer(ret)))>>3)&math.MaxUint8]
	memlist[(uint64(uintptr(unsafe.Pointer(ret)))>>3)&math.MaxUint8] = m
	return ret
}
func zrealloc(what *uint8, len_ int, file *byte, line int) *uint8 {
	var (
		ret    *uint8
		m      *meminfo
		prev_m *meminfo
	)
	if what != nil {
		for func() *meminfo {
			prev_m = nil
			return func() *meminfo {
				m = memlist[(uint64(uintptr(unsafe.Pointer(what)))>>3)&math.MaxUint8]
				return m
			}()
		}(); m != nil; func() *meminfo {
			prev_m = m
			return func() *meminfo {
				m = m.Next
				return m
			}()
		}() {
			if m.Addr == what {
				ret = (*uint8)(libc.Realloc(unsafe.Add(unsafe.Pointer(what), -int(4)), len_+int(4)+int(4)))
				if ret == nil {
					stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zrealloc: FAILED for 0x%p %d bytes mallocd at %s:%d,\n          %d bytes reallocd at %s:%d.\n", m.Addr, m.Size, m.File, m.Line, len_, file, line)
					if zmalloclogging > 1 {
						zdump(m)
					}
					return nil
				}
				libc.MemCpy(unsafe.Pointer(ret), unsafe.Pointer(&beginPad[0]), int(4))
				ret = (*uint8)(unsafe.Add(unsafe.Pointer(ret), 4))
				libc.MemCpy(unsafe.Add(unsafe.Pointer(ret), len_), unsafe.Pointer(&endPad[0]), int(4))
				if zmalloclogging > 2 {
					stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zrealloc: 0x%p %d bytes mallocd at %s:%d, %d bytes reallocd at %s:%d.\n", m.Addr, m.Size, m.File, m.Line, len_, file, line)
				}
				m.Addr = ret
				m.Size = len_
				if m.File != nil {
					libc.Free(unsafe.Pointer(m.File))
				}
				m.File = C.strdup(file)
				m.Line = line
				if prev_m != nil {
					prev_m.Next = m.Next
				} else {
					memlist[(uint64(uintptr(unsafe.Pointer(what)))>>3)&math.MaxUint8] = m.Next
				}
				m.Next = memlist[(uint64(uintptr(unsafe.Pointer(ret)))>>3)&math.MaxUint8]
				memlist[(uint64(uintptr(unsafe.Pointer(ret)))>>3)&math.MaxUint8] = m
				return ret
			}
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zrealloc: invalid pointer 0x%p, %d bytes to realloc at %s:%d.\n", what, len_, file, line)
	return zmalloc(len_, file, line)
}
func zfree2(what *uint8, file *byte, line int) {
	var (
		m     *meminfo
		gotit int = 0
	)
	if what == nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zfree2: ERR: Null pointer free'd: %s:%d.\n", file, line)
		return
	}
	for m = memlist[(uint64(uintptr(unsafe.Pointer(what)))>>3)&math.MaxUint8]; m != nil; m = m.Next {
		if m.Addr == what {
			if zmalloclogging > 2 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zfree2: Freed 0x%p %d bytes mallocd at %s:%d, freed at %s:%d\n", m.Addr, m.Size, m.File, m.Line, file, line)
			}
			pad_check(m)
			m.Frees++
			if m.Frees > 1 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zfree2: ERR: multiple frees! 0x%p %d bytes\n       mallocd at %s:%d, freed at %s:%d.\n", m.Addr, m.Size, m.File, m.Line, file, line)
				if zmalloclogging > 1 {
					zdump(m)
				}
			}
			gotit++
		}
	}
	if gotit == 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zfree2: ERR: attempt to free unallocated memory 0x%p at %s:%d.\n", what, file, line)
	}
	if gotit > 1 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zfree2: ERR: Multiply-allocd memory 0x%p.\n", what)
	}
}
func zC.strdup(src *byte, file *byte, line int) *byte {
	var result *byte
	result = (*byte)(unsafe.Pointer(zmalloc(int(C.strlen(src)+1), file, line)))
	if result == nil {
		return nil
	}
	C.strcpy(result, src)
	return result
}
func zmalloc_check() {
	var (
		m            *meminfo
		next_m       *meminfo
		admonishemnt *byte
		total_leak   int = 0
		num_leaks    int = 0
		i            int
	)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "\n------------ Checking leaks ------------\n\n")
	for i = 0; i < NUM_ZBUCKETS; i++ {
		for m = memlist[i]; m != nil; m = next_m {
			next_m = m.Next
			if m.Addr != nil && m.Frees <= 0 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: UNfreed memory 0x%p %d bytes mallocd at %s:%d\n", m.Addr, m.Size, m.File, m.Line)
				if zmalloclogging > 1 {
					zdump(m)
				}
				pad_check(m)
				total_leak += m.Size
				num_leaks++
			}
			if m.Addr != nil {
				libc.Free(unsafe.Add(unsafe.Pointer(m.Addr), -int(4)))
			}
			if m.File != nil {
				libc.Free(unsafe.Pointer(m.File))
			}
			libc.Free(unsafe.Pointer(m))
		}
	}
	if total_leak != 0 {
		if total_leak > 10000 {
			admonishemnt = libc.CString("you must work for Microsoft.")
		} else if total_leak > 5000 {
			admonishemnt = libc.CString("you should be ashamed!")
		} else if total_leak > 2000 {
			admonishemnt = libc.CString("you call yourself a programmer?")
		} else if total_leak > 1000 {
			admonishemnt = libc.CString("the X consortium has a job for you...")
		} else {
			admonishemnt = libc.CString("close, but not there yet.")
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: %d leaks totalling %d bytes... %s\n", num_leaks, total_leak, admonishemnt)
	} else {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "zmalloc: Congratulations: leak-free code!\n")
	}
	if zfd != nil {
		fflush(zfd)
		C.fclose(zfd)
	}
}
func pad_check(m *meminfo) {
	if memcmp(unsafe.Add(unsafe.Pointer(m.Addr), -int(4)), unsafe.Pointer(&beginPad[0]), uint64(4)) != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "pad_check: ERR: beginPad was modified! (mallocd@ %s:%d)\n", m.File, m.Line)
		if zmalloclogging > 1 {
			zdump(m)
		}
	}
	if memcmp(unsafe.Add(unsafe.Pointer(m.Addr), m.Size), unsafe.Pointer(&endPad[0]), uint64(4)) != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(zfd)), "pad_check: ERR: endPad was modified! (mallocd@ %s:%d)\n", m.File, m.Line)
		if zmalloclogging > 1 {
			zdump(m)
		}
	}
}
