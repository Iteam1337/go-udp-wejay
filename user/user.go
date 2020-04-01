package user

import (
	"log"
	"strconv"
	"time"

	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/Iteam1337/go-protobuf-wejay/message"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// ListenMsg …
type ListenMsg struct {
	Type message.ListenResponse_Change
	Meta []byte
	Ok   bool
}

// User …
type User struct {
	id        string
	client    *spotify.Client
	listen    *chan ListenMsg
	listening bool
	active    bool
	progress  int
	current   string
}

func (u *User) getClient() *spotify.Client {
	return u.client
}

func (u *User) toggleShuffleState() (state bool, err error) {
	client := u.client

	playerState, err := client.PlayerState()

	if err != nil {
		return
	}

	state = !playerState.ShuffleState
	playerState.ShuffleState = state

	return
}

// Destroy …
func (u *User) Destroy() {
	utils.SetNil(&u.id)
	utils.SetNil(&u.client)
	utils.SetNil(&u.listen)
	utils.SetNil(&u.listening)
	utils.SetNil(&u.active)
	utils.SetNil(&u.progress)
	utils.SetNil(&u.current)
}

// SetClient …
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

// SetListen …
func (u *User) SetListen(listen *chan ListenMsg, close *chan bool) {
	exists := u.listen != nil
	u.listen = listen

	if exists {
		return
	}

	go func() {
		for {
			if u.client == nil {
				*close <- true
				break
			}
			ps, err := u.client.PlayerState()
			if err != nil {
				log.Println(err)
				*close <- true
				break
			}

			m := ListenMsg{Type: message.ListenResponse_PLAYBACK, Ok: true}

			listening := ps.Playing
			if u.listening != listening {
				m.Meta = []byte{byte(message.ListenResponse_LISTENING), utils.BoolToByte(listening)}
				if u.listen != nil {
					*u.listen <- m
				}

				u.listening = listening
			}

			active := ps.Device.Active
			if u.active != active {
				m.Meta = []byte{byte(message.ListenResponse_ACTIVE), utils.BoolToByte(active)}
				if u.listen != nil {
					*u.listen <- m
				}
				u.active = active
			}

			progress := ps.Progress
			if u.active && u.progress != progress && u.listen != nil {
				var buf []byte
				buf = append(buf, byte(message.ListenResponse_PROGRESS))
				buf = append(buf, []byte(strconv.Itoa(progress))...)
				m.Meta = buf
				if u.listen != nil {
					*u.listen <- m
				}
				u.progress = progress
			}

			if ps.CurrentlyPlaying.Item != nil {
				current := string(ps.CurrentlyPlaying.Item.URI)
				if u.current != current && u.listen != nil {
					var buf []byte
					buf = append(buf, byte(message.ListenResponse_CURRENT))
					buf = append(buf, []byte(current)...)
					m.Meta = buf
					if u.listen != nil {
						*u.listen <- m
					}
					u.current = current
				}
			}

			if u.listen == nil {
				*close <- true
				break
			}

			waitTime := 5 * time.Second
			now := time.Now()

			time.Sleep(utils.RoundToNearestSecond(now, waitTime).Sub(now))
		}
	}()
}

// New …
func New(id string) (u User) {
	u.id = id
	return
}
