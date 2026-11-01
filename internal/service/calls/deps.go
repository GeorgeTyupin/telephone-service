package calls

import (
	"context"
	"telephone-service/internal/model"
)

type CallsRepository interface {
	GetCalls(ctx context.Context) ([]model.Call, error)
	CreateCall(ctx context.Context, call model.Call) (model.Call, error)
}
