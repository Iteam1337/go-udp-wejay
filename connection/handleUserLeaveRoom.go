package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/rooms"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleUserLeaveRoom() {
	msg := c.msg.(*message.UserLeaveRoom)
	log.Println("handleUserLeaveRoom", msg)

	res := message.UserLeaveRoomResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	if exists := users.Exists(msg.UserId); !exists {
		res.Error = "user does not exist"
		c.send(&res)
		return
	}

	user, _ := users.GetUser(msg.UserId)
	ex := rooms.Get(user.Room)

	if ex != nil {
		id, empty := ex.Evict(msg.UserId)
		if empty {
			rooms.Delete(id)
		}

		res.UserId = msg.UserId

		res.Ok = true
	}

	c.send(&res)
}
