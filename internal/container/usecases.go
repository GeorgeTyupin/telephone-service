package container

import (
	"telephone-service/internal/usecase/calls/create_call"
	"telephone-service/internal/usecase/calls/get_calls"
	"telephone-service/internal/usecase/subscribers/get_subscribers"
)

func (c *Container) ResolveGetCallsUsecase() *get_calls_uc.GetCallsUsecase {
	if c.getCallsUsecase == nil {
		c.getCallsUsecase = get_calls_uc.NewGetCallsUsecase(c.ResolveCallsService())
	}

	return c.getCallsUsecase
}

func (c *Container) ResolveCreateCallUsecase() *create_call_uc.CreateCallUsecase {
	if c.createCallUsecase == nil {
		c.createCallUsecase = create_call_uc.NewCreateCallUsecase(c.ResolveCallsService())
	}

	return c.createCallUsecase
}

func (c *Container) ResolveGetSubscribersUsecase() *get_subscribers_uc.GetSubscribersUsecase {
	if c.getSubscribersUsecase == nil {
		c.getSubscribersUsecase = get_subscribers_uc.NewGetSubscribersUsecase(c.ResolveSubscribersService())
	}

	return c.getSubscribersUsecase
}
