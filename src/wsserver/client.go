package wsserver

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	send     chan []byte
	server   *WSServer
}

func (c *Client) read() {
	defer func() {
		c.server.leaving <- c.username
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			break
		}

		fmt.Println(c.username + " received a message: " + string(message))

		c.server.requests <- message
	}
}

func (c *Client) write() {
	defer func() {
		c.server.leaving <- c.username
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
