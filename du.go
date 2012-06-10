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
		path, name    string
		data     ByteSize
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

func NewResult(path, name string, data ByteSize, children []NodeID) *Result {
	return &Result{path, name, data, NewID(), children}
}

func du(path string, name string, store Storage) *Result {
	result := du1(path, name, store)
	store.Store(result)
	return result
}

func du1(path string, name string, store Storage) *Result {
	f, err := os.Open(path)
	if err != nil {
		log.Println("error", err)
		return NewResult(path, name, 0, nil) // FIXME - indicate error
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		log.Println("error", err)
		return NewResult(path, name, 0, nil) // FIXME - indicate error
	}
	// Simple file
	if !fileInfo.IsDir() {
		return NewResult(path, name, ByteSize(fileInfo.Size()), nil)
	}
	childPaths, err := f.Readdirnames(0)
	f.Close()
	if err != nil {
		log.Println("error", err)
		return NewResult(path, name, 0, nil) // FIXME - indicate error
	}
	var myData ByteSize
	children := make([]NodeID, len(childPaths))
	for j, c := range childPaths {
		result := du(path + "/" + c, c, store)
		myData += result.data
		children[j] = result.id
	}
	return NewResult(path, name, myData, children)
}
