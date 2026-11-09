package connect_call_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ConnectsSubscriber(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockConnectCallService(gomock.NewController(t))
	usecase := NewConnectCallUsecase(service)
	expected := model.Call{Id: 42, Status: model.CallStatusActive, Participants: []model.Subscriber{{Id: 7, Name: "Иван"}}}
	service.EXPECT().ConnectCall(ctx, 42, 7).Return(expected, nil)

	// act
	call, err := usecase.Execute(ctx, 42, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}
