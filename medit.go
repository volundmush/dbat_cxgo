package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unicode"
	"unsafe"
)

func do_oasis_medit(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		send_to_char(ch, libc.CString("Specify a mobile VNUM to edit.\r\n"))
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
		if d.Connected == CON_MEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That mobile is currently being edited by %s.\r\n"), GET_NAME(d.Character))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_medit: Player already had olc structure."))
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
		send_cannot_edit(ch, zone_table[d.Olc.Zone_num].Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all mobiles in zone %d.\r\n"), zone_table[d.Olc.Zone_num].Number)
		mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(ch.Player_specials.Invis_level))), TRUE, libc.CString("OLC: %s saves mobile info for zone %d."), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number)
		save_mobiles(d.Olc.Zone_num)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() int {
		real_num = int(real_mobile(mob_vnum(number)))
		return real_num
	}()) == int(-1) {
		medit_setup_new(d)
	} else {
		medit_setup_existing(d, real_num)
	}
	medit_disp_menu(d)
	d.Connected = CON_MEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	SET_BIT_AR(ch.Act[:], PLR_WRITING)
	mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number, ch.Player_specials.Olc_zone)
}
func medit_save_to_disk(foo zone_vnum) {
	save_mobiles(real_zone(foo))
}
func medit_setup_new(d *descriptor_data) {
	var mob *char_data
	mob = new(char_data)
	init_mobile(mob)
	mob.Nr = -1
	mob.Name = libc.CString("mob unfinished")
	mob.Short_descr = libc.CString("the unfinished mob")
	mob.Long_descr = libc.CString("An unfinished mob stands here.\r\n")
	mob.Description = libc.CString("It looks unfinished.\r\n")
	mob.Script = nil
	mob.Proto_script = func() *trig_proto_list {
		p := &d.Olc.Script
		d.Olc.Script = nil
		return *p
	}()
	d.Olc.Mob = mob
	d.Olc.Value = FALSE
	d.Olc.Item_type = MOB_TRIGGER
}
func medit_setup_existing(d *descriptor_data, rmob_num int) {
	var mob *char_data
	mob = new(char_data)
	copy_mobile(mob, &mob_proto[rmob_num])
	d.Olc.Mob = mob
	d.Olc.Item_type = MOB_TRIGGER
	dg_olc_script_copy(d)
	mob.Script = nil
	d.Olc.Mob.Proto_script = nil
}
func init_mobile(mob *char_data) {
	clear_char(mob)
	mob.Hit = 0
	mob.Max_mana = 0
	mob.Mob_specials.Damnodice = 0
	mob.Sex = SEX_MALE
	mob.Race_level = 0
	mob.Chclass = CLASS_NPC_COMMONER
	mob.Weight = uint8(int8(rand_number(100, 200)))
	mob.Height = uint8(int8(rand_number(100, 200)))
	mob.Real_abils.Str = func() int8 {
		p := &mob.Real_abils.Intel
		mob.Real_abils.Intel = func() int8 {
			p := &mob.Real_abils.Wis
			mob.Real_abils.Wis = int8(rand_number(8, 16))
			return *p
		}()
		return *p
	}()
	mob.Real_abils.Dex = func() int8 {
		p := &mob.Real_abils.Con
		mob.Real_abils.Con = func() int8 {
			p := &mob.Real_abils.Cha
			mob.Real_abils.Cha = int8(rand_number(8, 16))
			return *p
		}()
		return *p
	}()
	mob.Aff_abils = mob.Real_abils
	SET_BIT_AR(mob.Act[:], MOB_ISNPC)
	mob.Player_specials = &dummy_mob
}
func medit_save_internally(d *descriptor_data) {
	var (
		i        int
		new_rnum mob_rnum
		dsc      *descriptor_data
		mob      *char_data
	)
	i = int(libc.BoolToInt(real_mobile(mob_vnum(d.Olc.Number)) == mob_rnum(-1)))
	if (func() mob_rnum {
		new_rnum = mob_rnum(add_mobile(d.Olc.Mob, mob_vnum(d.Olc.Number)))
		return new_rnum
	}()) == mob_rnum(-1) {
		basic_mud_log(libc.CString("medit_save_internally: add_mobile failed."))
		return
	}
	if mob_proto[new_rnum].Proto_script != nil && mob_proto[new_rnum].Proto_script != d.Olc.Script {
		free_proto_script(unsafe.Pointer(&mob_proto[new_rnum]), MOB_TRIGGER)
	}
	mob_proto[new_rnum].Proto_script = d.Olc.Script
	for mob = character_list; mob != nil; mob = mob.Next {
		if mob.Nr != new_rnum {
			continue
		}
		if mob.Script != nil {
			extract_script(unsafe.Pointer(mob), MOB_TRIGGER)
		}
		free_proto_script(unsafe.Pointer(mob), MOB_TRIGGER)
		copy_proto_script(unsafe.Pointer(&mob_proto[new_rnum]), unsafe.Pointer(mob), MOB_TRIGGER)
		assign_triggers(unsafe.Pointer(mob), MOB_TRIGGER)
	}
	if i == 0 {
		return
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		if dsc.Connected == CON_SEDIT {
			dsc.Olc.Shop.Keeper += mob_rnum(libc.BoolToInt(dsc.Olc.Shop.Keeper != mob_rnum(-1) && dsc.Olc.Shop.Keeper >= new_rnum))
		} else if dsc.Connected == CON_MEDIT {
			dsc.Olc.Mob.Nr += mob_rnum(libc.BoolToInt(dsc.Olc.Mob.Nr != mob_rnum(-1) && dsc.Olc.Mob.Nr >= new_rnum))
		}
	}
	for dsc = descriptor_list; dsc != nil; dsc = dsc.Next {
		if dsc.Connected == CON_ZEDIT {
			for i = 0; int(dsc.Olc.Zone.Cmd[i].Command) != 'S'; i++ {
				if int(dsc.Olc.Zone.Cmd[i].Command) == 'M' {
					if dsc.Olc.Zone.Cmd[i].Arg1 >= vnum(new_rnum) {
						dsc.Olc.Zone.Cmd[i].Arg1++
					}
				}
			}
		}
	}
}
func medit_disp_positions(d *descriptor_data) {
	var i int
	clear_screen(d)
	for i = 0; *position_types[i] != '\n'; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %s\r\n"), i, position_types[i])
	}
	write_to_output(d, libc.CString("Enter position number : "))
}
func medit_disp_sex(d *descriptor_data) {
	var i int
	clear_screen(d)
	for i = 0; i < NUM_GENDERS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %s\r\n"), i, genders[i])
	}
	write_to_output(d, libc.CString("Enter gender number : "))
}
func medit_disp_mob_flags(d *descriptor_data) {
	var (
		i       int
		columns int = 0
		flags   [64936]byte
	)
	clear_screen(d)
	for i = 0; i < NUM_MOB_FLAGS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s  %s"), i+1, action_bits[i], func() string {
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
	sprintbitarray(d.Olc.Mob.Act[:], action_bits[:], AF_ARRAY_MAX, &flags[0])
	write_to_output(d, libc.CString("\r\nCurrent flags : @c%s@n\r\nEnter mob flags (0 to quit) : "), &flags[0])
}
func medit_disp_personality(d *descriptor_data) {
	write_to_output(d, libc.CString("@GPersonalities\n"))
	write_to_output(d, libc.CString("@D--------------@n\n"))
	write_to_output(d, libc.CString("@w1@D) @WBasic@n\n"))
	write_to_output(d, libc.CString("@w1@D) @WCareful@n\n"))
	write_to_output(d, libc.CString("@w1@D) @WAggressive@n\n"))
	write_to_output(d, libc.CString("@w1@D) @WArrogant\n"))
	write_to_output(d, libc.CString("@w1@D) @WIntelligent@n\n"))
}
func medit_disp_aff_flags(d *descriptor_data) {
	var (
		i       int
		columns int = 0
		flags   [64936]byte
	)
	clear_screen(d)
	for i = 0; i < NUM_AFF_FLAGS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s  %s"), i+1, affected_bits[i+1], func() string {
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
	sprintbitarray(d.Olc.Mob.Affected_by[:], affected_bits[:], AF_ARRAY_MAX, &flags[0])
	write_to_output(d, libc.CString("\r\nCurrent flags   : @c%s@n\r\nEnter aff flags (0 to quit) : "), &flags[0])
}
func medit_disp_class(d *descriptor_data) {
	var (
		i   int
		buf [2048]byte
	)
	clear_screen(d)
	for i = 0; i < NUM_CLASSES; i++ {
		stdio.Sprintf(&buf[0], "@g%2d@n) %s\r\n", i, pc_class_types[i])
		write_to_output(d, &buf[0])
	}
	write_to_output(d, libc.CString("Enter class number : "))
}
func medit_disp_race(d *descriptor_data) {
	var (
		i       int
		columns int = 0
		buf     [2048]byte
	)
	clear_screen(d)
	for i = 0; i < NUM_RACES; i++ {
		stdio.Sprintf(&buf[0], "@g%2d@n) %-20.20s  %s", i, pc_race_types[i], func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 2) == 0 {
				return "\r\n"
			}
			return ""
		}())
		write_to_output(d, &buf[0])
	}
	write_to_output(d, libc.CString("Enter race number : "))
}
func medit_disp_size(d *descriptor_data) {
	var (
		i       int
		columns int = 0
		buf     [2048]byte
	)
	clear_screen(d)
	for i = -1; i < NUM_SIZES; i++ {
		stdio.Sprintf(&buf[0], "@g%2d@n) %-20.20s  %s", i, func() string {
			if i == int(-1) {
				return "DEFAULT"
			}
			return libc.GoString(size_names[i])
		}(), func() string {
			if (func() int {
				p := &columns
				*p++
				return *p
			}() % 2) == 0 {
				return "\r\n"
			}
			return ""
		}())
		write_to_output(d, &buf[0])
	}
	write_to_output(d, libc.CString("Enter size number (-1 for default): "))
}
func medit_disp_menu(d *descriptor_data) {
	var (
		mob   *char_data
		flags [64936]byte
		flag2 [64936]byte
	)
	mob = d.Olc.Mob
	clear_screen(d)
	write_to_output(d, libc.CString("-- Mob Number:  [@c%d@n]\r\n@g1@n) Sex: @y%-7.7s@n\t         @g2@n) Alias: @y%s\r\n@g3@n) S-Desc: @y%s\r\n@g4@n) L-Desc:-\r\n@y%s@g5@n) D-Desc:-\r\n@y%s@g6@n) Level:       [@c%4d@n],  @g7@n) Alignment:    [@c%5d@n]\r\n@g8@n) Accuracy Mod:[@c%4d@n],  @g9@n) Damage Mod:   [@c%5d@n]\r\n@gA@n) NumDamDice:  [@c%4d@n],  @gB@n) SizeDamDice:  [@c%5d@n]\r\n@gC@n) Num HP Dice: [@c%4lld@n],  @gD@n) Size HP Dice: [@c%5lld@n],  @gE@n) HP Bonus: [@c%5lld@n]\r\n@gF@n) Armor Class: [@c%4d@n],  @gG@n) Exp:      [@c%lld@n],  @gH@n) Gold:  [@c%8d@n]\r\n"), d.Olc.Number, genders[int(mob.Sex)], mob.Name, mob.Short_descr, mob.Long_descr, mob.Description, mob.Race_level, mob.Alignment, mob.Accuracy_mod, mob.Damage_mod, mob.Mob_specials.Damnodice, mob.Mob_specials.Damsizedice, mob.Hit, mob.Mana, mob.Move, mob.Armor, mob.Exp, mob.Gold)
	sprintbitarray(mob.Act[:], action_bits[:], AF_ARRAY_MAX, &flags[0])
	sprintbitarray(mob.Affected_by[:], affected_bits[:], AF_ARRAY_MAX, &flag2[0])
	write_to_output(d, libc.CString("@gI@n) Position   : @y%-10s@n,\t @gJ@n) Default   : @y%-10s\r\n@gK@n) Personality: @Y%s@n\r\n@gL@n) NPC Flags  : @c%s\r\n@gM@n) AFF Flags  : @c%s\r\n@gN@n) Class      : @y%-10s@n,\t @gO@n) Race      : @y%-10s\r\n@gS@n) Script     : @c%s\r\n@gW@n) Copy mob              ,\t @gX@n) Delete mob\r\n@gY@n) Size       : @y%s\r\n@gZ@n) Wiznet     :\r\n@gQ@n) Quit\r\nEnter choice : "), position_types[int(mob.Position)], position_types[int(mob.Mob_specials.Default_pos)], npc_personality[mob.Personality], &flags[0], &flag2[0], pc_class_types[int(mob.Chclass)], pc_race_types[int(mob.Race)], func() string {
		if d.Olc.Script != nil {
			return "Set."
		}
		return "Not Set."
	}(), size_names[get_size(mob)])
	d.Olc.Mode = MEDIT_MAIN_MENU
}
func medit_parse(d *descriptor_data, arg *byte) {
	var (
		i       int   = -1
		oldtext *byte = nil
	)
	if d.Olc.Mode > MEDIT_NUMERICAL_RESPONSE {
		i = libc.Atoi(libc.GoString(arg))
		if *arg == 0 || !unicode.IsDigit(rune(*arg)) && (*arg == '-' && !unicode.IsDigit(rune(*(*byte)(unsafe.Add(unsafe.Pointer(arg), 1))))) {
			write_to_output(d, libc.CString("Field must be numerical, try again : "))
			return
		}
	} else {
		if genolc_checkstring(d, arg) == 0 {
			return
		}
	}
	switch d.Olc.Mode {
	case MEDIT_CONFIRM_SAVESTRING:
		SET_BIT_AR(d.Olc.Mob.Act[:], MOB_ISNPC)
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			medit_save_internally(d)
			mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(d.Character.Player_specials.Invis_level))), TRUE, libc.CString("OLC: %s edits mob %d"), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				medit_save_to_disk(zone_table[real_zone_by_thing(d.Olc.Number)].Number)
				write_to_output(d, libc.CString("Mobile saved to disk.\r\n"))
			} else {
				write_to_output(d, libc.CString("Mobile saved to memory.\r\n"))
			}
			cleanup_olc(d, CLEANUP_ALL)
			return
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Mob.Proto_script = d.Olc.Script
			cleanup_olc(d, CLEANUP_ALL)
			return
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\n"))
			write_to_output(d, libc.CString("Do you wish to save the mobile? : "))
			return
		}
	case MEDIT_MAIN_MENU:
		i = 0
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = MEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '1':
			d.Olc.Mode = MEDIT_SEX
			medit_disp_sex(d)
			return
		case '2':
			d.Olc.Mode = MEDIT_ALIAS
			i--
		case '3':
			d.Olc.Mode = MEDIT_S_DESC
			i--
		case '4':
			d.Olc.Mode = MEDIT_L_DESC
			i--
		case '5':
			d.Olc.Mode = MEDIT_D_DESC
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter mob description:\r\n\r\n"))
			if d.Olc.Mob.Description != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Mob.Description)
				oldtext = libc.StrDup(d.Olc.Mob.Description)
			}
			string_write(d, &d.Olc.Mob.Description, MAX_MOB_DESC, 0, unsafe.Pointer(oldtext))
			d.Olc.Value = 1
			return
		case '6':
			d.Olc.Mode = MEDIT_LEVEL
			i++
		case '7':
			d.Olc.Mode = MEDIT_ALIGNMENT
			i++
		case '8':
			d.Olc.Mode = MEDIT_ACCURACY
			i++
		case '9':
			d.Olc.Mode = MEDIT_DAMAGE
			i++
		case 'a':
			fallthrough
		case 'A':
			d.Olc.Mode = MEDIT_NDD
			i++
		case 'b':
			fallthrough
		case 'B':
			d.Olc.Mode = MEDIT_SDD
			i++
		case 'c':
			fallthrough
		case 'C':
			d.Olc.Mode = MEDIT_NUM_HP_DICE
			i++
		case 'd':
			fallthrough
		case 'D':
			d.Olc.Mode = MEDIT_SIZE_HP_DICE
			i++
		case 'e':
			fallthrough
		case 'E':
			d.Olc.Mode = MEDIT_ADD_HP
			i++
		case 'f':
			fallthrough
		case 'F':
			d.Olc.Mode = MEDIT_AC
			i++
		case 'g':
			fallthrough
		case 'G':
			d.Olc.Mode = MEDIT_EXP
			i++
		case 'h':
			fallthrough
		case 'H':
			d.Olc.Mode = MEDIT_GOLD
			i++
		case 'i':
			fallthrough
		case 'I':
			d.Olc.Mode = MEDIT_POS
			medit_disp_positions(d)
			return
		case 'j':
			fallthrough
		case 'J':
			d.Olc.Mode = MEDIT_DEFAULT_POS
			medit_disp_positions(d)
			return
		case 'k':
			fallthrough
		case 'K':
			d.Olc.Mode = MEDIT_PERSONALITY
			medit_disp_personality(d)
			return
		case 'l':
			fallthrough
		case 'L':
			d.Olc.Mode = MEDIT_NPC_FLAGS
			medit_disp_mob_flags(d)
			return
		case 'm':
			fallthrough
		case 'M':
			d.Olc.Mode = MEDIT_AFF_FLAGS
			medit_disp_aff_flags(d)
			return
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Mode = MEDIT_CLASS
			medit_disp_class(d)
			return
		case 'o':
			fallthrough
		case 'O':
			d.Olc.Mode = MEDIT_RACE
			medit_disp_race(d)
			return
		case 's':
			fallthrough
		case 'S':
			d.Olc.Script_mode = SCRIPT_MAIN_MENU
			dg_script_menu(d)
			return
		case 'w':
			fallthrough
		case 'W':
			write_to_output(d, libc.CString("Copy what mob? "))
			d.Olc.Mode = MEDIT_COPY
			return
		case 'x':
			fallthrough
		case 'X':
			write_to_output(d, libc.CString("Are you sure you want to delete this mobile? "))
			d.Olc.Mode = MEDIT_DELETE
			return
		case 'y':
			fallthrough
		case 'Y':
			d.Olc.Mode = MEDIT_SIZE
			medit_disp_size(d)
			return
		case 'Z':
			fallthrough
		case 'z':
			search_replace(arg, libc.CString("z "), libc.CString(""))
			do_wiznet(d.Character, arg, 0, 0)
		default:
			medit_disp_menu(d)
			return
		}
		if i == 0 {
			break
		} else if i == 1 {
			write_to_output(d, libc.CString("\r\nEnter new value : "))
		} else if i == -1 {
			write_to_output(d, libc.CString("\r\nEnter new text :\r\n] "))
		} else {
			write_to_output(d, libc.CString("Oops...\r\n"))
		}
		return
	case OLC_SCRIPT_EDIT:
		if dg_script_edit_parse(d, arg) != 0 {
			return
		}
	case MEDIT_ALIAS:
		smash_tilde(arg)
		if d.Olc.Mob.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Mob.Name))
		}
		d.Olc.Mob.Name = str_udup(arg)
	case MEDIT_S_DESC:
		smash_tilde(arg)
		if d.Olc.Mob.Short_descr != nil {
			libc.Free(unsafe.Pointer(d.Olc.Mob.Short_descr))
		}
		d.Olc.Mob.Short_descr = str_udup(arg)
	case MEDIT_L_DESC:
		smash_tilde(arg)
		if d.Olc.Mob.Long_descr != nil {
			libc.Free(unsafe.Pointer(d.Olc.Mob.Long_descr))
		}
		if arg != nil && *arg != 0 {
			var buf [2048]byte
			stdio.Snprintf(&buf[0], int(2048), "%s\r\n", arg)
			d.Olc.Mob.Long_descr = libc.StrDup(&buf[0])
		} else {
			d.Olc.Mob.Long_descr = libc.CString("undefined")
		}
	case MEDIT_D_DESC:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: medit_parse(): Reached D_DESC case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	case MEDIT_NPC_FLAGS:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) <= 0 {
			break
		} else if i <= NUM_MOB_FLAGS {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], bitvector_t(int32(i-1)))
		}
		medit_disp_mob_flags(d)
		return
	case MEDIT_PERSONALITY:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) <= 0 {
			break
		} else if i <= MAX_PERSONALITIES {
			d.Olc.Mob.Personality = i
		}
		medit_disp_personality(d)
		return
	case MEDIT_AFF_FLAGS:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) <= 0 {
			break
		} else if i <= NUM_AFF_FLAGS {
			TOGGLE_BIT_AR(d.Olc.Mob.Affected_by[:], bitvector_t(int32(i)))
		}
		REMOVE_BIT_AR(d.Olc.Mob.Affected_by[:], bitvector_t(int32(int(AFF_CHARM|AFF_POISON)|AFF_GROUP|AFF_SLEEP)))
		medit_disp_aff_flags(d)
		return
	case MEDIT_SEX:
		d.Olc.Mob.Sex = int8(MIN(int64(int(NUM_GENDERS-1)), MAX(int64(i), 0)))
	case MEDIT_ACCURACY:
		d.Olc.Mob.Accuracy_mod = int(MIN(50, MAX(int64(i), 0)))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_DAMAGE:
		d.Olc.Mob.Damage_mod = int(MIN(50, MAX(int64(i), 0)))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_NDD:
		d.Olc.Mob.Mob_specials.Damnodice = int8(MIN(30, MAX(int64(i), 0)))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_SDD:
		d.Olc.Mob.Mob_specials.Damsizedice = int8(MIN(math.MaxInt8, MAX(int64(i), 0)))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_NUM_HP_DICE:
		d.Olc.Mob.Hit = MIN(int64(config_info.Play.Level_cap), MAX(int64(i), 0))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_SIZE_HP_DICE:
		d.Olc.Mob.Mana = MIN(1000, MAX(int64(i), 0))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_ADD_HP:
		d.Olc.Mob.Move = MIN(30000, MAX(int64(i), 0))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_AC:
		d.Olc.Mob.Armor = int(MIN(200000, MAX(int64(i), 10)))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_EXP:
		d.Olc.Mob.Exp = MIN(MAX_MOB_EXP, MAX(int64(i), 0))
		if MOB_FLAGGED(d.Olc.Mob, MOB_AUTOBALANCE) {
			TOGGLE_BIT_AR(d.Olc.Mob.Act[:], MOB_AUTOBALANCE)
		}
	case MEDIT_GOLD:
		d.Olc.Mob.Gold = int(MIN(MAX_MOB_GOLD, MAX(int64(i), 0)))
	case MEDIT_POS:
		d.Olc.Mob.Position = int8(MIN(int64(int(NUM_POSITIONS-1)), MAX(int64(i), 0)))
	case MEDIT_DEFAULT_POS:
		d.Olc.Mob.Mob_specials.Default_pos = int8(MIN(int64(int(NUM_POSITIONS-1)), MAX(int64(i), 0)))
	case MEDIT_ATTACK:
		d.Olc.Mob.Mob_specials.Attack_type = int8(MIN(int64(int(NUM_ATTACK_TYPES-1)), MAX(int64(i), 0)))
	case MEDIT_LEVEL:
		d.Olc.Mob.Race_level = int(MIN(150, MAX(int64(i), 1)))
	case MEDIT_ALIGNMENT:
		d.Olc.Mob.Alignment = int(MIN(1000, MAX(int64(i), -1000)))
	case MEDIT_CLASS:
		d.Olc.Mob.Chclass = int8(MIN(NUM_CLASSES, MAX(int64(i), 0)))
		d.Olc.Mob.Mana = int64(class_hit_die_size[d.Olc.Mob.Chclass])
	case MEDIT_COPY:
		if (func() int {
			i = int(real_mobile(mob_vnum(libc.Atoi(libc.GoString(arg)))))
			return i
		}()) != int(-1) {
			medit_setup_existing(d, i)
		} else {
			write_to_output(d, libc.CString("That mob does not exist.\r\n"))
		}
	case MEDIT_DELETE:
		if *arg == 'y' || *arg == 'Y' {
			if delete_mobile(d.Olc.Mob.Nr) != int(-1) {
				write_to_output(d, libc.CString("Mobile deleted.\r\n"))
			} else {
				write_to_output(d, libc.CString("Couldn't delete the mobile!\r\n"))
			}
			cleanup_olc(d, CLEANUP_ALL)
			return
		} else if *arg == 'n' || *arg == 'N' {
			medit_disp_menu(d)
			d.Olc.Mode = MEDIT_MAIN_MENU
			return
		} else {
			write_to_output(d, libc.CString("Please answer 'Y' or 'N': "))
		}
	case MEDIT_RACE:
		d.Olc.Mob.Race = int8(MIN(NUM_RACES, MAX(int64(i), 0)))
		d.Olc.Mob.Size = race_def_sizetable[d.Olc.Mob.Race]
	case MEDIT_SIZE:
		d.Olc.Mob.Size = int(MIN(int64(int(NUM_SIZES-1)), MAX(int64(i), -1)))
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: medit_parse(): Reached default case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
	d.Olc.Value = TRUE
	medit_disp_menu(d)
}
func medit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case MEDIT_D_DESC:
		fallthrough
	default:
		medit_disp_menu(d)
	}
}
