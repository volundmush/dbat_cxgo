package main

import (
	"github.com/gotranspile/cxgo/runtime/cnet"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"os"
	"unicode"
	"unsafe"
)

const INVALID_SOCKET = -1
const COMPRESS2 = 86
const NUM_RESERVED_DESCS = 8
const COPYOVER_FILE = "copyover.dat"
const TO_ROOM = 1
const TO_VICT = 2
const TO_NOTVICT = 3
const TO_CHAR = 4
const TO_GMOTE = 5
const TO_SLEEP = 256
const DG_NO_TRIG = 512
const TO_SNEAKRESIST = 1024
const TO_HIDERESIST = 2048
const PAGE_LENGTH = 22
const PAGE_WIDTH = 79

var descriptor_list *descriptor_data = nil
var bufpool *txt_block = nil
var buf_largecount int = 0
var buf_overflows int = 0
var buf_switches int = 0
var circle_shutdown int = 0
var circle_reboot int = 0
var no_specials int = 0
var max_players int = 0
var tics_passed int = 0
var scheck int = 0
var null_time libc.TimeVal
var reread_wizlist int8
var emergency_unban int8
var logfile *stdio.File = nil
var text_overflow *byte = libc.CString("**OVERFLOW**\r\n")
var dg_act_check int
var pulse uint = 0
var fCopyOver bool
var port uint16
var last_act_message *byte = nil

func copyover_recover() {
	var (
		d              *descriptor_data
		fp             *stdio.File
		host           [1024]byte
		desc           int
		player_i       int
		fOld           bool
		name           [2048]byte
		username       [100]byte
		saved_loadroom int = int(-1)
		set_loadroom   int = int(-1)
	)
	basic_mud_log(libc.CString("Copyover recovery initiated"))
	PCOUNTDAY = libc.GetTime(nil) + 60
	fp = stdio.FOpen(COPYOVER_FILE, "r")
	if fp == nil {
		perror(libc.CString("copyover_recover:fopen"))
		basic_mud_log(libc.CString("Copyover file not found. Exitting.\n\r"))
		os.Exit(1)
	}
	stdio.Unlink(libc.CString(COPYOVER_FILE))
	for {
		fOld = TRUE != 0
		stdio.Fscanf(fp, "%d %s %s %d %s\n", &desc, &name[0], &host[0], &saved_loadroom, &username[0])
		if desc == -1 {
			break
		}
		if write_to_descriptor(desc, libc.CString("\n\rFolding initiated...\n\r")) < 0 {
			stdio.ByFD(uintptr(desc)).Close()
			continue
		}
		d = new(descriptor_data)
		*(*descriptor_data)(unsafe.Pointer((*byte)(unsafe.Pointer(d)))) = descriptor_data{}
		init_descriptor(d, desc)
		libc.StrCpy(&d.Host[0], &host[0])
		d.Next = descriptor_list
		descriptor_list = d
		d.Connected = CON_CLOSE
		d.Character = new(char_data)
		clear_char(d.Character)
		d.Character.Player_specials = new(player_special_data)
		d.Character.Desc = d
		if (func() int {
			player_i = load_char(&name[0], d.Character)
			return player_i
		}()) >= 0 {
			d.Character.Pfilepos = player_i
			if !PLR_FLAGGED(d.Character, PLR_DELETED) {
				d.Character.Act[int(PLR_WRITING/32)] &= bitvector_t(int32(^(1 << (int(PLR_WRITING % 32)))))
				d.Character.Act[int(PLR_MAILING/32)] &= bitvector_t(int32(^(1 << (int(PLR_MAILING % 32)))))
				d.Character.Act[int(PLR_CRYO/32)] &= bitvector_t(int32(^(1 << (int(PLR_CRYO % 32)))))
				userLoad(d, &username[0])
			}
		} else {
			fOld = FALSE != 0
		}
		if !fOld {
			write_to_descriptor(desc, libc.CString("\n\rSomehow, your character was lost during the folding. Sorry.\n\r"))
			close_socket(d)
		} else {
			write_to_descriptor(desc, libc.CString("\n\rFolding complete.\n\r"))
			set_loadroom = int(d.Character.Player_specials.Load_room)
			d.Character.Player_specials.Load_room = room_vnum(saved_loadroom)
			enter_player_game(d)
			d.Character.Player_specials.Load_room = room_vnum(set_loadroom)
			d.Connected = CON_PLAYING
			look_at_room(d.Character.In_room, d.Character, 0)
			if AFF_FLAGGED(d.Character, AFF_HAYASA) {
				d.Character.Speedboost = int(float64(GET_SPEEDCALC(d.Character)) * 0.5)
			}
		}
	}
	fp.Close()
}
func init_game(cmport uint16) {
	touch(libc.CString(KILLSCRIPT_FILE))
	circle_srandom(uint(libc.GetTime(nil)))
	basic_mud_log(libc.CString("Finding player limit."))
	max_players = get_max_players()
	if !fCopyOver {
		basic_mud_log(libc.CString("Opening mother connection."))
	}
	event_init()
	init_lookup_table()
	boot_db()
	var mapfile *stdio.File
	var rowcounter int
	var colcounter int
	var vnum_read int
	basic_mud_log(libc.CString("Signal trapping."))
	signal_setup()
	basic_mud_log(libc.CString("Loading Space Map. "))
	mapfile = stdio.FOpen("../lib/surface.map", "r")
	for rowcounter = 0; rowcounter <= MAP_ROWS; rowcounter++ {
		for colcounter = 0; colcounter <= MAP_COLS; colcounter++ {
			stdio.Fscanf(mapfile, "%d", &vnum_read)
			mapnums[rowcounter][colcounter] = int(real_room(room_vnum(vnum_read)))
		}
	}
	mapfile.Close()
	topLoad()
	stdio.Remove(KILLSCRIPT_FILE)
	if fCopyOver {
		copyover_recover()
	}
	basic_mud_log(libc.CString("Entering game loop."))
	Crash_save_all()
	basic_mud_log(libc.CString("Closing all sockets."))
	for descriptor_list != nil {
		close_socket(descriptor_list)
	}
	if circle_reboot != 2 {
		save_all()
	}
	basic_mud_log(libc.CString("Saving current MUD time."))
	save_mud_time(&time_info)
	if circle_reboot != 0 {
		basic_mud_log(libc.CString("Rebooting."))
		os.Exit(52)
	}
	basic_mud_log(libc.CString("Normal termination of game."))
}
func init_socket(cmport uint16) int {
}
func get_max_players() int {
	return 1000
}
func game_loop(cmmother_desc int) {
}
func heartbeat(heart_pulse int) {
	var mins_since_crashsave int = 0
	event_process()
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 13)) == 0 {
		script_trigger_check()
	}
	if (heart_pulse % (config_info.Ticks.Pulse_zone * (int(1000000 / OPT_USEC)))) == 0 {
		zone_update()
	}
	if (heart_pulse % (config_info.Ticks.Pulse_idlepwd * (int(1000000 / OPT_USEC)))) == 0 {
		check_idle_passwords()
	}
	if (heart_pulse % (((int(1000000 / OPT_USEC)) * 1) * 60)) == 0 {
		check_idle_menu()
	}
	if (heart_pulse % ((config_info.Ticks.Pulse_idlepwd * (int(1000000 / OPT_USEC))) / 15)) == 0 {
		dball_load()
	}
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 2)) == 0 {
		base_update()
		fish_update()
	}
	if (heart_pulse % (((int(1000000 / OPT_USEC)) * 1) * 15)) == 0 {
		handle_songs()
	}
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 1)) == 0 {
		wishSYS()
	}
	if (heart_pulse % (config_info.Ticks.Pulse_mobile * (int(1000000 / OPT_USEC)))) == 0 {
		mobile_activity()
	}
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 15)) == 0 {
		check_auction()
	}
	if (heart_pulse % ((config_info.Ticks.Pulse_idlepwd * (int(1000000 / OPT_USEC))) / 15)) == 0 {
		fight_stack()
	}
	if (heart_pulse % (((config_info.Ticks.Pulse_idlepwd * (int(1000000 / OPT_USEC))) / 15) * 2)) == 0 {
		if rand_number(1, 2) == 2 {
			homing_update()
		}
		huge_update()
		broken_update()
	}
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 1)) == 0 {
		copyover_check()
	}
	if (heart_pulse % (config_info.Ticks.Pulse_violence * (int(1000000 / OPT_USEC)))) == 0 {
		affect_update_violence()
	}
	if (heart_pulse % (SECS_PER_MUD_HOUR * (int(1000000 / OPT_USEC)))) == 0 {
		weather_and_time(1)
		check_timed_triggers()
		affect_update()
	}
	if (heart_pulse % ((int(SECS_PER_MUD_HOUR / 3)) * (int(1000000 / OPT_USEC)))) == 0 {
		point_update()
	}
	if config_info.Csd.Auto_save != 0 && (heart_pulse%(config_info.Ticks.Pulse_autosave*(int(1000000/OPT_USEC)))) == 0 {
		clan_update()
		if func() int {
			p := &mins_since_crashsave
			*p++
			return *p
		}() >= config_info.Csd.Autosave_time {
			mins_since_crashsave = 0
			Crash_save_all()
			House_save_all()
		}
	}
	if (heart_pulse % (config_info.Ticks.Pulse_sanity * 300 * (int(1000000 / OPT_USEC)))) == 0 {
		record_usage()
	}
	if (heart_pulse % (config_info.Ticks.Pulse_timesave * 900 * (int(1000000 / OPT_USEC)))) == 0 {
		save_mud_time(&time_info)
	}
	if (heart_pulse % ((int(1000000 / OPT_USEC)) * 30)) == 0 {
		timed_dt(nil)
	}
	extract_pending_chars()
}
func timediff(rslt *libc.TimeVal, a *libc.TimeVal, b *libc.TimeVal) {
	if a.Sec < b.Sec {
		*rslt = null_time
	} else if a.Sec == b.Sec {
		if a.USec < b.USec {
			*rslt = null_time
		} else {
			rslt.Sec = 0
			rslt.USec = a.USec - b.USec
		}
	} else {
		rslt.Sec = a.Sec - b.Sec
		if a.USec < b.USec {
			rslt.USec = a.USec + 1000000 - b.USec
			rslt.Sec--
		} else {
			rslt.USec = a.USec - b.USec
		}
	}
}
func timeadd(rslt *libc.TimeVal, a *libc.TimeVal, b *libc.TimeVal) {
	rslt.Sec = a.Sec + b.Sec
	rslt.USec = a.USec + b.USec
	for rslt.USec >= 1000000 {
		rslt.USec -= 1000000
		rslt.Sec++
	}
}
func record_usage() {
	var (
		sockets_connected int = 0
		sockets_playing   int = 0
		d                 *descriptor_data
	)
	for d = descriptor_list; d != nil; d = d.Next {
		sockets_connected++
		if IS_PLAYING(d) {
			sockets_playing++
		}
	}
	basic_mud_log(libc.CString("nusage: %-3d sockets connected, %-3d sockets playing"), sockets_connected, sockets_playing)
}
func make_prompt(d *descriptor_data) *byte {
	var (
		prompt  [1024]byte
		chair   *obj_data = nil
		flagged int       = FALSE
	)
	if d.Showstr_count != 0 {
		stdio.Snprintf(&prompt[0], int(1024), "\r\n[ Return to continue, (q)uit, (r)efresh, (b)ack, or page number (%d/%d) ]", d.Showstr_page, d.Showstr_count)
	} else if d.Str != nil {
		if d.Connected == CON_EXDESC {
			libc.StrCpy(&prompt[0], libc.CString("Enter Description(/h for editor help)> "))
		} else if PLR_FLAGGED(d.Character, PLR_WRITING) && !PLR_FLAGGED(d.Character, PLR_MAILING) {
			libc.StrCpy(&prompt[0], libc.CString("Enter Message(/h for editor help)> "))
		} else if PLR_FLAGGED(d.Character, PLR_MAILING) {
			libc.StrCpy(&prompt[0], libc.CString("Enter Mail Message(/h for editor help)> "))
		} else {
			libc.StrCpy(&prompt[0], libc.CString("Enter Message> "))
		}
	} else if d.Connected == CON_PLAYING && !IS_NPC(d.Character) {
		var (
			count int
			len_  uint64 = 0
		)
		prompt[0] = '\x00'
		if int(d.Character.Player_specials.Invis_level) != 0 && len_ < uint64(1024) {
			count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "i%d ", d.Character.Player_specials.Invis_level)
			if count >= 0 {
				len_ += uint64(count)
			}
		}
		if PRF_FLAGGED(d.Character, PRF_DISPAUTO) && GET_LEVEL(d.Character) >= 500 && len_ < uint64(1024) {
			var ch *char_data = d.Character
			if ch.Hit<<2 < ch.Max_hit {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "PL: %lld ", ch.Hit)
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if ch.Move<<2 < ch.Max_move && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "STA: %lld ", ch.Move)
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if ch.Ki<<2 < ch.Max_ki && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "KI: %lld ", ch.Ki)
				if count >= 0 {
					len_ += uint64(count)
				}
			}
		} else {
			if len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@w")
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PLR_FLAGGED(d.Character, PLR_SELFD) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RSELF-D@r: @w%s@D]@n", func() string {
					if PLR_FLAGGED(d.Character, PLR_SELFD2) {
						return "READY"
					}
					return "PREP"
				}())
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if int(d.Character.Race) == RACE_HALFBREED && !PLR_FLAGGED(d.Character, PLR_FURY) && PRF_FLAGGED(d.Character, PRF_FURY) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mFury@W: @r%d@D]@w", d.Character.Fury)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if int(d.Character.Race) == RACE_HALFBREED && PLR_FLAGGED(d.Character, PLR_FURY) && PRF_FLAGGED(d.Character, PRF_FURY) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mFury@W: @rENGAGED@D]@w")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if has_mail(int(d.Character.Idnum)) != 0 && !PRF_FLAGGED(d.Character, PRF_NMWARN) && d.Character.Admlevel > 0 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "CHECK MAIL - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Kaioken > 0 && d.Character.Admlevel > 0 {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "KAIOKEN X%d - ", d.Character.Kaioken)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Powerattack > 0 {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "%s - ", song_types[d.Character.Powerattack])
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Kaioken > 0 && d.Character.Admlevel <= 0 {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "KAIOKEN X%d - ", d.Character.Kaioken)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if has_mail(int(d.Character.Idnum)) != 0 && d.Character.Admlevel <= 0 && !PRF_FLAGGED(d.Character, PRF_NMWARN) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "CHECK MAIL - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Snooping != nil && d.Snooping.Character != nil && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Snooping: (%s) - ", GET_NAME(d.Snooping.Character))
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Drag != nil && d.Character.Drag != nil && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Dragging: (%s) - ", GET_NAME(d.Character.Drag))
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_BUILDWALK) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "BUILDWALKING - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if AFF_FLAGGED(d.Character, AFF_FLYING) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "FLYING - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if AFF_FLAGGED(d.Character, AFF_HIDE) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "HIDING - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PLR_FLAGGED(d.Character, PLR_SPAR) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "SPARRING - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PLR_FLAGGED(d.Character, PLR_NOSHOUT) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "MUTED - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 51 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Bash) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 52 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Headbutt) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 56 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Tailwhip) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 0 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Punch) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 1 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Kick) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 2 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Elbow) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 3 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Knee) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 4 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Roundhouse) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 5 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Uppercut) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 6 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Slam) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Combo == 8 && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Combo (Heeldrop) - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_AFK) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "AFK - ")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PLR_FLAGGED(d.Character, PLR_FISHING) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "FISHING -")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if flagged == TRUE && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@n\n")
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Sits != nil && PLR_FLAGGED(d.Character, PLR_HEALT) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				chair = d.Character.Sits
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@c<@CFloating inside a healing tank@c>@n\r\n")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Sits != nil && int(d.Character.Position) == POS_SITTING && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				chair = d.Character.Sits
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Sitting on: %s\r\n", chair.Short_description)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Sits != nil && int(d.Character.Position) == POS_RESTING && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				chair = d.Character.Sits
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Resting on: %s\r\n", chair.Short_description)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if d.Character.Sits != nil && int(d.Character.Position) == POS_SLEEPING && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				chair = d.Character.Sits
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Sleeping on: %s\r\n", chair.Short_description)
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if AFF_FLAGGED(d.Character, AFF_POSITION) && len_ < uint64(1024) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
				chair = d.Character.Sits
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "(Best Position)\r\n")
				flagged = TRUE
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if float64(d.Character.Charge) < float64(d.Character.Max_mana)*0.01 && d.Character.Charge > 0 {
				d.Character.Charge = 0
			}
			if d.Character.Charge > 0 {
				var charge int64 = d.Character.Charge
				if !PRF_FLAGGED(d.Character, PRF_NODEC) && !PRF_FLAGGED(d.Character, PRF_DISPERC) {
					if charge >= d.Character.Max_mana {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G==@D<@RMAX@D>@G===@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.95 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=========-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.9 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=========@g-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.85 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G========-@g-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.8 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G========@g--@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.75 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=======-@g--@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.7 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=======@g---@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.65 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G======-@g---@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.6 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G======@g----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.55 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=====-@g----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.5 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=====@g-----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.45 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G====-@g-----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.4 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G====@g------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.35 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G===-@g------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.3 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G===@g-------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.25 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G==-@g-------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.2 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G==@g--------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.15 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=-@g--------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.1 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G=@g---------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) >= float64(d.Character.Max_mana)*0.05 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@G-@g---------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(charge) < float64(d.Character.Max_mana)*0.05 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@g----------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@CCharge @D[@g----------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
				if PRF_FLAGGED(d.Character, PRF_DISPERC) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
					if d.Character.Charge > 0 {
						var perc int64 = (d.Character.Charge * 100) / d.Character.Max_mana
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@BCharge@Y: @C%lld%s@D]@n\n", perc, "%")
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
				if PRF_FLAGGED(d.Character, PRF_NODEC) {
					if charge > 0 {
						var perc int64 = (charge * 100) / d.Character.Max_mana
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "Ki is charged to %lld percent.\n", perc)
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
			}
			if AFF_FLAGGED(d.Character, AFF_FIRESHIELD) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D(@rF@RI@YR@rE@RS@YH@rI@RE@YL@rD@D)@n\n")
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if AFF_FLAGGED(d.Character, AFF_SANCTUARY) {
				if PRF_FLAGGED(d.Character, PRF_DISPERC) && !PRF_FLAGGED(d.Character, PRF_NODEC) {
					if d.Character.Barrier > 0 {
						var perc int64 = (d.Character.Barrier * 100) / d.Character.Max_mana
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GBarrier@Y: @B%lld%s@D]@n\n", perc, "%")
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
				if !PRF_FLAGGED(d.Character, PRF_NODEC) && !PRF_FLAGGED(d.Character, PRF_DISPERC) {
					if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.75 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C==MAX==@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.7 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C=======@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.65 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C======-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.6 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C======@c-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.55 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C=====-@c-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.5 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C=====@c--@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.45 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C====-@c--@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.4 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C====@c---@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.35 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C===-@c---@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.3 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C===@c----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.25 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C==-@c----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.2 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C==@c-----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.15 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C=-@c-----@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.1 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C=@c------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) >= float64(d.Character.Max_mana)*0.05 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C-@c------@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					} else if float64(d.Character.Barrier) < float64(d.Character.Max_mana)*0.05 {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@BBarrier @D[@C--Low-@D]@n\n")
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
				if PRF_FLAGGED(d.Character, PRF_NODEC) {
					if d.Character.Barrier > 0 {
						var perc int64 = (d.Character.Barrier * 100) / d.Character.Max_mana
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "A barrier charged to %lld percent surrounds you.@n\n", perc)
						if count >= 0 {
							len_ += uint64(count)
						}
					}
				}
			}
			if !PRF_FLAGGED(d.Character, PRF_DISPERC) {
				if PRF_FLAGGED(d.Character, PRF_DISPHP) && len_ < uint64(1024) && d.Character.Hit >= gear_pl(d.Character) && d.Character.Hit < d.Character.Max_hit {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @m%s@D]@n", add_commas(d.Character.Hit))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPHP) && len_ < uint64(1024) && d.Character.Hit > d.Character.Max_hit {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @G%s@D]@n", add_commas(d.Character.Hit))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPHP) && len_ < uint64(1024) && d.Character.Hit > gear_pl(d.Character)/2 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @c%s@D]@n", add_commas(d.Character.Hit))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPHP) && len_ < uint64(1024) && d.Character.Hit > gear_pl(d.Character)/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @y%s@D]@n", add_commas(d.Character.Hit))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPHP) && len_ < uint64(1024) && d.Character.Hit <= gear_pl(d.Character)/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @r%s@D]@n", add_commas(d.Character.Hit))
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			} else if PRF_FLAGGED(d.Character, PRF_DISPHP) {
				var (
					power    int64 = d.Character.Hit
					maxpower int64 = d.Character.Max_hit
					perc     int   = 0
				)
				if power <= 0 {
					power = 1
				}
				if maxpower <= 0 {
					maxpower = 1
				}
				perc = int((power * 100) / maxpower)
				if perc > 100 {
					if power >= gear_pl(d.Character) && power < d.Character.Max_hit {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @m%d%s@D]@n", perc, "@w%")
					} else if power > d.Character.Max_hit {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @G%d%s@D]@n", perc, "@w%")
					} else {
						count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @g%d%s@D]@n", perc, "@w%")
					}
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 70 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @c%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 51 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @Y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@RPL@n@Y: @r%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			}
			if !PRF_FLAGGED(d.Character, PRF_DISPERC) {
				if PRF_FLAGGED(d.Character, PRF_DISPKI) && len_ < uint64(1024) && d.Character.Mana > d.Character.Max_mana/2 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@Y: @c%s@D]@n", add_commas(d.Character.Mana))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPKI) && len_ < uint64(1024) && d.Character.Mana > d.Character.Max_mana/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@Y: @y%s@D]@n", add_commas(d.Character.Mana))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPKI) && len_ < uint64(1024) && d.Character.Mana <= d.Character.Max_mana/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@Y: @r%s@D]@n", add_commas(d.Character.Mana))
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			} else if PRF_FLAGGED(d.Character, PRF_DISPKI) {
				var (
					power    int64 = d.Character.Mana
					maxpower int64 = d.Character.Max_mana
					perc     int   = 0
				)
				if power <= 0 {
					power = 1
				}
				if maxpower <= 0 {
					maxpower = 1
				}
				perc = int((power * 100) / maxpower)
				if perc > 100 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@n@Y: @G%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 70 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@n@Y: @c%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 51 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@n@Y: @Y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@n@Y: @y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@CKI@n@Y: @r%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			}
			if !PRF_FLAGGED(d.Character, PRF_DISPERC) {
				if PRF_FLAGGED(d.Character, PRF_DISPMOVE) && len_ < uint64(1024) && d.Character.Move > d.Character.Max_move/2 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@Y: @c%s@D]@n", add_commas(d.Character.Move))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPMOVE) && len_ < uint64(1024) && d.Character.Move > d.Character.Max_move/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@Y: @y%s@D]@n", add_commas(d.Character.Move))
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if PRF_FLAGGED(d.Character, PRF_DISPMOVE) && len_ < uint64(1024) && d.Character.Move <= d.Character.Max_move/10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@Y: @r%s@D]@n", add_commas(d.Character.Move))
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			} else if PRF_FLAGGED(d.Character, PRF_DISPMOVE) {
				var (
					power    int64 = d.Character.Move
					maxpower int64 = d.Character.Max_move
					perc     int   = 0
				)
				if power <= 0 {
					power = 1
				}
				if maxpower <= 0 {
					maxpower = 1
				}
				perc = int((power * 100) / maxpower)
				if perc > 100 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@n@Y: @G%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 70 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@n@Y: @c%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 51 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@n@Y: @Y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else if perc >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@n@Y: @y%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				} else {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@GSTA@n@Y: @r%d%s@D]@n", perc, "@w%")
					if count >= 0 {
						len_ += uint64(count)
					}
				}
			}
			if PRF_FLAGGED(d.Character, PRF_DISPTNL) && len_ < uint64(1024) && GET_LEVEL(d.Character) < 100 {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@yTNL@Y: @W%s@D]@n", add_commas(int64(level_exp(d.Character, GET_LEVEL(d.Character)+1)-int(d.Character.Exp))))
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_DISTIME) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@W%2d %s@D]@n", func() int {
					if time_info.Hours%12 == 0 {
						return 12
					}
					return time_info.Hours % 12
				}(), func() string {
					if time_info.Hours >= 12 {
						return "PM"
					}
					return "AM"
				}())
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_DISGOLD) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@YZen@y: @W%s@D]@n", add_commas(int64(d.Character.Gold)))
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_DISPRAC) && len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mPS@y: @W%s@D]@n", add_commas(int64(d.Character.Player_specials.Class_skill_points[d.Character.Chclass])))
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if PRF_FLAGGED(d.Character, PRF_DISHUTH) && len_ < uint64(1024) {
				var (
					hun  int = int(d.Character.Player_specials.Conditions[HUNGER])
					thir int = int(d.Character.Player_specials.Conditions[THIRST])
				)
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "\n@D[@mHung@y:")
				if count >= 0 {
					len_ += uint64(count)
				}
				if hun >= 48 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WFull@D]@n")
				} else if hun >= 40 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WAlmost Full@D]@n")
				} else if hun >= 30 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WNeed Snack@D]@n")
				} else if hun >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WHungry@D]@n")
				} else if hun >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WVery Hungry@D]@n")
				} else if hun >= 10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WAlmost Starving@D]@n")
				} else if hun >= 5 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WNear Starving@D]@n")
				} else if hun >= 0 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WStarving@D]@n")
				} else {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), " @WN/A@D]@n")
				}
				if count >= 0 {
					len_ += uint64(count)
				}
				if thir >= 48 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WQuenched@D]@n")
				} else if thir >= 40 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WNeed Sip@D]@n")
				} else if thir >= 30 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WNeed Drink@D]@n")
				} else if thir >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WThirsty@D]@n")
				} else if thir >= 20 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WVery Thirsty@D]@n")
				} else if thir >= 10 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WAlmost Dehydrated@D]@n")
				} else if thir >= 5 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WNear Dehydration@D]@n")
				} else if thir >= 0 {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WDehydrated@D]@n")
				} else {
					count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "@D[@mThir@y: @WN/A@D]@n")
				}
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if len_ < uint64(1024) && has_group(d.Character) != 0 && !PRF_FLAGGED(d.Character, PRF_GHEALTH) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "\n%s", report_party_health(d.Character))
				if d.Character.Temp_prompt != nil {
					libc.Free(unsafe.Pointer(d.Character.Temp_prompt))
				}
				if count >= 0 {
					len_ += uint64(count)
				}
			}
			if len_ < uint64(1024) {
				count = stdio.Snprintf(&prompt[len_], int(1024-uintptr(len_)), "\n")
			}
		}
		if len_ < uint64(1024) && len_ < 5 {
			libc.StrNCat(&prompt[0], libc.CString(">\n"), int(1024-uintptr(len_)-1))
		}
	} else if d.Connected == CON_PLAYING && IS_NPC(d.Character) {
		stdio.Snprintf(&prompt[0], int(1024), "%s>\n", CAP(GET_NAME(d.Character)))
	} else {
		prompt[0] = '\x00'
	}
	return &prompt[0]
}
func write_to_q(txt *byte, queue *txt_q, aliased int) {
	var newt *txt_block
	newt = new(txt_block)
	newt.Text = libc.StrDup(txt)
	newt.Aliased = aliased
	if queue.Head == nil {
		newt.Next = nil
		queue.Head = func() *txt_block {
			p := &queue.Tail
			queue.Tail = newt
			return *p
		}()
	} else {
		queue.Tail.Next = newt
		queue.Tail = newt
		newt.Next = nil
	}
}
func get_from_q(queue *txt_q, dest *byte, aliased *int) int {
	var tmp *txt_block
	if queue.Head == nil {
		return 0
	}
	libc.StrCpy(dest, queue.Head.Text)
	*aliased = queue.Head.Aliased
	tmp = queue.Head
	queue.Head = queue.Head.Next
	libc.Free(unsafe.Pointer(tmp.Text))
	libc.Free(unsafe.Pointer(tmp))
	return 1
}
func flush_queues(d *descriptor_data) {
	if d.Large_outbuf != nil {
		d.Large_outbuf.Next = bufpool
		bufpool = d.Large_outbuf
	}
	for d.Input.Head != nil {
		var tmp *txt_block = d.Input.Head
		d.Input.Head = d.Input.Head.Next
		libc.Free(unsafe.Pointer(tmp.Text))
		libc.Free(unsafe.Pointer(tmp))
	}
}
func write_to_output(t *descriptor_data, txt *byte, _rest ...interface{}) uint64 {
	var (
		args libc.ArgList
		left uint64
	)
	args.Start(txt, _rest)
	left = vwrite_to_output(t, txt, args)
	args.End()
	return left
}

var ANSI [31]*byte = [31]*byte{libc.CString("@"), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_NORMAL), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AA_BOLD), libc.CString(AB_BLACK), libc.CString(AB_BLUE), libc.CString(AB_GREEN), libc.CString(AB_CYAN), libc.CString(AB_RED), libc.CString(AB_MAGENTA), libc.CString(AB_YELLOW), libc.CString(AB_WHITE), libc.CString(AA_BLINK), libc.CString(AA_UNDERLINE), libc.CString(AA_BOLD), libc.CString(AA_REVERSE), libc.CString("!")}
var CCODE [33]byte = func() [33]byte {
	var t [33]byte
	copy(t[:], []byte("@ndbgcrmywDBGCRMYW01234567luoex!"))
	return t
}()
var RANDOM_COLORS [15]byte = func() [15]byte {
	var t [15]byte
	copy(t[:], []byte("bgcrmywBGCRMWY"))
	return t
}()

func proc_colors(txt *byte, maxlen uint64, parse int, choices **byte) uint64 {
	var (
		dest_char   *byte
		source_char *byte
		color_char  *byte
		save_pos    *byte
		replacement *byte = nil
		i           int
		temp_color  int
		wanted      uint64
	)
	if txt == nil || libc.StrChr(txt, '@') == nil {
		return uint64(libc.StrLen(txt))
	}
	source_char = txt
	dest_char = (*byte)(unsafe.Pointer(&make([]int8, int(maxlen))[0]))
	save_pos = dest_char
	for *source_char != 0 && uint64(int64(uintptr(unsafe.Pointer(dest_char))-uintptr(unsafe.Pointer(save_pos)))) < maxlen {
		if *source_char != '@' {
			*func() *byte {
				p := &dest_char
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &source_char
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
			continue
		}
		source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
		if *source_char == 'x' {
			temp_color = int(libc.Rand()) % 14
			*source_char = RANDOM_COLORS[temp_color]
		}
		if *source_char == '\x00' {
			*func() *byte {
				p := &dest_char
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = '@'
			continue
		}
		if parse == 0 {
			if *source_char == '@' {
				*func() *byte {
					p := &dest_char
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}() = '@'
			}
			if *source_char == '[' {
				source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
				for *source_char != 0 && unicode.IsDigit(rune(*source_char)) {
					source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
				}
				if *source_char == 0 {
					source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), -1))
				}
			}
			source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
			continue
		}
		if *source_char == '[' {
			source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
			if *source_char != 0 {
				i = libc.Atoi(libc.GoString(source_char))
				if i < 0 || i >= NUM_COLOR {
					i = COLOR_NORMAL
				}
				replacement = default_color_choices[i]
				if choices != nil && *(**byte)(unsafe.Add(unsafe.Pointer(choices), unsafe.Sizeof((*byte)(nil))*uintptr(i))) != nil {
					replacement = *(**byte)(unsafe.Add(unsafe.Pointer(choices), unsafe.Sizeof((*byte)(nil))*uintptr(i)))
				}
				for *source_char != 0 && unicode.IsDigit(rune(*source_char)) {
					source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
				}
				if *source_char == 0 {
					source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), -1))
				}
			}
		} else if *source_char == 'n' {
			replacement = default_color_choices[COLOR_NORMAL]
			if choices != nil && *(**byte)(unsafe.Add(unsafe.Pointer(choices), unsafe.Sizeof((*byte)(nil))*uintptr(COLOR_NORMAL))) != nil {
				replacement = *(**byte)(unsafe.Add(unsafe.Pointer(choices), unsafe.Sizeof((*byte)(nil))*uintptr(COLOR_NORMAL)))
			}
		} else {
			for i = 0; CCODE[i] != '!'; i++ {
				if (*source_char) == CCODE[i] {
					replacement = ANSI[i]
					break
				}
			}
		}
		if replacement != nil {
			if uint64(int64(uintptr(unsafe.Pointer(dest_char))-uintptr(unsafe.Pointer(save_pos))))+uint64(libc.StrLen(replacement))+uint64(libc.StrLen(libc.CString(ANSISTART)))+1 < maxlen {
				if unicode.IsDigit(rune(*replacement)) {
					for color_char = libc.CString(ANSISTART); *color_char != 0; {
						*func() *byte {
							p := &dest_char
							x := *p
							*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
							return x
						}() = *func() *byte {
							p := &color_char
							x := *p
							*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
							return x
						}()
					}
				}
				for color_char = replacement; *color_char != 0; {
					*func() *byte {
						p := &dest_char
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}() = *func() *byte {
						p := &color_char
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}()
				}
				if unicode.IsDigit(rune(*replacement)) {
					*func() *byte {
						p := &dest_char
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}() = ANSIEND
				}
			}
			replacement = nil
		}
		source_char = (*byte)(unsafe.Add(unsafe.Pointer(source_char), 1))
	}
	*dest_char = '\x00'
	wanted = uint64(libc.StrLen(source_char))
	libc.StrNCpy(txt, save_pos, int(maxlen-1))
	libc.Free(unsafe.Pointer(save_pos))
	return uint64(int64(uintptr(unsafe.Pointer(dest_char))-uintptr(unsafe.Pointer(save_pos)))) + wanted
}
func vwrite_to_output(t *descriptor_data, format *byte, args libc.ArgList) uint64 {
	var (
		txt      [64936]byte
		wantsize uint64
		size     int
	)
	if t.Bufspace == 0 {
		return 0
	}
	wantsize = uint64(func() int {
		size = stdio.Vsnprintf(&txt[0], int(64936), libc.GoString(format), args)
		return size
	}())
	if t.Character != nil {
		wantsize = uint64(func() int {
			size = int(proc_colors(&txt[0], uint64(64936), int(libc.BoolToInt((func() int {
				if !IS_NPC(t.Character) {
					if PRF_FLAGGED(t.Character, PRF_COLOR) {
						return 1
					}
					return 0
				}
				return 0
			}()) > 0)), &(func() [16]*byte {
				if IS_NPC(t.Character) {
					return ([16]*byte)(0)
				}
				if t.Character.Player_specials != nil {
					return t.Character.Player_specials.Color_choices
				}
				return ([16]*byte)(0)
			}())[0]))
			return size
		}())
	}
	if size < 0 || wantsize >= uint64(64936) {
		size = int(64936 - 1)
		libc.StrCpy((*byte)(unsafe.Add(unsafe.Pointer(&txt[size]), -libc.StrLen(text_overflow))), text_overflow)
	}
	if size+t.Bufptr+1 > (int((96*1024)-GARBAGE_SPACE) - MAX_PROMPT_LENGTH) {
		size = (int((96*1024)-GARBAGE_SPACE) - MAX_PROMPT_LENGTH) - t.Bufptr - 1
		txt[size] = '\x00'
		t.Character.Overf = TRUE
		buf_overflows++
	}
	if t.Bufspace > size {
		libc.StrCpy((*byte)(unsafe.Add(unsafe.Pointer(t.Output), t.Bufptr)), &txt[0])
		t.Bufspace -= size
		t.Bufptr += size
		return uint64(t.Bufspace)
	}
	buf_switches++
	if bufpool != nil {
		t.Large_outbuf = bufpool
		bufpool = bufpool.Next
	} else {
		t.Large_outbuf = new(txt_block)
		t.Large_outbuf.Text = (*byte)(unsafe.Pointer(&make([]int8, int((96*1024)-GARBAGE_SPACE)-MAX_PROMPT_LENGTH)[0]))
		buf_largecount++
	}
	libc.StrCpy(t.Large_outbuf.Text, t.Output)
	t.Output = t.Large_outbuf.Text
	libc.StrCat(t.Output, &txt[0])
	t.Bufptr = libc.StrLen(t.Output)
	t.Bufspace = (int((96*1024)-GARBAGE_SPACE) - MAX_PROMPT_LENGTH) - 1 - t.Bufptr
	return uint64(t.Bufspace)
}
func free_bufpool() {
	var tmp *txt_block
	for bufpool != nil {
		tmp = bufpool.Next
		if bufpool.Text != nil {
			libc.Free(unsafe.Pointer(bufpool.Text))
		}
		libc.Free(unsafe.Pointer(bufpool))
		bufpool = tmp
	}
}
func get_bind_addr() *cnet.Address {
}
func set_sendbuf(s int) int {
}
func init_descriptor(newd *descriptor_data, desc int) {
	var last_desc int = 0
	newd.Descriptor = desc
	newd.Idle_tics = 0
	newd.Output = &newd.Small_outbuf[0]
	newd.Bufspace = int(SMALL_BUFSIZE - 1)
	newd.Login_time = libc.GetTime(nil)
	*newd.Output = '\x00'
	newd.Bufptr = 0
	newd.Has_prompt = 1
	newd.Connected = CON_GET_USER
	newd.History = &make([]*byte, HISTORY_SIZE)[0]
	if func() int {
		p := &last_desc
		*p++
		return *p
	}() == 1000 {
		last_desc = 1
	}
	newd.Desc_num = last_desc
}
func set_color(d *descriptor_data) {
	if d.Character == nil {
		d.Character = new(char_data)
		clear_char(d.Character)
		d.Character.Player_specials = new(player_special_data)
		d.Character.Desc = d
	}
	d.Character.Player_specials.Pref[int(PRF_COLOR/32)] |= bitvector_t(int32(1 << (int(PRF_COLOR % 32))))
	write_to_output(d, GREETANSI)
	write_to_output(d, libc.CString("\r\n@w                  Welcome to Dragonball Advent Truth\r\n"))
	write_to_output(d, libc.CString("@D                 ---(@CPeak Logon Count Today@W: @w%4d@D)---@n\r\n"), PCOUNT)
	write_to_output(d, libc.CString("@D                 ---(@CHighest Logon Count   @W: @w%4d@D)---@n\r\n"), HIGHPCOUNT)
	write_to_output(d, libc.CString("@D                 ---(@CTotal Era %d Characters@W: @w%4s@D)---@n\r\n"), CURRENT_ERA, add_commas(int64(ERAPLAYERS)))
	write_to_output(d, libc.CString("\r\n@cEnter your desired username or the username you have already made.\n@CEnter Username:@n\r\n"))
	d.User = libc.CString("Empty")
	d.Pass = libc.CString("Empty")
	d.Email = libc.CString("Empty")
	d.Tmp1 = libc.CString("Empty")
	d.Tmp2 = libc.CString("Empty")
	d.Tmp3 = libc.CString("Empty")
	d.Tmp4 = libc.CString("Empty")
	d.Tmp5 = libc.CString("Empty")
	return
}
func new_descriptor(s int) int {
}
func process_output(t *descriptor_data) int {
	var (
		i      [98304]byte
		osb    *byte = &i[2]
		result int
	)
	libc.StrCpy(&i[0], libc.CString("\r\n"))
	libc.StrCpy(osb, t.Output)
	if t.Bufspace == 0 {
		libc.StrCat(osb, libc.CString("**OVERFLOW**\r\n"))
	}
	if t.Connected == CON_PLAYING && t.Character != nil && !IS_NPC(t.Character) && !PRF_FLAGGED(t.Character, PRF_COMPACT) {
		libc.StrCat(osb, libc.CString("\r\n"))
	}
	libc.StrCat(&i[0], make_prompt(t))
	if t.Connected == CON_PLAYING {
		proc_colors(&i[0], uint64(98304), int(libc.BoolToInt((func() int {
			if !IS_NPC(t.Character) {
				if PRF_FLAGGED(t.Character, PRF_COLOR) {
					return 1
				}
				return 0
			}
			return 0
		}()) > 0)), &(func() [16]*byte {
			if IS_NPC(t.Character) {
				return ([16]*byte)(0)
			}
			if t.Character.Player_specials != nil {
				return t.Character.Player_specials.Color_choices
			}
			return ([16]*byte)(0)
		}())[0])
	}
	if t.Has_prompt != 0 {
		t.Has_prompt = FALSE
		result = write_to_descriptor(t.Descriptor, &i[0])
		if result >= 2 {
			result -= 2
		}
	} else {
		result = write_to_descriptor(t.Descriptor, osb)
	}
	if result < 0 {
		close_socket(t)
		return -1
	} else if result == 0 {
		return 0
	}
	if t.Snoop_by != nil {
		write_to_output(t.Snoop_by, libc.CString("\nvvvvvvvvvvvvv[Snoop]vvvvvvvvvvvvv\n%s\n^^^^^^^^^^^^^[Snoop]^^^^^^^^^^^^^\n"), t.Output)
	}
	if result >= t.Bufptr {
		if t.Large_outbuf != nil {
			t.Large_outbuf.Next = bufpool
			bufpool = t.Large_outbuf
			t.Large_outbuf = nil
			t.Output = &t.Small_outbuf[0]
		}
		t.Bufspace = int(SMALL_BUFSIZE - 1)
		t.Bufptr = 0
		*t.Output = '\x00'
		if uint(result) < uint(libc.StrLen(osb)) {
			var savetextlen uint64 = uint64(libc.StrLen((*byte)(unsafe.Add(unsafe.Pointer(osb), result))))
			libc.StrCat(t.Output, (*byte)(unsafe.Add(unsafe.Pointer(osb), result)))
			t.Bufptr -= int(savetextlen)
			t.Bufspace += int(savetextlen)
		}
	} else {
		libc.StrCpy(t.Output, (*byte)(unsafe.Add(unsafe.Pointer(t.Output), result)))
		t.Bufptr -= result
		t.Bufspace += result
	}
	return result
}
func perform_socket_write(desc int, txt *byte, length uint64) int64 {
	var (
		result       int64 = 0
		compr_result int
	)
	_ = compr_result
	var tmp int
	_ = tmp
	var cnt int
	_ = cnt
	var bytes_copied int
	_ = bytes_copied
	result = int64(stdio.ByFD(uintptr(desc)).Write(txt, int(length)))
	if result > 0 {
		return result
	}
	if result == 0 {
		basic_mud_log(libc.CString("SYSERR: Huh??  write() returned 0???  Please report this!"))
		return -1
	}
	if libc.Errno == EAGAIN {
		return 0
	}
	return -1
}
func write_to_descriptor(desc int, txt *byte) int {
	var (
		bytes_written int64
		total         uint64 = uint64(libc.StrLen(txt))
		write_total   uint64 = 0
	)
	for total > 0 {
		bytes_written = perform_socket_write(desc, txt, total)
		if bytes_written < 0 {
			perror(libc.CString("SYSERR: Write to socket"))
			return -1
		} else if bytes_written == 0 {
			return int(write_total)
		} else {
			txt = (*byte)(unsafe.Add(unsafe.Pointer(txt), bytes_written))
			total -= uint64(bytes_written)
			write_total += uint64(bytes_written)
		}
	}
	return int(write_total)
}
func perform_socket_read(desc int, read_point *byte, space_left uint64) int64 {
	var ret int64
	ret = int64(stdio.ByFD(uintptr(desc)).Read(read_point, int(space_left)))
	if ret > 0 {
		return ret
	}
	if ret == 0 {
		basic_mud_log(libc.CString("WARNING: EOF on socket read (connection broken by peer)"))
		return -1
	}
	if libc.Errno == EINTR {
		return 0
	}
	if libc.Errno == EAGAIN {
		return 0
	}
	if libc.Errno == ECONNRESET {
		return -1
	}
	perror(libc.CString("SYSERR: perform_socket_read: about to lose connection"))
	return -1
}
func process_input(t *descriptor_data) int {
	var (
		buf_length   int
		failed_subst int
		bytes_read   int64
		space_left   uint64
		ptr          *byte
		read_point   *byte
		write_point  *byte
		nl_pos       *byte = nil
		tmp          [2048]byte
	)
	buf_length = libc.StrLen(&t.Inbuf[0])
	read_point = &t.Inbuf[buf_length]
	space_left = uint64(MAX_RAW_INPUT_LENGTH - buf_length - 1)
	for {
		if space_left <= 0 {
			basic_mud_log(libc.CString("WARNING: process_input: about to close connection: input overflow"))
			return -1
		}
		bytes_read = perform_socket_read(t.Descriptor, read_point, space_left)
		if bytes_read < 0 {
			return -1
		} else if bytes_read == 0 {
			return 0
		}
		*((*byte)(unsafe.Add(unsafe.Pointer(read_point), bytes_read))) = '\x00'
		for ptr = read_point; *ptr != 0 && nl_pos == nil; ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
			if ISNEWL(int8(*ptr)) {
				nl_pos = ptr
			}
		}
		read_point = (*byte)(unsafe.Add(unsafe.Pointer(read_point), bytes_read))
		space_left -= uint64(bytes_read)
		if nl_pos != nil {
			break
		}
	}
	read_point = &t.Inbuf[0]
	for nl_pos != nil {
		write_point = &tmp[0]
		space_left = uint64(int(MAX_INPUT_LENGTH - 1))
		for ptr = read_point; space_left > 1 && uintptr(unsafe.Pointer(ptr)) < uintptr(unsafe.Pointer(nl_pos)); ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
			if *ptr == '\b' || *ptr == math.MaxInt8 {
				if uintptr(unsafe.Pointer(write_point)) > uintptr(unsafe.Pointer(&tmp[0])) {
					if *(func() *byte {
						p := &write_point
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), -1))
						return *p
					}()) == '$' {
						write_point = (*byte)(unsafe.Add(unsafe.Pointer(write_point), -1))
						space_left += 2
					} else {
						space_left++
					}
				}
			} else if isascii(rune(*ptr)) && unicode.IsPrint(rune(*ptr)) {
				if (func() byte {
					p := (func() *byte {
						p := &write_point
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}())
					*(func() *byte {
						p := &write_point
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}()) = *ptr
					return *p
				}()) == '$' {
					*(func() *byte {
						p := &write_point
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}()) = '$'
					space_left -= 2
				} else {
					space_left--
				}
			}
		}
		*write_point = '\x00'
		if space_left <= 0 && uintptr(unsafe.Pointer(ptr)) < uintptr(unsafe.Pointer(nl_pos)) {
			var buffer [2112]byte
			stdio.Snprintf(&buffer[0], int(2112), "Line too long.  Truncated to:\r\n%s\r\n", &tmp[0])
			if write_to_descriptor(t.Descriptor, &buffer[0]) < 0 {
				return -1
			}
		}
		if t.Snoop_by != nil {
			write_to_output(t.Snoop_by, libc.CString("%% %s\r\n"), &tmp[0])
		}
		failed_subst = 0
		if tmp[0] == '!' && (tmp[1]) == 0 {
			libc.StrCpy(&tmp[0], &t.Last_input[0])
		} else if tmp[0] == '!' && tmp[1] != 0 {
			var (
				commandln    *byte = (&tmp[1])
				starting_pos int   = t.History_pos
				cnt          int   = (func() int {
					if t.History_pos == 0 {
						return int(HISTORY_SIZE - 1)
					}
					return t.History_pos - 1
				}())
			)
			skip_spaces(&commandln)
			for ; cnt != starting_pos; cnt-- {
				if *(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(cnt))) != nil && is_abbrev(commandln, *(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(cnt)))) != 0 {
					libc.StrCpy(&tmp[0], *(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(cnt))))
					libc.StrCpy(&t.Last_input[0], &tmp[0])
					write_to_output(t, libc.CString("%s\r\n"), &tmp[0])
					break
				}
				if cnt == 0 {
					cnt = HISTORY_SIZE
				}
			}
		} else if tmp[0] == '^' {
			if (func() int {
				failed_subst = perform_subst(t, &t.Last_input[0], &tmp[0])
				return failed_subst
			}()) == 0 {
				libc.StrCpy(&t.Last_input[0], &tmp[0])
			}
		} else {
			libc.StrCpy(&t.Last_input[0], &tmp[0])
			if *(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(t.History_pos))) != nil {
				libc.Free(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(t.History_pos)))))
			}
			*(**byte)(unsafe.Add(unsafe.Pointer(t.History), unsafe.Sizeof((*byte)(nil))*uintptr(t.History_pos))) = libc.StrDup(&tmp[0])
			if func() int {
				p := &t.History_pos
				*p++
				return *p
			}() >= HISTORY_SIZE {
				t.History_pos = 0
			}
		}
		if masadv(&tmp[0], t.Character) != 0 {
		}
		if tmp[0] == '-' && tmp[1] == '-' && (tmp[2]) == 0 {
			write_to_output(t, libc.CString("All queued commands cancelled.\r\n"))
			flush_queues(t)
		}
		if failed_subst == 0 {
			write_to_q(&tmp[0], &t.Input, 0)
		}
		for ISNEWL(int8(*nl_pos)) {
			nl_pos = (*byte)(unsafe.Add(unsafe.Pointer(nl_pos), 1))
		}
		read_point = func() *byte {
			ptr = nl_pos
			return ptr
		}()
		for nl_pos = nil; *ptr != 0 && nl_pos == nil; ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1)) {
			if ISNEWL(int8(*ptr)) {
				nl_pos = ptr
			}
		}
	}
	write_point = &t.Inbuf[0]
	for *read_point != 0 {
		*(func() *byte {
			p := &write_point
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(func() *byte {
			p := &read_point
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	*write_point = '\x00'
	return 1
}
func perform_subst(t *descriptor_data, orig *byte, subst *byte) int {
	var (
		newsub [2053]byte
		first  *byte
		second *byte
		strpos *byte
	)
	first = (*byte)(unsafe.Add(unsafe.Pointer(subst), 1))
	if (func() *byte {
		second = libc.StrChr(first, '^')
		return second
	}()) == nil {
		write_to_output(t, libc.CString("Invalid substitution.\r\n"))
		return 1
	}
	*(func() *byte {
		p := &second
		x := *p
		*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
		return x
	}()) = '\x00'
	if (func() *byte {
		strpos = libc.StrStr(orig, first)
		return strpos
	}()) == nil {
		write_to_output(t, libc.CString("Invalid substitution.\r\n"))
		return 1
	}
	libc.StrNCpy(&newsub[0], orig, int(int64(uintptr(unsafe.Pointer(strpos))-uintptr(unsafe.Pointer(orig)))))
	newsub[int64(uintptr(unsafe.Pointer(strpos))-uintptr(unsafe.Pointer(orig)))] = '\x00'
	libc.StrNCat(&newsub[0], second, MAX_INPUT_LENGTH-libc.StrLen(&newsub[0])-1)
	if ((int64(uintptr(unsafe.Pointer(strpos)) - uintptr(unsafe.Pointer(orig)))) + int64(libc.StrLen(first))) < int64(libc.StrLen(orig)) {
		libc.StrNCat(&newsub[0], (*byte)(unsafe.Add(unsafe.Pointer(strpos), libc.StrLen(first))), MAX_INPUT_LENGTH-libc.StrLen(&newsub[0])-1)
	}
	newsub[int(MAX_INPUT_LENGTH-1)] = '\x00'
	libc.StrCpy(subst, &newsub[0])
	return 0
}
func free_user(d *descriptor_data) {
	if d.User_freed == 1 {
		return
	}
	if d.User == nil {
		send_to_imm(libc.CString("ERROR: free_user called but no user to free!"))
		return
	}
	d.User_freed = 1
	if libc.StrCaseCmp(d.User, libc.CString("Empty")) == 0 {
		return
	}
	basic_mud_log(libc.CString("Freeing User: %s"), d.User)
	if d.User != nil {
		libc.Free(unsafe.Pointer(d.User))
	}
	if d.Pass != nil {
		libc.Free(unsafe.Pointer(d.Pass))
	}
	if d.Email != nil {
		libc.Free(unsafe.Pointer(d.Email))
	}
	if d.Tmp1 != nil {
		libc.Free(unsafe.Pointer(d.Tmp1))
	}
	if d.Tmp2 != nil {
		libc.Free(unsafe.Pointer(d.Tmp2))
	}
	if d.Tmp3 != nil {
		libc.Free(unsafe.Pointer(d.Tmp3))
	}
	if d.Tmp4 != nil {
		libc.Free(unsafe.Pointer(d.Tmp4))
	}
	if d.Tmp5 != nil {
		libc.Free(unsafe.Pointer(d.Tmp5))
	}
}
func close_socket(d *descriptor_data) {
	var temp *descriptor_data
	if d == descriptor_list {
		descriptor_list = d.Next
	} else {
		temp = descriptor_list
		for temp != nil && temp.Next != d {
			temp = temp.Next
		}
		if temp != nil {
			temp.Next = d.Next
		}
	}
	stdio.ByFD(uintptr(d.Descriptor)).Close()
	flush_queues(d)
	if d.Snooping != nil {
		d.Snooping.Snoop_by = nil
	}
	if d.Snoop_by != nil {
		write_to_output(d.Snoop_by, libc.CString("Your victim is no longer among us.\r\n"))
		d.Snoop_by.Snooping = nil
	}
	if d.Character != nil {
		d.Character.Desc = nil
		if !IS_NPC(d.Character) && PLR_FLAGGED(d.Character, PLR_MAILING) && d.Str != nil {
			if *d.Str != nil {
				libc.Free(unsafe.Pointer(*d.Str))
			}
			libc.Free(unsafe.Pointer(d.Str))
			d.Str = nil
		} else if d.Backstr != nil && !IS_NPC(d.Character) && !PLR_FLAGGED(d.Character, PLR_WRITING) {
			libc.Free(unsafe.Pointer(d.Backstr))
			d.Backstr = nil
		}
		if IS_PLAYING(d) || d.Connected == CON_DISCONNECT {
			var link_challenged *char_data
			if d.Original != nil {
				link_challenged = d.Original
			} else {
				link_challenged = d.Character
			}
			act(libc.CString("$n has lost $s link."), TRUE, link_challenged, nil, nil, TO_ROOM)
			save_char(link_challenged)
			mudlog(NRM, MAX(ADMLVL_IMMORT, int(link_challenged.Player_specials.Invis_level)), TRUE, libc.CString("Closing link to: %s."), GET_NAME(link_challenged))
		} else {
			free_char(d.Character)
		}
	} else {
		mudlog(CMP, ADMLVL_IMMORT, TRUE, libc.CString("Losing descriptor without char."))
	}
	if d.Original != nil && d.Original.Desc != nil {
		d.Original.Desc = nil
	}
	if d.History != nil {
		var cnt int
		for cnt = 0; cnt < HISTORY_SIZE; cnt++ {
			if *(**byte)(unsafe.Add(unsafe.Pointer(d.History), unsafe.Sizeof((*byte)(nil))*uintptr(cnt))) != nil {
				libc.Free(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(d.History), unsafe.Sizeof((*byte)(nil))*uintptr(cnt)))))
			}
		}
		libc.Free(unsafe.Pointer(d.History))
	}
	if d.Showstr_head != nil {
		libc.Free(unsafe.Pointer(d.Showstr_head))
	}
	if d.Showstr_count != 0 {
		libc.Free(unsafe.Pointer(d.Showstr_vector))
	}
	if d.Obj_name != nil {
		libc.Free(unsafe.Pointer(d.Obj_name))
	}
	if d.Obj_short != nil {
		libc.Free(unsafe.Pointer(d.Obj_short))
	}
	if d.Obj_long != nil {
		libc.Free(unsafe.Pointer(d.Obj_long))
	}
	free_user(d)
	switch d.Connected {
	case CON_OEDIT:
		fallthrough
	case CON_IEDIT:
		fallthrough
	case CON_REDIT:
		fallthrough
	case CON_ZEDIT:
		fallthrough
	case CON_MEDIT:
		fallthrough
	case CON_SEDIT:
		fallthrough
	case CON_TEDIT:
		fallthrough
	case CON_AEDIT:
		fallthrough
	case CON_TRIGEDIT:
		cleanup_olc(d, CLEANUP_ALL)
	default:
	}
	libc.Free(unsafe.Pointer(d))
}
func check_idle_passwords() {
	var (
		d      *descriptor_data
		next_d *descriptor_data
	)
	for d = descriptor_list; d != nil; d = next_d {
		next_d = d.Next
		if d.Connected != CON_PASSWORD && d.Connected != CON_GET_EMAIL && d.Connected != CON_NEWPASSWD {
			continue
		}
		if int(d.Idle_tics) == 0 {
			d.Idle_tics++
			continue
		} else {
			write_to_output(d, libc.CString("\r\nTimed out... goodbye.\r\n"))
			d.Connected = CON_CLOSE
		}
	}
}
func check_idle_menu() {
	var (
		d      *descriptor_data
		next_d *descriptor_data
	)
	for d = descriptor_list; d != nil; d = next_d {
		next_d = d.Next
		if d.Connected != CON_MENU && d.Connected != CON_GET_USER && d.Connected != CON_UMENU {
			continue
		}
		if int(d.Idle_tics) == 0 {
			d.Idle_tics++
			write_to_output(d, libc.CString("\r\nYou are about to be disconnected due to inactivity in 60 seconds.\r\n"))
			continue
		} else {
			write_to_output(d, libc.CString("\r\nTimed out... goodbye.\r\n"))
			d.Connected = CON_CLOSE
		}
	}
}
func nonblock(s int) {
}
func reread_wizlists(sig int) {
	reread_wizlist = TRUE
}
func unrestrict_game(sig int) {
	emergency_unban = TRUE
}
func reap(sig int) {
}
func checkpointing(sig int) {
	if tics_passed == 0 {
		basic_mud_log(libc.CString("SYSERR: CHECKPOINT shutdown: tics not updated. (Infinite loop suspected)"))
		panic("abort")
	} else {
		tics_passed = 0
	}
}
func hupsig(sig int) {
	basic_mud_log(libc.CString("SYSERR: Received SIGHUP, SIGINT, or SIGTERM.  Shutting down..."))
	os.Exit(1)
}
func signal_setup() {
}
func send_to_char(ch *char_data, messg *byte, _rest ...interface{}) uint64 {
	if ch.Desc != nil && messg != nil && *messg != 0 {
		var (
			left uint64
			args libc.ArgList
		)
		args.Start(messg, _rest)
		left = vwrite_to_output(ch.Desc, messg, args)
		args.End()
		return left
	}
	return 0
}
func arena_watch(ch *char_data) int {
	var (
		d     *descriptor_data
		found int = FALSE
		room  int = int(-1)
	)
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected != CON_PLAYING {
			continue
		}
		if IN_ARENA(d.Character) {
			if ch.Arenawatch == int(d.Character.Idnum) {
				found = TRUE
				if d.Character.In_room != room_rnum(-1) && d.Character.In_room <= top_of_world {
					room = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(d.Character.In_room)))).Number)
				} else {
					room = -1
				}
			}
		}
	}
	if found == FALSE {
		ch.Player_specials.Pref[int(PRF_ARENAWATCH/32)] &= bitvector_t(int32(^(1 << (int(PRF_ARENAWATCH % 32)))))
		ch.Arenawatch = -1
		return -1
	} else {
		return room
	}
}
func send_to_eaves(messg *byte, tch *char_data, _rest ...interface{}) {
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected != CON_PLAYING {
			continue
		}
		var roll int = rand_number(1, 101)
		if d.Character.Listenroom == (func() room_vnum {
			if tch.In_room != room_rnum(-1) && tch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(tch.In_room)))).Number
			}
			return -1
		}()) && GET_SKILL(d.Character, SKILL_EAVESDROP) > roll {
			var (
				buf  [1000]byte
				buf2 [1000]byte
			)
			buf[0] = '\x00'
			stdio.Sprintf(&buf2[0], "@W%s %s\r\n", PERS(d.Character, tch), messg)
			stdio.Sprintf(&buf[0], "-----Eavesdrop-----\r\n%s-----Eavesdrop-----\r\n", &buf2[0])
			send_to_char(d.Character, &buf[0])
		}
	}
}
func send_to_all(messg *byte, _rest ...interface{}) {
	var (
		i    *descriptor_data
		args libc.ArgList
	)
	if messg == nil {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		args.Start(messg, _rest)
		vwrite_to_output(i, messg, args)
		args.End()
	}
}
func send_to_outdoor(messg *byte, _rest ...interface{}) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		var args libc.ArgList
		if i.Connected != CON_PLAYING || i.Character == nil {
			continue
		}
		if !AWAKE(i.Character) || !OUTSIDE(i.Character) {
			continue
		}
		args.Start(messg, _rest)
		vwrite_to_output(i, messg, args)
		args.End()
	}
}
func send_to_moon(messg *byte, _rest ...interface{}) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		var args libc.ArgList
		if i.Connected != CON_PLAYING || i.Character == nil {
			continue
		}
		if !AWAKE(i.Character) || !HAS_MOON(i.Character) {
			continue
		}
		args.Start(messg, _rest)
		vwrite_to_output(i, messg, args)
		args.End()
	}
}
func send_to_planet(type_ int, planet int, messg *byte, _rest ...interface{}) {
	var i *descriptor_data
	if messg == nil || *messg == 0 {
		return
	}
	for i = descriptor_list; i != nil; i = i.Next {
		var args libc.ArgList
		if i.Connected != CON_PLAYING || i.Character == nil {
			continue
		}
		if !AWAKE(i.Character) || !ROOM_FLAGGED(i.Character.In_room, bitvector_t(int32(planet))) {
			continue
		} else {
			if type_ == 0 {
				args.Start(messg, _rest)
				vwrite_to_output(i, messg, args)
				args.End()
			} else if OUTSIDE(i.Character) && GET_SKILL(i.Character, SKILL_SPOT) >= axion_dice(-5) {
				args.Start(messg, _rest)
				vwrite_to_output(i, messg, args)
				args.End()
			}
		}
	}
}
func send_to_room(room room_rnum, messg *byte, _rest ...interface{}) {
	var (
		i    *char_data
		args libc.ArgList
	)
	if messg == nil {
		return
	}
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(room)))).People; i != nil; i = i.Next_in_room {
		if i.Desc == nil {
			continue
		}
		args.Start(messg, _rest)
		vwrite_to_output(i.Desc, messg, args)
		args.End()
	}
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected != CON_PLAYING {
			continue
		}
		if PRF_FLAGGED(d.Character, PRF_ARENAWATCH) {
			if arena_watch(d.Character) == int(room) {
				var buf [2000]byte
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "@c-----@CArena@c-----@n\r\n%s\r\n@c-----@CArena@c-----@n\r\n", messg)
				args.Start(messg, _rest)
				vwrite_to_output(d, &buf[0], args)
				args.End()
			}
		}
		if d.Character.Listenroom > 0 {
			var roll int = rand_number(1, 101)
			if d.Character.Listenroom == room_vnum(room) && GET_SKILL(d.Character, SKILL_EAVESDROP) > roll {
				var buf [1000]byte
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "-----Eavesdrop-----\r\n%s\r\n-----Eavesdrop-----\r\n", messg)
				args.Start(messg, _rest)
				vwrite_to_output(d, &buf[0], args)
				args.End()
			}
		}
	}
}

var ACTNULL *byte = libc.CString("<NULL>")

func perform_act(orig *byte, ch *char_data, obj *obj_data, vict_obj unsafe.Pointer, to *char_data) {
	var (
		i             *byte = nil
		lbuf          [64936]byte
		buf           *byte
		j             *byte
		uppercasenext bool       = FALSE != 0
		dg_victim     *char_data = nil
		dg_target     *obj_data  = nil
		dg_arg        *byte      = nil
	)
	buf = &lbuf[0]
	for {
		if *orig == '$' {
			switch *(func() *byte {
				p := &orig
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return *p
			}()) {
			case 'n':
				i = PERS(ch, to)
			case 'N':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = PERS((*char_data)(vict_obj), to)
				}
				dg_victim = (*char_data)(vict_obj)
			case 'm':
				i = HMHR(ch)
			case 'M':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = HMHR((*char_data)(vict_obj))
				}
				dg_victim = (*char_data)(vict_obj)
			case 's':
				i = HSHR(ch)
			case 'S':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = HSHR((*char_data)(vict_obj))
				}
				dg_victim = (*char_data)(vict_obj)
			case 'e':
				i = HSSH(ch)
			case 'E':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = HSSH((*char_data)(vict_obj))
				}
				dg_victim = (*char_data)(vict_obj)
			case 'o':
				if obj == nil {
					i = ACTNULL
				} else {
					i = OBJN(obj, to)
				}
			case 'O':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = OBJN((*obj_data)(vict_obj), to)
				}
				dg_target = (*obj_data)(vict_obj)
			case 'p':
				if obj == nil {
					i = ACTNULL
				} else {
					i = OBJS(obj, to)
				}
			case 'P':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = OBJS((*obj_data)(vict_obj), to)
				}
				dg_target = (*obj_data)(vict_obj)
			case 'a':
				if obj == nil {
					i = ACTNULL
				} else {
					i = SANA(obj)
				}
			case 'A':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = SANA((*obj_data)(vict_obj))
				}
				dg_target = (*obj_data)(vict_obj)
			case 'T':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = (*byte)(vict_obj)
				}
				dg_arg = (*byte)(vict_obj)
			case 't':
				if obj == nil {
					i = ACTNULL
				} else {
					i = (*byte)(unsafe.Pointer(obj))
				}
			case 'F':
				if vict_obj == nil {
					i = ACTNULL
				} else {
					i = fname((*byte)(vict_obj))
				}
			case 'u':
				for j = buf; uintptr(unsafe.Pointer(j)) > uintptr(unsafe.Pointer(&lbuf[0])) && !unicode.IsSpace(rune(int(*((*byte)(unsafe.Add(unsafe.Pointer(j), -1)))))); j = (*byte)(unsafe.Add(unsafe.Pointer(j), -1)) {
				}
				if j != buf {
					*j = byte(int8(unicode.ToUpper(rune(*j))))
				}
				i = libc.CString("")
			case 'U':
				uppercasenext = TRUE != 0
				i = libc.CString("")
			case '$':
				i = libc.CString("$")
			default:
				return
			}
			for (func() byte {
				p := buf
				*buf = *(func() *byte {
					p := &i
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}())
				return *p
			}()) != 0 {
				if uppercasenext && !unicode.IsSpace(rune(int(*buf))) {
					*buf = byte(int8(unicode.ToUpper(rune(*buf))))
					uppercasenext = FALSE != 0
				}
				buf = (*byte)(unsafe.Add(unsafe.Pointer(buf), 1))
			}
			orig = (*byte)(unsafe.Add(unsafe.Pointer(orig), 1))
		} else if (func() byte {
			p := (func() *byte {
				p := &buf
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			*(func() *byte {
				p := &buf
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *(func() *byte {
				p := &orig
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			return *p
		}()) == 0 {
			break
		} else if uppercasenext && !unicode.IsSpace(rune(int(*((*byte)(unsafe.Add(unsafe.Pointer(buf), -1)))))) {
			*((*byte)(unsafe.Add(unsafe.Pointer(buf), -1))) = byte(int8(unicode.ToUpper(rune(*((*byte)(unsafe.Add(unsafe.Pointer(buf), -1)))))))
			uppercasenext = FALSE != 0
		}
	}
	*(func() *byte {
		p := &buf
		*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), -1))
		return *p
	}()) = '\r'
	*(func() *byte {
		p := &buf
		*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
		return *p
	}()) = '\n'
	*(func() *byte {
		p := &buf
		*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
		return *p
	}()) = '\x00'
	if to.Desc != nil {
		write_to_output(to.Desc, libc.CString("%s"), CAP(&lbuf[0]))
	}
	if IS_NPC(to) && dg_act_check != 0 && to != ch {
		act_mtrigger(to, &lbuf[0], ch, dg_victim, obj, dg_target, dg_arg)
	}
	if last_act_message != nil {
		libc.Free(unsafe.Pointer(last_act_message))
	}
	last_act_message = libc.StrDup(&lbuf[0])
}

var to_sleeping int = 0

func act(str *byte, hide_invisible int, ch *char_data, obj *obj_data, vict_obj unsafe.Pointer, type_ int) *byte {
	var (
		to        *char_data
		res_sneak int
		res_hide  int
		dcval     int = 0
		resskill  int = 0
	)
	to_sleeping = 0
	if str == nil || *str == 0 {
		return nil
	}
	if (func() int {
		to_sleeping = type_ & (2 << 7)
		return to_sleeping
	}()) != 0 {
		type_ &= ^(2 << 7)
	}
	if (func() int {
		res_sneak = type_ & (2 << 9)
		return res_sneak
	}()) != 0 {
		type_ &= ^(2 << 9)
	}
	if (func() int {
		res_hide = type_ & (2 << 10)
		return res_hide
	}()) != 0 {
		type_ &= ^(2 << 10)
	}
	if res_sneak != 0 && AFF_FLAGGED(ch, AFF_SNEAK) {
		dcval = roll_skill(ch, SKILL_MOVE_SILENTLY)
		if GET_SKILL(ch, SKILL_BALANCE) != 0 {
			dcval += GET_SKILL(ch, SKILL_BALANCE) / 10
		}
		if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 5 || (ch.Genome[1]) == 5) {
			dcval += 10
		}
		resskill = SKILL_SPOT
	} else if res_hide != 0 && AFF_FLAGGED(ch, AFF_HIDE) {
		dcval = roll_skill(ch, SKILL_HIDE)
		if GET_SKILL(ch, SKILL_BALANCE) != 0 {
			dcval += GET_SKILL(ch, SKILL_BALANCE) / 10
		}
		resskill = SKILL_SPOT
	}
	if (func() int {
		dg_act_check = int(libc.BoolToInt((type_ & (2 << 8)) == 0))
		return dg_act_check
	}()) == 0 {
		type_ &= ^(2 << 8)
	}
	if type_ == TO_CHAR {
		if ch != nil && SENDOK(ch) && (resskill == 0 || roll_skill(ch, resskill) >= dcval) {
			perform_act(str, ch, obj, vict_obj, ch)
			return last_act_message
		}
		return nil
	}
	if type_ == TO_VICT {
		if (func() *char_data {
			to = (*char_data)(vict_obj)
			return to
		}()) != nil && SENDOK(to) && (resskill == 0 || roll_skill(to, resskill) >= dcval) {
			perform_act(str, ch, obj, vict_obj, to)
			return last_act_message
		}
		return nil
	}
	if type_ == TO_GMOTE {
		var (
			i   *descriptor_data
			buf [64936]byte
		)
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Connected == 0 && i.Character != nil && !PRF_FLAGGED(i.Character, PRF_NOGOSS) && !PLR_FLAGGED(i.Character, PLR_WRITING) && !ROOM_FLAGGED(i.Character.In_room, ROOM_SOUNDPROOF) {
				stdio.Sprintf(&buf[0], "@y%s@n", str)
				perform_act(&buf[0], ch, obj, vict_obj, i.Character)
				var buf2 [64936]byte
				stdio.Sprintf(&buf2[0], "%s\r\n", &buf[0])
				add_history(i.Character, &buf2[0], HIST_GOSSIP)
			}
		}
		return last_act_message
	}
	if ch != nil && ch.In_room != room_rnum(-1) {
		to = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People
	} else if obj != nil && obj.In_room != room_rnum(-1) {
		to = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).People
	} else {
		return nil
	}
	if (type_ & TO_ROOM) != 0 {
		var d *descriptor_data
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Connected != CON_PLAYING {
				continue
			}
			if ch != nil {
				if IN_ARENA(ch) {
					if PRF_FLAGGED(d.Character, PRF_ARENAWATCH) {
						if arena_watch(d.Character) == int(func() room_vnum {
							if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
							}
							return -1
						}()) {
							var buf3 [2000]byte
							buf3[0] = '\x00'
							stdio.Sprintf(&buf3[0], "@c-----@CArena@c-----@n\r\n%s\r\n@c-----@CArena@c-----@n\r\n", str)
							perform_act(&buf3[0], ch, obj, vict_obj, d.Character)
						}
					}
				}
			}
			if d.Character.Listenroom > 0 {
				var roll int = rand_number(1, 101)
				if resskill == 0 || roll_skill(d.Character, resskill) >= dcval {
					if ch != nil && d.Character.Listenroom == (func() room_vnum {
						if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
						}
						return -1
					}()) && GET_SKILL(d.Character, SKILL_EAVESDROP) > roll {
						var buf3 [1000]byte
						buf3[0] = '\x00'
						stdio.Sprintf(&buf3[0], "-----Eavesdrop-----\r\n%s\r\n-----Eavesdrop-----\r\n", str)
						perform_act(&buf3[0], ch, obj, vict_obj, d.Character)
					} else if obj != nil && d.Character.Listenroom == (func() room_vnum {
						if obj.In_room != room_rnum(-1) && obj.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(obj.In_room)))).Number
						}
						return -1
					}()) && GET_SKILL(d.Character, SKILL_EAVESDROP) > roll {
						var buf3 [1000]byte
						buf3[0] = '\x00'
						stdio.Sprintf(&buf3[0], "-----Eavesdrop-----\r\n%s\r\n-----Eavesdrop-----\r\n", str)
						perform_act(&buf3[0], ch, obj, vict_obj, d.Character)
					}
				}
			}
		}
	}
	for ; to != nil; to = to.Next_in_room {
		if !SENDOK(to) || to == ch {
			continue
		}
		if hide_invisible != 0 && ch != nil && !CAN_SEE(to, ch) {
			continue
		}
		if type_ != TO_ROOM && unsafe.Pointer(to) == vict_obj {
			continue
		}
		if resskill != 0 && roll_skill(to, resskill) < dcval {
			continue
		}
		perform_act(str, ch, obj, vict_obj, to)
	}
	return last_act_message
}
func setup_log(filename *byte, fd int) {
	var s_fp *stdio.File
	s_fp = stdio.Stderr()
	if filename == nil || *filename == '\x00' {
		logfile = s_fp
		puts(libc.CString("Using file descriptor for logging."))
		return
	}
	if open_logfile(filename, s_fp) != 0 {
		return
	}
	if open_logfile(libc.CString("log/syslog"), s_fp) != 0 {
		return
	}
	if open_logfile(libc.CString("syslog"), s_fp) != 0 {
		return
	}
	puts(libc.CString("SYSERR: Couldn't open anything to log to, giving up."))
	os.Exit(1)
}
func open_logfile(filename *byte, stderr_fp *stdio.File) int {
	if stderr_fp != nil {
		logfile = freopen(filename, libc.CString("w"), stderr_fp)
	} else {
		logfile = stdio.FOpen(libc.GoString(filename), "w")
	}
	if logfile != nil {
		stdio.Printf("Using log file '%s'%s.\n", filename, func() string {
			if stderr_fp != nil {
				return " with redirection"
			}
			return ""
		}())
		return TRUE
	}
	stdio.Printf("SYSERR: Error opening file '%s': %s\n", filename, libc.StrError(libc.Errno))
	return FALSE
}
func circle_sleep(timeout *libc.TimeVal) {
}
func show_help(t *descriptor_data, entry *byte) {
	var (
		chk    int
		bot    int
		top    int
		mid    int
		minlen int
		buf    [64936]byte
	)
	if help_table == nil {
		return
	}
	bot = 0
	top = top_of_helpt
	minlen = libc.StrLen(entry)
	for {
		mid = (bot + top) / 2
		if bot > top {
			return
		} else if (func() int {
			chk = libc.StrNCaseCmp(entry, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Keywords, minlen)
			return chk
		}()) == 0 {
			for mid > 0 && (func() int {
				chk = libc.StrNCaseCmp(entry, (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid-1)))).Keywords, minlen)
				return chk
			}()) == 0 {
				mid--
			}
			write_to_output(t, libc.CString("\r\n"))
			stdio.Snprintf(&buf[0], int(64936), "%s\r\n[ PRESS RETURN TO CONTINUE ]", (*(*help_index_element)(unsafe.Add(unsafe.Pointer(help_table), unsafe.Sizeof(help_index_element{})*uintptr(mid)))).Entry)
			page_string(t, &buf[0], 0)
			return
		} else {
			if chk > 0 {
				bot = mid + 1
			} else {
				top = mid - 1
			}
		}
	}
}
func send_to_range(start room_vnum, finish room_vnum, messg *byte, _rest ...interface{}) {
	var (
		i    *char_data
		args libc.ArgList
		j    int
	)
	if start > finish {
		basic_mud_log(libc.CString("send_to_range passed start room value greater then finish."))
		return
	}
	if messg == nil {
		return
	}
	for j = 0; j < int(top_of_world); j++ {
		if (func() room_vnum {
			if j != int(-1) && j <= int(top_of_world) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(j)))).Number
			}
			return -1
		}()) >= start && (func() room_vnum {
			if j != int(-1) && j <= int(top_of_world) {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(j)))).Number
			}
			return -1
		}()) <= finish {
			for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(j)))).People; i != nil; i = i.Next_in_room {
				if i.Desc == nil {
					continue
				}
				args.Start(messg, _rest)
				vwrite_to_output(i.Desc, messg, args)
				args.End()
			}
		}
	}
}
func passcomm(ch *char_data, comm *byte) int {
	if libc.StrCaseCmp(comm, libc.CString("score")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("sco")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("ooc")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("newbie")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("newb")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("look")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("lo")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("l")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("status")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("stat")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("sta")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("tell")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("reply")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("say")) == 0 {
		return TRUE
	} else if libc.StrCaseCmp(comm, libc.CString("osay")) == 0 {
		return TRUE
	} else {
		return FALSE
	}
}
