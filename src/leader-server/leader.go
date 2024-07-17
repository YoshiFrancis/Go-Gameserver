package leaderserver

type Leader struct {
	servers      map[int]bool         // server id, isExist
	clients      map[string]int       // username, server id
	hub          *Hub                 // room title, room id
	lobbies      map[int]*Lobby       // lobby id, Lobby pointer
	applications map[int]*Application // app id, application pointer
	idGen        func() int           // used to generate ids
}

func NewLeader() *Leader {
	return &Leader{
		servers:      make(map[int]bool),
		clients:      make(map[string]int),
		hub:          NewHub(),
		lobbies:      make(map[int]*Lobby),
		applications: make(map[int]*Application),
		idGen:        idGenerator(),
	}
}

func idGenerator() func() int {
	id := 1
	return func() int {
		id++
		return id
	}
}
