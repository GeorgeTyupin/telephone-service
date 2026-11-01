package model

import (
	"time"

	"github.com/itlightning/dateparse"
)

type Duration time.Duration

var WrongDuration = Duration(-1)

func (d Duration) DurationToString() string {
	return time.Duration(d).String()
}

func StringToDuration(duration string) (Duration, error) {
	durationString, err := time.ParseDuration(duration)
	if err != nil {
		return WrongDuration, DurationFormatErr
	}

	return Duration(durationString), nil
}

type Time time.Time

var WrongTime = Time(time.Time{})

func (a Time) TimeToString() string {
	return time.Time(a).String()
}

func StringToTime(date string) (Time, error) {
	duration, err := dateparse.ParseAny(date)
	if err != nil {
		return WrongTime, TimeFormatError
	}

	return Time(duration), nil
}
