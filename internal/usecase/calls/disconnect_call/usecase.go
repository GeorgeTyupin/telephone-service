package disconnect_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type DisconnectCallUsecase struct {
	callsService DisconnectCallService
}

func NewDisconnectCallUsecase(callsService DisconnectCallService) *DisconnectCallUsecase {
	return &DisconnectCallUsecase{callsService: callsService}
}

func (uc *DisconnectCallUsecase) Execute(ctx context.Context, id, subscriberId int) (model.Call, error) {
	return uc.callsService.DisconnectCall(ctx, id, subscriberId)
}
