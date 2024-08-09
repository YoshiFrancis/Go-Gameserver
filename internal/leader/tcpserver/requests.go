package tcpserver

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

func (s *TCPServer) handleFollowerRequest(req Request) {
	fmt.Println("New request!")
	args := req.args

	switch args[0] {
	case "register":
		fmt.Println("Registering user")
		s.userStorage.Set(args[2], *rooms.NewUser(args[2], s.hub))
	case "broadcast":
		fmt.Println("Broadcasting message to users room")
		roomId, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
			return
		}
		if roomId == -1 {
			roomId = 1
		}

		if room, ok := s.roomStorage.Get(roomId); !ok {
			fmt.Println("Not found roomId: ", roomId)
			return
		} else {
			broadcastMsg := room.Broadcast(args[2], args[3])
			fmt.Println("Broadcast msg: ", broadcastMsg)
			s.fbroadcast(broadcastMsg)
		}
	case "join":
		roomId, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
			return
		}
		new_room, ok := s.roomStorage.Get(roomId)
		if !ok {
			return
		}

		user, ok := s.userStorage.Get(args[2])
		if !ok {
			return
		}

		new_room.Join(user)
	case "lobby":
		creatorUsername := args[1]
		lobbyTitle := args[2]
		new_lobby := rooms.NewLobby(s.idGen(), s.hub, lobbyTitle, creatorUsername)
		user, ok := s.userStorage.Get(creatorUsername)
		if ok {
			new_lobby.Join(user)
		}
		broadcastMsg := s.hub.Broadcast("Server", lobbyTitle+" lobby created by "+creatorUsername)
		s.fbroadcast(broadcastMsg)
	}
}
