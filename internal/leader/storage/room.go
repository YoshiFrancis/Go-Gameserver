package storage

type Room interface {
	HandleMessage(message string, sender string)
	Join(user User)
	Leave(user User)
	DeliverAll(message string)
	GetInfo() string
	getUserStorage() *Storage[string, User]
}
