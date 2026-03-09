package repository

import (
	"errors"

	"gorm.io/gorm"

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

// GetLatestMemoryBySender returns the most recent memory for the sender.
// Returns (model.Memory{}, nil) when no memory exists.
func GetLatestMemoryBySender(from string) (model.Memory, error) {
	var m model.Memory
	err := database.DB.
		Where(`"from" = ?`, from).
		Order("created_at desc").
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Memory{}, nil
		}
		return model.Memory{}, err
	}
	return m, nil
}

// FindMemoriesByKeyword returns memories whose content contains the keyword, for the sender.
// Newest first. Returns empty slice when no match.
func FindMemoriesByKeyword(from string, keyword string) ([]model.Memory, error) {
	var memories []model.Memory
	err := database.DB.
		Where(`"from" = ? AND content ILIKE ?`, from, "%"+keyword+"%").
		Order("created_at desc").
		Find(&memories).Error
	if err != nil {
		return nil, err
	}
	return memories, nil
}
