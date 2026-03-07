package ai

import "time"

// ParseResult holds the canonical structured contract shared by AI and rule parsers.
type ParseResult struct {
	Intent string
	Lang   string
	Task   string
	Time   string
}

// Parse is the unified parser entry point.
// Always fills Intent and Lang.
// For create_reminder only: uses ParseReminder to fill Task and Time (RFC3339).
// If reminder parse fails, Task and Time stay empty.
func Parse(body string) ParseResult {
	intent := DetectIntent(body)
	lang := DetectLanguage(body)
	result := ParseResult{
		Intent: intent,
		Lang:   lang,
		Task:   "",
		Time:   "",
	}

	if intent == "create_reminder" {
		task, _, parsedTime, ok := ParseReminder(body)
		if ok {
			result.Task = task
			if !parsedTime.IsZero() {
				result.Time = parsedTime.Format(time.RFC3339)
			}
		}
	}

	return result
}
