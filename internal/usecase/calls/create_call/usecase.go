package create_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type CreateCallUsecase struct {
	callsService CreateCallService
}

func NewCreateCallUsecase(callsService CreateCallService) *CreateCallUsecase {
	return &CreateCallUsecase{callsService: callsService}
}

func (uc *CreateCallUsecase) Execute(ctx context.Context, call model.Call) (model.Call, error) {
	return uc.callsService.CreateCall(ctx, call)
}
