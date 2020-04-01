package user

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

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
	item := spotify.FullTrack{}
	item.URI = "uri://"
	p := spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Timestamp: 0,
			Progress:  1337,
			Playing:   true,
			Item:      &item,
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

func Test_listenQueriesPlayerstateOnFixedIntervals(t *testing.T) {
	listen := make(chan ListenMsg, math.MaxInt16)
	now := make(chan time.Time, math.MaxInt16)
	playerState := make(chan spotify.PlayerState, 3)
	close := make(chan bool, 1)
	u := New("x")
	u.client = &spotify.Client{}

	// setup fake timer
	go func() {
		hour := 13
		minute := 4

		for {
			// timer mock is fragile, dont change sec or nanosec
			now <- time.Date(2020, time.March, 23, hour, minute, 13, 37, time.UTC)
			if minute+5 > 60 {
				minute = minute % 60
				hour = hour + 1
			} else {
				minute = minute + 5
			}
		}
	}()
	// setup listen state
	go func() {
		item := spotify.FullTrack{}
		item.URI = "uri://"

		ps := spotify.PlayerState{
			CurrentlyPlaying: spotify.CurrentlyPlaying{
				Progress: 1337,
				Playing:  true,
				PlaybackContext: spotify.PlaybackContext{
					URI: "uri://",
				},
				Item: &item,
			},
			Device: spotify.PlayerDevice{
				Active: true,
			},
		}

		playerState <- ps
		ps.CurrentlyPlaying.Progress = ps.CurrentlyPlaying.Progress + int(10*time.Second)
		playerState <- ps
		ps.CurrentlyPlaying.Progress = ps.CurrentlyPlaying.Progress + int(10*time.Second)
		playerState <- ps
	}()

	// monkey patches
	playerStatesSent := 0
	var d *spotify.Client
	monkey.Patch(time.Now, func() time.Time { return <-now })
	monkey.PatchInstanceMethod(reflect.TypeOf(d), "PlayerState", func(*spotify.Client) (ps *spotify.PlayerState, e error) {
		state := <-playerState
		playerStatesSent = playerStatesSent + 1
		if playerStatesSent < 3 {
			ps = &state
		} else {
			e = errors.New("only send 3 player states")
			u.Destroy()
		}
		return
	})
	defer monkey.UnpatchAll()

	u.SetListen(&listen, &close)

	var messages []ListenMsg
	go func() {
		for {
			msg := <-listen
			messages = append(messages, msg)
		}
	}()

	select {
	case <-time.After(10 * time.Second):
		t.Fatal("Test should not have taken more than ~4 seconds (one per iteration)")
	case <-close:
	}

	if len(messages) != 5 {
		t.Error("[len(messages)] expected 5 messages, got", len(messages))
	}

	last := messages[len(messages)-1]
	expectedByte := byte(message.ListenResponse_PROGRESS)
	if last.Meta[0] != expectedByte {
		t.Error("[ListenMsg.Meta[0]] expected", expectedByte, "got", last.Meta[0])
	}

	for i, msg := range messages {
		if msg.Type != message.ListenResponse_PLAYBACK {
			t.Error(i, "[ListenMsg.Type] expected", message.ListenResponse_PLAYBACK, "got", msg.Type)
		}
	}
}
