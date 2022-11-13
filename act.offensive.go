package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func do_galikgun(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.15
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_GALIKGUN) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_GALIKGUN]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_GALIKGUN)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_GALIKGUN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_GALIKGUN]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		} else if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		} else if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Galik Gun before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Galik Gun before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Galik Gun before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your Galik Gun!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Galik Gun!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Galik Gun!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 16, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Galik Gun, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Galik Gun, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Galik Gun, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 16, skill, SKILL_GALIKGUN)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Galik Gun misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Galik Gun at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Galik Gun at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Galik Gun misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Galik Gun at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Galik Gun at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 16, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou take your charged ki and form a sparkling purple shroud of energy around your body! You swing your arms towards @c$N@W with your palms flatly facing $M and shout '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around your body forms into a beam and crashes into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards you with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards @c$N@W with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou take your charged ki and form a sparkling purple shroud of energy around your body! You swing your arms towards @c$N@W with your palms flatly facing $M and shout '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around your body forms into a beam and crashes into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards you with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards @c$N@W with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou take your charged ki and form a sparkling purple shroud of energy around your body! You swing your arms towards @c$N@W with your palms flatly facing $M and shout '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around your body forms into a beam and crashes into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards you with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards @c$N@W with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou take your charged ki and form a sparkling purple shroud of energy around your body! You swing your arms towards @c$N@W with your palms flatly facing $M and shout '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around your body forms into a beam and crashes into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards you with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards @c$N@W with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou take your charged ki and form a sparkling purple shroud of energy around your body! You swing your arms towards @c$N@W with your palms flatly facing $M and shout '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around your body forms into a beam and crashes into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards you with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W takes $s charged ki and forms a sparkling purple shroud of energy around $mself! $e swings $s arms towards @c$N@W with palms facing out flatly and shouts '@mG@Mal@wik @mG@Mu@wn@W!' as the energy around $s body forms into a beam and crashes into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if int(ch.Skillperfs[SKILL_GALIKGUN]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 16, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a Galik Gun at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Galik Gun at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_honoo(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_HONOO) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_HONOO]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_HONOO]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_HONOO)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_HONOO, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_HONOO]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Honoo before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Honoo before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Honoo before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
					send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
				} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
					send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
				}
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your honoo!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W honoo!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W honoo!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 21, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
					} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
					}
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your honoo, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W honoo, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W honoo, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 21, skill, SKILL_HONOO)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
					} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
					}
					return
				} else {
					act(libc.CString("@WYou can't believe it but your honoo misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a honoo at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a honoo at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
					} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
						send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
					}
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your honoo misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a honoo at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a honoo at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
				send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
				send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 21, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			if check_ruby(ch) == 1 {
				dmg += int64(float64(dmg) * 0.2)
			}
			if (vict.Bonuses[BONUS_FIREPROOF]) != 0 {
				dmg -= int64(float64(dmg) * 0.4)
			} else if (vict.Bonuses[BONUS_FIREPRONE]) != 0 {
				dmg += int64(float64(dmg) * 0.4)
			}
			vict.Affected_by[int(AFF_ASHED/32)] |= 1 << (int(AFF_ASHED % 32))
			switch hitspot {
			case 1:
				act(libc.CString("@WYou gather your charged ki and bring it up into your throat while mixing it with the air in your lungs. You grin evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from your lips! @c$N@W's body is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at you before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! YOUR body is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! @c$N@W's body is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou gather your charged ki and bring it up into your throat while mixing it with the air in your lungs. You grin evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from your lips! @c$N@W's face is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at you before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! YOUR face is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! @c$N@W's face is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou gather your charged ki and bring it up into your throat while mixing it with the air in your lungs. You grin evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from your lips! @c$N@W's gut is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at you before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! YOUR gut is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! @c$N@W's gut is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou gather your charged ki and bring it up into your throat while mixing it with the air in your lungs. You grin evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from your lips! @c$N@W's arm is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at you before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! YOUR arm is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! @c$N@W's arm is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou gather your charged ki and bring it up into your throat while mixing it with the air in your lungs. You grin evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from your lips! @c$N@W's leg is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at you before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! YOUR leg is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki and brings it up into $s throat while mixing it with the air in $s lungs. $e grins evily at @c$N@W before unleashing a massive jet of @rf@Rl@Ya@rm@Re@Ys@W from $s lips! @c$N@W's leg is engulfed!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_BURNED) && rand_number(1, 4) == 3 && int(vict.Race) != RACE_DEMON && (vict.Bonuses[BONUS_FIREPROOF]) == 0 {
				send_to_char(vict, libc.CString("@RYou are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are burned by the attack!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			} else if (vict.Bonuses[BONUS_FIREPROOF]) != 0 || int(vict.Race) == RACE_DEMON {
				send_to_char(ch, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			} else if (vict.Bonuses[BONUS_FIREPRONE]) != 0 {
				send_to_char(vict, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are easily burned!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			}
			if int(ch.Skillperfs[SKILL_HONOO]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect < -1 {
				send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates some!\r\n"))
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect += 1
			} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == -1 {
				send_to_room(ch.In_room, libc.CString("The water surrounding the area evaporates completely away!\r\n"))
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 0
			}
			vict.Affected_by[int(AFF_ASHED/32)] &= ^(1 << (int(AFF_ASHED % 32)))
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 21, skill, attperc)
		dmg /= 10
		if GET_OBJ_VNUM(obj) == 79 {
			dmg *= 3
		}
		act(libc.CString("@WYou fire a honoo at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a honoo at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_psyblast(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_PSYBLAST) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_PSYBLAST]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_PSYBLAST)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_PSYBLAST, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_PSYBLAST]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Psychic Blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Psychic Blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Psychic Blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your psychic blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W psychic blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W psychic blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 20, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your psychic blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W psychic blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W psychic blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 20, skill, SKILL_PSYBLAST)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your psychic blast misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a psychic blast at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a psychic blast at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your psychic blast misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a psychic blast at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a psychic blast at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 20, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou gather your charged ki into your brain as a flash of @bb@Bl@wue@W light shoots from your forehead and slams into @c$N@W's body! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into YOUR body! You scream for a moment as terrifying images sear through your mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into @c$N@W's body! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou gather your charged ki into your brain as a flash of @bb@Bl@wue@W light shoots from your forehead and slams into @c$N@W's head! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into YOUR head! You scream for a moment as terrifying images sear through your mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into @c$N@W's head! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou gather your charged ki into your brain as a flash of @bb@Bl@wue@W light shoots from your forehead and slams into @c$N@W's gut! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into YOUR gut! You scream for a moment as terrifying images sear through your mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into @c$N@W's gut! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou gather your charged ki into your brain as a flash of @bb@Bl@wue@W light shoots from your forehead and slams into @c$N@W's arm! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into YOUR arm! You scream for a moment as terrifying images sear through your mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into @c$N@W's arm! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou gather your charged ki into your brain as a flash of @bb@Bl@wue@W light shoots from your forehead and slams into @c$N@W's leg! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into YOUR leg! You scream for a moment as terrifying images sear through your mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s brain as a flash of @bb@Bl@wue@W light shoots from $s forehead and slams into @c$N@W's leg! $E screams for a moment as terrifying images sear through $S mind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if vict.Charge > 0 && rand_number(1, 3) == 2 {
				vict.Charge -= dmg / 5
				if vict.Charge < 0 {
					vict.Charge = 0
				}
				send_to_char(vict, libc.CString("@RYou lose some of your charged ki!@n\r\n"))
			}
			if !AFF_FLAGGED(vict, AFF_SHOCKED) && rand_number(1, 4) == 4 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
				act(libc.CString("@MYour mind has been shocked!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@M$n@m's mind has been shocked!@n"), TRUE, vict, nil, nil, TO_ROOM)
				vict.Affected_by[int(AFF_SHOCKED/32)] |= 1 << (int(AFF_SHOCKED % 32))
			}
			if int(ch.Skillperfs[SKILL_PSYBLAST]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 20, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a psychic blast at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a psychic blast at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_tslash(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_TSLASH) == 0 {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	} else if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 && (ch.Limb_condition[1]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right arm has damaged it more!@n\r\n"))
		ch.Limb_condition[0] -= rand_number(3, 5)
		if (ch.Limb_condition[0]) < 0 {
			act(libc.CString("@RYour right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 && (ch.Limb_condition[0]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left arm has damaged it more!@n\r\n"))
		ch.Limb_condition[1] -= rand_number(3, 5)
		if (ch.Limb_condition[1]) < 0 {
			act(libc.CString("@RYour left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) == nil {
		send_to_char(ch, libc.CString("You need to wield a sword to use this.\r\n"))
		return
	}
	if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) != int(TYPE_SLASH-TYPE_HIT) {
		send_to_char(ch, libc.CString("You are not wielding a sword, you need one to use this technique.\r\n"))
		return
	}
	var wobj *obj_data = (ch.Equipment[WEAR_WIELD1])
	var wlvl int = 0
	if OBJ_FLAGGED(wobj, ITEM_WEAPLVL1) {
		wlvl = 1
	} else if OBJ_FLAGGED(wobj, ITEM_WEAPLVL2) {
		wlvl = 2
	} else if OBJ_FLAGGED(wobj, ITEM_WEAPLVL3) {
		wlvl = 3
	} else if OBJ_FLAGGED(wobj, ITEM_WEAPLVL4) {
		wlvl = 4
	} else if OBJ_FLAGGED(wobj, ITEM_WEAPLVL5) {
		wlvl = 5
	}
	if int(ch.Skillperfs[SKILL_TSLASH]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_TSLASH]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_TSLASH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_TSLASH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_TSLASH]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Twin Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Twin Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Twin Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your twin slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W twin slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W twin slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 19, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your twin slash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W twin slash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W twin slash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 19, skill, SKILL_TSLASH)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your twin slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a twin slash at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a twin slash at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your twin slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a twin slash at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a twin slash at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 19, skill, attperc)
			if GET_SKILL(ch, SKILL_TSLASH) >= 100 {
				dmg += int64((float64(dmg) * 0.05) * float64(wlvl))
			} else if GET_SKILL(ch, SKILL_TSLASH) >= 60 {
				dmg += int64((float64(dmg) * 0.02) * float64(wlvl))
			} else if GET_SKILL(ch, SKILL_TSLASH) >= 40 {
				dmg += int64((float64(dmg) * 0.01) * float64(wlvl))
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as you draw it up to attack. Two blindingly quick slashes connect with @c$N@W's body as you fly past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with YOUR body as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with @c$N@W's body as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
				if (int(vict.Race) == RACE_ICER || int(vict.Race) == RACE_BIO) && PLR_FLAGGED(vict, PLR_TAIL) {
					act(libc.CString("@rYou cut off $S tail!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYour tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r's tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Act[int(PLR_TAIL/32)] &= bitvector_t(int32(^(1 << (int(PLR_TAIL % 32)))))
					remove_limb(vict, 6)
				}
				if (int(vict.Race) == RACE_SAIYAN || int(vict.Race) == RACE_HALFBREED) && PLR_FLAGGED(vict, PLR_STAIL) {
					act(libc.CString("@rYou cut off $S tail!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYour tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r's tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Act[int(PLR_STAIL/32)] &= bitvector_t(int32(^(1 << (int(PLR_STAIL % 32)))))
					remove_limb(vict, 5)
					if PLR_FLAGGED(vict, PLR_OOZARU) {
						act(libc.CString("@CYour body begins to shrink back to its normal form as the power of the Oozaru leaves you. You fall asleep shortly after returning to normal!@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@c$n@C's body begins to shrink and return to normal. Their giant ape features fading back into humanoid features until $e is left normal and asleep.@n"), TRUE, vict, nil, nil, TO_ROOM)
						vict.Act[int(PLR_OOZARU/32)] &= bitvector_t(int32(^(1 << (int(PLR_OOZARU % 32)))))
						vict.Hit = (vict.Hit / 2) - 10000
						vict.Mana = (vict.Mana / 2) - 10000
						vict.Move = (vict.Move / 2) - 10000
						vict.Max_hit = vict.Basepl
						vict.Max_mana = vict.Baseki
						vict.Max_move = vict.Basest
						if vict.Move < 1 {
							vict.Move = 1
						}
						if vict.Mana < 1 {
							vict.Mana = 1
						}
						if vict.Hit < 1 {
							vict.Hit = 1
						}
					}
				}
			case 2:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as you draw it up to attack. Two blindingly quick slashes connect with @c$N@W's face as you fly past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with YOUR face as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with @c$N@W's face as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				if dmg > vict.Max_hit/5 && (int(vict.Race) != RACE_MAJIN && int(vict.Race) != RACE_BIO) {
					act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou have your head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Death_type = DTYPE_HEAD
					remove_limb(vict, 0)
					die(vict, ch)
					if AFF_FLAGGED(ch, AFF_GROUP) {
						group_gain(ch, vict)
					} else {
						solo_gain(ch, vict)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
						do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
						do_get(ch, libc.CString("all corpse"), 0, 0)
					}
				} else if dmg > vict.Max_hit/5 && (int(vict.Race) == RACE_MAJIN || int(vict.Race) == RACE_BIO) {
					if GET_SKILL(vict, SKILL_REGENERATE) > rand_number(1, 101) && vict.Mana >= vict.Max_mana/40 {
						act(libc.CString("@R$N@r has $S head cut off by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou have your head cut off by the attack but regenerate a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r has $S head cut off by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Mana -= vict.Max_mana / 40
						hurt(0, 0, ch, vict, nil, dmg, 1)
					} else {
						act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou have your head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Death_type = DTYPE_HEAD
						die(vict, ch)
						if AFF_FLAGGED(ch, AFF_GROUP) {
							group_gain(ch, vict)
						} else {
							solo_gain(ch, vict)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
							do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
							do_get(ch, libc.CString("all corpse"), 0, 0)
						}
					}
				} else {
					hurt(0, 0, ch, vict, nil, dmg, 1)
				}
				dmg *= int64(calc_critical(ch, 0))
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as you draw it up to attack. Two blindingly quick slashes connect with @c$N@W's gut as you fly past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with YOUR gut as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with @c$N@W's gut as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as you draw it up to attack. Two blindingly quick slashes connect with @c$N@W's arm as you fly past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with YOUR arm as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with @c$N@W's arm as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 70 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[1]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severs $N's left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left arm is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[1] = 0
						remove_limb(vict, 2)
					} else if (vict.Limb_condition[0]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severs $N's right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right arm is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[0] = 0
						remove_limb(vict, 1)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as you draw it up to attack. Two blindingly quick slashes connect with @c$N@W's leg as you fly past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with YOUR leg as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wchannels $s charged ki into the blade of $s sword. @rF@Rl@Ya@rm@Ri@Yn@rg @gg@Gre@wen@W energy burns around the blade as $e draws it up to attack. Two blindingly quick slashes connect with @c$N@W's leg as $e flies past, leaving a green after-image behind!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 70 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[3]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severs $N's left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left leg is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[3] = 0
						remove_limb(vict, 4)
					} else if (vict.Limb_condition[2]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severs $N's right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right leg is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[2] = 0
						remove_limb(vict, 3)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if int(ch.Skillperfs[SKILL_TSLASH]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 19, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a twin slash at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a twin slash at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_eraser(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_ERASER) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_ERASER)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_ERASER, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Eraser Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Eraser Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Eraser Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				if GET_SKILL(ch, SKILL_ERASER) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
				} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your eraser cannon!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W eraser cannon!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W eraser cannon!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					if GET_SKILL(ch, SKILL_ERASER) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					dmg = damtype(ch, 18, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your eraser cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W eraser cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W eraser cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 18, skill, SKILL_ERASER)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					if GET_SKILL(ch, SKILL_ERASER) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your eraser cannon misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a eraser cannon at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a eraser cannon at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					if GET_SKILL(ch, SKILL_ERASER) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your eraser cannon misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a eraser cannon at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a eraser cannon at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
				if GET_SKILL(ch, SKILL_ERASER) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
				} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 18, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou face @c$N@W quickly and open your mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in your throat, growing brighter as it rises up and out your mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wfaces you quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wfaces @c$N@W quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				fallthrough
			case 3:
				act(libc.CString("@WYou face @c$N@W quickly and open your mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in your throat, growing brighter as it rises up and out your mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wfaces you quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wfaces @c$N@W quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 4:
				act(libc.CString("@WYou face @c$N@W quickly and open your mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in your throat, growing brighter as it rises up and out your mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wfaces you quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wfaces @c$N@W quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou face @c$N@W quickly and open your mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in your throat, growing brighter as it rises up and out your mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wfaces you quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wfaces @c$N@W quickly and opens $s mouth wide. A @Yg@yo@Yl@yd@Ye@yn@W glow forms from deep in $s throat, growing brighter as it rises up and out $s mouth. Suddenly a powerful Eraser Cannon erupts and slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			if GET_SKILL(ch, SKILL_ERASER) >= 100 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
			} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
			} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 18, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a eraser cannon at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a eraser cannon at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
		if GET_SKILL(ch, SKILL_ERASER) >= 100 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
		} else if GET_SKILL(ch, SKILL_ERASER) >= 60 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
		} else if GET_SKILL(ch, SKILL_ERASER) >= 40 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_pbarrage(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_PBARRAGE) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_PBARRAGE)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_PBARRAGE, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Psychic Barrage before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Psychic Barrage before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Psychic Barrage before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your psychic barrage!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W psychic barrage!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W psychic barrage!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 31, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your psychic barrage, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W psychic barrage, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W psychic barrage, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 31, skill, SKILL_PBARRAGE)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your psychic barrage misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a psychic barrage at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a psychic barrage at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your psychic barrage misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a psychic barrage at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a psychic barrage at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 31, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou close your eyes for a moment as your charged ki is pooled into your brain. Your eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and look at @c$N@W intensly. Invisible waves of psychic energy slam into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at YOU intensly. Invisible waves of psychic energy slam into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at @c$N@W intensly. Invisible waves of psychic energy slam into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou close your eyes for a moment as your charged ki is pooled into your brain. Your eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and look at @c$N@W intensly. Invisible waves of psychic energy slam into $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at YOU intensly. Invisible waves of psychic energy slam into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at @c$N@W intensly. Invisible waves of psychic energy slam into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou close your eyes for a moment as your charged ki is pooled into your brain. Your eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and look at @c$N@W intensly. Invisible waves of psychic energy slam into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at YOU intensly. Invisible waves of psychic energy slam into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at @c$N@W intensly. Invisible waves of psychic energy slam into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou close your eyes for a moment as your charged ki is pooled into your brain. Your eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and look at @c$N@W intensly. Invisible waves of psychic energy slam into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at YOU intensly. Invisible waves of psychic energy slam into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at @c$N@W intensly. Invisible waves of psychic energy slam into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou close your eyes for a moment as your charged ki is pooled into your brain. Your eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and look at @c$N@W intensly. Invisible waves of psychic energy slam into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at YOU intensly. Invisible waves of psychic energy slam into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W closes $s eyes for a moment as $s charged ki is pooled into $s brain. @C$n@W's eyes snap open, flashing with @bb@Bl@Wu@we @yl@Yi@Wg@wht@W, and looks at @c$N@W intensly. Invisible waves of psychic energy slam into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_MBREAK) && rand_number(1, 4) == 4 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
				act(libc.CString("@mYour mind's eye has been shattered, you can't charge ki until you recover!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@M$n@m's mind has been damaged by the attack!@n"), TRUE, vict, nil, nil, TO_ROOM)
				vict.Affected_by[int(AFF_MBREAK/32)] |= 1 << (int(AFF_MBREAK % 32))
			} else if !AFF_FLAGGED(vict, AFF_SHOCKED) && rand_number(1, 4) == 4 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
				act(libc.CString("@MYour mind has been shocked!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@M$n@m's mind has been shocked!@n"), TRUE, vict, nil, nil, TO_ROOM)
				vict.Affected_by[int(AFF_SHOCKED/32)] |= 1 << (int(AFF_SHOCKED % 32))
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 31, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a psychic barrage at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a psychic barrage at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_geno(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		prob    int
		attperc float64    = 0.5
		minimum float64    = 0.4
		vict    *char_data = nil
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_GENOCIDE) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	vict = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("No one around here by that name.\r\n"))
			return
		}
	}
	if can_kill(ch, vict, nil, 3) == 0 {
		return
	}
	if handle_defender(vict, ch) != 0 {
		var def *char_data = vict.Defender
		vict = def
	}
	prob = init_skill(ch, SKILL_GENOCIDE)
	perc = rand_number(1, 115)
	if prob < perc-20 {
		act(libc.CString("@WYou raise one arm above your head and pour your charged ki there. A large swirling pink ball of energy begins to form above your raised hand. You lose concentration and the ball of energy dissipates!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W raises one arm above $s head and pours $s charged ki there. A large swirling pink ball of energy begins to form above $s raised hand. @C$n@W loses concentration and the ball of energy dissipates!@n"), TRUE, ch, nil, nil, TO_ROOM)
		hurt(0, 0, ch, vict, nil, 0, 1)
		pcost(ch, attperc, 0)
		improve_skill(ch, SKILL_GENOCIDE, 2)
		return
	}
	var obj *obj_data
	var dista int = int(15 - float64(ch.Aff_abils.Intel)*0.1)
	if GET_SKILL(ch, SKILL_GENOCIDE) >= 100 {
		dista -= 3
	} else if GET_SKILL(ch, SKILL_GENOCIDE) >= 60 {
		dista -= 2
	} else if GET_SKILL(ch, SKILL_GENOCIDE) >= 40 {
		dista -= 1
	}
	obj = read_object(83, VIRTUAL)
	obj_to_room(obj, vict.In_room)
	ch.Charge += ch.Max_hit / 10
	obj.Target = vict
	obj.Kicharge = damtype(ch, 41, prob, attperc)
	obj.Kitype = SKILL_GENOCIDE
	obj.User = ch
	obj.Distance = dista
	pcost(ch, attperc, 0)
	act(libc.CString("@WYou raise one arm above your head and pour your charged ki there. A large swirling pink ball of energy begins to form above your raised hand. You grin viciously as the @mG@Me@wn@mo@Mc@wi@md@Me@W attack is complete and you toss it at @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
	act(libc.CString("@C$n@W raises one arm above $s head and pours $s charged ki there. A large swirling pink ball of energy begins to form above $s raised hand. @C$n@W grins viciously as the @mG@Me@wn@mo@Mc@wi@md@Me@W attack is complete and $e tosses it at YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
	act(libc.CString("@C$n@W raises one arm above $s head and pours $s charged ki there. A large swirling pink ball of energy begins to form above $s raised hand. @C$n@W grins viciously as the @mG@Me@wn@mo@Mc@wi@md@Me@W attack is complete and $e tosses it at @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	improve_skill(ch, SKILL_GENOCIDE, 2)
}
func do_genki(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		prob    int
		attperc float64    = 0.5
		minimum float64    = 0.4
		friend  *char_data = nil
		vict    *char_data = nil
		next_v  *char_data = nil
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_GENKIDAMA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	vict = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("No one around here by that name.\r\n"))
			return
		}
	}
	if can_kill(ch, vict, nil, 3) == 0 {
		return
	}
	if handle_defender(vict, ch) != 0 {
		var def *char_data = vict.Defender
		vict = def
	}
	prob = init_skill(ch, SKILL_GENKIDAMA)
	perc = rand_number(1, 115)
	if prob < perc-20 {
		act(libc.CString("@WYou raise both your arms upwards and begin to pool your charged ki there. You also start calling on the ki of all living beings in the vicinity who are willing to help. Your concentration wavers though and you waste the energy you have!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W raises both $s arms upwards and begin to pool $s charged ki there. @C$n@W also starts calling on the ki of all living beings in the vicinity who are willing to help. @C$n@W's concentration wavers though and $e wastes the energy $e has!@n"), TRUE, ch, nil, nil, TO_ROOM)
		hurt(0, 0, ch, vict, nil, 0, 1)
		pcost(ch, attperc, 0)
		handle_cooldown(ch, 10)
		improve_skill(ch, SKILL_GENKIDAMA, 2)
		return
	}
	for friend = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; friend != nil; friend = next_v {
		next_v = friend.Next_in_room
		if friend == ch {
			continue
		}
		if AFF_FLAGGED(friend, AFF_GROUP) && (friend.Master == ch || ch.Master == friend || friend.Master == ch.Master) {
			ch.Charge += ch.Mana / 10
			ch.Mana -= ch.Mana / 10
		}
	}
	var dista int = int(15 - float64(ch.Aff_abils.Intel)*0.1)
	if GET_SKILL(ch, SKILL_GENKIDAMA) >= 100 {
		dista -= 3
	} else if GET_SKILL(ch, SKILL_GENKIDAMA) >= 60 {
		dista -= 2
	} else if GET_SKILL(ch, SKILL_GENKIDAMA) >= 40 {
		dista -= 1
	}
	var obj *obj_data
	obj = read_object(82, VIRTUAL)
	obj_to_room(obj, vict.In_room)
	obj.Target = vict
	obj.Kicharge = damtype(ch, 40, prob, attperc)
	obj.Kitype = SKILL_GENKIDAMA
	obj.User = ch
	obj.Distance = dista
	pcost(ch, attperc, 0)
	act(libc.CString("@WYou raise both your arms upwards and begin to pool your charged ki there. You also start calling on the ki of all living beings in the vicinity who are willing to help. A large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W forms above your hands, when it is finished you lob it toward @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
	act(libc.CString("@C$n@W raises both $s arms upwards and begin to pool $s charged ki there. @C$n@W also starts calling on the ki of all living beings in the vicinity who are willing to help. A large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W forms above $s hands, when it is finished $e lobs it toward YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
	act(libc.CString("@C$n@W raises both $s arms upwards and begin to pool $s charged ki there. @C$n@W also starts calling on the ki of all living beings in the vicinity who are willing to help. A large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W forms above $s hands, when it is finished $e lobs it toward @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	handle_cooldown(ch, 10)
	improve_skill(ch, SKILL_GENKIDAMA, 2)
}
func do_spiritball(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.25
		minimum float64 = 0.2
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SPIRITBALL) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SPIRITBALL)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 8)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SPIRITBALL, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Spirit Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Spirit Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Spirit Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				dodge_ki(ch, vict, 2, 39, skill, SKILL_SPIRITBALL)
				hurt(0, 0, ch, vict, nil, 0, 1)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects your Spirit Ball, sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect @C$n's@W Spirit Ball sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects @c$n's@W Spirit Ball sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(vict, 0, vict.Max_hit/200)
					parry_ki(attperc, ch, vict, func() [1000]byte {
						var t [1000]byte
						copy(t[:], []byte("spiritball"))
						return t
					}(), prob, perc, skill, 39)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks the Spirit Ball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Spirit Ball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Spirit Ball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 39, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Spirit Ball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Spirit Ball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Spirit Ball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your spirit ball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a spirit ball at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a spirit ball at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 2, 39, skill, SKILL_SPIRITBALL)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your spirit ball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a spirit ball at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a spirit ball at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dodge_ki(ch, vict, 2, 39, skill, SKILL_SPIRITBALL)
				pcost(ch, attperc, 0)
				hurt(0, 0, ch, vict, nil, 0, 1)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 39, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise your palm upward with your arm bent in front of you. Your charged ki slowly begins to creep up your arm and form a large glowing orb of energy above your upraised palm. With your @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted you move your hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into $S body, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at YOU! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at you, slamming into YOUR body, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into @c$N@W's body, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou raise your palm upward with your arm bent in front of you. Your charged ki slowly begins to creep up your arm and form a large glowing orb of energy above your upraised palm. With your @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted you move your hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into $S head, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at YOU! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at you, slamming into YOUR head, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into @c$N@W's head, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou raise your palm upward with your arm bent in front of you. Your charged ki slowly begins to creep up your arm and form a large glowing orb of energy above your upraised palm. With your @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted you move your hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into $S gut, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at YOU! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at you, slamming into YOUR gut, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into @c$N@W's gut, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou raise your palm upward with your arm bent in front of you. Your charged ki slowly begins to creep up your arm and form a large glowing orb of energy above your upraised palm. With your @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted you move your hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into $S arm, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at YOU! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at you, slamming into YOUR arm, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into @c$N@W's arm, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou raise your palm upward with your arm bent in front of you. Your charged ki slowly begins to creep up your arm and form a large glowing orb of energy above your upraised palm. With your @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted you move your hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into $S leg, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at YOU! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at you, slamming into YOUR leg, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s palm upward with $s arm bent in front of $m. @C$n@W's charged ki slowly begins to creep up $s arm and form a large glowing orb of energy above $s upraised palm. With $s @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wcompleted $e moves $s hand with index and middle fingers pointing at @c$N@W! The @yS@Yp@Wi@wri@yt @YB@Wa@wll @Wflies at $M, slamming into @c$N@W's leg, and explodes with a load roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 39, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a spirit ball at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a spirit ball at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_deathball(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.3
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_DEATHBALL) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_DEATHBALL)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 8)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_DEATHBALL, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		prob -= rand_number(8, 10)
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Deathball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Deathball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Deathball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your deathball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W deathball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W deathball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 38, skill, SKILL_DEATHBALL)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 80 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 20
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your deathball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a deathball at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a deathball at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your deathball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a deathball at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a deathball at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 38, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise your hand with its index finger extended upwards. Your charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as you look down on @c$N@W. You move the hand that formed the attack and point at $M as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above you follows the movement! It descends on $M and crushes into $S body before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on YOU. @C$n@W moves the hand that formed the attack and points at YOU as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on YOU and crushes into YOUR body before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on @c$N@W. @C$n@W moves the hand that formed the attack and points at @c$N@W as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on @c$N@W and crushes into $S body before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou raise your hand with its index finger extended upwards. Your charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as you look down on @c$N@W. You move the hand that formed the attack and point at $M as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above you follows the movement! It descends on $M and crushes into $S head before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on YOU. @C$n@W moves the hand that formed the attack and points at YOU as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on YOU and crushes into YOUR head before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on @c$N@W. @C$n@W moves the hand that formed the attack and points at @c$N@W as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on @c$N@W and crushes into $S head before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= 4
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou raise your hand with its index finger extended upwards. Your charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as you look down on @c$N@W. You move the hand that formed the attack and point at $M as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above you follows the movement! It descends on $M and crushes into $S gut before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on YOU. @C$n@W moves the hand that formed the attack and points at YOU as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on YOU and crushes into YOUR gut before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on @c$N@W. @C$n@W moves the hand that formed the attack and points at @c$N@W as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on @c$N@W and crushes into $S gut before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou raise your hand with its index finger extended upwards. Your charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as you look down on @c$N@W. You move the hand that formed the attack and point at $M as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above you follows the movement! It descends on $M and crushes into $S arm before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on YOU. @C$n@W moves the hand that formed the attack and points at YOU as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on YOU and crushes into YOUR arm before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on @c$N@W. @C$n@W moves the hand that formed the attack and points at @c$N@W as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on @c$N@W and crushes into $S arm before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou raise your hand with its index finger extended upwards. Your charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as you look down on @c$N@W. You move the hand that formed the attack and point at $M as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above you follows the movement! It descends on $M and crushes into $S leg before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on YOU. @C$n@W moves the hand that formed the attack and points at YOU as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on YOU and crushes into YOUR leg before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand with its index finger extended upwards. @C$n@W's charged ki begins to pool above that finger, forming a small @rred@W orb of energy. The orb of energy quickly grows to an enormous size as $e looks down on @c$N@W. @C$n@W moves the hand that formed the attack and points at @c$N@W as the @rD@Re@Da@rt@Rh@Db@ra@Rl@Dl@W above $m follows the movement! It descends on @c$N@W and crushes into $S leg before exploding into a massive blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 38, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a deathball at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a deathball at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_pslash(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.3
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_PSLASH) == 0 {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	} else if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 && (ch.Limb_condition[1]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right arm has damaged it more!@n\r\n"))
		ch.Limb_condition[0] -= rand_number(3, 5)
		if (ch.Limb_condition[0]) < 0 {
			act(libc.CString("@RYour right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 && (ch.Limb_condition[0]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left arm has damaged it more!@n\r\n"))
		ch.Limb_condition[1] -= rand_number(3, 5)
		if (ch.Limb_condition[1]) < 0 {
			act(libc.CString("@RYour left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) == nil {
		send_to_char(ch, libc.CString("You need to wield a sword to use this.\r\n"))
		return
	}
	if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) != int(TYPE_SLASH-TYPE_HIT) {
		send_to_char(ch, libc.CString("You are not wielding a sword, you need one to use this technique.\r\n"))
		return
	}
	skill = init_skill(ch, SKILL_PSLASH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 8)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_PSLASH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Phoenix Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Phoenix Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Phoenix Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your Phoenix Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Phoenix Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Phoenix Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 37, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Phoenix Slash, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Phoenix Slash, letting it letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Phoenix Slash, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Phoenix Slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Phoenix Slash at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Phoenix Slash at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Phoenix Slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Phoenix Slash at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Phoenix Slash at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 37, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			if check_ruby(ch) == 1 {
				dmg += int64(float64(dmg) * 0.2)
			}
			if (vict.Bonuses[BONUS_FIREPROOF]) != 0 {
				dmg -= int64(float64(dmg) * 0.4)
			} else if (vict.Bonuses[BONUS_FIREPRONE]) != 0 {
				dmg += int64(float64(dmg) * 0.4)
			}
			vict.Affected_by[int(AFF_ASHED/32)] |= 1 << (int(AFF_ASHED % 32))
			switch hitspot {
			case 1:
				act(libc.CString("@WYou pour your charged ki into your sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surround the entire sword in the same instant that you pull the blade back and extend it behind your body. Suddenly you swing the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs $S body in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward YOU, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards YOU! The Phoenix Slash engulfs YOUR body in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs @c$N@W's body in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou pour your charged ki into your sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surround the entire sword in the same instant that you pull the blade back and extend it behind your body. Suddenly you swing the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs $S head in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward YOU, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards YOU! The Phoenix Slash engulfs YOUR head in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs @c$N@W's head in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou pour your charged ki into your sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surround the entire sword in the same instant that you pull the blade back and extend it behind your body. Suddenly you swing the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs $S gut in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward YOU, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards YOU! The Phoenix Slash engulfs YOUR gut in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs @c$N@W's gut in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou pour your charged ki into your sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surround the entire sword in the same instant that you pull the blade back and extend it behind your body. Suddenly you swing the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs $S arm in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward YOU, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards YOU! The Phoenix Slash engulfs YOUR arm in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs @c$N@W's arm in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou pour your charged ki into your sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surround the entire sword in the same instant that you pull the blade back and extend it behind your body. Suddenly you swing the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs $S leg in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward YOU, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards YOU! The Phoenix Slash engulfs YOUR leg in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged ki into $s sword's blade. @YF@ri@Re@Yr@ry @Rf@Yl@ra@Rm@Ye@rs@W surrounds the entire sword in the same instant that $e pulls the blade back and extends it behind $s body. Suddenly $e swings the blade forward toward @c$N@W, unleashing a large wave of flames! The flames take the shape of a large phoenix soaring towards @c$N@W! The Phoenix Slash engulfs @c$N@W's leg in flames a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_BURNED) && rand_number(1, 4) == 4 && int(vict.Race) != RACE_DEMON && (vict.Bonuses[BONUS_FIREPROOF]) == 0 {
				send_to_char(vict, libc.CString("@RYou are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are burned by the attack!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			} else if (vict.Bonuses[BONUS_FIREPROOF]) != 0 || int(vict.Race) == RACE_DEMON {
				send_to_char(ch, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			} else if (vict.Bonuses[BONUS_FIREPRONE]) != 0 {
				send_to_char(vict, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are easily burned!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			}
			pcost(ch, attperc, 0)
			vict.Affected_by[int(AFF_ASHED/32)] &= ^(1 << (int(AFF_ASHED % 32)))
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 37, skill, attperc)
		dmg /= 10
		if GET_OBJ_VNUM(obj) == 79 {
			dmg *= 3
		}
		act(libc.CString("@WYou fire a Phoenix Slash at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Phoenix Slash at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_bigbang(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.3
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BIGBANG) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_BIGBANG)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 8)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_BIGBANG, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Big Bang before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Big Bang before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Big Bang before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your Big Bang!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Big Bang!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Big Bang!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 36, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Big Bang, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Big Bang, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Big Bang, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 36, skill, SKILL_BIGBANG)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 80 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 20
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Big Bang misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Big Bang at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Big Bang at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Big Bang misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Big Bang at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Big Bang at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 36, skill, attperc)
			switch rand_number(1, 6) {
			case 1:
				act(libc.CString("@WYou aim your hand at @c$N@W palm flat up at a ninety degree angle. Your charged ki pools there rapidly before a massive ball of energy explodes from your palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's body, and explodes leaving behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand at YOU palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into YOUR body, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s hand at @c$N@W palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's body, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				fallthrough
			case 3:
				act(libc.CString("@WYou aim your hand at @c$N@W palm flat up at a ninety degree angle. Your charged ki pools there rapidly before a massive ball of energy explodes from your palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's head, and explodes leaving behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand at YOU palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into YOUR head, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s hand at @c$N@W palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's head, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 4:
				act(libc.CString("@WYou aim your hand at @c$N@W palm flat up at a ninety degree angle. Your charged ki pools there rapidly before a massive ball of energy explodes from your palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's gut, and explodes leaving behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand at YOU palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into YOUR gut, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s hand at @c$N@W palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's gut, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 5:
				act(libc.CString("@WYou aim your hand at @c$N@W palm flat up at a ninety degree angle. Your charged ki pools there rapidly before a massive ball of energy explodes from your palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's arm, and explodes leaving behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand at YOU palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into YOUR arm, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s hand at @c$N@W palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's arm, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 6:
				act(libc.CString("@WYou aim your hand at @c$N@W palm flat up at a ninety degree angle. Your charged ki pools there rapidly before a massive ball of energy explodes from your palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's leg, and explodes leaving behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand at YOU palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into YOUR leg, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s hand at @c$N@W palm flat up at a ninety degree angle. @C$n@W's charged ki pools there rapidly before a massive ball of energy explodes from $s palm! The @yB@Yi@wg @yB@Ya@wn@yg @Wattack crosses the distance rapidly, slams into @c$N@W's leg, and explodes leading behind a mushroom cloud shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 36, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a Big Bang at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Big Bang at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_scatter(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.3
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SCATTER) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SCATTER)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	var cool int = 8
	if int(ch.Chclass) == CLASS_PICCOLO {
		if int(ch.Skills[SKILL_STYLE]) >= 100 {
			cool -= 3
		} else if int(ch.Skills[SKILL_STYLE]) >= 75 {
			cool -= 2
		} else if int(ch.Skills[SKILL_STYLE]) >= 40 {
			cool -= 1
		}
	}
	if cool < 1 {
		cool = 1
	}
	handle_cooldown(ch, cool)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SCATTER, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		prob += rand_number(10, 20)
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Scatter Shot before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Scatter Shot before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Scatter Shot before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks every kiball of your scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block every kiball of @C$n's@W scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks every kiball of @c$n's@W scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 35, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your scatter shot kiballs, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W scatter shot kiballs, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W scatter shot kiballs, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wBright explosions erupts from the impacts!\r\n"))
					dodge_ki(ch, vict, 0, 35, skill, SKILL_SCATTER)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 80 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 20
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but all the kiballs of your scatter shot miss, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires scatter shot kiballs at you, but they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires scatter shot kiballs at @C$N@W, but somehow they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but all the kiballs of your scatter shot miss, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires scatter shot kiballs at you, but they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires scatter shot kiballs at @C$N@W, but somehow they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 35, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou bring your charged ki to your palms and furiously begin to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! You hold out your hand and clench it dramatically as your @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at YOU! The kiballs surround you from every direction in a sphere, preventing your escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou bring your charged ki to your palms and furiously begin to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! You hold out your hand and clench it dramatically as your @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at YOU! The kiballs surround you from every direction in a sphere, preventing your escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou bring your charged ki to your palms and furiously begin to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! You hold out your hand and clench it dramatically as your @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at YOU! The kiballs surround you from every direction in a sphere, preventing your escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou bring your charged ki to your palms and furiously begin to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! You hold out your hand and clench it dramatically as your @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at YOU! The kiballs surround you from every direction in a sphere, preventing your escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou bring your charged ki to your palms and furiously begin to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! You hold out your hand and clench it dramatically as your @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at YOU! The kiballs surround you from every direction in a sphere, preventing your escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W brings $s charged ki to $s palms and furiously begins to throw hundreds of kiballs at @C$N@W! The kiballs surround $M from every direction in a sphere, preventing $S escape! @C$n@W holds out $s hand and clenches it dramatically as $s @yS@Yc@ra@Rt@yt@Yer @yS@Yh@ro@Rt@W kiballs close in and explode against @C$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 35, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a scatter shot at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a scatter shot at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_balefire(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.3
		minimum float64 = 0.15
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BALEFIRE) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_BALEFIRE)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	var cool int = 8
	if int(ch.Chclass) == CLASS_PICCOLO {
		if int(ch.Skills[SKILL_STYLE]) >= 100 {
			cool -= 3
		} else if int(ch.Skills[SKILL_STYLE]) >= 75 {
			cool -= 2
		} else if int(ch.Skills[SKILL_STYLE]) >= 40 {
			cool -= 1
		}
	}
	if cool < 1 {
		cool = 1
	}
	handle_cooldown(ch, cool)
	if vict != nil {
		if can_kill(ch, vict, nil, 3) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_BALEFIRE, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		prob += rand_number(10, 20)
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Balefire before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Balefire before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Balefire before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks every kiball of your scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block every kiball of @C$n's@W scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks every kiball of @c$n's@W scatter shot!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 35, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your balefire, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W balefire, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W balefire, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wBright explosions erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 35, skill, SKILL_BALEFIRE)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 80 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 20
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but all the kiballs of your scatter shot miss, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires scatter shot kiballs at you, but they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires scatter shot kiballs at @C$N@W, but somehow they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but all of your balefire misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires balefire at you, but they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires balefire at @C$N@W, but somehow they miss!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 35, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou open yourself to the source and begin tracing weaves, the most complex combination, with just the right amounts of each! Your innate talent finishes the final weave and a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from your hands! Your @WB@Ra@Yl@Wef@Yi@Rr@We@W crosses the distance quickly and slams into @C$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into @C$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou open yourself to the source and begin tracing weaves, the most complex combination, with just the right amounts of each! Your innate talent finishes the final weave and a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from your hands! Your @WB@Ra@Yl@Wef@Yi@Rr@We@W crosses the distance quickly and slams into @C$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into @C$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou open yourself to the source and begin tracing weaves, the most complex combination, with just the right amounts of each! Your innate talent finishes the final weave and a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from your hands! Your @WB@Ra@Yl@Wef@Yi@Rr@We@W crosses the distance quickly and slams into @C$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into @C$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou open yourself to the source and begin tracing weaves, the most complex combination, with just the right amounts of each! Your innate talent finishes the final weave and a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from your hands! Your @WB@Ra@Yl@Wef@Yi@Rr@We@W crosses the distance quickly and slams into @C$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into @C$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou open yourself to the source and begin tracing weaves, the most complex combination, with just the right amounts of each! Your innate talent finishes the final weave and a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from your hands! Your @WB@Ra@Yl@Wef@Yi@Rr@We@W crosses the distance quickly and slams into @C$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W begins to glow with a soft, white aura! Their hands begin manipulating something unseen when suddenly a giant bar of @WB@Ra@Yl@Wef@Yi@Rr@We@W lances out from $s hands, crossing the distance and slamming into @C$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 35, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire balefire at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires balefire at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kakusanha(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		dge     int = 2
		count   int = 0
		skill   int
		dmg     int64
		attperc float64    = 0.3
		minimum float64    = 0.1
		vict    *char_data = nil
		next_v  *char_data = nil
		arg2    [2048]byte
	)
	one_argument(argument, &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KAKUSANHA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KAKUSANHA)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		}
		if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
			continue
		}
		if vict.Hit <= 0 {
			continue
		}
		if AFF_FLAGGED(vict, AFF_GROUP) && !IS_NPC(vict) {
			if vict.Master == ch {
				continue
			} else if ch.Master == vict {
				continue
			} else if vict.Master == ch.Master {
				continue
			}
		}
		if MOB_FLAGGED(vict, MOB_NOKILL) {
			continue
		}
		if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
			continue
		} else {
			count += 1
		}
	}
	if count <= 0 {
		send_to_char(ch, libc.CString("There is no one worth targeting around.\r\n"))
		return
	} else {
		handle_cooldown(ch, 5)
		var hits int = 0
		perc = chance_to_hit(ch)
		if skill < perc {
			act(libc.CString("@WYou pour your charged ki into your hands and bring them both forward quickly. @yG@Yo@Wl@wden@W orbs of energy form at the extent of your palms. You fire one massive beam of @yg@Yo@Wl@wden @Wenergy from combining both orbs. The beam flies forward a short distance before you swing your arms upward and the beam follows suit. The beam flies off harmlessly through the air as you lose control over it!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W pours $s charged ki into $s hands and brings them both forward quickly. @yG@Yo@Wl@wden@W orbs of energy form at the extent of $s palms. @C$n@W fires one massive beam of @yg@Yo@Wl@wden @Wenergy from combining both orbs. The beam flies forward a short distance before $e swings $s arms upward and the beam follows suit. The beam flies off harmlessly through the air as $e loses control over it!@n"), TRUE, ch, nil, nil, TO_ROOM)
			pcost(ch, attperc, 0)
			improve_skill(ch, SKILL_KAKUSANHA, 0)
			return
		}
		act(libc.CString("@WYou pour your charged ki into your hands and bring them both forward quickly. @yG@Yo@Wl@wden@W orbs of energy form at the extent of your palms. You fire one massive beam of @yg@Yo@Wl@wden @Wenergy from combining both orbs. The beam flies forward a short distance before you swing your arms upward and the beam follows suit. Above your targets the beam breaks apart into five seperate pieces that follow their victims!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W pours $s charged ki into $s hands and brings them both forward quickly. @yG@Yo@Wl@wden@W orbs of energy form at the extent of $s palms. @C$n@W fires one massive beam of @yg@Yo@Wl@wden @Wenergy from combining both orbs. The beam flies forward a short distance before $e swings $s arms upward and the beam follows suit. Above you the beam breaks apart into five seperate pieces that follow their victims!@n"), TRUE, ch, nil, nil, TO_ROOM)
		dmg = damtype(ch, 34, skill, attperc)
		if count >= 3 {
			dmg = (dmg / 100) * 40
		} else if count > 1 {
			dmg = (dmg / 100) * 60
		}
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
			next_v = vict.Next_in_room
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
				continue
			}
			if AFF_FLAGGED(vict, AFF_GROUP) && (vict.Master == ch || ch.Master == vict) {
				continue
			}
			if MOB_FLAGGED(vict, MOB_NOKILL) {
				continue
			}
			if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
				continue
			}
			if vict.Hit <= 0 {
				continue
			}
			if hits >= 5 {
				continue
			}
			dge = handle_dodge(vict)
			if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
				hits++
				act(libc.CString("@C$N@c disappears, avoiding the beam chasing $M!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding the beam chasing you!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding the beam chasing $M!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(vict, 0, vict.Max_hit/200)
				hurt(0, 0, ch, vict, nil, 0, 1)
				continue
			} else if dge+rand_number(-10, 5) > skill {
				hits++
				act(libc.CString("@c$N@W manages to escape the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@WYou manage to escape the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N@W manages to escape the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				improve_skill(vict, SKILL_DODGE, 0)
				hurt(0, 0, ch, vict, nil, 0, 1)
				continue
			} else {
				hits++
				act(libc.CString("@R$N@r is slammed by one of the beams!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@RYou are slammed by one of the beams!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is slammed by one of the beams!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				continue
			}
		}
		if count < 5 && !ROOM_FLAGGED(ch.In_room, ROOM_SPACE) {
			send_to_room(ch.In_room, libc.CString("The rest of the beams slam into the ground!@n\r\n"))
			send_to_room(ch.In_room, libc.CString("@wBright explosions erupt from the impacts!\r\n"))
			if (func() int {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) != SECT_INSIDE {
				impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
				switch rand_number(1, 8) {
				case 1:
					act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, nil, TO_ROOM)
				case 2:
					if rand_number(1, 4) == 4 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect == 0 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect = 5
						act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, nil, TO_CHAR)
						act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, nil, TO_ROOM)
					}
				case 3:
					act(libc.CString("A cloud of dust envelopes the entire area!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("A cloud of dust envelopes the entire area!"), TRUE, ch, nil, nil, TO_ROOM)
				case 4:
					act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				case 5:
					act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				case 6:
					act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, nil, TO_ROOM)
				default:
				}
			}
			if (func() int {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) == SECT_UNDERWATER {
				switch rand_number(1, 3) {
				case 1:
					act(libc.CString("The water churns violently!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The water churns violently!"), TRUE, ch, nil, nil, TO_ROOM)
				case 2:
					act(libc.CString("Large bubbles rise from the movement!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("Large bubbles rise from the movement!"), TRUE, ch, nil, nil, TO_ROOM)
				case 3:
					act(libc.CString("The water collapses in on the hole created!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The water collapses in on the hole create!"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			if (func() int {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) == SECT_WATER_SWIM || (func() int {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) == SECT_WATER_NOSWIM {
				switch rand_number(1, 3) {
				case 1:
					act(libc.CString("A huge column of water erupts from the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("A huge column of water erupts from the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				case 2:
					act(libc.CString("The impact briefly causes a swirling vortex of water!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The impact briefly causes a swirling vortex of water!"), TRUE, ch, nil, nil, TO_ROOM)
				case 3:
					act(libc.CString("A huge depression forms in the water and erupts into a wave from the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("A huge depression forms in the water and erupts into a wave from the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				}
			}
			if (func() int {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
				}
				return SECT_INSIDE
			}()) == SECT_INSIDE {
				impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
				switch rand_number(1, 8) {
				case 1:
					act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, nil, TO_ROOM)
				case 2:
					act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, nil, TO_ROOM)
				case 3:
					act(libc.CString("Parts of the ceiling collapse, crushing into the floor!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("Parts of the ceiling collapse, crushing into the floor!"), TRUE, ch, nil, nil, TO_ROOM)
				case 4:
					act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				case 5:
					act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, nil, TO_ROOM)
				case 6:
					act(libc.CString("The walls of the surrounding room crack in the same instant!"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("The walls of the surrounding room crack in the same instant!"), TRUE, ch, nil, nil, TO_ROOM)
				default:
				}
			}
			if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= (100 - (5-count)*5) {
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += (5 - count) * 5
			}
		}
		pcost(ch, attperc, 0)
		improve_skill(ch, SKILL_KAKUSANHA, 0)
		handle_cooldown(ch, 5)
		return
	}
}
func do_hellspear(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		dge     int = 2
		count   int = 0
		skill   int
		dmg     int64
		attperc float64    = 0.3
		minimum float64    = 0.1
		vict    *char_data = nil
		next_v  *char_data = nil
		arg2    [2048]byte
	)
	one_argument(argument, &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_HELLSPEAR) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
		send_to_char(ch, libc.CString("This room just has such a peaceful, easy feeling...\r\n"))
		return
	}
	skill = init_skill(ch, SKILL_HELLSPEAR)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		}
		if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
			continue
		}
		if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
			continue
		}
		if MOB_FLAGGED(vict, MOB_NOKILL) {
			continue
		} else {
			count += 1
		}
	}
	if count <= 0 {
		send_to_char(ch, libc.CString("There is no one worth targeting around.\r\n"))
		return
	} else {
		perc = chance_to_hit(ch)
		if skill < perc {
			act(libc.CString("@WYou fly up higher in the air while holding your hand up above your head. Your charged ki is condensed and materialized in the grasp of your raised hand forming a spear of energy. The spear disolves as you screw up the technique!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W flies up higher in the air while holding $s hand up above $s head. @C$n@W's charged ki is condensed and materialized in the grasp of $s raised hand forming a spear of energy. The spear disolves as $e screws up the technique!@n"), TRUE, ch, nil, nil, TO_ROOM)
			pcost(ch, attperc, 0)
			improve_skill(ch, SKILL_HELLSPEAR, 0)
			return
		}
		act(libc.CString("@WYou fly up higher in the air while holding your hand above your head. Your charged ki is condensed and materialized in the grasp of your raised hand forming a spear of energy. Grinning evily you aim the spear at the ground below and throw it. As the red spear of energy slams into the ground your laughter rings throughout the area, and the @rH@Re@Dl@Wl @rS@Rp@De@Wa@wr B@rl@Ra@Ds@wt@W erupts with a roar!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W flies up higher in the air while holding $s hand above $s head. @C$n@W's charged ki is condensed and materialized in the grasp of $s raised hand forming a spear of energy. Grinning evily $e aims the spear at the ground below and throws it. As the red spear of energy slams into the ground $s laughter rings throughout the area, and the @rH@Re@Dl@Wl @rS@Rp@De@Wa@wr B@rl@Ra@Ds@wt@W erupts with a roar!@n"), TRUE, ch, nil, nil, TO_ROOM)
		dmg = damtype(ch, 33, skill, attperc)
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
			next_v = vict.Next_in_room
			if vict == ch {
				continue
			}
			if MOB_FLAGGED(vict, MOB_NOKILL) {
				continue
			}
			if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
				continue
			}
			if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
				continue
			}
			dge = handle_dodge(vict)
			if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
				act(libc.CString("@C$N@c disappears, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(vict, 0, vict.Max_hit/200)
				continue
			} else if dge+rand_number(-10, 5) > skill {
				act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@WYou manage to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 1)
				improve_skill(vict, SKILL_DODGE, 0)
				continue
			} else {
				act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@RYou are caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 4) == 4 {
					handle_knockdown(vict)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				continue
			}
		}
		pcost(ch, attperc, 0)
		improve_skill(ch, SKILL_HELLSPEAR, 0)
		handle_cooldown(ch, 5)
		return
	}
}
func do_hellflash(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_HELLFLASH) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_HELLFLASH]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_HELLFLASH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_HELLFLASH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_HELLFLASH]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Hell Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Hell Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Hell Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your Hell Flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Hell Flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Hell Flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 32, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Hell Flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Hell Flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Hell Flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 32, skill, SKILL_HELLFLASH)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Hell Flash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Hell Flash at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Hell Flash at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Hell Flash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Hell Flash at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Hell Flash at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 32, skill, attperc)
			if vict.Barrier > 0 {
				vict.Barrier -= dmg
				if vict.Barrier <= 0 {
					vict.Barrier = 1
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou stick one of your hands in each of your armpits and by twisting you remove them. Then you aim your exposed wrist cannons at @c$N@W as your charged energy begins to be channeled into the cannons. Shouting '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' two large blasts of energy explode from the cannons and slam into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at YOU as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at @c$N@W as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou stick one of your hands in each of your armpits and by twisting you remove them. Then you aim your exposed wrist cannons at @c$N@W as your charged energy begins to be channeled into the cannons. Shouting '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' two large blasts of energy explode from the cannons and slam into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at YOU as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at @c$N@W as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou stick one of your hands in each of your armpits and by twisting you remove them. Then you aim your exposed wrist cannons at @c$N@W as your charged energy begins to be channeled into the cannons. Shouting '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' two large blasts of energy explode from the cannons and slam into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at YOU as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at @c$N@W as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou stick one of your hands in each of your armpits and by twisting you remove them. Then you aim your exposed wrist cannons at @c$N@W as your charged energy begins to be channeled into the cannons. Shouting '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' two large blasts of energy explode from the cannons and slam into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at YOU as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at @c$N@W as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou stick one of your hands in each of your armpits and by twisting you remove them. Then you aim your exposed wrist cannons at @c$N@W as your charged energy begins to be channeled into the cannons. Shouting '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' two large blasts of energy explode from the cannons and slam into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at YOU as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W sticks one of $s hands in each of $s armpits and by twisting $e removes them. Then $e aims $s exposed wrist cannons at @c$N@W as $s charged energy begins to be channeled into the cannons. @C$n@W shouts '@YH@re@Rl@Yl @rF@Rl@Ya@rs@Rh@W' as two large blasts of energy explode from the cannons and slam into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if int(ch.Skillperfs[SKILL_HELLFLASH]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 32, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a Hell Flash at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Hell Flash at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_ddslash(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_DDSLASH) == 0 {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	} else if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 && (ch.Limb_condition[1]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right arm has damaged it more!@n\r\n"))
		ch.Limb_condition[0] -= rand_number(3, 5)
		if (ch.Limb_condition[0]) < 0 {
			act(libc.CString("@RYour right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 && (ch.Limb_condition[0]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left arm has damaged it more!@n\r\n"))
		ch.Limb_condition[1] -= rand_number(3, 5)
		if (ch.Limb_condition[1]) < 0 {
			act(libc.CString("@RYour left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) == nil {
		send_to_char(ch, libc.CString("You need to wield a sword to use this.\r\n"))
		return
	}
	if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) != int(TYPE_SLASH-TYPE_HIT) {
		send_to_char(ch, libc.CString("You are not wielding a sword, you need one to use this technique.\r\n"))
		return
	}
	skill = init_skill(ch, SKILL_DDSLASH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_DDSLASH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Darkness Dragon Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Darkness Dragon Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Darkness Dragon Slash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your Darkness Dragon Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Darkness Dragon Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Darkness Dragon Slash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 30, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your Darkness Dragon Slash, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Darkness Dragon Slash, letting it letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Darkness Dragon Slash, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Darkness Dragon Slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Darkness Dragon Slash at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Darkness Dragon Slash at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Darkness Dragon Slash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Darkness Dragon Slash at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Darkness Dragon Slash at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 30, skill, attperc)
			if GET_SKILL(ch, SKILL_DDSLASH) >= 100 {
				dmg += int64(float64(dmg) * 0.15)
			} else if GET_SKILL(ch, SKILL_DDSLASH) >= 60 {
				dmg += int64(float64(dmg) * 0.1)
			} else if GET_SKILL(ch, SKILL_DDSLASH) >= 40 {
				dmg += int64(float64(dmg) * 0.05)
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging your sword at @c$N@W you unleash a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at YOU $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at @c$N@W $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging your sword at @c$N@W you unleash a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at YOU $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at @c$N@W $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging your sword at @c$N@W you unleash a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at YOU $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at @c$N@W $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging your sword at @c$N@W you unleash a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at YOU $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at @c$N@W $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou channel your charged ki into the blade of your sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging your sword at @c$N@W you unleash a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at YOU $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W channels $s charged ki into the blade of $s sword. The energy takes the form of @md@Ma@Wr@wk@W @mf@Ml@Wa@wmes@W burning along the very blade's edge. Swinging $s sword at @c$N@W $e unleashes a serpentine dragon of @md@Ma@Wr@wk @mf@Ml@Wa@wmes@W! The fiery dragon slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if rand_number(1, 3) == 3 && !AFF_FLAGGED(vict, AFF_BLIND) {
				act(libc.CString("@mYou are struck blind temporarily!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@c$n@m is struck blind by the attack!@n"), TRUE, vict, nil, nil, TO_ROOM)
				var duration int = 1
				assign_affect(vict, AFF_BLIND, SKILL_SOLARF, duration, 0, 0, 0, 0, 0, 0)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 30, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a Darkness Dragon Slash at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Darkness Dragon Slash at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_crusher(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_CRUSHER) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_CRUSHER]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_CRUSHER)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_CRUSHER, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_CRUSHER]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Crusher Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Crusher Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Crusher Ball before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your crusher ball!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W crusher ball!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W crusher ball!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 29, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your crusher ball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W crusher ball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W crusher ball, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 29, skill, SKILL_CRUSHER)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 90 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 10
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your crusher ball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a crusher ball at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a crusher ball at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your crusher ball misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a crusher ball at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a crusher ball at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 29, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou zoom up into the air higher than @c$N@W and raise your hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W you slam your hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W zooms up into the air higher than you and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on YOU $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W zooms up into the air higher than @c$N@W and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou zoom up into the air higher than @c$N@W and raise your hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W you slam your hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W zooms up into the air higher than you and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on YOU $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W zooms up into the air higher than @c$N@W and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou zoom up into the air higher than @c$N@W and raise your hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W you slam your hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W zooms up into the air higher than you and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on YOU $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W zooms up into the air higher than @c$N@W and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou zoom up into the air higher than @c$N@W and raise your hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W you slam your hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W zooms up into the air higher than you and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on YOU $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W zooms up into the air higher than @c$N@W and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou zoom up into the air higher than @c$N@W and raise your hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W you slam your hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W zooms up into the air higher than you and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on YOU $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W zooms up into the air higher than @c$N@W and raises $s hand. Charged @rred@W ki pools above that hand in a large blazing ball. Looking down on @c$N@W $e slams $s hand into the ball of energy while shouting '@rC@Rr@Wu@wsher @rB@Ra@Wl@wl!@W' Moments later the ball of energy slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 3) == 3 {
				handle_knockdown(vict)
			} else if AFF_FLAGGED(vict, AFF_FLYING) && rand_number(1, 3) == 3 {
				act(libc.CString("@mYou are knocked out of the air by the attack!"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@W$n@m is knocked out of the air by the attack!"), TRUE, vict, nil, nil, TO_ROOM)
				vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
				ch.Altitude = 0
			}
			if int(ch.Skillperfs[SKILL_CRUSHER]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 29, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a crusher ball at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a crusher ball at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_final(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_FINALFLASH) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_FINALFLASH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_FINALFLASH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Final Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Final Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Final Flash before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your final flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W final flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W final flash!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 28, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your final flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W final flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W final flash, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 28, skill, SKILL_FINALFLASH)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 90 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 10
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your final flash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a final flash at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a final flash at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your final flash misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a final flash at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a final flash at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 28, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou grin as you bring either hand to the sides of your body. Your charged ki begins to pool there as @bblue@W orbs of energy form in your palms. You quickly bring both hands forward slamming your wrists together with the palms flat out facing @c$N@W! You shout '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing YOU! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing @c$N@W! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou grin as you bring either hand to the sides of your body. Your charged ki begins to pool there as @bblue@W orbs of energy form in your palms. You quickly bring both hands forward slamming your wrists together with the palms flat out facing @c$N@W! You shout '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing YOU! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing @c$N@W! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_POWERHIT]) == 0 {
					dmg *= 3
				} else {
					dmg *= 5
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou grin as you bring either hand to the sides of your body. Your charged ki begins to pool there as @bblue@W orbs of energy form in your palms. You quickly bring both hands forward slamming your wrists together with the palms flat out facing @c$N@W! You shout '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing YOU! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing @c$N@W! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou grin as you bring either hand to the sides of your body. Your charged ki begins to pool there as @bblue@W orbs of energy form in your palms. You quickly bring both hands forward slamming your wrists together with the palms flat out facing @c$N@W! You shout '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing YOU! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing @c$N@W! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou grin as you bring either hand to the sides of your body. Your charged ki begins to pool there as @bblue@W orbs of energy form in your palms. You quickly bring both hands forward slamming your wrists together with the palms flat out facing @c$N@W! You shout '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing YOU! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grins as $e brings either hand to the sides of $s body. @C$n@W's charged ki begins to pool there as @bblue@W orbs of energy form in $s palms. @C$n@W quickly brings both hands forward slamming $s wrists together with the palms flat out facing @c$N@W! @C$n@W shouts '@DF@ci@Cn@Da@cl @CF@Dl@ca@Cs@Dh@W!' as a massive wave of energy erupts all over @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 28, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a final flash at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a final flash at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_sbc(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.15
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SBC) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SBC)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SBC, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Special Beam Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Special Beam Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Special Beam Cannon before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your special beam cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W special beam cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W special beam cannon, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 27, skill, SKILL_SBC)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 90 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 10
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your special beam cannon misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a special beam cannon at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a special beam cannon at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your special beam cannon misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a special beam cannon at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a special beam cannon at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 27, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou place your index and middle fingers against your forehead and pool your charged ki there. Sparks and light surround your fingertips as the technique becomes ready. You point your fingers at @c$N@W and yell '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from your fingers and slams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at YOU and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at @c$N@W and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou place your index and middle fingers against your forehead and pool your charged ki there. Sparks and light surround your fingertips as the technique becomes ready. You point your fingers at @c$N@W and yell '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from your fingers and slams into $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at YOU and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at @c$N@W and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou place your index and middle fingers against your forehead and pool your charged ki there. Sparks and light surround your fingertips as the technique becomes ready. You point your fingers at @c$N@W and yell '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from your fingers and slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at YOU and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at @c$N@W and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou place your index and middle fingers against your forehead and pool your charged ki there. Sparks and light surround your fingertips as the technique becomes ready. You point your fingers at @c$N@W and yell '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from your fingers and slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at YOU and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at @c$N@W and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou place your index and middle fingers against your forehead and pool your charged ki there. Sparks and light surround your fingertips as the technique becomes ready. You point your fingers at @c$N@W and yell '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from your fingers and slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at YOU and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s index and middle fingers against $s forehead and pools $s charged ki there. Sparks and light surround the fingertips as the technique becomes ready. @C$n@W points $s fingers at @c$N@W and yells '@YS@yp@De@Yc@yi@Da@Yl @yB@De@Ya@ym @DC@Ya@yn@Dn@Yo@yn@W!' A spiraling beam of energy fires from the fingers and slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 27, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a special beam cannon at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a special beam cannon at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_tribeam(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.2
		minimum float64 = 0.1
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_TRIBEAM) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_TRIBEAM)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_TRIBEAM, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Tribeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Tribeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Tribeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your tribeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W tribeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W tribeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 26, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your tribeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W tribeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W tribeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 26, skill, SKILL_TRIBEAM)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 90 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 10
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your tribeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a tribeam at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a tribeam at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your tribeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a tribeam at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a tribeam at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 26, skill, attperc)
			if GET_SKILL(ch, SKILL_TRIBEAM) >= 100 {
				dmg += int64(float64(GET_LIFEMAX(ch)) * 0.2)
			} else if GET_SKILL(ch, SKILL_TRIBEAM) >= 60 {
				dmg += int64(float64(GET_LIFEMAX(ch)) * 0.1)
			} else if GET_SKILL(ch, SKILL_TRIBEAM) >= 40 {
				dmg += int64(float64(GET_LIFEMAX(ch)) * 0.05)
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou gather your charged ki as you form a triangle with your hands. You aim the gap between your hands at @c$N@W and shout '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from your hands and slams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at YOU and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at @c$N@W and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou gather your charged ki as you form a triangle with your hands. You aim the gap between your hands at @c$N@W and shout '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from your hands and slams into $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at YOU and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at @c$N@W and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou gather your charged ki as you form a triangle with your hands. You aim the gap between your hands at @c$N@W and shout '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from your hands and slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at YOU and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at @c$N@W and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou gather your charged ki as you form a triangle with your hands. You aim the gap between your hands at @c$N@W and shout '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from your hands and slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at YOU and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at @c$N@W and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou gather your charged ki as you form a triangle with your hands. You aim the gap between your hands at @c$N@W and shout '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from your hands and slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at YOU and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki as $e forms a triangle with $s hands. @C$n@W aims the gap between $s hands at @c$N@W and shouts '@rT@RR@YI@rB@RE@YA@rM@W!'. A bright blast of energy explodes from $s hands and slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !IS_NPC(vict) {
				WAIT_STATE(vict, (int(1000000/OPT_USEC))*2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 26, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a tribeam at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a tribeam at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kienzan(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KIENZAN) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KIENZAN)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KIENZAN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Kienzan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Kienzan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Kienzan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your kienzan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kienzan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W kienzan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wThe kienzan expands, cutting through everything in its path!\r\n"))
					dodge_ki(ch, vict, 2, 25, skill, SKILL_KIENZAN)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your kienzan misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a kienzan at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a kienzan at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 2, 25, skill, SKILL_KIENZAN)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your kienzan misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a kienzan at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a kienzan at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dodge_ki(ch, vict, 2, 25, skill, SKILL_KIENZAN)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 25, skill, attperc)
			if GET_SKILL(ch, SKILL_KIENZAN) >= 100 {
				dmg += int64(float64(dmg) * 0.25)
			} else if GET_SKILL(ch, SKILL_KIENZAN) >= 60 {
				dmg += int64(float64(dmg) * 0.15)
			} else if GET_SKILL(ch, SKILL_KIENZAN) >= 40 {
				dmg += int64(float64(dmg) * 0.05)
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise your hand above your head and pool your charged ki above your flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete you shout '@yK@Yi@we@yn@Yz@wa@yn@W!' and throw it! You watch as it slices into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if dmg > vict.Max_hit/5 && (int(vict.Race) != RACE_MAJIN && int(vict.Race) != RACE_BIO) {
					act(libc.CString("@R$N@r is cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou are cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r is cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					die(vict, ch)
					if AFF_FLAGGED(ch, AFF_GROUP) {
						group_gain(ch, vict)
					} else {
						solo_gain(ch, vict)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
						do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
						do_get(ch, libc.CString("all corpse"), 0, 0)
					}
				} else if dmg > vict.Max_hit/5 && (int(vict.Race) == RACE_MAJIN || int(vict.Race) == RACE_BIO) {
					if GET_SKILL(vict, SKILL_REGENERATE) > rand_number(1, 101) && vict.Mana >= vict.Max_mana/40 {
						act(libc.CString("@R$N@r is cut in half by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou are cut in half by the attack but regenerate a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r is cut in half by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Mana -= vict.Max_mana / 40
						hurt(0, 0, ch, vict, nil, dmg, 1)
					} else {
						act(libc.CString("@R$N@r is cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou are cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r is cut in half by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						die(vict, ch)
						if AFF_FLAGGED(ch, AFF_GROUP) {
							group_gain(ch, vict)
						} else {
							solo_gain(ch, vict)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
							do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
							do_get(ch, libc.CString("all corpse"), 0, 0)
						}
					}
				} else {
					hurt(0, 0, ch, vict, nil, dmg, 1)
					if (int(vict.Race) == RACE_ICER || int(vict.Race) == RACE_BIO) && PLR_FLAGGED(vict, PLR_TAIL) {
						act(libc.CString("@rYou cut off $S tail!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYour tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r's tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Act[int(PLR_TAIL/32)] &= bitvector_t(int32(^(1 << (int(PLR_TAIL % 32)))))
						remove_limb(vict, 6)
					}
					if (int(vict.Race) == RACE_SAIYAN || int(vict.Race) == RACE_HALFBREED) && PLR_FLAGGED(vict, PLR_STAIL) {
						act(libc.CString("@rYou cut off $S tail!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYour tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r's tail is cut off!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Act[int(PLR_STAIL/32)] &= bitvector_t(int32(^(1 << (int(PLR_STAIL % 32)))))
						remove_limb(vict, 5)
						if PLR_FLAGGED(vict, PLR_OOZARU) {
							act(libc.CString("@CYour body begins to shrink back to its normal form as the power of the Oozaru leaves you. You fall asleep shortly after returning to normal!@n"), TRUE, vict, nil, nil, TO_CHAR)
							act(libc.CString("@c$n@C's body begins to shrink and return to normal. Their giant ape features fading back into humanoid features until $e is left normal and asleep.@n"), TRUE, vict, nil, nil, TO_ROOM)
							vict.Act[int(PLR_OOZARU/32)] &= bitvector_t(int32(^(1 << (int(PLR_OOZARU % 32)))))
							vict.Hit = (vict.Hit / 2) - 10000
							vict.Mana = (vict.Mana / 2) - 10000
							vict.Move = (vict.Move / 2) - 10000
							vict.Max_hit = vict.Basepl
							vict.Max_mana = vict.Baseki
							vict.Max_move = vict.Basest
							if vict.Move < 1 {
								vict.Move = 1
							}
							if vict.Mana < 1 {
								vict.Mana = 1
							}
							if vict.Hit < 1 {
								vict.Hit = 1
							}
						}
					}
				}
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou raise your hand above your head and pool your charged ki above your flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete you shout '@yK@Yi@we@yn@Yz@wa@yn@W!' and throw it! You watch as it slices into @c$N@W's neck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into YOUR neck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into @c$N@W's neck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				if dmg > vict.Max_hit/5 && (int(vict.Race) != RACE_MAJIN && int(vict.Race) != RACE_BIO) {
					act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou have your head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Death_type = DTYPE_HEAD
					remove_limb(vict, 0)
					die(vict, ch)
					if AFF_FLAGGED(ch, AFF_GROUP) {
						group_gain(ch, vict)
					} else {
						solo_gain(ch, vict)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
						do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
					}
					if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
						do_get(ch, libc.CString("all corpse"), 0, 0)
					}
				} else if dmg > vict.Max_hit/5 && (int(vict.Race) == RACE_MAJIN || int(vict.Race) == RACE_BIO) {
					if GET_SKILL(vict, SKILL_REGENERATE) > rand_number(1, 101) && vict.Mana >= vict.Max_mana/40 {
						act(libc.CString("@R$N@r has $S head cut off by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou have your head cut off by the attack but regenerate a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r has $S head cut off by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Mana -= vict.Max_mana / 40
						hurt(0, 0, ch, vict, nil, dmg, 1)
					} else {
						act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@rYou have your head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N@r has $S head cut off by the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Death_type = DTYPE_HEAD
						die(vict, ch)
						if AFF_FLAGGED(ch, AFF_GROUP) {
							group_gain(ch, vict)
						} else {
							solo_gain(ch, vict)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
							do_get(ch, libc.CString("all.zenni corpse"), 0, 0)
						}
						if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
							do_get(ch, libc.CString("all corpse"), 0, 0)
						}
					}
				} else {
					hurt(0, 0, ch, vict, nil, dmg, 1)
				}
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou raise your hand above your head and pool your charged ki above your flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete you shout '@yK@Yi@we@yn@Yz@wa@yn@W!' and throw it! You watch as it slices into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou raise your hand above your head and pool your charged ki above your flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete you shout '@yK@Yi@we@yn@Yz@wa@yn@W!' and throw it! You watch as it slices into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 70 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[1]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severes $N's left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severes your left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left arm is severed in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Limb_condition[1] = 0
						remove_limb(vict, 2)
					} else if (vict.Limb_condition[0]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severes $N's right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severes your right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right arm is severed in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Limb_condition[0] = 0
						remove_limb(vict, 1)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou raise your hand above your head and pool your charged ki above your flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete you shout '@yK@Yi@we@yn@Yz@wa@yn@W!' and throw it! You watch as it slices into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand above $s head and pools $s charged ki above $s flattened palm. Slowly a golden spinning disk of energy grows from the ki. With the attack complete $e shouts '@yK@Yi@we@yn@Yz@wa@yn@W!' and throws it! @C$n@W watches as it slices into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 70 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[3]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severes $N's left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severes your left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left leg is severed in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Limb_condition[3] = 0
						remove_limb(vict, 4)
					} else if (vict.Limb_condition[2]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severes $N's right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severes your right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right leg is severed in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						vict.Limb_condition[2] = 0
						remove_limb(vict, 3)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 25, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a kienzan at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a kienzan at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_baku(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		dge     int = 2
		count   int = 0
		skill   int
		dmg     int64
		attperc float64    = 0.15
		minimum float64    = 0.05
		vict    *char_data = nil
		next_v  *char_data = nil
		arg2    [2048]byte
	)
	one_argument(argument, &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BAKUHATSUHA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
		send_to_char(ch, libc.CString("This room just has such a peaceful, easy feeling...\r\n"))
		return
	}
	skill = init_skill(ch, SKILL_BAKUHATSUHA)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		}
		if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
			continue
		}
		if AFF_FLAGGED(vict, AFF_GROUP) && (vict.Master == ch || ch.Master == vict) {
			continue
		}
		if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
			continue
		}
		if MOB_FLAGGED(vict, MOB_NOKILL) {
			continue
		} else {
			count += 1
		}
	}
	if count <= 0 {
		send_to_char(ch, libc.CString("There is no one worth targeting around.\r\n"))
		return
	} else {
		perc = chance_to_hit(ch)
		handle_cooldown(ch, 6)
		if skill < perc {
			act(libc.CString("@WYou raise your hand with index and middle fingers extended upwards. You try releasing your charged ki in a @yB@Ya@Wk@wuh@ya@Yt@Ws@wuh@ya@W but mess up and waste the ki!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W raises $s hand with index and middle fingers extended upwards. @C$n@W tries releasing $s charged ki in a @yB@Ya@Wk@wuh@ya@Yt@Ws@wuh@ya@W but messes up and wastes the ki!@n"), TRUE, ch, nil, nil, TO_ROOM)
			pcost(ch, attperc, 0)
			improve_skill(ch, SKILL_BAKUHATSUHA, 0)
			return
		}
		act(libc.CString("@WYou raise your hand with index and middle fingers extended upwards. A sudden burst rushes up from the ground as your charged ki explodes in the form of a @yB@Ya@Wk@wuh@ya@Yt@Ws@wuh@ya@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W raises $s hand with index and middle fingers extended upwards. A sudden burst rushes up from the ground as $s charged ki explodes in the form of a @yB@Ya@Wk@wuh@ya@Yt@Ws@wuh@ya@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
		dmg = damtype(ch, 24, skill, attperc)
		if GET_SKILL(ch, SKILL_BAKUHATSUHA) >= 100 {
			dmg += int64(float64(ch.Max_hit) * 0.08)
		} else if GET_SKILL(ch, SKILL_BAKUHATSUHA) >= 60 {
			dmg += int64(float64(ch.Max_hit) * 0.04)
		} else if GET_SKILL(ch, SKILL_BAKUHATSUHA) >= 40 {
			dmg += int64(float64(ch.Max_hit) * 0.02)
		}
		switch count {
		case 1:
			dmg = dmg
		case 2:
			dmg = (dmg / 100) * 75
		case 3:
			dmg = (dmg / 100) * 50
		default:
			dmg = (dmg / 100) * 25
		}
		for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
			next_v = vict.Next_in_room
			if vict == ch {
				continue
			}
			if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
				continue
			}
			if AFF_FLAGGED(vict, AFF_GROUP) && (vict.Master == ch || ch.Master == vict || vict.Master == ch.Master) {
				continue
			}
			if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
				continue
			}
			if MOB_FLAGGED(vict, MOB_NOKILL) {
				continue
			}
			dge = handle_dodge(vict)
			if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
				act(libc.CString("@C$N@c disappears, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding the explosion before reappearing elsewhere!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(vict, 0, vict.Max_hit/200)
				hurt(0, 0, ch, vict, nil, 0, 1)
				continue
			} else if dge+rand_number(-15, 5) > skill {
				act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@WYou manage to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 1)
				improve_skill(vict, SKILL_DODGE, 0)
				continue
			} else {
				act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@RYou are caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING {
					handle_knockdown(vict)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				continue
			}
		}
		pcost(ch, attperc, 0)
		improve_skill(ch, SKILL_BAKUHATSUHA, 0)
		return
	}
}
func do_rogafufuken(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.05
		minimum float64 = 0.01
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_ROGAFUFUKEN) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), ch.Max_hit/50) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_ROGAFUFUKEN)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_ROGAFUFUKEN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Rogafufuken before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Rogafufuken before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Rogafufuken before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your rogafufuken with a punch of $s own!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W rogafufuken with a punch of your own!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W rogafufuken with a punch of $s own!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, vict.Max_hit/300)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 23, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W rogafufuken!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your rogafufuken misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a rogafufuken at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a rogafufuken at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your rogafufuken misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a rogafufuken at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a rogafufuken at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 23, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou pour your charged energy into your hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as you leap towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. You unleash a flurry of hand strikes on @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards YOU while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou pour your charged energy into your hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as you leap towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. You unleash a flurry of hand strikes on @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards YOU while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou pour your charged energy into your hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as you leap towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. You unleash a flurry of hand strikes on @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards YOU while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou pour your charged energy into your hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as you leap towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. You unleash a flurry of hand strikes on @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards YOU while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou pour your charged energy into your hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as you leap towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. You unleash a flurry of hand strikes on @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards YOU while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W pours $s charged energy into $s hands and feet. A @rr@Re@wd@W glow trails behind the movements of either as $e leaps towards @c$N@W while yelling '@cW@Co@Wl@wf @DFang @rFist@W!'. $e unleashes a flurry of hand strikes on @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your hands and feet!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W hands and feet are burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W hands and feet are burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.04)
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 23, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a rogafufuken at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a rogafufuken at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_dualbeam(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.1
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_DUALBEAM) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_DUALBEAM)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_DUALBEAM, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Dualbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Dualbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Dualbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		var hits int = 3
		for hits > 0 {
			hits -= 1
			if hits == 1 {
				prob -= prob / 5
			} else if hits <= 0 {
				return
			} else if vict == nil {
				return
			} else if (AFF_FLAGGED(vict, AFF_SPIRIT) || vict.Hit <= 0) && !IS_NPC(vict) {
				return
			}
			if prob < perc-20 {
				if vict.Move > 0 {
					if blk > axion_dice(10) {
						act(libc.CString("@C$N@W moves quickly and blocks your dualbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@WYou move quickly and block @C$n's@W dualbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W dualbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						if hits == 1 {
							pcost(ch, attperc, 0)
						}
						pcost(vict, 0, vict.Max_hit/500)
						dmg = damtype(ch, 22, skill, attperc)
						dmg /= 4
						hurt(0, 0, ch, vict, nil, dmg, 1)
						continue
					} else if dge > axion_dice(10) {
						act(libc.CString("@C$N@W manages to dodge your dualbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@WYou dodge @C$n's@W dualbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@C$N@W manages to dodge @c$n's@W dualbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
						dodge_ki(ch, vict, 0, 22, skill, SKILL_DUALBEAM)
						if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
							(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
						}
						improve_skill(vict, SKILL_DODGE, 0)
						if hits == 1 {
							pcost(ch, attperc, 0)
						}
						hurt(0, 0, ch, vict, nil, 0, 1)
						continue
					} else {
						act(libc.CString("@WYou can't believe it but your dualbeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@C$n@W fires a dualbeam at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@c$n@W fires a dualbeam at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
						if hits == 1 {
							pcost(ch, attperc, 0)
						}
						hurt(0, 0, ch, vict, nil, 0, 1)
						continue
					}
				} else {
					act(libc.CString("@WYou can't believe it but your dualbeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a dualbeam at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a dualbeam at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if hits == 1 {
						pcost(ch, attperc, 0)
					}
				}
				hurt(0, 0, ch, vict, nil, 0, 1)
				continue
			} else {
				dmg = damtype(ch, 22, skill, attperc)
				var hitspot int = 1
				hitspot = roll_hitloc(ch, vict, skill)
				if GET_SKILL(ch, SKILL_DUALBEAM) >= 100 {
					if rand_number(1, 100) >= 60 {
						hitspot = 2
					}
				} else if GET_SKILL(ch, SKILL_DUALBEAM) >= 60 {
					if rand_number(1, 100) >= 80 {
						hitspot = 2
					}
				} else if GET_SKILL(ch, SKILL_DUALBEAM) >= 40 {
					if rand_number(1, 100) >= 95 {
						hitspot = 2
					}
				}
				switch hitspot {
				case 1:
					act(libc.CString("@WYou gather your charged energy up through the circuits in your arms. A @gg@Gr@We@wen @Wglow appears around your hand right as you aim it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S body in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at YOU. A @gg@Gr@We@wen @Wbeam blasts out and slams into YOUR body in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S body in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (ch.Bonuses[BONUS_SOFT]) != 0 {
						dmg *= int64(calc_critical(ch, 2))
					}
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou gather your charged energy up through the circuits in your arms. A @gg@Gr@We@wen @Wglow appears around your hand right as you aim it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S face in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at YOU. A @gg@Gr@We@wen @Wbeam blasts out and slams into YOUR face in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S face in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 3)
				case 3:
					act(libc.CString("@WYou gather your charged energy up through the circuits in your arms. A @gg@Gr@We@wen @Wglow appears around your hand right as you aim it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S gut in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at YOU. A @gg@Gr@We@wen @Wbeam blasts out and slams into YOUR gut in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S gut in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (ch.Bonuses[BONUS_SOFT]) != 0 {
						dmg *= int64(calc_critical(ch, 2))
					}
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 4)
				case 4:
					act(libc.CString("@WYou gather your charged energy up through the circuits in your arms. A @gg@Gr@We@wen @Wglow appears around your hand right as you aim it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S arm in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at YOU. A @gg@Gr@We@wen @Wbeam blasts out and slams into YOUR arm in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S arm in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 190, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 1)
				case 5:
					act(libc.CString("@WYou gather your charged energy up through the circuits in your arms. A @gg@Gr@We@wen @Wglow appears around your hand right as you aim it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S leg in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at YOU. A @gg@Gr@We@wen @Wbeam blasts out and slams into YOUR leg in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W gathers $s charged energy up through the circuits in $s arms. A @gg@Gr@We@wen @Wglow appears around $s hand right as $e aims it at @c$N@W. A @gg@Gr@We@wen @Wbeam blasts out and slams into $S leg in that instant!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(1, 190, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 2)
				}
				if hits == 1 {
					pcost(ch, attperc, 0)
				}
				if vict.Hit <= 0 {
					if hits == 2 {
						pcost(ch, attperc, 0)
					}
					hits = 0
				}
				continue
			}
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 22, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a dualbeam at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a dualbeam at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_blessedhammer(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.05
		minimum float64 = 0.01
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BLESSEDHAMMER) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_BLESSEDHAMMER)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_BLESSEDHAMMER, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob += 15
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@C before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@C before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@C before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@n!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@n!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@n!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 42, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 17, skill, SKILL_BLESSEDHAMMER)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr@W at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 42, skill, attperc)
			if AFF_FLAGGED(vict, AFF_SANCTUARY) {
				dmg *= int64(calc_critical(ch, 1))
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WFocusing your attention on @R$N@W, you reach back and form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer before hurling it with all your might into their chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W narrows $s eyes and focuses $s energy. They reach back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at you with all their @Rmight@W into YOUR chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W narrows their eyes at @c$N@W. $n reaches back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at @c$N@W with all their might into @c$N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				if (ch.Bonuses[BONUS_POWERHIT]) == 0 {
					dmg *= 3
				} else {
					dmg *= 5
				}
				act(libc.CString("@WFocusing your attention on @R$N@W, you reach back and form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer before hurling it with all your might into their face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W narrows $s eyes and focuses $s energy. They reach back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at you with all their @Rmight@W into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W narrows their eyes at @c$N@W. $n reaches back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at @c$N@W with all their might into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WFocusing your attention on @R$N@W, you reach back and form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer before hurling it with all your might into their gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W narrows $s eyes and focuses $s energy. They reach back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at you with all their @Rmight@W into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W narrows their eyes at @c$N@W. $n reaches back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at @c$N@W with all their might into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WFocusing your attention on @R$N@W, you reach back and form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer before hurling it with all your might into their arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W narrows $s eyes and focuses $s energy. They reach back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at you with all their @Rmight@W into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W narrows their eyes at @c$N@W. $n reaches back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at @c$N@W with all their might into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WFocusing your attention on @R$N@W, you reach back and form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer before hurling it with all your might into their leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W narrows $s eyes and focuses $s energy. They reach back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at you with all their @Rmight@W into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W narrows their eyes at @c$N@W. $n reaches back to form an @ne@Dt@nh@De@nr@De@na@Dl @Whammer which they hurl at @c$N@W with all their might into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_BURNED) && rand_number(1, 4) == 3 && int(vict.Race) != RACE_DEMON {
				send_to_char(vict, libc.CString("@RYou are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are burned by the attack!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 42, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a @WB@Dl@We@Ds@Ws@De@Wd @DH@Wa@Dm@Wm@De@Wr at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kousengan(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.05
		minimum float64 = 0.01
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KOUSENGAN) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KOUSENGAN)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KOUSENGAN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob += 15
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your kousengan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c kousengan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c kousengan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your kousengan!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W kousengan!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W kousengan!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 42, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your kousengan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kousengan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W kousengan, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 17, skill, SKILL_KOUSENGAN)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your kousengan misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a kousengan at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a kousengan at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your kousengan misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a kousengan at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a kousengan at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 42, skill, attperc)
			if AFF_FLAGGED(vict, AFF_SANCTUARY) {
				dmg *= int64(calc_critical(ch, 1))
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou look at @c$N@W and grin. Then bright @Mpink@W lasers shoot from your eyes and slam into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W looks at YOU and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into YOUR chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W looks at @c$N@W and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into @c$N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				if (ch.Bonuses[BONUS_POWERHIT]) == 0 {
					dmg *= 3
				} else {
					dmg *= 5
				}
				act(libc.CString("@WYou look at @c$N@W and grin. Then bright @Mpink@W lasers shoot from your eyes and slam into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W looks at YOU and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W looks at @c$N@W and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou look at @c$N@W and grin. Then bright @Mpink@W lasers shoot from your eyes and slam into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W looks at YOU and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W looks at @c$N@W and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou look at @c$N@W and grin. Then bright @Mpink@W lasers shoot from your eyes and slam into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W looks at YOU and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W looks at @c$N@W and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou look at @c$N@W and grin. Then bright @Mpink@W lasers shoot from your eyes and slam into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W looks at YOU and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W looks at @c$N@W and grins. Then bright @Mpink@W lasers shoot from $s eyes and slam into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if !AFF_FLAGGED(vict, AFF_BURNED) && rand_number(1, 4) == 3 && int(vict.Race) != RACE_DEMON {
				send_to_char(vict, libc.CString("@RYou are burned by the attack!@n\r\n"))
				send_to_char(ch, libc.CString("@RThey are burned by the attack!@n\r\n"))
				vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 42, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a kousengan at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a kousengan at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_deathbeam(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.1
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_DEATHBEAM) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_DEATHBEAM)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_DEATHBEAM, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob += 15
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Deathbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Deathbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Deathbeam before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your deathbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W deathbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W deathbeam!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 17, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your deathbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W deathbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W deathbeam, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 17, skill, SKILL_DEATHBEAM)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					return
				} else {
					act(libc.CString("@WYou can't believe it but your deathbeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a deathbeam at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a deathbeam at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your deathbeam misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a deathbeam at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a deathbeam at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			}
			return
		} else {
			dmg = damtype(ch, 17, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou move swiftly, drawing your charged ki to your index finger, and point at @c$N@W! You fire a @Rred@W Deathbeam from your finger which slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and point at you@W! $e fires a @Rred@W Deathbeam from $s finger which slams into your body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and point at @c$N@W! $e fires a @Rred@W Deathbeam from $s finger which slams into @c$N@W's body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			case 2:
				act(libc.CString("@WYou move swiftly, drawing your charged ki to your index finger, and point at @c$N@W! You fire a @Rred@W Deathbeam from your finger which slams into $S face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at you@W! $e fires a @Rred@W Deathbeam from $s finger which slams into your face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at @c$N@W! $e fires a @Rred@W Deathbeam from $s finger which slams into @c$N@W's face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			case 3:
				act(libc.CString("@WYou move swiftly, drawing your charged ki to your index finger, and point at @c$N@W! You fire a @Rred@W Deathbeam from your finger which slams into $S gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at you@W! $e fires a @Rred@W Deathbeam from $s finger which slams into your gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at @c$N@W! $e fires a @Rred@W Deathbeam from $s finger which slams into @c$N@W's gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			case 4:
				act(libc.CString("@WYou move swiftly, drawing your charged ki to your index finger, and point at @c$N@W! You fire a @Rred@W Deathbeam from your finger which slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at you@W! $e fires a @Rred@W Deathbeam from $s finger which slams into your arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at @c$N@W! $e fires a @Rred@W Deathbeam from $s finger which slams into @c$N@W's arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			case 5:
				act(libc.CString("@WYou move swiftly, drawing your charged ki to your index finger, and point at @c$N@W! You fire a @Rred@W Deathbeam from your finger which slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at you@W! $e fires a @Rred@W Deathbeam from $s finger which slams into your leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W moves swiftly, drawing charged ki to $s index finger, and points at @c$N@W! $e fires a @Rred@W Deathbeam from $s finger which slams into @c$N@W's leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
				if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			}
			if GET_SKILL(ch, SKILL_DEATHBEAM) >= 100 && vict.Hit >= 2 {
				vict.Lifeforce -= int64(float64(dmg) * 0.4)
				if vict.Lifeforce < -1 {
					vict.Lifeforce = -1
				}
			} else if GET_SKILL(ch, SKILL_DEATHBEAM) >= 60 && vict.Hit >= 2 {
				vict.Lifeforce -= int64(float64(dmg) * 0.2)
				if vict.Lifeforce < -1 {
					vict.Lifeforce = -1
				}
			} else if GET_SKILL(ch, SKILL_DEATHBEAM) >= 40 && vict.Hit >= 2 {
				vict.Lifeforce -= int64(float64(dmg) * 0.05)
				if vict.Lifeforce < -1 {
					vict.Lifeforce = -1
				}
			}
			if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 17, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a deathbeam at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a deathbeam at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
		if int(ch.Skillperfs[SKILL_DEATHBEAM]) == 3 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_dodonpa(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.1
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_DODONPA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_DODONPA]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_DODONPA]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_DODONPA)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_DODONPA, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_DODONPA]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if int(ch.Skillperfs[SKILL_DODONPA]) == 3 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Dodonpa before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Dodonpa before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Dodonpa before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your dodonpa!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W dodonpa!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W dodonpa!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 15, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your dodonpa, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W dodonpa, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W dodonpa, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 15, skill, SKILL_DODONPA)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your dodonpa misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a dodonpa at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a dodonpa at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your dodonpa misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a dodonpa at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a dodonpa at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 15, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou gather your charged ki into your index finger and point at @c$N@W. A @Ygolden@W glow forms around your finger tip right before you shout '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fire a large golden beam! The beam slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at you. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into your body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at @c$N@W. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into @c$N@W's body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou gather your charged ki into your index finger and point at @c$N@W. A @Ygolden@W glow forms around your finger tip right before you shout '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fire a large golden beam! The beam slams into $S face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at you. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into your face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at @c$N@W. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into @c$N@W's face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou gather your charged ki into your index finger and point at @c$N@W. A @Ygolden@W glow forms around your finger tip right before you shout '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fire a large golden beam! The beam slams into $S gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at you. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into your gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at @c$N@W. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into @c$N@W's gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou gather your charged ki into your index finger and point at @c$N@W. A @Ygolden@W glow forms around your finger tip right before you shout '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fire a large golden beam! The beam slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at you. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into your arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at @c$N@W. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into @c$N@W's arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou gather your charged ki into your index finger and point at @c$N@W. A @Ygolden@W glow forms around your finger tip right before you shout '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fire a large golden beam! The beam slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at you. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into your leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W gathers $s charged ki into $s index finger and points at @c$N@W. A @Ygolden@W glow forms around the finger tip right before $e shouts '@YD@yo@Yd@yo@Yn@yp@Ya@W!' and fires a large golden beam! The beam slams into @c$N@W's leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if rand_number(1, 3) == 2 {
				vict.Mana -= dmg / 4
				if vict.Mana < 0 {
					vict.Mana = 0
				}
				send_to_char(vict, libc.CString("@RYou feel some of your ki drained away by the attack!@n\r\n"))
			}
			if int(ch.Skillperfs[SKILL_DODONPA]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 15, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a dodonpa at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a dodonpa at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_masenko(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.15
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_MASENKO) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_MASENKO]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_MASENKO]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_MASENKO)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_MASENKO, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_MASENKO]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Masenko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Masenko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Masenko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your masenko!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W masenko!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W masenko!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 14, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your masenko, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W masenko, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W masenko, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 14, skill, SKILL_MASENKO)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your masenko misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a masenko at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a masenko at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your masenko misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a masenko at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a masenko at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 14, skill, attperc)
			if GET_SKILL(ch, SKILL_MASENKO) >= 100 {
				dmg += int64(float64(dmg) * 0.08)
			} else if GET_SKILL(ch, SKILL_MASENKO) >= 100 {
				dmg += int64(float64(dmg) * 0.05)
			} else if GET_SKILL(ch, SKILL_MASENKO) >= 100 {
				dmg += int64(float64(dmg) * 0.03)
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			if int(ch.Chclass) == CLASS_PICCOLO {
				if int(ch.Skills[SKILL_STYLE]) >= 100 {
					dmg += int64(float64(ch.Max_mana) * 0.08)
				} else if int(ch.Skills[SKILL_STYLE]) >= 60 {
					dmg += int64(float64(ch.Max_mana) * 0.05)
				} else if int(ch.Skills[SKILL_STYLE]) >= 40 {
					dmg += int64(float64(ch.Max_mana) * 0.03)
				}
			}
			switch hitspot {
			case 1:
				act(libc.CString("@WYou place your hands on top of each other and raise them over your head. Energy that you have charged gathers into your palms just before you bring them down aimed at @c$N@W! You scream '@RMasenko @rHa!@W' as you release your attack. A bright red wave of energy flows from your hands and slams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at you! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at @c$N@W! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou place your hands on top of each other and raise them over your head. Energy that you have charged gathers into your palms just before you bring them down aimed at @c$N@W! You scream '@RMasenko @rHa!@W' as you release your attack. A bright red wave of energy flows from your hands and slams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at you! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at @c$N@W! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into @c$N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou place your hands on top of each other and raise them over your head. Energy that you have charged gathers into your palms just before you bring them down aimed at @c$N@W! You scream '@RMasenko @rHa!@W' as you release your attack. A bright red wave of energy flows from your hands and slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at you! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at @c$N@W! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou place your hands on top of each other and raise them over your head. Energy that you have charged gathers into your palms just before you bring them down aimed at @c$N@W! You scream '@RMasenko @rHa!@W' as you release your attack. A bright red wave of energy flows from your hands and slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at you! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at @c$N@W! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou place your hands on top of each other and raise them over your head. Energy that you have charged gathers into your palms just before you bring them down aimed at @c$N@W! You scream '@RMasenko @rHa!@W' as you release your attack. A bright red wave of energy flows from your hands and slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at you! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W places $s hands on top of each other and raises them over $s head. Energy that $e had charged gathers into $s palms just before $e brings them down aimed at @c$N@W! $e screams '@RMasenko @rHa!@W' as $e releases $s attack. A bright red wave of energy flows from $s hands and slams into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if rand_number(1, 2) == 2 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
				send_to_char(vict, libc.CString("@RThe attack seems to have taken a toll on your stamina!@n\r\n"))
				vict.Move -= dmg / 4
				if vict.Move < 0 {
					vict.Move = 0
				}
			}
			if int(ch.Skillperfs[SKILL_MASENKO]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 14, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a masenko at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a masenko at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kamehameha(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.15
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KAMEHAMEHA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KAMEHAMEHA)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KAMEHAMEHA, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 2 {
			prob += 5
		}
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Kamehameha Wave before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Kamehameha Wave before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Kamehameha Wave before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
				} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your kamehameha!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W kamehameha!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W kamehameha!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					pcost(vict, 0, vict.Max_hit/500)
					if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					dmg = damtype(ch, 13, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your kamehameha, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kamehameha, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W kamehameha, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 13, skill, SKILL_KAMEHAMEHA)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your kamehameha misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a kamehameha at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a kamehameha at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your kamehameha misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a kamehameha at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a kamehameha at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
				} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
			}
			if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 13, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou cup your hands at your side and begin to pool your charged ki there. As the ki pools there you begin to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly you bring your hands forward facing @c$N@W and shout '@BHAAAA!!!@W' while releasing a bright blue kamehameha at $M! It slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing you and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into your body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing @c$N@W and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou cup your hands at your side and begin to pool your charged ki there. As the ki pools there you begin to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly you bring your hands forward facing @c$N@W and shout '@BHAAAA!!!@W' while releasing a bright blue kamehameha at $M! It slams into $S face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing you and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into your face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing @c$N@W and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into $S face and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou cup your hands at your side and begin to pool your charged ki there. As the ki pools there you begin to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly you bring your hands forward facing @c$N@W and shout '@BHAAAA!!!@W' while releasing a bright blue kamehameha at $M! It slams into $S gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing you and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into your gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing @c$N@W and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into $S gut and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			case 4:
				act(libc.CString("@WYou cup your hands at your side and begin to pool your charged ki there. As the ki pools there you begin to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly you bring your hands forward facing @c$N@W and shout '@BHAAAA!!!@W' while releasing a bright blue kamehameha at $M! It slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing you and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into your arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing @c$N@W and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou cup your hands at your side and begin to pool your charged ki there. As the ki pools there you begin to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly you bring your hands forward facing @c$N@W and shout '@BHAAAA!!!@W' while releasing a bright blue kamehameha at $M! It slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing you and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into your leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W cups $s hands at $s side and begins to pool charged ki there. As the ki pools there $e begins to chant, '@BKaaaaa@bmeeee@Bhaaaaaaa@bmeeeee@W'. Suddenly $e brings $s hands forward facing @c$N@W and shouts '@BHAAAA!!!@W' while releasing a bright blue kamehameha! It slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
				if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
			}
			if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
			} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
			} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 13, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a kamehameha at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a kamehameha at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		if int(ch.Skillperfs[SKILL_KAMEHAMEHA]) == 3 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
		if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 100 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.25)
		} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 60 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
		} else if GET_SKILL(ch, SKILL_KAMEHAMEHA) >= 40 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_renzo(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_RENZO) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_RENZO)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_RENZO, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		var master_roll int = rand_number(1, 100)
		var master_chance int = 0
		var half_chance int = 0
		var master_pass int = 0
		if skill >= 100 {
			master_chance = 10
			half_chance = 20
		} else if skill >= 75 {
			master_chance = 5
			half_chance = 10
		} else if skill >= 50 {
			master_chance = 5
			half_chance = 5
		}
		if master_chance >= master_roll {
			master_pass = 1
		} else if half_chance >= master_roll {
			master_pass = 2
		}
		if master_pass == 1 {
			send_to_char(ch, libc.CString("@GYour mastery of the technique has made your use of energy more efficient!@n\r\n"))
		} else if master_pass == 2 {
			send_to_char(ch, libc.CString("@GYour mastery of the technique has made your use of energy as efficient as possible!@n\r\n"))
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Renzokou Energy Dan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Renzokou Energy Dan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Renzokou Energy Dan before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if master_pass == 1 {
					pcost(ch, attperc*0.25, 0)
				} else if master_pass == 2 {
					pcost(ch, attperc*0.5, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		var count int = 0
		if prob+20 < perc {
			count = 0
		} else if prob+15 < perc {
			count = 10
		} else if prob+10 < perc {
			count = 25
		} else if prob+5 < perc {
			count = 50
		} else if prob < perc-20 {
			count = 75
		} else if prob > perc {
			count = 100
		}
		if count > 0 {
			if int(ch.Chclass) == CLASS_NAIL {
				if GET_SKILL(ch, SKILL_RENZO) >= 100 {
					count += 200
				} else if GET_SKILL(ch, SKILL_RENZO) >= 60 {
					count += 100
				} else if GET_SKILL(ch, SKILL_RENZO) >= 40 {
					count += 40
				}
			}
			if rand_number(1, 5) >= 5 {
				count += rand_number(-15, 25)
			}
		}
		if count == 0 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects every shot of your renzokou energy dan, sending them flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect every shot @C$n's@W renzokou energy dan sending them flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects every shot of @c$n's@W renzokou energy dan sending them flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if master_pass == 1 {
						pcost(ch, attperc*0.25, 0)
					} else if master_pass == 2 {
						pcost(ch, attperc*0.5, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/200)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks every renzokou energy dan shot with $S arms!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W renzokou energy dan!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W renzokou energy dan with $S arms!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if master_pass == 1 {
						pcost(ch, attperc*0.25, 0)
					} else if master_pass == 2 {
						pcost(ch, attperc*0.5, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 12, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge all your renzokou energy dan shots, letting them slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge all of @C$n's @Wrenzokou energy dan shots, letting them slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge all @c$n's@W renzokou energy dan shots, letting them slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if master_pass == 1 {
						pcost(ch, attperc*0.25, 0)
					} else if master_pass == 2 {
						pcost(ch, attperc*0.5, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but all your renzokou energy dan shots miss, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires hundreds of renzokou energy dan shots at you, but misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires hundreds of renzokou energy dan shots at @C$N@W, but somehow misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if master_pass == 1 {
						pcost(ch, attperc*0.25, 0)
					} else if master_pass == 2 {
						pcost(ch, attperc*0.5, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but all your renzokou energy dan shots miss, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires hundreds of renzokou energy dan shots at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires hundreds of renzokou energy dan shots at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if master_pass == 1 {
					pcost(ch, attperc*0.25, 0)
				} else if master_pass == 2 {
					pcost(ch, attperc*0.5, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 12, skill, attperc)
			if count >= 100 {
				dmg = int64(float64(dmg) * 0.01)
				dmg *= int64(count)
				act(libc.CString("@WYou gather your charged energy into your hands as a golden glow appears around each. You slam your hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! All of the shots hit!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at you! All of the shots hit!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! All of the shots hit!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			if count >= 75 && count < 100 {
				dmg = int64(float64(dmg) * 0.01)
				dmg *= int64(count)
				act(libc.CString("@WYou gather your charged energy into your hands as a golden glow appears around each. You slam your hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Most of the shots hit, but some of them are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at you! Most of the shots hit, but some of them you manage to avoid!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Most of the shots hit, but some of them are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			if count >= 50 && count < 75 {
				dmg = int64(float64(dmg) * 0.01)
				dmg *= int64(count)
				act(libc.CString("@WYou gather your charged energy into your hands as a golden glow appears around each. You slam your hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! About half of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at you! About half of the shots hit, the rest you manage to avoid!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! About half of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			if count >= 25 && count < 50 {
				dmg = int64(float64(dmg) * 0.01)
				dmg *= int64(count)
				act(libc.CString("@WYou gather your charged energy into your hands as a golden glow appears around each. You slam your hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Few of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at you! Few of the shots hit, the rest you manage to avoid!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Few of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			if count >= 10 && count < 25 {
				dmg /= 100
				dmg *= int64(count)
				act(libc.CString("@WYou gather your charged energy into your hands as a golden glow appears around each. You slam your hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Very few of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at you! Very few of the shots hit, the rest you manage to avoid!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@w$n gathers charged energy into $s hands as a golden glow appears around each. $e slams $s hands forward rapidly firing hundreds of Renzokou Energy Dan shots at $N@W! Very few of the shots hit, the rest are avoided by $M!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			if master_pass == 1 {
				pcost(ch, attperc*0.25, 0)
			} else if master_pass == 2 {
				pcost(ch, attperc*0.5, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 12, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire hundreds of renzokou energy dan shots at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires hundreds of renzokou energy dan shots at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_heeldrop(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 0
		dge     int = 0
		blk     int = 0
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_HEELDROP)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if check_skill(ch, SKILL_HEELDROP) == 0 {
		return
	}
	if can_grav(ch) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/90) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_HEELDROP)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Chclass) == CLASS_PICCOLO {
		if int(ch.Skills[SKILL_STYLE]) >= 75 {
			handle_cooldown(ch, 7)
		} else {
			handle_cooldown(ch, 9)
		}
	} else {
		handle_cooldown(ch, 9)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_HEELDROP, 1)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		if !IS_NPC(vict) {
			pry = handle_parry(vict)
			blk = GET_SKILL(vict, SKILL_BLOCK)
			dge = GET_SKILL(vict, SKILL_DODGE)
		}
		if IS_NPC(vict) && GET_LEVEL(ch) <= 10 {
			if IS_HUMANOID(vict) {
				pry = rand_number(20, 40)
				blk = rand_number(20, 40)
			}
			dge = rand_number(20, 40)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 20 {
			if IS_HUMANOID(vict) {
				pry = rand_number(20, 60)
				blk = rand_number(20, 60)
			}
			dge = rand_number(20, 60)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 30 {
			if IS_HUMANOID(vict) {
				pry = rand_number(20, 70)
				blk = rand_number(20, 70)
			}
			dge = rand_number(20, 70)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 40 {
			if IS_HUMANOID(vict) {
				pry = rand_number(30, 70)
				blk = rand_number(30, 70)
			}
			dge = rand_number(30, 70)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 40 {
			if IS_HUMANOID(vict) {
				pry = rand_number(40, 70)
				blk = rand_number(40, 70)
			}
			dge = rand_number(40, 70)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 60 {
			if IS_HUMANOID(vict) {
				pry = rand_number(40, 80)
				blk = rand_number(40, 80)
			}
			dge = rand_number(40, 80)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 70 {
			if IS_HUMANOID(vict) {
				pry = rand_number(50, 80)
				blk = rand_number(50, 80)
			}
			dge = rand_number(50, 80)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 80 {
			if IS_HUMANOID(vict) {
				pry = rand_number(50, 90)
				blk = rand_number(50, 90)
			}
			dge = rand_number(50, 90)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 90 {
			if IS_HUMANOID(vict) {
				pry = rand_number(60, 90)
				blk = rand_number(60, 90)
			}
			dge = rand_number(60, 90)
		} else if IS_NPC(vict) && GET_LEVEL(ch) <= 100 {
			if IS_HUMANOID(vict) {
				pry = rand_number(70, 100)
				blk = rand_number(70, 100)
			}
			dge = rand_number(70, 100)
		}
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Heeldrop before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Heeldrop before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Heeldrop before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your heeldrop with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W heeldrop with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W heeldrop with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 8, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W heeldrop!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your heeldrop misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W moves to heeldrop you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to heeldrop @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your heeldrop misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves to heeldrop you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to heeldrop @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 8, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou disappear, appearing above @C$N@W you spin and heeldrop $M in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W disappears, only to appear above you, spinning quickly and heeldropping you in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W disappears, only to appear above @C$N@W, spinning quickly and heeldropping $M in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(vict, AFF_FLYING) {
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@wYou are knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
					vict.Altitude = 0
					vict.Position = POS_SITTING
				} else {
					handle_knockdown(vict)
				}
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou disappear, reappearing in front of @C$N@W, you flip upside down and slam your heel into the top of $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W disappears, reappearing in front of you, $e flips upside down and slams $s heel into the top of your head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W disappears, reappearing in front of @C$N@W, $e flips upside down and slams $s heel into the top of @C$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(vict, AFF_FLYING) {
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@wYou are knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
					vict.Altitude = 0
					vict.Position = POS_SITTING
				} else {
					handle_knockdown(vict)
				}
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou fly at @C$N@W, heeldropping $S gut as you fly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W flies at you, heeldropping your gut as $e flies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W flies at @C$N@W, heeldropping $S gut as $e flies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou heeldrop @C$N@W, hitting $M in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W heeldrops you, hitting you in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W heeldrops @C$N@W, hitting $M in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 180, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou heeldrop @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W heeldrops your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W heeldrops @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 180, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W leg is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W leg is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou heeldrop $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W heeldrops $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_attack(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob      int
		perc      int
		avo       int
		index     int   = 0
		pry       int   = 0
		dge       int   = 0
		blk       int   = 0
		skill     int   = 0
		wtype     int   = 0
		gun       int   = FALSE
		gun2      int   = FALSE
		dualwield int   = 0
		wielded   int   = 0
		guncost   int   = 0
		stcost    int64 = (ch.Max_hit / 150)
		dmg       int64
		vict      *char_data
		obj       *obj_data = nil
		arg       [2048]byte
		attperc   float64 = 0
	)
	if int(ch.Race) == RACE_ANDROID {
		stcost *= int64(0.25)
	}
	one_argument(argument, &arg[0])
	if (ch.Equipment[WEAR_WIELD1]) == nil && (ch.Equipment[WEAR_WIELD2]) == nil {
		send_to_char(ch, libc.CString("You need to wield a weapon to use this, without one try punch, kick, or other no weapon attacks.\r\n"))
		return
	}
	if can_grav(ch) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil {
		if int(ch.Race) != RACE_ANDROID {
			stcost += (ch.Equipment[WEAR_WIELD1]).Weight
		} else {
			stcost += int64(float64((ch.Equipment[WEAR_WIELD1]).Weight) * 0.25)
		}
		if check_points(ch, 0, stcost) == 0 {
			return
		}
		wielded = 1
	} else if (ch.Equipment[WEAR_WIELD2]) != nil {
		if int(ch.Race) != RACE_ANDROID {
			stcost += (ch.Equipment[WEAR_WIELD2]).Weight
		} else {
			stcost += int64(float64((ch.Equipment[WEAR_WIELD2]).Weight) * 0.25)
		}
		if check_points(ch, 0, stcost) == 0 {
			return
		}
	}
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil {
		if vict != nil {
			if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
				if can_kill(ch, vict, nil, 1) == 0 {
					return
				}
			} else {
				if can_kill(ch, vict, nil, 0) == 0 {
					return
				}
			}
		}
		if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
			skill = init_skill(ch, SKILL_DAGGER)
			improve_skill(ch, SKILL_DAGGER, 1)
			wtype = 1
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SWORD)
			improve_skill(ch, SKILL_SWORD, 1)
			wtype = 0
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_CRUSH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_CLUB)
			improve_skill(ch, SKILL_CLUB, 1)
			wtype = 2
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SPEAR)
			improve_skill(ch, SKILL_SPEAR, 1)
			wtype = 3
		} else if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
			gun = TRUE
			skill = init_skill(ch, SKILL_GUN)
			improve_skill(ch, SKILL_GUN, 1)
			wtype = 4
		} else {
			skill = init_skill(ch, SKILL_BRAWL)
			improve_skill(ch, SKILL_BRAWL, 1)
			wtype = 5
		}
	} else if (ch.Equipment[WEAR_WIELD2]) != nil {
		if vict != nil {
			if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
				if can_kill(ch, vict, nil, 1) == 0 {
					return
				}
			} else {
				if can_kill(ch, vict, nil, 0) == 0 {
					return
				}
			}
		}
		if wielded == 1 {
			wielded = 2
		}
		if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
			skill = init_skill(ch, SKILL_DAGGER)
			improve_skill(ch, SKILL_DAGGER, 1)
			wtype = 1
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SWORD)
			improve_skill(ch, SKILL_SWORD, 1)
			wtype = 0
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_CRUSH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_CLUB)
			improve_skill(ch, SKILL_CLUB, 1)
			wtype = 2
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SPEAR)
			improve_skill(ch, SKILL_SPEAR, 1)
			wtype = 3
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
			gun2 = TRUE
			skill = init_skill(ch, SKILL_GUN)
			improve_skill(ch, SKILL_GUN, 1)
			wtype = 4
		} else {
			skill = init_skill(ch, SKILL_BRAWL)
			improve_skill(ch, SKILL_BRAWL, 1)
			wtype = 5
		}
	}
	if wielded == 2 && gun == FALSE {
		if int(ch.Skills[SKILL_DUALWIELD]) >= 100 {
			dualwield = 3
			stcost -= int64(float64(stcost) * 0.3)
		} else if int(ch.Skills[SKILL_DUALWIELD]) >= 75 {
			dualwield = 2
			stcost -= int64(float64(stcost) * 0.25)
		}
	}
	var wlvl int = 0
	var weap *obj_data = nil
	if (ch.Equipment[WEAR_WIELD1]) != nil {
		weap = ch.Equipment[WEAR_WIELD1]
	} else {
		weap = ch.Equipment[WEAR_WIELD2]
	}
	if OBJ_FLAGGED(weap, ITEM_WEAPLVL1) {
		wlvl = 1
	} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL2) {
		wlvl = 2
	} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL3) {
		wlvl = 3
	} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL4) {
		wlvl = 4
	} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL5) {
		wlvl = 5
	}
	if ch.Preference != PREFERENCE_H2H {
		handle_cooldown(ch, 4)
	} else {
		handle_cooldown(ch, 8)
	}
	if wielded == 1 && (gun == TRUE || gun2 == TRUE) {
		if wlvl == 5 {
			guncost = 12
		} else if wlvl == 4 {
			guncost = 6
		} else if wlvl == 3 {
			guncost = 4
		} else if wlvl == 2 {
			guncost = 2
		}
		if ch.Gold < guncost {
			send_to_char(ch, libc.CString("You do not have enough zenni. You need %d zenni per shot for that level of gun.\r\n"), guncost)
			return
		} else {
			ch.Gold -= guncost
		}
	} else if wielded == 2 && gun == TRUE {
		if wlvl == 5 {
			guncost = 12
		} else if wlvl == 4 {
			guncost = 6
		} else if wlvl == 3 {
			guncost = 4
		} else if wlvl == 2 {
			guncost = 2
		}
		if ch.Gold < guncost {
			send_to_char(ch, libc.CString("You do not have enough zenni. You need %d zenni per shot for that level of gun.\r\n"), guncost)
			return
		} else {
			ch.Gold -= guncost
		}
	}
	if vict != nil {
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		if gun == TRUE {
			if dualwield >= 2 {
				prob += int(float64(prob) * 0.1)
			}
		}
		prob -= avo
		if PLR_FLAGGED(ch, PLR_THANDW) {
			perc += 15
		}
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
					pcost(ch, 0, stcost/3)
				}
				pcost(vict, 0, vict.Max_hit/150)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W intercepts and parries your attack with $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou intercept and parry @C$n's@W attack with one of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W intercepts and parries @c$n's@W attack with one of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if wtype != 4 {
						handle_disarm(ch, vict)
					}
					improve_skill(vict, SKILL_PARRY, 0)
					if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
						pcost(ch, 0, stcost)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
						pcost(ch, 0, stcost)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, -1, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
						pcost(ch, 0, stcost)
					}
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your attack misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W moves to attack you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to attack @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
						pcost(ch, 0, stcost/3)
					}
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your attack misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves to attack you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to attack @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
					pcost(ch, 0, stcost/3)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, -1, skill, attperc)
			if OBJ_FLAGGED(weap, ITEM_WEAPLVL1) {
				dmg += int64(float64(dmg) * 0.05)
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL2) {
				dmg += int64(float64(dmg) * 0.1)
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL3) {
				dmg += int64(float64(dmg) * 0.2)
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL4) {
				dmg += int64(float64(dmg) * 0.3)
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL5) {
				dmg += int64(float64(dmg) * 0.5)
			}
			if wtype == 5 {
				if GET_SKILL(ch, SKILL_BRAWL) >= 100 {
					dmg += int64(float64(dmg) * 0.5)
					wlvl = 5
				} else if GET_SKILL(ch, SKILL_BRAWL) >= 50 {
					dmg += int64(float64(dmg) * 0.2)
					wlvl = 3
				}
			}
			if wtype == 0 && int(ch.Race) == RACE_KONATSU {
				dmg += int64(float64(dmg) * 0.25)
			}
			if PLR_FLAGGED(ch, PLR_THANDW) {
				dmg += int64(float64(dmg) * 1.2)
			}
			if !IS_NPC(ch) {
				if PLR_FLAGGED(ch, PLR_THANDW) && gun == FALSE && gun2 == FALSE {
					if int(ch.Skills[SKILL_TWOHAND]) >= 100 {
						dmg += int64(float64(dmg) * 0.5)
					} else if int(ch.Skills[SKILL_TWOHAND]) >= 75 {
						dmg += int64(float64(dmg) * 0.25)
					} else if int(ch.Skills[SKILL_TWOHAND]) >= 50 {
						dmg += int64(float64(dmg) * 0.1)
					}
					if wtype == 3 {
						switch wlvl {
						case 1:
							dmg += int64(float64(dmg) * 0.04)
						case 2:
							dmg += int64(float64(dmg) * 0.08)
						case 3:
							dmg += int64(float64(dmg) * 0.12)
						case 4:
							dmg += int64(float64(dmg) * 0.2)
						case 5:
							dmg += int64(float64(dmg) * 0.25)
						}
					}
				}
			}
			if wtype == 3 {
				if skill >= 100 {
					dmg += int64(float64(dmg) * 0.04)
				} else if skill >= 50 {
					dmg += int64(float64(dmg) * 0.1)
				}
			}
			var hitspot int = 1
			if gun == TRUE {
				dmg = gun_dam(ch, wlvl)
			}
			hitspot = roll_hitloc(ch, vict, skill)
			var beforepl int64 = vict.Hit
			switch hitspot {
			case 1:
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
			case 2:
				hitspot = 4
			case 3:
				hitspot = 5
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
			case 4:
				hitspot = 2
			case 5:
				hitspot = 3
			}
			if PLR_FLAGGED(ch, PLR_THANDW) && gun == TRUE {
				if hitspot != 4 && boom_headshot(ch) != 0 {
					hitspot = 4
					send_to_char(ch, libc.CString("@GBoom headshot!@n\r\n"))
				}
			}
			switch wtype {
			case 0:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou slash @C$N@W across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou slash @C$N@W across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou slash @C$N@W across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou slash @C$N@W across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if !IS_NPC(ch) {
						if PLR_FLAGGED(ch, PLR_THANDW) && gun == FALSE && gun2 == FALSE {
							if int(ch.Skills[SKILL_TWOHAND]) >= 100 {
								var mult float64 = calc_critical(ch, 0)
								mult += 1.0
								dmg *= int64(mult)
							} else {
								dmg *= int64(calc_critical(ch, 0))
							}
						} else {
							dmg *= int64(calc_critical(ch, 0))
						}
					} else {
						dmg *= int64(calc_critical(ch, 0))
					}
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou slash @C$N@W across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
				if float64(beforepl-vict.Hit) >= float64(gear_pl(vict))*0.025 {
					cut_limb(ch, vict, wlvl, hitspot)
				}
			case 1:
				if ch.Fighting == nil && backstab(ch, vict, wlvl, dmg) != 0 {
					if vict != nil && vict.Hit > 1 && axion_dice(0) < GET_SKILL(ch, SKILL_DUALWIELD) && (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
						do_attack2(ch, nil, 0, 0)
					}
					pcost(ch, 0, stcost)
					return
				}
				dmg += int64((float64(dmg) * 0.01) * (float64(ch.Aff_abils.Dex) * 0.5))
				switch hitspot {
				case 1:
					act(libc.CString("@WYou pierce @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou pierce @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou pierce @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou pierce @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou pierce @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 2:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou crush @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W chest@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou crush @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou crush @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou crush @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes you in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou crush @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W stomach@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
				club_stamina(ch, vict, wlvl, dmg)
			case 3:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou stab @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou stab @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou stab @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou stab @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs you in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou stab @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W stomach@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 4:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou blast @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou blast @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou blast @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou blast @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts you in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou blast @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W stomach@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 5:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou whack @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou whack @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou whack @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou whack @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks you in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N@W in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if GET_SKILL(ch, SKILL_BRAWL) >= 100 {
						var mult float64 = calc_critical(ch, 0)
						mult += 1.0
						dmg *= int64(mult)
					} else {
						dmg *= int64(calc_critical(ch, 0))
					}
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou whack @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W stomach@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			}
		}
		if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
			if (ch.Equipment[WEAR_WIELD1]) != nil {
				if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
					act(libc.CString("@c$N's@W fireshield burns your weapon!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n's@W weapon is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n's@W weapon is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var damdam int = GET_SKILL(vict, SKILL_FIRESHIELD) / 2
					hurt(0, 0, vict, nil, ch.Equipment[WEAR_WIELD1], int64(damdam), 0)
				} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
					send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
				}
				pcost(ch, 0, stcost)
			} else {
				if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
					act(libc.CString("@c$N's@W fireshield burns your weapon!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n's@W weapon is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n's@W weapon is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var damdam int = GET_SKILL(vict, SKILL_FIRESHIELD) / 2
					hurt(0, 0, vict, nil, ch.Equipment[WEAR_WIELD2], int64(damdam), 0)
				}
				pcost(ch, 0, stcost)
			}
		}
		if gun == FALSE && gun2 == FALSE {
			damage_weapon(ch, weap, vict)
		}
		if !IS_NPC(ch) {
			if PLR_FLAGGED(ch, PLR_THANDW) {
				if GET_SKILL(ch, SKILL_TWOHAND) == 0 && slot_count(ch)+1 <= ch.Skill_slots {
					var numb int = rand_number(10, 15)
					for {
						ch.Skills[SKILL_TWOHAND] = int8(numb)
						if true {
							break
						}
					}
					send_to_char(ch, libc.CString("@GYou learn the very basics of two-handing your weapon!@n\r\n"))
				} else {
					improve_skill(ch, SKILL_TWOHAND, 0)
				}
			}
		}
		if (ch.Equipment[WEAR_WIELD2]) != nil {
			if GET_SKILL(ch, SKILL_DUALWIELD) == 0 && slot_count(ch)+1 <= ch.Skill_slots && int((ch.Equipment[WEAR_WIELD2]).Type_flag) != ITEM_LIGHT {
				var numb int = rand_number(10, 15)
				for {
					ch.Skills[SKILL_DUALWIELD] = int8(numb)
					if true {
						break
					}
				}
				send_to_char(ch, libc.CString("@GYou learn the very basics of dual-wielding!@n\r\n"))
			} else {
				improve_skill(ch, SKILL_DUALWIELD, 0)
			}
			if vict != nil && vict.Hit > 1 && axion_dice(0) < GET_SKILL(ch, SKILL_DUALWIELD) && (ch.Equipment[WEAR_WIELD1]) != nil {
				do_attack2(ch, nil, 0, 0)
			}
		}
		handle_multihit(ch, vict)
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou attack $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W attacks $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		if wielded == 2 && gun == FALSE || gun2 == FALSE && gun == FALSE {
			pcost(ch, 0, stcost)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_shogekiha(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.125
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SHOGEKIHA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SHOGEKIHA)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Chclass) != CLASS_KABITO {
		handle_cooldown(ch, 5)
	} else {
		if GET_SKILL(ch, SKILL_SHOGEKIHA) < 100 {
			handle_cooldown(ch, 5)
		}
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SHOGEKIHA, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		var master_roll int = rand_number(1, 100)
		var master_chance int = 0
		var master_pass int = FALSE
		if skill >= 100 {
			master_chance = 20
		} else if skill >= 75 {
			master_chance = 10
		} else if skill >= 50 {
			master_chance = 5
		}
		if master_chance >= master_roll {
			master_pass = TRUE
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Shogekiha before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Shogekiha before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Shogekiha before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks shogekiha!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W shogekiha!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W shogekiha!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 10, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your shogekiha, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W shogekiha, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W shogekiha, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your shogekiha misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W aims $s hand and releases a shogekiha at you, but misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W aims $s hand and releases a shogekiha at @C$N@W, but somehow misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your shogekiha misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W aims $s hand and releases a shogekiha at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s hand and releases a shogekiha at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 10, skill, attperc)
			if int(ch.Chclass) == CLASS_KABITO {
				if GET_SKILL(ch, SKILL_SHOGEKIHA) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.15)
				} else if GET_SKILL(ch, SKILL_SHOGEKIHA) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_SHOGEKIHA) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou aim your hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly you unleash a large shogekiha that slams into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into $N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou aim your hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly you unleash a large shogekiha that slams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into $N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou aim your hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly you unleash a large shogekiha that slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into $N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou aim your hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly you unleash a large shogekiha that slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into $N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou aim your hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly you unleash a large shogekiha that slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and nearby loose objects begin to be pushed out by an invisible force. Suddenly $e unleashes a large shogekiha that slams into $N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if master_pass == TRUE {
				act(libc.CString("@CYour skillful shogekiha disipated some of @c$N's@C charged ki!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@C's skillful shogekiha disipated some of YOUR charged ki!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@C's skillful shogekiha disipated some of @c$N's@C charged ki!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Charge -= int64(float64(vict.Charge) * 0.25)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		if obj.Kicharge > 0 && ch.Charge > obj.Kicharge {
			act(libc.CString("@WYou leap at $p@W with your arms spread out to your sides. As you are about to make contact with $p@W you scream and shatter the attack with your ki!@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("@C$n@W leaps out at $p@W with $s arms spead out to $s sides. As $e is about to make contact with $p@W $e screams and shatters the attack with $s ki!@n"), TRUE, ch, obj, nil, TO_ROOM)
			obj.Kicharge -= ch.Charge
			extract_obj(obj)
		} else if obj.Kicharge > 0 && ch.Charge < obj.Kicharge {
			act(libc.CString("@WYou leap at $p@W with your arms spread out to your sides. As you are about to make contact with $p@W you scream and weaken the attack with your ki before taking the rest of the attack in the chest!@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("@C$n@W leaps out at $p@W with $s arms spead out to $s sides. As $e is about to make contact with $p@W $e screams and weakens the attack with $s ki before taking the rest of the attack in the chest!@n"), TRUE, ch, obj, nil, TO_ROOM)
			obj.Kicharge -= ch.Charge
			ch.Charge = 0
			dmg = obj.Kicharge
			hurt(0, 0, obj.User, ch, nil, dmg, 0)
			extract_obj(obj)
		} else {
			dmg = damtype(ch, 10, skill, attperc)
			dmg /= 10
			act(libc.CString("@WYou fire a shogekiha at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("@C$n@W fires a shogekiha at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
			hurt(0, 0, ch, nil, obj, dmg, 0)
			pcost(ch, attperc, 0)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_tsuihidan(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		dmg     int64
		attperc float64 = 0.1
		minimum float64 = 0.05
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_TSUIHIDAN) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_TSUIHIDAN)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_TSUIHIDAN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		var master_roll int = rand_number(1, 100)
		var master_chance int = 0
		var master_pass int = FALSE
		if skill >= 100 {
			master_chance = 20
		} else if skill >= 75 {
			master_chance = 10
		} else if skill >= 50 {
			master_chance = 5
		}
		if master_chance >= master_roll {
			master_pass = TRUE
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Tsuihidan before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Tsuihidan before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Tsuihidan before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(vict, 0, vict.Max_hit/200)
				dodge_ki(ch, vict, 1, 11, skill, SKILL_TSUIHIDAN)
				hurt(0, 0, ch, vict, nil, 0, 1)
				pcost(ch, attperc, 0)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects your tsuihidan, sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect @C$n's@W tsuihidan sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects @c$n's@W tsuihidan sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(vict, 0, vict.Max_hit/200)
					parry_ki(attperc, ch, vict, func() [1000]byte {
						var t [1000]byte
						copy(t[:], []byte("tsuihidan"))
						return t
					}(), prob, perc, skill, 11)
					pcost(ch, attperc, 0)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks tsuihidan!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W tsuihidan!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W tsuihidan!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 11, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your tsuihidan, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W tsuihidan, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W tsuihidan, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your tsuihidan misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a tsuihidan at you, but misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a tsuihidan at @C$N@W, but somehow misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 1, 11, skill, SKILL_TSUIHIDAN)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your tsuihidan misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a tsuihidan at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a tsuihidan at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dodge_ki(ch, vict, 1, 11, skill, SKILL_TSUIHIDAN)
				hurt(0, 0, ch, vict, nil, 0, 1)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 11, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large tsuihidan that slams into $s chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into $N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large tsuihidan that slams into $s face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into $N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large tsuihidan that slams into $s gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into $N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large tsuihidan that slams into $s arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into $N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large tsuihidan that slams into $s leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large tsuihidan that slams into $N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if master_pass == TRUE {
				vict.Move -= dmg
				if vict.Move < 0 {
					vict.Move = 0
				}
				act(libc.CString("@CYour tsuihidan hits a vital spot and seems to sap some of @c$N's@C stamina!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@C tsuihidan hits a vital spot and saps some of your stamina!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n's@C tsuihidan hits a vital spot and saps some of @c$N's@C stamina!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 11, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a tsuihidan at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a tsuihidan at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_attack2(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob      int
		perc      int
		avo       int
		index     int = 0
		pry       int = 0
		dge       int = 0
		blk       int = 0
		skill     int = 0
		wtype     int = 0
		gun2      int = FALSE
		dualwield int = 0
		dmg       int64
		vict      *char_data = nil
		obj       *obj_data  = nil
		arg       [2048]byte
		attperc   float64 = 0
	)
	one_argument(argument, &arg[0])
	if (ch.Equipment[WEAR_WIELD2]) == nil {
		return
	}
	var stcost int64 = ((ch.Max_hit / 150) + (ch.Equipment[WEAR_WIELD2]).Weight)
	var kicost int64 = ((ch.Max_hit / 150) + (ch.Equipment[WEAR_WIELD2]).Weight)
	_ = kicost
	if int(ch.Race) == RACE_ANDROID {
		stcost *= int64(0.25)
	}
	if int(ch.Race) == RACE_ANDROID && gun2 == TRUE {
		kicost *= int64(0.25)
	}
	if can_grav(ch) == 0 {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("With what arms!?\r\n"))
		return
	} else if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 && (ch.Limb_condition[1]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right arm has damaged it more!@n\r\n"))
		ch.Limb_condition[0] -= rand_number(3, 5)
		if (ch.Limb_condition[0]) < 0 {
			act(libc.CString("@RYour right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 && (ch.Limb_condition[0]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left arm has damaged it more!@n\r\n"))
		ch.Limb_condition[1] -= rand_number(3, 5)
		if (ch.Limb_condition[1]) < 0 {
			act(libc.CString("@RYour left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if ch.Fighting == nil {
		return
	}
	if check_points(ch, 0, stcost) == 0 {
		return
	}
	if !IS_NPC(ch) || IS_NPC(ch) {
		if vict != nil {
			if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
				if can_kill(ch, vict, nil, 1) == 0 {
					return
				}
			} else {
				if can_kill(ch, vict, nil, 0) == 0 {
					return
				}
			}
		}
		if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
			skill = init_skill(ch, SKILL_DAGGER)
			improve_skill(ch, SKILL_DAGGER, 1)
			wtype = 1
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SWORD)
			improve_skill(ch, SKILL_SWORD, 1)
			wtype = 0
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_CRUSH-TYPE_HIT) {
			skill = init_skill(ch, SKILL_CLUB)
			improve_skill(ch, SKILL_CLUB, 1)
			wtype = 2
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
			skill = init_skill(ch, SKILL_SPEAR)
			improve_skill(ch, SKILL_SPEAR, 1)
			wtype = 3
		} else if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
			gun2 = TRUE
			skill = init_skill(ch, SKILL_GUN)
			improve_skill(ch, SKILL_GUN, 1)
			wtype = 4
		} else {
			skill = init_skill(ch, SKILL_BRAWL)
			improve_skill(ch, SKILL_BRAWL, 1)
			wtype = 5
		}
	}
	if gun2 == FALSE {
		if int(ch.Skills[SKILL_DUALWIELD]) >= 100 {
			dualwield = 3
			stcost -= int64(float64(stcost) * 0.3)
		} else if int(ch.Skills[SKILL_DUALWIELD]) >= 75 {
			dualwield = 2
			stcost -= int64(float64(stcost) * 0.25)
		} else if int(ch.Skills[SKILL_DUALWIELD]) >= 50 {
			dualwield = 1
			stcost -= int64(float64(stcost) * 0.25)
		}
	}
	if IS_NPC(ch) && GET_LEVEL(ch) <= 10 {
		skill = rand_number(30, 50)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 20 {
		skill = rand_number(30, 60)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 30 {
		skill = rand_number(30, 70)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 50 {
		skill = rand_number(40, 80)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 70 {
		skill = rand_number(50, 90)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 80 {
		skill = rand_number(60, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 90 {
		skill = rand_number(70, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 100 {
		skill = rand_number(80, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) > 100 {
		skill = rand_number(95, 100)
	}
	if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
		vict = ch.Fighting
	}
	if gun2 == TRUE {
		if ch.Gold < 1 {
			send_to_char(ch, libc.CString("You do not have enough zenni. You need 1 zenni per shot.\r\n"))
			return
		} else {
			ch.Gold -= 1
		}
	}
	if ch.Preference != PREFERENCE_H2H {
		handle_cooldown(ch, 4)
	} else {
		handle_cooldown(ch, 8)
	}
	if vict != nil {
		if ((ch.Equipment[WEAR_WIELD2]).Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
			if can_kill(ch, vict, nil, 1) == 0 {
				return
			}
		} else {
			if can_kill(ch, vict, nil, 0) == 0 {
				return
			}
		}
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		if dualwield == 3 {
			pry -= int(float64(pry) * 0.1)
			blk -= int(float64(pry) * 0.1)
			dge -= int(float64(pry) * 0.1)
		}
		if gun2 == TRUE {
			if dualwield >= 1 {
				prob += int(float64(prob) * 0.1)
			}
		}
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your attack before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c attack before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c attack before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if gun2 == FALSE {
					pcost(ch, 0, stcost/3)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W intercepts and parries your attack with $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou intercept and parry @C$n's@W attack with one of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W intercepts and parries @c$n's@W attack with one of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if wtype != 4 {
						handle_disarm(ch, vict)
					}
					improve_skill(vict, SKILL_PARRY, 0)
					if gun2 == FALSE {
						pcost(ch, 0, stcost)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					if gun2 == FALSE {
						pcost(ch, 0, stcost)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, -1, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					if gun2 == FALSE {
						pcost(ch, 0, stcost)
					}
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your attack misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W moves to attack you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to attack @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if gun2 == FALSE {
						pcost(ch, 0, stcost/3)
					}
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your attack misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves to attack you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to attack @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if gun2 == FALSE {
					pcost(ch, 0, stcost/3)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, -1, skill, attperc)
			var wlvl int = 0
			var weap *obj_data = (ch.Equipment[WEAR_WIELD2])
			if OBJ_FLAGGED(weap, ITEM_WEAPLVL1) {
				dmg += int64(float64(dmg) * 0.05)
				wlvl = 1
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL2) {
				dmg += int64(float64(dmg) * 0.1)
				wlvl = 2
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL3) {
				dmg += int64(float64(dmg) * 0.2)
				wlvl = 3
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL4) {
				dmg += int64(float64(dmg) * 0.3)
				wlvl = 4
			} else if OBJ_FLAGGED(weap, ITEM_WEAPLVL5) {
				dmg += int64(float64(dmg) * 0.5)
				wlvl = 5
			}
			if wtype == 5 {
				if GET_SKILL(ch, SKILL_BRAWL) >= 100 {
					dmg += int64(float64(dmg) * 0.5)
					wlvl = 5
				} else if GET_SKILL(ch, SKILL_BRAWL) >= 50 {
					dmg += int64(float64(dmg) * 0.2)
					wlvl = 3
				}
			}
			if wtype == 0 && int(ch.Race) == RACE_KONATSU {
				dmg += int64(float64(dmg) * 0.25)
			}
			var hitspot int = 1
			if gun2 == TRUE {
				dmg = gun_dam(ch, wlvl)
			}
			hitspot = roll_hitloc(ch, vict, skill)
			var beforepl int64 = vict.Hit
			if wtype == 3 {
				if skill >= 100 {
					dmg += int64(float64(dmg) * 0.04)
				} else if skill >= 50 {
					dmg += int64(float64(dmg) * 0.1)
				}
			}
			if gun2 == TRUE {
				if dualwield == 3 && rand_number(1, 3) == 3 {
					send_to_char(ch, libc.CString("@GYour masterful aim scores a critical!@n\r\n"))
					hitspot = 2
				}
			}
			switch hitspot {
			case 1:
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
			case 2:
				hitspot = 4
			case 3:
				hitspot = 5
			case 4:
				hitspot = 5
			case 5:
				hitspot = 5
			}
			switch wtype {
			case 0:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou slash @C$N@W across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou slash @C$N@W across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou slash @C$N@W across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou slash @C$N@W across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou slash @C$N@W across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W slashes you across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W slashes @C$N@W across the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
				if float64(beforepl-vict.Hit) >= float64(gear_pl(vict))*0.025 {
					cut_limb(ch, vict, wlvl, hitspot)
				}
			case 1:
				dmg += int64((float64(dmg) * 0.01) * (float64(ch.Aff_abils.Dex) * 0.5))
				switch hitspot {
				case 1:
					act(libc.CString("@WYou pierce @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou pierce @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N'@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou pierce @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou pierce @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou pierce @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W pierces your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W pierces @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 2:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou crush @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou crush @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou crush @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou crush @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N'@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou crush @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W crushes your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W crushes @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
				club_stamina(ch, vict, wlvl, dmg)
			case 3:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou stab @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou stab @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou stab @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou stab @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N'@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou stab @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W stabs your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W stabs @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 4:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou blast @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou blast @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou blast @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou blast @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N'@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou blast @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W blasts your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W blasts @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			case 5:
				switch hitspot {
				case 1:
					act(libc.CString("@WYou whack @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W stomach!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou whack @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 1)
				case 3:
					act(libc.CString("@WYou whack @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 2)
				case 4:
					act(libc.CString("@WYou whack @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N'@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if GET_SKILL(ch, SKILL_BRAWL) >= 100 {
						var mult float64 = calc_critical(ch, 0)
						mult += 1.0
						dmg *= int64(mult)
					} else {
						dmg *= int64(calc_critical(ch, 0))
					}
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 3)
				case 5:
					act(libc.CString("@WYou whack @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W whacks your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W whacks @C$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					dam_eq_loc(vict, 4)
				}
			}
		}
		if gun2 == FALSE {
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your weapon!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W weapon is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W weapon is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				var damdam int = GET_SKILL(vict, SKILL_FIRESHIELD) / 2
				hurt(0, 0, vict, nil, ch.Equipment[WEAR_WIELD2], int64(damdam), 0)
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou attack $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W attacks $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		if gun2 == FALSE {
			pcost(ch, 0, stcost)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_bite(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = (ch.Max_hit / 500) + 20
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) && (int(ch.Race) != RACE_MUTANT || (ch.Genome[0]) != 7 && (ch.Genome[1]) != 7) {
		send_to_char(ch, libc.CString("You don't want to put that in your mouth, you don't know where it has been!\r\n"))
		return
	}
	if can_grav(ch) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		return
	}
	if check_points(ch, 0, ch.Max_hit/200) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_PUNCH)
	if skill <= 0 {
		skill = 60
	}
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			return
		}
	}
	handle_cooldown(ch, 4)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your bite before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c bite before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c bite before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your bite, but your zanzoken is faster!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your bite with a punch of their own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W bite with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W bite with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -1, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 6, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou move to bite @C$N@W, but miss!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W moves to bite you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to bite @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@C$n@W moves to bite you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to bite @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 8, skill, attperc)
			if !IS_NPC(ch) {
				dmg = damtype(ch, 0, skill, attperc)
			}
			dmg += int64(float64(dmg) * 0.25)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou bite @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bites your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W bites into $N's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 2:
				act(libc.CString("@WYou bite @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bites you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bites @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 3:
				act(libc.CString("@WYou bite @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bites you on the body, sending blood flying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bites @c$N@W on the body, sending blood flying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou bite @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bites you on the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bites @c$N@W on the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou bite @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bites into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W bites into @c$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if !IS_NPC(ch) {
				if axion_dice(0) > int(vict.Aff_abils.Con) && rand_number(1, 5) == 5 {
					act(libc.CString("@R$N@r was poisoned by your bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@rYou were poisoned by the bite!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					vict.Poisonby = ch
					var duration int = (int(ch.Aff_abils.Intel) / 50) + 1
					assign_affect(vict, AFF_POISON, SKILL_POISON, duration, 0, 0, 0, 0, 0, 0)
				}
			}
			pcost(ch, 0, stcost)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@C$n@W bites $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kiball(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob     int
		perc     int
		avo      int
		index    int = 0
		pry      int = 2
		dge      int = 2
		blk      int = 2
		skill    int = 0
		frompool int = FALSE
		dmg      int64
		attperc  float64 = 0.05
		minimum  float64 = 0.01
		vict     *char_data
		obj      *obj_data
		arg      [2048]byte
		arg2     [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KIBALL) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if frompool == 0 && float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KIBALL)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KIBALL, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		var mult_roll int = rand_number(1, 100)
		var mult_count int = 1
		var mult_chance int = 0
		if skill >= 100 {
			mult_chance = 30
		} else if skill >= 75 {
			mult_chance = 15
		} else if skill >= 50 {
			mult_chance = 10
		}
		if mult_roll <= mult_chance {
			mult_count = rand_number(2, 3)
		}
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your kiball before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c kiball before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c kiball before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if frompool == TRUE {
					ch.Mana -= int64(float64(ch.Max_mana) * attperc)
				} else {
					pcost(ch, attperc, 0)
				}
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects your kiball, sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect @C$n's@W kiball sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects @c$n's@W kiball sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(vict, 0, vict.Max_hit/200)
					parry_ki(attperc, ch, vict, func() [1000]byte {
						var t [1000]byte
						copy(t[:], []byte("kiball"))
						return t
					}(), prob, perc, skill, 7)
					pcost(ch, attperc, 0)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 98 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 2
					}
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks kiball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W kiball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W kiball!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 7, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your kiball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kiball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W kiball, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 0, 7, skill, SKILL_KIBALL)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 98 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 2
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your kiball misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a bright yellow kiball at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a bright yellow kiball at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your kiball misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a bright yellow kiball at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a bright yellow kiball at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			if mult_count > 1 {
				act(libc.CString("@CYour expertise has allowed you to fire multiple shots in a row!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@C expertise has allowed $m to fire multiple shots in a row!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			}
			for mult_count > 0 {
				mult_count -= 1
				dmg = damtype(ch, 7, skill, attperc)
				var hitspot int = 1
				hitspot = roll_hitloc(ch, vict, skill)
				switch hitspot {
				case 1:
					act(libc.CString("@WYou hold out your hand towards @C$N@W, and fire a bright yellow kiball! The kiball slams into $M quickly and explodes with roaring light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W holds out $s hand towards you, and fires a bright yellow kiball! The kiball slams into you quickly and explodes with roaring light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W holds out $s hand towards @C$N@W, and fires a bright yellow kiball! The kiball slams into $M quickly and explodes with roaring light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (ch.Bonuses[BONUS_SOFT]) != 0 {
						dmg *= int64(calc_critical(ch, 2))
					}
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 4)
				case 2:
					act(libc.CString("@WYou hold out your hand towards @C$N@W, and fire a bright yellow kiball! The kiball slams into $S face and explodes, shrouding $S head with smoke!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W holds out $s hand towards you, and fires a bright yellow kiball! The kiball slams into your face and explodes, leaving you choking on smoke!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W holds out $s hand towards @C$N@W, and fires a bright yellow kiball! The kiball slams into $S face and explodes, shrouding $S head with smoke!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 0))
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 3)
				case 3:
					act(libc.CString("@WYou hold out your hand towards @C$N@W, and fire a bright yellow kiball! The kiball slams into $S body and explodes with a loud roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W holds out $s hand towards you, and fires a bright yellow kiball! The kiball slams into your body and explodes with a loud roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W holds out $s hand towards @C$N@W, and fires a bright yellow kiball! The kiball slams into $S body and explodes with a loud roar!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (ch.Bonuses[BONUS_SOFT]) != 0 {
						dmg *= int64(calc_critical(ch, 2))
					}
					hurt(0, 0, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 4)
				case 4:
					act(libc.CString("@WYou hold out your hand towards @C$N@W, and fire a bright yellow kiball! The kiball grazes $S arm and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W holds out $s hand towards you, and fires a bright yellow kiball! The kiball grazes your arm and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W holds out $s hand towards @C$N@W, and fires a bright yellow kiball! The kiball grazes $S arm and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 195, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 1)
				case 5:
					act(libc.CString("@WYou hold out your hand towards @C$N@W, and fire a bright yellow kiball! The kiball grazes $S leg and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@c$n@W holds out $s hand towards you, and fires a bright yellow kiball! The kiball grazes your leg and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W holds out $s hand towards @C$N@W, and fires a bright yellow kiball! The kiball grazes $S leg and explodes shortly after!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= int64(calc_critical(ch, 1))
					hurt(1, 195, ch, vict, nil, dmg, 1)
					dam_eq_loc(vict, 2)
				}
				if vict.Hit <= 0 {
					mult_count = 0
				}
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 7, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a kiball at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a kiball at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_beam(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		attperc float64 = 0.1
		minimum float64 = 0.01
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BEAM) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_BEAM)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_BEAM, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your beam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c beam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c beam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects your beam, sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect @C$n's@W beam sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects @c$n's@W beam sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(vict, 0, vict.Max_hit/200)
					parry_ki(attperc, ch, vict, func() [1000]byte {
						var t [1000]byte
						copy(t[:], []byte("beam"))
						return t
					}(), prob, perc, skill, 10)
					pcost(ch, attperc, 0)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks beam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W beam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W beam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 10, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your beam, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W beam, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W beam, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 0, 10, skill, SKILL_BEAM)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your beam misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a bright yellow beam at you, but misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a bright yellow beam at @C$N@W, but somehow misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your beam misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a bright yellow beam at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a bright yellow beam at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 10, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large beam that slams into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into $N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large beam that slams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into $N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large beam that slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into $N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large beam that slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into $N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large beam that slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large beam that slams into $N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, attperc, 0)
			var master_roll int = rand_number(1, 100)
			var master_chance int = 0
			var master_pass int = FALSE
			if skill >= 100 {
				master_chance = 20
			} else if skill >= 75 {
				master_chance = 10
			} else if skill >= 50 {
				master_chance = 5
			}
			if master_chance >= master_roll {
				master_pass = TRUE
			}
			if vict.Hit > 0 && dmg > vict.Max_hit/4 && master_pass == TRUE {
				var (
					attempt int = rand_number(0, NUM_OF_DIRS)
					count   int = 0
				)
				for count < 12 {
					attempt = count
					if CAN_GO(vict, attempt) {
						count = 12
					} else {
						count++
					}
				}
				if CAN_GO(vict, attempt) {
					act(libc.CString("$N@W is pushed away by the blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou are pushed away by the blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$N@W is pushed away by the blast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					do_simple_move(vict, attempt, TRUE)
				} else {
					act(libc.CString("$N@W is pushed away by the blast, but is slammed into an obstruction!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou are pushed away by the blast, but are slammed into an obstruction!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("$N@W is pushed away by the blast, but is slammed into an obstruction!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dmg *= 2
					hurt(1, 195, ch, vict, nil, dmg, 1)
				}
			}
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 10, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a beam at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a beam at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kiblast(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		attperc float64 = 0.075
		minimum float64 = 0.01
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KIBLAST) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0]))) * 0.01
		if adjust < 0.01 || adjust > 1.0 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust < attperc && adjust >= minimum {
			attperc = adjust
		} else if adjust < minimum {
			attperc = minimum
		}
	}
	if float64(ch.Max_mana)*attperc > float64(ch.Charge) {
		attperc = float64(ch.Charge) / float64(ch.Max_mana)
	}
	if check_points(ch, int64(float64(ch.Max_mana)*minimum), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KIBLAST)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Race) != RACE_ANDROID || GET_SKILL(ch, SKILL_KIBLAST) < 100 {
		handle_cooldown(ch, 5)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KIBLAST, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = chance_to_hit(ch)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		var mastery int = rand_number(1, 100)
		var master_pass int = FALSE
		var chance int = 0
		if skill >= 100 {
			chance = 30
		} else if skill >= 75 {
			chance = 20
		} else if skill >= 50 {
			chance = 15
		}
		if mastery <= chance {
			master_pass = TRUE
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your kiblast before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c kiblast before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c kiblast before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if int(vict.Race) == RACE_ANDROID && HAS_ARMS(vict) && GET_SKILL(vict, SKILL_ABSORB) > rand_number(1, 140) {
					act(libc.CString("@C$N@W absorbs your ki attack and all your charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou absorb @C$n's@W ki attack and all $s charged ki with your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W absorbs @c$n's@W ki attack and all $s charged ki with $S hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var amot int = int(ch.Charge)
					if IS_NPC(ch) {
						amot = int(ch.Max_mana / 20)
					}
					if vict.Charge+int64(amot) > vict.Max_mana {
						vict.Mana += vict.Max_mana - vict.Charge
						vict.Charge = vict.Max_mana
					} else {
						vict.Charge += int64(amot)
					}
					pcost(ch, 1, 0)
					return
				} else if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W deflects your kiblast, sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou deflect @C$n's@W kiblast sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W deflects @c$n's@W kiblast sending it flying away!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(vict, 0, vict.Max_hit/200)
					parry_ki(attperc, ch, vict, func() [1000]byte {
						var t [1000]byte
						copy(t[:], []byte("kiblast"))
						return t
					}(), prob, perc, skill, 9)
					pcost(ch, attperc, 0)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks kiblast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W kiblast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W kiblast!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 9, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your kiblast, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kiblast, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W kiblast, letting it slam into the surroundings!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					dodge_ki(ch, vict, 0, 9, skill, SKILL_KIBLAST)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your kiblast misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a bright yellow kiblast at you, but misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a bright yellow kiblast at @C$N@W, but somehow misses!@n "), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your kiblast misses, flying through the air harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a bright yellow kiblast at you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a bright yellow kiblast at @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 9, skill, attperc)
			if int(ch.Race) == RACE_ANDROID {
				if GET_SKILL(ch, SKILL_KIBLAST) >= 100 {
					dmg += int64(float64(dmg) * 0.15)
				} else if GET_SKILL(ch, SKILL_KIBLAST) >= 60 {
					dmg += int64(float64(dmg) * 0.1)
				} else if GET_SKILL(ch, SKILL_KIBLAST) >= 40 {
					dmg += int64(float64(dmg) * 0.05)
				}
			}
			var record int64 = vict.Hit
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large kiblast that slams into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into $N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large kiblast that slams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into $N@W's face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large kiblast that slams into $S gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into $N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large kiblast that slams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into $N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou aim your hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly you unleash a large kiblast that slams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W$n@W aims $s hand at you, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@W$n@W aims $s hand at $N@W, and bright @Yyellow@W energy begins to pool there. Suddenly $e unleashes a large kiblast that slams into $N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if master_pass == TRUE && record > vict.Hit && float64(record-vict.Hit) > float64(gear_pl(vict))*0.025 {
				if !AFF_FLAGGED(vict, AFF_KNOCKED) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou are knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_KNOCKED/32)] |= 1 << (int(AFF_KNOCKED % 32))
					if AFF_FLAGGED(vict, AFF_FLYING) {
						vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
						vict.Altitude = 0
					}
					vict.Position = POS_SLEEPING
				}
			}
			pcost(ch, attperc, 0)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 1) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = damtype(ch, 10, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a kiblast at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a kiblast at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_slam(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_SLAM)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SLAM) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/100) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SLAM)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Chclass) == CLASS_BARDOCK {
		if int(ch.Skills[SKILL_STYLE]) >= 75 {
			handle_cooldown(ch, 7)
		} else {
			handle_cooldown(ch, 9)
		}
	} else {
		handle_cooldown(ch, 9)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SLAM, 1)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your slam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c slam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c slam before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your slam with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W slam with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W slam with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 6, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W slam!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your slam misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W moves to slam you with both $s fists, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to slam @C$N@W with both $s fists, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your slam misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W moves to slam you with both $s fists, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to slam @C$N@W with both $s fists, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 6, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou disappear, appearing above @C$N@W and slam a double fisted blow into $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W disappears, only to appear above you, slamming a double fisted blow into you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W disappears, only to appear above @C$N@W, slamming a double fisted blow into $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou disappear, reappearing in front of @C$N@W, you grab $M! Spinning you send $M flying into the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W disappears, reappearing in front of you, and $e grabs you! Spinning quickly $e sends you flying into the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W disappears, reappearing in front of @C$N@W, and grabs $M! Spinning quickly $e sends $M flying into the ground!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				if !AFF_FLAGGED(vict, AFF_KNOCKED) && (rand_number(1, 4) >= 3 && vict.Hit > ch.Hit/5 && !AFF_FLAGGED(vict, AFF_SANCTUARY)) {
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou are knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_KNOCKED/32)] |= 1 << (int(AFF_KNOCKED % 32))
					if AFF_FLAGGED(vict, AFF_FLYING) {
						vict.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
						vict.Altitude = 0
					}
					vict.Position = POS_SLEEPING
				} else if (int(vict.Position) == POS_STANDING || int(vict.Position) == POS_FIGHTING) && !AFF_FLAGGED(vict, AFF_KNOCKED) {
					vict.Position = POS_SITTING
				}
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg <= 95 && !ROOM_FLAGGED(vict.In_room, ROOM_SPACE) {
					act(libc.CString("@W$N@W slams into the ground forming a large crater with $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou slam into the ground forming a large crater with your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@W$N@W slams into the ground forming a large crater with $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) != SECT_INSIDE && (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) != SECT_UNDERWATER && (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) != SECT_WATER_SWIM && (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) != SECT_WATER_NOSWIM {
						impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
						switch rand_number(1, 8) {
						case 1:
							act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 2:
							if rand_number(1, 4) == 4 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect == 0 {
								(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect = 1
								act(libc.CString("Lava leaks up through cracks in the crater!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("Lava leaks up through cracks in the crater!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
							}
						case 3:
							act(libc.CString("A cloud of dust envelopes the entire area!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("A cloud of dust envelopes the entire area!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 4:
							act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 5:
							act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 6:
						default:
						}
					}
					if (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) == SECT_UNDERWATER {
						switch rand_number(1, 3) {
						case 1:
							act(libc.CString("The water churns violently!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The water churns violently!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 2:
							act(libc.CString("Large bubbles rise from the movement!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("Large bubbles rise from the movement!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 3:
							act(libc.CString("The water collapses in on the hole created!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The water collapses in on the hole create!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						}
					}
					if (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) == SECT_WATER_SWIM || (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) == SECT_WATER_NOSWIM {
						switch rand_number(1, 3) {
						case 1:
							act(libc.CString("A huge column of water erupts from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("A huge column of water erupts from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 2:
							act(libc.CString("The impact briefly causes a swirling vortex of water!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The impact briefly causes a swirling vortex of water!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 3:
							act(libc.CString("A huge depression forms in the water and erupts into a wave from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("A huge depression forms in the water and erupts into a wave from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						}
					}
					if (func() int {
						if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
						}
						return SECT_INSIDE
					}()) == SECT_INSIDE {
						impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
						switch rand_number(1, 8) {
						case 1:
							act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 2:
							act(libc.CString("The structure of the surrounding room cracks and quakes from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The structure of the surrounding room cracks and quakes from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 3:
							act(libc.CString("Parts of the ceiling collapse, crushing into the floor!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("Parts of the ceiling collapse, crushing into the floor!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 4:
							act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The surrounding area roars and shudders from the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 5:
							act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The ground shatters apart from the stress of the impact!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						case 6:
							act(libc.CString("The walls of the surrounding room crack in the same instant!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("The walls of the surrounding room crack in the same instant!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						default:
						}
					}
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg += 5
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou fly at @C$N@W, slamming both your fists into $S gut as you fly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W flies at you, slamming both $s fists into your gut as $e flies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W flies at @C$N@W, slamming both $s fists into $S gut as $e flies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou slam both your fists into @C$N@W, hitting $M in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams both $s fists into you, hitting you in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams both $s fists into @C$N@W, hitting $M in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou slam both your fists into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams both $s fists into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams both $s fists into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your hands!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W hands are burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W hands are burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou slam $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W slams $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_uppercut(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_UPPERCUT)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_UPPERCUT) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/200) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_UPPERCUT)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Chclass) == CLASS_FRIEZA {
		if int(ch.Skills[SKILL_STYLE]) >= 75 {
			handle_cooldown(ch, 5)
		} else {
			handle_cooldown(ch, 7)
		}
	} else {
		handle_cooldown(ch, 7)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_UPPERCUT, 1)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your uppercut before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c uppercut before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c uppercut before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your uppercut with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W uppercut with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W uppercut with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 5, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W uppercut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your uppercut misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws an uppercut at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws an uppercut at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your uppercut misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws an uppercut at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws an uppercut at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 5, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou leap up and launch an uppercut into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W leaps up and launches an uppercut into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W leaps up and launches an uppercut into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou smash an uppercut into @C$N's@W chin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W smashes an uppercut into your chin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W smashes an uppercut into @C$N's@W chin!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_KNOCKED) && (rand_number(1, 8) >= 7 && vict.Hit > ch.Hit/5 && !AFF_FLAGGED(vict, AFF_SANCTUARY)) {
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou are knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W is knocked out!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_KNOCKED/32)] |= 1 << (int(AFF_KNOCKED % 32))
					vict.Position = POS_SLEEPING
				}
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou uppercut @C$N@W, hitting $M directly in chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W uppercuts you, hitting you directly in the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W uppercuts @C$N@W, hitting $M directly in the chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYour poorly aimed uppercut hits @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W poorly aims an uppercut and hits you in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W poorly aims an uppercut and hits @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou slam an uppercut into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams an uppercut into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams an uppercut into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W hand is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W hand is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou uppercut $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W uppercuts $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_tailwhip(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_TAILWHIP)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_TAILWHIP) == 0 {
		return
	}
	if !PLR_FLAGGED(ch, PLR_TAIL) && !IS_NPC(ch) {
		send_to_char(ch, libc.CString("You have no tail!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/120) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_TAILWHIP)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 7)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_TAILWHIP, 1)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				ch.Combo -= 1
				ch.Combhits = 0
				act(libc.CString("@C$N@c disappears, avoiding your tailwhip before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c tailwhip before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c tailwhip before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your tailwhip with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W tailwhip with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W tailwhip with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 3, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W tailwhip!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your tailwhip misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a tailwhip at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a tailwhip at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your tailwhip misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a tailwhip at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a tailwhip at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 56, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou spin to swing your tail and slam it into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W spins to swing $s tail and slams it into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W spins to swing $s tail and slams it into @c$N@W's body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 8) >= 7 {
					handle_knockdown(vict)
				}
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou flip forward and slam your tail into the top of @c$N@W's head brutally!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W flips forward and slams $s tail into the top of YOUR head brutally!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W flips forward and slams $s tail into the top of @c$N@W's head brutally!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 8) >= 6 {
					handle_knockdown(vict)
				}
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou swing your tail and manage to slam it into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into YOUR gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 8) >= 7 {
					handle_knockdown(vict)
				}
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou swing your tail and manage to slam it into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou swing your tail and manage to slam it into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W swings $s tail and manages to slam it into @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou tailwhip $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W tailwhips $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_roundhouse(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_ROUNDHOUSE)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_ROUNDHOUSE) == 0 {
		return
	}
	if limb_ok(ch, 1) == 0 {
		return
	} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 && (ch.Limb_condition[3]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right leg has damaged it more!@n\r\n"))
		ch.Limb_condition[2] -= rand_number(3, 5)
		if (ch.Limb_condition[2]) < 0 {
			act(libc.CString("@RYour right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 && (ch.Limb_condition[2]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left leg has damaged it more!@n\r\n"))
		ch.Limb_condition[3] -= rand_number(3, 5)
		if (ch.Limb_condition[3]) < 0 {
			act(libc.CString("@RYour left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/180) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_ROUNDHOUSE)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	if int(ch.Chclass) == CLASS_NAIL {
		if int(ch.Skills[SKILL_STYLE]) >= 75 {
			handle_cooldown(ch, 5)
		} else {
			handle_cooldown(ch, 7)
		}
	} else {
		handle_cooldown(ch, 7)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_ROUNDHOUSE, 1)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your roundhouse before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c roundhouse before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c roundhouse before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your roundhouse with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W roundhouse with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W roundhouse with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 4, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W roundhouse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your roundhouse misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a roundhouse at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a roundhouse at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your roundhouse misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a roundhouse at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a roundhouse at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 4, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou spin in mid air and land a kick into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W spins in mid air and lands a kick into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W spins in mid air and lands a kick into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou spin a fierce roundhouse into @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W spins a fierce roundhouse into your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W spins a fierce roundhouse into @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 8) >= 7 {
					handle_knockdown(vict)
				}
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou throw a roundhouse at @C$N@W, hitting $M directly in the neck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a roundhouse at you, hitting YOU directly in the neck!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a roundhouse at @C$N@W, hitting $M directly in the face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYour poorly aimed roundhouse hits @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W poorly aims a roundhouse and hits you in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W poorly aims a roundhouse and hits @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou slam a roundhouse into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams a roundhouse into your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams a roundhouse into @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 195, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W leg is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W leg is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou roundhouse $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W roundhouses $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_elbow(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_ELBOW)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_ELBOW) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/300) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_ELBOW)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_ELBOW, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your elbow before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c elbow before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c elbow before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your elbow with an elbow of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W elbow with one of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W elbow with an elbow of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W blocks your elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou block @C$n's@W elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W blocks @c$n's@W elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 2, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W dodges your elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W dodges @c$n's@W elbow strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it, your elbow strike misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws an elbow strike at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws an elbow strike at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it, your elbow strike misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws an elbow stike at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws an elbow strike at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 0)
				pcost(ch, 0, 0)
			}
			return
		} else {
			dmg = damtype(ch, 2, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou slam your elbow into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s elbow into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s elbow into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou slam your elbow into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s elbow into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s elbow into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou land your elbow against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s elbow against your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s elbow against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou land your elbow against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s elbow against your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s elbow against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WYou land your elbow against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s elbow against your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s elbow against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your elbow!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W elbow is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W elbow is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou elbow stike $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W elbow strikes $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_kick(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_KICK)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KICK) == 0 {
		return
	}
	if limb_ok(ch, 1) == 0 {
		return
	} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 && (ch.Limb_condition[3]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right leg has damaged it more!@n\r\n"))
		ch.Limb_condition[2] -= rand_number(3, 5)
		if (ch.Limb_condition[2]) < 0 {
			act(libc.CString("@RYour right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 && (ch.Limb_condition[2]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left leg has damaged it more!@n\r\n"))
		ch.Limb_condition[3] -= rand_number(3, 5)
		if (ch.Limb_condition[3]) < 0 {
			act(libc.CString("@RYour left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/400) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KICK)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 5)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KICK, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your kick before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c kick before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c kick before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your kick with a kick of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W kick with one of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W kick with a kick of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W blocks your kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou block @C$n's@W kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W blocks @c$n's@W kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 1, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W dodges your kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W dodges @c$n's@W kick!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it, your kick misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a kick at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a kick at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it, your kick misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a kick at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a kick at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 0)
				pcost(ch, 0, stcost/2)
				return
			}
			return
		} else {
			dmg = damtype(ch, 1, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou slam your foot into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s foot into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s foot into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou slam your foot into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s foot into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s foot into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou land your foot against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s foot against your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s foot against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou land your foot against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s foot against your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s foot against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WYou land your foot against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s foot against your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s foot against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W leg is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W leg is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou kick $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W kicks $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_knee(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_KNEE)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KNEE) == 0 {
		return
	}
	if limb_ok(ch, 1) == 0 {
		return
	} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 && (ch.Limb_condition[3]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right leg has damaged it more!@n\r\n"))
		ch.Limb_condition[2] -= rand_number(3, 5)
		if (ch.Limb_condition[2]) < 0 {
			act(libc.CString("@RYour right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 && (ch.Limb_condition[2]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left leg has damaged it more!@n\r\n"))
		ch.Limb_condition[3] -= rand_number(3, 5)
		if (ch.Limb_condition[3]) < 0 {
			act(libc.CString("@RYour left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's left leg has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/250) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_KNEE)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 6)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_KNEE, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your knee strike before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c knee strike before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c knee strike before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your knee with a knee of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W knee with one of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W knee with a knee of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W blocks your knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou block @C$n's@W knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W blocks @c$n's@W knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 3, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W dodges your knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W dodges @c$n's@W knee strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it, your knee strike misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a knee strike at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a knee strike at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it, your  knee strike misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a knee stike at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a knee strike at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 0)
				pcost(ch, 0, 0)
			}
			return
		} else {
			dmg = damtype(ch, 3, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou slam your knee into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s knee into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s knee into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou slam your knee into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s knee into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s knee into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou land your knee against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s knee against your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s knee against @C$N's@W gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou land your knee against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s knee against your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s knee against @C$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WYou land your knee against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W lands $s knee against your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W lands $s knee against @C$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your knee!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W knee is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W knee is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou knee stike $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W knee strikes $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_punch(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int = 0
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int = 0
		dmg     int64
		stcost  int64 = physical_cost(ch, SKILL_PUNCH)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_PUNCH) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
		send_to_char(ch, libc.CString("Your hands are full!\r\n"))
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/500) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_PUNCH)
	vict = nil
	obj = nil
	if arg[0] == 0 || (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else if (func() *obj_data {
			obj = get_obj_in_list_vis(ch, &arg[0], nil, (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents)
			return obj
		}()) == nil {
			send_to_char(ch, libc.CString("Nothing around here by that name.\r\n"))
			return
		}
	}
	handle_cooldown(ch, 4)
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_PUNCH, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, FALSE != 0)
		perc = chance_to_hit(ch)
		if int(ch.Chclass) == CLASS_KABITO && !IS_NPC(ch) {
			if int(ch.Skills[SKILL_STYLE]) >= 75 {
				perc -= int(float64(perc) * 0.2)
			}
		}
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		prob -= avo
		if int(vict.Position) == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if int(vict.Position) == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if int(vict.Position) == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your punch before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c punch before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c punch before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				ch.Combo = -1
				ch.Combhits = 0
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@C$N@W parries your punch with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W punch with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W punch with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@C$N@W moves quickly and blocks your punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 0, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@C$N@W manages to dodge your punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W punch!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your punch misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a punch at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a punch at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your punch misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a punch at you but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a punch at @C$N@W but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 0, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou slam your fist into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s fist into your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s fist into @C$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou slam your fist into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W slams $s fist into your face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W slams $s fist into @C$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou punch @C$N@W directly in the gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W punches you directly in the gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W punches @C$N@W directly in the gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou punch @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W punches you in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W punches @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@WYou punch @C$N@W in the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W punches you in the leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W punches @C$N@W in the arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && (ch.Bonuses[BONUS_FIREPROOF]) == 0 && int(ch.Race) != RACE_DEMON {
				act(libc.CString("@c$N's@W fireshield burns your hand!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W hand is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W hand is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg = int64(float64(vict.Max_mana) * 0.02)
				vict.Lastattack += 1000
				hurt(0, 0, vict, ch, nil, dmg, 0)
				if (ch.Bonuses[BONUS_FIREPRONE]) != 0 {
					send_to_char(ch, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are easily burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				} else if int(ch.Aff_abils.Con) < axion_dice(0) {
					send_to_char(ch, libc.CString("@RYou are badly burned!@n\r\n"))
					send_to_char(vict, libc.CString("@RThey are burned!@n\r\n"))
					ch.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
				}
			} else if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) && ((ch.Bonuses[BONUS_FIREPROOF]) != 0 || int(ch.Race) == RACE_DEMON) {
				send_to_char(vict, libc.CString("@RThey appear to be fireproof!@n\r\n"))
			}
			pcost(ch, 0, stcost)
			handle_multihit(ch, vict)
			return
		}
	} else if obj != nil {
		if can_kill(ch, nil, obj, 0) == 0 {
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is broken already!\r\n"))
			return
		}
		dmg = (ch.Hit / 10000) + int64(ch.Aff_abils.Str)
		act(libc.CString("@WYou punch $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W punches $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_charge(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg [2048]byte
		amt int
	)
	one_argument(argument, &arg[0])
	if PLR_FLAGGED(ch, PLR_AURALIGHT) {
		send_to_char(ch, libc.CString("@WYou are concentrating too much on your aura to be able to charge."))
		return
	}
	if PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_MBREAK) {
		send_to_char(ch, libc.CString("Your mind is still strained from psychic attacks...\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_POISON) {
		send_to_char(ch, libc.CString("You feel too sick from the poison to concentrate.\r\n"))
		return
	}
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Charge, yes. How much percent though?\r\n"))
		send_to_char(ch, libc.CString("[ 1 - 100 | cancel | release]\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("release"), &arg[0]) == 0 && (PLR_FLAGGED(ch, PLR_CHARGE) && ch.Charge <= 0) {
		send_to_char(ch, libc.CString("Try cancel instead, you have nothing charged up yet!\r\n"))
		return
	} else if libc.StrCaseCmp(libc.CString("release"), &arg[0]) == 0 && (PLR_FLAGGED(ch, PLR_CHARGE) && ch.Charge > 0) {
		send_to_char(ch, libc.CString("You stop charging and release your pent up energy.\r\n"))
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("$n@w's aura fades.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("$n@w's aura flickers brightly before disappearing.@n"), TRUE, ch, nil, nil, TO_ROOM)
		default:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		ch.Mana += ch.Charge
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		ch.Charge = 0
		ch.Chargeto = 0
		ch.Act[int(PLR_CHARGE/32)] &= bitvector_t(int32(^(1 << (int(PLR_CHARGE % 32)))))
		return
	} else if libc.StrCaseCmp(libc.CString("release"), &arg[0]) == 0 && ch.Charge > 0 {
		send_to_char(ch, libc.CString("You release your pent up energy.\r\n"))
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("$n@w's aura fades.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("$n@w's aura flickers brightly before disappearing.@n"), TRUE, ch, nil, nil, TO_ROOM)
		default:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		ch.Mana += ch.Charge
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		ch.Charge = 0
		ch.Chargeto = 0
		return
	} else if libc.StrCaseCmp(libc.CString("cancel"), &arg[0]) == 0 && PLR_FLAGGED(ch, PLR_CHARGE) {
		send_to_char(ch, libc.CString("You stop charging.\r\n"))
		switch rand_number(1, 3) {
		case 1:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 2:
			act(libc.CString("$n@w's aura fades.@n"), TRUE, ch, nil, nil, TO_ROOM)
		case 3:
			act(libc.CString("$n@w's aura flickers brightly before disappearing.@n"), TRUE, ch, nil, nil, TO_ROOM)
		default:
			act(libc.CString("$n@w's aura disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		ch.Act[int(PLR_CHARGE/32)] &= bitvector_t(int32(^(1 << (int(PLR_CHARGE % 32)))))
		ch.Chargeto = 0
		return
	} else if libc.StrCaseCmp(libc.CString("cancel"), &arg[0]) == 0 && !PLR_FLAGGED(ch, PLR_CHARGE) {
		send_to_char(ch, libc.CString("You are not even charging!\r\n"))
		return
	} else if ch.Mana < ch.Max_mana/100 {
		send_to_char(ch, libc.CString("You don't even have enough ki!\r\n"))
		return
	} else if (func() int {
		amt = libc.Atoi(libc.GoString(&arg[0]))
		return amt
	}()) > 0 {
		if amt >= 101 {
			send_to_char(ch, libc.CString("You have set it too high!\r\n"))
			return
		} else if AFF_FLAGGED(ch, AFF_SPIRITCONTROL) {
			var diff int64 = 0
			if float64(ch.Mana) < ((float64(ch.Max_mana)*0.01)*float64(amt))+1 {
				diff = int64((((float64(ch.Max_mana) * 0.01) * float64(amt)) + 1) - float64(ch.Mana))
			}
			var chance int = 15
			chance -= GET_SKILL(ch, SKILL_SPIRITCONTROL) / 10
			if chance < 10 {
				chance = 10
			} else if chance > 15 {
				chance = 15
			}
			if chance > rand_number(1, 100) {
				send_to_char(ch, libc.CString("The rush of ki that you try to pool temporarily overwhelms you and you lose control!\r\n"))
				null_affect(ch, AFF_SPIRITCONTROL)
				return
			} else {
				var spiritcost int64 = int64(float64(ch.Max_mana) * 0.05)
				if GET_SKILL(ch, SKILL_SPIRITCONTROL) >= 100 {
					spiritcost = int64(float64(ch.Max_mana) * 0.01)
				}
				reveal_hiding(ch, 0)
				send_to_char(ch, libc.CString("Your %s colored aura flashes up around you as you instantly take control of the ki you needed!\r\n"), aura_types[ch.Aura])
				send_to_char(ch, libc.CString("@D[@RCost@D: @r%s@D]@n\r\n"), add_commas(spiritcost))
				var bloom [2048]byte
				stdio.Sprintf(&bloom[0], "@wA %s aura flashes up brightly around $n@w!@n", aura_types[ch.Aura])
				act(&bloom[0], TRUE, ch, nil, nil, TO_ROOM)
				ch.Charge = int64((((float64(ch.Max_mana) * 0.01) * float64(amt)) + 1) - float64(diff))
				ch.Mana -= int64((((float64(ch.Max_mana) * 0.01) * float64(amt)) + 1) - float64(diff))
				ch.Mana -= spiritcost
				if ch.Mana < 0 {
					ch.Mana = 0
				}
			}
		} else {
			reveal_hiding(ch, 0)
			send_to_char(ch, libc.CString("You begin to charge some energy, as a %s aura begins to burn around you!\r\n"), aura_types[ch.Aura])
			var bloom [2048]byte
			stdio.Sprintf(&bloom[0], "@wA %s aura flashes up brightly around $n@w!@n", aura_types[ch.Aura])
			act(&bloom[0], TRUE, ch, nil, nil, TO_ROOM)
			ch.Chargeto = int64(((float64(ch.Max_mana) * 0.01) * float64(amt)) + 1)
			ch.Charge += 1
			ch.Act[int(PLR_CHARGE/32)] |= bitvector_t(int32(1 << (int(PLR_CHARGE % 32))))
		}
	} else if amt < 1 && (func() room_vnum {
		if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
			return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
		}
		return -1
	}()) != 1562 {
		send_to_char(ch, libc.CString("You have set it too low!\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("That is an invalid argument.\r\n"))
	}
}
func do_powerup(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		ch.Act[int(MOB_POWERUP/32)] |= bitvector_t(int32(1 << (int(MOB_POWERUP % 32))))
		if ch.Max_hit < 50000 {
			act(libc.CString("@RYou begin to powerup, and air billows outward around you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and air billows outward around $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 500000 {
			act(libc.CString("@RYou begin to powerup, and loose objects are lifted into the air!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and loose objects are lifted into the air!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 5000000 {
			act(libc.CString("@RYou begin to powerup, and torrents of energy crackle around you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and torrents of energy crackle around $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 50000000 {
			act(libc.CString("@RYou begin to powerup, and the entire area begins to shudder!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and the entire area begins to shudder!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 100000000 {
			act(libc.CString("@RYou begin to powerup, and massive cracks begin to form beneath you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and massive cracks begin to form beneath $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 300000000 {
			act(libc.CString("@RYou begin to powerup, and everything around you shudders from the power!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and everything around $m shudders from the power!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			act(libc.CString("@RYou begin to powerup, and the very air around you begins to burn!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and the very air around $m begins to burn!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		return
	}
	if PLR_FLAGGED(ch, PLR_AURALIGHT) {
		send_to_char(ch, libc.CString("@WYou are concentrating too much on your aura to be able to power up."))
		return
	}
	if int(ch.Race) == RACE_ANDROID {
		send_to_char(ch, libc.CString("@WYou are an android, you do not powerup.@n"))
		return
	}
	if ch.Suppression > 0 {
		send_to_char(ch, libc.CString("@WYou currently have your powerlevel suppressed to %lld percent.@n"), ch.Suppression)
		return
	}
	if PLR_FLAGGED(ch, PLR_POWERUP) {
		send_to_char(ch, libc.CString("@WYou stop powering up.@n"))
		ch.Act[int(PLR_POWERUP/32)] &= bitvector_t(int32(^(1 << (int(PLR_POWERUP % 32)))))
		return
	}
	if ch.Hit >= ch.Max_hit {
		send_to_char(ch, libc.CString("@WYou are already at max!@n"))
		return
	}
	if ch.Mana < ch.Max_mana/20 {
		send_to_char(ch, libc.CString("@WYou do not have enough ki to powerup!@n"))
		return
	} else {
		reveal_hiding(ch, 0)
		if ch.Max_hit < 50000 {
			act(libc.CString("@RYou begin to powerup, and air billows outward around you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and air billows outward around $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 500000 {
			act(libc.CString("@RYou begin to powerup, and loose objects are lifted into the air!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and loose objects are lifted into the air!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 5000000 {
			act(libc.CString("@RYou begin to powerup, and torrents of energy crackle around you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and torrents of energy crackle around $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 50000000 {
			act(libc.CString("@RYou begin to powerup, and the entire area begins to shudder!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and the entire area begins to shudder!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 100000000 {
			act(libc.CString("@RYou begin to powerup, and massive cracks begin to form beneath you!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and massive cracks begin to form beneath $m!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else if ch.Max_hit < 300000000 {
			act(libc.CString("@RYou begin to powerup, and everything around you shudders from the power!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and everything around $m shudders from the power!@n"), TRUE, ch, nil, nil, TO_ROOM)
		} else {
			act(libc.CString("@RYou begin to powerup, and the very air around you begins to burn!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@R$n begins to powerup, and the very air around $m begins to burn!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
		ch.Act[int(PLR_POWERUP/32)] |= bitvector_t(int32(1 << (int(PLR_POWERUP % 32))))
		return
	}
}
func do_rescue(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [100]byte
		helpee   *char_data
		opponent *char_data
	)
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to rescue?\r\n"))
	} else if (func() *char_data {
		helpee = get_char_vis(ch, &arg[0], nil, 1<<0)
		return helpee
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if helpee == ch {
		send_to_char(ch, libc.CString("You can't help yourself any more than this!\r\n"))
	} else if helpee.Fighting == nil {
		send_to_char(ch, libc.CString("They are not fighting anyone!\r\n"))
	} else if ch.Fighting != nil && !IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are a little too busy fighting for yourself!\r\n"))
	} else {
		opponent = helpee.Fighting
		var mobbonus int = 0
		if IS_NPC(ch) {
			mobbonus = int(float64(GET_SPEEDI(ch)) * 0.2)
		}
		if GET_SPEEDI(ch)+mobbonus < GET_SPEEDI(opponent) && rand_number(1, 3) != 3 {
			act(libc.CString("@GYou leap towards @g$N@G and try to rescue $M but are too slow!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_CHAR)
			act(libc.CString("@g$n@G leaps towards you! $n is too slow and fails to rescue you!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_VICT)
			act(libc.CString("@g$n@G leaps towards @g$N@G and tries to rescue $M but is too slow!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_NOTVICT)
			return
		}
		act(libc.CString("@GYou leap in front of @g$N@G and rescue $M!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_CHAR)
		act(libc.CString("@g$n@G leaps in front of you! You are rescued!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_VICT)
		act(libc.CString("@g$n@G leaps in front of @g$N@G and rescues $M!@n"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_NOTVICT)
		stop_fighting(opponent)
		hurt(0, 0, ch, opponent, nil, int64(rand_number(1, GET_LEVEL(ch))), 1)
		hurt(0, 0, opponent, ch, nil, int64(rand_number(1, GET_LEVEL(ch))), 1)
		return
	}
}
func do_assist(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg      [2048]byte
		helpee   *char_data
		opponent *char_data
	)
	if ch.Fighting != nil {
		send_to_char(ch, libc.CString("You're already fighting!  How can you assist someone else?\r\n"))
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Whom do you wish to assist?\r\n"))
	} else if (func() *char_data {
		helpee = get_char_vis(ch, &arg[0], nil, 1<<0)
		return helpee
	}()) == nil {
		send_to_char(ch, libc.CString("%s"), config_info.Play.NOPERSON)
	} else if helpee == ch {
		send_to_char(ch, libc.CString("You can't help yourself any more than this!\r\n"))
	} else {
		if helpee.Fighting != nil {
			opponent = helpee.Fighting
		} else {
			for opponent = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; opponent != nil && opponent.Fighting != helpee; opponent = opponent.Next_in_room {
			}
		}
		if opponent == nil {
			act(libc.CString("But nobody is fighting $M!"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_CHAR)
		} else if !CAN_SEE(ch, opponent) {
			act(libc.CString("You can't see who is fighting $M!"), TRUE, ch, nil, unsafe.Pointer(helpee), TO_CHAR)
		} else {
			reveal_hiding(ch, 0)
			send_to_char(ch, libc.CString("You join the fight!\r\n"))
			act(libc.CString("$N assists you!"), 0, helpee, nil, unsafe.Pointer(ch), TO_CHAR)
			act(libc.CString("$n assists $N."), TRUE, ch, nil, unsafe.Pointer(helpee), TO_NOTVICT)
			if ch.Fighting == nil {
				set_fighting(ch, opponent)
			}
			if opponent.Fighting == nil {
				set_fighting(opponent, ch)
			}
		}
	}
}
func do_kill(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg  [2048]byte
		vict *char_data
	)
	if IS_NPC(ch) || !ADM_FLAGGED(ch, ADM_INSTANTKILL) {
		return
	}
	one_argument(argument, &arg[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Kill who?\r\n"))
	} else {
		if (func() *char_data {
			vict = get_char_vis(ch, &arg[0], nil, 1<<0)
			return vict
		}()) == nil {
			send_to_char(ch, libc.CString("They aren't here.\r\n"))
		} else if ch == vict {
			send_to_char(ch, libc.CString("Your mother would be so sad.. :(\r\n"))
		} else {
			act(libc.CString("You chop $M to pieces!  Ah!  The blood!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("$N chops you to pieces!"), TRUE, vict, nil, unsafe.Pointer(ch), TO_CHAR)
			act(libc.CString("$n brutally slays $N!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			raw_kill(vict, ch)
		}
	}
}
func do_flee(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		i            int
		attempt      int
		was_fighting *char_data
	)
	_ = was_fighting
	var arg [2048]byte
	one_argument(argument, &arg[0])
	if int(ch.Position) < POS_RESTING {
		send_to_char(ch, libc.CString("You are in pretty bad shape, unable to flee!\r\n"))
		return
	}
	if ch.Grappling != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return
	}
	if ch.Grappled != nil {
		send_to_char(ch, libc.CString("You are grappling with someone!\r\n"))
		return
	}
	if ch.Absorbing != nil {
		send_to_char(ch, libc.CString("You are absorbing from someone!\r\n"))
		return
	}
	if ch.Absorbby != nil {
		send_to_char(ch, libc.CString("You are being absorbed from by someone!\r\n"))
		return
	}
	if !IS_NPC(ch) {
		var (
			fail     int = FALSE
			obj      *obj_data
			next_obj *obj_data
		)
		for obj = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; obj != nil; obj = next_obj {
			next_obj = obj.Next_content
			if obj.Kicharge > 0 && obj.User == ch {
				fail = TRUE
			}
		}
		if fail == TRUE {
			send_to_char(ch, libc.CString("You are too busy controlling your attack!\r\n"))
			return
		}
	}
	for i = 0; i < 12; i++ {
		if arg[0] != 0 {
			if (func() int {
				attempt = int(libc.BoolToInt(search_block(&arg[0], &dirs[0], FALSE) > -1))
				return attempt
			}()) != 0 {
				attempt = search_block(&arg[0], &dirs[0], FALSE)
			} else if (func() int {
				attempt = int(libc.BoolToInt(search_block(&arg[0], &abbr_dirs[0], FALSE) > -1))
				return attempt
			}()) != 0 {
				attempt = search_block(&arg[0], &abbr_dirs[0], FALSE)
			} else {
				attempt = rand_number(0, int(NUM_OF_DIRS-1))
			}
		}
		if arg[0] == 0 {
			attempt = rand_number(0, int(NUM_OF_DIRS-1))
		}
		if CAN_GO(ch, attempt) && !ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[attempt]).To_room, ROOM_DEATH) {
			act(libc.CString("$n panics, and attempts to flee!"), TRUE, ch, nil, nil, TO_ROOM)
			if IS_NPC(ch) && ROOM_FLAGGED(((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dir_option[attempt]).To_room, ROOM_NOMOB) {
				return
			}
			was_fighting = ch.Fighting
			var wall *obj_data
			for wall = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; wall != nil; wall = wall.Next_content {
				if GET_OBJ_VNUM(wall) == 79 {
					if wall.Cost == attempt {
						send_to_char(ch, libc.CString("That direction has a glacial wall blocking it.\r\n"))
						return
					}
				}
			}
			if block_calc(ch) == 0 {
				return
			}
			if ch.Absorbing != nil {
				send_to_char(ch, libc.CString("You are busy absorbing from %s!\r\n"), GET_NAME(ch.Absorbing))
				return
			}
			if ch.Absorbby != nil {
				if axion_dice(0) < GET_SKILL(ch.Absorbing, SKILL_ABSORB) {
					send_to_char(ch, libc.CString("You are being held by %s, they are absorbing you!\r\n"), GET_NAME(ch.Absorbby))
					send_to_char(ch.Absorbby, libc.CString("%s struggles in your grasp!\r\n"), GET_NAME(ch))
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
					return
				} else {
					act(libc.CString("@c$N@W manages to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_NOTVICT)
					act(libc.CString("@WYou manage to break loose of @C$n's@W hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_VICT)
					act(libc.CString("@c$N@W manages to break loose of your hold!@n"), TRUE, ch.Absorbby, nil, unsafe.Pointer(ch), TO_CHAR)
					ch.Absorbby.Absorbing = nil
					ch.Absorbby = nil
				}
			}
			if do_simple_move(ch, attempt, TRUE) != 0 {
				send_to_char(ch, libc.CString("You flee head over heels.\r\n"))
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			} else {
				act(libc.CString("$n tries to flee, but can't!"), TRUE, ch, nil, nil, TO_ROOM)
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
			}
			return
		}
	}
	send_to_char(ch, libc.CString("PANIC!  You couldn't escape!\r\n"))
}
