package tcpserver

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

func (s *TCPServer) handleFollowerRequest(req Request) {
	switch req.command {
	case "register":
		fmt.Println("Registering user")
		s.userStorage.Set(req.arg, *rooms.NewUser(req.arg, s.hub))
	case "broadcast":
		fmt.Println("Broadcasting message to users room")
		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			fmt.Println("Could not find user")
			return
		}

		room := user.GetRoom()
		broadcastMsg := room.Broadcast(req.sender, req.arg)
		fmt.Println("Broadcast msg: ", broadcastMsg)
		s.fbroadcast(broadcastMsg)

	case "join":
		roomId, err := strconv.Atoi(req.arg)
		if err != nil {
			log.Println(err)
			return
		}
		new_room, ok := s.roomStorage.Get(roomId)
		if !ok {
			return
		}
		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			return
		}
		new_room.Join(user)

	case "lobby":
		creatorUsername := req.sender
		lobbyTitle := req.arg
		new_lobby := rooms.NewLobby(s.idGen(), s.hub, lobbyTitle, creatorUsername)
		user, ok := s.userStorage.Get(creatorUsername)
		if ok {
			new_lobby.Join(user)
		}
		broadcastMsg := s.hub.Broadcast("Server", lobbyTitle+" lobby created by "+creatorUsername)
		s.fbroadcast(broadcastMsg)
	}
}
