package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	envPath    = "configs/.env"
	configPath = "configs/config.yaml"
)

func Load() (Config, error) {
	loader := koanf.New(".")
	if err := loader.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("не удалось прочитать YAML-конфигурацию: %w", err)
	}

	if err := loader.Load(file.Provider(envPath), dotenv.ParserEnv("", ".", environmentPath)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("не удалось прочитать .env: %w", err)
	}

	if err := loader.Load(env.Provider("", ".", environmentPath), nil); err != nil {
		return Config{}, fmt.Errorf("не удалось прочитать переменные окружения: %w", err)
	}

	var config Config
	if err := loader.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "path", FlatPaths: true}); err != nil {
		return Config{}, fmt.Errorf("не удалось преобразовать конфигурацию: %w", err)
	}

	return config, nil
}

func environmentPath(key string) string {
	switch key {
	case "POSTGRES_USER":
		return "db.user"
	case "POSTGRES_DB":
		return "db.name"
	case "POSTGRES_PASSWORD":
		return "db.password"
	case "TELEMETRY_ENVIRONMENT":
		return "telemetry.environment"
	case "LOG_LEVEL":
		return "telemetry.logs.level"
	case "LOG_FORMAT":
		return "telemetry.logs.format"
	case "METRICS_ENABLED":
		return "telemetry.metrics.enabled"
	case "TRACING_ENABLED":
		return "telemetry.tracing.enabled"
	case "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":
		return "telemetry.tracing.endpoint"
	case "TRACING_SAMPLE_RATIO":
		return "telemetry.tracing.sample_ratio"
	}

	if strings.Contains(key, "__") {
		return strings.ToLower(strings.ReplaceAll(key, "__", "."))
	}

	return ""
}
