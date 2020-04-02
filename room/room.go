package room

import (
	"github.com/Iteam1337/go-udp-wejay/utils"
)

// Room …
type Room struct {
	id string
}

// GetRoom …
func (r *Room) GetRoom() {}

// Destroy …
func (u *Room) Destroy() {
	utils.SetNil(&u.id)
}

// New …
func New(id string) (r Room) {
	r.id = id
	return
}
