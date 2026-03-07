package scheduler

import (
	"time"

	"nemoris/internal/service"
	"nemoris/internal/utils"
)

func Start() {
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