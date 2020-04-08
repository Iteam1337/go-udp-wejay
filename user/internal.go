package user

import (
	"bytes"
	"log"
	"time"

	"github.com/Iteam1337/go-udp-wejay/cover"
	"github.com/ankjevel/spotify"
)

func (u *User) playlistName(name string) string {
	return "[wejay] " + name
}

func (u *User) findPlaylist(n string) (playlist spotify.SimplePlaylist, err error) {
	name := u.playlistName(n)

	found := false
	if pl, e := u.client.CurrentUsersPlaylists(); e != nil {
		err = e
	} else {
		for _, pl := range pl.Playlists {
			if pl.Name != name {
				continue
			}

			if pl.IsPublic {
				err = u.client.ChangePlaylistAccess(pl.ID, false)
				if err != nil {
					log.Println("can't change access status", err)
					return
				}
			}

			if !pl.Collaborative {
				err = u.client.ChangePlaylistCollaborative(pl.ID, true)
				if err != nil {
					log.Println("can't change collaborative status", err)
					return
				}
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
		log.Println("can't get current user", e)
		err = e
		return
	}

	if pl, e := u.client.CreatePlaylistForUser(user.ID, name, "collaborative playlist for wejay", false); e != nil {
		err = e
	} else {
		err = u.client.ChangePlaylistCollaborative(pl.ID, true)
		if err != nil {
			log.Println("can't change collaborative status", err)
			return
		}
		playlist = pl.SimplePlaylist
	}

	if err != nil {
		return
	}

	go func() {
		cover, e := cover.Gen(n)
		if e != nil {
			log.Println("after get", e)
			err = e
			return
		}

		r := bytes.NewReader(cover)
		if e := u.client.SetPlaylistImage(playlist.ID, r); e != nil {
			log.Println("cant upload image", e, cover)
		}
	}()

	return
}

func (u *User) setDefaults() {
	if user, err := u.client.CurrentUser(); err != nil {
		log.Println("can't get user", err)
		return
	} else {
		u.ClientID = spotify.ID(user.ID)
	}
	u.playlist = spotify.SimplePlaylist{}
	u.active = false
}

func (u *User) getActiveDevice() (current spotify.PlayerDevice, ok bool) {
	ps, err := u.client.PlayerState()
	if err != nil {
		log.Println("can't get playerstate", err)
		return
	}

	if ps.Device.Active && !ps.Device.Restricted && ps.Device.ID != "" {
		ok = true
		current = ps.Device
		return
	}

	dev, err := u.client.PlayerDevices()
	if err != nil {
		log.Println("can't get devices", err)
		return
	}

	for _, device := range dev {
		if device.Restricted || device.ID == "" {
			continue
		}

		current = device
		ok = true
		break
	}

	return
}

func (u *User) handlePlayerState() {
	ps, err := u.client.PlayerState()
	if err != nil || ps == nil {
		return
	}

	if ps.ShuffleState {
		if err := u.client.Shuffle(false); err != nil {
			// log.Println("can't unshuffle", err)
			return
		}
	}

	if ps.RepeatState != "off" {
		if err := u.client.Repeat("off"); err != nil {
			// log.Println("can't set repeat state", err)
			return
		}
	}
}

func sleep() {
	time.Sleep(30 * time.Second)
}

func (u *User) loopState() {
	for {
		if u.client == nil || !u.active {
			break
		}

		if _, ok := u.getActiveDevice(); ok {
			go u.handlePlayerState()
		}

		sleep()
	}
}
