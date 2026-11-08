package adapter

import (
	"fmt"
	"telephone-service/internal/model"
	"telephone-service/internal/repository/postgres/record"
)

func SubscribersToRecord(subscribers []model.Subscriber) record.CreateSubscribersParams {
	params := record.CreateSubscribersParams{
		Names:        make([]string, len(subscribers)),
		PhoneNumbers: make([]string, len(subscribers)),
	}
	for i, subscriber := range subscribers {
		params.Names[i] = subscriber.Name
		params.PhoneNumbers[i] = model.PhoneNumberToString(subscriber.PhoneNumber)
	}

	return params
}

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

func SubscriberStatsFromRecord(stats record.GetSubscriberStatsRow) (model.SubscriberStats, error) {
	callsCount, err := idFromRecord(stats.CallsCount)
	if err != nil {
		return model.SubscriberStats{}, fmt.Errorf("не удалось преобразовать количество звонков: %w", err)
	}

	totalDuration, err := durationFromRecord(stats.TotalDuration)
	if err != nil {
		return model.SubscriberStats{}, fmt.Errorf("не удалось преобразовать общую длительность: %w", err)
	}

	averageDuration, err := durationFromRecord(stats.AverageDuration)
	if err != nil {
		return model.SubscriberStats{}, fmt.Errorf("не удалось преобразовать среднюю длительность: %w", err)
	}

	return model.SubscriberStats{
		CallsCount: callsCount, TotalDuration: totalDuration, AverageDuration: averageDuration,
	}, nil
}
