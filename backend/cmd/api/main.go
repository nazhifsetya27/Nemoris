package main

import (
	"log"
	"net/http"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/handler"
	"nemoris/internal/middleware"
	"nemoris/internal/scheduler"
)

func main() {
	config.Load()
	database.Init()
	scheduler.Start()

	http.HandleFunc("/health", handler.HealthCheck)
	http.Handle("/webhook", middleware.WebhookAuth(http.HandlerFunc(handler.WebhookHandler)))
	http.HandleFunc("/reminders", handler.GetReminders)
	http.HandleFunc("/reminders/pending", handler.GetPendingReminders)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}