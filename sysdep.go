package main

import "unsafe"

const CIRCLE_GNU_LIBC_MEMORY_TRACK = 0
const CIRCLE_UNSIGNED_INDEX = 0
const NOTHING = -1
const NOWHERE = -1
const NOBODY = -1
const NOFLAG = -1
const I64T = "lld"
const SZT = "lld"
const TMT = "ld"
const FALSE = 0
const TRUE = 1
const NO = 0
const YES = 1

type vnum = int64
type room_vnum = vnum
type obj_vnum = vnum
type mob_vnum = vnum
type zone_vnum = vnum
type shop_vnum = vnum
type trig_vnum = vnum
type guild_vnum = vnum
type room_rnum = vnum
type obj_rnum = vnum
type mob_rnum = vnum
type zone_rnum = vnum
type shop_rnum = vnum
type trig_rnum = vnum
type guild_rnum = vnum
type bitvector_t = uint32
type CommandFunc func(ch *char_data, argument *byte, cmd int, subcmd int)
type SpecialFunc func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int
