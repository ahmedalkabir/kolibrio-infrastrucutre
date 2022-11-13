package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// websocket
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	// the websocket connection
	conn *websocket.Conn
}

func (c *Client) write() {
	log.Println("WRITE CALLED")
	defer func() {
		c.conn.Close()
	}()
	for {

		select {
		case msg := <-webServerToMqtt:
			log.Println("WS -- ", c.conn.RemoteAddr().Network(), " -- ", msg)
			// ws.WriteMessage(websocket.TextMessage, msg.payload)
			msgToSend := struct {
				Topic string `json:"topic"`
				Msg   string `json:"msg"`
			}{
				msg.topic,
				string(msg.payload),
			}

			c.conn.WriteJSON(msgToSend)

		}
	}
}
