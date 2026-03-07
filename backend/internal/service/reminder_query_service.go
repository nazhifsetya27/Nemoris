package service

import (
	"nemoris/internal/model"
	"nemoris/internal/repository"
)

// ListAllReminders returns all reminders, optionally filtered by from.
func ListAllReminders(from string) ([]model.Reminder, error) {
	return repository.GetAllReminders(from)
}

// ListPendingReminders returns reminders that are pending (future, not yet sent).
func ListPendingReminders() ([]model.Reminder, error) {
	return repository.GetPendingReminders()
}
