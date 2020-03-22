package spotifyclient

import (
	"github.com/zmb3/spotify"
)

type MockSpotifyClient struct {
	PlayerStateResponse struct {
		p *spotify.PlayerState
		e error
	}
	ShuffleResponse  error
	PlayResponse     error
	PauseResponse    error
	NextResponse     error
	PreviousResponse error
}

func (c *MockSpotifyClient) PlayerState() (*spotify.PlayerState, error) {
	return c.PlayerStateResponse.p, c.PlayerStateResponse.e
}
func (c *MockSpotifyClient) Shuffle(bool) error {
	return c.ShuffleResponse
}
func (c *MockSpotifyClient) Play() error {
	return c.PlayResponse
}
func (c *MockSpotifyClient) Pause() error {
	return c.PauseResponse
}
func (c *MockSpotifyClient) Next() error {
	return c.NextResponse
}
func (c *MockSpotifyClient) Previous() error {
	return c.PreviousResponse
}
