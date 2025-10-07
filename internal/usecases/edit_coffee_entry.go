// file: internal/usecases/edit_coffee_entry.go
package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type EditCoffeeEntryUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewEditCoffeeEntryUseCase(coffeeRepo repositories.CoffeeEntryRepository) *EditCoffeeEntryUseCase {
	return &EditCoffeeEntryUseCase{
		coffeeRepo: coffeeRepo,
	}
}

type EditCoffeeEntryRequest struct {
	ID        uuid.UUID `json:"id"`
    Notes     *string   `json:"notes,omitempty"`
    Timestamp time.Time `json:"timestamp"`
	CoffeeType *int    	`json:"type,omitempty"`
	Size *int    		`json:"size,omitempty"`
}

func (uc *EditCoffeeEntryUseCase) Execute(ctx context.Context, req EditCoffeeEntryRequest, userID uuid.UUID) (*entities.CoffeeEntry, error) {
	// if req.CoffeeType == nil {
	// 	return nil, ErrInvalidInput
	// }
	
	// if req.Rating < 1 || req.Rating > 5 {
	// 	return nil, ErrInvalidInput
	// }

	entry := &entities.CoffeeEntry{
		ID:         	req.ID,
		UserID:     	userID,
		CoffeeTypeID: 	req.CoffeeType,
		SizeID:       	req.Size,
		// Caffeine:   req.Caffeine,
		Notes:      	req.Notes,
		// Price:      req.Price,
		// Rating:     req.Rating,
		Timestamp: 		req.Timestamp,
		UpdatedAt:  	time.Now().UTC(),
	}

	if err := uc.coffeeRepo.Update(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
