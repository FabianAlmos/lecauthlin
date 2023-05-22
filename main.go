package main

import (
	"auth/server"
)

func main() {
	cfg := server.NewConfig()

	server.Start(cfg)
}
