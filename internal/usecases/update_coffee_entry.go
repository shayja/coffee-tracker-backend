// file: internal/usecases/update_coffee_entry.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/infrastructure/http/models"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type UpdateCoffeeEntryUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewUpdateCoffeeEntryUseCase(coffeeRepo repositories.CoffeeEntryRepository) *UpdateCoffeeEntryUseCase {
	return &UpdateCoffeeEntryUseCase{
		coffeeRepo: coffeeRepo,
	}
}

func (uc *UpdateCoffeeEntryUseCase) Execute(ctx context.Context, userID uuid.UUID, entryID uuid.UUID, req *models.UpdateCoffeeEntryRequest) (*entities.CoffeeEntry, error) {
	// if req.CoffeeType == nil {
	// 	return nil, ErrInvalidInput
	// }
	
	// if req.Rating < 1 || req.Rating > 5 {
	// 	return nil, ErrInvalidInput
	// }

	entry := &entities.CoffeeEntry{
		ID:         	entryID,
		UserID:     	userID,
		CoffeeTypeID: 	req.CoffeeType,
		SizeID:       	req.Size,
		// Caffeine:   req.Caffeine,
		Notes:      	req.Notes,
		// Price:      req.Price,
		// Rating:     req.Rating,
		Timestamp: 		req.Timestamp,
		UpdatedAt:  	utils.NowUTC(),
	}

	if err := uc.coffeeRepo.Update(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
