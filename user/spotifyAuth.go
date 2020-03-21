package user

import (
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/zmb3/spotify"
)

var spotifyAuth = spotify.NewAuthenticator(
	utils.GetEnv("REDIRECT_URL", "http://localhost:8080/callback"),
	spotify.ScopeUserReadCurrentlyPlaying,
	spotify.ScopeUserReadPlaybackState,
	spotify.ScopeUserModifyPlaybackState,
)

// AuthURL …
func AuthURL(id string) string {
	return spotifyAuth.AuthURL(id)
}
