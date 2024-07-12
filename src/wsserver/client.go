package wsserver

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	room     *Room
	output   chan []byte
}

func (client *Client) read() {

}

func (client *Client) write() {

}
