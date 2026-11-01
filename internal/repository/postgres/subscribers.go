package postgres

import (
	"context"
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/adapter"
)

func (r *Repository) GetSubscribers(ctx context.Context) ([]model.Subscriber, error) {
	rows, err := r.queries.GetSubscribers(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить абонентов из БД: %w", err)
	}

	subscribers := make([]model.Subscriber, len(rows))
	for i, row := range rows {
		subscriber, err := adapter.SubscriberFromRecord(row)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать абонента %d из БД: %w", row.ID, err)
		}

		subscribers[i] = subscriber
	}

	return subscribers, nil
}
