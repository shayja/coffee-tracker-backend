package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
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
	UserID     uuid.UUID `json:"user_id"`
	CoffeeType string    `json:"coffee_type"`
	Size       string    `json:"size"`
	Caffeine   int       `json:"caffeine_mg"`
	Notes      string    `json:"notes"`
	Location   string    `json:"location"`
	Price      float64   `json:"price"`
	Rating     int       `json:"rating"`
}

func (uc *CreateCoffeeEntryUseCase) Execute(ctx context.Context, req CreateCoffeeEntryRequest) (*entities.CoffeeEntry, error) {
	if req.CoffeeType == "" {
		return nil, ErrInvalidInput
	}
	
	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidInput
	}

	entry := &entities.CoffeeEntry{
		ID:         uuid.New(),
		UserID:     req.UserID,
		CoffeeType: req.CoffeeType,
		Size:       req.Size,
		Caffeine:   req.Caffeine,
		Notes:      req.Notes,
		Location:   req.Location,
		Price:      req.Price,
		Rating:     req.Rating,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.coffeeRepo.Create(ctx, entry); err != nil {
		return nil, ErrInternalError
	}

	return entry, nil
}
