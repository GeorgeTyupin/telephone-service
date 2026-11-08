package get_subscriber_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberUsecase struct {
	subscribersService GetSubscriberService
}

func NewGetSubscriberUsecase(subscribersService GetSubscriberService) *GetSubscriberUsecase {
	return &GetSubscriberUsecase{subscribersService: subscribersService}
}

func (uc *GetSubscriberUsecase) Execute(ctx context.Context, id int) (model.Subscriber, error) {
	return uc.subscribersService.GetSubscriber(ctx, id)
}
