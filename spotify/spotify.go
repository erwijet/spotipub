package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:3000/callback"

type Unit struct{}

var (
	auth    = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying))
	channel = make(chan Unit)
	state   = "abc123"
	once    sync.Once
	client  *spotify.Client
)

func setClient(c *spotify.Client) {
	client = c
	close(channel)
}

func WaitForClient() *spotify.Client {
	<-channel
	return client
}

func Initialize() {
	url := auth.AuthURL("abc123")
	fmt.Println("Please log into Spotify: ", url)

	client := WaitForClient()

	currentSong, err := client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are listening to... ", currentSong)
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

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprint(w, "login complete!")
	setClient(client)
}

func GetCurrent(w http.ResponseWriter, r *http.Request) {
	// http.Redirect("asdf", w, r)
}
