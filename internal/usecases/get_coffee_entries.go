package usecases

import (
	"context"

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
		limit = 20 // default limit
	}

	entries, err := uc.coffeeRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, ErrInternalError
	}

	return entries, nil
}
