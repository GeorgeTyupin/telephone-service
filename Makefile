.DEFAULT_GOAL := help
.PHONY: help generate ogen_gen sqlc_gen mocks_gen migrate test test-integration build run vet

POSTGRES_PORT ?= 5400

help:
	@printf '%s\n' \
		'make ogen_gen  Сгенерировать Go-код из OpenAPI-контракта' \
		'make sqlc_gen  Сгенерировать Go-код из SQL-запросов' \
		'make mocks_gen Сгенерировать тестовые моки' \
		'make generate  Запустить все генераторы' \
		'make migrate   Применить миграции к локальной PostgreSQL' \
		'make test      Запустить тесты без Docker' \
		'make test-integration  Запустить интеграционные тесты с PostgreSQL' \
		'make build     Проверить сборку проекта' \
		'make vet       Запустить go vet'

ogen_gen:
	go generate -run 'ogen' ./tools.go

sqlc_gen:
	go generate -run 'sqlc' ./tools.go

mocks_gen:
	go generate -run 'mockgen' ./internal/service/... ./internal/usecase/...

generate: ogen_gen sqlc_gen mocks_gen

migrate:
	@set -eu; set -a; . ./configs/.env; set +a; \
	PGUSER="$${POSTGRES_USER}" \
	PGPASSWORD="$${POSTGRES_PASSWORD}" \
	PGDATABASE="$${POSTGRES_DB}" \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="host=localhost port=$(POSTGRES_PORT) sslmode=disable" \
	go run github.com/pressly/goose/v3/cmd/goose@v3.29.0 -dir migrations up

test:
	go test ./...

test-integration:
	go test -race -count=1 -tags=integration -timeout=10m -v ./internal/tests/cases

build: generate
	go build ./...

run:
	go run ./cmd

vet:
	go vet ./...
