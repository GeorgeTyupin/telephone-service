package start_call_uc

import (
	"testing"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecute_StartsCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	service := NewMockStartCallService(gomock.NewController(t))
	usecase := NewStartCallUsecase(service)
	input := model.Call{Participants: []model.Subscriber{{Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}}}}
	expected := input
	expected.Id = 42
	expected.Status = model.CallStatusActive
	service.EXPECT().StartCall(ctx, input).Return(expected, nil)

	// act
	call, err := usecase.Execute(ctx, input)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}
