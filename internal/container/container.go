package container

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"telephone-service/internal/config"
	"telephone-service/internal/db/postgres/migrator"
	"telephone-service/internal/db/postgres/pool"
	"telephone-service/internal/repository/postgres"
	callsservice "telephone-service/internal/service/calls"
	subscribersservice "telephone-service/internal/service/subscribers"
	"telephone-service/internal/transport/http/handlers/calls"
	"telephone-service/internal/transport/http/handlers/subscribers"
	"telephone-service/internal/usecase/calls/connect_call"
	"telephone-service/internal/usecase/calls/disconnect_call"
	"telephone-service/internal/usecase/calls/finish_call"
	"telephone-service/internal/usecase/calls/get_call"
	"telephone-service/internal/usecase/calls/get_calls"
	"telephone-service/internal/usecase/calls/start_call"
	"telephone-service/internal/usecase/subscribers/get_subscriber"
	"telephone-service/internal/usecase/subscribers/get_subscriber_calls"
	"telephone-service/internal/usecase/subscribers/get_subscriber_stats"
	"telephone-service/internal/usecase/subscribers/get_subscribers"

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

	postgresRepository *postgres.Repository

	callsService       *callsservice.Service
	subscribersService *subscribersservice.Service

	callsHandler       *calls_handlers.CallsHandler
	subscribersHandler *subscribers_handlers.SubscribersHandler

	getCallsUsecase           *get_calls_uc.GetCallsUsecase
	getCallUsecase            *get_call_uc.GetCallUsecase
	startCallUsecase          *start_call_uc.StartCallUsecase
	connectCallUsecase        *connect_call_uc.ConnectCallUsecase
	disconnectCallUsecase     *disconnect_call_uc.DisconnectCallUsecase
	finishCallUsecase         *finish_call_uc.FinishCallUsecase
	getSubscriberUsecase      *get_subscriber_uc.GetSubscriberUsecase
	getSubscriberCallsUsecase *get_subscriber_calls_uc.GetSubscriberCallsUsecase
	getSubscriberStatsUsecase *get_subscriber_stats_uc.GetSubscriberStatsUsecase
	getSubscribersUsecase     *get_subscribers_uc.GetSubscribersUsecase
}

func NewContainerWithDependencies(cfg config.Config, logger *slog.Logger, dbPool *pgxpool.Pool) *Container {
	return &Container{
		config: cfg,
		logger: logger,
		dbPool: dbPool,
	}
}

func NewContainer(ctx context.Context) (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("загрузка конфигурации: %w", err)
	}

	c := NewContainerWithDependencies(cfg, slog.Default(), nil)
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
	if err := migrator.Run(migrationsCtx, c.dbPool, os.DirFS("migrations")); err != nil {
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
