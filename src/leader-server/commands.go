package leaderserver

import (
	"fmt"
	"strconv"
	"strings"
)

func (l *Leader) handleArgs(flag byte, args []string) (res string) {
	args[0] = strings.ToLower(args[0])
	res = ""
	if flag == '-' { // server
		switch args[0] {
		case "create": // creating server
			lobby := NewLobby(l.idGen())
			l.mutex.Lock()
			defer l.mutex.Unlock()
			l.lobbies[lobby.lobbyId] = lobby
			go lobby.run()
		// messages.RoomCreation
		case "shutdown": // shutting server down
			break
		case "disc": // disconnecting user
			username := args[2]
			serverId, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("given invalid serverid")
			}
			user := l.Users[username]
			if user.serverId != serverId {
				fmt.Println("User does not belong to that server!")
				return
			}
			l.disconnectUser(username)
		case "join": // user is joining
			username := args[2]
			serverId, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("given invalid serverid")
			}
			user := NewUser(username, serverId, l.hub.hubId, l.WSServer.Clients[username])
			l.Users[username] = user
			l.hub.register <- user
		default:
			fmt.Println("Given an invalid server command")
			return
		}
	} else if flag == '+' { // hub
		hubId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("given an invalid room id")
			return
		}
		if hubId != 1 {
			fmt.Println("given invalid hub id")
			return
		}

		switch args[0] {
		case "broadcast":
			l.hub.broadcast <- []byte(args[2])
		default:
			fmt.Println("Given invalid hub command")
			return
		}
	} else if flag == '/' { // lobby
		lobbyId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("given an invalid room id")
			return
		}
		lobby, ok := l.lobbies[lobbyId]
		if !ok {
			fmt.Println("given invalid hub id")
			return
		}
		switch args[0] {
		case "broadcast":
			lobby.broadcast <- []byte(args[2])
		default:
			fmt.Println("Given invalid lobby id")
		}
	} else if flag == '*' { // user
		username := args[1]
		user, ok := l.Users[username]
		if !ok {
			fmt.Println("given invalid username")
			return
		}

		switch args[0] {
		case "pm":
			user.send([]byte(args[2]))
		default:
			fmt.Println("Given an invalid user command")
			return
		}
	}

	return ""
}
