package spotifyauth

import (
	"github.com/Iteam1337/go-udp-wejay/utils"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Interface …
type Interface interface {
	AuthURL(id string) string
	NewClient(*oauth2.Token) spotify.Client
	Exchange(string) (*oauth2.Token, error)
}

// SpotifyAuth …
type SpotifyAuth struct {
	auth spotify.Authenticator
}

// AuthURL …
func (s SpotifyAuth) AuthURL(id string) string {
	return s.auth.AuthURL(id)
}

// NewClient …
func (s SpotifyAuth) NewClient(token *oauth2.Token) spotify.Client {
	return s.auth.NewClient(token)
}

// Exchange …
func (s SpotifyAuth) Exchange(code string) (*oauth2.Token, error) {
	return s.auth.Exchange(code)
}

// Global exports
var (
	Struct = SpotifyAuth{
		spotify.NewAuthenticator(
			utils.GetEnv("REDIRECT_URL", "http://localhost:8080/callback"),
			spotify.ScopeUserReadCurrentlyPlaying,
			spotify.ScopeUserReadPlaybackState,
			spotify.ScopeUserModifyPlaybackState,
		),
	}
	AuthURL   = Struct.AuthURL
	NewClient = Struct.NewClient
	Exchange  = Struct.Exchange
)
