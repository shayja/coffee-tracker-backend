// file: internal/infrastructure/http/dto/coffee_dto.go
package dto

import (
	"time"

	"github.com/google/uuid"
)


type CreateCoffeeEntryRequest struct {
    Notes     *string   `json:"notes,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  *float64  `json:"latitude,omitempty"`
    Longitude *float64  `json:"longitude,omitempty"`
	CoffeeType *int    	`json:"type,omitempty"`
	Size *int    		`json:"size,omitempty"`
}

type UpdateCoffeeEntryRequest struct {
	ID        uuid.UUID `json:"id"`
    Notes     *string   `json:"notes,omitempty"`
    Timestamp time.Time `json:"timestamp"`
	CoffeeType *int    	`json:"type,omitempty"`
	Size *int    		`json:"size,omitempty"`
}