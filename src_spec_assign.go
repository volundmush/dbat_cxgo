package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func ASSIGNMOB(mob int, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool) {
	var rnum int
	if (func() int {
		rnum = real_mobile(mob)
		return rnum
	}()) != int(-1) {
		mob_index[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant mob #%d"), mob)
	}
}
func ASSIGNOBJ(obj int, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool) {
	var rnum int
	if (func() int {
		rnum = real_object(obj)
		return rnum
	}()) != int(-1) {
		obj_index[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant obj #%d"), obj)
	}
}
func ASSIGNROOM(room int, fname func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool) {
	var rnum int
	if (func() int {
		rnum = real_room(room)
		return rnum
	}()) != int(-1) {
		world[rnum].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return fname(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant room #%d"), room)
	}
}
func assign_mobiles() {
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return dziak(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return azimer(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(3010, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return postmaster(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return janitor(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return fido(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return lyrzaxyn(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_marduk(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return snake(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return magic_user(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cleric_ao(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return cityguard(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return receptionist(ch, me, cmd, argument)
	})
	ASSIGNMOB(1, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return guild_guard(ch, me, cmd, argument)
	})
}
func assign_objects() {
	ASSIGNOBJ(3034, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return bank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(3036, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return bank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(11, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return gravity(ch, me, cmd, argument)
	})
	ASSIGNOBJ(65, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return healtank(ch, me, cmd, argument)
	})
	ASSIGNOBJ(3, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return augmenter(ch, me, cmd, argument)
	})
}
func assign_rooms() {
	var i int
	ASSIGNROOM(5, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return dump(ch, me, cmd, argument)
	})
	ASSIGNROOM(3, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return pet_shops(ch, me, cmd, argument)
	})
	ASSIGNROOM(4, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return pet_shops(ch, me, cmd, argument)
	})
	ASSIGNROOM(81, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(82, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(83, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(84, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(85, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	ASSIGNROOM(86, func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
		return auction(ch, me, cmd, argument)
	})
	if config_info.Play.Dts_are_dumps != 0 {
		for i = 0; i <= top_of_world; i++ {
			if ROOM_FLAGGED(i, ROOM_DEATH) {
				world[i].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
					return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
						return dump(ch, me, cmd, argument)
					}(ch, me, cmd, argument)
				}
			}
		}
	}
}
