package adapter

import (
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/record"
)

func SubscriberFromRecord(subscriber record.Abonent) (model.Subscriber, error) {
	id, err := idFromRecord(subscriber.ID)
	if err != nil {
		return model.Subscriber{}, fmt.Errorf("не удалось преобразовать ID абонента: %w", err)
	}

	if !subscriber.Name.Valid || !subscriber.PhoneNumber.Valid {
		return model.Subscriber{}, fmt.Errorf("у абонента %d отсутствует имя или номер телефона", id)
	}

	phoneNumber, err := model.StringToPhoneNumber(subscriber.PhoneNumber.String)
	if err != nil {
		return model.Subscriber{}, fmt.Errorf("не удалось преобразовать номер телефона абонента %d: %w", id, err)
	}

	return model.Subscriber{Id: id, Name: subscriber.Name.String, PhoneNumber: phoneNumber}, nil
}
