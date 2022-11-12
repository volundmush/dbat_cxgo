package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func disp_ship_locations(ch *char_data, vehicle *obj_data) {
	if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 50 {
		send_to_char(ch, libc.CString("@D------------------[ @GEarth@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Nexus City, South Ocean, Nexus field, Cherry Blossom Mountain,\n"))
		send_to_char(ch, libc.CString("Sandy Desert, Northern Plains, Korin's Tower, Kami's Lookout,\n"))
		send_to_char(ch, libc.CString("Shadow Forest, Decrepit Area, West City, Hercule Beach, Satan City.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 51 {
		send_to_char(ch, libc.CString("@D------------------[ @CFrigid@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Ice Crown City, Ice Highway, Topica Snowfield, Glug's Volcano,\n"))
		send_to_char(ch, libc.CString("Platonic Sea, Slave City, Acturian Woods, Desolate Demesne,\n"))
		send_to_char(ch, libc.CString("Chateau Ishran, Wyrm Spine Mountain, Cloud Ruler Temple, Koltoan mine.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
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
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 53 {
		send_to_char(ch, libc.CString("@D------------------[ @YVegeta@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Vegetos City, Blood Dunes, Ancestral Mountains, Destopa Swamp,\n"))
		send_to_char(ch, libc.CString("Pride Forest, Pride tower, Ruby Cave.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 198 {
		send_to_char(ch, libc.CString("@D------------------[ @MCerria@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Cerria Colony, Fistarl Volcano, Crystalline Forest.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 54 {
		send_to_char(ch, libc.CString("@D------------------[ @gNamek@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Senzu Village, Guru's House, Crystalline Cave, Elder Village,\n"))
		send_to_char(ch, libc.CString("Frieza's Ship, Kakureta Village.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 55 {
		send_to_char(ch, libc.CString("@D------------------[ @BAether@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Haven City, Serenity Lake, Kaiju Forest, Ortusian Temple,\n"))
		send_to_char(ch, libc.CString("Silent Glade.\n"))
		send_to_char(ch, libc.CString("@D--------------------------------------------@n\n"))
		send_to_char(ch, libc.CString("@D------------------[ @BAether@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Haven City, Serenity Lake, Kaiju Forest, Ortusian Temple,\n"))
		send_to_char(ch, libc.CString("Silent Glade.\n"))
		send_to_char(ch, libc.CString("@D--------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 56 {
		send_to_char(ch, libc.CString("@D-----------------[ @mYardrat@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Yardra City, Jade Forest, Jade Cliffs, Mount Valaria.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 57 {
		send_to_char(ch, libc.CString("@D-----------------[ @CZennith@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Utatlan City, Zenith Jungle, Ancient Castle.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 58 {
		send_to_char(ch, libc.CString("@D-----------------[ @CKanassa@D ]-----------------@c\n"))
		send_to_char(ch, libc.CString("Aquis City, Yunkai Pirate Base.\n"))
		send_to_char(ch, libc.CString("@D-------------------------------------------@n\n"))
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 59 {
		send_to_char(ch, libc.CString("@D------------------[ @MArlia@D ]------------------@c\n"))
		send_to_char(ch, libc.CString("Janacre, Arlian Wasteland, Arlia Mine, Kemabra Wastes.\n"))
		send_to_char(ch, libc.CString("@D---------------------------------------------@n\n"))
	} else {
		send_to_char(ch, libc.CString("You are not above a planet!\r\n"))
	}
}
func ship_land_location(ch *char_data, vehicle *obj_data, arg *byte) int {
	var landspot int = 50
	if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 50 {
		if C.strcasecmp(arg, libc.CString("Nexus City")) == 0 {
			landspot = 300
			landspot += rand_number(0, 63)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("South Ocean")) == 0 {
			landspot = 800
			landspot += rand_number(0, 99)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Nexus Field")) == 0 {
			landspot = 1150
			landspot += rand_number(-16, 28)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Cherry Blossom Mountain")) == 0 {
			landspot = 1180
			landspot += rand_number(0, 19)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Sandy Desert")) == 0 {
			landspot = 1287
			landspot += rand_number(0, 64)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Northern Plains")) == 0 {
			landspot = 1428
			landspot += rand_number(0, 55)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Korin's Tower")) == 0 {
			return 1456
		} else if C.strcasecmp(arg, libc.CString("Kami's Lookout")) == 0 {
			landspot = 1506
			landspot += rand_number(0, 30)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Shadow Forest")) == 0 {
			landspot = 1600
			landspot += rand_number(0, 66)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Decrepit Area")) == 0 {
			return 1710
		} else if C.strcasecmp(arg, libc.CString("West City")) == 0 {
			return 19510
		} else if C.strcasecmp(arg, libc.CString("Hercule Beach")) == 0 {
			landspot = 2141
			landspot += rand_number(0, 53)
			return landspot
		} else if C.strcasecmp(arg, libc.CString("Satan City")) == 0 {
			landspot = 1150
			landspot += rand_number(-16, 28)
			return landspot
			return 13020
		} else {
			send_to_char(ch, libc.CString("You don't know where that made up place is, but decided to land anyway."))
			return 300
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 51 {
		if C.strcasecmp(arg, libc.CString("Ice Crown City")) == 0 {
			return 4264
		} else if C.strcasecmp(arg, libc.CString("Ice Highway")) == 0 {
			return 4300
		} else if C.strcasecmp(arg, libc.CString("Topica Snowfield")) == 0 {
			return 4351
		} else if C.strcasecmp(arg, libc.CString("Glug's Volcano")) == 0 {
			return 4400
		} else if C.strcasecmp(arg, libc.CString("Platonic Sea")) == 0 {
			return 4600
		} else if C.strcasecmp(arg, libc.CString("Slave City")) == 0 {
			return 4800
		} else if C.strcasecmp(arg, libc.CString("Acturian Woods")) == 0 {
			return 5100
		} else if C.strcasecmp(arg, libc.CString("Desolate Demesne")) == 0 {
			return 5150
		} else if C.strcasecmp(arg, libc.CString("Chateau Ishran")) == 0 {
			return 5165
		} else if C.strcasecmp(arg, libc.CString("Wyrm Spine Mountain")) == 0 {
			return 5200
		} else if C.strcasecmp(arg, libc.CString("Cloud Ruler Temple")) == 0 {
			return 5500
		} else if C.strcasecmp(arg, libc.CString("Koltoan Mine")) == 0 {
			return 4944
		} else {
			send_to_char(ch, libc.CString("You don't know where that made up place is, but decided to land anyway."))
			return 4264
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 52 {
		if C.strcasecmp(arg, libc.CString("Tiranoc City")) == 0 {
			return 8006
		} else if C.strcasecmp(arg, libc.CString("Great Oroist Temple")) == 0 {
			return 8300
		} else if C.strcasecmp(arg, libc.CString("Elzthuan Forest")) == 0 {
			return 8400
		} else if C.strcasecmp(arg, libc.CString("Mazori Farm")) == 0 {
			return 8447
		} else if C.strcasecmp(arg, libc.CString("Dres")) == 0 {
			return 8500
		} else if C.strcasecmp(arg, libc.CString("Colvian Farm")) == 0 {
			return 8600
		} else if C.strcasecmp(arg, libc.CString("St Alucia")) == 0 {
			return 8700
		} else if C.strcasecmp(arg, libc.CString("Meridius Memorial")) == 0 {
			return 8800
		} else if C.strcasecmp(arg, libc.CString("Desert of Illusion")) == 0 {
			return 8900
		} else if C.strcasecmp(arg, libc.CString("Plains of Confusion")) == 0 {
			return 8954
		} else if C.strcasecmp(arg, libc.CString("Turlon Fair")) == 0 {
			return 9200
		} else if C.strcasecmp(arg, libc.CString("Wetlands")) == 0 {
			return 9700
		} else if C.strcasecmp(arg, libc.CString("Kerberos")) == 0 {
			return 9855
		} else if C.strcasecmp(arg, libc.CString("Shaeras Mansion")) == 0 {
			return 9864
		} else if C.strcasecmp(arg, libc.CString("Slavinus Ravine")) == 0 {
			return 9900
		} else if C.strcasecmp(arg, libc.CString("Furian Citadel")) == 0 {
			return 9949
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 8006
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 53 {
		if C.strcasecmp(arg, libc.CString("Vegetos City")) == 0 {
			return 2226
		} else if C.strcasecmp(arg, libc.CString("Blood Dunes")) == 0 {
			return 2600
		} else if C.strcasecmp(arg, libc.CString("Ancestral Mountains")) == 0 {
			return 2616
		} else if C.strcasecmp(arg, libc.CString("Destopa Swamp")) == 0 {
			return 2709
		} else if C.strcasecmp(arg, libc.CString("Pride forest")) == 0 {
			return 2800
		} else if C.strcasecmp(arg, libc.CString("Pride Tower")) == 0 {
			return 2899
		} else if C.strcasecmp(arg, libc.CString("Ruby Cave")) == 0 {
			return 2615
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 2226
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 54 {
		if C.strcasecmp(arg, libc.CString("Senzu Village")) == 0 {
			return 11600
		} else if C.strcasecmp(arg, libc.CString("Guru's House")) == 0 {
			return 0x27C6
		} else if C.strcasecmp(arg, libc.CString("Crystalline Cave")) == 0 {
			return 0x28EA
		} else if C.strcasecmp(arg, libc.CString("Elder Village")) == 0 {
			return 13300
		} else if C.strcasecmp(arg, libc.CString("Frieza's Ship")) == 0 {
			return 0x27DB
		} else if C.strcasecmp(arg, libc.CString("Kakureta Village")) == 0 {
			return 0x2AAA
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 11600
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 55 {
		if C.strcasecmp(arg, libc.CString("Haven City")) == 0 {
			return 12010
		} else if C.strcasecmp(arg, libc.CString("Serenity Lake")) == 0 {
			return 0x2F47
		} else if C.strcasecmp(arg, libc.CString("Kaiju Forest")) == 0 {
			return 12300
		} else if C.strcasecmp(arg, libc.CString("Ortusian Temple")) == 0 {
			return 12400
		} else if C.strcasecmp(arg, libc.CString("Silent Glade")) == 0 {
			return 0x30C0
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 12010
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 56 {
		if C.strcasecmp(arg, libc.CString("Yardra City")) == 0 {
			return 0x36B8
		} else if C.strcasecmp(arg, libc.CString("Jade Forest")) == 0 {
			return 14100
		} else if C.strcasecmp(arg, libc.CString("Jade Cliffs")) == 0 {
			return 14200
		} else if C.strcasecmp(arg, libc.CString("Mount Valaria")) == 0 {
			return 14300
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x36B8
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 198 {
		if C.strcasecmp(arg, libc.CString("Cerria Colony")) == 0 {
			return 0x447B
		} else if C.strcasecmp(arg, libc.CString("Crystalline Forest")) == 0 {
			return 7950
		} else if C.strcasecmp(arg, libc.CString("Fistarl Volcano")) == 0 {
			return 17420
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x447B
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 57 {
		if C.strcasecmp(arg, libc.CString("Utatlan City")) == 0 {
			return 3412
		} else if C.strcasecmp(arg, libc.CString("Zenith Jungle")) == 0 {
			return 3520
		} else if C.strcasecmp(arg, libc.CString("Ancient Castle")) == 0 {
			return 19600
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 3412
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 58 {
		if C.strcasecmp(arg, libc.CString("Aquis City")) == 0 {
			return 0x3A38
		} else if C.strcasecmp(arg, libc.CString("Yunkai Pirate Base")) == 0 {
			return 0x3D27
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x3A38
		}
	} else if (func() room_vnum {
		if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
		}
		return -1
	}()) == 59 {
		if C.strcasecmp(arg, libc.CString("Janacre")) == 0 {
			return 0x3E89
		} else if C.strcasecmp(arg, libc.CString("Arlian Wasteland")) == 0 {
			return 0x40A0
		} else if C.strcasecmp(arg, libc.CString("Arlia Mine")) == 0 {
			return 16600
		} else if C.strcasecmp(arg, libc.CString("Kemabra Wastes")) == 0 {
			return 0x41B0
		} else {
			send_to_char(ch, libc.CString("you don't know where that made up place is, but decided to land anyway."))
			return 0x3E89
		}
	} else {
		send_to_char(ch, libc.CString("You are not above a planet!\r\n"))
		return -1
	}
}
func find_vehicle_by_vnum(vnum int) *obj_data {
	var i *obj_data
	for i = object_list; i != nil; i = i.Next {
		if i.Type_flag == ITEM_VEHICLE {
			if GET_OBJ_VNUM(i) == obj_vnum(vnum) {
				return i
			}
		}
	}
	return nil
}
func find_hatch_by_vnum(vnum int) *obj_data {
	var i *obj_data
	for i = object_list; i != nil; i = i.Next {
		if i.Type_flag == ITEM_HATCH {
			if GET_OBJ_VNUM(i) == obj_vnum(vnum) {
				return i
			}
		}
	}
	return nil
}
func get_obj_in_list_type(type_ int, list *obj_data) *obj_data {
	var i *obj_data
	for i = list; i != nil; i = i.Next_content {
		if int(i.Type_flag) == type_ {
			return i
		}
	}
	return nil
}
func find_control(ch *char_data) *obj_data {
	var (
		controls *obj_data
		obj      *obj_data
		j        int
	)
	controls = get_obj_in_list_type(ITEM_CONTROL, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
	if controls == nil {
		for obj = ch.Carrying; obj != nil && controls == nil; obj = obj.Next_content {
			if CAN_SEE_OBJ(ch, obj) && obj.Type_flag == ITEM_CONTROL {
				controls = obj
			}
		}
	}
	if controls == nil {
		for j = 0; j < NUM_WEARS && controls == nil; j++ {
			if (ch.Equipment[j]) != nil && CAN_SEE_OBJ(ch, ch.Equipment[j]) && (ch.Equipment[j]).Type_flag == ITEM_CONTROL {
				controls = ch.Equipment[j]
			}
		}
	}
	return controls
}
func drive_into_vehicle(ch *char_data, vehicle *obj_data, arg *byte) {
	var (
		vehicle_in_out *obj_data
		was_in         int
	)
	_ = was_in
	var is_in int
	var is_going_to int
	var buf [2048]byte
	if *arg == 0 {
		send_to_char(ch, libc.CString("@wDrive into what?\r\n"))
	} else if (func() *obj_data {
		vehicle_in_out = get_obj_in_list_vis(ch, arg, nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Contents)
		return vehicle_in_out
	}()) == nil {
		send_to_char(ch, libc.CString("@wNothing here by that name!\r\n"))
	} else if vehicle_in_out.Type_flag != ITEM_VEHICLE {
		send_to_char(ch, libc.CString("@wThat's not a ship.\r\n"))
	} else if vehicle == vehicle_in_out {
		send_to_char(ch, libc.CString("@wMy, we are in a clever mood today, aren't we.\r\n"))
	} else {
		is_going_to = int(real_room(room_vnum(vehicle_in_out.Value[0])))
		if !IS_SET_AR((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(is_going_to)))).Room_flags[:], ROOM_VEHICLE) {
			send_to_char(ch, libc.CString("@wThat ship can't carry other ships."))
		} else {
			stdio.Sprintf(&buf[0], "%s @wenters %s.\n\r", vehicle.Short_description, vehicle_in_out.Short_description)
			send_to_room(vehicle.In_room, &buf[0])
			was_in = int(vehicle.In_room)
			obj_from_room(vehicle)
			obj_to_room(vehicle, room_rnum(is_going_to))
			is_in = int(vehicle.In_room)
			if ch.Desc != nil {
				act(libc.CString(""), TRUE, ch, nil, nil, TO_ROOM)
			}
			send_to_char(ch, libc.CString("@wThe ship flies onward:\r\n"))
			look_at_room(vehicle.In_room, ch, 0)
			stdio.Sprintf(&buf[0], "%s @wenters.\r\n", vehicle.Short_description)
			send_to_room(room_rnum(is_in), &buf[0])
		}
	}
}
func drive_outof_vehicle(ch *char_data, vehicle *obj_data) {
	var (
		hatch          *obj_data
		vehicle_in_out *obj_data
		buf            [2048]byte
	)
	if (func() *obj_data {
		hatch = get_obj_in_list_type(ITEM_HATCH, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Contents)
		return hatch
	}()) == nil {
		send_to_char(ch, libc.CString("@wNowhere to pilot out of.\r\n"))
	} else if (func() *obj_data {
		vehicle_in_out = find_vehicle_by_vnum(hatch.Value[0])
		return vehicle_in_out
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou can't pilot out anywhere!\r\n"))
	} else {
		stdio.Sprintf(&buf[0], "%s @wexits %s.\r\n", vehicle.Short_description, vehicle_in_out.Short_description)
		send_to_room(vehicle.In_room, &buf[0])
		obj_from_room(vehicle)
		obj_to_room(vehicle, vehicle_in_out.In_room)
		if ch.Desc != nil {
			act(libc.CString("@wThe @De@Wn@wg@Di@wn@We@Ds@w of the ship @rr@Ro@ra@Rr@w as it moves."), TRUE, ch, nil, nil, TO_ROOM)
		}
		send_to_char(ch, libc.CString("@wThe ship flies onward:\r\n"))
		look_at_room(vehicle.In_room, ch, 0)
		var door int
		for door = 0; door < NUM_OF_DIRS; door++ {
			if CAN_GO(ch, door) {
				send_to_room((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door].To_room, libc.CString("@wThe @De@Wn@wg@Di@wn@We@Ds@w of the ship @rr@Ro@ra@Rr@w as it moves.\r\n"))
			}
		}
		stdio.Sprintf(&buf[0], "%s @wflies out of %s.\r\n", vehicle.Short_description, vehicle_in_out.Short_description)
		send_to_room(vehicle.In_room, &buf[0])
	}
}
func drive_in_direction(ch *char_data, vehicle *obj_data, dir int) {
	var buf [2048]byte
	if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]) == nil || ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).To_room == room_rnum(-1) {
		send_to_char(ch, libc.CString("@wApparently %s doesn't exist there.\r\n"), dirs[dir])
	} else if (((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).Exit_info & (1 << 1)) != 0 {
		if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).Keyword != nil {
			send_to_char(ch, libc.CString("@wThe %s seems to be closed.\r\n"), fname(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).Keyword))
		} else {
			send_to_char(ch, libc.CString("@wIt seems to be closed.\r\n"))
		}
	} else if !IS_SET_AR((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).To_room)))).Room_flags[:], ROOM_VEHICLE) && !IS_SET_AR((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[dir]).To_room)))).Room_flags[:], ROOM_SPACE) {
		send_to_char(ch, libc.CString("@wThe ship can't fit there!\r\n"))
	} else {
		var (
			was_in int
			is_in  int
		)
		stdio.Sprintf(&buf[0], "%s @wflies %s.\n\r", vehicle.Short_description, dirs[dir])
		send_to_room(vehicle.In_room, &buf[0])
		was_in = int(vehicle.In_room)
		obj_from_room(vehicle)
		obj_to_room(vehicle, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(was_in)))).Dir_option[dir].To_room)
		var controls *obj_data
		if (func() *obj_data {
			controls = find_control(ch)
			return controls
		}()) != nil {
			if (controls.Value[3]) < 5 {
				controls.Value[3] += 1
			} else {
				controls.Value[3] = 0
				controls.Value[2] -= 1
				if (controls.Value[2]) < 0 {
					controls.Value[2] = 0
				}
			}
		}
		var hatch *obj_data = nil
		for hatch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(real_room(room_vnum(vehicle.Value[0])))))).Contents; hatch != nil; hatch = hatch.Next_content {
			if hatch.Type_flag == ITEM_HATCH {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					hatch.Value[3] = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number)
				} else {
					hatch.Value[3] = -1
				}
			}
		}
		is_in = int(vehicle.In_room)
		if ch.Desc != nil {
			act(libc.CString("@wThe @De@Wn@wg@Di@wn@We@Ds@w of the ship @rr@Ro@ra@Rr@w as it moves."), TRUE, ch, nil, nil, TO_ROOM)
		}
		send_to_char(ch, libc.CString("@wThe ship flies onward:\r\n"))
		look_at_room(room_rnum(is_in), ch, 0)
		if controls != nil {
			send_to_char(ch, libc.CString("@RFUEL@D: %s%s@n\r\n"), func() string {
				if (controls.Value[2]) >= 200 {
					return "@G"
				}
				if (controls.Value[2]) >= 100 {
					return "@Y"
				}
				return "@r"
			}(), add_commas(int64(controls.Value[2])))
		}
		var door int
		for door = 0; door < NUM_OF_DIRS; door++ {
			if CAN_GO(ch, door) {
				send_to_room((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[door].To_room, libc.CString("@wThe @De@Wn@wg@Di@wn@We@Ds@w of the ship @rr@Ro@ra@Rr@w as it moves.\r\n"))
			}
		}
		stdio.Sprintf(&buf[0], "%s @wflies in from the %s.\r\n", vehicle.Short_description, dirs[rev_dir[dir]])
		send_to_room(room_rnum(is_in), &buf[0])
	}
}
func do_warp(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vehicle  *obj_data
		controls *obj_data
		arg      [2048]byte
	)
	one_argument(argument, &arg[0])
	if IS_NPC(ch) {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if !PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("@wYou need to be seated in the pilot's seat.\r\n[Enter: Pilot ready/unready]\r\n"))
		return
	}
	if (func() *obj_data {
		controls = find_control(ch)
		return controls
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou have nothing to control here!\r\n"))
		return
	}
	if (func() *obj_data {
		vehicle = find_vehicle_by_vnum(controls.Value[0])
		return vehicle
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou can't find anything to pilot.\r\n"))
		return
	} else if arg[0] == 0 {
		send_to_char(ch, libc.CString("Syntax: shipwarp [ earth | vegeta | namek | konack | aether | frigid | buoy1 | buoy2 | buoy3 ]\r\n"))
		return
	} else if C.strcasecmp(&arg[0], libc.CString("earth")) != 0 && C.strcasecmp(&arg[0], libc.CString("vegeta")) != 0 && C.strcasecmp(&arg[0], libc.CString("namek")) != 0 && C.strcasecmp(&arg[0], libc.CString("konack")) != 0 && C.strcasecmp(&arg[0], libc.CString("frigid")) != 0 && C.strcasecmp(&arg[0], libc.CString("aether")) != 0 && C.strcasecmp(&arg[0], libc.CString("buoy1")) != 0 && C.strcasecmp(&arg[0], libc.CString("buoy2")) != 0 && C.strcasecmp(&arg[0], libc.CString("buoy3")) != 0 {
		send_to_char(ch, libc.CString("Syntax: shipwarp [ earth | vegeta | namek | konack | aether | frigid | buoy1 | buoy2 | buoy3 ]\r\n"))
		return
	} else if ROOM_FLAGGED(room_rnum(libc.BoolToInt(vehicle.In_room == 0)), ROOM_SPACE) {
		send_to_char(ch, libc.CString("Your ship needs to be in space to utilize its Instant Travel Warp Accelerator.\r\n"))
		return
	} else if GET_OBJ_VNUM(vehicle) != 18400 {
		send_to_char(ch, libc.CString("Your ship is not outfitted with an Instant Travel Warp Accelerator.\r\n"))
		return
	} else {
		if C.strcasecmp(&arg[0], libc.CString("earth")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0xA013 || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 50 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0xA013))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("namek")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0xA780 || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 54 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0xA780))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("frigid")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0x78A9 || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 51 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0x78A9))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("konack")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0x69B9 || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 52 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0x69B9))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("vegeta")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0x7E6D || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 53 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0x7E6D))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("aether")) == 0 {
			if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 0xA3E7 || (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == 55 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(0xA3E7))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("buoy1")) == 0 {
			if ch.Radar1 <= 0 {
				send_to_char(ch, libc.CString("You have not launched that buoy!\r\n"))
				return
			} else if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == ch.Radar1 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(ch.Radar1))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("buoy2")) == 0 {
			if ch.Radar2 <= 0 {
				send_to_char(ch, libc.CString("You have not launched that buoy!\r\n"))
				return
			} else if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == ch.Radar2 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(ch.Radar2))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else if C.strcasecmp(&arg[0], libc.CString("buoy3")) == 0 {
			if ch.Radar3 <= 0 {
				send_to_char(ch, libc.CString("You have not launched that buoy!\r\n"))
				return
			} else if (func() room_vnum {
				if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
				}
				return -1
			}()) == ch.Radar3 {
				send_to_char(ch, libc.CString("Your ship is already there!\r\n"))
				return
			} else {
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find your ship in a new location!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@BA glow of blue light floods in through the window for an instant. You feel a strange shift as the light disappears and you find the ship in a new location!@n"), TRUE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("%s @Bbegins to glow bright blue before disappearing in a flash of light!@n\r\n"), vehicle.Short_description)
				obj_from_room(vehicle)
				obj_to_room(vehicle, real_room(ch.Radar3))
				send_to_room(vehicle.In_room, libc.CString("@BSuddenly in a flash of blue light @n%s @B appears instantly!@n\r\n"), vehicle.Short_description)
			}
		} else {
			basic_mud_log(libc.CString("ERROR: Ship Instant Warp Failure! Unknown argument!"))
			send_to_char(ch, libc.CString("ERROR\r\n"))
			return
		}
	}
}
func do_drive(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		dir       int
		confirmed int = FALSE
		count     int = 0
		vehicle   *obj_data
		controls  *obj_data
		arg3      [2048]byte
	)
	one_argument(argument, &arg3[0])
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no arms!\r\n"))
		return
	}
	if C.strcasecmp(&arg3[0], libc.CString("unready")) == 0 && !IS_NPC(ch) {
		if !PLR_FLAGGED(ch, PLR_PILOTING) {
			send_to_char(ch, libc.CString("You are already not flying the ship!\r\n"))
			return
		} else if PLR_FLAGGED(ch, PLR_PILOTING) {
			act(libc.CString("@w$n stands up and stops piloting the ship."), TRUE, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("@wYou stand up from the pilot's seat.\r\n"))
			ch.Position = POS_STANDING
			ch.Act[int(PLR_PILOTING/32)] &= bitvector_t(^(1 << (int(PLR_PILOTING % 32))))
			return
		}
	}
	if C.strcasecmp(&arg3[0], libc.CString("ready")) == 0 && !IS_NPC(ch) {
		if (func() *obj_data {
			controls = find_control(ch)
			return controls
		}()) == nil {
			send_to_char(ch, libc.CString("@wYou have nothing to control here!\r\n"))
			return
		}
		var d *descriptor_data
		if PLR_FLAGGED(ch, PLR_PILOTING) {
			send_to_char(ch, libc.CString("@wYou are already piloting the ship, try [pilot unready].\r\n"))
			return
		}
		if ch.Player_specials.Carrying != nil {
			send_to_char(ch, libc.CString("@wYou are busy carrying someone.\r\n"))
			return
		}
		if ch.Drag != nil {
			send_to_char(ch, libc.CString("@wYou are busy dragging someone.\r\n"))
			return
		}
		for d = descriptor_list; d != nil; d = d.Next {
			if !IS_PLAYING(d) {
				continue
			}
			if d.Character == ch {
				continue
			}
			if PLR_FLAGGED(d.Character, PLR_PILOTING) && d.Character.In_room == ch.In_room {
				send_to_char(ch, libc.CString("@w%s is already piloting the ship!\r\n"), GET_NAME(d.Character))
				count = 1
				return
			}
		}
		if count == 0 {
			confirmed = TRUE
		}
	}
	if confirmed == TRUE {
		ch.Act[int(PLR_PILOTING/32)] |= bitvector_t(1 << (int(PLR_PILOTING % 32)))
		act(libc.CString("@w$n sits down and begins piloting the ship."), TRUE, ch, nil, nil, TO_ROOM)
		ch.Position = POS_SITTING
		send_to_char(ch, libc.CString("@wYou take a seat in the pilot's chair.\r\n"))
		return
	} else if !PLR_FLAGGED(ch, PLR_PILOTING) {
		send_to_char(ch, libc.CString("@wYou need to be seated in the pilot's seat.\r\n[Enter: Pilot ready/unready]\r\n"))
	} else if ch.Position < POS_SLEEPING {
		send_to_char(ch, libc.CString("@wYou can't see anything but stars!\r\n"))
	} else if AFF_FLAGGED(ch, AFF_BLIND) {
		send_to_char(ch, libc.CString("@wYou can't see a damned thing, you're blind!\r\n"))
	} else if room_is_dark(ch.In_room) != 0 && !CAN_SEE_IN_DARK(ch) {
		send_to_char(ch, libc.CString("@wIt is pitch black...\r\n"))
	} else if (func() *obj_data {
		controls = find_control(ch)
		return controls
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou have nothing to control here!\r\n"))
	} else if invalid_align(ch, controls) != 0 || invalid_class(ch, controls) != 0 || invalid_race(ch, controls) != 0 {
		act(libc.CString("@wYou are zapped by $p@w and instantly step away from it."), FALSE, ch, controls, nil, TO_CHAR)
		act(libc.CString("@w$n@w is zapped by $p@w and instantly steps away from it."), FALSE, ch, controls, nil, TO_ROOM)
	} else if (func() *obj_data {
		vehicle = find_vehicle_by_vnum(controls.Value[0])
		return vehicle
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou can't find anything to pilot.\r\n"))
	} else {
		var (
			arg  [2048]byte
			arg2 [2048]byte
		)
		half_chop(argument, &arg[0], &arg2[0])
		if (controls.Value[2]) <= 0 {
			send_to_char(ch, libc.CString("Your ship doesn't have enough fuel to move.\r\n"))
			return
		}
		if arg[0] == 0 {
			send_to_char(ch, libc.CString("@wPilot, yes, but where?\r\n"))
		} else if is_abbrev(&arg[0], libc.CString("into")) != 0 || is_abbrev(&arg[0], libc.CString("onto")) != 0 {
			drive_into_vehicle(ch, vehicle, &arg2[0])
		} else if is_abbrev(&arg[0], libc.CString("out")) != 0 && ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dir_option[OUTDIR]) == nil {
			drive_outof_vehicle(ch, vehicle)
		} else {
			if !OBJVAL_FLAGGED(vehicle, 1<<2) {
				send_to_char(ch, libc.CString("@wThe hatch is open, are you insane!?\r\n"))
				return
			}
			if C.strcasecmp(&arg[0], libc.CString("north")) == 0 || C.strcasecmp(&arg[0], libc.CString("n")) == 0 {
				dir = 0
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("east")) == 0 || C.strcasecmp(&arg[0], libc.CString("e")) == 0 {
				dir = 1
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("south")) == 0 || C.strcasecmp(&arg[0], libc.CString("s")) == 0 {
				dir = 2
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("west")) == 0 || C.strcasecmp(&arg[0], libc.CString("w")) == 0 {
				dir = 3
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("up")) == 0 || C.strcasecmp(&arg[0], libc.CString("u")) == 0 {
				dir = 4
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("down")) == 0 || C.strcasecmp(&arg[0], libc.CString("d")) == 0 {
				dir = 5
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("northwest")) == 0 || C.strcasecmp(&arg[0], libc.CString("nw")) == 0 || C.strcasecmp(&arg[0], libc.CString("northw")) == 0 {
				dir = 6
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("northeast")) == 0 || C.strcasecmp(&arg[0], libc.CString("ne")) == 0 || C.strcasecmp(&arg[0], libc.CString("northe")) == 0 {
				dir = 7
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("southeast")) == 0 || C.strcasecmp(&arg[0], libc.CString("se")) == 0 || C.strcasecmp(&arg[0], libc.CString("southe")) == 0 {
				dir = 8
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("southwest")) == 0 || C.strcasecmp(&arg[0], libc.CString("sw")) == 0 || C.strcasecmp(&arg[0], libc.CString("southw")) == 0 {
				dir = 9
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("inside")) == 0 {
				dir = 10
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("outside")) == 0 {
				dir = 11
				drive_in_direction(ch, vehicle, dir)
				if (controls.Value[1]) == 1 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*2)*1.5))
				} else if (controls.Value[1]) == 2 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
				} else if (controls.Value[1]) == 3 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*1.5))
				} else if (controls.Value[1]) == 4 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
				} else if (controls.Value[1]) == 5 {
					WAIT_STATE(ch, int(float64((int(1000000/OPT_USEC))*1)*0.5))
				}
			} else if C.strcasecmp(&arg[0], libc.CString("land")) == 0 {
				if arg2[0] == 0 {
					if GET_OBJ_VNUM(vehicle) >= 46000 && GET_OBJ_VNUM(vehicle) <= 0xB413 {
						send_to_char(ch, libc.CString("@wLand on which pad? 1, 2, 3 or 4?\r\n"))
						send_to_char(ch, libc.CString("@CSpecial Ship Ability@D: @wpilot land (area name)\n@GExample@D: @wpilot land Nexus City\r\n"))
						disp_ship_locations(ch, vehicle)
					} else {
						send_to_char(ch, libc.CString("@wLand on which pad? 1, 2, 3 or 4?\r\n"))
					}
					return
				}
				var blah [2048]byte
				_ = blah
				var land_location int = 50
				if GET_OBJ_VNUM(vehicle) > 0xB413 {
					if C.strcasecmp(&arg2[0], libc.CString("1")) != 0 && C.strcasecmp(&arg2[0], libc.CString("2")) != 0 && C.strcasecmp(&arg2[0], libc.CString("3")) != 0 && C.strcasecmp(&arg2[0], libc.CString("4")) != 0 {
						send_to_char(ch, libc.CString("@wLand on which pad? 1, 2, 3 or 4?\r\n"))
						return
					}
				} else if C.strcasecmp(&arg2[0], libc.CString("1")) != 0 && C.strcasecmp(&arg2[0], libc.CString("2")) != 0 && C.strcasecmp(&arg2[0], libc.CString("3")) != 0 && C.strcasecmp(&arg2[0], libc.CString("4")) != 0 {
					land_location = ship_land_location(ch, vehicle, &arg2[0])
				}
				var buf3 [2048]byte
				buf3[0] = '\x00'
				act(libc.CString("@wYou set the controls to descend.@n"), FALSE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @wmanipulates the ship controls.@n"), FALSE, ch, nil, nil, TO_ROOM)
				act(libc.CString("@RThe ship rocks and shakes as it descends through the atmosphere!@n"), FALSE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@RThe ship rocks and shakes as it descends through the atmosphere!@n"), FALSE, ch, nil, nil, TO_ROOM)
				if land_location <= 50 {
					act(libc.CString("@wThe ship has landed.@n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@wThe ship has landed.@n"), FALSE, ch, nil, nil, TO_ROOM)
				}
				if land_location > 50 {
					act(libc.CString("@wThe ship slams into the ground and forms a small crater!@n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@wThe ship slams into the ground and forms a small crater!@n"), FALSE, ch, nil, nil, TO_ROOM)
					obj_from_room(vehicle)
					obj_to_room(vehicle, real_room(room_vnum(land_location)))
				} else if vehicle.In_room == real_room(50) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(409))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(411))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(412))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(410))
					} else if C.strcasecmp(&arg2[0], libc.CString("4365")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x49D8))
					} else if C.strcasecmp(&arg2[0], libc.CString("6329")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x49ED))
					} else if C.strcasecmp(&arg2[0], libc.CString("1983")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x4A33))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(57) {
					obj_from_room(vehicle)
					obj_to_room(vehicle, real_room(3508))
				} else if vehicle.In_room == real_room(198) {
					obj_from_room(vehicle)
					obj_to_room(vehicle, real_room(17420))
				} else if vehicle.In_room == real_room(58) {
					obj_from_room(vehicle)
					obj_to_room(vehicle, real_room(0x3A38))
				} else if vehicle.In_room == real_room(53) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(2319))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(2318))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(2320))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(2322))
					} else if C.strcasecmp(&arg2[0], libc.CString("4126")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x4724))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(56) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x36B3))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x36B4))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x36B5))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x36B6))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(55) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2EE3))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2EE4))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2EE6))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2EE5))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(59) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x3EC1))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x3EC2))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x3EC3))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x3EC4))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(51) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(4264))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(4263))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(4261))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(4262))
					} else if C.strcasecmp(&arg2[0], libc.CString("1337")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x46C4))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(54) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2D6C))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2D6D))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(11630))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(0x2D6B))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else if vehicle.In_room == real_room(52) {
					if C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(8195))
					} else if C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(8196))
					} else if C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(8197))
					} else if C.strcasecmp(&arg2[0], libc.CString("4")) == 0 {
						obj_from_room(vehicle)
						obj_to_room(vehicle, real_room(8198))
					} else {
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_CHAR)
						act(libc.CString("@wLanding sequence aborted, improper coordinates.@n"), FALSE, ch, nil, nil, TO_ROOM)
					}
				} else {
					send_to_char(ch, libc.CString("@wYou are not where you can land, you need to be in a planet's low orbit.@n\r\n"))
				}
				if land_location <= 50 {
					stdio.Sprintf(&buf3[0], "%s @wcomes in from above and slowly settles on the launch-pad.@n\r\n", vehicle.Short_description)
					look_at_room(vehicle.In_room, ch, 0)
					send_to_room(vehicle.In_room, &buf3[0])
				} else {
					stdio.Sprintf(&buf3[0], "%s @wcomes in from above and slams into the ground!@n\r\n", vehicle.Short_description)
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dmg += 1
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dmg >= 10 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Dmg = 10
					}
					look_at_room(vehicle.In_room, ch, 0)
					send_to_room(vehicle.In_room, &buf3[0])
				}
			} else if C.strcasecmp(&arg[0], libc.CString("launch")) == 0 {
				var (
					lnum int = 0
					rnum int = 0
				)
				if ROOM_FLAGGED(vehicle.In_room, ROOM_EARTH) {
					lnum = 1
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_FRIGID) {
					lnum = 2
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_KONACK) {
					lnum = 3
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_VEGETA) {
					lnum = 4
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_NAMEK) {
					lnum = 5
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_AETHER) {
					lnum = 6
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_YARDRAT) {
					lnum = 7
				} else if (func() room_vnum {
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					}
					return -1
				}()) >= 3400 && (func() room_vnum {
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					}
					return -1
				}()) <= 3599 || (func() room_vnum {
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					}
					return -1
				}()) >= 62900 && (func() room_vnum {
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					}
					return -1
				}()) <= 0xF617 || (func() room_vnum {
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					}
					return -1
				}()) == 19600 {
					lnum = 8
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_CERRIA) {
					lnum = 11
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_KANASSA) {
					lnum = 9
				} else if ROOM_FLAGGED(vehicle.In_room, ROOM_ARLIA) {
					lnum = 10
				} else {
					send_to_char(ch, libc.CString("@wYou are not on a planet.@n\r\n"))
					return
				}
				act(libc.CString("@wYou set the controls to launch.@n"), FALSE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @wmanipulates the ship controls.@n"), FALSE, ch, nil, nil, TO_ROOM)
				act(libc.CString("@RThe ship shudders as it launches up into the sky!@n"), FALSE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@RThe ship shudders as it launches up into the sky!@n"), FALSE, ch, nil, nil, TO_ROOM)
				act(libc.CString("@wThe ship has reached low orbit.@n"), FALSE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@wThe ship has reached low orbit.@n"), FALSE, ch, nil, nil, TO_ROOM)
				send_to_room(vehicle.In_room, libc.CString("@R%s @Rshudders before blasting off into the sky!@n"), vehicle.Short_description)
				if lnum == 1 {
					rnum = int(real_room(50))
				}
				if lnum == 2 {
					rnum = int(real_room(51))
				}
				if lnum == 3 {
					rnum = int(real_room(52))
				}
				if lnum == 4 {
					rnum = int(real_room(53))
				}
				if lnum == 5 {
					rnum = int(real_room(54))
				}
				if lnum == 6 {
					rnum = int(real_room(55))
				}
				if lnum == 7 {
					rnum = int(real_room(56))
				}
				if lnum == 8 {
					rnum = int(real_room(57))
				}
				if lnum == 9 {
					rnum = int(real_room(58))
				}
				if lnum == 10 {
					rnum = int(real_room(59))
				}
				if lnum == 11 {
					rnum = int(real_room(198))
				}
				if (controls.Value[3]) < 5 {
					controls.Value[3] += 1
				} else {
					controls.Value[3] = 0
					controls.Value[2] -= 1
					if (controls.Value[2]) < 0 {
						controls.Value[2] = 0
					}
				}
				obj_from_room(vehicle)
				obj_to_room(vehicle, room_rnum(rnum))
				look_at_room(vehicle.In_room, ch, 0)
				send_to_char(ch, libc.CString("@RFUEL@D: %s%s@n\r\n"), func() string {
					if (controls.Value[2]) >= 200 {
						return "@G"
					}
					if (controls.Value[2]) >= 100 {
						return "@Y"
					}
					return "@r"
				}(), add_commas(int64(controls.Value[2])))
			} else if C.strcasecmp(&arg[0], libc.CString("mark")) == 0 {
				var rnum int = 0
				_ = rnum
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("@wWhich marker are you wanting to launch? 1, 2, or 3?\r\n"))
					return
				}
				if C.strcasecmp(&arg2[0], libc.CString("1")) != 0 && C.strcasecmp(&arg2[0], libc.CString("2")) != 0 && C.strcasecmp(&arg2[0], libc.CString("3")) != 0 {
					send_to_char(ch, libc.CString("@wWhich marker are you wanting to launch? 1, 2, or 3?\r\n"))
					return
				}
				if !ROOM_FLAGGED(vehicle.In_room, ROOM_SPACE) {
					send_to_char(ch, libc.CString("@wYou need to be in space to launch a marker buoy.\r\n"))
					return
				}
				rnum = int(vehicle.In_room)
				if ch.Radar1 > 0 && C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
					send_to_char(ch, libc.CString("@wYou need to 'deactivate' that marker.\r\n"))
					return
				} else if ch.Radar1 <= 0 && C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
					act(libc.CString("@wYou enter a unique code and launch a marker buoy.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						ch.Radar1 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					} else {
						ch.Radar1 = -1
					}
				}
				if ch.Radar2 > 0 && C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
					send_to_char(ch, libc.CString("@wYou need to 'deactivate' that marker.\r\n"))
					return
				} else if ch.Radar2 <= 0 && C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
					act(libc.CString("@wYou enter a unique code and launch a marker buoy.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						ch.Radar2 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					} else {
						ch.Radar2 = -1
					}
				}
				if ch.Radar3 > 0 && C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
					send_to_char(ch, libc.CString("@wYou need to 'deactivate' that marker.\r\n"))
					return
				} else if ch.Radar3 <= 0 && C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
					act(libc.CString("@wYou enter a unique code and launch a marker buoy.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					if vehicle.In_room != room_rnum(-1) && vehicle.In_room <= top_of_world {
						ch.Radar3 = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vehicle.In_room)))).Number
					} else {
						ch.Radar3 = -1
					}
				}
			} else if C.strcasecmp(&arg[0], libc.CString("deactivate")) == 0 {
				if arg2[0] == 0 {
					send_to_char(ch, libc.CString("@wWhich marker are you wanting to launch? 1, 2, or 3?\r\n"))
					return
				}
				if C.strcasecmp(&arg2[0], libc.CString("1")) != 0 && C.strcasecmp(&arg2[0], libc.CString("2")) != 0 && C.strcasecmp(&arg2[0], libc.CString("3")) != 0 {
					send_to_char(ch, libc.CString("@wWhich marker are you wanting to deactivate? 1, 2, or 3?\r\n"))
					return
				}
				if ch.Radar1 <= 0 && C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
					send_to_char(ch, libc.CString("@wYou haven't launched that buoy yet.\r\n"))
					return
				} else if ch.Radar1 > 0 && C.strcasecmp(&arg2[0], libc.CString("1")) == 0 {
					act(libc.CString("@wYou enter buoy one's code and command it to deactivate.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					ch.Radar1 = 0
				}
				if ch.Radar2 <= 0 && C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
					send_to_char(ch, libc.CString("@wYou haven't launched that buoy yet.\r\n"))
					return
				} else if ch.Radar2 > 0 && C.strcasecmp(&arg2[0], libc.CString("2")) == 0 {
					act(libc.CString("@wYou enter buoy two's code and command it to deactivate.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					ch.Radar2 = 0
				}
				if ch.Radar3 <= 0 && C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
					send_to_char(ch, libc.CString("@wYou haven't launched that buoy yet.\r\n"))
					return
				} else if ch.Radar3 > 0 && C.strcasecmp(&arg2[0], libc.CString("3")) == 0 {
					act(libc.CString("@wYou enter buoy three's code and command it to deactivate.@n\r\n"), FALSE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@C$n@w manipulates the ship controls.@n\r\n"), FALSE, ch, nil, nil, TO_ROOM)
					ch.Radar3 = 0
				}
			} else {
				send_to_char(ch, libc.CString("@wThats not a valid direction.\r\n"))
				send_to_char(ch, libc.CString("Try one of these.\r\n"))
				send_to_char(ch, libc.CString("[ north/n  | south/s  | east/e  |  west/w  ]\r\n"))
				send_to_char(ch, libc.CString("[ up/u | down/d | northeast/ne/northe | northwest/nw/northw]\r\n"))
				send_to_char(ch, libc.CString("[  southeast/se/southe  |  southwest/sw/southw]\r\n"))
				send_to_char(ch, libc.CString("[  into  |  onto  |  inside  |  outside  ]@n\r\n"))
				send_to_char(ch, libc.CString("[ land | launch ]@n\r\n"))
			}
		}
	}
}
func do_ship_fire(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vehicle  *obj_data
		controls *obj_data
		arg1     [2048]byte
		arg2     [2048]byte
	)
	two_arguments(argument, &arg1[0], &arg2[0])
	if (func() *obj_data {
		controls = find_control(ch)
		return controls
	}()) == nil {
		send_to_char(ch, libc.CString("@wYou must be near the comm station in the cockpit.\r\n"))
		return
	}
	if (func() *obj_data {
		vehicle = find_vehicle_by_vnum(controls.Value[0])
		return vehicle
	}()) == nil {
		send_to_char(ch, libc.CString("@wSomething cosmic is jamming your signal! Quick call Iovan to repair it!\r\n"))
		return
	}
	var obj *obj_data = nil
	var obj2 *obj_data = nil
	_ = obj2
	var next_obj *obj_data = nil
	var shot int = FALSE
	for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if shot == FALSE {
			if obj.Type_flag == ITEM_VEHICLE && obj != vehicle {
				if C.strcasecmp(&arg1[0], obj.Name) == 0 {
					obj2 = obj
					shot = TRUE
				}
			}
		}
	}
}
