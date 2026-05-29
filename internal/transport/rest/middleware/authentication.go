package middleware

import (
	"context"
	"docintel/internal/domain"
	"net/http"
)

type MiddleWare struct {
	AuthService domain.AuthService
}

func NewMiddleWare(authService domain.AuthService) *MiddleWare {
	return &MiddleWare{
		AuthService: authService,
	}
}

func (m *MiddleWare) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bearerToken := r.Header.Get("Authorization")

		sessionID := "session:" + bearerToken[len("Bearer "):]

		userID, err := m.AuthService.ValidateUserSession(ctx, sessionID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
