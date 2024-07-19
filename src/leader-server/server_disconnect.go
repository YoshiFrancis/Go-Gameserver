package leaderserver

/*
programs here help with a graceful exit of a shutting down server.
when a server disconnects, it will send all of its clients urls to other servers in the server group.
if no servers exist, the client will be told that it is no longer in any server group.

The servers will try to split off the clients in a manner so that the distribution across the servers is even.
*/

type TupleHeap []ServerTuple

type ServerTuple struct {
	server_id int
	count     int
}

func (h TupleHeap) Len() int           { return len(h) }
func (h TupleHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h TupleHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TupleHeap) Push(x any) {
	*h = append(*h, x.(ServerTuple))
}

func (h *TupleHeap) Pop() ServerTuple {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (l *Leader) shutdown() {

	// getting the count of users in each server
	// could potentially do this while server is running so that we can renavigate users when they are connecting
	server_id_count := make(map[int]int)
	for _, user := range l.Users {
		server_id_count[user.serverId]++
	}

	tuples := make([]ServerTuple, 0)
	urls := ""
	for server_id, count := range server_id_count {
		tuples = append(tuples, ServerTuple{server_id, count})
		urls += "," + l.TCPServer.Servers[server_id].Url
	}

	var h TupleHeap

	for t := range tuples {
		h.Push(t)
	}

	for _, user := range l.Users {
		tuple := h.Pop()
		server_url := l.TCPServer.Servers[tuple.server_id].Url
		user.send([]byte("/join " + server_url + urls)) // sending all urls in case the user cannot connect to the first url
		tuple.count++
		h.Push(tuple)
	}

	l.WSServer.Shutdown()
	l.TCPServer.Shutdown()
	close(l.TCPrequests)
	close(l.WSrequests)
}
