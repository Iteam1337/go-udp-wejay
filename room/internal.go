package room

import (
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

func (r *Room) checkPlaylistSongs(client *spotify.Client) {
	pl, err := client.GetPlaylistTracks(r.playlist.ID)
	if err != nil {
		log.Println("can't get playlist tracks", err)
	}

	var tracksToRemove []spotify.PlaylistTrack
	for _, track := range pl.Tracks {
		if r.includesClient(track.AddedBy.ID) {
			continue
		}

		tracksToRemove = append(tracksToRemove, track)
	}

	log.Println(tracksToRemove)
}

func (r *Room) listen() {
	for {
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

		r.checkPlaylistSongs(client)

		sleep()
	}
}
