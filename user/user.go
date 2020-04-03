package user

import (
	"log"

	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type User struct {
	id       string
	client   *spotify.Client
	clientID spotify.ID
	active   bool
	playlist spotify.SimplePlaylist
	Room     string
}

func (u *User) GetClient() *spotify.Client {
	return u.client
}

func (u *User) Destroy() {
	utils.SetNil(&u.id)
	utils.SetNil(&u.client)
	utils.SetNil(&u.clientID)
	utils.SetNil(&u.active)
	utils.SetNil(&u.playlist)
	utils.SetNil(&u.Room)
}

func (u *User) JoinRoom(name string, playlist spotify.SimplePlaylist) {
	if u.playlist.ID != playlist.ID {
		u.playlist = playlist
	}

	if u.Room == name {
		return
	}

	u.Room = name
	u.active = true

	go u.loopState()
}

func (u *User) LeaveRoom() {
	if u.Room == "" {
		return
	}

	if u.playlist.ID != "" {
		if err := u.client.UnfollowPlaylist(u.clientID, u.playlist.ID); err != nil {
			log.Println(err)
		}

		u.playlist = spotify.SimplePlaylist{}
	}

	u.Room = ""
}

func (u *User) SetClient(token *oauth2.Token) {
	client := spotifyauth.NewClient(token)
	u.client = &client

	defer u.setDefaults()
}

func (u *User) NewPlaylist(name string) (playlist spotify.SimplePlaylist, ok bool) {
	pl, err := u.findPlaylist(name)

	if err != nil {
		log.Println(err)
		return
	}

	ok = true
	playlist = pl
	return
}

func New(id string) (u User) {
	u.id = id
	return
}
