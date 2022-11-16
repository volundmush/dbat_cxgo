package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const ASSM_MAKE = 0
const ASSM_BAKE = 1
const ASSM_BREW = 2
const ASSM_ASSEMBLE = 3
const ASSM_CRAFT = 4
const ASSM_FLETCH = 5
const ASSM_KNIT = 6
const ASSM_MIX = 7
const ASSM_THATCH = 8
const ASSM_WEAVE = 9
const ASSM_FORGE = 10

type assembly_data struct {
	LVnum           int
	LNumComponents  int
	UchAssemblyType uint8
	PComponents     []component_data
}

type component_data struct {
	BExtract bool
	BInRoom  bool
	LVnum    int
}

var g_lNumAssemblies int = 0
var g_pAssemblyTable *assembly_data = nil

func assemblyBootAssemblies() {
	var (
		szLine     [64936]byte = [64936]byte{0: '\x00'}
		szTag      [64936]byte = [64936]byte{0: '\x00'}
		szType     [64936]byte = [64936]byte{0: '\x00'}
		iExtract   int         = 0
		iInRoom    int         = 0
		iType      int         = 0
		lLineCount int         = 0
		lPartVnum  int         = int(-1)
		lVnum      int         = int(-1)
		pFile      *stdio.File = nil
	)
	if (func() *stdio.File {
		pFile = stdio.FOpen(LIB_ETC, "rt")
		return pFile
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyBootAssemblies(): Couldn't open file '%s' for reading."), LIB_ETC)
		return
	}
	for int(pFile.IsEOF()) == 0 {
		lLineCount += get_line(pFile, &szLine[0])
		half_chop(&szLine[0], &szTag[0], &szLine[0])
		if szTag[0] == '\x00' {
			continue
		}
		if libc.StrCaseCmp(&szTag[0], libc.CString("Component")) == 0 {
			if stdio.Sscanf(&szLine[0], "#%ld %d %d", &lPartVnum, &iExtract, &iInRoom) != 3 {
				basic_mud_log(libc.CString("SYSERR: bootAssemblies(): Invalid format in file %s, line %ld: szTag=%s, szLine=%s."), LIB_ETC, lLineCount, &szTag[0], &szLine[0])
			} else if !assemblyAddComponent(lVnum, lPartVnum, iExtract != 0, iInRoom != 0) {
				basic_mud_log(libc.CString("SYSERR: bootAssemblies(): Could not add component #%ld to assembly #%ld."), lPartVnum, lVnum)
			}
		} else if libc.StrCaseCmp(&szTag[0], libc.CString("Vnum")) == 0 {
			if stdio.Sscanf(&szLine[0], "#%ld %s", &lVnum, &szType[0]) != 2 {
				basic_mud_log(libc.CString("SYSERR: bootAssemblies(): Invalid format in file %s, line %ld."), LIB_ETC, lLineCount)
				lVnum = -1
			} else if (func() int {
				iType = search_block(&szType[0], &AssemblyTypes[0], TRUE)
				return iType
			}()) < 0 {
				basic_mud_log(libc.CString("SYSERR: bootAssemblies(): Invalid type '%s' for assembly vnum #%ld at line %ld."), &szType[0], lVnum, lLineCount)
				lVnum = -1
			} else if !assemblyCreate(lVnum, iType) {
				basic_mud_log(libc.CString("SYSERR: bootAssemblies(): Could not create assembly for vnum #%ld, type %s."), lVnum, &szType[0])
				lVnum = -1
			}
		} else {
			basic_mud_log(libc.CString("SYSERR: Invalid tag '%s' in file %s, line #%ld."), &szTag[0], LIB_ETC, lLineCount)
		}
		szLine[0] = '\x00'
		szTag[0] = '\x00'
	}
	pFile.Close()
}
func assemblySaveAssemblies() {
	var (
		szType    [64936]byte    = [64936]byte{0: '\x00'}
		i         int            = 0
		j         int            = 0
		pAssembly *assembly_data = nil
		pFile     *stdio.File    = nil
	)
	if (func() *stdio.File {
		pFile = stdio.FOpen(LIB_ETC, "wt")
		return pFile
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblySaveAssemblies(): Couldn't open file '%s' for writing."), LIB_ETC)
		return
	}
	for i = 0; i < g_lNumAssemblies; i++ {
		pAssembly = (*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))
		sprinttype(int(pAssembly.UchAssemblyType), AssemblyTypes[:], &szType[0], uint64(64936))
		stdio.Fprintf(pFile, "Vnum                #%ld %s\n", pAssembly.LVnum, &szType[0])
		for j = 0; j < pAssembly.LNumComponents; j++ {
			stdio.Fprintf(pFile, "Component           #%ld %d %d\n", pAssembly.PComponents[j].LVnum, func() int {
				if pAssembly.PComponents[j].BExtract {
					return 1
				}
				return 0
			}(), func() int {
				if pAssembly.PComponents[j].BInRoom {
					return 1
				}
				return 0
			}())
		}
		if i < g_lNumAssemblies-1 {
			stdio.Fprintf(pFile, "\n")
		}
	}
	pFile.Close()
}
func assemblyListToChar(pCharacter *char_data) {
	var (
		szBuffer   [64936]byte = [64936]byte{0: '\x00'}
		szAssmType [2048]byte  = [2048]byte{0: '\x00'}
		i          int         = 0
		j          int         = 0
		lRnum      int         = 0
	)
	if pCharacter == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyListAssembliesToChar(): NULL 'pCharacter'."))
		return
	} else if g_pAssemblyTable == nil {
		send_to_char(pCharacter, libc.CString("No assemblies exist.\r\n"))
		return
	}
	send_to_char(pCharacter, libc.CString("The following assemblies exists:\r\n"))
	for i = 0; i < g_lNumAssemblies; i++ {
		if (func() int {
			lRnum = int(real_object(obj_vnum((*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum)))
			return lRnum
		}()) < 0 {
			send_to_char(pCharacter, libc.CString("[-----] ***RESERVED***\r\n"))
			basic_mud_log(libc.CString("SYSERR: assemblyListToChar(): Invalid vnum #%ld in assembly table."), (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum)
		} else {
			sprinttype(int((*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).UchAssemblyType), AssemblyTypes[:], &szAssmType[0], uint64(2048))
			stdio.Sprintf(&szBuffer[0], "[%5ld] %s (%s)\r\n", (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum, obj_proto[lRnum].Short_description, &szAssmType[0])
			send_to_char(pCharacter, &szBuffer[0])
			for j = 0; j < (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LNumComponents; j++ {
				if (func() int {
					lRnum = int(real_object(obj_vnum((*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[j].LVnum)))
					return lRnum
				}()) < 0 {
					send_to_char(pCharacter, libc.CString(" -----: ***RESERVED***\r\n"))
					basic_mud_log(libc.CString("SYSERR: assemblyListToChar(): Invalid component vnum #%ld in assembly for vnum #%ld."), (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[j].LVnum, (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum)
				} else {
					stdio.Sprintf(&szBuffer[0], " %5ld: %-20.20s Extract=%-3.3s InRoom=%-3.3s\r\n", +(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[j].LVnum, obj_proto[lRnum].Short_description, func() string {
						if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[j].BExtract {
							return "Yes"
						}
						return "No"
					}(), func() string {
						if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[j].BInRoom {
							return "Yes"
						}
						return "No"
					}())
					send_to_char(pCharacter, &szBuffer[0])
				}
			}
		}
	}
}
func assemblyAddComponent(lVnum int, lComponentVnum int, bExtract bool, bInRoom bool) bool {
	var pAssembly *assembly_data = nil
	if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyAddComponent(): Invalid 'lVnum' #%ld."), lVnum)
		return FALSE != 0
	} else if real_object(obj_vnum(lComponentVnum)) <= obj_rnum(-1) {
		basic_mud_log(libc.CString("SYSERR: assemblyAddComponent(): Invalid 'lComponentVnum' #%ld."), lComponentVnum)
		return FALSE != 0
	}
	if pAssembly.PComponents == nil {
		pAssembly.PComponents = make([]component_data, pAssembly.LNumComponents+1)
	} else {
		//pAssembly.PComponents = []component_data((*component_data)(libc.Realloc(unsafe.Pointer(&pAssembly.PComponents[0]), pAssembly.LNumComponents*int(unsafe.Sizeof(component_data{}))+1)))
	}
	pAssembly.PComponents[pAssembly.LNumComponents].LVnum = lComponentVnum
	pAssembly.PComponents[pAssembly.LNumComponents].BExtract = bExtract
	pAssembly.PComponents[pAssembly.LNumComponents].BInRoom = bInRoom
	pAssembly.LNumComponents += 1
	return TRUE != 0
}
func assemblyCheckComponents(lVnum int, pCharacter *char_data, extract_yes int) bool {
	var (
		bOk                bool           = TRUE != 0
		i                  int            = 0
		lRnum              int            = 0
		ppComponentObjects **obj_data     = nil
		pAssembly          *assembly_data = nil
	)
	if pCharacter == nil {
		basic_mud_log(libc.CString("SYSERR: NULL assemblyCheckComponents(): 'pCharacter'."))
		return FALSE != 0
	} else if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: NULL assemblyCheckComponents(): Invalid 'lVnum' #%ld."), lVnum)
		return FALSE != 0
	}
	if pAssembly.PComponents == nil {
		return FALSE != 0
	} else if pAssembly.LNumComponents <= 0 {
		return FALSE != 0
	}
	ppComponentObjects = &make([]*obj_data, pAssembly.LNumComponents)[0]
	for i = 0; i < pAssembly.LNumComponents && bOk; i++ {
		if (func() int {
			lRnum = int(real_object(obj_vnum(pAssembly.PComponents[i].LVnum)))
			return lRnum
		}()) < 0 {
			bOk = FALSE != 0
		} else {
			if pAssembly.PComponents[i].BInRoom {
				if (func() *obj_data {
					p := (**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i)))
					*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))) = get_obj_in_list_num(lRnum, world[pCharacter.In_room].Contents)
					return *p
				}()) == nil {
					bOk = FALSE != 0
				} else {
					obj_from_room(*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))))
				}
			} else {
				if (func() *obj_data {
					p := (**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i)))
					*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))) = get_obj_in_list_num(lRnum, pCharacter.Carrying)
					return *p
				}()) == nil {
					bOk = FALSE != 0
				} else {
					obj_from_char(*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))))
				}
			}
		}
	}
	for i = 0; i < pAssembly.LNumComponents; i++ {
		if *(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))) == nil {
			continue
		}
		if pAssembly.PComponents[i].BExtract && bOk && extract_yes == TRUE {
			extract_obj(*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))))
		} else if pAssembly.PComponents[i].BInRoom {
			obj_to_room(*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))), pCharacter.In_room)
		} else {
			obj_to_char(*(**obj_data)(unsafe.Add(unsafe.Pointer(ppComponentObjects), unsafe.Sizeof((*obj_data)(nil))*uintptr(i))), pCharacter)
		}
	}
	libc.Free(unsafe.Pointer(ppComponentObjects))
	return bOk
}
func assemblyCreate(lVnum int, iAssembledType int) bool {
	var (
		lBottom           int            = 0
		lMiddle           int            = 0
		lTop              int            = 0
		pNewAssemblyTable *assembly_data = nil
	)
	if lVnum < 0 {
		return FALSE != 0
	} else if iAssembledType < 0 || iAssembledType >= MAX_ASSM {
		return FALSE != 0
	}
	if g_pAssemblyTable == nil {
		g_pAssemblyTable = new(assembly_data)
		g_lNumAssemblies = 1
	} else {
		lTop = g_lNumAssemblies - 1
		for {
			lMiddle = (lBottom + lTop) / 2
			if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum == lVnum {
				return FALSE != 0
			} else if lBottom >= lTop {
				break
			} else if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum > lVnum {
				lTop = lMiddle - 1
			} else {
				lBottom = lMiddle + 1
			}
		}
		if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum <= lVnum {
			lMiddle += 1
		}
		pNewAssemblyTable = &make([]assembly_data, g_lNumAssemblies+1)[0]
		if lMiddle > 0 {
			libc.MemMove(unsafe.Pointer(pNewAssemblyTable), unsafe.Pointer(g_pAssemblyTable), lMiddle*int(unsafe.Sizeof(assembly_data{})))
		}
		if lMiddle <= g_lNumAssemblies-1 {
			libc.MemMove(unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer(pNewAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))), unsafe.Sizeof(assembly_data{})*1))), unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))), (g_lNumAssemblies-lMiddle)*int(unsafe.Sizeof(assembly_data{})))
		}
		libc.Free(unsafe.Pointer(g_pAssemblyTable))
		g_pAssemblyTable = pNewAssemblyTable
		g_lNumAssemblies += 1
	}
	(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LNumComponents = 0
	(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum = lVnum
	(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).PComponents = nil
	(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).UchAssemblyType = uint8(int8(iAssembledType))
	return TRUE != 0
}
func assemblyDestroy(lVnum int) bool {
	var (
		lIndex            int            = 0
		pNewAssemblyTable *assembly_data = nil
	)
	if g_pAssemblyTable == nil || (func() int {
		lIndex = assemblyGetAssemblyIndex(lVnum)
		return lIndex
	}()) < 0 {
		basic_mud_log(libc.CString("SYSERR: assemblyDestroy(): Invalid 'lVnum' #%ld."), lVnum)
		return FALSE != 0
	}
	if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lIndex)))).PComponents != nil {
		libc.Free(unsafe.Pointer(&(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lIndex)))).PComponents[0]))
	}
	if g_lNumAssemblies > 1 {
		pNewAssemblyTable = &make([]assembly_data, g_lNumAssemblies-1)[0]
		if lIndex > 0 {
			libc.MemMove(unsafe.Pointer(pNewAssemblyTable), unsafe.Pointer(g_pAssemblyTable), lIndex*int(unsafe.Sizeof(assembly_data{})))
		}
		if lIndex < g_lNumAssemblies-1 {
			libc.MemMove(unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer(pNewAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lIndex)))), unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer((*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lIndex)))), unsafe.Sizeof(assembly_data{})*1))), (g_lNumAssemblies-lIndex-1)*int(unsafe.Sizeof(assembly_data{})))
		}
	}
	libc.Free(unsafe.Pointer(g_pAssemblyTable))
	g_lNumAssemblies -= 1
	g_pAssemblyTable = pNewAssemblyTable
	return TRUE != 0
}
func assemblyHasComponent(lVnum int, lComponentVnum int) bool {
	var pAssembly *assembly_data = nil
	if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyHasComponent(): Invalid 'lVnum' #%ld."), lVnum)
		return FALSE != 0
	}
	return assemblyGetComponentIndex(pAssembly, lComponentVnum) >= 0
}
func assemblyRemoveComponent(lVnum int, lComponentVnum int) bool {
	var (
		lIndex         int             = 0
		pAssembly      *assembly_data  = nil
		pNewComponents *component_data = nil
	)
	if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyRemoveComponent(): Invalid 'lVnum' #%ld."), lVnum)
		return FALSE != 0
	} else if (func() int {
		lIndex = assemblyGetComponentIndex(pAssembly, lComponentVnum)
		return lIndex
	}()) < 0 {
		basic_mud_log(libc.CString("SYSERR: assemblyRemoveComponent(): Vnum #%ld is not a component of assembled vnum #%ld."), lComponentVnum, lVnum)
		return FALSE != 0
	}
	if pAssembly.PComponents != nil && pAssembly.LNumComponents > 1 {
		pNewComponents = &make([]component_data, pAssembly.LNumComponents-1)[0]
		if lIndex > 0 {
			libc.MemMove(unsafe.Pointer(pNewComponents), unsafe.Pointer(&pAssembly.PComponents[0]), lIndex*int(unsafe.Sizeof(component_data{})))
		}
		if lIndex < pAssembly.LNumComponents-1 {
			libc.MemMove(unsafe.Pointer((*component_data)(unsafe.Add(unsafe.Pointer(pNewComponents), unsafe.Sizeof(component_data{})*uintptr(lIndex)))), unsafe.Pointer(&pAssembly.PComponents[lIndex+1]), (pAssembly.LNumComponents-lIndex-1)*int(unsafe.Sizeof(component_data{})))
		}
	}
	libc.Free(unsafe.Pointer(&pAssembly.PComponents[0]))
	//pAssembly.PComponents = []component_data(pNewComponents)
	pAssembly.LNumComponents -= 1
	return TRUE != 0
}
func assemblyGetType(lVnum int) int {
	var pAssembly *assembly_data = nil
	if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyGetType(): Invalid 'lVnum' #%ld."), lVnum)
		return -1
	}
	return int(pAssembly.UchAssemblyType)
}
func assemblyCountComponents(lVnum int) int {
	var pAssembly *assembly_data = nil
	if (func() *assembly_data {
		pAssembly = assemblyGetAssemblyPtr(lVnum)
		return pAssembly
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: assemblyCountComponents(): Invalid 'lVnum' #%ld."), lVnum)
		return 0
	}
	return pAssembly.LNumComponents
}
func assemblyFindAssembly(pszAssemblyName *byte) int {
	var (
		i     int = 0
		lRnum int = int(-1)
	)
	if g_pAssemblyTable == nil {
		return -1
	} else if pszAssemblyName == nil || *pszAssemblyName == '\x00' {
		return -1
	}
	for i = 0; i < g_lNumAssemblies; i++ {
		if (func() int {
			lRnum = int(real_object(obj_vnum((*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum)))
			return lRnum
		}()) < 0 {
			basic_mud_log(libc.CString("SYSERR: assemblyFindAssembly(): Invalid vnum #%ld in assembly table."), (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum)
		} else if isname(pszAssemblyName, obj_proto[lRnum].Name) != 0 {
			return (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LVnum
		}
	}
	return -1
}
func assemblyGetAssemblyIndex(lVnum int) int {
	var (
		lBottom int = 0
		lMiddle int = 0
		lTop    int = 0
	)
	lTop = g_lNumAssemblies - 1
	for {
		lMiddle = (lBottom + lTop) / 2
		if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum == lVnum {
			return lMiddle
		} else if lBottom >= lTop {
			return -1
		} else if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lMiddle)))).LVnum > lVnum {
			lTop = lMiddle - 1
		} else {
			lBottom = lMiddle + 1
		}
	}
}
func assemblyGetComponentIndex(pAssembly *assembly_data, lComponentVnum int) int {
	var i int = 0
	if pAssembly == nil {
		return -1
	}
	for i = 0; i < pAssembly.LNumComponents; i++ {
		if pAssembly.PComponents[i].LVnum == lComponentVnum {
			return i
		}
	}
	return -1
}
func assemblyGetAssemblyPtr(lVnum int) *assembly_data {
	var lIndex int = 0
	if g_pAssemblyTable == nil {
		return nil
	}
	if (func() int {
		lIndex = assemblyGetAssemblyIndex(lVnum)
		return lIndex
	}()) >= 0 {
		return (*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(lIndex)))
	}
	return nil
}
func free_assemblies() {
	var i int
	if g_pAssemblyTable == nil {
		return
	}
	for i = 0; i < g_lNumAssemblies; i++ {
		if (*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents != nil {
			libc.Free(unsafe.Pointer(&(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).PComponents[0]))
		}
		(*(*assembly_data)(unsafe.Add(unsafe.Pointer(g_pAssemblyTable), unsafe.Sizeof(assembly_data{})*uintptr(i)))).LNumComponents = 0
	}
	libc.Free(unsafe.Pointer(g_pAssemblyTable))
}
