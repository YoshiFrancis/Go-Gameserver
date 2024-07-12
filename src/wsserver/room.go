package wsserver

type Room struct {
	clients      map[string]*Client
	parentRoom   *Room
	title        string
	member_count int
	register     chan *Client
	unregister   chan *Client
}

func NewRoom(title string, parentRoom *Room) *Room {
	if title == "" {
		title = "Lobby"
	}

	return &Room{
		clients:      make(map[string]*Client),
		parentRoom:   parentRoom,
		title:        title,
		member_count: 0,
		register:     make(chan *Client),
		unregister:   make(chan *Client),
	}
}
