package playback

import (
	"context"
	"math"
	"time"

	"github.com/erwijet/spotipub/internal/logging"
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
	log := logging.GetLogger("Notifier updateAll")

	if update.Item != nil {
		log.Printf("notifying %d listener(s) [%s@%d ; paused?: %t]\n", len(n.listeners), update.Item.SimpleTrack.Name, update.Progress, !update.Playing)
	} else {
		log.Printf("notifying %d listener(s) [no item]\n", len(n.listeners))
	}

	for _, l := range n.listeners {
		go func() {
			l.Ch <- update
		}()
	}
}

func (n *Notifier) Run() {
	client := WaitForClient()
	log := logging.GetLogger("Notifier Run")

	for {
		func() {
			data, err := client.PlayerCurrentlyPlaying(context.Background())
			if err != nil {
				log.Printf("spotify notifier: failed to poll playback: %v", err)
				time.Sleep(5 * time.Second)
				return
			}

			if n.prev_data == nil {
				n.prev_data = data
				n.prev_time = time.Now()
				return
			}

			cur_time := time.Now()
			progress_delta_sec := (data.Progress / 1000) - (n.prev_data.Progress / 1000)
			time_delta_sec := cur_time.Unix() - n.prev_time.Unix()

			jitter := math.Abs(float64(progress_delta_sec - libspot.Numeric(time_delta_sec)))

			defer func() {
				n.prev_data = data
				n.prev_time = cur_time

				time.Sleep(2 * time.Second)
			}()

			if n.prev_data.Item == nil || data.Item == nil {
				return
			}

			if jitter > 1 || n.prev_data.Playing != data.Playing {
				n.notifyAll(data)
				return
			}

			if n.prev_data.Item.Name != data.Item.Name {
				n.notifyAll(data)
			}
		}()
	}
}
