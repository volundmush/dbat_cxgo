package main

import "github.com/gotranspile/cxgo/runtime/libc"

const MIN_MAIL_LEVEL = 3
const STAMP_PRICE = 10
const MAX_MAIL_SIZE = 4096
const BLOCK_SIZE = 256
const HEADER_BLOCK = -1
const LAST_BLOCK = -2
const DELETED_BLOCK = -3

type header_data_type struct {
	Next_block int
	From       int
	To         int
	Mail_time  libc.Time
}
type header_block_type_d struct {
	Block_type  int
	Header_data header_data_type
	Txt         [216]byte
}
type data_block_type_d struct {
	Block_type int
	Txt        [248]byte
}
type position_list_type_d struct {
	Position int
	Next     *position_list_type_d
}
type mail_index_type_d struct {
	Recipient  int
	List_start *position_list_type_d
	Next       *mail_index_type_d
}
