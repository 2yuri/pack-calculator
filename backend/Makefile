
#!make
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: run build test test-unit test-integration generate-mocks migrate-up migrate-down docker-up docker-down docker-build

DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
ifneq ($(CACHE_PASSWORD),)
REDIS_ADDR=redis://default:$(CACHE_PASSWORD)@$(CACHE_HOST):$(CACHE_PORT)/0
else
REDIS_ADDR=$(CACHE_HOST):$(CACHE_PORT)
endif

## Project
run:
	go run cmd/main.go

build:
	go build -o bin/sorting cmd/main.go

## Tests
test:
	go test ./... -v -count=1

test-integration:
	TEST_DATABASE_URL=$(DSN) TEST_REDIS_ADDR=$(REDIS_ADDR) go test ./internal/repository/ -v -count=1

## Mocks
generate-mocks:
	mockgen -source=internal/domain/product.go -destination=internal/mocks/product.go -package=mocks
	mockgen -source=internal/domain/pack.go -destination=internal/mocks/pack.go -package=mocks
	mockgen -source=internal/domain/order.go -destination=internal/mocks/order.go -package=mocks

## Migrations
migrate-up:
	migrate -path migrations -database $(DSN) up

migrate-down:
	migrate -path migrations -database $(DSN) down

## Docker
docker-up:
	docker compose up -d postgres redis

docker-down:
	docker compose down -v
