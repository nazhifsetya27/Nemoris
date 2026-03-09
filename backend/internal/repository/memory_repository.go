package repository

import (
	"nemoris/internal/database"
	"nemoris/internal/model"
)

func SaveMemory(m model.Memory) error {
	return database.DB.Create(&m).Error
}

func ListMemoriesByFrom(from string) ([]model.Memory, error) {
	var memories []model.Memory
	err := database.DB.
		Where(`"from" = ?`, from).
		Order("created_at desc").
		Find(&memories).Error
	return memories, err
}
