package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func dump(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		k     *obj_data
		value int = 0
	)
	for k = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; k != nil; k = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents {
		act(libc.CString("$p vanishes in a puff of smoke!"), FALSE, nil, k, nil, TO_ROOM)
		extract_obj(k)
	}
	if libc.StrCmp(libc.CString("drop"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) != 0 {
		return FALSE
	}
	do_drop(ch, argument, cmd, SCMD_DROP)
	for k = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; k != nil; k = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents {
		act(libc.CString("$p vanishes in a puff of smoke!"), FALSE, nil, k, nil, TO_ROOM)
		value += MAX(1, MIN(50, k.Cost/10))
		extract_obj(k)
	}
	if value != 0 {
		send_to_char(ch, libc.CString("You are awarded for outstanding performance.\r\n"))
		act(libc.CString("$n has been awarded for being a good citizen."), TRUE, ch, nil, nil, TO_ROOM)
		if GET_LEVEL(ch) < 3 {
			gain_exp(ch, int64(value))
		} else {
			ch.Gold += value
		}
	}
	return TRUE
}
func mayor(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		actbuf    [2048]byte
		open_path [53]byte = func() [53]byte {
			var t [53]byte
			copy(t[:], []byte("W3a3003b33000c111d0d111Oe333333Oe22c222112212111a1S."))
			return t
		}()
		close_path [53]byte = func() [53]byte {
			var t [53]byte
			copy(t[:], []byte("W3a3003b33000c111d0d111CE333333CE22c222112212111a1S."))
			return t
		}()
		path       *byte = nil
		path_index int
		move       bool = FALSE != 0
	)
	if !move {
		if time_info.Hours == 6 {
			move = TRUE != 0
			path = &open_path[0]
			path_index = 0
		} else if time_info.Hours == 20 {
			move = TRUE != 0
			path = &close_path[0]
			path_index = 0
		}
	}
	if cmd != 0 || !move || int(ch.Position) < POS_SLEEPING || int(ch.Position) == POS_FIGHTING {
		return FALSE
	}
	switch *(*byte)(unsafe.Add(unsafe.Pointer(path), path_index)) {
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		perform_move(ch, int(*(*byte)(unsafe.Add(unsafe.Pointer(path), path_index))-'0'), 1)
	case 'W':
		ch.Position = POS_STANDING
		act(libc.CString("$n awakens and groans loudly."), FALSE, ch, nil, nil, TO_ROOM)
	case 'S':
		ch.Position = POS_SLEEPING
		act(libc.CString("$n lies down and instantly falls asleep."), FALSE, ch, nil, nil, TO_ROOM)
	case 'a':
		act(libc.CString("$n says 'Hello Honey!'"), FALSE, ch, nil, nil, TO_ROOM)
		act(libc.CString("$n smirks."), FALSE, ch, nil, nil, TO_ROOM)
	case 'b':
		act(libc.CString("$n says 'What a view!  I must get something done about that dump!'"), FALSE, ch, nil, nil, TO_ROOM)
	case 'c':
		act(libc.CString("$n says 'Vandals!  Youngsters nowadays have no respect for anything!'"), FALSE, ch, nil, nil, TO_ROOM)
	case 'd':
		act(libc.CString("$n says 'Good day, citizens!'"), FALSE, ch, nil, nil, TO_ROOM)
	case 'e':
		act(libc.CString("$n says 'I hereby declare the bazaar open!'"), FALSE, ch, nil, nil, TO_ROOM)
	case 'E':
		act(libc.CString("$n says 'I hereby declare Midgaard closed!'"), FALSE, ch, nil, nil, TO_ROOM)
	case 'O':
		do_gen_door(ch, libc.StrCpy(&actbuf[0], libc.CString("gate")), 0, SCMD_UNLOCK)
		do_gen_door(ch, libc.StrCpy(&actbuf[0], libc.CString("gate")), 0, SCMD_OPEN)
	case 'C':
		do_gen_door(ch, libc.StrCpy(&actbuf[0], libc.CString("gate")), 0, SCMD_CLOSE)
		do_gen_door(ch, libc.StrCpy(&actbuf[0], libc.CString("gate")), 0, SCMD_LOCK)
	case '.':
		move = FALSE != 0
	}
	path_index++
	return FALSE
}
func num_players_in_room(room room_vnum) int {
	var (
		i           *descriptor_data
		num_players int = 0
	)
	for i = descriptor_list; i != nil; i = i.Next {
		if i.Connected != CON_PLAYING {
			continue
		}
		if i.Character == nil {
			continue
		}
		if i.Character.In_room == room_rnum(-1) || i.Character.In_room > top_of_world {
			continue
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.Character.In_room)))).Number != room {
			continue
		}
		if i.Character.Admlevel >= ADMLVL_IMMORT && PRF_FLAGGED(i.Character, PRF_NOHASSLE) {
			continue
		}
		num_players++
	}
	return num_players
}
func check_mob_in_room(mob mob_vnum, room room_vnum) bool {
	var (
		i     *char_data
		found bool = FALSE != 0
	)
	for i = character_list; i != nil; i = i.Next {
		if GET_MOB_VNUM(i) == mob {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Number == room {
				found = TRUE != 0
			}
		}
	}
	return found
}
func check_obj_in_room(obj obj_vnum, room room_vnum) bool {
	var (
		i      *obj_data
		list   *obj_data
		found  bool = FALSE != 0
		r_room room_rnum
	)
	r_room = real_room(room)
	list = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(r_room)))).Contents
	for i = list; i != nil; i = i.Next_content {
		if GET_OBJ_VNUM(i) == obj {
			found = TRUE != 0
		}
	}
	return found
}

var gauntlet_info [42][3]int = [42][3]int{{0, 2403, SCMD_SOUTH}, {1, 2404, SCMD_SOUTH}, {2, 2405, SCMD_SOUTH}, {3, 2406, SCMD_SOUTH}, {4, 2407, SCMD_SOUTH}, {5, 2408, SCMD_SOUTH}, {6, 2409, SCMD_SOUTH}, {7, 2410, SCMD_SOUTH}, {8, 2411, SCMD_SOUTH}, {9, 2412, SCMD_SOUTH}, {10, 2413, SCMD_SOUTH}, {11, 2414, SCMD_SOUTH}, {12, 2415, SCMD_SOUTH}, {13, 2416, SCMD_SOUTH}, {14, 2417, SCMD_SOUTH}, {15, 2418, SCMD_SOUTH}, {16, 2420, SCMD_SOUTH}, {17, 2421, SCMD_SOUTH}, {18, 2422, SCMD_SOUTH}, {19, 2423, SCMD_SOUTH}, {20, 2424, SCMD_SOUTH}, {21, 2425, SCMD_SOUTH}, {22, 2426, SCMD_SOUTH}, {23, 2427, SCMD_SOUTH}, {24, 2428, SCMD_SOUTH}, {25, 2429, SCMD_SOUTH}, {26, 2430, SCMD_SOUTH}, {27, 2431, SCMD_SOUTH}, {28, 2432, SCMD_SOUTH}, {29, 2433, SCMD_SOUTH}, {30, 2434, SCMD_SOUTH}, {31, 2435, SCMD_SOUTH}, {32, 2436, SCMD_SOUTH}, {33, 2437, SCMD_SOUTH}, {34, 2438, SCMD_SOUTH}, {35, 2439, SCMD_SOUTH}, {36, 2440, SCMD_SOUTH}, {37, 2441, SCMD_SOUTH}, {38, 2442, SCMD_SOUTH}, {39, 2443, SCMD_SOUTH}, {40, 2444, SCMD_SOUTH}, {-1, -1, -1}}

func gauntlet_room(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		i       int = 0
		proceed int = 1
		tch     *char_data
		buf2    *byte = libc.CString("$N tried to sneak past without a fight, and got nowhere.")
		buf     [64936]byte
		nomob   bool = TRUE != 0
	)
	for i = 0; gauntlet_info[i][0] != -1; i++ {
		if !IS_NPC(ch) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number == room_vnum(gauntlet_info[i][1]) {
			if ch.Player_specials.Gauntlet < (gauntlet_info[i][0]) {
				ch.Player_specials.Gauntlet = gauntlet_info[i][0]
			}
		}
	}
	if cmd == 0 {
		return FALSE
	}
	if libc.StrCmp(libc.CString("flee"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		send_to_char(ch, libc.CString("Fleeing is not allowed!  If you want to get out of here, type @Ysurrender@n while fighting to be returned to the start."))
		return TRUE
	}
	if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command_pointer) != libc.FuncAddr(do_move) && libc.StrCmp(libc.CString("surrender"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) != 0 {
		return FALSE
	}
	if IS_NPC(ch) {
		return FALSE
	}
	if libc.StrCmp(libc.CString("surrender"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if ch.Fighting != nil {
			char_from_room(ch)
			char_to_room(ch, real_room(room_vnum(gauntlet_info[0][1])))
			act(libc.CString("$n suddenly appears looking relieved after $s trial in the Gauntlet"), FALSE, ch, nil, unsafe.Pointer(ch), TO_NOTVICT)
			act(libc.CString("You are returned to the start of the Gauntlet"), FALSE, ch, nil, unsafe.Pointer(ch), TO_VICT)
			if ch.Hit > 2000 {
				ch.Hit = ch.Hit / 5
			} else if ch.Hit > 500 {
				ch.Hit = 100
			} else {
				ch.Hit = 1
			}
			look_at_room(ch.In_room, ch, 0)
			return TRUE
		} else {
			send_to_char(ch, libc.CString("You can only surrender while fighting, so at least TRY to make an effort"))
			return TRUE
		}
	}
	if ch.Admlevel >= ADMLVL_IMMORT {
		return FALSE
	}
	for i = 0; gauntlet_info[i][0] != -1; i++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number == room_vnum(gauntlet_info[i][1]) {
			if cmd == gauntlet_info[i][2] {
				for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = tch.Next_in_room {
					if IS_NPC(tch) && i > 0 {
						proceed = 0
						stdio.Sprintf(&buf[0], "%s wants to teach you a lesson first.\r\n", GET_NAME(tch))
					}
				}
				if proceed != 0 {
					nomob = TRUE != 0
					for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(room_vnum(gauntlet_info[i+1][1])))))).People; tch != nil; tch = tch.Next_in_room {
						if !IS_NPC(tch) {
							proceed = 0
							stdio.Sprintf(&buf[0], "%s is in the next room.  You must wait for them to finish.\r\n", GET_NAME(tch))
						} else {
							nomob = FALSE != 0
						}
					}
					if int(libc.BoolToInt(nomob)) == TRUE {
						proceed = 0
						stdio.Sprintf(&buf[0], "The next room is empty.  You must wait for your opponent to re-appear.\r\n")
					}
				}
				if proceed == 0 {
					send_to_char(ch, &buf[0])
					act(buf2, FALSE, ch, nil, unsafe.Pointer(ch), TO_ROOM)
					return TRUE
				}
			}
		}
	}
	return FALSE
}
func gauntlet_end(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var i int = 0
	if !IS_NPC(ch) {
		if ch.Player_specials.Gauntlet < GAUNTLET_END {
			ch.Player_specials.Gauntlet = GAUNTLET_END
		}
	}
	if cmd == 0 {
		return FALSE
	}
	if libc.StrCmp(libc.CString("flee"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if ch.Fighting != nil && int(ch.Position) == POS_FIGHTING {
			send_to_char(ch, libc.CString("You can't flee from this fight./r/nIt's your own fault for summoning creatures into the gauntlet!\r\n"))
			return TRUE
		} else {
			send_to_char(ch, libc.CString("There is nothing here to flee from\r\n"))
			return TRUE
		}
	}
	if libc.StrCmp(libc.CString("surrender"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		send_to_char(ch, libc.CString("You have completed the gauntlet, why would you need to surrender?\r\n"))
		return TRUE
	}
	if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command_pointer) != libc.FuncAddr(do_move) {
		return FALSE
	}
	if IS_NPC(ch) {
		return FALSE
	}
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[cmd-1]) == nil || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[cmd-1]).To_room == room_rnum(-1) {
		return FALSE
	}
	if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[cmd-1], 1<<1) {
		return FALSE
	}
	for i = 0; gauntlet_info[i][0] != -1; i++ {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[cmd-1]).To_room)))).Number == room_vnum(gauntlet_info[i][1]) {
			send_to_char(ch, libc.CString("You have completed the gauntlet, you cannot go backwards!\r\n"))
			return TRUE
		}
	}
	return FALSE
}
func gauntlet_rest(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		i       int = 0
		proceed int = 1
		door    int
		tch     *char_data
		buf2    *byte = libc.CString("$N tried to return to the gauntlet, and got nowhere.")
		buf     [64936]byte
		nomob   bool = TRUE != 0
	)
	_ = nomob
	if cmd == 0 {
		return FALSE
	}
	if libc.StrCmp(libc.CString("flee"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		send_to_char(ch, libc.CString("Fleeing is not allowed!  If you want to get out of here, type @Ysurrender@n while fighting to be returned to the start."))
		return TRUE
	}
	if libc.StrCmp(libc.CString("surrender"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		send_to_char(ch, libc.CString("You are in a rest-room.  Surrender is not an option.\r\nIf you want to leave the Gauntlet, you can surrender while fighting.\r\n"))
		return TRUE
	}
	if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command_pointer) != libc.FuncAddr(do_move) {
		return FALSE
	}
	if IS_NPC(ch) {
		return FALSE
	}
	if ch.Admlevel >= ADMLVL_IMMORT {
		return FALSE
	}
	for i = 0; gauntlet_info[i][0] != -1; i++ {
		for door = 0; door < NUM_OF_DIRS; door++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) == nil || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room == room_rnum(-1) {
				continue
			}
			if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1) {
				continue
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room)))).Number == room_vnum(gauntlet_info[i][1]) && door == (cmd-1) {
				nomob = TRUE != 0
				for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(room_vnum(gauntlet_info[i][1])))))).People; tch != nil; tch = tch.Next_in_room {
					if !IS_NPC(tch) {
						proceed = 0
						stdio.Sprintf(&buf[0], "%s has moved into the next room.  You must wait for them to finish.\r\n", GET_NAME(tch))
					} else {
						nomob = FALSE != 0
					}
				}
				if proceed == 0 {
					send_to_char(ch, &buf[0])
					act(buf2, FALSE, ch, nil, unsafe.Pointer(ch), TO_ROOM)
					return TRUE
				}
			}
		}
	}
	return FALSE
}
func npc_steal(ch *char_data, victim *char_data) {
	var gold int
	if IS_NPC(victim) {
		return
	}
	if IS_NPC(ch) {
		return
	}
	if ADM_FLAGGED(victim, ADM_NOSTEAL) {
		return
	}
	if !CAN_SEE(ch, victim) {
		return
	}
	if AWAKE(victim) && rand_number(0, GET_LEVEL(ch)) == 0 {
		act(libc.CString("You discover that $n has $s hands in your wallet."), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
		act(libc.CString("$n tries to steal zenni from $N."), TRUE, ch, nil, unsafe.Pointer(victim), TO_NOTVICT)
	} else {
		gold = (victim.Gold * rand_number(1, 10)) / 100
		if gold > 0 {
			ch.Gold += gold
			victim.Gold -= gold
		}
	}
}
func snake(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	if cmd != 0 || int(ch.Position) != POS_FIGHTING || ch.Fighting == nil {
		return FALSE
	}
	if ch.Fighting.In_room != ch.In_room || rand_number(0, GET_LEVEL(ch)) != 0 {
		return FALSE
	}
	act(libc.CString("$n bites $N!"), 1, ch, nil, unsafe.Pointer(ch.Fighting), TO_NOTVICT)
	act(libc.CString("$n bites you!"), 1, ch, nil, unsafe.Pointer(ch.Fighting), TO_VICT)
	call_magic(ch, ch.Fighting, nil, SPELL_POISON, GET_LEVEL(ch), CAST_SPELL, nil)
	return TRUE
}
func thief(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var cons *char_data
	if IS_NPC(ch) {
		return FALSE
	}
	if cmd != 0 || int(ch.Position) != POS_STANDING {
		return FALSE
	}
	for cons = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; cons != nil; cons = cons.Next_in_room {
		if !IS_NPC(cons) && !ADM_FLAGGED(cons, ADM_NOSTEAL) && rand_number(0, 4) == 0 {
			npc_steal(ch, cons)
			return TRUE
		}
	}
	return FALSE
}
func magic_user_orig(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var vict *char_data
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil && ch.Fighting.In_room == ch.In_room {
		vict = ch.Fighting
	}
	if vict == nil {
		return TRUE
	}
	if GET_LEVEL(ch) > 13 && rand_number(0, 10) == 0 {
		cast_spell(ch, vict, nil, SPELL_POISON, nil)
	}
	if GET_LEVEL(ch) > 7 && rand_number(0, 8) == 0 {
		cast_spell(ch, vict, nil, SPELL_BLINDNESS, nil)
	}
	if GET_LEVEL(ch) > 12 && rand_number(0, 12) == 0 {
		if IS_EVIL(ch) {
			cast_spell(ch, vict, nil, SPELL_ENERGY_DRAIN, nil)
		} else if IS_GOOD(ch) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_EVIL, nil)
		}
	}
	if rand_number(0, 4) != 0 {
		return TRUE
	}
	switch GET_LEVEL(ch) {
	case 4:
		fallthrough
	case 5:
		cast_spell(ch, vict, nil, SPELL_MAGIC_MISSILE, nil)
	case 6:
		fallthrough
	case 7:
		cast_spell(ch, vict, nil, SPELL_CHILL_TOUCH, nil)
	case 8:
		fallthrough
	case 9:
		cast_spell(ch, vict, nil, SPELL_BURNING_HANDS, nil)
	case 10:
		fallthrough
	case 11:
		cast_spell(ch, vict, nil, SPELL_SHOCKING_GRASP, nil)
	case 12:
		fallthrough
	case 13:
		cast_spell(ch, vict, nil, SPELL_LIGHTNING_BOLT, nil)
	case 14:
		fallthrough
	case 15:
		fallthrough
	case 16:
		fallthrough
	case 17:
		cast_spell(ch, vict, nil, SPELL_COLOR_SPRAY, nil)
	default:
		cast_spell(ch, vict, nil, SPELL_FIREBALL, nil)
	}
	return TRUE
}
func guild_guard(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		i     int
		guard *char_data = (*char_data)(me)
		buf   *byte      = libc.CString("The guard humiliates you, and blocks your way.\r\n")
		buf2  *byte      = libc.CString("The guard humiliates $n, and blocks $s way.")
	)
	if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command_pointer) != libc.FuncAddr(do_move) || AFF_FLAGGED(guard, AFF_BLIND) {
		return FALSE
	}
	if ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		return FALSE
	}
	for i = 0; guild_info[i].Guild_room != room_vnum(-1); i++ {
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != guild_info[i].Guild_room || cmd != guild_info[i].Direction {
			continue
		}
		if !IS_NPC(ch) && ((ch.Chclasses[guild_info[i].Pc_class])+(ch.Epicclasses[guild_info[i].Pc_class])) > 0 {
			continue
		}
		send_to_char(ch, libc.CString("%s"), buf)
		act(buf2, FALSE, ch, nil, nil, TO_ROOM)
		return TRUE
	}
	return FALSE
}
func puff(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var actbuf [2048]byte
	if cmd != 0 {
		return FALSE
	}
	switch rand_number(0, 60) {
	case 0:
		do_say(ch, libc.StrCpy(&actbuf[0], libc.CString("My god!  It's full of stars!")), 0, 0)
		return TRUE
	case 1:
		do_say(ch, libc.StrCpy(&actbuf[0], libc.CString("How'd all those fish get up here?")), 0, 0)
		return TRUE
	case 2:
		do_say(ch, libc.StrCpy(&actbuf[0], libc.CString("I'm a very female dragon.")), 0, 0)
		return TRUE
	case 3:
		do_say(ch, libc.StrCpy(&actbuf[0], libc.CString("I've got a peaceful, easy feeling.")), 0, 0)
		return TRUE
	default:
		return FALSE
	}
}
func fido(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		i        *obj_data
		temp     *obj_data
		next_obj *obj_data
	)
	if cmd != 0 || !AWAKE(ch) {
		return FALSE
	}
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
		if !IS_CORPSE(i) {
			continue
		}
		act(libc.CString("$n savagely devours a corpse."), FALSE, ch, nil, nil, TO_ROOM)
		for temp = i.Contains; temp != nil; temp = next_obj {
			next_obj = temp.Next_content
			obj_from_obj(temp)
			obj_to_room(temp, ch.In_room)
		}
		extract_obj(i)
		return TRUE
	}
	return FALSE
}
func janitor(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var i *obj_data
	if cmd != 0 || !AWAKE(ch) {
		return FALSE
	}
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
		if !OBJWEAR_FLAGGED(i, ITEM_WEAR_TAKE) {
			continue
		}
		if int(i.Type_flag) == ITEM_DRINKCON || i.Cost >= 100 {
			continue
		}
		act(libc.CString("$n picks up some trash."), FALSE, ch, nil, nil, TO_ROOM)
		obj_from_room(i)
		obj_to_char(i, ch)
		return TRUE
	}
	return FALSE
}
func cityguard(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tch      *char_data
		evil     *char_data
		spittle  *char_data
		max_evil int
		min_cha  int
	)
	if cmd != 0 || !AWAKE(ch) || ch.Fighting != nil {
		return FALSE
	}
	max_evil = 1000
	min_cha = 6
	spittle = func() *char_data {
		evil = nil
		return evil
	}()
	for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = tch.Next_in_room {
		if !CAN_SEE(ch, tch) {
			continue
		}
		if !IS_NPC(tch) && PLR_FLAGGED(tch, PLR_KILLER) {
			act(libc.CString("$n screams 'HEY!!!  You're one of those PLAYER KILLERS!!!!!!'"), FALSE, ch, nil, nil, TO_ROOM)
			return TRUE
		}
		if !IS_NPC(tch) && PLR_FLAGGED(tch, PLR_THIEF) {
			act(libc.CString("$n screams 'HEY!!!  You're one of those PLAYER THIEVES!!!!!!'"), FALSE, ch, nil, nil, TO_ROOM)
			return TRUE
		}
		if tch.Fighting != nil && tch.Alignment < max_evil && (IS_NPC(tch) || IS_NPC(tch.Fighting)) {
			max_evil = tch.Alignment
			evil = tch
		}
		if int(tch.Aff_abils.Cha) < min_cha {
			spittle = tch
			min_cha = int(tch.Aff_abils.Cha)
		}
	}
	if evil != nil && evil.Fighting.Alignment >= 0 {
		act(libc.CString("$n screams 'PROTECT THE INNOCENT!  BANZAI!  CHARGE!  ARARARAGGGHH!'"), FALSE, ch, nil, nil, TO_ROOM)
		return TRUE
	}
	if spittle != nil && rand_number(0, 9) == 0 {
		var spit_social int
		if spit_social == 0 {
			spit_social = find_command(libc.CString("spit"))
		}
		if spit_social > 0 {
			var spitbuf [21]byte
			libc.StrNCpy(&spitbuf[0], GET_NAME(spittle), int(21))
			spitbuf[21-1] = '\x00'
			do_action(ch, &spitbuf[0], spit_social, 0)
			return TRUE
		}
	}
	return FALSE
}
func pet_shops(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		buf      [64936]byte
		pet_name [256]byte
		pet_room room_rnum
		pet      *char_data
	)
	pet_room = ch.In_room + 1
	if libc.StrCmp(libc.CString("list"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		send_to_char(ch, libc.CString("Available pets are:\r\n"))
		for pet = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(pet_room)))).People; pet != nil; pet = pet.Next_in_room {
			if !IS_NPC(pet) {
				continue
			}
			send_to_char(ch, libc.CString("%8d - %s\r\n"), GET_LEVEL(pet)*300, GET_NAME(pet))
		}
		return TRUE
	} else if libc.StrCmp(libc.CString("buy"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		two_arguments(argument, &buf[0], &pet_name[0])
		if (func() *char_data {
			pet = get_char_room(&buf[0], nil, pet_room)
			return pet
		}()) == nil || !IS_NPC(pet) {
			send_to_char(ch, libc.CString("There is no such pet!\r\n"))
			return TRUE
		}
		if ch.Gold < (GET_LEVEL(pet) * 300) {
			send_to_char(ch, libc.CString("You don't have enough zenni!\r\n"))
			return TRUE
		}
		ch.Gold -= GET_LEVEL(pet) * 300
		pet = read_mobile(mob_vnum(pet.Nr), REAL)
		pet.Exp = 0
		pet.Affected_by[int(AFF_CHARM/32)] |= 1 << (int(AFF_CHARM % 32))
		if pet_name[0] != 0 {
			stdio.Snprintf(&buf[0], int(64936), "%s %s", pet.Name, &pet_name[0])
			pet.Name = libc.StrDup(&buf[0])
			stdio.Snprintf(&buf[0], int(64936), "%sA small sign on a chain around the neck says 'My name is %s'\r\n", pet.Description, &pet_name[0])
			pet.Description = libc.StrDup(&buf[0])
		}
		char_to_room(pet, ch.In_room)
		add_follower(pet, ch)
		pet.Master_id = ch.Idnum
		pet.Carry_weight = 1000
		pet.Carry_items = 100
		send_to_char(ch, libc.CString("May you enjoy your pet.\r\n"))
		act(libc.CString("$n buys $N as a pet."), FALSE, ch, nil, unsafe.Pointer(pet), TO_ROOM)
		return TRUE
	}
	return FALSE
}
func auction(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		auct_room room_rnum
		obj       *obj_data
		next_obj  *obj_data
		obj2      *obj_data = nil
		found     int       = FALSE
	)
	auct_room = real_room(80)
	if libc.StrCmp(libc.CString("cancel"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(auct_room)))).Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil && int(obj.Aucter) == int(ch.Id) {
				obj2 = obj
				found = TRUE
				if int(obj2.CurBidder) != -1 && obj2.AucTime+0x7E900 > libc.GetTime(nil) {
					send_to_char(ch, libc.CString("Unable to cancel. Someone has already bid on it and their bid license hasn't expired.\r\n"))
					var remain libc.Time = (obj2.AucTime + 0x7E900) - libc.GetTime(nil)
					var day int = int((remain % 604800) / 86400)
					var hour int = int((remain % 86400) / 3600)
					var minu int = int((remain % 3600) / 60)
					send_to_char(ch, libc.CString("Time Till License Expiration: %d day%s, %d hour%s, %d minute%s.\r\n"), day, func() string {
						if day > 1 {
							return "s"
						}
						return ""
					}(), hour, func() string {
						if hour > 1 {
							return "s"
						}
						return ""
					}(), minu, func() string {
						if minu > 1 {
							return "s"
						}
						return ""
					}())
					continue
				}
				send_to_char(ch, libc.CString("@wYou cancel the auction of %s@w and it is returned to you.@n\r\n"), obj2.Short_description)
				var d *descriptor_data
				for d = descriptor_list; d != nil; d = d.Next {
					if d.Connected != CON_PLAYING || IS_NPC(d.Character) {
						continue
					}
					if d.Character == ch {
						continue
					}
					if (d.Character.Equipment[WEAR_EYE]) != nil {
						send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GThe auction of @w%s@G has been canceled.\r\n"), obj2.Short_description)
					}
				}
				obj_from_room(obj2)
				obj_to_char(obj2, ch)
				auc_save()
			}
		}
		if found == FALSE {
			send_to_char(ch, libc.CString("There are no items being auctioned by you.\r\n"))
		}
		return TRUE
	} else if libc.StrCmp(libc.CString("pickup"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		var (
			d       *descriptor_data
			founded int = FALSE
		)
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(auct_room)))).Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj != nil && int(obj.CurBidder) == int(ch.Id) {
				obj2 = obj
				found = TRUE
				if int(obj.Aucter) <= 0 {
					continue
				}
				if obj2.Bid > ch.Gold {
					send_to_char(ch, libc.CString("Unable to purchase %s, you don't have enough money on hand.\r\n"), obj2.Short_description)
					continue
				}
				if obj2.AucTime+86400 > libc.GetTime(nil) {
					var (
						remain libc.Time = (obj2.AucTime + 86400) - libc.GetTime(nil)
						hour   int       = int((remain % 86400) / 3600)
						minu   int       = int((remain % 3600) / 60)
					)
					send_to_char(ch, libc.CString("Unable to purchase %s, minimum time to bid is 24 hours. %d hour%s and %d minute%s remain.\r\n"), obj2.Short_description, hour, func() string {
						if hour > 1 {
							return "s"
						}
						return ""
					}(), minu, func() string {
						if minu > 1 {
							return "s"
						}
						return ""
					}())
					continue
				}
				ch.Gold -= obj2.Bid
				obj_from_room(obj2)
				obj_to_char(obj2, ch)
				send_to_char(ch, libc.CString("You pay %s zenni and receive the item.\r\n"), add_commas(int64(obj2.Bid)))
				auc_save()
				for d = descriptor_list; d != nil; d = d.Next {
					if d.Connected != CON_PLAYING || IS_NPC(d.Character) {
						continue
					}
					if d.Character == ch {
						continue
					}
					if int(d.Character.Idnum) == int(obj2.Aucter) {
						founded = TRUE
						d.Character.Bank_gold += obj2.Bid
						if (d.Character.Equipment[WEAR_EYE]) != nil {
							send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GSomeone has purchased your @w%s@G and you had the money put in your bank account.\r\n"), obj2.Short_description)
						}
					} else if (d.Character.Equipment[WEAR_EYE]) != nil {
						send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GSomeone has purchased the @w%s@G that was on auction.\r\n"), obj2.Short_description)
					}
				}
				if founded == FALSE {
					var (
						vict     *char_data = nil
						is_file  int        = FALSE
						player_i int        = 0
					)
					vict = new(char_data)
					clear_char(vict)
					vict.Player_specials = new(player_special_data)
					var blam [50]byte
					stdio.Sprintf(&blam[0], "%s", obj2.Auctname)
					if (func() int {
						player_i = load_char(&blam[0], vict)
						return player_i
					}()) > -1 {
						is_file = TRUE
					} else {
						free_char(vict)
						continue
					}
					vict.Bank_gold += obj2.Bid
					vict.Pfilepos = player_i
					save_char(vict)
					if is_file == TRUE {
						free_char(vict)
					}
				}
			}
		}
		if found == FALSE {
			send_to_char(ch, libc.CString("There are no items that you have bid on.\r\n"))
		}
		return TRUE
	} else if libc.StrCmp(libc.CString("auction"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		var (
			arg   [2048]byte
			arg2  [2048]byte
			d     *descriptor_data
			value int = 0
		)
		two_arguments(argument, &arg[0], &arg2[0])
		if arg[0] == 0 || arg2[0] == 0 {
			send_to_char(ch, libc.CString("Auction what item and for how much?\r\n"))
			return TRUE
		}
		value = libc.Atoi(libc.GoString(&arg2[0]))
		if (func() *obj_data {
			obj2 = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
			return obj2
		}()) == nil {
			send_to_char(ch, libc.CString("You don't have that item to auction.\r\n"))
			return TRUE
		}
		if value <= 999 {
			send_to_char(ch, libc.CString("Do not auction anything for less than 1,000 zenni.\r\n"))
			return TRUE
		}
		if OBJ_FLAGGED(obj2, ITEM_BROKEN) {
			act(libc.CString("$P is broken and we will not accept it."), FALSE, ch, nil, unsafe.Pointer(obj2), TO_CHAR)
			return TRUE
		}
		if OBJ_FLAGGED(obj2, ITEM_NODONATE) {
			act(libc.CString("$P is junk and we will not accept it."), FALSE, ch, nil, unsafe.Pointer(obj2), TO_CHAR)
			return TRUE
		}
		obj2.Bid = value
		obj2.Startbid = 0
		obj2.Aucter = 0
		if obj2.Auctname != nil {
			libc.Free(unsafe.Pointer(obj2.Auctname))
		}
		obj2.AucTime = 0
		obj2.Bid = value
		obj2.Startbid = obj2.Bid
		obj2.Aucter = ch.Id
		obj2.Auctname = libc.StrDup(GET_NAME(ch))
		obj2.AucTime = libc.GetTime(nil)
		obj2.CurBidder = -1
		obj_from_char(obj2)
		obj_to_room(obj2, auct_room)
		auc_save()
		send_to_char(ch, libc.CString("You place %s on auction for %s zenni.\r\n"), obj2.Short_description, add_commas(int64(obj2.Bid)))
		basic_mud_log(libc.CString("AUCTION: %s places %s on auction for %s"), GET_NAME(ch), obj2.Short_description, add_commas(int64(obj2.Bid)))
		for d = descriptor_list; d != nil; d = d.Next {
			if d.Connected != CON_PLAYING || IS_NPC(d.Character) {
				continue
			}
			if d.Character == ch {
				continue
			}
			if (d.Character.Equipment[WEAR_EYE]) != nil {
				send_to_char(d.Character, libc.CString("@RScouter Auction News@D: @GThe item, @w%s@G, has been placed on auction for @Y%s@G zenni.@n\r\n"), obj2.Short_description, add_commas(int64(obj2.Bid)))
			}
		}
		return TRUE
	}
	return FALSE
}
func healtank(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		htank *obj_data = nil
		i     *obj_data
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
		if GET_OBJ_VNUM(i) == 65 {
			htank = i
		} else {
			continue
		}
	}
	if libc.StrCmp(libc.CString("htank"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if htank == nil {
			return FALSE
		}
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("@WHealing Tank Commands:\r\nhtank [ enter | exit | check ]@n"))
			return TRUE
		}
		if libc.StrCaseCmp(libc.CString("enter"), &arg[0]) == 0 {
			if PLR_FLAGGED(ch, PLR_HEALT) {
				send_to_char(ch, libc.CString("You are already inside a healing tank!\r\n"))
				return TRUE
			}
			if ch.Master != nil && ch.Master != ch {
				send_to_char(ch, libc.CString("You can't enter it while following someone!\r\n"))
				return TRUE
			} else if int(ch.Race) == RACE_ANDROID {
				send_to_char(ch, libc.CString("A healing tank will have no effect on you.\r\n"))
				return TRUE
			} else if htank.Healcharge <= 0 {
				send_to_char(ch, libc.CString("That healing tank needs to recharge, wait a while.\r\n"))
				return TRUE
			} else if OBJ_FLAGGED(htank, ITEM_BROKEN) {
				send_to_char(ch, libc.CString("It is broken! You will need to fix it yourself or wait for someone else to fix it.\r\n"))
				return TRUE
			} else if ch.Sits != nil {
				send_to_char(ch, libc.CString("You are already on something.\r\n"))
				return TRUE
			} else if htank.Sitting != nil {
				send_to_char(ch, libc.CString("Someone else is already inside that healing tank!\r\n"))
				return TRUE
			} else {
				ch.Charge = 0
				ch.Act[int(PLR_CHARGE/32)] &= bitvector_t(int32(^(1 << (int(PLR_CHARGE % 32)))))
				ch.Chargeto = 0
				ch.Barrier = 0
				act(libc.CString("@wYou step inside the healing tank and put on its breathing mask. A water like solution pours over your body until the tank is full.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@w steps inside the healing tank and puts on its breathing mask. A water like solution pours over $s body until the tank is full.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Act[int(PLR_HEALT/32)] |= bitvector_t(int32(1 << (int(PLR_HEALT % 32))))
				ch.Sits = htank
				htank.Sitting = ch
				return TRUE
			}
		} else if libc.StrCaseCmp(libc.CString("exit"), &arg[0]) == 0 {
			if !PLR_FLAGGED(ch, PLR_HEALT) {
				send_to_char(ch, libc.CString("You are not inside a healing tank.\r\n"))
				return TRUE
			} else {
				act(libc.CString("@wThe healing tank drains and you exit it shortly after."), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n@w exits the healing tank after letting it drain.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Act[int(PLR_HEALT/32)] &= bitvector_t(int32(^(1 << (int(PLR_HEALT % 32)))))
				htank.Sitting = nil
				ch.Sits = nil
				return TRUE
			}
		} else if libc.StrCaseCmp(libc.CString("check"), &arg[0]) == 0 {
			if htank.Healcharge < 20 && htank.Healcharge > 0 {
				send_to_char(ch, libc.CString("The healing tank has %d bars of energy displayed on its meter.\r\n"), htank.Healcharge)
			} else if htank.Healcharge <= 0 {
				send_to_char(ch, libc.CString("The healing tank has no energy displayed on its meter.\r\n"))
			} else {
				send_to_char(ch, libc.CString("The healing tank has full energy shown on its meter.\r\n"))
			}
			return TRUE
		} else {
			send_to_char(ch, libc.CString("@WHealing Tank Commands:\r\nhtank [ enter | exit | check ]@n"))
			return TRUE
		}
	} else {
		return FALSE
	}
}
func augmenter(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if libc.StrCmp(libc.CString("augment"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		var (
			strength int = int(ch.Real_abils.Str)
			intel    int = int(ch.Real_abils.Intel)
			wisdom   int = int(ch.Real_abils.Wis)
			speed    int = int(ch.Real_abils.Cha)
			consti   int = int(ch.Real_abils.Con)
			agility  int = int(ch.Real_abils.Dex)
			strcost  int = strength * 1200
			intcost  int = intel * 1200
			concost  int = consti * 1200
			wiscost  int = wisdom * 1200
			agicost  int = agility * 1200
			specost  int = speed * 1200
		)
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("@D                        -----@WBody Augmentations@D-----@n\r\n"))
			send_to_char(ch, libc.CString("@RStrength    @y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), strength, add_commas(int64(strcost)))
			send_to_char(ch, libc.CString("@BIntelligence@y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), intel, add_commas(int64(intcost)))
			send_to_char(ch, libc.CString("@CWisdom      @y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), wisdom, add_commas(int64(wiscost)))
			send_to_char(ch, libc.CString("@GConstitution@y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), consti, add_commas(int64(concost)))
			send_to_char(ch, libc.CString("@mAgility     @y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), agility, add_commas(int64(agicost)))
			send_to_char(ch, libc.CString("@YSpeed       @y: @WCurrently measured at @w%d@W, cost to augment @Y%s@W.@n\r\n"), speed, add_commas(int64(specost)))
			send_to_char(ch, libc.CString("\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("strength"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("str"), &arg[0]) == 0 {
			if strength >= 100 {
				send_to_char(ch, libc.CString("Your strength is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < strcost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Str += 1
				ch.Gold -= strcost
				save_char(ch)
			}
		} else if libc.StrCaseCmp(libc.CString("intelligence"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("int"), &arg[0]) == 0 {
			if intel >= 100 {
				send_to_char(ch, libc.CString("Your intelligence is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < intcost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Intel += 1
				ch.Gold -= intcost
				save_char(ch)
			}
		} else if libc.StrCaseCmp(libc.CString("constitution"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("con"), &arg[0]) == 0 {
			if consti >= 100 {
				send_to_char(ch, libc.CString("Your constitution is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < concost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Con += 1
				ch.Gold -= concost
				save_char(ch)
			}
		} else if libc.StrCaseCmp(libc.CString("speed"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("spe"), &arg[0]) == 0 {
			if speed >= 100 {
				send_to_char(ch, libc.CString("Your speed is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < specost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Cha += 1
				ch.Gold -= specost
				save_char(ch)
			}
		} else if libc.StrCaseCmp(libc.CString("agility"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("agi"), &arg[0]) == 0 {
			if agility >= 100 {
				send_to_char(ch, libc.CString("Your agility is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < agicost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Dex += 1
				ch.Gold -= agicost
				save_char(ch)
			}
		} else if libc.StrCaseCmp(libc.CString("wisdom"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("wis"), &arg[0]) == 0 {
			if wisdom >= 100 {
				send_to_char(ch, libc.CString("Your wisdom how somehow been measured is already as high as it can possibly go.\r\n"))
			} else if ch.Gold < wiscost {
				send_to_char(ch, libc.CString("You can not afford the price!\r\n"))
			} else {
				act(libc.CString("@WThe machine's arm moves out and quickly augments your body with microscopic attachments.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@WThe Augmenter 9001 moves its arm over to @C$n@W and quickly operates on $s body.@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Real_abils.Wis += 1
				ch.Gold -= wiscost
				save_char(ch)
			}
		} else {
			send_to_char(ch, libc.CString("Syntax: augment [str | con | int | wis | agi | spe]\r\n"))
		}
		return TRUE
	} else {
		return FALSE
	}
}
func gravity(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		i     *obj_data
		obj   *obj_data = nil
		arg   [2048]byte
		match int = FALSE
	)
	one_argument(argument, &arg[0])
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
		if GET_OBJ_VNUM(i) == 11 {
			obj = i
		} else {
			continue
		}
	}
	if libc.StrCmp(libc.CString("gravity"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 || libc.StrCmp(libc.CString("generator"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("@WGravity Commands:@n\r\n"))
			send_to_char(ch, libc.CString("@Wgravity [ 0 | N | 10 | 20 | 30 | 40 | 50 | 100 | 200 ]\r\n          [  300 | 400 | 500 | 1,000 | 5,000 | 10,000  ]@n\r\n"))
			return TRUE
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It's broken!\r\n"))
			return TRUE
		}
		if (libc.StrCaseCmp(libc.CString("N"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("n"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("0"), &arg[0]) == 0) && obj.Weight == 0 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("N"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("n"), &arg[0]) == 0 || libc.StrCaseCmp(libc.CString("0"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in normal gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) || ROOM_FLAGGED(ch.In_room, ROOM_GRAVITYX10) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 10
				obj.Weight = 0
			} else {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 0
				obj.Weight = 0
			}
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("10"), &arg[0]) == 0 && obj.Weight == 10 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("10"), &arg[0]) == 0 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && (ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) || ROOM_FLAGGED(ch.In_room, ROOM_GRAVITYX10)) {
			send_to_char(ch, libc.CString("The gravity around you is already at that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("10"), &arg[0]) == 0 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity != 10 && (ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) || ROOM_FLAGGED(ch.In_room, ROOM_GRAVITYX10)) {
			send_to_char(ch, libc.CString("You punch in normal gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 10
			obj.Weight = 0
			match = TRUE
		} else if libc.StrCaseCmp(libc.CString("10"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in ten times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 10
			obj.Weight = 10
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("20"), &arg[0]) == 0 && obj.Weight == 20 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("20"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in twenty times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 20
			obj.Weight = 20
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("30"), &arg[0]) == 0 && obj.Weight == 30 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("30"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in thirty times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 30
			obj.Weight = 30
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("40"), &arg[0]) == 0 && obj.Weight == 40 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("40"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in fourty times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 40
			obj.Weight = 40
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("50"), &arg[0]) == 0 && obj.Weight == 50 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("50"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in fifty times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 50
			obj.Weight = 50
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("100"), &arg[0]) == 0 && obj.Weight == 100 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("100"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in one hundred times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 100
			obj.Weight = 100
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("200"), &arg[0]) == 0 && obj.Weight == 200 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("200"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in two hundred times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 200
			obj.Weight = 200
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("300"), &arg[0]) == 0 && obj.Weight == 300 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("300"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in three hundred times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 300
			obj.Weight = 300
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("400"), &arg[0]) == 0 && obj.Weight == 400 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("400"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in four hundred times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 400
			obj.Weight = 400
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("500"), &arg[0]) == 0 && obj.Weight == 500 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("500"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in five hundred times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 500
			obj.Weight = 500
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("1000"), &arg[0]) == 0 && obj.Weight == 1000 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("1000"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in one thousand times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 1000
			obj.Weight = 1000
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("5000"), &arg[0]) == 0 && obj.Weight == 5000 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("5000"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in five thousand times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 5000
			obj.Weight = 5000
			match = TRUE
		}
		if libc.StrCaseCmp(libc.CString("10000"), &arg[0]) == 0 && obj.Weight == 10000 {
			send_to_char(ch, libc.CString("The gravity generator is already set to that.\r\n"))
			return TRUE
		} else if libc.StrCaseCmp(libc.CString("10000"), &arg[0]) == 0 {
			send_to_char(ch, libc.CString("You punch in ten thousand times gravity on the generator. It hums for a moment\r\nbefore you feel the pressure on your body change.\r\n"))
			act(libc.CString("@W$n@w pushes some buttons on the gravity generator, and you feel a change in pressure on your body.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ROOM_FLAGGED(ch.In_room, ROOM_AURA) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Room_flags[int(ROOM_AURA/32)] &= bitvector_t(int32(^(1 << (int(ROOM_AURA % 32)))))
				send_to_room(ch.In_room, libc.CString("The increased gravity forces the aura to disappear.\r\n"))
			}
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity = 10000
			obj.Weight = 10000
			match = TRUE
		} else if match == FALSE {
			send_to_char(ch, libc.CString("That is not a proper command for this device.\r\n"))
			send_to_char(ch, libc.CString("@WGravity Commands:@n\r\n"))
			send_to_char(ch, libc.CString("@Wgravity [ 0 | N | 10 | 20 | 30 | 40 | 50 | 100 | 200 ]\r\n          [  300 | 400 | 500 | 1,000 | 5,000 | 10,000  ]@n\r\n"))
			return TRUE
		}
		return TRUE
	} else {
		return FALSE
	}
}
func bank(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		amount int
		num    int = 0
		i      *obj_data
		obj    *obj_data = nil
	)
	for i = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; i != nil; i = i.Next_content {
		if GET_OBJ_VNUM(i) == 3034 {
			obj = i
		} else {
			continue
		}
	}
	if libc.StrCmp(libc.CString("balance"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("The ATM is broken!\r\n"))
			return TRUE
		}
		if ch.Bank_gold > 0 {
			send_to_char(ch, libc.CString("Your current balance is %d zenni.\r\n"), ch.Bank_gold)
		} else {
			send_to_char(ch, libc.CString("You currently have no money deposited.\r\n"))
		}
		return TRUE
	} else if libc.StrCmp(libc.CString("wire"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		var (
			arg  [2048]byte
			arg2 [2048]byte
			vict *char_data = nil
		)
		two_arguments(argument, &arg[0], &arg2[0])
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("The ATM is broken!\r\n"))
			return TRUE
		}
		if (func() int {
			amount = libc.Atoi(libc.GoString(&arg[0]))
			return amount
		}()) <= 0 {
			send_to_char(ch, libc.CString("How much do you want to transfer?\r\n"))
			return TRUE
		}
		if ch.Bank_gold < amount+amount/100 {
			send_to_char(ch, libc.CString("You don't have that much zenni in the bank (plus 1%s charge)!\r\n"), "%")
			return TRUE
		}
		if arg2[0] == 0 {
			send_to_char(ch, libc.CString("You want to transfer it to who?!\r\n"))
			return TRUE
		}
		if (func() *char_data {
			vict = get_player_vis(ch, &arg2[0], nil, 1<<1)
			return vict
		}()) == nil {
			var (
				is_file  int = FALSE
				player_i int = 0
				name     [2048]byte
			)
			vict = new(char_data)
			clear_char(vict)
			vict.Player_specials = new(player_special_data)
			stdio.Sprintf(&name[0], "%s", rIntro(ch, &arg2[0]))
			if (func() int {
				player_i = load_char(&name[0], vict)
				return player_i
			}()) > -1 {
				is_file = TRUE
			} else {
				free_char(vict)
				send_to_char(ch, libc.CString("That person doesn't exist.\r\n"))
				return TRUE
			}
			if ch.Desc.User == nil {
				send_to_char(ch, libc.CString("There is an error. Report to Iovan."))
				return TRUE
			}
			if libc.StrCaseCmp(GET_NAME(vict), ch.Desc.Tmp1) == 0 || libc.StrCaseCmp(GET_NAME(vict), ch.Desc.Tmp2) == 0 || libc.StrCaseCmp(GET_NAME(vict), ch.Desc.Tmp3) == 0 || libc.StrCaseCmp(GET_NAME(vict), ch.Desc.Tmp4) == 0 || libc.StrCaseCmp(GET_NAME(vict), ch.Desc.Tmp5) == 0 {
				send_to_char(ch, libc.CString("You can not transfer money to your own offline characters..."))
				if is_file == TRUE {
					free_char(vict)
				}
				return TRUE
			}
			vict.Bank_gold += amount
			ch.Bank_gold -= amount + amount/100
			vict.Pfilepos = player_i
			mudlog(NRM, MAX(ADMLVL_IMPL, int(ch.Player_specials.Invis_level)), TRUE, libc.CString("EXCHANGE: %s gave %s zenni to user %s"), GET_NAME(ch), add_commas(int64(amount)), GET_NAME(vict))
			save_char(vict)
			if is_file == TRUE {
				free_char(vict)
			}
		} else {
			vict.Bank_gold += amount
			ch.Bank_gold -= amount + amount/100
			send_to_char(vict, libc.CString("@WYou have just had @Y%s@W zenni wired into your bank account.@n\r\n"), add_commas(int64(amount)))
		}
		send_to_char(ch, libc.CString("You transfer %s zenni to them.\r\n"), add_commas(int64(amount)))
		act(libc.CString("$n makes a bank transaction."), TRUE, ch, nil, unsafe.Pointer(uintptr(FALSE)), TO_ROOM)
		return TRUE
	} else if libc.StrCmp(libc.CString("deposit"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("The ATM is broken!\r\n"))
			return TRUE
		}
		if (func() int {
			amount = libc.Atoi(libc.GoString(argument))
			return amount
		}()) <= 0 {
			send_to_char(ch, libc.CString("How much do you want to deposit?\r\n"))
			return TRUE
		}
		if ch.Gold < amount {
			send_to_char(ch, libc.CString("You don't have that much zenni!\r\n"))
			return TRUE
		}
		ch.Gold -= amount
		ch.Bank_gold += amount
		send_to_char(ch, libc.CString("You deposit %d zenni.\r\n"), amount)
		act(libc.CString("$n makes a bank transaction."), TRUE, ch, nil, unsafe.Pointer(uintptr(FALSE)), TO_ROOM)
		return TRUE
	} else if libc.StrCmp(libc.CString("withdraw"), (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(cmd)))).Command) == 0 {
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("The ATM is broken!\r\n"))
			return TRUE
		}
		if (func() int {
			amount = libc.Atoi(libc.GoString(argument))
			return amount
		}()) <= 0 {
			send_to_char(ch, libc.CString("How much do you want to withdraw?\r\n"))
			return TRUE
		}
		if ch.Bank_gold < amount {
			send_to_char(ch, libc.CString("You don't have that much zenni!\r\n"))
			return TRUE
		}
		if ch.Bank_gold-(amount+(amount/100+1)) < 0 {
			if amount >= 100 {
				amount = amount + amount/100
			} else if amount < 100 {
				amount = amount + 1
			}
			send_to_char(ch, libc.CString("You need at least %s in the bank with the 1 percent withdraw fee.\r\n"), add_commas(int64(amount)))
			return TRUE
		}
		if ch.Gold+amount > GOLD_CARRY(ch) {
			send_to_char(ch, libc.CString("You can only carry %s zenni, you left the rest.\r\n"), add_commas(int64(GOLD_CARRY(ch))))
			var diff int = (ch.Gold + amount) - GOLD_CARRY(ch)
			ch.Gold = GOLD_CARRY(ch)
			amount -= diff
			if amount >= 100 {
				num = amount / 100
				ch.Bank_gold -= amount + num
			} else if amount < 100 {
				ch.Bank_gold -= amount + 1
			}
			send_to_char(ch, libc.CString("You withdraw %s zenni,  and pay %s in withdraw fees.\r\n.\r\n"), add_commas(int64(amount)), add_commas(int64(num)))
			act(libc.CString("$n makes a bank transaction."), TRUE, ch, nil, unsafe.Pointer(uintptr(FALSE)), TO_ROOM)
			return TRUE
		}
		ch.Gold += amount
		if amount >= 100 {
			num = amount / 100
			ch.Bank_gold -= amount + num
		} else if amount < 100 {
			ch.Bank_gold -= amount + 1
		}
		send_to_char(ch, libc.CString("You withdraw %s zenni, and pay %s in withdraw fees.\r\n"), add_commas(int64(amount)), add_commas(int64(num)))
		act(libc.CString("$n makes a bank transaction."), TRUE, ch, nil, unsafe.Pointer(uintptr(FALSE)), TO_ROOM)
		return TRUE
	} else {
		return FALSE
	}
}
func cleric_marduk(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 12
	tmp = rand_number(1, 10)
	if tmp == 7 || tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 1 && GET_LEVEL(ch) > 13 {
			cast_spell(ch, vict, nil, SPELL_EARTHQUAKE, nil)
			return TRUE
		}
		if tmp == 2 && (GET_LEVEL(ch) > 8 && IS_EVIL(vict)) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_EVIL, nil)
			return TRUE
		}
		if tmp == 3 && GET_LEVEL(ch) > 4 {
			cast_spell(ch, vict, nil, SPELL_BESTOW_CURSE, nil)
			return TRUE
		}
		if tmp == 4 && (GET_LEVEL(ch) > 8 && IS_GOOD(vict)) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_GOOD, nil)
			return TRUE
		}
		if tmp == 5 && (GET_LEVEL(ch) > 4 && affected_by_spell(ch, SPELL_BESTOW_CURSE)) {
			cast_spell(ch, ch, nil, SPELL_REMOVE_CURSE, nil)
			return TRUE
		}
		if tmp == 6 && (GET_LEVEL(ch) > 6 && affected_by_spell(ch, SPELL_POISON)) {
			cast_spell(ch, ch, nil, SPELL_NEUTRALIZE_POISON, nil)
			return TRUE
		}
		if tmp == 7 {
			cast_spell(ch, ch, nil, SPELL_CURE_LIGHT, nil)
			return TRUE
		}
		if tmp == 8 && GET_LEVEL(ch) > 6 && !AFF_FLAGGED(vict, AFF_UNDEAD) {
			cast_spell(ch, vict, nil, SPELL_POISON, nil)
			return TRUE
		}
		if tmp == 9 && GET_LEVEL(ch) > 8 {
			cast_spell(ch, ch, nil, SPELL_CURE_CRITIC, nil)
			return TRUE
		}
		if tmp == 10 && GET_LEVEL(ch) > 10 {
			cast_spell(ch, vict, nil, SPELL_HARM, nil)
			return TRUE
		}
		if tmp == 11 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_LIGHT, nil)
			return TRUE
		}
		if tmp == 12 && GET_LEVEL(ch) > 8 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_CRITIC, nil)
			return TRUE
		}
	}
	return FALSE
}
func cleric_ao(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 8
	tmp = rand_number(1, 10)
	if tmp == 7 || tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 1 && GET_LEVEL(ch) > 13 {
			cast_spell(ch, vict, nil, SPELL_EARTHQUAKE, nil)
			return TRUE
		}
		if tmp == 2 && (GET_LEVEL(ch) > 8 && IS_EVIL(vict)) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_EVIL, nil)
			return TRUE
		}
		if tmp == 3 && (GET_LEVEL(ch) > 8 && IS_GOOD(vict)) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_GOOD, nil)
			return TRUE
		}
		if tmp == 4 && (GET_LEVEL(ch) > 4 && affected_by_spell(ch, SPELL_BESTOW_CURSE)) {
			cast_spell(ch, ch, nil, SPELL_REMOVE_CURSE, nil)
			return TRUE
		}
		if tmp == 5 && (GET_LEVEL(ch) > 6 && affected_by_spell(ch, SPELL_POISON)) {
			cast_spell(ch, ch, nil, SPELL_NEUTRALIZE_POISON, nil)
			return TRUE
		}
		if tmp == 6 {
			cast_spell(ch, ch, nil, SPELL_CURE_LIGHT, nil)
			return TRUE
		}
		if tmp == 7 && GET_LEVEL(ch) > 8 {
			cast_spell(ch, ch, nil, SPELL_CURE_CRITIC, nil)
			return TRUE
		}
		if tmp == 8 && GET_LEVEL(ch) > 10 {
			cast_spell(ch, ch, nil, SPELL_HEAL, nil)
			return TRUE
		}
		if tmp == 9 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_LIGHT, nil)
			return TRUE
		}
		if tmp == 10 && GET_LEVEL(ch) > 8 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_CRITIC, nil)
			return TRUE
		}
	}
	return FALSE
}
func dziak(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 9
	tmp = rand_number(3, 10)
	if tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 2 || tmp == 1 {
			cast_spell(ch, vict, nil, SPELL_SHOCKING_GRASP, nil)
			return TRUE
		}
		if tmp == 3 {
			cast_spell(ch, vict, nil, SPELL_MAGIC_MISSILE, nil)
			return TRUE
		}
		if tmp == 4 {
			cast_spell(ch, vict, nil, SPELL_LIGHTNING_BOLT, nil)
			return TRUE
		}
		if tmp == 5 {
			cast_spell(ch, vict, nil, SPELL_FIREBALL, nil)
			return TRUE
		}
		if tmp == 6 {
			cast_spell(ch, ch, nil, SPELL_CURE_CRITIC, nil)
			return TRUE
		}
		if tmp == 7 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_CRITIC, nil)
			return TRUE
		}
		if tmp == 8 && IS_GOOD(vict) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_GOOD, nil)
			return TRUE
		}
		if tmp == 9 {
			cast_spell(ch, ch, nil, SPELL_HEAL, nil)
			return TRUE
		}
	}
	return FALSE
}
func azimer(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 8
	tmp = rand_number(3, 10)
	if tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 2 || tmp == 1 {
			cast_spell(ch, vict, nil, SPELL_MAGIC_MISSILE, nil)
			return TRUE
		}
		if tmp == 3 {
			cast_spell(ch, vict, nil, SPELL_SHOCKING_GRASP, nil)
			return TRUE
		}
		if tmp == 4 {
			cast_spell(ch, vict, nil, SPELL_LIGHTNING_BOLT, nil)
			return TRUE
		}
		if tmp == 5 {
			cast_spell(ch, vict, nil, SPELL_FIREBALL, nil)
			return TRUE
		}
		if tmp == 6 {
			cast_spell(ch, ch, nil, SPELL_CURE_CRITIC, nil)
			return TRUE
		}
		if tmp == 7 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_CRITIC, nil)
			return TRUE
		}
		if tmp == 8 && IS_GOOD(vict) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_GOOD, nil)
			return TRUE
		}
	}
	return FALSE
}
func lyrzaxyn(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 8
	tmp = rand_number(3, 10)
	if tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 2 || tmp == 1 {
			cast_spell(ch, vict, nil, SPELL_MAGIC_MISSILE, nil)
			return TRUE
		}
		if tmp == 3 {
			cast_spell(ch, vict, nil, SPELL_SHOCKING_GRASP, nil)
			return TRUE
		}
		if tmp == 4 {
			cast_spell(ch, vict, nil, SPELL_LIGHTNING_BOLT, nil)
			return TRUE
		}
		if tmp == 5 {
			cast_spell(ch, vict, nil, SPELL_FIREBALL, nil)
			return TRUE
		}
		if tmp == 6 {
			cast_spell(ch, ch, nil, SPELL_CURE_CRITIC, nil)
			return TRUE
		}
		if tmp == 7 {
			cast_spell(ch, vict, nil, SPELL_INFLICT_CRITIC, nil)
			return TRUE
		}
		if tmp == 8 && IS_GOOD(vict) {
			cast_spell(ch, vict, nil, SPELL_DISPEL_GOOD, nil)
			return TRUE
		}
	}
	return FALSE
}
func magic_user(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
	var (
		tmp      int
		num_used int = 0
		vict     *char_data
	)
	if IS_NPC(ch) && int(ch.Position) > POS_SITTING && int(ch.Chclass) == CLASS_KABITO {
		if !affected_by_spell(ch, SPELL_MAGE_ARMOR) {
			cast_spell(ch, ch, nil, SPELL_MAGE_ARMOR, nil)
			return TRUE
		}
	}
	if cmd != 0 || int(ch.Position) != POS_FIGHTING {
		return FALSE
	}
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = vict.Next_in_room {
		if vict.Fighting == ch && rand_number(0, 4) == 0 {
			break
		}
	}
	if vict == nil {
		vict = ch.Fighting
	}
	num_used = 6
	tmp = rand_number(2, 10)
	if tmp == 8 || tmp == 9 || tmp == 10 {
		tmp = rand_number(1, num_used)
		if tmp == 1 && GET_LEVEL(ch) > 1 {
			cast_spell(ch, vict, nil, SPELL_CHILL_TOUCH, nil)
			return TRUE
		}
		if tmp == 2 && !affected_by_spell(ch, SPELL_MAGE_ARMOR) {
			cast_spell(ch, ch, nil, SPELL_MAGE_ARMOR, nil)
			return TRUE
		}
		if tmp == 3 && GET_LEVEL(ch) > 1 {
			cast_spell(ch, vict, nil, SPELL_BURNING_HANDS, nil)
			return TRUE
		}
		if tmp == 4 && GET_LEVEL(ch) > 1 {
			cast_spell(ch, vict, nil, SPELL_MAGIC_MISSILE, nil)
			return TRUE
		}
		if tmp == 5 && GET_LEVEL(ch) > 5 {
			cast_spell(ch, vict, nil, SPELL_SHOCKING_GRASP, nil)
			return TRUE
		}
		if tmp == 6 && GET_LEVEL(ch) > 9 {
			cast_spell(ch, vict, nil, SPELL_LIGHTNING_BOLT, nil)
			return TRUE
		}
	}
	return FALSE
}
