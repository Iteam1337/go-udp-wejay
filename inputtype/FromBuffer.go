package inputtype

import (
	"fmt"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-protobuf-wejay/version"
)

// FromBuffer …
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
	case byte(types.IUserExists):
		log.Println("new request: UserExists")
		mt = types.IUserExists
	case byte(types.IAction):
		log.Println("new request: Action")
		mt = types.IAction
	case byte(types.INewUser):
		log.Println("new request: NewUser")
		mt = types.INewUser
	case byte(types.ICallbackURL):
		log.Println("new request: CallbackURL")
		mt = types.ICallbackURL
	case byte(types.IPing):
		log.Println("new request: Ping")
		mt = types.IPing
	case byte(types.INowPlaying):
		log.Println("new request: NowPlaying")
		mt = types.INowPlaying
	case byte(types.IListen):
		log.Println("new request: Listen")
		mt = types.IListen
	case byte(types.IDeleteUser):
		log.Println("new request: DeleteUser")
		mt = types.IDeleteUser
	default:
		err = fmt.Errorf("unkown type: %s", string(buf[0]))
	}
	return
}
