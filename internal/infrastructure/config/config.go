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
}

func Load() *Config {
	return &Config{
		Env:         getEnv("ENV", "dev"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		MagicOtp:    getEnv("MAGIC_OTP", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
