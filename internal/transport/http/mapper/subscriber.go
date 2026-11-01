package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

func SubscriberFromRequest(subscriber api.Subscriber) (model.Subscriber, error) {
	number, err := model.StringToPhoneNumber(subscriber.PhoneNumber)
	if err != nil {
		return model.Subscriber{}, fmt.Errorf("преобразование абонента, phone_number: %w", err)
	}

	return model.Subscriber{
		Id:          subscriber.ID,
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
