package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/erwijet/spotipub/spotify"
	"github.com/erwijet/spotipub/sse"
)

func main() {
	notifier := spotify.NewNotifier()
	mux := http.NewServeMux()

	go notifier.Run()
	spotify.BeginAuthFlow()

	mux.HandleFunc("/current", spotify.GetCurrent)
	mux.HandleFunc("/callback", spotify.GetCallback)

	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		l := notifier.NewListener()
		defer l.Cleanup()

		sse := sse.NewSSEWriter(&w)
		sse.WriteHeaders()

		initial, err := spotify.WaitForClient().PlayerCurrentlyPlaying(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		sse.WriteEvent("initial", initial)

	loop:
		for {
			select {
			case <-w.(http.CloseNotifier).CloseNotify():
				break loop

			case data, ok := <-l.Ch:
				if !ok {
					break loop
				}

				sse.WriteEvent("update", data)
				time.Sleep(2 * time.Second)
			}
		}
	})

	//

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
