package database

import (
	"log"

	"nemoris/internal/model"
	"nemoris/internal/utils"
)

func Migrate() {
	err := DB.AutoMigrate(
		&model.Message{},
		&model.Reminder{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	utils.LogDB("Database Migrated")
}