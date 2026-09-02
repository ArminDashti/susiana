package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort         string
	DatabaseURL      string
	MarketAPIBaseURL string
	MarketAPIKey     string
	MarketProvider   string // stub | http
	SyncInterval     time.Duration
	GinMode          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:         getEnv("HTTP_PORT", "8080"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		MarketAPIBaseURL: os.Getenv("MARKET_API_BASE_URL"),
		MarketAPIKey:     os.Getenv("MARKET_API_KEY"),
		MarketProvider:   getEnv("MARKET_PROVIDER", "stub"),
		GinMode:          getEnv("GIN_MODE", "debug"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if interval := os.Getenv("SYNC_INTERVAL"); interval != "" {
		d, err := time.ParseDuration(interval)
		if err != nil {
			return nil, fmt.Errorf("invalid SYNC_INTERVAL: %w", err)
		}
		cfg.SyncInterval = d
	}

	switch cfg.MarketProvider {
	case "stub", "http":
	default:
		return nil, fmt.Errorf("MARKET_PROVIDER must be stub or http, got %q", cfg.MarketProvider)
	}

	if cfg.MarketProvider == "http" && cfg.MarketAPIBaseURL == "" {
		return nil, fmt.Errorf("MARKET_API_BASE_URL is required when MARKET_PROVIDER=http")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
