package get_calls_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallsUsecase struct {
	callsService GetCallsService
}

func NewGetCallsUsecase(callsService GetCallsService) *GetCallsUsecase {
	return &GetCallsUsecase{callsService: callsService}
}

func (uc *GetCallsUsecase) Execute(ctx context.Context) ([]model.Call, error) {
	return uc.callsService.GetCalls(ctx)
}
