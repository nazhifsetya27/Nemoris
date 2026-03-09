// Manual test for health: DB health, WAHA health, scheduler metrics availability.
// Run: go run ./cmd/tests/test_health.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/scheduler"
	"nemoris/internal/service"
)

func main() {
	config.Load()
	database.Init()

	var out strings.Builder
	out.WriteString("=== test_health.go ===\n\n")

	pass, fail := 0, 0

	// DB health
	hr := service.CheckHealth()
	out.WriteString(fmt.Sprintf("--- Health Check ---\n"))
	out.WriteString(fmt.Sprintf("status=%s database=%s waha=%s\n", hr.Status, hr.Database, hr.WAHA))

	if hr.Database == "ok" {
		out.WriteString("PASS DB health\n")
		pass++
	} else {
		out.WriteString("FAIL DB health\n")
		fail++
	}

	if hr.WAHA == "ok" || hr.WAHA == "fail" {
		out.WriteString("PASS WAHA health (checked)\n")
		pass++
	} else {
		out.WriteString("INFO WAHA unknown (DB may have failed first)\n")
	}

	// Scheduler metrics availability
	m := scheduler.NewSchedulerMetrics()
	m.CheckedReminders = 0
	m.DueReminders = 0
	m.SentReminders = 0
	m.RetryReminders = 0
	m.FailedReminders = 0
	m.ClaimConflicts = 0
	m.Log() // should not panic
	out.WriteString("\nPASS scheduler metrics (NewSchedulerMetrics, Log)\n")
	pass++

	// GetReminderBacklogStats (scheduler startup uses this)
	stats, err := service.GetReminderBacklogStats()
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL GetReminderBacklogStats: %v\n", err))
		fail++
	} else {
		out.WriteString(fmt.Sprintf("PASS GetReminderBacklogStats: pending=%d overdue=%d\n", stats.Pending, stats.Overdue))
		pass++
	}

	out.WriteString("\n--- Summary ---\n")
	out.WriteString(fmt.Sprintf("PASS: %d  FAIL: %d\n", pass, fail))

	fmt.Print(out.String())
	writeResult(out.String())
}

func writeResult(content string) {
	base, _ := os.Getwd()
	if _, err := os.Stat("backend"); err == nil {
		base = filepath.Join(base, "backend")
	}
	p := filepath.Join(base, "test_results", "test_health_result.txt")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte(content), 0644)
}
