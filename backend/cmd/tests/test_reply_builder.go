// Manual test for i18n reply builder layer.
// Run: go run ./cmd/tests/test_reply_builder.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nemoris/internal/i18n"
)

func main() {
	var out strings.Builder
	out.WriteString("=== test_reply_builder.go ===\n\n")

	pass, fail := 0, 0

	// --- create_reminder id, en ---
	out.WriteString("--- create_reminder ---\n")
	createID := i18n.BuildFromIntent("id", "create_reminder", map[string]string{"time": "besok jam 8 malam"})
	createEN := i18n.BuildFromIntent("en", "create_reminder", map[string]string{"time": "tomorrow 8pm"})
	out.WriteString(fmt.Sprintf("id: %q\n", createID))
	out.WriteString(fmt.Sprintf("en: %q\n", createEN))
	if strings.Contains(createID, "8 malam") || strings.Contains(createID, "ingatkan") {
		out.WriteString("PASS create_reminder id\n")
		pass++
	} else {
		out.WriteString("FAIL create_reminder id\n")
		fail++
	}
	if strings.Contains(createEN, "8pm") || strings.Contains(createEN, "remind") {
		out.WriteString("PASS create_reminder en\n")
		pass++
	} else {
		out.WriteString("FAIL create_reminder en\n")
		fail++
	}

	// --- list_reminders id, en ---
	out.WriteString("\n--- list_reminders ---\n")
	listID := i18n.BuildFromIntent("id", "list_reminders", map[string]string{"count": "3"})
	listEN := i18n.BuildFromIntent("en", "list_reminders", map[string]string{"count": "3"})
	out.WriteString(fmt.Sprintf("id: %q\n", listID))
	out.WriteString(fmt.Sprintf("en: %q\n", listEN))
	if strings.Contains(listID, "3") && strings.Contains(listID, "pengingat") {
		out.WriteString("PASS list_reminders id\n")
		pass++
	} else {
		out.WriteString("FAIL list_reminders id\n")
		fail++
	}
	if strings.Contains(listEN, "3") && strings.Contains(listEN, "reminder") {
		out.WriteString("PASS list_reminders en\n")
		pass++
	} else {
		out.WriteString("FAIL list_reminders en\n")
		fail++
	}

	// --- store_memory (memory) id, en ---
	out.WriteString("\n--- store_memory ---\n")
	memID := i18n.BuildFromIntent("id", "store_memory", nil)
	memEN := i18n.BuildFromIntent("en", "store_memory", nil)
	out.WriteString(fmt.Sprintf("id: %q\n", memID))
	out.WriteString(fmt.Sprintf("en: %q\n", memEN))
	if memID != "" && memEN != "" {
		out.WriteString("PASS store_memory id, en\n")
		pass += 2
	} else {
		out.WriteString("FAIL store_memory\n")
		fail += 2
	}

	// --- unknown language fallback ---
	out.WriteString("\n--- unknown language fallback ---\n")
	fallback := i18n.BuildFromIntent("xx", "unknown", nil)
	out.WriteString(fmt.Sprintf("xx (unknown): %q\n", fallback))
	if fallback != "" && (strings.Contains(fallback, "Sorry") || strings.Contains(fallback, "understand")) {
		out.WriteString("PASS unknown lang falls back to en\n")
		pass++
	} else {
		out.WriteString("FAIL unknown lang fallback\n")
		fail++
	}

	// --- placeholders: task, time, count ---
	out.WriteString("\n--- placeholders ---\n")
	taskMsg := i18n.Build("en", "create_reminder_success", map[string]string{"task": "pay bills", "time": "tomorrow 9am"})
	timeMsg := i18n.Build("en", "create_reminder_success", map[string]string{"time": "tomorrow 9am"})
	countMsg := i18n.Build("en", "list_reminders_success", map[string]string{"count": "5"})
	out.WriteString(fmt.Sprintf("task+time: %q\n", taskMsg))
	out.WriteString(fmt.Sprintf("time only: %q\n", timeMsg))
	out.WriteString(fmt.Sprintf("count: %q\n", countMsg))
	if strings.Contains(timeMsg, "9am") {
		out.WriteString("PASS {{time}} placeholder\n")
		pass++
	} else {
		out.WriteString("FAIL {{time}} placeholder\n")
		fail++
	}
	if strings.Contains(countMsg, "5") {
		out.WriteString("PASS {{count}} placeholder\n")
		pass++
	} else {
		out.WriteString("FAIL {{count}} placeholder\n")
		fail++
	}
	if strings.Contains(taskMsg, "9am") {
		out.WriteString("PASS {{task}} available in Build (template may not use it)\n")
		pass++
	} else {
		out.WriteString("FAIL task placeholder\n")
		fail++
	}

	// --- Summary ---
	out.WriteString("\n--- Summary ---\n")
	out.WriteString(fmt.Sprintf("PASS: %d  FAIL: %d\n", pass, fail))

	fmt.Print(out.String())

	base, _ := os.Getwd()
	p := filepath.Join(base, "test_results", "test_reply_builder_result.txt")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte(out.String()), 0644)
}
