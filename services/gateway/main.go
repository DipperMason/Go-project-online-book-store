package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	AuthURL    = "http://auth:8001"
	ProfileURL = "http://profile:8003"
	CatalogURL = "http://catalog:8081"
	OrderURL   = "http://order:8002"
)

func newReverseProxy(target string) http.Handler {
	p := http.NewSingleHostReverseProxy(&http.URL{Scheme: "http", Host: target})
	return p
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(log.Writer(), next)
}

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	router := mux.NewRouter()

	authProxy := newReverseProxy("auth:8001")
	router.PathPrefix("/api/v1/auth").Handler(authProxy)

	profileProxy := newReverseProxy("profile:8003")
	router.PathPrefix("/api/v1/profile").Handler(profileProxy)

	catalogProxy := newReverseProxy("catalog:8081")
	router.PathPrefix("/api/v1/catalog").Handler(catalogProxy)

	orderProxy := newReverseProxy("order:8002")
	router.PathPrefix("/api/v1/order").Handler(orderProxy)

	router.HandleFunc("/health", healthCheck).Methods("GET")

	handler := loggingMiddleware(timeoutMiddleware(router))

	log.Println("API Gateway запущен на :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Ошибка запуска шлюза: %v", err)
	}
}