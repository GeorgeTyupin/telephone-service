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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	c := container.NewContainer(ctx)

	if err := c.App.Run(ctx); err != nil {
		logger.Error("Ошибка запуска приложения", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
