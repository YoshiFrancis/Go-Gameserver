package tcpserver

import (
	"fmt"
	"log"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

type FollowerRequest struct {
	flag    byte
	command string
	arg     string
	sender  string
	server  *ExtenalTCPServer
}

type ApplicationRequest struct {
	command    string
	arg        string
	lobbyTitle string
	receivers  []string
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

	case "app":
		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			return
		}
		appName := user.GetRoom().GetApp()
		app, ok := s.aServers[appName]
		if !ok {
			return
		}
		appReq := messages.ApplicationRequestTo(req.arg, user.GetRoom().GetName(), req.sender, user.GetRoom().GetUsers())
		app.Send <- []byte(appReq)

	case "app-start":
		user, ok := s.userStorage.Get(req.sender)
		if !ok {
			log.Println("Could not find username: ", req.sender)
			return
		}

		app, ok := s.aServers[req.arg]
		if !ok {
			log.Println("Could not find app name: ", req.arg)
			return
		}

		if user.GetRoom().GetApp() == "" {
			appStartReq := messages.ApplicationStart("", user.GetRoom().GetName(), req.sender, user.GetRoom().GetUsers())
			app.Send <- []byte(appStartReq)
		}

	}
}

func (s *TCPServer) handleApplicationRequest(req ApplicationRequest) {
	if req.command == "shutdown" {
		fmt.Println("Application server is shutting down!")
		// TODO:
		// gracefully shutdown
	} else if req.command == "broadcast" { // to all users in lobby
		room, ok := s.roomStorage.Get(req.lobbyTitle)
		if ok {
			s.fbroadcast(messages.LeaderRoomBroadcast(req.arg, room.GetUsers()))
		}
	} else if req.command == "send" { // to directed users
		s.fbroadcast(messages.LeaderRoomBroadcast(req.arg, req.receivers))
	}

	// maybe implement some security to make sure that the application has access to the users

}
