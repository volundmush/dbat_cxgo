package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func do_lightgrenade(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc  int
		dge   int = 2
		count int = 0
	)
	_ = count
	var skill int
	var dmg int64
	var attperc float64 = 0.35
	var minimum float64 = 0.15
	var vict *char_data = nil
	var targ *char_data = nil
	var next_v *char_data = nil
	var arg [2048]byte
	var arg2 [2048]byte
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_LIGHTGRENADE) == 0 {
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
	targ = nil
	if arg[0] == 0 || (func() *char_data {
		targ = get_char_vis(ch, &arg[0], nil, 1<<0)
		return targ
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			targ = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("Nobody around here by that name.\r\n"))
			return
		}
	}
	skill = init_skill(ch, SKILL_LIGHTGRENADE)
	perc = chance_to_hit(ch)
	if skill < perc {
		act(libc.CString("@WYou quickly bring your hands in front of your body and cup them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as your ki is condensed between your hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. Your concentration waivers however and the energy slips away from your control harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_CHAR)
		act(libc.CString("@C$n@W quickly brings $s hands in front of $s body and cups them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as ki is condensed between $s hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. Suddenly $s concentration seems to waiver and the energy slips away from $s control harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_VICT)
		act(libc.CString("@C$n@W quickly brings $s hands in front of $s body and cups them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as ki is condensed between $s hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. Suddenly $s concentration seems to waiver and the energy slips away from $s control harmlessly!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_NOTVICT)
		pcost(ch, attperc, 0)
		improve_skill(ch, SKILL_LIGHTGRENADE, 0)
		return
	}
	act(libc.CString("@WYou quickly bring your hands in front of your body and cup them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as your ki is condensed between your hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. You shout @r'@YLIGHT GRENADE@r'@W as the orb launches from your hands at @C$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_CHAR)
	act(libc.CString("@C$n@W quickly brings $s hands in front of $s body and cups them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as ki is condensed between $s hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. @C$n shouts @r'@YLIGHT GRENADE@r'@W as the orb launches from $s hands at YOU!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_VICT)
	act(libc.CString("@C$n@W quickly brings $s hands in front of $s body and cups them a short distance from each other. A flash of @Ggreen @Ylight@W can be seen as ki is condensed between $s hands before a @Yg@yo@Yl@yd@Ye@yn@W orb of ki replaces the green light. @C$n shouts @r'@YLIGHT GRENADE@r'@W as the orb launches from $s hands at @c$N@W!@n"), TRUE, ch, nil, unsafe.Pointer(targ), TO_NOTVICT)
	dmg = damtype(ch, 57, skill, attperc)
	for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; vict != nil; vict = next_v {
		next_v = vict.Next_in_room
		if vict == ch {
			continue
		}
		if MOB_FLAGGED(vict, MOB_NOKILL) {
			continue
		}
		if AFF_FLAGGED(vict, AFF_GROUP) {
			if vict.Master == ch || ch.Master == vict || ch.Master == vict.Master {
				if vict == targ {
					send_to_char(ch, libc.CString("Leave the group if you want to murder them.\r\n"))
				}
				continue
			}
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
			if vict == targ {
				pcost(ch, attperc, 0)
				handle_cooldown(ch, 9)
				return
			} else {
				continue
			}
		} else if dge > axion_dice(int(float64(skill)*0.5)) && vict == targ {
			send_to_char(ch, libc.CString("DGE: %d\n"), dge)
			act(libc.CString("@c$N@W manages to dodge the light grenade!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@WYou manages to dodge the light grenade!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N@W manages to dodge the light grenade!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			hurt(0, 0, ch, vict, nil, 0, 1)
			improve_skill(vict, SKILL_DODGE, 0)
			pcost(ch, attperc, 0)
			handle_cooldown(ch, 9)
			return
		} else if dge > axion_dice(int(float64(skill)*0.5)) && vict != targ {
			act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@WYou manage to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			hurt(0, 0, ch, vict, nil, 0, 1)
			improve_skill(vict, SKILL_DODGE, 0)
			continue
		} else if vict == targ {
			act(libc.CString("@R$N@r is hit by the light grenade which explodes all around $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RYou are hit by the light grenade which explodes all around you!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$N@r is hit by the light grenade which explodes all around $m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 4) == 4 {
				handle_knockdown(vict)
			}
			hurt(0, 0, ch, vict, nil, dmg, 1)
			continue
		} else {
			act(libc.CString("@R$N@r is caught by the light grenade's explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RYou are caught by the light grenade's explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@R$N@r is caught by the light grenade's explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if !AFF_FLAGGED(vict, AFF_FLYING) && int(vict.Position) == POS_STANDING && rand_number(1, 4) == 4 {
				handle_knockdown(vict)
			}
			hurt(0, 0, ch, vict, nil, int64(float64(dmg)*0.5), 1)
			continue
		}
	}
	pcost(ch, attperc, 0)
	improve_skill(ch, SKILL_LIGHTGRENADE, 0)
	handle_cooldown(ch, 9)
}
func do_energize(ch *char_data, argument *byte, cmd int, subcmd int) {
	if IS_NPC(ch) {
		return
	}
	if ch.Preference != PREFERENCE_THROWING {
		send_to_char(ch, libc.CString("You aren't dedicated to throwing!\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_ENERGIZE) == 0 {
		if GET_SKILL(ch, SKILL_FOCUS) >= 30 {
			var result int = rand_number(10, 14)
			for {
				ch.Skills[SKILL_ENERGIZE] = int8(result)
				if true {
					break
				}
			}
			send_to_char(ch, libc.CString("You learn the basics for energizing thrown weapons! Now use the energize command again.\r\n"))
			return
		} else {
			send_to_char(ch, libc.CString("You need a Focus skill level of 30 to figure out the basics of this technique.\r\n"))
			return
		}
	} else {
		if PRF_FLAGGED(ch, PRF_ENERGIZE) {
			send_to_char(ch, libc.CString("You stop focusing ki into your fingertips.\r\n"))
			ch.Player_specials.Pref[int(PRF_ENERGIZE/32)] &= bitvector_t(int32(^(1 << (int(PRF_ENERGIZE % 32)))))
			return
		} else {
			send_to_char(ch, libc.CString("You start focusing your latent ki into your fingertips.\r\n"))
			ch.Player_specials.Pref[int(PRF_ENERGIZE/32)] |= bitvector_t(int32(1 << (int(PRF_ENERGIZE % 32))))
			return
		}
	}
}
func do_breath(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		stcost  int64 = ch.Max_hit / 5000
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
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
			return
		}
	}
	handle_cooldown(ch, 10)
	if vict != nil {
		if can_kill(ch, vict, nil, 1) == 0 {
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
		prob += 15
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@cYou disappear, avoiding @C$n's@c breath before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c breath before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@WYou move quickly and block @C$n's@W fiery breath!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W fiery breath!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 6, skill, attperc)
					dmg = int64(float64(dmg) * 0.8)
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
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@WYou dodge the fiery jets of flames coming from @C$n's@W mouth!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge the fiery jets of flames coming from @c$n's@W mouth!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@C$n@W moves to breath flames on you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to breath flames on @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@C$n@W moves to breath flames on you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to breath flames on @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 8, skill, attperc)
			dmg += dmg * 2
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@C$n@W aims $s mouth at you and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s mouth at @c$N@W and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto @c$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 2:
				act(libc.CString("@C$n@W aims $s mouth at you and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s mouth at @c$N@W and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto @c$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 3:
				act(libc.CString("@C$n@W aims $s mouth at you and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s mouth at @c$N@W and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto @c$N's@W body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@C$n@W aims $s mouth at you and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s mouth at @c$N@W and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto @c$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@C$n@W aims $s mouth at you and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W aims $s mouth at @c$N@W and opens it wide slowly. A high pitched sound can be heard as the mouth opens, and as the throat is exposed a bright white flame can be seen burning there. Suddenly @C$n@W breathes a jet of @rf@Ri@Ye@Rr@ry@W flames onto @c$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, 0, stcost)
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
		act(libc.CString("@C$n@W breathes flames on $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_ram(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		stcost  int64 = ch.Max_hit / 200
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
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
		prob -= 5
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@cYou disappear, avoiding @C$n's@c ram before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c ram before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > axion_dice(10) {
					act(libc.CString("@WYou move quickly and block @C$n's@W body as $e tries to ram YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W body as $e tries to ram $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 6, skill, attperc)
					dmg -= int64(float64(dmg) * 0.2)
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@WYou dodge @C$n's@W attempted ram!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W attempted ram!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@C$n@W moves to ram you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to ram @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@C$n@W moves to ram you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to ram @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 8, skill, attperc)
			dmg += int64(float64(dmg) * 1.1)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@C$n@W aims $s body at YOU and rams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s body at @C$N@W and rams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 2:
				act(libc.CString("@C$n@W aims $s body at YOU and rams into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s body at @C$N@W and rams into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 3:
				act(libc.CString("@C$n@W aims $s body at YOU and rams into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s body at @C$N@W and rams into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@C$n@W aims $s body at YOU and rams into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s body at @C$N@W and rams into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 190, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@C$n@W aims $s body at YOU and rams into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W aims $s body at @C$N@W and rams into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 190, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, 0, stcost)
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
		act(libc.CString("@C$n@W rams $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_strike(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		stcost  int64 = ch.Max_hit / 400
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if !IS_NPC(ch) {
		return
	}
	if can_grav(ch) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		return
	}
	if check_points(ch, 0, ch.Max_hit/400) == 0 {
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
		prob += 5
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@cYou disappear, avoiding @C$n's@c fang strike before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c fang strike before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, 0, stcost/2)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if pry > rand_number(1, 140) && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					act(libc.CString("@WYou parry @C$n's@W fang strike with a punch of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W fang strike with a punch of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > axion_dice(10) {
					act(libc.CString("@WYou move quickly and block @C$n's@W fang strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W fang strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 6, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > axion_dice(10) {
					act(libc.CString("@WYou dodge @C$n's@W fang strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W fang strike!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@C$n@W moves to fang strike you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W moves to fang strike @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@C$n@W moves to fang strike you, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W moves to fang strike @C$N@W, but somehow misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, 0, stcost/2)
			}
			hurt(0, 0, ch, vict, nil, 0, 0)
			return
		} else {
			dmg = damtype(ch, 8, skill, attperc)
			dmg += int64(float64(dmg) * 0.5)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@C$n@W launches $s body at YOU and sinks $s fang strike into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W launches $s body at @C$N@W and sinks $s fang strike into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 2:
				act(libc.CString("@C$n@W launches $s body at YOU and sinks $s fang strike into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W launches $s body at @C$N@W and sinks $s fang strike into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 3:
				act(libc.CString("@C$n@W launches $s body at YOU and sinks $s fang strike into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W launches $s body at @C$N@W and sinks $s fang strike into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@C$n@W launches $s body at YOU and sinks $s fang strike into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W launches $s body at @C$N@W and sinks $s fang strike into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@C$n@W launches $s body at YOU and sinks $s fang strike into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W launches $s body at @C$N@W and sinks $s fang strike into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			}
			pcost(ch, 0, stcost)
			vict.Move -= int64(float64(dmg) * 0.25)
			if vict.Move < 0 {
				vict.Move = 0
			}
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
		act(libc.CString("@C$n@W fang strikes $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_combine(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		arg   [2048]byte
		arg2  [2048]byte
		vict  *char_data
		f     *follow_type
		fire  int = FALSE
		temp  int = -1
		temp2 int = -1
	)
	_ = temp2
	two_arguments(argument, &arg[0], &arg2[0])
	if has_group(ch) == 0 {
		send_to_char(ch, libc.CString("You need to be in a group!\r\n"))
		return
	} else {
		if arg[0] == 0 || ch.Master == nil && arg2[0] == 0 {
			send_to_char(ch, libc.CString("Follower Syntax: combine (attack)\r\n"))
			send_to_char(ch, libc.CString("Leader Syntax: combine (attack) (target)\r\n"))
			send_to_char(ch, libc.CString("Cancel Syntax: combine stop\r\n"))
		} else {
			if libc.StrCaseCmp(&arg[0], libc.CString("stop")) == 0 && ch.Master != nil {
				if ch.Combine == -1 {
					send_to_char(ch, libc.CString("You are not trying to combine any attacks...\r\n"))
					return
				} else {
					send_to_char(ch, libc.CString("You stop your preparations to combine your attack with a group attack.\r\n"))
					send_to_char(ch.Master, libc.CString("@Y%s@C is no longer prepared to combine an attack with the group!@n\r\n"), get_i_name(ch.Master, ch))
					for f = ch.Master.Followers; f != nil; f = f.Next {
						if ch != f.Follower {
							send_to_char(f.Follower, libc.CString("@Y%s@C is no longer prepared to combine an attack with the group!@n\r\n"), get_i_name(f.Follower, ch))
						}
					}
					ch.Combine = -1
					return
				}
			} else if libc.StrCaseCmp(&arg[0], libc.CString("stop")) == 0 && ch.Master == nil {
				send_to_char(ch, libc.CString("You do not need to stop as you haven't prepared anything.\r\n"))
				return
			}
			var i int = 0
			for i = 0; i < 14; i++ {
				if libc.StrStr(&arg[0], attack_names[i]) != nil {
					if i == 5 {
						if (ch.Equipment[WEAR_WIELD1]) == nil {
							send_to_char(ch, libc.CString("You need to wield a sword to use this technique.\r\n"))
							return
						}
						if ((ch.Equipment[WEAR_WIELD1]).Value[VAL_WEAPON_DAMTYPE]) != int(TYPE_SLASH-TYPE_HIT) {
							send_to_char(ch, libc.CString("You are not wielding a sword, you need one to use this technique.\r\n"))
							return
						} else {
							temp = i
							i = 15
						}
					}
					temp = i
					i = 15
				}
			}
			if temp == -1 {
				send_to_char(ch, libc.CString("Follower Syntax: combine (attack)\r\n"))
				send_to_char(ch, libc.CString("Leader Syntax: combine (attack) (target)\r\n"))
				send_to_char(ch, libc.CString("Follower Cancel Syntax: combine stop\r\n"))
				return
			} else if GET_SKILL(ch, attack_skills[temp]) == 0 {
				send_to_char(ch, libc.CString("You do not know that skill.\r\n"))
				return
			} else if attack_skills[temp] == 440 && int(ch.Chclass) != CLASS_NAIL {
				send_to_char(ch, libc.CString("Only students of Nail know how to combine that attack effectively.\r\n"))
				return
			} else if float64(ch.Charge) < float64(ch.Max_mana)*0.05 {
				send_to_char(ch, libc.CString("You need to have the minimum of 5%s ki charged to combine.\r\n"), "%")
			}
			if ch.Master == nil {
				if (func() *char_data {
					vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
					return vict
				}()) == nil {
					send_to_char(ch, libc.CString("Who will your combined attack be targeting?\r\n"))
					return
				} else if vict == ch {
					send_to_char(ch, libc.CString("No targeting yourself...\r\n"))
					return
				}
				ch.Combine = temp
				for f = ch.Followers; f != nil; f = f.Next {
					if !AFF_FLAGGED(f.Follower, AFF_GROUP) {
						continue
					} else if f.Follower.Combine != -1 && float64(f.Follower.Charge) >= float64(f.Follower.Max_mana)*0.05 {
						fire = TRUE
					}
				}
				if fire == TRUE {
					combine_attacks(ch, vict)
					return
				} else {
					send_to_char(ch, libc.CString("You do not have any followers who have readied an attack to combine or they do not have enough ki anymore to combine said attack.\r\n"))
					return
				}
			} else if ch.Master != nil {
				if float64(ch.Charge) >= float64(ch.Max_mana)*0.05 {
					act(libc.CString("@C$n@c appears to be concentrating hard and focusing $s energy!@n\r\n"), TRUE, ch, nil, nil, TO_ROOM)
					send_to_char(ch.Master, libc.CString("@BCOMBINE@c: @Y%s@C has prepared to combine a @c'@G%s@c'@C with the next group attack!@n\r\n"), get_i_name(ch.Master, ch), attack_names[temp])
					for f = ch.Master.Followers; f != nil; f = f.Next {
						if ch != f.Follower {
							send_to_char(f.Follower, libc.CString("@BCOMBINE@c: @Y%s@C has prepared to combine a @c'@G%s@c'@C with the next group attack!@n\r\n"), get_i_name(f.Follower, ch), attack_names[temp])
						}
					}
					ch.Combine = temp
				} else {
					send_to_char(ch, libc.CString("You do not have the minimum 5%s ki charged.\r\n"), "%")
					return
				}
			}
		}
	}
}
func do_sunder(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if check_skill(ch, SKILL_SUNDER) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_SUNDER)
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
		improve_skill(ch, SKILL_SUNDER, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		perc = axion_dice(0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Sundering Force before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Sundering Force before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Sundering Force before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Sundering Force but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Sundering Force but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Sundering Force but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Sundering Force, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Sundering Force, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Sundering Force, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Sundering Force misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Sundering Force at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Sundering Force at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Sundering Force misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Sundering Force at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Sundering Force at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			var chance int = rand_number(2, 12)
			dmg = damtype(ch, 55, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise both hands and aim the flats of your palms toward @c$N@W. As you concentrate your charged ki long arcing beams of blue energy shoot out and form a field around $M. With a quick motion you move your hands in opposite directions and wrench @c$N's@W body with the force of your energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward YOU. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around YOU. With a quick motion $e moves $s hands in opposite directions and wrenches YOUR body with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward @c$N@W. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around @c$N@W. With a quick motion @C$n@W moves $s hands in opposite directions and wrenches @c$N's@W body with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				if chance >= 10 {
					hurt(0, 160, ch, vict, nil, dmg, 1)
				} else if chance >= 8 {
					hurt(1, 160, ch, vict, nil, dmg, 1)
				} else {
					hurt(0, 0, ch, vict, nil, dmg, 1)
				}
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou raise both hands and aim the flats of your palms toward @c$N@W. As you concentrate your charged ki long arcing beams of blue energy shoot out and form a field around $M. With a quick motion you move your hands in opposite directions and wrench @c$N's@W head with the force of your energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward YOU. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around YOU. With a quick motion $e moves $s hands in opposite directions and wrenches YOUR head with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward @c$N@W. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around @c$N@W. With a quick motion @C$n@W moves $s hands in opposite directions and wrenches @c$N's@W head with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou raise both hands and aim the flats of your palms toward @c$N@W. As you concentrate your charged ki long arcing beams of blue energy shoot out and form a field around $M. With a quick motion you move your hands in opposite directions and wrench @c$N's@W body with the force of your energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward YOU. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around YOU. With a quick motion $e moves $s hands in opposite directions and wrenches YOUR body with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward @c$N@W. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around @c$N@W. With a quick motion @C$n@W moves $s hands in opposite directions and wrenches @c$N's@W body with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				if chance >= 10 {
					hurt(0, 160, ch, vict, nil, dmg, 1)
				} else if chance >= 8 {
					hurt(1, 160, ch, vict, nil, dmg, 1)
				} else {
					hurt(0, 0, ch, vict, nil, dmg, 1)
				}
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou raise both hands and aim the flats of your palms toward @c$N@W. As you concentrate your charged ki long arcing beams of blue energy shoot out and form a field around $M. With a quick motion you move your hands in opposite directions and wrench @c$N's@W arm with the force of your energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward YOU. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around YOU. With a quick motion $e moves $s hands in opposite directions and wrenches YOUR arm with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward @c$N@W. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around @c$N@W. With a quick motion @C$n@W moves $s hands in opposite directions and wrenches @c$N's@W arm with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 160, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WYou raise both hands and aim the flats of your palms toward @c$N@W. As you concentrate your charged ki long arcing beams of blue energy shoot out and form a field around $M. With a quick motion you move your hands in opposite directions and wrench @c$N's@W leg with the force of your energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward YOU. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around YOU. With a quick motion $e moves $s hands in opposite directions and wrenches YOUR leg with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises both hands and aims the flats of $s palms toward @c$N@W. As $e concentrates $s charged ki long arcing beams of blue energy shoot out and form a field around @c$N@W. With a quick motion @C$n@W moves $s hands in opposite directions and wrenches @c$N's@W leg with the force of $s energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(1, 160, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
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
		act(libc.CString("@WYou fire a Sundering Force at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Sundering Force at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_zen(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if check_skill(ch, SKILL_ZEN) == 0 {
		return
	}
	if !HAS_ARMS(ch) {
		send_to_char(ch, libc.CString("You have no available arms!\r\n"))
		return
	} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 && (ch.Limb_condition[2]) < 0 {
		send_to_char(ch, libc.CString("Using your broken right arm has damaged it more!@n\r\n"))
		ch.Limb_condition[1] -= rand_number(3, 5)
		if (ch.Limb_condition[1]) < 0 {
			act(libc.CString("@RYour right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@R's right arm has fallen apart!@n"), TRUE, ch, nil, nil, TO_ROOM)
		}
	} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 && (ch.Limb_condition[1]) < 0 {
		send_to_char(ch, libc.CString("Using your broken left arm has damaged it more!@n\r\n"))
		ch.Limb_condition[2] -= rand_number(3, 5)
		if (ch.Limb_condition[2]) < 0 {
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
	if int(ch.Skillperfs[SKILL_ZEN]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_ZEN]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_ZEN)
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
		improve_skill(ch, SKILL_ZEN, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_ZEN]) == 2 {
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
				act(libc.CString("@C$N@c disappears, avoiding your Zen Blade Strike before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Zen Blade Strike before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Zen Blade Strike before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
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
				if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W moves quickly and blocks your Zen Blade Strike!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Zen Blade Strike!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Zen Blade Strike!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 54, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Zen Blade Strike, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Zen Blade Strike, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Zen Blade Strike, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 19, skill, SKILL_ZEN)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					improve_skill(vict, SKILL_DODGE, 0)
					if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Zen Blade Strike misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Zen Blade Strike at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Zen Blade Strike at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Zen Blade Strike misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Zen Blade Strike at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Zen Blade Strike at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 54, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@CRaising your blade above your head, and closing your eyes, you focus ki into its edge. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as you open your eyes and instantly move past @g$N's@C body while slashing with the pure energy of your resolve! A large explosion of energy erupts across $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past YOUR body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past @g$N's@C body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across @g$N's@C body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
				act(libc.CString("@CRaising your blade above your head, and closing your eyes, you focus ki into its edge. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as you open your eyes and instantly move past @g$N's@C body while slashing with the pure energy of your resolve! A large explosion of energy erupts across $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past YOUR body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past @g$N's@C body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across @g$N's@C head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
				act(libc.CString("@CRaising your blade above your head, and closing your eyes, you focus ki into its edge. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as you open your eyes and instantly move past @g$N's@C body while slashing with the pure energy of your resolve! A large explosion of energy erupts across $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past YOUR body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past @g$N's@C body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across @g$N's@C body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@CRaising your blade above your head, and closing your eyes, you focus ki into its edge. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as you open your eyes and instantly move past @g$N's@C body while slashing with the pure energy of your resolve! A large explosion of energy erupts across $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past YOUR body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past @g$N's@C body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across @g$N's@C arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 80 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[2]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severs $N's left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left arm is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[2] = 0
						remove_limb(vict, 2)
					} else if (vict.Limb_condition[1]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severs $N's right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right arm is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[1] = 0
						remove_limb(vict, 1)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 5:
				act(libc.CString("@CRaising your blade above your head, and closing your eyes, you focus ki into its edge. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as you open your eyes and instantly move past @g$N's@C body while slashing with the pure energy of your resolve! A large explosion of energy erupts across $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past YOUR body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@G$n @Craises $s blade above $s head, and closes $s eyes. The edge of the blade begins to glow a soft blue as the blade begins to throb with excess energy. Peels of lighting begin to arc from the blade in all directions as $e opens $s eyes and instantly moves past @g$N's@C body while slashing with the pure energy of $s resolve! A large explosion of energy erupts across @g$N's@C leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				if rand_number(1, 100) >= 80 && !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
					if (vict.Limb_condition[4]) > 0 && !is_sparring(ch) && rand_number(1, 2) == 2 {
						act(libc.CString("@RYour attack severs $N's left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's left leg is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[4] = 0
						remove_limb(vict, 4)
					} else if (vict.Limb_condition[3]) > 0 && !is_sparring(ch) {
						act(libc.CString("@RYour attack severs $N's right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@R$n's attack severs your right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@R$N's right leg is severered in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						vict.Limb_condition[3] = 0
						remove_limb(vict, 3)
					}
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 2)
			}
			if int(ch.Skillperfs[SKILL_ZEN]) == 3 && attperc > minimum {
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
		act(libc.CString("@WYou fire a Zen Blade Strike at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Zen Blade Strike at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_malice(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if check_skill(ch, SKILL_MALICE) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_MALICE)
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
		improve_skill(ch, SKILL_MALICE, 0)
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
		if time_info.Hours <= 15 || time_info.Hours > 22 {
			prob += 5
		}
		if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Malice Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Malice Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Malice Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
				if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W moves quickly and blocks your Malice Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Malice Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Malice Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 36, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Malice Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Malice Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Malice Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 36, skill, SKILL_MALICE)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 80 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 20
					}
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Malice Breaker misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Malice Breaker at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Malice Breaker at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Malice Breaker misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Malice Breaker at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Malice Breaker at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 36, skill, attperc)
			if time_info.Hours <= 15 {
				dmg *= int64(1.25)
			} else if time_info.Hours <= 22 {
				dmg *= int64(1.4)
			}
			switch rand_number(1, 6) {
			case 1:
				act(libc.CString("@WYou rush forward at @c$N@W, building ki into your arm. As you slam an open palm into $S chest, you send the charged energy into $S body. A few small explosions seem to hit across $S entire body, forcing $M to stumble back. Finally, you launch $M into the air, pointing a forefinger at $m like a pistol, and shout '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of your first strike on @c$N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wrushes forward at you, building ki into $s arm. As $e slams an open palm into your chest, $e sends the charged energy into your body. A few small explosions seem to hit across your entire body, forcing you to stumble back. Finally, @C$n@W launches you into the air, pointing a forefinger at your body like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W rushes forward at $N@W, building ki into $s arm. As $e slam an open palm into @c$N's@W chest, $e sends the charged energy into $S body. A few small explosions seem to hit across @c$N's@W entire body, forcing $M to stumble back. Finally, @C$n@W launches $M into the air, pointing a forefinger at $m like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on @c$N@W's chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				fallthrough
			case 3:
				act(libc.CString("@WYou rush forward at @c$N@W, building ki into your arm. As you slam an open palm into $S head, you send the charged energy into $S body. A few small explosions seem to hit across $S entire body, forcing $M to stumble back. Finally, you launch $M into the air, pointing a forefinger at $m like a pistol, and shout '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of your first strike on @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wrushes forward at you, building ki into $s arm. As $e slams an open palm into your head, $e sends the charged energy into your body. A few small explosions seem to hit across your entire body, forcing you to stumble back. Finally, @C$n@W launches you into the air, pointing a forefinger at your body like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on your head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W rushes forward at $N@W, building ki into $s arm. As $e slam an open palm into @c$N's@W head, $e sends the charged energy into $S body. A few small explosions seem to hit across @c$N's@W entire body, forcing $M to stumble back. Finally, @C$n@W launches $M into the air, pointing a forefinger at $m like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on @c$N@W's head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 4:
				act(libc.CString("@WYou rush forward at @c$N@W, building ki into your arm. As you slam an open palm into $S gut, you send the charged energy into $S body. A few small explosions seem to hit across $S entire body, forcing $M to stumble back. Finally, you launch $M into the air, pointing a forefinger at $m like a pistol, and shout '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of your first strike on @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wrushes forward at you, building ki into $s arm. As $e slams an open palm into your gut, $e sends the charged energy into your body. A few small explosions seem to hit across your entire body, forcing you to stumble back. Finally, @C$n@W launches you into the air, pointing a forefinger at your body like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on your gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W rushes forward at $N@W, building ki into $s arm. As $e slam an open palm into @c$N's@W gut, $e sends the charged energy into $S body. A few small explosions seem to hit across @c$N's@W entire body, forcing $M to stumble back. Finally, @C$n@W launches $M into the air, pointing a forefinger at $m like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on @c$N@W's gut!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 5:
				act(libc.CString("@WYou rush forward at @c$N@W, building ki into your arm. As you slam an open palm into $S arm, you send the charged energy into $S body. A few small explosions seem to hit across $S entire body, forcing $M to stumble back. Finally, you launch $M into the air, pointing a forefinger at $m like a pistol, and shout '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of your first strike on @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wrushes forward at you, building ki into $s arm. As $e slams an open palm into your arm, $e sends the charged energy into your body. A few small explosions seem to hit across your entire body, forcing you to stumble back. Finally, @C$n@W launches you into the air, pointing a forefinger at your body like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W rushes forward at $N@W, building ki into $s arm. As $e slam an open palm into @c$N's@W arm, $e sends the charged energy into $S body. A few small explosions seem to hit across @c$N's@W entire body, forcing $M to stumble back. Finally, @C$n@W launches $M into the air, pointing a forefinger at $m like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on @c$N@W's arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 170, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 6:
				act(libc.CString("@WYou rush forward at @c$N@W, building ki into your arm. As you slam an open palm into $S leg, you send the charged energy into $S body. A few small explosions seem to hit across $S entire body, forcing $M to stumble back. Finally, you launch $M into the air, pointing a forefinger at $m like a pistol, and shout '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of your first strike on @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wrushes forward at you, building ki into $s arm. As $e slams an open palm into your leg, $e sends the charged energy into your body. A few small explosions seem to hit across your entire body, forcing you to stumble back. Finally, @C$n@W launches you into the air, pointing a forefinger at your body like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W rushes forward at $N@W, building ki into $s arm. As $e slam an open palm into @c$N's@W leg, $e sends the charged energy into $S body. A few small explosions seem to hit across @c$N's@W entire body, forcing $M to stumble back. Finally, @C$n@W launches $M into the air, pointing a forefinger at $m like a pistol, and shouts '@MM@ma@Dl@wi@Wce Br@we@Dak@me@Mr@W!' as a dark, violet explosion erupts at the epicenter of $s first strike on @c$N@W's leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
		act(libc.CString("@WYou fire a Malice Breaker at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Malice Breaker at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_nova(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		perc    int
		dge     int = 2
		count   int = 0
		skill   int
		dmg     int64
		attperc float64    = 0.2
		minimum float64    = 0.1
		vict    *char_data = nil
		next_v  *char_data = nil
		arg2    [2048]byte
	)
	one_argument(argument, &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_STARNOVA) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if int(ch.Skillperfs[SKILL_STARNOVA]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_STARNOVA]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_STARNOVA)
	if int(ch.Skillperfs[SKILL_STARNOVA]) == 2 {
		skill += 5
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
		if time_info.Hours <= 15 && time_info.Hours > 22 {
			skill += 5
		}
		handle_cooldown(ch, 6)
		if skill < perc {
			act(libc.CString("@WYou gather your charged energy and clench your upheld fists at either side of your body while crouching down. A hot glow of energy begins to form around your body before you lose your concentration and fail to create a @yS@Yt@Wa@wr @cN@Co@Wv@wa@W!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@C$n@W gathers $s charged energy and clenches $s upheld fists at either side of $s body while crouching down. A hot glow of energy begins to form around $s body before $e seems to lose $s concentration and fail to create a @yS@Yt@Wa@wr @cN@Co@Wv@wa@W!@n"), TRUE, ch, nil, nil, TO_ROOM)
			if int(ch.Skillperfs[SKILL_STARNOVA]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			improve_skill(ch, SKILL_STARNOVA, 0)
			return
		}
		act(libc.CString("@WYou gather your charged energy and clench your upheld fists at either side of your body while crouching down. A hot glow of energy begins to form around your body in the shape of a sphere! Suddenly a shockwave of heat and energy erupts out into the surrounding area as your glorious @yS@Yt@Wa@wr @cN@Co@Wv@wa@W is born!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W gathers $s charged energy and clenches $s upheld fists at either side of $s body while crouching down. A hot glow of energy begins to form around $s body in the shape of a sphere! Suddenly a shockwave of heat and energy erupts out into the surrounding area as @C$n's@W glorious @yS@Yt@Wa@wr @cN@Co@Wv@wa@W is born!@n"), TRUE, ch, nil, nil, TO_ROOM)
		dmg = damtype(ch, 53, skill, attperc)
		if time_info.Hours <= 15 {
			dmg *= int64(1.25)
		} else if time_info.Hours <= 22 {
			dmg *= int64(1.4)
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
		if int(ch.Skillperfs[SKILL_STARNOVA]) == 3 && attperc > minimum {
			pcost(ch, attperc-0.05, 0)
		} else {
			pcost(ch, attperc, 0)
		}
		improve_skill(ch, SKILL_STARNOVA, 0)
		return
	}
}
func do_head(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		stcost  int64 = physical_cost(ch, SKILL_HEADBUTT)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_HEADBUTT) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/100) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_HEADBUTT)
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
	if int(ch.Chclass) != CLASS_KURZAK || GET_SKILL(ch, SKILL_HEADBUTT) < 100 {
		handle_cooldown(ch, 6)
	} else if GET_SKILL(ch, SKILL_HEADBUTT) >= 100 {
		handle_cooldown(ch, 5)
	}
	if vict != nil {
		if can_kill(ch, vict, nil, 0) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_HEADBUTT, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your headbutt before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c headbutt before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c headbutt before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
					act(libc.CString("@C$N@W parries your headbutt with an attack of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W headbutt with an attack of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W headbutt with an attack of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W blocks your headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou block @C$n's@W headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W blocks @c$n's@W headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 3, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W dodges your headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W dodges @c$n's@W headbutt!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it, your headbutt misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a headbutt at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a headbutt at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it, your headbutt misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a headbutt at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a headbutt at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 0)
				pcost(ch, 0, 0)
			}
			return
		} else {
			dmg = damtype(ch, 52, skill, attperc)
			if int(ch.Chclass) == CLASS_KURZAK {
				if GET_SKILL(ch, SKILL_HEADBUTT) >= 60 {
					dmg += int64(float64(dmg) * 0.1)
				} else if GET_SKILL(ch, SKILL_HEADBUTT) >= 40 {
					dmg += int64(float64(dmg) * 0.05)
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou grab @c$N@W by the shoulders and slam your head into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU by the shoulders and slams $s head into YOUR chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W by the shoulders and slams $s head into @c$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou grab @c$N@W by the shoulders and slam your head into $S face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU by the shoulders and slams $s head into YOUR face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W by the shoulders and slams $s head into @c$N's@W face!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if !AFF_FLAGGED(vict, AFF_KNOCKED) && (rand_number(1, 7) >= 4 && vict.Hit > ch.Hit/5 && !AFF_FLAGGED(vict, AFF_SANCTUARY)) {
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
				var mult int = int(calc_critical(ch, 0))
				if int(ch.Chclass) == CLASS_KURZAK && !IS_NPC(ch) {
					if int(ch.Skills[SKILL_STYLE]) >= 75 {
						mult += 1
					}
				}
				dmg *= int64(mult)
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou grab @c$N@W by the shoulders and slam your head into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU by the shoulders and slams $s head into YOUR chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W by the shoulders and slams $s head into @c$N's@W chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou grab @c$N@W and barely manage to slam your head into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU and barely manages to slam $s head into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W and barely manages to slam $s head into @c$N's@W leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WYou grab @c$N@W and barely manage to slam your head into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU and barely manages to slam $s head into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W and barely manages to slam $s head into @c$N's@W arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) {
				act(libc.CString("@c$N's@W fireshield burns your head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W head is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W head is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
		act(libc.CString("@WYou headbutt $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W headbutt $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_bash(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		stcost  int64 = physical_cost(ch, SKILL_BASH)
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		attperc float64 = 0
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BASH) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, 0, ch.Max_hit/70) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_BASH)
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
		improve_skill(ch, SKILL_BASH, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your bash before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c bash before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c bash before reappearing!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
					act(libc.CString("@C$N@W parries your bash with an attack of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou parry @C$n's@W bash with an attack of your own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W parries @c$n's@W bash with an attack of $S own!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_PARRY, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(vict, -2, skill, attperc)
					dmg *= int64(calc_critical(ch, 1))
					hurt(0, 0, vict, ch, nil, dmg, -1)
					return
				} else if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W blocks your bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou block @C$n's@W bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W blocks @c$n's@W bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_BLOCK, 0)
					pcost(ch, 0, stcost/2)
					pcost(vict, 0, vict.Max_hit/500)
					dmg = damtype(ch, 3, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 0)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W dodges your bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W dodges @c$n's@W bash!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					improve_skill(vict, SKILL_DODGE, 0)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				} else {
					act(libc.CString("@WYou can't believe it, your bash misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W throws a bash at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W throws a bash at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, 0, stcost/2)
					hurt(0, 0, ch, vict, nil, 0, 0)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it, your bash misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W throws a bash at you, but thankfully misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W throws a bash at @C$N@W, but misses!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, 0, 0)
				pcost(ch, 0, 0)
			}
			return
		} else {
			dmg = damtype(ch, 51, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WBending over slightly you aim your body at @c$N@W and instantly launch yourself toward $M at full speed! You slam into $S body with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at YOU and then instantly launches $mself toward YOU at full speed! @C$n@W slams into YOUR body with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at @c$N@W and then instantly launches $mself toward $M at full speed! @C$n@W slams into $S body with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WBending over slightly you aim your body at @c$N@W and instantly launch yourself toward $M at full speed! You slam into $S head with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at YOU and then instantly launches $mself toward YOU at full speed! @C$n@W slams into YOUR head with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at @c$N@W and then instantly launches $mself toward $M at full speed! @C$n@W slams into $S head with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WBending over slightly you aim your body at @c$N@W and instantly launch yourself toward $M at full speed! You slam into $S gut with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at YOU and then instantly launches $mself toward YOU at full speed! @C$n@W slams into YOUR gut with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at @c$N@W and then instantly launches $mself toward $M at full speed! @C$n@W slams into $S gut with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WBending over slightly you aim your body at @c$N@W and instantly launch yourself toward $M at full speed! You slam into $S leg with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at YOU and then instantly launches $mself toward YOU at full speed! @C$n@W slams into YOUR leg with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at @c$N@W and then instantly launches $mself toward $M at full speed! @C$n@W slams into $S leg with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 2)
			case 5:
				act(libc.CString("@WBending over slightly you aim your body at @c$N@W and instantly launch yourself toward $M at full speed! You slam into $S arm with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at YOU and then instantly launches $mself toward YOU at full speed! @C$n@W slams into YOUR arm with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W bends over slightly aiming $s body at @c$N@W and then instantly launches $mself toward $M at full speed! @C$n@W slams into $S arm with a crashing impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 0)
				dam_eq_loc(vict, 1)
			}
			if vict.Hit > 0 && !AFF_FLAGGED(vict, AFF_SPIRIT) && AFF_FLAGGED(vict, AFF_FIRESHIELD) {
				act(libc.CString("@c$N's@W fireshield burns your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n's@W body is burned by your fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n's@W body is burned by @C$N's@W fireshield!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
			}
			if vict != nil && rand_number(1, 5) >= 4 {
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
			}
			if rand_number(1, 5) >= 5 {
				if AFF_FLAGGED(ch, AFF_FLYING) {
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, vict, nil, unsafe.Pointer(ch), TO_CHAR)
					act(libc.CString("@wYou are knocked out of the air!@n"), TRUE, vict, nil, unsafe.Pointer(ch), TO_VICT)
					act(libc.CString("@w$N@w is knocked out of the air!@n"), TRUE, vict, nil, unsafe.Pointer(ch), TO_NOTVICT)
					ch.Affected_by[int(AFF_FLYING/32)] &= ^(1 << (int(AFF_FLYING % 32)))
					ch.Altitude = 0
					ch.Position = POS_SITTING
				} else {
					handle_knockdown(vict)
				}
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
		act(libc.CString("@WYou bash $p@W as hard as you can!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W bash $p@W extremely hard!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, 0, stcost)
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_seishou(ch *char_data, argument *byte, cmd int, subcmd int) {
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
	if check_skill(ch, SKILL_SEISHOU) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_SEISHOU)
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
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_SEISHOU, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Seishou Enko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Seishou Enko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Seishou Enko before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc/4, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Seishou Enko but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Seishou Enko but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Seishou Enko but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Seishou Enko, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Seishou Enko, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Seishou Enko, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc/4, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Seishou Enko misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Seishou Enko at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Seishou Enko at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc/4, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Seishou Enko misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Seishou Enko at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Seishou Enko at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc/4, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 50, skill, attperc)
			if ch.Moltlevel >= 150 {
				dmg *= 2
			}
			switch rand_number(1, 7) {
			case 1:
				act(libc.CString("@WYou aim your mouth at @C$N@W and focus your charged ki. In an instant you fire a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S body with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Waims $s mouth at YOU and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at YOU! Almost instantly it blasts into YOUR body with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Waims $s mouth at @c$N@W and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S body with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				act(libc.CString("@WYou aim your mouth at @C$N@W and focus your charged ki. In an instant you fire a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S head with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Waims $s mouth at YOU and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at YOU! Almost instantly it blasts into YOUR head with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Waims $s mouth at @c$N@W and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S head with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 5:
				act(libc.CString("@WYou aim your mouth at @C$N@W and focus your charged ki. In an instant you fire a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S gut with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Waims $s mouth at YOU and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at YOU! Almost instantly it blasts into YOUR gut with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Waims $s mouth at @c$N@W and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S gut with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 6:
				act(libc.CString("@WYou aim your mouth at @C$N@W and focus your charged ki. In an instant you fire a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S arm with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Waims $s mouth at YOU and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at YOU! Almost instantly it blasts into YOUR arm with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Waims $s mouth at @c$N@W and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S arm with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 1)
			case 7:
				act(libc.CString("@WYou aim your mouth at @C$N@W and focus your charged ki. In an instant you fire a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S leg with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Waims $s mouth at YOU and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at YOU! Almost instantly it blasts into YOUR leg with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Waims $s mouth at @c$N@W and seems to focus $s ki. In an instant $e fires a large @Rr@re@Rd@W @rS@Re@Wi@ws@rh@Ro@Wu @wE@rn@Rk@Wo at $M! Almost instantly it blasts into $S leg with searing heat!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
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
		dmg = damtype(ch, 10, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire a Seishou Enko at $p@W!@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires a Seishou Enko at $p@W!@n"), TRUE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_throw(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		vict   *char_data = nil
		tch    *char_data = nil
		obj    *obj_data  = nil
		arg    [2048]byte
		arg2   [1000]byte
		chunk  [2000]byte
		arg3   [1000]byte
		odam   int = 0
		miss   int = TRUE
		perc   int = 0
		prob   int = 0
		perc2  int = 0
		grab   int = FALSE
		damage int64
	)
	half_chop(argument, &arg[0], &chunk[0])
	if arg[0] == 0 {
		send_to_char(ch, libc.CString("Throw what?\r\n"))
		return
	}
	if is_sparring(ch) {
		send_to_char(ch, libc.CString("You can not spar with throw.\r\n"))
		return
	}
	if chunk[0] != 0 {
		two_arguments(&chunk[0], &arg2[0], &arg3[0])
	}
	if (func() *obj_data {
		obj = get_obj_in_list_vis(ch, &arg[0], nil, ch.Carrying)
		return obj
	}()) == nil {
		if (func() *char_data {
			tch = get_char_vis(ch, &arg[0], nil, 1<<0)
			return tch
		}()) == nil {
			send_to_char(ch, libc.CString("You do not have that object or character to throw!\r\n"))
			return
		}
	}
	if (func() *char_data {
		vict = get_char_vis(ch, &arg2[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil && ch.Fighting.In_room == ch.In_room {
			vict = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("Who do you want to target?\r\n"))
			return
		}
	}
	if vict.Hit <= 1 {
		return
	}
	if can_kill(ch, vict, nil, 1) == 0 {
		return
	}
	if handle_defender(vict, ch) != 0 {
		var def *char_data = vict.Defender
		vict = def
	}
	if obj != nil {
		if ch.Throws == -1 {
			ch.Throws = 0
			return
		}
		if ch.Move < ((ch.Max_hit / 200) + obj.Weight) {
			send_to_char(ch, libc.CString("You do not have enough stamina to do it...\r\n"))
			return
		}
		if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("That is broken and useless to throw!\r\n"))
			return
		}
		if obj.Weight+int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity) > max_carry_weight(ch) {
			send_to_char(ch, libc.CString("The gravity has made that too heavy for you to throw!\r\n"))
			return
		} else {
			var (
				penalty    int = 0
				chance     int = axion_dice(0) + axion_dice(0)
				wtype      int = 0
				wlvl       int = 1
				multithrow int = TRUE
			)
			handle_cooldown(ch, 5)
			improve_skill(ch, SKILL_THROW, 0)
			damage = int64(float64((obj.Weight/3)*int64(ch.Aff_abils.Str)*int64(int(ch.Aff_abils.Cha)/3)) + float64(ch.Max_hit)*0.01)
			damage += int64((float64(damage) * 0.01) * float64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity/4))
			if ch.Preference == PREFERENCE_THROWING {
				chance -= int(float64(chance) * 0.25)
			}
			if OBJ_FLAGGED(obj, ITEM_WEAPLVL1) {
				damage += int64(float64(damage) * 0.1)
				wlvl = 1
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL2) {
				damage += int64(float64(damage) * 0.2)
				wlvl = 2
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL3) {
				damage += int64(float64(damage) * 0.3)
				wlvl = 3
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL4) {
				damage += int64(float64(damage) * 0.4)
				wlvl = 4
			} else if OBJ_FLAGGED(obj, ITEM_WEAPLVL5) {
				damage += int64(float64(damage) * 0.5)
				wlvl = 5
			}
			if int(obj.Type_flag) == ITEM_WEAPON {
				if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_PIERCE-TYPE_HIT) {
					wtype = 1
				} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_SLASH-TYPE_HIT) {
					wtype = 2
				} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_CRUSH-TYPE_HIT) {
					wtype = 3
				} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_STAB-TYPE_HIT) {
					wtype = 4
				} else if (obj.Value[VAL_WEAPON_DAMTYPE]) == int(TYPE_BLAST-TYPE_HIT) {
					wtype = 5
					damage = int64(float64((obj.Weight*int64(ch.Aff_abils.Str))*int64(int(ch.Aff_abils.Cha)/3)) + float64(ch.Max_hit)*0.01)
					damage += int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity / 2))
				} else {
					wtype = 6
				}
			}
			if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_STEEL {
				odam = rand_number(5, 30)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_IRON {
				odam = rand_number(18, 50)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_MITHRIL {
				odam = rand_number(5, 15)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_KACHIN {
				odam = rand_number(5, 15)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_STONE {
				odam = rand_number(20, 50)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_DIAMOND {
				odam = rand_number(5, 20)
			} else if (obj.Value[VAL_ALL_MATERIAL]) == MATERIAL_ENERGY {
				if rand_number(1, 2) == 2 {
					odam = 0
				} else {
					odam = rand_number(1, 3)
				}
			} else {
				odam = rand_number(90, 100)
			}
			if !OBJ_FLAGGED(obj, ITEM_THROW) {
				penalty = 15
				multithrow = FALSE
				damage = int64(float64(damage) * 0.45)
			} else {
				odam = rand_number(0, 1)
				damage += int64(float64(ch.Aff_abils.Str) * ((float64(ch.Hit) * 0.00012) + float64(rand_number(1, 20))))
				damage += int64(float64(wlvl) * (float64(damage) * 0.1))
			}
			if wlvl == 5 {
				damage += 25000
			} else if wlvl == 4 {
				damage += 16000
			} else if wlvl == 3 {
				damage += 10000
			} else if wlvl == 2 {
				damage += 5000
			} else if wlvl == 1 {
				damage += 1000
			}
			var hot int = FALSE
			if OBJ_FLAGGED(obj, ITEM_HOT) {
				hot = TRUE
			}
			if wtype > 0 && wtype != 5 && odam > 1 {
				odam = 1
			}
			perc = init_skill(ch, SKILL_THROW)
			perc2 = init_skill(vict, SKILL_DODGE)
			prob = axion_dice(penalty)
			if arg3[0] != 0 {
				if libc.StrCaseCmp(&arg3[0], libc.CString("1")) == 0 || libc.StrCaseCmp(&arg3[0], libc.CString("single")) == 0 {
					multithrow = FALSE
				} else {
					send_to_char(ch, libc.CString("Syntax: throw (obj | character) (target) <-- This will multithrow if able\nSyntax: throw (obj) (target) (1 | single) <-- This will not multi throw)\r\n"))
					return
				}
			}
			if (!IS_NPC(vict) && int(vict.Race) == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && int(vict.Position) != POS_SLEEPING {
				if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
					act(libc.CString("@C$N@c disappears, avoiding your $p before reappearing!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@cYou disappear, avoiding @C$n's@c $p before reappearing!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@c disappears, avoiding @C$n's@c $p before reappearing!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					ch.Combo = -1
					ch.Combhits = 0
					if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
						ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
					}
					var stcost int = int((ch.Max_hit / 200) + obj.Weight)
					vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
					pcost(ch, 0, int64(stcost/2))
					pcost(vict, 0, vict.Max_hit/200)
					obj_from_char(obj)
					obj_to_room(obj, vict.In_room)
					return
				} else {
					act(libc.CString("@C$N@c disappears, trying to avoid your attack but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@cYou zanzoken to avoid the attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c attack but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
			}
			if perc-perc2/10 < prob {
				if OBJ_FLAGGED(obj, ITEM_ICE) && int(vict.Race) == RACE_DEMON {
					act(libc.CString("You throw $p at $N@n, but it melts before touching $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n@n throws $p at $N@n, but it melts before touching $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("$n@n throws $p at you, but it melts before touching you!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
					ch.Move -= (ch.Max_hit / 100) + obj.Weight
					extract_obj(obj)
					return
				}
				if perc2 > 0 {
					act(libc.CString("You throw $p at $N@n, but $E manages to dodge it easily!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n@n throws $p at $N@n, but $E manages to dodge it easily!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("$n@n throws $p at you, but you easily dodge it."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
				} else if perc2 <= 0 {
					act(libc.CString("You throw $p at $N@n, but unfortunatly miss!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("$n@n throws $p at $N@n, but unfortunatly misses!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
					act(libc.CString("$n@n throws $p at you, but thankfully misses you."), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
				}
				ch.Move -= (ch.Max_hit / 100) + obj.Weight
				if !OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) {
					obj.Value[VAL_ALL_HEALTH] -= odam / 2
				}
				ch.Lastattack = -50
				hurt(0, 0, ch, vict, nil, 0, 0)
				obj_from_char(obj)
				obj_to_room(obj, vict.In_room)
				ch.Move -= (ch.Max_hit / 200) + obj.Weight
				if (ch.Equipment[WEAR_WIELD1]) == nil && (ch.Equipment[WEAR_WIELD2]) == nil {
					perc += 20
				}
				if perc+int(ch.Aff_abils.Cha) >= chance+penalty && multithrow == TRUE && vict.Hit > 1 && ch.Throws > 1 {
					do_throw(ch, argument, 0, 0)
					ch.Throws -= 1
				} else if perc+int(ch.Aff_abils.Cha) >= chance+penalty && multithrow == TRUE && vict.Hit > 1 && ch.Throws == 1 {
					do_throw(ch, argument, 0, 0)
					ch.Throws = -1
				} else {
					ch.Throws = 0
				}
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				return
			} else if perc-perc2/10 > prob {
				miss = FALSE
			}
			if !IS_NPC(ch) && PRF_FLAGGED(ch, PRF_ENERGIZE) && float64(ch.Mana) >= float64(ch.Max_mana)*0.02 {
				damage += int64(float64(damage) * (float64(GET_SKILL(ch, SKILL_ENERGIZE)) * 0.0016))
				act(libc.CString("You charge $p with the energy in your fingertips! As it begins to @Yglow a bright hot @Rred@n you throw $p at $N@n full speed, and watch it smash into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n charges $p with the energy in $s fingertips! As it begins to @Yglow a bright hot @Rred@n $e throws $p at $N@n full speed, and watches it smash into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n charges $p with the energy in $s fingertips! As it begins to @Yglow a bright hot @Rred@n $e throws $p at YOU@n full speed, and watches it smash into YOU!!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
				if float64(ch.Max_mana)*0.02 > 0 {
					ch.Mana -= int64(float64(ch.Max_mana) * 0.02)
				} else {
					ch.Mana -= 1
				}
				improve_skill(ch, SKILL_ENERGIZE, 0)
			} else if wtype == 0 {
				act(libc.CString("You throw $p at $N@n full speed, and watch it smash into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n throws $p at $N@n full speed, and watches it smash into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n throws $p at you full speed. You reel as it smashes into your body!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			} else if wtype == 1 || wtype == 2 {
				act(libc.CString("You pull out and throw $p at $N@n full speed, and watch it sink into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n pulls out and throws $p at $N@n full speed, and watches it sink into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n pulls out and throws $p at you full speed. You reel as it sink into your body!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			} else if wtype == 3 {
				act(libc.CString("You swing $p overhead and throw it at $N@n full speed, and watch it slam into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n swings $p overhead and throws it at $N@n full speed, and watches it slam into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n swings $p overhead and throws it at you full speed. You reel as it slam into your body!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			} else if wtype == 4 {
				act(libc.CString("You bring $p over your shoulder and throw it at $N@n full speed, and watch it sink into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n brings $p over $s shoulder and throws it at $N@n full speed, and watches it sink into $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n brings $p over $s shoulder and throws $p at you full speed. You reel as it sink into your body!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			} else if wtype == 5 {
				act(libc.CString("You pull out and throw $p at $N@n full speed, and watch it hit $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n pulls out and throws $p at $N@n full speed, and watches it hit $M!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n pulls out and throws $p at you full speed. You reel as it hits your body!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			}
			if !OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) {
				obj.Value[VAL_ALL_HEALTH] -= odam
			}
			ch.Lastattack = -50
			if ((obj.Value[VAL_ALL_HEALTH])-odam) <= 0 && !OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) {
				act(libc.CString("You smile as $p breaks on $N's@n face!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("$n@n smiles as $p breaks on $N's@n face!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("$n@n smiles as $p breaks on your face!"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
				obj.Extra_flags[int(ITEM_BROKEN/32)] = bitvector_t(int32(int(obj.Extra_flags[int(ITEM_BROKEN/32)]) ^ 1<<(int(ITEM_BROKEN%32))))
			} else if int(ch.Aff_abils.Dex) >= axion_dice(0) {
				if int(vict.Race) == RACE_ANDROID || int(vict.Race) == RACE_WARHOST {
					act(libc.CString("@RSome pieces of metal are sent flying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RSome pieces of metal are sent flying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@RSome pieces of metal are sent flying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else if int(vict.Race) == RACE_MAJIN {
					act(libc.CString("@RA wide hole is left in $S gooey flesh!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RA wide hole is left is your gooey flesh!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@RA wide hole is left in $N@R's gooey flesh@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else {
					act(libc.CString("@RBlood flies out from the impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RBlood flies out from the impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@RBlood flies out from the impact!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				}
				if OBJ_FLAGGED(obj, ITEM_ICE) {
					if int(vict.Race) != RACE_ANDROID && int(vict.Race) != RACE_ICER {
						vict.Move -= int64((float64(vict.Max_move) * 0.005) + float64(obj.Weight))
						if vict.Move < 0 {
							vict.Move = 0
						}
						act(libc.CString("@mYou lose some stamina to the @ccold@m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						act(libc.CString("@C$N@m loses some stamina to the @ccold@m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@C$N@m loses some stamina to the @ccold@m!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					}
				}
				damage *= int64(calc_critical(ch, 0))
			}
			if hot == TRUE {
				if int(vict.Race) != RACE_DEMON && (vict.Bonuses[BONUS_FIREPROOF]) == 0 {
					act(libc.CString("@R$N@R is burned by it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYou are burned by it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@R is burned by it!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
					damage += int64(float64(damage) * 0.4)
				}
			}
			if ch.Preference == PREFERENCE_KI {
				damage -= int64(float64(damage) * 0.2)
			}
			if GET_OBJ_VNUM(obj) == 5899 || GET_OBJ_VNUM(obj) == 5898 {
				damage *= int64(0.35)
			}
			hurt(0, 0, ch, vict, nil, damage, 0)
			obj_from_char(obj)
			obj_to_room(obj, vict.In_room)
			ch.Move -= (ch.Max_hit / 200) + obj.Weight
			if (ch.Equipment[WEAR_WIELD1]) == nil && (ch.Equipment[WEAR_WIELD2]) == nil {
				perc += 12
			}
			if perc+int(ch.Aff_abils.Cha) >= chance+penalty && multithrow == TRUE && vict.Hit > 1 && ch.Throws > 1 {
				do_throw(ch, argument, 0, 0)
				ch.Throws -= 1
			} else if perc+int(ch.Aff_abils.Cha) >= chance+penalty && multithrow == TRUE && vict.Hit > 1 && ch.Throws == 1 {
				do_throw(ch, argument, 0, 0)
				ch.Throws = -1
			} else {
				ch.Throws = 0
			}
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			return
		}
	}
	if tch != nil {
		if tch == vict {
			send_to_char(ch, libc.CString("You can't throw someone at theirself.\r\n"))
			return
		}
		if can_kill(ch, tch, nil, 0) == 0 {
			send_to_char(ch, libc.CString("The one you are throwing can't be harmed.\r\n"))
			return
		}
		if GET_SPEEDI(tch) < GET_SPEEDI(ch) && rand_number(1, 106) < GET_SKILL(ch, SKILL_THROW) {
			grab = TRUE
		}
		if ch.Move < ((ch.Max_hit / 100) + int64(GET_PC_WEIGHT(tch))) {
			send_to_char(ch, libc.CString("You do not have enough stamina to do it...\r\n"))
			return
		}
		if GET_PC_WEIGHT(tch)+(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity > int(max_carry_weight(ch)) {
			send_to_char(ch, libc.CString("The gravity has made them too heavy for you to throw!\r\n"))
			return
		}
		if grab == FALSE {
			act(libc.CString("@WYou try to grab @C$N@W and throw them, but they manage to dodge your attempt!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
			act(libc.CString("@C$n@W tries to @RGRAB@W you and @RTHROW@W you, but you manage to dodge the attempt!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
			act(libc.CString("@C$n@W tries to @RGRAB@W @c$N@W and @RTHROW@W $M, but $E manages to dodge the attempt!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
			hurt(0, 0, ch, tch, nil, 0, 0)
			handle_cooldown(ch, 5)
			ch.Move -= (ch.Max_hit / 200) + int64(GET_PC_WEIGHT(tch))
			return
		} else {
			handle_cooldown(ch, 5)
			improve_skill(ch, SKILL_THROW, 0)
			damage = int64(((GET_PC_WEIGHT(tch) * int(ch.Aff_abils.Str)) * (int(ch.Aff_abils.Cha) / 3)) + int(ch.Max_hit/100))
			damage += int64((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity * ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity / 2))
			perc = init_skill(ch, SKILL_THROW)
			perc2 = init_skill(vict, SKILL_DODGE)
			prob = rand_number(1, 106)
			if perc-perc2/10 < prob {
				if perc2 > 0 {
					act(libc.CString("@WYou grab @C$N@W and spinning around quickly you throw $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
					act(libc.CString("@C$n@W grabs YOU and spinning around quickly $e throws you!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
					act(libc.CString("@C$n@W grabs @c$N@W and spinning around quickly $e throws $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
					act(libc.CString("@WThrown through the air, YOU fly at @c$N@W, but $E manages to dodge and you manage recover your bearings a moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WThrown through the air, @C$n@W flies at YOU, but you manage to dodge and @C$n@W recovers $s bearings a moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@WThrown through the air, @C$n@W flies at @c$N@W, but $E manages to dodge and @C$n@W recovers $s bearingsa moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else if perc2 <= 0 {
					act(libc.CString("@WYou grab @C$N@W and spinning around quickly you throw $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
					act(libc.CString("@C$n@W grabs YOU and spinning around quickly $e throws you!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
					act(libc.CString("@C$n@W grabs @c$N@W and spinning around quickly $e throws $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
					act(libc.CString("@WThrown through the air, YOU fly at @c$N@W, but the throw is a miss! You manage recover your bearings a moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WThrown through the air, @C$n@W flies at YOU, but the throw is a miss! @C$n@W recovers $s bearings a moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@WThrown through the air, @C$n@W flies at @c$N@W, but the throw is a miss! @C$n@W recovers $s bearingsa moment later!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				}
				ch.Move -= (ch.Max_hit / 100) + int64(GET_PC_WEIGHT(tch))
				act(libc.CString("@W--@R$N@W--@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W--@RYOU@W--@n"), TRUE, vict, nil, nil, TO_CHAR)
				hurt(0, 0, ch, vict, nil, 0, 0)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, vict, nil, unsafe.Pointer(tch), TO_CHAR)
				act(libc.CString("@W--@RYOU@W--@n"), TRUE, tch, nil, nil, TO_CHAR)
				hurt(0, 0, ch, tch, nil, 0, 0)
				return
			} else if perc-perc2/10 >= prob {
				miss = FALSE
			}
			if miss == FALSE {
				act(libc.CString("@WYou grab @C$N@W and spinning around quickly you throw $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				act(libc.CString("@C$n@W grabs YOU and spinning around quickly $e throws you!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
				act(libc.CString("@C$n@W grabs @c$N@W and spinning around quickly $e throws $M!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
				act(libc.CString("@WThrown through the air, YOU fly at @c$N@W and smash into $M!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@WThrown through the air, @C$n@W flies at YOU and smashes into YOU!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@WThrown through the air, @C$n@W flies at @c$N@W and smashes into $M!@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, tch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@W--@RYOU@W--@n"), TRUE, vict, nil, nil, TO_CHAR)
				hurt(0, 0, ch, vict, nil, damage, 0)
				act(libc.CString("@W--@R$N@W--@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				if vict != nil {
					act(libc.CString("@W--@R$N@W--@n"), TRUE, vict, nil, unsafe.Pointer(tch), TO_CHAR)
				}
				act(libc.CString("@W--@RYOU@W--@n"), TRUE, tch, nil, nil, TO_CHAR)
				hurt(0, 0, ch, tch, nil, damage, 0)
			}
			ch.Move -= (ch.Max_hit / 200) + int64(GET_PC_WEIGHT(tch))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
	}
	if obj == nil && tch == nil {
		send_to_imm(libc.CString("ERROR: Throw resolved without character or object."))
		return
	}
}
func do_selfd(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		tch    *char_data = nil
		next_v *char_data = nil
		dmg    int64      = 0
	)
	if IS_NPC(ch) {
		return
	}
	if IN_ARENA(ch) {
		send_to_char(ch, libc.CString("You can not use self destruct in the arena.\r\n"))
		return
	}
	if AFF_FLAGGED(ch, AFF_SPIRIT) {
		send_to_char(ch, libc.CString("You are already dead!\r\n"))
		return
	}
	if GET_LEVEL(ch) < 9 {
		send_to_char(ch, libc.CString("You can't self destruct while protected by the newbie shield!\r\n"))
		return
	}
	if int(libc.BoolToInt(ch.Con_sdcooldown == 0)) <= 0 {
		send_to_char(ch, libc.CString("Your body is still recovering from the last self destruct!\r\n"))
		return
	}
	if GET_SKILL(ch, SKILL_SELFD) == 0 {
		var num int = rand_number(10, 20)
		for {
			ch.Skills[SKILL_SELFD] = int8(num)
			if true {
				break
			}
		}
	}
	if !PLR_FLAGGED(ch, PLR_SELFD) {
		act(libc.CString("@RYour body starts to glow @wwhite@R and flash. The flashes start out slowly but steadilly increase in speed. Your aura begins to burn around your body at the same time in a violent fashion!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@R$n's body starts to glow @wwhite@R and flash. The flashes start out slowly but steadilly increase in speed. $n's aura begins to burn around $s body at the same time in a violent fashion!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_SELFD/32)] |= bitvector_t(int32(1 << (int(PLR_SELFD % 32))))
		return
	} else if !PLR_FLAGGED(ch, PLR_SELFD2) {
		act(libc.CString("@wYour body slowly stops flashing. Steam rises from your skin as you slowly let off the energy you built up in a safe manner.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@w$n's body slowly stops flashing. Steam rises from $s skin as $e slowly lets off the energy $e built up in a safe manner.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_SELFD/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD % 32)))))
		return
	} else if ch.Grappling != nil && can_kill(ch, ch.Grappling, nil, 3) == 0 {
		act(libc.CString("@wYour body slowly stops flashing. Steam rises from your skin as you slowly let off the energy you built up in a safe manner.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@w$n's body slowly stops flashing. Steam rises from $s skin as $e slowly lets off the energy $e built up in a safe manner.@n"), TRUE, ch, nil, nil, TO_ROOM)
		send_to_char(ch, libc.CString("You can't kill them, the immortals won't allow it!\r\n"))
		ch.Act[int(PLR_SELFD/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD % 32)))))
		return
	} else if ch.Grappling != nil {
		tch = ch.Grappling
		dmg += ch.Charge
		ch.Charge = 0
		dmg += int64(float64(ch.Basepl) * 0.6)
		dmg += ch.Basest
		ch.Hit = 1
		ch.Suppressed = 0
		ch.Suppression = 0
		act(libc.CString("@RYou EXPLODE! The explosion concentrates on @r$N@R, engulfing $M in a sphere of deadly energy!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
		act(libc.CString("@R$n EXPLODES! The explosion concentrates on YOU, engulfing your body in a sphere of deadly energy!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
		act(libc.CString("@R$n EXPLODES! The explosion concentrates on @r$N@R, engulfing $M in a sphere of deadly energy!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
		hurt(0, 0, ch, tch, nil, dmg, 1)
		ch.Act[int(PLR_SELFD/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD % 32)))))
		ch.Act[int(PLR_SELFD2/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD2 % 32)))))
		if PLR_FLAGGED(ch, PLR_IMMORTAL) {
			ch.Con_sdcooldown = 600
		}
		if (int(ch.Race) == RACE_MAJIN || int(ch.Race) == RACE_BIO) && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.5 {
			ch.Lifeforce = -1
			ch.Act[int(PLR_GOOP/32)] |= bitvector_t(int32(1 << (int(PLR_GOOP % 32))))
			ch.Gooptime = 70
		} else {
			die(ch, nil)
		}
		var num int = rand_number(10, 20) + GET_SKILL(ch, SKILL_SELFD)
		if GET_SKILL(ch, SKILL_SELFD)+num <= 100 {
			for {
				ch.Skills[SKILL_SELFD] = int8(num)
				if true {
					break
				}
			}
		} else {
			for {
				ch.Skills[SKILL_SELFD] = 100
				if true {
					break
				}
			}
		}
		return
	} else {
		dmg += ch.Charge
		ch.Charge = 0
		dmg += int64(float64(ch.Basepl) * 0.6)
		dmg += ch.Basest
		dmg *= int64(1.5)
		ch.Hit = 1
		ch.Suppressed = 0
		ch.Suppression = 0
		act(libc.CString("@RYou EXPLODE! The explosion expands outward burning up all surroundings for a large distance. The explosion takes on the shape of a large energy dome with you at its center!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@R$n EXPLODES! The explosion expands outward burning up all surroundings for a large distance. The explosion takes on the shape of a large energy dome with $n at its center!@n"), TRUE, ch, nil, nil, TO_ROOM)
		for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = next_v {
			next_v = tch.Next_in_room
			if tch == ch {
				continue
			}
			if can_kill(ch, tch, nil, 3) == 0 {
				continue
			}
			if MOB_FLAGGED(tch, MOB_NOKILL) {
				continue
			} else {
				act(libc.CString("@r$N@R is caught in the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				act(libc.CString("@RYou are caught in the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
				act(libc.CString("@r$N@R is caught in the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
				hurt(0, 0, ch, tch, nil, dmg, 1)
			}
		}
		if PLR_FLAGGED(ch, PLR_IMMORTAL) {
			ch.Con_sdcooldown = 600
		}
		die(ch, nil)
		ch.Act[int(PLR_SELFD/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD % 32)))))
		ch.Act[int(PLR_SELFD2/32)] &= bitvector_t(int32(^(1 << (int(PLR_SELFD2 % 32)))))
		var num int = rand_number(10, 20) + GET_SKILL(ch, SKILL_SELFD)
		if GET_SKILL(ch, SKILL_SELFD)+num <= 100 {
			for {
				ch.Skills[SKILL_SELFD] = int8(num)
				if true {
					break
				}
			}
		} else {
			for {
				ch.Skills[SKILL_SELFD] = 100
				if true {
					break
				}
			}
		}
		return
	}
}
func do_razor(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		attperc float64 = 0.14
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
	if check_skill(ch, SKILL_WRAZOR) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_WRAZOR)
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
		if int(vict.Race) == RACE_ANDROID {
			send_to_char(ch, libc.CString("There is not a necessary amount of water in cybernetic creatures.\r\n"))
			return
		}
		if can_kill(ch, vict, nil, 1) == 0 {
			return
		}
		if handle_defender(vict, ch) != 0 {
			var def *char_data = vict.Defender
			vict = def
		}
		improve_skill(ch, SKILL_WRAZOR, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Water Razor before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Water Razor before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Water Razor before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Water Razor but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Water Razor but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Water Razor but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Water Razor, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Water Razor, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Water Razor, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Water Razor misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Water Razor at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Water Razor at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Water Razor misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Water Razor at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Water Razor at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			var reduction int64 = vict.Hit
			dmg = damtype(ch, 47, skill, attperc)
			if AFF_FLAGGED(ch, AFF_SANCTUARY) {
				if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 100 {
					ch.Barrier += int64(float64(dmg) * 0.1)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 60 {
					ch.Barrier += int64(float64(dmg) * 0.05)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 40 {
					ch.Barrier += int64(float64(dmg) * 0.02)
				}
				if ch.Barrier > ch.Max_mana {
					ch.Barrier = ch.Max_mana
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise your hand toward @c$N@W with the palm open as if you are ready to grab the air between the two of you. You then focus your ki into $S body and close your hand in an instant! The @Bwater@W in $S body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand toward YOU with the palm open as if $e is ready to grab the air between the two of you. $e then focus $s ki into YOUR body and close $s hand in an instant! The @Bwater@W in your body instantly takes the shape of millions of microscopic @Dblades@W that cut up YOUR insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand toward @c$N@W with the palm open as if $e is ready to grab the air between the two of them. $e then focus $s ki into @c$N's@W body and close $s hand in an instant! The @Bwater@W in @c$N's@W body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
			case 2:
				act(libc.CString("@WYou raise your hand toward @c$N@W with the palm open as if you are ready to grab the air between the two of you. You then focus your ki into $S body and close your hand in an instant! The @Bwater@W in $S body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides! Blood sprays out in a mist from every pore of $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand toward YOU with the palm open as if $e is ready to grab the air between the two of you. $e then focus $s ki into YOUR body and close $s hand in an instant! The @Bwater@W in your body instantly takes the shape of millions of microscopic @Dblades@W that cut up YOUR insides! Blood sprays out of your pores into a fine mist!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand toward @c$N@W with the palm open as if $e is ready to grab the air between the two of them. $e then focus $s ki into @c$N's@W body and close $s hand in an instant! The @Bwater@W in @c$N's@W body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides! Blood sprays out of $S pores into a fine mist!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
			case 3:
				act(libc.CString("@WYou raise your hand toward @c$N@W with the palm open as if you are ready to grab the air between the two of you. You then focus your ki into $S body and close your hand in an instant! The @Bwater@W in $S body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides! Blood sprays out in a mist from every pore of $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand toward YOU with the palm open as if $e is ready to grab the air between the two of you. $e then focus $s ki into YOUR body and close $s hand in an instant! The @Bwater@W in your body instantly takes the shape of millions of microscopic @Dblades@W that cut up YOUR insides! Blood sprays out of your pores into a fine mist!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand toward @c$N@W with the palm open as if $e is ready to grab the air between the two of them. $e then focus $s ki into @c$N's@W body and close $s hand in an instant! The @Bwater@W in @c$N's@W body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides! Blood sprays out of $S pores into a fine mist!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
			case 4:
				act(libc.CString("@WYou raise your hand toward @c$N@W with the palm open as if you are ready to grab the air between the two of you. You then focus your ki into $S body and close your hand in an instant! The @Bwater@W in $S body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand toward YOU with the palm open as if $e is ready to grab the air between the two of you. $e then focus $s ki into YOUR body and close $s hand in an instant! The @Bwater@W in your body instantly takes the shape of millions of microscopic @Dblades@W that cut up YOUR insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand toward @c$N@W with the palm open as if $e is ready to grab the air between the two of them. $e then focus $s ki into @c$N's@W body and close $s hand in an instant! The @Bwater@W in @c$N's@W body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@WYou raise your hand toward @c$N@W with the palm open as if you are ready to grab the air between the two of you. You then focus your ki into $S body and close your hand in an instant! The @Bwater@W in $S body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s hand toward YOU with the palm open as if $e is ready to grab the air between the two of you. $e then focus $s ki into YOUR body and close $s hand in an instant! The @Bwater@W in your body instantly takes the shape of millions of microscopic @Dblades@W that cut up YOUR insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s hand toward @c$N@W with the palm open as if $e is ready to grab the air between the two of them. $e then focus $s ki into @c$N's@W body and close $s hand in an instant! The @Bwater@W in @c$N's@W body instantly takes the shape of millions of microscopic @Dblades@W that cut up $S insides!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				hurt(0, 0, ch, vict, nil, dmg, 1)
			}
			pcost(ch, attperc, 0)
			if vict != nil {
				reduction = reduction - vict.Hit
				if !IS_NPC(vict) && !AFF_FLAGGED(vict, AFF_SPIRIT) {
					if vict.Mana > reduction {
						vict.Mana -= reduction
					} else {
						vict.Mana = 0
					}
					if vict.Move > reduction {
						vict.Move -= reduction
					} else {
						vict.Move = 0
					}
				} else if IS_NPC(vict) && vict.Hit > 0 {
					if vict.Mana > reduction {
						vict.Mana -= reduction
					} else {
						vict.Mana = 0
					}
					if vict.Move > reduction {
						vict.Move -= reduction
					} else {
						vict.Move = 0
					}
				}
			}
			return
		}
	} else if obj != nil {
		send_to_char(ch, libc.CString("You can not hurt an inanimate object with this technique.\r\n"))
		return
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_spike(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		prob    int
		perc    int
		avo     int
		index   int
		pry     int = 2
		dge     int = 2
		blk     int = 2
		skill   int
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
		dmg     int64
		attperc float64 = 0.14
		minimum float64 = 0.05
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_WSPIKE) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if int(ch.Skillperfs[SKILL_WSPIKE]) == 1 {
		attperc += 0.05
	} else if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
		minimum -= 0.05
		if minimum <= 0.0 {
			minimum = 0.01
		}
	}
	if arg2[0] != 0 {
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_WSPIKE)
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
		improve_skill(ch, SKILL_WSPIKE, 0)
		index = check_def(vict)
		prob = roll_accuracy(ch, skill, TRUE != 0)
		if int(ch.Skillperfs[SKILL_WSPIKE]) == 2 {
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
				act(libc.CString("@C$N@c disappears, avoiding your attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c attack before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
				} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
				pcost(vict, 0, vict.Max_hit/200)
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
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
				if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W moves quickly and blocks your water spikes!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W water spikes!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W water spikes!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					dmg = damtype(ch, 10, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your water spikes, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W water spikes, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W water spikes, letting them slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 43, skill, SKILL_WSPIKE)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your water spikes miss, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fire water spikes at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fire water spikes at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
						pcost(ch, attperc-0.05, 0)
					} else {
						pcost(ch, attperc, 0)
					}
					if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
					} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
						ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
					}
					if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
						WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
					}
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your water spikes miss, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires water spikes at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires water spikes at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
					pcost(ch, attperc-0.05, 0)
				} else {
					pcost(ch, attperc, 0)
				}
				if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
				} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
				} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
					ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
				}
			}
			if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
				WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 43, skill, attperc)
			if AFF_FLAGGED(ch, AFF_SANCTUARY) {
				if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 100 {
					ch.Barrier += int64(float64(dmg) * 0.1)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 60 {
					ch.Barrier += int64(float64(dmg) * 0.05)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 40 {
					ch.Barrier += int64(float64(dmg) * 0.02)
				}
				if ch.Barrier > ch.Max_mana {
					ch.Barrier = ch.Max_mana
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@CYou slam both your hands together in front of you, palms flat. As you pull them apart ki flows from your palms and forms into a giant ball of water. You raise your hand above your head with the ball of water. You command the water to form several sharp spikes which freeze upon forming. The spikes then launch at @R$N@C and slam into $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @RYOU@C and slam into YOUR chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @R$N@C and slam into $s chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@CYou slam both your hands together in front of you, palms flat. As you pull them apart ki flows from your palms and forms into a giant ball of water. You raise your hand above your head with the ball of water. You command the water to form several sharp spikes which freeze upon forming. The spikes then launch at @R$N@C and slam into $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @RYOU@C and slam into YOUR head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @R$N@C and slam into $s head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
				if !AFF_FLAGGED(vict, AFF_KNOCKED) && (rand_number(1, 3) >= 2 && vict.Hit > ch.Hit/5 && !AFF_FLAGGED(vict, AFF_SANCTUARY)) {
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
			case 3:
				act(libc.CString("@CYou slam both your hands together in front of you, palms flat. As you pull them apart ki flows from your palms and forms into a giant ball of water. You raise your hand above your head with the ball of water. You command the water to form several sharp spikes which freeze upon forming. The spikes then launch at @R$N@C and slam into $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @RYOU@C and slam into YOUR body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @R$N@C and slam into $s body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@CYou slam both your hands together in front of you, palms flat. As you pull them apart ki flows from your palms and forms into a giant ball of water. You raise your hand above your head with the ball of water. You command the water to form several sharp spikes which freeze upon forming. The spikes then launch at @R$N@C and slam into $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @RYOU@C and slam into YOUR arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @R$N@C and slam into $s arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dam_eq_loc(vict, 1)
				hurt(0, 190, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@CYou slam both your hands together in front of you, palms flat. As you pull them apart ki flows from your palms and forms into a giant ball of water. You raise your hand above your head with the ball of water. You command the water to form several sharp spikes which freeze upon forming. The spikes then launch at @R$N@C and slam into $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @RYOU@C and slam into YOUR leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C slams both $s hands together in front of $mself, palms flat. As $e pulls them apart ki flows from the palms and forms into a giant ball of water. Then $e raises $s hand above $s head with the ball of water. The water forms several sharp spikes which freeze instantly as they take shape. The spikes then launch at @R$N@C and slam into $s leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dam_eq_loc(vict, 2)
				if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
					WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
				}
				hurt(1, 190, ch, vict, nil, dmg, 1)
			}
			if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 && attperc > minimum {
				pcost(ch, attperc-0.05, 0)
			} else {
				pcost(ch, attperc, 0)
			}
			if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
			} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
				ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
			} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
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
		dmg = damtype(ch, 43, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire water spikes at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires water spikes at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		if int(ch.Skillperfs[SKILL_WSPIKE]) == 3 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		}
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
		if GET_SKILL(ch, SKILL_WSPIKE) >= 100 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.3)
		} else if GET_SKILL(ch, SKILL_WSPIKE) >= 60 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.1)
		} else if GET_SKILL(ch, SKILL_WSPIKE) >= 40 {
			ch.Mana += int64((float64(ch.Max_mana) * attperc) * 0.05)
		}
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_koteiru(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
		attperc float64 = 0.3
		minimum float64 = 0.1
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_KOTEIRU) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_KOTEIRU)
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
		improve_skill(ch, SKILL_KOTEIRU, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Koteiru Bakuha before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Koteiru Bakuha before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Koteiru Bakuha before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Koteiru Bakuha but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Koteiru Bakuha but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Koteiru Bakuha but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Koteiru Bakuha, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Koteiru Bakuha, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Koteiru Bakuha, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Koteiru Bakuha misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Koteiru Bakuha at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Koteiru Bakuha at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Koteiru Bakuha misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Koteiru Bakuha at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Koteiru Bakuha at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 48, skill, attperc)
			if AFF_FLAGGED(ch, AFF_SANCTUARY) {
				if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 100 {
					ch.Barrier += int64(float64(dmg) * 0.1)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 60 {
					ch.Barrier += int64(float64(dmg) * 0.05)
				} else if GET_SKILL(ch, SKILL_AQUA_BARRIER) >= 40 {
					ch.Barrier += int64(float64(dmg) * 0.02)
				}
				if ch.Barrier > ch.Max_mana {
					ch.Barrier = ch.Max_mana
				}
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@CYou hold your hands outstretched in front of your body and with your ki begin to create turbulant waters that hover around your body. You begin a sort of dance as you control the water more and more until you have created a huge floating wave. In an instant you swing the wave at @W$N@C! As the wave slams into $S chest it freezes solid around $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at YOU! As the wave slams into YOUR chest it freezes solid around YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at @W$N@C! As the wave slams into $S chest it freezes solid around @W$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@CYou hold your hands outstretched in front of your body and with your ki begin to create turbulant waters that hover around your body. You begin a sort of dance as you control the water more and more until you have created a huge floating wave. In an instant you swing the wave at @W$N@C! As the wave slams into $S head it freezes solid around $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at YOU! As the wave slams into YOUR head it freezes solid around YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at @W$N@C! As the wave slams into $S head it freezes solid around @W$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@CYou hold your hands outstretched in front of your body and with your ki begin to create turbulant waters that hover around your body. You begin a sort of dance as you control the water more and more until you have created a huge floating wave. In an instant you swing the wave at @W$N@C! As the wave slams into $S gut it freezes solid around $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at YOU! As the wave slams into YOUR gut it freezes solid around YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at @W$N@C! As the wave slams into $S gut it freezes solid around @W$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@CYou hold your hands outstretched in front of your body and with your ki begin to create turbulant waters that hover around your body. You begin a sort of dance as you control the water more and more until you have created a huge floating wave. In an instant you swing the wave at @W$N@C! As the wave slams into $S arm it freezes solid around $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at YOU! As the wave slams into YOUR arm it freezes solid around YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at @W$N@C! As the wave slams into $S arm it freezes solid around @W$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 1)
				hurt(0, 190, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@CYou hold your hands outstretched in front of your body and with your ki begin to create turbulant waters that hover around your body. You begin a sort of dance as you control the water more and more until you have created a huge floating wave. In an instant you swing the wave at @W$N@C! As the wave slams into $S leg it freezes solid around $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at YOU! As the wave slams into YOUR leg it freezes solid around YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@C holds $s hands outstretched in front of $s body and with $s ki begins to create turbulant waters that hover around $s body. @c$n@C begins a sort of dance as $e controls the water more and more until $e has created a huge floating wave. In an instant $e swings the wave at @W$N@C! As the wave slams into $S leg it freezes solid around @W$N@C!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 2)
				hurt(1, 190, ch, vict, nil, dmg, 1)
			}
			pcost(ch, attperc, 0)
			if vict != nil {
				if vict.Hit > 1 {
					if rand_number(1, 4) == 1 && !AFF_FLAGGED(vict, AFF_FROZEN) && int(vict.Race) != RACE_DEMON {
						act(libc.CString("@CYour body completely freezes!@n"), TRUE, vict, nil, nil, TO_CHAR)
						act(libc.CString("@c$n's@C body completely freezes!@n"), TRUE, vict, nil, nil, TO_ROOM)
						vict.Affected_by[int(AFF_FROZEN/32)] |= 1 << (int(AFF_FROZEN % 32))
					}
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
		dmg = damtype(ch, 48, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire Koteiru Bakuha at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires Koteiru Bakuha at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_hspiral(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
		attperc float64 = 0.3
		minimum float64 = 0.15
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_HSPIRAL) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_HSPIRAL)
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
		improve_skill(ch, SKILL_HSPIRAL, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Hell Spiral before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Hell Spiral before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Hell Spiral before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Hell Spiral but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Hell Spiral but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Hell Spiral but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Hell Spiral, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Hell Spiral, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Hell Spiral, letting it fly harmlessly by!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Hell Spiral misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Hell Spiral at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Hell Spiral at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Hell Spiral misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Hell Spiral at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Hell Spiral at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 49, skill, attperc)
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou hold out your hand, palm upward, while looking toward @C$N@W and grinning as @Rred@W energy begins to pool in the center of your palm. As the large orb of energy swells to about three feet in diameter you move your hand away! You quickly punch the ball of energy and send it flying into @C$N's@W chest where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward YOU and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into YOUR chest where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward @c$N@W and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into @c$N's@W chest where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou hold out your hand, palm upward, while looking toward @C$N@W and grinning as @Rred@W energy begins to pool in the center of your palm. As the large orb of energy swells to about three feet in diameter you move your hand away! You quickly punch the ball of energy and send it flying into @C$N's@W head where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward YOU and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into YOUR head where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward @c$N@W and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into @c$N's@W head where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou hold out your hand, palm upward, while looking toward @C$N@W and grinning as @Rred@W energy begins to pool in the center of your palm. As the large orb of energy swells to about three feet in diameter you move your hand away! You quickly punch the ball of energy and send it flying into @C$N's@W body where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward YOU and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into YOUR body where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward @c$N@W and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into @c$N's@W body where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou hold out your hand, palm upward, while looking toward @C$N@W and grinning as @Rred@W energy begins to pool in the center of your palm. As the large orb of energy swells to about three feet in diameter you move your hand away! You quickly punch the ball of energy and send it flying into @C$N's@W arm where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward YOU and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into YOUR arm where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward @c$N@W and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into @c$N's@W arm where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 1)
				hurt(0, 120, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@WYou hold out your hand, palm upward, while looking toward @C$N@W and grinning as @Rred@W energy begins to pool in the center of your palm. As the large orb of energy swells to about three feet in diameter you move your hand away! You quickly punch the ball of energy and send it flying into @C$N's@W leg where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward YOU and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into YOUR leg where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n @Wholds out $s hand, palm upward, while looking toward @c$N@W and grinning as @Rred@W energy begins to pool in the center of $s palm. As the large orb of energy swells to about three feet in diameter $e moves $s hand away! Then $e quickly punches the ball of energy and sends it flying into @c$N's@W leg where it @re@Rx@Dp@rl@Ro@Dd@re@Rs@W in a flash of light!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 2)
				hurt(1, 120, ch, vict, nil, dmg, 1)
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
		dmg = damtype(ch, 49, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire Hell Spiral at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires Hell Spiral at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
func do_spiral(ch *char_data, argument *byte, cmd int, subcmd int) {
	var (
		skill int
		vict  *char_data
		arg   [2048]byte
	)
	one_argument(argument, &arg[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_SPIRAL) == 0 {
		return
	}
	if limb_ok(ch, 0) == 0 {
		return
	}
	if arg[0] == 0 && ch.Fighting == nil {
		send_to_char(ch, libc.CString("Direct it at who?\r\n"))
		return
	}
	if check_points(ch, int64(float64(ch.Max_mana)*0.5), 0) == 0 {
		return
	}
	skill = init_skill(ch, SKILL_SPIRAL)
	if (func() *char_data {
		vict = get_char_vis(ch, &arg[0], nil, 1<<0)
		return vict
	}()) == nil {
		if ch.Fighting != nil {
			vict = ch.Fighting
		} else {
			send_to_char(ch, libc.CString("Nothing around here by that name."))
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
	ch.Act[int(PLR_SPIRAL/32)] |= bitvector_t(int32(1 << (int(PLR_SPIRAL % 32))))
	improve_skill(ch, SKILL_SPIRAL, 0)
	act(libc.CString("@mFlying to a spot above your intended target you begin to move so fast all that can be seen of you are trails of color. You focus your movements into a vortex and prepare to attack!@n"), TRUE, ch, nil, nil, TO_CHAR)
	act(libc.CString("@w$n@m flies to a spot above and begins to move so fast all that can be seen of $m are trails of color. Suddenly $e focuses $s movements into a spinning vortex and you lose track of $s movements entirely!@n"), TRUE, ch, nil, nil, TO_ROOM)
	handle_spiral(ch, vict, skill, TRUE)
	handle_cooldown(ch, 8)
}
func do_breaker(ch *char_data, argument *byte, cmd int, subcmd int) {
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
		vict    *char_data
		obj     *obj_data
		arg     [2048]byte
		arg2    [2048]byte
		attperc float64 = 0.14
		minimum float64 = 0.05
	)
	two_arguments(argument, &arg[0], &arg2[0])
	if can_grav(ch) == 0 {
		return
	}
	if check_skill(ch, SKILL_BREAKER) == 0 {
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
		var adjust float64 = float64(libc.Atoi(libc.GoString(&arg2[0])))
		adjust *= 0.01
		if adjust <= 0 || adjust > 1 {
			send_to_char(ch, libc.CString("If you are going to supply a percentage of your charge to use then use an acceptable number (1-100)\r\n"))
			return
		} else if adjust <= attperc && adjust >= minimum {
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
	skill = init_skill(ch, SKILL_BREAKER)
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
		improve_skill(ch, SKILL_BREAKER, 0)
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
				act(libc.CString("@C$N@c disappears, avoiding your Star Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Star Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Star Breaker before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, attperc, 0)
				pcost(vict, 0, vict.Max_hit/200)
				return
			} else {
				act(libc.CString("@C$N@c disappears, trying to avoid your Star Breaker but your zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou zanzoken to avoid the Star Breaker but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, trying to avoid @C$n's@c Star Breaker but @C$n's@c zanzoken is faster!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
			}
		}
		if prob < perc-20 {
			if vict.Move > 0 {
				if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W moves quickly and blocks your Star Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Star Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Star Breaker!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					dmg = damtype(ch, 10, skill, attperc)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Star Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Star Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Star Breaker, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 46, skill, SKILL_BREAKER)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Star Breaker misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Star Breaker at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Star Breaker at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, attperc, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Star Breaker misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Star Breaker at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Star Breaker at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, attperc, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			dmg = damtype(ch, 46, skill, attperc)
			var theft int64 = 0
			if GET_LEVEL(ch)-30 > GET_LEVEL(vict) {
				theft = 1
				vict.Exp -= theft
			} else if GET_LEVEL(ch)-20 > GET_LEVEL(vict) {
				theft = vict.Exp / 1000
				vict.Exp -= theft
			} else if GET_LEVEL(ch)-10 > GET_LEVEL(vict) {
				theft = vict.Exp / 100
				vict.Exp -= theft
			} else if GET_LEVEL(ch) >= GET_LEVEL(vict) {
				theft = vict.Exp / 50
				vict.Exp -= theft
			} else if GET_LEVEL(ch)+10 >= GET_LEVEL(vict) {
				theft = vict.Exp / 500
				vict.Exp -= theft
			} else if GET_LEVEL(ch)+20 >= GET_LEVEL(vict) {
				theft = vict.Exp / 1000
				vict.Exp -= theft
			} else {
				theft = vict.Exp / 2000
				vict.Exp -= theft
			}
			var hitspot int = 1
			hitspot = roll_hitloc(ch, vict, skill)
			switch hitspot {
			case 1:
				act(libc.CString("@WYou raise your right hand above your head with it slightly cupped. Dark @rred@W energy from the Eldritch Star begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by your ki flows up your left arm as you raise it. You slam both of your hands together, combining the energy, and then toss your @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then toss a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at YOU! It engulfs your chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then tosses a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S chest!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou raise your right hand above your head with it slightly cupped. Dark @rred@W energy from the Eldritch Star begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by your ki flows up your left arm as you raise it. You slam both of your hands together, combining the energy, and then toss your @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then toss a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at YOU! It engulfs your head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then tosses a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S head!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 0))
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou raise your right hand above your head with it slightly cupped. Dark @rred@W energy from the Eldritch Star begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by your ki flows up your left arm as you raise it. You slam both of your hands together, combining the energy, and then toss your @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then toss a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at YOU! It engulfs your body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then tosses a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if (ch.Bonuses[BONUS_SOFT]) != 0 {
					dmg *= int64(calc_critical(ch, 2))
				}
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou raise your right hand above your head with it slightly cupped. Dark @rred@W energy from the Eldritch Star begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by your ki flows up your left arm as you raise it. You slam both of your hands together, combining the energy, and then toss your @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then toss a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at YOU! It engulfs your arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then tosses a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 1)
				hurt(0, 180, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@WYou raise your right hand above your head with it slightly cupped. Dark @rred@W energy from the Eldritch Star begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by your ki flows up your left arm as you raise it. You slam both of your hands together, combining the energy, and then toss your @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then toss a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at YOU! It engulfs your leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W raises $s right hand above $s head with it slightly cupped. Dark @rred@W energy begins to pool there and form a growing orb of energy. At the same time @mpurple@W arcs of electricity formed by $s ki flows up $s left arm as $e raises it. Suddenly $e slams both of $s hands together, combining the energy, and then tosses a @YS@yta@Yr @rB@Rr@De@ra@Rk@De@rr@W at @c$N@W! It engulfs $S leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= int64(calc_critical(ch, 1))
				dam_eq_loc(vict, 2)
				hurt(1, 180, ch, vict, nil, dmg, 1)
			}
			if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 || GET_LEVEL(ch) >= 100 {
				send_to_char(ch, libc.CString("The returning Eldritch energy blesses you with some experience. @D[@G%s@D]@n\r\n"), add_commas(theft))
				ch.Exp += theft * 2
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
		dmg = damtype(ch, 46, skill, attperc)
		dmg /= 10
		act(libc.CString("@WYou fire Star Breaker at $p@W!@n"), FALSE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@C$n@W fires Star Breaker at $p@W!@n"), FALSE, ch, obj, nil, TO_ROOM)
		hurt(0, 0, ch, nil, obj, dmg, 0)
		pcost(ch, attperc, 0)
	} else {
		send_to_char(ch, libc.CString("Error! Please report.\r\n"))
		return
	}
}
