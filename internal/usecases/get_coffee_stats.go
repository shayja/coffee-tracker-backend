// file: internal/usecases/get_coffee_stats.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type GetCoffeeStatsUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewGetCoffeeStatsUseCase(coffeeRepo repositories.CoffeeEntryRepository) *GetCoffeeStatsUseCase {
	return &GetCoffeeStatsUseCase{
		coffeeRepo: coffeeRepo,
	}
}

func (uc *GetCoffeeStatsUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entities.CoffeeStats, error) {
	stats, err := uc.coffeeRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, ErrInternalError
	}

	return stats, nil
}
