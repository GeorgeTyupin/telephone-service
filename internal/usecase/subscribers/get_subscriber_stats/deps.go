package get_subscriber_stats_uc

import (
	"context"
	"telephone-service/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=deps.go -destination=deps_mock_test.go -package=$GOPACKAGE

type GetSubscriberStatsService interface {
	GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error)
}
