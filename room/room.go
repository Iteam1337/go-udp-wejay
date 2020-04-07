package room

import (
	"time"

	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/users"
	"github.com/ankjevel/spotify"
)

type Room struct {
	active        bool
	clientIDs     map[spotify.ID]string
	currentTrack  spotify.PlaylistTrack
	elapsed       time.Time
	id            string
	owner         *user.User
	playlist      spotify.SimplePlaylist
	playlistOwner spotify.ID
	update        bool
	users         map[string]*user.User
}

func (r *Room) Evict(userID string) (id string, empty bool) {
	delete(r.users, userID)

	id = r.id
	empty = len(r.users) == 0

	u, err := users.GetUser(userID)
	if u != nil && err == nil {
		delete(r.clientIDs, u.ClientID)
		u.LeaveRoom()

		if r.owner.ClientID == u.ClientID {
			r.promoteNewOwner()
		}
	}

	if empty {
		r.active = false
	}

	return
}

func (r *Room) Includes(userID string) bool {
	if _, ok := r.users[userID]; ok {
		return true
	}

	return false
}

func (r *Room) Add(userID string) {
	if r.Includes(userID) {
		return
	}

	u, err := users.GetUser(userID)
	if u == nil || err != nil {
		return
	}

	r.users[userID] = u
	r.clientIDs[u.ClientID] = userID
	u.JoinRoom(r.id, r.playlist, r.owner.ClientID)
}

func (r *Room) Size() int {
	return len(r.users)
}

func (r *Room) Elapsed() time.Duration {
	return time.Since(r.elapsed)
}

func New(id string, userID string) (r Room) {
	u, err := users.GetUser(userID)

	if u == nil || err != nil {
		return
	}

	playlist, ok := u.NewPlaylist(id)
	if !ok {
		return
	}

	playlistOwner := u.ClientID

	r.users = map[string]*user.User{}
	r.users[userID] = u

	r.clientIDs = map[spotify.ID]string{}
	r.clientIDs[u.ClientID] = userID

	r.id = id
	r.playlist = playlist
	r.playlistOwner = playlistOwner
	r.owner = u
	r.active = true

	u.JoinRoom(id, playlist, u.ClientID)

	go r.ownerListen()
	go func() {
		for {
			if !r.active || r.currentTrack.Track.ID != "" {
				break
			}
			time.Sleep(1 * time.Second)
		}

		r.clientsListen()
	}()

	return
}
