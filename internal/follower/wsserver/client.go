package wsserver

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	Send     chan []byte
	server   *WSServer
}

func NewClient(username string, conn *websocket.Conn, server *WSServer) {
	client := &Client{
		username: username,
		conn:     conn,
		Send:     make(chan []byte),
		server:   server,
	}

	server.register <- client
	go client.read()
}

func (c *Client) read() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			break
		}

		fmt.Println(c.username + " received a message: " + string(message))
		c.server.TCPto <- []byte(c.handleCommand(string(message)))
	}
}

func (c *Client) write() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for message := range c.Send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(message)
		message_count := len(c.Send)
		for i := 0; i < message_count; i++ {
			w.Write([]byte("\n"))
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}

	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
