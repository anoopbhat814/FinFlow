package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	FraudURL     string
	FraudTimeout time.Duration
	RateLimitRPS int
}

func Load() (*Config, error) {
	c := &Config{
		Port:         getenv("LEDGER_PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("LEDGER_JWT_SECRET"),
		FraudURL:     getenv("FRAUD_SERVICE_URL", "http://localhost:8000"),
		FraudTimeout: 800 * time.Millisecond,
		RateLimitRPS: 20,
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}
	if len(c.JWTSecret) < 16 {
		return nil, fmt.Errorf("config: LEDGER_JWT_SECRET must be at least 16 characters")
	}
	return c, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
