package playback

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	libspot "github.com/zmb3/spotify/v2"
)

type Notifier struct {
	prev_data *libspot.CurrentlyPlaying
	prev_time time.Time
	listeners []Listener
}

type Listener struct {
	owner *Notifier
	Ch    chan *libspot.CurrentlyPlaying
}

func (n *Notifier) NewListener() *Listener {
	n.listeners = append(n.listeners, Listener{
		owner: n,
		Ch:    make(chan *libspot.CurrentlyPlaying),
	})
	return &n.listeners[len(n.listeners)-1]
}

func (l *Listener) Cleanup() {
	for i, cur := range l.owner.listeners {
		if *l == cur {
			l.owner.listeners = append(l.owner.listeners[:i], l.owner.listeners[i+1:]...)
			return
		}
	}
}

func NewNotifier() Notifier {
	return Notifier{
		listeners: make([]Listener, 0),
	}
}

func (n *Notifier) notifyAll(update *libspot.CurrentlyPlaying) {
	if update.Item != nil {
		fmt.Printf("notifying %d listener(s) [%s@%d ; paused?: %t]\n", len(n.listeners), update.Item.SimpleTrack.Name, update.Progress, !update.Playing)
	} else {
		fmt.Printf("notifying %d listener(s) [no item]\n", len(n.listeners))
	}

	for _, l := range n.listeners {
		go func() {
			l.Ch <- update
		}()
	}
}

func (n *Notifier) Run() {
	client := WaitForClient()

	for {
		data, err := client.PlayerCurrentlyPlaying(context.Background())
		if err != nil {
			log.Printf("spotify notifier: failed to poll playback: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if n.prev_data == nil {
			n.prev_data = data
			n.prev_time = time.Now()
			continue
		}

		cur_time := time.Now()
		progress_delta_sec := (data.Progress / 1000) - (n.prev_data.Progress / 1000)
		time_delta_sec := cur_time.Unix() - n.prev_time.Unix()

		jitter := math.Abs(float64(progress_delta_sec - libspot.Numeric(time_delta_sec)))

		if jitter > 1 || n.prev_data.Playing != data.Playing || n.prev_data.Item.Name != data.Item.Name {
			n.notifyAll(data)
		}

		n.prev_data = data
		n.prev_time = cur_time

		time.Sleep(2 * time.Second)
	}
}
