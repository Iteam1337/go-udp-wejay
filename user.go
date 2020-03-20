package main

import (
	"log"
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// User …
type User struct {
	client *spotify.Client
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

	return
}

// NowPlaying …
func (u User) NowPlaying() (artist string, track string) {
	client := u.getClient()

	currentlyPlaying, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Println(err.Error())
		return
	}

	item := currentlyPlaying.Item
	if item == nil {
		return
	}

	var artists []string
	for _, key := range item.Artists {
		artists = append(artists, key.Name)
	}

	artist = strings.Join(artists, ", ")
	track = item.SimpleTrack.Name
	return
}
