package main

import (
	"fmt"
	"sync"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

// User …
type User struct {
	id string
}

var (
	mutex = &sync.Mutex{}
	users = make(map[string]*User)
)

// GetUser …
func GetUser(id string) (user *User, err error) {
	mutex.Lock()
	if result, ok := users[id]; ok {
		user = result
	} else {
		err = utils.NewError(fmt.Sprintf("cant find %s", id))
	}
	mutex.Unlock()

	return
}
