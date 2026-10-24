package container

import (
	"telephone-service/internal/usecase/calls"
)

func (c *Container) ResolveGetCallsUsecase() *calls.GetCallsUsecase {
	if c.getCallsUsecase == nil {
		c.getCallsUsecase = calls.NewGetCallsUsecase(
			c.logger,
		)
	}
	return c.getCallsUsecase
}
