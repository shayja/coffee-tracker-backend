package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetCoffeeEntriesUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewGetCoffeeEntriesUseCase(coffeeRepo repositories.CoffeeEntryRepository) *GetCoffeeEntriesUseCase {
	return &GetCoffeeEntriesUseCase{
		coffeeRepo: coffeeRepo,
	}
}

func (uc *GetCoffeeEntriesUseCase) Execute(ctx context.Context, userID uuid.UUID, date *string, limit, offset int) ([]*entities.CoffeeEntry, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}

	startDate, err := time.Parse("2006-01-02T00:00:00.000", *date)
	if err != nil {
		return nil, ErrInvalidInput
	}
	endDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, startDate.Location())

	entries, err := uc.coffeeRepo.GetByUserIDAndDateRange(ctx, userID, limit, offset, startDate, endDate)
	if err != nil {
		return nil, ErrInternalError
	}

	return entries, nil
}
