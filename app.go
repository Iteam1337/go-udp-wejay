package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Iteam1337/go-udp-wejay/rooms"
	"github.com/Iteam1337/go-udp-wejay/users"
	"github.com/Iteam1337/go-udp-wejay/utils"
)

func main() {
	if utils.GetEnv("STORE_STATE", "1") == "1" {
		users.LoadState()
		rooms.Restore()

		signalChannel := make(chan os.Signal, 2)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGSEGV)

		go func() {
			<-signalChannel
			users.SaveState()
			os.Exit(0)
		}()
	}

	listen(utils.GetEnv("ADDR", ":8090"))
}
