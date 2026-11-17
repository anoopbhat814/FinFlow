package money

// Allocate splits a into n parts without losing a minor unit. The first
// (Minor mod n) parts receive one extra unit, mirroring Money.allocate in
// @finflow/money.
func (a Amount) Allocate(n int) []Amount {
	if n <= 0 {
		return nil
	}
	base := a.Minor / int64(n)
	rem := a.Minor - base*int64(n)
	parts := make([]Amount, n)
	for i := range parts {
		extra := int64(0)
		switch {
		case rem > 0:
			extra, rem = 1, rem-1
		case rem < 0:
			extra, rem = -1, rem+1
		}
		parts[i] = Amount{Minor: base + extra, Currency: a.Currency}
	}
	return parts
}
