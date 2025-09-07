// file: internal/usecases/create_tapering_journey.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateTaperingJourneyUseCase struct {
	Repo repositories.TaperingJourneyRepository
}

func NewCreateTaperingJourneyUseCase(repo repositories.TaperingJourneyRepository) *CreateTaperingJourneyUseCase {
	return &CreateTaperingJourneyUseCase{Repo: repo}
}

func (uc *CreateTaperingJourneyUseCase) Execute(ctx context.Context, userID uuid.UUID, goalFrequency, startLimit, targetLimit, reductionStep, stepPeriod int, startedAt time.Time, statusID int) (*entities.TaperingJourney, error) {
	
	j := &entities.TaperingJourney{
		ID:            uuid.New(),
		UserID:        userID,
		GoalFrequency: entities.GoalFrequency(goalFrequency),
		StartLimit:    startLimit,
		TargetLimit:   targetLimit,
		ReductionStep: reductionStep,
		StepPeriod:    stepPeriod,
		CurrentLimit:  startLimit,
		StatusID:      entities.GoalStatus(statusID),
		StartedAt:     startedAt,
	}
	err := uc.Repo.Create(ctx, j)
	if err != nil {
		return nil, err
	}
	return j, nil
}
