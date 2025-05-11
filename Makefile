include	.env
MIGRATIONS_PATH = ./cmd/migrate/migrations

build:
		@go build -o bin/ws_practice_1 ./cmd/api/*.go

run: build
		@./bin/ws_practice_1

.PHONY:	migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY:	migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY:	migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run ./cmd/seed/main.go
