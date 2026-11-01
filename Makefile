.DEFAULT_GOAL := help
.PHONY: help generate ogen_gen sqlc_gen test build run vet

help:
	@printf '%s\n' \
		'make ogen_gen  Сгенерировать Go-код из OpenAPI-контракта' \
		'make sqlc_gen  Сгенерировать Go-код из SQL-запросов' \
		'make generate  Запустить оба генератора' \
		'make test      Запустить все Go-тесты' \
		'make build     Проверить сборку проекта' \
		'make vet       Запустить go vet'

ogen_gen:
	go generate -run 'ogen' ./tools.go

sqlc_gen:
	go generate -run 'sqlc' ./tools.go

generate: ogen_gen sqlc_gen

test:
	go test ./...

build: generate
	go build ./...

run:
	go run ./...

vet:
	go vet ./...
