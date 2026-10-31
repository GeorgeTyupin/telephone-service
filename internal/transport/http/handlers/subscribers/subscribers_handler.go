package subscribers

import (
	"context"
	"telephone-service/codegen/oas"
)

type SubscribersHandler struct{}

func NewSubscribersHandler() *SubscribersHandler {
	return &SubscribersHandler{}
}

func (*SubscribersHandler) GetSubscribers(
	_ context.Context,
) (api.GetSubscribersRes, error) {
	response := api.GetSubscribersNotImplemented(api.Error{
		ErrorMessage: "Получение абонентов пока не реализовано",
	})
	return &response, nil
}
