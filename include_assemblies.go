package main

const ASSM_MAKE = 0
const ASSM_BAKE = 1
const ASSM_BREW = 2
const ASSM_ASSEMBLE = 3
const ASSM_CRAFT = 4
const ASSM_FLETCH = 5
const ASSM_KNIT = 6
const ASSM_MIX = 7
const ASSM_THATCH = 8
const ASSM_WEAVE = 9
const ASSM_FORGE = 10

type assembly_data struct {
	LVnum           int
	LNumComponents  int
	UchAssemblyType uint8
	PComponents     []component_data
}
type component_data struct {
	BExtract bool
	BInRoom  bool
	LVnum    int
}
