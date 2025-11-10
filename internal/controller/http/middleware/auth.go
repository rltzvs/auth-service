package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"auth/internal/security"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(jwtManager security.JWTManager, logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			userID, err := jwtManager.ValidateToken(token)
			if err != nil {
				logger.Debug("failed to validate token", zap.Error(err))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	return strings.TrimPrefix(authHeader, "Bearer ")
}
