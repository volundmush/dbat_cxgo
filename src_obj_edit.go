package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

var custom_types [17]*byte = [17]*byte{libc.CString("Error"), libc.CString("Armor Slot"), libc.CString("Gi Slot"), libc.CString("Wrist Slot"), libc.CString("Ear Slot"), libc.CString("Finger Slot"), libc.CString("Eye Slot"), libc.CString("Hands Slot"), libc.CString("Feet Slot"), libc.CString("Belt Slot"), libc.CString("Legs Slot"), libc.CString("Arms Slot"), libc.CString("Head Slot"), libc.CString("Neck Slot"), libc.CString("Back Slot"), libc.CString("Shoulder Slot"), libc.CString("Weapon")}
var custom_weapon [7]*byte = [7]*byte{libc.CString("Not A Weapon"), libc.CString("Sword"), libc.CString("Dagger"), libc.CString("Spear"), libc.CString("Club"), libc.CString("Gun"), libc.CString("Brawling")}

func disp_custom_menu(d *descriptor_data) {
	write_to_output(d, libc.CString("@GCustom Equipment Construction Menu@n\n"))
	write_to_output(d, libc.CString("@D-----------------------------------------@n\n"))
	write_to_output(d, libc.CString("@C1@D) @WKeyword Name     @D: @w%s@n\n"), d.Obj_name)
	write_to_output(d, libc.CString("@C2@D) @WShort Description@D: @w%s@n\n"), d.Obj_short)
	write_to_output(d, libc.CString("@C3@D) @WLong Description @D: @w%s@n\n"), d.Obj_long)
	write_to_output(d, libc.CString("@C4@D) @WEquipment Type   @D: @w%s@n\n"), custom_types[d.Obj_type])
	write_to_output(d, libc.CString("@C5@D) @WWeapon Type      @D: @w%s@n\n"), custom_weapon[d.Obj_weapon])
	write_to_output(d, libc.CString("@CQ@D) @WQuit Menu@n\n"))
}
func disp_restring_menu(d *descriptor_data) {
	write_to_output(d, libc.CString("@GEquipment Restring Menu@n\n"))
	write_to_output(d, libc.CString("@D-----------------------------------------@n\n"))
	write_to_output(d, libc.CString("@C1@D) @WKeyword Name     @D: @w%s@n\n"), d.Obj_name)
	write_to_output(d, libc.CString("@C2@D) @WShort Description@D: @w%s@n\n"), d.Obj_short)
	write_to_output(d, libc.CString("@C3@D) @WLong Description @D: @w%s@n\n"), d.Obj_long)
	write_to_output(d, libc.CString("@CQ@D) @WQuit Menu@n\n"))
}
func pobj_edit_parse(d *descriptor_data, arg *byte) {
	var (
		obj  *obj_data = nil
		buf  [2048]byte
		buf2 [2048]byte
		buf3 [2048]byte
	)
	if d.Obj_editflag == EDIT_RESTRING {
		switch d.Obj_editval {
		case EDIT_RESTRING_MAIN:
			switch *arg {
			case '1':
				write_to_output(d, libc.CString("@wCurrent Equipment's Keyword Name @D<@Y%s@D>@n\r\n"), d.Obj_name)
				write_to_output(d, libc.CString("@wEnter Name @D(@RDo not use colorcode!@D)@w:@n \n"))
				d.Obj_editval = EDIT_RESTRING_NAME
			case '2':
				write_to_output(d, libc.CString("@wCurrent Equipment's Short Desc @D<@Y%s@D>@n\r\n"), d.Obj_short)
				write_to_output(d, libc.CString("@wEnter Short Desc@w:@n \n"))
				d.Obj_editval = EDIT_RESTRING_SDESC
			case '3':
				write_to_output(d, libc.CString("@wCurrent Equipment's Long Desc @D<@Y%s@D>@n\r\n"), d.Obj_long)
				write_to_output(d, libc.CString("@wEnter Long Desc@w:@n \n"))
				d.Obj_editval = EDIT_RESTRING_LDESC
			case 'Q':
				fallthrough
			case 'q':
				write_to_output(d, libc.CString("Save current changes and be charged?\r\nYes or No\r\n"))
				d.Obj_editval = EDIT_RESTRING_QUIT
			default:
				disp_restring_menu(d)
				return
			}
		case EDIT_RESTRING_NAME:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_restring_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
				return
			} else if libc.StrStr(arg, libc.CString("@")) != nil {
				write_to_output(d, libc.CString("@RNO COLORCODE IN THE KEYWORD NAME.@n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else if libc.StrLen(arg) > 100 {
				write_to_output(d, libc.CString("@wToo long. Limit is 100 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_name != nil {
					libc.Free(unsafe.Pointer(d.Obj_name))
				}
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "%s", arg)
				d.Obj_name = libc.StrDup(&buf[0])
				disp_restring_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
			}
		case EDIT_RESTRING_SDESC:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
				return
			} else if libc.StrLen(arg) > 150 {
				write_to_output(d, libc.CString("@wToo long. Limit is 200 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_short != nil {
					libc.Free(unsafe.Pointer(d.Obj_short))
				}
				buf[0] = '\x00'
				stdio.Snprintf(&buf[0], int(2048), "%s", arg)
				d.Obj_short = libc.StrDup(&buf[0])
				disp_restring_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
			}
		case EDIT_RESTRING_LDESC:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
				return
			} else if libc.StrLen(arg) > 200 {
				write_to_output(d, libc.CString("@wToo long. Limit is 200 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_long != nil {
					libc.Free(unsafe.Pointer(d.Obj_long))
				}
				buf[0] = '\x00'
				stdio.Snprintf(&buf[0], int(2048), "%s", arg)
				d.Obj_long = libc.StrDup(&buf[0])
				disp_restring_menu(d)
				d.Obj_editval = EDIT_RESTRING_MAIN
			}
		case EDIT_RESTRING_QUIT:
			if *arg == 0 {
				write_to_output(d, libc.CString("Save current changes and be charged?\r\nYes or No\r\n"))
				return
			} else if libc.StrCaseCmp(arg, libc.CString("yes")) == 0 || libc.StrCaseCmp(arg, libc.CString("Yes")) == 0 || libc.StrCaseCmp(arg, libc.CString("y")) == 0 || libc.StrCaseCmp(arg, libc.CString("Y")) == 0 {
				obj = d.Obj_point
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "%s", d.Obj_name)
				obj.Name = libc.StrDup(&buf[0])
				buf2[0] = '\x00'
				stdio.Sprintf(&buf2[0], "%s", d.Obj_short)
				obj.Short_description = libc.StrDup(&buf2[0])
				buf3[0] = '\x00'
				stdio.Sprintf(&buf3[0], "%s", d.Obj_long)
				obj.Description = libc.StrDup(&buf3[0])
				d.Obj_editflag = EDIT_NONE
				d.Obj_editval = EDIT_NONE
				d.Character.Rbank -= 2
				d.Character.Rbank = d.Character.Rbank
				d.Rbank -= 2
				userWrite(d, 0, 0, 0, libc.CString("index"))
				SET_BIT_AR(obj.Extra_flags[:], ITEM_RESTRING)
				write_to_output(d, libc.CString("Purchase complete."))
				send_to_imm(libc.CString("Restring Eq: %s has bought: %s, which was %s."), GET_NAME(d.Character), obj.Short_description, d.Obj_was)
				d.Connected = CON_PLAYING
			} else if libc.StrCaseCmp(arg, libc.CString("No")) == 0 || libc.StrCaseCmp(arg, libc.CString("no")) == 0 || libc.StrCaseCmp(arg, libc.CString("n")) == 0 || libc.StrCaseCmp(arg, libc.CString("N")) == 0 {
				write_to_output(d, libc.CString("Canceling purchase at no cost.\r\n"))
				send_to_imm(libc.CString("Restring Eq: %s has canceled their equipment restring."), GET_NAME(d.Character))
				d.Obj_editval = EDIT_NONE
				d.Connected = CON_PLAYING
			} else {
				write_to_output(d, libc.CString("Save current changes and be charged?\r\nYes or No\r\n"))
				return
			}
		default:
			d.Obj_editval = EDIT_RESTRING_MAIN
			disp_restring_menu(d)
			return
		}
	}
	if d.Obj_editflag == EDIT_CUSTOM {
		switch d.Obj_editval {
		case EDIT_CUSTOM_MAIN:
			switch *arg {
			case '1':
				write_to_output(d, libc.CString("@wEnter Equipment's Keyword Name @D(@RDo not use colorcode!@D)@w.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				d.Obj_editval = EDIT_CUSTOM_NAME
			case '2':
				write_to_output(d, libc.CString("@wEnter Equipment's Short Description. This is the colored name you see in your inventory or when the eq is used.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				d.Obj_editval = EDIT_CUSTOM_SDESC
			case '3':
				write_to_output(d, libc.CString("@wEnter Equipment's Long Description. This is the colored name you see when it's seen in a room. @D(@RNo punctuation at the end!@D)@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				d.Obj_editval = EDIT_CUSTOM_LDESC
			case '4':
				write_to_output(d, libc.CString("@wEnter number to select type of equipment you want it to be: @n\n"))
				write_to_output(d, libc.CString("@D[ @C--1@W) @cArmor Slot  @C--2@W) @cGi Slot     @C--3@W) @cWrist Slot   @D]@n\n"))
				write_to_output(d, libc.CString("@D[ @C--4@W) @cEar Slot    @C--5@W) @cFinger Slot @C--6@W) @cEye Slot     @D]@n\n"))
				write_to_output(d, libc.CString("@D[ @C--7@W) @cHands Slot  @C--8@W) @cFeet Slot   @C--9@W) @cBelt Slot    @D]@n\n"))
				write_to_output(d, libc.CString("@D[ @C-10@W) @cLegs Slot   @C-11@W) @cArms Slot   @C-12@W) @cHead Slot    @D]@n\n"))
				write_to_output(d, libc.CString("@D[ @C-13@W) @cNeck Slot   @C-14@W) @cBack Slot   @C-15@W) @cShoulder Slot@D]@n\n"))
				write_to_output(d, libc.CString("@D[ @C-16@W) @cWeapon@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				d.Obj_editval = EDIT_CUSTOM_TYPE
			case '5':
				if d.Obj_type != 16 {
					write_to_output(d, libc.CString("@wYou can only use this part of the menu if you select the weapon type.@n\r\n"))
					return
				} else {
					write_to_output(d, libc.CString("@wEnter number to select type of weapon you want it to be: @n\n"))
					write_to_output(d, libc.CString("@D[ @C--1@W) @cSword       @C--2@W) @cDagger      @C--3@W) @cSpear        @D]@n\n"))
					write_to_output(d, libc.CString("@D[ @C--4@W) @cClub        @C--5@W) @cGun         @C--6@W) @cBrawling     @D]@n\n"))
					write_to_output(d, libc.CString("@wEnter: @n\n"))
					d.Obj_editval = EDIT_CUSTOM_WEAPON
				}
			case 'q':
				fallthrough
			case 'Q':
				write_to_output(d, libc.CString("@wPurchase this custom piece? (Y or N)@n\r\n"))
				d.Obj_editval = EDIT_CUSTOM_QUIT
			default:
				disp_custom_menu(d)
				return
			}
		case EDIT_CUSTOM_NAME:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
				return
			} else if libc.StrStr(arg, libc.CString("@")) != nil {
				write_to_output(d, libc.CString("@RNO COLORCODE IN THE KEYWORD NAME.@n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else if libc.StrLen(arg) > 100 {
				write_to_output(d, libc.CString("@wToo long. Limit is 100 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_name != nil {
					libc.Free(unsafe.Pointer(d.Obj_name))
				}
				buf[0] = '\x00'
				stdio.Sprintf(&buf[0], "%s", arg)
				d.Obj_name = libc.StrDup(&buf[0])
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
			}
		case EDIT_CUSTOM_SDESC:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
				return
			} else if libc.StrLen(arg) > 150 {
				write_to_output(d, libc.CString("@wToo long. Limit is 200 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_short != nil {
					libc.Free(unsafe.Pointer(d.Obj_short))
				}
				buf[0] = '\x00'
				stdio.Snprintf(&buf[0], int(2048), "%s", arg)
				d.Obj_short = libc.StrDup(&buf[0])
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
			}
		case EDIT_CUSTOM_LDESC:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Keeping what was previously entered.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
				return
			} else if libc.StrLen(arg) > 200 {
				write_to_output(d, libc.CString("@wToo long. Limit is 200 character.@n\r\n"))
				write_to_output(d, libc.CString("@wEnter: @n\n"))
				return
			} else {
				if d.Obj_long != nil {
					libc.Free(unsafe.Pointer(d.Obj_long))
				}
				buf[0] = '\x00'
				stdio.Snprintf(&buf[0], int(2048), "%s", arg)
				d.Obj_long = libc.StrDup(&buf[0])
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
			}
		case EDIT_CUSTOM_TYPE:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Returning to main menu.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
				return
			} else {
				d.Obj_type = libc.Atoi(libc.GoString(arg))
				if d.Obj_type < 1 || d.Obj_type > 16 {
					write_to_output(d, libc.CString("@wValue must be between 1 and 16.@n\r\n"))
					write_to_output(d, libc.CString("@wEnter: @n\n"))
					d.Obj_type = 1
					return
				} else {
					if d.Obj_type == 16 {
						d.Obj_weapon = 1
					}
					disp_custom_menu(d)
					d.Obj_editval = EDIT_CUSTOM_MAIN
				}
			}
		case EDIT_CUSTOM_WEAPON:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wNothing entered. Returning to main menu.@n\r\n"))
				disp_custom_menu(d)
				d.Obj_editval = EDIT_CUSTOM_MAIN
				return
			} else {
				d.Obj_weapon = libc.Atoi(libc.GoString(arg))
				if d.Obj_weapon < 1 || d.Obj_weapon > 6 {
					write_to_output(d, libc.CString("@wValue must be between 1 and 6.@n\r\n"))
					write_to_output(d, libc.CString("@wEnter: @n\n"))
					d.Obj_weapon = 0
					return
				} else {
					disp_custom_menu(d)
					d.Obj_editval = EDIT_CUSTOM_MAIN
				}
			}
		case EDIT_CUSTOM_QUIT:
			if *arg == 0 {
				write_to_output(d, libc.CString("@wPurchase this custom piece? (Y or N)@n\r\n"))
				return
			} else if libc.StrCaseCmp(arg, libc.CString("y")) == 0 || libc.StrCaseCmp(arg, libc.CString("Y")) == 0 {
				write_to_output(d, libc.CString("@wPurchase complete.@n\r\n"))
				d.Connected = CON_PLAYING
				if d.Obj_weapon == 0 {
					obj = read_object(0x4E83, VIRTUAL)
					obj_to_char(obj, d.Character)
					switch d.Obj_type {
					case 1:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_BODY)
					case 2:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_ABOUT)
					case 3:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_WRIST)
					case 4:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_EAR)
					case 5:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_FINGER)
					case 6:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_EYE)
					case 7:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_HANDS)
					case 8:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_FEET)
					case 9:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_WAIST)
					case 10:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_LEGS)
					case 11:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_ARMS)
					case 12:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_HEAD)
					case 13:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_NECK)
					case 14:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_PACK)
					case 15:
						SET_BIT_AR(obj.Wear_flags[:], ITEM_WEAR_SH)
					}
					buf[0] = '\x00'
					stdio.Sprintf(&buf[0], libc.GoString(d.Obj_name))
					obj.Name = libc.StrDup(&buf[0])
					buf2[0] = '\x00'
					stdio.Sprintf(&buf2[0], libc.GoString(d.Obj_short))
					obj.Short_description = libc.StrDup(&buf2[0])
					buf3[0] = '\x00'
					stdio.Sprintf(&buf3[0], libc.GoString(d.Obj_long))
					obj.Description = libc.StrDup(&buf3[0])
				} else {
					obj = read_object(0x4E82, VIRTUAL)
					obj_to_char(obj, d.Character)
					buf[0] = '\x00'
					stdio.Sprintf(&buf[0], "%s", d.Obj_name)
					obj.Name = libc.StrDup(&buf[0])
					buf2[0] = '\x00'
					stdio.Sprintf(&buf2[0], "%s", d.Obj_short)
					obj.Short_description = libc.StrDup(&buf2[0])
					buf3[0] = '\x00'
					stdio.Sprintf(&buf3[0], "%s", d.Obj_long)
					obj.Description = libc.StrDup(&buf3[0])
					switch d.Obj_weapon {
					case 1:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_SLASH - TYPE_HIT)
					case 2:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_PIERCE - TYPE_HIT)
					case 3:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_STAB - TYPE_HIT)
					case 4:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_CRUSH - TYPE_HIT)
					case 5:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_BLAST - TYPE_HIT)
					case 6:
						obj.Value[VAL_WEAPON_DAMTYPE] = int(TYPE_POUND - TYPE_HIT)
					}
					obj.Level = 20
				}
				SET_BIT_AR(obj.Extra_flags[:], ITEM_SLOT2)
				d.Obj_editflag = EDIT_NONE
				d.Obj_editval = EDIT_NONE
				d.Character.Rbank -= 30
				d.Character.Rbank = d.Character.Rbank
				d.Rbank -= 30
				userWrite(d, 0, 0, 0, libc.CString("index"))
				obj.Size = get_size(d.Character)
				SET_BIT_AR(obj.Extra_flags[:], ITEM_CUSTOM)
				send_to_imm(libc.CString("Custom Eq: %s has bought: %s."), GET_NAME(d.Character), obj.Short_description)
				customWrite(d.Character, obj)
				log_custom(d, obj)
			} else if libc.StrCaseCmp(arg, libc.CString("n")) == 0 || libc.StrCaseCmp(arg, libc.CString("N")) == 0 {
				write_to_output(d, libc.CString("Canceling purchase at no cost.\r\n"))
				send_to_imm(libc.CString("Custom Eq: %s has canceled their custom eq construction."), GET_NAME(d.Character))
				d.Obj_editval = EDIT_NONE
				d.Connected = CON_PLAYING
			}
		}
	}
}
