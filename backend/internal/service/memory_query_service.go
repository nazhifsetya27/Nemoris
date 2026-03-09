package service

import (
	"nemoris/internal/model"
	"nemoris/internal/repository"
)

// ListMemories returns memories filtered by sender, newest first.
func ListMemories(from string) ([]model.Memory, error) {
	return repository.ListMemoriesByFrom(from)
}
