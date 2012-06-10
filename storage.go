// Copyright 2012 Google Inc.  All Rights Reserved.
// Author: ook@google.com (Gary Capell)

package main

import (
	"fmt"
	"log"
)

type (
	Storage interface {
		Store(*Result)
		Retrieve(NodeID) (Result, error)
	}
	mapStore map[NodeID]Result
)

func (m mapStore) Store(r *Result) {
	log.Println("store", r)
	m[r.id] = *r
}

func (m mapStore) Retrieve(n NodeID) (Result, error) {
	r, ok := m[n]
	if !ok {
		return r, fmt.Errorf("%v not in store", n)
	}
	return r, nil
}
