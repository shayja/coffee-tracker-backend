// file: internal/contextkeys/context_helpers.go
package contextkeys

import (
	"coffee-tracker-backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

// UserFromContext retrieves the full User object from context
func UserFromContext(ctx context.Context) (*entities.User, bool) {
    user, ok := ctx.Value(CurrentUserKey).(*entities.User)
    return user, ok
}

// UserIDFromContext retrieves the user ID from context
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
