package report

import (
	"errors"
	"strconv"
	"time"
)

type Interval string

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

func parseTimeRange(input string) (Interval, int, error) {
	if len(input) < 2 {
		return "", 0, errors.New("invalid time range")
	}

	unit := input[len(input)-1]
	count, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		return "", 0, err
	}

	switch unit {
	case 'd':
		return Day, count, nil
	case 'w':
		return Week, count, nil
	case 'm':
		return Month, count, nil
	case 'y':
		return Year, count, nil
	default:
		return "", 0, errors.New("unsupported unit")
	}
}

// CalculateOldest calculates the oldest date from a given string input.
func CalculateOldest(input string) (time.Time, error) {
	unit, count, err := parseTimeRange(input)
	if err != nil {
		return time.Time{}, err
	}
	count = count - 1 // Adjust count to handle the current interval
	now := time.Now()
	switch unit {
	case Day:
		return now.AddDate(0, 0, -count), nil
	case Week:
		now = startOfWeek(now)
		return now.AddDate(0, 0, -7*count), nil
	case Month:
		now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return now.AddDate(0, -count, 0), nil
	case Year:
		now = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return now.AddDate(-count, 0, 0), nil
	default:
		return time.Time{}, errors.New("unsupported unit")
	}
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday → 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
}
