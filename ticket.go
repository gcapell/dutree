package main

import (
	"log"
)

const QUEUE_CHUNK = 1024

var TICKETMASTER = NewTicketMaster(1)

type (
	TicketRequest struct {
		level Level
		reply chan int
	}

	TicketMaster struct {
		tickets uint               // available tickets
		returns chan int           // channel to return tickets
		request chan TicketRequest // channel to request tickets

		maxLevel Level

		// queued requests (by level)
		queue map[Level][]chan int
	}
)

func NewTicketMaster(tickets uint) *TicketMaster {
	t := &TicketMaster{
		tickets,
		make(chan int),
		make(chan TicketRequest),
		0,
		make(map[Level][]chan int),
	}
	go t.run()
	return t
}

func (t *TicketMaster) getTicket(level Level) {
	reply := make(chan int)
	t.request <- TicketRequest{level, reply}
	_ = <-reply
}

func (t *TicketMaster) returnTicket() {
	t.returns <- 1
}

func (t *TicketMaster) run() {
	for {
		select {
		case _ = <-t.returns:
			if t.queueEmpty() {
				t.tickets += 1
			} else {
				t.pop() <- 1
			}
		case r := <-t.request:
			if t.tickets > 0 {
				r.reply <- 1
				t.tickets -= 1
			} else {
				t.push(r.reply, r.level)
			}
		}

	}
}

func (t *TicketMaster) queueEmpty() bool {
	return len(t.queue[t.maxLevel]) == 0
}

func (t *TicketMaster) pop() chan int {
	q := t.queue[t.maxLevel]
	var r chan int
	t.queue[t.maxLevel], r = q[:len(q)-1], q[len(q)-1]

	// Emptied the top queue?
	if len(t.queue[t.maxLevel]) == 0 {
		for {
			t.maxLevel -= 1
			if t.maxLevel < 0 || len(t.queue[t.maxLevel]) > 0 {
				break
			}
		}
	}
	return r
}

func (t *TicketMaster) push(reply chan int, level Level) {
	q, ok := t.queue[level]
	if !ok {
		q = make([]chan int, 0, QUEUE_CHUNK)
	}
	t.queue[level] = append(q, reply)
	if level >= t.maxLevel {
		t.maxLevel = level
	}

	log.Println("queue")
	for k := range t.queue {
		if len(t.queue[k]) > 0 {
			log.Println(k, len(t.queue[k]))
		}
	}

}
