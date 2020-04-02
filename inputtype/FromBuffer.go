package inputtype

import (
	"fmt"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-protobuf-wejay/version"
)

func FromBuffer(buf []byte) (mt types.InputType, err error) {
	if len(buf) != 2 {
		err = fmt.Errorf("wrong buffer length\n expected: 2, got: %d", len(buf))
		return
	}

	v := int8(buf[1])
	vv := int8(version.Version)

	if v != vv {
		err = fmt.Errorf("wrong message version\n expected: %d, got: %d", vv, v)
		return
	}

	switch buf[0] {
	case byte(types.ICallbackURL):
		log.Println("new request: CallbackURL")
		mt = types.ICallbackURL
	case byte(types.IUserExists):
		log.Println("new request: UserExists")
		mt = types.IUserExists
	case byte(types.INewUser):
		log.Println("new request: NewUser")
		mt = types.INewUser
	case byte(types.IDeleteUser):
		log.Println("new request: DeleteUser")
		mt = types.IDeleteUser
	case byte(types.IPing):
		log.Println("new request: Ping")
		mt = types.IPing
	case byte(types.IRoomExists):
		log.Println("new request: RoomExists")
		mt = types.IRoomExists
	case byte(types.IJoinRoom):
		log.Println("new request: JoinRoom")
		mt = types.IJoinRoom
	case byte(types.ILeaveRoom):
		log.Println("new request: LeaveRoom")
		mt = types.ILeaveRoom
	case byte(types.IQueryRooms):
		log.Println("new request: QueryRooms")
		mt = types.IQueryRooms
	case byte(types.IUserLeaveRoom):
		log.Println("new request: UserLeaveRoom")
		mt = types.IUserLeaveRoom
	default:
		err = fmt.Errorf("unkown type: %s", string(buf[0]))
	}
	return
}
