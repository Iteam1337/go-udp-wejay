package inputtype

import (
	"fmt"

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
		mt = types.ICallbackURL
	case byte(types.IUserExists):
		mt = types.IUserExists
	case byte(types.INewUser):
		mt = types.INewUser
	case byte(types.IDeleteUser):
		mt = types.IDeleteUser
	case byte(types.IPing):
		mt = types.IPing
	case byte(types.IRoomExists):
		mt = types.IRoomExists
	case byte(types.IJoinRoom):
		mt = types.IJoinRoom
	case byte(types.ILeaveRoom):
		mt = types.ILeaveRoom
	case byte(types.IQueryRooms):
		mt = types.IQueryRooms
	case byte(types.IUserLeaveRoom):
		mt = types.IUserLeaveRoom
	case byte(types.IUserRoom):
		mt = types.IUserRoom
	default:
		err = fmt.Errorf("unkown type: %s", string(buf[0]))
	}
	return
}
