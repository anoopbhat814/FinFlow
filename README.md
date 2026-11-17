<div align="center">

# FinFlow

**Payments and money movement on a double-entry ledger, with real-time fraud scoring.**

[![CI](https://img.shields.io/github/actions/workflow/status/your-org/FinFlow/ci.yml?branch=main&label=CI)](https://github.com/your-org/FinFlow/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.4-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Python](https://img.shields.io/badge/Python-3.12-3776AB?logo=python&logoColor=white)](https://www.python.org/)
[![Prisma](https://img.shields.io/badge/Prisma-5.10-2D3748?logo=prisma&logoColor=white)](https://www.prisma.io/)
[![Node](https://img.shields.io/badge/Node-20-339933?logo=node.js&logoColor=white)](https://nodejs.org/)

</div>

---

## Overview

Moving money correctly is harder than moving data. Transfers must be exactly-once under client retries, balances must never drift from the books, and risky payments must be caught before funds leave an account, all without making legitimate customers wait.

Go is the primary backend: the ledger service owns the money-movement API, while a Next.js frontend, a Python fraud service and a Node worker support it.

FinFlow is a payments platform that addresses these problems directly. It provides multi-currency accounts, transfers, invoicing and cards backed by a **double-entry ledger** whose invariants are enforced in code and in the database. Every transfer is scored by a dedicated **fraud service** before it posts, and customer notifications and webhooks are delivered asynchronously by a **worker**.

## Key features

- **Multi-currency accounts**: USD, EUR, GBP, CAD and JPY, with checking, savings and business account kinds.
- **Double-entry ledger**: balanced debit/credit entries, strictly positive amounts, and a running `balanceAfterMinor` on every entry. See [docs/ledger.md](docs/ledger.md).
- **Exact money arithmetic**: all amounts are integer minor units; floating point is never used (`packages/money`, `services/ledger/pkg/money`).
- **Idempotent transfers**: every transfer requires an `Idempotency-Key`; client retries are replayed from a Redis-backed store.
- **Real-time fraud scoring**: rule-based signals (large amount, velocity, new large destination, unusual hour, amount outlier) combined with a statistical scorer, returning `ALLOW`, `REVIEW` or `BLOCK`.
- **Graceful degradation**: if the fraud service is unavailable or exceeds its time budget, transfers are routed to manual review rather than failed.
- **Invoicing and cards**: invoices with line items and lifecycle states; virtual and physical cards with monthly limits.
- **Notifications and webhooks**: email delivery via SMTP, HMAC-SHA256 signed webhooks, retry with backoff and a dead-letter stream.
- **Per-user rate limiting** and JWT authentication on the ledger API.

## Architecture

```mermaid
flowchart LR
    Browser["Browser"] --> Web["apps/web<br/>Next.js / TypeScript"]
    Web -- "REST / JSON" --> Ledger["services/ledger<br/>Go"]
    Ledger -- "SQL" --> PG[("PostgreSQL")]
    Ledger -- "idempotency keys" --> Redis[("Redis")]
    Ledger -- "POST /v1/score" --> Fraud["services/fraud<br/>Python / FastAPI"]
    Ledger -- "notification jobs" --> Redis
    Redis --> Worker["apps/worker<br/>Node.js"]
    Worker --> Out["Email (SMTP) / Webhooks"]
```

Transfer request flow:

1. The web app generates an idempotency key and calls `POST /v1/transfers` on the ledger service.
2. The ledger service requests a verdict from the fraud service within an 800 ms budget. On timeout it degrades to `REVIEW`.
3. For `ALLOW`, a serializable Postgres transaction locks both accounts in sorted ID order, writes balanced ledger entries and updates balances.
4. A notification job is published to Redis; the worker renders the template and delivers email or webhooks.

The ledger service is the only writer of account balances and ledger tables. More detail in [docs/architecture.md](docs/architecture.md).

## Tech stack

| Layer | Technology |
|-------|------------|
| Web | Next.js 14, React 18, TypeScript, Zustand |
| Ledger service | Go 1.22, `net/http`, pgx, golang-jwt |
| Fraud service | Python 3.12, FastAPI, Pydantic, NumPy |
| Worker | Node.js 20, ioredis, Nodemailer |
| Data | PostgreSQL 16, Prisma 5 (schema and migrations), Redis 7 |
| Shared packages | `@finflow/types`, `@finflow/money` (npm workspaces) |
| Tooling | ESLint, Prettier, Ruff, pytest, `go test`, Docker Compose, GitHub Actions, Kubernetes |

## Repository structure

```
FinFlow/
├── apps/
│   ├── web/                 # Next.js dashboard (accounts, transfers, invoices, cards)
│   └── worker/              # Notification and webhook worker (Node.js)
├── packages/
│   ├── money/               # Integer minor-unit money arithmetic
│   └── types/               # Shared domain types
├── services/
│   ├── ledger/              # Go: accounts, transfers, double-entry ledger API
│   │   ├── cmd/ledger/      #   entrypoint
│   │   ├── internal/        #   handlers, middleware, service, repository, domain, config
│   │   └── pkg/             #   idempotency, money
│   └── fraud/               # Python: rule + statistical fraud scoring
│       ├── app/             #   FastAPI app and models
│       └── tests/
├── prisma/                  # schema.prisma, migrations, seed.ts
├── deploy/                  # Dockerfiles, docker-compose.yml, k8s/ manifests
├── docs/                    # architecture.md, api.md, ledger.md
├── scripts/                 # check-env, wait-for-db, smoke, OpenAPI generation, k8s rollout
├── .github/                 # CI, release, CodeQL and staging deploy workflows
└── Makefile
```

## Getting started

### Prerequisites

- Node.js 20 and npm
- Go 1.22
- Python 3.12 (3.11+ supported)
- Docker with Docker Compose

### Setup

```bash
git clone https://github.com/your-org/FinFlow.git
cd FinFlow

# Configure environment
cp .env.example .env

# Install JS, Python and Go dependencies
make bootstrap
```

`.env.example` defines the following variables:

| Variable | Default | Purpose |
|----------|---------|---------|
| `DATABASE_URL` | `postgresql://finflow:finflow@localhost:5432/finflow` | PostgreSQL connection |
| `REDIS_URL` | `redis://localhost:6379` | Idempotency store and job queue |
| `LEDGER_PORT` | `8080` | Ledger service port |
| `LEDGER_JWT_SECRET` | `change-me-in-local-dev` | JWT signing secret (set a strong value outside local dev) |
| `FRAUD_SERVICE_URL` | `http://localhost:8000` | Fraud service address |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | API base URL used by the web app |
| `SMTP_HOST` / `SMTP_PORT` / `SMTP_FROM` | `localhost` / `1025` / `no-reply@finflow.example` | Outbound email |

Validate your configuration with `npm run check-env`.

### Run the full stack

```bash
make dev        # docker compose -f deploy/docker-compose.yml up --build
```

This starts PostgreSQL, Redis, the ledger and fraud services, the worker, the web app and MailHog.

| Service | URL |
|---------|-----|
| Web | http://localhost:3000 |
| Ledger API | http://localhost:8080 |
| Fraud service | http://localhost:8000 |
| MailHog (outbound mail) | http://localhost:8025 |

### Database migrations and seed data

With PostgreSQL running:

```bash
make migrate    # prisma migrate deploy
make seed       # tsx prisma/seed.ts
```

### Running services individually

```bash
# Ledger (Go)
cd services/ledger && go run ./cmd/ledger

# Fraud (Python)
cd services/fraud && uvicorn app.main:app --port 8000

# Worker (Node.js)
npm start -w @finflow/worker

# Web (Next.js)
npm run dev -w @finflow/web
```

### Make targets

| Target | Description |
|--------|-------------|
| `make bootstrap` | Install npm, pip and Go dependencies |
| `make dev` | Build and start the full stack with Docker Compose |
| `make test` | Run JS, Go and Python test suites |
| `make lint` | ESLint, `go vet`, Ruff |
| `make build` | Build workspaces and the ledger binary (`services/ledger/bin/ledger`) |
| `make migrate` | Apply Prisma migrations |
| `make seed` | Seed the database |

## API overview

The ledger service exposes a versioned JSON API. All `/v1` endpoints require `Authorization: Bearer <jwt>`. Full reference: [docs/api.md](docs/api.md).

| Service | Method | Path | Description |
|---------|--------|------|-------------|
| Ledger | `GET` | `/healthz` | Liveness and database health |
| Ledger | `GET` | `/v1/accounts` | List the caller's accounts |
| Ledger | `GET` | `/v1/accounts/{id}/entries?limit=50` | Ledger statement for an account |
| Ledger | `POST` | `/v1/transfers` | Create a transfer; `Idempotency-Key` header required |
| Fraud | `GET` | `/healthz` | Liveness |
| Fraud | `POST` | `/v1/score` | Score a transfer (internal, called by the ledger service) |

Example:

```bash
curl -X POST $API/v1/transfers \
  -H "Authorization: Bearer $TOKEN" -H "Idempotency-Key: $(uuidgen)" \
  -d '{"sourceAccountId":"a1","destinationAccountId":"a2","amountMinor":1250,"currency":"USD"}'
```

Errors are JSON (`{ "code": "unprocessable", "message": "insufficient funds" }`); rate-limited calls return `429` with `Retry-After`.

## Data model

The schema lives in [`prisma/schema.prisma`](prisma/schema.prisma) and is migrated with Prisma Migrate.

| Model | Purpose |
|-------|---------|
| `User` | Customer identity, KYC flag |
| `Account` | Currency-denominated account with balance and hold in minor units |
| `Transfer` | Account-to-account payment with a unique idempotency key and status (`PENDING`, `IN_REVIEW`, `COMPLETED`, `FAILED`, `REVERSED`) |
| `LedgerTransaction` | Groups the balanced entries of one posting |
| `LedgerEntry` | Single debit or credit with the resulting balance |
| `FraudCheck` | Score, action (`ALLOW`/`REVIEW`/`BLOCK`) and signals for a transfer |
| `Invoice`, `InvoiceLine` | Invoices and their line items |
| `Card` | Virtual or physical card with a monthly spend limit |
| `Notification` | Notification record per channel (email, SMS, push, webhook) |
| `ApiKey` | Hashed API keys with prefix, revocation and last-used tracking |
| `Statement` | Monthly account statement with opening/closing balance and debit/credit totals |
| `Beneficiary` | Saved payee (IBAN, BIC) with a verification status |
| `Budget`, `BudgetLine` | Per-user spending budgets split into category limits |
| `Dispute`, `DisputeEvidence` | Transfer disputes with reason, status workflow and uploaded evidence |
| `Refund` | Full or partial refund of a completed transfer |
| `WebhookEndpoint`, `WebhookDelivery` | Customer webhook subscriptions and their signed, retried deliveries |
| `AuditLog` | Append-only record of actions by users and the system |
| `KycDocument` | Identity documents with review status and expiry |
| `ExchangeRate` | Time-stamped currency pair rates |
| `CardAuthorization` | Approved and declined card purchase attempts |

## Testing

```bash
make test       # everything
make lint       # static checks
```

| Component | Command |
|-----------|---------|
| JS/TS workspaces (web, worker, money) | `npm test` |
| Ledger | `cd services/ledger && go test -race ./...` |
| Fraud | `cd services/fraud && pytest -q` |

CI (`.github/workflows/ci.yml`) runs lint, tests and `prisma validate` for the web workspace, `go vet` and `go test -race` for the ledger, and `pytest` for the fraud service on every pull request.

## Deployment

- **Images**: Dockerfiles for each component live in `deploy/` and `services/*/Dockerfile`. Pushing a `v*.*.*` tag triggers the release workflow, which builds the web, worker, ledger and fraud images and publishes them to GitHub Container Registry.
- **Kubernetes**: manifests for the ledger and fraud services are in `deploy/k8s/`.
- **Staging**: after CI succeeds on `main`, the staging workflow applies the manifests and waits for the rollout (`scripts/k8s-rollout.js`).
- **Migrations**: run `make migrate` (`prisma migrate deploy`) against the target database before rolling out new versions.

## Security

- JWT bearer authentication and per-user token-bucket rate limiting on the ledger API.
- Idempotency keys prevent duplicate postings on client retries.
- Ledger invariants are enforced in the domain layer and by database constraints.
- Accounts are locked in sorted order within serializable transactions to avoid deadlocks.
- Webhook payloads are signed with HMAC-SHA256.
- `.env` is git-ignored; always set a strong `LEDGER_JWT_SECRET` outside local development.
- CodeQL analysis and Dependabot run on the repository.

To report a vulnerability, follow [SECURITY.md](SECURITY.md). Please do not open public issues for security reports.

## Roadmap

- Invoice creation and card issuing on the ledger API (read, pay, freeze and authorization endpoints exist)
- SMS and push notification channels (modeled in the schema)
- Transfer reversal workflow
- Published OpenAPI specification (generator in `scripts/generate-openapi.js`)

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) and our [Code of Conduct](CODE_OF_CONDUCT.md) before opening a pull request. Release notes are tracked in [CHANGELOG.md](CHANGELOG.md).

## License

Released under the [MIT License](LICENSE).
