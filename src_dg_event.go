package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

var event_q *queue

func event_init() {
	event_q = queue_init()
}
func event_create(func_ func(event_obj unsafe.Pointer) int, event_obj unsafe.Pointer, when int) *event {
	var new_event *event
	if when < 1 {
		when = 1
	}
	new_event = new(event)
	new_event.Func = func(event_obj unsafe.Pointer) int {
		return func(event_obj unsafe.Pointer) int {
			return func_(event_obj)
		}(event_obj)
	}
	new_event.Event_obj = event_obj
	new_event.Q_el = queue_enq(event_q, unsafe.Pointer(new_event), when+int(pulse))
	return new_event
}
func event_cancel(event *event) {
	if event == nil {
		basic_mud_log(libc.CString("SYSERR:  Attempted to cancel a NULL event"))
		return
	}
	if event.Q_el == nil {
		basic_mud_log(libc.CString("SYSERR:  Attempted to cancel a non-NULL unqueued event, freeing anyway"))
	} else {
		queue_deq(event_q, event.Q_el)
	}
	if event.Event_obj != nil {
		event.Event_obj = nil
	}
	libc.Free(unsafe.Pointer(event))
}
func event_process() {
	var (
		the_event *event
		new_time  int
	)
	for int(pulse) >= queue_key(event_q) {
		if (func() *event {
			the_event = (*event)(queue_head(event_q))
			return the_event
		}()) == nil {
			basic_mud_log(libc.CString("SYSERR: Attempt to get a NULL event"))
			return
		}
		the_event.Q_el = nil
		if (func() int {
			new_time = the_event.Func(the_event.Event_obj)
			return new_time
		}()) > 0 {
			the_event.Q_el = queue_enq(event_q, unsafe.Pointer(the_event), new_time+int(pulse))
		} else {
			libc.Free(unsafe.Pointer(the_event))
		}
	}
}
func event_time(event *event) int {
	var when int
	when = queue_elmt_key(event.Q_el)
	return when - int(pulse)
}
func event_free_all() {
	queue_free(event_q)
}
func event_is_queued(event *event) bool {
	return event != nil && event.Q_el != nil
}
func queue_init() *queue {
	var q *queue
	q = new(queue)
	return q
}
func queue_enq(q *queue, data unsafe.Pointer, key int) *q_element {
	var (
		qe     *q_element
		i      *q_element
		bucket int
	)
	qe = new(q_element)
	qe.Data = data
	qe.Key = key
	bucket = key % NUM_EVENT_QUEUES
	if q.Head[bucket] == nil {
		q.Head[bucket] = qe
		q.Tail[bucket] = qe
	} else {
		for i = q.Tail[bucket]; i != nil; i = i.Prev {
			if i.Key < key {
				if i == q.Tail[bucket] {
					q.Tail[bucket] = qe
				} else {
					qe.Next = i.Next
					i.Next.Prev = qe
				}
				qe.Prev = i
				i.Next = qe
				break
			}
		}
		if i == nil {
			qe.Next = q.Head[bucket]
			q.Head[bucket] = qe
			qe.Next.Prev = qe
		}
	}
	return qe
}
func queue_deq(q *queue, qe *q_element) {
	var i int
	if qe == nil {
		panic("assert failed")
	}
	i = qe.Key % NUM_EVENT_QUEUES
	if qe.Prev == nil {
		q.Head[i] = qe.Next
	} else {
		qe.Prev.Next = qe.Next
	}
	if qe.Next == nil {
		q.Tail[i] = qe.Prev
	} else {
		qe.Next.Prev = qe.Prev
	}
	libc.Free(unsafe.Pointer(qe))
}
func queue_head(q *queue) unsafe.Pointer {
	var (
		dg_data unsafe.Pointer
		i       int
	)
	i = int(pulse % NUM_EVENT_QUEUES)
	if q.Head[i] == nil {
		return nil
	}
	dg_data = q.Head[i].Data
	queue_deq(q, q.Head[i])
	return dg_data
}
func queue_key(q *queue) int {
	var i int
	i = int(pulse % NUM_EVENT_QUEUES)
	if q.Head[i] != nil {
		return q.Head[i].Key
	} else {
		return math.MaxInt32
	}
}
func queue_elmt_key(qe *q_element) int {
	return qe.Key
}
func queue_free(q *queue) {
	var (
		i       int
		qe      *q_element
		next_qe *q_element
		ev      *event
	)
	for i = 0; i < NUM_EVENT_QUEUES; i++ {
		for qe = q.Head[i]; qe != nil; qe = next_qe {
			next_qe = qe.Next
			if (func() *event {
				ev = (*event)(qe.Data)
				return ev
			}()) != nil {
				if ev.Event_obj != nil {
					ev.Event_obj = nil
				}
				libc.Free(unsafe.Pointer(ev))
			}
			libc.Free(unsafe.Pointer(qe))
		}
	}
	libc.Free(unsafe.Pointer(q))
}
