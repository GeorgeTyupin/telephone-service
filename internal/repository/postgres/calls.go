package postgres

import (
	"context"
	"errors"
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/adapter"
	"telephone-service/internal/repository/postgres/record"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	participantIDs := make([]int64, len(call.Participants))
	for i, participant := range call.Participants {
		participantIDs[i] = int64(participant.Id)
	}

	params := adapter.CallToRecord(call)
	var storedCall record.Call
	var participants []record.Abonent

	err := pgx.BeginFunc(ctx, r.dbPool, func(tx pgx.Tx) error {
		q := r.queries.WithTx(tx)

		var err error
		storedCall, err = q.CreateCall(ctx, params)
		if err != nil {
			return fmt.Errorf("не удалось вставить звонок: %w", err)
		}

		if err = q.AddCallParticipants(ctx, record.AddCallParticipantsParams{
			CallID: storedCall.ID, Column2: participantIDs,
		}); err != nil {
			return fmt.Errorf("не удалось вставить участников звонка: %w", err)
		}

		participants, err = q.GetCallParticipants(ctx, storedCall.ID)
		if err != nil {
			return fmt.Errorf("не удалось получить участников созданного звонка: %w", err)
		}

		return nil
	})

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch {
			case pgErr.Code == "23503" && pgErr.ConstraintName == "call_participants_subscriber_id_fkey":
				return model.Call{}, fmt.Errorf("не удалось создать звонок: %w", model.ErrSubscriberNotFound)
			case pgErr.Code == "23505" && pgErr.ConstraintName == "call_participants_pkey":
				return model.Call{}, fmt.Errorf("не удалось создать звонок: %w", model.ErrInvalidCall)
			}
		}
		return model.Call{}, fmt.Errorf("не удалось создать звонок в транзакции: %w", err)
	}

	createdCall, err := adapter.CallFromRecord(storedCall)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать созданный звонок: %w", err)
	}

	for _, participant := range participants {
		subscriber, err := adapter.SubscriberFromRecord(participant)
		if err != nil {
			return model.Call{}, fmt.Errorf("не удалось преобразовать участника созданного звонка: %w", err)
		}

		createdCall.Participants = append(createdCall.Participants, subscriber)
	}

	return createdCall, nil
}
