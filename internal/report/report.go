package report

import (
	"fmt"
	"time"

	"github.com/KEINOS/go-todotxt/todo"
)

// Interval represents the time interval for grouping tasks
func getGrouper(interval Interval) TimeGrouper {
	switch interval {
	case Day:
		return DayGrouper{}
	case Week:
		return WeekGrouper{}
	case Month:
		return MonthGrouper{}
	case Year:
		return YearGrouper{}
	default:
		return nil
	}
}

type Summary struct {
	Count     int
	TimeSpent time.Duration
}

func summarizeTasks(tasks todo.TaskList, grouper TimeGrouper, count int) map[string]Summary {
	now := time.Now()
	summary := make(map[string]Summary)

	for _, task := range tasks {
		t := task.CompletedDate
		if !grouper.WithinRange(t, now, count) {
			continue
		}

		key := grouper.Key(t)

		s := summary[key]
		s.Count++
		if timeSpent, exists := task.AdditionalTags["time"]; exists {
			duration, err := time.ParseDuration(timeSpent)
			if err == nil {
				s.TimeSpent += duration
			}
		}
		summary[key] = s
	}

	return summary
}

func Generate(tasks todo.TaskList, timeRange string) (map[string]Summary, error) {
	interval, count, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("error parsing time range: %w", err)
	}

	grouper := getGrouper(interval)
	if grouper == nil {
		return nil, fmt.Errorf("unsupported time range: %s", timeRange)
	}

	return summarizeTasks(tasks, grouper, count), nil
}
