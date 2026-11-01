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

	if err := run(); err != nil {
		slog.Default().With(slog.String("operation", op)).Error("Ошибка работы приложения", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	c, err := container.NewContainer(ctx)
	if err != nil {
		return err
	}
	defer c.Stop()

	return c.Run(ctx)
}
