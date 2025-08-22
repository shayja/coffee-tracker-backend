package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

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
    Notes     string    `json:"notes"`
    Timestamp time.Time `json:"timestamp"`
}

func (uc *EditCoffeeEntryUseCase) Execute(ctx context.Context, req EditCoffeeEntryRequest, userID uuid.UUID) (*entities.CoffeeEntry, error) {
	// if req.CoffeeType == "" {
	// 	return nil, ErrInvalidInput
	// }
	
	// if req.Rating < 1 || req.Rating > 5 {
	// 	return nil, ErrInvalidInput
	// }

	entry := &entities.CoffeeEntry{
		ID:         req.ID,
		UserID:     userID,
		// CoffeeType: req.CoffeeType,
		// Size:       req.Size,
		// Caffeine:   req.Caffeine,
		Notes:      req.Notes,
		// Location:   req.Location,
		// Price:      req.Price,
		// Rating:     req.Rating,
		Timestamp: req.Timestamp,
		UpdatedAt:  time.Now(),
	}

	if err := uc.coffeeRepo.Update(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
