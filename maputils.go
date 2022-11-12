package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const MAP_ROWS = 199
const MAP_COLS = 199

type mapstruct struct {
	X int
	Y int
}
type MapStruct mapstruct

var mapnums [200][200]int

func ping_ship(vnum int, vnum2 int) {
	var (
		tch      *char_data
		next_ch  *char_data
		controls *obj_data = nil
		obj      *obj_data = nil
		found    int       = FALSE
	)
	if vnum2 == -1 {
		return
	}
	for tch = character_list; tch != nil; tch = next_ch {
		next_ch = tch.Next
		if found == FALSE {
			if (func() *obj_data {
				obj = find_control(tch)
				return obj
			}()) == nil {
				continue
			} else {
				if (obj.Value[0]) == vnum && vnum != vnum2 {
					controls = obj
					found = TRUE
				}
			}
		}
	}
	if found == TRUE {
		send_to_room(controls.In_room, libc.CString("@D[@RALERT@D: @YAn unknown radar signal has been detected!@D]@n"))
	}
}
func checkship(rnum int, vnum int) int {
	var (
		i     *obj_data = nil
		there int       = FALSE
	)
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Contents; i != nil; i = i.Next_content {
		if !ROOM_FLAGGED(room_rnum(rnum), ROOM_NEBULA) {
			if i.Type_flag == ITEM_VEHICLE && there != TRUE {
				there = TRUE
				ping_ship(int(GET_OBJ_VNUM(i)), vnum)
			}
		}
	}
	i = nil
	return there
}
func getmapchar(rnum int, ch *char_data, start int, vnum int) *byte {
	var (
		mapchar [50]byte
		there   int = FALSE
		enemy   int = FALSE
	)
	if rnum == start {
		there = TRUE
	}
	if checkship(rnum, vnum) != 0 {
		enemy = TRUE
	}
	if rnum == int(real_room(ch.Radar1)) || rnum == int(real_room(ch.Radar2)) || rnum == int(real_room(ch.Radar3)) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@WB@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@WB@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@WBB")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_EORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@GE@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@GE@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@GEE")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_CORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@MC@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@MC@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@MCC")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_FORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@CF@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@CF@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@CFF")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_KORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@mK@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@mK@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@mKK")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_NORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@gN@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@gN@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@gNN")
		}
	} else if (func() room_vnum {
		if rnum != int(-1) && rnum <= int(top_of_world) {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Number
		}
		return -1
	}()) == 0xC654 {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@cZ@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@cZ@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@cZZ")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_VORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@YV@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@YV@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@YVV")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_AORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@BA@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@BA@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@BAA")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_YORBIT) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@MY@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@MY@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@MYY")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_KANORB) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@CK@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@CK@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@CKK")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_ARLORB) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@mA@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@mA@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@mAA")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_NEBULA) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@m&@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@m&@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@m&&")
		}
	} else if (func() room_vnum {
		if rnum != int(-1) && rnum <= int(top_of_world) {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Number
		}
		return -1
	}()) == 0x948C {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@yQ@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@yQ@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@yQQ")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_ASTERO) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@y:@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@y:@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@y::")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_WORMHO) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@b@1*@RX@n")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@b@1*@r#@n")
		} else {
			stdio.Sprintf(&mapchar[0], "@b@1**@n")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_STATION) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@DS@RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@DS@r#")
		} else {
			stdio.Sprintf(&mapchar[0], "@DSS")
		}
	} else if ROOM_FLAGGED(room_rnum(rnum), ROOM_STAR) {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@6 @RX@n")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@6 @r#@n")
		} else {
			stdio.Sprintf(&mapchar[0], "@6  @n")
		}
	} else {
		if there != 0 {
			stdio.Sprintf(&mapchar[0], "@w @RX")
		} else if enemy == TRUE {
			stdio.Sprintf(&mapchar[0], "@w @r#")
		} else {
			var color int = rand_number(1, 30)
			if rand_number(1, 40) == 2 {
				stdio.Sprintf(&mapchar[0], "%s. ", func() string {
					if color > 15 {
						return "@w"
					}
					if color >= 7 {
						return "@Y"
					}
					if color > 3 {
						return "@R"
					}
					return "@B"
				}())
			} else if rand_number(1, 40) == 2 {
				stdio.Sprintf(&mapchar[0], "%s .", func() string {
					if color > 15 {
						return "@w"
					}
					if color >= 7 {
						return "@Y"
					}
					if color > 3 {
						return "@R"
					}
					return "@B"
				}())
			} else {
				stdio.Sprintf(&mapchar[0], "@w  ")
			}
		}
	}
	return &mapchar[0]
}
func findcoord(rnum int) MapStruct {
	var (
		x      int
		y      int
		coords MapStruct
	)
	coords.X = 0
	coords.Y = 0
	for y = 0; y <= MAP_ROWS; y++ {
		for x = 0; x <= MAP_COLS; x++ {
			if mapnums[y][x] == rnum {
				coords.Y = y
				coords.X = x
				return coords
			}
		}
	}
	basic_mud_log(libc.CString("SYSERR: findcoord for non-map rnum"))
	return coords
}
func printmap(rnum int, ch *char_data, type_ int, vnum int) {
	var (
		x           int = 0
		lasty       int = -1
		y           int = 0
		sightradius int
		count       int = 0
		initline    int = 0
		buf         [129872]byte
		buf2        [512]byte
		coord       MapStruct
		start       int = rnum
	)
	coord = findcoord(rnum)
	C.strcpy(&buf[0], libc.CString("\n"))
	if type_ == 0 {
		sightradius = 12
	} else {
		sightradius = 4
	}
	if type_ == 0 {
		send_to_char(ch, libc.CString("@b______________________________________________________________________@n\n"))
	}
	for y = coord.Y - sightradius; y <= coord.Y+sightradius; y++ {
		if type_ == 0 {
			if count == initline {
				C.strcat(&buf[0], libc.CString("@b     [@CR. Key@b]     | "))
			} else if count == initline+1 {
				C.strcat(&buf[0], libc.CString("@GEE@D:@w Earth@b         | "))
			} else if count == initline+2 {
				C.strcat(&buf[0], libc.CString("@gNN@D:@w Namek@b         | "))
			} else if count == initline+3 {
				C.strcat(&buf[0], libc.CString("@YVV@D:@w Vegeta@b        | "))
			} else if count == initline+4 {
				C.strcat(&buf[0], libc.CString("@CFF@D:@w Frigid@b        | "))
			} else if count == initline+5 {
				C.strcat(&buf[0], libc.CString("@mKK@D:@w Konack@b        | "))
			} else if count == initline+6 {
				C.strcat(&buf[0], libc.CString("@BAA@D:@w Aether@b        | "))
			} else if count == initline+7 {
				C.strcat(&buf[0], libc.CString("@MYY@D:@w Yardrat@b       | "))
			} else if count == initline+8 {
				C.strcat(&buf[0], libc.CString("@CKK@D:@w Kanassa@b       | "))
			} else if count == initline+9 {
				C.strcat(&buf[0], libc.CString("@mAA@D:@w Arlia@b         | "))
			} else if count == initline+10 {
				C.strcat(&buf[0], libc.CString("@cZZ@D:@w Zenith@b        | "))
			} else if count == initline+11 {
				C.strcat(&buf[0], libc.CString("@MCC@D:@w Cerria@b        | "))
			} else if count == initline+12 {
				C.strcat(&buf[0], libc.CString("@WBB@D:@w Buoy@b          | "))
			} else if count == initline+13 {
				C.strcat(&buf[0], libc.CString("@m&&@D:@w Nebula@b        | "))
			} else if count == initline+14 {
				C.strcat(&buf[0], libc.CString("@yQQ@D:@w Asteroid@b      | "))
			} else if count == initline+15 {
				C.strcat(&buf[0], libc.CString("@y::@D:@w Asteroid Field@b| "))
			} else if count == initline+16 {
				C.strcat(&buf[0], libc.CString("@b@1**@n@D:@w Wormhole@b      | "))
			} else if count == initline+17 {
				C.strcat(&buf[0], libc.CString("@DSS@D:@w S. Station@b    | "))
			} else if count == initline+18 {
				C.strcat(&buf[0], libc.CString(" @r#@D:@w Unknown Ship@b  | "))
			} else if count == initline+19 {
				C.strcat(&buf[0], libc.CString("@6  @n@D:@w Star@b          | "))
			} else {
				C.strcat(&buf[0], libc.CString("                  @b| "))
			}
			count++
		} else {
			if count == 0 {
				C.strcat(&buf[0], libc.CString("      @RCompass@n           "))
			} else if count == 2 {
				stdio.Sprintf(&buf2[0], "@w       @w|%s@w|            ", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[0]) != nil {
						return " @CN "
					}
					return "   "
				}())
				C.strcat(&buf[0], &buf2[0])
			} else if count == 3 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|      ", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[6]) != nil {
						return " @CNW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[4]) != nil {
						return " @YU "
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[7]) != nil {
						return "@CNE "
					}
					return "   "
				}())
				C.strcat(&buf[0], &buf2[0])
			} else if count == 4 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|      ", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[3]) != nil {
						return "  @CW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[10]) != nil {
						return "@m I "
					}
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[11]) != nil {
						return "@mOUT"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[1]) != nil {
						return "@CE  "
					}
					return "   "
				}())
				C.strcat(&buf[0], &buf2[0])
			} else if count == 5 {
				stdio.Sprintf(&buf2[0], "@w @w|%s@w| |%s@w| |%s@w|      ", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[9]) != nil {
						return " @CSW"
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[5]) != nil {
						return " @YD "
					}
					return "   "
				}(), func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[8]) != nil {
						return "@CSE "
					}
					return "   "
				}())
				C.strcat(&buf[0], &buf2[0])
			} else if count == 6 {
				stdio.Sprintf(&buf2[0], "@w       @w|%s@w|            ", func() string {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Dir_option[2]) != nil {
						return " @CS "
					}
					return "   "
				}())
				C.strcat(&buf[0], &buf2[0])
			} else {
				C.strcat(&buf[0], libc.CString("                        "))
			}
			count++
		}
		for x = coord.X - sightradius; x <= coord.X+sightradius; x++ {
			if x == coord.X && y == coord.Y {
				C.strcat(&buf[0], getmapchar(mapnums[y][x], ch, start, vnum))
			} else if x > MAP_COLS || x < 0 {
				if lasty != TRUE && y > -1 && y < 200 {
					C.strcat(&buf[0], libc.CString("@D?"))
					lasty = TRUE
				}
			} else if y > MAP_ROWS || y < 0 {
				if y == -1 || y == 200 {
					C.strcat(&buf[0], libc.CString("@D??"))
				}
			} else {
				C.strcat(&buf[0], getmapchar(mapnums[y][x], ch, start, vnum))
			}
		}
		C.strcat(&buf[0], libc.CString("\n"))
		lasty = FALSE
	}
	send_to_char(ch, &buf[0])
	buf2[0] = '\x00'
	buf[0] = '\x00'
	if type_ == 0 {
		send_to_char(ch, libc.CString("\n@b______________________________________________________________________@n"))
	}
}
