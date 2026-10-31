package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	api "telephone-service/codegen/oas"
	"telephone-service/internal/config"
	"telephone-service/internal/transport/http/handlers"
)

type Server struct {
	logger          *slog.Logger
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(logger *slog.Logger, cfg config.Server, handler api.Handler) (*Server, error) {
	apiServer, err := api.NewServer(
		handler,
		api.WithErrorHandler(handlers.NewErrorHandler(logger)),
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать API-сервер: %w", err)
	}

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:              cfg.Addr,
			Handler:           apiServer,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	defer s.shutdown(ctx)

	serveErrors := make(chan error, 1)
	go func() {
		defer close(serveErrors)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrors <- fmt.Errorf("не удалось запустить HTTP-сервер: %w", err)
			return
		}
		serveErrors <- nil
	}()
	select {
	case err := <-serveErrors:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		closeErr := s.server.Close()
		if closeErr != nil {
			s.logger.ErrorContext(shutdownCtx, "Не удалось принудительно закрыть HTTP-сервер", slog.String("error", closeErr.Error()))
		}
		s.logger.ErrorContext(shutdownCtx, "Не удалось завершить HTTP-сервер", slog.String("error", err.Error()))
	}
}
