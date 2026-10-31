package calls

import (
	"context"
	"log/slog"
	"telephone-service/internal/domain"
)

type GetCallsUsecase struct {
	logger *slog.Logger
}

func NewGetCallsUsecase(logger *slog.Logger) *GetCallsUsecase {
	return &GetCallsUsecase{
		logger: logger,
	}
}

func (uc *GetCallsUsecase) Execute(_ context.Context) ([]domain.Call, error) {
	return []domain.Call{}, nil
}
