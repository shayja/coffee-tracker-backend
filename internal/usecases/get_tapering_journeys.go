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
	
	entries, err := uc.Repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrInternalError
	}

	if entries == nil {
		return []*entities.TaperingJourney{}, nil
	}
	
	return entries, nil
}
