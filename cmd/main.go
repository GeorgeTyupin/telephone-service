package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"telephone-service/internal/container"
)

func main() {
	const op = "main.main"
	logger := slog.Default().With(slog.String("operation", op))

	if err := run(); err != nil {
		logger.Error("Ошибка работы приложения", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	c := container.NewContainer(ctx)
	defer c.Stop()

	server, err := c.BuildApi()
	if err != nil {
		return err
	}

	return server.Run(ctx)
}
