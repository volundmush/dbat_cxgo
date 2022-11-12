package main

import "C"
import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const TOG_OFF = 0
const TOG_ON = 1

func log_imm_action(messg *byte, _rest ...interface{}) {
	var (
		fl       *C.FILE
		filename *byte
		fbuf     stat
	)
	filename = libc.CString(LIB_MISC)
	if C.stat(filename, &fbuf) < 0 {
		C.perror(libc.CString("SYSERR: Can't C.stat() file"))
		return
	}
	if fbuf.St_size >= __off_t(config_info.Operation.Max_filesize*4) {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(filename), "a")))
		return fl
	}()) == nil {
		C.perror(libc.CString("SYSERR: log_imm_action"))
		return
	}
	var ct int64 = C.time(nil)
	var time_s *byte = C.asctime(C.localtime(&ct))
	var args libc.ArgList
	args.Start(messg, _rest)
	*(*byte)(unsafe.Add(unsafe.Pointer(time_s), C.strlen(time_s)-1)) = '\x00'
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%-15.15s :: ", (*byte)(unsafe.Add(unsafe.Pointer(time_s), 4)))
	stdio.Vfprintf((*stdio.File)(unsafe.Pointer(fl)), libc.GoString(messg), args)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "\n")
	args.End()
	C.fclose(fl)
}
func log_custom(d *descriptor_data, obj *obj_data) {
	var (
		fl       *C.FILE
		filename *byte
		fbuf     stat
	)
	filename = libc.CString(LIB_MISC)
	if C.stat(filename, &fbuf) < 0 {
		C.perror(libc.CString("SYSERR: Can't C.stat() file"))
		return
	}
	if fbuf.St_size >= __off_t(config_info.Operation.Max_filesize*4) {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(filename), "a")))
		return fl
	}()) == nil {
		C.perror(libc.CString("SYSERR: log_custom"))
		return
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "@D[@cUser@W: @R%-20s @cName@W: @C%-20s @cCustom@W: @Y%s@D]\n", GET_USER(d.Character), GET_NAME(d.Character), obj.Short_description)
	C.fclose(fl)
}
func bring_to_cap(ch *char_data) {
	var (
		skippl int = FALSE
		skipki int = FALSE
		skipst int = FALSE
		mult   int = 1
	)
	if !soft_cap(ch, 0) {
		skippl = TRUE
	}
	if !soft_cap(ch, 1) {
		skipki = TRUE
	}
	if !soft_cap(ch, 2) {
		skipst = TRUE
	}
	if ch.Race == RACE_BIO {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = int(3.5)
		} else if PLR_FLAGGED(ch, PLR_TRANS4) {
			mult = 4
		}
	} else if ch.Race == RACE_MAJIN {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 2
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = int(4.5)
		}
	} else if ch.Race == RACE_TRUFFLE {
		if PLR_FLAGGED(ch, PLR_TRANS1) {
			mult = 3
		} else if PLR_FLAGGED(ch, PLR_TRANS2) {
			mult = 4
		} else if PLR_FLAGGED(ch, PLR_TRANS3) {
			mult = 5
		}
	}
	if ch.Race == RACE_KANASSAN || ch.Race == RACE_DEMON {
		if skippl == FALSE {
			var (
				base int64 = ch.Basepl
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 2000000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 2000000)
			}
			if base < int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 1000000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 1000000)
			}
			if base < int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 300000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 300000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 100000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 100000)
			}
			if base < int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 40000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 40000)
			}
			if base < int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 25000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 25000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 500)
			}
		}
		if skipki == FALSE {
			var (
				base int64 = ch.Baseki
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 2000000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 2000000)
			}
			if base < int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 1000000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 1000000)
			}
			if base < int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 300000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 300000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 100000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 100000)
			}
			if base < int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 40000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 40000)
			}
			if base < int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 25000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 25000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 500)
			}
		}
		if skipst == FALSE {
			var (
				base int64 = ch.Basest
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*2000000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 2000000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 2000000)
			}
			if base < int64(GET_LEVEL(ch)*1000000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 1000000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 1000000)
			}
			if base < int64(GET_LEVEL(ch)*300000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 300000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 300000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*100000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 100000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 100000)
			}
			if base < int64(GET_LEVEL(ch)*40000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 40000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 40000)
			}
			if base < int64(GET_LEVEL(ch)*25000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 25000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 25000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 500)
			}
		}
	} else {
		if skippl == FALSE {
			var (
				base int64 = ch.Basepl
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 1500000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 1500000)
			}
			if base < int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 800000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 800000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 200000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 200000)
			}
			if base < int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 80000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 80000)
			}
			if base < int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 20000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 20000)
			}
			if base < int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 15000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 15000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Basepl))
				ch.Max_hit += diff * int64(mult)
				ch.Basepl = int64(GET_LEVEL(ch) * 500)
			}
		}
		if skipki == FALSE {
			var (
				base int64 = ch.Baseki
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 1500000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 1500000)
			}
			if base < int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 800000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 800000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 200000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 200000)
			}
			if base < int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 80000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 80000)
			}
			if base < int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 20000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 20000)
			}
			if base < int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 15000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 15000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Baseki))
				ch.Max_mana += diff * int64(mult)
				ch.Baseki = int64(GET_LEVEL(ch) * 500)
			}
		}
		if skipst == FALSE {
			var (
				base int64 = ch.Basest
				diff int64
			)
			if base < int64(GET_LEVEL(ch)*1500000) && GET_LEVEL(ch) > 90 && GET_LEVEL(ch) <= 99 {
				diff = int64((GET_LEVEL(ch) * 1500000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 1500000)
			}
			if base < int64(GET_LEVEL(ch)*800000) && GET_LEVEL(ch) > 80 && GET_LEVEL(ch) <= 90 {
				diff = int64((GET_LEVEL(ch) * 800000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 800000)
			}
			if base < int64(GET_LEVEL(ch)*250000) && GET_LEVEL(ch) > 70 && GET_LEVEL(ch) <= 80 {
				diff = int64((GET_LEVEL(ch) * 250000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 250000)
			}
			if base < int64(GET_LEVEL(ch)*200000) && GET_LEVEL(ch) > 60 && GET_LEVEL(ch) <= 70 {
				diff = int64((GET_LEVEL(ch) * 200000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 200000)
			}
			if base < int64(GET_LEVEL(ch)*80000) && GET_LEVEL(ch) > 50 && GET_LEVEL(ch) <= 60 {
				diff = int64((GET_LEVEL(ch) * 80000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 80000)
			}
			if base < int64(GET_LEVEL(ch)*20000) && GET_LEVEL(ch) > 40 && GET_LEVEL(ch) <= 50 {
				diff = int64((GET_LEVEL(ch) * 20000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 20000)
			}
			if base < int64(GET_LEVEL(ch)*15000) && GET_LEVEL(ch) > 30 && GET_LEVEL(ch) <= 40 {
				diff = int64((GET_LEVEL(ch) * 15000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 15000)
			}
			if base < int64(GET_LEVEL(ch)*5000) && GET_LEVEL(ch) > 20 && GET_LEVEL(ch) <= 30 {
				diff = int64((GET_LEVEL(ch) * 5000) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 5000)
			}
			if base < int64(GET_LEVEL(ch)*1500) && GET_LEVEL(ch) > 10 && GET_LEVEL(ch) <= 20 {
				diff = int64((GET_LEVEL(ch) * 1500) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 1500)
			}
			if base < int64(GET_LEVEL(ch)*500) && GET_LEVEL(ch) <= 10 {
				diff = int64((GET_LEVEL(ch) * 500) - int(ch.Basest))
				ch.Max_move += diff * int64(mult)
				ch.Basest = int64(GET_LEVEL(ch) * 500)
			}
		}
	}
}
func do_rpp(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg        [2048]byte
		arg2       [2048]byte
		tnlcost    int = 1
		revcost    int = 1
		selection  int = 0
		pay        int = 0
		bpay       int = 0
		max_choice int = 15
		obj        *obj_data
	)
	half_chop(argument, &arg[0], &arg2[0])
	if IS_NPC(ch) {
		return
	}
	revcost = revcost * (GET_LEVEL(ch) / 15)
	tnlcost = (tnlcost * (GET_LEVEL(ch) / 40)) + 1
	if revcost < 1 {
		revcost = 1
	}
	if tnlcost < 1 {
		tnlcost = 1
	}
	if PLR_FLAGGED(ch, PLR_PDEATH) {
		revcost *= 6
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@C                             Rewards Menu\n"))
		send_to_char(ch, libc.CString("@b  ------------------------------------------------------------------\n"))
		send_to_char(ch, libc.CString("  @C1@D)@R Disabled            @D[@G -- RPP @D]  @C2@D)@R Disabled              @D[@G -- RPP @D]@n\n"))
		send_to_char(ch, libc.CString("  @C3@D)@c Custom Equipment    @D[@G 30 RPP @D]  @C4@D)@c Alignment Change      @D[@G 20 RPP @D]\n"))
		send_to_char(ch, libc.CString("  @C5@D)@c 7,500 zenni         @D[@G  1 RPP @D]  @C6@D)@c +2 To A Stat          @D[@G  2 RPP @D]\n"))
		send_to_char(ch, libc.CString("  @C7@D)@c +750 PS             @D[@G  4 RPP @D]  @C8@D)@c Revival               @D[@G%3d RPP @D]\n"), revcost)
		send_to_char(ch, libc.CString("  @C9@D)@c 50%s TNL Exp         @D[@G%3d RPP @D] @C10@D)@c Aura Change           @D[@G  2 RPP @D]\n"), "%", tnlcost)
		send_to_char(ch, libc.CString(" @C11@D)@c Reach Softcap       @D[@G  6 RPP @D] @C12@D)@c RPP Store             @D[@G ?? RPP @D]\n"))
		send_to_char(ch, libc.CString(" @C13@D)@c Extra Feature       @D[@G  1 RPP @D] @C14@D)@c Restring Equipment    @D[@G  1 RPP @D]\n"))
		send_to_char(ch, libc.CString(" @C15@D)@c Extra Skillslot     @D[@G  5 RPP @D] @C16@D)@R Disabled              @D[@G -- RPP @D]@n\n"))
		send_to_char(ch, libc.CString("@b  ------------------------------------------------------------------@n\n"))
		send_to_char(ch, libc.CString("@D                           [@YYour RPP@D:@G %3d@D]@n\n"), ch.Rp)
		send_to_char(ch, libc.CString("\nSyntax: rpp (num)\n"))
		return
	}
	selection = libc.Atoi(libc.GoString(&arg[0]))
	if selection <= 0 || selection > max_choice {
		send_to_char(ch, libc.CString("You must choose a number from the menu. Enter the command again with no arguments for the menu.\r\n"))
		return
	}
	if selection > 2 {
		if selection == 3 {
			if ch.Rp < 30 {
				send_to_char(ch, libc.CString("You need at least 30 RPP to initiate a custom equipment build.\r\n"))
				return
			} else {
				ch.Desc.Connected = CON_POBJ
				ch.Desc.Obj_name = C.strdup(libc.CString("Generic Armor Vest"))
				ch.Desc.Obj_short = C.strdup(libc.CString("@cGeneric @DArmor @WVest@n"))
				ch.Desc.Obj_long = C.strdup(libc.CString("@wA @cgeneric @Darmor @Wvest@w is lying here@n"))
				ch.Desc.Obj_type = 1
				ch.Desc.Obj_weapon = 0
				disp_custom_menu(ch.Desc)
				ch.Desc.Obj_editflag = EDIT_CUSTOM
				ch.Desc.Obj_editval = EDIT_CUSTOM_MAIN
				return
			}
		}
		if selection == 4 {
			pay = 20
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else {
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("What do you want to change your alignment to? (evil, sorta-evil, neutral, sorta-good, good)"))
					return
				}
				if C.strcasecmp(&arg2[0], libc.CString("evil")) == 0 {
					send_to_char(ch, libc.CString("You change your alignment to Evil.\r\n"))
					ch.Alignment = -750
				} else if C.strcasecmp(&arg2[0], libc.CString("sorta-evil")) == 0 {
					send_to_char(ch, libc.CString("You change your alignment to Sorta Evil.\r\n"))
					ch.Alignment = -50
				} else if C.strcasecmp(&arg2[0], libc.CString("neutral")) == 0 {
					send_to_char(ch, libc.CString("You change your alignment to Neutral.\r\n"))
					ch.Alignment = 0
				} else if C.strcasecmp(&arg2[0], libc.CString("sorta-good")) == 0 {
					send_to_char(ch, libc.CString("You change your alignment to Sorta Good.\r\n"))
					ch.Alignment = 51
				} else if C.strcasecmp(&arg2[0], libc.CString("good")) == 0 {
					send_to_char(ch, libc.CString("You change your alignment to Good.\r\n"))
					ch.Alignment = 300
				} else {
					send_to_char(ch, libc.CString("That is not an acceptable option for changing alignment.\r\n"))
					return
				}
			}
		}
		if selection == 5 {
			pay = 1
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else {
				ch.Bank_gold += 7500
				send_to_char(ch, libc.CString("Your bank zenni has been increased by 7,500\r\n"))
			}
		}
		if selection == 6 {
			pay = 2
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else {
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("What stat? (str, con, int, wis, spd, agl)"))
					return
				}
				if C.strcasecmp(&arg2[0], libc.CString("str")) == 0 {
					if (ch.Bonuses[BONUS_WIMP]) > 0 && ch.Real_abils.Str >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Str >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your strength by 2.\r\n"))
					ch.Real_abils.Str += 2
				} else if C.strcasecmp(&arg2[0], libc.CString("con")) == 0 {
					if (ch.Bonuses[BONUS_FRAIL]) > 0 && ch.Real_abils.Con >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Con >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your constitution by 2.\r\n"))
					ch.Real_abils.Con += 2
				} else if C.strcasecmp(&arg2[0], libc.CString("int")) == 0 {
					if (ch.Bonuses[BONUS_DULL]) > 0 && ch.Real_abils.Intel >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Intel >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your intelligence by 2.\r\n"))
					ch.Real_abils.Intel += 2
				} else if C.strcasecmp(&arg2[0], libc.CString("wis")) == 0 {
					if (ch.Bonuses[BONUS_FOOLISH]) > 0 && ch.Real_abils.Wis >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Wis >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your wisdom by 2.\r\n"))
					ch.Real_abils.Wis += 2
				} else if C.strcasecmp(&arg2[0], libc.CString("spd")) == 0 {
					if (ch.Bonuses[BONUS_SLOW]) > 0 && ch.Real_abils.Cha >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Cha >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your speed by 2.\r\n"))
					ch.Real_abils.Cha += 2
				} else if C.strcasecmp(&arg2[0], libc.CString("agl")) == 0 {
					if (ch.Bonuses[BONUS_CLUMSY]) > 0 && ch.Real_abils.Dex >= 45 {
						send_to_char(ch, libc.CString("You can't because that stat maxes at 45 due to a trait negative.\r\n"))
						return
					} else if ch.Aff_abils.Dex >= 100 {
						send_to_char(ch, libc.CString("100 is the maximum for any stat.\r\n"))
						return
					}
					send_to_char(ch, libc.CString("You increase your speed by 2.\r\n"))
					ch.Real_abils.Dex += 2
				} else {
					send_to_char(ch, libc.CString("That is not an acceptable option for changing alignment.\r\n"))
					return
				}
			}
		}
		if selection == 7 {
			pay = 4
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else {
				ch.Player_specials.Class_skill_points[ch.Chclass] += 750
				send_to_char(ch, libc.CString("Your practices have been increased by 750\r\n"))
			}
		}
		if selection == 8 {
			pay = revcost
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else if !AFF_FLAGGED(ch, AFF_SPIRIT) {
				send_to_char(ch, libc.CString("You aren't even dead!"))
				return
			} else {
				ch.Affected_by[int(AFF_ETHEREAL/32)] &= ^(1 << (int(AFF_ETHEREAL % 32)))
				ch.Affected_by[int(AFF_SPIRIT/32)] &= ^(1 << (int(AFF_SPIRIT % 32)))
				ch.Hit = gear_pl(ch)
				ch.Mana = ch.Max_mana
				ch.Move = ch.Max_move
				ch.Limb_condition[0] = 100
				ch.Limb_condition[1] = 100
				ch.Limb_condition[2] = 100
				ch.Limb_condition[3] = 100
				ch.Act[int(PLR_HEAD/32)] |= bitvector_t(1 << (int(PLR_HEAD % 32)))
				ch.Act[int(PLR_PDEATH/32)] &= bitvector_t(^(1 << (int(PLR_PDEATH % 32))))
				char_from_room(ch)
				if ch.Droom != room_vnum(-1) && ch.Droom != 0 && ch.Droom != 1 {
					char_to_room(ch, real_room(ch.Droom))
				} else if ch.Chclass == CLASS_ROSHI {
					char_to_room(ch, real_room(1130))
				} else if ch.Chclass == CLASS_KABITO {
					char_to_room(ch, real_room(0x2F42))
				} else if ch.Chclass == CLASS_NAIL {
					char_to_room(ch, real_room(0x2DA3))
				} else if ch.Chclass == CLASS_BARDOCK {
					char_to_room(ch, real_room(2268))
				} else if ch.Chclass == CLASS_KRANE {
					char_to_room(ch, real_room(0x32D1))
				} else if ch.Chclass == CLASS_TAPION {
					char_to_room(ch, real_room(8231))
				} else if ch.Chclass == CLASS_PICCOLO {
					char_to_room(ch, real_room(1659))
				} else if ch.Chclass == CLASS_ANDSIX {
					char_to_room(ch, real_room(1713))
				} else if ch.Chclass == CLASS_DABURA {
					char_to_room(ch, real_room(6486))
				} else if ch.Chclass == CLASS_FRIEZA {
					char_to_room(ch, real_room(4282))
				} else if ch.Chclass == CLASS_GINYU {
					char_to_room(ch, real_room(4289))
				} else if ch.Chclass == CLASS_JINTO {
					char_to_room(ch, real_room(3499))
				} else if ch.Chclass == CLASS_TSUNA {
					char_to_room(ch, real_room(15000))
				} else if ch.Chclass == CLASS_KURZAK {
					char_to_room(ch, real_room(16100))
				} else {
					char_to_room(ch, real_room(300))
					send_to_imm(libc.CString("ERROR: Player %s without acceptable sensei.\r\n"), GET_NAME(ch))
				}
				look_at_room(ch.In_room, ch, 0)
				ch.Deathtime = 0
				act(libc.CString("$n's body forms in a pool of @Bblue light@n."), TRUE, ch, nil, nil, TO_ROOM)
				send_to_char(ch, libc.CString("You have been revived.\r\n"))
				if GET_LEVEL(ch) > 9 {
					var losschance int = axion_dice(0)
					if GET_LEVEL(ch) > 9 {
						send_to_char(ch, libc.CString("@RThe the strain of this type of revival has caused you to be in a weakened state for 100 hours (Game time)! Strength, constitution, wisdom, intelligence, speed, and agility have been reduced by 8 points for the duration.@n\r\n"))
						var str int = -8
						var intel int = -8
						var wis int = -8
						var spd int = -8
						var con int = -8
						var agl int = -8
						if ch.Real_abils.Str <= 16 {
							str = -4
						}
						if ch.Real_abils.Intel <= 16 {
							intel = -4
						}
						if ch.Real_abils.Cha <= 16 {
							spd = -4
						}
						if ch.Real_abils.Dex <= 16 {
							agl = -4
						}
						if ch.Real_abils.Wis <= 16 {
							wis = -4
						}
						if ch.Real_abils.Con <= 16 {
							con = -4
						}
						assign_affect(ch, AFF_WEAKENED_STATE, SKILL_WARP, 100, str, con, intel, agl, wis, spd)
					}
					if losschance >= 100 {
						var psloss int = rand_number(100, 300)
						ch.Player_specials.Class_skill_points[ch.Chclass] -= psloss
						send_to_char(ch, libc.CString("@R...and a loss of @r%d@R PS!@n"), psloss)
						if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 0 {
							ch.Player_specials.Class_skill_points[ch.Chclass] = 0
						}
					}
				}
			}
		}
		if selection == 9 {
			pay = tnlcost
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else if GET_LEVEL(ch) >= 100 {
				send_to_char(ch, libc.CString("You can not buy experience anymore at your level. I think you know why.\r\n"))
				return
			} else if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) < 0 {
				send_to_char(ch, libc.CString("You can not buy experience anymore UNTIL you level.\r\n"))
				return
			} else {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.52)
				send_to_char(ch, libc.CString("You gained 50%s of the entire experience needed for your next level.\r\n"), "%")
			}
		}
		if selection == 10 {
			pay = 2
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else {
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("Change your aura to what? (white, blue, red, green, pink, purple, yellow, black, orange)"))
					return
				}
				if C.strcasecmp(&arg2[0], libc.CString("white")) == 0 {
					ch.Aura = 0
					send_to_char(ch, libc.CString("You change your aura to white.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("blue")) == 0 {
					ch.Aura = 1
					send_to_char(ch, libc.CString("You change your aura to blue.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("red")) == 0 {
					ch.Aura = 2
					send_to_char(ch, libc.CString("You change your aura to red.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("green")) == 0 {
					ch.Aura = 3
					send_to_char(ch, libc.CString("You change your aura to green.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("pink")) == 0 {
					ch.Aura = 4
					send_to_char(ch, libc.CString("You change your aura to pink.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("purple")) == 0 {
					ch.Aura = 5
					send_to_char(ch, libc.CString("You change your aura to purple.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("yellow")) == 0 {
					ch.Aura = 6
					send_to_char(ch, libc.CString("You change your aura to yellow.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("black")) == 0 {
					ch.Aura = 7
					send_to_char(ch, libc.CString("You change your aura to black.\r\n"))
				} else if C.strcasecmp(&arg2[0], libc.CString("orange")) == 0 {
					ch.Aura = 8
					send_to_char(ch, libc.CString("You change your aura to orange.\r\n"))
				} else {
					send_to_char(ch, libc.CString("That is not an acceptable option for changing alignment.\r\n"))
					return
				}
			}
		}
		if selection == 11 {
			pay = 6
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP for that selection.\r\n"))
				return
			} else if GET_LEVEL(ch) >= 100 {
				send_to_char(ch, libc.CString("You can't use this at level 100.\r\n"))
				return
			} else if ch.Race == RACE_ARLIAN {
				send_to_char(ch, libc.CString("This is not available to bugs.\r\n"))
				return
			} else {
				if !soft_cap(ch, 0) && !soft_cap(ch, 1) && !soft_cap(ch, 2) {
					send_to_char(ch, libc.CString("You are already above your softcap for this level.\r\n"))
					return
				}
				bring_to_cap(ch)
			}
		}
		if selection == 12 {
			if arg2[0] == 0 {
				disp_rpp_store(ch)
				return
			} else if libc.Atoi(libc.GoString(&arg2[0])) <= 0 {
				send_to_char(ch, libc.CString("That is not a choice in the RPP store!\r\n"))
				return
			} else {
				var choice int = libc.Atoi(libc.GoString(&arg2[0]))
				handle_rpp_store(ch, choice)
				return
			}
		}
		if selection == 13 {
			rpp_feature(ch, &arg2[0])
			return
		}
		if selection == 14 {
			pay = 1
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You need at least 1 RPP to initiate an equipment restring.\r\n"))
				return
			} else if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &arg2[0], nil, ch.Carrying)
				return obj
			}()) == nil {
				send_to_char(ch, libc.CString("You don't have a that equipment to restring in your inventory.\r\n"))
				send_to_char(ch, libc.CString("Syntax: rpp 14 (obj name)\r\n"))
				return
			} else if OBJ_FLAGGED(obj, ITEM_CUSTOM) {
				send_to_char(ch, libc.CString("You can not restring a custom piece. Why? Cause I say so. :P\r\n"))
				return
			} else {
				ch.Desc.Connected = CON_POBJ
				var thename [2048]byte
				var theshort [2048]byte
				var thelong [2048]byte
				thename[0] = '\x00'
				theshort[0] = '\x00'
				thelong[0] = '\x00'
				stdio.Sprintf(&thename[0], "%s", obj.Name)
				stdio.Sprintf(&theshort[0], "%s", obj.Short_description)
				stdio.Sprintf(&thelong[0], "%s", obj.Description)
				ch.Desc.Obj_name = C.strdup(&thename[0])
				ch.Desc.Obj_was = C.strdup(&theshort[0])
				ch.Desc.Obj_short = C.strdup(&theshort[0])
				ch.Desc.Obj_long = C.strdup(&thelong[0])
				ch.Desc.Obj_point = obj
				ch.Desc.Obj_type = 1
				ch.Desc.Obj_weapon = 0
				disp_restring_menu(ch.Desc)
				ch.Desc.Obj_editflag = EDIT_RESTRING
				ch.Desc.Obj_editval = EDIT_RESTRING_MAIN
				return
			}
		}
		if selection == 15 {
			pay = 5
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP in your bank for that selection.\r\n"))
				return
			} else if (ch.Bonuses[BONUS_GMEMORY]) != 0 && ch.Skill_slots >= 65 {
				send_to_char(ch, libc.CString("You are already at your skillslot cap.\r\n"))
				return
			} else if (ch.Bonuses[BONUS_GMEMORY]) == 0 && ch.Skill_slots >= 60 {
				send_to_char(ch, libc.CString("You are already at your skillslot cap.\r\n"))
				return
			} else {
				ch.Skill_slots += 1
			}
		}
		if selection == 16 {
			pay = 5000
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("You do not have enough RPP in your bank for that selection.\r\n"))
				return
			} else {
				var (
					found int       = FALSE
					k     *obj_data = nil
				)
				for k = object_list; k != nil; k = k.Next {
					if OBJ_FLAGGED(k, ITEM_FORGED) {
						continue
					}
					if GET_OBJ_VNUM(k) == 20 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 21 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 22 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 23 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 24 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 25 {
						found = TRUE
					} else if GET_OBJ_VNUM(k) == 26 {
						found = TRUE
					}
				}
				if found == FALSE {
					send_to_char(ch, libc.CString("You have reduced the Dragon Ball wait by a whole real life day!\r\n"))
					send_to_all(libc.CString("%s has just reduced the Dragon Ball wait by a whole real life day!\r\n"), GET_NAME(ch))
					dballtime -= 86400
					if dballtime <= 0 {
						dballtime = 1
					}
				} else if SELFISHMETER >= 10 {
					send_to_char(ch, libc.CString("Sorry, it seems there there are several powers interfering with the Dragon Balls.\r\n"))
					return
				} else {
					send_to_char(ch, libc.CString("Sorry, but there is already a set of Dragon Balls in existence.\r\n"))
					return
				}
			}
		}
	}
	if selection <= 2 {
		var (
			fl       *C.FILE
			filename *byte
			fbuf     stat
		)
		filename = libc.CString(LIB_MISC)
		if selection == 1 {
			pay = 6500
			if ch.Rp < pay {
				send_to_char(ch, libc.CString("Nice try but you don't have enough RPP for that.\r\n"))
				return
			} else {
				send_to_char(ch, libc.CString("You now have an Excel House Capsule!\r\n"))
				var hobj *obj_data = read_object(6, VIRTUAL)
				obj_to_char(hobj, ch)
				ch.Rp -= pay
				ch.Desc.Rpp = ch.Rp
				userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
				save_char(ch)
				send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), pay)
				send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), pay)
				return
			}
		} else if selection == 2 {
			return
			pay = 200
		}
		if ch.Rp < pay {
			send_to_char(ch, libc.CString("Nice try but you don't have enough RPP for that.\r\n"))
			return
		}
		if C.stat(filename, &fbuf) < 0 {
			C.perror(libc.CString("SYSERR: Can't C.stat() file"))
			return
		}
		if fbuf.St_size >= __off_t(config_info.Operation.Max_filesize) {
			send_to_char(ch, libc.CString("Sorry, the file is full right now.. try again later.\r\n"))
			return
		}
		if (func() *C.FILE {
			fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(filename), "a")))
			return fl
		}()) == nil {
			C.perror(libc.CString("SYSERR: do_reward_request"))
			send_to_char(ch, libc.CString("Could not open the file.  Sorry.\r\n"))
			return
		}
		if selection == 1 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "@D[@cName@W: @C%-20s @cRequest@W: @Y%-20s@D]\n", GET_NAME(ch), "House")
			send_to_imm(libc.CString("RPP Request: %s paid for house"), GET_NAME(ch))
			BOARDNEWCOD = C.time(nil)
			save_mud_time(&time_info)
		} else if selection == 2 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "@D[@cName@W: @C%-20s @cRequest@W: @Y%-20s@D]\n", GET_NAME(ch), "Custom Skill")
			send_to_imm(libc.CString("RPP Request: %s paid for Custom Skill, uhoh spaggettios"), GET_NAME(ch))
			BOARDNEWCOD = C.time(nil)
			save_mud_time(&time_info)
		}
		ch.Rp -= pay
		ch.Desc.Rpp = ch.Rp
		userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		save_char(ch)
		send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. An immortal will address the request soon enough. Be patient.@n\r\n"), pay)
		C.fclose(fl)
	}
	if selection >= 4 && selection < 12 && pay > 0 {
		ch.Rp -= pay
		ch.Desc.Rpp = ch.Rp
		userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		save_char(ch)
		send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), pay)
		send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), pay)
	}
	if selection > 12 && pay > 0 {
		ch.Rp -= pay
		ch.Desc.Rpp = ch.Rp
		userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		save_char(ch)
		send_to_char(ch, libc.CString("@R%d@W RPP paid for your selection. Enjoy!@n\r\n"), bpay)
		send_to_imm(libc.CString("RPP Purchase: %s %d"), GET_NAME(ch), bpay)
	}
}
func do_commune(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_COMMUNE) == 0 {
		return
	}
	if ch.Move >= ch.Max_move {
		send_to_char(ch, libc.CString("Your stamina is already at full.\r\n"))
		return
	}
	var prob int = GET_SKILL(ch, SKILL_COMMUNE)
	var perc int = axion_dice(0)
	var cost int64 = int64(float64(ch.Max_move) * 0.05)
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki to commune with the Eldritch Star.\r\n"))
		return
	}
	if prob < perc {
		ch.Mana -= cost
		reveal_hiding(ch, 0)
		act(libc.CString("@cYou close your eyes and try to commune with the Eldritch Star. You are unable to concentrate though.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n closes $s eyes for a moment. Then $e reopens them and frowns.@n"), TRUE, ch, nil, nil, TO_ROOM)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		ch.Mana -= cost
		ch.Move += cost
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		reveal_hiding(ch, 0)
		act(libc.CString("@cYou close your eyes and commune with the Eldritch Star spiritually. You feel your stamina replenish some.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n closes $s eyes for a moment. Then $e reopens them and smiles.@n"), TRUE, ch, nil, nil, TO_ROOM)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
}
func do_willpower(ch *char_data, argument *byte, cmd int, subcmd int) {
	var fail int = FALSE
	if IS_NPC(ch) {
		return
	}
	if ch.Majinize <= 0 {
		send_to_char(ch, libc.CString("You are not majinized and have no need to reclaim full control of your own will.\r\n"))
		return
	} else {
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 100 && GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("You do not have enough PS to focus your attempt to break free.\r\n"))
			fail = TRUE
		}
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 200 && GET_LEVEL(ch) >= 100 {
			send_to_char(ch, libc.CString("You do not have enough PS to focus your attempt to break free.\r\n"))
			fail = TRUE
		}
		if ch.Exp < int64(level_exp(ch, GET_LEVEL(ch)+1)) && GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("You need a full level's worth of experience stored up to try and break free.\r\n"))
			fail = TRUE
		}
		if fail == TRUE {
			return
		} else {
			ch.Exp = 0
			ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
			if GET_LEVEL(ch) >= 100 {
				ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
			}
			if rand_number(10, 100)-int(ch.Aff_abils.Intel) > 60 {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou focus all your knowledge and will on breaking free. Dark purple energy swirls around your body and the M on your forehead burns brightly. After a few moments you give up, having failed to overcome the majinization!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@W$n focuses hard with $s eyes closed. Dark purple energy swirls around $s body and the M on $s head burns brightly. After a few moments $n seems to give up and the commotion dies down.@n"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				ch.Exp = 0
				ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
				if GET_LEVEL(ch) >= 100 {
					ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
				}
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou focus all your knowledge and will on breaking free. Dark purple energy swirls around your body and the M on your forehead burns brightly. After a few moments the ground splits beneath you and while letting out a piercing scream the M disappears from your forehead! You are free while still keeping the boost you had recieved from the majinization!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@W$n focuses hard with $s eyes closed. Dark purple energy swirls around $s body and the M on $s head burns brightly. After a few moments the ground beneath $n splits and $e lets out a piercing scream. The M on $s forehead disappears!@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Majinize = 3
				return
			}
		}
	}
}
func do_grapple(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_GRAPPLE) == 0 {
		return
	}
	if PLR_FLAGGED(ch, PLR_THANDW) {
		send_to_char(ch, libc.CString("Your are too busy wielding your weapon with two hands!\r\n"))
		return
	}
	if ch.Absorbing != nil {
		send_to_char(ch, libc.CString("You are currently absorbing from someone!\r\n"))
		return
	}
	if ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are currently being absorbed by someone! Try 'escape'!\r\n"))
		return
	}
	if ch.Grappling != nil {
		act(libc.CString("@RYou stop grappling with @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappling), TO_CHAR)
		act(libc.CString("@r$n@R stops grappling with @rYOU!!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappling), TO_VICT)
		act(libc.CString("@r$n@R stops grappling with @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Grappling), TO_NOTVICT)
		ch.Grappling.Grap = -1
		ch.Grappling.Grappled = nil
		ch.Grappling = nil
		ch.Grap = -1
		return
	}
	if ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are currently a victim of grappling! Try 'escape' to break free!\r\n"))
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	}
	var vict *char_data
	var arg [200]byte
	var arg2 [200]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: grapple (target) (hold | choke | grab)\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting == nil {
			send_to_char(ch, libc.CString("That target isn't here.\r\n"))
			return
		} else {
			vict = ch.Fighting
		}
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	if AFF_FLAGGED(vict, AFF_KNOCKED) {
		send_to_char(ch, libc.CString("They are unconcious. What would be the point?\r\n"))
		return
	}
	if vict.Grappled != nil {
		send_to_char(ch, libc.CString("They are currently in someone else's grasp!\r\n"))
		return
	}
	if vict.Absorbby != nil {
		send_to_char(ch, libc.CString("They are currently in someone else's grasp!\r\n"))
		return
	}
	if vict.Absorbing != nil {
		send_to_char(ch, libc.CString("They are currently absorbing from someone!\r\n"))
		return
	}
	var pass int = FALSE
	_ = pass
	if C.strcasecmp(libc.CString("hold"), &arg2[0]) == 0 || C.strcasecmp(libc.CString("choke"), &arg2[0]) == 0 || C.strcasecmp(libc.CString("grab"), &arg2[0]) == 0 || C.strcasecmp(libc.CString("wrap"), &arg2[0]) == 0 {
		pass = TRUE
		var perc int = GET_SKILL(ch, SKILL_GRAPPLE)
		var prob int = axion_dice(0)
		var cost int = int(ch.Max_move / 100)
		if ch.Move < int64(cost) {
			send_to_char(ch, libc.CString("You do not have enough stamina to grapple!\r\n"))
			return
		}
		if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				reveal_hiding(ch, 0)
				act(libc.CString("@C$N@c disappears, avoiding your grapple attempt before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c grapple attempt before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c grapple attempt before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Move -= int64(cost)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				return
			} else {
				reveal_hiding(ch, 0)
				act(libc.CString("@C$N@c disappears, trying to avoid your grapple but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the grapple attempt but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c grapple attempt but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if GET_SPEEDI(ch) > GET_SPEEDI(vict)*2 {
			perc += 5
		} else if GET_SPEEDI(ch) > GET_SPEEDI(vict) {
			perc += 2
		} else if GET_SPEEDI(ch)*2 < GET_SPEEDI(vict) {
			perc -= 5
		} else if GET_SPEEDI(ch) < GET_SPEEDI(vict) {
			perc -= 2
		}
		if (float64(ch.Hit)*0.02)*float64(ch.Aff_abils.Str) < (float64(vict.Hit)*0.01)*float64(vict.Aff_abils.Str) {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou try to grapple with @r$N@R, but $E manages to overpower you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R tries to grapple with YOU, but you manage to overpower $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R tries to grapple with @r$N@R, but $E manages to overpower @r$n@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if (float64(ch.Hit)*0.01)*float64(ch.Aff_abils.Str) < (float64(vict.Hit)*0.01)*float64(vict.Aff_abils.Str) && rand_number(1, 4) == 1 {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou try to grapple with @r$N@R, but $E manages to overpower you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R tries to grapple with YOU, but you manage to overpower $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R tries to grapple with @r$N@R, but $E manages to overpower @r$n@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if perc < prob {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou try to grapple with @r$N@R, but $E manages to avoid it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R tries to grapple with YOU, but you manage to avoid it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R tries to grapple with @r$N@R, but $E manages to avoid it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if !HAS_ARMS(vict) && C.strcasecmp(libc.CString("grab"), &arg2[0]) == 0 {
			send_to_char(ch, libc.CString("They don't even have an arm to grab onto!\r\n"))
			return
		} else if C.strcasecmp(libc.CString("hold"), &arg2[0]) == 0 {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou rush at @r$N@R and manage to get $M in a hold from behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R rushes at YOU and manages to get you in a hold from behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R rushes at @r$N@R and manages to get $M in a hold from behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Grappling = vict
			ch.Grap = 1
			vict.Grappled = ch
			vict.Grap = 1
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if C.strcasecmp(libc.CString("choke"), &arg2[0]) == 0 {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou rush at @r$N@R and manage to grab $S throat with both hands!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R rushes at YOU and manages to grab your throat with both hands!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R rushes at @r$N@R and manages to grab $S throat with both hands!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Grappling = vict
			ch.Grap = 2
			vict.Grappled = ch
			vict.Grap = 2
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if C.strcasecmp(libc.CString("wrap"), &arg2[0]) == 0 {
			if ch.Race != RACE_MAJIN {
				send_to_char(ch, libc.CString("Your body is not flexible enough to wrap around a target!\r\n"))
				return
			}
			act(libc.CString("@MMoving quickly you stretch your body out and wrap it around the length of @c$N's@M body! You tighten your body until you begin crushing @c$N@M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@M quickly stretches out $s body and wraps it around @RYOU@M! You feel $s body begin to crush your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@M quickly stretches out $s body and wraps it around @c$N@M! It appears that @c$N's@M body is being crushed slowly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Grappling = vict
			ch.Grap = 4
			vict.Grappled = ch
			vict.Grap = 4
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		} else if C.strcasecmp(libc.CString("grab"), &arg2[0]) == 0 {
			reveal_hiding(ch, 0)
			act(libc.CString("@RYou rush at @r$N@R and manage to lock your arm onto $S!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n@R rushes at YOU and manages to lock $s arm onto your's!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R rushes at @r$N@R and manages to lock $s arm onto @r$N's@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Grappling = vict
			ch.Grap = 3
			vict.Grappled = ch
			vict.Grap = 3
			if !PLR_FLAGGED(vict, PLR_THANDW) {
				vict.Act[int(PLR_THANDW/32)] &= bitvector_t(^(1 << (int(PLR_THANDW % 32))))
			}
			ch.Move -= int64(cost)
			improve_skill(ch, SKILL_GRAPPLE, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			return
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: grapple (target) (hold | choke | grab | wrap)\r\n"))
		return
	}
}
func do_trip(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [200]byte
		vict *char_data = nil
	)
	one_argument(argument, &arg[0])
	if check_skill(ch, SKILL_TRIP) == 0 && !IS_NPC(ch) {
		return
	}
	var cost int = int(ch.Max_hit / 200)
	if cost > int(ch.Move) {
		send_to_char(ch, libc.CString("You don't have enough stamina.\r\n"))
		return
	}
	var perc int = init_skill(ch, SKILL_TRIP)
	var prob int = rand_number(1, 114)
	if perc == 0 {
		perc = GET_LEVEL(ch) + rand_number(1, 10)
	}
	vict = nil
	vict = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("That target isn't here.\r\n"))
			return
		}
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	if vict != nil {
		if AFF_FLAGGED(vict, AFF_FLYING) {
			send_to_char(ch, libc.CString("They are flying and are not on their feet!\r\n"))
			return
		}
		if vict.Position == POS_SITTING {
			send_to_char(ch, libc.CString("They are not on their feet!\r\n"))
			return
		}
		if PLR_FLAGGED(vict, PLR_HEALT) {
			send_to_char(ch, libc.CString("They are inside a healing tank!\r\n"))
			return
		}
		if GET_SPEEDI(ch) > GET_SPEEDI(vict)*2 {
			perc += 5
		} else if GET_SPEEDI(ch) > GET_SPEEDI(vict) {
			perc += 2
		} else if GET_SPEEDI(ch)*2 < GET_SPEEDI(vict) {
			perc -= 5
		} else if GET_SPEEDI(ch) < GET_SPEEDI(vict) {
			perc -= 2
		}
		if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				reveal_hiding(ch, 0)
				act(libc.CString("@C$N@c disappears, avoiding your trip before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c trip before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c trip before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Move -= int64(cost)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
				return
			} else {
				reveal_hiding(ch, 0)
				act(libc.CString("@C$N@c disappears, trying to avoid your trip but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the trip but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c trip but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if perc < prob {
			reveal_hiding(ch, 0)
			act(libc.CString("@mYou move to trip $N@m, but you screw up and $E keeps $S footing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@m$n@m moves to trip YOU, but $e screws up and you manage to keep your footing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@m$n@m moves to trip $N@m, but $e screws up and $N@m manages to keep $S footing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_TRIP, 0)
			ch.Move -= int64(cost)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			} else if ch.Fighting != vict {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			} else if vict.Fighting != ch {
				set_fighting(vict, ch)
			}
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@mYou move to trip $N@m, and manage to knock $M off $S feet!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@m$n@m moves to trip YOU, and manages to knock you off your feet!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@m$n@m moves to trip $N@m, and manages to knock $N@m off $S feet!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_TRIP, 0)
			ch.Move -= int64(cost)
			vict.Position = POS_SITTING
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			}
			return
		}
	} else {
		send_to_char(ch, libc.CString("ERROR: Report to Iovan.\r\n"))
		return
	}
}
func do_train(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if ch.Carry_weight > int(max_carry_weight(ch)) {
		send_to_char(ch, libc.CString("You are weighted down too much!\r\n"))
		return
	}
	var plus int = 0
	var total int64 = 0
	var weight int64 = 0
	var bonus int64 = 0
	var cost int64 = 0
	var arg [200]byte
	one_argument(argument, &arg[0])
	weight = int64(gear_weight(ch) - ch.Carry_weight)
	var strcap int = 5000
	var spdcap int = 5000
	var intcap int = 5000
	var wiscap int = 5000
	var concap int = 5000
	var aglcap int = 5000
	strcap += int(ch.Real_abils.Str * (-12))
	intcap += int(ch.Real_abils.Intel * (-12))
	wiscap += int(ch.Real_abils.Wis * (-12))
	spdcap += int(ch.Real_abils.Cha * (-12))
	concap += int(ch.Real_abils.Con * (-12))
	aglcap += int(ch.Real_abils.Dex * (-12))
	if ch.Race == RACE_HUMAN {
		intcap = int(float64(intcap) * 0.75)
		wiscap = int(float64(wiscap) * 0.75)
	} else if ch.Race == RACE_KANASSAN {
		intcap = int(float64(intcap) * 0.4)
		wiscap = int(float64(wiscap) * 0.4)
		aglcap = int(float64(aglcap) * 0.4)
	} else if ch.Race == RACE_HALFBREED {
		intcap = int(float64(intcap) * 0.75)
		strcap = int(float64(strcap) * 0.75)
	} else if ch.Race == RACE_TRUFFLE {
		strcap = int(float64(strcap) * 1.5)
		concap = int(float64(concap) * 1.5)
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("@D-------------[ @GTraining Status @D]-------------@n\r\n"))
		send_to_char(ch, libc.CString("  @mStrength Progress    @D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Trainstr)), func() string {
			if ch.Real_abils.Str >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(strcap)))
		}())
		send_to_char(ch, libc.CString("  @mSpeed Progress       @D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Trainspd)), func() string {
			if ch.Real_abils.Cha >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(spdcap)))
		}())
		send_to_char(ch, libc.CString("  @mConstitution Progress@D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Traincon)), func() string {
			if ch.Real_abils.Con >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(concap)))
		}())
		send_to_char(ch, libc.CString("  @mIntelligence Progress@D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Trainint)), func() string {
			if ch.Real_abils.Intel >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(intcap)))
		}())
		send_to_char(ch, libc.CString("  @mWisdom Progress      @D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Trainwis)), func() string {
			if ch.Real_abils.Wis >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(wiscap)))
		}())
		send_to_char(ch, libc.CString("  @mAgility Progress     @D: @R%6s/%6s@n\r\n"), add_commas(int64(ch.Player_specials.Trainagl)), func() string {
			if ch.Real_abils.Dex >= 80 {
				return "@rCAPPED"
			}
			return libc.GoString(add_commas(int64(aglcap)))
		}())
		send_to_char(ch, libc.CString("@D  -----------------------------------------  @n\r\n"))
		send_to_char(ch, libc.CString("  @CCurrent Weight Worn  @D: @c%s@n\r\n"), add_commas(weight))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\r\n"))
		send_to_char(ch, libc.CString("Syntax: train (str | spd | agl | wis | int | con)\r\n"))
		return
	}
	if weight < int64(GET_LEVEL(ch)*100) && GET_LEVEL(ch) <= 19 {
		send_to_char(ch, libc.CString("With so little weight on you like that it would be a joke to try and train.\r\n"))
		return
	} else if weight < int64(GET_LEVEL(ch)*110) && GET_LEVEL(ch) <= 45 {
		send_to_char(ch, libc.CString("With so little weight on you like that it would be a joke to try and train.\r\n"))
		return
	} else if weight < int64(GET_LEVEL(ch)*125) && GET_LEVEL(ch) > 45 {
		send_to_char(ch, libc.CString("With so little weight on you like that it would be a joke to try and train.\r\n"))
		return
	}
	total = weight * int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity+1)
	total += int64(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity + 1) * ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity + 1))
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 6100 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 6135 {
		total += int64(float64(total) * 0.15)
	}
	var sensei int = -1
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1131 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GRoshi begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_ROSHI
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1131 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1131 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2F46 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GKibito begins to instruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_KABITO
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2F96 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2F96 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1714 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@Sixteen begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_ANDSIX
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1714 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1714 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4283 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GFrieza begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_FRIEZA
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4283 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4283 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x32D4 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GKrane begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_KRANE
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x32D4 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x32D4 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 2267 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GBardock begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_BARDOCK
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 2267 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 2267 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1662 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GPiccolo begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_PICCOLO
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1662 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 1662 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2DA4 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GNail begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_NAIL
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2DA4 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x2DA4 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4290 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GGinyu begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_GINYU
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4290 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 4290 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 8233 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GTapion begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_TAPION
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 8233 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 8233 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 6487 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GDabura begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_DABURA
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 6487 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 6487 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 3499 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GJinto begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_JINTO
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 3499 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 3499 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x3AA1 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GTsuna begins to intruct you in training technique@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_TSUNA
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x3AA1 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 0x3AA1 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 16100 && ch.Gold >= 8 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 1 {
		send_to_char(ch, libc.CString("@GKurzak begins to intruct you in training technique.@n\r\n"))
		total += int64(float64(total) * 0.85)
		sensei = CLASS_KURZAK
		if GET_LEVEL(ch) >= 100 {
			total *= 15000
		} else if GET_LEVEL(ch) >= 80 {
			total *= 1500
		} else if GET_LEVEL(ch) >= 40 {
			total *= 600
		} else if GET_LEVEL(ch) >= 20 {
			total *= 300
		} else if GET_LEVEL(ch) >= 10 {
			total *= 150
		}
		ch.Gold -= 8
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 1
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 16100 && ch.Gold < 8 {
		send_to_char(ch, libc.CString("@YYou do not have enough zenni (5) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 16100 && (ch.Player_specials.Class_skill_points[ch.Chclass]) < 1 {
		send_to_char(ch, libc.CString("@YYou do not have enough PS (1) in order to train with your sensei. Train on your own elsewhere.@n\r\n"))
	}
	if total > ch.Max_hit*2 {
		bonus = 5
	} else if total > ch.Max_hit {
		bonus = 4
	} else if total > (ch.Max_hit / 2) {
		bonus = 3
	} else if total > (ch.Max_hit / 4) {
		bonus = 2
	} else if total > (ch.Max_hit / 8) {
		bonus = 1
	}
	if sensei < 0 {
		cost = (total / 20) + ch.Max_move/50
	} else {
		cost = (total / 25) + ch.Max_move/60
	}
	if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
		cost -= int64(float64(cost) * 0.25)
	}
	if ch.Relax_count >= 464 {
		cost *= 10
	} else if ch.Relax_count >= 232 {
		cost *= 5
	} else if ch.Relax_count >= 116 {
		cost *= 2
	}
	if C.strcasecmp(libc.CString("str"), &arg[0]) == 0 {
		if ch.Real_abils.Str == 80 {
			send_to_char(ch, libc.CString("Your base strength is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Str == 25 && (ch.Bonuses[BONUS_WIMP]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your strength beyond 25.\r\n"))
			return
		}
		if ch.Move < cost {
			send_to_char(ch, libc.CString("You do not have enough stamina with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WYou throw a flurry of punches into the air at an invisible opponent.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n throws a flurry of punches into the air.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou leap into the air and throw a wild kick at an invisible opponent@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n leaps into the air and throws a wild kick at nothing.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou leap high into the air and unleash a flurry of punches and kicks at an invisible opponent@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n leaps high into the air and unleashes a flurry of punches and kicks at nothing.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_move/50) * 100) / ch.Max_move)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_BRAWNY]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		if sensei > -1 {
			plus += int(float64(plus) * 0.2)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Trainstr += plus + 5
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Trainstr += plus + 10
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Trainstr += plus + 25
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Trainstr += plus + 50
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Trainstr += plus + 100
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Trainstr += 1
			ch.Move -= cost
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Trainstr >= strcap {
			ch.Player_specials.Trainstr -= strcap
			send_to_char(ch, libc.CString("You feel your strength improve!@n\r\n"))
			ch.Real_abils.Str += 1
			if ch.Chclass == CLASS_PICCOLO && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else if C.strcasecmp(libc.CString("spd"), &arg[0]) == 0 {
		if ch.Real_abils.Cha == 80 {
			send_to_char(ch, libc.CString("Your base speed is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Cha == 25 && (ch.Bonuses[BONUS_SLOW]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your speed beyond 25.\r\n"))
			return
		}
		if ch.Move < cost {
			send_to_char(ch, libc.CString("You do not have enough stamina with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WYou dash quickly around the surrounding area as fast as you can!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n dashes quickly around the surrounding area as fast as $e can!@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou dodge to the side as fast as you can!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n dodges to the side as fast as $e can!@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou dash backwards as fast as you can!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n dashes backwards as fast as $e can!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_move/50) * 100) / ch.Max_move)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_QUICK]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Trainspd += plus + 5
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Trainspd += plus + 10
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Trainspd += plus + 25
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Trainspd += plus + 50
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Trainspd += plus + 100
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Trainspd += 1
			ch.Move -= cost
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Trainspd >= spdcap {
			ch.Player_specials.Trainspd -= spdcap
			send_to_char(ch, libc.CString("You feel your speed improve!@n\r\n"))
			ch.Real_abils.Cha += 1
			if ch.Chclass == CLASS_PICCOLO && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else if C.strcasecmp(libc.CString("con"), &arg[0]) == 0 {
		if ch.Real_abils.Con == 80 {
			send_to_char(ch, libc.CString("Your base constitution is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Con == 25 && (ch.Bonuses[BONUS_WIMP]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your constitution beyond 25.\r\n"))
			return
		}
		if ch.Move < cost {
			send_to_char(ch, libc.CString("You do not have enough stamina with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WYou leap into the air and then slam into the ground with your feet outstretched!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n leaps into the air and then slams into the ground with $s feet outstretched!?@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou leap into the air and then slam into the ground with your fists!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n leaps into the air and then slams into the ground with $s fists!?@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou leap into the air and then slam into the ground with your body!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n leaps into the air and then slams into the ground with $s body!?@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_move/50) * 100) / ch.Max_move)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_STURDY]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Traincon += plus + 5
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Traincon += plus + 10
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Traincon += plus + 25
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Traincon += plus + 50
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Traincon += plus + 100
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Traincon += 1
			ch.Move -= cost
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Traincon >= concap {
			ch.Player_specials.Traincon -= concap
			send_to_char(ch, libc.CString("You feel your constitution improve!@n\r\n"))
			ch.Real_abils.Con += 1
			if ch.Chclass == CLASS_PICCOLO && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else if C.strcasecmp(libc.CString("agl"), &arg[0]) == 0 {
		if ch.Real_abils.Dex == 80 {
			send_to_char(ch, libc.CString("Your base agility is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Dex == 25 && (ch.Bonuses[BONUS_CLUMSY]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your agility beyond 25.\r\n"))
			return
		}
		if ch.Move < cost {
			send_to_char(ch, libc.CString("You do not have enough stamina with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WYou do a series of backflips through the air, landing gracefully on one foot a moment later.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n does a series of backflips through the air, landing gracefully on one foot a moment later.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou flip forward and launch off your hands into the air. You land gracefully on one foot a moment later.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n flips forward and launches off $s hands into the air. Then $e lands gracefully on one foot a moment later.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou flip to the side off one hand and then land on your feet.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n flips to the side off one hand and then lands on $s feet.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_move/50) * 100) / ch.Max_move)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_AGILE]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Trainagl += plus + 5
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Trainagl += plus + 10
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Trainagl += plus + 25
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Trainagl += plus + 50
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Trainagl += plus + 100
			ch.Move -= cost
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Trainagl += 1
			ch.Move -= cost
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Trainagl >= aglcap {
			ch.Player_specials.Trainagl -= aglcap
			send_to_char(ch, libc.CString("You feel your agility improve!@n\r\n"))
			ch.Real_abils.Dex += 1
			if ch.Chclass == CLASS_PICCOLO && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else if C.strcasecmp(libc.CString("int"), &arg[0]) == 0 {
		if ch.Real_abils.Intel == 80 {
			send_to_char(ch, libc.CString("Your base intelligence is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Intel == 25 && (ch.Bonuses[BONUS_DULL]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your intelligence beyond 25.\r\n"))
			return
		}
		if ch.Mana < ((total / 20) + ch.Max_mana/50) {
			send_to_char(ch, libc.CString("You do not have enough ki with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WConcentrating you fly high into the air as fast as you can before settling slowly back to the ground.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n flies high into the air as fast as $e can before settling slowly back to the ground.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou focus your ki at your outstretched hand and send a mild shockwave in that direction!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n focuses $s ki at $s outstretched hand and sends a mild shockwave in that direction!@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou concentrate on your ki and force torrents of it to rush out from your body randomly!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n seems to concentrate before torrents of ki randomly blasts out from $s body!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_mana/50) * 100) / ch.Max_mana)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_SCHOLARLY]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Trainint += plus + 5
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Trainint += plus + 10
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Trainint += plus + 25
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Trainint += plus + 50
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Trainint += plus + 100
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Trainint += 1
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Trainint >= intcap {
			ch.Player_specials.Trainint -= intcap
			send_to_char(ch, libc.CString("You feel your intelligence improve!@n\r\n"))
			ch.Real_abils.Intel += 1
			if ch.Chclass == CLASS_PICCOLO && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else if C.strcasecmp(libc.CString("wis"), &arg[0]) == 0 {
		if ch.Real_abils.Wis == 80 {
			send_to_char(ch, libc.CString("Your base wisdom is maxed!\r\n"))
			return
		}
		if ch.Real_abils.Wis == 25 && (ch.Bonuses[BONUS_FOOLISH]) > 0 {
			send_to_char(ch, libc.CString("You're not able to withstand increasing your wisdom beyond 25.\r\n"))
			return
		}
		if ch.Mana < ((total / 20) + ch.Max_mana/50) {
			send_to_char(ch, libc.CString("You do not have enough ki with the current weight worn and gravity!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("@WYou close your eyes and wage a mental battle against an imaginary opponent.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n closes $s eyes for a moment and an expression of intensity forms on it.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou look around and contemplate battle tactics for an imaginary scenario.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n looks around and appears to be imagining things that aren't there.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou invent a battle plan for a battle that doesn't exist!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n seems to have thought of something.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		plus = int((((total / 20) + ch.Max_mana/50) * 100) / ch.Max_mana)
		if GET_LEVEL(ch) > 80 {
			plus += 50
		} else if GET_LEVEL(ch) > 60 {
			plus += 25
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			plus *= 4
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			plus *= 3
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			plus += int(float64(plus) * 0.25)
		}
		if (ch.Bonuses[BONUS_SAGE]) != 0 {
			plus += int(float64(plus) * 0.75)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			plus += int(float64(plus) * 0.05)
		}
		switch bonus {
		case 1:
			ch.Player_specials.Trainwis += plus + 5
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel slight improvement. @D[@G+%d@D]@n\r\n"), plus+5)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 2:
			ch.Player_specials.Trainwis += plus + 10
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel some improvement. @D[@G+%d@D]@n\r\n"), plus+10)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 3:
			ch.Player_specials.Trainwis += plus + 25
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel good improvement. @D[@G+%d@D]@n\r\n"), plus+25)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		case 4:
			ch.Player_specials.Trainwis += plus + 50
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel great improvement! @D[@G+%d@D]@n\r\n"), plus+50)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		case 5:
			ch.Player_specials.Trainwis += plus + 100
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You feel awesome improvement! @D[@G+%d@D]@n\r\n"), plus+100)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		default:
			ch.Player_specials.Trainwis += 1
			if sensei < 0 {
				ch.Mana -= (total / 20) + ch.Max_mana/50
			} else {
				ch.Mana -= (total / 25) + ch.Max_mana/60
			}
			send_to_char(ch, libc.CString("You barely feel any improvement. @D[@G+1@D]@n\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if ch.Player_specials.Trainwis >= wiscap {
			ch.Player_specials.Trainwis -= wiscap
			send_to_char(ch, libc.CString("You feel your wisdom improve!@n\r\n"))
			ch.Real_abils.Wis += 1
			if ch.Chclass == CLASS_PICCOLO && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.25)
				send_to_char(ch, libc.CString("You gained quite a bit of experience from that!\r\n"))
			}
			save_char(ch)
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: train (str | spd | agl | wis | int | con)\r\n"))
		return
	}
}
func do_rip(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Rip the tail off who?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if !PLR_FLAGGED(vict, PLR_TAIL) && !PLR_FLAGGED(vict, PLR_STAIL) {
		send_to_char(ch, libc.CString("They do not have a tail to rip off!\r\n"))
		return
	}
	if ch != vict && ch.Position > POS_SLEEPING {
		if ch.Move < ch.Max_move/20 {
			send_to_char(ch, libc.CString("You are too tired to manage to grab their tail!\r\n"))
			return
		} else if GET_SPEEDI(ch) > GET_SPEEDI(vict) {
			ch.Move -= ch.Max_move / 20
			if ch.Hit > vict.Hit*2 {
				reveal_hiding(ch, 0)
				act(libc.CString("@rYou rush at @R$N@r and grab $S tail! With a powerful tug you pull it off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@R$n@r rushes at YOU and grabs your tail! With a powerful tug $e pulls it off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$n@R rushes at @R$N@r and grab $S tail! With a powerful tug $e pulls it off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Act[int(PLR_TAIL/32)] &= bitvector_t(^(1 << (int(PLR_TAIL % 32))))
				vict.Act[int(PLR_STAIL/32)] &= bitvector_t(^(1 << (int(PLR_STAIL % 32))))
				oozaru_drop(vict)
				return
			} else {
				reveal_hiding(ch, 0)
				act(libc.CString("@rYou rush at @R$N@r and grab $S tail! You are too weak to pull it off though!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@R$n@r rushes at YOU and grabs your tail! $e is too weak to pull it off though!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$n@R rushes at @R$N@r and grab $S tail! $e is too weak to pull it off though!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				return
			}
		} else {
			ch.Move -= ch.Max_move / 20
			reveal_hiding(ch, 0)
			act(libc.CString("@rYou rush at @R$N@r and try to grab $S tail, but fail!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@R$n@r rushes at YOU and tries to grab your tail, but fails!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$n@R rushes at @R$N@r and tries to grab $S tail, but fails!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			return
		}
	} else if ch == vict {
		reveal_hiding(ch, 0)
		act(libc.CString("@rYou grab your own tail and yank it off!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@R$n@r grabs $s own tail and yanks it off!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_TAIL/32)] &= bitvector_t(^(1 << (int(PLR_TAIL % 32))))
		ch.Act[int(PLR_STAIL/32)] &= bitvector_t(^(1 << (int(PLR_STAIL % 32))))
		oozaru_drop(vict)
	} else {
		if ch.Move < ch.Max_move/20 {
			send_to_char(ch, libc.CString("You are too tired to manage to grab their tail!\r\n"))
			return
		}
		ch.Move -= ch.Max_move / 20
		reveal_hiding(ch, 0)
		act(libc.CString("@rYou reach and grab @R$N's@r tail! With a powerful tug you pull it off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@RYou feel your tail pulled off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@R$n@R reaches and grabs @R$N's@r tail! With a powerful tug $e pulls it off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Act[int(PLR_TAIL/32)] &= bitvector_t(^(1 << (int(PLR_TAIL % 32))))
		vict.Act[int(PLR_STAIL/32)] &= bitvector_t(^(1 << (int(PLR_STAIL % 32))))
		oozaru_drop(vict)
		return
	}
}
func do_infuse(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_INFUSE) == 0 {
		return
	}
	if AFF_FLAGGED(ch, AFF_INFUSE) {
		act(libc.CString("You stop infusing ki into your attacks."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n stops infusing ki into $s attacks."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Affected_by[int(AFF_INFUSE/32)] &= ^(1 << (int(AFF_INFUSE % 32)))
		return
	}
	if ch.Mana < ch.Max_mana/100 {
		send_to_char(ch, libc.CString("You don't have enough ki to infuse into your attacks!\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	act(libc.CString("You start infusing ki into your attacks."), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("$n starts infusing ki into $s attacks."), TRUE, ch, nil, nil, TO_ROOM)
	ch.Affected_by[int(AFF_INFUSE/32)] |= 1 << (int(AFF_INFUSE % 32))
	ch.Mana -= ch.Max_mana / 100
}
func do_paralyze(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if know_skill(ch, SKILL_PARALYZE) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who are you wanting to paralyze?\r\n"))
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	if AFF_FLAGGED(vict, AFF_PARA) {
		send_to_char(ch, libc.CString("They are already partially paralyzed!\r\n"))
		return
	}
	if ch.Mana < vict.Hit/10+ch.Max_mana/20 {
		send_to_char(ch, libc.CString("You realize you can't paralyze them. You don't have enough ki to restrain them!\r\n"))
		return
	}
	var prob int = GET_SKILL(ch, SKILL_PARALYZE)
	var perc int = axion_dice(0)
	if GET_SPEEDI(ch)*2 < GET_SPEEDI(vict) {
		prob -= 10
	}
	if GET_SPEEDI(ch)+GET_SPEEDI(ch)/2 < GET_SPEEDI(vict) {
		prob -= 5
	}
	if (vict.Bonuses[BONUS_INSOMNIAC]) != 0 {
		ch.Mana -= vict.Hit/6 + ch.Max_mana/20
		act(libc.CString("@RYou focus ki and point both your arms at @r$N@R. However $N seems to shake off your paralysis attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at YOU! Your insomnia makes you immune to $s feeble paralysis attempt.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at @r$N@R. However $N seems to shake off $s paralysis attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		return
	} else if prob < perc {
		reveal_hiding(ch, 0)
		act(libc.CString("@RYou focus ki and point both your arms at @r$N@R. However $E manages to avoid your attempt to paralyze $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at YOU! You manage to avoid $s technique though...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at @r$N@R. However $E manages to avoid @r$n's@R attempted technique...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= vict.Hit/6 + ch.Max_mana/20
		improve_skill(ch, SKILL_PARALYZE, 0)
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("@RYou focus ki and point both your arms at @r$N@R. Your ki flows into $S body and partially paralyzes $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at YOU! You are caught in $s paralysis technique and now can barely move!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@r$n @Rfocuses ki and points both $s arms at @r$N@R. @r$n's@R ki flows into @r$N@R body and partially paralyzes $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		var duration int = int(ch.Aff_abils.Intel / 15)
		assign_affect(vict, AFF_PARA, SKILL_PARALYZE, duration, 0, 0, 0, 0, 0, 0)
		ch.Mana -= vict.Hit/6 + ch.Max_mana/20
		improve_skill(ch, SKILL_PARALYZE, 0)
	}
}
func do_taisha(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_TAISHA) == 0 {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
		send_to_char(ch, libc.CString("This area already has an aura of regeneration around it.\r\n"))
		return
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity > 0 {
		send_to_char(ch, libc.CString("This area's gravity is too hostile to an aura.\r\n"))
		return
	}
	if ch.Mana < ch.Max_mana/3 {
		send_to_char(ch, libc.CString("You don't have enough ki.\r\n"))
		return
	}
	var prob int = GET_SKILL(ch, SKILL_TAISHA)
	var perc int = axion_dice(0)
	if prob < perc {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou hold up your hands while channeling ki. Your technique fails to produce an aura though....@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n@W holds up $s hands while channeling ki. $s technique fails to produce an aura though...."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= ch.Max_mana / 3
		improve_skill(ch, SKILL_TAISHA, 1)
		return
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou hold up your hands while channeling ki. Suddenly a @wburst@W of calming @Cblue@W light covers the surrounding area!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@g$n holds up $s hands while channeling ki. Suddenly a @wburst@W of calming @Cblue@W light covers the surrounding area!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= ch.Max_mana / 3
		improve_skill(ch, SKILL_TAISHA, 1)
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] |= bitvector_t(1 << (int(ROOM_AURA % 32)))
		return
	}
}
func do_kura(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_KURA) == 0 {
		return
	}
	if ch.Mana >= ch.Max_mana {
		send_to_char(ch, libc.CString("Your ki is already maxed out!\r\n"))
		return
	}
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: kuraiiro (1-100).\r\n"))
		return
	}
	var num int = libc.Atoi(libc.GoString(&arg[0]))
	var skill int = GET_SKILL(ch, SKILL_KURA)
	var cost int64 = 0
	var bonus int64 = 0
	if num > skill {
		send_to_char(ch, libc.CString("The number can not be greater than your skill.\r\n"))
		return
	}
	if num <= 0 {
		send_to_char(ch, libc.CString("The number can not be less than 1.\r\n"))
		return
	}
	cost = (ch.Max_mana / 100) * int64(num)
	bonus = cost
	if ch.Move < cost {
		send_to_char(ch, libc.CString("You do not have enough stamina for that high a number.\r\n"))
		return
	}
	if skill <= axion_dice(0) {
		ch.Move -= cost
		reveal_hiding(ch, 0)
		act(libc.CString("You crouch down and scream as your eyes turn red. You attempt to tap into your dark energies but you fail!"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@w crouches down and screams as $s eyes turn red and $e attempts to tap into dark energies but fails!"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_KURA, 0)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		ch.Move -= cost
		ch.Mana += bonus
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		reveal_hiding(ch, 0)
		act(libc.CString("You crouch down and scream as your eyes turn red. You attempt to tap into your dark energies and succeed as a rush of energy explodes around you!"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@w crouches down and screams as $s eyes turn red. Suddenly $e manages to tap into dark energies and a rush of energy explodes around $m!"), TRUE, ch, nil, nil, TO_ROOM)
		improve_skill(ch, SKILL_KURA, 0)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
}
func do_candy(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		obj  *obj_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if ch.Race != RACE_MAJIN {
		send_to_char(ch, libc.CString("You are not a majin, how can you do that?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Turn who into candy?\r\n"))
		return
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	if !IS_NPC(vict) {
		send_to_char(ch, libc.CString("You can't turn them into candy.\r\n"))
		return
	}
	if vict.Max_hit > ch.Max_hit*2 {
		send_to_char(ch, libc.CString("They are too powerful.\r\n"))
		return
	}
	if float64(vict.Max_hit) < float64(ch.Max_hit)*0.25 && GET_LEVEL(ch) < 100 {
		send_to_char(ch, libc.CString("They are too weak.\r\n"))
		return
	}
	if float64(vict.Max_hit) < float64(ch.Max_hit)*0.09 && GET_LEVEL(ch) == 100 {
		send_to_char(ch, libc.CString("They are too weak.\r\n"))
		return
	}
	if ch.Mana < ch.Max_mana/15 {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	if rand_number(1, 6) == 6 {
		ch.Mana -= ch.Max_mana / 15
		reveal_hiding(ch, 0)
		act(libc.CString("@cYou aim your forelock at @R$N@c and fire a beam of energy but it is dodged!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@c aims $s forelock at @R$N@c and fires a beam of energy but the beam is dodged!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		if ch.Fighting == nil {
			set_fighting(ch, vict)
		}
		if vict.Fighting == nil {
			set_fighting(vict, ch)
		}
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		return
	} else {
		ch.Mana -= ch.Max_mana / 15
		reveal_hiding(ch, 0)
		act(libc.CString("@cYou aim your forelock at @R$N@c and fire a beam of energy that envelopes $S entire body and changes $M into candy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@c aims $s forelock at @R$N@c and fires a beam of energy that envelopes $S entire body and changes $M into candy!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		if float64(vict.Max_hit) >= float64(ch.Max_hit)*1.5 {
			send_to_char(ch, libc.CString("You grab the candy as it falls.\r\n"))
			obj = read_object(95, VIRTUAL)
			obj_to_char(obj, ch)
		} else if vict.Max_hit >= ch.Max_hit {
			send_to_char(ch, libc.CString("You grab the candy as it falls.\r\n"))
			obj = read_object(94, VIRTUAL)
			obj_to_char(obj, ch)
		} else if vict.Max_hit < ch.Max_hit && (float64(vict.Max_hit) >= float64(ch.Max_hit)*0.5 && GET_LEVEL(ch) < 100 || float64(vict.Max_hit) >= float64(ch.Max_hit)*0.1 && GET_LEVEL(ch) >= 100) {
			send_to_char(ch, libc.CString("You grab the candy as it falls.\r\n"))
			obj = read_object(93, VIRTUAL)
			obj_to_char(obj, ch)
		} else if vict.Max_hit < ch.Max_hit && (float64(vict.Max_hit) < float64(ch.Max_hit)*0.5 && GET_LEVEL(ch) < 100 || float64(vict.Max_hit) < float64(ch.Max_hit)*0.1 && GET_LEVEL(ch) >= 100) {
			send_to_char(ch, libc.CString("You grab the candy as it falls.\r\n"))
			obj = read_object(53, VIRTUAL)
			obj_to_char(obj, ch)
		}
		vict.Act[int(MOB_HUSK/32)] |= bitvector_t(1 << (int(MOB_HUSK % 32)))
		die(vict, ch)
	}
}
func do_future(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		vict *char_data = nil
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) || ch.Race != RACE_KANASSAN {
		send_to_char(ch, libc.CString("You are incapable of this ability.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Bestow advance future sight on who?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Bestow advance future sight on who?\r\n"))
		return
	}
	if AFF_FLAGGED(vict, AFF_FUTURE) {
		send_to_char(ch, libc.CString("They already can see the future.\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("You can't target them, there would be no point.\r\n"))
		return
	}
	if ch.Mana < ch.Max_mana/40 {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 100 {
		send_to_char(ch, libc.CString("You do not have enough PS to activate or pass on this ability.\r\n"))
		return
	}
	if vict != ch {
		if vict.Real_abils.Cha+5 > 25 && (vict.Bonuses[BONUS_SLOW]) > 0 {
			send_to_char(ch, libc.CString("They can't handle having their speed increased beyond 25.\r\n"))
			return
		}
		if vict.Real_abils.Intel+2 > 25 && (vict.Bonuses[BONUS_DULL]) > 0 {
			send_to_char(ch, libc.CString("They can't handle having their intelligence increased beyond 25.\r\n"))
			return
		}
		ch.Mana -= ch.Max_mana / 40
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
		reveal_hiding(ch, 0)
		act(libc.CString("@CYou focus your energy into your fingers before stabbing your claws into $N and bestowing the power of Future Sight upon $M. Shortly after $E passes out.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n focuses $s energy into $s fingers before stabbing $s claws into YOUR neck and bestowing the power of Future Sight upon you! Soon after you pass out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n focuses $s energy into $s fingers before stabbing $s claws into $N's neck and bestowing the power of Future Sight upon $M! Soon after $E passes out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Affected_by[int(AFF_FUTURE/32)] |= 1 << (int(AFF_FUTURE % 32))
		vict.Real_abils.Cha += 5
		vict.Real_abils.Intel += 2
		vict.Position = POS_SLEEPING
		save_char(vict)
	} else {
		if ch.Real_abils.Cha+5 > 25 && (ch.Bonuses[BONUS_SLOW]) > 0 {
			send_to_char(ch, libc.CString("You can't handle having your speed increased beyond 25.\r\n"))
			return
		}
		if ch.Real_abils.Intel+2 > 25 && (ch.Bonuses[BONUS_DULL]) > 0 {
			send_to_char(ch, libc.CString("You can't handle having your intelligence increased beyond 25.\r\n"))
			return
		}
		ch.Mana -= ch.Max_mana / 40
		ch.Player_specials.Class_skill_points[ch.Chclass] -= 100
		reveal_hiding(ch, 0)
		act(libc.CString("@CYou focus your energy into your mind and awaken your latent Future Sight powers!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n focuses $s energy while closing $s eyes for a moment.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n focuses $s energy while closing $s eyes for a moment.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Affected_by[int(AFF_FUTURE/32)] |= 1 << (int(AFF_FUTURE % 32))
		ch.Real_abils.Cha += 5
		ch.Real_abils.Intel += 2
		vict.Position = POS_SLEEPING
		save_char(ch)
	}
}
func do_drag(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if ch.Drag != nil {
		vict = ch.Drag
		ch.Drag = nil
		vict.Dragged = nil
		act(libc.CString("@wYou stop dragging @C$N@W.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W stops dragging @c$N@W.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
		return
	}
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are busy carrying someone at the moment.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who do you want to drag?\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are a bit busy fighting right now!\r\n"))
		return
	}
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM || (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_SWIM {
		send_to_char(ch, libc.CString("You decide to not be a tugboat instead.\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Drag who?\r\n"))
		return
	}
	if vict == ch {
		send_to_char(ch, libc.CString("You can't drag yourself.\r\n"))
		return
	}
	if vict.Dragged != nil {
		send_to_char(ch, libc.CString("They are already being dragged!\r\n"))
		return
	}
	if IS_NPC(vict) && MOB_FLAGGED(vict, MOB_NOKILL) {
		send_to_char(ch, libc.CString("They are not to be touched!\r\n"))
		return
	}
	if vict.Position != POS_SLEEPING {
		reveal_hiding(ch, 0)
		act(libc.CString("@wYou try to grab and pull @C$N@W with you, but $E resists!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W tries to grab and pull you! However you resist!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W tries to grab and pull @c$N@W but $E resists!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		if IS_NPC(vict) && vict.Fighting == nil {
			set_fighting(vict, ch)
		}
		return
	} else if GET_PC_WEIGHT(vict)+vict.Carry_weight > int(max_carry_weight(ch)) {
		reveal_hiding(ch, 0)
		act(libc.CString("@wYou try to grab and pull @C$N@W with you, but $E is too heavy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W tries to grab and pull @c$N@W but $E is too heavy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
		return
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("@wYou grab and start dragging @C$N@W.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W grabs and starts dragging @c$N@W.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Drag = vict
		vict.Dragged = ch
		if !AFF_FLAGGED(vict, AFF_KNOCKED) && !AFF_FLAGGED(vict, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(vict, libc.CString("You feel your sleeping body being moved.\r\n"))
			if IS_NPC(vict) && vict.Fighting == nil {
				set_fighting(vict, ch)
			}
		}
	}
}
func do_stop(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if ch.Fighting == nil {
		send_to_char(ch, libc.CString("You are not even fighting!\r\n"))
		return
	} else {
		act(libc.CString("@CYou move out of your fighting posture.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@c$n@C moves out of $s fighting posture.@n"), TRUE, ch, nil, nil, TO_ROOM)
		stop_fighting(ch)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
}
func do_suppress(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if ch.Race == RACE_ANDROID {
		send_to_char(ch, libc.CString("You are unable to suppress your powerlevel.\r\n"))
		return
	}
	if (ch.Bonuses[BONUS_ARROGANT]) > 0 {
		send_to_char(ch, libc.CString("You are far too arrogant to hide your strength.\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_POWERUP) {
		send_to_char(ch, libc.CString("You are currently powering up, can't suppress.\r\n"))
		return
	}
	if ch.Kaioken != 0 {
		send_to_char(ch, libc.CString("You are currently concentrating on kaioken!\r\n"))
		return
	}
	if IS_NONPTRANS(ch) && ch.Race != RACE_ICER && IS_TRANSFORMED(ch) {
		send_to_char(ch, libc.CString("You must revert before you try and suppress.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Suppress to what percent?\r\nSyntax: suppress (1 - 99 | release)\r\n"))
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("release")) == 0 {
		if ch.Suppression != 0 {
			reveal_hiding(ch, 0)
			act(libc.CString("@GYou stop suppressing your current powerlevel!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@G$n smiles as a rush of power erupts around $s body briefly.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Hit += ch.Suppressed
			ch.Suppressed = 0
			if ch.Hit > ch.Max_hit {
				ch.Hit = gear_pl(ch)
			}
			ch.Suppression = 0
			return
		} else {
			send_to_char(ch, libc.CString("You are not suppressing!\r\n"))
			return
		}
	}
	var num int = libc.Atoi(libc.GoString(&arg[0]))
	if num > 99 || num <= 0 {
		send_to_char(ch, libc.CString("Out of suppression range.\r\nSyntax: suppress (1 - 99 | release)\r\n"))
		return
	}
	var max int64 = gear_pl(ch)
	var amt int64 = int64((float64(max) * 0.01) * float64(num))
	if ch.Hit < amt && ch.Suppression == 0 {
		send_to_char(ch, libc.CString("You are already below %d percent of your max!\r\n"), num)
		return
	} else if ch.Hit < amt && ch.Suppression != 0 {
		if ch.Suppressed+ch.Hit < amt {
			send_to_char(ch, libc.CString("You do not have enough suppressed to raise your powerlevel by that much.\r\n"))
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@GYou alter your suppression level!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@G$n seems to concentrate for a moment.@n"), TRUE, ch, nil, nil, TO_ROOM)
			var diff int64 = amt - ch.Hit
			ch.Suppressed -= diff
			ch.Hit += diff
			ch.Suppression = int64(num)
			return
		}
	} else if ch.Hit > amt && ch.Suppression != 0 {
		reveal_hiding(ch, 0)
		act(libc.CString("@GYou alter your suppression level!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@G$n seems to concentrate for a moment.@n"), TRUE, ch, nil, nil, TO_ROOM)
		var diff int64 = ch.Hit - amt
		ch.Suppressed += diff
		ch.Hit -= diff
		ch.Suppression = int64(num)
		return
	}
	if ch.Hit == amt {
		send_to_char(ch, libc.CString("You are already at %d percent of your max!\r\n"), num)
		return
	}
	reveal_hiding(ch, 0)
	act(libc.CString("@GYou suppress your current powerlevel!@n"), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("@G$n seems to concentrate for a moment.@n"), TRUE, ch, nil, nil, TO_ROOM)
	var diff int64 = ch.Hit - amt
	ch.Suppressed = diff
	ch.Hit = amt
	ch.Suppression = int64(num)
	return
}
func do_hass(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc int = 0
		prob int = 0
	)
	if check_skill(ch, SKILL_HASSHUKEN) == 0 {
		return
	}
	if ch.Move < ch.Max_move/30 {
		send_to_char(ch, libc.CString("You do not have enough stamina.\r\n"))
		return
	}
	perc = init_skill(ch, SKILL_HASSHUKEN)
	prob = axion_dice(0)
	if perc < prob {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou try to move your arms at incredible speeds but screw up and waste some of your stamina.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W tries to move $s arms at incredible speeds but screws up and wastes some of $s stamina.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Move -= ch.Max_move / 30
		improve_skill(ch, SKILL_HASSHUKEN, 0)
		return
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou concentrate and start to move your arms at incredible speeds.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W concentrates and starts to move $s arms at incredible speeds.@n"), TRUE, ch, nil, nil, TO_ROOM)
		var duration int = perc / 15
		assign_affect(ch, AFF_HASS, SKILL_HASSHUKEN, duration, 0, 0, 0, 0, 0, 0)
		ch.Move -= ch.Max_move / 30
		improve_skill(ch, SKILL_HASSHUKEN, 0)
		return
	}
}
func do_implant(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		limb     *obj_data = nil
		obj      *obj_data = nil
		next_obj *obj_data
		vict     *char_data = nil
		arg      [2048]byte
		arg2     [2048]byte
		found    int = FALSE
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: implant (rarm | larm | rleg | lleg) (target)\r\n"))
		return
	}
	if arg2[0] == 0 {
		vict = ch
	} else if (func() *char_data {
		vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That person isn't here.\r\n"))
		return
	}
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if found == FALSE && GET_OBJ_VNUM(obj) == 66 && (!OBJ_FLAGGED(obj, ITEM_BROKEN) && !OBJ_FLAGGED(obj, ITEM_FORGED)) {
			found = TRUE
			limb = obj
		}
	}
	if found == FALSE {
		send_to_char(ch, libc.CString("You do not have a cybernetic limb to implant.\r\n"))
		return
	}
	if C.strcmp(&arg[0], libc.CString("rarm")) == 0 {
		if (vict.Limb_condition[0]) >= 1 {
			if vict != ch {
				send_to_char(ch, libc.CString("They already have a right arm!\r\n"))
			}
			if vict == ch {
				send_to_char(ch, libc.CString("You already have a right arm!\r\n"))
			}
			return
		} else {
			if vict != ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new right arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places a $p@W up to your body. It automaticly adjusts itself, becoming a new right arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places a $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new right arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_NOTVICT)
			}
			if vict == ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to your body. It automaticly adjusts itself, becoming a new right arm!@n"), TRUE, ch, limb, nil, TO_CHAR)
				act(libc.CString("@C$n@W places the $p@W up to $s body. It automaticly adjusts itself, becoming a new right arm!@n"), TRUE, ch, limb, nil, TO_ROOM)
			}
			vict.Act[int(PLR_CRARM/32)] |= bitvector_t(1 << (int(PLR_CRARM % 32)))
			obj_from_char(limb)
			extract_obj(limb)
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("larm")) == 0 {
		if (vict.Limb_condition[1]) >= 1 {
			if vict != ch {
				send_to_char(ch, libc.CString("They already have a left arm!\r\n"))
			}
			if vict == ch {
				send_to_char(ch, libc.CString("You already have a left arm!\r\n"))
			}
			return
		} else {
			if vict != nil && vict != ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new left arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places a $p@W up to your body. It automaticly adjusts itself, becoming a new left arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places a $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new left arm!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_NOTVICT)
			}
			if vict == ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to your body. It automaticly adjusts itself, becoming a new left arm!@n"), TRUE, ch, limb, nil, TO_CHAR)
				act(libc.CString("@C$n@W places the $p@W up to $s body. It automaticly adjusts itself, becoming a new left arm!@n"), TRUE, ch, limb, nil, TO_ROOM)
			}
			vict.Act[int(PLR_CLARM/32)] |= bitvector_t(1 << (int(PLR_CLARM % 32)))
			obj_from_char(limb)
			extract_obj(limb)
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("rleg")) == 0 {
		if (vict.Limb_condition[2]) >= 1 {
			if vict != ch {
				send_to_char(ch, libc.CString("They already have a right leg!\r\n"))
			}
			if vict == ch {
				send_to_char(ch, libc.CString("You already have a right leg!\r\n"))
			}
			return
		} else {
			if vict != nil && vict != ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new right leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places a $p@W up to your body. It automaticly adjusts itself, becoming a new right leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places a $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new right leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_NOTVICT)
			}
			if vict == ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to your body. It automaticly adjusts itself, becoming a new right leg!@n"), TRUE, ch, limb, nil, TO_CHAR)
				act(libc.CString("@C$n@W places the $p@W up to $s body. It automaticly adjusts itself, becoming a new right leg!@n"), TRUE, ch, limb, nil, TO_ROOM)
			}
			vict.Act[int(PLR_CRLEG/32)] |= bitvector_t(1 << (int(PLR_CRLEG % 32)))
			obj_from_char(limb)
			extract_obj(limb)
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("lleg")) == 0 {
		if (vict.Limb_condition[3]) >= 1 {
			if vict != ch {
				send_to_char(ch, libc.CString("They already have a left leg!\r\n"))
			}
			if vict == ch {
				send_to_char(ch, libc.CString("You already have a left leg!\r\n"))
			}
			return
		} else {
			if vict != nil && vict != ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new left leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places a $p@W up to your body. It automaticly adjusts itself, becoming a new left leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places a $p@W up to @c$N@W's body. It automaticly adjusts itself, becoming a new left leg!@n"), TRUE, ch, limb, unsafe.Pointer(vict), TO_NOTVICT)
			}
			if vict == nil || vict == ch {
				reveal_hiding(ch, 0)
				act(libc.CString("@WYou place the $p@W up to your body. It automaticly adjusts itself, becoming a new left leg!@n"), TRUE, ch, limb, nil, TO_CHAR)
				act(libc.CString("@C$n@W places the $p@W up to $s body. It automaticly adjusts itself, becoming a new left leg!@n"), TRUE, ch, limb, nil, TO_ROOM)
			}
			vict.Act[int(PLR_CLLEG/32)] |= bitvector_t(1 << (int(PLR_CLLEG % 32)))
			obj_from_char(limb)
			extract_obj(limb)
			return
		}
	} else {
		send_to_char(ch, libc.CString("Syntax: implant (rarm | larm | rleg | rleg)\r\n"))
		return
	}
}
func do_pose(ch *char_data, argument *byte, cmd int, subcmd int) {
	if know_skill(ch, SKILL_POSE) == 0 {
		return
	}
	if ch.Move < ch.Max_move/40 {
		send_to_char(ch, libc.CString("You do not have enough stamina to pull off such an exciting pose!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_POSE) {
		send_to_char(ch, libc.CString("You are already feeling good and confident from a previous pose.\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are too busy to pose right now!\r\n"))
		return
	}
	var prob int = GET_SKILL(ch, SKILL_POSE)
	var perc int = rand_number(1, 70)
	if ch.Real_abils.Str+8 > 25 && (ch.Bonuses[BONUS_WIMP]) > 0 {
		send_to_char(ch, libc.CString("You can't handle having your strength increased beyond 25.\r\n"))
		return
	}
	if ch.Real_abils.Dex+8 > 25 && (ch.Bonuses[BONUS_CLUMSY]) > 0 {
		send_to_char(ch, libc.CString("You can't handle having your agility increased beyond 25.\r\n"))
		return
	}
	if prob < perc {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou attempt to strike an awe inspiring pose, but end up falling on your face!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W attempts to strike an awe inspiring pose, but ends up falling on $s face!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Move -= ch.Max_move / 40
		improve_skill(ch, SKILL_POSE, 0)
		return
	} else {
		reveal_hiding(ch, 0)
		switch rand_number(1, 4) {
		case 1:
			act(libc.CString("@WYou turn around with your back to everyone. You bend forward dramatically and put your head between your legs!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W turns around with $s back to you. $e bends forward dramatically and puts $s head between $s legs. Strange...@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("@WYou turn to the side while flexing your muscles and extend your arms up at an angle dramatically!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W turns to the side while flexing $s muscles and extending $s arms up at an angle dramatically!@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("@WYou extend one leg outward while you bend forward, balancing on a single leg!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W extends one leg outward while $e bends forward, balancing on a single leg!@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 4:
			act(libc.CString("@WYou drop down to one knee while angling your arms up to either side and slanting your hands down like wings!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W drops down to one knee while angling $s arms up to either side and slanting $s hands down like wings!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		send_to_char(ch, libc.CString("@WYou feel your confidence increase! @G+3 Str @Wand@G +3 Agl!@n\r\n"))
		ch.Real_abils.Str += 8
		ch.Real_abils.Dex += 8
		save_char(ch)
		var before int64 = int64(GET_LIFEMAX(ch))
		ch.Act[int(PLR_POSE/32)] |= bitvector_t(1 << (int(PLR_POSE % 32)))
		ch.Lifeforce += int64(GET_LIFEMAX(ch) - int(before))
		if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
			ch.Lifeforce = int64(GET_LIFEMAX(ch))
		}
		ch.Move -= ch.Max_move / 40
		improve_skill(ch, SKILL_POSE, 0)
		return
	}
}
func do_fury(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if ch.Race != RACE_HALFBREED || IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are furious, but you'll get over it.\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_FURY) {
		send_to_char(ch, libc.CString("You are already furious, your next attack will devestate, hurry use it!\r\n"))
		return
	}
	if int(ch.Fury) < 100 {
		send_to_char(ch, libc.CString("You do not have enough anger to release your fury upon your foes!\r\n"))
		return
	}
	if arg[0] == 0 {
		if ch.Hit < gear_pl(ch) {
			if float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.2 {
				ch.Hit = gear_pl(ch)
				ch.Lifeforce -= int64(float64(GET_LIFEMAX(ch)) * 0.2)
			} else {
				ch.Hit += ch.Lifeforce
				ch.Lifeforce = -1
				if ch.Hit > gear_pl(ch) {
					ch.Hit = gear_pl(ch)
				}
			}
		}
		ch.Fury = 0
	} else if C.strcasecmp(&arg[0], libc.CString("attack")) == 0 {
		ch.Fury = 50
	} else {
		send_to_char(ch, libc.CString("Syntax: fury (attack) <--- this will not use up your LF to restore PL.\n        fury <--- fury by itself will do both LF to PL restore and attack boost.\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	act(libc.CString("You release your fury! Your very next attack is guaranteed to rip your foes a new one!"), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("$n screams furiously as a look of anger appears on $s face!"), TRUE, ch, nil, nil, TO_ROOM)
	ch.Act[int(PLR_FURY/32)] |= bitvector_t(1 << (int(PLR_FURY % 32)))
}
func hint_system(ch *char_data, num int) {
	var hints [22]*byte = [22]*byte{libc.CString("Remember to save often."), libc.CString("Remember to eat or drink if you want to stay alive."), libc.CString("It is a good idea to save up PS for learning skills instead of just practicing them."), libc.CString("A good way to save up money is with the bank."), libc.CString("If you want to stay alive in this rough world you will need to be mindful of your surroundings."), libc.CString("Knowing when to rest and recover can be the difference between life and death."), libc.CString("Not every battle can be won. Great warriors know how to pick their fights."), libc.CString("It is a good idea to experiment with skills fully before deciding their worth."), libc.CString("Having a well balanced repertoire of skills can help you out of any situation."), libc.CString("You can become hidden from your enemies on who and ooc with the whohide command."), libc.CString("You can value an item at a shopkeeper with the value command."), libc.CString("There are ways to earn money through jobs, try looking for a job. Bum."), libc.CString("You never know what may be hidden nearby. You should always check out anything you can."), libc.CString("You should check for a help file on any subject you can, you never know how the info may 'help' you."), libc.CString("Until you are capable of taking care of yourself for long periods of time you should stick near your sensei."), libc.CString("You shouldn't travel to other planets until you have a stable supply of money."), libc.CString("There is a vast galaxy out there that you may not be able to reach by public ship."), libc.CString("Score is used to view the various statistics about your character."), libc.CString("Status is used to view what is influencing your character and its characteristics."), libc.CString("You will need a scouter in order to use the Scouter Network (SNET)."), libc.CString("The DBAT forum is a great resource for finding out information and for conversing\r\nwith fellow players. http://advent-truth.com/forum"), libc.CString("Found a bug or have a suggestion? Log into our forums and post in the relevant section.")}
	if num == 0 {
		num = rand_number(0, 21)
	}
	if ch.Race != RACE_ANDROID && ch.Race != RACE_NAMEK {
		send_to_char(ch, libc.CString("@D[@GHint@D] @G%s@n\r\n"), hints[num])
	} else {
		if num == 1 {
			num = 0
		}
		send_to_char(ch, libc.CString("@D[@GHint@D] @G%s@n\r\n"), hints[num])
	}
	send_to_char(ch, libc.CString("@D(@gYou can turn off hints with the command 'hints'@D)@n\r\n"))
}
func do_think(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	if IS_NPC(ch) {
		return
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("Lol, no.\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_TELEPATHY) != 0 {
		send_to_char(ch, libc.CString("You can just use telepathy.\r\n"))
		return
	}
	if ch.Mindlink == nil {
		send_to_char(ch, libc.CString("No one has linked with your mind.\r\n"))
		return
	}
	if *argument == 0 {
		send_to_char(ch, libc.CString("Syntax: think (message)\r\n"))
		return
	} else {
		var tch *char_data
		tch = ch.Mindlink
		send_to_char(ch, libc.CString("@c%s@w reads your thoughts, '@C%s@w'@n\r\n"), GET_NAME(tch), argument)
		send_to_char(tch, libc.CString("@c%s@w thinks, '@C%s@w'@n\r\n"), GET_NAME(ch), argument)
		send_to_imm(libc.CString("@GTELEPATHY: @C%s@G telepaths @c%s, @W'@w%s@W'@n"), func() *byte {
			if ch.Admlevel > 0 {
				return GET_NAME(ch)
			}
			return GET_USER(ch)
		}(), func() *byte {
			if tch.Admlevel > 0 {
				return GET_NAME(tch)
			}
			return GET_USER(tch)
		}(), argument)
		return
	}
}
func do_telepathy(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
		arg2 [2048]byte
	)
	half_chop(argument, &arg[0], &arg2[0])
	if know_skill(ch, SKILL_TELEPATHY) == 0 {
		return
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("Lol, no.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: telepathy [ read ] (target)\r\n        telepathy [ link ] (target)\r\n        telepathy [  far ] (target)\r\n        telepathy (target) (message)\r\n"))
		return
	} else if ch.Mana < ch.Max_mana/40 {
		send_to_char(ch, libc.CString("You do not have enough ki to focus your mental abilities.\r\n"))
		return
	}
	if C.strcmp(&arg[0], libc.CString("far")) == 0 {
		if IS_NPC(ch) {
			return
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Look through who's eyes?\r\n"))
			return
		} else if vict == ch {
			send_to_char(ch, libc.CString("Oh that makes a lot of sense...\r\n"))
			return
		} else if IS_NPC(vict) {
			send_to_char(ch, libc.CString("You can't touch the mind of such a thing.\r\n"))
			return
		} else if vict.Admlevel > ch.Admlevel {
			send_to_char(ch, libc.CString("Their mental power oustrips your's by unfathomable measurements!\r\n"))
			return
		} else if AFF_FLAGGED(ch, AFF_SHOCKED) {
			send_to_char(ch, libc.CString("Your mind has been shocked by telepathic feedback! You are not able to use telepathy right now.\r\n"))
			return
		} else if vict.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't touch the mind of such an artificial being.\r\n"))
			return
		} else if GET_SKILL(vict, SKILL_TELEPATHY)+int(vict.Aff_abils.Intel) > GET_SKILL(ch, SKILL_TELEPATHY)+int(ch.Aff_abils.Intel) {
			send_to_char(ch, libc.CString("They throw off your attempt with their own telepathic abilities!\r\n"))
			return
		} else if ch.In_room == vict.In_room {
			send_to_char(ch, libc.CString("They are in the same room as you!\r\n"))
			return
		} else if AFF_FLAGGED(vict, AFF_BLIND) {
			send_to_char(ch, libc.CString("They are blind!\r\n"))
			return
		} else if PLR_FLAGGED(vict, PLR_EYEC) {
			send_to_char(ch, libc.CString("Their eyes are closed!\r\n"))
			return
		} else {
			look_at_room(vict.In_room, ch, 0)
			send_to_char(ch, libc.CString("You see all this through their eyes!\r\n"))
			if vict.Aff_abils.Intel > ch.Aff_abils.Intel {
				send_to_char(ch, libc.CString("You feel like someone was using your mind for something...\r\n"))
			}
			ch.Mana -= ch.Max_mana / 40
			return
		}
	}
	if C.strcmp(&arg[0], libc.CString("link")) == 0 {
		if IS_NPC(ch) {
			return
		}
		if ch.Mindlink != nil {
			act(libc.CString("@CYou remove the link your mind had with @w$N.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_CHAR)
			act(libc.CString("@w$n@C removes the link $s mind had with yours.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_VICT)
			ch.Mindlink.Mindlink = nil
			ch.Mindlink = nil
			ch.Linker = 0
			return
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Link with the mind of who?\r\n"))
			return
		} else if vict == ch {
			send_to_char(ch, libc.CString("Oh that makes a lot of sense...\r\n"))
			return
		} else if IS_NPC(vict) {
			send_to_char(ch, libc.CString("You can't touch the mind of such a thing.\r\n"))
			return
		} else if vict.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't touch the mind of such an artificial being.\r\n"))
			return
		} else if GET_SKILL(vict, SKILL_TELEPATHY) != 0 {
			send_to_char(ch, libc.CString("Kinda pointless when you are both telepathic huh?\r\n"))
			return
		} else if vict.Mindlink != nil {
			send_to_char(ch, libc.CString("Someone else is already telepathically linked with them.\r\n"))
			return
		} else if GET_SKILL(ch, SKILL_TELEPATHY) < axion_dice(int(float64(vict.Aff_abils.Intel)*0.1)) {
			act(libc.CString("@R$n@r tried to link $s mind with yours, but you manage to force a break in the link!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$N@r manages to sense the intrusion and with $S intelligence push you out!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			return
		} else {
			act(libc.CString("@CYou link your mind with @w$N.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@w$n@C links $s mind with yours. You can speak your thoughts to $m with 'think'.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			send_to_char(vict, libc.CString("@wIf this is undesirable, Try: meditate break@n\r\n"))
			vict.Mindlink = ch
			ch.Mindlink = vict
			ch.Linker = 1
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("read")) == 0 {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Read the mind of who?\r\n"))
			return
		} else if vict == ch {
			send_to_char(ch, libc.CString("Oh that makes a lot of sense...\r\n"))
			return
		} else if vict.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't touch the mind of such an artificial being.\r\n"))
			return
		} else {
			if axion_dice(0) > GET_SKILL(ch, SKILL_TELEPATHY) {
				ch.Mana -= ch.Max_mana / 40
				act(libc.CString("@wYou attempt to read $N's@w mind, but fail to see it clearly.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				if rand_number(1, 15) >= 14 && !AFF_FLAGGED(ch, AFF_SHOCKED) {
					act(libc.CString("@MYour mind has been shocked!@n"), TRUE, ch, nil, nil, TO_CHAR)
					ch.Affected_by[int(AFF_SHOCKED/32)] |= 1 << (int(AFF_SHOCKED % 32))
				} else {
					improve_skill(ch, SKILL_TELEPATHY, 0)
				}
				return
			} else if GET_SKILL(vict, SKILL_TELEPATHY) >= GET_SKILL(ch, SKILL_TELEPATHY) && rand_number(1, 2) == 2 {
				ch.Mana -= ch.Max_mana / 40
				act(libc.CString("@wYou fail to read @c$N's@w mind and they seemed to have noticed the attempt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@w attempts to read your mind, but you resist and force $m out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				improve_skill(ch, SKILL_TELEPATHY, 0)
				return
			} else {
				send_to_char(ch, libc.CString("@wYou peer into their mind:\r\n"))
				ch.Mana -= ch.Max_mana / 40
				send_to_char(ch, libc.CString("@GName      @D: @W%s@n\r\n"), GET_NAME(vict))
				send_to_char(ch, libc.CString("@GRace      @D: @W%s@n\r\n"), pc_race_types[int(vict.Race)])
				send_to_char(ch, libc.CString("@GSensei    @D: @W%s@n\r\n"), pc_class_types[int(vict.Chclass)])
				send_to_char(ch, libc.CString("@GStr       @D: @W%d@n\r\n"), vict.Aff_abils.Str)
				send_to_char(ch, libc.CString("@GCon       @D: @W%d@n\r\n"), vict.Aff_abils.Con)
				send_to_char(ch, libc.CString("@GInt       @D: @W%d@n\r\n"), vict.Aff_abils.Intel)
				send_to_char(ch, libc.CString("@GWis       @D: @W%d@n\r\n"), vict.Aff_abils.Wis)
				send_to_char(ch, libc.CString("@GSpd       @D: @W%d@n\r\n"), vict.Aff_abils.Cha)
				send_to_char(ch, libc.CString("@GAgi       @D: @W%d@n\r\n"), vict.Aff_abils.Dex)
				send_to_char(ch, libc.CString("@GZenni     @D: @W%s@n\r\n"), add_commas(int64(vict.Gold)))
				send_to_char(ch, libc.CString("@GBank Zenni@D: @W%s@n\r\n"), add_commas(int64(vict.Bank_gold)))
				if vict.Alignment >= 1000 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wSaint         @n\r\n"))
				} else if vict.Alignment > 750 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wExtremely Good@n\r\n"))
				} else if vict.Alignment > 500 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wReally Good   @n\r\n"))
				} else if vict.Alignment > 250 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wGood          @n\r\n"))
				} else if vict.Alignment > 100 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wPretty Good   @n\r\n"))
				} else if vict.Alignment > 50 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wSorta Good    @n\r\n"))
				} else if vict.Alignment > -50 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wNeutral       @n\r\n"))
				} else if vict.Alignment > -100 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wSorta Evil    @n\r\n"))
				} else if vict.Alignment > -500 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wPretty Evil   @n\r\n"))
				} else if vict.Alignment >= -750 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wEvil          @n\r\n"))
				} else if vict.Alignment < -750 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wExtremely Evil@n\r\n"))
				} else if vict.Alignment <= -1000 {
					send_to_char(ch, libc.CString("@GAlignment @D: @wDevil         @n\r\n"))
				} else {
					send_to_char(ch, libc.CString("@GAlignment @D: @wUnknown       @n\r\n"))
				}
			}
		}
		improve_skill(ch, SKILL_TELEPATHY, 0)
		return
	} else {
		if ch.Mindlink != nil {
			vict = ch.Mindlink
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Send your thoughts to who?\r\n"))
			return
		}
		if vict == ch {
			send_to_char(ch, libc.CString("Oh that makes a lot of sense...\r\n"))
			return
		} else if vict.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't touch the mind of such an artificial being.\r\n"))
			return
		} else {
			if ch.Mindlink == nil {
				send_to_char(ch, libc.CString("@WYou tell @c%s@W telepathically, @w'@C%s@w'@n\r\n"), GET_NAME(vict), &arg2[0])
				send_to_char(vict, libc.CString("@c%s@W talks to you telepathically, @w'@C%s@w'@n\r\n"), GET_NAME(ch), &arg2[0])
				send_to_imm(libc.CString("@GTELEPATHY: @C%s@G telepaths @c%s, @W'@w%s@W'@n"), func() *byte {
					if ch.Admlevel > 0 {
						return GET_NAME(ch)
					}
					return GET_USER(ch)
				}(), func() *byte {
					if vict.Admlevel > 0 {
						return GET_NAME(vict)
					}
					return GET_USER(vict)
				}(), &arg2[0])
			} else {
				send_to_char(ch, libc.CString("@WYou tell @c%s@W telepathically, @w'@C%s@w'@n\r\n"), GET_NAME(vict), argument)
				send_to_char(vict, libc.CString("@c%s@W talks to you telepathically, @w'@C%s@w'@n\r\n"), GET_NAME(ch), argument)
				send_to_imm(libc.CString("@GTELEPATHY: @C%s@G telepaths @c%s, @W'@w%s@W'@n"), func() *byte {
					if ch.Admlevel > 0 {
						return GET_NAME(ch)
					}
					return GET_USER(ch)
				}(), func() *byte {
					if vict.Admlevel > 0 {
						return GET_NAME(vict)
					}
					return GET_USER(vict)
				}(), argument)
			}
			ch.Mana -= ch.Max_mana / 40
		}
		return
	}
}
func do_potential(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		boost int = 0
		vict  *char_data
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if know_skill(ch, SKILL_POTENTIAL) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who's potential do you want to release?\r\n"))
		send_to_char(ch, libc.CString("Potential Releases: %d\r\n"), ch.Boosts)
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("Why would you waste your time releasing their potential?\r\n"))
		return
	}
	if vict == ch {
		send_to_char(ch, libc.CString("You can't release your own potential.\r\n"))
		return
	}
	if ch.Boosts == 0 {
		send_to_char(ch, libc.CString("You have no potential releases to perform.\r\n"))
		return
	}
	if PLR_FLAGGED(vict, PLR_PR) {
		send_to_char(ch, libc.CString("Their potential has already been released\r\n"))
		return
	}
	if vict.Race == RACE_ANDROID {
		send_to_char(ch, libc.CString("They are a machine and have no potential to release.\r\n"))
		return
	}
	if vict.Majinize > 0 {
		send_to_char(ch, libc.CString("They are already majinized and have no potential to release.\r\n"))
		return
	}
	if vict.Race == RACE_MAJIN {
		send_to_char(ch, libc.CString("They have no potential to release...\r\n"))
		return
	} else {
		boost = GET_SKILL(ch, SKILL_POTENTIAL) / 2
		var mult float64 = 1
		if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 3
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 4
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 5
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 1 {
			mult = 2
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 2 {
			mult = 3
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 3.5
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS4) {
			mult = 4
		}
		vict.Act[int(PLR_PR/32)] |= bitvector_t(1 << (int(PLR_PR % 32)))
		vict.Max_hit += int64(float64((vict.Basepl/100)*int64(boost)) * mult)
		vict.Basepl += (vict.Basepl / 100) * int64(boost)
		if vict.Race == RACE_HALFBREED {
			vict.Max_mana += int64(float64((vict.Baseki/100)*int64(boost)) * mult)
			vict.Baseki += (vict.Baseki / 100) * int64(boost)
			vict.Max_move += int64(float64((vict.Basest/100)*int64(boost)) * mult)
			vict.Basest += (vict.Basest / 100) * int64(boost)
		}
		reveal_hiding(ch, 0)
		act(libc.CString("You place your hand on top of $N's head. After a moment of concentrating you release their hidden potential."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n places $s hand on top of your head. After a moment you feel a rush of power as your hidden potential is released!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n places $s hand on $N's head. After a moment a rush of power explodes off of $N's body!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		improve_skill(ch, SKILL_POTENTIAL, 0)
		improve_skill(ch, SKILL_POTENTIAL, 0)
		improve_skill(ch, SKILL_POTENTIAL, 0)
		improve_skill(ch, SKILL_POTENTIAL, 0)
		ch.Boosts -= 1
		return
	}
}
func do_majinize(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if ch.Race != RACE_MAJIN {
		send_to_char(ch, libc.CString("You are not a majin and can not majinize anyone.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who do you want to majinize?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if IS_NPC(vict) {
		send_to_char(ch, libc.CString("Why would you waste your time majinizing them?\r\n"))
		return
	}
	if vict == ch {
		send_to_char(ch, libc.CString("You can't majinize yourself.\r\n"))
		return
	}
	if PLR_FLAGGED(vict, PLR_PR) {
		send_to_char(ch, libc.CString("You can't majinize them their potential has been released!\r\n"))
		return
	}
	var alignmentTotal int = ch.Alignment - vict.Alignment
	if vict.Majinize > 0 && vict.Majinize != int(ch.Id) {
		send_to_char(ch, libc.CString("They are already majinized before by someone else.\r\n"))
		return
	} else if vict.Master != ch {
		send_to_char(ch, libc.CString("They must be following you in order for you to majinize them.\r\n"))
		return
	} else if alignmentTotal < -1500 || alignmentTotal > 1500 {
		send_to_char(ch, libc.CString("Their alignment is so opposed to your's that they resist your attempts to enslave them!\r\n"))
		return
	} else if vict.Max_hit > ch.Max_hit*4 {
		send_to_char(ch, libc.CString("Their powerlevel is so much higher than yours they resist your attempts to enslave them!\r\n"))
		return
	} else if vict.Majinize > 0 && vict.Majinize == int(ch.Id) {
		reveal_hiding(ch, 0)
		act(libc.CString("You remove $N's majinization, freeing them from your influence, but also weakening them."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n removes your majinization, freeing you from their influence, and weakening you!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n waves a hand at $N, and instantly the glowing M on $S forehead disappears!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Majinize = 0
		ch.Boosts += 1
		var mult float64 = 1
		if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 3
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 4
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 5
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 3.5
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS4) {
			mult = 4
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 4.5
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 1 {
			mult = 2
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 2 {
			mult = 3
		}
		if vict.Majinizer == 0 {
			vict.Majinizer = int64((float64(vict.Basepl) * 0.4) * mult)
		}
		vict.Max_hit -= int64(float64(vict.Majinizer) * mult)
		if vict.Majinizer == 0 {
			vict.Majinizer = int64(float64(vict.Basepl) * 0.4)
		}
		vict.Basepl -= vict.Majinizer
		return
	} else if ch.Boosts == 0 {
		send_to_char(ch, libc.CString("You are incapable of majinizing%s.\r\n"), func() string {
			if GET_LEVEL(ch) < 100 {
				return " right now"
			}
			return " anymore"
		}())
		if GET_LEVEL(ch) < 25 {
			send_to_char(ch, libc.CString("Your next available majinize will be at level 25\r\n"))
		} else if GET_LEVEL(ch) < 50 {
			send_to_char(ch, libc.CString("Your next available majinize will be at level 50\r\n"))
		} else if GET_LEVEL(ch) < 75 {
			send_to_char(ch, libc.CString("Your next available majinize will be at level 75\r\n"))
		} else if GET_LEVEL(ch) < 100 {
			send_to_char(ch, libc.CString("Your next available majinize will be at level 100\r\n"))
		}
		return
	} else {
		reveal_hiding(ch, 0)
		act(libc.CString("You focus your power into $N, influencing their mind and increasing their strength! After the struggle ends in $S mind a glowing purple M forms on $S forehead."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n focuses power into you, influencing your mind and increasing your strength! After the struggle in your mind ends a glowing purple M forms on your forehead."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n focuses power into $N, influencing their mind and increasing their strength! After the struggle ends in $S mind a glowing purple M forms on $S forehead."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Majinize = int(ch.Id)
		ch.Boosts -= 1
		var mult float64 = 1
		if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 3
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 4
		} else if vict.Race == RACE_TRUFFLE && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 5
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 3.5
		} else if vict.Race == RACE_BIO && PLR_FLAGGED(vict, PLR_TRANS4) {
			mult = 4
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS1) {
			mult = 2
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS2) {
			mult = 3
		} else if vict.Race == RACE_MAJIN && PLR_FLAGGED(vict, PLR_TRANS3) {
			mult = 4.5
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 1 {
			mult = 2
		} else if vict.Race == RACE_HOSHIJIN && vict.Starphase == 2 {
			mult = 3
		}
		vict.Majinizer = int64(float64(vict.Basepl) * 0.4)
		vict.Max_hit += int64((float64(vict.Basepl) * 0.4) * mult)
		vict.Basepl += int64(float64(vict.Basepl) * 0.4)
		return
	}
}
func do_spit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		cost int = 0
		vict *char_data
		af   affected_type
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if know_skill(ch, SKILL_SPIT) == 0 {
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You can't manage to spit in this fight!\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Yes but who do you want to petrify?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("That target isn't here.\r\n"))
		return
	}
	if can_kill(ch, vict, nil, 0) == 0 {
		return
	}
	if AFF_FLAGGED(vict, AFF_PARALYZE) {
		act(libc.CString("$N has already been turned to stone."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		return
	}
	if vict.Fighting != nil {
		send_to_char(vict, libc.CString("You can't manage to spit on them, they are moving around too much!\r\n"))
		return
	}
	cost = int((ch.Max_mana / int64(GET_SKILL(ch, SKILL_SPIT)/4)) + ch.Max_mana/100)
	if ch.Mana < int64(cost) {
		send_to_char(ch, libc.CString("You do not have enough ki to petrifiy with your spit!\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_SPIT) < axion_dice(0) {
		ch.Mana -= int64(cost)
		reveal_hiding(ch, 0)
		act(libc.CString("@WGathering spit you concentrate ki into a wicked loogie and let it loose, but it falls short of hitting @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W seems to focus ki before hawking a loogie at you! Fortunatly the loogie falls short.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W seems to focus ki before hawking a loogie at @c$N@W! Fortunatly for @c$N@W the loogie falls short.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		improve_skill(ch, SKILL_SPIT, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else if AFF_FLAGGED(vict, AFF_ZANZOKEN) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
		ch.Mana -= int64(cost)
		reveal_hiding(ch, 0)
		act(libc.CString("@C$N@c disappears, avoiding your spit before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@cYou disappear, avoiding @C$n's@c @rstone spit@c before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$N@c disappears, avoiding @C$n's@c @rstone spit@c before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		improve_skill(ch, SKILL_SPIT, 1)
		return
	} else {
		af.Type = SPELL_PARALYZE
		af.Duration = int16(rand_number(1, 2))
		af.Modifier = 0
		af.Location = APPLY_NONE
		af.Bitvector = AFF_PARALYZE
		affect_join(vict, &af, FALSE != 0, FALSE != 0, FALSE != 0, FALSE != 0)
		ch.Mana -= int64(cost)
		reveal_hiding(ch, 0)
		act(libc.CString("@WGathering spit you concentrate ki into a wicked loogie and let it loose, and it smacks into @c$N@W turning $M into stone!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W seems to focus ki before hawking a loogie at you! It manages to hit and you instantly turn to stone!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W seems to focus ki before hawking a loogie at @c$N@W! It manages to hit and $E instantly turns to stone!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		improve_skill(ch, SKILL_SPIT, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
}
func boost_obj(obj *obj_data, ch *char_data, type_ int) {
	if obj == nil || ch == nil {
		return
	}
	var boost int = 0
	if GET_LEVEL(ch) >= 100 {
		boost = 100
	} else if GET_LEVEL(ch) >= 90 {
		boost = 90
	} else if GET_LEVEL(ch) >= 80 {
		boost = 80
	} else if GET_LEVEL(ch) >= 70 {
		boost = 70
	} else if GET_LEVEL(ch) >= 60 {
		boost = 60
	} else if GET_LEVEL(ch) >= 50 {
		boost = 50
	} else if GET_LEVEL(ch) >= 40 {
		boost = 40
	} else if GET_LEVEL(ch) >= 30 {
		boost = 30
	}
	switch type_ {
	case 0:
		if boost != 0 {
			obj.Level = boost
			obj.Affected[0].Location = 17
			obj.Affected[0].Modifier += boost * GET_LEVEL(ch)
			if GET_OBJ_VNUM(obj) == 91 {
				obj.Affected[1].Location = 1
				obj.Affected[1].Modifier = boost / 20
			} else {
				obj.Affected[1].Location = 3
				obj.Affected[1].Modifier = boost / 20
			}
		}
	case 1:
		switch boost {
		case 30:
			obj.Extra_flags[int(ITEM_WEAPLVL2/32)] |= bitvector_t(1 << (int(ITEM_WEAPLVL2 % 32)))
		case 40:
			fallthrough
		case 50:
			obj.Extra_flags[int(ITEM_WEAPLVL3/32)] |= bitvector_t(1 << (int(ITEM_WEAPLVL3 % 32)))
		case 60:
			fallthrough
		case 70:
			fallthrough
		case 80:
			fallthrough
		case 90:
			obj.Extra_flags[int(ITEM_WEAPLVL4/32)] |= bitvector_t(1 << (int(ITEM_WEAPLVL4 % 32)))
		case 100:
			obj.Extra_flags[int(ITEM_WEAPLVL5/32)] |= bitvector_t(1 << (int(ITEM_WEAPLVL5 % 32)))
		default:
			obj.Extra_flags[int(ITEM_WEAPLVL1/32)] |= bitvector_t(1 << (int(ITEM_WEAPLVL1 % 32)))
		}
		if boost != 0 {
			obj.Level = boost
			obj.Affected[0].Location = 1
			obj.Affected[0].Modifier = boost / 20
		}
	}
}
func do_form(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		skill    int     = 0
		senzu    int     = FALSE
		bag      int     = FALSE
		light    int     = FALSE
		sword    int     = FALSE
		mattress int     = FALSE
		gi       int     = FALSE
		pants    int     = FALSE
		kachin   int     = FALSE
		boost    int     = FALSE
		shuriken int     = FALSE
		clothes  int     = FALSE
		wrist    int     = FALSE
		boots    int     = FALSE
		level    int     = 0
		discount float64 = 1.0
		cost     int64   = 0
		obj      *obj_data
		arg      [2048]byte
		arg2     [2048]byte
		arg3     [2048]byte
		clam     [2048]byte
	)
	half_chop(argument, &arg[0], &clam[0])
	half_chop(&clam[0], &arg2[0], &arg3[0])
	if know_skill(ch, SKILL_CREATE) == 0 {
		return
	}
	if ch.Con_cooldown > 0 {
		send_to_char(ch, libc.CString("You must wait a short period before concentrating again.\r\n"))
		return
	}
	skill = GET_SKILL(ch, SKILL_CREATE)
	if skill >= 100 {
		boost = TRUE
	}
	if skill >= 90 {
		kachin = TRUE
	}
	if skill >= 80 {
		senzu = TRUE
	}
	if skill >= 70 {
		shuriken = TRUE
	}
	if skill >= 60 {
		clothes = TRUE
	}
	if skill >= 50 {
		sword = TRUE
		gi = TRUE
		pants = TRUE
		wrist = TRUE
		boots = TRUE
	}
	if skill >= 40 {
		mattress = TRUE
	}
	if skill >= 30 {
		bag = TRUE
	}
	if skill >= 20 {
		light = TRUE
	}
	if GET_SKILL(ch, SKILL_CONCENTRATION) != 0 {
		if GET_SKILL(ch, SKILL_CONCENTRATION) >= 100 {
			discount = 0.5
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 90 {
			discount = 0.6
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 80 {
			discount = 0.65
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 70 {
			discount = 0.7
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 60 {
			discount = 0.75
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 50 {
			discount = 0.8
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 40 {
			discount = 0.85
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 30 {
			discount = 0.9
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) >= 20 {
			discount = 0.95
		}
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What do you want to create?\r\n@GCreation @WMenu@n\r\n@D---------------@n\r\n@wcreate food\r\ncreate water\r\n%s%s%s%s%s%s%s%s%s%s%s%s%s\r\n"), func() string {
			if light != 0 {
				return "create light\r\n"
			}
			return ""
		}(), func() string {
			if bag != 0 {
				return "create bag\r\n"
			}
			return ""
		}(), func() string {
			if mattress != 0 {
				return "create mattress\r\n"
			}
			return ""
		}(), func() string {
			if sword != 0 {
				return "create weapon (sword | club | dagger | spear | gun )\r\n"
			}
			return ""
		}(), func() string {
			if pants != 0 {
				return "create pants\r\n"
			}
			return ""
		}(), func() string {
			if gi != 0 {
				return "create gi\r\n"
			}
			return ""
		}(), func() string {
			if wrist != 0 {
				return "create wristband\r\n"
			}
			return ""
		}(), func() string {
			if boots != 0 {
				return "create boots\r\n"
			}
			return ""
		}(), func() string {
			if clothes != 0 {
				return "create clothesbeam (target)\r\n"
			}
			return ""
		}(), func() string {
			if shuriken != 0 {
				return "create shuriken\r\n"
			}
			return ""
		}(), func() string {
			if senzu != 0 {
				return "create senzu\r\n"
			}
			return ""
		}(), func() string {
			if kachin != 0 {
				return "create kachin\r\n"
			}
			return ""
		}(), func() string {
			if boost != 0 {
				return "create elixir\r\n"
			}
			return ""
		}())
		return
	}
	reveal_hiding(ch, 0)
	if C.strcmp(&arg[0], libc.CString("food")) == 0 {
		cost = ch.Max_mana / int64(skill/2)
		cost *= int64(discount)
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			if arg2[0] == 0 {
				send_to_char(ch, libc.CString("Making lowest quality version of object. To make a higher quality use, Syntax: create (type) (mid | high | highest)\r\n"))
				send_to_char(ch, libc.CString("If you are capable you will make it. If not you will make a low quality version.\r\n"))
			} else if arg2[0] != 0 {
				if C.strcasecmp(&arg2[0], libc.CString("highest")) == 0 && skill >= 100 {
					level = 4
				} else if C.strcasecmp(&arg2[0], libc.CString("high")) == 0 && skill >= 75 {
					level = 3
				} else if C.strcasecmp(&arg2[0], libc.CString("mid")) == 0 && skill >= 50 {
					level = 2
				} else {
					level = 1
				}
			}
			if level == 4 {
				obj = read_object(1512, VIRTUAL)
				add_unique_id(obj)
			} else if level == 3 {
				obj = read_object(1511, VIRTUAL)
				add_unique_id(obj)
			} else if level == 2 {
				obj = read_object(1510, VIRTUAL)
				add_unique_id(obj)
			} else {
				obj = read_object(70, VIRTUAL)
				add_unique_id(obj)
			}
			obj_to_char(obj, ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("water")) == 0 {
		cost = ch.Max_mana / int64(skill*2)
		cost *= int64(discount)
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			if arg2[0] == 0 {
				send_to_char(ch, libc.CString("Making lowest quality version of object. To make a higher quality use, Syntax: create (type) (mid | high | highest)\r\n"))
				send_to_char(ch, libc.CString("If you are capable you will make it. If not you will make a low quality version.\r\n"))
			} else if arg2[0] != 0 {
				if C.strcasecmp(&arg2[0], libc.CString("highest")) == 0 && skill >= 100 {
					level = 4
				} else if C.strcasecmp(&arg2[0], libc.CString("high")) == 0 && skill >= 75 {
					level = 3
				} else if C.strcasecmp(&arg2[0], libc.CString("mid")) == 0 && skill >= 50 {
					level = 2
				} else {
					level = 1
				}
			}
			if level == 4 {
				obj = read_object(1515, VIRTUAL)
				add_unique_id(obj)
			} else if level == 3 {
				obj = read_object(1514, VIRTUAL)
				add_unique_id(obj)
			} else if level == 2 {
				obj = read_object(1513, VIRTUAL)
				add_unique_id(obj)
			} else {
				obj = read_object(71, VIRTUAL)
				add_unique_id(obj)
			}
			obj_to_char(obj, ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("bag")) == 0 {
		cost = ch.Max_mana / int64(skill*2)
		cost *= int64(discount)
		if bag == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(319, VIRTUAL)
			add_unique_id(obj)
			obj_to_char(obj, ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("mattress")) == 0 {
		cost = ch.Max_mana / int64(skill)
		cost *= int64(discount)
		if mattress == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(16, VIRTUAL)
			add_unique_id(obj)
			obj_to_char(obj, ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("weapon")) == 0 {
		cost = ch.Max_mana / 5
		cost *= int64(discount)
		if sword == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			if arg2[0] == 0 {
				send_to_char(ch, libc.CString("What type of weapon?\r\nSyntax: create weapon (sword | club | spear | dagger | gun)\r\n"))
				return
			}
			if arg3[0] == 0 {
				send_to_char(ch, libc.CString("Making lowest quality version of object. To make a higher quality use, Syntax: create (type) (mid | high | higher | highest)\r\n"))
				send_to_char(ch, libc.CString("If you are capable you will make it. If not you will make a low quality version.\r\n"))
			} else if arg3[0] != 0 {
				if C.strcasecmp(&arg3[0], libc.CString("highest")) == 0 && skill >= 100 {
					level = 5
				} else if C.strcasecmp(&arg3[0], libc.CString("higher")) == 0 && skill >= 75 {
					level = 4
				} else if C.strcasecmp(&arg3[0], libc.CString("high")) == 0 && skill >= 50 {
					level = 3
				} else if C.strcasecmp(&arg3[0], libc.CString("mid")) == 0 && skill >= 30 {
					level = 2
				} else {
					level = 1
				}
			}
			if C.strcasecmp(&arg2[0], libc.CString("sword")) == 0 {
				if level == 5 {
					obj = read_object(1519, VIRTUAL)
					add_unique_id(obj)
				} else if level == 4 {
					obj = read_object(1518, VIRTUAL)
					add_unique_id(obj)
				} else if level == 3 {
					obj = read_object(1517, VIRTUAL)
					add_unique_id(obj)
				} else if level == 2 {
					obj = read_object(1516, VIRTUAL)
					add_unique_id(obj)
				} else {
					obj = read_object(90, VIRTUAL)
					add_unique_id(obj)
				}
			} else if C.strcasecmp(&arg2[0], libc.CString("dagger")) == 0 {
				if level == 5 {
					obj = read_object(1540, VIRTUAL)
					add_unique_id(obj)
				} else if level == 4 {
					obj = read_object(1539, VIRTUAL)
					add_unique_id(obj)
				} else if level == 3 {
					obj = read_object(1538, VIRTUAL)
					add_unique_id(obj)
				} else if level == 2 {
					obj = read_object(1537, VIRTUAL)
					add_unique_id(obj)
				} else {
					obj = read_object(1536, VIRTUAL)
					add_unique_id(obj)
				}
			} else if C.strcasecmp(&arg2[0], libc.CString("club")) == 0 {
				if level == 5 {
					obj = read_object(1545, VIRTUAL)
					add_unique_id(obj)
				} else if level == 4 {
					obj = read_object(1544, VIRTUAL)
					add_unique_id(obj)
				} else if level == 3 {
					obj = read_object(1543, VIRTUAL)
					add_unique_id(obj)
				} else if level == 2 {
					obj = read_object(1542, VIRTUAL)
					add_unique_id(obj)
				} else {
					obj = read_object(1541, VIRTUAL)
					add_unique_id(obj)
				}
			} else if C.strcasecmp(&arg2[0], libc.CString("spear")) == 0 {
				if level == 5 {
					obj = read_object(1550, VIRTUAL)
					add_unique_id(obj)
				} else if level == 4 {
					obj = read_object(1549, VIRTUAL)
					add_unique_id(obj)
				} else if level == 3 {
					obj = read_object(1548, VIRTUAL)
					add_unique_id(obj)
				} else if level == 2 {
					obj = read_object(1547, VIRTUAL)
					add_unique_id(obj)
				} else {
					obj = read_object(1546, VIRTUAL)
					add_unique_id(obj)
				}
			} else if C.strcasecmp(&arg2[0], libc.CString("gun")) == 0 {
				if level == 5 {
					obj = read_object(1555, VIRTUAL)
					add_unique_id(obj)
				} else if level == 4 {
					obj = read_object(1554, VIRTUAL)
					add_unique_id(obj)
				} else if level == 3 {
					obj = read_object(1553, VIRTUAL)
					add_unique_id(obj)
				} else if level == 2 {
					obj = read_object(1552, VIRTUAL)
					add_unique_id(obj)
				} else {
					obj = read_object(1551, VIRTUAL)
					add_unique_id(obj)
				}
			} else {
				send_to_char(ch, libc.CString("What type of weapon?\r\nSyntax: create weapon (sword | club | spear | dagger | gun)\r\n"))
				return
			}
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("clothesbeam")) == 0 {
		cost = ch.Max_mana / 2
		cost *= int64(discount)
		if clothes == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		}
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("Who do you want to hit with clothesbeam?\r\nSyntax: create clothesbeam (target)\r\n"))
			return
		}
		var vict *char_data = nil
		if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Clothesbeam who?\r\nSyntax: create clothesbeam (target)\r\n"))
			return
		}
		if vict.Master != ch {
			send_to_char(ch, libc.CString("They must be following you first.\r\n"))
			return
		} else {
			obj = read_object(92, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, vict)
			obj.Size = get_size(vict)
			obj = read_object(91, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, vict)
			obj.Size = get_size(vict)
			obj = read_object(1528, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, vict)
			obj.Size = get_size(vict)
			obj = read_object(1528, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, vict)
			obj.Size = get_size(vict)
			obj = read_object(1532, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, vict)
			obj.Size = get_size(vict)
			do_wear(vict, libc.CString("all"), 0, 0)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("gi")) == 0 {
		cost = ch.Max_mana / 5
		cost *= int64(discount)
		if gi == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(92, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("shuriken")) == 0 {
		cost = ch.Max_mana / 4
		cost *= int64(discount)
		if shuriken == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(0x4A6D, VIRTUAL)
			add_unique_id(obj)
			obj_to_char(obj, ch)
			obj.Extra_flags[int(ITEM_NORENT/32)] |= bitvector_t(1 << (int(ITEM_NORENT % 32)))
			obj.Extra_flags[int(ITEM_NOSELL/32)] |= bitvector_t(1 << (int(ITEM_NOSELL % 32)))
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("pants")) == 0 {
		cost = ch.Max_mana / 5
		cost *= int64(discount)
		if pants == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(91, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("wristband")) == 0 {
		cost = ch.Max_mana / 5
		cost *= int64(discount)
		if wrist == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(1528, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("boots")) == 0 {
		cost = ch.Max_mana / 5
		cost *= int64(discount)
		if boots == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(1532, VIRTUAL)
			add_unique_id(obj)
			boost_obj(obj, ch, 0)
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("light")) == 0 {
		cost = ch.Max_mana / int64(skill*2)
		cost *= int64(discount)
		if light == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(72, VIRTUAL)
			add_unique_id(obj)
			obj_to_char(obj, ch)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("kachin")) == 0 {
		cost = ch.Max_mana - 1
		cost *= int64(discount)
		if kachin == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		} else {
			obj = read_object(87, VIRTUAL)
			add_unique_id(obj)
			obj_to_room(obj, ch.In_room)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("elixir")) == 0 {
		cost = ch.Max_mana - 1
		cost *= int64(discount)
		if boost == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s\r\n"), &arg[0])
			return
		}
		if ch.Hit < ch.Max_hit {
			send_to_char(ch, libc.CString("You need to be at full powerlevel to create %s\r\n"), &arg[0])
			return
		} else if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 10 {
			send_to_char(ch, libc.CString("You do not have enough PS to create %s, you need at least 10.\r\n"), &arg[0])
			return
		} else {
			obj = read_object(86, VIRTUAL)
			add_unique_id(obj)
			obj_to_room(obj, ch.In_room)
			obj.Size = get_size(ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			ch.Hit = 1
			ch.Player_specials.Class_skill_points[ch.Chclass] -= 10
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("senzu")) == 0 {
		cost = ch.Max_mana
		var cost2 int64 = gear_pl(ch) - 1
		if senzu == FALSE {
			send_to_char(ch, libc.CString("What do you want to create?\r\n"))
			return
		}
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to create %s, you need full ki.\r\n"), &arg[0])
			return
		} else if ch.Hit <= cost2 {
			send_to_char(ch, libc.CString("You do not have enough powerlevel to create %s, you need to be at full.\r\n"), &arg[0])
			return
		} else if ch.Move < ch.Max_move {
			send_to_char(ch, libc.CString("You do not have enough stamina to create %s, you need to be at full.\r\n"), &arg[0])
			return
		} else if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 50 {
			send_to_char(ch, libc.CString("You do not have enough PS to create %s, you need at least 50.\r\n"), &arg[0])
			return
		} else {
			obj = read_object(1, VIRTUAL)
			add_unique_id(obj)
			obj_to_char(obj, ch)
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You hold out your hand and create $p out of your ki!"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n holds out $s hand and creates $p out of thin air!"), TRUE, ch, obj, nil, TO_ROOM)
			ch.Mana -= cost
			ch.Hit -= cost2
			ch.Move = 1
			ch.Player_specials.Class_skill_points[ch.Chclass] -= 50
			return
		}
	} else {
		send_to_char(ch, libc.CString("Create what?\r\n"))
		return
	}
}
func do_recharge(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) || ch.Race != RACE_ANDROID {
		send_to_char(ch, libc.CString("Only androids can use recharge\r\n"))
		return
	}
	if ch.Con_cooldown > 0 {
		send_to_char(ch, libc.CString("You must wait a short period before your nanites can convert your ki.\r\n"))
		return
	}
	if !PLR_FLAGGED(ch, PLR_REPAIR) {
		send_to_char(ch, libc.CString("You are not a repair model android.\r\n"))
		return
	} else {
		var cost int64 = 0
		cost = ch.Max_move / 20
		if ch.Mana < cost {
			send_to_char(ch, libc.CString("You do not have enough ki to recharge your stamina.\r\n"))
			return
		} else if ch.Move >= ch.Max_move {
			send_to_char(ch, libc.CString("Your energy reserves are already full.\r\n"))
			return
		} else {
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You focus your ki into your energy reserves, recharging them some."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n stops and glows green briefly."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Mana -= cost
			if ch.Move+cost*2 < ch.Max_move {
				ch.Move += cost * 2
			} else {
				ch.Move = ch.Max_move
				send_to_char(ch, libc.CString("You are fully recharged now.\r\n"))
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
	}
}
func do_srepair(ch *char_data, argument *byte, cmd int, subcmd int) {
	var i int
	if ch.Race != RACE_ANDROID {
		send_to_char(ch, libc.CString("Only androids can use repair, maybe you want 'fix' instead?\r\n"))
		return
	}
	if ch.Con_cooldown > 0 {
		send_to_char(ch, libc.CString("You must wait a short period before your nanites can repair you.\r\n"))
		return
	}
	if !IS_NPC(ch) && !PLR_FLAGGED(ch, PLR_REPAIR) {
		send_to_char(ch, libc.CString("You are not a repair model android.\r\n"))
		return
	} else {
		var (
			cost int64 = 0
			heal int64 = 0
		)
		cost = ch.Max_hit / 40
		if ch.Move < cost {
			send_to_char(ch, libc.CString("You do not have enough stamina to repair yourself.\r\n"))
			return
		} else if ch.Hit >= gear_pl(ch) {
			send_to_char(ch, libc.CString("You are already at full functionality and do not require repairs.\r\n"))
			return
		} else {
			reveal_hiding(ch, 0)
			ch.Con_cooldown = 10
			act(libc.CString("You repair some of your outer casings and internal systems, with the small nano-robots contained in your body."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n stops a moment as small glowing particles move across $s body."), TRUE, ch, nil, nil, TO_ROOM)
			var repaired int = FALSE
			if !IS_NPC(ch) {
				for i = 0; i < NUM_WEARS; i++ {
					if (ch.Equipment[i]) != nil {
						if ((ch.Equipment[i]).Value[VAL_ALL_HEALTH]) < 100 {
							(ch.Equipment[i]).Value[VAL_ALL_HEALTH] += 20
							if ((ch.Equipment[i]).Value[VAL_ALL_HEALTH]) > 100 {
								(ch.Equipment[i]).Value[VAL_ALL_HEALTH] = 100
							}
							if OBJ_FLAGGED(ch.Equipment[i], ITEM_BROKEN) {
								(ch.Equipment[i]).Extra_flags[int(ITEM_BROKEN/32)] &= bitvector_t(^(1 << (int(ITEM_BROKEN % 32))))
							}
							repaired = TRUE
						}
					}
				}
			}
			if repaired == TRUE {
				send_to_char(ch, libc.CString("@GYour nano-robots also repair all of your equipment a little bit.@n\r\n"))
			}
			ch.Move -= cost
			heal = cost * 2
			if (ch.Bonuses[BONUS_HEALER]) > 0 {
				heal += int64(float64(heal) * 0.25)
			}
			if ch.Hit+heal < gear_pl(ch) {
				ch.Hit += heal
			} else {
				ch.Hit = gear_pl(ch)
				send_to_char(ch, libc.CString("You are fully repaired now.\r\n"))
			}
			if !IS_NPC(ch) && rand_number(1, 3) == 2 && ch.Mana < ch.Max_mana {
				send_to_char(ch, libc.CString("@GThe repairs have managed to relink power reserves and boost your current energy level.@n\r\n"))
				ch.Mana += cost
				if ch.Mana > ch.Max_mana {
					ch.Mana = ch.Max_mana
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
	}
}
func do_upgrade(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		count int = 0
		bonus int = 0
		cost  int = 0
		arg   [2048]byte
		arg2  [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if IS_NPC(ch) || ch.Race != RACE_ANDROID {
		send_to_char(ch, libc.CString("You are not an android!\r\n"))
		return
	}
	if arg[0] == 0 {
		if !PLR_FLAGGED(ch, PLR_ABSORB) {
			send_to_char(ch, libc.CString("@c--------@D[@rUpgrade Menu@D]@c--------\r\n@cUpgrade @RPowerlevel@D: @Y75 @WPoints\r\n@cUpgrade @CKi        @D: @Y40 @WPoints\r\n@cUpgrade @GStamina   @D: @Y50 @WPoints\r\n@D            -----------\r\n"))
		}
		send_to_char(ch, libc.CString("@cAugment @RPowerlevel\r\n@cAugment @CKi\r\n@cAugment @GStamina\r\n@WCurrent Upgrade Points @D[@y%s@D]@n\r\n"), add_commas(int64(ch.Upgrade)))
		return
	}
	if C.strcasecmp(libc.CString("augment"), &arg[0]) == 0 {
		var (
			obj  *obj_data = nil
			gain int64     = 0
		)
		if GET_LEVEL(ch) < 80 {
			send_to_char(ch, libc.CString("You need to be at least level 80 to use these kits.\r\n"))
			return
		}
		if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, libc.CString("Augmentation"), nil, ch.Carrying)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have a Circuit Augmentation Kit.\r\n"))
			return
		} else {
			switch GET_LEVEL(ch) {
			case 80:
				gain = int64(float64(ch.Max_hit) * 0.005)
			case 81:
				fallthrough
			case 82:
				fallthrough
			case 83:
				fallthrough
			case 84:
				fallthrough
			case 85:
				fallthrough
			case 86:
				fallthrough
			case 87:
				fallthrough
			case 88:
				fallthrough
			case 89:
				fallthrough
			case 90:
				gain = int64(float64(ch.Max_hit) * 0.005)
			case 91:
				fallthrough
			case 92:
				fallthrough
			case 93:
				fallthrough
			case 94:
				fallthrough
			case 95:
				fallthrough
			case 96:
				fallthrough
			case 97:
				fallthrough
			case 98:
				fallthrough
			case 99:
				gain = int64(float64(ch.Max_hit) * 0.005)
			case 100:
				gain = int64(float64(ch.Max_hit) * 0.005)
				if gain > 10000000 {
					gain = 10000000
				}
			}
			if C.strcasecmp(libc.CString("powerlevel"), &arg2[0]) == 0 {
				obj_from_char(obj)
				extract_obj(obj)
				act(libc.CString("@WYou install the circuits and upgrade your maximum powerlevel.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W installs some circuits and upgrades $s systems.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Max_hit += gain
				ch.Basepl += gain
				send_to_char(ch, libc.CString("@gGain @D[@G+%s@D]\r\n"), add_commas(gain))
				return
			} else if C.strcasecmp(libc.CString("ki"), &arg2[0]) == 0 {
				obj_from_char(obj)
				extract_obj(obj)
				act(libc.CString("@WYou install the circuits and upgrade your maximum ki.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W installs some circuits and upgrades $s systems.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Max_mana += gain
				ch.Baseki += gain
				send_to_char(ch, libc.CString("@gGain @D[@G+%s@D]\r\n"), add_commas(gain))
				return
			} else if C.strcasecmp(libc.CString("stamina"), &arg2[0]) == 0 {
				obj_from_char(obj)
				extract_obj(obj)
				act(libc.CString("@WYou install the circuits and upgrade your maximum stamina.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W installs some circuits and upgrades $s systems.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Max_move += gain
				ch.Basest += gain
				send_to_char(ch, libc.CString("@gGain @D[@G+%s@D]\r\n"), add_commas(gain))
				return
			} else {
				send_to_char(ch, libc.CString("What do you want to augment? Powerlevel, ki, or stamina?\r\n"))
				return
			}
		}
	}
	if PLR_FLAGGED(ch, PLR_ABSORB) {
		send_to_char(ch, libc.CString("You are an absorb model and can only upgrade with augmentation kits.\r\n"))
		return
	}
	if !soft_cap(ch, 0) {
		send_to_char(ch, libc.CString("@mYou are unable to spend anymore UGP right now (Softcap)@n\r\n"))
		return
	}
	if arg2[0] == 0 && (C.strcasecmp(libc.CString("powerlevel"), &arg[0]) == 0 || C.strcasecmp(libc.CString("ki"), &arg[0]) == 0 || C.strcasecmp(libc.CString("stamina"), &arg[0]) == 0) {
		send_to_char(ch, libc.CString("How many times do you want to increase %s?"), &arg[0])
		return
	}
	if libc.Atoi(libc.GoString(&arg2[0])) <= 0 && (C.strcasecmp(libc.CString("powerlevel"), &arg[0]) == 0 || C.strcasecmp(libc.CString("ki"), &arg[0]) == 0 || C.strcasecmp(libc.CString("stamina"), &arg[0]) == 0) {
		send_to_char(ch, libc.CString("It needs to be between 1-1000\r\n"))
		return
	}
	if libc.Atoi(libc.GoString(&arg2[0])) > 1000 && (C.strcasecmp(libc.CString("powerlevel"), &arg[0]) == 0 || C.strcasecmp(libc.CString("ki"), &arg[0]) == 0 || C.strcasecmp(libc.CString("stamina"), &arg[0]) == 0) {
		send_to_char(ch, libc.CString("It needs to be between 1-1000\r\n"))
		return
	}
	if C.strcasecmp(libc.CString("powerlevel"), &arg[0]) == 0 {
		count = libc.Atoi(libc.GoString(&arg2[0]))
		for count > 0 {
			if GET_LEVEL(ch) >= 90 {
				bonus += GET_LEVEL(ch) * 5000
			} else if GET_LEVEL(ch) >= 80 {
				bonus += GET_LEVEL(ch) * 2500
			} else if GET_LEVEL(ch) >= 70 {
				bonus += GET_LEVEL(ch) * 2000
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1300
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1200
			} else if GET_LEVEL(ch) >= 50 {
				bonus += GET_LEVEL(ch) * 500
			} else if GET_LEVEL(ch) >= 25 {
				bonus += GET_LEVEL(ch) * 250
			} else {
				bonus += GET_LEVEL(ch) * 150
			}
			cost += 75
			count--
		}
		if cost > ch.Upgrade {
			send_to_char(ch, libc.CString("You need %s upgrade points, and only have %s.\r\n"), add_commas(int64(cost)), add_commas(int64(ch.Upgrade)))
			return
		} else if !soft_cap(ch, int64(bonus)) {
			send_to_char(ch, libc.CString("@mYou can't spend that much UGP on it as it will go over your softcap.@n\r\n"))
			return
		} else {
			ch.Upgrade -= cost
			send_to_char(ch, libc.CString("You upgrade your system and gain %s %s!"), add_commas(int64(bonus)), &arg[0])
			ch.Max_hit += int64(bonus)
			ch.Basepl += int64(bonus)
		}
	} else if C.strcasecmp(libc.CString("ki"), &arg[0]) == 0 {
		count = libc.Atoi(libc.GoString(&arg2[0]))
		for count > 0 {
			if GET_LEVEL(ch) >= 90 {
				bonus += GET_LEVEL(ch) * 3650
			} else if GET_LEVEL(ch) >= 80 {
				bonus += GET_LEVEL(ch) * 2450
			} else if GET_LEVEL(ch) >= 70 {
				bonus += GET_LEVEL(ch) * 1800
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1250
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1150
			} else if GET_LEVEL(ch) >= 50 {
				bonus += GET_LEVEL(ch) * 400
			} else if GET_LEVEL(ch) >= 25 {
				bonus += GET_LEVEL(ch) * 200
			} else {
				bonus += GET_LEVEL(ch) * 120
			}
			cost += 40
			count--
		}
		if cost > ch.Upgrade {
			send_to_char(ch, libc.CString("You need %s upgrade points, and only have %s.\r\n"), add_commas(int64(cost)), add_commas(int64(ch.Upgrade)))
			return
		} else if !soft_cap(ch, int64(bonus)) {
			send_to_char(ch, libc.CString("@mYou can't spend that much UGP on it as it will go over your softcap.@n\r\n"))
			return
		} else {
			ch.Upgrade -= cost
			send_to_char(ch, libc.CString("You upgrade your system and gain %s %s!"), add_commas(int64(bonus)), &arg[0])
			ch.Max_mana += int64(bonus)
			ch.Baseki += int64(bonus)
		}
	} else if C.strcasecmp(libc.CString("stamina"), &arg[0]) == 0 {
		count = libc.Atoi(libc.GoString(&arg2[0]))
		for count > 0 {
			if GET_LEVEL(ch) >= 90 {
				bonus += GET_LEVEL(ch) * 3650
			} else if GET_LEVEL(ch) >= 80 {
				bonus += GET_LEVEL(ch) * 2450
			} else if GET_LEVEL(ch) >= 70 {
				bonus += GET_LEVEL(ch) * 1800
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1250
			} else if GET_LEVEL(ch) >= 60 {
				bonus += GET_LEVEL(ch) * 1150
			} else if GET_LEVEL(ch) >= 50 {
				bonus += GET_LEVEL(ch) * 500
			} else if GET_LEVEL(ch) >= 25 {
				bonus += GET_LEVEL(ch) * 200
			} else {
				bonus += GET_LEVEL(ch) * 120
			}
			cost += 50
			count--
		}
		if cost > ch.Upgrade {
			send_to_char(ch, libc.CString("You need %s upgrade points, and only have %s.\r\n"), add_commas(int64(cost)), add_commas(int64(ch.Upgrade)))
			return
		} else if !soft_cap(ch, int64(bonus)) {
			send_to_char(ch, libc.CString("@mYou can't spend that much UGP on it as it will go over your softcap.@n\r\n"))
			return
		} else {
			ch.Upgrade -= cost
			send_to_char(ch, libc.CString("You upgrade your system and gain %s %s!"), add_commas(int64(bonus)), &arg[0])
			ch.Max_move += int64(bonus)
			ch.Basest += int64(bonus)
		}
	} else {
		send_to_char(ch, libc.CString("That is not a valid upgrade option.\r\n"))
		return
	}
}
func do_ingest(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Race == RACE_MAJIN {
		var (
			vict *char_data
			arg  [2048]byte
		)
		one_argument(argument, &arg[0])
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("Who do you want to ingest?\r\n"))
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Ingest who?\r\n"))
			return
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if vict.Absorbby != nil {
			send_to_char(ch, libc.CString("%s is already absorbing from them!"), GET_NAME(vict.Absorbby))
			return
		}
		if ch.Absorbs > 3 {
			send_to_char(ch, libc.CString("You already have already ingested 4 people.\r\n"))
			return
		}
		if GET_LEVEL(ch) < 25 {
			send_to_char(ch, libc.CString("You can't ingest yet.\r\n"))
			return
		}
		if GET_LEVEL(ch) < 100 && GET_LEVEL(ch) >= 75 && ch.Absorbs == 3 {
			send_to_char(ch, libc.CString("You already have ingested as much as you can. You'll have to get more experienced.\r\n"))
			return
		}
		if GET_LEVEL(ch) < 75 && GET_LEVEL(ch) >= 50 && ch.Absorbs == 2 {
			send_to_char(ch, libc.CString("You already have ingested as much as you can. You'll have to get more experienced.\r\n"))
			return
		}
		if GET_LEVEL(ch) < 50 && GET_LEVEL(ch) >= 25 && ch.Absorbs == 1 {
			send_to_char(ch, libc.CString("You already have ingested as much as you can. You'll have to get more experienced.\r\n"))
			return
		}
		if vict.Max_hit >= ch.Basepl*3 {
			send_to_char(ch, libc.CString("You are too weak to ingest them into your body!\r\n"))
			return
		}
		if AFF_FLAGGED(vict, AFF_SANCTUARY) {
			send_to_char(ch, libc.CString("You can't ingest them, they have a barrier!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		if AFF_FLAGGED(vict, AFF_ZANZOKEN) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
			act(libc.CString("@C$N@c disappears, avoiding your attempted ingestion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@cYou disappear, avoiding @C$n's@c attempted @ringestion@c before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$N@c disappears, avoiding @C$n's@c attempted @ringestion@c before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		}
		if GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(ch)+rand_number(1, 5) {
			act(libc.CString("@WYou fling a piece of goo at @c$N@W, and try to ingest $M! $E manages to avoid your blob of goo though!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W flings a piece of goo at you, you manage to avoid it though!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w flings a piece of goo at @c$N@W, but the goo misses $M@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else {
			act(libc.CString("@WYou flings a piece of goo at @c$N@W! The goo engulfs $M and then returns to your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W flings a piece of goo at you! The goo engulfs your body and then returns to @C$n@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w flings a piece of goo at @c$N@W! The goo engulfs $M and then return to @C$n@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Absorbs += 1
			var pl int64 = vict.Basepl / 6
			var stam int64 = vict.Basest / 6
			var ki int64 = vict.Baseki / 6
			ch.Max_hit += pl
			ch.Basepl += pl
			ch.Max_mana += ki
			ch.Baseki += ki
			ch.Max_move += stam
			ch.Basest += stam
			if !IS_NPC(vict) && !IS_NPC(ch) {
				send_to_imm(libc.CString("[PK] %s killed %s at room [%d]\r\n"), GET_NAME(ch), GET_NAME(vict), func() room_vnum {
					if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Number
					}
					return -1
				}())
				vict.Act[int(PLR_ABSORBED/32)] |= bitvector_t(1 << (int(PLR_ABSORBED % 32)))
			}
			send_to_char(ch, libc.CString("@D[@mINGEST@D] @rPL@W: @D(@y%s@D) @cKi@W: @D(@y%s@D) @gSt@W: @D(@y%s@D)@n\r\n"), add_commas(pl), add_commas(ki), add_commas(stam))
			if rand_number(1, 3) == 3 {
				send_to_char(ch, libc.CString("You get %s's eye color.\r\n"), GET_NAME(vict))
				ch.Eye = vict.Eye
			} else if rand_number(1, 3) == 3 {
				send_to_char(ch, libc.CString("%s changes your height.\r\n"), GET_NAME(vict))
				if GET_PC_HEIGHT(ch) > GET_PC_HEIGHT(vict) {
					ch.Height -= uint8(int8((GET_PC_HEIGHT(ch) - GET_PC_HEIGHT(vict)) / 2))
				} else if GET_PC_HEIGHT(ch) < GET_PC_HEIGHT(vict) {
					ch.Height += uint8(int8((GET_PC_HEIGHT(vict) - GET_PC_HEIGHT(ch)) / 2))
				} else {
					ch.Height = uint8(int8(GET_PC_HEIGHT(vict)))
				}
			} else if rand_number(1, 3) == 3 {
				send_to_char(ch, libc.CString("%s changes your weight.\r\n"), GET_NAME(vict))
				if GET_PC_WEIGHT(ch) > GET_PC_WEIGHT(vict) {
					ch.Weight -= uint8(int8((GET_PC_WEIGHT(ch) - GET_PC_WEIGHT(vict)) / 2))
				} else if GET_PC_WEIGHT(ch) < GET_PC_WEIGHT(vict) {
					ch.Weight += uint8(int8((GET_PC_WEIGHT(vict) - GET_PC_WEIGHT(ch)) / 2))
				} else {
					ch.Weight = uint8(int8(GET_PC_WEIGHT(vict)))
				}
			} else {
				send_to_char(ch, libc.CString("Your forelock length changes because of %s.\r\n"), GET_NAME(vict))
				ch.Hairl = vict.Hairl
			}
			handle_ingest_learn(ch, vict)
			die(vict, nil)
			return
		}
	} else {
		send_to_char(ch, libc.CString("You are not a majin, you can not ingest.\r\n"))
		return
	}
}
func do_absorb(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		arg  [2048]byte
		arg2 [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if check_skill(ch, SKILL_ABSORB) == 0 {
		return
	}
	if ch.Race == RACE_ANDROID {
		if limb_ok(ch, 0) == 0 {
			return
		}
	}
	if !IS_NPC(ch) {
		if ch.Race == RACE_BIO && !PLR_FLAGGED(ch, PLR_TAIL) {
			send_to_char(ch, libc.CString("You have no tail!\r\n"))
			return
		}
	}
	if ch.Race != RACE_ANDROID && ch.Race != RACE_BIO {
		send_to_char(ch, libc.CString("You shouldn't have this skill, you are incapable of absorbing.\r\n"))
		send_to_imm(libc.CString("ERROR: Absorb skill on %s when they are not a bio or android."), GET_NAME(ch))
		return
	}
	if ch.Fighting != nil && ch.Race != RACE_ANDROID {
		send_to_char(ch, libc.CString("You are too busy fighting!\r\n"))
		return
	}
	if ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are currently being grappled with! Try 'escape'!\r\n"))
		return
	}
	if ch.Grappling != nil {
		send_to_char(ch, libc.CString("You are currently grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil {
		act(libc.CString("@WYou stop absorbing from @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Absorbing), TO_CHAR)
		act(libc.CString("$n stops absorbing from you!"), TRUE, ch, nil, unsafe.Pointer(ch.Absorbing), TO_VICT)
		act(libc.CString("$n stops absorbing from $N!"), TRUE, ch, nil, unsafe.Pointer(ch.Absorbing), TO_NOTVICT)
		if IS_NPC(ch.Absorbing) && ch.Absorbing.Fighting == nil {
			set_fighting(ch.Absorbing, ch)
		}
		ch.Absorbing.Absorbby = nil
		ch.Absorbing = nil
	}
	if arg[0] == 0 && ch.Race == RACE_ANDROID {
		send_to_char(ch, libc.CString("Who do you want to absorb?\r\n"))
		return
	}
	if ch.Race == RACE_ANDROID {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Absorb %s?\r\n"), func() string {
				if ch.Race == RACE_ANDROID {
					return "from who"
				}
				return "who"
			}())
			return
		}
	}
	if ch.Race == RACE_BIO {
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("Syntax: absorb (swallow | extract) (target)\r\n"))
			return
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Syntax: absorb (swallow | extract) (target)\r\n"))
			return
		}
	}
	if AFF_FLAGGED(vict, AFF_SANCTUARY) {
		send_to_char(ch, libc.CString("You can't absorb them, they have a barrier!\r\n"))
		return
	}
	if ch.Race == RACE_ANDROID {
		if !IS_NPC(ch) {
			if !PLR_FLAGGED(ch, PLR_ABSORB) {
				send_to_char(ch, libc.CString("You are not an absorbtion model.\r\n"))
				return
			}
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if vict.Absorbby != nil {
			send_to_char(ch, libc.CString("%s is already absorbing from them!"), GET_NAME(vict.Absorbby))
			return
		}
		if vict.Max_hit > ch.Max_hit*2 {
			send_to_char(ch, libc.CString("They are too strong for you to absorb from.\r\n"))
			return
		}
		if vict.Max_hit*20 < ch.Max_hit {
			send_to_char(ch, libc.CString("They are too weak for you to bother absorbing from.\r\n"))
			return
		}
		if vict.Move < (vict.Max_move/20) && vict.Mana < (vict.Max_mana/20) {
			send_to_char(ch, libc.CString("They have nothing to absorb right now, they are drained...\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		if init_skill(ch, SKILL_ABSORB) < axion_dice(0) {
			act(libc.CString("@WYou rush at @c$N@W and try to absorb from them, but $E manages to avoid you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W rushes at you and tries to grab you, but you manage to avoid $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w rushes at @c$N@W and tries to grab $M, but @c$N@W manages to avoid @c$n@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_ABSORB, 1)
			if IS_NPC(vict) && IS_HUMANOID(vict) && rand_number(1, 3) == 3 {
				if ch.Fighting == nil {
					set_fighting(ch, vict)
				}
				if vict.Fighting == nil {
					set_fighting(vict, ch)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else {
			act(libc.CString("@WYou rush at @c$N@W and try to absorb from them, and manage to grab on!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W rushes at you and tries to grab you, and manages to grab on!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w rushes at @c$N@W and tries to grab $M, and manages to grab on!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_ABSORB, 1)
			ch.Absorbing = vict
			vict.Absorbby = ch
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		}
	} else if ch.Race == RACE_BIO && C.strcmp(&arg[0], libc.CString("swallow")) == 0 {
		if vict.Absorbby != nil {
			send_to_char(ch, libc.CString("%s is already absorbing from them!"), GET_NAME(vict.Absorbby))
			return
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if ch.Absorbs < 1 {
			send_to_char(ch, libc.CString("You already have already absorbed 3 people.\r\n"))
			return
		}
		if vict.Max_hit >= ch.Basepl*3 {
			send_to_char(ch, libc.CString("You are too weak to absorb them into your cellular structure!\r\n"))
			return
		}
		reveal_hiding(ch, 0)
		if GET_SKILL(ch, SKILL_ABSORB) < axion_dice(0) {
			act(libc.CString("@WYou rush at @c$N@W and try to absorb from them, but $E manages to avoid you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W rushes at you and tries to grab you, but you manage to avoid $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w rushes at @c$N@W and tries to grab $M, but @c$N@W manages to avoid @c$n@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_ABSORB, 1)
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else {
			act(libc.CString("@WYou rush at @c$N@W and your tail engulfs $M! You quickly suck $S squirming body into your tail, absorbing $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W rushes at you and $s tail engulfs you! $e quickly sucks your squirming body into $s tail, absorbing you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w rushes at @c$N@W and $s tail engulfs $M! You quickly suck $S squirming body into your tail, absorbing @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Absorbs -= 1
			ch.Max_hit += vict.Basepl / 5
			ch.Basepl += vict.Basepl / 5
			ch.Max_mana += vict.Baseki / 5
			ch.Baseki += vict.Baseki / 5
			ch.Max_move += vict.Basest / 5
			ch.Basest += vict.Basest / 5
			if !IS_NPC(vict) && !IS_NPC(ch) {
				send_to_imm(libc.CString("[PK] %s killed %s at room [%d]\r\n"), GET_NAME(ch), GET_NAME(vict), func() room_vnum {
					if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Number
					}
					return -1
				}())
				vict.Act[int(PLR_ABSORBED/32)] |= bitvector_t(1 << (int(PLR_ABSORBED % 32)))
			}
			var stam int64 = vict.Basest / 5
			var ki int64 = vict.Baseki / 5
			var pl int64 = vict.Basepl / 5
			send_to_char(ch, libc.CString("@D[@gABSORB@D] @rPL@W: @D(@y%s@D) @cKi@W: @D(@y%s@D) @gSt@W: @D(@y%s@D)@n\r\n"), add_commas(pl), add_commas(ki), add_commas(stam))
			improve_skill(ch, SKILL_ABSORB, 1)
			die(vict, nil)
		}
	} else if ch.Race == RACE_BIO && C.strcmp(&arg[0], libc.CString("extract")) == 0 {
		var failthresh int = rand_number(1, 125)
		if GET_LEVEL(vict) > 99 {
			failthresh += (GET_LEVEL(vict) - 95) * 2
		}
		if vict.Absorbby != nil {
			send_to_char(ch, libc.CString("%s is already absorbing from them!"), GET_NAME(vict.Absorbby))
			return
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if vict.Max_hit >= ch.Max_hit {
			send_to_char(ch, libc.CString("You are too weak to absorb them into your cellular structure!\r\n"))
			return
		}
		if vict.Max_hit < ch.Max_hit/5 {
			send_to_char(ch, libc.CString("They would be worthless to you at your strength!\r\n"))
			return
		}
		if !IS_NPC(vict) {
			send_to_char(ch, libc.CString("You can't absorb their bio extract, you need to swallow them with your tail!\r\n"))
			return
		}
		if !soft_cap(ch, 0) {
			send_to_char(ch, libc.CString("You can not handle any more bio extract at your current level.\r\n"))
			return
		}
		if GET_SKILL(ch, SKILL_ABSORB) < failthresh {
			act(libc.CString("@WYou rush at @c$N@W and try to absorb from them, but $E manages to avoid you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W rushes at you and tries to grab you, but you manage to avoid $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@w rushes at @c$N@W and tries to grab $M, but @c$N@W manages to avoid @c$n@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			improve_skill(ch, SKILL_ABSORB, 0)
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		} else {
			act(libc.CString("@WYou rush at @c$N@W and stab them with your tail! You quickly suck out all the bio extract you need and leave the empty husk behind!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@w rushes at @c$N@W and stabs $M with $s tail! $e quickly sucks out all the bio extract and leaves the empty husk of @c$N@W behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			var stam int64 = vict.Basest / 5000
			var ki int64 = vict.Baseki / 5000
			var pl int64 = vict.Basepl / 5000
			stam += int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2))
			pl += int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2))
			ki += int64(rand_number(GET_LEVEL(ch), GET_LEVEL(ch)*2))
			if stam > 1500000 {
				stam = 1500000
			}
			if pl > 1500000 {
				pl = 1500000
			}
			if ki > 1500000 {
				ki = 1500000
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				ch.Max_hit += pl * 2
				ch.Max_move += stam * 2
				ch.Max_mana += ki * 2
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				ch.Max_hit += pl * 3
				ch.Max_move += stam * 3
				ch.Max_mana += ki * 3
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				ch.Max_hit += int64(float64(pl) * 3.5)
				ch.Max_move += int64(float64(stam) * 3.5)
				ch.Max_mana += int64(float64(ki) * 3.5)
			} else if PLR_FLAGGED(ch, PLR_TRANS4) {
				ch.Max_hit += pl * 4
				ch.Max_move += stam * 4
				ch.Max_mana += ki * 4
			} else {
				ch.Max_hit += pl
				ch.Max_move += stam
				ch.Max_mana += ki
			}
			ch.Basepl += pl
			ch.Basest += stam
			ch.Baseki += ki
			ch.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.05)
			if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
				ch.Lifeforce = int64(GET_LIFEMAX(ch))
			}
			send_to_char(ch, libc.CString("@D[@gABSORB@D] @rPL@W: @D(@y%s@D) @cKi@W: @D(@y%s@D) @gSt@W: @D(@y%s@D)@n\r\n"), add_commas(pl), add_commas(ki), add_commas(stam))
			improve_skill(ch, SKILL_ABSORB, 0)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
			vict.Act[int(MOB_HUSK/32)] |= bitvector_t(1 << (int(MOB_HUSK % 32)))
			die(vict, ch)
		}
	} else {
		if ch.Race != RACE_BIO && ch.Race != RACE_ANDROID {
			send_to_char(ch, libc.CString("You have the absorb skill but are incapable of absorbing. This error has been reported.\r\n"))
			send_to_imm(libc.CString("ERROR: Absorb attempted by %s even though they are not bio or android."), GET_NAME(ch))
		} else {
			send_to_char(ch, libc.CString("Syntax: absorb (extract | swallow) (target)\r\n"))
		}
		return
	}
}
func do_escape(ch *char_data, argument *byte, cmd int, subcmd int) {
	if ch.Absorbby == nil && ch.Grappled == nil {
		send_to_char(ch, libc.CString("You are not in anyone's grasp!\r\n"))
		return
	}
	var num int = int(ch.Aff_abils.Str)
	if ch.Absorbby != nil {
		var skill int = GET_SKILL(ch.Absorbby, SKILL_ABSORB)
		if ch.Hit > ch.Absorbby.Hit*10 {
			num += rand_number(10, 15)
		} else if ch.Hit > ch.Absorbby.Hit*5 {
			num += rand_number(6, 10)
		} else if ch.Hit > ch.Absorbby.Hit*2 {
			num += rand_number(4, 8)
		} else if ch.Hit > ch.Absorbby.Hit {
			num += rand_number(2, 5)
		} else if ch.Hit*10 <= ch.Absorbby.Hit {
			skill -= rand_number(10, 15)
		} else if ch.Hit*5 <= ch.Absorbby.Hit {
			skill -= rand_number(6, 10)
		} else if ch.Hit*2 <= ch.Absorbby.Hit {
			skill -= rand_number(4, 8)
		} else if ch.Hit < ch.Absorbby.Hit {
			skill -= rand_number(2, 5)
		}
		if num > skill {
			act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_NOTVICT)
			act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_CHAR)
			if ch.Fighting == nil {
				set_fighting(ch, ch.Absorbby)
			}
			if ch.Absorbby.Fighting == nil {
				set_fighting(ch.Absorbby, ch)
			}
			ch.Absorbby.Absorbing = nil
			ch.Absorbby = nil
		} else {
			act(libc.CString("@c$N@W struggles to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_NOTVICT)
			act(libc.CString("@WYou struggle to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("@c$N@W struggles to break loose of your hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_CHAR)
			if rand_number(1, 3) == 3 {
				var dmg int64 = int64(float64(ch.Max_hit) * 0.025)
				hurt(0, 0, ch, ch.Absorbby, nil, dmg, 0)
				if ch.Absorbby.Position == POS_SLEEPING {
					act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_VICT)
					act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_CHAR)
					ch.Absorbby.Absorbing = nil
					ch.Absorbby = nil
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
	}
	if ch.Grappled != nil {
		var skill int = GET_SKILL(ch.Grappled, SKILL_GRAPPLE)
		if ch.Hit > ch.Grappled.Hit*10 {
			num += rand_number(10, 15)
		} else if ch.Hit > ch.Grappled.Hit*5 {
			num += rand_number(6, 10)
		} else if ch.Hit > ch.Grappled.Hit*2 {
			num += rand_number(4, 8)
		} else if ch.Hit > ch.Grappled.Hit {
			num += rand_number(2, 5)
		} else if ch.Hit*10 <= ch.Grappled.Hit {
			skill -= rand_number(10, 15)
		} else if ch.Hit*5 <= ch.Grappled.Hit {
			skill -= rand_number(6, 10)
		} else if ch.Hit*2 <= ch.Grappled.Hit {
			skill -= rand_number(4, 8)
		} else if ch.Hit < ch.Grappled.Hit {
			skill -= rand_number(2, 5)
		}
		if num > skill {
			if ch.Grappled.Grap == 4 {
				act(libc.CString("@c$N@M flexes with all $S might and causes your body to explode outward into gooey chunks!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_CHAR)
				act(libc.CString("@MYou flex with all your might and cause @C$n's@M body to explode outward into gooey chunks!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_VICT)
				act(libc.CString("@c$N@M flexes with all $S might and causes @C$n's@M body to explode outward into gooey chunks!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_NOTVICT)
				act(libc.CString("@MYou reform your body mere moments later.@n"), TRUE, ch.Grappled, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@M reforms $s body mere moments later."), TRUE, ch.Grappled, nil, nil, TO_ROOM)
			} else {
				act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_NOTVICT)
				act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_VICT)
				act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_CHAR)
			}
			if ch.Fighting == nil {
				set_fighting(ch, ch.Grappled)
			}
			if ch.Grappled.Fighting == nil {
				set_fighting(ch.Grappled, ch)
			}
			ch.Grappled.Grap = -1
			ch.Grappled.Grappling = nil
			ch.Grappled = nil
			ch.Grap = -1
		} else {
			act(libc.CString("@c$N@W struggles to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_NOTVICT)
			act(libc.CString("@WYou struggle to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("@c$N@W struggles to break loose of your hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_CHAR)
			if rand_number(1, 3) == 3 {
				var dmg int64 = int64(float64(ch.Max_hit) * 0.025)
				hurt(0, 0, ch, ch.Grappled, nil, dmg, 0)
				if ch.Grappled.Position == POS_SLEEPING {
					act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_VICT)
					act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_CHAR)
					ch.Grappled.Grap = -1
					ch.Grappled.Grappling = nil
					ch.Grappled = nil
					ch.Grap = -1
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
	}
}
func do_regenerate(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		amt   int64 = 0
		skill int   = 0
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	skill = init_skill(ch, SKILL_REGENERATE)
	if skill < 1 {
		send_to_char(ch, libc.CString("You are incapable of regenerating.\r\n"))
		return
	}
	if ch.Suppression > 0 {
		skill = int(ch.Suppression)
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Regenerate how much PL?\r\nMax percent you can regen: %d\r\nSyntax: regenerate (1 - 100)\r\n"), skill)
		return
	}
	if ch.Hit >= gear_pl(ch) {
		send_to_char(ch, libc.CString("You do not need to regenerate, you are at full health.\r\n"))
		return
	}
	if ch.Suppression > 0 && ch.Hit >= ((gear_pl(ch)/100)*ch.Suppression) {
		send_to_char(ch, libc.CString("You do not need to regenerate, you are at full health.\r\n"))
		return
	}
	var num int = libc.Atoi(libc.GoString(&arg[0]))
	if num <= 0 {
		send_to_char(ch, libc.CString("What is the point of that?\r\nSyntax: regenerate (1 - 100)\r\n"))
		return
	}
	if num > 100 {
		send_to_char(ch, libc.CString("You can't regenerate that much!\r\nMax you can regen: %d\r\n"), skill)
		return
	}
	if num > skill {
		send_to_char(ch, libc.CString("You can't regenerate that much!\r\nMax you can regen: %d\r\n"), skill)
		return
	}
	if ch.Suppression > 0 && num > int(ch.Suppression) {
		send_to_char(ch, libc.CString("You can't regenerate that much!\r\nMax you can regen: %d\r\n"), skill)
		return
	}
	amt = int64((float64(gear_pl(ch)) * 0.01) * float64(num))
	if amt > 1 {
		amt /= 2
	}
	if ch.Race == RACE_BIO {
		amt = int64(float64(amt) * 0.9)
	}
	var life int64 = int64(float64(ch.Lifeforce) - float64(amt)*0.8)
	var energy int64 = int64(float64(ch.Mana) - float64(amt)*0.2)
	if (life <= 0 || energy <= 0) && !IS_NPC(ch) {
		send_to_char(ch, libc.CString("Your life force or ki are too low to regenerate that much.\r\n"))
		send_to_char(ch, libc.CString("@YLF Needed@D: @C%s@w, @YKi Needed@D: @C%s@w.@n\r\n"), add_commas(int64(float64(amt)*0.8)), add_commas(int64(float64(amt)*0.2)))
		return
	} else if IS_NPC(ch) && energy <= 0 {
		return
	}
	ch.Hit += amt * 2
	if !IS_NPC(ch) {
		ch.Lifeforce -= int64(float64(amt) * 0.8)
	}
	ch.Mana -= int64(float64(amt) * 0.2)
	if ch.Hit > gear_pl(ch) {
		ch.Hit = gear_pl(ch)
	}
	reveal_hiding(ch, 0)
	if ch.Suppression > 0 && ch.Hit > ((ch.Max_hit/100)*ch.Suppression) {
		ch.Hit = (ch.Max_hit / 100) * ch.Suppression
		send_to_char(ch, libc.CString("@mYou regenerate to your suppression limit.@n\r\n"))
	}
	if ch.Hit >= gear_pl(ch) {
		act(libc.CString("You concentrate your ki and regenerate your body completely."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n concentrates and regenerates $s body completely."), TRUE, ch, nil, nil, TO_ROOM)
	} else if amt < ch.Max_hit/10 {
		act(libc.CString("You concentrate your ki and regenerate your body a little."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n concentrates and regenerates $s body a little."), TRUE, ch, nil, nil, TO_ROOM)
	} else if amt < ch.Max_hit/5 {
		act(libc.CString("You concentrate your ki and regenerate your body some."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n concentrates and regenerates $s body some."), TRUE, ch, nil, nil, TO_ROOM)
	} else if amt < ch.Max_hit/2 {
		act(libc.CString("You concentrate your ki and regenerate your body a great deal."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n concentrates and regenerates $s body a great deal."), TRUE, ch, nil, nil, TO_ROOM)
	} else if ch.Hit < ch.Max_hit {
		act(libc.CString("You concentrate your ki and regenerate you nearly completely."), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n concentrates and regenerates $s body nearly completely."), TRUE, ch, nil, nil, TO_ROOM)
	}
	improve_skill(ch, SKILL_REGENERATE, 0)
	if AFF_FLAGGED(ch, AFF_BURNED) {
		send_to_char(ch, libc.CString("Your burns are healed now.\r\n"))
		act(libc.CString("$n@w's burns are now healed.@n"), TRUE, ch, nil, nil, TO_ROOM)
		null_affect(ch, AFF_BURNED)
	}
	if !IS_NPC(ch) {
		if (ch.Limb_condition[0]) <= 0 {
			act(libc.CString("You regrow your right arm!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regrows $s right arm!"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Limb_condition[0] = 100
		} else if (ch.Limb_condition[0]) >= 0 && (ch.Limb_condition[0]) < 50 {
			act(libc.CString("Your broken right arm mends itself!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regenerates $s broken right arm!"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Limb_condition[0] = 100
		}
		if (ch.Limb_condition[1]) <= 0 {
			ch.Limb_condition[1] = 100
			act(libc.CString("You regrow your left arm!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regrows $s left arm!"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 {
			act(libc.CString("Your broken left arm mends itself!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regenerates $s broken left arm!"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Limb_condition[1] = 100
		}
		if (ch.Limb_condition[3]) <= 0 {
			ch.Limb_condition[3] = 100
			act(libc.CString("You regrow your left leg!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regrows $s left leg!"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 {
			act(libc.CString("Your broken left leg mends itself!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regenerates $s broken left leg!"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Limb_condition[3] = 100
		}
		if (ch.Limb_condition[2]) <= 0 {
			ch.Limb_condition[2] = 100
			act(libc.CString("You regrow your right leg!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regrows $s right leg!"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 {
			ch.Limb_condition[2] = 100
			act(libc.CString("Your broken right leg mends itself!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regenerates $s broken right leg!"), TRUE, ch, nil, nil, TO_ROOM)
		}
		if !PLR_FLAGGED(ch, PLR_TAIL) && ch.Race == RACE_BIO {
			ch.Act[int(PLR_TAIL/32)] |= bitvector_t(1 << (int(PLR_TAIL % 32)))
			act(libc.CString("You regrow your tail!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n regrows $s tail!"), TRUE, ch, nil, nil, TO_ROOM)
		}
		improve_skill(ch, SKILL_REGENERATE, 0)
	}
}
func do_focus(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data = nil
		arg  [2048]byte
		name [2048]byte
	)
	name[0] = '\x00'
	arg[0] = '\x00'
	two_arguments(argument, &arg[0], &name[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Yes but what do you want to focus?\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_CURSE) {
		send_to_char(ch, libc.CString("You are cursed and can't focus!\r\n"))
		return
	} else if C.strcmp(&arg[0], libc.CString("tough")) == 0 {
		if know_skill(ch, SKILL_TSKIN) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_STONESKIN) {
				send_to_char(ch, libc.CString("You already have tough skin!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to infuse into your skin.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_TSKIN) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your skin, but fail in making it tough!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s skin, but fails in making it tough!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				var duration int = int(ch.Aff_abils.Intel / 20)
				assign_affect(ch, AFF_STONESKIN, SKILL_TSKIN, duration, 0, 0, 0, 0, 0, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your skin, making it tough!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s skin, making it tough!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Focus your ki into who's skin?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_STONESKIN) {
					send_to_char(ch, libc.CString("They already have tough skin!\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to infuse into their skin.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_TSKIN) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's skin, but fail in making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your skin, but fails in making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's skin, but fails in making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
					assign_affect(vict, AFF_STONESKIN, SKILL_TSKIN, duration, 0, 0, 0, 0, 0, 0)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's skin, making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your skin, making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's skin, making it tough!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("might")) == 0 {
		if know_skill(ch, SKILL_MIGHT) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_MIGHT) {
				send_to_char(ch, libc.CString("You already have mighty muscles!\r\n"))
				return
			} else if (ch.Bonuses[BONUS_WIMP]) > 0 && ch.Aff_abils.Str+10 > 25 {
				send_to_char(ch, libc.CString("Your body is not able to withstand increasing its strength beyond 25.\r\n"))
				return
			} else if (ch.Bonuses[BONUS_FRAIL]) > 0 && ch.Aff_abils.Str+2 > 25 {
				send_to_char(ch, libc.CString("Your body is not able to withstand increasing its strength beyond 25.\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to infuse into your muscles.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_MIGHT) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your muscles, but fail in making them mighty!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s muscles, but fails in making them mighty!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				ch.Affected_by[int(AFF_MIGHT/32)] |= 1 << (int(AFF_MIGHT % 32))
				ch.Mana -= ch.Max_mana / 20
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
				assign_affect(ch, AFF_MIGHT, SKILL_MIGHT, duration, 10, 2, 0, 0, 0, 0)
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your muscles, making them mighty!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s muscles, making them mighty!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Focus your ki into who's muscles?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_MIGHT) {
					send_to_char(ch, libc.CString("They already have mighty muscles!\r\n"))
					return
				} else if (vict.Bonuses[BONUS_WIMP]) > 0 && vict.Aff_abils.Str+10 > 25 {
					send_to_char(ch, libc.CString("Their body is not able to withstand increasing its strength beyond 25.\r\n"))
					return
				} else if (vict.Bonuses[BONUS_FRAIL]) > 0 && vict.Aff_abils.Con+2 > 25 {
					send_to_char(ch, libc.CString("Their body is not able to withstand increasing its constitution beyond 25.\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to infuse into their muscles.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_MIGHT) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's muscles, but fail in making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your muscles, but fails in making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's muscles, but fails in making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					ch.Mana -= ch.Max_mana / 20
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
					assign_affect(vict, AFF_MIGHT, SKILL_MIGHT, duration, 10, 2, 0, 0, 0, 0)
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's muscles, making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your muscles, making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's muscles, making them mighty!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("wither")) == 0 {
		if know_skill(ch, SKILL_WITHER) == 0 {
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Focus your ki into who's muscles?\r\n"))
			return
		}
		if ch == vict {
			send_to_char(ch, libc.CString("You don't want to wither your own body!\r\n"))
			return
		}
		if can_kill(ch, vict, nil, 2) == 0 {
			return
		}
		if AFF_FLAGGED(vict, AFF_WITHER) {
			send_to_char(ch, libc.CString("They already have been withered!\r\n"))
			return
		} else if ch.Mana < ch.Max_mana/20 {
			send_to_char(ch, libc.CString("You do not have enough ki to wither them.\r\n"))
			return
		} else if GET_SKILL(ch, SKILL_WITHER) < axion_dice(0) {
			ch.Mana -= ch.Max_mana / 20
			reveal_hiding(ch, 0)
			act(libc.CString("You focus ki into $N's body, but fail in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("$n focuses ki into your body, but fails in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n focuses ki into $N's body, but fails in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			return
		} else {
			vict.Affected_by[int(AFF_WITHER/32)] |= 1 << (int(AFF_WITHER % 32))
			ch.Mana -= ch.Max_mana / 20
			vict.Real_abils.Str -= 3
			vict.Real_abils.Cha -= 3
			save_char(vict)
			reveal_hiding(ch, 0)
			act(libc.CString("You focus ki into $N's body, and succeed in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("$n focuses ki into your body, and succeeds in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n focuses ki into $N's body, and succeeds in withering it!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			return
		}
	} else if C.strcmp(&arg[0], libc.CString("enlighten")) == 0 {
		if know_skill(ch, SKILL_ENLIGHTEN) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_ENLIGHTEN) {
				send_to_char(ch, libc.CString("You already have superior wisdom!\r\n"))
				return
			} else if (ch.Bonuses[BONUS_FOOLISH]) > 0 && ch.Aff_abils.Wis+10 > 25 {
				send_to_char(ch, libc.CString("You're not able to withstand increasing your wisdom beyond 25.\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to use this skill.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_ENLIGHTEN) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your mind, but fail in awakening it to cosmic wisdom!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s mind, but fails in awakening it to cosmic wisdom!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
				assign_affect(ch, AFF_ENLIGHTEN, SKILL_ENLIGHTEN, duration, 0, 0, 0, 0, 10, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your mind, awakening it to cosmic wisdom!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s mind, awakening it to cosmic wisdom!"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Chclass == CLASS_JINTO && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 15 && rand_number(1, 4) >= 3 {
					var gain int64 = 0
					ch.Player_specials.Class_skill_points[ch.Chclass] -= 15
					if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 100 {
						gain = int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.15)
						ch.Exp += gain
						send_to_char(ch, libc.CString("@GYou gain @g%s@G experience due to your excellence with this skill.@n\r\n"), add_commas(gain))
					} else if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 60 {
						gain = int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.1)
						ch.Exp += gain
						send_to_char(ch, libc.CString("@GYou gain @g%s@G experience due to your excellence with this skill.@n\r\n"), add_commas(gain))
					} else if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 40 {
						gain = int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.05)
						ch.Exp += gain
						send_to_char(ch, libc.CString("@GYou gain @g%s@G experience due to your excellence with this skill.@n\r\n"), add_commas(gain))
					}
				}
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Focus your ki into who's mind?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_ENLIGHTEN) {
					send_to_char(ch, libc.CString("They already have superior wisdom!\r\n"))
					return
				} else if (vict.Bonuses[BONUS_FOOLISH]) > 0 && vict.Aff_abils.Wis+10 > 25 {
					send_to_char(ch, libc.CString("They're not able to withstand increasing their wisdom beyond 25.\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to use this skill.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_ENLIGHTEN) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's mind, but fail in awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your mind, but fails in awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's mind, but fails in awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
					assign_affect(vict, AFF_ENLIGHTEN, SKILL_ENLIGHTEN, duration, 0, 0, 0, 0, 10, 0)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's mind, awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your mind, awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's mind, awakening it to cosmic wisdom!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if ch.Chclass == CLASS_JINTO && level_exp(vict, GET_LEVEL(vict)+1)-int(vict.Exp) > 0 && (ch.Player_specials.Class_skill_points[ch.Chclass]) >= 15 && rand_number(1, 4) >= 3 {
						var gain int64 = 0
						ch.Player_specials.Class_skill_points[ch.Chclass] -= 15
						if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 100 {
							gain = int64(float64(level_exp(vict, GET_LEVEL(vict)+1)) * 0.15)
							vict.Exp += gain
							send_to_char(vict, libc.CString("@GYou gain @g%s@G experience due to the level of enlightenment you have received!@n\r\n"), add_commas(gain))
						} else if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 60 {
							gain = int64(float64(level_exp(vict, GET_LEVEL(vict)+1)) * 0.1)
							vict.Exp += gain
							send_to_char(vict, libc.CString("@GYou gain @g%s@G experience due to the level of enlightenment you have received!@n\r\n"), add_commas(gain))
						} else if GET_SKILL(ch, SKILL_ENLIGHTEN) >= 40 {
							gain = int64(float64(level_exp(vict, GET_LEVEL(vict)+1)) * 0.05)
							vict.Exp += gain
							send_to_char(vict, libc.CString("@GYou gain @g%s@G experience due to the level of enlightenment you have received!@n\r\n"), add_commas(gain))
						}
					}
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("genius")) == 0 {
		if know_skill(ch, SKILL_GENIUS) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_GENIUS) {
				send_to_char(ch, libc.CString("You already have superior intelligence!\r\n"))
				return
			} else if (ch.Bonuses[BONUS_DULL]) > 0 && ch.Aff_abils.Intel+10 > 25 {
				send_to_char(ch, libc.CString("You're not able to withstand increasing your intelligence beyond 25.\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to infuse into your mind.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_GENIUS) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your mind, but fail in making it work faster!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s muscles, but fails in making it work faster!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
				assign_affect(ch, AFF_GENIUS, SKILL_GENIUS, duration, 0, 0, 10, 0, 0, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your mind, making it work faster!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s mind, making it work faster!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Focus your ki into who's mind?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_GENIUS) {
					send_to_char(ch, libc.CString("They already have superior intelligence!\r\n"))
					return
				} else if (vict.Bonuses[BONUS_DULL]) > 0 && vict.Aff_abils.Intel+10 > 25 {
					send_to_char(ch, libc.CString("They're not able to withstand increasing their intelligence beyond 25.\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to infuse into their mind.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_GENIUS) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's mind, but fail in making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your mind, but fails in making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's mind, but fails in making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
					assign_affect(vict, AFF_GENIUS, SKILL_GENIUS, duration, 0, 0, 10, 0, 0, 0)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's mind, making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your mind, making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's mind, making it work faster!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (vict.Master == ch || ch.Master == vict || ch.Master == vict.Master) && AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
						if ch.Race == RACE_KAI && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 && rand_number(1, 3) == 3 {
							ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.05)
						}
					}
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("flex")) == 0 {
		if know_skill(ch, SKILL_FLEX) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_FLEX) {
				send_to_char(ch, libc.CString("You already have superior agility!\r\n"))
				return
			} else if (ch.Bonuses[BONUS_CLUMSY]) > 0 && ch.Aff_abils.Dex+10 > 25 {
				send_to_char(ch, libc.CString("You're not able to withstand increasing your agility beyond 25.\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to infuse into your limbs.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_FLEX) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your limbs, but fail in making them more flexible!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s muscles, but fails in making them more flexible!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				ch.Mana -= ch.Max_mana / 20
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
				assign_affect(ch, AFF_FLEX, SKILL_FLEX, duration, 0, 0, 0, 10, 0, 0)
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your limbs, making them more flexible!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki into $s limbs, making them more flexible!"), TRUE, ch, nil, nil, TO_ROOM)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Focus your ki into who's limbs?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_FLEX) {
					send_to_char(ch, libc.CString("They already have superior agility!\r\n"))
					return
				} else if (vict.Bonuses[BONUS_CLUMSY]) > 0 && vict.Aff_abils.Dex+3 > 25 {
					send_to_char(ch, libc.CString("They're not able to withstand increasing their agility beyond 25.\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to infuse into their limbs.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_FLEX) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's limbs, but fail in making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your limbs, but fails in making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's limbs, but fails in making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					ch.Mana -= ch.Max_mana / 20
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 2)
					assign_affect(vict, AFF_FLEX, SKILL_FLEX, duration, 0, 0, 0, 10, 0, 0)
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's limbs, making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your limbs, making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki into $N's limbs, making them more flexible!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (vict.Master == ch || ch.Master == vict || ch.Master == vict.Master) && AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
						if ch.Race == RACE_KAI && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 && rand_number(1, 3) == 3 {
							ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.05)
						}
					}
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("bless")) == 0 {
		if know_skill(ch, SKILL_BLESS) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_BLESS) {
				send_to_char(ch, libc.CString("You already are blessed!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to bless.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_BLESS) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while chanting spiritual words. Your blessing does nothing though, you must have messed up!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki while chanting spiritual words. $n seems disappointed."), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 3)
				assign_affect(ch, AFF_BLESS, SKILL_BLESS, duration, 0, 0, 0, 0, 0, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				if ch.Chclass == CLASS_KABITO {
					ch.Blesslvl = GET_SKILL(ch, SKILL_BLESS)
				} else {
					ch.Blesslvl = 0
				}
				act(libc.CString("You focus ki while chanting spiritual words. You feel your body recovering at above normal speed!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki while chanting spiritual words. $n smiles after finishing $s chant."), TRUE, ch, nil, nil, TO_ROOM)
				if AFF_FLAGGED(ch, AFF_CURSE) {
					send_to_char(ch, libc.CString("Your cursing was nullified!\r\n"))
					null_affect(ch, AFF_CURSE)
				}
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Bless who?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_BLESS) {
					send_to_char(ch, libc.CString("They already have been blessed!\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to bless.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_BLESS) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki while chanting spiritual words. Your blessing fails!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("$n focuses ki while chanting spiritual words. $n places a hand on your head, but nothing happens!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki while chanting spiritual words. $n places a hand on $N's head, but nothing happens!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 3)
					assign_affect(vict, AFF_BLESS, SKILL_BLESS, duration, 0, 0, 0, 0, 0, 0)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					if ch.Race == RACE_KAI {
						vict.Blesslvl = GET_SKILL(ch, SKILL_BLESS)
					} else {
						vict.Blesslvl = 0
					}
					act(libc.CString("You focus ki while chanting spiritual words. Blessing $N with faster regeneration!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki while chanting spiritual words. $n then places a hand on your head, blessing you!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki while chanting spiritual words. $n then places a hand on $N's head, blessing them!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (vict.Master == ch || ch.Master == vict || ch.Master == vict.Master) && AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(vict, AFF_GROUP) {
						if ch.Race == RACE_KAI && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 && rand_number(1, 3) == 3 {
							ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.05)
						}
					}
					if AFF_FLAGGED(vict, AFF_CURSE) {
						send_to_char(vict, libc.CString("Your cursing was nullified!\r\n"))
						null_affect(vict, AFF_CURSE)
					}
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("curse")) == 0 {
		if know_skill(ch, SKILL_CURSE) == 0 {
			return
		}
		if name[0] == 0 {
			if AFF_FLAGGED(ch, AFF_CURSE) {
				send_to_char(ch, libc.CString("You already are cursed!\r\n"))
				return
			} else if ch.Race == RACE_DEMON {
				send_to_char(ch, libc.CString("You are immune to curses!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to CURSE.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_CURSE) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while chanting demonic words. Your cursing does nothing though, you must have messed up!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki while chanting demonic words. $n seems disappointed."), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 3)
				assign_affect(vict, AFF_CURSE, SKILL_CURSE, duration, 0, 0, 0, 0, 0, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while chanting demonic words. You feel your body recovering at below normal speed!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki while chanting demonic words. $n grins after finishing $s chant."), TRUE, ch, nil, nil, TO_ROOM)
				if AFF_FLAGGED(ch, AFF_BLESS) {
					send_to_char(ch, libc.CString("Your blessing was nullified!\r\n"))
					null_affect(ch, AFF_BLESS)
				}
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("Curse who?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 0) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if AFF_FLAGGED(vict, AFF_CURSE) {
					send_to_char(ch, libc.CString("They already have been cursed!\r\n"))
					return
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if vict.Race == RACE_DEMON {
					send_to_char(ch, libc.CString("They are immune to curses!\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to CURSE.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_CURSE) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki while chanting demonic words. Your cursing fails!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("$n focuses ki while chanting demonic words. $n places a hand on your head, but nothing happens!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki while chanting demonic words. $n places a hand on $N's head, but nothing happens!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					var duration int = roll_aff_duration(int(ch.Aff_abils.Intel), 3)
					assign_affect(vict, AFF_CURSE, SKILL_CURSE, duration, 0, 0, 0, 0, 0, 0)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki while chanting demonic words. cursing $N with slower regeneration!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki while chanting demonic words. $n then places a hand on your head, cursing you!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki while chanting demonic words. $n then places a hand on $N's head, cursing them!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if AFF_FLAGGED(vict, AFF_BLESS) {
						send_to_char(vict, libc.CString("Your blessing was nullified!\r\n"))
						null_affect(vict, AFF_BLESS)
					}
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("yoikominminken")) == 0 || C.strcmp(&arg[0], libc.CString("yoik")) == 0 {
		if know_skill(ch, SKILL_YOIK) == 0 {
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Use Yoikominminken on who?\r\n"))
			return
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		} else {
			if AFF_FLAGGED(vict, AFF_SLEEP) {
				send_to_char(ch, libc.CString("They already have been put to sleep!\r\n"))
				return
			} else if PLR_FLAGGED(vict, PLR_EYEC) {
				send_to_char(ch, libc.CString("Their eyes are closed!\r\n"))
				return
			} else if AFF_FLAGGED(vict, AFF_BLIND) {
				send_to_char(ch, libc.CString("They appear to be blind!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to use Yoikominminken.\r\n"))
				return
			} else if (vict.Bonuses[BONUS_INSOMNIAC]) != 0 {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while moving your hands in lulling patterns, but $N doesn't look the least bit sleepy!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, but you just don't feel tired."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, but $N doesn't look the least bit sleepy!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				return
			} else if GET_SKILL(ch, SKILL_YOIK) < axion_dice(0) || int(ch.Aff_abils.Intel)+rand_number(1, 3) < int(vict.Aff_abils.Intel)+rand_number(1, 5) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while moving your hands in lulling patterns, but fail to put $N to sleep!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, but you resist the technique!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, but $N resists the technique!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				return
			} else {
				var duration int = rand_number(1, 2)
				assign_affect(vict, AFF_SLEEP, SKILL_YOIK, duration, 0, 0, 0, 0, 0, 0)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki while moving your hands in lulling patterns, putting $N to sleep!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, before you realise it you are asleep!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n focuses ki while moving $s hands in a lulling pattern, putting $N to sleep!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Position = POS_SLEEPING
				if AFF_FLAGGED(vict, AFF_FLYING) {
					vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
					vict.Altitude = 0
				}
				return
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("vigor")) == 0 {
		if know_skill(ch, SKILL_VIGOR) == 0 {
			return
		}
		if name[0] == 0 {
			if ch.Mana < ch.Max_mana/10 {
				send_to_char(ch, libc.CString("You do not have enough ki to use vigor.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_VIGOR) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 10
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your very cells, but fail at re-engerizing them!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki and glows green for a moment, $e then frowns."), TRUE, ch, nil, nil, TO_ROOM)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			} else if ch.Move >= ch.Max_move {
				send_to_char(ch, libc.CString("You already have full stamina.\r\n"))
				return
			} else {
				if (ch.Bonuses[BONUS_HEALER]) > 0 {
					ch.Move += ch.Max_mana / 8
					ch.Mana -= ch.Max_mana / 8
				} else {
					ch.Move += ch.Max_mana / 10
					ch.Mana -= ch.Max_mana / 10
				}
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki into your very cells, and manage to re-energize them!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki and glows green for a moment, $e then smiles."), TRUE, ch, nil, nil, TO_ROOM)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("VIGOR who?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if IS_NPC(vict) {
					send_to_char(ch, libc.CString("Whatever would you waste your ki on them for?\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/10 {
					send_to_char(ch, libc.CString("You do not have enough ki to use vigor.\r\n"))
					return
				} else if vict.Move >= vict.Max_move {
					send_to_char(ch, libc.CString("They already have full stamina.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_VIGOR) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 10
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's very cells, and fail at re-energizing them!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your very cells, but nothing happens!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki and $N glows green for a moment, $N frowns."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
					return
				} else {
					if (ch.Bonuses[BONUS_HEALER]) > 0 {
						vict.Move += vict.Max_mana / 8
						ch.Mana -= ch.Max_mana / 8
					} else {
						vict.Move += vict.Max_mana / 10
						ch.Mana -= ch.Max_mana / 10
					}
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki into $N's very cells, and manage to re-energize them!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki into your very cells, and manages to re-energize them!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki and $N glows green for a moment, $N smiles."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("cure")) == 0 {
		if know_skill(ch, SKILL_CURE) == 0 {
			return
		}
		if name[0] == 0 {
			if !AFF_FLAGGED(ch, AFF_POISON) {
				send_to_char(ch, libc.CString("You are not poisoned!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to cure.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_CURE) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki and aim a pulsing light at your body. Nothing happens!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki and aims a pulsing light at $s body. Nothing seems to happen."), TRUE, ch, nil, nil, TO_ROOM)
				return
			} else {
				affect_from_char(ch, SPELL_POISON)
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki and aim a pulsing light at your body. You feel the poison in your blood disappear!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("$n focuses ki and aims a pulsing light at $s body. $n smiles."), TRUE, ch, nil, nil, TO_ROOM)
				null_affect(ch, AFF_POISON)
				return
			}
		} else {
			if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<0)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("cure who?\r\n"))
				return
			}
			if can_kill(ch, vict, nil, 2) == 0 {
				return
			} else {
				if ch == vict {
					send_to_char(ch, libc.CString("Use focus %s, not focus %s %s.\r\n"), &arg[0], &arg[0], GET_NAME(vict))
					return
				}
				if !AFF_FLAGGED(vict, AFF_POISON) {
					send_to_char(ch, libc.CString("They are not poisoned!\r\n"))
					return
				} else if ch.Mana < ch.Max_mana/20 {
					send_to_char(ch, libc.CString("You do not have enough ki to cure.\r\n"))
					return
				} else if GET_SKILL(ch, SKILL_CURE) < axion_dice(0) {
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki and aim a pulsing light at $N's body. Nothing happens."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki and aims a pulsing light at your body. You are STILL poisoned!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki and aims a pulsing light at $N's body. $N looks disappointed."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					return
				} else {
					affect_from_char(vict, SPELL_POISON)
					ch.Mana -= ch.Max_mana / 20
					reveal_hiding(ch, 0)
					act(libc.CString("You focus ki and aim a pulsing light at $N's body. $e is cured."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n focuses ki and aims a pulsing light at your body. You have been cured of your poison!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$n focuses ki and aims a pulsing light at $N's body. $N smiles."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					null_affect(vict, AFF_POISON)
					return
				}
			}
		}
	} else if C.strcmp(&arg[0], libc.CString("poison")) == 0 {
		if know_skill(ch, SKILL_POISON) == 0 {
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Poison who?\r\n"))
			return
		}
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		} else {
			if ch == vict {
				send_to_char(ch, libc.CString("Why poison yourself?\r\n"))
				return
			}
			if IS_NPC(vict) {
				if MOB_FLAGGED(vict, MOB_NOPOISON) {
					send_to_char(ch, libc.CString("You get the feeling that this being is immune to poison.\r\n"))
					return
				}
			}
			if AFF_FLAGGED(vict, AFF_POISON) {
				send_to_char(ch, libc.CString("They already have been poisoned!\r\n"))
				return
			} else if ch.Mana < ch.Max_mana/20 {
				send_to_char(ch, libc.CString("You do not have enough ki to poison.\r\n"))
				return
			} else if GET_SKILL(ch, SKILL_POISON) < axion_dice(0) {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki and fling poison at $N. You missed!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n focuses ki and flings poison at you, but misses!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n focuses ki and flings poison at $N, but misses!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				return
			} else {
				ch.Mana -= ch.Max_mana / 20
				reveal_hiding(ch, 0)
				act(libc.CString("You focus ki and fling poison at $N! The poison burns into $s skin!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n focuses ki and flings poison at you! The poison burns into your skin!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n focuses ki and flings poison at $N! The poison burns into $s skin!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if IS_NPC(vict) {
					set_fighting(vict, ch)
				}
				if vict.Race == RACE_MUTANT && ((vict.Genome[0]) == 7 || (vict.Genome[1]) == 7) {
					act(libc.CString("However $N seems unaffected by the poison."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("Your natural immunity to poison prevents it from affecting you."), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("However $N seems unaffected by the poison."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else {
					vict.Poisonby = ch
					if ch.Charge > 0 {
						send_to_char(ch, libc.CString("You lose your concentration and release your charged ki!\r\n"))
						do_charge(ch, libc.CString("release"), 0, 0)
					}
					var duration int = int(ch.Aff_abils.Intel / 20)
					assign_affect(vict, AFF_POISON, SKILL_POISON, duration, 0, 0, 0, 0, 0, 0)
				}
				return
			}
		}
	} else {
		send_to_char(ch, libc.CString("What do you want to focus?\r\n"))
		return
	}
}
func do_kaioken(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		roll  int   = axion_dice(0)
		x     int   = 0
		pass  int   = FALSE
		boost int64 = 0
	)
	one_argument(argument, &arg[0])
	if check_skill(ch, SKILL_KAIOKEN) == 0 {
		return
	}
	if ch.Alignment <= -50 {
		send_to_char(ch, libc.CString("Your heart is too corrupt to use that technique!\r\n"))
		return
	}
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_HEALT) {
			send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
			return
		}
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What level of kaioken do you want to try and achieve?\r\nSyntax: kaioken 1-20\r\n"))
		return
	}
	if ch.Kaioken > 0 {
		send_to_char(ch, libc.CString("You drop out of kaioken.\r\n"))
		act(libc.CString("$n@w drops out of kaioken.@n"), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Hit-(gear_pl(ch)/10)*int64(ch.Kaioken) > 0 {
			ch.Hit -= (gear_pl(ch) / 10) * int64(ch.Kaioken)
		} else {
			ch.Hit = 1
		}
		ch.Kaioken = 0
		return
	}
	x = libc.Atoi(libc.GoString(&arg[0]))
	if x <= 0 || x > 20 {
		send_to_char(ch, libc.CString("That level of kaioken dosn't exist...\r\nSyntax: kaioken 1-20\r\n"))
		return
	}
	if !IS_NPC(ch) {
		if (IS_TRANSFORMED(ch) || ch.Race == RACE_HOSHIJIN && ch.Starphase > 0) && x > 5 {
			send_to_char(ch, libc.CString("You can not manage a kaioken level higher than 5 when transformed.\r\n"))
			return
		}
	}
	if ch.Mana < ((ch.Max_mana / 50) * int64(x)) {
		send_to_char(ch, libc.CString("You do not have enough ki to focus into your body for that level.\r\n"))
		return
	}
	var xnum int = (x * 5) + 1
	roll = rand_number(1, xnum)
	reveal_hiding(ch, 0)
	if init_skill(ch, SKILL_KAIOKEN) < roll {
		send_to_char(ch, libc.CString("You try to focus your ki into your body but mess up somehow.\r\n"))
		act(libc.CString("$n tries to use kaioken but messes up somehow."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= (ch.Max_mana / 50) * int64(x)
		improve_skill(ch, SKILL_KAIOKEN, 1)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		return
	}
	switch x {
	case 1:
		fallthrough
	case 2:
		pass = TRUE
	case 3:
		if ch.Max_hit < 5000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 4:
		if ch.Max_hit < 10000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 5:
		if ch.Max_hit < 15000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 6:
		if ch.Max_hit < 25000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 7:
		if ch.Max_hit < 35000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 8:
		if ch.Max_hit < 50000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 9:
		if ch.Max_hit < 75000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 10:
		if ch.Max_hit < 100000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 11:
		if ch.Max_hit < 150000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 12:
		if ch.Max_hit < 200000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 13:
		if ch.Max_hit < 250000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 14:
		if ch.Max_hit < 300000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 15:
		if ch.Max_hit < 400000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 16:
		if ch.Max_hit < 500000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 17:
		if ch.Max_hit < 600000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 18:
		if ch.Max_hit < 700000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 19:
		if ch.Max_hit < 800000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	case 20:
		if ch.Max_hit < 1000000 {
			act(libc.CString("@rA blazing red aura bursts up around your body, flashing intensely before your body gives out and you release the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@rA blazing red aura bursts up around @R$n's @rbody, flashing intensely before $s body gives out and $e releases the kaioken because of the pressure!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			pass = TRUE
		}
	default:
		pass = TRUE
	}
	if pass == FALSE {
		improve_skill(ch, SKILL_KAIOKEN, 1)
		ch.Mana -= (ch.Max_mana / 50) * int64(x)
		return
	}
	boost = (gear_pl(ch) / 10) * int64(x)
	if ch.Hit > gear_pl(ch) {
		ch.Hit = gear_pl(ch)
	}
	ch.Hit += boost
	ch.Mana -= (ch.Max_mana / 50) * int64(x)
	ch.Kaioken = x
	send_to_char(ch, libc.CString("@rA dark red aura bursts up around your body as you achieve Kaioken x %d!@n\r\n"), ch.Kaioken)
	act(libc.CString("@rA dark red aura bursts up around @R$n@r as they achieve a level of Kaioken!@n"), TRUE, ch, nil, nil, TO_ROOM)
	improve_skill(ch, SKILL_KAIOKEN, 1)
	ch.Act[int(PLR_POWERUP/32)] &= bitvector_t(^(1 << (int(PLR_POWERUP % 32))))
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
}
func do_plant(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict      *char_data
		obj       *obj_data
		vict_name [100]byte
		obj_name  [100]byte
		roll      int = 0
		detect    int = 0
		fail      int = 0
	)
	if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
		send_to_char(ch, libc.CString("This room just has such a peaceful, easy feeling...\r\n"))
		return
	}
	two_arguments(argument, &obj_name[0], &vict_name[0])
	if (func() *char_data {
		vict = get_char_vis(ch, &vict_name[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Plant what on who?\r\n"))
		return
	} else if vict == ch {
		send_to_char(ch, libc.CString("Come on now, that's rather stupid!\r\n"))
		return
	}
	if MOB_FLAGGED(vict, MOB_NOKILL) && ch.Admlevel == ADMLVL_NONE {
		send_to_char(ch, libc.CString("That isn't such a good idea...\r\n"))
		return
	}
	roll = roll_skill(ch, SKILL_SLEIGHT_OF_HAND) + rand_number(1, 3)
	fail = rand_number(1, 105)
	if (ch.Feats[FEAT_DEFT_HANDS]) != 0 {
		roll += 2
	}
	if vict.Position < POS_SLEEPING {
		detect = 0
	} else {
		detect = roll_skill(vict, SKILL_SPOT) + rand_number(1, 3)
	}
	if (ADM_FLAGGED(vict, ADM_NOSTEAL) || libc.FuncAddr(GET_MOB_SPEC(vict)) == libc.FuncAddr(shop_keeper)) && ch.Admlevel < 5 {
		roll = -10
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &obj_name[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You don't have that to plant on them.\r\n"))
		return
	}
	if roll <= detect && roll <= fail {
		reveal_hiding(ch, 0)
		act(libc.CString("@C$n@w tries to plant $p@w on you!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@w tries to plant $p@w on @c$N@w!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		act(libc.CString("@wYou try and fail to plant $p@w on @c$N@w, and $E notices!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else if roll <= fail {
		act(libc.CString("@wYou try and fail to plant $p@w on @c$N@w! However no one seemed to notice.@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else if obj.Weight+int64(gear_weight(vict)) > max_carry_weight(vict) {
		reveal_hiding(ch, 0)
		act(libc.CString("@C$n@w tries to plant $p@w on you!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@w tries to plant $p@w on @c$N@w!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		act(libc.CString("@wYou try and fail to plant $p@w on @c$N@w because $E can't carry the weight. It seems $E noticed the attempt!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	} else if roll <= detect {
		act(libc.CString("@cYou feel like the weight of your inventory has changed.@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@c$N@w looks around after feeling $S pockets.@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		act(libc.CString("@wYou plant $p@w on @c$N@w! @c$N @wseems to notice the change in weight in their inventory.@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		obj_from_char(obj)
		obj_to_char(obj, vict)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		act(libc.CString("@wYou plant $p@w on @c$N@w! No one noticed, whew....@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
		obj_from_char(obj)
		obj_to_char(obj, vict)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
}
func do_forgery(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj2     *obj_data
		obj3     *obj_data = nil
		obj      *obj_data
		obj4     *obj_data = nil
		next_obj *obj_data
		found    int = FALSE
		arg      [2048]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_FORGERY) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Okay, make a forgery of what?\r\n"))
		return
	}
	if (func() *obj_data {
		obj2 = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj2
	}()) == nil {
		send_to_char(ch, libc.CString("You want to make a fake copy of what?\r\n"))
		return
	}
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if found == FALSE && GET_OBJ_VNUM(obj) == 19 && (!OBJ_FLAGGED(obj, ITEM_BROKEN) && !OBJ_FLAGGED(obj, ITEM_FORGED)) {
			found = TRUE
			obj4 = obj
		}
	}
	if found == FALSE || obj4 == nil {
		send_to_char(ch, libc.CString("You need a forgery kit.\r\n"))
		return
	}
	if GET_OBJ_VNUM(obj2) == 19 {
		send_to_char(ch, libc.CString("You can't duplicate a forgery kit.\r\n"))
		return
	}
	if OBJ_FLAGGED(obj2, ITEM_FORGED) {
		send_to_char(ch, libc.CString("%s is forgery, there is no reason to make a fake of a fake!\r\n"), obj2.Short_description)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	if OBJ_FLAGGED(obj2, ITEM_BROKEN) {
		send_to_char(ch, libc.CString("%s is broken, there is no reason to make a fake of this mess!\r\n"), obj2.Short_description)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	if GET_OBJ_VNUM(obj2) >= 60000 || GET_OBJ_VNUM(obj2) == 0 {
		send_to_char(ch, libc.CString("You can not make a forgery of that! It's far too squishy...."))
		return
	}
	if GET_OBJ_VNUM(obj2) >= 18800 && GET_OBJ_VNUM(obj2) <= 0x4A37 {
		send_to_char(ch, libc.CString("You can not make a forgery of that!\r\n"))
		return
	}
	if GET_OBJ_VNUM(obj2) >= 19080 && GET_OBJ_VNUM(obj2) <= 0x4AFF {
		send_to_char(ch, libc.CString("You can not make a forgery of that!\r\n"))
		return
	}
	if GET_OBJ_VNUM(obj2) >= 4 && GET_OBJ_VNUM(obj2) <= 6 {
		send_to_char(ch, libc.CString("You can not make a forgery of that!\r\n"))
		return
	}
	if OBJ_FLAGGED(obj2, ITEM_PROTECTED) {
		send_to_char(ch, libc.CString("You don't know where to begin with this work of ART.\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	act(libc.CString("@c$n@w looks at $p, begins to work on forging a fake copy of it.@n"), TRUE, ch, obj2, nil, TO_ROOM)
	improve_skill(ch, SKILL_FORGERY, 1)
	if GET_SKILL(ch, SKILL_FORGERY) < axion_dice(0) {
		if rand_number(1, 10) >= 9 {
			send_to_char(ch, libc.CString("In the middle of creating a forgery of %s you screw up. The fabrication unit built into the forgery kit melts and bonds with the original. You clumsy mistake with the Estex Titanium drill has broken both.\r\n"), obj2.Short_description)
			extract_obj(obj4)
			extract_obj(obj2)
			return
		}
		send_to_char(ch, libc.CString("You start to make a forgery of %s but screw up and waste your forgery kit..\r\n"), obj2.Short_description)
		act(libc.CString("@c$n@w tried to duplicate $p but screws up somehow.@n"), TRUE, ch, obj2, nil, TO_ROOM)
		obj_from_char(obj4)
		extract_obj(obj4)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	var loadn int = int(GET_OBJ_VNUM(obj2))
	obj3 = read_object(obj_vnum(loadn), VIRTUAL)
	add_unique_id(obj3)
	obj_to_char(obj3, ch)
	obj3.Extra_flags[int(ITEM_FORGED/32)] |= bitvector_t(1 << (int(ITEM_FORGED % 32)))
	obj3.Weight = int64(rand_number(int(obj3.Weight/2), int(obj3.Weight)))
	obj_from_char(obj4)
	extract_obj(obj4)
	send_to_char(ch, libc.CString("You make an excellent forgery of %s@n!\r\n"), obj2.Short_description)
	act(libc.CString("@c$n@w makes a perfect forgery of $p.@n"), TRUE, ch, obj2, nil, TO_ROOM)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
}
func do_appraise(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i     int
		found int
		obj   *obj_data
		arg   [2048]byte
		buf   [64936]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if know_skill(ch, SKILL_APPRAISE) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Okay, appraise what?\r\n"))
		return
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("You want to appraise what?\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	act(libc.CString("@c$n@w looks at $p, turning it over in $s hands.@n"), TRUE, ch, obj, nil, TO_ROOM)
	improve_skill(ch, SKILL_APPRAISE, 1)
	if GET_SKILL(ch, SKILL_APPRAISE) < axion_dice(-10) {
		send_to_char(ch, libc.CString("You fail to perceive the worth of %s..\r\n"), obj.Short_description)
		act(libc.CString("@c$n@w looks stumped about $p.@n"), TRUE, ch, obj, nil, TO_ROOM)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	if OBJ_FLAGGED(obj, ITEM_BROKEN) {
		send_to_char(ch, libc.CString("%s is broken!\r\n"), obj.Short_description)
		act(libc.CString("@c$n@w looks at $p and frowns.@n"), TRUE, ch, obj, nil, TO_ROOM)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	if OBJ_FLAGGED(obj, ITEM_FORGED) {
		send_to_char(ch, libc.CString("%s is fake and worthless!\r\n"), obj.Short_description)
		act(libc.CString("@c$n@w looks at $p with an angry face.@n"), TRUE, ch, obj, nil, TO_ROOM)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	found = FALSE
	var displevel int = obj.Level
	if obj.Type_flag == ITEM_WEAPON && OBJ_FLAGGED(obj, ITEM_CUSTOM) {
		displevel = 20
	}
	send_to_char(ch, libc.CString("%s is worth: %s\r\nMin Lvl: %d\r\n"), obj.Short_description, add_commas(int64(obj.Cost)), displevel)
	if obj.Type_flag == ITEM_WEAPON {
		if OBJ_FLAGGED(obj, ITEM_WEAPLVL1) {
			send_to_char(ch, libc.CString("Weapon Level: 1\nDamage Bonus: 5%s\r\n"), "%")
		} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL2) {
			send_to_char(ch, libc.CString("Weapon Level: 2\nDamage Bonus: 10%s\r\n"), "%")
		} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL3) {
			send_to_char(ch, libc.CString("Weapon Level: 3\nDamage Bonus: 20%s\r\n"), "%")
		} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL4) {
			send_to_char(ch, libc.CString("Weapon Level: 4\nDamage Bonus: 30%s\r\n"), "%")
		} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL5) {
			send_to_char(ch, libc.CString("Weapon Level: 5\nDamage Bonus: 50%s\r\n"), "%")
		}
	}
	send_to_char(ch, libc.CString("Size: %s\r\n"), size_names[obj.Size])
	if OBJ_FLAGGED(obj, ITEM_SLOT1) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("Token Slots  : @m0/1@n\n"))
	} else if OBJ_FLAGGED(obj, ITEM_SLOT1) && OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("Token Slots  : @m1/1@n\n"))
	} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && !OBJ_FLAGGED(obj, ITEM_SLOT_ONE) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("Token Slots  : @m0/2@n\n"))
	} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && OBJ_FLAGGED(obj, ITEM_SLOT_ONE) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("Token Slots  : @m1/2@n\n"))
	} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
		send_to_char(ch, libc.CString("Token Slots  : @m2/2@n\n"))
	}
	send_to_char(ch, libc.CString("Bonuses:"))
	act(libc.CString("@c$n@w looks at $p and nods, a satisfied look on $s face.@n"), TRUE, ch, obj, nil, TO_ROOM)
	var percent int = FALSE
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Modifier != 0 {
			if obj.Affected[i].Location == APPLY_REGEN || obj.Affected[i].Location == APPLY_TRAIN || obj.Affected[i].Location == APPLY_LIFEMAX {
				percent = TRUE
			}
			sprinttype(obj.Affected[i].Location, apply_types[:], &buf[0], uint64(64936))
			send_to_char(ch, libc.CString("%s %+d%s to %s"), func() string {
				if func() int {
					p := &found
					x := *p
					*p++
					return x
				}() != 0 {
					return ","
				}
				return ""
			}(), obj.Affected[i].Modifier, func() string {
				if percent == TRUE {
					return "%"
				}
				return ""
			}(), &buf[0])
			percent = FALSE
			switch obj.Affected[i].Location {
			case APPLY_FEAT:
				send_to_char(ch, libc.CString(" (%s)"), feat_list[obj.Affected[i].Specific].Name)
			case APPLY_SKILL:
				send_to_char(ch, libc.CString(" (%s)"), spell_info[obj.Affected[i].Specific].Name)
			}
		}
	}
	if found == 0 {
		send_to_char(ch, libc.CString(" None"))
	}
	var buf2 [64936]byte
	sprintbitarray(obj.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &buf2[0])
	send_to_char(ch, libc.CString("\nSpecial: %s\r\n"), &buf2[0])
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
}
func do_disguise(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		skill int = 0
		roll  int = 0
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You forgot your disguise off in mobland.\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_DISGUISED) {
		send_to_char(ch, libc.CString("You stop disguising yourself.\r\n"))
		ch.Act[int(PLR_DISGUISED/32)] &= bitvector_t(^(1 << (int(PLR_DISGUISED % 32))))
		act(libc.CString("@C$n @wpulls off $s disguise and reveals $mself!"), TRUE, ch, nil, nil, TO_ROOM)
		return
	}
	if know_skill(ch, SKILL_DISGUISE) == 0 {
		return
	}
	if (ch.Equipment[WEAR_HEAD]) == nil {
		send_to_char(ch, libc.CString("You can't disguise your identity without anything on your head.\r\n"))
		return
	}
	if ch.Move < ch.Move/50 {
		send_to_char(ch, libc.CString("You are too tired to try that right now.\r\n"))
		return
	}
	skill = GET_SKILL(ch, SKILL_DISGUISE)
	roll = axion_dice(-10)
	if skill > roll {
		send_to_char(ch, libc.CString("You managed to disguise yourself with some skilled manipulation of your headwear.\r\n"))
		act(libc.CString("@C$n @wmanages to disguise $mself with some skilled manipulation of $s headwear."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_DISGUISED/32)] |= bitvector_t(1 << (int(PLR_DISGUISED % 32)))
		return
	} else {
		send_to_char(ch, libc.CString("You finish attempting to disguise yourself, but realize you failed and need to try again.\r\n"))
		act(libc.CString("@C$n @wattempts and fails to disguise $mself properly and must try again."), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Move >= ch.Move/50 {
			ch.Move -= ch.Move / 50
		} else {
			ch.Move = 0
		}
		return
	}
}
func do_eavesdrop(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		dir int
		buf [100]byte
	)
	one_argument(argument, &buf[0])
	if ch.Listenroom > 0 {
		send_to_char(ch, libc.CString("You stop eavesdropping.\r\n"))
		ch.Listenroom = room_vnum(real_room(0))
		ch.Eavesdir = -1
		return
	}
	if buf[0] == 0 {
		send_to_char(ch, libc.CString("In which direction would you like to eavesdrop?\r\n"))
		return
	}
	if (func() int {
		dir = search_block(&buf[0], &dirs[0], FALSE)
		return dir
	}()) < 0 {
		send_to_char(ch, libc.CString("Which directions is that?\r\n"))
		return
	}
	if know_skill(ch, SKILL_EAVESDROP) == 0 {
		return
	}
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) != nil {
		if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Exit_info&(1<<1)) != 0 && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword != nil {
			stdio.Sprintf(&buf[0], "The %s is closed.\r\n", fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword))
			send_to_char(ch, &buf[0])
		} else {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
				ch.Listenroom = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Number
			} else {
				ch.Listenroom = -1
			}
			ch.Eavesdir = dir
			send_to_char(ch, libc.CString("Okay.\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("There is not a room there...\r\n"))
	}
}
func do_zanzoken(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob int   = 0
		perc int   = 0
		cost int64 = 0
	)
	if know_skill(ch, SKILL_ZANZOKEN) == 0 && !IS_NPC(ch) {
		return
	}
	if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
		ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
		send_to_char(ch, libc.CString("You release the ki you had prepared for a zanzoken.\r\n"))
		return
	}
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are busy in a grapple!\r\n"))
		return
	}
	if !IS_NPC(ch) {
		prob = GET_SKILL(ch, SKILL_ZANZOKEN)
	} else {
		prob = rand_number(80, 90)
	}
	perc = axion_dice(0)
	cost = ch.Max_mana / 50
	if prob > 75 {
		cost *= 2
	} else if prob > 50 {
		cost *= 4
	} else if prob >= 25 {
		cost *= 8
	} else if prob < 25 {
		cost *= 10
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	if prob < perc {
		send_to_char(ch, libc.CString("You focus your ki in preparation of a zanzoken but mess up and waste your ki!\r\n"))
		improve_skill(ch, SKILL_ZANZOKEN, 2)
		ch.Mana -= cost
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	}
	act(libc.CString("@wYou focus your ki, preparing to move at super speeds if necessary.@n"), TRUE, ch, nil, nil, TO_CHAR)
	ch.Mana -= cost
	ch.Affected_by[int(AFF_ZANZOKEN/32)] |= 1 << (int(AFF_ZANZOKEN % 32))
	improve_skill(ch, SKILL_ZANZOKEN, 2)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
}
func do_block(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict *char_data
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if arg[0] == 0 {
		if ch.Blocks == nil {
			send_to_char(ch, libc.CString("You want to block who?\r\n"))
			return
		}
		if ch.Blocks != nil {
			act(libc.CString("@wYou stop blocking @c$N@w.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_CHAR)
			act(libc.CString("@C$n@w stops blocking you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_VICT)
			act(libc.CString("@C$n@w stops blocking @c$N@w.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_NOTVICT)
			vict = ch.Blocks
			vict.Blocked = nil
			ch.Blocks = nil
			return
		}
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("You do not see the target here.\r\n"))
		return
	}
	if ch.Blocks == vict {
		send_to_char(ch, libc.CString("They are already blocked by you!\r\n"))
		return
	}
	if ch == vict {
		send_to_char(ch, libc.CString("You can't block yourself, are you mental?\r\n"))
		return
	}
	if vict.Blocked != nil {
		send_to_char(ch, libc.CString("They are already blocked by someone else!\r\n"))
		return
	}
	if ch.Blocks != nil {
		act(libc.CString("@wYou stop blocking @c$N@w.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_CHAR)
		act(libc.CString("@C$n@w stops blocking you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_VICT)
		act(libc.CString("@C$n@w stops blocking @c$N@w.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_NOTVICT)
		var oldv *char_data = ch.Blocks
		oldv.Blocked = nil
		ch.Blocks = vict
		vict.Blocked = ch
		reveal_hiding(ch, 0)
		act(libc.CString("@wYou start blocking @c$N's@w escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_CHAR)
		act(libc.CString("@C$n@w starts blocking your escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_VICT)
		act(libc.CString("@C$n@w starts blocking @c$N's@w escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_NOTVICT)
		return
	} else {
		ch.Blocks = vict
		vict.Blocked = ch
		reveal_hiding(ch, 0)
		act(libc.CString("@wYou start blocking @c$N's@w escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_CHAR)
		act(libc.CString("@C$n@w starts blocking your escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_VICT)
		act(libc.CString("@C$n@w starts blocking @c$N's@w escape.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Blocks), TO_NOTVICT)
		return
	}
}
func do_eyec(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if PLR_FLAGGED(ch, PLR_EYEC) {
		ch.Act[int(PLR_EYEC/32)] &= bitvector_t(^(1 << (int(PLR_EYEC % 32))))
		act(libc.CString("@wYou open your eyes.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@w opens $s eyes.@n"), TRUE, ch, nil, nil, TO_ROOM)
	} else if !PLR_FLAGGED(ch, PLR_EYEC) {
		ch.Act[int(PLR_EYEC/32)] |= bitvector_t(1 << (int(PLR_EYEC % 32)))
		act(libc.CString("@wYou close your eyes.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@w closes $s eyes.@n"), TRUE, ch, nil, nil, TO_ROOM)
	}
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
}
func do_solar(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict   *char_data = nil
		next_v *char_data = nil
		prob   int        = 0
		perc   int        = 0
		cost   int        = 0
		bonus  int        = 0
	)
	if know_skill(ch, SKILL_SOLARF) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	prob = GET_SKILL(ch, SKILL_SOLARF)
	perc = rand_number(0, 101)
	if prob >= 75 {
		cost = int(ch.Max_mana / 50)
	} else if prob >= 50 {
		cost = int(ch.Max_mana / 25)
	} else if prob >= 25 {
		cost = int(ch.Max_mana / 20)
	} else if prob < 25 {
		cost = int(ch.Max_mana / 15)
	}
	if ch.Mana < int64(cost) {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	bonus = int(ch.Aff_abils.Intel / 3)
	prob += bonus
	if prob < perc {
		act(libc.CString("@WYou raise both your hands to either side of your face, while closing your eyes, and shout '@YSolar Flare@W' but nothing happens!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W raises both $s hands to either side of $s face, while closing $s eyes, and shouts '@YSolar Flare@W' but nothing happens!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Mana -= int64(cost)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		improve_skill(ch, SKILL_SOLARF, 0)
		return
	}
	act(libc.CString("@WYou raise both your hands to either side of your face, while closing your eyes, and shout '@YSolar Flare@W' as a blinding light fills the area!@n"), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("@C$n@W raises both $s hands to either side of $s face, while closing $s eyes, and shouts '@YSolar Flare@W' as a blinding light fills the area!@n"), TRUE, ch, nil, nil, TO_ROOM)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		} else if PLR_FLAGGED(vict, PLR_EYEC) {
			continue
		} else if AFF_FLAGGED(vict, AFF_BLIND) {
			continue
		} else if vict.Position == POS_SLEEPING {
			continue
		} else {
			var duration int = 1
			assign_affect(vict, AFF_BLIND, SKILL_SOLARF, duration, 0, 0, 0, 0, 0, 0)
			act(libc.CString("@W$N@W is @YBLINDED@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RYou are @YBLINDED@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@W$N@W is @YBLINDED@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		}
	}
	improve_skill(ch, SKILL_SOLARF, 0)
	ch.Mana -= int64(cost)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
}
func do_heal(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		cost  int64 = 0
		prob  int64 = 0
		perc  int64 = 0
		heal  int64 = 0
		bonus int64 = 0
		vict  *char_data
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if check_skill(ch, SKILL_HEAL) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("You want to heal WHO?\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("You do not see the target here.\r\n"))
		return
	}
	prob = int64(init_skill(ch, SKILL_HEAL))
	perc = int64(rand_number(0, 110))
	if prob >= 100 {
		cost = ch.Max_mana / 20
		heal = vict.Max_hit / 5
	} else if prob >= 90 {
		cost = ch.Max_mana / 16
		heal = vict.Max_hit / 10
	} else if prob >= 75 {
		cost = ch.Max_mana / 14
		heal = vict.Max_hit / 12
	} else if prob >= 50 {
		cost = ch.Max_mana / 12
		heal = vict.Max_hit / 15
	} else if prob >= 25 {
		cost = ch.Max_mana / 10
		heal = vict.Max_hit / 20
	} else if prob < 25 {
		cost = ch.Max_mana / 6
		heal = vict.Max_hit / 20
	}
	if (ch.Bonuses[BONUS_HEALER]) > 0 {
		heal += int64(float64(heal) * 0.1)
	}
	if heal < gear_pl(vict) {
		heal += (heal / 100) * int64(ch.Aff_abils.Wis/4)
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You do not have enough ki.\r\n"))
		return
	}
	if vict.Hit >= gear_pl(vict) {
		if vict != ch {
			send_to_char(ch, libc.CString("They are already at full health.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You are already at full health.\r\n"))
		}
		return
	}
	if vict.Suppression > 0 && vict.Hit >= ((gear_pl(vict)/100)*vict.Suppression) {
		send_to_char(ch, libc.CString("They are already at full health.\r\n"))
		return
	}
	bonus = int64((ch.Aff_abils.Intel / 2) + ch.Aff_abils.Wis/3)
	prob += bonus
	if prob < perc {
		if vict != ch {
			act(libc.CString("@WYou place your hands near @c$N@W, but fail to concentrate enough to heal them!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W places $s hands near you, but nothing happens!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$n@W places $s hands near @c$N@W, but nothing happens."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Mana -= cost
			improve_skill(ch, SKILL_HEAL, 0)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		}
		if vict == ch {
			act(libc.CString("@WYou place your hands on your body, but fail to concentrate to heal yourself!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W places $s hands on $s body, but nothing happens."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			ch.Mana -= cost
			improve_skill(ch, SKILL_HEAL, 0)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		}
	}
	if vict != ch {
		if (ch.Bonuses[BONUS_HEALER]) > 0 {
			heal += int64(float64(heal) * 0.25)
		}
		act(libc.CString("@WYou place your hands near @c$N@W and an orange glow surrounds $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W places $s hands near you and an orange glow surrounds you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n@W places $s hands near @c$N@W and an orange glow surrounds $M."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= cost
		if vict.Hit+heal <= gear_pl(vict) {
			vict.Hit += heal
		} else {
			vict.Hit = gear_pl(vict)
		}
		if ch.Chclass == CLASS_NAIL {
			if GET_SKILL(ch, SKILL_HEAL) >= 100 {
				ch.Move += int64(float64(heal) * 0.4)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			} else if GET_SKILL(ch, SKILL_HEAL) >= 60 {
				ch.Move += int64(float64(heal) * 0.2)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			} else if GET_SKILL(ch, SKILL_HEAL) >= 40 {
				ch.Move += int64(float64(heal) * 0.1)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			}
		}
		if vict.Suppression > 0 && vict.Hit > ((gear_pl(vict)/100)*vict.Suppression) {
			vict.Hit = (gear_pl(vict) / 100) * vict.Suppression
			send_to_char(vict, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
		}
		null_affect(ch, AFF_POISON)
		null_affect(ch, AFF_BLIND)
		if AFF_FLAGGED(vict, AFF_BURNED) {
			send_to_char(vict, libc.CString("Your burns are healed now.\r\n"))
			act(libc.CString("$n@w's burns are now healed.@n"), TRUE, vict, nil, nil, TO_ROOM)
			vict.Affected_by[int(AFF_BURNED/32)] &= ^(1 << (int(AFF_BURNED % 32)))
		}
		if AFF_FLAGGED(vict, AFF_HYDROZAP) {
			send_to_char(vict, libc.CString("You no longer feel a great thirst.\r\n"))
			act(libc.CString("$n@w no longer looks as if they could drink an ocean.@n"), TRUE, vict, nil, nil, TO_ROOM)
			vict.Affected_by[int(AFF_HYDROZAP/32)] &= ^(1 << (int(AFF_HYDROZAP % 32)))
		}
		vict.Limb_condition[0] = 100
		vict.Limb_condition[1] = 100
		vict.Limb_condition[2] = 100
		vict.Limb_condition[3] = 100
		if float64(vict.Lifeforce) <= float64(GET_LIFEMAX(vict))*0.5 && vict.Race != RACE_ANDROID {
			vict.Lifeforce += int64(float64(GET_LIFEMAX(ch)) * 0.35)
			if vict.Lifeforce > int64(GET_LIFEMAX(ch)) {
				vict.Lifeforce = int64(GET_LIFEMAX(ch))
			}
			send_to_char(vict, libc.CString("You feel that your lifeforce has recovered some!\r\n"))
		}
		improve_skill(ch, SKILL_HEAL, 0)
		if vict.Master == ch || ch.Master == vict || ch.Master == vict.Master {
			if ch.Chclass == CLASS_NAIL && ch.Race == RACE_NAMEK && level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 && float64(vict.Hit) <= float64(gear_pl(vict))*0.85 && rand_number(1, 3) == 3 {
				ch.Exp += int64(float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.005)
			}
		}
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	}
	if vict == ch {
		if (ch.Bonuses[BONUS_HEALER]) > 0 {
			heal += int64(float64(heal) * 0.25)
		}
		act(libc.CString("@WYou place your hands on your body and an orange glow surrounds you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n@W places $s hands on $s body and an orange glow surrounds $m."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Mana -= cost
		if vict.Hit+heal <= gear_pl(vict) {
			vict.Hit += heal
		} else {
			vict.Hit = gear_pl(vict)
		}
		if vict.Suppression > 0 && vict.Hit > ((gear_pl(vict)/100)*vict.Suppression) {
			vict.Hit = (gear_pl(vict) / 100) * vict.Suppression
			send_to_char(vict, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
		}
		if ch.Chclass == CLASS_NAIL {
			if GET_SKILL(ch, SKILL_HEAL) >= 100 {
				ch.Move += int64(float64(heal) * 0.4)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			} else if GET_SKILL(ch, SKILL_HEAL) >= 60 {
				ch.Move += int64(float64(heal) * 0.2)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			} else if GET_SKILL(ch, SKILL_HEAL) >= 40 {
				ch.Move += int64(float64(heal) * 0.1)
				send_to_char(vict, libc.CString("@GYou feel some of your stamina return as well!@n\r\n"))
			}
		}
		vict.Affected_by[int(AFF_BLIND/32)] &= ^(1 << (int(AFF_BLIND % 32)))
		vict.Limb_condition[0] = 100
		vict.Limb_condition[1] = 100
		vict.Limb_condition[2] = 100
		vict.Limb_condition[3] = 100
		if !PLR_FLAGGED(vict, PLR_TAIL) && (vict.Race == RACE_BIO || vict.Race == RACE_ICER) {
			vict.Act[int(PLR_TAIL/32)] |= bitvector_t(1 << (int(PLR_TAIL % 32)))
		}
		if !PLR_FLAGGED(vict, PLR_STAIL) && (vict.Race == RACE_SAIYAN || vict.Race == RACE_HALFBREED) {
			vict.Act[int(PLR_STAIL/32)] |= bitvector_t(1 << (int(PLR_STAIL % 32)))
		}
		improve_skill(ch, SKILL_HEAL, 0)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	}
	return
}
func do_barrier(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob int = 0
		perc int = 0
		size int = 0
		arg  [2048]byte
	)
	one_argument(argument, &arg[0])
	if know_skill(ch, SKILL_BARRIER) == 0 && GET_SKILL(ch, SKILL_AQUA_BARRIER) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("[Syntax] barrier < 1-75 | release >\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SANCTUARY) && C.strcasecmp(libc.CString("release"), &arg[0]) == 0 {
		act(libc.CString("@BYou dispel your barrier, releasing its energy.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@B$n@B dispels $s barrier, releasing its energy.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Barrier = 0
		ch.Affected_by[int(AFF_SANCTUARY/32)] &= ^(1 << (int(AFF_SANCTUARY % 32)))
		return
	} else if C.strcasecmp(libc.CString("release"), &arg[0]) == 0 {
		send_to_char(ch, libc.CString("You don't have a barrier.\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SANCTUARY) {
		send_to_char(ch, libc.CString("You already have a barrier, try releasing it.\r\n"))
		return
	}
	if ch.Con_cooldown > 0 {
		send_to_char(ch, libc.CString("You must wait a short period before concentrating again.\r\n"))
		return
	}
	size = libc.Atoi(libc.GoString(&arg[0]))
	var cost int64 = 0
	prob = 0
	if GET_SKILL(ch, SKILL_BARRIER) != 0 {
		prob = init_skill(ch, SKILL_BARRIER)
	} else {
		prob = GET_SKILL(ch, SKILL_AQUA_BARRIER)
	}
	perc = axion_dice(0)
	cost = int64((float64(ch.Max_mana) * 0.01) * (float64(size) * 0.5))
	if size > prob {
		send_to_char(ch, libc.CString("You can not create a barrier that is stronger than your skill in barrier.\r\n"))
		return
	} else if size < 1 {
		send_to_char(ch, libc.CString("You have to put at least some ki into the barrier!\r\n"))
		return
	} else if size > 75 {
		send_to_char(ch, libc.CString("You can't control a barrier with more than 75 percent!\r\n"))
		return
	} else if ch.Charge < cost {
		send_to_char(ch, libc.CString("You do not have enough ki charged up!\r\n"))
		return
	} else if prob < perc {
		act(libc.CString("@BYou shout as you form a barrier of ki around your body, but you imbalance it and it explodes outward!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@B$n@B shouts as $e forms a barrier of ki around $s body, but it becomes imbalanced and explodes outward!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Charge -= cost
		if GET_SKILL(ch, SKILL_BARRIER) != 0 {
			improve_skill(ch, SKILL_BARRIER, 2)
		} else {
			improve_skill(ch, SKILL_AQUA_BARRIER, 2)
		}
		ch.Con_cooldown = 30
		return
	} else {
		if GET_SKILL(ch, SKILL_BARRIER) != 0 {
			act(libc.CString("@BYou shout as you form a barrier of ki around your body!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@B$n@B shouts as $e forms a barrier of ki around $s body!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			act(libc.CString("@BYou shout as you form a barrier of ki and raging waters around your body!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@B$n@B shouts as $e forms a barrier of ki and raging waters around $s body!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		ch.Barrier = (ch.Max_mana / 100) * int64(size)
		ch.Charge -= cost
		if GET_SKILL(ch, SKILL_BARRIER) != 0 {
			improve_skill(ch, SKILL_BARRIER, 2)
		} else {
			improve_skill(ch, SKILL_AQUA_BARRIER, 2)
		}
		ch.Affected_by[int(AFF_SANCTUARY/32)] |= 1 << (int(AFF_SANCTUARY % 32))
		ch.Con_cooldown = 20
		return
	}
}
func do_instant(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		skill     int        = 0
		perc      int        = 0
		skill_num int        = 0
		location  int        = 0
		cost      int64      = 0
		tar       *char_data = nil
		arg       [2048]byte = func() [2048]byte {
			var t [2048]byte
			copy(t[:], []byte(""))
			return t
		}()
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
		if PRF_FLAGGED(ch, PRF_ARENAWATCH) {
			ch.Player_specials.Pref[int(PRF_ARENAWATCH/32)] &= bitvector_t(^(1 << (int(PRF_ARENAWATCH % 32))))
			ch.Arenawatch = -1
			send_to_char(ch, libc.CString("You stop watching the arena action.\r\n"))
		}
	}
	if know_skill(ch, SKILL_INSTANTT) == 0 {
		return
	} else if GET_SKILL(ch, SKILL_SENSE) == 0 && !PLR_FLAGGED(ch, PLR_SENSEM) {
		send_to_char(ch, libc.CString("You can't sense them to go to there!\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	} else if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 19800 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 0x4DBB {
		send_to_char(ch, libc.CString("@rYou are in a pocket dimension!@n\r\n"))
		return
	} else if ROOM_FLAGGED(ch.In_room, ROOM_RHELL) || ROOM_FLAGGED(ch.In_room, ROOM_AL) || ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
		send_to_char(ch, libc.CString("You can not leave where you are at!\r\n"))
		return
	} else if arg[0] == 0 {
		send_to_char(ch, libc.CString("Who or where do you want to instant transmission to? [target | planet-(planet name)]\r\n"))
		send_to_char(ch, libc.CString("Example: instant goku\nExample 2: instant planet-earth\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_INSTANTT) > 75 {
		cost = ch.Max_mana / 40
	} else if GET_SKILL(ch, SKILL_INSTANTT) > 50 {
		cost = ch.Max_mana / 20
	} else if GET_SKILL(ch, SKILL_INSTANTT) > 25 {
		cost = ch.Max_mana / 15
	} else if GET_SKILL(ch, SKILL_INSTANTT) < 25 {
		cost = ch.Max_mana / 10
	}
	if ch.Mana-cost < 0 {
		send_to_char(ch, libc.CString("You do not have enough ki to instantaneously move.\r\n"))
		return
	}
	perc = axion_dice(0)
	skill = GET_SKILL(ch, SKILL_INSTANTT)
	skill_num = SKILL_INSTANTT
	if C.strcasecmp(&arg[0], libc.CString("planet-earth")) == 0 {
		location = 300
	} else if C.strcasecmp(&arg[0], libc.CString("planet-namek")) == 0 {
		location = 0x27EE
	} else if C.strcasecmp(&arg[0], libc.CString("planet-frigid")) == 0 {
		location = 4017
	} else if C.strcasecmp(&arg[0], libc.CString("planet-vegeta")) == 0 {
		location = 2200
	} else if C.strcasecmp(&arg[0], libc.CString("planet-konack")) == 0 {
		location = 8006
	} else if C.strcasecmp(&arg[0], libc.CString("planet-aether")) == 0 {
		location = 0x2EF8
	} else if (func() *char_data {
		tar = get_char_vis(ch, &arg[0], nil, 1<<1)
		return tar
	}()) == nil {
		send_to_char(ch, libc.CString("@RThat target was not found.@n\r\n"))
		send_to_char(ch, libc.CString("Who or where do you want to instant transmission to? [target | planet-(planet name)]\r\n"))
		send_to_char(ch, libc.CString("Example: instant goku\nExample 2: instant planet-earth\r\n"))
		return
	}
	if skill < perc || ch.Fighting != nil && rand_number(1, 2) <= 1 {
		if tar != nil {
			if tar != ch {
				send_to_char(ch, libc.CString("You prepare to move instantly but mess up the process and waste some of your ki!\r\n"))
				ch.Mana -= cost
				improve_skill(ch, skill_num, 2)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				return
			} else {
				send_to_char(ch, libc.CString("Moving to yourself would be kinda impossible wouldn't it? If not that then it would at least be pointless.\r\n"))
				return
			}
		} else {
			send_to_char(ch, libc.CString("You prepare to move instantly but mess up the process and waste some of your ki!\r\n"))
			ch.Mana -= cost
			improve_skill(ch, skill_num, 2)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		}
	}
	reveal_hiding(ch, 0)
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
	if tar != nil {
		if tar == ch {
			send_to_char(ch, libc.CString("Moving to yourself would be kinda impossible wouldn't it? If not that then it would at least be pointless.\r\n"))
			return
		} else if ch.Grappling != nil && ch.Grappling == tar {
			send_to_char(ch, libc.CString("You are already in the same room with them and are grappling with them!\r\n"))
			return
		} else if read_sense_memory(ch, tar) == 0 {
			send_to_char(ch, libc.CString("You've never sensed them up close so you do not have a good bearing on their ki signal.\r\n"))
			return
		} else if tar.Admlevel > 0 && ch.Admlevel < 1 {
			send_to_char(ch, libc.CString("That immortal prevents you from reaching them.\r\n"))
			return
		} else if tar.Race == RACE_ANDROID || float64(tar.Hit) < (float64(ch.Hit)*0.001)+1 {
			send_to_char(ch, libc.CString("You can't sense them well enough.\r\n"))
			return
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_AL) && ROOM_FLAGGED(tar.In_room, ROOM_AL) {
			send_to_char(ch, libc.CString("They are dead and can't be reached.\r\n"))
			return
		} else if !ROOM_FLAGGED(ch.In_room, ROOM_RHELL) && ROOM_FLAGGED(tar.In_room, ROOM_RHELL) {
			send_to_char(ch, libc.CString("They are dead and can't be reached.\r\n"))
			return
		} else if ROOM_FLAGGED(tar.In_room, ROOM_NOINSTANT) {
			send_to_char(ch, libc.CString("You can not go there as it is a protected area!\r\n"))
			return
		}
		ch.Mana -= cost
		act(libc.CString("@wPlacing two fingers on your forehead you close your eyes and concentrate. Accelerating to such a speed that you move through the molecules of the universe faster than the speed of light. You stop as you arrive at $N@w!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_CHAR)
		act(libc.CString("@w$n@w appears in an instant out of nowhere right next to you!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_VICT)
		act(libc.CString("@w$n@w places two fingers on $s forehead and disappears in an instant!@n"), TRUE, ch, nil, unsafe.Pointer(tar), TO_NOTVICT)
		ch.Act[int(PLR_TRANSMISSION/32)] |= bitvector_t(1 << (int(PLR_TRANSMISSION % 32)))
		handle_teleport(ch, tar, 0)
		improve_skill(ch, skill_num, 2)
	} else {
		ch.Mana -= cost
		act(libc.CString("@wPlacing two fingers on your forehead you close your eyes and concentrate. Accelerating to such a speed that you move faster than light and arrive almost instantly at your destination. Having located the planet by its collective population's ki.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@w$n@w places two fingers on $s forehead and disappears in an instant!@n"), TRUE, ch, nil, nil, TO_NOTVICT)
		handle_teleport(ch, nil, location)
		improve_skill(ch, skill_num, 2)
	}
}
func load_shadow_dragons() {
	var (
		mob   *char_data = nil
		r_num mob_rnum
	)
	if SHADOW_DRAGON1 > 0 {
		r_num = real_mobile(SHADOW_DRAGON1_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON1)))
		mob = nil
	}
	if SHADOW_DRAGON2 > 0 {
		r_num = real_mobile(SHADOW_DRAGON2_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON2)))
		mob = nil
	}
	if SHADOW_DRAGON3 > 0 {
		r_num = real_mobile(SHADOW_DRAGON3_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON3)))
		mob = nil
	}
	if SHADOW_DRAGON4 > 0 {
		r_num = real_mobile(SHADOW_DRAGON4_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON4)))
		mob = nil
	}
	if SHADOW_DRAGON5 > 0 {
		r_num = real_mobile(SHADOW_DRAGON5_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON5)))
		mob = nil
	}
	if SHADOW_DRAGON6 > 0 {
		r_num = real_mobile(SHADOW_DRAGON6_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON6)))
		mob = nil
	}
	if SHADOW_DRAGON7 > 0 {
		r_num = real_mobile(SHADOW_DRAGON7_VNUM)
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON7)))
	}
	save_mud_time(&time_info)
}
func wishSYS() {
	if SHENRON == TRUE {
		if SELFISHMETER < 10 {
			switch DRAGONC {
			case 300:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@WThe dragon balls on the ground begin to glow yellow in slow pulses.@n\r\n"))
				send_to_planet(0, ROOM_EARTH, libc.CString("@DThe sky begins to grow dark and cloudy suddenly.@n\r\n"))
				DRAGONC -= 1
			case 295:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@WSuddenly lightning shoots into the sky, twisting about as a roar can be heard for miles!@n\r\n"))
				send_to_planet(0, ROOM_EARTH, libc.CString("@DThe sky flashes with lightning.@n\r\n"))
				DRAGONC -= 1
			case 290:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@WThe lightning takes shape and slowly the Eternal Dragon, Shenron, can be made out from the glow!@n\r\n"))
				char_from_room(EDRAGON)
				char_to_room(EDRAGON, real_room(room_vnum(DRAGONR)))
				DRAGONC -= 1
			case 285:
				send_to_planet(0, ROOM_EARTH, libc.CString("@DThe lightning stops suddenly, but the sky remains mostly dark.@n\r\n"))
				DRAGONC -= 1
			case 280:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@WThe glow around Shenron becomes subdued as the Eternal Dragon coils so that his head is looking down on the dragon balls!@n\r\n"))
				DRAGONC -= 1
			case 275:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CWho summoned me? I will grant you any two wishes that are within my power.@w'@n\r\n"))
				DRAGONC -= 1
			case 180:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CMake your wish already, you only have 3 minutes remaining.@w'@n\r\n"))
				DRAGONC -= 1
			case 120:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CMake your wish. I am losing patience, you only have 2 minutes left.@w'@n\r\n"))
				DRAGONC -= 1
			case 60:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@wShenron says, '@CMake your wish now! You only have 1 minute left.@w'@n\r\n"))
				DRAGONC -= 1
			case 0:
				send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("Shenron growls and disappears with a blinding flash that is absorbed into the dragon balls. The glowing dragon balls then float high into the sky, splitting into several directions and streaking across the sky!@n\r\n"))
				send_to_planet(0, ROOM_EARTH, libc.CString("@DThe sky grows brighter again as the clouds disappear magicly.@n\r\n"))
				extract_char(EDRAGON)
				SHENRON = FALSE
				DRAGONC -= 1
				save_mud_time(&time_info)
			default:
				DRAGONC -= 1
			}
			if WISH[0] == 1 && WISH[1] == 1 {
				DRAGONC = 0
				WISH[0] = 0
				WISH[1] = 0
			}
		} else {
			send_to_room(real_room(room_vnum(DRAGONR)), libc.CString("@RThe dragon balls suddenly begin to crack and darkness begins to pour out through the cracks! Shenron begins to turn pitch black slowly as the darkness escapes. Suddenly Shenron explodes out into the distance in seven parts. Each part taking a dragon ball with it!@n\r\n"))
			var num int = rand_number(200, 20000)
			var done int = FALSE
			var place int = 1
			DRAGONC = 0
			WISH[0] = 0
			WISH[1] = 0
			for done == FALSE {
				switch place {
				case 1:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON1 = num
							place = 2
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 2:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON2 = num
							place = 3
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 3:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON3 = num
							place = 4
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 4:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON4 = num
							place = 5
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 5:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON5 = num
							place = 6
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 6:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON6 = num
							place = 7
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				case 7:
					if real_room(room_vnum(num)) != room_rnum(-1) {
						if ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_EARTH) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_VEGETA) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_FRIGID) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_AETHER) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_NAMEK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_KONACK) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) || ROOM_FLAGGED(real_room(room_vnum(num)), ROOM_YARDRAT) {
							SHADOW_DRAGON7 = num
							done = TRUE
							num = rand_number(200, 20000)
						} else {
							num = rand_number(200, 20000)
						}
					} else {
						num = rand_number(200, 20000)
					}
				}
				save_mud_time(&time_info)
			}
			var mob *char_data = nil
			var r_num mob_rnum
			r_num = real_mobile(SHADOW_DRAGON1_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON1)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON2_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON2)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON3_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON3)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON4_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON4)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON5_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON5)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON6_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON6)))
			mob = nil
			r_num = real_mobile(SHADOW_DRAGON7_VNUM)
			mob = read_mobile(mob_vnum(r_num), REAL)
			char_to_room(mob, real_room(room_vnum(SHADOW_DRAGON7)))
			mob = nil
			extract_char(EDRAGON)
			SHENRON = FALSE
			DRAGONC = 0
		}
	}
}
func do_summon(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		summoned int    = FALSE
		count    int    = 0
		dball    [7]int = [7]int{20, 21, 22, 23, 24, 25, 26}
		dball2   [7]int = [7]int{20, 21, 22, 23, 24, 25, 26}
		obj      *obj_data
		next_obj *obj_data
		mob      *char_data = nil
		r_num    mob_rnum
	)
	if !ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
		send_to_char(ch, libc.CString("@wYou can not summon Shenron when you are not on earth.@n\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_NOINSTANT) || ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
		send_to_char(ch, libc.CString("You can not summon shenron in this protected area!\r\n"))
		return
	}
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_INSIDE {
		send_to_char(ch, libc.CString("Go outside to summon Shenron!\r\n"))
		return
	}
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if OBJ_FLAGGED(obj, ITEM_FORGED) {
			continue
		}
		if GET_OBJ_VNUM(obj) == obj_vnum(dball[0]) {
			dball[0] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[1]) {
			dball[1] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[2]) {
			dball[2] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[3]) {
			dball[3] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[4]) {
			dball[4] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[5]) {
			dball[5] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[6]) {
			dball[6] = -1
			count++
			continue
		} else {
			continue
		}
	}
	if count == 7 {
		summoned = TRUE
	}
	if summoned == TRUE {
		reveal_hiding(ch, 0)
		act(libc.CString("@WYou place the dragon balls on the ground and with both hands outstretched towards them you say '@CArise Eternal Dragon Shenron!@W'@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n places the dragon balls on the ground and with both hands outstretched towards them $e says '@CArise Eternal Dragon Shenron!@W'@n"), TRUE, ch, nil, nil, TO_ROOM)
		SHENRON = TRUE
		DRAGONC = 300
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			DRAGONR = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number)
		} else {
			DRAGONR = -1
		}
		if real_room(room_vnum(DRAGONR)) == room_rnum(-1) {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				DRAGONR = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number)
			} else {
				DRAGONR = -1
			}
		}
		send_to_imm(libc.CString("Shenron summoned to room: %d\r\n"), DRAGONR)
		if (func() int {
			DRAGONZ = int(real_zone_by_thing(room_vnum(DRAGONR)))
			return DRAGONZ
		}()) != int(-1) {
			DRAGONZ = int(real_zone_by_thing(room_vnum(DRAGONR)))
		}
		for obj = ch.Carrying; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if GET_OBJ_VNUM(obj) == obj_vnum(dball2[0]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[0] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[1]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[1] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[2]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[2] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[3]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[3] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[4]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[4] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[5]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[5] = -1
				continue
			} else if GET_OBJ_VNUM(obj) == obj_vnum(dball2[6]) {
				obj_from_char(obj)
				extract_obj(obj)
				dball2[6] = -1
				continue
			} else {
				continue
			}
		}
		if (func() mob_rnum {
			r_num = real_mobile(21)
			return r_num
		}()) == mob_rnum(-1) {
			send_to_imm(libc.CString("Shenron doesn't exist!"))
			return
		}
		mob = read_mobile(mob_vnum(r_num), REAL)
		char_to_room(mob, 0)
		EDRAGON = mob
		return
	} else {
		send_to_char(ch, libc.CString("@wYou do not have all the dragon balls and can not summon the dragon!@n\r\n"))
		return
	}
}
func handle_transform(ch *char_data, add int64, mult float64, drain float64) {
	if ch == nil {
		return
	}
	if IS_NPC(ch) {
		return
	} else {
		var (
			cur       float64
			max       float64
			dapercent float64 = float64(ch.Lifeperc)
		)
		cur = float64(ch.Hit)
		max = float64(ch.Max_hit)
		ch.Max_hit = int64(float64(ch.Basepl+add) * mult)
		if android_can(ch) == 1 {
			ch.Max_hit += 50000000
		} else if android_can(ch) == 2 {
			ch.Max_hit += 20000000
		}
		if (float64(ch.Hit)+float64(add)*(cur/max))*mult <= float64(gear_pl(ch)) {
			ch.Hit = int64((float64(ch.Hit) + float64(add)*(cur/max)) * mult)
		} else if (float64(ch.Hit)+float64(add)*(cur/max))*mult > float64(gear_pl(ch)) {
			ch.Hit = gear_pl(ch)
		}
		cur = float64(ch.Mana)
		max = float64(ch.Max_mana)
		ch.Max_mana = int64(float64(ch.Baseki+add) * mult)
		if android_can(ch) == 1 {
			ch.Max_mana += 50000000
		} else if android_can(ch) == 2 {
			ch.Max_mana += 20000000
		}
		if (float64(ch.Mana)+float64(add)*(cur/max))*mult <= float64(ch.Max_mana) {
			ch.Mana = int64((float64(ch.Mana) + float64(add)*(cur/max)) * mult)
		} else if (float64(ch.Mana)+float64(add)*(cur/max))*mult > float64(ch.Max_mana) {
			ch.Mana = ch.Max_mana
		}
		ch.Move -= int64(float64(ch.Move) * drain)
		cur = float64(ch.Move)
		max = float64(ch.Max_move)
		ch.Max_move = int64(float64(ch.Basest+add) * mult)
		if android_can(ch) == 1 {
			ch.Max_move += 50000000
		} else if android_can(ch) == 2 {
			ch.Max_move += 20000000
		}
		if (float64(ch.Move)+float64(add)*(cur/max))*mult <= float64(ch.Max_move) {
			ch.Move = int64((float64(ch.Move) + float64(add)*(cur/max)) * mult)
		} else if (float64(ch.Move)+float64(add)*(cur/max))*mult > float64(ch.Max_move) {
			ch.Move = ch.Max_move
		}
		if ch.Race != RACE_ANDROID {
			ch.Lifeforce = int64(float64(GET_LIFEMAX(ch)) * dapercent)
			if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
				ch.Lifeforce = int64(GET_LIFEMAX(ch))
			}
		}
	}
}
func handle_revert(ch *char_data, add uint64, mult float64) {
	if ch == nil {
		return
	}
	if IS_NPC(ch) {
		return
	} else {
		var (
			convert   float64
			dapercent float64 = float64(ch.Lifeperc)
		)
		convert = float64(ch.Hit) / float64(ch.Max_hit)
		ch.Hit = int64((float64(ch.Hit) - (float64(add)*mult)*convert) / mult)
		convert = float64(ch.Mana) / float64(ch.Max_mana)
		ch.Mana = int64((float64(ch.Mana) - (float64(add)*mult)*convert) / mult)
		convert = float64(ch.Move) / float64(ch.Max_move)
		ch.Move = int64((float64(ch.Move) - (float64(add)*mult)*convert) / mult)
		ch.Lifeforce = int64(float64(GET_LIFEMAX(ch)) * dapercent)
		if ch.Lifeforce > int64(GET_LIFEMAX(ch)) {
			ch.Lifeforce = int64(GET_LIFEMAX(ch))
		}
		if ch.Move < 1 {
			ch.Move = 1
		}
		if ch.Mana < 1 {
			ch.Mana = 1
		}
		if ch.Hit < 1 {
			ch.Hit = 1
		}
		if ch.Lifeforce < 1 {
			ch.Lifeforce = 1
		}
		if ch.Hit > gear_pl(ch) {
			ch.Hit = gear_pl(ch)
		}
		ch.Max_hit = ch.Basepl
		ch.Max_mana = ch.Baseki
		ch.Max_move = ch.Basest
	}
}
func do_transform(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		arg2 [2048]byte
		buf3 [2048]byte
	)
	if IS_NPC(ch) {
		return
	}
	two_arguments(argument, &arg[0], &arg2[0])
	if PLR_FLAGGED(ch, PLR_OOZARU) {
		send_to_char(ch, libc.CString("You are the great Oozaru right now and can't transform!\r\n"))
		return
	}
	if ch.Kaioken > 0 {
		send_to_char(ch, libc.CString("You are in kaioken right now and can't transform!\r\n"))
		return
	}
	if ch.Suppression > 0 {
		send_to_char(ch, libc.CString("You are suppressing right now and can't transform!\r\n"))
		return
	}
	if int(ch.Clones) > 0 {
		send_to_char(ch, libc.CString("You can't concentrate on transforming while your body is split into multiple forms!\r\n"))
		return
	}
	reveal_hiding(ch, 0)
	if ch.Basepl < 50000 {
		send_to_char(ch, libc.CString("@RYou are too weak to comprehend transforming!@n\r\n"))
		return
	}
	if arg[0] == 0 {
		if ch.Race == RACE_HUMAN {
			send_to_char(ch, libc.CString("              @YSuper @CHuman@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YSuper @CHuman @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CHuman @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CHuman @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CHuman @WFourth @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_SAIYAN && !PLR_FLAGGED(ch, PLR_LSSJ) {
			send_to_char(ch, libc.CString("              @YSuper @CSaiyan@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WFourth @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_SAIYAN && PLR_FLAGGED(ch, PLR_LSSJ) {
			send_to_char(ch, libc.CString("                @YSuper @CSaiyan@n\r\n"))
			send_to_char(ch, libc.CString("@b-------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WFirst   @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YLegendary @CSuper Saiyan @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b-------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_HALFBREED {
			send_to_char(ch, libc.CString("              @YSuper @CSaiyan@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CSaiyan @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_NAMEK {
			send_to_char(ch, libc.CString("              @YSuper @CNamek@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YSuper @CNamek @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CNamek @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CNamek @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper @CNamek @WFourth @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_ICER {
			send_to_char(ch, libc.CString("              @YTransform@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YTransform @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YTransform @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YTransform @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YTransform @WFourth @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_MUTANT {
			send_to_char(ch, libc.CString("              @YMutate@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YMutate @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMutate @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMutate @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_KONATSU {
			send_to_char(ch, libc.CString("              @YShadow@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YShadow @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YShadow @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YShadow @WThird  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("              @YUpgrade@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@Y1.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@Y2.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@Y3.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@Y4.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@Y5.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 5)) * 0.75) {
					return add_commas(int64(trans_req(ch, 5)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@Y6.0 @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 6)) * 0.75) {
					return add_commas(int64(trans_req(ch, 6)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_BIO {
			send_to_char(ch, libc.CString("              @YPerfection@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YMature        @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSemi-Perfect  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YPerfect       @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YSuper Perfect @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 4)) * 0.75) {
					return add_commas(int64(trans_req(ch, 4)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_TRUFFLE {
			send_to_char(ch, libc.CString("              @YAscend@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YAscend @WFirst  @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YAscend @WSecond @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YAscend @WThird @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_MAJIN {
			send_to_char(ch, libc.CString("              @YMorph@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YMorph @WAffinity @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMorph @WSuper    @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMorph @WTrue     @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else if ch.Race == RACE_KAI {
			send_to_char(ch, libc.CString("              @YMystic@n\r\n"))
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
			send_to_char(ch, libc.CString("@YMystic @WFirst     @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 1)) * 0.75) {
					return add_commas(int64(trans_req(ch, 1)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMystic @WSecond    @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 2)) * 0.75) {
					return add_commas(int64(trans_req(ch, 2)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@YMystic @WThird     @R-@G %s BPL Req\r\n"), func() *byte {
				if float64(ch.Basepl) >= (float64(trans_req(ch, 3)) * 0.75) {
					return add_commas(int64(trans_req(ch, 3)))
				}
				return libc.CString("??????????")
			}())
			send_to_char(ch, libc.CString("@b------------------------------------------------@n\r\n"))
		} else {
			send_to_char(ch, libc.CString("You do not have a transformation.\r\n"))
			return
		}
		if trans_req(ch, 1) > 0 {
			if ch.Transclass == 1 {
				send_to_char(ch, libc.CString("\r\n@RYou have @rterrible@R transformation BPL Requirements.@n\r\n"))
			} else if ch.Transclass == 2 {
				send_to_char(ch, libc.CString("\r\n@CYou have @caverage@C transformation BPL Requirements.@n\r\n"))
			} else if ch.Transclass == 3 {
				send_to_char(ch, libc.CString("\r\n@GYou have @gGREAT@G transformation BPL Requirements.@n\r\n"))
			}
		}
		return
	} else if C.strcasecmp(libc.CString("first"), &arg[0]) == 0 {
		if ch.Race == RACE_HUMAN {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou spread your feet out and crouch slightly as a bright white aura bursts around your body. Torrents of white and blue energy burn upwards around your body while your muscles grow and become more defined at the same time. In a sudden rush of power you achieve @CSuper @cHuman @GFirst@W sending surrounding debris high into the sky!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W crouches slightly while spreading $s feet as a bright white aura bursts up around $s body. Torrents of white and blue energy burn upwards around $s body while $s muscles grow and become more defined at the same time. In a sudden rush of power debris is sent flying high into the air with $m achieving @CSuper @cHuman @GFirst@W!"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 1000000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_TRUFFLE {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYour mind accelerates working through the mysteries of the universe while at the same time your body begins to change! Innate nano-technology within your body begins to activate, forming flexible metal plating across parts of your skin!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W begins to write complicated calculations in the air as though $e were possessed while at the same time $s body begins to change! Innate nano-technology within $s body begins to activate, forming flexible metal plating across parts of $s skin!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 1300000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KAI {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WThoughts begin to flow through your mind of events throughout your life. The progression leads up to more recent events and finally to this very moment. All of it's significance overwhelms you momentarily and your motivation and drive increase. As your attention is drawn back to your surroundings, you feel as though your thinking, senses, and reflexes have sharpened dramatically.  At the core of your being, a greater depth of power can be felt.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@W$n@W's face tenses, it becoming clear momentarily that they are deep in thought. After a brief lapse in focus, their attention seems to return to their surroundings. Though it's not apparent why they were so distracted, something definitely seems different about $m.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 1100000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KONATSU {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA dark shadowy aura with flecks of white energy begins to burn around your body! Strength and agility can be felt rising up within as your form becomes blurred and ethereal looking. You smile as you realize your @GFirst @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA dark shadowy aura with flecks of white energy begins to burn around @C$n@W's body! $s form becomes blurred and ethereal-looking as $s muscles become strong and lithe. $e smiles as $e achieves $s @GFirst @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 1000000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou yell with pain as your body begins to grow and power surges within! Your legs expand outward to triple their previous length. Soon after your arms, chest, and head follow. Your horns grow longer and curve upwards while lastly your tail expands. You are left confidently standing, having completed your @GFirst @cTransformation@W.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W yells with pain as $s body begins to grow and power surges outward! $s legs expand outward to triple their previous length. Soon after $s arms, chest, and head follow. $s horns grow longer and curve upwards while lastly $s tail expands. $e is left confidently standing, having completed $s @GFirst @cTransformation@W.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 400000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_MUTANT {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYour flesh grows tougher as power surges up from within. Your fingernails grow longer, sharper, and more claw-like. Lastly your muscles double in size as you achieve your @GFirst @mMutation@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W flesh grows tougher as power surges up around $m. $s fingernails grow longer, sharper, and more claw-like. Lastly $s muscles double in size as $e achieves $s @GFirst @mMutation@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 100000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_HALFBREED {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				if !PLR_FLAGGED(ch, PLR_LSSJ) && !PLR_FLAGGED(ch, PLR_FPSSJ) && rand_number(1, 500) >= 491 && ch.Max_hit > 19200000 {
					send_to_char(ch, libc.CString("You have mastered the super saiyan first transformation and have achieved Full Power Super Saiyan! You will now no longer use stamina while in this form.\r\n"))
					ch.Act[int(PLR_FPSSJ/32)] |= bitvector_t(1 << (int(PLR_FPSSJ % 32)))
				}
				var zone int = 0
				act(libc.CString("@WSomething inside your mind snaps as your rage spills over! Lightning begins to strike the ground all around you as you feel torrents of power rushing through every fiber of your being. Your hair suddenly turns golden as your eyes change to the color of emeralds. In a final rush of power a golden aura rushes up around your body! You have become a @CSuper @YSaiyan@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W screams in rage as lightning begins to crash all around! $s hair turns golden and $s eyes change to an emerald color as a bright golden aura bursts up around $s body! As $s energy stabilizes $e wears a fierce look upon $s face, having transformed into a @CSuper @YSaiyan@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 900000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_NAMEK {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou crouch down and clench your fists as your muscles begin to bulge! Sweat pours down your body as the ground beneath your feet cracks and warps under the pressure of your rising ki! With a sudden burst that sends debris flying you realize a new plateau in your power, having achieved @CSuper @gNamek @GFirst@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wcrouches down and clenches $s fists as $s muscles begin to bulge! Sweat pours down $s body as the ground beneath $s feet cracks and warps under the pressure of  $s rising ki! With a sudden burst that sends debris flying $e seems to realize a new plateau in $s power, having achieved @CSuper @gNamek @GFirst@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 200000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_SAIYAN {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				if !PLR_FLAGGED(ch, PLR_LSSJ) && !PLR_FLAGGED(ch, PLR_FPSSJ) && rand_number(1, 500) >= 491 && ch.Max_hit > 19200000 {
					send_to_char(ch, libc.CString("You have mastered the super saiyan first transformation and have achieved Full Power Super Saiyan! You will now no longer use stamina while in this form.\r\n"))
					ch.Act[int(PLR_FPSSJ/32)] |= bitvector_t(1 << (int(PLR_FPSSJ % 32)))
				}
				var zone int = 0
				act(libc.CString("@WSomething inside your mind snaps as your rage spills over! Lightning begins to strike the ground all around you as you feel torrents of power rushing through every fiber of your being. Your hair suddenly turns golden as your eyes change to the color of emeralds. In a final rush of power a golden aura rushes up around your body! You have become a @CSuper @YSaiyan@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W screams in rage as lightning begins to crash all around! $s hair turns golden and $s eyes change to an emerald color as a bright golden aura bursts up around $s body! As $s energy stabilizes $e wears a fierce look upon $s face, having transformed into a @CSuper @YSaiyan@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 800000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.1)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		return
	} else if C.strcasecmp(libc.CString("second"), &arg[0]) == 0 {
		if ch.Race == RACE_HUMAN {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WSuddenly a bright white aura bursts into existance around your body, you feel the intensity of your hidden potential boil until it can't be contained any longer! Waves of ki shoot out from your aura streaking outwards in many directions. A roar that shakes everything in the surrounding area sounds right as your energy reaches its potential and you achieve @CSuper @cHuman @GSecond@W!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W is suddenly covered with a bright white aura as $e grits $s teeth, apparently struggling with the power boiling to the surface! Waves of ki shoot out from $s aura, streaking in several directions as a mighty roar shakes everything in the surrounding area. As $s aura calms $e smiles, having achieved @CSuper @cHuman @GSecond@W!"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				var add int = 12000000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KAI {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou feel a sudden rush of emotion, escalating almost to a loss of control as your thoughts race. Your heart begins to beat fast as memories mix with the raw emotion. A faint blue glow begins to surround you. As your emotions level off, you feel a deeper understanding of the universe as you know it. You visibly calm back down to an almost steely eyed resolve as you assess your surroundings. The blue aura wicks around you for a few moments and then dissipates. Thought it's full impact is not yet clear to you, you are left feeling as though both your power and inner strength have turned into nearly bottomless wells.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@W$n@W's appears to be hit by some sudden pangs of agony, their face contorted in pain.  After a moment a faint blue aura appears around them, glowing brighter as time passes. You can feel something in the pit of your stomach, letting you know that something very significant is changing around you. Before long $n@W's aura fades, leaving a very determined looking person in your presence.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1100000, 3)
				}
				var add int = 115000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_TRUFFLE {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WComplete understanding of every physical thing floods your mind as the nano-technology within you continues to change your body! Your eyes change; becoming glassy, hard, and glowing. Your muscles merge with a nano-fiber strengthening them at the molecular level! Finally your very bones become plated in nano-metals that have yet to be invented naturally!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@.s nano-technology continues to change $s body! $s eyes change; becoming glassy, hard, and glowing. $s muscles merge with a nano-fiber strengthening them at the molecular level! Finally $s very bones become plated in nano-metals that have yet to be invented naturally!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1300000, 3)
				}
				var add int = 80000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KONATSU {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WThe shadowy aura surrounding your body burns larger than ever as dark bolts of purple electricity crackles across your skin. Your eyes begin to glow white as shockwaves of power explode outward! All the shadows in the immediate area are absorbed into your aura in an instant as you achieve your @GSecond @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe shadowy aura surrounding @C$n@W's body burns larger than ever as dark bolts of purple electricity crackles across $s skin. $s eyes begin to glow white as shockwaves of power explode outward! All the shadows in the immediate area are absorbed into $s aura in an instant as $e achieves $s @GSecond @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				var add int = 56000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_MUTANT {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WSpikes grow out from your elbows as your power begins to climb to new heights. The muscles along your forearms grow to double their former size as the spikes growing from your elbows flatten and sharpen into blades. You have achieved your @GSecond @mMutation@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WSpikes grow out from @C$n@W's elbows as $s power begins to climb to new heights. The muscles along $s forearms grow to double their former size as the spikes growing from $s elbows flatten and sharpen into blades. $e has achieved your @GSecond @mMutation@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 100000, 2)
				}
				var add int = 8500000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou yell as pain envelopes your body and a dark aura bursts up! The back of your head stretches out growing multiple spikes along its edges as it grows. In addition to this your shoulders stretch out forming pointed ridges. You cackle as lastly your nose disappears and your face becomes more lizard like. Energy swirls around your body as you realize your @GSecond @cTransformation@W.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W yells as pain envelopes $s body and a dark aura bursts up around $m! The back of $s head stretches out growing multiple spikes along its edges as it grows. In addition to this $s shoulders stretch out forming pointed ridges. $e cackles as lastly $s nose disappears and $s face becomes more lizard like. Energy swirls around $s body as $e realizes $s @GSecond @cTransformation@W.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 400000, 2)
				}
				var add int = 7000000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_HALFBREED {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WBlinding rage burns through your mind as a sudden eruption of energy surges forth! A golden aura bursts up around your body, glowing as bright as the sun. Rushing winds rocket out from your body in every direction as bolts of electricity begin to crackle in your aura. As your aura dims you are left standing confidently, having achieved @CSuper @YSaiyan @GSecond@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W stands up straight with $s head back as $e releases an ear piercing scream! A blindingly bright golden aura bursts up around $s body, glowing as bright as the sun. As rushing winds begin to rocket out from $m in every direction, bolts of electricity flash and crackle in $s aura. As $s aura begins to dim $e is left standing confidently, having achieved @CSuper @YSaiyan @GSecond@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 900000, 2)
				}
				var add int = 16500000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_NAMEK {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou gasp in shock as a power within your body that you had not been aware of begins to surge to the surface! Your muscles grow larger as energy crackles between your antennae intensely! A shockwave of energy explodes outward as you achieve a new plateau in power, @CSuper @gNamek @GSecond@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wgasps in shock as a power within $s body begins to surge out! $s muscles grow larger as energy crackles between $s antennae intensely! A shockwave of energy explodes outward as $e achieves a new plateau in power, @CSuper @gNamek @GSecond@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 200000, 2)
				}
				var add int = 4000000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_SAIYAN {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.2 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				if !PLR_FLAGGED(ch, PLR_LSSJ) {
					act(libc.CString("@WBlinding rage burns through your mind as a sudden eruption of energy surges forth! A golden aura bursts up around your body, glowing as bright as the sun. Rushing winds rocket out from your body in every direction as bolts of electricity begin to crackle in your aura. As your aura dims you are left standing confidently, having achieved @CSuper @YSaiyan @GSecond@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@W stands up straight with $s head back as $e releases an ear piercing scream! A blindingly bright golden aura bursts up around $s body, glowing as bright as the sun. As rushing winds begin to rocket out from $m in every direction, bolts of electricity flash and crackle in $s aura. As $s aura begins to dim $e is left standing confidently, having achieved @CSuper @YSaiyan @GSecond@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				} else {
					act(libc.CString("@WYou roar and then stand at your full height. You flex every muscle in your body as you feel your strength grow! Your eyes begin to glow @wwhite@W with energy, your hair turns @Ygold@W, and at the same time a @wbright @Yg@yo@Yl@yd@Ye@yn@W aura flashes up around your body! You release your @YL@ye@Dg@We@wn@Yd@ya@Dr@Yy@W power upon the universe!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n @Wroars and then stands at $s full height. Then $s muscles start to buldge and grow as $e flexes them! Suddenly $s eyes begin to glow @wwhite@W with energy, $s hair turns @Ygold@W, and at the same time a @wbright @Yg@yo@Yl@yd@Ye@yn@W aura flashes up around $s body! @C$n@W releases $s @YL@ye@Dg@We@wn@Yd@ya@Dr@Yy@W power upon the universe!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 0
				var mult float64 = 0
				if !PLR_FLAGGED(ch, PLR_LSSJ) {
					add = 20000000
					mult = 3
				} else {
					add = 185000000
					mult = 6
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 800000, 2)
				}
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		return
	} else if C.strcasecmp(libc.CString("third"), &arg[0]) == 0 {
		if ch.Race == RACE_HUMAN {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou clench both of your fists as the bright white aura around your body is absorbed back into your flesh. As it is absorbed, your muscles triple in size and electricity crackles across your flesh. You grin as you feel the power of @CSuper @cHuman @GThird@W!"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W clenches both of $s fists as the bright white aura around $s body is absorbed back into $s flesh. As it is absorbed, $s muscles triple in size and bright electricity crackles across $s flesh. $e smiles as $e achieves the power of @CSuper @cHuman @GThird@W!"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 12000000, 3)
				}
				var add int = 50000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_TRUFFLE {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou have reached the final stage of enlightenment and the nano-technology thriving inside you begin to initiate the changes! Your neural pathways become refined, your reflexes honed, your auditory and ocular senses sharpening far beyond normal levels! Your gravitational awareness improves, increasing sensitivity and accuracy in your equilibrum!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n begins to mumble quietly, slowly at first and gradually picking up speed. A glint is seen from $s eyes and $s arms reach outwards briefly as $e appears to catch his balance. $s arms drop back to $s sides as balance is regained, a vicious smile on $s face.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 80000000, 4)
				}
				var add int = 300000000
				var mult float64 = 5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KAI {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYour minds' eye becomes overwhelmed by secrets unimaginable. The threads of the very universe become visible in your heightened state of awareness. Reaching out, a single thread vibrates, producing a @Rred @Wcolor -- yours. Your fingertips brush against it and your senses become clouded by a vast expanse of white color and noise. As your vision and hearing return, you understand the threads tying every living being together. Your awareness has expanded beyond comprehension!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W's eyes grow wide, mouth agape. $s body begins to shiver uncontrollably! $s arms reaches out cautiously before falling back down to $s side. $s face relaxes visibly, features returning to a normal state. $s irises remain larger than before, a slight smile softening $s gaze.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1100000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 115000000, 4)
				}
				var add int = 270000000
				var mult float64 = 5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_MUTANT {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA dark cyan aura bursts up around your body as the ground begins to crack beneath you! You scream out in pain as your power begins to explode! Two large spikes grow out from your shoulder blades as you reach your @GThird @mMutation!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA dark cyan aura bursts up around @C$n@W's body as the ground begins to crack beneath $m and $e screams out in pain as $s power begins to explode! Two large spikes grow out from $s shoulder blades as $e reaches $s @GThird @mMutation!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 100000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 8500000, 3)
				}
				var add int = 80000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_KONATSU {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WThe shadowy aura around you explodes outward as your power begins to rise!  You're overcome with a sudden realization, that the shadows are an extension of yourself, that light isn't needed for your shadows to bloom.  With this newfound wisdom comes ability and power!  The color in your aura drains as the shadows slide inward and cling to your body like a second, solid black skin!  Shockwaves roll off of you in quick succession, pelting the surrounding area harshly!  Accompanying the waves, a pool of darkness blossoms underneath you, slowly spreading the shadows to the whole area, projecting onto any surface nearby!  Purple and black electricity crackle in your solid white aura, and you grin as you realize your @GThird @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe shadowy aura around $n explodes outward as $s power begins to rise!  Realization dawns on $s face, followed shortly by confidence! The color in $s aura drains as the shadows slide inward to cling to $s body like a second, solid black skin! Shockwaves roll off of $n in quick succession, pelting the surrounding area harshly!  Accompanying the waves, a pool of darkness blossoms underneath them, slowly spreading the shadows to the whole area, projecting onto any surface nearby! Purple and black electricity crackle in $s solid white aura, and he grins as $e realizes $s @GThird @DShadow @BForm@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 56000000, 4)
				}
				var add int = 290000000
				var mult float64 = 5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA blinding light surrounds your body while your rising power begins to rip up the ground beneath you! Your skin and torso shell begin to crack as your new body struggles to free its self. Huge chunks of debris lift free of the ground as your power begins to rise to unbelievable heights. Suddenly your old skin and torso shell burst off from your body, leaving a sleek form glowing where they had been. Everything comes crashing down as your power evens out, leaving you with your @GThird @cTransformation @Wcompleted!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA blinding light surrounds @C$n@W's body while $s rising power begins to rip up the ground beneath $m! $s skin and torso shell begin to crack as $s new body struggles to free its self. Huge chunks of debris lift free of the ground as $s power begins to rise to unbelievable heights. Suddenly $s old skin and torso shell burst off from $s body, leaving a sleek form glowing where they had been. Everything comes crashing down as @C$n@W's power evens out, leaving $m with $s @GThird @cTransformation @Wcompleted!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 400000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 7000000, 3)
				}
				var add int = 45000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_HALFBREED {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WElectricity begins to crackle around your body as your aura grows explosively! You yell as your powerlevel begins to skyrocket while your hair grows to multiple times the length it was previously. Your muscles become incredibly dense instead of growing in size, preserving your speed. Finally your irises appear just as your transformation becomes complete, having achieved @CSuper @YSaiyan @GThird@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WElectricity begins to crackle around @C$n@W, as $s aura grows explosively! $e yells as the energy around $m skyrockets and $s hair grows to multiple times its previous length. $e smiles as $s irises appear and $s muscles tighten up. $s transformation complete, $e now stands confidently, having achieved @CSuper @YSaiyan @GThird@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 900000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 16500000, 4)
				}
				var add int = 240000000
				var mult float64 = 5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_NAMEK {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA fierce clear aura bursts up around your body as you struggle to control a growing power within! Energy leaks off of your aura at an astounding rate filling the air around you with small orbs of ki. As your power begins to level off the ambient ki hovering around you is absorbed inward in a sudden shock that leaves your skin glowing! You have achieved a rare power, @CSuper @gNamek @GThird@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA fierce clear aura bursts up around @C$n@W's body as $e struggles to control $s own power! Energy leaks off of $s aura at an astounding rate filling the air around $m with small orbs of ki. As $s power begins to level off the ambient ki hovering around $m is absorbed inward in a sudden shock that leaves $s skin glowing! $e has achieved a rare power, @CSuper @gNamek @GThird@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 200000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 4000000, 3)
				}
				var add int = 65000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_SAIYAN && !PLR_FLAGGED(ch, PLR_LSSJ) {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= float64(ch.Max_move)*0.1 {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WElectricity begins to crackle around your body as your aura grows explosively! You yell as your powerlevel begins to skyrocket while your hair grows to multiple times the length it was previously. Your muscles become incredibly dense instead of growing in size, preserving your speed. Finally your irises appear just as your transformation becomes complete, having achieved @CSuper @YSaiyan @GThird@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WElectricity begins to crackle around @C$n@W, as $s aura grows explosively! $e yells as the energy around $m skyrockets and $s hair grows to multiple times its previous length. $e smiles as $s irises appear and $s muscles tighten up. $s transformation complete, $e now stands confidently, having achieved @CSuper @YSaiyan @GThird@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 80000000
				var mult float64 = 4
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 800000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 20000000, 3)
				}
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		return
	} else if C.strcasecmp(libc.CString("fourth"), &arg[0]) == 0 {
		if ch.Race == RACE_HUMAN {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[4] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYou grit your teeth and clench your fists as a sudden surge of power begins to tear through your body! Your muscles lose volume and gain mass, condensing into sleek hyper efficiency as a spectacular shimmering white aura flows over you, flashes of multicolored light flaring up in rising stars around your new form, a corona of glory! You feel your ultimate potential realized as you ascend to @CSuper @cHuman @GFourth@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W grits $s teeth and clenches $s fists as a sudden surge of power begins to tear through $s body! $n@W's muscles lose volume and gain mass, condensing into sleek hyper efficiency as a spectacular shimmering white aura flows over $m, flashes of multicolored light flare up in rising stars around $s new form, a corona of glory! $n@W smiles as his ultimate potential is realized as $e ascends to @CSuper @cHuman @GFourth@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 12000000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
					handle_revert(ch, 50000000, 4)
				}
				var add int = 270000000
				var mult float64 = 4.5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[4] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA feeling of complete power courses through your viens as your body begins to change radically! You triple in height while a hard shell forms over your entire torso. Hard bones grow out from your head forming four ridges that jut outward. A hard covering grows up over your mouth and nose completing the transformation! A dark crimson aura flames around your body as you realize your @GFourth @cTransformation@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W's body begins to change radically! $e triples in height while a hard shell forms over $s entire torso. Hard bones grow out from $s head forming four ridges that jut outward. A hard covering grows up over $s mouth and nose completing the transformation! A dark crimson aura flames around @C$n@W's body as $e realizes $s @GFourth @cTransformation@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 400000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 7000000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
					handle_revert(ch, 45000000, 4)
				}
				var add int = 200000000
				var mult float64 = 5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_NAMEK {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[4] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WAn inner calm fills your mind as your power surges higher than ever before. Complete clarity puts everything once questioned into perspective. While this inner calm is filling your mind, an outer storm of energy erupts around your body! The storm of energy boils and crackles while growing larger. You have achieved @CSuper @gNamek @GFourth@W, a mystery of the ages.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W smiles calmly as a look of complete understand fills $s eyes. While $e remains perfectly calm and detached a massivly powerful storm of energy erupts from his body. This storm of energy shimmers with the colors of the rainbow and boils and crackles with awesome power! $s smile disappears as he realizes a mysterious power of the ages, @CSuper @gNamek @GFourth@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 200000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 4000000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
					handle_revert(ch, 65000000, 4)
				}
				var add int = 230000000
				var mult float64 = 4.5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		if ch.Race == RACE_SAIYAN && !PLR_FLAGGED(ch, PLR_LSSJ) {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that form! Try 'revert'.\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[4] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WHaving absorbed enough blutz waves, your body begins to transform! Red fur grows over certain parts of your skin as your hair grows longer and unkempt. A red outline forms around your eyes while the irises of those very same eyes change to an amber color. Energy crackles about your body violently as you achieve the peak of saiyan perfection, @CSuper @YSaiyan @GFourth@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WHaving absorbed enough blutz waves, @C$n@W's body begins to transform! Red fur grows over certain parts of $s skin as $s hair grows longer and unkempt. A red outline forms around $s eyes while the irises of those very same eyes change to an amber color. Energy crackles about $s body violently as $e achieves the peak of saiyan perfection, @CSuper @YSaiyan @GFourth@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 800000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 20000000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
					handle_revert(ch, 80000000, 4)
				}
				var add int = 182000000
				var mult float64 = 5.5
				handle_transform(ch, int64(add), mult, 0.2)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
		return
	} else if C.strcasecmp(libc.CString("mature"), &arg[0]) == 0 {
		if ch.Race == RACE_BIO {
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if ch.Absorbs > 2 {
				send_to_char(ch, libc.CString("You need to absorb something to transform!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@gYou bend over as @rpain@g wracks your body! Your limbs begin to grow out, becoming more defined and muscular. As your limbs finish growing outward you feel a painful sensation coming from your back as a long tail with a spike grows out of your back! As the pain subsides you stand up straight and a current of power shatters part of the ground beneath you. You have @rmatured@g beyond your @Gl@ga@Dr@gv@Ga@ge stage!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@W$n @gbends over as a @rpainful@g look covers $s face! $s limbs begin to grow out, becoming more defined and muscular. As $s limbs finish growing outward $e screams as a long tail with a spike grows rips out of $s back! As $e calms $e stands up straight and a current of power shatters part of the ground beneath $m. $e has @rmatured@g beyond $s @Gl@ga@Dr@gv@Ga@ge stage!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				var add int = 1000000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				return
			}
		}
		return
	} else if C.strcasecmp(libc.CString("semi-perfect"), &arg[0]) == 0 || C.strcasecmp(libc.CString("Semi-Perfect"), &arg[0]) == 0 {
		if ch.Race == RACE_BIO {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Absorbs > 1 {
				send_to_char(ch, libc.CString("You need to absorb something to transform!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYour exoskeleton begins to glow spectacularly while the shape of your body begins to change. Your tail shrinks slightly. Your hands, feet, and facial features become more refined. While your body colors change slightly. The crests on your head change, standing up straighter on either side of your head as well. As you finish transforming a wave of power floods your being. You have achieved your @gSemi@D-@GPerfect @BForm@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W's exoskeleton begins to glow spectacularly while the shape of $s body begins to change. $s tail shrinks slightly. $s hands, feet, and facial features become more refined. While $s body colors change slightly. The crests on $s head change, standing up straighter on either side of $s head as well. As $e finishes transforming a wave of power rushes out from $m. $e has achieved $s @gSemi@D-@GPerfect @BForm@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				var add int = 8000000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if C.strcasecmp(libc.CString("perfect"), &arg[0]) == 0 {
		if ch.Race == RACE_BIO {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Absorbs > 0 {
				send_to_char(ch, libc.CString("You need to absorb something to transform!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if float64(ch.Move) <= (float64(ch.Max_move) * 0.2) {
				send_to_char(ch, libc.CString("You do not have enough stamina!"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WYour whole body is engulfed in blinding light as your exoskeleton begins to change shape! Your hands, feet, and facial features become more refined and humanoid. While your colors change, becoming more subdued and neutral. A bright golden aura bursts up around your body as you achieve your @GPerfect @BForm@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W whole body is engulfed in blinding light as $s exoskeleton begins to change shape! $s hands, feet, and facial features become more refined and humanoid. While $s colors change, becoming more subdued and neutral. A bright golden aura bursts up around $s body as $e achieves $s @GPerfect @BForm@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 8000000, 3)
				}
				var add int = 70000000
				var mult float64 = 3.5
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if (C.strcasecmp(libc.CString("super"), &arg[0]) == 0 || C.strcasecmp(libc.CString("super perfect"), &arg[0]) == 0) && ch.Race == RACE_BIO {
		if ch.Race == RACE_BIO {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[4] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA rush of power explodes from your perfect body, crushing nearby debris and sending dust billowing in all directions. Electricity crackles throughout your aura intensely while your muscles grow slightly larger but incredibly dense. You smile as you realize that you have taken your perfect form beyond imagination. You are now @CSuper @GPerfect@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA rush of power explodes from @C$n@W's perfect body, crushing nearby debris and sending dust billowing in all directions. Electricity crackles throughout $s aura intensely while $s muscles grow slightly larger but incredibly dense. $e smiles as $e has taken $s perfect form beyond imagination. $e is now @CSuper @GPerfect@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1000000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 8000000, 3)
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
					handle_revert(ch, 70000000, 3.5)
				}
				var add int = 400000000
				var mult float64 = 4
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if C.strcasecmp(libc.CString("affinity"), &arg[0]) == 0 {
		if ch.Race == RACE_MAJIN {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[1] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WA dark pink aura bursts up around your body as images of good and evil fill your mind! You feel the power within your body growing intensely, reflecting your personal alignment as your body changes!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WA dark pink aura bursts up around @C$n@W's body as images of good and evil fill $s mind! $e feels the power within $s body growing intensely, reflecting $s personal alignment as $s body changes!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				var add int = 1250000
				var mult float64 = 2
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if C.strcasecmp(libc.CString("super"), &arg[0]) == 0 {
		if ch.Race == RACE_MAJIN {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			}
			if ch.Absorbs > 0 && GET_LEVEL(ch) < 50 {
				send_to_char(ch, libc.CString("You need to ingest someone before you can use that form.\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[2] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WAn intense pink aura surrounds your body as it begins to change, taking on the characteristics of those you have ingested! Explosions of pink energy burst into existence all around you as your power soars to sights unseen!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WAn intense pink aura surrounds @C$n@W's body as it begins to change, taking on the characteristics of those $e has ingested! Explosions of pink energy burst into existence all around $m as $s power soars to sights unseen!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1250000, 2)
				}
				var add int = 15000000
				var mult float64 = 3
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if C.strcasecmp(libc.CString("true"), &arg[0]) == 0 {
		if ch.Race == RACE_MAJIN {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that form!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that form!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that transformation!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[3] = TRUE
					}
				}
				var zone int = 0
				act(libc.CString("@WRipples of intense pink energy rush upwards around your body as it begins to morph into its truest form! The ground beneath your feet forms into a crater from the very pressure of your rising ki! Earthquakes shudder throughout the area as your finish morphing!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WRipples of intense pink energy rush upwards around @C$n@W's body as it begins to morph into its truest form! The ground beneath $s feet forms into a crater from the very pressure of $s rising ki! Earthquakes shudder throughout the area as $e finishes morphing!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if (func() int {
					zone = int(real_zone_by_thing(room_vnum(ch.In_room)))
					return zone
				}()) != int(-1) {
					send_to_zone(libc.CString("An explosion of power ripples through the surrounding area!\r\n"), zone_rnum(zone))
				}
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 1250000, 2)
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					handle_revert(ch, 15000000, 3)
				}
				var add int = 340000000
				var mult float64 = 4.5
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				send_to_sense(0, libc.CString("You sense a nearby power grow unbelievably!"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Transformed Powerlevel@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				return
			}
		}
	} else if C.strcasecmp(libc.CString("1.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS2) || PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that upgrade!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 1)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 1)) {
				if (ch.Transcost[1]) == FALSE {
					if (ch.Player_specials.Class_skill_points[ch.Chclass]) < 50 {
						send_to_char(ch, libc.CString("You need 50 practice points in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Player_specials.Class_skill_points[ch.Chclass] -= 50
						ch.Transcost[1] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				var add int = 5000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 7500000
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS1/32)] |= bitvector_t(1 << (int(PLR_TRANS1 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("2.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS3) || PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that upgrade!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 2)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 2)) {
				if (ch.Transcost[2]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[2] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				}
				var add int = 20000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 30000000
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS2/32)] |= bitvector_t(1 << (int(PLR_TRANS2 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("3.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already beyond that upgrade!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if !PLR_FLAGGED(ch, PLR_TRANS2) {
				send_to_char(ch, libc.CString("You havn't upgraded to 2.0 yet!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 3)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 3)) {
				if (ch.Transcost[3]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[3] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				}
				var add int = 125000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 187500000
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS3/32)] |= bitvector_t(1 << (int(PLR_TRANS3 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("4.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS5) || PLR_FLAGGED(ch, PLR_TRANS6) {
				send_to_char(ch, libc.CString("You are already beyond that upgrade!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if !PLR_FLAGGED(ch, PLR_TRANS3) {
				send_to_char(ch, libc.CString("You havn't upgraded to 3.0 yet!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 4)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 4)) {
				if (ch.Transcost[4]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[4] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				}
				var add uint64 = 1000000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 1500000000
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS4/32)] |= bitvector_t(1 << (int(PLR_TRANS4 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("5.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS6) {
				send_to_char(ch, libc.CString("You are already beyond that upgrade!\r\n"))
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS5) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if !PLR_FLAGGED(ch, PLR_TRANS4) {
				send_to_char(ch, libc.CString("You havn't upgraded to 4.0 yet!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 5)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 5)) {
				if (ch.Transcost[5]) == FALSE {
					if ch.Rp < 1 {
						send_to_char(ch, libc.CString("You need 1 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 1
						ch.Transcost[5] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS4) {
					ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				}
				var add uint64 = 25000000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 3750000000
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS5/32)] |= bitvector_t(1 << (int(PLR_TRANS5 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("6.0"), &arg[0]) == 0 {
		if ch.Race == RACE_ANDROID {
			if PLR_FLAGGED(ch, PLR_TRANS6) {
				send_to_char(ch, libc.CString("You are already in that upgrade!\r\n"))
				return
			}
			if !PLR_FLAGGED(ch, PLR_TRANS5) {
				send_to_char(ch, libc.CString("You havn't upgraded to 5.0 yet!\r\n"))
				return
			}
			if ch.Basepl < int64(trans_req(ch, 6)) {
				send_to_char(ch, libc.CString("You are not strong enough to handle that upgrade!\r\n"))
				return
			} else if ch.Basepl >= int64(trans_req(ch, 6)) {
				if (ch.Transcost[6]) == FALSE {
					if ch.Rp < 2 {
						send_to_char(ch, libc.CString("You need 2 RPP in order to obtain a transformation for the first time.\r\n"))
						return
					} else {
						ch.Rp -= 2
						ch.Transcost[6] = TRUE
					}
				}
				act(libc.CString("@WYou stop for a moment as the nano-machines within your body reprogram and restructure you. You are now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @Wstops for a moment as the nano-machines within $s body reprogram and restructure $m. $e is now more powerful and efficient!@n"), TRUE, ch, nil, nil, TO_ROOM)
				if PLR_FLAGGED(ch, PLR_TRANS1) {
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS2) {
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS3) {
					ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS4) {
					ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				}
				if PLR_FLAGGED(ch, PLR_TRANS5) {
					ch.Act[int(PLR_TRANS5/32)] &= bitvector_t(^(1 << (int(PLR_TRANS5 % 32))))
				}
				var add uint64 = 1000000000
				add += 1000000000
				add += 1000000000
				add += 1000000000
				add += 1000000000
				var mult float64 = 1
				if PLR_FLAGGED(ch, PLR_SENSEM) {
					add += 750000000
					add *= 100
				}
				handle_transform(ch, int64(add), mult, 0.0)
				ch.Act[int(PLR_TRANS6/32)] |= bitvector_t(1 << (int(PLR_TRANS6 % 32)))
				return
			}
		}
	} else if C.strcasecmp(libc.CString("revert"), &arg[0]) == 0 {
		if ch.Race == RACE_HUMAN {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				act(libc.CString("@wYou revert from @CSuper @cHuman @GFourth@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cHuman @GFourth@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 270000000, 4.5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CSuper @cHuman @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cHuman @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 50000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CSuper @cHuman @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cHuman @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 12000000, 3)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CSuper @cHuman @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cHuman @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 1000000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_SAIYAN {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GFourth@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GFourth@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 182000000, 5.5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 80000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				if !PLR_FLAGGED(ch, PLR_LSSJ) {
					act(libc.CString("@wYou revert from @CSuper @cSaiyan @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
					if ch.Charge > 0 {
						do_charge(ch, libc.CString("release"), 0, 0)
					}
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 20000000, 3)
				} else {
					act(libc.CString("@wYou revert from your @YLegendary @CSuper Saiyan@w form.@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@w$n@w reverts from $s @YLegendary @CSuper Saiyan@w form@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
					if ch.Charge > 0 {
						do_charge(ch, libc.CString("release"), 0, 0)
					}
					ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
					ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
					handle_revert(ch, 185000000, 6)
				}
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 800000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_HALFBREED {
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 240000000, 5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 16500000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CSuper @cSaiyan @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cSaiyan @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 900000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_NAMEK {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				act(libc.CString("@wYou revert from @CSuper @cNamek @GFourth@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cNamek @GFourth@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 230000000, 4.5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CSuper @cNamek @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cNamek @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 65000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CSuper @cNamek @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cNamek @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 4000000, 3)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CSuper @cNamek @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CSuper @cNamek @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 200000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS4) {
				act(libc.CString("@wYou revert from @CTransform @GFourth@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CTransform @GFourth@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS4/32)] &= bitvector_t(^(1 << (int(PLR_TRANS4 % 32))))
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 200000000, 5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CTransform @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CTransform @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 45000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CTransform @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CTransform @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 7000000, 3)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CTransform @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CTransform @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 400000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_MUTANT {
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CMutate @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMutate @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 80000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CMutate @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMutate @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 8500000, 3)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CMutate @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMutate @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 100000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_KONATSU {
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CShadow @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CShadow @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 290000000, 5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CShadow @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CShadow @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 56000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CShadow @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CShadow @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 1000000, 2)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else if ch.Race == RACE_KAI {
			if PLR_FLAGGED(ch, PLR_TRANS3) {
				act(libc.CString("@wYou revert from @CMystic @GThird@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMystic @GThird@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS3/32)] &= bitvector_t(^(1 << (int(PLR_TRANS3 % 32))))
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 270000000, 5)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			}
			if PLR_FLAGGED(ch, PLR_TRANS2) {
				act(libc.CString("@wYou revert from @CMystic @GSecond@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMystic @GSecond@w.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS2/32)] &= bitvector_t(^(1 << (int(PLR_TRANS2 % 32))))
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 115000000, 4)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else if PLR_FLAGGED(ch, PLR_TRANS1) {
				act(libc.CString("@wYou revert from @CMystic @GFirst@w.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@w$n@w reverts from @CMystic @GFirst.@n"), TRUE, ch, nil, nil, TO_ROOM)
				if ch.Charge > 0 {
					do_charge(ch, libc.CString("release"), 0, 0)
				}
				ch.Act[int(PLR_TRANS1/32)] &= bitvector_t(^(1 << (int(PLR_TRANS1 % 32))))
				handle_revert(ch, 1100000, 3)
				if arg2[0] != 0 {
					do_transform(ch, &arg2[0], 0, 0)
				}
				return
			} else {
				send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
				return
			}
		} else {
			send_to_char(ch, libc.CString("You do not need to revert from any form!\r\n"))
		}
		return
	} else if ch.Race == RACE_KANASSAN || ch.Race == RACE_DEMON {
		send_to_char(ch, libc.CString("You do not have a transformation.\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("What form?\r\n"))
		return
	}
}
func do_situp(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		cost  int = 1
		bonus int = 0
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are a mob fool!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
		send_to_char(ch, libc.CString("The fire makes it too hot!\r\n"))
		return
	}
	if ch.Drag != nil {
		send_to_char(ch, libc.CString("You are dragging someone!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_FISHING) {
		send_to_char(ch, libc.CString("Stop fishing first.\r\n"))
		return
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are carrying someone!\r\n"))
		return
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity != 10 && ch.Chclass == CLASS_BARDOCK || ch.Chclass != CLASS_BARDOCK {
		cost += (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * 2) + GET_LEVEL(ch))
	} else if ch.Chclass == CLASS_BARDOCK {
		cost += (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * GET_LEVEL(ch)
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
		cost += 1000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
		cost += 2000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
		cost += 7500000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
		cost += 15000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
		cost += 25000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
		cost += 50000000
	}
	if cost == 1 || cost == 0 {
		cost = 25
	}
	if (ch.Bonuses[BONUS_HARDWORKER]) > 0 {
		cost -= int(float64(cost) * 0.25)
	} else if (ch.Bonuses[BONUS_SLACKER]) > 0 {
		cost += int(float64(cost) * 0.25)
	}
	if ch.Relax_count >= 464 {
		cost *= 50
	} else if ch.Relax_count >= 232 {
		cost *= 15
	} else if ch.Relax_count >= 116 {
		cost *= 4
	}
	if ch.Race == RACE_ANDROID || ch.Race == RACE_BIO || ch.Race == RACE_MAJIN || ch.Race == RACE_ARLIAN {
		send_to_char(ch, libc.CString("You will gain nothing from exercising!\r\n"))
		return
	}
	if limb_ok(ch, 1) == 0 {
		return
	}
	if ch.Move < int64(cost) {
		send_to_char(ch, libc.CString("You are too tired!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_SPAR) {
		send_to_char(ch, libc.CString("You shouldn't be sparring if you want to work out, it could be dangerous.\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are fighting you moron!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_FLYING) {
		send_to_char(ch, libc.CString("You can't do situps in midair!\r\n"))
		return
	} else {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 0 {
			bonus = 1
			act(libc.CString("@gYou do a situp.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 {
			bonus = rand_number(3, 7)
			act(libc.CString("@gYou do a situp, feeling the strain of gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 20 {
			bonus = rand_number(8, 14)
			act(libc.CString("@gYou do a situp, and feel gravity's pull.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 30 {
			bonus = rand_number(14, 20)
			act(libc.CString("@gYou do a situp, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 40 {
			bonus = rand_number(20, 35)
			act(libc.CString("@gYou do a situp, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 50 {
			bonus = rand_number(40, 60)
			act(libc.CString("@gYou do a situp, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 100 {
			bonus = rand_number(180, 250)
			act(libc.CString("@gYou do a situp, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 200 {
			bonus = rand_number(400, 600)
			act(libc.CString("@gYou do a situp, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
			bonus = rand_number(800, 1200)
			act(libc.CString("@gYou do a situp, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
			bonus = rand_number(2000, 3000)
			act(libc.CString("@gYou do a situp, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
			bonus = rand_number(4000, 6000)
			act(libc.CString("@gYou do a situp, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
			bonus = rand_number(9000, 10000)
			act(libc.CString("@gYou do a situp, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
			bonus = rand_number(15000, 20000)
			act(libc.CString("@gYou do a situp, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
			bonus = rand_number(25000, 30000)
			act(libc.CString("@gYou do a situp, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a situp, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		if cost >= int(ch.Max_move/2) {
			send_to_char(ch, libc.CString("This gravity is a great challenge for you!\r\n"))
			bonus *= 10
		} else if cost >= int(ch.Max_move/4) {
			send_to_char(ch, libc.CString("This gravity is an awesome challenge for you!\r\n"))
			bonus *= 8
		} else if cost >= int(ch.Max_move/10) {
			send_to_char(ch, libc.CString("This gravity is a good challenge for you!\r\n"))
			bonus *= 6
		} else if cost < int(ch.Max_move/1000) {
			send_to_char(ch, libc.CString("This gravity is so easy to you, you could do it in your sleep...\r\n"))
			bonus /= 8
		} else if cost < int(ch.Max_move/100) {
			send_to_char(ch, libc.CString("This gravity is the opposite of a challenge for you...\r\n"))
			bonus /= 5
		} else if cost < int(ch.Max_move/50) {
			send_to_char(ch, libc.CString("This gravity is definitely not a challenge for you...\r\n"))
			bonus /= 4
		} else if cost < int(ch.Max_move/30) {
			send_to_char(ch, libc.CString("This gravity is barely a challenge for you...\r\n"))
			bonus /= 3
		} else if cost < int(ch.Max_move/20) {
			send_to_char(ch, libc.CString("This gravity is hardly a challenge for you...\r\n"))
			bonus /= 2
		} else {
			send_to_char(ch, libc.CString("This gravity is just perfect for you...\r\n"))
			bonus *= 4
		}
		if !soft_cap(ch, 2) {
			bonus = 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			send_to_char(ch, libc.CString("@rThis place feels like it operates on a different time frame, it feels great...@n\r\n"))
			bonus *= 10
		} else if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus *= 10
			} else {
				bonus *= 5
			}
		} else if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			send_to_char(ch, libc.CString("@rThis place feels like... Magic.@n\r\n"))
			bonus *= 20
		}
		if bonus <= 0 && !ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 1
		}
		if bonus <= 0 && ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 15
		}
		if bonus <= 0 && ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus = 12
			} else {
				bonus = 6
			}
		}
		bonus += GET_LEVEL(ch) / 20
		if ch.Race == RACE_NAMEK {
			bonus -= bonus / 4
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			bonus += int(float64(bonus) * 0.5)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			bonus += int(float64(bonus) * 0.1)
		}
		send_to_char(ch, libc.CString("You feel slightly more vigorous @D[@G+%s@D]@n.\r\n"), add_commas(int64(bonus)))
		if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_move += int64(bonus * 3)
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_move += int64(bonus * 4)
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_move += int64(bonus * 5)
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 1 {
			ch.Max_move += int64(bonus * 2)
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 2 {
			ch.Max_move += int64(bonus * 3)
		} else {
			ch.Max_move += int64(bonus)
		}
		ch.Basest += int64(bonus)
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 50 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 200 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 500 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 1000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 5000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 10000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*7)
		}
		ch.Move -= int64(cost)
		if ch.Move < 0 {
			ch.Move = 0
		}
	}
}
func do_meditate(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		bonus  int64 = 0
		cost   int64 = 1
		weight int64 = 0
		obj    *obj_data
		arg    [2048]byte
		arg2   [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are a mob fool!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
		send_to_char(ch, libc.CString("The fire makes it too hot!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_FISHING) {
		send_to_char(ch, libc.CString("Stop fishing first.\r\n"))
		return
	}
	if ch.Race == RACE_ANDROID || ch.Race == RACE_BIO || ch.Race == RACE_MAJIN || ch.Race == RACE_ARLIAN {
		send_to_char(ch, libc.CString("You will gain nothing from exercising!\r\n"))
		return
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are carrying someone!\r\n"))
		return
	}
	if ch.Drag != nil {
		send_to_char(ch, libc.CString("You are dragging someone!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_SPAR) {
		send_to_char(ch, libc.CString("You shouldn't be sparring if you want to work out, it could be dangerous.\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are fighting you moron!\r\n"))
		return
	}
	if ch.Position != POS_SITTING {
		send_to_char(ch, libc.CString("You need to be sitting to meditate.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: meditate (object)\nSyntax: meditate expand\nSyntax: meditate break\r\n"))
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("expand")) == 0 {
		var cost int = 3500
		if ch.Race == RACE_SAIYAN {
			cost = 7000
		}
		if (ch.Player_specials.Class_skill_points[ch.Chclass]) < cost {
			send_to_char(ch, libc.CString("You do not have enough practice sessions to expand your mind and ability to remember skills.\r\n"))
			send_to_char(ch, libc.CString("%s needed.\r\n"), add_commas(int64(cost)))
		} else if ch.Skill_slots >= 60 && (ch.Bonuses[BONUS_GMEMORY]) == 0 {
			send_to_char(ch, libc.CString("You can not have any more slots through this process.\r\n"))
		} else if ch.Skill_slots >= 65 && (ch.Bonuses[BONUS_GMEMORY]) == 1 {
			send_to_char(ch, libc.CString("You can not have any more slots through this process.\r\n"))
		} else {
			send_to_char(ch, libc.CString("During your meditation you manage to expand your mind and get the feeling you could learn some new skills.\r\n"))
			ch.Skill_slots += 1
			ch.Player_specials.Class_skill_points[ch.Chclass] -= cost
			return
		}
		return
	} else if C.strcasecmp(&arg[0], libc.CString("break")) == 0 {
		if ch.Mindlink == nil {
			send_to_char(ch, libc.CString("You are not mind linked with anyone.\r\n"))
			return
		} else if ch.Linker == 1 {
			send_to_char(ch, libc.CString("This is not how you break YOUR mind link.\r\n"))
			return
		} else if float64(ch.Mana) < float64(ch.Mindlink.Max_mana)*0.05 {
			send_to_char(ch, libc.CString("You do not have enough ki to manage a break.\r\n"))
			return
		} else if float64(int(ch.Aff_abils.Intel)+rand_number(-5, 10)) >= float64(ch.Mindlink.Aff_abils.Intel)+float64(GET_SKILL(ch.Mindlink, SKILL_TELEPATHY))*0.1 {
			act(libc.CString("@rYou manage to break the mind link between you and @R$N@r!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_CHAR)
			act(libc.CString("$n closes their eyes for a few seconds."), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_ROOM)
			send_to_char(ch.Mindlink, libc.CString("@rYour mind linked target manages to push you out!@n\r\n"))
			if int(ch.Mindlink.Aff_abils.Intel) < axion_dice(-10) && !AFF_FLAGGED(ch.Mindlink, AFF_SHOCKED) {
				send_to_char(ch.Mindlink, libc.CString("Your mind is shocked by the flood of mental energy that pushed it out!@n\r\n"))
				ch.Mindlink.Affected_by[int(AFF_SHOCKED/32)] &= ^(1 << (int(AFF_SHOCKED % 32)))
			}
			ch.Mindlink.Linker = 0
			ch.Mindlink.Mindlink = nil
			ch.Mindlink = nil
			return
		} else {
			act(libc.CString("@rYou struggle to free your mind of @R$N's@r link, but fail!@n"), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_CHAR)
			act(libc.CString("$n closes their eyes for a few seconds, and appears to struggle quite a bit."), FALSE, ch, nil, unsafe.Pointer(ch.Mindlink), TO_ROOM)
			send_to_char(ch.Mindlink, libc.CString("@rYour mind linked target struggles to free their mind, but fails!@n\r\n"))
			ch.Mana -= int64(float64(ch.Mindlink.Max_mana) * 0.05)
			return
		}
	} else if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
		return obj
	}()) == nil {
		send_to_char(ch, libc.CString("Syntax: meditate (object)\nSyntax: meditate expand\nSyntax: meditate break\r\n"))
		return
	}
	if GET_OBJ_VNUM(obj) == 79 {
		send_to_char(ch, libc.CString("It's frozen to the surface.\r\n"))
		return
	}
	weight = obj.Weight
	if obj.Sitting != nil {
		weight += int64(obj.Sitting.Weight)
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity != 10 && ch.Chclass == CLASS_BARDOCK || ch.Chclass != CLASS_BARDOCK {
		cost += int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * 2) + GET_LEVEL(ch)))
	} else if ch.Chclass == CLASS_BARDOCK {
		cost += int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * GET_LEVEL(ch))
	}
	cost += weight * int64(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity+1)/5)
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
		cost += 1000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
		cost += 2000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
		cost += 7500000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
		cost += 15000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
		cost += 25000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
		cost += 50000000
	}
	if cost < weight {
		cost = weight
	}
	if cost == 1 || cost == 0 {
		cost = 25
	}
	if (ch.Bonuses[BONUS_HARDWORKER]) > 0 {
		cost -= int64(float64(cost) * 0.25)
	} else if (ch.Bonuses[BONUS_SLACKER]) > 0 {
		cost += int64(float64(cost) * 0.25)
	}
	if ch.Relax_count >= 464 {
		cost *= 50
	} else if ch.Relax_count >= 232 {
		cost *= 15
	} else if ch.Relax_count >= 116 {
		cost *= 4
	}
	if ch.Mana < cost {
		send_to_char(ch, libc.CString("You don't have enough ki!\r\n"))
		return
	} else {
		act(libc.CString("@cYou close your eyes and concentrate, lifting $p@c with your ki.@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@c$n closes $s eyes and lifts $p@c with $s ki.@n"), TRUE, ch, obj, nil, TO_ROOM)
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 0 {
			bonus = 1
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 {
			bonus = int64(rand_number(2, 4))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 20 {
			bonus = int64(rand_number(5, 10))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 30 {
			bonus = int64(rand_number(10, 15))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 40 {
			bonus = int64(rand_number(15, 20))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 50 {
			bonus = int64(rand_number(40, 60))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 100 {
			bonus = int64(rand_number(180, 250))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 200 {
			bonus = int64(rand_number(400, 600))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
			bonus = int64(rand_number(800, 1200))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
			bonus = int64(rand_number(2000, 3000))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
			bonus = int64(rand_number(4000, 6000))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
			bonus = int64(rand_number(9000, 10000))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
			bonus = int64(rand_number(15000, 20000))
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
			bonus = int64(rand_number(25000, 30000))
		}
		bonus += ((weight + 1) / 500) + 1
		if cost >= ch.Max_mana/2 {
			send_to_char(ch, libc.CString("The object weight and gravity are a great challenge for you!\r\n"))
			bonus *= 10
		} else if cost >= ch.Max_mana/4 {
			send_to_char(ch, libc.CString("The object weight and gravity are an awesome challenge for you!\r\n"))
			bonus *= 8
		} else if cost >= ch.Max_mana/10 {
			send_to_char(ch, libc.CString("The object weight and gravity are a good challenge for you!\r\n"))
			bonus *= 6
		} else if cost < ch.Max_mana/1000 {
			send_to_char(ch, libc.CString("The object weight and gravity are so easy to you, you could do it in your sleep....\r\n"))
			bonus /= 8
		} else if cost < ch.Max_mana/100 {
			send_to_char(ch, libc.CString("The object weight and gravity are the opposite of a challenge for you....\r\n"))
			bonus /= 5
		} else if cost < ch.Max_mana/50 {
			send_to_char(ch, libc.CString("The object weight and gravity are definitely not a challenge for you....\r\n"))
			bonus /= 4
		} else if cost < ch.Max_mana/30 {
			send_to_char(ch, libc.CString("The object weight and gravity are barely a challenge for you....\r\n"))
			bonus /= 3
		} else if cost < ch.Max_mana/20 {
			send_to_char(ch, libc.CString("The object weight and gravity are hardly a challenge for you....\r\n"))
			bonus /= 2
		} else {
			send_to_char(ch, libc.CString("This gravity is just perfect for you....\r\n"))
			bonus *= 4
		}
		if !soft_cap(ch, 1) {
			bonus = 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			send_to_char(ch, libc.CString("@rThis place feels like it operates on a different time frame, it feels great...@n\r\n"))
			bonus *= 10
		} else if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus *= 10
			} else {
				bonus *= 5
			}
		} else if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			send_to_char(ch, libc.CString("@rThis place feels like... Magic.@n\r\n"))
			bonus *= 20
		}
		if bonus <= 0 && !ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 1
		}
		if bonus <= 0 && ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 15
		}
		if bonus <= 1 && ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus = 12
			} else {
				bonus = 6
			}
		}
		if bonus != 1 && ch.Race == RACE_DEMON && rand_number(1, 100) >= 80 {
			send_to_char(ch, libc.CString("Your spirit magnifies the strength of your body! @D[@G+%s@D]@n\r\n"), add_commas(bonus/2))
			ch.Max_hit += bonus / 2
			ch.Basepl += bonus / 2
		}
		bonus += int64(GET_LEVEL(ch) / 20)
		if ch.Race == RACE_NAMEK {
			bonus += bonus / 2
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			bonus += int64(float64(bonus) * 0.1)
		}
		send_to_char(ch, libc.CString("You feel your spirit grow stronger @D[@G+%s@D]@n.\r\n"), add_commas(bonus))
		if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_mana += bonus * 3
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_mana += bonus * 4
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_mana += bonus * 5
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 1 {
			ch.Max_mana += bonus * 2
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 2 {
			ch.Max_mana += bonus * 3
		} else {
			ch.Max_mana += bonus
		}
		ch.Baseki += bonus
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 50 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 200 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 500 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 1000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 5000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 10000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*7)
		}
		ch.Mana -= cost
		if ch.Mana < 0 {
			ch.Mana = 0
		}
	}
}
func do_pushup(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		cost  int = 1
		bonus int = 0
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are a mob fool!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_FISHING) {
		send_to_char(ch, libc.CString("Stop fishing first.\r\n"))
		return
	}
	if ch.Drag != nil {
		send_to_char(ch, libc.CString("You are dragging someone!\r\n"))
		return
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are carrying someone!\r\n"))
		return
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity != 10 && ch.Chclass == CLASS_BARDOCK || ch.Chclass != CLASS_BARDOCK {
		cost += (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * 2) + GET_LEVEL(ch))
	} else if ch.Chclass == CLASS_BARDOCK {
		cost += (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * GET_LEVEL(ch)
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
		cost += 1000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
		cost += 2000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
		cost += 7500000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
		cost += 15000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
		cost += 25000000
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
		cost += 50000000
	}
	if cost == 1 || cost == 0 {
		cost = 25
	}
	if (ch.Bonuses[BONUS_HARDWORKER]) > 0 {
		cost -= int(float64(cost) * 0.25)
	} else if (ch.Bonuses[BONUS_SLACKER]) > 0 {
		cost += int(float64(cost) * 0.25)
	}
	if ch.Relax_count >= 464 {
		cost *= 50
	} else if ch.Relax_count >= 232 {
		cost *= 15
	} else if ch.Relax_count >= 116 {
		cost *= 4
	}
	if ch.Race == RACE_ANDROID || ch.Race == RACE_BIO || ch.Race == RACE_MAJIN || ch.Race == RACE_ARLIAN {
		send_to_char(ch, libc.CString("You will gain nothing from exercising!\r\n"))
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if ch.Move < int64(cost) {
		send_to_char(ch, libc.CString("You are too tired!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_SPAR) {
		send_to_char(ch, libc.CString("You shouldn't be sparring if you want to work out, it could be dangerous.\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are fighting you moron!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_FLYING) {
		send_to_char(ch, libc.CString("You can't do pushups in midair!\r\n"))
		return
	} else {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 0 {
			bonus = 1
			act(libc.CString("@gYou do a pushup.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 {
			bonus = rand_number(3, 7)
			act(libc.CString("@gYou do a pushup, feeling the strain of gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 20 {
			bonus = rand_number(8, 14)
			act(libc.CString("@gYou do a pushup, and feel gravity's pull.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 30 {
			bonus = rand_number(14, 20)
			act(libc.CString("@gYou do a pushup, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 40 {
			bonus = rand_number(20, 35)
			act(libc.CString("@gYou do a pushup, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 50 {
			bonus = rand_number(40, 60)
			act(libc.CString("@gYou do a pushup, and feel the burn.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 100 {
			bonus = rand_number(180, 250)
			act(libc.CString("@gYou do a pushup, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 200 {
			bonus = rand_number(400, 600)
			act(libc.CString("@gYou do a pushup, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 {
			bonus = rand_number(800, 1200)
			act(libc.CString("@gYou do a pushup, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 {
			bonus = rand_number(2000, 3000)
			act(libc.CString("@gYou do a pushup, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 {
			bonus = rand_number(4000, 6000)
			act(libc.CString("@gYou do a pushup, and really strain against the gravity.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 {
			bonus = rand_number(9000, 10000)
			act(libc.CString("@gYou do a pushup, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 {
			bonus = rand_number(15000, 20000)
			act(libc.CString("@gYou do a pushup, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 {
			bonus = rand_number(25000, 30000)
			act(libc.CString("@gYou do a pushup, and it was a really hard one to finish.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@g$n does a pushup, while sweating profusely.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		if cost >= int(ch.Max_hit/2) {
			send_to_char(ch, libc.CString("This gravity is a great challenge for you!\r\n"))
			bonus *= 10
		} else if cost >= int(ch.Max_hit/4) {
			send_to_char(ch, libc.CString("This gravity is an awesome challenge for you!\r\n"))
			bonus *= 8
		} else if cost >= int(ch.Max_hit/10) {
			send_to_char(ch, libc.CString("This gravity is a good challenge for you!\r\n"))
			bonus *= 6
		} else if cost < int(ch.Max_hit/1000) {
			send_to_char(ch, libc.CString("This gravity is so easy to you, you could do it in your sleep....\r\n"))
			bonus /= 8
		} else if cost < int(ch.Max_hit/100) {
			send_to_char(ch, libc.CString("This gravity is the opposite of a challenge for you....\r\n"))
			bonus /= 5
		} else if cost < int(ch.Max_hit/50) {
			send_to_char(ch, libc.CString("This gravity is definitely not a challenge for you....\r\n"))
			bonus /= 4
		} else if cost < int(ch.Max_hit/30) {
			send_to_char(ch, libc.CString("This gravity is barely a challenge for you....\r\n"))
			bonus /= 3
		} else if cost < int(ch.Max_hit/20) {
			send_to_char(ch, libc.CString("This gravity is hardly a challenge for you....\r\n"))
			bonus /= 2
		} else {
			send_to_char(ch, libc.CString("This gravity is just perfect for you....\r\n"))
			bonus *= 4
		}
		if !soft_cap(ch, 0) {
			bonus = 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			send_to_char(ch, libc.CString("@rThis place feels like it operates on a different time frame, it feels great...@n\r\n"))
			bonus *= 10
		} else if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus *= 10
			} else {
				bonus *= 5
			}
		} else if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 19800 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0x4DBB {
			send_to_char(ch, libc.CString("@rThis place feels like... Magic.@n\r\n"))
			bonus *= 20
		}
		if bonus <= 0 && !ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 1
		}
		if bonus <= 0 && ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			bonus = 15
		}
		if bonus <= 1 && ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				bonus = 12
			} else {
				bonus = 6
			}
		}
		bonus += GET_LEVEL(ch) / 20
		if ch.Race == RACE_NAMEK {
			bonus -= bonus / 4
		}
		if (ch.Bonuses[BONUS_HARDWORKER]) != 0 {
			bonus += int(float64(bonus) * 0.5)
		}
		if (ch.Bonuses[BONUS_LONER]) != 0 {
			bonus += int(float64(bonus) * 0.1)
		}
		send_to_char(ch, libc.CString("You feel slightly stronger @D[@G+%s@D]@n.\r\n"), add_commas(int64(bonus)))
		if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS1) {
			ch.Max_hit += int64(bonus * 3)
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS2) {
			ch.Max_hit += int64(bonus * 4)
		} else if ch.Race == RACE_TRUFFLE && PLR_FLAGGED(ch, PLR_TRANS3) {
			ch.Max_hit += int64(bonus * 5)
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 1 {
			ch.Max_hit += int64(bonus * 2)
		} else if ch.Race == RACE_HOSHIJIN && ch.Starphase == 2 {
			ch.Max_hit += int64(bonus * 3)
		} else {
			if ch.Race == RACE_HUMAN {
				bonus = int(float64(bonus) * 0.8)
			}
			ch.Max_hit += int64(bonus)
		}
		ch.Basepl += int64(bonus)
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 50 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 200 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 500 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 1000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 5000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity <= 10000 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*7)
		}
		ch.Move -= int64(cost)
		if ch.Move < 0 {
			ch.Move = 0
		}
	}
}
func do_spar(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	} else {
		if PLR_FLAGGED(ch, PLR_SPAR) {
			act(libc.CString("@wYou cease your sparring stance.@n"), FALSE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w ceases $s sparring stance.@n"), FALSE, ch, nil, nil, TO_ROOM)
		}
		if !PLR_FLAGGED(ch, PLR_SPAR) {
			act(libc.CString("@wYou move into your sparring stance.@n"), FALSE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w moves into $s sparring stance.@n"), FALSE, ch, nil, nil, TO_ROOM)
		}
		ch.Act[int(PLR_SPAR/32)] = ch.Act[int(PLR_SPAR/32)] ^ bitvector_t(1<<(int(PLR_SPAR%32)))
	}
}
func check_eq(ch *char_data) {
	var (
		obj *obj_data
		i   int
	)
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			obj = ch.Equipment[i]
			if OBJ_FLAGGED(obj, ITEM_BROKEN) {
				act(libc.CString("@W$p@W falls apart and you remove it.@n"), FALSE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@W$p@W falls apart and @C$n@W remove it.@n"), FALSE, ch, obj, nil, TO_ROOM)
				perform_remove(ch, i)
				return
			}
			if obj == (ch.Equipment[WEAR_WIELD1]) && (ch.Limb_condition[0]) <= 0 {
				act(libc.CString("@WWithout your right arm you let go of @c$p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W lets go of @c$p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
				perform_remove(ch, i)
				return
			}
			if obj == (ch.Equipment[WEAR_WIELD2]) && (ch.Limb_condition[1]) <= 0 {
				act(libc.CString("@WWithout your left arm you let go of @c$p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
				act(libc.CString("@C$n@W lets go of @c$p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
				perform_remove(ch, i)
				return
			}
		}
	}
}
func base_update() {
	var (
		d       *descriptor_data
		cash    int = FALSE
		inc     int = 0
		countch int = FALSE
		pcoun   int = 0
	)
	if INTERESTTIME != 0 && INTERESTTIME <= C.time(nil) && C.time(nil) != 0 {
		INTERESTTIME = C.time(nil) + 86400
		LASTINTEREST = C.time(nil)
		save_mud_time(&time_info)
		cash = TRUE
		countch = TRUE
	}
	if TOPCOUNTDOWN > 0 {
		TOPCOUNTDOWN -= 4
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if !IS_PLAYING(d) {
			continue
		}
		if IS_NPC(d.Character) {
			if d.Character.Absorbing != nil && d.Character.In_room != d.Character.Absorbing.In_room {
				send_to_char(d.Character, libc.CString("You stop absorbing %s!\r\n"), GET_NAME(d.Character.Absorbing))
				d.Character.Absorbing.Absorbby = nil
				d.Character.Absorbing = nil
			}
			if d.Character.Race == RACE_ANDROID && d.Character.Absorbing != nil && rand_number(1, 10) >= 7 {
				var (
					drain1  int64      = int64(float64(d.Character.Max_mana) * 0.01)
					drain2  int64      = int64(float64(d.Character.Max_move) * 0.01)
					drained *char_data = d.Character.Absorbing
				)
				if drained.Move-drain2 < 0 {
					drain2 = drained.Move
				}
				if drained.Mana-drain1 < 0 {
					drain1 = drained.Mana
				}
				d.Character.Move += drain2
				d.Character.Mana += drain1
				d.Character.Hit += int64(float64(drain1) * 0.5)
				if d.Character.Mana >= d.Character.Max_mana {
					d.Character.Mana = d.Character.Max_mana
				}
				if d.Character.Move >= d.Character.Max_move {
					d.Character.Move = d.Character.Max_move
				}
				if d.Character.Hit >= d.Character.Max_hit {
					d.Character.Hit = d.Character.Max_hit
				}
				if d.Character.Mana == d.Character.Max_mana && d.Character.Move == d.Character.Max_move {
					do_absorb(d.Character, nil, 0, 0)
				}
			}
			continue
		}
		if countch == TRUE {
			pcoun += 1
		}
		if !IS_NPC(d.Character) && rand_number(1, 15) >= 14 {
			ash_burn(d.Character)
		}
		if AFF_FLAGGED(d.Character, AFF_CURSE) && float64(d.Character.Lifeforce) > float64(GET_LIFEMAX(d.Character))*0.4 {
			d.Character.Lifeforce -= int64(float64(GET_LIFEMAX(d.Character)) * 0.01)
			demon_refill_lf(d.Character, int64(float64(GET_LIFEMAX(d.Character))*0.01))
			if float64(d.Character.Lifeforce) < float64(GET_LIFEMAX(d.Character))*0.4 {
				d.Character.Lifeforce = int64(float64(GET_LIFEMAX(d.Character)) * 0.4)
			}
		}
		if d.Character.Backstabcool > 0 {
			d.Character.Backstabcool -= 1
		}
		if PLR_FLAGGED(d.Character, PLR_GOOP) && d.Character.Gooptime == 60 {
			if d.Character.Race == RACE_BIO {
				act(libc.CString("@GConciousness slowly returns to you. You realize quickly that some of your cells have survived. You take control of your regenerative processes and focus on growing a new body!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
			} else {
				act(libc.CString("@MSlowly you regain conciousness. The various split off chunks of your body begin to likewise stir.@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@MYou think you notice the chunks of @m$n@M's moving slightly.@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			}
			d.Character.Gooptime -= 1
		} else if PLR_FLAGGED(d.Character, PLR_GOOP) && d.Character.Gooptime == 30 {
			if d.Character.Race == RACE_BIO {
				act(libc.CString("@GFrom the collection of cells growing a crude form of your body starts to take shape!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@GYou start to notice a large mass of pulsing flesh growing before you!@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			} else {
				act(libc.CString("@MYou will the various chunks of your body to return and slowly more and more of them begin to fly into you. Your body begins to grow larger and larger as this process unfolds!@n "), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@MThe various chunks of @m$n@M's body start to fly into the largest chunk! As the chunks collide they begin to form a larger and still growing blob of goo!@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			}
			d.Character.Gooptime -= 1
		} else if PLR_FLAGGED(d.Character, PLR_GOOP) && d.Character.Gooptime == 15 {
			if d.Character.Race == RACE_BIO {
				act(libc.CString("@GYour body has almost reached its previous form! Only a little more regenerating is needed!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@GThe lump of flesh has now grown to the size where the likeness of @g$n@G can be seen of it! It appears that $e is regenerating $s body from what was only a few cells!@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			} else {
				act(libc.CString("@MYour body has reached half its previous size as your limbs ooze slowly out into their proper shape!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@m$n@M's body has regenerated to half its previous size! Slowly $s limbs ooze out into their proper shape! It won't be long now till $e has fully regenerated!@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			}
			d.Character.Gooptime -= 1
		} else if PLR_FLAGGED(d.Character, PLR_GOOP) && d.Character.Gooptime == 0 {
			if d.Character.Race == RACE_BIO {
				d.Character.Hit = gear_pl(d.Character)
				act(libc.CString("@GYour body has fully regenerated! You flex your arms and legs outward with a rush of renewed strength!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@g$n@G's body has fully regenerated! Suddenly $e flexes $s arms and legs and a rush of power erupts from off of $s body!@n"), TRUE, d.Character, nil, nil, TO_ROOM)
			} else if d.Character.Race == RACE_SAIYAN {
				var (
					zenkaiPL int
					zenkaiKi int
					zenkaiSt int
				)
				zenkaiPL = int(float64(d.Character.Basepl) * 1.03)
				zenkaiKi = int(float64(d.Character.Baseki) * 1.015)
				zenkaiSt = int(float64(d.Character.Basest) * 1.015)
				d.Character.Hit = int64(float64(gear_pl(d.Character)) * 0.5)
				d.Character.Mana = int64(float64(d.Character.Max_mana) * 0.2)
				d.Character.Move = int64(float64(d.Character.Max_move) * 0.2)
				if !IN_ARENA(d.Character) {
					d.Character.Basepl = int64(zenkaiPL)
					d.Character.Baseki = int64(zenkaiKi)
					d.Character.Basest = int64(zenkaiSt)
					d.Character.Max_hit = int64(zenkaiPL)
					d.Character.Max_mana = int64(zenkaiKi)
					d.Character.Max_move = int64(zenkaiSt)
					send_to_char(d.Character, libc.CString("@D[@YZ@ye@wn@Wk@Ya@yi @YB@yo@wo@Ws@Yt@D] @WYou feel much stronger!\r\n"))
					send_to_char(d.Character, libc.CString("@D[@RPL@Y:@n+%s@D] @D[@CKI@Y:@n+%s@D] @D[@GSTA@Y:@n+%s@D]@n\r\n"), add_commas(int64(zenkaiPL)), add_commas(int64(zenkaiKi)), add_commas(int64(zenkaiSt)))
				}
				act(libc.CString("@RYou collapse to the ground, body pushed beyond the typical limits of exhaustion. The passage of time distorts and an indescribable amount of time passes as raw emotions pass through your very being. Your eyes open and focus with a newfound clarity as your unadulterated emotions and feelings revive you for a second wind!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@r$n@R collapses to the ground, seemingly dead. After a brief moment, their eyes flash open with a determined look on their face!"), TRUE, d.Character, nil, nil, TO_ROOM)
			} else {
				d.Character.Hit = gear_pl(d.Character)
				act(libc.CString("@MYour body has fully regenerated! You scream out in triumph and a short gust of steam erupts from your pores!@n"), TRUE, d.Character, nil, nil, TO_CHAR)
				act(libc.CString("@m$n@M's body has fully regenerated! Suddenly $e screams out in gleeful triumph and short gust of steam erupts from $s skin pores!"), TRUE, d.Character, nil, nil, TO_ROOM)
			}
			d.Character.Act[int(PLR_GOOP/32)] &= bitvector_t(^(1 << (int(PLR_GOOP % 32))))
		} else {
			d.Character.Gooptime -= 1
		}
		if d.Character.Con_cooldown > 0 {
			d.Character.Con_cooldown -= 2
			if d.Character.Con_cooldown <= 0 {
				d.Character.Con_cooldown = 0
				send_to_char(d.Character, libc.CString("You can concentrate again.\r\n"))
			}
		}
		if d.Character.Con_sdcooldown > 0 {
			d.Character.Con_sdcooldown -= 10
			if d.Character.Con_sdcooldown <= 0 {
				d.Character.Con_sdcooldown = 0
				send_to_char(d.Character, libc.CString("Your body has recovered from your last selfdestruct.\r\n"))
			}
		}
		if d.Character.Player_specials.Carrying != nil {
			if d.Character.Player_specials.Carrying.In_room != d.Character.In_room {
				carry_drop(d.Character, 3)
			}
		}
		if d.Character.Defender != nil {
			if d.Character.In_room != d.Character.Defender.In_room {
				d.Character.Defender.Defending = nil
				d.Character.Defender = nil
			}
		}
		if d.Character.Defending != nil {
			if d.Character.In_room != d.Character.Defending.In_room {
				d.Character.Defending.Defender = nil
				d.Character.Defending = nil
			}
		}
		if PLR_FLAGGED(d.Character, PLR_TRANSMISSION) {
			d.Character.Act[int(PLR_TRANSMISSION/32)] &= bitvector_t(^(1 << (int(PLR_TRANSMISSION % 32))))
		}
		if d.Character.Fighting == nil && AFF_FLAGGED(d.Character, AFF_POSITION) {
			d.Character.Affected_by[int(AFF_POSITION/32)] &= ^(1 << (int(AFF_POSITION % 32)))
		}
		if d.Character.Sits != nil {
			if d.Character.In_room != d.Character.Sits.In_room {
				var chair *obj_data = d.Character.Sits
				chair.Sitting = nil
				d.Character.Sits = nil
			}
		}
		if d.Character.Ping >= 1 {
			d.Character.Ping -= 1
			if PLR_FLAGGED(d.Character, PLR_PILOTING) && d.Character.Ping == 0 {
				send_to_char(d.Character, libc.CString("Your radar is ready to calculate the direction of another destination.\r\n"))
			}
		}
		if d.Character.Admlevel < 1 && TOPCOUNTDOWN <= 0 && GET_LEVEL(d.Character) > 0 {
			topWrite(d.Character)
		}
		if PLR_FLAGGED(d.Character, PLR_SELFD) && !PLR_FLAGGED(d.Character, PLR_SELFD2) {
			if rand_number(4, 100) < GET_SKILL(d.Character, SKILL_SELFD) {
				send_to_char(d.Character, libc.CString("You feel you are ready to self destruct!\r\n"))
				d.Character.Act[int(PLR_SELFD2/32)] |= bitvector_t(1 << (int(PLR_SELFD2 % 32)))
			}
		}
		if d.Character.Fighting == nil && d.Character.Combo > -1 {
			d.Character.Combo = -1
			d.Character.Combhits = 0
		}
		if MOON_UP != 0 && (d.Character.Race == RACE_SAIYAN || d.Character.Race == RACE_HALFBREED) && !PLR_FLAGGED(d.Character, PLR_OOZARU) {
			oozaru_add(d.Character)
		}
		if cash == TRUE && d.Character.Bank_gold > 0 {
			inc = (d.Character.Bank_gold / 50) * 2
			d.Character.Lastint = LASTINTEREST
			if inc >= 25000 {
				inc = 25000
			}
			d.Character.Bank_gold += inc
			send_to_char(d.Character, libc.CString("@cBank Interest@D: @Y%s@n\r\n"), add_commas(int64(inc)))
		}
		if !IS_NPC(d.Character) {
			check_eq(d.Character)
		}
		if !IS_NPC(d.Character) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Geffect >= 1 && rand_number(1, 100) >= 96 {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Geffect <= 4 {
				switch rand_number(1, 4) {
				case 1:
					act(libc.CString("@RLava spews up violently from the cracks in the ground!@n"), FALSE, d.Character, nil, nil, TO_ROOM)
					act(libc.CString("@RLava spews up violently from the cracks in the ground!@n"), FALSE, d.Character, nil, nil, TO_CHAR)
				case 2:
					act(libc.CString("@RThe lava bubbles and gives off tremendous heat!@n"), FALSE, d.Character, nil, nil, TO_ROOM)
					act(libc.CString("@RThe lava bubbles and gives off tremendous heat!@n"), FALSE, d.Character, nil, nil, TO_CHAR)
				case 3:
					act(libc.CString("@RNoxious fumes rise from the bubbling lava!@n"), FALSE, d.Character, nil, nil, TO_ROOM)
					act(libc.CString("@RNoxious fumes rise from the bubbling lava!@n"), FALSE, d.Character, nil, nil, TO_CHAR)
				case 4:
					act(libc.CString("@RSome of the lava cools as it spreads further from the source!@n"), FALSE, d.Character, nil, nil, TO_ROOM)
					act(libc.CString("@RSome of the lava cools as it spreads further from the source!@n"), FALSE, d.Character, nil, nil, TO_CHAR)
				}
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Geffect += 1
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Geffect == 5 {
				act(libc.CString("@RLava covers the entire area now!@n"), FALSE, d.Character, nil, nil, TO_ROOM)
				act(libc.CString("@RLava covers the entire area now!@n"), FALSE, d.Character, nil, nil, TO_CHAR)
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Geffect += 1
			}
		}
		if d.Character.Absorbing != nil && d.Character.In_room != d.Character.Absorbing.In_room {
			send_to_char(d.Character, libc.CString("You stop absorbing %s!\r\n"), GET_NAME(d.Character.Absorbing))
			d.Character.Absorbing.Absorbby = nil
			d.Character.Absorbing = nil
		}
		if d.Character.Race == RACE_ANDROID && d.Character.Absorbing != nil {
			if d.Character.Absorbing.Move < (d.Character.Max_move/15) && d.Character.Absorbing.Mana < (d.Character.Max_mana/15) {
				act(libc.CString("@WYou stop absorbing stamina and ki from @c$N as they don't have enough for you to take@W!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_CHAR)
				act(libc.CString("@C$n@W stops absorbing stamina and ki from you!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_VICT)
				act(libc.CString("@C$n@W stops absorbing stamina and ki from @c$N@w!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_NOTVICT)
				if d.Character.Fighting == nil || d.Character.Fighting != d.Character.Absorbing {
					set_fighting(d.Character, d.Character.Absorbing.Absorbby)
				}
				if d.Character.Absorbing.Absorbby.Fighting == nil || d.Character.Absorbing.Absorbby.Fighting != d.Character {
					set_fighting(d.Character.Absorbing.Absorbby, d.Character)
				}
				d.Character.Absorbing.Absorbby = nil
				d.Character.Absorbing = nil
			}
		}
		if d.Character.Race == RACE_ANDROID && d.Character.Absorbing != nil && rand_number(1, 9) >= 6 {
			if d.Character.Absorbing.Move > (d.Character.Max_move/15) || d.Character.Absorbing.Mana > (d.Character.Max_mana/15) {
				d.Character.Move += int64(float64(d.Character.Max_move) * 0.08)
				d.Character.Mana += int64(float64(d.Character.Max_mana) * 0.08)
				d.Character.Absorbing.Move -= d.Character.Max_move / 20
				d.Character.Absorbing.Mana -= d.Character.Max_mana / 20
				act(libc.CString("@WYou absorb stamina and ki from @c$N@W!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_CHAR)
				act(libc.CString("@C$n@W absorbs stamina and ki from you!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_VICT)
				send_to_char(d.Character.Absorbing, libc.CString("@wTry 'escape'!@n\r\n"))
				act(libc.CString("@C$n@W absorbs stamina and ki from @c$N@w!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_NOTVICT)
				if d.Character.Hit < gear_pl(d.Character) {
					d.Character.Hit += int64(float64(d.Character.Max_mana) * 0.04)
					if d.Character.Hit > gear_pl(d.Character) {
						d.Character.Hit = gear_pl(d.Character)
					}
					send_to_char(d.Character, libc.CString("@CYou convert a portion of the absorbed energy into refilling your powerlevel.@n\r\n"))
				}
				if d.Character.Absorbing.Move <= 0 {
					d.Character.Absorbing.Move = 1
				}
				if d.Character.Absorbing.Mana <= 0 {
					d.Character.Absorbing.Mana = 1
				}
				if d.Character.Move > d.Character.Max_move && d.Character.Mana < d.Character.Max_mana {
					d.Character.Move = d.Character.Max_move
				} else if d.Character.Move < d.Character.Max_move && d.Character.Mana > d.Character.Max_mana {
					d.Character.Mana = d.Character.Max_mana
				} else if d.Character.Move >= d.Character.Max_move && d.Character.Mana >= d.Character.Max_mana {
					d.Character.Mana = d.Character.Max_mana
					d.Character.Move = d.Character.Max_move
					act(libc.CString("@WYou stop absorbing stamina and ki from @c$N as you are full@W!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_CHAR)
					act(libc.CString("@C$n@W stops absorbing stamina and ki from you!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_VICT)
					act(libc.CString("@C$n@W stops absorbing stamina and ki from @c$N@w!@n"), TRUE, d.Character, nil, unsafe.Pointer(d.Character.Absorbing), TO_NOTVICT)
					if d.Character.Fighting == nil || d.Character.Fighting != d.Character.Absorbing {
						set_fighting(d.Character, d.Character.Absorbing.Absorbby)
					}
					if d.Character.Absorbing.Absorbby.Fighting == nil || d.Character.Absorbing.Absorbby.Fighting != d.Character {
						set_fighting(d.Character.Absorbing.Absorbby, d.Character)
					}
					d.Character.Absorbing.Absorbby = nil
					d.Character.Absorbing = nil
				}
				var sum int = 1
				var mum int = 1
				var ium int = 1
				if !soft_cap(d.Character, 0) {
					sum = 0
				}
				if !soft_cap(d.Character, 2) {
					mum = 0
				}
				if !soft_cap(d.Character, 1) {
					ium = 0
				}
				if sum == 1 {
					if rand_number(1, 8) >= 6 {
						var gain int = rand_number(GET_LEVEL(d.Character)/2, GET_LEVEL(d.Character)*3) + GET_LEVEL(d.Character)*18
						if GET_LEVEL(d.Character) > 30 {
							gain += rand_number(GET_LEVEL(d.Character)*2, GET_LEVEL(d.Character)*4) + GET_LEVEL(d.Character)*50
						}
						if GET_LEVEL(d.Character) > 60 {
							gain *= 2
						}
						if GET_LEVEL(d.Character) > 80 {
							gain *= 3
						}
						if GET_LEVEL(d.Character) > 90 {
							gain *= 4
						}
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent powerlevel!@n\r\n"), gain)
						if group_bonus(d.Character, 2) == 7 {
							if PLR_FLAGGED(d.Character.Master, PLR_SENSEM) {
								var gbonus int = int(float64(gain) * 0.15)
								gain += gbonus
								send_to_char(d.Character, libc.CString("The leader of your group conveys an extra bonus! @D[@G+%s@D]@n \r\n"), add_commas(int64(gbonus)))
							}
						}
						d.Character.Max_hit += int64(gain)
						d.Character.Basepl += int64(gain)
					}
				}
				if mum == 1 {
					if rand_number(1, 8) >= 6 {
						var gain int = rand_number(GET_LEVEL(d.Character)/2, GET_LEVEL(d.Character)*3) + GET_LEVEL(d.Character)*18
						if GET_LEVEL(d.Character) > 30 {
							gain += rand_number(GET_LEVEL(d.Character)*2, GET_LEVEL(d.Character)*4) + GET_LEVEL(d.Character)*50
						}
						if GET_LEVEL(d.Character) > 60 {
							gain *= 2
						}
						if GET_LEVEL(d.Character) > 80 {
							gain *= 3
						}
						if GET_LEVEL(d.Character) > 90 {
							gain *= 4
						}
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent stamina!@n\r\n"), gain)
						if group_bonus(d.Character, 2) == 7 {
							if PLR_FLAGGED(d.Character.Master, PLR_SENSEM) {
								var gbonus int = int(float64(gain) * 0.15)
								gain += gbonus
								send_to_char(d.Character, libc.CString("The leader of your group conveys an extra bonus! @D[@G+%s@D]@n \r\n"), add_commas(int64(gbonus)))
							}
						}
						d.Character.Max_move += int64(gain)
						d.Character.Basest += int64(gain)
					}
				}
				if ium == 1 {
					if rand_number(1, 8) >= 6 {
						var gain int = rand_number(GET_LEVEL(d.Character)/2, GET_LEVEL(d.Character)*3) + GET_LEVEL(d.Character)*18
						if GET_LEVEL(d.Character) > 30 {
							gain += rand_number(GET_LEVEL(d.Character)*2, GET_LEVEL(d.Character)*4) + GET_LEVEL(d.Character)*50
						}
						if GET_LEVEL(d.Character) > 60 {
							gain *= 2
						}
						if GET_LEVEL(d.Character) > 80 {
							gain *= 3
						}
						if GET_LEVEL(d.Character) > 90 {
							gain *= 4
						}
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent ki!@n\r\n"), gain)
						if group_bonus(d.Character, 2) == 7 {
							if PLR_FLAGGED(d.Character.Master, PLR_SENSEM) {
								var gbonus int = int(float64(gain) * 0.15)
								gain += gbonus
								send_to_char(d.Character, libc.CString("The leader of your group conveys an extra bonus! @D[@G+%s@D]@n \r\n"), add_commas(int64(gbonus)))
							}
						}
						d.Character.Max_mana += int64(gain)
						d.Character.Baseki += int64(gain)
					}
				}
				if sum == 0 {
					if rand_number(1, 8) >= 6 {
						var gain int = 1
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent powerlevel. You may need to level.@n\r\n"), gain)
						d.Character.Max_hit += int64(gain)
						d.Character.Basepl += int64(gain)
					}
				}
				if mum == 0 {
					if rand_number(1, 8) >= 6 {
						var gain int = 1
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent stamina. You may need to level.@n\r\n"), gain)
						d.Character.Max_move += int64(gain)
						d.Character.Basest += int64(gain)
					}
				}
				if ium == 0 {
					if rand_number(1, 8) >= 6 {
						var gain int = 1
						send_to_char(d.Character, libc.CString("@gYou gain +@G%d@g permanent ki. You may need to level.@n\r\n"), gain)
						d.Character.Max_mana += int64(gain)
						d.Character.Baseki += int64(gain)
					}
				}
			}
		}
		if d.Character.Blocks != nil {
			var vict *char_data = d.Character.Blocks
			if vict.In_room != d.Character.In_room {
				vict.Blocked = nil
				d.Character.Blocks = nil
			}
		}
		if d.Character.Overf == TRUE {
			mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("OVERFLOW: %s has caused an overflow, check for illegal activity."), GET_NAME(d.Character))
			d.Character.Overf = FALSE
		}
		if d.Character.Spam > 0 {
			d.Character.Spam = 0
		} else {
			continue
		}
	}
	if countch == TRUE {
		PCOUNT = pcoun
		PCOUNTDAY = C.time(nil)
	}
	if TOPCOUNTDOWN <= 0 {
		TOPCOUNTDOWN = 60
	}
}
func has_scanner(ch *char_data) int {
	var (
		obj      *obj_data
		next_obj *obj_data
		success  int = 0
	)
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if obj != nil && GET_OBJ_VNUM(obj) == 13600 {
			success = 1
		}
	}
	return success
}
func do_snet(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		channel int = 0
		global  int = FALSE
		call    int = -1
		reached int = FALSE
		i       *descriptor_data
		voice   [150]byte
		arg     [2048]byte
		arg2    [2048]byte
		hist    [2048]byte
	)
	half_chop(argument, &arg[0], &arg2[0])
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	if ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
		send_to_char(ch, libc.CString("This is a different dimension!\r\n"))
		return
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("Lol, no.\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
		send_to_char(ch, libc.CString("This is the past, you can't talk on scouter net!\r\n"))
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
		send_to_char(ch, libc.CString("The fire eats your transmission!\r\n"))
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 19800 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 0x4DBB {
		send_to_char(ch, libc.CString("Your signal will not be able to escape the walls of the pocket dimension.\r\n"))
		return
	}
	if !IS_NPC(ch) {
		if (ch.Equipment[WEAR_EYE]) != nil {
			obj = ch.Equipment[WEAR_EYE]
		} else {
			send_to_char(ch, libc.CString("You do not have a scouter on.\r\n"))
			return
		}
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("[Syntax] snet < [1-999] | check | #(scouter number) | * message | message>\r\n"))
		return
	}
	if C.strstr(&arg[0], libc.CString("#")) != nil {
		search_replace(&arg[0], libc.CString("#"), libc.CString(""))
		call = libc.Atoi(libc.GoString(&arg[0]))
		if call <= -1 {
			send_to_char(ch, libc.CString("Call what personal scouter number?\r\n"))
			return
		}
	}
	if C.strcasecmp(&arg[0], libc.CString("check")) == 0 {
		send_to_char(ch, libc.CString("Your personal scouter number is: %d\r\n"), ch.Id)
		return
	}
	if call <= -1 {
		channel = libc.Atoi(libc.GoString(&arg[0]))
	}
	if channel > 0 {
		obj.Scoutfreq = channel
		if channel > 999 {
			obj.Scoutfreq = 999
		}
		act(libc.CString("@wYou push some buttons on $p@w and change its channel."), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@w pushes some buttons on $p@w and changes its channel."), TRUE, ch, obj, nil, TO_ROOM)
		return
	} else {
		if (ch.Bonuses[BONUS_MUTE]) > 0 {
			send_to_char(ch, libc.CString("You are unable to speak though.\r\n"))
			return
		}
		if obj.Scoutfreq == 0 {
			obj.Scoutfreq = 1
		}
		if C.strcasecmp(&arg[0], libc.CString("*")) == 0 && call <= -1 {
			global = TRUE
		}
		if ch.Voice != nil {
			stdio.Sprintf(&voice[0], "%s", ch.Voice)
		}
		if ch.Voice == nil {
			stdio.Sprintf(&voice[0], "A generic voice")
		}
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Connected != CON_PLAYING {
				continue
			}
			if i.Character == ch {
				continue
			}
			if i.Character.In_room == ch.In_room {
				continue
			}
			if ROOM_FLAGGED(i.Character.In_room, ROOM_HBTC) {
				continue
			}
			if ROOM_FLAGGED(i.Character.In_room, ROOM_PAST) {
				continue
			}
			if ROOM_FLAGGED(i.Character.In_room, ROOM_RHELL) && !ROOM_FLAGGED(ch.In_room, ROOM_RHELL) || ROOM_FLAGGED(i.Character.In_room, ROOM_AL) && !ROOM_FLAGGED(ch.In_room, ROOM_AL) {
				continue
			}
			if !ROOM_FLAGGED(i.Character.In_room, ROOM_RHELL) && ROOM_FLAGGED(ch.In_room, ROOM_RHELL) || !ROOM_FLAGGED(i.Character.In_room, ROOM_AL) && ROOM_FLAGGED(ch.In_room, ROOM_AL) {
				continue
			}
			if i.Character.Position == POS_SLEEPING {
				continue
			}
			if (func() room_vnum {
				if i.Character.In_room != room_rnum(-1) && i.Character.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.Character.In_room)))).Number
				}
				return -1
			}()) >= 19800 && (func() room_vnum {
				if i.Character.In_room != room_rnum(-1) && i.Character.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.Character.In_room)))).Number
				}
				return -1
			}()) <= 0x4DBB {
				continue
			}
			if (i.Character.Equipment[WEAR_EYE]) != nil {
				obj2 = i.Character.Equipment[WEAR_EYE]
				if obj2.Scoutfreq == 0 {
					obj2.Scoutfreq = 1
				}
				if global == FALSE && call <= -1 && obj2.Scoutfreq == obj.Scoutfreq && i.Character.Admlevel < 1 {
					send_to_char(i.Character, libc.CString("@C%s is heard @W(@c%s@W), @D[@WSNET FREQ@D: @Y%d@D] @G%s %s@n\r\n"), &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), obj.Scoutfreq, CAP(&arg[0]), func() string {
						if arg2[0] == 0 {
							return ""
						}
						return libc.GoString(&arg2[0])
					}())
					hist[0] = '\x00'
					stdio.Sprintf(&hist[0], "@C%s is heard @W(@c%s@W), @D[@WSNET FREQ@D: @Y%d@D] @G%s %s@n\r\n", &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), obj.Scoutfreq, CAP(&arg[0]), func() string {
						if arg2[0] == 0 {
							return ""
						}
						return libc.GoString(&arg2[0])
					}())
					add_history(i.Character, &hist[0], HIST_SNET)
					if has_scanner(i.Character) != 0 {
						var blah *byte = sense_location(ch)
						send_to_char(i.Character, libc.CString("@WScanner@D: @Y%s@n\r\n"), blah)
						libc.Free(unsafe.Pointer(blah))
					}
					continue
				} else if global == TRUE && call <= -1 && i.Character.Admlevel < 1 {
					send_to_char(i.Character, libc.CString("@C%s is heard @W(@c%s@W), @D[@WSNET FREQ@D: @Y%d @mBroadcast@D] @G%s@n\r\n"), &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), obj.Scoutfreq, CAP(&arg2[0]))
					hist[0] = '\x00'
					stdio.Sprintf(&hist[0], "@C%s is heard @W(@c%s@W), @D[@WSNET FREQ@D: @Y%d @mBroadcast@D] @G%s@n\r\n", &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), obj.Scoutfreq, CAP(&arg2[0]))
					add_history(i.Character, &hist[0], HIST_SNET)
					if has_scanner(i.Character) != 0 {
						var blah *byte = sense_location(ch)
						send_to_char(i.Character, libc.CString("@WScanner@D: @Y%s@n\r\n"), blah)
						libc.Free(unsafe.Pointer(blah))
					}
					continue
				} else if call > -1 && i.Character.Id == int32(call) {
					send_to_char(i.Character, libc.CString("@C%s is heard @W(@c%s@W), @D[@R#@W%d @Ycalling YOU@D] @G%s@n\r\n"), &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), ch.Id, func() string {
						if arg2[0] == 0 {
							return ""
						}
						return libc.GoString(CAP(&arg2[0]))
					}())
					hist[0] = '\x00'
					stdio.Sprintf(&hist[0], "@C%s is heard @W(@c%s@W), @D[@R#@W%d @Ycalling YOU@D] @G%s@n\r\n", &voice[0], func() *byte {
						if readIntro(i.Character, ch) == 1 {
							return get_i_name(i.Character, ch)
						}
						return libc.CString("Unknown")
					}(), ch.Id, func() string {
						if arg2[0] == 0 {
							return ""
						}
						return libc.GoString(CAP(&arg2[0]))
					}())
					add_history(i.Character, &hist[0], HIST_SNET)
					if has_scanner(i.Character) != 0 {
						var blah *byte = sense_location(ch)
						send_to_char(i.Character, libc.CString("@WScanner@D: @Y%s@n\r\n"), blah)
						libc.Free(unsafe.Pointer(blah))
					}
					reached = TRUE
				}
			}
			if i.Character.Admlevel > 0 && call <= -1 {
				send_to_char(i.Character, libc.CString("@C%s (%s) is heard, @D[@WSNET FREQ@D: @Y%d@D] @G%s %s@n\r\n"), &voice[0], GET_NAME(ch), obj.Scoutfreq, CAP(&arg[0]), func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(&arg2[0])
				}())
				continue
			} else if i.Character.Admlevel > 0 {
				send_to_char(i.Character, libc.CString("@C%s (%s) is heard, @D[@WCall to @R#@Y%d@D] @G%s@n\r\n"), &voice[0], GET_NAME(ch), call, func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(CAP(&arg2[0]))
				}())
				continue
			}
		}
		if call <= -1 {
			if global == FALSE {
				reveal_hiding(ch, 3)
				send_to_char(ch, libc.CString("@CYou @D[@WSNET FREQ@D: @Y%d@D] @G%s %s@n\r\n"), obj.Scoutfreq, &arg[0], func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(&arg2[0])
				}())
				hist[0] = '\x00'
				stdio.Sprintf(&hist[0], "@CYou @D[@WSNET FREQ@D: @Y%d@D] @G%s %s@n\r\n", obj.Scoutfreq, &arg[0], func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(&arg2[0])
				}())
				add_history(ch, &hist[0], HIST_SNET)
				var over [64936]byte
				stdio.Sprintf(&over[0], "@C$n@W says into $s scouter, '@G@G%s %s@W'@n\r\n", CAP(&arg[0]), func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(&arg2[0])
				}())
				act(&over[0], TRUE, ch, nil, nil, TO_ROOM)
				if ROOM_FLAGGED(ch.In_room, ROOM_RHELL) || ROOM_FLAGGED(ch.In_room, ROOM_AL) {
					send_to_char(ch, libc.CString("@mThe transmission only reaches those who are in the afterlife.@n\r\n"))
				}
			}
			if global == TRUE {
				reveal_hiding(ch, 3)
				send_to_char(ch, libc.CString("@CYou @D[@WSNET FREQ@D: @Y%d @mBroadcast@D] @G%s@n\r\n"), obj.Scoutfreq, func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(CAP(&arg2[0]))
				}())
				hist[0] = '\x00'
				stdio.Sprintf(&hist[0], "@CYou @D[@WSNET FREQ@D: @Y%d @mBroadcast@D] @G%s@n\r\n", obj.Scoutfreq, func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(CAP(&arg2[0]))
				}())
				add_history(ch, &hist[0], HIST_SNET)
				var over [64936]byte
				stdio.Sprintf(&over[0], "@C$n@W says into $s scouter, '@G@G%s@W'@n\r\n", func() string {
					if arg2[0] == 0 {
						return ""
					}
					return libc.GoString(CAP(&arg2[0]))
				}())
				act(&over[0], TRUE, ch, nil, nil, TO_ROOM)
				if ROOM_FLAGGED(ch.In_room, ROOM_RHELL) || ROOM_FLAGGED(ch.In_room, ROOM_AL) {
					send_to_char(ch, libc.CString("@mThe transmission only reaches those who are in the afterlife.@n\r\n"))
				}
			}
		} else {
			reveal_hiding(ch, 3)
			send_to_char(ch, libc.CString("@CYou call @D[@R#@W%d@D] @G%s@n\r\n"), call, func() string {
				if arg2[0] == 0 {
					return ""
				}
				return libc.GoString(CAP(&arg2[0]))
			}())
			hist[0] = '\x00'
			stdio.Sprintf(&hist[0], "@CYou call @D[@R#@W%d@D] @G%s@n\r\n", call, func() string {
				if arg2[0] == 0 {
					return ""
				}
				return libc.GoString(CAP(&arg2[0]))
			}())
			add_history(ch, &hist[0], HIST_SNET)
			var over [64936]byte
			stdio.Sprintf(&over[0], "@C$n@W says into $s scouter, '@G@G%s@W'@n\r\n", func() string {
				if arg2[0] == 0 {
					return ""
				}
				return libc.GoString(CAP(&arg2[0]))
			}())
			act(&over[0], TRUE, ch, nil, nil, TO_ROOM)
			if reached == FALSE {
				send_to_char(ch, libc.CString("@mThe transmission didn't reach them.@n\r\n"))
			}
		}
	}
}
func do_scouter(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict  *char_data = nil
		i     *descriptor_data
		count int = 0
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	}
	var obj *obj_data = (ch.Equipment[WEAR_EYE])
	if obj == nil {
		send_to_char(ch, libc.CString("You do not even have a scouter!"))
		obj = nil
		return
	} else {
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("[Syntax] scouter < target | scan>\r\n"))
			return
		}
		reveal_hiding(ch, 3)
		if C.strcasecmp(libc.CString("scan"), &arg[0]) == 0 {
			for i = descriptor_list; i != nil; i = i.Next {
				if i.Connected != CON_PLAYING {
					continue
				} else if i.Character == ch {
					continue
				} else if i.Character.Race == RACE_ANDROID {
					continue
				} else if planet_check(ch, i.Character) != 0 {
					var (
						dir     int = find_first_step(ch.In_room, i.Character.In_room)
						same    int = FALSE
						pathway [64936]byte
					)
					if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone)))).Number == (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.Character.In_room)))).Zone)))).Number {
						same = TRUE
					}
					switch dir {
					case (-1):
						stdio.Sprintf(&pathway[0], "@rERROR")
					case (-2):
						stdio.Sprintf(&pathway[0], "@RHERE")
					case (-4):
						send_to_char(ch, libc.CString("@MUNKNOWN"))
					default:
						send_to_char(ch, libc.CString("@G%s\r\n"), dirs[dir])
					}
					var blah *byte = sense_location(i.Character)
					if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && i.Character.Hit >= 150000 {
						send_to_char(ch, libc.CString("@D<@GPowerlevel Detected@D:@w ?????????@D> @w---> @C%s@n\r\n"), func() *byte {
							if same == TRUE {
								return &pathway[0]
							}
							return blah
						}())
					} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && i.Character.Hit >= 5000000 {
						send_to_char(ch, libc.CString("@D<@GPowerlevel Detected@D:@w ?????????@D> @w---> @C%s@n\r\n"), func() *byte {
							if same == TRUE {
								return &pathway[0]
							}
							return blah
						}())
					} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && i.Character.Hit >= 15000000 {
						send_to_char(ch, libc.CString("@D<@GPowerlevel Detected@D:@w ?????????@D> @w---> @C%s@n\r\n"), func() *byte {
							if same == TRUE {
								return &pathway[0]
							}
							return blah
						}())
					} else {
						send_to_char(ch, libc.CString("@D<@GPowerlevel Detected@D: [@Y%s@D]@w ---> @C%s@n\r\n"), add_commas(i.Character.Hit), func() *byte {
							if same == TRUE {
								return &pathway[0]
							}
							return blah
						}())
					}
					count++
					libc.Free(unsafe.Pointer(blah))
				}
			}
			if count == 0 {
				send_to_char(ch, libc.CString("You didn't detect anyone of notice.\r\n"))
				return
			} else if count >= 1 {
				send_to_char(ch, libc.CString("%d powerlevels detected.\r\n"), count)
				return
			}
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("They don't seem to be here.\r\n"))
			return
		}
		if vict.Race == RACE_ANDROID {
			act(libc.CString("$n points $s scouter at you."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n points $s scouter at $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			send_to_char(ch, libc.CString("@D,==================================|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1@RReading target...                 @n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1@RP@r@1o@Rw@r@1e@1@Rr L@r@1e@Rv@r@1e@1@Rl@1@D:                 @RERROR@n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1@CC@c@1ha@1@Cr@c@1ge@1@Cd Ki @1@D:                 @RERROR@n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D:                 @RERROR@n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1@GE@g@1x@Gt@g@1r@Ga I@g@1nf@Go @D:                 @RERROR@n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
			send_to_char(ch, libc.CString("@D`==================================|@n\r\n"))
			return
		} else {
			if OBJ_FLAGGED(obj, ITEM_BSCOUTER) && vict.Hit >= 150000 {
				act(libc.CString("$n points $s scouter at you."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n points $s scouter at $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				perform_remove(ch, WEAR_EYE)
				send_to_char(ch, libc.CString("Your scouter overloads and explodes!\r\n"))
				act(libc.CString("$n's scouter explodes!"), FALSE, ch, nil, nil, TO_ROOM)
				extract_obj(obj)
				save_char(ch)
				return
			} else if OBJ_FLAGGED(obj, ITEM_MSCOUTER) && vict.Hit >= 5000000 {
				act(libc.CString("$n points $s scouter at you."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n points $s scouter at $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				perform_remove(ch, WEAR_EYE)
				send_to_char(ch, libc.CString("Your scouter overloads and explodes!\r\n"))
				act(libc.CString("$n's scouter explodes!"), FALSE, ch, nil, nil, TO_ROOM)
				extract_obj(obj)
				save_char(ch)
				return
			} else if OBJ_FLAGGED(obj, ITEM_ASCOUTER) && vict.Hit >= 15000000 {
				act(libc.CString("$n points $s scouter at you."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n points $s scouter at $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				perform_remove(ch, WEAR_EYE)
				send_to_char(ch, libc.CString("Your scouter overloads and explodes!\r\n"))
				act(libc.CString("$n's scouter explodes!"), FALSE, ch, nil, nil, TO_ROOM)
				extract_obj(obj)
				save_char(ch)
				return
			} else {
				var (
					percent float64 = 0.0
					cur     float64 = 0.0
					max     float64 = 0.0
					stam    int64   = vict.Move
					mstam   int64   = vict.Max_move
				)
				if stam <= 0 {
					stam = 1
				}
				if mstam <= 0 {
					mstam = 1
				}
				cur = float64(stam)
				max = float64(mstam)
				percent = (cur / max) * 100
				act(libc.CString("$n points $s scouter at you."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("$n points $s scouter at $N."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				send_to_char(ch, libc.CString("@D,==================================|@n\r\n"))
				send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
				send_to_char(ch, libc.CString("@D|@1@RReading target...                 @n@D|@n\r\n"))
				send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
				send_to_char(ch, libc.CString("@D|@1@RP@r@1o@Rw@r@1e@1@Rr L@r@1e@Rv@r@1e@1@Rl@1@D: @Y%21s@n@D|@n\r\n"), add_commas(vict.Hit))
				if !IS_NPC(vict) {
					send_to_char(ch, libc.CString("@D|@1@CC@c@1ha@1@Cr@c@1ge@1@Cd Ki @1@D: @Y%21s@n@D|@n\r\n"), add_commas(vict.Charge))
				} else if IS_NPC(vict) {
					send_to_char(ch, libc.CString("@D|@1@CC@c@1ha@1@Cr@c@1ge@1@Cd Ki @1@D: @Y%21s@n@D|@n\r\n"), add_commas(int64(vict.Mobcharge*rand_number(GET_LEVEL(ch)*50, GET_LEVEL(ch)*200))))
				}
				if percent < 10 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Exhausted")
				} else if percent < 25 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Extremely Tired")
				} else if percent < 50 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Very Tired")
				} else if percent < 75 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Tired")
				} else if percent < 90 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Winded")
				} else if percent < 100 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Untired")
				} else if percent >= 100 {
					send_to_char(ch, libc.CString("@D|@1@YS@y@1ta@1@Ym@y@1in@1@Ya    @1@D: @Y%21s@n@D|@n\r\n"), "Energetic")
				}
				send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
				var check int = FALSE
				send_to_char(ch, libc.CString("@D|@1@GE@g@1x@Gt@g@1r@Ga I@g@1nf@Go @D: "))
				if AFF_FLAGGED(vict, AFF_ZANZOKEN) {
					send_to_char(ch, libc.CString("@Y%21s@n@D|@n\n"), "Zanzoken Prepared")
					check = TRUE
				}
				if AFF_FLAGGED(vict, AFF_HASS) {
					send_to_char(ch, libc.CString("%s@Y%21s@n@D|@n\n"), func() string {
						if check == TRUE {
							return "@D|@1             "
						}
						return ""
					}(), "Accelerated Arms")
					check = TRUE
				}
				if AFF_FLAGGED(vict, AFF_HEALGLOW) {
					send_to_char(ch, libc.CString("%s@Y%21s@n@D|@n\n"), func() string {
						if check == TRUE {
							return "@D|@1             "
						}
						return ""
					}(), "Healing Glow Prepared")
					check = TRUE
				}
				if AFF_FLAGGED(vict, AFF_POISON) {
					send_to_char(ch, libc.CString("%s@Y%21s@n@D|@n\n"), func() string {
						if check == TRUE {
							return "@D|@1             "
						}
						return ""
					}(), "Poisoned")
					check = TRUE
				}
				if PLR_FLAGGED(vict, PLR_SELFD) {
					send_to_char(ch, libc.CString("%s@Y%21s@n@D|@n\n"), func() string {
						if check == TRUE {
							return "@D|@1             "
						}
						return ""
					}(), "Explosive Energy")
					check = TRUE
				}
				if check == FALSE {
					send_to_char(ch, libc.CString("%s@Y%21s@n@D|@n\n"), func() string {
						if check == TRUE {
							return "@D|@1             "
						}
						return ""
					}(), "None Detected.")
				}
				send_to_char(ch, libc.CString("@D|@1                                  @n@D|@n\r\n"))
				send_to_char(ch, libc.CString("@D`==================================|@n\r\n"))
			}
		}
	}
}
func dball_count(ch *char_data) int {
	var (
		dball    [7]int    = [7]int{20, 21, 22, 23, 24, 25, 26}
		count    int       = 0
		obj      *obj_data = nil
		next_obj *obj_data = nil
	)
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if GET_OBJ_VNUM(obj) == obj_vnum(dball[0]) {
			dball[0] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[1]) {
			dball[1] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[2]) {
			dball[2] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[3]) {
			dball[3] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[4]) {
			dball[4] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[5]) {
			dball[5] = -1
			count++
			continue
		} else if GET_OBJ_VNUM(obj) == obj_vnum(dball[6]) {
			dball[6] = -1
			count++
			continue
		} else {
			continue
		}
	}
	if count >= 1 {
		return 1
	} else {
		return 0
	}
}
func do_quit(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) || ch.Desc == nil {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
		send_to_char(ch, libc.CString("This is the past, you can't quit here!\r\n"))
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 2002 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 2011 {
		send_to_char(ch, libc.CString("You can't quit in the arena!\r\n"))
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 101 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 139 {
		send_to_char(ch, libc.CString("You can't quit in the mud school!\r\n"))
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) >= 19800 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) <= 0x4DBB {
		send_to_char(ch, libc.CString("You can't quit in a pocket dimension!\r\n"))
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 2069 {
		send_to_char(ch, libc.CString("You can't quit here!\r\n"))
		return
	}
	if ch.Mindlink != nil && ch.Linker == 0 {
		send_to_char(ch, libc.CString("@RYou feel like the mind that is linked with yours is preventing you from quiting!@n\r\n"))
		if ch.Mindlink.In_room != room_rnum(-1) {
			look_at_room(ch.Mindlink.In_room, ch, 0)
			send_to_char(ch, libc.CString("You get an impression of where this interference is originating from.\r\n"))
		}
		return
	}
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 2070 {
		send_to_char(ch, libc.CString("You can't quit here!\r\n"))
		return
	}
	if dball_count(ch) != 0 {
		send_to_char(ch, libc.CString("You can not quit while you have dragon balls! Place them somewhere first."))
		return
	}
	if subcmd != SCMD_QUIT {
		send_to_char(ch, libc.CString("You have to type quit--no less, to quit!\r\n"))
	} else if ch.Position == POS_FIGHTING {
		send_to_char(ch, libc.CString("No way!  You're fighting for your life!\r\n"))
	} else if ch.Position < POS_STUNNED {
		send_to_char(ch, libc.CString("You die before your time...\r\n"))
		die(ch, nil)
	} else {
		act(libc.CString("$n has left the game."), TRUE, ch, nil, nil, TO_ROOM)
		mudlog(NRM, MAX(ADMLVL_IMMORT, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("%s has quit the game."), GET_NAME(ch))
		send_to_char(ch, libc.CString("Goodbye, friend.. Come back soon!\r\n"))
		if ch.Followers != nil || ch.Master != nil {
			die_follower(ch)
		}
		if ch == ch_selling {
			stop_auction(AUC_QUIT_CANCEL, nil)
		}
		if !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && ((func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) < 19800 || (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) > 0x4DBB) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != room_vnum(-1) && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != 0 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != 1 {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				} else {
					ch.Player_specials.Load_room = -1
				}
			}
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != room_vnum(-1) && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != 0 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) != 1 {
				if real_room(1561) != room_rnum(-1) && real_room(1561) <= top_of_world {
					ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(1561))))).Number
				} else {
					ch.Player_specials.Load_room = -1
				}
			}
		}
		Crash_rentsave(ch, 0)
		cp(ch)
		extract_char(ch)
	}
	if ch.Desc.Snoop_by != nil {
		write_to_output(ch.Desc.Snoop_by, libc.CString("Your victim is no longer among us.\r\n"))
		ch.Desc.Snoop_by.Snooping = nil
		ch.Desc.Snoop_by = nil
	}
}
func do_save(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) || ch.Desc == nil {
		return
	}
	if cmd != 0 {
		if config_info.Csd.Auto_save != 0 && ch.Admlevel < 1 {
			send_to_char(ch, libc.CString("Saving.\r\n"))
			write_aliases(ch)
			save_char(ch)
			Crash_crashsave(ch)
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) < 19800 || (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) > 0x4DBB {
				if (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) != room_vnum(-1) && (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) != 0 && (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) != 1 {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					} else {
						ch.Player_specials.Load_room = -1
					}
				}
			}
			return
		}
		send_to_char(ch, libc.CString("Saving.\r\n"))
	}
	write_aliases(ch)
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) < 19800 || (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) > 0x4DBB {
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != room_vnum(-1) && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != 0 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != 1 {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			} else {
				ch.Player_specials.Load_room = -1
			}
		}
	}
	save_char(ch)
	Crash_crashsave(ch)
	cp(ch)
}
func do_not_here(ch *char_data, argument *byte, cmd int, subcmd int) {
	send_to_char(ch, libc.CString("Sorry, but you cannot do that here!\r\n"))
}
func do_steal(ch *char_data, argument *byte, cmd int, subcmd int) {
	if GET_SKILL(ch, SKILL_SLEIGHT_OF_HAND) == 0 && slot_count(ch)+1 <= ch.Skill_slots {
		send_to_char(ch, libc.CString("You learn the very veeeery basics of theft. Which is don't get caught.\r\n"))
		for {
			ch.Skills[SKILL_SLEIGHT_OF_HAND] = 1
			if true {
				break
			}
		}
	} else if GET_SKILL(ch, SKILL_SLEIGHT_OF_HAND) == 0 && slot_count(ch)+1 > ch.Skill_slots {
		send_to_char(ch, libc.CString("You can't learn any more skills and thus can not steal right now!\r\n"))
		return
	}
	var vict *char_data
	var obj *obj_data
	var arg [500]byte
	var arg2 [500]byte
	var gold int = 0
	var prob int = GET_SKILL(ch, SKILL_SLEIGHT_OF_HAND)
	var perc int = 0
	var eq_pos int
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("An important basic of theft is actually having a victim!\r\n"))
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("Steal what from who?\r\n"))
		return
	} else if vict == ch {
		send_to_char(ch, libc.CString("Come on now, that's rather stupid!\r\n"))
		return
	} else if can_kill(ch, vict, nil, 0) == 0 {
		return
	} else if GET_LEVEL(ch) <= 8 {
		send_to_char(ch, libc.CString("You are trapped inside the newbie shield until level 9 and can't piss off those bigger and better than you. Awww...\r\n"))
		return
	} else if MOB_FLAGGED(vict, MOB_NOKILL) && ch.Admlevel == ADMLVL_NONE {
		send_to_char(ch, libc.CString("That isn't such a good idea...\r\n"))
		return
	}
	if ch.Move < (ch.Max_move/40)+int64(ch.Carry_weight) {
		send_to_char(ch, libc.CString("You do not have enough stamina.\r\n"))
		return
	}
	if !IS_NPC(vict) && GET_SKILL(vict, SKILL_SPOT) != 0 {
		perc = GET_SKILL(vict, SKILL_SPOT)
		perc += int(vict.Aff_abils.Intel)
	} else {
		perc = rand_number(int(vict.Aff_abils.Intel), int(vict.Aff_abils.Intel+10))
		if IS_NPC(vict) {
			perc += int(float64(GET_LEVEL(vict)) * 0.25)
		}
	}
	if vict.Position == POS_SITTING {
		perc -= 5
	}
	if vict.Position == POS_RESTING {
		perc -= 10
	}
	if vict.Position <= POS_SLEEPING {
		perc -= 25
	}
	prob += int(ch.Aff_abils.Dex)
	perc += rand_number(-5, 5)
	prob += rand_number(-5, 5)
	if axion_dice(0) > 100 && vict.Position != POS_SLEEPING {
		reveal_hiding(ch, 0)
		act(libc.CString("@r$N@R just happens to glance in your direction! What terrible luck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@RYou just happen to glance behind you and spot @r$n@R trying to STEAL from you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@r$N@R just happens to glance in @r$n's@R direction and catches $m trying to STEAL!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		prob = -1000
	}
	if prob+20 < perc && vict.Position != POS_SLEEPING {
		reveal_hiding(ch, 0)
		act(libc.CString("@rYou are caught trying to stick your hand in @R$N's@r possessions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@rYou catch @R$n@r trying to rummage through your possessions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@R$n@R is caught by @R$N@r as $e sticks $s hand in @R$N's@r possessions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		if IS_NPC(vict) {
			set_fighting(vict, ch)
		}
		improve_skill(vict, SKILL_SPOT, 1)
	} else {
		if C.strcasecmp(&arg[0], libc.CString("zenni")) == 0 {
			if prob > perc {
				if vict.Gold > 0 {
					if vict.Gold > 100 {
						gold = (vict.Gold / 100) * rand_number(1, 10)
					} else {
						gold = vict.Gold
					}
					if gold+ch.Gold > GOLD_CARRY(ch) {
						send_to_char(ch, libc.CString("You can't hold that much more zenni on your person!\r\n"))
						return
					}
					vict.Gold -= gold
					ch.Gold += gold
					if !IS_NPC(vict) {
						vict.Act[int(PLR_STOLEN/32)] |= bitvector_t(1 << (int(PLR_STOLEN % 32)))
						mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("THEFT: %s has stolen %s zenni@n from %s"), GET_NAME(ch), add_commas(int64(gold)), GET_NAME(vict))
					}
					if gold > 1 {
						send_to_char(ch, libc.CString("Bingo!  You got %d zenni.\r\n"), gold)
					} else {
						send_to_char(ch, libc.CString("You manage to swipe a solitary zenni.\r\n"))
					}
					if axion_dice(0) > prob {
						send_to_char(ch, libc.CString("You think that your movements might have been a bit obvious.\r\n"))
						reveal_hiding(ch, 0)
						act(libc.CString("@R$n@r just stole zenni from @R$N@r!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						send_to_char(vict, libc.CString("You feel like something may be missing...\r\n"))
						if IS_NPC(vict) && rand_number(1, 3) == 3 {
							set_fighting(vict, ch)
						}
						improve_skill(vict, SKILL_SPOT, 1)
					}
					improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 1)
					return
				} else {
					send_to_char(ch, libc.CString("It appears like they are broke...\r\n"))
					return
				}
			} else {
				reveal_hiding(ch, 0)
				act(libc.CString("@rYou are caught trying to steal zenni from @R$N@r!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou catch @R$n's@r hand trying to snatch your zenni!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r catches @R$n's@r hand trying to snatch $S zenni!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				if IS_NPC(vict) {
					set_fighting(vict, ch)
				}
				improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 2)
				improve_skill(vict, SKILL_SPOT, 1)
				return
			}
		} else {
			if (func() *obj_data {
				obj = get_obj_in_list_vis(ch, &arg[0], nil, vict.Carrying)
				return obj
			}()) == nil {
				for eq_pos = 0; eq_pos < NUM_WEARS; eq_pos++ {
					if (vict.Equipment[eq_pos]) != nil && isname(&arg[0], (vict.Equipment[eq_pos]).Name) != 0 && CAN_SEE_OBJ(ch, vict.Equipment[eq_pos]) {
						obj = vict.Equipment[eq_pos]
						break
					}
				}
				if obj == nil {
					act(libc.CString("$E isn't wearing that item."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					return
				} else if vict.Position > POS_SLEEPING {
					send_to_char(ch, libc.CString("Steal worn equipment from them while they are awake? That's a stupid idea...\r\n"))
					return
				} else if give_otrigger(obj, vict, ch) == 0 || receive_mtrigger(ch, vict, obj) == 0 {
					send_to_char(ch, libc.CString("Impossible!\r\n"))
					return
				} else if GET_OBJ_VNUM(obj) >= 20000 {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if GET_OBJ_VNUM(obj) >= 18800 && GET_OBJ_VNUM(obj) <= 0x4A37 {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if GET_OBJ_VNUM(obj) >= 19100 && GET_OBJ_VNUM(obj) <= 0x4AFF {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if obj.Type_flag == ITEM_KEY {
					send_to_char(ch, libc.CString("No stealing keys!\r\n"))
					return
				} else if OBJ_FLAGGED(obj, ITEM_NOSTEAL) {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if obj.Weight+int64(gear_weight(ch)) > max_carry_weight(ch) {
					send_to_char(ch, libc.CString("You can't carry that much weight.\r\n"))
					return
				} else if ch.Carry_items+1 > 50 {
					send_to_char(ch, libc.CString("You don't have the room for it right now!\r\n"))
					return
				} else if prob > perc {
					act(libc.CString("You unequip $p and steal it."), FALSE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					if axion_dice(0) > prob {
						send_to_char(ch, libc.CString("You think that your movements might have been a bit obvious.\r\n"))
						reveal_hiding(ch, 0)
						act(libc.CString("@R$n@r just stole $p@r from @R$N@r!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_ROOM)
						send_to_char(vict, libc.CString("You feel your body being disturbed.\r\n"))
						improve_skill(vict, SKILL_SPOT, 1)
					}
					obj_to_char(unequip_char(vict, eq_pos), ch)
					improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 1)
					return
				} else {
					reveal_hiding(ch, 0)
					vict.Position = POS_SITTING
					act(libc.CString("@rYou are caught trying to steal $p@r from @R$N@r!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou feel your body being shifted while you sleep and wake up to find @R$n@r trying to steal $p@r from you!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r catches @R$n's@r trying to $p@r from $M during $S sleep!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					if IS_NPC(vict) {
						vict.Position = POS_STANDING
						set_fighting(vict, ch)
					}
					improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 2)
					improve_skill(vict, SKILL_SPOT, 1)
					return
				}
			} else {
				if give_otrigger(obj, vict, ch) == 0 || receive_mtrigger(ch, vict, obj) == 0 {
					send_to_char(ch, libc.CString("Impossible!\r\n"))
					return
				} else if GET_OBJ_VNUM(obj) >= 20000 {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if OBJ_FLAGGED(obj, ITEM_NOSTEAL) {
					send_to_char(ch, libc.CString("You can't steal that!\r\n"))
					return
				} else if obj.Type_flag == ITEM_KEY {
					send_to_char(ch, libc.CString("No stealing keys!\r\n"))
					return
				} else if obj.Weight+int64(gear_weight(ch)) > max_carry_weight(ch) {
					send_to_char(ch, libc.CString("You can't carry that much weight.\r\n"))
					return
				} else if ch.Carry_items+1 > 50 {
					send_to_char(ch, libc.CString("You don't have the room for it right now!\r\n"))
					return
				} else if prob > perc {
					act(libc.CString("You steal $p from $N."), FALSE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					obj_from_char(obj)
					obj_to_char(obj, ch)
					if !IS_NPC(vict) {
						vict.Act[int(PLR_STOLEN/32)] |= bitvector_t(1 << (int(PLR_STOLEN % 32)))
						mudlog(NRM, MAX(ADMLVL_GRGOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("THEFT: %s has stolen %s@n from %s"), GET_NAME(ch), obj.Short_description, GET_NAME(vict))
					}
					if axion_dice(0) > prob {
						reveal_hiding(ch, 0)
						send_to_char(ch, libc.CString("You think that your movements might have been a bit obvious.\r\n"))
						act(libc.CString("@R$n@r just stole $p@r from @R$N@r!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_ROOM)
						send_to_char(vict, libc.CString("You feel like something may be missing...\r\n"))
						if IS_NPC(vict) && rand_number(1, 3) == 3 {
							set_fighting(vict, ch)
						}
						improve_skill(vict, SKILL_SPOT, 1)
					}
					improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 1)
					return
				} else {
					reveal_hiding(ch, 0)
					act(libc.CString("@rYou are caught trying to steal $p@r from @R$N@r!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou catch @R$n@r trying to steal $p@r from you!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r catches @R$n's@r trying to $p@r!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					if IS_NPC(vict) {
						vict.Position = POS_STANDING
						set_fighting(vict, ch)
					}
					improve_skill(ch, SKILL_SLEIGHT_OF_HAND, 2)
					improve_skill(vict, SKILL_SPOT, 1)
					return
				}
			}
		}
	}
}
func do_practice(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [200]byte
	one_argument(argument, &arg[0])
	if arg[0] != 0 {
		send_to_char(ch, libc.CString("You can only practice skills with your trainer.\r\n"))
	} else {
		send_to_char(ch, libc.CString("Use the skills command unless you are at your trainer.\r\n"))
	}
}
func do_skills(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [1000]byte
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	list_skills(ch, &arg[0])
}
func do_visible(ch *char_data, argument *byte, cmd int, subcmd int) {
	var appeared int = 0
	if ch.Admlevel != 0 {
		perform_immort_vis(ch)
		return
	}
	if AFF_FLAGGED(ch, AFF_INVISIBLE) {
		appear(ch)
		appeared = 1
		send_to_char(ch, libc.CString("You break the spell of invisibility.\r\n"))
	}
	if AFF_FLAGGED(ch, AFF_ETHEREAL) && affectedv_by_spell(ch, ART_EMPTY_BODY) {
		affectv_from_char(ch, ART_EMPTY_BODY)
		if AFF_FLAGGED(ch, AFF_ETHEREAL) {
			send_to_char(ch, libc.CString("Returning to the material plane will not be so easy.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You return to the material plane.\r\n"))
			if appeared == 0 {
				act(libc.CString("$n flashes into existence."), FALSE, ch, nil, nil, TO_ROOM)
			}
		}
		appeared = 1
	}
	if appeared == 0 {
		send_to_char(ch, libc.CString("You are already visible.\r\n"))
	}
}
func do_title(ch *char_data, argument *byte, cmd int, subcmd int) {
	skip_spaces(&argument)
	delete_doubledollar(argument)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Your title is fine... go away.\r\n"))
	} else if PLR_FLAGGED(ch, PLR_NOTITLE) {
		send_to_char(ch, libc.CString("You can't title yourself -- you shouldn't have abused it!\r\n"))
	} else if C.strstr(argument, libc.CString("(")) != nil || C.strstr(argument, libc.CString(")")) != nil {
		send_to_char(ch, libc.CString("Titles can't contain the ( or ) characters.\r\n"))
	} else if C.strlen(argument) > MAX_TITLE_LENGTH {
		send_to_char(ch, libc.CString("Sorry, titles can't be longer than %d characters.\r\n"), MAX_TITLE_LENGTH)
	} else {
		set_title(ch, argument)
		send_to_char(ch, libc.CString("Okay, you're now %s %s.\r\n"), GET_NAME(ch), GET_TITLE(ch))
	}
}
func perform_group(ch *char_data, vict *char_data, highlvl int, lowlvl int, highpl int64, lowpl int64) int {
	if AFF_FLAGGED(vict, AFF_GROUP) || !CAN_SEE(ch, vict) {
		return 0
	}
	if (vict.Bonuses[BONUS_LONER]) > 0 {
		act(libc.CString("$n is the loner type and refuses to be in your group."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
		return 0
	}
	if GET_LEVEL(vict)+12 < highlvl {
		act(libc.CString("$n isn't experienced enough to be in your group with its current members."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
		return 0
	}
	if GET_LEVEL(vict) > lowlvl+12 {
		act(libc.CString("$n is too experienced to be in your group with its current members."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
		return 0
	}
	if highlvl >= 100 {
		if float64(gear_pl(vict)) > float64(highpl)*1.5 {
			act(libc.CString("$n is too powerful right now to be in a level 100 group with you."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
			return 0
		}
		if float64(gear_pl(vict)) < float64(lowpl)*0.5 {
			act(libc.CString("$n is too weak right now to be in a level 100 group with you."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
			return 0
		}
	}
	vict.Affected_by[int(AFF_GROUP/32)] |= 1 << (int(AFF_GROUP % 32))
	if ch != vict {
		act(libc.CString("$N is now a member of your group."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
	}
	act(libc.CString("You are now a member of $n's group."), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
	act(libc.CString("$N is now a member of $n's group."), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	return 1
}
func print_group(ch *char_data) {
	var (
		k *char_data
		f *follow_type
	)
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("But you are not the member of a group!\r\n"))
	} else {
		var buf [64936]byte
		send_to_char(ch, libc.CString("Your group consists of:\r\n"))
		if ch.Master != nil {
			k = ch.Master
		} else {
			k = ch
		}
		if AFF_FLAGGED(k, AFF_GROUP) {
			send_to_char(ch, libc.CString("@D----------------@n\r\n"))
			if k.Hit > k.Max_hit/10 {
				stdio.Snprintf(&buf[0], int(64936), "@gL@D: @w$N @W- @D[@RPL@Y: @c%s @CKi@Y: @c%s @GST@Y: @c%s@D] [@w%2d %s %s@D]@n", add_commas(k.Hit), add_commas(k.Mana), add_commas(k.Move), GET_LEVEL(k), class_abbrevs[int(k.Chclass)], race_abbrevs[int(k.Race)])
			}
			if k.Hit <= (k.Max_hit-int64(gear_weight(k)))/10 {
				stdio.Snprintf(&buf[0], int(64936), "@gL@D: @w$N @W- @D[@RPL@Y: @r%s @CKi@Y: @c%s @GST@Y: @c%s@D] [@w%2d %s %s@D]@n", add_commas(k.Hit), add_commas(k.Mana), add_commas(k.Move), GET_LEVEL(k), class_abbrevs[int(k.Chclass)], race_abbrevs[int(k.Race)])
			}
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(k), TO_CHAR)
		}
		for f = k.Followers; f != nil; f = f.Next {
			if !AFF_FLAGGED(f.Follower, AFF_GROUP) {
				continue
			}
			send_to_char(ch, libc.CString("@D----------------@n\r\n"))
			if f.Follower.Hit > (f.Follower.Max_hit-int64(gear_weight(f.Follower)))/10 {
				stdio.Snprintf(&buf[0], int(64936), "@gF@D: @w$N @W- @D[@RPL@Y: @c%s @CKi@Y: @c%s @GST@Y: @c%s@D] [@w%2d %s %s@D]", add_commas(f.Follower.Hit), add_commas(f.Follower.Mana), add_commas(f.Follower.Move), GET_LEVEL(f.Follower), class_abbrevs[int(f.Follower.Chclass)], race_abbrevs[int(f.Follower.Race)])
			}
			if f.Follower.Hit <= (f.Follower.Max_hit-int64(gear_weight(f.Follower)))/10 {
				stdio.Snprintf(&buf[0], int(64936), "@gF@D: @w$N @W- @D[@RPL@Y: @r%s @CKi@Y: @c%s @GST@Y: @c%s@D] [@w%2d %s %s@D]", add_commas(f.Follower.Hit), add_commas(f.Follower.Mana), add_commas(f.Follower.Move), GET_LEVEL(f.Follower), class_abbrevs[int(f.Follower.Chclass)], race_abbrevs[int(f.Follower.Race)])
			}
			act(&buf[0], FALSE, ch, nil, unsafe.Pointer(f.Follower), TO_CHAR)
		}
		send_to_char(ch, libc.CString("@D----------------@n\r\n"))
	}
}
func do_group(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf     [64936]byte
		vict    *char_data
		f       *follow_type
		found   int
		highlvl int   = 0
		lowlvl  int   = 0
		highpl  int64 = 0
		lowpl   int64 = 0
	)
	one_argument(argument, &buf[0])
	if (ch.Bonuses[BONUS_LONER]) > 0 {
		send_to_char(ch, libc.CString("You can not group as you prefer to be alone.\r\n"))
		return
	}
	if buf[0] == 0 {
		print_group(ch)
		return
	}
	if ch.Master != nil {
		act(libc.CString("You cannot enroll group members without being head of a group."), FALSE, ch, nil, nil, TO_CHAR)
		return
	}
	highlvl = GET_LEVEL(ch)
	lowlvl = GET_LEVEL(ch)
	highpl = gear_pl(ch)
	lowpl = gear_pl(ch)
	for func() *follow_type {
		found = 0
		return func() *follow_type {
			f = ch.Followers
			return f
		}()
	}(); f != nil; f = f.Next {
		if AFF_FLAGGED(f.Follower, AFF_GROUP) {
			if GET_LEVEL(f.Follower) > highlvl {
				highlvl = GET_LEVEL(f.Follower)
			}
			if GET_LEVEL(f.Follower) < lowlvl {
				lowlvl = GET_LEVEL(f.Follower)
			}
		}
	}
	var foundwas int = 0
	if C.strcasecmp(&buf[0], libc.CString("all")) == 0 {
		perform_group(ch, ch, GET_LEVEL(ch), GET_LEVEL(ch), highpl, lowpl)
		for func() *follow_type {
			found = 0
			return func() *follow_type {
				f = ch.Followers
				return f
			}()
		}(); f != nil; f = f.Next {
			foundwas = found
			found += perform_group(ch, f.Follower, highlvl, lowlvl, highpl, lowpl)
			if found > foundwas {
				if GET_LEVEL(f.Follower) > highlvl {
					highlvl = GET_LEVEL(f.Follower)
				} else if GET_LEVEL(f.Follower) < lowlvl {
					lowlvl = GET_LEVEL(f.Follower)
				}
			}
		}
		if found == 0 {
			send_to_char(ch, libc.CString("Everyone following you is already in your group.\r\n"))
		}
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &buf[0], nil, 1<<0)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if vict.Master != ch && vict != ch {
		act(libc.CString("$N must follow you to enter your group."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
	} else {
		if !AFF_FLAGGED(vict, AFF_GROUP) {
			if !AFF_FLAGGED(ch, AFF_GROUP) {
				send_to_char(ch, libc.CString("You form a group, with you as leader.\r\n"))
				ch.Affected_by[int(AFF_GROUP/32)] |= 1 << (int(AFF_GROUP % 32))
			}
			perform_group(ch, vict, highlvl, lowlvl, highpl, lowpl)
		} else {
			if ch != vict {
				act(libc.CString("$N is no longer a member of your group."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			}
			act(libc.CString("You have been kicked out of $n's group!"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$N has been kicked out of $n's group!"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			vict.Affected_by[int(AFF_GROUP/32)] &= ^(1 << (int(AFF_GROUP % 32)))
		}
	}
}
func do_ungroup(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf      [2048]byte
		f        *follow_type
		next_fol *follow_type
		tch      *char_data
	)
	one_argument(argument, &buf[0])
	if buf[0] == 0 {
		if ch.Master != nil || !AFF_FLAGGED(ch, AFF_GROUP) {
			send_to_char(ch, libc.CString("But you lead no group!\r\n"))
			return
		}
		for f = ch.Followers; f != nil; f = next_fol {
			next_fol = f.Next
			if AFF_FLAGGED(f.Follower, AFF_GROUP) {
				f.Follower.Affected_by[int(AFF_GROUP/32)] &= ^(1 << (int(AFF_GROUP % 32)))
				act(libc.CString("$N has disbanded the group."), TRUE, f.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
				f.Follower.Combatexpertise = 0
				if !AFF_FLAGGED(f.Follower, AFF_CHARM) {
					stop_follower(f.Follower)
				}
			}
		}
		ch.Affected_by[int(AFF_GROUP/32)] &= ^(1 << (int(AFF_GROUP % 32)))
		ch.Combatexpertise = 0
		send_to_char(ch, libc.CString("You disband the group.\r\n"))
		return
	}
	if (func() *char_data {
		tch = get_char_vis(ch, &buf[0], nil, 1<<0)
		return tch
	}()) == nil {
		send_to_char(ch, libc.CString("There is no such person!\r\n"))
		return
	}
	if tch.Master != ch {
		send_to_char(ch, libc.CString("That person is not following you!\r\n"))
		return
	}
	if !AFF_FLAGGED(tch, AFF_GROUP) {
		send_to_char(ch, libc.CString("That person isn't in your group.\r\n"))
		return
	}
	tch.Affected_by[int(AFF_GROUP/32)] &= ^(1 << (int(AFF_GROUP % 32)))
	tch.Combatexpertise = 0
	act(libc.CString("$N is no longer a member of your group."), FALSE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
	act(libc.CString("You have been kicked out of $n's group!"), FALSE, ch, nil, unsafe.Pointer(tch), TO_VICT)
	act(libc.CString("$N has been kicked out of $n's group!"), FALSE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
	if !AFF_FLAGGED(tch, AFF_CHARM) {
		stop_follower(tch)
	}
}
func do_report(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf [64936]byte
		k   *char_data
		f   *follow_type
	)
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		send_to_char(ch, libc.CString("But you are not a member of any group!\r\n"))
		return
	}
	stdio.Snprintf(&buf[0], int(64936), "$n reports: %lld/%lldH, %lld/%lldM, %lld/%lldV\r\n", ch.Hit, ch.Max_hit, ch.Mana, ch.Max_mana, ch.Move, ch.Max_move)
	if ch.Master != nil {
		k = ch.Master
	} else {
		k = ch
	}
	for f = k.Followers; f != nil; f = f.Next {
		if AFF_FLAGGED(f.Follower, AFF_GROUP) && f.Follower != ch {
			act(&buf[0], TRUE, ch, nil, unsafe.Pointer(f.Follower), TO_VICT)
		}
	}
	if k != ch {
		act(&buf[0], TRUE, ch, nil, unsafe.Pointer(k), TO_VICT)
	}
	send_to_char(ch, libc.CString("You report to the group.\r\n"))
}
func do_split(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		amount int
		num    int
		share  int
		rest   int
		len_   uint64
		k      *char_data
		f      *follow_type
	)
	if IS_NPC(ch) {
		return
	}
	one_argument(argument, &buf[0])
	if is_number(&buf[0]) != 0 {
		amount = libc.Atoi(libc.GoString(&buf[0]))
		if amount <= 0 {
			send_to_char(ch, libc.CString("Sorry, you can't do that.\r\n"))
			return
		}
		if amount > ch.Gold {
			send_to_char(ch, libc.CString("You don't seem to have that much gold to split.\r\n"))
			return
		}
		ch.Gold -= amount
		if ch.Master != nil {
			k = ch.Master
		} else {
			k = ch
		}
		if AFF_FLAGGED(k, AFF_GROUP) && k.In_room == ch.In_room {
			num = 1
		} else {
			num = 0
		}
		for f = k.Followers; f != nil; f = f.Next {
			if AFF_FLAGGED(f.Follower, AFF_GROUP) && !IS_NPC(f.Follower) && f.Follower != ch && f.Follower.In_room == ch.In_room {
				num++
			}
		}
		if num > 0 && AFF_FLAGGED(ch, AFF_GROUP) {
			share = amount / num
			rest = amount % num
		} else {
			send_to_char(ch, libc.CString("With whom do you wish to share your gold?\r\n"))
			return
		}
		ch.Gold += share
		len_ = uint64(stdio.Snprintf(&buf[0], int(2048), "%s splits %d zenni; you receive %d.\r\n", GET_NAME(ch), amount, share))
		if rest != 0 && len_ < uint64(2048) {
			stdio.Snprintf(&buf[len_], int(2048-uintptr(len_)), "%d zenni %s not splitable, so %s keeps the money.\r\n", rest, func() string {
				if rest == 1 {
					return "was"
				}
				return "were"
			}(), GET_NAME(ch))
		}
		if AFF_FLAGGED(k, AFF_GROUP) && k.In_room == ch.In_room && !IS_NPC(k) && k != ch {
			k.Gold += share
			send_to_char(k, libc.CString("%s"), &buf[0])
		}
		for f = k.Followers; f != nil; f = f.Next {
			if AFF_FLAGGED(f.Follower, AFF_GROUP) && !IS_NPC(f.Follower) && f.Follower.In_room == ch.In_room && f.Follower != ch {
				f.Follower.Gold += share
				send_to_char(f.Follower, libc.CString("%s"), &buf[0])
			}
		}
		send_to_char(ch, libc.CString("You split %d zenni among %d members -- %d zenni each.\r\n"), amount, num, share)
		if rest != 0 {
			send_to_char(ch, libc.CString("%d zenni %s not splitable, so you keep the money.\r\n"), rest, func() string {
				if rest == 1 {
					return "was"
				}
				return "were"
			}())
			ch.Gold += rest
		}
	} else {
		send_to_char(ch, libc.CString("How much zenni do you wish to split with your group?\r\n"))
		return
	}
}
func do_use(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf      [100]byte
		arg      [2048]byte
		mag_item *obj_data = nil
	)
	half_chop(argument, &arg[0], &buf[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("What do you want to %s?\r\n"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command)
		return
	}
	if mag_item == nil {
		switch subcmd {
		case SCMD_RECITE:
			fallthrough
		case SCMD_QUAFF:
			if (func() *obj_data {
				mag_item = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
				return mag_item
			}()) == nil {
				send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
				return
			}
		case SCMD_USE:
			if (func() *obj_data {
				mag_item = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
				return mag_item
			}()) == nil {
				send_to_char(ch, libc.CString("You don't seem to have %s %s.\r\n"), AN(&arg[0]), &arg[0])
				return
			}
		default:
			basic_mud_log(libc.CString("SYSERR: Unknown subcmd %d passed to do_use."), subcmd)
			return
		}
	}
	switch subcmd {
	case SCMD_QUAFF:
		if mag_item.Type_flag != ITEM_POTION {
			send_to_char(ch, libc.CString("You can only swallow beans.\r\n"))
			return
		}
		if ch.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't swallow beans, you are an android.\r\n"))
			return
		}
		if OBJ_FLAGGED(mag_item, ITEM_FORGED) {
			send_to_char(ch, libc.CString("You can't swallow that, it is fake!\r\n"))
			return
		}
		if OBJ_FLAGGED(mag_item, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("You can't swallow that, it is broken!\r\n"))
			return
		}
	case SCMD_RECITE:
		if mag_item.Type_flag != ITEM_SCROLL {
			send_to_char(ch, libc.CString("You can only recite scrolls.\r\n"))
			return
		}
	case SCMD_USE:
		if ch.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("You are not biological enough to use these, Tincan.\r\n"))
			return
		} else {
			switch GET_OBJ_VNUM(mag_item) {
			case 381:
				if ch.Move >= ch.Max_move {
					send_to_char(ch, libc.CString("Your stamina is full.\r\n"))
					return
				}
				act(libc.CString("@WYou place the $p@W against your chest and feel a rush of stamina as it automatically administers the dose.@n"), TRUE, ch, mag_item, nil, TO_CHAR)
				act(libc.CString("@C$n@W places an $p@W against $s chest and a loud click is heard.@n"), TRUE, ch, mag_item, nil, TO_ROOM)
				if GET_SKILL(ch, SKILL_FIRST_AID) > 0 {
					send_to_char(ch, libc.CString("@CYour skill in First Aid has helped increase the use of the injector. You gain more stamina as a result.@n\r\n"))
					ch.Move += int64(float64(ch.Max_move) * 0.25)
					if ch.Move > ch.Max_move {
						ch.Move = ch.Max_move
					}
				} else {
					ch.Move += int64(float64(ch.Max_move) * 0.1)
					if ch.Move > ch.Max_move {
						ch.Move = ch.Max_move
					}
				}
				extract_obj(mag_item)
				return
			case 382:
				if AFF_FLAGGED(ch, AFF_BURNED) {
					act(libc.CString("@WYou gently apply the salve to your burns.@n"), TRUE, ch, mag_item, nil, TO_CHAR)
					act(libc.CString("@C$n@W gently applies a burn salve to $s burns.@n"), TRUE, ch, mag_item, nil, TO_ROOM)
					ch.Affected_by[int(AFF_BURNED/32)] &= ^(1 << (int(AFF_BURNED % 32)))
					extract_obj(mag_item)
				} else {
					send_to_char(ch, libc.CString("You are not burned.\r\n"))
				}
				return
			case 383:
				if AFF_FLAGGED(ch, AFF_POISON) {
					act(libc.CString("@WYou place the $p@W against your neck and feel a rush of relief as the antitoxiin enters your bloodstream.@n"), TRUE, ch, mag_item, nil, TO_CHAR)
					act(libc.CString("@C$n@W places an $p@W against $s neck and a loud click is heard.@n"), TRUE, ch, mag_item, nil, TO_ROOM)
					null_affect(ch, AFF_POISON)
					extract_obj(mag_item)
				} else {
					send_to_char(ch, libc.CString("You are not poisoned.\r\n"))
				}
				return
			case 385:
				act(libc.CString("@WYou drink the contents of the vial before disposing of it.@n"), TRUE, ch, mag_item, nil, TO_CHAR)
				act(libc.CString("@C$n@W dinks a $p and then disposes of it.@n"), TRUE, ch, mag_item, nil, TO_ROOM)
				if AFF_FLAGGED(ch, AFF_BLIND) {
					act(libc.CString("@WYour eyesight has returned!@n"), TRUE, ch, mag_item, nil, TO_CHAR)
					act(libc.CString("@C$n@W eyesight seems to have returned.@n"), TRUE, ch, mag_item, nil, TO_ROOM)
					null_affect(ch, AFF_BLIND)
				}
				var refreshed int = FALSE
				if float64(ch.Hit) <= float64(gear_pl(ch))*0.99 {
					ch.Hit += large_rand(int64(float64(gear_pl(ch))*0.08), int64(float64(gear_pl(ch))*0.16))
					if ch.Hit > gear_pl(ch) {
						ch.Hit = gear_pl(ch)
					}
					refreshed = TRUE
				} else if float64(ch.Mana) <= float64(gear_pl(ch))*0.99 {
					ch.Mana += large_rand(int64(float64(ch.Max_mana)*0.08), int64(float64(ch.Max_mana)*0.16))
					if ch.Mana > ch.Max_mana {
						ch.Mana = ch.Max_mana
					}
					refreshed = TRUE
				} else if float64(ch.Move) <= float64(ch.Max_move)*0.99 {
					ch.Move += large_rand(int64(float64(ch.Max_move)*0.08), int64(float64(ch.Max_move)*0.16))
					if ch.Move > ch.Max_move {
						ch.Move = ch.Max_move
					}
					refreshed = TRUE
				}
				if refreshed == TRUE {
					send_to_char(ch, libc.CString("@CYou feel refreshed!\r\n"))
				}
				extract_obj(mag_item)
				return
			default:
				send_to_char(ch, libc.CString("That is not something you can apparently use.\r\n"))
				return
			}
		}
	}
	mag_objectmagic(ch, mag_item, &buf[0])
}
func do_value(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg       [2048]byte
		value_lev int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		switch subcmd {
		case SCMD_WIMPY:
			if ch.Player_specials.Wimp_level != 0 {
				send_to_char(ch, libc.CString("Your current wimp level is %d powerlevel.\r\n"), ch.Player_specials.Wimp_level)
				return
			} else {
				send_to_char(ch, libc.CString("At the moment, you're not a wimp.  (sure, sure...)\r\n"))
				return
			}
		}
	}
	if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(arg[0]))))) & int(uint16(int16(_ISdigit)))) != 0 {
		switch subcmd {
		case SCMD_WIMPY:
			if IS_NPC(ch) {
				return
			}
			if (func() int {
				value_lev = libc.Atoi(libc.GoString(&arg[0]))
				return value_lev
			}()) != 0 {
				if value_lev < 0 {
					send_to_char(ch, libc.CString("Heh, heh, heh.. we are jolly funny today, eh?\r\n"))
				} else if value_lev > int(ch.Max_hit) {
					send_to_char(ch, libc.CString("That doesn't make much sense, now does it?\r\n"))
				} else if float64(value_lev) > (float64(ch.Max_hit) * 0.5) {
					send_to_char(ch, libc.CString("You can't set your wimp level above half your powerlevel.\r\n"))
				} else {
					send_to_char(ch, libc.CString("Okay, you'll wimp out if you drop below %d powerlevel.\r\n"), value_lev)
					ch.Player_specials.Wimp_level = value_lev
				}
			} else {
				send_to_char(ch, libc.CString("Okay, you'll now tough out fights to the bitter end.\r\n"))
				ch.Player_specials.Wimp_level = 0
			}
		default:
			basic_mud_log(libc.CString("Unknown subcmd to do_value %d called by %s"), subcmd, GET_NAME(ch))
		}
	} else {
		send_to_char(ch, libc.CString("Specify a value.  (0 to disable)\r\n"))
	}
}
func do_display(ch *char_data, argument *byte, cmd int, subcmd int) {
	var i uint64
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Mosters don't need displays.  Go away.\r\n"))
		return
	}
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("Usage: prompt { P | K | T | S | F | H | G | L | C | M | all/on | none/off }\r\n"))
		return
	}
	if C.strcasecmp(argument, libc.CString("on")) == 0 || C.strcasecmp(argument, libc.CString("all")) == 0 {
		ch.Player_specials.Pref[int(PRF_DISPHP/32)] |= bitvector_t(1 << (int(PRF_DISPHP % 32)))
		ch.Player_specials.Pref[int(PRF_DISPMOVE/32)] |= bitvector_t(1 << (int(PRF_DISPMOVE % 32)))
		ch.Player_specials.Pref[int(PRF_DISPKI/32)] |= bitvector_t(1 << (int(PRF_DISPKI % 32)))
		ch.Player_specials.Pref[int(PRF_DISPTNL/32)] |= bitvector_t(1 << (int(PRF_DISPTNL % 32)))
		ch.Player_specials.Pref[int(PRF_FURY/32)] |= bitvector_t(1 << (int(PRF_FURY % 32)))
		ch.Player_specials.Pref[int(PRF_DISTIME/32)] |= bitvector_t(1 << (int(PRF_DISTIME % 32)))
		ch.Player_specials.Pref[int(PRF_DISGOLD/32)] |= bitvector_t(1 << (int(PRF_DISGOLD % 32)))
		ch.Player_specials.Pref[int(PRF_DISPRAC/32)] |= bitvector_t(1 << (int(PRF_DISPRAC % 32)))
		ch.Player_specials.Pref[int(PRF_DISHUTH/32)] |= bitvector_t(1 << (int(PRF_DISHUTH % 32)))
		ch.Player_specials.Pref[int(PRF_DISPERC/32)] |= bitvector_t(1 << (int(PRF_DISPERC % 32)))
	} else if C.strcasecmp(argument, libc.CString("off")) == 0 || C.strcasecmp(argument, libc.CString("none")) == 0 {
		ch.Player_specials.Pref[int(PRF_DISPHP/32)] &= bitvector_t(^(1 << (int(PRF_DISPHP % 32))))
		ch.Player_specials.Pref[int(PRF_DISPKI/32)] &= bitvector_t(^(1 << (int(PRF_DISPKI % 32))))
		ch.Player_specials.Pref[int(PRF_DISPMOVE/32)] &= bitvector_t(^(1 << (int(PRF_DISPMOVE % 32))))
		ch.Player_specials.Pref[int(PRF_DISPTNL/32)] &= bitvector_t(^(1 << (int(PRF_DISPTNL % 32))))
		ch.Player_specials.Pref[int(PRF_FURY/32)] &= bitvector_t(^(1 << (int(PRF_FURY % 32))))
		ch.Player_specials.Pref[int(PRF_DISTIME/32)] &= bitvector_t(^(1 << (int(PRF_DISTIME % 32))))
		ch.Player_specials.Pref[int(PRF_DISGOLD/32)] &= bitvector_t(^(1 << (int(PRF_DISGOLD % 32))))
		ch.Player_specials.Pref[int(PRF_DISPRAC/32)] &= bitvector_t(^(1 << (int(PRF_DISPRAC % 32))))
		ch.Player_specials.Pref[int(PRF_DISHUTH/32)] &= bitvector_t(^(1 << (int(PRF_DISHUTH % 32))))
		ch.Player_specials.Pref[int(PRF_DISPERC/32)] &= bitvector_t(^(1 << (int(PRF_DISPERC % 32))))
	} else {
		for i = 0; i < uint64(C.strlen(argument)); i++ {
			switch C.tolower(int(*(*byte)(unsafe.Add(unsafe.Pointer(argument), i)))) {
			case 'p':
				ch.Player_specials.Pref[int(PRF_DISPHP/32)] = ch.Player_specials.Pref[int(PRF_DISPHP/32)] ^ bitvector_t(1<<(int(PRF_DISPHP%32)))
			case 's':
				ch.Player_specials.Pref[int(PRF_DISPMOVE/32)] = ch.Player_specials.Pref[int(PRF_DISPMOVE/32)] ^ bitvector_t(1<<(int(PRF_DISPMOVE%32)))
			case 'k':
				ch.Player_specials.Pref[int(PRF_DISPKI/32)] = ch.Player_specials.Pref[int(PRF_DISPKI/32)] ^ bitvector_t(1<<(int(PRF_DISPKI%32)))
			case 't':
				ch.Player_specials.Pref[int(PRF_DISPTNL/32)] = ch.Player_specials.Pref[int(PRF_DISPTNL/32)] ^ bitvector_t(1<<(int(PRF_DISPTNL%32)))
			case 'h':
				ch.Player_specials.Pref[int(PRF_DISTIME/32)] = ch.Player_specials.Pref[int(PRF_DISTIME/32)] ^ bitvector_t(1<<(int(PRF_DISTIME%32)))
			case 'g':
				ch.Player_specials.Pref[int(PRF_DISGOLD/32)] = ch.Player_specials.Pref[int(PRF_DISGOLD/32)] ^ bitvector_t(1<<(int(PRF_DISGOLD%32)))
			case 'l':
				ch.Player_specials.Pref[int(PRF_DISPRAC/32)] = ch.Player_specials.Pref[int(PRF_DISPRAC/32)] ^ bitvector_t(1<<(int(PRF_DISPRAC%32)))
			case 'c':
				ch.Player_specials.Pref[int(PRF_DISPERC/32)] = ch.Player_specials.Pref[int(PRF_DISPERC/32)] ^ bitvector_t(1<<(int(PRF_DISPERC%32)))
			case 'm':
				ch.Player_specials.Pref[int(PRF_DISHUTH/32)] = ch.Player_specials.Pref[int(PRF_DISHUTH/32)] ^ bitvector_t(1<<(int(PRF_DISHUTH%32)))
			case 'f':
				if ch.Race != RACE_HALFBREED {
					send_to_char(ch, libc.CString("Only halfbreeds use fury.\r\n"))
				}
				ch.Player_specials.Pref[int(PRF_FURY/32)] = ch.Player_specials.Pref[int(PRF_FURY/32)] ^ bitvector_t(1<<(int(PRF_FURY%32)))
			default:
				send_to_char(ch, libc.CString("Usage: prompt { P | K | T | S | F | H | G | L | all/on | none/off }\r\n"))
				return
			}
		}
	}
	send_to_char(ch, libc.CString("%s"), config_info.Play.OK)
}
func do_gen_write(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		fl       *C.FILE
		tmp      *byte
		filename *byte
		fbuf     stat
		ct       int64
	)
	switch subcmd {
	case SCMD_BUG:
		filename = libc.CString(LIB_MISC)
	case SCMD_TYPO:
		filename = libc.CString(LIB_MISC)
	case SCMD_IDEA:
		filename = libc.CString(LIB_MISC)
	default:
		return
	}
	ct = C.time(nil)
	tmp = C.asctime(C.localtime(&ct))
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Monsters can't have ideas - Go away.\r\n"))
		return
	}
	skip_spaces(&argument)
	delete_doubledollar(argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("That must be a mistake...\r\n"))
		return
	}
	send_to_imm(libc.CString("[A new %s has been filed by: %s]\r\n"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command, GET_NAME(ch))
	if C.stat(filename, &fbuf) < 0 {
		C.perror(libc.CString("SYSERR: Can't C.stat() file"))
		return
	}
	if fbuf.St_size >= __off_t(config_info.Operation.Max_filesize) {
		send_to_char(ch, libc.CString("Sorry, the file is full right now.. try again later.\r\n"))
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(filename), "a")))
		return fl
	}()) == nil {
		C.perror(libc.CString("SYSERR: do_gen_write"))
		send_to_char(ch, libc.CString("Could not open the file.  Sorry.\r\n"))
		return
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "@D[@WUser: @c%-10s@D] [@WChar: @C%-10s@D] [@WRoom: @G%-4d@D] [@WDate: @Y%6.6s@D]@b \n-----------@w\n%s\n", func() *byte {
		if GET_USER(ch) != nil {
			return GET_USER(ch)
		}
		return libc.CString("ERR")
	}(), GET_NAME(ch), func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}(), (*byte)(unsafe.Add(unsafe.Pointer(tmp), 4)), argument)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "@D-------------------------------@n\n")
	C.fclose(fl)
	send_to_char(ch, libc.CString("Okay.  Thanks!\r\n"))
}
func do_gen_tog(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		result       int
		tog_messages [44][2]*byte = [44][2]*byte{{libc.CString("You are now safe from summoning by other players.\r\n"), libc.CString("You may now be summoned by other players.\r\n")}, {libc.CString("Nohassle disabled.\r\n"), libc.CString("Nohassle enabled.\r\n")}, {libc.CString("Brief mode off.\r\n"), libc.CString("Brief mode on.\r\n")}, {libc.CString("Compact mode off.\r\n"), libc.CString("Compact mode on.\r\n")}, {libc.CString("You can now hear tells.\r\n"), libc.CString("You are now deaf to tells.\r\n")}, {libc.CString("You can now hear newbie.\r\n"), libc.CString("You are now deaf to newbie.\r\n")}, {libc.CString("You can now hear shouts.\r\n"), libc.CString("You are now deaf to shouts.\r\n")}, {libc.CString("You can now hear ooc.\r\n"), libc.CString("You are now deaf to ooc.\r\n")}, {libc.CString("You can now hear the congratulation messages.\r\n"), libc.CString("You are now deaf to the congratulation messages.\r\n")}, {libc.CString("You can now hear the Wiz-channel.\r\n"), libc.CString("You are now deaf to the Wiz-channel.\r\n")}, {libc.CString("You are no longer part of the Quest.\r\n"), libc.CString("Okay, you are part of the Quest!\r\n")}, {libc.CString("You will no longer see the room flags.\r\n"), libc.CString("You will now see the room flags.\r\n")}, {libc.CString("You will now have your communication repeated.\r\n"), libc.CString("You will no longer have your communication repeated.\r\n")}, {libc.CString("HolyLight mode off.\r\n"), libc.CString("HolyLight mode on.\r\n")}, {libc.CString("Nameserver_is_slow changed to NO; IP addresses will now be resolved.\r\n"), libc.CString("Nameserver_is_slow changed to YES; sitenames will no longer be resolved.\r\n")}, {libc.CString("Autoexits disabled.\r\n"), libc.CString("Autoexits enabled.\r\n")}, {libc.CString("Will no longer track through doors.\r\n"), libc.CString("Will now track through doors.\r\n")}, {libc.CString("Buildwalk Off.\r\n"), libc.CString("Buildwalk On.\r\n")}, {libc.CString("AFK flag is now off.\r\n"), libc.CString("AFK flag is now on.\r\n")}, {libc.CString("You will no longer Auto-Assist.\r\n"), libc.CString("You will now Auto-Assist.\r\n")}, {libc.CString("Autoloot disabled.\r\n"), libc.CString("Autoloot enabled.\r\n")}, {libc.CString("Autogold disabled.\r\n"), libc.CString("Autogold enabled.\r\n")}, {libc.CString("Will no longer clear screen in OLC.\r\n"), libc.CString("Will now clear screen in OLC.\r\n")}, {libc.CString("Autosplit disabled.\r\n"), libc.CString("Autosplit enabled.\r\n")}, {libc.CString("Autosac disabled.\r\n"), libc.CString("Autosac enabled.\r\n")}, {libc.CString("You will no longer attempt to be sneaky.\r\n"), libc.CString("You will try to move as silently as you can.\r\n")}, {libc.CString("You will no longer attempt to stay hidden.\r\n"), libc.CString("You will try to stay hidden.\r\n")}, {libc.CString("You will no longer automatically memorize spells in your list.\r\n"), libc.CString("You will automatically memorize spells in your list.\r\n")}, {libc.CString("Viewing newest board messages at top of list.\r\n"), libc.CString("Viewing newest board messages at bottom of list.\r\n")}, {libc.CString("Compression will be used if your client supports it.\r\n"), libc.CString("Compression will not be used even if your client supports it.\r\n")}, {libc.CString(""), libc.CString("")}, {libc.CString("You are no longer hidden from view on the who list and public channels.\r\n"), libc.CString("You are now hidden from view on the who list and public channels.\r\n")}, {libc.CString("You will now be told that you have mail on prompt.\r\n"), libc.CString("You will no longer be told that you have mail on prompt.\r\n")}, {libc.CString("You will no longer receive automatic hints.\r\n"), libc.CString("You will now receive automatic hints.\r\n")}, {libc.CString("Screen Reader Friendly Mode Deactivated.\r\n"), libc.CString("Screen Reader Friendly Mode Activated..\r\n")}, {libc.CString("You will now see equipment when looking at someone.\r\n"), libc.CString("You will no longer see equipment when looking at someone.\r\n")}, {libc.CString("You will now listen to the music channel.\r\n"), libc.CString("You will no longer listen to the music channel.\r\n")}, {libc.CString("You will now parry attacks.\r\n"), libc.CString("You will no longer parry attacks.\r\n")}, {libc.CString("You will no longer keep cybernetic limbs with death.\r\n"), libc.CString("You will now keep cybernetic limbs with death.\r\n")}, {libc.CString("You will no longer worry about acquiring steaks from animals.\r\n"), libc.CString("You will now acquire steaks from animal if you can.\r\n")}, {libc.CString("You will now accept things being given to you.\r\n"), libc.CString("You will no longer accept things being given to you.\r\n")}, {libc.CString("You will no longer instruct those you spar with.\r\n"), libc.CString("You will now instruct those you spar with.\r\n")}, {libc.CString("You will no longer view group health.\r\n"), libc.CString("You will now view group health.\r\n")}, {libc.CString("You will no longer view item health.\r\n"), libc.CString("You will now view item health.\r\n")}}
	)
	if IS_NPC(ch) {
		return
	}
	switch subcmd {
	case SCMD_NOSUMMON:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_SUMMONABLE/32)]
			ch.Player_specials.Pref[int(PRF_SUMMONABLE/32)] = ch.Player_specials.Pref[int(PRF_SUMMONABLE/32)] ^ bitvector_t(1<<(int(PRF_SUMMONABLE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_SUMMONABLE%32))))
	case SCMD_NOHASSLE:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOHASSLE/32)]
			ch.Player_specials.Pref[int(PRF_NOHASSLE/32)] = ch.Player_specials.Pref[int(PRF_NOHASSLE/32)] ^ bitvector_t(1<<(int(PRF_NOHASSLE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOHASSLE%32))))
	case SCMD_BRIEF:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_BRIEF/32)]
			ch.Player_specials.Pref[int(PRF_BRIEF/32)] = ch.Player_specials.Pref[int(PRF_BRIEF/32)] ^ bitvector_t(1<<(int(PRF_BRIEF%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_BRIEF%32))))
	case SCMD_COMPACT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_COMPACT/32)]
			ch.Player_specials.Pref[int(PRF_COMPACT/32)] = ch.Player_specials.Pref[int(PRF_COMPACT/32)] ^ bitvector_t(1<<(int(PRF_COMPACT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_COMPACT%32))))
	case SCMD_NOTELL:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOTELL/32)]
			ch.Player_specials.Pref[int(PRF_NOTELL/32)] = ch.Player_specials.Pref[int(PRF_NOTELL/32)] ^ bitvector_t(1<<(int(PRF_NOTELL%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOTELL%32))))
	case SCMD_NOAUCTION:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOAUCT/32)]
			ch.Player_specials.Pref[int(PRF_NOAUCT/32)] = ch.Player_specials.Pref[int(PRF_NOAUCT/32)] ^ bitvector_t(1<<(int(PRF_NOAUCT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOAUCT%32))))
	case SCMD_DEAF:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_DEAF/32)]
			ch.Player_specials.Pref[int(PRF_DEAF/32)] = ch.Player_specials.Pref[int(PRF_DEAF/32)] ^ bitvector_t(1<<(int(PRF_DEAF%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_DEAF%32))))
	case SCMD_NOGOSSIP:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOGOSS/32)]
			ch.Player_specials.Pref[int(PRF_NOGOSS/32)] = ch.Player_specials.Pref[int(PRF_NOGOSS/32)] ^ bitvector_t(1<<(int(PRF_NOGOSS%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOGOSS%32))))
	case SCMD_NOGRATZ:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOGRATZ/32)]
			ch.Player_specials.Pref[int(PRF_NOGRATZ/32)] = ch.Player_specials.Pref[int(PRF_NOGRATZ/32)] ^ bitvector_t(1<<(int(PRF_NOGRATZ%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOGRATZ%32))))
	case SCMD_NOWIZ:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOWIZ/32)]
			ch.Player_specials.Pref[int(PRF_NOWIZ/32)] = ch.Player_specials.Pref[int(PRF_NOWIZ/32)] ^ bitvector_t(1<<(int(PRF_NOWIZ%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOWIZ%32))))
	case SCMD_QUEST:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_QUEST/32)]
			ch.Player_specials.Pref[int(PRF_QUEST/32)] = ch.Player_specials.Pref[int(PRF_QUEST/32)] ^ bitvector_t(1<<(int(PRF_QUEST%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_QUEST%32))))
	case SCMD_ROOMFLAGS:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_ROOMFLAGS/32)]
			ch.Player_specials.Pref[int(PRF_ROOMFLAGS/32)] = ch.Player_specials.Pref[int(PRF_ROOMFLAGS/32)] ^ bitvector_t(1<<(int(PRF_ROOMFLAGS%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_ROOMFLAGS%32))))
	case SCMD_NOREPEAT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOREPEAT/32)]
			ch.Player_specials.Pref[int(PRF_NOREPEAT/32)] = ch.Player_specials.Pref[int(PRF_NOREPEAT/32)] ^ bitvector_t(1<<(int(PRF_NOREPEAT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOREPEAT%32))))
	case SCMD_HOLYLIGHT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_HOLYLIGHT/32)]
			ch.Player_specials.Pref[int(PRF_HOLYLIGHT/32)] = ch.Player_specials.Pref[int(PRF_HOLYLIGHT/32)] ^ bitvector_t(1<<(int(PRF_HOLYLIGHT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_HOLYLIGHT%32))))
	case SCMD_SLOWNS:
		result = func() int {
			p := &config_info.Operation.Nameserver_is_slow
			config_info.Operation.Nameserver_is_slow = int(libc.BoolToInt(config_info.Operation.Nameserver_is_slow == 0))
			return *p
		}()
	case SCMD_AUTOEXIT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)]
			ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)] = ch.Player_specials.Pref[int(PRF_AUTOEXIT/32)] ^ bitvector_t(1<<(int(PRF_AUTOEXIT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOEXIT%32))))
	case SCMD_TRACK:
		result = func() int {
			p := &config_info.Play.Track_through_doors
			config_info.Play.Track_through_doors = int(libc.BoolToInt(config_info.Play.Track_through_doors == 0))
			return *p
		}()
	case SCMD_AFK:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AFK/32)]
			ch.Player_specials.Pref[int(PRF_AFK/32)] = ch.Player_specials.Pref[int(PRF_AFK/32)] ^ bitvector_t(1<<(int(PRF_AFK%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AFK%32))))
		if PRF_FLAGGED(ch, PRF_AFK) {
			act(libc.CString("$n has gone AFK."), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			act(libc.CString("$n has come back from AFK."), TRUE, ch, nil, nil, TO_ROOM)
			if has_mail(int(ch.Idnum)) != 0 {
				send_to_char(ch, libc.CString("You have mail waiting.\r\n"))
			}
		}
	case SCMD_AUTOLOOT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOLOOT/32)]
			ch.Player_specials.Pref[int(PRF_AUTOLOOT/32)] = ch.Player_specials.Pref[int(PRF_AUTOLOOT/32)] ^ bitvector_t(1<<(int(PRF_AUTOLOOT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOLOOT%32))))
	case SCMD_AUTOGOLD:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOGOLD/32)]
			ch.Player_specials.Pref[int(PRF_AUTOGOLD/32)] = ch.Player_specials.Pref[int(PRF_AUTOGOLD/32)] ^ bitvector_t(1<<(int(PRF_AUTOGOLD%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOGOLD%32))))
	case SCMD_CLS:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_CLS/32)]
			ch.Player_specials.Pref[int(PRF_CLS/32)] = ch.Player_specials.Pref[int(PRF_CLS/32)] ^ bitvector_t(1<<(int(PRF_CLS%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_CLS%32))))
	case SCMD_BUILDWALK:
		if ch.Admlevel < ADMLVL_IMMORT {
			send_to_char(ch, libc.CString("Immortals only, sorry.\r\n"))
			return
		}
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_BUILDWALK/32)]
			ch.Player_specials.Pref[int(PRF_BUILDWALK/32)] = ch.Player_specials.Pref[int(PRF_BUILDWALK/32)] ^ bitvector_t(1<<(int(PRF_BUILDWALK%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_BUILDWALK%32))))
		if PRF_FLAGGED(ch, PRF_BUILDWALK) {
			mudlog(CMP, GET_LEVEL(ch), TRUE, libc.CString("OLC: %s turned buildwalk on. Allowed zone %d"), GET_NAME(ch), ch.Player_specials.Olc_zone)
		} else {
			mudlog(CMP, GET_LEVEL(ch), TRUE, libc.CString("OLC: %s turned buildwalk off. Allowed zone %d"), GET_NAME(ch), ch.Player_specials.Olc_zone)
		}
	case SCMD_AUTOSPLIT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOSPLIT/32)]
			ch.Player_specials.Pref[int(PRF_AUTOSPLIT/32)] = ch.Player_specials.Pref[int(PRF_AUTOSPLIT/32)] ^ bitvector_t(1<<(int(PRF_AUTOSPLIT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOSPLIT%32))))
	case SCMD_AUTOSAC:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOSAC/32)]
			ch.Player_specials.Pref[int(PRF_AUTOSAC/32)] = ch.Player_specials.Pref[int(PRF_AUTOSAC/32)] ^ bitvector_t(1<<(int(PRF_AUTOSAC%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOSAC%32))))
	case SCMD_SNEAK:
		result = (func() int {
			p := &ch.Affected_by[int(AFF_SNEAK/32)]
			ch.Affected_by[int(AFF_SNEAK/32)] = ch.Affected_by[int(AFF_SNEAK/32)] ^ 1<<(int(AFF_SNEAK%32))
			return *p
		}()) & (1 << (int(AFF_SNEAK % 32)))
	case SCMD_HIDE:
		if ch.Charge > 0 && ch.Preference != PREFERENCE_KI || float64(ch.Charge) > float64(ch.Max_mana)*0.1 && ch.Preference == PREFERENCE_KI || PLR_FLAGGED(ch, PLR_POWERUP) || AFF_FLAGGED(ch, AFF_FLYING) {
			send_to_char(ch, libc.CString("You stand out too much to hide right now!\r\n"))
			return
		} else if PLR_FLAGGED(ch, PLR_HEALT) {
			send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
			return
		}
		if GET_SKILL(ch, SKILL_HIDE) == 0 && slot_count(ch)+1 <= ch.Skill_slots {
			send_to_char(ch, libc.CString("@GYou learn the very minimal basics to hiding.@n\r\n"))
			for {
				ch.Skills[SKILL_HIDE] = int8(rand_number(1, 5))
				if true {
					break
				}
			}
		} else if GET_SKILL(ch, SKILL_HIDE) == 0 && slot_count(ch)+1 > ch.Skill_slots {
			send_to_char(ch, libc.CString("@RYou need more skill slots in order to learn this skill.@n\r\n"))
			return
		}
		result = (func() int {
			p := &ch.Affected_by[int(AFF_HIDE/32)]
			ch.Affected_by[int(AFF_HIDE/32)] = ch.Affected_by[int(AFF_HIDE/32)] ^ 1<<(int(AFF_HIDE%32))
			return *p
		}()) & (1 << (int(AFF_HIDE % 32)))
	case SCMD_AUTOMEM:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOMEM/32)]
			ch.Player_specials.Pref[int(PRF_AUTOMEM/32)] = ch.Player_specials.Pref[int(PRF_AUTOMEM/32)] ^ bitvector_t(1<<(int(PRF_AUTOMEM%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOMEM%32))))
	case SCMD_VIEWORDER:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_VIEWORDER/32)]
			ch.Player_specials.Pref[int(PRF_VIEWORDER/32)] = ch.Player_specials.Pref[int(PRF_VIEWORDER/32)] ^ bitvector_t(1<<(int(PRF_VIEWORDER%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_VIEWORDER%32))))
	case SCMD_TEST:
		if ch.Admlevel >= 1 {
			ch.Player_specials.Pref[int(PRF_TEST/32)] = ch.Player_specials.Pref[int(PRF_TEST/32)] ^ bitvector_t(1<<(int(PRF_TEST%32)))
			send_to_char(ch, libc.CString("Okay. Testing is now: %s\r\n"), func() string {
				if PRF_FLAGGED(ch, PRF_TEST) {
					return "On"
				}
				return "Off"
			}())
			if PRF_FLAGGED(ch, PRF_TEST) {
				send_to_char(ch, libc.CString("Make sure to remove nohassle as well.\r\n"))
			}
			return
		} else {
			send_to_char(ch, libc.CString("You are not an immortal.\r\n"))
			return
		}
	case SCMD_NOCOMPRESS:
		if config_info.Play.Enable_compression != 0 {
			result = int((func() bitvector_t {
				p := &ch.Player_specials.Pref[int(PRF_NOCOMPRESS/32)]
				ch.Player_specials.Pref[int(PRF_NOCOMPRESS/32)] = ch.Player_specials.Pref[int(PRF_NOCOMPRESS/32)] ^ bitvector_t(1<<(int(PRF_NOCOMPRESS%32)))
				return *p
			}()) & bitvector_t(1<<(int(PRF_NOCOMPRESS%32))))
			break
		} else {
			send_to_char(ch, libc.CString("Sorry, compression is globally disabled.\r\n"))
		}
		fallthrough
	case SCMD_AUTOASSIST:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_AUTOASSIST/32)]
			ch.Player_specials.Pref[int(PRF_AUTOASSIST/32)] = ch.Player_specials.Pref[int(PRF_AUTOASSIST/32)] ^ bitvector_t(1<<(int(PRF_AUTOASSIST%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_AUTOASSIST%32))))
	case SCMD_WHOHIDE:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_HIDE/32)]
			ch.Player_specials.Pref[int(PRF_HIDE/32)] = ch.Player_specials.Pref[int(PRF_HIDE/32)] ^ bitvector_t(1<<(int(PRF_HIDE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_HIDE%32))))
	case SCMD_NMWARN:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NMWARN/32)]
			ch.Player_specials.Pref[int(PRF_NMWARN/32)] = ch.Player_specials.Pref[int(PRF_NMWARN/32)] ^ bitvector_t(1<<(int(PRF_NMWARN%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NMWARN%32))))
	case SCMD_HINTS:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_HINTS/32)]
			ch.Player_specials.Pref[int(PRF_HINTS/32)] = ch.Player_specials.Pref[int(PRF_HINTS/32)] ^ bitvector_t(1<<(int(PRF_HINTS%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_HINTS%32))))
	case SCMD_NODEC:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NODEC/32)]
			ch.Player_specials.Pref[int(PRF_NODEC/32)] = ch.Player_specials.Pref[int(PRF_NODEC/32)] ^ bitvector_t(1<<(int(PRF_NODEC%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NODEC%32))))
	case SCMD_NOEQSEE:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOEQSEE/32)]
			ch.Player_specials.Pref[int(PRF_NOEQSEE/32)] = ch.Player_specials.Pref[int(PRF_NOEQSEE/32)] ^ bitvector_t(1<<(int(PRF_NOEQSEE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOEQSEE%32))))
	case SCMD_NOMUSIC:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOMUSIC/32)]
			ch.Player_specials.Pref[int(PRF_NOMUSIC/32)] = ch.Player_specials.Pref[int(PRF_NOMUSIC/32)] ^ bitvector_t(1<<(int(PRF_NOMUSIC%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOMUSIC%32))))
	case SCMD_NOPARRY:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOPARRY/32)]
			ch.Player_specials.Pref[int(PRF_NOPARRY/32)] = ch.Player_specials.Pref[int(PRF_NOPARRY/32)] ^ bitvector_t(1<<(int(PRF_NOPARRY%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOPARRY%32))))
	case SCMD_LKEEP:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_LKEEP/32)]
			ch.Player_specials.Pref[int(PRF_LKEEP/32)] = ch.Player_specials.Pref[int(PRF_LKEEP/32)] ^ bitvector_t(1<<(int(PRF_LKEEP%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_LKEEP%32))))
	case SCMD_CARVE:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_CARVE/32)]
			ch.Player_specials.Pref[int(PRF_CARVE/32)] = ch.Player_specials.Pref[int(PRF_CARVE/32)] ^ bitvector_t(1<<(int(PRF_CARVE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_CARVE%32))))
	case SCMD_NOGIVE:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_NOGIVE/32)]
			ch.Player_specials.Pref[int(PRF_NOGIVE/32)] = ch.Player_specials.Pref[int(PRF_NOGIVE/32)] ^ bitvector_t(1<<(int(PRF_NOGIVE%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_NOGIVE%32))))
	case SCMD_INSTRUCT:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_INSTRUCT/32)]
			ch.Player_specials.Pref[int(PRF_INSTRUCT/32)] = ch.Player_specials.Pref[int(PRF_INSTRUCT/32)] ^ bitvector_t(1<<(int(PRF_INSTRUCT%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_INSTRUCT%32))))
	case SCMD_GHEALTH:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_GHEALTH/32)]
			ch.Player_specials.Pref[int(PRF_GHEALTH/32)] = ch.Player_specials.Pref[int(PRF_GHEALTH/32)] ^ bitvector_t(1<<(int(PRF_GHEALTH%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_GHEALTH%32))))
	case SCMD_IHEALTH:
		result = int((func() bitvector_t {
			p := &ch.Player_specials.Pref[int(PRF_IHEALTH/32)]
			ch.Player_specials.Pref[int(PRF_IHEALTH/32)] = ch.Player_specials.Pref[int(PRF_IHEALTH/32)] ^ bitvector_t(1<<(int(PRF_IHEALTH%32)))
			return *p
		}()) & bitvector_t(1<<(int(PRF_IHEALTH%32))))
	default:
		basic_mud_log(libc.CString("SYSERR: Unknown subcmd %d in do_gen_toggle."), subcmd)
		return
	}
	if result != 0 {
		send_to_char(ch, libc.CString("%s"), tog_messages[subcmd][TOG_ON])
	} else {
		send_to_char(ch, libc.CString("%s"), tog_messages[subcmd][TOG_OFF])
	}
	return
}
func do_file(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		req_file  *C.FILE
		cur_line  int = 0
		num_lines int = 0
		req_lines int = 0
		i         int
		j         int
		l         int
		field     [100]byte
		value     [100]byte
		line      [256]byte
		buf       [64936]byte
	)
	type file_struct struct {
		Cmd   *byte
		Level int8
		File  *byte
	}
	var fields [17]file_struct = [17]file_struct{{Cmd: libc.CString("none"), Level: 6, File: libc.CString("Does Nothing")}, {Cmd: libc.CString("bug"), Level: ADMLVL_IMMORT, File: libc.CString("../lib/misc/bugs")}, {Cmd: libc.CString("typo"), Level: ADMLVL_IMMORT, File: libc.CString("../lib/misc/typos")}, {Cmd: libc.CString("report"), Level: ADMLVL_IMMORT, File: libc.CString("../lib/misc/ideas")}, {Cmd: libc.CString("xnames"), Level: 4, File: libc.CString("../lib/misc/xnames")}, {Cmd: libc.CString("levels"), Level: 4, File: libc.CString("../log/levels")}, {Cmd: libc.CString("rip"), Level: 4, File: libc.CString("../log/rip")}, {Cmd: libc.CString("players"), Level: 4, File: libc.CString("../log/newplayers")}, {Cmd: libc.CString("rentgone"), Level: 4, File: libc.CString("../log/rentgone")}, {Cmd: libc.CString("errors"), Level: 4, File: libc.CString("../log/errors")}, {Cmd: libc.CString("godcmds"), Level: 4, File: libc.CString("../log/godcmds")}, {Cmd: libc.CString("syslog"), Level: ADMLVL_IMMORT, File: libc.CString("../syslog")}, {Cmd: libc.CString("crash"), Level: ADMLVL_IMMORT, File: libc.CString("../syslog.CRASH")}, {Cmd: libc.CString("immlog"), Level: ADMLVL_IMMORT, File: libc.CString("../lib/misc/request")}, {Cmd: libc.CString("customs"), Level: ADMLVL_IMMORT, File: libc.CString("../lib/misc/customs")}, {Cmd: libc.CString("todo"), Level: 5, File: libc.CString("../todo")}, {Cmd: libc.CString("\n"), Level: 0, File: libc.CString("\n")}}
	skip_spaces(&argument)
	if *argument == 0 {
		C.strcpy(&buf[0], libc.CString("USAGE: file <option> <num lines>\r\n\r\nFile options:\r\n"))
		for func() int {
			j = 0
			return func() int {
				i = 1
				return i
			}()
		}(); int(fields[i].Level) != 0; i++ {
			if int(fields[i].Level) <= GET_LEVEL(ch) {
				stdio.Sprintf(&buf[C.strlen(&buf[0])], "%-15s%s\r\n", fields[i].Cmd, fields[i].File)
			}
		}
		send_to_char(ch, &buf[0])
		return
	}
	two_arguments(argument, &field[0], &value[0])
	for l = 0; *fields[l].Cmd != '\n'; l++ {
		if C.strncmp(&field[0], fields[l].Cmd, uint64(C.strlen(&field[0]))) == 0 {
			break
		}
	}
	if *fields[l].Cmd == '\n' {
		send_to_char(ch, libc.CString("That is not a valid option!\r\n"))
		return
	}
	if ch.Admlevel < int(fields[l].Level) {
		send_to_char(ch, libc.CString("You are not godly enough to view that file!\r\n"))
		return
	}
	if C.strcasecmp(&field[0], libc.CString("request")) == 0 {
		ch.Lboard[2] = C.time(nil)
	}
	if value[0] == 0 {
		req_lines = 15
	} else {
		req_lines = libc.Atoi(libc.GoString(&value[0]))
	}
	if (func() *C.FILE {
		req_file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(fields[l].File), "r")))
		return req_file
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: Error opening file %s using 'file' command."), fields[l].File)
		return
	}
	get_line(req_file, &line[0])
	for C.feof(req_file) == 0 {
		num_lines++
		get_line(req_file, &line[0])
	}
	C.rewind(req_file)
	req_lines = MIN(MIN(req_lines, num_lines), 5000)
	buf[0] = '\x00'
	get_line(req_file, &line[0])
	for C.feof(req_file) == 0 {
		cur_line++
		if cur_line > (num_lines - req_lines) {
			stdio.Sprintf(&buf[C.strlen(&buf[0])], "%s\r\n", &line[0])
		}
		get_line(req_file, &line[0])
	}
	C.fclose(req_file)
	page_string(ch.Desc, &buf[0], 1)
}
func do_compare(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1   [100]byte
		arg2   [2048]byte
		obj1   *obj_data
		obj2   *obj_data
		tchar  *char_data
		value1 int = 0
		value2 int = 0
		o1     int
		o2     int
		msg    *byte = nil
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 || arg2[0] == 0 {
		send_to_char(ch, libc.CString("Compare what to what?\n\r"))
		return
	}
	o1 = generic_find(&arg1[0], (1<<2)|1<<5, ch, &tchar, &obj1)
	o2 = generic_find(&arg2[0], (1<<2)|1<<5, ch, &tchar, &obj2)
	if o1 == 0 || o2 == 0 {
		send_to_char(ch, libc.CString("You do not have that item.\r\n"))
		return
	}
	if obj1 == obj2 {
		msg = libc.CString("You compare $p to itself.  It looks about the same.")
	} else if obj1.Type_flag != obj2.Type_flag {
		msg = libc.CString("You can't compare $p and $P.")
	} else {
		switch obj1.Type_flag {
		default:
			msg = libc.CString("You can't compare $p and $P.")
		case ITEM_ARMOR:
			value1 = obj1.Value[VAL_ARMOR_APPLYAC]
			value2 = obj2.Value[VAL_ARMOR_APPLYAC]
		case ITEM_WEAPON:
			value1 = ((obj1.Value[VAL_WEAPON_DAMSIZE]) + 1) * (obj1.Value[VAL_WEAPON_DAMDICE])
			value2 = ((obj2.Value[VAL_WEAPON_DAMSIZE]) + 1) * (obj2.Value[VAL_WEAPON_DAMDICE])
		}
	}
	if msg == nil {
		if value1 == value2 {
			msg = libc.CString("$p and $P look about the same.")
		} else if value1 > value2 {
			msg = libc.CString("$p looks better than $P.")
		} else {
			msg = libc.CString("$p looks worse than $P.")
		}
	}
	act(msg, FALSE, ch, obj1, unsafe.Pointer(obj2), TO_CHAR)
	return
}
func do_break(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		obj   *obj_data
		dummy *char_data = nil
		cmbrk int
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Usually you break SOMETHING.\r\n"))
		return
	}
	if (func() int {
		cmbrk = generic_find(&arg[0], (1<<2)|1<<5, ch, &dummy, &obj)
		return cmbrk
	}()) == 0 {
		send_to_char(ch, libc.CString("Can't seem to find what you want to break!\r\n"))
		return
	}
	if OBJ_FLAGGED(obj, ITEM_BROKEN) {
		send_to_char(ch, libc.CString("Seems like it's already broken!\r\n"))
		return
	}
	send_to_char(ch, libc.CString("You ruin %s.\r\n"), obj.Short_description)
	act(libc.CString("$n ruins $p."), FALSE, ch, obj, nil, TO_ROOM)
	obj.Value[VAL_ALL_HEALTH] = 0
	obj.Extra_flags[int(ITEM_BROKEN/32)] = obj.Extra_flags[int(ITEM_BROKEN/32)] ^ bitvector_t(1<<(int(ITEM_BROKEN%32)))
	return
}
func do_fix(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		obj      *obj_data
		obj4     *obj_data = nil
		rep      *obj_data
		next_obj *obj_data
		dummy    *char_data = nil
		cmbrk    int
		found    int = FALSE
		self     int = FALSE
		custom   int = FALSE
	)
	one_argument(argument, &arg[0])
	if know_skill(ch, SKILL_REPAIR) == 0 {
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Usually you fix SOMETHING.\r\n"))
		return
	}
	if C.strcasecmp(libc.CString("self"), &arg[0]) == 0 {
		if ch.Race != RACE_ANDROID {
			send_to_char(ch, libc.CString("Only androids can fix their bodies with repair kits.\r\n"))
			return
		} else {
			self = TRUE
		}
	}
	if self == FALSE {
		if (func() int {
			cmbrk = generic_find(&arg[0], (1<<2)|1<<5|1<<3, ch, &dummy, &obj)
			return cmbrk
		}()) == 0 {
			send_to_char(ch, libc.CString("Can't seem to find what you want to fix!\r\n"))
			return
		}
		if cmbrk != 0 && (obj.Value[VAL_ALL_HEALTH]) >= 100 {
			send_to_char(ch, libc.CString("But it isn't even damaged!\r\n"))
			return
		}
		if OBJ_FLAGGED(obj, ITEM_FORGED) {
			send_to_char(ch, libc.CString("That is fake, why bother fixing it?\r\n"))
			return
		}
		switch obj.Value[VAL_ALL_MATERIAL] {
		case MATERIAL_ORGANIC:
			fallthrough
		case MATERIAL_FOOD:
			fallthrough
		case MATERIAL_PAPER:
			fallthrough
		case MATERIAL_LIQUID:
			send_to_char(ch, libc.CString("You can't repair that.\r\n"))
			return
		}
		if GET_OBJ_VNUM(obj) == 0x4E83 || GET_OBJ_VNUM(obj) == 0x4E82 {
			custom = TRUE
		}
	}
	for rep = ch.Carrying; rep != nil; rep = next_obj {
		next_obj = rep.Next_content
		if custom == FALSE {
			if found == FALSE && GET_OBJ_VNUM(rep) == 48 && (!OBJ_FLAGGED(rep, ITEM_BROKEN) && !OBJ_FLAGGED(rep, ITEM_FORGED)) {
				found = TRUE
				obj4 = rep
			}
		} else {
			if found == FALSE && GET_OBJ_VNUM(rep) == 0x3519 && (!OBJ_FLAGGED(rep, ITEM_BROKEN) && !OBJ_FLAGGED(rep, ITEM_FORGED)) {
				found = TRUE
				obj4 = rep
			}
		}
	}
	if found == FALSE && custom == FALSE {
		send_to_char(ch, libc.CString("You do not even have a repair kit.\r\n"))
		return
	} else if found == FALSE && custom == TRUE {
		send_to_char(ch, libc.CString("You do not even have a Nano-tech Repair Orb.\r\n"))
		return
	}
	if self == FALSE {
		if GET_SKILL(ch, SKILL_REPAIR) < axion_dice(0) {
			act(libc.CString("You try to repair $p but screw up.."), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$n tries to repair $p but screws up.."), TRUE, ch, obj, nil, TO_ROOM)
			extract_obj(obj4)
			improve_skill(ch, SKILL_REPAIR, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		}
		if (obj.Value[VAL_ALL_HEALTH])+GET_SKILL(ch, SKILL_REPAIR) < 100 {
			send_to_char(ch, libc.CString("You repair %s a bit.\r\n"), obj.Short_description)
			act(libc.CString("$n repairs $p a bit."), FALSE, ch, obj, nil, TO_ROOM)
			obj.Value[VAL_ALL_HEALTH] += GET_SKILL(ch, SKILL_REPAIR)
			obj.Extra_flags[int(ITEM_BROKEN/32)] &= bitvector_t(^(1 << (int(ITEM_BROKEN % 32))))
		} else {
			send_to_char(ch, libc.CString("You repair %s completely.\r\n"), obj.Short_description)
			act(libc.CString("$n repairs $p completely."), FALSE, ch, obj, nil, TO_ROOM)
			obj.Value[VAL_ALL_HEALTH] = 100
			obj.Extra_flags[int(ITEM_BROKEN/32)] &= bitvector_t(^(1 << (int(ITEM_BROKEN % 32))))
		}
		if obj.Carried_by == nil && !PLR_FLAGGED(ch, PLR_REPLEARN) && (level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 || GET_LEVEL(ch) >= 100) {
			var gain int64 = int64((float64(level_exp(ch, GET_LEVEL(ch)+1)) * 0.0003) * float64(GET_SKILL(ch, SKILL_REPAIR)))
			send_to_char(ch, libc.CString("@mYou've learned a bit from repairing it. @D[@gEXP@W: @G+%s@D]@n\r\n"), add_commas(gain))
			ch.Act[int(PLR_REPLEARN/32)] |= bitvector_t(1 << (int(PLR_REPLEARN % 32)))
			gain_exp(ch, gain)
		} else if rand_number(2, 12) >= 10 && PLR_FLAGGED(ch, PLR_REPLEARN) {
			ch.Act[int(PLR_REPLEARN/32)] &= bitvector_t(^(1 << (int(PLR_REPLEARN % 32))))
			send_to_char(ch, libc.CString("@mYou think you might be on to something...@n\r\n"))
		}
		improve_skill(ch, SKILL_REPAIR, 1)
		extract_obj(obj4)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		if ch.Hit >= gear_pl(ch) {
			send_to_char(ch, libc.CString("Your body is already in peak condition.\r\n"))
			return
		} else if GET_SKILL(ch, SKILL_REPAIR) < axion_dice(0) {
			act(libc.CString("You try to repair your body but screw up.."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n tries to repair $s body but screws up.."), TRUE, ch, nil, nil, TO_ROOM)
			extract_obj(obj4)
			improve_skill(ch, SKILL_REPAIR, 1)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
			return
		} else {
			act(libc.CString("You use the repair kit to fix part of your body..."), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n works on their body with a repair kit."), TRUE, ch, nil, nil, TO_ROOM)
			var mult int64 = int64(GET_SKILL(ch, SKILL_REPAIR))
			var add int64 = int64(((float64(gear_pl(ch)) * 0.005) + 10) * float64(mult))
			ch.Hit += add
			extract_obj(obj4)
			if ch.Hit > gear_pl(ch) {
				ch.Hit = gear_pl(ch)
				send_to_char(ch, libc.CString("Your body has been totally repaired.\r\n"))
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
			} else {
				send_to_char(ch, libc.CString("Your body still needs some work done to it.\r\n"))
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
			}
		}
	}
}
func innate_remove(ch *char_data, inn *innate_node) {
	var temp *innate_node
	if ch.Innate == nil {
		core_dump_real(libc.CString(__FILE__), __LINE__)
		return
	}
	if inn == ch.Innate {
		ch.Innate = inn.Next
	} else {
		temp = ch.Innate
		for temp != nil && temp.Next != inn {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = inn.Next
		}
	}
	libc.Free(unsafe.Pointer(inn))
}
func innate_add(ch *char_data, innate int, timer int) {
	var inn *innate_node
	inn = new(innate_node)
	inn.Timer = timer
	inn.Spellnum = innate
	inn.Next = ch.Innate
	ch.Innate = inn
}
func is_innate(ch *char_data, spellnum int) int {
	switch spellnum {
	case SPELL_FAERIE_FIRE:
		if ch.Race == RACE_ANIMAL {
			return TRUE
		}
	}
	return FALSE
}
func is_innate_ready(ch *char_data, spellnum int) int {
	var (
		inn      *innate_node
		next_inn *innate_node
	)
	for inn = ch.Innate; inn != nil; inn = next_inn {
		next_inn = inn.Next
		if inn.Spellnum == spellnum {
			return FALSE
		}
	}
	return TRUE
}
func add_innate_timer(ch *char_data, spellnum int) {
	var timer int = 6
	switch spellnum {
	case SPELL_FAERIE_FIRE:
		timer = 6
	case ABIL_LAY_HANDS:
		timer = 12
	}
	if is_innate_ready(ch, spellnum) != 0 {
		innate_add(ch, spellnum, timer)
	} else {
		send_to_char(ch, libc.CString("BUG!\r\n"))
	}
}
func add_innate_affects(ch *char_data) {
	switch ch.Race {
	case RACE_DEMON:
		fallthrough
	case RACE_ICER:
		fallthrough
	case RACE_ANDROID:
		fallthrough
	case RACE_BIO:
		affect_modify(ch, APPLY_NONE, 0, 0, AFF_INFRAVISION, TRUE != 0)
	}
	affect_total(ch)
}
func update_innate(ch *char_data) {
	var (
		inn      *innate_node
		next_inn *innate_node
	)
	for inn = ch.Innate; inn != nil; inn = next_inn {
		next_inn = inn.Next
		if inn.Timer > 0 {
			inn.Timer--
		} else {
			switch inn.Spellnum {
			case ABIL_LAY_HANDS:
				send_to_char(ch, libc.CString("Your special healing abilities have returned.\r\n"))
			default:
				send_to_char(ch, libc.CString("You are now able to use your innate %s again.\r\n"), spell_info[inn.Spellnum].Name)
			}
			innate_remove(ch, inn)
		}
	}
}
func spell_in_book(obj *obj_data, spellnum int) int {
	var (
		i     int
		found bool = FALSE != 0
	)
	if obj.Sbinfo == nil {
		return FALSE
	}
	for i = 0; i < SPELLBOOK_SIZE; i++ {
		if (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(i)))).Spellname == spellnum {
			found = TRUE != 0
			break
		}
	}
	if found {
		return 1
	}
	return 0
}
func spell_in_scroll(obj *obj_data, spellnum int) int {
	if (obj.Value[VAL_SCROLL_SPELL1]) == spellnum {
		return TRUE
	}
	return FALSE
}
func spell_in_domain(ch *char_data, spellnum int) int {
	if spell_info[spellnum].Domain == -1 {
		return FALSE
	}
	return TRUE
}

var freeres [9]room_vnum = [9]room_vnum{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}

func do_resurrect(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		rm      room_rnum
		af      *affected_type
		next_af *affected_type
	)
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("Sorry, only players get spirits.\r\n"))
		return
	}
	if !AFF_FLAGGED(ch, AFF_SPIRIT) {
		send_to_char(ch, libc.CString("But you're not even dead!\r\n"))
		return
	}
	send_to_char(ch, libc.CString("You take an experience penalty and pray for charity resurrection.\r\n"))
	gain_exp(ch, int64(-(level_exp(ch, GET_LEVEL(ch)) - level_exp(ch, GET_LEVEL(ch)-1))))
	for af = ch.Affected; af != nil; af = next_af {
		next_af = af.Next
		if af.Location == APPLY_NONE && af.Type == -1 && (af.Bitvector == AFF_SPIRIT || af.Bitvector == AFF_ETHEREAL) {
			affect_remove(ch, af)
		}
	}
	if ch.Hit < 1 {
		ch.Hit = 1
	}
	if (func() room_rnum {
		rm = real_room(freeres[ALIGN_TYPE(ch)])
		return rm
	}()) == room_rnum(-1) {
		rm = real_room(config_info.Room_nums.Mortal_start_room)
	}
	if rm != room_rnum(-1) {
		char_from_room(ch)
		char_to_room(ch, rm)
		look_at_room(ch.In_room, ch, 0)
	}
	act(libc.CString("$n's body forms in a pool of @Bblue light@n."), TRUE, ch, nil, nil, TO_ROOM)
}
func show_clan_info(ch *char_data) {
	send_to_char(ch, libc.CString("@c----------------------------------------\r\n@c|@WProvided by@D: @YAlister of Aeonian Dreams@c|\r\n@c|@YWith many changes made by Iovan.      @c|\r\n@c----------------------------------------@w\r\n  Commands are:\r\n@c--@YClan Members Only@c-----@w\r\n  clan members\r\n  clan bank\r\n  clan deposit\r\n  clan leave    <clan>\r\n@c--@YClan Mod/Highrank@c-----@w\r\n  clan decline  <person> <clan>\r\n  clan enroll   <person> <clan>\r\n@c--@YClan Moderators Only@c--@w\r\n  clan withdraw\r\n  clan infow\r\n  clan setjoin  <free | restricted> <clan>\r\n  clan setleave <free | restricted> <clan>\r\n  clan expel    <person> <clan>\r\n  clan highrank <new highrank title>\r\n  clan midrank  <new midrank title>\r\n  clan rank     <person> < 0 / 1 / or 2>\r\n  clan makemod  <person> <clan>\r\n@c--@YAnyone@c----------------@w\r\n  clan list\r\n  clan info     <clan>\r\n  clan apply    <clan>\r\n  clan join     <clan>\r\n"))
	if ch.Admlevel >= ADMLVL_IMPL {
		send_to_char(ch, libc.CString("@c--@YImmort@c----------------@w\r\n  clan create   <clan>\r\n  clan destroy  <clan>\r\n  clan reload   <clan>\r\n  clan bset     <clan>\n"))
	}
}
func do_clan(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg1 [100]byte
		arg2 [2048]byte
	)
	if ch == nil || IS_NPC(ch) {
		return
	}
	half_chop(argument, &arg1[0], &arg2[0])
	if arg1[0] == 0 {
		show_clan_info(ch)
		return
	}
	if ch.Admlevel >= ADMLVL_IMPL {
		if C.strcmp(&arg1[0], libc.CString("create")) == 0 {
			if arg2[0] == 0 {
				show_clan_info(ch)
			} else if isClan(&arg2[0]) {
				send_to_char(ch, libc.CString("There is already a clan with the name, %s.\r\n"), &arg2[0])
			} else {
				send_to_char(ch, libc.CString("You create a clan with the name, %s.\r\n"), &arg2[0])
				clanCreate(&arg2[0])
				mudlog(BRF, MAX(ADMLVL_GOD, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("(GC) %s has created a clan named %s."), GET_NAME(ch), &arg2[0])
			}
			return
		} else if C.strcmp(&arg1[0], libc.CString("destroy")) == 0 {
			if arg2[0] == 0 {
				show_clan_info(ch)
			} else if !isClan(&arg2[0]) {
				send_to_char(ch, libc.CString("No clan with the name %s exists.\r\n"), &arg2[0])
			} else {
				send_to_char(ch, libc.CString("You destroy %s.\r\n"), &arg2[0])
				clanDestroy(&arg2[0])
			}
			return
		} else if C.strcmp(&arg1[0], libc.CString("bset")) == 0 {
			if arg2[0] == 0 {
				show_clan_info(ch)
			} else if !isClan(&arg2[0]) {
				send_to_char(ch, libc.CString("No clan with the name %s exists.\r\n"), &arg2[0])
			} else if ch.Admlevel < 5 {
				show_clan_info(ch)
			} else {
				clanBSET(&arg2[0], ch)
			}
			return
		} else if C.strcmp(&arg1[0], libc.CString("reload")) == 0 {
			if arg2[0] == 0 {
				show_clan_info(ch)
			} else if !isClan(&arg2[0]) {
				send_to_char(ch, libc.CString("No clan with the name %s exists.\r\n"), &arg2[0])
			} else if clanReload(&arg2[0]) {
				send_to_char(ch, libc.CString("Data for %s has been reloaded.\r\n"), &arg2[0])
			} else {
				send_to_char(ch, libc.CString("Failed to reload the data for %s.\r\n"), &arg2[0])
			}
			return
		}
	}
	if C.strcmp(&arg1[0], libc.CString("apply")) == 0 {
		if arg2[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg2[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg2[0])
		} else if ch.Clan != nil && clanIsMember(ch.Clan, ch) {
			send_to_char(ch, libc.CString("You are already a member of %s.\r\n"), ch.Clan)
		} else if clanOpenJoin(&arg2[0]) {
			send_to_char(ch, libc.CString("You can just join %s, it is open.\r\n"), &arg2[0])
			return
		} else {
			if ch.Clan != nil && checkCLAN(ch) == TRUE {
				checkAPP(ch)
				send_to_char(ch, libc.CString("You stop applying to %s\r\n"), ch.Clan)
				clanDecline(ch.Clan, ch)
				if ch.Clan != nil {
					libc.Free(unsafe.Pointer(ch.Clan))
				}
				ch.Clan = C.strdup(libc.CString("None."))
			}
			send_to_char(ch, libc.CString("You apply to become a member of %s.\r\n"), &arg2[0])
			clanApply(&arg2[0], ch)
			return
		}
		return
	}
	if C.strcmp(&arg1[0], libc.CString("join")) == 0 {
		if arg2[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg2[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg2[0])
		} else if clanIsMember(&arg2[0], ch) {
			send_to_char(ch, libc.CString("You are already a member of %s.\r\n"), &arg2[0])
		} else if ch.Clan != nil && checkCLAN(ch) == TRUE && unsafe.Pointer(C.strstr(ch.Clan, libc.CString("Applying"))) == unsafe.Pointer(uintptr(FALSE)) {
			send_to_char(ch, libc.CString("You are already a member of %s, you need to leave it first.\r\n"), ch.Clan)
		} else if clanOpenJoin(&arg2[0]) {
			if ch.Clan != nil && checkCLAN(ch) == TRUE {
				checkAPP(ch)
				send_to_char(ch, libc.CString("You stop applying to %s\r\n"), ch.Clan)
				clanDecline(ch.Clan, ch)
				if ch.Clan != nil {
					libc.Free(unsafe.Pointer(ch.Clan))
				}
				ch.Clan = C.strdup(libc.CString("None."))
			}
			send_to_char(ch, libc.CString("You are now a member of %s.\r\n"), &arg2[0])
			clanInduct(&arg2[0], ch)
			return
		} else {
			send_to_char(ch, libc.CString("%s isn't open, you must apply instead.\r\n"), &arg2[0])
			return
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("leave")) == 0 {
		if arg2[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg2[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg2[0])
		} else if !clanIsMember(&arg2[0], ch) {
			send_to_char(ch, libc.CString("You aren't even a member of %s.\r\n"), &arg2[0])
		} else if !clanOpenLeave(&arg2[0]) && !clanIsModerator(&arg2[0], ch) {
			send_to_char(ch, libc.CString("You must be expelled from %s in order to leave it.\r\n"), &arg2[0])
		} else {
			send_to_char(ch, libc.CString("You are no longer a member of %s.\r\n"), &arg2[0])
			clanExpel(&arg2[0], ch)
			return
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("infow")) == 0 {
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if int(libc.BoolToInt(clanIsModerator(ch.Clan, ch))) == FALSE {
				send_to_char(ch, libc.CString("You must be a moderator to edit the clan's information.\r\n"))
			} else {
				clanINFOW(ch.Clan, ch)
				ch.Act[int(PLR_WRITING/32)] &= bitvector_t(^(1 << (int(PLR_WRITING % 32))))
			}
			return
		}
	} else if C.strcmp(&arg1[0], libc.CString("deposit")) == 0 {
		var bank int = 0
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsMember(ch.Clan, ch) && !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
				return
			} else if !ROOM_FLAGGED(ch.In_room, ROOM_CBANK) && int(libc.BoolToInt(clanBANY(ch.Clan, ch))) == FALSE {
				send_to_char(ch, libc.CString("You are not in your clan bank and your clan doesn't have bank anywhere.\r\n"))
				return
			} else if arg2[0] == 0 {
				send_to_char(ch, libc.CString("How much do you want to deposit?\r\n"))
				return
			} else if libc.Atoi(libc.GoString(&arg2[0])) <= 0 {
				send_to_char(ch, libc.CString("It needs to be a value higher than 0...\r\n"))
				return
			} else if ch.Gold < libc.Atoi(libc.GoString(&arg2[0])) {
				send_to_char(ch, libc.CString("You do not have that much to deposit!\r\n"))
				return
			} else {
				bank = libc.Atoi(libc.GoString(&arg2[0]))
				ch.Gold -= bank
				clanBANKADD(ch.Clan, ch, bank)
				send_to_char(ch, libc.CString("You have deposited %s into the clan bank.\r\n"), add_commas(int64(bank)))
			}
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("highrank")) == 0 {
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not leading a clan.\r\n"))
				return
			} else if arg2[0] == 0 {
				send_to_char(ch, libc.CString("What name are you going to make the rank?\r\n"))
				return
			} else if C.strlen(&arg2[0]) > 20 {
				send_to_char(ch, libc.CString("The name length can't be longer than 20 characters.\r\n"))
				return
			} else if C.strstr(&arg2[0], libc.CString("@")) != nil {
				send_to_char(ch, libc.CString("No colorcode allowed in the ranks.\r\n"))
				return
			} else {
				clanHIGHRANK(ch.Clan, ch, &arg2[0])
				send_to_char(ch, libc.CString("High rank set.\r\n"))
			}
			return
		}
	} else if C.strcmp(&arg1[0], libc.CString("midrank")) == 0 {
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not leading a clan.\r\n"))
				return
			} else if arg2[0] == 0 {
				send_to_char(ch, libc.CString("What name are you going to make the rank?\r\n"))
				return
			} else if C.strlen(&arg2[0]) > 20 {
				send_to_char(ch, libc.CString("The name length can't be longer than 20 characters.\r\n"))
				return
			} else if C.strstr(&arg2[0], libc.CString("@")) != nil {
				send_to_char(ch, libc.CString("No colorcode allowed in the ranks.\r\n"))
				return
			} else {
				clanMIDRANK(ch.Clan, ch, &arg2[0])
				send_to_char(ch, libc.CString("Mid rank set.\r\n"))
			}
			return
		}
	} else if C.strcmp(&arg1[0], libc.CString("rank")) == 0 {
		var (
			vict *char_data = nil
			arg3 [100]byte
			name [100]byte
		)
		half_chop(&arg2[0], &name[0], &arg3[0])
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not leading a clan.\r\n"))
				return
			} else if arg2[0] == 0 {
				send_to_char(ch, libc.CString("Who's rank do you want to change?\r\n"))
				return
			} else if (func() *char_data {
				vict = get_char_vis(ch, &name[0], nil, 1<<1)
				return vict
			}()) == nil {
				send_to_char(ch, libc.CString("That person is no where to be found in the entire universe.\r\n"))
				return
			} else if vict.Clan == nil || C.strcmp(vict.Clan, libc.CString("None.")) == 0 {
				send_to_char(ch, libc.CString("That person is not even in a clan, let alone your's.\r\n"))
				return
			} else if !clanIsMember(ch.Clan, vict) {
				send_to_char(ch, libc.CString("You can only rank those in your clan and only if below leader.\r\n"))
				return
			} else if clanIsModerator(ch.Clan, vict) {
				send_to_char(ch, libc.CString("You can't rank a fellow leader, you require imm assistance.\r\n"))
				return
			} else if arg3[0] == 0 {
				send_to_char(ch, libc.CString("What rank do you want to set them to?\r\n[ 0 = Member, 1 = Midrank, 2 = Highrank]\r\n"))
				return
			}
			var num int = libc.Atoi(libc.GoString(&arg3[0]))
			if num < 0 || num > 2 {
				send_to_char(ch, libc.CString("It must be above zero and lower than three...\r\n"))
				return
			} else if vict.Crank == num {
				send_to_char(ch, libc.CString("They are already that rank!\r\n"))
				return
			} else if vict.Crank > num {
				clanRANK(ch.Clan, ch, vict, num)
				switch num {
				case 0:
					send_to_char(ch, libc.CString("You demote %s.\r\n"), GET_NAME(vict))
					send_to_char(vict, libc.CString("%s has demoted your clan rank to member!\r\n"), GET_NAME(ch))
				case 1:
					send_to_char(ch, libc.CString("You demote %s.\r\n"), GET_NAME(vict))
					send_to_char(vict, libc.CString("%s has demoted your clan rank to midrank!\r\n"), GET_NAME(ch))
				}
				return
			} else if vict.Crank < num {
				clanRANK(ch.Clan, ch, vict, num)
				switch num {
				case 1:
					send_to_char(ch, libc.CString("You promote %s.\r\n"), GET_NAME(vict))
					send_to_char(vict, libc.CString("%s has promoted your clan rank to midrank!\r\n"), GET_NAME(ch))
				case 2:
					send_to_char(ch, libc.CString("You promote %s.\r\n"), GET_NAME(vict))
					send_to_char(vict, libc.CString("%s has promoted your clan rank to highrank!\r\n"), GET_NAME(ch))
				}
				return
			}
		}
	} else if C.strcmp(&arg1[0], libc.CString("withdraw")) == 0 {
		var bank int = 0
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsMember(ch.Clan, ch) && !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
				return
			} else if arg2[0] == 0 {
				send_to_char(ch, libc.CString("How much do you want to withdraw?\r\n"))
				return
			} else if !clanIsModerator(ch.Clan, ch) && ch.Admlevel < ADMLVL_IMPL {
				send_to_char(ch, libc.CString("You do not have the power to withdraw from the clan bank.\r\n"))
				return
			} else if libc.Atoi(libc.GoString(&arg2[0])) <= 0 {
				send_to_char(ch, libc.CString("It needs to be a value higher than 0...\r\n"))
				return
			} else if ch.Gold+libc.Atoi(libc.GoString(&arg2[0])) > GOLD_CARRY(ch) {
				send_to_char(ch, libc.CString("You can not hold that much zenni!\r\n"))
				return
			} else {
				bank = libc.Atoi(libc.GoString(&arg2[0]))
				if clanBANKSUB(ch.Clan, ch, bank) {
					send_to_char(ch, libc.CString("You have withdrawn %s from the clan bank.\r\n"), add_commas(int64(bank)))
					ch.Gold += bank
				} else {
					send_to_char(ch, libc.CString("There isn't that much in the clan's bank!\r\n"))
				}
			}
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("bank")) == 0 {
		var bank int = 0
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
			return
		} else {
			if !clanIsMember(ch.Clan, ch) && !clanIsModerator(ch.Clan, ch) {
				send_to_char(ch, libc.CString("You are not in a clan.\r\n"))
				return
			}
			bank = clanBANK(ch.Clan, ch)
			send_to_char(ch, libc.CString("@W[ @C%-20s @W]@w has @D(@Y%s@D)@w zenni in its clan bank.\r\n"), ch.Clan, add_commas(int64(bank)))
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("members")) == 0 {
		if ch.Clan == nil || C.strcmp(ch.Clan, libc.CString("None.")) == 0 {
			send_to_char(ch, libc.CString("You are not even in a clan.\r\n"))
			return
		} else {
			handle_clan_member_list(ch)
		}
	} else if C.strcmp(&arg1[0], libc.CString("expel")) == 0 {
		var (
			is_file  int = FALSE
			player_i int = 0
			vict     *char_data
			arg3     [100]byte
			name     [2048]byte
			name1    [100]byte
		)
		half_chop(&arg2[0], &name1[0], &arg3[0])
		if arg3[0] == 0 || name1[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg3[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg3[0])
		} else if !clanIsModerator(&arg3[0], ch) && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("Only leaders can expel people from a clan.\r\n"))
		} else if clanOpenJoin(&arg3[0]) {
			send_to_char(ch, libc.CString("You can't kick someone out of an open-join clan.\r\n"))
		} else if (func() *char_data {
			vict = get_char_vis(ch, &name1[0], nil, 1<<1)
			return vict
		}()) == nil {
			vict = new(char_data)
			clear_char(vict)
			vict.Player_specials = new(player_special_data)
			stdio.Sprintf(&name[0], "%s", rIntro(ch, &name1[0]))
			if (func() int {
				player_i = load_char(&name[0], vict)
				return player_i
			}()) > -1 {
				is_file = TRUE
				if !clanIsMember(&arg3[0], vict) {
					send_to_char(ch, libc.CString("%s isn't even a member of %s.\r\n"), GET_NAME(vict), &arg3[0])
				} else if clanIsModerator(&arg3[0], vict) && ch.Admlevel < ADMLVL_IMPL {
					send_to_char(ch, libc.CString("You do not have the power to kick a leader out of %s.\r\n"), &arg3[0])
				} else {
					send_to_char(ch, libc.CString("You expel %s from %s.\r\n"), GET_NAME(vict), &arg3[0])
					clanExpel(&arg3[0], vict)
				}
			} else if (func() int {
				player_i = load_char(&name1[0], vict)
				return player_i
			}()) > -1 {
				is_file = TRUE
				if !clanIsMember(&arg3[0], vict) {
					send_to_char(ch, libc.CString("%s isn't even a member of %s.\r\n"), GET_NAME(vict), &arg3[0])
				} else if clanIsModerator(&arg3[0], vict) && ch.Admlevel < ADMLVL_IMPL {
					send_to_char(ch, libc.CString("You do not have the power to kick a leader out of %s.\r\n"), &arg3[0])
				} else {
					send_to_char(ch, libc.CString("You expel %s from %s.\r\n"), GET_NAME(vict), &arg3[0])
					clanExpel(&arg3[0], vict)
				}
			} else {
				free_char(vict)
				send_to_char(ch, libc.CString("%s does not seem to exist.\r\n"), &name1[0])
				return
			}
			if is_file == TRUE {
				free_char(vict)
			}
			return
		} else if !clanIsMember(&arg3[0], vict) {
			send_to_char(ch, libc.CString("%s isn't even a member of %s.\r\n"), GET_NAME(vict), &arg3[0])
		} else if clanIsModerator(&arg3[0], vict) && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("You do not have the power to kick a leader out of %s.\r\n"), &arg3[0])
		} else {
			send_to_char(ch, libc.CString("You expel %s from %s.\r\n"), GET_NAME(vict), &arg3[0])
			send_to_char(vict, libc.CString("You have been expelled from %s.\r\n"), &arg3[0])
			clanExpel(&arg3[0], vict)
			return
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("decline")) == 0 {
		var (
			vict *char_data
			arg3 [100]byte
			name [100]byte
		)
		half_chop(&arg2[0], &name[0], &arg3[0])
		if arg3[0] == 0 || name[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg3[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg3[0])
		} else if int(libc.BoolToInt(clanIsModerator(&arg3[0], ch))) == FALSE && int(libc.BoolToInt(clanIsMember(&arg3[0], ch))) == FALSE && ch.Admlevel < ADMLVL_IMPL || ch.Crank < 2 && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("Only leaders or highrank can decline people from entering a clan.\r\n"))
		} else if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s is not around at the moment.\r\n"), &name[0])
		} else if !clanIsApplicant(&arg3[0], vict) {
			send_to_char(ch, libc.CString("%s isn't applying to join %s.\r\n"), GET_NAME(vict), &arg3[0])
		} else {
			send_to_char(ch, libc.CString("You decline %s enterance to %s.\r\n"), GET_NAME(vict), &arg3[0])
			send_to_char(vict, libc.CString("You have been declined enterance to %s.\r\n"), &arg3[0])
			clanDecline(&arg3[0], vict)
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("enroll")) == 0 {
		var (
			vict *char_data
			arg3 [100]byte
			name [100]byte
		)
		half_chop(&arg2[0], &name[0], &arg3[0])
		if arg3[0] == 0 || name[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg3[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg3[0])
		} else if !clanIsMember(&arg3[0], ch) && ch.Admlevel < 1 {
			send_to_char(ch, libc.CString("You are not in that clan.\r\n"))
		} else if !clanIsModerator(&arg3[0], ch) && ch.Admlevel < ADMLVL_IMPL && ch.Crank < 2 {
			send_to_char(ch, libc.CString("Only leaders or captains can enroll people into their clan.\r\n"))
		} else if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s is not around at the moment.\r\n"), &name[0])
		} else if !clanIsApplicant(&arg3[0], vict) {
			send_to_char(ch, libc.CString("%s isn't applying to join %s.\r\n"), GET_NAME(vict), &arg3[0])
		} else if vict.Trp < 5 {
			send_to_char(ch, libc.CString("%s needs to have at least earned 5 RPP at some point to join a clan.\r\n"), GET_NAME(vict))
		} else {
			send_to_char(ch, libc.CString("You enroll %s into %s.\r\n"), GET_NAME(vict), &arg3[0])
			send_to_char(vict, libc.CString("You have been enrolled into %s.\r\n"), &arg3[0])
			clanInduct(&arg3[0], vict)
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("makemod")) == 0 {
		var (
			vict *char_data
			arg3 [100]byte
			name [100]byte
		)
		half_chop(&arg2[0], &name[0], &arg3[0])
		if arg3[0] == 0 || name[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&arg3[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &arg3[0])
		} else if !clanIsModerator(&arg3[0], ch) && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("Only leaders can make other people in a clan a leader.\r\n"))
		} else if (func() *char_data {
			vict = get_char_vis(ch, &name[0], nil, 1<<1)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s is not around at the moment.\r\n"), &name[0])
		} else {
			send_to_char(ch, libc.CString("You make %s a leader of %s.\r\n"), GET_NAME(vict), &arg3[0])
			send_to_char(vict, libc.CString("You have been made a leader of %s.\r\n"), &arg3[0])
			clanMakeModerator(&arg3[0], vict)
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("setleave")) == 0 {
		var (
			name    [100]byte
			setting [100]byte
		)
		half_chop(&arg2[0], &setting[0], &name[0])
		if name[0] == 0 || setting[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&name[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &name[0])
		} else if !clanIsModerator(&name[0], ch) && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("Only leaders can change that.\r\n"))
		} else if C.strcmp(&setting[0], libc.CString("free")) == 0 {
			send_to_char(ch, libc.CString("Members of %s are free to leave as they please.\r\n"), &name[0])
			clanSetOpenLeave(&name[0], TRUE)
		} else if C.strcmp(&setting[0], libc.CString("restricted")) == 0 {
			send_to_char(ch, libc.CString("Members of %s can no longer leave as they please.\r\n"), &name[0])
			clanSetOpenLeave(&name[0], FALSE)
		} else {
			send_to_char(ch, libc.CString("Leave access may only be set to free or restricted.\r\n"))
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("setjoin")) == 0 {
		var (
			name    [100]byte
			setting [100]byte
		)
		half_chop(&arg2[0], &setting[0], &name[0])
		if name[0] == 0 || setting[0] == 0 {
			show_clan_info(ch)
		} else if !isClan(&name[0]) {
			send_to_char(ch, libc.CString("%s is not a valid clan.\r\n"), &name[0])
		} else if !clanIsModerator(&name[0], ch) && ch.Admlevel < ADMLVL_IMPL {
			send_to_char(ch, libc.CString("Only leaders can change that.\r\n"))
		} else if C.strcmp(&setting[0], libc.CString("free")) == 0 {
			send_to_char(ch, libc.CString("People may now freely join %s.\r\n"), &name[0])
			clanSetOpenJoin(&name[0], TRUE)
		} else if C.strcmp(&setting[0], libc.CString("restricted")) == 0 {
			send_to_char(ch, libc.CString("People must be enrolled into %s to join.\r\n"), &name[0])
			clanSetOpenJoin(&name[0], FALSE)
		} else {
			send_to_char(ch, libc.CString("Leave access my only be set to free or restricted.\r\n"))
		}
		return
	} else if C.strcmp(&arg1[0], libc.CString("list")) == 0 {
		listClans(ch)
	} else if C.strcmp(&arg1[0], libc.CString("info")) == 0 {
		if arg2[0] == 0 {
			show_clan_info(ch)
		} else {
			listClanInfo(&arg2[0], ch)
		}
	} else {
		show_clan_info(ch)
		send_to_char(ch, libc.CString("These are viable options.\r\n"))
	}
}
func do_pagelength(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	if IS_NPC(ch) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("You current page length is set to %d lines.\r\n"), ch.Player_specials.Page_length)
	} else if is_number(&arg[0]) != 0 {
		ch.Player_specials.Page_length = uint8(int8(MIN(MAX(libc.Atoi(libc.GoString(&arg[0])), 5), 50)))
		send_to_char(ch, libc.CString("Okay, your page length is now set to %d lines.\r\n"), ch.Player_specials.Page_length)
	} else {
		send_to_char(ch, libc.CString("Please specify a number of lines (5 - 50).\r\n"))
	}
}
func do_aid(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict     *char_data
		obj      *obj_data = nil
		aid_obj  *obj_data = nil
		aid_prod *obj_data = nil
		next_obj *obj_data
		arg      [2048]byte
		arg2     [2048]byte
		dc       int = 0
		found    int = FALSE
		num      int = 47
		num2     int = 0
		survival int = 0
	)
	_ = survival
	if IS_NPC(ch) {
		return
	}
	if GET_SKILL(ch, SKILL_SURVIVAL) != 0 {
		survival = 1
	}
	two_arguments(argument, &arg[0], &arg2[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: aid heal (target)\r\n"))
		send_to_char(ch, libc.CString("        aid adrenex\r\n"))
		send_to_char(ch, libc.CString("        aid antitoxin\r\n"))
		send_to_char(ch, libc.CString("        aid salve\r\n"))
		send_to_char(ch, libc.CString("        aid formula-82\r\n"))
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("adrenex")) == 0 {
		num = 380
		num2 = 381
	} else if C.strcasecmp(&arg[0], libc.CString("antitoxin")) == 0 {
		num = 380
		num2 = 383
	} else if C.strcasecmp(&arg[0], libc.CString("salve")) == 0 {
		num = 380
		num2 = 382
	} else if C.strcasecmp(&arg[0], libc.CString("formula-82")) == 0 {
		num = 380
		num2 = 385
	}
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if found == FALSE && GET_OBJ_VNUM(obj) == obj_vnum(num) && (!OBJ_FLAGGED(obj, ITEM_BROKEN) && !OBJ_FLAGGED(obj, ITEM_FORGED)) {
			found = TRUE
			aid_obj = obj
		}
	}
	if found == FALSE || aid_obj == nil {
		if num == 47 {
			send_to_char(ch, libc.CString("You need bandages to be able to use first aid.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You need a TCX-Medical Equipment Construction Kit.\r\n"))
		}
		return
	}
	if num == 47 {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("Apply first aid to who?\r\n"))
			return
		} else if IS_NPC(vict) {
			send_to_char(ch, libc.CString("What ever for?\r\n"))
			return
		} else if vict.Race == RACE_ANDROID {
			send_to_char(ch, libc.CString("They are an android!\r\n"))
			return
		}
		if !AFF_FLAGGED(vict, AFF_SPIRIT) && !PLR_FLAGGED(vict, PLR_BANDAGED) {
			if vict != ch {
				send_to_char(ch, libc.CString("You attempt to lend first aid to %s.\r\n"), GET_NAME(vict))
			}
			act(libc.CString("$n attempts to bandage $N's wounds."), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			dc = axion_dice(0)
			if (GET_SKILL(ch, SKILL_FIRST_AID) + 1) > dc {
				send_to_char(ch, libc.CString("You bandage %s's wounds.\r\n"), GET_NAME(vict))
				var roll int64 = int64(float64((gear_pl(vict)/100)*int64(ch.Aff_abils.Wis/4)) + float64(gear_pl(vict))*0.25)
				if (ch.Bonuses[BONUS_HEALER]) > 0 {
					roll += int64(float64(roll) * 0.1)
				}
				vict.Hit += roll
				if vict.Hit > gear_pl(vict) {
					vict.Hit = gear_pl(vict)
				}
				if vict.Suppression > 0 && vict.Hit > ((vict.Max_hit/100)*vict.Suppression) {
					vict.Hit = (vict.Max_hit / 100) * vict.Suppression
					send_to_char(vict, libc.CString("@mYou are healed to your suppression limit.@n\r\n"))
				}
				send_to_char(vict, libc.CString("Your wounds are bandaged by %s!\r\n"), GET_NAME(ch))
				act(libc.CString("$n's wounds are stablized by $N!"), TRUE, vict, nil, unsafe.Pointer(ch), TO_NOTVICT)
				vict.Act[int(PLR_BANDAGED/32)] |= bitvector_t(1 << (int(PLR_BANDAGED % 32)))
				extract_obj(aid_obj)
			} else {
				if vict != ch {
					send_to_char(ch, libc.CString("You fail to bandage their wounds properly, wasting the set of bandages...\r\n"))
					act(libc.CString("$N fails to bandage $n's wounds properly, wasting an entire set of bandages..."), TRUE, vict, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("$N fails to bandage your wounds properly, wasting an entire set of bandages..."), TRUE, vict, nil, unsafe.Pointer(ch), TO_CHAR)
				} else {
					act(libc.CString("$N fails to bandage $s wounds properly, wasting an entire set of bandages..."), TRUE, vict, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("You fail to bandage your wounds properly, wasting an entire set of bandages..."), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
				}
				extract_obj(aid_obj)
			}
			improve_skill(ch, SKILL_FIRST_AID, 1)
		} else if PLR_FLAGGED(vict, PLR_BANDAGED) {
			send_to_char(ch, libc.CString("They are already bandaged!\r\n"))
		} else if AFF_FLAGGED(vict, AFF_SPIRIT) {
			send_to_char(ch, libc.CString("The dead don't need first aid.\r\n"))
		} else {
			send_to_char(ch, libc.CString("They apparently do not need bandaging.\r\n"))
		}
	} else if num2 == 381 {
		if GET_SKILL(ch, SKILL_FIRST_AID) < 65 {
			send_to_char(ch, libc.CString("You need at least a skill level of 65 in First Aid.\r\n"))
			return
		} else {
			if GET_SKILL(ch, SKILL_FIRST_AID) < axion_dice(15) {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. As you begin to construct an Andrenex Adreneline Injector you screw up and end up breaking the water tight seal. The adreneline leaks out and is wasted.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A frown forms on $s face as it appears that $e has failed.@n"), TRUE, ch, nil, nil, TO_ROOM)
				extract_obj(aid_obj)
			} else {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. Your knowledge in basic medical devices and treatments helps you as you successfully construct an Adrenex Adreneline Injector@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A moment later $e holds up a completed Adrenex Adreneline Injector!@n"), TRUE, ch, nil, nil, TO_ROOM)
				aid_prod = read_object(obj_vnum(num2), VIRTUAL)
				add_unique_id(aid_prod)
				obj_to_char(aid_prod, ch)
				extract_obj(aid_obj)
				improve_skill(ch, SKILL_FIRST_AID, 1)
			}
		}
	} else if num2 == 382 {
		if GET_SKILL(ch, SKILL_FIRST_AID) < 50 {
			send_to_char(ch, libc.CString("You need at least a skill level of 50 in First Aid.\r\n"))
			return
		} else {
			if GET_SKILL(ch, SKILL_FIRST_AID) < axion_dice(10) {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. As you go to put the salve ingredients into the kit's salve compartment and set the temperature you accidentally set it too high. The salve is burned and ruined. Yes you managed to burn a burn salve.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A frown forms on $s face as it appears that $e has failed.@n"), TRUE, ch, nil, nil, TO_ROOM)
				extract_obj(aid_obj)
			} else {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. Your knowledge in basic medical devices and treatments helps you as you successfully boil a burn salve to perfection and it is automatically placed in a jar.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A moment later $e holds up a jar of burn salve!@n"), TRUE, ch, nil, nil, TO_ROOM)
				aid_prod = read_object(obj_vnum(num2), VIRTUAL)
				add_unique_id(aid_prod)
				obj_to_char(aid_prod, ch)
				extract_obj(aid_obj)
				improve_skill(ch, SKILL_FIRST_AID, 1)
			}
		}
	} else if num2 == 383 {
		if GET_SKILL(ch, SKILL_FIRST_AID) < 40 {
			send_to_char(ch, libc.CString("You need at least a skill level of 40 in First Aid.\r\n"))
			return
		} else {
			if GET_SKILL(ch, SKILL_FIRST_AID) < axion_dice(5) {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. As you complete the Antitoxin Injector you notice that you didn't seal the syringe properly and it all leaks out.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A frown forms on $s face as it appears that $e has failed.@n"), TRUE, ch, nil, nil, TO_ROOM)
				extract_obj(aid_obj)
			} else {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. Your knowledge in basic medical devices and treatments helps you as you successfully assemble the Antitoxin Injector.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A moment later $e holds up a completed Antitoxin Injector!@n"), TRUE, ch, nil, nil, TO_ROOM)
				aid_prod = read_object(obj_vnum(num2), VIRTUAL)
				add_unique_id(aid_prod)
				obj_to_char(aid_prod, ch)
				extract_obj(aid_obj)
				improve_skill(ch, SKILL_FIRST_AID, 1)
			}
		}
	} else if num2 == 385 {
		if GET_SKILL(ch, SKILL_FIRST_AID) < 40 {
			send_to_char(ch, libc.CString("You need at least a skill level of 40 in First Aid.\r\n"))
			return
		} else {
			if GET_SKILL(ch, SKILL_FIRST_AID) < axion_dice(15) {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. As you complete a vial of Formula 82 you notice that you read the mixture measurements wrong. You dispose of the vile vial immediately.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A frown forms on $s face as it appears that $e has failed.@n"), TRUE, ch, nil, nil, TO_ROOM)
				extract_obj(aid_obj)
			} else {
				act(libc.CString("@WYou unlock and open the TCX-M.E.C.K. case. The case hisses as its lid opens. Your knowledge in basic medical devices and treatments helps you as you successfully assemble a vial of Formula 82.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@W holds a steel case up and opens it. The case hisses as its lid opens. @C$n@W wastes no time as $e reaches into the case and begins constructing something. A moment later $e holds up a completed Vial of Formula 82!@n"), TRUE, ch, nil, nil, TO_ROOM)
				aid_prod = read_object(obj_vnum(num2), VIRTUAL)
				add_unique_id(aid_prod)
				obj_to_char(aid_prod, ch)
				extract_obj(aid_obj)
				improve_skill(ch, SKILL_FIRST_AID, 1)
			}
		}
	}
	WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
}
func do_aura(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: aura light\r\n"))
		return
	}
	if ch.Charge != 0 {
		send_to_char(ch, libc.CString("You can't focus enough on this while charging."))
		return
	}
	if PLR_FLAGGED(ch, PLR_POWERUP) {
		send_to_char(ch, libc.CString("You are busy powering up!\r\n"))
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("light")) == 0 {
		if GET_SKILL(ch, SKILL_FOCUS) != 0 && GET_SKILL(ch, SKILL_CONCENTRATION) < 75 {
			send_to_char(ch, libc.CString("You need at least a skill level of 75 in Focus and Concentration to use this.\r\n"))
			return
		} else {
			if PLR_FLAGGED(ch, PLR_AURALIGHT) {
				send_to_char(ch, libc.CString("Your aura fades as you stop shining light.\r\n"))
				act(libc.CString("$n's aura fades as they stop shining light on the area."), TRUE, ch, nil, nil, TO_ROOM)
				ch.Act[int(PLR_AURALIGHT/32)] &= bitvector_t(^(1 << (int(PLR_AURALIGHT % 32))))
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Light--
			} else if float64(ch.Mana) > float64(ch.Max_mana)*0.12 {
				if ch.In_room != room_rnum(-1) {
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Light++
				}
				reveal_hiding(ch, 0)
				ch.Mana -= int64(float64(ch.Max_mana) * 0.12)
				send_to_char(ch, libc.CString("A bright %s aura begins to burn around you as you provide light to the surrounding area!\r\n"), aura_types[ch.Aura])
				var bloom [2048]byte
				stdio.Sprintf(&bloom[0], "@wA %s aura flashes up brightly around $n@w as they provide light to the area.@n", aura_types[ch.Aura])
				act(&bloom[0], TRUE, ch, nil, nil, TO_ROOM)
				ch.Act[int(PLR_AURALIGHT/32)] |= bitvector_t(1 << (int(PLR_AURALIGHT % 32)))
			} else {
				send_to_char(ch, libc.CString("You don't have enough KI to do that.\r\n"))
				return
			}
		}
	}
}
