package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump sekarang akan membaca pesan JSON dari klien
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}
		message, err := json.Marshal(msg) // Mengubah pesan ke dalam format JSON
		if err != nil {
			log.Println("json.Marshal:", err)
			break
		}
		c.hub.broadcast <- message
	}
}

// writePump sekarang akan menulis pesan JSON ke klien
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		var msg Message
		err := json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("json.Unmarshal:", err)
			return
		}

		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println("write:", err)
			return
		}
	}
}
