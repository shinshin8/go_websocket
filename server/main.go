package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

type Msg struct {
	IsNew   bool
	Message string
}

type Upgrader struct {
	Ugrader websocket.Upgrader
}

func (u *Upgrader) wsocket(w http.ResponseWriter, r *http.Request) {
	u.Ugrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := u.Ugrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("connected")

	var msg Msg
	done := make(chan bool)

	go func() {
		for i := 15; i >= 0; i-- {
			time.Sleep(time.Second * 1)

			if i == 0 {
				done <- true
				return
			}

			_, p, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			json.Unmarshal(p, &msg)

			if msg.IsNew {
				i += 10
			}

			log.Println(msg)

		}
	}()

	for {
		select {
		case <-done:
			log.Println("timeout")
			os.Exit(0)
		}
	}
}

func main() {
	up := &Upgrader{Ugrader: websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}}
	http.HandleFunc("/ws", up.wsocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
