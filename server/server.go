package server

import (
	"auth/config"
	"auth/server/handler"
	"log"
	"net/http"
)

func Start(cfg *config.Config) {
	authHandler := handler.NewAuthHandler(cfg)
	userHandler := handler.NewUserHandler(cfg)

	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/profile", userHandler.GetProfile)

	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}
