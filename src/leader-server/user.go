package leaderserver

type User struct {
	username string
	roomtype string // hub, lobby, application
	serverId int
	roomId   int
}

func NewUser(username string, serverId int, roomId int) User {
	return User{
		username: username,
		roomtype: "hub",
		serverId: serverId,
		roomId:   roomId,
	}
}
