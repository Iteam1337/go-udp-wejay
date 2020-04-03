package user

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"golang.org/x/oauth2"

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
	}

	pu := spotify.PrivateUser{
		User: spotify.User{
			ID: "hey",
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

	monkey.PatchInstanceMethod(reflect.TypeOf(d), "CurrentUser", func(*spotify.Client) (res *spotify.PrivateUser, e error) {
		res = &pu
		return
	})

	token := oauth2.Token{}

	token.AccessToken = "AccessToken"
	token.TokenType = "Bearer"
	token.RefreshToken = "RefreshToken"

	u.SetClient(&token)

	if u.active != false {
		t.Error("active not set\n", u.active)
	}

	if u.clientID != "hey" {
		t.Error("clientID not set\n", u.clientID)
	}

	defer monkey.UnpatchAll()
}
