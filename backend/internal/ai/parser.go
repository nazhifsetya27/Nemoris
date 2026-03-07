package ai

// ParseResult holds the result of parsing a message body.
type ParseResult struct {
	Intent string
	Lang   string
	Task   string
	Time   string
}

// Parse is the unified parser entry point.
// Calls DetectIntent and DetectLanguage.
// Task and Time are empty for now (no advanced extraction yet).
func Parse(body string) ParseResult {
	return ParseResult{
		Intent: DetectIntent(body),
		Lang:   DetectLanguage(body),
		Task:   "",
		Time:   "",
	}
}
