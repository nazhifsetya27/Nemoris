// Manual test for reminder flow: ProcessMessage, duplicate, self, valid reminder.
// Run: go run ./cmd/tests/test_reminder_flow.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nemoris/internal/config"
	"nemoris/internal/database"
	"nemoris/internal/service"
)

const testSender = "test-user@s.whatsapp.net"

func main() {
	config.Load()
	database.Init()

	var out strings.Builder
	out.WriteString("=== test_reminder_flow.go ===\n\n")

	pass, fail := 0, 0

	// --- Self message ---
	out.WriteString("--- Self message ---\n")
	selfResp := service.ProcessMessage(config.App.BotNumber, "remind me to self test tomorrow 8pm")
	if selfResp == "" {
		out.WriteString("PASS self message ignored (empty response)\n")
		pass++
	} else {
		out.WriteString(fmt.Sprintf("FAIL self message: expected empty, got %q\n", selfResp))
		fail++
	}

	// --- Duplicate message ---
	out.WriteString("\n--- Duplicate message ---\n")
	dupBody := "remind me to duplicate test tomorrow 8pm"
	first := service.ProcessMessage(testSender, dupBody)
	second := service.ProcessMessage(testSender, dupBody)
	if second == "duplicate ignored" {
		out.WriteString("PASS duplicate returns 'duplicate ignored'\n")
		pass++
	} else {
		out.WriteString(fmt.Sprintf("FAIL duplicate: got %q, want 'duplicate ignored'\n", second))
		fail++
	}
	out.WriteString(fmt.Sprintf("  first response: %q\n", first))

	// --- Valid reminder ---
	out.WriteString("\n--- Valid reminder ---\n")
	validBody := "remind me to valid flow test tomorrow 9pm"
	resp := service.ProcessMessage(testSender, validBody)
	if resp != "" && resp != "internal error" && !strings.Contains(resp, "specify") {
		out.WriteString("PASS valid reminder: localized reply returned\n")
		out.WriteString(fmt.Sprintf("  reply: %q\n", resp))
		pass++
	} else {
		out.WriteString(fmt.Sprintf("FAIL valid reminder: got %q\n", resp))
		fail++
	}

	// --- Summary ---
	out.WriteString("\n--- Summary ---\n")
	out.WriteString(fmt.Sprintf("PASS: %d  FAIL: %d\n", pass, fail))

	fmt.Print(out.String())

	base, _ := os.Getwd()
	p := filepath.Join(base, "test_results", "test_reminder_flow_result.txt")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte(out.String()), 0644)
}
