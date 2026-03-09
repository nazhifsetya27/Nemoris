package service

import (
	"strings"

	"nemoris/internal/ai"
	"nemoris/internal/config"
	"nemoris/internal/formatter"
	"nemoris/internal/i18n"
	"nemoris/internal/model"
	"nemoris/internal/repository"
	"nemoris/internal/utils"
)

func ProcessMessage(from string, body string) string {
	if from == config.App.BotNumber {
		utils.LogInbound("Self message ignored")
		return ""
	}

	normalizedBody := formatter.NormalizeInput(body)
	isDuplicate, err := repository.FindRecentDuplicate(from, normalizedBody)
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

	parsed := ai.Parse(normalizedBody)

	if parsed.Intent == "create_reminder" {
		if parsed.Time == "" {
			if parsed.Lang == "id" {
				return "tolong tentukan waktu pengingat"
			}
			return "please specify reminder time"
		}
		if response := CreateReminderFromMessage(from, normalizedBody, parsed.Lang); response != "" {
			return response
		}
	}

	if parsed.Intent == "list_reminders" {
		return ListRemindersFromMessage(from, parsed.Lang)
	}

	if parsed.Intent == "store_memory" {
		m := model.Memory{
			From:    from,
			Content: normalizedBody,
			Lang:    parsed.Lang,
		}
		if err := repository.SaveMemory(m); err != nil {
			utils.LogDB("Memory save failed: " + err.Error())
			return "internal error"
		}
		return i18n.BuildFromIntent(parsed.Lang, "store_memory", nil)
	}

	if parsed.Intent == "retrieve_memory" {
		return RetrieveMemoryReply(from, normalizedBody, parsed.Lang)
	}

	switch strings.ToLower(normalizedBody) {
	case "ping":
		return "pong"
	default:
		return i18n.BuildFromIntent(parsed.Lang, "unknown", nil)
	}
}