// file: internal/infrastructure/repositories/coffee_entry_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/utils"

	"github.com/google/uuid"
)

type CoffeeEntryRepositoryImpl struct {
	db *sql.DB
}

func NewCoffeeEntryRepositoryImpl(db *sql.DB) repositories.CoffeeEntryRepository {
	return &CoffeeEntryRepositoryImpl{db: db}
}

func (r *CoffeeEntryRepositoryImpl) Create(ctx context.Context, entry *entities.CoffeeEntry) error {
	query := `
		INSERT INTO coffee_entries (id, user_id, notes, timestamp)
		VALUES ($1, $2, $3, $4)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		entry.ID,
		entry.UserID,
		utils.NullIfEmpty(entry.Notes),
		entry.Timestamp,
	)
	
	return err
}

func (r *CoffeeEntryRepositoryImpl) Update(ctx context.Context, entry *entities.CoffeeEntry) error {
	query := `
		UPDATE coffee_entries 
		SET notes = $2, timestamp = $3, updated_at = $4
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query,
		entry.ID,
		utils.NullIfEmpty(entry.Notes),
		entry.Timestamp,
		time.Now(),
	)
	
	return err
}

func (r *CoffeeEntryRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.CoffeeEntry, error) {
	query := `
		SELECT id, user_id, notes, timestamp, created_at, updated_at
		FROM coffee_entries
		WHERE id = $1
		LIMIT 1
	`
	
	var entry entities.CoffeeEntry
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.Notes,
		&entry.Timestamp,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &entry, nil
}

func (r *CoffeeEntryRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.CoffeeEntry, error) {
	query := `
		SELECT id, user_id, notes, timestamp, created_at, updated_at
		FROM coffee_entries
		WHERE user_id = $1
		ORDER BY timestamp ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*entities.CoffeeEntry
	for rows.Next() {
		var entry entities.CoffeeEntry
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Notes,
			&entry.Timestamp,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	
	return entries, rows.Err()
}

func (r *CoffeeEntryRepositoryImpl) GetByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, limit int, offset int, startDate, endDate time.Time) ([]*entities.CoffeeEntry, error) {
	
	query := `
		SELECT id, user_id, notes, timestamp, created_at, updated_at
		FROM coffee_entries
		WHERE user_id = $1 
		AND timestamp >= $4 AND timestamp < $5
		ORDER BY timestamp ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*entities.CoffeeEntry
	for rows.Next() {
		var entry entities.CoffeeEntry
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Notes,
			&entry.Timestamp,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	
	return entries, rows.Err()
}

func (r *CoffeeEntryRepositoryImpl) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM coffee_entries WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}

func (r *CoffeeEntryRepositoryImpl) GetStats(ctx context.Context, userID uuid.UUID) (*entities.CoffeeStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_entries,
			--COALESCE(SUM(caffeine_mg), 0) as total_caffeine,
			--COALESCE(AVG(rating), 0) as average_rating,
			--COALESCE(SUM(price), 0) as total_spent,
			(SELECT COUNT(*) FROM coffee_entries WHERE user_id = $1 AND timestamp >= NOW() - INTERVAL '7 days') as entries_this_week,
			(SELECT COUNT(*) FROM coffee_entries WHERE user_id = $1 AND timestamp >= NOW() - INTERVAL '30 days') as entries_this_month
		FROM coffee_entries 
		WHERE user_id = $1
	`
	
	var stats entities.CoffeeStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.TotalEntries,
		// &stats.TotalCaffeine,
		// &stats.AverageRating,
		// &stats.TotalSpent,
		&stats.EntriesThisWeek,
		&stats.EntriesThisMonth,
	)
	
	if err != nil {
		return nil, err
	}
/*
	// Get favorite coffee type
	favoriteQuery := `
		SELECT coffee_type
		FROM coffee_entries 
		WHERE user_id = $1
		GROUP BY coffee_type
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`
	
	var favoriteCoffee sql.NullString
	err = r.db.QueryRowContext(ctx, favoriteQuery, userID).Scan(&favoriteCoffee)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	if favoriteCoffee.Valid {
		stats.FavoriteCoffee = favoriteCoffee.String
	}
*/	
	return &stats, nil
}

func (r *CoffeeEntryRepositoryImpl) GetCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM coffee_entries WHERE user_id = $1`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
