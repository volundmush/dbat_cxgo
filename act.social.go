package main

import "C"
import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
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
	action = (*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(act_nr)))
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
				targ = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
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
					if ch.Sex == SEX_MALE {
						if victim.Sex == SEX_MALE {
							act(libc.CString("$n accuses you of fighting like a woman!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						} else {
							act(libc.CString("$n says that women can't fight."), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
						}
					} else {
						if victim.Sex == SEX_MALE {
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
		fl           *C.FILE
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
		if (func() *C.FILE {
			fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_MISC, "r")))
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: can't open socials file '%s': %s"), LIB_MISC, C.strerror(*__errno_location()))
			C.exit(1)
		}
		next_soc[0] = '\x00'
		for C.feof(fl) == 0 {
			C.fgets(&next_soc[0], MAX_STRING_LENGTH, fl)
			if next_soc[0] == '~' {
				top_of_socialt++
			}
		}
	} else {
		if (func() *C.FILE {
			fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(LIB_MISC, "r")))
			return fl
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: can't open socials file '%s': %s"), LIB_MISC, C.strerror(*__errno_location()))
			C.exit(1)
		}
		for C.feof(fl) == 0 {
			C.fgets(&next_soc[0], MAX_STRING_LENGTH, fl)
			if next_soc[0] == '\n' || next_soc[0] == '\r' {
				top_of_socialt++
			}
		}
	}
	basic_mud_log(libc.CString("Social table contains %d socials."), top_of_socialt)
	C.rewind(fl)
	soc_mess_list = &make([]social_messg, top_of_socialt+1)[0]
	for {
		__isoc99_fscanf(fl, libc.CString(" %s "), &next_soc[0])
		if next_soc[0] == '$' {
			break
		}
		if config_info.Operation.Use_new_socials == TRUE {
			if __isoc99_fscanf(fl, libc.CString(" %s %d %d %d %d \n"), &sorted[0], &hide, &min_char_pos, &min_pos, &min_lvl) != 5 {
				basic_mud_log(libc.CString("SYSERR: format error in social file near social '%s'"), &next_soc[0])
				C.exit(1)
			}
			curr_soc++
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Command = C.strdup(&next_soc[1])
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Sort_as = C.strdup(&sorted[0])
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Hide = hide
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_char_position = min_char_pos
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_victim_position = min_pos
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_level_char = min_lvl
		} else {
			if __isoc99_fscanf(fl, libc.CString(" %d %d \n"), &hide, &min_pos) != 2 {
				basic_mud_log(libc.CString("SYSERR: format error in social file near social '%s'"), &next_soc[0])
				C.exit(1)
			}
			curr_soc++
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Command = C.strdup(&next_soc[0])
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Sort_as = C.strdup(&next_soc[0])
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Hide = hide
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_char_position = POS_RESTING
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_victim_position = min_pos
			(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Min_level_char = 0
		}
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_no_arg = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Others_no_arg = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_found = fread_action(fl, nr)
		if config_info.Operation.Use_new_socials == FALSE && (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_found == nil {
			continue
		}
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Others_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Vict_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Not_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_auto = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Others_auto = fread_action(fl, nr)
		if config_info.Operation.Use_new_socials == FALSE {
			continue
		}
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_body_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Others_body_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Vict_body_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Char_obj_found = fread_action(fl, nr)
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(curr_soc)))).Others_obj_found = fread_action(fl, nr)
	}
	C.fclose(fl)
	if curr_soc <= top_of_socialt {
	} else {
		__assert_fail(libc.CString("curr_soc <= top_of_socialt"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
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
			if C.strcasecmp((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))).Sort_as, (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(k)))).Sort_as) < 0 {
				k = i
			}
		}
		if j != k {
			temp = *(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j)))
			*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j))) = *(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(k)))
			*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(k))) = temp
		}
	}
	i = 0
	for *cmd_info[i].Command != '\n' {
		i++
	}
	i++
	complete_cmd_info = &make([]command_info, top_of_socialt+i+2)[0]
	i = 0
	j = 0
	k = 0
	for *cmd_info[i].Command != '\n' {
		*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(func() int {
			p := &k
			x := *p
			*p++
			return x
		}()))) = cmd_info[func() int {
			p := &i
			x := *p
			*p++
			return x
		}()]
	}
	for j <= top_of_socialt {
		(*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j)))).Act_nr = k
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Command = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j)))).Command
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Sort_as = (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j)))).Sort_as
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_position = int8((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(j)))).Min_char_position)
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Command_pointer = do_action
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_level = int16((*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(func() int {
			p := &j
			x := *p
			*p++
			return x
		}())))).Min_level_char)
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_admlevel = ADMLVL_NONE
		(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(func() int {
			p := &k
			x := *p
			*p++
			return x
		}())))).Subcmd = 0
	}
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Command = C.strdup(libc.CString("\n"))
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Sort_as = C.strdup(libc.CString("zzzzzzz"))
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_position = 0
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Command_pointer = nil
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_level = 0
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Minimum_admlevel = 0
	(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(k)))).Subcmd = 0
	basic_mud_log(libc.CString("Command info rebuilt, %d total commands."), k)
}
func free_command_list() {
	libc.Free(unsafe.Pointer(complete_cmd_info))
	complete_cmd_info = nil
}
func fread_action(fl *C.FILE, nr int) *byte {
	var buf [64936]byte
	C.fgets(&buf[0], MAX_STRING_LENGTH, fl)
	if C.feof(fl) != 0 {
		basic_mud_log(libc.CString("SYSERR: fread_action: unexpected EOF near action #%d"), nr)
		C.exit(1)
	}
	if buf[0] == '#' {
		return nil
	}
	buf[C.strlen(&buf[0])-1] = '\x00'
	return C.strdup(&buf[0])
}
func free_social_messages() {
	var (
		mess *social_messg
		i    int
	)
	for i = 0; i <= top_of_socialt; i++ {
		mess = (*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(i)))
		free_action(mess)
	}
	libc.Free(unsafe.Pointer(soc_mess_list))
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
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(mid)))).Act_nr == cmd {
			return mid
		}
		if bot >= top {
			return -1
		}
		if (*(*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(mid)))).Act_nr > cmd {
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
	return (*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(socidx)))
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
			length = int(C.strlen(&buf[0]))
			return func() int {
				cmd = 0
				return cmd
			}()
		}(); *(*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command != '\n'; cmd++ {
			if C.strncmp((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command, &buf[0], uint64(length)) == 0 {
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
	action = (*social_messg)(unsafe.Add(unsafe.Pointer(soc_mess_list), unsafe.Sizeof(social_messg{})*uintptr(act_nr)))
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
