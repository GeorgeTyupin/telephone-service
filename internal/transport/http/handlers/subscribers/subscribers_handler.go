package subscribers_handlers

import (
	"context"
	"errors"
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
	"telephone-service/internal/transport/http/mapper"
)

type SubscribersHandler struct {
	getSubscribersUsecase     GetSubscribersUsecase
	getSubscriberUsecase      GetSubscriberUsecase
	getSubscriberCallsUsecase GetSubscriberCallsUsecase
	getSubscriberStatsUsecase GetSubscriberStatsUsecase
}

func NewSubscribersHandler(
	getSubscribersUsecase GetSubscribersUsecase,
	getSubscriberUsecase GetSubscriberUsecase,
	getSubscriberCallsUsecase GetSubscriberCallsUsecase,
	getSubscriberStatsUsecase GetSubscriberStatsUsecase,
) *SubscribersHandler {
	return &SubscribersHandler{
		getSubscribersUsecase:     getSubscribersUsecase,
		getSubscriberUsecase:      getSubscriberUsecase,
		getSubscriberCallsUsecase: getSubscriberCallsUsecase,
		getSubscriberStatsUsecase: getSubscriberStatsUsecase,
	}
}

func (h *SubscribersHandler) GetSubscribers(ctx context.Context) (api.GetSubscribersRes, error) {
	subscribers, err := h.getSubscribersUsecase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить абонентов: %w", err)
	}

	response := make(api.GetSubscribersOKApplicationJSON, len(subscribers))
	for i, subscriber := range subscribers {
		response[i] = mapper.SubscriberToResponse(subscriber)
	}

	return &response, nil
}

func (h *SubscribersHandler) GetSubscriber(ctx context.Context, params api.GetSubscriberParams) (api.GetSubscriberRes, error) {
	subscriber, err := h.getSubscriberUsecase.Execute(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidSubscriber):
			return &api.GetSubscriberBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrSubscriberNotFound):
			return &api.GetSubscriberNotFound{ErrorMessage: model.ErrSubscriberNotFound.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось получить абонента: %w", err)
		}
	}

	response := mapper.SubscriberToResponse(subscriber)
	return &response, nil
}

func (h *SubscribersHandler) GetSubscriberCalls(ctx context.Context, params api.GetSubscriberCallsParams) (api.GetSubscriberCallsRes, error) {
	calls, err := h.getSubscriberCallsUsecase.Execute(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidSubscriber):
			return &api.GetSubscriberCallsBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrSubscriberNotFound):
			return &api.GetSubscriberCallsNotFound{ErrorMessage: model.ErrSubscriberNotFound.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось получить звонки абонента: %w", err)
		}
	}

	response := api.GetSubscriberCallsOKApplicationJSON(mapper.CallsToResponse(calls))
	return &response, nil
}

func (h *SubscribersHandler) GetSubscriberStats(ctx context.Context, params api.GetSubscriberStatsParams) (api.GetSubscriberStatsRes, error) {
	stats, err := h.getSubscriberStatsUsecase.Execute(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidSubscriber):
			return &api.GetSubscriberStatsBadRequest{ErrorMessage: err.Error()}, nil
		case errors.Is(err, model.ErrSubscriberNotFound):
			return &api.GetSubscriberStatsNotFound{ErrorMessage: model.ErrSubscriberNotFound.Error()}, nil
		default:
			return nil, fmt.Errorf("не удалось получить статистику абонента: %w", err)
		}
	}

	response := mapper.SubscriberStatsToResponse(stats)
	return &response, nil
}
