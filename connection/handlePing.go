package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/utils"
)

func (c *Connection) handlePing() {
	utils.SendM(c.it.Inv(), &message.Pong{Int: 0}, c.conn, c.addr)
}
