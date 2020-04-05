package room

import (
	"fmt"
	"log"
	"time"

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
		log.Println("can't get playlist tracks", err)
		return
	}

	if pl.Total == 0 {
		return
	}

	current = pl.Tracks[0]

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
			log.Println(err)
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

func (r *Room) clientsListen() {
	for {
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

		for _, u := range r.users {
			client := u.GetClient()
			ps, err := client.PlayerState()
			if err != nil || !ps.Playing {
				continue
			}
			if r.currentTrack.Track.ID == ps.CurrentlyPlaying.Item.ID && r.acceptableTimeDiff(ps.CurrentlyPlaying.Progress) {
				continue
			}
			po := spotify.PlayOptions{
				PlaybackContext: &r.playlist.URI,
				PositionMs:      int(r.Elapsed().Milliseconds()),
			}
			if err := client.PlayOpt(&po); err != nil {
				log.Println("could not set context", u.ClientID, err)
				return
			}
			if ps.Playing {
				return
			}
			if err := client.Play(); err != nil {
				log.Println("could not play", u.ClientID, err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (r *Room) ownerListen() {
	for {
		now := time.Now()
		r.elapsed = now
		if !r.active || r.Size() < 1 {
			r.destroy()
			break
		}

		client := r.owner.GetClient()
		if client == nil {
			log.Println("no client :(")
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
		}

		time.Sleep(sleep - time.Since(now))

		if removeTrack {
			log.Println("removing track", current.Track.ID)
			_, err := client.RemoveTracksFromPlaylist(r.playlist.ID, current.Track.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
