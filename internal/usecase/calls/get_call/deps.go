package get_call_uc

import (
	"context"
	"telephone-service/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=deps.go -destination=deps_mock_test.go -package=$GOPACKAGE

type GetCallService interface {
	GetCall(ctx context.Context, id int) (model.Call, error)
}
