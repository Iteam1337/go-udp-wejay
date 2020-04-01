package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/utils"
)

func (c *Connection) handlePing() {
	if err := utils.SendM(c.it.Inv(), &message.Pong{Int: 0}, c.conn, c.addr); err != nil {
		log.Println(err)
	}
}
