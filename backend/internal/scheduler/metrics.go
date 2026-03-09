package scheduler

import (
	"fmt"
	"time"

	"nemoris/internal/utils"
)

// SchedulerMetrics holds execution counts for a single scheduler run.
// Scheduler only observes; service layer owns the counts.
type SchedulerMetrics struct {
	// Checked currently equals due because repository returns due-only set.
	CheckedReminders  int
	DueReminders      int
	SentReminders     int
	RetryReminders    int
	FailedReminders   int
	ClaimConflicts    int
	ExecutionDuration time.Duration
}

// NewSchedulerMetrics returns a zero-valued metrics struct.
func NewSchedulerMetrics() *SchedulerMetrics {
	return &SchedulerMetrics{}
}

// Log writes metrics to the scheduler log in the required format.
func (m *SchedulerMetrics) Log() {
	msg := fmt.Sprintf("checked=%d due=%d sent=%d retry=%d failed=%d claim_conflict=%d duration=%s",
		m.CheckedReminders, m.DueReminders, m.SentReminders,
		m.RetryReminders, m.FailedReminders, m.ClaimConflicts, m.ExecutionDuration)
	utils.LogScheduler(msg)
}
