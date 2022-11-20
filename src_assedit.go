package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
	"unsafe"
)

var lRnum int = 0

func do_assedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d    *descriptor_data = ch.Desc
		buf  [64936]byte
		buf2 [64936]byte
	)
	buf[0] = '\x00'
	buf2[0] = '\x00'
	if IS_NPC(ch) {
		return
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_ASSEDIT {
			send_to_char(ch, libc.CString("Assemblies are already being editted by someone.\r\n"))
			return
		}
	}
	two_arguments(argument, &buf[0], &buf2[0])
	d = ch.Desc
	if buf[0] == 0 {
		nodigit(d)
		return
	}
	if !unicode.IsDigit(rune(buf[0])) {
		if libc.StrNCaseCmp(libc.CString("new"), &buf[0], 3) == 0 {
			if !unicode.IsDigit(rune(buf2[0])) {
				nodigit(d)
			} else if real_object(libc.Atoi(libc.GoString(&buf2[0]))) == int(-1) {
				send_to_char(d.Character, libc.CString("You need to create the assembly object before you can create the new assembly.\r\n"))
				return
			} else {
				assemblyCreate(libc.Atoi(libc.GoString(&buf2[0])), 0)
				send_to_char(d.Character, libc.CString("Assembly Created.\r\n"))
				assemblySaveAssemblies()
				return
			}
		} else if libc.StrNCaseCmp(libc.CString("delete"), &buf[0], 6) == 0 {
			if !unicode.IsDigit(rune(buf2[0])) {
				nodigit(d)
			} else {
				assemblyDestroy(libc.Atoi(libc.GoString(&buf2[0])))
				send_to_char(d.Character, libc.CString("Assembly Deleted.\r\n"))
				assemblySaveAssemblies()
				return
			}
		} else {
			nodigit(d)
			return
		}
	} else if unicode.IsDigit(rune(buf[0])) {
		d = ch.Desc
		d.Olc = new(oasis_olc_data)
		assedit_setup(d, libc.Atoi(libc.GoString(&buf[0])))
	}
}
func assedit_setup(d *descriptor_data, number int) {
	var pOldAssembly *assembly_data = nil
	// todo: fix this
	//d.Olc.OlcAssembly = []assembly_data(new(assembly_data))
	if (func() *assembly_data {
		pOldAssembly = assemblyGetAssemblyPtr(number)
		return pOldAssembly
	}()) == nil {
		send_to_char(d.Character, libc.CString("That assembly does not exist\r\n"))
		cleanup_olc(d, CLEANUP_ALL)
		return
	} else {
		d.Olc.OlcAssembly[0].LVnum = pOldAssembly.LVnum
		d.Olc.OlcAssembly[0].UchAssemblyType = pOldAssembly.UchAssemblyType
		d.Olc.OlcAssembly[0].LNumComponents = pOldAssembly.LNumComponents
		if d.Olc.OlcAssembly[0].LNumComponents > 0 {
			d.Olc.OlcAssembly[0].PComponents = make([]component_data, d.Olc.OlcAssembly[0].LNumComponents)
			libc.MemMove(unsafe.Pointer(&d.Olc.OlcAssembly[0].PComponents[0]), unsafe.Pointer(&pOldAssembly.PComponents[0]), d.Olc.OlcAssembly[0].LNumComponents*int(unsafe.Sizeof(component_data{})))
		}
	}
	if (func() int {
		lRnum = real_object(d.Olc.OlcAssembly[0].LVnum)
		return lRnum
	}()) < 0 {
		send_to_char(d.Character, libc.CString("Assembled item may not exist, check the vnum and assembles (show assemblies). \r\n"))
		cleanup_olc(d, CLEANUP_ALL)
		return
	}
	d.Connected = CON_ASSEDIT
	act(libc.CString("$n starts using OLC."), 1, d.Character, nil, nil, TO_ROOM)
	SET_BIT_AR(d.Character.Act[:], PLR_WRITING)
	assedit_disp_menu(d)
}
func assedit_disp_menu(d *descriptor_data) {
	var (
		i          int        = 0
		szAssmType [2048]byte = [2048]byte{0: '\x00'}
	)
	sprinttype(int(d.Olc.OlcAssembly[0].UchAssemblyType), AssemblyTypes[:], &szAssmType[0], uint64(2048))
	send_to_char(d.Character, libc.CString("Assembly Number: @c%ld@n\r\nAssembly Name  : @y%s@n\r\nAssembly Type  : @y%s@n\r\nComponents:\r\n"), d.Olc.OlcAssembly[0].LVnum, obj_proto[real_object(d.Olc.OlcAssembly[0].LVnum)].Short_description, &szAssmType[0])
	if d.Olc.OlcAssembly[0].LNumComponents <= 0 {
		send_to_char(d.Character, libc.CString("   < NONE > \r\n"))
	} else {
		for i = 0; i < d.Olc.OlcAssembly[0].LNumComponents; i++ {
			if (func() int {
				lRnum = real_object(d.Olc.OlcAssembly[0].PComponents[i].LVnum)
				return lRnum
			}()) < 0 {
				send_to_char(d.Character, libc.CString("@g%2d@n) @y ERROR --- Contact an Implementor @n\r\n "), i+1)
			} else {
				send_to_char(d.Character, libc.CString("@g%2d@n) [@c%5ld@n] %-20.20s  In room: @c%-3.3s@n    Extract: @y%-3.3s@n\r\n"), i+1, d.Olc.OlcAssembly[0].PComponents[i].LVnum, obj_proto[lRnum].Short_description, func() string {
					if d.Olc.OlcAssembly[0].PComponents[i].BInRoom {
						return "Yes"
					}
					return "No"
				}(), func() string {
					if d.Olc.OlcAssembly[0].PComponents[i].BExtract {
						return "Yes"
					}
					return "No"
				}())
			}
		}
	}
	send_to_char(d.Character, libc.CString("@gA@n) Add a new component.\r\n@gE@n) Edit a component.\r\n@gD@n) Delete a component.\r\n@gT@n) Change Assembly Type.\r\n@gQ@n) Quit.\r\n\r\nEnter your choice : "))
	d.Olc.Mode = ASSEDIT_MAIN_MENU
	return
}
func assedit_parse(d *descriptor_data, arg *byte) {
	var (
		pos          int = 0
		i            int = 0
		counter      int
		columns      int             = 0
		pTComponents *component_data = nil
	)
	switch d.Olc.Mode {
	case ASSEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			assemblyDestroy(d.Olc.OlcAssembly[0].LVnum)
			assemblyCreate(d.Olc.OlcAssembly[0].LVnum, int(d.Olc.OlcAssembly[0].UchAssemblyType))
			for i = 0; i < d.Olc.OlcAssembly[0].LNumComponents; i++ {
				assemblyAddComponent(d.Olc.OlcAssembly[0].LVnum, d.Olc.OlcAssembly[0].PComponents[i].LVnum, d.Olc.OlcAssembly[0].PComponents[i].BExtract, d.Olc.OlcAssembly[0].PComponents[i].BInRoom)
			}
			send_to_char(d.Character, libc.CString("\r\nSaving all assemblies\r\n"))
			assemblySaveAssemblies()
			cleanup_olc(d, CLEANUP_ALL)
		case 't':
			fallthrough
		case 'T':
			for counter = 0; counter < MAX_ASSM; counter++ {
				send_to_char(d.Character, libc.CString("@g%2d@n) %-20.20s %s"), counter+1, AssemblyTypes[counter], func() string {
					if (func() int {
						p := &columns
						*p++
						return *p
					}() % 2) == 0 {
						return "\r\n"
					}
					return ""
				}())
			}
			send_to_char(d.Character, libc.CString("Enter the assembly type : "))
			d.Olc.Mode = ASSEDIT_EDIT_TYPES
		case 'a':
			fallthrough
		case 'A':
			send_to_char(d.Character, libc.CString("\r\nWhat is the vnum of the new component?"))
			d.Olc.Mode = ASSEDIT_ADD_COMPONENT
		case 'e':
			fallthrough
		case 'E':
			send_to_char(d.Character, libc.CString("\r\nEdit which component? "))
			d.Olc.Mode = ASSEDIT_EDIT_COMPONENT
		case 'd':
			fallthrough
		case 'D':
			if pos < 0 || pos > d.Olc.OlcAssembly[0].LNumComponents {
				send_to_char(d.Character, libc.CString("\r\nWhich component do you wish to remove?"))
				assedit_disp_menu(d)
			} else {
				send_to_char(d.Character, libc.CString("\r\nWhich component do you wish to remove?"))
				d.Olc.Mode = ASSEDIT_DELETE_COMPONENT
			}
		default:
			assedit_disp_menu(d)
		}
	case ASSEDIT_EDIT_TYPES:
		if unicode.IsDigit(rune(*arg)) {
			pos = libc.Atoi(libc.GoString(arg)) - 1
			if pos >= 0 || pos < MAX_ASSM {
				d.Olc.OlcAssembly[0].UchAssemblyType = uint8(int8(pos))
				assedit_disp_menu(d)
				break
			}
		} else {
			assedit_disp_menu(d)
		}
	case ASSEDIT_ADD_COMPONENT:
		if unicode.IsDigit(rune(*arg)) {
			pos = libc.Atoi(libc.GoString(arg))
			if real_object(pos) <= int(-1) {
				break
			}
			for i = 0; i < d.Olc.OlcAssembly[0].LNumComponents; i++ {
				if d.Olc.OlcAssembly[0].PComponents[i].LVnum == pos {
					break
				}
			}
			pTComponents = &make([]component_data, d.Olc.OlcAssembly[0].LNumComponents+1)[0]
			if d.Olc.OlcAssembly[0].PComponents != nil {
				libc.MemMove(unsafe.Pointer(pTComponents), unsafe.Pointer(&d.Olc.OlcAssembly[0].PComponents[0]), d.Olc.OlcAssembly[0].LNumComponents*int(unsafe.Sizeof(component_data{})))
			}
			// todo: fix this
			//d.Olc.OlcAssembly[0].PComponents = []component_data(pTComponents)
			d.Olc.OlcAssembly[0].PComponents[d.Olc.OlcAssembly[0].LNumComponents].LVnum = pos
			d.Olc.OlcAssembly[0].PComponents[d.Olc.OlcAssembly[0].LNumComponents].BExtract = YES != 0
			d.Olc.OlcAssembly[0].PComponents[d.Olc.OlcAssembly[0].LNumComponents].BInRoom = NO != 0
			d.Olc.OlcAssembly[0].LNumComponents += 1
			assedit_disp_menu(d)
		} else {
			send_to_char(d.Character, libc.CString("That object does not exist. Please try again\r\n"))
			assedit_disp_menu(d)
		}
	case ASSEDIT_EDIT_COMPONENT:
		pos = libc.Atoi(libc.GoString(arg))
		if unicode.IsDigit(rune(*arg)) {
			pos--
			d.Olc.Value = pos
			assedit_edit_extract(d)
			break
		} else {
			assedit_disp_menu(d)
		}
	case ASSEDIT_DELETE_COMPONENT:
		if unicode.IsDigit(rune(*arg)) {
			pos = libc.Atoi(libc.GoString(arg))
			pos -= 1
			pTComponents = &make([]component_data, d.Olc.OlcAssembly[0].LNumComponents-1)[0]
			if pos > 0 {
				libc.MemMove(unsafe.Pointer(pTComponents), unsafe.Pointer(&d.Olc.OlcAssembly[0].PComponents[0]), pos*int(unsafe.Sizeof(component_data{})))
			}
			if pos < d.Olc.OlcAssembly[0].LNumComponents-1 {
				libc.MemMove(unsafe.Pointer((*component_data)(unsafe.Add(unsafe.Pointer(pTComponents), unsafe.Sizeof(component_data{})*uintptr(pos)))), unsafe.Pointer(&d.Olc.OlcAssembly[0].PComponents[pos+1]), (d.Olc.OlcAssembly[0].LNumComponents-pos-1)*int(unsafe.Sizeof(component_data{})))
			}
			libc.Free(unsafe.Pointer(&d.Olc.OlcAssembly[0].PComponents[0]))
			// todo: fix this
			//d.Olc.OlcAssembly[0].PComponents = []component_data(pTComponents)
			d.Olc.OlcAssembly[0].LNumComponents -= 1
			assedit_disp_menu(d)
			break
		} else {
			assedit_disp_menu(d)
		}
	case ASSEDIT_EDIT_EXTRACT:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.OlcAssembly[0].PComponents[d.Olc.Value].BExtract = true
			assedit_edit_inroom(d)
		case 'n':
			fallthrough
		case 'N':
			d.Olc.OlcAssembly[0].PComponents[d.Olc.Value].BExtract = false
			assedit_edit_inroom(d)
		default:
			send_to_char(d.Character, libc.CString("Is the item to be extracted when the assembly is created? (Y/N)"))
		}
	case ASSEDIT_EDIT_INROOM:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.OlcAssembly[0].PComponents[d.Olc.Value].BInRoom = true
			assedit_disp_menu(d)
		case 'n':
			fallthrough
		case 'N':
			d.Olc.OlcAssembly[0].PComponents[d.Olc.Value].BInRoom = false
			assedit_disp_menu(d)
		default:
			send_to_char(d.Character, libc.CString("Object in the room when assembly is created? (n =  in inventory):"))
		}
	default:
		mudlog(BRF, ADMLVL_GOD, 1, libc.CString("SYSERR: OLC assedit_parse(): Reached default case!"))
		send_to_char(d.Character, libc.CString("Opps...\r\n"))
		d.Connected = CON_PLAYING
	}
}
func assedit_delete(d *descriptor_data) {
	send_to_char(d.Character, libc.CString("Which item number do you wish to delete from this assembly?"))
	d.Olc.Mode = ASSEDIT_DELETE_COMPONENT
	return
}
func assedit_edit_extract(d *descriptor_data) {
	send_to_char(d.Character, libc.CString("Is the item to be extracted when the assembly is created? (Y/N):"))
	d.Olc.Mode = ASSEDIT_EDIT_EXTRACT
	return
}
func assedit_edit_inroom(d *descriptor_data) {
	send_to_char(d.Character, libc.CString("Should the object be in the room when assembly is created (n = in inventory)?"))
	d.Olc.Mode = ASSEDIT_EDIT_INROOM
	return
}
func nodigit(d *descriptor_data) {
	send_to_char(d.Character, libc.CString("Usage: assedit <vnum>\r\n"))
	send_to_char(d.Character, libc.CString("     : assedit new <vnum>\r\n"))
	send_to_char(d.Character, libc.CString("     : assedit delete <vnum>\r\n"))
	return
}
