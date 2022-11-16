package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

const MSG_TRAINER_NOT_OPEN = "I'm busy! Come back later!"
const MSG_TRAINER_NO_SEE_CH = "I don't train someone I can't see!"
const MSG_TRAINER_DISLIKE_ALIGN = "Get out of here before I get angry, you are not aligned with me!"
const MSG_TRAINER_DISLIKE_CLASS = "I won't train those of your discipline!"
const MSG_TRAINER_DISLIKE_RACE = "Get out of here, I don't help your kind!"
const MSG_TRAINER_MINLVL = "You are not of a skilled enough level to recieve my training."
const LEARNED_LEVEL = 0
const MAX_PER_PRAC = 1
const MIN_PER_PRAC = 2
const PRAC_TYPE = 3

func copy_guild(tgm *guild_data, fgm *guild_data) {
	var i int
	tgm.Vnum = fgm.Vnum
	tgm.Charge = fgm.Charge
	tgm.Gm = fgm.Gm
	for i = 0; i < SW_ARRAY_MAX; i++ {
		tgm.With_who[i] = fgm.With_who[i]
	}
	tgm.Open = fgm.Open
	tgm.Close = fgm.Close
	tgm.Minlvl = fgm.Minlvl
	tgm.Func = fgm.Func
	free_guild_strings(tgm)
	tgm.No_such_skill = str_udup(fgm.No_such_skill)
	tgm.Not_enough_gold = str_udup(fgm.Not_enough_gold)
	for i = 0; i < SKILL_TABLE_SIZE; i++ {
		tgm.Skills[i] = fgm.Skills[i]
	}
	for i = 0; i < NUM_FEATS_DEFINED; i++ {
		tgm.Feats[i] = fgm.Feats[i]
	}
}
func free_guild_strings(guild *guild_data) {
	if guild.No_such_skill != nil {
		libc.Free(unsafe.Pointer(guild.No_such_skill))
		guild.No_such_skill = nil
	}
	if guild.Not_enough_gold != nil {
		libc.Free(unsafe.Pointer(guild.Not_enough_gold))
		guild.Not_enough_gold = nil
	}
}
func free_guild(guild *guild_data) {
	free_guild_strings(guild)
	libc.Free(unsafe.Pointer(guild))
}
func real_guild(vnum guild_vnum) guild_rnum {
	var (
		bot      guild_rnum
		top      guild_rnum
		mid      guild_rnum
		last_top guild_rnum
	)
	if top_guild < 0 {
		return -1
	}
	bot = 0
	top = guild_rnum(top_guild)
	for {
		last_top = top
		mid = (bot + top) / 2
		if guild_index[mid].Vnum == room_vnum(vnum) {
			return mid
		}
		if bot >= top {
			return -1
		}
		if guild_index[mid].Vnum > room_vnum(vnum) {
			top = mid - 1
		} else {
			bot = mid + 1
		}
		if top > last_top {
			return -1
		}
	}
}
func gedit_modify_string(str **byte, new_g *byte) {
	var (
		pointer *byte
		buf     [64936]byte
	)
	if *new_g != '%' {
		stdio.Snprintf(&buf[0], int(64936), "%%s %s", new_g)
		pointer = &buf[0]
	} else {
		pointer = new_g
	}
	if *str != nil {
		libc.Free(unsafe.Pointer(*str))
	}
	*str = libc.StrDup(pointer)
}
func add_guild(ngld *guild_data) int {
	var (
		rguild guild_rnum
		found  int       = 0
		rznum  zone_rnum = real_zone_by_thing(ngld.Vnum)
	)
	if (func() guild_rnum {
		rguild = real_guild(guild_vnum(ngld.Vnum))
		return rguild
	}()) != guild_rnum(-1) {
		copy_guild(&guild_index[rguild], ngld)
		if rznum != zone_rnum(-1) {
			add_to_save_list(zone_table[rznum].Number, SL_GLD)
		} else {
			mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: GenOLC: Cannot determine guild zone."))
		}
		return int(rguild)
	}
	mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: GenOLC: Creating new guild."))
	top_guild++
	// todo: figure this out
	//guild_index = []guild_data((*guild_data)(libc.Realloc(unsafe.Pointer(&guild_index[0]), top_guild*int(unsafe.Sizeof(guild_data{}))+1)))
	for rguild = guild_rnum(top_guild); rguild > 0; rguild-- {
		if ngld.Vnum > guild_index[rguild-1].Vnum {
			found = int(rguild)
			copy_guild(&guild_index[rguild], ngld)
			break
		}
		guild_index[rguild] = guild_index[rguild-1]
	}
	if found == 0 {
		copy_guild(&guild_index[0], ngld)
	}
	if rznum != zone_rnum(-1) {
		add_to_save_list(zone_table[rznum].Number, SL_GLD)
	} else {
		mudlog(BRF, ADMLVL_BUILDER, TRUE, libc.CString("SYSERR: GenOLC: Cannot determine guild zone."))
	}
	return int(rguild)
}
func save_guilds(zone_num zone_rnum) int {
	var (
		i          int
		j          int
		rguild     int
		guild_file *stdio.File
		fname      [64]byte
		guild      *guild_data
	)
	if zone_num < 0 || zone_num > top_of_zone_table {
		basic_mud_log(libc.CString("SYSERR: GenOLC: save_guilds: Invalid real zone number %d. (0-%d)"), zone_num, top_of_zone_table)
		return FALSE
	}
	stdio.Snprintf(&fname[0], int(64), "%s%d.gld", LIB_WORLD, zone_table[zone_num].Number)
	if (func() *stdio.File {
		guild_file = stdio.FOpen(libc.GoString(&fname[0]), "w")
		return guild_file
	}()) == nil {
		mudlog(BRF, ADMLVL_GOD, TRUE, libc.CString("SYSERR: OLC: Cannot open Guild file!"))
		return FALSE
	}
	for i = int(genolc_zone_bottom(zone_num)); i <= int(zone_table[zone_num].Top); i++ {
		if (func() int {
			rguild = int(real_guild(guild_vnum(i)))
			return rguild
		}()) != int(-1) {
			stdio.Fprintf(guild_file, "#%d~\n", i)
			guild = &guild_index[rguild]
			for j = 0; j < SKILL_TABLE_SIZE; j++ {
				if (guild.Skills[j]) != 0 {
					stdio.Fprintf(guild_file, "%d 1\n", j)
				}
			}
			for j = 0; j < NUM_FEATS_DEFINED; j++ {
				if (guild.Feats[j]) != 0 {
					stdio.Fprintf(guild_file, "%d 2\n", j)
				}
			}
			stdio.Fprintf(guild_file, "-1\n")
			stdio.Fprintf(guild_file, "%1.2f\n", guild.Charge)
			stdio.Fprintf(guild_file, "%s~\n%s~\n", func() *byte {
				if guild.No_such_skill != nil {
					return guild.No_such_skill
				}
				return libc.CString("%s ERROR")
			}(), func() *byte {
				if guild.Not_enough_gold != nil {
					return guild.Not_enough_gold
				}
				return libc.CString("%s ERROR")
			}())
			stdio.Fprintf(guild_file, "%d\n", guild.Minlvl)
			stdio.Fprintf(guild_file, "%d\n%d\n%d\n%d\n", func() mob_vnum {
				if guild.Gm == mob_rnum(-1) {
					return -1
				}
				return mob_index[guild.Gm].Vnum
			}(), guild.With_who[0], guild.Open, guild.Close)
			for j = 1; j < SW_ARRAY_MAX; j++ {
				stdio.Fprintf(guild_file, "%s%d", func() string {
					if j == 1 {
						return ""
					}
					return " "
				}(), guild.With_who[j])
			}
			stdio.Fprintf(guild_file, "\n")
		}
	}
	stdio.Fprintf(guild_file, "$~\n")
	guild_file.Close()
	if in_save_list(zone_table[zone_num].Number, SL_GLD) != 0 {
		remove_from_save_list(zone_table[zone_num].Number, SL_GLD)
		create_world_index(int(zone_table[zone_num].Number), libc.CString("gld"))
		basic_mud_log(libc.CString("GenOLC: save_guilds: Saving guilds '%s'"), &fname[0])
	}
	return TRUE
}
