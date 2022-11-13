package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func do_oasis_cedit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		d    *descriptor_data
		buf1 [64936]byte
	)
	one_argument(argument, &buf1[0])
	if ch.Admlevel < 5 {
		send_to_char(ch, libc.CString("You can't modify the game configuration.\r\n"))
		return
	}
	d = ch.Desc
	if buf1[0] == 0 {
		d.Olc = new(oasis_olc_data)
		d.Olc.Zone = nil
		cedit_setup(d)
		d.Connected = CON_CEDIT
		act(libc.CString("$n starts using OLC."), TRUE, d.Character, nil, nil, TO_ROOM)
		ch.Act[int(PLR_WRITING/32)] |= bitvector_t(int32(1 << (int(PLR_WRITING % 32))))
		mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("OLC: %s starts editing the game configuration."), GET_NAME(ch))
		return
	} else if libc.StrCaseCmp(libc.CString("save"), &buf1[0]) != 0 {
		send_to_char(ch, libc.CString("Yikes!  Stop that, someone will get hurt!\r\n"))
		return
	}
	send_to_char(ch, libc.CString("Saving the game configuration.\r\n"))
	mudlog(CMP, MAX(ADMLVL_BUILDER, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s saves the game configuration."), GET_NAME(ch))
	cedit_save_to_disk()
}
func cedit_setup(d *descriptor_data) {
	d.Olc.Config = new(config_data)
	d.Olc.Config.Play.Pk_allowed = config_info.Play.Pk_allowed
	d.Olc.Config.Play.Pt_allowed = config_info.Play.Pt_allowed
	d.Olc.Config.Play.Level_can_shout = config_info.Play.Level_can_shout
	d.Olc.Config.Play.Holler_move_cost = config_info.Play.Holler_move_cost
	d.Olc.Config.Play.Tunnel_size = config_info.Play.Tunnel_size
	d.Olc.Config.Play.Max_exp_gain = config_info.Play.Max_exp_gain
	d.Olc.Config.Play.Max_exp_loss = config_info.Play.Max_exp_loss
	d.Olc.Config.Play.Max_npc_corpse_time = config_info.Play.Max_npc_corpse_time
	d.Olc.Config.Play.Max_pc_corpse_time = config_info.Play.Max_pc_corpse_time
	d.Olc.Config.Play.Idle_void = config_info.Play.Idle_void
	d.Olc.Config.Play.Idle_rent_time = config_info.Play.Idle_rent_time
	d.Olc.Config.Play.Idle_max_level = config_info.Play.Idle_max_level
	d.Olc.Config.Play.Dts_are_dumps = config_info.Play.Dts_are_dumps
	d.Olc.Config.Play.Load_into_inventory = config_info.Play.Load_into_inventory
	d.Olc.Config.Play.Track_through_doors = config_info.Play.Track_through_doors
	d.Olc.Config.Play.Level_cap = config_info.Play.Level_cap
	d.Olc.Config.Play.Stack_mobs = config_info.Play.Stack_mobs
	d.Olc.Config.Play.Stack_objs = config_info.Play.Stack_objs
	d.Olc.Config.Play.Mob_fighting = config_info.Play.Mob_fighting
	d.Olc.Config.Play.Disp_closed_doors = config_info.Play.Disp_closed_doors
	d.Olc.Config.Play.Reroll_player = config_info.Play.Reroll_player
	d.Olc.Config.Play.Initial_points = config_info.Play.Initial_points
	d.Olc.Config.Play.Enable_compression = config_info.Play.Enable_compression
	d.Olc.Config.Play.Enable_languages = config_info.Play.Enable_languages
	d.Olc.Config.Play.All_items_unique = config_info.Play.All_items_unique
	d.Olc.Config.Play.Exp_multiplier = config_info.Play.Exp_multiplier
	d.Olc.Config.Csd.Free_rent = config_info.Csd.Free_rent
	d.Olc.Config.Csd.Max_obj_save = config_info.Csd.Max_obj_save
	d.Olc.Config.Csd.Min_rent_cost = config_info.Csd.Min_rent_cost
	d.Olc.Config.Csd.Auto_save = config_info.Csd.Auto_save
	d.Olc.Config.Csd.Autosave_time = config_info.Csd.Autosave_time
	d.Olc.Config.Csd.Crash_file_timeout = config_info.Csd.Crash_file_timeout
	d.Olc.Config.Csd.Rent_file_timeout = config_info.Csd.Rent_file_timeout
	d.Olc.Config.Room_nums.Mortal_start_room = config_info.Room_nums.Mortal_start_room
	d.Olc.Config.Room_nums.Immort_start_room = config_info.Room_nums.Immort_start_room
	d.Olc.Config.Room_nums.Frozen_start_room = config_info.Room_nums.Frozen_start_room
	d.Olc.Config.Room_nums.Donation_room_1 = config_info.Room_nums.Donation_room_1
	d.Olc.Config.Room_nums.Donation_room_2 = config_info.Room_nums.Donation_room_2
	d.Olc.Config.Room_nums.Donation_room_3 = config_info.Room_nums.Donation_room_3
	d.Olc.Config.Operation.DFLT_PORT = config_info.Operation.DFLT_PORT
	d.Olc.Config.Operation.Max_playing = config_info.Operation.Max_playing
	d.Olc.Config.Operation.Max_filesize = config_info.Operation.Max_filesize
	d.Olc.Config.Operation.Max_bad_pws = config_info.Operation.Max_bad_pws
	d.Olc.Config.Operation.Siteok_everyone = config_info.Operation.Siteok_everyone
	d.Olc.Config.Operation.Use_new_socials = config_info.Operation.Use_new_socials
	d.Olc.Config.Operation.Auto_save_olc = config_info.Operation.Auto_save_olc
	d.Olc.Config.Operation.Imc_enabled = config_info.Operation.Imc_enabled
	d.Olc.Config.Operation.Nameserver_is_slow = config_info.Operation.Nameserver_is_slow
	d.Olc.Config.Autowiz.Use_autowiz = config_info.Autowiz.Use_autowiz
	d.Olc.Config.Autowiz.Min_wizlist_lev = config_info.Autowiz.Min_wizlist_lev
	d.Olc.Config.Advance.Allow_multiclass = config_info.Advance.Allow_multiclass
	d.Olc.Config.Advance.Allow_prestige = config_info.Advance.Allow_prestige
	d.Olc.Config.Play.OK = str_udup(config_info.Play.OK)
	d.Olc.Config.Play.NOPERSON = str_udup(config_info.Play.NOPERSON)
	d.Olc.Config.Play.NOEFFECT = str_udup(config_info.Play.NOEFFECT)
	if config_info.Operation.DFLT_IP != nil {
		d.Olc.Config.Operation.DFLT_IP = libc.StrDup(config_info.Operation.DFLT_IP)
	} else {
		d.Olc.Config.Operation.DFLT_IP = nil
	}
	if config_info.Operation.DFLT_DIR != nil {
		d.Olc.Config.Operation.DFLT_DIR = libc.StrDup(config_info.Operation.DFLT_DIR)
	} else {
		d.Olc.Config.Operation.DFLT_DIR = nil
	}
	if config_info.Operation.LOGNAME != nil {
		d.Olc.Config.Operation.LOGNAME = libc.StrDup(config_info.Operation.LOGNAME)
	} else {
		d.Olc.Config.Operation.LOGNAME = nil
	}
	if config_info.Operation.MENU != nil {
		d.Olc.Config.Operation.MENU = libc.StrDup(config_info.Operation.MENU)
	} else {
		d.Olc.Config.Operation.MENU = nil
	}
	if config_info.Operation.WELC_MESSG != nil {
		d.Olc.Config.Operation.WELC_MESSG = libc.StrDup(config_info.Operation.WELC_MESSG)
	} else {
		d.Olc.Config.Operation.WELC_MESSG = nil
	}
	if config_info.Operation.START_MESSG != nil {
		d.Olc.Config.Operation.START_MESSG = libc.StrDup(config_info.Operation.START_MESSG)
	} else {
		d.Olc.Config.Operation.START_MESSG = nil
	}
	cedit_disp_menu(d)
	d.Olc.Config.Ticks.Pulse_violence = config_info.Ticks.Pulse_violence
	d.Olc.Config.Ticks.Pulse_mobile = config_info.Ticks.Pulse_mobile
	d.Olc.Config.Ticks.Pulse_zone = config_info.Ticks.Pulse_zone
	d.Olc.Config.Ticks.Pulse_autosave = config_info.Ticks.Pulse_autosave
	d.Olc.Config.Ticks.Pulse_idlepwd = config_info.Ticks.Pulse_idlepwd
	d.Olc.Config.Ticks.Pulse_sanity = config_info.Ticks.Pulse_sanity
	d.Olc.Config.Ticks.Pulse_usage = config_info.Ticks.Pulse_usage
	d.Olc.Config.Ticks.Pulse_timesave = config_info.Ticks.Pulse_timesave
	d.Olc.Config.Ticks.Pulse_current = config_info.Ticks.Pulse_current
	d.Olc.Config.Creation.Method = config_info.Creation.Method
}
func cedit_save_internally(d *descriptor_data) {
	var copyover_needed bool = FALSE != 0
	_ = copyover_needed
	var reassign int = int(libc.BoolToInt(config_info.Play.Dts_are_dumps != d.Olc.Config.Play.Dts_are_dumps))
	config_info.Play.Pk_allowed = d.Olc.Config.Play.Pk_allowed
	config_info.Play.Pt_allowed = d.Olc.Config.Play.Pt_allowed
	config_info.Play.Level_can_shout = d.Olc.Config.Play.Level_can_shout
	config_info.Play.Holler_move_cost = d.Olc.Config.Play.Holler_move_cost
	config_info.Play.Tunnel_size = d.Olc.Config.Play.Tunnel_size
	config_info.Play.Max_exp_gain = d.Olc.Config.Play.Max_exp_gain
	config_info.Play.Max_exp_loss = d.Olc.Config.Play.Max_exp_loss
	config_info.Play.Max_npc_corpse_time = d.Olc.Config.Play.Max_npc_corpse_time
	config_info.Play.Max_pc_corpse_time = d.Olc.Config.Play.Max_pc_corpse_time
	config_info.Play.Idle_void = d.Olc.Config.Play.Idle_void
	config_info.Play.Idle_rent_time = d.Olc.Config.Play.Idle_rent_time
	config_info.Play.Idle_max_level = d.Olc.Config.Play.Idle_max_level
	config_info.Play.Dts_are_dumps = d.Olc.Config.Play.Dts_are_dumps
	config_info.Play.Load_into_inventory = d.Olc.Config.Play.Load_into_inventory
	config_info.Play.Track_through_doors = d.Olc.Config.Play.Track_through_doors
	config_info.Play.Level_cap = d.Olc.Config.Play.Level_cap
	config_info.Play.Stack_mobs = d.Olc.Config.Play.Stack_mobs
	config_info.Play.Stack_objs = d.Olc.Config.Play.Stack_objs
	config_info.Play.Mob_fighting = d.Olc.Config.Play.Mob_fighting
	config_info.Play.Disp_closed_doors = d.Olc.Config.Play.Disp_closed_doors
	config_info.Play.Reroll_player = d.Olc.Config.Play.Reroll_player
	config_info.Play.Initial_points = d.Olc.Config.Play.Initial_points
	config_info.Play.Enable_compression = d.Olc.Config.Play.Enable_compression
	config_info.Play.Enable_languages = d.Olc.Config.Play.Enable_languages
	config_info.Play.All_items_unique = d.Olc.Config.Play.All_items_unique
	config_info.Play.Exp_multiplier = d.Olc.Config.Play.Exp_multiplier
	config_info.Csd.Free_rent = d.Olc.Config.Csd.Free_rent
	config_info.Csd.Max_obj_save = d.Olc.Config.Csd.Max_obj_save
	config_info.Csd.Min_rent_cost = d.Olc.Config.Csd.Min_rent_cost
	config_info.Csd.Auto_save = d.Olc.Config.Csd.Auto_save
	config_info.Csd.Autosave_time = d.Olc.Config.Csd.Autosave_time
	config_info.Csd.Crash_file_timeout = d.Olc.Config.Csd.Crash_file_timeout
	config_info.Csd.Rent_file_timeout = d.Olc.Config.Csd.Rent_file_timeout
	config_info.Room_nums.Mortal_start_room = d.Olc.Config.Room_nums.Mortal_start_room
	config_info.Room_nums.Immort_start_room = d.Olc.Config.Room_nums.Immort_start_room
	config_info.Room_nums.Frozen_start_room = d.Olc.Config.Room_nums.Frozen_start_room
	config_info.Room_nums.Donation_room_1 = d.Olc.Config.Room_nums.Donation_room_1
	config_info.Room_nums.Donation_room_2 = d.Olc.Config.Room_nums.Donation_room_2
	config_info.Room_nums.Donation_room_3 = d.Olc.Config.Room_nums.Donation_room_3
	config_info.Operation.DFLT_PORT = d.Olc.Config.Operation.DFLT_PORT
	config_info.Operation.Max_playing = d.Olc.Config.Operation.Max_playing
	config_info.Operation.Max_filesize = d.Olc.Config.Operation.Max_filesize
	config_info.Operation.Max_bad_pws = d.Olc.Config.Operation.Max_bad_pws
	config_info.Operation.Siteok_everyone = d.Olc.Config.Operation.Siteok_everyone
	config_info.Operation.Use_new_socials = d.Olc.Config.Operation.Use_new_socials
	config_info.Operation.Nameserver_is_slow = d.Olc.Config.Operation.Nameserver_is_slow
	config_info.Operation.Auto_save_olc = d.Olc.Config.Operation.Auto_save_olc
	config_info.Operation.Imc_enabled = d.Olc.Config.Operation.Imc_enabled
	config_info.Autowiz.Use_autowiz = d.Olc.Config.Autowiz.Use_autowiz
	config_info.Autowiz.Min_wizlist_lev = d.Olc.Config.Autowiz.Min_wizlist_lev
	config_info.Advance.Allow_multiclass = d.Olc.Config.Advance.Allow_multiclass
	config_info.Advance.Allow_prestige = d.Olc.Config.Advance.Allow_prestige
	if config_info.Play.OK != nil {
		libc.Free(unsafe.Pointer(config_info.Play.OK))
	}
	config_info.Play.OK = str_udup(d.Olc.Config.Play.OK)
	if config_info.Play.NOPERSON != nil {
		libc.Free(unsafe.Pointer(config_info.Play.NOPERSON))
	}
	config_info.Play.NOPERSON = str_udup(d.Olc.Config.Play.NOPERSON)
	if config_info.Play.NOEFFECT != nil {
		libc.Free(unsafe.Pointer(config_info.Play.NOEFFECT))
	}
	config_info.Play.NOEFFECT = str_udup(d.Olc.Config.Play.NOEFFECT)
	if config_info.Operation.DFLT_IP != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.DFLT_IP))
	}
	if d.Olc.Config.Operation.DFLT_IP != nil {
		config_info.Operation.DFLT_IP = libc.StrDup(d.Olc.Config.Operation.DFLT_IP)
	} else {
		config_info.Operation.DFLT_IP = nil
	}
	if config_info.Operation.DFLT_DIR != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.DFLT_DIR))
	}
	if d.Olc.Config.Operation.DFLT_DIR != nil {
		config_info.Operation.DFLT_DIR = libc.StrDup(d.Olc.Config.Operation.DFLT_DIR)
	} else {
		config_info.Operation.DFLT_DIR = nil
	}
	if config_info.Operation.LOGNAME != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.LOGNAME))
	}
	if d.Olc.Config.Operation.LOGNAME != nil {
		config_info.Operation.LOGNAME = libc.StrDup(d.Olc.Config.Operation.LOGNAME)
	} else {
		config_info.Operation.LOGNAME = nil
	}
	if config_info.Operation.MENU != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.MENU))
	}
	if d.Olc.Config.Operation.MENU != nil {
		config_info.Operation.MENU = libc.StrDup(d.Olc.Config.Operation.MENU)
	} else {
		config_info.Operation.MENU = nil
	}
	if config_info.Operation.WELC_MESSG != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.WELC_MESSG))
	}
	if d.Olc.Config.Operation.WELC_MESSG != nil {
		config_info.Operation.WELC_MESSG = libc.StrDup(d.Olc.Config.Operation.WELC_MESSG)
	} else {
		config_info.Operation.WELC_MESSG = nil
	}
	if config_info.Operation.START_MESSG != nil {
		libc.Free(unsafe.Pointer(config_info.Operation.START_MESSG))
	}
	if d.Olc.Config.Operation.START_MESSG != nil {
		config_info.Operation.START_MESSG = libc.StrDup(d.Olc.Config.Operation.START_MESSG)
	} else {
		config_info.Operation.START_MESSG = nil
	}
	if reassign != 0 {
		reassign_rooms()
	}
	add_to_save_list(zone_vnum(-1), SL_CFG)
	config_info.Ticks.Pulse_violence = d.Olc.Config.Ticks.Pulse_violence
	config_info.Ticks.Pulse_violence = d.Olc.Config.Ticks.Pulse_violence
	config_info.Ticks.Pulse_mobile = d.Olc.Config.Ticks.Pulse_mobile
	config_info.Ticks.Pulse_zone = d.Olc.Config.Ticks.Pulse_zone
	config_info.Ticks.Pulse_autosave = d.Olc.Config.Ticks.Pulse_autosave
	config_info.Ticks.Pulse_idlepwd = d.Olc.Config.Ticks.Pulse_idlepwd
	config_info.Ticks.Pulse_sanity = d.Olc.Config.Ticks.Pulse_sanity
	config_info.Ticks.Pulse_usage = d.Olc.Config.Ticks.Pulse_usage
	config_info.Ticks.Pulse_timesave = d.Olc.Config.Ticks.Pulse_timesave
	config_info.Ticks.Pulse_current = d.Olc.Config.Ticks.Pulse_current
	config_info.Creation.Method = d.Olc.Config.Creation.Method
}
func cedit_save_to_disk() {
	save_config(int64(-1))
}
func save_config(nowhere int64) int {
	var (
		fl  *stdio.File
		buf [64936]byte
	)
	if (func() *stdio.File {
		fl = stdio.FOpen(libc.GoString(config_info.CONFFILE), "w")
		return fl
	}()) == nil {
		perror(libc.CString("SYSERR: save_config"))
		return FALSE
	}
	stdio.Fprintf(fl, "* This file is autogenerated by OasisOLC (CEdit).\n* Please note the following information about this file's format.\n*\n* - If variable is a yes/no or true/false based variable, use 1's and 0's\n*   where YES or TRUE = 1 and NO or FALSE = 0.\n* - Variable names in this file are case-insensitive.  Variable values\n*   are not case-insensitive.\n* -----------------------------------------------------------------------\n* Lines starting with * are comments, and are not parsed.\n* -----------------------------------------------------------------------\n\n* [ Game Play Options ]\n")
	stdio.Fprintf(fl, "* Is player killing allowed on the mud?\npk_allowed = %d\n\n", config_info.Play.Pk_allowed)
	stdio.Fprintf(fl, "* Is player thieving allowed on the mud?\npt_allowed = %d\n\n", config_info.Play.Pt_allowed)
	stdio.Fprintf(fl, "* What is the minimum level a player can shout/gossip/etc?\nlevel_can_shout = %d\n\n", config_info.Play.Level_can_shout)
	stdio.Fprintf(fl, "* How many movement points does shouting cost the player?\nholler_move_cost = %d\n\n", config_info.Play.Holler_move_cost)
	stdio.Fprintf(fl, "* How many players can fit in a tunnel?\ntunnel_size = %d\n\n", config_info.Play.Tunnel_size)
	stdio.Fprintf(fl, "* Maximum experience gainable per kill?\nmax_exp_gain = %d\n\n", config_info.Play.Max_exp_gain)
	stdio.Fprintf(fl, "* Maximum experience loseable per death?\nmax_exp_loss = %d\n\n", config_info.Play.Max_exp_loss)
	stdio.Fprintf(fl, "* Number of tics before NPC corpses decompose.\nmax_npc_corpse_time = %d\n\n", config_info.Play.Max_npc_corpse_time)
	stdio.Fprintf(fl, "* Number of tics before PC corpses decompose.\nmax_pc_corpse_time = %d\n\n", config_info.Play.Max_pc_corpse_time)
	stdio.Fprintf(fl, "* Number of tics before a PC is sent to the void.\nidle_void = %d\n\n", config_info.Play.Idle_void)
	stdio.Fprintf(fl, "* Number of tics before a PC is autorented.\nidle_rent_time = %d\n\n", config_info.Play.Idle_rent_time)
	stdio.Fprintf(fl, "* Admin level and above of players whom are immune to idle penalties.\nidle_max_level = %d\n\n", config_info.Play.Idle_max_level)
	stdio.Fprintf(fl, "* Should the items in death traps be junked automatically?\ndts_are_dumps = %d\n\n", config_info.Play.Dts_are_dumps)
	stdio.Fprintf(fl, "* When an immortal loads an object, should it load into their inventory?\nload_into_inventory = %d\n\n", config_info.Play.Load_into_inventory)
	stdio.Fprintf(fl, "* Should PC's be able to track through hidden or closed doors?\ntrack_through_doors = %d\n\n", config_info.Play.Track_through_doors)
	stdio.Fprintf(fl, "* What is the level that cannot be attained?\nlevel_cap = %d\n\n", config_info.Play.Level_cap)
	stdio.Fprintf(fl, "* Stack mobiles when showing contents of rooms?\nstack_mobs = %d\n\n", config_info.Play.Stack_mobs)
	stdio.Fprintf(fl, "* Stack objects when showing contents of rooms?\nstack_objs = %d\n\n", config_info.Play.Stack_objs)
	stdio.Fprintf(fl, "* Allow aggressive mobs to attack other mobs?\nmob_fighting = %d\n\n", config_info.Play.Mob_fighting)
	stdio.Fprintf(fl, "* Should closed doors be shown on autoexit / exit?\ndisp_closed_doors = %d\n\n", config_info.Play.Disp_closed_doors)
	stdio.Fprintf(fl, "* Should players be able to reroll stats at creation?\nreroll_stats = %d\n\n", config_info.Play.Reroll_player)
	stdio.Fprintf(fl, "* How many points in players initial points pool if using point pool creation?\ninitial_points = %d\n\n", config_info.Play.Initial_points)
	stdio.Fprintf(fl, "* Should compression be used if the client supports it?\ncompression = %d\n\n", config_info.Play.Enable_compression)
	stdio.Fprintf(fl, "* Should spoken languages be used?\nenable_languages = %d\n\n", config_info.Play.Enable_languages)
	stdio.Fprintf(fl, "* Should all items be treated as unique?\nall_items_unique = %d\n\n", config_info.Play.All_items_unique)
	stdio.Fprintf(fl, "* Amount of in game experience multiplier.\nexp_multiplier = %.2f\n\n", config_info.Play.Exp_multiplier)
	libc.StrCpy(&buf[0], config_info.Play.OK)
	strip_cr(&buf[0])
	stdio.Fprintf(fl, "* Text sent to players when OK is all that is needed.\nok = %s\n\n", &buf[0])
	libc.StrCpy(&buf[0], config_info.Play.NOPERSON)
	strip_cr(&buf[0])
	stdio.Fprintf(fl, "* Text sent to players when noone is available.\nnoperson = %s\n\n", &buf[0])
	libc.StrCpy(&buf[0], config_info.Play.NOEFFECT)
	strip_cr(&buf[0])
	stdio.Fprintf(fl, "* Text sent to players when an effect fails.\nnoeffect = %s\n", &buf[0])
	stdio.Fprintf(fl, "\n\n\n* [ Rent/Crashsave Options ]\n")
	stdio.Fprintf(fl, "* Should the MUD allow you to 'rent' for free?  (i.e. if you just quit,\n* your objects are saved at no cost, as in Merc-type MUDs.)\nfree_rent = %d\n\n", config_info.Csd.Free_rent)
	stdio.Fprintf(fl, "* Maximum number of items players are allowed to rent.\nmax_obj_save = %d\n\n", config_info.Csd.Max_obj_save)
	stdio.Fprintf(fl, "* Receptionist's surcharge on top of item costs.\nmin_rent_cost = %d\n\n", config_info.Csd.Min_rent_cost)
	stdio.Fprintf(fl, "* Should the game automatically save people?\nauto_save = %d\n\n", config_info.Csd.Auto_save)
	stdio.Fprintf(fl, "* If auto_save = 1, how often (in minutes) should the game save people's objects?\nautosave_time = %d\n\n", config_info.Csd.Autosave_time)
	stdio.Fprintf(fl, "* Lifetime of crashfiles and force-rent (idlesave) files in days.\ncrash_file_timeout = %d\n\n", config_info.Csd.Crash_file_timeout)
	stdio.Fprintf(fl, "* Lifetime of normal rent files in days.\nrent_file_timeout = %d\n\n", config_info.Csd.Rent_file_timeout)
	stdio.Fprintf(fl, "\n\n\n* [ Room Numbers ]\n")
	stdio.Fprintf(fl, "* The virtual number of the room that mortals should enter at.\nmortal_start_room = %d\n\n", config_info.Room_nums.Mortal_start_room)
	stdio.Fprintf(fl, "* The virtual number of the room that immorts should enter at.\nimmort_start_room = %d\n\n", config_info.Room_nums.Immort_start_room)
	stdio.Fprintf(fl, "* The virtual number of the room that frozen people should enter at.\nfrozen_start_room = %d\n\n", config_info.Room_nums.Frozen_start_room)
	stdio.Fprintf(fl, "* The virtual numbers of the donation rooms.  Note: Add donation rooms\n* sequentially (1 & 2 before 3). If you don't, you might not be able to\n* donate. Use -1 for 'no such room'.\ndonation_room_1 = %d\ndonation_room_2 = %d\ndonation_room_3 = %d\n\n", func() room_vnum {
		if config_info.Room_nums.Donation_room_1 != room_vnum(-1) {
			return config_info.Room_nums.Donation_room_1
		}
		return -1
	}(), func() room_vnum {
		if config_info.Room_nums.Donation_room_2 != room_vnum(-1) {
			return config_info.Room_nums.Donation_room_2
		}
		return -1
	}(), func() room_vnum {
		if config_info.Room_nums.Donation_room_3 != room_vnum(-1) {
			return config_info.Room_nums.Donation_room_3
		}
		return -1
	}())
	stdio.Fprintf(fl, "\n\n\n* [ Game Operation Options ]\n")
	stdio.Fprintf(fl, "* This is the default port on which the game should run if no port is\n* given on the command-line.  NOTE WELL: If you're using the\n* 'autorun' script, the port number there will override this setting.\n* Change the PORT= line in autorun instead of (or in addition to)\n* changing this.\nDFLT_PORT = %d\n\n", config_info.Operation.DFLT_PORT)
	if config_info.Operation.DFLT_IP != nil {
		libc.StrCpy(&buf[0], config_info.Operation.DFLT_IP)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* IP address to which the MUD should bind.\nDFLT_IP = %s\n\n", &buf[0])
	}
	if config_info.Operation.DFLT_DIR != nil {
		libc.StrCpy(&buf[0], config_info.Operation.DFLT_DIR)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* default directory to use as data directory.\nDFLT_DIR = %s\n\n", &buf[0])
	}
	if config_info.Operation.LOGNAME != nil {
		libc.StrCpy(&buf[0], config_info.Operation.LOGNAME)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* What file to log messages to (ex: 'log/syslog').\nLOGNAME = %s\n\n", &buf[0])
	}
	stdio.Fprintf(fl, "* Maximum number of players allowed before game starts to turn people away.\nmax_playing = %d\n\n", config_info.Operation.Max_playing)
	stdio.Fprintf(fl, "* Maximum size of bug, typo, and idea files in bytes (to prevent bombing).\nmax_filesize = %d\n\n", config_info.Operation.Max_filesize)
	stdio.Fprintf(fl, "* Maximum number of password attempts before disconnection.\nmax_bad_pws = %d\n\n", config_info.Operation.Max_bad_pws)
	stdio.Fprintf(fl, "* Is the site ok for everyone except those that are banned?\nsiteok_everyone = %d\n\n", config_info.Operation.Siteok_everyone)
	stdio.Fprintf(fl, "* If you want to use the original social file format\n* and disable Aedit, set to 0, otherwise, 1.\nuse_new_socials = %d\n\n", config_info.Operation.Use_new_socials)
	stdio.Fprintf(fl, "* If the nameserver is fast, set to 0, otherwise, 1.\nnameserver_is_slow = %d\n\n", config_info.Operation.Nameserver_is_slow)
	stdio.Fprintf(fl, "* Should OLC autosave to disk (1) or save internally (0).\nauto_save_olc = %d\n\n", config_info.Operation.Auto_save_olc)
	if config_info.Operation.MENU != nil {
		libc.StrCpy(&buf[0], config_info.Operation.MENU)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* The entrance/exit menu.\nMENU = \n%s~\n\n", &buf[0])
	}
	if config_info.Operation.WELC_MESSG != nil {
		libc.StrCpy(&buf[0], config_info.Operation.WELC_MESSG)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* The welcome message.\nWELC_MESSG = \n%s~\n\n", &buf[0])
	}
	if config_info.Operation.START_MESSG != nil {
		libc.StrCpy(&buf[0], config_info.Operation.START_MESSG)
		strip_cr(&buf[0])
		stdio.Fprintf(fl, "* NEWBIE start message.\nSTART_MESSG = \n%s~\n\n", &buf[0])
	}
	stdio.Fprintf(fl, "* Is the IMC global channel enabled (1) or not (0).\nimc_enabled = %d\n\n", config_info.Operation.Imc_enabled)
	stdio.Fprintf(fl, "\n\n\n* [ Autowiz Options ]\n")
	stdio.Fprintf(fl, "* Should the game automatically create a new wizlist/immlist every time\n* someone immorts, or is promoted to a higher (or lower) god level?\nuse_autowiz = %d\n\n", config_info.Autowiz.Use_autowiz)
	stdio.Fprintf(fl, "* If yes, what is the lowest level which should be on the wizlist?\nmin_wizlist_lev = %d\n\n", config_info.Autowiz.Min_wizlist_lev)
	stdio.Fprintf(fl, "\n\n\n* [ Character Advancement Options ]\n")
	stdio.Fprintf(fl, "* Should characters be allowed to advance in multiple classes?\nallow_multiclass = %d\n\n", config_info.Advance.Allow_multiclass)
	stdio.Fprintf(fl, "* Should characters be allowed to advance in prestige classes?\nallow_prestige = %d\n\n", config_info.Advance.Allow_prestige)
	stdio.Fprintf(fl, "\n\n\n* [ Game Speed Options ]\n")
	stdio.Fprintf(fl, "*Speed of Violence system. 0 - 10 (10) being the slowest.\npulse_viol = %d\n\n", config_info.Ticks.Pulse_violence)
	stdio.Fprintf(fl, "*Speed of Mobile actions. 0 - 20 (20) being the slowest.\npulse_mobile = %d\n\n", config_info.Ticks.Pulse_mobile)
	stdio.Fprintf(fl, "*Zone updates. 0 - 20 (20) being the slowest.\npulse_zone = %d\n\n", config_info.Ticks.Pulse_zone)
	stdio.Fprintf(fl, "*Autosave. 0 - 100 (100) being the slowest.\npulse_autosave = %d\n\n", config_info.Ticks.Pulse_autosave)
	stdio.Fprintf(fl, "*Idling passwords. 0 - 30 (30) being the slowest.\npulse_idlepwd = %d\n\n", config_info.Ticks.Pulse_idlepwd)
	stdio.Fprintf(fl, "*Sanity check time.. 0 - 50 (50) being the slowest.\npulse_sanity = %d\n\n", config_info.Ticks.Pulse_sanity)
	stdio.Fprintf(fl, "*Usage time.  0 - 10 (10) being the slowest in minutes.\npulse_usage = %d\n\n", config_info.Ticks.Pulse_usage)
	stdio.Fprintf(fl, "*Timesaving. 0 - 50 (50) being the slowest in minutes.\npulse_timesave = %d\n\n", config_info.Ticks.Pulse_timesave)
	stdio.Fprintf(fl, "*Current update time. 0 - 30 (30) being the slowest.\npulse_current = %d\n\n", config_info.Ticks.Pulse_current)
	stdio.Fprintf(fl, "*Character creation method.\nmethod = %d\n\n", config_info.Creation.Method)
	fl.Close()
	if in_save_list(zone_vnum(-1), SL_CFG) != 0 {
		remove_from_save_list(zone_vnum(-1), SL_CFG)
	}
	return TRUE
}
func cedit_disp_menu(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("OasisOLC MUD Configuration Editor\r\n@WA@B) @CAutowiz Options\r\n@WC@B) @CCrashsave/Rent Options\r\n@WD@B) @CaDvancement Options\r\n@WG@B) @CGame Play Options\r\n@WN@B) @CNew Character Creation Options\r\n@WO@B) @COperation Options\r\n@WR@B) @CRoom Numbers\r\n@WT@B) @CGame Tick Options\r\n@WQ@B) @CQuit\r\n@WEnter your choice : @n"))
	d.Olc.Mode = CEDIT_MAIN_MENU
}
func cedit_disp_game_play_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CPlayer Killing Allowed  : @c%s\r\n@WB@B) @CPlayer Thieving Allowed : @c%s\r\n@WC@B) @CMinimum Level To Shout  : @c%d\r\n@WD@B) @CHoller Move Cost        : @c%d\r\n@WE@B) @CTunnel Size             : @c%d\r\n@WF@B) @CMaximum Experience Gain : @c%d\r\n@WG@B) @CMaximum Experience Loss : @c%d\r\n@WH@B) @CMax Time for NPC Corpse : @c%d\r\n@WI@B) @CMax Time for PC Corpse  : @c%d\r\n@WJ@B) @CTics before PC sent to void : @c%d\r\n@WK@B) @CTics before PC is autosaved : @c%d\r\n@WL@B) @CLevel Immune To IDLE        : @c%d\r\n@WM@B) @CDeath Traps Junk Items      : @c%s\r\n@WN@B) @CObjects Load Into Inventory : @c%s\r\n@WO@B) @CTrack Through Doors         : @c%s\r\n@WP@B) @CDisplay Closed Doors        : @c%s\r\n@WR@B) @CUnattainable Level          : @c%d\r\n@WS@B) @CTreat all Objects as Unique : @c%s\r\n@WT@B) @CExperience multiplier       : @c%.2f\r\n@W1@B) @CStack Mobiles in room descs : @c%s\r\n@W2@B) @CStack Objects in room descs : @c%s\r\n@W3@B) @CAllow mobs to fight mobs    : @c%s\r\n@W4@B) @COK Message Text         : @c%s@W5@B) @CNOPERSON Message Text   : @c%s@W6@B) @CNOEFFECT Message Text   : @c%s@W8@B) @CAllow MCCP2 stream compression (recommended): @c%s\r\n@W9@B) @CAllow spoken languages : @c%s\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), func() string {
		if d.Olc.Config.Play.Pk_allowed == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Pt_allowed == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Play.Level_can_shout, d.Olc.Config.Play.Holler_move_cost, d.Olc.Config.Play.Tunnel_size, d.Olc.Config.Play.Max_exp_gain, d.Olc.Config.Play.Max_exp_loss, d.Olc.Config.Play.Max_npc_corpse_time, d.Olc.Config.Play.Max_pc_corpse_time, d.Olc.Config.Play.Idle_void, d.Olc.Config.Play.Idle_rent_time, d.Olc.Config.Play.Idle_max_level, func() string {
		if d.Olc.Config.Play.Dts_are_dumps == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Load_into_inventory == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Track_through_doors == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Disp_closed_doors == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Play.Level_cap, func() string {
		if d.Olc.Config.Play.All_items_unique == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Play.Exp_multiplier, func() string {
		if d.Olc.Config.Play.Stack_mobs == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Stack_objs == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Mob_fighting == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Play.OK, d.Olc.Config.Play.NOPERSON, d.Olc.Config.Play.NOEFFECT, func() string {
		if d.Olc.Config.Play.Enable_compression == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Play.Enable_languages == YES {
			return "Yes"
		}
		return "No"
	}())
	d.Olc.Mode = CEDIT_GAME_OPTIONS_MENU
}
func cedit_disp_crash_save_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CFree Rent          : @c%s\r\n@WB@B) @CMax Objects Saved  : @c%d\r\n@WC@B) @CMinimum Rent Cost  : @c%d\r\n@WD@B) @CAuto Save          : @c%s\r\n@WE@B) @CAuto Save Time     : @c%d minute(s)\r\n@WF@B) @CCrash File Timeout : @c%d day(s)\r\n@WG@B) @CRent File Timeout  : @c%d day(s)\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), func() string {
		if d.Olc.Config.Csd.Free_rent == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Csd.Max_obj_save, d.Olc.Config.Csd.Min_rent_cost, func() string {
		if d.Olc.Config.Csd.Auto_save == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Csd.Autosave_time, d.Olc.Config.Csd.Crash_file_timeout, d.Olc.Config.Csd.Rent_file_timeout)
	d.Olc.Mode = CEDIT_CRASHSAVE_OPTIONS_MENU
}
func cedit_disp_room_numbers(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CMortal Start Room   : @c%d\r\n@WB@B) @CImmortal Start Room : @c%d\r\n@WC@B) @CFrozen Start Room   : @c%d\r\n@W1@B) @CDonation Room #1    : @c%d\r\n@W2@B) @CDonation Room #2    : @c%d\r\n@W3@B) @CDonation Room #3    : @c%d\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), d.Olc.Config.Room_nums.Mortal_start_room, d.Olc.Config.Room_nums.Immort_start_room, d.Olc.Config.Room_nums.Frozen_start_room, d.Olc.Config.Room_nums.Donation_room_1, d.Olc.Config.Room_nums.Donation_room_2, d.Olc.Config.Room_nums.Donation_room_3)
	d.Olc.Mode = CEDIT_ROOM_NUMBERS_MENU
}
func cedit_disp_operation_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CDefault Port : @c%d\r\n@WB@B) @CDefault IP   : @c%s\r\n@WC@B) @CDefault Directory   : @c%s\r\n@WD@B) @CLogfile Name        : @c%s\r\n@WE@B) @CMax Players         : @c%d\r\n@WF@B) @CMax Filesize        : @c%d\r\n@WG@B) @CMax Bad Pws         : @c%d\r\n@WH@B) @CSite Ok Everyone    : @c%s\r\n@WI@B) @CName Server Is Slow : @c%s\r\n@WJ@B) @CUse new socials file: @c%s\r\n@WK@B) @COLC autosave to disk: @c%s\r\n@WL@B) @CMain Menu           : \r\n@n%s@n\r\n@WM@B) @CWelcome Message     : \r\n@n%s@n\r\n@WN@B) @CStart Message       : \r\n@n%s@n\r\n@WO@B) @CIMC Enabled         : @c%s@n\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), d.Olc.Config.Operation.DFLT_PORT, func() *byte {
		if d.Olc.Config.Operation.DFLT_IP != nil {
			return d.Olc.Config.Operation.DFLT_IP
		}
		return libc.CString("<None>")
	}(), func() *byte {
		if d.Olc.Config.Operation.DFLT_DIR != nil {
			return d.Olc.Config.Operation.DFLT_DIR
		}
		return libc.CString("<None>")
	}(), func() *byte {
		if d.Olc.Config.Operation.LOGNAME != nil {
			return d.Olc.Config.Operation.LOGNAME
		}
		return libc.CString("<None>")
	}(), d.Olc.Config.Operation.Max_playing, d.Olc.Config.Operation.Max_filesize, d.Olc.Config.Operation.Max_bad_pws, func() string {
		if d.Olc.Config.Operation.Siteok_everyone != 0 {
			return "YES"
		}
		return "NO"
	}(), func() string {
		if d.Olc.Config.Operation.Nameserver_is_slow != 0 {
			return "YES"
		}
		return "NO"
	}(), func() string {
		if d.Olc.Config.Operation.Use_new_socials != 0 {
			return "YES"
		}
		return "NO"
	}(), func() string {
		if d.Olc.Config.Operation.Auto_save_olc != 0 {
			return "YES"
		}
		return "NO"
	}(), func() *byte {
		if d.Olc.Config.Operation.MENU != nil {
			return d.Olc.Config.Operation.MENU
		}
		return libc.CString("<None>")
	}(), func() *byte {
		if d.Olc.Config.Operation.WELC_MESSG != nil {
			return d.Olc.Config.Operation.WELC_MESSG
		}
		return libc.CString("<None>")
	}(), func() *byte {
		if d.Olc.Config.Operation.START_MESSG != nil {
			return d.Olc.Config.Operation.START_MESSG
		}
		return libc.CString("<None>")
	}(), func() string {
		if d.Olc.Config.Operation.Imc_enabled != 0 {
			return "YES"
		}
		return "NO"
	}())
	d.Olc.Mode = CEDIT_OPERATION_OPTIONS_MENU
}
func cedit_disp_autowiz_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CUse the autowiz        : @c%s\r\n@WB@B) @CMinimum wizlist level  : @c%d\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), func() string {
		if d.Olc.Config.Autowiz.Use_autowiz == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Autowiz.Min_wizlist_lev)
	d.Olc.Mode = CEDIT_AUTOWIZ_OPTIONS_MENU
}
func cedit_disp_advance_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CAllow multiclass       : @c%s\r\n@WB@B) @CAllow Prestige Classes : @c%s\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), func() string {
		if d.Olc.Config.Advance.Allow_multiclass == YES {
			return "Yes"
		}
		return "No"
	}(), func() string {
		if d.Olc.Config.Advance.Allow_prestige == YES {
			return "Yes"
		}
		return "No"
	}())
	d.Olc.Mode = CEDIT_ADVANCE_OPTIONS_MENU
}
func cedit_disp_ticks_menu(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CPulse Violence Time  : @c%d\r\n@WB@B) @CMobile Action Time   : @c%d\r\n@WC@B) @CZone Time            : @c%d\r\n@WD@B) @CAutosave Time        : @c%d\r\n@WE@B) @CIdle Password Time   : @c%d\r\n@WF@B) @CSanity Check Time    : @c%d\r\n@WG@B) @CUsage Check Time     : @c%d\r\n@WH@B) @CTimesave Time        : @c%d\r\n@WI@B) @CCurrents Update Time : @c%d\r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), d.Olc.Config.Ticks.Pulse_violence, d.Olc.Config.Ticks.Pulse_mobile, d.Olc.Config.Ticks.Pulse_zone, d.Olc.Config.Ticks.Pulse_autosave, d.Olc.Config.Ticks.Pulse_idlepwd, d.Olc.Config.Ticks.Pulse_sanity, d.Olc.Config.Ticks.Pulse_usage, d.Olc.Config.Ticks.Pulse_timesave, d.Olc.Config.Ticks.Pulse_current)
	d.Olc.Mode = CEDIT_TICKS_OPTIONS_MENU
}
func cedit_disp_creation_options(d *descriptor_data) {
	clear_screen(d)
	write_to_output(d, libc.CString("\r\n\r\n@WA@B) @CCharacter Creation Method : @c%d \r\n   %s\r\n@WB@B) @CPlayers can reroll stats on creation : @c%s\r\n@WC@B) @CNumber of points in initial points pool for pool generation methods : @c%d \r\n@WQ@B) @CExit To The Main Menu\r\n@WEnter your choice : @n"), d.Olc.Config.Creation.Method, creation_methods[d.Olc.Config.Creation.Method], func() string {
		if d.Olc.Config.Play.Reroll_player == YES {
			return "Yes"
		}
		return "No"
	}(), d.Olc.Config.Play.Initial_points)
	d.Olc.Mode = CEDIT_CREATION_OPTIONS_MENU
}
func cedit_disp_creation_menu(d *descriptor_data) {
	var (
		i   int
		buf [2048]byte
	)
	clear_screen(d)
	for i = 0; i < NUM_CREATION_METHODS; i++ {
		stdio.Sprintf(&buf[0], "@W%d@B) @C%s@n\r\n", i, creation_methods[i])
		write_to_output(d, &buf[0])
	}
	write_to_output(d, libc.CString("Choose character creation type: "))
	d.Olc.Mode = CEDIT_CREATION_MENU
}
func cedit_disp_points_menu(d *descriptor_data) {
	write_to_output(d, libc.CString("Enter size of initial points pool. (0 or greater). "))
	d.Olc.Mode = CEDIT_POINTS_MENU
}
func cedit_parse(d *descriptor_data, arg *byte) {
	var oldtext *byte = nil
	switch d.Olc.Mode {
	case CEDIT_CONFIRM_SAVESTRING:
		switch *arg {
		case 'y':
			fallthrough
		case 'Y':
			cedit_save_internally(d)
			mudlog(CMP, MAX(ADMLVL_BUILDER, int(d.Character.Player_specials.Invis_level)), TRUE, libc.CString("OLC: %s modifies the game configuration."), GET_NAME(d.Character))
			cleanup_olc(d, CLEANUP_CONFIG)
			if config_info.Csd.Auto_save != 0 {
				cedit_save_to_disk()
				write_to_output(d, libc.CString("Game configuration saved to disk.\r\n"))
			} else {
				write_to_output(d, libc.CString("Game configuration saved to memory.\r\n"))
			}
			return
		case 'n':
			fallthrough
		case 'N':
			write_to_output(d, libc.CString("Game configuration not saved to memory.\r\n"))
			cleanup_olc(d, CLEANUP_CONFIG)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
			write_to_output(d, libc.CString("Do you wish to save your changes? : "))
			return
		}
		fallthrough
	case CEDIT_MAIN_MENU:
		switch *arg {
		case 'g':
			fallthrough
		case 'G':
			cedit_disp_game_play_options(d)
			d.Olc.Mode = CEDIT_GAME_OPTIONS_MENU
		case 'c':
			fallthrough
		case 'C':
			cedit_disp_crash_save_options(d)
			d.Olc.Mode = CEDIT_CRASHSAVE_OPTIONS_MENU
		case 'r':
			fallthrough
		case 'R':
			cedit_disp_room_numbers(d)
			d.Olc.Mode = CEDIT_ROOM_NUMBERS_MENU
		case 'o':
			fallthrough
		case 'O':
			cedit_disp_operation_options(d)
			d.Olc.Mode = CEDIT_OPERATION_OPTIONS_MENU
		case 'a':
			fallthrough
		case 'A':
			cedit_disp_autowiz_options(d)
			d.Olc.Mode = CEDIT_AUTOWIZ_OPTIONS_MENU
		case 'd':
			fallthrough
		case 'D':
			cedit_disp_advance_options(d)
			d.Olc.Mode = CEDIT_ADVANCE_OPTIONS_MENU
		case 't':
			fallthrough
		case 'T':
			cedit_disp_ticks_menu(d)
			d.Olc.Mode = CEDIT_TICKS_OPTIONS_MENU
		case 'q':
			fallthrough
		case 'Q':
			write_to_output(d, libc.CString("Do you wish to save your changes? : "))
			d.Olc.Mode = CEDIT_CONFIRM_SAVESTRING
		case 'n':
			fallthrough
		case 'N':
			cedit_disp_creation_options(d)
			d.Olc.Mode = CEDIT_CREATION_OPTIONS_MENU
		default:
			write_to_output(d, libc.CString("That is an invalid choice!\r\n"))
			cedit_disp_menu(d)
		}
	case CEDIT_GAME_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			if d.Olc.Config.Play.Pk_allowed == YES {
				d.Olc.Config.Play.Pk_allowed = NO
			} else {
				d.Olc.Config.Play.Pk_allowed = YES
			}
		case 'b':
			fallthrough
		case 'B':
			if d.Olc.Config.Play.Pt_allowed == YES {
				d.Olc.Config.Play.Pt_allowed = NO
			} else {
				d.Olc.Config.Play.Pt_allowed = YES
			}
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter the minimum level a player must be to shout, gossip, etc : "))
			d.Olc.Mode = CEDIT_LEVEL_CAN_SHOUT
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("Enter the amount it costs (in move points) to holler : "))
			d.Olc.Mode = CEDIT_HOLLER_MOVE_COST
			return
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Enter the maximum number of people allowed in a tunnel : "))
			d.Olc.Mode = CEDIT_TUNNEL_SIZE
			return
		case 'f':
			fallthrough
		case 'F':
			write_to_output(d, libc.CString("Enter the maximum gain of experience per kill for players : "))
			d.Olc.Mode = CEDIT_MAX_EXP_GAIN
			return
		case 'g':
			fallthrough
		case 'G':
			write_to_output(d, libc.CString("Enter the maximum loss of experience per death for players : "))
			d.Olc.Mode = CEDIT_MAX_EXP_LOSS
			return
		case 'h':
			fallthrough
		case 'H':
			write_to_output(d, libc.CString("Enter the number of tics before NPC corpses decompose : "))
			d.Olc.Mode = CEDIT_MAX_NPC_CORPSE_TIME
			return
		case 'i':
			fallthrough
		case 'I':
			write_to_output(d, libc.CString("Enter the number of tics before PC corpses decompose : "))
			d.Olc.Mode = CEDIT_MAX_PC_CORPSE_TIME
			return
		case 'j':
			fallthrough
		case 'J':
			write_to_output(d, libc.CString("Enter the number of tics before PC's are sent to the void (idle) : "))
			d.Olc.Mode = CEDIT_IDLE_VOID
			return
		case 'k':
			fallthrough
		case 'K':
			write_to_output(d, libc.CString("Enter the number of tics before PC's are automatically rented and forced to quit : "))
			d.Olc.Mode = CEDIT_IDLE_RENT_TIME
			return
		case 'l':
			fallthrough
		case 'L':
			write_to_output(d, libc.CString("Enter the level a player must be to become immune to IDLE : "))
			d.Olc.Mode = CEDIT_IDLE_MAX_LEVEL
			return
		case 'm':
			fallthrough
		case 'M':
			if d.Olc.Config.Play.Dts_are_dumps == YES {
				d.Olc.Config.Play.Dts_are_dumps = NO
			} else {
				d.Olc.Config.Play.Dts_are_dumps = YES
			}
		case 'n':
			fallthrough
		case 'N':
			if d.Olc.Config.Play.Load_into_inventory == YES {
				d.Olc.Config.Play.Load_into_inventory = NO
			} else {
				d.Olc.Config.Play.Load_into_inventory = YES
			}
		case 'o':
			fallthrough
		case 'O':
			if d.Olc.Config.Play.Track_through_doors == YES {
				d.Olc.Config.Play.Track_through_doors = NO
			} else {
				d.Olc.Config.Play.Track_through_doors = YES
			}
		case 'p':
			fallthrough
		case 'P':
			if d.Olc.Config.Play.Disp_closed_doors == YES {
				d.Olc.Config.Play.Disp_closed_doors = NO
			} else {
				d.Olc.Config.Play.Disp_closed_doors = YES
			}
		case 'r':
			fallthrough
		case 'R':
			write_to_output(d, libc.CString("Enter the number a character cannot level to : "))
			d.Olc.Mode = CEDIT_LEVEL_CAP
			return
		case 's':
			fallthrough
		case 'S':
			if d.Olc.Config.Play.All_items_unique == YES {
				d.Olc.Config.Play.All_items_unique = NO
			} else {
				d.Olc.Config.Play.All_items_unique = YES
			}
		case 't':
			fallthrough
		case 'T':
			write_to_output(d, libc.CString("Enter the multiplier a player will recieve to experience gained : "))
			d.Olc.Mode = CEDIT_EXP_MULTIPLIER
			return
		case '1':
			if d.Olc.Config.Play.Stack_mobs == YES {
				d.Olc.Config.Play.Stack_mobs = NO
			} else {
				d.Olc.Config.Play.Stack_mobs = YES
			}
		case '2':
			if d.Olc.Config.Play.Stack_objs == YES {
				d.Olc.Config.Play.Stack_objs = NO
			} else {
				d.Olc.Config.Play.Stack_objs = YES
			}
		case '3':
			if d.Olc.Config.Play.Mob_fighting == YES {
				d.Olc.Config.Play.Mob_fighting = NO
			} else {
				d.Olc.Config.Play.Mob_fighting = YES
			}
		case '4':
			write_to_output(d, libc.CString("Enter the OK message : "))
			d.Olc.Mode = CEDIT_OK
			return
		case '5':
			write_to_output(d, libc.CString("Enter the NOPERSON message : "))
			d.Olc.Mode = CEDIT_NOPERSON
			return
		case '6':
			write_to_output(d, libc.CString("Enter the NOEFFECT message : "))
			d.Olc.Mode = CEDIT_NOEFFECT
			return
		case '8':
			if d.Olc.Config.Play.Enable_compression == YES {
				d.Olc.Config.Play.Enable_compression = NO
			} else {
				d.Olc.Config.Play.Enable_compression = YES
			}
		case '9':
			if d.Olc.Config.Play.Enable_languages == YES {
				d.Olc.Config.Play.Enable_languages = NO
			} else {
				d.Olc.Config.Play.Enable_languages = YES
			}
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
			cedit_disp_game_play_options(d)
		}
		cedit_disp_game_play_options(d)
		return
	case CEDIT_CRASHSAVE_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			if d.Olc.Config.Csd.Free_rent == YES {
				d.Olc.Config.Csd.Free_rent = NO
			} else {
				d.Olc.Config.Csd.Free_rent = YES
			}
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter the maximum number of items players can rent : "))
			d.Olc.Mode = CEDIT_MAX_OBJ_SAVE
			return
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter the surcharge on top of item costs : "))
			d.Olc.Mode = CEDIT_MIN_RENT_COST
			return
		case 'd':
			fallthrough
		case 'D':
			if d.Olc.Config.Csd.Auto_save == YES {
				d.Olc.Config.Csd.Auto_save = NO
			} else {
				d.Olc.Config.Csd.Auto_save = YES
			}
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Enter how often (in minutes) should the MUD save players : "))
			d.Olc.Mode = CEDIT_AUTOSAVE_TIME
			return
		case 'f':
			fallthrough
		case 'F':
			write_to_output(d, libc.CString("Enter the lifetime of crash and idlesave files (days) : "))
			d.Olc.Mode = CEDIT_CRASH_FILE_TIMEOUT
			return
		case 'g':
			fallthrough
		case 'G':
			write_to_output(d, libc.CString("Enter the lifetime of normal rent files (days) : "))
			d.Olc.Mode = CEDIT_RENT_FILE_TIMEOUT
			return
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_crash_save_options(d)
		return
	case CEDIT_ROOM_NUMBERS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter the room's vnum where mortals should load into : "))
			d.Olc.Mode = CEDIT_MORTAL_START_ROOM
			return
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter the room's vnum where immortals should load into : "))
			d.Olc.Mode = CEDIT_IMMORT_START_ROOM
			return
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter the room's vnum where frozen people should load into : "))
			d.Olc.Mode = CEDIT_FROZEN_START_ROOM
			return
		case '1':
			write_to_output(d, libc.CString("Enter the vnum for donation room #1 : "))
			d.Olc.Mode = CEDIT_DONATION_ROOM_1
			return
		case '2':
			write_to_output(d, libc.CString("Enter the vnum for donation room #2 : "))
			d.Olc.Mode = CEDIT_DONATION_ROOM_2
			return
		case '3':
			write_to_output(d, libc.CString("Enter the vnum for donation room #3 : "))
			d.Olc.Mode = CEDIT_DONATION_ROOM_3
			return
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_room_numbers(d)
		return
	case CEDIT_OPERATION_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter the default port number : "))
			d.Olc.Mode = CEDIT_DFLT_PORT
			return
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter the default IP Address : "))
			d.Olc.Mode = CEDIT_DFLT_IP
			return
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter the default directory : "))
			d.Olc.Mode = CEDIT_DFLT_DIR
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("Enter the name of the logfile : "))
			d.Olc.Mode = CEDIT_LOGNAME
			return
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Enter the maximum number of players : "))
			d.Olc.Mode = CEDIT_MAX_PLAYING
			return
		case 'f':
			fallthrough
		case 'F':
			write_to_output(d, libc.CString("Enter the maximum size of the logs : "))
			d.Olc.Mode = CEDIT_MAX_FILESIZE
			return
		case 'g':
			fallthrough
		case 'G':
			write_to_output(d, libc.CString("Enter the maximum number of password attempts : "))
			d.Olc.Mode = CEDIT_MAX_BAD_PWS
			return
		case 'h':
			fallthrough
		case 'H':
			if d.Olc.Config.Operation.Siteok_everyone == YES {
				d.Olc.Config.Operation.Siteok_everyone = NO
			} else {
				d.Olc.Config.Operation.Siteok_everyone = YES
			}
		case 'i':
			fallthrough
		case 'I':
			if d.Olc.Config.Operation.Nameserver_is_slow == YES {
				d.Olc.Config.Operation.Nameserver_is_slow = NO
			} else {
				d.Olc.Config.Operation.Nameserver_is_slow = YES
			}
		case 'j':
			fallthrough
		case 'J':
			if d.Olc.Config.Operation.Use_new_socials == YES {
				d.Olc.Config.Operation.Use_new_socials = NO
			} else {
				d.Olc.Config.Operation.Use_new_socials = YES
			}
			send_to_char(d.Character, libc.CString("Please note that using the stock social file will disable AEDIT.\r\n"))
		case 'k':
			fallthrough
		case 'K':
			if d.Olc.Config.Operation.Auto_save_olc == YES {
				d.Olc.Config.Operation.Auto_save_olc = NO
			} else {
				d.Olc.Config.Operation.Auto_save_olc = YES
			}
		case 'l':
			fallthrough
		case 'L':
			d.Olc.Mode = CEDIT_MENU
			clear_screen(d)
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter the new MENU :\r\n\r\n"))
			if d.Olc.Config.Operation.MENU != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Config.Operation.MENU)
				oldtext = libc.StrDup(d.Olc.Config.Operation.MENU)
			}
			string_write(d, &d.Olc.Config.Operation.MENU, MAX_INPUT_LENGTH, 0, unsafe.Pointer(oldtext))
			return
		case 'm':
			fallthrough
		case 'M':
			d.Olc.Mode = CEDIT_WELC_MESSG
			clear_screen(d)
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter the new welcome message :\r\n\r\n"))
			if d.Olc.Config.Operation.WELC_MESSG != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Config.Operation.WELC_MESSG)
				oldtext = str_udup(d.Olc.Config.Operation.WELC_MESSG)
			}
			string_write(d, &d.Olc.Config.Operation.WELC_MESSG, MAX_INPUT_LENGTH, 0, unsafe.Pointer(oldtext))
			return
		case 'n':
			fallthrough
		case 'N':
			d.Olc.Mode = CEDIT_START_MESSG
			clear_screen(d)
			send_editor_help(d)
			write_to_output(d, libc.CString("Enter the new newbie start message :\r\n\r\n"))
			if d.Olc.Config.Operation.START_MESSG != nil {
				write_to_output(d, libc.CString("%s"), d.Olc.Config.Operation.START_MESSG)
				oldtext = libc.StrDup(d.Olc.Config.Operation.START_MESSG)
			}
			string_write(d, &d.Olc.Config.Operation.START_MESSG, MAX_INPUT_LENGTH, 0, unsafe.Pointer(oldtext))
			return
		case 'o':
			fallthrough
		case 'O':
			if d.Olc.Config.Operation.Imc_enabled == YES {
				d.Olc.Config.Operation.Imc_enabled = NO
			} else {
				d.Olc.Config.Operation.Imc_enabled = YES
			}
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_operation_options(d)
		return
	case CEDIT_AUTOWIZ_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			if d.Olc.Config.Autowiz.Use_autowiz == YES {
				d.Olc.Config.Autowiz.Use_autowiz = NO
			} else {
				d.Olc.Config.Autowiz.Use_autowiz = YES
			}
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter the minimum level for players to appear on the wizlist : "))
			d.Olc.Mode = CEDIT_MIN_WIZLIST_LEV
			return
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_autowiz_options(d)
		return
	case CEDIT_ADVANCE_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			if d.Olc.Config.Advance.Allow_multiclass == YES {
				d.Olc.Config.Advance.Allow_multiclass = NO
			} else {
				d.Olc.Config.Advance.Allow_multiclass = YES
			}
		case 'b':
			fallthrough
		case 'B':
			if d.Olc.Config.Advance.Allow_prestige == YES {
				d.Olc.Config.Advance.Allow_prestige = NO
			} else {
				d.Olc.Config.Advance.Allow_prestige = YES
			}
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_advance_options(d)
		return
	case CEDIT_TICKS_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			write_to_output(d, libc.CString("Enter the Speed of the violence system (1 - 10) 10 = Slowest : "))
			d.Olc.Mode = CEDIT_PULSE_VIOLENCE
			return
		case 'b':
			fallthrough
		case 'B':
			write_to_output(d, libc.CString("Enter the Speed of the Mobile Actions (1 - 20) 20 = Slowest : "))
			d.Olc.Mode = CEDIT_PULSE_MOBILE
			return
		case 'c':
			fallthrough
		case 'C':
			write_to_output(d, libc.CString("Enter Zone update time : "))
			d.Olc.Mode = CEDIT_PULSE_ZONE
			return
		case 'd':
			fallthrough
		case 'D':
			write_to_output(d, libc.CString("Enter time for autosaving : "))
			d.Olc.Mode = CEDIT_PULSE_AUTOSAVE
			return
		case 'e':
			fallthrough
		case 'E':
			write_to_output(d, libc.CString("Enter the time to kill connection waiting for password : "))
			d.Olc.Mode = CEDIT_PULSE_IDLEPWD
			return
		case 'f':
			fallthrough
		case 'F':
			write_to_output(d, libc.CString("Enter the maximum size of the logs : "))
			d.Olc.Mode = CEDIT_PULSE_SANITY
			return
		case 'g':
			fallthrough
		case 'G':
			write_to_output(d, libc.CString("Enter the maximum number of password attempts : "))
			d.Olc.Mode = CEDIT_PULSE_USAGE
			return
		case 'h':
			fallthrough
		case 'H':
			write_to_output(d, libc.CString("Enter timesave : "))
			d.Olc.Mode = CEDIT_PULSE_TIMESAVE
			return
		case 'i':
			fallthrough
		case 'I':
			write_to_output(d, libc.CString("Enter Current update time : "))
			d.Olc.Mode = CEDIT_PULSE_CURRENT
			return
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
		}
		cedit_disp_ticks_menu(d)
		return
	case CEDIT_LEVEL_CAN_SHOUT:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the minimum level a player must be to shout, gossip, etc : "))
		} else {
			d.Olc.Config.Play.Level_can_shout = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_HOLLER_MOVE_COST:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the amount it costs (in move points) to holler : "))
		} else {
			d.Olc.Config.Play.Holler_move_cost = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_TUNNEL_SIZE:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the maximum number of people allowed in a tunnel : "))
		} else {
			d.Olc.Config.Play.Tunnel_size = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_MAX_EXP_GAIN:
		if *arg != 0 {
			d.Olc.Config.Play.Max_exp_gain = libc.Atoi(libc.GoString(arg))
		}
		cedit_disp_game_play_options(d)
	case CEDIT_MAX_EXP_LOSS:
		if *arg != 0 {
			d.Olc.Config.Play.Max_exp_loss = libc.Atoi(libc.GoString(arg))
		}
		cedit_disp_game_play_options(d)
	case CEDIT_MAX_NPC_CORPSE_TIME:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the number of tics before NPC corpses decompose : "))
		} else {
			d.Olc.Config.Play.Max_npc_corpse_time = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_MAX_PC_CORPSE_TIME:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the number of tics before PC corpses decompose : "))
		} else {
			d.Olc.Config.Play.Max_pc_corpse_time = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_IDLE_VOID:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the number of tics before PC's are sent to the void (idle) : "))
		} else {
			d.Olc.Config.Play.Idle_void = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_IDLE_RENT_TIME:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the number of tics before PC's are automatically rented and forced to quit : "))
		} else {
			d.Olc.Config.Play.Idle_rent_time = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_IDLE_MAX_LEVEL:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the level a player must be to become immune to IDLE : "))
		} else {
			d.Olc.Config.Play.Idle_max_level = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_LEVEL_CAP:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the level a player cannot reach : "))
		} else {
			d.Olc.Config.Play.Level_cap = libc.Atoi(libc.GoString(arg))
			cedit_disp_game_play_options(d)
		}
	case CEDIT_OK:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Config.Play.OK != nil {
			libc.Free(unsafe.Pointer(d.Olc.Config.Play.OK))
		}
		d.Olc.Config.Play.OK = str_udup(arg)
		libc.StrCat(d.Olc.Config.Play.OK, libc.CString("\r\n"))
		cedit_disp_game_play_options(d)
	case CEDIT_NOPERSON:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Config.Play.NOPERSON != nil {
			libc.Free(unsafe.Pointer(d.Olc.Config.Play.NOPERSON))
		}
		d.Olc.Config.Play.NOPERSON = str_udup(arg)
		libc.StrCat(d.Olc.Config.Play.NOPERSON, libc.CString("\r\n"))
		cedit_disp_game_play_options(d)
	case CEDIT_NOEFFECT:
		if genolc_checkstring(d, arg) == 0 {
			break
		}
		if d.Olc.Config.Play.NOEFFECT != nil {
			libc.Free(unsafe.Pointer(d.Olc.Config.Play.NOEFFECT))
		}
		d.Olc.Config.Play.NOEFFECT = str_udup(arg)
		libc.StrCat(d.Olc.Config.Play.NOEFFECT, libc.CString("\r\n"))
		cedit_disp_game_play_options(d)
	case CEDIT_MAX_OBJ_SAVE:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the maximum objects a player can save : "))
		} else {
			d.Olc.Config.Csd.Max_obj_save = libc.Atoi(libc.GoString(arg))
			cedit_disp_crash_save_options(d)
		}
	case CEDIT_MIN_RENT_COST:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the minimum amount it costs to rent : "))
		} else {
			d.Olc.Config.Csd.Min_rent_cost = libc.Atoi(libc.GoString(arg))
			cedit_disp_crash_save_options(d)
		}
	case CEDIT_AUTOSAVE_TIME:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the interval for player's being autosaved : "))
		} else {
			d.Olc.Config.Csd.Autosave_time = libc.Atoi(libc.GoString(arg))
			cedit_disp_crash_save_options(d)
		}
	case CEDIT_CRASH_FILE_TIMEOUT:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the lifetime of crash and idlesave files (days) : "))
		} else {
			d.Olc.Config.Csd.Crash_file_timeout = libc.Atoi(libc.GoString(arg))
			cedit_disp_crash_save_options(d)
		}
	case CEDIT_RENT_FILE_TIMEOUT:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the lifetime of rent files (days) : "))
		} else {
			d.Olc.Config.Csd.Rent_file_timeout = libc.Atoi(libc.GoString(arg))
			cedit_disp_crash_save_options(d)
		}
	case CEDIT_MORTAL_START_ROOM:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the room's vnum where mortals should load into : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the room's vnum where mortals should load into : "))
		} else {
			d.Olc.Config.Room_nums.Mortal_start_room = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_IMMORT_START_ROOM:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the room's vnum where immortals should load into : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the room's vnum where immortals should load into : "))
		} else {
			d.Olc.Config.Room_nums.Immort_start_room = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_FROZEN_START_ROOM:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the room's vnum where frozen people should load into : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the room's vnum where frozen people should load into : "))
		} else {
			d.Olc.Config.Room_nums.Frozen_start_room = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_DONATION_ROOM_1:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the vnum for donation room #1 : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the vnum for donation room #1 : "))
		} else {
			d.Olc.Config.Room_nums.Donation_room_1 = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_DONATION_ROOM_2:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the vnum for donation room #2 : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the vnum for donation room #2 : "))
		} else {
			d.Olc.Config.Room_nums.Donation_room_2 = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_DONATION_ROOM_3:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the vnum for donation room #3 : "))
		} else if real_room(room_vnum(libc.Atoi(libc.GoString(arg)))) == room_rnum(-1) {
			write_to_output(d, libc.CString("That room doesn't exist!\r\nEnter the vnum for donation room #3 : "))
		} else {
			d.Olc.Config.Room_nums.Donation_room_3 = room_vnum(libc.Atoi(libc.GoString(arg)))
			cedit_disp_room_numbers(d)
		}
	case CEDIT_DFLT_PORT:
		d.Olc.Config.Operation.DFLT_PORT = uint16(int16(libc.Atoi(libc.GoString(arg))))
		cedit_disp_operation_options(d)
	case CEDIT_DFLT_IP:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the default ip address : "))
		} else {
			d.Olc.Config.Operation.DFLT_IP = str_udup(arg)
			cedit_disp_operation_options(d)
		}
	case CEDIT_DFLT_DIR:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the default directory : "))
		} else {
			d.Olc.Config.Operation.DFLT_DIR = str_udup(arg)
			cedit_disp_operation_options(d)
		}
	case CEDIT_LOGNAME:
		if *arg == 0 {
			write_to_output(d, libc.CString("That is an invalid choice!\r\nEnter the name of the logfile : "))
		} else {
			d.Olc.Config.Operation.LOGNAME = str_udup(arg)
			cedit_disp_operation_options(d)
		}
	case CEDIT_MAX_PLAYING:
		d.Olc.Config.Operation.Max_playing = libc.Atoi(libc.GoString(arg))
		cedit_disp_operation_options(d)
	case CEDIT_MAX_FILESIZE:
		d.Olc.Config.Operation.Max_filesize = libc.Atoi(libc.GoString(arg))
		cedit_disp_operation_options(d)
	case CEDIT_MAX_BAD_PWS:
		d.Olc.Config.Operation.Max_bad_pws = libc.Atoi(libc.GoString(arg))
		cedit_disp_operation_options(d)
	case CEDIT_MIN_WIZLIST_LEV:
		if libc.Atoi(libc.GoString(arg)) > ADMLVL_IMPL {
			write_to_output(d, libc.CString("The minimum wizlist level can't be greater than %d.\r\nEnter the minimum level for players to appear on the wizlist : "), ADMLVL_IMPL)
		} else {
			d.Olc.Config.Autowiz.Min_wizlist_lev = libc.Atoi(libc.GoString(arg))
			cedit_disp_autowiz_options(d)
		}
	case CEDIT_PULSE_VIOLENCE:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 10 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 10.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_violence = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_MOBILE:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 20 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 20.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_mobile = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_IDLEPWD:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 30 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 30.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_idlepwd = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_ZONE:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 20 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 20.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_zone = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_AUTOSAVE:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 100 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 100.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_autosave = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_CURRENT:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 30 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 30.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_current = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_SANITY:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 50 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 50.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_sanity = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_PULSE_USAGE:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > 10 {
			write_to_output(d, libc.CString("Please enter a number between 0 - 10.\r\n"))
			cedit_disp_ticks_menu(d)
		} else {
			d.Olc.Config.Ticks.Pulse_usage = libc.Atoi(libc.GoString(arg))
			cedit_disp_ticks_menu(d)
		}
	case CEDIT_EXP_MULTIPLIER:
		if *arg != 0 {
			d.Olc.Config.Play.Exp_multiplier = float32(libc.Atof(libc.GoString(arg)))
		}
		cedit_disp_game_play_options(d)
	case CEDIT_CREATION_OPTIONS_MENU:
		switch *arg {
		case 'a':
			fallthrough
		case 'A':
			cedit_disp_creation_menu(d)
			d.Olc.Mode = CEDIT_CREATION_MENU
		case 'b':
			fallthrough
		case 'B':
			if d.Olc.Config.Play.Reroll_player == YES {
				d.Olc.Config.Play.Reroll_player = NO
			} else {
				d.Olc.Config.Play.Reroll_player = YES
			}
			cedit_disp_creation_options(d)
		case 'c':
			fallthrough
		case 'C':
			cedit_disp_points_menu(d)
			d.Olc.Mode = CEDIT_POINTS_MENU
		case 'q':
			fallthrough
		case 'Q':
			cedit_disp_menu(d)
			return
		default:
			write_to_output(d, libc.CString("\r\nThat is an invalid choice!\r\n"))
			cedit_disp_creation_menu(d)
		}
		return
	case CEDIT_CREATION_MENU:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 || libc.Atoi(libc.GoString(arg)) > int(NUM_CREATION_METHODS-1) {
			write_to_output(d, libc.CString("Please enter a number between 0 - %d.\r\n"), int(NUM_CREATION_METHODS-1))
			cedit_disp_creation_menu(d)
		} else {
			d.Olc.Config.Creation.Method = libc.Atoi(libc.GoString(arg))
			cedit_disp_creation_options(d)
		}
	case CEDIT_POINTS_MENU:
		if !unicode.IsDigit(rune(*arg)) || libc.Atoi(libc.GoString(arg)) < 0 {
			write_to_output(d, libc.CString("Please enter a number 0 or higher.\r\n"))
			cedit_disp_points_menu(d)
		} else {
			d.Olc.Config.Play.Initial_points = libc.Atoi(libc.GoString(arg))
			cedit_disp_creation_options(d)
		}
	default:
		cleanup_olc(d, CLEANUP_CONFIG)
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: OLC: cedit_parse(): Reached default case!"))
		write_to_output(d, libc.CString("Oops...\r\n"))
	}
}
func reassign_rooms() {
	var i int
	for i = 0; i < int(top_of_world); i++ {
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Func = nil
	}
	assign_rooms()
}
func cedit_string_cleanup(d *descriptor_data, terminator int) {
	switch d.Olc.Mode {
	case CEDIT_MENU:
		fallthrough
	case CEDIT_WELC_MESSG:
		fallthrough
	case CEDIT_START_MESSG:
		cedit_disp_operation_options(d)
	}
}
