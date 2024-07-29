package storage

type Room interface {
	join(user User)
	leave(user User)
	deliverAll(message string)
	handleMessage(message string, sender string)
	getInfo() string
	getUserStorage() *Storage[string, User]
}
