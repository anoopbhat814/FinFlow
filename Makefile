.PHONY: bootstrap dev test lint build migrate seed

bootstrap:
	npm install
	cd services/fraud && pip install -r requirements.txt
	cd services/ledger && go mod download

dev:
	docker compose -f deploy/docker-compose.yml up --build

test:
	npm test --workspaces --if-present
	cd services/ledger && go test ./...
	cd services/fraud && pytest -q

lint:
	npx eslint . --ext .js,.ts,.tsx
	cd services/ledger && go vet ./...
	cd services/fraud && ruff check app tests

build:
	npm run build --workspaces --if-present
	cd services/ledger && go build -o bin/ledger ./cmd/ledger

migrate:
	npx prisma migrate deploy --schema prisma/schema.prisma

seed:
	npx tsx prisma/seed.ts
