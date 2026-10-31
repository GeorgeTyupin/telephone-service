package container

import (
	"telephone-service/codegen/oas"
	"telephone-service/internal/transport/http/handlers/calls"
	"telephone-service/internal/transport/http/handlers/subscribers"
)

type apiHandler struct {
	*calls.CallsHandler
	*subscribers.SubscribersHandler
}

var _ api.Handler = (*apiHandler)(nil)

func (c *Container) ResolveCallsHandler() *calls.CallsHandler {
	if c.callsHandler == nil {
		c.callsHandler = calls.NewCallsHandler(c.ResolveGetCallsUsecase())
	}
	return c.callsHandler
}

func (c *Container) ResolveSubscribersHandler() *subscribers.SubscribersHandler {
	if c.subscribersHandler == nil {
		c.subscribersHandler = subscribers.NewSubscribersHandler()
	}
	return c.subscribersHandler
}
