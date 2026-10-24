package container

import "telephone-service/internal/app"

func (c *Container) BuildApp() *app.App {
	if c.app == nil {
		c.app = app.NewApp(
			c.logger,
			c.config,
			c.ResolveCallsHandler(),
		)
	}
	return c.app
}
