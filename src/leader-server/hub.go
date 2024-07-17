package leaderserver

type Hub struct {
	hubId        int
	users        map[User]bool
	member_count int
}

func NewHub() *Hub {
	return &Hub{}
}
