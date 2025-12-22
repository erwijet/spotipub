package playback

import (
	"encoding/json"
	"net/http"
)

func GetCurrent(w http.ResponseWriter, r *http.Request) {
	client := WaitForClient()
	playback, err := client.PlayerCurrentlyPlaying(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(playback)
}
