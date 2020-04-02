package rooms

import (
	"log"

	"github.com/Iteam1337/go-udp-wejay/room"
	"github.com/Iteam1337/go-udp-wejay/users"
)

type Rooms struct {
	rooms map[string]*room.Room
}

func (r *Rooms) Get(id string) (room *room.Room) {
	log.Println(id, r.rooms)

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
		newRoom := room.New(id, userID)
		r.rooms[id] = &newRoom
		out = newRoom
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

func (r *Rooms) Evict(userID string) (ok bool) {
	user, err := users.GetUser(userID)
	if err != nil {
		log.Println(err)
		return
	}

	ok = true
	userRoom := Get(user.Room)
	if userRoom != nil {
		id, empty := userRoom.Evict(userID)
		if empty {
			Delete(id)
		}
	}

	return
}

var (
	rooms = Rooms{
		rooms: make(map[string]*room.Room),
	}
	Get    = rooms.Get
	Add    = rooms.Add
	Exists = rooms.Exists
	Delete = rooms.Delete
	Evict  = rooms.Evict
)
