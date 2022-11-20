package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

func one_phrase(arg *byte, first_arg *byte) *byte {
	skip_spaces(&arg)
	if *arg == 0 {
		*first_arg = '\x00'
	} else if *arg == '"' {
		var (
			p *byte
			c int8
		)
		p = matching_quote(arg)
		c = int8(*p)
		*p = '\x00'
		libc.StrCpy(first_arg, (*byte)(unsafe.Add(unsafe.Pointer(arg), 1)))
		if int(c) == '\x00' {
			return p
		} else {
			return (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
	} else {
		var (
			s *byte
			p *byte
		)
		s = first_arg
		p = arg
		for *p != 0 && !unicode.IsSpace(rune(*p)) && *p != '"' {
			*func() *byte {
				p := &s
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &p
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		}
		*s = '\x00'
		return p
	}
	return arg
}
func is_substring(sub *byte, string_ *byte) bool {
	var s *byte
	if (func() *byte {
		s = str_str(string_, sub)
		return s
	}()) != nil {
		var (
			len_   int = libc.StrLen(string_)
			sublen int = libc.StrLen(sub)
		)
		if (s == string_ || unicode.IsSpace(rune(*((*byte)(unsafe.Add(unsafe.Pointer(s), -1))))) || unicode.IsPunct(rune(*((*byte)(unsafe.Add(unsafe.Pointer(s), -1)))))) && ((*byte)(unsafe.Add(unsafe.Pointer(s), sublen)) == (*byte)(unsafe.Add(unsafe.Pointer(string_), len_)) || unicode.IsSpace(rune(*(*byte)(unsafe.Add(unsafe.Pointer(s), sublen)))) || unicode.IsPunct(rune(*(*byte)(unsafe.Add(unsafe.Pointer(s), sublen))))) {
			return true
		}
	}
	return false
}
func word_check(str *byte, wordlist *byte) bool {
	var (
		words  [2048]byte
		phrase [2048]byte
		s      *byte
	)
	if *wordlist == '*' {
		return true
	}
	libc.StrCpy(&words[0], wordlist)
	for s = one_phrase(&words[0], &phrase[0]); phrase[0] != 0; s = one_phrase(s, &phrase[0]) {
		if is_substring(&phrase[0], str) {
			return true
		}
	}
	return false
}
func random_mtrigger(ch *char_data) {
	var t *trig_data
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<1) || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<1) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func bribe_mtrigger(ch *char_data, actor *char_data, amount int) {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<12) || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<12) && t.Depth == 0 && amount >= t.Narg {
			stdio.Snprintf(&buf[0], int(2048), "%d", amount)
			add_var(&t.Var_list, libc.CString("amount"), &buf[0], 0)
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func greet_memory_mtrigger(actor *char_data) {
	var (
		t                 *trig_data
		ch                *char_data
		mem               *script_memory
		buf               [2048]byte
		command_performed int = 0
	)
	if !valid_dg_target(actor, 1<<0) {
		return
	}
	for ch = world[actor.In_room].People; ch != nil; ch = ch.Next_in_room {
		if ch.Memory == nil || !AWAKE(ch) || ch.Fighting != nil || ch == actor || AFF_FLAGGED(ch, AFF_CHARM) {
			continue
		}
		for mem = ch.Memory; mem != nil && ch.Memory != nil; mem = mem.Next {
			if int(actor.Id) != mem.Id {
				continue
			}
			if mem.Cmd != nil {
				command_interpreter(ch, mem.Cmd)
				command_performed = 1
				break
			}
			if mem != nil && command_performed == 0 {
				for t = ch.Script.Trig_list; t != nil; t = t.Next {
					if IS_SET(uint32(int32(t.Trigger_type)), 1<<14) && CAN_SEE(ch, actor) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
						for {
							stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
							add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
							if true {
								break
							}
						}
						script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
						break
					}
				}
			}
			if mem != nil {
				if ch.Memory == mem {
					ch.Memory = mem.Next
				} else {
					var prev *script_memory
					prev = ch.Memory
					for prev.Next != mem {
						prev = prev.Next
					}
					prev.Next = mem.Next
				}
				if mem.Cmd != nil {
					libc.Free(unsafe.Pointer(mem.Cmd))
				}
				libc.Free(unsafe.Pointer(mem))
			}
		}
	}
}
func greet_mtrigger(actor *char_data, dir int) bool {
	var (
		t            *trig_data
		ch           *char_data
		buf          [2048]byte
		intermediate int
		final        int = 1
	)
	if !valid_dg_target(actor, 1<<0) {
		return true
	}
	for ch = world[actor.In_room].People; ch != nil; ch = ch.Next_in_room {
		if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), (1<<6)|1<<7) || !AWAKE(ch) || ch.Fighting != nil || ch == actor || AFF_FLAGGED(ch, AFF_CHARM) {
			continue
		}
		if (ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<7)) && IS_NPC(actor) {
			continue
		}
		for t = ch.Script.Trig_list; t != nil; t = t.Next {
			if (IS_SET(uint32(int32(t.Trigger_type)), 1<<6) && CAN_SEE(ch, actor) || IS_SET(uint32(int32(t.Trigger_type)), 1<<7)) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
				if dir >= 0 && dir < NUM_OF_DIRS {
					add_var(&t.Var_list, libc.CString("direction"), dirs[rev_dir[dir]], 0)
				} else {
					add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
				}
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
				intermediate = script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
				if intermediate == 0 {
					final = 0
				}
				continue
			}
		}
	}
	return final != 0
}
func entry_memory_mtrigger(ch *char_data) {
	var (
		t     *trig_data
		actor *char_data
		mem   *script_memory
		buf   [2048]byte
	)
	if ch.Memory == nil || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for actor = world[ch.In_room].People; actor != nil && ch.Memory != nil; actor = actor.Next_in_room {
		if actor != ch && ch.Memory != nil {
			for mem = ch.Memory; mem != nil && ch.Memory != nil; mem = mem.Next {
				if int(actor.Id) == mem.Id {
					var prev *script_memory
					if mem.Cmd != nil {
						command_interpreter(ch, mem.Cmd)
					} else {
						for t = ch.Script.Trig_list; t != nil; t = t.Next {
							if IS_SET(uint32(int32(t.Trigger_type)), 1<<14) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
								for {
									stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
									add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
									if true {
										break
									}
								}
								script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
								break
							}
						}
					}
					if ch.Memory == mem {
						ch.Memory = mem.Next
					} else {
						prev = ch.Memory
						for prev.Next != mem {
							prev = prev.Next
						}
						prev.Next = mem.Next
					}
					if mem.Cmd != nil {
						libc.Free(unsafe.Pointer(mem.Cmd))
					}
					libc.Free(unsafe.Pointer(mem))
				}
			}
		}
	}
}
func entry_mtrigger(ch *char_data) bool {
	var t *trig_data
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<8) || AFF_FLAGGED(ch, AFF_CHARM) {
		return true
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<8) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			return script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0
			break
		}
	}
	return true
}
func command_mtrigger(actor *char_data, cmd *byte, argument *byte) bool {
	var (
		ch      *char_data
		ch_next *char_data
		t       *trig_data
		buf     [2048]byte
	)
	if !valid_dg_target(actor, 0) {
		return false
	}
	for ch = world[actor.In_room].People; ch != nil; ch = ch_next {
		ch_next = ch.Next_in_room
		if ch.Script != nil && IS_SET(uint32(int32(ch.Script.Types)), 1<<2) && !AFF_FLAGGED(ch, AFF_CHARM) && actor != ch {
			for t = ch.Script.Trig_list; t != nil; t = t.Next {
				if !IS_SET(uint32(int32(t.Trigger_type)), 1<<2) || t.Depth != 0 {
					continue
				}
				if t.Arglist == nil || *t.Arglist == 0 {
					mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: Command Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
					continue
				}
				if *t.Arglist == '*' || libc.StrNCaseCmp(t.Arglist, cmd, libc.StrLen(t.Arglist)) == 0 {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
						add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
						if true {
							break
						}
					}
					skip_spaces(&argument)
					add_var(&t.Var_list, libc.CString("arg"), argument, 0)
					skip_spaces(&cmd)
					add_var(&t.Var_list, libc.CString("cmd"), cmd, 0)
					if script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0 {
						return true
					}
				}
			}
		}
	}
	return false
}
func speech_mtrigger(actor *char_data, str *byte) {
	var (
		ch      *char_data
		ch_next *char_data
		t       *trig_data
		buf     [2048]byte
	)
	for ch = world[actor.In_room].People; ch != nil; ch = ch_next {
		ch_next = ch.Next_in_room
		if ch.Script != nil && IS_SET(uint32(int32(ch.Script.Types)), 1<<3) && AWAKE(ch) && !AFF_FLAGGED(ch, AFF_CHARM) && actor != ch {
			for t = ch.Script.Trig_list; t != nil; t = t.Next {
				if !IS_SET(uint32(int32(t.Trigger_type)), 1<<3) || t.Depth != 0 {
					continue
				}
				if t.Arglist == nil || *t.Arglist == 0 {
					mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: Speech Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
					continue
				}
				if t.Narg != 0 && word_check(str, t.Arglist) || t.Narg == 0 && is_substring(t.Arglist, str) {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
						add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
						if true {
							break
						}
					}
					add_var(&t.Var_list, libc.CString("speech"), str, 0)
					script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
					break
				}
			}
		}
	}
}
func act_mtrigger(ch *char_data, str *byte, actor *char_data, victim *char_data, object *obj_data, target *obj_data, arg *byte) {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if ch.Script != nil && IS_SET(uint32(int32(ch.Script.Types)), 1<<4) && !AFF_FLAGGED(ch, AFF_CHARM) && actor != ch {
		for t = ch.Script.Trig_list; t != nil; t = t.Next {
			if !IS_SET(uint32(int32(t.Trigger_type)), 1<<4) || t.Depth != 0 {
				continue
			}
			if t.Arglist == nil || *t.Arglist == 0 {
				mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: Act Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
				continue
			}
			if t.Narg != 0 && word_check(str, t.Arglist) || t.Narg == 0 && is_substring(t.Arglist, str) {
				if actor != nil {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
						add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
						if true {
							break
						}
					}
				}
				if victim != nil {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, victim.Id)
						add_var(&t.Var_list, libc.CString("victim"), &buf[0], 0)
						if true {
							break
						}
					}
				}
				if object != nil {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, object.Id)
						add_var(&t.Var_list, libc.CString("object"), &buf[0], 0)
						if true {
							break
						}
					}
				}
				if target != nil {
					for {
						stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, target.Id)
						add_var(&t.Var_list, libc.CString("target"), &buf[0], 0)
						if true {
							break
						}
					}
				}
				if str != nil {
					var (
						nstr *byte = libc.StrDup(str)
						fstr *byte = nstr
						p    *byte = libc.StrChr(nstr, '\r')
					)
					_ = p
					skip_spaces(&nstr)
					*p = '\x00'
					add_var(&t.Var_list, libc.CString("arg"), nstr, 0)
					libc.Free(unsafe.Pointer(fstr))
				}
				script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
				break
			}
		}
	}
}
func fight_mtrigger(ch *char_data) {
	var (
		actor *char_data
		t     *trig_data
		buf   [2048]byte
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<10) || ch.Fighting == nil || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<10) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			actor = ch.Fighting
			if actor != nil {
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
			} else {
				add_var(&t.Var_list, libc.CString("actor"), libc.CString("nobody"), 0)
			}
			script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func hitprcnt_mtrigger(ch *char_data) {
	var (
		actor *char_data
		t     *trig_data
		buf   [2048]byte
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<11) || ch.Fighting == nil || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<11) && t.Depth == 0 && ch.Max_hit != 0 && ch.Hit <= (ch.Max_hit/100)*int64(t.Narg) {
			actor = ch.Fighting
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func receive_mtrigger(ch *char_data, actor *char_data, obj *obj_data) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<9) || AFF_FLAGGED(ch, AFF_CHARM) {
		return true
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<9) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, obj.Id)
				add_var(&t.Var_list, libc.CString("object"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			if PLR_FLAGGED(actor, PLR_NOTDEADYET) || MOB_FLAGGED(actor, MOB_NOTDEADYET) || (PLR_FLAGGED(ch, PLR_NOTDEADYET) || MOB_FLAGGED(ch, MOB_NOTDEADYET)) || obj.Carried_by != actor {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func death_mtrigger(ch *char_data, actor *char_data) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<5) || AFF_FLAGGED(ch, AFF_CHARM) {
		return true
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<5) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			if actor != nil {
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
			}
			return script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func load_mtrigger(ch *char_data) {
	var (
		t      *trig_data
		result int = 0
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<13) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<13) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			result = script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
	if result == int(-9999999) {
		if ch.Nr != int(-1) {
			free_proto_script(unsafe.Pointer(&mob_proto[ch.Nr]), MOB_TRIGGER)
		}
	}
}
func cast_mtrigger(actor *char_data, ch *char_data, spellnum int) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if ch == nil {
		return true
	}
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<15) || AFF_FLAGGED(ch, AFF_CHARM) {
		return true
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<15) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			stdio.Sprintf(&buf[0], "%d", spellnum)
			add_var(&t.Var_list, libc.CString("spell"), &buf[0], 0)
			add_var(&t.Var_list, libc.CString("spellname"), skill_name(spellnum), 0)
			return script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func leave_mtrigger(actor *char_data, dir int) bool {
	var (
		t   *trig_data
		ch  *char_data
		buf [2048]byte
	)
	if !valid_dg_target(actor, 1<<0) {
		return true
	}
	for ch = world[actor.In_room].People; ch != nil; ch = ch.Next_in_room {
		if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<16) || !AWAKE(ch) || ch.Fighting != nil || ch == actor || AFF_FLAGGED(ch, AFF_CHARM) {
			continue
		}
		for t = ch.Script.Trig_list; t != nil; t = t.Next {
			if IS_SET(uint32(int32(t.Trigger_type)), 1<<16) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
				if dir >= 0 && dir < NUM_OF_DIRS {
					add_var(&t.Var_list, libc.CString("direction"), dirs[dir], 0)
				} else {
					add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
				}
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
				return script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0
			}
		}
	}
	return true
}
func door_mtrigger(actor *char_data, subcmd int, dir int) bool {
	var (
		t   *trig_data
		ch  *char_data
		buf [2048]byte
	)
	for ch = world[actor.In_room].People; ch != nil; ch = ch.Next_in_room {
		if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<17) || !AWAKE(ch) || ch.Fighting != nil || ch == actor || AFF_FLAGGED(ch, AFF_CHARM) {
			continue
		}
		for t = ch.Script.Trig_list; t != nil; t = t.Next {
			if IS_SET(uint32(int32(t.Trigger_type)), 1<<17) && CAN_SEE(ch, actor) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
				add_var(&t.Var_list, libc.CString("cmd"), cmd_door[subcmd], 0)
				if dir >= 0 && dir < NUM_OF_DIRS {
					add_var(&t.Var_list, libc.CString("direction"), dirs[dir], 0)
				} else {
					add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
				}
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
				return script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW) != 0
			}
		}
	}
	return true
}
func time_mtrigger(ch *char_data) {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if ch.Script == nil || !IS_SET(uint32(int32(ch.Script.Types)), 1<<19) || AFF_FLAGGED(ch, AFF_CHARM) {
		return
	}
	for t = ch.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<19) && t.Depth == 0 && time_info.Hours == t.Narg {
			stdio.Sprintf(&buf[0], "%d", time_info.Hours)
			add_var(&t.Var_list, libc.CString("time"), &buf[0], 0)
			script_driver(unsafe.Pointer(&ch), t, MOB_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func random_otrigger(obj *obj_data) {
	var t *trig_data
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<1) {
		return
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<1) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func timer_otrigger(obj *obj_data) {
	var t *trig_data
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<5) {
		return
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<5) && t.Depth == 0 {
			script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
		}
	}
	return
}
func get_otrigger(obj *obj_data, actor *char_data) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<6) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<6) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if PLR_FLAGGED(actor, PLR_NOTDEADYET) || MOB_FLAGGED(actor, MOB_NOTDEADYET) || obj == nil {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func cmd_otrig(obj *obj_data, actor *char_data, cmd *byte, argument *byte, type_ int) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if obj != nil && (obj.Script != nil && IS_SET(uint32(int32(obj.Script.Types)), 1<<2)) {
		for t = obj.Script.Trig_list; t != nil; t = t.Next {
			if !IS_SET(uint32(int32(t.Trigger_type)), 1<<2) || t.Depth != 0 {
				continue
			}
			if IS_SET(uint32(int32(t.Narg)), uint32(int32(type_))) && (t.Arglist == nil || *t.Arglist == 0) {
				mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: O-Command Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
				continue
			}
			if IS_SET(uint32(int32(t.Narg)), uint32(int32(type_))) && (*t.Arglist == '*' || libc.StrNCaseCmp(t.Arglist, cmd, libc.StrLen(t.Arglist)) == 0) {
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
				skip_spaces(&argument)
				add_var(&t.Var_list, libc.CString("arg"), argument, 0)
				skip_spaces(&cmd)
				add_var(&t.Var_list, libc.CString("cmd"), cmd, 0)
				if script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW) != 0 {
					return true
				}
			}
		}
	}
	return false
}
func command_otrigger(actor *char_data, cmd *byte, argument *byte) bool {
	var (
		obj *obj_data
		i   int
	)
	if !valid_dg_target(actor, 0) {
		return false
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (actor.Equipment[i]) != nil {
			if cmd_otrig(actor.Equipment[i], actor, cmd, argument, 1<<0) && !OBJ_FLAGGED(actor.Equipment[i], ITEM_FORGED) {
				return true
			}
		}
	}
	for obj = actor.Carrying; obj != nil; obj = obj.Next_content {
		if cmd_otrig(obj, actor, cmd, argument, 1<<1) && !OBJ_FLAGGED(obj, ITEM_FORGED) {
			return true
		}
	}
	for obj = world[actor.In_room].Contents; obj != nil; obj = obj.Next_content {
		if cmd_otrig(obj, actor, cmd, argument, 1<<2) && !OBJ_FLAGGED(obj, ITEM_FORGED) {
			return true
		}
	}
	return false
}
func wear_otrigger(obj *obj_data, actor *char_data, where int) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<9) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<9) && t.Depth == 0 {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if obj == nil {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func remove_otrigger(obj *obj_data, actor *char_data) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<11) {
		return true
	}
	if !valid_dg_target(actor, 0) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<11) && t.Depth == 0 {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if obj == nil {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func drop_otrigger(obj *obj_data, actor *char_data) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<7) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<7) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if obj == nil {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func give_otrigger(obj *obj_data, actor *char_data, victim *char_data) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<8) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<8) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, victim.Id)
				add_var(&t.Var_list, libc.CString("victim"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if obj == nil || obj.Carried_by != actor {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func load_otrigger(obj *obj_data) {
	var (
		t      *trig_data
		result int = 0
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<13) {
		return
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<13) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			result = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			break
		}
	}
	if result == int(-9999999) {
		if obj.Item_number != int(-1) {
			free_proto_script(unsafe.Pointer(&obj_proto[obj.Item_number]), OBJ_TRIGGER)
		}
	}
}
func cast_otrigger(actor *char_data, obj *obj_data, spellnum int) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if obj == nil {
		return true
	}
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<15) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<15) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			stdio.Sprintf(&buf[0], "%d", spellnum)
			add_var(&t.Var_list, libc.CString("spell"), &buf[0], 0)
			add_var(&t.Var_list, libc.CString("spellname"), skill_name(spellnum), 0)
			return script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func leave_otrigger(room *room_data, actor *char_data, dir int) bool {
	var (
		t        *trig_data
		buf      [2048]byte
		temp     int
		final    int = 1
		obj      *obj_data
		obj_next *obj_data
	)
	if !valid_dg_target(actor, 1<<0) {
		return true
	}
	for obj = room.Contents; obj != nil; obj = obj_next {
		obj_next = obj.Next_content
		if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<16) {
			continue
		}
		for t = obj.Script.Trig_list; t != nil; t = t.Next {
			if IS_SET(uint32(int32(t.Trigger_type)), 1<<16) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
				if dir >= 0 && dir < NUM_OF_DIRS {
					add_var(&t.Var_list, libc.CString("direction"), dirs[dir], 0)
				} else {
					add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
				}
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
					add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
					if true {
						break
					}
				}
				temp = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
				if temp == 0 {
					final = 0
				}
			}
		}
	}
	return final != 0
}
func consume_otrigger(obj *obj_data, actor *char_data, cmd int) bool {
	var (
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<18) {
		return true
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<18) && t.Depth == 0 {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			switch cmd {
			case OCMD_EAT:
				add_var(&t.Var_list, libc.CString("command"), libc.CString("eat"), 0)
			case OCMD_DRINK:
				add_var(&t.Var_list, libc.CString("command"), libc.CString("drink"), 0)
			case OCMD_QUAFF:
				add_var(&t.Var_list, libc.CString("command"), libc.CString("quaff"), 0)
			}
			ret_val = script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			if obj == nil {
				return false
			} else {
				return ret_val != 0
			}
		}
	}
	return true
}
func time_otrigger(obj *obj_data) {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if obj.Script == nil || !IS_SET(uint32(int32(obj.Script.Types)), 1<<19) {
		return
	}
	for t = obj.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<19) && t.Depth == 0 && time_info.Hours == t.Narg {
			stdio.Sprintf(&buf[0], "%d", time_info.Hours)
			add_var(&t.Var_list, libc.CString("time"), &buf[0], 0)
			script_driver(unsafe.Pointer(&obj), t, OBJ_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func reset_wtrigger(room *room_data) {
	var t *trig_data
	if room.Script == nil || !IS_SET(uint32(int32(room.Script.Types)), 1<<5) {
		return
	}
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<5) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func random_wtrigger(room *room_data) {
	var t *trig_data
	if room.Script == nil || !IS_SET(uint32(int32(room.Script.Types)), 1<<1) {
		return
	}
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<1) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func enter_wtrigger(room *room_data, actor *char_data, dir int) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if room.Script == nil || !IS_SET(uint32(int32(room.Script.Types)), 1<<6) {
		return true
	}
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<6) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			if dir >= 0 && dir < NUM_OF_DIRS {
				add_var(&t.Var_list, libc.CString("direction"), dirs[rev_dir[dir]], 0)
			} else {
				add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			return script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func command_wtrigger(actor *char_data, cmd *byte, argument *byte) bool {
	var (
		room *room_data
		t    *trig_data
		buf  [2048]byte
	)
	if actor == nil || (world[actor.In_room].Script == nil || !IS_SET(uint32(int32(world[actor.In_room].Script.Types)), 1<<2)) {
		return false
	}
	if !valid_dg_target(actor, 0) {
		return false
	}
	room = &world[actor.In_room]
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if !IS_SET(uint32(int32(t.Trigger_type)), 1<<2) || t.Depth != 0 {
			continue
		}
		if t.Arglist == nil || *t.Arglist == 0 {
			mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: W-Command Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
			continue
		}
		if *t.Arglist == '*' || libc.StrNCaseCmp(t.Arglist, cmd, libc.StrLen(t.Arglist)) == 0 {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			skip_spaces(&argument)
			add_var(&t.Var_list, libc.CString("arg"), argument, 0)
			skip_spaces(&cmd)
			add_var(&t.Var_list, libc.CString("cmd"), cmd, 0)
			return script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW) != 0
		}
	}
	return false
}
func speech_wtrigger(actor *char_data, str *byte) {
	var (
		room *room_data
		t    *trig_data
		buf  [2048]byte
	)
	if actor == nil || (world[actor.In_room].Script == nil || !IS_SET(uint32(int32(world[actor.In_room].Script.Types)), 1<<3)) {
		return
	}
	room = &world[actor.In_room]
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if !IS_SET(uint32(int32(t.Trigger_type)), 1<<3) || t.Depth != 0 {
			continue
		}
		if t.Arglist == nil || *t.Arglist == 0 {
			mudlog(NRM, ADMLVL_BUILDER, 1, libc.CString("SYSERR: W-Speech Trigger #%d has no text argument!"), trig_index[t.Nr].Vnum)
			continue
		}
		if *t.Arglist == '*' || t.Narg != 0 && word_check(str, t.Arglist) || t.Narg == 0 && is_substring(t.Arglist, str) {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			add_var(&t.Var_list, libc.CString("speech"), str, 0)
			script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW)
			break
		}
	}
}
func drop_wtrigger(obj *obj_data, actor *char_data) bool {
	var (
		room    *room_data
		t       *trig_data
		buf     [2048]byte
		ret_val int
	)
	if actor == nil || (world[actor.In_room].Script == nil || !IS_SET(uint32(int32(world[actor.In_room].Script.Types)), 1<<7)) {
		return true
	}
	room = &world[actor.In_room]
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<7) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, obj.Id)
				add_var(&t.Var_list, libc.CString("object"), &buf[0], 0)
				if true {
					break
				}
			}
			ret_val = script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW)
			if obj.Carried_by != actor {
				return false
			} else {
				return ret_val != 0
			}
			break
		}
	}
	return true
}
func cast_wtrigger(actor *char_data, vict *char_data, obj *obj_data, spellnum int) bool {
	var (
		room *room_data
		t    *trig_data
		buf  [2048]byte
	)
	if actor == nil || (world[actor.In_room].Script == nil || !IS_SET(uint32(int32(world[actor.In_room].Script.Types)), 1<<15)) {
		return true
	}
	room = &world[actor.In_room]
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<15) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			if vict != nil {
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, vict.Id)
					add_var(&t.Var_list, libc.CString("victim"), &buf[0], 0)
					if true {
						break
					}
				}
			}
			if obj != nil {
				for {
					stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, obj.Id)
					add_var(&t.Var_list, libc.CString("object"), &buf[0], 0)
					if true {
						break
					}
				}
			}
			stdio.Sprintf(&buf[0], "%d", spellnum)
			add_var(&t.Var_list, libc.CString("spell"), &buf[0], 0)
			add_var(&t.Var_list, libc.CString("spellname"), skill_name(spellnum), 0)
			return script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func leave_wtrigger(room *room_data, actor *char_data, dir int) bool {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if !valid_dg_target(actor, 1<<0) {
		return true
	}
	if room.Script == nil || !IS_SET(uint32(int32(room.Script.Types)), 1<<16) {
		return true
	}
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<16) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			if dir >= 0 && dir < NUM_OF_DIRS {
				add_var(&t.Var_list, libc.CString("direction"), dirs[dir], 0)
			} else {
				add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			return script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func door_wtrigger(actor *char_data, subcmd int, dir int) bool {
	var (
		room *room_data
		t    *trig_data
		buf  [2048]byte
	)
	if actor == nil || (world[actor.In_room].Script == nil || !IS_SET(uint32(int32(world[actor.In_room].Script.Types)), 1<<17)) {
		return true
	}
	room = &world[actor.In_room]
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<17) && t.Depth == 0 && rand_number(1, 100) <= t.Narg {
			add_var(&t.Var_list, libc.CString("cmd"), cmd_door[subcmd], 0)
			if dir >= 0 && dir < NUM_OF_DIRS {
				add_var(&t.Var_list, libc.CString("direction"), dirs[dir], 0)
			} else {
				add_var(&t.Var_list, libc.CString("direction"), libc.CString("none"), 0)
			}
			for {
				stdio.Sprintf(&buf[0], "%c%d", UID_CHAR, actor.Id)
				add_var(&t.Var_list, libc.CString("actor"), &buf[0], 0)
				if true {
					break
				}
			}
			return script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW) != 0
		}
	}
	return true
}
func time_wtrigger(room *room_data) {
	var (
		t   *trig_data
		buf [2048]byte
	)
	if room.Script == nil || !IS_SET(uint32(int32(room.Script.Types)), 1<<19) {
		return
	}
	for t = room.Script.Trig_list; t != nil; t = t.Next {
		if IS_SET(uint32(int32(t.Trigger_type)), 1<<19) && t.Depth == 0 && time_info.Hours == t.Narg {
			stdio.Sprintf(&buf[0], "%d", time_info.Hours)
			add_var(&t.Var_list, libc.CString("time"), &buf[0], 0)
			script_driver(unsafe.Pointer(&room), t, WLD_TRIGGER, TRIG_NEW)
			break
		}
	}
}
