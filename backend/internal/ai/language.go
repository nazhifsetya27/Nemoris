package ai

import "strings"

// DetectLanguage returns a language code from the message body.
// Uses lightweight keyword heuristic. Returns: id, en.
// Default fallback: en.
func DetectLanguage(body string) string {
	lower := strings.ToLower(strings.TrimSpace(body))

	// Indonesian indicators
	indonesianKeywords := []string{"ingatkan", "besok", "jam", "hari"}
	for _, kw := range indonesianKeywords {
		if strings.Contains(lower, kw) {
			return "id"
		}
	}

	return "en"
}
