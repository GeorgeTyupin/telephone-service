package get_subscriber_calls_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberCallsUsecase struct {
	callsService GetSubscriberCallsService
}

func NewGetSubscriberCallsUsecase(callsService GetSubscriberCallsService) *GetSubscriberCallsUsecase {
	return &GetSubscriberCallsUsecase{callsService: callsService}
}

func (uc *GetSubscriberCallsUsecase) Execute(ctx context.Context, id int) ([]model.Call, error) {
	return uc.callsService.GetSubscriberCalls(ctx, id)
}
