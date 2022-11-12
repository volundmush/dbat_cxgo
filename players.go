package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const LOAD_HIT = 0
const LOAD_MANA = 1
const LOAD_MOVE = 2
const LOAD_KI = 3
const LOAD_LIFE = 4
const ASCII_SAVE_POOFS = 0
const NUM_OF_SAVE_THROWS = 3

var player_table *player_index_element = nil
var top_of_p_table int = 0
var top_of_p_file int = 0
var top_idnum int = 0

func save_imc_pfile(ch *char_data) {
	var (
		fl       *C.FILE
		filename [4096]byte
	)
	if ch.Player_specials.Imcchardata == nil {
		return
	}
	if get_filename(&filename[0], uint64(4096), IMC_FILE, ch.Name) == 0 {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&filename[0]), "w")))
		return fl
	}()) == nil {
		if (*__errno_location()) != ENOENT {
			basic_mud_log(libc.CString("SYSERR: opening IMC2 file '%s' for writing: %s"), &filename[0], C.strerror(*__errno_location()))
		}
		return
	}
	imc_savechar(ch, fl)
	C.fclose(fl)
}
func load_imc_pfile(ch *char_data) {
	var (
		fl       *C.FILE
		filename [4096]byte
		word     *byte
	)
	if get_filename(&filename[0], uint64(4096), IMC_FILE, ch.Name) == 0 {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&filename[0]), "r")))
		return fl
	}()) == nil {
		if (*__errno_location()) != ENOENT {
			basic_mud_log(libc.CString("SYSERR: opening IMC2 file '%s' for reading: %s"), &filename[0], C.strerror(*__errno_location()))
		}
		return
	}
	for {
		word = imcfread_word(fl)
		if *word != 'I' {
			break
		}
		imc_loadchar(ch, fl, word)
	}
	C.fclose(fl)
}
func build_player_index() {
	var (
		rec_count  int = 0
		i          int
		plr_index  *C.FILE
		index_name [40]byte
		line       [256]byte
		bits       [64]byte
		arg2       [80]byte
	)
	stdio.Sprintf(&index_name[0], "%s%s", LIB_PLRC.FILES, INDEX_FILE)
	if (func() *C.FILE {
		plr_index = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&index_name[0]), "r")))
		return plr_index
	}()) == nil {
		top_of_p_table = -1
		basic_mud_log(libc.CString("No player index file!  First new char will be IMP!"))
		return
	}
	for get_line(plr_index, &line[0]) != 0 {
		if line[0] != '~' {
			rec_count++
		}
	}
	C.rewind(plr_index)
	if rec_count == 0 {
		player_table = nil
		top_of_p_table = -1
		return
	}
	player_table = &make([]player_index_element, rec_count)[0]
	for i = 0; i < rec_count; i++ {
		get_line(plr_index, &line[0])
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Admlevel = ADMLVL_NONE
		__isoc99_sscanf(&line[0], libc.CString("%ld %s %d %s %ld %d %d %d %ld"), &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id, &arg2[0], &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level, &bits[0], &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Last, &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Admlevel, &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Ship, &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Shiproom, &(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Played)
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name = (*byte)(unsafe.Pointer(&make([]int8, int(C.strlen(&arg2[0])+1))[0]))
		C.strcpy((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name, &arg2[0])
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Flags = int(asciiflag_conv(&bits[0]))
		top_idnum = MAX(top_idnum, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id)
	}
	C.fclose(plr_index)
	top_of_p_file = func() int {
		top_of_p_table = i - 1
		return top_of_p_table
	}()
}
func create_entry(name *byte) int {
	var (
		i   int
		pos int
	)
	if top_of_p_table == -1 {
		player_table = new(player_index_element)
		pos = func() int {
			top_of_p_table = 0
			return top_of_p_table
		}()
	} else if (func() int {
		pos = get_ptable_by_name(name)
		return pos
	}()) == -1 {
		i = func() int {
			p := &top_of_p_table
			*p++
			return *p
		}() + 1
		player_table = (*player_index_element)(libc.Realloc(unsafe.Pointer(player_table), i*int(unsafe.Sizeof(player_index_element{}))))
		pos = top_of_p_table
	}
	(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pos)))).Name = (*byte)(unsafe.Pointer(&make([]int8, int(C.strlen(name)+1))[0]))
	for i = 0; (func() byte {
		p := (*byte)(unsafe.Add(unsafe.Pointer((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pos)))).Name), i))
		*(*byte)(unsafe.Add(unsafe.Pointer((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pos)))).Name), i)) = byte(int8(C.tolower(int(*(*byte)(unsafe.Add(unsafe.Pointer(name), i))))))
		return *p
	}()) != 0; i++ {
	}
	(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pos)))).Flags = 0
	return pos
}
func save_player_index() {
	var (
		i          int
		index_name [50]byte
		bits       [64]byte
		index_file *C.FILE
	)
	stdio.Sprintf(&index_name[0], "%s%s", LIB_PLRC.FILES, INDEX_FILE)
	if (func() *C.FILE {
		index_file = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&index_name[0]), "w")))
		return index_file
	}()) == nil {
		basic_mud_log(libc.CString("SYSERR: Could not write player index file"))
		return
	}
	for i = 0; i <= top_of_p_table; i++ {
		if *(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name != 0 {
			sprintascii(&bits[0], bitvector_t((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Flags))
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(index_file)), "%ld %s %d %s %ld %d %d %d %ld\n", (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level, &func() [64]byte {
				if bits[0] != 0 {
					return bits
				}
				return func() [64]byte {
					var t [64]byte
					copy(t[:], []byte("0"))
					return t
				}()
			}()[0], (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Last, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Admlevel, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Ship, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Shiproom, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Played)
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(index_file)), "~\n")
	C.fclose(index_file)
}
func free_player_index() {
	var tp int
	if player_table == nil {
		return
	}
	for tp = 0; tp <= top_of_p_table; tp++ {
		if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(tp)))).Name != nil {
			libc.Free(unsafe.Pointer((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(tp)))).Name))
		}
	}
	libc.Free(unsafe.Pointer(player_table))
	player_table = nil
	top_of_p_table = 0
}
func get_ptable_by_name(name *byte) int {
	var i int
	for i = 0; i <= top_of_p_table; i++ {
		if C.strcasecmp((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name, name) == 0 {
			return i
		}
	}
	return -1
}
func get_id_by_name(name *byte) int {
	var i int
	for i = 0; i <= top_of_p_table; i++ {
		if C.strcasecmp((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name, name) == 0 {
			return (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id
		}
	}
	return -1
}
func get_name_by_id(id int) *byte {
	var i int
	for i = 0; i <= top_of_p_table; i++ {
		if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Id == id {
			return (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name
		}
	}
	return nil
}
func load_follower_from_file(fl *C.FILE, ch *char_data) {
	var (
		nr    int
		line  [2049]byte
		newch *char_data
	)
	if get_line(fl, &line[0]) == 0 {
		return
	}
	if line[0] != '#' || line[1] == 0 {
		return
	}
	nr = libc.Atoi(libc.GoString(&line[1]))
	newch = create_char()
	newch.Nr = real_mobile(mob_vnum(nr))
	if parse_mobile_from_file(fl, newch) == 0 {
		libc.Free(unsafe.Pointer(newch))
	} else {
		add_follower(newch, ch)
		newch.Master_id = ch.Idnum
		newch.Position = POS_STANDING
	}
}
func load_char(name *byte, ch *char_data) int {
	var (
		id    int
		i     int
		num   int = 0
		num2  int = 0
		num3  int = 0
		fl    *C.FILE
		fname [256]byte
		buf   [128]byte
		buf2  [128]byte
		line  [2049]byte
		tag   [6]byte
		f1    [128]byte
		f2    [128]byte
		f3    [128]byte
		f4    [128]byte
	)
	if (func() int {
		id = get_ptable_by_name(name)
		return id
	}()) < 0 {
		return -1
	} else {
		if get_filename(&fname[0], uint64(256), PLR_FILE, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Name) == 0 {
			return -1
		}
		if (func() *C.FILE {
			fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fname[0]), "r")))
			return fl
		}()) == nil {
			mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: Couldn't open player file %s"), &fname[0])
			return -1
		}
		ch.Affected = nil
		ch.Affectedv = nil
		for i = 1; i <= SKILL_TABLE_SIZE; i++ {
			for {
				ch.Skills[i] = 0
				if true {
					break
				}
			}
			for {
				ch.Skillmods[i] = 0
				if true {
					break
				}
			}
			for {
				ch.Skillperfs[i] = 0
				if true {
					break
				}
			}
		}
		ch.Sex = PFDEF_SEX
		ch.Size = -1
		ch.Chclass = PFDEF_CLASS
		for i = 0; i < NUM_CLASSES; i++ {
			ch.Chclasses[i] = 0
			ch.Epicclasses[i] = 0
		}
		ch.Loguser = C.strdup(libc.CString("NOUSER"))
		ch.Race = PFDEF_RACE
		ch.Admlevel = PFDEF_LEVEL
		ch.Level = PFDEF_LEVEL
		ch.Race_level = PFDEF_LEVEL
		ch.Rbank = PFDEF_SKIN
		ch.Rp = PFDEF_SKIN
		ch.Trp = PFDEF_SKIN
		ch.Suppression = PFDEF_SKIN
		ch.Suppressed = PFDEF_SKIN
		ch.Fury = PFDEF_HAIRL
		ch.Clan = C.strdup(libc.CString("None."))
		ch.Level_adj = PFDEF_LEVEL
		ch.Hometown = PFDEF_HOMETOWN
		ch.Height = PFDEF_HEIGHT
		ch.Weight = PFDEF_WEIGHT
		ch.Basepl = PFDEF_BASEPL
		ch.Relax_count = PFDEF_EYE
		ch.Blesslvl = PFDEF_HEIGHT
		ch.Lifeforce = PFDEF_BASEPL
		ch.Lifeperc = PFDEF_WEIGHT
		ch.Stupidkiss = 0
		ch.Position = POS_STANDING
		ch.Majinizer = PFDEF_BASEPL
		ch.Player_specials.Gauntlet = PFDEF_GAUNTLET
		ch.Baseki = PFDEF_BASEKI
		ch.Basest = PFDEF_BASEST
		ch.Hairl = PFDEF_HAIRL
		ch.Hairc = PFDEF_HAIRC
		ch.Skin = PFDEF_SKIN
		ch.Eye = PFDEF_EYE
		ch.Hairs = PFDEF_HAIRS
		ch.Distfea = PFDEF_DISTFEA
		ch.Radar1 = PFDEF_RADAR1
		ch.Ship = PFDEF_SHIP
		ch.Lastpl = PFDEF_LPLAY
		ch.Boosts = PFDEF_DISTFEA
		ch.Majinize = PFDEF_DISTFEA
		ch.Lastint = PFDEF_LPLAY
		ch.Deathtime = PFDEF_LPLAY
		ch.Starphase = PFDEF_EYE
		ch.Mimic = PFDEF_EYE
		ch.Skill_slots = 0
		ch.Tail_growth = 0
		ch.Player_specials.Trainstr = PFDEF_EYE
		ch.Player_specials.Trainspd = PFDEF_EYE
		ch.Player_specials.Trainwis = PFDEF_EYE
		ch.Player_specials.Trainagl = PFDEF_EYE
		ch.Player_specials.Traincon = PFDEF_EYE
		ch.Player_specials.Trainint = PFDEF_EYE
		ch.Rewtime = PFDEF_LPLAY
		ch.Dcount = PFDEF_EYE
		ch.Genome[0] = PFDEF_EYE
		ch.Preference = PFDEF_EYE
		ch.Genome[1] = PFDEF_EYE
		ch.Aura = PFDEF_SKIN
		for i = 0; i < 52; i++ {
			ch.Bonuses[i] = PFDEF_BOARD
		}
		ch.Combatexpertise = 0
		ch.Powerattack = 0
		ch.Limb_condition[0] = 0
		ch.Limb_condition[0] = 0
		ch.Limb_condition[1] = 0
		ch.Limb_condition[2] = 0
		ch.Lboard[0] = PFDEF_BOARD
		ch.Lboard[1] = PFDEF_BOARD
		ch.Lboard[2] = PFDEF_BOARD
		ch.Lboard[3] = PFDEF_BOARD
		ch.Lboard[4] = PFDEF_BOARD
		ch.Shipr = PFDEF_SHIPROOM
		ch.Radar2 = PFDEF_RADAR2
		ch.Radar3 = PFDEF_RADAR3
		ch.Droom = PFDEF_DROOM
		ch.Crank = PFDEF_CRANK
		ch.Alignment = PFDEF_ALIGNMENT
		ch.Alignment_ethic = PFDEF_ETHIC_ALIGNMENT
		for i = 0; i < AF_ARRAY_MAX; i++ {
			ch.Affected_by[i] = PFDEF_AFFFLAGS
		}
		for i = 0; i < PM_ARRAY_MAX; i++ {
			ch.Act[i] = PFDEF_PLRFLAGS
		}
		for i = 0; i < PR_ARRAY_MAX; i++ {
			ch.Player_specials.Pref[i] = PFDEF_PREFFLAGS
		}
		for i = 0; i < AD_ARRAY_MAX; i++ {
			ch.Admflags[i] = 0
		}
		for i = 0; i < NUM_OF_SAVE_THROWS; i++ {
			ch.Apply_saving_throw[i] = PFDEF_SAVETHROW
			ch.Saving_throw[i] = PFDEF_SAVETHROW
		}
		ch.Player_specials.Load_room = PFDEF_LOADROOM
		ch.Player_specials.Invis_level = PFDEF_INVISLEV
		ch.Player_specials.Freeze_level = PFDEF_FREEZELEV
		ch.Player_specials.Wimp_level = PFDEF_WIMPLEV
		ch.Powerattack = PFDEF_POWERATT
		ch.Player_specials.Conditions[HUNGER] = PFDEF_HUNGER
		ch.Player_specials.Conditions[THIRST] = PFDEF_THIRST
		ch.Player_specials.Conditions[DRUNK] = PFDEF_DRUNK
		ch.Player_specials.Bad_pws = PFDEF_BADPWS
		ch.Player_specials.Skill_points = PFDEF_PRACTICES
		for i = 0; i < NUM_CLASSES; i++ {
			ch.Player_specials.Class_skill_points[i] = PFDEF_PRACTICES
		}
		ch.Gold = PFDEF_GOLD
		ch.Backstabcool = 0
		ch.Con_cooldown = 0
		ch.Con_sdcooldown = 0
		ch.Bank_gold = PFDEF_BANK
		ch.Absorbs = PFDEF_BANK
		ch.IngestLearned = PFDEF_BANK
		ch.Player_specials.Racial_pref = PFDEF_BANK
		ch.Upgrade = PFDEF_BANK
		ch.Forgeting = PFDEF_BANK
		ch.Forgetcount = PFDEF_BANK
		ch.Kaioken = PFDEF_BANK
		ch.Exp = PFDEF_EXP
		ch.Transclass = PFDEF_EXP
		for i = 0; i < 6; i++ {
			ch.Transcost[i] = FALSE
		}
		ch.Moltexp = PFDEF_EXP
		ch.Accuracy_mod = PFDEF_ACCURACY
		ch.Accuracy = PFDEF_ACCURACY
		ch.Damage_mod = PFDEF_DAMAGE
		ch.Armor = PFDEF_AC
		ch.Real_abils.Str = PFDEF_STR
		ch.Real_abils.Dex = PFDEF_DEX
		ch.Real_abils.Intel = PFDEF_INT
		ch.Real_abils.Wis = PFDEF_WIS
		ch.Real_abils.Con = PFDEF_CON
		ch.Real_abils.Cha = PFDEF_CHA
		ch.Hit = PFDEF_HIT
		ch.Max_hit = PFDEF_MAXHIT
		ch.Mana = PFDEF_MANA
		ch.Max_mana = PFDEF_MAXMANA
		ch.Move = PFDEF_MOVE
		ch.Max_move = PFDEF_MAXMOVE
		ch.Ki = PFDEF_KI
		ch.Max_ki = PFDEF_MAXKI
		ch.Player_specials.Speaking = PFDEF_SPEAKING
		ch.Player_specials.Olc_zone = -1
		ch.Player_specials.Host = nil
		for i = 1; i < (int(MAX_SPELL_LEVEL * 10)); i++ {
			ch.Player_specials.Spellmem[i] = 0
		}
		for i = 0; i < MAX_SPELL_LEVEL; i++ {
			ch.Player_specials.Spell_level[i] = 0
		}
		ch.Player_specials.Memcursor = 0
		ch.Time.Birth = func() int64 {
			p := &ch.Time.Created
			ch.Time.Created = func() int64 {
				p := &ch.Time.Maxage
				ch.Time.Maxage = 0
				return *p
			}()
			return *p
		}()
		ch.Followers = nil
		ch.Player_specials.Page_length = PFDEF_PAGELENGTH
		for i = 0; i < NUM_COLOR; i++ {
			ch.Player_specials.Color_choices[i] = nil
		}
		for get_line(fl, &line[0]) != 0 {
			tag_argument(&line[0], &tag[0])
			switch tag[0] {
			case 'A':
				if C.strcmp(&tag[0], libc.CString("Ac  ")) == 0 {
					ch.Armor = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Act ")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%s %s %s %s"), &f1[0], &f2[0], &f3[0], &f4[0])
					ch.Act[0] = asciiflag_conv(&f1[0])
					ch.Act[1] = asciiflag_conv(&f2[0])
					ch.Act[2] = asciiflag_conv(&f3[0])
					ch.Act[3] = asciiflag_conv(&f4[0])
				} else if C.strcmp(&tag[0], libc.CString("Aff ")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%s %s %s %s"), &f1[0], &f2[0], &f3[0], &f4[0])
					ch.Affected_by[0] = int(asciiflag_conv(&f1[0]))
					ch.Affected_by[1] = int(asciiflag_conv(&f2[0]))
					ch.Affected_by[2] = int(asciiflag_conv(&f3[0]))
					ch.Affected_by[3] = int(asciiflag_conv(&f4[0]))
				} else if C.strcmp(&tag[0], libc.CString("Affs")) == 0 {
					load_affects(fl, ch, 0)
				} else if C.strcmp(&tag[0], libc.CString("Affv")) == 0 {
					load_affects(fl, ch, 1)
				} else if C.strcmp(&tag[0], libc.CString("AdmL")) == 0 {
					ch.Admlevel = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Abso")) == 0 {
					ch.Absorbs = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("AdmF")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%s %s %s %s"), &f1[0], &f2[0], &f3[0], &f4[0])
					ch.Admflags[0] = asciiflag_conv(&f1[0])
					ch.Admflags[1] = asciiflag_conv(&f2[0])
					ch.Admflags[2] = asciiflag_conv(&f3[0])
					ch.Admflags[3] = asciiflag_conv(&f4[0])
				} else if C.strcmp(&tag[0], libc.CString("Alin")) == 0 {
					ch.Alignment = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Aura")) == 0 {
					ch.Aura = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'B':
				if C.strcmp(&tag[0], libc.CString("Badp")) == 0 {
					ch.Player_specials.Bad_pws = uint8(int8(libc.Atoi(libc.GoString(&line[0]))))
				} else if C.strcmp(&tag[0], libc.CString("Bank")) == 0 {
					ch.Bank_gold = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Bki ")) == 0 {
					load_BASE(ch, &line[0], LOAD_MANA)
				} else if C.strcmp(&tag[0], libc.CString("Blss")) == 0 {
					ch.Blesslvl = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Boam")) == 0 {
					ch.Lboard[0] = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Boai")) == 0 {
					ch.Lboard[1] = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Boac")) == 0 {
					ch.Lboard[2] = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Boad")) == 0 {
					ch.Lboard[3] = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Boab")) == 0 {
					ch.Lboard[4] = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Bonu")) == 0 {
					load_bonuses(fl, ch, FALSE != 0)
				} else if C.strcmp(&tag[0], libc.CString("Boos")) == 0 {
					ch.Boosts = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Bpl ")) == 0 {
					load_BASE(ch, &line[0], LOAD_HIT)
				} else if C.strcmp(&tag[0], libc.CString("Brth")) == 0 {
					ch.Time.Birth = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Bst ")) == 0 {
					load_BASE(ch, &line[0], LOAD_MOVE)
				}
			case 'C':
				if C.strcmp(&tag[0], libc.CString("Cha ")) == 0 {
					ch.Real_abils.Cha = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Clan")) == 0 {
					ch.Clan = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("Clar")) == 0 {
					ch.Crank = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Clas")) == 0 {
					ch.Chclass = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Colr")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%d %s"), &num, &buf2[0])
					ch.Player_specials.Color_choices[num] = C.strdup(&buf2[0])
				} else if C.strcmp(&tag[0], libc.CString("Con ")) == 0 {
					ch.Real_abils.Con = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Cool")) == 0 {
					ch.Con_cooldown = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Crtd")) == 0 {
					ch.Time.Created = int64(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'D':
				if C.strcmp(&tag[0], libc.CString("Deat")) == 0 {
					ch.Deathtime = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Deac")) == 0 {
					ch.Dcount = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Desc")) == 0 {
					ch.Description = fread_string(fl, &buf2[0])
				} else if C.strcmp(&tag[0], libc.CString("Dex ")) == 0 {
					ch.Real_abils.Dex = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Drnk")) == 0 {
					ch.Player_specials.Conditions[DRUNK] = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Damg")) == 0 {
					ch.Damage_mod = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Droo")) == 0 {
					ch.Droom = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'E':
				if C.strcmp(&tag[0], libc.CString("Exp ")) == 0 {
					ch.Exp = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Eali")) == 0 {
					ch.Alignment_ethic = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Ecls")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%d=%d"), &num, &num2)
					ch.Epicclasses[num] = num2
				} else if C.strcmp(&tag[0], libc.CString("Eye ")) == 0 {
					ch.Eye = int8(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'F':
				if C.strcmp(&tag[0], libc.CString("Fisd")) == 0 {
					ch.Accuracy_mod = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Frez")) == 0 {
					ch.Player_specials.Freeze_level = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Forc")) == 0 {
					ch.Forgetcount = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Forg")) == 0 {
					ch.Forgeting = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Fury")) == 0 {
					ch.Fury = int16(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'G':
				if C.strcmp(&tag[0], libc.CString("Gold")) == 0 {
					ch.Gold = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Gaun")) == 0 {
					ch.Player_specials.Gauntlet = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Geno")) == 0 {
					ch.Genome[0] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Gen1")) == 0 {
					ch.Genome[1] = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'H':
				if C.strcmp(&tag[0], libc.CString("Hit ")) == 0 {
					load_HMVS(ch, &line[0], LOAD_HIT)
				} else if C.strcmp(&tag[0], libc.CString("HitD")) == 0 {
					ch.Race_level = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Hite")) == 0 {
					ch.Height = uint8(int8(libc.Atoi(libc.GoString(&line[0]))))
				} else if C.strcmp(&tag[0], libc.CString("Home")) == 0 {
					ch.Hometown = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Host")) == 0 {
					if ch.Player_specials.Host != nil {
						libc.Free(unsafe.Pointer(ch.Player_specials.Host))
					}
					ch.Player_specials.Host = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("Hrc ")) == 0 {
					ch.Hairc = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Hrl ")) == 0 {
					ch.Hairl = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Hrs ")) == 0 {
					ch.Hairs = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Hung")) == 0 {
					ch.Player_specials.Conditions[HUNGER] = int8(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'I':
				if C.strcmp(&tag[0], libc.CString("Id  ")) == 0 {
					ch.Idnum = int32(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("INGl")) == 0 {
					ch.IngestLearned = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Int ")) == 0 {
					ch.Real_abils.Intel = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Invs")) == 0 {
					ch.Player_specials.Invis_level = int16(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'K':
				if C.strcmp(&tag[0], libc.CString("Ki  ")) == 0 {
					load_HMVS(ch, &line[0], LOAD_KI)
				} else if C.strcmp(&tag[0], libc.CString("Kaio")) == 0 {
					ch.Kaioken = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'L':
				if C.strcmp(&tag[0], libc.CString("Last")) == 0 {
					ch.Time.Logon = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Lern")) == 0 {
					ch.Player_specials.Class_skill_points[ch.Chclass] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Levl")) == 0 {
					ch.Level = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("LF  ")) == 0 {
					load_BASE(ch, &line[0], LOAD_LIFE)
				} else if C.strcmp(&tag[0], libc.CString("LFPC")) == 0 {
					ch.Lifeperc = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Lila")) == 0 {
					ch.Limb_condition[1] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Lill")) == 0 {
					ch.Limb_condition[3] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Lira")) == 0 {
					ch.Limb_condition[0] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Lirl")) == 0 {
					ch.Limb_condition[2] = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Lint")) == 0 {
					ch.Lastint = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Lpla")) == 0 {
					ch.Lastpl = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("LvlA")) == 0 {
					ch.Level_adj = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'M':
				if C.strcmp(&tag[0], libc.CString("Mana")) == 0 {
					load_HMVS(ch, &line[0], LOAD_MANA)
				} else if C.strcmp(&tag[0], libc.CString("Mexp")) == 0 {
					load_molt(ch, &line[0])
				} else if C.strcmp(&tag[0], libc.CString("Mlvl")) == 0 {
					ch.Moltlevel = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Move")) == 0 {
					load_HMVS(ch, &line[0], LOAD_MOVE)
				} else if C.strcmp(&tag[0], libc.CString("Mcls")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%d=%d"), &num, &num2)
					ch.Chclasses[num] = num2
				} else if C.strcmp(&tag[0], libc.CString("Maji")) == 0 {
					ch.Majinize = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Majm")) == 0 {
					load_majin(ch, &line[0])
				} else if C.strcmp(&tag[0], libc.CString("Mimi")) == 0 {
					ch.Mimic = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("MxAg")) == 0 {
					ch.Time.Maxage = int64(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'N':
				if C.strcmp(&tag[0], libc.CString("Name")) == 0 {
					ch.Name = C.strdup(&line[0])
				}
			case 'O':
				if C.strcmp(&tag[0], libc.CString("Olc ")) == 0 {
					ch.Player_specials.Olc_zone = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'P':
				if C.strcmp(&tag[0], libc.CString("Page")) == 0 {
					ch.Player_specials.Page_length = uint8(int8(libc.Atoi(libc.GoString(&line[0]))))
				} else if C.strcmp(&tag[0], libc.CString("Phas")) == 0 {
					ch.Distfea = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Phse")) == 0 {
					ch.Starphase = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Plyd")) == 0 {
					ch.Time.Played = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("PfIn")) == 0 {
					ch.Player_specials.Poofin = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("PfOt")) == 0 {
					ch.Player_specials.Poofout = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("Pole")) == 0 {
					ch.Accuracy = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Posi")) == 0 {
					ch.Position = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("PwrA")) == 0 {
					ch.Powerattack = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Pref")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%s %s %s %s"), &f1[0], &f2[0], &f3[0], &f4[0])
					ch.Player_specials.Pref[0] = asciiflag_conv(&f1[0])
					ch.Player_specials.Pref[1] = asciiflag_conv(&f2[0])
					ch.Player_specials.Pref[2] = asciiflag_conv(&f3[0])
					ch.Player_specials.Pref[3] = asciiflag_conv(&f4[0])
				} else if C.strcmp(&tag[0], libc.CString("Prff")) == 0 {
					ch.Preference = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'R':
				if C.strcmp(&tag[0], libc.CString("Race")) == 0 {
					ch.Race = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Raci")) == 0 {
					ch.Player_specials.Racial_pref = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("RBan")) == 0 {
					ch.Rbank = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("rDis")) == 0 {
					ch.Rdisplay = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("Rela")) == 0 {
					ch.Relax_count = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Rtim")) == 0 {
					ch.Rewtime = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Rad1")) == 0 {
					ch.Radar1 = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Rad2")) == 0 {
					ch.Radar2 = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Rad3")) == 0 {
					ch.Radar3 = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Room")) == 0 {
					ch.Player_specials.Load_room = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("RPfe")) == 0 {
					ch.Feature = C.strdup(&line[0])
				} else if C.strcmp(&tag[0], libc.CString("RPP ")) == 0 {
					ch.Rp = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'S':
				if C.strcmp(&tag[0], libc.CString("Sex ")) == 0 {
					ch.Sex = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Ship")) == 0 {
					ch.Ship = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Scoo")) == 0 {
					ch.Con_sdcooldown = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Shpr")) == 0 {
					ch.Shipr = room_vnum(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Skil")) == 0 {
					load_skills(fl, ch, FALSE != 0)
				} else if C.strcmp(&tag[0], libc.CString("Skn ")) == 0 {
					ch.Skin = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Size")) == 0 {
					ch.Size = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("SklB")) == 0 {
					load_skills(fl, ch, TRUE != 0)
				} else if C.strcmp(&tag[0], libc.CString("SkRc")) == 0 {
					ch.Player_specials.Skill_points = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("SkCl")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%d %d"), &num2, &num3)
					ch.Player_specials.Class_skill_points[num2] = num3
				} else if C.strcmp(&tag[0], libc.CString("Slot")) == 0 {
					ch.Skill_slots = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Spek")) == 0 {
					ch.Player_specials.Speaking = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Str ")) == 0 {
					ch.Real_abils.Str = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Stuk")) == 0 {
					ch.Stupidkiss = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Supp")) == 0 {
					ch.Suppression = int64(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Sups")) == 0 {
					ch.Suppressed = int64(libc.Atoi(libc.GoString(&line[0])))
				}
			case 'T':
				if C.strcmp(&tag[0], libc.CString("Tgro")) == 0 {
					ch.Tail_growth = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Tcla")) == 0 {
					ch.Transclass = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Tcos")) == 0 {
					__isoc99_sscanf(&line[0], libc.CString("%d %d"), &num2, &num3)
					ch.Transcost[num2] = num3
				} else if C.strcmp(&tag[0], libc.CString("Thir")) == 0 {
					ch.Player_specials.Conditions[THIRST] = int8(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Thr1")) == 0 {
					ch.Apply_saving_throw[0] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Thr2")) == 0 {
					ch.Apply_saving_throw[1] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Thr3")) == 0 {
					ch.Apply_saving_throw[2] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Thr4")) == 0 || C.strcmp(&tag[0], libc.CString("Thr5")) == 0 {
				} else if C.strcmp(&tag[0], libc.CString("ThB1")) == 0 {
					ch.Saving_throw[0] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("ThB2")) == 0 {
					ch.Saving_throw[1] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("ThB3")) == 0 {
					ch.Saving_throw[2] = int16(libc.Atoi(libc.GoString(&line[0])))
				} else if C.strcmp(&tag[0], libc.CString("Trns")) == 0 {
					ch.Player_specials.Ability_trains = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trag")) == 0 {
					ch.Player_specials.Trainagl = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trco")) == 0 {
					ch.Player_specials.Traincon = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trin")) == 0 {
					ch.Player_specials.Trainint = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trsp")) == 0 {
					ch.Player_specials.Trainspd = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trst")) == 0 {
					ch.Player_specials.Trainstr = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trwi")) == 0 {
					ch.Player_specials.Trainwis = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Trp ")) == 0 {
					ch.Trp = libc.Atoi(libc.GoString(&line[0]))
				}
			case 'U':
				if C.strcmp(&tag[0], libc.CString("Upgr")) == 0 {
					ch.Upgrade = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("User")) == 0 {
					if ch.Loguser != nil {
						libc.Free(unsafe.Pointer(ch.Loguser))
					}
					ch.Loguser = C.strdup(&line[0])
				}
			case 'V':
				if C.strcmp(&tag[0], libc.CString("Voic")) == 0 {
					ch.Voice = C.strdup(&line[0])
				}
			case 'W':
				if C.strcmp(&tag[0], libc.CString("Wate")) == 0 {
					ch.Weight = uint8(int8(libc.Atoi(libc.GoString(&line[0]))))
				} else if C.strcmp(&tag[0], libc.CString("Wimp")) == 0 {
					ch.Player_specials.Wimp_level = libc.Atoi(libc.GoString(&line[0]))
				} else if C.strcmp(&tag[0], libc.CString("Wis ")) == 0 {
					ch.Real_abils.Wis = int8(libc.Atoi(libc.GoString(&line[0])))
				}
			default:
				stdio.Sprintf(&buf[0], "SYSERR: Unknown tag %s in pfile %s", &tag[0], name)
			}
		}
	}
	if ch.Time.Created == 0 {
		basic_mud_log(libc.CString("No creation timestamp for user %s, using current time"), GET_NAME(ch))
		ch.Time.Created = C.time(nil)
	}
	if ch.Time.Birth == 0 {
		basic_mud_log(libc.CString("No birthday for user %s, using standard starting age determination"), GET_NAME(ch))
		ch.Time.Birth = C.time(nil) - birth_age(ch)
	}
	if ch.Time.Maxage == 0 {
		basic_mud_log(libc.CString("No max age for user %s, using standard max age determination"), GET_NAME(ch))
		ch.Time.Maxage = ch.Time.Birth + max_age(ch)
	}
	affect_total(ch)
	if ch.Admlevel >= ADMLVL_IMMORT {
		for i = 1; i <= SKILL_TABLE_SIZE; i++ {
			for {
				ch.Skills[i] = 100
				if true {
					break
				}
			}
		}
		ch.Player_specials.Conditions[HUNGER] = -1
		ch.Player_specials.Conditions[THIRST] = -1
		ch.Player_specials.Conditions[DRUNK] = -1
	}
	if config_info.Operation.Imc_enabled != 0 {
		imc_initchar(ch)
		load_imc_pfile(ch)
	}
	if ch.Race == RACE_ANDROID {
		ch.Player_specials.Conditions[HUNGER] = -1
		ch.Player_specials.Conditions[THIRST] = -1
		ch.Player_specials.Conditions[DRUNK] = -1
	}
	C.fclose(fl)
	return id
}
func kill_ems(str *byte) {
	var (
		ptr1 *byte
		ptr2 *byte
		tmp  *byte
	)
	_ = tmp
	tmp = str
	ptr1 = str
	ptr2 = str
	for *ptr1 != 0 {
		if (func() byte {
			p := (func() *byte {
				p := &ptr2
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			*(func() *byte {
				p := &ptr2
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()) = *(func() *byte {
				p := &ptr1
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}())
			return *p
		}()) == '\r' {
			if *ptr1 == '\r' {
				ptr1 = (*byte)(unsafe.Add(unsafe.Pointer(ptr1), 1))
			}
		}
	}
	*ptr2 = '\x00'
}
func save_char_pets(ch *char_data) {
	var (
		foll  *follow_type
		fname [40]byte
		fl    *C.FILE
	)
	if IS_NPC(ch) || ch.Pfilepos < 0 {
		return
	}
	if get_filename(&fname[0], uint64(40), PET_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fname[0]), "w")))
		return fl
	}()) == nil {
		mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: Couldn't open pet file %s for write"), &fname[0])
		return
	}
	for foll = ch.Followers; foll != nil; foll = foll.Next {
		write_mobile_record(GET_MOB_VNUM(foll.Follower), foll.Follower, fl)
	}
	C.fclose(fl)
}
func load_char_pets(ch *char_data) {
	var (
		fname     [40]byte
		fl        *C.FILE
		load_room int64
		foll      *follow_type
	)
	if IS_NPC(ch) || ch.Pfilepos < 0 {
		return
	}
	if get_filename(&fname[0], uint64(40), PET_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fname[0]), "r")))
		return fl
	}()) == nil {
		return
	}
	for C.feof(fl) == 0 {
		load_follower_from_file(fl, ch)
	}
	for foll = ch.Followers; foll != nil; foll = foll.Next {
		load_room = int64(real_room(1))
		if load_room == int64(-1) {
			load_room = int64(real_room(room_vnum(ch.In_room)))
		}
		char_to_room(foll.Follower, room_rnum(load_room))
		act(libc.CString("You are joined by $N."), FALSE, ch, nil, unsafe.Pointer(foll.Follower), TO_CHAR)
	}
}
func save_char(ch *char_data) {
	var (
		fl         *C.FILE
		fname      [40]byte
		buf        [64936]byte
		i          int
		id         int
		save_index int = FALSE
		aff        *affected_type
		tmp_aff    [32]affected_type
		tmp_affv   [32]affected_type
		char_eq    [23]*obj_data
		fbuf1      [64936]byte
		fbuf2      [64936]byte
		fbuf3      [64936]byte
		fbuf4      [64936]byte
	)
	if IS_NPC(ch) || ch.Pfilepos < 0 {
		return
	}
	if ch.Desc != nil {
		if ch.Desc.Host != nil && ch.Desc.Host[0] != 0 {
			if ch.Player_specials.Host == nil {
				ch.Player_specials.Host = C.strdup(&ch.Desc.Host[0])
			} else if ch.Player_specials.Host != nil && C.strcmp(ch.Player_specials.Host, &ch.Desc.Host[0]) == 0 {
				libc.Free(unsafe.Pointer(ch.Player_specials.Host))
				ch.Player_specials.Host = C.strdup(&ch.Desc.Host[0])
			}
		}
		if ch.Desc.Connected == CON_PLAYING {
			ch.Time.Played += C.time(nil) - ch.Time.Logon
			ch.Time.Logon = C.time(nil)
		}
	}
	if get_filename(&fname[0], uint64(40), PLR_FILE, GET_NAME(ch)) == 0 {
		return
	}
	if (func() *C.FILE {
		fl = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&fname[0]), "w")))
		return fl
	}()) == nil {
		mudlog(NRM, ADMLVL_GOD, TRUE, libc.CString("SYSERR: Couldn't open player file %s for write"), &fname[0])
		return
	}
	if ch.Trp < ch.Rp {
		ch.Trp = ch.Rp
	}
	if ch.Desc != nil && ch.Desc.User != nil {
		userWrite(ch.Desc, 0, 0, 0, libc.CString("index"))
	}
	for i = 0; i < NUM_WEARS; i++ {
		if (ch.Equipment[i]) != nil {
			char_eq[i] = unequip_char(ch, i)
		} else {
			char_eq[i] = nil
		}
	}
	for func() int {
		aff = ch.Affected
		return func() int {
			i = 0
			return i
		}()
	}(); i < MAX_AFFECT; i++ {
		if aff != nil {
			tmp_aff[i] = *aff
			tmp_aff[i].Next = nil
			aff = aff.Next
		} else {
			tmp_aff[i].Type = 0
			tmp_aff[i].Duration = 0
			tmp_aff[i].Modifier = 0
			tmp_aff[i].Specific = 0
			tmp_aff[i].Location = 0
			tmp_aff[i].Bitvector = 0
			tmp_aff[i].Next = nil
		}
	}
	for func() int {
		aff = ch.Affectedv
		return func() int {
			i = 0
			return i
		}()
	}(); i < MAX_AFFECT; i++ {
		if aff != nil {
			tmp_affv[i] = *aff
			tmp_affv[i].Next = nil
			aff = aff.Next
		} else {
			tmp_affv[i].Type = 0
			tmp_affv[i].Duration = 0
			tmp_affv[i].Modifier = 0
			tmp_affv[i].Location = 0
			tmp_affv[i].Specific = 0
			tmp_affv[i].Bitvector = 0
			tmp_affv[i].Next = nil
		}
	}
	save_char_vars(ch)
	for ch.Affected != nil {
		affect_remove(ch, ch.Affected)
	}
	for ch.Affectedv != nil {
		affectv_remove(ch, ch.Affectedv)
	}
	if i >= MAX_AFFECT && aff != nil && aff.Next != nil {
		basic_mud_log(libc.CString("SYSERR: WARNING: OUT OF STORE ROOM FOR AFFECTED TYPES!!!"))
	}
	ch.Aff_abils = ch.Real_abils
	if GET_NAME(ch) != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Name: %s\n", GET_NAME(ch))
	}
	if GET_USER(ch) != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "User: %s\n", GET_USER(ch))
	}
	if ch.Voice != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Voic: %s\n", ch.Voice)
	}
	if ch.Clan != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Clan: %s\n", ch.Clan)
	}
	if ch.Feature != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "RPfe: %s\n", ch.Feature)
	}
	if ch.Player_specials.Ability_trains != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trns: %d\n", ch.Player_specials.Ability_trains)
	}
	if ch.Description != nil && *ch.Description != 0 {
		C.strcpy(&buf[0], ch.Description)
		kill_ems(&buf[0])
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Desc:\n%s~\n", &buf[0])
	}
	if ch.Player_specials.Poofin != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "PfIn: %s\n", ch.Player_specials.Poofin)
	}
	if ch.Player_specials.Poofout != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "PfOt: %s\n", ch.Player_specials.Poofout)
	}
	if ch.Sex != PFDEF_SEX {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Sex : %d\n", ch.Sex)
	}
	if ch.Size != -1 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Size: %d\n", ch.Size)
	}
	if ch.Chclass != PFDEF_CLASS {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Clas: %d\n", ch.Chclass)
	}
	if ch.Race != PFDEF_RACE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Race: %d\n", ch.Race)
	}
	if ch.Player_specials.Racial_pref != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Raci: %d\n", ch.Player_specials.Racial_pref)
	}
	if ch.Admlevel != PFDEF_LEVEL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "AdmL: %d\n", ch.Admlevel)
	}
	if ch.Level != PFDEF_LEVEL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Levl: %d\n", ch.Level)
	}
	if ch.Race_level != PFDEF_LEVEL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "HitD: %d\n", ch.Race_level)
	}
	if ch.Level_adj != PFDEF_LEVEL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "LvlA: %d\n", ch.Level_adj)
	}
	if ch.Hometown != PFDEF_HOMETOWN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Home: %d\n", ch.Hometown)
	}
	for i = 0; i < NUM_CLASSES; i++ {
		if ((ch.Chclasses[i]) + (ch.Epicclasses[i])) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Mcls: %d=%d\n", i, ch.Chclasses[i])
		}
		if (ch.Epicclasses[i]) != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Ecls: %d=%d\n", i, ch.Epicclasses[i])
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Id  : %d\n", ch.Idnum)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Brth: %ld\n", ch.Time.Birth)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Crtd: %ld\n", ch.Time.Created)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "MxAg: %ld\n", ch.Time.Maxage)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Plyd: %ld\n", ch.Time.Played)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Last: %ld\n", ch.Time.Logon)
	if ch.Player_specials.Host != nil {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Host: %s\n", ch.Player_specials.Host)
	}
	if ch.Height != PFDEF_HEIGHT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hite: %d\n", ch.Height)
	}
	if ch.Weight != PFDEF_HEIGHT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Wate: %d\n", ch.Weight)
	}
	if ch.Alignment != PFDEF_ALIGNMENT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Alin: %d\n", ch.Alignment)
	}
	if ch.Aura != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Aura: %d\n", ch.Aura)
	}
	if ch.Alignment_ethic != PFDEF_ETHIC_ALIGNMENT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Eali: %d\n", ch.Alignment_ethic)
	}
	sprintascii(&fbuf1[0], ch.Act[0])
	sprintascii(&fbuf2[0], ch.Act[1])
	sprintascii(&fbuf3[0], ch.Act[2])
	sprintascii(&fbuf4[0], ch.Act[3])
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Act : %s %s %s %s\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0])
	sprintascii(&fbuf1[0], bitvector_t(ch.Affected_by[0]))
	sprintascii(&fbuf2[0], bitvector_t(ch.Affected_by[1]))
	sprintascii(&fbuf3[0], bitvector_t(ch.Affected_by[2]))
	sprintascii(&fbuf4[0], bitvector_t(ch.Affected_by[3]))
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Aff : %s %s %s %s\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0])
	sprintascii(&fbuf1[0], ch.Player_specials.Pref[0])
	sprintascii(&fbuf2[0], ch.Player_specials.Pref[1])
	sprintascii(&fbuf3[0], ch.Player_specials.Pref[2])
	sprintascii(&fbuf4[0], ch.Player_specials.Pref[3])
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Pref: %s %s %s %s\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0])
	sprintascii(&fbuf1[0], ch.Admflags[0])
	sprintascii(&fbuf2[0], ch.Admflags[1])
	sprintascii(&fbuf3[0], ch.Admflags[2])
	sprintascii(&fbuf4[0], ch.Admflags[3])
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "AdmF: %s %s %s %s\n", &fbuf1[0], &fbuf2[0], &fbuf3[0], &fbuf4[0])
	if (ch.Saving_throw[0]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "ThB1: %d\n", ch.Saving_throw[0])
	}
	if (ch.Saving_throw[1]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "ThB2: %d\n", ch.Saving_throw[1])
	}
	if (ch.Saving_throw[2]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "ThB3: %d\n", ch.Saving_throw[2])
	}
	if (ch.Apply_saving_throw[0]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Thr1: %d\n", ch.Apply_saving_throw[0])
	}
	if (ch.Apply_saving_throw[1]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Thr2: %d\n", ch.Apply_saving_throw[1])
	}
	if (ch.Apply_saving_throw[2]) != PFDEF_SAVETHROW {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Thr3: %d\n", ch.Apply_saving_throw[2])
	}
	if ch.Player_specials.Wimp_level != PFDEF_WIMPLEV {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Wimp: %d\n", ch.Player_specials.Wimp_level)
	}
	if ch.Powerattack != PFDEF_POWERATT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "PwrA: %d\n", ch.Powerattack)
	}
	if ch.Player_specials.Freeze_level != PFDEF_FREEZELEV {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Frez: %d\n", ch.Player_specials.Freeze_level)
	}
	if ch.Player_specials.Invis_level != PFDEF_INVISLEV {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Invs: %d\n", ch.Player_specials.Invis_level)
	}
	if ch.Player_specials.Load_room != PFDEF_LOADROOM {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Room: %d\n", ch.Player_specials.Load_room)
	}
	if ch.Player_specials.Bad_pws != PFDEF_BADPWS {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Badp: %d\n", ch.Player_specials.Bad_pws)
	}
	if ch.Player_specials.Skill_points != PFDEF_PRACTICES {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "SkRc: %d\n", ch.Player_specials.Skill_points)
	}
	for i = 0; i < 6; i++ {
		if (ch.Transcost[i]) != FALSE {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Tcos: %d %d\n", i, ch.Transcost[i])
		}
	}
	for i = 0; i < NUM_CLASSES; i++ {
		if (ch.Player_specials.Class_skill_points[i]) != PFDEF_PRACTICES {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "SkCl: %d %d\n", i, ch.Player_specials.Class_skill_points[i])
		}
	}
	if (ch.Player_specials.Conditions[HUNGER]) != PFDEF_HUNGER && ch.Admlevel < ADMLVL_IMMORT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hung: %d\n", ch.Player_specials.Conditions[HUNGER])
	}
	if (ch.Player_specials.Conditions[THIRST]) != PFDEF_THIRST && ch.Admlevel < ADMLVL_IMMORT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Thir: %d\n", ch.Player_specials.Conditions[THIRST])
	}
	if (ch.Player_specials.Conditions[DRUNK]) != PFDEF_DRUNK && ch.Admlevel < ADMLVL_IMMORT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Drnk: %d\n", ch.Player_specials.Conditions[DRUNK])
	}
	if ch.Hit != PFDEF_HIT || ch.Max_hit != PFDEF_MAXHIT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hit : %lld/%lld\n", ch.Hit, ch.Max_hit)
	}
	if ch.Mana != PFDEF_MANA || ch.Max_mana != PFDEF_MAXMANA {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Mana: %lld/%lld\n", ch.Mana, ch.Max_mana)
	}
	if ch.Move != PFDEF_MOVE || ch.Max_move != PFDEF_MAXMOVE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Move: %lld/%lld\n", ch.Move, ch.Max_move)
	}
	if ch.Ki != PFDEF_KI || ch.Max_ki != PFDEF_MAXKI {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Ki  : %lld/%lld\n", ch.Ki, ch.Max_ki)
	}
	if ch.Aff_abils.Str != PFDEF_STR {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Str : %d\n", ch.Aff_abils.Str)
	}
	if ch.Aff_abils.Intel != PFDEF_INT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Int : %d\n", ch.Aff_abils.Intel)
	}
	if ch.Aff_abils.Wis != PFDEF_WIS {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Wis : %d\n", ch.Aff_abils.Wis)
	}
	if ch.Aff_abils.Dex != PFDEF_DEX {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Dex : %d\n", ch.Aff_abils.Dex)
	}
	if ch.Aff_abils.Con != PFDEF_CON {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Con : %d\n", ch.Aff_abils.Con)
	}
	if ch.Aff_abils.Cha != PFDEF_CHA {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Cha : %d\n", ch.Aff_abils.Cha)
	}
	if ch.Con_cooldown != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Cool: %d\n", ch.Con_cooldown)
	}
	if ch.Con_cooldown != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Scoo: %d\n", ch.Con_sdcooldown)
	}
	if ch.Armor != PFDEF_AC {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Ac  : %d\n", ch.Armor)
	}
	if ch.Absorbs != PFDEF_GOLD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Abso: %d\n", ch.Absorbs)
	}
	if ch.IngestLearned != PFDEF_GOLD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "INGl: %d\n", ch.IngestLearned)
	}
	if ch.Upgrade != PFDEF_GOLD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Upgr: %d\n", ch.Upgrade)
	}
	if ch.Forgeting != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Forg: %d\n", ch.Forgeting)
	}
	if ch.Forgetcount != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Forc: %d\n", ch.Forgetcount)
	}
	if ch.Kaioken != PFDEF_GOLD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Kaio: %d\n", ch.Kaioken)
	}
	if ch.Gold != PFDEF_GOLD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Gold: %d\n", ch.Gold)
	}
	if ch.Bank_gold != PFDEF_BANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Bank: %d\n", ch.Bank_gold)
	}
	if ch.Exp != PFDEF_EXP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Exp : %lld\n", ch.Exp)
	}
	if ch.Transclass != PFDEF_EXP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Tcla: %d\n", ch.Transclass)
	}
	if ch.Moltexp != PFDEF_EXP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Mexp: %lld\n", ch.Moltexp)
	}
	if ch.Majinizer != PFDEF_EXP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Majm: %lld\n", ch.Majinizer)
	}
	if ch.Moltlevel != PFDEF_EXP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Mlvl: %d\n", ch.Moltlevel)
	}
	if ch.Accuracy_mod != PFDEF_ACCURACY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Fisd: %d\n", ch.Accuracy_mod)
	}
	if ch.Accuracy != PFDEF_ACCURACY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Pole: %d\n", ch.Accuracy)
	}
	if ch.Preference != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Prff: %d\n", ch.Preference)
	}
	if ch.Damage_mod != PFDEF_DAMAGE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Damg: %d\n", ch.Damage_mod)
	}
	if ch.Player_specials.Speaking != PFDEF_SPEAKING {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Spek: %d\n", ch.Player_specials.Speaking)
	}
	if ch.Player_specials.Olc_zone != int(-1) {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Olc : %d\n", ch.Player_specials.Olc_zone)
	}
	if ch.Player_specials.Page_length != PFDEF_PAGELENGTH {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Page: %d\n", ch.Player_specials.Page_length)
	}
	if ch.Player_specials.Gauntlet != PFDEF_GAUNTLET {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Gaun: %d\n", ch.Player_specials.Gauntlet)
	}
	if (ch.Genome[0]) != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Geno: %d\n", ch.Genome[0])
	}
	if (ch.Genome[1]) != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Gen1: %d\n", ch.Genome[1])
	}
	if ch.Position != POS_STANDING {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Posi: %d\n", ch.Position)
	}
	if ch.Lifeforce != PFDEF_BASEPL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "LF  : %lld\n", ch.Lifeforce)
	}
	if ch.Lifeperc != PFDEF_WEIGHT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "LFPC: %d\n", ch.Lifeperc)
	}
	if ch.Basepl != PFDEF_BASEPL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Bpl : %lld\n", ch.Basepl)
	}
	if ch.Baseki != PFDEF_BASEKI {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Bki : %lld\n", ch.Baseki)
	}
	if ch.Basest != PFDEF_BASEST {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Bst : %lld\n", ch.Basest)
	}
	if ch.Droom != PFDEF_DROOM {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Droo: %d\n", ch.Droom)
	}
	if ch.Hairl != PFDEF_HAIRL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hrl : %d\n", ch.Hairl)
	}
	if ch.Hairs != PFDEF_HAIRS {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hrs : %d\n", ch.Hairs)
	}
	if ch.Hairc != PFDEF_HAIRC {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Hrc : %d\n", ch.Hairc)
	}
	if ch.Skin != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Skn : %d\n", ch.Skin)
	}
	if ch.Eye != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Eye : %d\n", ch.Eye)
	}
	if ch.Distfea != PFDEF_DISTFEA {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Phas: %d\n", ch.Distfea)
	}
	if int(ch.Fury) != PFDEF_HAIRL {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Fury: %d\n", ch.Fury)
	}
	if ch.Radar1 != PFDEF_RADAR1 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Rad1: %d\n", ch.Radar1)
	}
	if ch.Radar2 != PFDEF_RADAR2 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Rad2: %d\n", ch.Radar2)
	}
	if ch.Radar3 != PFDEF_RADAR3 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Rad3: %d\n", ch.Radar3)
	}
	if ch.Ship != PFDEF_SHIP {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Ship: %d\n", ch.Ship)
	}
	if ch.Shipr != PFDEF_SHIPROOM {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Shpr: %d\n", ch.Shipr)
	}
	if ch.Lastpl != PFDEF_LPLAY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lpla: %ld\n", ch.Lastpl)
	}
	if ch.Lastint != PFDEF_LPLAY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lint: %ld\n", ch.Lastint)
	}
	if ch.Deathtime != PFDEF_LPLAY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Deat: %ld\n", ch.Deathtime)
	}
	if ch.Rewtime != PFDEF_LPLAY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Rtim: %ld\n", ch.Rewtime)
	}
	if ch.Boosts != PFDEF_DISTFEA {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boos: %d\n", ch.Boosts)
	}
	if ch.Majinize != PFDEF_LPLAY {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Maji: %d\n", ch.Majinize)
	}
	if ch.Blesslvl != PFDEF_HEIGHT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Blss: %d\n", ch.Blesslvl)
	}
	if (ch.Lboard[0]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boam: %ld\n", ch.Lboard[0])
	}
	if (ch.Lboard[1]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boai: %ld\n", ch.Lboard[1])
	}
	if (ch.Lboard[2]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boac: %ld\n", ch.Lboard[2])
	}
	if (ch.Lboard[3]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boad: %ld\n", ch.Lboard[3])
	}
	if (ch.Lboard[4]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Boab: %ld\n", ch.Lboard[4])
	}
	if (ch.Limb_condition[0]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lira: %d\n", ch.Limb_condition[0])
	}
	if (ch.Limb_condition[1]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lila: %d\n", ch.Limb_condition[1])
	}
	if (ch.Limb_condition[2]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lirl: %d\n", ch.Limb_condition[2])
	}
	if (ch.Limb_condition[3]) != PFDEF_BOARD {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Lill: %d\n", ch.Limb_condition[3])
	}
	if ch.Crank != PFDEF_CRANK {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Clar: %d\n", ch.Crank)
	}
	if ch.Rp != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "RPP : %d\n", ch.Rp)
	}
	if ch.Rbank != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "RBan: %d\n", ch.Rbank)
	}
	if ch.Suppression != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Supp: %lld\n", ch.Suppression)
	}
	if ch.Suppressed != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Sups: %lld\n", ch.Suppressed)
	}
	if ch.Trp != PFDEF_SKIN {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trp : %d\n", ch.Trp)
	}
	if ch.Dcount != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Deac: %d\n", ch.Dcount)
	}
	if ch.Player_specials.Trainagl != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trag: %d\n", ch.Player_specials.Trainagl)
	}
	if ch.Player_specials.Traincon != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trco: %d\n", ch.Player_specials.Traincon)
	}
	if ch.Player_specials.Trainint != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trin: %d\n", ch.Player_specials.Trainint)
	}
	if ch.Player_specials.Trainspd != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trsp: %d\n", ch.Player_specials.Trainspd)
	}
	if ch.Player_specials.Trainstr != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trst: %d\n", ch.Player_specials.Trainstr)
	}
	if ch.Player_specials.Trainwis != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Trwi: %d\n", ch.Player_specials.Trainwis)
	}
	if ch.Starphase != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Phse: %d\n", ch.Starphase)
	}
	if ch.Mimic != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Mimi: %d\n", ch.Mimic)
	}
	if ch.Skill_slots != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Slot: %d\n", ch.Skill_slots)
	}
	if ch.Tail_growth != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Tgro: %d\n", ch.Tail_growth)
	}
	if int(ch.Stupidkiss) != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Stuk: %d\n", ch.Stupidkiss)
	}
	if unsafe.Pointer(ch.Rdisplay) != unsafe.Pointer(uintptr(PFDEF_EYE)) {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "rDis: %s\n", ch.Rdisplay)
	}
	if ch.Relax_count != PFDEF_EYE {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Rela: %d\n", ch.Relax_count)
	}
	if ch.Admlevel < ADMLVL_IMMORT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Skil:\n")
		for i = 1; i <= SKILL_TABLE_SIZE; i++ {
			if (ch.Skills[i]) != 0 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%d %d %d\n", i, ch.Skills[i], ch.Skillperfs[i])
			}
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "0 0\n")
	}
	var buff [200]byte
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Bonu:\n")
	for i = 0; i < 52; i++ {
		if (ch.Bonuses[i]) != 0 && i == 0 {
			stdio.Sprintf(&buff[0], "%d", ch.Bonuses[i])
		} else if (ch.Bonuses[i]) != 0 && i != 0 {
			stdio.Sprintf(&buff[C.strlen(&buff[0])], " %d", ch.Bonuses[i])
		} else if i == 0 {
			stdio.Sprintf(&buff[0], "0")
		} else {
			stdio.Sprintf(&buff[C.strlen(&buff[0])], " 0")
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%s\n", &buff[0])
	buff[0] = '\x00'
	if ch.Admlevel < ADMLVL_IMMORT {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "SklB:\n")
		for i = 1; i <= SKILL_TABLE_SIZE; i++ {
			if (ch.Skillmods[i]) != 0 {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%d %d %d\n", i, ch.Skillmods[i], ch.Skillperfs[i])
			}
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "0 0\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Affs:\n")
	for i = 0; i < MAX_AFFECT; i++ {
		aff = &tmp_aff[i]
		if aff.Type != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "0 0 0 0 0 0\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Affv:\n")
	for i = 0; i < MAX_AFFECT; i++ {
		aff = &tmp_affv[i]
		if aff.Type != 0 {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "%d %d %d %d %d %d\n", aff.Type, aff.Duration, aff.Modifier, aff.Location, int(aff.Bitvector), aff.Specific)
		}
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "0 0 0 0 0 0\n")
	for i = 0; i < NUM_COLOR; i++ {
		if ch.Player_specials.Color_choices[i] != nil {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fl)), "Colr: %d %s\r\n", i, ch.Player_specials.Color_choices[i])
		}
	}
	C.fclose(fl)
	for i = 0; i < MAX_AFFECT; i++ {
		if tmp_aff[i].Type != 0 {
			affect_to_char(ch, &tmp_aff[i])
		}
	}
	if config_info.Operation.Imc_enabled != 0 {
		save_imc_pfile(ch)
	}
	for i = 0; i < MAX_AFFECT; i++ {
		if tmp_affv[i].Type != 0 {
			affectv_to_char(ch, &tmp_affv[i])
		}
	}
	for i = 0; i < NUM_WEARS; i++ {
		if char_eq[i] != nil {
			equip_char(ch, char_eq[i], i)
		}
	}
	if (func() int {
		id = get_ptable_by_name(GET_NAME(ch))
		return id
	}()) < 0 {
		return
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Level != GET_LEVEL(ch) {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Level = GET_LEVEL(ch)
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Admlevel != ch.Admlevel {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Admlevel = ch.Admlevel
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Last != ch.Time.Logon {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Last = ch.Time.Logon
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Played != ch.Lastpl {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Played = ch.Lastpl
	}
	if ch.Clan != nil && (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Clan != ch.Clan {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Clan = C.strdup(ch.Clan)
	}
	if ch.Clan == nil {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Clan = C.strdup(libc.CString("None."))
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Ship != ch.Ship {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Ship = ch.Ship
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Shiproom != int(ch.Shipr) {
		save_index = TRUE
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Shiproom = int(ch.Shipr)
	}
	i = (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags
	if PLR_FLAGGED(ch, PLR_DELETED) {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags |= 1 << 0
	} else {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags &= ^(1 << 0)
	}
	if PLR_FLAGGED(ch, PLR_NODELETE) || PLR_FLAGGED(ch, PLR_CRYO) {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags |= 1 << 1
	} else {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags &= ^(1 << 1)
	}
	if PLR_FLAGGED(ch, PLR_FROZEN) || PLR_FLAGGED(ch, PLR_NOWIZLIST) {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags |= 1 << 3
	} else {
		(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags &= ^(1 << 3)
	}
	if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(id)))).Flags != i || save_index != 0 {
		save_player_index()
	}
}
func save_etext(ch *char_data) {
}
func tag_argument(argument *byte, tag *byte) {
	var (
		tmp  *byte = argument
		ttag *byte = tag
		wrt  *byte = argument
		i    int
	)
	for i = 0; i < 4; i++ {
		*(func() *byte {
			p := &ttag
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(func() *byte {
			p := &tmp
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	*ttag = '\x00'
	for *tmp == ':' || *tmp == ' ' {
		tmp = (*byte)(unsafe.Add(unsafe.Pointer(tmp), 1))
	}
	for *tmp != 0 {
		*(func() *byte {
			p := &wrt
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()) = *(func() *byte {
			p := &tmp
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	*wrt = '\x00'
}
func load_affects(fl *C.FILE, ch *char_data, violence int) {
	var (
		num  int
		num2 int
		num3 int
		num4 int
		num5 int
		num6 int
		i    int
		line [2049]byte
		af   affected_type
	)
	i = 0
	for {
		get_line(fl, &line[0])
		num = func() int {
			num2 = func() int {
				num3 = func() int {
					num4 = func() int {
						num5 = func() int {
							num6 = 0
							return num6
						}()
						return num5
					}()
					return num4
				}()
				return num3
			}()
			return num2
		}()
		__isoc99_sscanf(&line[0], libc.CString("%d %d %d %d %d %d"), &num, &num2, &num3, &num4, &num5, &num6)
		if num != 0 {
			af.Type = int16(num)
			af.Duration = int16(num2)
			af.Modifier = num3
			af.Location = num4
			af.Bitvector = bitvector_t(num5)
			af.Specific = num6
			if violence != 0 {
				affectv_to_char(ch, &af)
			} else {
				affect_to_char(ch, &af)
			}
			i++
		}
		if num == 0 {
			break
		}
	}
}
func load_skills(fl *C.FILE, ch *char_data, mods bool) {
	var (
		num  int = 0
		num2 int = 0
		num3 int = 0
		line [2049]byte
	)
	for {
		get_line(fl, &line[0])
		__isoc99_sscanf(&line[0], libc.CString("%d %d %d"), &num, &num2, &num3)
		if num != 0 {
			if mods {
				for {
					ch.Skillmods[num] = int8(num2)
					if true {
						break
					}
				}
			} else {
				for {
					ch.Skills[num] = int8(num2)
					if true {
						break
					}
				}
			}
			for {
				ch.Skillperfs[num] = int8(num3)
				if true {
					break
				}
			}
		}
		if num == 0 {
			break
		}
	}
}
func load_bonuses(fl *C.FILE, ch *char_data, mods bool) {
	var (
		num  [52]int = [52]int{}
		i    int
		line [2049]byte
	)
	get_line(fl, &line[0])
	__isoc99_sscanf(&line[0], libc.CString("%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"), &num[0], &num[1], &num[2], &num[3], &num[4], &num[5], &num[6], &num[7], &num[8], &num[9], &num[10], &num[11], &num[12], &num[13], &num[14], &num[15], &num[16], &num[17], &num[18], &num[19], &num[20], &num[21], &num[22], &num[23], &num[24], &num[25], &num[26], &num[27], &num[28], &num[29], &num[30], &num[31], &num[32], &num[33], &num[34], &num[35], &num[36], &num[37], &num[38], &num[39], &num[40], &num[41], &num[42], &num[43], &num[44], &num[45], &num[46], &num[47], &num[48], &num[49], &num[50], &num[51])
	for i = 0; i < 52; i++ {
		if num[i] > 0 {
			ch.Bonuses[i] = num[i]
		}
	}
}
func load_feats(fl *C.FILE, ch *char_data) {
	var (
		num  int = 0
		num2 int = 0
		line [2049]byte
	)
	for {
		get_line(fl, &line[0])
		__isoc99_sscanf(&line[0], libc.CString("%d %d"), &num, &num2)
		if num != 0 {
			ch.Feats[num] = int8(num2)
		}
		if num == 0 {
			break
		}
	}
}
func load_HMVS(ch *char_data, line *byte, mode int) {
	var (
		num  int64 = 0
		num2 int64 = 0
	)
	__isoc99_sscanf(line, libc.CString("%lld/%lld"), &num, &num2)
	switch mode {
	case LOAD_HIT:
		ch.Hit = num
		ch.Max_hit = num2
	case LOAD_MANA:
		ch.Mana = num
		ch.Max_mana = num2
	case LOAD_MOVE:
		ch.Move = num
		ch.Max_move = num2
	case LOAD_KI:
		ch.Ki = num
		ch.Max_ki = num2
	}
}
func load_BASE(ch *char_data, line *byte, mode int) {
	var num int64 = 0
	__isoc99_sscanf(line, libc.CString("%lld"), &num)
	switch mode {
	case LOAD_HIT:
		ch.Basepl = num
	case LOAD_MANA:
		ch.Baseki = num
	case LOAD_MOVE:
		ch.Basest = num
	case LOAD_LIFE:
		ch.Lifeforce = num
	}
}
func load_majin(ch *char_data, line *byte) {
	var num int64 = 0
	__isoc99_sscanf(line, libc.CString("%lld"), &num)
	ch.Majinizer = num
}
func load_molt(ch *char_data, line *byte) {
	var num int64 = 0
	__isoc99_sscanf(line, libc.CString("%lld"), &num)
	ch.Moltexp = num
}
func remove_player(pfilepos int) {
	var (
		fname [40]byte
		i     int
	)
	if *(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Name == 0 {
		return
	}
	for i = 0; i < MAX_FILES; i++ {
		if get_filename(&fname[0], uint64(40), i, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Name) != 0 {
			unlink(&fname[0])
		}
		if get_filename(&fname[0], uint64(40), i, CAP((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Name)) != 0 {
			unlink(&fname[0])
		}
	}
	basic_mud_log(libc.CString("PCLEAN: %s Lev: %d Last: %s"), (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Name, (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Level, C.asctime(C.localtime(&(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Last)))
	*(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(pfilepos)))).Name = '\x00'
	save_player_index()
}
func clean_pfiles() {
	var (
		i  int
		ci int
	)
	for i = 0; i <= top_of_p_table; i++ {
		if ((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Flags&(1<<1)) == 0 && *(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Name != 0 {
			if ((*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Flags & (1 << 0)) != 0 {
				remove_player(i)
			} else {
				for ci = 0; pclean_criteria[ci].Level > -1; ci++ {
					if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Admlevel > 1 {
						continue
					} else if (*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Level <= pclean_criteria[ci].Level && (C.time(nil)-(*(*player_index_element)(unsafe.Add(unsafe.Pointer(player_table), unsafe.Sizeof(player_index_element{})*uintptr(i)))).Last) >= int64(pclean_criteria[ci].Days*((int(SECS_PER_REAL_MIN*60))*24)) {
						remove_player(i)
						break
					}
				}
			}
		}
	}
}
