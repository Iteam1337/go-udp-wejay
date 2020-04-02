package rooms

import (
	"github.com/Iteam1337/go-udp-wejay/room"
	"github.com/Iteam1337/go-udp-wejay/users"
)

type Rooms struct {
	rooms map[string]*room.Room
}

func (r *Rooms) Get(id string) (room *room.Room) {
	if res, ok := r.rooms[id]; ok {
		room = res
	}

	return
}

func (r *Rooms) Add(userID string, id string) (out room.Room, ok bool) {
	user, _ := users.GetUser(userID)
	ex := rooms.Get(user.Room)

	if ex != nil {
		if id, empty := ex.Evict(userID); empty {
			delete(r.rooms, id)
		}
	}

	if res, ok := r.rooms[id]; ok {
		res.Add(userID)
		out = *res
	} else {
		user.Room = id
		out = room.New(id, userID)
	}

	if out.Size() > 0 {
		ok = true
	}

	return
}

func (r *Rooms) Exists(id string) bool {
	if _, ok := r.rooms[id]; ok {
		return true
	}

	return false
}

func (r *Rooms) Delete(id string) {
	delete(r.rooms, id)
}

var (
	rooms = Rooms{
		rooms: make(map[string]*room.Room),
	}
	Get    = rooms.Get
	Add    = rooms.Add
	Exists = rooms.Exists
	Delete = rooms.Delete
)
