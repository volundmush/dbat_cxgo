package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unsafe"
)

func do_action(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		part   [2048]byte
		act_nr int
		action *social_messg
		vict   *char_data
		targ   *obj_data
	)
	if (func() int {
		act_nr = find_action(cmd)
		return act_nr
	}()) < 0 {
		send_to_char(ch, libc.CString("That action is not supported.\r\n"))
		return
	}
	action = &soc_mess_list[act_nr]
	if argument == nil || *argument == 0 {
		send_to_char(ch, libc.CString("%s\r\n"), action.Char_no_arg)
		act(action.Others_no_arg, action.Hide, ch, nil, nil, TO_ROOM)
		return
	}
	two_arguments(argument, &arg[0], &part[0])
	if action.Char_body_found == nil && (part[0]) != 0 {
		send_to_char(ch, libc.CString("Sorry, this social does not support body parts.\r\n"))
		return
	}
	if action.Char_found == nil {
		arg[0] = '\x00'
	}
	if action.Char_found != nil && argument != nil {
		one_argument(argument, &arg[0])
	} else {
		arg[0] = '\x00'
	}
	vict = get_char_vis(ch, &arg[0], nil, 1<<0)
	if vict == nil {
		if action.Char_obj_found != nil {
			targ = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			if targ == nil {
				targ = get_obj_in_list_vis(ch, &arg[0], nil, world[ch.In_room].Contents)
			}
			if targ != nil {
				act(action.Char_obj_found, action.Hide, ch, targ, nil, TO_CHAR)
				act(action.Others_obj_found, action.Hide, ch, targ, nil, TO_ROOM)
				return
			}
		}
		if action.Not_found != nil {
			send_to_char(ch, libc.CString("%s\r\n"), action.Not_found)
		} else {
			send_to_char(ch, libc.CString("I don't see anything by that name here.\r\n"))
		}
		return
	}
	if vict == ch {
		if action.Char_auto != nil {
			send_to_char(ch, libc.CString("%s\r\n"), action.Char_auto)
		} else {
			send_to_char(ch, libc.CString("Erm, no.\r\n"))
		}
		act(action.Others_auto, action.Hide, ch, nil, nil, TO_ROOM)
		return
	}
	if int(vict.Position) < action.Min_victim_position {
		act(libc.CString("$N is not in a proper position for that."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
	} else {
		if part[0] != 0 {
			act(action.Char_body_found, 0, ch, (*obj_data)(unsafe.Pointer(&part[0])), unsafe.Pointer(vict), int(TO_CHAR|2<<7))
			act(action.Others_body_found, action.Hide, ch, (*obj_data)(unsafe.Pointer(&part[0])), unsafe.Pointer(vict), TO_NOTVICT)
			act(action.Vict_body_found, action.Hide, ch, (*obj_data)(unsafe.Pointer(&part[0])), unsafe.Pointer(vict), TO_VICT)
		} else {
			act(action.Char_found, 0, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
			act(action.Others_found, action.Hide, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			act(action.Vict_found, action.Hide, ch, nil, unsafe.Pointer(vict), TO_VICT)
		}
	}
}
func do_insult(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg    [2048]byte
		victim *char_data
	)
	one_argument(argument, &arg[0])
	if arg[0] != 0 {
		if (func() *char_data {
			victim = get_char_vis(ch, &arg[0], nil, 1<<0)
			return victim
		}()) == nil {
			send_to_char(ch, libc.CString("Can't hear you!\r\n"))
		} else {
			if victim != ch {
				send_to_char(ch, libc.CString("You insult %s.\r\n"), GET_NAME(victim))
				switch rand_number(0, 2) {
				case 0:
					if int(ch.Sex) == SEX_MALE {
						if int(victim.Sex) == SEX_MALE {
							act(libc.CString("$n accuses you of fighting like a woman!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						} else {
							act(libc.CString("$n says that women can't fight."), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						}
					} else {
						if int(victim.Sex) == SEX_MALE {
							act(libc.CString("$n accuses you of having the smallest... (brain?)"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						} else {
							act(libc.CString("$n tells you that you'd lose a beauty contest against a troll."), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						}
					}
				case 1:
					act(libc.CString("$n calls your mother a bitch!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
				default:
					act(libc.CString("$n tells you to get lost!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
				}
				act(libc.CString("$n insults $N."), TRUE, ch, nil, unsafe.Pointer(victim), TO_NOTVICT)
			} else {
				send_to_char(ch, libc.CString("You feel insulted.\r\n"))
			}
		}
	} else {
		send_to_char(ch, libc.CString("I'm sure you don't want to insult *everybody*...\r\n"))
	}
}
func boot_social_messages() {
	var (
		fl           *stdio.File
		nr           int = 0
		hide         int
		min_char_pos int
		min_pos      int
		min_lvl      int
		curr_soc     int = -1
		next_soc     [64936]byte
		sorted       [2048]byte
	)
	if config_info.Operation.Use_new_socials == TRUE {
		if (func() *stdio.File {
			fl = stdio.FOpen(LIB_MISC, "r")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: can't open socials file '%s': %s"), LIB_MISC, libc.StrError(libc.Errno))
			os.Exit(1)
		}
		next_soc[0] = '\x00'
		for int(fl.IsEOF()) == 0 {
			fl.GetS(&next_soc[0], MAX_STRING_LENGTH)
			if next_soc[0] == '~' {
				top_of_socialt++
			}
		}
	} else {
		if (func() *stdio.File {
			fl = stdio.FOpen(LIB_MISC, "r")
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: can't open socials file '%s': %s"), LIB_MISC, libc.StrError(libc.Errno))
			os.Exit(1)
		}
		for int(fl.IsEOF()) == 0 {
			fl.GetS(&next_soc[0], MAX_STRING_LENGTH)
			if next_soc[0] == '\n' || next_soc[0] == '\r' {
				top_of_socialt++
			}
		}
	}
	basic_mud_log(libc.CString("Social table contains %d socials."), top_of_socialt)
	fl.Seek(0, 0)
	soc_mess_list = make([]social_messg, top_of_socialt+1)
	for {
		stdio.Fscanf(fl, " %s ", &next_soc[0])
		if next_soc[0] == '$' {
			break
		}
		if config_info.Operation.Use_new_socials == TRUE {
			if stdio.Fscanf(fl, " %s %d %d %d %d \n", &sorted[0], &hide, &min_char_pos, &min_pos, &min_lvl) != 5 {
				basic_mud_log(libc.CString("SYSERR: format error in social file near social '%s'"), &next_soc[0])
				os.Exit(1)
			}
			curr_soc++
			soc_mess_list[curr_soc].Command = libc.StrDup(&next_soc[1])
			soc_mess_list[curr_soc].Sort_as = libc.StrDup(&sorted[0])
			soc_mess_list[curr_soc].Hide = hide
			soc_mess_list[curr_soc].Min_char_position = min_char_pos
			soc_mess_list[curr_soc].Min_victim_position = min_pos
			soc_mess_list[curr_soc].Min_level_char = min_lvl
		} else {
			if stdio.Fscanf(fl, " %d %d \n", &hide, &min_pos) != 2 {
				basic_mud_log(libc.CString("SYSERR: format error in social file near social '%s'"), &next_soc[0])
				os.Exit(1)
			}
			curr_soc++
			soc_mess_list[curr_soc].Command = libc.StrDup(&next_soc[0])
			soc_mess_list[curr_soc].Sort_as = libc.StrDup(&next_soc[0])
			soc_mess_list[curr_soc].Hide = hide
			soc_mess_list[curr_soc].Min_char_position = POS_RESTING
			soc_mess_list[curr_soc].Min_victim_position = min_pos
			soc_mess_list[curr_soc].Min_level_char = 0
		}
		soc_mess_list[curr_soc].Char_no_arg = fread_action(fl, nr)
		soc_mess_list[curr_soc].Others_no_arg = fread_action(fl, nr)
		soc_mess_list[curr_soc].Char_found = fread_action(fl, nr)
		if config_info.Operation.Use_new_socials == FALSE && soc_mess_list[curr_soc].Char_found == nil {
			continue
		}
		soc_mess_list[curr_soc].Others_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Vict_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Not_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Char_auto = fread_action(fl, nr)
		soc_mess_list[curr_soc].Others_auto = fread_action(fl, nr)
		if config_info.Operation.Use_new_socials == FALSE {
			continue
		}
		soc_mess_list[curr_soc].Char_body_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Others_body_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Vict_body_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Char_obj_found = fread_action(fl, nr)
		soc_mess_list[curr_soc].Others_obj_found = fread_action(fl, nr)
	}
	fl.Close()
	if curr_soc > top_of_socialt {
		panic("assert failed")
	}
	top_of_socialt = curr_soc
}
func create_command_list() {
	var (
		i    int
		j    int
		k    int
		temp social_messg
	)
	if complete_cmd_info != nil {
		free_command_list()
	}
	for j = 0; j < top_of_socialt; j++ {
		k = j
		for i = j + 1; i <= top_of_socialt; i++ {
			if libc.StrCaseCmp(soc_mess_list[i].Sort_as, soc_mess_list[k].Sort_as) < 0 {
				k = i
			}
		}
		if j != k {
			temp = soc_mess_list[j]
			soc_mess_list[j] = soc_mess_list[k]
			soc_mess_list[k] = temp
		}
	}
	i = 0
	for *cmd_info[i].Command != '\n' {
		i++
	}
	i++
	complete_cmd_info = make([]command_info, top_of_socialt+i+2)
	i = 0
	j = 0
	k = 0
	for *cmd_info[i].Command != '\n' {
		complete_cmd_info[func() int {
			p := &k
			x := *p
			*p++
			return x
		}()] = cmd_info[func() int {
			p := &i
			x := *p
			*p++
			return x
		}()]
	}
	for j <= top_of_socialt {
		soc_mess_list[j].Act_nr = k
		complete_cmd_info[k].Command = soc_mess_list[j].Command
		complete_cmd_info[k].Sort_as = soc_mess_list[j].Sort_as
		complete_cmd_info[k].Minimum_position = int8(soc_mess_list[j].Min_char_position)
		complete_cmd_info[k].Command_pointer = func(ch *char_data, argument *byte, cmd int, subcmd int) {
			func(ch *char_data, argument *byte, cmd int, subcmd int) {
				do_action(ch, argument, cmd, subcmd)
			}(ch, argument, cmd, subcmd)
		}
		complete_cmd_info[k].Minimum_level = int16(soc_mess_list[func() int {
			p := &j
			x := *p
			*p++
			return x
		}()].Min_level_char)
		complete_cmd_info[k].Minimum_admlevel = ADMLVL_NONE
		complete_cmd_info[func() int {
			p := &k
			x := *p
			*p++
			return x
		}()].Subcmd = 0
	}
	complete_cmd_info[k].Command = libc.CString("\n")
	complete_cmd_info[k].Sort_as = libc.CString("zzzzzzz")
	complete_cmd_info[k].Minimum_position = 0
	complete_cmd_info[k].Command_pointer = nil
	complete_cmd_info[k].Minimum_level = 0
	complete_cmd_info[k].Minimum_admlevel = 0
	complete_cmd_info[k].Subcmd = 0
	basic_mud_log(libc.CString("Command info rebuilt, %d total commands."), k)
}
func free_command_list() {
	libc.Free(unsafe.Pointer(&complete_cmd_info[0]))
	complete_cmd_info = nil
}
func fread_action(fl *stdio.File, nr int) *byte {
	var buf [64936]byte
	fl.GetS(&buf[0], MAX_STRING_LENGTH)
	if int(fl.IsEOF()) != 0 {
		basic_mud_log(libc.CString("SYSERR: fread_action: unexpected EOF near action #%d"), nr)
		os.Exit(1)
	}
	if buf[0] == '#' {
		return nil
	}
	buf[libc.StrLen(&buf[0])-1] = '\x00'
	return libc.StrDup(&buf[0])
}
func free_social_messages() {
	var (
		mess *social_messg
		i    int
	)
	for i = 0; i <= top_of_socialt; i++ {
		mess = &soc_mess_list[i]
		free_action(mess)
	}
	libc.Free(unsafe.Pointer(&soc_mess_list[0]))
}
func free_action(mess *social_messg) {
	if mess.Command != nil {
		libc.Free(unsafe.Pointer(mess.Command))
	}
	if mess.Sort_as != nil {
		libc.Free(unsafe.Pointer(mess.Sort_as))
	}
	if mess.Char_no_arg != nil {
		libc.Free(unsafe.Pointer(mess.Char_no_arg))
	}
	if mess.Others_no_arg != nil {
		libc.Free(unsafe.Pointer(mess.Others_no_arg))
	}
	if mess.Char_found != nil {
		libc.Free(unsafe.Pointer(mess.Char_found))
	}
	if mess.Others_found != nil {
		libc.Free(unsafe.Pointer(mess.Others_found))
	}
	if mess.Vict_found != nil {
		libc.Free(unsafe.Pointer(mess.Vict_found))
	}
	if mess.Char_body_found != nil {
		libc.Free(unsafe.Pointer(mess.Char_body_found))
	}
	if mess.Others_body_found != nil {
		libc.Free(unsafe.Pointer(mess.Others_body_found))
	}
	if mess.Vict_body_found != nil {
		libc.Free(unsafe.Pointer(mess.Vict_body_found))
	}
	if mess.Not_found != nil {
		libc.Free(unsafe.Pointer(mess.Not_found))
	}
	if mess.Char_auto != nil {
		libc.Free(unsafe.Pointer(mess.Char_auto))
	}
	if mess.Others_auto != nil {
		libc.Free(unsafe.Pointer(mess.Others_auto))
	}
	if mess.Char_obj_found != nil {
		libc.Free(unsafe.Pointer(mess.Char_obj_found))
	}
	if mess.Others_obj_found != nil {
		libc.Free(unsafe.Pointer(mess.Others_obj_found))
	}
	*mess = social_messg{}
}
func find_action(cmd int) int {
	var (
		bot int
		top int
		mid int
	)
	bot = 0
	top = top_of_socialt
	if top < 0 {
		return -1
	}
	for {
		mid = (bot + top) / 2
		if soc_mess_list[mid].Act_nr == cmd {
			return mid
		}
		if bot >= top {
			return -1
		}
		if soc_mess_list[mid].Act_nr > cmd {
			top = func() int {
				p := &mid
				*p--
				return *p
			}()
		} else {
			bot = func() int {
				p := &mid
				*p++
				return *p
			}()
		}
	}
}
func find_social(name *byte) *social_messg {
	var (
		cmd    int
		socidx int
	)
	if (func() int {
		cmd = find_command(name)
		return cmd
	}()) < 0 {
		return nil
	}
	if (func() int {
		socidx = find_action(cmd)
		return socidx
	}()) < 0 {
		return nil
	}
	return &soc_mess_list[socidx]
}
func do_gmote(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		act_nr int
		length int
		arg    [2048]byte
		buf    [2048]byte
		action *social_messg
		vict   *char_data = nil
	)
	half_chop(argument, &buf[0], &arg[0])
	if subcmd != 0 {
		for func() int {
			length = libc.StrLen(&buf[0])
			return func() int {
				cmd = 0
				return cmd
			}()
		}(); *complete_cmd_info[cmd].Command != '\n'; cmd++ {
			if libc.StrNCmp(complete_cmd_info[cmd].Command, &buf[0], length) == 0 {
				break
			}
		}
	}
	if (func() int {
		act_nr = find_action(cmd)
		return act_nr
	}()) < 0 {
		stdio.Snprintf(&buf[0], int(2048), "@D[@BOOC@D: @g%s %s@n@D]", func() *byte {
			if ch.Admlevel < 1 {
				return ch.Desc.User
			}
			return GET_NAME(ch)
		}(), argument)
		act(&buf[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_GMOTE)
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_SOUNDPROOF) {
		send_to_char(ch, libc.CString("The walls seem to absorb your actions.\r\n"))
		return
	}
	action = &soc_mess_list[act_nr]
	if action.Char_found == nil {
		arg[0] = '\x00'
	}
	if arg[0] == 0 {
		if action.Others_no_arg == nil || *action.Others_no_arg == 0 {
			send_to_char(ch, libc.CString("Who are you going to do that to?\r\n"))
			return
		}
		stdio.Snprintf(&buf[0], int(2048), "@D[@BOOC@D: @g%s@D]@n", action.Others_no_arg)
	} else if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("%s\r\n"), action.Not_found)
		return
	} else if vict == ch {
		if action.Others_auto == nil || *action.Others_auto == 0 {
			send_to_char(ch, libc.CString("%s\r\n"), action.Char_auto)
			return
		}
		stdio.Snprintf(&buf[0], int(2048), "@D[@BOOC@D: @g%s@D]@n", action.Others_auto)
	} else {
		if int(vict.Position) < action.Min_victim_position {
			act(libc.CString("$N is not in a proper position for that."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_CHAR|2<<7))
			return
		}
		stdio.Snprintf(&buf[0], int(2048), "@D[@BOOC@D: @g%s@D]@n", action.Others_found)
	}
	act(&buf[0], FALSE, ch, nil, unsafe.Pointer(vict), TO_GMOTE)
}
