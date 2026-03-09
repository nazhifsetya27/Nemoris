package provider

// Result holds the canonical structured contract for external AI responses.
// Matches the shape expected by the parser layer: intent, task, time, lang.
// Contract: explicit typed struct only. No raw prose. No map[string]interface{}.
type Result struct {
	Intent         string
	Task           string
	Time           string
	Lang           string
	RecurrenceType string
}

// ZeroResult returns the safe zero-value canonical struct for unsupported or empty provider results.
func ZeroResult() Result {
	return Result{
		Intent:         "unknown",
		Lang:           "en",
		Task:           "",
		Time:           "",
		RecurrenceType: "",
	}
}

// PrepareExternalCall prepares an external AI call and returns a stub-safe result.
// No real provider integration yet. No HTTP client, no SDK.
// Always returns typed Result. Future: will accept body, call external API, parse into Result.
func PrepareExternalCall(body string) (Result, error) {
	return ZeroResult(), nil
}
