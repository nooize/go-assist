package assist

import (
	"errors"
	"time"
)

const DateFormat = "2006-01-02"

func IsTimeZero(t *time.Time) bool {
	return t != nil && t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}

func StartOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func IsSameDay(d1 time.Time, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.YearDay() == d1.YearDay()
}

func ParseFromTo(fstr string, tstr string) (from time.Time, to time.Time, err error) {
	from, err = time.Parse(DateFormat, fstr)
	if err != nil {
		return
	}
	from = StartOfTheDay(from)
	if len(tstr) > 0 {
		to, err = time.Parse(DateFormat, tstr)
		if err != nil {
			return
		}
	}
	to = EndOfTheDay(to)
	if from.After(to) {
		err = errors.New("to date must be after from")
		return
	}
	return
}


