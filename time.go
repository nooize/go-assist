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

// spetial type for universal json parse unix time stamp or time string
type JsonTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *JsonTime) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}
	str := strings.Trim(string(bytes), "\"")
	fmt := ""
	strLen := len(str)
	switch {
	case strLen == 0:
		return nil
	case strLen == 13:
		i, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			*t = JsonTime{time.Unix(0, i*int64(time.Millisecond))}
		}
		return err
	case strLen > 7 && strLen < 10:
		str += "T00:00:00" + time.Now().Format("Z07:00")
		fmt = "2006-1-2T15:04:05Z07:00"
	case strLen == len(JsonDateFormat):
		str += "T00:00:00" + time.Now().Format("Z07:00")
		fmt = time.RFC3339
	case strLen == len(JsonDateTimeFormat):
		str = str[:10] + "T" + str[11:] + time.Now().Format("Z07:00")
		fmt = time.RFC3339
	case strLen == len(time.RFC3339):
		fmt = time.RFC3339
	case strLen == len(time.RFC3339Nano):
		fmt = time.RFC3339Nano
	default:
		return errors.New("time must be in RFC3339 or timestamp format")
	}
	nt, err := time.Parse(fmt, str)
	if err == nil {
		*t = JsonTime{nt}
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))), nil
}

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
	vd1 := d1.UTC()
	vd2 := d2.UTC()
	return vd1.Year() == vd2.Year() && vd1.YearDay() == vd2.YearDay()
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
	from := time.Date(year, month-1, 1, 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return from, to
}

func ParseFromTo(fromStr string, toStr string) (from time.Time, to time.Time, err error) {
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
