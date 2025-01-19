package main

import "github.com/joho/godotenv"

func main() {
	envFile, _ := godotenv.Read("../.env")
	listenAddr, addrValid := envFile["LISTEN_ADDR"]

	if !addrValid {
		panic("No address set for service")
	}

	app := App{}
	app.Initialize()
	app.Run(listenAddr)
}
