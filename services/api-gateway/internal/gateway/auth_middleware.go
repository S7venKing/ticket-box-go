package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
)

type userIDContextKey struct{}
type userEmailContextKey struct{}

func AuthMiddleware(tokenService *auth.TokenService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if header == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		claims, err := tokenService.Validate(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey{}, claims.UserID)
		ctx = context.WithValue(ctx, userEmailContextKey{}, claims.Email)
		next(w, r.WithContext(ctx))
	}
}
