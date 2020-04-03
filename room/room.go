package room

import (
	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/users"
	"github.com/zmb3/spotify"
)

type Room struct {
	id       string
	users    map[string]*user.User
	playlist spotify.SimplePlaylist
}

func (r *Room) Evict(userID string) (id string, empty bool) {
	delete(r.users, userID)

	id = r.id
	empty = len(r.users) == 0

	u, err := users.GetUser(userID)
	if u != nil && err == nil {
		u.LeaveRoom()
	}

	return
}

func (r *Room) Add(userID string) {
	if _, ok := r.users[userID]; ok {
		return
	}

	u, err := users.GetUser(userID)
	if u == nil || err != nil {
		return
	}

	r.users[userID] = u
	u.JoinRoom(r.id, r.playlist)
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

	r.users = map[string]*user.User{}
	r.users[userID] = u
	r.id = id
	r.playlist = playlist

	u.JoinRoom(id, playlist)
	return
}
