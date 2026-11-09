package disconnect_call_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_DisconnectsSubscriber(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockDisconnectCallService(gomock.NewController(t))
	usecase := NewDisconnectCallUsecase(service)
	expected := model.Call{Id: 42, Status: model.CallStatusActive, Version: 2}
	service.EXPECT().DisconnectCall(ctx, 42, 7).Return(expected, nil)

	// act
	call, err := usecase.Execute(ctx, 42, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}
