package provider

// Result holds the canonical structured contract for external AI responses.
// Matches the shape expected by the parser layer: intent, task, time, lang.
type Result struct {
	Intent         string
	Task           string
	Time           string
	Lang           string
	RecurrenceType string
}

// PrepareExternalCall prepares an external AI call and returns a stub-safe result.
// No real provider integration yet. No HTTP client, no SDK.
// Future: will accept body, call external API, parse response into Result.
func PrepareExternalCall(body string) (Result, error) {
	return Result{
		Intent:         "unknown",
		Task:           "",
		Time:           "",
		Lang:           "en",
		RecurrenceType: "",
	}, nil
}
