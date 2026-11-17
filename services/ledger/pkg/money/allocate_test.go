package money

import "testing"

func TestAllocatePreservesTotal(t *testing.T) {
	for _, minor := range []int64{1999, 1, 100, -1001} {
		parts := New(minor, "USD").Allocate(3)
		var sum int64
		for _, p := range parts {
			sum += p.Minor
		}
		if sum != minor {
			t.Errorf("allocate(%d) sums to %d", minor, sum)
		}
	}
}

func TestAllocateDistributesRemainderFirst(t *testing.T) {
	parts := New(1999, "USD").Allocate(3)
	want := []int64{667, 666, 666}
	for i, p := range parts {
		if p.Minor != want[i] {
			t.Errorf("part %d = %d, want %d", i, p.Minor, want[i])
		}
	}
}
