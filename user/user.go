package user

import (
	"errors"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// ListenMsg …
type ListenMsg struct {
	Type message.ListenResponse_ActionType
	Meta []byte
}

// User …
type User struct {
	id     string
	client *spotify.Client
	listen *chan ListenMsg
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
	client := spotifyAuth.NewClient(token)

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
		*u.listen <- ListenMsg{message.ListenResponse_STATE_CHANGE, []byte{byte(action)}}
	}

	return
}

// SetListen …
func (u *User) SetListen(listen *chan ListenMsg) {
	u.listen = listen
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
