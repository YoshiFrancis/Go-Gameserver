package rooms

import "github.com/yoshifrancis/go-gameserver/internal/containers"

type Room interface {
	HandleMessage(message string, sender string)
	Join(user User)
	Leave(user User)
	DeliverAll(message string)
	GetInfo() string
	getUserStorage() *containers.Storage[string, User]
}
