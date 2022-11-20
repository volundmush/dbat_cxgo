package main

import "github.com/gotranspile/cxgo/runtime/libc"

const MAX_MESSAGE_LENGTH = 4096
const BOARD_MAGIC = 0xFFFFF
const CURRENT_BOARD_VER = 2

type board_msg struct {
	Poster    int
	Timestamp libc.Time
	Subject   *byte
	Data      *byte
	Next      *board_msg
	Prev      *board_msg
	Name      *byte
}
type board_memory struct {
	Timestamp int
	Reader    int
	Next      *board_memory
	Name      *byte
}
type board_info struct {
	Read_lvl     int
	Write_lvl    int
	Remove_lvl   int
	Num_messages int
	Vnum         int
	Next         *board_info
	Messages     *board_msg
	Version      int
	Memory       [301]*board_memory
}
