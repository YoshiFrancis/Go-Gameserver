package main

import (
	"log"
	"net/http"

	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

func main() {
	wsserver := wsserver.NewServer()
	go wsserver.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsserver.Serve(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
