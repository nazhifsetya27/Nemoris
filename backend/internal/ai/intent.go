package ai

import "strings"

// DetectIntent returns a canonical intent from the message body.
// Supports English and Indonesian keywords.
// Returns: create_reminder, list_reminders, store_memory, retrieve_memory, ignore_smalltalk, or unknown.
//
// Detection order prevents collisions: create_reminder → list_reminders → store_memory → retrieve_memory → ignore_smalltalk → unknown.
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

	// 4. retrieve_memory
	if strings.Contains(lower, "what did i tell") || strings.Contains(lower, "what did i say") ||
		strings.Contains(lower, "what is my") || strings.Contains(lower, "what do you remember") ||
		strings.Contains(lower, "what did you remember") ||
		strings.Contains(lower, "apa yang saya bilang") || strings.Contains(lower, "apa yang aku bilang") ||
		strings.Contains(lower, "apa yang saya katakan") || strings.Contains(lower, "apa yang aku katakan") ||
		strings.Contains(lower, "berapa nomor rekening") {
		return "retrieve_memory"
	}

	// 5. ignore_smalltalk (exact match to avoid "hi" in "this")
	switch lower {
	case "hi", "hello", "thanks", "halo", "hai", "makasih", "terima kasih":
		return "ignore_smalltalk"
	}

	return "unknown"
}
