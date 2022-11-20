package main

import "unsafe"

const SW_ARRAY_MAX = 4
const MAX_TRADE = 5
const MAX_PROD = 5
const VERSION3_TAG = "v3.0"
const MAX_SHOP_OBJ = 100
const OBJECT_DEAD = 0
const OBJECT_NOTOK = 1
const OBJECT_OK = 2
const OBJECT_NOVAL = 3
const LIST_PRODUCE = 0
const LIST_TRADE = 1
const LIST_ROOM = 2
const TRADE_NOGOOD = 0
const TRADE_NOEVIL = 1
const TRADE_NONEUTRAL = 2
const TRADE_NOWIZARD = 3
const TRADE_NOCLERIC = 4
const TRADE_NOROGUE = 5
const TRADE_NOFIGHTER = 6
const TRADE_NOHUMAN = 7
const TRADE_NOICER = 8
const TRADE_NOSAIYAN = 9
const TRADE_NOKONATSU = 10
const TRADE_NONAMEK = 11
const TRADE_NOMUTANT = 12
const TRADE_NOKANASSAN = 13
const TRADE_NOBIO = 14
const TRADE_NOANDROID = 15
const TRADE_NODEMON = 16
const TRADE_NOMAJIN = 17
const TRADE_NOKAI = 18
const TRADE_NOTRUFFLE = 19
const TRADE_NOGOBLIN = 20
const TRADE_NOANIMAL = 21
const TRADE_NOORC = 22
const TRADE_NOSNAKE = 23
const TRADE_NOTROLL = 24
const TRADE_NOHALFBREED = 25
const TRADE_NOMINOTAUR = 26
const TRADE_NOKOBOLD = 27
const TRADE_NOLIZARDFOLK = 28
const TRADE_NOMONK = 29
const TRADE_NOPALADIN = 30
const TRADE_UNUSED = 31
const TRADE_ONLYWIZARD = 32
const TRADE_ONLYCLERIC = 33
const TRADE_ONLYROGUE = 34
const TRADE_ONLYFIGHTER = 35
const TRADE_ONLYMONK = 36
const TRADE_ONLYPALADIN = 37
const TRADE_NOSORCERER = 38
const TRADE_NODRUID = 39
const TRADE_NOBARD = 40
const TRADE_NORANGER = 41
const TRADE_NOBARBARIAN = 42
const TRADE_ONLYSORCERER = 43
const TRADE_ONLYDRUID = 44
const TRADE_ONLYBARD = 45
const TRADE_ONLYRANGER = 46
const TRADE_ONLYBARBARIAN = 47
const TRADE_ONLYARCANE_ARCHER = 48
const TRADE_ONLYARCANE_TRICKSTER = 49
const TRADE_ONLYARCHMAGE = 50
const TRADE_ONLYASSASSIN = 51
const TRADE_ONLYBLACKGUARD = 52
const TRADE_ONLYDRAGON_DISCIPLE = 53
const TRADE_ONLYDUELIST = 54
const TRADE_ONLYDWARVEN_DEFENDER = 55
const TRADE_ONLYELDRITCH_KNIGHT = 56
const TRADE_ONLYHIEROPHANT = 57
const TRADE_ONLYHORIZON_WALKER = 58
const TRADE_ONLYLOREMASTER = 59
const TRADE_ONLYMYSTIC_THEURGE = 60
const TRADE_ONLYSHADOWDANCER = 61
const TRADE_ONLYTHAUMATURGIST = 62
const TRADE_NOARCANE_ARCHER = 63
const TRADE_NOARCANE_TRICKSTER = 64
const TRADE_NOARCHMAGE = 65
const TRADE_NOASSASSIN = 66
const TRADE_NOBLACKGUARD = 67
const TRADE_NODRAGON_DISCIPLE = 68
const TRADE_NODUELIST = 69
const TRADE_NODWARVEN_DEFENDER = 70
const TRADE_NOELDRITCH_KNIGHT = 71
const TRADE_NOHIEROPHANT = 72
const TRADE_NOHORIZON_WALKER = 73
const TRADE_NOLOREMASTER = 74
const TRADE_NOMYSTIC_THEURGE = 75
const TRADE_NOSHADOWDANCER = 76
const TRADE_NOTHAUMATURGIST = 77
const TRADE_NOBROKEN = 78
const OPER_OPEN_PAREN = 0
const OPER_CLOSE_PAREN = 1
const OPER_OR = 2
const OPER_AND = 3
const OPER_NOT = 4
const MAX_OPER = 4
const WILL_START_FIGHT = 1
const WILL_BANK_MONEY = 2
const WILL_ALLOW_STEAL = 4
const MIN_OUTSIDE_BANK = 5000
const MAX_OUTSIDE_BANK = 15000
const MSG_NOT_OPEN_YET = "Come back later!"
const MSG_NOT_REOPEN_YET = "Sorry, we have closed, but come back later."
const MSG_CLOSED_FOR_DAY = "Sorry, come back tomorrow."
const MSG_NO_STEAL_HERE = "$n is a bloody thief!"
const MSG_NO_SEE_CHAR = "I don't trade with someone I can't see!"
const MSG_NO_SELL_ALIGN = "Get out of here before I call the guards!"
const MSG_NO_SELL_CLASS = "We don't serve your kind here!"
const MSG_NO_SELL_RACE = "Get lost! We don't serve you kind here!"
const MSG_NO_USED_WANDSTAFF = "I don't buy used up wands or staves!"
const MSG_CANT_KILL_KEEPER = "Get out of here before I call the guards!"
const MSG_NO_BUY_BROKEN = "Sorry, but I don't deal in broken items."

type shop_buy_data struct {
	Type     int
	Keywords *byte
}
type shop_data struct {
	Vnum          int
	Producing     []int
	Profit_buy    float32
	Profit_sell   float32
	Type          []shop_buy_data
	No_such_item1 *byte
	No_such_item2 *byte
	Missing_cash1 *byte
	Missing_cash2 *byte
	Do_not_buy    *byte
	Message_buy   *byte
	Message_sell  *byte
	Temper1       int
	Bitvector     uint32
	Keeper        int
	With_who      [4]uint32
	In_room       []int
	Open1         int
	Open2         int
	Close1        int
	Close2        int
	BankAccount   int
	Lastsort      int
	Func          func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool
}
type stack_data struct {
	Data [100]int
	Len  int
}
