package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type olc_scmd_info_t struct {
	Text     *byte
	Con_type int
}

var olc_scmd_info [12]olc_scmd_info_t = [12]olc_scmd_info_t{{Text: libc.CString("room"), Con_type: CON_REDIT}, {Text: libc.CString("object"), Con_type: CON_OEDIT}, {Text: libc.CString("zone"), Con_type: CON_ZEDIT}, {Text: libc.CString("mobile"), Con_type: CON_MEDIT}, {Text: libc.CString("shop"), Con_type: CON_SEDIT}, {Text: libc.CString("config"), Con_type: CON_CEDIT}, {Text: libc.CString("trigger"), Con_type: CON_TRIGEDIT}, {Text: libc.CString("action"), Con_type: CON_AEDIT}, {Text: libc.CString("guild"), Con_type: CON_GEDIT}, {Text: libc.CString("help"), Con_type: CON_HEDIT}, {Text: libc.CString("house"), Con_type: CON_HSEDIT}, {Text: libc.CString("\n"), Con_type: -1}}

func clear_screen(d *descriptor_data) {
	if PRF_FLAGGED(d.Character, PRF_CLS) {
		write_to_output(d, libc.CString("\x1b[H\x1b[J"))
	}
}
func do_oasis(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) || ch.Desc == nil {
		return
	}
	if ch.Desc.Connected != CON_PLAYING {
		return
	}
	switch subcmd {
	case SCMD_OASIS_CEDIT:
		do_oasis_cedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_ZEDIT:
		do_oasis_zedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_REDIT:
		do_oasis_redit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_OEDIT:
		do_oasis_oedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_MEDIT:
		do_oasis_medit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_SEDIT:
		do_oasis_sedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_AEDIT:
		do_oasis_aedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_HEDIT:
		do_oasis_hedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_HSEDIT:
		do_oasis_hsedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_RLIST:
		fallthrough
	case SCMD_OASIS_MLIST:
		fallthrough
	case SCMD_OASIS_OLIST:
		fallthrough
	case SCMD_OASIS_SLIST:
		fallthrough
	case SCMD_OASIS_ZLIST:
		fallthrough
	case SCMD_OASIS_TLIST:
		fallthrough
	case SCMD_OASIS_GLIST:
		do_oasis_list(ch, argument, cmd, subcmd)
	case SCMD_OASIS_TRIGEDIT:
		do_oasis_trigedit(ch, argument, cmd, subcmd)
	case SCMD_OASIS_LINKS:
		do_oasis_links(ch, argument, cmd, subcmd)
	case SCMD_OASIS_GEDIT:
		do_oasis_gedit(ch, argument, cmd, subcmd)
	default:
		basic_mud_log(libc.CString("SYSERR: (OLC) Invalid subcmd passed to do_oasis, subcmd - (%d)"), subcmd)
		return
	}
	return
}
func cleanup_olc(d *descriptor_data, cleanup_type int8) {
	if d.Olc == nil {
		return
	}
	if d.Olc.Room != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_proto_script(unsafe.Pointer(d.Olc.Room), WLD_TRIGGER)
			free_room(d.Olc.Room)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Room))
		case CLEANUP_CONFIG:
			free_config(d.Olc.Config)
		default:
			basic_mud_log(libc.CString("SYSERR: cleanup_olc: Unknown type!"))
		}
	}
	if d.Olc.Obj != nil {
		free_object_strings(d.Olc.Obj)
		libc.Free(unsafe.Pointer(d.Olc.Obj))
	}
	if d.Olc.Mob != nil {
		free_mobile(d.Olc.Mob)
	}
	if d.Olc.Zone != nil {
		if d.Olc.Zone.Builders != nil {
			libc.Free(unsafe.Pointer(d.Olc.Zone.Builders))
		}
		if d.Olc.Zone.Name != nil {
			libc.Free(unsafe.Pointer(d.Olc.Zone.Name))
		}
		if d.Olc.Zone.Cmd != nil {
			libc.Free(unsafe.Pointer(&d.Olc.Zone.Cmd[0]))
		}
		libc.Free(unsafe.Pointer(d.Olc.Zone))
	}
	if d.Olc.Shop != nil {
		free_shop(d.Olc.Shop)
	}
	if d.Olc.Guild != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_guild(d.Olc.Guild)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Guild))
		default:
		}
	}
	if d.Olc.House != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_house(d.Olc.House)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.House))
		default:
		}
	}
	if d.Olc.Action != nil {
		switch cleanup_type {
		case CLEANUP_ALL:
			free_action(d.Olc.Action)
		case CLEANUP_STRUCTS:
			libc.Free(unsafe.Pointer(d.Olc.Action))
		default:
		}
	}
	if d.Olc.Storage != nil {
		libc.Free(unsafe.Pointer(d.Olc.Storage))
		d.Olc.Storage = nil
	}
	if d.Olc.Trig != nil {
		free_trigger(d.Olc.Trig)
		d.Olc.Trig = nil
	}
	if d.Character != nil {
		REMOVE_BIT_AR(d.Character.Act[:], PLR_WRITING)
		act(libc.CString("$n stops using OLC."), 1, d.Character, nil, nil, TO_ROOM)
		if int(cleanup_type) == CLEANUP_CONFIG {
			mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("OLC: %s stops editing the game configuration"), GET_NAME(d.Character))
		} else if d.Connected == CON_TEDIT {
			mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("OLC: %s stops editing text files."), GET_NAME(d.Character))
		} else if d.Connected == CON_HEDIT {
			mudlog(CMP, ADMLVL_IMMORT, 1, libc.CString("OLC: %s stops editing help files."), GET_NAME(d.Character))
		} else {
			mudlog(BRF, ADMLVL_IMMORT, 1, libc.CString("OLC: %s stops editing zone %d allowed zone %d"), GET_NAME(d.Character), zone_table[d.Olc.Zone_num].Number, d.Character.Player_specials.Olc_zone)
		}
		d.Connected = CON_PLAYING
	}
	libc.Free(unsafe.Pointer(d.Olc))
	d.Olc = nil
}
func split_argument(argument *byte, tag *byte) {
	var (
		tmp  *byte = argument
		ttag *byte = tag
		wrt  *byte = argument
		i    int
	)
	for i = 0; *tmp != 0; func() int {
		tmp = (*byte)(unsafe.Add(unsafe.Pointer(tmp), 1))
		return func() int {
			p := &i
			x := *p
			*p++
			return x
		}()
	}() {
		if *tmp != ' ' && *tmp != '=' {
			*(func() *byte {
				p := &ttag
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *tmp
		} else if *tmp == '=' {
			break
		}
	}
	*ttag = '\x00'
	for *tmp == '=' || *tmp == ' ' {
		tmp = (*byte)(unsafe.Add(unsafe.Pointer(tmp), 1))
	}
	for *tmp != 0 {
		*(func() *byte {
			p := &wrt
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(func() *byte {
			p := &tmp
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	*wrt = '\x00'
}
func free_config(data *config_data) {
	free_strings(unsafe.Pointer(data), OASIS_CFG)
	libc.Free(unsafe.Pointer(data))
}
func can_edit_zone(ch *char_data, rnum int) bool {
	if ch.Desc == nil || IS_NPC(ch) || rnum == int(-1) {
		return false
	}
	if rnum == HEDIT_PERMISSION {
		return true
	}
	if ch.Admlevel >= ADMLVL_GRGOD {
		return true
	}
	if is_name(GET_NAME(ch), zone_table[rnum].Builders) {
		return true
	}
	if ch.Player_specials.Olc_zone == int(-1) {
		return false
	}
	if ch.Admlevel < ADMLVL_BUILDER {
		return false
	}
	if real_zone(ch.Player_specials.Olc_zone) == rnum {
		return true
	}
	return false
}
func send_cannot_edit(ch *char_data, zone int) {
	send_to_char(ch, libc.CString("You do not have permission to edit zone %d."), zone)
	if ch.Player_specials.Olc_zone != int(-1) {
		send_to_char(ch, libc.CString("  Try zone %d."), ch.Player_specials.Olc_zone)
	}
	send_to_char(ch, libc.CString("\r\n"))
	mudlog(BRF, ADMLVL_IMPL, 1, libc.CString("OLC: %s tried to edit zone %d allowed zone %d"), GET_NAME(ch), zone, ch.Player_specials.Olc_zone)
}
