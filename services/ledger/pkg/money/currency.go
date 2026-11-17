package money

// exponents maps ISO 4217 codes to the number of minor-unit decimal places.
var exponents = map[string]int{
	"USD": 2,
	"EUR": 2,
	"GBP": 2,
	"CAD": 2,
	"JPY": 0,
}

// Supported reports whether the currency is handled by the ledger.
func Supported(code string) bool {
	_, ok := exponents[code]
	return ok
}

// Exponent returns the number of decimal places for code, or 2 when unknown.
func Exponent(code string) int {
	if e, ok := exponents[code]; ok {
		return e
	}
	return 2
}
