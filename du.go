// Copyright 2012 Google Inc.  All Rights Reserved.
// Author: ook@google.com (Gary Capell)

package main

import (
	"flag"
	"log"
	"os"
	"fmt"
)

var (
	nWorkers = flag.Int("workers", 50, "concurrent workers")
)

type (
	Level int
	Task  struct {
		path  string
		level Level
		reply chan *Result
	}
	Result struct {
		path string
		data int64
		children []*Result
	}

	// Worker represented as a channel that accepts tasks
	Worker chan *Task
)

func (r *Result) show(level int) {
	fmt.Printf("%*s %s %d\n", level*4, " ", r.path, r.data)
	if r.children == nil {
		return
	}
	for _,c := range r.children {
		c.show(level+1)
	}
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
	log.Println(result)
	// store(result)
	task.reply <- result
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

func du1(task *Task, queue chan *Task) *Result {
	fileInfo, err := os.Stat(task.path)
	if err != nil {
		log.Fatal(err) // FIXME
	}
	// Simple file
	if !fileInfo.IsDir() {
		return &Result{task.path, fileInfo.Size(), nil}
	}
	myChan := make(chan *Result)
	childPaths := listDir(task.path)
	for _, c := range childPaths {
		queue <- &Task{task.path + "/" + c, task.level + 1, myChan}
	}
	var myData int64
	children := make([]*Result, len(childPaths))
	for j := range childPaths {
		result := <-myChan
		myData += result.data
		children[j] = result
	}
	return &Result{task.path, myData, children}
}
