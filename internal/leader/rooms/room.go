package rooms

type Room interface {
	Join(user User) []byte
	Leave(user User)
	Broadcast(sender, message string) string
	GetInfo() string
	GetName() string
}

type Message struct {
	username string
	text     string
}
