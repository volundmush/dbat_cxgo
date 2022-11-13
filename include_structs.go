package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type extra_descr_data struct {
	Keyword     *byte
	Description *byte
	Next        *extra_descr_data
}
type obj_affected_type struct {
	Location int
	Specific int
	Modifier int
}
type obj_spellbook_spell struct {
	Spellname int
	Pages     int
}
type obj_data struct {
	Item_number        obj_vnum
	In_room            room_rnum
	Room_loaded        room_vnum
	Value              [16]int
	Type_flag          int8
	Level              int
	Wear_flags         [4]int
	Extra_flags        [4]bitvector_t
	Weight             int64
	Cost               int
	Cost_per_day       int
	Timer              int
	Bitvector          [4]bitvector_t
	Size               int
	Affected           [6]obj_affected_type
	Name               *byte
	Description        *byte
	Short_description  *byte
	Action_description *byte
	Ex_description     *extra_descr_data
	Carried_by         *char_data
	Worn_by            *char_data
	Worn_on            int16
	In_obj             *obj_data
	Contains           *obj_data
	Id                 int32
	Generation         libc.Time
	Unique_id          int64
	Proto_script       *trig_proto_list
	Script             *script_data
	Next_content       *obj_data
	Next               *obj_data
	Sbinfo             *obj_spellbook_spell
	Sitting            *char_data
	Scoutfreq          int
	Lload              libc.Time
	Healcharge         int
	Kicharge           int64
	Kitype             int
	User               *char_data
	Target             *char_data
	Distance           int
	Foob               int
	Aucter             int32
	CurBidder          int32
	AucTime            libc.Time
	Bid                int
	Startbid           int
	Auctname           *byte
	Posttype           int
	Posted_to          *obj_data
	Fellow_wall        *obj_data
}
type room_direction_data struct {
	General_description *byte
	Keyword             *byte
	Exit_info           bitvector_t
	Key                 obj_vnum
	To_room             room_rnum
	Dclock              int
	Dchide              int
	Dcskill             int
	Dcmove              int
	Failsavetype        int
	Dcfailsave          int
	Failroom            int
	Totalfailroom       int
}
type room_data struct {
	Number         room_vnum
	Zone           zone_rnum
	Sector_type    int
	Name           *byte
	Description    *byte
	Ex_description *extra_descr_data
	Dir_option     [12]*room_direction_data
	Room_flags     [4]bitvector_t
	Proto_script   *trig_proto_list
	Script         *script_data
	Light          int8
	Func           func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int
	Contents       *obj_data
	People         *char_data
	Timed          int
	Dmg            int
	Gravity        int
	Geffect        int
}
type memory_rec_struct struct {
	Id   int32
	Next *memory_rec_struct
}
type memory_rec memory_rec_struct
type time_info_data struct {
	Hours int
	Day   int
	Month int
	Year  int16
}
type time_data struct {
	Birth   libc.Time
	Created libc.Time
	Maxage  libc.Time
	Logon   libc.Time
	Played  libc.Time
}
type pclean_criteria_data struct {
	Level int
	Days  int
}
type abil_data struct {
	Str   int8
	Intel int8
	Wis   int8
	Dex   int8
	Con   int8
	Cha   int8
}
type player_special_data struct {
	Poofin                 *byte
	Poofout                *byte
	Aliases                *alias_data
	Last_tell              int32
	Last_olc_targ          unsafe.Pointer
	Last_olc_mode          int
	Host                   *byte
	Spell_level            [10]int
	Memcursor              int
	Wimp_level             int
	Freeze_level           int8
	Invis_level            int16
	Load_room              room_vnum
	Pref                   [4]bitvector_t
	Bad_pws                uint8
	Conditions             [3]int8
	Skill_points           int
	Class_skill_points     [31]int
	Comm_hist              [10]*txt_block
	Olc_zone               int
	Gauntlet               int
	Speaking               int
	Tlevel                 int
	Ability_trains         int
	Spellmem               [100]int
	Feat_points            int
	Epic_feat_points       int
	Class_feat_points      [31]int
	Epic_class_feat_points [31]int
	Domain                 [37]int
	School                 [10]int
	Deity                  int
	Spell_mastery_points   int
	Color_choices          [16]*byte
	Page_length            uint8
	Murder                 int
	Trainstr               int
	Trainint               int
	Traincon               int
	Trainwis               int
	Trainagl               int
	Trainspd               int
	Carrying               *char_data
	Carried_by             *char_data
	Racial_pref            int
}
type memorize_node struct {
	Timer int
	Spell int
	Next  *memorize_node
}
type innate_node struct {
	Timer    int
	Spellnum int
	Next     *innate_node
}
type mob_special_data struct {
	Memory      *memory_rec
	Attack_type int8
	Default_pos int8
	Damnodice   int8
	Damsizedice int8
	Newitem     int
}
type affected_type struct {
	Type      int16
	Duration  int16
	Modifier  int
	Location  int
	Specific  int
	Bitvector bitvector_t
	Next      *affected_type
}
type queued_act struct {
	Level    int
	Spellnum int
}
type follow_type struct {
	Follower *char_data
	Next     *follow_type
}
type level_learn_entry struct {
	Next     *level_learn_entry
	Location int
	Specific int
	Value    int8
}
type levelup_data struct {
	Next                 *levelup_data
	Prev                 *levelup_data
	Type                 int8
	Spec                 int8
	Level                int8
	Hp_roll              int8
	Mana_roll            int8
	Ki_roll              int8
	Move_roll            int8
	Accuracy             int8
	Fort                 int8
	Reflex               int8
	Will                 int8
	Add_skill            int8
	Add_gen_feats        int8
	Add_epic_feats       int8
	Add_class_feats      int8
	Add_class_epic_feats int8
	Skills               *level_learn_entry
	Feats                *level_learn_entry
}
type char_data struct {
	Pfilepos           int
	Nr                 mob_rnum
	In_room            room_rnum
	Was_in_room        room_rnum
	Wait               int
	Name               *byte
	Short_descr        *byte
	Long_descr         *byte
	Description        *byte
	Title              *byte
	Size               int
	Sex                int8
	Race               int8
	Hairl              int8
	Hairs              int8
	Hairc              int8
	Skin               int8
	Eye                int8
	Distfea            int8
	Race_level         int
	Level_adj          int
	Chclass            int8
	Chclasses          [31]int
	Epicclasses        [31]int
	Level_info         *levelup_data
	Level              int
	Admlevel           int
	Admflags           [4]bitvector_t
	Hometown           room_vnum
	Time               time_data
	Weight             uint8
	Height             uint8
	Real_abils         abil_data
	Aff_abils          abil_data
	Player_specials    *player_special_data
	Mob_specials       mob_special_data
	Affected           *affected_type
	Affectedv          *affected_type
	Actq               *queued_act
	Equipment          [23]*obj_data
	Carrying           *obj_data
	Desc               *descriptor_data
	Id                 int32
	Proto_script       *trig_proto_list
	Script             *script_data
	Memory             *script_memory
	Next_in_room       *char_data
	Next               *char_data
	Next_fighting      *char_data
	Next_affect        *char_data
	Next_affectv       *char_data
	Followers          *follow_type
	Master             *char_data
	Master_id          int32
	Memorized          *memorize_node
	Innate             *innate_node
	Fighting           *char_data
	Position           int8
	Carry_weight       int
	Carry_items        int8
	Timer              int
	Sits               *obj_data
	Blocks             *char_data
	Blocked            *char_data
	Absorbing          *char_data
	Absorbby           *char_data
	Feats              [751]int8
	Combat_feats       [6][4]int
	School_feats       [2]int
	Skills             [1001]int8
	Skillmods          [1001]int8
	Skillperfs         [1001]int8
	Alignment          int
	Alignment_ethic    int
	Idnum              int32
	Act                [4]bitvector_t
	Affected_by        [4]int
	Bodyparts          [4]int
	Saving_throw       [3]int16
	Apply_saving_throw [3]int16
	Powerattack        int
	Combatexpertise    int
	Mana               int64
	Max_mana           int64
	Hit                int64
	Max_hit            int64
	Move               int64
	Max_move           int64
	Ki                 int64
	Max_ki             int64
	Armor              int
	Shield_bonus       int16
	Gold               int
	Bank_gold          int
	Exp                int64
	Accuracy           int
	Accuracy_mod       int
	Damage_mod         int
	Spellfail          int16
	Armorcheck         int16
	Armorcheckall      int16
	Basepl             int64
	Baseki             int64
	Basest             int64
	Charge             int64
	Chargeto           int64
	Barrier            int64
	Clan               *byte
	Droom              room_vnum
	Choice             int
	Sleeptime          int
	Foodr              int
	Altitude           int
	Overf              int
	Spam               int
	Radar1             room_vnum
	Radar2             room_vnum
	Radar3             room_vnum
	Ship               int
	Shipr              room_vnum
	Lastpl             libc.Time
	Lboard             [5]libc.Time
	Listenroom         room_vnum
	Crank              int
	Kaioken            int
	Absorbs            int
	Boosts             int
	Upgrade            int
	Lastint            libc.Time
	Majinize           int
	Fury               int16
	Btime              int16
	Eavesdir           int
	Deathtime          libc.Time
	Rp                 int
	Suppression        int64
	Suppressed         int64
	Drag               *char_data
	Dragged            *char_data
	Trp                int
	Mindlink           *char_data
	Lasthit            int
	Dcount             int
	Voice              *byte
	Limbs              [4]int
	Aura               int
	Rewtime            libc.Time
	Grappling          *char_data
	Grappled           *char_data
	Grap               int
	Genome             [2]int
	Combo              int
	Lastattack         int
	Combhits           int
	Ping               int
	Starphase          int
	Mimic              int
	Bonuses            [52]int
	Ccpoints           int
	Negcount           int
	Cooldown           int
	Death_type         int
	Moltexp            int64
	Moltlevel          int
	Loguser            *byte
	Arenawatch         int
	Majinizer          int64
	Speedboost         int
	Skill_slots        int
	Tail_growth        int
	Rage_meter         int
	Feature            *byte
	Transclass         int
	Transcost          [6]int
	Armor_last         int
	Forgeting          int
	Forgetcount        int
	Backstabcool       int
	Con_cooldown       int
	Stupidkiss         int16
	Temp_prompt        *byte
	Personality        int
	Combine            int
	Linker             int
	Fishstate          int
	Throws             int
	Defender           *char_data
	Defending          *char_data
	Lifeforce          int64
	Lifeperc           int
	Gooptime           int
	Blesslvl           int
	Poisonby           *char_data
	Mobcharge          int
	Preference         int
	Aggtimer           int
	Lifebonus          int
	Asb                int
	Regen              int
	Rbank              int
	Con_sdcooldown     int
	Limb_condition     [4]int
	Rdisplay           *byte
	Song               int16
	Original           *char_data
	Clones             int16
	Relax_count        int
	IngestLearned      int
}
type txt_block struct {
	Text    *byte
	Aliased int
	Next    *txt_block
}
type txt_q struct {
	Head *txt_block
	Tail *txt_block
}
type descriptor_data struct {
	Descriptor     int
	Host           [41]byte
	Bad_pws        int8
	Idle_tics      int8
	Connected      int
	Desc_num       int
	Login_time     libc.Time
	Showstr_head   *byte
	Showstr_vector **byte
	Showstr_count  int
	Showstr_page   int
	Str            **byte
	Backstr        *byte
	Max_str        uint64
	Mail_to        int32
	Has_prompt     int
	Inbuf          [4096]byte
	Last_input     [2048]byte
	Small_outbuf   [6020]byte
	Output         *byte
	History        **byte
	History_pos    int
	Bufptr         int
	Bufspace       int
	Large_outbuf   *txt_block
	Input          txt_q
	Character      *char_data
	Original       *char_data
	Snooping       *descriptor_data
	Snoop_by       *descriptor_data
	Next           *descriptor_data
	Olc            *oasis_olc_data
	User           *byte
	Email          *byte
	Pass           *byte
	Loadplay       *byte
	Writenew       int
	Total          int
	Rpp            int
	Tmp1           *byte
	Tmp2           *byte
	Tmp3           *byte
	Tmp4           *byte
	Tmp5           *byte
	Level          int
	Newsbuf        *byte
	Obj_editval    int
	Obj_editflag   int
	Obj_was        *byte
	Obj_name       *byte
	Obj_short      *byte
	Obj_long       *byte
	Obj_type       int
	Obj_weapon     int
	Obj_point      *obj_data
	Shipmenu       int
	Shipsize       int
	Ship_name      *byte
	Shipextra      [4]int
	Shields        int
	Armor          int
	Drive          int
	Shipweap       int
	User_freed     int
	Customfile     int
	Title          *byte
	Rbank          int
}
type msg_type struct {
	Attacker_msg *byte
	Victim_msg   *byte
	Room_msg     *byte
}
type message_type struct {
	Die_msg  msg_type
	Miss_msg msg_type
	Hit_msg  msg_type
	God_msg  msg_type
	Next     *message_type
}
type message_list struct {
	A_type            int
	Number_of_attacks int
	Msg               *message_type
}
type social_messg struct {
	Act_nr              int
	Command             *byte
	Sort_as             *byte
	Hide                int
	Min_victim_position int
	Min_char_position   int
	Min_level_char      int
	Char_no_arg         *byte
	Others_no_arg       *byte
	Char_found          *byte
	Others_found        *byte
	Vict_found          *byte
	Char_body_found     *byte
	Others_body_found   *byte
	Vict_body_found     *byte
	Not_found           *byte
	Char_auto           *byte
	Others_auto         *byte
	Char_obj_found      *byte
	Others_obj_found    *byte
}
type weather_data struct {
	Pressure int
	Change   int
	Sky      int
	Sunlight int
}
type index_data struct {
	Vnum   mob_vnum
	Number int
	Func   func(ch *char_data, me unsafe.Pointer, cmd int, argument *byte) int
	Farg   *byte
	Proto  *trig_data
}
type trig_proto_list struct {
	Vnum int
	Next *trig_proto_list
}
type guild_info_type struct {
	Pc_class   int
	Guild_room room_vnum
	Direction  int
}
type game_data struct {
	Pk_allowed          int
	Pt_allowed          int
	Level_can_shout     int
	Holler_move_cost    int
	Tunnel_size         int
	Max_exp_gain        int
	Max_exp_loss        int
	Max_npc_corpse_time int
	Max_pc_corpse_time  int
	Idle_void           int
	Idle_rent_time      int
	Idle_max_level      int
	Dts_are_dumps       int
	Load_into_inventory int
	Track_through_doors int
	Level_cap           int
	Stack_mobs          int
	Stack_objs          int
	Mob_fighting        int
	OK                  *byte
	NOPERSON            *byte
	NOEFFECT            *byte
	Disp_closed_doors   int
	Reroll_player       int
	Initial_points      int
	Enable_compression  int
	Enable_languages    int
	All_items_unique    int
	Exp_multiplier      float32
}
type crash_save_data struct {
	Free_rent          int
	Max_obj_save       int
	Min_rent_cost      int
	Auto_save          int
	Autosave_time      int
	Crash_file_timeout int
	Rent_file_timeout  int
}
type room_numbers struct {
	Mortal_start_room room_vnum
	Immort_start_room room_vnum
	Frozen_start_room room_vnum
	Donation_room_1   room_vnum
	Donation_room_2   room_vnum
	Donation_room_3   room_vnum
}
type game_operation struct {
	DFLT_PORT          uint16
	DFLT_IP            *byte
	DFLT_DIR           *byte
	LOGNAME            *byte
	Max_playing        int
	Max_filesize       int
	Max_bad_pws        int
	Siteok_everyone    int
	Nameserver_is_slow int
	Use_new_socials    int
	Auto_save_olc      int
	MENU               *byte
	WELC_MESSG         *byte
	START_MESSG        *byte
	Imc_enabled        int
}
type autowiz_data struct {
	Use_autowiz     int
	Min_wizlist_lev int
}
type tick_data struct {
	Pulse_violence int
	Pulse_mobile   int
	Pulse_zone     int
	Pulse_autosave int
	Pulse_idlepwd  int
	Pulse_sanity   int
	Pulse_usage    int
	Pulse_timesave int
	Pulse_current  int
}
type advance_data struct {
	Allow_multiclass int
	Allow_prestige   int
}
type creation_data struct {
	Method int
}
type config_data struct {
	CONFFILE  *byte
	Play      game_data
	Csd       crash_save_data
	Room_nums room_numbers
	Operation game_operation
	Autowiz   autowiz_data
	Advance   advance_data
	Ticks     tick_data
	Creation  creation_data
}
type aging_data struct {
	Adult     int
	Classdice [3][2]int
	Middle    int
	Old       int
	Venerable int
	Maxdice   [2]int
}
