package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

var save_list *save_list_data
var save_types [10]struct {
	Save_type int
	Func      func(rnum int) int
	Message   *byte
} = [10]struct {
	Save_type int
	Func      func(rnum int) int
	Message   *byte
}{{Save_type: SL_MOB, Func: save_mobiles, Message: libc.CString("mobile")}, {Save_type: SL_OBJ, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_objects(rnum)))
}, Message: libc.CString("object")}, {Save_type: SL_SHP, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_shops(rnum)))
}, Message: libc.CString("shop")}, {Save_type: SL_WLD, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_rooms(rnum)))
}, Message: libc.CString("room")}, {Save_type: SL_ZON, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_zone(rnum)))
}, Message: libc.CString("zone")}, {Save_type: SL_CFG, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_config(rnum)))
}, Message: libc.CString("config")}, {Save_type: SL_GLD, Func: func(rnum int) int {
	return int(libc.BoolToInt(save_guilds(rnum)))
}, Message: libc.CString("guild")}, {Save_type: int(SL_GLD + 1), Func: nil, Message: libc.CString("social")}, {Save_type: int(SL_GLD + 2), Func: nil, Message: libc.CString("help")}, {Save_type: -1, Func: nil, Message: nil}}

func genolc_checkstring(d *descriptor_data, arg *byte) bool {
	smash_tilde(arg)
	return true
}
func str_udup(txt *byte) *byte {
	return libc.StrDup(func() *byte {
		if txt != nil && *txt != 0 {
			return txt
		}
		return libc.CString("undefined")
	}())
}
func save_all() bool {
	for save_list != nil {
		if save_list.Type < 0 || save_list.Type > SL_GLD {
			switch save_list.Type {
			case int(SL_GLD + 1):
				basic_mud_log(libc.CString("Actions not saved - can not autosave. Use 'aedit save'."))
				save_list = save_list.Next
			case int(SL_GLD + 2):
				basic_mud_log(libc.CString("Help not saved - can not autosave. Use 'hedit save'."))
				save_list = save_list.Next
			default:
				basic_mud_log(libc.CString("SYSERR: GenOLC: Invalid save type %d in save list.\n"), save_list.Type)
			}
		} else if (save_types[save_list.Type].Func)(real_zone(save_list.Zone)) < 0 {
			save_list = save_list.Next
		}
	}
	return true
}
func strip_cr(buffer *byte) {
	var (
		rpos int
		wpos int
	)
	if buffer == nil {
		return
	}
	for func() int {
		rpos = 0
		return func() int {
			wpos = 0
			return wpos
		}()
	}(); *(*byte)(unsafe.Add(unsafe.Pointer(buffer), rpos)) != 0; rpos++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(buffer), wpos)) = *(*byte)(unsafe.Add(unsafe.Pointer(buffer), rpos))
		wpos += int(libc.BoolToInt(*(*byte)(unsafe.Add(unsafe.Pointer(buffer), rpos)) != '\r'))
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(buffer), wpos)) = '\x00'
}
func copy_ex_descriptions(to **extra_descr_data, from *extra_descr_data) {
	var wpos *extra_descr_data
	*to = new(extra_descr_data)
	wpos = *to
	for ; from != nil; func() *extra_descr_data {
		from = from.Next
		return func() *extra_descr_data {
			wpos = wpos.Next
			return wpos
		}()
	}() {
		wpos.Keyword = str_udup(from.Keyword)
		wpos.Description = str_udup(from.Description)
		if from.Next != nil {
			wpos.Next = new(extra_descr_data)
		}
	}
}
func free_ex_descriptions(head *extra_descr_data) {
	var (
		thised   *extra_descr_data
		next_one *extra_descr_data
	)
	if head == nil {
		basic_mud_log(libc.CString("free_ex_descriptions: NULL pointer or NULL data."))
		return
	}
	for thised = head; thised != nil; thised = next_one {
		next_one = thised.Next
		if thised.Keyword != nil {
			libc.Free(unsafe.Pointer(thised.Keyword))
		}
		if thised.Description != nil {
			libc.Free(unsafe.Pointer(thised.Description))
		}
		libc.Free(unsafe.Pointer(thised))
	}
}
func remove_from_save_list(zone int, type_ int) bool {
	var (
		ritem *save_list_data
		temp  *save_list_data
	)
	for ritem = save_list; ritem != nil; ritem = ritem.Next {
		if ritem.Zone == zone && ritem.Type == type_ {
			break
		}
	}
	if ritem == nil {
		basic_mud_log(libc.CString("SYSERR: remove_from_save_list: Saved item not found. (%d/%d)"), zone, type_)
		return false
	}
	if ritem == save_list {
		save_list = ritem.Next
	} else {
		temp = save_list
		for temp != nil && temp.Next != ritem {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = ritem.Next
		}
	}
	libc.Free(unsafe.Pointer(ritem))
	return true
}
func add_to_save_list(zone int, type_ int) bool {
	var (
		nitem *save_list_data
		rznum int
	)
	if type_ == SL_CFG {
		return false
	}
	rznum = real_zone(zone)
	if rznum == int(-1) || rznum > top_of_zone_table {
		if zone != AEDIT_PERMISSION && zone != HEDIT_PERMISSION {
			basic_mud_log(libc.CString("SYSERR: add_to_save_list: Invalid zone number passed. (%d => %d, 0-%d)"), zone, rznum, top_of_zone_table)
			return false
		}
	}
	for nitem = save_list; nitem != nil; nitem = nitem.Next {
		if nitem.Zone == zone && nitem.Type == type_ {
			return false
		}
	}
	nitem = new(save_list_data)
	nitem.Zone = zone
	nitem.Type = type_
	nitem.Next = save_list
	save_list = nitem
	return true
}
func in_save_list(zone int, type_ int) bool {
	var nitem *save_list_data
	for nitem = save_list; nitem != nil; nitem = nitem.Next {
		if nitem.Zone == zone && nitem.Type == type_ {
			return true
		}
	}
	return false
}
func free_save_list() {
	var (
		sld      *save_list_data
		next_sld *save_list_data
	)
	for sld = save_list; sld != nil; sld = next_sld {
		next_sld = sld.Next
		libc.Free(unsafe.Pointer(sld))
	}
}
func do_show_save_list(ch *char_data, argument *byte, cmd int, subcmd int) {
	if save_list == nil {
		send_to_char(ch, libc.CString("All world files are up to date.\r\n"))
	} else {
		var item *save_list_data
		send_to_char(ch, libc.CString("The following files need saving:\r\n"))
		for item = save_list; item != nil; item = item.Next {
			if item.Type != SL_CFG {
				send_to_char(ch, libc.CString(" - %s data for zone %d.\r\n"), save_types[item.Type].Message, item.Zone)
			} else {
				send_to_char(ch, libc.CString(" - Game configuration data.\r\n"))
			}
		}
	}
}
func genolc_zonep_bottom(zone *zone_data) int {
	return zone.Bot
}
func genolc_zone_bottom(rznum int) int {
	return zone_table[rznum].Bot
}
func sprintascii(out *byte, bits uint32) int {
	var (
		i     int
		j     int   = 0
		flags *byte = libc.CString("abcdefghijklmnopqrstuvwxyzABCDEF")
	)
	for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(flags), i)) != '\x00'; i++ {
		if int(bits)&(1<<i) != 0 {
			*(*byte)(unsafe.Add(unsafe.Pointer(out), func() int {
				p := &j
				x := *p
				*p++
				return x
			}())) = *(*byte)(unsafe.Add(unsafe.Pointer(flags), i))
		}
	}
	if j == 0 {
		*(*byte)(unsafe.Add(unsafe.Pointer(out), func() int {
			p := &j
			x := *p
			*p++
			return x
		}())) = '0'
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(out), func() int {
		p := &j
		x := *p
		*p++
		return x
	}())) = '\x00'
	return j
}
