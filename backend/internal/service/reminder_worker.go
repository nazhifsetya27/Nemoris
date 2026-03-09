package service

import (
	"fmt"

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

// ProcessDueReminders fetches due reminders, sends them, and returns execution counts.
func ProcessDueReminders() (ProcessDueRemindersResult, error) {
	result := ProcessDueRemindersResult{}

	reminders, err := repository.GetDueReminders()
	if err != nil {
		return result, err
	}

	result.Checked = len(reminders)
	result.Due = len(reminders)

	for _, reminder := range reminders {
		utils.LogScheduler("due reminder found: " + reminder.Task)

		claimed := repository.ClaimReminder(reminder.ID)
		if !claimed {
			utils.LogScheduler("already claimed: " + reminder.Task)
			result.ClaimConflicts++
			continue
		}

		text := fmt.Sprintf("Reminder: %s", reminder.Task)

		sendResult := whatsapp.SendText(reminder.From, text)
		if sendResult.Err != nil {
			reason := sendResult.Err.Error()
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=" + reason)
			handleRetry(reminder.ID, reminder.RetryCount, reason)
			trackRetryOrFailed(reminder.RetryCount, &result)
			continue
		}
		if !sendResult.Accepted {
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=send rejected")
			handleRetry(reminder.ID, reminder.RetryCount, "send rejected")
			trackRetryOrFailed(reminder.RetryCount, &result)
			continue
		}

		err = repository.MarkReminderSent(reminder.ID)
		if err != nil {
			utils.LogScheduler("failed mark sent: " + err.Error())
			continue
		}
		result.Sent++
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
