package service

import (
	"context"
	"time"

	"nemoris/internal/database"
	"nemoris/internal/utils"
	"nemoris/internal/whatsapp"
)

const healthTimeout = 3 * time.Second

// HealthResult holds the aggregate health status and operational metrics.
type HealthResult struct {
	Status            string `json:"status"`
	Database          string `json:"database"`
	WAHA              string `json:"waha"`
	DBLatencyMs       int64  `json:"db_latency_ms"`
	WAHALatencyMs     int64  `json:"waha_latency_ms"`
	SchedulerLastTick string `json:"scheduler_last_tick"`
}

// CheckHealth verifies PostgreSQL and WAHA.
// Uses database.DB directly (no repository).
// getLastTick is called to obtain scheduler last tick; pass nil to omit.
// If DB fails, WAHA is set to "unknown" and not checked.
func CheckHealth(getLastTick func() time.Time) HealthResult {
	lastTick := time.Time{}
	if getLastTick != nil {
		lastTick = getLastTick()
	}

	dbOK, dbLatency := checkDatabaseLatency()
	if !dbOK {
		utils.LogSystem("health db fail")
		return HealthResult{
			Status:            "fail",
			Database:          "fail",
			WAHA:              "unknown",
			DBLatencyMs:       dbLatency.Milliseconds(),
			WAHALatencyMs:     0,
			SchedulerLastTick: formatSchedulerLastTick(lastTick),
		}
	}

	wahaOK, wahaLatency := whatsapp.CheckWAHALatency()
	if !wahaOK {
		utils.LogSystem("health waha fail")
		return HealthResult{
			Status:            "degraded",
			Database:          "ok",
			WAHA:              "fail",
			DBLatencyMs:       dbLatency.Milliseconds(),
			WAHALatencyMs:     wahaLatency.Milliseconds(),
			SchedulerLastTick: formatSchedulerLastTick(lastTick),
		}
	}

	return HealthResult{
		Status:            "ok",
		Database:          "ok",
		WAHA:              "ok",
		DBLatencyMs:       dbLatency.Milliseconds(),
		WAHALatencyMs:     wahaLatency.Milliseconds(),
		SchedulerLastTick: formatSchedulerLastTick(lastTick),
	}
}

func formatSchedulerLastTick(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.UTC().Format(time.RFC3339)
}

func checkDatabaseLatency() (bool, time.Duration) {
	if database.DB == nil {
		return false, 0
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		return false, 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()
	start := time.Now()
	err = sqlDB.PingContext(ctx)
	latency := time.Since(start)
	return err == nil, latency
}
