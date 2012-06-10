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
	if err != nil {
		leaf(w, store, root)
	} else {
		leaves(w, store, NodeID(id))
	}
}

func leaf(w http.ResponseWriter, store Storage, n NodeID) {
	val, err := store.Retrieve(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	node := Node {fmt.Sprintf("%s %d", val.path, val.data),"closed", Attr{val.id}}

	b, err := json.Marshal(node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func leaves(w http.ResponseWriter, store Storage, n NodeID) {
	val, err := store.Retrieve(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	children := make([]Node, len(val.children))
	for j, c := range val.children {
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
	w.Write(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	flag.Parse()
	var store mapStore = make(map[NodeID]Result)
	data := du(flag.Arg(0), store)

	http.HandleFunc("/", index)
	http.HandleFunc("/tree", func(w http.ResponseWriter, r *http.Request) {
		tree(w, r, store, data.id)
	})
	addr := ":8080"
	log.Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}
