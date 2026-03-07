package service

import (
	"nemoris/internal/repository"
	"nemoris/internal/utils"
)

// CreateReminderFromMessage parses reminder from body, validates, saves, and returns user-facing response.
// Returns empty string when body is not a reminder (caller handles ping/default).
func CreateReminderFromMessage(from string, body string) string {
	task, rawTime, remindTime, ok := ParseReminder(body)
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