// file: internal/infrastructure/database/supabase.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// NewSupabaseDB creates a new database connection using configuration stored in env.
func NewSupabaseDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
