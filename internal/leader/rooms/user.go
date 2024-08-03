package rooms

import "net"

type User struct {
	username string
	conn     net.Conn
	room     Room
}

func NewUser(username string, serverConn net.Conn, hub *Hub) *User {
	return &User{
		username: username,
		conn:     serverConn,
		room:     hub,
	}
}

func GetRoom(user User) Room {
	return user.room
}
