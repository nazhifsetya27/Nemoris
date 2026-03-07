package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"nemoris/internal/database"
	"nemoris/internal/handler"
	"nemoris/internal/model"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	log.Fatal("Error loading .env file")
	}

	database.Connect()

	database.DB.AutoMigrate(&model.Message{}, &model.Reminder{})

	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/webhook", handler.WebhookHandler)
	http.HandleFunc("/reminders", handler.GetReminders)
	http.HandleFunc("/reminders/pending", handler.GetPendingReminders)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}