package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func real_zone_by_thing(vznum room_vnum) zone_rnum {
	var (
		bot      zone_rnum
		top      zone_rnum
		mid      zone_rnum
		last_top zone_rnum
	)
	_ = last_top
	var low int
	var high int
	bot = 0
	top = top_of_zone_table
	if genolc_zone_bottom(bot) > vznum || zone_table[top].Top < vznum {
		return -1
	}
	for bot <= top {
		last_top = top
		mid = (bot + top) / 2
		low = int(genolc_zone_bottom(mid))
		high = int(zone_table[mid].Top)
		if low <= int(vznum) && vznum <= room_vnum(high) {
			return mid
		}
		if low > int(vznum) {
			top = mid - 1
		} else {
			bot = mid + 1
		}
	}
	return -1
}
func create_new_zone(vzone_num zone_vnum, bottom room_vnum, top room_vnum, error **byte) zone_rnum {
	var (
		fp    *stdio.File
		zone  *zone_data
		i     int
		rznum zone_rnum
		buf   [64936]byte
	)
	if vzone_num < 0 {
		*error = libc.CString("You can't make negative zones.\r\n")
		return -1
	} else if bottom > top {
		*error = libc.CString("Bottom room cannot be greater than top room.\r\n")
		return -1
	}
	for i = 0; i < int(top_of_zone_table); i++ {
		if zone_table[i].Number == vzone_num {
			*error = libc.CString("That virtual zone already exists.\r\n")
			return -1
		}
	}
	stdio.Snprintf(&buf[0], int(64936), "%s%d.zon", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new zone file."))
		*error = libc.CString("Could not write zone file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "#%d\nNone~\nNew Zone~\n%d %d 30 2\nS\n$\n", vzone_num, bottom, top)
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.wld", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new world file."))
		*error = libc.CString("Could not write world file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "#%d\nThe Beginning~\nNot much here.\n~\n%d 0 0\nS\n$\n", bottom, vzone_num)
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.mob", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new mob file."))
		*error = libc.CString("Could not write mobile file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "$\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.obj", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new obj file."))
		*error = libc.CString("Could not write object file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "$\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.shp", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new shop file."))
		*error = libc.CString("Could not write shop file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.trg", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new trigger file"))
		*error = libc.CString("Could not write trigger file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s/%i.gld", LIB_WORLD, vzone_num)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&buf[0]), "w")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Can't write new guild file"))
		*error = libc.CString("Could not write guild file.\r\n")
		return -1
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	create_world_index(int(vzone_num), libc.CString("zon"))
	create_world_index(int(vzone_num), libc.CString("wld"))
	create_world_index(int(vzone_num), libc.CString("mob"))
	create_world_index(int(vzone_num), libc.CString("obj"))
	create_world_index(int(vzone_num), libc.CString("shp"))
	create_world_index(int(vzone_num), libc.CString("trg"))
	create_world_index(int(vzone_num), libc.CString("gld"))
	// todo : figure this out
	//zone_table = []zone_data((*zone_data)(libc.Realloc(unsafe.Pointer(&zone_table[0]), int(top_of_zone_table*zone_rnum(unsafe.Sizeof(zone_data{}))+2))))
	zone_table[top_of_zone_table+1].Number = 32000
	if vzone_num > zone_table[top_of_zone_table].Number {
		rznum = top_of_zone_table + 1
	} else {
		var (
			j    int
			room int
		)
		for i = int(top_of_zone_table + 1); i > 0 && vzone_num < zone_table[i-1].Number; i-- {
			zone_table[i] = zone_table[i-1]
			for j = int(zone_table[i].Bot); j <= int(zone_table[i].Top); j++ {
				if (func() int {
					room = int(real_room(room_vnum(j)))
					return room
				}()) != int(-1) {
					world[room].Zone++
				}
			}
		}
		rznum = zone_rnum(i)
	}
	zone = &zone_table[rznum]
	zone.Name = libc.CString("New Zone")
	zone.Number = vzone_num
	zone.Builders = libc.CString("None")
	zone.Bot = bottom
	zone.Top = top
	zone.Lifespan = 30
	zone.Max_level = ADMLVL_IMPL
	zone.Cmd = make([]reset_com, 1)
	zone.Cmd[0].Command = 'S'
	top_of_zone_table++
	for i = int(top_of_world); i > 0; i-- {
		if world[i].Zone < real_zone(zone_vnum(rznum)) {
			break
		} else {
			world[i].Zone = real_zone_by_thing(room_vnum(libc.BoolToInt(GET_ROOM_VNUM(room_rnum(i)))))
		}
	}
	add_to_save_list(zone.Number, SL_ZON)
	return rznum
}
func create_world_index(znum int, type_ *byte) {
	var (
		newfile  *stdio.File
		oldfile  *stdio.File
		new_name [32]byte
		old_name [32]byte
		prefix   *byte
		num      int
		found    int = FALSE
		buf      [64936]byte
		buf1     [64936]byte
	)
	switch *type_ {
	case 'z':
		prefix = libc.CString(LIB_WORLD)
	case 'w':
		prefix = libc.CString(LIB_WORLD)
	case 'o':
		prefix = libc.CString(LIB_WORLD)
	case 'm':
		prefix = libc.CString(LIB_WORLD)
	case 's':
		prefix = libc.CString(LIB_WORLD)
	case 't':
		prefix = libc.CString(LIB_WORLD)
	case 'g':
		prefix = libc.CString(LIB_WORLD)
	default:
		return
	}
	stdio.Snprintf(&old_name[0], int(32), "%s/index", prefix)
	stdio.Snprintf(&new_name[0], int(32), "%s/newindex", prefix)
	if (func() *stdio.File {
		oldfile = stdio.FOpen(libc.GoString(&old_name[0]), "r")
		return oldfile
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Failed to open %s."), &old_name[0])
		return
	} else if (func() *stdio.File {
		newfile = stdio.FOpen(libc.GoString(&new_name[0]), "w")
		return newfile
	}()) == nil {
		mudlog(BRF, ADMLVL_IMPL, TRUE, libc.CString("SYSERR: OLC: Failed to open %s."), &new_name[0])
		oldfile.Close()
		return
	}
	stdio.Snprintf(&buf1[0], int(64936), "%d.%s", znum, type_)
	for get_line(oldfile, &buf[0]) != 0 {
		if buf[0] == '$' {
			stdio.Fprintf(newfile, "%s", func() *byte {
				if found == 0 {
					return libc.StrNCat(&buf1[0], libc.CString("\n$\n"), int(64936-1))
				}
				return libc.CString("$\n")
			}())
			break
		} else if found == 0 {
			stdio.Sscanf(&buf[0], "%d", &num)
			if num == znum {
				found = TRUE
			} else if num > znum {
				found = TRUE
				stdio.Fprintf(newfile, "%s\n", &buf1[0])
			}
		}
		stdio.Fprintf(newfile, "%s\n", &buf[0])
	}
	newfile.Close()
	oldfile.Close()
	stdio.Remove(libc.GoString(&old_name[0]))
	stdio.Rename(libc.GoString(&new_name[0]), libc.GoString(&old_name[0]))
}
func remove_room_zone_commands(zone zone_rnum, room_num room_rnum) {
	var (
		subcmd   int = 0
		cmd_room int = -2
	)
	for int(zone_table[zone].Cmd[subcmd].Command) != 'S' {
		switch zone_table[zone].Cmd[subcmd].Command {
		case 'M':
			fallthrough
		case 'O':
			fallthrough
		case 'T':
			fallthrough
		case 'V':
			cmd_room = int(zone_table[zone].Cmd[subcmd].Arg3)
		case 'D':
			fallthrough
		case 'R':
			cmd_room = int(zone_table[zone].Cmd[subcmd].Arg1)
		default:
		}
		if cmd_room == int(room_num) {
			remove_cmd_from_list((**reset_com)(unsafe.Pointer(&zone_table[zone].Cmd[0])), subcmd)
		} else {
			subcmd++
		}
	}
}
func save_zone(zone_num zone_rnum) int {
	var (
		subcmd  int
		arg1    int = -1
		arg2    int = -1
		arg3    int = -1
		arg4    int = -1
		arg5    int = -1
		fname   [128]byte
		oldname [128]byte
		comment *byte = nil
		zfile   *stdio.File
		zbuf1   [64936]byte
		zbuf2   [64936]byte
		zbuf3   [64936]byte
		zbuf4   [64936]byte
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_zone: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return FALSE
	}
	stdio.Snprintf(&fname[0], int(128), "%s%d.new", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		zfile = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return zfile
	}()) == nil {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: save_zones:  Can't write zone %d."), zone_table[zone_num].Number)
		return FALSE
	}
	sprintascii(&zbuf1[0], zone_table[zone_num].Zone_flags[0])
	sprintascii(&zbuf2[0], zone_table[zone_num].Zone_flags[1])
	sprintascii(&zbuf3[0], zone_table[zone_num].Zone_flags[2])
	sprintascii(&zbuf4[0], zone_table[zone_num].Zone_flags[3])
	stdio.Fprintf(zfile, "@Version: %d\n", CUR_ZONE_VERSION)
	stdio.Fprintf(zfile, "#%d\n%s~\n%s~\n%d %d %d %d %s %s %s %s %d %d\n", zone_table[zone_num].Number, func() *byte {
		if zone_table[zone_num].Builders != nil && *zone_table[zone_num].Builders != 0 {
			return zone_table[zone_num].Builders
		}
		return libc.CString("None.")
	}(), func() *byte {
		if zone_table[zone_num].Name != nil && *zone_table[zone_num].Name != 0 {
			return zone_table[zone_num].Name
		}
		return libc.CString("undefined")
	}(), genolc_zone_bottom(zone_num), zone_table[zone_num].Top, zone_table[zone_num].Lifespan, zone_table[zone_num].Reset_mode, &zbuf1[0], &zbuf2[0], &zbuf3[0], &zbuf4[0], zone_table[zone_num].Min_level, zone_table[zone_num].Max_level)
	for subcmd = 0; int(zone_table[zone_num].Cmd[subcmd].Command) != 'S'; subcmd++ {
		switch zone_table[zone_num].Cmd[subcmd].Command {
		case 'M':
			arg1 = int(mob_index[zone_table[zone_num].Cmd[subcmd].Arg1].Vnum)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(world[zone_table[zone_num].Cmd[subcmd].Arg3].Number)
			arg4 = int(zone_table[zone_num].Cmd[subcmd].Arg4)
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = mob_proto[zone_table[zone_num].Cmd[subcmd].Arg1].Short_descr
		case 'O':
			arg1 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg1].Vnum)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(world[zone_table[zone_num].Cmd[subcmd].Arg3].Number)
			arg4 = int(zone_table[zone_num].Cmd[subcmd].Arg4)
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = obj_proto[zone_table[zone_num].Cmd[subcmd].Arg1].Short_description
		case 'G':
			arg1 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg1].Vnum)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = -1
			arg4 = -1
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = obj_proto[zone_table[zone_num].Cmd[subcmd].Arg1].Short_description
		case 'E':
			arg1 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg1].Vnum)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(zone_table[zone_num].Cmd[subcmd].Arg3)
			arg4 = -1
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = obj_proto[zone_table[zone_num].Cmd[subcmd].Arg1].Short_description
		case 'P':
			arg1 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg1].Vnum)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg3].Vnum)
			arg4 = -1
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = obj_proto[zone_table[zone_num].Cmd[subcmd].Arg1].Short_description
		case 'D':
			arg1 = int(world[zone_table[zone_num].Cmd[subcmd].Arg1].Number)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(zone_table[zone_num].Cmd[subcmd].Arg3)
			comment = world[zone_table[zone_num].Cmd[subcmd].Arg1].Name
		case 'R':
			arg1 = int(world[zone_table[zone_num].Cmd[subcmd].Arg1].Number)
			arg2 = int(obj_index[zone_table[zone_num].Cmd[subcmd].Arg2].Vnum)
			comment = obj_proto[zone_table[zone_num].Cmd[subcmd].Arg2].Short_description
			arg3 = -1
		case 'T':
			arg1 = int(zone_table[zone_num].Cmd[subcmd].Arg1)
			arg2 = int(trig_index[zone_table[zone_num].Cmd[subcmd].Arg2].Vnum)
			arg3 = int(world[zone_table[zone_num].Cmd[subcmd].Arg3].Number)
			arg4 = -1
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
			comment = trig_index[real_trigger(trig_vnum(arg2))].Proto.Name
		case 'V':
			arg1 = int(zone_table[zone_num].Cmd[subcmd].Arg1)
			arg2 = int(zone_table[zone_num].Cmd[subcmd].Arg2)
			arg3 = int(world[zone_table[zone_num].Cmd[subcmd].Arg3].Number)
			arg4 = -1
			arg5 = int(zone_table[zone_num].Cmd[subcmd].Arg5)
		case '*':
			continue
		default:
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: z_save_to_disk(): Unknown cmd '%c' - NOT saving"), zone_table[zone_num].Cmd[subcmd].Command)
			continue
		}
		if int(zone_table[zone_num].Cmd[subcmd].Command) != 'V' {
			stdio.Fprintf(zfile, "%c %d %d %d %d %d %d \t(%s)\n", zone_table[zone_num].Cmd[subcmd].Command, zone_table[zone_num].Cmd[subcmd].If_flag, arg1, arg2, arg3, arg4, arg5, comment)
		} else {
			stdio.Fprintf(zfile, "%c %d %d %d %d %d %d %s %s\n", zone_table[zone_num].Cmd[subcmd].Command, zone_table[zone_num].Cmd[subcmd].If_flag, arg1, arg2, arg3, arg4, arg5, zone_table[zone_num].Cmd[subcmd].Sarg1, zone_table[zone_num].Cmd[subcmd].Sarg2)
		}
	}
	zfile.PutS(libc.CString("S\n$\n"))
	zfile.Close()
	stdio.Snprintf(&oldname[0], int(128), "%s%d.zon", LIB_WORLD, zone_table[zone_num].Number)
	stdio.Remove(libc.GoString(&oldname[0]))
	stdio.Rename(libc.GoString(&fname[0]), libc.GoString(&oldname[0]))
	if in_save_list(zone_table[zone_num].Number, SL_ZON) != 0 {
		remove_from_save_list(zone_table[zone_num].Number, SL_ZON)
		create_world_index(int(zone_table[zone_num].Number), libc.CString("zon"))
		basic_mud_log(libc.CString("GenOLC: save_zone: Saving zone '%s'"), &oldname[0])
	}
	return TRUE
}
func count_commands(list []reset_com) int {
	var count int = 0
	for int(list[count].Command) != 'S' {
		count++
	}
	return count
}
func add_cmd_to_list(list **reset_com, newcmd *reset_com, pos int) {
	var (
		count   int
		i       int
		l       int
		newlist *reset_com
	)
	// todo: figure thi sout
	//count = count_commands([]reset_com(*list))
	newlist = &make([]reset_com, count+2)[0]
	for func() int {
		i = 0
		return func() int {
			l = 0
			return l
		}()
	}(); i <= count; i++ {
		if i == pos {
			*(*reset_com)(unsafe.Add(unsafe.Pointer(newlist), unsafe.Sizeof(reset_com{})*uintptr(i))) = *newcmd
		} else {
			*(*reset_com)(unsafe.Add(unsafe.Pointer(newlist), unsafe.Sizeof(reset_com{})*uintptr(i))) = *(*reset_com)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(reset_com{})*uintptr(func() int {
				p := &l
				x := *p
				*p++
				return x
			}())))
		}
	}
	(*(*reset_com)(unsafe.Add(unsafe.Pointer(newlist), unsafe.Sizeof(reset_com{})*uintptr(count+1)))).Command = 'S'
	libc.Free(unsafe.Pointer(*list))
	*list = newlist
}
func remove_cmd_from_list(list **reset_com, pos int) {
	var (
		count   int
		i       int
		l       int
		newlist *reset_com
	)
	// todo : figure this out
	//count = count_commands([]reset_com(*list))
	newlist = &make([]reset_com, count)[0]
	for func() int {
		i = 0
		return func() int {
			l = 0
			return l
		}()
	}(); i < count; i++ {
		if i != pos {
			*(*reset_com)(unsafe.Add(unsafe.Pointer(newlist), unsafe.Sizeof(reset_com{})*uintptr(func() int {
				p := &l
				x := *p
				*p++
				return x
			}()))) = *(*reset_com)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(reset_com{})*uintptr(i)))
		}
	}
	(*(*reset_com)(unsafe.Add(unsafe.Pointer(newlist), unsafe.Sizeof(reset_com{})*uintptr(count-1)))).Command = 'S'
	libc.Free(unsafe.Pointer(*list))
	*list = newlist
}
func new_command(zone *zone_data, pos int) int {
	var (
		subcmd  int = 0
		new_com reset_com
	)
	for int(zone.Cmd[subcmd].Command) != 'S' {
		subcmd++
	}
	if pos < 0 || pos > subcmd {
		return 0
	}
	new_com.Command = 'N'
	add_cmd_to_list((**reset_com)(unsafe.Pointer(&zone.Cmd[0])), &new_com, pos)
	return 1
}
func delete_zone_command(zone *zone_data, pos int) {
	var subcmd int = 0
	for int(zone.Cmd[subcmd].Command) != 'S' {
		subcmd++
	}
	if pos < 0 || pos >= subcmd {
		return
	}
	remove_cmd_from_list((**reset_com)(unsafe.Pointer(&zone.Cmd[0])), pos)
}
