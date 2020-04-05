package users

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

func (u *Users) SaveState() {
	path := utils.GetEnv("SAVE_STATE_LOCATION", "/tmp/wejay")

	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}

	write := func(buffer []byte) {
		if _, err := f.Write(buffer); err != nil {
			log.Println(err)
		}

		if err := f.Sync(); err != nil {
			log.Println(err)
		}
	}

	defer f.Close()

	for _, user := range u.users {
		state, err := user.GetSaveState()
		if err != nil {
			log.Println(err)
			continue
		}

		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(state)))

		write(buf)
		write(state)
	}
}
