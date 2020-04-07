package room

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/ankjevel/spotify"
)

func (r *Room) destroy() {
	utils.SetNil(&r.active)
	utils.SetNil(&r.id)
	utils.SetNil(&r.users)
	utils.SetNil(&r.clientIDs)
	utils.SetNil(&r.playlist)
	utils.SetNil(&r.playlistOwner)
	utils.SetNil(&r.owner)
	utils.SetNil(&r.elapsed)
	utils.SetNil(&r.update)
}

func (r *Room) promoteNewOwner() {
	for _, user := range r.users {
		user.Promote()
		r.owner = user
		break
	}
}

func sleep() {
	time.Sleep(10 * time.Second)
}

func (r *Room) includesClient(clientID string) bool {
	if _, ok := r.clientIDs[spotify.ID(clientID)]; ok {
		return true
	}

	return false
}

func (r *Room) checkPlaylistSongs(client *spotify.Client) (current spotify.PlaylistTrack) {
	pl, err := client.GetPlaylistTracks(r.playlist.ID)
	if err != nil {
		log.Printf(`[%s] can't get playlist tracks: %s`, r.id, err)
		return
	}

	if pl.Total == 0 {
		return
	}

	c := make([]spotify.PlaylistTrack, 1)
	copy(c, pl.Tracks[:1])
	current = c[0]

	var trackIDs []spotify.ID
	for _, track := range pl.Tracks {
		if r.includesClient(track.AddedBy.ID) {
			continue
		}

		trackIDs = append(trackIDs, track.Track.ID)
	}

	if len(trackIDs) > 0 {
		_, err := client.RemoveTracksFromPlaylist(r.playlist.ID, trackIDs...)
		if err != nil {
			log.Printf(`[%s] cant remove tracks: %s`, r.id, err)
		}
	}

	return
}

func (r *Room) acceptableTimeDiff(progress int) (ok bool) {
	var diff time.Duration
	d, _ := time.ParseDuration(fmt.Sprintf("%dms", progress))
	e := r.Elapsed()

	if e > d {
		diff = e - d
	} else if e < d {
		diff = d - e
	} else {
		return
	}

	ok = diff < 5*time.Second

	return
}

func (r *Room) updateUser(u *user.User, wg *sync.WaitGroup) {
	defer wg.Done()

	client := u.GetClient()
	if client == nil {
		return
	}

	ps, err := client.PlayerState()
	if err != nil || ps == nil || !ps.Playing {
		return
	}

	if ps.CurrentlyPlaying.Item == nil {
		log.Printf("[%s](%s) nothing is playing", r.id, u.ClientID)
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

func (r *Room) clientsListen() {
	for {
		if r.update {
			time.Sleep(5 * time.Second)
			continue
		}

		if !r.active || r.Size() < 1 {
			r.destroy()
			break
		}

		client := r.owner.GetClient()
		if client == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		current := r.checkPlaylistSongs(client)
		if current.Track.ID == "" || current.Track.ID != r.currentTrack.Track.ID {
			sleep()
			continue
		}

		var wg sync.WaitGroup
		for _, u := range r.users {
			if r.update {
				continue
			}
			wg.Add(1)
			go r.updateUser(u, &wg)
		}
		wg.Wait()
		time.Sleep(10 * time.Second)
	}
}

func (r *Room) ownerListen() {
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

		current := r.checkPlaylistSongs(client)
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

		time.Sleep(sleep - time.Since(now))

		if !removeTrack {
			continue
		}

		log.Printf(`[%s] removing: "%s"`, r.id, current.Track.Name)
		_, err := client.RemoveTracksFromPlaylist(r.playlist.ID, current.Track.ID)
		if err != nil {
			log.Println(err)
		}
	}
}
