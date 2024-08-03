package wsserver

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	send     chan []byte
	ws       *WSServer
}

type Message struct {
	Message string `json:"message"`
}

func NewClient(username string, conn *websocket.Conn, ws *WSServer) *Client {
	client := &Client{
		username: username,
		conn:     conn,
		send:     make(chan []byte),
		ws:       ws,
	}

	client.ws.register <- client
	go client.read()
	go client.write()
	return client
}

func (c *Client) read() {
	defer func() {
		c.ws.unregister <- c
		c.conn.Close()
	}()

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			break
		}
		var message Message
		json.Unmarshal(jsonMessage, &message)
		handled := []byte(c.handleCommand(string(message.Message)))
		c.ws.TCPto <- handled
	}
}

func (c *Client) write() {
	defer func() {
		c.ws.unregister <- c
		c.conn.Close()
	}()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(message)
		message_count := len(c.send)
		for i := 0; i < message_count; i++ {
			w.Write([]byte("\n"))
			w.Write(<-c.send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}

	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
