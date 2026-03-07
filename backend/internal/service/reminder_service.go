package service

import (
	"fmt"

	"nemoris/internal/ai"
	"nemoris/internal/repository"
	"nemoris/internal/utils"
)

// ListRemindersFromMessage fetches reminders for the user and returns formatted response.
func ListRemindersFromMessage(from string, lang string) string {
	reminders, err := ListAllReminders(from)
	if err != nil {
		utils.LogDB("List reminders failed: " + err.Error())
		return "internal error"
	}

	if len(reminders) == 0 {
		return "no reminders"
	}

	msg := fmt.Sprintf("you have %d reminder(s)", len(reminders))
	for i, r := range reminders {
		msg += fmt.Sprintf("\n%d. %s", i+1, r.Task)
	}
	return msg
}

// CreateReminderFromMessage parses reminder from body, validates, saves, and returns user-facing response.
// Returns empty string when body is not a reminder (caller handles ping/default).
// Caller must ensure intent is create_reminder before calling.
func CreateReminderFromMessage(from string, body string) string {
	task, rawTime, remindTime, ok := ai.ParseReminder(body)
	if !ok {
		return ""
	}

	utils.LogAI("Reminder detected")
	utils.LogAI("Task: " + task)
	utils.LogAI("RawTime: " + rawTime)
	utils.LogAI("Parsed Time: " + remindTime.String())

	if remindTime.IsZero() {
		utils.LogAI("invalid reminder time")
		return "please specify reminder time"
	}

	err := repository.SaveReminder(from, task, rawTime, remindTime)
	if err != nil {
		utils.LogDB("Reminder save failed: " + err.Error())
		return "internal error"
	}

	return "reminder noted"
}