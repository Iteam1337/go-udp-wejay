package rooms

import (
	"log"
	"strings"

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

func (r *Rooms) Restore() {
	for _, user := range users.GetAll() {
		id := user.Room
		userID := user.ID

		if id == "" {
			continue
		}

		if existingRoom, ok := r.rooms[id]; ok {
			existingRoom.Add(userID)
		} else {
			newRoom := room.New(id, userID)
			r.rooms[id] = &newRoom
		}
	}
}

func (r *Rooms) Add(userID string, id string) (out *room.Room, ok bool) {
	user, _ := users.GetUser(userID)
	ex := rooms.Get(user.Room)

	if ex != nil {
		if id, empty := ex.Evict(userID); empty {
			delete(r.rooms, id)
		}
	}

	if existingRoom, ok := r.rooms[id]; ok {
		if existingRoom.Size() < 30 {
			existingRoom.Add(userID)
			out = existingRoom
		}
	} else {
		newRoom := room.New(id, userID)
		r.rooms[id] = &newRoom
		out = &newRoom
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

func (r *Rooms) InRoom(id string, userID string) (isInRoom bool) {
	if res, ok := r.rooms[id]; ok {
		isInRoom = res.Includes(userID)
	}

	return
}

type QueryResult struct {
	Name string
	Size int
}

func (r *Rooms) Available(filter string) (result []QueryResult) {
	for key, res := range r.rooms {
		if filter == "" || strings.Contains(key, filter) {
			result = append(result, QueryResult{
				Name: key,
				Size: res.Size(),
			})
		}
	}

	return
}

var (
	rooms     = Rooms{rooms: make(map[string]*room.Room)}
	Get       = rooms.Get
	Add       = rooms.Add
	Exists    = rooms.Exists
	Delete    = rooms.Delete
	Evict     = rooms.Evict
	Available = rooms.Available
	Restore   = rooms.Restore
)
