package leaderserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

func (l *Leader) handleArgs(flag byte, args []string) (res string) { // response is only used if the message is from the websocket server
	args[0] = strings.ToLower(args[0])
	res = ""
	if flag == '-' { // server
		switch args[0] {
		case "disc": // disconnecting user
			username := args[1]
			l.disconnectUser(username)
			res = messages.ServerDisconnectUser(username)
		case "join": // user is joining
			fmt.Println("User is joining leader!")
			username := args[2]
			serverId, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("given invalid serverid")
			}

			l.mutex.Lock()
			defer l.mutex.Unlock()

			user := NewUser(username, serverId, l.hub.hubId, l.WSServer.Clients[username])
			l.Users[username] = user
			l.hub.register <- user
			res = messages.ServerJoinUser(username, serverId)
		default:
			fmt.Println("Given an invalid server command")
			return
		}
	} else if flag == '+' { // hub
		switch args[0] {
		case "broadcast":
			res = messages.HubBroadcast(args[2], l.hub.hubId, args[3])
			username := args[1]
			broadcast_msg := args[3]
			if l.Users[username].serverId != l.TCPServer.ServerId() {
				l.hub.broadcast <- []byte(broadcast_msg)
			}
		case "join":
			username := args[2]
			userRoomId := l.Users[username].roomId
			if userRoomId != l.hub.hubId {
				l.lobbies[userRoomId].unregister <- l.Users[username]
				res = messages.HubJoinUser(username, l.hub.hubId)
				l.hub.register <- l.Users[username]
			}
		case "lobby":
			l.mutex.Lock()
			defer l.mutex.Unlock()
			roomId, _ := strconv.Atoi(args[1])
			if args[2] == "ws" {
				roomId = l.idGen()
				res = messages.HubCreateLobby("tcp", roomId)
			} else {
				l.idGen = idGenerator(roomId)
			}
			lobby := NewLobby(roomId)
			l.lobbies[lobby.lobbyId] = lobby
			go lobby.run()
			fmt.Println("New lobby id:", lobby.lobbyId)
		default:
			fmt.Println("Given invalid hub command")
			return
		}
	} else if flag == '/' { // lobby
		l.mutex.Lock()
		defer l.mutex.Unlock()
		lobbyId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("given an invalid room id")
			return
		}
		_, ok := l.lobbies[lobbyId]
		if !ok {
			fmt.Println("given invalid lobby id")
			return
		}
		switch args[0] {
		case "broadcast":
			username := args[2]
			broadcast_msg := args[3]
			l.lobbies[lobbyId].broadcast <- []byte(username + broadcast_msg) // need a better api for this
			res = messages.LobbyBroadcast(username, l.Users[username].roomId, args[3])
		case "join":
			lobby, ok := l.lobbies[lobbyId]
			if !ok {
				fmt.Println("given invalid room id")
				return
			}
			username := args[2]
			userRoomId := l.Users[username].roomId
			if userRoomId == l.hub.hubId {
				l.hub.unregister <- l.Users[username]
			} else if lobbyId != userRoomId {
				lobby, ok := l.lobbies[userRoomId]
				if ok {
					lobby.unregister <- l.Users[username]
				}
			} else {
				return
			}
			res = messages.LobbyJoinUser(username, lobby.lobbyId)
			lobby.register <- l.Users[username]
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

	return
}
