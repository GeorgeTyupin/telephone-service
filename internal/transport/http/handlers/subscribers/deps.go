package subscribers_handlers

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscribersUsecase interface {
	Execute(ctx context.Context) ([]model.Subscriber, error)
}

type GetSubscriberUsecase interface {
	Execute(ctx context.Context, id int) (model.Subscriber, error)
}

type GetSubscriberCallsUsecase interface {
	Execute(ctx context.Context, id int) ([]model.Call, error)
}

type GetSubscriberStatsUsecase interface {
	Execute(ctx context.Context, id int) (model.SubscriberStats, error)
}
