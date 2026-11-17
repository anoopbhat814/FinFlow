# Changelog

## 0.4.0

- Added `FraudDecision` and `levelFromScore` helpers.
- Added `DomainEvent` union and `eventUserFacing` formatter.
- `Card` now exposes `spentThisMonthMinor`.

## 0.3.0

- Added `Invoice`, `InvoiceLine`, `invoiceTotalMinor` and `isOverdue`.
- `Page<T>` now returns `nextCursor` instead of numeric offsets.

## 0.2.0

- Introduced `LedgerTransaction` and `isBalanced`.
- `ApiError` carries the server `requestId` for support tickets.

## 0.1.0

- Initial release with `Account` and `Transfer` types.
