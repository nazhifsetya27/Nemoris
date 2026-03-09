package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nemoris/internal/model"
	"nemoris/internal/queue"
	"nemoris/internal/redis"
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

// ProcessDueReminders runs hybrid flow: detect → enqueue → process from queue.
// Falls back to claim loop when queue unavailable.
func ProcessDueReminders() (ProcessDueRemindersResult, error) {
	due, err := repository.GetDueReminders()
	if err != nil {
		return ProcessDueRemindersResult{}, err
	}
	if len(due) == 0 {
		return processViaClaimLoop()
	}

	// Try queue path; fallback only when first enqueue fails (no partial enqueue)
	if err := queue.Enqueue(due[0].ID); err != nil {
		utils.LogSchedulerWarn("queue enqueue failed, falling back to claim loop: " + err.Error())
		return processViaClaimLoop()
	}
	for _, d := range due[1:] {
		_ = queue.Enqueue(d.ID)
	}
	return processFromQueue()
}

// processFromQueue dequeues and executes. Execution reads through queue path.
func processFromQueue() (ProcessDueRemindersResult, error) {
	result := ProcessDueRemindersResult{}

	for {
		reminderID, err := queue.Dequeue()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			utils.LogSchedulerWarn("queue dequeue failed: " + err.Error())
			break
		}

		reminder, err := repository.ClaimDueReminderByID(reminderID)
		if err != nil {
			return result, err
		}
		if reminder == nil {
			continue
		}

		acquired, lockErr := queue.TryLock(reminder.ID)
		if lockErr != nil {
			utils.LogSchedulerWarn("queue lock failed, proceeding (fallback): " + lockErr.Error())
		}
		if !acquired && lockErr == nil {
			utils.LogSchedulerWarn("Redis lock held for id=" + reminder.ID + ", proceeding (fallback)")
		}

		result.Checked++
		result.Due++
		executeOneReminder(reminder, &result)
		_ = queue.Release(reminder.ID)
	}

	return result, nil
}

// processViaClaimLoop is the fallback path when queue unavailable.
func processViaClaimLoop() (ProcessDueRemindersResult, error) {
	result := ProcessDueRemindersResult{}

	for {
		reminder, err := repository.ClaimOneDueReminder()
		if err != nil {
			return result, err
		}
		if reminder == nil {
			break
		}

		acquired, _ := redis.TryLock(reminder.ID)
		if !acquired {
			utils.LogSchedulerWarn("Redis lock held for id=" + reminder.ID + ", proceeding (fallback)")
		}

		result.Checked++
		result.Due++
		executeOneReminder(reminder, &result)
		redis.Release(reminder.ID)
	}

	return result, nil
}

func executeOneReminder(reminder *model.Reminder, result *ProcessDueRemindersResult) {
	utils.LogScheduler("due reminder found: " + reminder.Task)

	text := fmt.Sprintf("Reminder: %s", reminder.Task)
	sendResult := whatsapp.SendText(reminder.From, text)

	if sendResult.Err != nil {
		reason := sendResult.Err.Error()
		utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=" + reason)
		ft := classifyFailure(sendResult.Err, sendResult.Accepted)
		handleRetry(reminder.ID, reminder.RetryCount, reason, ft)
		trackRetryOrFailed(reminder.RetryCount, result)
		return
	}
	if !sendResult.Accepted {
		utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=send rejected")
		ft := classifyFailure(nil, false)
		handleRetry(reminder.ID, reminder.RetryCount, "send rejected", ft)
		trackRetryOrFailed(reminder.RetryCount, result)
		return
	}

	if err := repository.MarkReminderSent(reminder.ID); err != nil {
		utils.LogScheduler("failed mark sent: " + err.Error())
		return
	}
	result.Sent++
	if reminder.RecurrenceType != "" {
		nextAt, calcErr := NextOccurrence(reminder.RemindAt, reminder.RecurrenceType, time.Now())
		if calcErr != nil {
			utils.LogScheduler("recurrence calc failed: id=" + reminder.ID + " err=" + calcErr.Error())
			return
		}
		if createErr := repository.SaveRecurringReminder(reminder.From, reminder.Task, reminder.RawTime, nextAt, reminder.RecurrenceType, reminder.RecurrenceInterval); createErr != nil {
			utils.LogScheduler("recurrence create failed: id=" + reminder.ID + " err=" + createErr.Error())
		}
	}
}

// trackRetryOrFailed increments Retry or Failed based on retry count.
func trackRetryOrFailed(retryCount int, r *ProcessDueRemindersResult) {
	if retryCount+1 >= 3 {
		r.Failed++
	} else {
		r.Retry++
	}
}
