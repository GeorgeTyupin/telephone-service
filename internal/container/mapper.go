package container

import "telephone-service/internal/mapper"

func (c *Container) ResolveMapper() *mapper.Mapper {
	if c.mapper == nil {
		c.mapper = mapper.NewMapper(c.logger)
	}
	return c.mapper
}
