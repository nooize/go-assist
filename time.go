package assist

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const JsonDateFormat = "2006-01-02"
const JsonTimeFormat = "15:04:05" // "2006-01-02T15:04:05.999Z"
const JsonDateTimeFormat = JsonDateFormat + "T" + JsonTimeFormat
const JsonDateTimeWithZoneFormat = JsonDateTimeFormat + " MST"
const jsonDateWithZoneFormat = JsonDateFormat + " MST"

type JsonTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *JsonTime) UnmarshalJSON(bytes []byte) error {
	str := strings.Trim(string(bytes), "\"")
	if str == "null" {
		return nil
	}
	timeFormat := time.RFC3339
	strLen := len(str)
	switch {
	case strLen == 0:
		return nil
	case strLen == 13:
		// unix timestamp
		i, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			*t = JsonTime{time.Unix(0, i*int64(time.Millisecond))}
		}
		return err
	case strLen > 7 && strLen < 10:
		offset, err := dateOffset(str)
		if err != nil {
			return err
		}
		str += "T00:00:00" + offset
		timeFormat = "2006-1-2T15:04:05Z07:00"
	case strLen == len(JsonDateFormat):
		// detect time offset depend from date
		// same locations change offset during the year
		offset, err := dateOffset(str)
		if err != nil {
			return err
		}
		str += "T00:00:00" + offset
	case strLen == len(JsonDateTimeFormat):
		offset, err := dateOffset(str[:10])
		if err != nil {
			return err
		}
		str = str[:10] + "T" + str[11:] + " " + offset
	case strLen == len(time.RFC3339Nano):
		timeFormat = time.RFC3339Nano
	default:
		// for all other cases we will try RFC 3339
	}
	nt, err := time.Parse(timeFormat, str)
	if err != nil {
		return errors.New("expect time in RFC3339 or timestamp, has: " + str)
	}
	*t = JsonTime{nt}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))), nil
}

// IsTimeZero check if tiem is 00:00:00
func IsTimeZero(t *time.Time) bool {
	return t != nil && t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}

// StartOfTheDay return new time.Time with 00:00:00 time
func StartOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfTheDay return new time.Time with 23.59.59.999999999 time
func EndOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func TrimToMilliseconds(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000000*1000000000, t.Location())
}

// TrimToMicroseconds trim time to microseconds
func TrimToMicroseconds(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000*1000000, t.Location())
}

func IsSameDay(d1, d2 time.Time) bool {
	ld2 := d2.In(d1.Location())
	return d1.Year() == ld2.Year() && d1.YearDay() == ld2.YearDay()
}

func DayPeriod(t time.Time) (time.Time, time.Time) {
	return StartOfTheDay(t), EndOfTheDay(t)
}

func WeekPeriod(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()
	weekDay := t.Weekday() - 1
	if weekDay < 0 {
		weekDay = 6
	}
	from := time.Date(year, month, day-int(weekDay), 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return from, to
}

func MonthPeriod(t time.Time) (time.Time, time.Time) {
	year, month, _ := t.Date()
	from := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return from, to
}

func ParseFromTo(fromStr, toStr string) (from, to time.Time, err error) {
	tmp := JsonTime{}
	if err = tmp.UnmarshalJSON([]byte(fromStr)); err != nil {
		return
	}
	from = tmp.Time
	if len(toStr) == 0 {
		from = StartOfTheDay(from)
		to = EndOfTheDay(from)
	} else {
		tmp = JsonTime{}
		if err = tmp.UnmarshalJSON([]byte(toStr)); err != nil {
			return
		}
		to = tmp.Time
	}
	if from.After(to) {
		err = errors.New("to date must be after from")
		return
	}
	return
}

// this func return time offset for date
// same location may change time zone depends
// from time of the year
func dateOffset(date string) (string, error) {
	zone, _ := time.Now().Zone()
	d, err := time.Parse(jsonDateWithZoneFormat, date+" "+zone)
	if err != nil {
		return "", err
	}
	return d.Format("Z07:00"), nil
}
