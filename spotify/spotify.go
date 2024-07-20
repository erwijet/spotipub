package spotify

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	libspot "github.com/zmb3/spotify/v2"
	libspotauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth    = libspotauth.New(libspotauth.WithRedirectURL(getRedirectUri()), libspotauth.WithScopes(libspotauth.ScopeUserReadCurrentlyPlaying))
	channel = make(chan int)
	state   = randState()
	owned   *libspot.Client
)

func getRedirectUri() string {
	uri, exists := os.LookupEnv("REDIRECT_URI")
	if exists {
		return uri
	} else {
		return "http://localhost:3000/callback"
	}
}

func randState() string {
	s := "abcdef123"
	runes := []rune(s)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}

func setClient(c *libspot.Client) {
	owned = c
	close(channel)
}

func WaitForClient() *libspot.Client {
	<-channel
	return owned
}

func BeginAuthFlow() {
	url := auth.AuthURL(state)
	fmt.Println("Please log into Spotify: ", url)
}

func GetCallback(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := libspot.New(auth.Client(r.Context(), tok))
	fmt.Fprint(w, "login complete!")
	setClient(client)
}

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
