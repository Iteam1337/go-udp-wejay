package user

import (
	"log"

	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type User struct {
	id          string
	client      *spotify.Client
	clientID    spotify.ID
	listening   bool
	active      bool
	progress    int
	currentItem string
	playlist    spotify.SimplePlaylist
	Room        string
}

func (u *User) playlistName(name string) string {
	return "[wejay] " + name
}

func (u *User) findPlaylist(name string) (playlist spotify.SimplePlaylist, err error) {
	name = u.playlistName(name)

	found := false
	if pl, e := u.client.CurrentUsersPlaylists(); e != nil {
		err = e
	} else {
		for _, pl := range pl.Playlists {
			if pl.Name != name {
				continue
			}

			playlist = pl
			found = true
			break
		}
	}

	if found {
		return
	}

	user, e := u.client.CurrentUser()
	if e != nil {
		err = e
		return
	}

	if pl, e := u.client.CreatePlaylistForUser(user.ID, name, "collaborative playlist for wejay", true); e != nil {
		err = e
	} else {
		playlist = pl.SimplePlaylist
	}

	return
}

func (u *User) setDefaults() {
	ps, err := u.client.PlayerState()
	if err != nil {
		log.Println(err)
		return
	}

	if user, err := u.client.CurrentUser(); err != nil {
		log.Println(err)
		return
	} else {
		u.clientID = spotify.ID(user.ID)
	}

	u.listening = ps.Playing
	u.active = ps.Device.Active
	u.progress = ps.Progress
	if ps.CurrentlyPlaying.Item != nil {
		u.currentItem = string(ps.CurrentlyPlaying.Item.URI)
	}
}

func (u *User) GetClient() *spotify.Client {
	return u.client
}

func (u *User) Destroy() {
	utils.SetNil(&u.id)
	utils.SetNil(&u.client)
	utils.SetNil(&u.clientID)
	utils.SetNil(&u.listening)
	utils.SetNil(&u.active)
	utils.SetNil(&u.progress)
	utils.SetNil(&u.currentItem)
	utils.SetNil(&u.playlist)
	utils.SetNil(&u.Room)
}

func (u *User) JoinRoom(name string) {
	if u.Room == name {
		return
	}

	u.Room = name

	if pl, err := u.findPlaylist(name); err != nil {
		log.Println(err)
	} else {
		u.playlist = pl
	}
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

func New(id string) (u User) {
	u.id = id
	return
}
