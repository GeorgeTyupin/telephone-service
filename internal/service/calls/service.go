package calls

import (
	"context"
	"fmt"
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

	participants := make(map[int]struct{}, len(call.Participants))
	for _, participant := range call.Participants {
		id := participant.Id
		if _, exists := participants[id]; exists {
			return model.Call{}, fmt.Errorf("%w: участник %d указан повторно", model.ErrInvalidCall, id)
		}

		participants[id] = struct{}{}
	}

	return s.repository.CreateCall(ctx, call)
}
