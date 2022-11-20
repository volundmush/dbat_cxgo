package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

type bfs_queue_struct struct {
	Room int
	Dir  int8
	Next *bfs_queue_struct
}

var bfs_queue_head *bfs_queue_struct = nil
var bfs_queue_tail *bfs_queue_struct = nil

func VALID_EDGE(x int, y int) bool {
	if world[x].Dir_option[y] == nil || world[x].Dir_option[y].To_room == int(-1) {
		return false
	}
	if config_info.Play.Track_through_doors == 0 && EXIT_FLAGGED(world[x].Dir_option[y], 1<<1) {
		return false
	}
	if ROOM_FLAGGED(world[x].Dir_option[y].To_room, ROOM_NOTRACK) || ROOM_FLAGGED(world[x].Dir_option[y].To_room, ROOM_BFS_MARK) {
		return false
	}
	return true
}
func bfs_enqueue(room int, dir int) {
	var curr *bfs_queue_struct
	curr = new(bfs_queue_struct)
	curr.Room = room
	curr.Dir = int8(dir)
	curr.Next = nil
	if bfs_queue_tail != nil {
		bfs_queue_tail.Next = curr
		bfs_queue_tail = curr
	} else {
		bfs_queue_head = func() *bfs_queue_struct {
			bfs_queue_tail = curr
			return bfs_queue_tail
		}()
	}
}
func bfs_dequeue() {
	var curr *bfs_queue_struct
	curr = bfs_queue_head
	if (func() *bfs_queue_struct {
		bfs_queue_head = bfs_queue_head.Next
		return bfs_queue_head
	}()) == nil {
		bfs_queue_tail = nil
	}
	libc.Free(unsafe.Pointer(curr))
}
func bfs_clear_queue() {
	for bfs_queue_head != nil {
		bfs_dequeue()
	}
}
func find_first_step(src int, target int) int {
	var (
		curr_dir  int
		curr_room int
	)
	if src == int(-1) || target == int(-1) || src > top_of_world || target > top_of_world {
		basic_mud_log(libc.CString("SYSERR: Illegal value %d or %d passed to find_first_step. (%s)"), src, target, "__THISFILE__")
		return -1
	}
	if src == target {
		return -2
	}
	for curr_room = 0; curr_room <= top_of_world; curr_room++ {
		REMOVE_BIT_AR(world[curr_room].Room_flags[:], ROOM_BFS_MARK)
	}
	SET_BIT_AR(world[src].Room_flags[:], ROOM_BFS_MARK)
	for curr_dir = 0; curr_dir < NUM_OF_DIRS; curr_dir++ {
		if VALID_EDGE(src, curr_dir) {
			SET_BIT_AR(world[world[src].Dir_option[curr_dir].To_room].Room_flags[:], ROOM_BFS_MARK)
			bfs_enqueue(world[src].Dir_option[curr_dir].To_room, curr_dir)
		}
	}
	for bfs_queue_head != nil {
		if bfs_queue_head.Room == target {
			curr_dir = int(bfs_queue_head.Dir)
			bfs_clear_queue()
			return curr_dir
		} else {
			for curr_dir = 0; curr_dir < NUM_OF_DIRS; curr_dir++ {
				if VALID_EDGE(bfs_queue_head.Room, curr_dir) {
					SET_BIT_AR(world[world[bfs_queue_head.Room].Dir_option[curr_dir].To_room].Room_flags[:], ROOM_BFS_MARK)
					bfs_enqueue(world[bfs_queue_head.Room].Dir_option[curr_dir].To_room, int(bfs_queue_head.Dir))
				}
			}
			bfs_dequeue()
		}
	}
	return -4
}
func do_sradar(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vehicle  *obj_data = nil
		controls *obj_data = nil
		dir      int       = 0
		noship   int       = 0
		arg      [2048]byte
		planet   [20]byte
	)
	one_argument(argument, &arg[0])
	if !PLR_FLAGGED(ch, PLR_PILOTING) && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("You are not flying a ship, maybe you want detect?\r\n"))
		return
	}
	if (func() *obj_data {
		controls = find_control(ch)
		return controls
	}()) == nil && ch.Admlevel < 1 {
		send_to_char(ch, libc.CString("@wYou have nothing to control here!\r\n"))
		return
	}
	if !PLR_FLAGGED(ch, PLR_PILOTING) && ch.Admlevel >= 1 {
		noship = 1
	} else if (func() *obj_data {
		vehicle = find_vehicle_by_vnum(controls.Value[0])
		return vehicle
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou can't find anything to pilot.\r\n"))
		return
	}
	if noship == 0 && SECT(vehicle.In_room) != SECT_SPACE {
		send_to_char(ch, libc.CString("@wYour ship is not in space!\r\n"))
		return
	}
	if noship == 1 && SECT(ch.In_room) != SECT_SPACE {
		send_to_char(ch, libc.CString("@wYou are not even in space!\r\n"))
		return
	}
	if arg[0] == 0 {
		if ch.Admlevel >= 1 && noship == 1 {
			printmap(ch.In_room, ch, 0, -1)
		} else {
			printmap(vehicle.In_room, ch, 0, GET_OBJ_VNUM(vehicle))
		}
		return
	}
	if ch.Ping > 0 {
		send_to_char(ch, libc.CString("@wYou need to wait a few more seconds before pinging a destination again.\r\n"))
		return
	}
	if noship == 0 {
		if libc.StrCaseCmp(&arg[0], libc.CString("earth")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Earth")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0xA013))
			stdio.Sprintf(&planet[0], "Earth")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("frigid")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Frigid")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0x78A9))
			stdio.Sprintf(&planet[0], "Frigid")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("konack")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Konack")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0x69B9))
			stdio.Sprintf(&planet[0], "Konack")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("vegeta")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Vegeta")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0x7E6D))
			stdio.Sprintf(&planet[0], "Vegeta")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("aether")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Aether")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0xA3E7))
			stdio.Sprintf(&planet[0], "Aether")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("namek")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Namek")) == 0 {
			dir = find_first_step(vehicle.In_room, real_room(0xA780))
			stdio.Sprintf(&planet[0], "Namek")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy1")) == 0 && ch.Radar1 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy2")) == 0 && ch.Radar2 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy3")) == 0 && ch.Radar3 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy1")) == 0 && ch.Radar1 > 0 {
			var rad int = ch.Radar1
			dir = find_first_step(vehicle.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy One")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy2")) == 0 && ch.Radar2 > 0 {
			var rad int = ch.Radar2
			dir = find_first_step(vehicle.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy Two")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy3")) == 0 && ch.Radar3 > 0 {
			var rad int = ch.Radar3
			dir = find_first_step(vehicle.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy Three")
		} else {
			send_to_char(ch, libc.CString("@wThat is not an existing planet.@n\r\n"))
			return
		}
	}
	if noship == 1 {
		if libc.StrCaseCmp(&arg[0], libc.CString("earth")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Earth")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0xA013))
			stdio.Sprintf(&planet[0], "Earth")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("frigid")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Frigid")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0x78A9))
			stdio.Sprintf(&planet[0], "Frigid")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("konack")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Konack")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0x69B9))
			stdio.Sprintf(&planet[0], "Konack")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("vegeta")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Vegeta")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0x7E6D))
			stdio.Sprintf(&planet[0], "Vegeta")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("aether")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Aether")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0xA3E7))
			stdio.Sprintf(&planet[0], "Aether")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("namek")) == 0 || libc.StrCaseCmp(&arg[0], libc.CString("Namek")) == 0 {
			dir = find_first_step(ch.In_room, real_room(0xA780))
			stdio.Sprintf(&planet[0], "Namek")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy1")) == 0 && ch.Radar1 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy2")) == 0 && ch.Radar2 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy3")) == 0 && ch.Radar3 <= 0 {
			send_to_char(ch, libc.CString("@wYou haven't launched that buoy.\r\n"))
			return
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy1")) == 0 && ch.Radar1 > 0 {
			var rad int = ch.Radar1
			dir = find_first_step(ch.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy One")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy2")) == 0 && ch.Radar2 > 0 {
			var rad int = ch.Radar2
			dir = find_first_step(ch.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy Two")
		} else if libc.StrCaseCmp(&arg[0], libc.CString("buoy3")) == 0 && ch.Radar3 > 0 {
			var rad int = ch.Radar3
			dir = find_first_step(ch.In_room, real_room(rad))
			stdio.Sprintf(&planet[0], "Buoy Three")
		} else {
			send_to_char(ch, libc.CString("@wThat is not an existing planet.@n\r\n"))
			return
		}
	}
	switch dir {
	case (-1):
		send_to_char(ch, libc.CString("Hmm.. something seems to be wrong.\r\n"))
	case (-2):
		send_to_char(ch, libc.CString("@wThe radar shows that your are already there.@n\r\n"))
	case (-4):
		send_to_char(ch, libc.CString("@wYou should be in space to use the radar.@n\r\n"))
	default:
		send_to_char(ch, libc.CString("@wYour radar detects @C%s@w to the @G%s@n\r\n"), &planet[0], dirs[dir])
	}
	ch.Ping = 5
}
func do_radar(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		room   int = 0
		dir    int
		num    int = 0
		found  int = 0
		found2 int = 0
	)
	_ = found2
	var fcount int = 0
	var tch *char_data
	var obj *obj_data
	var radar *obj_data = find_obj_in_list_vnum_good(ch.Carrying, 12)
	if radar == nil {
		send_to_char(ch, libc.CString("You do not even have a dragon radar!\r\n"))
		return
	}
	if IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are a freaking mob!\r\n"))
		return
	} else {
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		act(libc.CString("$n holds up a dragon radar and pushes its button."), 0, ch, nil, nil, TO_ROOM)
		for num < 20000 {
			if real_room(room) != int(-1) {
				for obj = world[real_room(room)].Contents; obj != nil; obj = obj.Next_content {
					if OBJ_FLAGGED(obj, ITEM_FORGED) {
						continue
					} else if GET_OBJ_VNUM(obj) == 20 || GET_OBJ_VNUM(obj) == 21 || GET_OBJ_VNUM(obj) == 22 || GET_OBJ_VNUM(obj) == 23 || GET_OBJ_VNUM(obj) == 24 || GET_OBJ_VNUM(obj) == 25 || GET_OBJ_VNUM(obj) == 26 {
						dir = find_first_step(ch.In_room, obj.In_room)
						fcount += 1
						switch dir {
						case (-1):
							send_to_char(ch, libc.CString("Hmm.. something seems to be wrong.\r\n"))
						case (-2):
							send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a dragonball right here!\r\n"), fcount)
						case (-4):
							send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a faint dragonball signal, but can not direct you further.\r\n"), fcount)
						default:
							send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a dragonball %s of here.\r\n"), fcount, dirs[dir])
						}
						found = 1
					}
				}
				for tch = world[real_room(room)].People; tch != nil; tch = tch.Next_in_room {
					if tch == ch {
						continue
					}
					for obj = tch.Carrying; obj != nil; obj = obj.Next_content {
						if OBJ_FLAGGED(obj, ITEM_FORGED) {
							continue
						} else if GET_OBJ_VNUM(obj) == 20 || GET_OBJ_VNUM(obj) == 21 || GET_OBJ_VNUM(obj) == 22 || GET_OBJ_VNUM(obj) == 23 || GET_OBJ_VNUM(obj) == 24 || GET_OBJ_VNUM(obj) == 25 || GET_OBJ_VNUM(obj) == 26 {
							dir = find_first_step(ch.In_room, tch.In_room)
							fcount += 1
							switch dir {
							case (-1):
								send_to_char(ch, libc.CString("Hmm.. something seems to be wrong.\r\n"))
							case (-2):
								send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a dragonball right here!\r\n"), fcount)
							case (-4):
								send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a faint dragonball signal, but can not direct you further.\r\n"), fcount)
							default:
								send_to_char(ch, libc.CString("@D<@G%d@D>@w The radar detects a dragonball %s of here.\r\n"), fcount, dirs[dir])
							}
							found = 1
						}
					}
				}
			}
			num += 1
			room += 1
		}
		if found == 0 {
			send_to_char(ch, libc.CString("The radar didn't detect any dragonballs on the planet.\r\n"))
			return
		}
	}
}
func do_track(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		vict  *char_data
		i     *descriptor_data
		count int = 0
		dir   int
	)
	if IS_NPC(ch) || GET_SKILL(ch, SKILL_SENSE) == 0 {
		send_to_char(ch, libc.CString("You have no idea how.\r\n"))
		return
	}
	if ch.Suppression <= 20 && ch.Suppression > 0 {
		send_to_char(ch, libc.CString("You are concentrating too hard on suppressing your powerlevel at this level of suppression.\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Whom are you trying to sense?\r\n"))
		return
	} else if arg[0] == 0 && ch.Fighting != nil {
		vict = ch.Fighting
		send_to_char(ch, libc.CString("You focus on the one your are fighting.\r\n"))
		if AFF_FLAGGED(vict, AFF_NOTRACK) || int(vict.Race) == RACE_ANDROID {
			send_to_char(ch, libc.CString("You can't sense them.\r\n"))
			return
		}
		if read_sense_memory(ch, vict) == 0 {
			send_to_char(ch, libc.CString("You will remember their ki signal from now on.\r\n"))
			sense_memory_write(ch, vict)
		}
		act(libc.CString("You look at $N@n intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("$n looks at you intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("$n looks at $N@n intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		if int(vict.Race) != RACE_ANDROID {
			if vict.Alignment > 50 && vict.Alignment < 200 {
				send_to_char(ch, libc.CString("You sense slightly pure and good ki from them.\r\n"))
			} else if vict.Alignment > 200 && vict.Alignment < 500 {
				send_to_char(ch, libc.CString("You sense a pure and good ki from them.\r\n"))
			} else if vict.Alignment >= 500 {
				send_to_char(ch, libc.CString("You sense an extremely pure and good ki from them.\r\n"))
			} else if vict.Alignment < -50 && vict.Alignment > -200 {
				send_to_char(ch, libc.CString("You sense slightly sour and evil ki from them.\r\n"))
			} else if vict.Alignment < -200 && vict.Alignment > -500 {
				send_to_char(ch, libc.CString("You sense a sour and evil ki from them.\r\n"))
			} else if vict.Alignment <= -500 {
				send_to_char(ch, libc.CString("You sense an extremely evil ki from them.\r\n"))
			} else if vict.Alignment > -50 && vict.Alignment < 50 {
				send_to_char(ch, libc.CString("You sense slightly mild indefinable ki from them.\r\n"))
			}
			if vict.Hit > ch.Hit*50 {
				send_to_char(ch, libc.CString("Their power is so huge it boggles your mind and crushes your spirit to fight!\n"))
			} else if vict.Hit > ch.Hit*25 {
				send_to_char(ch, libc.CString("Their power is so much larger than you that you would die like an insect.\n"))
			} else if vict.Hit > ch.Hit*10 {
				send_to_char(ch, libc.CString("Their power is many times larger than your own.\n"))
			} else if vict.Hit > ch.Hit*5 {
				send_to_char(ch, libc.CString("Their power is a great deal larger than your own.\n"))
			} else if vict.Hit > ch.Hit*2 {
				send_to_char(ch, libc.CString("Their power is more than twice as large as your own.\n"))
			} else if vict.Hit > ch.Hit {
				send_to_char(ch, libc.CString("Their power is about twice as large as your own.\n"))
			} else if vict.Hit == ch.Hit {
				send_to_char(ch, libc.CString("Their power is exactly as strong as you.\n"))
			} else if float64(vict.Hit) >= float64(ch.Hit)*0.75 {
				send_to_char(ch, libc.CString("Their power is about a quarter of your own or larger.\n"))
			} else if float64(vict.Hit) >= float64(ch.Hit)*0.5 {
				send_to_char(ch, libc.CString("Their power is about half of your own or larger.\n"))
			} else if float64(vict.Hit) >= float64(ch.Hit)*0.25 {
				send_to_char(ch, libc.CString("Their power is about a quarter of your own or larger.\n"))
			} else if float64(vict.Hit) >= float64(ch.Hit)*0.1 {
				send_to_char(ch, libc.CString("Their power is about a tenth of your own or larger.\n"))
			} else if float64(vict.Hit) >= float64(ch.Hit)*0.01 {
				send_to_char(ch, libc.CString("Their power is less than a tenth of your own.\n"))
			} else if float64(vict.Hit) < float64(ch.Hit)*0.01 {
				send_to_char(ch, libc.CString("Their power is less than 1 percent of your own. What a weakling...\n"))
			}
		} else {
			send_to_char(ch, libc.CString("You can't sense their powerlevel as they are a machine.\r\n"))
		}
		return
	}
	if libc.StrCaseCmp(&arg[0], libc.CString("scan")) == 0 {
		for i = descriptor_list; i != nil; i = i.Next {
			if i.Connected != CON_PLAYING {
				continue
			} else if int(i.Character.Race) == RACE_ANDROID {
				continue
			} else if i.Character == ch {
				continue
			} else if float64(i.Character.Hit) < (float64(ch.Hit)*0.001)+1 {
				continue
			} else if planet_check(ch, i.Character) != 0 {
				if readIntro(ch, i.Character) == 1 {
					send_to_char(ch, libc.CString("@D[@Y%d@D] @CYou sense @c%s@C with "), count+1, get_i_name(ch, i.Character))
				} else {
					send_to_char(ch, libc.CString("@D[@Y%d@D] @CYou sense "), count+1)
				}
				if i.Character.Hit > ch.Hit*50 {
					send_to_char(ch, libc.CString("a power so huge it boggles your mind and crushes your spirit to fight!\n"))
				} else if i.Character.Hit > ch.Hit*25 {
					send_to_char(ch, libc.CString("a power so much larger than you that you would die like an insect.\n"))
				} else if i.Character.Hit > ch.Hit*10 {
					send_to_char(ch, libc.CString("a power that is many times larger than your own.\n"))
				} else if i.Character.Hit > ch.Hit*5 {
					send_to_char(ch, libc.CString("a power that is a great deal larger than your own.\n"))
				} else if i.Character.Hit > ch.Hit*2 {
					send_to_char(ch, libc.CString("a power that is more than twice as large as your own.\n"))
				} else if i.Character.Hit > ch.Hit {
					send_to_char(ch, libc.CString("a power that is about twice as large as your own.\n"))
				} else if i.Character.Hit == ch.Hit {
					send_to_char(ch, libc.CString("a power that is exactly as strong as you.\n"))
				} else if float64(i.Character.Hit) >= float64(ch.Hit)*0.75 {
					send_to_char(ch, libc.CString("a power that is about a quarter of your own or larger.\n"))
				} else if float64(i.Character.Hit) >= float64(ch.Hit)*0.5 {
					send_to_char(ch, libc.CString("a power that is about half of your own or larger.\n"))
				} else if float64(i.Character.Hit) >= float64(ch.Hit)*0.25 {
					send_to_char(ch, libc.CString("a power that is about a quarter of your own or larger.\n"))
				} else if float64(i.Character.Hit) >= float64(ch.Hit)*0.1 {
					send_to_char(ch, libc.CString("a power that is about a tenth of your own or larger.\n"))
				} else if float64(i.Character.Hit) >= float64(ch.Hit)*0.01 {
					send_to_char(ch, libc.CString("a power that is less than a tenth of your own.\n"))
				} else if float64(i.Character.Hit) < float64(ch.Hit)*0.01 {
					send_to_char(ch, libc.CString("a power that is less than 1 percent of your own. What a weakling...\n"))
				}
				if i.Character.Alignment >= 500 {
					send_to_char(ch, libc.CString("@wYou sense an extremely pure and good ki from them.@n\n"))
				} else if i.Character.Alignment > 200 {
					send_to_char(ch, libc.CString("@wYou sense a pure and good ki from them.@n\n"))
				} else if i.Character.Alignment > 50 {
					send_to_char(ch, libc.CString("@wYou sense slightly pure and good ki from them.@n\n"))
				} else if i.Character.Alignment > -50 {
					send_to_char(ch, libc.CString("@wYou sense a slightly mild indefinable ki from them.@n\n"))
				} else if i.Character.Alignment > -200 {
					send_to_char(ch, libc.CString("@wYou sense a slightly sour and evil ki from them.@n\n"))
				} else if i.Character.Alignment > -500 {
					send_to_char(ch, libc.CString("@wYou sense a sour and evil ki from them.@n\n"))
				} else if i.Character.Alignment <= -500 {
					send_to_char(ch, libc.CString("@wYou sense an extremely evil ki from them.@n\n"))
				}
				var blah *byte = sense_location(i.Character)
				send_to_char(ch, libc.CString("@wLastly you sense that they are at... @C%s@n\n"), blah)
				count++
				libc.Free(unsafe.Pointer(blah))
			}
		}
		if count == 0 {
			send_to_char(ch, libc.CString("You sense that there is no one important around.@n\n"))
		}
		return
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<1)
		return vict
	}()) == nil {
		send_to_char(ch, libc.CString("No one is around by that name.\r\n"))
		return
	}
	if AFF_FLAGGED(vict, AFF_NOTRACK) || int(vict.Race) == RACE_ANDROID {
		send_to_char(ch, libc.CString("You can't sense them.\r\n"))
		return
	}
	if float64(vict.Hit) < (float64(ch.Hit)*0.001)+1 {
		if ch.In_room == vict.In_room {
			if read_sense_memory(ch, vict) == 0 {
				send_to_char(ch, libc.CString("Their powerlevel is too weak for you to sense properly, but you will recognise their ki signal from now on.\r\n"))
				sense_memory_write(ch, vict)
			} else {
				send_to_char(ch, libc.CString("Their powerlevel is too weak for you to sense properly.\r\n"))
			}
		} else {
			send_to_char(ch, libc.CString("Their powerlevel is too weak for you to sense properly.\r\n"))
		}
		return
	}
	if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_YARDRAT) && ROOM_FLAGGED(vict.In_room, ROOM_YARDRAT)) {
		send_to_char(ch, libc.CString("@WSense@D: @YYardrat@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_EARTH) && ROOM_FLAGGED(vict.In_room, ROOM_EARTH)) {
		send_to_char(ch, libc.CString("@WSense@D: @GEarth@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) && ROOM_FLAGGED(vict.In_room, ROOM_VEGETA)) {
		send_to_char(ch, libc.CString("@WSense@D: @YVegeta@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) && ROOM_FLAGGED(vict.In_room, ROOM_NAMEK)) {
		send_to_char(ch, libc.CString("@WSense@D: @gNamek@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) && ROOM_FLAGGED(vict.In_room, ROOM_FRIGID)) {
		send_to_char(ch, libc.CString("@WSense@D: @CFrigid@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_AETHER) && ROOM_FLAGGED(vict.In_room, ROOM_AETHER)) {
		send_to_char(ch, libc.CString("@WSense@D: @mAetherh@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_KONACK) && ROOM_FLAGGED(vict.In_room, ROOM_KONACK)) {
		send_to_char(ch, libc.CString("@WSense@D: @MKonack@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) && ROOM_FLAGGED(vict.In_room, ROOM_KANASSA)) {
		send_to_char(ch, libc.CString("@WSense@D: @cKanassa@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!ROOM_FLAGGED(ch.In_room, ROOM_ARLIA) && ROOM_FLAGGED(vict.In_room, ROOM_ARLIA)) {
		send_to_char(ch, libc.CString("@WSense@D: @yArlia@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else if int(ch.Skills[SKILL_SENSE]) == 100 && (!PLANET_ZENITH(ch.In_room) && PLANET_ZENITH(vict.In_room)) {
		send_to_char(ch, libc.CString("@WSense@D: @CZenith@n\r\n"))
		if func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<1)
			return vict
		}() != nil {
			var blah *byte = sense_location(vict)
			send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
			libc.Free(unsafe.Pointer(blah))
		}
	} else {
		if GET_SKILL(ch, SKILL_SENSE) < rand_number(1, 101) {
			var tries int = 10
			for {
				dir = rand_number(0, int(NUM_OF_DIRS-1))
				if CAN_GO(ch, dir) || func() int {
					p := &tries
					*p--
					return *p
				}() == 0 {
					break
				}
			}
			send_to_char(ch, libc.CString("You sense them %s faintly from here, but are unsure....\r\n"), dirs[dir])
			improve_skill(ch, SKILL_SENSE, 1)
			improve_skill(ch, SKILL_SENSE, 1)
			improve_skill(ch, SKILL_SENSE, 1)
			return
		}
		dir = find_first_step(ch.In_room, vict.In_room)
		switch dir {
		case (-1):
			send_to_char(ch, libc.CString("Hmm.. something seems to be wrong.\r\n"))
		case (-2):
			act(libc.CString("You look at $N@n intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("$n looks at you intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("$n looks at $N intently for a moment."), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if int(vict.Race) != RACE_ANDROID {
				if vict.Alignment > 50 && vict.Alignment < 200 {
					send_to_char(ch, libc.CString("You sense slightly pure and good ki from them.\r\n"))
				} else if vict.Alignment > 200 && vict.Alignment < 500 {
					send_to_char(ch, libc.CString("You sense a pure and good ki from them.\r\n"))
				} else if vict.Alignment >= 500 {
					send_to_char(ch, libc.CString("You sense an extremely pure and good ki from them.\r\n"))
				} else if vict.Alignment < -50 && vict.Alignment > -200 {
					send_to_char(ch, libc.CString("You sense slightly sour and evil ki from them.\r\n"))
				} else if vict.Alignment < -200 && vict.Alignment > -500 {
					send_to_char(ch, libc.CString("You sense a sour and evil ki from them.\r\n"))
				} else if vict.Alignment <= -500 {
					send_to_char(ch, libc.CString("You sense an extremely evil ki from them.\r\n"))
				} else if vict.Alignment > -50 && vict.Alignment < 50 {
					send_to_char(ch, libc.CString("You sense slightly mild indefinable ki from them.\r\n"))
				}
			}
			if int(vict.Race) != RACE_ANDROID {
				if vict.Hit > ch.Hit*50 {
					send_to_char(ch, libc.CString("Their power is so huge it boggles your mind and crushes your spirit to fight!\n"))
				} else if vict.Hit > ch.Hit*25 {
					send_to_char(ch, libc.CString("Their power is so much larger than you that you would die like an insect.\n"))
				} else if vict.Hit > ch.Hit*10 {
					send_to_char(ch, libc.CString("Their power is many times larger than your own.\n"))
				} else if vict.Hit > ch.Hit*5 {
					send_to_char(ch, libc.CString("Their power is a great deal larger than your own.\n"))
				} else if vict.Hit > ch.Hit*2 {
					send_to_char(ch, libc.CString("Their power is more than twice as large as your own.\n"))
				} else if vict.Hit > ch.Hit {
					send_to_char(ch, libc.CString("Their power is about twice as large as your own.\n"))
				} else if vict.Hit == ch.Hit {
					send_to_char(ch, libc.CString("Their power is exactly as strong as you.\n"))
				} else if float64(vict.Hit) >= float64(ch.Hit)*0.75 {
					send_to_char(ch, libc.CString("Their power is about a quarter of your own or larger.\n"))
				} else if float64(vict.Hit) >= float64(ch.Hit)*0.5 {
					send_to_char(ch, libc.CString("Their power is about half of your own or larger.\n"))
				} else if float64(vict.Hit) >= float64(ch.Hit)*0.25 {
					send_to_char(ch, libc.CString("Their power is about a quarter of your own or larger.\n"))
				} else if float64(vict.Hit) >= float64(ch.Hit)*0.1 {
					send_to_char(ch, libc.CString("Their power is about a tenth of your own or larger.\n"))
				} else if float64(vict.Hit) >= float64(ch.Hit)*0.01 {
					send_to_char(ch, libc.CString("Their power is less than a tenth of your own.\n"))
				} else if float64(vict.Hit) < float64(ch.Hit)*0.01 {
					send_to_char(ch, libc.CString("Their power is less than 1 percent of your own. What a weakling...\n"))
				}
				if read_sense_memory(ch, vict) == 0 {
					send_to_char(ch, libc.CString("You will remember their ki signal from now on.\r\n"))
					sense_memory_write(ch, vict)
				}
			} else {
				send_to_char(ch, libc.CString("You can't sense their powerlevel as they are a machine.\r\n"))
			}
		case (-3):
			send_to_char(ch, libc.CString("You are too far to sense %s accurately from here.\r\n"), HMHR(vict))
		case (-4):
			send_to_char(ch, libc.CString("You can't sense %s from here.\r\n"), HMHR(vict))
		default:
			if int(ch.Skills[SKILL_SENSE]) >= 75 {
				var blah *byte = sense_location(vict)
				send_to_char(ch, libc.CString("You sense them %s from here!\r\n"), dirs[dir])
				send_to_char(ch, libc.CString("@WSense@D: @Y%s@n\r\n"), blah)
				libc.Free(unsafe.Pointer(blah))
			} else {
				send_to_char(ch, libc.CString("You sense them %s from here!\r\n"), dirs[dir])
				break
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			improve_skill(ch, SKILL_SENSE, 1)
			improve_skill(ch, SKILL_SENSE, 1)
			improve_skill(ch, SKILL_SENSE, 1)
		}
	}
}
