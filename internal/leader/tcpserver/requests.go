package tcpserver

import (
	"log"
	"strconv"

	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

func (s *TCPServer) handleFollowerRequest(req Request) {
	args := req.args

	switch args[0] {
	case "register":
		s.userStorage.Set(args[2], *rooms.NewUser(args[2], &s.hub))
	case "broadcast":
		roomId, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
			return
		}
		room, ok := s.roomStorage.Get(roomId)
		if !ok {
			return
		}
		broadcastMsg := room.Broadcast(args[2], args[3])
		s.fbroadcast(broadcastMsg)
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
	}
}
