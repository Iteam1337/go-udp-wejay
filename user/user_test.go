package user

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"golang.org/x/oauth2"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
	"github.com/zmb3/spotify"
)

func Test_canCreateUser(t *testing.T) {
	u := New("x")

	if u.id != "x" {
		t.Error("user was not created", u)
	}
}

func Test_setClient(t *testing.T) {
	u := New("x")

	if u.id != "x" {
		t.Error("wrong id")
		return
	}

	c := spotify.Client{}
	p := spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Timestamp: 0,
			Progress:  1337,
			Playing:   true,
			Item:      &spotify.FullTrack{},
			PlaybackContext: spotify.PlaybackContext{
				URI: "uri://",
			},
		},
		Device: spotify.PlayerDevice{
			Active: true,
		},
	}

	var a spotifyauth.SpotifyAuth
	var d *spotify.Client

	monkey.PatchInstanceMethod(reflect.TypeOf(a), "NewClient", func(spotifyauth.SpotifyAuth, *oauth2.Token) spotify.Client {
		return c
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(d), "PlayerState", func(*spotify.Client) (ps *spotify.PlayerState, e error) {
		ps = &p
		return
	})

	token := oauth2.Token{}

	token.AccessToken = "AccessToken"
	token.TokenType = "Bearer"
	token.RefreshToken = "RefreshToken"

	u.SetClient(&token)

	if u.active != true {
		t.Error("active not set\n", u.active)
	}

	if u.listening != true {
		t.Error("listening not set\n", u.listening)
	}

	if u.progress != 1337 {
		t.Error("progress not set\n", u.progress)
	}

	if u.current != "uri://" {
		t.Error("current not set\n", u.current)
	}

	defer monkey.UnpatchAll()
}

func Test_runActions(t *testing.T) {
	u := New("x")
	u.client = &spotify.Client{}

	var d *spotify.Client
	for _, s := range []struct {
		t message.Action_ActionType
		m string
	}{
		{message.Action_PLAY, "Play"},
		{message.Action_PAUSE, "Pause"},
		{message.Action_NEXT, "Next"},
		{message.Action_PREVIOUS, "Previous"},
	} {

		err := fmt.Errorf("%s called", s.m)

		p := monkey.PatchInstanceMethod(reflect.TypeOf(d), s.m, func(*spotify.Client) error {
			return err
		})

		if e := u.RunAction(s.t); e != err {
			t.Error(e)
		}

		defer p.Restore()
	}
}

func Test_runActionShuffle(t *testing.T) {
	u := New("x")
	u.client = &spotify.Client{}
	p := spotify.PlayerState{
		ShuffleState: false,
	}
	err := errors.New("Shuffle called last")
	res := make(chan bool, 1)

	var d *spotify.Client
	monkey.PatchInstanceMethod(reflect.TypeOf(d), "Shuffle", func(c *spotify.Client, state bool) error {
		res <- state
		return err
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(d), "PlayerState", func(*spotify.Client) (ps *spotify.PlayerState, e error) {
		ps = &p
		return
	})

	defer monkey.UnpatchAll()

	if e := u.RunAction(message.Action_SHUFFLE); e != err {
		t.Error(e)
		return
	}

	if true != <-res {
		t.Error("new shuffle state should have been `true`")
		return
	}
}
