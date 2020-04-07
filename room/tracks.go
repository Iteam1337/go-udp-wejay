package room

import (
	"fmt"
	"log"

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

	t := *tracks
	c := make([]spotify.PlaylistTrack, 1)

	if prev == nil {
		copy(c, t[:1])
		current = c[0]
		return
	}

	for _, track := range *tracks {
		if prev.Track.ID == track.Track.ID {
			continue
		}
		current = track
		break
	}

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

	current = r.getCurrentTrack(tracks, p)

	var trackIDs []spotify.ID
	for _, track := range *tracks {
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
