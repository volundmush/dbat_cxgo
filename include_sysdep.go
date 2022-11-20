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
const NO = 0
const YES = 1

type CommandFunc func(ch *char_data, argument *byte, cmd int, subcmd int)
type SpecialFunc func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool
