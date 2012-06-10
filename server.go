package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type (
	Attr struct {
		Id NodeID `json:"id"`
	}
	Node struct {
		Data  string `json:"data"`
		State string `json:"state"`
		Attr  Attr   `json:"attr"`
	}
)

func tree(w http.ResponseWriter, r *http.Request, store Storage, root NodeID) {
	log.Printf("TREE: %v => %s\n", r.URL, r.FormValue("id"))

	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	var nodeID NodeID
	if err != nil {
		nodeID = root
	} else {
		nodeID = NodeID(id)
	}
	node, err := store.Retrieve(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	children := make([]Node, len(node.children))
	for j, c := range node.children {
		cval, err := store.Retrieve(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var state string
		if len(cval.children) != 0 {
			state = "closed"
		} else {
			state = "x"
		}
		children[j] = Node{
			fmt.Sprintf("%s %d", cval.path, cval.data),
			state, Attr{cval.id}}
	}

	b, err := json.Marshal(children)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("returning: %v\n", string(b))
	w.Write(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	flag.Parse()
	reply := make(chan *Result)
	var store mapStore = make(map[NodeID]Result)
	manager() <- &Task{flag.Arg(0), 0, reply, store}
	data := <-reply
	log.Printf("data: %v", data)
	log.Printf("store: %v", store)

	http.HandleFunc("/", index)
	http.HandleFunc("/tree", func(w http.ResponseWriter, r *http.Request) {
		tree(w, r, store, data.id)
	})
	addr := ":8080"
	log.Printf("Listening on", addr)
	http.ListenAndServe(addr, nil)
}
