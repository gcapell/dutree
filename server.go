package main

import (
	"log"
	"net/http"
)

func tree(w http.ResponseWriter, r *http.Request) {
	log.Printf("\n\nTREE: %v\n", *r)
	http.ServeFile(w, r, "blah.js")
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/tree", tree)
	http.ListenAndServe(":8080", nil)
}
