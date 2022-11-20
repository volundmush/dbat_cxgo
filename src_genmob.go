package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unsafe"
)

func add_mobile(mob *char_data, vnum int) int {
	basic_mud_log(libc.CString("REIMPLEMENT THIS!"))
	os.Exit(-1)
	return 0
}
func copy_mobile(to *char_data, from *char_data) {
	free_mobile_strings(to)
	*to = *from
	check_mobile_strings(from)
	copy_mobile_strings(to, from)
}
func extract_mobile_all(vnum int) {
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
func delete_mobile(refpt int) int {
	var (
		live_mob *char_data
		counter  int
		cmd_no   int
		vnum     int
		zone     int
	)
	if refpt < 0 || refpt > top_of_mobt {
		basic_mud_log(libc.CString("SYSERR: GenOLC: delete_mobile: Invalid rnum %d."), refpt)
		return -1
	}
	vnum = mob_index[refpt].Vnum
	extract_mobile_all(vnum)
	for counter = refpt; counter < top_of_mobt; counter++ {
		mob_index[counter] = mob_index[counter+1]
		mob_proto[counter] = mob_proto[counter+1]
		mob_proto[counter].Nr = counter
	}
	top_of_mobt--
	// todo: fix this
	//mob_index = []index_data((*index_data)(libc.Realloc(unsafe.Pointer(&mob_index[0]), top_of_mobt*int(unsafe.Sizeof(index_data{}))+1)))
	//mob_proto = []char_data((*char_data)(libc.Realloc(unsafe.Pointer(&mob_proto[0]), top_of_mobt*int(unsafe.Sizeof(char_data{}))+1)))
	for live_mob = character_list; live_mob != nil; live_mob = live_mob.Next {
		live_mob.Nr -= int(libc.BoolToInt(live_mob.Nr >= refpt))
	}
	for zone = 0; zone <= top_of_zone_table; zone++ {
		for cmd_no = 0; int(zone_table[zone].Cmd[cmd_no].Command) != 'S'; cmd_no++ {
			if int(zone_table[zone].Cmd[cmd_no].Command) == 'M' && zone_table[zone].Cmd[cmd_no].Arg1 == refpt {
				delete_zone_command(&zone_table[zone], cmd_no)
			}
		}
	}
	if shop_index != nil {
		for counter = 0; counter <= top_shop; counter++ {
			if shop_index[counter].Keeper == refpt {
				shop_index[counter].Keeper = -1
			}
		}
	}
	if guild_index != nil {
		for counter = 0; counter <= top_guild; counter++ {
			if guild_index[counter].Gm == refpt {
				guild_index[counter].Gm = -1
			}
		}
	}
	save_mobiles(real_zone_by_thing(vnum))
	return refpt
}
func copy_mobile_strings(t *char_data, f *char_data) {
	if f.Name != nil {
		t.Name = libc.StrDup(f.Name)
	}
	if f.Title != nil {
		t.Title = libc.StrDup(f.Title)
	}
	if f.Short_descr != nil {
		t.Short_descr = libc.StrDup(f.Short_descr)
	}
	if f.Long_descr != nil {
		t.Long_descr = libc.StrDup(f.Long_descr)
	}
	if f.Description != nil {
		t.Description = libc.StrDup(f.Description)
	}
}
func update_mobile_strings(t *char_data, f *char_data) {
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
}
func free_mobile_strings(mob *char_data) {
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
}
func free_mobile(mob *char_data) bool {
	var i int
	if mob == nil {
		return false
	}
	if (func() int {
		i = mob.Nr
		return i
	}()) == int(-1) {
		free_mobile_strings(mob)
		free_proto_script(unsafe.Pointer(mob), MOB_TRIGGER)
	} else {
		if mob.Name != nil && mob.Name != mob_proto[i].Name {
			libc.Free(unsafe.Pointer(mob.Name))
		}
		if mob.Title != nil && mob.Title != mob_proto[i].Title {
			libc.Free(unsafe.Pointer(mob.Title))
		}
		if mob.Short_descr != nil && mob.Short_descr != mob_proto[i].Short_descr {
			libc.Free(unsafe.Pointer(mob.Short_descr))
		}
		if mob.Long_descr != nil && mob.Long_descr != mob_proto[i].Long_descr {
			libc.Free(unsafe.Pointer(mob.Long_descr))
		}
		if mob.Description != nil && mob.Description != mob_proto[i].Description {
			libc.Free(unsafe.Pointer(mob.Description))
		}
		if mob.Proto_script != nil && mob.Proto_script != mob_proto[i].Proto_script {
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
	return true
}
func save_mobiles(zone_num int) int {
	var (
		mobfd     *stdio.File
		i         int
		rmob      int
		written   int
		mobfname  [64]byte
		usedfname [64]byte
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_mobiles: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return 0
	}
	stdio.Snprintf(&mobfname[0], int(64), "%s%d.new", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		mobfd = stdio.FOpen(libc.GoString(&mobfname[0]), "w")
		return mobfd
	}()) == nil {
		mudlog(BRF, ADMLVL_GOD, 1, libc.CString("SYSERR: GenOLC: Cannot open mob file for writing."))
		return 0
	}
	for i = genolc_zone_bottom(zone_num); i <= zone_table[zone_num].Top; i++ {
		if (func() int {
			rmob = real_mobile(i)
			return rmob
		}()) == int(-1) {
			continue
		}
		check_mobile_strings(&mob_proto[rmob])
		if int(libc.BoolToInt(write_mobile_record(i, &mob_proto[rmob], mobfd))) < 0 {
			basic_mud_log(libc.CString("SYSERR: GenOLC: Error writing mobile #%d."), i)
		}
	}
	mobfd.PutS(libc.CString("$\n"))
	written = int(mobfd.Tell())
	mobfd.Close()
	stdio.Snprintf(&usedfname[0], int(64), "%s%d.mob", LIB_WORLD, zone_table[zone_num].Number)
	stdio.Remove(libc.GoString(&usedfname[0]))
	stdio.Rename(libc.GoString(&mobfname[0]), libc.GoString(&usedfname[0]))
	if in_save_list(zone_table[zone_num].Number, SL_MOB) {
		remove_from_save_list(zone_table[zone_num].Number, SL_MOB)
		create_world_index(zone_table[zone_num].Number, libc.CString("mob"))
		basic_mud_log(libc.CString("GenOLC: save_mobiles: Saving mobiles '%s'"), &usedfname[0])
	}
	return written
}
func write_mobile_espec(mvnum int, mob *char_data, fd *stdio.File) bool {
	var (
		aff *affected_type
		i   int
	)
	if get_size(mob) != race_def_sizetable[mob.Race] {
		stdio.Fprintf(fd, "Size: %d\n", get_size(mob))
	}
	if int(mob.Mob_specials.Attack_type) != 0 {
		stdio.Fprintf(fd, "BareHandAttack: %d\n", mob.Mob_specials.Attack_type)
	}
	if int(mob.Aff_abils.Str) != 0 {
		stdio.Fprintf(fd, "Str: %d\n", mob.Aff_abils.Str)
	}
	if int(mob.Aff_abils.Dex) != 0 {
		stdio.Fprintf(fd, "Dex: %d\n", mob.Aff_abils.Dex)
	}
	if int(mob.Aff_abils.Intel) != 0 {
		stdio.Fprintf(fd, "Int: %d\n", mob.Aff_abils.Intel)
	}
	if int(mob.Aff_abils.Wis) != 0 {
		stdio.Fprintf(fd, "Wis: %d\n", mob.Aff_abils.Wis)
	}
	if int(mob.Aff_abils.Con) != 0 {
		stdio.Fprintf(fd, "Con: %d\n", mob.Aff_abils.Con)
	}
	if int(mob.Aff_abils.Cha) != 0 {
		stdio.Fprintf(fd, "Cha: %d\n", mob.Aff_abils.Cha)
	}
	if &mob_proto[real_mobile(mvnum)] != mob {
		stdio.Fprintf(fd, "Hit: %lld\nMaxHit: %lld\nMana: %lld\nMaxMana: %lld\nMoves: %lld\nMaxMoves: %lld\n", mob.Hit, mob.Max_hit, mob.Mana, mob.Max_mana, mob.Move, mob.Max_move)
		for aff = mob.Affected; aff != nil; aff = aff.Next {
			if int(aff.Type) != 0 {
				stdio.Fprintf(fd, "Affect: %d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
			}
		}
		for aff = mob.Affectedv; aff != nil; aff = aff.Next {
			if int(aff.Type) != 0 {
				stdio.Fprintf(fd, "AffectV: %d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
			}
		}
	}
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		if int(mob.Feats[i]) != 0 {
			stdio.Fprintf(fd, "Feat: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if int(mob.Skills[i]) != 0 {
			stdio.Fprintf(fd, "Skill: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		if int(mob.Skillmods[i]) != 0 {
			stdio.Fprintf(fd, "SkillMod: %d %d\n", i, mob.Feats[i])
		}
	}
	for i = 0; i < NUM_CLASSES; i++ {
		if (mob.Chclasses[i]) != 0 {
			stdio.Fprintf(fd, "Class: %d %d\n", i, mob.Chclasses[i])
		}
		if (mob.Epicclasses[i]) != 0 {
			stdio.Fprintf(fd, "EpicClass: %d %d\n", i, mob.Epicclasses[i])
		}
	}
	fd.PutS(libc.CString("E\n"))
	return true
}
func write_mobile_record(mvnum int, mob *char_data, fd *stdio.File) bool {
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
	strip_cr(libc.StrNCpy(&ldesc[0], mob.Long_descr, int(MAX_STRING_LENGTH-1)))
	strip_cr(libc.StrNCpy(&ddesc[0], mob.Description, int(MAX_STRING_LENGTH-1)))
	stdio.Fprintf(fd, "#%d\n%s%c\n%s%c\n%s%c\n%s%c\n", mvnum, mob.Name, STRING_TERMINATOR, mob.Short_descr, STRING_TERMINATOR, &ldesc[0], STRING_TERMINATOR, &ddesc[0], STRING_TERMINATOR)
	sprintascii(&fbuf1[0], mob.Act[0])
	sprintascii(&fbuf2[0], mob.Act[1])
	sprintascii(&fbuf3[0], mob.Act[2])
	sprintascii(&fbuf4[0], mob.Act[3])
	sprintascii(&abuf1[0], mob.Affected_by[0])
	sprintascii(&abuf2[0], mob.Affected_by[1])
	sprintascii(&abuf3[0], mob.Affected_by[2])
	sprintascii(&abuf4[0], mob.Affected_by[3])
	stdio.Fprintf(fd, "%s %s %s %s %s %s %s %s %d E\n%d %d %d %lldd%lld+%lld %dd%d+%d\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0], &abuf1[0], &abuf2[0], &abuf3[0], &abuf4[0], mob.Alignment, mob.Race_level, mob.Accuracy_mod, 10-mob.Armor/10, mob.Hit, mob.Mana, mob.Move, mob.Mob_specials.Damnodice, mob.Mob_specials.Damsizedice, mob.Damage_mod)
	stdio.Fprintf(fd, "%d 0 %d %d\n%d %d %d\n", mob.Gold, mob.Race, mob.Chclass, mob.Position, mob.Mob_specials.Default_pos, mob.Sex)
	if int(libc.BoolToInt(write_mobile_espec(mvnum, mob, fd))) < 0 {
		basic_mud_log(libc.CString("SYSERR: GenOLC: Error writing E-specs for mobile #%d."), mvnum)
	}
	script_save_to_disk(fd, unsafe.Pointer(mob), MOB_TRIGGER)
	return true
}
func check_mobile_strings(mob *char_data) {
	var mvnum int = mob_index[mob.Nr].Vnum
	check_mobile_string(mvnum, &mob.Long_descr, libc.CString("long description"))
	check_mobile_string(mvnum, &mob.Description, libc.CString("detailed description"))
	check_mobile_string(mvnum, &mob.Name, libc.CString("alias list"))
	check_mobile_string(mvnum, &mob.Short_descr, libc.CString("short description"))
}
func check_mobile_string(i int, string_ **byte, dscr *byte) {
	if *string_ == nil || **string_ == '\x00' {
		var smbuf [128]byte
		stdio.Sprintf(&smbuf[0], "GenOLC: Mob #%d has an invalid %s.", i, dscr)
		mudlog(BRF, ADMLVL_GOD, 1, &smbuf[0])
		if *string_ != nil {
			libc.Free(unsafe.Pointer(*string_))
		}
		*string_ = libc.CString("An undefined string.")
	}
}
