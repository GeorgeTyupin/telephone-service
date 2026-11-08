package postgres

import (
	"context"
	"errors"
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/adapter"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetSubscriber(ctx context.Context, id int) (model.Subscriber, error) {
	row, err := r.queries.GetSubscriber(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscriber{}, fmt.Errorf("%w: %w", model.ErrSubscriberNotFound, err)
		}

		return model.Subscriber{}, fmt.Errorf("не удалось получить абонента из БД: %w", err)
	}

	subscriber, err := adapter.SubscriberFromRecord(row)
	if err != nil {
		return model.Subscriber{}, fmt.Errorf("не удалось преобразовать абонента из БД: %w", err)
	}

	return subscriber, nil
}

func (r *Repository) GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error) {
	row, err := r.queries.GetSubscriberStats(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.SubscriberStats{}, fmt.Errorf("%w: %w", model.ErrSubscriberNotFound, err)
		}

		return model.SubscriberStats{}, fmt.Errorf("не удалось получить статистику абонента из БД: %w", err)
	}

	stats, err := adapter.SubscriberStatsFromRecord(row)
	if err != nil {
		return model.SubscriberStats{}, fmt.Errorf("не удалось преобразовать статистику абонента: %w", err)
	}

	return stats, nil
}
