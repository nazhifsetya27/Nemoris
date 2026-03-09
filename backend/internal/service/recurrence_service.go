package service

import (
	"errors"
	"time"
)

var ErrUnsupportedRecurrenceType = errors.New("unsupported recurrence type")

// NextOccurrence returns the next valid future time for a recurring reminder.
// baseTime: the current reminder time (timezone preserved).
// recurrenceType: "daily", "weekly", or "monthly".
// after: cutoff; returned time is strictly after this (never past).
func NextOccurrence(baseTime time.Time, recurrenceType string, after time.Time) (time.Time, error) {
	switch recurrenceType {
	case "daily":
		return nextDaily(baseTime, after), nil
	case "weekly":
		return nextWeekly(baseTime, after), nil
	case "monthly":
		return nextMonthly(baseTime, after), nil
	default:
		return time.Time{}, ErrUnsupportedRecurrenceType
	}
}

func nextDaily(baseTime time.Time, after time.Time) time.Time {
	next := baseTime
	for !next.After(after) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

func nextWeekly(baseTime time.Time, after time.Time) time.Time {
	next := baseTime
	for !next.After(after) {
		next = next.AddDate(0, 0, 7)
	}
	return next
}

func nextMonthly(baseTime time.Time, after time.Time) time.Time {
	next := baseTime
	for !next.After(after) {
		next = next.AddDate(0, 1, 0)
	}
	return next
}
