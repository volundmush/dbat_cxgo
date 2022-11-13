package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const NEED_OPEN = 1
const NEED_CLOSED = 2
const NEED_UNLOCKED = 4
const NEED_LOCKED = 8

func handle_teleport(ch *char_data, tar *char_data, location int) {
	var success int = FALSE
	if location != 0 {
		char_from_room(ch)
		char_to_room(ch, real_room(room_vnum(location)))
		success = TRUE
	} else if tar != nil {
		char_from_room(ch)
		char_to_room(ch, tar.In_room)
		success = TRUE
	}
	if success == TRUE {
		act(libc.CString("@w$n@w appears in an instant out of nowhere!@n"), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Drag != nil && !IS_NPC(ch.Drag) {
			char_from_room(ch.Drag)
			char_to_room(ch.Drag, ch.In_room)
			act(libc.CString("@w$n@w appears in an instant out of nowhere being dragged by $N!@n"), TRUE, ch.Drag, nil, unsafe.Pointer(ch), TO_NOTVICT)
		}
		if ch.Grappling != nil && !IS_NPC(ch.Grappling) {
			char_from_room(ch.Grappling)
			char_to_room(ch.Grappling, ch.In_room)
			act(libc.CString("@w$n@w appears in an instant out of nowhere being grappled by $N!@n"), TRUE, ch.Grappling, nil, unsafe.Pointer(ch), TO_NOTVICT)
		}
		if ch.Player_specials.Carrying != nil {
			char_from_room(ch.Player_specials.Carrying)
			char_to_room(ch.Player_specials.Carrying, ch.In_room)
			act(libc.CString("@w$n@w appears in an instant out of nowhere being carried by $N!@n"), TRUE, ch.Player_specials.Carrying, nil, unsafe.Pointer(ch), TO_NOTVICT)
		}
		if ch.Grappled != nil && !IS_NPC(ch.Grappled) {
			char_from_room(ch.Grappled)
			char_to_room(ch.Grappled, ch.In_room)
			act(libc.CString("@w$n@w appears in an instant out of nowhere being grappled by $N!@n"), TRUE, ch.Grappled, nil, unsafe.Pointer(ch), TO_NOTVICT)
		}
		if ch.Drag != nil && IS_NPC(ch.Drag) {
			act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
			act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
			ch.Drag.Dragged = nil
			ch.Drag = nil
		}
		if ch.Grappling != nil && IS_NPC(ch.Grappling) {
			ch.Grappling.Grap = -1
			ch.Grappling.Grappled = nil
			ch.Grappling = nil
			ch.Grap = -1
		}
		if ch.Grappled != nil && IS_NPC(ch.Grappled) {
			ch.Grappled.Grap = -1
			ch.Grappled.Grappling = nil
			ch.Grappled = nil
			ch.Grap = -1
		}
	} else {
		basic_mud_log(libc.CString("ERROR: handle_teleport called without a destination."))
		return
	}
}
func do_carry(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	var vict *char_data = nil
	var arg [2048]byte
	if ch.Drag != nil {
		send_to_char(ch, libc.CString("You are busy dragging someone at the moment.\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	if ch.Player_specials.Carrying != nil {
		if ch.Alignment > 50 {
			carry_drop(ch, 0)
		} else {
			carry_drop(ch, 1)
		}
		return
	} else {
		one_argument(argument, &arg[0])
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("You want to carry who?\r\n"))
			return
		}
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("That person isn't here.\r\n"))
			return
		}
		if IS_NPC(vict) {
			send_to_char(ch, libc.CString("There's no point in carrying them.\r\n"))
			return
		}
		if vict.Player_specials.Carried_by != nil {
			send_to_char(ch, libc.CString("Someone is already carrying them!\r\n"))
			return
		}
		if int(vict.Position) > POS_SLEEPING {
			send_to_char(ch, libc.CString("They are not unconcious.\r\n"))
			return
		}
		if GET_PC_WEIGHT(vict)+vict.Carry_weight > int(max_carry_weight(ch)) {
			act(libc.CString("@WYou try to pick up @C$N@W but have to put them down. They are too heavy for you at the moment.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W tries to pick up @c$N@W. After struggling for a moment $e has to put $M down.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			return
		} else {
			act(libc.CString("@WYou pick up @C$N@W and put $M over your shoulder.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@W picks up $c$N@W and puts $M over $s shoulder.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if vict.Sits != nil {
				var chair *obj_data = vict.Sits
				chair.Sitting = nil
				vict.Sits = nil
			}
			ch.Player_specials.Carrying = vict
			vict.Player_specials.Carried_by = ch
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			return
		}
	}
}
func carry_drop(ch *char_data, type_ int) {
	var vict *char_data = nil
	vict = ch.Player_specials.Carrying
	switch type_ {
	case 0:
		act(libc.CString("@WYou gently set @C$N@W down on the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n @Wgently sets you down on the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n @Wgently sets @c$N@W down on the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	case 1:
		act(libc.CString("@WYou set @C$N@W hastily onto the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n @Wsets you hastily onto the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n @Wsets @c$N@W hastily onto the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	case 2:
		act(libc.CString("@WYou have @C$N@W knocked out of your arms and onto the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@WYou are knocked out of @C$n's@W arms and onto the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n @Whas @c$N@W knocked out of $s arms and onto the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	case 3:
		act(libc.CString("@WYou stop carrying @C$N@W for some reason.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@C$n @Wstops carrying you for some reason.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@C$n @Wstops carrying @c$N@W for some reason.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	}
	ch.Player_specials.Carrying = nil
	vict.Player_specials.Carried_by = nil
}
func land_location(ch *char_data, arg *byte) int {
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 50 {
		if libc.StrCaseCmp(arg, libc.CString("Nexus City")) == 0 {
			return 300
		} else if libc.StrCaseCmp(arg, libc.CString("South Ocean")) == 0 {
			return 800
		} else if libc.StrCaseCmp(arg, libc.CString("Nexus Field")) == 0 {
			return 1150
		} else if libc.StrCaseCmp(arg, libc.CString("Cherry Blossom Mountain")) == 0 {
			return 1180
		} else if libc.StrCaseCmp(arg, libc.CString("Sandy Desert")) == 0 {
			return 1287
		} else if libc.StrCaseCmp(arg, libc.CString("Northern Plains")) == 0 {
			return 1428
		} else if libc.StrCaseCmp(arg, libc.CString("Korin's Tower")) == 0 {
			return 1456
		} else if libc.StrCaseCmp(arg, libc.CString("Kami's Lookout")) == 0 {
			return 1506
		} else if libc.StrCaseCmp(arg, libc.CString("Shadow Forest")) == 0 {
			return 1636
		} else if libc.StrCaseCmp(arg, libc.CString("Decrepit Area")) == 0 {
			return 1710
		} else if libc.StrCaseCmp(arg, libc.CString("West City")) == 0 {
			return 19510
		} else if libc.StrCaseCmp(arg, libc.CString("Hercule Beach")) == 0 {
			return 2141
		} else if libc.StrCaseCmp(arg, libc.CString("Satan City")) == 0 {
			return 13020
		} else {
			send_to_char(ch, libc.CString("You don't know where that made up place is, but decided to land anyway."))
			return 300
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 51 {
		if libc.StrCaseCmp(arg, libc.CString("Ice Crown City")) == 0 {
			return 4264
		} else if libc.StrCaseCmp(arg, libc.CString("Ice Highway")) == 0 {
			return 4300
		} else if libc.StrCaseCmp(arg, libc.CString("Topica Snowfield")) == 0 {
			return 4351
		} else if libc.StrCaseCmp(arg, libc.CString("Glug's Volcano")) == 0 {
			return 4400
		} else if libc.StrCaseCmp(arg, libc.CString("Platonic Sea")) == 0 {
			return 4600
		} else if libc.StrCaseCmp(arg, libc.CString("Slave City")) == 0 {
			return 4800
		} else if libc.StrCaseCmp(arg, libc.CString("Acturian Woods")) == 0 {
			return 5100
		} else if libc.StrCaseCmp(arg, libc.CString("Desolate Demesne")) == 0 {
			return 5150
		} else if libc.StrCaseCmp(arg, libc.CString("Chateau Ishran")) == 0 {
			return 5165
		} else if libc.StrCaseCmp(arg, libc.CString("Wyrm Spine Mountain")) == 0 {
			return 5200
		} else if libc.StrCaseCmp(arg, libc.CString("Cloud Ruler Temple")) == 0 {
			return 5500
		} else if libc.StrCaseCmp(arg, libc.CString("Koltoan Mine")) == 0 {
			return 4944
		} else {
			send_to_char(ch, libc.CString("You don't know where that made up place is, but decided to land anyway."))
			return 4264
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 52 {
		if libc.StrCaseCmp(arg, libc.CString("Tiranoc City")) == 0 {
			return 8006
		} else if libc.StrCaseCmp(arg, libc.CString("Great Oroist Temple")) == 0 {
			return 8300
		} else if libc.StrCaseCmp(arg, libc.CString("Elzthuan Forest")) == 0 {
			return 8400
		} else if libc.StrCaseCmp(arg, libc.CString("Mazori Farm")) == 0 {
			return 8447
		} else if libc.StrCaseCmp(arg, libc.CString("Dres")) == 0 {
			return 8500
		} else if libc.StrCaseCmp(arg, libc.CString("Colvian Farm")) == 0 {
			return 8600
		} else if libc.StrCaseCmp(arg, libc.CString("St Alucia")) == 0 {
			return 8700
		} else if libc.StrCaseCmp(arg, libc.CString("Meridius Memorial")) == 0 {
			return 8800
		} else if libc.StrCaseCmp(arg, libc.CString("Desert of Illusion")) == 0 {
			return 8900
		} else if libc.StrCaseCmp(arg, libc.CString("Plains of Confusion")) == 0 {
			return 8954
		} else if libc.StrCaseCmp(arg, libc.CString("Turlon Fair")) == 0 {
			return 9200
		} else if libc.StrCaseCmp(arg, libc.CString("Wetlands")) == 0 {
			return 9700
		} else if libc.StrCaseCmp(arg, libc.CString("Kerberos")) == 0 {
			return 9855
		} else if libc.StrCaseCmp(arg, libc.CString("Shaeras Mansion")) == 0 {
			return 9864
		} else if libc.StrCaseCmp(arg, libc.CString("Slavinus Ravine")) == 0 {
			return 9900
		} else if libc.StrCaseCmp(arg, libc.CString("Furian Citadel")) == 0 {
			return 9949
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 8006
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 53 {
		if libc.StrCaseCmp(arg, libc.CString("Vegetos City")) == 0 {
			return 2226
		} else if libc.StrCaseCmp(arg, libc.CString("Blood Dunes")) == 0 {
			return 2600
		} else if libc.StrCaseCmp(arg, libc.CString("Ancestral Mountains")) == 0 {
			return 2616
		} else if libc.StrCaseCmp(arg, libc.CString("Destopa Swamp")) == 0 {
			return 2709
		} else if libc.StrCaseCmp(arg, libc.CString("Pride forest")) == 0 {
			return 2800
		} else if libc.StrCaseCmp(arg, libc.CString("Pride Tower")) == 0 {
			return 2899
		} else if libc.StrCaseCmp(arg, libc.CString("Ruby Cave")) == 0 {
			return 2615
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 2226
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 54 {
		if libc.StrCaseCmp(arg, libc.CString("Senzu Village")) == 0 {
			return 11600
		} else if libc.StrCaseCmp(arg, libc.CString("Guru's House")) == 0 {
			return 0x27C6
		} else if libc.StrCaseCmp(arg, libc.CString("Crystalline Cave")) == 0 {
			return 0x28EA
		} else if libc.StrCaseCmp(arg, libc.CString("Elder Village")) == 0 {
			return 13300
		} else if libc.StrCaseCmp(arg, libc.CString("Frieza's Ship")) == 0 {
			return 0x27DB
		} else if libc.StrCaseCmp(arg, libc.CString("Kakureta Village")) == 0 {
			return 0x2AAA
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 11600
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 55 {
		if libc.StrCaseCmp(arg, libc.CString("Haven City")) == 0 {
			return 12010
		} else if libc.StrCaseCmp(arg, libc.CString("Serenity Lake")) == 0 {
			return 0x2F47
		} else if libc.StrCaseCmp(arg, libc.CString("Kaiju Forest")) == 0 {
			return 12300
		} else if libc.StrCaseCmp(arg, libc.CString("Ortusian Temple")) == 0 {
			return 12400
		} else if libc.StrCaseCmp(arg, libc.CString("Silent Glade")) == 0 {
			return 0x30C0
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 12010
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 56 {
		if libc.StrCaseCmp(arg, libc.CString("Yardra City")) == 0 {
			return 0x36B8
		} else if libc.StrCaseCmp(arg, libc.CString("Jade Forest")) == 0 {
			return 14100
		} else if libc.StrCaseCmp(arg, libc.CString("Jade Cliffs")) == 0 {
			return 14200
		} else if libc.StrCaseCmp(arg, libc.CString("Mount Valaria")) == 0 {
			return 14300
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x36B8
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 198 {
		if libc.StrCaseCmp(arg, libc.CString("Cerria Colony")) == 0 {
			return 0x447B
		} else if libc.StrCaseCmp(arg, libc.CString("Crystalline Forest")) == 0 {
			return 7950
		} else if libc.StrCaseCmp(arg, libc.CString("Fistarl Volcano")) == 0 {
			return 17420
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x447B
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 57 {
		if libc.StrCaseCmp(arg, libc.CString("Utatlan City")) == 0 {
			return 3412
		} else if libc.StrCaseCmp(arg, libc.CString("Zenith Jungle")) == 0 {
			return 3520
		} else if libc.StrCaseCmp(arg, libc.CString("Ancient Castle")) == 0 {
			return 19600
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 3412
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 58 {
		if libc.StrCaseCmp(arg, libc.CString("Aquis City")) == 0 {
			return 0x3A38
		} else if libc.StrCaseCmp(arg, libc.CString("Yunkai Pirate Base")) == 0 {
			return 0x3D27
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x3A38
		}
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 59 {
		if libc.StrCaseCmp(arg, libc.CString("Janacre")) == 0 {
			return 0x3E89
		} else if libc.StrCaseCmp(arg, libc.CString("Arlian Wasteland")) == 0 {
			return 0x40A0
		} else if libc.StrCaseCmp(arg, libc.CString("Arlia Mine")) == 0 {
			return 16600
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x3E89
		}
	} else {
		send_to_char(ch, libc.CString("You are not above a planet!\r\n"))
		return -1
	}
}
func disp_locations(ch *char_data) {
	if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 50 {
		send_to_char(ch, libc.CString("@D------------------[ @GEarth@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Nexus City, South Ocean, Nexus field, Cherry Blossom Mountain,\n"))
		send_to_char(ch, libc.CString("Sandy Desert, Northern Plains, Korin's Tower, Kami's Lookout,\n"))
		send_to_char(ch, libc.CString("Shadow Forest, Decrepit Area, West City, Hercule Beach, Satan City.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 51 {
		send_to_char(ch, libc.CString("@D------------------[ @CFrigid@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Ice Crown City, Ice Highway, Topica Snowfield, Glug's Volcano,\n"))
		send_to_char(ch, libc.CString("Platonic Sea, Slave City, Acturian Woods, Desolate Demesne,\n"))
		send_to_char(ch, libc.CString("Chateau Ishran, Wyrm Spine Mountain, Cloud Ruler Temple, Koltoan mine.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 52 {
		send_to_char(ch, libc.CString("@D------------------[ @MKonack@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Great Oroist Temple, Elzthuan Forest, Mazori Farm, Dres,\n"))
		send_to_char(ch, libc.CString("Colvian Farm, St Alucia, Meridius Memorial, Desert of Illusion,\n"))
		send_to_char(ch, libc.CString("Plains of Confusion, Turlon Fair, Wetlands, Kerberos,\n"))
		send_to_char(ch, libc.CString("Shaeras Mansion, Slavinus Ravine, Furian Citadel.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 53 {
		send_to_char(ch, libc.CString("@D------------------[ @YVegeta@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Vegetos City, Blood Dunes, Ancestral Mountains, Destopa Swamp,\n"))
		send_to_char(ch, libc.CString("Pride Forest, Pride tower, Ruby Cave.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 198 {
		send_to_char(ch, libc.CString("@D------------------[ @MCerria@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Cerria Colony, Fistarl Volcano, Crystalline Forest.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 54 {
		send_to_char(ch, libc.CString("@D------------------[ @gNamek@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Senzu Village, Guru's House, Crystalline Cave, Elder Village,\n"))
		send_to_char(ch, libc.CString("Frieza's Ship, Kakureta Village.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 55 {
		send_to_char(ch, libc.CString("@D------------------[ @BAether@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Haven City, Serenity Lake, Kaiju Forest, Ortusian Temple,\n"))
		send_to_char(ch, libc.CString("Silent Glade.\n"))
		send_to_char(ch, libc.CString("@D--------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 56 {
		send_to_char(ch, libc.CString("@D-----------------[ @mYardrat@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Yardra City, Jade Forest, Jade Cliffs, Mount Valaria.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 57 {
		send_to_char(ch, libc.CString("@D-----------------[ @CZennith@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Utatlan City, Zenith Jungle, Ancient Castle.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 58 {
		send_to_char(ch, libc.CString("@D-----------------[ @CKanassa@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Aquis City, Yunkai Pirate Base.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) == 59 {
		send_to_char(ch, libc.CString("@D------------------[ @MArlia@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Janacre, Arlian Wasteland, Arlia Mine.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else {
		send_to_char(ch, libc.CString("You are not above a planet!\r\n"))
	}
}
func do_land(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		above_planet int = TRUE
		inroom       int = int(func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}())
	)
	skip_spaces(&argument)
	if inroom != 50 && inroom != 198 && inroom != 51 && inroom != 52 && inroom != 53 && inroom != 54 && inroom != 55 && inroom != 56 && inroom != 57 && inroom != 58 && inroom != 59 {
		above_planet = FALSE
	}
	if *argument == 0 {
		if above_planet == TRUE {
			send_to_char(ch, libc.CString("Land where?\n"))
			disp_locations(ch)
			return
		} else {
			send_to_char(ch, libc.CString("You are not even in the lower atmosphere of a planet!\r\n"))
			return
		}
	}
	var landing int = land_location(ch, argument)
	if landing != -1 {
		var was_in int = int(func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}())
		send_to_char(ch, libc.CString("You descend through the upper atmosphere, and coming down through the clouds you land quickly on the ground below.\r\n"))
		char_from_room(ch)
		char_to_room(ch, real_room(room_vnum(landing)))
		var blah *byte = sense_location(ch)
		var sendback [2048]byte
		char_from_room(ch)
		char_to_room(ch, real_room(room_vnum(was_in)))
		stdio.Sprintf(&sendback[0], "@C$n@Y flies down through the atmosphere toward @G%s@Y!@n", blah)
		act(&sendback[0], TRUE, ch, nil, nil, TO_ROOM)
		char_from_room(ch)
		char_to_room(ch, real_room(room_vnum(landing)))
		var zone int = 0
		if (func() int {
			zone = int(real_zone_by_thing(func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()))
			return zone
		}()) != int(-1) {
			fly_zone(zone_rnum(zone), libc.CString("can be seen landing from space nearby!@n\r\n"), ch)
		}
		send_to_sense(1, libc.CString("landing on the planet"), ch)
		send_to_scouter(libc.CString("A powerlevel signal has been detected landing on the planet"), ch, 0, 1)
		act(libc.CString("$n comes down from high above in the sky and quickly lands on the ground."), TRUE, ch, nil, nil, TO_ROOM)
		return
	}
}
func has_boat(ch *char_data) int {
	var (
		obj *obj_data
		i   int
	)
	if ADM_FLAGGED(ch, ADM_WALKANYWHERE) || ch.Admlevel > 4 {
		return 1
	}
	if AFF_FLAGGED(ch, AFF_WATERWALK) {
		return 1
	}
	for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
		if int(obj.Type_flag) == ITEM_BOAT && find_eq_pos(ch, obj, nil) < 0 {
			return 1
		}
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil && int((ch.Equipment[i]).Type_flag) == ITEM_BOAT {
			return 1
		}
	}
	return 0
}
func has_flight(ch *char_data) int {
	var obj *obj_data
	if ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		return 1
	}
	if AFF_FLAGGED(ch, AFF_FLYING) && ch.Mana >= int64(GET_LEVEL(ch)+int(ch.Max_mana/int64(GET_LEVEL(ch)*30))) && int(ch.Race) != RACE_ANDROID && !IS_NPC(ch) {
		return 1
	}
	if AFF_FLAGGED(ch, AFF_FLYING) && ch.Mana < int64(GET_LEVEL(ch)+int(ch.Max_mana/int64(GET_LEVEL(ch)*30))) && int(ch.Race) != RACE_ANDROID && !IS_NPC(ch) {
		act(libc.CString("@WYou crash to the ground, too tired to fly anymore!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n@W crashes to the ground!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
		handle_fall(ch)
		return 0
	}
	if AFF_FLAGGED(ch, AFF_FLYING) && int(ch.Race) == RACE_ANDROID {
		return 1
	}
	if AFF_FLAGGED(ch, AFF_FLYING) && IS_NPC(ch) {
		return 1
	}
	for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
		if OBJAFF_FLAGGED(obj, AFF_FLYING) && find_eq_pos(ch, obj, nil) < 0 {
			return 1
		}
	}
	return 0
}
func has_o2(ch *char_data) int {
	if ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		return 1
	}
	if AFF_FLAGGED(ch, AFF_WATERBREATH) {
		return 1
	}
	if int(ch.Race) == RACE_KANASSAN || int(ch.Race) == RACE_ANDROID || int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_MAJIN {
		return 1
	}
	return 0
}
func do_simple_move(ch *char_data, dir int, need_specials_check int) int {
	var (
		throwaway [2048]byte = func() [2048]byte {
			var t [2048]byte
			copy(t[:], []byte(""))
			return t
		}()
		buf2          [64936]byte
		buf3          [64936]byte
		was_in        room_rnum = ch.In_room
		need_movement int
		rm            *room_data
	)
	if need_specials_check != 0 && special(ch, dir+1, &throwaway[0]) != 0 {
		return 0
	}
	if leave_mtrigger(ch, dir) == 0 || ch.In_room != was_in {
		return 0
	}
	if leave_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, dir) == 0 || ch.In_room != was_in {
		return 0
	}
	if leave_otrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, dir) == 0 || ch.In_room != was_in {
		return 0
	}
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master != nil && ch.In_room == ch.Master.In_room {
		send_to_char(ch, libc.CString("The thought of leaving your master makes you weep.\r\n"))
		act(libc.CString("$n bursts into tears."), FALSE, ch, nil, nil, TO_ROOM)
		return 0
	}
	var willfall int = FALSE
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_FLYING || (func() int {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_FLYING {
		if has_flight(ch) == 0 {
			if dir != 4 {
				willfall = TRUE
			} else {
				send_to_char(ch, libc.CString("You need to fly to go there!\r\n"))
				return 0
			}
		}
	}
	if ((func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM || (func() int {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM) && IS_HUMANOID(ch) {
		if int(ch.Race) == RACE_KANASSAN && has_flight(ch) == 0 {
			act(libc.CString("@CYou swim swiftly.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C swims swiftly.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if int(ch.Race) == RACE_ICER && has_flight(ch) == 0 {
			act(libc.CString("@CYou swim swiftly.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@C swims swiftly.@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if int(ch.Race) != RACE_KANASSAN && int(ch.Race) != RACE_ICER && has_flight(ch) == 0 {
			if check_swim(ch) == 0 {
				return 0
			} else {
				act(libc.CString("@CYou swim through the cold water.@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@C swim through the cold water.@n"), TRUE, ch, nil, nil, TO_ROOM)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			}
		}
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		if int(ch.Race) != RACE_ANDROID {
			if check_swim(ch) == 0 {
				return 0
			}
		}
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Geffect == 6 && !IS_HUMANOID(ch) && IS_NPC(ch) {
		return 0
	}
	if IS_NPC(ch) && ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room, ROOM_NOMOB) && ch.Master == nil {
		return 0
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < 0 || (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_UNDERWATER || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Geffect < 0 || (func() int {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_UNDERWATER) {
		if has_o2(ch) == 0 && (group_bonus(ch, 2) != 10 && ch.Mana < ch.Max_mana/200 || group_bonus(ch, 2) == 10 && ch.Mana < ch.Max_mana/800) {
			if ch.Hit >= ch.Max_hit/20 {
				send_to_char(ch, libc.CString("@RYou struggle to breath!@n\r\n"))
				ch.Hit -= ch.Max_hit / 20
			}
			if ch.Hit < ch.Max_hit/20 {
				send_to_char(ch, libc.CString("@rYou drown!@n\r\n"))
				die(ch, nil)
				return 0
			}
		}
		if has_o2(ch) == 0 && (group_bonus(ch, 2) != 10 && ch.Mana >= ch.Max_mana/200 || group_bonus(ch, 2) == 10 && ch.Mana >= ch.Max_mana/800) {
			send_to_char(ch, libc.CString("@CYou hold your breath!@n\r\n"))
			if group_bonus(ch, 2) == 10 {
				ch.Mana -= ch.Max_mana / 800
			} else {
				ch.Mana -= ch.Max_mana / 200
			}
		}
	}
	need_movement = 1
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity > 10 {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && int(ch.Chclass) != CLASS_BARDOCK && !IS_NPC(ch) {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	}
	if GET_LEVEL(ch) <= 1 {
		need_movement = 0
	}
	if AFF_FLAGGED(ch, AFF_HIDE) {
		if roll_skill(ch, SKILL_HIDE) > 15 {
			need_movement *= 2
		} else {
			need_movement *= 4
		}
	}
	if AFF_FLAGGED(ch, AFF_SNEAK) {
		if roll_skill(ch, SKILL_MOVE_SILENTLY) > 15 {
			need_movement *= int(1.2)
		} else {
			need_movement *= 2
		}
	}
	var flight_cost int = 0
	if AFF_FLAGGED(ch, AFF_FLYING) && int(ch.Race) != RACE_ANDROID {
		if GET_SKILL(ch, SKILL_CONCENTRATION) == 0 && GET_SKILL(ch, SKILL_FOCUS) == 0 {
			flight_cost = int(ch.Max_mana / 100)
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) != 0 && GET_SKILL(ch, SKILL_FOCUS) == 0 {
			flight_cost = int(ch.Max_mana / int64(GET_SKILL(ch, SKILL_CONCENTRATION)*2))
		} else if GET_SKILL(ch, SKILL_CONCENTRATION) == 0 && GET_SKILL(ch, SKILL_FOCUS) != 0 {
			flight_cost = int(ch.Max_mana / int64(GET_SKILL(ch, SKILL_FOCUS)*3))
		} else {
			flight_cost = int(ch.Max_mana / int64((GET_SKILL(ch, SKILL_CONCENTRATION)*2)+GET_SKILL(ch, SKILL_FOCUS)*3))
		}
	}
	if AFF_FLAGGED(ch, AFF_FLYING) && ch.Mana < int64(flight_cost) && int(ch.Race) != RACE_ANDROID {
		ch.Mana = 0
		act(libc.CString("@WYou crash to the ground, too tired to fly anymore!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n@W crashes to the ground!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
	} else if AFF_FLAGGED(ch, AFF_FLYING) && int(ch.Race) != RACE_ANDROID {
		ch.Mana -= int64(flight_cost)
	}
	if ch.Move < int64(need_movement) && !AFF_FLAGGED(ch, AFF_FLYING) && !IS_NPC(ch) {
		if need_specials_check != 0 && ch.Master != nil {
			send_to_char(ch, libc.CString("You are too exhausted to follow.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You are too exhausted.\r\n"))
		}
		return 0
	}
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dcskill != 0 {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dcmove > roll_skill(ch, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dcskill) {
			send_to_char(ch, libc.CString("Your skill in %s isn't enough to move that way!\r\n"), spell_info[((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dcskill].Name)
			if !ADM_FLAGGED(ch, ADM_WALKANYWHERE) && !IS_NPC(ch) && !AFF_FLAGGED(ch, AFF_FLYING) {
				ch.Move -= int64(need_movement)
			}
			return 0
		} else {
			send_to_char(ch, libc.CString("Your skill in %s aids in your movement.\r\n"), spell_info[((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Dcskill].Name)
		}
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_ATRIUM) {
		if House_can_enter(ch, func() room_vnum {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Number
			}
			return -1
		}()) == 0 {
			send_to_char(ch, libc.CString("That's private property -- no trespassing!\r\n"))
			return 0
		}
	}
	if ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room, ROOM_TUNNEL) && num_pc_in_room((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))) >= config_info.Play.Tunnel_size {
		if config_info.Play.Tunnel_size > 1 {
			send_to_char(ch, libc.CString("There isn't enough room for you to go there!\r\n"))
		} else {
			send_to_char(ch, libc.CString("There isn't enough room there for more than one person!\r\n"))
		}
		return 0
	}
	if ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room, ROOM_GODROOM) && ch.Admlevel < ADMLVL_GRGOD {
		send_to_char(ch, libc.CString("You aren't godly enough to use that room!\r\n"))
		return 0
	}
	rm = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))
	if !IS_NPC(ch) && ch.Admlevel < ADMLVL_IMMORT && GET_LEVEL(ch) < (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rm.Zone)))).Min_level && (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rm.Zone)))).Min_level > 0 {
		send_to_char(ch, libc.CString("Sorry, you are too low a level to enter this zone.\r\n"))
		return 0
	}
	if ch.Admlevel < ADMLVL_IMMORT && GET_LEVEL(ch) > (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rm.Zone)))).Max_level && (*(*zone_data)(unsafe.Add(unsafe.Pointer(zone_table), unsafe.Sizeof(zone_data{})*uintptr(rm.Zone)))).Max_level > 0 {
		send_to_char(ch, libc.CString("Sorry, you are too high a level to enter this zone.\r\n"))
		return 0
	}
	if ch.Admlevel < ADMLVL_IMMORT && ZONE_FLAGGED(rm.Zone, ZONE_CLOSED) {
		send_to_char(ch, libc.CString("This zone is currently closed to mortals.\r\n"))
		return 0
	}
	if ch.Admlevel >= ADMLVL_IMMORT && ch.Admlevel < ADMLVL_GRGOD && ZONE_FLAGGED(rm.Zone, ZONE_NOIMMORT) {
		send_to_char(ch, libc.CString("This zone is closed to all.\r\n"))
		return 0
	}
	if ch.Admlevel >= ADMLVL_IMMORT && ch.Admlevel < ADMLVL_GOD && can_edit_zone(ch, rm.Zone) == 0 && ZONE_FLAGGED(rm.Zone, ZONE_QUEST) {
		send_to_char(ch, libc.CString("This is a Quest zone.\r\n"))
		return 0
	}
	if !ADM_FLAGGED(ch, ADM_WALKANYWHERE) && !IS_NPC(ch) && !AFF_FLAGGED(ch, AFF_FLYING) {
		ch.Move -= int64(need_movement)
	}
	if AFF_FLAGGED(ch, AFF_SNEAK) && !IS_NPC(ch) {
		stdio.Sprintf(&buf2[0], "$n sneaks %s.", dirs[dir])
		if GET_SKILL(ch, SKILL_MOVE_SILENTLY) != 0 {
			improve_skill(ch, SKILL_MOVE_SILENTLY, 0)
		} else if slot_count(ch)+1 > ch.Skill_slots {
			send_to_char(ch, libc.CString("@RYour skill slots are full. You can not learn Move Silently.\r\n"))
			ch.Affected_by[int(AFF_SNEAK/32)] &= ^(1 << (int(AFF_SNEAK % 32)))
		} else {
			send_to_char(ch, libc.CString("@GYou learn the very basics of moving silently.@n\r\n"))
			for {
				ch.Skills[SKILL_MOVE_SILENTLY] = int8(rand_number(5, 10))
				if true {
					break
				}
			}
			act(&buf2[0], TRUE, ch, nil, nil, int(TO_ROOM|2<<9))
			if int(ch.Aff_abils.Dex) < rand_number(1, 30) {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			}
		}
	}
	if !AFF_FLAGGED(ch, AFF_SNEAK) && !AFF_FLAGGED(ch, AFF_FLYING) {
		stdio.Sprintf(&buf2[0], "$n leaves %s.", dirs[dir])
		act(&buf2[0], TRUE, ch, nil, nil, TO_ROOM)
	}
	if !AFF_FLAGGED(ch, AFF_SNEAK) && AFF_FLAGGED(ch, AFF_FLYING) {
		stdio.Sprintf(&buf2[0], "$n flies %s.", dirs[dir])
		act(&buf2[0], TRUE, ch, nil, nil, TO_ROOM)
	}
	was_in = ch.In_room
	if ch.Drag != nil {
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
	}
	ch.Affected_by[int(AFF_PURSUIT/32)] |= 1 << (int(AFF_PURSUIT % 32))
	char_from_room(ch)
	char_to_room(ch, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room)
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Zone != (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Zone && !IS_NPC(ch) && int(ch.Race) != RACE_ANDROID {
		send_to_sense(0, libc.CString("You sense someone"), ch)
		stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@Y %s\r\n@RSomeone has entered your scouter detection range@n.", add_commas(ch.Hit))
		send_to_scouter(&buf3[0], ch, 0, 0)
	}
	if entry_mtrigger(ch) == 0 || enter_wtrigger((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room))), ch, dir) == 0 {
		char_from_room(ch)
		char_to_room(ch, was_in)
		ch.Affected_by[int(AFF_PURSUIT/32)] &= ^(1 << (int(AFF_PURSUIT % 32)))
		return 0
	}
	stdio.Snprintf(&buf2[0], int(64936), "%s%s", func() string {
		if dir == UP || dir == DOWN {
			return ""
		}
		return "the "
	}(), func() string {
		if dir == UP {
			return "below"
		}
		if dir == DOWN {
			return "above"
		}
		return libc.GoString(dirs[rev_dir[dir]])
	}())
	act(libc.CString("$n arrives from $T."), TRUE, ch, nil, unsafe.Pointer(&buf2[0]), int(TO_ROOM|2<<9))
	if ch.Fighting != nil {
		if (func() int {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_FLYING && (func() int {
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room != room_rnum(-1) && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_WATER_NOSWIM && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room)))).Geffect == 0 {
			roll_pursue(ch.Fighting, ch)
		}
		ch.Affected_by[int(AFF_PURSUIT/32)] &= ^(1 << (int(AFF_PURSUIT % 32)))
	}
	if ch.Drag != nil {
		act(libc.CString("@wYou drag @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		char_from_room(ch.Drag)
		char_to_room(ch.Drag, ch.In_room)
		if ch.Drag.Sits != nil {
			obj_from_room(ch.Drag.Sits)
			obj_to_room(ch.Drag.Sits, ch.In_room)
		}
		if !AFF_FLAGGED(ch.Drag, AFF_KNOCKED) && !AFF_FLAGGED(ch.Drag, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Drag, libc.CString("You feel your sleeping body being moved.\r\n"))
			if IS_NPC(ch.Drag) && ch.Drag.Fighting == nil {
				set_fighting(ch.Drag, ch)
			}
		}
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@wYou carry @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_CHAR)
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
		char_from_room(ch.Player_specials.Carrying)
		char_to_room(ch.Player_specials.Carrying, ch.In_room)
		if !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_KNOCKED) && !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Player_specials.Carrying, libc.CString("You feel your sleeping body being moved.\r\n"))
		}
	}
	if ch.Desc != nil {
		look_at_room(ch.In_room, ch, 0)
		if AFF_FLAGGED(ch, AFF_SNEAK) && !IS_NPC(ch) && GET_SKILL(ch, SKILL_MOVE_SILENTLY) != 0 && GET_SKILL(ch, SKILL_MOVE_SILENTLY) < rand_number(1, 101) {
			send_to_char(ch, libc.CString("@wYou make a noise as you arrive and are no longer sneaking!@n\r\n"))
			act(libc.CString("@c$n@w makes a noise revealing $s sneaking!@n"), TRUE, ch, nil, nil, int(TO_ROOM|2<<9))
			reveal_hiding(ch, 0)
			ch.Affected_by[int(AFF_SNEAK/32)] &= ^(1 << (int(AFF_SNEAK % 32)))
		}
	}
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == 6 || (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Geffect == 6 {
		if int(ch.Race) != RACE_DEMON && !AFF_FLAGGED(ch, AFF_FLYING) && group_bonus(ch, 2) != 14 {
			act(libc.CString("@rYour legs are burned by the lava!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n@r's legs are burned by the lava!@n"), TRUE, ch, nil, nil, TO_ROOM)
			if IS_NPC(ch) && IS_HUMANOID(ch) && rand_number(1, 2) == 2 {
				do_fly(ch, nil, 0, 0)
			}
			if ch.Suppressed >= ch.Max_hit/20 {
				ch.Suppressed -= ch.Max_hit / 20
			} else {
				ch.Suppression = 0
				ch.Hit -= (ch.Max_hit / 20) - ch.Suppressed
				ch.Suppressed = 0
				if ch.Hit < 0 {
					act(libc.CString("@rYou have burned to death!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n@r has burned to death!@n"), TRUE, ch, nil, nil, TO_ROOM)
					die(ch, nil)
				}
			}
		}
		if ch.Drag != nil && int(ch.Drag.Race) != RACE_DEMON {
			act(libc.CString("@R$N@r gets burned!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
			act(libc.CString("@R$N@r gets burned!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
			ch.Drag.Hit -= ch.Drag.Max_hit / 20
			if ch.Drag.Hit < 0 {
				act(libc.CString("@rYou have burned to death!@n"), TRUE, ch.Drag, nil, nil, TO_CHAR)
				act(libc.CString("@R$n@r has burned to death!@n"), TRUE, ch.Drag, nil, nil, TO_ROOM)
				die(ch.Drag, nil)
			}
		}
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_TIMED_DT) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		timed_dt(nil)
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_DEATH) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		log_death_trap(ch)
		death_cry(ch)
		extract_char(ch)
		return 0
	}
	entry_memory_mtrigger(ch)
	if greet_mtrigger(ch, dir) == 0 {
		char_from_room(ch)
		char_to_room(ch, was_in)
		look_at_room(ch.In_room, ch, 0)
	} else {
		greet_memory_mtrigger(ch)
	}
	if willfall == TRUE {
		handle_fall(ch)
		if ch.Drag != nil {
			handle_fall(ch.Drag)
		}
	}
	return 1
}
func perform_move(ch *char_data, dir int, need_specials_check int) int {
	var (
		was_in room_rnum
		k      *follow_type
		next   *follow_type
	)
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return 0
	}
	if ch.Absorbing != nil || ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are struggling with someone!\r\n"))
		return 0
	}
	if !AFF_FLAGGED(ch, AFF_SNEAK) || AFF_FLAGGED(ch, AFF_SNEAK) && GET_SKILL(ch, SKILL_MOVE_SILENTLY) < axion_dice(0) {
		reveal_hiding(ch, 0)
	}
	if ch == nil || dir < 0 || dir >= NUM_OF_DIRS {
		return 0
	} else if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]) == nil && buildwalk(ch, dir) == 0 || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room == room_rnum(-1) || EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<4) && EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<1) {
		send_to_char(ch, libc.CString("Alas, you cannot go that way...\r\n"))
	} else if EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir], 1<<1) {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword != nil {
			send_to_char(ch, libc.CString("The %s seems to be closed.\r\n"), fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).Keyword))
		} else {
			send_to_char(ch, libc.CString("It seems to be closed.\r\n"))
		}
	} else if (func() room_vnum {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Number
		}
		return -1
	}()) == 0 || (func() room_vnum {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[dir]).To_room)))).Number
		}
		return -1
	}()) == 1 {
		send_to_char(ch, libc.CString("Report this direction, it is illegal.\r\n"))
	} else {
		var wall *obj_data
		for wall = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; wall != nil; wall = wall.Next_content {
			if GET_OBJ_VNUM(wall) == 79 {
				if wall.Cost == dir {
					send_to_char(ch, libc.CString("That direction has a glacial wall blocking it.\r\n"))
					return 0
				}
			}
		}
		if ch.Followers == nil {
			return do_simple_move(ch, dir, need_specials_check)
		}
		was_in = ch.In_room
		if do_simple_move(ch, dir, need_specials_check) == 0 {
			return 0
		}
		for k = ch.Followers; k != nil; k = next {
			next = k.Next
			if k.Follower.In_room == was_in && int(k.Follower.Position) >= POS_STANDING && (!AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_GROUP) && AFF_FLAGGED(k.Follower, AFF_GROUP)) {
				act(libc.CString("You follow $N.\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
				perform_move(k.Follower, dir, 1)
			} else if k.Follower.In_room == was_in && int(k.Follower.Position) >= POS_STANDING && (AFF_FLAGGED(ch, AFF_ZANZOKEN) && AFF_FLAGGED(k.Follower, AFF_ZANZOKEN)) && (!AFF_FLAGGED(ch, AFF_GROUP) || !AFF_FLAGGED(k.Follower, AFF_GROUP)) {
				act(libc.CString("$N tries to zanzoken and escape, but your zanzoken matches $S!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
				act(libc.CString("$N tries to zanzoken and escape, but $n's zanzoken matches $S!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_NOTVICT)
				act(libc.CString("You zanzoken to try and escape, but $n's zanzoken matches yours!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_VICT)
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				k.Follower.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				perform_move(k.Follower, dir, 1)
			} else if k.Follower.In_room == was_in && int(k.Follower.Position) >= POS_STANDING && (AFF_FLAGGED(ch, AFF_ZANZOKEN) && !AFF_FLAGGED(k.Follower, AFF_ZANZOKEN)) {
				act(libc.CString("You try to follow $N, but $E disappears in a flash of movement!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
				act(libc.CString("$n tries to follow $N, but $E disappears in a flash of movement!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_NOTVICT)
				act(libc.CString("$n tries to follow you, but you manage to zanzoken away!\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_VICT)
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		return 1
	}
	return 0
}
func do_move(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		perform_move(ch, subcmd-1, 0)
		return
	}
	if PLR_FLAGGED(ch, PLR_SELFD) {
		send_to_char(ch, libc.CString("You are preparing to blow up!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
		send_to_char(ch, libc.CString("You are liquefied right now!\r\n"))
		return
	}
	if float64(ch.Charge) >= float64(ch.Max_mana)*0.51 {
		send_to_char(ch, libc.CString("You have too much ki charged. You can't concentrate on keeping it charged while also traveling.\r\n"))
		return
	} else if float64(ch.Charge) >= float64(ch.Max_mana)*0.5 && float64(ch.Charge) < float64(ch.Max_mana)*0.51 && GET_SKILL(ch, SKILL_CONCENTRATION) < 100 {
		send_to_char(ch, libc.CString("You have too much ki charged. You can't concentrate on keeping it charged while also traveling.\r\n"))
		return
	} else if float64(ch.Charge) >= float64(ch.Max_mana)*0.4 && float64(ch.Charge) < float64(ch.Max_mana)*0.5 && GET_SKILL(ch, SKILL_CONCENTRATION) < 80 {
		send_to_char(ch, libc.CString("You have too much ki charged. You can't concentrate on keeping it charged while also traveling.\r\n"))
		return
	} else if float64(ch.Charge) >= float64(ch.Max_mana)*0.3 && float64(ch.Charge) < float64(ch.Max_mana)*0.4 && GET_SKILL(ch, SKILL_CONCENTRATION) < 70 {
		send_to_char(ch, libc.CString("You have too much ki charged. You can't concentrate on keeping it charged while also traveling.\r\n"))
		return
	} else if float64(ch.Charge) >= float64(ch.Max_mana)*0.2 && float64(ch.Charge) < float64(ch.Max_mana)*0.3 && GET_SKILL(ch, SKILL_CONCENTRATION) < 60 {
		send_to_char(ch, libc.CString("You have too much ki charged. You can't concentrate on keeping it charged while also traveling.\r\n"))
		return
	}
	if int(ch.Player_specials.Conditions[DRUNK]) > 4 && (rand_number(1, 9)+int(ch.Player_specials.Conditions[DRUNK])) >= rand_number(14, 20) {
		send_to_char(ch, libc.CString("You wobble around and then fall on your ass.\r\n"))
		act(libc.CString("@C$n@W wobbles around before falling on $s ass@n."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Position = POS_SITTING
		return
	}
	if ch.Fighting != nil && !IS_NPC(ch) {
		var blah [2048]byte
		stdio.Sprintf(&blah[0], "%s", dirs[subcmd-1])
		do_flee(ch, &blah[0], 0, 0)
		return
	}
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		var (
			vehicle  *obj_data = nil
			controls *obj_data = nil
			noship   int       = FALSE
		)
		if (func() *obj_data {
			controls = find_control(ch)
			return controls
		}()) == nil && ch.Admlevel < 1 {
			noship = TRUE
		} else if (func() *obj_data {
			vehicle = find_vehicle_by_vnum(controls.Value[0])
			return vehicle
		}()) == nil {
			noship = TRUE
		}
		if noship == TRUE {
			send_to_char(ch, libc.CString("Your ship controls are not here or your ship was not found, report to Iovan!\r\n"))
			return
		} else if controls != nil && vehicle != nil {
			if (controls.Value[2]) <= 0 {
				send_to_char(ch, libc.CString("The ship is out of fuel!\r\n"))
				return
			}
			drive_in_direction(ch, vehicle, subcmd-1)
			if (controls.Value[1]) == 1 {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			} else if (controls.Value[1]) == 2 {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			}
			controls = nil
			vehicle = nil
			return
		}
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if !IS_NPC(ch) {
		var (
			fail     int = FALSE
			obj      *obj_data
			next_obj *obj_data
		)
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj.Kicharge > 0 && obj.User == ch {
				fail = TRUE
			}
		}
		if fail == TRUE {
			send_to_char(ch, libc.CString("You are too busy controlling your attack!\r\n"))
			return
		}
	}
	if !IS_NPC(ch) && (ch.Limb_condition[0]) <= 0 && (ch.Limb_condition[1]) <= 0 && (ch.Limb_condition[2]) <= 0 && (ch.Limb_condition[3]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
		send_to_char(ch, libc.CString("Unless you fly, you can't get far with no limbs.\r\n"))
		return
	}
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil {
		send_to_char(ch, libc.CString("You are busy absorbing from %s!\r\n"), GET_NAME(ch.Absorbing))
		return
	}
	if ch.Absorbby != nil {
		if axion_dice(0) < GET_SKILL(ch.Absorbby, SKILL_ABSORB) {
			send_to_char(ch, libc.CString("You are being held by %s, they are absorbing you!\r\n"), GET_NAME(ch.Absorbby))
			send_to_char(ch.Absorbby, libc.CString("%s struggles in your grasp!\r\n"), GET_NAME(ch))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			return
		} else {
			act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_NOTVICT)
			act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_CHAR)
			ch.Absorbby.Absorbing = nil
			ch.Absorbby = nil
		}
	}
	if block_calc(ch) == 0 {
		return
	}
	if ch.Listenroom > 0 {
		send_to_char(ch, libc.CString("You stop eavesdropping.\r\n"))
		ch.Listenroom = room_vnum(real_room(0))
	}
	if !IS_NPC(ch) {
		if PRF_FLAGGED(ch, PRF_ARENAWATCH) {
			ch.Player_specials.Pref[int(PRF_ARENAWATCH/32)] &= bitvector_t(int32(^(1 << (int(PRF_ARENAWATCH % 32)))))
			ch.Arenawatch = -1
		}
		if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != room_vnum(-1) && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != 0 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) != 1 {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				ch.Player_specials.Load_room = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			} else {
				ch.Player_specials.Load_room = -1
			}
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && ch.Max_hit <= 10000 && int(ch.Chclass) != CLASS_BARDOCK && !IS_NPC(ch) {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 20 && ch.Max_hit <= 30000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 30 && ch.Max_hit <= 100000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 40 && ch.Max_hit <= 200000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 50 && ch.Max_hit <= 300000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 100 && ch.Max_hit <= 500000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 200 && ch.Max_hit <= 1000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 && ch.Max_hit <= 8000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 && ch.Max_hit <= 15000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 && ch.Max_hit <= 25000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 && ch.Max_hit <= 35000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 && ch.Max_hit <= 100000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		}
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 && ch.Max_hit <= 200000000 {
			send_to_char(ch, libc.CString("The gravity slows you down some.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) && ch.Admlevel < 1 {
			send_to_char(ch, libc.CString("You struggle to cross the vast distance.\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		} else if (ch.Limb_condition[2]) <= 0 && (ch.Limb_condition[3]) <= 0 && (ch.Limb_condition[0]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
			act(libc.CString("@wYou slowly pull yourself along with your arm...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w slowly pulls $mself along with one arm...@n"), TRUE, ch, nil, nil, TO_ROOM)
			if (ch.Limb_condition[1]) < 50 {
				send_to_char(ch, libc.CString("@RYour left arm is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[1] -= rand_number(1, 5)
				if (ch.Limb_condition[0]) <= 0 {
					act(libc.CString("@RYour left arm falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R left arm falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if (ch.Limb_condition[2]) <= 0 && (ch.Limb_condition[3]) <= 0 && (ch.Limb_condition[1]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
			act(libc.CString("@wYou slowly pull yourself along with your arm...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w slowly pulls $mself along with one arm...@n"), TRUE, ch, nil, nil, TO_ROOM)
			if (ch.Limb_condition[0]) < 50 {
				send_to_char(ch, libc.CString("@RYour right arm is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[0] -= rand_number(1, 5)
				if (ch.Limb_condition[0]) <= 0 {
					act(libc.CString("@RYour right arm falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R right arm falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if (ch.Limb_condition[2]) <= 0 && (ch.Limb_condition[3]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
			act(libc.CString("@wYou slowly pull yourself along with your arms...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w slowly pulls $mself along with one arms...@n"), TRUE, ch, nil, nil, TO_ROOM)
			if (ch.Limb_condition[1]) < 50 {
				send_to_char(ch, libc.CString("@RYour left arm is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[1] -= rand_number(1, 5)
				if (ch.Limb_condition[1]) <= 0 {
					act(libc.CString("@RYour left arm falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R left arm falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			if (ch.Limb_condition[0]) < 50 {
				send_to_char(ch, libc.CString("@RYour right arm is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[0] -= rand_number(1, 5)
				if (ch.Limb_condition[0]) <= 0 {
					act(libc.CString("@RYour right arm falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R right arm falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if (ch.Limb_condition[2]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
			act(libc.CString("@wYou hop on one leg...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w hops on one leg...@n"), TRUE, ch, nil, nil, TO_ROOM)
			if (ch.Limb_condition[3]) < 50 {
				send_to_char(ch, libc.CString("@RYour left leg is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[3] -= rand_number(1, 5)
				if (ch.Limb_condition[3]) <= 0 {
					act(libc.CString("@RYour left leg falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R left leg falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if (ch.Limb_condition[3]) <= 0 && !AFF_FLAGGED(ch, AFF_FLYING) {
			act(libc.CString("@wYou hop on one leg...@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w hops on one leg...@n"), TRUE, ch, nil, nil, TO_ROOM)
			if (ch.Limb_condition[2]) < 50 {
				send_to_char(ch, libc.CString("@RYour right leg is damaged by the forced use!@n\r\n"))
				ch.Limb_condition[2] -= rand_number(1, 5)
				if (ch.Limb_condition[2]) <= 0 {
					act(libc.CString("@RYour right leg falls apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@r$n's@R right leg falls apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if int(ch.Position) == POS_RESTING {
			act(libc.CString("@wYou crawl on your hands and knees.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w crawls on $s hands and knees.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Sits != nil {
				var chair *obj_data = ch.Sits
				chair.Sitting = nil
				ch.Sits = nil
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if int(ch.Position) == POS_SITTING {
			act(libc.CString("@wYou shuffle on your hands and knees.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@w shuffles on $s hands and knees.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Sits != nil {
				var chair *obj_data = ch.Sits
				chair.Sitting = nil
				ch.Sits = nil
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else if int(ch.Position) < POS_RESTING {
			send_to_char(ch, libc.CString("You are in no condition to move! Try standing...\r\n"))
			return
		}
	}
	perform_move(ch, subcmd-1, 0)
	if ch.Rdisplay != nil {
		if ch.Rdisplay != libc.CString("Empty") {
			ch.Rdisplay = libc.CString("Empty")
		}
	}
}
func find_door(ch *char_data, type_ *byte, dir *byte, cmdname *byte) int {
	var door int
	if *dir != 0 {
		if (func() int {
			door = search_block(dir, &dirs[0], FALSE)
			return door
		}()) < 0 && (func() int {
			door = search_block(dir, &abbr_dirs[0], FALSE)
			return door
		}()) < 0 {
			send_to_char(ch, libc.CString("That's not a direction.\r\n"))
			return -1
		}
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
				if is_name(type_, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword) != 0 {
					return door
				} else {
					send_to_char(ch, libc.CString("I see no %s there.\r\n"), type_)
					return -1
				}
			} else {
				return door
			}
		} else {
			send_to_char(ch, libc.CString("I really don't see how you can %s anything there.\r\n"), cmdname)
			return -1
		}
	} else {
		if *type_ == 0 {
			send_to_char(ch, libc.CString("What is it you want to %s?\r\n"), cmdname)
			return -1
		}
		for door = 0; door < NUM_OF_DIRS; door++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
					if is_name(type_, ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword) != 0 {
						return door
					}
				}
			}
		}
		send_to_char(ch, libc.CString("There doesn't seem to be %s %s that could be manipulated in that way here.\r\n"), AN(type_), type_)
		return -1
	}
}
func has_key(ch *char_data, key obj_vnum) int {
	var (
		o *obj_data
		i int
	)
	if key == 1 {
		return 1
	}
	for o = ch.Carrying; o != nil; o = o.Next_content {
		if GET_OBJ_VNUM(o) == key {
			return 1
		}
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			if GET_OBJ_VNUM(ch.Equipment[i]) == key {
				return 1
			}
		}
	}
	return 0
}

var cmd_door [5]*byte = [5]*byte{libc.CString("open"), libc.CString("close"), libc.CString("unlock"), libc.CString("lock"), libc.CString("pick")}
var flags_door [5]int = [5]int{(1 << 1) | 1<<2, (1 << 0), (1 << 1) | 1<<3, (1 << 1) | 1<<2, (1 << 1) | 1<<3}

func do_doorcmd(ch *char_data, obj *obj_data, door int, scmd int) {
	var (
		buf        [64936]byte
		len_       uint64
		num        int                  = 0
		other_room room_rnum            = room_rnum(-1)
		back       *room_direction_data = nil
		hatch      *obj_data            = nil
		obj2       *obj_data            = nil
		next_obj   *obj_data
		vehicle    *obj_data = nil
	)
	if obj != nil && int(obj.Type_flag) == ITEM_HATCH {
		vehicle = find_vehicle_by_vnum(obj.Value[VAL_HATCH_DEST])
	} else if obj != nil && int(obj.Type_flag) == ITEM_VEHICLE {
		if real_room(room_vnum(obj.Value[VAL_PORTAL_DEST])) != room_rnum(-1) {
			num = int(ch.In_room)
			char_from_room(ch)
			char_to_room(ch, real_room(room_vnum(obj.Value[VAL_PORTAL_DEST])))
		}
		for obj2 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj2 != nil; obj2 = next_obj {
			next_obj = obj2.Next_content
			if int(obj2.Type_flag) == ITEM_HATCH {
				hatch = obj2
			}
		}
		obj2 = nil
	}
	if door_mtrigger(ch, scmd, door) == 0 {
		return
	}
	if door_wtrigger(ch, scmd, door) == 0 {
		return
	}
	len_ = uint64(stdio.Snprintf(&buf[0], int(64936), "$n %ss ", cmd_door[scmd]))
	if obj == nil && (func() room_rnum {
		other_room = ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room
		return other_room
	}()) != room_rnum(-1) {
		if (func() *room_direction_data {
			back = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]
			return back
		}()) != nil {
			if back.To_room != ch.In_room {
				back = nil
			}
		}
	}
	switch scmd {
	case SCMD_OPEN:
		if obj != nil {
			if obj != nil && int(obj.Type_flag) == ITEM_HATCH && vehicle != nil {
				if vehicle != nil {
					vehicle.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 2)
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 1)
				}
				if GET_OBJ_VNUM(obj) > 0x4AFF {
					send_to_room(ch.In_room, libc.CString("@wThe ship hatch opens slowly and settles onto the ground outside.\r\n"))
					send_to_room(vehicle.In_room, libc.CString("@wThe ship hatch opens slowly and settles onto the ground.\r\n"))
					if ROOM_FLAGGED(vehicle.In_room, ROOM_SPACE) {
						send_to_room(ch.In_room, libc.CString("@wA great vortex forms as air begins to get sucked out into the void!\r\n"))
					}
				} else {
					act(libc.CString("@wYou open @c$p@w."), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@w opens @c$p@w."), TRUE, ch, obj, nil, TO_ROOM)
					send_to_room(vehicle.In_room, libc.CString("@wThe door to %s@w is opened from the other side.\r\n"), vehicle.Short_description)
				}
				vehicle = nil
			}
			if obj != nil && int(obj.Type_flag) == ITEM_VEHICLE && hatch != nil {
				if hatch != nil {
					hatch.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 2)
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 1)
				}
				char_from_room(ch)
				char_to_room(ch, room_rnum(num))
				if GET_OBJ_VNUM(obj) > 0x4AFF {
					send_to_room(room_rnum(num), libc.CString("@wThe ship hatch opens slowly and settles onto the ground.\r\n"))
					send_to_room(hatch.In_room, libc.CString("@wThe ship hatch opens slowly.\r\n"))
					if ROOM_FLAGGED(obj.In_room, ROOM_SPACE) {
						send_to_room(room_rnum(num), libc.CString("@wThe air starts getting sucked out into space as the hatch opens!\r\n"))
					}
				} else {
					act(libc.CString("@wYou open @c$p@w."), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@w opens @c$p@w."), TRUE, ch, obj, nil, TO_ROOM)
					send_to_room(hatch.In_room, libc.CString("@wThe door is opened from the other side.\r\n"))
				}
				hatch = nil
			}
		}
		if obj != nil {
			obj.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 2)
		} else {
			((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 1)
		}
		if back != nil {
			if obj != nil {
				obj.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 2)
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]).Exit_info &= ^(1 << 1)
			}
		}
		if obj == nil {
			send_to_char(ch, libc.CString("You open the %s that leads %s.\r\n"), func() *byte {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
					return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword
				}
				return libc.CString("door")
			}(), dirs[door])
		} else if int(obj.Type_flag) != ITEM_VEHICLE && int(obj.Type_flag) != ITEM_HATCH {
			send_to_char(ch, libc.CString("You open %s.\r\n"), obj.Short_description)
		}
	case SCMD_CLOSE:
		if obj != nil {
			if obj != nil && int(obj.Type_flag) == ITEM_HATCH && vehicle != nil {
				if vehicle != nil {
					vehicle.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 1
				}
				if GET_OBJ_VNUM(obj) > 0x4AFF {
					send_to_room(ch.In_room, libc.CString("@wThe ship hatch slowly closes, sealing the ship from the outside.\r\n"))
					send_to_room(vehicle.In_room, libc.CString("@wThe ship hatch slowly closes, sealing the ship.\r\n"))
					if ROOM_FLAGGED(vehicle.In_room, ROOM_SPACE) {
						send_to_room(ch.In_room, libc.CString("@wThe air stops getting sucked out into space as the hatch seals!\r\n"))
					}
				} else {
					act(libc.CString("@wYou close @c$p@w."), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@w closes @c$p@w."), TRUE, ch, obj, nil, TO_ROOM)
					send_to_room(vehicle.In_room, libc.CString("@wThe door to %s@w is closed from the other side.\r\n"), vehicle.Short_description)
				}
				vehicle = nil
			}
			if obj != nil && int(obj.Type_flag) == ITEM_VEHICLE && hatch != nil {
				if hatch != nil {
					hatch.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 1
				}
				char_from_room(ch)
				char_to_room(ch, room_rnum(num))
				if GET_OBJ_VNUM(obj) > 0x4AFF {
					send_to_room(room_rnum(num), libc.CString("@wThe ship hatch slowly closes, sealing the ship.\r\n"))
					send_to_room(hatch.In_room, libc.CString("@wThe ship hatch slowly closes, sealing the ship from the outside.\r\n"))
					if ROOM_FLAGGED(obj.In_room, ROOM_SPACE) {
						send_to_room(room_rnum(num), libc.CString("@wAir stops getting sucked out into space as the hatch seals!\r\n"))
					}
				} else {
					act(libc.CString("@wYou close @c$p@w."), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("@C$n@w closes @c$p@w."), TRUE, ch, obj, nil, TO_ROOM)
					send_to_room(hatch.In_room, libc.CString("@wThe door to %s@w is closed from the other side.\r\n"), hatch.Short_description)
				}
				hatch = nil
			}
		}
		if obj != nil {
			obj.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
		} else {
			((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 1
		}
		if back != nil {
			if obj != nil {
				obj.Value[VAL_CONTAINER_FLAGS] |= 1 << 2
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]).Exit_info |= 1 << 1
			}
		}
		if obj == nil {
			send_to_char(ch, libc.CString("You close the %s that leads %s.\r\n"), func() *byte {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
					return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword
				}
				return libc.CString("door")
			}(), dirs[door])
		} else if int(obj.Type_flag) != ITEM_VEHICLE && int(obj.Type_flag) != ITEM_HATCH {
			send_to_char(ch, libc.CString("You close %s.\r\n"), obj.Short_description)
		}
	case SCMD_LOCK:
		if obj != nil {
			if obj != nil && int(obj.Type_flag) == ITEM_HATCH && vehicle != nil {
				if vehicle != nil {
					vehicle.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 2
				}
				vehicle = nil
			}
			if obj != nil && int(obj.Type_flag) == ITEM_VEHICLE && hatch != nil {
				if hatch != nil {
					hatch.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 2
				}
				char_from_room(ch)
				char_to_room(ch, room_rnum(num))
				hatch = nil
			}
		}
		if obj != nil {
			obj.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
		} else {
			((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info |= 1 << 2
		}
		if back != nil {
			if obj != nil {
				obj.Value[VAL_CONTAINER_FLAGS] |= 1 << 3
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]).Exit_info |= 1 << 2
			}
		}
		if obj == nil {
			send_to_char(ch, libc.CString("You lock the %s that leads %s.\r\n"), func() *byte {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
					return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword
				}
				return libc.CString("door")
			}(), dirs[door])
		} else {
			send_to_char(ch, libc.CString("You lock %s.\r\n"), obj.Short_description)
		}
	case SCMD_UNLOCK:
		if obj != nil {
			if obj != nil && int(obj.Type_flag) == ITEM_HATCH && vehicle != nil {
				if vehicle != nil {
					vehicle.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 3)
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 2)
				}
				vehicle = nil
			}
			if obj != nil && int(obj.Type_flag) == ITEM_VEHICLE && hatch != nil {
				if hatch != nil {
					hatch.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 3)
				} else {
					((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 2)
				}
				char_from_room(ch)
				char_to_room(ch, room_rnum(num))
				hatch = nil
			}
		}
		if obj != nil {
			obj.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 3)
		} else {
			((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info &= ^(1 << 2)
		}
		if back != nil {
			if obj != nil {
				obj.Value[VAL_CONTAINER_FLAGS] &= ^(1 << 3)
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]).Exit_info &= ^(1 << 2)
			}
		}
		if obj == nil {
			send_to_char(ch, libc.CString("You unlock the %s that leads %s.\r\n"), func() *byte {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
					return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword
				}
				return libc.CString("door")
			}(), dirs[door])
		} else {
			send_to_char(ch, libc.CString("You unlock %s.\r\n"), obj.Short_description)
		}
	case SCMD_PICK:
		if obj != nil {
			obj.Value[VAL_CONTAINER_FLAGS] ^= 1 << 3
		} else {
			((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Exit_info ^= 1 << 2
		}
		if back != nil {
			if obj != nil {
				obj.Value[VAL_CONTAINER_FLAGS] ^= 1 << 3
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(other_room)))).Dir_option[rev_dir[door]]).Exit_info ^= 1 << 2
			}
		}
		send_to_char(ch, libc.CString("The lock quickly yields to your skills.\r\n"))
		len_ = strlcpy(&buf[0], libc.CString("$n skillfully picks the lock on "), uint64(64936))
	}
	var dbuf [100]byte
	if obj == nil {
		stdio.Sprintf(&dbuf[0], "%s", dirs[door])
	}
	if len_ < uint64(64936) {
		stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%s%s%s%s.", func() string {
			if obj != nil {
				return ""
			}
			return "the "
		}(), func() string {
			if obj != nil {
				return "$p"
			}
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
				return "$F"
			}
			return "door"
		}(), func() string {
			if obj != nil {
				return ""
			}
			return " that leads "
		}(), func() string {
			if obj != nil {
				return ""
			}
			return libc.GoString(&dbuf[0])
		}())
	}
	if obj == nil || obj.In_room != room_rnum(-1) {
		act(&buf[0], FALSE, ch, obj, unsafe.Pointer(func() *byte {
			if obj != nil {
				return nil
			}
			return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword
		}()), TO_ROOM)
	}
	if back != nil && (scmd == SCMD_OPEN || scmd == SCMD_CLOSE) {
		send_to_room(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, libc.CString("The %s that leads %s is %s%s from the other side.\r\n"), func() *byte {
			if back.Keyword != nil {
				return fname(back.Keyword)
			}
			return libc.CString("door")
		}(), &dbuf[0], cmd_door[scmd], func() string {
			if scmd == SCMD_CLOSE {
				return "d"
			}
			return "ed"
		}())
	} else if back != nil && (scmd == SCMD_LOCK || scmd == SCMD_UNLOCK) {
		send_to_room(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, libc.CString("The %s that leads %s is %sed from the other side.\r\n"), func() *byte {
			if back.Keyword != nil {
				return fname(back.Keyword)
			}
			return libc.CString("door")
		}(), &dbuf[0], cmd_door[scmd])
	}
	dbuf[0] = '\x00'
}
func ok_pick(ch *char_data, keynum obj_vnum, pickproof int, dclock int, scmd int, hatch *obj_data) int {
	var (
		skill_lvl int
		found     int = FALSE
		obj       *obj_data
		next_obj  *obj_data
	)
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if GET_OBJ_VNUM(obj) == 18 && (!OBJ_FLAGGED(obj, ITEM_BROKEN) && !OBJ_FLAGGED(obj, ITEM_FORGED)) {
			found = TRUE
		}
	}
	if scmd != SCMD_PICK {
		return 1
	}
	if GET_SKILL(ch, SKILL_OPEN_LOCK) == 0 {
		send_to_char(ch, libc.CString("You have no idea how!\r\n"))
		return 0
	}
	if found == FALSE {
		send_to_char(ch, libc.CString("You need a lock picking kit.\r\n"))
		return 0
	}
	if hatch != nil && (int(hatch.Type_flag) == ITEM_HATCH || int(hatch.Type_flag) == ITEM_VEHICLE) {
		send_to_char(ch, libc.CString("No picking ship hatches.\r\n"))
		hatch = nil
		return 0
	}
	skill_lvl = roll_skill(ch, SKILL_OPEN_LOCK)
	if dclock == 0 {
		dclock = rand_number(1, 101)
	}
	if keynum == obj_vnum(-1) {
		send_to_char(ch, libc.CString("Odd - you can't seem to find a keyhole.\r\n"))
	} else if pickproof != 0 {
		send_to_char(ch, libc.CString("It resists your attempts to pick it.\r\n"))
		act(libc.CString("@c$n@w puts a set of lockpick tools away.@n"), TRUE, ch, nil, nil, TO_ROOM)
	} else if ch.Move < ch.Max_move/30 {
		send_to_char(ch, libc.CString("You don't have the stamina to try, it takes percision to pick locks.Not shaking tired hands.\r\n"))
	} else if dclock > (skill_lvl - 2) {
		send_to_char(ch, libc.CString("You failed to pick the lock...\r\n"))
		act(libc.CString("@c$n@w puts a set of lockpick tools away.@n"), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Move > ch.Move/30 {
			ch.Move -= ch.Move / 30
		} else {
			ch.Move = 0
		}
	} else {
		if ch.Move > ch.Move/30 {
			ch.Move -= ch.Move / 30
		} else {
			ch.Move = 0
		}
		return 1
	}
	return 0
}
func do_gen_door(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		door   int = -1
		keynum obj_vnum
		type_  [2048]byte
		dir    [2048]byte
		obj    *obj_data  = nil
		victim *char_data = nil
	)
	skip_spaces(&argument)
	if *argument == 0 {
		send_to_char(ch, libc.CString("%c%s what?\r\n"), unicode.ToUpper(rune(*cmd_door[subcmd])), (*byte)(unsafe.Add(unsafe.Pointer(cmd_door[subcmd]), 1)))
		return
	}
	two_arguments(argument, &type_[0], &dir[0])
	if generic_find(&type_[0], (1<<2)|1<<3, ch, &victim, &obj) == 0 {
		door = find_door(ch, &type_[0], &dir[0], cmd_door[subcmd])
	}
	if obj != nil && (int(obj.Type_flag) != ITEM_CONTAINER && int(obj.Type_flag) != ITEM_VEHICLE && int(obj.Type_flag) != ITEM_HATCH) {
		obj = nil
		door = find_door(ch, &type_[0], &dir[0], cmd_door[subcmd])
	}
	if obj != nil || door >= 0 {
		if obj != nil {
			keynum = obj_vnum(obj.Value[VAL_KEY_KEYCODE])
		} else {
			keynum = ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Key
		}
		if (func() int {
			if obj != nil {
				return obj.Value[VAL_DOOR_DCLOCK]
			}
			return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Dclock
		}()) == 0 {
			if obj != nil {
				obj.Value[VAL_DOOR_DCLOCK] = 20
			} else {
				((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Dclock = 20
			}
		}
		if !(func() bool {
			if obj != nil {
				return int(obj.Type_flag) == ITEM_CONTAINER && OBJVAL_FLAGGED(obj, 1<<0) || int(obj.Type_flag) == ITEM_VEHICLE && OBJVAL_FLAGGED(obj, 1<<0) || int(obj.Type_flag) == ITEM_HATCH && OBJVAL_FLAGGED(obj, 1<<0) || int(obj.Type_flag) == ITEM_WINDOW && OBJVAL_FLAGGED(obj, 1<<0) || int(obj.Type_flag) == ITEM_PORTAL && OBJVAL_FLAGGED(obj, 1<<0)
			}
			return EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<0)
		}()) {
			act(libc.CString("You can't $F that!"), FALSE, ch, nil, unsafe.Pointer(cmd_door[subcmd]), TO_CHAR)
		} else if !(func() bool {
			if obj != nil {
				return !OBJVAL_FLAGGED(obj, 1<<2)
			}
			return !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1)
		}()) && ((flags_door[subcmd])&(1<<0)) != 0 {
			send_to_char(ch, libc.CString("But it's already closed!\r\n"))
		} else if (func() bool {
			if obj != nil {
				return !OBJVAL_FLAGGED(obj, 1<<2)
			}
			return !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1)
		}()) && ((flags_door[subcmd])&(1<<1)) != 0 {
			send_to_char(ch, libc.CString("But it's currently open!\r\n"))
		} else if (func() bool {
			if obj != nil {
				return !OBJVAL_FLAGGED(obj, 1<<3)
			}
			return !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<2)
		}()) && ((flags_door[subcmd])&(1<<3)) != 0 {
			send_to_char(ch, libc.CString("Oh.. it wasn't locked, after all..\r\n"))
		} else if !(func() bool {
			if obj != nil {
				return !OBJVAL_FLAGGED(obj, 1<<3)
			}
			return !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<2)
		}()) && ((flags_door[subcmd])&(1<<2)) != 0 {
			send_to_char(ch, libc.CString("It seems to be locked.\r\n"))
		} else if has_key(ch, keynum) == 0 && !ADM_FLAGGED(ch, ADM_NOKEYS) && (subcmd == SCMD_LOCK || subcmd == SCMD_UNLOCK) {
			send_to_char(ch, libc.CString("You don't seem to have the proper key.\r\n"))
		} else if obj == nil && ok_pick(ch, keynum, int(libc.BoolToInt(func() bool {
			if obj != nil {
				return OBJVAL_FLAGGED(obj, 1<<1)
			}
			return EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<3)
		}())), func() int {
			if obj != nil {
				return obj.Value[VAL_DOOR_DCLOCK]
			}
			return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Dclock
		}(), subcmd, nil) != 0 {
			do_doorcmd(ch, obj, door, subcmd)
		} else if ok_pick(ch, keynum, int(libc.BoolToInt(func() bool {
			if obj != nil {
				return OBJVAL_FLAGGED(obj, 1<<1)
			}
			return EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<3)
		}())), func() int {
			if obj != nil {
				return obj.Value[VAL_DOOR_DCLOCK]
			}
			return ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Dclock
		}(), subcmd, obj) != 0 && obj != nil {
			do_doorcmd(ch, obj, door, subcmd)
		}
	}
	return
}
func do_simple_enter(ch *char_data, obj *obj_data, need_specials_check int) int {
	var (
		dest_room     room_rnum = real_room(room_vnum(obj.Value[VAL_PORTAL_DEST]))
		was_in        room_rnum = ch.In_room
		need_movement int       = 0
	)
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master != nil && ch.In_room == ch.Master.In_room {
		send_to_char(ch, libc.CString("The thought of leaving your master makes you weep.\r\n"))
		act(libc.CString("$n bursts into tears."), FALSE, ch, nil, nil, TO_ROOM)
		return 0
	}
	need_movement = 1
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity > 10 {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && int(ch.Chclass) != CLASS_BARDOCK && !IS_NPC(ch) {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	}
	if GET_LEVEL(ch) <= 1 {
		need_movement = 0
	}
	if ch.Move < int64(need_movement) && !AFF_FLAGGED(ch, AFF_FLYING) && !IS_NPC(ch) {
		if need_specials_check != 0 && ch.Master != nil {
			send_to_char(ch, libc.CString("You are too exhausted to follow.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You are too exhausted.\r\n"))
		}
		return 0
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_ATRIUM) {
		if House_can_enter(ch, func() room_vnum {
			if dest_room != room_rnum(-1) && dest_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(dest_room)))).Number
			}
			return -1
		}()) == 0 {
			send_to_char(ch, libc.CString("That's private property -- no trespassing!\r\n"))
			return 0
		}
	}
	if ROOM_FLAGGED(dest_room, ROOM_TUNNEL) && num_pc_in_room((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(dest_room)))) >= config_info.Play.Tunnel_size {
		if config_info.Play.Tunnel_size > 1 {
			send_to_char(ch, libc.CString("There isn't enough room for you to go there!\r\n"))
		} else {
			send_to_char(ch, libc.CString("There isn't enough room there for more than one person!\r\n"))
		}
		return 0
	}
	if ROOM_FLAGGED(dest_room, ROOM_GODROOM) && ch.Admlevel < ADMLVL_GRGOD {
		send_to_char(ch, libc.CString("You aren't godly enough to use that room!\r\n"))
		return 0
	}
	if !IS_NPC(ch) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) && !AFF_FLAGGED(ch, AFF_FLYING) {
		ch.Move -= int64(need_movement)
	}
	act(libc.CString("$n enters $p."), TRUE, ch, obj, nil, int(TO_ROOM|2<<9))
	if ch.Drag != nil {
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
	}
	char_from_room(ch)
	char_to_room(ch, dest_room)
	if entry_mtrigger(ch) == 0 {
		char_from_room(ch)
		char_to_room(ch, was_in)
		return 0
	}
	if int(obj.Type_flag) == ITEM_PORTAL {
		act(libc.CString("$n arrives from $p."), FALSE, ch, obj, nil, int(TO_ROOM|2<<9))
	} else {
		act(libc.CString("$n arrives from outside."), FALSE, ch, nil, nil, int(TO_ROOM|2<<9))
	}
	if ch.Drag != nil {
		act(libc.CString("@wYou drag @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		if !AFF_FLAGGED(ch.Drag, AFF_KNOCKED) && !AFF_FLAGGED(ch.Drag, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Drag, libc.CString("You feel your sleeping body being moved.\r\n"))
			if IS_NPC(ch.Drag) && ch.Drag.Fighting == nil {
				set_fighting(ch.Drag, ch)
			}
		}
		char_from_room(ch.Drag)
		char_to_room(ch.Drag, ch.In_room)
		if ch.Drag.Sits != nil {
			obj_from_room(ch.Drag.Sits)
			obj_to_room(ch.Drag.Sits, ch.In_room)
		}
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@wYou carry @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_CHAR)
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
		if !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_KNOCKED) && !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Player_specials.Carrying, libc.CString("You feel your sleeping body being moved.\r\n"))
		}
		char_from_room(ch.Player_specials.Carrying)
		char_to_room(ch.Player_specials.Carrying, ch.In_room)
		if ch.Player_specials.Carrying.Sits != nil {
			obj_from_room(ch.Player_specials.Carrying.Sits)
			obj_to_room(ch.Player_specials.Carrying.Sits, ch.In_room)
		}
	}
	if ch.Desc != nil {
		look_at_room(ch.In_room, ch, 0)
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_DEATH) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		log_death_trap(ch)
		death_cry(ch)
		extract_char(ch)
		return 0
	}
	entry_memory_mtrigger(ch)
	greet_memory_mtrigger(ch)
	return 1
}
func perform_enter_obj(ch *char_data, obj *obj_data, need_specials_check int) int {
	var (
		was_in     room_rnum = ch.In_room
		could_move int       = FALSE
		k          *follow_type
	)
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return 0
	}
	if int(obj.Type_flag) == ITEM_VEHICLE || int(obj.Type_flag) == ITEM_PORTAL {
		if OBJVAL_FLAGGED(obj, 1<<2) {
			send_to_char(ch, libc.CString("But it's closed!\r\n"))
		} else if (obj.Value[VAL_PORTAL_DEST]) != int(-1) && real_room(room_vnum(obj.Value[VAL_PORTAL_DEST])) != room_rnum(-1) {
			if (obj.Value[VAL_PORTAL_DEST]) >= 45000 && (obj.Value[VAL_PORTAL_DEST]) <= 0xB02B {
				var (
					tch    *char_data
					next_v *char_data
					filled int = FALSE
				)
				for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(room_vnum(obj.Value[VAL_PORTAL_DEST])))))).People; tch != nil; tch = next_v {
					next_v = tch.Next_in_room
					if tch != nil {
						filled = TRUE
					}
				}
				if filled == TRUE {
					send_to_char(ch, libc.CString("Only one person can fit in there at a time.\r\n"))
					return 0
				}
			}
			if (func() int {
				could_move = do_simple_enter(ch, obj, need_specials_check)
				return could_move
			}()) != 0 {
				for k = ch.Followers; k != nil; k = k.Next {
					if k.Follower.In_room == was_in && int(k.Follower.Position) >= POS_STANDING {
						act(libc.CString("You follow $N.\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
						perform_enter_obj(k.Follower, obj, 1)
					}
				}
			}
		} else {
			send_to_char(ch, libc.CString("It doesn't look like you can enter it at the moment.\r\n"))
		}
	} else {
		send_to_char(ch, libc.CString("You can't enter that!\r\n"))
	}
	return could_move
}
func do_enter(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		obj      *obj_data = nil
		buf      [2048]byte
		door     int
		move_dir int = -1
	)
	one_argument(argument, &buf[0])
	if buf[0] != 0 {
		obj = get_obj_in_list_vis(ch, &buf[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
		if obj == nil {
			obj = get_obj_in_list_vis(ch, &buf[0], nil, ch.Carrying)
		}
		if obj == nil {
			obj = get_obj_in_equip_vis(ch, &buf[0], nil, ch.Equipment[:])
		}
		if obj != nil {
			perform_enter_obj(ch, obj, 0)
		} else {
			for door = 0; door < NUM_OF_DIRS; door++ {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
					if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword != nil {
						if isname(&buf[0], ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).Keyword) != 0 {
							move_dir = door
						}
					}
				}
			}
			if move_dir > -1 {
				perform_move(ch, move_dir, 1)
			} else {
				send_to_char(ch, libc.CString("There is no %s here.\r\n"), &buf[0])
			}
		}
	} else if ROOM_FLAGGED(ch.In_room, ROOM_INDOORS) {
		send_to_char(ch, libc.CString("You are already indoors.\r\n"))
	} else {
		for door = 0; door < NUM_OF_DIRS; door++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room != room_rnum(-1) {
					if !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1) && ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ROOM_INDOORS) {
						move_dir = door
					}
				}
			}
		}
		if move_dir > -1 {
			perform_move(ch, move_dir, 1)
		} else {
			send_to_char(ch, libc.CString("You can't seem to find anything to enter.\r\n"))
		}
	}
}
func do_simple_leave(ch *char_data, obj *obj_data, need_specials_check int) int {
	var (
		was_in        room_rnum = ch.In_room
		dest_room     room_rnum = room_rnum(-1)
		need_movement int       = 0
		vehicle       *obj_data = nil
	)
	if int(obj.Type_flag) != ITEM_PORTAL {
		vehicle = find_vehicle_by_vnum(obj.Value[VAL_HATCH_DEST])
	}
	if vehicle == nil && int(obj.Type_flag) != ITEM_PORTAL {
		send_to_char(ch, libc.CString("That doesn't appear to lead anywhere.\r\n"))
		return 0
	}
	if int(obj.Type_flag) == ITEM_PORTAL && OBJVAL_FLAGGED(obj, 1<<2) {
		send_to_char(ch, libc.CString("But it's closed!\r\n"))
		return 0
	}
	if vehicle != nil {
		if (func() room_rnum {
			dest_room = vehicle.In_room
			return dest_room
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("That doesn't appear to lead anywhere.\r\n"))
			return 0
		}
	}
	if vehicle == nil {
		if (func() room_rnum {
			dest_room = real_room(room_vnum(obj.Value[VAL_PORTAL_DEST]))
			return dest_room
		}()) == room_rnum(-1) {
			send_to_char(ch, libc.CString("That doesn't appear to lead anywhere.\r\n"))
			return 0
		}
	}
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master != nil && ch.In_room == ch.Master.In_room {
		send_to_char(ch, libc.CString("The thought of leaving your master makes you weep.\r\n"))
		act(libc.CString("$n bursts into tears."), FALSE, ch, nil, nil, TO_ROOM)
		return 0
	}
	need_movement = 1
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity > 10 {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && int(ch.Chclass) != CLASS_BARDOCK && !IS_NPC(ch) {
		need_movement = (need_movement + (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) * (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
	}
	if GET_LEVEL(ch) <= 1 {
		need_movement = 0
	}
	if ch.Move < int64(need_movement) && !AFF_FLAGGED(ch, AFF_FLYING) && !IS_NPC(ch) {
		if need_specials_check != 0 && ch.Master != nil {
			send_to_char(ch, libc.CString("You are too exhausted to follow.\r\n"))
		} else {
			send_to_char(ch, libc.CString("You are too exhausted.\r\n"))
		}
		return 0
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_ATRIUM) {
		if House_can_enter(ch, func() room_vnum {
			if dest_room != room_rnum(-1) && dest_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(dest_room)))).Number
			}
			return -1
		}()) == 0 {
			send_to_char(ch, libc.CString("That's private property -- no trespassing!\r\n"))
			return 0
		}
	}
	if ROOM_FLAGGED(dest_room, ROOM_TUNNEL) && num_pc_in_room((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(dest_room)))) >= config_info.Play.Tunnel_size {
		if config_info.Play.Tunnel_size > 1 {
			send_to_char(ch, libc.CString("There isn't enough room for you to go there!\r\n"))
		} else {
			send_to_char(ch, libc.CString("There isn't enough room there for more than one person!\r\n"))
		}
		return 0
	}
	if !IS_NPC(ch) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) && !AFF_FLAGGED(ch, AFF_FLYING) {
		ch.Move -= int64(need_movement)
	}
	act(libc.CString("$n leaves $p."), TRUE, ch, vehicle, nil, int(TO_ROOM|2<<9))
	if ch.Drag != nil {
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
	}
	char_from_room(ch)
	char_to_room(ch, dest_room)
	if entry_mtrigger(ch) == 0 {
		char_from_room(ch)
		char_to_room(ch, was_in)
		return 0
	}
	if vehicle != nil {
		act(libc.CString("$n arrives from inside $p."), TRUE, ch, vehicle, nil, int(TO_ROOM|2<<9))
	} else {
		act(libc.CString("$n arrives from inside"), TRUE, ch, nil, nil, int(TO_ROOM|2<<9))
	}
	if ch.Drag != nil {
		act(libc.CString("@wYou drag @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@w drags @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		char_from_room(ch.Drag)
		char_to_room(ch.Drag, ch.In_room)
		if ch.Drag.Sits != nil {
			obj_from_room(ch.Drag.Sits)
			obj_to_room(ch.Drag.Sits, ch.In_room)
		}
		if !AFF_FLAGGED(ch.Drag, AFF_KNOCKED) && !AFF_FLAGGED(ch.Drag, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Drag, libc.CString("You feel your sleeping body being moved.\r\n"))
			if IS_NPC(ch.Drag) && ch.Drag.Fighting == nil {
				set_fighting(ch.Drag, ch)
			}
		}
	}
	if ch.Player_specials.Carrying != nil {
		act(libc.CString("@wYou carry @C$N@w with you.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_CHAR)
		act(libc.CString("@C$n@w carries @c$N@w with $m.@n"), TRUE, ch, nil, unsafe.Pointer(ch.Player_specials.Carrying), TO_ROOM)
		char_from_room(ch.Player_specials.Carrying)
		char_to_room(ch.Player_specials.Carrying, ch.In_room)
		if ch.Player_specials.Carrying.Sits != nil {
			obj_from_room(ch.Player_specials.Carrying.Sits)
			obj_to_room(ch.Player_specials.Carrying.Sits, ch.In_room)
		}
		if !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_KNOCKED) && !AFF_FLAGGED(ch.Player_specials.Carrying, AFF_SLEEP) && rand_number(1, 3) != 0 {
			send_to_char(ch.Player_specials.Carrying, libc.CString("You feel your sleeping body being moved.\r\n"))
		}
	}
	var buf3 [64936]byte
	send_to_sense(0, libc.CString("You sense someone "), ch)
	stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@Y %s\r\n@RSomeone has entered your scouter detection range.@n", add_commas(ch.Hit))
	send_to_scouter(&buf3[0], ch, 0, 0)
	if ch.Desc != nil {
		act(obj.Action_description, TRUE, ch, obj, nil, TO_CHAR)
		look_at_room(ch.In_room, ch, 0)
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_DEATH) && !ADM_FLAGGED(ch, ADM_WALKANYWHERE) {
		log_death_trap(ch)
		death_cry(ch)
		extract_char(ch)
		return 0
	}
	entry_memory_mtrigger(ch)
	greet_memory_mtrigger(ch)
	return 1
}
func perform_leave_obj(ch *char_data, obj *obj_data, need_specials_check int) int {
	var (
		was_in     room_rnum = ch.In_room
		could_move int       = FALSE
		k          *follow_type
	)
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return 0
	}
	if OBJVAL_FLAGGED(obj, 1<<2) {
		send_to_char(ch, libc.CString("But the way out is closed.\r\n"))
	} else {
		if (obj.Value[VAL_HATCH_DEST]) != int(-1) {
			if (func() int {
				could_move = do_simple_leave(ch, obj, need_specials_check)
				return could_move
			}()) != 0 {
				for k = ch.Followers; k != nil; k = k.Next {
					if k.Follower.In_room == was_in && int(k.Follower.Position) >= POS_STANDING {
						act(libc.CString("You follow $N.\r\n"), FALSE, k.Follower, nil, unsafe.Pointer(ch), TO_CHAR)
						perform_leave_obj(k.Follower, obj, 1)
					}
				}
			}
		}
	}
	return could_move
}
func do_leave(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		door int
		obj  *obj_data = nil
	)
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = obj.Next_content {
		if CAN_SEE_OBJ(ch, obj) {
			if int(obj.Type_flag) == ITEM_HATCH || int(obj.Type_flag) == ITEM_PORTAL {
				perform_leave_obj(ch, obj, 0)
				return
			}
		}
	}
	if OUTSIDE(ch) {
		send_to_char(ch, libc.CString("You are outside.. where do you want to go?\r\n"))
	} else {
		for door = 0; door < NUM_OF_DIRS; door++ {
			if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]) != nil {
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room != room_rnum(-1) {
					if !EXIT_FLAGGED((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door], 1<<1) && !ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door]).To_room, ROOM_INDOORS) {
						perform_move(ch, door, 1)
						return
					}
				}
			}
		}
		send_to_char(ch, libc.CString("I see no obvious exits to the outside.\r\n"))
	}
}
func handle_fall(ch *char_data) {
	var room int = -1
	for ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]) != nil && (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_FLYING {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]).To_room != room_rnum(-1) && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]).To_room <= top_of_world {
			room = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]).To_room)))).Number)
		} else {
			room = -1
		}
		char_from_room(ch)
		char_to_room(ch, real_room(room_vnum(room)))
		if ch.Player_specials.Carrying != nil {
			char_from_room(ch.Player_specials.Carrying)
			char_to_room(ch.Player_specials.Carrying, real_room(room_vnum(room)))
		}
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[5]) == nil || (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_FLYING {
			act(libc.CString("@r$n slams into the ground!@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Hit-ch.Max_hit/20 <= 0 {
				ch.Hit = 1
			} else {
				ch.Hit -= ch.Max_hit / 20
			}
			act(libc.CString("@rYou slam into the ground!@n"), TRUE, ch, nil, nil, TO_CHAR)
			look_at_room(ch.In_room, ch, 0)
		} else {
			act(libc.CString("@r$n pummets down toward the ground below!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM && ch.Player_specials.Carried_by == nil && int(ch.Race) != RACE_KANASSAN {
		if ch.Move >= int64(gear_weight(ch)) {
			act(libc.CString("@bYou swim in place.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@b swims in place.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Move -= int64(gear_weight(ch))
		} else {
			ch.Move -= int64(gear_weight(ch))
			if ch.Move < 0 {
				ch.Move = 0
			}
			act(libc.CString("@RYou are drowning!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@b gulps water as $e struggles to stay above the water line.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Hit-gear_pl(ch)/3 <= 0 {
				act(libc.CString("@rYou drown!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n@r drowns!@n"), TRUE, ch, nil, nil, TO_ROOM)
				die(ch, nil)
				ch.Hit = 1
			} else {
				ch.Hit -= gear_pl(ch) / 3
			}
		}
	}
}
func check_swim(ch *char_data) int {
	if ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
		if ch.Mana >= (ch.Max_mana/1000)+int64(gear_weight(ch)/2) {
			ch.Mana -= (ch.Max_mana / 1000) + int64(gear_weight(ch)/2)
			return 1
		} else {
			ch.Mana = 0
			send_to_char(ch, libc.CString("You do not have enough ki to fly through space. You are drifting helplessly.\r\n"))
			return 0
		}
	} else {
		if ch.Move >= int64(gear_weight(ch)-1) {
			ch.Move -= int64(gear_weight(ch) - 1)
			return 1
		} else {
			send_to_char(ch, libc.CString("You are too tired to swim!\r\n"))
			return 0
		}
	}
}
func do_fly(ch *char_data, argument *byte, cmd int, subcmd int) {
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if ch.Absorbing != nil || ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You can't fly, you are struggling with someone right now!"))
		return
	}
	if ch.Grappling != nil || ch.Grappled != nil {
		send_to_char(ch, libc.CString("You can't fly, you are struggling with someone right now!"))
		return
	}
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_HEALT) {
			send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
			return
		}
		if PLR_FLAGGED(ch, PLR_PILOTING) {
			send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
			return
		}
	}
	if !IS_NPC(ch) && GET_SKILL(ch, SKILL_FOCUS) < 30 && int(ch.Race) != RACE_ANDROID {
		send_to_char(ch, libc.CString("You do not have enough focus to hold yourself aloft.\r\n"))
		send_to_char(ch, libc.CString("@wOOC@D: @WYou need the skill Focus at @m30@W.@n\r\n"))
		return
	}
	if arg[0] == 0 {
		if AFF_FLAGGED(ch, AFF_FLYING) && (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_FLYING && (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_SPACE {
			act(libc.CString("@WYou slowly settle down to the ground.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n slowly settles down to the ground.@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			ch.Altitude = 0
			return
		}
		if AFF_FLAGGED(ch, AFF_FLYING) && (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FLYING {
			act(libc.CString("@WYou begin to plummet to the ground!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n starts to pummet to the ground below!@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			ch.Altitude = 0
			handle_fall(ch)
			return
		}
		if AFF_FLAGGED(ch, AFF_FLYING) && (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_SPACE {
			act(libc.CString("@WYou let yourself drift aimlessly through space.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n starts to drift slowly.!@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			ch.Altitude = 0
			return
		}
		if ch.Mana < ch.Max_mana/100 && int(ch.Race) != RACE_ANDROID {
			send_to_char(ch, libc.CString("You do not have the ki to fly."))
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@WYou slowly take off into the sky.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n slowly takes off into the sky.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Sits != nil {
				ch.Sits.Sitting = nil
				ch.Sits = nil
			}
			if int(ch.Position) < POS_STANDING {
				ch.Position = POS_STANDING
			}
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			ch.Altitude = 1
			ch.Mana -= ch.Max_mana / 100
		}
	}
	if libc.StrCaseCmp(libc.CString("high"), &arg[0]) == 0 {
		if ch.Mana < ch.Max_mana/100 && int(ch.Race) != RACE_ANDROID {
			send_to_char(ch, libc.CString("You do not have the ki to fly."))
			return
		} else {
			reveal_hiding(ch, 0)
			act(libc.CString("@WYou rocket high into the sky.@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@W$n rockets high into the sky.@n"), TRUE, ch, nil, nil, TO_ROOM)
			if ch.Sits != nil {
				ch.Sits.Sitting = nil
				ch.Sits = nil
			}
			if int(ch.Position) < POS_STANDING {
				ch.Position = POS_STANDING
			}
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			ch.Altitude = 2
			ch.Mana -= ch.Max_mana / 100
		}
	}
	if libc.StrCaseCmp(libc.CString("space"), &arg[0]) == 0 {
		if !OUTSIDE(ch) {
			send_to_char(ch, libc.CString("You are not outside!"))
			return
		}
		if ch.Mana < ch.Max_mana/10 && int(ch.Race) != RACE_ANDROID {
			send_to_char(ch, libc.CString("You do not have the ki to fly to space."))
			return
		}
		if ch.Fighting != nil {
			send_to_char(ch, libc.CString("You are too busy fighting!"))
			return
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_EARTH) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(50))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_CERRIA) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(198))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_VEGETA) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(53))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_KANASSA) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) == 0x3A38 {
				reveal_hiding(ch, 0)
				ch.Altitude = 2
				ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
				if block_calc(ch) == 0 {
					return
				}
				ch.Altitude = 0
				ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
				var zone int = 0
				if (func() int {
					zone = int(real_zone_by_thing(func() room_vnum {
						if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
						}
						return -1
					}()))
					return zone
				}()) != int(-1) {
					fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
				}
				send_to_sense(1, libc.CString("leaving the planet"), ch)
				send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
				act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
				char_from_room(ch)
				char_to_room(ch, real_room(58))
				act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
				if int(ch.Race) != RACE_ANDROID {
					ch.Mana -= ch.Max_mana / 10
				}
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			} else {
				send_to_char(ch, libc.CString("You can only fly off the planet from the launchpad of Aquis.\r\n"))
			}
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_FRIGID) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(51))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_KONACK) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(52))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_NAMEK) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(54))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_AETHER) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(55))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_YARDRAT) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(56))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if ROOM_FLAGGED(ch.In_room, ROOM_ARLIA) {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(59))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else if (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 3400 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 3599 || (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) >= 62900 && (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) <= 0xF617 || (func() room_vnum {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
			}
			return -1
		}()) == 19600 {
			reveal_hiding(ch, 0)
			ch.Altitude = 2
			ch.Affected_by[int(AFF_FLYING/32)] |= 1 << (int(AFF_FLYING % 32))
			if block_calc(ch) == 0 {
				return
			}
			ch.Altitude = 0
			ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
			var zone int = 0
			if (func() int {
				zone = int(real_zone_by_thing(func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()))
				return zone
			}()) != int(-1) {
				fly_zone(zone_rnum(zone), libc.CString("can be seen blasting off into space!@n\r\n"), ch)
			}
			send_to_sense(1, libc.CString("leaving the planet"), ch)
			send_to_scouter(libc.CString("A powerlevel signal has left the planet"), ch, 0, 2)
			act(libc.CString("@CYou blast off from the ground and rocket through the air. Your speed increases until you manage to reach the brink of space!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n blasts off from the ground and rockets through the air. You quickly lose sight of $m as $e continues upward!@n"), TRUE, ch, nil, nil, TO_ROOM)
			char_from_room(ch)
			char_to_room(ch, real_room(57))
			act(libc.CString("@C$n blasts up from the atmosphere below and then comes to a stop.@n"), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@mOOC: Use the command 'land' to return to the planet from here.@n\r\n"))
			if int(ch.Race) != RACE_ANDROID {
				ch.Mana -= ch.Max_mana / 10
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		} else {
			send_to_char(ch, libc.CString("You are not on a planet.\r\n"))
			return
		}
	}
}
func do_stand(ch *char_data, argument *byte, cmd int, subcmd int) {
	var chair *obj_data
	if AFF_FLAGGED(ch, AFF_KNOCKED) {
		send_to_char(ch, libc.CString("You are knocked out cold for right now!\r\n"))
		return
	}
	if !IS_NPC(ch) && (ch.Limb_condition[2]) <= 0 && (ch.Limb_condition[3]) <= 0 {
		send_to_char(ch, libc.CString("With what legs will you be standing up on?\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	switch ch.Position {
	case POS_STANDING:
		send_to_char(ch, libc.CString("You are already standing.\r\n"))
	case POS_SITTING:
		reveal_hiding(ch, 0)
		send_to_char(ch, libc.CString("You stand up.\r\n"))
		act(libc.CString("$n clambers to $s feet."), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Sits != nil {
			if OBJWEAR_FLAGGED(ch.Sits, ITEM_WEAR_TAKE) && int(ch.Sits.Type_flag) != ITEM_CHAIR && ch.Carry_weight+int(ch.Sits.Weight) <= int(max_carry_weight(ch)) {
				obj_from_room(ch.Sits)
				obj_to_char(ch.Sits, ch)
				act(libc.CString("You pick up $p."), TRUE, ch, ch.Sits, nil, TO_CHAR)
				act(libc.CString("$n picks up $p."), TRUE, ch, ch.Sits, nil, TO_ROOM)
			}
			chair = ch.Sits
			chair.Sitting = nil
			ch.Sits = nil
		}
		if ch.Fighting != nil {
			ch.Position = POS_FIGHTING
		} else {
			ch.Position = POS_STANDING
		}
	case POS_RESTING:
		send_to_char(ch, libc.CString("You stop resting, and stand up.\r\n"))
		act(libc.CString("$n stops resting, and clambers to $s feet."), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Sits != nil {
			if OBJWEAR_FLAGGED(ch.Sits, ITEM_WEAR_TAKE) && ch.Carry_weight+int(ch.Sits.Weight) <= int(max_carry_weight(ch)) {
				obj_from_room(ch.Sits)
				obj_to_char(ch.Sits, ch)
				act(libc.CString("You pick up $p."), TRUE, ch, ch.Sits, nil, TO_CHAR)
				act(libc.CString("$n picks up $p."), TRUE, ch, ch.Sits, nil, TO_ROOM)
			}
			chair = ch.Sits
			chair.Sitting = nil
			ch.Sits = nil
		}
		ch.Position = POS_STANDING
	case POS_SLEEPING:
		send_to_char(ch, libc.CString("You have to wake up first!\r\n"))
	default:
		send_to_char(ch, libc.CString("You stop floating around, and put your feet on the ground.\r\n"))
		act(libc.CString("$n stops floating around, and puts $s feet on the ground."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Position = POS_STANDING
	}
}
func do_sit(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		chair *obj_data = nil
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if ch.Drag != nil {
		act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		ch.Drag.Dragged = nil
		ch.Drag = nil
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are busy carrying someone!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_FLYING) {
		do_fly(ch, nil, 0, 0)
	}
	if arg[0] == 0 {
		switch ch.Position {
		case POS_STANDING:
			reveal_hiding(ch, 0)
			send_to_char(ch, libc.CString("You sit down.\r\n"))
			act(libc.CString("$n sits down."), FALSE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SITTING
		case POS_SITTING:
			send_to_char(ch, libc.CString("You're sitting already.\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("You stop resting, and sit up.\r\n"))
			act(libc.CString("$n stops resting."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SITTING
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You have to wake up first.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Sit down while fighting? Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and sit down.\r\n"))
			act(libc.CString("$n stops floating around, and sits down."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SITTING
		}
	} else {
		if ch.Sits != nil {
			send_to_char(ch, libc.CString("You are already on something!\r\n"))
			return
		}
		if (func() *obj_data {
			chair = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return chair
		}()) == nil {
			send_to_char(ch, libc.CString("That isn't here.\r\n"))
			return
		}
		if GET_OBJ_VNUM(chair) == 65 {
			send_to_char(ch, libc.CString("You can't get on that!\r\n"))
			return
		}
		if chair.Sitting != nil {
			send_to_char(ch, libc.CString("Someone is already on that one!\r\n"))
			return
		}
		if int(chair.Type_flag) != ITEM_CHAIR && int(chair.Type_flag) != ITEM_BED {
			send_to_char(ch, libc.CString("You can't sit on that!\r\n"))
			return
		}
		if chair.Size+1 < get_size(ch) {
			send_to_char(ch, libc.CString("You are too large for it!\r\n"))
			return
		}
		switch ch.Position {
		case POS_STANDING:
			reveal_hiding(ch, 0)
			act(libc.CString("You sit down on $p."), FALSE, ch, chair, nil, TO_CHAR)
			act(libc.CString("$n sits down on $p."), FALSE, ch, chair, nil, TO_ROOM)
			ch.Position = POS_SITTING
			ch.Sits = chair
			chair.Sitting = ch
		case POS_SITTING:
			send_to_char(ch, libc.CString("You should stand up first.\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("You should stand up first.\r\n"))
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You have to wake up first.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Sit down while fighting? Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and sit down.\r\n"))
			act(libc.CString("$n stops floating around, and sits down."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SITTING
		}
	}
}
func do_rest(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		chair *obj_data = nil
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM {
		send_to_char(ch, libc.CString("You can't rest here!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SANCTUARY) {
		if GET_SKILL(ch, SKILL_BARRIER) != 0 {
			send_to_char(ch, libc.CString("You have a barrier around you and can't rest.\r\n"))
			return
		} else {
			ch.Barrier = 0
			ch.Affected_by[int(AFF_SANCTUARY/32)] &= ^(1 << (int(AFF_SANCTUARY % 32)))
		}
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are a bit busy at the moment!\r\n"))
		return
	}
	if ch.Kaioken > 0 {
		send_to_char(ch, libc.CString("You are utilizing kaioken and can't rest!\r\n"))
		return
	}
	if ch.Drag != nil {
		act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		ch.Drag.Dragged = nil
		ch.Drag = nil
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are carrying someone!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_FLYING) {
		do_fly(ch, nil, 0, 0)
	}
	if arg[0] == 0 {
		if ch.Sits != nil {
			chair = ch.Sits
			if int(chair.Type_flag) != ITEM_BED {
				send_to_char(ch, libc.CString("You can't lay on that!\r\n"))
				return
			}
		}
		switch ch.Position {
		case POS_STANDING:
			reveal_hiding(ch, 0)
			send_to_char(ch, libc.CString("You lay down and rest your tired bones.\r\n"))
			act(libc.CString("$n lays down and rests."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_RESTING
		case POS_SITTING:
			send_to_char(ch, libc.CString("You rest your tired bones.\r\n"))
			act(libc.CString("$n rests."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_RESTING
		case POS_RESTING:
			send_to_char(ch, libc.CString("You are already resting.\r\n"))
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You have to wake up first.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Rest while fighting?  Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and stop to rest your tired bones.\r\n"))
			act(libc.CString("$n stops floating around, and rests."), FALSE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_RESTING
		}
	} else {
		if ch.Sits != nil {
			send_to_char(ch, libc.CString("You are already on something!\r\n"))
			return
		}
		if (func() *obj_data {
			chair = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return chair
		}()) == nil {
			send_to_char(ch, libc.CString("That isn't here.\r\n"))
			return
		}
		if GET_OBJ_VNUM(chair) == 65 {
			send_to_char(ch, libc.CString("You can't get on that!\r\n"))
			return
		}
		if chair.Sitting != nil {
			send_to_char(ch, libc.CString("Someone is already on that one!\r\n"))
			return
		}
		if int(chair.Type_flag) != ITEM_BED {
			send_to_char(ch, libc.CString("You can't lay on that!\r\n"))
			return
		}
		if chair.Size+1 < get_size(ch) {
			send_to_char(ch, libc.CString("You are too large for it!\r\n"))
			return
		}
		switch ch.Position {
		case POS_STANDING:
			reveal_hiding(ch, 0)
			act(libc.CString("You lay down and rest on $p."), TRUE, ch, chair, nil, TO_CHAR)
			act(libc.CString("$n lays down and rests on $p."), TRUE, ch, chair, nil, TO_ROOM)
			ch.Sits = chair
			chair.Sitting = ch
			ch.Position = POS_RESTING
		case POS_SITTING:
			send_to_char(ch, libc.CString("You should get up first.\r\n"))
		case POS_RESTING:
			send_to_char(ch, libc.CString("You are already resting.\r\n"))
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You have to wake up first.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Rest while fighting?  Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and stop to rest your tired bones.\r\n"))
			act(libc.CString("$n stops floating around, and rests."), FALSE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_RESTING
		}
	}
}
func do_sleep(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		chair *obj_data = nil
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
		if PRF_FLAGGED(ch, PRF_ARENAWATCH) {
			ch.Player_specials.Pref[int(PRF_ARENAWATCH/32)] &= bitvector_t(int32(^(1 << (int(PRF_ARENAWATCH % 32)))))
			ch.Arenawatch = -1
			send_to_char(ch, libc.CString("You stop watching the arena action.\r\n"))
		}
	}
	if (ch.Bonuses[BONUS_INSOMNIAC]) != 0 {
		send_to_char(ch, libc.CString("You don't feel the least bit tired.\r\n"))
		return
	}
	if (func() int {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
		}
		return SECT_INSIDE
	}()) == SECT_WATER_NOSWIM {
		send_to_char(ch, libc.CString("You can't rest here!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("You are busy piloting a ship!\r\n"))
		return
	}
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You are a bit busy at the moment!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if PLR_FLAGGED(ch, PLR_POWERUP) {
		send_to_char(ch, libc.CString("You are busy powering up!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SANCTUARY) {
		if GET_SKILL(ch, SKILL_BARRIER) > 0 {
			send_to_char(ch, libc.CString("You have a barrier around you and can't sleep.\r\n"))
			return
		} else {
			ch.Barrier = 0
			ch.Affected_by[int(AFF_SANCTUARY/32)] &= ^(1 << (int(AFF_SANCTUARY % 32)))
		}
	}
	if ch.Kaioken > 0 {
		send_to_char(ch, libc.CString("You are utilizing kaioken and can't sleep!\r\n"))
		return
	}
	if ch.Sleeptime > 0 {
		send_to_char(ch, libc.CString("You aren't sleepy enough.\r\n"))
		return
	}
	if ch.Drag != nil {
		act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
		act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
		ch.Drag.Dragged = nil
		ch.Drag = nil
	}
	if ch.Player_specials.Carrying != nil {
		send_to_char(ch, libc.CString("You are carrying someone!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_FLYING) {
		do_fly(ch, nil, 0, 0)
	}
	if arg[0] == 0 {
		if ch.Sits != nil {
			chair = ch.Sits
			if int(chair.Type_flag) != ITEM_BED {
				send_to_char(ch, libc.CString("You can't sleep on %s.\r\n"), chair.Short_description)
				return
			}
		}
		switch ch.Position {
		case POS_STANDING:
			fallthrough
		case POS_SITTING:
			fallthrough
		case POS_RESTING:
			reveal_hiding(ch, 0)
			send_to_char(ch, libc.CString("You go to sleep.\r\n"))
			act(libc.CString("$n lies down and falls asleep."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SLEEPING
			if PLR_FLAGGED(ch, PLR_FURY) {
				send_to_char(ch, libc.CString("Your fury subsides for now. Next time try to take advantage of it before you calm down.\r\n"))
				ch.Act[int(PLR_FURY/32)] &= bitvector_t(int32(^(1 << (int(PLR_FURY % 32)))))
			}
			if int(ch.Stupidkiss) > 0 {
				ch.Stupidkiss = 0
				send_to_char(ch, libc.CString("You forget about that stupid kiss.\r\n"))
			}
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You are already sound asleep.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Sleep while fighting?  Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and lie down to sleep.\r\n"))
			act(libc.CString("$n stops floating around, and lie down to sleep."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SLEEPING
		}
	} else {
		if ch.Sits != nil {
			send_to_char(ch, libc.CString("You are already on something!\r\n"))
			return
		}
		if (func() *obj_data {
			chair = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return chair
		}()) == nil {
			send_to_char(ch, libc.CString("That isn't here.\r\n"))
			return
		}
		if GET_OBJ_VNUM(chair) == 65 {
			send_to_char(ch, libc.CString("You can't get on that!\r\n"))
			return
		}
		if chair.Sitting != nil {
			send_to_char(ch, libc.CString("Someone is already on that one!\r\n"))
			return
		}
		if int(chair.Type_flag) != ITEM_BED {
			send_to_char(ch, libc.CString("You can't sleep on that!\r\n"))
			return
		}
		if chair.Size+1 < get_size(ch) {
			send_to_char(ch, libc.CString("You are too large for it!\r\n"))
			return
		}
		switch ch.Position {
		case POS_RESTING:
			fallthrough
		case POS_SITTING:
			send_to_char(ch, libc.CString("You need to get up first!\r\n"))
		case POS_STANDING:
			reveal_hiding(ch, 0)
			act(libc.CString("You lay down on $p and sleep."), FALSE, ch, chair, nil, TO_CHAR)
			act(libc.CString("$n lays down on $p and sleeps."), FALSE, ch, chair, nil, TO_ROOM)
			if PLR_FLAGGED(ch, PLR_FURY) {
				send_to_char(ch, libc.CString("Your fury subsides for now. Next time try to take advantage of it before you calm down.\r\n"))
				ch.Act[int(PLR_FURY/32)] &= bitvector_t(int32(^(1 << (int(PLR_FURY % 32)))))
			}
			if int(ch.Stupidkiss) > 0 {
				ch.Stupidkiss = 0
				send_to_char(ch, libc.CString("You forget about that stupid kiss.\r\n"))
			}
			ch.Sits = chair
			chair.Sitting = ch
			ch.Position = POS_SLEEPING
		case POS_SLEEPING:
			send_to_char(ch, libc.CString("You are already sound asleep.\r\n"))
		case POS_FIGHTING:
			send_to_char(ch, libc.CString("Sleep while fighting?  Are you MAD?\r\n"))
		default:
			send_to_char(ch, libc.CString("You stop floating around, and lie down to sleep.\r\n"))
			act(libc.CString("$n stops floating around, and lie down to sleep."), TRUE, ch, nil, nil, TO_ROOM)
			ch.Position = POS_SLEEPING
		}
	}
}
func do_wake(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		vict *char_data
		self int = 0
	)
	one_argument(argument, &arg[0])
	if AFF_FLAGGED(ch, AFF_KNOCKED) {
		send_to_char(ch, libc.CString("You are knocked out cold for right now!\r\n"))
		return
	}
	if (ch.Bonuses[BONUS_LATE]) != 0 && int(ch.Position) == POS_SLEEPING {
		send_to_char(ch, libc.CString("Nah you're enjoying sleeping too much.\r\n"))
		return
	}
	if arg[0] != 0 {
		if int(ch.Position) == POS_SLEEPING {
			send_to_char(ch, libc.CString("Maybe you should wake yourself up first.\r\n"))
		} else if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
		} else if vict == ch {
			self = 1
		} else if AWAKE(vict) {
			act(libc.CString("$E is already awake."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		} else if AFF_FLAGGED(vict, AFF_SLEEP) {
			act(libc.CString("You can't wake $M up!"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		} else if int(vict.Position) < POS_SLEEPING {
			act(libc.CString("$E's in pretty bad shape!"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		} else if AFF_FLAGGED(vict, AFF_KNOCKED) {
			send_to_char(ch, libc.CString("They are knocked out cold for right now!\r\n"))
		} else if (ch.Bonuses[BONUS_LATE]) != 0 {
			send_to_char(ch, libc.CString("They say 'Yeah yeah...' and then roll back over.\r\n"))
		} else {
			act(libc.CString("You wake $M up."), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("You are awakened by $n."), FALSE, ch, nil, unsafe.Pointer(vict), int(TO_VICT|2<<7))
			vict.Position = POS_SITTING
			if vict.Dragged != nil {
				act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, vict.Dragged, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, vict.Dragged, nil, unsafe.Pointer(vict), TO_ROOM)
				vict.Dragged.Drag = nil
				vict.Dragged = nil
			}
			if vict.Player_specials.Carried_by != nil {
				if vict.Player_specials.Carried_by.Alignment > 50 {
					carry_drop(vict.Player_specials.Carried_by, 0)
				} else {
					carry_drop(vict.Player_specials.Carried_by, 1)
				}
			}
		}
		if self == 0 {
			return
		}
	}
	if AFF_FLAGGED(ch, AFF_SLEEP) {
		send_to_char(ch, libc.CString("You can't wake up!\r\n"))
	} else if int(ch.Position) > POS_SLEEPING {
		send_to_char(ch, libc.CString("You are already awake...\r\n"))
	} else {
		send_to_char(ch, libc.CString("You awaken, and sit up.\r\n"))
		act(libc.CString("$n awakens."), TRUE, ch, nil, nil, TO_ROOM)
		if ch.Dragged != nil {
			act(libc.CString("@WYou stop dragging @C$N@W!@n"), TRUE, ch.Dragged, nil, unsafe.Pointer(ch), TO_CHAR)
			act(libc.CString("@C$n@W stops dragging you!@n"), TRUE, ch.Dragged, nil, unsafe.Pointer(ch), TO_VICT)
			act(libc.CString("@C$n@W stops dragging @c$N@W!@n"), TRUE, ch.Dragged, nil, unsafe.Pointer(ch), TO_NOTVICT)
			ch.Dragged.Drag = nil
			ch.Dragged = nil
		}
		if ch.Player_specials.Carried_by != nil {
			if ch.Player_specials.Carried_by.Alignment > 50 {
				carry_drop(ch.Player_specials.Carried_by, 0)
			} else {
				carry_drop(ch.Player_specials.Carried_by, 1)
			}
		}
		ch.Position = POS_SITTING
	}
}
func do_follow(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		buf    [2048]byte
		leader *char_data
	)
	one_argument(argument, &buf[0])
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if buf[0] != 0 {
		if (func() *char_data {
			leader = get_char_vis(ch, &buf[0], nil, 1<<0)
			return leader
		}()) == nil {
			send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
			return
		}
	} else {
		send_to_char(ch, libc.CString("Whom do you wish to follow?\r\n"))
		return
	}
	if ch.Master == leader {
		act(libc.CString("You are already following $M."), FALSE, ch, nil, unsafe.Pointer(leader), TO_CHAR)
		return
	}
	if AFF_FLAGGED(ch, AFF_CHARM) && ch.Master != nil {
		act(libc.CString("But you only feel like following $N!"), FALSE, ch, nil, unsafe.Pointer(ch.Master), TO_CHAR)
	} else {
		if leader == ch {
			if ch.Master == nil {
				send_to_char(ch, libc.CString("You are already following yourself.\r\n"))
				return
			}
			stop_follower(ch)
		} else {
			if circle_follow(ch, leader) {
				send_to_char(ch, libc.CString("Sorry, but following in loops is not allowed.\r\n"))
				return
			}
			if ch.Master != nil {
				stop_follower(ch)
			}
			ch.Affected_by[int(AFF_GROUP/32)] &= ^(1 << (int(AFF_GROUP % 32)))
			reveal_hiding(ch, 0)
			add_follower(ch, leader)
		}
	}
}
