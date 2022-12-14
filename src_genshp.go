package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func copy_shop(tshop *shop_data, fshop *shop_data, free_old_strings int) {
	var i int
	tshop.Vnum = fshop.Vnum
	tshop.Keeper = fshop.Keeper
	tshop.Open1 = fshop.Open1
	tshop.Close1 = fshop.Close1
	tshop.Open2 = fshop.Open2
	tshop.Close2 = fshop.Close2
	tshop.BankAccount = fshop.BankAccount
	tshop.Temper1 = fshop.Temper1
	tshop.Bitvector = fshop.Bitvector
	for i = 0; i < SW_ARRAY_MAX; i++ {
		tshop.With_who[i] = fshop.With_who[i]
	}
	tshop.Lastsort = fshop.Lastsort
	tshop.Profit_buy = fshop.Profit_buy
	tshop.Profit_sell = fshop.Profit_sell
	tshop.Func = fshop.Func
	copy_list((**int)(unsafe.Pointer(&tshop.In_room[0])), &fshop.In_room[0])
	copy_list((**int)(unsafe.Pointer(&tshop.Producing[0])), &fshop.Producing[0])
	copy_type_list((**shop_buy_data)(unsafe.Pointer(&tshop.Type[0])), &fshop.Type[0])
	if free_old_strings != 0 {
		free_shop_strings(tshop)
	}
	tshop.No_such_item1 = str_udup(fshop.No_such_item1)
	tshop.No_such_item2 = str_udup(fshop.No_such_item2)
	tshop.Missing_cash1 = str_udup(fshop.Missing_cash1)
	tshop.Missing_cash2 = str_udup(fshop.Missing_cash2)
	tshop.Do_not_buy = str_udup(fshop.Do_not_buy)
	tshop.Message_buy = str_udup(fshop.Message_buy)
	tshop.Message_sell = str_udup(fshop.Message_sell)
}
func copy_list(tlist **int, flist *int) {
	var (
		num_items int
		i         int
	)
	if *tlist != nil {
		libc.Free(unsafe.Pointer(*tlist))
	}
	for i = 0; *(*int)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(int(0))*uintptr(i))) != int(-1); i++ {
	}
	num_items = i + 1
	*tlist = &make([]int, num_items)[0]
	for i = 0; i < num_items; i++ {
		*(*int)(unsafe.Add(unsafe.Pointer(*tlist), unsafe.Sizeof(int(0))*uintptr(i))) = *(*int)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(int(0))*uintptr(i)))
	}
}
func copy_type_list(tlist **shop_buy_data, flist *shop_buy_data) {
	var (
		num_items int
		i         int
	)
	if *tlist != nil {
		free_type_list(tlist)
	}
	for i = 0; (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type != int(-1); i++ {
	}
	num_items = i + 1
	*tlist = &make([]shop_buy_data, num_items)[0]
	for i = 0; i < num_items; i++ {
		(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*tlist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type = (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type
		if (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Keywords != nil {
			(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*tlist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Keywords = libc.StrDup((*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(flist), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Keywords)
		}
	}
}
func remove_from_type_list(list **shop_buy_data, num int) {
	var (
		i         int
		num_items int
		nlist     *shop_buy_data
	)
	for i = 0; (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type != int(-1); i++ {
	}
	if num < 0 || num >= i {
		return
	}
	num_items = i
	nlist = &make([]shop_buy_data, num_items)[0]
	for i = 0; i < num_items; i++ {
		if i < num {
			*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(shop_buy_data{})*uintptr(i))) = *(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))
		} else {
			*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(shop_buy_data{})*uintptr(i))) = *(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i+1)))
		}
	}
	libc.Free(unsafe.Pointer((*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(num)))).Keywords))
	libc.Free(unsafe.Pointer(*list))
	*list = nlist
}
func add_to_type_list(list **shop_buy_data, newl *shop_buy_data) {
	var (
		i         int
		num_items int
		nlist     *shop_buy_data
	)
	for i = 0; (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type != int(-1); i++ {
	}
	num_items = i
	nlist = &make([]shop_buy_data, num_items+2)[0]
	for i = 0; i < num_items; i++ {
		*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(shop_buy_data{})*uintptr(i))) = *(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))
	}
	*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(shop_buy_data{})*uintptr(num_items))) = *newl
	(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(shop_buy_data{})*uintptr(num_items+1)))).Type = -1
	libc.Free(unsafe.Pointer(*list))
	*list = nlist
}
func add_to_int_list(list **int, newi int) {
	var (
		i         int
		num_items int
		nlist     *int
	)
	for i = 0; *(*int)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(int(0))*uintptr(i))) != int(-1); i++ {
	}
	num_items = i
	nlist = &make([]int, num_items+2)[0]
	for i = 0; i < num_items; i++ {
		*(*int)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(int(0))*uintptr(i))) = *(*int)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(int(0))*uintptr(i)))
	}
	*(*int)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(int(0))*uintptr(num_items))) = newi
	*(*int)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(int(0))*uintptr(num_items+1))) = -1
	libc.Free(unsafe.Pointer(*list))
	*list = nlist
}
func remove_from_int_list(list **int, num int) {
	var (
		i         int
		num_items int
		nlist     *int
	)
	for i = 0; *(*int)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(int(0))*uintptr(i))) != int(-1); i++ {
	}
	if num < 0 || num >= i {
		return
	}
	num_items = i
	nlist = &make([]int, num_items)[0]
	for i = 0; i < num_items; i++ {
		if i < num {
			*(*int)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(int(0))*uintptr(i))) = *(*int)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(int(0))*uintptr(i)))
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(nlist), unsafe.Sizeof(int(0))*uintptr(i))) = *(*int)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(int(0))*uintptr(i+1)))
		}
	}
	libc.Free(unsafe.Pointer(*list))
	*list = nlist
}
func free_shop_strings(shop *shop_data) {
	if shop.No_such_item1 != nil {
		libc.Free(unsafe.Pointer(shop.No_such_item1))
		shop.No_such_item1 = nil
	}
	if shop.No_such_item2 != nil {
		libc.Free(unsafe.Pointer(shop.No_such_item2))
		shop.No_such_item2 = nil
	}
	if shop.Missing_cash1 != nil {
		libc.Free(unsafe.Pointer(shop.Missing_cash1))
		shop.Missing_cash1 = nil
	}
	if shop.Missing_cash2 != nil {
		libc.Free(unsafe.Pointer(shop.Missing_cash2))
		shop.Missing_cash2 = nil
	}
	if shop.Do_not_buy != nil {
		libc.Free(unsafe.Pointer(shop.Do_not_buy))
		shop.Do_not_buy = nil
	}
	if shop.Message_buy != nil {
		libc.Free(unsafe.Pointer(shop.Message_buy))
		shop.Message_buy = nil
	}
	if shop.Message_sell != nil {
		libc.Free(unsafe.Pointer(shop.Message_sell))
		shop.Message_sell = nil
	}
}
func free_type_list(list **shop_buy_data) {
	var i int
	for i = 0; (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Type != int(-1); i++ {
		if (*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Keywords != nil {
			libc.Free(unsafe.Pointer((*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(*list), unsafe.Sizeof(shop_buy_data{})*uintptr(i)))).Keywords))
		}
	}
	libc.Free(unsafe.Pointer(*list))
	*list = nil
}
func free_shop(shop *shop_data) {
	free_shop_strings(shop)
	free_type_list((**shop_buy_data)(unsafe.Pointer(&shop.Type[0])))
	libc.Free(unsafe.Pointer(&shop.In_room[0]))
	libc.Free(unsafe.Pointer(&shop.Producing[0]))
	libc.Free(unsafe.Pointer(shop))
}
func real_shop(vnum int) int {
	var (
		bot      int
		top      int
		mid      int
		last_top int
	)
	if top_shop < 0 {
		return -1
	}
	bot = 0
	top = top_shop
	for {
		last_top = top
		mid = (bot + top) / 2
		if shop_index[mid].Vnum == vnum {
			return mid
		}
		if bot >= top {
			return -1
		}
		if shop_index[mid].Vnum > vnum {
			top = mid
		} else {
			bot = mid + 1
		}
		if top > last_top {
			return -1
		}
	}
}
func modify_string(str **byte, new_s *byte) {
	var (
		buf     [64936]byte
		pointer *byte
	)
	if *new_s != '%' {
		stdio.Snprintf(&buf[0], int(64936), "%%s %s", new_s)
		pointer = &buf[0]
	} else {
		pointer = new_s
	}
	if *str != nil {
		libc.Free(unsafe.Pointer(*str))
	}
	*str = libc.StrDup(pointer)
}
func add_shop(nshp *shop_data) int {
	var (
		rshop int
		found int = 0
		rznum int = real_zone_by_thing(nshp.Vnum)
	)
	if (func() int {
		rshop = real_shop(nshp.Vnum)
		return rshop
	}()) != int(-1) {
		copy_shop(&shop_index[rshop], nshp, 1)
		if rznum != int(-1) {
			add_to_save_list(zone_table[rznum].Number, SL_SHP)
		} else {
			mudlog(BRF, ADMLVL_BUILDER, 1, libc.CString("SYSERR: GenOLC: Cannot determine shop zone."))
		}
		return rshop
	}
	top_shop++
	// todo: fix this
	//shop_index = []shop_data((*shop_data)(libc.Realloc(unsafe.Pointer(&shop_index[0]), top_shop*int(unsafe.Sizeof(shop_data{}))+1)))
	for rshop = top_shop; rshop > 0; rshop-- {
		if nshp.Vnum > shop_index[rshop-1].Vnum {
			found = rshop
			shop_index[rshop].In_room = nil
			shop_index[rshop].Producing = nil
			shop_index[rshop].Type = nil
			copy_shop(&shop_index[rshop], nshp, 0)
			break
		}
		shop_index[rshop] = shop_index[rshop-1]
	}
	if found == 0 {
		shop_index[rshop].In_room = nil
		shop_index[rshop].Producing = nil
		shop_index[rshop].Type = nil
		copy_shop(&shop_index[0], nshp, 0)
	}
	if rznum != int(-1) {
		add_to_save_list(zone_table[rznum].Number, SL_SHP)
	} else {
		mudlog(BRF, ADMLVL_BUILDER, 1, libc.CString("SYSERR: GenOLC: Cannot determine shop zone."))
	}
	return rshop
}
func save_shops(zone_num int) bool {
	var (
		i         int
		j         int
		rshop     int
		shop_file *stdio.File
		fname     [128]byte
		oldname   [128]byte
		shop      *shop_data
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_shops: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return false
	}
	stdio.Snprintf(&fname[0], int(128), "%s%d.new", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		shop_file = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return shop_file
	}()) == nil {
		mudlog(BRF, ADMLVL_GOD, 1, libc.CString("SYSERR: OLC: Cannot open shop file!"))
		return false
	} else if stdio.Fprintf(shop_file, "CircleMUD v3.0 Shop File~\n") < 0 {
		mudlog(BRF, ADMLVL_GOD, 1, libc.CString("SYSERR: OLC: Cannot write to shop file!"))
		shop_file.Close()
		return false
	}
	for i = genolc_zone_bottom(zone_num); i <= zone_table[zone_num].Top; i++ {
		if (func() int {
			rshop = real_shop(i)
			return rshop
		}()) != int(-1) {
			stdio.Fprintf(shop_file, "#%d~\n", i)
			shop = &shop_index[rshop]
			for j = 0; (shop.Producing[j]) != int(-1); j++ {
				stdio.Fprintf(shop_file, "%d\n", obj_index[shop.Producing[j]].Vnum)
			}
			stdio.Fprintf(shop_file, "-1\n")
			stdio.Fprintf(shop_file, "%1.2f\n%1.2f\n", shop.Profit_buy, shop.Profit_sell)
			for j = 0; (shop.Type[j]).Type != int(-1); j++ {
				stdio.Fprintf(shop_file, "%d%s\n", (shop.Type[j]).Type, func() *byte {
					if (shop.Type[j]).Keywords != nil {
						return (shop.Type[j]).Keywords
					}
					return libc.CString("")
				}())
			}
			stdio.Fprintf(shop_file, "-1\n")
			stdio.Fprintf(shop_file, "%s~\n%s~\n%s~\n%s~\n%s~\n%s~\n%s~\n%d\n%d\n%d\n", func() *byte {
				if shop.No_such_item1 != nil {
					return shop.No_such_item1
				}
				return libc.CString("%s Ke?!")
			}(), func() *byte {
				if shop.No_such_item2 != nil {
					return shop.No_such_item2
				}
				return libc.CString("%s Ke?!")
			}(), func() *byte {
				if shop.Do_not_buy != nil {
					return shop.Do_not_buy
				}
				return libc.CString("%s Ke?!")
			}(), func() *byte {
				if shop.Missing_cash1 != nil {
					return shop.Missing_cash1
				}
				return libc.CString("%s Ke?!")
			}(), func() *byte {
				if shop.Missing_cash2 != nil {
					return shop.Missing_cash2
				}
				return libc.CString("%s Ke?!")
			}(), func() *byte {
				if shop.Message_buy != nil {
					return shop.Message_buy
				}
				return libc.CString("%s Ke?! %d?")
			}(), func() *byte {
				if shop.Message_sell != nil {
					return shop.Message_sell
				}
				return libc.CString("%s Ke?! %d?")
			}(), shop.Temper1, shop.Bitvector, func() int {
				if shop.Keeper == int(-1) {
					return -1
				}
				return mob_index[shop.Keeper].Vnum
			}())
			for j = 0; j < SW_ARRAY_MAX; j++ {
				stdio.Fprintf(shop_file, "%s%d", func() string {
					if j != 0 {
						return " "
					}
					return ""
				}(), shop.With_who[j])
			}
			stdio.Fprintf(shop_file, "\n")
			for j = 0; (shop.In_room[j]) != int(-1); j++ {
				stdio.Fprintf(shop_file, "%d\n", shop.In_room[j])
			}
			stdio.Fprintf(shop_file, "-1\n")
			stdio.Fprintf(shop_file, "%d\n%d\n%d\n%d\n", shop.Open1, shop.Close1, shop.Open2, shop.Close2)
		}
	}
	stdio.Fprintf(shop_file, "$~\n")
	shop_file.Close()
	stdio.Snprintf(&oldname[0], int(128), "%s%d.shp", LIB_WORLD, zone_table[zone_num].Number)
	stdio.Remove(libc.GoString(&oldname[0]))
	stdio.Rename(libc.GoString(&fname[0]), libc.GoString(&oldname[0]))
	if in_save_list(zone_table[zone_num].Number, SL_SHP) {
		remove_from_save_list(zone_table[zone_num].Number, SL_SHP)
		create_world_index(zone_table[zone_num].Number, libc.CString("shp"))
		basic_mud_log(libc.CString("GenOLC: save_shops: Saving shops '%s'"), &oldname[0])
	}
	return true
}
