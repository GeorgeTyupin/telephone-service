.DEFAULT_GOAL := help
.PHONY: help generate ogen_gen test build vet

help:
	@printf '%s\n' \
		'make ogen_gen  Сгенерировать Go-код из OpenAPI-контракта' \
		'make generate  Алиас команды make ogen_gen' \
		'make test      Запустить все Go-тесты' \
		'make build     Проверить сборку проекта' \
		'make vet       Запустить go vet'

ogen_gen:
	go generate ./...

generate: ogen_gen
test:
	go test ./...

build: ogen_gen
	go build ./...

run:
	go run ./...

vet:
	go vet ./...
