package service

import (
	"strings"

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

	task, rawTime, remindTime, ok := ParseReminder(body)
	if ok {
	utils.LogAI("Reminder detected")
	utils.LogAI("Task: " + task)
	utils.LogAI("RawTime: " + rawTime)
	utils.LogAI("Parsed Time: " + remindTime.String())

	err = repository.SaveReminder(from, task, rawTime, remindTime)
	if err != nil {
		utils.LogDB("Reminder save failed: " + err.Error())
		return "internal error"
	}

	return "reminder noted"
	}

	switch strings.ToLower(body) {
	case "ping":
		return "pong"
	default:
		return "message received"
	}
}