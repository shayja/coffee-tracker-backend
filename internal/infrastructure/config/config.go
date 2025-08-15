package config

import (
	"os"
)

type Config struct {
	Env		 	string
	Port        string
	DatabaseURL string
	SupabaseURL string
	SupabaseKey string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Env:         getEnv("ENV", "dev"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		SupabaseURL: getEnv("SUPABASE_URL", ""),
		SupabaseKey: getEnv("SUPABASE_KEY", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
