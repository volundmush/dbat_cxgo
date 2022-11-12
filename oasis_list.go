package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func do_oasis_list(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		rzone zone_rnum = zone_rnum(-1)
		vmin  room_rnum = room_rnum(-1)
		vmax  room_rnum = room_rnum(-1)
		smin  [2048]byte
		smax  [2048]byte
	)
	two_arguments(argument, &smin[0], &smax[0])
	if subcmd == SCMD_OASIS_ZLIST {
		if smin != nil && smin[0] != 0 && is_number(&smin[0]) != 0 {
			print_zone(ch, zone_vnum(libc.Atoi(libc.GoString(&smin[0]))))
		} else {
			list_zones(ch)
		}
		return
	}
	if smin[0] == 0 || smin[0] == '.' {
		rzone = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone
	} else if smax[0] == 0 {
		rzone = real_zone(zone_vnum(libc.Atoi(libc.GoString(&smin[0]))))
		if rzone == zone_rnum(-1) {
			send_to_char(ch, libc.CString("Sorry, there's no zone with that number\r\n"))
			return
		}
	} else {
		vmin = room_rnum(libc.Atoi(libc.GoString(&smin[0])))
		vmax = room_rnum(libc.Atoi(libc.GoString(&smax[0])))
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
		list_rooms(ch, rzone, room_vnum(vmin), room_vnum(vmax))
	case SCMD_OASIS_OLIST:
		list_objects(ch, rzone, obj_vnum(vmin), obj_vnum(vmax))
	case SCMD_OASIS_MLIST:
		list_mobiles(ch, rzone, mob_vnum(vmin), mob_vnum(vmax))
	case SCMD_OASIS_TLIST:
		list_triggers(ch, rzone, trig_vnum(vmin), trig_vnum(vmax))
	case SCMD_OASIS_SLIST:
		list_shops(ch, rzone, shop_vnum(vmin), shop_vnum(vmax))
	case SCMD_OASIS_GLIST:
		list_guilds(ch, rzone, guild_vnum(vmin), guild_vnum(vmax))
	default:
		send_to_char(ch, libc.CString("You can't list that!\r\n"))
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("SYSERR: do_oasis_list: Unknown list option: %d"), subcmd)
	}
}
func do_oasis_links(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		zrnum   zone_rnum
		zvnum   zone_vnum
		nr      room_rnum
		to_room room_rnum
		first   int
		last    int
		j       int
		arg     [2048]byte
	)
	skip_spaces(&argument)
	one_argument(argument, &arg[0])
	if C.strcmp(&arg[0], libc.CString(".")) == 0 || (arg == nil || arg[0] == 0) {
		zrnum = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone
		zvnum = (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Number
	} else {
		zvnum = zone_vnum(libc.Atoi(libc.GoString(&arg[0])))
		zrnum = real_zone(zvnum)
	}
	if zrnum == zone_rnum(-1) || zvnum == zone_vnum(-1) {
		send_to_char(ch, libc.CString("No zone was found with that number.\n\r"))
		return
	}
	last = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Top)
	first = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(zrnum)))).Bot)
	send_to_char(ch, libc.CString("Zone %d is linked to the following zones:\r\n"), zvnum)
	for nr = 0; nr <= top_of_world && (func() room_vnum {
		if nr != room_rnum(-1) && nr <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Number
		}
		return -1
	}()) <= room_vnum(last); nr++ {
		if (func() room_vnum {
			if nr != room_rnum(-1) && nr <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Number
			}
			return -1
		}()) >= room_vnum(first) {
			for j = 0; j < NUM_OF_DIRS; j++ {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Dir_option[j] != nil {
					to_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Dir_option[j].To_room
					if to_room != room_rnum(-1) && zrnum != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(to_room)))).Zone {
						send_to_char(ch, libc.CString("%3d %-30s at %5d (%-5s) ---> %5d\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(to_room)))).Zone)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(to_room)))).Zone)))).Name, func() room_vnum {
							if nr != room_rnum(-1) && nr <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(nr)))).Number
							}
							return -1
						}(), dirs[j], (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(to_room)))).Number)
					}
				}
			}
		}
	}
}
func list_rooms(ch *char_data, rnum zone_rnum, vmin zone_vnum, vmax zone_vnum) {
	var (
		i       int
		j       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != zone_rnum(-1) {
		bottom = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot)
		top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Room Name                                Exits\r\n----- ------- ---------------------------------------- -----@n\r\n"))
	if top_of_world == 0 {
		return
	}
	for i = 0; i <= int(top_of_world); i++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number >= room_vnum(bottom) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number <= room_vnum(top) {
			counter++
			send_to_char(ch, libc.CString("%4d) [@g%-5d@n] @[1]%-*s@n %s"), counter, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Number, count_color_chars((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name)+44, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Name, func() string {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Proto_script != nil {
					return "[TRIG] "
				}
				return ""
			}())
			for j = 0; j < NUM_OF_DIRS; j++ {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]) == nil {
					continue
				}
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room == room_rnum(-1) {
					continue
				}
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Zone {
					send_to_char(ch, libc.CString("(@y%d@n)"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Dir_option[j]).To_room)))).Number)
				}
			}
			send_to_char(ch, libc.CString("\r\n"))
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("No rooms found for zone/range specified.\r\n"))
	}
}
func list_mobiles(ch *char_data, rnum zone_rnum, vmin zone_vnum, vmax zone_vnum) {
	var (
		i       int
		bottom  int
		top     int
		counter int = 0
		admg    int
	)
	_ = admg
	if rnum != zone_rnum(-1) {
		bottom = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot)
		top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Mobile Name                    Race      Class     Level\r\n----- ------- -------------------------      --------- --------- -----\r\n"))
	if top_of_mobt == 0 {
		return
	}
	for i = 0; i <= int(top_of_mobt); i++ {
		if (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum(top) {
			counter++
			admg = int((float64((*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Mob_specials.Damsizedice+1) / 2.0) * float64((*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Mob_specials.Damnodice))
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n] @[3]%-*s @C%-9s @c%-9s @y[%4d]@n %s\r\n"), counter, (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum, count_color_chars((*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Short_descr)+30, (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Short_descr, pc_race_types[(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Race], pc_class_types[(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Chclass], (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Level+(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Level_adj+(*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Race_level, func() string {
				if (*(*char_data)(unsafe.Add(unsafe.Pointer(mob_proto), unsafe.Sizeof(char_data{})*uintptr(i)))).Proto_script != nil {
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
func list_objects(ch *char_data, rnum zone_rnum, vmin room_vnum, vmax room_vnum) {
	var (
		i       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != zone_rnum(-1) {
		bottom = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot)
		top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("@nIndex VNum    Object Name                                  Object Type\r\n----- ------- -------------------------------------------- ----------------\r\n"))
	if top_of_objt == 0 {
		return
	}
	for i = 0; i <= int(top_of_objt); i++ {
		if (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum(top) {
			counter++
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n] @[2]%-*s @y[%s]@n%s\r\n"), counter, (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum, count_color_chars((*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Short_description)+44, (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Short_description, item_types[(*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Type_flag], func() string {
				if (*(*obj_data)(unsafe.Add(unsafe.Pointer(obj_proto), unsafe.Sizeof(obj_data{})*uintptr(i)))).Proto_script != nil {
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
func list_shops(ch *char_data, rnum zone_rnum, vmin shop_vnum, vmax shop_vnum) {
	var (
		i       int
		j       int
		bottom  int
		top     int
		counter int = 0
	)
	if rnum != zone_rnum(-1) {
		bottom = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot)
		top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("Index VNum    Shop Room(s)\r\n----- ------- ---------------------------------------------\r\n"))
	for i = 0; i <= top_shop; i++ {
		if (*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).Vnum >= room_vnum(bottom) && (*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).Vnum <= room_vnum(top) {
			counter++
			send_to_char(ch, libc.CString("@g%4d@n) [@g%-5d@n]"), counter, (*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).Vnum)
			for j = 0; (*(*room_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).In_room), unsafe.Sizeof(room_vnum(0))*uintptr(j)))) != room_vnum(-1); j++ {
				send_to_char(ch, libc.CString("%s@c[@y%d@c]@n"), func() string {
					if j > 0 && j%8 == 0 {
						return "\r\n              "
					}
					return " "
				}(), *(*room_vnum)(unsafe.Add(unsafe.Pointer((*(*shop_data)(unsafe.Add(unsafe.Pointer(shop_index), unsafe.Sizeof(shop_data{})*uintptr(i)))).In_room), unsafe.Sizeof(room_vnum(0))*uintptr(j))))
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
	for i = 0; i <= int(top_of_zone_table); i++ {
		send_to_char(ch, libc.CString("[@g%3d@n] @c%-*s @y%-1s@n\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Number, count_color_chars((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Name)+30, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Name, func() *byte {
			if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Builders != nil {
				return (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(i)))).Builders
			}
			return libc.CString("None.")
		}())
	}
}
func print_zone(ch *char_data, vnum zone_vnum) {
	var (
		rnum          zone_rnum
		size_rooms    int
		size_objects  int
		size_mobiles  int
		i             int
		size_guilds   int
		size_triggers int
		size_shops    int
		top           room_vnum
		bottom        room_vnum
		largest_table int
		bits          [64936]byte
	)
	if (func() zone_rnum {
		rnum = real_zone(vnum)
		return rnum
	}()) == zone_rnum(-1) {
		send_to_char(ch, libc.CString("Zone #%d does not exist in the database.\r\n"), vnum)
		return
	}
	sprintbitarray((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Zone_flags[:], zone_bits[:], ZF_ARRAY_MAX, &bits[0])
	if top_of_world >= room_rnum(top_of_objt) && top_of_world >= room_rnum(top_of_mobt) {
		largest_table = int(top_of_world)
	} else if top_of_objt >= obj_rnum(top_of_mobt) && top_of_objt >= obj_rnum(top_of_world) {
		largest_table = int(top_of_objt)
	} else {
		largest_table = int(top_of_mobt)
	}
	size_rooms = 0
	size_objects = 0
	size_mobiles = 0
	top = (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top
	bottom = (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot
	size_shops = 0
	size_triggers = 0
	size_guilds = 0
	size_shops = 0
	for i = 0; i <= largest_table; i++ {
		if i <= int(top_of_world) {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Zone == rnum {
				size_rooms++
			}
		}
		if i <= int(top_of_objt) {
			if (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum(top) {
				size_objects++
			}
		}
		if i <= int(top_of_mobt) {
			if (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(i)))).Vnum <= mob_vnum(top) {
				size_mobiles++
			}
		}
		size_shops = count_shops(shop_vnum(bottom), shop_vnum(top))
		if i < top_of_trigt {
			if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum <= mob_vnum(top) {
				size_triggers++
			}
		}
		size_guilds = count_guilds(guild_vnum(bottom), guild_vnum(top))
	}
	send_to_char(ch, libc.CString("@gVirtual Number = @c%d\r\n@gName of zone   = @c%s\r\n@gBuilders       = @c%s\r\n@gLifespan       = @c%d\r\n@gAge            = @c%d\r\n@gBottom of Zone = @c%d\r\n@gTop of Zone    = @c%d\r\n@gReset Mode     = @c%s\r\n@gMin Level      = @c%d\r\n@gMax Level      = @c%d\r\n@gZone Flags     = @c%s\r\n@gSize\r\n@g   Rooms       = @c%d\r\n@g   Objects     = @c%d\r\n@g   Mobiles     = @c%d\r\n@g   Shops       = @c%d\r\n@g   Triggers    = @c%d\r\n@g   Guilds      = @c%d@n\r\n"), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Number, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Name, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Builders, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Lifespan, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Age, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top, func() string {
		if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Reset_mode != 0 {
			if (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Reset_mode == 1 {
				return "Reset when no players are in zone."
			}
			return "Normal reset."
		}
		return "Never reset"
	}(), (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Min_level, (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Max_level, &bits[0], size_rooms, size_objects, size_mobiles, size_shops, size_triggers, size_guilds)
}
func list_triggers(ch *char_data, rnum zone_rnum, vmin trig_vnum, vmax trig_vnum) {
	var (
		i        int
		bottom   int
		top      int
		counter  int = 0
		trgtypes [256]byte
	)
	if rnum != zone_rnum(-1) {
		bottom = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Bot)
		top = int((*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rnum)))).Top)
	} else {
		bottom = int(vmin)
		top = int(vmax)
	}
	send_to_char(ch, libc.CString("Index VNum    Trigger Name                        Type\r\n----- ------- -------------------------------------------------------\r\n"))
	for i = 0; i < top_of_trigt; i++ {
		if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum >= mob_vnum(bottom) && (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum <= mob_vnum(top) {
			counter++
			send_to_char(ch, libc.CString("%4d) [@g%5d@n] @[1]%-45.45s "), counter, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Vnum, (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Name)
			if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Attach_type == OBJ_TRIGGER {
				sprintbit(bitvector_t((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Trigger_type), otrig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("obj @y%s@n\r\n"), &trgtypes[0])
			} else if (*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Attach_type == WLD_TRIGGER {
				sprintbit(bitvector_t((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Trigger_type), wtrig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("wld @y%s@n\r\n"), &trgtypes[0])
			} else {
				sprintbit(bitvector_t((*(**index_data)(unsafe.Add(unsafe.Pointer(trig_index), unsafe.Sizeof((*index_data)(nil))*uintptr(i)))).Proto.Trigger_type), trig_types[:], &trgtypes[0], uint64(256))
				send_to_char(ch, libc.CString("mob @y%s@n\r\n"), &trgtypes[0])
			}
		}
	}
	if counter == 0 {
		send_to_char(ch, libc.CString("No triggers found from %d to %d.\r\n"), vmin, vmax)
	}
}
