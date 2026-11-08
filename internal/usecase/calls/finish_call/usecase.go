package finish_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type FinishCallUsecase struct {
	callsService FinishCallService
}

func NewFinishCallUsecase(callsService FinishCallService) *FinishCallUsecase {
	return &FinishCallUsecase{callsService: callsService}
}

func (uc *FinishCallUsecase) Execute(ctx context.Context, id int) (model.Call, error) {
	return uc.callsService.FinishCall(ctx, id)
}
