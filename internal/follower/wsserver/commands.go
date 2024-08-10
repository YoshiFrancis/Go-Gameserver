package wsserver

import (
	"fmt"
	"strings"

	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

func (c *Client) handleCommand(message string) (request_msg string) {
	if message[0] != '/' {
		request_msg = messages.FollowerRoomBroadcast(message, c.username)
		return
	}
	split_msg := strings.Split(message, " ")
	split_msg[0] = strings.ToLower(split_msg[0])
	if split_msg[0] == "/echo" {
		c.echo(message)
	} else if split_msg[0] == "/lobby" {
		fmt.Println("Websocket: Attempting to join room: ", split_msg[1])
		// move user to new room
		request_msg = messages.RoomJoinUser(split_msg[1], c.username)

	} else if split_msg[0] == "/create" {
		request_msg = messages.CreateLobby("ws", c.username)
	} else if split_msg[0] == "/leave" {
		request_msg = messages.DisconnectUser(c.username)
	} else if split_msg[0] == "/help" {
		fmt.Println("User is trying to get help!")
	} else if split_msg[0] == "/msg" {
		fmt.Println("User is trying to msg the user: ", split_msg[1])
	} else if split_msg[0] == "/hub" {
		request_msg = messages.RoomJoinUser("hub", c.username)
	} else {
		request_msg = ""
	}
	return
}

func (c *Client) echo(message string) {
	c.send <- []byte(message)
}
