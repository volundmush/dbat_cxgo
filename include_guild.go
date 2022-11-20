package main

const GW_ARRAY_MAX = 4

type guild_data struct {
	Vnum            int
	Skills          [1000]int
	Charge          float32
	No_such_skill   *byte
	Not_enough_gold *byte
	Minlvl          int
	Gm              int
	With_who        [4]uint32
	Open            int
	Close           int
	Func            SpecialFunc
	Feats           [252]int
}
