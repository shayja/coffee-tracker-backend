// file: internal/usecases/delete_coffee_entry.go
package usecases

import (
	"coffee-tracker-backend/internal/repositories"
	"context"

	"github.com/google/uuid"
)

type DeleteCoffeeEntryUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewDeleteCoffeeEntryUseCase(coffeeRepo repositories.CoffeeEntryRepository) *DeleteCoffeeEntryUseCase {
	return &DeleteCoffeeEntryUseCase{
		coffeeRepo: coffeeRepo,
	}
}

// Execute deletes a coffee entry for a given user
func (uc *DeleteCoffeeEntryUseCase) Execute(ctx context.Context, entryID, userID uuid.UUID) error {
	err := uc.coffeeRepo.Delete(ctx, userID, entryID)
	if err != nil {
		if err.Error() == "no coffee entry found with id "+entryID.String()+" for this user" {
			return ErrNotFound
		}
		return ErrInternalError
	}

	return nil
}
