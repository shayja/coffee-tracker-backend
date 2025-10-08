// file: internal/usecases/get_coffee_entries.go
package usecases

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

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

func (uc *GetCoffeeEntriesUseCase) Execute(ctx context.Context, userID uuid.UUID, dateStr *string, tzOffsetMinutes *int, limit, offset int) ([]*entities.CoffeeEntry, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}


	//log.Printf("GetCoffeeEntriesUseCase: dateStr=%v, tzOffsetMinutes=%v", *dateStr, *tzOffsetMinutes)

	var parsedTime time.Time
	var err error

	// Parse the base date without timezone
	baseTime, err := time.Parse("2006-01-02", *dateStr) // For "2025-08-21"
	if err != nil {
		return nil, ErrInvalidInput
	}

	// Handle timezone offset if provided
	if tzOffsetMinutes != nil {
		parsedTime = adjustTimeWithOffsetMinutes(baseTime, *tzOffsetMinutes)
	} else {
		// Default to UTC if no offset provided
		parsedTime = baseTime.UTC()
	}

	// Convert to UTC and get start/end of day
	utcStart := parsedTime.UTC()
	utcEnd := utcStart.Add(24 * time.Hour) // This is the start of the NEXT day
	

	entries, err := uc.coffeeRepo.GetByUserIDAndDateRange(ctx, userID, limit, offset, utcStart, utcEnd)
	if err != nil {
		return nil, ErrInternalError
	}
	// If entries is nil, return empty slice instead
	if entries == nil {
		return []*entities.CoffeeEntry{}, nil
	}

	return entries, nil
}

// Helper function to adjust time based on timezone offset in minutes
func adjustTimeWithOffsetMinutes(baseTime time.Time, offsetMinutes int) time.Time {
	// Create fixed location based on offset in minutes
	location := time.FixedZone("UserOffset", offsetMinutes*60)

	// Convert base time to the user's timezone
	return time.Date(
		baseTime.Year(),
		baseTime.Month(),
		baseTime.Day(),
		0, 0, 0, 0, // Start of day
		location,
	)
}