package get_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallUsecase struct {
	callsService GetCallService
}

func NewGetCallUsecase(callsService GetCallService) *GetCallUsecase {
	return &GetCallUsecase{callsService: callsService}
}

func (uc *GetCallUsecase) Execute(ctx context.Context, id int) (model.Call, error) {
	return uc.callsService.GetCall(ctx, id)
}
