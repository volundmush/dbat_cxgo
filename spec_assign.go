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
		(*(*index_data)(unsafe.Add(unsafe.Pointer(mob_index), unsafe.Sizeof(index_data{})*uintptr(rnum)))).Func = fname
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
		(*(*index_data)(unsafe.Add(unsafe.Pointer(obj_index), unsafe.Sizeof(index_data{})*uintptr(rnum)))).Func = fname
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
		(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(rnum)))).Func = fname
	} else if mini_mud == 0 {
		basic_mud_log(libc.CString("SYSERR: Attempt to assign spec to non-existant room #%d"), room)
	}
}
func assign_mobiles() {
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, dziak)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, azimer)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(3010, postmaster)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, janitor)
	ASSIGNMOB(1, fido)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, lyrzaxyn)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_marduk)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, snake)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, magic_user)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, cleric_ao)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, guild_guard)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, cityguard)
	ASSIGNMOB(1, receptionist)
	ASSIGNMOB(1, guild_guard)
}
func assign_objects() {
	ASSIGNOBJ(3034, bank)
	ASSIGNOBJ(3036, bank)
	ASSIGNOBJ(11, gravity)
	ASSIGNOBJ(65, healtank)
	ASSIGNOBJ(3, augmenter)
}
func assign_rooms() {
	var i room_rnum
	ASSIGNROOM(5, dump)
	ASSIGNROOM(3, pet_shops)
	ASSIGNROOM(4, pet_shops)
	ASSIGNROOM(81, auction)
	ASSIGNROOM(82, auction)
	ASSIGNROOM(83, auction)
	ASSIGNROOM(84, auction)
	ASSIGNROOM(85, auction)
	ASSIGNROOM(86, auction)
	if config_info.Play.Dts_are_dumps != 0 {
		for i = 0; i <= top_of_world; i++ {
			if ROOM_FLAGGED(i, ROOM_DEATH) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i)))).Func = dump
			}
		}
	}
}
