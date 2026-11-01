package adapter

import (
	"fmt"
	"math"
	"telephone-service/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func idFromRecord(id int64) (int, error) {
	if int64(int(id)) != id {
		return 0, fmt.Errorf("ID %d не помещается в int", id)
	}

	return int(id), nil
}

func durationFromRecord(interval pgtype.Interval) (model.Duration, error) {
	if !interval.Valid {
		return 0, fmt.Errorf("длительность отсутствует")
	}

	if interval.Months != 0 {
		return 0, fmt.Errorf("интервал с месяцами нельзя преобразовать в длительность")
	}

	const microsecondsPerDay = int64(24 * time.Hour / time.Microsecond)
	days := int64(interval.Days)
	if days > math.MaxInt64/microsecondsPerDay || days < math.MinInt64/microsecondsPerDay {
		return 0, fmt.Errorf("длительность превышает диапазон time.Duration")
	}

	microseconds := days * microsecondsPerDay
	if interval.Microseconds > 0 && microseconds > math.MaxInt64-interval.Microseconds ||
		interval.Microseconds < 0 && microseconds < math.MinInt64-interval.Microseconds {
		return 0, fmt.Errorf("длительность превышает диапазон time.Duration")
	}

	microseconds += interval.Microseconds
	if microseconds > math.MaxInt64/int64(time.Microsecond) || microseconds < math.MinInt64/int64(time.Microsecond) {
		return 0, fmt.Errorf("длительность превышает диапазон time.Duration")
	}

	return model.Duration(microseconds * int64(time.Microsecond)), nil
}
