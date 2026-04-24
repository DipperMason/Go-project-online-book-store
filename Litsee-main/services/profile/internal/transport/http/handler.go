package http

import (
	"encoding/json"
	"jwt"
	"net/http"
	"profile/internal/domain"
	"profile/internal/services"
	"strconv"

	"middlewares"
)

// Handler обработчик HTTP запросов для профиля.
type Handler struct {
	profileService *services.ProfileService
	jwtService     jwt.Service
}

func New(profileService *services.ProfileService, jwtService jwt.Service) *Handler {
	return &Handler{profileService: profileService, jwtService: jwtService}
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/profile", h.getProfile)
	mux.HandleFunc("PUT /api/v1/profile", h.updateProfile)
	mux.HandleFunc("GET /api/v1/profile/books", h.getPurchasedBooks)
	mux.HandleFunc("GET /api/v1/profile/activity", h.getActivityHistory)
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	h.registerRoutes(mux)

	var handler http.Handler = mux
	handler = AuthMiddleware(h.jwtService, handler)
	handler = middlewares.WithLogging(handler)

	return handler
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		_ = WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, profile)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		_ = WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.profileService.UpdateProfile(r.Context(), userID, &req); err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{"message": "profile updated successfully"})
}

func (h *Handler) getPurchasedBooks(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		_ = WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := parseIntWithDefault(r.URL.Query().Get("limit"), 50)
	offset := parseIntWithDefault(r.URL.Query().Get("offset"), 0)

	books, err := h.profileService.GetPurchasedBooks(r.Context(), userID, limit, offset)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, books)
}

func (h *Handler) getActivityHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		_ = WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := parseIntWithDefault(r.URL.Query().Get("limit"), 50)
	offset := parseIntWithDefault(r.URL.Query().Get("offset"), 0)

	activities, err := h.profileService.GetActivityHistory(r.Context(), userID, limit, offset)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if activities == nil {
		activities = []domain.ActivityLogResponse{}
	}

	_ = WriteJSON(w, http.StatusOK, activities)
}

func userIDFromContext(r *http.Request) (int, bool) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0, false
	}
	value, ok := userID.(int)
	return value, ok
}

func parseIntWithDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
