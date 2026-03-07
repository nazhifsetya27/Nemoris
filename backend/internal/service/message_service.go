package service

import (
	"log"
	"os"
	"strings"

	"nemoris/internal/repository"
)

func ProcessMessage(from string, body string) string {
	if from == os.Getenv("BOT_NUMBER") {
	log.Println("Self message ignored")
	return "" 
	}

	isDuplicate, err := repository.FindRecentDuplicate(from, body)
	if err != nil {
		log.Println("Duplicate check failed:", err)
		return "internal error"
	}

	if isDuplicate {
		log.Println("Duplicate ignored")
		return "duplicate ignored"
	}

	err = repository.SaveMessage(from, body)
	if err != nil {
		log.Println("Save failed:", err)
		return "internal error"
	}

	task, rawTime, remindTime, ok := ParseReminder(body)
	if ok {
	log.Println("Reminder detected")
	log.Println("Task:", task)
	log.Println("RawTime:", rawTime)
	log.Println("Parsed Time:", remindTime)

	err = repository.SaveReminder(from, task, rawTime, remindTime)
	if err != nil {
		log.Println("Reminder save failed:", err)
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