package main

import (
	"fmt"
	"sync"
)

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
		err = fmt.Errorf("cant find %s", id)
	}
	mutex.Unlock()

	return
}
