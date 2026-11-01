package subscribers_handlers

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscribersUsecase interface {
	Execute(ctx context.Context) ([]model.Subscriber, error)
}
