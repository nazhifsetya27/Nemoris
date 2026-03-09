package service

import (
	"fmt"
	"strings"
	"time"

	"nemoris/internal/model"
	"nemoris/internal/repository"
	"nemoris/internal/utils"
	"nemoris/internal/whatsapp"
)

func classifyFailure(err error, accepted bool) string {
	if !accepted && err == nil {
		return model.FailureInvalidTarget
	}
	if err == nil {
		return model.FailureUnknownSendErr
	}
	s := strings.ToLower(err.Error())
	if strings.Contains(s, "timeout") || strings.Contains(s, "deadlineexceeded") {
		return model.FailureWahaTimeout
	}
	if strings.Contains(s, "blocked") || strings.Contains(s, "invalid") || strings.Contains(s, "chatid") {
		return model.FailureInvalidTarget
	}
	if strings.Contains(s, "lock") || strings.Contains(s, "deadlock") {
		return model.FailureDBLockFail
	}
	return model.FailureUnknownSendErr
}

func handleRetry(reminderID string, retryCount int, lastError string, failureType string) {
	nextRetry := retryCount + 1

	if nextRetry >= 3 {
		err := repository.MarkReminderFailed(reminderID, nextRetry, lastError, failureType)
		if err != nil {
			utils.LogScheduler("failed mark failed: " + err.Error())
		}
		return
	}

	err := repository.MarkReminderRetrying(reminderID, nextRetry, lastError, failureType)
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
			ft := classifyFailure(sendResult.Err, sendResult.Accepted)
			handleRetry(reminder.ID, reminder.RetryCount, reason, ft)
			updateDuration = time.Since(t2)
			trackRetryOrFailed(reminder.RetryCount, &result)
		} else if !sendResult.Accepted {
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=send rejected")
			t2 := time.Now()
			ft := classifyFailure(nil, false)
			handleRetry(reminder.ID, reminder.RetryCount, "send rejected", ft)
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
