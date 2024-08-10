package tcpserver

import (
	"fmt"

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
		new_room, ok := s.roomStorage.Get(req.arg)
		if !ok {
			return
		}
		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			return
		}

		new_room.Join(user)

		s.fbroadcast(user.GetRoom().LeavingMessage(req.sender))
		s.fbroadcast(user.GetRoom().JoiningMessage(req.sender))

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
