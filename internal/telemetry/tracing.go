package telemetry

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"telephone-service/internal/config"
)

func NewTracer(ctx context.Context, cfg config.Telemetry) (*sdktrace.TracerProvider, error) {
	if !cfg.Tracing.Enabled {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil, nil
	}
	if math.IsNaN(cfg.Tracing.SampleRatio) || cfg.Tracing.SampleRatio < 0 || cfg.Tracing.SampleRatio > 1 {
		return nil, fmt.Errorf("sample_ratio должен быть в диапазоне [0, 1]")
	}
	endpoint, err := url.Parse(cfg.Tracing.Endpoint)
	if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, fmt.Errorf("endpoint трейсов должен быть HTTP(S) URL без credentials, query и fragment")
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(cfg.Tracing.Endpoint),
		otlptracehttp.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("создание OTLP exporter: %w", err)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Tracing.SampleRatio))),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(provider)
	return provider, nil
}
