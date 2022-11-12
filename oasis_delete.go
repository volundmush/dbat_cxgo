package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func free_strings(data unsafe.Pointer, type_ int) int {
	var config *config_data
	switch type_ {
	case OASIS_MOB:
		fallthrough
	case OASIS_OBJ:
		return FALSE
	case OASIS_CFG:
		config = (*config_data)(data)
		if config.Play.OK != nil {
			libc.Free(unsafe.Pointer(config.Play.OK))
		}
		if config.Play.NOPERSON != nil {
			libc.Free(unsafe.Pointer(config.Play.NOPERSON))
		}
		if config.Play.NOEFFECT != nil {
			libc.Free(unsafe.Pointer(config.Play.NOEFFECT))
		}
		if config.Operation.DFLT_IP != nil {
			libc.Free(unsafe.Pointer(config.Operation.DFLT_IP))
		}
		if config.Operation.DFLT_DIR != nil {
			libc.Free(unsafe.Pointer(config.Operation.DFLT_DIR))
		}
		if config.Operation.LOGNAME != nil {
			libc.Free(unsafe.Pointer(config.Operation.LOGNAME))
		}
		if config.Operation.MENU != nil {
			libc.Free(unsafe.Pointer(config.Operation.MENU))
		}
		if config.Operation.WELC_MESSG != nil {
			libc.Free(unsafe.Pointer(config.Operation.WELC_MESSG))
		}
		if config.Operation.START_MESSG != nil {
			libc.Free(unsafe.Pointer(config.Operation.START_MESSG))
		}
		return TRUE
	default:
		mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("SYSERR: oasis_delete.c: free_strings: Invalid type handled (Type %d)."), type_)
		return FALSE
	}
}
