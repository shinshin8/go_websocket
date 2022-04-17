package main

import (
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type Msg struct {
	IsNew   bool
	Message string
}

var (
	origin = "http://localhost:8080"
	url    = "ws://localhost:8080/ws"
)

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatalln(err)
	}

	ticker := time.NewTicker(time.Second * 10)
	reNew := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				reNew <- true
			}
		}
	}()

	for {
		select {
		case <-reNew:
			sendMsg(ws, true, "message")
			time.Sleep(time.Second * 1)
		default:
			sendMsg(ws, false, "message")
			time.Sleep(time.Second * 1)
		}
	}
}

func sendMsg(ws *websocket.Conn, flg bool, msg string) {
	var snd = Msg{flg, msg}
	err := websocket.JSON.Send(ws, snd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(snd)
}
