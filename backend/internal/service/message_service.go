package service

import (
	"strings"

	"nemoris/internal/ai"
	"nemoris/internal/config"
	"nemoris/internal/repository"
	"nemoris/internal/utils"
)

func ProcessMessage(from string, body string) string {
	if from == config.App.BotNumber {
		utils.LogInbound("Self message ignored")
		return ""
	}

	isDuplicate, err := repository.FindRecentDuplicate(from, body)
	if err != nil {
		utils.LogDB("Duplicate check failed: " + err.Error())
		return "internal error"
	}

	if isDuplicate {
		utils.LogInbound("Duplicate ignored")
		return "duplicate ignored"
	}

	err = repository.SaveMessage(from, body)
	if err != nil {
		utils.LogDB("Message save failed: " + err.Error())
		return "internal error"
	}

	parsed := ai.Parse(body)

	if parsed.Intent == "create_reminder" {
		if response := CreateReminderFromMessage(from, body); response != "" {
			return response
		}
	}

	if parsed.Intent == "list_reminders" {
		return ListRemindersFromMessage(from, parsed.Lang)
	}

	switch strings.ToLower(body) {
	case "ping":
		return "pong"
	default:
		return "message received"
	}
}