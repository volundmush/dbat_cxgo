package main

import "github.com/gotranspile/cxgo/runtime/libc"

func pre_reset(znum int) bool {
	var ret_value bool = PERFORM_NORMAL_RESET != 0
	switch znum {
	case RESET_GAUNTLET:
		ret_value = prereset_gauntlet_zone()
	default:
		ret_value = PERFORM_NORMAL_RESET != 0
	}
	return ret_value
}
func post_reset(znum int) {
	switch znum {
	default:
	}
}
func prereset_gauntlet_zone() bool {
	var (
		i                int
		gauntlet_players int = 0
		mob              *char_data
		gauntlet_mobs    [20]gauntlet_mob = [20]gauntlet_mob{{Vroom: 2403, Vmob: 2400}, {Vroom: 2405, Vmob: 2401}, {Vroom: 2407, Vmob: 2402}, {Vroom: 2409, Vmob: 2403}, {Vroom: 2411, Vmob: 2404}, {Vroom: 2413, Vmob: 2405}, {Vroom: 2415, Vmob: 2406}, {Vroom: 2417, Vmob: 2407}, {Vroom: 2419, Vmob: 2408}, {Vroom: 2421, Vmob: 2409}, {Vroom: 2423, Vmob: 2410}, {Vroom: 2425, Vmob: 2411}, {Vroom: 2427, Vmob: 2412}, {Vroom: 2429, Vmob: 2413}, {Vroom: 2431, Vmob: 2414}, {Vroom: 2433, Vmob: 2415}, {Vroom: 2435, Vmob: 2416}, {Vroom: 2437, Vmob: 2417}, {Vroom: 2439, Vmob: 2418}, {Vroom: 2441, Vmob: 2419}}
	)
	basic_mud_log(libc.CString("Special Reset: zone %d: Resetting Gauntlet"), RESET_GAUNTLET)
	for i = 0; i < NUM_GAUNTLET_ROOMS; i++ {
		gauntlet_players += num_players_in_room(gauntlet_mobs[i].Vroom)
	}
	if gauntlet_players == 0 {
		basic_mud_log(libc.CString("Special Reset: zone %d: No players in Gauntlet - executing normal reset"), RESET_GAUNTLET)
		return PERFORM_NORMAL_RESET != 0
	}
	basic_mud_log(libc.CString("Special Reset: zone %d: %d players in Gauntlet - special reset only"), RESET_GAUNTLET, gauntlet_players)
	for i = 0; i < NUM_GAUNTLET_ROOMS; i++ {
		if !check_mob_in_room(gauntlet_mobs[i].Vmob, gauntlet_mobs[i].Vroom) {
			if num_players_in_room(gauntlet_mobs[i].Vroom) == 0 {
				if real_mobile(gauntlet_mobs[i].Vmob) != 0 && real_room(gauntlet_mobs[i].Vroom) != 0 {
					if (func() *char_data {
						mob = read_mobile(gauntlet_mobs[i].Vmob, VIRTUAL)
						return mob
					}()) != nil {
						char_to_room(mob, real_room(gauntlet_mobs[i].Vroom))
						basic_mud_log(libc.CString("Special Reset: zone %d: Gauntlet mob reset (%d, %s)"), RESET_GAUNTLET, gauntlet_mobs[i].Vmob, GET_NAME(mob))
					}
				}
			}
		}
	}
	return BLOCK_NORMAL_RESET != 0
}
