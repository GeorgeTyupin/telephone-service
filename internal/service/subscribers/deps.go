package subscribers

import (
	"context"
	"telephone-service/internal/model"
)

type SubscribersRepository interface {
	GetSubscribers(ctx context.Context) ([]model.Subscriber, error)
	GetSubscriber(ctx context.Context, id int) (model.Subscriber, error)
	GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error)
}
