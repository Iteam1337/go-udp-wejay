package user

import (
	"log"

	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/ankjevel/spotify"
	"golang.org/x/oauth2"
)

type User struct {
	id            string
	client        *spotify.Client
	active        bool
	playlist      spotify.SimplePlaylist
	playlistOwner bool
	ClientID      spotify.ID
	Room          string
}

func (u *User) GetClient() *spotify.Client {
	return u.client
}

func (u *User) Destroy() {
	utils.SetNil(&u.id)
	utils.SetNil(&u.client)
	utils.SetNil(&u.active)
	utils.SetNil(&u.playlist)
	utils.SetNil(&u.playlistOwner)
	utils.SetNil(&u.ClientID)
	utils.SetNil(&u.Room)
}

func (u *User) JoinRoom(name string, playlist spotify.SimplePlaylist, owner spotify.ID) {
	if u.playlist.ID != playlist.ID {
		u.playlist = playlist
	}

	playlistOwner := u.ClientID == owner
	u.playlistOwner = playlistOwner
	if !playlistOwner {
		if err := u.client.FollowPlaylist(owner, playlist.ID, true); err != nil {
			log.Println("follow failed", err)
			return
		}
	}

	u.Room = name
	if u.active {
		return
	}
	u.active = true
	go u.loopState()
}

func (u *User) LeaveRoom() {
	u.active = false

	if u.Room == "" {
		return
	}

	if u.playlist.ID != "" {
		if err := u.client.UnfollowPlaylist(u.ClientID, u.playlist.ID); err != nil {
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
