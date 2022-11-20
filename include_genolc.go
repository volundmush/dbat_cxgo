package main

const STRING_TERMINATOR = 126
const CONFIG_GENOLC_MOBPROG = 0
const SL_MOB = 0
const SL_OBJ = 1
const SL_SHP = 2
const SL_WLD = 3
const SL_ZON = 4
const SL_CFG = 5
const SL_GLD = 6
const SL_MAX = 6
const SL_ACT = 7
const SL_HLP = 8

type save_list_data struct {
	Zone int
	Type int
	Next *save_list_data
}
