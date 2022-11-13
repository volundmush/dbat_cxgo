package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unicode"
)

func parse_stats(d *descriptor_data, arg *byte) int {
	var ch *char_data
	ch = d.Character
	switch d.Olc.Mode {
	case STAT_QUIT:
		return 1
	case STAT_PARSE_MENU:
		if parse_stat_menu(d, arg) != 0 {
			return 1
		}
	case STAT_GET_STR:
		ch.Real_abils.Str = int8(stats_assign_stat(int(ch.Real_abils.Str), arg, d))
		stats_disp_menu(d)
	case STAT_GET_INT:
		ch.Real_abils.Intel = int8(stats_assign_stat(int(ch.Real_abils.Intel), arg, d))
		stats_disp_menu(d)
	case STAT_GET_WIS:
		ch.Real_abils.Wis = int8(stats_assign_stat(int(ch.Real_abils.Wis), arg, d))
		stats_disp_menu(d)
	case STAT_GET_DEX:
		ch.Real_abils.Dex = int8(stats_assign_stat(int(ch.Real_abils.Dex), arg, d))
		stats_disp_menu(d)
	case STAT_GET_CON:
		ch.Real_abils.Con = int8(stats_assign_stat(int(ch.Real_abils.Con), arg, d))
		stats_disp_menu(d)
	case STAT_GET_CHA:
		ch.Real_abils.Cha = int8(stats_assign_stat(int(ch.Real_abils.Cha), arg, d))
		stats_disp_menu(d)
	default:
		d.Olc.Mode = stats_disp_menu(d)
	}
	return 0
}
func stats_disp_menu(d *descriptor_data) int {
	send_to_char(d.Character, libc.CString("\r\n@W-<[@y==========@B[ @YCWG @B]@y==========@W]>-\r\n <| Total Points Left: @m%3d@W    |>     You should select the letter of the score you\r\n <|                           |>     wish to adjust.  When prompted, enter the new score,\r\n <| = Select a stat:          |>     NOT the amount to add.  NOTE: If you quit before you\r\n <| @BS@W) @rStrength     : @m%3d@W     |>     assign all the points, you will lose them forever.\r\n <| @BD@W) @rDexterity    : @m%3d@W     |>     If your points are at zero, you may still reassign\r\n <| @BN@W) @rConstitution : @m%3d@W     |>     points by lowering any statistic, then add those\r\n <| @BI@W) @rIntelligence : @m%3d@W     |>     points to the statistic of your choice.\r\n <| @BW@W) @rWisdom       : @m%3d@W     |>\r\n <| @BC@W) @rCharisma     : @m%3d@W     |>\r\n <| @BQ@W) @CQuit@W                   |>\r\n-<[@y===========================@W]>-@n\r\n\r\n"), d.Olc.Value, d.Character.Real_abils.Str, d.Character.Real_abils.Dex, d.Character.Real_abils.Con, d.Character.Real_abils.Intel, d.Character.Real_abils.Wis, d.Character.Real_abils.Cha)
	send_to_char(d.Character, libc.CString("Enter Letter to Change: "))
	d.Olc.Mode = STAT_PARSE_MENU
	return 1
}
func parse_stat_menu(d *descriptor_data, arg *byte) int {
	*arg = byte(int8(unicode.ToLower(rune(*arg))))
	switch *arg {
	case 's':
		d.Olc.Mode = STAT_GET_STR
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'i':
		d.Olc.Mode = STAT_GET_INT
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'w':
		d.Olc.Mode = STAT_GET_WIS
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'd':
		d.Olc.Mode = STAT_GET_DEX
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'n':
		d.Olc.Mode = STAT_GET_CON
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'c':
		d.Olc.Mode = STAT_GET_CHA
		send_to_char(d.Character, libc.CString("Enter new value: "))
	case 'q':
		d.Olc.Mode = STAT_QUIT
		return 1
	default:
		stats_disp_menu(d)
	}
	return 0
}
func stats_assign_stat(abil int, arg *byte, d *descriptor_data) int {
	var temp int
	if abil > 0 {
		d.Olc.Value = d.Olc.Value + abil
		abil = 0
	}
	if libc.Atoi(libc.GoString(arg)) > d.Olc.Value {
		temp = d.Olc.Value
	} else {
		temp = libc.Atoi(libc.GoString(arg))
	}
	if temp > 100 {
		if d.Olc.Value < 100 {
			temp = d.Olc.Value
		} else {
			temp = 100
		}
	}
	if temp < 3 {
		temp = 3
	}
	if d.Olc.Value <= 0 {
		temp = 0
		d.Olc.Value = 0
		mudlog(NRM, ADMLVL_IMMORT, TRUE, libc.CString("Stat total below 0: possible code error"))
	}
	abil = temp
	d.Olc.Value -= temp
	return abil
}
