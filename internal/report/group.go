package report

import (
	"fmt"
	"time"
)

type TimeGrouper interface {
	Key(t time.Time) string      // label like "2025-08-01" or "2025-W31"
	Start(t time.Time) time.Time // normalized start of interval
	WithinRange(t, now time.Time, count int) bool
}

// DayGrouper
type DayGrouper struct{}

func (DayGrouper) Key(t time.Time) string {
	return t.Format("2006-01-02")
}

func (DayGrouper) Start(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func (DayGrouper) WithinRange(t, now time.Time, count int) bool {
	return !t.Before(now.AddDate(0, 0, -count))
}

// WeekGrouper groups by week, starting on Monday
type WeekGrouper struct{}

func (WeekGrouper) Key(t time.Time) string {
	monday := startOfWeek(t)
	return monday.Format("2006-01-02") // Or use ISO week: "2025-W31"
}

func (WeekGrouper) Start(t time.Time) time.Time {
	return startOfWeek(t)
}

func (WeekGrouper) WithinRange(t, now time.Time, count int) bool {
	tStart := startOfWeek(t)
	nowStart := startOfWeek(now)
	diff := int(nowStart.Sub(tStart).Hours() / 24 / 7)
	return diff >= 0 && diff < count
}

// MonthGrouper groups by month, starting on the first day of the month
type MonthGrouper struct{}

func (MonthGrouper) Key(t time.Time) string {
	return t.Format("2006-01")
}

func (MonthGrouper) Start(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func (MonthGrouper) WithinRange(t, now time.Time, count int) bool {
	tStart := MonthGrouper{}.Start(t)
	nowStart := MonthGrouper{}.Start(now)
	diff := (nowStart.Year()-tStart.Year())*12 + int(nowStart.Month()-tStart.Month())
	return diff >= 0 && diff < count
}

// YearGrouper groups by year, starting on January 1st
type YearGrouper struct{}

func (YearGrouper) Key(t time.Time) string {
	return fmt.Sprintf("%d", t.Year())
}

func (YearGrouper) Start(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func (YearGrouper) WithinRange(t, now time.Time, count int) bool {
	diff := now.Year() - t.Year()
	return diff >= 0 && diff < count
}
