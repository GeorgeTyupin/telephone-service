package get_calls_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_ReturnsFilteredCalls(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockGetCallsService(gomock.NewController(t))
	usecase := NewGetCallsUsecase(service)
	filter := model.CallsFilter{Status: model.CallStatusActive}
	expected := []model.Call{{Id: 42, Status: model.CallStatusActive}}
	service.EXPECT().GetCalls(ctx, filter).Return(expected, nil)

	// act
	calls, err := usecase.Execute(ctx, filter)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, calls)
}
