package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/rooms"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleJoinRoom() {
	msg := c.msg.(*message.JoinRoom)
	log.Println("handleJoinRoom", msg)
	res := message.JoinRoomResponse{Ok: false}

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

	if room, ok := rooms.Add(msg.UserId, msg.RoomId); ok {
		res.Room = &message.RefRoom{
			Id:   msg.RoomId,
			Size: int32(room.Size()),
		}
		res.UserId = msg.UserId
		res.Ok = true
	}

	c.send(&res)
}
