package ai

import (
	"strconv"
	"strings"
	"time"
)

func detectReminderLanguage(body string) string {
	lower := strings.ToLower(body)
	if strings.HasPrefix(lower, "ingatkan saya untuk ") {
		return "id"
	}
	if strings.HasPrefix(lower, "remind me to ") {
		return "en"
	}
	return ""
}

// ambiguousRecurrenceModifiers: presence of any blocks recurrence (e.g. "sometimes every month", "maybe every week", "tiap bulan mungkin", "kadang setiap minggu").
var ambiguousRecurrenceModifiers = []string{
	"sometimes", "maybe",
	"mungkin", "kadang",
}

func isAmbiguousRecurrence(contentLower string) bool {
	for _, m := range ambiguousRecurrenceModifiers {
		if strings.Contains(contentLower, m) {
			return true
		}
	}
	return false
}

// detectRecurrenceType returns canonical recurrence type from explicit phrase match, or "" if none.
// Ambiguous content (e.g. "sometimes every month") returns "" for safety.
// English: every day, every week, every month.
// Indonesian: setiap hari, setiap minggu, setiap bulan.
func detectRecurrenceType(contentLower string) string {
	if isAmbiguousRecurrence(contentLower) {
		return ""
	}
	phrases := map[string]string{
		"every day":     "daily",
		"every week":   "weekly",
		"every month":  "monthly",
		"setiap hari":  "daily",
		"setiap minggu": "weekly",
		"setiap bulan":  "monthly",
	}
	for phrase, canonical := range phrases {
		if strings.Contains(contentLower, phrase) {
			return canonical
		}
	}
	return ""
}

func ParseReminder(body string) (string, string, time.Time, string, bool) {
	lang := detectReminderLanguage(body)
	if lang == "" {
		return "", "", time.Time{}, "", false
	}

	var prefix, timeSep, rawTimePrefix string
	if lang == "en" {
		prefix = "remind me to "
		timeSep = " tomorrow "
		rawTimePrefix = "tomorrow "
	} else {
		prefix = "ingatkan saya untuk "
		timeSep = " besok "
		rawTimePrefix = "besok "
	}

	content := body[len(prefix):]
	contentLower := strings.ToLower(content)
	recurrenceType := detectRecurrenceType(contentLower)

	if strings.Contains(contentLower, timeSep) {
		idx := strings.Index(contentLower, timeSep)
		task := strings.TrimSpace(content[:idx])
		timePart := content[idx+len(timeSep):]
		rawTime := rawTimePrefix + timePart

		parsedTime, ok := parseTomorrowTime(timePart, lang)
		if !ok {
			return task, rawTime, time.Time{}, recurrenceType, false
		}
		return task, rawTime, parsedTime, recurrenceType, true
	}

	return content, "", time.Time{}, recurrenceType, true
}

func parseTomorrowTime(raw string, lang string) (time.Time, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))

	if lang == "id" {
		raw = strings.TrimPrefix(raw, "jam ")
		raw = strings.ReplaceAll(raw, " malam", "pm")
		raw = strings.ReplaceAll(raw, " pagi", "am")
	}

	raw = strings.TrimPrefix(raw, "at ")

	var hour int
	var err error

	if strings.HasSuffix(raw, "pm") {
		hourStr := strings.TrimSuffix(raw, "pm")
		hour, err = strconv.Atoi(strings.TrimSpace(hourStr))
		if err != nil {
			return time.Time{}, false
		}

		if hour < 12 {
			hour += 12
		}
	} else if strings.HasSuffix(raw, "am") {
		hourStr := strings.TrimSuffix(raw, "am")
		hour, err = strconv.Atoi(strings.TrimSpace(hourStr))
		if err != nil {
			return time.Time{}, false
		}
	} else {
		return time.Time{}, false
	}

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	result := time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		hour,
		0,
		0,
		0,
		now.Location(),
	)

	return result, true
}
