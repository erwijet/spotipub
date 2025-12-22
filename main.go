package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/erwijet/spotipub/internal/playback"
	"github.com/erwijet/spotipub/internal/sse"
)

func main() {
	notifier := playback.NewNotifier()
	mux := http.NewServeMux()

	go notifier.Run()
	playback.BeginAuthFlow()

	mux.HandleFunc("/current", playback.GetCurrent)
	mux.HandleFunc("/callback", playback.GetCallback)

	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		l := notifier.NewListener()
		defer l.Cleanup()

		sse := sse.NewSSEWriter(w)
		sse.WriteHeaders()

		initial, err := playback.WaitForClient().PlayerCurrentlyPlaying(context.Background())
		if err != nil {
			http.Error(w, "failed to fetch current playback", http.StatusBadGateway)
			log.Printf("sse initial playback fetch failed: %v", err)
			return
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
