package get_subscriber_calls_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberCallsService interface {
	GetSubscriberCalls(ctx context.Context, id int) ([]model.Call, error)
}
