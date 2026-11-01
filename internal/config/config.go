package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	envPath    = "configs/.env"
	configPath = "configs/config.yaml"
)

type Config struct {
	Server    Server    `yaml:"server"`
	Database  database  `yaml:"db"`
	Telemetry Telemetry `yaml:"telemetry"`
}

type database struct {
	User              string        `env:"POSTGRES_USER"`
	Name              string        `env:"POSTGRES_DB"`
	Password          string        `env:"POSTGRES_PASSWORD"`
	Port              string        `yaml:"port"`
	InitTimeout       time.Duration `yaml:"init_timeout"`
	MigrationsTimeout time.Duration `yaml:"migrations_timeout"`
}

func (d database) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", d.User, d.Password, d.Port, d.Name)
}

type Server struct {
	Addr              string        `yaml:"addr"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
}

type Telemetry struct {
	Environment     string        `yaml:"environment" env:"TELEMETRY_ENVIRONMENT"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	Logs            Logs          `yaml:"logs"`
	Metrics         Metrics       `yaml:"metrics"`
	Tracing         Tracing       `yaml:"tracing"`
}

type Logs struct {
	Level  string `yaml:"level" env:"LOG_LEVEL"`
	Format string `yaml:"format" env:"LOG_FORMAT"`
}

type Metrics struct {
	Enabled bool `yaml:"enabled" env:"METRICS_ENABLED"`
}

type Tracing struct {
	Enabled     bool    `yaml:"enabled" env:"TRACING_ENABLED"`
	Endpoint    string  `yaml:"endpoint" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	SampleRatio float64 `yaml:"sample_ratio" env:"TRACING_SAMPLE_RATIO"`
}

func MustLoad() Config {
	const op = "config.MustLoad"
	logger := slog.Default().With(slog.String("operation", op))

	err := godotenv.Load(envPath)
	if err != nil {
		logger.Error("Ошибка загрузки .env файла", slog.String("error", err.Error()))
		os.Exit(1)
	}

	var config Config
	err = cleanenv.ReadEnv(&config)
	if err != nil {
		logger.Error("Ошибка чтения .end файла", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		logger.Error("Ошибка чтения конфига", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return config
}
