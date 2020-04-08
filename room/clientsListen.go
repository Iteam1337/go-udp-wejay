package room

import (
	"sync"
	"time"
)

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

		client := r.getOwnerClient()
		if client == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		current := r.getCurrent(client)
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
