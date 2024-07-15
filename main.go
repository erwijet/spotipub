package main

import (
	"log"
	"net/http"

	"github.com/erwijet/spotipub/spotify"
)

func main() {
	hub := newHub()
	mux := http.NewServeMux()

	go hub.run()
	go spotify.Initialize()

	mux.HandleFunc("/current", spotify.GetCurrent)
	mux.HandleFunc("/callback", spotify.GetCallback)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	//

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
