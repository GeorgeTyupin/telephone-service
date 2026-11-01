package get_subscribers_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscribersService interface {
	GetSubscribers(ctx context.Context) ([]model.Subscriber, error)
}
