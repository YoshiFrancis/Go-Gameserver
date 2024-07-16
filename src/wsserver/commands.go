package wsserver

import (
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) handleCommand(message string) {
	split_msg := strings.Split(message, " ")
	split_msg[0] = strings.ToLower(split_msg[0])
	if split_msg[0] == "echo" {
		c.echo(message)
	} else if split_msg[0] == "join" {
		fmt.Println("Attempting to join room")
		roomId, err := strconv.Atoi(split_msg[1])
		if err != nil {
			fmt.Println("Given an invalid room id")
		}

		new_room, ok := getRoom(c.room.server, roomId)
		if !ok {
			fmt.Println("Room id does not exist")
			return
		}
		c.switchRoom(new_room)

	} else if split_msg[0] == "create" {
		roomTitle := split_msg[1]
		new_room := NewRoom(roomTitle, c.room, c.room.server)
		go new_room.run()
		c.switchRoom(new_room)
		fmt.Println("Created room: ", roomTitle)
	} else if split_msg[0] == "leave" {
		c.switchRoom(c.room.parentRoom)
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
