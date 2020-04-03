package room

import (
	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/users"
)

type Room struct {
	id    string
	users map[string]*user.User
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
}

func (r *Room) Size() int {
	return len(r.users)
}

func New(id string, userID string) (r Room) {
	u, err := users.GetUser(userID)

	if u == nil || err != nil {
		return
	}

	r.users = map[string]*user.User{}
	r.users[userID] = u
	r.id = id
	return
}
