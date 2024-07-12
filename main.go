package main

import (
	"net/http"

	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

func main() {
	wsserver := wsserver.NewServer()

	http.HandleFunc("/go-server", func(w http.ResponseWriter, r *http.Request) {
		wsserver.Serve(w, r)
	})
}
