package service

import (
	"nemoris/internal/database"
	"nemoris/internal/utils"
	"nemoris/internal/whatsapp"
)

// HealthResult holds the aggregate health status.
type HealthResult struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	WAHA     string `json:"waha"`
}

// CheckHealth verifies PostgreSQL and WAHA.
// Uses database.DB directly (no repository).
// If DB fails, WAHA is set to "unknown" and not checked.
func CheckHealth() HealthResult {
	dbOK := checkDatabase()
	if !dbOK {
		utils.LogSystem("health db fail")
		return HealthResult{
			Status:   "fail",
			Database: "fail",
			WAHA:     "unknown",
		}
	}

	wahaOK := whatsapp.CheckWAHA()
	if !wahaOK {
		utils.LogSystem("health waha fail")
		return HealthResult{
			Status:   "degraded",
			Database: "ok",
			WAHA:     "fail",
		}
	}

	return HealthResult{
		Status:   "ok",
		Database: "ok",
		WAHA:     "ok",
	}
}

func checkDatabase() bool {
	if database.DB == nil {
		return false
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		return false
	}
	err = sqlDB.Ping()
	return err == nil
}
