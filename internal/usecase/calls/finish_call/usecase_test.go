package finish_call_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_FinishesCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockFinishCallService(gomock.NewController(t))
	usecase := NewFinishCallUsecase(service)
	expected := model.Call{Id: 42, Status: model.CallStatusFinished}
	service.EXPECT().FinishCall(ctx, 42).Return(expected, nil)

	// act
	call, err := usecase.Execute(ctx, 42)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}
