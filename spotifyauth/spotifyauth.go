package spotifyauth

import (
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var auth = spotify.NewAuthenticator(
	utils.GetEnv("REDIRECT_URL", "http://localhost:8080/callback"),
	spotify.ScopeUserReadCurrentlyPlaying,
	spotify.ScopeUserReadPlaybackState,
	spotify.ScopeUserModifyPlaybackState,
)

// AuthURL …
func AuthURL(id string) string {
	return auth.AuthURL(id)
}

// NewClient …
func NewClient(token *oauth2.Token) spotify.Client {
	return auth.NewClient(token)
}

// Exchange …
func Exchange(code string) (*oauth2.Token, error) {
	return auth.Exchange(code)
}
