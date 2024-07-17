package leaderserver

type Lobby struct {
	lobbyId      int
	users        map[User]bool
	member_count int
	appId        int
}

func NewLobby(lobbyId, appId int) *Lobby {
	return &Lobby{
		lobbyId:      lobbyId,
		users:        make(map[User]bool),
		member_count: 0,
		appId:        appId,
	}
}
