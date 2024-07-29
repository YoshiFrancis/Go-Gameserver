package storage

type Room interface {
	join(username string)
	leave(username string)
	deliverAll(message string)
	handleMessage(message string, sender string)
	getInfo() string
}
