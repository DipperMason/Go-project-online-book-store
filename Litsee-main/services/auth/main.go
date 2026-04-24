package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	application := app.New()

	log.Printf("starting Auth Service on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, application.Handler()); err != nil {
		log.Fatalf("auth server error: %v", err)
	}
}
