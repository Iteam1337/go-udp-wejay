package room

import (
	"log"
	"sync"

	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/ankjevel/spotify"
)

func (r *Room) addUsersTopTrack(client *spotify.Client) {
	var wg sync.WaitGroup
	for _, u := range r.users {
		wg.Add(1)
		go func(u *user.User) {
			defer wg.Done()

			c := u.GetClient()
			if c == nil {
				return
			}

			t, err := c.CurrentUsersTopTracks()
			if err != nil || t == nil {
				log.Printf(`[%s](%s) cant get users top tracks: %s`, r.id, u.ClientID, err.Error())
			}

			if t.Total == 0 {
				return
			}

			track := t.Tracks[utils.Random(len(t.Tracks))]
			if track.URI == "" {
				return
			}

			if _, err := client.AddTracksToPlaylist(r.playlist.ID, track.ID); err != nil {
				log.Printf(`[%s](%s) add track failed: %s`, r.id, u.ClientID, err.Error())
			} else {
				log.Printf(`[%s](%s) added track: %s`, r.id, u.ClientID, track.Name)
			}
		}(u)
	}

	wg.Wait()
}
