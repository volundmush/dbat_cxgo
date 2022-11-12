package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func gedit_save_internally(d *descriptor_data) {
	d.Olc.Guild.Vnum = d.Olc.Number
	add_guild(d.Olc.Guild)
}
func gedit_save_to_disk(num int) {
	save_guilds(zone_rnum(num))
}
func do_oasis_gedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int = int(-1)
		save     int = 0
		real_num guild_rnum
		d        *descriptor_data
		buf3     *byte
	)
	_ = buf3
	var buf1 [2048]byte
	var buf2 [2048]byte
	buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		send_to_char(ch, libc.CString("Specify a guild VNUM to edit.\r\n"))
		return
	} else if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(buf1[0]))))) & int(uint16(int16(_ISdigit)))) == 0 {
		if C.strcasecmp(libc.CString("save"), &buf1[0]) != 0 {
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
		if d.Connected == CON_GEDIT {
			if d.Olc != nil && d.Olc.Number == room_vnum(number) {
				send_to_char(ch, libc.CString("That guild is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_gedit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if (func() zone_rnum {
		p := &d.Olc.Zone_num
		d.Olc.Zone_num = real_zone_by_thing(room_vnum(number))
		return *p
	}()) == zone_rnum(-1) {
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
		send_to_char(ch, libc.CString("Saving all guilds in zone %d.\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves guild info for zone %d."), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number)
		gedit_save_to_disk(int(d.Olc.Zone_num))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = room_vnum(number)
	if (func() guild_rnum {
		real_num = real_guild(guild_vnum(number))
		return real_num
	}()) != guild_rnum(-1) {
		gedit_setup_existing(d, int(real_num))
	} else {
		gedit_setup_new(d)
	}
	d.Connected = CON_GEDIT
	act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
	ch.Act[int(PLR_WRITING/32)] |= bitvector_t(1 << (int(PLR_WRITING % 32)))
	mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(d.Olc.Zone_num)))).Number, ch.Player_specials.Olc_zone)
}
func gedit_setup_new(d *descriptor_data) {
	var (
		i         int
		guilddata *guild_data
	)
	guilddata = new(guild_data)
	guilddata.Gm = -1
	guilddata.Open = 0
	guilddata.Close = 28
	guilddata.Charge = 1.0
	for i = 0; i < GW_ARRAY_MAX; i++ {
		guilddata.With_who[i] = 0
	}
	guilddata.Func = nil
	guilddata.Minlvl = 0
	guilddata.No_such_skill = C.strdup(libc.CString("%s Sorry, but I don't know that one."))
	guilddata.Not_enough_gold = C.strdup(libc.CString("%s Sorry, but I'm gonna need more zenni first."))
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if spell_info[i].Skilltype == (1<<1) && C.strcmp(spell_info[i].Name, libc.CString("!UNUSED!")) != 0 {
			guilddata.Skills[i] = 0
		}
	}
	for i = 0; i < NUM_FEATS_DEFINED; i++ {
		if feat_list[i].In_game != 0 {
			guilddata.Feats[i] = 0
		}
	}
	d.Olc.Guild = guilddata
	gedit_disp_menu(d)
}
func gedit_setup_existing(d *descriptor_data, rgm_num int) {
	d.Olc.Guild = new(guild_data)
	copy_guild(d.Olc.Guild, (*guild_data)(unsafe.Add(unsafe.Pointer(guild_index), unsafe.Sizeof(guild_data{})*uintptr(rgm_num))))
	gedit_disp_menu(d)
}
func gedit_select_skills_menu(d *descriptor_data) {
	var (
		i         int
		j         int = 0
		found     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	write_to_output(d, libc.CString("Skills known:\r\n"))
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if spell_info[i].Skilltype == (1<<1) && C.strcmp(spell_info[i].Name, libc.CString("!UNUSED!")) != 0 {
			write_to_output(d, libc.CString("@n[@c%-3s@n] %-3d %-20.20s  "), func() string {
				if (guilddata.Skills[i]) != 0 {
					return "YES"
				}
				return "NO"
			}(), i, spell_info[i].Name)
			j++
			found = 1
		}
		if found != 0 && (j%3) == 0 {
			found = 0
			write_to_output(d, libc.CString("\r\n"))
		}
	}
	write_to_output(d, libc.CString("\r\nEnter skill num, 0 to quit:  "))
	d.Olc.Mode = GEDIT_SELECT_SKILLS
}
func gedit_select_spells_menu(d *descriptor_data) {
	var (
		i         int
		j         int = 0
		found     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	write_to_output(d, libc.CString("Spells known:\r\n"))
	for i = 0; i <= SKILL_TABLE_SIZE; i++ {
		if (spell_info[i].Skilltype&(1<<0)) != 0 && C.strcmp(spell_info[i].Name, libc.CString("!UNUSED!")) != 0 {
			write_to_output(d, libc.CString("@n[@c%-3s@n] %-3d %-20.20s  "), func() string {
				if (guilddata.Skills[i]) != 0 {
					return "YES"
				}
				return "NO"
			}(), i, spell_info[i].Name)
			j++
			found = 1
		}
		if found != 0 && (j%3) == 0 {
			found = 0
			write_to_output(d, libc.CString("\r\n"))
		}
	}
	write_to_output(d, libc.CString("\r\nEnter spell num, 0 to quit:  "))
	d.Olc.Mode = GEDIT_SELECT_SPELLS
}
func gedit_select_feats_menu(d *descriptor_data) {
	var (
		i         int
		j         int = 0
		found     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	write_to_output(d, libc.CString("Feats known:\r\n"))
	for i = 0; i <= NUM_FEATS_DEFINED; i++ {
		if feat_list[i].In_game != 0 {
			write_to_output(d, libc.CString("@n[@c%-3s@n] %-3d %-20.20s  "), func() string {
				if (guilddata.Feats[i]) != 0 {
					return "YES"
				}
				return "NO"
			}(), i, feat_list[i].Name)
			j++
			found = 1
		}
		if found != 0 && (j%3) == 0 {
			found = 0
			write_to_output(d, libc.CString("\r\n"))
		}
	}
	write_to_output(d, libc.CString("\r\nEnter feat num, 0 to quit:  "))
	d.Olc.Mode = GEDIT_SELECT_FEATS
}
func gedit_select_lang_menu(d *descriptor_data) {
	var (
		i         int
		j         int = 0
		found     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	write_to_output(d, libc.CString("Skills known:\r\n"))
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if (spell_info[i].Skilltype&(1<<2)) != 0 && C.strcmp(spell_info[i].Name, libc.CString("!UNUSED!")) != 0 {
			write_to_output(d, libc.CString("@n[@c%-3s@n] %-3d %-20.20s  "), func() string {
				if (guilddata.Skills[i]) != 0 {
					return "YES"
				}
				return "NO"
			}(), i, spell_info[i].Name)
			j++
			found = 1
		}
		if found != 0 && (j%3) == 0 {
			found = 0
			write_to_output(d, libc.CString("\r\n"))
		}
	}
	write_to_output(d, libc.CString("\r\nEnter skill num, 0 to quit:  "))
	d.Olc.Mode = GEDIT_SELECT_LANGS
}
func gedit_select_wp_menu(d *descriptor_data) {
	var (
		i         int
		j         int = 0
		found     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	write_to_output(d, libc.CString("Skills known:\r\n"))
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		if (spell_info[i].Skilltype&(1<<3)) != 0 && C.strcmp(spell_info[i].Name, libc.CString("!UNUSED!")) != 0 {
			write_to_output(d, libc.CString("@n[@c%-3s@n] %-3d %-20.20s  "), func() string {
				if (guilddata.Skills[i]) != 0 {
					return "YES"
				}
				return "NO"
			}(), i, spell_info[i].Name)
			j++
			found = 1
		}
		if found != 0 && (j%3) == 0 {
			found = 0
			write_to_output(d, libc.CString("\r\n"))
		}
	}
	write_to_output(d, libc.CString("\r\nEnter skill num, 0 to quit:  "))
	d.Olc.Mode = GEDIT_SELECT_WPS
}
func gedit_no_train_menu(d *descriptor_data) {
	var (
		bits      [64936]byte
		i         int
		count     int = 0
		guilddata *guild_data
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	for i = 0; i < NUM_TRADERS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s   %s"), i+1, trade_letters[i], func() string {
			if (func() int {
				p := &count
				*p++
				return *p
			}() % 2) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	sprintbitarray(guilddata.With_who[:], trade_letters[:], int(64936), &bits[0])
	write_to_output(d, libc.CString("\r\nCurrent train flags: @c%s@n\r\nEnter choice, 0 to quit : "), &bits[0])
	d.Olc.Mode = GEDIT_NO_TRAIN
}
func gedit_disp_menu(d *descriptor_data) {
	var (
		guilddata *guild_data
		buf1      [64936]byte
	)
	guilddata = d.Olc.Guild
	clear_screen(d)
	sprintbitarray(guilddata.With_who[:], trade_letters[:], int(64936), &buf1[0])
	write_to_output(d, libc.CString("-- Guild Number: [@c%d@n]\r\n@g 0@n) Guild Master : [@c%d@n] @y%s\r\n@g 1@n) Doesn't know skill:\r\n @y%s\r\n@g 2@n) Player no gold:\r\n @y%s\r\n@g 3@n) Open   :  [@c%d@n]\r\n@g 4@n) Close  :  [@c%d@n]\r\n@g 5@n) Charge :  [@c%3.1f@n]\r\n@g 6@n) Minlvl :  [@c%d@n]\r\n@g 7@n) Who to Train:  @c%s\r\n@g 8@n) Feats Menu\r\n@g 9@n) Skills Menu\r\n@g B@n) Languages Menu\r\n@g Q@n) Quit\r\nEnter Choice : "), d.Olc.Number, func() mob_vnum {
		if guilddata.Gm == mob_rnum(-1) {
			return -1
		}
		return (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(guilddata.Gm)))).Vnum
	}(), func() string {
		if guilddata.Gm == mob_rnum(-1) {
			return "None"
		}
		return libc.GoString((*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(guilddata.Gm)))).Short_descr)
	}(), guilddata.No_such_skill, guilddata.Not_enough_gold, guilddata.Open, guilddata.Close, guilddata.Charge, guilddata.Minlvl, &buf1[0])
	d.Olc.Mode = GEDIT_MAIN_MENU
}
func gedit_parse(d *descriptor_data, arg *byte) {
	var i int
	if d.Olc.Mode > GEDIT_NUMERICAL_RESPONSE {
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*arg)))))&int(uint16(int16(_ISdigit)))) == 0 && (*arg == '-' && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*(*byte)(unsafe.Add(unsafe.Pointer(arg), 1)))))))&int(uint16(int16(_ISdigit)))) == 0) {
			write_to_output(d, libc.CString("Field must be numerical, try again : "))
			return
		}
	}
	switch d.Olc.Mode {
	case GEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			send_to_char(d.Character, libc.CString("Saving Guild to memory.\r\n"))
			gedit_save_internally(d)
			mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s edits guild %d"), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				gedit_save_to_disk(int(real_zone_by_thing(d.Olc.Number)))
				write_to_output(d, libc.CString("Guild %d saved to disk.\r\n"), d.Olc.Number)
			} else {
				write_to_output(d, libc.CString("Guild %d saved to memory.\r\n"), d.Olc.Number)
			}
			cleanup_olc(d, CLEANUP_STRUCTS)
			return
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
			return
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to save the guild? : "))
			return
		}
	case GEDIT_MAIN_MENU:
		i = 0
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save the changes to the Guild? (y/n) : "))
				d.Olc.Mode = GEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '0':
			d.Olc.Mode = GEDIT_TRAINER
			write_to_output(d, libc.CString("Enter vnum of guild master : "))
			return
		case '1':
			d.Olc.Mode = GEDIT_NO_SKILL
			i--
		case '2':
			d.Olc.Mode = GEDIT_NO_CASH
			i--
		case '3':
			d.Olc.Mode = GEDIT_OPEN
			write_to_output(d, libc.CString("When does this shop open (a day has 28 hours) ? "))
			i++
		case '4':
			d.Olc.Mode = GEDIT_CLOSE
			write_to_output(d, libc.CString("When does this shop close (a day has 28 hours) ? "))
			i++
		case '5':
			d.Olc.Mode = GEDIT_CHARGE
			i++
		case '6':
			d.Olc.Mode = GEDIT_MINLVL
			write_to_output(d, libc.CString("Minumum Level will Train: "))
			i++
			return
		case '7':
			d.Olc.Mode = GEDIT_NO_TRAIN
			gedit_no_train_menu(d)
			return
		case '8':
			d.Olc.Mode = GEDIT_SELECT_FEATS
			gedit_select_feats_menu(d)
			return
		case '9':
			d.Olc.Mode = GEDIT_SELECT_SKILLS
			gedit_select_skills_menu(d)
			return
		case 'b':
			fallthrough
		case 'B':
			d.Olc.Mode = GEDIT_SELECT_LANGS
			gedit_select_lang_menu(d)
			return
		default:
			gedit_disp_menu(d)
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
	case GEDIT_NO_SKILL:
		gedit_modify_string(&d.Olc.Guild.No_such_skill, arg)
	case GEDIT_NO_CASH:
		gedit_modify_string(&d.Olc.Guild.Not_enough_gold, arg)
	case GEDIT_TRAINER:
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*arg))))) & int(uint16(int16(_ISdigit)))) != 0 {
			i = libc.Atoi(libc.GoString(arg))
			if (func() int {
				i = libc.Atoi(libc.GoString(arg))
				return i
			}()) != -1 {
				if (func() int {
					i = int(real_mobile(mob_vnum(i)))
					return i
				}()) == int(-1) {
					write_to_output(d, libc.CString("That mobile does not exist, try again : "))
					return
				}
			}
			d.Olc.Guild.Gm = mob_rnum(i)
			if i == -1 {
				break
			}
			if libc.FuncAddr((*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Func) != libc.FuncAddr(guild) {
				d.Olc.Guild.Func = (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Func
			} else {
				d.Olc.Guild.Func = nil
			}
			(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Func = guild
			break
		} else {
			write_to_output(d, libc.CString("Invalid response.\r\n"))
			gedit_disp_menu(d)
			return
		}
		fallthrough
	case GEDIT_OPEN:
		d.Olc.Guild.Open = MIN(28, MAX(libc.Atoi(libc.GoString(arg)), 0))
	case GEDIT_CLOSE:
		d.Olc.Guild.Close = MIN(28, MAX(libc.Atoi(libc.GoString(arg)), 0))
	case GEDIT_CHARGE:
		__isoc99_sscanf(arg, libc.CString("%f"), &d.Olc.Guild.Charge)
	case GEDIT_NO_TRAIN:
		if (func() int {
			i = MIN(int(NUM_TRADERS-1), MAX(libc.Atoi(libc.GoString(arg)), 0))
			return i
		}()) > 0 {
			d.Olc.Guild.With_who[(i-1)/32] = d.Olc.Guild.With_who[(i-1)/32] ^ 1<<((i-1)%32)
			gedit_no_train_menu(d)
			return
		}
	case GEDIT_MINLVL:
		d.Olc.Guild.Minlvl = MAX(libc.Atoi(libc.GoString(arg)), 0)
	case GEDIT_SELECT_SPELLS:
		i = libc.Atoi(libc.GoString(arg))
		if i == 0 {
			break
		}
		i = MAX(1, MIN(i, SKILL_TABLE_SIZE))
		d.Olc.Guild.Skills[i] = int(libc.BoolToInt((d.Olc.Guild.Skills[i]) == 0))
		gedit_select_spells_menu(d)
		return
	case GEDIT_SELECT_FEATS:
		i = libc.Atoi(libc.GoString(arg))
		if i == 0 {
			break
		}
		i = MAX(1, MIN(i, NUM_FEATS_DEFINED))
		d.Olc.Guild.Feats[i] = int(libc.BoolToInt((d.Olc.Guild.Feats[i]) == 0))
		gedit_select_feats_menu(d)
		return
	case GEDIT_SELECT_SKILLS:
		i = libc.Atoi(libc.GoString(arg))
		if i == 0 {
			break
		}
		i = MAX(1, MIN(i, SKILL_TABLE_SIZE))
		d.Olc.Guild.Skills[i] = int(libc.BoolToInt((d.Olc.Guild.Skills[i]) == 0))
		gedit_select_skills_menu(d)
		return
	case GEDIT_SELECT_WPS:
		i = libc.Atoi(libc.GoString(arg))
		if i == 0 {
			break
		}
		i = MAX(1, MIN(i, SKILL_TABLE_SIZE))
		d.Olc.Guild.Skills[i] = int(libc.BoolToInt((d.Olc.Guild.Skills[i]) == 0))
		gedit_select_wp_menu(d)
		return
	case GEDIT_SELECT_LANGS:
		i = libc.Atoi(libc.GoString(arg))
		if i == 0 {
			break
		}
		i = MAX(1, MIN(i, SKILL_TABLE_SIZE))
		d.Olc.Guild.Skills[i] = int(libc.BoolToInt((d.Olc.Guild.Skills[i]) == 0))
		gedit_select_lang_menu(d)
		return
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: gedit_parse(): Reached default case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
	d.Olc.Value = 1
	gedit_disp_menu(d)
}
