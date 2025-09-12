// file: internal/usecases/create_coffee_entry.go
package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type CreateCoffeeEntryUseCase struct {
	coffeeRepo repositories.CoffeeEntryRepository
}

func NewCreateCoffeeEntryUseCase(coffeeRepo repositories.CoffeeEntryRepository) *CreateCoffeeEntryUseCase {
	return &CreateCoffeeEntryUseCase{
		coffeeRepo: coffeeRepo,
	}
}

type CreateCoffeeEntryRequest struct {
    Notes     *string    `json:"notes"`
    Timestamp time.Time `json:"timestamp"`
}

func (uc *CreateCoffeeEntryUseCase) Execute(ctx context.Context, req CreateCoffeeEntryRequest, userID uuid.UUID) (*entities.CoffeeEntry, error) {
	// if req.CoffeeType == "" {
	// 	return nil, ErrInvalidInput
	// }
	
	// if req.Rating < 1 || req.Rating > 5 {
	// 	return nil, ErrInvalidInput
	// }

	entry := &entities.CoffeeEntry{
		ID:         uuid.New(),
		UserID:     userID,
		// CoffeeType: req.CoffeeType,
		// Size:       req.Size,
		// Caffeine:   req.Caffeine,
		Notes:      req.Notes,
		// Location:   req.Location,
		// Price:      req.Price,
		// Rating:     req.Rating,
		Timestamp: req.Timestamp,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.coffeeRepo.Create(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
