package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var webSocketConnections = make(map[string]*websocket.Conn)

func readWritePump(username string) {
	for {
		webSocketConnection, ok := webSocketConnections[username]
		if ok {
			_, rdr, err := webSocketConnection.NextReader()
			if err != nil {
				log.Println(err)
			}
			if rdr != nil {
				_, err := io.ReadAll(rdr)
				if err != nil {
					log.Println(err)
				}
			}
			wtr, err := webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
			}
			if wtr != nil {
				err = wtr.Close()
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(250 * time.Millisecond)
		} else {
			return
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		var username string
		for _, cookie := range cookies {
			if cookie.Name == "username" {
				username = cookie.Value
			}
		}
		if username == "" {
			return
		}
		webSocketConnection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		webSocketConnections[username] = webSocketConnection
		go readWritePump("bob")
	})

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
