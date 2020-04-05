package user

import (
	"encoding/json"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/golang/protobuf/proto"
	"golang.org/x/oauth2"
)

func (u *User) GetSaveState() (bytes []byte, err error) {
	var token *oauth2.Token
	var client, playlist []byte

	token, err = u.client.Token()
	if err != nil {
		return
	}

	client, err = json.Marshal(token)
	if err != nil {
		return
	}

	if u.playlist.URI != "" {
		playlist, err = json.Marshal(u.playlist)
		if err != nil {
			return
		}
	}

	return proto.Marshal(&message.RefUserSave{
		Id:            u.ID,
		Client:        client,
		Active:        u.active,
		Playlist:      playlist,
		PlaylistOwner: u.playlistOwner,
		ClientId:      string(u.ClientID),
		Room:          u.Room,
	})
}
