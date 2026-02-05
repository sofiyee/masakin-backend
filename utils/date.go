package utils

import "time"
import "errors"

func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("empty date")
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}