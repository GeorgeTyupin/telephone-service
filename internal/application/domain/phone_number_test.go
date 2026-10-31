package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToPhoneNumber(t *testing.T) {
	t.Parallel()

	//arrange
	table := []struct {
		testName       string
		phoneNumberStr string
		expectedCode   uint16
		expectedNumber uint64
		expectedErr    error
	}{
		{
			testName:       "обычный номер",
			phoneNumberStr: "+79876543210",
			expectedCode:   7,
			expectedNumber: 9876543210,
		},
		{
			testName:       "длинный код страны",
			phoneNumberStr: "+1239876543210",
			expectedCode:   123,
			expectedNumber: 9876543210,
		},
		{
			testName:       "с другими символами",
			phoneNumberStr: "+7 (987) 654 32-10",
			expectedCode:   7,
			expectedNumber: 9876543210,
		},
		{
			testName:       "с текстом",
			phoneNumberStr: "+79876X3210",
			expectedErr:    PhoneNumberFormatErr,
		},
		{
			testName:       "cлишком длинный код страны",
			phoneNumberStr: "+10009876543210",
			expectedErr:    CountryCodeErr,
		},
		{
			testName:       "телефон без +",
			phoneNumberStr: "89876543210",
			expectedCode:   8,
			expectedNumber: 9876543210,
		},
	}

	for _, tt := range table {
		t.Run(tt.testName, func(t *testing.T) {
			//act
			phoneNumber, err := StringToPhoneNumber(tt.phoneNumberStr)

			//assert
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr, "неверный формат ошибки")
			}
			assert.Equal(t, tt.expectedCode, phoneNumber.CountryCode, "неверный код страны")
			assert.Equal(t, tt.expectedNumber, phoneNumber.Number, "неверный номер телефона")
		})
	}
}

func TestPhoneNumberToString(t *testing.T) {
	t.Parallel()

	//arrange
	phoneNumber := PhoneNumber{
		CountryCode: 7,
		Number:      9876543210,
	}

	//act
	phoneNumberStr := PhoneNumberToString(phoneNumber)

	//assert
	assert.Equal(t, "+7 (987) 654-32-10", phoneNumberStr)
}
