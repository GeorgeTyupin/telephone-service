//go:build integration

package fixtures

import (
	"context"
	"fmt"
	"slices"
	"time"

	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/adapter"
	"telephone-service/internal/repository/postgres/record"
	"telephone-service/internal/tests/builders"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type callRelations struct {
	participants   []model.Subscriber
	participations []model.CallParticipation
}

type CallOption func(*callRelations)

func WithParticipants(participants ...model.Subscriber) CallOption {
	return func(relations *callRelations) {
		relations.participants = slices.Clone(participants)
	}
}

func WithParticipations(participations ...model.CallParticipation) CallOption {
	return func(relations *callRelations) {
		relations.participations = slices.Clone(participations)
	}
}

func CreateCall(ctx context.Context, db *pgxpool.Pool, call model.Call, options ...CallOption) (model.Call, error) {
	relations := callRelations{participants: call.Participants, participations: call.Participations}
	for _, option := range options {
		option(&relations)
	}

	err := pgx.BeginTxFunc(ctx, db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `
			insert into calls (duration, created_at, status, started_at, finished_at, version)
			values ($1, $2, $3, $4, $5, $6) returning id`,
			pgtype.Interval{Microseconds: int64(time.Duration(call.Duration) / time.Microsecond), Valid: true},
			time.Time(call.CreatedAt), string(call.Status), timestamp(call.StartedAt), timestamp(call.FinishedAt), call.Version,
		).Scan(&call.Id); err != nil {
			return fmt.Errorf("сохранение звонка: %w", err)
		}

		queries := record.New(tx)
		params := adapter.SubscribersToRecord(relations.participants)
		if err := queries.CreateSubscribers(ctx, params); err != nil {
			return fmt.Errorf("сохранение абонентов: %w", err)
		}
		subscribers, err := queries.GetSubscribersByPhoneNumbers(ctx, params.PhoneNumbers)
		if err != nil {
			return fmt.Errorf("получение сохранённых абонентов: %w", err)
		}
		call.Participants = make([]model.Subscriber, len(subscribers))
		participantIDs := make([]int64, len(subscribers))
		for i, subscriber := range subscribers {
			call.Participants[i], err = adapter.SubscriberFromRecord(subscriber)
			if err != nil {
				return fmt.Errorf("преобразование абонента: %w", err)
			}
			participantIDs[i] = subscriber.ID
		}
		if err := queries.AddCallParticipants(ctx, record.AddCallParticipantsParams{
			CallID: int64(call.Id), SubscriberIDs: participantIDs,
		}); err != nil {
			return fmt.Errorf("сохранение участников звонка: %w", err)
		}

		if relations.participations == nil {
			relations.participations = make([]model.CallParticipation, len(call.Participants))
			for i, participant := range call.Participants {
				relations.participations[i] = builders.NewParticipation(call, participant.Id)
			}
		}
		call.Participations = slices.Clone(relations.participations)
		for i := range call.Participations {
			participation := &call.Participations[i]
			if err := tx.QueryRow(ctx, `
				insert into call_participations (call_id, subscriber_id, joined_at, left_at)
				values ($1, $2, $3, $4) returning id`,
				call.Id, participation.SubscriberId, time.Time(participation.JoinedAt), timestamp(participation.LeftAt),
			).Scan(&participation.Id); err != nil {
				return fmt.Errorf("сохранение периода участия: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return model.Call{}, fmt.Errorf("создание фикстуры звонка: %w", err)
	}

	return call, nil
}

func timestamp(value *model.Time) *time.Time {
	if value == nil {
		return nil
	}
	result := time.Time(*value)
	return &result
}
