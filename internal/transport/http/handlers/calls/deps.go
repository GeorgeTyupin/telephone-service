package calls_handlers

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallsUsecase interface {
	Execute(ctx context.Context) ([]model.Call, error)
}

type CreateCallUsecase interface {
	Execute(ctx context.Context, call model.Call) (model.Call, error)
}
