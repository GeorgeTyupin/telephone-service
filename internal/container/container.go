package container

import (
	"context"
	"fmt"
	"log/slog"
	"telephone-service/internal/config"
	"telephone-service/internal/db/postgres/migrator"
	"telephone-service/internal/db/postgres/pool"
	"telephone-service/internal/transport/http/handlers/calls"
	callsHandler "telephone-service/internal/transport/http/handlers/subscribers"
	getCallsUC "telephone-service/internal/usecase/calls"

	"github.com/jackc/pgx/v5/pgxpool"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Container struct {
	config config.Config
	logger *slog.Logger

	dbPool *pgxpool.Pool

	meterProvider  *sdkmetric.MeterProvider
	tracerProvider *sdktrace.TracerProvider

	//handlers
	callsHandler       *calls.CallsHandler
	subscribersHandler *callsHandler.SubscribersHandler

	//usecases
	getCallsUsecase *getCallsUC.GetCallsUsecase
}

func NewContainer(ctx context.Context) (*Container, error) {
	cfg := config.MustLoad()
	c := &Container{
		config: cfg,
		logger: slog.Default(),
	}
	ready := false
	defer func() {
		if !ready {
			c.Stop()
		}
	}()
	if err := c.initTelemetry(ctx); err != nil {
		return nil, fmt.Errorf("инициализация телеметрии: %w", err)
	}

	poolInitCtx, cancel := context.WithTimeout(ctx, cfg.Database.InitTimeout)
	defer cancel()
	dbPool, err := pool.NewPgxPool(poolInitCtx, cfg.Database.GetConnectionString())
	if err != nil {
		return nil, err
	}
	c.dbPool = dbPool

	migrationsCtx, cancel := context.WithTimeout(ctx, cfg.Database.MigrationsTimeout)
	defer cancel()
	if err := migrator.Run(migrationsCtx, c.dbPool); err != nil {
		return nil, fmt.Errorf("применение миграций: %w", err)
	}

	ready = true
	return c, nil
}

func (c *Container) Run(ctx context.Context) error {
	apiServer, err := c.BuildApi()
	if err != nil {
		return err
	}
	return apiServer.Run(ctx)
}

func (c *Container) Stop() {
	if c.dbPool != nil {
		c.dbPool.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Telemetry.ShutdownTimeout)
	defer cancel()
	if c.tracerProvider != nil {
		if err := c.tracerProvider.Shutdown(ctx); err != nil {
			c.logger.ErrorContext(ctx, "Не удалось завершить tracing", slog.Any("error", err))
		}
	}
	if c.meterProvider != nil {
		if err := c.meterProvider.Shutdown(ctx); err != nil {
			c.logger.ErrorContext(ctx, "Не удалось завершить metrics", slog.Any("error", err))
		}
	}
}
