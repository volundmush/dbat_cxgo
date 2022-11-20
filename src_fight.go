package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

var combat_list *char_data = nil
var next_combat_list *char_data = nil

func group_bonus(ch *char_data, type_ int) bool {
	var (
		k    *follow_type
		next *follow_type
	)
	if !AFF_FLAGGED(ch, AFF_GROUP) {
		return false
	}
	if ch.Followers != nil {
		for k = ch.Followers; k != nil; k = next {
			next = k.Next
			if !AFF_FLAGGED(k.Follower, AFF_GROUP) {
				continue
			} else {
				if type_ == 0 {
					k.Follower.Lifeforce += int64(float64(GET_LIFEMAX(k.Follower)) * 0.25)
					if k.Follower.Lifeforce > int64(GET_LIFEMAX(k.Follower)) {
						k.Follower.Lifeforce = int64(GET_LIFEMAX(k.Follower))
					}
					send_to_char(k.Follower, libc.CString("@CIncensed by the death of your comrade your life force swells!@n"))
					return true
				} else if type_ == 1 {
					k.Follower.Lifeforce += int64(float64(GET_LIFEMAX(k.Follower)) * 0.4)
					if k.Follower.Lifeforce > int64(GET_LIFEMAX(k.Follower)) {
						k.Follower.Lifeforce = int64(GET_LIFEMAX(k.Follower))
					}
					send_to_char(k.Follower, libc.CString("@CIncensed by the death of your comrade your life force swells!@n"))
					return true
				} else if type_ == 2 {
					if int(ch.Chclass) == CLASS_ROSHI {
						return 2 != 0
					} else if int(ch.Chclass) == CLASS_KRANE {
						return 3 != 0
					} else if int(ch.Chclass) == CLASS_BARDOCK {
						return 4 != 0
					} else if int(ch.Chclass) == CLASS_NAIL {
						return 5 != 0
					} else if int(ch.Chclass) == CLASS_KABITO {
						return 6 != 0
					} else if int(ch.Chclass) == CLASS_ANDSIX {
						return 7 != 0
					} else if int(ch.Chclass) == CLASS_TAPION {
						return 8 != 0
					} else if int(ch.Chclass) == CLASS_FRIEZA {
						return 9 != 0
					} else if int(ch.Chclass) == CLASS_TSUNA {
						return 10 != 0
					} else if int(ch.Chclass) == CLASS_PICCOLO {
						return 11 != 0
					} else if int(ch.Chclass) == CLASS_KURZAK {
						return 12 != 0
					} else if int(ch.Chclass) == CLASS_JINTO {
						return 13 != 0
					} else if int(ch.Chclass) == CLASS_DABURA {
						return 14 != 0
					}
				}
				return false
			}
		}
	} else if ch.Master != nil {
		if !AFF_FLAGGED(ch.Master, AFF_GROUP) {
			return false
		} else {
			if type_ == 0 {
				group_bonus(ch.Master, 0)
			} else if type_ == 2 {
				if int(ch.Master.Chclass) == CLASS_ROSHI {
					return 2 != 0
				} else if int(ch.Master.Chclass) == CLASS_KRANE {
					return 3 != 0
				} else if int(ch.Master.Chclass) == CLASS_BARDOCK {
					return 4 != 0
				} else if int(ch.Master.Chclass) == CLASS_NAIL {
					return 5 != 0
				} else if int(ch.Master.Chclass) == CLASS_KABITO {
					return 6 != 0
				} else if int(ch.Master.Chclass) == CLASS_ANDSIX {
					return 7 != 0
				} else if int(ch.Master.Chclass) == CLASS_TAPION {
					return 8 != 0
				} else if int(ch.Master.Chclass) == CLASS_FRIEZA {
					return 9 != 0
				} else if int(ch.Master.Chclass) == CLASS_TSUNA {
					return 10 != 0
				} else if int(ch.Master.Chclass) == CLASS_PICCOLO {
					return 11 != 0
				} else if int(ch.Master.Chclass) == CLASS_KURZAK {
					return 12 != 0
				} else if int(ch.Master.Chclass) == CLASS_JINTO {
					return 13 != 0
				} else if int(ch.Master.Chclass) == CLASS_DABURA {
					return 14 != 0
				}
			}
			return true
		}
	}
	return false
}
func mutant_limb_regen(ch *char_data) {
	if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 {
		act(libc.CString("The bones in your right arm have mended them selves."), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n starts moving $s right arm gingerly for a moment."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[0] = 100
	} else if (ch.Limb_condition[0]) <= 0 {
		act(libc.CString("Your right arm begins to grow back very quickly. Within moments it is whole again!"), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n's right arm starts to regrow! Within moments the arm is whole again!."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[0] = 100
	}
	if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 {
		act(libc.CString("The bones in your left arm have mended them selves."), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n starts moving $s left arm gingerly for a moment."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[1] = 100
	} else if (ch.Limb_condition[1]) <= 0 {
		act(libc.CString("Your right arm begins to grow back very quickly. Within moments it is whole again!"), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n's right arm starts to regrow! Within moments the arm is whole again!."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[1] = 100
	}
	if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 {
		act(libc.CString("The bones in your right leg have mended them selves."), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n starts moving $s right leg gingerly for a moment."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[2] = 100
	} else if (ch.Limb_condition[2]) <= 0 {
		act(libc.CString("Your right arm begins to grow back very quickly. Within moments it is whole again!"), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n's right arm starts to regrow! Within moments the arm is whole again!."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[2] = 100
	}
	if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 {
		act(libc.CString("The bones in your left leg have mended them selves."), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n starts moving $s left leg gingerly for a moment."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[3] = 100
	} else if (ch.Limb_condition[3]) <= 0 {
		act(libc.CString("Your right arm begins to grow back very quickly. Within moments it is whole again!"), 1, ch, nil, nil, TO_CHAR)
		act(libc.CString("$n's right arm starts to regrow! Within moments the arm is whole again!."), 1, ch, nil, nil, TO_ROOM)
		ch.Limb_condition[3] = 100
	}
}
func pick_n_throw(ch *char_data, buf *byte) bool {
	var (
		cont *obj_data
		buf2 [2048]byte
		buf3 [2048]byte
	)
	if rand_number(1, 20) < 18 {
		return false
	}
	for cont = world[ch.In_room].Contents; cont != nil; cont = cont.Next_content {
		if cont.Weight <= max_carry_weight(ch)+int64(ch.Carry_weight) {
			stdio.Sprintf(&buf2[0], "%s", cont.Name)
			do_get(ch, &buf2[0], 0, 0)
			stdio.Sprintf(&buf3[0], "%s %s", &buf2[0], buf)
			do_throw(ch, &buf3[0], 0, 0)
			return true
		}
	}
	return false
}
func mob_attack(ch *char_data, buf *byte) {
	var (
		power   int = rand_number(1, 5)
		bonus   int = int(float64(GET_LEVEL(ch)) * 0.1)
		special int = 0
		buf2    [2048]byte
	)
	power += bonus
	if rand_number(1, 4) == 4 {
		power += 10
	}
	if power > 20 {
		power = 20
	}
	if int(ch.Chclass) == CLASS_NPC_COMMONER {
		special = 0
	}
	var dragonpass int = 1
	if int(ch.Race) == RACE_LIZARDFOLK {
		if GET_MOB_VNUM(ch) == 81 || GET_MOB_VNUM(ch) == 82 || GET_MOB_VNUM(ch) == 83 || GET_MOB_VNUM(ch) == 84 || GET_MOB_VNUM(ch) == 85 || GET_MOB_VNUM(ch) == 86 || GET_MOB_VNUM(ch) == 87 {
			dragonpass = 1
			special = rand_number(40, 100)
		} else {
			dragonpass = 0
		}
	}
	if axion_dice(-10) > 90 && ch.Hit <= ch.Max_hit/2 && !MOB_FLAGGED(ch, MOB_POWERUP) && GET_MOB_VNUM(ch) != 25 {
		do_powerup(ch, nil, 0, 0)
		return
	}
	if float64(ch.Mana) >= float64(ch.Max_mana)*0.05 && IS_HUMANOID(ch) && (int(ch.Race) != RACE_LIZARDFOLK || dragonpass == 1) {
		if ch.Mobcharge <= 0 && rand_number(1, 10) >= 8 {
			act(libc.CString("@wAn aura flares up around @R$n@w!@n"), 1, ch, nil, nil, TO_ROOM)
			ch.Mobcharge += 1
			if GET_LEVEL(ch) > 80 {
				ch.Mobcharge += 1
			}
		} else if ch.Mobcharge <= 5 {
			act(libc.CString("@wThe aura burns brighter around @R$n@w!@n"), 1, ch, nil, nil, TO_ROOM)
			ch.Mobcharge += 1
			if GET_LEVEL(ch) > 80 {
				ch.Mobcharge += 1
			}
		} else if ch.Mobcharge == 6 {
			act(libc.CString("@wThe aura around @R$n@w flashes!@n"), 1, ch, nil, nil, TO_ROOM)
			ch.Mobcharge += 1
			special = 100
		}
	}
	if IS_HUMANOID(ch) && dragonpass == 1 {
		if AFF_FLAGGED(ch, AFF_PARALYZE) {
			return
		} else if AFF_FLAGGED(ch, AFF_ENSNARED) {
			return
		} else if special < 100 {
			if int(ch.Chclass) == CLASS_SHADOWDANCER && rand_number(1, 3) == 3 {
				stdio.Sprintf(&buf2[0], "ass %s", buf)
				do_throw(ch, &buf2[0], 0, 0)
			} else if int(ch.Race) == RACE_ANDROID && MOB_FLAGGED(ch, MOB_REPAIR) && float64(ch.Hit) <= float64(gear_pl(ch))*0.5 && rand_number(1, 20) >= 16 {
				do_srepair(ch, nil, 0, 0)
			} else if int(ch.Race) == RACE_ANDROID && MOB_FLAGGED(ch, MOB_ABSORB) && rand_number(1, 20) >= 19 {
				do_absorb(ch, &buf2[0], 0, 0)
			} else if (int(ch.Race) == RACE_BIO || int(ch.Race) == RACE_MAJIN) && float64(ch.Hit) <= float64(gear_pl(ch))*0.5 && rand_number(1, 20) >= 17 {
				do_regenerate(ch, libc.CString("25"), 0, 0)
			} else if int(ch.Race) == RACE_NAMEK && float64(ch.Hit) <= float64(gear_pl(ch))*0.5 && rand_number(1, 20) == 20 {
				do_regenerate(ch, libc.CString("25"), 0, 0)
			} else if pick_n_throw(ch, buf) {
			} else if MOB_FLAGGED(ch, MOB_KNOWKAIO) && rand_number(1, 50) >= 46 {
				if rand_number(1, 10) == 10 {
					do_kaioken(ch, libc.CString("20"), 0, 0)
				} else if rand_number(1, 10) >= 8 {
					do_kaioken(ch, libc.CString("10"), 0, 0)
				} else {
					do_kaioken(ch, libc.CString("5"), 0, 0)
				}
			} else {
				switch power {
				case 1:
					fallthrough
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					fallthrough
				case 5:
					if (ch.Equipment[WEAR_WIELD1]) != nil {
						do_attack(ch, buf, 0, 0)
					} else if rand_number(1, 5) == 5 {
						do_kick(ch, buf, 0, 0)
					} else if rand_number(1, 10) == 10 {
						do_elbow(ch, buf, 0, 0)
					} else {
						do_punch(ch, buf, 0, 0)
					}
				case 6:
					fallthrough
				case 7:
					fallthrough
				case 8:
					if (ch.Equipment[WEAR_WIELD1]) != nil {
						do_attack(ch, buf, 0, 0)
					} else if rand_number(1, 5) == 5 {
						do_punch(ch, buf, 0, 0)
					} else if rand_number(1, 10) == 10 {
						do_knee(ch, buf, 0, 0)
					} else {
						do_kick(ch, buf, 0, 0)
					}
				case 9:
					fallthrough
				case 10:
					if rand_number(1, 5) == 5 {
						do_knee(ch, buf, 0, 0)
					} else if rand_number(1, 10) == 10 {
						do_uppercut(ch, buf, 0, 0)
					} else {
						do_elbow(ch, buf, 0, 0)
					}
				case 11:
					fallthrough
				case 12:
					if rand_number(1, 5) == 5 {
						do_elbow(ch, buf, 0, 0)
					} else if rand_number(1, 10) == 10 {
						do_roundhouse(ch, buf, 0, 0)
					} else if rand_number(1, 8) == 8 {
						do_trip(ch, buf, 0, 0)
					} else {
						do_knee(ch, buf, 0, 0)
					}
				case 13:
					fallthrough
				case 14:
					if (int(ch.Chclass) == CLASS_BARDOCK || int(ch.Chclass) == CLASS_KURZAK) && rand_number(1, 2) == 2 {
						do_head(ch, buf, 0, 0)
					} else if (int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO) && rand_number(1, 2) == 2 {
						do_tailwhip(ch, buf, 0, 0)
					} else if rand_number(1, 8) == 8 {
						do_trip(ch, buf, 0, 0)
					} else {
						do_uppercut(ch, buf, 0, 0)
					}
				case 15:
					fallthrough
				case 16:
					if (int(ch.Chclass) == CLASS_BARDOCK || int(ch.Chclass) == CLASS_KURZAK) && rand_number(1, 2) == 2 {
						do_head(ch, buf, 0, 0)
					} else if (int(ch.Race) == RACE_ICER || int(ch.Race) == RACE_BIO) && rand_number(1, 2) == 2 {
						do_tailwhip(ch, buf, 0, 0)
					} else if rand_number(1, 8) >= 7 {
						do_trip(ch, buf, 0, 0)
					} else {
						do_roundhouse(ch, buf, 0, 0)
					}
				case 17:
					fallthrough
				case 18:
					do_slam(ch, buf, 0, 0)
				case 19:
					fallthrough
				case 20:
					do_heeldrop(ch, buf, 0, 0)
				}
			}
		} else {
			mob_specials_used += 1
			switch power {
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 3:
				fallthrough
			case 4:
				if special > 80 {
					do_zanzoken(ch, buf, 0, 0)
				}
				if ch.Mobcharge == 7 {
					ch.Mobcharge = 0
					do_kiball(ch, buf, 0, 0)
				}
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				fallthrough
			case 8:
				if special > 80 {
					do_zanzoken(ch, buf, 0, 0)
				}
				if ch.Mobcharge == 7 {
					ch.Mobcharge = 0
					do_kiblast(ch, buf, 0, 0)
				}
			case 9:
				fallthrough
			case 10:
				fallthrough
			case 11:
				if special > 80 {
					do_zanzoken(ch, buf, 0, 0)
				}
				if int(ch.Race) == RACE_LIZARDFOLK && rand_number(1, 4) == 4 {
					do_breath(ch, buf, 0, 0)
				} else {
					if ch.Mobcharge == 7 {
						ch.Mobcharge = 0
						do_beam(ch, buf, 0, 0)
					}
				}
			case 12:
				fallthrough
			case 13:
				fallthrough
			case 14:
				if special > 80 {
					do_zanzoken(ch, buf, 0, 0)
				}
				if int(ch.Race) == RACE_LIZARDFOLK && rand_number(1, 4) == 4 {
					do_breath(ch, buf, 0, 0)
				} else {
					if ch.Mobcharge == 7 {
						ch.Mobcharge = 0
						do_renzo(ch, buf, 0, 0)
					}
				}
			case 15:
				fallthrough
			case 16:
				if int(ch.Race) == RACE_LIZARDFOLK && rand_number(1, 4) == 4 {
					do_breath(ch, buf, 0, 0)
				} else {
					if ch.Mobcharge == 7 {
						ch.Mobcharge = 0
						do_tsuihidan(ch, buf, 0, 0)
					}
				}
			case 17:
				fallthrough
			case 18:
				if int(ch.Race) == RACE_LIZARDFOLK && rand_number(1, 4) == 4 {
					do_breath(ch, buf, 0, 0)
				} else {
					if ch.Mobcharge == 7 {
						ch.Mobcharge = 0
						do_shogekiha(ch, buf, 0, 0)
					}
				}
			case 19:
				fallthrough
			case 20:
				if int(ch.Race) == RACE_LIZARDFOLK {
					do_breath(ch, buf, 0, 0)
				}
				if ch.Mobcharge == 7 {
					ch.Mobcharge = 0
					switch ch.Chclass {
					case CLASS_ROSHI:
						if special >= 100 {
							do_kakusanha(ch, buf, 0, 0)
						} else if special >= 80 {
							do_kienzan(ch, buf, 0, 0)
						} else if special >= 70 {
							do_kamehameha(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_FRIEZA:
						if special >= 100 {
							do_deathball(ch, buf, 0, 0)
						} else if special >= 80 {
							do_kienzan(ch, buf, 0, 0)
						} else if special >= 70 {
							do_deathbeam(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_KRANE:
						if special >= 100 {
							do_tribeam(ch, buf, 0, 0)
						} else if special >= 80 {
							do_hass(ch, nil, 0, 0)
						} else if special >= 70 {
							do_dodonpa(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_PICCOLO:
						if special >= 100 {
							do_scatter(ch, buf, 0, 0)
						} else if special >= 80 {
							do_sbc(ch, buf, 0, 0)
						} else if special >= 70 {
							do_masenko(ch, buf, 0, 0)
						} else if special >= 100 {
							do_balefire(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_BARDOCK:
						if special >= 100 {
							do_final(ch, buf, 0, 0)
						} else if special >= 80 {
							do_bigbang(ch, buf, 0, 0)
						} else if special >= 70 {
							do_galikgun(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_ANDSIX:
						if special >= 100 {
							do_hellflash(ch, buf, 0, 0)
						} else if special >= 80 {
							do_kousengan(ch, buf, 0, 0)
						} else if special >= 70 {
							do_dualbeam(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_NAIL:
						if special >= 100 {
							do_regenerate(ch, libc.CString("50"), 0, 0)
						} else if special >= 80 {
							do_heal(ch, libc.CString("self"), 0, 0)
						} else if special >= 70 {
							do_masenko(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, nil, 0, 0)
						}
					case CLASS_KURZAK:
						if special >= 100 {
							do_ensnare(ch, buf, 0, 0)
						} else if special >= 80 {
							do_seishou(ch, buf, 0, 0)
						} else if special >= 70 {
							do_renzo(ch, buf, 0, 0)
						} else if special >= 50 {
							do_barrier(ch, libc.CString("40"), 0, 0)
						} else {
							do_barrier(ch, libc.CString("25"), 0, 0)
						}
					case CLASS_JINTO:
						if special >= 100 {
							do_nova(ch, buf, 0, 0)
						} else if special >= 80 {
							do_breaker(ch, buf, 0, 0)
						} else if special >= 70 {
							do_trip(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, libc.CString("40"), 0, 0)
						}
					case CLASS_TSUNA:
						if special >= 100 {
							do_koteiru(ch, buf, 0, 0)
						} else if special >= 80 {
							do_razor(ch, buf, 0, 0)
						} else if special >= 70 {
							do_spike(ch, buf, 0, 0)
						} else {
							do_barrier(ch, libc.CString("20"), 0, 0)
						}
					case CLASS_TAPION:
						if special >= 100 {
							do_pslash(ch, buf, 0, 0)
						} else if special >= 80 {
							do_ddslash(ch, buf, 0, 0)
						} else if special >= 70 {
							do_tslash(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, libc.CString("40"), 0, 0)
						}
					case CLASS_KABITO:
						if special >= 100 {
							do_pbarrage(ch, buf, 0, 0)
						} else if special >= 80 {
							do_psyblast(ch, buf, 0, 0)
						} else if special >= 70 {
							do_heal(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, libc.CString("40"), 0, 0)
						}
					case CLASS_DABURA:
						if special >= 100 {
							do_hellspear(ch, buf, 0, 0)
						} else if special >= 80 {
							do_honoo(ch, buf, 0, 0)
						} else if special >= 70 {
							do_fireshield(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, libc.CString("40"), 0, 0)
						}
					case CLASS_GINYU:
						if special >= 100 {
							do_spiral(ch, buf, 0, 0)
						} else if special >= 80 {
							do_crusher(ch, buf, 0, 0)
						} else if special >= 70 {
							do_eraser(ch, buf, 0, 0)
						} else {
							do_zanzoken(ch, libc.CString("40"), 0, 0)
						}
					}
				}
			}
		}
	} else if !IS_HUMANOID(ch) || dragonpass == 0 {
		if int(ch.Race) == RACE_SNAKE && rand_number(1, 5) == 5 {
			do_strike(ch, buf, 0, 0)
		} else if int(ch.Race) == RACE_LIZARDFOLK && rand_number(1, 12) >= 10 && GET_MOB_VNUM(ch) != 0x45FD {
			do_breath(ch, buf, 0, 0)
		} else {
			if rand_number(1, 10) >= 7 && GET_LEVEL(ch) >= 10 {
				do_ram(ch, buf, 0, 0)
			} else {
				do_bite(ch, buf, 0, 0)
			}
		}
	}
	fight_mtrigger(ch)
}
func cleanup_arena_watch(ch *char_data) {
	var d *descriptor_data
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected != CON_PLAYING {
			continue
		}
		if PRF_FLAGGED(d.Character, PRF_ARENAWATCH) {
			if d.Character.Arenawatch == int(ch.Idnum) {
				REMOVE_BIT_AR(d.Character.Player_specials.Pref[:], PRF_ARENAWATCH)
				d.Character.Arenawatch = -1
			}
		}
	}
}
func shadow_dragons_live() {
	var value int = 0
	if SHADOW_DRAGON1 != -1 || SHADOW_DRAGON2 != -1 || SHADOW_DRAGON3 != -1 || SHADOW_DRAGON4 != -1 || SHADOW_DRAGON5 != -1 || SHADOW_DRAGON6 != -1 || SHADOW_DRAGON7 != -1 {
		value = 1
	}
	if value == 0 {
		SELFISHMETER = 0
		save_mud_time(&time_info)
	}
}
func impact_sound(ch *char_data, mssg *byte) {
	var door int
	for door = 0; door < NUM_OF_DIRS; door++ {
		if CAN_GO(ch, door) {
			send_to_room(world[ch.In_room].Dir_option[door].To_room, libc.CString("%s"), mssg)
		}
	}
}
func remove_limb(vict *char_data, num int) {
	var (
		body_part *obj_data
		part      [1000]byte
		buf       [1000]byte
		buf2      [1000]byte
	)
	body_part = create_obj()
	body_part.Item_number = -1
	body_part.In_room = -1
	switch num {
	case 0:
		stdio.Snprintf(&part[0], int(1000), "@C%s@w's bloody head@n", GET_NAME(vict))
		stdio.Snprintf(&buf[0], int(1000), "%s bloody head", GET_NAME(vict))
	case 1:
		stdio.Snprintf(&part[0], int(1000), "@w%s right arm@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "right arm")
		if PLR_FLAGGED(vict, PLR_CRARM) {
			REMOVE_BIT_AR(vict.Act[:], PLR_CRARM)
		}
	case 2:
		stdio.Snprintf(&part[0], int(1000), "@w%s left arm@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "left arm")
		if PLR_FLAGGED(vict, PLR_CLARM) {
			REMOVE_BIT_AR(vict.Act[:], PLR_CLARM)
		}
	case 3:
		stdio.Snprintf(&part[0], int(1000), "@w%s right leg@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "right leg")
		if PLR_FLAGGED(vict, PLR_CRLEG) {
			REMOVE_BIT_AR(vict.Act[:], PLR_CRLEG)
		}
	case 4:
		stdio.Snprintf(&part[0], int(1000), "@w%s left leg@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "left leg")
		if PLR_FLAGGED(vict, PLR_CLLEG) {
			REMOVE_BIT_AR(vict.Act[:], PLR_CLLEG)
		}
	case 5:
		stdio.Snprintf(&part[0], int(1000), "@wA %s tail@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "tail")
	case 6:
		stdio.Snprintf(&buf[0], int(1000), "tail")
	default:
		stdio.Snprintf(&part[0], int(1000), "@w%s body part@n", pc_race_types[int(vict.Race)])
		stdio.Snprintf(&buf[0], int(1000), "body part")
	}
	body_part.Name = libc.StrDup(&buf[0])
	if num > 0 {
		stdio.Snprintf(&buf2[0], int(1000), "@wA %s is lying here@n", &part[0])
	} else {
		stdio.Snprintf(&buf2[0], int(1000), "%s@w is lying here@n", &part[0])
	}
	body_part.Description = libc.StrDup(&buf2[0])
	body_part.Short_description = libc.StrDup(&part[0])
	body_part.Type_flag = ITEM_OTHER
	SET_BIT_AR(body_part.Wear_flags[:], ITEM_WEAR_TAKE)
	SET_BIT_AR(body_part.Extra_flags[:], ITEM_UNIQUE_SAVE)
	body_part.Value[0] = 0
	body_part.Value[1] = 0
	body_part.Value[2] = 0
	body_part.Value[3] = 0
	body_part.Value[4] = 1
	body_part.Value[5] = 1
	body_part.Weight = int64(rand_number(4, 10))
	body_part.Cost_per_day = 0
	add_unique_id(body_part)
	obj_to_room(body_part, vict.In_room)
}

var attack_hit_text [15]attack_hit_type = [15]attack_hit_type{{Singular: libc.CString("hit"), Plural: libc.CString("hits")}, {Singular: libc.CString("sting"), Plural: libc.CString("stings")}, {Singular: libc.CString("whip"), Plural: libc.CString("whips")}, {Singular: libc.CString("slash"), Plural: libc.CString("slashes")}, {Singular: libc.CString("bite"), Plural: libc.CString("bites")}, {Singular: libc.CString("bludgeon"), Plural: libc.CString("bludgeons")}, {Singular: libc.CString("crush"), Plural: libc.CString("crushes")}, {Singular: libc.CString("pound"), Plural: libc.CString("pounds")}, {Singular: libc.CString("claw"), Plural: libc.CString("claws")}, {Singular: libc.CString("maul"), Plural: libc.CString("mauls")}, {Singular: libc.CString("thrash"), Plural: libc.CString("thrashes")}, {Singular: libc.CString("pierce"), Plural: libc.CString("pierces")}, {Singular: libc.CString("blast"), Plural: libc.CString("blasts")}, {Singular: libc.CString("punch"), Plural: libc.CString("punches")}, {Singular: libc.CString("stab"), Plural: libc.CString("stabs")}}

func fight_stack() {
	var (
		perc int = 0
		ch   *char_data
		tch  *char_data
		wch  *char_data
	)
	for tch = character_list; tch != nil; tch = tch.Next {
		ch = tch
		if int(ch.Position) == POS_FIGHTING {
			ch.Position = POS_STANDING
		}
		if PLR_FLAGGED(ch, PLR_SPIRAL) {
			handle_spiral(ch, nil, GET_SKILL(ch, SKILL_SPIRAL), 0)
		}
		if IS_NPC(ch) && ch.Cooldown > 0 {
			ch.Cooldown -= 1
			if rand_number(1, 2) == 2 && ch.Cooldown > 0 {
				ch.Cooldown -= 1
			}
			if ch.Cooldown > 0 {
				continue
			}
		}
		if IS_NPC(ch) && MOB_FLAGGED(ch, MOB_POWERUP) && axion_dice(0) >= 90 {
			if ch.Hit >= ch.Max_hit {
				act(libc.CString("@g$n@ finishes powering up as $s aura flashes brightly filling the entire area briefly with its light!@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Hit = ch.Max_hit
				REMOVE_BIT_AR(ch.Act[:], MOB_POWERUP)
			} else if ch.Hit >= ch.Max_hit/2 {
				act(libc.CString("@g$n@G continues powering up as torrents of energy crackle within $s aura.@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Hit += ch.Max_hit / 10
			} else if ch.Hit > ch.Max_hit/4 {
				act(libc.CString("@g$n@G powers up as a steady aura around $s body grow brighter.@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Hit += ch.Max_hit / 8
			} else if ch.Hit > 0 {
				act(libc.CString("@g$n@G powers up, as a weak aura flickers around $s body.@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Hit += ch.Max_hit / 5
			}
		}
		if IS_NPC(ch) && AFF_FLAGGED(ch, AFF_FROZEN) {
			continue
		}
		if ch.Grappling == nil && ch.Grappled == nil && ch.Fighting == nil && !PLR_FLAGGED(ch, PLR_CHARGE) && !PLR_FLAGGED(ch, PLR_POWERUP) && ch.Charge <= 0 && !IS_TRANSFORMED(ch) {
			continue
		}
		if ch.Fighting != nil && ch.Fighting.In_room != ch.In_room {
			wch = ch.Fighting
			stop_fighting(wch)
			stop_fighting(ch)
		}
		if ch.Fighting != nil && ch.Drag != nil {
			act(libc.CString("@WYou are forced to stop dragging @C$N@W!@n"), 1, ch, nil, unsafe.Pointer(ch.Drag), TO_CHAR)
			act(libc.CString("@C$n@W is forced to stop dragging @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(ch.Drag), TO_ROOM)
			ch.Drag.Dragged = nil
			ch.Drag = nil
		}
		if float64(ch.Hit) <= (float64(gear_pl(ch))*0.01)*float64(ch.Lifeperc) && ch.Lifeforce > 0 && int(ch.Race) != RACE_ANDROID {
			if rand_number(1, 15) >= 14 {
				if float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.05 || AFF_FLAGGED(ch, AFF_HEALGLOW) || int(ch.Race) == RACE_KANASSAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.03 {
					var (
						refill int64 = 0
						lfcost int64 = int64(float64(GET_LIFEMAX(ch)) * 0.05)
					)
					if (ch.Bonuses[BONUS_DIEHARD]) > 0 && (int(ch.Race) != RACE_MUTANT || (ch.Genome[0]) != 2 && (ch.Genome[1]) != 2) {
						refill = int64(float64(GET_LIFEMAX(ch)) * 0.1)
					} else if (ch.Bonuses[BONUS_DIEHARD]) > 0 && int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 2 || (ch.Genome[1]) == 2) {
						refill = int64(float64(GET_LIFEMAX(ch)) * 0.17)
					} else if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 2 || (ch.Genome[1]) == 2) {
						refill = int64(float64(GET_LIFEMAX(ch)) * 0.12)
					} else if int(ch.Race) == RACE_KANASSAN {
						lfcost = int64(float64(GET_LIFEMAX(ch)) * 0.03)
						refill = int64(float64(GET_LIFEMAX(ch)) * 0.03)
					} else {
						refill = int64(float64(GET_LIFEMAX(ch)) * 0.05)
					}
					ch.Hit += refill
					if !AFF_FLAGGED(ch, AFF_HEALGLOW) {
						ch.Lifeforce -= lfcost
					}
				} else {
					ch.Hit += ch.Lifeforce
					ch.Lifeforce = -1
				}
				if ch.Suppression > 0 && float64(ch.Hit) > ((float64(ch.Max_hit)*0.01)*float64(ch.Suppression)) {
					ch.Hit = int64((float64(ch.Max_hit) * 0.01) * float64(ch.Suppression))
				} else if ch.Hit > gear_pl(ch) {
					ch.Hit = gear_pl(ch)
				}
				send_to_char(ch, libc.CString("@YYour life force has kept you strong@n!\r\n"))
			}
		}
		if !AFF_FLAGGED(ch, AFF_POSITION) {
			if roll_balance(ch) > axion_dice(0) && rand_number(1, 10) >= 7 {
				if ch.Fighting != nil {
					if !AFF_FLAGGED(ch.Fighting, AFF_POSITION) {
						act(libc.CString("@YYou manage to move into an advantageous position!@n"), 1, ch, nil, nil, TO_CHAR)
						act(libc.CString("@y$n@Y manages to move into an advantageous position!@n"), 1, ch, nil, nil, TO_ROOM)
						SET_BIT_AR(ch.Affected_by[:], AFF_POSITION)
					} else {
						var vict *char_data = ch.Fighting
						if roll_balance(ch) > roll_balance(vict) {
							act(libc.CString("@YYou struggle to gain a better position than @y$N@Y and succeed!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_CHAR)
							act(libc.CString("@y$n@Y struggles to gain a better position than you and succeeds!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_VICT)
							act(libc.CString("@y$n@Y struggles to gain a better position than @y$N@Y and succeeds!@n"), 1, ch, nil, unsafe.Pointer(vict), TO_NOTVICT)
							REMOVE_BIT_AR(vict.Affected_by[:], AFF_POSITION)
							SET_BIT_AR(ch.Affected_by[:], AFF_POSITION)
						}
					}
				}
			}
		} else {
			if roll_balance(ch) < axion_dice(-30) || int(ch.Position) < POS_STANDING {
				act(libc.CString("@YYou are moved out of your position!@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@y$n@Y is moved out of $s position!@n"), 1, ch, nil, nil, TO_ROOM)
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_POSITION)
			}
		}
		if ch.Grappling != nil && ch.Grap == 2 && rand_number(1, 11) >= 8 {
			if ch.Grappling.Move >= ch.Grappling.Max_move/8 {
				act(libc.CString("@WYou choke @C$N@W!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_CHAR)
				act(libc.CString("@C$n@W chokes YOU@W!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_VICT)
				act(libc.CString("@C$n@W chokes @c$N@W!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_NOTVICT)
				ch.Grappling.Move -= ch.Grappling.Max_move / 8
			} else {
				act(libc.CString("@WYou choke @C$N@W, and $E passes out!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_CHAR)
				act(libc.CString("@C$n@W chokes YOU@W, and you pass out!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_VICT)
				act(libc.CString("@C$n@W chokes @c$N@W, and $E passes out!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_NOTVICT)
				SET_BIT_AR(ch.Grappling.Affected_by[:], AFF_KNOCKED)
				ch.Grappling.Position = POS_SLEEPING
				ch.Grappling.Grap = -1
				ch.Grappling.Grappled = nil
				ch.Grappling = nil
				ch.Grap = -1
			}
		} else if ch.Grappling != nil && ch.Grap == 4 && rand_number(1, 12) >= 8 {
			act(libc.CString("@WYou crush @C$N@W some more!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_CHAR)
			act(libc.CString("@C$n@W crushes YOU@W some more!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_VICT)
			act(libc.CString("@C$n@W crushes @c$N@W some more!@n"), 1, ch, nil, unsafe.Pointer(ch.Grappling), TO_NOTVICT)
			var damg int64 = int64(float64(ch.Aff_abils.Str) * ((float64(ch.Max_hit) * 0.005) + 10))
			hurt(0, 0, ch, ch.Grappling, nil, damg, 0)
		}
		if ch.Grappled != nil && rand_number(1, 2) == 2 {
			send_to_char(ch, libc.CString("@CTry 'escape' to break free from the hold!@n\r\n"))
		}
		if int(ch.Race) == RACE_HALFBREED && PLR_FLAGGED(ch, PLR_FURY) {
			ch.Rage_meter += 1
			if ch.Rage_meter >= 1000 {
				ch.Hit += int64(float64(gear_pl(ch)) * 0.15)
				ch.Mana += int64(float64(ch.Max_mana) * 0.15)
				ch.Move += int64(float64(ch.Max_move) * 0.15)
				if ch.Hit > gear_pl(ch) {
					ch.Hit = gear_pl(ch)
				}
				if ch.Mana > ch.Max_mana {
					ch.Mana = ch.Max_mana
				}
				if ch.Move > ch.Max_move {
					ch.Move = ch.Max_move
				}
				send_to_char(ch, libc.CString("Your fury has called forth more of your hidden power and you feel better!\r\n"))
			}
		}
		if !IS_NPC(ch) && IS_TRANSFORMED(ch) {
			if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move < ch.Max_move/60 {
				act(libc.CString("@mExhausted of stamina, your body forcibly reverts from its form.@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@C$n @wbreathing heavily, reverts from $s form, returning to normal.@n"), 1, ch, nil, nil, TO_ROOM)
				if ch.Kaioken < 1 {
					do_transform(ch, libc.CString("revert"), 0, 0)
				} else if ch.Kaioken >= 1 {
					do_kaioken(ch, libc.CString("0"), 0, 0)
					do_transform(ch, libc.CString("revert"), 0, 0)
				}
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/900 && PLR_FLAGGED(ch, PLR_TRANS1) && int(ch.Race) != RACE_KONATSU && int(ch.Race) != RACE_KAI && int(ch.Race) != RACE_NAMEK {
				if int(ch.Race) == RACE_SAIYAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.7 {
					ch.Move -= ch.Max_move / 1000
				} else {
					ch.Move -= ch.Max_move / 900
				}
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/800 && PLR_FLAGGED(ch, PLR_TRANS1) {
				if int(ch.Race) == RACE_SAIYAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.7 {
					ch.Move -= ch.Max_move / 900
				} else {
					ch.Move -= ch.Max_move / 800
				}
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/600 && PLR_FLAGGED(ch, PLR_TRANS2) && int(ch.Race) != RACE_KONATSU && int(ch.Race) != RACE_KAI && int(ch.Race) != RACE_NAMEK {
				if int(ch.Race) == RACE_SAIYAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.7 {
					ch.Move -= ch.Max_move / 700
				} else {
					ch.Move -= ch.Max_move / 600
				}
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/500 && PLR_FLAGGED(ch, PLR_TRANS2) {
				ch.Move -= ch.Max_move / 500
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/400 && PLR_FLAGGED(ch, PLR_TRANS3) && int(ch.Race) != RACE_SAIYAN {
				ch.Move -= ch.Max_move / 400
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/250 && PLR_FLAGGED(ch, PLR_TRANS3) {
				if int(ch.Race) == RACE_SAIYAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.7 {
					ch.Move -= ch.Max_move / 300
				} else {
					ch.Move -= ch.Max_move / 250
				}
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/200 && PLR_FLAGGED(ch, PLR_TRANS4) && int(ch.Race) != RACE_SAIYAN {
				ch.Move -= ch.Max_move / 200
			} else if IS_NONPTRANS(ch) && int(ch.Race) != RACE_ICER && ch.Move >= ch.Max_move/170 && PLR_FLAGGED(ch, PLR_TRANS4) {
				if int(ch.Race) == RACE_SAIYAN && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.7 {
					ch.Move -= ch.Max_move / 240
				} else {
					ch.Move -= ch.Max_move / 170
				}
			}
		}
		if !IS_NPC(ch) && ch.Player_specials.Wimp_level != 0 && ch.Hit < int64(ch.Player_specials.Wimp_level) && ch.Hit > 0 && ch.Fighting != nil {
			send_to_char(ch, libc.CString("You wimp out, and attempt to flee!\r\n"))
			do_flee(ch, nil, 0, 0)
		}
		if IS_NPC(ch) && ch.Hit < ch.Max_hit/10 && ch.Hit > 0 && ch.Fighting != nil && !MOB_FLAGGED(ch, MOB_SENTINEL) {
			if rand_number(1, 30) >= 25 && int(ch.Position) > POS_SITTING {
				do_flee(ch, nil, 0, 0)
			}
		}
		if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 6 || (ch.Genome[1]) == 6) && rand_number(1, 200) >= 175 {
			mutant_limb_regen(ch)
		}
		if !IS_NPC(ch) && PLR_FLAGGED(ch, PLR_DISGUISED) && GET_SKILL(ch, SKILL_DISGUISE) < rand_number(1, 125) {
			send_to_char(ch, libc.CString("Your disguise comes off because of your swift movements!\r\n"))
			REMOVE_BIT_AR(ch.Act[:], PLR_DISGUISED)
			act(libc.CString("@W$n's@W disguise comes off because of $s swift movements!@n"), 0, ch, nil, nil, TO_ROOM)
		}
		if IS_NPC(ch) && AFF_FLAGGED(ch, AFF_BLIND) && rand_number(1, 200) >= 190 {
			act(libc.CString("@W$n@W is no longer blind.@n"), 0, ch, nil, nil, TO_ROOM)
			REMOVE_BIT_AR(ch.Affected_by[:], AFF_BLIND)
		}
		if AFF_FLAGGED(ch, AFF_KNOCKED) && rand_number(1, 200) >= 195 {
			act(libc.CString("@W$n@W is no longer senseless, and wakes up.@n"), 0, ch, nil, nil, TO_ROOM)
			send_to_char(ch, libc.CString("You are no longer knocked out, and wake up!@n\r\n"))
			REMOVE_BIT_AR(ch.Affected_by[:], AFF_KNOCKED)
			ch.Position = POS_SITTING
			if IS_NPC(ch) && rand_number(1, 20) >= 12 {
				act(libc.CString("@W$n@W stands up.@n"), 0, ch, nil, nil, TO_ROOM)
				ch.Position = POS_STANDING
			}
		}
		if !IS_NPC(ch) && ch.Desc == nil && int(ch.Position) > POS_STUNNED && !AFF_FLAGGED(ch, AFF_FROZEN) {
			if ch.Fighting != nil {
				do_flee(ch, nil, 0, 0)
			}
		}
		if IS_NPC(ch) && ch.Grappled != nil && !MOB_FLAGGED(ch, MOB_DUMMY) && rand_number(1, 5) >= 4 {
			do_escape(ch, nil, 0, 0)
			continue
		}
		if ch.Fighting != nil && IS_NPC(ch) && !MOB_FLAGGED(ch, MOB_DUMMY) {
			if AFF_FLAGGED(ch.Fighting, AFF_FLYING) && !AFF_FLAGGED(ch, AFF_FLYING) && IS_HUMANOID(ch) && GET_LEVEL(ch) > 10 {
				do_fly(ch, nil, 0, 0)
				continue
			}
			if !AFF_FLAGGED(ch.Fighting, AFF_FLYING) && AFF_FLAGGED(ch, AFF_FLYING) {
				do_fly(ch, nil, 0, 0)
				continue
			}
			if AFF_FLAGGED(ch.Fighting, AFF_FLYING) && AFF_FLAGGED(ch, AFF_FLYING) && ch.Altitude < ch.Fighting.Altitude {
				do_fly(ch, libc.CString("high"), 0, 0)
				continue
			}
			if AFF_FLAGGED(ch.Fighting, AFF_FLYING) && !IS_HUMANOID(ch) && !AFF_FLAGGED(ch, AFF_FLYING) && int(ch.Position) > POS_RESTING {
				if rand_number(1, 30) >= 22 && !block_calc(ch) {
					act(libc.CString("$n@G flees in terror and you lose sight of $m!"), 1, ch, nil, nil, TO_ROOM)
					for ch.Carrying != nil {
						extract_obj(ch.Carrying)
					}
					extract_char(ch)
					continue
				}
			}
			if AFF_FLAGGED(ch.Fighting, AFF_FLYING) && IS_HUMANOID(ch) && GET_LEVEL(ch) <= 10 {
				if rand_number(1, 30) >= 22 && !block_calc(ch) {
					act(libc.CString("$n@G turns and runs away. You lose sight of $m!"), 1, ch, nil, nil, TO_ROOM)
					for ch.Carrying != nil {
						extract_obj(ch.Carrying)
					}
					extract_char(ch)
					continue
				}
			}
			if int(ch.Position) == POS_SITTING && sec_roll_check(ch) == 1 {
				do_stand(ch, nil, 0, 0)
				continue
			}
			if int(ch.Position) == POS_RESTING && sec_roll_check(ch) == 1 {
				do_stand(ch, nil, 0, 0)
				continue
			}
			if AFF_FLAGGED(ch, AFF_PARA) && IS_NPC(ch) && int(ch.Aff_abils.Intel)+10 < rand_number(1, 60) {
				act(libc.CString("@yYou fail to overcome your paralysis!@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@Y$n @ystruggles with $s paralysis!@n"), 1, ch, nil, nil, TO_ROOM)
				continue
			}
			if int(ch.Position) == POS_SLEEPING && !AFF_FLAGGED(ch, AFF_KNOCKED) && sec_roll_check(ch) == 1 {
				do_wake(ch, nil, 0, 0)
				do_stand(ch, nil, 0, 0)
				continue
			}
			var vict *char_data
			var buf [100]byte
			vict = ch.Fighting
			stdio.Sprintf(&buf[0], "%s", GET_NAME(vict))
			if ch.In_room == vict.In_room && !MOB_FLAGGED(ch, MOB_DUMMY) && !AFF_FLAGGED(ch, AFF_KNOCKED) && int(ch.Position) != POS_SITTING && int(ch.Position) != POS_RESTING && int(ch.Position) != POS_SLEEPING {
				if IS_NPC(ch) && rand_number(1, 30) <= 12 {
					continue
				}
				mob_attack(ch, &buf[0])
			} else {
				continue
			}
		}
		if int(ch.Position) <= POS_RESTING && PLR_FLAGGED(ch, PLR_POWERUP) {
			REMOVE_BIT_AR(ch.Act[:], PLR_POWERUP)
		}
		if PLR_FLAGGED(ch, PLR_POWERUP) && rand_number(1, 3) == 3 {
			var buf3 [64936]byte
			if ch.Hit >= gear_pl(ch) && ch.Mana >= ch.Max_mana/20 && ch.Preference != PREFERENCE_KI {
				if float64(ch.Mana) >= float64(ch.Max_mana)*0.5 {
					var raise int64 = int64(float64(ch.Max_move) * 0.02)
					if ch.Move+raise < ch.Max_move {
						ch.Move += raise
					} else {
						ch.Move = ch.Max_move
					}
				}
				ch.Hit = gear_pl(ch)
				ch.Mana -= ch.Max_mana / 20
				dispel_ash(ch)
				act(libc.CString("@RYou have reached your maximum!@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n stops powering up in a flash of light!@n"), 1, ch, nil, nil, TO_ROOM)
				send_to_sense(0, libc.CString("You sense someone stop powering up"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Rising Powerlevel Final@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				REMOVE_BIT_AR(ch.Act[:], PLR_POWERUP)
			} else if ch.Hit >= gear_pl(ch) && float64(ch.Mana) >= (float64(ch.Max_mana)*0.0375)+1 && ch.Preference == PREFERENCE_KI {
				if float64(ch.Mana) >= (float64(ch.Max_mana)*0.0375)+1 {
					var raise int64 = int64(float64(ch.Max_move) * 0.02)
					if ch.Move+raise < ch.Max_move {
						ch.Move += raise
					} else {
						ch.Move = ch.Max_move
					}
				}
				ch.Hit = gear_pl(ch)
				ch.Mana -= int64((float64(ch.Max_mana) * 0.0375) + 1)
				dispel_ash(ch)
				act(libc.CString("@RYou have reached your maximum!@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n stops powering up in a flash of light!@n"), 1, ch, nil, nil, TO_ROOM)
				send_to_sense(0, libc.CString("You sense someone stop powering up"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Rising Powerlevel Final@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				REMOVE_BIT_AR(ch.Act[:], PLR_POWERUP)
			}
			if ch.Mana < ch.Max_mana/20 && ch.Preference != PREFERENCE_KI {
				ch.Mana = 0
				act(libc.CString("@RYou have run out of ki.@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n stops powering up in a flash of light!@n"), 1, ch, nil, nil, TO_ROOM)
				send_to_sense(0, libc.CString("You sense someone stop powering up"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Rising Powerlevel Final@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				REMOVE_BIT_AR(ch.Act[:], PLR_POWERUP)
			} else if float64(ch.Mana) < (float64(ch.Max_mana)*0.0375)+1 && ch.Preference == PREFERENCE_KI {
				ch.Mana = 0
				act(libc.CString("@RYou have run out of ki.@n"), 1, ch, nil, nil, TO_CHAR)
				act(libc.CString("@R$n stops powering up in a flash of light!@n"), 1, ch, nil, nil, TO_ROOM)
				send_to_sense(0, libc.CString("You sense someone stop powering up"), ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Rising Powerlevel Final@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				REMOVE_BIT_AR(ch.Act[:], PLR_POWERUP)
			}
			if ch.Hit < gear_pl(ch) && (ch.Preference != PREFERENCE_KI && ch.Mana >= ch.Max_mana/20 || ch.Preference == PREFERENCE_KI && float64(ch.Mana) >= (float64(ch.Max_mana)*0.0375)+1) {
				ch.Hit += gear_pl(ch) / 10
				if ch.Preference != PREFERENCE_KI {
					ch.Mana -= ch.Max_mana / 20
				} else {
					ch.Mana -= int64(float64(ch.Max_mana) * 0.0375)
				}
				if float64(ch.Mana) >= float64(ch.Max_mana)*0.5 {
					var raise int64 = int64(float64(ch.Max_move) * 0.02)
					if ch.Move+raise < ch.Max_move {
						ch.Move += raise
					} else {
						ch.Move = ch.Max_move
					}
				}
				if ch.Max_hit < 50000 {
					act(libc.CString("@RYou continue to powerup, as wind billows out from around you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as wind billows out from around $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else if ch.Max_hit < 500000 {
					act(libc.CString("@RYou continue to powerup, as the ground splits beneath you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as the ground splits beneath $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else if ch.Max_hit < 5000000 {
					act(libc.CString("@RYou continue to powerup, as the ground shudders and splits beneath you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as the ground shudders and splits beneath $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else if ch.Max_hit < 50000000 {
					act(libc.CString("@RYou continue to powerup, as a huge depression forms beneath you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as a huge depression forms beneath $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else if ch.Max_hit < 100000000 {
					act(libc.CString("@RYou continue to powerup, as the entire area quakes around you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as the entire area quakes around $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else if ch.Max_hit < 300000000 {
					act(libc.CString("@RYou continue to powerup, as huge chunks of ground are ripped apart beneath you!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as huge chunks of ground are ripped apart beanth $m!@n"), 1, ch, nil, nil, TO_ROOM)
				} else {
					act(libc.CString("@RYou continue to powerup, as the very air around you crackles and burns!@n"), 1, ch, nil, nil, TO_CHAR)
					act(libc.CString("@R$n continues to powerup, as the very air around $m crackles and burns!@n"), 1, ch, nil, nil, TO_ROOM)
				}
				send_to_sense(0, libc.CString("You sense someone powering up"), ch)
				send_to_worlds(ch)
				stdio.Sprintf(&buf3[0], "@D[@GBlip@D]@r Rising Powerlevel Detected@D: [@Y%s@D]", add_commas(ch.Hit))
				send_to_scouter(&buf3[0], ch, 1, 0)
				dispel_ash(ch)
			}
		}
		if (int(ch.Position) == POS_SLEEPING || int(ch.Position) == POS_RESTING) && (PLR_FLAGGED(ch, PLR_CHARGE) || ch.Charge >= 1) {
			send_to_char(ch, libc.CString("You stop charging and release all your pent up energy!\r\n"))
			switch rand_number(1, 3) {
			case 1:
				act(libc.CString("$n@w's aura disappears.@n"), 1, ch, nil, nil, TO_ROOM)
			case 2:
				act(libc.CString("$n@w's aura fades.@n"), 1, ch, nil, nil, TO_ROOM)
			case 3:
				act(libc.CString("$n@w's aura flickers brightly before disappearing.@n"), 1, ch, nil, nil, TO_ROOM)
			default:
				act(libc.CString("$n@w's aura disappears.@n"), 1, ch, nil, nil, TO_ROOM)
			}
			REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
			ch.Mana += ch.Charge
			if ch.Mana > ch.Max_mana {
				ch.Mana = ch.Max_mana
			}
			ch.Charge = 0
			ch.Chargeto = 0
		}
		if PLR_FLAGGED(ch, PLR_CHARGE) && (ch.Bonuses[BONUS_UNFOCUSED]) > 0 && rand_number(1, 80) >= 70 {
			send_to_char(ch, libc.CString("You lose concentration due to your unfocused mind and release your charged energy!\r\n"))
			switch rand_number(1, 3) {
			case 1:
				act(libc.CString("$n@w's aura disappears.@n"), 1, ch, nil, nil, TO_ROOM)
			case 2:
				act(libc.CString("$n@w's aura fades.@n"), 1, ch, nil, nil, TO_ROOM)
			case 3:
				act(libc.CString("$n@w's aura flickers brightly before disappearing.@n"), 1, ch, nil, nil, TO_ROOM)
			default:
				act(libc.CString("$n@w's aura disappears.@n"), 1, ch, nil, nil, TO_ROOM)
			}
			REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
			ch.Mana += ch.Charge
			if ch.Mana > ch.Max_mana {
				ch.Mana = ch.Max_mana
			}
			ch.Charge = 0
			ch.Chargeto = 0
		}
		if !PLR_FLAGGED(ch, PLR_CHARGE) && rand_number(1, 40) >= 38 && ch.Fighting == nil && (ch.Preference != PREFERENCE_KI || float64(ch.Charge) > float64(ch.Max_mana)*0.1) {
			if ch.Charge >= ch.Max_mana/100 {
				var loss int64 = 0
				send_to_char(ch, libc.CString("You lose some of your energy slowly.\r\n"))
				switch rand_number(1, 3) {
				case 1:
					act(libc.CString("$n@w's aura flickers weakly.@n"), 1, ch, nil, nil, TO_ROOM)
				case 2:
					act(libc.CString("$n@w's aura sheds energy.@n"), 1, ch, nil, nil, TO_ROOM)
				case 3:
					act(libc.CString("$n@w's aura flickers brightly before growing dimmer.@n"), 1, ch, nil, nil, TO_ROOM)
				default:
					act(libc.CString("$n@w's aura shrinks some.@n"), 1, ch, nil, nil, TO_ROOM)
				}
				loss = ch.Charge / 20
				ch.Charge -= loss
			} else if ch.Charge < ch.Max_mana/100 && ch.Charge != 0 {
				send_to_char(ch, libc.CString("Your charged energy is completely gone as your aura fades.\r\n"))
				act(libc.CString("$n@w's aura fades away dimmly.@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Charge = 0
			}
		}
		if PLR_FLAGGED(ch, PLR_CHARGE) {
			if GET_SKILL(ch, SKILL_CONCENTRATION) > 74 {
				perc = 10
			} else if GET_SKILL(ch, SKILL_CONCENTRATION) > 49 {
				perc = 5
			} else if GET_SKILL(ch, SKILL_CONCENTRATION) > 24 {
				perc = 2
			} else {
				perc = 1
			}
			if int(ch.Race) == RACE_TRUFFLE && perc == 10 {
				perc += 10
			}
			if int(ch.Race) == RACE_TRUFFLE && perc == 5 {
				perc += 5
			}
			if int(ch.Race) == RACE_TRUFFLE && perc == 2 {
				perc += 3
			}
			if int(ch.Race) == RACE_TRUFFLE && perc == 1 {
				perc += 1
			}
			if perc > 1 && ch.Preference == PREFERENCE_H2H {
				perc = int(float64(perc) * 0.5)
			}
		}
		if PLR_FLAGGED(ch, PLR_CHARGE) {
			if GET_SKILL(ch, SKILL_CONCENTRATION) > 74 {
				perc = 10
			} else if GET_SKILL(ch, SKILL_CONCENTRATION) > 49 {
				perc = 5
			} else if GET_SKILL(ch, SKILL_CONCENTRATION) > 24 {
				perc = 2
			} else {
				perc = 1
			}
			if int(ch.Race) == RACE_MUTANT && perc == 10 {
				perc -= 1
			}
			if int(ch.Race) == RACE_MUTANT && perc == 5 {
				perc -= 1
			}
			if int(ch.Race) == RACE_MUTANT && perc == 2 {
				perc -= 1
			}
			if perc > 1 && ch.Preference == PREFERENCE_H2H {
				perc = int(float64(perc) * 0.5)
			}
			if ch.Mana <= 0 {
				send_to_char(ch, libc.CString("You can not charge anymore, you have charged all your energy!\r\n"))
				act(libc.CString("$n@w's aura grows calm.@n"), 1, ch, nil, nil, TO_ROOM)
				REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
			} else if ((float64(ch.Max_mana) * 0.01) * float64(perc)) >= float64(ch.Mana) {
				send_to_char(ch, libc.CString("You have charged the last that you can.\r\n"))
				act(libc.CString("$n@w's aura @Yflashes@w spectacularly, rushing upwards in torrents!@n"), 1, ch, nil, nil, TO_ROOM)
				ch.Charge += ch.Mana
				ch.Mana = 0
				ch.Chargeto = 0
				REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
			} else {
				if ch.Charge >= ch.Chargeto {
					send_to_char(ch, libc.CString("You have already reached the maximum that you wished to charge.\r\n"))
					act(libc.CString("$n@w's aura burns steadily.@n"), 1, ch, nil, nil, TO_ROOM)
					ch.Chargeto = 0
					REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
				} else if float64(ch.Charge)+(((float64(ch.Max_mana)*0.01)*float64(perc))+1) >= float64(ch.Chargeto) {
					ch.Mana -= ch.Chargeto - ch.Charge
					ch.Charge = ch.Chargeto
					send_to_char(ch, libc.CString("You stop charging as you reach the maximum that you wished to charge.\r\n"))
					act(libc.CString("$n@w's aura flares up brightly and then burns steadily.@n"), 1, ch, nil, nil, TO_ROOM)
					ch.Chargeto = 0
					REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
				} else {
					ch.Mana -= int64(((float64(ch.Max_mana) * 0.01) * float64(perc)) + 1)
					ch.Charge += int64(((float64(ch.Max_mana) * 0.01) * float64(perc)) + 1)
					switch rand_number(1, 3) {
					case 1:
						act(libc.CString("$n@w's aura ripples magnificantly while growing brighter!@n"), 1, ch, nil, nil, TO_ROOM)
						send_to_char(ch, libc.CString("Your aura grows bright as you charge more ki.\r\n"))
					case 2:
						act(libc.CString("$n@w's aura ripples with power as it grows larger!@n"), 1, ch, nil, nil, TO_ROOM)
						send_to_char(ch, libc.CString("Your aura ripples with power as you charge more ki.\r\n"))
					case 3:
						act(libc.CString("$n@w's aura throws sparks off violently!.@n"), 1, ch, nil, nil, TO_ROOM)
						send_to_char(ch, libc.CString("Your aura throws sparks off violently as you charge more ki.\r\n"))
					default:
					}
					if ch.Charge >= ch.Chargeto {
						ch.Charge = ch.Chargeto
						ch.Charge += int64(GET_LEVEL(ch))
						send_to_char(ch, libc.CString("You have finished charging!\r\n"))
						act(libc.CString("$n@w's aura burns brightly and then evens out.@n"), 1, ch, nil, nil, TO_ROOM)
						REMOVE_BIT_AR(ch.Act[:], PLR_CHARGE)
						ch.Chargeto = 0
					}
				}
				if GET_SKILL(ch, SKILL_CONCENTRATION) != 0 {
					improve_skill(ch, SKILL_CONCENTRATION, 1)
				}
			}
		}
	}
}
func appear(ch *char_data) {
	if affected_by_spell(ch, SPELL_INVISIBLE) {
		affect_from_char(ch, SPELL_INVISIBLE)
	}
	if AFF_FLAGGED(ch, AFF_INVISIBLE) {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_INVISIBLE)
	}
	if AFF_FLAGGED(ch, AFF_HIDE) {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_HIDE)
	}
	act(libc.CString("$n slowly fades into existence."), 0, ch, nil, nil, TO_ROOM)
}
func update_pos(victim *char_data) {
	if AFF_FLAGGED(victim, AFF_KNOCKED) {
		return
	}
	if victim.Hit > 0 && int(victim.Position) > POS_STUNNED {
		return
	} else if int(victim.Position) == POS_SITTING && victim.Fighting != nil {
		return
	} else if int(victim.Position) == POS_SITTING && victim.Fighting != nil {
		return
	} else if victim.Hit > 0 {
		victim.Position = POS_STANDING
	} else if victim.Hit <= -11 {
		victim.Position = POS_DEAD
	} else if victim.Hit <= -6 {
		victim.Position = POS_MORTALLYW
	} else if victim.Hit <= -3 {
		victim.Position = POS_INCAP
	} else {
		victim.Position = POS_STUNNED
	}
}
func check_killer(ch *char_data, vict *char_data) {
	if PLR_FLAGGED(vict, PLR_KILLER) || PLR_FLAGGED(vict, PLR_THIEF) {
		return
	}
	if PLR_FLAGGED(ch, PLR_KILLER) || IS_NPC(ch) || IS_NPC(vict) || ch == vict {
		return
	}
}
func set_fighting(ch *char_data, vict *char_data) {
	if ch == vict {
		return
	}
	if ch.Fighting != nil {
		panic("THIS HAPPENED!")
		return
	}
	ch.Next_fighting = combat_list
	combat_list = ch
	ch.Fighting = vict
	if int(ch.Position) == POS_SITTING {
		ch.Position = POS_SITTING
	} else if int(ch.Position) == POS_SLEEPING {
		ch.Position = POS_SLEEPING
	}
	if config_info.Play.Pk_allowed == 0 {
		check_killer(ch, vict)
	}
}
func stop_fighting(ch *char_data) {
	var temp *char_data
	if ch == next_combat_list {
		next_combat_list = ch.Next_fighting
	}
	if IS_NPC(ch) {
		ch.Combo = -1
		ch.Combhits = 0
	}
	if ch == combat_list {
		combat_list = ch.Next_fighting
	} else {
		temp = combat_list
		for temp != nil && temp.Next_fighting != ch {
			temp = temp.Next_fighting
		}
		if temp != nil {
			temp.Next_fighting = ch.Next_fighting
		}
	}
	ch.Next_fighting = nil
	ch.Fighting = nil
	if AFF_FLAGGED(ch, AFF_POSITION) {
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_POSITION)
	}
	update_pos(ch)
}
func make_pcorpse(ch *char_data) {
	var (
		corpse *obj_data
		money  *obj_data
		x      int
		y      int
	)
	corpse = create_obj()
	corpse.Item_number = -1
	corpse.In_room = -1
	handle_corpse_condition(corpse, ch)
	if AFF_FLAGGED(ch, AFF_ASHED) {
		act(libc.CString("@WSome ashes fall off the corpse.@n"), 1, ch, nil, nil, TO_ROOM)
		var ashes *obj_data
		if rand_number(1, 3) == 2 {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		} else if rand_number(1, 2) == 2 {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		} else {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		}
	}
	corpse.Type_flag = ITEM_CONTAINER
	corpse.Size = get_size(ch)
	for x = func() int {
		y = 0
		return y
	}(); x < EF_ARRAY_MAX || y < TW_ARRAY_MAX; func() int {
		x++
		return func() int {
			p := &y
			x := *p
			*p++
			return x
		}()
	}() {
		if x < EF_ARRAY_MAX {
			corpse.Extra_flags[x] = 0
		}
		if y < TW_ARRAY_MAX {
			corpse.Wear_flags[y] = 0
		}
	}
	SET_BIT_AR(corpse.Wear_flags[:], ITEM_WEAR_TAKE)
	SET_BIT_AR(corpse.Extra_flags[:], ITEM_NODONATE)
	corpse.Value[VAL_CONTAINER_CAPACITY] = 0
	corpse.Value[VAL_CONTAINER_CORPSE] = 1
	corpse.Value[VAL_CONTAINER_OWNER] = ch.Pfilepos
	corpse.Weight = int64(GET_PC_WEIGHT(ch) + ch.Carry_weight)
	corpse.Cost_per_day = 100000
	corpse.Timer = config_info.Play.Max_pc_corpse_time
	SET_BIT_AR(corpse.Extra_flags[:], ITEM_UNIQUE_SAVE)
	var obj *obj_data
	var next_obj *obj_data
	_ = next_obj
	for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
		if obj != nil && GET_OBJ_VNUM(obj) < 19900 && GET_OBJ_VNUM(obj) != 0x464E {
			if GET_OBJ_VNUM(obj) >= 18800 && GET_OBJ_VNUM(obj) <= 0x4A37 || GET_OBJ_VNUM(obj) >= 19100 && GET_OBJ_VNUM(obj) <= 0x4AFF {
				continue
			} else {
				obj_from_char(obj)
				obj_to_obj(obj, corpse)
				continue
			}
		} else {
			continue
		}
	}
	if ch.Gold > 0 {
		if IS_NPC(ch) || ch.Desc != nil {
			money = create_money(ch.Gold)
			obj_to_obj(money, corpse)
		}
		ch.Gold = 0
	}
	obj_to_room(corpse, ch.In_room)
}
func handle_corpse_condition(corpse *obj_data, ch *char_data) {
	var (
		buf2    [148]byte
		descBuf [512]byte
	)
	corpse.Value[VAL_CORPSE_HEAD] = 1
	corpse.Value[VAL_CORPSE_RARM] = 1
	corpse.Value[VAL_CORPSE_LARM] = 1
	corpse.Value[VAL_CORPSE_RLEG] = 1
	corpse.Value[VAL_CORPSE_LLEG] = 1
	switch ch.Death_type {
	case DTYPE_HEAD:
		buf2[0] = '\x00'
		stdio.Snprintf(&buf2[0], int(148), "headless corpse %s", GET_NAME(ch))
		corpse.Name = libc.StrDup(&buf2[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The headless corpse of %s is lying here", GET_NAME(ch))
		corpse.Description = libc.StrDup(&descBuf[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The headless remains of %s's corpse", GET_NAME(ch))
		corpse.Short_description = libc.StrDup(&descBuf[0])
		corpse.Value[VAL_CORPSE_HEAD] = 0
	case DTYPE_HALF:
		buf2[0] = '\x00'
		stdio.Snprintf(&buf2[0], int(148), "half corpse %s", GET_NAME(ch))
		corpse.Name = libc.StrDup(&buf2[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "Half of %s's corpse is lying here", GET_NAME(ch))
		corpse.Description = libc.StrDup(&descBuf[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "Half of %s's corpse", GET_NAME(ch))
		corpse.Short_description = libc.StrDup(&descBuf[0])
	case DTYPE_VAPOR:
		buf2[0] = '\x00'
		stdio.Snprintf(&buf2[0], int(148), "burnt chunks corpse %s", GET_NAME(ch))
		corpse.Name = libc.StrDup(&buf2[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The burnt chunks of %s's corpse are scattered here", GET_NAME(ch))
		corpse.Description = libc.StrDup(&descBuf[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The burnt chunks of %s's corpse", GET_NAME(ch))
		corpse.Short_description = libc.StrDup(&descBuf[0])
	case DTYPE_PULP:
		buf2[0] = '\x00'
		stdio.Snprintf(&buf2[0], int(148), "beaten bloody corpse %s", GET_NAME(ch))
		corpse.Name = libc.StrDup(&buf2[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The bloody and beaten corpse of %s is lying here", GET_NAME(ch))
		corpse.Description = libc.StrDup(&descBuf[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The bloody and beaten remains of %s's corpse", GET_NAME(ch))
		corpse.Short_description = libc.StrDup(&descBuf[0])
	default:
		stdio.Snprintf(&buf2[0], int(148), "corpse %s", GET_NAME(ch))
		corpse.Name = libc.StrDup(&buf2[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "The corpse of %s is lying here", GET_NAME(ch))
		corpse.Description = libc.StrDup(&descBuf[0])
		descBuf[0] = '\x00'
		stdio.Snprintf(&descBuf[0], int(512), "the remains of %s's corpse", GET_NAME(ch))
		corpse.Short_description = libc.StrDup(&descBuf[0])
	}
	if !IS_NPC(ch) {
		if (ch.Limb_condition[0]) <= 0 {
			corpse.Value[VAL_CORPSE_RARM] = 0
		} else if (ch.Limb_condition[0]) > 0 && (ch.Limb_condition[0]) < 50 {
			corpse.Value[VAL_CORPSE_RARM] = 2
		}
		if (ch.Limb_condition[1]) <= 0 {
			corpse.Value[VAL_CORPSE_LARM] = 0
		} else if (ch.Limb_condition[1]) > 0 && (ch.Limb_condition[1]) < 50 {
			corpse.Value[VAL_CORPSE_LARM] = 2
		}
		if (ch.Limb_condition[2]) <= 0 {
			corpse.Value[VAL_CORPSE_RLEG] = 0
		} else if (ch.Limb_condition[2]) > 0 && (ch.Limb_condition[2]) < 50 {
			corpse.Value[VAL_CORPSE_RLEG] = 2
		}
		if (ch.Limb_condition[3]) <= 0 {
			corpse.Value[VAL_CORPSE_LLEG] = 0
		} else if (ch.Limb_condition[3]) > 0 && (ch.Limb_condition[3]) < 50 {
			corpse.Value[VAL_CORPSE_LLEG] = 2
		}
		return
	} else {
		return
	}
}
func make_corpse(ch *char_data, tch *char_data) {
	var (
		corpse   *obj_data
		o        *obj_data
		money    *obj_data
		obj      *obj_data
		next_obj *obj_data
	)
	_ = next_obj
	var meat *obj_data
	var i int
	var x int
	var y int
	corpse = create_obj()
	corpse.Item_number = -1
	corpse.In_room = -1
	handle_corpse_condition(corpse, ch)
	if AFF_FLAGGED(ch, AFF_ASHED) {
		act(libc.CString("@WSome ashes fall off the corpse.@n"), 1, ch, nil, nil, TO_ROOM)
		var ashes *obj_data
		if rand_number(1, 3) == 2 {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		} else if rand_number(1, 2) == 2 {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		} else {
			ashes = read_object(1305, VIRTUAL)
			obj_to_room(ashes, ch.In_room)
		}
	}
	if tch != nil {
		if !IS_NPC(tch) && GET_SKILL(tch, SKILL_SURVIVAL) != 0 {
			var skill int = GET_SKILL(tch, SKILL_SURVIVAL)
			if !IS_HUMANOID(ch) && PRF_FLAGGED(tch, PRF_CARVE) && axion_dice(0) < skill {
				send_to_char(tch, libc.CString("The choice edible meat is preserved because of your skill.\r\n"))
				meat = read_object(1612, VIRTUAL)
				obj_to_char(meat, ch)
				var nick [2048]byte
				var nick2 [2048]byte
				var nick3 [2048]byte
				stdio.Sprintf(&nick[0], "@RRaw %s@R Steak@n", GET_NAME(ch))
				stdio.Sprintf(&nick2[0], "Raw %s Steak", ch.Name)
				stdio.Sprintf(&nick3[0], "@wA @Rraw %s@R steak@w is lying here@n", GET_NAME(ch))
				meat.Short_description = libc.StrDup(&nick[0])
				meat.Name = libc.StrDup(&nick2[0])
				meat.Description = libc.StrDup(&nick3[0])
				meat.Value[VAL_ALL_MATERIAL] = 14
			}
		}
	}
	corpse.Type_flag = ITEM_CONTAINER
	corpse.Size = get_size(ch)
	for x = func() int {
		y = 0
		return y
	}(); x < EF_ARRAY_MAX || y < TW_ARRAY_MAX; func() int {
		x++
		return func() int {
			p := &y
			x := *p
			*p++
			return x
		}()
	}() {
		if x < EF_ARRAY_MAX {
			corpse.Extra_flags[x] = 0
		}
		if y < TW_ARRAY_MAX {
			corpse.Wear_flags[y] = 0
		}
	}
	SET_BIT_AR(corpse.Wear_flags[:], ITEM_WEAR_TAKE)
	SET_BIT_AR(corpse.Extra_flags[:], ITEM_NODONATE)
	corpse.Value[VAL_CONTAINER_CAPACITY] = 0
	corpse.Value[VAL_CONTAINER_CORPSE] = 1
	corpse.Value[VAL_CONTAINER_OWNER] = ch.Pfilepos
	corpse.Weight = int64(GET_PC_WEIGHT(ch) + ch.Carry_weight)
	corpse.Cost_per_day = 100000
	if IS_NPC(ch) {
		corpse.Timer = config_info.Play.Max_npc_corpse_time
	} else {
		corpse.Timer = rand_number(config_info.Play.Max_pc_corpse_time/2, config_info.Play.Max_pc_corpse_time)
	}
	SET_BIT_AR(corpse.Extra_flags[:], ITEM_UNIQUE_SAVE)
	if MOB_FLAGGED(ch, MOB_HUSK) {
		for obj = ch.Carrying; obj != nil; obj = obj.Next_content {
			obj_from_char(obj)
			extract_obj(obj)
		}
	}
	if !MOB_FLAGGED(ch, MOB_HUSK) {
		corpse.Contains = ch.Carrying
		for o = corpse.Contains; o != nil; o = o.Next_content {
			o.In_obj = corpse
		}
		object_list_new_owner(corpse, nil)
		var eqdrop int = 0
		_ = eqdrop
		for i = 0; i < NUM_WEARS; i++ {
			if (ch.Equipment[i]) != nil {
				remove_otrigger(ch.Equipment[i], ch)
				obj_to_obj(unequip_char(ch, i), corpse)
				eqdrop = 1
			}
		}
	}
	if ch.Gold > 0 && !MOB_FLAGGED(ch, MOB_HUSK) {
		if IS_NPC(ch) || ch.Desc != nil {
			money = create_money(ch.Gold)
			obj_to_obj(money, corpse)
		}
		ch.Gold = 0
	}
	if !MOB_FLAGGED(ch, MOB_HUSK) {
		ch.Carrying = nil
		ch.Carry_items = 0
		ch.Carry_weight = 0
	}
	obj_to_room(corpse, ch.In_room)
	if !IS_NPC(ch) {
		Crash_rentsave(ch, 0)
	}
}
func loadmap(ch *char_data) {
	var obj *obj_data
	if !IS_NPC(ch) {
		obj = read_object(17, VIRTUAL)
		obj_to_char(obj, ch)
	}
}
func change_alignment(ch *char_data, victim *char_data) {
}
func death_cry(ch *char_data) {
	var door int
	for door = 0; door < NUM_OF_DIRS; door++ {
		if CAN_GO(ch, door) {
			send_to_room(world[ch.In_room].Dir_option[door].To_room, libc.CString("Your blood freezes as you hear someone's death cry.\r\n"))
		}
	}
}
func final_combat_resolve(ch *char_data) {
	var chair *obj_data
	if ch.Sits != nil {
		chair = ch.Sits
		ch.Sits = nil
		chair.Sitting = nil
	}
	if !IS_NPC(ch) && int(ch.Clones) > 0 {
		var clone *char_data = nil
		for clone = character_list; clone != nil; clone = clone.Next {
			if IS_NPC(clone) {
				if GET_MOB_VNUM(clone) == 25 {
					if clone.Original == ch {
						handle_multi_merge(clone)
					}
				}
			}
		}
	}
	if ch.Player_specials.Carrying != nil {
		carry_drop(ch, 2)
	}
	if ch.Player_specials.Carried_by != nil {
		carry_drop(ch.Player_specials.Carried_by, 2)
	}
	if ch.Drag != nil {
		ch.Drag.Dragged = nil
		ch.Drag = nil
	}
	if ch.Dragged != nil {
		ch.Dragged.Drag = nil
		ch.Dragged = nil
	}
	if ch.Grappling != nil {
		ch.Grappling.Grap = -1
		ch.Grappling.Grappled = nil
		ch.Grappling = nil
		ch.Grap = -1
	}
	if ch.Grappled != nil {
		ch.Grappled.Grap = -1
		ch.Grappled.Grappling = nil
		ch.Grappled = nil
		ch.Grap = -1
	}
	if ch.Blocked != nil {
		ch.Blocked.Blocks = nil
		ch.Blocked = nil
	}
	if ch.Blocks != nil {
		ch.Blocks.Blocked = nil
		ch.Blocks = nil
	}
	if ch.Absorbing != nil {
		ch.Absorbing.Absorbby = nil
		ch.Absorbing = nil
	}
	if ch.Absorbby != nil {
		ch.Absorbby.Absorbing = nil
		ch.Absorbby = nil
	}
}
func raw_kill(ch *char_data, killer *char_data) {
	var (
		k    *char_data
		temp *char_data
	)
	if ch.Fighting != nil {
		stop_fighting(ch)
	}
	for ch.Affected != nil {
		affect_remove(ch, ch.Affected)
	}
	for ch.Affectedv != nil {
		affectv_remove(ch, ch.Affectedv)
	}
	if int(ch.Position) != POS_SITTING && int(ch.Position) != POS_SLEEPING && int(ch.Position) != POS_RESTING {
		ch.Position = POS_STANDING
	}
	if killer != nil && !IS_NPC(killer) {
		if !IS_NPC(killer) && !IS_NPC(ch) {
			send_to_imm(libc.CString("[PK] %s killed %s at room [%d]\r\n"), GET_NAME(killer), GET_NAME(ch), GET_ROOM_VNUM(killer.In_room))
		}
		if int(killer.Race) == RACE_SAIYAN && rand_number(1, 2) == 2 || int(killer.Race) != RACE_SAIYAN {
			if rand_number(1, 6) >= 5 && (level_exp(killer, GET_LEVEL(killer)+1)-int(killer.Exp) > 0 || GET_LEVEL(killer) == 100) {
				var psreward float64 = (float64(killer.Aff_abils.Wis) * 0.35)
				if GET_LEVEL(killer) > GET_LEVEL(ch)+5 {
					psreward *= (0.2)
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+2 {
					psreward *= (0.5)
				}
				if int(killer.Race) == RACE_HUMAN || int(killer.Race) == RACE_BIO && ((killer.Genome[0]) == 1 || (killer.Genome[1]) == 1) {
					psreward *= (1.25)
				}
				if int(ch.Race) == RACE_HALFBREED {
					psreward -= (float64(psreward) * 0.4)
				}
				if IS_NPC(ch) && MOB_FLAGGED(ch, MOB_HUSK) && (killer.Player_specials.Class_skill_points[killer.Chclass]) > 50 && int(ch.Race) == RACE_BIO {
					psreward = 0
					send_to_char(killer, libc.CString("@D[@G+0 @BPS @cCapped at 50 for Absorb@D]@n\r\n"))
				} else {
					killer.Player_specials.Class_skill_points[killer.Chclass] += int(psreward)
					send_to_char(killer, libc.CString("@D[@G+%d @BPS@D]@n\r\n"), psreward)
				}
			}
		}
		if int(killer.Race) == RACE_ANDROID && !IS_NPC(killer) && !PLR_FLAGGED(killer, PLR_ABSORB) {
			if PLR_FLAGGED(killer, PLR_REPAIR) {
				if GET_LEVEL(killer) > GET_LEVEL(ch)+15 {
					send_to_char(killer, libc.CString("@D[@G+0 @mUpgrade Point @r-WEAK-@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+10 {
					killer.Upgrade += 3
					send_to_char(killer, libc.CString("@D[@G+3 @mUpgrade Point@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+8 {
					killer.Upgrade += 6
					send_to_char(killer, libc.CString("@D[@G+6 @mUpgrade Points@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+4 {
					killer.Upgrade += 12
					send_to_char(killer, libc.CString("@D[@G+12 @mUpgrade Points@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+2 {
					killer.Upgrade += 16
					send_to_char(killer, libc.CString("@D[@G+16 @mUpgrade Points@D]@n\r\n"))
				} else {
					killer.Upgrade += 28
					send_to_char(killer, libc.CString("@D[@G+28 @mUpgrade Points@D]@n\r\n"))
				}
			} else {
				if GET_LEVEL(killer) > GET_LEVEL(ch)+15 {
					send_to_char(killer, libc.CString("@D[@G+0 @mUpgrade Point @r-WEAK-@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+10 {
					killer.Upgrade += 5
					send_to_char(killer, libc.CString("@D[@G+5 @mUpgrade Point@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+6 {
					killer.Upgrade += 12
					send_to_char(killer, libc.CString("@D[@G+12 @mUpgrade Points@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+4 {
					killer.Upgrade += 18
					send_to_char(killer, libc.CString("@D[@G+18 @mUpgrade Points@D]@n\r\n"))
				} else if GET_LEVEL(killer) > GET_LEVEL(ch)+2 {
					killer.Upgrade += 28
					send_to_char(killer, libc.CString("@D[@G+28 @mUpgrade Points@D]@n\r\n"))
				} else {
					killer.Upgrade += 36
					send_to_char(killer, libc.CString("@D[@G+36 @mUpgrade Points@D]@n\r\n"))
				}
			}
		}
		if death_mtrigger(ch, killer) {
			death_cry(ch)
		}
	} else {
		death_cry(ch)
	}
	update_pos(ch)
	if IS_NPC(ch) && !MOB_FLAGGED(ch, MOB_DUMMY) {
		var shadowed int = 0
		ch.Hit = 0
		if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON1_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON1 = -1
			send_to_room(ch.In_room, libc.CString("@YThe one star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(20, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON2_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON2 = -1
			send_to_room(ch.In_room, libc.CString("@YThe two star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(21, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON3_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON3 = -1
			send_to_room(ch.In_room, libc.CString("@YThe three star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(22, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON4_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON4 = -1
			send_to_room(ch.In_room, libc.CString("@YThe four star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(23, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON5_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON5 = -1
			send_to_room(ch.In_room, libc.CString("@YThe five star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(24, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON6_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON6 = -1
			send_to_room(ch.In_room, libc.CString("@YThe six star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(25, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		} else if IS_NPC(ch) && GET_MOB_VNUM(ch) == SHADOW_DRAGON7_VNUM {
			var obj *obj_data = nil
			SHADOW_DRAGON7 = -1
			send_to_room(ch.In_room, libc.CString("@YThe seven star dragon ball falls to the ground!@n\r\n"))
			obj = read_object(26, VIRTUAL)
			obj_to_room(obj, ch.In_room)
			shadowed = 1
		}
		make_corpse(ch, killer)
		purge_homing(ch)
		extract_char(ch)
		if shadowed == 1 {
			shadow_dragons_live()
		}
	} else if IS_NPC(ch) && MOB_FLAGGED(ch, MOB_DUMMY) {
		ch.Hit = 0
		extract_char(ch)
	} else {
		if !AFF_FLAGGED(ch, AFF_SPIRIT) && !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && (int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 17900 || int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) > 0x464F) {
			if !PLR_FLAGGED(ch, PLR_ABSORBED) {
				make_pcorpse(ch)
				loadmap(ch)
			} else {
				REMOVE_BIT_AR(ch.Act[:], PLR_ABSORBED)
			}
		}
		final_combat_resolve(ch)
		if ch.Fighting != nil {
			stop_fighting(ch)
		}
		for k = combat_list; k != nil; k = temp {
			temp = k.Next_fighting
			if k.Fighting == ch {
				stop_fighting(k)
			}
		}
		if GET_LEVEL(ch) >= 9 && !IS_NPC(ch) && !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && (int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 17900 || int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) > 0x464F) {
			SET_BIT_AR(ch.Affected_by[:], AFF_SPIRIT)
			SET_BIT_AR(ch.Affected_by[:], AFF_ETHEREAL)
			ch.Limb_condition[0] = 100
			ch.Limb_condition[1] = 100
			ch.Limb_condition[2] = 100
			ch.Limb_condition[3] = 100
			SET_BIT_AR(ch.Act[:], PLR_HEAD)
			if !PRF_FLAGGED(ch, PRF_LKEEP) {
				if PLR_FLAGGED(ch, PLR_CLLEG) {
					REMOVE_BIT_AR(ch.Act[:], PLR_CLLEG)
				}
				if PLR_FLAGGED(ch, PLR_CRLEG) {
					REMOVE_BIT_AR(ch.Act[:], PLR_CRLEG)
				}
				if PLR_FLAGGED(ch, PLR_CRARM) {
					REMOVE_BIT_AR(ch.Act[:], PLR_CRARM)
				}
				if PLR_FLAGGED(ch, PLR_CLARM) {
					REMOVE_BIT_AR(ch.Act[:], PLR_CLARM)
				}
			}
			if AFF_FLAGGED(ch, AFF_FROZEN) {
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_FROZEN)
			}
			ch.Hit = 1
			purge_homing(ch)
			if has_group(ch) {
			}
			char_from_room(ch)
			char_to_room(ch, real_room(6000))
			if GET_LEVEL(ch) > 0 && has_group(ch) {
				if ch.Master != nil {
					group_bonus(ch, 1)
				} else {
					group_bonus(ch, 0)
				}
			}
			if AFF_FLAGGED(ch, AFF_BLIND) {
				REMOVE_BIT_AR(ch.Affected_by[:], AFF_BLIND)
			}
			look_at_room(ch.In_room, ch, 0)
			Crash_delete_crashfile(ch)
			update_pos(ch)
			save_char(ch)
		} else if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 17900 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 0x464F {
			ch.Hit = ch.Max_hit - int64(gear_weight(ch))
			char_from_room(ch)
			char_to_room(ch, real_room(17900))
			look_at_room(ch.In_room, ch, 0)
			send_to_char(ch, libc.CString("You wake up and realise that you didn't die, how or why are a mystery.\r\n"))
			Crash_delete_crashfile(ch)
			update_pos(ch)
			save_char(ch)
		} else if ROOM_FLAGGED(ch.In_room, ROOM_PAST) {
			ch.Hit = ch.Max_hit - int64(gear_weight(ch))
			char_from_room(ch)
			char_to_room(ch, real_room(1561))
			look_at_room(ch.In_room, ch, 0)
			send_to_char(ch, libc.CString("You wake up and realise that you died, but only in your mind.\r\n"))
			final_combat_resolve(ch)
			Crash_delete_crashfile(ch)
			update_pos(ch)
			save_char(ch)
		} else if GET_LEVEL(ch) <= 8 && !IS_NPC(ch) {
			ch.Hit = 1
			ch.Mana = ch.Max_mana / 10
			ch.Move = ch.Max_move / 10
			char_from_room(ch)
			if int(ch.Chclass) == CLASS_ROSHI {
				char_to_room(ch, real_room(1130))
			}
			if int(ch.Chclass) == CLASS_KABITO {
				char_to_room(ch, real_room(0x2F42))
			}
			if int(ch.Chclass) == CLASS_NAIL {
				char_to_room(ch, real_room(0x2DA3))
			}
			if int(ch.Chclass) == CLASS_BARDOCK {
				char_to_room(ch, real_room(2268))
			}
			if int(ch.Chclass) == CLASS_KRANE {
				char_to_room(ch, real_room(0x32D1))
			}
			if int(ch.Chclass) == CLASS_TAPION {
				char_to_room(ch, real_room(8231))
			}
			if int(ch.Chclass) == CLASS_PICCOLO {
				char_to_room(ch, real_room(1659))
			}
			if int(ch.Chclass) == CLASS_ANDSIX {
				char_to_room(ch, real_room(1713))
			}
			if int(ch.Chclass) == CLASS_DABURA {
				char_to_room(ch, real_room(6486))
			}
			if int(ch.Chclass) == CLASS_FRIEZA {
				char_to_room(ch, real_room(4282))
			}
			if int(ch.Chclass) == CLASS_GINYU {
				char_to_room(ch, real_room(4289))
			}
			if int(ch.Chclass) == CLASS_JINTO {
				char_to_room(ch, real_room(3499))
			}
			if int(ch.Chclass) == CLASS_TSUNA {
				char_to_room(ch, real_room(15000))
			}
			if int(ch.Chclass) == CLASS_KURZAK {
				char_to_room(ch, real_room(16100))
			}
			look_at_room(ch.In_room, ch, 0)
			ch.Limb_condition[0] = 100
			ch.Limb_condition[1] = 100
			ch.Limb_condition[2] = 100
			ch.Limb_condition[3] = 100
			SET_BIT_AR(ch.Act[:], PLR_HEAD)
			Crash_delete_crashfile(ch)
			update_pos(ch)
			save_char(ch)
			send_to_char(ch, libc.CString("\r\n@RYou should beware, when you reach level 9, you will actually die. So you\r\nshould learn to be more careful. Since when you die past that point and\r\nactually reach the afterlife you need to realise that being revived will\r\nnot be very easy. So treat your character's dying with as much care as\r\npossible.@n\r\n"))
		}
		if int(ch.Race) == RACE_ANDROID && !PLR_FLAGGED(ch, PLR_ABSORB) && !AFF_FLAGGED(ch, AFF_SPIRIT) && !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && (int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) < 17900 || int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) > 0x464F) && ch.Upgrade > 5 {
			var loss int = ch.Upgrade / 5
			ch.Upgrade -= loss
			send_to_char(ch, libc.CString("@rYou lose @R%s@r upgrade points!@n\r\n"), add_commas(int64(loss)))
		}
		WAIT_STATE(ch, config_info.Ticks.Pulse_violence*(int(1000000/OPT_USEC)))
	}
}
func die(ch *char_data, killer *char_data) {
	if !IS_NPC(ch) {
		if PLR_FLAGGED(ch, PLR_HEALT) {
			REMOVE_BIT_AR(ch.Act[:], PLR_HEALT)
		}
		if (int(ch.Race) == RACE_MAJIN || int(ch.Race) == RACE_BIO) && (float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.75 || PLR_FLAGGED(ch, PLR_SELFD2) && float64(ch.Lifeforce) >= float64(GET_LIFEMAX(ch))*0.5) {
			ch.Lifeforce = -1
			ch.Hit = 1
			SET_BIT_AR(ch.Act[:], PLR_GOOP)
			ch.Gooptime = 32
			return
		}
		if PLR_FLAGGED(ch, PLR_IMMORTAL) {
			act(libc.CString("@c$n@w disappears right before dying. $n appears to be immortal.@n"), 1, ch, nil, nil, TO_CHAR)
			act(libc.CString("@c$n@w disappears right before dying. $n appears to be immortal.@n."), 1, ch, nil, nil, TO_ROOM)
			ch.Hit = 1
			ch.Mana = 1
			ch.Move = 1
			null_affect(ch, AFF_POISON)
			if int(ch.Player_specials.Conditions[HUNGER]) >= 0 {
				ch.Player_specials.Conditions[HUNGER] = 48
			}
			if int(ch.Player_specials.Conditions[THIRST]) >= 0 {
				ch.Player_specials.Conditions[THIRST] = 48
			}
			if ch.Fighting != nil {
				stop_fighting(ch)
			}
			ch.Position = POS_SITTING
			char_from_room(ch)
			if int(ch.Chclass) == CLASS_ROSHI {
				char_to_room(ch, real_room(1130))
			}
			if int(ch.Chclass) == CLASS_KABITO {
				char_to_room(ch, real_room(0x2F42))
			}
			if int(ch.Chclass) == CLASS_NAIL {
				char_to_room(ch, real_room(0x2DA3))
			}
			if int(ch.Chclass) == CLASS_BARDOCK {
				char_to_room(ch, real_room(2268))
			}
			if int(ch.Chclass) == CLASS_KRANE {
				char_to_room(ch, real_room(0x32D1))
			}
			if int(ch.Chclass) == CLASS_TAPION {
				char_to_room(ch, real_room(8231))
			}
			if int(ch.Chclass) == CLASS_PICCOLO {
				char_to_room(ch, real_room(1659))
			}
			if int(ch.Chclass) == CLASS_ANDSIX {
				char_to_room(ch, real_room(1713))
			}
			if int(ch.Chclass) == CLASS_DABURA {
				char_to_room(ch, real_room(6486))
			}
			if int(ch.Chclass) == CLASS_FRIEZA {
				char_to_room(ch, real_room(4282))
			}
			if int(ch.Chclass) == CLASS_GINYU {
				char_to_room(ch, real_room(4289))
			}
			if int(ch.Chclass) == CLASS_JINTO {
				char_to_room(ch, real_room(3499))
			}
			if int(ch.Chclass) == CLASS_KURZAK {
				char_to_room(ch, real_room(16100))
			}
			return
		}
		REMOVE_BIT_AR(ch.Act[:], PLR_KILLER)
		REMOVE_BIT_AR(ch.Act[:], PLR_THIEF)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_KNOCKED)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_SLEEP)
		REMOVE_BIT_AR(ch.Affected_by[:], AFF_PARALYZE)
		if !AFF_FLAGGED(ch, AFF_SPIRIT) && !ROOM_FLAGGED(ch.In_room, ROOM_PAST) && GET_LEVEL(ch) > 8 {
			if int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) >= 2002 && int(libc.BoolToInt(GET_ROOM_VNUM(ch.In_room))) <= 2011 {
				ch.Deathtime = libc.GetTime(nil)
			} else if ROOM_FLAGGED(ch.In_room, ROOM_AL) || ROOM_FLAGGED(ch.In_room, ROOM_HELL) {
				send_to_char(ch, libc.CString("Your soul is saved from destruction by King Yemma. Why? Who knows.\r\n"))
			} else if IN_ARENA(ch) {
				cleanup_arena_watch(ch)
				if killer != nil {
					cleanup_arena_watch(killer)
					send_to_all(libc.CString("@R%s@r manages to defeat @R%s@r in the Arena!@n\r\n"), GET_NAME(killer), GET_NAME(ch))
					char_from_room(killer)
					char_to_room(killer, real_room(0x45D3))
					look_at_room(killer.In_room, killer, 0)
					final_combat_resolve(killer)
					final_combat_resolve(ch)
				} else {
					send_to_all(libc.CString("@R%s@r dies in the water of the Arena and is disqualified!@n\r\n"), GET_NAME(ch))
				}
				char_from_room(ch)
				char_to_room(ch, real_room(0x45D3))
				ch.Hit = 1
				look_at_room(ch.In_room, ch, 0)
				final_combat_resolve(ch)
				return
			} else {
				if killer != nil && IS_NPC(killer) {
					ch.Deathtime = libc.GetTime(nil) + 28800
					ch.Dcount += 1
				} else if killer != nil && !IS_NPC(killer) {
					ch.Deathtime = libc.GetTime(nil) + 1123200
					SET_BIT_AR(ch.Act[:], PLR_PDEATH)
					ch.Dcount += 1
				} else {
					if ch.Dcount <= 0 {
						ch.Deathtime = libc.GetTime(nil) + 28800
					} else if ch.Dcount <= 1 {
						ch.Deathtime = libc.GetTime(nil) + 43200
					} else if ch.Dcount <= 3 {
						ch.Deathtime = libc.GetTime(nil) + 86400
					} else if ch.Dcount <= 5 {
						ch.Deathtime = libc.GetTime(nil) + 0x2A300
					} else if ch.Dcount > 5 {
						ch.Deathtime = libc.GetTime(nil) + 604800
					}
					ch.Dcount += 1
				}
			}
			if int(ch.Player_specials.Conditions[HUNGER]) >= 0 {
				ch.Player_specials.Conditions[HUNGER] = 48
			}
			if int(ch.Player_specials.Conditions[THIRST]) >= 0 {
				ch.Player_specials.Conditions[THIRST] = 48
			}
		}
	}
	raw_kill(ch, killer)
}
func perform_group_gain(ch *char_data, base int, victim *char_data) {
	var share float64
	if IN_ARENA(ch) {
		return
	}
	share = float64(MIN(2000000, int64(base*GET_LEVEL(ch))))
	if !IS_NPC(ch) {
		if GET_LEVEL(ch) >= 100 && float64(ch.Max_hit)*0.025 >= float64(victim.Max_hit) {
			share *= (0.05)
		} else if float64(ch.Max_hit)*0.025 >= float64(victim.Max_hit) {
			share = 1
		} else if float64(ch.Max_hit)*0.05 >= float64(victim.Max_hit) {
			share *= (0.05)
		} else if float64(ch.Max_hit)*0.1 >= float64(victim.Max_hit) {
			share *= (0.1)
		} else if float64(ch.Max_hit)*0.15 >= float64(victim.Max_hit) {
			share *= (0.15)
		} else if float64(ch.Max_hit)*0.25 >= float64(victim.Max_hit) {
			share *= (0.25)
		} else if float64(ch.Max_hit)*0.5 >= float64(victim.Max_hit) {
			share *= (0.5)
		} else if float64(ch.Max_hit)*0.9 >= float64(victim.Max_hit) {
			share *= (0.65)
		} else if ch.Max_hit >= victim.Max_hit {
			share *= (0.7)
		}
	}
	if victim.Lasthit != 0 && victim.Lasthit != int(ch.Idnum) {
		var (
			f       *follow_type
			checkit int = 0
		)
		for f = ch.Followers; f != nil; f = f.Next {
			if AFF_FLAGGED(f.Follower, AFF_GROUP) && victim.Lasthit == int(f.Follower.Idnum) {
				checkit = 1
			}
		}
		if checkit == 0 && ch.Master != nil && int(ch.Master.Idnum) == victim.Lasthit {
			checkit = 1
		}
		if checkit == 0 && ch.Master != nil {
			var master *char_data = ch.Master
			for f = master.Followers; f != nil; f = f.Next {
				if f.Follower != ch {
					if AFF_FLAGGED(f.Follower, AFF_GROUP) && victim.Lasthit == int(f.Follower.Idnum) {
						checkit = 1
					}
				}
			}
		}
		if checkit == 0 {
			send_to_char(ch, libc.CString("@RYou didn't do most of the work for this kill.@n\r\n"))
			share = 1
		}
	}
	if IS_NPC(victim) && MOB_FLAGGED(victim, MOB_HUSK) {
		share /= 10
	}
	if (ch.Bonuses[BONUS_PRODIGY]) > 0 {
		share = (share + float64(share)*0.25)
	}
	if int(ch.Race) == RACE_SAIYAN {
		share = (float64(share) + float64(share)*0.5)
	}
	if int(ch.Race) == RACE_HALFBREED {
		share = (float64(share) + float64(share)*0.4)
	}
	if int(ch.Race) == RACE_ICER {
		share = (float64(share) - float64(share)*0.2)
	}
	if (ch.Bonuses[BONUS_LOYAL]) > 0 && ch.Master != nil {
		share += (float64(share) * 0.2)
	}
	if ch.Master != nil && ch.Master != ch {
		share += (float64(share) * 0.15)
	}
	if MOB_FLAGGED(victim, MOB_KNOWKAIO) {
		share += (float64(share) * 0.25)
	}
	ch.Combatexpertise += 1
	if float64((ch.Combatexpertise+1)/20) > float64(share)*0.16 {
		share += (float64(share) * 0.16)
	} else {
		share += ((float64(share) * 0.02) * float64((ch.Combatexpertise+1)/20))
	}
	if int(libc.BoolToInt(group_bonus(ch, 2))) == 2 {
		send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G+2 PS!@D]@n\r\n"))
		ch.Player_specials.Class_skill_points[ch.Chclass] += 2
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 3 {
		send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G+5%s Exp!@D]@n\r\n"), "%")
		share += (float64(share) * 0.05)
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 5 {
		ch.Mana += int64(float64(ch.Max_mana) * 0.04)
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		}
		send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G4%s Ki Regenerated!@D]@n\r\n"), "%")
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 6 {
		ch.Mana += int64(float64(ch.Max_mana) * 0.02)
		ch.Move += int64(float64(ch.Max_move) * 0.02)
		ch.Hit += int64(float64(gear_pl(ch)) * 0.02)
		if ch.Mana > ch.Max_mana {
			ch.Mana = ch.Max_mana
		} else if ch.Hit > gear_pl(ch) {
			ch.Hit = gear_pl(ch)
		} else if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G2%s PL/ST/Ki Regenerated!@D]@n\r\n"), "%")
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 7 && int(ch.Race) == RACE_ANDROID {
		if PLR_FLAGGED(ch.Master, PLR_ABSORB) {
			ch.Mana += int64(float64(ch.Max_mana) * 0.02)
			ch.Move += int64(float64(ch.Max_move) * 0.02)
			if ch.Mana > ch.Max_mana {
				ch.Mana = ch.Max_mana
			} else if ch.Hit > ch.Max_hit {
				ch.Hit = ch.Max_hit
			} else if ch.Move > ch.Max_move {
				ch.Move = ch.Max_move
			}
			send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G2%s PL/ST/Ki Recovered!@D]@n\r\n"), "%")
		} else if PLR_FLAGGED(ch.Master, PLR_REPAIR) {
			ch.Hit += int64(float64(gear_pl(ch)) * 0.02)
			if ch.Hit > gear_pl(ch) {
				ch.Hit = gear_pl(ch)
			}
			send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G5%s PL Repaired@D]@n\r\n"), "%")
		} else if PLR_FLAGGED(ch.Master, PLR_SENSEM) && !PLR_FLAGGED(ch, PLR_ABSORB) {
			ch.Upgrade += 5
			send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G+5 @mUpgrade Points@D]@n\r\n"))
		}
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 11 {
		ch.Move += int64(float64(ch.Max_move) * 0.04)
		if ch.Move > ch.Max_move {
			ch.Move = ch.Max_move
		}
		send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G4%s ST Regenerated!@D]@n\r\n"), "%")
	} else if int(libc.BoolToInt(group_bonus(ch, 2))) == 13 {
		if ch.Master.Starphase == 1 {
			share += (float64(share) * 0.05)
			send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G+5%s Exp!@D]@n\r\n"), "%")
		} else if ch.Master.Starphase == 2 {
			share += (float64(share) * 0.1)
			send_to_char(ch, libc.CString("You receive a bonus from your group's leader! @D[@G+10%s Exp!@D]@n\r\n"), "%")
		}
	}
	share = float64(gear_exp(ch, int64(share)))
	if share > 1 {
		send_to_char(ch, libc.CString("You receive your share of experience -- %s points.\r\n"), add_commas(int64(share)))
	} else {
		send_to_char(ch, libc.CString("You receive your share of experience -- one measly little point!\r\n"))
	}
	gain_exp(ch, int64(share))
}
func group_gain(ch *char_data, victim *char_data) {
	var (
		tot_levels  int
		tot_members int
		tot_gain    int64
		base        int64
		k           *char_data
		f           *follow_type
	)
	if (func() *char_data {
		k = ch.Master
		return k
	}()) == nil {
		k = ch
	}
	if AFF_FLAGGED(k, AFF_GROUP) && k.In_room == ch.In_room {
		tot_levels = GET_LEVEL(k)
		tot_members = 1
	} else {
		tot_levels = 0
		tot_members = 0
	}
	for f = k.Followers; f != nil; f = f.Next {
		if AFF_FLAGGED(f.Follower, AFF_GROUP) && f.Follower.In_room == ch.In_room {
			if !IS_WEIGHTED(f.Follower) {
				tot_levels += GET_LEVEL(f.Follower)
				tot_members++
			} else if float64(gear_pl(f.Follower)) >= float64(gear_pl(ch))*0.5 {
				tot_levels += GET_LEVEL(f.Follower)
				tot_members++
			}
		}
	}
	if tot_members == 1 || IN_ARENA(ch) {
		solo_gain(ch, victim)
		return
	}
	tot_gain = victim.Exp + int64(tot_members) - 1
	if !IS_NPC(victim) {
		tot_gain = MIN(int64(config_info.Play.Max_exp_loss*2), tot_gain)
	}
	if tot_levels >= 1 {
		base = MAX(1, tot_gain/int64(tot_levels))
		var perc int = tot_members * 20
		if perc >= 80 {
			perc = 60
		}
		base += (base / 100) * int64(perc)
	} else {
		base = 0
	}
	if AFF_FLAGGED(k, AFF_GROUP) && k.In_room == ch.In_room {
		if !IS_WEIGHTED(k) {
			perform_group_gain(k, int(base), victim)
		} else if k != ch && float64(gear_pl(k)) >= float64(gear_pl(ch))*0.5 {
			perform_group_gain(k, int(base), victim)
		} else if k == ch && float64(gear_pl(k)) >= float64(ch.Max_hit)*0.5 {
			perform_group_gain(k, int(base), victim)
		} else {
			if k == ch {
				send_to_char(ch, libc.CString("You can not group gain while your powerlevel is weighted down more than half of your max.\r\n"))
			} else {
				send_to_char(ch, libc.CString("You can not group gain while your powerlevel is weighted down more than half of the leader's adjusted powerlevel.\r\n"))
			}
		}
	}
	for f = k.Followers; f != nil; f = f.Next {
		if AFF_FLAGGED(f.Follower, AFF_GROUP) && f.Follower.In_room == ch.In_room {
			if float64(gear_pl(f.Follower)) >= float64(ch.Max_hit)*0.5 {
				perform_group_gain(f.Follower, int(base), victim)
			}
		}
	}
}
func solo_gain(ch *char_data, victim *char_data) {
	if IS_NPC(ch) {
		if ch.Original != nil {
			ch = ch.Original
		}
	}
	var exp float64
	exp = float64(MIN(2000000, victim.Exp))
	if !IS_NPC(ch) {
		if GET_LEVEL(ch) >= 100 && float64(ch.Max_hit)*0.025 >= float64(victim.Max_hit) {
			exp *= (0.04)
		} else if float64(ch.Max_hit)*0.025 >= float64(victim.Max_hit) {
			exp = 1
		} else if float64(ch.Max_hit)*0.05 >= float64(victim.Max_hit) {
			exp *= (0.05)
		} else if float64(ch.Max_hit)*0.1 >= float64(victim.Max_hit) {
			exp *= (0.1)
		} else if float64(ch.Max_hit)*0.15 >= float64(victim.Max_hit) {
			exp *= (0.15)
		} else if float64(ch.Max_hit)*0.25 >= float64(victim.Max_hit) {
			exp *= (0.25)
		} else if float64(ch.Max_hit)*0.5 >= float64(victim.Max_hit) {
			exp *= (0.5)
		}
	}
	if victim.Lasthit != 0 && victim.Lasthit != int(ch.Idnum) {
		send_to_char(ch, libc.CString("@RYou didn't do most of the work for this victory.@n\r\n"))
		exp = 1
	}
	if IS_NPC(victim) && MOB_FLAGGED(victim, MOB_HUSK) {
		exp /= 10
	}
	if (ch.Bonuses[BONUS_PRODIGY]) > 0 {
		exp = (float64(exp) + float64(exp)*0.25)
	}
	if int(ch.Race) == RACE_SAIYAN {
		exp = (float64(exp) + float64(exp)*0.5)
	}
	if int(ch.Race) == RACE_HALFBREED {
		exp = (float64(exp) + float64(exp)*0.4)
	}
	if int(ch.Race) == RACE_ICER {
		exp = (float64(exp) - float64(exp)*0.2)
	}
	if MOB_FLAGGED(victim, MOB_KNOWKAIO) {
		exp += (float64(exp) * 0.25)
	}
	exp = float64(gear_exp(ch, int64(exp)))
	exp = float64(MAX(int64(exp), 1))
	if exp > 1 {
		send_to_char(ch, libc.CString("You receive %s experience points.\r\n"), add_commas(int64(exp)))
	} else {
		send_to_char(ch, libc.CString("You receive one lousy experience point. That fight was hardly worth it...\r\n"))
	}
	if !IS_NPC(ch) {
		gain_exp(ch, int64(exp))
	}
	if IS_NPC(victim) {
		gain_exp(victim, int64(-exp))
	}
	if !IS_NPC(victim) {
		exp = exp / 5
		gain_exp(victim, int64(-exp))
	}
}
