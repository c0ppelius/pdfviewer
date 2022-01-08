package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/c0ppelius/pdfviewer/src"
	"github.com/gorilla/websocket"
)

const (
	port = ":8888"
)

var assets = src.Assets

var upgrader = websocket.Upgrader{}

type Message struct {
	Action string `json:"action"`
	Page   string `json:"page"`
}

type pageStruct struct {
	Page int
}

var cmd Message

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

func main() {

	trigger := false

	cmd = Message{}

	path := filePath()

	http.HandleFunc("/pdf/xyz.pdf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})

	url := "http://localhost:8888/web/viewer.html?file=http%3A%2F%2Flocalhost%3A8888%2Fpdf%2Fxyz.pdf"

	fmt.Println("Pdf at: ", url)

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		trigger = true
		cmd.Action = "reload"
		w.Write(nil)
	})

	// hit with `curl -h "Content-Type: application/json" -d '{"Page":1}' http://localhost:8888/reload`
	http.HandleFunc("/forward", func(w http.ResponseWriter, r *http.Request) {
		trigger = true
		cmd.Action = "forward"
		decoder := json.NewDecoder(r.Body)
		var p pageStruct
		_ = decoder.Decode(&p)
		cmd.Page = url + "#page=" + strconv.Itoa(p.Page)
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
			if trigger {
				msg, err := json.Marshal(cmd)
				if err != nil {
					log.Println("json marshall failed:", err)
				}
				err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Println("write failed:", err)
				}
				trigger = false
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
