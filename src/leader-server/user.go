package leaderserver

import "github.com/yoshifrancis/go-gameserver/src/wsserver"

type User struct {
	username string
	roomtype string // hub, lobby, application
	serverId int
	roomId   int
	client   *wsserver.Client
}

func NewUser(username string, serverId int, roomId int, client *wsserver.Client) *User {
	return &User{
		username: username,
		roomtype: "hub",
		serverId: serverId,
		roomId:   roomId,
		client:   client,
	}
}

func (user *User) send(msg []byte) {
	user.client.Send <- msg
}
