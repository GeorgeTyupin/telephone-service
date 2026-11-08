package get_subscriber_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberService interface {
	GetSubscriber(ctx context.Context, id int) (model.Subscriber, error)
}
