package scheduler

import (
	"fmt"
	"time"

	"nemoris/internal/service"
	"nemoris/internal/utils"
)

const (
	tickInterval   = 60 * time.Second
	driftThreshold = 5 * time.Second
)

func Start() {
	// Startup backlog logging before first tick
	if stats, err := service.GetReminderBacklogStats(); err != nil {
		utils.LogScheduler("startup backlog error: " + err.Error())
	} else {
		utils.LogScheduler(fmt.Sprintf("startup backlog pending=%d overdue=%d", stats.Pending, stats.Overdue))
	}

	runTick()

	ticker := time.NewTicker(tickInterval)
	prevTick := time.Now()
	go func() {
		for range ticker.C {
			now := time.Now()
			actualInterval := now.Sub(prevTick)
			if actualInterval > tickInterval+driftThreshold {
				utils.LogScheduler(fmt.Sprintf("scheduler drift detected: %s", actualInterval))
			}
			prevTick = now
			runTick()
		}
	}()
}

func runTick() {
	utils.LogScheduler("checking due reminders")

	metrics := NewSchedulerMetrics()
	start := time.Now()

	result, err := service.ProcessDueReminders()
	if err != nil {
		utils.LogScheduler("scheduler error: " + err.Error())
	}

	metrics.CheckedReminders = result.Checked
	metrics.DueReminders = result.Due
	metrics.SentReminders = result.Sent
	metrics.RetryReminders = result.Retry
	metrics.FailedReminders = result.Failed
	metrics.ClaimConflicts = result.ClaimConflicts
	metrics.ExecutionDuration = time.Since(start)

	metrics.Log()
}