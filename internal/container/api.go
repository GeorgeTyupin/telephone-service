package container

import (
	"log/slog"

	httpserver "telephone-service/internal/transport/http"
	"telephone-service/internal/transport/http/handlers"
)

func (c *Container) BuildApi() (*httpserver.Server, error) {
	logger := c.logger.With(slog.String("component", "http"))

	var opts []handlers.Option
	opts = append(opts, handlers.WithPing())
	if c.config.Telemetry.Metrics.Enabled {
		opts = append(opts, handlers.WithMetrics())
	}

	handler, err := handlers.NewHandler(
		logger,
		c.ResolveCallsHandler(),
		c.ResolveSubscribersHandler(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return httpserver.NewServer(logger, c.config.Server, handler), nil
}
