package model

import "errors"

var (
	TimeFormatError   = errors.New("неверный формат времени")
	DurationFormatErr = errors.New("неверный формат временного периода")

	PhoneNumberFormatErr = errors.New("неверный формат номера телефона")
	CountryCodeErr       = errors.New("неверный код страны")
)
