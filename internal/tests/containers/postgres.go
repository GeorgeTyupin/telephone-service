//go:build integration

package containers

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func StartPostgres(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(ctx, "postgres:18.6",
		postgres.WithDatabase("telephone_test"),
		postgres.WithUsername("telephone_test"),
		postgres.WithPassword("telephone_test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return container, fmt.Errorf("запуск тестовой PostgreSQL: %w", err)
	}

	return container, nil
}
