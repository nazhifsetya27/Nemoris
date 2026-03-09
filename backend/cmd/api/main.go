package main

import (
	"net/http"
	"os"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/handler"
	"nemoris/internal/middleware"
	"nemoris/internal/scheduler"
	"nemoris/internal/utils"
)

func main() {
	config.Load()
	database.Init()
	scheduler.Start()

	http.Handle("/health", middleware.Recover(http.HandlerFunc(handler.HealthCheck)))
	http.Handle("/webhook", middleware.Recover(middleware.RateLimit(middleware.WebhookAuth(http.HandlerFunc(handler.WebhookHandler)))))
	http.Handle("/reminders", middleware.Recover(http.HandlerFunc(handler.GetReminders)))
	http.Handle("/reminders/pending", middleware.Recover(http.HandlerFunc(handler.GetPendingReminders)))
	http.Handle("/reminders/failed", middleware.Recover(http.HandlerFunc(handler.GetFailedReminders)))
	http.Handle("/memories", middleware.Recover(http.HandlerFunc(handler.GetMemories)))

	utils.LogSystem("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		utils.LogSystem("server failed: " + err.Error())
		os.Exit(1)
	}
}