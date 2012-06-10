package main

import (
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
	http.HandleFunc("/", index)
	http.HandleFunc("/tree", tree)
	http.ListenAndServe(":8080", nil)
}
