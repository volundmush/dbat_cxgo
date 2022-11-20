package main

import "github.com/gotranspile/cxgo/runtime/libc"

const MAX_HOUSES = 1000
const MAX_GUESTS = 10
const HOUSE_PRIVATE = 0
const HOUSE_GOD = 1
const HOUSE_CLAN = 2
const HOUSE_UNOWNED = 3
const NUM_HOUSE_TYPES = 4
const HOUSE_NOGUESTS = 1
const HOUSE_FREE = 2
const HOUSE_NOIMMS = 4
const HOUSE_IMPONLY = 8
const HOUSE_RENTFREE = 16
const HOUSE_SAVENORENT = 32
const HOUSE_NOSAVE = 64
const HOUSE_NUM_FLAGS = 7

type house_control_rec struct {
	Vnum          int
	Atrium        int
	Exit_num      int16
	Built_on      libc.Time
	Mode          int
	Owner         int
	Num_of_guests int
	Guests        [10]int
	Last_payment  libc.Time
	Bitvector     int
	Builtby       int
	Spare2        int
	Spare3        int
	Spare4        int
	Spare5        int
	Spare6        int
	Spare7        int
}
