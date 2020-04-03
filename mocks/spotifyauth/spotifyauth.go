package spotifyauth

import (
	"github.com/ankjevel/spotify"
	"golang.org/x/oauth2"
)

type MockSpotifyAuth struct {
	AuthURLResponse   string
	NewClientResponse spotify.Client
	ExchangeResponse  struct {
		t *oauth2.Token
		e error
	}
}

func (s MockSpotifyAuth) AuthURL(id string) string {
	return s.AuthURLResponse
}

func (s MockSpotifyAuth) NewClient(token *oauth2.Token) spotify.Client {
	return s.NewClientResponse
}

func (s MockSpotifyAuth) Exchange(code string) (*oauth2.Token, error) {
	return s.ExchangeResponse.t, s.ExchangeResponse.e
}
