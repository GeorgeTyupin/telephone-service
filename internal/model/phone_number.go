package model

import (
	"fmt"
	"regexp"
	"strconv"
)

type PhoneNumber struct {
	CountryCode uint16
	Number      uint64
}

var WrongPhoneNumber = PhoneNumber{}

func StringToPhoneNumber(phoneNumberStr string) (PhoneNumber, error) {
	re := regexp.MustCompile(`\D`)
	withoutSymbols := re.ReplaceAllString(phoneNumberStr, "")
	if len(withoutSymbols) < 11 {
		return WrongPhoneNumber, PhoneNumberFormatErr
	}

	onlyNumberStr := withoutSymbols[len(withoutSymbols)-10:]
	number, _ := strconv.ParseUint(onlyNumberStr, 10, 64)
	code, err := strconv.ParseUint(withoutSymbols[:len(withoutSymbols)-len(onlyNumberStr)], 10, 16)
	if err != nil || code > 999 {
		return WrongPhoneNumber, CountryCodeErr
	}
	return PhoneNumber{CountryCode: uint16(code), Number: number}, nil
}

func PhoneNumberToString(phoneNumber PhoneNumber) string {
	numberStr := fmt.Sprintf("%010d", phoneNumber.Number)

	part1, part2, part3, part4 := numberStr[:3], numberStr[3:6], numberStr[6:8], numberStr[8:]
	return fmt.Sprintf("+%d (%s) %s-%s-%s",
		phoneNumber.CountryCode,
		part1,
		part2,
		part3,
		part4,
	)
}
