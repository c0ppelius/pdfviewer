package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	port = ":8888"
	src  = "/Users/matt/GitHub/pdfviewer/src"
)

func badFile() {
	fmt.Println("Enter the pdf file name.")
}

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		badFile()
		return
	}

	file := args[0]

	if !strings.Contains(file, "pdf") {
		badFile()
		return
	} else if file[len(file)-3:] != "pdf" {
		badFile()
		return
	}

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
