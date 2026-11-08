package container

import (
	"telephone-service/internal/usecase/calls/connect_call"
	"telephone-service/internal/usecase/calls/disconnect_call"
	"telephone-service/internal/usecase/calls/finish_call"
	"telephone-service/internal/usecase/calls/get_call"
	"telephone-service/internal/usecase/calls/get_calls"
	"telephone-service/internal/usecase/calls/start_call"
	"telephone-service/internal/usecase/subscribers/get_subscriber"
	"telephone-service/internal/usecase/subscribers/get_subscriber_calls"
	"telephone-service/internal/usecase/subscribers/get_subscriber_stats"
	"telephone-service/internal/usecase/subscribers/get_subscribers"
)

func (c *Container) ResolveGetCallsUsecase() *get_calls_uc.GetCallsUsecase {
	if c.getCallsUsecase == nil {
		c.getCallsUsecase = get_calls_uc.NewGetCallsUsecase(c.ResolveCallsService())
	}

	return c.getCallsUsecase
}

func (c *Container) ResolveGetCallUsecase() *get_call_uc.GetCallUsecase {
	if c.getCallUsecase == nil {
		c.getCallUsecase = get_call_uc.NewGetCallUsecase(c.ResolveCallsService())
	}

	return c.getCallUsecase
}

func (c *Container) ResolveStartCallUsecase() *start_call_uc.StartCallUsecase {
	if c.startCallUsecase == nil {
		c.startCallUsecase = start_call_uc.NewStartCallUsecase(c.ResolveCallsService())
	}

	return c.startCallUsecase
}

func (c *Container) ResolveConnectCallUsecase() *connect_call_uc.ConnectCallUsecase {
	if c.connectCallUsecase == nil {
		c.connectCallUsecase = connect_call_uc.NewConnectCallUsecase(c.ResolveCallsService())
	}

	return c.connectCallUsecase
}

func (c *Container) ResolveDisconnectCallUsecase() *disconnect_call_uc.DisconnectCallUsecase {
	if c.disconnectCallUsecase == nil {
		c.disconnectCallUsecase = disconnect_call_uc.NewDisconnectCallUsecase(c.ResolveCallsService())
	}

	return c.disconnectCallUsecase
}

func (c *Container) ResolveFinishCallUsecase() *finish_call_uc.FinishCallUsecase {
	if c.finishCallUsecase == nil {
		c.finishCallUsecase = finish_call_uc.NewFinishCallUsecase(c.ResolveCallsService())
	}

	return c.finishCallUsecase
}

func (c *Container) ResolveGetSubscriberUsecase() *get_subscriber_uc.GetSubscriberUsecase {
	if c.getSubscriberUsecase == nil {
		c.getSubscriberUsecase = get_subscriber_uc.NewGetSubscriberUsecase(c.ResolveSubscribersService())
	}

	return c.getSubscriberUsecase
}

func (c *Container) ResolveGetSubscriberCallsUsecase() *get_subscriber_calls_uc.GetSubscriberCallsUsecase {
	if c.getSubscriberCallsUsecase == nil {
		c.getSubscriberCallsUsecase = get_subscriber_calls_uc.NewGetSubscriberCallsUsecase(c.ResolveCallsService())
	}

	return c.getSubscriberCallsUsecase
}

func (c *Container) ResolveGetSubscriberStatsUsecase() *get_subscriber_stats_uc.GetSubscriberStatsUsecase {
	if c.getSubscriberStatsUsecase == nil {
		c.getSubscriberStatsUsecase = get_subscriber_stats_uc.NewGetSubscriberStatsUsecase(c.ResolveSubscribersService())
	}

	return c.getSubscriberStatsUsecase
}

func (c *Container) ResolveGetSubscribersUsecase() *get_subscribers_uc.GetSubscribersUsecase {
	if c.getSubscribersUsecase == nil {
		c.getSubscribersUsecase = get_subscribers_uc.NewGetSubscribersUsecase(c.ResolveSubscribersService())
	}

	return c.getSubscribersUsecase
}
