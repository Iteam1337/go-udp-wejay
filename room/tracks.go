package room

import (
	"fmt"
	"log"
	"sort"
	"strings"
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
	if dur > (10*time.Minute) || dur < (30*time.Second) {
		return false
	}

	if !r.includesClient(track.AddedBy.ID) {
		return false
	}

	return true
}

type byAdded []spotify.PlaylistTrack

func (a byAdded) Len() int           { return len(a) }
func (a byAdded) Less(i, j int) bool { return a[i].AddedAt < a[j].AddedAt }
func (a byAdded) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func weave(tracksMeta map[string][]spotify.PlaylistTrack, size int) (result []spotify.PlaylistTrack) {
	for _, tracks := range tracksMeta {
		sort.Sort(byAdded(tracks))
	}

	for i := 0; i < size; i++ {
		for _, tracks := range tracksMeta {
			if len(tracks) <= i {
				continue
			}
			result = append(result, tracks[i])
		}
	}
	return
}

func (r *Room) orderTracks(client *spotify.Client, tracks *[]spotify.PlaylistTrack) (out *[]spotify.PlaylistTrack) {
	var (
		preIDs, postIDs []string
		sorted          []spotify.PlaylistTrack
		numTracks       int
	)

	if tracks == nil {
		return
	}

	tracksObject := make(map[spotify.ID]int)
	tracksMeta := make(map[string][]spotify.PlaylistTrack)

	for i, track := range *tracks {
		preIDs = append(preIDs, track.Track.ID.String())
		tracksObject[track.Track.ID] = i
		if !r.okTrack(track) {
			continue
		}

		uID := track.AddedBy.ID
		uTracks := tracksMeta[uID]
		uTracks = append(uTracks, track)

		tracksMeta[uID] = uTracks
		numTracks += 1
	}

	sorted = weave(tracksMeta, numTracks)

	for _, track := range sorted {
		postIDs = append(postIDs, track.Track.ID.String())
	}

	if strings.Join(preIDs, ",") != strings.Join(postIDs, ",") {
		for i, track := range sorted {
			_, err := client.ReorderPlaylistTracks(r.playlist.ID, spotify.PlaylistReorderOptions{
				RangeStart:   tracksObject[track.Track.ID],
				InsertBefore: i,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}

	out = &sorted
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

	ordered := r.orderTracks(client, tracks)
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
