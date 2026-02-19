package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/models"
)

type contextKey string

const ownerKey contextKey = "owner"

// OwnerFromContext retrieves the authenticated Owner from the request context.
func OwnerFromContext(ctx context.Context) *models.Owner {
	if v, ok := ctx.Value(ownerKey).(*models.Owner); ok {
		return v
	}
	return nil
}

// RequireAPIKey is middleware that enforces API key authentication.
// It expects an Authorization header of the form "Bearer ae_<hex>".
func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, `{"error":"missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		key := strings.TrimPrefix(auth, "Bearer ")
		if key == auth { // no "Bearer " prefix found
			http.Error(w, `{"error":"invalid Authorization format, expected: Bearer <api-key>"}`, http.StatusUnauthorized)
			return
		}

		hash := models.HashAPIKey(key)
		var owner models.Owner
		if err := database.DB.Where("api_key_hash = ? AND is_active = ?", hash, true).First(&owner).Error; err != nil {
			http.Error(w, `{"error":"invalid or inactive API key"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ownerKey, &owner)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
