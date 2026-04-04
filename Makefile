include .env
export

APP_NAME=movie-checkin-service
MAIN_PATH=./cmd/api/main.go

MIGRATION_PATH=./migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_EXTERNAL_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help run build test lint fmt tidy docker-up docker-down docker-logs clean

run:
	go run $(MAIN_PATH)

tidy:
	go mod tidy

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

migrate-up-one:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up 1

migrate-down-one:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" version

seed:
	docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) < migrations/seeds/seed.sql