# Changelog

All notable changes to this project are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Account closure timestamp (`closedAt`) and API key labels.
- Redis-backed idempotency store for multi-replica ledger deployments.

### Changed
- Web pages now ship loading and error boundaries.

## [0.4.0] - 2025-06-12

### Added
- Virtual and physical cards with freeze and unfreeze actions.
- Fraud scoring service with rule and model signals; transfers can be held for review.
- Webhook delivery with HMAC-SHA256 signatures.

### Fixed
- Transfer history no longer flickers while polling.

## [0.3.0] - 2025-04-02

### Added
- Invoicing: draft, send, and mark invoices as paid.
- CSV export for transaction history.

### Changed
- Rate limiting is now applied per authenticated user.

## [0.2.0] - 2025-02-15

### Added
- Double-entry ledger with database-level invariants.
- Idempotent transfer creation via the `Idempotency-Key` header.

## [0.1.0] - 2025-01-01

### Added
- Initial release: accounts, transfers and the web dashboard.
