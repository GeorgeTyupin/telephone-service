package finish_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type FinishCallService interface {
	FinishCall(ctx context.Context, id int) (model.Call, error)
}
