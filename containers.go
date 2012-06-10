// Copyright 2012 Google Inc.  All Rights Reserved.
// Author: ook@google.com (Gary Capell)

package main

import (
	"log"
)

type (
	TaskHeap struct {
		// Slices of all tasks of a given level
		task map[Level]([]*Task)
		// largest non-empty level
		top Level
	}
	WorkerSet []Worker
)

func NewTaskHeap() TaskHeap {
	return TaskHeap{make(map[Level][]*Task), -1}
}

func (th *TaskHeap) push(t *Task) {
	if _, ok := th.task[t.level]; !ok {
		th.task[t.level] = make([]*Task, 1, 1024)
		th.task[t.level][0] = t
	} else {
		th.task[t.level] = append(th.task[t.level], t)
	}
	if t.level > th.top {
		th.top = t.level
		log.Println("New top is", th.top)
	}
}

// Return any task from the highest non-empty level in heap
func (th *TaskHeap) pop() *Task {
	if th.top < 0 {
		return nil
	}
	ts := th.task[th.top] // top slice
	if len(ts) == 0 {
		log.Fatalln(*th, ts, th.top)
	}
	var t *Task
	t, th.task[th.top] = ts[len(ts)-1], ts[:len(ts)-1]
	if len(th.task[th.top]) == 0 {
		// Emptied top level, find next non-empty level
		for {
			th.top -= 1
			if th.top < 0 {
				break
			}
			if ts, ok := th.task[th.top]; ok {
				if len(ts) > 0 {
					break
				}
			}
		}
		log.Println("New top is", th.top)
	}
	return t
}

func (ws *WorkerSet) push(w Worker) {
	*ws = append(*ws, w)
}

func (ws *WorkerSet) pop() Worker {
	size := len(*ws)
	if size == 0 {
		return nil
	}
	w := (*ws)[size-1]
	*ws = (*ws)[:size-1]
	return w
}
