// Manual test for repository: SaveReminder, GetDueReminders, ClaimOneDueReminder,
// MarkReminderSent, MarkReminderRetrying, MarkReminderFailed.
// Run: go run ./cmd/tests/test_repository.go
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
)

const testSender = "test-user@s.whatsapp.net"

func main() {
	config.Load()
	database.Init()

	var out strings.Builder
	out.WriteString("=== test_repository.go ===\n\n")

	pass, fail := 0, 0

	// SaveReminder
	task := "repo test task"
	rawTime := "tomorrow 10am"
	remindAt := time.Now().Add(24 * time.Hour)
	err := repository.SaveReminder(testSender, task, rawTime, remindAt)
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL SaveReminder: %v\n", err))
		fail++
	} else {
		out.WriteString("PASS SaveReminder\n")
		pass++
	}

	// GetDueReminders (our reminder is future, so may not appear)
	due, err := repository.GetDueReminders()
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL GetDueReminders: %v\n", err))
		fail++
	} else {
		out.WriteString(fmt.Sprintf("PASS GetDueReminders (count=%d)\n", len(due)))
		pass++
	}

	// Find our test reminder
	all, _ := repository.GetAllReminders(testSender)
	var testReminder *model.Reminder
	for i := range all {
		if all[i].Task == task && all[i].From == testSender {
			testReminder = &all[i]
			break
		}
	}
	if testReminder == nil {
		out.WriteString("FAIL could not find test reminder\n")
		fail++
		fmt.Print(out.String())
		writeResult(out.String())
		return
	}

	// MarkReminderRetrying
	err = repository.MarkReminderRetrying(testReminder.ID, 1, "test retry", model.FailureUnknownSendErr)
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL MarkReminderRetrying: %v\n", err))
		fail++
	} else {
		out.WriteString("PASS MarkReminderRetrying\n")
		pass++
	}

	// MarkReminderFailed (simulate after 3 retries)
	err = repository.MarkReminderFailed(testReminder.ID, 3, "test failed", model.FailureUnknownSendErr)
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL MarkReminderFailed: %v\n", err))
		fail++
	} else {
		out.WriteString("PASS MarkReminderFailed\n")
		pass++
	}

	// Verify failed state
	failedList, _ := repository.GetFailedReminders(testSender)
	found := false
	for _, r := range failedList {
		if r.ID == testReminder.ID {
			found = true
			break
		}
	}
	if found {
		out.WriteString("PASS reminder in failed list\n")
		pass++
	} else {
		out.WriteString("FAIL reminder not in failed list\n")
		fail++
	}

	// Create another for ClaimOneDueReminder and MarkReminderSent
	remindNow := time.Now().Add(-1 * time.Second)
	err = repository.SaveReminder(testSender, "claim test", "now", remindNow)
	if err != nil {
		out.WriteString(fmt.Sprintf("FAIL SaveReminder for claim test: %v\n", err))
		fail++
	} else {
		claimed, err := repository.ClaimOneDueReminder()
		if err != nil {
			out.WriteString(fmt.Sprintf("FAIL ClaimOneDueReminder: %v\n", err))
			fail++
		} else if claimed != nil {
			out.WriteString("PASS ClaimOneDueReminder\n")
			pass++
			err = repository.MarkReminderSent(claimed.ID)
			if err != nil {
				out.WriteString(fmt.Sprintf("FAIL MarkReminderSent: %v\n", err))
				fail++
			} else {
				out.WriteString("PASS MarkReminderSent\n")
				pass++
			}
		} else {
			out.WriteString("INFO ClaimOneDueReminder returned nil (no due reminders)\n")
		}
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
	p := filepath.Join(base, "test_results", "test_repository_result.txt")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte(content), 0644)
}
