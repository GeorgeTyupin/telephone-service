package container

import (
	"context"
	"log/slog"
	"os"
	"telephone-service/internal/config"
	"telephone-service/internal/db/postgres/migrator"
	"telephone-service/internal/db/postgres/pool"
	"telephone-service/internal/transport/http/handlers/calls"
	callsHandler "telephone-service/internal/transport/http/handlers/subscribers"
	getCallsUC "telephone-service/internal/usecase/calls"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	config config.Config
	logger *slog.Logger

	dbPool *pgxpool.Pool

	//handlers
	callsHandler       *calls.CallsHandler
	subscribersHandler *callsHandler.SubscribersHandler

	//usecases
	getCallsUsecase *getCallsUC.GetCallsUsecase
}

func NewContainer(ctx context.Context) *Container {
	const op = "container.NewContainer"
	baseLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger := baseLogger.With(slog.String("operation", op))

	cfg := config.MustLoad()

	poolInitCtx, cancel := context.WithTimeout(ctx, cfg.Database.InitTimeout)
	defer cancel()
	dbPool := pool.NewPgxPool(poolInitCtx, cfg.Database.GetConnectionString())

	migrationsCtx, cancel := context.WithTimeout(ctx, cfg.Database.MigrationsTimeout)
	defer cancel()
	if err := migrator.Run(migrationsCtx, dbPool); err != nil {
		logger.Error("Ошибка выполнения миграции", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return &Container{
		config: cfg,
		logger: baseLogger,

		dbPool: dbPool,
	}
}

func (c *Container) Stop() {
	c.dbPool.Close()
}
