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
	getCallsUsecase       GetCallsUsecase
	getCallUsecase        GetCallUsecase
	startCallUsecase      StartCallUsecase
	connectCallUsecase    ConnectCallUsecase
	disconnectCallUsecase DisconnectCallUsecase
	finishCallUsecase     FinishCallUsecase
}

func NewCallsHandler(
	getCallsUsecase GetCallsUsecase,
	getCallUsecase GetCallUsecase,
	startCallUsecase StartCallUsecase,
	connectCallUsecase ConnectCallUsecase,
	disconnectCallUsecase DisconnectCallUsecase,
	finishCallUsecase FinishCallUsecase,
) *CallsHandler {
	return &CallsHandler{
		getCallsUsecase:       getCallsUsecase,
		getCallUsecase:        getCallUsecase,
		startCallUsecase:      startCallUsecase,
		connectCallUsecase:    connectCallUsecase,
		disconnectCallUsecase: disconnectCallUsecase,
		finishCallUsecase:     finishCallUsecase,
	}
}

func (h *CallsHandler) GetCalls(ctx context.Context, params api.GetCallsParams) (api.GetCallsRes, error) {
	calls, err := h.getCallsUsecase.Execute(ctx, mapper.CallsFilterFromParams(params))
	if err != nil {
		if errors.Is(err, model.ErrInvalidCall) {
			return &api.GetCallsBadRequest{ErrorMessage: err.Error()}, nil
		}

		return nil, fmt.Errorf("не удалось получить звонки: %w", err)
	}

	response := api.GetCallsOKApplicationJSON(mapper.CallsToResponse(calls))
	return &response, nil
}

func (h *CallsHandler) GetCall(ctx context.Context, params api.GetCallParams) (api.GetCallRes, error) {
	call, err := h.getCallUsecase.Execute(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCall):
			return &api.GetCallBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrCallNotFound):
			return &api.GetCallNotFound{ErrorMessage: model.ErrCallNotFound.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось получить звонок: %w", err)
		}
	}

	response := mapper.CallToResponse(call)
	return &response, nil
}

func (h *CallsHandler) StartCall(ctx context.Context, request *api.StartCallRequest) (api.StartCallRes, error) {
	call, err := mapper.CallFromRequest(*request)
	if err != nil {
		return &api.StartCallBadRequest{ErrorMessage: err.Error()}, nil
	}

	startedCall, err := h.startCallUsecase.Execute(ctx, call)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCall) {
			return &api.StartCallBadRequest{ErrorMessage: err.Error()}, nil
		}

		return nil, fmt.Errorf("не удалось начать звонок: %w", err)
	}

	response := mapper.CallToResponse(startedCall)
	return &response, nil
}

func (h *CallsHandler) ConnectCall(ctx context.Context, request *api.CallParticipantRequest, params api.ConnectCallParams) (api.ConnectCallRes, error) {
	call, err := h.connectCallUsecase.Execute(ctx, params.ID, request.SubscriberID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCall), errors.Is(err, model.ErrInvalidSubscriber):
			return &api.ConnectCallBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrCallNotFound), errors.Is(err, model.ErrSubscriberNotFound):
			return &api.ConnectCallNotFound{ErrorMessage: "звонок или абонент не найден"}, nil
		case errors.Is(err, model.ErrCallConflict):
			return &api.ConnectCallConflict{ErrorMessage: model.ErrCallConflict.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось подключить участника: %w", err)
		}
	}

	response := mapper.CallToResponse(call)
	return &response, nil
}

func (h *CallsHandler) DisconnectCall(ctx context.Context, request *api.CallParticipantRequest, params api.DisconnectCallParams) (api.DisconnectCallRes, error) {
	call, err := h.disconnectCallUsecase.Execute(ctx, params.ID, request.SubscriberID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCall), errors.Is(err, model.ErrInvalidSubscriber):
			return &api.DisconnectCallBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrCallNotFound), errors.Is(err, model.ErrSubscriberNotFound):
			return &api.DisconnectCallNotFound{ErrorMessage: "звонок или абонент не найден"}, nil
		case errors.Is(err, model.ErrCallConflict):
			return &api.DisconnectCallConflict{ErrorMessage: model.ErrCallConflict.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось отключить участника: %w", err)
		}
	}

	response := mapper.CallToResponse(call)
	return &response, nil
}

func (h *CallsHandler) FinishCall(ctx context.Context, params api.FinishCallParams) (api.FinishCallRes, error) {
	call, err := h.finishCallUsecase.Execute(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCall):
			return &api.FinishCallBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrCallNotFound):
			return &api.FinishCallNotFound{ErrorMessage: model.ErrCallNotFound.Error()}, nil
		case errors.Is(err, model.ErrCallConflict):
			return &api.FinishCallConflict{ErrorMessage: model.ErrCallConflict.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось завершить звонок: %w", err)
		}
	}

	response := mapper.CallToResponse(call)
	return &response, nil
}
