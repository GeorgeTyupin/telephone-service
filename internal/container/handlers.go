package container

import calls_handler "telephone-service/internal/handler/calls"

func (c *Container) ResolveCallsHandler() *calls_handler.CallsHandler {
	if c.callsHandler == nil {
		c.callsHandler = calls_handler.NewCallsHandler(
			c.logger,
			c.ResolveGetCallsUsecase(),
		)
	}
	return c.callsHandler
}
