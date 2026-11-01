package adapter

import (
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/record"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func CallToRecord(call model.Call) record.CreateCallParams {
	return record.CreateCallParams{
		Duration: pgtype.Interval{
			Microseconds: int64(time.Duration(call.Duration) / time.Microsecond),
			Valid:        true,
		},
		CreatedAt: pgtype.Timestamptz{Time: time.Time(call.CreatedAt), Valid: true},
	}
}

func CallFromRecord(call record.Call) (model.Call, error) {
	id, err := idFromRecord(call.ID)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать ID звонка: %w", err)
	}

	duration, err := durationFromRecord(call.Duration)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать длительность звонка %d: %w", id, err)
	}

	if !call.CreatedAt.Valid || call.CreatedAt.InfinityModifier != pgtype.Finite {
		return model.Call{}, fmt.Errorf("у звонка %d отсутствует конечное время создания", id)
	}

	return model.Call{
		Id:           id,
		Participants: make([]model.Subscriber, 0),
		Duration:     duration,
		CreatedAt:    model.Time(call.CreatedAt.Time),
	}, nil
}

func CallsFromRecords(rows []record.GetCallsRow) ([]model.Call, error) {
	calls := make([]model.Call, 0)
	for _, row := range rows {
		if len(calls) == 0 || int64(calls[len(calls)-1].Id) != row.ID {
			call, err := CallFromRecord(record.Call{
				ID:        row.ID,
				Duration:  row.Duration,
				CreatedAt: row.CreatedAt,
			})
			if err != nil {
				return nil, err
			}

			calls = append(calls, call)
		}

		if !row.SubscriberID.Valid {
			continue
		}

		participant, err := SubscriberFromRecord(record.Abonent{
			ID: row.SubscriberID.Int64, Name: row.Name, PhoneNumber: row.PhoneNumber,
		})
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать участника звонка %d: %w", row.ID, err)
		}

		calls[len(calls)-1].Participants = append(calls[len(calls)-1].Participants, participant)
	}

	return calls, nil
}
