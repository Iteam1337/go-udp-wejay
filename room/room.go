package room

import (
	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/users"
	"github.com/ankjevel/spotify"
)

type Room struct {
	id            string
	users         map[string]*user.User
	playlist      spotify.SimplePlaylist
	playlistOwner spotify.ID
	owner         *user.User
}

func (r *Room) Evict(userID string) (id string, empty bool) {
	delete(r.users, userID)

	id = r.id
	empty = len(r.users) == 0

	u, err := users.GetUser(userID)
	if u != nil && err == nil {
		u.LeaveRoom()

		if r.owner.ClientID == u.ClientID {
			r.promoteNewOwner()
		}
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
	u.JoinRoom(r.id, r.playlist, r.owner.ClientID)
}

func (r *Room) Size() int {
	return len(r.users)
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
	r.id = id
	r.playlist = playlist
	r.playlistOwner = playlistOwner
	r.owner = u

	u.JoinRoom(id, playlist, u.ClientID)
	return
}
