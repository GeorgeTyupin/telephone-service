package handlers

import (
	api "telephone-service/codegen/oas"
	"telephone-service/internal/transport/http/handlers/calls"
	"telephone-service/internal/transport/http/handlers/subscribers"
)

var _ api.Handler = (*handler)(nil)

type handler struct {
	*calls.CallsHandler
	*subscribers.SubscribersHandler
}

func NewHandler(
	callsHandler *calls.CallsHandler,
	subscribersHandler *subscribers.SubscribersHandler,
) api.Handler {
	return &handler{
		CallsHandler:       callsHandler,
		SubscribersHandler: subscribersHandler,
	}
}
