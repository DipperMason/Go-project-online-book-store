package http

import (
	"auth/internal/domain"
	"errors"
	"net/http"
	"serde"
)

type authenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authenticateResponse struct {
	Token string `json:"token"`
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) {
	authReq, err := serde.ReadJSONHttp[authenticateRequest](r.Body, w)
	if err != nil {
		return
	}

	tokenPair, err := h.auth.Authenticate(authReq.Email, authReq.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.UserNotFound) {
			status = http.StatusUnauthorized
		}

		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = serde.EncodeJSONHTTP[authenticateResponse](authenticateResponse{Token: tokenPair.Token}, w)
}
