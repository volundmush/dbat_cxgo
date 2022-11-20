package main

import "github.com/gotranspile/cxgo/runtime/libc"

var race_names [25]*byte = [25]*byte{libc.CString("human"), libc.CString("saiyan"), libc.CString("icer"), libc.CString("konatsu"), libc.CString("namekian"), libc.CString("mutant"), libc.CString("kanassan"), libc.CString("halfbreed"), libc.CString("bioandroid"), libc.CString("android"), libc.CString("demon"), libc.CString("majin"), libc.CString("kai"), libc.CString("truffle"), libc.CString("hoshijin"), libc.CString("animal"), libc.CString("saiba"), libc.CString("serpent"), libc.CString("ogre"), libc.CString("yardratian"), libc.CString("arlian"), libc.CString("dragon"), libc.CString("mechanical"), libc.CString("spirit"), libc.CString("\n")}
var race_abbrevs [25]*byte = [25]*byte{libc.CString("Hum"), libc.CString("Sai"), libc.CString("Ice"), libc.CString("Kon"), libc.CString("Nam"), libc.CString("Mut"), libc.CString("Kan"), libc.CString("H-B"), libc.CString("Bio"), libc.CString("And"), libc.CString("Dem"), libc.CString("Maj"), libc.CString("Kai"), libc.CString("Tru"), libc.CString("Hos"), libc.CString("Ict"), libc.CString("Sab"), libc.CString("Ser"), libc.CString("Trl"), libc.CString("Dra"), libc.CString("Arl"), libc.CString("Mnd"), libc.CString("Mec"), libc.CString("Spi"), libc.CString("\n")}
var pc_race_types [25]*byte = [25]*byte{libc.CString("Human"), libc.CString("Saiyan"), libc.CString("Icer"), libc.CString("Konatsu"), libc.CString("Namekian"), libc.CString("Mutant"), libc.CString("Kanassan"), libc.CString("Halfbreed"), libc.CString("Bioandroid"), libc.CString("Android"), libc.CString("Demon"), libc.CString("Majin"), libc.CString("Kai"), libc.CString("Truffle"), libc.CString("Hoshijin"), libc.CString("animal"), libc.CString("Saiba"), libc.CString("Serpent"), libc.CString("Ogre"), libc.CString("Yardratian"), libc.CString("Arlian"), libc.CString("Dragon"), libc.CString("mechanical"), libc.CString("Spirit"), libc.CString("\n")}
var d_race_types [25]*byte = [25]*byte{libc.CString("A Disguised Human"), libc.CString("A Disguised Saiyan"), libc.CString("A Disguised Icer"), libc.CString("A Disguised Konatsu"), libc.CString("A Disguised Namekian"), libc.CString("A Disguised Mutant"), libc.CString("A Disguised Kanassan"), libc.CString("A Disguised Halfbreed"), libc.CString("A Disguised Bioandroid"), libc.CString("A Disguised Android"), libc.CString("A Disguised Demon"), libc.CString("A Disguised Majin"), libc.CString("A Disguised Kai"), libc.CString("A Disguised Truffle"), libc.CString("A Disguised Hoshijin"), libc.CString("A Disguised Animal"), libc.CString("Saiba"), libc.CString("Serpent"), libc.CString("Ogre"), libc.CString("Yardratian"), libc.CString("A Disguised Arlian"), libc.CString("Dragon"), libc.CString("mechanical"), libc.CString("Spirit"), libc.CString("\n")}
var race_ok_gender [3][24]bool = [3][24]bool{{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, true, false, false, false}, {true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, true, false, false, false}, {true, true, true, true, false, true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, true, false, false, false}}
var race_display [24]*byte = [24]*byte{libc.CString("@B1@W) @cHuman\r\n"), libc.CString("@B2@W) @cSaiyan\r\n"), libc.CString("@B3@W) @cIcer\r\n"), libc.CString("@B4@W) @cKonatsu\r\n"), libc.CString("@B5@W) @cNamekian\r\n"), libc.CString("@B6@W) @cMutant\r\n"), libc.CString("@B7@W) @cKanassan\r\n"), libc.CString("@B8@W) @cHalf Breed\r\n"), libc.CString("@B9@W) @cBio-Android\r\n"), libc.CString("@B10@W) @cAndroid\r\n"), libc.CString("@B11@W) @cDemon\r\n"), libc.CString("@B12@W) @cMajin\r\n"), libc.CString("@B13@W) @cKai\r\n"), libc.CString("@B14@W) @cTruffle\r\n"), libc.CString("@B15@W) @cHoshijin\r\n"), libc.CString("@B16@W) @YArlian\r\n"), libc.CString("@B17@W) @GAnimal\r\n"), libc.CString("@B18@W) @MSaiba\r\n"), libc.CString("@B19@W) @BSerpent\r\n"), libc.CString("@B20@W) @ROgre\r\n"), libc.CString("@B21@W) @CYardratian\r\n"), libc.CString("@B22@W) @GLizardfolk\r\n"), libc.CString("@B23@W) @GMechanical\r\n"), libc.CString("@B24@W) @MSpirit\r\n")}

func parse_race(ch *char_data, arg int) int {
	var race int = -1
	switch arg {
	case 1:
		race = RACE_HUMAN
	case 2:
		if ch.Desc.Rpp >= 60 {
			race = RACE_SAIYAN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else if ch.Desc.Rbank >= 60 {
			race = RACE_SAIYAN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			race = -1
		}
	case 3:
		race = RACE_ICER
	case 4:
		race = RACE_KONATSU
	case 5:
		race = RACE_NAMEK
	case 6:
		race = RACE_MUTANT
	case 7:
		race = RACE_KANASSAN
	case 8:
		race = RACE_HALFBREED
	case 9:
		if ch.Desc.Rpp >= 35 {
			race = RACE_BIO
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else if ch.Desc.Rbank >= 35 {
			race = RACE_BIO
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			race = -1
		}
	case 10:
		race = RACE_ANDROID
	case 11:
		race = RACE_DEMON
	case 12:
		if ch.Desc.Rpp >= 55 {
			race = RACE_MAJIN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else if ch.Desc.Rbank >= 55 {
			race = RACE_MAJIN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			race = -1
		}
	case 13:
		race = RACE_KAI
	case 14:
		race = RACE_TRUFFLE
	case 15:
		race = RACE_HOSHIJIN
		if ch.Desc.Rpp >= 30 {
			race = RACE_HOSHIJIN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else if ch.Desc.Rbank >= 30 {
			race = RACE_HOSHIJIN
			userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
		} else {
			race = -1
		}
	case 16:
		race = RACE_ARLIAN
	case 17:
		race = RACE_ANIMAL
	case 18:
		race = RACE_ORC
	case 19:
		race = RACE_SNAKE
	case 20:
		race = RACE_TROLL
	case 21:
		race = RACE_MINOTAUR
	case 22:
		race = RACE_LIZARDFOLK
	case 23:
		race = RACE_WARHOST
	case 24:
		race = RACE_FAERIE
	default:
		race = -1
	}
	if race >= 0 && race < NUM_RACES {
		if !race_ok_gender[int(ch.Sex)][race] {
			race = -1
		}
	}
	return race
}

var racial_ability_mods [25][6]int = [25][6]int{{}, {0, -2, 0, 0, 2, 0}, {-2, 2, 0, 0, 0, 0}, {0, 2, 0, 0, 0, -2}, {}, {-2, 0, 0, 0, 2, 0}, {0, -2, 2, 0, 2, 2}, {2, 0, -2, 0, 0, -2}, {}, {}, {}, {}, {}, {14, 8, -4, 0, -2, -4}, {-2, 0, 0, 0, 2, -2}, {}, {4, 0, -2, -2, 0, -2}, {}, {12, 12, -4, -2, 4, -4}, {8, 4, -4, 0, 0, -2}, {-4, -2, 0, 0, 2, 0}, {}, {}, {}, {}}

func racial_ability_modifiers(ch *char_data) {
	var chrace int = 0
	_ = chrace
	if int(ch.Race) >= NUM_RACES || int(ch.Race) < 0 {
		basic_mud_log(libc.CString("SYSERR: Unknown race %d in racial_ability_modifiers"), ch.Race)
	} else {
		chrace = int(ch.Race)
	}
}

var hw_info [24]struct {
	Height    [3]int
	Heightdie int
	Weight    [3]int
	Weightfac int
} = [24]struct {
	Height    [3]int
	Heightdie int
	Weight    [3]int
	Weightfac int
}{{Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{140, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{100, 111, 95}, Heightdie: 10, Weight: [3]int{17, 18, 16}, Weightfac: 18}, {Height: [3]int{121, 124, 109}, Heightdie: 20, Weight: [3]int{52, 59, 45}, Weightfac: 125}, {Height: [3]int{137, 140, 135}, Heightdie: 20, Weight: [3]int{40, 45, 36}, Weightfac: 89}, {Height: [3]int{141, 150, 140}, Heightdie: 10, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{135, 135, 135}, Heightdie: 15, Weight: [3]int{37, 39, 36}, Weightfac: 63}, {Height: [3]int{141, 147, 135}, Heightdie: 30, Weight: [3]int{59, 68, 50}, Weightfac: 125}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{40, 50, 45}, Heightdie: 16, Weight: [3]int{16, 24, 9}, Weightfac: 8}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}, {Height: [3]int{141, 147, 135}, Heightdie: 26, Weight: [3]int{46, 54, 39}, Weightfac: 89}}

func set_height_and_weight_by_race(ch *char_data) {
	var (
		race int
		sex  int
		mod  int
	)
	if !IS_NPC(ch) {
		return
	}
	race = int(ch.Race)
	sex = int(ch.Sex)
	if sex < SEX_NEUTRAL || sex >= NUM_SEX {
		basic_mud_log(libc.CString("Invalid gender in set_height_and_weight_by_race: %d"), sex)
		sex = SEX_NEUTRAL
	}
	if race <= -1 || race >= NUM_RACES {
		basic_mud_log(libc.CString("Invalid gender in set_height_and_weight_by_race: %d"), ch.Sex)
		race = int(-1 + 1)
	}
	mod = dice(2, hw_info[race].Heightdie)
	ch.Height = uint8(int8(hw_info[race].Height[sex] + mod))
	mod *= hw_info[race].Weightfac
	mod /= 100
	ch.Weight = uint8(int8(hw_info[race].Weight[sex] + mod))
}
func invalid_race(ch *char_data, obj *obj_data) bool {
	if ch.Admlevel >= ADMLVL_IMMORT {
		return false
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_HUMAN) && int(ch.Race) == RACE_HUMAN {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_SAIYAN) && int(ch.Race) == RACE_SAIYAN {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_ICER) && int(ch.Race) == RACE_ICER {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ANTI_KONATSU) && int(ch.Race) == RACE_KONATSU {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_HUMAN) && int(ch.Race) != RACE_HUMAN {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_ICER) && int(ch.Race) != RACE_ICER {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_SAIYAN) && int(ch.Race) != RACE_SAIYAN {
		return true
	}
	if OBJ_FLAGGED(obj, ITEM_ONLY_KONATSU) && int(ch.Race) != RACE_KONATSU {
		return true
	}
	return false
}

var race_def_sizetable [25]int = [25]int{0: SIZE_MEDIUM, 1: SIZE_MEDIUM, 2: SIZE_MEDIUM, 3: SIZE_MEDIUM, 4: SIZE_MEDIUM, 5: SIZE_MEDIUM, 6: SIZE_MEDIUM, 7: SIZE_MEDIUM, 8: SIZE_MEDIUM, 9: SIZE_MEDIUM, 10: SIZE_MEDIUM, 11: SIZE_MEDIUM, 12: SIZE_MEDIUM, 13: SIZE_SMALL, 14: SIZE_MEDIUM, 15: SIZE_FINE, 16: SIZE_LARGE, 17: SIZE_MEDIUM, 18: SIZE_LARGE, 19: SIZE_MEDIUM, 20: SIZE_MEDIUM, 21: SIZE_MEDIUM, 22: SIZE_MEDIUM, 23: SIZE_TINY}

func get_size(ch *char_data) int {
	var racenum int
	if ch.Size != int(-1) {
		return ch.Size
	} else {
		racenum = int(ch.Race)
		if racenum < 0 || racenum >= NUM_RACES {
			return SIZE_MEDIUM
		}
		return func() int {
			p := &ch.Size
			ch.Size = race_def_sizetable[racenum]
			return *p
		}()
	}
}

var size_bonus_table [9]int = [9]int{8, 4, 2, 1, 0, -1, -2, -4, -8}

func get_size_bonus(sz int) int {
	if sz < 0 || sz >= NUM_SIZES {
		sz = SIZE_MEDIUM
	}
	return size_bonus_table[sz]
}
func wield_type(chsize int, weap *obj_data) int {
	if int(weap.Type_flag) != ITEM_WEAPON {
		if OBJ_FLAGGED(weap, ITEM_2H) {
			return WIELD_TWOHAND
		}
		return WIELD_ONEHAND
	} else if chsize > weap.Size {
		return WIELD_LIGHT
	} else if chsize == weap.Size {
		return WIELD_ONEHAND
	} else if chsize == weap.Size-1 {
		return WIELD_TWOHAND
	} else if chsize < weap.Size-1 {
		return WIELD_NONE
	} else {
		basic_mud_log(libc.CString("unknown size vector in wield_type: chsize=%d, weapsize=%d"), chsize, weap.Size)
		return WIELD_NONE
	}
}

var race_bodyparts [24][23]int = [24][23]int{{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}

func racial_body_parts(ch *char_data) {
	var i int
	for i = 1; i < NUM_WEARS; i++ {
		if race_bodyparts[ch.Race][i] != 0 {
			SET_BIT_AR(ch.Bodyparts[:], uint32(int32(i)))
		} else {
			if BODY_FLAGGED(ch, uint32(int32(i))) {
				REMOVE_BIT_AR(ch.Bodyparts[:], uint32(int32(i)))
			}
		}
	}
}
