package connect_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

type ConnectCallService interface {
	ConnectCall(ctx context.Context, id, subscriberId int) (model.Call, error)
}
