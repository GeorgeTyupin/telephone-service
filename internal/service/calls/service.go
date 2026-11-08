package calls

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"telephone-service/internal/model"
	"time"
)

type Service struct {
	repository CallsRepository
}

func NewService(repository CallsRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetCalls(ctx context.Context, filter model.CallsFilter) ([]model.Call, error) {
	if filter.Status != "" && filter.Status != model.CallStatusActive && filter.Status != model.CallStatusFinished {
		return nil, fmt.Errorf("%w: неизвестное состояние звонка", model.ErrInvalidCall)
	}

	return s.repository.GetCalls(ctx, filter)
}

func (s *Service) GetCall(ctx context.Context, id int) (model.Call, error) {
	if id <= 0 {
		return model.Call{}, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidCall)
	}

	return s.repository.GetCall(ctx, id)
}

func (s *Service) GetSubscriberCalls(ctx context.Context, id int) ([]model.Call, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidSubscriber)
	}

	if _, err := s.repository.GetSubscriber(ctx, id); err != nil {
		return nil, fmt.Errorf("не удалось получить абонента: %w", err)
	}

	return s.repository.GetSubscriberCalls(ctx, id)
}

func (s *Service) StartCall(ctx context.Context, call model.Call) (model.Call, error) {
	if len(call.Participants) == 0 {
		return model.Call{}, fmt.Errorf("%w: нужен хотя бы один участник", model.ErrInvalidCall)
	}

	participants := make(map[model.PhoneNumber]struct{}, len(call.Participants))
	for _, participant := range call.Participants {
		if strings.TrimSpace(participant.Name) == "" {
			return model.Call{}, fmt.Errorf("%w: имя участника не может быть пустым", model.ErrInvalidCall)
		}

		if strings.ContainsRune(participant.Name, '\x00') {
			return model.Call{}, fmt.Errorf("%w: имя участника не может содержать нулевой символ", model.ErrInvalidCall)
		}

		if _, exists := participants[participant.PhoneNumber]; exists {
			return model.Call{}, fmt.Errorf("%w: участник с телефоном %s указан повторно", model.ErrInvalidCall, model.PhoneNumberToString(participant.PhoneNumber))
		}

		participants[participant.PhoneNumber] = struct{}{}
	}

	return s.repository.StartCall(ctx, model.Call{Participants: call.Participants, Status: model.CallStatusActive})
}

func (s *Service) ConnectCall(ctx context.Context, id, subscriberId int) (model.Call, error) {
	if subscriberId <= 0 {
		return model.Call{}, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidSubscriber)
	}

	return s.updateCall(ctx, id, func(call model.Call) (*model.CallUpdate, error) {
		if call.Status != model.CallStatusActive {
			return nil, fmt.Errorf("%w: звонок завершён", model.ErrCallConflict)
		}

		for _, participation := range call.Participations {
			if participation.SubscriberId == subscriberId && participation.LeftAt == nil {
				return nil, fmt.Errorf("%w: участник уже подключён", model.ErrCallConflict)
			}
		}

		subscriber, err := s.repository.GetSubscriber(ctx, subscriberId)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить участника: %w", err)
		}

		update := newCallUpdate(call)
		update.ConnectedParticipant = &subscriber
		return update, nil
	})
}

func (s *Service) DisconnectCall(ctx context.Context, id, subscriberId int) (model.Call, error) {
	if subscriberId <= 0 {
		return model.Call{}, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidSubscriber)
	}

	return s.updateCall(ctx, id, func(call model.Call) (*model.CallUpdate, error) {
		if call.Status != model.CallStatusActive {
			return nil, fmt.Errorf("%w: звонок завершён", model.ErrCallConflict)
		}

		if _, err := s.repository.GetSubscriber(ctx, subscriberId); err != nil {
			return nil, fmt.Errorf("не удалось получить участника: %w", err)
		}

		update := newCallUpdate(call)
		connected := 0
		for _, participation := range call.Participations {
			if participation.LeftAt != nil {
				continue
			}

			connected++
			if participation.SubscriberId == subscriberId {
				update.DisconnectedParticipationIds = append(update.DisconnectedParticipationIds, participation.Id)
			}
		}

		if len(update.DisconnectedParticipationIds) == 0 {
			return nil, fmt.Errorf("%w: участник не подключён", model.ErrCallConflict)
		}

		if connected == len(update.DisconnectedParticipationIds) {
			finishCallUpdate(update)
		}

		return update, nil
	})
}

func (s *Service) FinishCall(ctx context.Context, id int) (model.Call, error) {
	return s.updateCall(ctx, id, func(call model.Call) (*model.CallUpdate, error) {
		if call.Status == model.CallStatusFinished {
			return nil, nil
		}

		update := newCallUpdate(call)
		for _, participation := range call.Participations {
			if participation.LeftAt == nil {
				update.DisconnectedParticipationIds = append(update.DisconnectedParticipationIds, participation.Id)
			}
		}

		finishCallUpdate(update)
		return update, nil
	})
}

func (s *Service) updateCall(ctx context.Context, id int, change func(model.Call) (*model.CallUpdate, error)) (model.Call, error) {
	var lastErr error
	for range 3 {
		call, err := s.GetCall(ctx, id)
		if err != nil {
			return model.Call{}, err
		}

		update, err := change(call)
		if err != nil {
			return model.Call{}, err
		}
		if update == nil {
			return call, nil
		}

		updatedCall, err := s.repository.UpdateCall(ctx, *update)
		if errors.Is(err, model.ErrCallChanged) {
			lastErr = err
			continue
		}

		return updatedCall, err
	}

	return model.Call{}, fmt.Errorf("%w: %w", model.ErrCallConflict, lastErr)
}

func newCallUpdate(call model.Call) *model.CallUpdate {
	changedAt := call.ObservedAt
	if call.StartedAt != nil && time.Time(changedAt).Before(time.Time(*call.StartedAt)) {
		changedAt = *call.StartedAt
	}
	for _, participation := range call.Participations {
		if time.Time(changedAt).Before(time.Time(participation.JoinedAt)) {
			changedAt = participation.JoinedAt
		}
		if participation.LeftAt != nil && time.Time(changedAt).Before(time.Time(*participation.LeftAt)) {
			changedAt = *participation.LeftAt
		}
	}

	return &model.CallUpdate{Call: call, ChangedAt: changedAt}
}

func finishCallUpdate(update *model.CallUpdate) {
	update.Call.Status = model.CallStatusFinished
	update.Call.FinishedAt = &update.ChangedAt
	update.Call.Duration = model.Duration(time.Time(update.ChangedAt).Sub(time.Time(*update.Call.StartedAt)))
}
