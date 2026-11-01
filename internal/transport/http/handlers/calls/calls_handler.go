package calls_handlers

import (
	"context"
	"errors"
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
	"telephone-service/internal/transport/http/mapper"
)

type CallsHandler struct {
	getCallsUsecase   GetCallsUsecase
	createCallUsecase CreateCallUsecase
}

func NewCallsHandler(getCallsUsecase GetCallsUsecase, createCallUsecase CreateCallUsecase) *CallsHandler {
	return &CallsHandler{
		getCallsUsecase:   getCallsUsecase,
		createCallUsecase: createCallUsecase,
	}
}

func (h *CallsHandler) GetCalls(ctx context.Context) (api.GetCallsRes, error) {
	calls, err := h.getCallsUsecase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить звонки: %w", err)
	}

	response := api.GetCallsOKApplicationJSON(mapper.CallsToResponse(calls))
	return &response, nil
}

func (h *CallsHandler) CreateCall(ctx context.Context, request *api.CreateCallRequest) (api.CreateCallRes, error) {
	call, err := mapper.CallFromRequest(*request)
	if err != nil {
		return &api.CreateCallBadRequest{ErrorMessage: err.Error()}, nil
	}

	createdCall, err := h.createCallUsecase.Execute(ctx, call)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCall) || errors.Is(err, model.ErrSubscriberNotFound) {
			return &api.CreateCallBadRequest{ErrorMessage: err.Error()}, nil
		}

		return nil, fmt.Errorf("не удалось создать звонок: %w", err)
	}

	response := mapper.CallToResponse(createdCall)
	return &response, nil
}
