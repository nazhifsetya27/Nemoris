package database

import (
	"os"

	"nemoris/internal/model"
	"nemoris/internal/utils"
)

func Migrate() {
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)
	err := DB.AutoMigrate(
		&model.Message{},
		&model.Reminder{},
	)

	if err != nil {
		utils.LogSystem("database migration failed: " + err.Error())
		os.Exit(1)
	}

	utils.LogDB("Database Migrated")
}