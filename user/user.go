package user

import (
	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type User struct {
	id        string
	client    *spotify.Client
	listening bool
	active    bool
	progress  int
	current   string
	Room      string
}

func (u *User) GetClient() *spotify.Client {
	return u.client
}

func (u *User) Destroy() {
	utils.SetNil(&u.id)
	utils.SetNil(&u.client)
	utils.SetNil(&u.listening)
	utils.SetNil(&u.active)
	utils.SetNil(&u.progress)
	utils.SetNil(&u.current)
	utils.SetNil(&u.Room)
}

func (u *User) SetClient(token *oauth2.Token) {
	client := spotifyauth.NewClient(token)
	u.client = &client

	ps, err := client.PlayerState()
	if err != nil {
		return
	}

	u.listening = ps.Playing
	u.active = ps.Device.Active
	u.progress = ps.Progress
	if ps.CurrentlyPlaying.Item != nil {
		u.current = string(ps.CurrentlyPlaying.Item.URI)
	}
}

func New(id string) (u User) {
	u.id = id
	return
}
