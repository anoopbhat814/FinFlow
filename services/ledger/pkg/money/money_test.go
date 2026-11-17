package money

import (
	"errors"
	"math"
	"testing"
)

func TestAddSameCurrency(t *testing.T) {
	got, err := New(1050, "USD").Add(New(250, "USD"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Minor != 1300 {
		t.Fatalf("got %d, want 1300", got.Minor)
	}
}

func TestAddCurrencyMismatch(t *testing.T) {
	if _, err := New(1, "USD").Add(New(1, "EUR")); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("err = %v", err)
	}
}

func TestAddOverflow(t *testing.T) {
	if _, err := New(math.MaxInt64, "USD").Add(New(1, "USD")); !errors.Is(err, ErrOverflow) {
		t.Fatalf("err = %v", err)
	}
}

func TestString(t *testing.T) {
	cases := map[Amount]string{
		New(123456, "USD"): "1234.56 USD",
		New(-5, "EUR"):     "-0.05 EUR",
		New(0, "GBP"):      "0.00 GBP",
	}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	}
}
