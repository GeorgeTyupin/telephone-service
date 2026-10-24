package app

import (
	"context"
	"log/slog"
	"net/http"
	"telephone-service/internal/config"
)

type App struct {
	addr   string
	logger *slog.Logger

	//handlers
	callsHandler    CallsHandler
	abonentsHandler AbonentsHandler
}

func NewApp(
	logger *slog.Logger,
	cfg config.Config,
	callsHandler CallsHandler,
	// abonentsHandler AbonentsHandler,
) *App {
	return &App{
		logger:       logger,
		addr:         cfg.Server.Addr,
		callsHandler: callsHandler,
		//abonentsHandler: abonentsHandler,
	}
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)
	mux := a.registerRoutes()

	go func() {
		if err := http.ListenAndServe(a.addr, mux); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}
	return nil
}
