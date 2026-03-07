package database

import (
	"fmt"
	"log"

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
		log.Fatal("Failed to connect database:", err)
	}

	utils.LogDB("PostgreSQL connected")
}