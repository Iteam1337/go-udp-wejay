package room

import (
	"log"
	"sync"
	"time"

	"github.com/ankjevel/spotify"
)

func (r *Room) ownerListen() {
	// checkIfEmpy := func(client *spotify.Client) {
	// 	tracks, err := r.getCurrentTracks(client)
	// 	if tracks == nil || (err != nil && err.Error() == "no tracks") || len(*tracks) <= 1 {
	// 		r.addUsersTopTrack(client)
	// 	}
	// }

	doRemoveTrack := func(current spotify.PlaylistTrack) {
		client := r.owner.GetClient()
		if client == nil {
			log.Printf(`[%s] no client`, r.id)
			return
		}

		log.Printf(`[%s] removing: "%s"`, r.id, current.Track.Name)
		_, err := client.RemoveTracksFromPlaylist(r.playlist.ID, current.Track.ID)
		if err != nil {
			log.Println(err)
		}
	}

	forceUpdate := func() {
		r.update = true
		var wg sync.WaitGroup
		for _, u := range r.users {
			wg.Add(1)
			go r.updateUser(u, &wg)
		}
		wg.Wait()
		r.update = false
	}

	var current spotify.PlaylistTrack
	for {
		now := time.Now()
		r.elapsed = now
		if !r.active || r.Size() < 1 {
			r.destroy()
			break
		}

		client := r.owner.GetClient()
		if client == nil {
			log.Printf(`[%s] no client`, r.id)
			sleep()
			continue
		}

		current = r.getCurrentAndRemoveOldPlaylistTracks(client, current)
		var sleep time.Duration
		var removeTrack bool

		if current.Track.ID == "" {
			sleep = 10 * time.Second
		} else {
			r.currentTrack = current
			sleep = current.Track.TimeDuration()
			removeTrack = true
			log.Printf(`[%s] playing: "%s"`, r.id, current.Track.Name)
		}

		forceUpdate()

		// time.Sleep(sleep - (5 * time.Second) - time.Since(now))
		// go checkIfEmpy(client)
		// time.Sleep((5 * time.Second) - time.Since(now))

		time.Sleep(sleep - (3 * time.Second) - time.Since(now))

		if removeTrack {
			doRemoveTrack(current)
		}

		time.Sleep((3 * time.Second) - time.Since(now))
	}
}
