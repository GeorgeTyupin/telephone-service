package tools

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@v1.24.0 --config ./api/ogen.yaml --target ./codegen/oas --package api --clean ./api/openapi.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate
