package mapper

import "log/slog"

type Mapper struct {
	logger *slog.Logger
}

func NewMapper(logger *slog.Logger) *Mapper {
	return &Mapper{
		logger: logger,
	}
}
