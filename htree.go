package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const HTREE_NODE_BITS = 4
const HTREE_NODE_SUBS = 16
const HTREE_NODE_MASK = 15

type htree_node struct {
	Content int64
	Parent  *htree_node
	Subs    [16]*htree_node
}

var HTREE_NULL *htree_node = nil
var htree_total_nodes int = 0
var htree_depth_used int = 0

func htree_shutdown() {
	libc.Free(unsafe.Pointer(HTREE_NULL))
	HTREE_NULL = nil
}
func htree_init() *htree_node {
	var (
		newnode *htree_node
		i       int
	)
	if HTREE_NULL == nil {
		htree_total_nodes++
		HTREE_NULL = new(htree_node)
		for i = 0; i < (int(1 << HTREE_NODE_BITS)); i++ {
			HTREE_NULL.Subs[i] = HTREE_NULL
		}
		HTREE_NULL.Content = -1
		HTREE_NULL.Parent = nil
	}
	if htree_depth_used == 0 {
		htree_depth_used = 1
	}
	htree_total_nodes++
	newnode = new(htree_node)
	libc.MemCpy(unsafe.Pointer(&newnode.Subs[0]), unsafe.Pointer(&HTREE_NULL.Subs[0]), (int(1<<HTREE_NODE_BITS))*int(unsafe.Sizeof((*htree_node)(nil))))
	newnode.Content = -1
	newnode.Parent = HTREE_NULL
	return newnode
}
func htree_free(root *htree_node) {
	var i int
	if root == nil || root == HTREE_NULL {
		return
	}
	for i = 0; i < (int(1 << HTREE_NODE_BITS)); i++ {
		htree_free(root.Subs[i])
	}
	libc.Free(unsafe.Pointer(root))
}
func htree_add(root *htree_node, htindex int64, content int64) {
	var (
		tmp   *htree_node
		i     int
		depth int
	)
	if root == nil {
		return
	}
	tmp = root
	depth = 0
	for htindex != 0 {
		depth++
		i = int(htindex & int64((int(1<<HTREE_NODE_BITS))-1))
		htindex >>= HTREE_NODE_BITS
		if tmp.Subs[i] == HTREE_NULL {
			htree_total_nodes++
			tmp.Subs[i] = new(htree_node)
			libc.MemCpy(unsafe.Pointer(&tmp.Subs[i].Subs[0]), unsafe.Pointer(&HTREE_NULL.Subs[0]), (int(1<<HTREE_NODE_BITS))*int(unsafe.Sizeof((*htree_node)(nil))))
			tmp.Subs[i].Content = -1
			tmp.Subs[i].Parent = HTREE_NULL
		}
		tmp = tmp.Subs[i]
	}
	if tmp == HTREE_NULL {
		return
	}
	if depth > htree_depth_used {
		htree_depth_used = depth
	}
	tmp.Content = content
}
func htree_find_node(root *htree_node, htindex int64) *htree_node {
	var (
		tmp *htree_node
		i   int
	)
	tmp = root
	for htindex != 0 {
		i = int(htindex & int64((int(1<<HTREE_NODE_BITS))-1))
		htindex >>= HTREE_NODE_BITS
		tmp = tmp.Subs[i]
	}
	return tmp
}
func htree_del(root *htree_node, htindex int64) {
	var tmp *htree_node
	tmp = htree_find_node(root, htindex)
	tmp.Content = -1
}
func htree_find(root *htree_node, htindex int64) int64 {
	var tmp *htree_node
	tmp = htree_find_node(root, htindex)
	return tmp.Content
}
func real_room_old(vnum room_vnum) room_rnum {
	var (
		bot room_rnum
		top room_rnum
		mid room_rnum
	)
	bot = 0
	top = top_of_world
	for {
		mid = (bot + top) / 2
		if ((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(mid)))).Number == vnum {
			return mid
		}
		if bot >= top {
			return -1
		}
		if ((*room_data)(unsafe.Add(unsafe.Pointer(world), unsafe.Sizeof(room_data{})*uintptr(mid)))).Number > vnum {
			top = mid - 1
		} else {
			bot = mid + 1
		}
	}
}
func htree_test() {
	basic_mud_log(libc.CString("htree stats (global): %d nodes, %lld bytes (depth %d/%lld used/possible)"), htree_total_nodes, htree_total_nodes*int(unsafe.Sizeof(htree_node{})), htree_depth_used, ((unsafe.Sizeof(int64(0))*8)/HTREE_NODE_BITS)+1)
}
