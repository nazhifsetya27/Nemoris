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

func ParseReminder(body string) (string, string, time.Time, bool) {
	lang := detectReminderLanguage(body)
	if lang == "" {
		return "", "", time.Time{}, false
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

	if strings.Contains(contentLower, timeSep) {
		idx := strings.Index(contentLower, timeSep)
		task := strings.TrimSpace(content[:idx])
		timePart := content[idx+len(timeSep):]
		rawTime := rawTimePrefix + timePart

		parsedTime, ok := parseTomorrowTime(timePart, lang)
		if !ok {
			return task, rawTime, time.Time{}, false
		}
		return task, rawTime, parsedTime, true
	}

	return content, "", time.Time{}, true
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
