package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func add_mobile(mob *char_data, vnum mob_vnum) int {
	var (
		rnum     int
		i        int
		found    int = FALSE
		shop     int
		guild    int
		cmd_no   int
		zone     zone_rnum
		live_mob *char_data
	)
	if (func() int {
		rnum = int(real_mobile(vnum))
		return rnum
	}()) != int(-1) {
		copy_mobile((*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(rnum))), mob)
		for live_mob = character_list; live_mob != nil; live_mob = live_mob.Next {
			if rnum == int(live_mob.Nr) {
				update_mobile_strings(live_mob, (*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(rnum))))
			}
		}
		add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(real_zone_by_thing(room_vnum(vnum)))))).Number, SL_MOB)
		basic_mud_log(libc.CString("GenOLC: add_mobile: Updated existing mobile #%d."), vnum)
		return rnum
	}
	mob_proto = (*char_data)(libc.Realloc(unsafe.Pointer(mob_proto), int(top_of_mobt*mob_rnum(unsafe.Sizeof(char_data{}))+2)))
	mob_index = (*index_data)(libc.Realloc(unsafe.Pointer(mob_index), int(top_of_mobt*mob_rnum(unsafe.Sizeof(index_data{}))+2)))
	top_of_mobt++
	for i = int(top_of_mobt); i > 0; i-- {
		if vnum > (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i-1)))).Vnum {
			*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i))) = *mob
			(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Nr = mob_rnum(i)
			copy_mobile_strings((*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i))), mob)
			(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum = vnum
			(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Number = 0
			(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Func = nil
			found = i
			break
		}
		*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i))) = *(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i-1)))
		*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i))) = *(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i-1)))
		(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Nr++
		htree_add(mob_htree, int64((*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum), int64(i))
	}
	if found == 0 {
		*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*0)) = *mob
		(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*0))).Nr = 0
		copy_mobile_strings((*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*0)), mob)
		(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*0))).Vnum = vnum
		(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*0))).Number = 0
		(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*0))).Func = nil
		htree_add(mob_htree, int64((*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*0))).Vnum), 0)
	}
	basic_mud_log(libc.CString("GenOLC: add_mobile: Added mobile %d at index #%d."), vnum, found)
	for live_mob = character_list; live_mob != nil; live_mob = live_mob.Next {
		live_mob.Nr += mob_rnum(libc.BoolToInt(live_mob.Nr != mob_rnum(-1) && live_mob.Nr >= mob_rnum(found)))
	}
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			if int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) == 'M' {
				(*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 >= vnum(found)))
			}
		}
	}
	if shop_index != nil {
		for shop = 0; shop <= top_shop; shop++ {
			(*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Keeper += mob_rnum(libc.BoolToInt((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Keeper != mob_rnum(-1) && (*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(shop)))).Keeper >= mob_rnum(found)))
		}
	}
	if guild_index != nil {
		for guild = 0; guild <= top_guild; guild++ {
			(*(*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(guild)))).Gm += mob_rnum(libc.BoolToInt((*(*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(guild)))).Gm != mob_rnum(-1) && (*(*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(guild)))).Gm >= mob_rnum(found)))
		}
	}
	add_to_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(real_zone_by_thing(room_vnum(vnum)))))).Number, SL_MOB)
	return found
}
func copy_mobile(to *char_data, from *char_data) int {
	free_mobile_strings(to)
	*to = *from
	check_mobile_strings(from)
	copy_mobile_strings(to, from)
	return TRUE
}
func extract_mobile_all(vnum mob_vnum) {
	var (
		next *char_data
		ch   *char_data
	)
	for ch = character_list; ch != nil; ch = next {
		next = ch.Next
		if GET_MOB_VNUM(ch) == vnum {
			extract_char(ch)
		}
	}
}
func delete_mobile(refpt mob_rnum) int {
	var (
		live_mob *char_data
		counter  int
		cmd_no   int
		vnum     mob_vnum
		zone     zone_rnum
	)
	if refpt < 0 || refpt > top_of_mobt {
		basic_mud_log(libc.CString("SYSERR: GenOLC: delete_mobile: Invalid rnum %d."), refpt)
		return -1
	}
	vnum = (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(refpt)))).Vnum
	extract_mobile_all(vnum)
	for counter = int(refpt); counter < int(top_of_mobt); counter++ {
		*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(counter))) = *(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(counter+1)))
		*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(counter))) = *(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(counter+1)))
		(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(counter)))).Nr = mob_rnum(counter)
	}
	top_of_mobt--
	mob_index = (*index_data)(libc.Realloc(unsafe.Pointer(mob_index), int(top_of_mobt*mob_rnum(unsafe.Sizeof(index_data{}))+1)))
	mob_proto = (*char_data)(libc.Realloc(unsafe.Pointer(mob_proto), int(top_of_mobt*mob_rnum(unsafe.Sizeof(char_data{}))+1)))
	for live_mob = character_list; live_mob != nil; live_mob = live_mob.Next {
		live_mob.Nr -= mob_rnum(libc.BoolToInt(live_mob.Nr >= refpt))
	}
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) != 'S'; cmd_no++ {
			if int((*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Command) == 'M' && (*(*reset_com)(unsafe.Add(unsafe.Pointer((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone)))).Cmd), unsafe.Sizeof(reset_com{})*uintptr(cmd_no)))).Arg1 == vnum(refpt) {
				delete_zone_command((*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone))), cmd_no)
			}
		}
	}
	if shop_index != nil {
		for counter = 0; counter <= top_shop; counter++ {
			if (*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(counter)))).Keeper == refpt {
				(*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(counter)))).Keeper = -1
			}
		}
	}
	if guild_index != nil {
		for counter = 0; counter <= top_guild; counter++ {
			if (*(*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(counter)))).Gm == refpt {
				(*(*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(counter)))).Gm = -1
			}
		}
	}
	save_mobiles(real_zone_by_thing(room_vnum(vnum)))
	return int(refpt)
}
func copy_mobile_strings(t *char_data, f *char_data) int {
	if f.Name != nil {
		t.Name = C.strdup(f.Name)
	}
	if f.Title != nil {
		t.Title = C.strdup(f.Title)
	}
	if f.Short_descr != nil {
		t.Short_descr = C.strdup(f.Short_descr)
	}
	if f.Long_descr != nil {
		t.Long_descr = C.strdup(f.Long_descr)
	}
	if f.Description != nil {
		t.Description = C.strdup(f.Description)
	}
	return TRUE
}
func update_mobile_strings(t *char_data, f *char_data) int {
	if f.Name != nil {
		t.Name = f.Name
	}
	if f.Title != nil {
		t.Title = f.Title
	}
	if f.Short_descr != nil {
		t.Short_descr = f.Short_descr
	}
	if f.Long_descr != nil {
		t.Long_descr = f.Long_descr
	}
	if f.Description != nil {
		t.Description = f.Description
	}
	return TRUE
}
func free_mobile_strings(mob *char_data) int {
	if mob.Name != nil {
		libc.Free(unsafe.Pointer(mob.Name))
	}
	if mob.Title != nil {
		libc.Free(unsafe.Pointer(mob.Title))
	}
	if mob.Short_descr != nil {
		libc.Free(unsafe.Pointer(mob.Short_descr))
	}
	if mob.Long_descr != nil {
		libc.Free(unsafe.Pointer(mob.Long_descr))
	}
	if mob.Description != nil {
		libc.Free(unsafe.Pointer(mob.Description))
	}
	return TRUE
}
func free_mobile(mob *char_data) int {
	var i mob_rnum
	if mob == nil {
		return FALSE
	}
	if (func() mob_rnum {
		i = mob.Nr
		return i
	}()) == mob_rnum(-1) {
		free_mobile_strings(mob)
		free_proto_script(unsafe.Pointer(mob), MOB_TRIGGER)
	} else {
		if mob.Name != nil && mob.Name != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Name {
			libc.Free(unsafe.Pointer(mob.Name))
		}
		if mob.Title != nil && mob.Title != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Title {
			libc.Free(unsafe.Pointer(mob.Title))
		}
		if mob.Short_descr != nil && mob.Short_descr != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Short_descr {
			libc.Free(unsafe.Pointer(mob.Short_descr))
		}
		if mob.Long_descr != nil && mob.Long_descr != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Long_descr {
			libc.Free(unsafe.Pointer(mob.Long_descr))
		}
		if mob.Description != nil && mob.Description != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Description {
			libc.Free(unsafe.Pointer(mob.Description))
		}
		if mob.Proto_script != nil && mob.Proto_script != (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Proto_script {
			free_proto_script(unsafe.Pointer(mob), MOB_TRIGGER)
		}
	}
	for mob.Affected != nil {
		affect_remove(mob, mob.Affected)
	}
	if mob.Script != nil {
		extract_script(unsafe.Pointer(mob), MOB_TRIGGER)
	}
	libc.Free(unsafe.Pointer(mob))
	return TRUE
}
func save_mobiles(zone_num zone_rnum) int {
	var (
		mobfd     *C.FILE
		i         room_vnum
		rmob      mob_rnum
		written   int
		mobfname  [64]byte
		usedfname [64]byte
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_mobiles: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return FALSE
	}
	stdio.Snprintf(&mobfname[0], int(64), "%s%d.new", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	if (func() *C.FILE {
		mobfd = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&mobfname[0]), "w")))
		return mobfd
	}()) == nil {
		mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("SYSERR: GenOLC: Cannot open mob file for writing."))
		return FALSE
	}
	for i = genolc_zone_bottom(zone_num); i <= (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Top; i++ {
		if (func() mob_rnum {
			rmob = real_mobile(mob_vnum(i))
			return rmob
		}()) == mob_rnum(-1) {
			continue
		}
		check_mobile_strings((*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(rmob))))
		if write_mobile_record(mob_vnum(i), (*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(rmob))), mobfd) < 0 {
			basic_mud_log(libc.CString("SYSERR: GenOLC: Error writing mobile #%d."), i)
		}
	}
	fputs(libc.CString("$\n"), mobfd)
	written = ftell(mobfd)
	C.fclose(mobfd)
	stdio.Snprintf(&usedfname[0], int(64), "%s%d.mob", LIB_WORLD, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number)
	stdio.Remove(libc.GoString(&usedfname[0]))
	stdio.Rename(libc.GoString(&mobfname[0]), libc.GoString(&usedfname[0]))
	if in_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_MOB) != 0 {
		remove_from_save_list((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number, SL_MOB)
		create_world_index(int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zone_num)))).Number), libc.CString("mob"))
		basic_mud_log(libc.CString("GenOLC: save_mobiles: Saving mobiles '%s'"), &usedfname[0])
	}
	return written
}
func write_mobile_espec(mvnum mob_vnum, mob *char_data, fd *C.FILE) int {
	var (
		aff *affected_type
		i   int
	)
	if get_size(mob) != race_def_sizetable[mob.Race] {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Size: %d\n", get_size(mob))
	}
	if mob.Mob_specials.Attack_type != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "BareHandAttack: %d\n", mob.Mob_specials.Attack_type)
	}
	if mob.Aff_abils.Str != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Str: %d\n", mob.Aff_abils.Str)
	}
	if mob.Aff_abils.Dex != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Dex: %d\n", mob.Aff_abils.Dex)
	}
	if mob.Aff_abils.Intel != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Int: %d\n", mob.Aff_abils.Intel)
	}
	if mob.Aff_abils.Wis != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Wis: %d\n", mob.Aff_abils.Wis)
	}
	if mob.Aff_abils.Con != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Con: %d\n", mob.Aff_abils.Con)
	}
	if mob.Aff_abils.Cha != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Cha: %d\n", mob.Aff_abils.Cha)
	}
	if (*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(real_mobile(mvnum)))) != mob {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Hit: %lld\nMaxHit: %lld\nMana: %lld\nMaxMana: %lld\nMoves: %lld\nMaxMoves: %lld\n", mob.Hit, mob.Max_hit, mob.Mana, mob.Max_mana, mob.Move, mob.Max_move)
		for aff = mob.Affected; aff != nil; aff = aff.Next {
			if aff.Type != 0 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Affect: %d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
			}
		}
		for aff = mob.Affectedv; aff != nil; aff = aff.Next {
			if aff.Type != 0 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "AffectV: %d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
			}
		}
	}
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		if (mob.Feats[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Feat: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if (mob.Skills[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Skill: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		if (mob.Skillmods[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "SkillMod: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i < NUM_CLASSES; i++ {
		if (mob.Chclasses[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "Class: %d %d\n", i, mob.Chclasses[i])
		}
		if (mob.Epicclasses[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "EpicClass: %d %d\n", i, mob.Epicclasses[i])
		}
	}
	fputs(libc.CString("E\n"), fd)
	return TRUE
}
func write_mobile_record(mvnum mob_vnum, mob *char_data, fd *C.FILE) int {
	var (
		ldesc [64936]byte
		ddesc [64936]byte
		fbuf1 [64936]byte
		fbuf2 [64936]byte
		fbuf3 [64936]byte
		fbuf4 [64936]byte
		abuf1 [64936]byte
		abuf2 [64936]byte
		abuf3 [64936]byte
		abuf4 [64936]byte
	)
	ldesc[int(MAX_STRING_LENGTH-1)] = '\x00'
	ddesc[int(MAX_STRING_LENGTH-1)] = '\x00'
	strip_cr(C.strncpy(&ldesc[0], mob.Long_descr, uint64(int(MAX_STRING_LENGTH-1))))
	strip_cr(C.strncpy(&ddesc[0], mob.Description, uint64(int(MAX_STRING_LENGTH-1))))
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "#%d\n%s%c\n%s%c\n%s%c\n%s%c\n", mvnum, mob.Name, STRING_TERMINATOR, mob.Short_descr, STRING_TERMINATOR, &ldesc[0], STRING_TERMINATOR, &ddesc[0], STRING_TERMINATOR)
	sprintascii(&fbuf1[0], mob.Act[0])
	sprintascii(&fbuf2[0], mob.Act[1])
	sprintascii(&fbuf3[0], mob.Act[2])
	sprintascii(&fbuf4[0], mob.Act[3])
	sprintascii(&abuf1[0], bitvector_t(mob.Affected_by[0]))
	sprintascii(&abuf2[0], bitvector_t(mob.Affected_by[1]))
	sprintascii(&abuf3[0], bitvector_t(mob.Affected_by[2]))
	sprintascii(&abuf4[0], bitvector_t(mob.Affected_by[3]))
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "%s %s %s %s %s %s %s %s %d E\n%d %d %d %lldd%lld+%lld %dd%d+%d\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0], &abuf1[0], &abuf2[0], &abuf3[0], &abuf4[0], mob.Alignment, mob.Race_level, mob.Accuracy_mod, 10-mob.Armor/10, mob.Hit, mob.Mana, mob.Move, mob.Mob_specials.Damnodice, mob.Mob_specials.Damsizedice, mob.Damage_mod)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fd)), "%d 0 %d %d\n%d %d %d\n", mob.Gold, mob.Race, mob.Chclass, mob.Position, mob.Mob_specials.Default_pos, mob.Sex)
	if write_mobile_espec(mvnum, mob, fd) < 0 {
		basic_mud_log(libc.CString("SYSERR: GenOLC: Error writing E-specs for mobile #%d."), mvnum)
	}
	script_save_to_disk(fd, unsafe.Pointer(mob), MOB_TRIGGER)
	return TRUE
}
func check_mobile_strings(mob *char_data) {
	var mvnum mob_vnum = (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(mob.Nr)))).Vnum
	check_mobile_string(mvnum, &mob.Long_descr, libc.CString("long description"))
	check_mobile_string(mvnum, &mob.Description, libc.CString("detailed description"))
	check_mobile_string(mvnum, &mob.Name, libc.CString("alias list"))
	check_mobile_string(mvnum, &mob.Short_descr, libc.CString("short description"))
}
func check_mobile_string(i mob_vnum, string_ **byte, dscr *byte) {
	if *string_ == nil || **string_ == '\x00' {
		var smbuf [128]byte
		stdio.Sprintf(&smbuf[0], "GenOLC: Mob #%d has an invalid %s.", i, dscr)
		mudlog(BRF, ADMLVL_GOD, TRUE, &smbuf[0])
		if *string_ != nil {
			libc.Free(unsafe.Pointer(*string_))
		}
		*string_ = C.strdup(libc.CString("An undefined string."))
	}
}
