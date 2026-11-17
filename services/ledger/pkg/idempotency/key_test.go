package idempotency

import (
	"strings"
	"testing"
)

func TestScopedKeyDiffersPerUser(t *testing.T) {
	a, _ := ScopedKey("u1", "POST /v1/transfers", "abc")
	b, _ := ScopedKey("u2", "POST /v1/transfers", "abc")
	if a == b {
		t.Fatal("keys for different users must differ")
	}
}

func TestScopedKeyRejectsInvalid(t *testing.T) {
	for _, k := range []string{"", "   ", strings.Repeat("x", 129), "bad\nkey"} {
		if _, err := ScopedKey("u1", "r", k); err == nil {
			t.Errorf("expected error for %q", k)
		}
	}
}
