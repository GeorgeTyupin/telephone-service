package subscribers

import (
	"context"
	"telephone-service/internal/model"
)

type SubscribersRepository interface {
	GetSubscribers(ctx context.Context) ([]model.Subscriber, error)
}
