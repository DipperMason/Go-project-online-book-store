package app

import (
	"auth/internal/config"
	"auth/internal/repo"
	"auth/internal/services"
	httptransport "auth/internal/transport/http"
	"net/http"

	"jwt"
)

// App связывает зависимости auth-сервиса.
type App struct {
	handler http.Handler
}

func New() *App {
	cfg := config.NewConfig()

	userRepo := repo.New()
	jwtService := jwt.NewService(jwt.Config{
		Secret:         cfg.JWTSecret,
		AccessTokenTTL: cfg.AccessTokenTTL,
	})
	authService := services.NewAuthService(userRepo, jwtService)
	handler := httptransport.New(authService)

	return &App{handler: handler.Router()}
}

func (a *App) Handler() http.Handler {
	return a.handler
}
