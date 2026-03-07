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

func ProcessDueReminders() error {
	reminders, err := repository.GetDueReminders()
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		utils.LogScheduler("due reminder found: " + reminder.Task)

		claimed := repository.ClaimReminder(reminder.ID)
		if !claimed {
			utils.LogScheduler("already claimed: " + reminder.Task)
			continue
		}

		text := fmt.Sprintf("Reminder: %s", reminder.Task)

		result := whatsapp.SendText(reminder.From, text)
		if result.Err != nil {
			reason := result.Err.Error()
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=" + reason)
			handleRetry(reminder.ID, reminder.RetryCount, reason)
			continue
		}
		if !result.Accepted {
			utils.LogScheduler("send failed: id=" + reminder.ID + " task=" + reminder.Task + " err=send rejected")
			handleRetry(reminder.ID, reminder.RetryCount, "send rejected")
			continue
		}

		err = repository.MarkReminderSent(reminder.ID)
		if err != nil {
			utils.LogScheduler("failed mark sent: " + err.Error())
		}
	}

	return nil
}
