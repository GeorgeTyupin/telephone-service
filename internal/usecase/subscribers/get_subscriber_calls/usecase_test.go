package get_subscriber_calls_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ReturnsSubscriberCalls(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockGetSubscriberCallsService(gomock.NewController(t))
	usecase := NewGetSubscriberCallsUsecase(service)
	expected := []model.Call{{Id: 42, Status: model.CallStatusFinished, Participants: []model.Subscriber{{Id: 7, Name: "Иван"}}}}
	service.EXPECT().GetSubscriberCalls(ctx, 7).Return(expected, nil)

	// act
	calls, err := usecase.Execute(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, calls)
}
