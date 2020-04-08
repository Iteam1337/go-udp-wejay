package room

import (
	"fmt"
	"log"
	"time"

	"github.com/ankjevel/spotify"
)

func (r *Room) getCurrentTracks(client *spotify.Client) (tracks *[]spotify.PlaylistTrack, err error) {
	pl, err := client.GetPlaylistTracks(r.playlist.ID)
	if err != nil {
		log.Printf(`[%s] can't get playlist tracks: %s`, r.id, err)
		return
	}

	if pl.Total == 0 {
		err = fmt.Errorf("no tracks")
		return
	}

	tracks = &pl.Tracks

	return
}

func (r *Room) getCurrent(client *spotify.Client) (current spotify.PlaylistTrack) {
	tracks, err := r.getCurrentTracks(client)
	if err != nil {
		return
	}

	return r.getCurrentTrack(tracks, nil)
}

func (r *Room) getCurrentTrack(tracks *[]spotify.PlaylistTrack, prev *spotify.PlaylistTrack) (current spotify.PlaylistTrack) {
	if tracks == nil {
		return
	}

	for _, track := range *tracks {
		if prev != nil && prev.Track.ID == track.Track.ID {
			continue
		}

		if !r.okTrack(track) {
			continue
		}

		current = track
		break
	}

	return
}

func (r *Room) okTrack(track spotify.PlaylistTrack) bool {
	dur := track.Track.TimeDuration()
	if dur > (7*time.Minute) || dur < (30*time.Second) {
		return false
	}

	if !r.includesClient(track.AddedBy.ID) {
		return false
	}

	return true
}

func (r *Room) orderTracks(client *spotify.Client, tracks *[]spotify.PlaylistTrack, prev *spotify.PlaylistTrack) (out *[]spotify.PlaylistTrack) {
	if tracks == nil {
		return
	}

	var ordered []spotify.PlaylistTrack

	log.Println(len(*tracks))

	for _, track := range *tracks {
		if !r.okTrack(track) {
			continue
		}

		ordered = append(ordered, track)
	}

	log.Println(len(ordered))

	out = &ordered

	return
}

func (r *Room) getCurrentAndRemoveOldPlaylistTracks(client *spotify.Client, prev spotify.PlaylistTrack) (current spotify.PlaylistTrack) {
	tracks, err := r.getCurrentTracks(client)
	if err != nil {
		return
	}

	var p *spotify.PlaylistTrack

	if prev.Track.ID != "" {
		p = &prev
	}

	ordered := r.orderTracks(client, tracks, p)
	current = r.getCurrentTrack(ordered, p)

	var trackIDs []spotify.ID
	for _, track := range *tracks {
		if r.okTrack(track) {
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
