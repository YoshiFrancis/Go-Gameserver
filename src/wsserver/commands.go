package wsserver

import (
	"fmt"
	"strings"
)

func (c *Client) handleCommand(message string) {
	split_msg := strings.Split(message, " ")
	split_msg[0] = strings.ToLower(split_msg[0])
	if split_msg[0] == "echo" {
		c.echo(message)
	} else if split_msg[0] == "join" {
		fmt.Println("User is trying yto join a lobby: ", split_msg[1])
	} else if split_msg[0] == "create" {
		fmt.Println("User is trying to create a lobby of the name", split_msg[1])
	} else if split_msg[0] == "leave" {
		fmt.Println("User is trying to leave his current lobby")
	} else if split_msg[0] == "help" {
		fmt.Println("User is trying to get help!")
	} else if split_msg[0] == "msg" {
		fmt.Println("User is trying to msg the user: ", split_msg[1])
	} else {
		fmt.Println("Invalid command by user!")
	}
}

func (c *Client) echo(message string) {
	c.send <- []byte(message)
}
