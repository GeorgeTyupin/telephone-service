package get_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetCallService interface {
	GetCall(ctx context.Context, id int) (model.Call, error)
}
