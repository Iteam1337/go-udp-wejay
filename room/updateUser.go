package room

import (
	"log"
	"sync"

	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/ankjevel/spotify"
)

func (r *Room) updateUser(u *user.User, wg *sync.WaitGroup) {
	defer wg.Done()

	client := u.GetClient()
	if client == nil {
		return
	}

	ps, err := client.PlayerState()
	if err != nil || ps == nil || !ps.Playing || ps.Device.Restricted {
		return
	}

	if ps.CurrentlyPlaying.Item == nil {
		return
	}

	if r.currentTrack.Track.ID == ps.CurrentlyPlaying.Item.ID && r.acceptableTimeDiff(ps.CurrentlyPlaying.Progress) {
		return
	}

	po := spotify.PlayOptions{
		PlaybackContext: &r.playlist.URI,
		PositionMs:      int(r.Elapsed().Milliseconds()),
	}

	if err := client.PlayOpt(&po); err != nil {
		log.Printf(`[%s](%s) could set context: %s`, r.id, u.ClientID, err)
		return
	}

	if ps.Playing {
		return
	}

	if err := client.Play(); err != nil {
		log.Printf(`[%s](%s) could not play: %s`, r.id, u.ClientID, err)
	}
}
