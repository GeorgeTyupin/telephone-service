package container

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"

	"telephone-service/internal/telemetry"
)

func (c *Container) initTelemetry(ctx context.Context) error {
	cfg := c.config.Telemetry
	logger, err := telemetry.NewLogger(cfg, os.Stdout)
	if err != nil {
		return err
	}
	c.logger = logger
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("Ошибка OpenTelemetry", slog.Any("error", err))
	}))

	c.meterProvider, err = telemetry.NewMetrics(cfg)
	if err != nil {
		return err
	}
	c.tracerProvider, err = telemetry.NewTracer(ctx, cfg)
	return err
}
