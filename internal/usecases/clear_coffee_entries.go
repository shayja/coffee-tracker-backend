// file: internal/usecases/clear_coffee_entries.go
package usecases

import (
	"coffee-tracker-backend/internal/repositories"
	"context"

	"github.com/google/uuid"
)

type ClearCoffeeEntriesUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewClearCoffeeEntriesUseCase(coffeeRepo repositories.CoffeeEntryRepository) *ClearCoffeeEntriesUseCase {
	return &ClearCoffeeEntriesUseCase{
		coffeeRepo: coffeeRepo,
	}
}

// Execute deletes a coffee entry for a given user
func (uc *ClearCoffeeEntriesUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	err := uc.coffeeRepo.DeleteAll(ctx, userID)
	if err != nil {
		if err.Error() == "no coffee entries for userid "+userID.String() {
			return ErrNotFound
		}
		return ErrInternalError
	}

	return nil
}
