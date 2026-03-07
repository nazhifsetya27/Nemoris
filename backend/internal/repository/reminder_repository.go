package repository

import (
	"time"

	"nemoris/internal/database"
	"nemoris/internal/model"
)

func SaveReminder(from string, task string, rawTime string, remindAt time.Time) error {
	reminder := model.Reminder{
		From:     from,
		Task:     task,
		RawTime:  rawTime,
		RemindAt: remindAt,
	}

	return database.DB.Create(&reminder).Error
}

func GetAllReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder

	err := database.DB.Order("id desc").Find(&reminders).Error

	return reminders, err
}

func GetPendingReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder

	err := database.DB.
		Where("remind_at > NOW()").
		Order("remind_at asc").
		Find(&reminders).Error

	return reminders, err
}