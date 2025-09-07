// file: internal/usecases/get_tapering_journeys.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"

	"github.com/google/uuid"
)

type GetTaperingJourneysUseCase struct {
	Repo repositories.TaperingJourneyRepository
}

func NewGetTaperingJourneysUseCase(repo repositories.TaperingJourneyRepository) *GetTaperingJourneysUseCase {
	return &GetTaperingJourneysUseCase{Repo: repo}
}

func (uc *GetTaperingJourneysUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]*entities.TaperingJourney, error) {
	return uc.Repo.GetByUserID(ctx, userID)
}
