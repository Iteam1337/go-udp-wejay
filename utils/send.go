package utils

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/golang/protobuf/proto"
)

func send(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	if _, err := conn.WriteToUDP([]byte(msg), addr); err != nil {
		log.Println(err)
	}
}

// Send …
func Send(m types.MessageType, msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	b := m.ByteAndVersion()
	send(append(b[:], msg...), conn, addr)
}

// SendM …
func SendM(m types.MessageType, pb proto.Message, conn *net.UDPConn, addr *net.UDPAddr) {
	b := m.ByteAndVersion()
	if msg, err := proto.Marshal(pb); err != nil {
		log.Println(err)
	} else {
		send(append(b[:], msg...), conn, addr)
	}
}

// SendEmpty …
func SendEmpty(conn *net.UDPConn, addr *net.UDPAddr) {
	send([]byte{}, conn, addr)
}
