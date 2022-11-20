package main

import "math"

var seed uint

func circle_srandom(initial_seed uint) {
	seed = initial_seed
}
func circle_random() uint {
	var (
		lo   int
		hi   int
		test int
	)
	hi = int(seed / 0x1F31D)
	lo = int(seed % 0x1F31D)
	test = lo*0x41A7 - hi*2836
	if test > 0 {
		seed = uint(test)
	} else {
		seed = uint(test + math.MaxInt32)
	}
	return seed
}
