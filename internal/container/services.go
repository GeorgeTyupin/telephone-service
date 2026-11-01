package container

import (
	callsservice "telephone-service/internal/service/calls"
	subscribersservice "telephone-service/internal/service/subscribers"
)

func (c *Container) ResolveCallsService() *callsservice.Service {
	if c.callsService == nil {
		c.callsService = callsservice.NewService(c.ResolvePostgresRepository())
	}

	return c.callsService
}

func (c *Container) ResolveSubscribersService() *subscribersservice.Service {
	if c.subscribersService == nil {
		c.subscribersService = subscribersservice.NewService(c.ResolvePostgresRepository())
	}

	return c.subscribersService
}
