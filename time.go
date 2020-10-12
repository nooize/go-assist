package assist

import (
	"time"
)

const DateFormat = "2006-01-02"

func IsTimeZero(t *time.Time) bool {
	return t != nil && t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}

func StartOfTheDay(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	res := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &res
}

func EndOfTheDay(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	res := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	return &res
}
