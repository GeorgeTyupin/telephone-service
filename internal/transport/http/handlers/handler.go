package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	api "telephone-service/codegen/oas"
	"telephone-service/internal/transport/http/handlers/calls"
	"telephone-service/internal/transport/http/handlers/subscribers"
)

var _ api.Handler = (*handler)(nil)

type handler struct {
	*calls.CallsHandler
	*subscribers.SubscribersHandler
}

type Option func(*http.ServeMux)

func NewHandler(
	logger *slog.Logger,
	callsHandler *calls.CallsHandler,
	subscribersHandler *subscribers.SubscribersHandler,
	opts ...Option,
) (http.Handler, error) {
	apiServer, err := api.NewServer(
		&handler{
			CallsHandler:       callsHandler,
			SubscribersHandler: subscribersHandler,
		},
		api.WithErrorHandler(NewErrorHandler(logger)),
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать API-сервер: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", apiServer)
	for _, opt := range opts {
		opt(mux)
	}
	return mux, nil
}

func WithPing() Option {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /ping", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}
}

func WithMetrics() Option {
	return func(mux *http.ServeMux) {
		mux.Handle("GET /metrics", promhttp.Handler())
	}
}
