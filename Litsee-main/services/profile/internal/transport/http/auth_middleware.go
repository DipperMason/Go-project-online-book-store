package http

import (
	"context"
	"jwt"
	"net/http"
	"strings"
)

// AuthMiddleware проверяет JWT токен от auth-сервиса.
func AuthMiddleware(jwtService jwt.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			_ = WriteError(w, http.StatusUnauthorized, "missing authorization token")
			return
		}

		claims, err := jwtService.ValidateAndParseToken(jwt.TokenPair{Token: token})
		if err != nil {
			_ = WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}
