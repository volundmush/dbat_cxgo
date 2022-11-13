package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func add_object(newobj *obj_data, ovnum obj_vnum) obj_rnum {
	var (
		found int       = int(-1)
		rznum zone_rnum = real_zone_by_thing(room_vnum(ovnum))
	)
	if (func() obj_vnum {
		p := &newobj.Item_number
		newobj.Item_number = obj_vnum(real_object(ovnum))
		return *p
	}()) != obj_vnum(-1) {
		copy_object((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(newobj.Item_number))), newobj)
		update_objects((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(newobj.Item_number))))
		add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rznum)))).Number, SL_OBJ)
		return obj_rnum(newobj.Item_number)
	}
	found = int(insert_object(newobj, ovnum))
	adjust_objects(obj_rnum(found))
	add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rznum)))).Number, SL_OBJ)
	return obj_rnum(found)
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
func adjust_objects(refpt obj_rnum) obj_rnum {
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
		obj.Item_number += obj_vnum(libc.BoolToInt(obj.Item_number != obj_vnum(-1) && obj.Item_number >= obj_vnum(refpt)))
	}
	for zone = 0; zone <= int(top_of_zone_table); zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command {
			case 'P':
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3 >= vnum(refpt)))
				fallthrough
			case 'O':
				fallthrough
			case 'G':
				fallthrough
			case 'E':
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 >= vnum(refpt)))
			case 'R':
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 >= vnum(refpt)))
			}
		}
	}
	for shop = 0; shop <= top_shop; shop++ {
		for i = 0; (*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i)))) != obj_vnum(-1); i++ {
			*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i))) += obj_vnum(libc.BoolToInt((*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i)))) >= obj_vnum(refpt)))
		}
	}
	return refpt
}
func insert_object(obj *obj_data, ovnum obj_vnum) obj_rnum {
	var i obj_rnum
	top_of_objt++
	obj_index = (*index_data)(libc.Realloc(unsafe.Pointer(obj_index), int(top_of_objt*obj_rnum(unsafe.Sizeof(index_data{}))+1)))
	obj_proto = (*obj_data)(libc.Realloc(unsafe.Pointer(obj_proto), int(top_of_objt*obj_rnum(unsafe.Sizeof(obj_data{}))+1)))
	for i = top_of_objt; i > 0; i-- {
		if ovnum > obj_vnum((*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i-1)))).Vnum) {
			return index_object(obj, ovnum, i)
		}
		*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i))) = *(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i-1)))
		*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i))) = *(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i-1)))
		(*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Item_number = obj_vnum(i)
		htree_add(obj_htree, int64((*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum), int64(i))
	}
	return index_object(obj, ovnum, 0)
}
func index_object(obj *obj_data, ovnum obj_vnum, ornum obj_rnum) obj_rnum {
	if obj == nil || ovnum < 0 || ornum < 0 || ornum > top_of_objt {
		return -1
	}
	obj.Item_number = obj_vnum(ornum)
	(*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(ornum)))).Vnum = mob_vnum(ovnum)
	(*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(ornum)))).Number = 0
	(*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(ornum)))).Func = nil
	copy_object_preserve((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(ornum))), obj)
	(*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(ornum)))).In_room = -1
	htree_add(obj_htree, int64((*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(ornum)))).Vnum), int64(ornum))
	return ornum
}
func save_objects(zone_num zone_rnum) int {
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
		return FALSE
	}
	stdio.Snprintf(&cmfname[0], int(128), "%s%d.new", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	if (func() *stdio.File {
		fp = stdio.FOpen(libc.GoString(&cmfname[0]), "w+")
		return fp
	}()) == nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: OLC: Cannot open objects file %s!"), &cmfname[0])
		return FALSE
	}
	for counter = int(genolc_zone_bottom(zone_num)); counter <= int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Top); counter++ {
		if (func() int {
			realcounter = int(real_object(obj_vnum(counter)))
			return realcounter
		}()) != int(-1) {
			if (func() *obj_data {
				obj = (*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(realcounter)))
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
			sprintascii(&wbuf1[0], bitvector_t(int32(obj.Wear_flags[0])))
			sprintascii(&wbuf2[0], bitvector_t(int32(obj.Wear_flags[1])))
			sprintascii(&wbuf3[0], bitvector_t(int32(obj.Wear_flags[2])))
			sprintascii(&wbuf4[0], bitvector_t(int32(obj.Wear_flags[3])))
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
						mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: OLC: oedit_save_to_disk: Corrupt ex_desc!"))
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
					if (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter2)))).Spellname == 0 {
						break
					}
					stdio.Fprintf(fp, "S\n%d %d\n", (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter2)))).Spellname, (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter2)))).Pages)
					continue
				}
			}
		}
	}
	stdio.Fprintf(fp, "$~\n")
	fp.Close()
	stdio.Snprintf(&buf[0], int(64936), "%s%d.obj", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	stdio.Remove(libc.GoString(&buf[0]))
	stdio.Rename(libc.GoString(&cmfname[0]), libc.GoString(&buf[0]))
	if in_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_OBJ) != 0 {
		remove_from_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_OBJ)
		create_world_index(int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number), libc.CString("obj"))
		basic_mud_log(libc.CString("GenOLC: save_objects: Saving objects '%s'"), &buf[0])
	}
	return TRUE
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
	var robj_num int = int(obj.Item_number)
	if obj.Name != nil && obj.Name != (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Name {
		libc.Free(unsafe.Pointer(obj.Name))
	}
	if obj.Description != nil && obj.Description != (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Description {
		libc.Free(unsafe.Pointer(obj.Description))
	}
	if obj.Short_description != nil && obj.Short_description != (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Short_description {
		libc.Free(unsafe.Pointer(obj.Short_description))
	}
	if obj.Action_description != nil && obj.Action_description != (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Action_description {
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
					plist = (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Ex_description
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
	return copy_object_main(to, from, TRUE)
}
func copy_object_preserve(to *obj_data, from *obj_data) int {
	return copy_object_main(to, from, FALSE)
}
func copy_object_main(to *obj_data, from *obj_data, free_object int) int {
	*to = *from
	copy_object_strings(to, from)
	return TRUE
}
func delete_object(rnum obj_rnum) int {
	var (
		i      obj_rnum
		zrnum  zone_rnum
		obj    *obj_data
		tmp    *obj_data
		shop   int
		j      int
		zone   int
		cmd_no int
	)
	if rnum == obj_rnum(-1) || rnum > top_of_objt {
		return -1
	}
	obj = (*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(rnum)))
	zrnum = real_zone_by_thing(room_vnum(GET_OBJ_VNUM(obj)))
	htree_del(obj_htree, int64(obj.Item_number))
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
	if (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(rnum)))).Number != 0 {
		panic("assert failed")
	}
	for tmp = object_list; tmp != nil; tmp = tmp.Next {
		tmp.Item_number -= obj_vnum(libc.BoolToInt(tmp.Item_number > obj_vnum(rnum)))
	}
	for i = rnum; i < top_of_objt; i++ {
		*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i))) = *(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i+1)))
		*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i))) = *(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i+1)))
		(*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Item_number = obj_vnum(i)
	}
	top_of_objt--
	obj_index = (*index_data)(libc.Realloc(unsafe.Pointer(obj_index), int(top_of_objt*obj_rnum(unsafe.Sizeof(index_data{}))+1)))
	obj_proto = (*obj_data)(libc.Realloc(unsafe.Pointer(obj_proto), int(top_of_objt*obj_rnum(unsafe.Sizeof(obj_data{}))+1)))
	for shop = 0; shop <= top_shop; shop++ {
		for j = 0; (*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(j)))) != obj_vnum(-1); j++ {
			*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(j))) -= obj_vnum(libc.BoolToInt((*(*obj_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(j)))) > obj_vnum(rnum)))
		}
	}
	for zone = 0; zone <= int(top_of_zone_table); zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			switch (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command {
			case 'P':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3 == vnum(rnum) {
					delete_zone_command((*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone))), cmd_no)
				} else {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3 -= vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg3 > vnum(rnum)))
				}
			case 'O':
				fallthrough
			case 'G':
				fallthrough
			case 'E':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(rnum) {
					delete_zone_command((*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone))), cmd_no)
				} else {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 -= vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 > vnum(rnum)))
				}
			case 'R':
				if (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 == vnum(rnum) {
					delete_zone_command((*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone))), cmd_no)
				} else {
					(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 -= vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg2 > vnum(rnum)))
				}
			}
		}
	}
	save_objects(zrnum)
	return int(rnum)
}
