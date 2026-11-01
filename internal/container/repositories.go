package container

import (
	"telephone-service/internal/repository/postgres"
)

func (c *Container) ResolvePostgresRepository() *postgres.Repository {
	if c.postgresRepository == nil {
		c.postgresRepository = postgres.NewRepository(c.logger, c.dbPool)
	}

	return c.postgresRepository
}
