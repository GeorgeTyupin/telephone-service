package subscribers

import (
	"context"
	"telephone-service/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=deps.go -destination=deps_mock_test.go -package=$GOPACKAGE

type SubscribersRepository interface {
	GetSubscribers(ctx context.Context) ([]model.Subscriber, error)
	GetSubscriber(ctx context.Context, id int) (model.Subscriber, error)
	GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error)
}
