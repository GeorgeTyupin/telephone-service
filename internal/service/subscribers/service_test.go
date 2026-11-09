package subscribers

import (
	"testing"
	"time"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSubscribers_ReturnsSubscribers(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockSubscribersRepository(gomock.NewController(t))
	service := NewService(repository)
	expected := []model.Subscriber{{Id: 7, Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}}}
	repository.EXPECT().GetSubscribers(ctx).Return(expected, nil)

	// act
	subscribers, err := service.GetSubscribers(ctx)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, subscribers)
}

func TestGetSubscriber_ReturnsSubscriber(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockSubscribersRepository(gomock.NewController(t))
	service := NewService(repository)
	expected := model.Subscriber{Id: 7, Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}}
	repository.EXPECT().GetSubscriber(ctx, 7).Return(expected, nil)

	// act
	subscriber, err := service.GetSubscriber(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, subscriber)
}

func TestGetSubscriberStats_ReturnsStats(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockSubscribersRepository(gomock.NewController(t))
	service := NewService(repository)
	expected := model.SubscriberStats{
		CallsCount:      3,
		TotalDuration:   model.Duration(6 * time.Minute),
		AverageDuration: model.Duration(2 * time.Minute),
	}
	repository.EXPECT().GetSubscriberStats(ctx, 7).Return(expected, nil)

	// act
	stats, err := service.GetSubscriberStats(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, stats)
}
