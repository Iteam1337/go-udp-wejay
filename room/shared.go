package room

import (
	"fmt"
	"log"
	"time"

	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/ankjevel/spotify"
)

func sleep() {
	time.Sleep(10 * time.Second)
}

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

func (r *Room) getOwnerClient() (client *spotify.Client) {
	ownerID := r.playlist.Owner.ID

	for _, user := range r.users {
		if string(user.ClientID) != ownerID {
			continue
		}

		client = user.GetClient()
		break
	}

	if client == nil {
		log.Printf(`[%s] no client`, r.id)
	}

	return
}

func (r *Room) promoteNewOwner() {
	for _, user := range r.users {
		user.Promote()
		r.owner = user
		break
	}
}

func (r *Room) includesClient(clientID string) bool {
	if _, ok := r.clientIDs[spotify.ID(clientID)]; ok {
		return true
	}

	return false
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
