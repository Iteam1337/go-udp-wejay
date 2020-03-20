package main

import (
	"fmt"
	"log"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

// InputType …
type InputType byte

// InputTypeEnum …
const (
	InputMessage InputType = 'm'
	InputListen  InputType = 'l'
	InputPing    InputType = 'p'
	InputPong    InputType = 'P'
)

// MessageType …
type MessageType struct {
	Type    InputType
	Version int8
}

// NewMessageType …
func NewMessageType(buff []byte) (mt MessageType, err error) {
	if len(buff) != 2 {
		err = utils.NewError(fmt.Sprintf("wrong buffer length\n expected: 2, got: %d\n", len(buff)))
		return
	}

	switch buff[0] {
	case 'm':
		mt = MessageType{InputMessage, int8(buff[1])}
	case 'l':
		mt = MessageType{InputListen, 0}
	case 'p':
		log.Println("Ping")
		mt = MessageType{InputPing, 0}
	case 'P':
		log.Println("Pong")
		mt = MessageType{InputPong, 0}
	default:
		err = utils.NewError(fmt.Sprintf("unkown type: %b\n", buff[0]))
	}
	return
}
