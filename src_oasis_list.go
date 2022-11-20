package main

import "github.com/gotranspile/cxgo/runtime/libc"

func do_oasis_list(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		rzone int = int(-1)
		vmin  int = int(-1)
		vmax  int = int(-1)
		smin  [2048]byte
		smax  [2048]byte
	)
	two_arguments(argument, &smin[0], &smax[0])
	if subcmd == SCMD_OASIS_ZLIST {
		if smin[0] != 0 && is_number(&smin[0]) {
			print_zone(ch, libc.Atoi(libc.GoString(&smin[0])))
		} else {
			list_zones(ch)
		}
		return
	}
	if smin[0] == 0 || smin[0] == '.' {
		rzone = world[ch.In_room].Zone
	} else if smax[0] == 0 {
		rzone = real_zone(libc.Atoi(libc.GoString(&smin[0])))
		if rzone == int(-1) {
			send_to_char(ch, libc.CString("Sorry, there's no zone with that number\r\n"))
			return
		}
	} else {
		vmin = libc.Atoi(libc.GoString(&smin[0]))
		vmax = libc.Atoi(libc.GoString(&smax[0]))
		if vmin+500 < vmax {
			send_to_char(ch, libc.CString("Really? Over 500?! You need to view that many at once? Come on...\r\n"))
			return
		}
		if vmin > vmax {
			send_to_char(ch, libc.CString("List from %d to %d - Aren't we funny today!\r\n"), vmin, vmax)
			return
		}
	}
	switch subcmd {
	case SCMD_OASIS_RLIST:
		list_rooms(ch, rzone, vmin, vmax)
	case SCMD_OASIS_OLIST:
		list_objects(ch, rzone, vmin, vmax)
	case SCMD_OASIS_MLIST:
		list_mobiles(ch, rzone, vmin, vmax)
	case SCMD_OASIS_TLIST:
		list_triggers(ch, rzone, vmin, vmax)
	case SCMD_OASIS_SLIST:
		list_shops(ch, rzone, vmin, vmax)
	case SCMD_OASIS_GLIST:
		list_guilds(ch, rzone, vmin, vmax)
	default:
		send_to_char(ch, libc.CString("You can't list that!\r\n"))
		mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("SYSERR: do_oasis_list: Unknown list option: %d"), subcmd)
	}
}
func do_oasis_links(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		zrnum   int
		zvnum   int
		nr      int
		to_room int
		first   int
		last    int
		j       int
		arg     [2048]byte
	)
	skip_spaces(&argument)
	one_argument(argument, &arg[0])
	if libc.StrCmp(&arg[0], libc.CString(".")) == 0 || (arg[0] == 0) {
		zrnum = world[ch.In_room].Zone
		zvnum = zone_table[zrnum].Number
	} else {
		zvnum = libc.Atoi(libc.GoString(&arg[0]))
		zrnum = real_zone(zvnum)
	}
	if zrnum == int(-1) || zvnum == int(-1) {
		send_to_char(ch, libc.CString("No zone was found with that number.\n\r"))
		return
	}
	last = zone_table[zrnum].Top
	first = zone_table[zrnum].Bot
	send_to_char(ch, libc.CString("Zone %d is linked to the following zones:\r\n"), zvnum)
	for nr = 0; nr <= top_of_world && int(libc.BoolToInt(GET_ROOM_VNUM(nr))) <= last; nr++ {
		if int(libc.BoolToInt(GET_ROOM_VNUM(nr))) >= first {
			for j = 0; j < NUM_OF_DIRS; j++ {
				if world[nr].Dir_option[j] != nil {
					to_room = world[nr].Dir_option[j].To_room
					if to_room != int(-1) && zrnum != world[to_room].Zone {
						send_to_char(ch, libc.CString("%3d %-30s at %5d (%-5s) ---> %5d\r\n"), zone_table[world[to_room].Zone].Number, zone_table[world[to_room].Zone].Name, GET_ROOM_VNUM(nr), dirs[j], world[to_room].Number)
					}
				}
			}
		}
	}
}
func list_rooms(ch *char_data, rnum int, vmin int, vmax int) {
	var (
		i       int
		j       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != int(-1) {
		bottom = zone_table[rnum].Bot
		top = zone_table[rnum].Top
	} else {
		bottom = vmin
		top = vmax
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Room Name                                Exits\r\n----- ------- ---------------------------------------- -----@n\r\n"))
	if top_of_world == 0 {
		return
	}
	for i = 0; i <= top_of_world; i++ {
		if world[i].Number >= bottom && world[i].Number <= top {
			counter++
			send_to_char(ch, libc.CString("%4d) [@g%-5d@n] @[1]%-*s@n %s"), counter, world[i].Number, count_color_chars(world[i].Name)+44, world[i].Name, func() string {
				if world[i].Proto_script != nil {
					return "[TRIG] "
				}
				return ""
			}())
			for j = 0; j < NUM_OF_DIRS; j++ {
				if (world[i].Dir_option[j]) == nil {
					continue
				}
				if (world[i].Dir_option[j]).To_room == int(-1) {
					continue
				}
				if world[(world[i].Dir_option[j]).To_room].Zone != world[i].Zone {
					send_to_char(ch, libc.CString("(@y%d@n)"), world[(world[i].Dir_option[j]).To_room].Number)
				}
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("No rooms found for zone/range specified.\r\n"))
	}
}
func list_mobiles(ch *char_data, rnum int, vmin int, vmax int) {
	var (
		i       int
		bottom  int
		top     int
		counter int = 0
		admg    int
	)
	_ = admg
	if rnum != int(-1) {
		bottom = zone_table[rnum].Bot
		top = zone_table[rnum].Top
	} else {
		bottom = vmin
		top = vmax
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Mobile Name                    Race      Class     Level\r\n----- ------- -------------------------      --------- --------- -----\r\n"))
	if top_of_mobt == 0 {
		return
	}
	for i = 0; i <= top_of_mobt; i++ {
		if mob_index[i].Vnum >= bottom && mob_index[i].Vnum <= top {
			counter++
			admg = int((float64(int(mob_proto[i].Mob_specials.Damsizedice)+1) / 2.0) * float64(mob_proto[i].Mob_specials.Damnodice))
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n] @[3]%-*s @C%-9s @c%-9s @y[%4d]@n %s\r\n"), counter, mob_index[i].Vnum, count_color_chars(mob_proto[i].Short_descr)+30, mob_proto[i].Short_descr, pc_race_types[mob_proto[i].Race], pc_class_types[mob_proto[i].Chclass], mob_proto[i].Level+mob_proto[i].Level_adj+mob_proto[i].Race_level, func() string {
				if mob_proto[i].Proto_script != nil {
					return " [TRIG]"
				}
				return ""
			}())
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("None found.\r\n"))
	}
}
func list_objects(ch *char_data, rnum int, vmin int, vmax int) {
	var (
		i       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != int(-1) {
		bottom = zone_table[rnum].Bot
		top = zone_table[rnum].Top
	} else {
		bottom = vmin
		top = vmax
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Object Name                                  Object Type\r\n----- ------- -------------------------------------------- ----------------\r\n"))
	if top_of_objt == 0 {
		return
	}
	for i = 0; i <= top_of_objt; i++ {
		if obj_index[i].Vnum >= bottom && obj_index[i].Vnum <= top {
			counter++
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n] @[2]%-*s @y[%s]@n%s\r\n"), counter, obj_index[i].Vnum, count_color_chars(obj_proto[i].Short_description)+44, obj_proto[i].Short_description, item_types[obj_proto[i].Type_flag], func() string {
				if obj_proto[i].Proto_script != nil {
					return " [TRIG]"
				}
				return ""
			}())
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("None found.\r\n"))
	}
}
func list_shops(ch *char_data, rnum int, vmin int, vmax int) {
	var (
		i       int
		j       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != int(-1) {
		bottom = zone_table[rnum].Bot
		top = zone_table[rnum].Top
	} else {
		bottom = vmin
		top = vmax
	}
	send_to_char(ch, libc.CString("Index VNum    Shop Room(s)\r\n----- ------- ---------------------------------------------\r\n"))
	for i = 0; i <= top_shop; i++ {
		if shop_index[i].Vnum >= bottom && shop_index[i].Vnum <= top {
			counter++
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n]"), counter, shop_index[i].Vnum)
			for j = 0; (shop_index[i].In_room[j]) != int(-1); j++ {
				send_to_char(ch, libc.CString("%s@c[@y%d@c]@n"), func() string {
					if j > 0 && j%8 == 0 {
						return "\r\n              "
					}
					return " "
				}(), shop_index[i].In_room[j])
			}
			if j == 0 {
				send_to_char(ch, libc.CString("@cNone.@n"))
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("None found.\r\n"))
	}
}
func list_zones(ch *char_data) {
	var i int
	send_to_char(ch, libc.CString("VNum  Zone Name                      Builder(s)\r\n----- ------------------------------ --------------------------------------\r\n"))
	if top_of_zone_table == 0 {
		return
	}
	for i = 0; i <= top_of_zone_table; i++ {
		send_to_char(ch, libc.CString("[@g%3d@n] @c%-*s @y%-1s@n\r\n"), zone_table[i].Number, count_color_chars(zone_table[i].Name)+30, zone_table[i].Name, func() *byte {
			if zone_table[i].Builders != nil {
				return zone_table[i].Builders
			}
			return libc.CString("None.")
		}())
	}
}
func print_zone(ch *char_data, vnum int) {
	var (
		rnum          int
		size_rooms    int
		size_objects  int
		size_mobiles  int
		i             int
		size_guilds   int
		size_triggers int
		size_shops    int
		top           int
		bottom        int
		largest_table int
		bits          [64936]byte
	)
	if (func() int {
		rnum = real_zone(vnum)
		return rnum
	}()) == int(-1) {
		send_to_char(ch, libc.CString("Zone #%d does not exist in the database.\r\n"), vnum)
		return
	}
	sprintbitarray(zone_table[rnum].Zone_flags[:], zone_bits[:], ZF_ARRAY_MAX, &bits[0])
	if top_of_world >= top_of_objt && top_of_world >= top_of_mobt {
		largest_table = top_of_world
	} else if top_of_objt >= top_of_mobt && top_of_objt >= top_of_world {
		largest_table = top_of_objt
	} else {
		largest_table = top_of_mobt
	}
	size_rooms = 0
	size_objects = 0
	size_mobiles = 0
	top = zone_table[rnum].Top
	bottom = zone_table[rnum].Bot
	size_shops = 0
	size_triggers = 0
	size_guilds = 0
	size_shops = 0
	for i = 0; i <= largest_table; i++ {
		if i <= top_of_world {
			if world[i].Zone == rnum {
				size_rooms++
			}
		}
		if i <= top_of_objt {
			if obj_index[i].Vnum >= bottom && obj_index[i].Vnum <= top {
				size_objects++
			}
		}
		if i <= top_of_mobt {
			if mob_index[i].Vnum >= bottom && mob_index[i].Vnum <= top {
				size_mobiles++
			}
		}
		size_shops = count_shops(bottom, top)
		if i < top_of_trigt {
			if trig_index[i].Vnum >= bottom && trig_index[i].Vnum <= top {
				size_triggers++
			}
		}
		size_guilds = count_guilds(bottom, top)
	}
	send_to_char(ch, libc.CString("@gVirtual Number = @c%d\r\n@gName of zone   = @c%s\r\n@gBuilders       = @c%s\r\n@gLifespan       = @c%d\r\n@gAge            = @c%d\r\n@gBottom of Zone = @c%d\r\n@gTop of Zone    = @c%d\r\n@gReset Mode     = @c%s\r\n@gMin Level      = @c%d\r\n@gMax Level      = @c%d\r\n@gZone Flags     = @c%s\r\n@gSize\r\n@g   Rooms       = @c%d\r\n@g   Objects     = @c%d\r\n@g   Mobiles     = @c%d\r\n@g   Shops       = @c%d\r\n@g   Triggers    = @c%d\r\n@g   Guilds      = @c%d@n\r\n"), zone_table[rnum].Number, zone_table[rnum].Name, zone_table[rnum].Builders, zone_table[rnum].Lifespan, zone_table[rnum].Age, zone_table[rnum].Bot, zone_table[rnum].Top, func() string {
		if zone_table[rnum].Reset_mode != 0 {
			if zone_table[rnum].Reset_mode == 1 {
				return "Reset when no players are in zone."
			}
			return "Normal reset."
		}
		return "Never reset"
	}(), zone_table[rnum].Min_level, zone_table[rnum].Max_level, &bits[0], size_rooms, size_objects, size_mobiles, size_shops, size_triggers, size_guilds)
}
func list_triggers(ch *char_data, rnum int, vmin int, vmax int) {
	var (
		i        int
		bottom   int
		top      int
		counter  int = 0
		trgtypes [256]byte
	)
	if rnum != int(-1) {
		bottom = zone_table[rnum].Bot
		top = zone_table[rnum].Top
	} else {
		bottom = vmin
		top = vmax
	}
	send_to_char(ch, libc.CString("Index VNum    Trigger Name                        Type\r\n----- ------- -------------------------------------------------------\r\n"))
	for i = 0; i < top_of_trigt; i++ {
		if trig_index[i].Vnum >= bottom && trig_index[i].Vnum <= top {
			counter++
			send_to_char(ch, libc.CString("%4d) [@g%5d@n] @[1]%-45.45s "), counter, trig_index[i].Vnum, trig_index[i].Proto.Name)
			if int(trig_index[i].Proto.Attach_type) == OBJ_TRIGGER {
				sprintbit(uint32(int32(trig_index[i].Proto.Trigger_type)), otrig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("obj @y%s@n\r\n"), &trgtypes[0])
			} else if int(trig_index[i].Proto.Attach_type) == WLD_TRIGGER {
				sprintbit(uint32(int32(trig_index[i].Proto.Trigger_type)), wtrig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("wld @y%s@n\r\n"), &trgtypes[0])
			} else {
				sprintbit(uint32(int32(trig_index[i].Proto.Trigger_type)), trig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("mob @y%s@n\r\n"), &trgtypes[0])
			}
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("No triggers found from %d to %d.\r\n"), vmin, vmax)
	}
}
