package pool

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, connString string) *pgxpool.Pool {
	const op = "pool.NewPgxPool"
	logger := slog.Default().With(slog.String("operation", op))

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Error("Не удалось создать пул", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return pool
}
