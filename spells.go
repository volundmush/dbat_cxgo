package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unicode"
	"unsafe"
)

const DEFAULT_STAFF_LVL = 12
const DEFAULT_WAND_LVL = 12
const CAST_UNDEFINED = -1
const CAST_SPELL = 0
const CAST_POTION = 1
const CAST_WAND = 2
const CAST_STAFF = 3
const CAST_SCROLL = 4
const CAST_STRIKE = 5
const MAG_DAMAGE = 1
const MAG_AFFECTS = 2
const MAG_UNAFFECTS = 4
const MAG_POINTS = 8
const MAG_ALTER_OBJS = 16
const MAG_GROUPS = 32
const MAG_MASSES = 64
const MAG_AREAS = 128
const MAG_SUMMONS = 256
const MAG_CREATIONS = 512
const MAG_MANUAL = 1024
const MAG_AFFECTSV = 2048
const MAG_ACTION_FREE = 4096
const MAG_ACTION_PARTIAL = 8192
const MAG_ACTION_FULL = 0x4000
const MAG_NEXTSTRIKE = 0x8000
const MAG_TOUCH_MELEE = 0x10000
const MAG_TOUCH_RANGED = 0x20000
const MAGSAVE_FORT = 1
const MAGSAVE_REFLEX = 2
const MAGSAVE_WILL = 4
const MAGSAVE_HALF = 8
const MAGSAVE_NONE = 16
const MAGSAVE_PARTIAL = 32
const MAGCOMP_DIVINE_FOCUS = 1
const MAGCOMP_EXP_COST = 2
const MAGCOMP_FOCUS = 4
const MAGCOMP_MATERIAL = 8
const MAGCOMP_SOMATIC = 16
const MAGCOMP_VERBAL = 32
const TYPE_UNDEFINED = -1
const SPELL_RESERVED_DBC = 0
const SPELL_MAGE_ARMOR = 1
const SPELL_TELEPORT = 2
const SPELL_BLESS = 3
const SPELL_BLINDNESS = 4
const SPELL_BURNING_HANDS = 5
const SPELL_CALL_LIGHTNING = 6
const SPELL_CHARM = 7
const SPELL_CHILL_TOUCH = 8
const SPELL_COLOR_SPRAY = 10
const SPELL_CONTROL_WEATHER = 11
const SPELL_CREATE_FOOD = 12
const SPELL_CREATE_WATER = 13
const SPELL_REMOVE_BLINDNESS = 14
const SPELL_CURE_CRITIC = 15
const SPELL_CURE_LIGHT = 16
const SPELL_BANE = 17
const SPELL_DETECT_ALIGN = 18
const SPELL_SEE_INVIS = 19
const SPELL_DETECT_MAGIC = 20
const SPELL_DETECT_POISON = 21
const SPELL_DISPEL_EVIL = 22
const SPELL_EARTHQUAKE = 23
const SPELL_ENCHANT_WEAPON = 24
const SPELL_ENERGY_DRAIN = 25
const SPELL_FIREBALL = 26
const SPELL_HARM = 27
const SPELL_HEAL = 28
const SPELL_INVISIBLE = 29
const SPELL_LIGHTNING_BOLT = 30
const SPELL_LOCATE_OBJECT = 31
const SPELL_MAGIC_MISSILE = 32
const SPELL_POISON = 33
const SPELL_PROT_FROM_EVIL = 34
const SPELL_REMOVE_CURSE = 35
const SPELL_SANCTUARY = 36
const SPELL_SHOCKING_GRASP = 37
const SPELL_SLEEP = 38
const SPELL_BULL_STRENGTH = 39
const SPELL_SUMMON = 40
const SPELL_VENTRILOQUATE = 41
const SPELL_WORD_OF_RECALL = 42
const SPELL_NEUTRALIZE_POISON = 43
const SPELL_SENSE_LIFE = 44
const SPELL_ANIMATE_DEAD = 45
const SPELL_DISPEL_GOOD = 46
const SPELL_GROUP_ARMOR = 47
const SPELL_MASS_HEAL = 48
const SPELL_GROUP_RECALL = 49
const SPELL_DARKVISION = 50
const SPELL_WATERWALK = 51
const SPELL_PORTAL = 52
const SPELL_PARALYZE = 53
const SPELL_INFLICT_LIGHT = 54
const SPELL_INFLICT_CRITIC = 55
const SPELL_IDENTIFY = 56
const SPELL_FAERIE_FIRE = 57
const ABIL_TURNING = 58
const ABIL_LAY_HANDS = 59
const SPELL_RESISTANCE = 60
const SPELL_ACID_SPLASH = 61
const SPELL_DAZE = 62
const SPELL_FLARE = 63
const SPELL_RAY_OF_FROST = 64
const SPELL_DISRUPT_UNDEAD = 65
const SPELL_LESSER_GLOBE_OF_INVUL = 66
const SPELL_STONESKIN = 67
const SPELL_MINOR_CREATION = 68
const SPELL_SUMMON_MONSTER_I = 69
const SPELL_SUMMON_MONSTER_II = 70
const SPELL_SUMMON_MONSTER_III = 71
const SPELL_SUMMON_MONSTER_IV = 72
const SPELL_SUMMON_MONSTER_V = 73
const SPELL_SUMMON_MONSTER_VI = 74
const SPELL_SUMMON_MONSTER_VII = 75
const SPELL_SUMMON_MONSTER_VIII = 76
const SPELL_SUMMON_MONSTER_IX = 77
const SPELL_FIRE_SHIELD = 78
const SPELL_ICE_STORM = 79
const SPELL_SHOUT = 80
const SPELL_FEAR = 81
const SPELL_CLOUDKILL = 82
const SPELL_MAJOR_CREATION = 83
const SPELL_HOLD_MONSTER = 84
const SPELL_CONE_OF_COLD = 85
const SPELL_ANIMAL_GROWTH = 86
const SPELL_BALEFUL_POLYMORPH = 87
const SPELL_PASSWALL = 88
const SPELL_BESTOW_CURSE = 89
const SPELL_SENSU = 90
const SPELL_HAYASA = 91
const MIN_LANGUAGES = 141
const SKILL_LANG_COMMON = 141
const SKILL_LANG_ELVEN = 142
const SKILL_LANG_GNOME = 143
const SKILL_LANG_DWARVEN = 144
const SKILL_LANG_HALFLING = 145
const SKILL_LANG_ORC = 146
const SKILL_LANG_DRUIDIC = 147
const SKILL_LANG_DRACONIC = 148
const MAX_LANGUAGES = 148
const SKILL_WP_UNARMED = 179
const SPELL_FIRE_BREATH = 202
const SPELL_GAS_BREATH = 203
const SPELL_FROST_BREATH = 204
const SPELL_ACID_BREATH = 205
const SPELL_LIGHTNING_BREATH = 206
const MAX_SPELLS = 90
const SPELL_DG_AFFECT = 298
const TYPE_HIT = 300
const TYPE_STING = 301
const TYPE_WHIP = 302
const TYPE_SLASH = 303
const TYPE_BITE = 304
const TYPE_BLUDGEON = 305
const TYPE_CRUSH = 306
const TYPE_POUND = 307
const TYPE_CLAW = 308
const TYPE_MAUL = 309
const TYPE_THRASH = 310
const TYPE_PIERCE = 311
const TYPE_BLAST = 312
const TYPE_PUNCH = 313
const TYPE_STAB = 314
const TYPE_SUFFERING = 399
const SKILL_FLEX = 400
const SKILL_GENIUS = 401
const SKILL_SOLARF = 402
const SKILL_MIGHT = 403
const SKILL_BALANCE = 404
const SKILL_BUILD = 405
const SKILL_TSKIN = 406
const SKILL_CONCENTRATION = 407
const SKILL_KAIOKEN = 408
const SKILL_SPOT = 409
const SKILL_FIRST_AID = 410
const SKILL_DISGUISE = 411
const SKILL_ESCAPE_ARTIST = 412
const SKILL_APPRAISE = 413
const SKILL_HEAL = 414
const SKILL_FORGERY = 415
const SKILL_HIDE = 416
const SKILL_BLESS = 417
const SKILL_CURSE = 418
const SKILL_LISTEN = 419
const SKILL_EAVESDROP = 420
const SKILL_POISON = 421
const SKILL_CURE = 422
const SKILL_OPEN_LOCK = 423
const SKILL_VIGOR = 424
const SKILL_REGENERATE = 425
const SKILL_KEEN = 426
const SKILL_SEARCH = 427
const SKILL_MOVE_SILENTLY = 428
const SKILL_ABSORB = 429
const SKILL_SLEIGHT_OF_HAND = 430
const SKILL_INGEST = 431
const SKILL_REPAIR = 432
const SKILL_SENSE = 433
const SKILL_SURVIVAL = 434
const SKILL_YOIK = 435
const SKILL_CREATE = 436
const SKILL_SPIT = 437
const SKILL_POTENTIAL = 438
const SKILL_TELEPATHY = 439
const SKILL_RENZO = 440
const SKILL_MASENKO = 441
const SKILL_DODONPA = 442
const SKILL_BARRIER = 443
const SKILL_GALIKGUN = 444
const SKILL_THROW = 445
const SKILL_DODGE = 446
const SKILL_PARRY = 447
const SKILL_BLOCK = 448
const SKILL_PUNCH = 449
const SKILL_KICK = 450
const SKILL_ELBOW = 451
const SKILL_KNEE = 452
const SKILL_ROUNDHOUSE = 453
const SKILL_UPPERCUT = 454
const SKILL_SLAM = 455
const SKILL_HEELDROP = 456
const SKILL_FOCUS = 457
const SKILL_KIBALL = 458
const SKILL_KIBLAST = 459
const SKILL_BEAM = 460
const SKILL_TSUIHIDAN = 461
const SKILL_SHOGEKIHA = 462
const SKILL_ZANZOKEN = 463
const SKILL_KAMEHAMEHA = 464
const SKILL_DAGGER = 465
const SKILL_SWORD = 466
const SKILL_CLUB = 467
const SKILL_SPEAR = 468
const SKILL_GUN = 469
const SKILL_BRAWL = 470
const SKILL_INSTANTT = 471
const SKILL_DEATHBEAM = 472
const SKILL_ERASER = 473
const SKILL_TSLASH = 474
const SKILL_PSYBLAST = 475
const SKILL_HONOO = 476
const SKILL_DUALBEAM = 477
const SKILL_ROGAFUFUKEN = 478
const SKILL_POSE = 479
const SKILL_BAKUHATSUHA = 480
const SKILL_KIENZAN = 481
const SKILL_TRIBEAM = 482
const SKILL_SBC = 483
const SKILL_FINALFLASH = 484
const SKILL_CRUSHER = 485
const SKILL_DDSLASH = 486
const SKILL_PBARRAGE = 487
const SKILL_HELLFLASH = 488
const SKILL_HELLSPEAR = 489
const SKILL_KAKUSANHA = 490
const SKILL_HASSHUKEN = 491
const SKILL_SCATTER = 492
const SKILL_BIGBANG = 493
const SKILL_PSLASH = 494
const SKILL_DEATHBALL = 495
const SKILL_SPIRITBALL = 496
const SKILL_GENKIDAMA = 497
const SKILL_GENOCIDE = 498
const SKILL_DUALWIELD = 499
const SKILL_KURA = 500
const SKILL_TAILWHIP = 501
const SKILL_KOUSENGAN = 502
const SKILL_TAISHA = 503
const SKILL_PARALYZE = 505
const SKILL_INFUSE = 506
const SKILL_ROLL = 507
const SKILL_TRIP = 508
const SKILL_GRAPPLE = 509
const SKILL_WSPIKE = 510
const SKILL_SELFD = 511
const SKILL_SPIRAL = 512
const SKILL_BREAKER = 513
const SKILL_ENLIGHTEN = 514
const SKILL_COMMUNE = 515
const SKILL_MIMIC = 516
const SKILL_WRAZOR = 517
const SKILL_KOTEIRU = 518
const SKILL_DIMIZU = 519
const SKILL_HYOGA_KABE = 520
const SKILL_WELLSPRING = 521
const SKILL_AQUA_BARRIER = 522
const SKILL_WARP = 523
const SKILL_HSPIRAL = 524
const SKILL_ARMOR = 525
const SKILL_FIRESHIELD = 526
const SKILL_COOKING = 527
const SKILL_SEISHOU = 528
const SKILL_SILK = 529
const SKILL_BASH = 530
const SKILL_HEADBUTT = 531
const SKILL_ENSNARE = 532
const SKILL_STARNOVA = 533
const SKILL_PURSUIT = 534
const SKILL_ZEN = 535
const SKILL_SUNDER = 536
const SKILL_WITHER = 537
const SKILL_TWOHAND = 538
const SKILL_STYLE = 539
const SKILL_METAMORPH = 540
const SKILL_HEALGLOW = 541
const SKILL_RUNIC = 542
const SKILL_EXTRACT = 543
const SKILL_GARDENING = 544
const SKILL_ENERGIZE = 545
const SKILL_MALICE = 549
const SKILL_HAYASA = 550
const SKILL_HANDLING = 551
const SKILL_MYSTICMUSIC = 552
const SKILL_LIGHTGRENADE = 553
const SKILL_MULTIFORM = 554
const SKILL_SPIRITCONTROL = 555
const SKILL_BALEFIRE = 556
const SKILL_BLESSEDHAMMER = 557
const ART_STUNNING_FIST = 1000
const ART_WHOLENESS_OF_BODY = 1001
const ART_ABUNDANT_STEP = 1002
const ART_QUIVERING_PALM = 1003
const ART_EMPTY_BODY = 1004
const SAVING_FORTITUDE = 0
const SAVING_REFLEX = 1
const SAVING_WILL = 2
const SAVING_OBJ_IMPACT = 0
const SAVING_OBJ_HEAT = 1
const SAVING_OBJ_COLD = 2
const SAVING_OBJ_BREATH = 3
const SAVING_OBJ_SPELL = 4
const TAR_IGNORE = 1
const TAR_CHAR_ROOM = 2
const TAR_CHAR_WORLD = 4
const TAR_FIGHT_SELF = 8
const TAR_FIGHT_VICT = 16
const TAR_SELF_ONLY = 32
const TAR_NOT_SELF = 64
const TAR_OBJ_INV = 128
const TAR_OBJ_ROOM = 256
const TAR_OBJ_WORLD = 512
const TAR_OBJ_EQUIP = 1024
const SKTYPE_NONE = 0
const SKTYPE_SPELL = 1
const SKTYPE_SKILL = 2
const SKTYPE_LANG = 4
const SKTYPE_WEAPON = 8
const SKTYPE_ART = 16
const SKFLAG_NEEDTRAIN = 1
const SKFLAG_STRMOD = 2
const SKFLAG_DEXMOD = 4
const SKFLAG_CONMOD = 8
const SKFLAG_INTMOD = 16
const SKFLAG_WISMOD = 32
const SKFLAG_CHAMOD = 64
const SKFLAG_ARMORBAD = 128
const SKFLAG_ARMORALL = 256
const SKFLAG_TIER1 = 512
const SKFLAG_TIER2 = 1024
const SKFLAG_TIER3 = 2048
const SKFLAG_TIER4 = 4096
const SKFLAG_TIER5 = 8192
const SKLEARN_CANT = 0
const SKLEARN_CROSSCLASS = 1
const SKLEARN_CLASS = 2
const SKLEARN_BOOL = 3
const SPELL_TYPE_SPELL = 0
const SPELL_TYPE_POTION = 1
const SPELL_TYPE_WAND = 2
const SPELL_TYPE_STAFF = 3
const SPELL_TYPE_SCROLL = 4

type spell_info_type struct {
	Min_position    int8
	Mana_min        int
	Mana_max        int
	Mana_change     int
	Ki_min          int
	Ki_max          int
	Ki_change       int
	Min_level       [31]int
	Routines        int
	Violent         int8
	Targets         int
	Name            *byte
	Wear_off_msg    *byte
	Race_can_learn  [24]int
	Skilltype       int
	Flags           int
	Save_flags      int
	Comp_flags      int
	Can_learn_skill [31]int8
	Spell_level     int
	School          int
	Domain          int
}
type attack_hit_type struct {
	Singular *byte
	Plural   *byte
}

const SUMMON_FAIL = "You failed.\r\n"

func spell_create_water(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var water int
	if ch == nil || obj == nil {
		return
	}
	if int(obj.Type_flag) == ITEM_DRINKCON {
		if (obj.Value[VAL_DRINKCON_LIQUID]) != LIQ_WATER && (obj.Value[VAL_DRINKCON_HOWFULL]) != 0 {
			name_from_drinkcon(obj)
			obj.Value[VAL_DRINKCON_LIQUID] = LIQ_SLIME
			name_to_drinkcon(obj, LIQ_SLIME)
		} else {
			water = MAX((obj.Value[VAL_DRINKCON_CAPACITY])-(obj.Value[VAL_DRINKCON_HOWFULL]), 0)
			if water > 0 {
				if (obj.Value[VAL_DRINKCON_HOWFULL]) >= 0 {
					name_from_drinkcon(obj)
				}
				obj.Value[VAL_DRINKCON_LIQUID] = LIQ_WATER
				obj.Value[VAL_DRINKCON_HOWFULL] += water
				name_to_drinkcon(obj, LIQ_WATER)
				weight_change_object(obj, water)
				act(libc.CString("$p is filled."), FALSE, ch, obj, nil, TO_CHAR)
			}
		}
	}
}
func spell_recall(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if victim == nil || IS_NPC(victim) {
		return
	}
	act(libc.CString("$n disappears."), TRUE, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, real_room(config_info.Room_nums.Mortal_start_room))
	act(libc.CString("$n appears in the middle of the room."), TRUE, victim, nil, nil, TO_ROOM)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_teleport(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var to_room room_rnum
	if victim == nil || IS_NPC(victim) {
		return
	}
	for {
		to_room = room_rnum(rand_number(0, int(top_of_world)))
		if !ROOM_FLAGGED(to_room, bitvector_t(int32(int(ROOM_PRIVATE|ROOM_DEATH)|ROOM_GODROOM))) {
			break
		}
	}
	act(libc.CString("$n slowly fades out of existence and is gone."), FALSE, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, to_room)
	act(libc.CString("$n slowly fades into existence."), FALSE, victim, nil, nil, TO_ROOM)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_summon(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if ch == nil || victim == nil {
		return
	}
	if GET_LEVEL(victim) > level+3 {
		send_to_char(ch, libc.CString("%s"), SUMMON_FAIL)
		return
	}
	if config_info.Play.Pk_allowed == 0 {
		if MOB_FLAGGED(victim, MOB_AGGRESSIVE) {
			act(libc.CString("As the words escape your lips and $N travels\r\nthrough time and space towards you, you realize that $E is\r\naggressive and might harm you, so you wisely send $M back."), FALSE, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			return
		}
		if !IS_NPC(victim) && !PRF_FLAGGED(victim, PRF_SUMMONABLE) && !PLR_FLAGGED(victim, PLR_KILLER) {
			send_to_char(victim, libc.CString("%s just tried to summon you to: %s.\r\n%s failed because you have summon protection on.\r\nType NOSUMMON to allow other players to summon you.\r\n"), GET_NAME(ch), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Name, func() string {
				if int(ch.Sex) == SEX_MALE {
					return "He"
				}
				return "She"
			}())
			send_to_char(ch, libc.CString("You failed because %s has summon protection on.\r\n"), GET_NAME(victim))
			mudlog(BRF, ADMLVL_IMMORT, TRUE, libc.CString("%s failed summoning %s to %s."), GET_NAME(ch), GET_NAME(victim), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Name)
			return
		}
	}
	if MOB_FLAGGED(victim, MOB_NOSUMMON) || IS_NPC(victim) && mag_newsaves(ch, victim, SPELL_SUMMON, level, int(ch.Aff_abils.Intel)) != 0 {
		send_to_char(ch, libc.CString("%s"), SUMMON_FAIL)
		return
	}
	act(libc.CString("$n disappears suddenly."), TRUE, victim, nil, nil, TO_ROOM)
	char_from_room(victim)
	char_to_room(victim, ch.In_room)
	act(libc.CString("$n arrives suddenly."), TRUE, victim, nil, nil, TO_ROOM)
	act(libc.CString("$n has summoned you!"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
	look_at_room(victim.In_room, victim, 0)
	entry_memory_mtrigger(victim)
	greet_mtrigger(victim, -1)
	greet_memory_mtrigger(victim)
}
func spell_locate_object(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		i    *obj_data
		name [2048]byte
		j    int
	)
	if obj == nil {
		send_to_char(ch, libc.CString("You sense nothing.\r\n"))
		return
	}
	strlcpy(&name[0], fname(obj.Name), uint64(2048))
	j = level / 2
	for i = object_list; i != nil && j > 0; i = i.Next {
		if isname(&name[0], i.Name) == 0 {
			continue
		}
		send_to_char(ch, libc.CString("%c%s"), unicode.ToUpper(rune(*i.Short_description)), (*byte)(unsafe.Add(unsafe.Pointer(i.Short_description), 1)))
		if i.Carried_by != nil {
			send_to_char(ch, libc.CString(" is being carried by %s.\r\n"), PERS(i.Carried_by, ch))
		} else if i.In_room != room_rnum(-1) {
			send_to_char(ch, libc.CString(" is in %s.\r\n"), (*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(i.In_room)))).Name)
		} else if i.In_obj != nil {
			send_to_char(ch, libc.CString(" is in %s.\r\n"), i.In_obj.Short_description)
		} else if i.Worn_by != nil {
			send_to_char(ch, libc.CString(" is being worn by %s.\r\n"), PERS(i.Worn_by, ch))
		} else {
			send_to_char(ch, libc.CString("'s location is uncertain.\r\n"))
		}
		j--
	}
	if j == level/2 {
		send_to_char(ch, libc.CString("You sense nothing.\r\n"))
	}
}
func spell_charm(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var af affected_type
	if victim == nil || ch == nil {
		return
	}
	if victim == ch {
		send_to_char(ch, libc.CString("You like yourself even better!\r\n"))
	} else if !IS_NPC(victim) && !PRF_FLAGGED(victim, PRF_SUMMONABLE) {
		send_to_char(ch, libc.CString("You fail because SUMMON protection is on!\r\n"))
	} else if AFF_FLAGGED(victim, AFF_SANCTUARY) {
		send_to_char(ch, libc.CString("Your victim is protected by sanctuary!\r\n"))
	} else if MOB_FLAGGED(victim, MOB_NOCHARM) {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else if AFF_FLAGGED(ch, AFF_CHARM) {
		send_to_char(ch, libc.CString("You can't have any followers of your own!\r\n"))
	} else if AFF_FLAGGED(victim, AFF_CHARM) || level < GET_LEVEL(victim) {
		send_to_char(ch, libc.CString("You fail.\r\n"))
	} else if config_info.Play.Pk_allowed == 0 && !IS_NPC(victim) {
		send_to_char(ch, libc.CString("You fail - shouldn't be doing it anyway.\r\n"))
	} else if int(victim.Race) == RACE_SAIYAN && rand_number(1, 100) <= 90 {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else if circle_follow(victim, ch) {
		send_to_char(ch, libc.CString("Sorry, following in circles cannot be allowed.\r\n"))
	} else if mag_newsaves(ch, victim, SPELL_CHARM, level, int(ch.Aff_abils.Intel)) != 0 {
		send_to_char(ch, libc.CString("Your victim resists!\r\n"))
	} else {
		if victim.Master != nil {
			stop_follower(victim)
		}
		add_follower(victim, ch)
		victim.Master_id = ch.Idnum
		af.Type = SPELL_CHARM
		af.Duration = 24 * 2
		if int(ch.Aff_abils.Cha) != 0 {
			af.Duration *= int16(ch.Aff_abils.Cha)
		}
		if int(victim.Aff_abils.Intel) != 0 {
			af.Duration /= int16(victim.Aff_abils.Intel)
		}
		af.Modifier = 0
		af.Location = 0
		af.Bitvector = AFF_CHARM
		affect_to_char(victim, &af)
		act(libc.CString("Isn't $n just such a nice fellow?"), FALSE, ch, nil, unsafe.Pointer(victim), TO_VICT)
		if IS_NPC(victim) {
			victim.Act[int(MOB_SPEC/32)] &= bitvector_t(int32(^(1 << (int(MOB_SPEC % 32)))))
		}
	}
}
func spell_identify(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		i     int
		found int
		len_  uint64
	)
	if obj != nil {
		var (
			bitbuf [64936]byte
			buf2   [64936]byte
		)
		sprinttype(int(obj.Type_flag), item_types[:], &bitbuf[0], uint64(64936))
		send_to_char(ch, libc.CString("You feel informed:\r\nObject '%s', Item type: %s\r\n"), obj.Short_description, &bitbuf[0])
		if obj.Bitvector != nil {
			sprintbitarray(obj.Bitvector[:], affected_bits[:], AF_ARRAY_MAX, &bitbuf[0])
			send_to_char(ch, libc.CString("Item will give you following abilities:  %s\r\n"), &bitbuf[0])
		}
		sprintbitarray(obj.Extra_flags[:], extra_bits[:], EF_ARRAY_MAX, &bitbuf[0])
		send_to_char(ch, libc.CString("Item is: %s\r\n"), &bitbuf[0])
		send_to_char(ch, libc.CString("Weight: %lld, Value: %d, Rent: %d, Min Level: %d\r\n"), obj.Weight, obj.Cost, obj.Cost_per_day, obj.Level)
		switch obj.Type_flag {
		case ITEM_SCROLL:
			fallthrough
		case ITEM_POTION:
			len_ = uint64(func() int {
				i = 0
				return i
			}())
			if (obj.Value[VAL_SCROLL_SPELL1]) >= 1 {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL1]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			if (obj.Value[VAL_SCROLL_SPELL2]) >= 1 && len_ < uint64(64936) {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL2]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			if (obj.Value[VAL_SCROLL_SPELL3]) >= 1 && len_ < uint64(64936) {
				i = stdio.Snprintf(&bitbuf[len_], int(64936-uintptr(len_)), " %s", skill_name(obj.Value[VAL_SCROLL_SPELL3]))
				if i >= 0 {
					len_ += uint64(i)
				}
			}
			send_to_char(ch, libc.CString("This %s casts: %s\r\n"), item_types[int(obj.Type_flag)], &bitbuf[0])
		case ITEM_WAND:
			fallthrough
		case ITEM_STAFF:
			send_to_char(ch, libc.CString("This %s casts: %s\r\nIt has %d maximum charge%s and %d remaining.\r\n"), item_types[int(obj.Type_flag)], skill_name(obj.Value[VAL_WAND_SPELL]), obj.Value[VAL_WAND_MAXCHARGES], func() string {
				if (obj.Value[VAL_WAND_MAXCHARGES]) == 1 {
					return ""
				}
				return "s"
			}(), obj.Value[VAL_WAND_CHARGES])
		case ITEM_WEAPON:
			send_to_char(ch, libc.CString("Damage Dice is '%dD%d' for an average per-round damage of %.1f.\r\n"), obj.Value[VAL_WEAPON_DAMDICE], obj.Value[VAL_WEAPON_DAMSIZE], (float64((obj.Value[VAL_WEAPON_DAMSIZE])+1)/2.0)*float64(obj.Value[VAL_WEAPON_DAMDICE]))
		case ITEM_ARMOR:
			send_to_char(ch, libc.CString("AC-apply is %.1f\r\n"), (float32(obj.Value[VAL_ARMOR_APPLYAC]))/10)
		}
		found = FALSE
		for i = 0; i < MAX_OBJ_AFFECT; i++ {
			if obj.Affected[i].Location != APPLY_NONE && obj.Affected[i].Modifier != 0 {
				if found == 0 {
					send_to_char(ch, libc.CString("Can affect you as :\r\n"))
					found = TRUE
				}
				sprinttype(obj.Affected[i].Location, apply_types[:], &bitbuf[0], uint64(64936))
				switch obj.Affected[i].Location {
				case APPLY_FEAT:
					stdio.Snprintf(&buf2[0], int(64936), " (%s)", feat_list[obj.Affected[i].Specific].Name)
				case APPLY_SKILL:
					stdio.Snprintf(&buf2[0], int(64936), " (%s)", spell_info[obj.Affected[i].Specific].Name)
				default:
					buf2[0] = 0
				}
				send_to_char(ch, libc.CString("   Affects: %s%s By %d\r\n"), &bitbuf[0], &buf2[0], obj.Affected[i].Modifier)
			}
		}
	}
}
func spell_enchant_weapon(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var i int
	if ch == nil || obj == nil {
		return
	}
	if int(obj.Type_flag) != ITEM_WEAPON || OBJ_FLAGGED(obj, ITEM_MAGIC) {
		return
	}
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location != APPLY_NONE {
			return
		}
	}
	obj.Extra_flags[int(ITEM_MAGIC/32)] |= bitvector_t(int32(1 << (int(ITEM_MAGIC % 32))))
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location == APPLY_NONE {
			obj.Affected[i].Location = APPLY_ACCURACY
			obj.Affected[i].Modifier = int(libc.BoolToInt(level >= 18)) + 1
			break
		}
	}
	for i = 0; i < MAX_OBJ_AFFECT; i++ {
		if obj.Affected[i].Location == APPLY_NONE {
			obj.Affected[i].Location = APPLY_DAMAGE
			obj.Affected[i].Modifier = int(libc.BoolToInt(level >= 20)) + 1
			break
		}
	}
	if IS_GOOD(ch) {
		obj.Extra_flags[int(ITEM_ANTI_EVIL/32)] |= bitvector_t(int32(1 << (int(ITEM_ANTI_EVIL % 32))))
		act(libc.CString("$p glows blue."), FALSE, ch, obj, nil, TO_CHAR)
	} else if IS_EVIL(ch) {
		obj.Extra_flags[int(ITEM_ANTI_GOOD/32)] |= bitvector_t(int32(1 << (int(ITEM_ANTI_GOOD % 32))))
		act(libc.CString("$p glows red."), FALSE, ch, obj, nil, TO_CHAR)
	} else {
		act(libc.CString("$p glows yellow."), FALSE, ch, obj, nil, TO_CHAR)
	}
}
func spell_detect_poison(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	if victim != nil {
		if victim == ch {
			if AFF_FLAGGED(victim, AFF_POISON) {
				send_to_char(ch, libc.CString("You can sense poison in your blood.\r\n"))
			} else {
				send_to_char(ch, libc.CString("You feel healthy.\r\n"))
			}
		} else {
			if AFF_FLAGGED(victim, AFF_POISON) {
				act(libc.CString("You sense that $E is poisoned."), FALSE, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			} else {
				act(libc.CString("You sense that $E is healthy."), FALSE, ch, nil, unsafe.Pointer(victim), TO_CHAR)
			}
		}
	}
	if obj != nil {
		switch obj.Type_flag {
		case ITEM_DRINKCON:
			fallthrough
		case ITEM_FOUNTAIN:
			fallthrough
		case ITEM_FOOD:
			if (obj.Value[VAL_FOOD_POISON]) != 0 {
				act(libc.CString("You sense that $p has been contaminated."), FALSE, ch, obj, nil, TO_CHAR)
			} else {
				act(libc.CString("You sense that $p is safe for consumption."), FALSE, ch, obj, nil, TO_CHAR)
			}
		default:
			send_to_char(ch, libc.CString("You sense that it should not be consumed.\r\n"))
		}
	}
}
func spell_portal(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		portal  *obj_data
		tportal *obj_data
		rm      *room_data
	)
	rm = (*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(victim.In_room)))
	if ch == nil || victim == nil {
		return
	}
	if can_edit_zone(ch, rm.Zone) == 0 && ZONE_FLAGGED(rm.Zone, ZONE_QUEST) {
		send_to_char(ch, libc.CString("That target is in a quest zone.\r\n"))
		return
	}
	if ZONE_FLAGGED(rm.Zone, ZONE_CLOSED) && ch.Admlevel < ADMLVL_IMMORT {
		send_to_char(ch, libc.CString("That target is in a closed zone.\r\n"))
		return
	}
	if ZONE_FLAGGED(rm.Zone, ZONE_NOIMMORT) && ch.Admlevel < ADMLVL_GRGOD {
		send_to_char(ch, libc.CString("That target is in a zone closed to all.\r\n"))
		return
	}
	portal = read_object(portal_object, VIRTUAL)
	if victim.In_room != room_rnum(-1) && victim.In_room <= top_of_world {
		portal.Value[VAL_PORTAL_DEST] = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(victim.In_room)))).Number)
	} else {
		portal.Value[VAL_PORTAL_DEST] = -1
	}
	portal.Value[VAL_PORTAL_HEALTH] = 100
	portal.Value[VAL_PORTAL_MAXHEALTH] = 100
	portal.Timer = level / 10
	add_unique_id(portal)
	obj_to_room(portal, ch.In_room)
	act(libc.CString("$n opens a portal in thin air."), TRUE, ch, nil, nil, TO_ROOM)
	act(libc.CString("You open a portal out of thin air."), TRUE, ch, nil, nil, TO_CHAR)
	tportal = read_object(portal_object, VIRTUAL)
	if ch.In_room != room_rnum(-1) && ch.In_room <= top_of_world {
		tportal.Value[VAL_PORTAL_DEST] = int((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(ch.In_room)))).Number)
	} else {
		tportal.Value[VAL_PORTAL_DEST] = -1
	}
	tportal.Value[VAL_PORTAL_HEALTH] = 100
	tportal.Value[VAL_PORTAL_MAXHEALTH] = 100
	tportal.Timer = level / 10
	add_unique_id(portal)
	obj_to_room(tportal, victim.In_room)
	act(libc.CString("A shimmering portal appears out of thin air."), TRUE, victim, nil, nil, TO_ROOM)
	act(libc.CString("A shimmering portal opens here for you."), TRUE, victim, nil, nil, TO_CHAR)
}
func art_abundant_step(level int, ch *char_data, victim *char_data, obj *obj_data, arg *byte) {
	var (
		steps    int
		i        int = 0
		j        int
		rep      int
		max      int
		r        room_rnum
		nextroom room_rnum
		buf      [2048]byte
		tc       int8
		p        *byte
	)
	steps = 0
	r = ch.In_room
	p = arg
	max = ((ch.Chclasses[CLASS_KABITO])+(ch.Epicclasses[CLASS_KABITO]))/2 + 10
	for p != nil && *p != 0 && !unicode.IsDigit(rune(*p)) && !libc.IsAlpha(rune(*p)) {
		p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
	}
	if p == nil || *p == 0 {
		send_to_char(ch, libc.CString("You must give directions from your current location. Examples:\r\n  w w nw n e\r\n  2w nw n e\r\n"))
		return
	}
	for *p != 0 {
		for *p != 0 && !unicode.IsDigit(rune(*p)) && !libc.IsAlpha(rune(*p)) {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
		if unicode.IsDigit(rune(*p)) {
			rep = libc.Atoi(libc.GoString(p))
			for unicode.IsDigit(rune(*p)) {
				p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
			}
		} else {
			rep = 1
		}
		if libc.IsAlpha(rune(*p)) {
			for i = 0; libc.IsAlpha(rune(*p)); func() *byte {
				i++
				return func() *byte {
					p := &p
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()
			}() {
				buf[i] = byte(int8(unicode.ToLower(rune(*p))))
			}
			j = i
			tc = int8(buf[i])
			buf[i] = 0
			for i = 1; libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command_pointer) == libc.FuncAddr(do_move) && libc.StrCmp((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Sort_as, &buf[0]) != 0; i++ {
			}
			if libc.FuncAddr((*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Command_pointer) == libc.FuncAddr(do_move) {
				i = (*(*command_info)(unsafe.Add(unsafe.Pointer(complete_cmd_info), unsafe.Sizeof(command_info{})*uintptr(i)))).Subcmd - 1
			} else {
				i = -1
			}
			buf[j] = byte(tc)
		}
		if i > -1 {
			for func() int {
				p := &rep
				x := *p
				*p--
				return x
			}() != 0 {
				if func() int {
					p := &steps
					*p++
					return *p
				}() > max {
					break
				}
				if ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(r)))).Dir_option[i]) == nil {
					send_to_char(ch, libc.CString("Invalid step. Skipping.\r\n"))
					break
				}
				nextroom = ((*(*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(r)))).Dir_option[i]).To_room
				if nextroom == room_rnum(-1) {
					break
				}
				r = nextroom
			}
		}
		if steps > max {
			break
		}
	}
	send_to_char(ch, libc.CString("Your will bends reality as you travel through the ethereal plane.\r\n"))
	act(libc.CString("$n is suddenly absent."), TRUE, ch, nil, nil, TO_ROOM)
	char_from_room(ch)
	char_to_room(ch, r)
	act(libc.CString("$n is suddenly present."), TRUE, ch, nil, nil, TO_ROOM)
	look_at_room(ch.In_room, ch, 0)
	return
}
func roll_skill(ch *char_data, snum int) int {
	var (
		roll  int
		skval int
		i     int
	)
	if !IS_NPC(ch) {
		skval = GET_SKILL(ch, snum)
		if SKILL_SPOT == snum {
			if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 4 || (ch.Genome[1]) == 4) {
				skval += 5
			}
		} else if SKILL_HIDE == snum {
			if AFF_FLAGGED(ch, AFF_LIQUEFIED) {
				skval += 5
			} else if int(ch.Race) == RACE_MUTANT && ((ch.Genome[0]) == 5 || (ch.Genome[1]) == 5) {
				skval += 10
			}
		}
	} else if IS_NPC(ch) {
		var numb int = 0
		if GET_LEVEL(ch) <= 10 {
			numb = rand_number(15, 30)
		}
		if GET_LEVEL(ch) <= 20 {
			numb = rand_number(20, 40)
		}
		if GET_LEVEL(ch) <= 30 {
			numb = rand_number(40, 60)
		}
		if GET_LEVEL(ch) <= 60 {
			numb = rand_number(60, 80)
		}
		if GET_LEVEL(ch) <= 80 {
			numb = rand_number(70, 90)
		}
		if GET_LEVEL(ch) <= 90 {
			numb = rand_number(80, 95)
		}
		if GET_LEVEL(ch) <= 100 {
			numb = rand_number(90, 100)
		}
		skval = numb
	}
	if snum == SKILL_SPOT && GET_SKILL(ch, SKILL_LISTEN) != 0 {
		skval += GET_SKILL(ch, SKILL_LISTEN) / 10
	}
	if snum < 0 || snum >= SKILL_TABLE_SIZE {
		return 0
	}
	if (spell_info[snum].Skilltype & (1 << 0)) != 0 {
		for func() int {
			i = 0
			return func() int {
				roll = 0
				return roll
			}()
		}(); i < NUM_CLASSES; i++ {
			if ((ch.Chclasses[i])+(ch.Epicclasses[i])) != 0 && spell_info[snum].Min_level[i] < ((ch.Chclasses[i])+(ch.Epicclasses[i])) {
				roll += (ch.Chclasses[i]) + (ch.Epicclasses[i])
			}
		}
		return roll + rand_number(1, 20)
	} else if (spell_info[snum].Skilltype & (1 << 1)) != 0 {
		if skval == 0 && (spell_info[snum].Flags&(1<<0)) != 0 {
			return -1
		} else {
			roll = skval
			if (spell_info[snum].Flags & (1 << 1)) != 0 {
				roll += int(ability_mod_value(int(ch.Aff_abils.Str)))
			}
			if (spell_info[snum].Flags & (1 << 2)) != 0 {
				roll += int(dex_mod_capped(ch))
			}
			if (spell_info[snum].Flags & (1 << 3)) != 0 {
				roll += int(ability_mod_value(int(ch.Aff_abils.Con)))
			}
			if (spell_info[snum].Flags & (1 << 4)) != 0 {
				roll += int(ability_mod_value(int(ch.Aff_abils.Intel)))
			}
			if (spell_info[snum].Flags & (1 << 5)) != 0 {
				roll += int(ability_mod_value(int(ch.Aff_abils.Wis)))
			}
			if (spell_info[snum].Flags & (1 << 6)) != 0 {
				roll += int(ability_mod_value(int(ch.Aff_abils.Cha)))
			}
			if (spell_info[snum].Flags & (1 << 8)) != 0 {
				roll -= int(ch.Armorcheckall)
			} else if (spell_info[snum].Flags & (1 << 7)) != 0 {
				roll -= int(ch.Armorcheck)
			}
			return roll + rand_number(1, 20)
		}
	} else {
		basic_mud_log(libc.CString("Trying to roll uncategorized skill/spell #%d for %s"), snum, GET_NAME(ch))
		return 0
	}
}
func roll_resisted(actor *char_data, sact int, resistor *char_data, sres int) int {
	return int(libc.BoolToInt(roll_skill(actor, sact) >= roll_skill(resistor, sres)))
}
