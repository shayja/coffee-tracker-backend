// file: internal/config/config.go
package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type OtpStrength string

const (
	OTP_EASY   OtpStrength = "easy"
	OTP_STRONG OtpStrength = "strong"
)

type Config struct {
	Env                string
	Port               string
	DatabaseURL        string
	JWTSecret          string
	OtpStrength        OtpStrength
	MagicOtp           string
	StorageURL         string
	ServiceRoleKey     string
	ProfileImageBucket string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

func Load() (*Config, error) {
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	if v := os.Getenv("ACCESS_TOKEN_TTL"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			accessTTL = dur
		} else {
			return nil, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %v", err)
		}
	}
	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			refreshTTL = dur
		} else {
			return nil, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %v", err)
		}
	}

	cfg := &Config{
		Env:                getEnv("ENV", "dev"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		OtpStrength:        OtpStrength(getEnv("OTP_STRENGTH", "easy")),
		MagicOtp:           getEnv("MAGIC_OTP", ""),
		StorageURL:         getEnv("SUPABASE_STORAGE_URL", ""),
		ServiceRoleKey:     getEnv("SUPABASE_SERVICE_KEY_ID", ""),
		ProfileImageBucket: getEnv("PROFILE_IMAGE_BUCKET", ""),
		AccessTokenTTL:     accessTTL,
		RefreshTokenTTL:    refreshTTL,
	}

	// Validate immediately
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) validate() error {
	// Environment
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	if _, ok := validEnvs[c.Env]; !ok {
		return fmt.Errorf("invalid ENV: %s (must be one of: dev, staging, prod)", c.Env)
	}

	// Port
	if c.Port == "" {
		return errors.New("PORT is required")
	}

	// Database
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	// JWT
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}

	// OTP strength
	switch c.OtpStrength {
	case OTP_EASY, OTP_STRONG:
		// ok
	default:
		return fmt.Errorf("invalid OTP_STRENGTH: %s (must be 'easy' or 'strong')", c.OtpStrength)
	}

	// Supabase Storage
	if c.StorageURL == "" {
		return errors.New("SUPABASE_STORAGE_URL is required")
	}
	if c.ServiceRoleKey == "" {
		return errors.New("SUPABASE_SERVICE_KEY_ID is required")
	}
	if c.ProfileImageBucket == "" {
		return errors.New("PROFILE_IMAGE_BUCKET is required")
	}

	// TTLs
	if c.AccessTokenTTL <= 0 {
		return errors.New("ACCESS_TOKEN_TTL must be greater than 0")
	}
	if c.RefreshTokenTTL <= 0 {
		return errors.New("REFRESH_TOKEN_TTL must be greater than 0")
	}

	return nil
}
