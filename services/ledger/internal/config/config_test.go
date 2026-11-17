package config

import (
	"testing"
	"time"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("LEDGER_JWT_SECRET", "0123456789abcdef")
	if _, err := Load(); err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func TestLoadRejectsShortSecret(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/finflow")
	t.Setenv("LEDGER_JWT_SECRET", "short")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for short JWT secret")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/finflow")
	t.Setenv("LEDGER_JWT_SECRET", "0123456789abcdef")
	t.Setenv("LEDGER_PORT", "")
	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Port != "8080" {
		t.Errorf("port = %q, want 8080", c.Port)
	}
	if c.FraudTimeout != 800*time.Millisecond {
		t.Errorf("fraud timeout = %v", c.FraudTimeout)
	}
}
