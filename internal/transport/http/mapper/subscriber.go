package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/domain"
)

func SubscriberFromRequest(subscriber api.Subscriber) (domain.Subscriber, error) {
	number, err := domain.StringToPhoneNumber(subscriber.PhoneNumber)
	if err != nil {
		return domain.Subscriber{}, fmt.Errorf("преобразование абонента, phone_number: %w", err)
	}

	return domain.Subscriber{
		Id:          subscriber.ID,
		Name:        subscriber.Name,
		PhoneNumber: number,
	}, nil
}

func SubscriberToResponse(subscriber domain.Subscriber) api.Subscriber {
	return api.Subscriber{
		ID:          subscriber.Id,
		Name:        subscriber.Name,
		PhoneNumber: domain.PhoneNumberToString(subscriber.PhoneNumber),
	}
}
