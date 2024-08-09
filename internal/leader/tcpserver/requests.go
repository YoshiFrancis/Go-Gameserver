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
		room, ok := s.roomStorage.Get(roomId)
		if !ok {
			return
		}
		fmt.Println("Preparing")
		broadcastMsg := room.Broadcast(args[2], args[3])
		fmt.Println(broadcastMsg)
		s.fbroadcast(broadcastMsg)
	case "join":
		fmt.Println("Joining user to room")
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
	}
}
