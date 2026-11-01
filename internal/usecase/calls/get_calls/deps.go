package get_calls_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallsService interface {
	GetCalls(ctx context.Context) ([]model.Call, error)
}
