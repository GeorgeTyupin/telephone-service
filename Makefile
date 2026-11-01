.DEFAULT_GOAL := help
.PHONY: help generate ogen_gen sqlc_gen migrate test build run vet

POSTGRES_PORT ?= 5400

help:
	@printf '%s\n' \
		'make ogen_gen  Сгенерировать Go-код из OpenAPI-контракта' \
		'make sqlc_gen  Сгенерировать Go-код из SQL-запросов' \
		'make generate  Запустить оба генератора' \
		'make migrate   Применить миграции к локальной PostgreSQL' \
		'make test      Запустить все Go-тесты' \
		'make build     Проверить сборку проекта' \
		'make vet       Запустить go vet'

ogen_gen:
	go generate -run 'ogen' ./tools.go

sqlc_gen:
	go generate -run 'sqlc' ./tools.go

generate: ogen_gen sqlc_gen

migrate:
	@set -eu; set -a; . ./configs/.env; set +a; \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@localhost:$(POSTGRES_PORT)/$${POSTGRES_DB}?sslmode=disable" \
	go run github.com/pressly/goose/v3/cmd/goose@v3.29.0 -dir migrations up

test:
	go test ./...

build: generate
	go build ./...

run:
	go run ./...

vet:
	go vet ./...
