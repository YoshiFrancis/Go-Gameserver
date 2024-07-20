package leaderserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

func (l *Leader) handleArgs(flag byte, args []string) (res string) {
	args[0] = strings.ToLower(args[0])
	res = ""
	if flag == '-' { // server
		switch args[0] {
		case "creation":
			// the server that was originally connected now must broadcast to all other servers rhat there is a neew server
			// i have to come up with new key word to signal that the new server has already been accepted by one of the nodes in the group already
			// the new servers will connect with the already connected node
			// this original node that accepted has the send all data about the servers to tje new node
		case "serverid":
		case "shutdown": // shutting server down
			break
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
			user := NewUser(username, serverId, l.hub.hubId, l.WSServer.Clients[username])
			l.Users[username] = user
			l.hub.register <- user
			res = messages.ServerJoinUser(username, serverId)
			fmt.Println("User joining message:", res)
		default:
			fmt.Println("Given an invalid server command")
			return
		}
	} else if flag == '+' { // hub
		fmt.Println(args)
		switch args[0] {
		case "broadcast":
			res = messages.HubBroadcast(args[2], l.hub.hubId, args[3])
		case "join":
			username := args[2]
			userRoomId := l.Users[username].roomId
			if userRoomId != l.hub.hubId {
				l.lobbies[userRoomId].unregister <- l.Users[username]
			} else {
				return
			}
			res = messages.HubJoinUser(username, l.hub.hubId)
			l.hub.register <- l.Users[username]
		case "lobby":
			fmt.Println("Creating new lobby!", args)
			roomId, _ := strconv.Atoi(args[1])
			if args[2] == "ws" {
				roomId = l.idGen()
				res = messages.HubCreateLobby("tcp", roomId)
			} else {
				l.idGen = idGenerator(roomId)
			}
			lobby := NewLobby(roomId)
			l.mutex.Lock()
			defer l.mutex.Unlock()
			l.lobbies[lobby.lobbyId] = lobby
			go lobby.run()
			fmt.Println(lobby.lobbyId)
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
		_, ok := l.lobbies[lobbyId]
		if !ok {
			fmt.Println("given invalid lobby id")
			return
		}
		switch args[0] {
		case "broadcast":
			username := args[2]
			res = messages.LobbyBroadcast(username, l.Users[username].roomId, args[3])
		case "join":
			lobby, ok := l.lobbies[lobbyId]
			if !ok {
				fmt.Println("given invalid room id")
				return
			}
			username := args[2]
			userRoomId := l.Users[username].roomId
			fmt.Println("hub's id:", l.hub.hubId)
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
