package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"nemoris/internal/database"
	"nemoris/internal/model"
	"nemoris/internal/utils"
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

// SaveRecurringReminder creates a pending reminder with recurrence fields for the next occurrence.
func SaveRecurringReminder(from, task, rawTime string, remindAt time.Time, recurrenceType string, recurrenceInterval int) error {
	reminder := model.Reminder{
		From:                from,
		Task:                task,
		RawTime:             rawTime,
		RemindAt:            remindAt,
		Status:              model.ReminderPending,
		RetryCount:          0,
		RecurrenceType:      recurrenceType,
		RecurrenceInterval:  recurrenceInterval,
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

// ClaimOneDueReminder atomically selects one due reminder, locks it with FOR UPDATE SKIP LOCKED,
// updates status to processing, and returns it. Returns (nil, nil) when no due reminder exists.
// Safe for concurrent workers and restart overlap.
func ClaimOneDueReminder() (*model.Reminder, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var reminder model.Reminder
	err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("remind_at <= ? AND status IN ?", time.Now(), []string{model.ReminderPending, model.ReminderRetrying}).
		Order("remind_at asc").
		Limit(1).
		First(&reminder).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogDB("no due reminder")
			return nil, nil
		}
		return nil, err
	}

	err = tx.Model(&reminder).Update("status", model.ReminderProcessing).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	utils.LogDB("reminder claimed safely id=" + reminder.ID)
	return &reminder, nil
}

// ClaimDueReminderByID atomically claims a specific reminder by ID if still due.
// Returns (nil, nil) when already claimed, not due, or not found.
func ClaimDueReminderByID(id string) (*model.Reminder, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var reminder model.Reminder
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND remind_at <= ? AND status IN ?", id, time.Now(), []string{model.ReminderPending, model.ReminderRetrying}).
		First(&reminder).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	err = tx.Model(&reminder).Update("status", model.ReminderProcessing).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	utils.LogDB("reminder claimed by id id=" + reminder.ID)
	return &reminder, nil
}

func MarkReminderSent(id string) error {
	now := time.Now()

	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       model.ReminderSent,
			"sent_at":      &now,
			"last_error":   "",
			"failure_type": "",
		}).Error
}

func MarkReminderRetrying(id string, retryCount int, lastError string, failureType string) error {
	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       model.ReminderRetrying,
			"retry_count":  retryCount,
			"last_error":   lastError,
			"failure_type": failureType,
		}).Error
}

func MarkReminderFailed(id string, retryCount int, lastError string, failureType string) error {
	return database.DB.Model(&model.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       model.ReminderFailed,
			"retry_count":  retryCount,
			"last_error":   lastError,
			"failure_type": failureType,
		}).Error
}

func GetFailedReminders(from, status, failureType, date string) ([]model.Reminder, error) {
	var reminders []model.Reminder

	query := database.DB.Where("status = ?", model.ReminderFailed).Order("created_at desc")

	if from != "" {
		query = query.Where(`"from" = ?`, from)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if failureType != "" {
		query = query.Where("failure_type = ?", failureType)
	}
	if date != "" {
		if t, err := time.Parse("2006-01-02", date); err == nil {
			start := t
			end := t.AddDate(0, 0, 1)
			query = query.Where("created_at >= ? AND created_at < ?", start, end)
		}
	}

	err := query.Find(&reminders).Error

	return reminders, err
}