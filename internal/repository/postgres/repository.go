package postgres

import (
	"log/slog"
	"telephone-service/internal/repository/postgres/record"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger  *slog.Logger
	dbPool  *pgxpool.Pool
	queries *record.Queries
}

func NewRepository(logger *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{
		logger:  logger,
		dbPool:  db,
		queries: record.New(db),
	}
}
