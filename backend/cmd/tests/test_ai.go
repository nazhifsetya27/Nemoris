// Manual test for AI layer: DetectLanguage, DetectIntent, ParseReminder.
// Run: go run ./cmd/tests/test_ai.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nemoris/internal/ai"
)

func main() {
	var out strings.Builder
	out.WriteString("=== test_ai.go ===\n\n")

	pass, fail := 0, 0

	// --- DetectLanguage ---
	out.WriteString("--- DetectLanguage ---\n")
	cases := []struct {
		body   string
		expect string
	}{
		{"remind me to pay electricity tomorrow 8pm", "en"},
		{"ingatkan saya untuk bayar listrik besok jam 8 malam", "id"}, // "untuk" required by parser
		{"remind aku tomorrow jam 8", "id"}, // "jam" triggers id
		{"hello bro", "en"},
	}
	for _, c := range cases {
		got := ai.DetectLanguage(c.body)
		if got == c.expect {
			out.WriteString(fmt.Sprintf("PASS DetectLanguage: %q -> %s\n", c.body, got))
			pass++
		} else {
			out.WriteString(fmt.Sprintf("FAIL DetectLanguage: %q -> got %s, want %s\n", c.body, got, c.expect))
			fail++
		}
	}

	// --- DetectIntent ---
	out.WriteString("\n--- DetectIntent ---\n")
	intentCases := []struct {
		body   string
		expect string
	}{
		{"remind me to pay electricity tomorrow 8pm", "create_reminder"},
		{"ingatkan saya untuk bayar listrik besok jam 8 malam", "create_reminder"},
		{"remind aku tomorrow jam 8", "create_reminder"},
		{"hello bro", "unknown"},
		{"list my tasks", "list_reminders"},     // "list" without "reminder"
		{"tampilkan tugas saya", "list_reminders"},
		{"remember that I like coffee", "store_memory"},
		{"hi", "ignore_smalltalk"},
	}
	for _, c := range intentCases {
		got := ai.DetectIntent(c.body)
		if got == c.expect {
			out.WriteString(fmt.Sprintf("PASS DetectIntent: %q -> %s\n", c.body, got))
			pass++
		} else {
			out.WriteString(fmt.Sprintf("FAIL DetectIntent: %q -> got %s, want %s\n", c.body, got, c.expect))
			fail++
		}
	}

	// --- ParseReminder ---
	out.WriteString("\n--- ParseReminder ---\n")
	parseCases := []struct {
		body   string
		expectOK bool
		desc   string
	}{
		{"remind me to pay electricity tomorrow 8pm", true, "English"},
		{"ingatkan saya untuk bayar listrik besok jam 8 malam", true, "Indonesian"},
		{"remind aku tomorrow jam 8", false, "Mixed (no 'me to' prefix)"},
		{"hello bro", false, "Unknown"},
	}
	for _, c := range parseCases {
		task, rawTime, remindAt, ok := ai.ParseReminder(c.body)
		if ok == c.expectOK {
			out.WriteString(fmt.Sprintf("PASS ParseReminder %s: %q -> ok=%v task=%q rawTime=%q\n", c.desc, c.body, ok, task, rawTime))
			if ok && !remindAt.IsZero() {
				out.WriteString(fmt.Sprintf("       remindAt=%s\n", remindAt.Format("2006-01-02 15:04")))
			}
			pass++
		} else {
			out.WriteString(fmt.Sprintf("FAIL ParseReminder %s: %q -> got ok=%v, want %v\n", c.desc, c.body, ok, c.expectOK))
			fail++
		}
	}

	// --- Summary ---
	out.WriteString("\n--- Summary ---\n")
	out.WriteString(fmt.Sprintf("PASS: %d  FAIL: %d\n", pass, fail))

	fmt.Print(out.String())

	// Write to result file (run from backend/ so test_results/ is backend/test_results/)
	base, _ := os.Getwd()
	p := filepath.Join(base, "test_results", "test_ai_result.txt")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte(out.String()), 0644)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err == nil {
		_ = os.WriteFile(p, []byte(out.String()), 0644)
	}
}
