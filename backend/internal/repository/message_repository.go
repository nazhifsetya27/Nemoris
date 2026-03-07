package repository

import (
	"nemoris/internal/database"
	"nemoris/internal/model"
)

func SaveMessage(from string, body string) error {
	message := model.Message{
		From: from,
		Body: body,
	}

	return database.DB.Create(&message).Error
}

func FindRecentDuplicate(from string, body string) (bool, error) {
	var count int64

	err := database.DB.Model(&model.Message{}).
		Where(`"from" = ? AND body = ? AND created_at >= NOW() - INTERVAL '10 seconds'`, from, body).
		Count(&count).Error

	return count > 0, err
}