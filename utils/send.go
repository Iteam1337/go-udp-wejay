package utils

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"

	"github.com/golang/protobuf/proto"
)

func send(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) (err error) {
	if _, err = conn.WriteToUDP([]byte(msg), addr); err != nil {
		log.Println(err)
	}
	return
}

func Send(m types.MessageType, msg []byte, conn *net.UDPConn, addr *net.UDPAddr) (e error) {
	b := m.ByteAndVersion()
	e = send(append(b[:], msg...), conn, addr)
	return
}

func SendM(m types.MessageType, pb proto.Message, conn *net.UDPConn, addr *net.UDPAddr) (e error) {
	b := m.ByteAndVersion()
	if msg, err := proto.Marshal(pb); err != nil {
		e = err
		log.Println(err)
	} else {
		e = send(append(b[:], msg...), conn, addr)
	}
	return
}

func SendEmpty(conn *net.UDPConn, addr *net.UDPAddr) (e error) {
	e = send([]byte{}, conn, addr)
	return
}
