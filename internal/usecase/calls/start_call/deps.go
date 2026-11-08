package start_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type StartCallService interface {
	StartCall(ctx context.Context, call model.Call) (model.Call, error)
}
