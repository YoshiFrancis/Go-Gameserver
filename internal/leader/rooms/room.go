package rooms

type Room interface {
	Join(user *User) (leavingTmpl, joiningTmpl []byte)
	Leave(user *User) []byte
	BroadcastMessage(sender, message string) string
	BroadcastTemplate(tmpl string) string
	LeavingMessage(leavingUser string) string
	JoiningMessage(joiningUser string) string
	GetInfo() string
	GetName() string
}

type Message struct {
	Username string
	Text     string
}
