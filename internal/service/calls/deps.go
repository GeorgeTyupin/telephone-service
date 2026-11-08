package calls

import (
	"context"
	"telephone-service/internal/model"
)

type CallsRepository interface {
	GetCalls(ctx context.Context, filter model.CallsFilter) ([]model.Call, error)
	GetCall(ctx context.Context, id int) (model.Call, error)
	StartCall(ctx context.Context, call model.Call) (model.Call, error)
	UpdateCall(ctx context.Context, update model.CallUpdate) (model.Call, error)
	GetSubscriber(ctx context.Context, id int) (model.Subscriber, error)
	GetSubscriberCalls(ctx context.Context, id int) ([]model.Call, error)
}
