package create_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type CreateCallService interface {
	CreateCall(ctx context.Context, call model.Call) (model.Call, error)
}
