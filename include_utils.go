package main

const READ_SIZE = 256
const OFF = 0
const BRF = 1
const NRM = 2
const CMP = 3
const CRASH_FILE = 0
const ETEXT_FILE = 1
const ALIAS_FILE = 2
const SCRIPT_VARS_FILE = 3
const NEW_OBJ_FILES = 4
const PLR_FILE = 5
const PET_FILE = 6
const IMC_FILE = 7
const USER_FILE = 8
const INTRO_FILE = 9
const SENSE_FILE = 10
const CUSTOME_FILE = 11
const MAX_FILES = 12
const BFS_ERROR = -1
const BFS_ALREADY_THERE = -2
const BFS_TO_FAR = -3
const BFS_NO_PATH = -4
const SECS_PER_MUD_HOUR = 900
const SECS_PER_MUD_DAY = 21600
const SECS_PER_MUD_MONTH = 648000
const SECS_PER_MUD_YEAR = 7776000
const SECS_PER_REAL_MIN = 60
const SECS_PER_REAL_HOUR = 3600
const SECS_PER_REAL_DAY = 86400
const SECS_PER_REAL_YEAR = 31536000
const FALSE = 0
const TRUE = 1
const SEEK_SET = 0
const SEEK_CUR = 1
const SEEK_END = 2

type xap_dir struct {
	Total    int
	Current  int
	Namelist bool
}
