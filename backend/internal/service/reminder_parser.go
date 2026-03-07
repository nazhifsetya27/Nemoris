package service

import (
	"strconv"
	"strings"
	"time"
)

func ParseReminder(body string) (string, string, time.Time, bool) {
	lower := strings.ToLower(body)

	if !strings.HasPrefix(lower, "remind me to ") {
		return "", "", time.Time{}, false
	}

	content := strings.TrimPrefix(body, "remind me to ")

	if strings.Contains(content, " tomorrow ") {
		parts := strings.SplitN(content, " tomorrow ", 2)

		task := parts[0]
		rawTime := "tomorrow " + parts[1]

		parsedTime, ok := parseTomorrowTime(parts[1])
		if !ok {
			return task, rawTime, time.Time{}, false
		}

		return task, rawTime, parsedTime, true
	}

	return content, "", time.Time{}, true
}

func parseTomorrowTime(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))

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