package wsserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

func (c *Client) handleCommand(message string) (request_msg string) {
	split_msg := strings.Split(message, " ")
	split_msg[0] = strings.ToLower(split_msg[0])
	if split_msg[0] == "echo" {
		c.echo(message)
	} else if split_msg[0] == "join" {
		fmt.Println("Websocket: Attempting to join room")
		roomId, err := strconv.Atoi(split_msg[1])
		if err != nil {
			fmt.Println("Given an invalid room id")
		}
		// move user to new room
		request_msg = messages.LobbyJoinUser(c.username, roomId)

	} else if split_msg[0] == "create" {

		// roomTitle := split_msg[1]
		request_msg = messages.ServerCreateLobby("ws", -1)
	} else if split_msg[0] == "leave" {
		request_msg = messages.ServerDisconnectUser(c.username)
	} else if split_msg[0] == "help" {
		fmt.Println("User is trying to get help!")
	} else if split_msg[0] == "msg" {
		fmt.Println("User is trying to msg the user: ", split_msg[1])
	} else {
		if len(split_msg) == 1 {
			request_msg = messages.HubBroadcast(c.username, 1, split_msg[0])
			fmt.Println("Broadcasting to hub", request_msg)
		}
	}
	return
}

func (c *Client) echo(message string) {
	c.Send <- []byte(message)
}
