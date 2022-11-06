package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var webSocketConnection *websocket.Conn

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var err error
		webSocketConnection, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
	})
	go func() {
		for {
			if webSocketConnection != nil {
				_, rdr, err := webSocketConnection.NextReader()
				if err != nil {
					log.Println(err)
				}
				data, err := io.ReadAll(rdr)
				if err != nil {
					log.Println(err)
				}
				wtr, err := webSocketConnection.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Println(err)
				}
				wtr.Write(data)
				err = wtr.Close()
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
