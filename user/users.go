package user

import (
	"fmt"
)

var (
	users = make(map[string]*User)
)

// GetUser …
func GetUser(id string) (user *User, err error) {
	if result, ok := users[id]; ok {
		user = result
	} else {
		err = fmt.Errorf("cant find %s", id)
	}

	return
}

// NewUser …
func NewUser(id string, code string) {
	token, err := spotifyAuth.Exchange(code)
	if err != nil {
		return
	}

	user := User{id: id}
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
