package subscribers

import (
	"context"
	"fmt"
	"telephone-service/internal/model"
)

type Service struct {
	repository SubscribersRepository
}

func NewService(repository SubscribersRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetSubscribers(ctx context.Context) ([]model.Subscriber, error) {
	return s.repository.GetSubscribers(ctx)
}

func (s *Service) GetSubscriber(ctx context.Context, id int) (model.Subscriber, error) {
	if id <= 0 {
		return model.Subscriber{}, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidSubscriber)
	}

	return s.repository.GetSubscriber(ctx, id)
}

func (s *Service) GetSubscriberStats(ctx context.Context, id int) (model.SubscriberStats, error) {
	if id <= 0 {
		return model.SubscriberStats{}, fmt.Errorf("%w: ID должен быть положительным", model.ErrInvalidSubscriber)
	}

	return s.repository.GetSubscriberStats(ctx, id)
}
