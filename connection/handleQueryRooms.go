package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/rooms"
)

func (c *Connection) handleQueryRooms() {
	msg := c.msg.(*message.QueryRooms)
	res := message.QueryRoomsResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	res.Ok = true
	for _, result := range rooms.Available(msg.Name) {
		res.Room = append(res.Room, &message.RefRoom{
			Id:   result.Name,
			Size: int32(result.Size),
		})
	}

	c.send(&res)
}
