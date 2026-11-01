package calls

import (
	"context"
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
	"telephone-service/internal/transport/http/mapper"
)

type GetCallsUsecase interface {
	Execute(ctx context.Context) ([]model.Call, error)
}

type CallsHandler struct {
	getCallsUsecase GetCallsUsecase
}

func NewCallsHandler(getCallsUsecase GetCallsUsecase) *CallsHandler {
	return &CallsHandler{getCallsUsecase: getCallsUsecase}
}

func (h *CallsHandler) GetCalls(ctx context.Context) (api.GetCallsRes, error) {
	calls, err := h.getCallsUsecase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("получить звонки: %w", err)
	}

	response := api.GetCallsOKApplicationJSON(mapper.CallsToResponse(calls))
	return &response, nil
}

func (*CallsHandler) CreateCall(
	_ context.Context,
	_ *api.CreateCallRequest,
) (api.CreateCallRes, error) {
	// Бизнес-сценарий создания звонка не входит в текущий рефакторинг.
	response := api.CreateCallNotImplemented(api.Error{
		ErrorMessage: "Создание звонка пока не реализовано",
	})
	return &response, nil
}
