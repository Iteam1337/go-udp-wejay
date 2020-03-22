package users

import (
	"testing"

	mock "github.com/Iteam1337/go-udp-wejay/mocks/spotifyauth"

	"github.com/Iteam1337/go-udp-wejay/user"
)

var auth = mock.MockSpotifyAuth{}

func Test_CreateNewUser(t *testing.T) {
	users := Users{
		users:       make(map[string]user.User),
		spotifyauth: auth,
	}

	users.New("user_id", "code")

	if !users.Exists("user_id") {
		t.Error("user could not be created")
	}
}

func Test_GetUser(t *testing.T) {
	var u *user.User
	var e error

	users := Users{
		users:       make(map[string]user.User),
		spotifyauth: auth,
	}

	users.New("user_id", "code")

	u, e = users.GetUser("user_id")

	if e != nil {
		t.Error("user could not be fetched", e.Error())
		return
	}

	if u == nil {
		t.Error("user not fetchedS")
	}
}

func Test_DestroyUser(t *testing.T) {
	users := Users{
		users:       make(map[string]user.User),
		spotifyauth: auth,
	}

	users.New("user_id", "code")

	if !users.Exists("user_id") {
		t.Error("user was never created")
		return
	}

	users.Delete("user_id")

	if users.Exists("user_id") {
		t.Error("user did not det deleted")
	}
}
