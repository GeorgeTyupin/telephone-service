package migrator

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

func Run(ctx context.Context, pool *pgxpool.Pool, migrations fs.FS) error {
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return fmt.Errorf("не удалось создать блокировку миграций: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		stdlib.OpenDBFromPool(pool),
		migrations,
		goose.WithSessionLocker(locker),
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
