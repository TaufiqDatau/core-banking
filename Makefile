# Load environment variables from .env file
include .env

postgres:
	docker run --name postgres17 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -p 5432:5432 -d postgres:17.2-alpine3.21

createdb:
	docker exec -it postgres17 createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) core_bank

dropdb:
	docker exec -it postgres17 dropdb core_bank

migrateup:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/core_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/core_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open -a "Google Chrome" coverage.html

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test coverage
