package main

import (
	"log"
	"net/http"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/handler"
	"nemoris/internal/scheduler"
)

func main() {
	config.Load()

	database.Connect()
	database.Migrate()

	scheduler.Start()

	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/webhook", handler.WebhookHandler)
	http.HandleFunc("/reminders", handler.GetReminders)
	http.HandleFunc("/reminders/pending", handler.GetPendingReminders)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}