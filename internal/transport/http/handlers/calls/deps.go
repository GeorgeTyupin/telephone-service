package calls_handlers

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallsUsecase interface {
	Execute(ctx context.Context, filter model.CallsFilter) ([]model.Call, error)
}

type GetCallUsecase interface {
	Execute(ctx context.Context, id int) (model.Call, error)
}

type StartCallUsecase interface {
	Execute(ctx context.Context, call model.Call) (model.Call, error)
}

type ConnectCallUsecase interface {
	Execute(ctx context.Context, id, subscriberId int) (model.Call, error)
}

type DisconnectCallUsecase interface {
	Execute(ctx context.Context, id, subscriberId int) (model.Call, error)
}

type FinishCallUsecase interface {
	Execute(ctx context.Context, id int) (model.Call, error)
}
