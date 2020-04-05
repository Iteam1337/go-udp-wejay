package users

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/golang/protobuf/proto"
)

func (u *Users) LoadState() {
	var pathIsOk bool

	path := utils.GetEnv("SAVE_STATE_LOCATION", "/tmp/wejay")
	// tmp, err := ioutil.ReadFile(path)
	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer func() {
		if pathIsOk {
			os.Remove(path)
		}
	}()
	defer f.Close()

	for {
		var pb message.RefUserSave

		buf := make([]byte, 4)
		if _, err := f.Read(buf); err != nil {
			break
		}
		size := binary.LittleEndian.Uint32(buf)
		msg := make([]byte, size)
		if _, err := f.Read(msg); err != nil {
			log.Println("msg read error", err)
			continue
		}

		if err := proto.Unmarshal(msg, &pb); err != nil {
			log.Println("unmarshal error", err)
			continue
		}

		u.restore(pb)

		pathIsOk = true
	}
}
