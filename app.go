package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

func main() {
	listen(utils.GetEnv("ADDR", ":8090"))
}
