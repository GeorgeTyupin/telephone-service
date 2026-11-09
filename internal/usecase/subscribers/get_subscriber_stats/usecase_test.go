package get_subscriber_stats_uc

import (
	"testing"
	"time"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ReturnsSubscriberStats(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockGetSubscriberStatsService(gomock.NewController(t))
	usecase := NewGetSubscriberStatsUsecase(service)
	expected := model.SubscriberStats{
		CallsCount:      3,
		TotalDuration:   model.Duration(6 * time.Minute),
		AverageDuration: model.Duration(2 * time.Minute),
	}
	service.EXPECT().GetSubscriberStats(ctx, 7).Return(expected, nil)

	// act
	stats, err := usecase.Execute(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}
