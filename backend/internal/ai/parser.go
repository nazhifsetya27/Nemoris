package ai

import (
	"time"

	"nemoris/internal/ai/provider"
)

// ParseResult holds the canonical structured contract shared by AI and rule parsers.
type ParseResult struct {
	Intent         string
	Lang           string
	Task           string
	Time           string
	RecurrenceType string
}

// Parse is the unified parser entry point.
// Dual parsing strategy: rule parser first, fallback second.
// Step 1: Detect intent and language.
// Step 2: Build default ParseResult.
// Step 3: For create_reminder, run rule parser; if success, return immediately.
// Step 4: If rule parse fails, call fallbackParse (placeholder).
func Parse(body string) ParseResult {
	// Step 1: Detect intent and language
	intent := DetectIntent(body)
	lang := DetectLanguage(body)

	// Step 2: Build default ParseResult
	result := ParseResult{
		Intent: intent,
		Lang:   lang,
		Task:   "",
		Time:   "",
	}

	// Step 3: Rule parser for create_reminder
	if intent == "create_reminder" {
		task, _, parsedTime, recurrenceType, ok := ParseReminder(body)
		if ok {
			result.Task = task
			if !parsedTime.IsZero() {
				result.Time = parsedTime.Format(time.RFC3339)
			}
			result.RecurrenceType = recurrenceType
			return result
		}
	}

	// Step 4: Rule parse failed — use fallback (placeholder)
	return fallbackParse(body, intent, lang)
}

// fallbackParse is the fallback path when rule parsing fails.
// For unknown intent only: calls provider stub. Known intents keep rule result.
func fallbackParse(body string, intent string, lang string) ParseResult {
	if intent == "unknown" {
		pr, err := provider.PrepareExternalCall(body)
		if err != nil {
			return ParseResult{Intent: "unknown", Lang: lang, Task: "", Time: ""}
		}
		return ParseResult{
			Intent:         pr.Intent,
			Lang:          pr.Lang,
			Task:          pr.Task,
			Time:          pr.Time,
			RecurrenceType: pr.RecurrenceType,
		}
	}
	return ParseResult{
		Intent: intent,
		Lang:   lang,
		Task:   "",
		Time:   "",
	}
}
