package ai

import "strings"

// DetectIntent returns a canonical intent from the message body.
// Supports English (remind, list) and Indonesian (ingatkan, daftar).
// Returns: create_reminder, list_reminders, or unknown.
func DetectIntent(body string) string {
	lower := strings.ToLower(strings.TrimSpace(body))

	// Create reminder: remind (en), ingatkan (id)
	if strings.Contains(lower, "remind") || strings.Contains(lower, "ingatkan") {
		return "create_reminder"
	}

	// List reminders: list (en), daftar (id)
	if strings.Contains(lower, "list") || strings.Contains(lower, "daftar") {
		return "list_reminders"
	}

	return "unknown"
}
