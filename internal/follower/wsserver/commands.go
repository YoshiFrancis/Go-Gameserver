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

	switch split_msg[0] {
	case "/echo":
		c.echo(message)
	case "/join":
		if len(split_msg) == 1 {
			return
		}
		request_msg = messages.RoomJoinUser(split_msg[1], c.username)
	case "/lobby":
		if len(split_msg) == 1 {
			return
		}
		request_msg = messages.CreateLobby(split_msg[1], c.username)
	case "/leave":
		request_msg = messages.DisconnectUser(c.username)
	case "/help":
		fmt.Println("User is trying to get help!")
	case "/msg":
		fmt.Println("User is trying to msg the user: ", split_msg[1])
	case "/hub":
		request_msg = messages.RoomJoinUser("hub", c.username)
	case "/app":
		if len(split_msg) == 1 {
			return
		}
		request_msg = messages.FollowerAppRequest(split_msg[1], c.username)
	case "/app-start":
		if len(split_msg) == 1 {
			return
		}
		request_msg = messages.FollowerAppStart(split_msg[1], c.username)
	default:
		request_msg = ""
	}
	return
}

func (ws *WSServer) handleLeaderRequest(req LeaderRequest) {
	switch req.command {
	case "broadcast":
		for _, username := range req.usernames {
			c, ok := ws.Clients.Get(username)
			if !ok {
				return
			}
			c.send <- []byte(req.arg)
		}
	}
}

func (c *Client) echo(message string) {
	c.send <- []byte(message)
}
