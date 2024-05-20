package utils

import (
	"errors"
	"time"
)

func GetDateRange(data string) (time.Time, time.Time, error) {
	now := time.Now()
	var start, end time.Time

	switch data {
	case "m":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "y":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	default:
		return time.Time{}, time.Time{}, errors.New("invalid date")
	}

	return start, end, nil
}
