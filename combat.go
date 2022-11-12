package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func damage_weapon(ch *char_data, obj *obj_data, vict *char_data) {
	if obj != nil {
		if OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) {
			return
		}
	}
	var ranking int = 0
	var material int = 1
	var PL10 int64 = 2000000000
	var PL9 int64 = 2000000000
	var PL8 int64 = 2000000000
	PL10 = PL10 * 5
	PL9 = PL9 * 4
	PL8 = PL8 * 2
	if vict.Hit >= PL10 {
		ranking = 10
	} else if vict.Hit >= PL9 {
		ranking = 9
	} else if vict.Hit >= PL8 {
		ranking = 8
	} else if vict.Hit >= 2000000000 {
		ranking = 7
	} else if vict.Hit >= 500000000 {
		ranking = 6
	} else if vict.Hit >= 250000000 {
		ranking = 5
	} else if vict.Hit >= 100000000 {
		ranking = 4
	} else if vict.Hit >= 50000000 {
		ranking = 3
	} else if vict.Hit >= 25000000 {
		ranking = 2
	} else if vict.Hit >= 1000000 {
		ranking = 1
	}
	switch obj.Value[VAL_ALL_MATERIAL] {
	case MATERIAL_STEEL:
		material = 4
	case MATERIAL_IRON:
		fallthrough
	case MATERIAL_COPPER:
		fallthrough
	case MATERIAL_BRASS:
		fallthrough
	case MATERIAL_METAL:
		material = 2
	case MATERIAL_SILVER:
		material = 5
	case MATERIAL_KACHIN:
		material = 9
	case MATERIAL_CRYSTAL:
		material = 7
	case MATERIAL_DIAMOND:
		material = 8
	case MATERIAL_PAPER:
		fallthrough
	case MATERIAL_COTTON:
		fallthrough
	case MATERIAL_SATIN:
		fallthrough
	case MATERIAL_SILK:
		fallthrough
	case MATERIAL_BURLAP:
		fallthrough
	case MATERIAL_VELVET:
		fallthrough
	case MATERIAL_HEMP:
		fallthrough
	case MATERIAL_WAX:
		material = 0
	default:
	}
	var result int = ranking - material
	if AFF_FLAGGED(ch, AFF_CURSE) {
		result += 3
	} else if AFF_FLAGGED(ch, AFF_BLESS) && rand_number(1, 3) == 3 {
		if result > 1 {
			result = 1
		}
	} else if AFF_FLAGGED(ch, AFF_BLESS) {
		result = 0
	}
	if GET_SKILL(ch, SKILL_HANDLING) >= axion_dice(10) {
		act(libc.CString("@GYour superior handling prevents @C$p@G from being damaged.@n"), TRUE, ch, obj, nil, TO_CHAR)
		act(libc.CString("@g$n's@G superior handling prevents @C$p@G from being damaged.@n"), TRUE, ch, obj, nil, TO_ROOM)
		result = 0
	}
	if result > 0 {
		obj.Value[VAL_ALL_HEALTH] -= result
		if (obj.Value[VAL_ALL_HEALTH]) <= 0 {
			act(libc.CString("@RYour @C$p@R shatters on @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n's@R @C$p@R shatters on YOUR body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n's@R @C$p@R shatters on @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			obj.Extra_flags[int(ITEM_BROKEN/32)] |= bitvector_t(1 << (int(ITEM_BROKEN % 32)))
			perform_remove(vict, 16)
			perform_remove(vict, 17)
		} else if result >= 8 {
			act(libc.CString("@RYour @C$p@R cracks loudly from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n's@R @C$p@R cracks loudly from striking YOUR body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n's@R @C$p@R cracks loudly from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		} else if result >= 6 {
			act(libc.CString("@RYour @C$p@R chips from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n's@R @C$p@R cracks from striking YOUR body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n's@R @C$p@R cracks from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		} else if result >= 3 {
			act(libc.CString("@RYour @C$p@R loses a piece from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n's@R @C$p@R loses a piece from striking YOUR body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n's@R @C$p@R loses a piece from striking @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		} else if result >= 6 {
			act(libc.CString("@RYour @C$p@R has a nick in it from hitting @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@r$n's@R @C$p@R has a nick in it from hitting YOUR body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n's@R @C$p@R has a nick in it from hitting @r$N's@R body!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
		}
	}
}
func handle_multihit(ch *char_data, vict *char_data) {
	var (
		perc int = int(ch.Aff_abils.Dex)
		prob int = int(vict.Aff_abils.Dex)
	)
	prob += rand_number(1, 15)
	perc += rand_number(-5, 5)
	if ch.Throws >= 3 {
		ch.Throws = 0
		return
	} else if ch.Throws == -1 {
		ch.Throws = 0
	}
	if ch.Race == RACE_KONATSU {
		perc *= int(1.5)
	}
	if ch.Race == RACE_BIO && ((ch.Genome[0]) == 8 || (ch.Genome[1]) == 8) {
		perc *= int(1.4)
	}
	if IS_NPC(ch) {
		perc -= int(float64(perc) * 0.3)
	}
	if ch.Lastattack == -1 {
		perc *= int(0.75)
	}
	var amt int = 70
	if GET_SKILL(ch, SKILL_STYLE) >= 100 {
		amt -= int(float64(amt) * 0.1)
	} else if GET_SKILL(ch, SKILL_STYLE) >= 80 {
		amt -= int(float64(amt) * 0.08)
	} else if GET_SKILL(ch, SKILL_STYLE) >= 60 {
		amt -= int(float64(amt) * 0.06)
	} else if GET_SKILL(ch, SKILL_STYLE) >= 40 {
		amt -= int(float64(amt) * 0.04)
	} else if GET_SKILL(ch, SKILL_STYLE) >= 20 {
		amt -= int(float64(amt) * 0.02)
	}
	if axion_dice(0) < amt {
		prob += 500
	}
	if perc >= prob {
		var buf [2048]byte
		act(libc.CString("@Y...in a lightning flash of speed you attack @y$N@Y again!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@Y...in a lightning flash of speed @y$n@Y attacks YOU again!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@Y...in a lightning flash of speed @y$n@Y attacks @y$N@Y again!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		ch.Throws += 1
		ch.Act[int(PLR_MULTIHIT/32)] |= bitvector_t(1 << (int(PLR_MULTIHIT % 32)))
		if ch.Combo > -1 {
			switch ch.Combo {
			case 0:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_punch(ch, &buf[0], 0, 0)
			case 1:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_kick(ch, &buf[0], 0, 0)
			case 2:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_elbow(ch, &buf[0], 0, 0)
			case 3:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_knee(ch, &buf[0], 0, 0)
			case 4:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_roundhouse(ch, &buf[0], 0, 0)
			case 5:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_uppercut(ch, &buf[0], 0, 0)
			case 6:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_slam(ch, &buf[0], 0, 0)
				ch.Throws += 1
			case 8:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_heeldrop(ch, &buf[0], 0, 0)
				ch.Throws += 1
			case 51:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_bash(ch, &buf[0], 0, 0)
				ch.Throws += 1
			case 52:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_head(ch, &buf[0], 0, 0)
			case 56:
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_tailwhip(ch, &buf[0], 0, 0)
				ch.Throws += 1
			}
		} else {
			if ch.Lastattack == -1 {
				stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
				do_attack(ch, &buf[0], 0, 0)
			} else {
				if rand_number(1, 3) == 2 && GET_SKILL(ch, SKILL_KICK) > 0 {
					stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
					do_kick(ch, &buf[0], 0, 0)
				} else {
					stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
					do_punch(ch, &buf[0], 0, 0)
				}
			}
		}
	}
}
func handle_defender(vict *char_data, ch *char_data) int {
	var result int = FALSE
	if vict.Defender != nil {
		var (
			def    *char_data = vict.Defender
			defnum int64      = int64((float64(GET_SPEEDI(def)) * 0.01) * float64(rand_number(-10, 10)))
			chnum  int64      = int64((float64(GET_SPEEDI(ch)) * 0.01) * float64(rand_number(-5, 10)))
		)
		if GET_SPEEDI(def)+int(defnum) > GET_SPEEDI(ch)+int(chnum) && def.In_room == vict.In_room && def.Position > POS_SITTING {
			act(libc.CString("@YYou move to and manage to intercept the attack aimed at @y$N@Y!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$n@Y moves to and manages to intercept the attack aimed at YOU!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@y$n@Y moves to and manages to intercept the attack aimed at @y$N@Y!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_NOTVICT)
			result = TRUE
		} else if def.In_room == vict.In_room && def.Position > POS_SITTING {
			act(libc.CString("@YYou move to intercept the attack aimed at @y$N@Y, but just not fast enough!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$n@Y moves to intercept the attack aimed at YOU, but $e wasn't fast enough!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@y$n@Y moves to intercept the attack aimed at @y$N@Y, but $e wasn't fast enough!@n"), TRUE, def, nil, unsafe.Pointer(vict), TO_NOTVICT)
		}
	}
	return result
}
func handle_disarm(ch *char_data, vict *char_data) {
	var (
		roll1   int = rand_number(-10, 10)
		roll2   int = rand_number(-10, 10)
		handled int = FALSE
	)
	roll1 += int(ch.Aff_abils.Str + ch.Aff_abils.Dex)
	roll2 += int(vict.Aff_abils.Str + vict.Aff_abils.Dex)
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_THANDW) {
			roll1 += 5
		}
	}
	if rand_number(1, 100) <= 50 && ch.Race != RACE_KONATSU {
		roll1 = -500
	} else if rand_number(1, 100) <= 75 {
		roll1 *= int(1.5)
	}
	if vict.Race == RACE_KONATSU {
		roll1 *= int(0.75)
	}
	if GET_SKILL(ch, SKILL_HANDLING) >= axion_dice(10) {
		handled = TRUE
	}
	if roll1 < roll2 {
		var obj *obj_data
		if (ch.Equipment[WEAR_WIELD1]) != nil && handled == FALSE {
			obj = ch.Equipment[WEAR_WIELD1]
			act(libc.CString("@y$N@Y manages to disarm you! The @w$p@Y falls from your grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$N@Y manages to disarm @R$n@Y! The @w$p@Y falls from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@YYou manage to disarm @R$n@Y! The @w$p@Y falls from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			perform_remove(ch, 16)
			if GET_OBJ_VNUM(obj) != 0x4E82 {
				obj_from_char(obj)
				obj_to_room(obj, ch.In_room)
			}
		} else if (ch.Equipment[WEAR_WIELD1]) != nil && handled == TRUE {
			obj = ch.Equipment[WEAR_WIELD1]
			act(libc.CString("@y$N@Y almosts disarms you, but your handling of @w$p@Y saves it from slipping from your grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$N@Y almost disarms @R$n@Y, but $s handling of @w$p@Y saves it from slipping from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@YYou almost disarm @R$n@Y, but $s handling of @w$p@Y saves it from slipping from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		} else if (ch.Equipment[WEAR_WIELD2]) != nil && handled == TRUE {
			obj = ch.Equipment[WEAR_WIELD2]
			act(libc.CString("@y$N@Y almosts disarms you, but your handling of @w$p@Y saves it from slipping from your grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$N@Y almost disarms @R$n@Y, but $s handling of @w$p@Y saves it from slipping from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@YYou almost disarm @R$n@Y, but $s handling of @w$p@Y saves it from slipping from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
		} else if (ch.Equipment[WEAR_WIELD2]) != nil {
			obj = ch.Equipment[WEAR_WIELD2]
			act(libc.CString("@y$N@Y manages to disarm you! The @w$p@Y falls from your grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@y$N@Y manages to disarm @R$n@Y! The @w$p@Y falls from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_NOTVICT)
			act(libc.CString("@YYou manage to disarm @R$n@Y! The @w$p@Y falls from $s grasp!@n"), TRUE, ch, obj, unsafe.Pointer(vict), TO_VICT)
			perform_remove(ch, 17)
			if GET_OBJ_VNUM(obj) != 0x4E82 {
				obj_from_char(obj)
				obj_to_room(obj, ch.In_room)
			}
		}
	}
}
func combine_attacks(ch *char_data, vict *char_data) {
	var (
		f         *follow_type
		chbuf     [2048]byte
		victbuf   [2048]byte
		rmbuf     [2048]byte
		bonus     int64   = 0
		maxki     float64 = 0.0
		totalmem  int     = 1
		attspd    int     = 0
		blockable int     = TRUE
		same      int     = TRUE
		attsk     int     = 0
		attavg    int     = GET_SKILL(ch, attack_skills[ch.Combine])
		burn      int     = FALSE
		shocked   int     = FALSE
	)
	switch ch.Combine {
	case 0:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You cup your hands at your sides and a ball of @Benergy@W forms there. You chant @B'@CKaaaaameeeehaaaameeee@B'@W and then fire a @RKamehameha@W wave at @r$N@W while screaming @B'@CHAAAAAAAAAAAAAAAAAAAAA!@B'@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W cups $s hands be $s side and chants @B'@CKaaaameeeehaaaameee@B'@W. A ball of energy forms in $s hands and he quickly brings them forward and fires a @RKamehameha @Wat @rYOU@W while screaming @B'@CHAAAAAAAAAAAAAAAAAAAAA!@B'@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W cups $s hands be $s side and chants @B'@CKaaaameeeehaaaameee@B'@W. A ball of energy forms in $s hands and he quickly brings them forward and fires a @RKamehameha @Wat @r$N@W while screaming @B'@CHAAAAAAAAAAAAAAAAAAAAA!@B'@n@n")
		maxki = 0.15
		attspd += 2
		bonus += int64(float64(ch.Max_move) * 0.02)
	case 1:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You throw your hands forward and launch a purple beam of energy at @r$N@n while shouting @B'@mGalik Gun!@B'@W")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W throws $s arms forward and launches a purple beam at @rYOU@W while shouting @B'@mGalik Gun!@B'@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W throws $s arms forward and launches a purple beam at @r$N@W while shouting @B'@mGalik Gun!@B'@n")
		maxki = 0.15
		attspd += 1
		bonus += int64(float64(ch.Max_mana) * 0.5)
	case 2:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You raise your hands above your head with one resting atop the other and begin to pour your charged energy to that point. As soon as the energy is ready you shout @B'@RMasenko Ha!@B'@W and bringing your hands down you launch a bright reddish orange beam at @r$N@W!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W raises $s hands above $s head and energy quicly pools there. Suddenly $e brings $s hands down and shouts @B'@RMasenko Ha!@B'@W as a bright reddish orange beam launches toward @rYOU!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W raises $s hands above $s head and energy quicly pools there. Suddenly $e brings $s hands down and shouts @B'@RMasenko Ha!@B'@W as a bright reddish orange beam launches toward @r$N!@n")
		maxki = 0.15
		attspd += 1
		bonus += int64(float64(ch.Max_mana) * 0.5)
	case 3:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! With a quick motion you point at @r$N@W and launch a lightning fast @MDeathbeam@W!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! With a quick motion $e points $s finger at @rYOU@W and launches a lightning fast @MDeathbeam@W!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! With a quick motion $e points $s finger at @r$N@W and launches a lightning fast @MDeathbeam@W!@n")
		maxki = 0.1
		attspd += 4
	case 4:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! With your energy ready you breath out toward @r$N@W jets of incredibly hot flames in the form of a deadly @rHonoo@W!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! Sudden jets of flame burst forth from $s mouth at @RYOU@W in the form of a deadly @rHonoo@W!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! Sudden jets of flame burst forth from $s mouth at @R$N@W in the form of a deadly @rHonoo@W!@n")
		maxki = 0.125
		attspd += 2
		burn = TRUE
	case 5:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! With your energy prepared you poor it into your blade and accelerate your body to incredible speeds toward @r$N! You leave two glowing green marks behind on $S body in a single instant as your @GTwin Slash@W hits!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! Raising $s sword @Y$n@W accelerates toward @rYOU@W with incredible speed! Two glowing green slashes are left on YOUR body from $s successful @GTwin Slash@W!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! Raising $s sword @Y$n@W accelerates toward @r$N@W with incredible speed! Two glowing green slashes are left on @R$N's@W body from the successful @GTwin Slash@W!@n")
		maxki = 0.125
		attspd += 2
		blockable = FALSE
	case 6:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You stick one of each of your hands in your armpits and detach them. With your hands detached your point the exposed arm cannons at @r$N@W and launch a massive @RHell Flash@W at $M!")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W sticks one of each of $s hands in $s armpits and detaches them there. With the hands detached $e aims $s exposed arm cannons at @RYOU@W and launches a massive @RHell Flash@W!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W sticks one of each of $s hands in $s armpits and detaches them there. With the hands detached $e aims $s exposed arm cannons at @R$N@W and launches a massive @RHell Flash@W!@n")
		maxki = 0.2
		attspd += 1
	case 7:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! With your energy ready you look at @R$N@W as the blue light of your @CPsychic Blast@W launches from your head toward $S!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! A blue light, identifying a @CPsychic Blast@W, launches from @Y$n's@W toward @RYOUR HEAD@W!")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! A blue light, identifying a @CPsychic Blast@W, launches from @Y$n's@W toward @R$N's@W head!")
		maxki = 0.125
		attspd += 1
		shocked = TRUE
	case 8:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! Pooling your energy you form a large ball of red energy above an upraised palm. Slamming your other hand into it you launch it toward @r$N@W while shouting @B'@RCrusher Ball@B'@W!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W raises a palm above his head and red energy begins to pool there. As the energy completes the formation of a ball @Y$n@W slams $s other hand into it and launches it at @rYOU@W while shouting @B'@RCrusher Ball@B'@W!")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W raises a palm above his head and red energy begins to pool there. As the energy completes the formation of a ball @Y$n@W slams $s other hand into it and launches it at @r$N@W while shouting @B'@RCrusher Ball@B'@W!")
		maxki = 0.2
	case 9:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! Using your energy to form a ball of water between your hands you then raise the ball above your head. Several spiked of ice form from the ball of water and you hurl them at @r$N@W!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! Forming a ball of water between $s palms with $s energy @Y$n@W then raises the ball of water above $s head. Suddenly several spikes of ice form from the water and $e launches them at @rYOU@W!")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! Forming a ball of water between $s palms with $s energy @Y$n@W then raises the ball of water above $s head. Suddenly several spikes of ice form from the water and $e launches them at @r$N@W!")
		maxki = 0.14
	case 10:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You form a triangle with your hands and aim the center of the triangle at @r$N@W. With the sudden shout @B'@YTribeam@B'@W you release your prepared energy at $M in the form of a beam!")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W forms a triangle with $s hands and aims the center at @rYOU@W! With the sudden shout @B'@YTribeam@B'@W a large beam of energy flashes toward @rYOU!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W forms a triangle with $s hands and aims the center at @r$N@W! With the sudden shout @B'@YTribeam@B'@W a large beam of energy flashes toward @r$N!@n")
		maxki = 0.2
		attspd += 2
		bonus += int64(float64(ch.Max_hit) * 0.5)
	case 11:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group, you call out to your allies to launch a group attack! You raise your right hand above your head as dark red energy begins to pool in your slightly cupped hand, while purple arcs of electricity flow up your left arm. Slamming both hands together, you shout @B'@YStarbreaker@B'@W and release your prepared energy at $M in the form of a ball!@n")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W raises $s right hand, pooling dark red energy in the palm. @Y$n@W slams both their hands together, shouting @B'@YStarbreaker@B'@W, a ball of energy flashes toward @rYOU!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W raises their right hand above their head, pooling dark red energy. @Y$n@W slams both their hands together, shouting @B'@YStarbreaker@B'@W, a ball of energy flashes toward @r$N!@n")
		maxki = 0.2
		bonus += int64(float64(ch.Max_mana) * 0.6)
	case 12:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You open your mouth and aim at @r$N@W. You grunt as you release your prepared energy at $M in the form of a beam!")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W opens $s mouth and aims at @rYOUW! With the sudden grunt, a large beam flashes towards @rYOU@n!")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W opens $s mouth and aims at @rYOUW! With the sudden grunt, a large beam flashes toward @r$N@n!")
		maxki = 0.125
		attspd += 35
	case 13:
		stdio.Sprintf(&chbuf[0], "@WPositioning yourself in the center of your group you call out to your allies to launch a group attack! You slam your hands together and aim at @r$n@W. With the sudden shout @B'@YRenzoku Energy Dan@B'@W you release your prepared energy in the form of hundreds of ki blasts!")
		stdio.Sprintf(&victbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @rYOU@W! @Y$n@W slams both $s hands together and aims at @rYOUW! With the sudden shout @B'@YRenzoku Energy Dan@B'@W hundreds of ki blasts flash towards @rYOU!@n")
		stdio.Sprintf(&rmbuf[0], "@Y$n@W calls out to $s allies to launch a group attack against @r$N@W! @Y$n@W slams $s hands together and aims at @r$N@W! With the sudden shout @B'@YRenzoku Energy Dan@B'@W hundreds of ki blasts flash towards @r$N!@n")
		maxki = 0.125
		attspd += 6
	default:
		send_to_imm(libc.CString("ERROR: Combine attacks failure for: %s"), GET_NAME(ch))
		send_to_char(ch, libc.CString("An error has been logged. Be patient while waiting for Iovan's response.\r\n"))
		return
	}
	var totki int64 = 0
	act(&chbuf[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
	act(&victbuf[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
	act(&rmbuf[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
	if float64(ch.Charge) >= float64(ch.Max_mana)*maxki {
		totki += int64(float64(ch.Max_mana) * maxki)
		ch.Charge -= int64(float64(ch.Max_mana) * maxki)
	} else {
		totki += ch.Charge
		ch.Charge = 0
	}
	for f = ch.Followers; f != nil; f = f.Next {
		if !AFF_FLAGGED(f.Follower, AFF_GROUP) {
			continue
		} else {
			if f.Follower.Combine != ch.Combine {
				same = FALSE
			}
			if float64(f.Follower.Charge) >= float64(f.Follower.Max_mana)*maxki {
				totki += int64(float64(f.Follower.Max_mana) * maxki)
				f.Follower.Charge -= int64(float64(f.Follower.Max_mana) * maxki)
			} else {
				totki += f.Follower.Charge
				f.Follower.Charge = 0
			}
			totalmem += 1
			attavg += GET_SKILL(f.Follower, attack_skills[f.Follower.Combine])
			var folbuf [2048]byte
			var folbuf2 [2048]byte
			stdio.Sprintf(&folbuf[0], "@Y$n@W times and merges $s @B'@R%s@B'@W into the group attack!@n", attack_names_comp[f.Follower.Combine])
			stdio.Sprintf(&folbuf2[0], "@WYou time and merge your @B'@R%s@B'@W into the group attack!@n", attack_names_comp[f.Follower.Combine])
			act(&folbuf[0], TRUE, f.Follower, nil, nil, TO_ROOM)
			act(&folbuf2[0], TRUE, f.Follower, nil, nil, TO_CHAR)
		}
	}
	totki += bonus
	if same == TRUE {
		totki += bonus
	}
	attsk = attavg / totalmem
	if ch.Combine != 5 {
		if attspd+attsk < GET_SKILL(vict, SKILL_DODGE)+int(ch.Aff_abils.Cha/10) {
			act(libc.CString("@GYou manage to dodge nimbly through the combined attack of your enemies!@n"), TRUE, vict, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@G manages to dodge nimbly through the combined attack!@n"), TRUE, vict, nil, nil, TO_ROOM)
			return
		} else if blockable == TRUE && attspd+attsk < GET_SKILL(vict, SKILL_BLOCK)+int(ch.Aff_abils.Str/10) {
			act(libc.CString("@GYou manage to effectivly block the combined attack of your enemies with the help of your great strength!@n"), TRUE, vict, nil, nil, TO_CHAR)
			act(libc.CString("@r$n@G manages to dodge nimbly through the combined attack!@n"), TRUE, vict, nil, nil, TO_ROOM)
			return
		}
	}
	if burn == TRUE {
		if !AFF_FLAGGED(vict, AFF_BURNED) && rand_number(1, 4) == 3 && vict.Race != RACE_DEMON && (vict.Bonuses[BONUS_FIREPROOF]) == 0 {
			send_to_char(vict, libc.CString("@RYou are burned by the attack!@n\r\n"))
			send_to_char(ch, libc.CString("@RThey are burned by the attack!@n\r\n"))
			vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
		} else if (vict.Bonuses[BONUS_FIREPROOF]) != 0 || vict.Race == RACE_DEMON {
			send_to_char(ch, libc.CString("@RThey appear to be fireproof!@n\r\n"))
		} else if (vict.Bonuses[BONUS_FIREPRONE]) != 0 {
			send_to_char(vict, libc.CString("@RYou are extremely flammable and are burned by the attack!@n\r\n"))
			send_to_char(ch, libc.CString("@RThey are easily burned!@n\r\n"))
			vict.Affected_by[int(AFF_BURNED/32)] |= 1 << (int(AFF_BURNED % 32))
		}
	}
	if shocked == TRUE {
		if !AFF_FLAGGED(vict, AFF_SHOCKED) && rand_number(1, 4) == 4 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
			act(libc.CString("@MYour mind has been shocked!@n"), TRUE, vict, nil, nil, TO_CHAR)
			act(libc.CString("@M$n@m's mind has been shocked!@n"), TRUE, vict, nil, nil, TO_ROOM)
			vict.Affected_by[int(AFF_SHOCKED/32)] |= 1 << (int(AFF_SHOCKED % 32))
		}
	}
	hurt(0, 0, ch, vict, nil, totki, 1)
	if same == TRUE {
		for f = ch.Followers; f != nil; f = f.Next {
			send_to_char(f.Follower, libc.CString("@YS@yy@Yn@ye@Yr@yg@Yi@ys@Yt@yi@Yc @yB@Yo@yn@Yu@ys@Y!@n\r\n"))
		}
		send_to_char(ch, libc.CString("@YS@yy@Yn@ye@Yr@yg@Yi@ys@Yt@yi@Yc @yB@Yo@yn@Yu@ys@Y!@n\r\n"))
	}
}
func check_ruby(ch *char_data) int {
	var (
		obj      *obj_data
		next_obj *obj_data = nil
		ruby     *obj_data = nil
		found    int       = 0
	)
	for obj = ch.Carrying; obj != nil; obj = next_obj {
		next_obj = obj.Next_content
		if found == 0 && GET_OBJ_VNUM(obj) == 6600 {
			if OBJ_FLAGGED(obj, ITEM_HOT) {
				found = 1
				ruby = obj
			}
		}
	}
	if found > 0 {
		act(libc.CString("@RYour $p@R flares up and disappears. Your fire attack has been aided!@n"), TRUE, ch, ruby, nil, TO_CHAR)
		act(libc.CString("@R$n's@R $p@R flares up and disappears!@n"), TRUE, ch, ruby, nil, TO_ROOM)
		extract_obj(ruby)
		return 1
	} else {
		return 0
	}
}
func combo_damage(ch *char_data, damage int64, type_ int) int64 {
	var bonus int64 = 0
	if type_ == 0 {
		var hits int = ch.Combhits
		if hits >= 30 {
			bonus += int64(float64(damage) * (float64(hits) * 0.15))
			bonus += damage * 12
		} else if hits >= 20 {
			bonus = int64(float64(damage) * (float64(hits) * 0.1))
			bonus += damage * 10
		} else if hits >= 10 {
			bonus = int64(float64(damage) * (float64(hits) * 0.1))
			bonus += damage * 5
		} else if hits >= 6 {
			bonus = int64(float64(damage) * (float64(hits) * 0.1))
			bonus += int64(float64(damage) * 1.5)
		} else if hits >= 2 {
			bonus = int64(float64(damage) * (float64(hits) * 0.05))
			bonus += int64(float64(damage) * 0.2)
		}
	} else if type_ == 1 {
		bonus = damage * 15
	}
	return bonus
}
func roll_balance(ch *char_data) int {
	var chance int = 0
	if IS_NPC(ch) {
		if GET_LEVEL(ch) >= 100 {
			chance = rand_number(80, 100)
		} else if GET_LEVEL(ch) >= 80 {
			chance = rand_number(75, 90)
		} else if GET_LEVEL(ch) >= 70 {
			chance = rand_number(70, 80)
		} else if GET_LEVEL(ch) >= 60 {
			chance = rand_number(65, 75)
		} else if GET_LEVEL(ch) >= 50 {
			chance = rand_number(50, 60)
		}
	} else {
		if GET_SKILL(ch, SKILL_BALANCE) > 50 {
			chance = GET_SKILL(ch, SKILL_BALANCE)
		}
	}
	return chance
}
func handle_knockdown(ch *char_data) {
	var chance int = 0
	if IS_NPC(ch) {
		if GET_LEVEL(ch) >= 100 {
			chance = rand_number(35, 45)
		} else if GET_LEVEL(ch) >= 90 {
			chance = rand_number(25, 35)
		} else if GET_LEVEL(ch) >= 70 {
			chance = rand_number(15, 25)
		} else if GET_LEVEL(ch) >= 50 {
			chance = rand_number(10, 15)
		} else if GET_LEVEL(ch) >= 30 {
			chance = rand_number(5, 10)
		}
	} else {
		chance = int(float64(GET_SKILL(ch, SKILL_BALANCE)) * 0.5)
	}
	if chance > axion_dice(0) {
		act(libc.CString("@mYou are @GALMOST@m knocked off your feet, but your great balance helps you keep your footing!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n@m is @GALMOST@m knocked off $s feet, but $s great balance helps $m keep $s footing!@n"), TRUE, ch, nil, nil, TO_ROOM)
	} else {
		act(libc.CString("@mYou are knocked off your feet!@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@W$n@m is knocked off $s feet!@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Position = POS_SITTING
	}
}
func boom_headshot(ch *char_data) int {
	var skill int = int(ch.Skills[SKILL_TWOHAND])
	if skill >= 100 && rand_number(1, 5) >= 3 {
		return 1
	} else if skill < 100 && skill >= 75 && rand_number(1, 5) == 5 {
		return 1
	} else if skill < 75 && skill >= 50 && rand_number(1, 6) == 6 {
		return 1
	} else {
		return 0
	}
}
func gun_dam(ch *char_data, wlvl int) int64 {
	var dmg int64 = 100
	switch wlvl {
	case 1:
		dmg = 50
	case 2:
		dmg = 200
	case 3:
		dmg = 750
	case 4:
		dmg = 2000
	case 5:
		dmg = 7500
	}
	if GET_SKILL(ch, SKILL_GUN) >= 100 {
		dmg *= 2
	} else if GET_SKILL(ch, SKILL_GUN) >= 50 {
		dmg += int64(float64(dmg) * 0.5)
	}
	var dmg_prior int64 = 0
	dmg_prior = (dmg * int64(ch.Aff_abils.Dex)) * int64((GET_LEVEL(ch)/5)+1)
	if float64(dmg_prior) <= float64(ch.Max_hit)*0.4 {
		dmg = dmg_prior
	} else {
		dmg = int64(float64(ch.Max_hit) * 0.4)
	}
	return dmg
}
func club_stamina(ch *char_data, vict *char_data, wlvl int, dmg int64) {
	var (
		drain   float64 = 0.0
		drained int64   = 0
	)
	switch wlvl {
	case 1:
		drain = 0.05
	case 2:
		drain = 0.1
	case 3:
		drain = 0.15
	case 4:
		drain = 0.2
	case 5:
		drain = 0.25
	}
	if GET_SKILL(ch, SKILL_CLUB) >= 100 {
		drain += 0.1
	} else if GET_SKILL(ch, SKILL_CLUB) >= 50 {
		drain += 0.05
	}
	drained = int64(float64(dmg) * drain)
	vict.Move -= drained
	if vict.Move < 0 {
		vict.Move = 0
	}
	send_to_char(ch, libc.CString("@D[@YVictim's @GStamina @cLoss@W: @g%s@D]@n\r\n"), add_commas(drained))
	send_to_char(vict, libc.CString("@D[@rYour @GStamina @cLoss@W: @g%s@D]@n\r\n"), add_commas(drained))
}
func backstab(ch *char_data, vict *char_data, wlvl int, dmg int64) int {
	var (
		chance       int     = 0
		roll_to_beat int     = rand_number(1, 100)
		bonus        float64 = 0.0
	)
	if ch.Backstabcool > 0 {
		return 0
	}
	switch wlvl {
	case 1:
		chance = 10
		bonus += 0.5
	case 2:
		chance = 15
		bonus += 2
	case 3:
		chance = 20
		bonus += 3
	case 4:
		chance = 25
		bonus += 4
	case 5:
		chance = 30
		bonus += 4
	}
	if (ch.Bonuses[BONUS_POWERHIT]) != 0 {
		bonus += 2
	}
	if GET_SKILL(ch, SKILL_DAGGER) >= 100 {
		chance += 20
	} else if GET_SKILL(ch, SKILL_DAGGER) >= 50 {
		chance += 10
	}
	ch.Backstabcool = 10
	if chance >= roll_to_beat {
		var (
			attacker_roll int = GET_SKILL(ch, SKILL_MOVE_SILENTLY) + GET_SKILL(ch, SKILL_SPOT) + int(ch.Aff_abils.Dex) + rand_number(-5, 5)
			defender_roll int = GET_SKILL(vict, SKILL_SPOT) + GET_SKILL(vict, SKILL_LISTEN) + int(ch.Aff_abils.Dex) + rand_number(-5, 5)
		)
		if attacker_roll > defender_roll {
			act(libc.CString("@RYou manage to sneak behind @r$N@R and stab $M in the back!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RYou feel @r$n's@R dagger thrust into your back unexpectantly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@r$n@R sneaks up behind @r$N@R and stabs $M in the back!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			dmg += int64(float64(dmg) * bonus)
			hurt(0, 0, ch, vict, nil, dmg, 0)
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}
func cut_limb(ch *char_data, vict *char_data, wlvl int, hitspot int) {
	var (
		chance       int = 0
		decap        int = 0
		decapitate   int = FALSE
		roll_to_beat int = rand_number(1, 10000)
	)
	if wlvl == 1 {
		chance = 25
	} else if wlvl == 2 {
		chance = 50
	} else if wlvl == 3 {
		chance = 100
		decap = 5
	} else if wlvl == 4 {
		chance = 200
		decap = 10
	} else if wlvl == 5 {
		chance = 200
		decap = 50
	}
	if GET_SKILL(ch, SKILL_SWORD) >= 100 {
		chance += 100
	} else if GET_SKILL(ch, SKILL_SWORD) >= 50 {
		chance += 50
	}
	if decap >= roll_to_beat && hitspot == 4 {
		decapitate = TRUE
	} else if chance < roll_to_beat {
		return
	}
	if vict.Hit <= 1 {
		return
	}
	if decapitate == TRUE {
		act(libc.CString("@R$N's@r head is cut off in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@RYOUR head is cut off in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@R$N's@rhead is cut off in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		remove_limb(vict, 0)
		die(vict, ch)
		if AFF_FLAGGED(ch, AFF_GROUP) {
			group_gain(ch, vict)
		} else {
			solo_gain(ch, vict)
		}
		var corp [256]byte
		if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
			stdio.Sprintf(&corp[0], "all.zenni corpse")
			do_get(ch, &corp[0], 0, 0)
		}
		if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
			stdio.Sprintf(&corp[0], "all corpse")
			do_get(ch, &corp[0], 0, 0)
		}
		return
	} else {
		if !IS_NPC(vict) {
			if HAS_ARMS(vict) && rand_number(1, 2) == 2 {
				if (vict.Limb_condition[1]) > 0 {
					vict.Limb_condition[1] = 0
					if PLR_FLAGGED(vict, PLR_CLARM) {
						vict.Act[int(PLR_CLARM/32)] &= bitvector_t(^(1 << (int(PLR_CLARM % 32))))
					}
					act(libc.CString("@R$N@r loses $s left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					remove_limb(vict, 2)
				} else if (vict.Limb_condition[0]) > 0 {
					vict.Limb_condition[0] = 100
					if PLR_FLAGGED(vict, PLR_CRARM) {
						vict.Act[int(PLR_CRARM/32)] &= bitvector_t(^(1 << (int(PLR_CRARM % 32))))
					}
					act(libc.CString("@R$N@r loses $s right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					remove_limb(vict, 1)
				}
			} else {
				if (vict.Limb_condition[3]) > 0 {
					vict.Limb_condition[3] = 100
					if PLR_FLAGGED(vict, PLR_CLLEG) {
						vict.Act[int(PLR_CLLEG/32)] &= bitvector_t(^(1 << (int(PLR_CLLEG % 32))))
					}
					act(libc.CString("@R$N@r loses $s left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					remove_limb(vict, 4)
				} else if (vict.Limb_condition[2]) > 0 {
					vict.Limb_condition[2] = 100
					if PLR_FLAGGED(vict, PLR_CRLEG) {
						vict.Act[int(PLR_CRLEG/32)] &= bitvector_t(^(1 << (int(PLR_CRLEG % 32))))
					}
					act(libc.CString("@R$N@r loses $s right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					remove_limb(vict, 3)
				}
			}
		} else {
			if HAS_ARMS(vict) && rand_number(1, 2) == 2 {
				if MOB_FLAGGED(vict, MOB_LARM) {
					vict.Act[int(MOB_LARM/32)] &= bitvector_t(^(1 << (int(MOB_LARM % 32))))
					remove_limb(vict, 2)
					act(libc.CString("@R$N@r loses $s left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else if MOB_FLAGGED(vict, MOB_RARM) {
					vict.Act[int(MOB_RARM/32)] &= bitvector_t(^(1 << (int(MOB_RARM % 32))))
					remove_limb(vict, 1)
					act(libc.CString("@R$N@r loses $s right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				}
			} else {
				if MOB_FLAGGED(vict, MOB_LLEG) {
					vict.Act[int(MOB_LLEG/32)] &= bitvector_t(^(1 << (int(MOB_LLEG % 32))))
					remove_limb(vict, 4)
					act(libc.CString("@R$N@r loses $s left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				} else if MOB_FLAGGED(vict, MOB_RLEG) {
					vict.Act[int(MOB_RLEG/32)] &= bitvector_t(^(1 << (int(MOB_RLEG % 32))))
					remove_limb(vict, 3)
					act(libc.CString("@R$N@r loses $s right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@RYOU lose your right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@R$N@r loses $s right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				}
			}
		}
	}
}
func count_physical(ch *char_data) int {
	var count int = 0
	if GET_SKILL(ch, SKILL_PUNCH) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_KICK) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_KNEE) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_ELBOW) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_ROUNDHOUSE) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_SLAM) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_UPPERCUT) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_TAILWHIP) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_BASH) >= 1 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_HEADBUTT) >= 1 {
		count += 1
	}
	return count
}
func physical_mastery(ch *char_data) int {
	var count int = 22
	if GET_SKILL(ch, SKILL_PUNCH) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_KICK) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_KNEE) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_ELBOW) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_ROUNDHOUSE) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_SLAM) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_UPPERCUT) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_TAILWHIP) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_BASH) >= 100 {
		count += 1
	}
	if GET_SKILL(ch, SKILL_HEADBUTT) >= 100 {
		count += 1
	}
	if count == 26 {
		count += 1
	} else if count >= 27 {
		count += 2
	}
	return count
}
func advanced_energy(ch *char_data, dmg int64) int64 {
	if ch == nil {
		return FALSE
	}
	var rate float64 = 0.0
	var count int = GET_LEVEL(ch)
	var add int64 = 0
	if (ch.Bonuses[BONUS_LEECH]) != 0 {
		rate = float64(count) * 0.2
		if rate > 0.0 {
			add = int64(float64(dmg) * rate)
			if ch.Charge+add > ch.Max_mana {
				if ch.Charge < ch.Max_mana {
					ch.Charge = ch.Max_mana
					act(libc.CString("@MYou leech some of the energy away!@n"), TRUE, ch, nil, nil, TO_CHAR)
					act(libc.CString("@m$n@M leeches some of the energy away!@n"), TRUE, ch, nil, nil, TO_ROOM)
				} else {
					send_to_char(ch, libc.CString("@MYou can't leech because there is too much charged energy for you to handle!@n\r\n"))
				}
			} else {
				ch.Charge += add
				act(libc.CString("@MYou leech some of the energy away!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@m$n@M leeches some of the energy away!@n"), TRUE, ch, nil, nil, TO_ROOM)
			}
		}
	}
	if (ch.Bonuses[BONUS_INTOLERANT]) != 0 {
		rate = float64(count) * 0.2
		if rate > 0.0 {
			if ch.Charge > 0 && rand_number(1, 100) <= 10 {
				act(libc.CString("@MThe attack causes your weak control to slip and you are shocked by your own charged energy!@n"), TRUE, ch, nil, nil, TO_CHAR)
				act(libc.CString("@m$n@M suffers shock from their own charged energy!@n"), TRUE, ch, nil, nil, TO_ROOM)
				ch.Hit -= ch.Charge / 4
				if ch.Hit <= 0 {
					ch.Hit = 1
				}
			}
			add = int64(float64(dmg) * rate)
		}
	}
	return add
}
func roll_accuracy(ch *char_data, skill int, kiatt bool) int {
	if !IS_NPC(ch) {
		if (ch.Bonuses[BONUS_ACCURATE]) != 0 {
			if int(libc.BoolToInt(kiatt)) == TRUE {
				skill += int(float64(skill) * 0.1)
			} else {
				skill += int(float64(skill) * 0.2)
			}
		} else if (ch.Bonuses[BONUS_POORDEPTH]) != 0 {
			if int(libc.BoolToInt(kiatt)) == TRUE {
				skill -= int(float64(skill) * 0.1)
			} else {
				skill -= int(float64(skill) * 0.2)
			}
		}
	}
	if skill < 40 {
		skill += rand_number(3, 10)
	}
	return skill
}
func calc_critical(ch *char_data, loc int) float64 {
	var (
		roll  int     = rand_number(1, 100)
		multi float64 = 1
	)
	if loc == 0 {
		if (ch.Bonuses[BONUS_POWERHIT]) != 0 {
			if roll >= 15 {
				multi = 4
			} else if (ch.Bonuses[BONUS_SOFT]) != 0 {
				multi = 1
			} else {
				multi = 2
			}
		} else if (ch.Bonuses[BONUS_SOFT]) != 0 {
			multi = 1
		} else {
			multi = 2
		}
	} else if loc == 1 {
		if (ch.Bonuses[BONUS_SOFT]) != 0 {
			multi = 0.25
		} else {
			multi = 0.5
		}
	} else {
		if (ch.Bonuses[BONUS_SOFT]) != 0 {
			multi = 0.5
		}
	}
	return multi
}
func roll_hitloc(ch *char_data, vict *char_data, skill int) int {
	var (
		location int = 4
		critmax  int = 1000
		critical int = 0
	)
	if IS_NPC(ch) {
		if GET_LEVEL(ch) > 100 {
			skill = rand_number(GET_LEVEL(ch), GET_LEVEL(ch)+10)
		} else {
			skill = rand_number(GET_LEVEL(ch), 100)
		}
	}
	if ch.Chclass == CLASS_DABURA && !IS_NPC(ch) {
		if (ch.Skills[SKILL_STYLE]) >= 75 {
			critmax -= 200
		}
	}
	critical = rand_number(80, critmax)
	if skill >= critical {
		location = 2
	} else if skill >= rand_number(50, 750) {
		location = 2
	} else if skill >= rand_number(50, 350) {
		location = 1
	} else if skill >= rand_number(30, 200) {
		location = 3
	} else {
		location = rand_number(4, 5)
	}
	if !IS_NPC(vict) {
		if location == 4 && (vict.Limb_condition[0]) <= 0 && (vict.Limb_condition[1]) <= 0 {
			location = 5
		}
		if location == 5 && (vict.Limb_condition[2]) <= 0 && (vict.Limb_condition[3]) <= 0 {
			location = 4
		}
		if location == 4 && (vict.Limb_condition[0]) <= 0 && (vict.Limb_condition[1]) <= 0 {
			location = 1
		}
	}
	if IS_NPC(vict) {
		if location == 4 && !MOB_FLAGGED(vict, MOB_RARM) && !MOB_FLAGGED(vict, MOB_LARM) {
			location = 5
		}
		if location == 5 && !MOB_FLAGGED(vict, MOB_RLEG) && !MOB_FLAGGED(vict, MOB_LLEG) {
			location = 4
		}
		if location == 5 && !MOB_FLAGGED(vict, MOB_RARM) && !MOB_FLAGGED(vict, MOB_LARM) {
			location = 1
		}
	}
	return location
}
func armor_calc(ch *char_data, dmg int64, type_ int) int64 {
	if IS_NPC(ch) {
		return 0
	}
	var reduce int64 = 0
	if ch.Armor < 1000 {
		reduce = int64(float64(ch.Armor) * 0.5)
	} else if ch.Armor < 2000 {
		reduce = int64(float64(ch.Armor) * 0.75)
	} else if ch.Armor < 5000 {
		reduce = int64(ch.Armor)
	} else if ch.Armor < 10000 {
		reduce = int64(ch.Armor * 2)
	} else if ch.Armor < 20000 {
		reduce = int64(ch.Armor * 4)
	} else if ch.Armor < 30000 {
		reduce = int64(ch.Armor * 8)
	} else if ch.Armor < 40000 {
		reduce = int64(ch.Armor * 12)
	} else if ch.Armor < 60000 {
		reduce = int64(ch.Armor * 25)
	} else if ch.Armor < 100000 {
		reduce = int64(ch.Armor * 50)
	} else if ch.Armor < 150000 {
		reduce = int64(ch.Armor * 75)
	} else if ch.Armor < 200000 {
		reduce = int64(ch.Armor * 150)
	} else if ch.Armor >= 200000 {
		reduce = int64(ch.Armor * 250)
	}
	var i int
	var loc int = -1
	var bonus float64 = 0.0
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			var obj *obj_data = (ch.Equipment[i])
			switch obj.Value[VAL_ALL_MATERIAL] {
			case MATERIAL_STEEL:
				loc = 0
				bonus = 0.05
			case MATERIAL_IRON:
				loc = 0
				bonus = 0.025
			case MATERIAL_COPPER:
				fallthrough
			case MATERIAL_BRASS:
				fallthrough
			case MATERIAL_METAL:
				loc = 0
				bonus = 0.01
			case MATERIAL_SILVER:
				loc = 1
				bonus = 0.1
			case MATERIAL_KACHIN:
				loc = 2
				bonus = 0.2
			case MATERIAL_CRYSTAL:
				loc = 1
				bonus = 0.05
			case MATERIAL_DIAMOND:
				loc = 2
				bonus = 0.05
			case MATERIAL_PAPER:
				fallthrough
			case MATERIAL_COTTON:
				fallthrough
			case MATERIAL_SATIN:
				fallthrough
			case MATERIAL_SILK:
				fallthrough
			case MATERIAL_BURLAP:
				fallthrough
			case MATERIAL_VELVET:
				fallthrough
			case MATERIAL_HEMP:
				fallthrough
			case MATERIAL_WAX:
				loc = 2
				bonus = -0.05
			default:
			}
		}
	}
	if bonus > 0.95 {
		bonus = 0.95
	}
	if loc != -1 {
		switch type_ {
		case 0:
			if loc == 0 || loc == 2 {
				reduce += int64(float64(reduce) * bonus)
			}
		case 1:
			if loc == 1 || loc == 2 {
				reduce += int64(float64(reduce) * bonus)
				reduce /= 2
			}
		}
	}
	return reduce
}
func chance_to_hit(ch *char_data) int {
	var num int = axion_dice(0)
	if IS_NPC(ch) {
		return num
	}
	if (ch.Player_specials.Conditions[DRUNK]) > 4 {
		num += int(ch.Player_specials.Conditions[DRUNK])
	}
	return num
}
func handle_speed(ch *char_data, vict *char_data) int {
	if ch == nil || vict == nil {
		return 0
	}
	if GET_SPEEDI(ch) > GET_SPEEDI(vict)*4 {
		return 15
	} else if GET_SPEEDI(ch) > GET_SPEEDI(vict)*2 {
		return 10
	} else if GET_SPEEDI(ch) > GET_SPEEDI(vict) {
		return 5
	} else if GET_SPEEDI(ch)*4 < GET_SPEEDI(vict) {
		return -15
	} else if GET_SPEEDI(ch)*2 < GET_SPEEDI(vict) {
		return -10
	} else if GET_SPEEDI(ch) < GET_SPEEDI(vict) {
		return -5
	}
	return 0
}
func hurt_limb(ch *char_data, vict *char_data, chance int, area int, power int64) {
	if vict == nil || IS_NPC(vict) {
		return
	}
	var dmg int = 0
	if chance > axion_dice(100) {
		return
	}
	if float64(power) > float64(gear_pl(vict))*0.5 {
		dmg = rand_number(25, 40)
	} else if float64(power) > float64(gear_pl(vict))*0.25 {
		dmg = rand_number(15, 24)
	} else if float64(power) > float64(gear_pl(vict))*0.1 {
		dmg = rand_number(8, 14)
	} else if float64(power) > float64(gear_pl(vict))*0.05 {
		dmg = rand_number(4, 7)
	} else {
		dmg = rand_number(1, 3)
	}
	if vict.Armor > 50000 {
		dmg -= 5
	} else if vict.Armor > 40000 {
		dmg -= 4
	} else if vict.Armor > 30000 {
		dmg -= 3
	} else if vict.Armor > 20000 {
		dmg -= 2
	} else if vict.Armor > 10000 {
		dmg -= 1
	} else if vict.Armor > 5000 {
		dmg -= rand_number(0, 1)
	}
	if dmg <= 0 {
		return
	}
	if !is_sparring(ch) {
		if area == 0 {
			if (vict.Limb_condition[1])-dmg <= 0 {
				act(libc.CString("@RYour attack @YDESTROYS @r$N's@R left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack @YDESTROYS@R YOUR left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack @YDESTROYS @r$N's@R left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Limb_condition[1] = 0
				if PLR_FLAGGED(vict, PLR_THANDW) {
					vict.Act[int(PLR_THANDW/32)] &= bitvector_t(^(1 << (int(PLR_THANDW % 32))))
				}
				if PLR_FLAGGED(vict, PLR_CLARM) {
					vict.Act[int(PLR_CLARM/32)] &= bitvector_t(^(1 << (int(PLR_CLARM % 32))))
				}
				remove_limb(vict, 2)
			} else if (vict.Limb_condition[1]) > 0 {
				vict.Limb_condition[1] -= dmg
				act(libc.CString("@RYour attack hurts @r$N's@R left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack hurts YOUR left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack hurts @r$N's@R left arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			} else if (vict.Limb_condition[0])-dmg <= 0 {
				act(libc.CString("@RYour attack @YDESTROYS @r$N's@R right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack @YDESTROYS@R YOUR right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack @YDESTROYS @r$N's@R right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Limb_condition[0] = 0
				if PLR_FLAGGED(vict, PLR_THANDW) {
					vict.Act[int(PLR_THANDW/32)] &= bitvector_t(^(1 << (int(PLR_THANDW % 32))))
				}
				if PLR_FLAGGED(vict, PLR_CLARM) {
					vict.Act[int(PLR_CRARM/32)] &= bitvector_t(^(1 << (int(PLR_CRARM % 32))))
				}
				remove_limb(vict, 2)
			} else if (vict.Limb_condition[0]) > 0 {
				vict.Limb_condition[0] -= dmg
				act(libc.CString("@RYour attack hurts @r$N's@R right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack hurts YOUR right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack hurts @r$N's@R right arm!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else if area == 1 {
			if (vict.Limb_condition[3])-dmg <= 0 {
				act(libc.CString("@RYour attack @YDESTROYS @r$N's@R left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack @YDESTROYS@R YOUR left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack @YDESTROYS @r$N's@R left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Limb_condition[3] = 0
				if PLR_FLAGGED(vict, PLR_THANDW) {
					vict.Act[int(PLR_THANDW/32)] &= bitvector_t(^(1 << (int(PLR_THANDW % 32))))
				}
				if PLR_FLAGGED(vict, PLR_CLLEG) {
					vict.Act[int(PLR_CLLEG/32)] &= bitvector_t(^(1 << (int(PLR_CLLEG % 32))))
				}
				remove_limb(vict, 2)
			} else if (vict.Limb_condition[3]) > 0 {
				vict.Limb_condition[3] -= dmg
				act(libc.CString("@RYour attack hurts @r$N's@R left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack hurts YOUR left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack hurts @r$N's@R left leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			} else if (vict.Limb_condition[2])-dmg <= 0 {
				act(libc.CString("@RYour attack @YDESTROYS @r$N's@R right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack @YDESTROYS@R YOUR right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack @YDESTROYS @r$N's@R right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Limb_condition[2] = 0
				if PLR_FLAGGED(vict, PLR_THANDW) {
					vict.Act[int(PLR_THANDW/32)] &= bitvector_t(^(1 << (int(PLR_THANDW % 32))))
				}
				if PLR_FLAGGED(vict, PLR_CLLEG) {
					vict.Act[int(PLR_CLLEG/32)] &= bitvector_t(^(1 << (int(PLR_CLLEG % 32))))
				}
				remove_limb(vict, 2)
			} else if (vict.Limb_condition[2]) > 0 {
				vict.Limb_condition[2] -= dmg
				act(libc.CString("@RYour attack hurts @r$N's@R right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@r$n's@R attack hurts YOUR right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@r$n's@R attack hurts @r$N's@R right leg!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		}
	}
}
func dam_eq_loc(vict *char_data, area int) {
	var (
		location int = 0
		num      int = 0
	)
	if vict == nil || vict == nil || vict.Hit <= 0 {
		return
	}
	switch area {
	case 1:
		num = rand_number(1, 8)
		switch num {
		case 1:
			location = WEAR_FINGER_R
		case 2:
			location = WEAR_FINGER_L
		case 3:
			location = WEAR_ARMS
		case 4:
			location = WEAR_WRIST_R
		case 5:
			location = WEAR_WRIST_L
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			location = WEAR_HANDS
		}
	case 2:
		num = rand_number(1, 3)
		switch num {
		case 1:
			location = WEAR_LEGS
		case 2:
			location = WEAR_FEET
		case 3:
			location = WEAR_WAIST
		}
	case 3:
		num = rand_number(1, 6)
		switch num {
		case 1:
			location = WEAR_HEAD
		case 2:
			location = WEAR_NECK_1
		case 3:
			location = WEAR_NECK_2
		case 4:
			location = WEAR_EAR_R
		case 5:
			location = WEAR_EAR_L
		case 6:
			location = WEAR_EYE
		}
	case 4:
		num = rand_number(1, 4)
		switch num {
		case 1:
			location = WEAR_BODY
		case 2:
			location = WEAR_ABOUT
		case 3:
			location = WEAR_BACKPACK
		case 4:
			location = WEAR_SH
		}
	default:
		location = WEAR_BODY
	}
	damage_eq(vict, location)
}
func damage_eq(vict *char_data, location int) {
	if (vict.Equipment[location]) != nil && rand_number(1, 20) >= 19 && !AFF_FLAGGED(vict, AFF_SANCTUARY) {
		var eq *obj_data = (vict.Equipment[location])
		if OBJ_FLAGGED(eq, ITEM_UNBREAKABLE) {
			return
		}
		var loss int = rand_number(2, 5)
		if GET_OBJ_VNUM(eq) == 0x4E83 || GET_OBJ_VNUM(eq) == 0x4E82 {
			loss = 1
		}
		if AFF_FLAGGED(vict, AFF_CURSE) {
			loss *= 3
		} else if AFF_FLAGGED(vict, AFF_BLESS) && rand_number(1, 3) == 3 {
			loss = 1
		} else if AFF_FLAGGED(vict, AFF_BLESS) {
			return
		}
		eq.Value[VAL_ALL_HEALTH] -= loss
		if (eq.Value[VAL_ALL_HEALTH]) <= 0 {
			eq.Value[VAL_ALL_HEALTH] = 0
			eq.Extra_flags[int(ITEM_BROKEN/32)] |= bitvector_t(1 << (int(ITEM_BROKEN % 32)))
			act(libc.CString("@WYour $p@W completely breaks!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$N's@W $p@W completely breaks!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_NOTVICT)
			perform_remove(vict, location)
			if !IS_NPC(vict) {
				save_char(vict)
			}
		} else if (eq.Value[VAL_ALL_MATERIAL]) == MATERIAL_LEATHER || (eq.Value[VAL_ALL_MATERIAL]) == MATERIAL_COTTON || (eq.Value[VAL_ALL_MATERIAL]) == MATERIAL_SILK {
			act(libc.CString("@WYour $p@W rips a little!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$N's@W $p@W rips a little!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_NOTVICT)
			if AFF_FLAGGED(vict, AFF_BLESS) {
				send_to_char(vict, libc.CString("@c...But your blessing seems to have partly mended this damage.@n\r\n"))
				act(libc.CString("@c...but @C$N's@c body glows blue for a moment and the damage mends a little.@n"), TRUE, nil, nil, unsafe.Pointer(vict), TO_NOTVICT)
			} else if AFF_FLAGGED(vict, AFF_CURSE) {
				send_to_char(vict, libc.CString("@r...and your curse seems to have made the damage three times worse!@n\r\n"))
				act(libc.CString("@c...but @C$N's@c body glows red for a moment and the damage grow three times worse!@n"), TRUE, nil, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else {
			act(libc.CString("@WYour $p@W cracks a little!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@C$N's@W $p@W cracks a little!@n"), FALSE, nil, eq, unsafe.Pointer(vict), TO_NOTVICT)
			if AFF_FLAGGED(vict, AFF_BLESS) {
				send_to_char(vict, libc.CString("@c...But your blessing seems to have partly mended this damage.@n\r\n"))
				act(libc.CString("@c...but @C$N's@c body glows blue for a moment and the damage mends a little.@n"), TRUE, nil, nil, unsafe.Pointer(vict), TO_NOTVICT)
			} else if AFF_FLAGGED(vict, AFF_CURSE) {
				send_to_char(vict, libc.CString("@r...and your curse seems to have made the damage three times worse!@n\r\n"))
				act(libc.CString("@c...but @C$N's@c body glows red for a moment and the damage grow three times worse!@n"), TRUE, nil, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		}
	}
}
func update_mob_absorb() {
	var (
		roll int = 0
		i    *char_data
		vict *char_data
	)
	for i = character_list; i != nil; i = i.Next {
		roll = int(float64(axion_dice(0)) + float64(GET_LEVEL(i))*0.25)
		if !IS_NPC(i) {
			continue
		}
		if i.Race != RACE_ANDROID {
			continue
		}
		if !MOB_FLAGGED(i, MOB_ABSORB) {
			continue
		}
		if i.Absorbing == nil || i.Absorbing == nil {
			continue
		} else if GET_LEVEL(i) < roll {
			continue
		} else if i.Absorbing != nil {
			vict = i.Absorbing
			var ki int = int(float64(vict.Max_mana) * 0.01)
			var stam int = int(float64(vict.Max_move) * 0.01)
			var pl int = int(float64(vict.Max_hit) * 0.01)
			var maxed int = 0
			if float64(roll) < float64(GET_LEVEL(i)+1)*0.75 {
				ki += ki * rand_number(2, 4)
				stam += stam * rand_number(2, 4)
				pl += pl * rand_number(2, 4)
			}
			vict.Hit -= int64(pl)
			i.Hit += int64(pl)
			vict.Mana -= int64(ki)
			i.Mana += int64(ki)
			vict.Move -= int64(stam)
			i.Move += int64(stam)
			if vict.Mana < 0 {
				vict.Mana = 0
			}
			if vict.Move < 0 {
				vict.Move = 0
			}
			if i.Hit > i.Max_hit {
				i.Hit = i.Max_hit
				maxed += 1
			}
			if i.Move > i.Max_move {
				i.Move = i.Max_move
				maxed += 1
			}
			if i.Mana > i.Max_mana {
				i.Mana = i.Max_mana
				maxed += 1
			}
			if vict.Hit <= 0 {
				act(libc.CString("@R$n@r absorbs the last of YOUR energy and you die...@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$n@r absorbs the last of @R$N's@r energy and $E dies...@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_NOTVICT)
				die(vict, i)
			} else if maxed >= 3 {
				act(libc.CString("@R$n@r absorbs some of YOUR energy...but $e seems to be full now and releases YOU!@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$n@r absorbs some of @R$N's@r energy...but $e seems to be full now and lets go.@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_NOTVICT)
			} else {
				act(libc.CString("@R$n@r absorbs some of YOUR energy!@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$n@r absorbs some of @R$N's@r energy.@n"), TRUE, i, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		}
	}
}
func huge_update() {
	var (
		dge    int   = 0
		skill  int   = 0
		bonus  int   = 1
		count  int   = 0
		dmg    int64 = 0
		k      *obj_data
		ch     *char_data
		vict   *char_data
		next_v *char_data
	)
	for k = object_list; k != nil; k = k.Next {
		if k.Aucter > 0 && k.AucTime+604800 <= C.time(nil) {
			if k.In_room != 0 && (func() room_vnum {
				if k.In_room != room_rnum(-1) && k.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Number
				}
				return -1
			}()) == 80 {
				var inroom room_vnum = room_vnum(k.In_room)
				(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(inroom)))).Room_flags[int(ROOM_HOUSE_CRASH/32)] &= bitvector_t(^(1 << (int(ROOM_HOUSE_CRASH % 32))))
				extract_obj(k)
				continue
			}
		}
		if k.Kicharge <= 0 {
			continue
		}
		if GET_OBJ_VNUM(k) != 82 && GET_OBJ_VNUM(k) != 83 {
			continue
		} else if k.Distance <= 0 {
			if k.Kitype == 497 {
				if k.Target.In_room == k.In_room {
					ch = k.User
					if (func() room_vnum {
						if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
						}
						return -1
					}()) == (func() room_vnum {
						if k.In_room != room_rnum(-1) && k.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Number
						}
						return -1
					}()) {
						bonus = 2
					}
					act(libc.CString("@WThe large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W descends on YOU! It eclipses everything above you as it crushes down into you! You struggle against it with all your might!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
					act(libc.CString("@WThe large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W descends on @C$n@W! It completely obscures $m from view as it crushes into $s body! It appears to be facing some resistance from $m!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
					send_to_room(k.In_room, libc.CString("\r\n"))
					if k.Target.Hit*int64(bonus) < k.Kicharge*5 {
						act(libc.CString("@WYour strength is no match for the power of the attack! It slowly grinds into you before exploding into a massive blast!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@W's strength is no match for the power of the attack! It slowly grinds into $m before exploding into a massive blast!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
						skill = init_skill(ch, SKILL_GENKIDAMA)
						dmg = int64(float64(k.Kicharge) * 1.25)
						hurt(0, 0, ch, k.Target, nil, dmg, 1)
						dmg /= 2
						for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).People; vict != nil; vict = next_v {
							next_v = vict.Next_in_room
							if vict == ch {
								continue
							}
							if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
								continue
							}
							if vict == k.Target {
								continue
							}
							if AFF_FLAGGED(vict, AFF_GROUP) {
								if vict.Master == ch {
									continue
								} else if ch.Master == vict {
									continue
								} else if vict.Master == ch.Master {
									continue
								}
							}
							if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
								continue
							}
							if MOB_FLAGGED(vict, MOB_NOKILL) {
								continue
							}
							dge = handle_dodge(vict)
							if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
								act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@cYou disappear, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
								act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
								vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
								pcost(vict, 0, vict.Max_hit/200)
								hurt(0, 0, ch, vict, nil, 0, 1)
								continue
							} else if dge+rand_number(-10, 5) > skill {
								act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@WYou manage to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
								act(libc.CString("@c$N@W manages to escape the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
								hurt(0, 0, ch, vict, nil, 0, 1)
								improve_skill(vict, SKILL_DODGE, 0)
								continue
							} else {
								count += 1
								if IS_NPC(vict) && count > 10 {
									if vict.Hit < dmg {
										var loss float64 = 0.0
										if count >= 30 {
											loss = 0.8
										} else if count >= 20 {
											loss = 0.6
										} else if count >= 15 {
											loss = 0.4
										} else if count >= 10 {
											loss = 0.25
										}
										vict.Exp -= int64(float64(vict.Exp) * loss)
									}
								}
								act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@RYou are caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
								act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
								hurt(0, 0, ch, vict, nil, dmg, 1)
								continue
							}
						}
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Dmg = 100
						var zone int = 0
						if (func() int {
							zone = int(real_zone_by_thing(func() room_vnum {
								if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
									return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
								}
								return -1
							}()))
							return zone
						}()) != int(-1) {
							send_to_zone(libc.CString("A MASSIVE explosion shakes the entire area!\r\n"), zone_rnum(zone))
						}
						extract_obj(k)
						continue
					} else {
						act(libc.CString("@WYou manage to overpower the attack! You lift up into the sky slowly with it and toss it up and away out of sight!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@W manages to unbelievably overpower the attack! It is lifted up into the sky and tossed away dramaticly!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
						hurt(0, 0, ch, k.Target, nil, 0, 1)
						k.Target.Move -= k.Kicharge / 4
						extract_obj(k)
						continue
					}
				} else if k.Target.In_room != k.In_room {
					ch = k.User
					send_to_room(k.In_room, libc.CString("@WThe large @cS@Cp@wi@cr@Ci@wt @cB@Co@wm@cb@W descends on the area! It slowly burns into the ground before exploding magnificently!@n\r\n"))
					skill = init_skill(ch, SKILL_GENKIDAMA)
					dmg = k.Kicharge
					dmg /= 2
					for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).People; vict != nil; vict = next_v {
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
						if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
							act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("@cYou disappear, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
							vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
							pcost(vict, 0, vict.Max_hit/200)
							hurt(0, 0, ch, vict, nil, 0, 1)
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
							hurt(0, 0, ch, vict, nil, dmg, 1)
							continue
						}
					}
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Dmg = 100
					var zone int = 0
					if (func() int {
						zone = int(real_zone_by_thing(func() room_vnum {
							if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
							}
							return -1
						}()))
						return zone
					}()) != int(-1) {
						send_to_zone(libc.CString("A MASSIVE explosion shakes the entire area!\r\n"), zone_rnum(zone))
					}
					extract_obj(k)
					continue
				}
				continue
			}
			if k.Kitype == 498 {
				if k.Target.In_room == k.In_room {
					ch = k.User
					if (func() room_vnum {
						if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
						}
						return -1
					}()) == (func() room_vnum {
						if k.In_room != room_rnum(-1) && k.In_room <= top_of_world {
							return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Number
						}
						return -1
					}()) {
						bonus = 2
					}
					act(libc.CString("@WThe large @mG@Me@wn@mo@Mc@wi@md@Me@W descends on YOU! It eclipses everything above you as it crushes down into you! You struggle against it with all your might!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
					act(libc.CString("@WThe large @mG@Me@wn@mo@Mc@wi@md@Me@W descends on @C$n@W! It completely obscures $m from view as it crushes into $s body! It appears to be facing some resistance from $m!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
					send_to_room(k.In_room, libc.CString("\r\n"))
					if k.Target.Hit*int64(bonus) < k.Kicharge*10 {
						act(libc.CString("@WYour strength is no match for the power of the attack! It slowly grinds into you before exploding into a massive blast!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@W's strength is no match for the power of the attack! It slowly grinds into $m before exploding into a massive blast!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
						skill = init_skill(ch, SKILL_GENOCIDE)
						dmg = k.Kicharge
						hurt(0, 0, ch, k.Target, nil, dmg, 1)
						dmg /= 2
						for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).People; vict != nil; vict = next_v {
							next_v = vict.Next_in_room
							if vict == ch {
								continue
							}
							if AFF_FLAGGED(vict, AFF_SPIRIT) && !IS_NPC(vict) {
								continue
							}
							if vict == k.Target {
								continue
							}
							if AFF_FLAGGED(vict, AFF_GROUP) {
								if vict.Master == ch {
									continue
								} else if ch.Master == vict {
									continue
								} else if vict.Master == ch.Master {
									continue
								}
							}
							if GET_LEVEL(vict) <= 8 && !IS_NPC(vict) {
								continue
							}
							if MOB_FLAGGED(vict, MOB_NOKILL) {
								continue
							}
							dge = handle_dodge(vict)
							if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
								act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@cYou disappear, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
								act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
								count += 1
								if IS_NPC(vict) && count > 10 {
									if vict.Hit < dmg {
										var loss float64 = 0.0
										if count >= 30 {
											loss = 0.8
										} else if count >= 20 {
											loss = 0.6
										} else if count >= 15 {
											loss = 0.4
										} else if count >= 10 {
											loss = 0.25
										}
										vict.Exp -= int64(float64(vict.Exp) * loss)
									}
								}
								act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@RYou are caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
								act(libc.CString("@R$N@r is caught by the explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
								hurt(0, 0, ch, vict, nil, dmg, 1)
								continue
							}
						}
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Dmg = 100
						var zone int = 0
						if (func() int {
							zone = int(real_zone_by_thing(func() room_vnum {
								if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
									return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
								}
								return -1
							}()))
							return zone
						}()) != int(-1) {
							send_to_zone(libc.CString("A MASSIVE explosion shakes the entire area!\r\n"), zone_rnum(zone))
						}
						extract_obj(k)
						continue
					} else {
						act(libc.CString("@WYou manage to overpower the attack! You lift up into the sky slowly with it and toss it up and away out of sight!@n"), TRUE, k.Target, nil, nil, TO_CHAR)
						act(libc.CString("@C$n@W manages to unbelievably overpower the attack! It is lifted up into the sky and tossed away dramaticly!@n"), TRUE, k.Target, nil, nil, TO_ROOM)
						hurt(0, 0, ch, k.Target, nil, 0, 1)
						k.Target.Move -= k.Kicharge / 4
						extract_obj(k)
						continue
					}
				} else if k.Target.In_room != k.In_room {
					ch = k.User
					send_to_room(k.In_room, libc.CString("@WThe large @mG@Me@wn@mo@Mc@wi@md@Me@W descends on the area! It slowly burns into the ground before exploding magnificantly!@n\r\n"))
					skill = init_skill(ch, SKILL_GENOCIDE)
					dmg = k.Kicharge
					dmg /= 2
					for vict = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).People; vict != nil; vict = next_v {
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
						if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
							act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("@cYou disappear, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@C$N@c disappears, avoiding the explosion!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
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
							hurt(0, 0, ch, vict, nil, dmg, 1)
							continue
						}
					}
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(k.In_room)))).Dmg = 100
					var zone int = 0
					if (func() int {
						zone = int(real_zone_by_thing(func() room_vnum {
							if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
								return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
							}
							return -1
						}()))
						return zone
					}()) != int(-1) {
						send_to_zone(libc.CString("A MASSIVE explosion shakes the entire area!\r\n"), zone_rnum(zone))
					}
					extract_obj(k)
					continue
				}
				continue
			} else {
				extract_obj(k)
				continue
			}
		}
		act(libc.CString("$p@W descends slowly towards the ground!@n"), TRUE, nil, k, nil, TO_ROOM)
		k.Distance--
	}
}
func homing_update() {
	var k *obj_data
	for k = object_list; k != nil; k = k.Next {
		if k == nil || k == nil {
			continue
		}
		if k.Kicharge <= 0 {
			continue
		}
		if GET_OBJ_VNUM(k) != 80 && GET_OBJ_VNUM(k) != 81 && GET_OBJ_VNUM(k) != 84 {
			continue
		} else if k.Target != nil && k.User != nil {
			var (
				ch   *char_data = k.User
				vict *char_data = k.Target
			)
			if GET_OBJ_VNUM(k) == 80 {
				if k.In_room != vict.In_room {
					act(libc.CString("@wThe $p@w pursues after you!@n"), TRUE, vict, k, nil, TO_CHAR)
					act(libc.CString("@wThe $p@W pursues after @C$n@w!@n"), TRUE, vict, k, nil, TO_ROOM)
					obj_from_room(k)
					obj_to_room(k, vict.In_room)
					continue
				} else {
					act(libc.CString("@RThe $p@R makes a tight turn and rockets straight for you!@n"), TRUE, vict, k, nil, TO_CHAR)
					act(libc.CString("@RThe $p@R makes a tight turn and rockets straight for @r$n@n"), TRUE, vict, k, nil, TO_ROOM)
					if handle_parry(vict) < rand_number(1, 140) {
						act(libc.CString("@rThe $p@r slams into your body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_CHAR)
						act(libc.CString("@rThe $p@r slams into @R$n's@r body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_ROOM)
						var dmg int64 = k.Kicharge
						extract_obj(k)
						hurt(0, 0, ch, vict, nil, dmg, 1)
						continue
					} else if rand_number(1, 3) > 1 {
						act(libc.CString("@wYou manage to deflect the $p@W sending it flying away and depleting some of its energy.@n"), TRUE, vict, k, nil, TO_CHAR)
						act(libc.CString("@C$n @wmanages to deflect the $p@w sending it flying away and depleting some of its energy.@n"), TRUE, vict, k, nil, TO_ROOM)
						k.Kicharge -= k.Kicharge / 10
						if k.Kicharge <= 0 {
							send_to_room(k.In_room, libc.CString("%s has lost all its energy and disappears.\r\n"), k.Short_description)
							extract_obj(k)
						}
						continue
					} else {
						act(libc.CString("@wYou manage to deflect the $p@w sending it flying away into the nearby surroundings!@n"), TRUE, vict, k, nil, TO_CHAR)
						act(libc.CString("@C$n @wmanages to deflect the $p@w sending it flying away into the nearby surroundings!@n"), TRUE, vict, k, nil, TO_ROOM)
						if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg <= 95 {
							(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg += 5
						}
						extract_obj(k)
						continue
					}
				}
				continue
			} else if GET_OBJ_VNUM(k) == 81 || GET_OBJ_VNUM(k) == 84 {
				if k.In_room != vict.In_room {
					act(libc.CString("@wYou lose sight of @C$N@W and let $p@W fly away!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@wYou manage to escape @C$n's@W $p@W!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$n@W loses sight of @c$N@W and lets $s $p@W fly away!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_NOTVICT)
					extract_obj(k)
					continue
				} else {
					act(libc.CString("@RYou move your hand and direct $p@R after @r$N@R!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@r$n@R moves $s hand and directs $p@R after YOU!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@r$n@R moves $s hand and directs $p@R after @r$N@R!@n"), TRUE, ch, k, unsafe.Pointer(vict), TO_NOTVICT)
					if handle_parry(vict) < rand_number(1, 140) {
						if GET_OBJ_VNUM(k) != 84 {
							act(libc.CString("@rThe $p@r slams into your body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_CHAR)
							act(libc.CString("@rThe $p@r slams into @R$n's@r body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_ROOM)
							var dmg int64 = k.Kicharge
							extract_obj(k)
							hurt(0, 0, ch, vict, nil, dmg, 1)
						} else if GET_OBJ_VNUM(k) == 84 {
							var dmg int64 = k.Kicharge
							if dmg > vict.Max_hit/5 && (vict.Race != RACE_MAJIN && vict.Race != RACE_BIO) {
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
							} else if dmg > vict.Max_hit/5 && (vict.Race == RACE_MAJIN || vict.Race == RACE_BIO) {
								if GET_SKILL(vict, SKILL_REGENERATE) > rand_number(1, 101) && vict.Mana >= vict.Max_mana/40 {
									act(libc.CString("@R$N@r is cut in half by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
									act(libc.CString("@rYou are cut in half by the attack but regenerate a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
									act(libc.CString("@R$N@r is cut in half by the attack but regenerates a moment later!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
									vict.Mana -= vict.Max_mana / 40
									hurt(0, 0, ch, vict, nil, dmg, 1)
								} else if dmg > vict.Max_hit/5 && (vict.Race == RACE_MAJIN || vict.Race == RACE_BIO) {
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
								}
							} else {
								act(libc.CString("@rThe $p@r slams into your body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_CHAR)
								act(libc.CString("@rThe $p@r slams into @R$n's@r body, exploding in a flash of bright light!@n"), TRUE, vict, k, nil, TO_ROOM)
								hurt(0, 0, ch, vict, nil, dmg, 1)
							}
							extract_obj(k)
						}
						continue
					} else if rand_number(1, 3) > 1 {
						act(libc.CString("@wYou manage to deflect the $p@W sending it flying away and depleting some of its energy.@n"), TRUE, vict, k, nil, TO_CHAR)
						act(libc.CString("@C$n @wmanages to deflect the $p@w sending it flying away and depleting some of its energy.@n"), TRUE, vict, k, nil, TO_ROOM)
						k.Kicharge -= k.Kicharge / 10
						if k.Kicharge <= 0 {
							send_to_room(k.In_room, libc.CString("%s has lost all its energy and disappears.\r\n"), k.Short_description)
							extract_obj(k)
						}
						continue
					} else {
						act(libc.CString("@wYou manage to deflect the $p@w sending it flying away into the nearby surroundings!@n"), TRUE, vict, k, nil, TO_CHAR)
						act(libc.CString("@C$n @wmanages to deflect the $p@w sending it flying away into the nearby surroundings!@n"), TRUE, vict, k, nil, TO_ROOM)
						if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg <= 95 {
							(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Dmg += 5
						}
						extract_obj(k)
						continue
					}
				}
			}
		}
	}
}
func limb_ok(ch *char_data, type_ int) int {
	if IS_NPC(ch) {
		if AFF_FLAGGED(ch, AFF_ENSNARED) && rand_number(1, 100) <= 90 {
			return FALSE
		}
		return TRUE
	}
	if ch.Grappling != nil && ch.Grap != 3 {
		send_to_char(ch, libc.CString("You are too busy grappling!\r\n"))
		return FALSE
	}
	if ch.Grappled != nil && (ch.Grap == 1 || ch.Grap == 4) {
		send_to_char(ch, libc.CString("You are unable to move while in this hold! Try using 'escape' to get out of it!\r\n"))
		return FALSE
	}
	if ch.Powerattack > 0 {
		send_to_char(ch, libc.CString("You are currently playing a song! Enter the song command in order to stop!\r\n"))
		return FALSE
	}
	if type_ == 0 {
		if !HAS_ARMS(ch) {
			send_to_char(ch, libc.CString("You have no available arms!\r\n"))
			return FALSE
		}
		if AFF_FLAGGED(ch, AFF_ENSNARED) && rand_number(1, 100) <= 90 {
			send_to_char(ch, libc.CString("You are unable to move your arms while bound by this strong silk!\r\n"))
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
			return FALSE
		} else if AFF_FLAGGED(ch, AFF_ENSNARED) {
			act(libc.CString("You manage to break the silk ensnaring your arms!"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("$n manages to break the silk ensnaring $s arms!"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Affected_by[int(AFF_ENSNARED/32)] &= ^(1 << (int(AFF_ENSNARED % 32)))
		}
		if (ch.Equipment[WEAR_WIELD1]) != nil && (ch.Equipment[WEAR_WIELD2]) != nil {
			send_to_char(ch, libc.CString("Your hands are full!\r\n"))
			return FALSE
		}
	} else if type_ > 0 {
		if !HAS_LEGS(ch) {
			send_to_char(ch, libc.CString("You have no working legs!\r\n"))
			return FALSE
		}
	}
	return TRUE
}
func init_skill(ch *char_data, snum int) int {
	var skill int = 0
	if !IS_NPC(ch) {
		skill = GET_SKILL(ch, snum)
		if PLR_FLAGGED(ch, PLR_TRANSMISSION) {
			skill += 4
		}
		if skill > 118 {
			skill = 118
		}
		return skill
	}
	if IS_NPC(ch) && GET_LEVEL(ch) <= 10 {
		skill = rand_number(30, 50)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 20 {
		skill = rand_number(45, 65)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 30 {
		skill = rand_number(55, 70)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 50 {
		skill = rand_number(65, 80)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 70 {
		skill = rand_number(75, 90)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 80 {
		skill = rand_number(85, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 90 {
		skill = rand_number(90, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 100 {
		skill = rand_number(95, 100)
	} else if IS_NPC(ch) && GET_LEVEL(ch) <= 110 {
		skill = rand_number(95, 105)
	} else {
		skill = rand_number(100, 110)
	}
	return skill
}
func handle_block(ch *char_data) int {
	if axion_dice(0) <= 4 {
		return 1
	}
	if !IS_NPC(ch) {
		if GET_SKILL(ch, SKILL_BLOCK) == 0 {
			return 0
		} else {
			var num int = GET_SKILL(ch, SKILL_BLOCK)
			if ch.Race == RACE_MUTANT && ((ch.Genome[0]) == 3 || (ch.Genome[1]) == 3) {
				num += 10
			}
			if (ch.Skills[SKILL_STYLE]) >= 100 {
				num += 5
			} else if (ch.Skills[SKILL_STYLE]) >= 80 {
				num += 4
			} else if (ch.Skills[SKILL_STYLE]) >= 60 {
				num += 3
			} else if (ch.Skills[SKILL_STYLE]) >= 40 {
				num += 2
			} else if (ch.Skills[SKILL_STYLE]) >= 20 {
				num += 1
			}
			return num
		}
	} else {
		if !IS_HUMANOID(ch) {
			var top int = GET_LEVEL(ch) / 4
			if top < 5 {
				top = 6
			}
			return rand_number(5, top)
		} else {
			if GET_LEVEL(ch) >= 110 {
				return rand_number(95, 105)
			} else if GET_LEVEL(ch) >= 100 {
				return rand_number(85, 95)
			} else if GET_LEVEL(ch) >= 90 {
				return rand_number(70, 85)
			} else if GET_LEVEL(ch) >= 75 {
				return rand_number(50, 70)
			} else if GET_LEVEL(ch) >= 40 {
				return rand_number(40, 50)
			} else {
				var top int = GET_LEVEL(ch)
				if top < 15 {
					top = 16
				}
				return rand_number(15, top)
			}
		}
	}
}
func handle_dodge(ch *char_data) int {
	if axion_dice(0) <= 4 {
		return 1
	}
	if !IS_NPC(ch) {
		if GET_SKILL(ch, SKILL_DODGE) == 0 {
			return 0
		} else {
			var num int = GET_SKILL(ch, SKILL_DODGE)
			if ch.Race == RACE_MUTANT && ((ch.Genome[0]) == 3 || (ch.Genome[1]) == 3) {
				num += 10
			}
			if (ch.Skills[SKILL_STYLE]) >= 100 {
				num += 5
			} else if (ch.Skills[SKILL_STYLE]) >= 80 {
				num += 4
			} else if (ch.Skills[SKILL_STYLE]) >= 60 {
				num += 3
			} else if (ch.Skills[SKILL_STYLE]) >= 40 {
				num += 2
			} else if (ch.Skills[SKILL_STYLE]) >= 20 {
				num += 1
			}
			if (ch.Skills[SKILL_SURVIVAL]) >= 100 {
				num += 3
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 75 {
				num += 2
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 50 {
				num += 1
			}
			if (ch.Skills[SKILL_ROLL]) >= 100 {
				num += 5
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 80 {
				num += 4
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 60 {
				num += 3
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 40 {
				num += 2
			} else if (ch.Skills[SKILL_SURVIVAL]) >= 20 {
				num += 1
			}
			if group_bonus(ch, 2) == 8 {
				num += int(float64(num) * 0.05)
			}
			return num
		}
	} else {
		if !IS_HUMANOID(ch) {
			var top int = (GET_LEVEL(ch) + 1) / 8
			if top < 5 {
				top = 6
			}
			return rand_number(5, top)
		} else {
			if GET_LEVEL(ch) >= 110 {
				return rand_number(95, 105)
			} else if GET_LEVEL(ch) >= 100 {
				return rand_number(75, 95)
			} else if GET_LEVEL(ch) >= 90 {
				return rand_number(50, 85)
			} else if GET_LEVEL(ch) >= 75 {
				return rand_number(30, 70)
			} else if GET_LEVEL(ch) >= 40 {
				return rand_number(20, 50)
			} else {
				var top int = GET_LEVEL(ch)
				if top < 15 {
					top = 16
				}
				return rand_number(15, top)
			}
		}
	}
}
func check_def(vict *char_data) int {
	var (
		index int = 0
		pry   int = handle_parry(vict)
		dge   int = handle_dodge(vict)
		blk   int = handle_block(vict)
	)
	index = pry + dge + blk
	if index > 0 {
		index /= 3
	}
	if AFF_FLAGGED(vict, AFF_KNOCKED) {
		index = 0
	}
	return index
}
func handle_defense(vict *char_data, pry *int, blk *int, dge *int) {
	if !IS_NPC(vict) {
		*pry = handle_parry(vict)
		*blk = handle_block(vict)
		*dge = handle_dodge(vict)
		if (vict.Bonuses[BONUS_WALL]) != 0 {
			*blk += int(float64(GET_SKILL(vict, SKILL_BLOCK)) * 0.2)
		}
		if (vict.Bonuses[BONUS_PUSHOVER]) != 0 {
			*blk -= int(float64(GET_SKILL(vict, SKILL_BLOCK)) * 0.2)
		}
		if (vict.Equipment[WEAR_WIELD1]) == nil && (vict.Equipment[WEAR_WIELD2]) == nil {
			*blk += 4
		}
		if *blk > 110 {
			*blk = 110
		}
		if (vict.Bonuses[BONUS_EVASIVE]) != 0 {
			*dge += int(float64(GET_SKILL(vict, SKILL_DODGE)) * 0.15)
		}
		if (vict.Bonuses[BONUS_PUNCHINGBAG]) != 0 {
			*dge -= int(float64(GET_SKILL(vict, SKILL_DODGE)) * 0.15)
		}
		if *dge > 110 {
			*dge = 110
		}
		if *pry > 110 {
			*pry = 110
		}
		if PLR_FLAGGED(vict, PLR_GOOP) && rand_number(1, 100) >= 15 {
			*dge += 100
			*blk += 100
			*pry += 100
		}
	} else {
		*pry = handle_parry(vict)
		*blk = handle_block(vict)
		*dge = handle_dodge(vict)
	}
	return
}
func parry_ki(attperc float64, ch *char_data, vict *char_data, sname [1000]byte, prob int, perc int, skill int, type_ int) {
	var (
		buf      [200]byte
		buf2     [200]byte
		buf3     [200]byte
		foundv   int   = FALSE
		foundo   int   = FALSE
		dmg      int64 = 0
		tob      *obj_data
		next_obj *obj_data
		tch      *char_data
		next_v   *char_data
	)
	for tch = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).People; tch != nil; tch = next_v {
		next_v = tch.Next_in_room
		if tch == ch {
			continue
		}
		if tch == vict {
			continue
		}
		if can_kill(ch, tch, nil, 1) == 0 {
			continue
		}
		if rand_number(1, 101) >= 90 && foundv == FALSE {
			if handle_parry(tch) > rand_number(1, 140) {
				stdio.Sprintf(&buf[0], "@C$N@W deflects your %s, sending it flying away!@n", &sname[0])
				stdio.Sprintf(&buf2[0], "@WYou deflect @C$n's@W %s sending it flying away!@n", &sname[0])
				stdio.Sprintf(&buf3[0], "@C$N@W deflects @c$n's@W %s sending it flying away!@n", &sname[0])
				act(&buf[0], TRUE, ch, nil, unsafe.Pointer(tch), TO_CHAR)
				act(&buf2[0], TRUE, ch, nil, unsafe.Pointer(tch), TO_VICT)
				act(&buf3[0], TRUE, ch, nil, unsafe.Pointer(tch), TO_NOTVICT)
				foundv = FALSE
			} else {
				foundv = TRUE
				stdio.Sprintf(&buf[0], "@WYou watch as the deflected %s slams into @C$N@W, exploding with a roar of blinding light!@n", &sname[0])
				stdio.Sprintf(&buf2[0], "@c$n@W watches as the deflected %s slams into you! The %s explodes with a roar of blinding light!@n", &sname[0], &sname[0])
				stdio.Sprintf(&buf3[0], "@c$n@W watches as the deflected %s slams into @C$N@W! The %s explodes with a roar of blinding light!@n", &sname[0], &sname[0])
				act(&buf[0], TRUE, vict, nil, unsafe.Pointer(tch), TO_CHAR)
				act(&buf2[0], TRUE, vict, nil, unsafe.Pointer(tch), TO_VICT)
				act(&buf3[0], TRUE, vict, nil, unsafe.Pointer(tch), TO_NOTVICT)
				dmg = damtype(ch, type_, skill, attperc)
				hurt(0, 0, ch, tch, nil, dmg, 1)
				return
			}
		}
	}
	for tob = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Contents; tob != nil; tob = next_obj {
		next_obj = tob.Next_content
		if OBJ_FLAGGED(tob, ITEM_UNBREAKABLE) {
			continue
		}
		if foundo == TRUE {
			continue
		}
		if rand_number(1, 101) >= 80 {
			foundo = TRUE
			stdio.Sprintf(&buf[0], "@WYou watch as the deflected %s slams into @g$p@W, exploding with a roar of blinding light!@n", &sname[0])
			stdio.Sprintf(&buf2[0], "@c$n@W watches as the deflected %s slams into @g$p@W, exploding with a roar of blinding light!@n", &sname[0])
			act(&buf[0], TRUE, vict, tob, nil, TO_CHAR)
			act(&buf2[0], TRUE, vict, tob, nil, TO_ROOM)
			hurt(0, 0, ch, nil, tob, 25, 1)
			return
		}
	}
	if (foundo == FALSE || foundv == FALSE) && !ROOM_FLAGGED(vict.In_room, ROOM_SPACE) {
		stdio.Sprintf(&buf[0], "@WYou watch as the deflected %s slams into the ground, exploding with a roar of blinding light!@n", &sname[0])
		stdio.Sprintf(&buf2[0], "@WThe deflected %s slams into the ground, exploding with a roar of blinding light!@n", &sname[0])
		act(&buf[0], TRUE, vict, nil, nil, TO_CHAR)
		act(&buf2[0], TRUE, vict, nil, nil, TO_ROOM)
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
		}()) != SECT_WATER_NOSWIM && (func() int {
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
					act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
				act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
				act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
		}
		var zone int = 0
		if (func() int {
			zone = int(real_zone_by_thing(func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()))
			return zone
		}()) != int(-1) {
			send_to_zone(libc.CString("An explosion shakes the entire area!\r\n"), zone_rnum(zone))
		}
		return
	}
}
func dodge_ki(ch *char_data, vict *char_data, type_ int, type2 int, skill int, skill2 int) {
	if type_ == 0 && !ROOM_FLAGGED(vict.In_room, ROOM_SPACE) {
		if (func() int {
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_INSIDE {
			impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
			switch rand_number(1, 8) {
			case 1:
				act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			case 2:
				if rand_number(1, 4) == 4 && (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect == 0 {
					(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect = 1
					act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("Lava spews up through cracks in the ground, roaring into the sky as a large column of molten rock!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
				act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("The explosion continues to burn spreading out and devouring some more of the ground before dying out."), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
			if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_INSIDE {
			impact_sound(ch, libc.CString("@wA loud roar is heard nearby!@n\r\n"))
			switch rand_number(1, 8) {
			case 1:
				act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("Debris is thrown into the air and showers down thunderously!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			case 2:
				act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("The structure of the surrounding room cracks and quakes from the blast!"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
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
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
			(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
		}
		var zone int = 0
		if (func() int {
			zone = int(real_zone_by_thing(func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()))
			return zone
		}()) != int(-1) {
			send_to_zone(libc.CString("An explosion shakes the entire area!\r\n"), zone_rnum(zone))
		}
	}
	if type_ == 1 {
		if rand_number(1, 3) != 2 {
			act(libc.CString("@RIt turns around at the last second and begins to pursue @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RIt turns around at the last second and begins to pursue YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@RIt turns around at the last second and begins to pursue @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			var obj *obj_data
			var num int = 0
			switch skill2 {
			case 461:
				num = 80
			default:
				num = 80
			}
			obj = read_object(obj_vnum(num), VIRTUAL)
			obj_to_room(obj, ch.In_room)
			obj.Target = vict
			obj.Kicharge = damtype(ch, type2, skill, 0.2)
			obj.Kitype = skill2
			obj.User = ch
		} else {
			act(libc.CString("@RIt fails to follow after @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@RIt fails to follow after YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@RIt fails to follow after @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		}
	}
	if type_ == 2 && (skill2 != 481 || ch.Chclass == CLASS_FRIEZA) {
		if skill2 == 481 {
			var (
				chance int = rand_number(25, 50)
				prob   int = axion_dice(0)
			)
			if GET_SKILL(ch, SKILL_KIENZAN) >= 100 {
				chance += int(float64(chance) * 0.8)
			} else if GET_SKILL(ch, SKILL_KIENZAN) >= 60 {
				chance += int(float64(chance) * 0.5)
			} else if GET_SKILL(ch, SKILL_KIENZAN) >= 40 {
				chance += int(float64(chance) * 0.25)
			}
			if chance < prob {
				return
			}
		}
		act(libc.CString("@RYou turn it around and send it back after @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
		act(libc.CString("@W$n @Rturns it around and sends it back after YOU!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		act(libc.CString("@W$n @Rturns it around and sends it back after @r$N@R!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		var obj *obj_data
		var num int = 0
		switch skill2 {
		case 496:
			num = 81
		case 481:
			num = 84
		default:
			num = 81
		}
		obj = read_object(obj_vnum(num), VIRTUAL)
		obj_to_room(obj, ch.In_room)
		obj.Target = vict
		obj.Kicharge = damtype(ch, type2, skill, 0.3)
		obj.Kitype = skill2
		obj.User = ch
	}
}
func damtype(ch *char_data, type_ int, skill int, percent float64) int64 {
	var (
		dam   int64 = 0
		cou1  int64 = 0
		cou2  int64 = 0
		focus int64 = 0
	)
	if !IS_NPC(ch) {
		if GET_SKILL(ch, SKILL_FOCUS) != 0 {
			focus = int64(GET_SKILL(ch, SKILL_FOCUS))
		}
		if type_ != -2 {
			ch.Lastattack = type_
		} else {
			type_ = 0
		}
		switch type_ {
		case -1:
			cou1 = int64(((skill / 4) * int((ch.Hit/1200)+int64(ch.Aff_abils.Str))) + 1)
			cou2 = int64(((skill / 4) * int((ch.Hit/1000)+int64(ch.Aff_abils.Str))) + 1)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if AFF_FLAGGED(ch, AFF_HASS) && !PLR_FLAGGED(ch, PLR_THANDW) {
				dam *= 2
				if ch.Chclass == CLASS_KRANE {
					if GET_SKILL(ch, SKILL_HASSHUKEN) >= 100 {
						dam += int64(float64(dam) * 0.3)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 60 {
						dam += int64(float64(dam) * 0.2)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 40 {
						dam += int64(float64(dam) * 0.1)
					}
				}
			} else if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Preference == PREFERENCE_KI {
				dam -= int64(float64(dam) * 0.2)
			}
			if ch.Preference == PREFERENCE_WEAPON && float64(ch.Charge) >= float64(ch.Max_mana)*0.05 {
				dam += int64(float64(ch.Max_mana) * 0.05)
				ch.Charge -= int64(float64(ch.Max_mana) * 0.05)
			} else if ch.Preference == PREFERENCE_WEAPON && ch.Charge > 0 {
				dam += ch.Charge
				ch.Charge -= 0
			}
			if group_bonus(ch, 2) == 8 {
				dam += int64(float64(dam) * 0.02)
			}
		case 0:
			cou1 = int64(((skill / 4) * int((ch.Hit/1600)+int64(ch.Aff_abils.Str))) + 15)
			cou2 = int64(((skill / 4) * int((ch.Hit/1300)+int64(ch.Aff_abils.Str))) + 15)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_HASS) {
				dam *= 2
				if ch.Chclass == CLASS_KRANE {
					if GET_SKILL(ch, SKILL_HASSHUKEN) >= 100 {
						dam += int64(float64(dam) * 0.3)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 60 {
						dam += int64(float64(dam) * 0.2)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 40 {
						dam += int64(float64(dam) * 0.1)
					}
				}
			} else if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ROSHI {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.2)
				}
			} else if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 1:
			cou1 = int64(((skill / 4) * int((ch.Hit/1200)+int64(ch.Aff_abils.Str))) + 40)
			cou2 = int64(((skill / 4) * int((ch.Hit/1000)+int64(ch.Aff_abils.Str))) + 40)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Chclass == CLASS_KRANE {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.2)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 2:
			cou1 = int64(((skill / 4) * int((ch.Hit/1300)+int64(ch.Aff_abils.Str))) + 100)
			cou2 = int64(((skill / 4) * int((ch.Hit/1050)+int64(ch.Aff_abils.Str))) + 100)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_HASS) {
				dam *= 2
				if ch.Chclass == CLASS_KRANE {
					if GET_SKILL(ch, SKILL_HASSHUKEN) >= 100 {
						dam += int64(float64(dam) * 0.3)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 60 {
						dam += int64(float64(dam) * 0.2)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 40 {
						dam += int64(float64(dam) * 0.1)
					}
				}
			} else if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 3:
			cou1 = int64(((skill / 4) * int((ch.Hit/1100)+int64(ch.Aff_abils.Str))) + 150)
			cou2 = int64(((skill / 4) * int((ch.Hit/1000)+int64(ch.Aff_abils.Str))) + 150)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 4:
			cou1 = int64(((skill / 4) * int((ch.Hit/1000)+int64(ch.Aff_abils.Str))) + 500)
			cou2 = int64(((skill / 4) * int((ch.Hit/800)+int64(ch.Aff_abils.Str))) + 500)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 5:
			cou1 = int64(((skill / 4) * int((ch.Hit/1100)+int64(ch.Aff_abils.Str))) + 350)
			cou2 = int64(((skill / 4) * int((ch.Hit/900)+int64(ch.Aff_abils.Str))) + 350)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_HASS) {
				dam *= 2
				if ch.Chclass == CLASS_KRANE {
					if GET_SKILL(ch, SKILL_HASSHUKEN) >= 100 {
						dam += int64(float64(dam) * 0.3)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 60 {
						dam += int64(float64(dam) * 0.2)
					} else if GET_SKILL(ch, SKILL_HASSHUKEN) >= 40 {
						dam += int64(float64(dam) * 0.1)
					}
				}
			} else if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 6:
			cou1 = int64(((skill / 4) * int((ch.Hit/800)+int64(ch.Aff_abils.Str))) + 8000)
			cou2 = int64(((skill / 4) * int((ch.Hit/500)+int64(ch.Aff_abils.Str))) + 8000)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 7:
			dam = int64(float64(ch.Max_mana) * percent)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 1000)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 25
			}
		case 8:
			cou1 = int64(((skill / 4) * int((ch.Hit/700)+int64(ch.Aff_abils.Str))) + 12500)
			cou2 = int64(((skill / 4) * int((ch.Hit/400)+int64(ch.Aff_abils.Str))) + 12500)
			dam = large_rand(cou1, cou2)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 9:
			dam = int64(float64(ch.Max_mana) * percent)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 500)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 25
			}
		case 10:
			dam = int64(float64(ch.Max_mana) * percent)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 25
			}
		case 11:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 500)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 25
			}
		case 12:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 500)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 25
			}
		case 13:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 800)
			dam += int64(float64(dam) * 0.25)
			if focus > 0 {
				dam += int64((float64(dam) * 0.005) * float64(focus))
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 14:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 15:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 650)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 16:
			dam += int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 800)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 17:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 650)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 18:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 19:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 650)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 20:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1200)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 21:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 900)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
				if skill == 101 {
					dam = int64(float64(dam) * 1.1)
				} else if skill == 102 {
					dam = int64(float64(dam) * 1.2)
				} else if skill == 103 {
					dam = int64(float64(dam) * 1.3)
				}
			}
		case 22:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 600)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		case 23:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 500)
			dam *= int64(1.25)
			dam += (dam / 100) * int64(ch.Aff_abils.Str)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 24:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 600)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 25:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 500)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 26:
			if !IS_NPC(ch) && percent > 0.15 {
				var (
					hitperc    float64 = (percent - 0.15) * 5
					amount     int64   = int64(float64(gear_pl(ch)) * hitperc)
					difference int64   = ch.Hit - amount
				)
				if difference < 1 {
					dam = int64(float64(ch.Max_mana) * percent)
					dam += int64(GET_LEVEL(ch) * 800)
					dam *= int64(1.25)
					dam += ch.Hit - 1
					ch.Hit = 1
				} else {
					ch.Hit = difference
					dam = int64(float64(ch.Max_mana) * percent)
					dam += int64(GET_LEVEL(ch) * 800)
					dam *= int64(1.25)
					dam += amount
				}
				if focus > 0 {
					dam += focus * (dam / 200)
				}
			} else {
				dam = int64(float64(ch.Max_mana) * percent)
				dam += int64(GET_LEVEL(ch) * 800)
				dam *= int64(1.25)
				if focus > 0 {
					dam += focus * (dam / 200)
				}
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 27:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1200)
			dam *= int64(1.25)
			dam += (dam / 100) * int64(ch.Aff_abils.Intel)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 28:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1500)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 29:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1200)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 30:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 31:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1100)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 32:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1400)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 33:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
		case 34:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1050)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 35:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1600)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 36:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1100)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 37:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1200)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 38:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 39:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 900)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 40:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 2000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_KAI {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 41:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 2000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_KAI {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 42:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 550)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 43:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 44:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 45:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1000)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 46:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1400)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		case 47:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 900)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		case 48:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 900)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		case 49:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 900)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if !IS_NPC(ch) {
				if PLR_FLAGGED(ch, PLR_TRANS6) {
					dam += dam
				} else if PLR_FLAGGED(ch, PLR_TRANS5) {
					dam += int64((float64(dam) * 0.01) * 75)
				} else if PLR_FLAGGED(ch, PLR_TRANS4) {
					dam += int64((float64(dam) * 0.01) * 50)
				} else if PLR_FLAGGED(ch, PLR_TRANS3) {
					dam += int64((float64(dam) * 0.01) * 25)
				} else if PLR_FLAGGED(ch, PLR_TRANS2) {
					dam += int64((float64(dam) * 0.01) * 15)
				} else if PLR_FLAGGED(ch, PLR_TRANS1) {
					dam += int64((float64(dam) * 0.01) * 5)
				}
			}
		case 50:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 800)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		case 51:
			cou1 = int64(((skill / 4) * int((ch.Hit/700)+int64(ch.Aff_abils.Str))) + 1000)
			cou2 = int64(((skill / 4) * int((ch.Hit/550)+int64(ch.Aff_abils.Str))) + 1000)
			dam = large_rand(cou1, cou2)
			dam += int64(GET_LEVEL(ch) * 100)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
		case 52:
			cou1 = int64(((skill / 4) * int((ch.Hit/900)+int64(ch.Aff_abils.Str))) + 800)
			cou2 = int64(((skill / 4) * int((ch.Hit/650)+int64(ch.Aff_abils.Str))) + 800)
			dam = large_rand(cou1, cou2)
			dam += int64(GET_LEVEL(ch) * 100)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if ch.Race == RACE_ARLIAN {
				dam += int64(float64(dam) * 0.02)
			}
			if (ch.Bonuses[BONUS_BRAWLER]) > 0 {
				dam += int64(float64(dam) * 0.2)
			}
			if ch.Chclass == CLASS_ANDSIX {
				if (ch.Skills[SKILL_STYLE]) >= 75 {
					dam += int64(float64(dam) * 0.1)
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
		case 53:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1600)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN {
				dam += (dam / 100) * 15
			}
		case 54:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_SAIYAN {
				dam += (dam / 100) * 20
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 55:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
			if ch.Race == RACE_HUMAN && skill == 101 {
				dam = int64(float64(dam) * 1.1)
			} else if ch.Race == RACE_HUMAN && skill == 102 {
				dam = int64(float64(dam) * 1.2)
			} else if ch.Race == RACE_HUMAN && skill == 103 {
				dam = int64(float64(dam) * 1.3)
			}
		case 56:
			cou1 = int64(((skill / 4) * int((ch.Hit/1100)+int64(ch.Aff_abils.Str))) + 400)
			cou2 = int64(((skill / 4) * int((ch.Hit/1000)+int64(ch.Aff_abils.Str))) + 400)
			dam = large_rand(cou1, cou2)
			dam += int64(GET_LEVEL(ch) * 100)
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
			if AFF_FLAGGED(ch, AFF_INFUSE) {
				dam += (dam / 100) * int64(GET_SKILL(ch, SKILL_INFUSE)/2)
				if ch.Chclass == CLASS_JINTO {
					if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.5)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.25)
					} else if GET_SKILL(ch, SKILL_INFUSE) >= 100 {
						dam += int64(((float64(dam) * 0.01) * float64(GET_SKILL(ch, SKILL_INFUSE)/2)) * 0.05)
					}
				}
			}
			if ch.Preference == PREFERENCE_THROWING {
				dam -= int64(float64(dam) * 0.15)
			} else if ch.Preference == PREFERENCE_H2H {
				dam += int64(float64(dam) * 0.2)
			}
		case 57:
			dam = int64(float64(ch.Max_mana) * percent)
			dam += int64(GET_LEVEL(ch) * 1700)
			dam *= int64(1.25)
			if focus > 0 {
				dam += focus * (dam / 200)
			}
		}
	} else {
		dam = int64((float64(ch.Hit) * 0.05) + float64(ch.Max_hit)*0.025)
		dam += int64((float64(dam) * 0.005) * float64(ch.Aff_abils.Str))
		if GET_LEVEL(ch) >= 120 {
			dam *= int64(0.25)
		} else if GET_LEVEL(ch) >= 110 {
			dam *= int64(0.45)
		} else if GET_LEVEL(ch) >= 100 {
			dam *= int64(0.75)
		}
	}
	if IS_NPC(ch) {
		if type_ == 0 || type_ == 1 || type_ == 2 || type_ == 3 || type_ == 4 || type_ == 5 || type_ == 6 || type_ == 8 || type_ == 51 || type_ == 52 || type_ == 56 {
			dam += int64(float64(ch.Aff_abils.Str) * (float64(dam) * 0.005))
		} else {
			dam += int64(float64(ch.Aff_abils.Intel) * (float64(dam) * 0.005))
		}
		var mobperc int64 = (ch.Hit * 100) / ch.Max_hit
		if mobperc < 98 && mobperc >= 90 {
			dam = int64(float64(dam) * 0.95)
		} else if mobperc < 90 && mobperc >= 80 {
			dam = int64(float64(dam) * 0.9)
		} else if mobperc < 80 && mobperc >= 790 {
			dam = int64(float64(dam) * 0.85)
		} else if mobperc < 70 && mobperc >= 50 {
			dam = int64(float64(dam) * 0.8)
		} else if mobperc < 50 && mobperc >= 30 {
			dam = int64(float64(dam) * 0.7)
		} else if mobperc <= 29 {
			dam = int64(float64(dam) * 0.6)
		}
		if ch.Chclass != CLASS_NPC_COMMONER {
			dam += int64(float64(dam) * 0.3)
		}
	}
	if ch.Kaioken > 0 {
		dam += (dam / 200) * int64(ch.Kaioken)
	}
	if PLR_FLAGGED(ch, PLR_FURY) && (type_ == 0 || type_ == 1 || type_ == 2 || type_ == 3 || type_ == 4 || type_ == 5 || type_ == 6 || type_ == 8 || type_ == 51 || type_ == 52) {
		dam *= int64(1.5)
		act(libc.CString("Your rage magnifies your attack power!"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("Swirling energy flows around $n as $e releases $s rage in the attack!"), TRUE, ch, nil, nil, TO_ROOM)
		if rand_number(1, 10) >= 7 {
			send_to_char(ch, libc.CString("You feel less angry.\r\n"))
			ch.Act[int(PLR_FURY/32)] &= bitvector_t(^(1 << (int(PLR_FURY % 32))))
		}
	} else if PLR_FLAGGED(ch, PLR_FURY) {
		dam *= 2
		act(libc.CString("Your rage magnifies your attack power!"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("Swirling energy flows around $n as $e releases $s rage in the attack!"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_FURY/32)] &= bitvector_t(^(1 << (int(PLR_FURY % 32))))
	}
	if type_ == -1 || type_ == 0 || type_ == 1 || type_ == 2 || type_ == 3 || type_ == 4 || type_ == 5 || type_ == 6 || type_ == 8 {
		if !IS_NPC(ch) {
			dam -= int64(float64(dam) * 0.08)
		}
		if !IS_NPC(ch) && float64(dam) > float64(ch.Max_hit)*0.1 {
			dam *= int64(0.6)
		}
	} else {
		dam += int64((float64(dam) * 0.005) * float64(ch.Aff_abils.Intel))
		if ch.Preference == PREFERENCE_WEAPON {
			dam -= int64(float64(dam) * 0.25)
		} else if ch.Preference == PREFERENCE_THROWING {
			dam -= int64(float64(dam) * 0.15)
		}
	}
	return dam
}
func saiyan_gain(ch *char_data, vict *char_data) {
	var (
		gain int = rand_number(GET_LEVEL(ch)*6, GET_LEVEL(ch)*8)
		weak int = FALSE
	)
	if vict == nil {
		return
	}
	if IS_NPC(ch) {
		return
	}
	if vict.Max_hit < ch.Max_hit/10 {
		weak = TRUE
	}
	if GET_LEVEL(ch) > 99 {
		gain += rand_number(GET_LEVEL(ch)*300, GET_LEVEL(ch)*500)
	} else if GET_LEVEL(ch) > 80 {
		gain += rand_number(GET_LEVEL(ch)*150, GET_LEVEL(ch)*200)
	} else if GET_LEVEL(ch) > 60 {
		gain += rand_number(GET_LEVEL(ch)*80, GET_LEVEL(ch)*100)
	} else if GET_LEVEL(ch) > 50 {
		gain += rand_number(GET_LEVEL(ch)*20, GET_LEVEL(ch)*25)
	} else if GET_LEVEL(ch) > 40 {
		gain += rand_number(GET_LEVEL(ch)*8, GET_LEVEL(ch)*10)
	} else if GET_LEVEL(ch) > 30 {
		gain += rand_number(GET_LEVEL(ch)*5, GET_LEVEL(ch)*8)
	} else {
	}
	if ch.Race == RACE_BIO && ((ch.Genome[0]) == 2 || (ch.Genome[1]) == 2) {
		gain /= 2
	}
	if rand_number(1, 22) >= 18 && (GET_LEVEL(ch) == 100 || level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0) {
		if weak == TRUE {
			send_to_char(ch, libc.CString("@D[@YSaiyan @RBlood@D] @WThey are too weak to inspire your saiyan soul!@n\r\n"))
		} else {
			switch rand_number(1, 3) {
			case 1:
				ch.Max_hit += int64(gain)
				ch.Basepl += int64(gain)
				send_to_char(ch, libc.CString("@D[@YSaiyan @RBlood@D] @WYou feel slightly stronger. @D[@G+%s@D]@n\r\n"), add_commas(int64(gain)))
			case 2:
				ch.Max_mana += int64(gain)
				ch.Baseki += int64(gain)
				send_to_char(ch, libc.CString("@D[@YSaiyan @RBlood@D] @WYou feel your spirit grow. @D[@G+%s@D]@n\r\n"), add_commas(int64(gain)))
			case 3:
				ch.Max_move += int64(gain)
				ch.Basest += int64(gain)
				send_to_char(ch, libc.CString("@D[@YSaiyan @RBlood@D] @WYou feel slightly more vigorous. @D[@G+%s@D]@n\r\n"), add_commas(int64(gain)))
			}
		}
	}
}
func spar_gain(ch *char_data, vict *char_data, type_ int, dmg int64) {
	var (
		chance     int = 0
		gmult      int
		gravity    int
		bonus      int   = 1
		pscost     int   = 2
		difference int   = 0
		gain       int64 = 0
		pl         int64 = 0
		ki         int64 = 0
		st         int64 = 0
		gaincalc   int64 = 0
	)
	if ch != nil && !IS_NPC(ch) {
		if dmg > vict.Max_hit/10 {
			chance = rand_number(20, 100)
		} else if dmg <= vict.Max_hit/10 {
			chance = rand_number(1, 75)
		}
		if ch.Relax_count >= 464 {
			chance = 0
		} else if ch.Relax_count >= 232 {
			chance -= int(float64(chance) * 0.5)
		} else if ch.Relax_count >= 116 {
			chance -= int(float64(chance) * 0.2)
		}
		gravity = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
		gmult = GET_LEVEL(ch) * ((gravity / 5) + 6)
		if (ch.Equipment[WEAR_SH]) != nil {
			var obj *obj_data = (ch.Equipment[WEAR_SH])
			if GET_OBJ_VNUM(obj) == 1127 {
				gmult *= 4
			}
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) || ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			if (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				gmult *= int(1.75)
				pscost += 2
			} else {
				gmult *= int(1.25)
				pscost += 1
			}
			pl = large_rand(int64(float64(gmult)*0.8), int64(float64(gmult)*1.2))
			ki = large_rand(int64(float64(gmult)*0.8), int64(float64(gmult)*1.2))
		} else {
			pl = large_rand(int64(float64(gmult)*0.4), int64(float64(gmult)*0.8))
			ki = large_rand(int64(float64(gmult)*0.4), int64(float64(gmult)*0.8))
		}
		if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) < 0 && GET_LEVEL(ch) < 100 {
			pl = 0
		}
		if level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) < 0 && GET_LEVEL(ch) < 100 {
			ki = 0
		}
		if chance >= rand_number(60, 75) {
			var (
				num    int64 = 0
				maxnum int64 = 500000
			)
			if GET_LEVEL(ch) >= 70 {
				num += int64(GET_LEVEL(ch) * 10000)
			} else if GET_LEVEL(ch) >= 60 {
				num += int64(GET_LEVEL(ch) * 6000)
			} else if GET_LEVEL(ch) >= 50 {
				num += int64(GET_LEVEL(ch) * 5000)
			} else if GET_LEVEL(ch) >= 45 {
				num += int64(GET_LEVEL(ch) * 2500)
			} else if GET_LEVEL(ch) >= 40 {
				num += int64(GET_LEVEL(ch) * 2200)
			} else if GET_LEVEL(ch) >= 35 {
				num += int64(GET_LEVEL(ch) * 1500)
			} else if GET_LEVEL(ch) >= 30 {
				num += int64(GET_LEVEL(ch) * 1200)
			} else if GET_LEVEL(ch) >= 25 {
				num += int64(GET_LEVEL(ch) * 550)
			} else if GET_LEVEL(ch) >= 20 {
				num += int64(GET_LEVEL(ch) * 400)
			} else if GET_LEVEL(ch) >= 15 {
				num += int64(GET_LEVEL(ch) * 250)
			} else if GET_LEVEL(ch) >= 10 {
				num += int64(GET_LEVEL(ch) * 100)
			} else if GET_LEVEL(ch) >= 5 {
				num += int64(GET_LEVEL(ch) * 50)
			} else {
				num += int64(GET_LEVEL(ch) * 30)
			}
			if num > maxnum {
				num = maxnum
			}
			if vict != nil && IS_NPC(vict) {
				num = int64(float64(num) * 0.7)
				gaincalc = int64(float64(num) * 1.5)
				type_ = 3
			} else if vict != nil && !IS_NPC(vict) {
				gaincalc = large_rand(int64(float64(num)*0.7), num)
				if GET_LEVEL(ch) > GET_LEVEL(vict) {
					difference = GET_LEVEL(ch) - GET_LEVEL(vict)
				} else if GET_LEVEL(ch) < GET_LEVEL(vict) {
					difference = GET_LEVEL(vict) - GET_LEVEL(ch)
				}
			} else {
				gaincalc = 0
			}
			if vict != nil {
				if difference >= 51 {
					send_to_char(ch, libc.CString("The difference in your levels is too great for you to gain anything.\r\n"))
					return
				} else if difference >= 40 {
					gaincalc = int64(float64(gaincalc) * 0.05)
				} else if difference >= 30 {
					gaincalc = int64(float64(gaincalc) * 0.1)
				} else if difference >= 20 {
					gaincalc = int64(float64(gaincalc) * 0.25)
				} else if difference >= 10 {
					gaincalc = int64(float64(gaincalc) * 0.5)
				}
				if !IS_NPC(vict) {
					if PRF_FLAGGED(vict, PRF_INSTRUCT) {
						if (vict.Player_specials.Class_skill_points[vict.Chclass]) > 10 {
							send_to_char(vict, libc.CString("You instruct them in proper fighting techniques and strategies.\r\n"))
							act(libc.CString("You take $N's instruction to heart and gain more experience.\r\n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							vict.Player_specials.Class_skill_points[vict.Chclass] -= 10
							bonus = 2
						}
					}
				}
			}
			if ch.Race == RACE_SAIYAN {
				gaincalc = int64(float64(gaincalc) + float64(gaincalc)*0.5)
			}
			if ch.Race == RACE_HALFBREED {
				gaincalc = int64(float64(gaincalc) + float64(gaincalc)*0.4)
			}
			if ch.Race == RACE_ICER || ch.Race == RACE_BIO && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				gaincalc = int64(float64(gaincalc) - float64(gaincalc)*0.2)
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) || ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
				if (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) >= 19100 && (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) <= 0x4AFF {
					gaincalc *= int64(1.5)
				} else {
					gaincalc *= int64(1.25)
				}
			}
			gain = gear_exp(ch, gaincalc)
			if (ch.Player_specials.Class_skill_points[ch.Chclass]) >= pscost {
				ch.Player_specials.Class_skill_points[ch.Chclass] -= pscost
				gain = gain * int64(bonus)
				gain_exp(ch, gain)
				send_to_char(ch, libc.CString("@D[@Y+ @G%s @mExp@D]@n "), add_commas(gain))
				if type_ == 0 && rand_number(1, 5) >= 4 {
					send_to_char(ch, libc.CString("@D[@Y+ @R%s @rPL@D]@n "), func() *byte {
						if pl > 0 {
							return add_commas(pl)
						}
						return libc.CString("SOFT-CAP")
					}())
					ch.Max_hit += pl
					ch.Basepl += pl
				} else if type_ == 1 && rand_number(1, 5) >= 4 {
					send_to_char(ch, libc.CString("@D[@Y+ @C%s @cKi@D]@n "), func() *byte {
						if ki > 0 {
							return add_commas(ki)
						}
						return libc.CString("SOFT-CAP")
					}())
					ch.Max_mana += ki
					ch.Baseki += ki
				}
				send_to_char(ch, libc.CString("@D[@R- @M%d @mPS@D]@n "), pscost)
				send_to_char(ch, libc.CString("\r\n"))
			} else {
				send_to_char(ch, libc.CString("@RYou need at least %d Practice Sessions in order to gain while sparring here.@n\r\n"), pscost)
			}
		}
	}
	if vict != nil && !IS_NPC(vict) && !IS_NPC(ch) {
		if dmg > vict.Max_hit/4 {
			chance = rand_number(1, 100)
		} else if dmg <= vict.Max_hit/4 {
			chance = rand_number(1, 70)
		}
		if vict.Relax_count >= 464 {
			chance = 0
		} else if vict.Relax_count >= 232 {
			chance -= int(float64(chance) * 0.5)
		} else if vict.Relax_count >= 116 {
			chance -= int(float64(chance) * 0.2)
		}
		gravity = (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity
		gmult = GET_LEVEL(vict) * ((gravity / 5) + 6)
		if ROOM_FLAGGED(vict.In_room, ROOM_WORKOUT) || ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
			if (func() room_vnum {
				if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Number
				}
				return -1
			}()) >= 19100 && (func() room_vnum {
				if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
					return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Number
				}
				return -1
			}()) <= 0x4AFF {
				gmult *= int(1.75)
			} else {
				gmult *= int(1.25)
			}
			st = large_rand(int64(float64(gmult)*0.8), int64(float64(gmult)*1.2))
		} else {
			st = large_rand(int64(float64(gmult)*0.4), int64(float64(gmult)*0.8))
		}
		if level_exp(vict, GET_LEVEL(vict)+1)-int(vict.Exp) < 0 && GET_LEVEL(vict) < 100 {
			st = 0
		}
		if chance >= rand_number(60, 75) {
			var (
				num    int64 = 0
				maxnum int64 = 500000
			)
			if GET_LEVEL(vict) >= 70 {
				num += int64(GET_LEVEL(vict) * 10000)
			} else if GET_LEVEL(vict) >= 60 {
				num += int64(GET_LEVEL(vict) * 6000)
			} else if GET_LEVEL(vict) >= 50 {
				num += int64(GET_LEVEL(vict) * 5000)
			} else if GET_LEVEL(vict) >= 45 {
				num += int64(GET_LEVEL(vict) * 2500)
			} else if GET_LEVEL(vict) >= 40 {
				num += int64(GET_LEVEL(vict) * 2200)
			} else if GET_LEVEL(vict) >= 35 {
				num += int64(GET_LEVEL(vict) * 1500)
			} else if GET_LEVEL(vict) >= 30 {
				num += int64(GET_LEVEL(vict) * 1200)
			} else if GET_LEVEL(vict) >= 25 {
				num += int64(GET_LEVEL(vict) * 550)
			} else if GET_LEVEL(vict) >= 20 {
				num += int64(GET_LEVEL(vict) * 400)
			} else if GET_LEVEL(vict) >= 15 {
				num += int64(GET_LEVEL(vict) * 250)
			} else if GET_LEVEL(vict) >= 10 {
				num += int64(GET_LEVEL(vict) * 100)
			} else if GET_LEVEL(vict) >= 5 {
				num += int64(GET_LEVEL(vict) * 50)
			} else {
				num += int64(GET_LEVEL(vict) * 30)
			}
			if num > maxnum {
				num = maxnum
			}
			gain = large_rand(int64(float64(num)*0.7), num)
			if difference > 50 {
				send_to_char(ch, libc.CString("The difference in your levels is too great for you to gain anything.\r\n"))
				return
			} else if difference >= 40 {
				gain = int64(float64(gain) * 0.05)
			} else if difference >= 30 {
				gain = int64(float64(gain) * 0.1)
			} else if difference >= 20 {
				gain = int64(float64(gain) * 0.25)
			} else if difference >= 10 {
				gain = int64(float64(gain) * 0.5)
			}
			if vict.Race == RACE_SAIYAN || vict.Race == RACE_HALFBREED {
				gain = int64(float64(gain) + float64(gain)*0.3)
			}
			if vict.Race == RACE_ICER {
				gain = int64(float64(gain) - float64(gain)*0.1)
			}
			if ROOM_FLAGGED(ch.In_room, ROOM_WORKOUT) || ROOM_FLAGGED(ch.In_room, ROOM_HBTC) {
				if (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) >= 19100 && (func() room_vnum {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number
					}
					return -1
				}()) <= 0x4AFF {
					gain *= int64(1.5)
				} else {
					gain *= int64(1.25)
				}
			}
			if (vict.Player_specials.Class_skill_points[vict.Chclass]) >= pscost {
				vict.Player_specials.Class_skill_points[vict.Chclass] -= pscost
				send_to_char(vict, libc.CString("@D[@Y+ @G%s @mExp@D]@n "), add_commas(gain))
				gain = gear_exp(vict, gain)
				gain_exp(vict, gain)
				if rand_number(1, 5) >= 4 {
					send_to_char(vict, libc.CString("@D[@Y+ @G%s @gSt@D]@n "), func() *byte {
						if st > 0 {
							return add_commas(st)
						}
						return libc.CString("SOFT-CAP")
					}())
					vict.Max_move += st
					vict.Basest += st
				}
				send_to_char(vict, libc.CString("@D[@R- @M%d @mPS@D]@n "), pscost)
				send_to_char(vict, libc.CString("\r\n"))
			} else {
				send_to_char(vict, libc.CString("@RYou need at least %d Practice Sessions in order to gain while sparring here.@n\r\n"), pscost)
			}
		}
	}
}
func can_grav(ch *char_data) int {
	if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10 && ch.Max_hit < 5000 && ch.Chclass != CLASS_BARDOCK && !IS_NPC(ch) {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 20 && ch.Max_hit < 20000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 30 && ch.Max_hit < 50000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 40 && ch.Max_hit < 100000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 50 && ch.Max_hit < 200000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 100 && ch.Max_hit < 400000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 200 && ch.Max_hit < 1000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 300 && ch.Max_hit < 5000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 400 && ch.Max_hit < 8000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 500 && ch.Max_hit < 15000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 1000 && ch.Max_hit < 25000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 5000 && ch.Max_hit < 100000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Gravity == 10000 && ch.Max_hit < 200000000 {
		send_to_char(ch, libc.CString("You are hardly able to move in this gravity!\r\n"))
		return 0
	} else {
		return 1
	}
}
func can_kill(ch *char_data, vict *char_data, obj *obj_data, num int) int {
	if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_HEALT) {
		send_to_char(ch, libc.CString("You are inside a healing tank!\r\n"))
		return 0
	}
	if ch.Carry_weight > int(max_carry_weight(ch)) {
		send_to_char(ch, libc.CString("You are weighted down too much!\r\n"))
		return 0
	}
	if vict != nil {
		if vict.Hit <= 0 && vict.Fighting != nil {
			return 0
		}
		if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
			send_to_char(ch, libc.CString("This room just has such a peaceful, easy feeling...\r\n"))
			return 0
		} else if vict == ch {
			send_to_char(ch, libc.CString("That's insane, don't hurt yourself. Hurt others! That's the key to life ^_^\r\n"))
			return 0
		} else if vict.Gooptime > 0 {
			send_to_char(ch, libc.CString("It seems like it'll be hard to kill them right now...\r\n"))
			return 0
		} else if ch.Player_specials.Carrying != nil {
			send_to_char(ch, libc.CString("You are too busy protecting the person on your shoulder!\r\n"))
			return 0
		} else if vict.Player_specials.Carried_by != nil {
			send_to_char(ch, libc.CString("They are being protected by someone else!\r\n"))
			return 0
		} else if AFF_FLAGGED(vict, AFF_PARALYZE) {
			send_to_char(ch, libc.CString("They are a statue, just leave them alone...\r\n"))
			return 0
		} else if MOB_FLAGGED(vict, MOB_NOKILL) {
			send_to_char(ch, libc.CString("But they are not to be killed!\r\n"))
			return 0
		} else if ch.Majinize == int(vict.Id) {
			send_to_char(ch, libc.CString("You can not harm your master!\r\n"))
			return 0
		} else if (ch.Bonuses[BONUS_COWARD]) > 0 && float64(vict.Max_hit) > float64(ch.Max_hit)+float64(ch.Max_hit)*0.5 && ch.Fighting == nil {
			send_to_char(ch, libc.CString("You are too cowardly to start anything with someone so much stronger than yourself!\r\n"))
			return 0
		} else if vict.Majinize == int(ch.Id) {
			send_to_char(ch, libc.CString("You can not harm your servant.\r\n"))
			return 0
		} else if ch.Grappling != nil && ch.Grap != 3 || ch.Grappled != nil && (ch.Grap == 1 || ch.Grap == 4) {
			send_to_char(ch, libc.CString("You are too busy grappling!%s\r\n"), func() string {
				if ch.Grappled != nil {
					return " Try 'escape'!"
				}
				return ""
			}())
			return 0
		} else if ch.Grappling != nil && ch.Grappling != vict {
			send_to_char(ch, libc.CString("You can't reach that far in your current position!\r\n"))
			return 0
		} else if ch.Grappled != nil && ch.Grappled != vict {
			send_to_char(ch, libc.CString("You can't reach that far in your current position!\r\n"))
			return 0
		} else if !IS_NPC(ch) && !IS_NPC(vict) && AFF_FLAGGED(ch, AFF_SPIRIT) && (!is_sparring(ch) || !is_sparring(vict)) && num != 2 {
			send_to_char(ch, libc.CString("You can not fight other players in AL/Hell.\r\n"))
			return 0
		} else if GET_LEVEL(vict) <= 8 && !IS_NPC(ch) && !IS_NPC(vict) && (!is_sparring(ch) || !is_sparring(vict)) {
			send_to_char(ch, libc.CString("Newbie Shield Protects them!\r\n"))
			return 0
		} else if GET_LEVEL(ch) <= 8 && !IS_NPC(ch) && !IS_NPC(vict) && (!is_sparring(ch) || !is_sparring(vict)) {
			send_to_char(ch, libc.CString("Newbie Shield Protects you until level 8.\r\n"))
			return 0
		} else if PLR_FLAGGED(vict, PLR_SPIRAL) && num != 3 {
			send_to_char(ch, libc.CString("Due to the nature of their current technique anything less than a Tier 4 or AOE attack will not work on them.\r\n"))
			return 0
		} else if ch.Absorbing != nil {
			send_to_char(ch, libc.CString("You are too busy absorbing %s!\r\n"), GET_NAME(ch.Absorbing))
			return 0
		} else if ch.Absorbby != nil {
			send_to_char(ch, libc.CString("You are too busy being absorbed by %s!\r\n"), GET_NAME(ch.Absorbby))
			return 0
		} else if (vict.Altitude-1 > ch.Altitude || vict.Altitude < ch.Altitude-1) && ch.Race == RACE_NAMEK {
			act(libc.CString("@GYou stretch your limbs toward @g$N@G in an attempt to hit $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@g$n@G stretches $s limbs toward @RYOU@G in an attempt to land a hit!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@g$n@G stretches $s limbs toward @g$N@G in an attempt to hit $M!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			return 1
		} else if AFF_FLAGGED(ch, AFF_FLYING) && !AFF_FLAGGED(vict, AFF_FLYING) && num == 0 {
			send_to_char(ch, libc.CString("You are too far above them.\r\n"))
			return 0
		} else if !AFF_FLAGGED(ch, AFF_FLYING) && AFF_FLAGGED(vict, AFF_FLYING) && num == 0 {
			send_to_char(ch, libc.CString("They are too far above you.\r\n"))
			return 0
		} else if !IS_NPC(ch) && ch.Altitude > vict.Altitude && !IS_NPC(vict) && num == 0 {
			if vict.Altitude < 0 {
				vict.Altitude = ch.Altitude
				return 1
			} else {
				send_to_char(ch, libc.CString("You are too far above them.\r\n"))
				return 0
			}
		} else if !IS_NPC(ch) && ch.Altitude < vict.Altitude && !IS_NPC(vict) && num == 0 {
			if vict.Altitude > 2 {
				vict.Altitude = ch.Altitude
				return 1
			} else {
				send_to_char(ch, libc.CString("They are too far above you.\r\n"))
				return 0
			}
		} else {
			return 1
		}
	}
	if obj != nil {
		if ROOM_FLAGGED(ch.In_room, ROOM_PEACEFUL) {
			send_to_char(ch, libc.CString("This room just has such a peaceful, easy feeling...\r\n"))
			return 0
		} else if OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) && GET_OBJ_VNUM(obj) != 87 && GET_OBJ_VNUM(obj) != 80 && GET_OBJ_VNUM(obj) != 81 && GET_OBJ_VNUM(obj) != 82 && GET_OBJ_VNUM(obj) != 83 {
			send_to_char(ch, libc.CString("You can't hit that, it is protected by the immortals!\r\n"))
			return 0
		} else if AFF_FLAGGED(ch, AFF_FLYING) {
			send_to_char(ch, libc.CString("You are too far above it.\r\n"))
			return 0
		} else if OBJ_FLAGGED(obj, ITEM_BROKEN) {
			send_to_char(ch, libc.CString("It is already broken!\r\n"))
			return 0
		} else {
			return 1
		}
	} else {
		send_to_char(ch, libc.CString("Error: Report to imm."))
		return 0
	}
}
func check_skill(ch *char_data, skill int) int {
	if know_skill(ch, skill) == 0 && !IS_NPC(ch) {
		return 0
	} else {
		return 1
	}
}
func check_points(ch *char_data, ki int64, st int64) int {
	if ch.Preference == PREFERENCE_H2H && float64(ch.Charge) >= float64(ch.Max_mana)*0.1 {
		st -= int64(float64(st) * 0.5)
	}
	var fail int = FALSE
	if IS_NPC(ch) {
		if ch.Mana < ki {
			send_to_char(ch, libc.CString("You do not have enough ki!\r\n"))
			fail = TRUE
		}
		if ch.Move < st {
			send_to_char(ch, libc.CString("You do not have enough stamina!\r\n"))
			fail = TRUE
		}
	} else {
		if grav_cost(ch, st) == 0 && ki <= 0 {
			send_to_char(ch, libc.CString("You do not have enough stamina to perform it in this gravity!\r\n"))
			return 0
		}
		if ch.Charge < ki {
			send_to_char(ch, libc.CString("You do not have enough ki charged.\r\n"))
			var perc int64 = int64(float64(ch.Max_mana) * 0.01)
			if ki >= perc*49 {
				send_to_char(ch, libc.CString("You need at least 50 percent charged.\r\n"))
			} else if ki >= perc*44 {
				send_to_char(ch, libc.CString("You need at least 45 percent charged.\r\n"))
			} else if ki >= perc*39 {
				send_to_char(ch, libc.CString("You need at least 40 percent charged.\r\n"))
			} else if ki >= perc*34 {
				send_to_char(ch, libc.CString("You need at least 35 percent charged.\r\n"))
			} else if ki >= perc*29 {
				send_to_char(ch, libc.CString("You need at least 30 percent charged.\r\n"))
			} else if ki >= perc*24 {
				send_to_char(ch, libc.CString("You need at least 25 percent charged.\r\n"))
			} else if ki >= perc*19 {
				send_to_char(ch, libc.CString("You need at least 20 percent charged.\r\n"))
			} else if ki >= perc*14 {
				send_to_char(ch, libc.CString("You need at least 15 percent charged.\r\n"))
			} else if ki >= perc*9 {
				send_to_char(ch, libc.CString("You need at least 10 percent charged.\r\n"))
			} else if ki >= perc*4 {
				send_to_char(ch, libc.CString("You need at least 5 percent charged.\r\n"))
			} else if ki >= 1 {
				send_to_char(ch, libc.CString("You need at least 1 percent charged.\r\n"))
			}
			fail = TRUE
		}
		if IS_NONPTRANS(ch) {
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				st -= int64(float64(st) * 0.2)
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				st -= int64(float64(st) * 0.4)
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				st -= int64(float64(st) * 0.6)
			} else if PLR_FLAGGED(ch, PLR_TRANS4) {
				st -= int64(float64(st) * 0.8)
			}
		}
		if ch.Move < st {
			send_to_char(ch, libc.CString("You do not have enough stamina.\r\n@C%s@n needed.\r\n"), add_commas(st))
			fail = TRUE
		}
	}
	if fail == TRUE {
		return 0
	} else {
		return 1
	}
}
func pcost(ch *char_data, ki float64, st int64) {
	var before int = 0
	if GET_LEVEL(ch) > 1 && !IS_NPC(ch) {
		if ki == 0 {
			before = int(ch.Move)
			if grav_cost(ch, 0) != 0 {
				if before > int(ch.Move) {
					send_to_char(ch, libc.CString("You exert more stamina in this gravity.\r\n"))
				}
			}
		}
		if float64(ch.Charge) <= (float64(ch.Max_mana) * ki) {
			ch.Charge = 0
		}
		if float64(ch.Charge) > (float64(ch.Max_mana) * ki) {
			ch.Charge -= int64(float64(ch.Max_mana) * ki)
		}
		if ch.Charge < 0 {
			ch.Charge = 0
		}
		if ch.Kaioken > 0 {
			st += (st / 20) * int64(ch.Kaioken)
		}
		if AFF_FLAGGED(ch, AFF_HASS) {
			st += int64(float64(st) * 0.3)
		}
		if !IS_NPC(ch) && (ch.Bonuses[BONUS_HARDWORKER]) > 0 {
			st -= int64(float64(st) * 0.25)
		} else if !IS_NPC(ch) && (ch.Bonuses[BONUS_SLACKER]) > 0 {
			st += int64(float64(st) * 0.25)
		}
		if ch.Race == RACE_ICER {
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				st = int64(float64(st) * 1.05)
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				st = int64(float64(st) * 1.1)
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				st = int64(float64(st) * 1.15)
			} else if PLR_FLAGGED(ch, PLR_TRANS4) {
				st = int64(float64(st) * 1.2)
			}
		}
		if ch.Preference == PREFERENCE_H2H && float64(ch.Charge) >= float64(ch.Max_mana)*0.1 {
			st -= int64(float64(st) * 0.5)
			ch.Charge -= st
			if ch.Charge < 0 {
				ch.Charge = 0
			}
		}
		if IS_NONPTRANS(ch) {
			if PLR_FLAGGED(ch, PLR_TRANS1) {
				st -= int64(float64(st) * 0.2)
			} else if PLR_FLAGGED(ch, PLR_TRANS2) {
				st -= int64(float64(st) * 0.4)
			} else if PLR_FLAGGED(ch, PLR_TRANS3) {
				st -= int64(float64(st) * 0.6)
			} else if PLR_FLAGGED(ch, PLR_TRANS4) {
				st -= int64(float64(st) * 0.8)
			}
		}
		ch.Move -= st
	}
	if IS_NPC(ch) {
		ch.Mana -= int64(ki)
		ch.Move -= st
	}
}
func hurt(limb int, chance int, ch *char_data, vict *char_data, obj *obj_data, dmg int64, type_ int) {
	var (
		index     int64 = 0
		maindmg   int64 = dmg
		beforered int64 = dmg
		dead      int   = FALSE
	)
	if type_ <= 0 {
		if ch.Race == RACE_SAIYAN && PLR_FLAGGED(ch, PLR_STAIL) {
			dmg += int64(float64(dmg) * 0.15)
		}
		if ch.Race == RACE_NAMEK && (ch.Equipment[WEAR_HEAD]) == nil {
			dmg += int64(float64(dmg) * 0.25)
		}
		if group_bonus(ch, 2) == 4 {
			dmg += int64(float64(dmg) * 0.1)
		} else if group_bonus(ch, 2) == 12 {
			dmg -= int64(float64(dmg) * 0.1)
		}
	} else {
		dmg = int64(float64(dmg) * 0.6)
		if group_bonus(ch, 2) == 9 {
			dmg -= int64(float64(dmg) * 0.1)
		}
		if AFF_FLAGGED(ch, AFF_POTENT) {
			dmg += int64(float64(dmg) * 0.3)
			send_to_room(ch.In_room, libc.CString("@wThere is a bright flash of @Yyellow@w light in the wake of the attack!@n\r\n"))
		}
	}
	if AFF_FLAGGED(ch, AFF_INFUSE) && !AFF_FLAGGED(ch, AFF_HASS) && type_ <= 0 {
		if float64(ch.Mana)-float64(ch.Max_mana)*0.005 > 0 && dmg > 1 {
			ch.Mana -= int64(float64(ch.Max_mana) * 0.005)
			send_to_room(ch.In_room, libc.CString("@CA swirl of ki explodes from the attack!@n\r\n"))
		} else if float64(ch.Mana)-float64(ch.Max_mana)*0.005 <= 0 {
			act(libc.CString("@wYou can no longer infuse ki into your attacks!@n"), TRUE, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@w can no longer infuse ki into $s attacks!@n"), TRUE, ch, nil, nil, TO_ROOM)
			ch.Affected_by[int(AFF_INFUSE/32)] &= ^(1 << (int(AFF_INFUSE % 32)))
		}
	}
	if vict != nil {
		if (func() room_vnum {
			if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Number
			}
			return -1
		}()) == 0x45D3 {
			return
		}
		reveal_hiding(vict, 0)
		if AFF_FLAGGED(vict, AFF_PARALYZE) {
			send_to_char(ch, libc.CString("They are a statue and can't be harmed\r\n"))
			return
		}
		if ch.Kaioken > 0 {
			dmg += (dmg / 100) * int64(ch.Kaioken*2)
		}
		if vict.Race == RACE_MUTANT && ((vict.Genome[0]) == 8 || (vict.Genome[1]) == 8) && type_ == 0 {
			var drain int64 = int64(float64(dmg) * 0.1)
			dmg -= drain
			ch.Move -= drain
			if ch.Move < 0 {
				ch.Move = 1
			}
			act(libc.CString("@Y$N's rubbery body makes hitting it tiring!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@Y$n's stamina is sapped a bit by hitting your rubbery body!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
		}
		if !IS_NPC(ch) {
			if PLR_FLAGGED(ch, PLR_OOZARU) {
				dmg += int64(float64(dmg) * 0.3)
			}
		}
		if !IS_NPC(vict) {
			if PLR_FLAGGED(vict, PLR_OOZARU) {
				dmg -= int64(float64(dmg) * 0.3)
			}
		}
		if type_ > -1 {
			if ch.Lastattack != 11 && ch.Lastattack != 39 && ch.Lastattack != 500 && ch.Lastattack < 1000 {
				if handle_combo(ch, vict) > 0 {
					if beforered <= 1 {
						ch.Combo = -1
						ch.Combhits = 0
						send_to_char(ch, libc.CString("@RYou have cut your combo short because you missed your last hit!@n\r\n"))
					} else if ch.Combhits < physical_mastery(ch) {
						dmg += combo_damage(ch, dmg, 0)
						if (ch.Combhits == 10 || ch.Combhits == 20 || ch.Combhits == 30) && (level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 || GET_LEVEL(ch) == 100) {
							var gain int64 = int64(GET_LEVEL(ch) * 1000)
							if GET_SKILL(ch, SKILL_STYLE) >= 100 {
								gain += gain * 2
							} else if GET_SKILL(ch, SKILL_STYLE) >= 80 {
								gain += int64(float64(gain) * 0.4)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 60 {
								gain += int64(float64(gain) * 0.3)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 40 {
								gain += int64(float64(gain) * 0.2)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 20 {
								gain += int64(float64(gain) * 0.1)
							}
							gain_exp(ch, gain)
							send_to_char(ch, libc.CString("@D[@mExp@W: @G%s@D]@n\r\n"), add_commas(gain))
						}
					} else {
						dmg += combo_damage(ch, dmg, 1)
						if (ch.Combhits == 10 || ch.Combhits == 20 || ch.Combhits == 30) && (level_exp(ch, GET_LEVEL(ch)+1)-int(ch.Exp) > 0 || GET_LEVEL(ch) == 100) {
							var gain int64 = int64(GET_LEVEL(ch) * 1000)
							if GET_SKILL(ch, SKILL_STYLE) >= 100 {
								gain += gain * 2
							} else if GET_SKILL(ch, SKILL_STYLE) >= 80 {
								gain += int64(float64(gain) * 0.4)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 60 {
								gain += int64(float64(gain) * 0.3)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 40 {
								gain += int64(float64(gain) * 0.2)
							} else if GET_SKILL(ch, SKILL_STYLE) >= 20 {
								gain += int64(float64(gain) * 0.1)
							}
							gain_exp(ch, gain)
							send_to_char(ch, libc.CString("@D[@mExp@W: @G%s@D]@n\r\n"), add_commas(gain))
						}
						ch.Combo = -1
						ch.Combhits = 0
					}
				}
			} else if ch.Combhits > 0 && ch.Lastattack < 1000 {
				send_to_char(ch, libc.CString("@RYou have cut your combo short because you used the wrong attack!@n\r\n"))
				ch.Combo = -1
				ch.Combhits = 0
			}
		}
		if ch.Lastattack >= 1000 {
			ch.Lastattack -= 1000
		}
		if ch.Preference == PREFERENCE_KI && ch.Charge > 0 {
			dmg -= int64(float64(dmg) * 0.08)
		}
		if AFF_FLAGGED(vict, AFF_SANCTUARY) {
			if GET_SKILL(vict, SKILL_AQUA_BARRIER) != 0 {
				if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Geffect >= 0 && (func() int {
					if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
						return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Sector_type
					}
					return SECT_INSIDE
				}()) != SECT_UNDERWATER {
					dmg = int64(float64(dmg) * 0.85)
				} else {
					dmg = int64(float64(dmg) * 0.75)
				}
			}
			if vict.Barrier-dmg > 0 {
				act(libc.CString("@c$N's@C barrier absorbs the damage!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				var barr [2048]byte
				stdio.Sprintf(&barr[0], "@CYour barrier absorbs the damage! @D[@B%s@D]@n", add_commas(dmg))
				act(&barr[0], TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N's@C barrier absorbs the damage!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Barrier -= dmg
				dmg = 0
			} else if vict.Barrier-dmg <= 0 {
				dmg -= vict.Barrier
				vict.Barrier = 0
				act(libc.CString("@c$N's@C barrier bursts!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@CYour barrier bursts!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N's@C barrier bursts!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_SANCTUARY/32)] &= ^(1 << (int(AFF_SANCTUARY % 32)))
			}
		}
		if AFF_FLAGGED(vict, AFF_FIRESHIELD) && rand_number(1, 200) < GET_SKILL(vict, SKILL_FIRESHIELD) {
			act(libc.CString("@c$N's@C fireshield repels the damage!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@CYour fireshield repels the damage!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N's@C fireshield repels the damage!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			if rand_number(1, 3) == 3 {
				act(libc.CString("@c$N's@C fireshield disappears...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@CYour fireshield disappears...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N's@C fireshield disappears...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Affected_by[int(AFF_FIRESHIELD/32)] &= ^(1 << (int(AFF_FIRESHIELD % 32)))
			}
			dmg = 0
		}
		var conlimit int64 = 2000000000
		if type_ == 0 {
			if vict.Max_hit < conlimit {
				index += (vict.Max_hit / 1500) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*2 {
				index += (vict.Max_hit / 2500) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*3 {
				index += (vict.Max_hit / 3500) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*5 {
				index += (vict.Max_hit / 6000) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*10 {
				index += (vict.Max_hit / 8500) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*15 {
				index += (vict.Max_hit / 10000) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*20 {
				index += (vict.Max_hit / 12500) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*25 {
				index += (vict.Max_hit / 16000) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit < conlimit*30 {
				index += (vict.Max_hit / 22000) * int64(vict.Aff_abils.Con/2)
			} else if vict.Max_hit > conlimit*30 {
				index += (vict.Max_hit / 25000) * int64(vict.Aff_abils.Con/2)
			}
		}
		if IS_NPC(vict) && GET_LEVEL(vict) > 0 {
			index /= 3
		} else if IS_NPC(vict) && GET_LEVEL(vict) < 40 {
			index /= 3
		}
		index += armor_calc(vict, dmg, type_)
		if AFF_FLAGGED(vict, AFF_STONESKIN) {
			if GET_LEVEL(vict) < 20 {
				index += int64(GET_LEVEL(vict) * 250)
			} else if GET_LEVEL(vict) < 30 {
				index += int64(GET_LEVEL(vict) * 500)
			} else if GET_LEVEL(vict) < 50 {
				index += int64(GET_LEVEL(vict) * 1000)
			} else if GET_LEVEL(vict) < 60 {
				index += int64(GET_LEVEL(vict) * 2000)
			} else if GET_LEVEL(vict) < 70 {
				index += int64(GET_LEVEL(vict) * 5000)
			} else if GET_LEVEL(vict) < 90 {
				index += int64(GET_LEVEL(vict) * 10000)
			} else if GET_LEVEL(vict) <= 100 {
				index += int64(GET_LEVEL(vict) * 25000)
			}
		}
		if AFF_FLAGGED(vict, AFF_SHELL) {
			dmg = int64(float64(dmg) * 0.25)
		}
		if AFF_FLAGGED(vict, AFF_WITHER) {
			dmg += int64((float64(dmg) * 0.01) * 20)
		}
		if !IS_NPC(vict) && (vict.Player_specials.Conditions[DRUNK]) > 4 {
			dmg -= int64((float64(dmg) * 0.001) * float64(vict.Player_specials.Conditions[DRUNK]))
		}
		if AFF_FLAGGED(vict, AFF_EARMOR) {
			dmg -= int64(float64(dmg) * 0.1)
		}
		if type_ > 0 {
			advanced_energy(vict, dmg)
			dmg -= int64((float64(dmg) * 0.0005) * float64(vict.Aff_abils.Wis))
		}
		if vict.Race == RACE_MUTANT {
			if type_ <= 0 {
				dmg -= int64(float64(dmg) * 0.3)
			} else if type_ > 0 {
				dmg -= int64(float64(dmg) * 0.25)
			}
		}
		if (vict.Bonuses[BONUS_THICKSKIN]) != 0 {
			if type_ <= 0 {
				dmg -= int64(float64(dmg) * 0.2)
			} else {
				dmg -= int64(float64(dmg) * 0.1)
			}
		} else if (vict.Bonuses[BONUS_THINSKIN]) != 0 {
			if type_ <= 0 {
				dmg += int64(float64(dmg) * 0.2)
			} else {
				dmg += int64(float64(dmg) * 0.1)
			}
		}
		if PLR_FLAGGED(vict, PLR_FURY) {
			dmg -= int64(float64(dmg) * 0.1)
		}
		if vict.Race == RACE_MAJIN {
			if type_ <= 0 {
				dmg -= int64(float64(dmg) * 0.5)
			}
		}
		if vict.Race == RACE_KAI {
			dmg += int64(float64(dmg) * 0.15)
		}
		if ch.Grappling == vict && ch.Grap == 3 {
			dmg += (dmg / 100) * 20
		}
		if vict.Clan != nil && C.strcasecmp(vict.Clan, libc.CString("Heavenly Kaios")) == 0 {
			if vict.Mana >= vict.Max_mana/2 {
				dmg -= (dmg / 100) * 20
				act(libc.CString("@wYou are covered in a pristine @Cglow@w.@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@w$n is covered in a pristine @Cglow@w!@n"), TRUE, vict, nil, nil, TO_ROOM)
			}
		}
		if !IS_NPC(vict) && GET_SKILL(vict, SKILL_ARMOR) != 0 {
			var (
				nanite int = GET_SKILL(vict, SKILL_ARMOR)
				perc   int = rand_number(1, 220)
			)
			if PLR_FLAGGED(vict, PLR_SENSEM) {
				perc = rand_number(1, 176)
			}
			if nanite >= perc {
				if PLR_FLAGGED(vict, PLR_TRANS6) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block MOST of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block MOST of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 50)
				} else if PLR_FLAGGED(vict, PLR_TRANS5) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block some of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block some of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 40)
				} else if PLR_FLAGGED(vict, PLR_TRANS4) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a lot of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a lot of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 30)
				} else if PLR_FLAGGED(vict, PLR_TRANS3) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a good deal of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a good deal of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 25)
				} else if PLR_FLAGGED(vict, PLR_TRANS2) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block some of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block some of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 20)
				} else if PLR_FLAGGED(vict, PLR_TRANS1) {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a bit of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a bit of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 10)
				} else {
					act(libc.CString("@WYour @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a tiny bit of the damage!@n"), TRUE, vict, nil, nil, TO_CHAR)
					act(libc.CString("@W$n's @gn@Ga@Wn@wite @Da@Wr@wm@Do@wr@W reacts in time to block a tiny bit of the damage!@n"), TRUE, vict, nil, nil, TO_ROOM)
					dmg -= int64((float64(dmg) * 0.01) * 5)
				}
			}
		}
		if !AFF_FLAGGED(vict, AFF_KNOCKED) && (vict.Position == POS_SITTING || vict.Position == POS_RESTING) && GET_SKILL(vict, SKILL_ROLL) > axion_dice(0) {
			var rollcost int64 = (vict.Max_hit / 300) * int64(ch.Aff_abils.Str/2)
			if vict.Move >= rollcost {
				act(libc.CString("@GYou roll to your feet in an agile fashion!@n"), TRUE, vict, nil, nil, TO_CHAR)
				act(libc.CString("@G$n rolls to $s feet in an agile fashion!@n"), TRUE, vict, nil, nil, TO_ROOM)
				do_stand(vict, nil, 0, 0)
				vict.Move -= rollcost
			}
		}
		if IS_NPC(vict) {
			hitprcnt_mtrigger(vict)
		}
		if IS_HUMANOID(vict) && !IS_NPC(ch) && IS_NPC(vict) && (!is_sparring(ch) || !is_sparring(vict)) {
			remember(vict, ch)
		}
		if IS_NPC(vict) && vict.Hit > gear_pl(vict)/4 {
			vict.Lasthit = int(ch.Idnum)
		}
		if AFF_FLAGGED(vict, AFF_SLEEP) && rand_number(1, 2) == 2 {
			affect_from_char(vict, SPELL_SLEEP)
			act(libc.CString("@c$N@W seems to be more aware now.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@WYou are no longer so sleepy.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N@W seems to be more aware now.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
		}
		if AFF_FLAGGED(vict, AFF_KNOCKED) && rand_number(1, 12) >= 11 {
			act(libc.CString("@W$n@W is no longer senseless, and wakes up.@n"), FALSE, vict, nil, nil, TO_ROOM)
			send_to_char(vict, libc.CString("You are no longer knocked out, and wake up!@n\r\n"))
			if vict.Player_specials.Carried_by != nil {
				if vict.Player_specials.Carried_by.Alignment > 50 {
					carry_drop(vict.Player_specials.Carried_by, 0)
				} else {
					carry_drop(vict.Player_specials.Carried_by, 1)
				}
			}
			vict.Affected_by[int(AFF_KNOCKED/32)] &= ^(1 << (int(AFF_KNOCKED % 32)))
			vict.Position = POS_SITTING
			if IS_NPC(vict) && rand_number(1, 20) >= 12 {
				act(libc.CString("@W$n@W stands up.@n"), FALSE, vict, nil, nil, TO_ROOM)
				vict.Position = POS_STANDING
			}
		}
		if IS_NPC(ch) {
			if GET_LEVEL(ch) > 10 {
				if dmg-index > 0 {
					dmg -= index
				} else if dmg-index <= 0 && dmg >= 1 {
					dmg = 1
				}
			} else if GET_LEVEL(ch) <= 10 {
				dmg = int64(float64(dmg) * 0.8)
			}
		}
		if !IS_NPC(ch) {
			if dmg >= 1 {
				if (float64(dmg)+float64(dmg)*0.5)-float64(index) <= 0 {
					dmg = 1
				} else if (float64(dmg)+float64(dmg)*0.4)-float64(index) <= 0 {
					dmg = int64(float64(dmg) * 0.04)
				} else if (float64(dmg)+float64(dmg)*0.3)-float64(index) <= 0 {
					dmg = int64(float64(dmg) * 0.08)
				} else if (float64(dmg)+float64(dmg)*0.2)-float64(index) <= 0 {
					dmg = int64(float64(dmg) * 0.12)
				} else if (float64(dmg)+float64(dmg)*0.1)-float64(index) <= 0 {
					dmg = int64(float64(dmg) * 0.16)
				} else if dmg-index <= 0 {
					dmg = int64(float64(dmg) * 0.2)
				} else if float64(dmg-index) > float64(dmg)*0.25 {
					dmg -= index
				} else {
					dmg = int64(float64(dmg) * 0.25)
				}
			}
		}
		if dmg < 1 {
			dmg = 0
		}
		if dmg >= 50 && chance > 0 {
			hurt_limb(ch, vict, chance, limb, dmg)
		}
		if IS_NPC(vict) && float64(dmg) > float64(vict.Max_hit)*0.7 && (ch.Bonuses[BONUS_SADISTIC]) > 0 {
			vict.Exp /= 2
		} else if IS_NPC(vict) && dmg > vict.Hit && float64(vict.Hit) >= float64(vict.Max_hit)*0.5 && (ch.Bonuses[BONUS_SADISTIC]) > 0 {
			vict.Exp /= 2
		}
		if vict.Player_specials.Carrying != nil && float64(dmg) > (float64(gear_pl(vict))*0.01) && rand_number(1, 10) >= 8 {
			carry_drop(vict, 2)
		}
		if vict.Position == POS_SITTING && IS_NPC(vict) && float64(vict.Hit) >= float64(gear_pl(vict))*0.98 {
			do_stand(vict, nil, 0, 0)
		}
		var suppresso int = FALSE
		if is_sparring(ch) && is_sparring(vict) && (vict.Suppressed+vict.Hit)-dmg <= 0 {
			if !IS_NPC(vict) {
				act(libc.CString("@c$N@w falls down unconscious, and you stop sparring with $M.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@w stops sparring with you as you fall unconscious.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N@w falls down unconscious, and @C$n@w stops sparring with $M.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Hit = 1
				if vict.Suppressed > 0 {
					vict.Suppressed = 0
					vict.Suppression = 0
				}
				if vict.Fighting != nil {
					stop_fighting(vict)
				}
				if ch.Fighting != nil {
					stop_fighting(ch)
				}
				vict.Position = POS_SLEEPING
				if !IS_NPC(ch) {
					vict.Affected_by[int(AFF_KNOCKED/32)] |= 1 << (int(AFF_KNOCKED % 32))
				}
			} else {
				act(libc.CString("@c$N@w admits defeat to you, stops sparring, and stumbles away.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@c$N@w admits defeat to $n, stops sparring, and stumbles away.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				solo_gain(ch, vict)
				var rew *obj_data
				var next_rew *obj_data
				var founded int = 0
				for rew = vict.Carrying; rew != nil; rew = next_rew {
					next_rew = rew.Next_content
					if rew != nil {
						obj_from_char(rew)
						obj_to_room(rew, vict.In_room)
						founded = 1
					}
				}
				if founded == 1 {
					act(libc.CString("@c$N@w leaves a reward behind out of respect.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				}
				vict.Hit = 0
				extract_char(vict)
				return
			}
		} else if is_sparring(ch) && (vict.Suppressed+vict.Hit)-dmg <= 0 {
			act(libc.CString("@c$N@w falls down unconscious, and you spare $S life.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
			act(libc.CString("@C$n@w spares your life as you fall unconscious.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
			act(libc.CString("@c$N@w falls down unconscious, and @C$n@w spares $S life.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			vict.Hit = 1
			if vict.Fighting != nil {
				stop_fighting(vict)
			}
			if ch.Fighting != nil {
				stop_fighting(ch)
			}
			vict.Position = POS_SLEEPING
			if !IS_NPC(ch) {
				vict.Affected_by[int(AFF_KNOCKED/32)] |= 1 << (int(AFF_KNOCKED % 32))
			}
		} else if is_sparring(ch) && !is_sparring(vict) && IS_NPC(ch) {
			act(libc.CString("@w$n@w stops sparring!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			ch.Act[int(MOB_SPAR/32)] &= bitvector_t(^(1 << (int(MOB_SPAR % 32))))
		} else if !is_sparring(ch) && is_sparring(vict) && IS_NPC(vict) {
			act(libc.CString("@w$n@w stops sparring!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
			vict.Act[int(MOB_SPAR/32)] &= bitvector_t(^(1 << (int(MOB_SPAR % 32))))
		}
		if vict.Suppressed > 0 && vict.Suppression > 0 {
			if vict.Suppressed > dmg {
				vict.Suppressed -= dmg
				suppresso = TRUE
			} else if vict.Suppressed <= dmg {
				send_to_char(vict, libc.CString("@GYou no longer have any reserve powerlevel suppressed.@n\r\n"))
				dmg -= vict.Suppressed
				vict.Suppressed = 0
				vict.Suppression = 0
			}
		}
		if PLR_FLAGGED(vict, PLR_IMMORTAL) && !is_sparring(ch) && vict.Hit-dmg <= 0 {
			if IN_ARENA(vict) {
				send_to_all(libc.CString("@R%s@r manages to defeat @R%s@r in the Arena!@n\r\n"), GET_NAME(ch), GET_NAME(vict))
				char_from_room(ch)
				char_to_room(ch, real_room(0x45D3))
				look_at_room(ch.In_room, ch, 0)
				char_from_room(vict)
				char_to_room(vict, real_room(0x45D3))
				vict.Hit = 1
				look_at_room(vict.In_room, vict, 0)
				if vict.Fighting != nil {
					stop_fighting(vict)
				}
				if ch.Fighting != nil {
					stop_fighting(ch)
				}
				return
			} else {
				act(libc.CString("@c$N@w disappears right before dying. $N appears to be immortal.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@CYou disappear right before death, having been saved by your immortality.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$N@w disappears right before dying. $N appears to be immortal.@n."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Hit = 1
				vict.Mana = 1
				vict.Move = 1
				if vict.Fighting != nil {
					stop_fighting(vict)
				}
				if ch.Fighting != nil {
					stop_fighting(ch)
				}
				vict.Position = POS_SITTING
				char_from_room(vict)
				if vict.Chclass == CLASS_ROSHI {
					char_to_room(vict, real_room(1130))
				}
				if vict.Chclass == CLASS_KABITO {
					char_to_room(vict, real_room(0x2F42))
				}
				if vict.Chclass == CLASS_NAIL {
					char_to_room(vict, real_room(0x2DA3))
				}
				if vict.Chclass == CLASS_BARDOCK {
					char_to_room(vict, real_room(2268))
				}
				if vict.Chclass == CLASS_KRANE {
					char_to_room(vict, real_room(0x32D1))
				}
				if vict.Chclass == CLASS_TAPION {
					char_to_room(vict, real_room(8231))
				}
				if vict.Chclass == CLASS_PICCOLO {
					char_to_room(vict, real_room(1659))
				}
				if vict.Chclass == CLASS_ANDSIX {
					char_to_room(vict, real_room(1713))
				}
				if vict.Chclass == CLASS_DABURA {
					char_to_room(vict, real_room(6486))
				}
				if vict.Chclass == CLASS_FRIEZA {
					char_to_room(vict, real_room(4282))
				}
				if vict.Chclass == CLASS_GINYU {
					char_to_room(vict, real_room(4289))
				}
				if vict.Chclass == CLASS_JINTO {
					char_to_room(vict, real_room(3499))
				}
				if vict.Chclass == CLASS_TSUNA {
					char_to_room(vict, real_room(15000))
				}
				if vict.Chclass == CLASS_KURZAK {
					char_to_room(vict, real_room(16100))
				}
			}
			return
		}
		if vict.Grappling != nil && vict.Grappling != ch && type_ == 1 {
			act(libc.CString("@YThe attack hurts YOU as well because you are grappling with $M!@n"), TRUE, vict, nil, unsafe.Pointer(vict.Grappling), TO_VICT)
			act(libc.CString("@YThe attack hurts @y$N@Y as well because $n is grappling with $m!@n"), TRUE, vict, nil, unsafe.Pointer(vict.Grappling), TO_NOTVICT)
			maindmg = maindmg / 2
			hurt(0, 0, ch, vict.Grappling, nil, maindmg, 3)
		}
		if vict.Grappled != nil && vict.Grappled != ch && type_ == 1 {
			act(libc.CString("@YThe attack hurts YOU as well because you are being grappled by $M!@n"), TRUE, vict, nil, unsafe.Pointer(vict.Grappled), TO_VICT)
			act(libc.CString("@YThe attack hurts @y$N@Y as well because $n is being grappled by $m!@n"), TRUE, vict, nil, unsafe.Pointer(vict.Grappled), TO_NOTVICT)
			maindmg = maindmg / 2
			hurt(0, 0, ch, vict.Grappled, nil, maindmg, 3)
		}
		if !is_sparring(ch) && !PLR_FLAGGED(vict, PLR_IMMORTAL) && vict.Hit-dmg <= 0 {
			if vict.Hit-dmg <= 0 && suppresso == FALSE {
				vict.Hit = 0
				if !IS_NPC(vict) && vict.Lifeforce-(dmg-vict.Hit) >= 0 {
					act(libc.CString("@c$N@w barely clings to life!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@CYou barely cling to life!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$N@w barely clings to life!@n."), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					var lifeloss int64 = dmg - vict.Hit
					vict.Lifeforce -= lifeloss
					send_to_char(vict, libc.CString("@D[@CLifeforce@D: @R-%s@D]\n"), add_commas(lifeloss))
					if float64(vict.Lifeforce) >= float64(GET_LIFEMAX(vict))*0.05 {
						send_to_char(vict, libc.CString("@YYou recover a bit thanks to your strong life force.@n\r\n"))
						vict.Hit = int64(float64(GET_LIFEMAX(vict)) * 0.05)
						vict.Lifeforce -= int64(float64(GET_LIFEMAX(vict)) * 0.05)
					} else {
						vict.Hit = int64(GET_LEVEL(vict) * 100)
					}
					return
				}
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = 0
				}
				if type_ <= 0 && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					handle_death_msg(ch, vict, 0)
				} else if type_ > 0 && (!IS_NPC(vict) || !MOB_FLAGGED(vict, MOB_DUMMY)) {
					handle_death_msg(ch, vict, 1)
				} else {
					act(libc.CString("@R$N@w self destructs with a mild explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@R$N@w self destructs with a mild explosion!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
				}
				if dmg > 1 {
					if type_ <= 0 && float64(ch.Hit) >= (float64(gear_pl(ch))*0.5) {
						var raise int64 = int64((float64(ch.Max_mana) * 0.005) + 1)
						if ch.Mana+raise < ch.Max_mana {
							ch.Mana += raise
						} else {
							ch.Mana = ch.Max_mana
						}
					}
					send_to_char(ch, libc.CString("@D[@GDamage@W: @R%s@D]@n\r\n"), add_commas(dmg))
					send_to_char(vict, libc.CString("@D[@rDamage@W: @R%s@D]@n\r\n"), add_commas(dmg))
					var healhp int64 = int64(float64(vict.Max_hit) * 0.12)
					if AFF_FLAGGED(ch, AFF_METAMORPH) && ch.Hit <= ch.Max_hit {
						act(libc.CString("@RYour dark aura saps some of @r$N's@R life energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
						act(libc.CString("@r$n@R's dark aura saps some of your life energy!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
						ch.Hit += healhp
					}
					if ch.Race == RACE_MUTANT && ((ch.Genome[0]) == 10 || (ch.Genome[1]) == 10) {
						ch.Mana += int64(float64(dmg) * 0.05)
						if ch.Mana > ch.Max_mana {
							ch.Mana = ch.Max_mana
						}
					}
					if !is_sparring(ch) && IS_NPC(vict) {
						if type_ == 0 && rand_number(1, 100) >= 97 {
							send_to_char(ch, libc.CString("@YY@yo@Yu @yg@Ya@yi@Yn@y s@Yo@ym@Ye @yb@Yo@yn@Yu@ys @Ye@yx@Yp@ye@Yr@yi@Ye@yn@Yc@ye@Y!@n\r\n"))
							var gain int64 = int64(float64(vict.Exp) * 0.05)
							gain += 1
							gain_exp(ch, gain)
						} else if type_ != 0 && rand_number(1, 100) >= 93 {
							var gain int64 = int64(float64(vict.Exp) * 0.05)
							gain += 1
							gain_exp(ch, gain)
						}
					}
					if AFF_FLAGGED(vict, AFF_ECHAINS) {
						if IS_NPC(ch) && type_ == 0 {
							ch.Real_abils.Cha -= 2
							if ch.Real_abils.Cha < 5 {
								ch.Real_abils.Cha = 5
							} else {
								act(libc.CString("@CEthereal chains burn into existence! They quickly latch onto @RYOUR@C body and begin temporarily hampering $s actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
								act(libc.CString("@CEthereal chains burn into existence! They quickly latch onto @c$n's@C body and begin temporarily hampering $s actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
							}
						} else if type_ == 0 {
							WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
							act(libc.CString("@CEthereal chains burn into existence! They quickly latch onto @RYOUR@C body and begin temporarily hampering $s actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("@CEthereal chains burn into existence! They quickly latch onto @c$n's@C body and begin temporarily hampering $s actions!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_ROOM)
						}
					}
				} else if dmg <= 1 {
					send_to_char(ch, libc.CString("@D[@GDamage@W: @BPitiful...@D]@n\r\n"))
					send_to_char(vict, libc.CString("@D[@rDamage@W: @BPitiful...@D]@n\r\n"))
				}
				vict.Hit = 0
				if AFF_FLAGGED(ch, AFF_GROUP) {
					group_gain(ch, vict)
				} else {
					solo_gain(ch, vict)
				}
				if ch.Race == RACE_DEMON && type_ == 1 {
					vict.Affected_by[int(AFF_ASHED/32)] |= 1 << (int(AFF_ASHED % 32))
				}
				die(vict, ch)
				dead = TRUE
			}
		} else if vict.Hit-dmg > 0 || suppresso == TRUE {
			if suppresso == FALSE {
				vict.Hit -= dmg
			}
			if ch.Fighting == nil {
				set_fighting(ch, vict)
			} else if ch.Fighting != vict {
				set_fighting(ch, vict)
			}
			if vict.Fighting == nil {
				set_fighting(vict, ch)
			} else if vict.Fighting != ch {
				set_fighting(vict, ch)
			}
			if dmg > 1 && suppresso == FALSE {
				if type_ == 0 && float64(ch.Hit) >= (float64(gear_pl(ch))*0.5) {
					var raise int64 = int64((float64(ch.Max_mana) * 0.005) + 1)
					if ch.Mana+raise < ch.Max_mana {
						ch.Mana += raise
					} else {
						ch.Mana = ch.Max_mana
					}
				}
				if ch.Race == RACE_MUTANT && ((ch.Genome[0]) == 10 || (ch.Genome[1]) == 10) {
					ch.Mana += int64(float64(dmg) * 0.05)
					if ch.Mana > ch.Max_mana {
						ch.Mana = ch.Max_mana
					}
				}
				send_to_char(ch, libc.CString("@D[@GDamage@W: @R%s@D]@n"), add_commas(dmg))
				send_to_char(vict, libc.CString("@D[@rDamage@W: @R%s@D]@n\r\n"), add_commas(dmg))
				if (ch.Equipment[WEAR_EYE]) != nil && vict != nil && !PRF_FLAGGED(ch, PRF_NODEC) {
					if vict.Race == RACE_ANDROID {
						send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
					} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_BSCOUTER) && vict.Hit >= 150000 {
						send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
					} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_MSCOUTER) && vict.Hit >= 5000000 {
						send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
					} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_ASCOUTER) && vict.Hit >= 15000000 {
						send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
					} else {
						send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c%s@D>@n\r\n"), add_commas(vict.Hit))
					}
				} else {
					send_to_char(ch, libc.CString("\r\n"))
				}
			} else if !IS_NPC(ch) {
				if dmg <= 1 && suppresso == FALSE && !PRF_FLAGGED(ch, PRF_NODEC) {
					send_to_char(ch, libc.CString("@D[@GDamage@W: @BPitiful...@D]@n"))
					send_to_char(vict, libc.CString("@D[@rDamage@W: @BPitiful...@D]@n\r\n"))
					if (ch.Equipment[WEAR_EYE]) != nil && vict != nil && !PRF_FLAGGED(ch, PRF_NODEC) {
						if vict.Race == RACE_ANDROID {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_BSCOUTER) && vict.Hit >= 150000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_MSCOUTER) && vict.Hit >= 5000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_ASCOUTER) && vict.Hit >= 15000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c%s@D>@n\r\n"), add_commas(vict.Hit))
						}
					} else {
						send_to_char(ch, libc.CString("\r\n"))
					}
				} else if dmg > 1 && suppresso == TRUE && !PRF_FLAGGED(ch, PRF_NODEC) {
					send_to_char(ch, libc.CString("@D[@GDamage@W: @R%s@D]@n"), add_commas(dmg))
					send_to_char(vict, libc.CString("@D[@rDamage@W: @R%s @c-Suppression-@D]@n\r\n"), add_commas(dmg))
					send_to_char(vict, libc.CString("@D[Suppression@W: @G%s@D]@n\r\n"), add_commas(vict.Suppressed))
					if (ch.Equipment[WEAR_EYE]) != nil && vict != nil && !PRF_FLAGGED(ch, PRF_NODEC) {
						if vict.Race == RACE_ANDROID {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_BSCOUTER) && vict.Hit >= 150000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_MSCOUTER) && vict.Hit >= 5000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_ASCOUTER) && vict.Hit >= 15000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c%s@D>@n\r\n"), add_commas(vict.Hit))
						}
					} else {
						send_to_char(ch, libc.CString("\r\n"))
					}
				} else if dmg <= 1 && suppresso == TRUE && !PRF_FLAGGED(ch, PRF_NODEC) {
					send_to_char(ch, libc.CString("@D[@GDamage@W: @BPitiful...@D]@n"))
					send_to_char(vict, libc.CString("@D[@rDamage@W: @BPitiful... @c-Suppression-@D]@n\r\n"))
					send_to_char(vict, libc.CString("@D[Suppression@W: @G%s@D]@n\r\n"), add_commas(vict.Suppressed))
					if (ch.Equipment[WEAR_EYE]) != nil && vict != nil {
						if vict.Race == RACE_ANDROID && !PRF_FLAGGED(ch, PRF_NODEC) {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_BSCOUTER) && vict.Hit >= 150000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_MSCOUTER) && vict.Hit >= 5000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else if OBJ_FLAGGED(ch.Equipment[WEAR_EYE], ITEM_ASCOUTER) && vict.Hit >= 15000000 {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c?????????????@D>@n\r\n"))
						} else {
							send_to_char(ch, libc.CString(" @D<@YProcessing@D: @c%s@D>@n\r\n"), add_commas(vict.Hit))
						}
					} else {
						send_to_char(ch, libc.CString("\r\n"))
					}
				}
			}
		}
		if GET_SKILL(ch, SKILL_FOCUS) != 0 && type_ == 1 {
			improve_skill(ch, SKILL_FOCUS, 1)
		}
		if !is_sparring(ch) && vict.Race == RACE_HALFBREED && int(vict.Fury) < 100 && !PLR_FLAGGED(vict, PLR_FURY) {
			send_to_char(vict, libc.CString("@RYour fury increases a little bit!@n\r\n"))
			vict.Fury += 1
		}
		if is_sparring(ch) && is_sparring(vict) && ch.Lastattack != -1 {
			spar_gain(ch, vict, type_, dmg)
		}
		if (ch.Race == RACE_SAIYAN || ch.Race == RACE_BIO && ((ch.Genome[0]) == 2 || (ch.Genome[1]) == 2)) && !IS_NPC(ch) && (is_sparring(ch) && is_sparring(vict) || !is_sparring(ch) && !is_sparring(vict)) {
			if ch.Position != POS_RESTING && vict.Position != POS_RESTING && dmg > 1 {
				saiyan_gain(ch, vict)
			}
		}
		if vict.Race == RACE_ARLIAN && dead != TRUE && !is_sparring(vict) && !is_sparring(ch) {
			handle_evolution(vict, dmg)
		}
		if dead == TRUE {
			var corp [256]byte
			if !PLR_FLAGGED(ch, PLR_SELFD2) {
				if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOGOLD) {
					stdio.Sprintf(&corp[0], "all.zenni corpse")
					do_get(ch, &corp[0], 0, 0)
				}
				if !IS_NPC(ch) && ch != vict && PRF_FLAGGED(ch, PRF_AUTOLOOT) {
					stdio.Sprintf(&corp[0], "all corpse")
					do_get(ch, &corp[0], 0, 0)
				}
			}
		}
	} else if obj != nil {
		switch obj.Value[VAL_ALL_MATERIAL] {
		case MATERIAL_STEEL:
			dmg = dmg / 4
		case MATERIAL_MITHRIL:
			dmg = dmg / 6
		case MATERIAL_IRON:
			dmg = dmg / 3
		case MATERIAL_STONE:
			dmg = dmg / 2
		case MATERIAL_DIAMOND:
			dmg = dmg / 10
		}
		if dmg <= 0 {
			dmg = 1
		}
		if OBJ_FLAGGED(obj, ITEM_UNBREAKABLE) {
			act(libc.CString("$p@w seems unaffected.@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$p@w seems unaffected.@n"), TRUE, ch, obj, nil, TO_ROOM)
		} else if GET_OBJ_VNUM(obj) == 79 {
			if obj.Weight-dmg > 0 {
				if type_ <= 0 {
					if AFF_FLAGGED(ch, AFF_INFUSE) {
						dmg *= 10
					}
					act(libc.CString("$p@w cracks some.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("$p@w cracks some.@n"), TRUE, ch, obj, nil, TO_ROOM)
					obj.Weight -= dmg
					if obj.Fellow_wall != nil {
						var wall *obj_data
						wall = obj.Fellow_wall
						wall.Weight -= dmg
						act(libc.CString("$p@w cracks some. A humanoid shadow can be seen moving on the other side.@n"), TRUE, nil, obj, nil, TO_ROOM)
					}
				} else {
					dmg *= 30
					act(libc.CString("$p@w melts some.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("$p@w melts some.@n"), TRUE, ch, obj, nil, TO_ROOM)
					obj.Weight -= dmg
					if obj.Fellow_wall != nil {
						var wall *obj_data
						wall = obj.Fellow_wall
						wall.Weight -= dmg
						act(libc.CString("$p@w melts some.@n"), TRUE, ch, obj, nil, TO_ROOM)
					}
				}
			} else {
				if type_ <= 0 {
					act(libc.CString("$p@w breaks completely apart and then melts away.@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("$p@w breaks completely apart and then melts away.@n"), TRUE, ch, obj, nil, TO_ROOM)
					extract_obj(obj)
				} else {
					act(libc.CString("$p@w is blown away into snow and water!@n"), TRUE, ch, obj, nil, TO_CHAR)
					act(libc.CString("$p@w is blown away into snow and water!@n"), TRUE, ch, obj, nil, TO_ROOM)
					extract_obj(obj)
				}
			}
		} else if (obj.Value[VAL_ALL_HEALTH])-int(dmg) > 0 {
			act(libc.CString("$p@w cracks some.@n"), TRUE, ch, obj, nil, TO_CHAR)
			act(libc.CString("$p@w cracks some.@n"), TRUE, ch, obj, nil, TO_ROOM)
			obj.Value[VAL_ALL_HEALTH] -= int(dmg)
		} else {
			if type_ <= 0 {
				act(libc.CString("$p@w shatters apart!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("$p@w shatters apart!@n"), TRUE, ch, obj, nil, TO_ROOM)
				obj.Value[VAL_ALL_HEALTH] = 0
				obj.Extra_flags[int(ITEM_BROKEN/32)] |= bitvector_t(1 << (int(ITEM_BROKEN % 32)))
				if obj.Type_flag == ITEM_DRINKCON && obj.Type_flag == ITEM_FOUNTAIN {
					obj.Value[VAL_DRINKCON_HOWFULL] = 0
				}
			} else if type_ != 0 {
				act(libc.CString("$p@w is disintegrated!@n"), TRUE, ch, obj, nil, TO_CHAR)
				act(libc.CString("$p@w is disintegrated!@n"), TRUE, ch, obj, nil, TO_ROOM)
				extract_obj(obj)
			}
		}
	} else {
		basic_mud_log(libc.CString("Log: Error with hurt.\n"))
	}
}
func handle_cooldown(ch *char_data, cooldown int) {
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_MULTIHIT) {
			ch.Act[int(PLR_MULTIHIT/32)] &= bitvector_t(^(1 << (int(PLR_MULTIHIT % 32))))
			return
		}
	}
	if IS_NPC(ch) {
		ch.Cooldown = 0
	}
	reveal_hiding(ch, 0)
	var waitCalc int = 10
	var base int = cooldown
	var cspd int64 = 0
	cspd = int64(GET_SPEEDI(ch))
	if cspd > 10000000 {
		waitCalc -= 9
	} else if cspd > 5000000 {
		waitCalc -= 8
	} else if cspd > 2500000 {
		waitCalc -= 7
	} else if cspd > 1000000 {
		waitCalc -= 6
	} else if cspd > 500000 {
		waitCalc -= 5
	} else if cspd > 100000 {
		waitCalc -= 4
	} else if cspd > 50000 {
		waitCalc -= 3
	} else if cspd > 25000 {
		waitCalc -= 2
	} else if cspd > 50 {
		waitCalc -= 1
	}
	base *= 10
	if base >= 100 {
		base = 30
	} else if base >= 70 {
		base = 20
	} else if base >= 30 {
		base = 10
	}
	if !IS_NPC(ch) {
		cooldown *= waitCalc
		cooldown += base
		if cooldown <= 0 {
			cooldown = 10
		}
		if cooldown >= 120 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*12)
		} else if cooldown >= 110 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*11)
		} else if cooldown >= 100 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*10)
		} else if cooldown >= 90 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*9)
		} else if cooldown >= 80 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*8)
		} else if cooldown >= 70 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*7)
		} else if cooldown >= 60 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*6)
		} else if cooldown >= 50 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*5)
		} else if cooldown >= 40 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*4)
		} else if cooldown >= 30 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*3)
		} else if cooldown >= 20 {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*2)
		} else {
			WAIT_STATE(ch, (int(1000000/OPT_USEC))*1)
		}
	} else {
		cooldown *= waitCalc
		cooldown += base
		if cooldown >= 120 {
			ch.Cooldown = 12
		} else if cooldown >= 110 {
			ch.Cooldown = 11
		} else if cooldown >= 100 {
			ch.Cooldown = 10
		} else if cooldown >= 90 {
			ch.Cooldown = 9
		} else if cooldown >= 80 {
			ch.Cooldown = 8
		} else if cooldown >= 70 {
			ch.Cooldown = 7
		} else if cooldown >= 60 {
			ch.Cooldown = 6
		} else if cooldown >= 50 {
			ch.Cooldown = 5
		} else if cooldown >= 40 {
			ch.Cooldown = 4
		} else if cooldown >= 30 {
			ch.Cooldown = 3
		} else if cooldown >= 20 {
			ch.Cooldown = 2
		} else {
			ch.Cooldown = 1
		}
	}
}
func handle_parry(ch *char_data) int {
	if axion_dice(0) <= 4 {
		return 1
	}
	if IS_NPC(ch) {
		if !IS_HUMANOID(ch) {
			return rand_number(0, 5)
		} else {
			if GET_LEVEL(ch) >= 110 {
				return rand_number(90, 105)
			} else if GET_LEVEL(ch) >= 100 {
				return rand_number(85, 95)
			} else if GET_LEVEL(ch) >= 80 {
				return rand_number(70, 85)
			} else if GET_LEVEL(ch) >= 40 {
				return rand_number(50, 70)
			} else {
				var top int = GET_LEVEL(ch)
				if top < 15 {
					top = 16
				}
				return rand_number(15, top)
			}
		}
	}
	if PRF_FLAGGED(ch, PRF_NOPARRY) {
		return -2
	} else {
		var num int = GET_SKILL(ch, SKILL_PARRY)
		if ch.Race == RACE_MUTANT && ((ch.Genome[0]) == 3 || (ch.Genome[1]) == 3) {
			num += 10
		}
		if (ch.Skills[SKILL_STYLE]) >= 100 {
			num += 5
		} else if (ch.Skills[SKILL_STYLE]) >= 80 {
			num += 4
		} else if (ch.Skills[SKILL_STYLE]) >= 60 {
			num += 3
		} else if (ch.Skills[SKILL_STYLE]) >= 40 {
			num += 2
		} else if (ch.Skills[SKILL_STYLE]) >= 20 {
			num += 1
		}
		return num
	}
}
func handle_combo(ch *char_data, vict *char_data) int {
	var (
		success int = FALSE
		pass    int = FALSE
	)
	if IS_NPC(ch) {
		return 0
	}
	switch ch.Lastattack {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 8:
		fallthrough
	case 51:
		fallthrough
	case 52:
		fallthrough
	case 56:
		pass = TRUE
	default:
		if ch.Combo != -1 {
			send_to_char(ch, libc.CString("@RYou have cut your combo short with the wrong attack!@n\r\n"))
		}
		ch.Combo = -1
		ch.Combhits = 0
		pass = FALSE
	}
	if pass == FALSE {
		return 0
	}
	if count_physical(ch) < 3 {
		return 0
	}
	var chance int = int((ch.Aff_abils.Cha * 2) - vict.Aff_abils.Cha)
	var bottom int = chance / 2
	if ch.Lastattack == 0 || ch.Lastattack == 1 {
		chance += 10
	}
	if ch.Lastattack == 2 || ch.Lastattack == 3 {
		chance += 5
	}
	if chance > 110 {
		chance = 110
	}
	if ch.Combo <= -1 && rand_number(bottom, 125) < chance {
		for success == FALSE {
			switch rand_number(1, 24) {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				fallthrough
			case 5:
				if GET_SKILL(ch, SKILL_PUNCH) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R punch @Gnext!@n\r\n"))
					ch.Combo = 0
					success = TRUE
				}
			case 6:
				fallthrough
			case 7:
				fallthrough
			case 8:
				fallthrough
			case 9:
				fallthrough
			case 10:
				if GET_SKILL(ch, SKILL_KICK) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R kick @Gnext!@n\r\n"))
					ch.Combo = 1
					success = TRUE
				}
			case 11:
				fallthrough
			case 12:
				fallthrough
			case 13:
				fallthrough
			case 14:
				if GET_SKILL(ch, SKILL_ELBOW) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try an@R elbow @Gnext!@n\r\n"))
					ch.Combo = 2
					success = TRUE
				}
			case 15:
				fallthrough
			case 16:
				fallthrough
			case 17:
				if GET_SKILL(ch, SKILL_KNEE) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R knee @Gnext!@n\r\n"))
					ch.Combo = 3
					success = TRUE
				}
			case 18:
				fallthrough
			case 19:
				if GET_SKILL(ch, SKILL_ROUNDHOUSE) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R roundhouse @Gnext!@n\r\n"))
					ch.Combo = 4
					success = TRUE
				}
			case 20:
				fallthrough
			case 21:
				if GET_SKILL(ch, SKILL_UPPERCUT) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try an@R uppercut @Gnext!@n\r\n"))
					ch.Combo = 5
					success = TRUE
				}
			case 22:
				if GET_SKILL(ch, SKILL_HEELDROP) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R heeldrop @Gnext!@n\r\n"))
					ch.Combo = 8
					success = TRUE
				}
			case 24:
				if GET_SKILL(ch, SKILL_SLAM) > 0 {
					send_to_char(ch, libc.CString("@GYou have a chance for a COMBO! Try a@R slam @Gnext!@n\r\n"))
					ch.Combo = 6
					success = TRUE
				}
			}
		}
		return 0
	} else if ch.Lastattack == ch.Combo && ch.Combhits < physical_mastery(ch) {
		ch.Combhits += 1
		for success == FALSE {
			if ch.Combhits >= 20 {
				switch rand_number(1, 34) {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					fallthrough
				case 5:
					fallthrough
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					if GET_SKILL(ch, SKILL_ELBOW) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R elbow@G!@n\r\n"), ch.Combhits)
						ch.Combo = 2
						success = TRUE
					}
				case 9:
					fallthrough
				case 10:
					fallthrough
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					if GET_SKILL(ch, SKILL_KNEE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rknee@G!@n\r\n"), ch.Combhits)
						ch.Combo = 3
						success = TRUE
					}
				case 17:
					fallthrough
				case 18:
					fallthrough
				case 19:
					fallthrough
				case 20:
					fallthrough
				case 21:
					if GET_SKILL(ch, SKILL_UPPERCUT) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R uppercut@G!@n\r\n"), ch.Combhits)
						ch.Combo = 5
						success = TRUE
					}
				case 22:
					fallthrough
				case 23:
					fallthrough
				case 24:
					fallthrough
				case 25:
					fallthrough
				case 26:
					if GET_SKILL(ch, SKILL_ROUNDHOUSE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rroundhouse@G!@n\r\n"), ch.Combhits)
						ch.Combo = 4
						success = TRUE
					}
				case 27:
					fallthrough
				case 28:
					fallthrough
				case 29:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_SLAM) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rslam@G!@n\r\n"), ch.Combhits)
						ch.Combo = 6
						success = TRUE
					}
				case 30:
					fallthrough
				case 31:
					fallthrough
				case 32:
					fallthrough
				case 33:
					fallthrough
				case 34:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEELDROP) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheeldrop@G!@n\r\n"), ch.Combhits)
						ch.Combo = 8
						success = TRUE
					}
				}
			} else if ch.Combhits >= 15 {
				switch rand_number(1, 36) {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					fallthrough
				case 5:
					fallthrough
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					fallthrough
				case 9:
					fallthrough
				case 10:
					if GET_SKILL(ch, SKILL_ELBOW) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R elbow@G!@n\r\n"), ch.Combhits)
						ch.Combo = 2
						success = TRUE
					}
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					fallthrough
				case 18:
					fallthrough
				case 19:
					fallthrough
				case 20:
					if GET_SKILL(ch, SKILL_KNEE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rknee@G!@n\r\n"), ch.Combhits)
						ch.Combo = 3
						success = TRUE
					}
				case 21:
					fallthrough
				case 22:
					fallthrough
				case 23:
					if GET_SKILL(ch, SKILL_PUNCH) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rpunch@G!@n\r\n"), ch.Combhits)
						ch.Combo = 0
						success = TRUE
					}
				case 25:
					fallthrough
				case 26:
					fallthrough
				case 27:
					if GET_SKILL(ch, SKILL_KICK) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rkick@G!@n\r\n"), ch.Combhits)
						ch.Combo = 1
						success = TRUE
					}
				case 29:
					fallthrough
				case 30:
					if GET_SKILL(ch, SKILL_UPPERCUT) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R uppercut@G!@n\r\n"), ch.Combhits)
						ch.Combo = 5
						success = TRUE
					}
				case 31:
					fallthrough
				case 32:
					fallthrough
				case 33:
					fallthrough
				case 34:
					if GET_SKILL(ch, SKILL_ROUNDHOUSE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rroundhouse@G!@n\r\n"), ch.Combhits)
						ch.Combo = 4
						success = TRUE
					}
				case 35:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_SLAM) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rslam@G!@n\r\n"), ch.Combhits)
						ch.Combo = 6
						success = TRUE
					}
				case 36:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEELDROP) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheeldrop@G!@n\r\n"), ch.Combhits)
						ch.Combo = 8
						success = TRUE
					}
				}
			} else if ch.Combhits >= 10 {
				switch rand_number(1, 34) {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					fallthrough
				case 5:
					if GET_SKILL(ch, SKILL_ELBOW) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R elbow@G!@n\r\n"), ch.Combhits)
						ch.Combo = 2
						success = TRUE
					}
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					fallthrough
				case 9:
					fallthrough
				case 10:
					if GET_SKILL(ch, SKILL_KNEE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rknee@G!@n\r\n"), ch.Combhits)
						ch.Combo = 3
						success = TRUE
					}
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					fallthrough
				case 18:
					if GET_SKILL(ch, SKILL_PUNCH) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rpunch@G!@n\r\n"), ch.Combhits)
						ch.Combo = 0
						success = TRUE
					}
				case 19:
					fallthrough
				case 20:
					fallthrough
				case 21:
					fallthrough
				case 22:
					fallthrough
				case 23:
					fallthrough
				case 24:
					fallthrough
				case 25:
					fallthrough
				case 26:
					if GET_SKILL(ch, SKILL_KICK) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rkick@G!@n\r\n"), ch.Combhits)
						ch.Combo = 1
						success = TRUE
					}
				case 27:
					fallthrough
				case 28:
					fallthrough
				case 29:
					if GET_SKILL(ch, SKILL_UPPERCUT) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R uppercut@G!@n\r\n"), ch.Combhits)
						ch.Combo = 5
						success = TRUE
					}
				case 30:
					fallthrough
				case 31:
					if GET_SKILL(ch, SKILL_ROUNDHOUSE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rroundhouse@G!@n\r\n"), ch.Combhits)
						ch.Combo = 4
						success = TRUE
					}
				case 32:
					fallthrough
				case 33:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_SLAM) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rslam@G!@n\r\n"), ch.Combhits)
						ch.Combo = 6
						success = TRUE
					}
				case 34:
					if GET_SKILL(ch, SKILL_BASH) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try bash@G!@n\r\n"), ch.Combhits)
						ch.Combo = 51
						success = TRUE
					} else if GET_SKILL(ch, SKILL_TAILWHIP) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rtailwhip@G!@n\r\n"), ch.Combhits)
						ch.Combo = 56
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEADBUTT) > 0 && rand_number(1, 2) == 2 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheadbutt@G!@n\r\n"), ch.Combhits)
						ch.Combo = 52
						success = TRUE
					} else if GET_SKILL(ch, SKILL_HEELDROP) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rheeldrop@G!@n\r\n"), ch.Combhits)
						ch.Combo = 8
						success = TRUE
					}
				}
			} else if ch.Combhits >= 5 {
				switch rand_number(1, 30) {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					if GET_SKILL(ch, SKILL_ELBOW) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R elbow@G!@n\r\n"), ch.Combhits)
						ch.Combo = 2
						success = TRUE
					}
				case 5:
					fallthrough
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					if GET_SKILL(ch, SKILL_KNEE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rknee@G!@n\r\n"), ch.Combhits)
						ch.Combo = 3
						success = TRUE
					}
				case 9:
					fallthrough
				case 10:
					fallthrough
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					fallthrough
				case 18:
					if GET_SKILL(ch, SKILL_PUNCH) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rpunch@G!@n\r\n"), ch.Combhits)
						ch.Combo = 0
						success = TRUE
					}
				case 19:
					fallthrough
				case 20:
					fallthrough
				case 21:
					fallthrough
				case 22:
					fallthrough
				case 23:
					fallthrough
				case 24:
					fallthrough
				case 25:
					fallthrough
				case 26:
					fallthrough
				case 27:
					fallthrough
				case 28:
					if GET_SKILL(ch, SKILL_KICK) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rkick@G!@n\r\n"), ch.Combhits)
						ch.Combo = 1
						success = TRUE
					}
				case 29:
					if GET_SKILL(ch, SKILL_UPPERCUT) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R uppercut@G!@n\r\n"), ch.Combhits)
						ch.Combo = 5
						success = TRUE
					}
				case 30:
					if GET_SKILL(ch, SKILL_ROUNDHOUSE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rroundhouse@G!@n\r\n"), ch.Combhits)
						ch.Combo = 4
						success = TRUE
					}
				}
			} else {
				switch rand_number(1, 30) {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					if GET_SKILL(ch, SKILL_ELBOW) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try an@R elbow@G!@n\r\n"), ch.Combhits)
						ch.Combo = 2
						success = TRUE
					}
				case 4:
					fallthrough
				case 5:
					fallthrough
				case 6:
					if GET_SKILL(ch, SKILL_KNEE) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rknee@G!@n\r\n"), ch.Combhits)
						ch.Combo = 3
						success = TRUE
					}
				case 7:
					fallthrough
				case 8:
					fallthrough
				case 9:
					fallthrough
				case 10:
					fallthrough
				case 11:
					fallthrough
				case 12:
					fallthrough
				case 13:
					fallthrough
				case 14:
					fallthrough
				case 15:
					fallthrough
				case 16:
					fallthrough
				case 17:
					fallthrough
				case 18:
					if GET_SKILL(ch, SKILL_PUNCH) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rpunch@G!@n\r\n"), ch.Combhits)
						ch.Combo = 0
						success = TRUE
					}
				case 19:
					fallthrough
				case 20:
					fallthrough
				case 21:
					fallthrough
				case 22:
					fallthrough
				case 23:
					fallthrough
				case 24:
					fallthrough
				case 25:
					fallthrough
				case 26:
					fallthrough
				case 27:
					fallthrough
				case 28:
					fallthrough
				case 29:
					fallthrough
				case 30:
					if GET_SKILL(ch, SKILL_KICK) > 0 {
						send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Next try a @Rkick@G!@n\r\n"), ch.Combhits)
						ch.Combo = 1
						success = TRUE
					}
				}
			}
		}
		return ch.Combhits
	} else if ch.Lastattack == ch.Combo && ch.Combhits >= physical_mastery(ch) {
		ch.Combhits += 1
		send_to_char(ch, libc.CString("@D(@GC-c-combo Bonus @gx%d@G!@D)@C Combo FINISHED for massive damage@G!@n\r\n"), ch.Combhits)
	} else if ch.Combo != ch.Lastattack && ch.Combo > -1 {
		send_to_char(ch, libc.CString("@GCombo failed! Try harder next time!@n\r\n"))
		ch.Combo = -1
		ch.Combhits = 0
		return 0
	} else {
		ch.Combo = -1
		ch.Combhits = 0
		return 0
	}
	return 0
}
func handle_spiral(ch *char_data, vict *char_data, skill int, first int) {
	var (
		prob   int
		perc   int
		avo    int
		index  int
		pry    int = 2
		dge    int = 2
		blk    int = 2
		dmg    int64
		amount float64 = 0.0
	)
	if first == FALSE {
		amount = 0.05
	} else {
		amount = 0.5
	}
	if vict == nil && ch.Fighting != nil {
		vict = ch.Fighting
	} else if vict == nil {
		act(libc.CString("@WHaving lost your target you slow down until your vortex disappears, and end your attack.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W slows down until $s vortex disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_SPIRAL/32)] &= bitvector_t(^(1 << (int(PLR_SPIRAL % 32))))
		return
	}
	if ch.Charge <= 0 {
		act(libc.CString("@WHaving no more charged ki you slow down until your vortex disappears, and end your attack.@n"), TRUE, ch, nil, nil, TO_CHAR)
		act(libc.CString("@C$n@W slows down until $s vortex disappears.@n"), TRUE, ch, nil, nil, TO_ROOM)
		ch.Act[int(PLR_SPIRAL/32)] &= bitvector_t(^(1 << (int(PLR_SPIRAL % 32))))
		return
	}
	if vict != nil {
		index = check_def(vict)
		prob = skill
		perc = axion_dice(0)
		index -= handle_speed(ch, vict)
		avo = index / 4
		handle_defense(vict, &pry, &blk, &dge)
		if avo > 0 && avo < 70 {
			prob -= avo
		} else if avo >= 70 {
			prob -= 69
		}
		if vict.Position == POS_SLEEPING {
			pry = 0
			blk = 0
			dge = 0
			prob += 50
		}
		if vict.Position == POS_RESTING {
			pry /= 4
			blk /= 4
			dge /= 4
			prob += 25
		}
		if vict.Position == POS_SITTING {
			pry /= 2
			blk /= 2
			dge /= 2
			prob += 10
		}
		if (!IS_NPC(vict) && vict.Race == RACE_ICER && rand_number(1, 30) >= 28 || AFF_FLAGGED(vict, AFF_ZANZOKEN)) && vict.Move >= 1 && vict.Position != POS_SLEEPING {
			if !AFF_FLAGGED(ch, AFF_ZANZOKEN) || AFF_FLAGGED(ch, AFF_ZANZOKEN) && GET_SPEEDI(ch)+rand_number(1, 5) < GET_SPEEDI(vict)+rand_number(1, 5) {
				act(libc.CString("@C$N@c disappears, avoiding your Spiral Comet blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@cYou disappear, avoiding @C$n's@c Spiral Comet blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$N@c disappears, avoiding @C$n's@c Spiral Comet blast before reappearing!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if AFF_FLAGGED(ch, AFF_ZANZOKEN) {
					ch.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				}
				vict.Affected_by[int(AFF_ZANZOKEN/32)] &= ^(1 << (int(AFF_ZANZOKEN % 32)))
				pcost(ch, amount, 0)
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
		if prob < perc {
			if vict.Move > 0 {
				if blk > rand_number(1, 130) {
					act(libc.CString("@C$N@W moves quickly and blocks your Spiral Comet blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou move quickly and block @C$n's@W Spiral Comet blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W moves quickly and blocks @c$n's@W Spiral Comet blast!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, amount, 0)
					dmg = damtype(ch, 10, skill, 0.05)
					dmg /= 4
					hurt(0, 0, ch, vict, nil, dmg, 1)
					return
				} else if dge > rand_number(1, 130) {
					act(libc.CString("@C$N@W manages to dodge your Spiral Comet blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@WYou dodge @C$n's@W Spiral Comet blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@C$N@W manages to dodge @c$n's@W Spiral Comet blast, letting it slam into the surroundings!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					send_to_room(vict.In_room, libc.CString("@wA bright explosion erupts from the impact!\r\n"))
					dodge_ki(ch, vict, 0, 45, skill, SKILL_SPIRAL)
					if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg <= 95 {
						(*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Dmg += 5
					}
					pcost(ch, amount, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				} else {
					act(libc.CString("@WYou can't believe it but your Spiral Comet blast misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
					act(libc.CString("@C$n@W fires a Spiral Comet blast at you, but misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
					act(libc.CString("@c$n@W fires a Spiral Comet blast at @C$N@W, but somehow misses!@n "), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
					pcost(ch, amount, 0)
					hurt(0, 0, ch, vict, nil, 0, 1)
					return
				}
			} else {
				act(libc.CString("@WYou can't believe it but your Spiral Comet blast misses, flying through the air harmlessly!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W fires a Spiral Comet blast at you, but misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@c$n@W fires a Spiral Comet blast at @C$N@W, but somehow misses!@n"), FALSE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				pcost(ch, amount, 0)
			}
			hurt(0, 0, ch, vict, nil, 0, 1)
			return
		} else {
			if first == TRUE {
				dmg = damtype(ch, 44, skill, 0.5)
			} else {
				dmg = damtype(ch, 45, skill, 0.01)
			}
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@WYou launch a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S chest and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at YOU! It slams into YOUR chest and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S chest and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 2:
				act(libc.CString("@WYou launch a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S head and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at YOU! It slams into YOUR head and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S head and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg *= 2
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 3)
			case 3:
				act(libc.CString("@WYou launch a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at YOU! It slams into YOUR body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S body and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				hurt(0, 0, ch, vict, nil, dmg, 1)
				dam_eq_loc(vict, 4)
			case 4:
				act(libc.CString("@WYou launch a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at YOU! It slams into YOUR arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S arm and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg /= 2
				dam_eq_loc(vict, 1)
				hurt(0, 190, ch, vict, nil, dmg, 1)
			case 5:
				act(libc.CString("@WYou launch a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at YOU! It slams into YOUR leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@C$n@W launches a bright @mp@Mu@mr@Mp@ml@Me@W ball of energy down at @c$N@W! It slams into $S leg and explodes!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				dmg /= 2
				dam_eq_loc(vict, 2)
				hurt(1, 190, ch, vict, nil, dmg, 1)
			}
			pcost(ch, amount, 0)
			return
		}
	} else {
		return
	}
}
func handle_death_msg(ch *char_data, vict *char_data, type_ int) {
	if type_ == 0 {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect >= 0 && (func() int {
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
		}()) != SECT_WATER_NOSWIM && !ROOM_FLAGGED(vict.In_room, ROOM_SPACE) && (func() int {
			if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_FLYING {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r coughs up blood before falling to the ground dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cough up blood before falling to the ground dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r coughs up blood before falling down dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 2:
				act(libc.CString("@R$N@r crumples to the ground dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou crumple to the ground dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r crumples to the ground dead.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 3:
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cry out your last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 4:
				act(libc.CString("@R$N@r writhes on the ground screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou writhe on the ground screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r writhes on the ground screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			case 5:
				act(libc.CString("@R$N@r hits the ground dead with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou hit the ground dead with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r hits the ground dead with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			}
		} else if (func() int {
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
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r coughs up blood and dies before falling down to the water. A large splash accompanies $S body hitting the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cough up blood and die before falling down to the water. A large splash accompanies your body hitting the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r coughs up blood and dies before falling down to the water. A large splash accompanies $S body hitting the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 2:
				act(libc.CString("@R$N@r crumples down to the water, with the signs of life leaving $S eyes as $E floats in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou crumple down to the water and die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r crumples down to the water, with the signs of life leaving $S eyes as $E floats in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 3:
				act(libc.CString("@R$N@r cries out $S last breath before dying and leaving a floating corpse behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cry out your last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r cries out $S last breath before dying and leaving a floating corpse behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 4:
				act(libc.CString("@R$N@r writhes in the water screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou writhe in the water screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r writhes in the water screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			case 5:
				act(libc.CString("@R$N@r hits the water dead with such force that blood mixed with water flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou hit the water dead with such force that blood mixed with water flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r hits the water dead with such force that blood mixed with water flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			}
		} else if ROOM_FLAGGED(vict.In_room, ROOM_SPACE) {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r coughs up blood and dies. The blood freezes and floats freely through space...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cough up blood and die. The blood freezes and floats freely through space...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r coughs up blood and dies. The blood freezes and floats freely through space...@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 2:
				act(libc.CString("@R$N@r dies and leaves $S corpse floating freely in space.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou die and leave your corpse floating freely in space.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r dies and leaves $S corpse floating freely in space.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 3:
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cry out your last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 4:
				act(libc.CString("@R$N@r writhes in space trying to scream in pain before $e finally dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou writhe in space trying to scream in pain before you finally die!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r writhes in space trying to scream in pain before $e finally dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			case 5:
				act(libc.CString("@R$N@r dies suddenly leaving behind a badly damaged corpse floating in space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou die suddenly leaving behind a badly damaged corpse floating in space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r dies suddenly leaving behind a badly damaged corpse floating in space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			}
		} else if (func() int {
			if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FLYING {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r coughs up blood before $s corpse starts to fall to the ground far below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou coughs up blood before your corpse starts to fall to the ground far below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r coughs up blood before $s corpse starts to fall to the ground far below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 2:
				act(libc.CString("@R$N@r dies and $S corpse begins to fall to the ground below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou die and your corpse begins to fall to the ground below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r dies and $S corpse begins to fall to the ground below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 3:
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cry out your last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 4:
				act(libc.CString("@R$N@r writhes in midair screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou writhe in midair screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r writhes in midair screaming in pain before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			case 5:
				act(libc.CString("@R$N@r snaps back and dies with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou snap back and die with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r hits the ground dead with such force that blood flies into the air briefly!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			}
		} else {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r coughs up blood before $s corpse starts to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou coughs up blood before your corpse starts to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r coughs up blood before $s corpse starts to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 2:
				act(libc.CString("@R$N@r dies and $S corpse begins to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou die and your corpse begins to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r dies and $S corpse begins to float limply in the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 3:
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou cry out your last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r cries out $S last breath before dying.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 4:
				act(libc.CString("@R$N@r writhes and thrases in the water trying to scream before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou writhe and thrash in the water trying to scream before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r writhes and thrashes in the water trying to scream before finally dying!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			case 5:
				act(libc.CString("@R$N@r snaps back and dies with such force that blood floods out of $S body into the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou snap back and die with such force that blood floods out of your body into the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r hits the ground dead with such force that blood floods out of $S body into the water!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				if vict.Death_type != DTYPE_HEAD {
					vict.Death_type = DTYPE_PULP
				}
			}
		}
	} else {
		if (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Geffect >= 0 && (func() int {
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
		}()) != SECT_WATER_NOSWIM && !ROOM_FLAGGED(vict.In_room, ROOM_SPACE) && (func() int {
			if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) != SECT_FLYING {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r explodes and chunks of $M shower to the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou explode leaving only chunks behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r explodes and chunks of $M shower to the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 2:
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rThe bottom half of your body is all that remains as you die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_HALF
			case 3:
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body completely disintegrates in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 4:
				act(libc.CString("@R$N@r falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 5:
				act(libc.CString("@rWhat's left of @R$N@r's body slams into the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rWhat's left of your body slams into the ground as you die!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rWhat's left of @R$N@r's body slams into the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else if (func() int {
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
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r explodes and chunks of $M shower to the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou explode leaving only chunks behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r explodes and chunks of $M shower to the ground.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 2:
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rThe bottom half of your body is all that remains as you die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_HALF
			case 3:
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body completely disintegrates in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 4:
				act(libc.CString("@R$N@r falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r falls down as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 5:
				act(libc.CString("@rWhat's left of @R$N@r's body slams into the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rWhat's left of your body slams into the ground as you die!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rWhat's left of @R$N@r's body slams into the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else if ROOM_FLAGGED(vict.In_room, ROOM_SPACE) {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r explodes and chunks of $M shower out into every direction of space.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou explode leaving only chunks behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r explodes and chunks of $M shower out into every direction of space.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 2:
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rThe bottom half of your body is all that remains as you die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_HALF
			case 3:
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body completely disintegrates in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 4:
				act(libc.CString("@R$N@r floats away as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body floats away as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r floats away as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 5:
				act(libc.CString("@rWhat's left of @R$N@r's body floats away through space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rWhat's left of your body floats away through space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rWhat's left of @R$N@r's body floats away through space!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else if (func() int {
			if vict.In_room != room_rnum(-1) && vict.In_room <= top_of_world {
				return (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(vict.In_room)))).Sector_type
			}
			return SECT_INSIDE
		}()) == SECT_FLYING {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r explodes and chunks of $M shower towards the ground far below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou explode leaving only chunks behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r explodes and chunks of $M shower toward the ground far below.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 2:
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rThe bottom half of your body is all that remains as you die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_HALF
			case 3:
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body completely disintegrates in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 4:
				act(libc.CString("@R$N@r falls down toward the ground as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body falls down toward the ground as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r falls down toward the ground as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 5:
				act(libc.CString("@rWhat's left of @R$N@r's body falls toward the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rWhat's left of yor body falls toward the ground as you die!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rWhat's left of @R$N@r's body falls toward the ground as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		} else {
			switch rand_number(1, 5) {
			case 1:
				act(libc.CString("@R$N@r explodes and chunks of $M float freely through the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYou explode leaving only chunks behind.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r explodes and chunks of $M float freely through the water.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 2:
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rThe bottom half of your body is all that remains as you die.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rThe bottom half of @R$N@r is all that remains as $E dies.@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_HALF
			case 3:
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body completely disintegrates in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r is completely disintegrated in the attack!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
				vict.Death_type = DTYPE_VAPOR
			case 4:
				act(libc.CString("@R$N@r falls back as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rYour body falls back as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@R$N@r falls back as a smoldering corpse!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			case 5:
				act(libc.CString("@rWhat's left of @R$N@r's body floats limply as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_CHAR)
				act(libc.CString("@rWhat's left of yor body floats limply as you die!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_VICT)
				act(libc.CString("@rWhat's left of @R$N@r's body floats limply as $E dies!@n"), TRUE, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
			}
		}
	}
}
