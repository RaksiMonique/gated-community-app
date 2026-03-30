package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Env         string
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:pass@localhost:5432/gated_db?sslmode=disable"
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: dbURL,
		Env:         getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
