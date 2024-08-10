package tcpserver

import (
	"fmt"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

func (s *TCPServer) handleFollowerRequest(req Request) {
	switch req.command {
	case "register":

		new_user := rooms.NewUser(req.arg, s.hub)
		s.userStorage.Set(req.arg, new_user)
		s.hub.Join(new_user)
		s.fbroadcast(s.hub.JoiningMessage(req.arg))

	case "broadcast":

		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			fmt.Println("Could not find user")
			return
		}

		room := user.GetRoom()
		fmt.Println(room.GetName())
		broadcastMsg := room.BroadcastMessage(req.sender, req.arg)
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

		prevRoom := user.GetRoom()
		leavingTmpl, joiningTmpl := new_room.Join(user)

		s.fbroadcast(prevRoom.BroadcastTemplate(string(leavingTmpl)))
		s.fbroadcast(user.GetRoom().BroadcastTemplate(string(joiningTmpl)))

	case "lobby":
		creatorUsername := req.sender
		lobbyTitle := req.arg
		user, ok := s.userStorage.Get(creatorUsername)
		if !ok {
			fmt.Println("Could not find the user: ", creatorUsername)
			return
		}
		new_lobby := rooms.NewLobby(s.idGen(), s.hub, lobbyTitle, creatorUsername)
		s.roomStorage.Set(lobbyTitle, new_lobby)
		new_lobby.Join(user)
		broadcastMsg := s.hub.BroadcastMessage("Server", lobbyTitle+" lobby created by "+creatorUsername)
		s.fbroadcast(broadcastMsg)
	}
}
