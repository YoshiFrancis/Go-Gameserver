package rooms

import "github.com/yoshifrancis/go-gameserver/internal/containers"

type Room interface {
	Join(user User)
	Leave(user User)
	Broadcast(sender, message string) string
	GetInfo() string
	getUserStorage() *containers.Storage[string, User]
}
