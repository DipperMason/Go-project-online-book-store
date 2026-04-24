package http

import (
	"auth/internal/domain"
	"errors"
	"net/http"
	"serde"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	registerReq, err := serde.ReadJSONHttp[registerRequest](r.Body, w)
	if err != nil {
		return
	}

	err = h.auth.Register(registerReq.Email, registerReq.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.UserAlreadyExists) {
			status = http.StatusConflict
		}

		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = serde.EncodeJSONHTTP[registerResponse](registerResponse{Status: "ok", Details: "registration successful"}, w)
}
