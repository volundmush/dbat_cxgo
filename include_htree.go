package main

const HTREE_NODE_BITS = 4
const HTREE_NODE_SUBS = 16
const HTREE_NODE_MASK = 15

type htree_node struct {
	Content int
	Parent  *htree_node
	Subs    [16]*htree_node
}
