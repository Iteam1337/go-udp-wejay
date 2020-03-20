package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Iteam1337/go-udp-wejay/message"
	"github.com/golang/protobuf/proto"
)

func listen(conn *net.UDPConn, addr *net.UDPAddr) {
	i := 0
loop:
	for {
		if _, err := conn.WriteToUDP([]byte(fmt.Sprintf(">> %d\n", i)), addr); err != nil {
			continue
		}

		i = i + 1
		if i > 3 {
			break loop
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}

func send(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	if _, err := conn.WriteToUDP([]byte(msg), addr); err != nil {
		log.Println(err)
	}
}

func sendString(msg string, conn *net.UDPConn, addr *net.UDPAddr) {
	send([]byte(msg), conn, addr)
}

// NewConnection …
func NewConnection(conn *net.UDPConn) {
	buffer := make([]byte, 4096)
	length, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	mt, err := NewMessageType(buffer[:2])
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mt.Type {
	case InputMessage:
		msg := message.Message{}
		if err = proto.Unmarshal(buffer[2:length], &msg); err != nil {
			fmt.Println(err)
			return
		}
		out := message.Message{
			Text:      fmt.Sprintf("got \"%s\"", msg.Text),
			Timestamp: time.Now().Unix(),
		}
		if data, err := proto.Marshal(&out); err != nil {
			fmt.Println(err)
		} else {
			send(append([]byte{'m', 0}, data...), conn, addr)
		}

	case InputListen:
		listen(conn, addr)
	case InputPing:
		send([]byte{'P', 0, '\n'}, conn, addr)
	case InputPong:
		send([]byte{'p', 0, '\n'}, conn, addr)
	}

}
