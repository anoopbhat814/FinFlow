package service

import "sort"

// ExternalRecord is a line from a bank or processor settlement file.
type ExternalRecord struct {
	Reference   string
	AmountMinor int64
	Currency    string
}

// LedgerRecord is the internal view of the same posting.
type LedgerRecord struct {
	Reference   string
	AmountMinor int64
	Currency    string
}

type MismatchKind string

const (
	MissingInLedger   MismatchKind = "MISSING_IN_LEDGER"
	MissingExternally MismatchKind = "MISSING_EXTERNALLY"
	AmountDiffers     MismatchKind = "AMOUNT_DIFFERS"
	CurrencyDiffers   MismatchKind = "CURRENCY_DIFFERS"
)

type Mismatch struct {
	Reference string       `json:"reference"`
	Kind      MismatchKind `json:"kind"`
	External  int64        `json:"externalMinor"`
	Internal  int64        `json:"internalMinor"`
}

type ReconciliationReport struct {
	Matched    int        `json:"matched"`
	Mismatches []Mismatch `json:"mismatches"`
}

// Reconcile compares settlement lines with ledger postings by reference.
// Duplicate references on one side are summed, which mirrors how partial
// captures appear in processor files. Output is sorted for stable reports.
func Reconcile(external []ExternalRecord, internal []LedgerRecord) ReconciliationReport {
	type agg struct {
		amount   int64
		currency string
	}
	ext := map[string]agg{}
	for _, e := range external {
		a := ext[e.Reference]
		a.amount += e.AmountMinor
		a.currency = e.Currency
		ext[e.Reference] = a
	}
	int_ := map[string]agg{}
	for _, l := range internal {
		a := int_[l.Reference]
		a.amount += l.AmountMinor
		a.currency = l.Currency
		int_[l.Reference] = a
	}

	var rep ReconciliationReport
	for ref, e := range ext {
		l, ok := int_[ref]
		switch {
		case !ok:
			rep.Mismatches = append(rep.Mismatches, Mismatch{Reference: ref, Kind: MissingInLedger, External: e.amount})
		case e.currency != l.currency:
			rep.Mismatches = append(rep.Mismatches, Mismatch{Reference: ref, Kind: CurrencyDiffers, External: e.amount, Internal: l.amount})
		case e.amount != l.amount:
			rep.Mismatches = append(rep.Mismatches, Mismatch{Reference: ref, Kind: AmountDiffers, External: e.amount, Internal: l.amount})
		default:
			rep.Matched++
		}
	}
	for ref, l := range int_ {
		if _, ok := ext[ref]; !ok {
			rep.Mismatches = append(rep.Mismatches, Mismatch{Reference: ref, Kind: MissingExternally, Internal: l.amount})
		}
	}
	sort.Slice(rep.Mismatches, func(i, j int) bool { return rep.Mismatches[i].Reference < rep.Mismatches[j].Reference })
	return rep
}
