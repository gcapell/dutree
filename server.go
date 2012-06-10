package main

import (
	"flag"
	"log"
	"net/http"
	"encoding/json"
)

type (
	Attr struct {Id string `json:"id"`}
	Node struct {
		Data string `json:"data"`
		State string `json:"state"`
		Attr Attr `json:"attr"`
	}
)

func tree(w http.ResponseWriter, r *http.Request) {
	log.Printf("TREE: %v\n", r.URL)

	n := []Node {
		{"c", "d",  Attr{"3"}},
		{"a", "closed",  Attr{"1"}},
	}

	b, err := json.Marshal(n)
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
	manager() <- &Task{flag.Arg(0), 0, reply}
	data := <-reply
	log.Printf("data: %v", data)
	data.show(0)

	http.HandleFunc("/", index)
	http.HandleFunc("/tree", tree)
	addr := ":8080"
	log.Printf("Listening on", addr)
	http.ListenAndServe(addr, nil)
}
