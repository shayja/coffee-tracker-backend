// file: internal/infrastructure/repositories/tapering_journey_repository.go
package repositories

import (
	"coffee-tracker-backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type TaperingJourneyRepository interface {
	Create(ctx context.Context, journey *entities.TaperingJourney) error
	Update(ctx context.Context, journey *entities.TaperingJourney) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.TaperingJourney, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TaperingJourney, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
