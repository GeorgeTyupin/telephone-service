package postgres

import (
	"context"
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/adapter"
	"telephone-service/internal/repository/postgres/record"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetCalls(ctx context.Context) ([]model.Call, error) {
	rows, err := r.queries.GetCalls(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить звонки из БД: %w", err)
	}

	calls, err := adapter.CallsFromRecords(rows)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать звонки из БД: %w", err)
	}

	return calls, nil
}

func (r *Repository) CreateCall(ctx context.Context, call model.Call) (model.Call, error) {
	params := adapter.SubscribersToRecord(call.Participants)
	var createdCall model.Call

	err := pgx.BeginTxFunc(ctx, r.dbPool, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		q := r.queries.WithTx(tx)

		if err := q.CreateSubscribers(ctx, params); err != nil {
			return fmt.Errorf("не удалось создать абонентов: %w", err)
		}

		storedSubscribers, err := q.GetSubscribersByPhoneNumbers(ctx, params.PhoneNumbers)
		if err != nil {
			return fmt.Errorf("не удалось получить участников звонка: %w", err)
		}
		if len(storedSubscribers) != len(call.Participants) {
			return fmt.Errorf("не удалось получить всех участников звонка: ожидалось %d, получено %d", len(call.Participants), len(storedSubscribers))
		}

		participants := make([]model.Subscriber, len(storedSubscribers))
		participantIDs := make([]int64, len(storedSubscribers))
		for i, storedSubscriber := range storedSubscribers {
			subscriber, err := adapter.SubscriberFromRecord(storedSubscriber)
			if err != nil {
				return fmt.Errorf("не удалось преобразовать участника звонка: %w", err)
			}

			participants[i] = subscriber
			participantIDs[i] = storedSubscriber.ID
		}

		storedCall, err := q.CreateCall(ctx, adapter.CallToRecord(call))
		if err != nil {
			return fmt.Errorf("не удалось вставить звонок: %w", err)
		}

		createdCall, err = adapter.CallFromRecord(storedCall)
		if err != nil {
			return fmt.Errorf("не удалось преобразовать созданный звонок: %w", err)
		}
		createdCall.Participants = participants

		if err := q.AddCallParticipants(ctx, record.AddCallParticipantsParams{
			CallID: storedCall.ID, Column2: participantIDs,
		}); err != nil {
			return fmt.Errorf("не удалось вставить участников звонка: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось создать звонок в транзакции: %w", err)
	}

	return createdCall, nil
}
