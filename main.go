package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/c0ppelius/pdfviewer/src"
	"github.com/gorilla/websocket"
)

const (
	port = ":8888"
)

var upgrader = websocket.Upgrader{}

var cmd string

func filePath() (path string) {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Enter the pdf file name.")
		log.Fatal("No file entered")
	}
	arg := args[0]

	if !strings.Contains(arg, "pdf") {
		fmt.Println("This doesn't look like a pdf file.")
		log.Fatal("Not a pdf")
	}

	if arg[len(arg)-3:] != "pdf" {
		fmt.Println("This doesn't look like a pdf file.")
		log.Fatal("Not a pdf")
	}

	_, err := os.Open(arg)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Can't find that file.")
		log.Fatal(err)
	}

	path, _ = filepath.Abs(arg)
	return
}

var assets = src.Assets

func main() {

	path := filePath()

	http.HandleFunc("/pdf/xyz.pdf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})

	url := "http://localhost:8888/web/viewer.html?file=http%3A%2F%2Flocalhost%3A8888%2Fpdf%2Fxyz.pdf"

	fmt.Println("Pdf at: ", url)

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		cmd = "reload"
		w.Write(nil)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		for {
			if cmd != "" {
				err := conn.WriteMessage(websocket.TextMessage, []byte(cmd))
				if err != nil {
					log.Println("write failed:", err)
				}
				cmd = ""
			}
		}
	})

	fs := http.FileServer(http.FS(assets))
	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
