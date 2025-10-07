package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Env                 string
	Port                string
	DatabaseURL         string
	JWTSecret           string
	MagicOtp            string
	StorageURL          string
	ServiceRoleKey      string
	ProfileImageBucket  string
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
}

func Load() *Config {
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	if v := os.Getenv("ACCESS_TOKEN_TTL"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			accessTTL = dur
		}
	}
	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			refreshTTL = dur
		}
	}

	return &Config{
		Env:                getEnv("ENV", "dev"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		MagicOtp:           getEnv("MAGIC_OTP", ""),
		StorageURL:         getEnv("SUPABASE_STORAGE_URL", ""),
		ServiceRoleKey:     getEnv("SUPABASE_SERVICE_KEY_ID", ""),
		ProfileImageBucket: getEnv("PROFILE_IMAGE_BUCKET", ""),
		AccessTokenTTL:     accessTTL,
		RefreshTokenTTL:    refreshTTL,
	}
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
