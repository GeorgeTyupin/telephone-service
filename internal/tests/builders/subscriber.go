package builders

import "telephone-service/internal/model"

type SubscriberOption func(*model.Subscriber)

func NewSubscriber(options ...SubscriberOption) model.Subscriber {
	subscriber := model.Subscriber{
		Name:        "Иван",
		PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567},
	}
	for _, option := range options {
		option(&subscriber)
	}
	return subscriber
}

func WithName(name string) SubscriberOption {
	return func(subscriber *model.Subscriber) {
		subscriber.Name = name
	}
}

func WithPhoneNumber(number model.PhoneNumber) SubscriberOption {
	return func(subscriber *model.Subscriber) {
		subscriber.PhoneNumber = number
	}
}
