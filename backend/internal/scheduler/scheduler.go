package scheduler

import (
	"time"

	"nemoris/internal/service"
	"nemoris/internal/utils"
)

func Start() {
	// Run immediately on startup, then every 60 seconds
	utils.LogScheduler("checking due reminders")
	if err := service.ProcessDueReminders(); err != nil {
		utils.LogScheduler("scheduler error: " + err.Error())
	}

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			utils.LogScheduler("checking due reminders")

			err := service.ProcessDueReminders()
			if err != nil {
				utils.LogScheduler("scheduler error: " + err.Error())
			}
		}
	}()
}