package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
	"unsafe"
)

func do_oasis_oedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int = int(-1)
		save     int = 0
		real_num int
		d        *descriptor_data
		buf3     *byte
	)
	_ = buf3
	var buf1 [64936]byte
	var buf2 [64936]byte
	buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		send_to_char(ch, libc.CString("Specify an object VNUM to edit.\r\n"))
		return
	} else if !unicode.IsDigit(rune(buf1[0])) {
		if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) != 0 {
			send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
			return
		}
		save = TRUE
		if is_number(&buf2[0]) != 0 {
			number = libc.Atoi(libc.GoString(&buf2[0]))
		} else if ch.Player_specials.Olc_zone > 0 {
			var zlok zone_rnum
			if (func() zone_rnum {
				zlok = real_zone(zone_vnum(ch.Player_specials.Olc_zone))
				return zlok
			}()) == zone_rnum(-1) {
				number = -1
			} else {
				number = int(genolc_zone_bottom(zlok))
			}
		}
		if number == int(-1) {
			send_to_char(ch, libc.CString("Save which zone?\r\n"))
			return
		}
	}
	if number == int(-1) {
		number = libc.Atoi(libc.GoString(&buf1[0]))
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_OEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That object is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if save != 0 {
		d.Olc.Zone_num = real_zone(zone_vnum(number))
	} else {
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(number))
	}
	if d.Olc.Zone_num == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for that number!\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if can_edit_zone(ch, d.Olc.Zone_num) == 0 {
		send_cannot_edit(ch, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all objects in zone %d.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves object info for zone %d."), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		save_objects(d.Olc.Zone_num)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() int {
		real_num = int(real_object(obj_vnum(number)))
		return real_num
	}()) != int(-1) {
		oedit_setup_existing(d, real_num)
	} else {
		oedit_setup_new(d)
	}
	oedit_disp_menu(d)
	d.Connected = CON_OEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
}
func oedit_setup_new(d *descriptor_data) {
	d.Olc.Obj = new(obj_data)
	clear_object(d.Olc.Obj)
	d.Olc.Obj.Name = libc.CString("unfinished object")
	d.Olc.Obj.Description = libc.CString("An unfinished object is lying here.")
	d.Olc.Obj.Short_description = libc.CString("an unfinished object")
	d.Olc.Obj.Wear_flags[int(ITEM_WEAR_TAKE/32)] |= 1 << (int(ITEM_WEAR_TAKE % 32))
	d.Olc.Value = 0
	d.Olc.Item_type = OBJ_TRIGGER
	d.Olc.Obj.Type_flag = ITEM_WORN
	d.Olc.Obj.Value[VAL_ALL_HEALTH] = 100
	d.Olc.Obj.Value[VAL_ALL_MAXHEALTH] = 100
	d.Olc.Obj.Value[VAL_ALL_MATERIAL] = MATERIAL_STEEL
	d.Olc.Obj.Size = SIZE_MEDIUM
	d.Olc.Obj.Script = nil
	d.Olc.Obj.Proto_script = func() *trig_proto_list {
		p := &d.Olc.Script
		d.Olc.Script = nil
		return *p
	}()
}
func oedit_setup_existing(d *descriptor_data, real_num int) {
	var obj *obj_data
	obj = new(obj_data)
	copy_object(obj, (*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(real_num))))
	d.Olc.Obj = obj
	d.Olc.Value = 0
	d.Olc.Item_type = OBJ_TRIGGER
	dg_olc_script_copy(d)
	obj.Script = nil
	d.Olc.Obj.Proto_script = nil
}
func oedit_save_internally(d *descriptor_data) {
	var (
		i        int
		robj_num obj_rnum
		dsc      *descriptor_data
		obj      *obj_data
	)
	i = int(libc.BoolToInt(real_object(obj_vnum(d.Olc.Number)) == obj_rnum(-1)))
	if (func() obj_rnum {
		robj_num = add_object(d.Olc.Obj, obj_vnum(d.Olc.Number))
		return robj_num
	}()) == obj_rnum(-1) {
		basic_mud_log(libc.CString("oedit_save_internally: add_object failed."))
		return
	}
	if (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Proto_script != nil && (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Proto_script != d.Olc.Script {
		free_proto_script(unsafe.Pointer((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))), OBJ_TRIGGER)
	}
	(*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))).Proto_script = d.Olc.Script
	for obj = object_list; obj != nil; obj = obj.Next {
		if obj.Item_number != obj_vnum(robj_num) {
			continue
		}
		if obj.Script != nil {
			extract_script(unsafe.Pointer(obj), OBJ_TRIGGER)
		}
		free_proto_script(unsafe.Pointer(obj), OBJ_TRIGGER)
		copy_proto_script(unsafe.Pointer((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj_num)))), unsafe.Pointer(obj), OBJ_TRIGGER)
		assign_triggers(unsafe.Pointer(obj), OBJ_TRIGGER)
	}
	if i == 0 {
		return
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		if dsc.Connected == CON_SEDIT {
			for i = 0; (*(*obj_vnum)(unsafe.Add(unsafe.Pointer(dsc.Olc.Shop.Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i)))) != obj_vnum(-1); i++ {
				if (*(*obj_vnum)(unsafe.Add(unsafe.Pointer(dsc.Olc.Shop.Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i)))) >= obj_vnum(robj_num) {
					(*(*obj_vnum)(unsafe.Add(unsafe.Pointer(dsc.Olc.Shop.Producing), unsafe.Sizeof(obj_vnum(0))*uintptr(i))))++
				}
			}
		}
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		if dsc.Connected == CON_ZEDIT {
			for i = 0; int((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Command) != 'S'; i++ {
				switch (*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Command {
				case 'P':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg3 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg3 >= vnum(robj_num)))
					fallthrough
				case 'E':
					fallthrough
				case 'G':
					fallthrough
				case 'O':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg1 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg1 >= vnum(robj_num)))
				case 'R':
					(*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg2 += vnum(libc.BoolToInt((*(*reset_com)(unsafe.Add(unsafe.Pointer(dsc.Olc.Zone.Cmd), unsafe.Sizeof(reset_com{})*uintptr(i)))).Arg2 >= vnum(robj_num)))
				default:
				}
			}
		}
	}
}
func oedit_save_to_disk(zone_num int) {
	save_objects(zone_rnum(zone_num))
}
func oedit_disp_container_flags_menu(d *descriptor_data) {
	var bits [64936]byte
	clear_screen(d)
	sprintbit(bitvector_t(int32(d.Olc.Obj.Value[1])), container_bits[:], &bits[0], uint64(64936))
	write_to_output(d, libc.CString("@g1@n) CLOSEABLE\r\n@g2@n) PICKPROOF\r\n@g3@n) CLOSED\r\n@g4@n) LOCKED\r\nContainer flags: @c%s@n\r\nEnter flag, 0 to quit : "), &bits[0])
}
func oedit_disp_extradesc_menu(d *descriptor_data) {
	var extra_desc *extra_descr_data = d.Olc.Desc
	clear_screen(d)
	write_to_output(d, libc.CString("Extra desc menu\r\n@g1@n) Keyword: @y%s@n\r\n@g2@n) Description:\r\n@y%s@n\r\n@g3@n) Goto next description: %s\r\n@g0@n) Quit\r\nEnter choice : "), func() *byte {
		if extra_desc.Keyword != nil && *extra_desc.Keyword != 0 {
			return extra_desc.Keyword
		}
		return libc.CString("<NONE>")
	}(), func() *byte {
		if extra_desc.Description != nil && *extra_desc.Description != 0 {
			return extra_desc.Description
		}
		return libc.CString("<NONE>")
	}(), func() string {
		if extra_desc.Next == nil {
			return "Not set."
		}
		return "Set."
	}())
	d.Olc.Mode = OEDIT_EXTRADESC_MENU
}
func oedit_disp_prompt_apply_menu(d *descriptor_data) {
	var (
		apply_buf [64936]byte
		counter   int
	)
	clear_screen(d)
	for counter = 0; counter < MAX_OBJ_AFFECT; counter++ {
		if d.Olc.Obj.Affected[counter].Modifier != 0 {
			sprinttype(d.Olc.Obj.Affected[counter].Location, apply_types[:], &apply_buf[0], uint64(64936))
			write_to_output(d, libc.CString(" @g%d@n) %+d to @b%s@n"), counter+1, d.Olc.Obj.Affected[counter].Modifier, &apply_buf[0])
			switch d.Olc.Obj.Affected[counter].Location {
			case APPLY_FEAT:
				write_to_output(d, libc.CString(" (%s)"), feat_list[d.Olc.Obj.Affected[counter].Specific].Name)
			case APPLY_SKILL:
				write_to_output(d, libc.CString(" (%s)"), spell_info[d.Olc.Obj.Affected[counter].Specific].Name)
			}
			write_to_output(d, libc.CString("\r\n"))
		} else {
			write_to_output(d, libc.CString(" @g%d@n) None.\r\n"), counter+1)
		}
	}
	write_to_output(d, libc.CString("\r\nEnter affection to modify (0 to quit) : "))
	d.Olc.Mode = OEDIT_PROMPT_APPLY
}
func oedit_disp_prompt_spellbook_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < SPELLBOOK_SIZE; counter++ {
		if d.Olc.Obj.Sbinfo != nil && (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter)))).Spellname != 0 && (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter)))).Spellname < SPELL_SENSU {
			write_to_output(d, libc.CString(" @g%3d@n) %-20.20s %s"), counter+1, spell_info[(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter)))).Spellname].Name, func() string {
				if (func() int {
					p := &columns
					*p++
					return *p
				}() % 3) == 0 {
					return "\r\n"
				}
				return ""
			}())
		} else {
			write_to_output(d, libc.CString(" @g%3d@n) None.%s"), counter+1, func() string {
				if (func() int {
					p := &columns
					*p++
					return *p
				}() % 3) == 0 {
					return "\r\n"
				}
				return ""
			}())
		}
	}
	write_to_output(d, libc.CString("\r\nEnter spell to modify (0 to quit) : "))
	d.Olc.Mode = OEDIT_PROMPT_SPELLBOOK
}
func oedit_disp_spellbook_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < SKILL_TABLE_SIZE; counter++ {
		if spell_info[counter].Skilltype == (1 << 0) {
			write_to_output(d, libc.CString("@g%3d@n) @y%-20.20s@n%s"), counter, spell_info[counter].Name, func() string {
				if (func() int {
					p := &columns
					*p++
					return *p
				}() % 3) == 0 {
					return "\r\n"
				}
				return ""
			}())
		}
	}
	write_to_output(d, libc.CString("@n\r\nEnter spell number (0 is no spell) : "))
	d.Olc.Mode = OEDIT_SPELLBOOK
}
func oedit_disp_apply_spec_menu(d *descriptor_data) {
	var buf *byte
	switch d.Olc.Obj.Affected[d.Olc.Value].Location {
	case APPLY_FEAT:
		buf = libc.CString("What feat should be modified : ")
	case APPLY_SKILL:
		buf = libc.CString("What skill should be modified : ")
	default:
		oedit_disp_prompt_apply_menu(d)
		return
	}
	write_to_output(d, libc.CString("\r\n%s"), buf)
	d.Olc.Mode = OEDIT_APPLYSPEC
}
func oedit_liquid_type(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_LIQ_TYPES; counter++ {
		write_to_output(d, libc.CString(" @g%2d@n) @y%-20.20s@n%s"), counter, drinks[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\n@nEnter drink type : "))
	d.Olc.Mode = OEDIT_VALUE_3
}
func oedit_disp_apply_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_APPLIES; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, apply_types[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter apply type (0 is no apply) : "))
	d.Olc.Mode = OEDIT_APPLY
}
func oedit_disp_crittype_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter <= CRIT_X4; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, crit_type[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter critical type : "))
}
func oedit_disp_weapon_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ATTACK_TYPES; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, attack_hit_text[counter].Singular, func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter weapon type : "))
}
func oedit_disp_armor_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter <= MAX_ARMOR_TYPES; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, armor_type[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter armor proficiency type : "))
}
func oedit_disp_spells_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < SKILL_TABLE_SIZE; counter++ {
		if (skill_type(counter) & (1 << 0)) != 0 {
			write_to_output(d, libc.CString("@g%2d@n) @y%-20.20s@n%s"), counter, spell_info[counter].Name, func() string {
				if (func() int {
					p := &columns
					*p++
					return *p
				}() % 3) == 0 {
					return "\r\n"
				}
				return ""
			}())
		}
	}
	write_to_output(d, libc.CString("\r\n@nEnter spell choice (-1 for none) : "))
}
func oedit_disp_material_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_MATERIALS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter, material_names[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\n@nEnter material type : "))
}
func oedit_disp_val1_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	d.Olc.Mode = OEDIT_VALUE_1
	switch d.Olc.Obj.Type_flag {
	case ITEM_LIGHT:
		oedit_disp_val3_menu(d)
	case ITEM_SCROLL:
		fallthrough
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		fallthrough
	case ITEM_POTION:
		write_to_output(d, libc.CString("Spell level : "))
	case ITEM_WEAPON:
		for counter = 0; counter <= MAX_WEAPON_TYPES; counter++ {
			write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, weapon_type[counter], func() string {
				if (func() int {
					p := &columns
					*p++
					return *p
				}() % 3) == 0 {
					return "\r\n"
				}
				return ""
			}())
		}
		write_to_output(d, libc.CString("\r\nEnter the weapon type for determining proficiencies: \r\n"))
	case ITEM_ARMOR:
		write_to_output(d, libc.CString("Apply to AC : "))
	case ITEM_CONTAINER:
		write_to_output(d, libc.CString("Max weight to contain (-1 for unlimited) : "))
	case ITEM_DRINKCON:
		fallthrough
	case ITEM_FOUNTAIN:
		write_to_output(d, libc.CString("Max drink units (-1 for unlimited) : "))
	case ITEM_FOOD:
		write_to_output(d, libc.CString("Hours to fill stomach : "))
	case ITEM_MONEY:
		write_to_output(d, libc.CString("Number of zenni : "))
	case ITEM_NOTE:
	case ITEM_VEHICLE:
		write_to_output(d, libc.CString("Enter room vnum of vehicle interior : "))
	case ITEM_HATCH:
		write_to_output(d, libc.CString("Enter vnum of the vehicle this hatch belongs to : "))
	case ITEM_WINDOW:
		write_to_output(d, libc.CString("Enter vnum of the vehicle this window belongs to, or -1 to specify the viewport room : "))
	case ITEM_CONTROL:
		write_to_output(d, libc.CString("Enter vnum of the vehicle these controls belong to : "))
	case ITEM_PORTAL:
		write_to_output(d, libc.CString("Which room number is the destination? : "))
	case ITEM_BOARD:
		write_to_output(d, libc.CString("Enter the minimum admin level to read this board (0 for mortals) : "))
	default:
		oedit_disp_val5_menu(d)
	}
}
func oedit_disp_val2_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_2
	switch d.Olc.Obj.Type_flag {
	case ITEM_SCROLL:
		fallthrough
	case ITEM_POTION:
		oedit_disp_spells_menu(d)
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		write_to_output(d, libc.CString("Max number of charges : "))
	case ITEM_WEAPON:
		write_to_output(d, libc.CString("Number of damage dice : "))
	case ITEM_ARMOR:
		oedit_disp_armor_menu(d)
	case ITEM_FOOD:
		oedit_disp_val4_menu(d)
	case ITEM_CONTROL:
		write_to_output(d, libc.CString("Enter Engine Type ( 1, 2, 3) : "))
	case ITEM_CONTAINER:
		fallthrough
	case ITEM_VEHICLE:
		fallthrough
	case ITEM_HATCH:
		fallthrough
	case ITEM_WINDOW:
		fallthrough
	case ITEM_PORTAL:
		oedit_disp_container_flags_menu(d)
	case ITEM_DRINKCON:
		fallthrough
	case ITEM_FOUNTAIN:
		write_to_output(d, libc.CString("Initial drink units : "))
	case ITEM_BOARD:
		write_to_output(d, libc.CString("Minimum admin level to write (0 for mortals) : "))
	default:
		oedit_disp_val5_menu(d)
	}
}
func oedit_disp_val3_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_3
	switch d.Olc.Obj.Type_flag {
	case ITEM_LIGHT:
		write_to_output(d, libc.CString("Number of hours (0 = burnt, -1 is infinite) : "))
		break
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		write_to_output(d, libc.CString("Number of charges remaining : "))
	case ITEM_WEAPON:
		write_to_output(d, libc.CString("Size of damage dice : "))
	case ITEM_ARMOR:
		write_to_output(d, libc.CString("Max dex bonus : "))
	case ITEM_CONTAINER:
		write_to_output(d, libc.CString("Vnum of key to open container (-1 for no key) : "))
	case ITEM_DRINKCON:
		fallthrough
	case ITEM_FOUNTAIN:
		oedit_liquid_type(d)
	case ITEM_VEHICLE:
		write_to_output(d, libc.CString("Vnum of key to unlock vehicle (-1 for no key) : "))
	case ITEM_HATCH:
		write_to_output(d, libc.CString("Vnum of key to unlock hatch (-1 for no key) : "))
	case ITEM_WINDOW:
		write_to_output(d, libc.CString("Vnum of key to unlock window (-1 for no key) : "))
	case ITEM_PORTAL:
		write_to_output(d, libc.CString("Vnum of the key to unlock portal (-1 for no key) : "))
	case ITEM_BOARD:
		write_to_output(d, libc.CString("Minimum admin level to remove messages (0 for mortals) : "))
	default:
		oedit_disp_val5_menu(d)
	}
}
func oedit_disp_val4_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_4
	switch d.Olc.Obj.Type_flag {
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		oedit_disp_spells_menu(d)
	case ITEM_WEAPON:
		oedit_disp_weapon_menu(d)
	case ITEM_ARMOR:
		write_to_output(d, libc.CString("Armor check penalty : "))
	case ITEM_DRINKCON:
		fallthrough
	case ITEM_FOUNTAIN:
		fallthrough
	case ITEM_FOOD:
		write_to_output(d, libc.CString("Poisoned (0 = not poison) : "))
	case ITEM_VEHICLE:
		write_to_output(d, libc.CString("What is the vehicle's appearance? (-1 for transparent) : "))
	case ITEM_HATCH:
		write_to_output(d, libc.CString("Enter default vehicle load room : "))
	case ITEM_PORTAL:
		write_to_output(d, libc.CString("What is the portal's appearance? (-1 for transparent) : "))
	case ITEM_WINDOW:
		if (d.Olc.Obj.Value[0]) < 0 {
			write_to_output(d, libc.CString("What is the viewport room vnum (-1 for default location) : "))
		} else {
			oedit_disp_menu(d)
		}
	default:
		oedit_disp_val5_menu(d)
	}
}
func oedit_disp_val5_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_5
	write_to_output(d, libc.CString("Enter object default quality percentage (100%% MAX): "))
}
func oedit_disp_val7_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_7
	switch d.Olc.Obj.Type_flag {
	case ITEM_WEAPON:
		oedit_disp_crittype_menu(d)
	case ITEM_ARMOR:
		write_to_output(d, libc.CString("Arcane spell failure %% : "))
	default:
		oedit_disp_val9_menu(d)
	}
}
func oedit_disp_val9_menu(d *descriptor_data) {
	d.Olc.Mode = OEDIT_VALUE_9
	switch d.Olc.Obj.Type_flag {
	case ITEM_WEAPON:
		write_to_output(d, libc.CString("Default crit is only on natural 20. Extend this range by: "))
	default:
		oedit_disp_menu(d)
	}
}
func oedit_disp_type_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ITEM_TYPES; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter, item_types[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter object type : "))
}
func oedit_disp_extra_menu(d *descriptor_data) {
	var (
		bits    [64936]byte
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ITEM_FLAGS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter+1, extra_bits[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbitarray(d.Olc.Obj.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, &bits[0])
	write_to_output(d, libc.CString("\r\nObject flags: @c%s@n\r\nEnter object extra flag (0 to quit) : "), &bits[0])
}
func oedit_disp_perm_menu(d *descriptor_data) {
	var (
		bitbuf  [64936]byte
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 1; counter < NUM_AFF_FLAGS; counter++ {
		if counter == AFF_CHARM {
			continue
		}
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter, affected_bits[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbitarray(d.Olc.Obj.Bitvector[:], affected_bits[:], EF_ARRAY_MAX, &bitbuf[0])
	write_to_output(d, libc.CString("\r\nObject permanent flags: @c%s@n\r\nEnter object perm flag (0 to quit) : "), &bitbuf[0])
}
func oedit_disp_size_menu(d *descriptor_data) {
	var (
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_SIZES; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s%s"), counter+1, size_names[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\nEnter object size : "))
}
func oedit_disp_wear_menu(d *descriptor_data) {
	var (
		bits    [64936]byte
		counter int
		columns int = 0
	)
	clear_screen(d)
	for counter = 0; counter < NUM_ITEM_WEARS; counter++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s %s"), counter+1, wear_bits[counter], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbitarray(d.Olc.Obj.Wear_flags[:], wear_bits[:], TW_ARRAY_MAX, &bits[0])
	write_to_output(d, libc.CString("\r\nWear flags: @c%s@n\r\nEnter wear flag, 0 to quit : "), &bits[0])
}
func oedit_disp_menu(d *descriptor_data) {
	var (
		tbitbuf [2048]byte
		ebitbuf [2048]byte
		obj     *obj_data
	)
	obj = d.Olc.Obj
	clear_screen(d)
	sprinttype(int(obj.Type_flag), item_types[:], &tbitbuf[0], uint64(2048))
	sprintbitarray(obj.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, &ebitbuf[0])
	write_to_output(d, libc.CString("-- Item number : [@c%d@n]\r\n@g1@n) Namelist : @y%s@n\r\n@g2@n) S-Desc   : @y%s@n\r\n@g3@n) L-Desc   :-\r\n@y%s@n\r\n@g4@n) A-Desc   :-\r\n@y%s@n@g5@n) Type        : @c%s@n\r\n@g6@n) Extra flags : @c%s@n\r\n"), d.Olc.Number, func() *byte {
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
	}(), func() *byte {
		if obj.Action_description != nil && *obj.Action_description != 0 {
			return obj.Action_description
		}
		return libc.CString("Not Set.\r\n")
	}(), &tbitbuf[0], &ebitbuf[0])
	sprintbitarray(d.Olc.Obj.Wear_flags[:], wear_bits[:], EF_ARRAY_MAX, &tbitbuf[0])
	sprintbitarray(d.Olc.Obj.Bitvector[:], affected_bits[:], EF_ARRAY_MAX, &ebitbuf[0])
	write_to_output(d, libc.CString("@g7@n) Wear flags  : @c%s@n\r\n@g8@n) Weight      : @c%-4lld@n, \t@g9@n) Cost        : @c%-4d@n\r\n@gA@n) Cost/Day    : @c%-4d@n, \t@gB@n) Timer       : @c%-4d@n\r\n@gC@n) Values      : @c%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d@n\r\n@gD@n) Applies menu@n\r\n@gE@n) Extra descriptions menu %s\r\n@gM@n) Min Level   : @c%d@n\r\n@gN@n) Material    : @c%s@n\r\n@gP@n) Perm Affects: @c%s@n\r\n@gS@n) Script      : @c%s@n\r\n@gT@n) Spellbook menu\r\n@gW@n) Copy object        ,\t@gX@n) Delete object\r\n@gY@n) Size        : @c%s@n\r\n@gZ@n) Wiznet      :\r\n@gQ@n) Quit\r\nEnter choice : "), &tbitbuf[0], obj.Weight, obj.Cost, obj.Cost_per_day, obj.Timer, obj.Value[0], obj.Value[1], obj.Value[2], obj.Value[3], obj.Value[4], obj.Value[5], obj.Value[6], obj.Value[7], obj.Value[8], obj.Value[9], obj.Value[10], obj.Value[11], obj.Value[12], obj.Value[13], obj.Value[14], obj.Value[15], func() string {
		if obj.Extra_flags != nil {
			return "Set."
		}
		return "Not Set."
	}(), obj.Level, material_names[obj.Value[VAL_ALL_MATERIAL]], &ebitbuf[0], func() string {
		if d.Olc.Script != nil {
			return "Set."
		}
		return "Not Set."
	}(), size_names[obj.Size])
	d.Olc.Mode = OEDIT_MAIN_MENU
}
func oedit_parse(d *descriptor_data, arg *byte) {
	var (
		number  int
		max_val int
		min_val int
		oldtext *byte = nil
		tmp     *board_info
		obj     *obj_data
		robj    obj_rnum
	)
	switch d.Olc.Mode {
	case OEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			oedit_save_internally(d)
			mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s edits obj %d"), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				oedit_save_to_disk(int(real_zone_by_thing(d.Olc.Number)))
				write_to_output(d, libc.CString("Object saved to disk.\r\n"))
			} else {
				write_to_output(d, libc.CString("Object saved to memory.\r\n"))
			}
			if int(d.Olc.Obj.Type_flag) == ITEM_BOARD {
				if (func() *board_info {
					tmp = locate_board(GET_OBJ_VNUM(d.Olc.Obj))
					return tmp
				}()) != nil {
					save_board(tmp)
				} else {
					tmp = create_new_board(GET_OBJ_VNUM(d.Olc.Obj))
					tmp.Next = bboards
					bboards = tmp
				}
			}
			cleanup_olc(d, CLEANUP_ALL)
			return
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Obj.Proto_script = d.Olc.Script
			free_proto_script(unsafe.Pointer(d.Olc.Obj), OBJ_TRIGGER)
			cleanup_olc(d, CLEANUP_ALL)
			return
		case 'a':
			fallthrough
		case 'A':
			oedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\n"))
			write_to_output(d, libc.CString("Do you wish to save your changes? : \r\n"))
			return
		}
		fallthrough
	case OEDIT_COPY:
		if (func() int {
			number = int(real_object(obj_vnum(libc.Atoi(libc.GoString(arg)))))
			return number
		}()) != int(-1) {
			oedit_setup_existing(d, number)
		} else {
			write_to_output(d, libc.CString("That object does not exist.\r\n"))
		}
	case OEDIT_DELETE:
		if *arg == 'y' || *arg == 'Y' {
			if delete_object(obj_rnum(d.Olc.Obj.Item_number)) != int(-1) {
				write_to_output(d, libc.CString("Object deleted.\r\n"))
			} else {
				write_to_output(d, libc.CString("Couldn't delete the object!\r\n"))
			}
			cleanup_olc(d, CLEANUP_ALL)
		} else if *arg == 'n' || *arg == 'N' {
			oedit_disp_menu(d)
			d.Olc.Mode = OEDIT_MAIN_MENU
		} else {
			write_to_output(d, libc.CString("Please answer 'Y' or 'N': "))
		}
		return
	case OEDIT_MAIN_MENU:
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Connected != CON_IEDIT {
				if d.Olc.Value != 0 {
					write_to_output(d, libc.CString("Do you wish to save your changes? : "))
					d.Olc.Mode = OEDIT_CONFIRM_SAVESTRING
				} else {
					cleanup_olc(d, CLEANUP_ALL)
				}
			} else {
				send_to_char(d.Character, libc.CString("\r\nCommitting iedit changes.\r\n"))
				obj = d.Olc.Iobj
				*obj = *d.Olc.Obj
				obj.Id = int32(func() int {
					p := &max_obj_id
					x := *p
					*p++
					return x
				}())
				add_to_lookup_table(int(obj.Id), unsafe.Pointer(obj))
				if GET_OBJ_VNUM(obj) != obj_vnum(-1) {
					if obj.Script != nil {
						extract_script(unsafe.Pointer(obj), OBJ_TRIGGER)
						obj.Script = nil
					}
					free_proto_script(unsafe.Pointer(obj), OBJ_TRIGGER)
					robj = real_object(GET_OBJ_VNUM(obj))
					copy_proto_script(unsafe.Pointer((*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(robj)))), unsafe.Pointer(obj), OBJ_TRIGGER)
					assign_triggers(unsafe.Pointer(obj), OBJ_TRIGGER)
				}
				obj.Extra_flags[int(ITEM_UNIQUE_SAVE/32)] |= bitvector_t(int32(1 << (int(ITEM_UNIQUE_SAVE % 32))))
				mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s iedit a unique #%d"), GET_NAME(d.Character), GET_OBJ_VNUM(obj))
				if d.Character != nil {
					d.Character.Act[int(PLR_WRITING/32)] &= bitvector_t(int32(^(1 << (int(PLR_WRITING % 32)))))
					d.Connected = CON_PLAYING
					act(libc.CString("$n stops using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
				}
				libc.Free(unsafe.Pointer(d.Olc))
				d.Olc = nil
			}
			return
		case '1':
			write_to_output(d, libc.CString("Enter namelist : "))
			d.Olc.Mode = OEDIT_EDIT_NAMELIST
		case '2':
			write_to_output(d, libc.CString("Enter short desc : "))
			d.Olc.Mode = OEDIT_SHORTDESC
		case '3':
			write_to_output(d, libc.CString("Enter long desc :-\r\n| "))
			d.Olc.Mode = OEDIT_LONGDESC
		case '4':
			d.Olc.Mode = OEDIT_ACTDESC
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter action description:\r\n\r\n"))
			if d.Olc.Obj.Action_description != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Obj.Action_description)
				oldtext = libc.StrDup(d.Olc.Obj.Action_description)
			}
			string_write(d, &d.Olc.Obj.Action_description, MAX_MESSAGE_LENGTH, 0, unsafe.Pointer(oldtext))
			d.Olc.Value = 1
		case '5':
			oedit_disp_type_menu(d)
			d.Olc.Mode = OEDIT_TYPE
		case '6':
			oedit_disp_extra_menu(d)
			d.Olc.Mode = OEDIT_EXTRAS
		case '7':
			oedit_disp_wear_menu(d)
			d.Olc.Mode = OEDIT_WEAR
		case '8':
			write_to_output(d, libc.CString("Enter weight : "))
			d.Olc.Mode = OEDIT_WEIGHT
		case '9':
			write_to_output(d, libc.CString("Enter cost : "))
			d.Olc.Mode = OEDIT_COST
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter cost per day : "))
			d.Olc.Mode = OEDIT_COSTPERDAY
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter timer : "))
			d.Olc.Mode = OEDIT_TIMER
		case 'c':
			fallthrough
		case 'C':
			d.Olc.Obj.Value[0] = 0
			d.Olc.Obj.Value[1] = 0
			d.Olc.Obj.Value[2] = 0
			d.Olc.Obj.Value[3] = 0
			d.Olc.Obj.Value[4] = 0
			d.Olc.Obj.Value[5] = 0
			d.Olc.Obj.Value[6] = 0
			d.Olc.Value = 1
			oedit_disp_val1_menu(d)
		case 'd':
			fallthrough
		case 'D':
			oedit_disp_prompt_apply_menu(d)
		case 'e':
			fallthrough
		case 'E':
			if d.Olc.Obj.Ex_description == nil {
				d.Olc.Obj.Ex_description = new(extra_descr_data)
				d.Olc.Obj.Ex_description.Next = nil
			}
			d.Olc.Desc = d.Olc.Obj.Ex_description
			oedit_disp_extradesc_menu(d)
		case 'm':
			fallthrough
		case 'M':
			write_to_output(d, libc.CString("Enter new minimum level: "))
			d.Olc.Mode = OEDIT_LEVEL
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Mode = OEDIT_MATERIAL
			oedit_disp_material_menu(d)
		case 'p':
			fallthrough
		case 'P':
			oedit_disp_perm_menu(d)
			d.Olc.Mode = OEDIT_PERM
		case 's':
			fallthrough
		case 'S':
			if d.Connected != CON_IEDIT {
				d.Olc.Script_mode = SCRIPT_MAIN_MENU
				dg_script_menu(d)
			} else {
				write_to_output(d, libc.CString("\r\nScripts cannot be modified on individual objects.\r\nEnter choice : "))
			}
			return
		case 't':
			fallthrough
		case 'T':
			oedit_disp_prompt_spellbook_menu(d)
		case 'w':
			fallthrough
		case 'W':
			write_to_output(d, libc.CString("Copy what object? "))
			d.Olc.Mode = OEDIT_COPY
		case 'x':
			fallthrough
		case 'X':
			write_to_output(d, libc.CString("Are you sure you want to delete this object? "))
			d.Olc.Mode = OEDIT_DELETE
		case 'y':
			fallthrough
		case 'Y':
			oedit_disp_size_menu(d)
			d.Olc.Mode = OEDIT_SIZE
		case 'Z':
			fallthrough
		case 'z':
			search_replace(arg, libc.CString("z "), libc.CString(""))
			do_wiznet(d.Character, arg, 0, 0)
		default:
			oedit_disp_menu(d)
		}
		return
	case OLC_SCRIPT_EDIT:
		if dg_script_edit_parse(d, arg) != 0 {
			return
		}
	case OEDIT_EDIT_NAMELIST:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Obj.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Obj.Name))
		}
		d.Olc.Obj.Name = str_udup(arg)
	case OEDIT_SHORTDESC:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Obj.Short_description != nil {
			libc.Free(unsafe.Pointer(d.Olc.Obj.Short_description))
		}
		d.Olc.Obj.Short_description = str_udup(arg)
	case OEDIT_LONGDESC:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Obj.Description != nil {
			libc.Free(unsafe.Pointer(d.Olc.Obj.Description))
		}
		d.Olc.Obj.Description = str_udup(arg)
	case OEDIT_TYPE:
		number = libc.Atoi(libc.GoString(arg))
		if number < 1 || number >= NUM_ITEM_TYPES {
			write_to_output(d, libc.CString("Invalid choice, try again : "))
			return
		} else {
			d.Olc.Obj.Type_flag = int8(number)
		}
		d.Olc.Obj.Value[0] = func() int {
			p := &d.Olc.Obj.Value[1]
			d.Olc.Obj.Value[1] = func() int {
				p := &d.Olc.Obj.Value[2]
				d.Olc.Obj.Value[2] = func() int {
					p := &d.Olc.Obj.Value[3]
					d.Olc.Obj.Value[3] = 0
					return *p
				}()
				return *p
			}()
			return *p
		}()
	case OEDIT_EXTRAS:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > NUM_ITEM_FLAGS {
			oedit_disp_extra_menu(d)
			return
		} else if number == 0 {
			break
		} else {
			d.Olc.Obj.Extra_flags[(number-1)/32] = bitvector_t(int32(int(d.Olc.Obj.Extra_flags[(number-1)/32]) ^ 1<<((number-1)%32)))
			oedit_disp_extra_menu(d)
			return
		}
		fallthrough
	case OEDIT_WEAR:
		number = libc.Atoi(libc.GoString(arg))
		if number < 0 || number > NUM_ITEM_WEARS {
			write_to_output(d, libc.CString("That's not a valid choice!\r\n"))
			oedit_disp_wear_menu(d)
			return
		} else if number == 0 {
			break
		} else {
			d.Olc.Obj.Wear_flags[(number-1)/32] = d.Olc.Obj.Wear_flags[(number-1)/32] ^ 1<<((number-1)%32)
			oedit_disp_wear_menu(d)
			return
		}
		fallthrough
	case OEDIT_WEIGHT:
		d.Olc.Obj.Weight = int64(MIN(MAX_OBJ_WEIGHT, MAX(libc.Atoi(libc.GoString(arg)), 0)))
	case OEDIT_COST:
		d.Olc.Obj.Cost = MIN(MAX_OBJ_COST, MAX(libc.Atoi(libc.GoString(arg)), 0))
	case OEDIT_COSTPERDAY:
		d.Olc.Obj.Cost_per_day = MIN(MAX_OBJ_RENT, MAX(libc.Atoi(libc.GoString(arg)), 0))
	case OEDIT_TIMER:
		switch d.Olc.Obj.Type_flag {
		case ITEM_PORTAL:
			d.Olc.Obj.Timer = MIN(MAX_OBJ_TIMER, MAX(libc.Atoi(libc.GoString(arg)), -1))
		default:
			d.Olc.Obj.Timer = MIN(MAX_OBJ_TIMER, MAX(libc.Atoi(libc.GoString(arg)), 0))
		}
	case OEDIT_LEVEL:
		d.Olc.Obj.Level = MAX(libc.Atoi(libc.GoString(arg)), 0)
	case OEDIT_MATERIAL:
		d.Olc.Obj.Value[VAL_ALL_MATERIAL] = MIN(NUM_MATERIALS, MAX(libc.Atoi(libc.GoString(arg)), 0))
	case OEDIT_PERM:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) == 0 {
			break
		}
		if number > 0 && number <= NUM_AFF_FLAGS {
			if number != AFF_CHARM {
				d.Olc.Obj.Bitvector[number/32] = bitvector_t(int32(int(d.Olc.Obj.Bitvector[number/32]) ^ 1<<(number%32)))
			}
		}
		oedit_disp_perm_menu(d)
		return
	case OEDIT_SIZE:
		number = libc.Atoi(libc.GoString(arg)) - 1
		d.Olc.Obj.Size = MIN(int(NUM_SIZES-1), MAX(number, 0))
	case OEDIT_VALUE_1:
		switch d.Olc.Obj.Type_flag {
		case ITEM_WEAPON:
			d.Olc.Obj.Value[0] = MIN(MAX_WEAPON_TYPES, MAX(libc.Atoi(libc.GoString(arg)), WEAPON_TYPE_UNARMED))
		case ITEM_CONTAINER:
			d.Olc.Obj.Value[0] = MIN(MAX_CONTAINER_SIZE, MAX(libc.Atoi(libc.GoString(arg)), -1))
		default:
			d.Olc.Obj.Value[0] = libc.Atoi(libc.GoString(arg))
		}
		oedit_disp_val2_menu(d)
		return
	case OEDIT_VALUE_2:
		number = libc.Atoi(libc.GoString(arg))
		switch d.Olc.Obj.Type_flag {
		case ITEM_SCROLL:
			fallthrough
		case ITEM_POTION:
			if number == 0 || number == -1 {
				d.Olc.Obj.Value[1] = -1
			} else {
				d.Olc.Obj.Value[1] = MIN(SKILL_TABLE_SIZE, MAX(number, 1))
			}
			oedit_disp_val3_menu(d)
		case ITEM_CONTROL:
			if number <= 0 {
				d.Olc.Obj.Value[1] = 1
			} else if number > 3 {
				d.Olc.Obj.Value[1] = 3
			} else {
				d.Olc.Obj.Value[1] = number
			}
			oedit_disp_val5_menu(d)
		case ITEM_CONTAINER:
			fallthrough
		case ITEM_VEHICLE:
			fallthrough
		case ITEM_WINDOW:
			fallthrough
		case ITEM_HATCH:
			fallthrough
		case ITEM_PORTAL:
			if number < 0 || number > 4 {
				oedit_disp_container_flags_menu(d)
			} else if number != 0 {
				d.Olc.Obj.Value[1] ^= 1 << (number - 1)
				d.Olc.Value = 1
				oedit_disp_val2_menu(d)
			} else {
				oedit_disp_val3_menu(d)
			}
		case ITEM_WEAPON:
			d.Olc.Obj.Value[1] = MIN(MAX_WEAPON_NDICE, MAX(number, 1))
			oedit_disp_val3_menu(d)
		default:
			d.Olc.Obj.Value[1] = number
			oedit_disp_val3_menu(d)
		}
		return
	case OEDIT_VALUE_3:
		number = libc.Atoi(libc.GoString(arg))
		switch d.Olc.Obj.Type_flag {
		case ITEM_WEAPON:
			min_val = 1
			max_val = MAX_WEAPON_SDICE
		case ITEM_ARMOR:
			min_val = 0
			max_val = 100
			fallthrough
		case ITEM_WAND:
			fallthrough
		case ITEM_STAFF:
			min_val = 0
			max_val = 20
		case ITEM_DRINKCON:
			fallthrough
		case ITEM_FOUNTAIN:
			min_val = 0
			max_val = int(NUM_LIQ_TYPES - 1)
		case ITEM_KEY:
			min_val = 0
			max_val = 60000
		default:
			min_val = -32000
			max_val = 60000
		}
		d.Olc.Obj.Value[2] = MIN(max_val, MAX(number, min_val))
		oedit_disp_val4_menu(d)
		return
	case OEDIT_VALUE_4:
		number = libc.Atoi(libc.GoString(arg))
		switch d.Olc.Obj.Type_flag {
		case ITEM_HATCH:
			min_val = 1
			max_val = 60000
		case ITEM_WAND:
			fallthrough
		case ITEM_STAFF:
			min_val = 1
			max_val = int(SKILL_TABLE_SIZE - 1)
		case ITEM_WEAPON:
			min_val = 0
			max_val = int(NUM_ATTACK_TYPES - 1)
		case ITEM_ARMOR:
			if number < 0 {
				number = 0 - number
			}
			min_val = 0
			max_val = 20
		default:
			min_val = -32000
			max_val = 32000
		}
		d.Olc.Obj.Value[3] = MIN(max_val, MAX(number, min_val))
		oedit_disp_val5_menu(d)
		return
	case OEDIT_VALUE_5:
		min_val = 1
		max_val = 100
		d.Olc.Obj.Value[4] = MIN(max_val, MAX(libc.Atoi(libc.GoString(arg)), min_val))
		d.Olc.Obj.Value[5] = max_val
		oedit_disp_val7_menu(d)
		return
	case OEDIT_VALUE_7:
		number = libc.Atoi(libc.GoString(arg))
		switch d.Olc.Obj.Type_flag {
		case ITEM_WEAPON:
			min_val = 0
			max_val = CRIT_X4
		case ITEM_ARMOR:
			min_val = -100
			max_val = 100
		default:
			min_val = -32000
			max_val = 32000
		}
		d.Olc.Obj.Value[6] = MIN(max_val, MAX(libc.Atoi(libc.GoString(arg)), min_val))
		oedit_disp_val9_menu(d)
		return
	case OEDIT_VALUE_9:
		number = libc.Atoi(libc.GoString(arg))
		switch d.Olc.Obj.Type_flag {
		case ITEM_WEAPON:
			min_val = 0
			max_val = 19
		default:
			min_val = -32000
			max_val = 32000
		}
		d.Olc.Obj.Value[8] = MIN(max_val, MAX(libc.Atoi(libc.GoString(arg)), min_val))
	case OEDIT_PROMPT_APPLY:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) == 0 {
			break
		} else if number < 0 || number > MAX_OBJ_AFFECT {
			oedit_disp_prompt_apply_menu(d)
			return
		}
		d.Olc.Value = number - 1
		d.Olc.Mode = OEDIT_APPLY
		oedit_disp_apply_menu(d)
		return
	case OEDIT_APPLY:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) == 0 {
			d.Olc.Obj.Affected[d.Olc.Value].Location = 0
			d.Olc.Obj.Affected[d.Olc.Value].Modifier = 0
			oedit_disp_prompt_apply_menu(d)
		} else if number < 0 || number >= NUM_APPLIES {
			oedit_disp_apply_menu(d)
		} else {
			var counter int
			if d.Character.Admlevel < ADMLVL_GRGOD {
				for counter = 0; counter < MAX_OBJ_AFFECT; counter++ {
					if d.Olc.Obj.Affected[counter].Location == number {
						write_to_output(d, libc.CString("Object already has that apply."))
						return
					}
				}
			}
			d.Olc.Obj.Affected[d.Olc.Value].Location = number
			write_to_output(d, libc.CString("Modifier : "))
			d.Olc.Mode = OEDIT_APPLYMOD
		}
		return
	case OEDIT_APPLYMOD:
		d.Olc.Obj.Affected[d.Olc.Value].Modifier = libc.Atoi(libc.GoString(arg))
		oedit_disp_apply_spec_menu(d)
		return
	case OEDIT_APPLYSPEC:
		if unicode.IsDigit(rune(*arg)) {
			d.Olc.Obj.Affected[d.Olc.Value].Specific = libc.Atoi(libc.GoString(arg))
		} else {
			switch d.Olc.Obj.Affected[d.Olc.Value].Location {
			case APPLY_SKILL:
				number = find_skill_num(arg, 1<<1)
				if number > -1 {
					d.Olc.Obj.Affected[d.Olc.Value].Specific = number
				}
			case APPLY_FEAT:
				number = find_feat_num(arg)
				if number > -1 {
					d.Olc.Obj.Affected[d.Olc.Value].Specific = number
				}
			default:
				d.Olc.Obj.Affected[d.Olc.Value].Specific = 0
			}
		}
		oedit_disp_prompt_apply_menu(d)
		return
	case OEDIT_EXTRADESC_KEY:
		if genolc_checkstring(d, arg) != 0 {
			if d.Olc.Desc.Keyword != nil {
				libc.Free(unsafe.Pointer(d.Olc.Desc.Keyword))
			}
			d.Olc.Desc.Keyword = str_udup(arg)
		}
		oedit_disp_extradesc_menu(d)
		return
	case OEDIT_EXTRADESC_MENU:
		switch func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}() {
		case 0:
			if d.Olc.Desc.Keyword == nil || d.Olc.Desc.Description == nil {
				var temp *extra_descr_data
				if d.Olc.Desc.Keyword != nil {
					libc.Free(unsafe.Pointer(d.Olc.Desc.Keyword))
				}
				if d.Olc.Desc.Description != nil {
					libc.Free(unsafe.Pointer(d.Olc.Desc.Description))
				}
				if d.Olc.Desc == d.Olc.Obj.Ex_description {
					d.Olc.Obj.Ex_description = d.Olc.Desc.Next
				} else {
					temp = d.Olc.Obj.Ex_description
					for temp != nil && temp.Next != d.Olc.Desc {
						temp = temp.Next
					}
					if temp != nil {
						temp.Next = d.Olc.Desc.Next
					}
				}
				libc.Free(unsafe.Pointer(d.Olc.Desc))
				d.Olc.Desc = nil
			}
		case 1:
			d.Olc.Mode = OEDIT_EXTRADESC_KEY
			write_to_output(d, libc.CString("Enter keywords, separated by spaces :-\r\n| "))
			return
		case 2:
			d.Olc.Mode = OEDIT_EXTRADESC_DESCRIPTION
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter the extra description:\r\n\r\n"))
			if d.Olc.Desc.Description != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Desc.Description)
				oldtext = libc.StrDup(d.Olc.Desc.Description)
			}
			string_write(d, &d.Olc.Desc.Description, MAX_MESSAGE_LENGTH, 0, unsafe.Pointer(oldtext))
			d.Olc.Value = 1
			return
		case 3:
			if d.Olc.Desc.Keyword != nil && d.Olc.Desc.Description != nil {
				var new_extra *extra_descr_data
				if d.Olc.Desc.Next != nil {
					d.Olc.Desc = d.Olc.Desc.Next
				} else {
					new_extra = new(extra_descr_data)
					d.Olc.Desc.Next = new_extra
					d.Olc.Desc = d.Olc.Desc.Next
				}
			}
			fallthrough
		default:
			oedit_disp_extradesc_menu(d)
			return
		}
	case OEDIT_PROMPT_SPELLBOOK:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) == 0 {
			break
		} else if number < 0 || number > SKILL_TABLE_SIZE {
			oedit_disp_prompt_spellbook_menu(d)
			return
		}
		d.Olc.Value = number - 1
		d.Olc.Mode = OEDIT_SPELLBOOK
		oedit_disp_spellbook_menu(d)
		return
	case OEDIT_SPELLBOOK:
		if (func() int {
			number = libc.Atoi(libc.GoString(arg))
			return number
		}()) == 0 {
			if d.Olc.Obj.Sbinfo != nil {
				(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Spellname = 0
				(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Pages = 0
			} else {
				d.Olc.Obj.Sbinfo = &make([]obj_spellbook_spell, SPELLBOOK_SIZE)[0]
				(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Spellname = 0
				(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Pages = 0
			}
			oedit_disp_prompt_spellbook_menu(d)
		} else if number < 0 || number >= SKILL_TABLE_SIZE {
			oedit_disp_spellbook_menu(d)
		} else {
			var counter int
			if GET_LEVEL(d.Character) < ADMLVL_IMPL {
				for counter = 0; counter < SKILL_TABLE_SIZE; counter++ {
					if d.Olc.Obj.Sbinfo != nil && (*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(counter)))).Spellname == number {
						write_to_output(d, libc.CString("Object already has that spell."))
						return
					}
				}
			}
			if d.Olc.Obj.Sbinfo == nil {
				d.Olc.Obj.Sbinfo = &make([]obj_spellbook_spell, SPELLBOOK_SIZE)[0]
			}
			(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Spellname = number
			(*(*obj_spellbook_spell)(unsafe.Add(unsafe.Pointer(d.Olc.Obj.Sbinfo), unsafe.Sizeof(obj_spellbook_spell{})*uintptr(d.Olc.Value)))).Pages = MAX(1, spell_info[number].Spell_level*2)
			oedit_disp_prompt_spellbook_menu(d)
		}
		return
	default:
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: Reached default case in oedit_parse()!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
	d.Olc.Value = 1
	oedit_disp_menu(d)
}
func oedit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case OEDIT_ACTDESC:
		oedit_disp_menu(d)
	case OEDIT_EXTRADESC_DESCRIPTION:
		oedit_disp_extradesc_menu(d)
	}
}
func iedit_setup_existing(d *descriptor_data, real_num *obj_data) {
	var obj *obj_data
	d.Olc.Iobj = real_num
	obj = create_obj()
	copy_object(obj, real_num)
	if obj.Script != nil {
		extract_script(unsafe.Pointer(obj), OBJ_TRIGGER)
	}
	obj.Script = nil
	remove_from_lookup_table(int(obj.Id))
	d.Olc.Obj = obj
	d.Olc.Iobj = real_num
	d.Olc.Value = 0
	oedit_disp_menu(d)
}
func do_iedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		k     *obj_data
		found int = 0
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 || *argument == 0 {
		send_to_char(ch, libc.CString("You must supply an object name.\r\n"))
	}
	if (func() *obj_data {
		k = get_obj_in_equip_vis(ch, &arg[0], nil, ch.Equipment[:])
		return k
	}()) != nil {
		found = 1
	} else if (func() *obj_data {
		k = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return k
	}()) != nil {
		found = 1
	} else if (func() *obj_data {
		k = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
		return k
	}()) != nil {
		found = 1
	} else if (func() *obj_data {
		k = get_obj_vis(ch, &arg[0], nil)
		return k
	}()) != nil {
		found = 1
	}
	if found == 0 {
		send_to_char(ch, libc.CString("Couldn't find that object. Sorry.\r\n"))
		return
	}
	ch.Desc.Olc = new(oasis_olc_data)
	k.Extra_flags[int(ITEM_UNIQUE_SAVE/32)] |= bitvector_t(int32(1 << (int(ITEM_UNIQUE_SAVE % 32))))
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
	iedit_setup_existing(ch.Desc, k)
	ch.Desc.Olc.Value = 0
	act(libc.CString("$n starts using OLC."), TRUE, ch, nil, nil, TO_ROOM)
	ch.Desc.Connected = CON_IEDIT
	return
}
