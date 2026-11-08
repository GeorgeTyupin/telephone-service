package config

import (
	"net"
	"net/url"
	"time"
)

type Config struct {
	Server    Server    `path:",squash"`
	Database  database  `path:",squash"`
	Telemetry Telemetry `path:",squash"`
}

type database struct {
	Host              string        `path:"db.host"`
	User              string        `path:"db.user"`
	Name              string        `path:"db.name"`
	Password          string        `path:"db.password"`
	Port              string        `path:"db.port"`
	InitTimeout       time.Duration `path:"db.init_timeout"`
	MigrationsTimeout time.Duration `path:"db.migrations_timeout"`
}

func (d database) GetConnectionString() string {
	connectionURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     net.JoinHostPort(d.Host, d.Port),
		Path:     "/" + d.Name,
		RawPath:  "/" + url.PathEscape(d.Name),
		RawQuery: "sslmode=disable",
	}

	return connectionURL.String()
}

type Server struct {
	Addr              string        `path:"server.addr"`
	ReadHeaderTimeout time.Duration `path:"server.read_header_timeout"`
	ReadTimeout       time.Duration `path:"server.read_timeout"`
	WriteTimeout      time.Duration `path:"server.write_timeout"`
	IdleTimeout       time.Duration `path:"server.idle_timeout"`
	MaxBodyBytes      int64         `path:"server.max_body_bytes"`
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
