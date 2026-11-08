package connect_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type ConnectCallUsecase struct {
	callsService ConnectCallService
}

func NewConnectCallUsecase(callsService ConnectCallService) *ConnectCallUsecase {
	return &ConnectCallUsecase{callsService: callsService}
}

func (uc *ConnectCallUsecase) Execute(ctx context.Context, id, subscriberId int) (model.Call, error) {
	return uc.callsService.ConnectCall(ctx, id, subscriberId)
}
