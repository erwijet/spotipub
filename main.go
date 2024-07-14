package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name   string  `json:"name"`
	Gizmos []Gizmo `json:"gizmos"`
}

type Gizmo struct {
	Title string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	hub := newHub()
	mux := http.NewServeMux()

	go hub.run()

	mux.HandleFunc("/user/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		user := &User{
			Name: name,
			Gizmos: []Gizmo{
				Gizmo{Title: "dingus", Price: 12},
				Gizmo{Title: "eh", Price: 5},
			},
		}

		b, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.Write(b)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
