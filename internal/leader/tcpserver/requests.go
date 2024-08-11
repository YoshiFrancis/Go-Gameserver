package tcpserver

import (
	"fmt"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

type Request interface {
	getFlag() byte
	getCommand() string
}

type FollowerRequest struct {
	flag    byte
	command string
	arg     string
	sender  string
	server  *ExtenalTCPServer
}

type ApplicationRequest struct {
	flag       byte
	command    string
	arg        string
	lobbyTitle string
	receivers  []string
}

func (f FollowerRequest) getFlag() byte {
	return f.flag
}

func (f FollowerRequest) getCommand() string {
	return f.command
}

func (a ApplicationRequest) getFlag() byte {
	return a.flag
}

func (a ApplicationRequest) getCommand() string {
	return a.command
}

func (s *TCPServer) handleFollowerRequest(req FollowerRequest) {
	switch req.command {
	case "register":

		new_user := rooms.NewUser(req.arg, s.hub)
		s.userStorage.Set(req.arg, new_user)
		_, joiningTmpl := s.hub.Join(new_user)
		s.fbroadcast(s.hub.BroadcastTemplate(string(joiningTmpl)))

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
		leaveTmpl, joinTmpl := new_lobby.Join(user)
		s.fbroadcast(user.GetRoom().BroadcastTemplate(string(joinTmpl)))
		s.fbroadcast(s.hub.BroadcastTemplate(string(leaveTmpl)))
		broadcastMsg := s.hub.BroadcastMessage("Server", lobbyTitle+" lobby created by "+creatorUsername)
		s.fbroadcast(broadcastMsg)
	}
}

func (s *TCPServer) handleApplicationRequest(req ApplicationRequest) {
	if req.command == "shutdown" {
		fmt.Println("Application server is shutting down!")
		// TODO:
		// gracefully shutdown
	}

	// all an application does is send htmx to the users in the lobby
	// req.arg -> htmx
	// req.sender -> lobbyTitle

	// maybe implement some security to make sure that the application has access to the lobby

}
