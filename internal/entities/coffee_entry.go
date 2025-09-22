// file: internal/entities/coffee_entry.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type CoffeeEntry struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CoffeeTypeID *int     `json:"type" db:"coffee_type_id"`     
	SizeID      *int  	  `json:"size" db:"size_id"`
	//Caffeine    *int     `json:"caffeine_mg" db:"caffeine_mg"`
	Notes       *string    `json:"notes" db:"notes"`
	//Price       float64  `json:"price" db:"price"`
	//Rating      int      `json:"rating" db:"rating"` // 1-5 scale
	Latitude  	*float64   `json:"latitude" db:"latitude"`
    Longitude 	*float64   `json:"longitude" db:"longitude"`
	Timestamp   time.Time  `json:"timestamp" db:"created_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CoffeeStats struct {
	TotalEntries     int     `json:"total_entries"`
	//TotalCaffeine    int     `json:"total_caffeine_mg"`
	//AverageRating    float64 `json:"average_rating"`
	//FavoriteCoffee   string  `json:"favorite_coffee"`
	//TotalSpent       float64 `json:"total_spent"`
	EntriesThisWeek  int     `json:"entries_this_week"`
	EntriesThisMonth int     `json:"entries_this_month"`
}
