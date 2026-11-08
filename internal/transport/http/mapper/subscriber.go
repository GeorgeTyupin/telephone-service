package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

func SubscriberFromRequest(subscriber api.StartCallParticipant) (model.Subscriber, error) {
	number, err := model.StringToPhoneNumber(subscriber.PhoneNumber)
	if err != nil {
		return model.Subscriber{}, fmt.Errorf("преобразование абонента, phone_number: %w", err)
	}

	return model.Subscriber{
		Name:        subscriber.Name,
		PhoneNumber: number,
	}, nil
}

func SubscriberToResponse(subscriber model.Subscriber) api.Subscriber {
	return api.Subscriber{
		ID:          subscriber.Id,
		Name:        subscriber.Name,
		PhoneNumber: model.PhoneNumberToString(subscriber.PhoneNumber),
	}
}

func SubscriberStatsToResponse(stats model.SubscriberStats) api.SubscriberStats {
	return api.SubscriberStats{
		CallsCount:      stats.CallsCount,
		TotalDuration:   stats.TotalDuration.DurationToString(),
		AverageDuration: stats.AverageDuration.DurationToString(),
	}
}
