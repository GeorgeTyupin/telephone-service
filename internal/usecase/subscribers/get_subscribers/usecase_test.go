package get_subscribers_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ReturnsSubscribers(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockGetSubscribersService(gomock.NewController(t))
	usecase := NewGetSubscribersUsecase(service)
	expected := []model.Subscriber{{Id: 7, Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}}}
	service.EXPECT().GetSubscribers(ctx).Return(expected, nil)

	// act
	subscribers, err := usecase.Execute(ctx)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, subscribers)
}
