package subscribers

import (
	"context"
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
