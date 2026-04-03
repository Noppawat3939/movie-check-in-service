APP_NAME=movie-checkin-service
MAIN_PATH=./cmd/api/main.go

.PHONY: help run build test lint fmt tidy docker-up docker-down docker-logs clean

run:
	go run $(MAIN_PATH)

tidy:
	go mod tidy

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down