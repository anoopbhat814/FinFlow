package domain

import (
	"testing"
	"time"
)

func TestBuildStatement(t *testing.T) {
	start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	at := func(day int) time.Time { return time.Date(2025, 3, day, 9, 0, 0, 0, time.UTC) }
	entries := []Entry{
		{Direction: Credit, AmountMinor: 5_000, BalanceAfter: 15_000, PostedAt: at(2)},
		{Direction: Debit, AmountMinor: 2_000, BalanceAfter: 13_000, PostedAt: at(10)},
		{Direction: Debit, AmountMinor: 500, BalanceAfter: 12_500, PostedAt: at(20)},
		{Direction: Credit, AmountMinor: 999, BalanceAfter: 13_499, PostedAt: end}, // excluded: end is exclusive
		{Direction: Credit, AmountMinor: 1, BalanceAfter: 10_001, PostedAt: start.Add(-time.Second)},
	}
	tests := []struct {
		name                          string
		start, end                    time.Time
		entries                       []Entry
		wantErr                       error
		opening, closing, deb, cred   int64
		count                         int
	}{
		{"mixed month", start, end, entries, nil, 10_000, 12_500, 2_500, 5_000, 3},
		{"empty period", start, end, nil, nil, 0, 0, 0, 0, 0},
		{"inverted period", end, start, nil, ErrInvalidPeriod, 0, 0, 0, 0, 0},
		{"zero length", start, start, nil, ErrInvalidPeriod, 0, 0, 0, 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st, err := BuildStatement("acc", tc.start, tc.end, tc.entries, end)
			if err != tc.wantErr {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			if st.OpeningMinor != tc.opening || st.ClosingMinor != tc.closing ||
				st.TotalDebits != tc.deb || st.TotalCredits != tc.cred || st.EntryCount != tc.count {
				t.Fatalf("got %+v", st)
			}
			if st.NetChange() != tc.closing-tc.opening {
				t.Fatalf("NetChange = %d", st.NetChange())
			}
		})
	}
}
