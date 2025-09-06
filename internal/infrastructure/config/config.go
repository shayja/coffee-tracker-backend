// file: internal/infrastructure/config/config.go
// Package config handles application configuration loading and management.
package config

import (
	"os"
)

type Config struct {
	Env		 	string
	Port        string
	DatabaseURL string
	JWTSecret   string
	MagicOtp   	string
	StorageURL  string
	ServiceRoleKey string
}

func Load() *Config {
	return &Config{
		Env:         getEnv("ENV", "dev"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		MagicOtp:    getEnv("MAGIC_OTP", ""),
		StorageURL:  getEnv("SUPABASE_STORAGE_URL", ""),
		ServiceRoleKey:getEnv("SUPABASE_SERVICE_KEY_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
