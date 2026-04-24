package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"profile/internal/app"
)

func main() {
	port := getEnv("PORT", "8003")

	application, err := app.New()
	if err != nil {
		log.Fatalf("app init error: %v", err)
	}
	defer func() {
		if closeErr := application.Close(); closeErr != nil {
			log.Printf("app close error: %v", closeErr)
		}
	}()

	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting Profile Service on %s", addr)

	if err = http.ListenAndServe(addr, application.GetHandler()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
