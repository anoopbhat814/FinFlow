package domain

import (
	"testing"
	"time"
)

func TestCardAuthorize(t *testing.T) {
	now := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	base := Card{Status: CardActive, ExpMonth: 12, ExpYear: 2028, MonthlyLimitMinor: 100_000, SpentMinor: 40_000}
	with := func(f func(*Card)) Card { c := base; f(&c); return c }
	tests := []struct {
		name   string
		card   Card
		amount int64
		want   error
	}{
		{"within limit", base, 10_000, nil},
		{"exactly at limit", base, 60_000, nil},
		{"over limit", base, 60_001, ErrLimitExceeded},
		{"zero amount", base, 0, ErrInvalidAmount},
		{"frozen", with(func(c *Card) { c.Status = CardFrozen }), 100, ErrCardNotActive},
		{"inactive", with(func(c *Card) { c.Status = CardInactive }), 100, ErrCardNotActive},
		{"expired last month", with(func(c *Card) { c.ExpMonth, c.ExpYear = 5, 2025 }), 100, ErrCardExpired},
		{"valid through expiry month", with(func(c *Card) { c.ExpMonth, c.ExpYear = 6, 2025 }), 100, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.card.Authorize(tc.amount, now); got != tc.want {
				t.Fatalf("Authorize() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardRemaining(t *testing.T) {
	tests := []struct{ limit, spent, want int64 }{{100, 30, 70}, {100, 100, 0}, {100, 150, 0}}
	for _, tc := range tests {
		c := Card{MonthlyLimitMinor: tc.limit, SpentMinor: tc.spent}
		if got := c.Remaining(); got != tc.want {
			t.Errorf("Remaining(%d,%d) = %d, want %d", tc.limit, tc.spent, got, tc.want)
		}
	}
}
