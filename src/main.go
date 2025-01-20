package main

import (
	"log"

	"github.com/joho/godotenv"
)

var FALLBACK_LISTEN_ADDR = "localhost:8080"

func main() {
	envFile, _ := godotenv.Read("../.env")
	listenAddr, addrValid := envFile["LISTEN_ADDR"]

	if !addrValid {
		listenAddr = FALLBACK_LISTEN_ADDR
		log.Printf("Couldn't find address for service in .env, falling back to '%s'", FALLBACK_LISTEN_ADDR)
	}

	app := App{}
	app.Initialize()
	app.Run(listenAddr)
}
