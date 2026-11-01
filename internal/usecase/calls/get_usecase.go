package calls

import (
	"context"
	"log/slog"
	"telephone-service/internal/model"
)

type GetCallsUsecase struct {
	logger *slog.Logger
}

func NewGetCallsUsecase(logger *slog.Logger) *GetCallsUsecase {
	return &GetCallsUsecase{
		logger: logger,
	}
}

func (uc *GetCallsUsecase) Execute(_ context.Context) ([]model.Call, error) {
	return []model.Call{}, nil
}
