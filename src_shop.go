package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unicode"
	"unsafe"
)

var shop_index []shop_data
var top_shop int = -1
var cmd_say int
var cmd_tell int
var cmd_emote int
var cmd_slap int
var cmd_puke int
var operator_str [5]*byte = [5]*byte{libc.CString("[({"), libc.CString("])}"), libc.CString("|+"), libc.CString("&*"), libc.CString("^'")}
var trade_letters [79]*byte = [79]*byte{0: libc.CString("No Good"), 1: libc.CString("No Evil"), 2: libc.CString("No Neutral"), 3: libc.CString("No Roshi"), 4: libc.CString("No Piccolo"), 5: libc.CString("No Krane"), 6: libc.CString("No Nail"), 7: libc.CString("No Human"), 8: libc.CString("No Icer"), 9: libc.CString("No Saiyan"), 10: libc.CString("No Konatsu"), 11: libc.CString("No Namek"), 12: libc.CString("No Mutant"), 13: libc.CString("No Kanassan"), 14: libc.CString("No Bio"), 15: libc.CString("No Android"), 16: libc.CString("No Demon"), 17: libc.CString("No Majin"), 18: libc.CString("No Kai"), 19: libc.CString("No Truffle"), 20: libc.CString("No Goblin"), 21: libc.CString("No Animal"), 22: libc.CString("No Orc"), 23: libc.CString("No Snake"), 24: libc.CString("No Halfbreed"), 25: libc.CString("No Minotaur"), 26: libc.CString("No Kobold"), 27: libc.CString("No Lizardfolk"), 28: libc.CString("No Bardock"), 29: libc.CString("No Ginyu"), 30: libc.CString("UNUSED"), 31: libc.CString("Must be Roshi"), 32: libc.CString("Must be Piccolo"), 33: libc.CString("Must be Krane"), 34: libc.CString("Must be Nail"), 35: libc.CString("Must be Bardock"), 36: libc.CString("Must be Ginyu"), 37: libc.CString("No Frieza"), 38: libc.CString("No Tapion"), 39: libc.CString("No Android 16"), 40: libc.CString("No Dabura"), 41: libc.CString("No Kabito"), 42: libc.CString("Must be Frieza"), 43: libc.CString("Must be Tapion"), 44: libc.CString("Must be Android 16"), 45: libc.CString("Must be Dabura"), 46: libc.CString("Must be Kabito"), 47: libc.CString("Must be Jinto"), 48: libc.CString("Must be Tsuna"), 49: libc.CString("Must be Kurzak"), 50: libc.CString("Must be Assassin"), 51: libc.CString("Must be Blackguard"), 52: libc.CString("Must be Dragon Disciple"), 53: libc.CString("Must be Duelist"), 54: libc.CString("Must be Dwarven Defender"), 55: libc.CString("Must be Eldritch Knight"), 56: libc.CString("Must be Hierophant"), 57: libc.CString("Must be Horizon Walker"), 58: libc.CString("Must be Loremaster"), 59: libc.CString("Must be Mystic Theurge"), 60: libc.CString("Must be Shadowdancer"), 61: libc.CString("Must be Thaumaturgist"), 62: libc.CString("No Jinto"), 63: libc.CString("No Tsuna"), 64: libc.CString("No Kurzak"), 65: libc.CString("No Assassin"), 66: libc.CString("No Blackguard"), 67: libc.CString("No Dragon Disciple"), 68: libc.CString("No Duelist"), 69: libc.CString("No Dwarven Defender"), 70: libc.CString("No Eldritch Knight"), 71: libc.CString("No Hierophant"), 72: libc.CString("No Horizon Walker"), 73: libc.CString("No Loremaster"), 74: libc.CString("No Mystic Theurge"), 75: libc.CString("No Shadowdancer"), 76: libc.CString("No Thaumaturgist"), 77: libc.CString("\n")}
var shop_bits [4]*byte = [4]*byte{libc.CString("WILL_FIGHT"), libc.CString("USES_BANK"), libc.CString("ALLOW_STEAL"), libc.CString("\n")}

func is_ok_char(keeper *char_data, ch *char_data, shop_nr int) bool {
	var buf [2048]byte
	if !CAN_SEE(keeper, ch) {
		var actbuf [2048]byte = func() [2048]byte {
			var t [2048]byte
			copy(t[:], []byte(MSG_NO_SEE_CHAR))
			return t
		}()
		do_say(keeper, &actbuf[0], cmd_say, 0)
		return false
	}
	if ADM_FLAGGED(ch, ADM_ALLSHOPS) {
		return true
	}
	if ch.Alignment > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOGOOD) || ch.Alignment < 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOEVIL) || ch.Alignment == 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NONEUTRAL) {
		stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_NO_SELL_ALIGN)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return false
	}
	if IS_NPC(ch) {
		return true
	}
	if int(ch.Chclass) == CLASS_ROSHI && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOWIZARD) || int(ch.Chclass) == CLASS_PICCOLO && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOCLERIC) || int(ch.Chclass) == CLASS_KRANE && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOROGUE) || int(ch.Chclass) == CLASS_BARDOCK && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOMONK) || int(ch.Chclass) == CLASS_GINYU && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOPALADIN) || int(ch.Chclass) == CLASS_NAIL && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOFIGHTER) || int(ch.Chclass) == CLASS_KABITO && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOBARBARIAN) || int(ch.Chclass) == CLASS_FRIEZA && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOSORCERER) || int(ch.Chclass) == CLASS_ANDSIX && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOBARD) || int(ch.Chclass) == CLASS_DABURA && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NORANGER) || int(ch.Chclass) == CLASS_TAPION && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NODRUID) || int(ch.Chclass) == CLASS_NAIL && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOFIGHTER) || int(ch.Chclass) == CLASS_JINTO && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOARCANE_ARCHER) || int(ch.Chclass) == CLASS_TSUNA && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOARCANE_TRICKSTER) || int(ch.Chclass) == CLASS_KURZAK && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOARCHMAGE) || ((ch.Chclasses[CLASS_ASSASSIN])+(ch.Epicclasses[CLASS_ASSASSIN])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOASSASSIN) || ((ch.Chclasses[CLASS_BLACKGUARD])+(ch.Epicclasses[CLASS_BLACKGUARD])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOBLACKGUARD) || ((ch.Chclasses[CLASS_DRAGON_DISCIPLE])+(ch.Epicclasses[CLASS_DRAGON_DISCIPLE])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NODRAGON_DISCIPLE) || ((ch.Chclasses[CLASS_DUELIST])+(ch.Epicclasses[CLASS_DUELIST])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NODUELIST) || ((ch.Chclasses[CLASS_DWARVEN_DEFENDER])+(ch.Epicclasses[CLASS_DWARVEN_DEFENDER])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NODWARVEN_DEFENDER) || ((ch.Chclasses[CLASS_ELDRITCH_KNIGHT])+(ch.Epicclasses[CLASS_ELDRITCH_KNIGHT])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOELDRITCH_KNIGHT) || ((ch.Chclasses[CLASS_HIEROPHANT])+(ch.Epicclasses[CLASS_HIEROPHANT])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOHIEROPHANT) || ((ch.Chclasses[CLASS_HORIZON_WALKER])+(ch.Epicclasses[CLASS_HORIZON_WALKER])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOHORIZON_WALKER) || ((ch.Chclasses[CLASS_LOREMASTER])+(ch.Epicclasses[CLASS_LOREMASTER])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOLOREMASTER) || ((ch.Chclasses[CLASS_MYSTIC_THEURGE])+(ch.Epicclasses[CLASS_MYSTIC_THEURGE])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOMYSTIC_THEURGE) || ((ch.Chclasses[CLASS_SHADOWDANCER])+(ch.Epicclasses[CLASS_SHADOWDANCER])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOSHADOWDANCER) || ((ch.Chclasses[CLASS_THAUMATURGIST])+(ch.Epicclasses[CLASS_THAUMATURGIST])) > 0 && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOTHAUMATURGIST) {
		stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_NO_SELL_CLASS)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return false
	}
	if int(ch.Race) == RACE_HUMAN && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOHUMAN) || int(ch.Race) == RACE_ICER && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOICER) || int(ch.Race) == RACE_SAIYAN && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOSAIYAN) || int(ch.Race) == RACE_KONATSU && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOKONATSU) {
		stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_NO_SELL_RACE)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return false
	}
	return true
}
func is_ok_obj(keeper *char_data, ch *char_data, obj *obj_data, shop_nr int) bool {
	var buf [2048]byte
	if OBJ_FLAGGED(obj, ITEM_BROKEN) && IS_SET_AR(shop_index[shop_nr].With_who[:], TRADE_NOBROKEN) {
		stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_NO_BUY_BROKEN)
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return false
	}
	if OBJ_FLAGGED(obj, ITEM_FORGED) {
		stdio.Snprintf(&buf[0], int(2048), "%s that piece of junk is an obvious forgery!", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return false
	}
	return true
}
func is_open(keeper *char_data, shop_nr int, msg int) bool {
	var buf [2048]byte
	buf[0] = '\x00'
	if shop_index[shop_nr].Open1 > time_info.Hours {
		strlcpy(&buf[0], libc.CString(MSG_NOT_OPEN_YET), uint64(2048))
	} else if shop_index[shop_nr].Close1 < time_info.Hours {
		if shop_index[shop_nr].Open2 > time_info.Hours {
			strlcpy(&buf[0], libc.CString(MSG_NOT_REOPEN_YET), uint64(2048))
		} else if shop_index[shop_nr].Close2 < time_info.Hours {
			strlcpy(&buf[0], libc.CString(MSG_CLOSED_FOR_DAY), uint64(2048))
		}
	}
	if buf[0] == 0 {
		return true
	}
	if msg != 0 {
		do_say(keeper, &buf[0], cmd_tell, 0)
	}
	return false
}
func is_ok(keeper *char_data, ch *char_data, shop_nr int) bool {
	if is_open(keeper, shop_nr, 1) {
		return is_ok_char(keeper, ch, shop_nr)
	}
	return false
}
func push(stack *stack_data, pushval int) {
	stack.Data[func() int {
		p := &stack.Len
		x := *p
		*p++
		return x
	}()] = pushval
}
func top(stack *stack_data) int {
	if stack.Len > 0 {
		return stack.Data[stack.Len-1]
	}
	return -1
}
func pop(stack *stack_data) int {
	if stack.Len > 0 {
		return stack.Data[func() int {
			p := &stack.Len
			*p--
			return *p
		}()]
	} else {
		basic_mud_log(libc.CString("SYSERR: Illegal expression %d in shop keyword list."), stack.Len)
		return 0
	}
}
func evaluate_operation(ops *stack_data, vals *stack_data) {
	var oper int
	if (func() int {
		oper = pop(ops)
		return oper
	}()) == OPER_NOT {
		push(vals, int(libc.BoolToInt(pop(vals) == 0)))
	} else {
		var (
			val1 int = pop(vals)
			val2 int = pop(vals)
		)
		if oper == OPER_AND {
			push(vals, int(libc.BoolToInt(val1 != 0 && val2 != 0)))
		} else if oper == OPER_OR {
			push(vals, int(libc.BoolToInt(val1 != 0 || val2 != 0)))
		}
	}
}
func find_oper_num(token int8) int {
	var oindex int
	for oindex = 0; oindex <= MAX_OPER; oindex++ {
		if libc.StrChr(operator_str[oindex], byte(token)) != nil {
			return oindex
		}
	}
	return -1
}
func evaluate_expression(obj *obj_data, expr *byte) int {
	var (
		ops    stack_data
		vals   stack_data
		ptr    *byte
		end    *byte
		name   [64936]byte
		temp   int
		eindex int
	)
	if expr == nil || *expr == 0 {
		return 1
	}
	ops.Len = func() int {
		p := &vals.Len
		vals.Len = 0
		return *p
	}()
	ptr = expr
	for *ptr != 0 {
		if unicode.IsSpace(rune(*ptr)) {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		} else {
			if (func() int {
				temp = find_oper_num(int8(*ptr))
				return temp
			}()) == int(-1) {
				end = ptr
				for *ptr != 0 && !unicode.IsSpace(rune(*ptr)) && find_oper_num(int8(*ptr)) == int(-1) {
					ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
				}
				libc.StrNCpy(&name[0], end, int(int64(uintptr(unsafe.Pointer(ptr))-uintptr(unsafe.Pointer(end)))))
				name[int64(uintptr(unsafe.Pointer(ptr))-uintptr(unsafe.Pointer(end)))] = '\x00'
				for eindex = 0; *extra_bits[eindex] != '\n'; eindex++ {
					if libc.StrCaseCmp(&name[0], extra_bits[eindex]) == 0 {
						push(&vals, int(libc.BoolToInt(OBJ_FLAGGED(obj, uint32(int32(eindex))))))
						break
					}
				}
				if *extra_bits[eindex] == '\n' {
					push(&vals, int(libc.BoolToInt(isname(&name[0], obj.Name))))
				}
			} else {
				if temp != OPER_OPEN_PAREN {
					for top(&ops) > temp {
						evaluate_operation(&ops, &vals)
					}
				}
				if temp == OPER_CLOSE_PAREN {
					if (func() int {
						temp = pop(&ops)
						return temp
					}()) != OPER_OPEN_PAREN {
						basic_mud_log(libc.CString("SYSERR: Illegal parenthesis in shop keyword expression."))
						return 0
					}
				} else {
					push(&ops, temp)
				}
				ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
			}
		}
	}
	for top(&ops) != -1 {
		evaluate_operation(&ops, &vals)
	}
	temp = pop(&vals)
	if top(&vals) != -1 {
		basic_mud_log(libc.CString("SYSERR: Extra operands left on shop keyword expression stack."))
		return 0
	}
	return temp
}
func trade_with(item *obj_data, shop_nr int) int {
	var counter int
	if item.Cost < 1 {
		return OBJECT_NOVAL
	}
	if OBJ_FLAGGED(item, ITEM_NOSELL) {
		return OBJECT_NOTOK
	}
	for counter = 0; (shop_index[shop_nr].Type[counter]).Type != int(-1); counter++ {
		if (shop_index[shop_nr].Type[counter]).Type == int(item.Type_flag) {
			if (item.Value[VAL_WAND_CHARGES]) == 0 && (int(item.Type_flag) == ITEM_WAND || int(item.Type_flag) == ITEM_STAFF) {
				return OBJECT_DEAD
			} else if evaluate_expression(item, (shop_index[shop_nr].Type[counter]).Keywords) != 0 {
				return OBJECT_OK
			}
		}
	}
	return OBJECT_NOTOK
}
func same_obj(obj1 *obj_data, obj2 *obj_data) bool {
	var (
		aindex int
		i      int
		ef1    int
		ef2    int
	)
	if obj1 == nil || obj2 == nil {
		return obj1 == obj2
	}
	if obj1.Item_number != obj2.Item_number {
		return false
	}
	if obj1.Cost != obj2.Cost {
		return false
	}
	for i = 0; i < EF_ARRAY_MAX; i++ {
		ef1 = int(obj1.Extra_flags[i])
		ef2 = int(obj2.Extra_flags[i])
		if i == (int(ITEM_UNIQUE_SAVE / 32)) {
			ef1 &= ^(1 << (int(ITEM_UNIQUE_SAVE % 32)))
			ef2 &= ^(1 << (int(ITEM_UNIQUE_SAVE % 32)))
		}
		if ef1 != ef2 {
			return false
		}
	}
	for aindex = 0; aindex < MAX_OBJ_AFFECT; aindex++ {
		if obj1.Affected[aindex].Location != obj2.Affected[aindex].Location || obj1.Affected[aindex].Modifier != obj2.Affected[aindex].Modifier {
			return false
		}
	}
	return true
}
func shop_producing(item *obj_data, shop_nr int) bool {
	var counter int
	if item.Item_number == int(-1) {
		return false
	}
	for counter = 0; (shop_index[shop_nr].Producing[counter]) != int(-1); counter++ {
		if same_obj(item, &obj_proto[shop_index[shop_nr].Producing[counter]]) {
			return true
		}
	}
	return false
}
func transaction_amt(arg *byte) int {
	var (
		buf     [2048]byte
		buywhat *byte
	)
	buywhat = one_argument(arg, &buf[0])
	if *buywhat != 0 && buf[0] != 0 && is_number(&buf[0]) {
		libc.StrCpy(arg, (*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(arg), libc.StrLen(&buf[0])))), 1)))
		return libc.Atoi(libc.GoString(&buf[0]))
	}
	return 1
}
func times_message(obj *obj_data, name *byte, num int) *byte {
	var (
		buf  [256]byte
		len_ uint64
		ptr  *byte
	)
	if obj != nil {
		len_ = strlcpy(&buf[0], obj.Short_description, uint64(256))
	} else {
		if (func() *byte {
			ptr = libc.StrChr(name, '.')
			return ptr
		}()) == nil {
			ptr = name
		} else {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
		len_ = uint64(stdio.Snprintf(&buf[0], int(256), "%s %s", AN(ptr), ptr))
	}
	if num > 1 && len_ < uint64(256) {
		stdio.Snprintf(&buf[len_], int(256-uintptr(len_)), " (x %d)", num)
	}
	return &buf[0]
}
func get_slide_obj_vis(ch *char_data, name *byte, list *obj_data) *obj_data {
	var (
		i          *obj_data
		last_match *obj_data = nil
		j          int
		number     int
		tmpname    [2048]byte
		tmp        *byte
	)
	strlcpy(&tmpname[0], name, uint64(2048))
	tmp = &tmpname[0]
	if (func() int {
		number = get_number(&tmp)
		return number
	}()) == 0 {
		return nil
	}
	for func() int {
		i = list
		return func() int {
			j = 1
			return j
		}()
	}(); i != nil && j <= number; i = i.Next_content {
		if isname(tmp, i.Name) {
			if CAN_SEE_OBJ(ch, i) && !same_obj(last_match, i) {
				if j == number {
					return i
				}
				last_match = i
				j++
			}
		}
	}
	return nil
}
func get_hash_obj_vis(ch *char_data, name *byte, list *obj_data) *obj_data {
	var (
		loop     *obj_data
		last_obj *obj_data = nil
		qindex   int
	)
	if is_number(name) {
		qindex = libc.Atoi(libc.GoString(name))
	} else if is_number((*byte)(unsafe.Add(unsafe.Pointer(name), 1))) {
		qindex = libc.Atoi(libc.GoString((*byte)(unsafe.Add(unsafe.Pointer(name), 1))))
	} else {
		return nil
	}
	for loop = list; loop != nil; loop = loop.Next_content {
		if CAN_SEE_OBJ(ch, loop) && loop.Cost > 0 {
			if !same_obj(last_obj, loop) {
				if func() int {
					p := &qindex
					*p--
					return *p
				}() == 0 {
					return loop
				}
				last_obj = loop
			}
		}
	}
	return nil
}
func get_purchase_obj(ch *char_data, arg *byte, keeper *char_data, shop_nr int, msg int) *obj_data {
	var (
		name [2048]byte
		obj  *obj_data
	)
	one_argument(arg, &name[0])
	for {
		if name[0] == '#' || is_number(&name[0]) {
			obj = get_hash_obj_vis(ch, &name[0], keeper.Carrying)
		} else {
			obj = get_slide_obj_vis(ch, &name[0], keeper.Carrying)
		}
		if obj == nil {
			if msg != 0 {
				var buf [2048]byte
				stdio.Snprintf(&buf[0], int(2048), libc.GoString(shop_index[shop_nr].No_such_item1), GET_NAME(ch))
				do_tell(keeper, &buf[0], cmd_tell, 0)
			}
			return nil
		}
		if obj.Cost <= 0 {
			extract_obj(obj)
			obj = nil
		}
		if obj != nil {
			break
		}
	}
	return obj
}
func buy_price(obj *obj_data, shop_nr int, keeper *char_data, buyer *char_data) int {
	var (
		cost   int     = int(float32(obj.Cost) * shop_index[shop_nr].Profit_buy)
		adjust float64 = 1.0
		k      *obj_data
	)
	for k = object_list; k != nil; k = k.Next {
		if GET_OBJ_VNUM(k) == GET_OBJ_VNUM(obj) {
			adjust -= 0.00025
		}
	}
	if adjust < 0.015 {
		adjust = 0.5
	}
	cost = int(float64(cost) * adjust)
	if !IS_NPC(buyer) && (buyer.Bonuses[BONUS_THRIFTY]) > 0 {
		if int(buyer.Race) == RACE_ARLIAN {
			cost += int(float64(cost) * 0.2)
		}
		cost -= int(float64(cost) * 0.1)
		return cost
	} else if !IS_NPC(buyer) && (buyer.Bonuses[BONUS_IMPULSE]) != 0 {
		cost += int(float64(cost) * 0.25)
		return cost
	} else if !IS_NPC(buyer) && int(buyer.Race) == RACE_ARLIAN {
		cost += int(float64(cost) * 0.2)
		return cost
	} else {
		return int(float32(obj.Cost) * shop_index[shop_nr].Profit_buy)
	}
}
func sell_price(obj *obj_data, shop_nr int, keeper *char_data, seller *char_data) int {
	var (
		sell_cost_modifier float32 = shop_index[shop_nr].Profit_sell
		buy_cost_modifier  float32 = shop_index[shop_nr].Profit_buy
	)
	if sell_cost_modifier > buy_cost_modifier {
		sell_cost_modifier = buy_cost_modifier
	}
	var adjust float64 = 1.0
	var k *obj_data
	for k = object_list; k != nil; k = k.Next {
		if GET_OBJ_VNUM(k) == GET_OBJ_VNUM(obj) {
			adjust -= 0.00025
		}
	}
	if adjust < 0.15 {
		adjust = 0.15
	}
	if !IS_NPC(seller) && (seller.Bonuses[BONUS_THRIFTY]) > 0 {
		var haggle int = int(float32(obj.Cost) * (sell_cost_modifier / 2))
		if int(seller.Race) == RACE_ARLIAN {
			haggle -= int(float64(haggle) * 0.2)
		}
		haggle += int(float64(haggle) * 0.1)
		haggle = int(float64(haggle) * adjust)
		return haggle
	} else if !IS_NPC(seller) && int(seller.Race) == RACE_ARLIAN {
		var haggle int = int(float32(obj.Cost) * (sell_cost_modifier / 2))
		haggle -= int(float64(haggle) * 0.2)
		haggle = int(float64(haggle) * adjust)
		return haggle
	} else {
		return int(float64(float32(obj.Cost)*(sell_cost_modifier/2)) * adjust)
	}
}
func shopping_app(arg *byte, ch *char_data, keeper *char_data, shop_nr int) {
	var (
		obj   *obj_data
		i     int
		found int = 0
		buf   [64936]byte
	)
	if !is_ok(keeper, ch, shop_nr) {
		return
	}
	if shop_index[shop_nr].Lastsort < int(keeper.Carry_items) {
		sort_keeper_objs(keeper, shop_nr)
	}
	if *arg == 0 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s What do you want to appraise?", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if (func() *obj_data {
		obj = get_purchase_obj(ch, arg, keeper, shop_nr, 1)
		return obj
	}()) == nil {
		do_appraise(ch, arg, 0, 0)
		return
	}
	act(libc.CString("@C$N@W hands you @G$p@W for a moment and let's you examine it before taking it back.@n"), 1, ch, obj, unsafe.Pointer(keeper), TO_CHAR)
	act(libc.CString("@c$N@W hands @C$n@W @G$p@W for a moment and let's $m examine it before taking it back.@n"), 1, ch, obj, unsafe.Pointer(keeper), TO_ROOM)
	if GET_SKILL(ch, SKILL_APPRAISE) == 0 {
		send_to_char(ch, libc.CString("You are unskilled at appraising.\r\n"))
		return
	}
	improve_skill(ch, SKILL_APPRAISE, 1)
	if GET_SKILL(ch, SKILL_APPRAISE) < rand_number(1, 101) {
		send_to_char(ch, libc.CString("@wYou were completely stumped about the worth of %s@n\r\n"), obj.Short_description)
		WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		return
	} else {
		var displevel int = obj.Level
		if int(obj.Type_flag) == ITEM_WEAPON && OBJ_FLAGGED(obj, ITEM_CUSTOM) {
			displevel = 20
		}
		send_to_char(ch, libc.CString("@c---------------------------------------------------------------@n\n"))
		send_to_char(ch, libc.CString("@GItem Name   @W: @w%s@n\n"), obj.Short_description)
		send_to_char(ch, libc.CString("@GTrue Value  @W: @Y%s@n\n"), add_commas(int64(obj.Cost)))
		send_to_char(ch, libc.CString("@GItem Min LVL@W: @w%d@n\n"), displevel)
		if (obj.Value[VAL_ALL_HEALTH]) >= 100 {
			send_to_char(ch, libc.CString("@GCondition   @W: @C%d%s@n\n"), obj.Value[VAL_ALL_HEALTH], "%")
		} else if (obj.Value[VAL_ALL_HEALTH]) >= 50 {
			send_to_char(ch, libc.CString("@GCondition   @W: @y%d%s@n\n"), obj.Value[VAL_ALL_HEALTH], "%")
		} else if (obj.Value[VAL_ALL_HEALTH]) >= 1 {
			send_to_char(ch, libc.CString("@GCondition   @W: @r%d%s@n\n"), obj.Value[VAL_ALL_HEALTH], "%")
		} else {
			send_to_char(ch, libc.CString("@GCondition   @W: @D%d%s@n\n"), obj.Value[VAL_ALL_HEALTH], "%")
		}
		send_to_char(ch, libc.CString("@GItem Size   @W:@w %s@n\n"), size_names[obj.Size])
		send_to_char(ch, libc.CString("@GItem Weight @W: @w%s@n\n"), add_commas(obj.Weight))
		if OBJ_FLAGGED(obj, ITEM_SLOT1) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString("GToken Slots  @W: @m0/1@n\n"))
		} else if OBJ_FLAGGED(obj, ITEM_SLOT1) && OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString("GToken Slots  @W: @m1/1@n\n"))
		} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && !OBJ_FLAGGED(obj, ITEM_SLOT_ONE) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString("GToken Slots  @W: @m0/2@n\n"))
		} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && OBJ_FLAGGED(obj, ITEM_SLOT_ONE) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString("GToken Slots  @W: @m1/2@n\n"))
		} else if OBJ_FLAGGED(obj, ITEM_SLOT2) && !OBJ_FLAGGED(obj, ITEM_SLOTS_FILLED) {
			send_to_char(ch, libc.CString("GToken Slots  @W: @m2/2@n\n"))
		}
		var bits [64936]byte
		sprintbitarray(obj.Wear_flags[:], wear_bits[:], TW_ARRAY_MAX, &bits[0])
		search_replace(&bits[0], libc.CString("TAKE"), libc.CString(""))
		send_to_char(ch, libc.CString("@GWear Loc.   @W:@w%s\n"), &bits[0])
		if int(obj.Type_flag) == ITEM_WEAPON {
			if OBJ_FLAGGED(obj, ITEM_WEAPLVL1) {
				send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C1@D]\n@GDamage Bonus@W: @D[@w5%s@D]@n\r\n"), "%")
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL2) {
				send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C2@D]\n@GDamage Bonus@W: @D[@w10%s@D]@n\r\n"), "%")
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL3) {
				send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C3@D]\n@GDamage Bonus@W: @D[@w20%s@D]@n\r\n"), "%")
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL4) {
				send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C4@D]\n@GDamage Bonus@W: @D[@w30%s@D]@n\r\n"), "%")
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL5) {
				send_to_char(ch, libc.CString("@GWeapon Level@W: @D[@C5@D]\n@GDamage Bonus@W: @D[@w50%s@D]@n\r\n"), "%")
			}
		}
		send_to_char(ch, libc.CString("@GItem Bonuses@W:@w"))
		for i = 0; i < MAX_OBJ_AFFECT; i++ {
			if obj.Affected[i].Modifier != 0 {
				sprinttype(obj.Affected[i].Location, apply_types[:], &buf[0], uint64(64936))
				send_to_char(ch, libc.CString("%s %+d to %s"), func() string {
					if func() int {
						p := &found
						x := *p
						*p++
						return x
					}() != 0 {
						return ","
					}
					return ""
				}(), obj.Affected[i].Modifier, &buf[0])
				switch obj.Affected[i].Location {
				case APPLY_FEAT:
					send_to_char(ch, libc.CString(" (%s)"), feat_list[obj.Affected[i].Specific].Name)
				case APPLY_SKILL:
					send_to_char(ch, libc.CString(" (%s)"), spell_info[obj.Affected[i].Specific].Name)
				}
			}
		}
		if found == 0 {
			send_to_char(ch, libc.CString(" None@n"))
		} else {
			send_to_char(ch, libc.CString("@n"))
		}
		var buf2 [64936]byte
		sprintbitarray(obj.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &buf2[0])
		send_to_char(ch, libc.CString("\n@GSpecial     @W:@w %s"), &buf2[0])
		send_to_char(ch, libc.CString("\n@c---------------------------------------------------------------@n\n"))
	}
}
func shopping_buy(arg *byte, ch *char_data, keeper *char_data, shop_nr int) {
	var (
		tempstr  [2048]byte
		tempbuf  [2048]byte
		obj      *obj_data
		last_obj *obj_data = nil
		goldamt  int       = 0
		buynum   int
		bought   int = 0
	)
	if !is_ok(keeper, ch, shop_nr) {
		return
	}
	if shop_index[shop_nr].Lastsort < int(keeper.Carry_items) {
		sort_keeper_objs(keeper, shop_nr)
	}
	if (func() int {
		buynum = transaction_amt(arg)
		return buynum
	}()) < 0 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s A negative amount?  Try selling me something.", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if *arg == 0 || buynum == 0 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s What do you want to buy?", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if (func() *obj_data {
		obj = get_purchase_obj(ch, arg, keeper, shop_nr, 1)
		return obj
	}()) == nil {
		return
	}
	if buy_price(obj, shop_nr, keeper, ch) > ch.Gold && !ADM_FLAGGED(ch, ADM_MONEY) {
		var actbuf [2048]byte
		stdio.Snprintf(&actbuf[0], int(2048), libc.GoString(shop_index[shop_nr].Missing_cash2), GET_NAME(ch))
		do_tell(keeper, &actbuf[0], cmd_tell, 0)
		switch shop_index[shop_nr].Temper1 {
		case 0:
			do_action(keeper, libc.StrCpy(&actbuf[0], GET_NAME(ch)), cmd_puke, 0)
			return
		case 1:
			do_echo(keeper, libc.StrCpy(&actbuf[0], libc.CString("smokes on his joint.")), cmd_emote, SCMD_EMOTE)
			return
		default:
			return
		}
	}
	if int(ch.Carry_items)+1 > 50 {
		send_to_char(ch, libc.CString("%s: You can't carry any more items.\r\n"), fname(obj.Name))
		return
	}
	if ch.Carry_weight+int(obj.Weight) > int(max_carry_weight(ch)) {
		send_to_char(ch, libc.CString("%s: You can't carry that much weight.\r\n"), fname(obj.Name))
		return
	}
	for obj != nil && (ch.Gold >= buy_price(obj, shop_nr, keeper, ch) || ADM_FLAGGED(ch, ADM_MONEY)) && int(ch.Carry_items) < 50 && bought < buynum && ch.Carry_weight+int(obj.Weight) <= int(max_carry_weight(ch)) {
		var charged int
		bought++
		if shop_producing(obj, shop_nr) {
			obj = read_object(obj.Item_number, REAL)
			add_unique_id(obj)
		} else {
			obj_from_char(obj)
			shop_index[shop_nr].Lastsort--
		}
		obj_to_char(obj, ch)
		if OBJ_FLAGGED(obj, ITEM_MATURE) {
			obj.Value[VAL_MAXMATURE] = 6
		}
		charged = buy_price(obj, shop_nr, keeper, ch)
		goldamt += charged
		if !ADM_FLAGGED(ch, ADM_MONEY) {
			ch.Gold -= charged
		} else {
			send_to_imm(libc.CString("IMM PURCHASE: %s has purchased %s for free."), GET_NAME(ch), obj.Short_description)
			log_imm_action(libc.CString("IMM PURCHASE: %s has purchased %s for free."), GET_NAME(ch), obj.Short_description)
		}
		last_obj = obj
		obj = get_purchase_obj(ch, arg, keeper, shop_nr, 0)
		if !same_obj(obj, last_obj) {
			break
		}
	}
	if bought < buynum {
		var buf [2048]byte
		if obj == nil || !same_obj(last_obj, obj) {
			stdio.Snprintf(&buf[0], int(2048), "%s I only have %d to sell you.", GET_NAME(ch), bought)
		} else if ch.Gold < buy_price(obj, shop_nr, keeper, ch) {
			stdio.Snprintf(&buf[0], int(2048), "%s You can only afford %d.", GET_NAME(ch), bought)
		} else if int(ch.Carry_items) >= 50 {
			stdio.Snprintf(&buf[0], int(2048), "%s You can only hold %d.", GET_NAME(ch), bought)
		} else if ch.Carry_weight+int(obj.Weight) > int(max_carry_weight(ch)) {
			stdio.Snprintf(&buf[0], int(2048), "%s You can only carry %d.", GET_NAME(ch), bought)
		} else {
			stdio.Snprintf(&buf[0], int(2048), "%s Something screwy only gave you %d.", GET_NAME(ch), bought)
		}
		do_tell(keeper, &buf[0], cmd_tell, 0)
	}
	if !ADM_FLAGGED(ch, ADM_MONEY) {
		keeper.Gold += goldamt
	}
	strlcpy(&tempstr[0], times_message(ch.Carrying, nil, bought), uint64(2048))
	stdio.Snprintf(&tempbuf[0], int(2048), "$n buys %s.", &tempstr[0])
	act(&tempbuf[0], 0, ch, obj, nil, TO_ROOM)
	stdio.Snprintf(&tempbuf[0], int(2048), libc.GoString(shop_index[shop_nr].Message_buy), GET_NAME(ch), goldamt)
	do_tell(keeper, &tempbuf[0], cmd_tell, 0)
	send_to_char(ch, libc.CString("You now have %s.\r\n"), &tempstr[0])
	if IS_SET(shop_index[shop_nr].Bitvector, 1<<1) {
		if keeper.Gold > MAX_OUTSIDE_BANK {
			shop_index[shop_nr].BankAccount += keeper.Gold - MAX_OUTSIDE_BANK
			keeper.Gold = MAX_OUTSIDE_BANK
		}
	}
}
func get_selling_obj(ch *char_data, name *byte, keeper *char_data, shop_nr int, msg int) *obj_data {
	var (
		buf    [2048]byte
		obj    *obj_data
		result int
	)
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, name, nil, ch.Carrying)
		return obj
	}()) == nil {
		if msg != 0 {
			var tbuf [2048]byte
			stdio.Snprintf(&tbuf[0], int(2048), libc.GoString(shop_index[shop_nr].No_such_item2), GET_NAME(ch))
			do_tell(keeper, &tbuf[0], cmd_tell, 0)
		}
		return nil
	}
	if (func() int {
		result = trade_with(obj, shop_nr)
		return result
	}()) == OBJECT_OK {
		return obj
	}
	if msg == 0 {
		return nil
	}
	switch result {
	case OBJECT_NOVAL:
		stdio.Snprintf(&buf[0], int(2048), "%s You've got to be kidding, that thing is worthless!", GET_NAME(ch))
	case OBJECT_NOTOK:
		stdio.Snprintf(&buf[0], int(2048), libc.GoString(shop_index[shop_nr].Do_not_buy), GET_NAME(ch))
	case OBJECT_DEAD:
		stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_NO_USED_WANDSTAFF)
	default:
		basic_mud_log(libc.CString("SYSERR: Illegal return value of %d from trade_with() (%s)"), result, "__THISFILE__")
		stdio.Snprintf(&buf[0], int(2048), "%s An error has occurred.", GET_NAME(ch))
	}
	do_tell(keeper, &buf[0], cmd_tell, 0)
	return nil
}
func slide_obj(obj *obj_data, keeper *char_data, shop_nr int) *obj_data {
	var (
		loop *obj_data
		temp int
	)
	if shop_index[shop_nr].Lastsort < int(keeper.Carry_items) {
		sort_keeper_objs(keeper, shop_nr)
	}
	if shop_producing(obj, shop_nr) {
		temp = obj.Item_number
		extract_obj(obj)
		return &obj_proto[temp]
	}
	shop_index[shop_nr].Lastsort++
	loop = keeper.Carrying
	obj_to_char(obj, keeper)
	keeper.Carrying = loop
	for loop != nil {
		if same_obj(obj, loop) {
			obj.Next_content = loop.Next_content
			loop.Next_content = obj
			return obj
		}
		loop = loop.Next_content
	}
	keeper.Carrying = obj
	return obj
}
func sort_keeper_objs(keeper *char_data, shop_nr int) {
	var (
		list *obj_data = nil
		temp *obj_data
	)
	for shop_index[shop_nr].Lastsort < int(keeper.Carry_items) {
		temp = keeper.Carrying
		obj_from_char(temp)
		temp.Next_content = list
		list = temp
	}
	for list != nil {
		temp = list
		list = list.Next_content
		if shop_producing(temp, shop_nr) && get_obj_in_list_num(temp.Item_number, keeper.Carrying) == nil {
			obj_to_char(temp, keeper)
			shop_index[shop_nr].Lastsort++
		} else {
			slide_obj(temp, keeper, shop_nr)
		}
	}
}
func shopping_sell(arg *byte, ch *char_data, keeper *char_data, shop_nr int) {
	var (
		tempstr [2048]byte
		name    [2048]byte
		tempbuf [2048]byte
		obj     *obj_data
		sellnum int
		sold    int = 0
		goldamt int = 0
	)
	if !is_ok(keeper, ch, shop_nr) {
		return
	}
	if (func() int {
		sellnum = transaction_amt(arg)
		return sellnum
	}()) < 0 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s A negative amount?  Try buying something.", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if *arg == 0 || sellnum == 0 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s What do you want to sell??", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	one_argument(arg, &name[0])
	if (func() *obj_data {
		obj = get_selling_obj(ch, &name[0], keeper, shop_nr, 1)
		return obj
	}()) == nil {
		return
	}
	if int(obj.Type_flag) == ITEM_PLANT && (obj.Value[VAL_WATERLEVEL]) <= -10 {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), "%s That thing is dead!", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	if !is_ok_obj(keeper, ch, obj, shop_nr) {
		return
	}
	if keeper.Gold+shop_index[shop_nr].BankAccount < sell_price(obj, shop_nr, keeper, ch) {
		var buf [2048]byte
		stdio.Snprintf(&buf[0], int(2048), libc.GoString(shop_index[shop_nr].Missing_cash1), GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	for obj != nil && keeper.Gold+shop_index[shop_nr].BankAccount >= sell_price(obj, shop_nr, keeper, ch) && sold < sellnum {
		var charged int = sell_price(obj, shop_nr, keeper, ch)
		goldamt += charged
		keeper.Gold -= charged
		sold++
		obj_from_char(obj)
		slide_obj(obj, keeper, shop_nr)
		obj = get_selling_obj(ch, &name[0], keeper, shop_nr, 0)
	}
	if sold < sellnum {
		var buf [2048]byte
		if obj == nil {
			stdio.Snprintf(&buf[0], int(2048), "%s You only have %d of those.", GET_NAME(ch), sold)
		} else if keeper.Gold+shop_index[shop_nr].BankAccount < sell_price(obj, shop_nr, keeper, ch) {
			stdio.Snprintf(&buf[0], int(2048), "%s I can only afford to buy %d of those.", GET_NAME(ch), sold)
		} else {
			stdio.Snprintf(&buf[0], int(2048), "%s Something really screwy made me buy %d.", GET_NAME(ch), sold)
		}
		do_tell(keeper, &buf[0], cmd_tell, 0)
	}
	strlcpy(&tempstr[0], times_message(nil, &name[0], sold), uint64(2048))
	stdio.Snprintf(&tempbuf[0], int(2048), "$n sells something to %s.\r\n", GET_NAME(keeper))
	act(&tempbuf[0], 0, ch, obj, nil, TO_ROOM)
	stdio.Snprintf(&tempbuf[0], int(2048), libc.GoString(shop_index[shop_nr].Message_sell), GET_NAME(ch), goldamt)
	do_tell(keeper, &tempbuf[0], cmd_tell, 0)
	send_to_char(ch, libc.CString("The shopkeeper gives you %s zenni.\r\n"), add_commas(int64(goldamt)))
	if ch.Gold+goldamt > GOLD_CARRY(ch) {
		goldamt = (ch.Gold + goldamt) - GOLD_CARRY(ch)
		ch.Gold = GOLD_CARRY(ch)
		ch.Bank_gold += goldamt
		send_to_char(ch, libc.CString("You couldn't hold all of the money. The rest was deposited for you.\r\n"))
	} else {
		ch.Gold += goldamt
	}
	if keeper.Gold < MIN_OUTSIDE_BANK {
		goldamt = int(MIN(int64(MAX_OUTSIDE_BANK-keeper.Gold), int64(shop_index[shop_nr].BankAccount)))
		shop_index[shop_nr].BankAccount -= goldamt
		keeper.Gold += goldamt
	}
}
func shopping_value(arg *byte, ch *char_data, keeper *char_data, shop_nr int) {
	var (
		buf  [64936]byte
		name [2048]byte
		obj  *obj_data
	)
	if !is_ok(keeper, ch, shop_nr) {
		return
	}
	if *arg == 0 {
		stdio.Snprintf(&buf[0], int(64936), "%s What do you want me to evaluate??", GET_NAME(ch))
		do_tell(keeper, &buf[0], cmd_tell, 0)
		return
	}
	one_argument(arg, &name[0])
	if (func() *obj_data {
		obj = get_selling_obj(ch, &name[0], keeper, shop_nr, 1)
		return obj
	}()) == nil {
		return
	}
	if !is_ok_obj(keeper, ch, obj, shop_nr) {
		return
	}
	stdio.Snprintf(&buf[0], int(64936), "%s I'll give you %d zenni for that!", GET_NAME(ch), sell_price(obj, shop_nr, keeper, ch))
	do_tell(keeper, &buf[0], cmd_tell, 0)
}
func list_object(obj *obj_data, cnt int, aindex int, shop_nr int, keeper *char_data, ch *char_data) *byte {
	var (
		result   [256]byte
		itemname [128]byte
		quantity [16]byte
	)
	if shop_producing(obj, shop_nr) {
		libc.StrCpy(&quantity[0], libc.CString("Unlimited"))
	} else {
		stdio.Sprintf(&quantity[0], "%d", cnt)
	}
	switch obj.Type_flag {
	case ITEM_DRINKCON:
		if (obj.Value[VAL_DRINKCON_HOWFULL]) != 0 {
			stdio.Snprintf(&itemname[0], int(128), "%s", obj.Short_description)
		} else {
			strlcpy(&itemname[0], obj.Short_description, uint64(128))
		}
	case ITEM_WAND:
		fallthrough
	case ITEM_STAFF:
		stdio.Snprintf(&itemname[0], int(128), "%s%s", obj.Short_description, func() string {
			if (obj.Value[VAL_WAND_CHARGES]) < (obj.Value[VAL_WAND_MAXCHARGES]) {
				return " (partially used)"
			}
			return ""
		}())
	default:
		strlcpy(&itemname[0], obj.Short_description, uint64(128))
	}
	if OBJ_FLAGGED(obj, ITEM_BROKEN) {
		var titemname [128]byte
		strlcpy(&titemname[0], &itemname[0], uint64(128))
		stdio.Snprintf(&itemname[0], int(128), "%s [broken]", &titemname[0])
	}
	CAP(&itemname[0])
	var displevel int = obj.Level
	if int(obj.Type_flag) == ITEM_WEAPON && OBJ_FLAGGED(obj, ITEM_CUSTOM) {
		displevel = 20
	}
	stdio.Snprintf(&result[0], int(256), " %2d)  %9s %-*s %3d %13s\r\n", aindex, &quantity[0], count_color_chars(&itemname[0])+36, &itemname[0], displevel, add_commas(int64(buy_price(obj, shop_nr, keeper, ch))))
	return &result[0]
}
func shopping_list(arg *byte, ch *char_data, keeper *char_data, shop_nr int) {
	var (
		buf      [259744]byte
		name     [2048]byte
		obj      *obj_data
		last_obj *obj_data = nil
		cnt      int       = 0
		lindex   int       = 0
		found    int       = 0
		len_     uint64
	)
	if !is_ok(keeper, ch, shop_nr) {
		return
	}
	if shop_index[shop_nr].Lastsort < int(keeper.Carry_items) {
		sort_keeper_objs(keeper, shop_nr)
	}
	one_argument(arg, &name[0])
	len_ = strlcpy(&buf[0], libc.CString(" ##   Available   Item                             Min. Lvl       Cost\r\n----------------------------------------------------------------------\r\n"), uint64(259744))
	if keeper.Carrying != nil {
		for obj = keeper.Carrying; obj != nil; obj = obj.Next_content {
			if CAN_SEE_OBJ(ch, obj) && obj.Cost > 0 {
				if last_obj == nil {
					last_obj = obj
					cnt = 1
				} else if same_obj(last_obj, obj) {
					cnt++
				} else {
					lindex++
					if name[0] == 0 || isname(&name[0], last_obj.Name) {
						libc.StrNCat(&buf[0], list_object(last_obj, cnt, lindex, shop_nr, keeper, ch), int(259744-uintptr(len_)-1))
						len_ = uint64(libc.StrLen(&buf[0]))
						if len_+1 >= uint64(259744) {
							break
						}
						found = 1
					}
					cnt = 1
					last_obj = obj
				}
			}
		}
	}
	lindex++
	if last_obj == nil {
		send_to_char(ch, libc.CString("Currently, there is nothing for sale.\r\n"))
	} else if name[0] != 0 && found == 0 {
		send_to_char(ch, libc.CString("Presently, none of those are for sale.\r\n"))
	} else {
		var zen [80]byte
		if name[0] == 0 || isname(&name[0], last_obj.Name) {
			if len_ < uint64(259744) {
				libc.StrNCat(&buf[0], list_object(last_obj, cnt, lindex, shop_nr, keeper, ch), int(259744-uintptr(len_)-1))
			}
		}
		if len_ < uint64(259744) {
			stdio.Sprintf(&zen[0], "@W[@wYour Zenni@D: @Y%s@W]", add_commas(int64(ch.Gold)))
			libc.StrNCat(&buf[0], &zen[0], int(259744-uintptr(len_)-1))
		}
		page_string(ch.Desc, &buf[0], 1)
	}
}
func ok_shop_room(shop_nr int, room int) bool {
	var mindex int
	for mindex = 0; (shop_index[shop_nr].In_room[mindex]) != int(-1); mindex++ {
		if (shop_index[shop_nr].In_room[mindex]) == room {
			return true
		}
	}
	return false
}
func shop_keeper(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
	var (
		keeper  *char_data = (*char_data)(me)
		shop_nr int
	)
	for shop_nr = 0; shop_nr <= top_shop; shop_nr++ {
		if shop_index[shop_nr].Keeper == keeper.Nr {
			break
		}
	}
	if shop_nr > top_shop {
		return false
	}
	if shop_index[shop_nr].Func != nil {
		if shop_index[shop_nr].Func(ch, me, cmd, argument) {
			return true
		}
	}
	if keeper == ch {
		if cmd != 0 {
			shop_index[shop_nr].Lastsort = 0
		}
		return false
	}
	if !ok_shop_room(shop_nr, int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))) {
		return false
	}
	if !AWAKE(keeper) {
		return false
	}
	if libc.StrCmp(libc.CString("steal"), complete_cmd_info[cmd].Command) == 0 {
		var argm [2048]byte
		if !IS_SET(shop_index[shop_nr].Bitvector, 1<<2) {
			stdio.Snprintf(&argm[0], int(2048), "$N shouts '%s'", MSG_NO_STEAL_HERE)
			act(&argm[0], 0, ch, nil, unsafe.Pointer(keeper), TO_CHAR)
			act(&argm[0], 0, ch, nil, unsafe.Pointer(keeper), TO_ROOM)
			do_action(keeper, GET_NAME(ch), cmd_slap, 0)
			return true
		} else {
			return false
		}
	}
	if libc.StrCmp(libc.CString("buy"), complete_cmd_info[cmd].Command) == 0 {
		shopping_buy(argument, ch, keeper, shop_nr)
		return true
	} else if libc.StrCmp(libc.CString("sell"), complete_cmd_info[cmd].Command) == 0 {
		shopping_sell(argument, ch, keeper, shop_nr)
		return true
	} else if libc.StrCmp(libc.CString("value"), complete_cmd_info[cmd].Command) == 0 {
		shopping_value(argument, ch, keeper, shop_nr)
		return true
	} else if libc.StrCmp(libc.CString("list"), complete_cmd_info[cmd].Command) == 0 {
		shopping_list(argument, ch, keeper, shop_nr)
		return true
	} else if libc.StrCmp(libc.CString("appraise"), complete_cmd_info[cmd].Command) == 0 {
		shopping_app(argument, ch, keeper, shop_nr)
		return true
	}
	return false
}
func ok_damage_shopkeeper(ch *char_data, victim *char_data) bool {
	var sindex int
	if !IS_MOB(victim) || libc.FuncAddr(mob_index[victim.Nr].Func) != libc.FuncAddr(shop_keeper) {
		return true
	}
	if AFF_FLAGGED(victim, AFF_CHARM) {
		return true
	}
	for sindex = 0; sindex <= top_shop; sindex++ {
		if victim.Nr == shop_index[sindex].Keeper && !IS_SET(shop_index[sindex].Bitvector, 1<<0) {
			var buf [2048]byte
			stdio.Snprintf(&buf[0], int(2048), "%s %s", GET_NAME(ch), MSG_CANT_KILL_KEEPER)
			do_tell(victim, &buf[0], cmd_tell, 0)
			do_action(victim, GET_NAME(ch), cmd_slap, 0)
			return false
		}
	}
	return true
}
func add_to_list(list *shop_buy_data, type_ int, len_ *int, val *int) bool {
	if *val != int(-1) && *val >= 0 {
		if *len_ < MAX_SHOP_OBJ {
			if type_ == LIST_PRODUCE {
				*val = real_object(*val)
			}
			if *val != int(-1) {
				(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof(shop_buy_data{})*uintptr(*len_)))).Type = *val
				(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof(shop_buy_data{})*uintptr(func() int {
					p := len_
					x := *p
					*p++
					return x
				}())))).Keywords = nil
			} else {
				*val = -1
			}
			return false
		} else {
			return true
		}
	}
	return false
}
func end_read_list(list *shop_buy_data, len_ int, error int) int {
	if error != 0 {
		basic_mud_log(libc.CString("SYSERR: Raise MAX_SHOP_OBJ constant in shop.h to %d"), len_+error)
	}
	(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof(shop_buy_data{})*uintptr(len_)))).Keywords = nil
	(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof(shop_buy_data{})*uintptr(func() int {
		p := &len_
		x := *p
		*p++
		return x
	}())))).Type = -1
	return len_
}
func read_line(shop_f *stdio.File, string_ *byte, data unsafe.Pointer) {
	var buf [256]byte
	if get_line(shop_f, &buf[0]) == 0 || stdio.Sscanf(&buf[0], libc.GoString(string_), data) == 0 {
		basic_mud_log(libc.CString("SYSERR: Error in shop #%d, near '%s' with '%s'"), shop_index[top_shop].Vnum, &buf[0], string_)
		os.Exit(1)
	}
}
func read_list(shop_f *stdio.File, list *shop_buy_data, new_format int, max int, type_ int) int {
	var (
		count int
		temp  int
		len_  int = 0
		error int = 0
	)
	if new_format != 0 {
		for {
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&temp))
			if temp < 0 {
				break
			}
			error += int(libc.BoolToInt(add_to_list(list, type_, &len_, &temp)))
		}
	} else {
		for count = 0; count < max; count++ {
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&temp))
			error += int(libc.BoolToInt(add_to_list(list, type_, &len_, &temp)))
		}
	}
	return end_read_list(list, len_, error)
}
func read_type_list(shop_f *stdio.File, list *shop_buy_data, new_format int, max int) int {
	var (
		tindex int
		num    int
		len_   int = 0
		error  int = 0
		ptr    *byte
		buf    [64936]byte
	)
	if new_format == 0 {
		return read_list(shop_f, list, 0, max, LIST_TRADE)
	}
	for {
		shop_f.GetS(&buf[0], int32(uint32(64936)))
		if (func() *byte {
			ptr = libc.StrChr(&buf[0], ';')
			return ptr
		}()) != nil {
			*ptr = '\x00'
		} else {
			*((*byte)(unsafe.Add(unsafe.Pointer(&buf[libc.StrLen(&buf[0])]), -1))) = '\x00'
		}
		num = -1
		if libc.StrNCmp(&buf[0], libc.CString("-1"), 2) != 0 {
			for tindex = 0; *item_types[tindex] != '\n'; tindex++ {
				if libc.StrNCaseCmp(item_types[tindex], &buf[0], libc.StrLen(item_types[tindex])) == 0 {
					num = tindex
					libc.StrCpy(&buf[0], &buf[libc.StrLen(item_types[tindex])])
					break
				}
			}
		}
		ptr = &buf[0]
		if num == -1 {
			stdio.Sscanf(&buf[0], "%d", &num)
			for !unicode.IsDigit(rune(*ptr)) {
				ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
			}
			for unicode.IsDigit(rune(*ptr)) {
				ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
			}
		}
		for unicode.IsSpace(rune(*ptr)) {
			ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), 1))
		}
		for unicode.IsSpace(rune(*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(ptr), libc.StrLen(ptr)))), -1))))) {
			*((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(ptr), libc.StrLen(ptr)))), -1))) = '\x00'
		}
		error += int(libc.BoolToInt(add_to_list(list, LIST_TRADE, &len_, &num)))
		if *ptr != 0 {
			(*(*shop_buy_data)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof(shop_buy_data{})*uintptr(len_-1)))).Keywords = libc.StrDup(ptr)
		}
		if num < 0 {
			break
		}
	}
	return end_read_list(list, len_, error)
}
func read_shop_message(mnum int, shr int, shop_f *stdio.File, why *byte) *byte {
	var (
		cht  int
		ss   int = 0
		ds   int = 0
		err  int = 0
		tbuf *byte
	)
	if (func() *byte {
		tbuf = fread_string(shop_f, why)
		return tbuf
	}()) == nil {
		return nil
	}
	for cht = 0; *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht)) != 0; cht++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht)) != '%' {
			continue
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht+1)) == 's' {
			ss++
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht+1)) == 'd' && (mnum == 5 || mnum == 6) {
			if ss == 0 {
				basic_mud_log(libc.CString("SYSERR: Shop #%d has %%d before %%s, message #%d."), shr, mnum)
				err++
			}
			ds++
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht+1)) != '%' {
			basic_mud_log(libc.CString("SYSERR: Shop #%d has invalid format '%%%c' in message #%d."), shr, *(*byte)(unsafe.Add(unsafe.Pointer(tbuf), cht+1)), mnum)
			err++
		}
	}
	if ss > 1 || ds > 1 {
		basic_mud_log(libc.CString("SYSERR: Shop #%d has too many specifiers for message #%d. %%s=%d %%d=%d"), shr, mnum, ss, ds)
		err++
	}
	if err != 0 {
		libc.Free(unsafe.Pointer(tbuf))
		return nil
	}
	return tbuf
}
func boot_the_shops(shop_f *stdio.File, filename *byte, rec_count int) {
	var (
		buf        *byte
		buf2       [256]byte
		p          *byte
		temp       int
		count      int
		new_format int = 0
		list       [101]shop_buy_data
		done       int = 0
	)
	stdio.Snprintf(&buf2[0], int(256), "beginning of shop file %s", filename)
	for done == 0 {
		buf = fread_string(shop_f, &buf2[0])
		if *buf == '#' {
			stdio.Sscanf(buf, "#%d\n", &temp)
			stdio.Snprintf(&buf2[0], int(256), "shop #%d in shop file %s", temp, filename)
			libc.Free(unsafe.Pointer(buf))
			top_shop++
			if top_shop == 0 {
				shop_index = make([]shop_data, rec_count)
			}
			shop_index[top_shop].Vnum = temp
			temp = read_list(shop_f, &list[0], new_format, MAX_PROD, LIST_PRODUCE)
			shop_index[top_shop].Producing = make([]int, temp)
			for count = 0; count < temp; count++ {
				shop_index[top_shop].Producing[count] = (list[count]).Type
			}
			read_line(shop_f, libc.CString("%f"), unsafe.Pointer(&shop_index[top_shop].Profit_buy))
			read_line(shop_f, libc.CString("%f"), unsafe.Pointer(&shop_index[top_shop].Profit_sell))
			temp = read_type_list(shop_f, &list[0], new_format, MAX_TRADE)
			shop_index[top_shop].Type = make([]shop_buy_data, temp)
			for count = 0; count < temp; count++ {
				(shop_index[top_shop].Type[count]).Type = (list[count]).Type
				(shop_index[top_shop].Type[count]).Keywords = (list[count]).Keywords
			}
			shop_index[top_shop].No_such_item1 = read_shop_message(0, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].No_such_item2 = read_shop_message(1, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].Do_not_buy = read_shop_message(2, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].Missing_cash1 = read_shop_message(3, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].Missing_cash2 = read_shop_message(4, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].Message_buy = read_shop_message(5, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			shop_index[top_shop].Message_sell = read_shop_message(6, shop_index[top_shop].Vnum, shop_f, &buf2[0])
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&shop_index[top_shop].Temper1))
			read_line(shop_f, libc.CString("%ld"), unsafe.Pointer(&shop_index[top_shop].Bitvector))
			read_line(shop_f, libc.CString("%hd"), unsafe.Pointer(&shop_index[top_shop].Keeper))
			shop_index[top_shop].Keeper = real_mobile(shop_index[top_shop].Keeper)
			buf = (*byte)(unsafe.Pointer(&make([]int8, READ_SIZE)[0]))
			get_line(shop_f, buf)
			p = buf
			for temp = 0; temp < SW_ARRAY_MAX; temp++ {
				if p == nil || *p == 0 {
					break
				}
				if stdio.Sscanf(p, "%d", &count) != 1 {
					basic_mud_log(libc.CString("SYSERR: Can't parse TRADE_WITH line in %s: '%s'"), &buf2[0], buf)
					break
				}
				shop_index[top_shop].With_who[temp] = uint32(count)
				for unicode.IsDigit(rune(*p)) || *p == '-' {
					p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
				}
				for *p != 0 && (!unicode.IsDigit(rune(*p)) && *p != '-') {
					p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
				}
			}
			libc.Free(unsafe.Pointer(buf))
			for temp < SW_ARRAY_MAX {
				shop_index[top_shop].With_who[func() int {
					p := &temp
					x := *p
					*p++
					return x
				}()] = 0
			}
			temp = read_list(shop_f, &list[0], new_format, 1, LIST_ROOM)
			shop_index[top_shop].In_room = make([]int, temp)
			for count = 0; count < temp; count++ {
				shop_index[top_shop].In_room[count] = (list[count]).Type
			}
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&shop_index[top_shop].Open1))
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&shop_index[top_shop].Close1))
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&shop_index[top_shop].Open2))
			read_line(shop_f, libc.CString("%d"), unsafe.Pointer(&shop_index[top_shop].Close2))
			shop_index[top_shop].BankAccount = 0
			shop_index[top_shop].Lastsort = 0
			shop_index[top_shop].Func = nil
		} else {
			if *buf == '$' {
				done = 1
			} else if libc.StrStr(buf, libc.CString(VERSION3_TAG)) != nil {
				new_format = 1
			}
			libc.Free(unsafe.Pointer(buf))
		}
	}
}
func assign_the_shopkeepers() {
	var cindex int
	cmd_say = find_command(libc.CString("say"))
	cmd_tell = find_command(libc.CString("tell"))
	cmd_emote = find_command(libc.CString("emote"))
	cmd_slap = find_command(libc.CString("slap"))
	cmd_puke = find_command(libc.CString("puke"))
	for cindex = 0; cindex <= top_shop; cindex++ {
		if shop_index[cindex].Keeper == int(-1) {
			continue
		}
		if mob_index[shop_index[cindex].Keeper].Func != nil && libc.FuncAddr(mob_index[shop_index[cindex].Keeper].Func) != libc.FuncAddr(shop_keeper) {
			shop_index[cindex].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return mob_index[shop_index[cindex].Keeper].Func(ch, me, cmd, argument)
			}
		}
		mob_index[shop_index[cindex].Keeper].Func = func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
			return func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) bool {
				return shop_keeper(ch, me, cmd, argument)
			}(ch, me, cmd, argument)
		}
	}
}
func customer_string(shop_nr int, detailed int) *byte {
	var (
		sindex int = 0
		flag   int = 0
		nlen   int
		len_   uint64 = 0
		buf    [256]byte
	)
	for *trade_letters[sindex] != '\n' && len_+1 < uint64(256) {
		if detailed != 0 {
			if !IS_SET_AR(shop_index[shop_nr].With_who[:], uint32(int32(flag))) {
				nlen = stdio.Snprintf(&buf[len_], int(256-uintptr(len_)), ", %s", trade_letters[sindex])
				if len_+uint64(nlen) >= uint64(256) || nlen < 0 {
					break
				}
				len_ += uint64(nlen)
			}
		} else {
			buf[func() uint64 {
				p := &len_
				x := *p
				*p++
				return x
			}()] = func() byte {
				if IS_SET_AR(shop_index[shop_nr].With_who[:], uint32(int32(flag))) {
					return '_'
				}
				return *trade_letters[sindex]
			}()
			buf[len_] = '\x00'
			if len_ >= uint64(256) {
				break
			}
		}
		sindex++
		flag += 1
	}
	buf[256-1] = '\x00'
	return &buf[0]
}
func list_all_shops(ch *char_data) {
	var (
		list_all_shops_header *byte = libc.CString(" ##   Virtual   Where    Keeper    Buy   Sell   Customers\r\n---------------------------------------------------------\r\n")
		shop_nr               int
		headerlen             int    = libc.StrLen(list_all_shops_header)
		len_                  uint64 = 0
		buf                   [64936]byte
		buf1                  [16]byte
	)
	buf[0] = '\x00'
	for shop_nr = 0; shop_nr <= top_shop && len_ < uint64(64936); shop_nr++ {
		if (shop_nr % (int(PAGE_LENGTH - 2))) == 0 {
			if len_+uint64(headerlen)+1 >= uint64(64936) {
				break
			}
			libc.StrCpy(&buf[len_], list_all_shops_header)
			len_ += uint64(headerlen)
		}
		if shop_index[shop_nr].Keeper == int(-1) {
			libc.StrCpy(&buf1[0], libc.CString("<NONE>"))
		} else {
			stdio.Sprintf(&buf1[0], "%6d", mob_index[shop_index[shop_nr].Keeper].Vnum)
		}
		len_ += uint64(stdio.Snprintf(&buf[len_], int(64936-uintptr(len_)), "%3d   %6d   %6d    %s   %3.2f   %3.2f    %s\r\n", shop_nr+1, shop_index[shop_nr].Vnum, shop_index[shop_nr].In_room[0], &buf1[0], shop_index[shop_nr].Profit_sell, shop_index[shop_nr].Profit_buy, customer_string(shop_nr, 0)))
	}
	page_string(ch.Desc, &buf[0], 1)
}
func list_detailed_shop(ch *char_data, shop_nr int) {
	var (
		k       *char_data
		sindex  int
		column  int
		ptrsave *byte
	)
	send_to_char(ch, libc.CString("Vnum:       [%5d], Rnum: [%5d]\r\n"), shop_index[shop_nr].Vnum, shop_nr+1)
	send_to_char(ch, libc.CString("Rooms:      "))
	column = 12
	for sindex = 0; (shop_index[shop_nr].In_room[sindex]) != int(-1); sindex++ {
		var (
			buf1    [128]byte
			linelen int
			temp    int
		)
		if sindex != 0 {
			send_to_char(ch, libc.CString(", "))
			column += 2
		}
		if (func() int {
			temp = real_room(shop_index[shop_nr].In_room[sindex])
			return temp
		}()) != int(-1) {
			linelen = stdio.Snprintf(&buf1[0], int(128), "%s (#%d)", world[temp].Name, GET_ROOM_VNUM(temp))
		} else {
			linelen = stdio.Snprintf(&buf1[0], int(128), "<UNKNOWN> (#%d)", shop_index[shop_nr].In_room[sindex])
		}
		if linelen+column >= 78 && column >= 20 {
			send_to_char(ch, libc.CString("\r\n            "))
			column = 12
		}
		if send_to_char(ch, libc.CString("%s"), &buf1[0]) == 0 {
			return
		}
		column += linelen
	}
	if sindex == 0 {
		send_to_char(ch, libc.CString("Rooms:      None!"))
	}
	send_to_char(ch, libc.CString("\r\nShopkeeper: "))
	if shop_index[shop_nr].Keeper != int(-1) {
		send_to_char(ch, libc.CString("%s (#%d), Special Function: %s\r\n"), GET_NAME(&mob_proto[shop_index[shop_nr].Keeper]), mob_index[shop_index[shop_nr].Keeper].Vnum, func() string {
			if shop_index[shop_nr].Func != nil {
				return "YES"
			}
			return "NO"
		}())
		if (func() *char_data {
			k = get_char_num(shop_index[shop_nr].Keeper)
			return k
		}()) != nil {
			send_to_char(ch, libc.CString("Coins:      [%9d], Bank: [%9d] (Total: %d)\r\n"), k.Gold, shop_index[shop_nr].BankAccount, k.Gold+shop_index[shop_nr].BankAccount)
		}
	} else {
		send_to_char(ch, libc.CString("<NONE>\r\n"))
	}
	send_to_char(ch, libc.CString("Customers:  %s\r\n"), func() *byte {
		if (func() *byte {
			ptrsave = customer_string(shop_nr, 1)
			return ptrsave
		}()) != nil {
			return ptrsave
		}
		return libc.CString("None")
	}())
	send_to_char(ch, libc.CString("Produces:   "))
	column = 12
	for sindex = 0; (shop_index[shop_nr].Producing[sindex]) != int(-1); sindex++ {
		var (
			buf1    [128]byte
			linelen int
		)
		if sindex != 0 {
			send_to_char(ch, libc.CString(", "))
			column += 2
		}
		linelen = stdio.Snprintf(&buf1[0], int(128), "%s (#%d)", obj_proto[shop_index[shop_nr].Producing[sindex]].Short_description, obj_index[shop_index[shop_nr].Producing[sindex]].Vnum)
		if linelen+column >= 78 && column >= 20 {
			send_to_char(ch, libc.CString("\r\n            "))
			column = 12
		}
		if send_to_char(ch, libc.CString("%s"), &buf1[0]) == 0 {
			return
		}
		column += linelen
	}
	if sindex == 0 {
		send_to_char(ch, libc.CString("Produces:   Nothing!"))
	}
	send_to_char(ch, libc.CString("\r\nBuys:       "))
	column = 12
	for sindex = 0; (shop_index[shop_nr].Type[sindex]).Type != int(-1); sindex++ {
		var (
			buf1    [128]byte
			linelen uint64
		)
		if sindex != 0 {
			send_to_char(ch, libc.CString(", "))
			column += 2
		}
		linelen = uint64(stdio.Snprintf(&buf1[0], int(128), "%s (#%d) [%s]", item_types[(shop_index[shop_nr].Type[sindex]).Type], (shop_index[shop_nr].Type[sindex]).Type, func() *byte {
			if (shop_index[shop_nr].Type[sindex]).Keywords != nil {
				return (shop_index[shop_nr].Type[sindex]).Keywords
			}
			return libc.CString("all")
		}()))
		if linelen+uint64(column) >= 78 && column >= 20 {
			send_to_char(ch, libc.CString("\r\n            "))
			column = 12
		}
		if send_to_char(ch, libc.CString("%s"), &buf1[0]) == 0 {
			return
		}
		column += int(linelen)
	}
	if sindex == 0 {
		send_to_char(ch, libc.CString("Buys:       Nothing!"))
	}
	send_to_char(ch, libc.CString("\r\nBuy at:     [%4.2f], Sell at: [%4.2f], Open: [%d-%d, %d-%d]\r\n"), shop_index[shop_nr].Profit_sell, shop_index[shop_nr].Profit_buy, shop_index[shop_nr].Open1, shop_index[shop_nr].Close1, shop_index[shop_nr].Open2, shop_index[shop_nr].Close2)
	{
		var buf1 [128]byte
		// todo: fix this
		//sprintbit(shop_index[shop_nr].Bitvector, ([]*byte)(shop_bits), &buf1[0], uint64(128))
		send_to_char(ch, libc.CString("Bits:       %s\r\n"), &buf1[0])
	}
}
func show_shops(ch *char_data, arg *byte) {
	var shop_nr int
	if *arg == 0 {
		list_all_shops(ch)
	} else {
		if libc.StrCaseCmp(arg, libc.CString(".")) == 0 {
			for shop_nr = 0; shop_nr <= top_shop; shop_nr++ {
				if ok_shop_room(shop_nr, int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room)))) {
					break
				}
			}
			if shop_nr > top_shop {
				send_to_char(ch, libc.CString("This isn't a shop!\r\n"))
				return
			}
		} else if is_number(arg) {
			shop_nr = libc.Atoi(libc.GoString(arg)) - 1
		} else {
			shop_nr = -1
		}
		if shop_nr < 0 || shop_nr > top_shop {
			send_to_char(ch, libc.CString("Illegal shop number.\r\n"))
			return
		}
		list_detailed_shop(ch, shop_nr)
	}
}
func destroy_shops() {
	var (
		cnt int64
		itr int64
	)
	if shop_index == nil {
		return
	}
	for cnt = 0; cnt <= int64(top_shop); cnt++ {
		if shop_index[cnt].No_such_item1 != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].No_such_item1))
		}
		if shop_index[cnt].No_such_item2 != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].No_such_item2))
		}
		if shop_index[cnt].Missing_cash1 != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].Missing_cash1))
		}
		if shop_index[cnt].Missing_cash2 != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].Missing_cash2))
		}
		if shop_index[cnt].Do_not_buy != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].Do_not_buy))
		}
		if shop_index[cnt].Message_buy != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].Message_buy))
		}
		if shop_index[cnt].Message_sell != nil {
			libc.Free(unsafe.Pointer(shop_index[cnt].Message_sell))
		}
		if shop_index[cnt].In_room != nil {
			libc.Free(unsafe.Pointer(&shop_index[cnt].In_room[0]))
		}
		if shop_index[cnt].Producing != nil {
			libc.Free(unsafe.Pointer(&shop_index[cnt].Producing[0]))
		}
		if shop_index[cnt].Type != nil {
			for itr = 0; (shop_index[cnt].Type[itr]).Type != int(-1); itr++ {
				if (shop_index[cnt].Type[itr]).Keywords != nil {
					libc.Free(unsafe.Pointer((shop_index[cnt].Type[itr]).Keywords))
				}
			}
			libc.Free(unsafe.Pointer(&shop_index[cnt].Type[0]))
		}
	}
	libc.Free(unsafe.Pointer(&shop_index[0]))
	shop_index = nil
	top_shop = -1
}
func count_shops(low int, high int) int {
	var (
		i int
		j int
	)
	for i = func() int {
		j = 0
		return j
	}(); shop_index[i].Vnum <= high; i++ {
		if shop_index[i].Vnum >= low {
			j++
		}
	}
	return j
}
