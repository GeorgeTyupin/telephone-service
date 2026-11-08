package get_subscriber_stats_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberStatsUsecase struct {
	subscribersService GetSubscriberStatsService
}

func NewGetSubscriberStatsUsecase(subscribersService GetSubscriberStatsService) *GetSubscriberStatsUsecase {
	return &GetSubscriberStatsUsecase{subscribersService: subscribersService}
}

func (uc *GetSubscriberStatsUsecase) Execute(ctx context.Context, id int) (model.SubscriberStats, error) {
	return uc.subscribersService.GetSubscriberStats(ctx, id)
}
