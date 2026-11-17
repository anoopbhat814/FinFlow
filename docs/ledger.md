# Ledger invariants

1. Every ledger transaction has at least two entries.
2. For each currency, the sum of debits equals the sum of credits.
3. Entry amounts are strictly positive; direction carries the sign.
4. `balance_after` on an entry equals the account balance immediately after it was posted.
5. Money is always an integer number of minor units. Floating point is never used.

The invariants are enforced in `domain.Transaction.Validate`, mirrored by a database CHECK on `ledger_entries.amountMinor`, and re-verified nightly by a reconciliation job.
