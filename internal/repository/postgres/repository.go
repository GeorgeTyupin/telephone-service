package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger *slog.Logger
	dbPool *pgxpool.Pool
}

func NewRepository(logger *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{
		logger: logger,
		dbPool: db,
	}
}
