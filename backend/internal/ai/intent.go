package ai

import "strings"

// DetectIntent returns a canonical intent from the message body.
// Supports English and Indonesian keywords.
// Returns: create_reminder, list_reminders, store_memory, ignore_smalltalk, or unknown.
//
// Detection order prevents collisions: create_reminder → list_reminders → store_memory → ignore_smalltalk → unknown.
func DetectIntent(body string) string {
	lower := strings.ToLower(strings.TrimSpace(body))

	// 1. create_reminder
	if strings.Contains(lower, "remind") || strings.Contains(lower, "reminder") ||
		strings.Contains(lower, "ingatkan") || strings.Contains(lower, "pengingat") {
		return "create_reminder"
	}

	// 2. list_reminders
	if strings.Contains(lower, "list") || strings.Contains(lower, "show reminders") ||
		strings.Contains(lower, "pending tasks") || strings.Contains(lower, "daftar") ||
		strings.Contains(lower, "tampilkan") || strings.Contains(lower, "tugas saya") {
		return "list_reminders"
	}

	// 3. store_memory
	if strings.Contains(lower, "remember that") || strings.Contains(lower, "remember my") ||
		strings.Contains(lower, "ingat bahwa") || strings.Contains(lower, "catat bahwa") {
		return "store_memory"
	}

	// 4. ignore_smalltalk (exact match to avoid "hi" in "this")
	switch lower {
	case "hi", "hello", "thanks", "halo", "hai", "makasih", "terima kasih":
		return "ignore_smalltalk"
	}

	return "unknown"
}
