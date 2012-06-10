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
	Level  int
	NodeID uint64
	Task   struct {
		path  string
		level Level
		reply chan *Result
		store Storage
	}
	Result struct {
		path     string
		data     int64
		id       NodeID
		children []NodeID
	}
)

var (
	maxNodeID     NodeID
	maxNodeIDLock sync.Mutex
)

func NewID() NodeID {
	maxNodeIDLock.Lock()
	defer maxNodeIDLock.Unlock()
	maxNodeID += 1
	return maxNodeID
}

func NewResult(path string, data int64, children []NodeID) *Result {
	return &Result{path, data, NewID(), children}
}

func du(task *Task) {
	result := du1(task)
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

func du1(task *Task) *Result {
	f, err := os.Open(task.path)
	if err != nil {
		log.Println("error", err)
		return NewResult(task.path, 0, nil) // FIXME - indicate error
	}
	defer f.Close()
	TICKETMASTER.getTicket(task.level)
	fileInfo, err := f.Stat()
	TICKETMASTER.returnTicket()
	if err != nil {
		log.Println("error", err)
		return NewResult(task.path, 0, nil) // FIXME - indicate error
	}
	// Simple file
	if !fileInfo.IsDir() {
		return NewResult(task.path, fileInfo.Size(), nil)
	}
	myChan := make(chan *Result)
	childPaths, err := f.Readdirnames(0)
	f.Close()
	if err != nil {
		log.Println("error", err)
		return NewResult(task.path, 0, nil) // FIXME - indicate error
	}
	for _, c := range childPaths {
		go du(task.duChild(c, myChan))
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
