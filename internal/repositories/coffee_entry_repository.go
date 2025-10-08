// file: internal/repositories/coffee_entry_repository.go
package repositories

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/entities"

	"github.com/google/uuid"
)

type CoffeeEntryRepository interface {
	Create(ctx context.Context, entry *entities.CoffeeEntry) error
	Update(ctx context.Context, entry *entities.CoffeeEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.CoffeeEntry, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.CoffeeEntry, error)
	GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, limit int, offset int, startDate, endDate time.Time) ([]*entities.CoffeeEntry, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
	GetStats(ctx context.Context, userID uuid.UUID) (*entities.CoffeeStats, error)
	GetCount(ctx context.Context, userID uuid.UUID) (int, error)
}
