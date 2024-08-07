package rooms

type User struct {
	username string
	room     Room
}

func NewUser(username string, hub *Hub) *User {
	return &User{
		username: username,
		room:     hub,
	}
}

func (user *User) GetRoom() Room {
	return user.room
}
