package get_subscribers_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscribersUsecase struct {
	subscribersService GetSubscribersService
}

func NewGetSubscribersUsecase(subscribersService GetSubscribersService) *GetSubscribersUsecase {
	return &GetSubscribersUsecase{subscribersService: subscribersService}
}

func (uc *GetSubscribersUsecase) Execute(ctx context.Context) ([]model.Subscriber, error) {
	return uc.subscribersService.GetSubscribers(ctx)
}
