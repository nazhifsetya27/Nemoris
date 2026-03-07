package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	BotNumber     string
	WAHABaseURL   string
	WAHAAPIKey    string
	WebhookSecret string
}

var App AppConfig

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wahaURL := os.Getenv("WAHA_BASE_URL")
	if wahaURL == "" {
		wahaURL = "http://localhost:3000"
	}
	App = AppConfig{
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		BotNumber:     os.Getenv("BOT_NUMBER"),
		WAHABaseURL:   wahaURL,
		WAHAAPIKey:    os.Getenv("WAHA_API_KEY"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
	}
}