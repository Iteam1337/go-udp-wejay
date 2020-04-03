package user

import (
	"log"
	"time"

	"github.com/ankjevel/spotify"
)

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

	if ps.Device.Active && !ps.Device.Restricted {
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
		if device.Restricted {
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
	if err != nil {
		return
	}

	if ps.ShuffleState {
		if err := u.client.Shuffle(false); err != nil {
			log.Println("can't unshuffle", err)
		}
	}
}

func (u *User) setContext(device spotify.PlayerDevice) {
	if err := u.client.PlayOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
	}); err != nil {
		log.Println("can't set context", err)
	}
}

func sleep() {
	time.Sleep(15 * time.Second)
}

func (u *User) loopState() {
	for {
		if u.client == nil {
			break
		}

		if !u.active {
			break
		}

		device, ok := u.getActiveDevice()
		if !ok {
			sleep()
			continue
		}

		u.setContext(device)
		go u.handlePlayerState()

		log.Println("do things")
		sleep()
	}
}
