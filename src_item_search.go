package main

func item_is_usable(obj *obj_data) bool {
	return !OBJ_FLAGGED(obj, ITEM_BROKEN) && !OBJ_FLAGGED(obj, ITEM_FORGED)
}
func find_obj_in_list_lambda(head *obj_data, f func(obj *obj_data) bool) *obj_data {
	for head != nil {
		if f(head) {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_vnum(head *obj_data, vn int) *obj_data {
	for head != nil {
		if GET_OBJ_VNUM(head) == vn {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_vnum_good(head *obj_data, vn int) *obj_data {
	for head != nil {
		if GET_OBJ_VNUM(head) == vn && item_is_usable(head) {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_type(head *obj_data, item_type int) *obj_data {
	for head != nil {
		if int(head.Type_flag) == item_type {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_flag(head *obj_data, flag uint32) *obj_data {
	for head != nil {
		if OBJ_FLAGGED(head, flag) {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_name(head *obj_data, name *byte) *obj_data {
	for head != nil {
		if isname(name, head.Name) {
			return head
		}
		head = head.Next_content
	}
	return nil
}
func find_obj_in_list_id(head *obj_data, item_id int) *obj_data {
	for head != nil {
		if int(head.Id) == item_id {
			return head
		}
		head = head.Next_content
	}
	return nil
}
