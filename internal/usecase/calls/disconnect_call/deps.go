package disconnect_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type DisconnectCallService interface {
	DisconnectCall(ctx context.Context, id, subscriberId int) (model.Call, error)
}
