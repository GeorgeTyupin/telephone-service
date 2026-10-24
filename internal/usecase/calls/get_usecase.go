package calls

import (
	"log/slog"
	"net/url"
	"telephone-service/internal/dto"
)

type GetCallsUsecase struct {
	logger *slog.Logger
}

func NewGetCallsUsecase(logger *slog.Logger) *GetCallsUsecase {
	return &GetCallsUsecase{
		logger: logger,
	}
}

func (uc *GetCallsUsecase) Execute(params url.Values) ([]dto.Call, error) {
	return []dto.Call{}, nil
}
