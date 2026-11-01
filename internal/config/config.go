package config

import (
	"fmt"
	"time"
)

type Config struct {
	Server    Server    `path:",squash"`
	Database  database  `path:",squash"`
	Telemetry Telemetry `path:",squash"`
}

type database struct {
	User              string        `path:"db.user"`
	Name              string        `path:"db.name"`
	Password          string        `path:"db.password"`
	Port              string        `path:"db.port"`
	InitTimeout       time.Duration `path:"db.init_timeout"`
	MigrationsTimeout time.Duration `path:"db.migrations_timeout"`
}

func (d database) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", d.User, d.Password, d.Port, d.Name)
}

type Server struct {
	Addr              string        `path:"server.addr"`
	ReadHeaderTimeout time.Duration `path:"server.read_header_timeout"`
	ShutdownTimeout   time.Duration `path:"server.shutdown_timeout"`
}

type Telemetry struct {
	Environment     string        `path:"telemetry.environment"`
	ShutdownTimeout time.Duration `path:"telemetry.shutdown_timeout"`
	Logs            Logs          `path:",squash"`
	Metrics         Metrics       `path:",squash"`
	Tracing         Tracing       `path:",squash"`
}

type Logs struct {
	Level  string `path:"telemetry.logs.level"`
	Format string `path:"telemetry.logs.format"`
}

type Metrics struct {
	Enabled bool `path:"telemetry.metrics.enabled"`
}

type Tracing struct {
	Enabled     bool    `path:"telemetry.tracing.enabled"`
	Endpoint    string  `path:"telemetry.tracing.endpoint"`
	SampleRatio float64 `path:"telemetry.tracing.sample_ratio"`
}
