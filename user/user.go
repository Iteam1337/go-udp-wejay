package user

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
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

// SetClient …
func (u *User) SetClient(token *oauth2.Token) {
	client := spotifyauth.NewClient(token)

	u.client = &client
}

// RunAction …
func (u *User) RunAction(action message.Action_ActionType) (err error) {
	client := u.getClient()
	switch action {
	case message.Action_PLAY:
		err = client.Play()
	case message.Action_PAUSE:
		err = client.Pause()
	case message.Action_NEXT:
		err = client.Next()
	case message.Action_PREVIOUS:
		err = client.Previous()
	case message.Action_SHUFFLE:
		state, maybeErr := u.toggleShuffleState()
		if maybeErr != nil {
			err = maybeErr
		} else {
			err = client.Shuffle(state)
		}
	}

	if err != nil {
		log.Print(err)
	}

	if u.listen != nil {
		*u.listen <- ListenMsg{
			message.ListenResponse_ACTION,
			[]byte{byte(action)},
			err == nil,
		}
	}

	return
}

func (u *User) boolToByte(b bool) byte {
	if b {
		return byte(1)
	}

	return byte(0)
}

// SetListen …
func (u *User) SetListen(listen *chan ListenMsg) {
	u.listen = listen

	go func() {
		for {
			ps, err := u.client.PlayerState()
			if err != nil {
				log.Println(err)
				break
			}

			m := ListenMsg{Type: message.ListenResponse_PLAYBACK, Ok: true}

			listening := ps.Playing
			if u.listening != listening {
				m.Meta = []byte{byte(message.ListenResponse_LISTENING), u.boolToByte(listening)}
				*u.listen <- m
				u.listening = listening
			}

			active := ps.Device.Active
			if u.active != active {
				m.Meta = []byte{byte(message.ListenResponse_ACTIVE), u.boolToByte(active)}
				*u.listen <- m
				u.active = active
			}

			progress := ps.Progress
			if u.active && u.progress != progress {
				var buf []byte
				buf = append(buf, byte(message.ListenResponse_PROGRESS))
				buf = append(buf, []byte(strconv.Itoa(progress))...)
				m.Meta = buf
				*u.listen <- m
				u.progress = progress
			}

			current := string(ps.PlaybackContext.URI)
			if u.current != current {
				var buf []byte
				buf = append(buf, byte(message.ListenResponse_CURRENT))
				buf = append(buf, []byte(current)...)
				m.Meta = buf
				*u.listen <- m
				u.current = current
			}

			time.Sleep(10 * time.Second)
		}
	}()
}

// NowPlaying …
func (u *User) NowPlaying() (track message.Track, e error) {
	client := u.getClient()
	currentlyPlaying, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		e = err
		return
	}

	item := currentlyPlaying.Item
	if item == nil {
		e = errors.New("could not get current item")
		return
	}

	track.Duration = int64(item.Duration)
	track.Id = item.ID.String()
	track.Name = item.Name
	track.Uri = string(item.URI)

	var artists []*message.Artist
	for _, key := range item.Artists {
		var artist message.Artist
		artist.Id = key.ID.String()
		artist.Name = key.Name
		artist.Uri = string(key.URI)
		artists = append(artists, &artist)
	}
	track.Artists = artists
	return
}

// New …
func New(id string) (u User) {
	u.id = id
	return
}
