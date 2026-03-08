.PHONY: dev run build air infra-dev infra-prod backup-postgres backup-redis backup-waha

dev run:
	go run ./cmd/api

air:
	air

build:
	go build -o bin/nemoris ./cmd/api

# Development: postgres, redis, waha only. Backend runs locally (air/go run).
infra-dev:
	cd .. && docker-compose -f deployments/docker-compose.dev.yml up -d

# Production: full stack including backend container.
infra-prod:
	cd .. && docker-compose -f deployments/docker-compose.prod.yml up -d --build

# Backup scripts (run from project root). Manual only — no cron.
backup-postgres:
	cd .. && ./scripts/backup/backup_postgres.sh

backup-redis:
	cd .. && ./scripts/backup/backup_redis.sh

backup-waha:
	cd .. && ./scripts/backup/backup_waha.sh