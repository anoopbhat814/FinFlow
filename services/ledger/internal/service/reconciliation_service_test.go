package service

import (
	"reflect"
	"testing"
)

func TestReconcile(t *testing.T) {
	tests := []struct {
		name        string
		ext         []ExternalRecord
		in          []LedgerRecord
		wantMatched int
		want        []Mismatch
	}{
		{
			name:        "all match",
			ext:         []ExternalRecord{{"a", 100, "USD"}, {"b", 200, "USD"}},
			in:          []LedgerRecord{{"b", 200, "USD"}, {"a", 100, "USD"}},
			wantMatched: 2,
		},
		{
			name:        "split captures are summed",
			ext:         []ExternalRecord{{"a", 60, "EUR"}, {"a", 40, "EUR"}},
			in:          []LedgerRecord{{"a", 100, "EUR"}},
			wantMatched: 1,
		},
		{
			name: "every mismatch kind, sorted",
			ext:  []ExternalRecord{{"d", 5, "USD"}, {"c", 70, "USD"}, {"b", 10, "EUR"}, {"a", 9, "USD"}},
			in:   []LedgerRecord{{"c", 75, "USD"}, {"b", 10, "USD"}, {"e", 3, "USD"}, {"a", 9, "USD"}},
			want: []Mismatch{
				{Reference: "b", Kind: CurrencyDiffers, External: 10, Internal: 10},
				{Reference: "c", Kind: AmountDiffers, External: 70, Internal: 75},
				{Reference: "d", Kind: MissingInLedger, External: 5},
				{Reference: "e", Kind: MissingExternally, Internal: 3},
			},
			wantMatched: 1,
		},
		{name: "both empty"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rep := Reconcile(tc.ext, tc.in)
			if rep.Matched != tc.wantMatched {
				t.Errorf("Matched = %d, want %d", rep.Matched, tc.wantMatched)
			}
			if !reflect.DeepEqual(rep.Mismatches, tc.want) {
				t.Errorf("Mismatches = %+v, want %+v", rep.Mismatches, tc.want)
			}
		})
	}
}
