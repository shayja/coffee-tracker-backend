// file: internal/usecases/delete_tapering_journey.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"context"

	"github.com/google/uuid"
)


type DeleteTaperingJourneyUseCase struct {
	Repo repositories.TaperingJourneyRepository
}

func NewDeleteTaperingJourneyUseCase(repo repositories.TaperingJourneyRepository) *DeleteTaperingJourneyUseCase {
	return &DeleteTaperingJourneyUseCase{Repo: repo}
}

func (uc *DeleteTaperingJourneyUseCase) Execute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return uc.Repo.Delete(ctx, id, userID)
}
