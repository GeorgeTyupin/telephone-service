package telemetry

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"telephone-service/internal/config"
)

func NewMetrics(cfg config.Telemetry) (*sdkmetric.MeterProvider, error) {
	if !cfg.Metrics.Enabled {
		otel.SetMeterProvider(noop.NewMeterProvider())
		return nil, nil
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("создание Prometheus exporter: %w", err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithView(sdkmetric.NewView(
			sdkmetric.Instrument{Name: "ogen.server.duration", Kind: sdkmetric.InstrumentKindHistogram},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000},
			}},
		)),
	)
	otel.SetMeterProvider(provider)
	return provider, nil
}
