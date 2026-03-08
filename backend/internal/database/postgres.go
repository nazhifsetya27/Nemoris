package database

import (
	"fmt"
	"os"

	"nemoris/internal/config"
	"nemoris/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.App.DBHost,
		config.App.DBUser,
		config.App.DBPassword,
		config.App.DBName,
		config.App.DBPort,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		utils.LogSystem("database connect failed: " + err.Error())
		os.Exit(1)
	}

	utils.LogDB("PostgreSQL connected")
}