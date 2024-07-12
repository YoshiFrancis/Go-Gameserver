package wsserver

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	room     *Room
	send     chan []byte
	prompt   string
}

func (c *Client) read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			break
		}
		// need to use messages.go and place message in a struct
		c.room.messages <- string(message)
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		// need to use messages.go and place message in a struct

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
	// hub closed channel
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})

}
