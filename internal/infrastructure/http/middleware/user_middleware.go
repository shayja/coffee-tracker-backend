// file: internal/infrastructure/http/middleware/user_middleware.go
package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)



func UserMiddleware(repo repositories.UserRepository, ttl time.Duration) func(http.Handler) http.Handler {
	cache := make(map[uuid.UUID]*entities.User)
	var mu sync.RWMutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := utils.GetUserIDOrAbort(w, r)
			if !ok { return }

			var user *entities.User

			// 1. Check cache
			mu.RLock()
			cachedUser, exists := cache[userID]
			mu.RUnlock()

			if exists {
				user = cachedUser
			} else {
				u, err := repo.GetByID(r.Context(), userID)
				if err != nil {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				user = u

				mu.Lock()
				cache[userID] = user
				mu.Unlock()

				// Optional: add expiration using TTL if needed
			}

			// 4. Check status via entity method
			if !user.IsActive() {
				http.Error(w, "User account inactive", http.StatusForbidden)
				return
			}

			// 5. Pass user in context
			ctx := context.WithValue(r.Context(), contextkeys.CurrentUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}