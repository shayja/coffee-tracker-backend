// file: internal/usecases/create_coffee_entry.go
package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/infrastructure/http/dto"
	"coffee-tracker-backend/internal/repositories"

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

func (uc *CreateCoffeeEntryUseCase) Execute(ctx context.Context, userID uuid.UUID, req *dto.CreateCoffeeEntryRequest) (*entities.CoffeeEntry, error) {

	// if req.Rating < 1 || req.Rating > 5 {
	// 	return nil, ErrInvalidInput
	// }

	entry := &entities.CoffeeEntry{
		ID:         uuid.New(),
		UserID:     userID,
		CoffeeTypeID: req.CoffeeType,
		SizeID:       req.Size,
		// Caffeine:   req.Caffeine,
		Notes:      req.Notes,
		// Price:      req.Price,
		// Rating:     req.Rating,
		Latitude:   req.Latitude,
    	Longitude:  req.Longitude,
		Timestamp: req.Timestamp,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if err := uc.coffeeRepo.Create(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
