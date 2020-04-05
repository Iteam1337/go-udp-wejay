package users

import (
	"fmt"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/user"
)

type Users struct {
	users       map[string]*user.User
	spotifyauth spotifyauth.Interface
}

func (u *Users) restore(pb message.RefUserSave) {
	user, err := user.Restore(pb)
	if err != nil {
		log.Println(err)
		return
	}

	if _, exists := u.users[pb.Id]; !exists {
		u.users[pb.Id] = &user
	}

	if user.Room == "" {
		return
	}
}

func (u *Users) GetAll() map[string]*user.User {
	return u.users
}

func (u *Users) GetUser(id string) (user *user.User, err error) {
	if result, ok := u.users[id]; ok {
		user = result
	} else {
		err = fmt.Errorf("can't find %s", id)
	}
	return
}

func (u *Users) New(id string, code string) {
	token, err := u.spotifyauth.Exchange(code)

	if err != nil {
		return
	}

	user := user.New(id)
	user.SetClient(token)

	if _, exists := u.users[id]; !exists {
		u.users[id] = &user
	}
}

func (u Users) Exists(id string) bool {
	_, ok := u.users[id]

	return ok
}

func (u *Users) Delete(id string) {
	if user, ok := u.users[id]; ok {
		user.Destroy()
		delete(u.users, id)
	}
}

var (
	users = Users{
		users:       make(map[string]*user.User),
		spotifyauth: spotifyauth.Struct,
	}
	GetUser   = users.GetUser
	GetAll    = users.GetAll
	New       = users.New
	Exists    = users.Exists
	Delete    = users.Delete
	SaveState = users.SaveState
	LoadState = users.LoadState
)
