# Architecture

FinFlow is split along operational boundaries rather than technical layers.

## Request flow for a transfer

1. The web app generates an idempotency key and POSTs `/v1/transfers` to the ledger service.
2. The ledger service asks the fraud service for a verdict (`ALLOW`, `REVIEW`, `BLOCK`) with an 800 ms budget. On timeout it degrades to `REVIEW` instead of failing the payment.
3. For `ALLOW`, a serializable Postgres transaction locks both accounts in sorted ID order, writes balanced ledger entries and adjusts balances.
4. A notification job is published to Redis; the Node worker renders a template and sends email or webhooks with quadratic backoff and a dead-letter stream.

## Data ownership

- Postgres is the system of record; the schema lives in `prisma/schema.prisma` and is migrated with Prisma Migrate.
- The ledger service is the only writer of `accounts.balanceMinor` and `ledger_*` tables.
- The fraud service keeps rolling per-account history in memory and is safe to restart.

## Failure modes

| Failure | Behaviour |
|---------|-----------|
| Fraud service down | Transfers go to manual review |
| Duplicate client retry | Replayed from the idempotency store |
| Deadlock risk | Account locks acquired in sorted order |
