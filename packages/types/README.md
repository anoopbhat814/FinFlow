# @finflow/types

Shared TypeScript contracts for the FinFlow platform. Everything exported here mirrors a payload produced by the ledger service (`services/ledger`) or the fraud service (`services/fraud`).

## Usage

```ts
import type { Transfer, CreateTransferRequest } from '@finflow/types';
import { isTerminal } from '@finflow/types';

if (isTerminal(transfer.status)) {
  // stop polling
}
```

## Conventions

- Monetary values are always integers in minor units and carry a `Minor` suffix (`amountMinor`).
- Timestamps are ISO-8601 strings in UTC.
- Enumerations are string unions whose members match the Prisma enums in `prisma/schema.prisma`.
- Helpers in this package must be pure and free of runtime dependencies.

## Adding a type

1. Create `src/<name>.ts` and export the interface plus any pure helpers.
2. Re-export it from `src/index.ts`.
3. If it maps to a database table, update `prisma/schema.prisma` and add a migration.
