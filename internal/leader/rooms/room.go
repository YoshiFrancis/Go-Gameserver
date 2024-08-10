package rooms

type Room interface {
	Join(user User) []byte
	Leave(user User)
	Broadcast(sender, message string) string
	LeavingMessage(leavingUser string) string
	JoiningMessage(joiningUser string) string
	GetInfo() string
	GetName() string
}

type Message struct {
	username string
	text     string
}
