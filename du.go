// Copyright 2012 Google Inc.  All Rights Reserved.
// Author: ook@google.com (Gary Capell)

package main

import (
	"log"
	"os"
	"sync"
)

type (
	Level  int
	NodeID uint64
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

func du(path string, store Storage) *Result {
	result := du1(path, store)
	store.Store(result)
	return result
}

func du1(path string, store Storage) *Result {
	f, err := os.Open(path)
	if err != nil {
		log.Println("error", err)
		return NewResult(path, 0, nil) // FIXME - indicate error
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		log.Println("error", err)
		return NewResult(path, 0, nil) // FIXME - indicate error
	}
	// Simple file
	if !fileInfo.IsDir() {
		return NewResult(path, fileInfo.Size(), nil)
	}
	childPaths, err := f.Readdirnames(0)
	f.Close()
	if err != nil {
		log.Println("error", err)
		return NewResult(path, 0, nil) // FIXME - indicate error
	}
	var myData int64
	children := make([]NodeID, len(childPaths))
	for j, c := range childPaths {
		result := du(path+"/"+c, store)
		myData += result.data
		children[j] = result.id
	}
	return NewResult(path, myData, children)
}
