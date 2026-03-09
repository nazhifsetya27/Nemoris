package repository

import (
	"time"

	"nemoris/internal/database"
	"nemoris/internal/model"
)

// Create

func SaveReminder(from string, task string, rawTime string, remindAt time.Time) error {
	reminder := model.Reminder{
		From:       from,
		Task:       task,
		RawTime:    rawTime,
		RemindAt:   remindAt,
		Status:     model.ReminderPending,
		RetryCount: 0,
	}

	return database.DB.Create(&reminder).Error
}

// Read

func GetAllReminders(from string) ([]model.Reminder, error) {
	var reminders []model.Reminder

	query := database.DB.Order("created_at desc")

	if from != "" {
		query = query.Where(`"from" = ?`, from)
	}

	err := query.Find(&reminders).Error

	return reminders, err
}

func GetPendingReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder

	err := database.DB.
		Where("remind_at > NOW() AND status = ?", model.ReminderPending).
		Order("remind_at asc").
		Find(&reminders).Error

	return reminders, err
}

func GetDueReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder

	err := database.DB.
		Where("remind_at <= ? AND status IN ?", time.Now(), []string{model.ReminderPending, model.ReminderRetrying}).
		Order("remind_at asc").
		Find(&reminders).Error

	return reminders, err
}

func CountPendingReminders() (int, error) {
	var count int64
	err := database.DB.Model(&model.Reminder{}).
		Where("remind_at > ? AND status = ?", time.Now(), model.ReminderPending).
		Count(&count).Error
	return int(count), err
}

func CountOverdueReminders() (int, error) {
	var count int64
	err := database.DB.Model(&model.Reminder{}).
		Where("remind_at <= ? AND status IN ?", time.Now(), []string{model.ReminderPending, model.ReminderRetrying}).
		Count(&count).Error
	return int(count), err
}

// Lifecycle Updates

func ClaimReminder(id string) bool {
	result := database.DB.Model(&model.Reminder{}).
		Where("id = ? AND status IN ?", id, []string{model.ReminderPending, model.ReminderRetrying}).
		Update("status", model.ReminderProcessing)

	return result.RowsAffected == 1
}

func MarkReminderSent(id string) error {
	now := time.Now()

	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.ReminderSent,
			"sent_at":     &now,
			"last_error":  "",
		}).Error
}

func MarkReminderRetrying(id string, retryCount int, lastError string) error {
	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.ReminderRetrying,
			"retry_count": retryCount,
			"last_error":  lastError,
		}).Error
}

func MarkReminderFailed(id string, retryCount int, lastError string) error {
	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.ReminderFailed,
			"retry_count": retryCount,
			"last_error":  lastError,
		}).Error
}

func GetFailedReminders(from string) ([]model.Reminder, error) {
	var reminders []model.Reminder

	query := database.DB.Where("status = ?", model.ReminderFailed).Order("created_at desc")

	if from != "" {
		query = query.Where(`"from" = ?`, from)
	}

	err := query.Find(&reminders).Error

	return reminders, err
}