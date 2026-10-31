package migrator

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationsPath = "migrations"

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		stdlib.OpenDBFromPool(pool),
		os.DirFS(migrationsPath),
	)
	if err != nil {
		return fmt.Errorf("не получилось инициализировать провайдер: %w", err)
	}
	defer provider.Close()

	_, err = provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	return nil
}
