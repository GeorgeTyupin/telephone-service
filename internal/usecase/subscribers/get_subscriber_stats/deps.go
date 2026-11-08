package get_subscriber_stats_uc

import (
	"context"
	"telephone-service/internal/model"
)

type GetSubscriberStatsService interface {
	GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error)
}
