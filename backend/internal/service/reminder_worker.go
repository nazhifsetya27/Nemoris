package service

import (
	"fmt"
	"time"

	"nemoris/internal/repository"
	"nemoris/internal/utils"
	"nemoris/internal/whatsapp"
)

func handleRetry(reminderID string, retryCount int, lastError string) {
	nextRetry := retryCount + 1

	if nextRetry >= 3 {
		err := repository.MarkReminderFailed(reminderID, nextRetry, lastError)
		if err != nil {
			utils.LogScheduler("failed mark failed: " + err.Error())
		}
		return
	}

	err := repository.MarkReminderRetrying(reminderID, nextRetry, lastError)
	if err != nil {
		utils.LogScheduler("failed mark retrying: " + err.Error())
	}
}

// ReminderBacklogStats holds pending and overdue counts for startup logging.
type ReminderBacklogStats struct {
	Pending int
	Overdue int
}

// GetReminderBacklogStats returns pending and overdue reminder counts.
func GetReminderBacklogStats() (ReminderBacklogStats, error) {
	stats := ReminderBacklogStats{}
	pending, err := repository.CountPendingReminders()
	if err != nil {
		return stats, err
	}
	overdue, err := repository.CountOverdueReminders()
	if err != nil {
		return stats, err
	}
	stats.Pending = pending
	stats.Overdue = overdue
	return stats, nil
}

// ProcessDueRemindersResult holds counts from a single run for scheduler metrics.
type ProcessDueRemindersResult struct {
	Checked       int
	Due          int
	Sent         int
	Retry        int
	Failed       int
	ClaimConflicts int
}

// ProcessDueReminders claims due reminders one-by-one (transaction-safe), sends them, and returns execution counts.
func ProcessDueReminders() (ProcessDueRemindersResult, error) {
	result := ProcessDueRemindersResult{}

	for {
		t0 := time.Now()
		reminder, err := repository.ClaimOneDueReminder()
		claimDuration := time.Since(t0)
		if err != nil {
			return result, err
		}
		if reminder == nil {
			break
		}

		result.Checked++
		result.Due++
		utils.LogScheduler("due reminder found: " + reminder.Task)

		text := fmt.Sprintf("Reminder: %s", reminder.Task)

		t1 := time.Now()
		sendResult := whatsapp.SendText(reminder.From, text)
		sendDuration := time.Since(t1)

		var updateDuration time.Duration
		if sendResult.Err != nil {
			reason := sendResult.Err.Error()
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=" + reason)
			t2 := time.Now()
			handleRetry(reminder.ID, reminder.RetryCount, reason)
			updateDuration = time.Since(t2)
			trackRetryOrFailed(reminder.RetryCount, &result)
		} else if !sendResult.Accepted {
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=send rejected")
			t2 := time.Now()
			handleRetry(reminder.ID, reminder.RetryCount, "send rejected")
			updateDuration = time.Since(t2)
			trackRetryOrFailed(reminder.RetryCount, &result)
		} else {
			t2 := time.Now()
			err = repository.MarkReminderSent(reminder.ID)
			updateDuration = time.Since(t2)
			if err != nil {
				utils.LogScheduler("failed mark sent: " + err.Error())
			} else {
				result.Sent++
			}
		}
		utils.LogScheduler(fmt.Sprintf("profile claim=%v send=%v update=%v id=%s", claimDuration, sendDuration, updateDuration, reminder.ID))
	}

	return result, nil
}

// trackRetryOrFailed increments Retry or Failed based on retry count.
func trackRetryOrFailed(retryCount int, r *ProcessDueRemindersResult) {
	if retryCount+1 >= 3 {
		r.Failed++
	} else {
		r.Retry++
	}
}
