# Contributing to FinFlow

Thanks for helping improve FinFlow. This guide covers how to set up the repository, the conventions we follow, and how changes get merged.

## Getting started

1. Fork the repository and clone your fork.
2. Install Node.js 20, Go 1.22 and Python 3.12.
3. Copy `.env.example` to `.env` and run `npm run check-env`.
4. Start dependencies with `docker compose -f deploy/docker-compose.yml up -d postgres redis`.
5. Apply the schema with `npm run db:migrate` and seed demo data with `npm run db:seed`.

## Repository layout

| Path | What lives there |
|------|------------------|
| `apps/web` | Next.js customer dashboard |
| `apps/worker` | Notification consumer (Node.js) |
| `services/ledger` | Double-entry ledger API (Go) |
| `services/fraud` | Fraud scoring API (Python) |
| `packages/*` | Shared TypeScript libraries |
| `prisma` | Database schema and migrations |

## Branches and commits

- Branch from `main` using `feat/<topic>`, `fix/<topic>` or `docs/<topic>`.
- Write commit messages in the imperative mood and follow [Conventional Commits](https://www.conventionalcommits.org/), for example `fix(ledger): reject zero-amount entries`.
- Keep pull requests focused. Large refactors should be discussed in an issue first.

## Quality checks

Run these before opening a pull request:

```
npm run lint
npm test
(cd services/ledger && go vet ./... && go test -race ./...)
(cd services/fraud && pytest -q)
```

## Database changes

- Never edit an applied migration. Add a new one with `npx prisma migrate dev --name <summary>`.
- Money columns are `BIGINT` minor units. Do not introduce floating-point amounts.
- Update `docs/ledger.md` if an invariant changes.

## Code style

Formatting is enforced by Prettier, `gofmt` and Black-compatible settings in `pyproject.toml`. An `.editorconfig` is provided for editors.

## Reviews

At least one maintainer approval is required, and CI must be green. Security-sensitive changes (authentication, ledger posting, fraud thresholds) need two approvals.

## Reporting security issues

Please do not open public issues for vulnerabilities. See [SECURITY.md](SECURITY.md).
