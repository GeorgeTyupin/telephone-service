package start_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type StartCallUsecase struct {
	callsService StartCallService
}

func NewStartCallUsecase(callsService StartCallService) *StartCallUsecase {
	return &StartCallUsecase{callsService: callsService}
}

func (uc *StartCallUsecase) Execute(ctx context.Context, call model.Call) (model.Call, error) {
	return uc.callsService.StartCall(ctx, call)
}
