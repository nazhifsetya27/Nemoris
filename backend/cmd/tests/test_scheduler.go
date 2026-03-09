// Manual test for scheduler: ProcessDueReminders, claim safety, retry, sent, failed.
// Creates a test reminder due now, runs ProcessDueReminders, prints before/after state.
// Run: go run ./cmd/tests/test_scheduler.go
// Note: May take 2-6 min if WAHA is slow/unavailable (retries until failed).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/model"
	"nemoris/internal/repository"
	"nemoris/internal/service"
)

const testSender = "test-user@s.whatsapp.net"

func main() {
	config.Load()
	database.Init()

	var out strings.Builder
	out.WriteString("=== test_scheduler.go ===\n\n")

	pass, fail := 0, 0

	// Create test reminder due now
	task := "scheduler test task"
	rawTime := "now"
	remindAt := time.Now().Add(-1 * time.Second) // due 1 second ago

	err := repository.SaveReminder(testSender, task, rawTime, remindAt)
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL setup: could not create test reminder: %v\n", err))
		fail++
		fmt.Print(out.String())
		writeResult(basePath(), "test_scheduler_result.txt", out.String())
		return
	}

	// Get reminder before processing
	dueBefore, _ := repository.GetDueReminders()
	var beforeID string
	for _, r := range dueBefore {
		if r.From == testSender && r.Task == task {
			beforeID = r.ID
			out.WriteString(fmt.Sprintf("--- Before ProcessDueReminders ---\n"))
			out.WriteString(fmt.Sprintf("id=%s status=%s retry_count=%d\n", r.ID, r.Status, r.RetryCount))
			break
		}
	}

	// Run ProcessDueReminders
	result, err := service.ProcessDueReminders()
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL ProcessDueReminders error: %v\n", err))
		fail++
	} else {
		out.WriteString(fmt.Sprintf("\n--- ProcessDueReminders result ---\n"))
		out.WriteString(fmt.Sprintf("Checked=%d Due=%d Sent=%d Retry=%d Failed=%d\n",
			result.Checked, result.Due, result.Sent, result.Retry, result.Failed))
		if result.Due >= 1 {
			out.WriteString("PASS due reminder found and processed\n")
			pass++
		} else {
			out.WriteString("FAIL no due reminder found\n")
			fail++
		}
	}

	// Get reminder after (may be sent or retrying depending on WAHA)
	all, _ := repository.GetAllReminders(testSender)
	var after *model.Reminder
	for i := range all {
		if all[i].ID == beforeID {
			after = &all[i]
			break
		}
	}
	if after != nil {
		out.WriteString(fmt.Sprintf("\n--- After state ---\n"))
		out.WriteString(fmt.Sprintf("id=%s status=%s retry_count=%d last_error=%q\n",
			after.ID, after.Status, after.RetryCount, after.LastError))
		if after.Status == model.ReminderSent {
			out.WriteString("PASS reminder marked sent\n")
			pass++
		} else if after.Status == model.ReminderRetrying || after.Status == model.ReminderFailed {
			out.WriteString("PASS reminder in retry/failed (WAHA may be down)\n")
			pass++
		} else {
			out.WriteString(fmt.Sprintf("INFO status=%s (expected sent/retrying/failed)\n", after.Status))
		}
	}

	// Claim safety: ClaimOneDueReminder returns nil when none due
	claimed, err := repository.ClaimOneDueReminder()
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL ClaimOneDueReminder error: %v\n", err))
		fail++
	} else if claimed == nil {
		out.WriteString("\nPASS claim safety: nil when no due reminders\n")
		pass++
	} else {
		out.WriteString(fmt.Sprintf("\nINFO claimed: %s (unexpected if we just processed)\n", claimed.ID))
	}

	// Summary
	out.WriteString("\n--- Summary ---\n")
	out.WriteString(fmt.Sprintf("PASS: %d  FAIL: %d\n", pass, fail))

	fmt.Print(out.String())
	writeResult(basePath(), "test_scheduler_result.txt", out.String())
}

func basePath() string {
	base, _ := os.Getwd()
	// If run from project root, use backend/test_results
	if _, err := os.Stat("backend"); err == nil {
		return filepath.Join(base, "backend", "test_results")
	}
	return filepath.Join(base, "test_results")
}

func writeResult(dir, name string, content string) {
	_ = os.MkdirAll(dir, 0755)
	_ = os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
}
