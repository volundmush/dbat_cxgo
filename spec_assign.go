package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func ASSIGNMOB(mob mob_vnum, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int) {
	var rnum mob_rnum
	if (func() mob_rnum {
		rnum = real_mobile(mob)
		return rnum
	}()) != mob_rnum(-1) {
		mob_index[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant mob #%d"), mob)
	}
}
func ASSIGNOBJ(obj obj_vnum, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int) {
	var rnum obj_rnum
	if (func() obj_rnum {
		rnum = real_object(obj)
		return rnum
	}()) != obj_rnum(-1) {
		obj_index[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant obj #%d"), obj)
	}
}
func ASSIGNROOM(room room_vnum, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int) {
	var rnum room_rnum
	if (func() room_rnum {
		rnum = real_room(room)
		return rnum
	}()) != room_rnum(-1) {
		world[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant room #%d"), room)
	}
}
func assign_mobiles() {
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return dziak(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return azimer(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(3010, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return postmaster(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return janitor(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return fido(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return lyrzaxyn(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return guild_guard(ch, me, cmd, argument)
	})
}
func assign_objects() {
	ASSIGNOBJ(3034, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return bank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(3036, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return bank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(11, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return gravity(ch, me, cmd, argument)
	})
	ASSIGNOBJ(65, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return healtank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(3, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return augmenter(ch, me, cmd, argument)
	})
}
func assign_rooms() {
	var i room_rnum
	ASSIGNROOM(5, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return dump(ch, me, cmd, argument)
	})
	ASSIGNROOM(3, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return pet_shops(ch, me, cmd, argument)
	})
	ASSIGNROOM(4, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return pet_shops(ch, me, cmd, argument)
	})
	ASSIGNROOM(81, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(82, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(83, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(84, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(85, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(86, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
		return auction(ch, me, cmd, argument)
	})
	if config_info.Play.Dts_are_dumps != 0 {
		for i = 0; i <= top_of_world; i++ {
			if ROOM_FLAGGED(i, ROOM_DEATH) {
				world[i].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
					return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int {
						return dump(ch, me, cmd, argument)
					}(ch, me, cmd, argument)
				}
			}
		}
	}
}
