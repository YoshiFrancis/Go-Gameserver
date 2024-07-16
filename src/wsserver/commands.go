package wsserver

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *Server) handleCommand(args []string) {

}

func (r *Room) handleCommand(args []string) {

}

func (c *Client) handleCommand(args []string) {

	args[0] = strings.ToLower(args[0])
	if args[0] == "echo" {
		c.echo(args[1])
	} else if args[0] == "join" {
		fmt.Println("Attempting to join room")
		roomId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Given an invalid room id")
		}

		new_room, ok := getRoom(c.room.server, roomId)
		if !ok {
			fmt.Println("Room id does not exist")
			return
		}
		c.switchRoom(new_room)

	} else if args[0] == "create" {
		roomTitle := args[1]
		new_room := NewRoom(roomTitle, c.room, c.room.server)
		go new_room.run()
		c.switchRoom(new_room)
		fmt.Println("Created room: ", roomTitle)
	} else if args[0] == "leave" {
		c.switchRoom(c.room.parentRoom)
	} else if args[0] == "help" {
		fmt.Println("User is trying to get help!")
	} else if args[0] == "msg" {
		fmt.Println("User is trying to msg the user: ", args[1])
	} else {
		fmt.Println("Invalid command by user!")
	}
}

func (c *Client) echo(message string) {
	c.send <- []byte(message)
}
