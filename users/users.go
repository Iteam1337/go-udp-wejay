package users

import (
	"fmt"

	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/user"
)

var users = make(map[string]*user.User)

// GetUser …
func GetUser(id string) (user *user.User, err error) {
	if result, ok := users[id]; ok {
		user = result
	} else {
		err = fmt.Errorf("cant find %s", id)
	}

	return
}

// New …
func New(id string, code string) {
	token, err := spotifyauth.Exchange(code)
	if err != nil {
		return
	}

	user := user.New(id)
	user.SetClient(token)

	if _, exists := users[id]; !exists {
		users[id] = &user
	}
}

// Exists …
func Exists(id string) bool {
	_, ok := users[id]
	return ok
}

// Delete …
func Delete(id string) {
	delete(users, id)
}
