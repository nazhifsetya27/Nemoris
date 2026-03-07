package service

import (
	"fmt"

	"nemoris/internal/repository"
	"nemoris/internal/utils"
	"nemoris/internal/whatsapp"
)

func handleRetry(reminderID string, retryCount int) {
	nextRetry := retryCount + 1

	if nextRetry >= 3 {
		err := repository.MarkReminderFailed(reminderID, nextRetry)
		if err != nil {
			utils.LogScheduler("failed mark failed: " + err.Error())
		}
		return
	}

	err := repository.MarkReminderRetrying(reminderID, nextRetry)
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
			utils.LogScheduler("send failed: to=" + reminder.From + " task=" + reminder.Task + " err=" + result.Err.Error())
			handleRetry(reminder.ID, reminder.RetryCount)
			continue
		}
		if !result.Accepted {
			utils.LogScheduler("send failed: to=" + reminder.From + " task=" + reminder.Task)
			handleRetry(reminder.ID, reminder.RetryCount)
			continue
		}

		err = repository.MarkReminderSent(reminder.ID)
		if err != nil {
			utils.LogScheduler("failed mark sent: " + err.Error())
		}
	}

	return nil
}
