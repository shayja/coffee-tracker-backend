// file: internal/infrastructure/database/supabase.go
// Package database provides utilities for connecting to and interacting with the Supabase database.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type SupabaseConfig struct {
	DatabaseURL string
	SupabaseURL string
	SupabaseKey string
}

func NewSupabaseConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings for Supabase
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// NewSupabaseClient creates a new Supabase client configuration
func NewSupabaseClient(cfg SupabaseConfig) *SupabaseConfig {
	return &SupabaseConfig{
		DatabaseURL: cfg.DatabaseURL,
		SupabaseURL: cfg.SupabaseURL,
		SupabaseKey: cfg.SupabaseKey,
	}
}
