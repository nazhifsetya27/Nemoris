package config

import (
	"os"
	"strings"

	"nemoris/internal/utils"

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
	AppEnv        string // "dev" enables query-token fallback for webhook auth
}

var App AppConfig

func Load() {
	err := godotenv.Load()
	if err != nil {
		utils.LogSystem("Error loading .env file")
		os.Exit(1)
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
		AppEnv:        strings.TrimSpace(strings.ToLower(os.Getenv("APP_ENV"))),
	}
}