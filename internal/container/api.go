package container

import (
	httpserver "telephone-service/internal/transport/http"
	"telephone-service/internal/transport/http/handlers"
)

func (c *Container) BuildApi() (*httpserver.Server, error) {
	handler := handlers.NewHandler(
		c.ResolveCallsHandler(),
		c.ResolveSubscribersHandler(),
	)

	return httpserver.NewServer(c.logger, c.config.Server, handler)
}
