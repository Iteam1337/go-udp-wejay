package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/zmb3/spotify"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

var (
	redirectURL = utils.Getenv("REDIRECT_URL", "http://localhost:8080/callback")
	spotifyAuth = spotify.NewAuthenticator(
		redirectURL,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
	)
)

func main() {
	Listen(utils.Getenv("ADDR", ":8090"))
}
