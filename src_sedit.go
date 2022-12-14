package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func sedit_save_internally(d *descriptor_data) {
	d.Olc.Shop.Vnum = d.Olc.Number
	add_shop(d.Olc.Shop)
}
func sedit_save_to_disk(num int) {
	save_shops(num)
}
func do_oasis_sedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		number   int = int(-1)
		save     int = 0
		real_num int
		d        *descriptor_data
		buf3     *byte
	)
	_ = buf3
	var buf1 [2048]byte
	var buf2 [2048]byte
	buf3 = two_arguments(argument, &buf1[0], &buf2[0])
	if buf1[0] == 0 {
		send_to_char(ch, libc.CString("Specify a shop VNUM to edit.\r\n"))
		return
	} else if !unicode.IsDigit(rune(buf1[0])) {
		if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) != 0 {
			send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
			return
		}
		save = 1
		if is_number(&buf2[0]) {
			number = libc.Atoi(libc.GoString(&buf2[0]))
		} else if ch.Player_specials.Olc_zone > 0 {
			var zlok int
			if (func() int {
				zlok = real_zone(ch.Player_specials.Olc_zone)
				return zlok
			}()) == int(-1) {
				number = -1
			} else {
				number = genolc_zone_bottom(zlok)
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
		if d.Connected == CON_SEDIT {
			if d.Olc != nil && d.Olc.Number == number {
				send_to_char(ch, libc.CString("That shop is currently being edited by %s.\r\n"), PERS(d.Character, ch))
				return
			}
		}
	}
	d = ch.Desc
	if d.Olc != nil {
		mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("SYSERR: do_oasis_sedit: Player already had olc structure."))
		libc.Free(unsafe.Pointer(d.Olc))
	}
	d.Olc = new(oasis_olc_data)
	if save != 0 {
		d.Olc.Zone_num = real_zone(number)
	} else {
		d.Olc.Zone_num = real_zone_by_thing(number)
	}
	if d.Olc.Zone_num == int(-1) {
		send_to_char(ch, libc.CString("Sorry, there is no zone for that number!\r\n"))
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if !can_edit_zone(ch, d.Olc.Zone_num) {
		send_cannot_edit(ch, zone_table[d.Olc.Zone_num].Number)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	if save != 0 {
		send_to_char(ch, libc.CString("Saving all shops in zone %d.\r\n"), zone_table[d.Olc.Zone_num].Number)
		mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(ch.Player_specials.Invis_level))), 1, libc.CString("OLC: %s saves shop info for zone %d."), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number)
		save_shops(d.Olc.Zone_num)
		libc.Free(unsafe.Pointer(d.Olc))
		d.Olc = nil
		return
	}
	d.Olc.Number = number
	if (func() int {
		real_num = real_shop(number)
		return real_num
	}()) != int(-1) {
		sedit_setup_existing(d, real_num)
	} else {
		sedit_setup_new(d)
	}
	sedit_disp_menu(d)
	d.Connected = CON_SEDIT
	act(libc.CString("$n starts using OLC."), 1, d.Character, nil, nil, TO_ROOM)
	SET_BIT_AR(ch.Act[:], PLR_WRITING)
	mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("OLC: %s starts editing zone %d allowed zone %d"), GET_NAME(ch), zone_table[d.Olc.Zone_num].Number, ch.Player_specials.Olc_zone)
}
func sedit_setup_new(d *descriptor_data) {
	var shop *shop_data
	shop = new(shop_data)
	shop.Keeper = -1
	shop.Close1 = 28
	shop.Profit_buy = 1.0
	shop.Profit_sell = 1.0
	shop.No_such_item1 = libc.CString("%s Sorry, I don't stock that item.")
	shop.No_such_item2 = libc.CString("%s You don't seem to have that.")
	shop.Missing_cash1 = libc.CString("%s I can't afford that!")
	shop.Missing_cash2 = libc.CString("%s You are too poor!")
	shop.Do_not_buy = libc.CString("%s I don't trade in such items.")
	shop.Message_buy = libc.CString("%s That'll be %d zenni, thanks.")
	shop.Message_sell = libc.CString("%s I'll give you %d zenni for that.")
	shop.Producing = make([]int, 0)
	shop.Producing[0] = -1
	shop.In_room = make([]int, 0)
	shop.In_room[0] = -1
	shop.Type = make([]shop_buy_data, 0)
	(shop.Type[0]).Type = -1
	SET_BIT_AR(shop.With_who[:], TRADE_NOBROKEN)
	d.Olc.Shop = shop
}
func sedit_setup_existing(d *descriptor_data, rshop_num int) {
	d.Olc.Shop = new(shop_data)
	copy_shop(d.Olc.Shop, &shop_index[rshop_num], 0)
}
func sedit_products_menu(d *descriptor_data) {
	var (
		shop *shop_data
		i    int
	)
	shop = d.Olc.Shop
	clear_screen(d)
	write_to_output(d, libc.CString("##     VNUM     Product\r\n"))
	for i = 0; (shop.Producing[i]) != int(-1); i++ {
		write_to_output(d, libc.CString("%2d - [@c%5d@n] - @y%s@n\r\n"), i, obj_index[shop.Producing[i]].Vnum, obj_proto[shop.Producing[i]].Short_description)
	}
	write_to_output(d, libc.CString("\r\n@gA@n) Add a new product.\r\n@gD@n) Delete a product.\r\n@gQ@n) Quit\r\nEnter choice : "))
	d.Olc.Mode = SEDIT_PRODUCTS_MENU
}
func sedit_compact_rooms_menu(d *descriptor_data) {
	var (
		shop  *shop_data
		i     int
		count int = 0
	)
	shop = d.Olc.Shop
	clear_screen(d)
	for i = 0; (shop.In_room[i]) != int(-1); i++ {
		write_to_output(d, libc.CString("%2d - [@c%5d@n]  | %s"), i, shop.In_room[i], func() string {
			if (func() int {
				p := &count
				*p++
				return *p
			}() % 5) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("\r\n@gA@n) Add a new room.\r\n@gD@n) Delete a room.\r\n@gL@n) Long display.\r\n@gQ@n) Quit\r\nEnter choice : "))
	d.Olc.Mode = SEDIT_ROOMS_MENU
}
func sedit_rooms_menu(d *descriptor_data) {
	var (
		shop *shop_data
		i    int
	)
	shop = d.Olc.Shop
	clear_screen(d)
	write_to_output(d, libc.CString("##     VNUM     Room\r\n\r\n"))
	for i = 0; (shop.In_room[i]) != int(-1); i++ {
		if real_room(shop.In_room[i]) != int(-1) {
			write_to_output(d, libc.CString("%2d - [@c%5d@n] - @y%s@n\r\n"), i, shop.In_room[i], world[real_room(shop.In_room[i])].Name)
		} else {
			write_to_output(d, libc.CString("%2d - [@R!Removed Room!@n]\r\n"), i)
		}
	}
	write_to_output(d, libc.CString("\r\n@gA@n) Add a new room.\r\n@gD@n) Delete a room.\r\n@gC@n) Compact Display.\r\n@gQ@n) Quit\r\nEnter choice : "))
	d.Olc.Mode = SEDIT_ROOMS_MENU
}
func sedit_namelist_menu(d *descriptor_data) {
	var (
		shop *shop_data
		i    int
	)
	shop = d.Olc.Shop
	clear_screen(d)
	write_to_output(d, libc.CString("##              Type   Namelist\r\n\r\n"))
	for i = 0; (shop.Type[i]).Type != int(-1); i++ {
		write_to_output(d, libc.CString("%2d - @c%15s@n - @y%s@n\r\n"), i, item_types[(shop.Type[i]).Type], func() *byte {
			if (shop.Type[i]).Keywords != nil {
				return (shop.Type[i]).Keywords
			}
			return libc.CString("<None>")
		}())
	}
	write_to_output(d, libc.CString("\r\n@gA@n) Add a new entry.\r\n@gD@n) Delete an entry.\r\n@gQ@n) Quit\r\nEnter choice : "))
	d.Olc.Mode = SEDIT_NAMELIST_MENU
}
func sedit_shop_flags_menu(d *descriptor_data) {
	var (
		bits  [64936]byte
		i     int
		count int = 0
	)
	clear_screen(d)
	for i = 0; i < NUM_SHOP_FLAGS; i++ {
		write_to_output(d, libc.CString("@g%2d@n) %-20.20s   %s"), i+1, shop_bits[i], func() string {
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
	// todo: fix this
	//sprintbit(d.Olc.Shop.Bitvector, ([]*byte)(shop_bits), &bits[0], uint64(64936))
	write_to_output(d, libc.CString("\r\nCurrent Shop Flags : @c%s@n\r\nEnter choice : "), &bits[0])
	d.Olc.Mode = SEDIT_SHOP_FLAGS
}
func sedit_no_trade_menu(d *descriptor_data) {
	var (
		bits  [64936]byte
		i     int
		count int = 0
	)
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
	sprintbitarray(d.Olc.Shop.With_who[:], trade_letters[:], int(64936), &bits[0])
	write_to_output(d, libc.CString("\r\nCurrently won't trade with: @c%s@n\r\nEnter choice : "), &bits[0])
	d.Olc.Mode = SEDIT_NOTRADE
}
func sedit_types_menu(d *descriptor_data) {
	var shop *shop_data
	_ = shop
	var i int
	var count int = 0
	shop = d.Olc.Shop
	clear_screen(d)
	for i = 0; i < NUM_ITEM_TYPES; i++ {
		write_to_output(d, libc.CString("@g%2d@n) @c%-20s@n  %s"), i, item_types[i], func() string {
			if (func() int {
				p := &count
				*p++
				return *p
			}() % 3) == 0 {
				return "\r\n"
			}
			return ""
		}())
	}
	write_to_output(d, libc.CString("@nEnter choice : "))
	d.Olc.Mode = SEDIT_TYPE_MENU
}
func sedit_disp_menu(d *descriptor_data) {
	var (
		buf1 [64936]byte
		buf2 [64936]byte
		shop *shop_data
	)
	shop = d.Olc.Shop
	clear_screen(d)
	sprintbitarray(shop.With_who[:], trade_letters[:], int(64936), &buf1[0])
	// todo: fix this
	//sprintbit(shop.Bitvector, ([]*byte)(shop_bits), &buf2[0], uint64(64936))
	write_to_output(d, libc.CString("-- Shop Number : [@c%d@n]\r\n@g0@n) Keeper      : [@c%d@n] @y%s\r\n@g1@n) Open 1      : @c%4d@n          @g2@n) Close 1     : @c%4d\r\n@g3@n) Open 2      : @c%4d@n          @g4@n) Close 2     : @c%4d\r\n@g5@n) Sell rate   : @c%1.2f@n          @g6@n) Buy rate    : @c%1.2f\r\n@g7@n) Keeper no item : @y%s\r\n@g8@n) Player no item : @y%s\r\n@g9@n) Keeper no cash : @y%s\r\n@gA@n) Player no cash : @y%s\r\n@gB@n) Keeper no buy  : @y%s\r\n@gC@n) Buy success    : @y%s\r\n@gD@n) Sell success   : @y%s\r\n@gE@n) No Trade With  : @c%s\r\n@gF@n) Shop flags     : @c%s\r\n@gR@n) Rooms Menu\r\n@gP@n) Products Menu\r\n@gT@n) Accept Types Menu\r\n@gW@n) Copy Shop\r\n@gQ@n) Quit\r\nEnter Choice : "), d.Olc.Number, func() int {
		if shop.Keeper == int(-1) {
			return -1
		}
		return mob_index[shop.Keeper].Vnum
	}(), func() string {
		if shop.Keeper == int(-1) {
			return "None"
		}
		return libc.GoString(mob_proto[shop.Keeper].Short_descr)
	}(), shop.Open1, shop.Close1, shop.Open2, shop.Close2, shop.Profit_buy, shop.Profit_sell, shop.No_such_item1, shop.No_such_item2, shop.Missing_cash1, shop.Missing_cash2, shop.Do_not_buy, shop.Message_buy, shop.Message_sell, &buf1[0], &buf2[0])
	d.Olc.Mode = SEDIT_MAIN_MENU
}
func sedit_parse(d *descriptor_data, arg *byte) {
	var i int
	if d.Olc.Mode > SEDIT_NUMERICAL_RESPONSE {
		if !unicode.IsDigit(rune(*arg)) && (*arg == '-' && !unicode.IsDigit(rune(*(*byte)(unsafe.Add(unsafe.Pointer(arg), 1))))) {
			write_to_output(d, libc.CString("Field must be numerical, try again : "))
			return
		}
	}
	switch d.Olc.Mode {
	case SEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			sedit_save_internally(d)
			mudlog(CMP, int(MAX(ADMLVL_BUILDER, int64(d.Character.Player_specials.Invis_level))), 1, libc.CString("OLC: %s edits shop %d"), GET_NAME(d.Character), d.Olc.Number)
			if config_info.Operation.Auto_save_olc != 0 {
				sedit_save_to_disk(real_zone_by_thing(d.Olc.Number))
				write_to_output(d, libc.CString("Shop saved to disk.\r\n"))
			} else {
				write_to_output(d, libc.CString("Shop saved to memory.\r\n"))
			}
			cleanup_olc(d, CLEANUP_STRUCTS)
			return
		case 'n':
			fallthrough
		case 'N':
			cleanup_olc(d, CLEANUP_ALL)
			return
		default:
			write_to_output(d, libc.CString("Invalid choice!\r\nDo you wish to save your changes? : "))
			return
		}
	case SEDIT_MAIN_MENU:
		i = 0
		switch *arg {
		case 'q':
			fallthrough
		case 'Q':
			if d.Olc.Value != 0 {
				write_to_output(d, libc.CString("Do you wish to save your changes? : "))
				d.Olc.Mode = SEDIT_CONFIRM_SAVESTRING
			} else {
				cleanup_olc(d, CLEANUP_ALL)
			}
			return
		case '0':
			d.Olc.Mode = SEDIT_KEEPER
			write_to_output(d, libc.CString("Enter vnum number of shop keeper : "))
			return
		case '1':
			d.Olc.Mode = SEDIT_OPEN1
			i++
		case '2':
			d.Olc.Mode = SEDIT_CLOSE1
			i++
		case '3':
			d.Olc.Mode = SEDIT_OPEN2
			i++
		case '4':
			d.Olc.Mode = SEDIT_CLOSE2
			i++
		case '5':
			d.Olc.Mode = SEDIT_BUY_PROFIT
			i++
		case '6':
			d.Olc.Mode = SEDIT_SELL_PROFIT
			i++
		case '7':
			d.Olc.Mode = SEDIT_NOITEM1
			i--
		case '8':
			d.Olc.Mode = SEDIT_NOITEM2
			i--
		case '9':
			d.Olc.Mode = SEDIT_NOCASH1
			i--
		case 'a':
			fallthrough
		case 'A':
			d.Olc.Mode = SEDIT_NOCASH2
			i--
		case 'b':
			fallthrough
		case 'B':
			d.Olc.Mode = SEDIT_NOBUY
			i--
		case 'c':
			fallthrough
		case 'C':
			d.Olc.Mode = SEDIT_BUY
			i--
		case 'd':
			fallthrough
		case 'D':
			d.Olc.Mode = SEDIT_SELL
			i--
		case 'e':
			fallthrough
		case 'E':
			sedit_no_trade_menu(d)
			return
		case 'f':
			fallthrough
		case 'F':
			sedit_shop_flags_menu(d)
			return
		case 'r':
			fallthrough
		case 'R':
			sedit_rooms_menu(d)
			return
		case 'p':
			fallthrough
		case 'P':
			sedit_products_menu(d)
			return
		case 't':
			fallthrough
		case 'T':
			sedit_namelist_menu(d)
			return
		case 'w':
			fallthrough
		case 'W':
			write_to_output(d, libc.CString("Copy what shop? "))
			d.Olc.Mode = SEDIT_COPY
			return
		default:
			sedit_disp_menu(d)
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
	case SEDIT_NAMELIST_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			sedit_types_menu(d)
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("\r\nDelete which entry? : "))
			d.Olc.Mode = SEDIT_DELETE_TYPE
			return
		case 'q':
			fallthrough
		case 'Q':
		}
	case SEDIT_PRODUCTS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("\r\nEnter new product vnum number : "))
			d.Olc.Mode = SEDIT_NEW_PRODUCT
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("\r\nDelete which product? : "))
			d.Olc.Mode = SEDIT_DELETE_PRODUCT
			return
		case 'q':
			fallthrough
		case 'Q':
		}
	case SEDIT_ROOMS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("\r\nEnter new room vnum number : "))
			d.Olc.Mode = SEDIT_NEW_ROOM
			return
		case 'c':
			fallthrough
		case 'C':
			sedit_compact_rooms_menu(d)
			return
		case 'l':
			fallthrough
		case 'L':
			sedit_rooms_menu(d)
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("\r\nDelete which room? : "))
			d.Olc.Mode = SEDIT_DELETE_ROOM
			return
		case 'q':
			fallthrough
		case 'Q':
		}
	case SEDIT_NOITEM1:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.No_such_item1, arg)
		}
	case SEDIT_NOITEM2:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.No_such_item2, arg)
		}
	case SEDIT_NOCASH1:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.Missing_cash1, arg)
		}
	case SEDIT_NOCASH2:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.Missing_cash2, arg)
		}
	case SEDIT_NOBUY:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.Do_not_buy, arg)
		}
	case SEDIT_BUY:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.Message_buy, arg)
		}
	case SEDIT_SELL:
		if genolc_checkstring(d, arg) {
			modify_string(&d.Olc.Shop.Message_sell, arg)
		}
	case SEDIT_NAMELIST:
		if genolc_checkstring(d, arg) {
			var new_entry shop_buy_data
			new_entry.Type = d.Olc.Value
			new_entry.Keywords = libc.StrDup(arg)
			add_to_type_list((**shop_buy_data)(unsafe.Pointer(&d.Olc.Shop.Type[0])), &new_entry)
		}
		sedit_namelist_menu(d)
		return
	case SEDIT_KEEPER:
		i = libc.Atoi(libc.GoString(arg))
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) != -1 {
			if (func() int {
				i = real_mobile(i)
				return i
			}()) == int(-1) {
				write_to_output(d, libc.CString("That mobile does not exist, try again : "))
				return
			}
		}
		d.Olc.Shop.Keeper = i
		if i == -1 {
			break
		}
		d.Olc.Shop.Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func() SpecialFunc {
				if libc.FuncAddr(mob_index[i].Func) != libc.FuncAddr(shop_keeper) {
					return mob_index[i].Func
				}
				return nil
			}()(ch, me, cmd, argument)
		}
		mob_index[i].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return shop_keeper(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	case SEDIT_OPEN1:
		d.Olc.Shop.Open1 = int(MIN(28, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
	case SEDIT_OPEN2:
		d.Olc.Shop.Open2 = int(MIN(28, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
	case SEDIT_CLOSE1:
		d.Olc.Shop.Close1 = int(MIN(28, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
	case SEDIT_CLOSE2:
		d.Olc.Shop.Close2 = int(MIN(28, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
	case SEDIT_BUY_PROFIT:
		stdio.Sscanf(arg, "%f", &d.Olc.Shop.Profit_buy)
	case SEDIT_SELL_PROFIT:
		stdio.Sscanf(arg, "%f", &d.Olc.Shop.Profit_sell)
	case SEDIT_TYPE_MENU:
		d.Olc.Value = int(MIN(int64(int(NUM_ITEM_TYPES-1)), MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
		write_to_output(d, libc.CString("Enter namelist (return for none) :-\r\n] "))
		d.Olc.Mode = SEDIT_NAMELIST
		return
	case SEDIT_DELETE_TYPE:
		remove_from_type_list((**shop_buy_data)(unsafe.Pointer(&d.Olc.Shop.Type[0])), libc.Atoi(libc.GoString(arg)))
		sedit_namelist_menu(d)
		return
	case SEDIT_NEW_PRODUCT:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) != -1 {
			if (func() int {
				i = real_object(i)
				return i
			}()) == int(-1) {
				write_to_output(d, libc.CString("That object does not exist, try again : "))
				return
			}
		}
		if i > 0 {
			add_to_int_list((**int)(unsafe.Pointer(&d.Olc.Shop.Producing[0])), i)
		}
		sedit_products_menu(d)
		return
	case SEDIT_DELETE_PRODUCT:
		remove_from_int_list((**int)(unsafe.Pointer(&d.Olc.Shop.Producing[0])), libc.Atoi(libc.GoString(arg)))
		sedit_products_menu(d)
		return
	case SEDIT_NEW_ROOM:
		if (func() int {
			i = libc.Atoi(libc.GoString(arg))
			return i
		}()) != -1 {
			if (func() int {
				i = real_room(i)
				return i
			}()) == int(-1) {
				write_to_output(d, libc.CString("That room does not exist, try again : "))
				return
			}
		}
		if i >= 0 {
			add_to_int_list((**int)(unsafe.Pointer(&d.Olc.Shop.In_room[0])), libc.Atoi(libc.GoString(arg)))
		}
		sedit_rooms_menu(d)
		return
	case SEDIT_DELETE_ROOM:
		remove_from_int_list((**int)(unsafe.Pointer(&d.Olc.Shop.In_room[0])), libc.Atoi(libc.GoString(arg)))
		sedit_rooms_menu(d)
		return
	case SEDIT_SHOP_FLAGS:
		if (func() int {
			i = int(MIN(NUM_SHOP_FLAGS, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
			return i
		}()) > 0 {
			d.Olc.Shop.Bitvector ^= uint32(int32(1 << (i - 1)))
			sedit_shop_flags_menu(d)
			return
		}
	case SEDIT_NOTRADE:
		if (func() int {
			i = int(MIN(NUM_TRADERS, MAX(int64(libc.Atoi(libc.GoString(arg))), 0)))
			return i
		}()) > 0 {
			TOGGLE_BIT_AR(d.Olc.Shop.With_who[:], uint32(int32(i-1)))
			sedit_no_trade_menu(d)
			return
		}
	case SEDIT_COPY:
		if (func() int {
			i = real_room(libc.Atoi(libc.GoString(arg)))
			return i
		}()) != int(-1) {
			sedit_setup_existing(d, i)
		} else {
			write_to_output(d, libc.CString("That shop does not exist.\r\n"))
		}
	default:
		cleanup_olc(d, CLEANUP_ALL)
		mudlog(BRF, ADMLVL_BUILDER, 1, libc.CString("SYSERR: OLC: sedit_parse(): Reached default case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
	d.Olc.Value = 1
	sedit_disp_menu(d)
}
