package calls

import (
	"context"
	"fmt"
	"strings"
	"telephone-service/internal/model"
)

type Service struct {
	repository CallsRepository
}

func NewService(repository CallsRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetCalls(ctx context.Context) ([]model.Call, error) {
	return s.repository.GetCalls(ctx)
}

func (s *Service) CreateCall(ctx context.Context, call model.Call) (model.Call, error) {
	if call.Duration < 0 {
		return model.Call{}, fmt.Errorf("%w: длительность не может быть отрицательной", model.ErrInvalidCall)
	}

	if len(call.Participants) == 0 {
		return model.Call{}, fmt.Errorf("%w: нужен хотя бы один участник", model.ErrInvalidCall)
	}

	participants := make(map[model.PhoneNumber]struct{}, len(call.Participants))
	for _, participant := range call.Participants {
		if strings.TrimSpace(participant.Name) == "" {
			return model.Call{}, fmt.Errorf("%w: имя участника не может быть пустым", model.ErrInvalidCall)
		}

		if _, exists := participants[participant.PhoneNumber]; exists {
			return model.Call{}, fmt.Errorf("%w: участник с телефоном %s указан повторно", model.ErrInvalidCall, model.PhoneNumberToString(participant.PhoneNumber))
		}

		participants[participant.PhoneNumber] = struct{}{}
	}

	return s.repository.CreateCall(ctx, call)
}
