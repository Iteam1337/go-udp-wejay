package rooms

import "github.com/Iteam1337/go-udp-wejay/room"

// Rooms …
type Rooms struct {
	rooms map[string]room.Room
}

// Get …
func (r *Rooms) Get(id string) {}

// Add …
func (r *Rooms) Add(userId string, id string) {}

// Exists …
func (r *Rooms) Exists(id string) bool {
	if _, ok := r.rooms[id]; ok {
		return true
	}

	return false
}

// Delete …
func (r *Rooms) Delete(id string) {
	delete(r.rooms, id)
}

// Global values
var (
	rooms = Rooms{
		rooms: make(map[string]room.Room),
	}
	Get    = rooms.Get
	Add    = rooms.Add
	Exists = rooms.Exists
	Delete = rooms.Delete
)
