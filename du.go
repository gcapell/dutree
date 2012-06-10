// Copyright 2012 Google Inc.  All Rights Reserved.
// Author: ook@google.com (Gary Capell)

package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

var (
	nWorkers = flag.Int("workers", 50, "concurrent workers")
)

type (
	Level int
	NodeID uint64
	Task  struct {
		path  string
		level Level
		reply chan *Result
		store Storage
	}
	Result struct {
		path string
		data int64
		id NodeID
		children []NodeID
	}

	// Worker represented as a channel that accepts tasks
	Worker chan *Task
)

var (
	maxNodeID NodeID
	maxNodeIDLock sync.Mutex
)

func NewID() NodeID {
	maxNodeIDLock.Lock()
	defer maxNodeIDLock.Unlock()
	maxNodeID+=1
	return maxNodeID
}

func NewResult(path string, data int64, children []NodeID) *Result {
	return &Result{path, data, NewID(), children}
}

func manager() Worker {
	queue := make(chan *Task)
	go manageQueues(queue)
	return queue
}

func manageQueues(queue chan *Task) {
	workerAvailable := make(chan (chan *Task))
	tasks := NewTaskHeap()
	var workers WorkerSet
	for i := 0; i < *nWorkers; i++ {
		go worker(workerAvailable, queue)
	}

	for {
		select {
		case w := <-workerAvailable:
			t := tasks.pop()
			if t != nil {
				w <- t
			} else {
				workers.push(w)
			}
		case t := <-queue:
			w := workers.pop()
			if w != nil {
				w <- t
			} else {
				tasks.push(t)
			}
		}
	}

}
func worker(workerAvailable chan chan *Task, tasks chan *Task) {
	myChan := make(chan *Task)
	for {
		workerAvailable <- myChan
		task := <-myChan
		du(task, tasks)
	}
}


func du(task *Task, queue chan *Task) {
	result := du1(task, queue)
	task.store.Store(result)
	task.reply <- result
}

// Create new task to du du childPath, send result to reply
func (t *Task) duChild(childPath string, reply chan *Result) *Task {
	return &Task{
		t.path + "/" + childPath,
		t.level + 1,
		reply,
		t.store}
}

func du1(task *Task, queue chan *Task) *Result {
	fileInfo, err := os.Stat(task.path)
	if err != nil {
		log.Fatal(err) // FIXME
	}
	// Simple file
	if !fileInfo.IsDir() {
		return NewResult(task.path, fileInfo.Size(), nil)
	}
	myChan := make(chan *Result)
	childPaths := listDir(task.path)
	for _, c := range childPaths {
		queue <- task.duChild(c, myChan)
	}
	var myData int64
	children := make([]NodeID, len(childPaths))
	for j := range childPaths {
		result := <-myChan
		myData += result.data
		children[j] = result.id
	}
	return NewResult(task.path, myData, children)
}

func listDir(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	names, err := f.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	return names
}

