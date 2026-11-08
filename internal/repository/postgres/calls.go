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

func (r *Repository) GetCalls(ctx context.Context, filter model.CallsFilter) ([]model.Call, error) {
	return getCalls(ctx, r.queries, record.GetCallsParams{Status: string(filter.Status)})
}

func (r *Repository) GetCall(ctx context.Context, id int) (model.Call, error) {
	return getCall(ctx, r.queries, id)
}

func (r *Repository) GetSubscriberCalls(ctx context.Context, id int) ([]model.Call, error) {
	return getCalls(ctx, r.queries, record.GetCallsParams{SubscriberID: int64(id)})
}

func (r *Repository) StartCall(ctx context.Context, call model.Call) (model.Call, error) {
	params := adapter.SubscribersToRecord(call.Participants)
	var startedCall model.Call

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

		participantIDs := make([]int64, len(storedSubscribers))
		for i, subscriber := range storedSubscribers {
			participantIDs[i] = subscriber.ID
		}

		storedCall, err := q.StartCall(ctx)
		if err != nil {
			return fmt.Errorf("не удалось создать звонок: %w", err)
		}

		if err := q.AddCallParticipants(ctx, record.AddCallParticipantsParams{
			CallID: storedCall.ID, SubscriberIDs: participantIDs,
		}); err != nil {
			return fmt.Errorf("не удалось сохранить участников звонка: %w", err)
		}

		if err := q.StartCallParticipations(ctx, record.StartCallParticipationsParams{
			CallID: storedCall.ID, SubscriberIDs: participantIDs, JoinedAt: storedCall.StartedAt,
		}); err != nil {
			return fmt.Errorf("не удалось сохранить подключения участников: %w", err)
		}

		createdCall, err := adapter.CallFromRecord(storedCall)
		if err != nil {
			return fmt.Errorf("не удалось преобразовать созданный звонок: %w", err)
		}

		startedCall, err = getCall(ctx, q, createdCall.Id)
		return err
	})
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось начать звонок: %w", err)
	}

	return startedCall, nil
}

func (r *Repository) UpdateCall(ctx context.Context, update model.CallUpdate) (model.Call, error) {
	var updatedCall model.Call
	err := pgx.BeginTxFunc(ctx, r.dbPool, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		q := r.queries.WithTx(tx)
		if _, err := q.UpdateCall(ctx, adapter.CallUpdateToRecord(update)); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: %w", model.ErrCallChanged, err)
			}

			return fmt.Errorf("не удалось обновить звонок: %w", err)
		}

		if update.ConnectedParticipant != nil {
			participantIDs := []int64{int64(update.ConnectedParticipant.Id)}
			if err := q.AddCallParticipants(ctx, record.AddCallParticipantsParams{
				CallID: int64(update.Call.Id), SubscriberIDs: participantIDs,
			}); err != nil {
				return fmt.Errorf("не удалось сохранить участника: %w", err)
			}

			if err := q.StartCallParticipations(ctx, record.StartCallParticipationsParams{
				CallID: int64(update.Call.Id), SubscriberIDs: participantIDs,
				JoinedAt: adapter.CallTimeToRecord(update.ChangedAt),
			}); err != nil {
				return fmt.Errorf("не удалось сохранить подключение участника: %w", err)
			}
		}

		if len(update.DisconnectedParticipationIds) != 0 {
			participationIDs := make([]int64, len(update.DisconnectedParticipationIds))
			for i, id := range update.DisconnectedParticipationIds {
				participationIDs[i] = int64(id)
			}

			if err := q.CloseCallParticipations(ctx, record.CloseCallParticipationsParams{
				CallID: int64(update.Call.Id), ParticipationIDs: participationIDs,
				LeftAt: adapter.CallTimeToRecord(update.ChangedAt),
			}); err != nil {
				return fmt.Errorf("не удалось сохранить отключение участников: %w", err)
			}
		}

		var err error
		updatedCall, err = getCall(ctx, q, update.Call.Id)
		return err
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23503" && pgErr.ConstraintName == "call_participants_subscriber_id_fkey":
				return model.Call{}, fmt.Errorf("%w: %w", model.ErrSubscriberNotFound, err)
			case pgErr.Code == "23505" && pgErr.ConstraintName == "call_participations_connected_key":
				return model.Call{}, fmt.Errorf("%w: %w", model.ErrCallConflict, err)
			}
		}

		return model.Call{}, fmt.Errorf("не удалось сохранить изменения звонка: %w", err)
	}

	return updatedCall, nil
}

func getCalls(ctx context.Context, q *record.Queries, params record.GetCallsParams) ([]model.Call, error) {
	rows, err := q.GetCalls(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить звонки из БД: %w", err)
	}

	calls, err := adapter.CallsFromRecords(rows)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать звонки из БД: %w", err)
	}

	return calls, nil
}

func getCall(ctx context.Context, q *record.Queries, id int) (model.Call, error) {
	calls, err := getCalls(ctx, q, record.GetCallsParams{CallID: int64(id)})
	if err != nil {
		return model.Call{}, err
	}
	if len(calls) == 0 {
		return model.Call{}, fmt.Errorf("%w: %d", model.ErrCallNotFound, id)
	}

	return calls[0], nil
}
