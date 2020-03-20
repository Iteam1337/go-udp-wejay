package main

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/golang/protobuf/proto"
)

func parse(conn *net.UDPConn) (addr *net.UDPAddr, it types.InputType, buffer []byte, err error) {
	var length int
	buffer = make([]byte, 4096)
	length, addr, err = conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return
	}
	it, err = inputType(buffer[:2])
	if err != nil {
		log.Println(err)
		utils.SendEmpty(conn, addr)
		return
	}
	buffer = buffer[2:length]
	return
}

func read(msg proto.Message, buf []byte) (err error) {
	if err = proto.Unmarshal(buf[:], msg); err != nil {
		log.Println(err)
		return
	}

	return
}

func convert(it types.InputType, buf []byte) (pb proto.Message) {
	pb = it.Message()
	read(pb, buf)
	return
}

func newConnection(conn *net.UDPConn) {
	addr, it, buffer, err := parse(conn)
	if err != nil {
		log.Println(err)
		return
	}

	switch it {
	case types.IAction:
		msg := convert(it, buffer).(*message.Action)
		log.Println("IAction", msg)
		utils.SendM(it.Inv(), &message.ActionResponse{
			UserId: msg.UserId,
			Ok:     true,
		}, conn, addr)

	case types.IUserExists:
		msg := convert(it, buffer).(*message.UserExists)
		log.Println("IUserExists", msg)
		utils.SendM(it.Inv(), &message.UserExistsResponse{
			UserId: msg.UserId,
			Exists: false,
			Ok:     true,
		}, conn, addr)

	case types.INewUser:
		msg := convert(it, buffer).(*message.NewUser)
		log.Println("INewUser", msg)
		utils.SendM(it.Inv(), &message.NewUserResponse{
			UserId: msg.UserId,
			Ok:     true,
		}, conn, addr)

	case types.ICallbackURL:
		msg := convert(it, buffer).(*message.CallbackURL)
		log.Println("ICallbackURL", msg)
		utils.SendM(it.Inv(), &message.CallbackURLResponse{
			UserId: msg.UserId,
			Url:    spotifyAuth.AuthURL(msg.UserId),
			Ok:     true,
		}, conn, addr)

	case types.IPing:
		utils.SendM(it.Inv(), &message.Pong{Int: 0}, conn, addr)

	default:
		utils.SendEmpty(conn, addr)
	}

}
