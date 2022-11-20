package main

import "unsafe"

const PULSE_DG_EVENT = 1
const NUM_EVENT_QUEUES = 10

type EventFunc func(event_obj unsafe.Pointer) int
type event struct {
	Func      EventFunc
	Event_obj unsafe.Pointer
	Q_el      *q_element
}
type queue struct {
	Head [10]*q_element
	Tail [10]*q_element
}
type q_element struct {
	Data unsafe.Pointer
	Key  int
	Prev *q_element
	Next *q_element
}
