package http

import (
	"auth/internal/services"
	"middlewares"
	"net/http"
)

type Handler struct {
	auth *services.AuthService
}

func New(auth *services.AuthService) *Handler {
	return &Handler{
		auth: auth,
	}
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/register", h.register)
	mux.HandleFunc("POST /api/v1/auth", h.authenticate)
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	h.registerRoutes(mux)

	var handler http.Handler = mux
	handler = middlewares.WithLogging(handler)

	return handler
}
