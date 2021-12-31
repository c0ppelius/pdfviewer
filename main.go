package main

import (
	"flag"
	"log"
	"net/http"
)

const (
	port = ":8888"
	src  = "/Users/matt/GitHub/pdfviewer/src"
)

var file string

func init() {
	flag.StringVar(&file, "file", "", "the pdf file")
	flag.Parse()
}

func main() {

	http.HandleFunc("/pdf/xyz.pdf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./"+file)
	})

	fs := http.FileServer(http.Dir(src))
	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
