package container

import (
	"telephone-service/codegen/oas"
	"telephone-service/internal/transport/http/handlers/calls"
	"telephone-service/internal/transport/http/handlers/subscribers"
)

type apiHandler struct {
	*calls_handlers.CallsHandler
	*subscribers_handlers.SubscribersHandler
}

var _ api.Handler = (*apiHandler)(nil)

func (c *Container) ResolveCallsHandler() *calls_handlers.CallsHandler {
	if c.callsHandler == nil {
		c.callsHandler = calls_handlers.NewCallsHandler(
			c.ResolveGetCallsUsecase(),
			c.ResolveGetCallUsecase(),
			c.ResolveStartCallUsecase(),
			c.ResolveConnectCallUsecase(),
			c.ResolveDisconnectCallUsecase(),
			c.ResolveFinishCallUsecase(),
		)
	}

	return c.callsHandler
}

func (c *Container) ResolveSubscribersHandler() *subscribers_handlers.SubscribersHandler {
	if c.subscribersHandler == nil {
		c.subscribersHandler = subscribers_handlers.NewSubscribersHandler(
			c.ResolveGetSubscribersUsecase(),
			c.ResolveGetSubscriberUsecase(),
			c.ResolveGetSubscriberCallsUsecase(),
			c.ResolveGetSubscriberStatsUsecase(),
		)
	}

	return c.subscribersHandler
}
