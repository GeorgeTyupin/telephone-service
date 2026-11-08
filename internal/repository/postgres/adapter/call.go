package adapter

import (
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/record"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func CallUpdateToRecord(update model.CallUpdate) record.UpdateCallParams {
	duration := update.Call.Duration
	if update.Call.Status == model.CallStatusActive {
		duration = 0
	}

	return record.UpdateCallParams{
		CallID:     int64(update.Call.Id),
		Version:    update.Call.Version,
		Status:     string(update.Call.Status),
		Duration:   pgtype.Interval{Microseconds: int64(time.Duration(duration) / time.Microsecond), Valid: true},
		FinishedAt: timeToRecord(update.Call.FinishedAt),
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

	createdAt, err := timeFromRecord(call.CreatedAt)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать время создания звонка %d: %w", id, err)
	}
	if createdAt == nil {
		return model.Call{}, fmt.Errorf("у звонка %d отсутствует время создания", id)
	}

	startedAt, err := timeFromRecord(call.StartedAt)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать время начала звонка %d: %w", id, err)
	}

	finishedAt, err := timeFromRecord(call.FinishedAt)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать время окончания звонка %d: %w", id, err)
	}

	status := model.CallStatus(call.Status)
	if status != model.CallStatusActive && status != model.CallStatusFinished {
		return model.Call{}, fmt.Errorf("у звонка %d неизвестное состояние %q", id, call.Status)
	}
	if status == model.CallStatusActive && startedAt == nil {
		return model.Call{}, fmt.Errorf("у активного звонка %d отсутствует время начала", id)
	}

	return model.Call{
		Id:             id,
		Participants:   make([]model.Subscriber, 0),
		Participations: make([]model.CallParticipation, 0),
		Duration:       duration,
		CreatedAt:      *createdAt,
		Status:         status,
		StartedAt:      startedAt,
		FinishedAt:     finishedAt,
		Version:        call.Version,
	}, nil
}

func CallsFromRecords(rows []record.GetCallsRow) ([]model.Call, error) {
	calls := make([]model.Call, 0)
	for _, row := range rows {
		if len(calls) == 0 || int64(calls[len(calls)-1].Id) != row.ID {
			call, err := CallFromRecord(record.Call{
				ID:         row.ID,
				Duration:   row.Elapsed,
				CreatedAt:  row.CreatedAt,
				Status:     row.Status,
				StartedAt:  row.StartedAt,
				FinishedAt: row.FinishedAt,
				Version:    row.Version,
			})
			if err != nil {
				return nil, err
			}

			observedAt, err := timeFromRecord(row.ObservedAt)
			if err != nil {
				return nil, fmt.Errorf("не удалось преобразовать время чтения звонка: %w", err)
			}
			if observedAt == nil {
				return nil, fmt.Errorf("отсутствует время чтения звонка %d", call.Id)
			}

			call.ObservedAt = *observedAt
			calls = append(calls, call)
		}

		if !row.SubscriberID.Valid {
			continue
		}

		call := &calls[len(calls)-1]
		if len(call.Participants) == 0 || int64(call.Participants[len(call.Participants)-1].Id) != row.SubscriberID.Int64 {
			participant, err := SubscriberFromRecord(record.Abonent{
				ID: row.SubscriberID.Int64, Name: row.Name, PhoneNumber: row.PhoneNumber,
			})
			if err != nil {
				return nil, fmt.Errorf("не удалось преобразовать участника звонка %d: %w", call.Id, err)
			}

			call.Participants = append(call.Participants, participant)
		}

		if !row.ParticipationID.Valid {
			continue
		}

		participation, err := CallParticipationFromRecord(record.CallParticipation{
			ID:           row.ParticipationID.Int64,
			SubscriberID: row.SubscriberID.Int64,
			JoinedAt:     row.JoinedAt,
			LeftAt:       row.LeftAt,
		})
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать период участия в звонке %d: %w", call.Id, err)
		}

		call.Participations = append(call.Participations, participation)
	}

	return calls, nil
}

func CallParticipationFromRecord(participation record.CallParticipation) (model.CallParticipation, error) {
	id, err := idFromRecord(participation.ID)
	if err != nil {
		return model.CallParticipation{}, fmt.Errorf("не удалось преобразовать ID периода участия: %w", err)
	}

	subscriberId, err := idFromRecord(participation.SubscriberID)
	if err != nil {
		return model.CallParticipation{}, fmt.Errorf("не удалось преобразовать ID участника: %w", err)
	}

	joinedAt, err := timeFromRecord(participation.JoinedAt)
	if err != nil {
		return model.CallParticipation{}, fmt.Errorf("не удалось преобразовать время подключения: %w", err)
	}
	if joinedAt == nil {
		return model.CallParticipation{}, fmt.Errorf("отсутствует время подключения участника %d", subscriberId)
	}

	leftAt, err := timeFromRecord(participation.LeftAt)
	if err != nil {
		return model.CallParticipation{}, fmt.Errorf("не удалось преобразовать время отключения: %w", err)
	}

	return model.CallParticipation{Id: id, SubscriberId: subscriberId, JoinedAt: *joinedAt, LeftAt: leftAt}, nil
}

func timeFromRecord(value pgtype.Timestamptz) (*model.Time, error) {
	if !value.Valid {
		return nil, nil
	}
	if value.InfinityModifier != pgtype.Finite {
		return nil, fmt.Errorf("время должно быть конечным")
	}

	result := model.Time(value.Time.UTC())
	return &result, nil
}

func timeToRecord(value *model.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{Time: time.Time(*value), Valid: true}
}

func CallTimeToRecord(value model.Time) pgtype.Timestamptz {
	return timeToRecord(&value)
}
