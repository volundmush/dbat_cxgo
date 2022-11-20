package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func add_object(newobj *obj_data, ovnum int) int {
	var (
		found int = int(-1)
		rznum int = real_zone_by_thing(ovnum)
	)
	if (func() int {
		p := &newobj.Item_number
		newobj.Item_number = real_object(ovnum)
		return *p
	}()) != int(-1) {
		copy_object(&obj_proto[newobj.Item_number], newobj)
		update_objects(&obj_proto[newobj.Item_number])
		add_to_save_list(zone_table[rznum].Number, SL_OBJ)
		return newobj.Item_number
	}
	found = insert_object(newobj, ovnum)
	adjust_objects(found)
	add_to_save_list(zone_table[rznum].Number, SL_OBJ)
	return found
}
func update_objects(refobj *obj_data) int {
	var (
		obj   *obj_data
		swap  obj_data
		count int = 0
	)
	for obj = object_list; obj != nil; obj = obj.Next {
		if obj.Item_number != refobj.Item_number {
			continue
		}
		count++
		swap = *obj
		*obj = *refobj
		obj.In_room = swap.In_room
		obj.Carried_by = swap.Carried_by
		obj.Worn_by = swap.Worn_by
		obj.Worn_on = swap.Worn_on
		obj.In_obj = swap.In_obj
		obj.Contains = swap.Contains
		obj.Next_content = swap.Next_content
		obj.Next = swap.Next
	}
	return count
}
func adjust_objects(refpt int) int {
	var (
		shop   int
		i      int
		zone   int
		cmd_no int
		obj    *obj_data
	)
	if refpt < 0 || refpt > top_of_objt {
		return -1
	}
	for obj = object_list; obj != nil; obj = obj.Next {
		obj.Item_number += int(libc.BoolToInt(obj.Item_number != int(-1) && obj.Item_number >= refpt))
	}
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int(zone_table[zone].Cmd[cmd_no].Command) != 'S'; cmd_no++ {
			switch zone_table[zone].Cmd[cmd_no].Command {
			case 'P':
				zone_table[zone].Cmd[cmd_no].Arg3 += int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg3 >= refpt))
				fallthrough
			case 'O':
				fallthrough
			case 'G':
				fallthrough
			case 'E':
				zone_table[zone].Cmd[cmd_no].Arg1 += int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg1 >= refpt))
			case 'R':
				zone_table[zone].Cmd[cmd_no].Arg2 += int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg2 >= refpt))
			}
		}
	}
	for shop = 0; shop <= top_shop; shop++ {
		for i = 0; (shop_index[shop].Producing[i]) != int(-1); i++ {
			shop_index[shop].Producing[i] += int(libc.BoolToInt((shop_index[shop].Producing[i]) >= refpt))
		}
	}
	return refpt
}
func insert_object(obj *obj_data, ovnum int) int {
	var i int
	top_of_objt++
	// todo: fix this
	//obj_index = []index_data((*index_data)(libc.Realloc(unsafe.Pointer(&obj_index[0]), top_of_objt*int(unsafe.Sizeof(index_data{}))+1)))
	//obj_proto = []obj_data((*obj_data)(libc.Realloc(unsafe.Pointer(&obj_proto[0]), top_of_objt*int(unsafe.Sizeof(obj_data{}))+1)))
	for i = top_of_objt; i > 0; i-- {
		if ovnum > obj_index[i-1].Vnum {
			return index_object(obj, ovnum, i)
		}
		obj_index[i] = obj_index[i-1]
		obj_proto[i] = obj_proto[i-1]
		obj_proto[i].Item_number = i
		htree_add(obj_htree, obj_index[i].Vnum, i)
	}
	return index_object(obj, ovnum, 0)
}
func index_object(obj *obj_data, ovnum int, ornum int) int {
	if obj == nil || ovnum < 0 || ornum < 0 || ornum > top_of_objt {
		return -1
	}
	obj.Item_number = ornum
	obj_index[ornum].Vnum = ovnum
	obj_index[ornum].Number = 0
	obj_index[ornum].Func = nil
	copy_object_preserve(&obj_proto[ornum], obj)
	obj_proto[ornum].In_room = -1
	htree_add(obj_htree, obj_index[ornum].Vnum, ornum)
	return ornum
}
func save_objects(zone_num int) bool {
	var (
		cmfname     [128]byte
		buf         [64936]byte
		ebuf1       [64936]byte
		ebuf2       [64936]byte
		ebuf3       [64936]byte
		ebuf4       [64936]byte
		wbuf1       [64936]byte
		wbuf2       [64936]byte
		wbuf3       [64936]byte
		wbuf4       [64936]byte
		pbuf1       [64936]byte
		pbuf2       [64936]byte
		pbuf3       [64936]byte
		pbuf4       [64936]byte
		counter     int
		counter2    int
		realcounter int
		fp          *stdio.File
		obj         *obj_data
		ex_desc     *extra_descr_data
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: OasisOLC: save_objects: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return false
	}
	stdio.Snprintf(&cmfname[0], int(128), "%s%d.new", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&cmfname[0]), "w+")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("SYSERR: OLC: Cannot open objects file %s!"), &cmfname[0])
		return false
	}
	for counter = genolc_zone_bottom(zone_num); counter <= zone_table[zone_num].Top; counter++ {
		if (func() int {
			realcounter = real_object(counter)
			return realcounter
		}()) != int(-1) {
			if (func() *obj_data {
				obj = &obj_proto[realcounter]
				return obj
			}()).Action_description != nil {
				libc.StrNCpy(&buf[0], obj.Action_description, int(64936-1))
				strip_cr(&buf[0])
			} else {
				buf[0] = '\x00'
			}
			stdio.Fprintf(fp, "#%d\n%s~\n%s~\n%s~\n%s~\n", GET_OBJ_VNUM(obj), func() *byte {
				if obj.Name != nil && *obj.Name != 0 {
					return obj.Name
				}
				return libc.CString("undefined")
			}(), func() *byte {
				if obj.Short_description != nil && *obj.Short_description != 0 {
					return obj.Short_description
				}
				return libc.CString("undefined")
			}(), func() *byte {
				if obj.Description != nil && *obj.Description != 0 {
					return obj.Description
				}
				return libc.CString("undefined")
			}(), &buf[0])
			sprintascii(&ebuf1[0], obj.Extra_flags[0])
			sprintascii(&ebuf2[0], obj.Extra_flags[1])
			sprintascii(&ebuf3[0], obj.Extra_flags[2])
			sprintascii(&ebuf4[0], obj.Extra_flags[3])
			sprintascii(&wbuf1[0], obj.Wear_flags[0])
			sprintascii(&wbuf2[0], obj.Wear_flags[1])
			sprintascii(&wbuf3[0], obj.Wear_flags[2])
			sprintascii(&wbuf4[0], obj.Wear_flags[3])
			sprintascii(&pbuf1[0], obj.Bitvector[0])
			sprintascii(&pbuf2[0], obj.Bitvector[1])
			sprintascii(&pbuf3[0], obj.Bitvector[2])
			sprintascii(&pbuf4[0], obj.Bitvector[3])
			stdio.Fprintf(fp, "%d %s %s %s %s %s %s %s %s %s %s %s %s\n%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n%lld %d %d %d\n", obj.Type_flag, &ebuf1[0], &ebuf2[0], &ebuf3[0], &ebuf4[0], &wbuf1[0], &wbuf2[0], &wbuf3[0], &wbuf4[0], &pbuf1[0], &pbuf2[0], &pbuf3[0], &pbuf4[0], obj.Value[0], obj.Value[1], obj.Value[2], obj.Value[3], obj.Value[4], obj.Value[5], obj.Value[6], obj.Value[7], obj.Value[8], obj.Value[9], obj.Value[10], obj.Value[11], obj.Value[12], obj.Value[13], obj.Value[14], obj.Value[15], obj.Weight, obj.Cost, obj.Cost_per_day, obj.Level)
			script_save_to_disk(fp, unsafe.Pointer(obj), OBJ_TRIGGER)
			stdio.Fprintf(fp, "Z\n%d\n", obj.Size)
			if obj.Ex_description != nil {
				for ex_desc = obj.Ex_description; ex_desc != nil; ex_desc = ex_desc.Next {
					if ex_desc.Keyword == nil || ex_desc.Description == nil || *ex_desc.Keyword == 0 || *ex_desc.Description == 0 {
						mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("SYSERR: OLC: oedit_save_to_disk: Corrupt ex_desc!"))
						continue
					}
					libc.StrNCpy(&buf[0], ex_desc.Description, int(64936-1))
					strip_cr(&buf[0])
					stdio.Fprintf(fp, "E\n%s~\n%s~\n", ex_desc.Keyword, &buf[0])
				}
			}
			for counter2 = 0; counter2 < MAX_OBJ_AFFECT; counter2++ {
				if obj.Affected[counter2].Modifier != 0 {
					stdio.Fprintf(fp, "A\n%d %d %d\n", obj.Affected[counter2].Location, obj.Affected[counter2].Modifier, obj.Affected[counter2].Specific)
				}
			}
			if obj.Sbinfo != nil {
				for counter2 = 0; counter2 < SKILL_TABLE_SIZE; counter2++ {
					if obj.Sbinfo[counter2].Spellname == 0 {
						break
					}
					stdio.Fprintf(fp, "S\n%d %d\n", obj.Sbinfo[counter2].Spellname, obj.Sbinfo[counter2].Pages)
					continue
				}
			}
		}
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.obj", LIB_WORLD, zone_table[zone_num].Number)
	stdio.Remove(libc.GoString(&buf[0]))
	stdio.Rename(libc.GoString(&cmfname[0]), libc.GoString(&buf[0]))
	if in_save_list(zone_table[zone_num].Number, SL_OBJ) {
		remove_from_save_list(zone_table[zone_num].Number, SL_OBJ)
		create_world_index(zone_table[zone_num].Number, libc.CString("obj"))
		basic_mud_log(libc.CString("GenOLC: save_objects: Saving objects '%s'"), &buf[0])
	}
	return true
}
func free_object_strings(obj *obj_data) {
	if obj.Name != nil {
		libc.Free(unsafe.Pointer(obj.Name))
	}
	if obj.Description != nil {
		libc.Free(unsafe.Pointer(obj.Description))
	}
	if obj.Short_description != nil {
		libc.Free(unsafe.Pointer(obj.Short_description))
	}
	if obj.Action_description != nil {
		libc.Free(unsafe.Pointer(obj.Action_description))
	}
	if obj.Ex_description != nil {
		free_ex_descriptions(obj.Ex_description)
	}
}
func free_object_strings_proto(obj *obj_data) {
	var robj_num int = obj.Item_number
	if obj.Name != nil && obj.Name != obj_proto[robj_num].Name {
		libc.Free(unsafe.Pointer(obj.Name))
	}
	if obj.Description != nil && obj.Description != obj_proto[robj_num].Description {
		libc.Free(unsafe.Pointer(obj.Description))
	}
	if obj.Short_description != nil && obj.Short_description != obj_proto[robj_num].Short_description {
		libc.Free(unsafe.Pointer(obj.Short_description))
	}
	if obj.Action_description != nil && obj.Action_description != obj_proto[robj_num].Action_description {
		libc.Free(unsafe.Pointer(obj.Action_description))
	}
	if obj.Ex_description != nil {
		var (
			thised   *extra_descr_data
			plist    *extra_descr_data
			next_one *extra_descr_data
			ok_key   int
			ok_desc  int
			ok_item  int
		)
		for thised = obj.Ex_description; thised != nil; thised = next_one {
			next_one = thised.Next
			for func() *extra_descr_data {
				ok_item = func() int {
					ok_key = func() int {
						ok_desc = 1
						return ok_desc
					}()
					return ok_key
				}()
				return func() *extra_descr_data {
					plist = obj_proto[robj_num].Ex_description
					return plist
				}()
			}(); plist != nil; plist = plist.Next {
				if plist.Keyword == thised.Keyword {
					ok_key = 0
				}
				if plist.Description == thised.Description {
					ok_desc = 0
				}
				if plist == thised {
					ok_item = 0
				}
			}
			if thised.Keyword != nil && ok_key != 0 {
				libc.Free(unsafe.Pointer(thised.Keyword))
			}
			if thised.Description != nil && ok_desc != 0 {
				libc.Free(unsafe.Pointer(thised.Description))
			}
			if ok_item != 0 {
				libc.Free(unsafe.Pointer(thised))
			}
		}
	}
}
func copy_object_strings(to *obj_data, from *obj_data) {
	if from.Name != nil {
		to.Name = libc.StrDup(from.Name)
	} else {
		to.Name = nil
	}
	if from.Description != nil {
		to.Description = libc.StrDup(from.Description)
	} else {
		to.Description = nil
	}
	if from.Short_description != nil {
		to.Short_description = libc.StrDup(from.Short_description)
	} else {
		to.Short_description = nil
	}
	if from.Action_description != nil {
		to.Action_description = libc.StrDup(from.Action_description)
	} else {
		to.Action_description = nil
	}
	if from.Ex_description != nil {
		copy_ex_descriptions(&to.Ex_description, from.Ex_description)
	} else {
		to.Ex_description = nil
	}
}
func copy_object(to *obj_data, from *obj_data) int {
	free_object_strings(to)
	return int(libc.BoolToInt(copy_object_main(to, from, 1)))
}
func copy_object_preserve(to *obj_data, from *obj_data) int {
	return int(libc.BoolToInt(copy_object_main(to, from, 0)))
}
func copy_object_main(to *obj_data, from *obj_data, free_object int) bool {
	*to = *from
	copy_object_strings(to, from)
	return true
}
func delete_object(rnum int) int {
	var (
		i      int
		zrnum  int
		obj    *obj_data
		tmp    *obj_data
		shop   int
		j      int
		zone   int
		cmd_no int
	)
	if rnum == int(-1) || rnum > top_of_objt {
		return -1
	}
	obj = &obj_proto[rnum]
	zrnum = real_zone_by_thing(GET_OBJ_VNUM(obj))
	htree_del(obj_htree, obj.Item_number)
	basic_mud_log(libc.CString("GenOLC: delete_object: Deleting object #%d (%s)."), GET_OBJ_VNUM(obj), obj.Short_description)
	for tmp = object_list; tmp != nil; tmp = tmp.Next {
		if tmp.Item_number != obj.Item_number {
			continue
		}
		if tmp.Contains != nil {
			var (
				this_content *obj_data
				next_content *obj_data
			)
			for this_content = tmp.Contains; this_content != nil; this_content = next_content {
				next_content = this_content.Next_content
				if tmp.In_room != 0 {
					obj_from_obj(this_content)
					obj_to_room(this_content, tmp.In_room)
				} else if tmp.Worn_by != nil || tmp.Carried_by != nil {
					obj_from_char(this_content)
					obj_to_char(this_content, tmp.Carried_by)
				} else if tmp.In_obj != nil {
					obj_from_obj(this_content)
					obj_to_obj(this_content, tmp.In_obj)
				}
			}
		}
		extract_obj(tmp)
	}
	if obj_index[rnum].Number != 0 {
		panic("assert failed")
	}
	for tmp = object_list; tmp != nil; tmp = tmp.Next {
		tmp.Item_number -= int(libc.BoolToInt(tmp.Item_number > rnum))
	}
	for i = rnum; i < top_of_objt; i++ {
		obj_index[i] = obj_index[i+1]
		obj_proto[i] = obj_proto[i+1]
		obj_proto[i].Item_number = i
	}
	top_of_objt--
	// todo: fix this
	//obj_index = []index_data((*index_data)(libc.Realloc(unsafe.Pointer(&obj_index[0]), top_of_objt*int(unsafe.Sizeof(index_data{}))+1)))
	//obj_proto = []obj_data((*obj_data)(libc.Realloc(unsafe.Pointer(&obj_proto[0]), top_of_objt*int(unsafe.Sizeof(obj_data{}))+1)))
	for shop = 0; shop <= top_shop; shop++ {
		for j = 0; (shop_index[shop].Producing[j]) != int(-1); j++ {
			shop_index[shop].Producing[j] -= int(libc.BoolToInt((shop_index[shop].Producing[j]) > rnum))
		}
	}
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int(zone_table[zone].Cmd[cmd_no].Command) != 'S'; cmd_no++ {
			switch zone_table[zone].Cmd[cmd_no].Command {
			case 'P':
				if zone_table[zone].Cmd[cmd_no].Arg3 == rnum {
					delete_zone_command(&zone_table[zone], cmd_no)
				} else {
					zone_table[zone].Cmd[cmd_no].Arg3 -= int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg3 > rnum))
				}
			case 'O':
				fallthrough
			case 'G':
				fallthrough
			case 'E':
				if zone_table[zone].Cmd[cmd_no].Arg1 == rnum {
					delete_zone_command(&zone_table[zone], cmd_no)
				} else {
					zone_table[zone].Cmd[cmd_no].Arg1 -= int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg1 > rnum))
				}
			case 'R':
				if zone_table[zone].Cmd[cmd_no].Arg2 == rnum {
					delete_zone_command(&zone_table[zone], cmd_no)
				} else {
					zone_table[zone].Cmd[cmd_no].Arg2 -= int(libc.BoolToInt(zone_table[zone].Cmd[cmd_no].Arg2 > rnum))
				}
			}
		}
	}
	save_objects(zrnum)
	return rnum
}
