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
		c.room.server.leaving <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			break
		}

		c.handleMessage(string(message))
		// need to use messages.go and place message in a struct
	}
}

func (c *Client) write() {
	defer func() {
		c.room.unregister <- c
		c.room.server.leaving <- c
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

func (c *Client) handleMessage(message string) {
	if c.prompt != NONE {
		c.handlePrompt(message)
	} else if message[0] == '/' {
		c.handleCommand(message[1:])
	} else { // broadcast it
		c.room.messages <- message
	}
}

func (c *Client) switchRoom(r *Room) {
	c.room.unregister <- c
	r.register <- c
	fmt.Println("client joined the room: ", r.roomId)
}
