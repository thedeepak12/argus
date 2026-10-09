include .env
export

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

db-connect:
	psql "$(DATABASE_URL)"

db-dump:
	pg_dump "$(DATABASE_URL)" > dump.sql

sqlc-gen:
	sqlc generate