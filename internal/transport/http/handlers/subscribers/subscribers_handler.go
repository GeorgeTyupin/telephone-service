package subscribers_handlers

import (
	"context"
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/transport/http/mapper"
)

type SubscribersHandler struct {
	getSubscribersUsecase GetSubscribersUsecase
}

func NewSubscribersHandler(getSubscribersUsecase GetSubscribersUsecase) *SubscribersHandler {
	return &SubscribersHandler{getSubscribersUsecase: getSubscribersUsecase}
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
