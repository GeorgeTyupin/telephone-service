package get_subscriber_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ReturnsSubscriber(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockGetSubscriberService(gomock.NewController(t))
	usecase := NewGetSubscriberUsecase(service)
	expected := model.Subscriber{Id: 7, Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}}
	service.EXPECT().GetSubscriber(ctx, 7).Return(expected, nil)

	// act
	subscriber, err := usecase.Execute(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, subscriber)
}
