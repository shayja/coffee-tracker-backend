// file: internal/usecases/update_tapering_journey.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
)

type UpdateTaperingJourneyUseCase struct {
	Repo repositories.TaperingJourneyRepository
}

func NewUpdateTaperingJourneyUseCase(repo repositories.TaperingJourneyRepository) *UpdateTaperingJourneyUseCase {
	return &UpdateTaperingJourneyUseCase{Repo: repo}
}

func (uc *UpdateTaperingJourneyUseCase) Execute(ctx context.Context, journey *entities.TaperingJourney) error {
	return uc.Repo.Update(ctx, journey)
}
