package repository

import (
	"nemoris/internal/database"
	"nemoris/internal/model"
)

func SaveMemory(m model.Memory) error {
	return database.DB.Create(&m).Error
}
