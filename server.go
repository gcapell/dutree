package main

import (
	"net/http"
	"log"
	"os"
	"io"
	"flag"
)

func tree(w http.ResponseWriter, r *http.Request) {
	log.Printf("\n\nTREE: %v\n", *r)
	fp, err := os.Open("blah.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, fp)
}

func index(w http.ResponseWriter, r *http.Request) {

	req := r.URL.Path[1:]
	fp, err := os.Open(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	written, err := io.Copy(w, fp)
	log.Println(req, written, err)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/tree", tree)
	http.ListenAndServe(":8080", nil)
}
